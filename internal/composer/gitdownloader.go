package composer

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Downloader/VcsDownloader.php,
// src/Composer/Downloader/GitDownloader.php, the source part of
// src/Composer/Downloader/DownloadManager.php and Package::getSourceUrls).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
	"gocomposer/internal/vcs"
)

// gitDownloader is GitDownloader (with its VcsDownloader parent).
type gitDownloader struct {
	im                  *InstallationManager
	env                 *vcs.Env
	git                 *vcs.Git
	cachedPackages      map[*pkg.Package]map[string]bool
	hasStashedChanges   map[string]bool
	hasDiscardedChanges map[string]bool
	hasCleanedChanges   map[string]bool
}

func (im *InstallationManager) gitDL() *gitDownloader {
	im.mu.Lock()
	defer im.mu.Unlock()
	if im.gitDownloader == nil {
		env := im.c.VcsEnv()
		im.gitDownloader = &gitDownloader{im: im, env: env, git: env.Git(),
			cachedPackages: map[*pkg.Package]map[string]bool{}, hasStashedChanges: map[string]bool{},
			hasDiscardedChanges: map[string]bool{}, hasCleanedChanges: map[string]bool{}}
	}
	return im.gitDownloader
}

func (d *gitDownloader) io() *cio.IO         { return d.im.c.IO }
func (d *gitDownloader) proc() *vcs.Process  { return d.env.Process }
func (d *gitDownloader) cacheVcsDir() string { return d.im.c.Config.GetString("cache-vcs-dir") }

func (d *gitDownloader) cachePath(u string) string {
	return d.cacheVcsDir() + "/" + vcs.CacheKey(u) + "/"
}

var reNonGitURLChars = regexp.MustCompile(`(?i)[^a-z0-9_.-]`)
var reGitHubMirrorSrc = regexp.MustCompile(`^(?:(?:https?|git)://github\.com/|git@github\.com:)([^/]+)/(.+?)(?:\.git)?$`)
var reBitbucketMirrorSrc = regexp.MustCompile(`^https://bitbucket\.org/([^/]+)/(.+?)(?:\.git)?/?$`)

// sourceURLs is Package::getSourceUrls.
func sourceURLs(p *pkg.Package) []string {
	u := p.SourceURL()
	if u == "" {
		return nil
	}
	urls := []string{u}
	for _, mv := range p.SourceMirrors() {
		m := phpjson.ToMap(mv)
		if m == nil {
			continue
		}
		var mirrorURL string
		switch p.SourceType() {
		case "git":
			norm := u
			if mm := reGitHubMirrorSrc.FindStringSubmatch(u); mm != nil {
				norm = "gh-" + mm[1] + "/" + mm[2]
			} else if mm := reBitbucketMirrorSrc.FindStringSubmatch(u); mm != nil {
				norm = "bb-" + mm[1] + "/" + mm[2]
			} else {
				norm = reNonGitURLChars.ReplaceAllString(strings.Trim(u, "/"), "-")
			}
			mirrorURL = strings.NewReplacer("%package%", p.Name, "%normalizedUrl%", norm, "%type%", "git").Replace(m.Str("url"))
		default:
			continue
		}
		if containsStr(urls, mirrorURL) {
			continue
		}
		if pref, _ := m.Get("preferred"); phpjson.Truthy(pref) {
			urls = append([]string{mirrorURL}, urls...)
		} else {
			urls = append(urls, mirrorURL)
		}
	}
	return urls
}

// prepareUrls is VcsDownloader::prepareUrls: local paths are realpath()ed
// ("" when they do not exist, which ends the url loop like PHP's false).
func prepareUrls(urls []string) []string {
	out := make([]string, len(urls))
	for i, u := range urls {
		if vcs.IsLocalPath(u) {
			isFile := false
			if strings.HasPrefix(u, "file://") {
				u = u[len("file://"):]
				isFile = true
			}
			if strings.Contains(u, "%") {
				u = rawURLDecode(u)
			}
			rp := fsutil.Realpath(u)
			if rp != "" && isFile {
				rp = "file://" + rp
			}
			u = rp
		}
		out[i] = u
	}
	return out
}

