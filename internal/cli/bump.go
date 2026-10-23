package cli

// Port of Composer\Command\BumpCommand and Composer\Package\Version\VersionBumper
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

const (
	bumpErrorGeneric     = 1
	bumpErrorLockOutdate = 2
)

func init() { registerCommand(bumpCommand) }

func bumpCommand() *Command {
	return &Command{Name: "bump", Run: func(ctx *Context) (int, error) {
		in := ctx.In
		return doBump(ctx, in.Bool("dev-only"), in.Bool("no-dev-only"), in.Bool("dry-run"), in.Args("packages"), "--dev-only")
	}}
}

// bumpAfterUpdate is the bump-after-update tail of UpdateCommand::execute:
// the packages installed or updated by the lock transaction are bumped.
// mode is the option/config value ("dev", "no-dev", or anything else for all).
func bumpAfterUpdate(ctx *Context, inst *composer.Installer, mode string) (int, error) {
	var updated []string
	if lt := inst.LockTransaction(); lt != nil {
		for _, op := range lt.Operations() {
			switch op.Type {
			case solver.OpInstall:
				updated = append(updated, op.Package.Name)
			case solver.OpUpdate:
				updated = append(updated, op.Target.Name)
			}
		}
	}
	if len(updated) == 0 {
		return 0, nil
	}
	ctx.IO.WriteError("<info>Bumping dependencies</info>")
	return doBump(ctx, mode == "dev", mode == "no-dev", ctx.In.Bool("dry-run"), updated, "--bump-after-update=dev")
}

var reBumpFilterStrip = regexp.MustCompile(`[:= ].+`)

