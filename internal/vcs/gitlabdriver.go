package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/GitLabDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// GitLabDriver is Repository\Vcs\GitLabDriver.
type GitLabDriver struct {
	base
	scheme               string
	namespace            string
	repository           string
	project              *phpjson.Map
	projectLoaded        bool
	commits              map[string]*phpjson.Map
	tags, branches       []Ref
	tagsOK, branchesOK   bool
	gitDriver            *GitDriver
	protocol             string
	isPrivate            bool
	hasNonstandardOrigin bool
}

// NewGitLabDriver is new GitLabDriver($repoConfig, ...).
func NewGitLabDriver(repoConfig *phpjson.Map, env *Env) Driver {
	d := &GitLabDriver{base: newBase(repoConfig, env), commits: map[string]*phpjson.Map{}, isPrivate: true}
	d.reader = d
	return d
}

// getContents is VcsDriver::getContents (HttpDownloader::get with the
// repository options).
func (b *base) getContents(u string) (*httpdl.Response, error) {
	if b.HTTP == nil {
		return nil, &RuntimeError{Msg: "No HTTP downloader available to fetch " + Sanitize(u)}
	}
	if opts := b.repoConfig.Map("options"); opts != nil {
		return b.HTTP.Request(u, opts)
	}
	return b.HTTP.Get(u)
}

func decodeJSON(resp *httpdl.Response) (phpjson.Value, error) {
	return jsonfile.Parse(resp.Body, resp.URL)
}

// Initialize is GitLabDriver::initialize.
func (d *GitLabDriver) Initialize() error {
	g := parseGitLabURL(d.url)
	if g == nil {
		return &InvalidArgumentError{Msg: fmt.Sprintf("The GitLab repository URL %s is invalid. It must be the HTTP URL of a GitLab project.", Sanitize(d.url))}
	}
	switch g.scheme {
	case "https", "http":
		d.scheme = g.scheme
	default:
		if v, ok := d.repoConfig.Get("secure-http"); ok && v == false {
			d.scheme = "http"
		} else {
			d.scheme = "https"
		}
	}
	parts := g.parts
	origin, ok := gitlabDetermineOrigin(configList(d.Config, "gitlab-domains"), g.domain, &parts, g.port, g.hasPort)
	if !ok {
		return &RuntimeError{Msg: "It should not be possible to create a gitlab driver with an unparsable origin URL (" + Sanitize(d.url) + ")"}
	}
	d.originURL = origin
	if p, ok := d.Config.Get("gitlab-protocol").(string); ok {
		if p != "git" && p != "http" && p != "https" {
			return &RuntimeError{Msg: "gitlab-protocol must be one of git, http."}
		}
		if p == "git" {
			d.protocol = "ssh"
		} else {
			d.protocol = "http"
		}
	}
	if strings.ContainsAny(d.originURL, ":/") {
		d.hasNonstandardOrigin = true
	}
	d.namespace = strings.Join(parts, "/")
	d.repository = strings.TrimSuffix(g.repo, ".git")
	d.cache = cache.New(d.Config.GetString("cache-repo-dir")+"/"+d.originURL+"/"+d.namespace+"/"+d.repository, "a-z0-9._", d.Config.GetBool("cache-read-only"))
	return d.fetchProject()
}

