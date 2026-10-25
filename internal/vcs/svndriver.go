package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/SvnDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
)

// SvnDriver is Repository\Vcs\SvnDriver.
type SvnDriver struct {
	base
	baseURL        string
	tags           []Ref
	tagsLoaded     bool
	branches       []Ref
	branchesLoaded bool
	rootIdentifier string
	// trunkPath false (repo config "trunk-path": false) is trunkFalse
	trunkPath        string
	trunkFalse       bool
	branchesPath     string
	branchesFalse    bool
	tagsPath         string
	tagsFalse        bool
	packagePath      string
	cacheCredentials bool
	util             *Svn
}

// NewSvnDriver is new SvnDriver($repoConfig, ...).
func NewSvnDriver(repoConfig *phpjson.Map, env *Env) *SvnDriver {
	d := &SvnDriver{base: newBase(repoConfig, env), trunkPath: "trunk", branchesPath: "branches", tagsPath: "tags", cacheCredentials: true}
	d.reader = d
	return d
}

var (
	reSvnLsLine   = regexp.MustCompile(`^\s*(\S+).*?(\S+)\s*$`)
	reSvnRevSufx  = regexp.MustCompile(`@\d+$`)
	reSvnIdentRev = regexp.MustCompile(`^(.+?)(@\d+)?/$`)
	reSvnChanged  = regexp.MustCompile(`^Last Changed Date: ([^(]+)`)
)

// svnPathOption reads a *-path repo option: isset() skips null, false
// disables the path, other scalars are used as strings.
func svnPathOption(m *phpjson.Map, key string, s *string, off *bool) {
	v, ok := m.Get(key)
	if !ok || v == nil {
		return
	}
	if b, isBool := v.(bool); isBool && !b {
		*off = true
		*s = ""
		return
	}
	str, _ := phpjson.ScalarString(v)
	*s = str
	*off = false
}

// Initialize is SvnDriver::initialize.
func (d *SvnDriver) Initialize() error {
	d.url = strings.TrimRight(svnNormalizeURL(d.url), "/")
	d.baseURL = d.url
	SvnCleanEnv()
	svnPathOption(d.repoConfig, "trunk-path", &d.trunkPath, &d.trunkFalse)
	svnPathOption(d.repoConfig, "branches-path", &d.branchesPath, &d.branchesFalse)
	svnPathOption(d.repoConfig, "tags-path", &d.tagsPath, &d.tagsFalse)
	if v, ok := d.repoConfig.Get("svn-cache-credentials"); ok {
		d.cacheCredentials = phpjson.Truthy(v)
	}
	if v, ok := d.repoConfig.Get("package-path"); ok && v != nil {
		s, _ := phpjson.ScalarString(v)
		d.packagePath = "/" + strings.Trim(s, "/")
	}
	// strrpos($this->url, '/' . $this->trunkPath) ('/' . false is '/')
	if pos := strings.LastIndex(d.url, "/"+d.trunkPath); pos >= 0 {
		d.baseURL = d.url[:pos]
	}
	d.cache = cache.New(d.Config.GetString("cache-repo-dir")+"/"+CacheKey(d.baseURL), "a-z0-9._", d.Config.GetBool("cache-read-only"))
	if _, err := d.Branches(); err != nil {
		return err
	}
	if _, err := d.Tags(); err != nil {
		return err
	}
	return nil
}

// RootIdentifier is SvnDriver::getRootIdentifier.
func (d *SvnDriver) RootIdentifier() (string, error) {
	if d.rootIdentifier != "" {
		return d.rootIdentifier, nil
	}
	return d.trunkPath, nil
}

// URL is SvnDriver::getUrl.
func (d *SvnDriver) URL() string { return d.url }

// Source is SvnDriver::getSource.
func (d *SvnDriver) Source(identifier string) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", "svn")
	m.Set("url", d.baseURL)
	m.Set("reference", identifier)
	return m
}

// Dist is SvnDriver::getDist.
func (d *SvnDriver) Dist(string) (*phpjson.Map, error) { return nil, nil }

