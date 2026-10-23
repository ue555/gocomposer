package cli

// Port of Composer\Command\RequireCommand and the non-interactive parts of
// PackageDiscoveryTrait (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// Interactive prompts are not supported: gocomposer behaves as if every
// question got its default answer (like --no-interaction).

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

func requireCommand() *Command {
	return &Command{Name: "require", Run: runRequire}
}

type requireCmd struct {
	ctx            *Context
	file           string
	lock           string
	newlyCreated   bool
	firstRequire   bool
	composerBackup []byte
	lockBackup     []byte
	resolved       bool
	discovery
}

// lockFileFor is Factory::getLockFile.
func lockFileFor(composerFile string) string {
	if strings.HasSuffix(composerFile, ".json") {
		return strings.TrimSuffix(composerFile, "json") + "lock"
	}
	return composerFile + ".lock"
}

func runRequire(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	file, err := composer.ComposerFileName()
	if err != nil {
		return 1, err
	}
	rc := &requireCmd{ctx: ctx, file: file}
	rc.discovery.ctx = ctx
	rc.discovery.minStability = rc.minimumStability
	if in.Bool("no-suggest") {
		io.WriteError(`<warning>You are using the deprecated option "--no-suggest". It has no effect and will break in Composer 3.</warning>`)
	}
	if _, err := os.Stat(file); err != nil {
		rc.newlyCreated = true
		if err := os.WriteFile(file, []byte("{\n}\n"), 0o666); err != nil {
			io.WriteError("<error>" + file + " could not be created.</error>")
			return 1, nil
		}
	}
	if f, err := os.Open(file); err != nil {
		io.WriteError("<error>" + file + " is not readable.</error>")
		return 1, nil
	} else {
		f.Close()
	}
	if st, err := os.Stat(file); err == nil && st.Size() == 0 {
		if err := jsonfile.SafeFilePutContents(file, []byte("{\n}\n")); err != nil {
			return 1, err
		}
	}
	rc.lock = lockFileFor(file)
	rc.composerBackup, _ = os.ReadFile(file)
	if b, err := os.ReadFile(rc.lock); err == nil {
		rc.lockBackup = b
	}
	if f, err := os.OpenFile(file, os.O_WRONLY, 0); err != nil {
		io.WriteError("<error>" + file + " is not writable.</error>")
		return 1, nil
	} else {
		f.Close()
	}
	if in.Bool("fixed") {
		cfg, err := rc.readJSON()
		if err != nil {
			return 1, err
		}
		typ := cfg.Str("type")
		if typ == "" {
			typ = "library"
		}
		if typ != "project" && !in.Bool("dev") {
			io.WriteError(`<error>The "--fixed" option is only allowed when the root package has a "project" type or for dev dependencies to prevent possible misuses.</error>`)
			if !cfg.Isset("type") {
				io.WriteError(`<error>If your package is not a library, you can explicitly specify the "type" by using "composer config type project".</error>`)
			}
			return 1, nil
		}
	}
	if rc.newlyCreated {
		// the instance from initialize() predates composer.json
		ctx.cache = nil
	}
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	platformRepo, err := c.CreatePlatformRepo(c.Config.Get("platform"))
	if err != nil {
		return 1, err
	}
	rc.platformRepo = platformRepo
	rc.repos = append([]repo.Repository{platformRepo}, c.Repositories...)
	preferredStability := c.Root().MinimumStability
	if c.Root().PreferStable {
		preferredStability = "stable"
	}
	requirements, err := rc.determineRequirements(c, in.Args("packages"), preferredStability, in.Bool("no-update"), in.Bool("fixed"))
	if err != nil {
		if rc.newlyCreated {
			rc.revertComposerFile()
			return 1, &composer.UserError{Msg: "No composer.json present in the current directory (" + rc.file + "), this may be the cause of the following exception.", Code: 1, Prev: err}
		}
		return 1, err
	}
	names, reqs, err := formatRequirements(requirements)
	if err != nil {
		return 1, err
	}
	if !in.Bool("dev") && io.Interactive {
		if err := rc.recommendDev(names); err != nil {
			return 1, err
		}
	}
	requireKey, removeKey := "require", "require-dev"
	if in.Bool("dev") {
		requireKey, removeKey = "require-dev", "require"
	}
	var toGuess []string
	for _, n := range names {
		if reqs[n] == "guess" {
			reqs[n] = "*"
			toGuess = append(toGuess, n)
		}
	}
	for _, n := range names {
		if strings.ToLower(n) == c.Package.Name {
			io.WriteError(fmt.Sprintf("<error>Root package '%s' cannot require itself in its composer.json</error>", n))
			return 1, nil
		}
		if reqs[n] == "self.version" {
			continue
		}
		if _, err := semver.ParseConstraints(reqs[n]); err != nil {
			return 1, err
		}
	}
	byKey, err := rc.packagesByRequireKey()
	if err != nil {
		return 1, err
	}
	inconsistent := 0
	for _, n := range names {
		if k, ok := byKey[n]; ok && k != requireKey {
			inconsistent++
			with := "without"
			if in.Bool("dev") {
				with = "with"
			}
			io.WriteError(fmt.Sprintf("<warning>%s is currently present in the %s key and you ran the command %s the --dev flag, which will move it to the %s key.</warning>", n, removeKey, with, requireKey))
		}
	}
	if inconsistent > 0 && io.Interactive {
		what := "this requirement"
		if inconsistent > 1 {
			what = "these requirements"
		}
		move, err := io.AskConfirmation("<info>Do you want to move "+what+"?</info> [<comment>no</comment>]? ", false)
		if err != nil {
			return 1, err
		}
		if !move {
			with := "with"
			if in.Bool("dev") {
				with = "without"
			}
			rerun, err := io.AskConfirmation("<info>Do you want to re-run the command "+with+" --dev?</info> [<comment>yes</comment>]? ", true)
			if err != nil {
				return 1, err
			}
			if !rerun {
				return 0, nil
			}
			in.setTrue("dev")
			requireKey, removeKey = removeKey, requireKey
		}
	}
	sortPackages := in.Bool("sort-packages") || c.Config.GetBool("sort-packages")
	rc.firstRequire = rc.newlyCreated
	if !rc.firstRequire {
		def, err := rc.readJSON()
		if err != nil {
			return 1, err
		}
		if phpjson.ArrayLen(get(def, "require")) == 0 && phpjson.ArrayLen(get(def, "require-dev")) == 0 {
			rc.firstRequire = true
		}
	}
	if !in.Bool("dry-run") {
		if err := rc.updateFile(names, reqs, requireKey, removeKey, sortPackages); err != nil {
			return 1, err
		}
	}
	verb := "updated"
	if rc.newlyCreated {
		verb = "created"
	}
	io.WriteError("<info>" + rc.file + " has been " + verb + "</info>")
	if in.Bool("no-update") {
		return 0, nil
	}
	code, err := rc.doUpdate(names, reqs, requireKey, removeKey)
	if err == nil && code == 0 && len(toGuess) > 0 {
		code, err = rc.updateRequirementsAfterResolution(toGuess, requireKey, removeKey, sortPackages, in.Bool("dry-run"), in.Bool("fixed"))
	}
	if err != nil && !rc.resolved {
		rc.revertComposerFile()
	}
	if in.Bool("dry-run") && rc.newlyCreated {
		os.Remove(rc.file)
	}
	return code, err
}

