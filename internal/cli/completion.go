package cli

// Port of the shell completion of Symfony Console and Composer:
//   - Symfony\Component\Console\Command\CompleteCommand (the hidden
//     "_complete" command), DumpCompletionCommand ("completion"),
//     Completion\CompletionInput, CompletionSuggestions,
//     Completion\Output\BashCompletionOutput, Resources/completion.bash,
//     the completion parts of Input\ArgvInput (parse, getFirstArgument),
//     Application::find/complete, HelpCommand::complete and
//     ListCommand::complete (symfony/console v5.4.54, MIT,
//     (c) Fabien Potencier), as locked by composer/composer
//     f83588392df43d4fd027d84d26de459946ef4e80;
//   - Composer's completion process handling in Console\Application
//     (plugin and script commands, NullIO Composer instance) and
//     GlobalCommand::complete (composer/composer
//     f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
//     Jordi Boggiano).
// The per-command suggestions (CompletionTrait and friends) are in
// completion_suggest.go. Only bash is supported, like this Symfony version.

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/platform"
)

//go:embed completion.bash
var completionBash string

func init() {
	registerCommand(completionCommand)
	registerCommand(completeCommand)
}

func completionCommand() *Command { return &Command{Name: "completion", Run: runCompletion} }
func completeCommand() *Command   { return &Command{Name: "_complete", Run: runComplete} }

// completionCommandName is basename($_SERVER['argv'][0]).
func completionCommandName() string {
	if len(os.Args) == 0 {
		return ""
	}
	return filepath.Base(os.Args[0])
}

// supportedShells is DumpCompletionCommand::getSupportedShells.
var supportedShells = []string{"bash"}

// runCompletion is DumpCompletionCommand::execute.
func runCompletion(ctx *Context) (int, error) {
	io := ctx.IO
	commandName := completionCommandName()
	if ctx.In.Bool("debug") {
		return 0, tailCompletionLog(ctx, commandName)
	}
	shell := ctx.In.Arg("shell")
	if _, given := ctx.In.args["shell"]; !given {
		// guessShell: basename($_SERVER['SHELL'] ?? '')
		shell = phpBasename(os.Getenv("SHELL"))
	}
	if shell != "bash" {
		if shell != "" {
			io.WriteError(fmt.Sprintf(`<error>Detected shell "%s", which is not supported by Symfony shell completion (supported shells: "%s").</>`, shell, strings.Join(supportedShells, `", "`)))
		} else {
			io.WriteError(fmt.Sprintf(`<error>Shell not detected, Symfony shell completion only supports "%s").</>`, strings.Join(supportedShells, `", "`)))
		}
		return 2, nil
	}
	script := strings.NewReplacer("{{ COMMAND_NAME }}", commandName, "{{ VERSION }}", platform.ComposerVersion).Replace(completionBash)
	io.WriteV(script, false, cio.Normal)
	return 0, nil
}

// phpBasename is PHP's basename() for unix paths.
func phpBasename(p string) string {
	p = strings.TrimRight(p, "/")
	if i := strings.LastIndex(p, "/"); i >= 0 {
		p = p[i+1:]
	}
	return p
}

// completionLogFile is sys_get_temp_dir().'/sf_'.$commandName.'.log'.
func completionLogFile(commandName string) string {
	return filepath.Join(os.TempDir(), "sf_"+commandName+".log")
}

// tailCompletionLog is DumpCompletionCommand::tailDebugLog.
func tailCompletionLog(ctx *Context, commandName string) error {
	file := completionLogFile(commandName)
	if _, err := os.Stat(file); err != nil {
		if f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0o666); err == nil {
			f.Close()
		}
	}
	cmd := exec.Command("tail", "-f", file)
	cmd.Stdout = ctx.IO.Out
	cmd.Stderr = ctx.IO.Err
	return cmd.Run()
}

// completion types (CompletionInput::TYPE_*)
const (
	compArgumentValue = "argument_value"
	compOptionValue   = "option_value"
	compOptionName    = "option_name"
	compNone          = "none"
)

// errCompletion is a failure of the completion process: Composer turns
// PHP warnings into exceptions, which CompleteCommand catches (exit 2).
type errCompletion struct{ msg string }

func (e *errCompletion) Error() string { return e.msg }

// compDef is an InputDefinition (arguments, options with their shortcuts
// and negations).
type compDef struct {
	args []ArgDef
	opts []OptDef
}

func (d *compDef) option(name string) (OptDef, bool) {
	for _, o := range d.opts {
		if o.Name == name {
			return o, true
		}
	}
	return OptDef{}, false
}

func (d *compDef) shortcut(s string) (OptDef, bool) {
	if s == "" {
		return OptDef{}, false
	}
	for _, o := range d.opts {
		if o.Shortcut == "" {
			continue
		}
		for _, sc := range strings.Split(o.Shortcut, "|") {
			if sc == s {
				return o, true
			}
		}
	}
	return OptDef{}, false
}