func rawURLDecode(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && isHexByte(s[i+1]) && isHexByte(s[i+2]) {
			b.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func isHexByte(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func unhex(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	}
	return c - 'A' + 10
}

// urlLoop runs fn for each prepared url like VcsDownloader's while loops.
func (d *gitDownloader) urlLoop(urls []string, fn func(u string) error) error {
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
			d.io().WriteError("Failed: [" + phpExceptionClass(err) + "] " + err.Error())
		} else if remaining > 0 {
			d.io().WriteError("    Failed, trying the next URL")
		}
		if remaining == 0 {
			return err
		}
	}
	return nil
}

func phpExceptionClass(err error) string {
	switch err.(type) {
	case *vcs.InvalidArgumentError:
		return "InvalidArgumentException"
	}
	return "RuntimeException"
}

func missingReference(p *pkg.Package) error {
	return &invalidArgumentError{"Package " + p.PrettyName + " is missing reference information"}
}

// invalidArgumentError is an \InvalidArgumentException (not a
// RuntimeException: DownloadManager does not fall back to another source).
type invalidArgumentError struct{ msg string }

func (e *invalidArgumentError) Error() string { return e.msg }

// download is VcsDownloader::download + GitDownloader::doDownload.
func (d *gitDownloader) download(p, prev *pkg.Package) error {
	if p.SourceReference() == "" {
		return missingReference(p)
	}
	return d.urlLoop(prepareUrls(sourceURLs(p)), func(u string) error {
		if vcs.IsLocalPath(u) {
			return nil
		}
		vcs.CleanEnv(d.proc())
		cachePath := d.cachePath(u)
		gitVersion := vcs.GitVersion(d.proc())
		if gitVersion != "" && semver.VersionCompareOp(gitVersion, "2.3.0-rc0", ">=") && vcs.CacheUsable(cachePath) {
			d.io().WriteError("  - Syncing <info>" + p.Name + "</info> (<comment>" + p.FullPrettyVersion(true) + "</comment>) into cache")
			d.io().WriteErrorV("    Cloning to cache at "+cachePath, true, cio.Debug)
			ref := p.SourceReference()
			ok, err := d.git.FetchRefOrSyncMirror(u, cachePath, ref, p.PrettyVersion())
			if err != nil {
				return err
			}
			if ok && isDir(cachePath) {
				if d.cachedPackages[p] == nil {
					d.cachedPackages[p] = map[string]bool{}
				}
				d.cachedPackages[p][ref] = true
			}
		} else if gitVersion == "" {
			return &vcs.RuntimeError{Msg: "git was not found in your PATH, skipping source download"}
		}
		return nil
	})
}

func (d *gitDownloader) shortHash(ref string) string {
	if !d.io().IsVerbose() && reSha40.MatchString(ref) {
		return ref[:10]
	}
	return ref
}

var reSha40 = regexp.MustCompile(`^[0-9a-f]{40}$`)

func networkDisabled() bool {
	v := os.Getenv("COMPOSER_DISABLE_NETWORK")
	return v != "" && v != "0"
}

// prepare is VcsDownloader::prepare.
func (d *gitDownloader) prepare(typ string, p, prev *pkg.Package, path string) error {
	switch typ {
	case "update":
		if err := d.cleanChanges(prev, path, true); err != nil {
			return err
		}
		d.hasCleanedChanges[prev.UniqueName()] = true
	case "install":
		return emptyDirectory(path)
	case "uninstall":
		return d.cleanChanges(p, path, false)
	}
	return nil
}

// cleanup is VcsDownloader::cleanup.
func (d *gitDownloader) cleanup(typ string, p, prev *pkg.Package, path string) error {
	if typ == "update" && d.hasCleanedChanges[prev.UniqueName()] {
		delete(d.hasCleanedChanges, prev.UniqueName())
		return d.reapplyChanges(path)
	}
	return nil
}

