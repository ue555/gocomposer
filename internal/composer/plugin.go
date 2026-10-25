package composer

// Composer plugin support: a Go port of Composer
// f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Plugin/PluginManager.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano) together with the plugin
// parts of Factory::createComposer, PluginInstaller and Util/PackageSorter.
// The plugin classes themselves run in the PHP plugin host (pluginhost.go).

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"gocomposer/internal/autoload"
	"gocomposer/internal/config"
	"gocomposer/internal/configsource"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/installed"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// PluginBlockedError is PluginBlockedException.
type PluginBlockedError struct{ Msg string }

func (e *PluginBlockedError) Error() string { return e.Msg }

type allowRule struct {
	re    *regexp.Regexp
	allow bool
}

// PluginManager is Composer's PluginManager for one Composer instance.
type PluginManager struct {
	c        *Composer
	disabled bool
	host     *pluginHost
	// rules is nil in the BC mode for pre-2.2 lock files without allow-plugins
	rules         []allowRule
	rulesBC       bool
	registered    map[string]bool
	registerOrder []string
	autoloaded    map[string]bool
	// devMode is the AutoloadGenerator dev mode used for plugin autoloaders
	devMode bool

	installers      map[string]*pluginInstallerInfo
	installPaths    map[*pkg.Package]string
	curIM           *InstallationManager
	curRepo         *installed.Repository
	curOps          []*solver.Operation
	curOpsKey       string
	opsSeq          int
	eventSeq        int
	dispatchers     []*EventDispatcher
	listenerKeys    map[string]bool
	loaderHash      string
	loaderSeq       int
	runScripts      *bool
	pendingErr      error
	global          *globalComposer
	closed          bool
	lastCallbackErr error

	// Composer instances created by plugins (Factory::create)
	instances map[int64]*Composer
	instSeq   int64
	// repositories handed to the runtime
	repoIDs     map[repo.Repository]int64
	repoByID    map[int64]repo.Repository
	repoConfigs map[repo.Repository]*phpjson.Map
	// outputProgress is InstallationManager::setOutputProgress by a plugin
	outputProgress *bool
	// eventExtras complete the event specs of stateful events by eid
	eventExtras map[string]*phpjson.Map
	// cleanupInstallers are the installers of the operations at download
	// time; pendingOps the operations of plugin installers being settled
	cleanupInstallers map[*pkg.Package]cleanupChoice
	pendingOps        map[string]*pendingOp
	opSeq             int
	// pluginIM is the installation manager of downloader calls made by
	// plugins outside of an install
	pluginIM *InstallationManager
	// txResult/txPresent are the packages of the LocalRepoTransaction of
	// pre-operations-exec
	txResult, txPresent []*pkg.Package
	// services replaced by plugins (plugin_services.go)
	dlReplaced map[string]bool
	alReplaced bool
	alCtx      *autoloadCtx
	// repositories implemented by plugin code (plugin_repos.go)
	phpRepos map[int64]*phpRepository
	// lock transactions of Installer objects of plugin code
	lockTxs map[int64]*solver.LockTransaction
	// configurations of Config objects of plugin code, instances created
	// from a configuration array (plugin_services.go)
	detachedCfgs map[int64]*config.Config
	cfgSeq       int64
	arrayConfigs map[*Composer]bool
	// HTTP clients for the TLS options of plugin requests (plugin_http.go)
	tlsClients map[string]*httpdl.Downloader
	sslWarned  map[string]bool
}

var (
	pluginMu       sync.Mutex
	pluginManagers = map[*Composer]*PluginManager{}
	livePMs        []*PluginManager
)

// pluginManagerOf returns the plugin manager of c (nil when none was created).
func pluginManagerOf(c *Composer) *PluginManager {
	if c == nil {
		return nil
	}
	pluginMu.Lock()
	defer pluginMu.Unlock()
	return pluginManagers[c]
}

// active reports whether plugin code may be running for this manager.
func (pm *PluginManager) active() bool {
	return pm != nil && !pm.disabled && pm.host != nil && !pm.host.dead
}

// ShutdownPlugins stops the plugin hosts of the process (called when the
// command ends).
func ShutdownPlugins() {
	pluginMu.Lock()
	list := livePMs
	livePMs = nil
	pluginMu.Unlock()
	for _, pm := range list {
		pm.close()
	}
}

