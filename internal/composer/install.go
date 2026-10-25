package composer

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/cio"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/installed"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

// InstallationManager ports InstallationManager + LibraryInstaller +
// DownloadManager + the File/Archive/Path downloaders.
type InstallationManager struct {
	c            *Composer
	vendorDir    string
	binDir       string
	filesCache   *cache.Cache
	preferSource bool
	preferDist   bool
	preferences  []preference
	tmpFiles     map[*pkg.Package]string
	cacheWrites  map[string]string
	mu           sync.Mutex
	// ScriptHook runs package events (pre/post-package-*).
	PackageEvent func(event string, op *solver.Operation) error
	// OutputProgress is InstallationManager::setOutputProgress (!--no-progress).
	OutputProgress bool
	// Trace, when set, turns Execute into InstallationManagerMock::execute
	// (used by the upstream installer fixture tests): operations are only
	// recorded and applied to the repository.
	Trace *[]string
	// InstallPathOverride replaces InstallPath (mock: vendor/<name>).
	InstallPathOverride func(p *pkg.Package) string
	// gitDownloader handles packages installed from git source.
	gitDownloader *gitDownloader
	// downloadPaths are the install paths file downloads were made for
	// (FileDownloader::cleanup uses the path of the installer chosen when
	// downloading, before plugins installed in the same run were activated).
	downloadPaths map[*pkg.Package]string
	// deferredBins are binary installs run after the batch's async
	// extraction processes were waited on (executeBatch)
	deferredBins []func()
}

type preference struct {
	re    *regexp.Regexp
	value string
}

// NewInstallationManager creates the manager from config.
func NewInstallationManager(c *Composer) *InstallationManager {
	im := &InstallationManager{c: c, tmpFiles: map[*pkg.Package]string{}, cacheWrites: map[string]string{}}
	im.vendorDir = strings.TrimRight(c.Config.GetString("vendor-dir"), "/")
	im.binDir = strings.TrimRight(c.Config.GetString("bin-dir"), "/")
	im.filesCache = cache.New(c.Config.GetString("cache-files-dir"), "a-z0-9_./", c.Config.GetBool("cache-read-only"))
	im.filesCache.Warn = func(s string) { c.IO.WriteError(s) }
	switch v := c.Config.Get("preferred-install").(type) {
	case string:
		if v == "source" {
			im.preferSource = true
		} else if v == "dist" {
			im.preferDist = true
		}
	case *phpjson.Map:
		for _, k := range v.Keys() {
			s, _ := phpjson.ScalarString(mustGet(v, k))
			im.preferences = append(im.preferences, preference{regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(k), `\*`, `.*`) + `$`), s})
		}
	}
	return im
}

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

// SetPreferSource / SetPreferDist mirror the command line flags.
func (im *InstallationManager) SetPreferSource(b bool) {
	if b {
		im.preferSource = true
	}
}
func (im *InstallationManager) SetPreferDist(b bool) {
	if b {
		im.preferDist = true
	}
}

func (im *InstallationManager) realVendorDir() string {
	_ = os.MkdirAll(im.vendorDir, 0o777)
	if r := fsutil.Realpath(im.vendorDir); r != "" {
		return r
	}
	return im.vendorDir
}

// InstallPath is LibraryInstaller::getInstallPath ("" for metapackages).
func (im *InstallationManager) InstallPath(p *pkg.Package) string {
	if im.InstallPathOverride != nil {
		return im.InstallPathOverride(p)
	}
	// installers added by Composer plugins are asked first (plugin_install.go)
	if path, ok := im.pluginInstallPath(p); ok {
		return path
	}
	if p.Type() == "metapackage" {
		return ""
	}
	base := im.realVendorDir() + "/" + p.PrettyName
	if td := p.TargetDir(); td != "" {
		return base + "/" + td
	}
	return base
}

func (im *InstallationManager) packageBasePath(p *pkg.Package) string {
	ip := im.InstallPath(p)
	if td := p.TargetDir(); td != "" {
		re := regexp.MustCompile(`/*` + strings.ReplaceAll(regexp.QuoteMeta(td), "/", "/+") + `/?$`)
		return re.ReplaceAllString(ip, "")
	}
	return ip
}

func (im *InstallationManager) resolvePreference(p *pkg.Package) string {
	for _, pr := range im.preferences {
		if pr.re.MatchString(p.Name) {
			if pr.value == "dist" || (!p.IsDev() && pr.value == "auto") {
				return "dist"
			}
			return "source"
		}
	}
	if p.IsDev() {
		return "source"
	}
	return "dist"
}

func (im *InstallationManager) availableSources(p, prev *pkg.Package) ([]string, error) {
	var sources []string
	if p.SourceType() != "" {
		sources = append(sources, "source")
	}
	if p.DistType() != "" {
		sources = append(sources, "dist")
	}
	if len(sources) == 0 {
		return nil, &UserError{Msg: "Package " + p.String() + " must have a source or dist specified", Code: 1}
	}
	if prev != nil && containsStr(sources, prev.InstallationSource()) &&
		!(!prev.IsDev() && prev.InstallationSource() == "dist" && p.IsDev()) {
		prevSource := prev.InstallationSource()
		if sources[0] != prevSource {
			sources[0], sources[1] = sources[1], sources[0]
		}
		return sources, nil
	}
	if !im.preferSource && (im.preferDist || im.resolvePreference(p) == "dist") {
		for i, j := 0, len(sources)-1; i < j; i, j = i+1, j-1 {
			sources[i], sources[j] = sources[j], sources[i]
		}
	}
	return sources, nil
}