func urlEncodeAll(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func (d *GitLabDriver) apiURL() string {
	return d.scheme + "://" + d.originURL + "/api/v4/projects/" + urlEncodeAll(d.namespace) + "%2F" + urlEncodeAll(d.repository)
}

// ComposerInformation is GitLabDriver::getComposerInformation.
func (d *GitLabDriver) ComposerInformation(identifier string) (*phpjson.Map, error) {
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
			composer = phpjson.ToMap(v)
			loaded = true
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
		if v, ok := composer.Get("support"); ok {
			if m := phpjson.ToMap(v); m != nil && phpjson.IsArray(v) {
				support = m.Clone()
			}
		}
		webURL := ""
		if d.project != nil && d.project.Isset("web_url") {
			webURL = d.project.Str("web_url")
		}
		if !support.Isset("source") && webURL != "" {
			label, err := refLabelLazy(identifier, d.Tags, d.Branches)
			if err != nil {
				return nil, err
			}
			support.Set("source", webURL+"/-/tree/"+label)
		}
		if !support.Isset("issues") && d.project != nil && phpjson.Truthy(mustGetMap(d.project, "issues_enabled")) && webURL != "" {
			support.Set("issues", webURL+"/-/issues")
		}
		if composer.Has("support") || support.Len() > 0 {
			composer.Set("support", support)
		}
		if !composer.Isset("abandoned") && d.project != nil && phpjson.Truthy(mustGetMap(d.project, "archived")) {
			composer.Set("abandoned", true)
		}
	}
	d.infoCache[identifier] = infoEntry{composer}
	return composer, nil
}

func mustGetMap(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

// refLabel is array_search($id, tags) ?: array_search($id, branches) ?: $id.
func refLabel(identifier string, tags, branches []Ref) string {
	l, _ := refLabelLazy(identifier, func() ([]Ref, error) { return tags, nil }, func() ([]Ref, error) { return branches, nil })
	return l
}

// refLabelLazy is `array_search($identifier, $this->getTags()) ?:
// array_search($identifier, $this->getBranches()) ?: $identifier`: the
// first matching name of each list counts, a falsy one ("" or "0") falls
// through, and the branches are only fetched when the tags have no match.
func refLabelLazy(identifier string, tags, branches func() ([]Ref, error)) (string, error) {
	for _, get := range []func() ([]Ref, error){tags, branches} {
		list, err := get()
		if err != nil {
			return "", err
		}
		for _, r := range list {
			if r.ID == identifier {
				if r.Name != "" && r.Name != "0" {
					return r.Name, nil
				}
				break
			}
		}
	}
	return identifier, nil
}

func (d *GitLabDriver) mustTags() []Ref     { t, _ := d.Tags(); return t }
func (d *GitLabDriver) mustBranches() []Ref { b, _ := d.Branches(); return b }

var reSha40Anywhere = regexp.MustCompile(`(?i)[a-f0-9]{40}`)

// FileContent is GitLabDriver::getFileContent.
func (d *GitLabDriver) FileContent(file, identifier string) (string, bool, error) {
	if d.gitDriver != nil {
		return d.gitDriver.FileContent(file, identifier)
	}
	if !reSha40Anywhere.MatchString(identifier) {
		branches, err := d.Branches()
		if err != nil {
			return "", false, err
		}
		for _, b := range branches {
			if b.Name == identifier {
				identifier = b.ID
				break
			}
		}
	}
	resp, err := d.getContentsGitLab(d.apiURL()+"/repository/files/"+urlEncodeAll(file)+"/raw?ref="+identifier, false)
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) && te.StatusCode == 404 {
			return "", false, nil
		}
		return "", false, err
	}
	return string(resp.Body), true, nil
}

// ChangeDate is GitLabDriver::getChangeDate.
func (d *GitLabDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	if d.gitDriver != nil {
		return d.gitDriver.ChangeDate(identifier)
	}
	if c, ok := d.commits[identifier]; ok {
		t, err := parsePHPDateTime(c.Str("committed_date"))
		if err != nil {
			return time.Time{}, false, err
		}
		return t, true, nil
	}
	return time.Time{}, false, nil
}

// parsePHPDateTime parses the ISO 8601 dates returned by the forge APIs.
func parsePHPDateTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.000-07:00", "2006-01-02T15:04:05Z07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	// the other formats new \DateTimeImmutable() accepts (e.g. "2012-09-10")
	if t, ok := pkg.ParsePHPDate(s); ok {
		return t, nil
	}
	return time.Time{}, &RuntimeError{Msg: "Failed to parse time string (" + s + ")"}
}

