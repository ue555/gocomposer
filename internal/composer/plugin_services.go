package composer

// Composer services replaced or extended by plugins, bridged to the PHP
// objects of the plugin runtime (internal/res/pluginrt/services.php):
// downloaders replacing gocomposer's built-in ones
// (DownloadManager::setDownloader) and the downloader internals subclasses
// call (getFileName, extract, getLocalChanges, getVcsReference). gocomposer
// keeps doing the work (transfers, extraction, file moves); the plugin's
// PHP object decides what is called, like Composer's DownloadManager calling
// the registered downloader. Semantics follow Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/Downloader/*.php,
// MIT license, (c) Nils Adermann, Jordi Boggiano).

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/autoload"
	"gocomposer/internal/config"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/installed"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

// builtinDownloaderTypes are the downloaders Factory::createDownloadManager
// registers.
var builtinDownloaderTypes = map[string]bool{
	"git": true, "svn": true, "fossil": true, "hg": true, "perforce": true,
	"zip": true, "rar": true, "tar": true, "gzip": true, "xz": true,
	"phar": true, "file": true, "path": true,
}

// servicesCallback serves the callbacks of services.php (ok is false for
// other ops).
func (pm *PluginManager) servicesCallback(op string, a *phpjson.Map) (phpjson.Value, bool, error) {
	if v, ok, err := pm.reposCallback(op, a); ok {
		return v, true, err
	}
	if v, ok, err := pm.configCallback(op, a); ok {
		return v, true, err
	}
	if v, ok, err := pm.archiveCallback(op, a); ok {
		return v, true, err
	}
	switch op {
	case "dlReplaced":
		// DownloadManager::setDownloader replacing a built-in downloader
		// of the Composer instance running the plugins
		if n, _ := mustGet(a, "inst").(int64); n == 0 {
			if pm.dlReplaced == nil {
				pm.dlReplaced = map[string]bool{}
			}
			pm.dlReplaced[strings.ToLower(a.Str("type"))] = true
		}
		return nil, true, nil
	case "dlFileName":
		p, err := pm.host.pkgFromRef(mustGet(a, "p"))
		if err != nil {
			return nil, true, err
		}
		return pm.dlIM().tmpFileName(p), true, nil
	case "dlExtract":
		async, err := pm.dlExtract(a)
		if err != nil {
			return nil, true, err
		}
		m := phpjson.NewMap()
		m.Set("async", async)
		return m, true, nil
	case "dlArchive":
		v, err := pm.dlArchive(a)
		return v, true, err
	case "dlLocalChanges":
		v, err := pm.dlLocalChanges(a)
		return v, true, err
	case "serviceSet":
		// Composer::set*() of the instance running the plugins
		if n, _ := mustGet(a, "inst").(int64); n == 0 {
			switch a.Str("name") {
			case "autoloadGenerator":
				pm.alReplaced = true
			}
		}
		return nil, true, nil
	case "nestedDispatch", "nestedHasListeners":
		// the event dispatcher of a Composer instance created with
		// Factory::create()
		c, err := pm.instance(a)
		if err != nil {
			return nil, true, err
		}
		d := NewEventDispatcher(c)
		d.RunScripts = mustGet(a, "runScripts") != false
		if op == "nestedHasListeners" {
			return d.HasListeners(a.Str("name")), true, nil
		}
		kind := a.Str("kind")
		if kind == "unsupported" {
			return nil, true, &UserError{Msg: "Dispatching package or installer events from plugins is not supported by gocomposer. Run the command with --no-plugins to skip plugins.", Code: 1}
		}
		ev := &ScriptEvent{Name: a.Str("name"), Kind: kind, DevMode: mustGet(a, "dev") == true, Command: a.Str("command")}
		phpjson.Each(mustGet(a, "args"), func(_ string, v phpjson.Value) {
			s, _ := phpjson.ScalarString(v)
			ev.Args = append(ev.Args, s)
		})
		if f := phpjson.ToMap(mustGet(a, "flags")); f != nil && f.Len() > 0 {
			ev.Flags = map[string]any{}
			for _, k := range f.Keys() {
				ev.Flags[k] = mustGet(f, k)
			}
		}
		code, err := d.dispatch(ev)
		if err != nil {
			return nil, true, err
		}
		return int64(code), true, nil
	case "replacedPost":
		return nil, true, pm.replacedPost(a.Str("key"))
	case "dlVcsReference":
		p, err := pm.host.pkgFromRef(mustGet(a, "p"))
		if err != nil {
			return nil, true, err
		}
		_, commit, ok := GuessVersionCommit(strings.TrimRight(a.Str("path"), "/"), p, pm.c)
		if !ok || commit == "" {
			return nil, true, nil
		}
		return commit, true, nil
	}
	return nil, false, nil
}

