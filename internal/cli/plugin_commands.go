package cli

// Plugin integration of the application (Application::doRunComposer and
// BaseCommand::initialize of Composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT license, (c) Nils Adermann, Jordi Boggiano): commands provided by
// plugins, the pre-command-run and command events.

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
)

// preloadedComposer is the Composer instance created to look up plugin
// commands; the command then reuses it like Application::getComposer.
var preloadedComposer *composer.Composer

// pluginCommands marks the commands provided by plugins.
var pluginCommands = map[*Command]bool{}

func argvHas(argv []string, names ...string) bool {
	for _, t := range argv {
		if t == "--" {
			return false
		}
		for _, n := range names {
			if t == n {
				return true
			}
		}
	}
	return false
}

// addPluginCommands adds the commands of the active plugins' CommandProvider
// capabilities when the command line may need them (unknown command, list,
// help), like Application::doRunComposer.
func (a *App) addPluginCommands(io *cio.IO, argv []string, cmdName string) error {
	a.setPluginApp(argv, cmdName)
	if argvHas(argv, "--version", "-V") || argvHas(argv, "--no-plugins") {
		return nil
	}
	if c, err := a.findCommand(cmdName); err == nil && c.Name != "list" && c.Name != "help" {
		return nil
	}
	if err := rootPluginGuard(io); err != nil {
		return err
	}
	if wd := workingDirArg(argv); wd != "" {
		prev, err := os.Getwd()
		if err != nil {
			return nil
		}
		abs, err := filepath.Abs(wd)
		if err != nil || os.Chdir(abs) != nil {
			return nil
		}
		defer os.Chdir(prev)
	}
	file, err := composer.ComposerFileName()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(file); err != nil {
		return nil
	}
	// the IO is configured from the input before plugins are loaded
	saved := io.Verbosity
	if argvHas(argv, "-q", "--quiet") {
		io.Verbosity = cio.Quiet
	} else if v := verbosityFrom(argv); v > 0 {
		io.Verbosity = v
	}
	c, err := composer.Create(io, composer.Options{DisableScripts: argvHas(argv, "--no-scripts"), FullLoad: true})
	io.Verbosity = saved
	if err != nil {
		var pe *jsonfile.ParsingError
		if errors.As(err, &pe) {
			return err
		}
		return nil
	}
	preloadedComposer = c
	cmds, err := composer.PluginCommands(c)
	if err != nil {
		return err
	}
	for _, pc := range cmds {
		if _, exists := a.byName[pc.Name]; exists {
			io.WriteError("<warning>Plugin command " + pc.Name + " (" + pc.Class + ") would override a Composer command and has been skipped</warning>")
			continue
		}
		cmd := pluginCommand(pc)
		pluginCommands[cmd] = true
		a.commands = append(a.commands, cmd)
		a.byName[cmd.Name] = cmd
		for _, al := range cmd.Aliases {
			if _, exists := a.byName[al]; !exists {
				a.byName[al] = cmd
			}
		}
	}
	return nil
}

func pluginCommand(pc composer.PluginCommand) *Command {
	cmd := &Command{Name: pc.Name, Aliases: pc.Aliases, Description: pc.Description, Help: pc.Help, Hidden: pc.Hidden, Synopsis: pc.Synopsis, Usages: pc.Usages}
	for _, a := range pc.Args {
		cmd.Args = append(cmd.Args, ArgDef{Name: a.Name, Required: a.Required, Array: a.Array, Desc: a.Desc, RawDefault: a.Default})
	}
	for _, o := range pc.Opts {
		od := OptDef{Name: o.Name, Shortcut: o.Shortcut, Desc: o.Desc, Negatable: o.Negatable, RawDefault: o.Default}
		switch {
		case o.Array:
			od.Kind = OptArray
		case o.Required:
			od.Kind = OptRequired
		case o.Optional:
			od.Kind = OptOptional
		default:
			od.Kind = OptNone
		}
		if s, ok := o.Default.(string); ok && (od.Kind == OptRequired || od.Kind == OptOptional) {
			od.Default = s
		}
		cmd.Options = append(cmd.Options, od)
	}
	index := pc.Index
	cmd.Run = func(ctx *Context) (int, error) {
		c, err := ctx.composerFor(true)
		if err != nil {
			return 1, err
		}
		return composer.RunPluginCommand(c, index, ctx.In.Raw)
	}
	return cmd
}

// takePreloadedComposer hands the Composer instance created for the plugin
// commands to the command context.
func takePreloadedComposer() *composer.Composer {
	c := preloadedComposer
	preloadedComposer = nil
	return c
}

// earlyCommandEvents are the commands dispatching PluginEvents::COMMAND
// right after obtaining the Composer instance; install, update, require and
// remove dispatch it from checkPlugins at their own point.
var earlyCommandEvents = map[string]bool{
	"archive": true, "licenses": true, "dump-autoload": true, "depends": true, "prohibits": true,
	"diagnose": true, "status": true, "show": true, "search": true,
}

