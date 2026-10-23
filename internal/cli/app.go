// Package cli implements the gocomposer command line (Composer's Application
// and commands), with Symfony Console style option parsing and errors.
package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/platform"
	"gocomposer/internal/solver"
)

// isTerminal is stream_isatty() (Platform::isTty).
func isTerminal(f *os.File) bool { return cio.IsTTY(f) }

// Version of gocomposer itself.
const Version = "0.1.0-dev"

// OptKind is how an option takes values.
type OptKind int

const (
	OptNone OptKind = iota
	OptRequired
	OptOptional
	OptArray
)

// OptDef defines an option.
type OptDef struct {
	Name       string
	Shortcut   string
	Kind       OptKind
	Default    string
	Desc       string
	Negatable  bool
	RawDefault interface{}
}

// ArgDef defines an argument.
type ArgDef struct {
	Name       string
	Required   bool
	Array      bool
	Desc       string
	RawDefault interface{}
}

// Command is one CLI command.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Options     []OptDef
	Args        []ArgDef
	Run         func(ctx *Context) (int, error)
	Hidden      bool
	Help        string
	Usages      []string
	Synopsis    string
	// IgnoreValidationErrors is Command::ignoreValidationErrors(): input
	// errors stop parsing but do not abort (ScriptAliasCommand).
	IgnoreValidationErrors bool
}

// Input is the parsed command line.
type Input struct {
	opts    map[string][]string
	present map[string]bool
	args    map[string][]string
	Raw     []string
	// missing is Input::validate()'s error, raised after the command was
	// initialized like Symfony's Command::run.
	missing error
}

func (in *Input) Bool(name string) bool { return in.present[name] }
func (in *Input) String(name string) (string, bool) {
	v, ok := in.opts[name]
	if !ok || len(v) == 0 {
		return "", in.present[name]
	}
	return v[len(v)-1], true
}
func (in *Input) Strings(name string) []string { return in.opts[name] }
func (in *Input) Arg(name string) string {
	if v := in.args[name]; len(v) > 0 {
		return v[0]
	}
	return ""
}
func (in *Input) Args(name string) []string { return in.args[name] }

// Context is passed to commands.
type Context struct {
	IO    *cio.IO
	In    *Input
	App   *App
	Cmd   *Command
	cache *composer.Composer
}

// HasOption reports whether the command (or the application) defines name.
func (c *Command) HasOption(name string) bool {
	for _, o := range allOptions(c) {
		if o.Name == name {
			return true
		}
	}
	return false
}

func (in *Input) setTrue(name string) {
	in.present[name] = true
}

// envTruthy is PHP's (bool) cast of Platform::getEnv().
func envTruthy(name string) bool {
	v, ok := os.LookupEnv(name)
	return ok && v != "" && v != "0"
}

// applyEnvOptions ports the env-to-option mapping of BaseCommand::initialize.
func applyEnvOptions(io *cio.IO, cmd *Command, in *Input) {
	if in.present["no-ansi"] && cmd.HasOption("no-progress") {
		in.setTrue("no-progress")
	}
	envOptions := []struct {
		env  string
		opts []string
	}{
		{"COMPOSER_NO_AUDIT", []string{"no-audit"}},
		{"COMPOSER_NO_DEV", []string{"no-dev", "update-no-dev"}},
		{"COMPOSER_PREFER_STABLE", []string{"prefer-stable"}},
		{"COMPOSER_PREFER_LOWEST", []string{"prefer-lowest"}},
		{"COMPOSER_MINIMAL_CHANGES", []string{"minimal-changes"}},
		{"COMPOSER_WITH_DEPENDENCIES", []string{"with-dependencies"}},
		{"COMPOSER_WITH_ALL_DEPENDENCIES", []string{"with-all-dependencies"}},
		{"COMPOSER_NO_SECURITY_BLOCKING", []string{"no-security-blocking"}},
		{"COMPOSER_NO_BLOCKING", []string{"no-blocking"}},
	}
	for _, e := range envOptions {
		for _, o := range e.opts {
			if cmd.HasOption(o) && !in.Bool(o) && envTruthy(e.env) {
				in.setTrue(o)
			}
		}
	}
	if cmd.HasOption("ignore-platform-reqs") && !in.Bool("ignore-platform-reqs") && envTruthy("COMPOSER_IGNORE_PLATFORM_REQS") {
		in.setTrue("ignore-platform-reqs")
		io.WriteError("<warning>COMPOSER_IGNORE_PLATFORM_REQS is set. You may experience unexpected errors.</warning>")
	}
	if cmd.HasOption("ignore-platform-req") && (!cmd.HasOption("ignore-platform-reqs") || !in.Bool("ignore-platform-reqs")) {
		if v := os.Getenv("COMPOSER_IGNORE_PLATFORM_REQ"); v != "" && len(in.Strings("ignore-platform-req")) == 0 {
			in.opts["ignore-platform-req"] = strings.Split(v, ",")
			in.present["ignore-platform-req"] = true
			io.WriteError("<warning>COMPOSER_IGNORE_PLATFORM_REQ is set to ignore " + v + ". You may experience unexpected errors.</warning>")
		}
	}
}