func containsStr(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

func (im *InstallationManager) downloaderType(p *pkg.Package) string {
	if p.Type() == "metapackage" {
		return ""
	}
	if p.InstallationSource() == "source" {
		return p.SourceType()
	}
	return p.DistType()
}

// download prepares the package files (DownloadManager::download). It prints
// the "Downloading"/cache messages synchronously and returns the network
// transfer (if any) to be run concurrently.
func (im *InstallationManager) download(p, prev *pkg.Package) (func() error, error) {
	// installers added by plugins that download themselves (plugin_install.go)
	if handled, err := im.pluginInstallerDownload(p, prev); handled {
		return nil, err
	}
	if p.Type() == "metapackage" {
		return nil, nil
	}
	sources, err := im.availableSources(p, prev)
	if err != nil {
		return nil, err
	}
	source := sources[0]
	p.SetInstallationSource(source)
	transfer, err := im.downloadWith(p, prev)
	// DownloadManager::download falls back from source to dist (never the
	// other way around) when the source download fails with a
	// RuntimeException.
	for err != nil && source == "source" && len(sources) > 1 && sources[1] == "dist" && isRuntimeException(err) {
		im.c.IO.WriteError("    <warning>Failed to download " + p.PrettyName + " from " + source + ": " + err.Error() + "</warning>")
		sources = sources[1:]
		source = sources[0]
		im.c.IO.WriteError("    <warning>Now trying to download from " + source + "</warning>")
		p.SetInstallationSource(source)
		transfer, err = im.downloadWith(p, prev)
	}
	if err != nil {
		if isRuntimeException(err) {
			im.c.IO.WriteError("    <warning>Failed to download " + p.PrettyName + " from " + source + ": " + err.Error() + "</warning>")
		}
		return nil, err
	}
	if transfer == nil {
		return nil, nil
	}
	return func() error {
		err := transfer()
		if err != nil {
			im.c.IO.WriteError("    <warning>Failed to download " + p.PrettyName + " from " + source + ": " + err.Error() + "</warning>")
		}
		return err
	}, nil
}

func (im *InstallationManager) downloadWith(p, prev *pkg.Package) (func() error, error) {
	typ := im.downloaderType(p)
	// a built-in downloader replaced by a plugin (plugin_services.go)
	if handled, err := im.replacedDownloaderOp(typ, "download", p, prev); handled {
		return nil, err
	}
	switch typ {
	case "zip", "tar", "gzip", "file", "phar", "xz", "rar":
		if typ == "zip" {
			// ZipDownloader::download (archive_cmds.go)
			if err := im.zipDownloadCheck(); err != nil {
				return nil, err
			}
		}
		return im.fileDownload(p)
	case "path":
		return nil, im.pathDownload(p)
	case "git":
		return nil, im.gitDL().download(p, prev)
	case "svn", "hg", "fossil", "perforce":
		return nil, im.sourceDL(typ).download(p, prev)
	}
	// downloaders added by plugins (DownloadManager::setDownloader)
	if handled, err := im.pluginDownloaderOp(typ, "download", p, prev); handled {
		return nil, err
	}
	return nil, &UserError{Msg: "Unknown downloader type: " + typ + ". Available types: zip, tar, gzip, file, phar, path (gocomposer)", Code: 1}
}

var (
	reGitHubLegacy  = regexp.MustCompile(`(?i)^https?://(?:www\.)?github\.com/([^/]+)/([^/]+)/(zip|tar)ball/(.+)$`)
	reGitHubArchive = regexp.MustCompile(`(?i)^https?://(?:www\.)?github\.com/([^/]+)/([^/]+)/archive/.+\.(zip|tar)(?:\.gz)?$`)
	reGitHubAPI     = regexp.MustCompile(`(?i)^https?://api\.github\.com/repos/([^/]+)/([^/]+)/(zip|tar)ball(?:/.+)?$`)
	reBitbucketGet  = regexp.MustCompile(`(?i)^https?://(?:www\.)?bitbucket\.org/([^/]+)/([^/]+)/get/(.+)\.(zip|tar\.gz|tar\.bz2)$`)
	reGitLabArchive = regexp.MustCompile(`(?i)^https?://(?:www\.)?gitlab\.com/api/v[34]/projects/([^/]+)/repository/archive\.(zip|tar\.gz|tar\.bz2|tar)\?sha=.+$`)
)

// UpdateDistReference is Url::updateDistReference.
func (im *InstallationManager) UpdateDistReference(u, ref string) string {
	host := ""
	if pu, err := url.Parse(u); err == nil {
		host = pu.Hostname()
	}
	switch host {
	case "api.github.com", "github.com", "www.github.com":
		if m := reGitHubLegacy.FindStringSubmatch(u); m != nil {
			return "https://api.github.com/repos/" + m[1] + "/" + m[2] + "/" + m[3] + "ball/" + ref
		} else if m := reGitHubArchive.FindStringSubmatch(u); m != nil {
			return "https://api.github.com/repos/" + m[1] + "/" + m[2] + "/" + m[3] + "ball/" + ref
		} else if m := reGitHubAPI.FindStringSubmatch(u); m != nil {
			return "https://api.github.com/repos/" + m[1] + "/" + m[2] + "/" + m[3] + "ball/" + ref
		}
		return u
	case "bitbucket.org", "www.bitbucket.org":
		if m := reBitbucketGet.FindStringSubmatch(u); m != nil {
			return "https://bitbucket.org/" + m[1] + "/" + m[2] + "/get/" + ref + "." + m[4]
		}
		return u
	case "gitlab.com", "www.gitlab.com":
		if m := reGitLabArchive.FindStringSubmatch(u); m != nil {
			return "https://gitlab.com/api/v4/projects/" + m[1] + "/repository/archive." + m[2] + "?sha=" + ref
		}
		return u
	}
	if inList(im.c.Config.Get("github-domains"), host) {
		return regexp.MustCompile(`(?i)(/repos/[^/]+/[^/]+/(zip|tar)ball)(?:/.+)?$`).ReplaceAllString(u, "${1}/"+ref)
	}
	if inList(im.c.Config.Get("gitlab-domains"), host) {
		return regexp.MustCompile(`(?i)(/api/v[34]/projects/[^/]+/repository/archive\.(?:zip|tar\.gz|tar\.bz2|tar)\?sha=).+$`).ReplaceAllString(u, "${1}"+ref)
	}
	return u
}

func inList(v phpjson.Value, s string) bool {
	found := false
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		if x, _ := e.(string); x == s {
			found = true
		}
	})
	return found
}