// pluginStartupEvents dispatches pre-command-run (BaseCommand::initialize)
// and the early command events once the Composer instance exists.
func (ctx *Context) pluginStartupEvents() error {
	c := ctx.cache
	if c == nil {
		return nil
	}
	if err := ctx.applyPluginInput(); err != nil {
		return err
	}
	d := composer.NewEventDispatcher(c)
	d.RunScripts = !c.DisableScripts
	if !pluginCommands[ctx.Cmd] {
		if err := d.DispatchPreCommandRun(ctx.Cmd.Name, commandInput(ctx)); err != nil {
			return err
		}
	}
	if earlyCommandEvents[ctx.Cmd.Name] {
		return d.DispatchCommandEvent(ctx.Cmd.Name, commandInput(ctx))
	}
	return nil
}

// commandInput describes the command line for the Symfony input of plugin
// events.
func commandInput(ctx *Context) *composer.CommandInput {
	in := &composer.CommandInput{Argv: ctx.In.Raw}
	for _, a := range ctx.Cmd.Args {
		mode := 2
		if a.Required {
			mode = 1
		}
		var def phpjson.Value
		if a.Array {
			mode |= 4
			def = []phpjson.Value{}
		} else if s, ok := a.RawDefault.(string); ok && !a.Required {
			def = s
		}
		in.Args = append(in.Args, composer.InputArgDef{Name: a.Name, Mode: mode, Desc: a.Desc, Default: def})
	}
	for _, o := range allOptions(ctx.Cmd) {
		var mode int
		var def phpjson.Value
		switch o.Kind {
		case OptNone:
			mode = 1
			if o.Negatable {
				mode |= 16
			}
		case OptRequired:
			mode = 2
		case OptOptional:
			mode = 4
		case OptArray:
			mode = 2 | 8
			def = []phpjson.Value{}
		}
		if (o.Kind == OptRequired || o.Kind == OptOptional) && o.Default != "" {
			def = o.Default
		}
		in.Opts = append(in.Opts, composer.InputOptDef{Name: o.Name, Shortcut: o.Shortcut, Mode: mode, Desc: o.Desc, Default: def})
	}
	return in
}

// commandNotFoundHint is Application::hintCommonErrors for missing commands
// while plugins are disabled.
func commandNotFoundHint(io *cio.IO, argv []string, err error) {
	if pluginsDisabledByDefault(argv) && os.Getuid() == 0 && !io.Interactive {
		io.WriteErrorV("<error>Plugins have been disabled automatically as you are running as root, this may be the cause of the following exception. See also https://getcomposer.org/root</error>", true, cio.Quiet)
		return
	}
	var ce *CLIError
	if !errors.As(err, &ce) || !pluginsDisabledByDefault(argv) {
		return
	}
	if strings.HasSuffix(ce.Msg, "\" is not defined.") || strings.Contains(ce.Msg, "Command \"") && strings.Contains(ce.Msg, "is not defined.") || strings.HasPrefix(ce.Msg, "There are no commands defined in the") {
		io.WriteErrorV("<error>Plugins have been disabled, which may be why some commands are missing, unless you made a typo</error>", true, cio.Quiet)
	}
}

// findNamespace is Symfony Application::findNamespace (commands provided by
// plugins and scripts may be namespaced like "plugin:command").
func (a *App) findNamespace(namespace string) error {
	all := map[string]bool{}
	for n := range a.byName {
		parts := strings.Split(n, ":")
		for i := 1; i < len(parts); i++ {
			all[strings.Join(parts[:i], ":")] = true
		}
	}
	var quoted []string
	for _, p := range strings.Split(namespace, ":") {
		quoted = append(quoted, regexpQuote(p))
	}
	re, err := regexpCompile(`^` + strings.Join(quoted, `[^:]*:`) + `[^:]*`)
	if err != nil {
		return nil
	}
	var names []string
	for n := range all {
		names = append(names, n)
		if re.MatchString(n) {
			return nil
		}
	}
	msg := `There are no commands defined in the "` + namespace + `" namespace.`
	if alts := findAlternatives(namespace, names); len(alts) > 0 {
		if len(alts) == 1 {
			msg += "\n\nDid you mean this?\n    "
		} else {
			msg += "\n\nDid you mean one of these?\n    "
		}
		msg += strings.Join(alts, "\n    ")
	}
	return &CLIError{Msg: msg}
}

// findAlternatives is Symfony Application::findAlternatives.
func findAlternatives(name string, collection []string) []string {
	const threshold = 1000
	alternatives := map[string]int{}
	collectionParts := map[string][]string{}
	for _, item := range collection {
		collectionParts[item] = strings.Split(item, ":")
	}
	for i, subname := range strings.Split(name, ":") {
		for _, cname := range collection {
			parts := collectionParts[cname]
			_, exists := alternatives[cname]
			if i >= len(parts) {
				if exists {
					alternatives[cname] += threshold
				}
				continue
			}
			lev := levenshtein(subname, parts[i])
			if float64(lev) <= float64(len(subname))/3 || subname != "" && strings.Contains(parts[i], subname) {
				if exists {
					alternatives[cname] += lev
				} else {
					alternatives[cname] = lev
				}
			} else if exists {
				alternatives[cname] += threshold
			}
		}
	}
	for _, item := range collection {
		lev := levenshtein(name, item)
		if float64(lev) <= float64(len(name))/3 || strings.Contains(item, name) {
			if v, ok := alternatives[item]; ok {
				alternatives[item] = v - lev
			} else {
				alternatives[item] = lev
			}
		}
	}
	var out []string
	for k, v := range alternatives {
		if v < 2*threshold {
			out = append(out, k)
		}
	}
	sortNaturalFold(out)
	return out
}

