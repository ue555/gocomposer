package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

var installOptions = []OptDef{
	{Name: "prefer-source", Kind: OptNone, Desc: "Forces installation from package sources when possible, including VCS information."},
	{Name: "prefer-dist", Kind: OptNone, Desc: "Forces installation from package dist (default behavior)."},
	{Name: "prefer-install", Kind: OptRequired, Desc: "Forces installation from package dist|source|auto (auto chooses source for dev versions, dist for the rest)."},
	{Name: "dry-run", Kind: OptNone, Desc: "Outputs the operations but will not execute anything (implicitly enables --verbose)."},
	{Name: "download-only", Kind: OptNone, Desc: "Download only, do not install packages."},
	{Name: "dev", Kind: OptNone, Desc: "DEPRECATED: Enables installation of require-dev packages (enabled by default, only present for BC)."},
	{Name: "no-suggest", Kind: OptNone, Desc: "DEPRECATED: This flag does not exist anymore."},
	{Name: "no-dev", Kind: OptNone, Desc: "Disables installation of require-dev packages."},
	{Name: "no-security-blocking", Kind: OptNone, Desc: "DEPRECATED: use --no-blocking instead."},
	{Name: "no-blocking", Kind: OptNone, Desc: "Disables all policy blocking during this command."},
	{Name: "no-autoloader", Kind: OptNone, Desc: "Skips autoloader generation"},
	{Name: "no-progress", Kind: OptNone, Desc: "Do not output download progress."},
	{Name: "no-install", Kind: OptNone, Desc: "Do not use, only defined here to catch misuse of the install command."},
	{Name: "audit", Kind: OptNone, Desc: "Run an audit after installation is complete."},
	{Name: "audit-format", Kind: OptRequired, Default: "summary", Desc: `Audit output format. Must be "table", "plain", "json", or "summary".`},
	{Name: "optimize-autoloader", Shortcut: "o", Kind: OptNone, Desc: "Optimize autoloader during autoloader dump"},
	{Name: "classmap-authoritative", Shortcut: "a", Kind: OptNone, Desc: "Autoload classes from the classmap only. Implicitly enables `--optimize-autoloader`."},
	{Name: "strict-psr-autoloader", Kind: OptNone, Desc: "Return a failed status code (6) if PSR-4 or PSR-0 mapping errors are present."},
	{Name: "apcu-autoloader", Kind: OptNone, Desc: "Use APCu to cache found/not-found classes."},
	{Name: "apcu-autoloader-prefix", Kind: OptRequired, Desc: "Use a custom prefix for the APCu autoloader cache. Implicitly enables --apcu-autoloader"},
	{Name: "ignore-platform-req", Kind: OptArray, Desc: "Ignore a specific platform requirement (php & ext- packages)."},
	{Name: "ignore-platform-reqs", Kind: OptNone, Desc: "Ignore all platform requirements (php & ext- packages)."},
}

func preferredInstall(c *composer.Composer, in *Input) (bool, bool, error) {
	preferSource, preferDist := false, false
	switch c.Config.Get("preferred-install") {
	case "source":
		preferSource = true
	case "dist":
		preferDist = true
	}
	if pi, ok := in.String("prefer-install"); ok && pi != "" {
		if in.Bool("prefer-source") {
			return false, false, &composer.UserError{Msg: "--prefer-source can not be used together with --prefer-install", Code: 1}
		}
		if in.Bool("prefer-dist") {
			return false, false, &composer.UserError{Msg: "--prefer-dist can not be used together with --prefer-install", Code: 1}
		}
		switch pi {
		case "dist":
			return false, true, nil
		case "source":
			return true, false, nil
		case "auto":
			return false, false, nil
		default:
			return false, false, &composer.UserError{Msg: "--prefer-install accepts one of \"dist\", \"source\" or \"auto\", got " + pi, Code: 1}
		}
	}
	if in.Bool("prefer-source") || in.Bool("prefer-dist") {
		return in.Bool("prefer-source"), in.Bool("prefer-dist"), nil
	}
	return preferSource, preferDist, nil
}

func pluginCheck(ctx *Context, c *composer.Composer) error {
	return nil
}