// recommendDev is RequireCommand's check for packages tagged as dev tools
// (interactive mode, without --dev).
func (rc *requireCmd) recommendDev(names []string) error {
	io := rc.ctx.IO
	devTags := []string{"dev", "testing", "static analysis"}
	current, err := rc.packagesByRequireKey()
	if err != nil {
		return err
	}
	var devPackages [][]string
	for _, n := range names {
		if _, ok := current[n]; ok {
			continue
		}
		var found []*pkg.Package
		for _, r := range rc.repos {
			ps, _ := r.FindPackages(n, nil)
			found = append(found, ps...)
		}
		p := mostCurrentVersion(found)
		if p == nil {
			continue
		}
		kw := map[string]bool{}
		for _, k := range p.Keywords() {
			kw[strings.ToLower(k)] = true
		}
		var tags []string
		for _, t := range devTags {
			if kw[t] {
				tags = append(tags, t)
			}
		}
		if len(tags) > 0 {
			devPackages = append(devPackages, tags)
		}
	}
	if len(devPackages) != len(names) {
		return nil
	}
	plural, plural2, plural3 := "", "is", "it is"
	if len(names) > 1 {
		plural, plural2, plural3 = "s", "are", "they are"
	}
	var tags []string
	seen := map[string]bool{}
	for _, l := range devPackages {
		for _, t := range l {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	io.WriteError("<warning>The package" + plural + " you required " + plural2 + " recommended to be placed in require-dev (because " + plural3 + ` tagged as "` + strings.Join(tags, `", "`) + `") but you did not use --dev.</warning>`)
	ok, err := io.AskConfirmation("<info>Do you want to re-run the command with --dev?</> [<comment>yes</>]? ", true)
	if err != nil {
		return err
	}
	if ok {
		rc.ctx.In.setTrue("dev")
	}
	return nil
}

// mostCurrentVersion is PackageSorter::getMostCurrentVersion.
func mostCurrentVersion(ps []*pkg.Package) *pkg.Package {
	if len(ps) == 0 {
		return nil
	}
	highest := ps[0]
	for _, c := range ps {
		if c.IsDefaultBranch() {
			return c
		}
		if semver.VersionCompare(highest.Version(), c.Version()) < 0 {
			highest = c
		}
	}
	return highest
}

func get(m *phpjson.Map, k string) phpjson.Value {
	if m == nil {
		return nil
	}
	v, _ := m.Get(k)
	return v
}

func (rc *requireCmd) readJSON() (*phpjson.Map, error) {
	raw, err := os.ReadFile(rc.file)
	if err != nil {
		return nil, err
	}
	v, err := jsonfile.Parse(raw, rc.file)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	return m, nil
}

// packagesByRequireKey is RequireCommand::getPackagesByRequireKey.
func (rc *requireCmd) packagesByRequireKey() (map[string]string, error) {
	def, err := rc.readJSON()
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	phpjson.Each(get(def, "require"), func(k string, _ phpjson.Value) { out[k] = "require" })
	phpjson.Each(get(def, "require-dev"), func(k string, _ phpjson.Value) { out[k] = "require-dev" })
	return out, nil
}

// minimumStability is PackageDiscoveryTrait::getMinimumStability.
func (rc *requireCmd) minimumStability() string {
	def, err := rc.readJSON()
	if err == nil {
		if ms, ok := def.Get("minimum-stability"); ok && ms != nil {
			s, _ := phpjson.ScalarString(ms)
			if n, err := semver.NormalizeStability(s); err == nil {
				return n
			}
		}
	}
	return "stable"
}

// levenshtein is PHP levenshtein() (byte based, unit costs).
func levenshtein(a, b string) int {
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

// updateFile is RequireCommand::updateFile with the JsonManipulator fast
// path and the full rewrite fallback.
func (rc *requireCmd) updateFile(names []string, reqs map[string]string, requireKey, removeKey string, sortPackages bool) error {
	ok, err := rc.updateFileCleanly(names, reqs, requireKey, removeKey, sortPackages)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	raw, err := os.ReadFile(rc.file)
	if err != nil {
		return err
	}
	def, err := rc.readJSON()
	if err != nil {
		return err
	}
	for _, n := range names {
		rk := def.Map(requireKey)
		if rk == nil {
			rk = phpjson.NewMap()
			def.Set(requireKey, rk)
		}
		rk.Set(n, reqs[n])
		if rm := def.Map(removeKey); rm != nil {
			rm.Delete(n)
			if rm.Len() == 0 {
				def.Delete(removeKey)
			}
		}
	}
	return jsonfile.Write(rc.file, def, jsonfile.DetectIndenting(string(raw)))
}

func (rc *requireCmd) updateFileCleanly(names []string, reqs map[string]string, requireKey, removeKey string, sortPackages bool) (bool, error) {
	raw, err := os.ReadFile(rc.file)
	if err != nil {
		return false, err
	}
	m, err := jsonmanip.New(string(raw))
	if err != nil {
		return false, err
	}
	ok, err := jsonmanip.Try(func() bool {
		for _, n := range names {
			if !m.AddLink(requireKey, n, reqs[n], sortPackages) {
				return false
			}
			if !m.RemoveSubNode(removeKey, n) {
				return false
			}
		}
		m.RemoveMainKeyIfEmpty(removeKey)
		return true
	})
	if err != nil || !ok {
		return false, err
	}
	return true, jsonfile.SafeFilePutContents(rc.file, []byte(m.Contents()))
}

func (rc *requireCmd) revertComposerFile() {
	io := rc.ctx.IO
	if rc.newlyCreated {
		io.WriteErrorV("\n<error>Installation failed, deleting "+rc.file+".</error>", true, cio.Normal)
		os.Remove(rc.file)
		if _, err := os.Stat(rc.lock); err == nil {
			os.Remove(rc.lock)
		}
		return
	}
	msg := " to its "
	if rc.lockBackup != nil {
		msg = " and " + rc.lock + " to their "
	}
	io.WriteError("\n<error>Installation failed, reverting " + rc.file + msg + "original content.</error>")
	_ = jsonfile.SafeFilePutContents(rc.file, rc.composerBackup)
	if rc.lockBackup != nil {
		_ = jsonfile.SafeFilePutContents(rc.lock, rc.lockBackup)
	}
}

// doUpdate is RequireCommand::doUpdate.
func (rc *requireCmd) doUpdate(names []string, reqs map[string]string, requireKey, removeKey string) (int, error) {
	ctx, io, in := rc.ctx, rc.ctx.IO, rc.ctx.In
	ctx.cache = nil
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	rc.resolved = false
	if in.Bool("dry-run") {
		root := c.Root()
		linkType := pkg.TypeRequire
		if requireKey == "require-dev" {
			linkType = pkg.TypeDevRequire
		}
		var newLinks pkg.Links
		for _, n := range names {
			l, err := pkg.CreateLink(root.Name, root.PrettyVersion(), linkType, strings.ToLower(n), reqs[n])
			if err != nil {
				return 1, err
			}
			newLinks = append(newLinks, l)
		}
		merge := func(existing pkg.Links) pkg.Links {
			out := append(pkg.Links(nil), existing...)
			for _, l := range newLinks {
				replaced := false
				for i, e := range out {
					if !e.Indexed && e.Target == l.Target {
						out[i] = l
						replaced = true
					}
				}
				if !replaced {
					out = append(out, l)
				}
			}
			return out
		}
		remove := func(existing pkg.Links) pkg.Links {
			var out pkg.Links
			for _, e := range existing {
				if _, drop := reqs[e.Target]; !drop {
					out = append(out, e)
				}
			}
			return out
		}
		if requireKey == "require" {
			root.SetRequires(merge(root.Requires()))
			root.SetDevRequires(remove(root.DevRequires()))
		} else {
			root.SetDevRequires(merge(root.DevRequires()))
			root.SetRequires(remove(root.Requires()))
		}
		lower := map[string]string{}
		var lnames []string
		for _, n := range names {
			lower[strings.ToLower(n)] = reqs[n]
			lnames = append(lnames, strings.ToLower(n))
		}
		root.References = composer.ExtractReferences(lnames, lower, root.References)
		root.StabilityFlags = composer.ExtractStabilityFlags(lnames, lower, root.MinimumStability, root.StabilityFlags)
	}
	flags := ""
	transitive := solver.UpdateOnlyListed
	if in.Bool("update-with-all-dependencies") || in.Bool("with-all-dependencies") {
		transitive = solver.UpdateListedWithTransitiveDeps
		flags = " --with-all-dependencies"
	} else if in.Bool("update-with-dependencies") || in.Bool("with-dependencies") {
		transitive = solver.UpdateListedWithTransitiveDepsNoRoot
		flags = " --with-dependencies"
	}
	io.WriteError("<info>Running composer update " + strings.Join(names, " ") + flags + "</info>")
	if err := checkPlugins(ctx, c); err != nil {
		return 1, err
	}
	inst := composer.NewInstaller(c)
	if err := configureInstaller(ctx, c, inst); err != nil {
		return 1, err
	}
	inst.DevMode = !in.Bool("update-no-dev")
	inst.Update = true
	inst.Install = !in.Bool("no-install")
	inst.UpdateAllowTransitive = transitive
	inst.PreferStable = in.Bool("prefer-stable")
	inst.PreferLowest = in.Bool("prefer-lowest")
	inst.MinimalUpdate = in.Bool("minimal-changes") || c.Config.GetBool("update-with-minimal-changes")
	inst.OnPreOperationsExec = func() { rc.resolved = true }
	if locked, _ := c.Locker.IsLocked(); !rc.firstRequire && locked {
		inst.UpdateAllowList = names
		inst.UpdateAllowListSet = true
	}
	status, err := inst.Run()
	if err != nil {
		return status, err
	}
	if status != 0 && status != composer.ErrAuditFailed {
		if status == composer.ErrDependencyResolution {
			for _, r := range parseNameVersionPairs(in.Args("packages")) {
				if !r.HasVersion {
					io.WriteError(`You can also try re-running composer require with an explicit version constraint, e.g. "composer require ` + r.Name + `:*" to figure out if any version is installable, or "composer require ` + r.Name + `:^2.1" if you know which you need.`)
					break
				}
			}
		}
		rc.revertComposerFile()
	}
	return status, nil
}

var reFeatureBranch = regexp.MustCompile(`^dev-`)
var reMainBranches = regexp.MustCompile(`^dev-(?:main|master|trunk|latest)$`)

// updateRequirementsAfterResolution is RequireCommand::updateRequirementsAfterResolution.
func (rc *requireCmd) updateRequirementsAfterResolution(toUpdate []string, requireKey, removeKey string, sortPackages, dryRun, fixed bool) (int, error) {
	ctx, io := rc.ctx, rc.ctx.IO
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	var r repo.Repository = c.LocalRepo
	locked, _ := c.Locker.IsLocked()
	if locked {
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
		r = lr
	}
	vs := composer.NewVersionSelector(solver.NewRepositorySet("stable", nil, nil, nil, nil, nil), nil)
	if _, facts, _ := c.PHP(); facts != nil {
		vp := strings.Split(strings.SplitN(facts.PHPVersion, "-", 2)[0], ".")
		if len(vp) >= 3 {
			vs.PHPVersion = strings.Join(vp[:3], ".")
		}
	}
	var names []string
	reqs := map[string]string{}
	for _, name := range toUpdate {
		p := repo.FindPackage(r, name, &semver.MatchAll{})
		for p != nil && p.AliasOf != nil {
			p = p.AliasOf
		}
		if p == nil {
			continue
		}
		v := vs.FindRecommendedRequireVersion(p)
		if fixed {
			v = p.PrettyVersion()
		}
		names = append(names, name)
		reqs[name] = v
		io.WriteError(fmt.Sprintf("Using version <info>%s</info> for <info>%s</info>", v, name))
		if reFeatureBranch.MatchString(v) && !reMainBranches.MatchString(v) {
			io.WriteError("<warning>Version " + v + " looks like it may be a feature branch which is unlikely to keep working in the long run and may be in an unstable state</warning>")
		}
	}
	if !dryRun {
		if err := rc.updateFile(names, reqs, requireKey, removeKey, sortPackages); err != nil {
			return 1, err
		}
		if locked && c.Config.GetBool("lock") {
			lower := map[string]string{}
			var lnames []string
			for _, n := range names {
				lower[n] = reqs[n]
				lnames = append(lnames, n)
			}
			flags := composer.ExtractStabilityFlags(lnames, lower, c.Root().MinimumStability, map[string]int{})
			raw, err := os.ReadFile(rc.file)
			if err != nil {
				return 1, err
			}
			err = c.Locker.UpdateHash(raw, func(data *phpjson.Map) *phpjson.Map {
				sf := data.Map("stability-flags")
				if sf == nil {
					sf = phpjson.NewMap()
				}
				for _, n := range lnames {
					if f, ok := flags[n]; ok {
						sf.Set(n, int64(f))
					}
				}
				data.Set("stability-flags", sf)
				return data
			})
			if err != nil {
				return 1, err
			}
		}
	}
	return 0, nil
}