// doBump is BumpCommand::doBump.
func doBump(ctx *Context, devOnly, noDevOnly, dryRun bool, packagesFilter []string, devOnlyFlagHint string) (int, error) {
	io := ctx.IO
	path, err := composer.ComposerFileName()
	if err != nil {
		return 1, err
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		io.WriteError("<error>" + path + " is not readable.</error>")
		return bumpErrorGeneric, nil
	}
	if f, err := os.OpenFile(path, os.O_WRONLY, 0); err != nil {
		if werr := os.WriteFile(path, contents, 0o666); werr != nil {
			io.WriteError("<error>" + path + " is not writable.</error>")
			return bumpErrorGeneric, nil
		}
	} else {
		f.Close()
	}
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	// $hasLockfileDisabled is (despite its name) true when locking is enabled
	lockEnabled := c.Config.GetBool("lock")
	var r repo.Repository
	if !lockEnabled {
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
		r = lr
	} else if locked, _ := c.Locker.IsLocked(); locked {
		fresh, err := bumpLockIsFresh(c)
		if err != nil {
			return 1, err
		}
		if !fresh {
			io.WriteError("<error>The lock file is not up to date with the latest changes in composer.json. Run the appropriate `update` to fix that before you use the `bump` command.</error>")
			return bumpErrorLockOutdate, nil
		}
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
		r = lr
	} else {
		r = c.LocalRepo
	}
	root := c.Package
	if root.Type() != "project" && !devOnly {
		io.WriteError("<warning>Warning: Bumping dependency constraints is not recommended for libraries as it will narrow down your dependencies and may cause problems for your users.</warning>")
		raw, err := os.ReadFile(path)
		if err != nil {
			return 1, err
		}
		v, err := jsonfile.Parse(raw, path)
		if err != nil {
			return 1, err
		}
		if m := phpjson.ToMap(v); m == nil || !m.Isset("type") {
			io.WriteError(`<warning>If your package is not a library, you can explicitly specify the "type" by using "composer config type project".</warning>`)
			io.WriteError("<warning>Alternatively you can use " + devOnlyFlagHint + ` to only bump dependencies within "require-dev".</warning>`)
		}
	}
	type task struct {
		key   string
		links pkg.Links
	}
	var tasks []task
	if !devOnly {
		tasks = append(tasks, task{"require", root.Requires()})
	}
	if !noDevOnly {
		tasks = append(tasks, task{"require-dev", root.DevRequires()})
	}
	if len(packagesFilter) > 0 {
		seen := map[string]bool{}
		var names []string
		for _, f := range packagesFilter {
			n := strings.ToLower(reBumpFilterStrip.ReplaceAllString(f, ""))
			if !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
		pattern := solver.PackageNamesToRegexp(names)
		for i := range tasks {
			var keep pkg.Links
			for _, l := range tasks[i].links {
				if pattern.MatchString(l.Target) {
					keep = append(keep, l)
				}
			}
			tasks[i].links = keep
		}
	}
	type update struct{ key, name, version string }
	var updates []update
	for _, t := range tasks {
		for _, l := range t.links {
			if repo.IsPlatformPackage(l.Target) {
				continue
			}
			p := repo.FindPackage(r, l.Target, &semver.MatchAll{})
			if p == nil {
				continue
			}
			for p.AliasOf != nil {
				p = p.AliasOf
			}
			bumped := bumpRequirement(l.Constraint, l.PrettyConstraint, p)
			if bumped == l.PrettyConstraint {
				continue
			}
			updates = append(updates, update{t.key, l.Target, bumped})
		}
	}
	if !dryRun {
		ok, err := bumpUpdateFileCleanly(path, func(m *jsonmanip.Manipulator) bool {
			for _, u := range updates {
				if !m.AddLink(u.key, u.name, u.version, false) {
					return false
				}
			}
			return true
		})
		if err != nil {
			return 1, err
		}
		if !ok {
			raw, err := os.ReadFile(path)
			if err != nil {
				return 1, err
			}
			v, err := jsonfile.Parse(raw, path)
			if err != nil {
				return 1, err
			}
			def := phpjson.ToMap(v)
			if def == nil {
				def = phpjson.NewMap()
			}
			for _, u := range updates {
				sub := def.Map(u.key)
				if sub == nil {
					sub = phpjson.NewMap()
					def.Set(u.key, sub)
				}
				sub.Set(u.name, u.version)
			}
			if err := jsonfile.Write(path, def, jsonfile.DetectIndenting(string(raw))); err != nil {
				return 1, err
			}
		}
	}
	changeCount := len(updates)
	if changeCount > 0 {
		if dryRun {
			io.Write("<info>" + path + " would be updated with:</info>")
			// $updates is grouped by require type
			for _, key := range []string{"require", "require-dev"} {
				for _, u := range updates {
					if u.key == key {
						io.Write(fmt.Sprintf("<info> - %s.%s: %s</info>", u.key, u.name, u.version))
					}
				}
			}
		} else {
			io.Write(fmt.Sprintf("<info>%s has been updated (%d changes).</info>", path, changeCount))
		}
	} else {
		io.Write("<info>No requirements to update in " + path + ".</info>")
	}
	if locked, _ := c.Locker.IsLocked(); !dryRun && locked && lockEnabled && changeCount > 0 {
		raw, err := os.ReadFile(path)
		if err != nil {
			return 1, err
		}
		if err := c.Locker.UpdateHash(raw, nil); err != nil {
			return 1, err
		}
	}
	if dryRun && changeCount > 0 {
		return bumpErrorGeneric, nil
	}
	return 0, nil
}

// bumpLockIsFresh is Locker::isFresh, which reads the lock file from disk
// (not the virtual lock data a dry-run update leaves behind).
func bumpLockIsFresh(c *composer.Composer) (bool, error) {
	file := c.Locker.File
	raw, err := os.ReadFile(file)
	if err != nil {
		return false, &composer.UserError{Msg: "Could not read " + file + "\n\nThe file \"" + file + "\" is not readable.", Code: 1}
	}
	v, err := jsonfile.Parse(raw, file)
	if err != nil {
		return false, err
	}
	if m := phpjson.ToMap(v); m != nil {
		if ch, ok := m.Get("content-hash"); ok && phpjson.Truthy(ch) {
			return ch == c.Locker.ContentHashValue(), nil
		}
	}
	return c.Locker.IsFresh(), nil
}

func bumpUpdateFileCleanly(path string, apply func(*jsonmanip.Manipulator) bool) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, &composer.UserError{Msg: "Unable to read " + path + " contents.", Code: 1}
	}
	m, err := jsonmanip.New(string(raw))
	if err != nil {
		return false, err
	}
	ok, err := jsonmanip.Try(func() bool { return apply(m) })
	if err != nil || !ok {
		return false, err
	}
	if err := jsonfile.SafeFilePutContents(path, []byte(m.Contents())); err != nil {
		return false, &composer.UserError{Msg: "Unable to write new " + path + " contents.", Code: 1}
	}
	return true, nil
}

var (
	reBumpMajor         = regexp.MustCompile(`^([1-9][0-9]*|0\.\d+).*`)
	reBumpSuffix        = regexp.MustCompile(`(?:\.(?:0|9999999))+(-dev)?$`)
	reBumpSimpleStable  = regexp.MustCompile(`^\^\d+(\.\d+)*$`)
	reBumpGreaterEquals = regexp.MustCompile(`^>=v?\d+(?:\.\d+)*$`)
	reBumpWildcard      = regexp.MustCompile(`^\*$`)
)

