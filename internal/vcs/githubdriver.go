package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/GitHubDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
)

// GitHubDriver is Repository\Vcs\GitHubDriver.
type GitHubDriver struct {
	base
	owner, repository  string
	tags, branches     []Ref
	tagsOK, branchesOK bool
	rootIdentifier     string
	repoData           *phpjson.Map
	hasIssues          bool
	isPrivate          bool
	isArchived         bool
	fundingInfo        []phpjson.Value
	fundingLoaded      bool
	gitDriver          *GitDriver
	allowGitFallback   bool
}

// NewGitHubDriver is new GitHubDriver($repoConfig, ...).
func NewGitHubDriver(repoConfig *phpjson.Map, env *Env) Driver {
	d := &GitHubDriver{base: newBase(repoConfig, env), allowGitFallback: true}
	d.reader = d
	return d
}

// Initialize is GitHubDriver::initialize.
func (d *GitHubDriver) Initialize() error {
	m := reGitHubURL.FindStringSubmatch(d.url)
	if m == nil {
		return &InvalidArgumentError{Msg: fmt.Sprintf("The GitHub repository URL %s is invalid.", Sanitize(d.url))}
	}
	d.owner, d.repository = m[4], m[5]
	origin := m[2]
	if origin == "" {
		origin = m[3]
	}
	d.originURL = strings.ToLower(origin)
	if d.originURL == "www.github.com" {
		d.originURL = "github.com"
	}
	d.cache = cache.New(d.Config.GetString("cache-repo-dir")+"/"+d.originURL+"/"+d.owner+"/"+d.repository, "a-z0-9._", d.Config.GetBool("cache-read-only"))
	if v, ok := d.repoConfig.Get("allow-git-fallback"); ok && v == false {
		d.allowGitFallback = false
	}
	noAPI := false
	if v, ok := d.repoConfig.Get("no-api"); ok && phpjson.Truthy(v) {
		noAPI = true
	}
	if d.Config.Get("use-github-api") == false || noAPI {
		return d.setupGitDriver(d.url)
	}
	return d.fetchRootIdentifier()
}

func (d *GitHubDriver) apiURL() string {
	if d.originURL == "github.com" {
		return "https://api.github.com"
	}
	return "https://" + d.originURL + "/api/v3"
}

// RootIdentifier is GitHubDriver::getRootIdentifier.
func (d *GitHubDriver) RootIdentifier() (string, error) {
	if d.gitDriver != nil {
		return d.gitDriver.RootIdentifier()
	}
	return d.rootIdentifier, nil
}

// URL is GitHubDriver::getUrl.
func (d *GitHubDriver) URL() string {
	if d.gitDriver != nil {
		return d.gitDriver.URL()
	}
	return "https://" + d.originURL + "/" + d.owner + "/" + d.repository + ".git"
}

func (d *GitHubDriver) sshURL() string {
	if strings.Contains(d.originURL, ":") {
		return "ssh://git@" + d.originURL + "/" + d.owner + "/" + d.repository + ".git"
	}
	return "git@" + d.originURL + ":" + d.owner + "/" + d.repository + ".git"
}

// Source is GitHubDriver::getSource.
func (d *GitHubDriver) Source(identifier string) *phpjson.Map {
	if d.gitDriver != nil {
		return d.gitDriver.Source(identifier)
	}
	u := d.URL()
	if d.isPrivate {
		u = d.sshURL()
	}
	m := phpjson.NewMap()
	m.Set("type", "git")
	m.Set("url", u)
	m.Set("reference", identifier)
	return m
}

// Dist is GitHubDriver::getDist.
func (d *GitHubDriver) Dist(identifier string) (*phpjson.Map, error) {
	m := phpjson.NewMap()
	m.Set("type", "zip")
	m.Set("url", d.apiURL()+"/repos/"+d.owner+"/"+d.repository+"/zipball/"+identifier)
	m.Set("reference", identifier)
	m.Set("shasum", "")
	return m, nil
}