// install is VcsDownloader::install + GitDownloader::doInstall.
func (d *gitDownloader) install(p *pkg.Package, path string) error {
	if p.SourceReference() == "" {
		return missingReference(p)
	}
	d.io().WriteErrorV("  - "+(&solver.Operation{Type: solver.OpInstall, Package: p}).Show(false)+": ", false, cio.Normal)
	return d.urlLoop(prepareUrls(sourceURLs(p)), func(u string) error {
		vcs.CleanEnv(d.proc())
		cachePath := d.cachePath(u)
		ref := p.SourceReference()
		var msg string
		var commands [][]string
		if d.cachedPackages[p][ref] {
			msg = "Cloning " + d.shortHash(ref) + " from cache"
			cloneFlags := []string{"--dissociate", "--reference", cachePath}
			if to := p.TransportOptions(); to != nil {
				if g := to.Map("git"); g != nil {
					if v, _ := g.Get("single_use_clone"); phpjson.Truthy(v) {
						cloneFlags = nil
					}
				}
			}
			commands = [][]string{
				append([]string{"git", "clone", "--no-checkout", cachePath, path}, cloneFlags...),
				{"git", "remote", "set-url", "origin", "--", "%sanitizedUrl%"},
				{"git", "remote", "add", "composer", "--", "%sanitizedUrl%"},
			}
		} else {
			msg = "Cloning " + d.shortHash(ref)
			commands = [][]string{
				{"git", "clone", "--no-checkout", "--", "%url%", path},
				{"git", "remote", "add", "composer", "--", "%url%"},
				{"git", "fetch", "composer"},
				{"git", "remote", "set-url", "origin", "--", "%sanitizedUrl%"},
				{"git", "remote", "set-url", "composer", "--", "%sanitizedUrl%"},
			}
			if networkDisabled() {
				return &vcs.RuntimeError{Msg: "The required git reference for " + p.Name + " is not in cache and network is disabled, aborting"}
			}
		}
		d.io().WriteError(msg)
		if _, err := d.git.RunCommands(commands, u, path, true); err != nil {
			return err
		}
		if su := p.SourceURL(); u != su && su != "" {
			d.updateOriginURL(path, su)
		} else {
			d.setPushURL(path, u)
		}
		return d.updateToCommit(p, path, ref, p.PrettyVersion())
	})
}