// dlIM is the installation manager of downloader calls made by plugins: the
// running one, or one kept for the calls outside of an install (download,
// then install use the same temporary file).
func (pm *PluginManager) dlIM() *InstallationManager {
	im := pm.curIM
	if im == nil || im.c != pm.c {
		if pm.pluginIM == nil {
			pm.pluginIM = NewInstallationManager(pm.c)
		}
		im = pm.pluginIM
		if pm.outputProgress != nil {
			im.OutputProgress = *pm.outputProgress
		}
	}
	return im
}

// replacedDownloader reports whether a plugin replaced gocomposer's built-in
// downloader typ (the plugin's PHP downloader then performs the operations
// of that type, calling back into gocomposer for what it inherits).
func (im *InstallationManager) replacedDownloader(typ string) bool {
	if typ == "" || !builtinDownloaderTypes[typ] {
		return false
	}
	pm := pluginManagerOf(im.c)
	return pm.active() && pm.dlReplaced[typ]
}

// replacedDownloaderOp runs DownloadManager's call of a replaced downloader
// (handled is false when typ was not replaced).
func (im *InstallationManager) replacedDownloaderOp(typ, method string, p, prev *pkg.Package) (bool, error) {
	if !im.replacedDownloader(typ) {
		return false, nil
	}
	return im.pluginDownloaderCall(typ, method, p, prev, "")
}

// pluginDownloaderCall calls method of the downloader registered for typ in
// the plugin runtime (opType: the operation of prepare/cleanup).
func (im *InstallationManager) pluginDownloaderCall(typ, method string, p, prev *pkg.Package, opType string) (bool, error) {
	pm := pluginManagerOf(im.c)
	args := phpjson.NewMap()
	args.Set("type", typ)
	args.Set("m", method)
	args.Set("p", pm.host.ref(p))
	if prev != nil {
		args.Set("i", pm.host.ref(prev))
	}
	if opType != "" {
		args.Set("op", opType)
	}
	args.Set("path", im.InstallPath(p))
	args.Set("frames", pm.frames())
	_, err := pm.host.call("downloaderOp", args)
	if err == nil {
		err = pm.takePendingErr()
	}
	return true, err
}

// prepareStage is DownloadManager::prepare for an operation: a replaced
// downloader's prepare(), else gocomposer's own (prepareOp).
func (im *InstallationManager) prepareStage(op *solver.Operation) error {
	if handled, err := im.replacedStage(op, "prepare"); handled {
		return err
	}
	return im.prepareOp(op)
}

// cleanupStage is DownloadManager::cleanup for an operation.
func (im *InstallationManager) cleanupStage(op *solver.Operation) error {
	if handled, err := im.replacedStage(op, "cleanup"); handled {
		return err
	}
	return im.cleanupOp(op)
}

func (im *InstallationManager) replacedStage(op *solver.Operation, stage string) (bool, error) {
	var p, prev *pkg.Package
	switch op.Type {
	case solver.OpInstall, solver.OpUninstall:
		p = op.Package
	case solver.OpUpdate:
		p, prev = op.Target, op.Initial
	default:
		return false, nil
	}
	if p.Type() == "metapackage" {
		return false, nil
	}
	typ := im.downloaderType(p)
	if !im.replacedDownloader(typ) {
		return false, nil
	}
	return im.pluginDownloaderCall(typ, stage, p, prev, op.Type)
}

