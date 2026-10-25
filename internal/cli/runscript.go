package cli

// Port of Composer\Command\RunScriptCommand, Composer\Command\ScriptAliasCommand
// and the script command registration of Composer\Console\Application::doRun
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).
// Scripts naming a Symfony Command class are registered as that command
// (described and run by the PHP command runner, internal/composer/scriptcmd.go).

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
)

func init() { registerCommand(runScriptCommand) }

// scriptEventNames are the Composer\Script\ScriptEvents constants (values).
var scriptEventNames = []string{
	"pre-install-cmd", "post-install-cmd", "pre-update-cmd", "post-update-cmd",
	"pre-status-cmd", "post-status-cmd", "pre-autoload-dump", "post-autoload-dump",
	"post-root-package-install", "post-create-project-cmd", "pre-archive-cmd", "post-archive-cmd",
}

// isScriptEventConstant is defined('Composer\Script\ScriptEvents::'.str_replace('-', '_', strtoupper($name))).
func isScriptEventConstant(name string) bool {
	c := strings.ReplaceAll(strings.ToUpper(name), "-", "_")
	for _, e := range scriptEventNames {
		if strings.ReplaceAll(strings.ToUpper(e), "-", "_") == c {
			return true
		}
	}
	return false
}

func runScriptCommand() *Command {
	return &Command{Name: "run-script", Run: runRunScript}
}

var reCtypeDigit = regexp.MustCompile(`^[0-9]+$`)

func runRunScript(ctx *Context) (int, error) {
	in := ctx.In
	if ctx.IO.Interactive {
		// Command::run calls interact() before execute() (select_cmds.go)
		if err := runScriptInteract(ctx); err != nil {
			return exitFor(err)
		}
	}
	if in.Bool("list") {
		return listScripts(ctx)
	}
	script := in.Arg("script")
	if _, ok := in.args["script"]; !ok {
		return 1, &composer.UserError{Msg: `Missing required argument "script"`, Code: 1}
	}
	isEvent := false
	for _, e := range scriptEventNames {
		if e == script {
			isEvent = true
		}
	}
	if !isEvent && isScriptEventConstant(script) {
		return 1, &composer.UserError{Msg: fmt.Sprintf(`Script "%s" cannot be run with this command`, script), Code: 1}
	}
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	devMode := in.Bool("dev") || !in.Bool("no-dev")
	d := composer.NewEventDispatcher(c)
	d.RunScripts = !c.DisableScripts
	if !d.HasListeners(script) {
		return 1, &composer.UserError{Msg: fmt.Sprintf(`Script "%s" is not defined in this package`, script), Code: 1}
	}
	if t, ok := in.String("timeout"); ok && in.present["timeout"] {
		if !reCtypeDigit.MatchString(t) {
			return 1, &composer.UserError{Msg: "Timeout value must be numeric and positive if defined, or 0 for forever", Code: 1}
		}
		var n int64
		fmt.Sscan(t, &n)
		d.Timeout = secondsDuration(n)
	}
	setDevModeEnv(devMode)
	return d.DispatchScript(script, devMode, in.Args("args"), nil)
}

func setDevModeEnv(devMode bool) {
	if devMode {
		os.Setenv("COMPOSER_DEV_MODE", "1")
	} else {
		os.Setenv("COMPOSER_DEV_MODE", "0")
	}
}

func listScripts(ctx *Context) (int, error) {
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	scripts := c.Root().Scripts()
	if scripts == nil || scripts.Len() == 0 {
		return 0, nil
	}
	ctx.IO.WriteError("<info>scripts:</info>")
	t := cio.NewTable("compact")
	for _, name := range scripts.Keys() {
		desc := ""
		if cmd, err := ctx.App.findCommand(name); err == nil {
			desc = cmd.Description
		}
		t.Rows = append(t.Rows, []string{"  " + name, desc})
	}
	t.Render(ctx.IO)
	return 0, nil
}