// update is VcsDownloader::update + GitDownloader::doUpdate.
func (d *gitDownloader) update(initial, target *pkg.Package, path string) error {
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
		exception = d.doUpdate(initial, target, path, u)
		if exception == nil {
			break
		}
		if d.io().IsDebug() {
			d.io().WriteError("Failed: [" + phpExceptionClass(exception) + "] " + exception.Error())
		} else if i < len(urls)-1 {
			d.io().WriteError("    Failed, trying the next URL")
		}
	}
	if exception == nil && d.io().IsVerbose() && isDir(path+"/.git") {
		message := "Pulling in changes:"
		logs, err := d.commitLogs(initial.SourceReference(), target.SourceReference(), path)
		if err != nil {
			return err
		}
		if strings.TrimSpace(logs) == "" {
			message = "Rolling back changes:"
			if logs, err = d.commitLogs(target.SourceReference(), initial.SourceReference(), path); err != nil {
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

func (d *gitDownloader) doUpdate(initial, target *pkg.Package, path, u string) error {
	vcs.CleanEnv(d.proc())
	if !isDir(path + "/.git") {
		return &vcs.RuntimeError{Msg: "The .git directory is missing from " + path + ", see https://getcomposer.org/commit-deps for more information"}
	}
	cachePath := d.cachePath(u)
	ref := target.SourceReference()
	var msg, remoteURL string
	if d.cachedPackages[target][ref] {
		msg = "Checking out " + d.shortHash(ref) + " from cache"
		remoteURL = cachePath
	} else {
		msg = "Checking out " + d.shortHash(ref)
		remoteURL = "%url%"
		if networkDisabled() {
			return &vcs.RuntimeError{Msg: "The required git reference for " + target.Name + " is not in cache and network is disabled, aborting"}
		}
	}
	d.io().WriteError(msg)
	if code, _, err := d.proc().Execute([]string{"git", "rev-parse", "--quiet", "--verify", ref + "^{commit}"}, path); err != nil {
		return err
	} else if code != 0 {
		if _, err := d.git.RunCommands([][]string{
			{"git", "remote", "set-url", "composer", "--", remoteURL},
			{"git", "fetch", "composer"},
			{"git", "fetch", "--tags", "composer"},
		}, u, path, false); err != nil {
			return err
		}
	}
	if _, err := d.git.RunCommands([][]string{{"git", "remote", "set-url", "composer", "--", "%sanitizedUrl%"}}, u, path, false); err != nil {
		return err
	}
	if err := d.updateToCommit(target, path, ref, target.PrettyVersion()); err != nil {
		return err
	}
	updateOrigin := false
	if code, out, _ := d.proc().Execute([]string{"git", "remote", "-v"}, path); code == 0 {
		om := reOriginURL.FindStringSubmatch(out)
		cm := reComposerURL.FindStringSubmatch(out)
		if om != nil && cm != nil && om[1] == cm[1] && cm[1] != target.SourceURL() {
			updateOrigin = true
		}
	}
	if updateOrigin && target.SourceURL() != "" {
		d.updateOriginURL(path, target.SourceURL())
	}
	return nil
}

var (
	reOriginURL   = regexp.MustCompile(`(?m)^origin\s+(\S+)`)
	reComposerURL = regexp.MustCompile(`(?m)^composer\s+(\S+)`)
)

// LocalChanges is GitDownloader::getLocalChanges ("" when clean).
func (d *gitDownloader) LocalChanges(p *pkg.Package, path string) (string, error) {
	vcs.CleanEnv(d.proc())
	if !isDir(path + "/.git") {
		return "", nil
	}
	cmd := []string{"git", "status", "--porcelain", "--untracked-files=no"}
	code, out, err := d.proc().Execute(cmd, path)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(cmd, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	return strings.TrimSpace(out), nil
}

var reHeadRef = regexp.MustCompile(`(?mi)^([a-f0-9]+) HEAD$`)

// UnpushedChanges is GitDownloader::getUnpushedChanges ("" when none).
func (d *gitDownloader) UnpushedChanges(p *pkg.Package, path string) (string, error) {
	vcs.CleanEnv(d.proc())
	if !isDir(path + "/.git") {
		return "", nil
	}
	showRef := func() (string, error) {
		cmd := []string{"git", "show-ref", "--head", "-d"}
		code, out, err := d.proc().Execute(cmd, path)
		if err != nil {
			return "", err
		}
		if code != 0 {
			return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(cmd, " ") + "\n\n" + d.proc().ErrorOutput()}
		}
		return strings.TrimSpace(out), nil
	}
	refs, err := showRef()
	if err != nil {
		return "", err
	}
	hm := reHeadRef.FindStringSubmatch(refs)
	if hm == nil {
		return "", nil
	}
	headRef := hm[1]
	bms := regexp.MustCompile(`(?mi)^`+regexp.QuoteMeta(headRef)+` refs/heads/(.+)$`).FindAllStringSubmatch(refs, -1)
	if len(bms) == 0 {
		return "", nil
	}
	var candidates []string
	for _, m := range bms {
		candidates = append(candidates, m[1])
	}
	branch := candidates[0]
	var unpushed *string
	branchNotFound := false
	for i := 0; i <= 1; i++ {
		var remoteBranches []string
		for _, c := range candidates {
			ms := regexp.MustCompile(`(?mi)^[a-f0-9]+ refs/remotes/((?:[^/]+)/`+regexp.QuoteMeta(c)+`)$`).FindAllStringSubmatch(refs, -1)
			if len(ms) > 0 {
				for _, m := range ms {
					branch = c
					remoteBranches = append(remoteBranches, m[1])
				}
				break
			}
		}
		if len(remoteBranches) == 0 {
			s := "Branch " + branch + " could not be found on any remote and appears to be unpushed"
			unpushed = &s
			branchNotFound = true
		} else {
			if branchNotFound {
				unpushed = nil
			}
			for _, rb := range remoteBranches {
				cmd := []string{"git", "diff", "--name-status", rb + "..." + branch, "--"}
				code, out, err := d.proc().Execute(cmd, path)
				if err != nil {
					return "", err
				}
				if code != 0 {
					return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(cmd, " ") + "\n\n" + d.proc().ErrorOutput()}
				}
				out = strings.TrimSpace(out)
				if unpushed == nil || len(out) < len(*unpushed) {
					o := out
					unpushed = &o
				}
			}
		}
		if unpushed != nil && *unpushed != "" && i == 0 {
			d.proc().Execute([]string{"git", "fetch", "--all"}, path)
			if refs, err = showRef(); err != nil {
				return "", err
			}
		}
		if unpushed == nil || *unpushed == "" {
			break
		}
	}
	if unpushed == nil {
		return "", nil
	}
	return *unpushed, nil
}

// cleanChanges is GitDownloader::cleanChanges: without interaction the
// discard-changes config decides, otherwise the user is asked what to do
// with the local modifications (y/n/v/d/s/?).
func (d *gitDownloader) cleanChanges(p *pkg.Package, path string, update bool) error {
	vcs.CleanEnv(d.proc())
	unpushed, err := d.UnpushedChanges(p, path)
	if err != nil {
		return err
	}
	discard := d.im.c.Config.Get("discard-changes")
	if unpushed != "" && (d.io().Interactive || discard != true) {
		return &vcs.RuntimeError{Msg: "Source directory " + path + " has unpushed changes on the current branch: \n" + unpushed}
	}
	changes, err := d.LocalChanges(p, path)
	if err != nil {
		return err
	}
	if changes == "" {
		return nil
	}
	if !d.io().Interactive {
		if discard == true {
			return d.discardChanges(path)
		}
		if s, ok := discard.(string); ok && s == "stash" {
			if !update {
				return &vcs.RuntimeError{Msg: "Source directory " + path + " has uncommitted changes."}
			}
			return d.stashChanges(path)
		}
		return &vcs.RuntimeError{Msg: "Source directory " + path + " has uncommitted changes."}
	}
	var lines []string
	for _, l := range reChangeSplit.Split(changes, -1) {
		lines = append(lines, "    "+l)
	}
	io := d.io()
	io.WriteError("    <error>" + p.PrettyName + " has modified files:</error>")
	if len(lines) > 10 {
		io.WriteError(strings.Join(lines[:10], "\n"))
		io.WriteError("    <info>" + strconv.Itoa(len(lines)-10) + " more files modified, choose \"v\" to view the full list</info>")
	} else {
		io.WriteError(strings.Join(lines, "\n"))
	}
	choices, action := "y,n,v,d,", "uninstall"
	if update {
		choices, action = "y,n,v,d,s,", "update"
	}
	for {
		ans, err := io.Ask("    <info>Discard changes ["+choices+"?]?</info> ", "?")
		if err != nil {
			return err
		}
		s, _ := cio.AnswerString(ans)
		switch {
		case s == "y":
			return d.discardChanges(path)
		case s == "s" && update:
			return d.stashChanges(path)
		case s == "n":
			return &vcs.RuntimeError{Msg: "Update aborted"}
		case s == "v":
			io.WriteError(strings.Join(lines, "\n"))
		case s == "d":
			if err := d.viewDiff(path); err != nil {
				return err
			}
		default:
			help := []string{
				"    y - discard changes and apply the " + action,
				"    n - abort the " + action + " and let you manually clean things up",
				"    v - view modified files",
				"    d - view local modifications (diff)",
			}
			io.WriteError(strings.Join(help, "\n"))
			if update {
				io.WriteError("    s - stash changes and try to reapply them after the update")
			}
			io.WriteError("    ? - print help")
		}
	}
}

var reChangeSplit = regexp.MustCompile(`\s*\r?\n\s*`)

// viewDiff is GitDownloader::viewDiff.
func (d *gitDownloader) viewDiff(path string) error {
	code, out, err := d.proc().Execute([]string{"git", "diff", "HEAD"}, path)
	if err != nil {
		return err
	}
	if code != 0 {
		return &vcs.RuntimeError{Msg: "Could not view diff\n\n:" + out}
	}
	d.io().WriteError(out)
	return nil
}

func (d *gitDownloader) reapplyChanges(path string) error {
	if d.hasStashedChanges[path] {
		delete(d.hasStashedChanges, path)
		d.io().WriteError("    <info>Re-applying stashed changes</info>")
		if code, _, err := d.proc().Execute([]string{"git", "stash", "pop"}, path); err != nil {
			return err
		} else if code != 0 {
			return &vcs.RuntimeError{Msg: "Failed to apply stashed changes:\n\n" + d.proc().ErrorOutput()}
		}
	}
	delete(d.hasDiscardedChanges, path)
	return nil
}

func (d *gitDownloader) discardChanges(path string) error {
	for _, cmd := range [][]string{{"git", "clean", "-df"}, {"git", "reset", "--hard"}} {
		code, out, err := d.proc().Execute(cmd, path)
		if err != nil {
			return err
		}
		if code != 0 {
			return &vcs.RuntimeError{Msg: "Could not reset changes\n\n:" + out}
		}
	}
	d.hasDiscardedChanges[path] = true
	return nil
}

func (d *gitDownloader) stashChanges(path string) error {
	code, out, err := d.proc().Execute([]string{"git", "stash", "--include-untracked"}, path)
	if err != nil {
		return err
	}
	if code != 0 {
		return &vcs.RuntimeError{Msg: "Could not stash changes\n\n:" + out}
	}
	d.hasStashedChanges[path] = true
	return nil
}

// updateToCommit is GitDownloader::updateToCommit.
func (d *gitDownloader) updateToCommit(p *pkg.Package, path, reference, prettyVersion string) error {
	var force []string
	if d.hasDiscardedChanges[path] || d.hasStashedChanges[path] {
		force = []string{"-f"}
	}
	branch := reDevAffixBranch.ReplaceAllString(prettyVersion, "")
	var output string
	execute := func(cmd []string) bool {
		code, out, err := d.proc().Execute(cmd, path)
		output = out
		return err == nil && code == 0
	}
	var branches *string
	if execute([]string{"git", "branch", "-r"}) {
		b := output
		branches = &b
	}
	gitRef := reference
	if !reSha40.MatchString(reference) && branches != nil &&
		regexp.MustCompile(`(?m)^\s+composer/`+regexp.QuoteMeta(reference)+`$`).MatchString(*branches) {
		cmd1 := append(append([]string{"git", "checkout"}, force...), "-B", branch, "composer/"+reference, "--")
		cmd2 := []string{"git", "reset", "--hard", "composer/" + reference, "--"}
		if execute(cmd1) && execute(cmd2) {
			return nil
		}
	}
	if reSha40.MatchString(reference) {
		if branches != nil && !regexp.MustCompile(`(?m)^\s+composer/`+regexp.QuoteMeta(branch)+`$`).MatchString(*branches) &&
			regexp.MustCompile(`(?m)^\s+composer/v`+regexp.QuoteMeta(branch)+`$`).MatchString(*branches) {
			branch = "v" + branch
		}
		command := []string{"git", "checkout", branch, "--"}
		fallback := append(append([]string{"git", "checkout"}, force...), "-B", branch, "composer/"+branch, "--")
		reset := []string{"git", "reset", "--hard", reference, "--"}
		if (execute(command) || execute(fallback)) && execute(reset) {
			return nil
		}
	}
	cmd1 := append(append([]string{"git", "checkout"}, force...), gitRef, "--")
	cmd2 := []string{"git", "reset", "--hard", gitRef, "--"}
	if execute(cmd1) && execute(cmd2) {
		return nil
	}
	extra := ""
	if strings.Contains(d.proc().ErrorOutput(), reference) {
		d.io().WriteError("    <warning>" + reference + " is gone (history was rewritten?)</warning>")
		what := "the tag was recreated"
		if p.IsDev() {
			what = "the commit was removed from the branch"
		}
		extra = "\nIt looks like the commit hash is not available in the repository, maybe " + what + "? Run \"composer update " + p.PrettyName + "\" to resolve this."
	}
	return &vcs.RuntimeError{Msg: vcs.Sanitize("Failed to execute " + strings.Join(cmd1, " ") + " && " + strings.Join(cmd2, " ") + "\n\n" + d.proc().ErrorOutput() + extra)}
}

var reDevAffixBranch = regexp.MustCompile(`(?i)(?:^dev-|(?:\.x)?-dev$)`)

func (d *gitDownloader) updateOriginURL(path, u string) {
	d.proc().Execute([]string{"git", "remote", "set-url", "origin", "--", u}, path)
	d.setPushURL(path, u)
}

func (d *gitDownloader) setPushURL(path, u string) {
	re := regexp.MustCompile(`^(?:https?|git)://` + vcs.GitHubDomainsRegex(d.im.c.Config) + `/([^/]+)/([^/]+?)(?:\.git)?$`)
	m := re.FindStringSubmatch(u)
	if m == nil {
		return
	}
	pushURL := "git@" + m[1] + ":" + m[2] + "/" + m[3] + ".git"
	protocols := []string{}
	phpjson.Each(d.im.c.Config.Get("github-protocols"), func(_ string, v phpjson.Value) {
		if s, ok := v.(string); ok {
			protocols = append(protocols, s)
		}
	})
	if !containsStr(protocols, "ssh") {
		pushURL = "https://" + m[1] + "/" + m[2] + "/" + m[3] + ".git"
	}
	d.proc().Execute([]string{"git", "remote", "set-url", "--push", "origin", "--", pushURL}, path)
}

func (d *gitDownloader) commitLogs(from, to, path string) (string, error) {
	cmd := vcs.RevListCommand(d.proc(), append([]string{"--format=%h - %an: %s", from + ".." + to}, vcs.NoShowSignatureFlags(d.proc())...)...)
	code, out, err := d.proc().Execute(cmd, path)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", &vcs.RuntimeError{Msg: "Failed to execute " + strings.Join(cmd, " ") + "\n\n" + d.proc().ErrorOutput()}
	}
	return vcs.ParseRevListOutput(d.proc(), out), nil
}

// VcsReference is VcsDownloader::getVcsReference (the commit of the
// checkout at path).
func (d *gitDownloader) VcsReference(p *pkg.Package, path string) string {
	cfg := phpjson.NewMap()
	cfg.Set("version", p.PrettyVersion())
	if res, ok := guessVersionCommit(path, cfg, d.im.c); ok {
		return res
	}
	return ""
}

func guessVersionCommit(path string, cfg *phpjson.Map, c *Composer) (string, bool) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	code, out, _ := vcs.NewProcess(nil, 0).Execute([]string{"git", "log", "--pretty=%H", "-n1", "HEAD"}, abs)
	if code != 0 {
		return "", false
	}
	if s := strings.TrimSpace(out); s != "" {
		return s, true
	}
	return "", false
}