// ComposerInformation is GitHubDriver::getComposerInformation.
func (d *GitHubDriver) ComposerInformation(identifier string) (*phpjson.Map, error) {
	if d.gitDriver != nil {
		return d.gitDriver.ComposerInformation(identifier)
	}
	if e, ok := d.infoCache[identifier]; ok {
		return e.data, nil
	}
	var composer *phpjson.Map
	loaded := false
	if d.shouldCache(identifier) {
		if res, ok := d.cache.Read(identifier); ok && len(res) > 0 {
			v, err := jsonfile.Parse(res, "")
			if err != nil {
				return nil, err
			}
			composer, loaded = phpjson.ToMap(v), true
		}
	}
	if !loaded {
		c, err := d.baseComposerInformation(identifier)
		if err != nil {
			return nil, err
		}
		composer = c
		if d.shouldCache(identifier) {
			var v phpjson.Value
			if composer != nil {
				v = composer
			}
			if enc, err := phpjson.Encode(v, phpjson.UnescapedUnicode|phpjson.UnescapedSlashes); err == nil {
				d.cache.Write(identifier, []byte(enc))
			}
		}
	}
	if composer != nil {
		composer = composer.Clone()
		support := phpjson.NewMap()
		if v, ok := composer.Get("support"); ok && phpjson.IsArray(v) {
			support = phpjson.ToMap(v).Clone()
		}
		if !support.Isset("source") {
			label, err := refLabelLazy(identifier, d.Tags, d.Branches)
			if err != nil {
				return nil, err
			}
			support.Set("source", fmt.Sprintf("https://%s/%s/%s/tree/%s", d.originURL, d.owner, d.repository, label))
		}
		if !support.Isset("issues") && d.hasIssues {
			support.Set("issues", fmt.Sprintf("https://%s/%s/%s/issues", d.originURL, d.owner, d.repository))
		}
		composer.Set("support", support)
		if !composer.Isset("abandoned") && d.isArchived {
			composer.Set("abandoned", true)
		}
		if !composer.Isset("funding") {
			if f := d.funding(); len(f) > 0 {
				composer.Set("funding", f)
			}
		}
	}
	d.infoCache[identifier] = infoEntry{composer}
	return composer, nil
}

var (
	reFundKV      = regexp.MustCompile(`^(\w+)\s*:\s*(.+)$`)
	reFundList    = regexp.MustCompile(`^\[(.*?)\](?:\s*#.*)?$`)
	reFundSplit   = regexp.MustCompile(`['"]?\s*,\s*['"]?`)
	reFundValue   = regexp.MustCompile(`^([^#].*?)(?:\s+#.*)?$`)
	reFundKeyOnly = regexp.MustCompile(`^(\w+)\s*:\s*#\s*$`)
	reFundDash    = regexp.MustCompile(`^-\s*(.+)(?:\s+#.*)?$`)
	reFundComma   = regexp.MustCompile(`^(.+),(?:\s*#.*)?$`)
	reFundDomain  = regexp.MustCompile(`^[a-z0-9-]+\.[a-z]{2,3}$`)
)

