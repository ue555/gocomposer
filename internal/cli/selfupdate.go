package cli

// Composer\Command\SelfUpdateCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano)
// replaces the running composer.phar. gocomposer is not a phar and has no
// self-update mechanism: the command accepts Composer's options and fails
// explicitly with exit code 1, like a Composer instance without self-update.

func init() { registerCommand(selfUpdateCommand) }

func selfUpdateCommand() *Command {
	return &Command{Name: "self-update", Run: func(ctx *Context) (int, error) {
		ctx.IO.Write("<error>This instance of Composer does not have the self-update command.</error>")
		ctx.IO.Write("<comment>gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.</comment>")
		return 1, nil
	}}
}