// App holds registered commands.
type App struct {
	commands []*Command
	byName   map[string]*Command
	// running is Symfony's runningCommand: set once the command was found,
	// so errors raised from then on print the command synopsis.
	running *Command
}

var globalOptions = defOptions(commandDefs.GlobalOptions)

// registeredCommands holds constructors added from init() in per-command
// files, so commands can be added without editing NewApp.
var registeredCommands []func() *Command

func registerCommand(f func() *Command) { registeredCommands = append(registeredCommands, f) }

// NewApp registers all commands.
func NewApp() *App {
	a := &App{byName: map[string]*Command{}}
	impl := map[string]*Command{}
	for _, c := range []*Command{installCommand(), updateCommand(), dumpAutoloadCommand(), validateCommand(), showCommand(), checkPlatformCommand(),
		requireCommand(), removeCommand(), configCommand(), aboutCommand(), listCommand(), helpCommand(), diagnoseCommand(), licensesCommand(), suggestsCommand(), dependsCommand(false), dependsCommand(true)} {
		impl[c.Name] = c
	}
	for _, f := range registeredCommands {
		c := f()
		impl[c.Name] = c
	}
	for _, def := range commandDefs.Commands {
		c := impl[def.Name]
		if c == nil {
			c = &Command{Run: notImplemented(def.Name)}
		}
		applyDefs(c, def)
		a.commands = append(a.commands, c)
		a.byName[c.Name] = c
		for _, al := range c.Aliases {
			a.byName[al] = c
		}
	}
	return a
}

// CLIError mirrors Symfony console exceptions (rendered without file/line).
type CLIError struct{ Msg string }

func (e *CLIError) Error() string { return e.Msg }

// Run executes the application and returns the exit code.
func (a *App) Run(argv []string) int {
	io := cio.New()
	// ConsoleOutput: decorated when both streams support colors
	io.Decorated = cio.HasColorSupport(os.Stdout) && cio.HasColorSupport(os.Stderr)
	io.Interactive = isTerminal(os.Stdin)
	// the PHP processes running Composer plugins end with the command
	defer composer.ShutdownPlugins()
	// Application::doRun reports the program status to the terminal
	io.InitProgramStatus()
	reportStatus := programStatusReported(io, argv)
	if reportStatus {
		io.StartProgramStatus()
		orig, had := os.LookupEnv("COMPOSER_TERMINAL_STATUS")
		// nested Composer processes inherit disabled reporting
		os.Setenv("COMPOSER_TERMINAL_STATUS", "off")
		defer func() {
			if had {
				os.Setenv("COMPOSER_TERMINAL_STATUS", orig)
			} else {
				os.Unsetenv("COMPOSER_TERMINAL_STATUS")
			}
		}()
	}
	code, err := a.run(io, argv)
	var se *composer.ScriptError
	if errors.As(err, &se) {
		rootScriptHint(io, argv)
		if reportStatus {
			io.WriteProgramStatus(cio.StatusFromExitCode(se.Code))
		}
		// Application::doRun returns the code of a ScriptExecutionException
		// without rendering it
		if se.Code > 255 {
			return 255
		}
		return se.Code
	}
	var fia *composer.FactoryInvalidArgumentError
	if errors.As(err, &fia) {
		// Application::getComposer(true) writes the message and exits
		io.WriteError(fia.Msg)
		if reportStatus {
			io.WriteProgramStatus(cio.StatusError())
		}
		return 1
	}
	if err != nil {
		// Application::hintCommonErrors
		commandNotFoundHint(io, argv, err)
		var hinted interface{ ExceptionHints() []string }
		if errors.As(err, &hinted) {
			for _, h := range hinted.ExceptionHints() {
				io.WriteErrorV(h, true, cio.Quiet)
			}
		}
		// HttpDownloader::getExceptionHints and the offline-mode hint
		var hte *httpdl.TransportError
		if errors.As(err, &hte) {
			for _, h := range exceptionHints(err) {
				io.WriteErrorV(h, true, cio.Quiet)
			}
			if (a.running == nil || a.running.Name != "self-update") && (strings.Contains(hte.Msg, "curl error 28 ") || strings.Contains(hte.Msg, "Resolving timed out") || strings.Contains(hte.Msg, "Could not resolve host")) {
				io.WriteErrorV("<warning>If you intend to run Composer without connecting to the internet, run the command again prefixed with COMPOSER_DISABLE_NETWORK=1 to make Composer run in offline mode.</warning>", true, cio.Quiet)
			}
		}
		if reportStatus {
			io.WriteProgramStatus(cio.StatusError())
		}
		renderError(io, err)
		if a.running != nil {
			io.WriteErrorV("<info>"+escapeTags(synopsisFor(a.running))+"</info>", true, cio.Quiet)
			io.WriteErrorV("", true, cio.Quiet)
		}
		// Symfony uses the exception code as exit code; Composer maps
		// TransportException to Installer::ERROR_TRANSPORT_EXCEPTION.
		var te *httpdl.TransportError
		if errors.As(err, &te) {
			code = composer.ErrTransport
		} else if code == 0 {
			code = 1
		}
	}
	if err == nil && reportStatus {
		io.WriteProgramStatus(cio.StatusFromExitCode(code))
	}
	if code > 255 {
		code = 255
	}
	return code
}

