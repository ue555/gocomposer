package cli

// Port of the Symfony Console 5.4 (v5.4.54, as locked by composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80; MIT, (c) Fabien Potencier) parts
// Composer's list/help output is made of: Command\ListCommand,
// Command\HelpCommand, Descriptor\ApplicationDescription, Descriptor\
// TextDescriptor, JsonDescriptor, XmlDescriptor, MarkdownDescriptor,
// Helper\DescriptorHelper and Application::find/findNamespace, together with
// Composer\Console\Application::getHelp/getLongVersion (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi
// Boggiano).
//
// By design gocomposer does not print Composer's ASCII logo and version
// line: its application name is "gocomposer" with its own version, and the
// banner of "list" is gocomposer's version line (versionLine).

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/platform"
)

// appName and appVersion are Application::getName/getVersion.
const appName = "gocomposer"

func appVersion() string { return Version }

func listCommand() *Command {
	return &Command{Name: "list", Run: func(ctx *Context) (int, error) {
		format, ok := ctx.In.String("format")
		if !ok {
			format = "txt"
		}
		return 0, ctx.App.describe(ctx.IO, nil, descOptions{
			format:    format,
			rawText:   ctx.In.Bool("raw"),
			namespace: ctx.In.Arg("namespace"),
			short:     ctx.In.Bool("short"),
		})
	}}
}

func helpCommand() *Command {
	return &Command{Name: "help", IgnoreValidationErrors: true, Run: func(ctx *Context) (int, error) {
		name := ctx.In.Arg("command_name")
		if len(ctx.In.Args("command_name")) == 0 {
			name = "help"
		}
		cmd, err := ctx.App.findCommand(name)
		if err != nil {
			return 1, err
		}
		return 0, ctx.App.describe(ctx.IO, cmd, helpOptions(ctx.In))
	}}
}

func helpOptions(in *Input) descOptions {
	format, ok := in.String("format")
	if !ok {
		format = "txt"
	}
	return descOptions{format: format, rawText: in.Bool("raw")}
}

// renderHelp describes cmd as text (HelpCommand with the default options).
func (a *App) renderHelp(io *cio.IO, cmd *Command) int {
	if err := a.describe(io, cmd, descOptions{format: "txt"}); err != nil {
		renderError(io, err)
		return 1
	}
	return 0
}

// hasParameterOption is ArgvInput::hasParameterOption.
func hasParameterOption(argv []string, onlyParams bool, values ...string) bool {
	for _, t := range argv {
		if onlyParams && t == "--" {
			return false
		}
		for _, v := range values {
			leading := v
			if strings.HasPrefix(v, "--") {
				leading = v + "="
			}
			if t == v || leading != "" && strings.HasPrefix(t, leading) {
				return true
			}
		}
	}
	return false
}

// runHelpFlag is Application::doRun with --help/-h: without a command name
// the help of the default (list) command is shown with default options,
// otherwise the help command runs for cmd with the original input bound to
// its definition (validation errors ignored).
func (a *App) runHelpFlag(io *cio.IO, argv []string, cmd *Command, noCommand bool, cmdToken string) (int, error) {
	help := a.byName["help"]
	a.running = help
	if noCommand {
		return 0, a.describe(io, a.byName["list"], descOptions{format: "txt"})
	}
	in, err := parseInput(help, argv, cmdToken)
	if err != nil {
		in, _ = parseInputLenient(help, argv, cmdToken)
	}
	return 0, a.describe(io, cmd, helpOptions(in))
}

// printVersion is the --version/-V output of Application::doRun plus the
// PHP version lines Composer's doRun writes to stderr (for gocomposer: the
// PHP binary it runs scripts and plugins with).
func printVersion(io *cio.IO) {
	io.Write(versionLine())
	if php, err := platform.FindPHP(); err == nil {
		if f, err := platform.Probe(php); err == nil {
			io.WriteError(fmt.Sprintf("<info>PHP</info> version <comment>%s</comment> (%s)", f.PHPVersion, f.PHPBinary))
		}
	}
	io.WriteError(`Run the "diagnose" command to get more detailed diagnostics output.`)
}

// ---------------------------------------------------------------------------
// Application::find / findNamespace / getNamespaces / all

