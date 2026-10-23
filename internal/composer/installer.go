package composer

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"gocomposer/internal/autoload"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/installed"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/locker"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	pol "gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// Installer exit codes (Installer::ERROR_*).
const (
	ErrNone                 = 0
	ErrGeneric              = 1
	ErrDependencyResolution = 2
	ErrNoLockForPartial     = 3
	ErrLockFileInvalid      = 4
	ErrAuditFailed          = 5
	ErrPsrAutoloadViolation = 6
	ErrTransport            = 100
)

// Installer ports Composer\Installer.
type Installer struct {
	C                     *Composer
	IM                    *InstallationManager
	Events                *EventDispatcher
	PreferSource          bool
	PreferDist            bool
	OptimizeAutoloader    bool
	ClassMapAuthoritative bool
	ApcuAutoloader        bool
	ApcuPrefix            string
	StrictPsr             bool
	DevMode               bool
	DryRun                bool
	DownloadOnly          bool
	Update                bool
	Install               bool
	DumpAutoloader        bool
	RunScripts            bool
	PreferStable          bool
	PreferLowest          bool
	MinimalUpdate         bool
	WriteLock             bool
	ExecuteOperations     bool
	Audit                 bool
	ErrorOnAudit          bool
	// MockAutoloadGenerator skips writing the autoloader (tests, like the
	// AutoloadGenerator mock of InstallerTest) while keeping the output.
	MockAutoloadGenerator bool
	// OnPreOperationsExec is called where Composer dispatches
	// InstallerEvents::PRE_OPERATIONS_EXEC.
	OnPreOperationsExec func()
	// PolicyConfig is Installer::setPolicyConfig (nil: from config).
	PolicyConfig           *pol.PolicyConfig
	AuditFormat            string
	UpdateMirrors          bool
	UpdateAllowList        []string
	UpdateAllowListSet     bool
	UpdateAllowTransitive  int
	PlatformFilter         solver.PlatformFilter
	TemporaryConstraints   map[string]semver.Constraint
	IgnoredTypes           []string
	AllowedTypes           []string
	SuggestionsFromInstall []*pkg.Package
	fixedRootPackage       *pkg.Package
	lockTransaction        *solver.LockTransaction
	suggestions            *SuggestedPackagesReporter
	localRepo              *installed.Repository
}

// NewInstaller is Installer::create with Composer's defaults.
func NewInstaller(c *Composer) *Installer {
	in := &Installer{C: c, IM: NewInstallationManager(c), Events: NewEventDispatcher(c), Install: true, DumpAutoloader: true, RunScripts: true,
		WriteLock: c.Config.GetBool("lock"), ExecuteOperations: true, Audit: true, AuditFormat: "summary", PlatformFilter: solver.IgnoreNothing{},
		IgnoredTypes: []string{"php-ext", "php-ext-zend"}, suggestions: NewSuggestedPackagesReporter(c.IO)}
	in.IM.PackageEvent = func(event string, op *solver.Operation) error {
		return in.Events.DispatchPackageEvent(event, in.DevMode, op)
	}
	return in
}