// programStatusReported is the $reportStatus of Application::doRun, decided
// after Application::configureIO applied --ansi/--no-ansi and the quiet
// verbosity (not for the _complete input completion process).
func programStatusReported(io *cio.IO, argv []string) bool {
	decorated := io.Decorated
	defer func() { io.Decorated = decorated }()
	probe := io
	quiet := os.Getenv("SHELL_VERBOSITY") == "-1" && verbosityFrom(argv) == 0
	for _, arg := range argv {
		if arg == "--" {
			break
		}
		switch arg {
		case "--ansi":
			probe.Decorated = true
		case "--no-ansi":
			probe.Decorated = false
		case "-q", "--quiet":
			quiet = true
		case "_complete":
			return false
		}
	}
	if quiet {
		return false
	}
	return probe.IsProgramStatusEnabled()
}

func (a *App) findCommand(name string) (*Command, error) {
	if c, ok := a.byName[name]; ok {
		return c, nil
	}
	// abbreviations like Symfony: unique prefix match
	var matches []string
	for n := range a.byName {
		if strings.HasPrefix(n, name) {
			matches = append(matches, n)
		}
	}
	uniq := map[*Command]bool{}
	for _, m := range matches {
		uniq[a.byName[m]] = true
	}
	if len(uniq) == 1 {
		for c := range uniq {
			return c, nil
		}
	}
	if len(uniq) > 1 {
		sort.Strings(matches)
		return nil, &CLIError{Msg: fmt.Sprintf("Command \"%s\" is ambiguous.\nDid you mean one of these?\n    %s", name, strings.Join(matches, "\n    "))}
	}
	if pos := strings.LastIndex(name, ":"); pos >= 0 {
		// Symfony's Application::findNamespace (plugin_commands.go)
		if err := a.findNamespace(name[:pos]); err != nil {
			return nil, err
		}
	}
	return nil, &CLIError{Msg: a.notDefinedMessage(name)}
}

