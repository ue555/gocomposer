package composer

// Installer selection for packages handled by plugins: Composer's
// InstallationManager asks the installers added by plugins first
// (InstallationManager::getInstaller of Composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT license, (c) Nils Adermann,
// Jordi Boggiano) and PluginInstaller registers/deactivates plugins around
// the operations on composer-plugin/composer-installer packages.

import (
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/fsutil"

	"gocomposer/internal/installed"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

// pluginInstallerInfo describes an installer added by a plugin.
type pluginInstallerInfo struct {
	ID        int64
	Class     string
	Library   bool
	Overrides []string
}

// installerFor returns the plugin installer supporting typ (nil when the
// package is handled by gocomposer's own installers).
func (pm *PluginManager) installerFor(typ string) (*pluginInstallerInfo, error) {
	if !pm.active() {
		return nil, nil
	}
	typ = strings.ToLower(typ)
	if pm.installers == nil {
		pm.installers = map[string]*pluginInstallerInfo{}
	}
	if info, ok := pm.installers[typ]; ok {
		return info, nil
	}
	args := phpjson.NewMap()
	args.Set("type", typ)
	v, err := pm.host.call("installerFor", args)
	if err != nil {
		return nil, err
	}
	var info *pluginInstallerInfo
	if m := phpjson.ToMap(v); m != nil {
		info = &pluginInstallerInfo{Class: m.Str("class")}
		if id, ok := mustGet(m, "id").(int64); ok {
			info.ID = id
		}
		info.Library = mustGet(m, "library") == true
		phpjson.Each(mustGet(m, "overrides"), func(_ string, o phpjson.Value) {
			if s, ok := o.(string); ok {
				info.Overrides = append(info.Overrides, s)
			}
		})
	}
	pm.installers[typ] = info
	return info, nil
}

// pluginInstallPath is the install path chosen by a plugin installer for p
// (ok is false when no plugin installer handles the package type).
func (im *InstallationManager) pluginInstallPath(p *pkg.Package) (string, bool) {
	pm := pluginManagerOf(im.c)
	if !pm.active() {
		return "", false
	}
	info, err := pm.installerFor(p.Type())
	if err != nil {
		pm.setPendingErr(err)
		return "", false
	}
	if info == nil {
		return "", false
	}
	if path, ok := pm.installPaths[p]; ok {
		return path, true
	}
	args := phpjson.NewMap()
	args.Set("p", pm.host.ref(p))
	v, err := pm.host.call("installPath", args)
	if err != nil {
		pm.setPendingErr(err)
		return "", false
	}
	path, _ := phpjson.ScalarString(v)
	// installer-paths templates end with a slash; gocomposer's file
	// operations (like Composer's Filesystem) work on normalized paths
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	if pm.installPaths == nil {
		pm.installPaths = map[*pkg.Package]string{}
	}
	pm.installPaths[p] = path
	return path, true
}

func (pm *PluginManager) setPendingErr(err error) {
	if pm.pendingErr == nil {
		pm.pendingErr = err
	}
}

// takePendingErr returns an error raised by plugin code where gocomposer
// could not propagate it immediately (install path lookups).
func (pm *PluginManager) takePendingErr() error {
	if pm == nil {
		return nil
	}
	err := pm.pendingErr
	pm.pendingErr = nil
	return err
}

// pluginDownloadCheck is PluginInstaller::download's extra.class check.
func pluginDownloadCheck(p *pkg.Package) error {
	if isPluginType(p.Type()) && len(pluginClasses(p)) == 0 {
		return &UserError{Msg: "Error while installing " + p.PrettyName + ", composer-plugin packages should have a class defined in their extra key to be usable.", Code: 1}
	}
	return nil
}

// pluginPrepare is PluginInstaller::prepare: installing a plugin fails early
// when it is blocked by allow-plugins.
func (im *InstallationManager) pluginPrepare(op *solver.Operation) error {
	if op.Type != solver.OpInstall && op.Type != solver.OpUpdate {
		return nil
	}
	p := op.MainPackage()
	pm := pluginManagerOf(im.c)
	if pm == nil || pm.disabled {
		return nil
	}
	if info, err := pm.installerFor(p.Type()); err != nil || info != nil {
		return err
	}
	if !isPluginType(p.Type()) {
		return nil
	}
	_, err := pm.isPluginAllowed(p.Name, false, pluginOptional(p), true)
	return err
}

// beginExecute records the operations of InstallationManager::execute for
// package events and installer callbacks.
func (im *InstallationManager) beginExecute(repo *installed.Repository, ops []*solver.Operation) func() {
	pm := pluginManagerOf(im.c)
	if pm == nil {
		return func() {}
	}
	prevIM, prevRepo, prevOps, prevKey := pm.curIM, pm.curRepo, pm.curOps, pm.curOpsKey
	pm.curIM, pm.curRepo, pm.curOps = im, repo, ops
	pm.opsSeq++
	pm.curOpsKey = "ops" + itoa(pm.opsSeq)
	key := pm.curOpsKey
	return func() {
		if pm.active() {
			args := phpjson.NewMap()
			args.Set("key", key)
			_, _ = pm.host.call("freeOps", args)
		}
		pm.curIM, pm.curRepo, pm.curOps, pm.curOpsKey = prevIM, prevRepo, prevOps, prevKey
	}
}

// runOp performs an install/update/uninstall operation with the installer
// Composer would use for the package type (InstallationManager::install/
// update/uninstall) including the PluginInstaller behaviour.
func (im *InstallationManager) runOp(repo *installed.Repository, op *solver.Operation) error {
	pm := pluginManagerOf(im.c)
	if pm != nil {
		prevIM, prevRepo := pm.curIM, pm.curRepo
		pm.curIM, pm.curRepo = im, repo
		defer func() { pm.curIM, pm.curRepo = prevIM, prevRepo }()
	}
	var err error
	switch op.Type {
	case solver.OpInstall:
		err = im.opInstall(pm, repo, op.Package)
	case solver.OpUpdate:
		err = im.opUpdate(pm, repo, op.Initial, op.Target)
	case solver.OpUninstall:
		err = im.opUninstall(pm, repo, op.Package)
	}
	if err == nil {
		err = pm.takePendingErr()
	}
	return err
}

func (im *InstallationManager) customInstaller(pm *PluginManager, typ string) (*pluginInstallerInfo, error) {
	if pm == nil {
		return nil, nil
	}
	return pm.installerFor(typ)
}

func (im *InstallationManager) opInstall(pm *PluginManager, repo *installed.Repository, p *pkg.Package) error {
	info, err := im.customInstaller(pm, p.Type())
	if err != nil {
		return err
	}
	if info != nil {
		return pm.installerCall(info, p.Type(), "install", nil, p)
	}
	if err := im.install(repo, p); err != nil {
		return err
	}
	if pm != nil && isPluginType(p.Type()) {
		if err := pm.registerPackage(p, true, false); err != nil {
			return im.rollbackPluginInstall(repo, p, err)
		}
	}
	return nil
}

// rollbackPluginInstall is PluginInstaller::rollbackInstall.
func (im *InstallationManager) rollbackPluginInstall(repo *installed.Repository, p *pkg.Package, err error) error {
	if _, ok := err.(*ScriptError); ok {
		return err
	}
	im.c.IO.WriteError("Plugin initialization failed (" + err.Error() + "), uninstalling plugin")
	if uerr := im.uninstall(repo, p); uerr != nil {
		return uerr
	}
	return err
}

func (im *InstallationManager) opUpdate(pm *PluginManager, repo *installed.Repository, initial, target *pkg.Package) error {
	ii, err := im.customInstaller(pm, initial.Type())
	if err != nil {
		return err
	}
	ti, err := im.customInstaller(pm, target.Type())
	if err != nil {
		return err
	}
	if ii != nil || ti != nil {
		if initial.Type() != target.Type() {
			// InstallationManager::update with differing package types
			if err := im.opUninstall(pm, repo, initial); err != nil {
				return err
			}
			return im.opInstall(pm, repo, target)
		}
		return pm.installerCall(ti, target.Type(), "update", initial, target)
	}
	if err := im.update(repo, initial, target); err != nil {
		return err
	}
	if pm != nil && isPluginType(target.Type()) {
		if err := pm.deactivatePackage(initial); err != nil {
			return err
		}
		if err := pm.registerPackage(target, true, false); err != nil {
			return im.rollbackPluginInstall(repo, target, err)
		}
	}
	return nil
}

func (im *InstallationManager) opUninstall(pm *PluginManager, repo *installed.Repository, p *pkg.Package) error {
	info, err := im.customInstaller(pm, p.Type())
	if err != nil {
		return err
	}
	if info != nil {
		return pm.installerCall(info, p.Type(), "uninstall", nil, p)
	}
	if pm != nil && isPluginType(p.Type()) {
		if err := pm.uninstallPackage(p); err != nil {
			return err
		}
	}
	return im.uninstall(repo, p)
}

// customFlow reports whether the plugin installer runs Composer's whole
// operation sequence itself (download, prepare, the operation, cleanup):
// installers not built on LibraryInstaller and LibraryInstaller subclasses
// replacing a download hook. The other LibraryInstaller subclasses only get
// the operation, the download and prepare/cleanup steps being Composer's own
// (done by gocomposer).
func (info *pluginInstallerInfo) customFlow(stage string) bool {
	if !info.Library {
		return true
	}
	for _, o := range info.Overrides {
		if o == stage {
			return true
		}
	}
	return false
}

// installerCall runs a method of a plugin installer. LibraryInstaller code
// (also when called through parent::) performs its work through gocomposer's
// LibraryInstaller and DownloadManager (callbacks libOp and dl).
func (pm *PluginManager) installerCall(info *pluginInstallerInfo, typ, method string, initial, target *pkg.Package) error {
	return pm.installerCallOp(info, method, "", initial, target)
}

func (pm *PluginManager) installerCallOp(info *pluginInstallerInfo, method, opType string, initial, target *pkg.Package) error {
	args := pm.installerOpArgs(info, method, initial, target)
	if opType != "" {
		args.Set("type", opType)
	}
	_, err := pm.host.call("installerOp", args)
	if err == nil {
		err = pm.takePendingErr()
	}
	if err != nil && pm.lastCallbackErr != nil {
		// a gocomposer error raised under the plugin installer
		if pe, ok := err.(*PHPException); ok && pe.Message == pm.lastCallbackErr.Error() {
			err = pm.lastCallbackErr
		}
	}
	pm.lastCallbackErr = nil
	return err
}

func (pm *PluginManager) installerOpArgs(info *pluginInstallerInfo, method string, initial, target *pkg.Package) *phpjson.Map {
	args := phpjson.NewMap()
	args.Set("id", info.ID)
	args.Set("m", method)
	args.Set("p", pm.host.ref(target))
	args.Set("psrc", installationSourceValue(target))
	if initial != nil {
		args.Set("i", pm.host.ref(initial))
		args.Set("isrc", installationSourceValue(initial))
	}
	args.Set("frames", pm.frames())
	return args
}

func installationSourceValue(p *pkg.Package) phpjson.Value {
	if s := p.InstallationSource(); s != "" {
		return s
	}
	return nil
}

// cleanup installers: InstallationManager::downloadAndExecuteBatch creates
// the cleanup step of every operation with the installer of the package
// type at download time (a plugin installed in the same run does not clean
// up the packages downloaded before its activation)

type cleanupChoice struct{ info *pluginInstallerInfo }

func (pm *PluginManager) cleanupInstaller(p *pkg.Package) (*pluginInstallerInfo, bool) {
	if pm.cleanupInstallers == nil {
		return nil, false
	}
	c, ok := pm.cleanupInstallers[p]
	if !ok {
		return nil, false
	}
	return c.info, true
}

// pluginRecordInstallers records the installers of the operations at
// download time.
func (im *InstallationManager) pluginRecordInstallers(ops []*solver.Operation) error {
	pm := pluginManagerOf(im.c)
	if pm == nil {
		return nil
	}
	if pm.cleanupInstallers == nil {
		pm.cleanupInstallers = map[*pkg.Package]cleanupChoice{}
	}
	for _, op := range ops {
		switch op.Type {
		case solver.OpInstall, solver.OpUpdate, solver.OpUninstall:
		default:
			continue
		}
		p := op.MainPackage()
		var info *pluginInstallerInfo
		if pm.active() {
			var err error
			if info, err = pm.installerFor(p.Type()); err != nil {
				return err
			}
		}
		pm.cleanupInstallers[p] = cleanupChoice{info}
	}
	return nil
}

// cleanupAfterFailure is InstallationManager::runCleanup for one operation:
// the cleanup of the installer chosen at download time (errors ignored).
func (im *InstallationManager) cleanupAfterFailure(op *solver.Operation) {
	if pm := pluginManagerOf(im.c); pm.active() {
		// the cleanup of a downloader replaced by a plugin
		if handled, _ := im.replacedStage(op, "cleanup"); handled {
			pm.lastCallbackErr = nil
			im.cleanup(op.MainPackage())
			return
		}
		if info, ok := pm.cleanupInstaller(op.MainPackage()); ok && info != nil && info.customFlow("cleanup") {
			var initial *pkg.Package
			if op.Type == solver.OpUpdate {
				initial = op.Initial
			}
			_ = pm.installerCallOp(info, "cleanup", op.Type, initial, op.MainPackage())
			pm.lastCallbackErr = nil
			return
		}
	}
	im.cleanup(op.MainPackage())
}

// pendingOp is an operation run by a plugin installer as a promise chain
// (settled with the batch, see pluginStartOp).
type pendingOp struct {
	im        *InstallationManager
	repo      *installed.Repository
	op        *solver.Operation
	devMode   bool
	goCleanup bool
	// replaced: the operation of a downloader replaced by a plugin
	// (plugin_services.go): its cleanup ran in the chain, gocomposer only
	// removes its own temporary files
	replaced bool
}

// pluginStartOp starts an operation of a plugin installer like
// InstallationManager::executeBatch: the installer's install/update/
// uninstall, then the cleanup of the installer chosen at download time,
// then the repository write, chained on the operation's promise and
// settled by pluginSettleOps. started is false for gocomposer's installers.
func (im *InstallationManager) pluginStartOp(repo *installed.Repository, op *solver.Operation, devMode bool) (bool, error) {
	pm := pluginManagerOf(im.c)
	if !pm.active() {
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
	info, err := pm.installerFor(p.Type())
	if err != nil || info == nil {
		return false, err
	}
	if op.Type == solver.OpUpdate {
		if ii, err := pm.installerFor(initial.Type()); err != nil || ii != info {
			// InstallationManager::update with differing installers
			return false, err
		}
	}
	cleanup, recorded := pm.cleanupInstaller(p)
	if !recorded {
		cleanup = info
	}
	pm.opSeq++
	key := "op" + itoa(pm.opSeq)
	if pm.pendingOps == nil {
		pm.pendingOps = map[string]*pendingOp{}
	}
	po := &pendingOp{im: im, repo: repo, op: op, devMode: devMode}
	args := pm.installerOpArgs(info, op.Type, initial, p)
	args.Set("defer", true)
	args.Set("key", key)
	if cleanup != nil && cleanup.customFlow("cleanup") {
		args.Set("cleanup", cleanup.ID)
	} else {
		po.goCleanup = true
	}
	pm.pendingOps[key] = po
	prevIM, prevRepo := pm.curIM, pm.curRepo
	pm.curIM, pm.curRepo = im, repo
	_, err = pm.host.call("installerOp", args)
	pm.curIM, pm.curRepo = prevIM, prevRepo
	if err == nil {
		err = pm.takePendingErr()
	}
	return true, err
}

// opDone is the end of a plugin installer operation's chain: gocomposer's
// cleanup (when the download was gocomposer's) and the repository write.
func (pm *PluginManager) opDone(key string) error {
	po := pm.pendingOps[key]
	if po == nil {
		return nil
	}
	delete(pm.pendingOps, key)
	im := po.im
	if po.goCleanup {
		if err := im.cleanupOp(po.op); err != nil {
			return err
		}
		im.cleanup(po.op.MainPackage())
	}
	if po.replaced {
		im.cleanup(po.op.MainPackage())
	}
	return im.writeRepo(po.repo, po.devMode)
}

// pluginSettleOps waits for the operations started by pluginStartOp
// (InstallationManager::waitOnPromises).
func (im *InstallationManager) pluginSettleOps(repo *installed.Repository) error {
	pm := pluginManagerOf(im.c)
	if !pm.active() {
		return nil
	}
	prevIM, prevRepo := pm.curIM, pm.curRepo
	pm.curIM, pm.curRepo = im, repo
	defer func() { pm.curIM, pm.curRepo = prevIM, prevRepo }()
	_, err := pm.host.call("settleOps", nil)
	if err == nil {
		err = pm.takePendingErr()
	}
	if err != nil && pm.lastCallbackErr != nil {
		if pe, ok := err.(*PHPException); ok && pe.Message == pm.lastCallbackErr.Error() {
			err = pm.lastCallbackErr
		}
	}
	pm.lastCallbackErr = nil
	return err
}

// pluginInstallerDownload is InstallationManager's download for packages of
// a plugin installer that downloads itself (handled is false otherwise).
func (im *InstallationManager) pluginInstallerDownload(p, prev *pkg.Package) (bool, error) {
	pm := pluginManagerOf(im.c)
	if !pm.active() {
		return false, nil
	}
	info, err := pm.installerFor(p.Type())
	if err != nil {
		return true, err
	}
	if info == nil || !info.customFlow("download") {
		return false, nil
	}
	return true, pm.installerCall(info, p.Type(), "download", prev, p)
}

// pluginInstallerStage runs the prepare/cleanup step of an operation through
// a plugin installer that replaces it (handled is false otherwise).
func (im *InstallationManager) pluginInstallerStage(op *solver.Operation, stage string) (bool, error) {
	pm := pluginManagerOf(im.c)
	if !pm.active() {
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
	info, err := pm.installerFor(p.Type())
	if err != nil {
		return true, err
	}
	if stage == "cleanup" {
		if recorded, ok := pm.cleanupInstaller(p); ok {
			info = recorded
		}
	}
	if info == nil || !info.customFlow(stage) {
		return false, nil
	}
	return true, pm.installerCallOp(info, stage, op.Type, initial, p)
}

// pluginDownloaderOp runs a downloader a plugin registered for a dist/source
// type gocomposer does not know (DownloadManager::setDownloader).
func (im *InstallationManager) pluginDownloaderOp(typ, method string, p, prev *pkg.Package) (bool, error) {
	pm := pluginManagerOf(im.c)
	if !pm.active() || typ == "" {
		return false, nil
	}
	switch typ {
	case "zip", "tar", "gzip", "file", "phar", "path", "git", "svn", "hg", "fossil", "perforce", "rar", "xz":
		// gocomposer's own downloaders (a plugin replacing one of them with
		// DownloadManager::setDownloader is not supported)
		return false, nil
	}
	args := phpjson.NewMap()
	args.Set("type", typ)
	v, err := pm.host.call("downloaderFor", args)
	if err != nil {
		return true, err
	}
	if v == nil {
		return false, nil
	}
	args = phpjson.NewMap()
	args.Set("type", typ)
	args.Set("m", method)
	args.Set("p", pm.host.ref(p))
	if prev != nil {
		args.Set("i", pm.host.ref(prev))
	}
	args.Set("path", im.InstallPath(p))
	args.Set("frames", pm.frames())
	_, err = pm.host.call("downloaderOp", args)
	if err == nil {
		err = pm.takePendingErr()
	}
	return true, err
}

// libOp performs the LibraryInstaller work requested by the plugin runtime.
func (pm *PluginManager) libOp(a *phpjson.Map) (phpjson.Value, error) {
	im, repo := pm.curIM, pm.curRepo
	if im == nil {
		im = NewInstallationManager(pm.c)
	}
	if repo == nil {
		repo = pm.c.LocalRepo
	}
	p, err := pm.host.pkgFromRef(mustGet(a, "p"))
	if err != nil {
		return nil, err
	}
	switch a.Str("m") {
	case "install":
		return nil, im.install(repo, p)
	case "update":
		initial, err := pm.host.pkgFromRef(mustGet(a, "i"))
		if err != nil {
			return nil, err
		}
		return nil, im.update(repo, initial, p)
	case "uninstall":
		return nil, im.uninstall(repo, p)
	case "installBinaries":
		warn := mustGet(a, "warn") != false
		im.installBinaries(p, a.Str("path"), warn)
		return nil, nil
	case "removeBinaries":
		im.removeBinaries(p)
		return nil, nil
	}
	return nil, &UserError{Msg: "unknown library installer operation " + a.Str("m"), Code: 1}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

// dropDownloadDir removes the (still empty) install directory created by a
// file download made for another destination (archive): Composer downloads
// those into the destination directly.
func (im *InstallationManager) dropDownloadDir(p *pkg.Package) {
	im.mu.Lock()
	dir, ok := im.downloadPaths[p]
	delete(im.downloadPaths, p)
	im.mu.Unlock()
	if !ok {
		return
	}
	for _, d := range []string{dir, filepath.Dir(dir)} {
		if fsutil.Realpath(d) == fsutil.Realpath(im.vendorDir) {
			break
		}
		if isDir(d) && fsutil.IsDirEmpty(d) {
			os.Remove(d)
		}
	}
}