func (pm *PluginManager) close() {
	if pm.closed {
		return
	}
	pm.closed = true
	if pm.host != nil {
		pm.host.close()
	}
	pluginMu.Lock()
	delete(pluginManagers, pm.c)
	pluginMu.Unlock()
}

// initPlugins is the plugin part of Factory::createComposer: the plugin
// manager is created, the installed plugins are loaded and the init event is
// dispatched for fully loaded instances.
func (c *Composer) initPlugins(fullLoad bool) error {
	// Require/RemoveCommand deactivate the installed plugins before creating
	// the Composer instance used for the update
	pluginMu.Lock()
	prev := append([]*PluginManager(nil), livePMs...)
	pluginMu.Unlock()
	for _, old := range prev {
		if old.c != c && old.active() && !creatingFromPlugin {
			if err := old.deactivateInstalledPlugins(); err != nil {
				return err
			}
		}
	}
	pm := &PluginManager{c: c, disabled: c.DisablePlugins, registered: map[string]bool{}, autoloaded: map[string]bool{}, listenerKeys: map[string]bool{}}
	pm.parseAllowedPlugins()
	// the repositories dispatch the file download events of their metadata
	pm.installRepoHooks(c)
	pluginMu.Lock()
	pluginManagers[c] = pm
	livePMs = append(livePMs, pm)
	pluginMu.Unlock()
	if err := pm.loadInstalledPlugins(); err != nil {
		pm.close()
		return err
	}
	if fullLoad && pm.active() {
		d := NewEventDispatcher(c)
		d.RunScripts = !c.DisableScripts
		if _, err := d.dispatch(&ScriptEvent{Name: "init", Kind: "plain"}); err != nil {
			pm.close()
			return err
		}
	}
	return nil
}

// parseAllowedPlugins is PluginManager::parseAllowedPlugins for the local
// configuration.
func (pm *PluginManager) parseAllowedPlugins() {
	v := pm.c.Config.Get("allow-plugins")
	if isEmptyArray(v) {
		if locked, _ := pm.c.Locker.IsLocked(); locked {
			api := "1.1.0"
			if data, err := pm.c.Locker.LockData(); err == nil {
				if s, ok := mustGet(data, "plugin-api-version").(string); ok {
					api = s
				}
			}
			if semver.VersionCompare(api, "2.2.0") < 0 {
				pm.rulesBC = true
				return
			}
		}
	}
	switch t := v.(type) {
	case bool:
		pm.rules = []allowRule{{regexp.MustCompile(`^`), t}}
		return
	}
	phpjson.Each(v, func(pattern string, allow phpjson.Value) {
		b, _ := allow.(bool)
		pm.rules = append(pm.rules, allowRule{packageNameToRegexp(pattern), b})
	})
}

func isEmptyArray(v phpjson.Value) bool {
	switch t := v.(type) {
	case []phpjson.Value:
		return len(t) == 0
	case *phpjson.Map:
		return t.Len() == 0
	}
	return false
}

// packageNameToRegexp is BasePackage::packageNameToRegexp.
func packageNameToRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(pattern), `\*`, `.*`) + `$`)
}