// Run is Installer::run.
func (in *Installer) Run() (int, error) {
	// plugin code run meanwhile finds this Installer on its call stack (plugin_api.go)
	defer pluginInstallerScope(in)()
	c := in.C
	io := c.IO
	if in.UpdateAllowListSet && in.UpdateMirrors {
		return 1, &UserError{Msg: "The installer options updateMirrors and updateAllowList are mutually exclusive.", Code: 1}
	}
	localRepo := c.LocalRepo
	in.localRepo = localRepo
	// suggestions collected before the run (create-project's root package,
	// given to Installer::setSuggestedPackagesReporter)
	for _, p := range in.SuggestionsFromInstall {
		in.suggestions.AddSuggestionsFromPackage(p)
	}
	isFreshInstall := localRepo.IsFresh()
	locked, lockErr := c.Locker.IsLocked()
	if !in.Update && (lockErr != nil || !locked) {
		if lockErr != nil {
			return 1, lockErr
		}
		msg := "Updating dependencies to latest instead of installing from lock file. See https://getcomposer.org/install for more information."
		if !in.WriteLock {
			io.WriteError("<warning>Lockfile creation is disabled in config. " + msg + "</warning>")
		} else {
			io.WriteError("<warning>No composer.lock file present. " + msg + "</warning>")
		}
		in.Update = true
	}
	if in.DryRun {
		in.RunScripts = false
		in.ExecuteOperations = false
		in.WriteLock = false
		in.DumpAutoloader = false
		in.mockLocalRepository()
		localRepo = in.localRepo
	}
	if in.DownloadOnly {
		in.DumpAutoloader = false
	}
	if in.Update && !in.Install {
		in.DumpAutoloader = false
	}
	in.Events.RunScripts = in.RunScripts && !c.DisableScripts
	if in.RunScripts {
		if in.DevMode {
			os.Setenv("COMPOSER_DEV_MODE", "1")
		} else {
			os.Setenv("COMPOSER_DEV_MODE", "0")
		}
		ev := PreInstallCmd
		if in.Update {
			ev = PreUpdateCmd
		}
		if _, err := in.Events.DispatchScript(ev, in.DevMode, nil, nil); err != nil {
			return scriptExitCode(err), err
		}
	}
	in.IM.SetPreferSource(in.PreferSource)
	in.IM.SetPreferDist(in.PreferDist)
	var res int
	var err error
	if in.Update {
		res, err = in.doUpdate(localRepo, in.Install)
	} else {
		res, err = in.doInstall(localRepo, false)
	}
	if err != nil {
		return res, err
	}
	if res != 0 {
		return res, nil
	}
	if in.Update {
		if isFreshInstall {
			in.suggestions.AddSuggestionsFromPackage(c.Root())
		}
		if err := in.outputSuggestions(); err != nil {
			return 1, err
		}
	}
	lockedRepo, err := c.Locker.LockedRepository(true)
	if err == nil {
		for _, p := range lockedRepo.MustPackages() {
			if p.AliasOf != nil || !p.IsAbandoned() {
				continue
			}
			replacement := "No replacement was suggested"
			if r := p.ReplacementPackage(); r != "" {
				replacement = "Use " + r + " instead"
			}
			io.WriteError(fmt.Sprintf("<warning>Package %s is abandoned, you should avoid using it. %s.</warning>", p.PrettyName, replacement))
		}
	}
	if in.DumpAutoloader {
		if in.OptimizeAutoloader {
			io.WriteError("<info>Generating optimized autoload files</info>")
		} else {
			io.WriteError("<info>Generating autoload files</info>")
		}
		if !in.MockAutoloadGenerator {
			result, err := in.dumpAutoload(localRepo)
			if err != nil {
				return 1, err
			}
			if in.StrictPsr && len(result.PsrViolations) > 0 {
				return ErrPsrAutoloadViolation, nil
			}
		}
	}
	if in.Install && in.ExecuteOperations {
		for _, p := range localRepo.MustPackages() {
			in.IM.EnsureBinariesPresence(p)
		}
	}
	showFunding := true
	if v := os.Getenv("COMPOSER_FUND"); v != "" {
		if n, err := fmt.Sscanf(v, "%d", new(int)); err == nil && n == 1 {
			var x int
			fmt.Sscanf(v, "%d", &x)
			showFunding = x != 0
		}
	}
	if showFunding {
		count := 0
		for _, p := range localRepo.MustPackages() {
			if p.AliasOf == nil && len(p.Funding()) > 0 {
				count++
			}
		}
		if count > 0 {
			s, verb := "s", "are"
			if count == 1 {
				s, verb = "", "is"
			}
			io.WriteError(fmt.Sprintf("<info>%d package%s you are using %s looking for funding.</info>", count, s, verb))
			io.WriteError("<info>Use the `composer fund` command to find out more!</info>")
		}
	}
	if in.RunScripts {
		ev := PostInstallCmd
		if in.Update {
			ev = PostUpdateCmd
		}
		if _, err := in.Events.DispatchScript(ev, in.DevMode, nil, nil); err != nil {
			return scriptExitCode(err), err
		}
	}
	if in.Audit {
		var packages []*pkg.Package
		target := "installed"
		if in.Update && !in.Install {
			if lr, err := c.Locker.LockedRepository(true); err == nil {
				packages = lr.MustPackages()
			}
			target = "locked"
		} else {
			packages = localRepo.MustPackages()
		}
		if len(packages) > 0 {
			code, err := in.runPolicyAudit(packages)
			var te *httpdl.TransportError
			if errors.As(err, &te) {
				io.WriteError("<error>Failed to audit " + target + " packages.</error>")
				if io.IsVerbose() {
					io.WriteError("<error>[Composer\\Downloader\\TransportException] " + te.Error() + "</error>")
				}
			} else if err != nil {
				return 1, err
			} else if code > 0 && in.ErrorOnAudit {
				return ErrAuditFailed, nil
			}
		} else {
			io.WriteError("No " + target + " packages - skipping audit.")
		}
	}
	return 0, nil
}

func scriptExitCode(err error) int {
	if se, ok := err.(*ScriptError); ok {
		return se.Code
	}
	return 1
}

func (in *Installer) mockLocalRepository() {
	mock := installed.NewInMemory()
	byString := map[string]*pkg.Package{}
	var order []*pkg.Package
	for _, p := range in.C.LocalRepo.MustPackages() {
		c := p.Clone()
		byString[p.String()] = c
		order = append(order, c)
	}
	for _, p := range order {
		mock.AddPackage(p)
	}
	in.localRepo = mock
}

// DumpAutoload runs the autoload generator for the installer state.
func (in *Installer) dumpAutoload(localRepo *installed.Repository) (*autoload.Result, error) {
	// like Composer the generator infers dev mode from the installed.json just written
	if err := localRepo.Reload(); err != nil {
		return nil, err
	}
	return RunAutoloadDump(in.C, in.Events, localRepo, AutoloadOptions{
		DevMode:               nil,
		Optimize:              in.OptimizeAutoloader,
		ClassMapAuthoritative: in.ClassMapAuthoritative,
		Apcu:                  in.ApcuAutoloader,
		ApcuPrefix:            in.ApcuPrefix,
		RunScripts:            in.RunScripts,
		PlatformFilter:        in.PlatformFilter,
		IM:                    in.IM,
	})
}