func distPathComponent(p *pkg.Package, ext bool) string {
	u := strings.ReplaceAll(p.DistURL(), `\`, "/")
	pth := u
	if pu, err := url.Parse(u); err == nil {
		pth = pu.Path
	}
	base := path.Base(pth)
	if pth == "" || pth == "/" {
		base = ""
	}
	if ext {
		if i := strings.LastIndexByte(base, '.'); i >= 0 {
			return base[i+1:]
		}
		return ""
	}
	return base
}

func (im *InstallationManager) tmpFileName(p *pkg.Package) string {
	im.mu.Lock()
	defer im.mu.Unlock()
	if f, ok := im.tmpFiles[p]; ok {
		return f
	}
	ext := distPathComponent(p, true)
	if ext == "" {
		ext = p.DistType()
	}
	var rnd [8]byte
	_, _ = rand.Read(rnd[:])
	sum := md5.Sum([]byte(p.String() + hex.EncodeToString(rnd[:])))
	name := strings.TrimRight(im.vendorDir+"/composer/tmp-"+hex.EncodeToString(sum[:])+"."+ext, ".")
	im.tmpFiles[p] = name
	return name
}

func transportHeaders(p *pkg.Package) []string {
	to := p.TransportOptions()
	var hs []string
	if to == nil {
		return nil
	}
	if http := to.Map("http"); http != nil {
		if h, ok := http.Get("header"); ok {
			switch t := h.(type) {
			case string:
				hs = append(hs, t)
			default:
				phpjson.Each(t, func(_ string, v phpjson.Value) {
					if s, ok := v.(string); ok {
						hs = append(hs, s)
					}
				})
			}
		}
	}
	return hs
}

// fileDownload is FileDownloader::download (dist archives and files).
func (im *InstallationManager) fileDownload(p *pkg.Package) (func() error, error) {
	if p.DistURL() == "" {
		return nil, &UserError{Msg: "The given package is missing url information", Code: 1}
	}
	var urls []string
	urls = append(urls, p.DistURL())
	for _, m := range p.DistMirrors() {
		_ = m // dist mirrors are not used by gocomposer (none are configured by default)
	}
	fileName := im.tmpFileName(p)
	// FileDownloader::download ensures the install path exists
	if target := im.InstallPath(p); target != "" {
		if err := os.MkdirAll(target, 0o777); err != nil {
			return nil, err
		}
		im.mu.Lock()
		if im.downloadPaths == nil {
			im.downloadPaths = map[*pkg.Package]string{}
		}
		im.downloadPaths[p] = target
		im.mu.Unlock()
	}
	if err := os.MkdirAll(filepath.Dir(fileName), 0o777); err != nil {
		return nil, err
	}
	base := urls[0]
	processed := base
	if p.DistReference() != "" {
		processed = im.UpdateDistReference(base, p.DistReference())
	}
	sum := sha1.Sum([]byte(processed))
	cacheKey := p.Name + "/" + hex.EncodeToString(sum[:]) + "." + p.DistType()
	// PluginEvents::PRE_FILE_DOWNLOAD (plugin_dlevents.go)
	processed, cacheKey, err := im.pluginPreFileDownload(p, processed, cacheKey)
	if err != nil {
		return nil, err
	}
	checksum := p.DistSha1Checksum()
	verify := func() error {
		if _, err := os.Stat(fileName); err != nil {
			return &UserError{Msg: httpdl.Sanitize(base) + " could not be saved to " + fileName + ", make sure the directory is writable and you have internet connectivity", Code: 1}
		}
		if checksum != "" {
			h := sha1.New()
			f, err := os.Open(fileName)
			if err != nil {
				return err
			}
			_, err = io.Copy(h, f)
			f.Close()
			if err != nil {
				return err
			}
			if hex.EncodeToString(h.Sum(nil)) != checksum {
				im.clearLastCacheWrite(p)
				os.Remove(fileName)
				return &UserError{Msg: "The checksum verification of the file failed (downloaded from " + httpdl.Sanitize(base) + ")", Code: 1}
			}
		}
		return nil
	}
	if (checksum == "" || checksum == im.filesCache.Sha1(cacheKey)) && im.filesCache.CopyTo(cacheKey, fileName) {
		im.c.IO.WriteErrorV("  - Loading <info>"+p.Name+"</info> (<comment>"+p.FullPrettyVersion(true)+"</comment>) from cache", true, 128)
		if !im.filesCache.ReadOnly {
			im.mu.Lock()
			im.cacheWrites[p.Name] = cacheKey
			im.mu.Unlock()
		}
		if err := verify(); err != nil {
			return nil, err
		}
		// PluginEvents::POST_FILE_DOWNLOAD (plugin_dlevents.go)
		return nil, im.pluginPostFileDownload(p, fileName, checksum, processed, false)
	}
	im.c.IO.WriteError("  - Downloading <info>" + p.Name + "</info> (<comment>" + p.FullPrettyVersion(true) + "</comment>)")
	// HttpDownloader::addCopy checks the URL against secure-http when the
	// job is queued, so the insecure-protocol warning precedes the next
	// package's Downloading line (errors are reported by the transfer).
	_ = im.c.HTTP.ProhibitURL(processed)
	return func() error {
		var err error
		for attempt := 0; attempt < 4; attempt++ {
			if attempt > 0 {
				// FileDownloader's reject handler waits and re-enters
				// $download(), which dispatches pre-file-download and prints
				// the Downloading line again.
				time.Sleep(500 * time.Millisecond)
				var perr error
				if processed, cacheKey, perr = im.pluginPreFileDownloadRetry(p, processed, cacheKey); perr != nil {
					return perr
				}
				im.c.IO.WriteError("  - Downloading <info>" + p.Name + "</info> (<comment>" + p.FullPrettyVersion(true) + "</comment>)")
			}
			_, err = im.c.HTTP.Download(processed, fileName, transportHeaders(p)...)
			if err == nil && p.DistType() == "zip" {
				if ok, size := ZipHasEndOfCentralDirectory(fileName); !ok {
					err = &httpdl.TransportError{Msg: fmt.Sprintf("The downloaded archive for %s is truncated or corrupt, zip end of central directory not found (file size: %d bytes)", p.Name, size)}
					if regexp.MustCompile(`(?i)^https?://`).MatchString(processed) && attempt < 3 {
						im.c.IO.WriteError("    <warning>" + err.Error() + ", retrying</warning>")
					}
					os.Remove(fileName)
					if !regexp.MustCompile(`(?i)^https?://`).MatchString(processed) {
						break
					}
					continue
				}
			}
			if err != nil {
				os.Remove(fileName)
				if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode != 0 && te.StatusCode != 500 && te.StatusCode != 502 && te.StatusCode != 503 && te.StatusCode != 504 {
					break
				}
				if strings.HasPrefix(err.Error(), "Your configuration does not allow") {
					break
				}
				continue
			}
			break
		}
		if err != nil {
			return err
		}
		if !im.filesCache.ReadOnly {
			im.mu.Lock()
			im.cacheWrites[p.Name] = cacheKey
			im.mu.Unlock()
			im.filesCache.CopyFrom(cacheKey, fileName)
		}
		if err := verify(); err != nil {
			return err
		}
		// PluginEvents::POST_FILE_DOWNLOAD, dispatched once the transfers are done
		return im.pluginPostFileDownload(p, fileName, checksum, processed, true)
	}, nil
}