// isPluginAllowed is PluginManager::isPluginAllowed (local plugins).
func (pm *PluginManager) isPluginAllowed(name string, isGlobal, optional, prompt bool) (bool, error) {
	io := pm.c.IO
	rules := &pm.rules
	if isGlobal && pm.global != nil {
		rules = &pm.global.rules
	}
	if !isGlobal && pm.rulesBC {
		if !io.Interactive {
			return false, &UserError{Msg: `Your composer.lock was generated before the allow-plugins security feature was introduced and your composer.json does not define allow-plugins. Run "composer update --lock" locally and commit the updated composer.lock, then add an explicit allow-plugins section to composer.json. See https://getcomposer.org/allow-plugins`, Code: 1}
		}
		pm.rulesBC = false
		pm.rules = []allowRule{}
	}
	for _, r := range *rules {
		if r.re.MatchString(name) {
			return r.allow, nil
		}
	}
	if name == "composer/package-versions-deprecated" {
		return false, nil
	}
	global := ""
	globalCmd := ""
	if isGlobal || pm.c.isGlobal() {
		global = " (installed globally)"
		globalCmd = "global "
	}
	if io.Interactive && prompt {
		io.WriteError("<warning>" + name + global + " contains a Composer plugin which is currently not in your allow-plugins config. See https://getcomposer.org/allow-plugins</warning>")
		attempts := 0
		for {
			if attempts > 5 {
				io.WriteError("Too many failed prompts, aborting.")
				break
			}
			answer, _, err := io.AskString(`Do you trust "<fg=green;options=bold>`+name+`</>" to execute code and wish to enable it now? (writes "allow-plugins" to composer.json) [<comment>y,n,d,?</comment>] `, "?")
			if err != nil {
				return false, err
			}
			switch answer {
			case "y", "n", "d":
				allow := answer == "y"
				*rules = append([]allowRule{{packageNameToRegexp(name), allow}}, *rules...)
				if (answer == "y" || answer == "n") && !isGlobal {
					if err := pm.persistAllowPlugin(name, allow); err != nil {
						return false, err
					}
				}
				return allow, nil
			default:
				attempts++
				io.WriteError("y - add package to allow-plugins in composer.json and let it run immediately")
				io.WriteError("n - add package (as disallowed) to allow-plugins in composer.json to suppress further prompts")
				io.WriteError("d - discard this, do not change composer.json and do not allow the plugin to run")
				io.WriteError("? - print help")
			}
		}
	} else if optional {
		return false, nil
	}
	return false, &PluginBlockedError{Msg: name + global + " contains a Composer plugin which is blocked by your allow-plugins config. You may add it to the list if you consider it safe.\n" +
		`You can run "composer ` + globalCmd + `config --no-plugins allow-plugins.` + name + ` [true|false]" to enable it (true) or disable it explicitly and suppress this exception (false)` + "\n" +
		"See https://getcomposer.org/allow-plugins"}
}

// persistAllowPlugin writes the prompt answer to composer.json like
// isPluginAllowed does through the config source.
func (pm *PluginManager) persistAllowPlugin(name string, allow bool) error {
	cur := pm.c.Config.Get("allow-plugins")
	m := phpjson.NewMap()
	switch t := cur.(type) {
	case *phpjson.Map:
		m = t.Clone()
	case []phpjson.Value:
	default:
		return nil
	}
	m.Set(name, allow)
	if pm.c.Config.GetBool("sort-packages") {
		m.KSort()
	}
	if err := configsource.New(pm.c.ComposerFile, false).AddConfigSetting("allow-plugins", m); err != nil {
		return err
	}
	wrap := phpjson.NewMap()
	cfg := phpjson.NewMap()
	cfg.Set("allow-plugins", m)
	wrap.Set("config", cfg)
	pm.c.Config.Merge(wrap)
	return nil
}

func pluginOptional(p *pkg.Package) bool {
	if e := p.Extra(); e != nil {
		if v, ok := e.Get("plugin-optional"); ok && v == true {
			return true
		}
	}
	return false
}

func isPluginType(t string) bool { return t == "composer-plugin" || t == "composer-installer" }