// notDefinedMessage is the CommandNotFoundException message of Symfony's
// Application::find with its alternatives (hidden commands left out).
func (a *App) notDefinedMessage(name string) string {
	msg := `Command "` + name + `" is not defined.`
	var all []string
	for n := range a.byName {
		all = append(all, n)
	}
	var alts []string
	for _, n := range findAlternatives(name, all) {
		if c := a.byName[n]; c != nil && !c.Hidden {
			alts = append(alts, n)
		}
	}
	if len(alts) == 0 {
		return msg
	}
	if len(alts) == 1 {
		msg += "\n\nDid you mean this?\n    "
	} else {
		msg += "\n\nDid you mean one of these?\n    "
	}
	return msg + strings.Join(alts, "\n    ")
}

// pluginsDisabledByRoot is set when plugins were disabled because Composer
// runs as root/super user in a non-interactive session.
var pluginsDisabledByRoot bool

// isNonAllowedRoot is the super user check of Application::doRunComposer
// (posix_getuid() === 0, COMPOSER_ALLOW_SUPERUSER not set, not in a
// container).
func isNonAllowedRoot() bool {
	if runtime.GOOS == "windows" || os.Getuid() != 0 {
		return false
	}
	if v := os.Getenv("COMPOSER_ALLOW_SUPERUSER"); v != "" && v != "0" {
		return false
	}
	return !isDocker()
}

// isDocker is Platform::isDocker.
func isDocker() bool {
	for _, f := range []string{"/.dockerenv", "/run/.containerenv", "/var/run/.containerenv"} {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}
	for _, f := range []string{"/proc/self/mountinfo", "/proc/1/cgroup"} {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		s := string(data)
		if strings.Contains(s, "/var/lib/docker/") || strings.Contains(s, "/io.containerd.snapshotter") {
			return true
		}
	}
	return false
}

func askContinueAsRoot(io *cio.IO) bool {
	ok, err := io.AskConfirmation("<info>Continue as root/super user</info> [<comment>yes</comment>]? ", true)
	return err == nil && ok
}

// rootPluginGuard runs before plugin commands are loaded: plugins must not be
// loaded as root unless explicitly allowed.
func rootPluginGuard(io *cio.IO) error {
	if !isNonAllowedRoot() {
		return nil
	}
	io.WriteError("<warning>Do not run Composer as root/super user! See https://getcomposer.org/root for details</warning>")
	if io.Interactive && askContinueAsRoot(io) {
		rootConfirmed = true
		return nil
	}
	io.WriteError("<warning>Aborting as no plugin should be loaded if running as super user is not explicitly allowed</warning>")
	return &composer.ScriptError{Code: 1}
}

// rootConfirmed avoids a second prompt once the user agreed.
var rootConfirmed bool

// rootChecks is the rest of the super user handling of
// Application::doRunComposer once the command is known; it returns false
// when the user declined to continue.
func rootChecks(io *cio.IO, argv []string, cmd *Command) bool {
	if !isNonAllowedRoot() || rootConfirmed {
		return true
	}
	if !argvHas(argv, "--no-plugins") && !io.Interactive {
		io.WriteError("<error>Composer plugins have been disabled for safety in this non-interactive session.</error>")
		io.WriteError("<error>Set COMPOSER_ALLOW_SUPERUSER=1 if you want to allow plugins to run as root/super user.</error>")
		pluginsDisabledByRoot = true
	}
	if cmd.Name != "self-update" && cmd.Name != "selfupdate" && cmd.Name != "_complete" {
		io.WriteError("<warning>Do not run Composer as root/super user! See https://getcomposer.org/root for details</warning>")
		if io.Interactive && !askContinueAsRoot(io) {
			return false
		}
	}
	return true
}

func pluginsDisabledByDefault(argv []string) bool {
	return pluginsDisabledByRoot || argvHas(argv, "--no-plugins")
}

// rootScriptHint is printed when a script fails while plugins were disabled
// as root (Application::doRun).
func rootScriptHint(io *cio.IO, argv []string) {
	if pluginsDisabledByDefault(argv) && os.Getuid() == 0 && !io.Interactive {
		io.WriteErrorV("<error>Plugins have been disabled automatically as you are running as root, this may be the cause of the script failure.</error>", true, cio.Quiet)
		io.WriteErrorV("<error>See also https://getcomposer.org/root</error>", true, cio.Quiet)
	}
}