func (im *InstallationManager) clearLastCacheWrite(p *pkg.Package) {
	im.mu.Lock()
	key, ok := im.cacheWrites[p.Name]
	im.mu.Unlock()
	if ok {
		im.filesCache.Remove(key)
	}
}

func (im *InstallationManager) pathDownload(p *pkg.Package) error {
	u := p.DistURL()
	if u == "" {
		return &UserError{Msg: "The package " + p.PrettyName + " has no dist url configured, cannot download.", Code: 1}
	}
	realURL := fsutil.Realpath(im.absFromCwd(u))
	if realURL == "" {
		return &UserError{Msg: fmt.Sprintf(`Source path "%s" is not found for package %s`, u, p.Name), Code: 1}
	}
	if st, err := os.Stat(realURL); err != nil || !st.IsDir() {
		return &UserError{Msg: fmt.Sprintf(`Source path "%s" is not found for package %s`, u, p.Name), Code: 1}
	}
	target := im.InstallPath(p)
	if rp := fsutil.Realpath(target); rp != "" && rp == realURL {
		return nil
	}
	if rp := fsutil.Realpath(target); rp != "" && strings.HasPrefix(rp+"/", realURL+"/") {
		return &UserError{Msg: fmt.Sprintf(`Package %s cannot install to "%s" inside its source at "%s"`, p.Name, rp, realURL), Code: 1}
	}
	return nil
}

func (im *InstallationManager) absFromCwd(p string) string {
	if fsutil.IsAbsolutePath(p) {
		return p
	}
	return im.c.Cwd + "/" + p
}

// installCode is DownloadManager::install for the chosen downloader.
func (im *InstallationManager) installCode(p *pkg.Package, output bool) error {
	target := im.InstallPath(p)
	// a built-in downloader replaced by a plugin (plugin_services.go)
	if handled, err := im.replacedDownloaderOp(im.downloaderType(p), "install", p, nil); handled {
		return err
	}
	switch im.downloaderType(p) {
	case "":
		return nil
	case "path":
		return im.pathInstall(p, target, output)
	case "zip", "tar", "gzip", "xz", "rar":
		return im.archiveInstall(p, target, output)
	case "file", "phar":
		return im.fileInstall(p, target, output)
	case "git":
		return im.gitDL().install(p, target)
	case "svn", "hg", "fossil", "perforce":
		return im.sourceDL(im.downloaderType(p)).install(p, target)
	}
	// downloaders added by plugins (DownloadManager::setDownloader)
	if handled, err := im.pluginDownloaderOp(im.downloaderType(p), "install", p, nil); handled {
		return err
	}
	return &UserError{Msg: "Unknown downloader type: " + im.downloaderType(p), Code: 1}
}

func (im *InstallationManager) archiveInstall(p *pkg.Package, target string, output bool) error {
	return im.archiveInstallAs(im.downloaderType(p), p, target, output)
}

// archiveInstallAs extracts with the archive downloader kind (zip, tar,
// gzip): downloaders added by plugins for other dist types may extend them.
func (im *InstallationManager) archiveInstallAs(kind string, p *pkg.Package, target string, output bool) error {
	if output {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false) + ": Extracting archive")
	}
	vendorDir := im.realVendorDir()
	if !strings.Contains(fsutil.NormalizePath(vendorDir), fsutil.NormalizePath(target+"/")) {
		if err := emptyDirectory(target); err != nil {
			return err
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
		return err
	}
	fileName := im.tmpFileName(p)
	cleanup := func() {
		im.clearLastCacheWrite(p)
		os.RemoveAll(tempDir)
		if st, err := os.Stat(target); err == nil && st.IsDir() && fsutil.Realpath(target) != im.c.Cwd {
			os.RemoveAll(target)
		}
	}
	var err error
	switch kind {
	case "zip":
		// unzip (native), 7z or ZipArchive like ZipDownloader (archive_cmds.go)
		err = im.extractZipWith(p, fileName, tempDir)
	case "xz":
		err = im.extractXz(fileName, tempDir)
	case "rar":
		err = im.extractRar(fileName, tempDir)
	case "tar":
		err = ExtractTar(fileName, tempDir)
		if err != nil {
			err = &UserError{Msg: "Failed to extract " + p.Name + ": " + err.Error(), Code: 1}
		}
	case "gzip":
		name := distPathComponent(p, false)
		if i := strings.LastIndexByte(name, '.'); i >= 0 {
			name = name[:i]
		}
		err = ExtractGzip(fileName, filepath.Join(tempDir, name))
	}
	if err != nil {
		cleanup()
		return err
	}
	os.Remove(fileName)
	entries := folderContent(tempDir)
	renameAsOne := false
	if _, err := os.Lstat(target); os.IsNotExist(err) {
		renameAsOne = true
	} else if fsutil.IsDirEmpty(target) {
		if os.Remove(target) == nil {
			renameAsOne = true
		}
	}
	single := len(entries) == 1 && isDir(entries[0])
	if renameAsOne {
		from := tempDir
		if single {
			from = entries[0]
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
			cleanup()
			return err
		}
		if err := os.Rename(from, target); err != nil {
			cleanup()
			return err
		}
	} else {
		from := tempDir
		if single {
			from = entries[0]
		}
		if err := renameRecursively(p, from, target); err != nil {
			cleanup()
			return err
		}
	}
	os.RemoveAll(tempDir)
	return nil
}

func isDir(p string) bool {
	st, err := os.Lstat(p)
	return err == nil && st.IsDir()
}

