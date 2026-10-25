package composer

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Downloader/VcsDownloader.php, SvnDownloader.php,
// HgDownloader.php, FossilDownloader.php, PerforceDownloader.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
	"gocomposer/internal/vcs"
)

// sourceDownloader is the VcsDownloader API shared by the git, svn, hg,
// fossil and perforce downloaders.
type sourceDownloader interface {
	download(p, prev *pkg.Package) error
	prepare(typ string, p, prev *pkg.Package, path string) error
	cleanup(typ string, p, prev *pkg.Package, path string) error
	install(p *pkg.Package, path string) error
	update(initial, target *pkg.Package, path string) error
	LocalChanges(p *pkg.Package, path string) (string, error)
}

// isVcsDownloaderType reports whether typ names one of the VCS downloaders
// registered by Factory::createDownloadManager.
func isVcsDownloaderType(typ string) bool {
	switch typ {
	case "git", "svn", "hg", "fossil", "perforce":
		return true
	}
	return false
}

type vcsDLKey struct {
	im  *InstallationManager
	typ string
}

var vcsDownloaders sync.Map // vcsDLKey -> sourceDownloader

// sourceDL returns the downloader instance for a VCS type (one per
// installation manager like the DownloadManager's).
func (im *InstallationManager) sourceDL(typ string) sourceDownloader {
	if typ == "git" {
		return im.gitDL()
	}
	key := vcsDLKey{im, typ}
	if d, ok := vcsDownloaders.Load(key); ok {
		return d.(sourceDownloader)
	}
	b := &vcsDownloader{im: im, env: im.c.VcsEnv(), hasCleanedChanges: map[string]bool{}}
	switch typ {
	case "svn":
		b.impl = &svnDownloader{vcsDownloader: b, cacheCredentials: true}
	case "hg":
		b.impl = &hgDownloader{b}
	case "fossil":
		b.impl = &fossilDownloader{b}
	case "perforce":
		b.impl = &perforceDownloader{vcsDownloader: b}
	default:
		return nil
	}
	d, _ := vcsDownloaders.LoadOrStore(key, sourceDownloader(b))
	return d.(sourceDownloader)
}

// vcsImpl are the abstract methods of VcsDownloader.
type vcsImpl interface {
	doDownload(p, prev *pkg.Package, path, u string) error
	doInstall(p *pkg.Package, path, u string) error
	doUpdate(initial, target *pkg.Package, path, u string) error
	commitLogs(from, to, path string) (string, error)
	hasMetadataRepository(path string) bool
	// localChanges is getLocalChanges (ok false: null)
	localChanges(p *pkg.Package, path string) (string, bool, error)
	cleanChanges(p *pkg.Package, path string, update bool) error
}

// vcsDownloader is the abstract VcsDownloader.
type vcsDownloader struct {
	im                *InstallationManager
	env               *vcs.Env
	impl              vcsImpl
	hasCleanedChanges map[string]bool
}

func (d *vcsDownloader) io() *cio.IO        { return d.im.c.IO }
func (d *vcsDownloader) proc() *vcs.Process { return d.env.Process }

// loop runs fn over the prepared urls like VcsDownloader's while loops
// (an empty url, PHP's false from realpath(), ends it).
func (d *vcsDownloader) loop(urls []string, fn func(u string) error) error {
	for i, u := range urls {
		if u == "" {
			return nil
		}
		err := fn(u)
		if err == nil {
			return nil
		}
		remaining := len(urls) - i - 1
		if d.io().IsDebug() {
			d.io().WriteError("Failed: [" + vcsExceptionClass(err) + "] " + err.Error())
		} else if remaining > 0 {
			d.io().WriteError("    Failed, trying the next URL")
		}
		if remaining == 0 {
			return err
		}
	}
	return nil
}

func vcsExceptionClass(err error) string {
	switch err.(type) {
	case *vcs.SecurityError:
		return "Composer\\Exception\\SecurityException"
	case *vcs.PlainError:
		return "Exception"
	}
	return phpExceptionClass(err)
}