// dlExtract is ArchiveDownloader::extract of the built-in archive downloader
// kind (called by a plugin's subclass overriding extract()): file is
// extracted into path. async reports whether Composer extracts in a
// background process (unzip), whose promise settles later.
func (pm *PluginManager) dlExtract(a *phpjson.Map) (bool, error) {
	p, err := pm.host.pkgFromRef(mustGet(a, "p"))
	if err != nil {
		return false, err
	}
	kind := a.Str("kind")
	file, path := a.Str("file"), a.Str("path")
	if err := os.MkdirAll(path, 0o777); err != nil {
		return false, err
	}
	im := pm.dlIM()
	switch kind {
	case "zip":
		// unzip (native), 7z or ZipArchive like ZipDownloader::extract
		if err := im.extractZipWith(p, file, path); err != nil {
			return false, err
		}
		return unzipAvailable(), nil
	case "xz":
		return false, im.extractXz(file, path)
	case "rar":
		return false, im.extractRar(file, path)
	case "tar":
		if err := ExtractTar(file, path); err != nil {
			return false, &UserError{Msg: "Failed to extract " + p.Name + ": " + err.Error(), Code: 1}
		}
		return false, nil
	case "gzip":
		name := distPathComponent(p, false)
		if i := strings.LastIndexByte(name, '.'); i >= 0 {
			name = name[:i]
		}
		return false, ExtractGzip(file, filepath.Join(path, name))
	}
	return false, &UserError{Msg: "The " + kind + " archive format is not supported by gocomposer", Code: 1}
}

// dlArchive performs the file system steps of ArchiveDownloader::install
// around the extract() of a plugin's downloader subclass: "begin" empties
// the target and creates the temporary directory, "finish" moves the
// extracted files into place, "fail" cleans up.
func (pm *PluginManager) dlArchive(a *phpjson.Map) (phpjson.Value, error) {
	im := pm.dlIM()
	p, err := pm.host.pkgFromRef(mustGet(a, "p"))
	if err != nil {
		return nil, err
	}
	target := strings.TrimRight(a.Str("path"), "/")
	vendorDir := im.realVendorDir()
	switch a.Str("m") {
	case "begin":
		if !strings.Contains(fsutil.NormalizePath(vendorDir), fsutil.NormalizePath(target+"/")) {
			if err := emptyDirectory(target); err != nil {
				return nil, err
			}
		}
		var tempDir string
		for {
			var rnd [4]byte
			_, _ = rand.Read(rnd[:])
			tempDir = vendorDir + "/composer/" + hex.EncodeToString(rnd[:])
			if _, err := os.Stat(tempDir); os.IsNotExist(err) {
				break
			}
		}
		if err := os.MkdirAll(tempDir, 0o777); err != nil {
			return nil, err
		}
		out := phpjson.NewMap()
		out.Set("tmp", tempDir)
		out.Set("file", im.tmpFileName(p))
		return out, nil
	case "fail":
		im.clearLastCacheWrite(p)
		os.RemoveAll(a.Str("tmp"))
		if st, err := os.Stat(target); err == nil && st.IsDir() && fsutil.Realpath(target) != im.c.Cwd {
			os.RemoveAll(target)
		}
		return nil, nil
	case "finish":
		tempDir := a.Str("tmp")
		os.Remove(im.tmpFileName(p))
		entries := folderContent(tempDir)
		renameAsOne := false
		if _, err := os.Lstat(target); os.IsNotExist(err) {
			renameAsOne = true
		} else if fsutil.IsDirEmpty(target) {
			if os.Remove(target) == nil {
				renameAsOne = true
			}
		}
		from := tempDir
		if len(entries) == 1 && isDir(entries[0]) {
			from = entries[0]
		}
		if renameAsOne {
			if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
				return nil, err
			}
			if err := os.Rename(from, target); err != nil {
				return nil, err
			}
		} else if err := renameRecursively(p, from, target); err != nil {
			return nil, err
		}
		os.RemoveAll(tempDir)
		return nil, nil
	}
	return nil, &UserError{Msg: "Unknown archive step " + a.Str("m"), Code: 1}
}

// dlLocalChanges is ChangeReportInterface::getLocalChanges of a built-in
// downloader kind (nil when nothing changed).
func (pm *PluginManager) dlLocalChanges(a *phpjson.Map) (phpjson.Value, error) {
	im := pm.dlIM()
	p, err := pm.host.pkgFromRef(mustGet(a, "p"))
	if err != nil {
		return nil, err
	}
	path := strings.TrimRight(a.Str("path"), "/")
	var out string
	switch kind := a.Str("kind"); kind {
	case "zip", "tar", "gzip", "xz", "rar", "file", "phar", "path":
		out, err = im.FileLocalChanges(p, path)
	default:
		if !isVcsDownloaderType(kind) {
			return nil, &UserError{Msg: "gocomposer cannot detect local changes with the " + kind + " downloader", Code: 1}
		}
		out, err = im.sourceDL(kind).LocalChanges(p, path)
	}
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return out, nil
}