func folderContent(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.Name() == ".DS_Store" {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out
}

func renameRecursively(p *pkg.Package, from, to string) error {
	for _, f := range folderContent(from) {
		dest := filepath.Join(to, filepath.Base(f))
		if isDir(dest) {
			if !isDir(f) {
				return &UserError{Msg: "Installing " + p.String() + " would lead to overwriting the " + to + "/" + filepath.Base(f) + " directory with a file from the package, invalid operation.", Code: 1}
			}
			if err := renameRecursively(p, f, dest); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(to, 0o777); err != nil {
			return err
		}
		os.Remove(dest)
		if err := os.Rename(f, dest); err != nil {
			return err
		}
	}
	return nil
}

func emptyDirectory(dir string) error {
	st, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o777)
	}
	if err != nil {
		return err
	}
	if st.Mode()&os.ModeSymlink != 0 {
		if err := os.Remove(dir); err != nil {
			return err
		}
		return os.MkdirAll(dir, 0o777)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (im *InstallationManager) fileInstall(p *pkg.Package, target string, output bool) error {
	if output {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false))
	}
	vendorDir := im.realVendorDir()
	if !strings.Contains(fsutil.NormalizePath(vendorDir), fsutil.NormalizePath(target+"/")) {
		if err := emptyDirectory(target); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(target, 0o777); err != nil {
		return err
	}
	base := distPathComponent(p, false)
	if err := os.Rename(im.tmpFileName(p), filepath.Join(target, base)); err != nil {
		return err
	}
	for _, bin := range p.Binaries() {
		bp := filepath.Join(target, bin)
		st, err := os.Stat(bp)
		if err != nil || st.Mode()&0o111 != 0 {
			continue
		}
		if !binInsidePackage(target, bp) {
			continue
		}
		_ = os.Chmod(bp, 0o777&^currentUmask())
	}
	return nil
}

func (im *InstallationManager) pathInstall(p *pkg.Package, target string, output bool) error {
	u := p.DistURL()
	realURL := fsutil.Realpath(im.absFromCwd(u))
	if realURL == "" {
		return &UserError{Msg: "Failed to realpath " + u, Code: 1}
	}
	if rp := fsutil.Realpath(target); rp == realURL {
		if output {
			im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false) + ": Source already present")
		}
		return nil
	}
	to := p.TransportOptions()
	relative := true
	var symlinkOpt *bool
	if to != nil {
		if v, ok := to.Get("relative"); ok {
			if b, ok := v.(bool); ok {
				relative = b
			}
		}
		if v, ok := to.Get("symlink"); ok {
			if b, ok := v.(bool); ok {
				symlinkOpt = &b
			}
		}
	}
	strategy := "symlink"
	allowMirror := true
	if v := os.Getenv("COMPOSER_MIRROR_PATH_REPOS"); v != "" && v != "0" {
		strategy = "mirror"
	}
	if symlinkOpt != nil {
		if *symlinkOpt {
			strategy = "symlink"
			allowMirror = false
		} else {
			strategy = "mirror"
		}
	}
	os.RemoveAll(target)
	if output {
		im.c.IO.WriteErrorV("  - "+(&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false)+": ", false, 32)
	}
	fallback := false
	if strategy == "symlink" {
		if output {
			im.c.IO.WriteErrorV("Symlinking from "+u, false, 32)
		}
		linkTarget := realURL + "/"
		if relative {
			sp, err := fsutil.FindShortestPath(target, realURL, false, true)
			if err == nil {
				linkTarget = sp + "/"
			}
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
			return err
		}
		if err := os.Symlink(linkTarget, target); err != nil {
			if !allowMirror {
				return &UserError{Msg: fmt.Sprintf(`Symlink from "%s" to "%s" failed!`, realURL, target), Code: 1}
			}
			if output {
				im.c.IO.WriteError("")
				im.c.IO.WriteError("    <error>Symlink failed, fallback to use mirroring!</error>")
			}
			strategy = "mirror"
			fallback = true
		}
	}
	if strategy == "mirror" {
		if output {
			prefix := ""
			if fallback {
				prefix = "    "
			}
			im.c.IO.WriteErrorV(prefix+"Mirroring from "+u, false, 32)
		}
		if err := mirrorDir(realURL, target); err != nil {
			return err
		}
	}
	if output {
		im.c.IO.WriteError("")
	}
	return nil
}

// mirrorDir copies a directory tree (Symfony Filesystem::mirror with the
// archivable files finder; VCS metadata directories are skipped).
func mirrorDir(src, dst string) error {
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if rel == "." {
			return os.MkdirAll(dst, 0o777)
		}
		base := filepath.Base(p)
		if info.IsDir() && (base == ".git" || base == ".svn" || base == ".hg") {
			return filepath.SkipDir
		}
		target := filepath.Join(dst, rel)
		if info.Mode()&os.ModeSymlink != 0 {
			l, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(l, target)
		}
		if info.IsDir() {
			return os.MkdirAll(target, 0o777)
		}
		return fsutil.CopyFile(p, target)
	})
}

// removeCode is FileDownloader/PathDownloader::remove.
func (im *InstallationManager) removeCode(p *pkg.Package, output bool) error {
	target := im.packageBasePath(p)
	if im.downloaderType(p) == "path" {
		u := p.DistURL()
		absP := target
		absD := im.absFromCwd(u)
		if fsutil.NormalizePath(absP) == fsutil.NormalizePath(absD) {
			if output {
				im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUninstall, Package: p}).Show(false) + ", source is still present in " + target)
			}
			return nil
		}
	}
	if im.downloaderType(p) == "" {
		return nil
	}
	if output {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUninstall, Package: p}).Show(false))
	}
	if err := fsutil.RemoveAll(target); err != nil {
		return &UserError{Msg: "Could not completely delete " + target + ", aborting.", Code: 1}
	}
	return nil
}

func (im *InstallationManager) cleanup(p *pkg.Package) {
	im.mu.Lock()
	fileName, ok := im.tmpFiles[p]
	im.mu.Unlock()
	if ok {
		os.Remove(fileName)
	}
	target := im.InstallPath(p)
	im.mu.Lock()
	if d, ok := im.downloadPaths[p]; ok {
		target = d
	}
	im.mu.Unlock()
	dirs := []string{target, im.vendorDir + "/" + strings.SplitN(p.PrettyName, "/", 2)[0], im.vendorDir + "/composer/", im.vendorDir}
	for _, d := range dirs {
		if isDir(d) && fsutil.IsDirEmpty(d) && fsutil.Realpath(d) != im.c.Cwd {
			os.Remove(d)
		}
	}
}