func (a *App) run(io *cio.IO, argv []string) (int, error) {
	// first pass: find command name and global flags affecting IO
	cmdName := ""
	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		if arg == "--" {
			if i+1 < len(argv) && cmdName == "" {
				cmdName = argv[i+1]
			}
			break
		}
		if strings.HasPrefix(arg, "-") {
			if (arg == "-d" || arg == "--working-dir") && i+1 < len(argv) {
				i++
			}
			continue
		}
		if cmdName == "" {
			cmdName = arg
		}
	}
	for _, arg := range argv {
		switch arg {
		case "--no-ansi":
			io.Decorated = false
		case "--ansi":
			io.Decorated = true
		case "-n", "--no-interaction":
			io.Interactive = false
		case "-q", "--quiet":
			io.Verbosity = cio.Quiet
			io.Interactive = false
		}
	}
	// Application::doRunComposer: COMPOSER_NO_INTERACTION is a PHP truthy value
	if v := os.Getenv("COMPOSER_NO_INTERACTION"); v != "" && v != "0" {
		io.Interactive = false
	}
	if cmdName == "" {
		for _, arg := range argv {
			if arg == "-V" || arg == "--version" {
				io.Write(versionLine())
				return 0, nil
			}
		}
		cmdName = "list"
	}
	// Application::doRunComposer: a composer.json of a parent directory
	// may be used (parentdir.go)
	if err := a.useParentDir(io, argv, cmdName); err != nil {
		return 1, err
	}
	// Application::doRunComposer adds plugin commands (plugin_commands.go)
	// and composer.json scripts as commands
	if err := a.addPluginCommands(io, argv, cmdName); err != nil {
		return 1, err
	}
	if err := a.addScriptCommands(io, argv, cmdName); err != nil {
		return 1, err
	}
	cmd, err := a.findCommand(cmdName)
	if err != nil {
		return 1, err
	}
	// super user checks (plugin_commands.go)
	if !rootChecks(io, argv, cmd) {
		return 1, nil
	}
	a.running = cmd
	in, err := parseInput(cmd, argv, cmdName)
	strictErr := err
	if err != nil && cmd.IgnoreValidationErrors {
		in, err = parseInputLenient(cmd, argv, cmdName)
	}
	if err != nil {
		return 1, err
	}
	if in.Bool("version") {
		io.Write(versionLine())
		return 0, nil
	}
	if in.Bool("help") {
		return a.renderHelp(io, cmd), nil
	}
	if in.missing != nil && cmd.Name == "global" {
		// GlobalCommand overrides Command::run and fails before initialize() (after --version/--help);
		// without a sub-command it binds its input strictly (option errors
		// come before the missing argument)
		if strictErr != nil {
			return 1, strictErr
		}
		return 1, in.missing
	}
	if in.Bool("quiet") {
		io.Verbosity = cio.Quiet
	} else if v := verbosityFrom(argv); v > 0 {
		io.Verbosity = v
	} else if sv := os.Getenv("SHELL_VERBOSITY"); sv != "" {
		switch sv {
		case "-1":
			io.Verbosity = cio.Quiet
		case "1":
			io.Verbosity = cio.Verbose
		case "2":
			io.Verbosity = cio.VeryVerbose
		case "3":
			io.Verbosity = cio.Debug
		}
	}
	if in.Bool("no-interaction") {
		io.Interactive = false
	}
	// Application::configureIO: quiet runs are not interactive
	if io.Verbosity == cio.Quiet {
		io.Interactive = false
	}
	if wd, ok := in.String("working-dir"); ok {
		if err := os.Chdir(wd); err != nil {
			return 1, &composer.UserError{Msg: "Invalid working directory specified, " + wd + " does not exist.", Code: 1}
		}
	}
	if in.Bool("no-cache") {
		os.Setenv("COMPOSER_CACHE_DIR", os.DevNull)
	}
	applyEnvOptions(io, cmd, in)
	ctx := &Context{IO: io, In: in, App: a, Cmd: cmd}
	if cmd.Name != "list" && cmd.Name != "help" && cmd.Name != "_complete" {
		// BaseCommand::initialize() -> tryComposer(): the instance is created
		// (with its warnings) before the command runs; Application::getComposer
		// swallows Runtime/JsonValidation/InvalidArgument exceptions there,
		// while JSON parse errors abort the command.
		if _, err := ctx.composerFor(true); err != nil {
			var pe *jsonfile.ParsingError
			if errors.As(err, &pe) {
				return 1, err
			}
			ctx.cache = nil
		} else if err := ctx.pluginStartupEvents(); err != nil {
			return 1, err
		}
	}
	if in.missing != nil {
		return 1, in.missing
	}
	return cmd.Run(ctx)
}

func verbosityFrom(argv []string) int {
	v := 0
	for _, a := range argv {
		switch a {
		case "-v", "--verbose":
			if v < cio.Verbose {
				v = cio.Verbose
			}
		case "-vv":
			if v < cio.VeryVerbose {
				v = cio.VeryVerbose
			}
		case "-vvv":
			v = cio.Debug
		default:
			if strings.HasPrefix(a, "--verbose=") {
				switch a[10:] {
				case "1":
					v = cio.Verbose
				case "2":
					v = cio.VeryVerbose
				case "3":
					v = cio.Debug
				}
			}
		}
	}
	return v
}