// removeCodeDM is DownloadManager::remove: a replaced downloader's remove(),
// else gocomposer's.
func (im *InstallationManager) removeCodeDM(p *pkg.Package) error {
	if handled, err := im.replacedDownloaderOp(im.downloaderType(p), "remove", p, nil); handled {
		return err
	}
	return im.removeCode(p, true)
}

// replacedStartOp starts an operation whose package is handled by a
// downloader replaced by a plugin like InstallationManager::executeBatch
// does with LibraryInstaller: the install/update/uninstall promise of the
// DownloadManager (the plugin's downloader), then the binaries and the
// local repository, the cleanup of the downloader and the repository
// write, chained in the plugin runtime and settled with the batch by
// pluginSettleOps (the downloader's background steps complete like
// Composer's ProcessExecutor jobs). started is false for other operations.
func (im *InstallationManager) replacedStartOp(repo *installed.Repository, op *solver.Operation, devMode bool) (bool, error) {
	pm := pluginManagerOf(im.c)
	if !pm.active() || len(pm.dlReplaced) == 0 {
		return false, nil
	}
	var p, initial *pkg.Package
	switch op.Type {
	case solver.OpInstall, solver.OpUninstall:
		p = op.Package
	case solver.OpUpdate:
		p, initial = op.Target, op.Initial
	default:
		return false, nil
	}
	// metapackages and plugins have their own installers
	for _, q := range []*pkg.Package{p, initial} {
		if q == nil {
			continue
		}
		switch q.Type() {
		case "metapackage", "composer-plugin", "composer-installer":
			return false, nil
		}
	}
	if !im.replacedDownloader(im.downloaderType(p)) && (initial == nil || !im.replacedDownloader(im.downloaderType(initial))) {
		return false, nil
	}
	// LibraryInstaller's steps before the code operation
	switch op.Type {
	case solver.OpInstall:
		if target := im.InstallPath(p); target != "" {
			if _, err := os.Stat(target); err != nil && repo.HasPackage(p) {
				im.removeBinaries(p)
			}
		}
	case solver.OpUpdate:
		if !repo.HasPackage(initial) {
			return true, &UserError{Msg: "Package is not installed: " + initial.String(), Code: 1}
		}
		im.removeBinaries(initial)
		initialPath, targetPath := im.InstallPath(initial), im.InstallPath(p)
		if initialPath != targetPath {
			if strings.HasPrefix(initialPath, targetPath) || strings.HasPrefix(targetPath, initialPath) {
				// LibraryInstaller::updateCode's remove + install
				return false, nil
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o777); err != nil {
				return true, err
			}
			if err := os.Rename(initialPath, targetPath); err != nil {
				return true, err
			}
		}
	case solver.OpUninstall:
		if !repo.HasPackage(p) {
			return true, &UserError{Msg: "Package is not installed: " + p.String(), Code: 1}
		}
	}
	pm.opSeq++
	key := "op" + itoa(pm.opSeq)
	if pm.pendingOps == nil {
		pm.pendingOps = map[string]*pendingOp{}
	}
	pm.pendingOps[key] = &pendingOp{im: im, repo: repo, op: op, devMode: devMode, replaced: true}
	args := phpjson.NewMap()
	args.Set("m", op.Type)
	args.Set("p", pm.host.ref(p))
	args.Set("psrc", installationSourceValue(p))
	if initial != nil {
		args.Set("i", pm.host.ref(initial))
		args.Set("isrc", installationSourceValue(initial))
	}
	args.Set("path", im.InstallPath(p))
	args.Set("key", key)
	args.Set("frames", pm.frames())
	prevIM, prevRepo := pm.curIM, pm.curRepo
	pm.curIM, pm.curRepo = im, repo
	_, err := pm.host.call("replacedOp", args)
	pm.curIM, pm.curRepo = prevIM, prevRepo
	if err == nil {
		err = pm.takePendingErr()
	}
	return true, err
}

