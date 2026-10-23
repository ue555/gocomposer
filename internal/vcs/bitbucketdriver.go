package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/GitBitbucketDriver.php).
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
)

// BitbucketAPIBase is the Bitbucket API origin (overridden by tests that
// talk to a loopback stub server).
var BitbucketAPIBase = "https://api.bitbucket.org"

// BitbucketDriver is Repository\Vcs\GitBitbucketDriver.
type BitbucketDriver struct {
	base
	owner, repository  string
	hasIssues          bool
	rootIdentifier     string
	rootOK             bool
	tags, branches     []Ref
	tagsOK, branchesOK bool
	branchesURL        string
	tagsURL            string
	homeURL            string
	website            string
	cloneHTTPSURL      string
	repoData           *phpjson.Map
	fallback           *GitDriver
	vcsType            string
}

// NewBitbucketDriver is new GitBitbucketDriver($repoConfig, ...).
func NewBitbucketDriver(repoConfig *phpjson.Map, env *Env) Driver {
	d := &BitbucketDriver{base: newBase(repoConfig, env)}
	d.reader = d
	return d
}

var reBitbucketRepo = regexp.MustCompile(`(?i)^https?://bitbucket\.org/([^/]+)/([^/]+?)(?:\.git|/?)?$`)

// Initialize is GitBitbucketDriver::initialize.
func (d *BitbucketDriver) Initialize() error {
	m := reBitbucketRepo.FindStringSubmatch(d.url)
	if m == nil {
		return &InvalidArgumentError{Msg: fmt.Sprintf("The Bitbucket repository URL %s is invalid. It must be the HTTPS URL of a Bitbucket repository.", Sanitize(d.url))}
	}
	d.owner, d.repository = m[1], m[2]
	d.originURL = "bitbucket.org"
	d.cache = cache.New(strings.Join([]string{d.Config.GetString("cache-repo-dir"), d.originURL, d.owner, d.repository}, "/"), "a-z0-9._", d.Config.GetBool("cache-read-only"))
	return nil
}

// URL is GitBitbucketDriver::getUrl.
func (d *BitbucketDriver) URL() string {
	if d.fallback != nil {
		return d.fallback.URL()
	}
	return d.cloneHTTPSURL
}

var reCloneUser = regexp.MustCompile(`https://([^@]+@)?`)

func (d *BitbucketDriver) getRepoData() (bool, error) {
	resource := fmt.Sprintf("%s/2.0/repositories/%s/%s?fields=-project%%2C-owner", BitbucketAPIBase, d.owner, d.repository)
	resp, err := d.fetch(resource, true)
	if err != nil {
		return false, err
	}
	if d.fallback != nil {
		return false, nil
	}
	v, err := decodeJSON(resp)
	if err != nil {
		return false, err
	}
	data := phpjson.ToMap(v)
	if data == nil {
		data = phpjson.NewMap()
	}
	links := data.Map("links")
	if links == nil {
		links = phpjson.NewMap()
	}
	phpjson.Each(mustGetMap(links, "clone"), func(_ string, cv phpjson.Value) {
		if c := phpjson.ToMap(cv); c != nil && c.Str("name") == "https" {
			d.cloneHTTPSURL = reCloneUser.ReplaceAllString(c.Str("href"), "https://")
		}
	})
	d.hasIssues = phpjson.Truthy(mustGetMap(data, "has_issues"))
	href := func(k string) string {
		if m := links.Map(k); m != nil {
			return m.Str("href")
		}
		return ""
	}
	d.branchesURL, d.tagsURL, d.homeURL = href("branches"), href("tags"), href("html")
	d.website = data.Str("website")
	d.vcsType = data.Str("scm")
	d.repoData = data
	return true, nil
}