func (d *GitLabDriver) repositoryURL() string {
	if d.protocol != "" {
		return d.project.Str(d.protocol + "_url_to_repo")
	}
	if d.isPrivate {
		return d.project.Str("ssh_url_to_repo")
	}
	return d.project.Str("http_url_to_repo")
}

// URL is GitLabDriver::getUrl.
func (d *GitLabDriver) URL() string {
	if d.gitDriver != nil {
		return d.gitDriver.URL()
	}
	return d.project.Str("web_url")
}

// Dist is GitLabDriver::getDist.
func (d *GitLabDriver) Dist(identifier string) (*phpjson.Map, error) {
	m := phpjson.NewMap()
	m.Set("type", "zip")
	m.Set("url", d.apiURL()+"/repository/archive.zip?sha="+identifier)
	m.Set("reference", identifier)
	m.Set("shasum", "")
	return m, nil
}

// Source is GitLabDriver::getSource.
func (d *GitLabDriver) Source(identifier string) *phpjson.Map {
	if d.gitDriver != nil {
		return d.gitDriver.Source(identifier)
	}
	m := phpjson.NewMap()
	m.Set("type", "git")
	m.Set("url", d.repositoryURL())
	m.Set("reference", identifier)
	return m
}

// RootIdentifier is GitLabDriver::getRootIdentifier.
func (d *GitLabDriver) RootIdentifier() (string, error) {
	if d.gitDriver != nil {
		return d.gitDriver.RootIdentifier()
	}
	return d.project.Str("default_branch"), nil
}

// Branches is GitLabDriver::getBranches.
func (d *GitLabDriver) Branches() ([]Ref, error) {
	if d.gitDriver != nil {
		return d.gitDriver.Branches()
	}
	if !d.branchesOK {
		refs, err := d.references("branches")
		if err != nil {
			return nil, err
		}
		d.branches, d.branchesOK = refs, true
	}
	return d.branches, nil
}

// Tags is GitLabDriver::getTags.
func (d *GitLabDriver) Tags() ([]Ref, error) {
	if d.gitDriver != nil {
		return d.gitDriver.Tags()
	}
	if !d.tagsOK {
		refs, err := d.references("tags")
		if err != nil {
			return nil, err
		}
		d.tags, d.tagsOK = refs, true
	}
	return d.tags, nil
}

func (d *GitLabDriver) references(typ string) ([]Ref, error) {
	const perPage = 100
	resource := d.apiURL() + "/repository/" + typ + "?per_page=100"
	var refs []Ref
	for resource != "" {
		resp, err := d.getContentsGitLab(resource, false)
		if err != nil {
			return nil, err
		}
		v, err := decodeJSON(resp)
		if err != nil {
			return nil, err
		}
		n := 0
		phpjson.Each(v, func(_ string, dv phpjson.Value) {
			n++
			m := phpjson.ToMap(dv)
			if m == nil {
				return
			}
			c := m.Map("commit")
			if c == nil {
				return
			}
			name, _ := phpjson.ScalarString(mustGetMap(m, "name"))
			refs = setRef(refs, name, c.Str("id"))
			d.commits[c.Str("id")] = c
		})
		if n >= perPage {
			resource = nextPage(resp)
		} else {
			resource = ""
		}
	}
	return refs, nil
}

var reNextLink = regexp.MustCompile(`<(.+?)>; *rel="next"`)

// nextPage is getNextPage (Link header rel="next").
func nextPage(resp *httpdl.Response) string {
	h := resp.Header.Get("Link")
	if h == "" {
		return ""
	}
	for _, link := range strings.Split(h, ",") {
		if m := reNextLink.FindStringSubmatch(link); m != nil {
			return m[1]
		}
	}
	return ""
}

func (d *GitLabDriver) fetchProject() error {
	if d.projectLoaded {
		return nil
	}
	resp, err := d.getContentsGitLab(d.apiURL(), true)
	if err != nil {
		return err
	}
	d.projectLoaded = true
	if resp == nil {
		// clone fallback: Response 'null'
		d.project = nil
		return nil
	}
	v, err := decodeJSON(resp)
	if err != nil {
		return err
	}
	d.project = phpjson.ToMap(v)
	if d.project != nil && d.project.Isset("visibility") {
		d.isPrivate = d.project.Str("visibility") != "public"
	} else {
		d.isPrivate = false
	}
	return nil
}