func versionLine() string {
	return "<info>gocomposer</info> version <comment>" + Version + "</comment> (Composer-compatible, based on Composer " + platform.ComposerVersion + " f83588392df43d4fd027d84d26de459946ef4e80)"
}

func allOptions(cmd *Command) []OptDef {
	return append(append([]OptDef(nil), cmd.Options...), globalOptions...)
}

func parseInput(cmd *Command, argv []string, cmdToken string) (*Input, error) {
	in := &Input{opts: map[string][]string{}, present: map[string]bool{}, args: map[string][]string{}, Raw: argv}
	opts := allOptions(cmd)
	byLong := map[string]OptDef{}
	byShort := map[string]OptDef{}
	for _, o := range opts {
		byLong[o.Name] = o
		if o.Shortcut != "" {
			for _, s := range strings.Split(o.Shortcut, "|") {
				byShort[s] = o
			}
		}
	}
	var positional []string
	seenCmd := false
	onlyArgs := false
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if onlyArgs {
			positional = append(positional, a)
			continue
		}
		if a == "--" {
			onlyArgs = true
			continue
		}
		if strings.HasPrefix(a, "--") {
			name, val, hasVal := strings.Cut(a[2:], "=")
			o, ok := byLong[name]
			if !ok && strings.HasPrefix(name, "no-") {
				if neg, nok := byLong[name[3:]]; nok && neg.Negatable {
					if hasVal {
						return nil, &CLIError{Msg: fmt.Sprintf(`The "--%s" option does not accept a value.`, name)}
					}
					in.present[name] = true
					continue
				}
			}
			if !ok {
				return nil, &CLIError{Msg: fmt.Sprintf(`The "--%s" option does not exist.`, name)}
			}
			in.present[name] = true
			switch o.Kind {
			case OptNone:
				if hasVal {
					return nil, &CLIError{Msg: fmt.Sprintf(`The "--%s" option does not accept a value.`, name)}
				}
			case OptRequired, OptArray:
				if !hasVal {
					if i+1 < len(argv) && !strings.HasPrefix(argv[i+1], "-") {
						val = argv[i+1]
						i++
					} else {
						return nil, &CLIError{Msg: fmt.Sprintf(`The "--%s" option requires a value.`, name)}
					}
				}
				in.opts[name] = append(in.opts[name], val)
			case OptOptional:
				if hasVal {
					in.opts[name] = append(in.opts[name], val)
				} else if i+1 < len(argv) && !strings.HasPrefix(argv[i+1], "-") {
					in.opts[name] = append(in.opts[name], argv[i+1])
					i++
				}
			}
			continue
		}
		if strings.HasPrefix(a, "-") && len(a) > 1 {
			short := a[1:]
			if o, ok := byShort[short]; ok && o.Name == "verbose" {
				in.present["verbose"] = true
				continue
			}
			for j := 0; j < len(short); j++ {
				o, ok := byShort[string(short[j])]
				if !ok {
					return nil, &CLIError{Msg: fmt.Sprintf(`The "-%s" option does not exist.`, string(short[j]))}
				}
				in.present[o.Name] = true
				if o.Kind == OptRequired || o.Kind == OptArray {
					val := short[j+1:]
					if val == "" {
						if i+1 < len(argv) {
							val = argv[i+1]
							i++
						} else {
							return nil, &CLIError{Msg: fmt.Sprintf(`The "--%s" option requires a value.`, o.Name)}
						}
					}
					in.opts[o.Name] = append(in.opts[o.Name], val)
					break
				}
				if o.Kind == OptOptional && j+1 < len(short) {
					in.opts[o.Name] = append(in.opts[o.Name], short[j+1:])
					break
				}
			}
			continue
		}
		if !seenCmd && a == cmdToken {
			seenCmd = true
			continue
		}
		positional = append(positional, a)
	}
	for idx, ad := range cmd.Args {
		if ad.Array {
			if idx < len(positional) {
				in.args[ad.Name] = positional[idx:]
			} else if ad.Required && in.missing == nil {
				in.missing = &CLIError{Msg: fmt.Sprintf(`Not enough arguments (missing: "%s").`, ad.Name)}
			}
			positional = nil
			break
		}
		if idx < len(positional) {
			in.args[ad.Name] = []string{positional[idx]}
		} else if ad.Required && in.missing == nil {
			in.missing = &CLIError{Msg: fmt.Sprintf(`Not enough arguments (missing: "%s").`, ad.Name)}
		}
	}
	nonArray := 0
	hasArray := false
	for _, ad := range cmd.Args {
		if ad.Array {
			hasArray = true
		} else {
			nonArray++
		}
	}
	if !hasArray && len(positional) > nonArray {
		extra := positional[nonArray]
		if len(cmd.Args) == 0 {
			return nil, &CLIError{Msg: fmt.Sprintf(`No arguments expected for "%s" command, got "%s".`, cmd.Name, extra)}
		}
		return nil, &CLIError{Msg: fmt.Sprintf(`Too many arguments to "%s" command, expected arguments "%s".`, cmd.Name, argNames(cmd))}
	}
	for _, o := range opts {
		if _, ok := in.opts[o.Name]; !ok && o.Default != "" && (o.Kind == OptRequired || o.Kind == OptOptional) {
			in.opts[o.Name] = []string{o.Default}
		}
	}
	return in, nil
}