func configureInstaller(ctx *Context, c *composer.Composer, inst *composer.Installer) error {
	in := ctx.In
	ps, pd, err := preferredInstall(c, in)
	if err != nil {
		return err
	}
	optimize := in.Bool("optimize-autoloader") || c.Config.GetBool("optimize-autoloader")
	authoritative := in.Bool("classmap-authoritative") || c.Config.GetBool("classmap-authoritative")
	strictPsr := in.Bool("strict-psr-autoloader") || c.Config.GetBool("strict-psr-autoloader")
	apcuPrefix, hasPrefix := in.String("apcu-autoloader-prefix")
	apcu := (hasPrefix && apcuPrefix != "") || in.Bool("apcu-autoloader") || c.Config.GetBool("apcu-autoloader")
	if strictPsr && !optimize && !authoritative {
		return &composer.UserError{Msg: "--strict-psr-autoloader mode only works with optimized autoloader, use --optimize-autoloader or --classmap-authoritative if you want a strict return value.", Code: 1}
	}
	inst.DryRun = in.Bool("dry-run")
	inst.IM.OutputProgress = !in.Bool("no-progress")
	pc, err := composer.CreatePolicyConfig(c, in.Bool("no-blocking") || in.Bool("no-security-blocking"), in.Bool("no-cooldown"))
	if err != nil {
		return err
	}
	inst.PolicyConfig = pc
	inst.DownloadOnly = in.Bool("download-only")
	inst.PreferSource = ps
	inst.PreferDist = pd
	inst.DevMode = !in.Bool("no-dev")
	inst.DumpAutoloader = !in.Bool("no-autoloader")
	inst.OptimizeAutoloader = optimize
	inst.ClassMapAuthoritative = authoritative
	inst.StrictPsr = strictPsr
	inst.ApcuAutoloader = apcu
	inst.ApcuPrefix = apcuPrefix
	inst.PlatformFilter = platformFilter(in)
	// BaseCommand::createAuditConfig: commands defining --audit only audit
	// when it is given; the others audit unless --no-audit is given.
	if ctx.Cmd.HasOption("audit") {
		inst.Audit = in.Bool("audit")
	} else {
		inst.Audit = !in.Bool("no-audit")
	}
	inst.ErrorOnAudit = in.Bool("audit")
	if f, ok := in.String("audit-format"); ok && f != "" {
		inst.AuditFormat = f
	}
	// --no-scripts only disables the root package scripts (Factory:
	// $dispatcher->setRunScripts(!$disableScripts), see Composer.DisableScripts);
	// the installer still dispatches its events to plugin listeners
	return nil
}

func installCommand() *Command {
	return &Command{
		Name:        "install",
		Aliases:     []string{"i"},
		Description: "Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json",
		Options:     installOptions,
		Args:        []ArgDef{{Name: "packages", Array: true, Desc: "Should not be provided, use composer require instead to add a given package to composer.json."}},
		Run: func(ctx *Context) (int, error) {
			io := ctx.IO
			in := ctx.In
			if in.Bool("dev") {
				io.WriteError(`<warning>You are using the deprecated option "--dev". It has no effect and will break in Composer 3.</warning>`)
			}
			if in.Bool("no-suggest") {
				io.WriteError(`<warning>You are using the deprecated option "--no-suggest". It has no effect and will break in Composer 3.</warning>`)
			}
			if args := in.Args("packages"); len(args) > 0 {
				io.WriteError(`<error>Invalid argument ` + strings.Join(args, " ") + `. Use "composer require ` + strings.Join(args, " ") + `" instead to add packages to your composer.json.</error>`)
				return 1, nil
			}
			if in.Bool("no-install") {
				io.WriteError(`<error>Invalid option "--no-install". Use "composer update --no-install" instead if you are trying to update the composer.lock file.</error>`)
				return 1, nil
			}
			c, err := ctx.composerFor(true)
			if err != nil {
				return 1, err
			}
			if err := checkPlugins(ctx, c); err != nil {
				return 1, err
			}
			inst := composer.NewInstaller(c)
			if err := configureInstaller(ctx, c, inst); err != nil {
				return 1, err
			}
			return inst.Run()
		},
	}
}

