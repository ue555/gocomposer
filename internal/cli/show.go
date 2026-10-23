package cli

// Port of Composer\Command\ShowCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// The --latest/--outdated listing and lookup live in show_latest.go.

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
	"gocomposer/internal/spdx"
)

func showCommand() *Command {
	return &Command{Name: "show", Run: runShow}
}

var treeColors = []string{"green", "yellow", "cyan", "magenta", "blue"}

type showCtx struct {
	ctx       *Context
	c         *composer.Composer
	installed *repo.Installed
	// repos is the repository list getPackage() builds its pool from
	// ($repos, a CompositeRepository flattened one level by RepositorySet).
	repos []repo.Repository
	// latestRS is the cached ShowCommand::getRepositorySet().
	latestRS *solver.RepositorySet
}

// showEntry is one $packages[$type][$name] value: a package, or only a name
// for packages listed through ComposerRepository::getPackageNames().
type showEntry struct {
	name string
	p    *pkg.Package
}

type showTyped struct {
	order []string
	m     map[string]showEntry
}

func (t *showTyped) set(name string, e showEntry) {
	if _, ok := t.m[name]; !ok {
		t.order = append(t.order, name)
	}
	t.m[name] = e
}

// errShowExit1 marks an error already reported by requireComposer.
var errShowExit1 = errors.New("exit 1")