// replacedPost is LibraryInstaller's work once the code operation of a
// replaced downloader completed: binaries and the local repository.
func (pm *PluginManager) replacedPost(key string) error {
	po := pm.pendingOps[key]
	if po == nil {
		return nil
	}
	im, repo, op := po.im, po.repo, po.op
	switch op.Type {
	case solver.OpInstall:
		if target := im.InstallPath(op.Package); target != "" {
			im.installBinaries(op.Package, target, true)
		}
		if !repo.HasPackage(op.Package) {
			repo.AddPackage(op.Package.Clone())
		}
	case solver.OpUpdate:
		im.installBinaries(op.Target, im.InstallPath(op.Target), true)
		repo.RemovePackage(op.Initial)
		if !repo.HasPackage(op.Target) {
			repo.AddPackage(op.Target.Clone())
		}
	case solver.OpUninstall:
		p := op.Package
		im.removeBinaries(p)
		repo.RemovePackage(p)
		if strings.Contains(p.Name, "/") {
			vendor := filepath.Dir(im.packageBasePath(p))
			if isDir(vendor) && fsutil.IsDirEmpty(vendor) {
				os.Remove(vendor)
			}
		}
	}
	return nil
}

// autoloadCtx is a dump gocomposer handed to an autoload generator replaced
// by a plugin: the generator's dump() calling parent::dump() comes back to
// gocomposer, which then dumps with this context.
type autoloadCtx struct {
	c         *Composer
	events    *EventDispatcher
	localRepo *installed.Repository
	o         AutoloadOptions
	result    *autoload.Result
}

// pluginAutoloadDump runs AutoloadGenerator::dump through the generator a
// plugin set with Composer::setAutoloadGenerator(), with the setters
// Composer's callers use (Installer::run, dump-autoload, reinstall).
func pluginAutoloadDump(c *Composer, events *EventDispatcher, localRepo *installed.Repository, o AutoloadOptions) (*autoload.Result, bool, error) {
	pm := pluginManagerOf(c)
	if !pm.active() || !pm.alReplaced || pm.c != c || pm.alCtx != nil {
		return nil, false, nil
	}
	ctx := &autoloadCtx{c: c, events: events, localRepo: localRepo, o: o}
	pm.alCtx = ctx
	defer func() { pm.alCtx = nil }()
	args := phpjson.NewMap()
	if o.DevMode != nil {
		args.Set("dev", *o.DevMode)
	} else {
		args.Set("dev", nil)
	}
	args.Set("optimize", o.Optimize)
	args.Set("authoritative", o.ClassMapAuthoritative)
	args.Set("apcu", o.Apcu)
	if o.ApcuPrefix != "" {
		args.Set("apcuPrefix", o.ApcuPrefix)
	} else {
		args.Set("apcuPrefix", nil)
	}
	args.Set("runScripts", o.RunScripts)
	args.Set("dryRun", o.DryRun)
	args.Set("strictAmbiguous", o.StrictAmbiguous)
	if o.Suffix != "" {
		args.Set("suffix", o.Suffix)
	} else {
		args.Set("suffix", nil)
	}
	args.Set("filter", goFilterSpec(o.PlatformFilter))
	args.Set("frames", pm.frames())
	prevRepo := pm.curRepo
	pm.curRepo = localRepo
	_, err := pm.host.call("autoloadVia", args)
	pm.curRepo = prevRepo
	if err == nil {
		err = pm.takePendingErr()
	}
	if err != nil {
		return nil, true, err
	}
	if ctx.result == nil {
		// the generator did not dump (Composer returns its ClassMap)
		ctx.result = &autoload.Result{ClassMap: autoload.NewClassMap()}
	}
	return ctx.result, true, nil
}

// goFilterSpec describes a platform requirement filter of gocomposer for the
// runtime (which creates the Composer filter object).
func goFilterSpec(f solver.PlatformFilter) phpjson.Value {
	m := phpjson.NewMap()
	switch t := f.(type) {
	case nil, solver.IgnoreNothing:
		m.Set("kind", "nothing")
	case solver.IgnoreAll:
		m.Set("kind", "all")
	case *solver.IgnoreList:
		m.Set("kind", "list")
		reqs := []phpjson.Value{}
		for _, r := range t.Reqs() {
			reqs = append(reqs, r)
		}
		m.Set("reqs", reqs)
	case *phpPlatformFilter:
		m.Set("kind", "php")
		m.Set("id", t.id)
	default:
		m.Set("kind", "nothing")
	}
	return m
}

