package cli

// Port of Composer\Command\OutdatedCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// The command is a proxy: it builds the arguments of `show --latest` and runs
// that command through the application, like OutdatedCommand::execute does
// with an ArrayInput.

func init() { registerCommand(outdatedCommand) }

func outdatedCommand() *Command {
	return &Command{Name: "outdated", Run: runOutdated}
}

func runOutdated(ctx *Context) (int, error) {
	in := ctx.In
	args := []string{"show", "--latest"}
	for _, o := range []string{"no-interaction", "no-plugins", "no-scripts", "no-cache"} {
		if in.Bool(o) {
			args = append(args, "--"+o)
		}
	}
	if !in.Bool("all") {
		args = append(args, "--outdated")
	}
	if in.Bool("direct") {
		args = append(args, "--direct")
	}
	for _, o := range []string{"strict", "major-only", "minor-only", "patch-only", "locked", "no-dev", "sort-by-age"} {
		if in.Bool(o) {
			args = append(args, "--"+o)
		}
	}
	if v, ok := in.String("cooldown-period"); ok {
		args = append(args, "--cooldown-period="+v)
	}
	for _, r := range in.Strings("ignore-platform-req") {
		args = append(args, "--ignore-platform-req="+r)
	}
	if in.Bool("ignore-platform-reqs") {
		args = append(args, "--ignore-platform-reqs")
	}
	format, ok := in.String("format")
	if !ok {
		format = "text"
	}
	args = append(args, "--format="+format)
	for _, i := range in.Strings("ignore") {
		args = append(args, "--ignore="+i)
	}
	// the package argument comes after "--" so that names starting with
	// a dash are not parsed as options, as with ArrayInput
	if pkgArg := in.Arg("package"); pkgArg != "" {
		args = append(args, "--", pkgArg)
	}

	app := ctx.App
	show, err := app.findCommand("show")
	if err != nil {
		return 1, err
	}
	app.running = show
	showIn, err := parseInput(show, args, "show")
	if err != nil {
		return 1, err
	}
	applyEnvOptions(ctx.IO, show, showIn)
	// the Application keeps its Composer instance across the nested run
	return show.Run(&Context{IO: ctx.IO, In: showIn, App: app, Cmd: show, cache: ctx.cache})
}