func runShow(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	if in.Bool("tree") {
		for _, c := range treeColors {
			io.SetStyle(c, c)
		}
	}
	c, _ := ctx.composerFor(true)
	if c == nil {
		ctx.cache = nil
	}
	if in.Bool("installed") && !in.Bool("self") {
		io.WriteError(`<warning>You are using the deprecated option "installed". Only installed packages are shown by default now. The --all option can be used to show all packages.</warning>`)
	}
	if in.Bool("outdated") {
		in.setTrue("latest")
	} else if len(in.Strings("ignore")) > 0 {
		io.WriteError(`<warning>You are using the option "ignore" for action other than "outdated", it will be ignored.</warning>`)
	}
	if in.Bool("direct") && (in.Bool("all") || in.Bool("available") || in.Bool("platform")) {
		io.WriteError("The --direct (-D) option is not usable in combination with --all, --platform (-p) or --available (-a)")
		return 1, nil
	}
	if in.Bool("tree") && (in.Bool("all") || in.Bool("available")) {
		io.WriteError("The --tree (-t) option is not usable in combination with --all or --available (-a)")
		return 1, nil
	}
	n := 0
	for _, o := range []string{"patch-only", "minor-only", "major-only"} {
		if in.Bool(o) {
			n++
		}
	}
	if n > 1 {
		io.WriteError("Only one of --major-only, --minor-only or --patch-only can be used at once")
		return 1, nil
	}
	if in.Bool("tree") && in.Bool("latest") {
		io.WriteError("The --tree (-t) option is not usable in combination with --latest (-l)")
		return 1, nil
	}
	if in.Bool("tree") && in.Bool("path") {
		io.WriteError("The --tree (-t) option is not usable in combination with --path (-P)")
		return 1, nil
	}
	format, _ := in.String("format")
	if format == "" {
		format = "text"
	}
	if format != "text" && format != "json" {
		io.WriteError(fmt.Sprintf(`Unsupported format "%s". See help for supported formats.`, format))
		return 1, nil
	}
	platformReqFilter := platformFilter(in)

	var platformOverrides phpjson.Value
	if c != nil {
		platformOverrides = c.Config.Get("platform")
	}
	var platformRepo repo.Repository
	var platformTyped *platform.Repository
	{
		// new PlatformRepository([], $platformOverrides) also works without a composer.json
		pc := c
		if pc == nil {
			pc = &composer.Composer{}
		}
		pr, err := pc.CreatePlatformRepo(platformOverrides)
		if err != nil {
			return 1, err
		}
		platformRepo, platformTyped = pr, pr
	}
	sc := &showCtx{ctx: ctx, c: c}
	var lockedRepo *repo.Array
	var selfPackage *pkg.Package
	requireComposer := func() (*composer.Composer, error) {
		if c != nil {
			return c, nil
		}
		cc, err := ctx.composerFor(true)
		if err != nil {
			// Application::getComposer(true): the InvalidArgumentException of a
			// missing composer.json is printed as a plain message, then exit(1)
			if strings.HasPrefix(err.Error(), "Composer could not find ") {
				io.WriteError(err.Error())
				return nil, errShowExit1
			}
			return nil, err
		}
		c = cc
		sc.c = cc
		return cc, nil
	}
	// listRepos is RepositoryUtils::flattenRepositories($repos).
	var listRepos []repo.Repository

	switch {
	case in.Bool("self") && !in.Bool("installed") && !in.Bool("locked"):
		cc, err := requireComposer()
		if err == errShowExit1 {
			return 1, nil
		}
		if err != nil {
			return 1, err
		}
		selfPackage = cc.Package
		if in.Bool("name-only") {
			io.Write(selfPackage.Name)
			return 0, nil
		}
		if in.Arg("package") != "" {
			return 1, &composer.UserError{Msg: "You cannot use --self together with a package name", Code: 1}
		}
		sc.installed = repo.NewInstalled(repo.NewRootPackageRepository(selfPackage))
		sc.repos = sc.installed.Repositories()
	case in.Bool("platform"):
		if platformRepo == nil {
			cc, err := requireComposer()
			if err == errShowExit1 {
				return 1, nil
			}
			if err != nil {
				return 1, err
			}
			pr, err := cc.CreatePlatformRepo(nil)
			if err != nil {
				return 1, err
			}
			platformRepo, platformTyped = pr, pr
		}
		sc.installed = repo.NewInstalled(platformRepo)
		sc.repos = sc.installed.Repositories()
	case in.Bool("available"):
		if c == nil {
			// RepositoryFactory::defaultReposWithDefaultManager (packagist.org)
			return 1, &composer.UserError{Msg: "show --available without a composer.json (listing the default packagist.org repository) is not supported by gocomposer yet.", Code: 1}
		}
		sc.installed = repo.NewInstalled(platformRepo, c.LocalRepo)
		sc.repos = append([]repo.Repository(nil), c.Repositories...)
	case in.Bool("all") && c != nil:
		locked, _ := c.Locker.IsLocked()
		if locked {
			lr, err := c.Locker.LockedRepository(true)
			if err != nil {
				return 1, err
			}
			lockedRepo = lr
			sc.installed = repo.NewInstalled(lockedRepo, c.LocalRepo, platformRepo)
		} else {
			sc.installed = repo.NewInstalled(c.LocalRepo, platformRepo)
		}
		// new FilterRepository($installedRepo, ['canonical' => false])
		opts := phpjson.NewMap()
		opts.Set("canonical", false)
		for _, r := range sc.installed.Repositories() {
			f, err := repo.NewFilter(r, opts)
			if err != nil {
				return 1, err
			}
			sc.repos = append(sc.repos, f)
		}
		sc.repos = append(sc.repos, c.Repositories...)
	case in.Bool("all"):
		return 1, &composer.UserError{Msg: "show --all without a composer.json (listing the default packagist.org repository) is not supported by gocomposer yet.", Code: 1}
	case in.Bool("locked"):
		locked := false
		if c != nil {
			locked, _ = c.Locker.IsLocked()
		}
		if !locked {
			return 1, &composer.UserError{Msg: "A valid composer.json and composer.lock files is required to run this command with --locked", Code: 1}
		}
		lr, err := c.Locker.LockedRepository(!in.Bool("no-dev"))
		if err != nil {
			return 1, err
		}
		lockedRepo = lr
		if in.Bool("self") {
			lockedRepo.AddPackage(c.Package.Clone())
		}
		sc.installed = repo.NewInstalled(lockedRepo)
		sc.repos = sc.installed.Repositories()
	default:
		cc, err := requireComposer()
		if err == errShowExit1 {
			return 1, nil
		}
		if err != nil {
			return 1, err
		}
		root := cc.Package
		var rootRepo repo.Repository = repo.NewInstalledArray()
		if in.Bool("self") {
			rootRepo = repo.NewRootPackageRepository(root.Clone())
		}
		if in.Bool("no-dev") {
			packages := filterRequiredPackages(cc.LocalRepo.MustPackages(), root, false, nil)
			arr := repo.NewInstalledArray()
			for _, p := range packages {
				arr.AddPackage(p.Clone())
			}
			sc.installed = repo.NewInstalled(rootRepo, arr)
		} else {
			sc.installed = repo.NewInstalled(rootRepo, cc.LocalRepo)
		}
		sc.repos = sc.installed.Repositories()
		if len(sc.installed.Packages()) == 0 {
			hasNonPlatform := func(links pkg.Links) bool {
				for _, l := range links {
					if !repo.IsPlatformPackage(l.Target) {
						return true
					}
				}
				return false
			}
			if hasNonPlatform(root.Requires()) || hasNonPlatform(root.DevRequires()) {
				io.WriteError("<warning>No dependencies installed. Try running composer install or update.</warning>")
			}
		}
	}
	for _, r := range sc.repos {
		if f, ok := r.(*repo.Filter); ok {
			r = f.Repo
		}
		listRepos = append(listRepos, r)
	}

	if in.Bool("latest") && c == nil {
		io.WriteError(`No composer.json found in the current directory, disabling "latest" option`)
		in.present["latest"] = false
	}
	var cooldown *policy.CooldownPolicyConfig
	if in.Bool("latest") && c != nil {
		cd, err := sc.cooldownPolicyConfig()
		if err != nil {
			return 1, err
		}
		cooldown = cd
	}

	packageFilter := in.Arg("package")
	var single *pkg.Package
	var versions *phpjson.Map
	if selfPackage != nil {
		single = selfPackage
		versions = phpjson.NewMap()
		versions.Set(single.PrettyVersion(), single.Version())
	} else if packageFilter != "" && !strings.Contains(packageFilter, "*") {
		p, vs, err := sc.getPackage(packageFilter, in.Arg("version"))
		if err != nil {
			return 1, err
		}
		single, versions = p, vs
		if single != nil && in.Bool("direct") && !inList(single.Name, sc.rootRequires()) {
			return 1, &composer.UserError{Msg: `Package "` + single.Name + `" is installed but not a direct dependent of the root package.`, Code: 1}
		}
		if single == nil {
			hint := ""
			if in.Bool("locked") {
				hint += " in lock file"
			}
			if wd, ok := in.String("working-dir"); ok {
				hint += " in " + wd + "/composer.json"
			}
			if repo.IsPlatformPackage(packageFilter) && !in.Bool("platform") {
				hint += ", try using --platform (-p) to show platform packages"
			}
			if !in.Bool("all") && !in.Bool("available") {
				hint += ", try using --available (-a) to show all available packages"
			}
			return 1, &composer.UserError{Msg: `Package "` + packageFilter + `" not found` + hint + ".", Code: 1}
		}
	}

	if single != nil {
		exitCode := 0
		if in.Bool("tree") {
			tree := sc.generatePackageTree(single)
			if format == "json" {
				out := phpjson.NewMap()
				out.Set("installed", []phpjson.Value{tree})
				io.Write(phpjson.MustEncode(out, phpjson.ComposerFlags))
			} else {
				sc.displayPackageTree([]*phpjson.Map{tree})
			}
			return exitCode, nil
		}
		var latest *pkg.Package
		var withheld *showWithheld
		if in.Bool("latest") {
			l, w, err := sc.findLatestPackage(single, platformRepo, in.Bool("major-only"), in.Bool("minor-only"), in.Bool("patch-only"), platformReqFilter, cooldown)
			if err != nil {
				return 1, err
			}
			latest, withheld = l, w
		}
		if in.Bool("outdated") && in.Bool("strict") && latest != nil && latest.FullPrettyVersion(true) != single.FullPrettyVersion(true) && !latest.IsAbandoned() {
			exitCode = 1
		}
		if in.Bool("path") {
			io.WriteV(single.Name, false, 32)
			path := sc.installPath(single)
			if path != "" {
				io.Write(" " + firstLine(realpath(path)))
			} else {
				io.Write(" null")
			}
			return exitCode, nil
		}
		if format == "json" {
			sc.printPackageInfoAsJSON(single, versions, latest, withheld)
		} else {
			sc.printPackageInfo(single, versions, latest, withheld)
		}
		return exitCode, nil
	}

	if in.Bool("tree") {
		rootRequires := sc.rootRequires()
		packages := append([]*pkg.Package(nil), sc.installed.Packages()...)
		sort.SliceStable(packages, func(i, j int) bool { return packages[i].String() < packages[j].String() })
		var trees []*phpjson.Map
		for _, p := range packages {
			if inList(p.Name, rootRequires) {
				trees = append(trees, sc.generatePackageTree(p))
			}
		}
		if format == "json" {
			out := phpjson.NewMap()
			list := []phpjson.Value{}
			for _, t := range trees {
				list = append(list, t)
			}
			out.Set("installed", list)
			io.Write(phpjson.MustEncode(out, phpjson.ComposerFlags))
		} else {
			sc.displayPackageTree(trees)
		}
		return 0, nil
	}

	// list packages
	var filterRe *regexp.Regexp
	if packageFilter != "" {
		filterRe = regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(packageFilter), `\*`, `.*?`) + `$`)
	}
	var listFilter []string
	if in.Bool("direct") {
		listFilter = sc.rootRequires()
	}
	if in.Bool("path") && c == nil {
		io.WriteError(`No composer.json found in the current directory, disabling "path" option`)
		in.present["path"] = false
	}
	packages := map[string]*showTyped{}
	typedFor := func(typ string) *showTyped {
		t := packages[typ]
		if t == nil {
			t = &showTyped{m: map[string]showEntry{}}
			packages[typ] = t
		}
		return t
	}
	installedRepos := sc.installed.Repositories()
	for _, r := range listRepos {
		typ := "available"
		switch {
		case platformRepo != nil && r == platformRepo:
			typ = "platform"
		case lockedRepo != nil && r == repo.Repository(lockedRepo):
			typ = "locked"
		default:
			for _, ir := range installedRepos {
				if ir == r {
					typ = "installed"
					break
				}
			}
		}
		if cr, ok := r.(*repo.Composer); ok {
			names, err := cr.PackageNames(packageFilter)
			if err != nil {
				return 1, err
			}
			for _, name := range names {
				typedFor(typ).set(name, showEntry{name: name})
			}
			continue
		}
		ps, err := r.Packages()
		if err != nil {
			return 1, err
		}
		for _, p := range ps {
			t := packages[typ]
			var cur showEntry
			ok := false
			if t != nil {
				cur, ok = t.m[p.Name]
			}
			if !ok || cur.p == nil || semver.VersionCompare(cur.p.Version(), p.Version()) < 0 {
				for p.AliasOf != nil {
					p = p.AliasOf
				}
				if filterRe == nil || filterRe.MatchString(p.Name) {
					if listFilter == nil || inList(p.Name, listFilter) {
						typedFor(typ).set(p.Name, showEntry{name: p.Name, p: p})
					}
				}
			}
		}
		if platformTyped != nil && r == platformRepo {
			disabled := platformTyped.Disabled()
			names := make([]string, 0, len(disabled))
			for name := range disabled {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				typedFor(typ).set(name, showEntry{name: name, p: disabled[name]})
			}
		}
	}
	return sc.listPackages(packages, format, platformRepo, platformReqFilter, cooldown)
}