// Execute is InstallationManager::execute.
func (im *InstallationManager) Execute(repo *installed.Repository, ops []*solver.Operation, devMode, runScripts, downloadOnly bool) error {
	if im.Trace != nil {
		for _, op := range ops {
			*im.Trace = append(*im.Trace, cio.StripTags(op.Show(false)))
			switch op.Type {
			case solver.OpInstall:
				repo.AddPackage(op.Package.Clone())
			case solver.OpUpdate:
				repo.RemovePackage(op.Initial)
				if !repo.HasPackage(op.Target) {
					repo.AddPackage(op.Target.Clone())
				}
			case solver.OpUninstall:
				repo.RemovePackage(op.Package)
			case solver.OpMarkAliasInstalled:
				if !repo.HasPackage(op.Package) {
					repo.AddPackage(op.Package.Clone())
				}
			case solver.OpMarkAliasUninstalled:
				repo.RemovePackage(op.Package)
			}
		}
		return nil
	}
	defer im.beginExecute(repo, ops)()
	// batch around download-modifying plugins (they would run before other downloads)
	var batches [][]*solver.Operation
	var batch []*solver.Operation
	for _, op := range ops {
		if op.Type == solver.OpInstall || op.Type == solver.OpUpdate {
			p := op.MainPackage()
			if p.Type() == "composer-plugin" {
				if extra := p.Extra(); extra != nil {
					if v, ok := extra.Get("plugin-modifies-downloads"); ok && v == true {
						if len(batch) > 0 {
							batches = append(batches, batch)
						}
						batches = append(batches, []*solver.Operation{op})
						batch = nil
						continue
					}
				}
			}
		}
		batch = append(batch, op)
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	for _, b := range batches {
		if err := im.downloadAndExecuteBatch(repo, b, devMode, runScripts, downloadOnly, ops); err != nil {
			for _, op := range ops {
				if op.Type == solver.OpInstall || op.Type == solver.OpUpdate || op.Type == solver.OpUninstall {
					if op.MainPackage().InstallationSource() != "" {
						// the cleanup step of the download-time installer (plugin_install.go)
						im.cleanupAfterFailure(op)
					}
				}
			}
			return err
		}
	}
	if downloadOnly {
		return nil
	}
	return im.writeRepo(repo, devMode)
}

func (im *InstallationManager) writeRepo(repo *installed.Repository, devMode bool) error {
	return repo.Write(devMode, im.InstallPath, im.c.Package, im.c.Cwd)
}

func (im *InstallationManager) downloadAndExecuteBatch(repo *installed.Repository, ops []*solver.Operation, devMode, runScripts, downloadOnly bool, all []*solver.Operation) error {
	// the installers the cleanup steps use (plugin_install.go)
	if err := im.pluginRecordInstallers(ops); err != nil {
		return err
	}
	type job struct {
		p, prev *pkg.Package
	}
	var jobs []job
	for _, op := range ops {
		switch op.Type {
		case solver.OpInstall:
			jobs = append(jobs, job{op.Package, nil})
		case solver.OpUpdate:
			jobs = append(jobs, job{op.Target, op.Initial})
		}
	}
	// "Downloading" lines are printed in operation order like Composer; the
	// transfers themselves run concurrently.
	var transfers []func() error
	httpJobs := 0
	for _, j := range jobs {
		if err := pluginDownloadCheck(j.p); err != nil {
			return err
		}
		t, err := im.download(j.p, j.prev)
		if err != nil {
			return err
		}
		if t != nil {
			transfers = append(transfers, t)
			if reHTTPURL.MatchString(j.p.DistURL()) {
				httpJobs++
			}
		}
	}
	if len(jobs) > 0 {
		errs := make([]error, len(transfers))
		err := im.waitOnJobs(len(jobs), httpJobs, len(transfers), func(i int) error {
			errs[i] = transfers[i]()
			return errs[i]
		})
		if err != nil {
			return err
		}
		// Loop::wait runs every transfer to its end before rethrowing a
		// failure (the post-file-download events of the transfers were
		// dispatched as they completed, plugin_dlevents.go)
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}
	if downloadOnly {
		for _, j := range jobs {
			im.cleanup(j.p)
		}
		return nil
	}
	var batches [][]*solver.Operation
	var batch []*solver.Operation
	for _, op := range ops {
		if op.Type == solver.OpInstall || op.Type == solver.OpUpdate {
			t := op.MainPackage().Type()
			if t == "composer-plugin" || t == "composer-installer" {
				if len(batch) > 0 {
					batches = append(batches, batch)
				}
				batches = append(batches, []*solver.Operation{op})
				batch = nil
				continue
			}
		}
		batch = append(batch, op)
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	for _, b := range batches {
		if err := im.executeBatch(repo, b, devMode, runScripts); err != nil {
			return err
		}
	}
	return nil
}

func (im *InstallationManager) executeBatch(repo *installed.Repository, ops []*solver.Operation, devMode, runScripts bool) error {
	var post []*solver.Operation
	promises, procJobs := 0, 0
	im.deferredBins = nil
	for _, op := range ops {
		switch op.Type {
		case solver.OpInstall, solver.OpUpdate, solver.OpUninstall:
			promises++
			procJobs += im.asyncProcessJobs(op)
		}
	}
	for _, op := range ops {
		switch op.Type {
		case solver.OpMarkAliasInstalled:
			if im.c.IO.IsDebug() {
				im.c.IO.WriteError("  - " + op.Show(false))
			}
			if !repo.HasPackage(op.Package) {
				repo.AddPackage(op.Package.Clone())
			}
			continue
		case solver.OpMarkAliasUninstalled:
			if im.c.IO.IsDebug() {
				im.c.IO.WriteError("  - " + op.Show(false))
			}
			repo.RemovePackage(op.Package)
			continue
		}
		pre := map[string]string{solver.OpInstall: "pre-package-install", solver.OpUpdate: "pre-package-update", solver.OpUninstall: "pre-package-uninstall"}[op.Type]
		if runScripts && im.PackageEvent != nil {
			if err := im.PackageEvent(pre, op); err != nil {
				return err
			}
		}
		// DownloadManager::prepare throws synchronously in Composer, before
		// the promise chain whose rejection handler reports "X of Y failed"
		if err := im.pluginPrepare(op); err != nil {
			return err
		}
		// prepare of installers added by plugins (plugin_install.go)
		if handled, err := im.pluginInstallerStage(op, "prepare"); err != nil {
			return err
		} else if !handled {
			if err := im.prepareStage(op); err != nil {
				return err
			}
		}
		// operations of installers added by plugins are promise chains
		// settled with the batch (plugin_install.go)
		if started, err := im.pluginStartOp(repo, op, devMode); err != nil {
			verb := map[string]string{solver.OpInstall: "Install", solver.OpUpdate: "Update", solver.OpUninstall: "Uninstall"}[op.Type]
			im.c.IO.WriteError("    <error>" + verb + " of " + op.MainPackage().PrettyName + " failed</error>")
			return err
		} else if started {
			post = append(post, op)
			continue
		}
		// operations of downloaders replaced by plugins (plugin_services.go)
		if started, err := im.replacedStartOp(repo, op, devMode); err != nil {
			verb := map[string]string{solver.OpInstall: "Install", solver.OpUpdate: "Update", solver.OpUninstall: "Uninstall"}[op.Type]
			im.c.IO.WriteError("    <error>" + verb + " of " + op.MainPackage().PrettyName + " failed</error>")
			return err
		} else if started {
			post = append(post, op)
			continue
		}
		// installer chosen like InstallationManager::getInstaller (plugin_install.go)
		err := im.runOp(repo, op)
		handled := false
		if err == nil {
			handled, err = im.pluginInstallerStage(op, "cleanup")
		}
		if err == nil && !handled {
			err = im.cleanupStage(op)
		}
		if err == nil {
			if !handled {
				im.cleanup(op.MainPackage())
			}
			err = im.writeRepo(repo, devMode)
		}
		if err != nil {
			verb := map[string]string{solver.OpInstall: "Install", solver.OpUpdate: "Update", solver.OpUninstall: "Uninstall"}[op.Type]
			im.c.IO.WriteError("    <error>" + verb + " of " + op.MainPackage().PrettyName + " failed</error>")
			return err
		}
		post = append(post, op)
	}
	if err := im.pluginSettleOps(repo); err != nil {
		return err
	}
	if promises > 0 {
		// Composer starts all operations and then waits on their async
		// processes (unzip, rm -rf) with a progress bar; gocomposer performs
		// them in-process, so the bar is drawn once they completed.
		if err := im.waitOnJobs(promises, procJobs, 0, nil); err != nil {
			return err
		}
	}
	for _, f := range im.deferredBins {
		f()
	}
	im.deferredBins = nil
	if runScripts && im.PackageEvent != nil {
		for _, op := range post {
			ev := map[string]string{solver.OpInstall: "post-package-install", solver.OpUpdate: "post-package-update", solver.OpUninstall: "post-package-uninstall"}[op.Type]
			if err := im.PackageEvent(ev, op); err != nil {
				return err
			}
		}
	}
	return nil
}

func (im *InstallationManager) install(repo *installed.Repository, p *pkg.Package) error {
	target := im.InstallPath(p)
	if target != "" {
		if _, err := os.Stat(target); err != nil && repo.HasPackage(p) {
			im.removeBinaries(p)
		}
	}
	if p.Type() == "metapackage" {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false))
	} else if err := im.installCode(p, true); err != nil {
		return err
	}
	if target != "" {
		if im.asyncProcessJobs(&solver.Operation{Type: solver.OpInstall, Package: p}) > 0 {
			// Composer extracts zips with an async unzip process; the
			// binaries are installed in the promise callback, which runs
			// once Loop::wait collects the finished processes
			im.deferredBins = append(im.deferredBins, func() { im.installBinaries(p, target, true) })
		} else {
			im.installBinaries(p, target, true)
		}
	}
	if !repo.HasPackage(p) {
		repo.AddPackage(p.Clone())
	}
	return nil
}

