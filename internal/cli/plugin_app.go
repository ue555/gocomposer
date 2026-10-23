package cli

// The command line as seen by plugin code: Composer runs plugins inside
// Application::doRun($input) (Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Console/Application.php, MIT license, (c) Nils Adermann, Jordi
// Boggiano), where plugins such as symfony/flex find the application and the
// command input and may rewrite the input (package aliases) before the
// command reads it.

import (
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
)

// setPluginApp describes the command line about to run to the plugin
// runtime (before plugins are loaded for it).
func (a *App) setPluginApp(argv []string, cmdName string) {
	global := false
	if prev := composer.CurrentPluginApp(); prev != nil && prev.Cmd != nil && prev.Cmd.Name == "global" {
		global = true
	}
	app := &composer.PluginApp{Argv: append([]string(nil), argv...), Global: global, Find: a.pluginFind}
	if c, err := a.findCommand(cmdName); err == nil && c.Name != "list" && c.Name != "help" {
		app.Cmd = describeForPlugins(c)
	}
	composer.SetPluginApp(app)
}

func (a *App) pluginFind(name string) (*composer.PluginAppCommand, error) {
	c, err := a.findCommand(name)
	if err != nil {
		return nil, err
	}
	return describeForPlugins(c), nil
}

func describeForPlugins(c *Command) *composer.PluginAppCommand {
	pc := &composer.PluginAppCommand{Name: c.Name, Aliases: c.Aliases, Description: c.Description}
	for _, arg := range c.Args {
		pc.Args = append(pc.Args, composer.PluginCommandArg{Name: arg.Name, Required: arg.Required, Array: arg.Array, Desc: arg.Desc, Default: toPHP(arg.RawDefault)})
	}
	for _, o := range c.Options {
		po := composer.PluginCommandOpt{Name: o.Name, Shortcut: o.Shortcut, Desc: o.Desc, Negatable: o.Negatable, Default: toPHP(o.RawDefault)}
		switch o.Kind {
		case OptArray:
			po.Array = true
		case OptRequired:
			po.Required = true
		case OptOptional:
			po.Optional = true
		default:
			po.None = true
		}
		if po.Default == nil && o.Default != "" && (o.Kind == OptRequired || o.Kind == OptOptional) {
			po.Default = o.Default
		}
		pc.Opts = append(pc.Opts, po)
	}
	return pc
}

func toPHP(v interface{}) phpjson.Value {
	switch t := v.(type) {
	case nil, bool, string, int64, float64, *phpjson.Map, []phpjson.Value:
		return t
	case int:
		return int64(t)
	case []string:
		out := []phpjson.Value{}
		for _, s := range t {
			out = append(out, s)
		}
		return out
	}
	return nil
}

// applyPluginInput re-reads the command input once plugins were activated:
// plugins may have rewritten it (symfony/flex resolves package aliases and
// versions of require/update/remove arguments in place).
func (ctx *Context) applyPluginInput() error {
	argv, changed, err := composer.PluginAppArgv(ctx.cache)
	if err != nil || !changed {
		return err
	}
	cmdToken := ""
	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		if arg == "--" {
			if i+1 < len(argv) {
				cmdToken = argv[i+1]
			}
			break
		}
		if strings.HasPrefix(arg, "-") {
			if (arg == "-d" || arg == "--working-dir") && i+1 < len(argv) {
				i++
			}
			continue
		}
		cmdToken = arg
		break
	}
	in, err := parseInput(ctx.Cmd, argv, cmdToken)
	if err != nil {
		return err
	}
	ctx.In = in
	return nil
}
