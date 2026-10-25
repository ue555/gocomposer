package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/FossilDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
)

// FossilDriver is Repository\Vcs\FossilDriver.
type FossilDriver struct {
	base
	tags           []Ref
	tagsLoaded     bool
	branches       []Ref
	branchesLoaded bool
	repoFile       string
	checkoutDir    string
}

// NewFossilDriver is new FossilDriver($repoConfig, ...).
func NewFossilDriver(repoConfig *phpjson.Map, env *Env) *FossilDriver {
	d := &FossilDriver{base: newBase(repoConfig, env)}
	d.reader = d
	return d
}

// Initialize is FossilDriver::initialize.
func (d *FossilDriver) Initialize() error {
	// checkFossil
	if code, _, err := d.Process.Execute([]string{"fossil", "version"}, ""); err != nil {
		return err
	} else if code != 0 {
		return &RuntimeError{Msg: "fossil was not found, check that it is installed and in your PATH env.\n\n" + d.Process.ErrorOutput()}
	}
	if err := ProhibitURL(d.Env, d.url); err != nil {
		return err
	}
	if IsLocalPath(d.url) && isDir(d.url) {
		d.checkoutDir = d.url
	} else {
		if !CacheUsable(d.Config.GetString("cache-repo-dir")) || !CacheUsable(d.Config.GetString("cache-vcs-dir")) {
			return &RuntimeError{Msg: "FossilDriver requires a usable cache directory, and it looks like you set it to be disabled"}
		}
		localName := reNonAlnum.ReplaceAllString(d.url, "-")
		d.repoFile = d.Config.GetString("cache-repo-dir") + "/" + localName + ".fossil"
		d.checkoutDir = d.Config.GetString("cache-vcs-dir") + "/" + localName + "/"
		if err := d.updateLocalRepo(); err != nil {
			return err
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

func isFile(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}

// updateLocalRepo is FossilDriver::updateLocalRepo.
func (d *FossilDriver) updateLocalRepo() error {
	if err := ensureDirectoryExists(d.checkoutDir); err != nil {
		return err
	}
	if !dirWritable(filepath.Dir(strings.TrimRight(d.checkoutDir, "/"))) {
		return &RuntimeError{Msg: "Can not clone " + Sanitize(d.url) + ` to access package information. The "` + d.checkoutDir + `" directory is not writable by the current user.`}
	}
	p := d.Process
	if isFile(d.repoFile) && isDir(d.checkoutDir) && fossilOK(p, []string{"fossil", "info"}, d.checkoutDir) {
		if !fossilOK(p, []string{"fossil", "pull"}, d.checkoutDir) {
			d.IO.WriteError("<error>Failed to update " + Sanitize(d.url) + ", package information from this repository may be outdated (" + p.ErrorOutput() + ")</error>")
		}
		return nil
	}
	_ = fsutil.RemoveAll(d.checkoutDir)
	_ = os.Remove(d.repoFile)
	if err := ensureDirectoryExists(d.checkoutDir); err != nil {
		return err
	}
	if code, _, err := p.Execute([]string{"fossil", "clone", "--", d.url, d.repoFile}, ""); err != nil {
		return err
	} else if code != 0 {
		return &RuntimeError{Msg: "Failed to clone " + Sanitize(d.url) + " to repository " + d.repoFile + "\n\n" + p.ErrorOutput()}
	}
	if code, _, err := p.Execute([]string{"fossil", "open", "--nested", "--", d.repoFile}, d.checkoutDir); err != nil {
		return err
	} else if code != 0 {
		return &RuntimeError{Msg: "Failed to open repository " + d.repoFile + " in " + d.checkoutDir + "\n\n" + p.ErrorOutput()}
	}
	return nil
}

func fossilOK(p *Process, cmd []string, cwd string) bool {
	code, _, err := p.Execute(cmd, cwd)
	return err == nil && code == 0
}

// RootIdentifier is FossilDriver::getRootIdentifier.
func (d *FossilDriver) RootIdentifier() (string, error) { return "trunk", nil }

// URL is FossilDriver::getUrl.
func (d *FossilDriver) URL() string { return d.url }

// Source is FossilDriver::getSource.
func (d *FossilDriver) Source(identifier string) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", "fossil")
	m.Set("url", d.URL())
	m.Set("reference", identifier)
	return m
}

// Dist is FossilDriver::getDist.
func (d *FossilDriver) Dist(string) (*phpjson.Map, error) { return nil, nil }

// FileContent is FossilDriver::getFileContent.
func (d *FossilDriver) FileContent(file, identifier string) (string, bool, error) {
	if strings.HasPrefix(identifier, "-") {
		return "", false, &RuntimeError{Msg: "Invalid fossil identifier detected. Identifier must not start with a -, given: " + identifier}
	}
	_, content, err := d.Process.Execute([]string{"fossil", "cat", "-r", identifier, "--", file}, d.checkoutDir)
	if err != nil {
		return "", false, err
	}
	if strings.TrimSpace(content) == "" {
		return "", false, nil
	}
	return content, true, nil
}

// ChangeDate is FossilDriver::getChangeDate (the date of the last change of
// composer.json, whatever the identifier).
func (d *FossilDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	_, out, err := d.Process.Execute([]string{"fossil", "finfo", "-b", "-n", "1", "composer.json"}, d.checkoutDir)
	if err != nil {
		return time.Time{}, false, err
	}
	parts := strings.SplitN(strings.TrimSpace(out), " ", 3)
	if len(parts) < 2 {
		return time.Time{}, false, &RuntimeError{Msg: "Undefined array key 1"}
	}
	for _, layout := range []string{"2006-01-02", "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, parts[1]); err == nil {
			return t, true, nil
		}
	}
	return time.Time{}, false, &RuntimeError{Msg: "Failed to parse time string (" + parts[1] + ")"}
}

var reFossilCurrent = regexp.MustCompile(`^\*`)

// Tags is FossilDriver::getTags.
func (d *FossilDriver) Tags() ([]Ref, error) {
	if !d.tagsLoaded {
		_, out, err := d.Process.Execute([]string{"fossil", "tag", "list"}, d.checkoutDir)
		if err != nil {
			return nil, err
		}
		var tags []Ref
		for _, tag := range SplitLines(out) {
			tags = setRef(tags, tag, tag)
		}
		d.tags, d.tagsLoaded = tags, true
	}
	return d.tags, nil
}

// Branches is FossilDriver::getBranches.
func (d *FossilDriver) Branches() ([]Ref, error) {
	if !d.branchesLoaded {
		_, out, err := d.Process.Execute([]string{"fossil", "branch", "list"}, d.checkoutDir)
		if err != nil {
			return nil, err
		}
		var branches []Ref
		for _, b := range SplitLines(out) {
			b = strings.TrimSpace(reFossilCurrent.ReplaceAllString(strings.TrimSpace(b), ""))
			branches = setRef(branches, b, b)
		}
		d.branches, d.branchesLoaded = branches, true
	}
	return d.branches, nil
}