func (im *InstallationManager) update(repo *installed.Repository, initial, target *pkg.Package) error {
	if !repo.HasPackage(initial) {
		return &UserError{Msg: "Package is not installed: " + initial.String(), Code: 1}
	}
	im.removeBinaries(initial)
	if initial.Type() == "metapackage" || target.Type() == "metapackage" {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUpdate, Initial: initial, Target: target}).Show(false))
		if initial.Type() != "metapackage" {
			if err := fsutil.RemoveAll(im.packageBasePath(initial)); err != nil {
				return err
			}
		}
		if target.Type() != "metapackage" {
			if err := im.installCode(target, false); err != nil {
				return err
			}
		}
	} else {
		initialPath := im.InstallPath(initial)
		targetPath := im.InstallPath(target)
		if initialPath != targetPath {
			if strings.HasPrefix(initialPath, targetPath) || strings.HasPrefix(targetPath, initialPath) {
				if err := im.removeCodeDM(initial); err != nil {
					return err
				}
				if err := im.installCode(target, true); err != nil {
					return err
				}
				goto done
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o777); err != nil {
				return err
			}
			if err := os.Rename(initialPath, targetPath); err != nil {
				return err
			}
		}
		if typ := im.downloaderType(target); typ == im.downloaderType(initial) && im.replacedDownloader(typ) {
			// a downloader replaced by a plugin (DownloadManager::update)
			if _, err := im.pluginDownloaderCall(typ, "update", target, initial, ""); err != nil {
				return err
			}
		} else if isVcsDownloaderType(im.downloaderType(target)) && im.downloaderType(initial) == im.downloaderType(target) {
			if err := im.sourceDL(im.downloaderType(target)).update(initial, target, targetPath); err != nil {
				return err
			}
		} else if im.downloaderType(initial) == im.downloaderType(target) {
			// downloaders added by plugins (DownloadManager::setDownloader)
			if handled, err := im.pluginDownloaderOp(im.downloaderType(target), "update", target, initial); handled {
				if err != nil {
					return err
				}
				goto done
			}
			appendix := ""
			switch im.downloaderType(target) {
			case "zip", "tar", "gzip", "xz", "rar":
				appendix = ": Extracting archive"
			case "path":
				appendix = im.pathAppendix(target)
			}
			im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUpdate, Initial: initial, Target: target}).Show(false) + appendix)
			if im.downloaderType(target) == "path" {
				if err := im.removeCode(initial, false); err != nil {
					return err
				}
				if err := im.pathInstall(target, targetPath, false); err != nil {
					return err
				}
			} else {
				if err := fsutil.RemoveAll(targetPath); err != nil {
					return err
				}
				if err := im.installCode(target, false); err != nil {
					return err
				}
			}
		} else {
			// on a type change the existing source install is about to be
			// wiped, so run its uninstall guard first
			if isVcsDownloaderType(im.downloaderType(initial)) {
				if err := im.sourceDL(im.downloaderType(initial)).prepare("uninstall", initial, nil, targetPath); err != nil {
					return err
				}
			}
			if err := im.removeCodeDM(initial); err != nil {
				return err
			}
			if err := im.installCode(target, true); err != nil {
				return err
			}
		}
	}
done:
	im.installBinaries(target, im.InstallPath(target), true)
	repo.RemovePackage(initial)
	if !repo.HasPackage(target) {
		repo.AddPackage(target.Clone())
	}
	return nil
}