// commandEntries is Application::all() (allEntries) completed with any
// registered name that is not reachable from a.commands.
func (a *App) commandEntries() []compEntry {
	entries := a.allEntries()
	seen := map[string]bool{}
	for _, e := range entries {
		seen[e.key] = true
	}
	var extra []string
	for k := range a.byName {
		if !seen[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	for _, k := range extra {
		entries = append(entries, compEntry{k, a.byName[k]})
	}
	return entries
}

// find is Application::find.
func (a *App) find(name string) (*Command, error) {
	if c, ok := a.byName[name]; ok {
		return c, nil
	}
	entries := a.commandEntries()
	cmdOf := map[string]*Command{}
	var all []string
	for _, e := range entries {
		all = append(all, e.key)
		cmdOf[e.key] = e.cmd
	}
	expr := abbrevExpr(name)
	re := regexp.MustCompile(`^` + expr)
	commands := grep(re, all)
	if len(commands) == 0 {
		commands = grep(regexp.MustCompile(`(?i)^`+expr), all)
	}
	if len(commands) == 0 || len(grep(regexp.MustCompile(`(?i)^`+expr+`$`), commands)) < 1 {
		if pos := strings.LastIndex(name, ":"); pos >= 0 {
			if _, err := a.findNamespaceName(name[:pos]); err != nil {
				return nil, err
			}
		}
		return nil, &CLIError{Msg: a.notDefinedMessage(name)}
	}
	if len(commands) > 1 {
		in := map[string]bool{}
		for _, c := range commands {
			in[c] = true
		}
		var kept []string
		seen := map[string]bool{}
		for _, k := range commands {
			cn := cmdOf[k].Name
			if (cn == k || !in[cn]) && !seen[k] {
				kept = append(kept, k)
				seen[k] = true
			}
		}
		commands = kept
	}
	if len(commands) > 1 {
		usable := cio.TerminalWidth() - 10
		maxLen := 0
		for _, k := range commands {
			if w := strWidth(k); w > maxLen {
				maxLen = w
			}
		}
		var abbrevs, visible []string
		for _, k := range commands {
			c := cmdOf[k]
			if c.Hidden {
				continue
			}
			visible = append(visible, k)
			ab := k + strings.Repeat(" ", max(0, maxLen-len(k))) + " " + c.Description
			if strWidth(ab) > usable {
				ab = mbSubstr(ab, usable-3) + "..."
			}
			abbrevs = append(abbrevs, ab)
		}
		commands = visible
		if len(commands) > 1 {
			return nil, &CLIError{Msg: fmt.Sprintf("Command \"%s\" is ambiguous.\nDid you mean one of these?\n%s.", name, "    "+strings.Join(abbrevs, "\n    "))}
		}
	}
	if len(commands) == 0 {
		return nil, &CLIError{Msg: a.notDefinedMessage(name)}
	}
	c := cmdOf[commands[0]]
	if c.Hidden {
		return nil, &CLIError{Msg: fmt.Sprintf(`The command "%s" does not exist.`, name)}
	}
	return c, nil
}

func abbrevExpr(name string) string {
	var q []string
	for _, p := range strings.Split(name, ":") {
		q = append(q, regexp.QuoteMeta(p))
	}
	return strings.Join(q, `[^:]*:`) + `[^:]*`
}

func grep(re *regexp.Regexp, list []string) []string {
	var out []string
	for _, s := range list {
		if re.MatchString(s) {
			out = append(out, s)
		}
	}
	return out
}

func mbSubstr(s string, n int) string {
	if n < 0 {
		n = 0
	}
	r := []rune(s)
	if n > len(r) {
		n = len(r)
	}
	return string(r[:n])
}

// extractNamespace is Application::extractNamespace.
func extractNamespace(name string, limit int) string {
	parts := strings.Split(name, ":")
	parts = parts[:len(parts)-1]
	if limit >= 0 && limit < len(parts) {
		parts = parts[:limit]
	}
	return strings.Join(parts, ":")
}

func extractAllNamespaces(name string) []string {
	parts := strings.Split(name, ":")
	parts = parts[:len(parts)-1]
	var out []string
	for _, p := range parts {
		if len(out) > 0 {
			out = append(out, out[len(out)-1]+":"+p)
		} else {
			out = append(out, p)
		}
	}
	return out
}

// getNamespaces is Application::getNamespaces.
func (a *App) getNamespaces() []string {
	var out []string
	seen := map[string]bool{}
	add := func(list []string) {
		for _, n := range list {
			if n != "" && !seen[n] {
				seen[n] = true
				out = append(out, n)
			}
		}
	}
	for _, e := range a.commandEntries() {
		if e.cmd.Hidden {
			continue
		}
		add(extractAllNamespaces(e.cmd.Name))
		for _, al := range e.cmd.Aliases {
			add(extractAllNamespaces(al))
		}
	}
	return out
}

// findNamespaceName is Application::findNamespace.
func (a *App) findNamespaceName(namespace string) (string, error) {
	all := a.getNamespaces()
	namespaces := grep(regexp.MustCompile(`^`+abbrevExpr(namespace)), all)
	if len(namespaces) == 0 {
		msg := `There are no commands defined in the "` + namespace + `" namespace.`
		if alts := findAlternatives(namespace, all); len(alts) > 0 {
			if len(alts) == 1 {
				msg += "\n\nDid you mean this?\n    "
			} else {
				msg += "\n\nDid you mean one of these?\n    "
			}
			msg += strings.Join(alts, "\n    ")
		}
		return "", &CLIError{Msg: msg}
	}
	exact := false
	for _, n := range namespaces {
		if n == namespace {
			exact = true
		}
	}
	if len(namespaces) > 1 && !exact {
		return "", &CLIError{Msg: fmt.Sprintf("The namespace \"%s\" is ambiguous.\nDid you mean one of these?\n%s.", namespace, "    "+strings.Join(namespaces, "\n    "))}
	}
	if exact {
		return namespace, nil
	}
	return namespaces[0], nil
}

// ---------------------------------------------------------------------------
// ApplicationDescription

type nsDescription struct {
	id       string
	commands []string
}

type appDescription struct {
	names      []string // getCommands() keys, in order
	commands   map[string]*Command
	aliases    map[string]*Command
	namespaces []nsDescription
}

func (d *appDescription) command(name string) *Command {
	if c, ok := d.commands[name]; ok {
		return c
	}
	return d.aliases[name]
}

const globalNamespace = "_global"

func (a *App) applicationDescription(namespace string, showHidden bool) (*appDescription, error) {
	entries := a.commandEntries()
	if namespace != "" {
		ns, err := a.findNamespaceName(namespace)
		if err != nil {
			return nil, err
		}
		limit := strings.Count(ns, ":") + 1
		var in []compEntry
		for _, e := range entries {
			if extractNamespace(e.key, limit) == ns {
				in = append(in, e)
			}
		}
		entries = in
	}
	// sortCommands
	var global []compEntry
	nsed := map[string][]compEntry{}
	var nsKeys []string
	for _, e := range entries {
		key := extractNamespace(e.key, 1)
		if key == "" || key == globalNamespace {
			global = append(global, e)
		} else {
			if _, ok := nsed[key]; !ok {
				nsKeys = append(nsKeys, key)
			}
			nsed[key] = append(nsed[key], e)
		}
	}
	ksort := func(list []compEntry) {
		sort.SliceStable(list, func(i, j int) bool { return phpjson.PHPCompareKeys(list[i].key, list[j].key) < 0 })
	}
	type group struct {
		id      string
		entries []compEntry
	}
	var groups []group
	if len(global) > 0 {
		ksort(global)
		groups = append(groups, group{globalNamespace, global})
	}
	sort.Strings(nsKeys)
	for _, k := range nsKeys {
		ksort(nsed[k])
		groups = append(groups, group{k, nsed[k]})
	}
	d := &appDescription{commands: map[string]*Command{}, aliases: map[string]*Command{}}
	for _, g := range groups {
		var names []string
		for _, e := range g.entries {
			if e.cmd.Name == "" || (!showHidden && e.cmd.Hidden) {
				continue
			}
			if e.cmd.Name == e.key {
				d.commands[e.key] = e.cmd
				d.names = append(d.names, e.key)
			} else {
				d.aliases[e.key] = e.cmd
			}
			names = append(names, e.key)
		}
		d.namespaces = append(d.namespaces, nsDescription{g.id, names})
	}
	return d, nil
}

// ---------------------------------------------------------------------------
// definitions

// defOption is an InputOption as the descriptors see it.
type defOption struct {
	OptDef
}

func (o defOption) acceptValue() bool { return o.Kind != OptNone }
func (o defOption) isArray() bool     { return o.Kind == OptArray }

// defaultValue is InputOption::getDefault.
func (o defOption) defaultValue() phpjson.Value {
	if o.Kind == OptNone && !o.Negatable {
		return false
	}
	v := toPHPValue(o.RawDefault)
	if v == nil && o.Kind == OptArray {
		return []phpjson.Value{}
	}
	if v == nil && o.Default != "" {
		return o.Default
	}
	return v
}

func argDefault(a ArgDef) phpjson.Value {
	v := toPHPValue(a.RawDefault)
	if v == nil && a.Array {
		return []phpjson.Value{}
	}
	return v
}

// toPHPValue converts encoding/json values of commanddefs.json (and the
// phpjson values of plugin commands) to phpjson values.
func toPHPValue(v interface{}) phpjson.Value {
	switch t := v.(type) {
	case float64:
		if t == float64(int64(t)) {
			return int64(t)
		}
		return t
	case int:
		return int64(t)
	case []interface{}:
		out := make([]phpjson.Value, len(t))
		for i, e := range t {
			out[i] = toPHPValue(e)
		}
		return out
	case map[string]interface{}:
		m := phpjson.NewMap()
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			m.Set(k, toPHPValue(t[k]))
		}
		return m
	}
	return v
}

// fullOptions is the merged definition of Command::mergeApplicationDefinition:
// the command's options followed by the application's; an application option
// with the name of a command option replaces it in place.
func fullOptions(cmd *Command) []defOption {
	var out []defOption
	pos := map[string]int{}
	for _, o := range cmd.Options {
		pos[o.Name] = len(out)
		out = append(out, defOption{o})
	}
	for _, o := range globalOptions {
		if i, ok := pos[o.Name]; ok {
			out[i] = defOption{o}
			continue
		}
		pos[o.Name] = len(out)
		out = append(out, defOption{o})
	}
	return out
}

func globalDefOptions() []defOption {
	out := make([]defOption, len(globalOptions))
	for i, o := range globalOptions {
		out[i] = defOption{o}
	}
	return out
}

// shortSynopsis is Command::getSynopsis(true).
func shortSynopsis(cmd *Command) string {
	var el []string
	if len(cmd.Options) > 0 {
		el = append(el, "[options]")
	}
	if len(el) > 0 && len(cmd.Args) > 0 {
		el = append(el, "[--]")
	}
	tail := ""
	for _, a := range cmd.Args {
		e := "<" + a.Name + ">"
		if a.Array {
			e += "..."
		}
		if !a.Required {
			e = "[" + e
			tail += "]"
		}
		el = append(el, e)
	}
	return strings.TrimSpace(cmd.Name + " " + strings.Join(el, " ") + tail)
}

// phpSelf is $_SERVER['PHP_SELF'].
func phpSelf() string {
	if len(os.Args) == 0 {
		return ""
	}
	return os.Args[0]
}

// processedHelp is Command::getProcessedHelp.
func processedHelp(cmd *Command) string {
	help := cmd.Help
	if cmd.Name == "completion" && strings.Contains(help, "command-dump.php") {
		// DumpCompletionCommand::configure embeds the running script
		// (commanddefs.json was dumped by oracle/scripts/command-dump.php)
		full := phpSelf()
		if abs, err := filepath.Abs(full); err == nil {
			if real, err := filepath.EvalSymlinks(abs); err == nil {
				full = real
			}
		}
		help = strings.ReplaceAll(help, "/etc/bash_completion.d/command-dump.php", "/etc/bash_completion.d/"+filepath.Base(phpSelf()))
		help = strings.ReplaceAll(help, "/scripts/command-dump.php", full)
	}
	if help == "" {
		help = cmd.Description
	}
	help = strings.ReplaceAll(help, "%command.name%", cmd.Name)
	return strings.ReplaceAll(help, "%command.full_name%", phpSelf()+" "+cmd.Name)
}

// ---------------------------------------------------------------------------
// DescriptorHelper

type descOptions struct {
	format    string
	rawText   bool
	namespace string
	short     bool
}

// describe is DescriptorHelper::describe for a command (cmd != nil) or the
// application.
func (a *App) describe(o *cio.IO, cmd *Command, opts descOptions) error {
	w := &descWriter{io: o}
	switch opts.format {
	case "txt":
		t := &textDescriptor{w: w, app: a}
		if cmd != nil {
			t.describeCommand(cmd, opts)
			return nil
		}
		return t.describeApplication(opts)
	case "json":
		if cmd != nil {
			w.raw(phpjson.MustEncode(jsonCommandData(cmd, opts.short), 0))
			return nil
		}
		return a.describeJSON(w, opts)
	case "xml":
		if cmd != nil {
			w.raw(xmlDocument(xmlCommand(cmd, opts.short)))
			return nil
		}
		return a.describeXML(w, opts)
	case "md":
		decorated := o.Decorated
		o.Decorated = false
		defer func() { o.Decorated = decorated }()
		m := &mdDescriptor{w: w, app: a}
		if cmd != nil {
			m.describeCommand(cmd, opts)
			return nil
		}
		return m.describeApplication(opts)
	}
	return &CLIError{Msg: fmt.Sprintf(`Unsupported format "%s".`, opts.format)}
}

// descWriter is Descriptor::write: formatted (OUTPUT_NORMAL) or raw output.
type descWriter struct{ io *cio.IO }

func (w *descWriter) fmt(s string) { w.io.WriteV(s, false, cio.Normal) }

func (w *descWriter) raw(s string) {
	if w.io.Verbosity <= cio.Quiet {
		return
	}
	io.WriteString(w.io.Out, s)
}

// ---------------------------------------------------------------------------
// TextDescriptor

type textDescriptor struct {
	w   *descWriter
	app *App
}

// writeText: options are passed (opts != nil) where TextDescriptor passes
// $options, so raw_text strips the tags only there.
func (t *textDescriptor) writeText(s string, opts *descOptions) {
	if opts != nil && opts.rawText {
		s = stripTags(s)
	}
	t.w.fmt(s)
}

var reDescNewline = regexp.MustCompile(`\s*[\r\n]\s*`)

func (t *textDescriptor) describeInputArgument(a ArgDef, totalWidth int, opts *descOptions) {
	def := ""
	if d := argDefault(a); d != nil && !isEmptyArray(d) {
		def = "<comment> [default: " + formatDefaultValue(d) + "]</comment>"
	}
	spacing := totalWidth - len(a.Name)
	t.writeText(fmt.Sprintf("  <info>%s</info>  %s%s%s", a.Name, strings.Repeat(" ", max(0, spacing)),
		reDescNewline.ReplaceAllLiteralString(a.Desc, "\n"+strings.Repeat(" ", totalWidth+4)), def), opts)
}

func optionSynopsisText(o defOption) string {
	value := ""
	if o.acceptValue() {
		value = "=" + strings.ToUpper(o.Name)
		if o.Kind == OptOptional {
			value = "[" + value + "]"
		}
	}
	short := "    "
	if o.Shortcut != "" {
		short = "-" + o.Shortcut + ", "
	}
	if o.Negatable {
		return short + "--" + o.Name + "|--no-" + o.Name
	}
	return short + "--" + o.Name + value
}

func (t *textDescriptor) describeInputOption(o defOption, totalWidth int, opts *descOptions) {
	def := ""
	if d := o.defaultValue(); o.acceptValue() && d != nil && !isEmptyArray(d) {
		def = "<comment> [default: " + formatDefaultValue(d) + "]</comment>"
	}
	syn := optionSynopsisText(o)
	multi := ""
	if o.isArray() {
		multi = "<comment> (multiple values allowed)</comment>"
	}
	t.writeText(fmt.Sprintf("  <info>%s</info>  %s%s%s%s", syn, strings.Repeat(" ", max(0, totalWidth-strWidth(syn))),
		reDescNewline.ReplaceAllLiteralString(o.Desc, "\n"+strings.Repeat(" ", totalWidth+4)), def, multi), opts)
}

func calculateTotalWidthForOptions(options []defOption) int {
	total := 0
	for _, o := range options {
		n := 1 + max(strWidth(o.Shortcut), 1) + 4 + strWidth(o.Name)
		if o.Negatable {
			n += 6 + strWidth(o.Name)
		} else if o.acceptValue() {
			v := 1 + strWidth(o.Name)
			if o.Kind == OptOptional {
				v += 2
			}
			n += v
		}
		total = max(total, n)
	}
	return total
}

func (t *textDescriptor) describeInputDefinition(args []ArgDef, options []defOption, opts *descOptions) {
	total := calculateTotalWidthForOptions(options)
	for _, a := range args {
		total = max(total, strWidth(a.Name))
	}
	if len(args) > 0 {
		t.writeText("<comment>Arguments:</comment>", opts)
		t.writeText("\n", nil)
		for _, a := range args {
			t.describeInputArgument(a, total, opts)
			t.writeText("\n", nil)
		}
	}
	if len(args) > 0 && len(options) > 0 {
		t.writeText("\n", nil)
	}
	if len(options) > 0 {
		var later []defOption
		t.writeText("<comment>Options:</comment>", opts)
		for _, o := range options {
			if len(o.Shortcut) > 1 {
				later = append(later, o)
				continue
			}
			t.writeText("\n", nil)
			t.describeInputOption(o, total, opts)
		}
		for _, o := range later {
			t.writeText("\n", nil)
			t.describeInputOption(o, total, opts)
		}
	}
}

func (t *textDescriptor) describeCommand(cmd *Command, o descOptions) {
	opts := &o
	description := cmd.Description
	if description != "" {
		t.writeText("<comment>Description:</comment>", opts)
		t.writeText("\n", nil)
		t.writeText("  "+description, nil)
		t.writeText("\n\n", nil)
	}
	t.writeText("<comment>Usage:</comment>", opts)
	usages := append(append([]string{shortSynopsis(cmd)}, cmd.Aliases...), cmd.Usages...)
	for _, u := range usages {
		t.writeText("\n", nil)
		t.writeText("  "+escapeTags(u), opts)
	}
	t.writeText("\n", nil)
	options := fullOptions(cmd)
	if len(options) > 0 || len(cmd.Args) > 0 {
		t.writeText("\n", nil)
		t.describeInputDefinition(cmd.Args, options, opts)
		t.writeText("\n", nil)
	}
	help := processedHelp(cmd)
	if help != "" && help != description {
		t.writeText("\n", nil)
		t.writeText("<comment>Help:</comment>", opts)
		t.writeText("\n", nil)
		t.writeText("  "+strings.ReplaceAll(help, "\n", "\n  "), opts)
		t.writeText("\n", nil)
	}
}

// columnWidth is TextDescriptor::getColumnWidth for command names.
func columnWidth(names []string) int {
	if len(names) == 0 {
		return 0
	}
	w := 0
	for _, n := range names {
		w = max(w, strWidth(n))
	}
	return w + 2
}

func (t *textDescriptor) describeApplication(o descOptions) error {
	opts := &o
	if o.rawText {
		d, err := t.app.applicationDescription(o.namespace, false)
		if err != nil {
			return err
		}
		var widths []string
		for _, n := range d.names {
			widths = append(widths, n)
			widths = append(widths, d.commands[n].Aliases...)
		}
		width := columnWidth(widths)
		for _, n := range d.names {
			t.writeText(fmt.Sprintf("%-*s %s", width, n, d.commands[n].Description), opts)
			t.writeText("\n", nil)
		}
		return nil
	}
	// Application::getHelp (gocomposer's version line instead of
	// Composer's logo and long version)
	if help := versionLine(); help != "" {
		t.writeText(help+"\n\n", opts)
	}
	t.writeText("<comment>Usage:</comment>\n", opts)
	t.writeText("  command [options] [arguments]\n\n", opts)
	t.describeInputDefinition(nil, globalDefOptions(), opts)
	t.writeText("\n", nil)
	t.writeText("\n", nil)
	d, err := t.app.applicationDescription(o.namespace, false)
	if err != nil {
		return err
	}
	commands := map[string]*Command{}
	for k, v := range d.commands {
		commands[k] = v
	}
	if o.namespace != "" && len(d.namespaces) > 0 {
		for _, n := range d.namespaces[0].commands {
			commands[n] = d.command(n)
		}
	}
	var names []string
	for _, ns := range d.namespaces {
		for _, n := range ns.commands {
			if commands[n] != nil {
				names = append(names, n)
			}
		}
	}
	width := columnWidth(names)
	if o.namespace != "" {
		t.writeText(fmt.Sprintf(`<comment>Available commands for the "%s" namespace:</comment>`, o.namespace), opts)
	} else {
		t.writeText("<comment>Available commands:</comment>", opts)
	}
	for _, ns := range d.namespaces {
		var list []string
		for _, n := range ns.commands {
			if commands[n] != nil {
				list = append(list, n)
			}
		}
		if len(list) == 0 {
			continue
		}
		if o.namespace == "" && ns.id != globalNamespace {
			t.writeText("\n", nil)
			t.writeText(" <comment>"+ns.id+"</comment>", opts)
		}
		for _, n := range list {
			t.writeText("\n", nil)
			c := commands[n]
			aliases := ""
			if n == c.Name && len(c.Aliases) > 0 {
				aliases = "[" + strings.Join(c.Aliases, "|") + "] "
			}
			t.writeText(fmt.Sprintf("  <info>%s</info>%s%s", n, strings.Repeat(" ", max(0, width-strWidth(n))), aliases+c.Description), opts)
		}
	}
	t.writeText("\n", nil)
	return nil
}

func isEmptyArray(v phpjson.Value) bool {
	switch t := v.(type) {
	case []phpjson.Value:
		return len(t) == 0
	case *phpjson.Map:
		return t == nil || t.Len() == 0
	}
	return false
}

// formatDefaultValue is TextDescriptor::formatDefaultValue.
func formatDefaultValue(v phpjson.Value) string {
	esc := func(x phpjson.Value) phpjson.Value {
		if s, ok := x.(string); ok {
			return escapeTags(s)
		}
		return x
	}
	switch t := v.(type) {
	case string:
		v = esc(t)
	case []phpjson.Value:
		out := make([]phpjson.Value, len(t))
		for i, e := range t {
			out[i] = esc(e)
		}
		v = out
	}
	s, err := phpjson.Encode(v, phpjson.UnescapedSlashes|phpjson.UnescapedUnicode)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(s, `\\`, `\`)
}

// stripTags is PHP's strip_tags() for the plain text written by the
// descriptors: "<" followed by whitespace stays, anything else up to the
// matching ">" (outside quotes) is removed.
func stripTags(s string) string {
	var sb strings.Builder
	depth := 0
	inTag := false
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !inTag {
			if c == '<' {
				if i+1 < len(s) && (s[i+1] == ' ' || s[i+1] == '\t' || s[i+1] == '\n' || s[i+1] == '\r' || s[i+1] == '\v' || s[i+1] == '\f') {
					sb.WriteByte(c)
					continue
				}
				if i+1 >= len(s) {
					continue
				}
				inTag = true
				depth = 0
				quote = 0
				continue
			}
			sb.WriteByte(c)
			continue
		}
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '<':
			depth++
		case c == '>':
			if depth > 0 {
				depth--
			} else {
				inTag = false
			}
		}
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// JsonDescriptor

func jsonArgumentData(a ArgDef) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("name", a.Name)
	m.Set("is_required", a.Required)
	m.Set("is_array", a.Array)
	m.Set("description", reDescNewline.ReplaceAllLiteralString(a.Desc, " "))
	m.Set("default", argDefault(a))
	return m
}

func jsonOptionData(o defOption, negated bool) *phpjson.Map {
	m := phpjson.NewMap()
	if negated {
		m.Set("name", "--no-"+o.Name)
		m.Set("shortcut", "")
		m.Set("accept_value", false)
		m.Set("is_value_required", false)
		m.Set("is_multiple", false)
		m.Set("description", `Negate the "--`+o.Name+`" option`)
		m.Set("default", false)
		return m
	}
	m.Set("name", "--"+o.Name)
	short := ""
	if o.Shortcut != "" {
		short = "-" + strings.ReplaceAll(o.Shortcut, "|", "|-")
	}
	m.Set("shortcut", short)
	m.Set("accept_value", o.acceptValue())
	m.Set("is_value_required", o.Kind == OptRequired || o.Kind == OptArray)
	m.Set("is_multiple", o.isArray())
	m.Set("description", reDescNewline.ReplaceAllLiteralString(o.Desc, " "))
	m.Set("default", o.defaultValue())
	return m
}

func jsonCommandData(cmd *Command, short bool) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("name", cmd.Name)
	m.Set("description", cmd.Description)
	strs := func(l []string) []phpjson.Value {
		out := make([]phpjson.Value, len(l))
		for i, s := range l {
			out[i] = s
		}
		return out
	}
	if short {
		m.Set("usage", strs(cmd.Aliases))
	} else {
		m.Set("usage", strs(append(append([]string{synopsisFor(cmd)}, cmd.Usages...), cmd.Aliases...)))
		m.Set("help", processedHelp(cmd))
		def := phpjson.NewMap()
		args := phpjson.NewMap()
		for _, a := range cmd.Args {
			args.Set(a.Name, jsonArgumentData(a))
		}
		options := phpjson.NewMap()
		for _, o := range fullOptions(cmd) {
			options.Set(o.Name, jsonOptionData(o, false))
			if o.Negatable {
				options.Set("no-"+o.Name, jsonOptionData(o, true))
			}
		}
		def.Set("arguments", args)
		def.Set("options", options)
		m.Set("definition", def)
	}
	m.Set("hidden", cmd.Hidden)
	return m
}

func (a *App) describeJSON(w *descWriter, o descOptions) error {
	d, err := a.applicationDescription(o.namespace, true)
	if err != nil {
		return err
	}
	data := phpjson.NewMap()
	app := phpjson.NewMap()
	app.Set("name", appName)
	app.Set("version", appVersion())
	data.Set("application", app)
	var cmds []phpjson.Value
	for _, n := range d.names {
		cmds = append(cmds, jsonCommandData(d.commands[n], o.short))
	}
	if cmds == nil {
		cmds = []phpjson.Value{}
	}
	data.Set("commands", cmds)
	if o.namespace != "" {
		data.Set("namespace", o.namespace)
	} else {
		var nss []phpjson.Value
		for _, ns := range d.namespaces {
			m := phpjson.NewMap()
			m.Set("id", ns.id)
			list := make([]phpjson.Value, len(ns.commands))
			for i, c := range ns.commands {
				list[i] = c
			}
			m.Set("commands", list)
			nss = append(nss, m)
		}
		if nss == nil {
			nss = []phpjson.Value{}
		}
		data.Set("namespaces", nss)
	}
	w.raw(phpjson.MustEncode(data, 0))
	return nil
}

// ---------------------------------------------------------------------------
// XmlDescriptor (DOMDocument::saveXML with formatOutput)

type xmlNode struct {
	name     string
	attrs    [][2]string
	children []*xmlNode
	text     *string // text node when name == ""
}

func xel(name string, attrs ...string) *xmlNode {
	n := &xmlNode{name: name}
	for i := 0; i+1 < len(attrs); i += 2 {
		n.attrs = append(n.attrs, [2]string{attrs[i], attrs[i+1]})
	}
	return n
}

func (n *xmlNode) add(c ...*xmlNode) *xmlNode { n.children = append(n.children, c...); return n }

func (n *xmlNode) addText(s string) *xmlNode {
	n.children = append(n.children, &xmlNode{text: &s})
	return n
}

func (n *xmlNode) setAttr(k, v string) { n.attrs = append(n.attrs, [2]string{k, v}) }

var xmlTextEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\r", "&#13;")
var xmlAttrEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "\n", "&#10;", "\r", "&#13;", "\t", "&#9;")

func (n *xmlNode) write(sb *strings.Builder, level int, format bool) {
	if n.text != nil {
		sb.WriteString(xmlTextEscaper.Replace(*n.text))
		return
	}
	sb.WriteString("<" + n.name)
	for _, a := range n.attrs {
		sb.WriteString(" " + a[0] + `="` + xmlAttrEscaper.Replace(a[1]) + `"`)
	}
	if len(n.children) == 0 {
		sb.WriteString("/>")
		return
	}
	sb.WriteString(">")
	childFormat := format
	for _, c := range n.children {
		if c.text != nil {
			childFormat = false
		}
	}
	if childFormat {
		sb.WriteString("\n")
	}
	for _, c := range n.children {
		if childFormat {
			sb.WriteString(strings.Repeat("  ", level+1))
		}
		c.write(sb, level+1, childFormat)
		if childFormat {
			sb.WriteString("\n")
		}
	}
	if childFormat {
		sb.WriteString(strings.Repeat("  ", level))
	}
	sb.WriteString("</" + n.name + ">")
}

func xmlDocument(roots ...*xmlNode) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	for _, r := range roots {
		r.write(&sb, 0, true)
		sb.WriteString("\n")
	}
	return sb.String()
}

func boolAttr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// xmlDefaults is the $defaults list of the XmlDescriptor.
func xmlDefaults(v phpjson.Value) []string {
	switch t := v.(type) {
	case []phpjson.Value:
		var out []string
		for _, e := range t {
			out = append(out, phpScalarString(e))
		}
		return out
	case *phpjson.Map:
		var out []string
		for _, k := range t.Keys() {
			e, _ := t.Get(k)
			out = append(out, phpScalarString(e))
		}
		return out
	case bool:
		if t {
			return []string{"true"}
		}
		return []string{"false"}
	}
	if phpjson.Truthy(v) {
		return []string{phpScalarString(v)}
	}
	return nil
}

func phpScalarString(v phpjson.Value) string {
	s, _ := phpjson.ScalarString(v)
	return s
}

func xmlArgument(a ArgDef) *xmlNode {
	n := xel("argument", "name", a.Name, "is_required", boolAttr(a.Required), "is_array", boolAttr(a.Array))
	n.add(xel("description").addText(a.Desc))
	defs := xel("defaults")
	for _, d := range xmlDefaults(argDefault(a)) {
		defs.add(xel("default").addText(d))
	}
	return n.add(defs)
}

func xmlOption(o defOption) []*xmlNode {
	n := xel("option", "name", "--"+o.Name)
	if pos := strings.Index(o.Shortcut, "|"); pos >= 0 {
		n.setAttr("shortcut", "-"+o.Shortcut[:pos])
		n.setAttr("shortcuts", "-"+strings.ReplaceAll(o.Shortcut, "|", "|-"))
	} else if o.Shortcut != "" {
		n.setAttr("shortcut", "-"+o.Shortcut)
	} else {
		n.setAttr("shortcut", "")
	}
	n.setAttr("accept_value", boolAttr(o.acceptValue()))
	n.setAttr("is_value_required", boolAttr(o.Kind == OptRequired || o.Kind == OptArray))
	n.setAttr("is_multiple", boolAttr(o.isArray()))
	n.add(xel("description").addText(o.Desc))
	if o.acceptValue() {
		defs := xel("defaults")
		for _, d := range xmlDefaults(o.defaultValue()) {
			defs.add(xel("default").addText(d))
		}
		n.add(defs)
	}
	out := []*xmlNode{n}
	if o.Negatable {
		neg := xel("option", "name", "--no-"+o.Name, "shortcut", "", "accept_value", "0", "is_value_required", "0", "is_multiple", "0")
		neg.add(xel("description").addText(`Negate the "--` + o.Name + `" option`))
		out = append(out, neg)
	}
	return out
}

func xmlCommand(cmd *Command, short bool) *xmlNode {
	n := xel("command", "id", cmd.Name, "name", cmd.Name, "hidden", boolAttr(cmd.Hidden))
	usages := xel("usages")
	n.add(usages)
	n.add(xel("description").addText(strings.ReplaceAll(cmd.Description, "\n", "\n ")))
	if short {
		for _, al := range cmd.Aliases {
			usages.add(xel("usage").addText(al))
		}
		return n
	}
	for _, u := range append(append([]string{synopsisFor(cmd)}, cmd.Aliases...), cmd.Usages...) {
		usages.add(xel("usage").addText(u))
	}
	n.add(xel("help").addText(strings.ReplaceAll(processedHelp(cmd), "\n", "\n ")))
	args := xel("arguments")
	for _, a := range cmd.Args {
		args.add(xmlArgument(a))
	}
	options := xel("options")
	for _, o := range fullOptions(cmd) {
		options.add(xmlOption(o)...)
	}
	return n.add(args, options)
}

func (a *App) describeXML(w *descWriter, o descOptions) error {
	root := xel("symfony", "name", appName, "version", appVersion())
	commands := xel("commands")
	root.add(commands)
	d, err := a.applicationDescription(o.namespace, true)
	if err != nil {
		return err
	}
	if o.namespace != "" {
		commands.setAttr("namespace", o.namespace)
	}
	for _, n := range d.names {
		commands.add(xmlCommand(d.commands[n], o.short))
	}
	if o.namespace == "" {
		nss := xel("namespaces")
		root.add(nss)
		for _, ns := range d.namespaces {
			x := xel("namespace", "id", ns.id)
			for _, c := range ns.commands {
				x.add(xel("command").addText(c))
			}
			nss.add(x)
		}
	}
	w.raw(xmlDocument(root))
	return nil
}

// ---------------------------------------------------------------------------
// MarkdownDescriptor (written formatted on an undecorated output)

type mdDescriptor struct {
	w   *descWriter
	app *App
}

func (m *mdDescriptor) write(s string) { m.w.fmt(s) }

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// varExport is var_export($v, true) with the newlines removed.
func varExport(v phpjson.Value) string {
	var s string
	switch t := v.(type) {
	case nil:
		s = "NULL"
	case bool:
		s = strconv.FormatBool(t)
	case int64:
		s = strconv.FormatInt(t, 10)
	case float64:
		s = phpjson.FormatFloatVarExport(t)
	case string:
		s = "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(t) + "'"
	case []phpjson.Value:
		s = "array (\n"
		for i, e := range t {
			s += "  " + strconv.Itoa(i) + " => " + varExport(e) + ",\n"
		}
		s += ")"
	case *phpjson.Map:
		s = "array (\n"
		for _, k := range t.Keys() {
			e, _ := t.Get(k)
			key := "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(k) + "'"
			if _, ok := phpjson.IntKey(k); ok {
				key = k
			}
			s += "  " + key + " => " + varExport(e) + ",\n"
		}
		s += ")"
	default:
		s = fmt.Sprint(v)
	}
	return strings.ReplaceAll(s, "\n", "")
}

func (m *mdDescriptor) describeInputArgument(a ArgDef) {
	name := a.Name
	if name == "" {
		name = "<none>"
	}
	desc := ""
	if a.Desc != "" {
		desc = reDescNewline.ReplaceAllLiteralString(a.Desc, "\n") + "\n\n"
	}
	m.write("#### `" + name + "`\n\n" + desc +
		"* Is required: " + yesNo(a.Required) + "\n" +
		"* Is array: " + yesNo(a.Array) + "\n" +
		"* Default: `" + varExport(argDefault(a)) + "`")
}

func (m *mdDescriptor) describeInputOption(o defOption) {
	name := "--" + o.Name
	if o.Negatable {
		name += "|--no-" + o.Name
	}
	if o.Shortcut != "" {
		name += "|-" + strings.ReplaceAll(o.Shortcut, "|", "|-")
	}
	desc := ""
	if o.Desc != "" {
		desc = reDescNewline.ReplaceAllLiteralString(o.Desc, "\n") + "\n\n"
	}
	m.write("#### `" + name + "`" + "\n\n" + desc +
		"* Accept value: " + yesNo(o.acceptValue()) + "\n" +
		"* Is value required: " + yesNo(o.Kind == OptRequired || o.Kind == OptArray) + "\n" +
		"* Is multiple: " + yesNo(o.isArray()) + "\n" +
		"* Is negatable: " + yesNo(o.Negatable) + "\n" +
		"* Default: `" + varExport(o.defaultValue()) + "`")
}

func (m *mdDescriptor) describeInputDefinition(args []ArgDef, options []defOption) {
	show := len(args) > 0
	if show {
		m.write("### Arguments")
		for _, a := range args {
			m.write("\n\n")
			m.describeInputArgument(a)
		}
	}
	if len(options) > 0 {
		if show {
			m.write("\n\n")
		}
		m.write("### Options")
		for _, o := range options {
			m.write("\n\n")
			m.describeInputOption(o)
		}
	}
}

func (m *mdDescriptor) describeCommand(cmd *Command, o descOptions) {
	head := "`" + cmd.Name + "`\n" + strings.Repeat("-", strWidth(cmd.Name)+2) + "\n\n"
	if cmd.Description != "" {
		head += cmd.Description + "\n\n"
	}
	head += "### Usage\n\n"
	usages := cmd.Aliases
	if !o.short {
		usages = append(append([]string{synopsisFor(cmd)}, cmd.Aliases...), cmd.Usages...)
	}
	for _, u := range usages {
		head += "* `" + u + "`\n"
	}
	m.write(head)
	if o.short {
		return
	}
	if help := processedHelp(cmd); help != "" {
		m.write("\n")
		m.write(help)
	}
	options := fullOptions(cmd)
	if len(options) > 0 || len(cmd.Args) > 0 {
		m.write("\n\n")
		m.describeInputDefinition(cmd.Args, options)
	}
}

func (m *mdDescriptor) describeApplication(o descOptions) error {
	d, err := m.app.applicationDescription(o.namespace, false)
	if err != nil {
		return err
	}
	title := appName + " " + appVersion()
	m.write(title + "\n" + strings.Repeat("=", strWidth(title)))
	for _, ns := range d.namespaces {
		if ns.id != globalNamespace {
			m.write("\n\n")
			m.write("**" + ns.id + ":**")
		}
		m.write("\n\n")
		var lines []string
		for _, n := range ns.commands {
			lines = append(lines, fmt.Sprintf("* [`%s`](#%s)", n, strings.ReplaceAll(d.command(n).Name, ":", "")))
		}
		m.write(strings.Join(lines, "\n"))
	}
	for _, n := range d.names {
		m.write("\n\n")
		m.describeCommand(d.commands[n], o)
	}
	return nil
}