// download is VcsDownloader::download.
func (d *vcsDownloader) download(p, prev *pkg.Package) error {
	if p.SourceReference() == "" {
		return missingReference(p)
	}
	path := strings.TrimRight(d.im.InstallPath(p), "/")
	return d.loop(prepareUrls(sourceURLs(p)), func(u string) error {
		return d.impl.doDownload(p, prev, path, u)
	})
}

// prepare is VcsDownloader::prepare.
func (d *vcsDownloader) prepare(typ string, p, prev *pkg.Package, path string) error {
	switch typ {
	case "update":
		if err := d.impl.cleanChanges(prev, path, true); err != nil {
			return err
		}
		d.hasCleanedChanges[prev.UniqueName()] = true
	case "install":
		return emptyDirectory(path)
	case "uninstall":
		return d.impl.cleanChanges(p, path, false)
	}
	return nil
}

// cleanup is VcsDownloader::cleanup (none of these downloaders reapply
// anything).
func (d *vcsDownloader) cleanup(typ string, p, prev *pkg.Package, path string) error {
	if typ == "update" && d.hasCleanedChanges[prev.UniqueName()] {
		delete(d.hasCleanedChanges, prev.UniqueName())
	}
	return nil
}

// install is VcsDownloader::install.
func (d *vcsDownloader) install(p *pkg.Package, path string) error {
	if p.SourceReference() == "" {
		return missingReference(p)
	}
	d.io().WriteErrorV("  - "+(&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false)+": ", false, cio.Normal)
	return d.loop(prepareUrls(sourceURLs(p)), func(u string) error {
		return d.impl.doInstall(p, path, u)
	})
}

// update is VcsDownloader::update.
func (d *vcsDownloader) update(initial, target *pkg.Package, path string) error {
	if target.SourceReference() == "" {
		return missingReference(target)
	}
	d.io().WriteErrorV("  - "+(&solver.Operation{Type: solver.OpUpdate, Initial: initial, Target: target}).Show(false)+": ", false, cio.Normal)
	urls := prepareUrls(sourceURLs(target))
	var exception error
	for i, u := range urls {
		if u == "" {
			break
		}
		exception = d.impl.doUpdate(initial, target, path, u)
		if exception == nil {
			break
		}
		if d.io().IsDebug() {
			d.io().WriteError("Failed: [" + vcsExceptionClass(exception) + "] " + exception.Error())
		} else if i < len(urls)-1 {
			d.io().WriteError("    Failed, trying the next URL")
		}
	}
	if exception == nil && d.io().IsVerbose() && d.impl.hasMetadataRepository(path) {
		message := "Pulling in changes:"
		logs, err := d.impl.commitLogs(initial.SourceReference(), target.SourceReference(), path)
		if err != nil {
			return err
		}
		if strings.TrimSpace(logs) == "" {
			message = "Rolling back changes:"
			if logs, err = d.impl.commitLogs(target.SourceReference(), initial.SourceReference(), path); err != nil {
				return err
			}
		}
		if strings.TrimSpace(logs) != "" {
			lines := strings.Split(logs, "\n")
			for i, l := range lines {
				lines[i] = "      " + l
			}
			d.io().WriteError("    " + message)
			d.io().WriteError(strings.ReplaceAll(strings.Join(lines, "\n"), "<", `\<`))
		}
	}
	return exception
}

// LocalChanges is getLocalChanges ("" for null).
func (d *vcsDownloader) LocalChanges(p *pkg.Package, path string) (string, error) {
	s, _, err := d.impl.localChanges(p, path)
	return s, err
}

// defaultCleanChanges is VcsDownloader::cleanChanges.
func (d *vcsDownloader) defaultCleanChanges(p *pkg.Package, path string) error {
	_, ok, err := d.impl.localChanges(p, path)
	if err != nil {
		return err
	}
	if ok {
		return &vcs.RuntimeError{Msg: "Source directory " + path + " has uncommitted changes."}
	}
	return nil
}