var updateOptions = append(append([]OptDef(nil), installOptions...),
	OptDef{Name: "with", Kind: OptArray, Desc: "Temporary version constraint to add, e.g. foo/bar:1.0.0 or foo/bar=1.0.0"},
	OptDef{Name: "lock", Kind: OptNone, Desc: "Overwrites the lock file hash to suppress warning about the lock file being out of date without updating package versions."},
	OptDef{Name: "no-audit", Kind: OptNone, Desc: "Skip the audit step after updating the composer.lock file."},
	OptDef{Name: "with-dependencies", Shortcut: "w", Kind: OptNone, Desc: "Update also dependencies of packages in the argument list, except those which are root requirements."},
	OptDef{Name: "with-all-dependencies", Shortcut: "W", Kind: OptNone, Desc: "Update also dependencies of packages in the argument list, including those which are root requirements."},
	OptDef{Name: "minimal-changes", Shortcut: "m", Kind: OptNone, Desc: "During a partial update with -w/-W, only perform absolutely necessary changes to transitive dependencies."},
	OptDef{Name: "patch-only", Kind: OptNone, Desc: "Only allow patch version updates for currently installed dependencies."},
	OptDef{Name: "interactive", Shortcut: "i", Kind: OptNone, Desc: "Interactive interface with autocompletion to select the packages to update."},
	OptDef{Name: "root-reqs", Kind: OptNone, Desc: "Restricts the update to your first degree dependencies."},
	OptDef{Name: "prefer-stable", Kind: OptNone, Desc: "Prefer stable versions of dependencies."},
	OptDef{Name: "prefer-lowest", Kind: OptNone, Desc: "Prefer lowest versions of dependencies."},
	OptDef{Name: "bump-after-update", Kind: OptOptional, Desc: "Runs bump after performing the update."},
)

var reWithShorthand = regexp.MustCompile(`\S+[ =:]\S+`)

// formatRequirements is BaseCommand::formatRequirements (name[:= ]constraint pairs).
func formatRequirements(reqs []string) ([]string, map[string]string, error) {
	var names []string
	out := map[string]string{}
	for _, p := range parseNameVersionPairs(reqs) {
		if !p.HasVersion {
			return nil, nil, &composer.UserError{Msg: "Option " + p.Name + " is missing a version constraint, use e.g. " + p.Name + ":^1.0", Code: 1}
		}
		if _, ok := out[p.Name]; !ok {
			names = append(names, p.Name)
		}
		out[p.Name] = p.Version
	}
	return names, out, nil
}

var rePairSplit = regexp.MustCompile(`^([^=: ]+)[=: ](.*)$`)

// parseNameVersionPairs is VersionParser::parseNameVersionPairs.
// nameVersion is one parseNameVersionPairs entry (Version is set when
// HasVersion, possibly to "").
type nameVersion struct {
	Name, Version string
	HasVersion    bool
}

func parseNameVersionPairs(pairs []string) []nameVersion {
	var out []nameVersion
	for i := 0; i < len(pairs); i++ {
		pair := rePairSplit.ReplaceAllString(strings.TrimSpace(pairs[i]), "$1 $2")
		if !strings.Contains(pair, " ") && i+1 < len(pairs) && !strings.Contains(pairs[i+1], "/") &&
			!regexp.MustCompile(`(?i)(?:[a-z0-9_/-])\*|\*(?:[a-z0-9_/-])`).MatchString(pairs[i+1]) && !isPlatform(pairs[i+1]) {
			pair += " " + pairs[i+1]
			i++
		}
		if idx := strings.Index(pair, " "); idx > 0 {
			out = append(out, nameVersion{pair[:idx], pair[idx+1:], true})
		} else {
			out = append(out, nameVersion{Name: pair})
		}
	}
	return out
}

func isPlatform(s string) bool {
	return regexp.MustCompile(`(?i)^(?:php(?:-64bit|-ipv6|-zts|-debug)?|hhvm|(?:ext|lib)-[a-z0-9](?:[_.-]?[a-z0-9]+)*|composer(?:-(?:plugin|runtime)-api)?)$`).MatchString(s)
}

