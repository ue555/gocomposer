package repo

import (
	sha256pkg "crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// Composer is ComposerRepository for the v1 inline/includes format and the
// v2 metadata-url protocol, and the v1 providers-url / provider-includes /
// providers-lazy-url protocol (composer_providers.go).
type Composer struct {
	*Array
	url                string
	baseURL            string
	options            *phpjson.Map
	repoConfig         *phpjson.Map
	cache              *cache.Cache
	dl                 *httpdl.Downloader
	loader             *pkg.Loader
	Warn               func(string)
	rootData           *phpjson.Map
	rootLoaded         bool
	notifyURL          string
	lazyProvidersURL   string
	providersURL       string
	providersAPIURL    string
	listURL            string
	providers          *providerListing
	hasProviders       bool
	hasPartialPackages bool
	partialByName      map[string][]*phpjson.Map
	partialOrder       []string
	availablePackages  map[string]bool
	availableNames     []string // availablePackages keys in order
	availablePatterns  []*regexp.Regexp
	hasAvailableList   bool
	degradedMode       bool
	notFound           map[string]bool
	fresh              map[string]bool
	sourceMirrors      map[string][]phpjson.Value
	distMirrors        []phpjson.Value
	displayedIndexWarn bool
	mu                 sync.Mutex
	// Advisories collected from v2 metadata (security-advisories key).
	Advisories  map[string][]*phpjson.Map
	advMetadata bool
	advAPIURL   string
	hasAdvCfg   bool
	// PreFileDownload and PostFileDownload dispatch the plugin events
	// PluginEvents::PRE_FILE_DOWNLOAD/POST_FILE_DOWNLOAD for the metadata
	// files (type "metadata") when the repository has an event dispatcher:
	// the pre event returns the URL and transport options to use.
	PreFileDownload  func(url string, options *phpjson.Map) (string, *phpjson.Map, error)
	PostFileDownload func(url, sha256 string, resp *httpdl.Response) error
}

// metaReq is a metadata request after the pre-file-download event.
type metaReq struct {
	url     string
	headers []string
	// deferPost leaves the post-file-download event to the caller (resp).
	deferPost bool
	resp      *httpdl.Response
}

// preDownload dispatches the pre-file-download event for a metadata URL.
func (r *Composer) preDownload(u string) (*metaReq, error) {
	q := &metaReq{url: u, headers: r.headers()}
	if r.PreFileDownload != nil {
		nu, opts, err := r.PreFileDownload(u, r.options)
		if err != nil {
			return nil, err
		}
		q.url = nu
		q.headers = headersOf(opts)
	}
	return q, nil
}

func (r *Composer) postDownload(q *metaReq, sha256 string, resp *httpdl.Response) error {
	if r.PostFileDownload == nil {
		return nil
	}
	if q.deferPost {
		q.resp = resp
		return nil
	}
	return r.PostFileDownload(q.url, sha256, resp)
}

var (
	reHasScheme     = regexp.MustCompile(`^[\w.]+\??://`)
	rePackagistRoot = regexp.MustCompile(`(?i)^(https?)://packagist\.org/?$`)
	reBaseURL       = regexp.MustCompile(`(?:/[^/\\]+\.json)?(?:[?#].*)?$`)
	reCacheDirName  = regexp.MustCompile(`(?i)[^a-z0-9.]`)
	reDevSuffix     = regexp.MustCompile(`~dev$`)
)

// NewComposer is ComposerRepository::__construct.
func NewComposer(repoConfig *phpjson.Map, cacheRepoDir string, readOnlyCache bool, dl *httpdl.Downloader) (*Composer, error) {
	u := repoConfig.Str("url")
	if !reHasScheme.MatchString(u) {
		if abs, err := filepath.Abs(u); err == nil {
			if _, err := os.Stat(abs); err == nil {
				if real, err := filepath.EvalSymlinks(abs); err == nil {
					abs = real
				}
				u = "file://" + abs
			} else {
				u = "http://" + u
			}
		}
	}
	u = strings.TrimRight(u, "/")
	if u == "" {
		return nil, fmt.Errorf("The repository url must not be an empty string")
	}
	if strings.HasPrefix(u, "https?") {
		u = "https" + u[6:]
	}
	if !strings.Contains(u, "://") {
		return nil, fmt.Errorf("Invalid url given for Composer repository: %s", httpdl.Sanitize(u))
	}
	options := repoConfig.Map("options")
	if options == nil {
		options = phpjson.NewMap()
	}
	r := &Composer{Array: &Array{kind: KindComposer}, repoConfig: repoConfig, options: options, dl: dl, loader: &pkg.Loader{},
		notFound: map[string]bool{}, fresh: map[string]bool{}, sourceMirrors: map[string][]phpjson.Value{}, Advisories: map[string][]*phpjson.Map{}}
	r.url = u
	if m := rePackagistRoot.FindStringSubmatch(u); m != nil {
		r.url = m[1] + "://repo.packagist.org"
	}
	r.baseURL = strings.TrimRight(reBaseURL.ReplaceAllString(r.url, ""), "/")
	r.cache = cache.New(cacheRepoDir+"/"+reCacheDirName.ReplaceAllString(httpdl.Sanitize(r.url), "-"), "a-z0-9.$~_", readOnlyCache)
	r.Array.init = r.initialize
	r.Array.self = r
	r.Array.name = "composer repo (" + httpdl.Sanitize(r.url) + ")"
	return r, nil
}

func (r *Composer) RepoName() string { return "composer repo (" + httpdl.Sanitize(r.url) + ")" }

func (r *Composer) warn(msg string) {
	if r.Warn != nil {
		r.Warn(msg)
	}
}

// URL returns the normalized repository URL.
func (r *Composer) URL() string { return r.url }

func (r *Composer) packagesJSONURL() string {
	path := r.url
	if i := strings.Index(path, "://"); i >= 0 {
		path = path[i+3:]
	}
	if j := strings.IndexByte(path, '/'); j >= 0 && strings.Contains(path[j:], ".json") {
		return r.url
	}
	return r.url + "/packages.json"
}

var reURLOrigin = regexp.MustCompile(`^[^:]+://[^/]*`)

func (r *Composer) canonicalizeURL(u string) string {
	if strings.HasPrefix(u, "/") {
		if m := reURLOrigin.FindString(r.url); m != "" {
			return m + u
		}
		return r.url
	}
	return u
}

func (r *Composer) loadRootServerFile() (*phpjson.Map, error) {
	if r.rootLoaded {
		return r.rootData, nil
	}
	var data *phpjson.Map
	if cached, ok := r.cache.Read("packages.json"); ok {
		if v, err := phpjson.Decode(cached); err == nil {
			cm := phpjson.ToMap(v)
			if lm := cm.Str("last-modified"); lm != "" {
				resp, err := r.fetchFileIfLastModified(r.packagesJSONURL(), "packages.json", lm)
				if err != nil {
					return nil, err
				}
				if resp == nil {
					data = cm
				} else {
					data = resp
				}
			}
		}
	}
	if data == nil {
		d, err := r.fetchFile(r.packagesJSONURL(), "packages.json", "", true)
		if err != nil {
			return nil, err
		}
		data = d
	}
	if s := data.Str("notify-batch"); s != "" {
		r.notifyURL = r.canonicalizeURL(s)
	} else if s := data.Str("notify"); s != "" {
		r.notifyURL = r.canonicalizeURL(s)
	}
	if mirrors, ok := data.Get("mirrors"); ok && phpjson.Truthy(mirrors) {
		phpjson.Each(mirrors, func(_ string, mv phpjson.Value) {
			m := phpjson.ToMap(mv)
			preferred, _ := m.Get("preferred")
			if g := m.Str("git-url"); g != "" {
				e := phpjson.NewMap()
				e.Set("url", g)
				e.Set("preferred", phpjson.Truthy(preferred))
				r.sourceMirrors["git"] = append(r.sourceMirrors["git"], e)
			}
			if h := m.Str("hg-url"); h != "" {
				e := phpjson.NewMap()
				e.Set("url", h)
				e.Set("preferred", phpjson.Truthy(preferred))
				r.sourceMirrors["hg"] = append(r.sourceMirrors["hg"], e)
			}
			if d := m.Str("dist-url"); d != "" {
				e := phpjson.NewMap()
				e.Set("url", r.canonicalizeURL(d))
				e.Set("preferred", phpjson.Truthy(preferred))
				r.distMirrors = append(r.distMirrors, e)
			}
		})
	}
	if s := data.Str("providers-lazy-url"); s != "" {
		r.lazyProvidersURL = r.canonicalizeURL(s)
		r.hasProviders = true
		pk, _ := data.Get("packages")
		r.hasPartialPackages = phpjson.Truthy(pk) && phpjson.IsArray(pk)
	}
	if s := data.Str("metadata-url"); s != "" {
		r.lazyProvidersURL = r.canonicalizeURL(s)
		r.providersURL = ""
		r.hasProviders = false
		pk, _ := data.Get("packages")
		r.hasPartialPackages = phpjson.Truthy(pk) && phpjson.IsArray(pk)
		if ap, ok := data.Get("available-packages"); ok && phpjson.Truthy(ap) {
			r.availablePackages = map[string]bool{}
			r.availableNames = nil
			phpjson.Each(ap, func(_ string, v phpjson.Value) {
				s, _ := v.(string)
				if !r.availablePackages[strings.ToLower(s)] {
					r.availableNames = append(r.availableNames, strings.ToLower(s))
				}
				r.availablePackages[strings.ToLower(s)] = true
			})
			r.hasAvailableList = true
		}
		if ap, ok := data.Get("available-package-patterns"); ok && phpjson.Truthy(ap) {
			phpjson.Each(ap, func(_ string, v phpjson.Value) {
				s, _ := v.(string)
				r.availablePatterns = append(r.availablePatterns, regexp.MustCompile(`(?i)^`+strings.ReplaceAll(regexp.QuoteMeta(s), `\*`, `.*`)+`$`))
			})
			r.hasAvailableList = true
		}
		data.Delete("providers-url")
		data.Delete("providers")
		data.Delete("providers-includes")
		if sa := data.Map("security-advisories"); sa != nil {
			r.hasAdvCfg = true
			if v, ok := sa.Get("metadata"); ok {
				r.advMetadata = phpjson.Truthy(v)
			}
			if u, ok := sa.Get("api-url"); ok {
				if us, ok := u.(string); ok {
					r.advAPIURL = r.canonicalizeURL(us)
				}
			}
		}
	}
	if s := data.Str("providers-url"); s != "" {
		r.providersURL = r.canonicalizeURL(s)
		r.hasProviders = true
	}
	if p, ok := data.Get("providers"); ok && phpjson.Truthy(p) {
		r.hasProviders = true
	}
	if p, ok := data.Get("providers-includes"); ok && phpjson.Truthy(p) {
		r.hasProviders = true
	}
	if s := data.Str("list"); s != "" {
		r.listURL = r.canonicalizeURL(s)
	}
	if s := data.Str("providers-api"); s != "" {
		r.providersAPIURL = r.canonicalizeURL(s)
	}
	r.rootData = data
	r.rootLoaded = true
	return data, nil
}

// UnsupportedError is returned for repository features gocomposer does not implement.
type UnsupportedError struct{ Msg string }

func (e *UnsupportedError) Error() string { return e.Msg }

// checkSupported reports repository protocols gocomposer does not implement
// (none left: v1 providers are handled by whatProvides).
func (r *Composer) checkSupported() error { return nil }

func (r *Composer) fetchFile(u, cacheKey, sha256 string, storeLastModified bool) (*phpjson.Map, error) {
	if i := strings.IndexByte(u, '$'); i > 0 && regexp.MustCompile(`(?i)^https?://`).MatchString(u) {
		u = u[:i] + "%24" + u[i+1:]
	}
	var resp *httpdl.Response
	var q *metaReq
	var err error
	for retries := 2; ; retries-- {
		q, err = r.preDownload(u)
		if err != nil {
			return nil, err
		}
		resp, err = r.dl.Get(q.url, q.headers...)
		if err != nil || sha256 == "" || sha256Hex(resp.Body) == sha256 {
			break
		}
		if retries > 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return nil, &RepositorySecurityError{"The contents of " + httpdl.Sanitize(q.url) + " do not match its signature. This could indicate a man-in-the-middle attack or e.g. antivirus software corrupting files. Try running composer again and report this if you think it is a mistake."}
	}
	u = q.url
	if err != nil {
		if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode == 404 {
			return nil, err
		}
		if _, ok := err.(*httpdl.TransportError); ok && strings.HasPrefix(err.Error(), "Your configuration does not allow") {
			return nil, err
		}
		if cacheKey != "" {
			if contents, ok := r.cache.Read(cacheKey); ok {
				if !r.degradedMode {
					r.warn("<warning>" + httpdl.Sanitize(r.url) + " could not be fully loaded (" + err.Error() + "), package information was loaded from the local cache and may be out of date</warning>")
				}
				r.degradedMode = true
				v, derr := phpjson.Decode(contents)
				if derr != nil {
					return nil, derr
				}
				return phpjson.ToMap(v), nil
			}
		}
		return nil, err
	}
	if err := r.postDownload(q, sha256, resp); err != nil {
		return nil, err
	}
	v, err := jsonfile.Parse(resp.Body, httpdl.Sanitize(u))
	if err != nil {
		return nil, err
	}
	r.responseWarnings(v)
	data := phpjson.ToMap(v)
	if data == nil {
		data = phpjson.NewMap()
	}
	if cacheKey != "" && !r.cache.ReadOnly {
		body := resp.Body
		if storeLastModified {
			if lm := resp.Header.Get("Last-Modified"); lm != "" {
				data.Set("last-modified", lm)
				enc, _ := phpjson.Encode(data, 0)
				body = []byte(enc)
			}
		}
		r.cache.Write(cacheKey, body)
	}
	return data, nil
}

// fetchFileIfLastModified returns nil when the cached copy is fresh (304).
func (r *Composer) fetchFileIfLastModified(u, cacheKey, lastModified string) (*phpjson.Map, error) {
	q, err := r.preDownload(u)
	if err != nil {
		return nil, err
	}
	u = q.url
	resp, err := r.dl.Get(u, append(q.headers, "If-Modified-Since: "+lastModified)...)
	if err != nil {
		if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode == 404 {
			return nil, err
		}
		if !r.degradedMode {
			r.warn("<warning>" + httpdl.Sanitize(r.url) + " could not be fully loaded (" + err.Error() + "), package information was loaded from the local cache and may be out of date</warning>")
		}
		r.degradedMode = true
		return nil, nil
	}
	if resp.StatusCode == 304 && len(resp.Body) == 0 {
		return nil, nil
	}
	if err := r.postDownload(q, "", resp); err != nil {
		return nil, err
	}
	v, err := phpjson.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	r.responseWarnings(v)
	data := phpjson.ToMap(v)
	body := resp.Body
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		data.Set("last-modified", lm)
		enc, _ := phpjson.Encode(data, 0)
		body = []byte(enc)
	}
	if !r.cache.ReadOnly {
		r.cache.Write(cacheKey, body)
	}
	return data, nil
}