// repoConfigOf is $package->getRepository()->getRepoConfig() when the
// package comes from a VcsRepository (nil otherwise).
func repoConfigOf(p *pkg.Package) *phpjson.Map {
	if p.AliasOf != nil && p.Repo == nil {
		p = p.AliasOf
	}
	if r, ok := p.Repo.(*repo.Vcs); ok {
		return r.RepoConfig()
	}
	return nil
}

// platformRealpath is Platform::realpath.
func platformRealpath(path string) (string, error) {
	rp := fsutil.Realpath(path)
	if rp == "" {
		return "", &vcs.RuntimeError{Msg: "Could not realpath() the path: " + path}
	}
	return rp, nil
}

// --- svn ---

type svnDownloader struct {
	*vcsDownloader
	cacheCredentials bool
}

func (d *svnDownloader) doDownload(p, prev *pkg.Package, path, u string) error {
	vcs.SvnCleanEnv()
	if vcs.NewSvn(u, d.env).BinaryVersion() == "" {
		return &vcs.RuntimeError{Msg: "svn was not found in your PATH, skipping source download"}
	}
	return nil
}

func (d *svnDownloader) doInstall(p *pkg.Package, path, u string) error {
	vcs.SvnCleanEnv()
	ref := p.SourceReference()
	if rc := repoConfigOf(p); rc != nil {
		if v, ok := rc.Get("svn-cache-credentials"); ok {
			d.cacheCredentials = phpjson.Truthy(v)
		}
	}
	d.io().WriteError(" Checking out " + ref)
	_, err := d.execute(p, u, []string{"svn", "co"}, u+"/"+ref, "", &path)
	return err
}

func (d *svnDownloader) doUpdate(initial, target *pkg.Package, path, u string) error {
	vcs.SvnCleanEnv()
	ref := target.SourceReference()
	if !d.hasMetadataRepository(path) {
		return &vcs.RuntimeError{Msg: "The .svn directory is missing from " + path + ", see https://getcomposer.org/commit-deps for more information"}
	}
	util := vcs.NewSvn(u, d.env)
	var flags []string
	if semver.VersionCompareOp(util.BinaryVersion(), "1.7.0", ">=") {
		flags = append(flags, "--ignore-ancestry")
	}
	d.io().WriteError(" Checking out " + ref)
	_, err := d.execute(target, u, append([]string{"svn", "switch"}, flags...), u+"/"+ref, path, nil)
	return err
}

var reSvnStatusChange = regexp.MustCompile(`(?m)^ *[^X ] +`)

func (d *svnDownloader) localChanges(p *pkg.Package, path string) (string, bool, error) {
	if !d.hasMetadataRepository(path) {
		return "", false, nil
	}
	_, out, err := d.proc().Execute([]string{"svn", "status", "--ignore-externals"}, path)
	if err != nil {
		return "", false, err
	}
	if reSvnStatusChange.MatchString(out) {
		return out, true, nil
	}
	return "", false, nil
}

// execute is SvnDownloader::execute.
func (d *svnDownloader) execute(p *pkg.Package, baseURL string, command []string, u, cwd string, path *string) (string, error) {
	util := vcs.NewSvn(baseURL, d.env)
	util.SetCacheCredentials(d.cacheCredentials)
	out, err := util.Execute(command, u, cwd, path, d.io().IsVerbose())
	if err != nil {
		return "", &vcs.RuntimeError{Msg: p.PrettyName + " could not be downloaded, " + err.Error()}
	}
	return out, nil
}

var reLineSplitTrim = regexp.MustCompile(`\s*\r?\n\s*`)