func nullable(s string) phpjson.Value {
	if s == "" {
		return nil
	}
	return s
}

// viewSourceURL is PackageInfo::getViewSourceUrl.
func viewSourceURL(p *pkg.Package) string {
	if s := p.Support(); s != nil {
		if v, ok := s.Get("source"); ok {
			if str, ok := v.(string); ok && str != "" {
				return str
			}
		}
	}
	return p.SourceURL()
}

// baseTerminalWidth is BaseCommand::getTerminalWidth (non-Windows).
func baseTerminalWidth() int {
	w := terminalWidth()
	if w < 80 {
		w = 80
	}
	return w
}

func firstLine(s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		return s[:i]
	}
	return s
}

// realpath is PHP realpath(); "" stands for false (missing path).
func realpath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	r, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return ""
	}
	return r
}

func (sc *showCtx) installPath(p *pkg.Package) string {
	if sc.c == nil {
		return ""
	}
	return composer.NewInstallationManager(sc.c).InstallPath(p)
}

// rootRequires is ShowCommand::getRootRequires.
func (sc *showCtx) rootRequires() []string {
	if sc.c == nil {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, l := range append(append(pkg.Links(nil), sc.c.Package.Requires()...), sc.c.Package.DevRequires()...) {
		n := strings.ToLower(l.Target)
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

// filterRequiredPackages is RepositoryUtils::filterRequiredPackages.
func filterRequiredPackages(packages []*pkg.Package, requirer *pkg.Package, includeDev bool, bucket []*pkg.Package) []*pkg.Package {
	requires := requirer.Requires()
	if includeDev {
		requires = append(append(pkg.Links(nil), requires...), requirer.DevRequires()...)
	}
	for _, cand := range packages {
		for _, name := range cand.Names(true) {
			if requires.Has(name) {
				found := false
				for _, b := range bucket {
					if b == cand {
						found = true
						break
					}
				}
				if !found {
					bucket = append(bucket, cand)
					bucket = filterRequiredPackages(packages, cand, false, bucket)
				}
				break
			}
		}
	}
	return bucket
}

// getPackage is ShowCommand::getPackage.
func (sc *showCtx) getPackage(name, version string) (*pkg.Package, *phpjson.Map, error) {
	name = strings.ToLower(name)
	var constraint semver.Constraint
	if version != "" {
		c, err := semver.ParseConstraints(version)
		if err != nil {
			return nil, nil, err
		}
		constraint = c
	}
	return sc.getPackageConstraint(name, constraint, version != "")
}

func (sc *showCtx) getPackageConstraint(name string, constraint semver.Constraint, versionGiven bool) (*pkg.Package, *phpjson.Map, error) {
	policy := solver.NewDefaultPolicy(false, false, nil)
	rs := solver.NewRepositorySet("dev", nil, nil, nil, nil, nil)
	rs.AllowInstalledRepositories()
	for _, r := range sc.repos {
		rs.AddRepository(r)
	}
	var pool *solver.Pool
	var err error
	if repo.IsPlatformPackage(name) {
		pool, err = rs.CreatePoolWithAllPackages()
	} else {
		req := solver.NewRequest(nil)
		req.RequireName(name, nil)
		req.RestrictPackages([]string{name})
		pool, err = rs.CreatePool(req, solver.NewPoolBuilder())
	}
	if err != nil {
		return nil, nil, err
	}
	var matched *pkg.Package
	versions := phpjson.NewMap()
	var literals []int
	for _, p := range pool.WhatProvides(name, constraint) {
		if p.AliasOf != nil && p.Version() == semver.DefaultBranchAlias {
			p = p.AliasOf
		}
		if !versionGiven && sc.installed.HasPackage(p) {
			matched = p
		}
		versions.Set(p.PrettyVersion(), p.Version())
		literals = append(literals, p.ID)
	}
	if matched == nil && len(literals) > 0 {
		preferred := policy.SelectPreferredPackages(pool, literals, "")
		matched = pool.LiteralToPackage(preferred[0])
	}
	return matched, versions, nil
}

func (sc *showCtx) printPackageInfo(p *pkg.Package, versions *phpjson.Map, latest *pkg.Package, withheld *showWithheld) {
	io := sc.ctx.IO
	sc.printMeta(p, versions, latest, withheld)
	sc.printLinks(p, p.Requires(), "requires")
	sc.printLinks(p, p.DevRequires(), "requires (dev)")
	if s := p.Suggests(); s != nil && s.Len() > 0 {
		io.Write("\n<info>suggests</info>")
		for _, k := range s.Keys() {
			v, _ := s.Get(k)
			str, _ := phpjson.ScalarString(v)
			io.Write(k + " <comment>" + str + "</comment>")
		}
	}
	sc.printLinks(p, p.Provides(), "provides")
	sc.printLinks(p, p.Conflicts(), "conflicts")
	sc.printLinks(p, p.Replaces(), "replaces")
}

func (sc *showCtx) printMeta(p *pkg.Package, versions *phpjson.Map, latest *pkg.Package, withheld *showWithheld) {
	io := sc.ctx.IO
	isInstalled := !repo.IsPlatformPackage(p.Name) && sc.installed.HasPackage(p)
	io.Write("<info>name</info>     : " + p.PrettyName)
	io.Write("<info>descrip.</info> : " + p.Description())
	io.Write("<info>keywords</info> : " + strings.Join(p.Keywords(), ", "))
	sc.printVersions(p, versions)
	if isInstalled && p.ReleaseDate() != "" {
		if t, ok := pkg.ParsePHPDate(p.ReleaseDate()); ok {
			io.Write("<info>released</info> : " + t.Format("2006-01-02") + ", " + relativeTime(t))
		}
	}
	hasInstallable := latest != nil && latest.FullPrettyVersion(true) != p.FullPrettyVersion(true)
	if withheld != nil && !hasInstallable {
		io.Write("<info>latest</info>   : " + withheld.p.PrettyVersion() + " (still in the cooldown period, " + withheld.availableIn + " left)")
		if latest == nil {
			latest = p
		}
	} else if latest != nil {
		style := updateStatusToVersionStyle(updateStatus(latest, p))
		released := ""
		if rt, ok := releaseTime(latest); ok {
			released = " released " + rt.Format("2006-01-02") + ", " + relativeTime(rt)
		}
		io.Write("<info>latest</info>   : <" + style + ">" + latest.PrettyVersion() + "</" + style + ">" + released)
		if withheld != nil {
			io.Write("<info>cooldown</info> : " + withheld.p.PrettyVersion() + " (still in the cooldown period, " + withheld.availableIn + " left)")
		}
	} else {
		latest = p
	}
	io.Write("<info>type</info>     : " + p.Type())
	sc.printLicenses(p)
	io.Write("<info>homepage</info> : " + p.Homepage())
	io.Write("<info>source</info>   : " + fmt.Sprintf("[%s] <comment>%s</comment> %s", p.SourceType(), p.SourceURL(), p.SourceReference()))
	io.Write("<info>dist</info>     : " + fmt.Sprintf("[%s] <comment>%s</comment> %s", p.DistType(), p.DistURL(), p.DistReference()))
	if isInstalled {
		path := sc.installPath(p)
		if path != "" {
			io.Write("<info>path</info>     : " + realpath(path))
		} else {
			io.Write("<info>path</info>     : null")
		}
	}
	io.Write("<info>names</info>    : " + strings.Join(p.Names(true), ", "))
	if latest.IsAbandoned() {
		repl := ""
		if r := latest.ReplacementPackage(); r != "" {
			repl = " The author suggests using the " + r + " package instead."
		}
		io.WriteError("<warning>Attention: This package is abandoned and no longer maintained." + repl + "</warning>")
	}
	if s := p.Support(); s != nil && s.Len() > 0 {
		io.Write("\n<info>support</info>")
		for _, k := range s.Keys() {
			v, _ := s.Get(k)
			str, _ := phpjson.ScalarString(v)
			io.Write("<comment>" + k + "</comment> : " + str)
		}
	}
	if al := p.Autoload(); al != nil && al.Len() > 0 {
		io.Write("\n<info>autoload</info>")
		for _, typ := range al.Keys() {
			io.Write("<comment>" + typ + "</comment>")
			v, _ := al.Get(typ)
			switch typ {
			case "psr-0", "psr-4":
				phpjson.Each(v, func(ns string, path phpjson.Value) {
					name := ns
					if !phpjson.Truthy(ns) {
						name = "*"
					}
					var pathStr string
					if phpjson.IsArray(path) {
						var parts []string
						phpjson.Each(path, func(_ string, e phpjson.Value) { s, _ := phpjson.ScalarString(e); parts = append(parts, s) })
						pathStr = strings.Join(parts, ", ")
					} else if phpjson.Truthy(path) {
						pathStr, _ = phpjson.ScalarString(path)
					} else {
						pathStr = "."
					}
					io.Write(name + " => " + pathStr)
				})
			case "classmap":
				var parts []string
				phpjson.Each(v, func(_ string, e phpjson.Value) { s, _ := phpjson.ScalarString(e); parts = append(parts, s) })
				io.Write(strings.Join(parts, ", "))
			}
		}
		if inc := p.IncludePaths(); len(inc) > 0 {
			io.Write("<comment>include-path</comment>")
			io.Write(strings.Join(inc, ", "))
		}
	}
}

// semverRsort is composer/semver Semver::rsort.
func semverRsort(versions []string) []string {
	type nv struct {
		norm string
		idx  int
	}
	var list []nv
	for i, v := range versions {
		n, err := semver.Normalize(v)
		if err != nil {
			n = v
		}
		list = append(list, nv{semver.NormalizeDefaultBranch(n), i})
	}
	sort.SliceStable(list, func(i, j int) bool {
		a, b := list[i].norm, list[j].norm
		if a == b {
			return false
		}
		// direction -1: greater first
		return !semver.NewConstraint("<", b).MatchSpecific(semver.NewConstraint("==", a), true)
	})
	out := make([]string, len(list))
	for i, e := range list {
		out[i] = versions[e.idx]
	}
	return out
}

func (sc *showCtx) printVersions(p *pkg.Package, versions *phpjson.Map) {
	vs := semverRsort(versions.Keys())
	for _, ip := range sc.installed.FindPackages(p.Name, nil) {
		iv := ip.PrettyVersion()
		for i, v := range vs {
			if v == iv {
				vs[i] = "<info>* " + iv + "</info>"
				break
			}
		}
	}
	sc.ctx.IO.Write("<info>versions</info> : " + strings.Join(vs, ", "))
}

func (sc *showCtx) printLinks(p *pkg.Package, links pkg.Links, title string) {
	if len(links) == 0 {
		return
	}
	io := sc.ctx.IO
	io.Write("\n<info>" + title + "</info>")
	for _, l := range links {
		io.Write(l.Target + " <comment>" + l.PrettyConstraint + "</comment>")
	}
}

func (sc *showCtx) printLicenses(p *pkg.Package) {
	for _, id := range p.License() {
		out := id
		if l := spdx.LicenseByIdentifier(id); l != nil {
			if l.OsiApproved {
				out = fmt.Sprintf("%s (%s) (OSI approved) %s", l.Name, id, l.URL)
			} else {
				out = fmt.Sprintf("%s (%s) %s", l.Name, id, l.URL)
			}
		}
		sc.ctx.IO.Write("<info>license</info>  : " + out)
	}
}

func strList(ss []string) []phpjson.Value {
	out := []phpjson.Value{}
	for _, s := range ss {
		out = append(out, s)
	}
	return out
}

func (sc *showCtx) printPackageInfoAsJSON(p *pkg.Package, versions *phpjson.Map, latest *pkg.Package, withheld *showWithheld) {
	j := phpjson.NewMap()
	j.Set("name", p.PrettyName)
	j.Set("description", nullable(p.Description()))
	j.Set("keywords", strList(p.Keywords()))
	j.Set("type", p.Type())
	j.Set("homepage", nullable(p.Homepage()))
	j.Set("names", strList(p.Names(true)))
	// appendVersions: uasort(version_compare) then reverse
	keys := versions.Keys()
	sort.SliceStable(keys, func(a, b int) bool {
		va, _ := versions.Get(keys[a])
		vb, _ := versions.Get(keys[b])
		sa, _ := va.(string)
		sb, _ := vb.(string)
		return semver.VersionCompare(sa, sb) < 0
	})
	rev := []phpjson.Value{}
	for i := len(keys) - 1; i >= 0; i-- {
		rev = append(rev, keys[i])
	}
	j.Set("versions", rev)
	if lic := p.License(); len(lic) > 0 {
		var out []phpjson.Value
		for _, id := range lic {
			if l := spdx.LicenseByIdentifier(id); l != nil {
				m := phpjson.NewMap()
				m.Set("name", l.Name)
				m.Set("osi", id)
				m.Set("url", l.URL)
				out = append(out, m)
			} else {
				out = append(out, id)
			}
		}
		j.Set("licenses", out)
	}
	hasInstallable := latest != nil && latest.FullPrettyVersion(true) != p.FullPrettyVersion(true)
	if withheld != nil && !hasInstallable {
		j.Set("latest", withheld.p.PrettyVersion())
		j.Set("latest-status", "cooldown")
		j.Set("cooldown-available-in", withheld.availableIn)
		if latest == nil {
			latest = p
		}
	} else if latest != nil {
		j.Set("latest", latest.PrettyVersion())
		if withheld != nil {
			j.Set("cooldown-latest", withheld.p.PrettyVersion())
			j.Set("cooldown-available-in", withheld.availableIn)
		}
	} else {
		latest = p
	}
	if p.SourceType() != "" {
		m := phpjson.NewMap()
		m.Set("type", p.SourceType())
		m.Set("url", nullable(p.SourceURL()))
		m.Set("reference", nullable(p.SourceReference()))
		j.Set("source", m)
	}
	if p.DistType() != "" {
		m := phpjson.NewMap()
		m.Set("type", p.DistType())
		m.Set("url", nullable(p.DistURL()))
		m.Set("reference", nullable(p.DistReference()))
		j.Set("dist", m)
	}
	if !repo.IsPlatformPackage(p.Name) && sc.installed.HasPackage(p) {
		if path := sc.installPath(p); path != "" {
			if rp := realpath(path); rp != "" {
				j.Set("path", rp)
			}
		} else {
			j.Set("path", nil)
		}
		if p.ReleaseDate() != "" {
			if t, ok := pkg.ParsePHPDate(p.ReleaseDate()); ok {
				j.Set("released", t.Format("2006-01-02T15:04:05-07:00"))
			}
		}
	}
	if latest.IsAbandoned() {
		j.Set("replacement", nullable(latest.ReplacementPackage()))
	}
	if s := p.Suggests(); s != nil && s.Len() > 0 {
		j.Set("suggests", s)
	}
	if s := p.Support(); s != nil && s.Len() > 0 {
		j.Set("support", s)
	}
	if al := p.Autoload(); al != nil && al.Len() > 0 {
		auto := phpjson.NewMap()
		for _, typ := range al.Keys() {
			v, _ := al.Get(typ)
			switch typ {
			case "psr-0", "psr-4":
				psr := phpjson.NewMap()
				phpjson.Each(v, func(ns string, path phpjson.Value) {
					if !phpjson.Truthy(path) {
						path = "."
					}
					key := ns
					if !phpjson.Truthy(ns) {
						key = "*"
					}
					psr.Set(key, path)
				})
				auto.Set(typ, psr)
			case "classmap":
				auto.Set("classmap", v)
			}
		}
		j.Set("autoload", auto)
	}
	if inc := p.IncludePaths(); len(inc) > 0 {
		j.Set("include_path", strList(inc))
	}
	for _, lt := range []struct {
		key   string
		links pkg.Links
	}{{"requires", p.Requires()}, {"devRequires", p.DevRequires()}, {"provides", p.Provides()}, {"conflicts", p.Conflicts()}, {"replaces", p.Replaces()}} {
		if len(lt.links) == 0 {
			continue
		}
		m := phpjson.NewMap()
		for _, l := range lt.links {
			m.Set(l.Target, l.PrettyConstraint)
		}
		j.Set(lt.key, m)
	}
	sc.ctx.IO.Write(phpjson.MustEncode(j, phpjson.ComposerFlags))
}

func sortedLinks(links pkg.Links) pkg.Links {
	out := append(pkg.Links(nil), links...)
	sort.SliceStable(out, func(i, j int) bool { return phpjson.PHPCompareKeys(out[i].Target, out[j].Target) < 0 })
	return out
}

func (sc *showCtx) installedVersion(name string) (string, bool) {
	for _, p := range sc.installed.FindPackagesWithReplacersAndProviders(name, nil) {
		for p.AliasOf != nil {
			p = p.AliasOf
		}
		return p.PrettyVersion(), true
	}
	return "", false
}

func (sc *showCtx) generatePackageTree(p *pkg.Package) *phpjson.Map {
	var children []phpjson.Value
	for _, req := range sortedLinks(p.Requires()) {
		inTree := []string{p.Name, req.Target}
		d := phpjson.NewMap()
		d.Set("name", req.Target)
		d.Set("version", req.PrettyConstraint)
		if v, ok := sc.installedVersion(req.Target); ok {
			d.Set("installed", v)
		}
		if deep := sc.addTree(req.Target, req, inTree); len(deep) > 0 {
			d.Set("requires", deep)
		}
		children = append(children, d)
	}
	tree := phpjson.NewMap()
	tree.Set("name", p.PrettyName)
	tree.Set("version", p.PrettyVersion())
	tree.Set("description", p.Description())
	if len(children) > 0 {
		tree.Set("requires", children)
	}
	return tree
}

func (sc *showCtx) addTree(name string, link *pkg.Link, inTree []string) []phpjson.Value {
	var children []phpjson.Value
	var p *pkg.Package
	if link.PrettyConstraint == "self.version" {
		p, _, _ = sc.getPackageConstraint(strings.ToLower(name), link.Constraint, true)
	} else {
		p, _, _ = sc.getPackage(name, link.PrettyConstraint)
	}
	if p == nil {
		return nil
	}
	for _, req := range sortedLinks(p.Requires()) {
		cur := append([]string(nil), inTree...)
		d := phpjson.NewMap()
		d.Set("name", req.Target)
		d.Set("version", req.PrettyConstraint)
		if v, ok := sc.installedVersion(req.Target); ok {
			d.Set("installed", v)
		}
		if !inList(req.Target, cur) {
			cur = append(cur, req.Target)
			if deep := sc.addTree(req.Target, req, cur); len(deep) > 0 {
				d.Set("requires", deep)
			}
		}
		children = append(children, d)
	}
	return children
}

func (sc *showCtx) writeTreeLine(line string) {
	io := sc.ctx.IO
	if !io.Decorated {
		line = strings.NewReplacer("└", "`-", "├", "|-", "──", "-", "│", "|").Replace(line)
	}
	io.Write(line)
}

func (sc *showCtx) displayPackageTree(trees []*phpjson.Map) {
	io := sc.ctx.IO
	for _, t := range trees {
		io.WriteV("<info>"+t.Str("name")+"</info>", false, 32)
		io.WriteV(" "+t.Str("version"), false, 32)
		if t.Has("description") {
			io.Write(" " + strtok(t.Str("description")))
		} else {
			io.Write("")
		}
		reqs, _ := t.Get("requires")
		list, _ := reqs.([]phpjson.Value)
		treeBar := "├"
		for j, r := range list {
			rm := phpjson.AsMap(r)
			if j == len(list)-1 {
				treeBar = "└"
			}
			color := treeColors[1]
			sc.writeTreeLine(fmt.Sprintf("%s──<%s>%s</%s> %s", treeBar, color, rm.Str("name"), color, rm.Str("version")))
			treeBar = strings.ReplaceAll(treeBar, "└", " ")
			sc.displayTree(rm, []string{t.Str("name"), rm.Str("name")}, treeBar, 2)
		}
	}
}

// strtok($s, "\r\n") returns the first non-empty token.
func strtok(s string) string {
	s = strings.TrimLeft(s, "\r\n")
	return firstLine(s)
}

func (sc *showCtx) displayTree(p *phpjson.Map, inTree []string, previousTreeBar string, level int) {
	previousTreeBar = strings.ReplaceAll(previousTreeBar, "├", "│")
	reqs, ok := p.Get("requires")
	if !ok {
		return
	}
	list, _ := reqs.([]phpjson.Value)
	treeBar := previousTreeBar + "  ├"
	for i, r := range list {
		rm := phpjson.AsMap(r)
		cur := append([]string(nil), inTree...)
		if i == len(list)-1 {
			treeBar = previousTreeBar + "  └"
		}
		color := treeColors[level%len(treeColors)]
		circular := ""
		if inList(rm.Str("name"), cur) {
			circular = "(circular dependency aborted here)"
		}
		info := strings.TrimRight(fmt.Sprintf("%s──<%s>%s</%s> %s %s", treeBar, color, rm.Str("name"), color, rm.Str("version"), circular), " \t\n\r\x00\x0B")
		sc.writeTreeLine(info)
		treeBar = strings.ReplaceAll(treeBar, "└", " ")
		cur = append(cur, rm.Str("name"))
		sc.displayTree(rm, cur, treeBar, level+1)
	}
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

func strPad(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

func strWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

func runeWidth(r rune) int {
	if (r >= 0x1100 && r <= 0x115F) || (r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) || (r >= 0xAC00 && r <= 0xD7A3) || (r >= 0xF900 && r <= 0xFAFF) || (r >= 0xFE30 && r <= 0xFE4F) || (r >= 0xFF00 && r <= 0xFF60) || (r >= 0xFFE0 && r <= 0xFFE6) || (r >= 0x20000 && r <= 0x3FFFD) {
		return 2
	}
	return 1
}

// strimwidth is mb_strimwidth($s, 0, $width, $marker).
func strimwidth(s string, width int, marker string) string {
	if strWidth(s) <= width {
		return s
	}
	target := width - strWidth(marker)
	w := 0
	var sb strings.Builder
	for _, r := range s {
		rw := runeWidth(r)
		if w+rw > target {
			break
		}
		w += rw
		sb.WriteRune(r)
	}
	return sb.String() + marker
}

// relativeTime is ShowCommand::getRelativeTime.
func relativeTime(t time.Time) string {
	now := time.Now()
	if t.Format("20060102") == now.Format("20060102") {
		return "today"
	}
	days := int(now.Sub(t).Hours() / 24)
	y, m := diffYM(t, now)
	switch {
	case days < 7:
		return "this week"
	case days < 14:
		return "last week"
	case m < 1 && y < 1 && days < 31:
		return strconv.Itoa(days/7) + " weeks ago"
	case y < 1:
		s := ""
		if m > 1 {
			s = "s"
		}
		return strconv.Itoa(m) + " month" + s + " ago"
	}
	s := ""
	if y > 1 {
		s = "s"
	}
	return strconv.Itoa(y) + " year" + s + " ago"
}

// diffYM is the y/m part of DateTime::diff.
func diffYM(a, b time.Time) (int, int) {
	a, b = a.UTC(), b.UTC()
	months := (b.Year()-a.Year())*12 + int(b.Month()) - int(a.Month())
	if b.Day() < a.Day() || (b.Day() == a.Day() && b.Sub(time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)) < a.Sub(time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, time.UTC))) {
		months--
	}
	if months < 0 {
		months = 0
	}
	return months / 12, months % 12
}