// ComposerInformation is GitBitbucketDriver::getComposerInformation.
func (d *BitbucketDriver) ComposerInformation(identifier string) (*phpjson.Map, error) {
	if d.fallback != nil {
		return d.fallback.ComposerInformation(identifier)
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
			// array_key_exists($label, $tags = getTags()) elseif branches
			hash := ""
			for _, get := range []func() ([]Ref, error){d.Tags, d.Branches} {
				list, err := get()
				if err != nil {
					return nil, err
				}
				for _, r := range list {
					if r.Name == label {
						hash = r.ID
						break
					}
				}
				if hash != "" {
					break
				}
			}
			if hash == "" {
				support.Set("source", fmt.Sprintf("https://%s/%s/%s/src", d.originURL, d.owner, d.repository))
			} else {
				support.Set("source", fmt.Sprintf("https://%s/%s/%s/src/%s/?at=%s", d.originURL, d.owner, d.repository, hash, label))
			}
		}
		if !support.Isset("issues") && d.hasIssues {
			support.Set("issues", fmt.Sprintf("https://%s/%s/%s/issues", d.originURL, d.owner, d.repository))
		}
		composer.Set("support", support)
		if !composer.Isset("homepage") {
			if d.website == "" {
				composer.Set("homepage", d.homeURL)
			} else {
				composer.Set("homepage", d.website)
			}
		}
	}
	d.infoCache[identifier] = infoEntry{composer}
	return composer, nil
}

func (d *BitbucketDriver) resolveBranch(identifier string) (string, error) {
	if strings.Contains(identifier, "/") {
		branches, err := d.Branches()
		if err != nil {
			return "", err
		}
		for _, b := range branches {
			if b.Name == identifier {
				return b.ID, nil
			}
		}
	}
	return identifier, nil
}

// FileContent is GitBitbucketDriver::getFileContent.
func (d *BitbucketDriver) FileContent(file, identifier string) (string, bool, error) {
	if d.fallback != nil {
		return d.fallback.FileContent(file, identifier)
	}
	identifier, err := d.resolveBranch(identifier)
	if err != nil {
		return "", false, err
	}
	resp, err := d.fetch(fmt.Sprintf("%s/2.0/repositories/%s/%s/src/%s/%s", BitbucketAPIBase, d.owner, d.repository, identifier, file), false)
	if err != nil {
		return "", false, err
	}
	return string(resp.Body), true, nil
}

