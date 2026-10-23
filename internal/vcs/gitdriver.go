package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/GitDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/phpjson"
)

// GitDriver is Repository\Vcs\GitDriver.
type GitDriver struct {
	base
	tags           []Ref
	tagsLoaded     bool
	branches       []Ref
	branchesLoaded bool
	rootIdentifier string
	repoDir        string
}

// NewGitDriver is new GitDriver($repoConfig, ...).
func NewGitDriver(repoConfig *phpjson.Map, env *Env) *GitDriver {
	d := &GitDriver{base: newBase(repoConfig, env)}
	d.reader = d
	return d
}

var (
	reDotGitSuffix = regexp.MustCompile(`[\\/]\.git/?$`)
	reSSHPort      = regexp.MustCompile(`^ssh://[^@]+@[^:]+:[^0-9]+`)
	reShowRefTag   = regexp.MustCompile(`^([a-f0-9]{40}) refs/tags/(\S+?)(\^\{\})?$`)
	reRemoteHEAD   = regexp.MustCompile(`^ *[^/]+/HEAD `)
	reBranchLine   = regexp.MustCompile(`^(?:\* )? *(\S+) *([a-f0-9]+)(?: .*)?$`)
	reCurrentBr    = regexp.MustCompile(`^\* +(\S+)`)
	reGitSupports  = regexp.MustCompile(`(?i)(^git://|\.git/?$|git(?:olite)?@|//git\.|//github.com/)`)
)

// Initialize is GitDriver::initialize.
func (d *GitDriver) Initialize() error {
	var cacheURL string
	if IsLocalPath(d.url) {
		d.url = reDotGitSuffix.ReplaceAllString(d.url, "")
		if st, err := os.Stat(d.url); err != nil || !st.IsDir() {
			return &RuntimeError{Msg: "Failed to read package information from " + Sanitize(d.url) + " as the path does not exist"}
		}
		d.repoDir = d.url
		cacheURL = d.url
		if r, err := filepath.EvalSymlinks(d.url); err == nil {
			if a, err := filepath.Abs(r); err == nil {
				cacheURL = a
			}
		}
	} else {
		vcsDir := d.Config.GetString("cache-vcs-dir")
		if !CacheUsable(vcsDir) {
			return &RuntimeError{Msg: "GitDriver requires a usable cache directory, and it looks like you set it to be disabled"}
		}
		d.repoDir = vcsDir + "/" + CacheKey(d.url) + "/"
		CleanEnv(d.Process)
		if err := os.MkdirAll(filepath.Dir(strings.TrimRight(d.repoDir, "/")), 0o777); err != nil {
			return &RuntimeError{Msg: "Can not clone " + Sanitize(d.url) + " to access package information. The \"" + filepath.Dir(strings.TrimRight(d.repoDir, "/")) + "\" directory is not writable by the current user."}
		}
		if reSSHPort.MatchString(d.url) {
			return &InvalidArgumentError{Msg: "The source URL " + Sanitize(d.url) + " is invalid, ssh URLs should have a port number after \":\".\nUse ssh://git@example.com:22/path or just git@example.com:path if you do not want to provide a password or custom port."}
		}
		ok, err := d.Git().SyncMirror(d.url, d.repoDir)
		if err != nil {
			return err
		}
		if !ok {
			if st, err := os.Stat(d.repoDir); err != nil || !st.IsDir() {
				return &RuntimeError{Msg: "Failed to clone " + Sanitize(d.url) + " to read package information from it"}
			}
			d.IO.WriteError("<error>Failed to update " + Sanitize(d.url) + ", package information from this repository may be outdated</error>")
		}
		cacheURL = d.url
	}
	if _, err := d.Tags(); err != nil {
		return err
	}
	if _, err := d.Branches(); err != nil {
		return err
	}
	d.cache = cache.New(d.Config.GetString("cache-repo-dir")+"/"+CacheKey(cacheURL), "a-z0-9._", d.Config.GetBool("cache-read-only"))
	return nil
}

// RootIdentifier is GitDriver::getRootIdentifier.
func (d *GitDriver) RootIdentifier() (string, error) {
	if d.rootIdentifier == "" {
		d.rootIdentifier = "master"
		if !IsLocalPath(d.url) {
			if def := d.Git().MirrorDefaultBranch(d.url, d.repoDir, false); def != "" {
				d.rootIdentifier = def
				return def, nil
			}
		}
		_, out, err := d.Process.Execute([]string{"git", "branch", "--no-color"}, d.repoDir)
		if err != nil {
			return "", err
		}
		lines := SplitLines(out)
		hasMaster := false
		for _, l := range lines {
			if l == "* master" {
				hasMaster = true
			}
		}
		if !hasMaster {
			for _, l := range lines {
				if m := reCurrentBr.FindStringSubmatch(l); l != "" && m != nil {
					d.rootIdentifier = m[1]
					break
				}
			}
		}
	}
	return d.rootIdentifier, nil
}