// addScriptCommands registers the non-event scripts of composer.json as
// commands, under the conditions of Application::doRunComposer.
func (a *App) addScriptCommands(io *cio.IO, argv []string, cmdName string) error {
	hasVersion := false
	for _, t := range argv {
		if t == "--" {
			break
		}
		if t == "--version" || t == "-V" {
			hasVersion = true
		}
	}
	mayNeed := false
	if c, err := a.findCommand(cmdName); err != nil {
		mayNeed = !hasVersion
	} else if c.Name == "list" || c.Name == "help" {
		mayNeed = !hasVersion
	} else if c.Name == "run-script" || c.Name != cmdName {
		mayNeed = true
	}
	if !mayNeed {
		return nil
	}
	file, err := composer.ComposerFileName()
	if err != nil {
		return nil
	}
	if wd := workingDirArg(argv); wd != "" && !strings.HasPrefix(file, "/") {
		file = wd + "/" + file
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	v, err := phpjson.Decode(raw)
	if err != nil || !phpjson.IsArray(v) {
		return nil
	}
	root := phpjson.ToMap(v)
	scriptsV, _ := root.Get("scripts")
	if !phpjson.IsArray(scriptsV) {
		return nil
	}
	scripts := phpjson.ToMap(scriptsV)
	descs := phpjson.ToMap(mapGet(root, "scripts-descriptions"))
	aliasesMap := phpjson.ToMap(mapGet(root, "scripts-aliases"))
	scriptDescription := func(script string) string {
		description := "Runs the " + script + " script as defined in composer.json"
		if descs != nil {
			if d, ok := descs.Get(script); ok && d != nil {
				description, _ = phpjson.ScalarString(d)
			}
		}
		return description
	}
	classCommands := a.scriptCommandClasses(io, argv, scripts, scriptDescription)
	for _, script := range scripts.Keys() {
		if isScriptEventConstant(script) {
			continue
		}
		if _, ok := a.byName[script]; ok {
			io.WriteError("<warning>A script named " + script + " would override a Composer command and has been skipped</warning>")
			continue
		}
		description := scriptDescription(script)
		if sc, ok := classCommands[script]; ok {
			// the script points to a valid Command subclass: import its details directly
			for _, w := range sc.Warnings {
				io.WriteError(w)
			}
			v, _ := scripts.Get(script)
			className, _ := v.(string)
			cmd := scriptClassCommand(sc.Command, script, className)
			a.commands = append(a.commands, cmd)
			a.byName[cmd.Name] = cmd
			for _, al := range cmd.Aliases {
				a.byName[al] = cmd
			}
			continue
		}
		var aliases []string
		if aliasesMap != nil {
			if al, ok := aliasesMap.Get(script); ok && al != nil {
				var bad bool
				phpjson.Each(al, func(_ string, e phpjson.Value) {
					s, ok := e.(string)
					if !ok {
						bad = true
						return
					}
					aliases = append(aliases, s)
				})
				if bad {
					return &composer.UserError{Msg: `"scripts-aliases" element array values should contain only strings`, Code: 1}
				}
			}
		}
		cmd := scriptAliasCommand(script, description, aliases)
		a.commands = append(a.commands, cmd)
		a.byName[cmd.Name] = cmd
		for _, al := range aliases {
			a.byName[al] = cmd
		}
	}
	return nil
}

func mapGet(m *phpjson.Map, k string) phpjson.Value {
	if m == nil {
		return nil
	}
	v, _ := m.Get(k)
	if !phpjson.IsArray(v) {
		return nil
	}
	return v
}

// workingDirArg finds --working-dir/-d before the argument separator.
func workingDirArg(argv []string) string {
	for i := 0; i < len(argv); i++ {
		t := argv[i]
		if t == "--" {
			break
		}
		switch {
		case t == "-d" || t == "--working-dir":
			if i+1 < len(argv) {
				return argv[i+1]
			}
		case strings.HasPrefix(t, "--working-dir="):
			return t[len("--working-dir="):]
		case strings.HasPrefix(t, "-d") && !strings.HasPrefix(t, "--"):
			return t[2:]
		}
	}
	return ""
}

func scriptAliasCommand(script, description string, aliases []string) *Command {
	return &Command{
		Name:        script,
		Aliases:     aliases,
		Description: description,
		Help:        "The <info>run-script</info> command runs scripts defined in composer.json:\n\n<info>php composer.phar run-script post-update-cmd</info>\n\nRead more at https://getcomposer.org/doc/03-cli.md#run-script--run",
		Options: []OptDef{
			{Name: "dev", Desc: "Sets the dev mode."},
			{Name: "no-dev", Desc: "Disables the dev mode."},
		},
		Args:                   []ArgDef{{Name: "args", Array: true}},
		IgnoreValidationErrors: true,
		Run: func(ctx *Context) (int, error) {
			c, err := requireComposer(ctx)
			if err != nil {
				return exitFor(err)
			}
			devMode := ctx.In.Bool("dev") || !ctx.In.Bool("no-dev")
			setDevModeEnv(devMode)
			d := composer.NewEventDispatcher(c)
			d.RunScripts = !c.DisableScripts
			return d.DispatchScript(script, devMode, ctx.In.Args("args"), map[string]any{"script-alias-input": scriptAliasInput(ctx.In.Raw)})
		},
	}
}

var (
	reEqOption = regexp.MustCompile(`^(-[^=]+=)(.+)`)
	reWordTok  = regexp.MustCompile(`^[\w-]+$`)
	reFirstTok = regexp.MustCompile(`^\S+ ?`)
)

// scriptAliasInput is Preg::replace('{^\S+ ?}', ”, (string) $input, 1)
// with ArgvInput::__toString's token escaping.
func scriptAliasInput(argv []string) string {
	esc := func(t string) string {
		if reWordTok.MatchString(t) {
			return t
		}
		return "'" + strings.ReplaceAll(t, "'", `'\''`) + "'"
	}
	toks := make([]string, len(argv))
	for i, t := range argv {
		if m := reEqOption.FindStringSubmatch(t); m != nil {
			toks[i] = m[1] + esc(m[2])
		} else if t != "" && t != "0" && t[0] != '-' {
			toks[i] = esc(t)
		} else {
			toks[i] = t
		}
	}
	s := strings.Join(toks, " ")
	if loc := reFirstTok.FindStringIndex(s); loc != nil {
		s = s[loc[1]:]
	}
	return s
}

// parseInputLenient mirrors Symfony binding with ignored validation errors:
// tokens are bound up to the first invalid one.
func parseInputLenient(cmd *Command, argv []string, cmdToken string) (*Input, error) {
	for k := len(argv) - 1; k >= 0; k-- {
		in, err := parseInput(cmd, argv[:k], cmdToken)
		if err == nil {
			in.Raw = argv
			return in, nil
		}
	}
	return &Input{opts: map[string][]string{}, present: map[string]bool{}, args: map[string][]string{}, Raw: argv}, nil
}

func secondsDuration(n int64) time.Duration { return time.Duration(n) * time.Second }

// errComposerMissing marks that requireComposer already reported a missing
// composer.json (Application::getComposer writes the message and exits 1).
var errComposerMissing = errors.New("composer.json missing")

// requireComposer is BaseCommand::requireComposer: Factory's
// InvalidArgumentException for a missing config file is written as a plain
// message and the command exits with 1.
func requireComposer(ctx *Context) (*composer.Composer, error) {
	c, err := ctx.composerFor(true)
	if err == nil {
		return c, nil
	}
	var ue *composer.UserError
	if errors.As(err, &ue) && strings.HasPrefix(ue.Msg, "Composer could not find ") {
		ctx.IO.WriteError(ue.Msg)
		return nil, errComposerMissing
	}
	return nil, err
}

// exitFor converts a requireComposer error into a command result.
func exitFor(err error) (int, error) {
	if err == errComposerMissing {
		return 1, nil
	}
	return 1, err
}

// reClassName is a PHP class name as accepted by class_exists().
var reClassName = regexp.MustCompile(`^\\?[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*(?:\\[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*)*$`)

// scriptCommandClasses finds the scripts Application::doRun imports as
// commands: a string value naming an autoloadable Symfony Console command
// class (checked by PHP with the project autoloader of the Composer instance
// it creates for that, which the command then reuses).
func (a *App) scriptCommandClasses(io *cio.IO, argv []string, scripts *phpjson.Map, description func(string) string) map[string]composer.ScriptCommandClass {
	var names []string
	classes := map[string]string{}
	descriptions := map[string]string{}
	for _, script := range scripts.Keys() {
		if isScriptEventConstant(script) {
			continue
		}
		if _, ok := a.byName[script]; ok {
			continue
		}
		v, _ := scripts.Get(script)
		cls, ok := v.(string)
		if !ok || !reClassName.MatchString(cls) {
			continue
		}
		names = append(names, script)
		classes[script] = cls
		descriptions[script] = description(script)
	}
	if len(names) == 0 {
		return nil
	}
	c := preloadedComposer
	if c == nil {
		if wd := workingDirArg(argv); wd != "" {
			if prev, err := os.Getwd(); err == nil && os.Chdir(wd) == nil {
				defer os.Chdir(prev)
			}
		}
		created, err := composer.Create(io, composer.Options{DisablePlugins: pluginsDisabledByDefault(argv), DisableScripts: argvHas(argv, "--no-scripts"), FullLoad: true})
		if err != nil {
			return nil
		}
		c = created
		preloadedComposer = c
	}
	out, err := composer.DescribeScriptCommandClasses(c, names, classes, descriptions)
	if err != nil {
		return nil
	}
	return out
}

// scriptClassCommand is a script registered as its Symfony Console command
// class: it runs in the PHP command runner with Composer's global options.
func scriptClassCommand(pc composer.PluginCommand, script, className string) *Command {
	cmd := pluginCommand(pc)
	cmd.Run = func(ctx *Context) (int, error) {
		c, err := requireComposer(ctx)
		if err != nil {
			return exitFor(err)
		}
		return composer.RunScriptCommandClass(c, script, className, scriptAliasInput(ctx.In.Raw))
	}
	return cmd
}