func (r *Composer) headers() []string { return headersOf(r.options) }

// headersOf returns the http.header transport option as header lines.
func headersOf(options *phpjson.Map) []string {
	var hs []string
	if options == nil {
		return nil
	}
	if http := options.Map("http"); http != nil {
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

func (r *Composer) initialize() error {
	data, err := r.loadRootServerFile()
	if err != nil {
		return err
	}
	if err := r.checkSupported(); err != nil {
		return err
	}
	versions, err := r.loadIncludes(data)
	if err != nil {
		return err
	}
	pkgs, err := r.createPackages(versions, "root file ("+httpdl.Sanitize(r.packagesJSONURL())+")")
	if err != nil {
		return err
	}
	for _, p := range pkgs {
		r.Array.AddPackage(p)
	}
	return nil
}

func (r *Composer) loadIncludes(data *phpjson.Map) ([]*phpjson.Map, error) {
	var out []*phpjson.Map
	if !data.Has("packages") && !data.Has("includes") {
		for _, k := range data.Keys() {
			p := data.Map(k)
			if p == nil {
				continue
			}
			if vs, ok := p.Get("versions"); ok && phpjson.IsArray(vs) {
				phpjson.Each(vs, func(_ string, v phpjson.Value) {
					if m := phpjson.ToMap(v); m != nil {
						out = append(out, m)
					}
				})
			}
		}
		return out, nil
	}
	if pk, ok := data.Get("packages"); ok {
		phpjson.Each(pk, func(name string, versions phpjson.Value) {
			lname := strings.ToLower(name)
			phpjson.Each(versions, func(_ string, v phpjson.Value) {
				m := phpjson.ToMap(v)
				if m == nil {
					return
				}
				out = append(out, m)
				if !r.displayedIndexWarn && lname != strings.ToLower(m.Str("name")) {
					r.displayedIndexWarn = true
					r.warn(fmt.Sprintf("<warning>Warning: the packages key '%s' doesn't match the name defined in the package metadata '%s' in repository %s</warning>", name, m.Str("name"), r.baseURL))
				}
			})
		})
	}
	if inc, ok := data.Get("includes"); ok {
		var err error
		phpjson.Each(inc, func(include string, meta phpjson.Value) {
			if err != nil {
				return
			}
			mm := phpjson.ToMap(meta)
			var included *phpjson.Map
			if mm != nil && mm.Isset("sha1") && r.cache.Sha1(include) == mm.Str("sha1") {
				c, _ := r.cache.Read(include)
				v, derr := phpjson.Decode(c)
				if derr == nil {
					included = phpjson.ToMap(v)
				}
			}
			if included == nil {
				included, err = r.fetchFile(r.baseURL+"/"+include, include, "", false)
				if err != nil {
					return
				}
			}
			more, lerr := r.loadIncludes(included)
			if lerr != nil {
				err = lerr
				return
			}
			out = append(out, more...)
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *Composer) createPackages(versions []*phpjson.Map, source string) ([]*pkg.Package, error) {
	if len(versions) == 0 {
		return nil, nil
	}
	prepared := make([]*phpjson.Map, len(versions))
	for i, v := range versions {
		c := v.Clone()
		if !c.Has("notification-url") {
			if r.notifyURL != "" {
				c.Set("notification-url", r.notifyURL)
			} else {
				c.Set("notification-url", nil)
			}
		}
		prepared[i] = c
	}
	pkgs, err := r.loader.LoadPackages(prepared)
	if err != nil {
		src := ""
		if source != "" {
			src = " from " + source
		}
		return nil, fmt.Errorf("Could not load packages in %s%s: [UnexpectedValueException] %s", r.RepoName(), src, err.Error())
	}
	for _, p := range pkgs {
		if m, ok := r.sourceMirrors[p.SourceType()]; ok {
			p.SetSourceMirrors(m)
		}
		p.SetDistMirrors(r.distMirrors)
		r.configureTransportOptions(p)
	}
	return pkgs, nil
}

func (r *Composer) configureTransportOptions(p *pkg.Package) {
	if r.options.Len() == 0 {
		return
	}
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(r.baseURL) + `(?:[/?#]|$)`)
	if p.DistURL() != "" && re.MatchString(p.DistURL()) {
		p.SetTransportOptions(r.options)
	}
}

func (r *Composer) initializePartialPackages() error {
	if r.partialByName != nil {
		return nil
	}
	data, err := r.loadRootServerFile()
	if err != nil {
		return err
	}
	r.partialByName = map[string][]*phpjson.Map{}
	pk, _ := data.Get("packages")
	phpjson.Each(pk, func(name string, versions phpjson.Value) {
		phpjson.Each(versions, func(_ string, v phpjson.Value) {
			m := phpjson.ToMap(v)
			if m == nil {
				return
			}
			vn := strings.ToLower(m.Str("name"))
			if _, ok := r.partialByName[vn]; !ok {
				r.partialOrder = append(r.partialOrder, vn)
			}
			r.partialByName[vn] = append(r.partialByName[vn], m)
			if !r.displayedIndexWarn && vn != strings.ToLower(name) {
				r.warn(fmt.Sprintf("<warning>Warning: the packages key '%s' doesn't match the name defined in the package metadata '%s' in repository %s</warning>", name, m.Str("name"), r.baseURL))
				r.displayedIndexWarn = true
			}
		})
	})
	return nil
}

func (r *Composer) partialPackages() (bool, error) {
	if r.hasPartialPackages && r.partialByName == nil {
		if err := r.initializePartialPackages(); err != nil {
			return false, err
		}
	}
	return r.hasPartialPackages, nil
}

func (r *Composer) lazyRepoContains(name string) bool {
	if r.availablePackages[name] {
		return true
	}
	for _, re := range r.availablePatterns {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

// LoadPackages is ComposerRepository::loadPackages.
func (r *Composer) LoadPackages(names []NameConstraint, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) (LoadResult, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return LoadResult{}, err
	}
	if err := r.checkSupported(); err != nil {
		return LoadResult{}, err
	}
	hasPartial, err := r.partialPackages()
	if err != nil {
		return LoadResult{}, err
	}
	if !r.hasProviders && !hasPartial && r.lazyProvidersURL == "" {
		return r.Array.LoadPackages(names, acceptable, flags, alreadyLoaded)
	}
	var packages []*pkg.Package
	inPackages := map[*pkg.Package]bool{}
	var namesFound []string
	found := map[string]bool{}
	remaining := append([]NameConstraint(nil), names...)
	if r.hasProviders || hasPartial {
		var keep []NameConstraint
		for _, nc := range remaining {
			if !r.hasProviders {
				if _, ok := r.partialByName[nc.Name]; !ok {
					keep = append(keep, nc)
					continue
				}
			}
			candidates, err := r.whatProvides(nc.Name, acceptable, flags, alreadyLoaded)
			if err != nil {
				return LoadResult{}, err
			}
			var matches []*pkg.Package
			inMatches := map[*pkg.Package]bool{}
			for _, c := range candidates {
				if !found[nc.Name] {
					found[nc.Name] = true
					namesFound = append(namesFound, nc.Name)
				}
				if nc.Constraint == nil || nc.Constraint.Matches(semver.NewConstraint("==", c.Version())) {
					if !inMatches[c] {
						inMatches[c] = true
						matches = append(matches, c)
					}
					if c.AliasOf != nil && !inMatches[c.AliasOf] {
						inMatches[c.AliasOf] = true
						matches = append(matches, c.AliasOf)
					}
				}
			}
			for _, c := range candidates {
				if c.AliasOf != nil && inMatches[c.AliasOf] && !inMatches[c] {
					inMatches[c] = true
					matches = append(matches, c)
				}
			}
			for _, m := range matches {
				if !inPackages[m] {
					inPackages[m] = true
					packages = append(packages, m)
				}
			}
		}
		remaining = keep
	}
	if r.lazyProvidersURL != "" && len(remaining) > 0 {
		if r.hasAvailableList {
			var keep []NameConstraint
			for _, nc := range remaining {
				if r.lazyRepoContains(strings.ToLower(nc.Name)) {
					keep = append(keep, nc)
				}
			}
			remaining = keep
		}
		res, err := r.loadAsyncPackages(remaining, acceptable, flags, alreadyLoaded, true)
		if err != nil {
			return LoadResult{}, err
		}
		for _, p := range res.Packages {
			if !inPackages[p] {
				inPackages[p] = true
				packages = append(packages, p)
			}
		}
		for _, n := range res.NamesFound {
			if !found[n] {
				found[n] = true
				namesFound = append(namesFound, n)
			}
		}
	}
	return LoadResult{NamesFound: namesFound, Packages: packages}, nil
}

// whatProvidesPartial is whatProvides() for inline (partial) packages of a v2 repo.
func (r *Composer) whatProvidesPartial(name string, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) ([]*pkg.Package, error) {
	versions, ok := r.partialByName[name]
	if !ok {
		return nil, nil
	}
	var toLoad []*phpjson.Map
	for _, v := range versions {
		if strings.ToLower(v.Str("name")) != name {
			continue
		}
		v = v.Clone()
		if err := ensureNormalized(v); err != nil {
			return nil, err
		}
		if alreadyLoaded[name][v.Str("version_normalized")] {
			continue
		}
		ok, err := r.isVersionAcceptable(nil, name, v, acceptable, flags)
		if err != nil {
			return nil, err
		}
		if ok {
			toLoad = append(toLoad, v)
		}
	}
	loaded, err := r.createPackages(toLoad, "root file ("+httpdl.Sanitize(r.packagesJSONURL())+")")
	if err != nil {
		return nil, err
	}
	var out []*pkg.Package
	for _, p := range loaded {
		p.Repo = r
		if p.AliasOf != nil {
			p.AliasOf.Repo = r
			out = append(out, p.AliasOf, p)
		} else {
			out = append(out, p)
		}
	}
	return out, nil
}

func ensureNormalized(v *phpjson.Map) error {
	vn, has := v.Get("version_normalized")
	vs, _ := phpjson.ScalarString(mustGetVal(v, "version"))
	if !has || vn == nil {
		n, err := semver.Normalize(vs)
		if err != nil {
			return err
		}
		v.Set("version_normalized", n)
	} else if s, _ := vn.(string); s == semver.DefaultBranchAlias {
		n, err := semver.Normalize(vs)
		if err != nil {
			return err
		}
		v.Set("version_normalized", n)
	}
	return nil
}

func mustGetVal(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

func (r *Composer) isVersionAcceptable(c semver.Constraint, name string, v *phpjson.Map, acceptable, flags map[string]int) (bool, error) {
	versions := []string{v.Str("version_normalized")}
	alias, err := pkg.GetBranchAlias(v)
	if err != nil {
		return false, err
	}
	if alias != "" {
		versions = append(versions, alias)
	}
	for _, ver := range versions {
		if acceptable != nil && flags != nil && !IsPackageAcceptable(acceptable, flags, []string{name}, semver.ParseStability(ver)) {
			continue
		}
		if c != nil && !semver.MatchVersion(c, ver) {
			continue
		}
		return true, nil
	}
	return false, nil
}

type asyncTask struct {
	fileName string
	realName string
	c        semver.Constraint
}

// loadAsyncPackages is ComposerRepository::loadAsyncPackages; withStability
// distinguishes null acceptable stabilities (findPackages) from provided ones.
func (r *Composer) loadAsyncPackages(names []NameConstraint, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool, withStability bool) (LoadResult, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return LoadResult{}, err
	}
	if r.lazyProvidersURL == "" {
		return LoadResult{}, fmt.Errorf("loadAsyncPackages only supports v2 protocol composer repos with a metadata-url")
	}
	type entry struct {
		name string
		c    semver.Constraint
	}
	var list []entry
	index := map[string]int{}
	set := func(n string, c semver.Constraint) {
		if i, ok := index[n]; ok {
			list[i].c = c
			return
		}
		index[n] = len(list)
		list = append(list, entry{n, c})
	}
	for _, nc := range names {
		set(nc.Name, nc.Constraint)
	}
	removed := map[string]bool{}
	for _, nc := range names {
		if !withStability || IsPackageAcceptable(acceptable, flags, []string{nc.Name}, "dev") {
			set(nc.Name+"~dev", nc.Constraint)
		}
		if withStability {
			if _, ok := acceptable["dev"]; ok && len(acceptable) == 1 && len(flags) == 0 {
				removed[nc.Name] = true
			}
		}
	}
	var tasks []asyncTask
	for _, e := range list {
		if removed[e.name] {
			continue
		}
		name := strings.ToLower(e.name)
		realName := reDevSuffix.ReplaceAllString(name, "")
		if IsPlatformPackage(realName) || realName == "__root__" {
			continue
		}
		tasks = append(tasks, asyncTask{name, realName, e.c})
	}
	results := make([]*phpjson.Map, len(tasks))
	sources := make([]string, len(tasks))
	errs := make([]error, len(tasks))
	// Composer dispatches the pre-file-download events while queueing the
	// requests (in order) and the post-file-download events as the
	// responses arrive
	var reqs []*metaReq
	if r.PreFileDownload != nil || r.PostFileDownload != nil {
		reqs = make([]*metaReq, len(tasks))
		for i, t := range tasks {
			q, err := r.prepareMetaReq(t.fileName)
			if err != nil {
				return LoadResult{}, err
			}
			reqs[i] = q
		}
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 12)
	for i, t := range tasks {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, t asyncTask) {
			defer wg.Done()
			defer func() { <-sem }()
			var q *metaReq
			if reqs != nil {
				q = reqs[i]
			}
			results[i], sources[i], errs[i] = r.cachedAsyncDownloadReq(t.fileName, t.realName, q)
		}(i, t)
	}
	wg.Wait()
	for _, q := range reqs {
		if q != nil && q.resp != nil {
			if err := r.PostFileDownload(q.url, "", q.resp); err != nil {
				return LoadResult{}, err
			}
		}
	}
	var packages []*pkg.Package
	inPackages := map[*pkg.Package]bool{}
	var namesFound []string
	found := map[string]bool{}
	for i, t := range tasks {
		if errs[i] != nil {
			return LoadResult{}, errs[i]
		}
		resp := results[i]
		if resp == nil {
			continue
		}
		pk := resp.Map("packages")
		if pk == nil {
			continue
		}
		versionsV, ok := pk.Get(t.realName)
		if !ok {
			continue
		}
		var versions []*phpjson.Map
		phpjson.Each(versionsV, func(_ string, v phpjson.Value) {
			if m := phpjson.ToMap(v); m != nil {
				versions = append(versions, m)
			}
		})
		if resp.Str("minified") == "composer/2.0" {
			versions = expandMinified(versions)
		}
		if adv, ok := resp.Get("security-advisories"); ok {
			r.mu.Lock()
			phpjson.Each(adv, func(_ string, a phpjson.Value) {
				if m := phpjson.ToMap(a); m != nil {
					r.Advisories[t.realName] = append(r.Advisories[t.realName], m)
				}
			})
			r.mu.Unlock()
		}
		if !found[t.realName] {
			found[t.realName] = true
			namesFound = append(namesFound, t.realName)
		}
		var toLoad []*phpjson.Map
		for _, v := range versions {
			v = v.Clone()
			if err := ensureNormalized(v); err != nil {
				return LoadResult{}, err
			}
			if alreadyLoaded[t.realName][v.Str("version_normalized")] {
				continue
			}
			var acc, fl map[string]int
			if withStability {
				acc, fl = acceptable, flags
			}
			ok, err := r.isVersionAcceptable(t.c, t.realName, v, acc, fl)
			if err != nil {
				return LoadResult{}, err
			}
			if ok {
				toLoad = append(toLoad, v)
			}
		}
		loaded, err := r.createPackages(toLoad, sources[i])
		if err != nil {
			return LoadResult{}, err
		}
		for _, p := range loaded {
			p.Repo = r
			if !inPackages[p] {
				inPackages[p] = true
				packages = append(packages, p)
			}
			if p.AliasOf != nil && !inPackages[p.AliasOf] {
				p.AliasOf.Repo = r
				inPackages[p.AliasOf] = true
				packages = append(packages, p.AliasOf)
			}
		}
	}
	return LoadResult{NamesFound: namesFound, Packages: packages}, nil
}

// expandMinified is MetadataMinifier::expand.
func expandMinified(versions []*phpjson.Map) []*phpjson.Map {
	var out []*phpjson.Map
	var expanded *phpjson.Map
	for _, v := range versions {
		if expanded == nil {
			expanded = v.Clone()
			out = append(out, expanded.Clone())
			continue
		}
		for _, k := range v.Keys() {
			val, _ := v.Get(k)
			if s, ok := val.(string); ok && s == "__unset" {
				expanded.Delete(k)
			} else {
				expanded.Set(k, phpjson.CloneValue(val))
			}
		}
		out = append(out, expanded.Clone())
	}
	return out
}

func (r *Composer) cachedAsyncDownload(fileName, realName string) (*phpjson.Map, string, error) {
	return r.cachedAsyncDownloadReq(fileName, realName, nil)
}

// prepareMetaReq dispatches the pre-file-download event of a package
// metadata request unless asyncFetchFile answers it without a request (nil).
func (r *Composer) prepareMetaReq(fileName string) (*metaReq, error) {
	name := strings.ToLower(fileName)
	u := strings.ReplaceAll(r.lazyProvidersURL, "%package%", name)
	cacheKey := "provider-" + strings.ReplaceAll(name, "/", "~") + ".json"
	lastModified := ""
	if c, ok := r.cache.Read(cacheKey); ok {
		if v, err := phpjson.Decode(c); err == nil {
			lastModified = phpjson.ToMap(v).Str("last-modified")
		}
	}
	r.mu.Lock()
	skip := r.notFound[u] || r.fresh[u] && lastModified != ""
	r.mu.Unlock()
	if skip {
		return nil, nil
	}
	q, err := r.preDownload(u)
	if err != nil {
		return nil, err
	}
	q.deferPost = true
	return q, nil
}

func (r *Composer) cachedAsyncDownloadReq(fileName, realName string, q *metaReq) (*phpjson.Map, string, error) {
	name := strings.ToLower(fileName)
	u := strings.ReplaceAll(r.lazyProvidersURL, "%package%", name)
	cacheKey := "provider-" + strings.ReplaceAll(name, "/", "~") + ".json"
	var contents *phpjson.Map
	lastModified := ""
	if c, ok := r.cache.Read(cacheKey); ok {
		if v, err := phpjson.Decode(c); err == nil {
			contents = phpjson.ToMap(v)
			lastModified = contents.Str("last-modified")
		}
	}
	resp, fresh, err := r.asyncFetchFileReq(u, cacheKey, lastModified, q)
	if err != nil {
		return nil, "", err
	}
	source := "downloaded file (" + httpdl.Sanitize(u) + ")"
	if fresh {
		source = "cached file (" + cacheKey + " originating from " + httpdl.Sanitize(u) + ")"
		resp = contents
	}
	if resp == nil {
		return nil, source, nil
	}
	pk := resp.Map("packages")
	if (pk == nil || !pk.Has(realName)) && !resp.Has("security-advisories") && !resp.Has("filter") {
		return nil, source, nil
	}
	return resp, source, nil
}

// asyncFetchFile returns (data, fresh=true when the cache is current, error).
func (r *Composer) asyncFetchFile(u, cacheKey, lastModified string) (*phpjson.Map, bool, error) {
	return r.asyncFetchFileReq(u, cacheKey, lastModified, nil)
}

// asyncFetchFileReq is asyncFetchFile; q carries the pre-file-download
// event result when the caller dispatched it already.
func (r *Composer) asyncFetchFileReq(u, cacheKey, lastModified string, q *metaReq) (*phpjson.Map, bool, error) {
	r.mu.Lock()
	if r.notFound[u] {
		r.mu.Unlock()
		return phpjson.NewMap(), false, nil
	}
	if r.fresh[u] && lastModified != "" {
		r.mu.Unlock()
		return nil, true, nil
	}
	r.mu.Unlock()
	if q == nil {
		var err error
		if q, err = r.preDownload(u); err != nil {
			return nil, false, err
		}
	}
	u = q.url
	headers := append([]string(nil), q.headers...)
	if lastModified != "" {
		headers = append(headers, "If-Modified-Since: "+lastModified)
	}
	resp, err := r.dl.Get(u, headers...)
	if err != nil {
		if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode == 404 {
			r.mu.Lock()
			r.notFound[u] = true
			r.mu.Unlock()
			return phpjson.NewMap(), false, nil
		}
		if strings.HasPrefix(err.Error(), "Your configuration does not allow") {
			return nil, false, err
		}
		r.mu.Lock()
		if !r.degradedMode {
			r.warn("<warning>" + httpdl.Sanitize(r.url) + " could not be fully loaded (" + err.Error() + "), package information was loaded from the local cache and may be out of date</warning>")
		}
		r.degradedMode = true
		r.mu.Unlock()
		if lastModified != "" {
			return nil, true, nil
		}
		return nil, false, err
	}
	if resp.StatusCode == 304 && len(resp.Body) == 0 {
		r.mu.Lock()
		r.fresh[u] = true
		r.mu.Unlock()
		return nil, true, nil
	}
	if err := r.postDownload(q, "", resp); err != nil {
		return nil, false, err
	}
	v, err := jsonfile.Parse(resp.Body, httpdl.Sanitize(u))
	if err != nil {
		return nil, false, err
	}
	r.responseWarnings(v)
	data := phpjson.ToMap(v)
	if data == nil {
		data = phpjson.NewMap()
	}
	body := resp.Body
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		data.Set("last-modified", lm)
		enc, _ := phpjson.Encode(data, phpjson.UnescapedSlashes|phpjson.UnescapedUnicode)
		body = []byte(enc)
	}
	if !r.cache.ReadOnly {
		r.cache.Write(cacheKey, body)
	}
	r.mu.Lock()
	r.fresh[u] = true
	r.mu.Unlock()
	return data, false, nil
}

// FindPackages is ComposerRepository::findPackages.
func (r *Composer) FindPackages(name string, c semver.Constraint) ([]*pkg.Package, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return nil, err
	}
	if err := r.checkSupported(); err != nil {
		return nil, err
	}
	name = strings.ToLower(name)
	if r.lazyProvidersURL != "" {
		hasPartial, err := r.partialPackages()
		if err != nil {
			return nil, err
		}
		if hasPartial {
			if _, ok := r.partialByName[name]; ok {
				cands, err := r.whatProvidesPartial(name, nil, nil, nil)
				if err != nil {
					return nil, err
				}
				return filterByConstraint(cands, c), nil
			}
		}
		if r.hasAvailableList && !r.lazyRepoContains(name) {
			return nil, nil
		}
		res, err := r.loadAsyncPackages([]NameConstraint{{name, c}}, nil, nil, nil, false)
		if err != nil {
			return nil, err
		}
		return res.Packages, nil
	}
	if r.hasProviders {
		names, err := r.providerNames()
		if err != nil {
			return nil, err
		}
		for _, n := range names {
			if n == name {
				cands, err := r.whatProvides(n, nil, nil, nil)
				if err != nil {
					return nil, err
				}
				return filterByConstraint(cands, c), nil
			}
		}
		return nil, nil
	}
	return r.Array.FindPackages(name, c)
}

func filterByConstraint(ps []*pkg.Package, c semver.Constraint) []*pkg.Package {
	if c == nil {
		return ps
	}
	var out []*pkg.Package
	for _, p := range ps {
		if c.Matches(semver.NewConstraint("==", p.Version())) {
			out = append(out, p)
		}
	}
	return out
}

// Packages is ComposerRepository::getPackages.
func (r *Composer) Packages() ([]*pkg.Package, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return nil, err
	}
	if err := r.checkSupported(); err != nil {
		return nil, err
	}
	if r.lazyProvidersURL != "" {
		if r.availablePackages != nil && len(r.availablePatterns) == 0 {
			var names []NameConstraint
			for n := range r.availablePackages {
				names = append(names, NameConstraint{Name: n, Constraint: &semver.MatchAll{}})
			}
			res, err := r.loadAsyncPackages(names, nil, nil, nil, false)
			if err != nil {
				return nil, err
			}
			return res.Packages, nil
		}
		hasPartial, err := r.partialPackages()
		if err != nil {
			return nil, err
		}
		if hasPartial {
			var all []*phpjson.Map
			for _, n := range r.partialOrder {
				all = append(all, r.partialByName[n]...)
			}
			return r.createPackages(all, "packages.json inline packages")
		}
		return nil, &UnsupportedError{"Composer repositories that have lazy providers and no available-packages list can not load the complete list of packages, use getPackageNames instead."}
	}
	if r.hasProviders {
		return nil, &UnsupportedError{"Composer repositories that have providers can not load the complete list of packages, use getPackageNames instead."}
	}
	return r.Array.Packages()
}

// Providers is ComposerRepository::getProviders.
func (r *Composer) Providers(name string) ([]ProviderInfo, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return nil, err
	}
	if r.providersAPIURL != "" {
		resp, err := r.dl.Get(strings.ReplaceAll(r.providersAPIURL, "%package%", name), r.headers()...)
		if err != nil {
			if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode == 404 {
				return nil, nil
			}
			return nil, err
		}
		v, err := phpjson.Decode(resp.Body)
		if err != nil {
			return nil, err
		}
		r.responseWarnings(v)
		var out []ProviderInfo
		idx := map[string]int{}
		phpjson.Each(phpjson.ToMap(v).Map("providers"), func(_ string, pv phpjson.Value) {})
		if provs, ok := phpjson.ToMap(v).Get("providers"); ok {
			phpjson.Each(provs, func(_ string, pv phpjson.Value) {
				m := phpjson.ToMap(pv)
				pi := ProviderInfo{Name: m.Str("name"), Description: m.Str("description"), Type: m.Str("type")}
				if i, ok := idx[pi.Name]; ok {
					out[i] = pi
					return
				}
				idx[pi.Name] = len(out)
				out = append(out, pi)
			})
		}
		return out, nil
	}
	var out []ProviderInfo
	hasPartial, err := r.partialPackages()
	if err != nil {
		return nil, err
	}
	if hasPartial {
		seen := map[string]bool{}
		for _, n := range r.partialOrder {
			for _, cand := range r.partialByName[n] {
				cn := cand.Str("name")
				if seen[cn] {
					continue
				}
				prov := cand.Map("provide")
				if prov == nil || !prov.Has(name) {
					continue
				}
				seen[cn] = true
				out = append(out, ProviderInfo{Name: cn, Description: cand.Str("description"), Type: cand.Str("type")})
			}
		}
	}
	if r.lazyProvidersURL == "" && !r.hasProviders {
		pkgs, err := r.Array.Packages()
		if err != nil {
			return nil, err
		}
		out = append(out, ProvidersInList(pkgs, name)...)
	}
	return out, nil
}