// AutoloadOptions configures RunAutoloadDump.
type AutoloadOptions struct {
	DevMode               *bool
	Optimize              bool
	ClassMapAuthoritative bool
	Apcu                  bool
	ApcuPrefix            string
	RunScripts            bool
	PlatformFilter        solver.PlatformFilter
	IM                    *InstallationManager
	DryRun                bool
	StrictAmbiguous       bool
	Suffix                string
}

// RunAutoloadDump wraps AutoloadGenerator::dump with its script events.
func RunAutoloadDump(c *Composer, events *EventDispatcher, localRepo *installed.Repository, o AutoloadOptions) (*autoload.Result, error) {
	devMode := false
	if o.DevMode != nil {
		devMode = *o.DevMode
	} else if localRepo.DevMode != nil {
		devMode = *localRepo.DevMode
	}
	if o.RunScripts {
		if os.Getenv("COMPOSER_DEV_MODE") == "" {
			if devMode {
				os.Setenv("COMPOSER_DEV_MODE", "1")
			} else {
				os.Setenv("COMPOSER_DEV_MODE", "0")
			}
		}
		flags := map[string]any{"optimize": o.Optimize || o.ClassMapAuthoritative}
		if _, err := events.DispatchScript(PreAutoloadDump, devMode, nil, flags); err != nil {
			return nil, err
		}
	}
	im := o.IM
	if im == nil {
		im = NewInstallationManager(c)
	}
	contentHash := ""
	if locked, _ := c.Locker.IsLocked(); locked {
		if data, err := c.Locker.LockData(); err == nil {
			contentHash = data.Str("content-hash")
		}
	}
	filter := o.PlatformFilter
	if filter == nil {
		filter = solver.IgnoreNothing{}
	}
	_, ignoreAll := filter.(solver.IgnoreAll)
	suffix := ""
	if s, ok := c.Config.Get("autoloader-suffix").(string); ok {
		suffix = s
	}
	prepend := c.Config.Get("prepend-autoloader") != false
	opts := autoload.Options{
		VendorDir:             c.VendorDir(),
		Cwd:                   c.Cwd,
		DevMode:               &devMode,
		ClassMapAuthoritative: o.ClassMapAuthoritative,
		Apcu:                  o.Apcu,
		ApcuPrefix:            o.ApcuPrefix,
		Optimize:              o.Optimize,
		Suffix:                o.Suffix,
		ConfigSuffix:          suffix,
		LockContentHash:       contentHash,
		UseIncludePath:        c.Config.GetBool("use-include-path"),
		PrependAutoloader:     prepend,
		PlatformCheck:         c.Config.Get("platform-check"),
		IgnorePlatform:        filter.IsIgnored,
		IgnoreAllPlatform:     ignoreAll,
		StrictAmbiguous:       o.StrictAmbiguous,
		DryRun:                o.DryRun,
		Warn:                  func(s string) { c.IO.WriteError(s) },
	}
	result, err := autoload.Dump(opts, c.Root(), localRepo.CanonicalPackages(), im.InstallPath, localRepo.DevPackageNames)
	if err != nil {
		return nil, err
	}
	if o.RunScripts && !o.DryRun {
		flags := map[string]any{"optimize": o.Optimize || o.ClassMapAuthoritative}
		if _, err := events.DispatchScript(PostAutoloadDump, devMode, nil, flags); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (in *Installer) platformOverrides() phpjson.Value {
	v := in.C.Config.Get("platform")
	if v == nil {
		return phpjson.NewMap()
	}
	return v
}

func (in *Installer) createPlatformRepo(forUpdate bool) (*platform.Repository, error) {
	var overrides phpjson.Value
	if forUpdate {
		overrides = in.platformOverrides()
	} else {
		overrides = in.C.Locker.PlatformOverrides()
	}
	r, err := in.C.CreatePlatformRepo(overrides)
	if err != nil {
		return nil, err
	}
	if r.Facts() == nil {
		if _, ok := in.PlatformFilter.(solver.IgnoreAll); !ok {
			if !hasPHPOverride(overrides) {
				_, _, perr := in.C.PHP()
				msg := "gocomposer needs a PHP binary to determine platform packages (php, ext-*). Set GOCOMPOSER_PHP or add php to PATH, configure config.platform.php, or use --ignore-platform-reqs."
				if perr != nil {
					msg += " (" + perr.Error() + ")"
				}
				return nil, &UserError{Msg: msg, Code: 1}
			}
		}
	}
	return r, nil
}

func hasPHPOverride(v phpjson.Value) bool {
	found := false
	phpjson.Each(v, func(k string, val phpjson.Value) {
		if strings.EqualFold(k, "php") {
			if _, ok := val.(string); ok {
				found = true
			}
		}
	})
	return found
}

func (in *Installer) createRepositorySet(forUpdate bool, platformRepo *platform.Repository, rootAliases []pkg.RootAlias, lockedRepo *repo.Array) *solver.RepositorySet {
	root := in.C.Root()
	var minimumStability string
	var stabilityFlags map[string]int
	var requires []solver.RequireEntry
	if forUpdate {
		minimumStability = root.MinimumStability
		stabilityFlags = copyFlags(root.StabilityFlags)
		for _, l := range append(append(pkg.Links(nil), root.Requires()...), root.DevRequires()...) {
			requires = setRequire(requires, l.Target, l.Constraint)
		}
	} else {
		minimumStability = in.C.Locker.MinimumStability()
		stabilityFlags = in.C.Locker.StabilityFlags()
		for _, p := range lockedRepo.MustPackages() {
			c := semver.NewConstraint("=", p.Version())
			c.SetPrettyString(p.PrettyVersion())
			requires = setRequire(requires, p.Name, c)
		}
	}
	var rootRequires []solver.RequireEntry
	for _, r := range requires {
		if in.PlatformFilter.IsIgnored(r.Name) {
			continue
		}
		c := r.Constraint
		if il, ok := in.PlatformFilter.(*solver.IgnoreList); ok {
			c = il.FilterConstraint(r.Name, c, true)
		}
		rootRequires = append(rootRequires, solver.RequireEntry{Name: r.Name, Constraint: c})
	}
	fixed := in.C.Package.Clone()
	fixed.SetRequires(nil)
	fixed.SetDevRequires(nil)
	if fixed.AliasOf != nil {
		base := fixed.AliasOf.Clone()
		base.SetRequires(nil)
		base.SetDevRequires(nil)
		fixed = pkg.NewAlias(base, fixed.Version(), fixed.PrettyVersion())
		fixed.IsRoot = true
	}
	in.fixedRootPackage = fixed
	stabilityFlags[root.Name] = pkg.Stabilities[semver.ParseStability(in.C.Package.Version())]
	rs := solver.NewRepositorySet(minimumStability, stabilityFlags, rootAliases, root.References, rootRequires, in.TemporaryConstraints)
	rs.AddRepository(repo.NewRootPackageRepository(fixed))
	rs.AddRepository(platformRepo)
	return rs
}

func copyFlags(m map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func setRequire(list []solver.RequireEntry, name string, c semver.Constraint) []solver.RequireEntry {
	for i, e := range list {
		if e.Name == name {
			list[i].Constraint = c
			return list
		}
	}
	return append(list, solver.RequireEntry{Name: name, Constraint: c})
}

func (in *Installer) createPolicy(forUpdate bool, req *solver.Request) *solver.DefaultPolicy {
	var preferStable, preferLowest *bool
	if !forUpdate {
		preferStable = in.C.Locker.PreferStable()
		preferLowest = in.C.Locker.PreferLowest()
	}
	ps := in.PreferStable || in.C.Root().PreferStable
	if preferStable != nil {
		ps = *preferStable
	}
	pl := in.PreferLowest
	if preferLowest != nil {
		pl = *preferLowest
	}
	var preferred map[string]string
	if forUpdate && in.MinimalUpdate && req != nil && req.LockedRepository != nil {
		preferred = map[string]string{}
		for _, p := range req.LockedRepository.MustPackages() {
			if p.AliasOf != nil || req.IsInUpdateAllowList(p.Name) {
				continue
			}
			preferred[p.Name] = p.Version()
		}
	}
	return solver.NewDefaultPolicy(ps, pl, preferred)
}

func (in *Installer) createRequest(root *pkg.Package, platformRepo *platform.Repository, lockedRepo *repo.Array) *solver.Request {
	req := solver.NewRequest(lockedRepo)
	req.FixPackage(root)
	if root.AliasOf != nil {
		req.FixPackage(root.AliasOf)
	}
	provided := root.Provides()
	for _, p := range platformRepo.MustPackages() {
		l := provided.Get(p.Name)
		if l == nil || !l.Constraint.Matches(semver.NewConstraint("=", p.Version())) {
			req.FixPackage(p)
		}
	}
	return req
}

func (in *Installer) requirePackagesForUpdate(req *solver.Request, lockedRepo *repo.Array, includeDev bool) {
	if in.UpdateMirrors {
		excluded := map[string]bool{}
		if !includeDev {
			for _, n := range in.C.Locker.DevPackageNames() {
				excluded[n] = true
			}
		}
		for _, p := range lockedRepo.MustPackages() {
			if p.AliasOf == nil && !excluded[p.Name] {
				req.RequireName(p.Name, semver.NewConstraint("==", p.Version()))
			}
		}
		return
	}
	root := in.C.Root()
	links := append(pkg.Links(nil), root.Requires()...)
	if includeDev {
		for _, l := range root.DevRequires() {
			replaced := false
			for i, e := range links {
				if e.Target == l.Target {
					links[i] = l
					replaced = true
				}
			}
			if !replaced {
				links = append(links, l)
			}
		}
	}
	for _, l := range links {
		req.RequireName(l.Target, l.Constraint)
	}
}

func (in *Installer) problemContext(rs *solver.RepositorySet, req *solver.Request, pool *solver.Pool) *solver.Context {
	_, facts, _ := in.C.PHP()
	return &solver.Context{RepoSet: rs, Request: req, Pool: pool, Verbose: in.C.IO.IsVerbose(), Platform: PlatformInfo{facts}, CompleteAdvisories: in.completeAdvisories(rs)}
}

func (in *Installer) newPoolBuilder(pol *solver.DefaultPolicy, optimize bool, scope string, update bool) (*solver.PoolBuilder, error) {
	b := solver.NewPoolBuilder()
	if optimize && os.Getenv("COMPOSER_POOL_OPTIMIZER") != "0" {
		b.Optimizer = solver.NewPoolOptimizer(pol)
	}
	f, err := in.poolFilter(scope, update)
	if err != nil {
		return nil, err
	}
	b.PolicyFilter = f
	b.IgnoredTypes = in.IgnoredTypes
	b.AllowedTypes = in.AllowedTypes
	b.Warn = func(s string) { in.C.IO.WriteError(s) }
	b.Write = func(s string, v int) { in.C.IO.WriteV(s, true, v) }
	// PluginEvents::PRE_POOL_CREATE (plugin_dlevents.go)
	b.PrePoolCreate = prePoolCreateHook(in.C)
	return b, nil
}

func (in *Installer) doUpdate(localRepo *installed.Repository, doInstall bool) (int, error) {
	c := in.C
	io := c.IO
	platformRepo, err := in.createPlatformRepo(true)
	if err != nil {
		return 1, err
	}
	aliases := c.RootAliases()
	var lockedRepo *repo.Array
	if locked, err := c.Locker.IsLocked(); err == nil && locked {
		lockedRepo, err = c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
	} else if err != nil {
		if _, ok := err.(*jsonfile.ParsingError); !ok || in.UpdateAllowListSet || in.UpdateMirrors {
			return 1, err
		}
	}
	if (in.UpdateAllowListSet || in.UpdateMirrors) && lockedRepo == nil {
		what := "only a partial set of packages"
		if in.UpdateMirrors {
			what = "lock file information"
		}
		io.WriteErrorV("<error>Cannot update "+what+" without a lock file present. Run `composer update` to generate a lock file.</error>", true, 16)
		return ErrNoLockForPartial, nil
	}
	io.WriteError("<info>Loading composer repositories with package information</info>")
	rs := in.createRepositorySet(true, platformRepo, aliases, nil)
	for _, r := range c.Repositories {
		rs.AddRepository(r)
	}
	if lockedRepo != nil {
		rs.AddRepository(lockedRepo)
	}
	req := in.createRequest(in.fixedRootPackage, platformRepo, lockedRepo)
	in.requirePackagesForUpdate(req, lockedRepo, true)
	if in.UpdateAllowListSet {
		req.SetUpdateAllowList(in.UpdateAllowList, in.UpdateAllowTransitive)
	}
	policy := in.createPolicy(true, req)
	pb, err := in.newPoolBuilder(policy, true, pol.BlockScopeUpdate, true)
	if err != nil {
		return 1, err
	}
	pool, err := rs.CreatePool(req, pb)
	if err != nil {
		return transportExit(err), err
	}
	io.WriteError("<info>Updating dependencies</info>")
	sv := solver.NewSolver(policy, pool)
	sv.Log = func(m string, v int) { in.C.IO.WriteErrorV(m, true, v) }
	lt, err := sv.Solve(req, in.PlatformFilter)
	if err != nil {
		if pe, ok := err.(*solver.ProblemsError); ok {
			msg := "Your requirements could not be resolved to an installable set of packages."
			io.WriteErrorV("<error>"+msg+"</error>", true, 16)
			io.WriteError(pe.PrettyString(in.problemContext(rs, req, pool), false))
			if !in.DevMode {
				io.WriteErrorV("<warning>Running update with --no-dev does not mean require-dev is ignored, it just means the packages will not be installed. If dev requirements are blocking the update you have to resolve those problems.</warning>", true, 16)
			}
			return ErrDependencyResolution, nil
		}
		return 1, err
	}
	in.lockTransaction = lt
	io.WriteErrorV(fmt.Sprintf("Analyzed %d packages to resolve dependencies", pool.Count()), true, 64)
	io.WriteErrorV(fmt.Sprintf("Analyzed %d rules to resolve dependencies", sv.RuleSetSize()), true, 64)
	if len(lt.Operations()) == 0 {
		io.WriteError("Nothing to modify in lock file")
		if in.MinimalUpdate && !in.UpdateAllowListSet && c.Locker.IsFresh() {
			io.WriteError("<warning>The --minimal-changes option should be used with package arguments or after modifying composer.json requirements, otherwise it will likely not yield any dependency changes.</warning>")
		}
	}
	if code, err := in.extractDevPackages(lt, platformRepo, aliases, policy, lockedRepo); err != nil || code != 0 {
		return code, err
	}
	semver.ClearMatchCache()
	platformReqs := extractPlatformRequirements(c.Root().Requires())
	platformDevReqs := extractPlatformRequirements(c.Root().DevRequires())
	var installsUpdates, uninstalls []*solver.Operation
	if ops := lt.Operations(); len(ops) > 0 {
		var installNames, updateNames, uninstallNames []string
		for _, op := range ops {
			switch op.Type {
			case solver.OpInstall:
				installsUpdates = append(installsUpdates, op)
				installNames = append(installNames, op.Package.PrettyName+":"+op.Package.FullPrettyVersion(true))
			case solver.OpUpdate:
				if in.UpdateMirrors && op.Initial.Name == op.Target.Name && op.Initial.Version() == op.Target.Version() {
					continue
				}
				installsUpdates = append(installsUpdates, op)
				updateNames = append(updateNames, op.Target.PrettyName+":"+op.Target.FullPrettyVersion(true))
			case solver.OpUninstall:
				uninstalls = append(uninstalls, op)
				uninstallNames = append(uninstallNames, op.Package.PrettyName)
			}
		}
		if c.Config.GetBool("lock") {
			io.WriteError(fmt.Sprintf("<info>Lock file operations: %d install%s, %d update%s, %d removal%s</info>", len(installNames), plural(len(installNames)), len(updateNames), plural(len(updateNames)), len(uninstalls), plural(len(uninstalls))))
			if len(installNames) > 0 {
				io.WriteErrorV("Installs: "+strings.Join(installNames, ", "), true, 64)
			}
			if len(updateNames) > 0 {
				io.WriteErrorV("Updates: "+strings.Join(updateNames, ", "), true, 64)
			}
			if len(uninstalls) > 0 {
				io.WriteErrorV("Removals: "+strings.Join(uninstallNames, ", "), true, 64)
			}
		}
	}
	byName := func(ops []*solver.Operation) {
		sort.SliceStable(ops, func(i, j int) bool { return ops[i].MainPackage().Name < ops[j].MainPackage().Name })
	}
	byName(uninstalls)
	byName(installsUpdates)
	for _, op := range append(append([]*solver.Operation(nil), uninstalls...), installsUpdates...) {
		if op.Type == solver.OpInstall {
			in.suggestions.AddSuggestionsFromPackage(op.Package)
		}
		if c.Config.GetBool("lock") && (!strings.Contains(op.Type, "Alias") || io.IsDebug()) {
			src := ""
			if io.IsVeryVerbose() && !strings.Contains(op.Type, "Alias") {
				if r, ok := op.MainPackage().Repo.(repo.Repository); ok {
					src = " from " + r.RepoName()
				}
			}
			io.WriteError("  - " + op.Show(true) + src)
		}
	}
	devPackages := lt.NewLockPackagesMirrors(true, in.UpdateMirrors)
	stab := c.Root().StabilityFlags
	updated, err := c.Locker.SetLockData(locker.LockOptions{
		Packages:          lt.NewLockPackagesMirrors(false, in.UpdateMirrors),
		DevPackages:       devPackages,
		DevPackagesSet:    true,
		PlatformReqs:      platformReqs,
		PlatformDevReqs:   platformDevReqs,
		Aliases:           lt.Aliases(aliases),
		MinimumStability:  c.Root().MinimumStability,
		StabilityFlags:    stab,
		PreferStable:      in.PreferStable || c.Root().PreferStable,
		PreferLowest:      in.PreferLowest,
		PlatformOverrides: phpjson.ToMap(in.platformOverrides()),
		Write:             in.WriteLock && in.ExecuteOperations,
	})
	if err != nil {
		return 1, err
	}
	if updated && in.WriteLock && in.ExecuteOperations {
		io.WriteError("<info>Writing lock file</info>")
	}
	if doInstall {
		return in.doInstall(localRepo, true)
	}
	return 0, nil
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func transportExit(err error) int {
	return 1
}

func extractPlatformRequirements(links pkg.Links) *phpjson.Map {
	m := phpjson.NewMap()
	for _, l := range links {
		if repo.IsPlatformPackage(l.Target) {
			m.Set(l.Target, l.PrettyConstraint)
		}
	}
	return m
}

func (in *Installer) extractDevPackages(lt *solver.LockTransaction, platformRepo *platform.Repository, aliases []pkg.RootAlias, policy *solver.DefaultPolicy, lockedRepo *repo.Array) (int, error) {
	if len(in.C.Root().DevRequires()) == 0 {
		return 0, nil
	}
	resultRepo := repo.NewArray()
	loader := &pkg.Loader{LoadOptions: true}
	for _, p := range lt.NewLockPackages(false) {
		lp, err := loader.Load(pkg.Dump(p))
		if err != nil {
			return 1, err
		}
		resultRepo.AddPackage(lp)
	}
	rs := in.createRepositorySet(true, platformRepo, aliases, nil)
	rs.AddRepository(resultRepo)
	req := in.createRequest(in.fixedRootPackage, platformRepo, nil)
	in.requirePackagesForUpdate(req, lockedRepo, false)
	pool, err := rs.CreatePoolWithAllPackages()
	if err != nil {
		return 1, err
	}
	sv := solver.NewSolver(policy, pool)
	sv.Log = func(m string, v int) { in.C.IO.WriteErrorV(m, true, v) }
	nonDev, err := sv.Solve(req, in.PlatformFilter)
	if err != nil {
		if pe, ok := err.(*solver.ProblemsError); ok {
			io := in.C.IO
			io.WriteErrorV("<error>Unable to find a compatible set of packages based on your non-dev requirements alone.</error>", true, 16)
			io.WriteError("Your requirements can be resolved successfully when require-dev packages are present.")
			io.WriteError("You may need to move packages from require-dev or some of their dependencies to require.")
			io.WriteError(pe.PrettyString(in.problemContext(rs, req, pool), true))
			return ErrDependencyResolution, nil
		}
		return 1, err
	}
	lt.SetNonDevPackages(nonDev)
	return 0, nil
}

func (in *Installer) doInstall(localRepo *installed.Repository, alreadySolved bool) (int, error) {
	c := in.C
	io := c.IO
	if !alreadySolved {
		if data, err := c.Locker.LockData(); err == nil && data.Str("content-hash") == locker.ConflictedContentHash {
			io.WriteErrorV("<error>"+locker.ConflictedContentHash+"</error>", true, 16)
			return ErrLockFileInvalid, nil
		}
	}
	if c.Config.GetBool("lock") {
		suffix := ""
		if in.DevMode {
			suffix = " (including require-dev)"
		}
		io.WriteError("<info>Installing dependencies from lock file" + suffix + "</info>")
	}
	lockedRepo, err := c.Locker.LockedRepository(in.DevMode)
	if err != nil {
		return 1, err
	}
	if !alreadySolved {
		io.WriteError("<info>Verifying lock file contents can be installed on current platform.</info>")
		platformRepo, err := in.createPlatformRepo(false)
		if err != nil {
			return 1, err
		}
		policy := in.createPolicy(false, nil)
		rs := in.createRepositorySet(false, platformRepo, nil, lockedRepo)
		rs.AddRepository(lockedRepo)
		req := in.createRequest(in.fixedRootPackage, platformRepo, lockedRepo)
		if !c.Locker.IsFresh() {
			io.WriteErrorV("<warning>Warning: The lock file is not up to date with the latest changes in composer.json. You may be getting outdated dependencies. It is recommended that you run `composer update` or `composer update <package name>`.</warning>", true, 16)
		}
		missing, err := in.missingRequirementInfo()
		if err != nil {
			return 1, err
		}
		if len(missing) > 0 {
			for _, m := range missing {
				io.WriteError(m)
			}
			if !c.Config.GetBool("allow-missing-requirements") {
				return ErrLockFileInvalid, nil
			}
		}
		for _, p := range lockedRepo.MustPackages() {
			req.FixLockedPackage(p)
		}
		rootRequires := append(pkg.Links(nil), c.Root().Requires()...)
		if in.DevMode {
			for _, l := range c.Root().DevRequires() {
				replaced := false
				for i, e := range rootRequires {
					if e.Target == l.Target {
						rootRequires[i] = l
						replaced = true
					}
				}
				if !replaced {
					rootRequires = append(rootRequires, l)
				}
			}
		}
		for _, l := range rootRequires {
			if repo.IsPlatformPackage(l.Target) {
				req.RequireName(l.Target, l.Constraint)
			}
		}
		platformLinks, err := c.Locker.PlatformRequirements(in.DevMode)
		if err != nil {
			return 1, err
		}
		for _, l := range platformLinks {
			if rootRequires.Get(l.Target) == nil {
				req.RequireName(l.Target, l.Constraint)
			}
		}
		pb, err := in.newPoolBuilder(policy, false, pol.BlockScopeInstall, false)
		if err != nil {
			return 1, err
		}
		pool, err := rs.CreatePool(req, pb)
		if err != nil {
			return 1, err
		}
		sv := solver.NewSolver(policy, pool)
		sv.Log = func(m string, v int) { in.C.IO.WriteErrorV(m, true, v) }
		lt, err := sv.Solve(req, in.PlatformFilter)
		if err != nil {
			if pe, ok := err.(*solver.ProblemsError); ok {
				io.WriteErrorV("<error>Your lock file does not contain a compatible set of packages. Please run composer update.</error>", true, 16)
				io.WriteError(pe.PrettyString(in.problemContext(rs, req, pool), false))
				return ErrDependencyResolution, nil
			}
			return 1, err
		}
		if len(lt.Operations()) != 0 {
			io.WriteErrorV("<error>Your lock file cannot be installed on this system without changes. Please run composer update.</error>", true, 16)
			return ErrLockFileInvalid, nil
		}
	}
	tx := solver.NewLocalRepoTransaction(lockedRepo.MustPackages(), localRepo.MustPackages())
	// the transaction's packages for plugins (plugin_events.go)
	pluginTransactionPackages(c, lockedRepo.MustPackages(), localRepo.MustPackages())
	// InstallerEvents::PRE_OPERATIONS_EXEC
	if in.OnPreOperationsExec != nil {
		in.OnPreOperationsExec()
	}
	ops := tx.Operations()
	if err := in.Events.DispatchInstallerEvent("pre-operations-exec", in.DevMode, in.ExecuteOperations, ops); err != nil {
		return scriptExitCode(err), err
	}
	var installs, updates, uninstalls []string
	for _, op := range ops {
		switch op.Type {
		case solver.OpInstall:
			installs = append(installs, op.Package.PrettyName+":"+op.Package.FullPrettyVersion(true))
		case solver.OpUpdate:
			updates = append(updates, op.Target.PrettyName+":"+op.Target.FullPrettyVersion(true))
		case solver.OpUninstall:
			uninstalls = append(uninstalls, op.Package.PrettyName)
		}
	}
	if len(installs) == 0 && len(updates) == 0 && len(uninstalls) == 0 {
		io.WriteError("Nothing to install, update or remove")
	} else {
		io.WriteError(fmt.Sprintf("<info>Package operations: %d install%s, %d update%s, %d removal%s</info>", len(installs), plural(len(installs)), len(updates), plural(len(updates)), len(uninstalls), plural(len(uninstalls))))
		if len(installs) > 0 {
			io.WriteErrorV("Installs: "+strings.Join(installs, ", "), true, 64)
		}
		if len(updates) > 0 {
			io.WriteErrorV("Updates: "+strings.Join(updates, ", "), true, 64)
		}
		if len(uninstalls) > 0 {
			io.WriteErrorV("Removals: "+strings.Join(uninstalls, ", "), true, 64)
		}
	}
	if in.ExecuteOperations {
		localRepo.DevPackageNames = c.Locker.DevPackageNames()
		if err := in.IM.Execute(localRepo, ops, in.DevMode, in.RunScripts, in.DownloadOnly); err != nil {
			if se, ok := err.(*ScriptError); ok {
				return se.Code, err
			}
			return 1, err
		}
		if len(ops) > 0 {
			vendor := c.VendorDir()
			if st, err := os.Stat(vendor); err == nil && st.IsDir() {
				now := timeNow()
				_ = os.Chtimes(vendor, now, now)
			}
		}
	} else {
		for _, op := range ops {
			if !strings.Contains(op.Type, "Alias") || io.IsDebug() {
				io.WriteError("  - " + op.Show(false))
			}
		}
	}
	return 0, nil
}

func (in *Installer) missingRequirementInfo() ([]string, error) {
	return MissingRequirementInfo(in.C, in.DevMode)
}

// MissingRequirementInfo is Locker::getMissingRequirementInfo.
func MissingRequirementInfo(c *Composer, includeDev bool) ([]string, error) {
	root := c.Root()
	type set struct {
		repoDev     bool
		links       pkg.Links
		description string
	}
	sets := []set{{false, root.Requires(), "Required"}}
	if includeDev {
		sets = append(sets, set{true, root.DevRequires(), "Required (in require-dev)"})
	}
	var info []string
	missing := false
	for _, s := range sets {
		lr, err := c.Locker.LockedRepository(s.repoDev)
		if err != nil {
			return nil, err
		}
		candidates := append(append([]*pkg.Package(nil), lr.MustPackages()...), root)
		for _, l := range s.links {
			if repo.IsPlatformPackage(l.Target) || l.PrettyConstraint == "self.version" {
				continue
			}
			if len(findWithReplacersAndProviders(candidates, l.Target, l.Constraint)) == 0 {
				results := findWithReplacersAndProviders(candidates, l.Target, nil)
				if len(results) > 0 {
					provider := results[0]
					desc := provider.PrettyVersion()
					if provider.Name != l.Target {
						done := false
						for _, kind := range []struct {
							links pkg.Links
							text  string
						}{{provider.Replaces(), "replaced as %s by %s"}, {provider.Provides(), "provided as %s by %s"}} {
							for _, pl := range kind.links {
								if pl.Target == l.Target {
									desc = fmt.Sprintf(kind.text, pl.PrettyConstraint, provider.PrettyName+" "+provider.PrettyVersion())
									done = true
									break
								}
							}
							if done {
								break
							}
						}
					}
					info = append(info, "- "+s.description+` package "`+l.Target+`" is in the lock file as "`+desc+`" but that does not satisfy your constraint "`+l.PrettyConstraint+`".`)
				} else {
					info = append(info, "- "+s.description+` package "`+l.Target+`" is not present in the lock file.`)
				}
				missing = true
			}
		}
	}
	if missing {
		info = append(info, "This usually happens when composer files are incorrectly merged or the composer.json file is manually edited.",
			"Read more about correctly resolving merge conflicts https://getcomposer.org/doc/articles/resolving-merge-conflicts.md",
			"and prefer using the \"require\" command over editing the composer.json file directly https://getcomposer.org/doc/03-cli.md#require--r")
	}
	return info, nil
}

func findWithReplacersAndProviders(pkgs []*pkg.Package, name string, c semver.Constraint) []*pkg.Package {
	var out []*pkg.Package
	for _, p := range pkgs {
		if p.Name == name {
			if c == nil || c.Matches(semver.NewConstraint("==", p.Version())) {
				out = append(out, p)
			}
			continue
		}
		for _, l := range append(append(pkg.Links(nil), p.Provides()...), p.Replaces()...) {
			if l.Target == name && (c == nil || c.Matches(l.Constraint)) {
				out = append(out, p)
				break
			}
		}
	}
	return out
}

func (in *Installer) outputSuggestions() error {
	c := in.C
	lr, err := c.Locker.LockedRepository(in.DevMode)
	if err != nil {
		return nil
	}
	installed := providers{}
	add := func(p *pkg.Package) {
		for _, n := range p.Names(true) {
			installed[n] = append(installed[n], p.Name)
		}
	}
	for _, p := range lr.MustPackages() {
		add(p)
	}
	if plat, err := in.createPlatformRepo(false); err == nil {
		for _, p := range plat.MustPackages() {
			add(p)
		}
	}
	add(c.Root())
	if f := in.suggestions.filtered(installed); len(f) > 0 {
		c.IO.WriteError(fmt.Sprintf("<info>%d package suggestions were added by new dependencies, use `composer suggest` to see details.</info>", len(f)))
	}
	return nil
}