// isRuntimeException reports whether err would be a PHP \RuntimeException
// (DownloadManager only reports/falls back for those).
func isRuntimeException(err error) bool {
	switch err.(type) {
	case *invalidArgumentError, *vcs.InvalidArgumentError:
		return false
	}
	return true
}

func unsupportedVcsDownloader(p *pkg.Package, typ string) error {
	return &vcs.RuntimeError{Msg: "gocomposer cannot install " + p.PrettyName + " from " + typ + " source: the " + typ + " downloader is not supported (only git sources and dist archives are)"}
}

// prepareOp is DownloadManager::prepare for the operation's package (only
// the VCS downloaders have a prepare step with side effects).
func (im *InstallationManager) prepareOp(op *solver.Operation) error {
	switch op.Type {
	case solver.OpInstall:
		if op.Package.Type() != "metapackage" && im.downloaderType(op.Package) == "git" {
			return im.gitDL().prepare("install", op.Package, nil, strings.TrimRight(im.InstallPath(op.Package), "/"))
		}
	case solver.OpUpdate:
		if op.Target.Type() != "metapackage" && im.downloaderType(op.Target) == "git" {
			return im.gitDL().prepare("update", op.Target, op.Initial, strings.TrimRight(im.InstallPath(op.Target), "/"))
		}
	case solver.OpUninstall:
		if op.Package.Type() != "metapackage" && im.downloaderType(op.Package) == "git" {
			return im.gitDL().prepare("uninstall", op.Package, nil, strings.TrimRight(im.InstallPath(op.Package), "/"))
		}
	}
	return nil
}

// cleanupOp is DownloadManager::cleanup (re-applies stashed changes).
func (im *InstallationManager) cleanupOp(op *solver.Operation) error {
	if op.Type == solver.OpUpdate && op.Target.Type() != "metapackage" && im.downloaderType(op.Target) == "git" {
		return im.gitDL().cleanup("update", op.Target, op.Initial, strings.TrimRight(im.InstallPath(op.Target), "/"))
	}
	return nil
}

// LocalChanges is DownloadManager/GitDownloader::getLocalChanges for an
// installed package ("" when the checkout is clean or not a git checkout).
func (im *InstallationManager) LocalChanges(p *pkg.Package) (string, error) {
	if im.downloaderType(p) != "git" {
		return "", nil
	}
	return im.gitDL().LocalChanges(p, strings.TrimRight(im.InstallPath(p), "/"))
}

// UnpushedChanges is GitDownloader::getUnpushedChanges for an installed package.
func (im *InstallationManager) UnpushedChanges(p *pkg.Package) (string, error) {
	if im.downloaderType(p) != "git" {
		return "", nil
	}
	return im.gitDL().UnpushedChanges(p, strings.TrimRight(im.InstallPath(p), "/"))
}
