package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/Hg.php, src/Composer/Repository/Vcs/HgDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
)

var (
	hgVersionOnce sync.Once
	hgVersion     string
	reHgVersion   = regexp.MustCompile(`^.+? (\d+(?:\.\d+)+)(?:\+.*?)?\)?\r?\n`)
	reHgAuthURL   = regexp.MustCompile(`(?mi)^(ssh|https?)://(?:([^:@]+)(?::([^:@]+))?@)?([^/]+)(/.*)?`)
)

// HgVersion is Hg::getVersion ("" when hg is not available; looked up once
// per process like the static property).
func HgVersion(p *Process) string {
	hgVersionOnce.Do(func() {
		if code, out, _ := p.Execute([]string{"hg", "--version"}, ""); code == 0 {
			if m := reHgVersion.FindStringSubmatch(out); m != nil {
				hgVersion = m[1]
			}
		}
	})
	return hgVersion
}

// Hg is Composer\Util\Hg.
type Hg struct{ env *Env }

// NewHg is new Hg($io, $config, $process).
func NewHg(env *Env) *Hg { return &Hg{env: env} }

// RunCommand is Hg::runCommand: the command is tried as is, then with the
// credentials known for the host.
func (h *Hg) RunCommand(command func(u string) []string, u, cwd string) error {
	if err := ProhibitURL(h.env, u); err != nil {
		return err
	}
	p := h.env.Process
	code, _, err := p.Execute(command(u), cwd)
	if err != nil {
		return err
	}
	if code == 0 {
		return nil
	}
	var errOut string
	if m := reHgAuthURL.FindStringSubmatch(u); m != nil && h.env.Auth != nil && h.env.Auth.Has(m[4]) {
		var authURL string
		if m[1] == "ssh" {
			user := ""
			if m[2] != "" {
				user = rawURLEncode(m[2]) + "@"
			}
			authURL = m[1] + "://" + user + m[4] + m[5]
		} else {
			a := h.env.Auth.Get(m[4])
			authURL = m[1] + "://" + rawURLEncode(a.Username) + ":" + rawURLEncode(a.Password) + "@" + m[4] + m[5]
		}
		code, _, err := p.Execute(command(authURL), cwd)
		if err != nil {
			return err
		}
		if code == 0 {
			return nil
		}
		errOut = p.ErrorOutput()
	} else {
		errOut = "The given URL (" + Sanitize(u) + ") does not match the required format (ssh|http(s)://(username:password@)example.com/path-to-repository)"
	}
	return h.throwException("Failed to clone "+u+", \n\n"+errOut, u)
}

func (h *Hg) throwException(message, u string) error {
	if HgVersion(h.env.Process) == "" {
		return &RuntimeError{Msg: Sanitize("Failed to clone " + u + ", hg was not found, check that it is installed and in your PATH env.\n\n" + h.env.Process.ErrorOutput())}
	}
	return &RuntimeError{Msg: Sanitize(message)}
}

// HgDriver is Repository\Vcs\HgDriver.
type HgDriver struct {
	base
	tags           []Ref
	tagsLoaded     bool
	branches       []Ref
	branchesLoaded bool
	rootIdentifier string
	rootLoaded     bool
	repoDir        string
}

// NewHgDriver is new HgDriver($repoConfig, ...).
func NewHgDriver(repoConfig *phpjson.Map, env *Env) *HgDriver {
	d := &HgDriver{base: newBase(repoConfig, env)}
	d.reader = d
	return d
}

var reNonAlnum = regexp.MustCompile(`(?i)[^a-z0-9]`)

// Initialize is HgDriver::initialize.
func (d *HgDriver) Initialize() error {
	if IsLocalPath(d.url) {
		d.repoDir = d.url
	} else {
		cacheDir := d.Config.GetString("cache-vcs-dir")
		if !CacheUsable(cacheDir) {
			return &RuntimeError{Msg: "HgDriver requires a usable cache directory, and it looks like you set it to be disabled"}
		}
		d.repoDir = cacheDir + "/" + reNonAlnum.ReplaceAllString(Sanitize(d.url), "-") + "/"
		if err := ensureDirectoryExists(cacheDir); err != nil {
			return err
		}
		if !dirWritable(filepath.Dir(strings.TrimRight(d.repoDir, "/"))) {
			return &RuntimeError{Msg: "Can not clone " + Sanitize(d.url) + ` to access package information. The "` + cacheDir + `" directory is not writable by the current user.`}
		}
		if err := ProhibitURL(d.Env, d.url); err != nil {
			return err
		}
		p := d.Process
		if isDir(d.repoDir) && hgSummaryOK(p, d.repoDir) {
			if code, _, _ := p.Execute([]string{"hg", "pull"}, d.repoDir); code != 0 {
				d.IO.WriteError("<error>Failed to update " + Sanitize(d.url) + ", package information from this repository may be outdated (" + p.ErrorOutput() + ")</error>")
			}
		} else {
			_ = fsutil.RemoveAll(d.repoDir)
			repoDir := d.repoDir
			if err := NewHg(d.Env).RunCommand(func(u string) []string {
				return []string{"hg", "clone", "--noupdate", "--", u, repoDir}
			}, d.url, ""); err != nil {
				return err
			}
		}
	}
	if _, err := d.Tags(); err != nil {
		return err
	}
	if _, err := d.Branches(); err != nil {
		return err
	}
	return nil
}

