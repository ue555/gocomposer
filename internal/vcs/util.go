package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/Url.php, src/Composer/Util/Filesystem.php,
// src/Composer/IO/BaseIO.php authentication store).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"regexp"
	"strings"
	"sync"

	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
)

var (
	reLocalPath      = regexp.MustCompile(`(?i)^(file://|/|/?[a-z]:[\\/]|\.\.[\\/]|[a-z0-9_.-]+[\\/])`)
	reFileScheme     = regexp.MustCompile(`(?i)^file://`)
	reStripCreds     = regexp.MustCompile(`(?i)([a-z0-9][a-z0-9+.-]*://)([^:/\s?#]*)(?::([^\s/?#]*))?@`)
	reHTTPPrefixOnly = regexp.MustCompile(`(?i)^https?://$`)
	reCacheKeyChars  = regexp.MustCompile(`(?i)[^a-z0-9.]`)
	reUnusableCache  = regexp.MustCompile(`(^|[\\/])(\$null|nul|NUL|/dev/null)([\\/]|$)`)
)

// IsLocalPath is Filesystem::isLocalPath (non-Windows).
func IsLocalPath(p string) bool { return reLocalPath.MatchString(p) }

// PlatformPath is Filesystem::getPlatformPath (non-Windows).
func PlatformPath(p string) string { return reFileScheme.ReplaceAllString(p, "") }

// Sanitize is Url::sanitize.
func Sanitize(u string) string { return httpdl.Sanitize(u) }

// StripCredentials is Url::stripCredentials.
func StripCredentials(u string) string {
	return reStripCreds.ReplaceAllStringFunc(u, func(m string) string {
		sm := reStripCreds.FindStringSubmatch(m)
		if sm[3] == "" && !reHTTPPrefixOnly.MatchString(sm[1]) {
			return sm[1] + sm[2] + "@"
		}
		return sm[1]
	})
}

// CacheKey is Preg::replace('{[^a-z0-9.]}i', '-', Url::sanitize($url)),
// the directory name used below cache-vcs-dir and cache-repo-dir.
func CacheKey(u string) string { return reCacheKeyChars.ReplaceAllString(Sanitize(u), "-") }

// CacheUsable is Cache::isUsable.
func CacheUsable(dir string) bool { return !reUnusableCache.MatchString(dir) }

// Auth is one IO authentication entry.
type Auth struct {
	Username string
	Password string
	// HasPassword distinguishes a null password from "".
	HasPassword bool
}

// AuthStore is the IOInterface authentication store (BaseIO::$authentications)
// loaded from the config like BaseIO::loadConfiguration.
type AuthStore struct {
	mu    sync.Mutex
	auths map[string]Auth
	// http is the downloader whose credentials this store shares once
	// linked (Composer has a single IO authentication store).
	http *httpdl.Downloader
}

// Link makes the store share the downloader's credentials: entries only
// known here are copied over and every later read and write goes to the
// downloader, so credentials entered for a git command are used by API
// requests and the other way round.
func (s *AuthStore) Link(d *httpdl.Downloader) *AuthStore {
	if d == nil {
		return s
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.http != nil {
		return s
	}
	for domain, a := range s.auths {
		if !d.HasAuth(domain) {
			d.SetAuth(domain, a.Username, a.Password)
		}
	}
	s.http = d
	return s
}

// NewAuthStore loads the credentials configured in cfg.
func NewAuthStore(cfg *config.Config) *AuthStore {
	s := &AuthStore{auths: map[string]Auth{}}
	if cfg == nil {
		return s
	}
	set := func(domain, user, pass string) { s.auths[domain] = Auth{user, pass, true} }
	if m := phpjson.ToMap(cfg.Get("bitbucket-oauth")); m != nil {
		for _, d := range m.Keys() {
			if c := m.Map(d); c != nil {
				set(d, c.Str("consumer-key"), c.Str("consumer-secret"))
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("github-oauth")); m != nil {
		for _, d := range m.Keys() {
			set(d, m.Str(d), "x-oauth-basic")
		}
	}
	if m := phpjson.ToMap(cfg.Get("gitlab-oauth")); m != nil {
		for _, d := range m.Keys() {
			v, _ := m.Get(d)
			if mm := phpjson.ToMap(v); mm != nil {
				set(d, mm.Str("token"), "oauth2")
			} else {
				s, _ := v.(string)
				set(d, s, "oauth2")
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("gitlab-token")); m != nil {
		for _, d := range m.Keys() {
			v, _ := m.Get(d)
			if mm := phpjson.ToMap(v); mm != nil {
				set(d, mm.Str("username"), mm.Str("token"))
			} else {
				s, _ := v.(string)
				set(d, s, "private-token")
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("forgejo-token")); m != nil {
		for _, d := range m.Keys() {
			if c := m.Map(d); c != nil {
				set(d, c.Str("username"), c.Str("token"))
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("http-basic")); m != nil {
		for _, d := range m.Keys() {
			if c := m.Map(d); c != nil {
				set(d, c.Str("username"), c.Str("password"))
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("bearer")); m != nil {
		for _, d := range m.Keys() {
			set(d, m.Str(d), "bearer")
		}
	}
	return s
}

// Has is IOInterface::hasAuthentication.
func (s *AuthStore) Has(domain string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.http != nil {
		return s.http.HasAuth(domain)
	}
	_, ok := s.auths[domain]
	return ok
}

// Get is IOInterface::getAuthentication.
func (s *AuthStore) Get(domain string) Auth {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.http != nil {
		a, ok := s.http.GetAuth(domain)
		if !ok {
			return Auth{}
		}
		return Auth{a.Username, a.Password, true}
	}
	return s.auths[domain]
}

// Set is IOInterface::setAuthentication.
func (s *AuthStore) Set(domain, user, pass string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.http != nil {
		s.http.SetAuth(domain, user, pass)
		return
	}
	s.auths[domain] = Auth{user, pass, true}
}

// SetPrompted is setAuthentication with IO answers (nil: PHP null).
func (s *AuthStore) SetPrompted(domain string, user, pass any) {
	s.mu.Lock()
	h := s.http
	s.mu.Unlock()
	if h != nil {
		h.SetPromptedAuth(domain, user, pass)
		return
	}
	u, _ := user.(string)
	p, _ := pass.(string)
	s.Set(domain, u, p)
}

// configList returns a config value as a list of strings.
func configList(cfg *config.Config, key string) []string {
	var out []string
	phpjson.Each(cfg.Get(key), func(_ string, v phpjson.Value) {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	})
	return out
}

func domainsRegex(cfg *config.Config, key string) string {
	ds := configList(cfg, key)
	q := make([]string, len(ds))
	for i, d := range ds {
		q[i] = regexp.QuoteMeta(d)
	}
	return "(" + strings.Join(q, "|") + ")"
}

// GitHubDomainsRegex is Git::getGitHubDomainsRegex.
func GitHubDomainsRegex(cfg *config.Config) string { return domainsRegex(cfg, "github-domains") }

// GitLabDomainsRegex is Git::getGitLabDomainsRegex.
func GitLabDomainsRegex(cfg *config.Config) string { return domainsRegex(cfg, "gitlab-domains") }