// URL is GitDriver::getUrl.
func (d *GitDriver) URL() string { return d.url }

// Source is GitDriver::getSource.
func (d *GitDriver) Source(identifier string) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", "git")
	m.Set("url", d.URL())
	m.Set("reference", identifier)
	return m
}

// Dist is GitDriver::getDist (git repositories have no dist).
func (d *GitDriver) Dist(identifier string) (*phpjson.Map, error) { return nil, nil }

// FileContent is GitDriver::getFileContent.
func (d *GitDriver) FileContent(file, identifier string) (string, bool, error) {
	if strings.HasPrefix(identifier, "-") {
		return "", false, &RuntimeError{Msg: "Invalid git identifier detected. Identifier must not start with a -, given: " + identifier}
	}
	_, content, err := d.Process.Execute([]string{"git", "show", identifier + ":" + file}, d.repoDir)
	if err != nil {
		return "", false, err
	}
	if strings.TrimSpace(content) == "" {
		return "", false, nil
	}
	return content, true, nil
}

// ChangeDate is GitDriver::getChangeDate.
func (d *GitDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	if strings.HasPrefix(identifier, "-") {
		return time.Time{}, false, &RuntimeError{Msg: "Invalid git identifier detected. Identifier must not start with a -, given: " + identifier}
	}
	_, out, err := d.Process.Execute(RevListCommand(d.Process, "-n1", "--format=%at", identifier), d.repoDir)
	if err != nil {
		return time.Time{}, false, err
	}
	ts, perr := strconv.ParseInt(strings.TrimSpace(ParseRevListOutput(d.Process, out)), 10, 64)
	if perr != nil {
		return time.Time{}, false, &RuntimeError{Msg: "Failed to parse the commit date of " + identifier}
	}
	return time.Unix(ts, 0).UTC(), true, nil
}

// Tags is GitDriver::getTags.
func (d *GitDriver) Tags() ([]Ref, error) {
	if !d.tagsLoaded {
		d.tagsLoaded = true
		_, out, err := d.Process.Execute([]string{"git", "show-ref", "--tags", "--dereference"}, d.repoDir)
		if err != nil {
			return nil, err
		}
		for _, line := range SplitLines(out) {
			if m := reShowRefTag.FindStringSubmatch(line); line != "" && m != nil {
				d.tags = setRef(d.tags, m[2], m[1])
			}
		}
	}
	return d.tags, nil
}

// setRef assigns $arr[$name] = $id keeping PHP array insertion order.
func setRef(refs []Ref, name, id string) []Ref {
	for i := range refs {
		if refs[i].Name == name {
			refs[i].ID = id
			return refs
		}
	}
	return append(refs, Ref{name, id})
}

// Branches is GitDriver::getBranches.
func (d *GitDriver) Branches() ([]Ref, error) {
	if !d.branchesLoaded {
		d.branchesLoaded = true
		_, out, err := d.Process.Execute([]string{"git", "branch", "--no-color", "--no-abbrev", "-v"}, d.repoDir)
		if err != nil {
			return nil, err
		}
		var branches []Ref
		for _, line := range SplitLines(out) {
			if line != "" && !reRemoteHEAD.MatchString(line) {
				if m := reBranchLine.FindStringSubmatch(line); m != nil && m[1][0] != '-' {
					branches = setRef(branches, m[1], m[2])
				}
			}
		}
		d.branches = branches
	}
	return d.branches, nil
}

// GitSupports is GitDriver::supports.
func GitSupports(env *Env, u string, deep bool) (bool, error) {
	if reGitSupports.MatchString(u) {
		return true, nil
	}
	if IsLocalPath(u) {
		u = PlatformPath(u)
		if st, err := os.Stat(u); err != nil || !st.IsDir() {
			return false, nil
		}
		p := NewProcess(env.IO, env.Process.Timeout)
		if code, _, _ := p.Execute([]string{"git", "tag"}, u); code == 0 {
			return true, nil
		}
		if err := CheckForRepoOwnershipError(p.ErrorOutput(), u); err != nil {
			return false, err
		}
	}
	if !deep {
		return false, nil
	}
	p := NewProcess(env.IO, env.Process.Timeout)
	g := &Git{IO: env.IO, Config: env.Config, Process: p, HTTP: env.HTTP, Auth: env.Auth}
	CleanEnv(p)
	if _, err := g.RunCommands([][]string{{"git", "ls-remote", "--heads", "--", "%url%"}}, u, os.TempDir(), false); err != nil {
		return false, nil
	}
	return true, nil
}