// cleanChanges is SvnDownloader::cleanChanges.
func (d *svnDownloader) cleanChanges(p *pkg.Package, path string, update bool) error {
	changes, ok, err := d.localChanges(p, path)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	io := d.io()
	if !io.Interactive {
		if d.im.c.Config.Get("discard-changes") == true {
			return d.discardChanges(path)
		}
		return d.defaultCleanChanges(p, path)
	}
	var lines []string
	for _, l := range reLineSplitTrim.Split(changes, -1) {
		lines = append(lines, "    "+l)
	}
	n := len(lines)
	plural := "s"
	if n == 1 {
		plural = ""
	}
	io.WriteError("    <error>" + p.PrettyName + " has modified file" + plural + ":</error>")
	shown := lines
	if n > 10 {
		shown = lines[:10]
	}
	for _, l := range shown {
		io.WriteError(l)
	}
	if n > 10 {
		rem := n - 10
		pl := "s"
		if rem == 1 {
			pl = ""
		}
		io.WriteError("    <info>" + strconv.Itoa(rem) + " more file" + pl + " modified, choose \"v\" to view the full list</info>")
	}
	action := "uninstall"
	if update {
		action = "update"
	}
	for {
		ans, err := io.Ask("    <info>Discard changes [y,n,v,?]?</info> ", "?")
		if err != nil {
			return err
		}
		s, _ := cio.AnswerString(ans)
		switch s {
		case "y":
			return d.discardChanges(path)
		case "n":
			return &vcs.RuntimeError{Msg: "Update aborted"}
		case "v":
			for _, l := range lines {
				io.WriteError(l)
			}
		default:
			for _, l := range []string{
				"    y - discard changes and apply the " + action,
				"    n - abort the " + action + " and let you manually clean things up",
				"    v - view modified files",
				"    ? - print help",
			} {
				io.WriteError(l)
			}
		}
	}
}

var (
	reSvnRev     = regexp.MustCompile(`@(\d+)$`)
	reSvnRevOnly = regexp.MustCompile(`.*@(\d+)$`)
	reSvnInfoURL = regexp.MustCompile(`<url>(.*)</url>`)
)

// commitLogs is SvnDownloader::getCommitLogs.
func (d *svnDownloader) commitLogs(from, to, path string) (string, error) {
	if !reSvnRev.MatchString(from) || !reSvnRev.MatchString(to) {
		return "Could not retrieve changes between " + from + " and " + to + " due to missing revision information", nil
	}
	command := []string{"svn", "info", "--non-interactive", "--xml", "--", path}
	code, out, err := d.proc().Execute(command, path)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	m := reSvnInfoURL.FindStringSubmatch(out)
	if m == nil {
		return "", &vcs.RuntimeError{Msg: "Unable to determine svn url for path " + path}
	}
	fromRev := reSvnRevOnly.ReplaceAllString(from, "$1")
	toRev := reSvnRevOnly.ReplaceAllString(to, "$1")
	command = []string{"svn", "log", "-r", fromRev + ":" + toRev, "--incremental"}
	util := vcs.NewSvn(m[1], d.env)
	util.SetCacheCredentials(d.cacheCredentials)
	logs, err := util.ExecuteLocal(command, path, "", d.io().IsVerbose())
	if err != nil {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + err.Error()}
	}
	return logs, nil
}

func (d *svnDownloader) discardChanges(path string) error {
	code, _, err := d.proc().Execute([]string{"svn", "revert", "-R", "."}, path)
	if err != nil {
		return err
	}
	if code != 0 {
		return &vcs.RuntimeError{Msg: "Could not reset changes\n\n:" + d.proc().ErrorOutput()}
	}
	return nil
}

func (d *svnDownloader) hasMetadataRepository(path string) bool { return vcsIsDir(path + "/.svn") }