// bumpRequirement is VersionBumper::bumpRequirement.
func bumpRequirement(constraint semver.Constraint, prettyConstraint string, p *pkg.Package) string {
	if strings.HasPrefix(prettyConstraint, "dev-") {
		return prettyConstraint
	}
	version := p.Version()
	if strings.HasPrefix(version, "dev-") {
		extra, err := pkg.GetBranchAlias(pkg.Dump(p))
		if err != nil || extra == "" || extra == semver.DefaultBranchAlias {
			return prettyConstraint
		}
		version = extra
	}
	if len(semver.BranchNames(constraint)) > 0 {
		return prettyConstraint
	}
	major := version
	if m := reBumpMajor.FindStringSubmatch(version); m != nil {
		major = m[1]
	}
	versionWithoutSuffix := reBumpSuffix.ReplaceAllString(version, "")
	newPrettyConstraint := "^" + versionWithoutSuffix
	if !reBumpSimpleStable.MatchString(newPrettyConstraint) {
		return prettyConstraint
	}
	// the pattern's alternatives, tried in order at each position (the
	// major is inserted unquoted like in the PHP pattern)
	alts := []*regexp.Regexp{
		regexp.MustCompile(`^\^v?` + major + `(?:\.\d+)*$`),
		regexp.MustCompile(`^~v?` + major + `(?:\.\d+){1,3}$`),
		regexp.MustCompile(`^v?` + major + `(?:\.[*x])+$`),
		reBumpGreaterEquals,
		reBumpWildcard,
	}
	type match struct {
		text   string
		offset int
	}
	var matches []match
	s := prettyConstraint
	for pos := 0; pos <= len(s); {
		// (?<=,|\ |\||^) leading separator
		if pos > 0 && s[pos-1] != ',' && s[pos-1] != ' ' && s[pos-1] != '|' {
			pos++
			continue
		}
		end := -1
		for _, alt := range alts {
			// backtracking yields the longest candidate first
			for e := len(s); e > pos; e-- {
				if !alt.MatchString(s[pos:e]) {
					continue
				}
				// (?=,|$|\ |\||@) trailing separator
				if e == len(s) || s[e] == ',' || s[e] == ' ' || s[e] == '|' || s[e] == '@' {
					end = e
					break
				}
			}
			if end >= 0 {
				break
			}
		}
		if end < 0 {
			pos++
			continue
		}
		matches = append(matches, match{s[pos:end], pos})
		pos = end
	}
	if len(matches) == 0 {
		return prettyConstraint
	}
	modified := prettyConstraint
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		if strings.HasPrefix(m.text, ">=") {
			partStart := strings.LastIndexByte(prettyConstraint[:m.offset], '|')
			if partStart < 0 {
				partStart = 0
			} else {
				partStart++
			}
			partEnd := strings.IndexByte(prettyConstraint[m.offset:], '|')
			if partEnd < 0 {
				partEnd = len(prettyConstraint)
			} else {
				partEnd += m.offset
			}
			part := strings.TrimSpace(prettyConstraint[partStart:partEnd])
			pc, err := semver.ParseConstraints(part)
			if err != nil || !pc.Matches(semver.NewConstraint("==", version)) {
				continue
			}
		}
		suffix := ""
		if strings.Count(m.text, ".") == 2 && strings.Count(versionWithoutSuffix, ".") == 1 {
			suffix = ".0"
		}
		var replacement string
		if strings.HasPrefix(m.text, "~") && strings.Count(m.text, ".") != 1 {
			bits := strings.Split(versionWithoutSuffix, ".")
			n := strings.Count(m.text, ".") + 1
			for len(bits) < n {
				bits = append(bits, "0")
			}
			replacement = "~" + strings.Join(bits[:n], ".")
		} else if m.text == "*" || strings.HasPrefix(m.text, ">=") {
			replacement = ">=" + versionWithoutSuffix + suffix
		} else {
			replacement = newPrettyConstraint + suffix
		}
		modified = modified[:m.offset] + replacement + modified[m.offset+len(m.text):]
	}
	newConstraint, err := semver.ParseConstraints(modified)
	if err != nil {
		return prettyConstraint
	}
	if semver.IsSubsetOf(newConstraint, constraint) && semver.IsSubsetOf(constraint, newConstraint) {
		return prettyConstraint
	}
	return modified
}