// filterFromSpec is the platform requirement filter of a runtime filter
// description (true/false, {kind: ...}).
func (pm *PluginManager) filterFromSpec(v phpjson.Value) solver.PlatformFilter {
	switch t := v.(type) {
	case bool:
		if t {
			return solver.IgnoreAll{}
		}
		return solver.IgnoreNothing{}
	}
	m := phpjson.ToMap(v)
	if m == nil {
		return solver.IgnoreNothing{}
	}
	switch m.Str("kind") {
	case "all":
		return solver.IgnoreAll{}
	case "list":
		var reqs []string
		phpjson.Each(mustGet(m, "reqs"), func(_ string, x phpjson.Value) {
			if s, ok := x.(string); ok {
				reqs = append(reqs, s)
			}
		})
		return solver.NewIgnoreList(reqs)
	case "php":
		id, _ := mustGet(m, "id").(int64)
		return &phpPlatformFilter{pm: pm, id: id, cache: map[string]bool{}}
	}
	return solver.IgnoreNothing{}
}

// phpPlatformFilter is a PlatformRequirementFilterInterface implemented by
// plugin code: gocomposer asks the PHP object (answers are cached).
type phpPlatformFilter struct {
	pm    *PluginManager
	id    int64
	cache map[string]bool
}

func (f *phpPlatformFilter) ask(m, req string) bool {
	key := m + "\x00" + req
	if v, ok := f.cache[key]; ok {
		return v
	}
	if !f.pm.active() {
		return false
	}
	args := phpjson.NewMap()
	args.Set("id", f.id)
	args.Set("m", m)
	args.Set("req", req)
	v, err := f.pm.host.call("filterCall", args)
	if err != nil {
		f.pm.setPendingErr(err)
		return false
	}
	f.cache[key] = v == true
	return v == true
}

func (f *phpPlatformFilter) IsIgnored(req string) bool { return f.ask("isIgnored", req) }

func (f *phpPlatformFilter) IsUpperBoundIgnored(req string) bool {
	return f.ask("isUpperBoundIgnored", req)
}

// pluginAutoloadOptions applies the state of the runtime's AutoloadGenerator
// (setters and dump() arguments) to o.
func (pm *PluginManager) pluginAutoloadOptions(a *phpjson.Map, o AutoloadOptions) AutoloadOptions {
	if v, ok := mustGet(a, "dev").(bool); ok {
		o.DevMode = &v
	} else {
		o.DevMode = nil
	}
	o.Optimize = mustGet(a, "optimize") == true
	o.ClassMapAuthoritative = mustGet(a, "authoritative") == true
	o.Apcu = mustGet(a, "apcu") == true
	o.ApcuPrefix, _ = mustGet(a, "apcuPrefix").(string)
	o.RunScripts = mustGet(a, "runScripts") == true
	o.DryRun = mustGet(a, "dryRun") == true
	o.StrictAmbiguous = mustGet(a, "strictAmbiguous") == true
	o.Suffix, _ = mustGet(a, "suffix").(string)
	if t, _ := mustGet(a, "targetDir").(string); t != "" && t != "composer" {
		o.TargetDir = t
	} else {
		o.TargetDir = ""
	}
	if spec := phpjson.ToMap(mustGet(a, "platformFilter")); spec != nil && spec.Str("kind") == "go" && pm.alCtx != nil {
		o.PlatformFilter = pm.alCtx.o.PlatformFilter
	} else {
		o.PlatformFilter = pm.filterFromSpec(mustGet(a, "platformFilter"))
	}
	return o
}

// configs created by plugin code (new Config(), Factory::createConfig()):
// gocomposer keeps a configuration per object.

func (pm *PluginManager) addDetachedConfig(cfg *config.Config) int64 {
	if pm.detachedCfgs == nil {
		pm.detachedCfgs = map[int64]*config.Config{}
	}
	pm.cfgSeq++
	pm.detachedCfgs[pm.cfgSeq] = cfg
	return pm.cfgSeq
}

