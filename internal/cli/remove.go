package cli

// Port of Composer\Command\RemoveCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// In interactive mode the "remove it from require(-dev)?" confirmation is
// asked like Composer (internal/cio QuestionHelper port).

import (
	"os"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/configsource"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

func removeCommand() *Command {
	return &Command{Name: "remove", Run: runRemove}
}

func runRemove(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	args := in.Args("packages")
	if len(args) == 0 && !in.Bool("unused") {
		return 1, &CLIError{Msg: `Not enough arguments (missing: "packages").`}
	}
	var packages []string
	for _, p := range args {
		packages = append(packages, strings.ToLower(p))
	}
	if in.Bool("unused") {
		c, err := ctx.composerFor(true)
		if err != nil {
			return 1, err
		}
		if locked, _ := c.Locker.IsLocked(); !locked {
			return 1, &composer.UserError{Msg: "A valid composer.lock file is required to run this command with --unused", Code: 1}
		}
		lr, err := c.Locker.LockedRepository(false)
		if err != nil {
			return 1, err
		}
		lockedPackages := append([]*pkg.Package(nil), lr.MustPackages()...)
		required := map[string]bool{}
		for _, l := range append(append(pkg.Links(nil), c.Package.Requires()...), c.Package.DevRequires()...) {
			required[l.Target] = true
		}
		for {
			found := false
			var rest []*pkg.Package
			for _, p := range lockedPackages {
				hit := false
				for _, n := range p.Names(true) {
					if required[n] {
						hit = true
						break
					}
				}
				if hit {
					for _, l := range p.Requires() {
						required[l.Target] = true
					}
					found = true
					continue
				}
				rest = append(rest, p)
			}
			lockedPackages = rest
			if !found {
				break
			}
		}
		for _, p := range lockedPackages {
			packages = append(packages, p.Name)
		}
		if len(packages) == 0 {
			io.WriteError("<info>No unused packages to remove</info>")
			return 0, nil
		}
	}
	file, err := composer.ComposerFileName()
	if err != nil {
		return 1, err
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return 1, &composer.UserError{Msg: "Could not read " + file + "\n\n" + err.Error(), Code: 1}
	}
	value, err := jsonfile.Parse(raw, file)
	if err != nil {
		return 1, err
	}
	def := phpjson.ToMap(value)
	if def == nil {
		def = phpjson.NewMap()
	}
	backup := raw
	source := configsource.New(file, false)
	typ, altType := "require", "require-dev"
	if in.Bool("dev") {
		typ, altType = "require-dev", "require"
	}
	if in.Bool("update-with-dependencies") {
		io.WriteError(`<warning>You are using the deprecated option "update-with-dependencies". This is now default behaviour. The --no-update-with-dependencies option can be used to remove a package without its dependencies.</warning>`)
	}
	// lower-case name => name as written, per link type
	names := map[string]map[string]string{}
	keys := map[string][]string{}
	for _, lt := range []string{"require", "require-dev"} {
		if !def.Isset(lt) {
			continue
		}
		names[lt] = map[string]string{}
		phpjson.Each(get(def, lt), func(n string, _ phpjson.Value) {
			names[lt][strings.ToLower(n)] = n
			keys[lt] = append(keys[lt], n)
		})
		// $composer[$linkType][strtolower($name)] = $name adds lower-cased keys
		for _, n := range append([]string(nil), keys[lt]...) {
			if l := strings.ToLower(n); l != n && !inList(l, keys[lt]) {
				keys[lt] = append(keys[lt], l)
			}
		}
	}
	dryRun := in.Bool("dry-run")
	toRemove := map[string][]string{}
	remove := func(lt, name string) error {
		if dryRun {
			toRemove[lt] = append(toRemove[lt], name)
			return nil
		}
		return source.RemoveLink(lt, name)
	}
	askAlt := func(name string) error {
		if !io.Interactive {
			return nil
		}
		ok, err := io.AskConfirmation("Do you want to remove it from "+altType+" [<comment>yes</comment>]? ", true)
		if err != nil || !ok {
			return err
		}
		return remove(altType, name)
	}
	for _, p := range packages {
		switch {
		case names[typ] != nil && names[typ][p] != "":
			if err := remove(typ, names[typ][p]); err != nil {
				return 1, err
			}
		case names[altType] != nil && names[altType][p] != "":
			io.WriteError("<warning>" + names[altType][p] + " could not be found in " + typ + " but it is present in " + altType + "</warning>")
			if err := askAlt(names[altType][p]); err != nil {
				return 1, err
			}
		case names[typ] != nil && len(grepNames(p, keys[typ])) > 0:
			for _, m := range grepNames(p, keys[typ]) {
				if err := remove(typ, m); err != nil {
					return 1, err
				}
			}
		case names[altType] != nil && len(grepNames(p, keys[altType])) > 0:
			for _, m := range grepNames(p, keys[altType]) {
				io.WriteError("<warning>" + m + " could not be found in " + typ + " but it is present in " + altType + "</warning>")
				if err := askAlt(m); err != nil {
					return 1, err
				}
			}
		default:
			io.WriteError("<warning>" + p + " is not required in your composer.json and has not been removed</warning>")
		}
	}
	io.WriteError("<info>" + file + " has been updated</info>")
	if in.Bool("no-update") {
		return 0, nil
	}
	ctx.cache = nil
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	if dryRun {
		root := c.Root()
		drop := func(links pkg.Links, list []string) pkg.Links {
			var out pkg.Links
			for _, l := range links {
				if !inList(l.Target, list) {
					out = append(out, l)
				}
			}
			return out
		}
		root.SetRequires(drop(root.Requires(), toRemove["require"]))
		root.SetDevRequires(drop(root.DevRequires(), toRemove["require-dev"]))
	}
	if ap := phpjson.ToMap(c.Config.Get("allow-plugins")); ap != nil && !dryRun {
		var removed []string
		for _, k := range ap.Keys() {
			if inList(k, packages) {
				removed = append(removed, k)
			}
		}
		if len(removed) > 0 {
			if len(removed) == ap.Len() {
				if err := source.RemoveConfigSetting("allow-plugins"); err != nil {
					return 1, err
				}
			} else {
				for _, p := range removed {
					if err := source.RemoveConfigSetting("allow-plugins." + p); err != nil {
						return 1, err
					}
				}
			}
		}
	}
	if err := checkPlugins(ctx, c); err != nil {
		return 1, err
	}
	inst := composer.NewInstaller(c)
	if err := configureInstaller(ctx, c, inst); err != nil {
		return 1, err
	}
	transitive := solver.UpdateListedWithTransitiveDepsNoRoot
	flags := ""
	if in.Bool("update-with-all-dependencies") || in.Bool("with-all-dependencies") {
		transitive = solver.UpdateListedWithTransitiveDeps
		flags = " --with-all-dependencies"
	} else if in.Bool("no-update-with-dependencies") {
		transitive = solver.UpdateOnlyListed
		flags = " --with-dependencies"
	}
	io.WriteError("<info>Running composer update " + strings.Join(packages, " ") + flags + "</info>")
	inst.DevMode = !in.Bool("update-no-dev")
	inst.Update = true
	inst.Install = !in.Bool("no-install")
	inst.UpdateAllowTransitive = transitive
	inst.DryRun = dryRun
	inst.MinimalUpdate = in.Bool("minimal-changes") || c.Config.GetBool("update-with-minimal-changes")
	if locked, _ := c.Locker.IsLocked(); locked {
		inst.UpdateAllowList = packages
		inst.UpdateAllowListSet = true
	}
	status, err := inst.Run()
	if err != nil {
		return status, err
	}
	if status != 0 {
		io.WriteError("\n<error>Removal failed, reverting " + file + " to its original content.</error>")
		if err := jsonfile.SafeFilePutContents(file, backup); err != nil {
			return 1, err
		}
	}
	if !dryRun {
		for _, p := range packages {
			if found, _ := c.LocalRepo.FindPackages(p, nil); len(found) > 0 {
				io.WriteError("<error>Removal failed, " + p + " is still present, it may be required by another package. See `composer why " + p + "`.</error>")
				return 2, nil
			}
		}
	}
	return status, nil
}

// grepNames is Preg::grep(BasePackage::packageNameToRegexp($pattern), $names).
func grepNames(pattern string, names []string) []string {
	re := solver.PackageNameToRegexp(pattern)
	var out []string
	for _, n := range names {
		if re.MatchString(n) {
			out = append(out, n)
		}
	}
	return out
}
