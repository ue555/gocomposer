package cli

// Port of Composer\Command\ReinstallCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"os"
	"sort"

	"gocomposer/internal/composer"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

func init() { registerCommand(reinstallCommand) }

func reinstallCommand() *Command {
	return &Command{Name: "reinstall", Run: runReinstall}
}

func runReinstall(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	localRepo := c.LocalRepo
	var toReinstall []*pkg.Package
	names := map[string]bool{}
	types := in.Strings("type")
	patterns := in.Args("packages")
	if len(types) > 0 {
		if len(patterns) > 0 {
			return 1, &composer.UserError{Msg: "You cannot specify package names and filter by type at the same time.", Code: 1}
		}
		for _, p := range localRepo.CanonicalPackages() {
			for _, t := range types {
				if p.Type() == t {
					toReinstall = append(toReinstall, p)
					names[p.Name] = true
					break
				}
			}
		}
	} else {
		if len(patterns) == 0 {
			return 1, &composer.UserError{Msg: "You must pass one or more package names to be reinstalled.", Code: 1}
		}
		for _, pattern := range patterns {
			re := solver.PackageNameToRegexp(pattern)
			matched := false
			for _, p := range localRepo.CanonicalPackages() {
				if re.MatchString(p.Name) {
					matched = true
					toReinstall = append(toReinstall, p)
					names[p.Name] = true
				}
			}
			if !matched {
				io.WriteError(`<warning>Pattern "` + pattern + `" does not match any currently installed packages.</warning>`)
			}
		}
	}
	if len(toReinstall) == 0 {
		io.WriteError("<warning>Found no packages to reinstall, aborting.</warning>")
		return 1, nil
	}

	var uninstallOps []*solver.Operation
	for _, p := range toReinstall {
		uninstallOps = append(uninstallOps, &solver.Operation{Type: solver.OpUninstall, Package: p})
	}
	// install operations ordered by dependencies/plugins
	all := localRepo.MustPackages()
	var present []*pkg.Package
	for _, p := range all {
		if !names[p.Name] {
			present = append(present, p)
		}
	}
	installOps := solver.NewTransaction(present, all).Operations()
	installOrder := map[string]int{}
	for i, op := range installOps {
		if op.Type == solver.OpInstall && op.Package.AliasOf == nil {
			installOrder[op.Package.Name] = i
		}
	}
	sort.SliceStable(uninstallOps, func(i, j int) bool {
		return installOrder[uninstallOps[j].Package.Name]-installOrder[uninstallOps[i].Package.Name] < 0
	})

	events := composer.NewEventDispatcher(c)
	events.RunScripts = !c.DisableScripts
	if err := events.DispatchCommandEvent("reinstall", commandInput(ctx)); err != nil {
		return 1, err
	}
	ps, pd, err := preferredInstall(c, in)
	if err != nil {
		return 1, err
	}
	im := composer.NewInstallationManager(c)
	im.OutputProgress = !in.Bool("no-progress")
	im.SetPreferSource(ps)
	im.SetPreferDist(pd)
	devMode := true
	if localRepo.DevMode != nil {
		devMode = *localRepo.DevMode
	}
	im.PackageEvent = func(event string, op *solver.Operation) error {
		return events.DispatchPackageEvent(event, devMode, op)
	}
	if devMode {
		os.Setenv("COMPOSER_DEV_MODE", "1")
	} else {
		os.Setenv("COMPOSER_DEV_MODE", "0")
	}
	if _, err := events.DispatchScript(composer.PreInstallCmd, devMode, nil, nil); err != nil {
		return scriptExit(err), err
	}
	if err := im.Execute(localRepo, uninstallOps, devMode, true, false); err != nil {
		return 1, err
	}
	if err := im.Execute(localRepo, installOps, devMode, true, false); err != nil {
		return 1, err
	}
	if !in.Bool("no-autoloader") {
		optimize := in.Bool("optimize-autoloader") || c.Config.GetBool("optimize-autoloader")
		authoritative := in.Bool("classmap-authoritative") || c.Config.GetBool("classmap-authoritative")
		apcuPrefix, hasPrefix := in.String("apcu-autoloader-prefix")
		apcu := (hasPrefix && in.present["apcu-autoloader-prefix"]) || in.Bool("apcu-autoloader") || c.Config.GetBool("apcu-autoloader")
		if err := localRepo.Reload(); err != nil {
			return 1, err
		}
		if _, err := composer.RunAutoloadDump(c, events, localRepo, composer.AutoloadOptions{
			Optimize: optimize, ClassMapAuthoritative: authoritative, Apcu: apcu, ApcuPrefix: apcuPrefix,
			RunScripts: false, PlatformFilter: platformFilter(in), IM: im,
		}); err != nil {
			return 1, err
		}
	}
	if _, err := events.DispatchScript(composer.PostInstallCmd, devMode, nil, nil); err != nil {
		return scriptExit(err), err
	}
	return 0, nil
}

// scriptExit is the exit code of a failed script event.
func scriptExit(err error) int {
	if se, ok := err.(*composer.ScriptError); ok {
		return se.Code
	}
	return 1
}