func (d *SvnDriver) shouldCache(identifier string) bool {
	return d.cache != nil && reSvnRevSufx.MatchString(identifier)
}

// ComposerInformation is SvnDriver::getComposerInformation: revisions are
// cached as <identifier>.json and a missing composer.json is remembered.
func (d *SvnDriver) ComposerInformation(identifier string) (*phpjson.Map, error) {
	if e, ok := d.infoCache[identifier]; ok {
		return e.data, nil
	}
	if d.shouldCache(identifier) {
		if res, ok := d.cache.Read(identifier + ".json"); ok && len(res) > 0 && string(res) != "0" {
			if string(res) == `""` {
				res = []byte("null")
				d.cache.Write(identifier+".json", res)
			}
			v, err := jsonfile.Parse(res, "")
			if err != nil {
				return nil, err
			}
			m := phpjson.ToMap(v)
			if !phpjson.IsArray(v) {
				m = nil
			}
			d.infoCache[identifier] = infoEntry{m}
			return m, nil
		}
	}
	composer, err := d.baseComposerInformation(identifier)
	if err != nil {
		var te *httpdl.TransportError
		if !errors.As(err, &te) {
			return nil, err
		}
		msg := strings.ToLower(te.Msg)
		if !strings.Contains(msg, "path not found") && !strings.Contains(msg, "svn: warning: w160013") {
			return nil, err
		}
		composer = nil
	}
	if d.shouldCache(identifier) {
		var v phpjson.Value
		if composer != nil {
			v = composer
		}
		if enc, err := phpjson.Encode(v, phpjson.UnescapedUnicode|phpjson.UnescapedSlashes); err == nil {
			d.cache.Write(identifier+".json", []byte(enc))
		}
	}
	d.infoCache[identifier] = infoEntry{composer}
	return composer, nil
}

// HasComposerFile is VcsDriver::hasComposerFile over this driver's
// getComposerInformation.
func (d *SvnDriver) HasComposerFile(identifier string) (bool, error) {
	m, err := d.ComposerInformation(identifier)
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) {
			return false, nil
		}
		return false, err
	}
	return m != nil, nil
}

func svnSplitIdentifier(identifier string) (path, rev string) {
	identifier = "/" + strings.Trim(identifier, "/") + "/"
	if m := reSvnIdentRev.FindStringSubmatch(identifier); m != nil && m[2] != "" {
		return m[1], m[2]
	}
	return identifier, ""
}

// FileContent is SvnDriver::getFileContent: failures become transport
// errors (TransportException).
func (d *SvnDriver) FileContent(file, identifier string) (string, bool, error) {
	path, rev := svnSplitIdentifier(identifier)
	out, err := d.execute([]string{"svn", "cat"}, d.baseURL+path+file+rev)
	if err != nil {
		return "", false, &httpdl.TransportError{Msg: err.Error()}
	}
	if strings.TrimSpace(out) == "" {
		return "", false, nil
	}
	return out, true, nil
}

// ChangeDate is SvnDriver::getChangeDate.
func (d *SvnDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	path, rev := svnSplitIdentifier(identifier)
	out, err := d.execute([]string{"svn", "info"}, d.baseURL+path+rev)
	if err != nil {
		return time.Time{}, false, err
	}
	for _, line := range SplitLines(out) {
		if m := reSvnChanged.FindStringSubmatch(line); line != "" && m != nil {
			s := strings.TrimSpace(m[1])
			for _, layout := range []string{"2006-01-02 15:04:05 -0700", "2006-01-02 15:04:05"} {
				if t, err := time.Parse(layout, s); err == nil {
					return t, true, nil
				}
			}
			return time.Time{}, false, &RuntimeError{Msg: "Failed to parse time string (" + s + ")"}
		}
	}
	return time.Time{}, false, nil
}