// vcsIsDir is is_dir() (symlinks followed).
func vcsIsDir(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

// vcsIsFile is is_file().
func vcsIsFile(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}

// --- hg ---

type hgDownloader struct{ *vcsDownloader }

func (d *hgDownloader) doDownload(p, prev *pkg.Package, path, u string) error {
	if vcs.HgVersion(d.proc()) == "" {
		return &vcs.RuntimeError{Msg: "hg was not found in your PATH, skipping source download"}
	}
	return nil
}

func (d *hgDownloader) doInstall(p *pkg.Package, path, u string) error {
	if err := vcs.NewHg(d.env).RunCommand(func(u string) []string {
		return []string{"hg", "clone", "--", u, path}
	}, u, path); err != nil {
		return err
	}
	command := []string{"hg", "up", "--", p.SourceReference()}
	code, _, err := d.proc().Execute(command, fsutil.Realpath(path))
	if err != nil {
		return err
	}
	if code != 0 {
		return &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	return nil
}

func (d *hgDownloader) doUpdate(initial, target *pkg.Package, path, u string) error {
	hg := vcs.NewHg(d.env)
	ref := target.SourceReference()
	d.io().WriteError(" Updating to " + ref)
	if !d.hasMetadataRepository(path) {
		return &vcs.RuntimeError{Msg: "The .hg directory is missing from " + path + ", see https://getcomposer.org/commit-deps for more information"}
	}
	if err := hg.RunCommand(func(u string) []string { return []string{"hg", "pull", "--", u} }, u, path); err != nil {
		return err
	}
	return hg.RunCommand(func(string) []string { return []string{"hg", "up", "--", ref} }, u, path)
}

func (d *hgDownloader) localChanges(p *pkg.Package, path string) (string, bool, error) {
	if !vcsIsDir(path + "/.hg") {
		return "", false, nil
	}
	_, out, err := d.proc().Execute([]string{"hg", "st"}, fsutil.Realpath(path))
	if err != nil {
		return "", false, err
	}
	out = strings.TrimSpace(out)
	return out, out != "", nil
}

func (d *hgDownloader) cleanChanges(p *pkg.Package, path string, update bool) error {
	return d.defaultCleanChanges(p, path)
}

func (d *hgDownloader) commitLogs(from, to, path string) (string, error) {
	command := []string{"hg", "log", "-r", from + ":" + to, "--style", "compact"}
	code, out, err := d.proc().Execute(command, fsutil.Realpath(path))
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	return out, nil
}

func (d *hgDownloader) hasMetadataRepository(path string) bool { return vcsIsDir(path + "/.hg") }

// --- fossil ---

type fossilDownloader struct{ *vcsDownloader }

func (d *fossilDownloader) doDownload(p, prev *pkg.Package, path, u string) error { return nil }

func (d *fossilDownloader) execute(command []string, cwd string) (string, error) {
	code, out, err := d.proc().Execute(command, cwd)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	return out, nil
}

func (d *fossilDownloader) doInstall(p *pkg.Package, path, u string) error {
	if err := vcs.ProhibitURL(d.env, u); err != nil {
		return err
	}
	repoFile := path + ".fossil"
	realPath, err := platformRealpath(path)
	if err != nil {
		return err
	}
	d.io().WriteError("Cloning " + p.SourceReference())
	if _, err := d.execute([]string{"fossil", "clone", "--", u, repoFile}, ""); err != nil {
		return err
	}
	if _, err := d.execute([]string{"fossil", "open", "--nested", "--", repoFile}, realPath); err != nil {
		return err
	}
	_, err = d.execute([]string{"fossil", "update", "--", p.SourceReference()}, realPath)
	return err
}

func (d *fossilDownloader) doUpdate(initial, target *pkg.Package, path, u string) error {
	if err := vcs.ProhibitURL(d.env, u); err != nil {
		return err
	}
	d.io().WriteError(" Updating to " + target.SourceReference())
	if !d.hasMetadataRepository(path) {
		return &vcs.RuntimeError{Msg: "The .fslckout file is missing from " + path + ", see https://getcomposer.org/commit-deps for more information"}
	}
	realPath, err := platformRealpath(path)
	if err != nil {
		return err
	}
	if _, err := d.execute([]string{"fossil", "pull"}, realPath); err != nil {
		return err
	}
	_, err = d.execute([]string{"fossil", "up", "--", target.SourceReference()}, realPath)
	return err
}

func (d *fossilDownloader) localChanges(p *pkg.Package, path string) (string, bool, error) {
	if !d.hasMetadataRepository(path) {
		return "", false, nil
	}
	realPath, err := platformRealpath(path)
	if err != nil {
		return "", false, err
	}
	_, out, err := d.proc().Execute([]string{"fossil", "changes"}, realPath)
	if err != nil {
		return "", false, err
	}
	out = strings.TrimSpace(out)
	return out, out != "", nil
}

func (d *fossilDownloader) cleanChanges(p *pkg.Package, path string, update bool) error {
	return d.defaultCleanChanges(p, path)
}

func (d *fossilDownloader) commitLogs(from, to, path string) (string, error) {
	realPath, err := platformRealpath(path)
	if err != nil {
		return "", err
	}
	out, err := d.execute([]string{"fossil", "timeline", "-t", "ci", "-W", "0", "-n", "0", "before", to}, realPath)
	if err != nil {
		return "", err
	}
	re, rerr := regexp.Compile(`\d\d:\d\d:\d\d\s+\[` + to + `\]`)
	var log strings.Builder
	for _, line := range vcs.SplitLines(out) {
		if rerr == nil && re.MatchString(line) {
			break
		}
		log.WriteString(line)
	}
	return log.String(), nil
}

func (d *fossilDownloader) hasMetadataRepository(path string) bool {
	return vcsIsFile(path+"/.fslckout") || vcsIsFile(path+"/_FOSSIL_")
}

// --- perforce ---

type perforceDownloader struct {
	*vcsDownloader
	perforce *vcs.Perforce
}

func (d *perforceDownloader) doDownload(p, prev *pkg.Package, path, u string) error { return nil }

func (d *perforceDownloader) doInstall(p *pkg.Package, path, u string) error {
	ref := p.SourceReference()
	var label *string
	if i := strings.Index(ref, "@"); i >= 0 {
		l := ref[i+1:]
		label = &l
	}
	d.io().WriteError("Cloning " + ref)
	if err := d.initPerforce(p, path, u); err != nil {
		return err
	}
	d.perforce.SetStream(ref)
	if err := d.perforce.P4Login(); err != nil {
		return err
	}
	if err := d.perforce.WriteP4ClientSpec(); err != nil {
		return err
	}
	d.perforce.ConnectClient()
	if err := d.perforce.SyncCodeBase(label); err != nil {
		return err
	}
	return d.perforce.CleanupClientSpec()
}

// initPerforce is PerforceDownloader::initPerforce.
func (d *perforceDownloader) initPerforce(p *pkg.Package, path, u string) error {
	if d.perforce != nil {
		return d.perforce.InitializePath(path)
	}
	pf, err := vcs.NewPerforce(repoConfigOf(p), u, path, d.env)
	if err != nil {
		return err
	}
	d.perforce = pf
	return nil
}

func (d *perforceDownloader) doUpdate(initial, target *pkg.Package, path, u string) error {
	return d.doInstall(target, path, u)
}

func (d *perforceDownloader) localChanges(p *pkg.Package, path string) (string, bool, error) {
	d.io().WriteError("Perforce driver does not check for local changes before overriding")
	return "", false, nil
}

func (d *perforceDownloader) cleanChanges(p *pkg.Package, path string, update bool) error {
	return d.defaultCleanChanges(p, path)
}

func (d *perforceDownloader) commitLogs(from, to, path string) (string, error) {
	if d.perforce == nil {
		return "", nil
	}
	return d.perforce.CommitLogs(from, to)
}

func (d *perforceDownloader) hasMetadataRepository(string) bool { return true }