func updateCommand() *Command {
	return &Command{
		Name:        "update",
		Aliases:     []string{"u", "upgrade"},
		Description: "Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file",
		Options:     updateOptions,
		Args:        []ArgDef{{Name: "packages", Array: true, Desc: "Packages that should be updated, if not provided all packages are."}},
		Run: func(ctx *Context) (int, error) {
			io := ctx.IO
			in := ctx.In
			if in.Bool("dev") {
				io.WriteError(`<warning>You are using the deprecated option "--dev". It has no effect and will break in Composer 3.</warning>`)
			}
			if in.Bool("no-suggest") {
				io.WriteError(`<warning>You are using the deprecated option "--no-suggest". It has no effect and will break in Composer 3.</warning>`)
			}
			c, err := ctx.composerFor(true)
			if err != nil {
				return 1, err
			}
			if err := checkPlugins(ctx, c); err != nil {
				return 1, err
			}
			packages := in.Args("packages")
			reqNames, reqs, err := formatRequirements(in.Strings("with"))
			if err != nil {
				return 1, err
			}
			if len(packages) > 0 {
				var withShort []string
				for _, p := range packages {
					if reWithShorthand.MatchString(p) {
						withShort = append(withShort, p)
					}
				}
				n2, r2, err := formatRequirements(withShort)
				if err != nil {
					return 1, err
				}
				for _, n := range n2 {
					if _, ok := reqs[n]; !ok {
						reqNames = append(reqNames, n)
					}
					reqs[n] = r2[n]
				}
				for i, p := range packages {
					if reWithShorthand.MatchString(p) {
						packages[i] = regexp.MustCompile(`^([^ =:]+)[ =:].*$`).ReplaceAllString(p, "$1")
					}
				}
			}
			root := c.Root()
			links := map[string]string{}
			for _, n := range reqNames {
				links[n] = reqs[n]
			}
			root.StabilityFlags = composer.ExtractStabilityFlags(reqNames, links, root.MinimumStability, root.StabilityFlags)
			temporary := map[string]semver.Constraint{}
			rootReqs := append(append(pkg.Links(nil), root.Requires()...), root.DevRequires()...)
			for _, name := range reqNames {
				constraint := reqs[name]
				lname := strings.ToLower(name)
				parsed, err := semver.ParseConstraints(constraint)
				if err != nil {
					return 1, err
				}
				if strings.Contains(lname, "*") {
					re := solver.PackageNameToRegexp(lname)
					for _, rr := range rootReqs {
						if !re.MatchString(rr.Target) {
							continue
						}
						temporary[rr.Target] = parsed
						if !semver.HaveIntersections(parsed, rr.Constraint) {
							io.WriteError(`<error>The temporary constraint "` + constraint + `" for "` + lname + `" matching "` + rr.Target + `" must be a subset of the constraint in your composer.json (` + rr.PrettyConstraint + `)</error>`)
							return 1, nil
						}
					}
				} else {
					temporary[lname] = parsed
					if rr := rootReqs.Get(lname); rr != nil && !semver.HaveIntersections(parsed, rr.Constraint) {
						io.WriteError(`<error>The temporary constraint "` + constraint + `" for "` + lname + `" must be a subset of the constraint in your composer.json (` + rr.PrettyConstraint + `)</error>`)
						io.Write("<info>Run `composer require " + lname + "` or `composer require " + lname + ":" + constraint + "` instead to replace the constraint</info>")
						return 1, nil
					}
				}
			}
			if in.Bool("patch-only") {
				locked, _ := c.Locker.IsLocked()
				if !locked {
					return 1, &composer.UserError{Msg: "patch-only can only be used with a lock file present", Code: 1}
				}
				lr, err := c.Locker.LockedRepository(true)
				if err != nil {
					return 1, err
				}
				for _, p := range lr.MustPackages() {
					if p.AliasOf != nil || p.IsDev() {
						continue
					}
					m := regexp.MustCompile(`^(\d+\.\d+\.\d+)`).FindStringSubmatch(p.Version())
					if m == nil {
						continue
					}
					cst := semver.MustParseConstraints("~" + m[1])
					if t, ok := temporary[p.Name]; ok {
						temporary[p.Name] = semver.CreateMulti([]semver.Constraint{t, cst}, true)
					} else {
						temporary[p.Name] = cst
					}
				}
			}
			if in.Bool("interactive") {
				// UpdateCommand::getPackagesInteractively (update_interactive.go)
				if packages, err = updatePackagesInteractively(ctx, c, packages); err != nil {
					return 1, err
				}
			}
			if in.Bool("root-reqs") {
				var requires []string
				for _, l := range root.Requires() {
					requires = append(requires, l.Target)
				}
				if !in.Bool("no-dev") {
					for _, l := range root.DevRequires() {
						requires = append(requires, l.Target)
					}
				}
				if len(packages) > 0 {
					var keep []string
					for _, p := range packages {
						for _, r := range requires {
							if p == r {
								keep = append(keep, p)
								break
							}
						}
					}
					packages = keep
				} else {
					packages = requires
				}
			}
			var filtered []string
			for _, p := range packages {
				if p != "lock" && p != "nothing" && p != "mirrors" {
					filtered = append(filtered, p)
				}
			}
			updateMirrors := in.Bool("lock") || len(filtered) != len(packages)
			packages = filtered
			if updateMirrors && len(packages) > 0 {
				io.WriteError("<error>You cannot simultaneously update only a selection of packages and regenerate the lock file metadata.</error>")
				return 255, nil
			}
			inst := composer.NewInstaller(c)
			if err := configureInstaller(ctx, c, inst); err != nil {
				return 1, err
			}
			inst.MinimalUpdate = in.Bool("minimal-changes") || c.Config.GetBool("update-with-minimal-changes")
			inst.Update = true
			inst.Install = !in.Bool("no-install")
			inst.UpdateMirrors = updateMirrors
			if len(packages) > 0 {
				inst.UpdateAllowList = packages
				inst.UpdateAllowListSet = true
			}
			transitive := solver.UpdateOnlyListed
			if in.Bool("with-all-dependencies") {
				transitive = solver.UpdateListedWithTransitiveDeps
			} else if in.Bool("with-dependencies") {
				transitive = solver.UpdateListedWithTransitiveDepsNoRoot
			}
			inst.UpdateAllowTransitive = transitive
			inst.PreferStable = in.Bool("prefer-stable")
			inst.PreferLowest = in.Bool("prefer-lowest")
			inst.TemporaryConstraints = temporary
			code, err := inst.Run()
			if err != nil || code != 0 {
				return code, err
			}
			bump, hasBump := in.String("bump-after-update")
			if !hasBump {
				if v, ok := c.Config.Get("bump-after-update").(string); ok {
					bump, hasBump = v, true
				} else if c.Config.Get("bump-after-update") == true {
					hasBump = true
				}
			}
			if hasBump && !in.Bool("lock") {
				return bumpAfterUpdate(ctx, inst, bump)
			}
			return 0, nil
		},
	}
}

