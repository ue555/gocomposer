package cli

import (
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