func argNames(cmd *Command) string {
	var n []string
	for _, a := range cmd.Args {
		n = append(n, a.Name)
	}
	return strings.Join(n, `" "`)
}

// previousError is implemented by errors that carry a PHP-style previous
// exception, which Symfony renders as a further box.
type previousError interface{ Previous() error }

// renderError writes an exception box like Symfony's Application::renderThrowable.
func renderError(io *cio.IO, err error) {
	io.WriteErrorV("", true, cio.Quiet)
	for e := err; e != nil; {
		renderErrorBox(io, e)
		pe, ok := e.(previousError)
		if !ok {
			break
		}
		e = pe.Previous()
	}
}

func renderErrorBox(io *cio.IO, err error) {
	msg := strings.TrimSpace(err.Error())
	var ue *composer.UserError
	var ce *CLIError
	isConsole := errors.As(err, &ce)
	// Symfony Console exceptions (ExceptionInterface) get no "In file line" title
	var mi *cio.MissingInputError
	var cre *cio.RuntimeError
	var cie *cio.InvalidArgumentError
	isConsole = isConsole || errors.As(err, &mi) || errors.As(err, &cre) || (errors.As(err, &cie) && cie.Console)
	_ = errors.As(err, &ue)
	width := terminalWidth() - 1
	var lines []string
	maxLen := 0
	for _, line := range regexp.MustCompile(`\r?\n`).Split(msg, -1) {
		for _, part := range splitByWidth(line, width-4) {
			lines = append(lines, part)
			if l := len([]rune(part)) + 4; l > maxLen {
				maxLen = l
			}
		}
	}
	var out []string
	if !isConsole {
		out = append(out, "<comment>In gocomposer:</comment>")
	}
	empty := "<error>" + strings.Repeat(" ", maxLen) + "</error>"
	out = append(out, empty)
	for _, l := range lines {
		out = append(out, "<error>  "+escapeTags(l)+"  "+strings.Repeat(" ", maxLen-len([]rune(l))-4)+"</error>")
	}
	out = append(out, empty, "")
	for _, l := range out {
		io.WriteErrorV(l, true, cio.Quiet)
	}
}

func escapeTags(s string) string { return strings.ReplaceAll(s, "<", `\<`) }

func splitByWidth(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	r := []rune(s)
	if len(r) == 0 {
		return []string{""}
	}
	var out []string
	for len(r) > width {
		out = append(out, string(r[:width]))
		r = r[width:]
	}
	return append(out, string(r))
}

func terminalWidth() int { return cio.TerminalWidth() }

// composerFor creates the Composer instance for a command.
func (ctx *Context) composerFor(fullLoad bool) (*composer.Composer, error) {
	if ctx.cache != nil {
		return ctx.cache, nil
	}
	if c := takePreloadedComposer(); c != nil {
		ctx.cache = c
		return c, nil
	}
	c, err := composer.Create(ctx.IO, composer.Options{DisablePlugins: ctx.In.Bool("no-plugins") || pluginsDisabledByRoot, DisableScripts: ctx.In.Bool("no-scripts"), FullLoad: fullLoad})
	if err != nil {
		return nil, err
	}
	ctx.cache = c
	return c, nil
}

// platformFilter builds the platform requirement filter from options.
func platformFilter(in *Input) solver.PlatformFilter {
	if in.Bool("ignore-platform-reqs") {
		return solver.IgnoreAll{}
	}
	if reqs := in.Strings("ignore-platform-req"); len(reqs) > 0 {
		return solver.NewIgnoreList(reqs)
	}
	return solver.IgnoreNothing{}
}