func (im *InstallationManager) pathAppendix(p *pkg.Package) string {
	u := p.DistURL()
	realURL := fsutil.Realpath(im.absFromCwd(u))
	if fsutil.Realpath(im.InstallPath(p)) == realURL {
		return ": Source already present"
	}
	to := p.TransportOptions()
	strategy := "symlink"
	if v := os.Getenv("COMPOSER_MIRROR_PATH_REPOS"); v != "" && v != "0" {
		strategy = "mirror"
	}
	if to != nil {
		if v, ok := to.Get("symlink"); ok {
			if b, ok := v.(bool); ok {
				if b {
					strategy = "symlink"
				} else {
					strategy = "mirror"
				}
			}
		}
	}
	if strategy == "symlink" {
		return ": Symlinking from " + u
	}
	return ": Mirroring from " + u
}

func (im *InstallationManager) uninstall(repo *installed.Repository, p *pkg.Package) error {
	if !repo.HasPackage(p) {
		return &UserError{Msg: "Package is not installed: " + p.String(), Code: 1}
	}
	if p.Type() == "metapackage" {
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUninstall, Package: p}).Show(false))
	} else if err := im.removeCodeDM(p); err != nil {
		return err
	}
	im.removeBinaries(p)
	repo.RemovePackage(p)
	if strings.Contains(p.Name, "/") && p.Type() != "metapackage" {
		vendor := filepath.Dir(im.packageBasePath(p))
		if isDir(vendor) && fsutil.IsDirEmpty(vendor) {
			os.Remove(vendor)
		}
	}
	return nil
}

// EnsureBinariesPresence re-creates missing bin proxies.
func (im *InstallationManager) EnsureBinariesPresence(p *pkg.Package) {
	if p.Type() == "metapackage" {
		return
	}
	im.installBinaries(p, im.InstallPath(p), false)
}

var reHTTPURL = regexp.MustCompile(`(?i)^https?://`)

// progressBar mirrors the conditions of InstallationManager::waitOnPromises.
func (im *InstallationManager) progressBar(promises int) *cio.ProgressBar {
	if !im.OutputProgress || envTruthy("CI") || im.c.IO.IsDebug() || promises <= 1 {
		return nil
	}
	return im.c.IO.NewProgressBar(0)
}

func envTruthy(name string) bool {
	v, ok := os.LookupEnv(name)
	return ok && v != "" && v != "0"
}

// maxParallelHTTP is HttpDownloader's maximum number of concurrent
// transfers: 12, or COMPOSER_MAX_PARALLEL_HTTP clamped to 1..50.
func maxParallelHTTP() int {
	v := strings.TrimSpace(os.Getenv("COMPOSER_MAX_PARALLEL_HTTP"))
	if v == "" {
		return 12
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 12
	}
	n := int(f)
	if n < 1 {
		n = 1
	}
	if n > 50 {
		n = 50
	}
	return n
}

// waitOnJobs is InstallationManager::waitOnPromises + Loop::wait: it runs n
// jobs concurrently (at most maxParallelHTTP() at a time like the curl multi handle) while
// drawing the progress bar over activeJobs, the number of async jobs Composer
// would have in flight when the wait starts.
func (im *InstallationManager) waitOnJobs(promises, activeJobs, n int, run func(i int) error) error {
	bar := im.progressBar(promises)
	// ConsoleIO::writeProgress reports the percentage as program status
	onProgress := func() {
		if bar.MaxSteps() > 0 {
			pct := 100 * bar.Progress() / bar.MaxSteps()
			im.c.IO.WriteProgress(&pct)
		}
	}
	if bar != nil {
		bar.Start(activeJobs)
		// Loop::runJobs always updates the bar on its first iteration
		// (an indeterminate bar is redrawn with the next offset)
		bar.SetProgress(bar.MaxSteps() - activeJobs)
		onProgress()
	}
	if n > 0 {
		var wg sync.WaitGroup
		var mu sync.Mutex
		remaining := n
		sem := make(chan struct{}, maxParallelHTTP())
		done := make(chan struct{})
		// the jobs start in order as slots free up, like HttpDownloader's
		// queue (with COMPOSER_MAX_PARALLEL_HTTP=1 one after the other)
		wg.Add(n)
		go func() {
			for i := 0; i < n; i++ {
				sem <- struct{}{}
				go func(i int) {
					defer wg.Done()
					defer func() { <-sem }()
					run(i)
					mu.Lock()
					remaining--
					mu.Unlock()
				}(i)
			}
		}()
		go func() { wg.Wait(); close(done) }()
		tick := time.NewTicker(100 * time.Millisecond)
	loop:
		for {
			select {
			case <-done:
				break loop
			case <-tick.C:
				if bar != nil {
					mu.Lock()
					active := remaining
					if active > bar.MaxSteps() {
						active = bar.MaxSteps()
					}
					mu.Unlock()
					bar.SetProgress(bar.MaxSteps() - active)
					onProgress()
				}
			}
		}
		tick.Stop()
	}
	if bar != nil {
		bar.Finish()
		hundred := 100
		im.c.IO.WriteProgress(&hundred)
		im.c.IO.WriteProgress(nil)
		bar.Clear()
		// ProgressBar in non-decorated output does not output a final line-break
		if !im.c.IO.Decorated {
			im.c.IO.WriteError("")
		}
	}
	return nil
}

var unzipAvailable = sync.OnceValue(func() bool {
	if _, err := exec.LookPath("unzip"); err == nil {
		return true
	}
	if runtime.GOOS != "windows" {
		for _, c := range []string{"7z", "7zz", "7za"} {
			if _, err := exec.LookPath(c); err == nil {
				return true
			}
		}
	}
	return false
})

// asyncProcessJobs is the number of ProcessExecutor jobs Composer has queued
// for op when it starts waiting: unzip for zip extraction and rm -rf for
// removing an existing package directory of a file-based dist.
func (im *InstallationManager) asyncProcessJobs(op *solver.Operation) int {
	isFileDist := func(p *pkg.Package) bool {
		return p.InstallationSource() == "dist" && p.DistType() != "path"
	}
	existingDir := func(p *pkg.Package) bool {
		st, err := os.Lstat(im.InstallPath(p))
		return err == nil && st.IsDir() && st.Mode()&os.ModeSymlink == 0
	}
	switch op.Type {
	case solver.OpInstall:
		p := op.Package
		if p.Type() != "metapackage" && isFileDist(p) && p.DistType() == "zip" && unzipAvailable() {
			return 1
		}
	case solver.OpUpdate:
		if op.Initial.Type() != "metapackage" && isFileDist(op.Initial) && existingDir(op.Initial) {
			return 1
		}
	case solver.OpUninstall:
		// the rm -rf of a removed package has normally finished (and left
		// ProcessExecutor's job list) when Loop::wait counts the jobs
		return 0
	}
	return 0
}