// Advisory is a (partial) security advisory.
type Advisory struct {
	PackageName      string
	AdvisoryID       string
	AffectedVersions semver.Constraint
	Title            string
	CVE              string
	Link             string
	Severity         string
	ReportedAt       string
}

// HasSecurityAdvisories is ComposerRepository::hasSecurityAdvisories.
func (r *Composer) HasSecurityAdvisories() bool {
	if _, err := r.loadRootServerFile(); err != nil {
		return false
	}
	return r.hasAdvCfg && (r.advMetadata || r.advAPIURL != "")
}

func newAdvisory(name string, data *phpjson.Map) Advisory {
	affected := data.Str("affectedVersions")
	c, err := semver.ParseConstraints(affected)
	if err != nil {
		re := regexp.MustCompile(`(^[>=<^~]*[\d.]+).*`)
		c, err = semver.ParseConstraints(re.ReplaceAllString(affected, "$1"))
		if err != nil {
			c = semver.NewConstraint("==", "0.0.0-invalid-version")
		}
	}
	return Advisory{PackageName: name, AdvisoryID: data.Str("advisoryId"), AffectedVersions: c, Title: data.Str("title"), CVE: data.Str("cve"), Link: data.Str("link"), Severity: data.Str("severity"), ReportedAt: data.Str("reportedAt")}
}