// loadInstalledPlugins is PluginManager::loadInstalledPlugins for the
// project's local repository (loadRepository).
func (pm *PluginManager) loadInstalledPlugins() error {
	if pm.disabled {
		return nil
	}
	pm.global = pm.loadGlobalComposer()
	if packages, err := pm.c.LocalRepo.Packages(); err == nil {
		if err := pm.loadRepository(packages, false); err != nil {
			return err
		}
	}
	if pm.global != nil {
		if packages, err := pm.global.repo.Packages(); err == nil {
			if err := pm.loadRepository(packages, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadRepository is PluginManager::loadRepository.
func (pm *PluginManager) loadRepository(packages []*pkg.Package, isGlobal bool) error {
	hasPlugin := false
	for _, p := range packages {
		if p.AliasOf == nil && isPluginType(p.Type()) {
			hasPlugin = true
		}
	}
	if !hasPlugin {
		return nil
	}
	weights := map[string]int{}
	for _, p := range packages {
		if p.Type() == "composer-plugin" {
			if p.Name == "composer/installers" {
				weights[p.Name] = -10000
			} else if e := p.Extra(); e != nil {
				if v, ok := e.Get("plugin-modifies-install-path"); ok && v == true {
					weights[p.Name] = -10000
				}
			}
		}
	}
	sorted := sortPackagesWeighted(packages, weights)
	isRequired := map[*pkg.Package]bool{}
	if !isGlobal {
		for _, p := range filterRequiredPackages(packages, pm.c.Root(), true, nil) {
			isRequired[p] = true
		}
	}
	for _, p := range sorted {
		if p.AliasOf != nil || !isPluginType(p.Type()) {
			continue
		}
		if !isGlobal && !isRequired[p] {
			ok, err := pm.isPluginAllowed(p.Name, false, true, false)
			if err != nil {
				return err
			}
			if !ok {
				pm.c.IO.WriteError(`<warning>The "` + p.Name + `" plugin was not loaded as it is not listed in allow-plugins and is not required by the root package anymore.</warning>`)
				continue
			}
		}
		if err := pm.registerPackage(p, false, isGlobal); err != nil {
			return err
		}
	}
	return nil
}

// sortPackagesWeighted is PackageSorter::sortPackages.
func sortPackagesWeighted(packages []*pkg.Package, weights map[string]int) []*pkg.Package {
	usage := map[string][]string{}
	for _, p := range packages {
		links := p.Requires()
		if p.IsRoot {
			links = append(append(pkg.Links(nil), links...), p.DevRequires()...)
		}
		for _, l := range links {
			usage[l.Target] = append(usage[l.Target], p.Name)
		}
	}
	computing := map[string]bool{}
	computed := map[string]int{}
	var importance func(name string) int
	importance = func(name string) int {
		if w, ok := computed[name]; ok {
			return w
		}
		if computing[name] {
			return 0
		}
		computing[name] = true
		w := weights[name]
		for _, user := range usage[name] {
			w -= 1 - importance(user)
		}
		delete(computing, name)
		computed[name] = w
		return w
	}
	type weighted struct {
		name   string
		weight int
		index  int
	}
	list := make([]weighted, len(packages))
	for i, p := range packages {
		list[i] = weighted{p.Name, importance(p.Name), i}
	}
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].weight != list[j].weight {
			return list[i].weight < list[j].weight
		}
		return semver.NatCompare(strings.ToLower(list[i].name), strings.ToLower(list[j].name)) < 0
	})
	out := make([]*pkg.Package, len(packages))
	for i, w := range list {
		out[i] = packages[w.index]
	}
	return out
}

// filterRequiredPackages is RepositoryUtils::filterRequiredPackages.
func filterRequiredPackages(packages []*pkg.Package, requirer *pkg.Package, includeDev bool, bucket []*pkg.Package) []*pkg.Package {
	requires := map[string]bool{}
	for _, l := range requirer.Requires() {
		requires[l.Target] = true
	}
	if includeDev {
		for _, l := range requirer.DevRequires() {
			requires[l.Target] = true
		}
	}
	for _, candidate := range packages {
		for _, name := range candidate.Names(true) {
			if requires[name] {
				found := false
				for _, b := range bucket {
					if b == candidate {
						found = true
						break
					}
				}
				if !found {
					bucket = append(bucket, candidate)
					bucket = filterRequiredPackages(packages, candidate, false, bucket)
				}
				break
			}
		}
	}
	return bucket
}

// registerPackage is PluginManager::registerPackage.
func (pm *PluginManager) registerPackage(p *pkg.Package, failOnMissing, isGlobal bool) error {
	io := pm.c.IO
	if pm.disabled {
		io.WriteError(`<warning>The "` + p.Name + `" plugin was not loaded as plugins are disabled.</warning>`)
		return nil
	}
	globalNote := ""
	if isGlobal || pm.c.isGlobal() {
		globalNote = "(installed globally) "
	}
	if p.Type() == "composer-plugin" {
		var req *pkg.Link
		for _, l := range p.Requires() {
			if l.Target == "composer-plugin-api" {
				req = l
				break
			}
		}
		if req == nil {
			return &UserError{Msg: "Plugin " + p.Name + " is missing a require statement for a version of the composer-plugin-api package.", Code: 1}
		}
		if req.PrettyConstraint == platform.PluginAPIVersion {
			io.WriteError(`<warning>The "` + p.Name + `" plugin requires composer-plugin-api ` + platform.PluginAPIVersion + `, this *WILL* break in the future and it should be fixed ASAP (require ^` + platform.PluginAPIVersion + ` instead for example).</warning>`)
		} else if req.Constraint != nil && !req.Constraint.Matches(semver.NewConstraint("==", platform.PluginAPIVersion+".0")) {
			io.WriteError(`<warning>The "` + p.Name + `" plugin ` + globalNote + `was skipped because it requires a Plugin API version ("` + req.PrettyConstraint + `") that does not match your Composer installation ("` + platform.PluginAPIVersion + `"). You may need to run composer update with the "--no-plugins" option.</warning>`)
			return nil
		}
		if p.Name == "symfony/flex" && regexp.MustCompile(`^[0-9.]+$`).MatchString(p.Version()) && semver.VersionCompare(p.Version(), "1.9.8") < 0 {
			io.WriteError(`<warning>The "` + p.Name + `" plugin ` + globalNote + `was skipped because it is not compatible with Composer 2+. Make sure to update it to version 1.9.8 or greater.</warning>`)
			return nil
		}
	}
	ok, err := pm.isPluginAllowed(p.Name, isGlobal, pluginOptional(p), true)
	if err != nil {
		return err
	}
	if !ok {
		io.WriteErrorV(`Skipped loading "`+p.Name+`" `+globalNote+`as it is not in config.allow-plugins`, true, 256)
		return nil
	}
	legacy := p.Type() == "composer-installer"
	if pm.registered[p.Name] {
		return nil
	}
	classes := pluginClasses(p)
	if len(classes) == 0 {
		return &UserError{Msg: "Error while installing " + p.PrettyName + ", composer-plugin packages should have a class defined in their extra key to be usable.", Code: 1}
	}
	al, err := pm.pluginAutoload(p, isGlobal)
	if err != nil {
		return err
	}
	if err := pm.ensureHost(); err != nil {
		return err
	}
	h := pm.host
	args := phpjson.NewMap()
	args.Set("pkg", h.ref(p))
	args.Set("global", isGlobal || pm.c.isGlobal())
	args.Set("legacy", legacy)
	args.Set("failOnMissing", failOnMissing)
	args.Set("vendorDir", pm.c.VendorDir())
	cls := make([]phpjson.Value, len(classes))
	for i, s := range classes {
		cls[i] = s
	}
	args.Set("classes", cls)
	nsList := func(list []autoload.NamespacePaths) []phpjson.Value {
		out := []phpjson.Value{}
		for _, e := range list {
			paths := make([]phpjson.Value, len(e.Paths))
			for i, s := range e.Paths {
				paths[i] = s
			}
			out = append(out, []phpjson.Value{e.Namespace, paths})
		}
		return out
	}
	args.Set("psr0", nsList(al.PSR0))
	args.Set("psr4", nsList(al.PSR4))
	cm := phpjson.NewMap()
	for _, e := range al.ClassMap {
		cm.Set(e[0], e[1])
	}
	args.Set("classmap", cm)
	files := []phpjson.Value{}
	for _, f := range al.Files {
		files = append(files, []phpjson.Value{f[0], f[1]})
	}
	args.Set("files", files)
	args.Set("frames", pm.frames())
	if _, err := h.call("register", args); err != nil {
		return err
	}
	pm.registered[p.Name] = true
	pm.registerOrder = append(pm.registerOrder, p.Name)
	pm.installers = nil
	pm.installPaths = nil
	return nil
}

func pluginClasses(p *pkg.Package) []string {
	e := p.Extra()
	if e == nil {
		return nil
	}
	v, ok := e.Get("class")
	if !ok || !phpjson.Truthy(v) {
		return nil
	}
	var out []string
	if s, ok := v.(string); ok {
		return []string{s}
	}
	phpjson.Each(v, func(_ string, c phpjson.Value) {
		if s, ok := phpjson.ScalarString(c); ok {
			out = append(out, s)
		}
	})
	return out
}

// pluginAutoload builds the plugin's runtime autoloader like registerPackage:
// the plugin, its dependencies from the local repository and the root
// package (without "files" autoloads), each package autoloaded only once.
func (pm *PluginManager) pluginAutoload(p *pkg.Package, isGlobal bool) (*autoload.PluginAutoload, error) {
	c := pm.c
	root := c.Root().Clone()
	clearFiles := func(m *phpjson.Map) *phpjson.Map {
		if m == nil {
			m = phpjson.NewMap()
		} else {
			m = m.Clone()
		}
		m.Set("files", []phpjson.Value{})
		return m
	}
	root.SetAutoload(clearFiles(c.Root().Autoload()))
	root.SetDevAutoload(clearFiles(c.Root().DevAutoload()))
	local := c.LocalRepo.MustPackages()
	repoPkgs := append(append([]*pkg.Package(nil), local...), root)
	if pm.global != nil {
		repoPkgs = append(repoPkgs, pm.global.repo.MustPackages()...)
	}
	collected := map[string]*pkg.Package{p.Name: p}
	order := []string{p.Name}
	var collect func(q *pkg.Package)
	collect = func(q *pkg.Package) {
		for _, l := range q.Requires() {
			for _, dep := range findWithReplacersAndProviders(repoPkgs, l.Target, nil) {
				if _, ok := collected[dep.Name]; !ok {
					collected[dep.Name] = dep
					order = append(order, dep.Name)
					collect(dep)
				}
			}
		}
	}
	collect(p)
	im := NewInstallationManager(c)
	var entries []autoload.PackageEntry
	should := func(q *pkg.Package, path string) bool {
		key := q.Name + ":" + q.Version() + ":" + path
		if pm.autoloaded[key] {
			return false
		}
		pm.autoloaded[key] = true
		return true
	}
	if should(root, "") {
		empty := ""
		entries = append(entries, autoload.PackageEntry{Pkg: root, InstallPath: &empty})
	}
	for _, name := range order {
		q := collected[name]
		if q == root {
			continue
		}
		if q.Type() == "metapackage" {
			continue
		}
		path := im.InstallPath(q)
		if pm.global.has(q) {
			path = pm.global.installPath(q)
		}
		if should(q, path) || q == p {
			pp := path
			entries = append(entries, autoload.PackageEntry{Pkg: q, InstallPath: &pp})
		}
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("At least the plugin package should always be autoloaded for the code below to work")
	}
	return autoload.PluginLoader(root, entries, c.Cwd, c.VendorDir(), pm.devMode, func(s string) { c.IO.WriteError(s) }), nil
}

func (pm *PluginManager) ensureHost() error {
	if pm.host != nil {
		if pm.host.dead {
			return &ScriptError{Msg: "", Code: 255}
		}
		return nil
	}
	h, err := startPluginHost(pm)
	if err != nil {
		return err
	}
	pm.host = h
	args := phpjson.NewMap()
	args.Set("verbosity", int64(pm.c.IO.Verbosity))
	args.Set("decorated", pm.c.IO.Decorated)
	args.Set("interactive", pm.c.IO.Interactive)
	args.Set("tmpDir", h.dir)
	args.Set("global", pm.c.isGlobal())
	args.Set("root", pm.rootDump())
	args.Set("app", appSpec())
	h.rootJSON, _ = phpjson.Encode(pm.rootDump(), 0)
	if _, err := h.call("init", args); err != nil {
		return err
	}
	return nil
}

// rootDump is the root package data sent to the plugin host.
func (pm *PluginManager) rootDump() *phpjson.Map { return rootDumpOf(pm.c) }

// rootDumpOf is the root package data of c for the plugin runtime.
func rootDumpOf(c *Composer) *phpjson.Map {
	root := c.Root()
	d := pkg.Dump(root)
	if c.Package != nil && c.Package.AliasOf != nil {
		d = pkg.Dump(c.Package.Base())
	}
	extra := phpjson.NewMap()
	extra.Set("minimum-stability", root.MinimumStability)
	extra.Set("prefer-stable", root.PreferStable)
	flags := phpjson.NewMap()
	var names []string
	for n := range root.StabilityFlags {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		flags.Set(n, int64(root.StabilityFlags[n]))
	}
	extra.Set("stability-flags", flags)
	refs := phpjson.NewMap()
	names = names[:0]
	for n := range root.References {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		refs.Set(n, root.References[n])
	}
	extra.Set("references", refs)
	aliases := []phpjson.Value{}
	for _, a := range root.Aliases {
		am := phpjson.NewMap()
		am.Set("package", a.Package)
		am.Set("version", a.Version)
		am.Set("alias", a.Alias)
		am.Set("alias_normalized", a.AliasNormalized)
		aliases = append(aliases, am)
	}
	extra.Set("aliases", aliases)
	if root.RootConfig != nil {
		extra.Set("config", root.RootConfig)
	} else {
		extra.Set("config", phpjson.NewMap())
	}
	d.Set("_root", extra)
	return d
}

// applyRootChanges takes over root package modifications made by plugins
// (e.g. autoload rules added in pre-autoload-dump listeners).
func (pm *PluginManager) applyRootChanges(m *phpjson.Map) {
	if m == nil {
		return
	}
	root := pm.c.Root()
	tmp, err := (&pkg.Loader{}).Load(m)
	if err != nil {
		pm.pendingErr = err
		return
	}
	root.SetAutoload(tmp.Autoload())
	root.SetDevAutoload(tmp.DevAutoload())
	root.SetExtra(tmp.Extra())
	root.SetScripts(tmp.Scripts())
	for _, kind := range pkg.LinkKinds {
		links := tmp.LinksOf(kind.Type)
		for _, l := range links {
			l.Source = root.Name
		}
		root.SetLinksOf(kind.Type, links)
	}
	if pm.host != nil {
		pm.host.rootJSON, _ = phpjson.Encode(pm.rootDump(), 0)
	}
	pm.installPaths = nil
}

// deactivateInstalledPlugins is PluginManager::deactivateInstalledPlugins.
func (pm *PluginManager) deactivateInstalledPlugins() error {
	if !pm.active() {
		return nil
	}
	repos := []*installed.Repository{pm.c.LocalRepo}
	if pm.global != nil {
		repos = append(repos, pm.global.repo)
	}
	for _, r := range repos {
		packages, err := r.Packages()
		if err != nil {
			continue
		}
		sorted := sortPackagesWeighted(packages, nil)
		for i := len(sorted) - 1; i >= 0; i-- {
			p := sorted[i]
			if p.AliasOf != nil || !isPluginType(p.Type()) {
				continue
			}
			if err := pm.deactivatePackage(p); err != nil {
				return err
			}
		}
	}
	return nil
}

// deactivatePackage is PluginManager::deactivatePackage.
func (pm *PluginManager) deactivatePackage(p *pkg.Package) error {
	if !pm.active() || !pm.registered[p.Name] {
		return nil
	}
	args := phpjson.NewMap()
	args.Set("name", p.Name)
	if _, err := pm.host.call("deactivatePackage", args); err != nil {
		return err
	}
	pm.forget(p.Name)
	pm.installers = nil
	pm.installPaths = nil
	return nil
}

// uninstallPackage is PluginManager::uninstallPackage.
func (pm *PluginManager) uninstallPackage(p *pkg.Package) error {
	if !pm.active() || !pm.registered[p.Name] {
		return nil
	}
	args := phpjson.NewMap()
	args.Set("name", p.Name)
	if _, err := pm.host.call("uninstallPackage", args); err != nil {
		return err
	}
	pm.forget(p.Name)
	pm.installers = nil
	pm.installPaths = nil
	return nil
}

func (pm *PluginManager) forget(name string) {
	delete(pm.registered, name)
	for i, n := range pm.registerOrder {
		if n == name {
			pm.registerOrder = append(pm.registerOrder[:i], pm.registerOrder[i+1:]...)
			break
		}
	}
}

// RegisteredPlugins is PluginManager::getRegisteredPlugins (diagnose).
func RegisteredPlugins(c *Composer) []string {
	pm := pluginManagerOf(c)
	if pm == nil {
		return nil
	}
	var out []string
	for _, n := range pm.registerOrder {
		if pm.registered[n] {
			out = append(out, n)
		}
	}
	return out
}