// RootIdentifier is HgDriver::getRootIdentifier.
func (d *HgDriver) RootIdentifier() (string, error) {
	if !d.rootLoaded {
		_, out, err := d.Process.Execute([]string{"hg", "tip", "--template", "{node}"}, d.repoDir)
		if err != nil {
			return "", err
		}
		lines := SplitLines(out)
		if len(lines) == 0 {
			return "", &RuntimeError{Msg: "Undefined array key 0"}
		}
		d.rootIdentifier, d.rootLoaded = lines[0], true
	}
	return d.rootIdentifier, nil
}

// URL is HgDriver::getUrl.
func (d *HgDriver) URL() string { return d.url }

// Source is HgDriver::getSource.
func (d *HgDriver) Source(identifier string) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", "hg")
	m.Set("url", d.URL())
	m.Set("reference", identifier)
	return m
}

// Dist is HgDriver::getDist.
func (d *HgDriver) Dist(string) (*phpjson.Map, error) { return nil, nil }

// FileContent is HgDriver::getFileContent.
func (d *HgDriver) FileContent(file, identifier string) (string, bool, error) {
	if strings.HasPrefix(identifier, "-") {
		return "", false, &RuntimeError{Msg: "Invalid hg identifier detected. Identifier must not start with a -, given: " + identifier}
	}
	_, content, err := d.Process.Execute([]string{"hg", "cat", "-r", identifier, "--", file}, d.repoDir)
	if err != nil {
		return "", false, err
	}
	if strings.TrimSpace(content) == "" || strings.TrimSpace(content) == "0" {
		return "", false, nil
	}
	return content, true, nil
}

// ChangeDate is HgDriver::getChangeDate.
func (d *HgDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	if strings.HasPrefix(identifier, "-") {
		return time.Time{}, false, &RuntimeError{Msg: "Invalid hg identifier detected. Identifier must not start with a -, given: " + identifier}
	}
	_, out, err := d.Process.Execute([]string{"hg", "log", "--template", "{date|rfc3339date}", "-r", identifier}, d.repoDir)
	if err != nil {
		return time.Time{}, false, err
	}
	s := strings.TrimSpace(out)
	if s == "" {
		// new \DateTimeImmutable('') is the current time
		return time.Now().UTC().Truncate(time.Second), true, nil
	}
	t, perr := time.Parse(time.RFC3339, s)
	if perr != nil {
		return time.Time{}, false, &RuntimeError{Msg: "Failed to parse time string (" + s + ")"}
	}
	return t, true, nil
}

var (
	reHgTag      = regexp.MustCompile(`^([^\s]+)\s+\d+:(.*)$`)
	reHgBranch   = regexp.MustCompile(`^([^\s]+)\s+\d+:([a-f0-9]+)`)
	reHgBookmark = regexp.MustCompile(`^(?:[\s*]*)([^\s]+)\s+\d+:(.*)$`)
)

// Tags is HgDriver::getTags.
func (d *HgDriver) Tags() ([]Ref, error) {
	if !d.tagsLoaded {
		_, out, err := d.Process.Execute([]string{"hg", "tags"}, d.repoDir)
		if err != nil {
			return nil, err
		}
		var tags []Ref
		for _, line := range SplitLines(out) {
			if m := reHgTag.FindStringSubmatch(line); line != "" && line != "0" && m != nil {
				tags = setRef(tags, m[1], m[2])
			}
		}
		var out2 []Ref
		for _, t := range tags {
			if t.Name != "tip" {
				out2 = append(out2, t)
			}
		}
		d.tags, d.tagsLoaded = out2, true
	}
	return d.tags, nil
}

// Branches is HgDriver::getBranches (bookmarks included, branches win).
func (d *HgDriver) Branches() ([]Ref, error) {
	if !d.branchesLoaded {
		b, err := HgBranches(d.Process, d.repoDir)
		if err != nil {
			return nil, err
		}
		d.branches, d.branchesLoaded = b, true
	}
	return d.branches, nil
}

// HgBranches is the body of HgDriver::getBranches run in cwd ("" for the
// current directory, as VersionGuesser uses a driver it never initialized).
func HgBranches(p *Process, cwd string) ([]Ref, error) {
	_, out, err := p.Execute([]string{"hg", "branches"}, cwd)
	if err != nil {
		return nil, err
	}
	var branches, bookmarks []Ref
	for _, line := range SplitLines(out) {
		if m := reHgBranch.FindStringSubmatch(line); line != "" && line != "0" && m != nil && m[1][0] != '-' {
			branches = setRef(branches, m[1], m[2])
		}
	}
	_, out, err = p.Execute([]string{"hg", "bookmarks"}, cwd)
	if err != nil {
		return nil, err
	}
	for _, line := range SplitLines(out) {
		if m := reHgBookmark.FindStringSubmatch(line); line != "" && line != "0" && m != nil && m[1][0] != '-' {
			bookmarks = setRef(bookmarks, m[1], m[2])
		}
	}
	merged := append([]Ref{}, bookmarks...)
	for _, b := range branches {
		merged = setRef(merged, b.Name, b.ID)
	}
	return merged, nil
}

// ensureDirectoryExists is Filesystem::ensureDirectoryExists.
func ensureDirectoryExists(dir string) error {
	if isDir(dir) {
		return nil
	}
	if _, err := os.Lstat(dir); err == nil {
		return &RuntimeError{Msg: dir + " exists and is not a directory."}
	}
	if err := os.MkdirAll(dir, 0o777); err != nil && !isDir(dir) {
		return &RuntimeError{Msg: dir + " does not exist and could not be created: " + err.Error()}
	}
	return nil
}

func hgSummaryOK(p *Process, dir string) bool {
	code, _, _ := p.Execute([]string{"hg", "summary"}, dir)
	return code == 0
}