func (d *GitLabDriver) attemptCloneFallback() error {
	var u string
	if !d.isPrivate {
		u = d.scheme + "://" + d.originURL + "/" + d.namespace + "/" + d.repository + ".git"
	} else if d.hasNonstandardOrigin {
		u = "ssh://git@" + d.originURL + "/" + d.namespace + "/" + d.repository + ".git"
	} else {
		u = "git@" + d.originURL + ":" + d.namespace + "/" + d.repository + ".git"
	}
	cfg := phpjson.NewMap()
	cfg.Set("url", u)
	gd := NewGitDriver(cfg, d.Env)
	if err := gd.Initialize(); err != nil {
		d.IO.WriteError("<error>Failed to clone the " + Sanitize(u) + " repository, try running in interactive mode so that you can enter your credentials</error>")
		return err
	}
	d.gitDriver = gd
	return nil
}

// getContentsGitLab is GitLabDriver::getContents (nil response = the
// dummy 'null' response after a git clone fallback).
func (d *GitLabDriver) getContentsGitLab(u string, fetchingRepoData bool) (*httpdl.Response, error) {
	resp, err := d.getContents(u)
	if err == nil && fetchingRepoData {
		v, jerr := decodeJSON(resp)
		if jerr != nil {
			return nil, jerr
		}
		json := phpjson.ToMap(v)
		if json == nil {
			json = phpjson.NewMap()
		}
		if !json.Isset("default_branch") && json.Isset("permissions") {
			d.isPrivate = json.Str("visibility") != "public"
			more := false
			phpjson.Each(mustGetMap(json, "permissions"), func(_ string, pv phpjson.Value) {
				if pm := phpjson.ToMap(pv); pm != nil {
					if lvl, ok := pm.Get("access_level"); ok {
						if n, ok := lvl.(int64); ok && n >= 20 {
							more = true
						}
					}
				}
			})
			if !more {
				d.IO.WriteError("<warning>GitLab token with Guest or Planner only access detected</warning>")
				if err := d.attemptCloneFallback(); err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
		if !json.Isset("default_branch") {
			if json.Str("repository_access_level") == "disabled" {
				err = &httpdl.TransportError{Msg: "The GitLab repository is disabled in the project", StatusCode: 400}
			} else {
				if phpjson.Truthy(mustGetMap(json, "id")) {
					d.isPrivate = false
				}
				err = &httpdl.TransportError{Msg: "GitLab API seems to not be authenticated as it did not return a default_branch", StatusCode: 401}
			}
		}
	}
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
		if g.authorizeGitLab(d.originURL) {
			return d.getContents(u)
		}
		if d.HTTP.GitLabIsOAuthExpired(d.originURL) {
			ok, rerr := d.HTTP.GitLabAuthorizeOAuthRefresh(d.scheme, d.originURL)
			if rerr != nil {
				return nil, rerr
			}
			if ok {
				return d.getContents(u)
			}
		}
		if !d.IO.Interactive {
			if err := d.attemptCloneFallback(); err != nil {
				return nil, err
			}
			return nil, nil
		}
		d.IO.WriteError("<warning>Failed to download " + d.namespace + "/" + d.repository + ":" + err.Error() + "</warning>")
		if _, aerr := d.HTTP.GitLabAuthorizeInteractively(d.scheme, d.originURL, "Your credentials are required to fetch private repository metadata (<info>"+Sanitize(d.url)+"</info>)"); aerr != nil {
			return nil, aerr
		}
		return d.getContents(u)
	case 403:
		if !d.Auth.Has(d.originURL) && g.authorizeGitLab(d.originURL) {
			return d.getContents(u)
		}
		if !d.IO.Interactive && fetchingRepoData {
			if err := d.attemptCloneFallback(); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	return nil, err
}