// funding is GitHubDriver::getFundingInfo.
func (d *GitHubDriver) funding() []phpjson.Value {
	if d.fundingLoaded {
		return d.fundingInfo
	}
	d.fundingLoaded = true
	if d.originURL != "github.com" {
		return nil
	}
	var content string
	for _, file := range []string{d.apiURL() + "/repos/" + d.owner + "/" + d.repository + "/contents/.github/FUNDING.yml", d.apiURL() + "/repos/" + d.owner + "/.github/contents/FUNDING.yml"} {
		resp, err := d.HTTP.GetNoAuthRetry(file) // 'retry-auth-failure' => false
		if err != nil {
			continue
		}
		v, err := decodeJSON(resp)
		m := phpjson.ToMap(v)
		if err != nil || m == nil || !phpjson.Truthy(mustGetMap(m, "content")) || m.Str("encoding") != "base64" {
			continue
		}
		b, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(m.Str("content"), "\n", ""))
		if err != nil || len(b) == 0 {
			continue
		}
		content = string(b)
		break
	}
	if content == "" {
		return nil
	}
	type item struct{ typ, url string }
	var result []item
	key := ""
	trimQ := func(s string) string { return strings.Trim(s, "\"' ") }
	for _, line := range regexp.MustCompile(`\r?\n`).Split(content, -1) {
		line = strings.TrimSpace(line)
		if m := reFundKV.FindStringSubmatch(line); m != nil {
			if m[2] == "[" {
				key = m[1]
				continue
			}
			if m2 := reFundList.FindStringSubmatch(m[2]); m2 != nil {
				for _, it := range reFundSplit.Split(m2[1], -1) {
					result = append(result, item{m[1], trimQ(strings.TrimSpace(it))})
				}
			} else if m2 := reFundValue.FindStringSubmatch(m[2]); m2 != nil {
				result = append(result, item{m[1], trimQ(m2[1])})
			}
			key = ""
		} else if m := reFundKeyOnly.FindStringSubmatch(line); m != nil {
			key = m[1]
		} else if key != "" {
			if m := reFundDash.FindStringSubmatch(line); m != nil {
				result = append(result, item{key, trimQ(m[1])})
			} else if m := reFundComma.FindStringSubmatch(line); m != nil {
				result = append(result, item{key, trimQ(m[1])})
			} else if line == "]" {
				key = ""
			}
		}
	}
	var out []phpjson.Value
	for _, it := range result {
		u := it.url
		switch it.typ {
		case "community_bridge":
			u = "https://funding.communitybridge.org/projects/" + path.Base(u)
		case "github":
			u = "https://github.com/" + path.Base(u)
		case "issuehunt":
			u = "https://issuehunt.io/r/" + u
		case "ko_fi":
			u = "https://ko-fi.com/" + path.Base(u)
		case "liberapay":
			u = "https://liberapay.com/" + path.Base(u)
		case "open_collective":
			u = "https://opencollective.com/" + path.Base(u)
		case "patreon":
			u = "https://www.patreon.com/" + path.Base(u)
		case "tidelift":
			u = "https://tidelift.com/funding/github/" + u
		case "polar":
			u = "https://polar.sh/" + path.Base(u)
		case "buy_me_a_coffee":
			u = "https://www.buymeacoffee.com/" + path.Base(u)
		case "thanks_dev":
			u = "https://thanks.dev/" + u
		case "otechie":
			u = "https://otechie.com/" + path.Base(u)
		case "custom":
			pu, err := url.Parse(u)
			if err != nil {
				continue
			}
			if pu.Scheme == "" && pu.Host == "" {
				if reFundDomain.MatchString(u) {
					u = "https://" + u
				} else {
					d.IO.WriteError("<warning>Funding URL " + u + " not in a supported format.</warning>")
					continue
				}
			}
		}
		e := phpjson.NewMap()
		e.Set("type", it.typ)
		e.Set("url", u)
		out = append(out, e)
	}
	d.fundingInfo = out
	return out
}

// FileContent is GitHubDriver::getFileContent.
func (d *GitHubDriver) FileContent(file, identifier string) (string, bool, error) {
	if d.gitDriver != nil {
		return d.gitDriver.FileContent(file, identifier)
	}
	resource := d.apiURL() + "/repos/" + d.owner + "/" + d.repository + "/contents/" + file + "?ref=" + url.QueryEscape(identifier)
	resp, err := d.getContentsGitHub(resource, false)
	if err != nil {
		return "", false, err
	}
	v, err := decodeJSON(resp)
	if err != nil {
		return "", false, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	if (!m.Isset("content") || m.Str("content") == "") && m.Str("encoding") == "none" && m.Isset("git_url") {
		resp, err := d.getContentsGitHub(m.Str("git_url"), false)
		if err != nil {
			return "", false, err
		}
		v, err := decodeJSON(resp)
		if err != nil {
			return "", false, err
		}
		if m = phpjson.ToMap(v); m == nil {
			m = phpjson.NewMap()
		}
	}
	if !m.Isset("content") || m.Str("encoding") != "base64" {
		return "", false, &RuntimeError{Msg: "Could not retrieve " + file + " for " + identifier}
	}
	b, err := base64.StdEncoding.DecodeString(strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == ' ' {
			return -1
		}
		return r
	}, m.Str("content")))
	if err != nil {
		return "", false, &RuntimeError{Msg: "Could not retrieve " + file + " for " + identifier}
	}
	return string(b), true, nil
}