func (d *compDef) negation(name string) (OptDef, bool) {
	if !strings.HasPrefix(name, "no-") {
		return OptDef{}, false
	}
	o, ok := d.option(name[3:])
	if !ok || !o.Negatable {
		return OptDef{}, false
	}
	return o, true
}

func (d *compDef) argument(name string) (ArgDef, bool) {
	for _, a := range d.args {
		if a.Name == name {
			return a, true
		}
	}
	return ArgDef{}, false
}

func optAcceptsValue(o OptDef) bool { return o.Kind != OptNone }

// appDefinition is Application::getDefinition() (the "command" argument and
// the global options).
func appDefinition() *compDef {
	return &compDef{args: []ArgDef{{Name: "command", Required: true}}, opts: globalOptions}
}

// mergedDefinition is the command definition after
// Command::mergeApplicationDefinition().
func mergedDefinition(cmd *Command) *compDef {
	args := append([]ArgDef{{Name: "command", Required: true}}, cmd.Args...)
	return &compDef{args: args, opts: mergedOptions(cmd)}
}

// mergedOptions are the command options followed by the application
// options it does not define itself (InputDefinition::addOptions keeps an
// existing option in place).
func mergedOptions(cmd *Command) []OptDef {
	opts := append([]OptDef(nil), cmd.Options...)
	for _, g := range globalOptions {
		dup := false
		for _, o := range cmd.Options {
			if o.Name == g.Name {
				dup = true
				break
			}
		}
		if !dup {
			opts = append(opts, g)
		}
	}
	return opts
}

// compArg is an argument value: a string, or a list for array arguments.
type compArg struct {
	s    string
	list []string
	arr  bool
}

// compInput is CompletionInput (an ArgvInput over COMP_WORDS).
type compInput struct {
	// tokens are all input tokens including the program name
	tokens []string
	cur    int
	def    *compDef

	arguments map[string]compArg
	argCount  int
	// options hold string, bool, nil or []any (array options) values
	options map[string]any
	parsed  []string

	ctype  string
	cname  *string
	cvalue string
}

func newCompInput(tokens []string, cur int) *compInput {
	return &compInput{tokens: tokens, cur: cur}
}

// argvTokens are ArgvInput's tokens (without the application name).
func (ci *compInput) argvTokens() []string {
	if len(ci.tokens) == 0 {
		return nil
	}
	return ci.tokens[1:]
}

// parse is ArgvInput::parse with CompletionInput::parseToken suppressing
// the input errors.
func (ci *compInput) parse() {
	parseOptions := true
	ci.parsed = append([]string(nil), ci.argvTokens()...)
	for len(ci.parsed) > 0 {
		tok := ci.parsed[0]
		ci.parsed = ci.parsed[1:]
		next, err := ci.parseToken(tok, parseOptions)
		if err == nil {
			parseOptions = next
		}
	}
}

func (ci *compInput) parseToken(tok string, parseOptions bool) (bool, error) {
	switch {
	case parseOptions && tok == "":
		return parseOptions, ci.parseArgument(tok)
	case parseOptions && tok == "--":
		return false, nil
	case parseOptions && strings.HasPrefix(tok, "--"):
		return parseOptions, ci.parseLongOption(tok)
	case parseOptions && tok[0] == '-' && tok != "-":
		return parseOptions, ci.parseShortOption(tok)
	}
	return parseOptions, ci.parseArgument(tok)
}

func (ci *compInput) parseShortOption(tok string) error {
	name := tok[1:]
	if len(name) > 1 {
		if o, ok := ci.def.shortcut(name[:1]); ok && optAcceptsValue(o) {
			// an option with a value (with no space)
			return ci.addShortOption(name[:1], strPtr(name[1:]))
		}
		return ci.parseShortOptionSet(name)
	}
	return ci.addShortOption(name, nil)
}

func (ci *compInput) parseShortOptionSet(name string) error {
	for i := 0; i < len(name); i++ {
		o, ok := ci.def.shortcut(name[i : i+1])
		if !ok {
			return fmt.Errorf(`The "-%s" option does not exist.`, name[i:i+1])
		}
		if optAcceptsValue(o) {
			if i == len(name)-1 {
				return ci.addLongOption(o.Name, nil)
			}
			return ci.addLongOption(o.Name, strPtr(name[i+1:]))
		}
		if err := ci.addLongOption(o.Name, nil); err != nil {
			return err
		}
	}
	return nil
}

func (ci *compInput) parseLongOption(tok string) error {
	name := tok[2:]
	if pos := strings.Index(name, "="); pos >= 0 {
		value := name[pos+1:]
		if value == "" {
			ci.parsed = append([]string{value}, ci.parsed...)
		}
		return ci.addLongOption(name[:pos], &value)
	}
	return ci.addLongOption(name, nil)
}