func dumpAutoloadCommand() *Command {
	return &Command{
		Name:        "dump-autoload",
		Aliases:     []string{"dumpautoload"},
		Description: "Dumps the autoloader",
		Options: []OptDef{
			{Name: "optimize", Shortcut: "o", Kind: OptNone, Desc: "Optimizes PSR0 and PSR4 packages to be loaded with classmaps too, good for production."},
			{Name: "classmap-authoritative", Shortcut: "a", Kind: OptNone, Desc: "Autoload classes from the classmap only. Implicitly enables `--optimize`."},
			{Name: "apcu", Kind: OptNone, Desc: "Use APCu to cache found/not-found classes."},
			{Name: "apcu-prefix", Kind: OptRequired, Desc: "Use a custom prefix for the APCu autoloader cache. Implicitly enables --apcu"},
			{Name: "dry-run", Kind: OptNone, Desc: "Outputs the operations but will not execute anything."},
			{Name: "dev", Kind: OptNone, Desc: "Enables autoload-dev rules."},
			{Name: "no-dev", Kind: OptNone, Desc: "Disables autoload-dev rules."},
			{Name: "ignore-platform-req", Kind: OptArray, Desc: "Ignore a specific platform requirement (php & ext- packages)."},
			{Name: "ignore-platform-reqs", Kind: OptNone, Desc: "Ignore all platform requirements (php & ext- packages)."},
			{Name: "strict-psr", Kind: OptNone, Desc: "Return a failed status code (1) if PSR-4 or PSR-0 mapping errors are present. Requires --optimize to work."},
			{Name: "strict-ambiguous", Kind: OptNone, Desc: "Return a failed status code (2) if the same class is found in multiple files. Requires --optimize to work."},
		},
		Run: func(ctx *Context) (int, error) {
			io := ctx.IO
			in := ctx.In
			c, err := ctx.composerFor(true)
			if err != nil {
				return 1, err
			}
			im := composer.NewInstallationManager(c)
			missing := false
			for _, p := range c.LocalRepo.CanonicalPackages() {
				ip := im.InstallPath(p)
				if ip != "" {
					if _, err := os.Stat(ip); err != nil {
						missing = true
						io.Write(`<warning>Not all dependencies are installed. Make sure to run a "composer install" to install missing dependencies</warning>`)
						break
					}
				}
			}
			optimize := in.Bool("optimize") || c.Config.GetBool("optimize-autoloader")
			authoritative := in.Bool("classmap-authoritative") || c.Config.GetBool("classmap-authoritative")
			strictPsr := in.Bool("strict-psr") || c.Config.GetBool("strict-psr-autoloader")
			apcuPrefix, hasPrefix := in.String("apcu-prefix")
			apcu := (hasPrefix && apcuPrefix != "") || in.Bool("apcu") || c.Config.GetBool("apcu-autoloader")
			if strictPsr && !optimize && !authoritative {
				return 1, &composer.UserError{Msg: "--strict-psr mode only works with optimized autoloader, use --optimize or --classmap-authoritative if you want a strict return value.", Code: 1}
			}
			if in.Bool("strict-ambiguous") && !optimize && !authoritative {
				return 1, &composer.UserError{Msg: "--strict-ambiguous mode only works with optimized autoloader, use --optimize or --classmap-authoritative if you want a strict return value.", Code: 1}
			}
			switch {
			case authoritative:
				io.Write("<info>Generating optimized autoload files (authoritative)</info>")
			case optimize:
				io.Write("<info>Generating optimized autoload files</info>")
			default:
				io.Write("<info>Generating autoload files</info>")
			}
			var devMode *bool
			if in.Bool("no-dev") {
				f := false
				devMode = &f
			}
			if in.Bool("dev") {
				if in.Bool("no-dev") {
					return 1, &composer.UserError{Msg: "You can not use both --no-dev and --dev as they conflict with each other.", Code: 1}
				}
				t := true
				devMode = &t
			}
			events := composer.NewEventDispatcher(c)
			events.RunScripts = !c.DisableScripts
			res, err := composer.RunAutoloadDump(c, events, c.LocalRepo, composer.AutoloadOptions{
				DevMode: devMode, Optimize: optimize, ClassMapAuthoritative: authoritative, Apcu: apcu, ApcuPrefix: apcuPrefix,
				RunScripts: !c.DisableScripts, PlatformFilter: platformFilter(in), IM: im, DryRun: in.Bool("dry-run"), StrictAmbiguous: in.Bool("strict-ambiguous"),
			})
			if err != nil {
				return 1, err
			}
			n := len(res.ClassMap.Map)
			switch {
			case authoritative:
				io.Write(fmt.Sprintf("<info>Generated optimized autoload files (authoritative) containing %d classes</info>", n))
			case optimize:
				io.Write(fmt.Sprintf("<info>Generated optimized autoload files containing %d classes</info>", n))
			default:
				io.Write("<info>Generated autoload files</info>")
			}
			if missing || (strictPsr && len(res.PsrViolations) > 0) {
				return 1, nil
			}
			if in.Bool("strict-ambiguous") {
				order, _ := res.ClassMap.AmbiguousClasses(false)
				if len(order) > 0 || len(res.ClassMap.AmbiguousFolders(false)) > 0 {
					return 2, nil
				}
			}
			return 0, nil
		},
	}
}