// ChangeDate is GitHubDriver::getChangeDate.
func (d *GitHubDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	if d.gitDriver != nil {
		return d.gitDriver.ChangeDate(identifier)
	}
	resp, err := d.getContentsGitHub(d.apiURL()+"/repos/"+d.owner+"/"+d.repository+"/commits/"+url.QueryEscape(identifier), false)
	if err != nil {
		return time.Time{}, false, err
	}
	v, err := decodeJSON(resp)
	if err != nil {
		return time.Time{}, false, err
	}
	date := ""
	if m := phpjson.ToMap(v); m != nil {
		if c := m.Map("commit"); c != nil {
			if cm := c.Map("committer"); cm != nil {
				date = cm.Str("date")
			}
		}
	}
	t, err := parsePHPDateTime(date)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

// Tags is GitHubDriver::getTags.
func (d *GitHubDriver) Tags() ([]Ref, error) {
	if d.gitDriver != nil {
		return d.gitDriver.Tags()
	}
	if !d.tagsOK {
		var tags []Ref
		resource := d.apiURL() + "/repos/" + d.owner + "/" + d.repository + "/tags?per_page=100"
		for resource != "" {
			resp, err := d.getContentsGitHub(resource, false)
			if err != nil {
				return nil, err
			}
			v, err := decodeJSON(resp)
			if err != nil {
				return nil, err
			}
			phpjson.Each(v, func(_ string, tv phpjson.Value) {
				if m := phpjson.ToMap(tv); m != nil {
					name, _ := phpjson.ScalarString(mustGetMap(m, "name"))
					sha := ""
					if c := m.Map("commit"); c != nil {
						sha = c.Str("sha")
					}
					tags = setRef(tags, name, sha)
				}
			})
			resource = nextPage(resp)
		}
		d.tags, d.tagsOK = tags, true
	}
	return d.tags, nil
}

// Branches is GitHubDriver::getBranches.
func (d *GitHubDriver) Branches() ([]Ref, error) {
	if d.gitDriver != nil {
		return d.gitDriver.Branches()
	}
	if !d.branchesOK {
		var branches []Ref
		resource := d.apiURL() + "/repos/" + d.owner + "/" + d.repository + "/git/refs/heads?per_page=100"
		for resource != "" {
			resp, err := d.getContentsGitHub(resource, false)
			if err != nil {
				return nil, err
			}
			v, err := decodeJSON(resp)
			if err != nil {
				return nil, err
			}
			phpjson.Each(v, func(_ string, bv phpjson.Value) {
				if m := phpjson.ToMap(bv); m != nil {
					ref := m.Str("ref")
					name := ""
					if len(ref) > 11 {
						name = ref[11:]
					}
					if name != "gh-pages" {
						sha := ""
						if o := m.Map("object"); o != nil {
							sha = o.Str("sha")
						}
						branches = setRef(branches, name, sha)
					}
				}
			})
			resource = nextPage(resp)
		}
		d.branches, d.branchesOK = branches, true
	}
	return d.branches, nil
}

func (d *GitHubDriver) fetchRootIdentifier() error {
	if d.repoData != nil {
		return nil
	}
	resp, err := d.getContentsGitHub(d.apiURL()+"/repos/"+d.owner+"/"+d.repository, true)
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) && te.StatusCode == 499 {
			if err := d.attemptCloneFallback(err); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if resp != nil {
		v, err := decodeJSON(resp)
		if err != nil {
			return err
		}
		d.repoData = phpjson.ToMap(v)
	}
	if d.repoData == nil && d.gitDriver != nil {
		return nil
	}
	if d.repoData == nil {
		d.repoData = phpjson.NewMap()
	}
	if o := d.repoData.Map("owner"); o != nil {
		d.owner = o.Str("login")
	}
	d.repository = d.repoData.Str("name")
	d.isPrivate = phpjson.Truthy(mustGetMap(d.repoData, "private"))
	switch {
	case d.repoData.Isset("default_branch"):
		d.rootIdentifier = d.repoData.Str("default_branch")
	case d.repoData.Isset("master_branch"):
		d.rootIdentifier = d.repoData.Str("master_branch")
	default:
		d.rootIdentifier = "master"
	}
	d.hasIssues = phpjson.Truthy(mustGetMap(d.repoData, "has_issues"))
	d.isArchived = phpjson.Truthy(mustGetMap(d.repoData, "archived"))
	return nil
}

func (d *GitHubDriver) attemptCloneFallback(cause error) error {
	if !d.allowGitFallback {
		return &RuntimeError{Msg: "Fallback to git driver disabled", Prev: cause}
	}
	d.isPrivate = true
	if err := d.setupGitDriver(d.sshURL()); err != nil {
		d.gitDriver = nil
		d.IO.WriteError("<error>Failed to clone the " + d.sshURL() + " repository, try running in interactive mode so that you can enter your GitHub credentials</error>")
		return err
	}
	return nil
}

func (d *GitHubDriver) setupGitDriver(u string) error {
	if !d.allowGitFallback {
		return &RuntimeError{Msg: "Fallback to git driver disabled"}
	}
	cfg := phpjson.NewMap()
	cfg.Set("url", u)
	gd := NewGitDriver(cfg, d.Env)
	d.gitDriver = gd
	return gd.Initialize()
}

// getContentsGitHub is GitHubDriver::getContents (nil response = the dummy
// 'null' response after a git clone fallback).
func (d *GitHubDriver) getContentsGitHub(u string, fetchingRepoData bool) (*httpdl.Response, error) {
	resp, err := d.getContents(u)
	if err == nil {
		return resp, nil
	}
	var te *httpdl.TransportError
	if !errors.As(err, &te) {
		return nil, err
	}
	g := d.Git()
	switch te.StatusCode {
	case 401, 404:
		if !fetchingRepoData {
			return nil, err
		}
		if g.authorizeGitHub(d.originURL) {
			return d.getContents(u)
		}
		if !d.IO.Interactive {
			if err := d.attemptCloneFallback(err); err != nil {
				return nil, err
			}
			return nil, nil
		}
		var issued, needed []string
		if v := te.Header.Get("X-OAuth-Scopes"); v != "" {
			issued = strings.Split(v, " ")
		}
		if v := te.Header.Get("X-Accepted-OAuth-Scopes"); v != "" {
			needed = strings.Split(v, " ")
		}
		failed := false
		for _, n := range needed {
			if !containsString(issued, n) {
				failed = true
			}
		}
		if len(te.Header) == 0 || len(needed) == 0 || failed {
			if _, err := d.HTTP.GitHubAuthorizeInteractively(d.originURL, "Your GitHub credentials are required to fetch private repository metadata (<info>"+Sanitize(d.url)+"</info>)"); err != nil {
				return nil, err
			}
		}
		return d.getContents(u)
	case 403:
		if !d.Auth.Has(d.originURL) && g.authorizeGitHub(d.originURL) {
			return d.getContents(u)
		}
		rateLimited := httpdl.GitHubIsRateLimited(te.Header)
		if !d.IO.Interactive && fetchingRepoData {
			if cerr := d.attemptCloneFallback(err); cerr != nil {
				if rateLimited {
					d.writeRateLimitError(te)
				}
				return nil, cerr
			}
			return nil, nil
		}
		if !d.Auth.Has(d.originURL) {
			if !d.IO.Interactive {
				d.writeRateLimitError(te)
				return nil, err
			}
			if _, err := d.HTTP.GitHubAuthorizeInteractively(d.originURL, "API limit exhausted. Enter your GitHub credentials to get a larger API limit (<info>"+Sanitize(d.url)+"</info>)"); err != nil {
				return nil, err
			}
			return d.getContents(u)
		}
		if rateLimited {
			d.writeRateLimitError(te)
		}
	}
	return nil, err
}

// writeRateLimitError is GitHubDriver::writeRateLimitError.
func (d *GitHubDriver) writeRateLimitError(te *httpdl.TransportError) {
	if !d.Auth.Has(d.originURL) {
		d.IO.WriteError("<error>GitHub API limit exhausted. Failed to get metadata for the " + Sanitize(d.url) + " repository, try running in interactive mode so that you can enter your GitHub credentials to increase the API limit</error>")
		return
	}
	limit, reset := httpdl.GitHubRateLimit(te.Header)
	d.IO.WriteError(fmt.Sprintf("<error>GitHub API limit (%d calls/hr) is exhausted. You are already authorized so you have to wait until %s before doing more requests</error>", limit, reset))
}