func (ci *compInput) parseArgument(tok string) error {
	c := ci.argCount
	if c < len(ci.def.args) {
		a := ci.def.args[c]
		if a.Array {
			ci.arguments[a.Name] = compArg{list: []string{tok}, arr: true}
		} else {
			ci.arguments[a.Name] = compArg{s: tok}
		}
		ci.argCount++
		return nil
	}
	if c-1 >= 0 && c-1 < len(ci.def.args) && ci.def.args[c-1].Array {
		a := ci.def.args[c-1]
		v := ci.arguments[a.Name]
		v.list = append(v.list, tok)
		ci.arguments[a.Name] = v
		return nil
	}
	return errors.New("too many arguments")
}

func (ci *compInput) addShortOption(shortcut string, value *string) error {
	o, ok := ci.def.shortcut(shortcut)
	if !ok {
		return fmt.Errorf(`The "-%s" option does not exist.`, shortcut)
	}
	return ci.addLongOption(o.Name, value)
}

func (ci *compInput) addLongOption(name string, value *string) error {
	o, ok := ci.def.option(name)
	if !ok {
		neg, nok := ci.def.negation(name)
		if !nok {
			return fmt.Errorf(`The "--%s" option does not exist.`, name)
		}
		if value != nil {
			return fmt.Errorf(`The "--%s" option does not accept a value.`, name)
		}
		ci.options[neg.Name] = false
		return nil
	}
	if value != nil && !optAcceptsValue(o) {
		return fmt.Errorf(`The "--%s" option does not accept a value.`, name)
	}
	if (value == nil || *value == "") && optAcceptsValue(o) && len(ci.parsed) > 0 {
		// if option accepts an optional or mandatory argument
		// let's see if there is one provided
		next := ci.parsed[0]
		ci.parsed = ci.parsed[1:]
		if next != "" && next[0] != '-' || next == "" {
			value = &next
		} else {
			ci.parsed = append([]string{next}, ci.parsed...)
		}
	}
	var v any
	if value == nil {
		if o.Kind == OptRequired || o.Kind == OptArray {
			return fmt.Errorf(`The "--%s" option requires a value.`, name)
		}
		if o.Kind == OptNone {
			v = true
		}
	} else {
		v = *value
	}
	if o.Kind == OptArray {
		list, _ := ci.options[name].([]any)
		ci.options[name] = append(list, v)
	} else {
		ci.options[name] = v
	}
	return nil
}

// isCursorFree is CompletionInput::isCursorFree.
func (ci *compInput) isCursorFree() (bool, error) {
	n := len(ci.tokens)
	if ci.cur > n {
		return false, &errCompletion{"Current index is invalid, it must be the number of input tokens or one more."}
	}
	return ci.cur >= n, nil
}

// token returns $this->tokens[$i] (an undefined key is a PHP warning).
func (ci *compInput) token(i int) (string, error) {
	if i < 0 || i >= len(ci.tokens) {
		return "", &errCompletion{fmt.Sprintf("Undefined array key %d", i)}
	}
	return ci.tokens[i], nil
}

// firstByte is $token[0] (an empty string is a PHP warning).
func firstByte(s string) (byte, error) {
	if s == "" {
		return 0, &errCompletion{"Uninitialized string offset 0"}
	}
	return s[0], nil
}