// SecurityAdvisories is ComposerRepository::getSecurityAdvisories.
func (r *Composer) SecurityAdvisories(names []NameConstraint, allowPartial bool) (map[string][]Advisory, error) {
	out := map[string][]Advisory{}
	if !r.HasSecurityAdvisories() {
		return out, nil
	}
	remaining := map[string]semver.Constraint{}
	var order []string
	for _, nc := range names {
		if r.hasAvailableList && !r.lazyRepoContains(strings.ToLower(nc.Name)) {
			continue
		}
		remaining[nc.Name] = nc.Constraint
		order = append(order, nc.Name)
	}
	if r.advMetadata && (allowPartial || r.advAPIURL == "") {
		for _, name := range order {
			lname := strings.ToLower(name)
			if IsPlatformPackage(lname) || lname == "__root__" {
				continue
			}
			resp, _, err := r.cachedAsyncDownload(lname, lname)
			if err != nil {
				return nil, err
			}
			if resp == nil {
				continue
			}
			adv, ok := resp.Get("security-advisories")
			if !ok || !phpjson.IsArray(adv) {
				continue
			}
			phpjson.Each(adv, func(_ string, a phpjson.Value) {
				m := phpjson.ToMap(a)
				if m == nil {
					return
				}
				ad := newAdvisory(lname, m)
				if ad.AffectedVersions.Matches(remaining[lname]) {
					out[lname] = append(out[lname], ad)
				}
			})
			delete(remaining, lname)
		}
	}
	if r.advAPIURL != "" && len(remaining) > 0 {
		return nil, &UnsupportedError{"gocomposer does not query the security-advisories API of " + r.RepoName() + " (only metadata advisories are supported)"}
	}
	return out, nil
}

func sha256Hex(b []byte) string {
	sum := sha256pkg.Sum256(b)
	return hex.EncodeToString(sum[:])
}