// configCallback serves the Config callbacks of the runtime (ok is false
// for other ops): the configurations created by plugin code ("cfg"), and
// the config sources and insecure protocol warnings of instances.
func (pm *PluginManager) configCallback(op string, a *phpjson.Map) (phpjson.Value, bool, error) {
	switch op {
	case "configNew":
		baseDir, _ := mustGet(a, "baseDir").(string)
		return pm.addDetachedConfig(config.New(mustGet(a, "useEnv") != false, baseDir)), true, nil
	case "configCreate":
		cwd, _ := mustGet(a, "cwd").(string)
		if cwd == "" {
			cwd = pm.c.Cwd
		}
		cfg, err := config.CreateConfigDebug(cwd, func(m string) { pm.c.IO.WriteError(m) }, debugLogger(pm.c.IO))
		if err != nil {
			return nil, true, err
		}
		if err := cfg.LoadAuthEnv(); err != nil {
			return nil, true, err
		}
		out := phpjson.NewMap()
		out.Set("cfg", pm.addDetachedConfig(cfg))
		home := cfg.GetString("home")
		out.Set("config", home+"/config.json")
		out.Set("auth", home+"/auth.json")
		return out, true, nil
	case "configSources":
		c, err := pm.instance(a)
		if err != nil {
			return nil, true, err
		}
		out := phpjson.NewMap()
		home := c.Config.GetString("home")
		out.Set("config", home+"/config.json")
		out.Set("auth", home+"/auth.json")
		if c.ComposerFile != "" && !pm.arrayConfigs[c] {
			file := absPath(c.ComposerFile)
			out.Set("config", file)
			if auth := filepath.Join(filepath.Dir(file), "auth.json"); isFile(auth) {
				out.Set("localAuth", auth)
			}
		}
		return out, true, nil
	case "configWarnSsl":
		ssl := phpjson.NewMap()
		ssl.Set(a.Str("option"), false)
		if a.Str("option") == "both" {
			ssl = phpjson.NewMap()
			ssl.Set("verify_peer", false)
			ssl.Set("verify_peer_name", false)
		}
		pm.sslVerifyWarning(a.Str("host"), ssl)
		return nil, true, nil
	case "configWarnInsecure":
		c, err := pm.instance(a)
		if err != nil {
			return nil, true, err
		}
		_ = c.HTTP.ProhibitURL(a.Str("url"))
		return nil, true, nil
	}
	id, ok := mustGet(a, "cfg").(int64)
	if !ok {
		return nil, false, nil
	}
	cfg := pm.detachedCfgs[id]
	if cfg == nil {
		return nil, true, &UserError{Msg: fmt.Sprintf("Unknown configuration %d", id), Code: 1}
	}
	switch op {
	case "configSetBaseDir":
		baseDir, _ := mustGet(a, "baseDir").(string)
		cfg.SetBaseDir(baseDir)
		return nil, true, nil
	case "configMerge":
		if data := phpjson.ToMap(mustGet(a, "config")); data != nil {
			source := a.Str("source")
			if source == "" {
				source = "unknown"
			}
			cfg.MergeSource(data, source)
		}
		return nil, true, nil
	}
	v, err := configValue(cfg, op, a)
	return v, true, err
}

// configValue answers the Config getters of the runtime for cfg.
func configValue(cfg *config.Config, op string, a *phpjson.Map) (phpjson.Value, error) {
	relative := false
	if n, _ := mustGet(a, "flags").(int64); n&1 != 0 {
		relative = true
	}
	switch op {
	case "config":
		if relative {
			return cfg.GetRelative(a.Str("key")), nil
		}
		return cfg.Get(a.Str("key")), nil
	case "configHas":
		return cfg.Raw().Has(a.Str("key")), nil
	case "configAll", "configRaw":
		out := phpjson.NewMap()
		out.Set("repositories", cfg.Repositories())
		if op == "configRaw" {
			out.Set("config", cfg.Raw())
			return out, nil
		}
		all := phpjson.NewMap()
		for _, k := range cfg.Raw().Keys() {
			if relative {
				all.Set(k, cfg.GetRelative(k))
			} else {
				all.Set(k, cfg.Get(k))
			}
		}
		out.Set("config", all)
		return out, nil
	case "configSource":
		return cfg.SourceOfValue(a.Str("key"))
	}
	return nil, &UserError{Msg: "Unknown configuration operation " + op, Code: 1}
}

// validateWithSchemaFile is JsonFile::validateJsonSchema with a schema file
// given by plugin code: {"$ref": file} (file:// prepended to paths), the
// errors formatted like JsonValidationException's.
func validateWithSchemaFile(data phpjson.Value, file string) (phpjson.Value, error) {
	if !strings.Contains(file, "://") {
		file = "file://" + file
	}
	schema := phpjson.NewMap()
	schema.Set("$ref", file)
	schema.Set("$schema", "https://json-schema.org/draft-04/schema#")
	vd := &jsonschema.Validator{}
	errs, err := vd.Validate(jsonschema.JSONRoundTrip(data), schema)
	if err != nil {
		return nil, err
	}
	out := []phpjson.Value{}
	for _, m := range jsonschema.FormatErrors(errs) {
		out = append(out, m)
	}
	return out, nil
}