// Tags is SvnDriver::getTags.
func (d *SvnDriver) Tags() ([]Ref, error) {
	if d.tagsLoaded {
		return d.tags, nil
	}
	var tags []Ref
	if !d.tagsFalse {
		out, err := d.execute([]string{"svn", "ls", "--verbose"}, d.baseURL+"/"+d.tagsPath)
		if err != nil {
			return nil, err
		}
		if out != "" {
			lastRev := 0
			for _, line := range SplitLines(out) {
				line = strings.TrimSpace(line)
				m := reSvnLsLine.FindStringSubmatch(line)
				if line == "" || m == nil {
					continue
				}
				rev := phpIntval(m[1])
				if m[2] == "./" {
					lastRev = rev
					continue
				}
				if lastRev > rev {
					rev = lastRev
				}
				tags = setRef(tags, strings.TrimRight(m[2], "/"), d.buildIdentifier("/"+d.tagsPath+"/"+m[2], rev))
			}
		}
	}
	d.tags, d.tagsLoaded = tags, true
	return d.tags, nil
}

// Branches is SvnDriver::getBranches.
func (d *SvnDriver) Branches() ([]Ref, error) {
	if d.branchesLoaded {
		return d.branches, nil
	}
	var branches []Ref
	trunkParent := d.baseURL + "/" + d.trunkPath
	if d.trunkFalse {
		trunkParent = d.baseURL + "/"
	}
	out, err := d.execute([]string{"svn", "ls", "--verbose"}, trunkParent)
	if err != nil {
		return nil, err
	}
	if out != "" {
		for _, line := range SplitLines(out) {
			line = strings.TrimSpace(line)
			m := reSvnLsLine.FindStringSubmatch(line)
			if line == "" || m == nil {
				continue
			}
			if m[2] == "./" {
				id := d.buildIdentifier("/"+d.trunkPath, phpIntval(m[1]))
				branches = setRef(branches, "trunk", id)
				d.rootIdentifier = id
				break
			}
		}
	}
	if !d.branchesFalse {
		out, err := d.execute([]string{"svn", "ls", "--verbose"}, d.baseURL+"/"+d.branchesPath)
		if err != nil {
			return nil, err
		}
		if out != "" {
			lastRev := 0
			for _, line := range SplitLines(strings.TrimSpace(out)) {
				line = strings.TrimSpace(line)
				m := reSvnLsLine.FindStringSubmatch(line)
				if line == "" || m == nil {
					continue
				}
				rev := phpIntval(m[1])
				if m[2] == "./" {
					lastRev = rev
					continue
				}
				if lastRev > rev {
					rev = lastRev
				}
				branches = setRef(branches, strings.TrimRight(m[2], "/"), d.buildIdentifier("/"+d.branchesPath+"/"+m[2], rev))
			}
		}
	}
	d.branches, d.branchesLoaded = branches, true
	return d.branches, nil
}

// phpIntval is (int) $string for the leading digits.
func phpIntval(s string) int {
	s = strings.TrimSpace(s)
	end := 0
	if end < len(s) && (s[end] == '-' || s[end] == '+') {
		end++
	}
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	n, _ := strconv.Atoi(s[:end])
	return n
}

// execute is SvnDriver::execute.
func (d *SvnDriver) execute(command []string, u string) (string, error) {
	if d.util == nil {
		d.util = NewSvn(d.baseURL, d.Env)
		d.util.SetCacheCredentials(d.cacheCredentials)
	}
	out, err := d.util.Execute(command, u, "", nil, false)
	if err == nil {
		return out, nil
	}
	if d.util.BinaryVersion() == "" {
		return "", &RuntimeError{Msg: "Failed to load " + Sanitize(d.url) + ", svn was not found, check that it is installed and in your PATH env.\n\n" + d.Process.ErrorOutput()}
	}
	return "", &RuntimeError{Msg: "Repository " + Sanitize(d.url) + " could not be processed, " + Sanitize(err.Error())}
}

// buildIdentifier is SvnDriver::buildIdentifier.
func (d *SvnDriver) buildIdentifier(baseDir string, revision int) string {
	return strings.TrimRight(baseDir, "/") + d.packagePath + "/@" + strconv.Itoa(revision)
}
