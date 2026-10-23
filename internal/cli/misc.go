package cli

import (
	"fmt"
	"sort"
	"strings"

	"gocomposer/internal/composer"
)

// checkPlugins is called by install, update, require and remove where
// Composer dispatches PluginEvents::COMMAND (the plugins themselves are
// loaded with the Composer instance, see internal/composer/plugin.go). It
// fails explicitly when an active plugin relies on events gocomposer does not
// dispatch to plugins.
func checkPlugins(ctx *Context, c *composer.Composer) error {
	if err := composer.CheckUnsupportedPluginEvents(c); err != nil {
		return err
	}
	d := composer.NewEventDispatcher(c)
	d.RunScripts = !c.DisableScripts
	return d.DispatchCommandEvent(ctx.Cmd.Name, commandInput(ctx))
}

func notImplemented(name string) func(ctx *Context) (int, error) {
	return func(ctx *Context) (int, error) {
		return 1, &composer.UserError{Msg: "The " + name + " command is not implemented by gocomposer yet.", Code: 1}
	}
}

func aboutCommand() *Command {
	return &Command{Name: "about", Description: "Shows a short information about Composer", Run: func(ctx *Context) (int, error) {
		ctx.IO.Write("<info>gocomposer - Dependency Manager for PHP (Composer-compatible, written in Go)</info>\n<comment>gocomposer re-implements Composer's dependency resolution, lock files, installation and autoload generation in Go.\nIt is based on Composer " + "f83588392df43d4fd027d84d26de459946ef4e80" + " and is not affiliated with the Composer project.</comment>")
		return 0, nil
	}}
}

func listCommand() *Command {
	return &Command{Name: "list", Description: "List commands", Args: []ArgDef{{Name: "namespace"}}, Run: func(ctx *Context) (int, error) {
		io := ctx.IO
		io.Write(versionLine())
		io.Write("")
		io.Write("<comment>Usage:</comment>")
		io.Write("  command [options] [arguments]")
		io.Write("")
		io.Write("<comment>Options:</comment>")
		for _, o := range globalOptions {
			short := "    "
			if o.Shortcut != "" {
				short = "-" + strings.ReplaceAll(o.Shortcut, "|", "|-") + ", "
			}
			io.Write(fmt.Sprintf("  <info>%-25s</info> %s", short+"--"+o.Name, o.Desc))
		}
		io.Write("")
		io.Write("<comment>Available commands:</comment>")
		var cmds []*Command
		for _, c := range ctx.App.commands {
			if !c.Hidden {
				cmds = append(cmds, c)
			}
		}
		sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name < cmds[j].Name })
		for _, c := range cmds {
			io.Write(fmt.Sprintf("  <info>%-20s</info> %s", c.Name, c.Description))
		}
		return 0, nil
	}}
}

func helpCommand() *Command {
	return &Command{Name: "help", Description: "Display help for a command", Args: []ArgDef{{Name: "command_name"}}, Run: func(ctx *Context) (int, error) {
		name := ctx.In.Arg("command_name")
		if name == "" {
			name = "help"
		}
		cmd, err := ctx.App.findCommand(name)
		if err != nil {
			return 1, err
		}
		return ctx.App.renderHelp(ctx.IO, cmd), nil
	}}
}

func (a *App) renderHelp(io interface {
	Write(string)
}, cmd *Command) int {
	io.Write("<comment>Description:</comment>")
	io.Write("  " + cmd.Description)
	io.Write("")
	io.Write("<comment>Usage:</comment>")
	usage := "  " + cmd.Name + " [options]"
	for _, ar := range cmd.Args {
		if ar.Array {
			usage += " [--] [<" + ar.Name + ">...]"
		} else if ar.Required {
			usage += " <" + ar.Name + ">"
		} else {
			usage += " [<" + ar.Name + ">]"
		}
	}
	io.Write(usage)
	for _, al := range cmd.Aliases {
		io.Write("  " + al)
	}
	if len(cmd.Args) > 0 {
		io.Write("")
		io.Write("<comment>Arguments:</comment>")
		for _, ar := range cmd.Args {
			io.Write(fmt.Sprintf("  <info>%-28s</info> %s", ar.Name, ar.Desc))
		}
	}
	io.Write("")
	io.Write("<comment>Options:</comment>")
	for _, o := range allOptions(cmd) {
		short := "    "
		if o.Shortcut != "" {
			short = "-" + strings.ReplaceAll(o.Shortcut, "|", "|-") + ", "
		}
		name := short + "--" + o.Name
		if o.Kind == OptRequired || o.Kind == OptArray {
			name += "=" + strings.ToUpper(o.Name)
		}
		io.Write(fmt.Sprintf("  <info>%-28s</info> %s", name, o.Desc))
	}
	return 0
}