// ChangeDate is GitBitbucketDriver::getChangeDate.
func (d *BitbucketDriver) ChangeDate(identifier string) (time.Time, bool, error) {
	if d.fallback != nil {
		return d.fallback.ChangeDate(identifier)
	}
	identifier, err := d.resolveBranch(identifier)
	if err != nil {
		return time.Time{}, false, err
	}
	resp, err := d.fetch(fmt.Sprintf("%s/2.0/repositories/%s/%s/commit/%s?fields=date", BitbucketAPIBase, d.owner, d.repository, identifier), false)
	if err != nil {
		return time.Time{}, false, err
	}
	v, err := decodeJSON(resp)
	if err != nil {
		return time.Time{}, false, err
	}
	date := ""
	if m := phpjson.ToMap(v); m != nil {
		date = m.Str("date")
	}
	t, err := parsePHPDateTime(date)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

// Source is GitBitbucketDriver::getSource.
func (d *BitbucketDriver) Source(identifier string) *phpjson.Map {
	if d.fallback != nil {
		return d.fallback.Source(identifier)
	}
	m := phpjson.NewMap()
	m.Set("type", d.vcsType)
	m.Set("url", d.URL())
	m.Set("reference", identifier)
	return m
}

// Dist is GitBitbucketDriver::getDist.
func (d *BitbucketDriver) Dist(identifier string) (*phpjson.Map, error) {
	if d.fallback != nil {
		return d.fallback.Dist(identifier)
	}
	m := phpjson.NewMap()
	m.Set("type", "zip")
	m.Set("url", fmt.Sprintf("https://bitbucket.org/%s/%s/get/%s.zip", d.owner, d.repository, identifier))
	m.Set("reference", identifier)
	m.Set("shasum", "")
	return m, nil
}

func (d *BitbucketDriver) refs(startURL, fields string) ([]Ref, error) {
	var refs []Ref
	resource := startURL + "?pagelen=100&fields=" + strings.ReplaceAll(fields, ",", "%2C") + "&sort=-target.date"
	for resource != "" {
		resp, err := d.fetch(resource, false)
		if err != nil {
			return nil, err
		}
		v, err := decodeJSON(resp)
		if err != nil {
			return nil, err
		}
		data := phpjson.ToMap(v)
		if data == nil {
			data = phpjson.NewMap()
		}
		phpjson.Each(mustGetMap(data, "values"), func(_ string, ev phpjson.Value) {
			if e := phpjson.ToMap(ev); e != nil {
				name, _ := phpjson.ScalarString(mustGetMap(e, "name"))
				hash := ""
				if t := e.Map("target"); t != nil {
					hash = t.Str("hash")
				}
				refs = setRef(refs, name, hash)
			}
		})
		if next := data.Str("next"); phpjson.Truthy(mustGetMap(data, "next")) {
			resource = next
		} else {
			resource = ""
		}
	}
	return refs, nil
}

// Tags is GitBitbucketDriver::getTags.
func (d *BitbucketDriver) Tags() ([]Ref, error) {
	if d.fallback != nil {
		return d.fallback.Tags()
	}
	if !d.tagsOK {
		refs, err := d.refs(d.tagsURL, "values.name,values.target.hash,next")
		if err != nil {
			return nil, err
		}
		d.tags, d.tagsOK = refs, true
	}
	return d.tags, nil
}

// Branches is GitBitbucketDriver::getBranches.
func (d *BitbucketDriver) Branches() ([]Ref, error) {
	if d.fallback != nil {
		return d.fallback.Branches()
	}
	if !d.branchesOK {
		refs, err := d.refs(d.branchesURL, "values.name,values.target.hash,values.heads,next")
		if err != nil {
			return nil, err
		}
		d.branches, d.branchesOK = refs, true
	}
	return d.branches, nil
}

// fetch is fetchWithOAuthCredentials (nil response = dummy 'null' after a
// clone fallback).
func (d *BitbucketDriver) fetch(u string, fetchingRepoData bool) (*httpdl.Response, error) {
	resp, err := d.getContents(u)
	if err == nil {
		return resp, nil
	}
	var te *httpdl.TransportError
	if !errors.As(err, &te) {
		return nil, err
	}
	if te.StatusCode == 403 || te.StatusCode == 404 || (te.StatusCode == 401 && strings.HasPrefix(te.Msg, "Could not authenticate against")) {
		if !d.Auth.Has(d.originURL) && d.Git().authorizeBitbucket(d.originURL) {
			return d.getContents(u)
		}
		if !d.IO.Interactive && fetchingRepoData {
			ssh := "git@ssh." + d.originURL + ":" + d.owner + "/" + d.repository + ".git"
			cfg := phpjson.NewMap()
			cfg.Set("url", ssh)
			gd := NewGitDriver(cfg, d.Env)
			if ierr := gd.Initialize(); ierr != nil {
				d.IO.WriteError("<error>Failed to clone the " + ssh + " repository, try running in interactive mode so that you can enter your Bitbucket OAuth consumer credentials</error>")
				return nil, ierr
			}
			d.fallback = gd
			return nil, nil
		}
	}
	return nil, err
}

// RootIdentifier is GitBitbucketDriver::getRootIdentifier.
func (d *BitbucketDriver) RootIdentifier() (string, error) {
	if d.fallback != nil {
		return d.fallback.RootIdentifier()
	}
	if !d.rootOK {
		ok, err := d.getRepoData()
		if err != nil {
			return "", err
		}
		if !ok {
			if d.fallback == nil {
				return "", &RuntimeError{Msg: "A fallback driver should be setup if getRepoData returns false"}
			}
			return d.fallback.RootIdentifier()
		}
		if d.vcsType != "git" {
			return "", &RuntimeError{Msg: Sanitize(d.url) + " does not appear to be a git repository, use " + d.cloneHTTPSURL + " but remember that Bitbucket no longer supports the mercurial repositories. https://bitbucket.org/blog/sunsetting-mercurial-support-in-bitbucket"}
		}
		d.rootIdentifier = "master"
		if mb := d.repoData.Map("mainbranch"); mb != nil && mb.Isset("name") {
			d.rootIdentifier = mb.Str("name")
		}
		d.rootOK = true
	}
	return d.rootIdentifier, nil
}