// bind is CompletionInput::bind.
func (ci *compInput) bind(def *compDef) error {
	ci.def = def
	ci.arguments = map[string]compArg{}
	ci.argCount = 0
	ci.options = map[string]any{}
	ci.parse()
	ci.ctype, ci.cname, ci.cvalue = "", nil, ""

	free, err := ci.isCursorFree()
	if err != nil {
		return err
	}
	idx := ci.cur
	if free {
		idx = ci.cur - 1
	}
	relevant, err := ci.token(idx)
	if err != nil {
		return err
	}
	b, err := firstByte(relevant)
	if err != nil {
		return err
	}
	if b == '-' {
		// the current token is an input option: complete either option name or option value
		optionToken, optionValue, _ := strings.Cut(relevant, "=")
		o, ok := ci.optionFromToken(optionToken)
		if !ok && !free {
			ci.ctype = compOptionName
			ci.cvalue = relevant
			return nil
		}
		if ok && optAcceptsValue(o) {
			ci.ctype = compOptionValue
			name := o.Name
			ci.cname = &name
			switch {
			case optionValue != "" && optionValue != "0":
				ci.cvalue = optionValue
			case !strings.HasPrefix(optionToken, "--"):
				if len(optionToken) > 2 {
					ci.cvalue = optionToken[2:]
				}
			}
			return nil
		}
	}

	prev, err := ci.token(ci.cur - 1)
	if err != nil {
		return err
	}
	pb, err := firstByte(prev)
	if err != nil {
		return err
	}
	if pb == '-' && strings.Trim(prev, "-") != "" {
		// check if previous option accepted a value
		if o, ok := ci.optionFromToken(prev); ok && optAcceptsValue(o) {
			ci.ctype = compOptionValue
			name := o.Name
			ci.cname = &name
			ci.cvalue = relevant
			return nil
		}
	}

	// complete argument value
	ci.ctype = compArgumentValue
	argumentName := ""
	for _, a := range def.args {
		argumentName = a.Name
		v, set := ci.arguments[a.Name]
		if !set {
			break
		}
		name := a.Name
		ci.cname = &name
		if v.arr {
			if len(v.list) == 0 {
				return &errCompletion{"getCompletionValue(): Return value must be of type string, null returned"}
			}
			ci.cvalue = v.list[len(v.list)-1]
		} else {
			ci.cvalue = v.s
		}
	}
	if len(def.args) == 0 {
		return &errCompletion{"Undefined variable $argumentName"}
	}
	if ci.cur >= len(ci.tokens) {
		a, _ := def.argument(argumentName)
		if _, set := ci.arguments[argumentName]; !set || a.Array {
			name := argumentName
			ci.cname = &name
			ci.cvalue = ""
		} else {
			// we've reached the end
			ci.ctype = compNone
			ci.cname = nil
			ci.cvalue = ""
		}
	}
	return nil
}

// optionFromToken is CompletionInput::getOptionFromToken.
func (ci *compInput) optionFromToken(optionToken string) (OptDef, bool) {
	optionName := strings.TrimLeft(optionToken, "-")
	if optionName == "" || optionName == "0" {
		return OptDef{}, false
	}
	if len(optionToken) > 1 && optionToken[1] == '-' {
		// long option name
		return ci.def.option(optionName)
	}
	// short option name
	return ci.def.shortcut(optionName[:1])
}

func (ci *compInput) completionName() string {
	if ci.cname == nil {
		return ""
	}
	return *ci.cname
}

func (ci *compInput) mustSuggestOptionValuesFor(name string) bool {
	return ci.ctype == compOptionValue && ci.cname != nil && *ci.cname == name
}

func (ci *compInput) mustSuggestArgumentValuesFor(name string) bool {
	return ci.ctype == compArgumentValue && ci.cname != nil && *ci.cname == name
}

// getArgument is Input::getArgument (nil: the null default).
func (ci *compInput) getArgument(name string) *string {
	if v, ok := ci.arguments[name]; ok {
		if v.arr {
			if len(v.list) == 0 {
				return nil
			}
			s := v.list[0]
			return &s
		}
		s := v.s
		return &s
	}
	if a, ok := ci.def.argument(name); ok {
		if s, ok := a.RawDefault.(string); ok {
			return &s
		}
	}
	return nil
}

// optionTruthy is (bool) $input->getOption($name).
func (ci *compInput) optionTruthy(name string) bool {
	if v, ok := ci.options[name]; ok {
		switch t := v.(type) {
		case []any:
			return len(t) > 0
		case string:
			return t != "" && t != "0"
		case bool:
			return t
		}
		return false
	}
	if o, ok := ci.def.option(name); ok {
		return phpjson.Truthy(o.RawDefault)
	}
	return false
}

// optionString is $input->getOption($name) for an option taking a value
// (nil when unset and without default).
func (ci *compInput) optionString(name string) *string {
	if v, ok := ci.options[name]; ok {
		if s, ok := v.(string); ok {
			return &s
		}
		return nil
	}
	if o, ok := ci.def.option(name); ok {
		if s, ok := o.RawDefault.(string); ok {
			return &s
		}
	}
	return nil
}

// firstArgument is ArgvInput::getFirstArgument.
func (ci *compInput) firstArgument() (string, bool, error) {
	toks := ci.argvTokens()
	isOption := false
	for i, tok := range toks {
		if tok != "" && tok != "0" && tok[0] == '-' {
			if strings.Contains(tok, "=") || i+1 >= len(toks) {
				continue
			}
			if len(tok) < 2 {
				return "", false, &errCompletion{"Uninitialized string offset 1"}
			}
			var name string
			if tok[1] == '-' {
				name = tok[2:]
			} else {
				name = tok[len(tok)-1:]
			}
			v, set := ci.options[name]
			set = set && v != nil
			_, hasShort := ci.def.shortcut(name)
			if !set && !hasShort {
				// noop
			} else {
				if !set {
					o, _ := ci.def.shortcut(name)
					name = o.Name
					v, set = ci.options[name]
					set = set && v != nil
				}
				if set {
					if s, ok := v.(string); ok && s == toks[i+1] {
						isOption = true
					}
				}
			}
			continue
		}
		if isOption {
			isOption = false
			continue
		}
		return tok, true, nil
	}
	return "", false, nil
}

// String is CompletionInput::__toString ("|" marks the cursor).
func (ci *compInput) String() string {
	var sb strings.Builder
	i := -1
	for j, t := range ci.tokens {
		i = j
		sb.WriteString(t)
		if ci.cur == j {
			sb.WriteString("|")
		}
		sb.WriteString(" ")
	}
	if ci.cur > i {
		sb.WriteString("|")
	}
	return strings.TrimRight(sb.String(), " \t\n\r\x00\x0B")
}

// compInputFromString is CompletionInput::fromString: the tokens matched by
// /(?<=^|\s)(['"]?)(.+?)(?<!\\)\1(?=$|\s)/.
func compInputFromString(s string, cur int) *compInput {
	isSpace := func(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f' }
	matchAt := func(pos int) (int, bool) {
		if pos > 0 && !isSpace(s[pos-1]) {
			return 0, false
		}
		quotes := []string{""}
		if pos < len(s) && (s[pos] == '\'' || s[pos] == '"') {
			quotes = []string{s[pos : pos+1], ""}
		}
		for _, q := range quotes {
			start := pos + len(q)
			for e := start + 1; e <= len(s); e++ {
				if s[e-1] == '\n' {
					break
				}
				if s[e-1] == '\\' || !strings.HasPrefix(s[e:], q) {
					continue
				}
				after := e + len(q)
				if after == len(s) || isSpace(s[after]) || (after == len(s)-1 && s[after] == '\n') {
					return after, true
				}
			}
		}
		return 0, false
	}
	var tokens []string
	for pos := 0; pos < len(s); {
		if end, ok := matchAt(pos); ok {
			tokens = append(tokens, s[pos:end])
			pos = end
			continue
		}
		pos++
	}
	return newCompInput(tokens, cur)
}

// compSuggestions is CompletionSuggestions.
type compSuggestions struct {
	values  []string
	options []OptDef
}

func (s *compSuggestions) suggestValues(vs []string) { s.values = append(s.values, vs...) }

// suggestNonEmpty is suggestValues(array_filter($values)).
func (s *compSuggestions) suggestNonEmpty(vs []string) {
	for _, v := range vs {
		if v != "" && v != "0" {
			s.values = append(s.values, v)
		}
	}
}

// compSession is the state of one _complete run (the Application and the
// Composer instance created for it with a NullIO).
type compSession struct {
	ctx      *Context
	app      *App
	composer *composer.Composer
}

func nullIO() *cio.IO { return &cio.IO{Out: io.Discard, Err: io.Discard, Verbosity: cio.Normal} }

// pluginLoadIO is the IO plugin commands are loaded with in the completion
// process: the Composer instance gets a NullIO, while Application::doRun
// writes the skipped plugin command warnings to the real output.
func pluginLoadIO(real *cio.IO) *cio.IO {
	o := nullIO()
	o.Decorated = real.Decorated
	o.Err = &lineFilter{dst: real.Err, keep: func(l string) bool {
		return strings.Contains(l, "would override a Composer command and has been skipped")
	}}
	return o
}

// lineFilter forwards the complete lines accepted by keep.
type lineFilter struct {
	dst  io.Writer
	keep func(string) bool
	buf  []byte
}

func (f *lineFilter) Write(p []byte) (int, error) {
	f.buf = append(f.buf, p...)
	for {
		i := strings.IndexByte(string(f.buf), '\n')
		if i < 0 {
			break
		}
		line := string(f.buf[:i+1])
		f.buf = f.buf[i+1:]
		if f.keep(line) {
			if _, err := io.WriteString(f.dst, line); err != nil {
				return len(p), err
			}
		}
	}
	return len(p), nil
}

// requireComposer is BaseCommand::requireComposer in the completion
// process: Application::getComposer(true) with a NullIO; a missing
// composer.json is written to the real output and exits 1.
func (s *compSession) requireComposer() (*composer.Composer, error) {
	if s.composer != nil {
		return s.composer, nil
	}
	if c := takePreloadedComposer(); c != nil {
		s.composer = c
		return c, nil
	}
	in := s.ctx.In
	c, err := composer.Create(nullIO(), composer.Options{DisablePlugins: in.Bool("no-plugins") || pluginsDisabledByRoot, DisableScripts: in.Bool("no-scripts"), FullLoad: true})
	if err != nil {
		var ue *composer.UserError
		if errors.As(err, &ue) && strings.HasPrefix(ue.Msg, "Composer could not find ") {
			return nil, &composer.FactoryInvalidArgumentError{Msg: ue.Msg}
		}
		return nil, err
	}
	s.composer = c
	return c, nil
}

// compEntry is one entry of Application::all(): a command under its name or
// one of its aliases.
type compEntry struct {
	key string
	cmd *Command
}

// composerCommandOrder is the order Composer's Application registers its
// commands (Symfony's default commands, then getDefaultCommands()).
var composerCommandOrder = []string{"help", "list", "_complete", "completion", "about", "config", "depends", "prohibits", "init", "install", "create-project", "update", "search", "validate", "audit", "show", "suggests", "require", "dump-autoload", "status", "archive", "diagnose", "run-script", "licenses", "global", "clear-cache", "remove", "browse", "exec", "outdated", "check-platform-reqs", "fund", "reinstall", "bump", "repository", "policy", "self-update"}

// allEntries is Application::all(): the commands in registration order, each
// followed by its aliases (plugin and script commands come last).
func (a *App) allEntries() []compEntry {
	var ordered []*Command
	seen := map[*Command]bool{}
	for _, n := range composerCommandOrder {
		for _, c := range a.commands {
			if c.Name == n && !seen[c] {
				ordered = append(ordered, c)
				seen[c] = true
			}
		}
	}
	for _, c := range a.commands {
		if !seen[c] {
			ordered = append(ordered, c)
			seen[c] = true
		}
	}
	var out []compEntry
	keys := map[string]int{}
	add := func(k string, c *Command) {
		if i, ok := keys[k]; ok {
			out[i].cmd = c
			return
		}
		keys[k] = len(out)
		out = append(out, compEntry{k, c})
	}
	for _, c := range ordered {
		add(c.Name, c)
		for _, al := range c.Aliases {
			add(al, c)
		}
	}
	return out
}

// symfonyFind is Symfony's Application::find (nil when no single command
// matches).
func (a *App) symfonyFind(name string) *Command {
	all := a.allEntries()
	byKey := map[string]*Command{}
	for _, e := range all {
		byKey[e.key] = e.cmd
	}
	if c, ok := byKey[name]; ok {
		return c
	}
	parts := strings.Split(name, ":")
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	expr := strings.Join(parts, "[^:]*:") + "[^:]*"
	grep := func(pattern string) []string {
		re := regexp.MustCompile(pattern)
		var out []string
		for _, e := range all {
			if re.MatchString(e.key) {
				out = append(out, e.key)
			}
		}
		return out
	}
	commands := grep("^" + expr)
	if len(commands) == 0 {
		commands = grep("(?i)^" + expr)
	}
	exact := regexp.MustCompile("(?i)^" + expr + "$")
	nExact := 0
	for _, c := range commands {
		if exact.MatchString(c) {
			nExact++
		}
	}
	if len(commands) == 0 || nExact < 1 {
		return nil
	}
	if len(commands) > 1 {
		// filter out aliases for commands which are already on the list
		in := map[string]bool{}
		for _, c := range commands {
			in[c] = true
		}
		var filtered []string
		for _, c := range commands {
			cn := byKey[c].Name
			if cn == c || !in[cn] {
				filtered = append(filtered, c)
			}
		}
		commands = filtered
	}
	if len(commands) > 1 {
		var visible []string
		for _, c := range commands {
			if !byKey[c].Hidden {
				visible = append(visible, c)
			}
		}
		commands = visible
		if len(commands) > 1 || len(commands) == 0 {
			return nil
		}
	}
	c := byKey[commands[0]]
	if c.Hidden {
		return nil
	}
	return c
}

// appComplete is Application::complete.
func (a *App) appComplete(ci *compInput, s *compSuggestions) {
	if ci.ctype == compArgumentValue && ci.completionName() == "command" {
		var names []string
		for _, e := range a.allEntries() {
			// skip hidden commands and aliased commands as they already get added below
			if e.cmd.Hidden || e.cmd.Name != e.key {
				continue
			}
			names = append(names, e.cmd.Name)
			names = append(names, e.cmd.Aliases...)
		}
		s.suggestNonEmpty(names)
		return
	}
	if ci.ctype == compOptionName {
		s.options = append(s.options, globalOptions...)
	}
}

// runComplete is CompleteCommand::execute.
func runComplete(ctx *Context) (int, error) {
	debugLog := envFilterBool(os.Getenv("SYMFONY_COMPLETION_DEBUG"))
	logf := func(lines ...string) {
		if !debugLog {
			return
		}
		f, err := os.OpenFile(completionLogFile(completionCommandName()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			return
		}
		defer f.Close()
		for _, l := range lines {
			fmt.Fprintln(f, l)
		}
	}
	code, err := completeExecute(ctx, logf)
	if err == nil {
		return code, nil
	}
	var fia *composer.FactoryInvalidArgumentError
	if errors.As(err, &fia) {
		return 1, err
	}
	logf("<error>Error!</error>", err.Error())
	if ctx.IO.IsDebug() {
		return 1, err
	}
	return 2, nil
}

// envFilterBool is filter_var($v, FILTER_VALIDATE_BOOLEAN).
func envFilterBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

func completeExecute(ctx *Context, logf func(...string)) (int, error) {
	in := ctx.In
	shell, _ := in.String("shell")
	if shell == "" || shell == "0" {
		return 2, &errCompletion{`The "--shell" option must be set.`}
	}
	if shell != "bash" {
		return 2, &errCompletion{fmt.Sprintf(`Shell completion is not supported for your shell: "%s" (supported: "bash").`, shell)}
	}
	current, _ := in.String("current")
	if current == "" || current == "0" || strings.Trim(current, "0123456789") != "" {
		return 2, &errCompletion{`The "--current" option must be set and it must be an integer.`}
	}
	cur, err := strconv.Atoi(current)
	if err != nil {
		return 2, &errCompletion{`The "--current" option must be set and it must be an integer.`}
	}
	app := ctx.App
	// Application::doRun loads plugin and script commands for the
	// completion process (plugins with a NullIO Composer instance)
	if !isNonAllowedRoot() {
		if err := app.addPluginCommands(pluginLoadIO(ctx.IO), in.Raw, "list"); err != nil {
			return 1, err
		}
	}
	if err := app.addScriptCommands(ctx.IO, in.Raw, "list"); err != nil {
		return 1, err
	}
	s := &compSession{ctx: ctx, app: app}

	ci := newCompInput(in.Strings("input"), cur)
	if err := ci.bind(appDefinition()); err != nil {
		return 2, err
	}
	sugg := &compSuggestions{}
	logf("", "<comment>"+time.Now().Format("2006-01-02 15:04:05")+"</>", `<info>Input:</> <comment>("|" indicates the cursor position)</>`, "  "+ci.String(), "<info>Command:</>", "  "+strings.Join(os.Args, " "), "<info>Messages:</>")

	var cmd *Command
	if name, ok, err := ci.firstArgument(); err != nil {
		return 2, err
	} else if ok {
		cmd = app.symfonyFind(name)
	}
	switch {
	case cmd == nil:
		logf("  No command found, completing using the Application class.")
		app.appComplete(ci, sugg)
	case ci.mustSuggestArgumentValuesFor("command") && cmd.Name != ci.cvalue && !containsStr(cmd.Aliases, ci.cvalue):
		logf("  No command found, completing using the Application class.")
		// expand shortcut names ("cache:cl<TAB>") into their full name ("cache:clear")
		sugg.suggestNonEmpty(append([]string{cmd.Name}, cmd.Aliases...))
	default:
		if err := ci.bind(mergedDefinition(cmd)); err != nil {
			return 2, err
		}
		if ci.ctype == compOptionName {
			logf("  Completing option names for the <comment>" + cmd.Name + "</> command.")
			sugg.options = append(sugg.options, mergedOptions(cmd)...)
		} else {
			logf("  Completing using the <comment>"+cmd.Name+"</> class.", "  Completing <comment>"+ci.ctype+"</> for <comment>"+ci.completionName()+"</>")
			if ci.cname != nil || ci.ctype != compNone {
				logf("  Current value: <comment>" + ci.cvalue + "</>")
			}
			if err := s.commandComplete(cmd, ci, sugg); err != nil {
				return 2, err
			}
		}
	}

	logf("<info>Suggestions:</>")
	var lines []string
	lines = append(lines, sugg.values...)
	for _, o := range sugg.options {
		lines = append(lines, "--"+o.Name)
		if o.Negatable {
			lines = append(lines, "--no-"+o.Name)
		}
	}
	if len(sugg.options) > 0 {
		var names []string
		for _, o := range sugg.options {
			names = append(names, o.Name)
		}
		logf("  --" + strings.Join(names, " --"))
	} else if len(sugg.values) > 0 {
		logf("  " + strings.Join(sugg.values, " "))
	} else {
		logf("  <comment>No suggestions were provided</>")
	}
	ctx.IO.Write(strings.Join(lines, "\n"))
	return 0, nil
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// commandComplete is $command->complete() for the found command.
func (s *compSession) commandComplete(cmd *Command, ci *compInput, sugg *compSuggestions) error {
	switch cmd.Name {
	case "help":
		// HelpCommand::complete
		if ci.mustSuggestArgumentValuesFor("command_name") {
			sugg.suggestValues(s.app.describedCommandNames())
			return nil
		}
		if ci.mustSuggestOptionValuesFor("format") {
			sugg.suggestValues([]string{"txt", "xml", "json", "md"})
		}
		return nil
	case "list":
		// ListCommand::complete
		if ci.mustSuggestArgumentValuesFor("namespace") {
			sugg.suggestValues(s.app.describedNamespaces())
			return nil
		}
		if ci.mustSuggestOptionValuesFor("format") {
			sugg.suggestValues([]string{"txt", "xml", "json", "md"})
		}
		return nil
	case "completion":
		// DumpCompletionCommand::complete
		if ci.mustSuggestArgumentValuesFor("shell") {
			sugg.suggestValues(supportedShells)
		}
		return nil
	case "_complete":
		return nil
	case "global":
		return s.globalComplete(ci, sugg)
	}
	if pluginCommands[cmd] {
		if ci.ctype == compOptionValue || ci.ctype == compArgumentValue {
			// the plugin's (PHP) complete() implementation is not run
			return &errCompletion{"gocomposer does not run the completion of plugin command " + cmd.Name}
		}
		return nil
	}
	// BaseCommand::complete: the suggested values of the option/argument
	values, err := s.definitionSuggestions(cmd, ci)
	if err != nil {
		return err
	}
	sugg.suggestValues(values)
	return nil
}

// describedCommandNames are the keys of ApplicationDescription::getCommands():
// the visible command names, sorted by namespace and name.
func (a *App) describedCommandNames() []string {
	var out []string
	for _, ns := range a.describedNamespaceGroups() {
		for _, e := range ns.entries {
			if e.cmd.Hidden || e.cmd.Name != e.key {
				continue
			}
			out = append(out, e.key)
		}
	}
	return out
}

// describedNamespaces are the keys of ApplicationDescription::getNamespaces().
func (a *App) describedNamespaces() []string {
	var out []string
	for _, ns := range a.describedNamespaceGroups() {
		out = append(out, ns.name)
	}
	return out
}

type describedNamespace struct {
	name    string
	entries []compEntry
}

// describedNamespaceGroups is ApplicationDescription::sortCommands.
func (a *App) describedNamespaceGroups() []describedNamespace {
	groups := map[string][]compEntry{}
	for _, e := range a.allEntries() {
		key := "_global"
		if parts := strings.Split(e.key, ":"); len(parts) > 1 && parts[0] != "" {
			key = parts[0]
		}
		groups[key] = append(groups[key], e)
	}
	var names []string
	for k := range groups {
		names = append(names, k)
	}
	sort.SliceStable(names, func(i, j int) bool { return phpjson.PHPCompareKeys(names[i], names[j]) < 0 })
	var out []describedNamespace
	for _, n := range names {
		entries := groups[n]
		sort.SliceStable(entries, func(i, j int) bool { return phpjson.PHPCompareKeys(entries[i].key, entries[j].key) < 0 })
		out = append(out, describedNamespace{n, entries})
	}
	return out
}

// globalComplete is GlobalCommand::complete.
func (s *compSession) globalComplete(ci *compInput, sugg *compSuggestions) error {
	app := s.app
	if ci.mustSuggestArgumentValuesFor("command-name") {
		var names []string
		for _, e := range app.allEntries() {
			if !e.cmd.Hidden {
				names = append(names, e.cmd.Name)
			}
		}
		sugg.suggestValues(names)
		return nil
	}
	commandName := ci.getArgument("command-name")
	if commandName == nil {
		return &errCompletion{"Composer\\Console\\Application::has(): Argument #1 ($name) must be of type string, null given"}
	}
	has := false
	for _, e := range app.allEntries() {
		if e.key == *commandName {
			has = true
		}
	}
	if !has {
		return nil
	}
	// prepareSubcommandInput($input, true)
	if os.Getenv("COMPOSER") != "" {
		os.Unsetenv("COMPOSER")
	}
	cwd, _ := os.Getwd()
	cfg, err := createConfigNoIO(cwd)
	if err != nil {
		return err
	}
	home := cfg.GetString("home")
	if st, err := os.Stat(home); err != nil || !st.IsDir() {
		_ = os.MkdirAll(home, 0o777)
		if st, err := os.Stat(home); err != nil || !st.IsDir() {
			return &errCompletion{"Could not create home directory"}
		}
	}
	if err := os.Chdir(home); err != nil {
		return &errCompletion{`Could not switch to home directory "` + home + `"`}
	}
	str := ci.String()
	if loc := reGlobalWord.FindStringIndex(str); loc != nil {
		str = str[:loc[0]] + str[loc[1]:]
	}
	toks, err := stringInputTokenize(str)
	if err != nil {
		return err
	}
	// Application::resetComposer
	s.composer = nil
	preloadedComposer = nil
	sub := compInputFromString(argvInputString(toks), 2)
	cmd := app.symfonyFind(*commandName)
	if cmd == nil {
		return &errCompletion{`Command "` + *commandName + `" is not defined.`}
	}
	if err := sub.bind(mergedDefinition(cmd)); err != nil {
		return err
	}
	return s.commandComplete(cmd, sub, sugg)
}
