// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Repository/ComposerRepository.php (whatProvides,
// loadProviderListings, getProviderNames, loadPackageList and the
// ResponseWarnings::output calls of the repository requests).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package repo

import (
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// RepositorySecurityError is Composer\Repository\RepositorySecurityException.
type RepositorySecurityError struct{ Msg string }

func (e *RepositorySecurityError) Error() string { return e.Msg }

// providerListing is ComposerRepository::$providerListing (null until the
// listings are loaded), keeping PHP's array order.
type providerListing struct {
	names []string
	hash  map[string]string
}

func (l *providerListing) set(name, sha256 string) {
	if _, ok := l.hash[name]; !ok {
		l.names = append(l.names, name)
	}
	l.hash[name] = sha256
}

// responseWarnings is ResponseWarnings::output($this->io, $this->url, $data).
func (r *Composer) responseWarnings(data phpjson.Value) bool {
	if r.dl == nil {
		return false
	}
	return r.dl.OutputResponseWarnings(r.url, data)
}

// ResponseWarnings shows the warnings/infos of a response the policy layer
// fetched for this repository (filter api-url, security-advisories API).
func (r *Composer) ResponseWarnings(data phpjson.Value) bool { return r.responseWarnings(data) }

// loadProviderListings is ComposerRepository::loadProviderListings.
func (r *Composer) loadProviderListings(data *phpjson.Map) error {
	if data == nil {
		return nil
	}
	if p, ok := data.Get("providers"); ok && p != nil {
		if r.providers == nil {
			r.providers = &providerListing{hash: map[string]string{}}
		}
		phpjson.Each(p, func(name string, meta phpjson.Value) {
			sha := ""
			if m := phpjson.ToMap(meta); m != nil {
				sha = m.Str("sha256")
			}
			r.providers.set(name, sha)
		})
	}
	if inc, ok := data.Get("provider-includes"); ok && inc != nil && r.providersURL != "" {
		var err error
		phpjson.Each(inc, func(include string, meta phpjson.Value) {
			if err != nil {
				return
			}
			sha := ""
			if m := phpjson.ToMap(meta); m != nil {
				sha = m.Str("sha256")
			}
			u := r.baseURL + "/" + strings.ReplaceAll(include, "%hash%", sha)
			cacheKey := strings.NewReplacer("%hash%", "", "$", "").Replace(include)
			var included *phpjson.Map
			if r.cache.Sha256(cacheKey) == sha {
				c, _ := r.cache.Read(cacheKey)
				v, derr := phpjson.Decode(c)
				if derr == nil {
					included = phpjson.ToMap(v)
				}
			}
			if included == nil {
				included, err = r.fetchFile(u, cacheKey, sha, false)
				if err != nil {
					return
				}
			}
			err = r.loadProviderListings(included)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ensureProviderListing loads the provider listings from the root file when
// they are not loaded yet (the root data is consumed by the partial packages
// once those were initialized, as in PHP).
func (r *Composer) ensureProviderListing() error {
	if r.providers != nil {
		return nil
	}
	data, err := r.loadRootServerFile()
	if err != nil {
		return err
	}
	if r.partialByName != nil {
		return nil
	}
	return r.loadProviderListings(data)
}

// providerNames is ComposerRepository::getProviderNames.
func (r *Composer) providerNames() ([]string, error) {
	if err := r.ensureProviderListing(); err != nil {
		return nil, err
	}
	if r.lazyProvidersURL != "" {
		return nil, nil
	}
	if r.providersURL != "" && r.providers != nil {
		return append([]string(nil), r.providers.names...), nil
	}
	return nil, nil
}

// whatProvides is ComposerRepository::whatProvides.
func (r *Composer) whatProvides(name string, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) ([]*pkg.Package, error) {
	hasPartial, err := r.partialPackages()
	if err != nil {
		return nil, err
	}
	if hasPartial {
		if _, ok := r.partialByName[name]; ok {
			return r.whatProvidesPartial(name, acceptable, flags, alreadyLoaded)
		}
	}
	if IsPlatformPackage(name) || name == "__root__" {
		return nil, nil
	}
	if err := r.ensureProviderListing(); err != nil {
		return nil, err
	}
	var u, cacheKey, hash string
	useLastModifiedCheck := false
	inListing := false
	if r.providers != nil {
		_, inListing = r.providers.hash[name]
	}
	switch {
	case r.lazyProvidersURL != "" && !inListing:
		u = strings.ReplaceAll(r.lazyProvidersURL, "%package%", name)
		cacheKey = "provider-" + strings.ReplaceAll(name, "/", "$") + ".json"
		useLastModifiedCheck = true
	case r.providersURL != "":
		if !inListing {
			return nil, nil
		}
		hash = r.providers.hash[name]
		u = strings.NewReplacer("%package%", name, "%hash%", hash).Replace(r.providersURL)
		cacheKey = "provider-" + strings.ReplaceAll(name, "/", "$") + ".json"
	default:
		return nil, nil
	}
	var packages *phpjson.Map
	source := ""
	cachedSource := "cached file (" + cacheKey + " originating from " + httpdl.Sanitize(u) + ")"
	if !useLastModifiedCheck && hash != "" && r.cache.Sha256(cacheKey) == hash {
		c, _ := r.cache.Read(cacheKey)
		if v, derr := phpjson.Decode(c); derr == nil {
			packages = phpjson.ToMap(v)
		}
		source = cachedSource
	} else if useLastModifiedCheck {
		if c, ok := r.cache.Read(cacheKey); ok && len(c) > 0 {
			v, _ := phpjson.Decode(c)
			contents := phpjson.ToMap(v)
			if _, loaded := alreadyLoaded[name]; loaded {
				packages = contents
				source = cachedSource
			} else if contents != nil && contents.Isset("last-modified") {
				resp, err := r.fetchFileIfLastModified(u, cacheKey, contents.Str("last-modified"))
				if err != nil {
					return nil, err
				}
				if resp == nil {
					packages = contents
					source = cachedSource
				} else {
					packages = resp
					source = "downloaded file (" + httpdl.Sanitize(u) + ")"
				}
			}
		}
	}
	if packages == nil || packages.Len() == 0 {
		packages, err = r.fetchFile(u, cacheKey, hash, useLastModifiedCheck)
		source = "downloaded file (" + httpdl.Sanitize(u) + ")"
		if err != nil {
			te, ok := err.(*httpdl.TransportError)
			if !ok || r.lazyProvidersURL == "" || te.StatusCode != 404 {
				return nil, err
			}
			packages = phpjson.NewMap()
			source = "not-found file (" + httpdl.Sanitize(u) + ")"
		}
	}

	var toLoad []*phpjson.Map
	uids := map[string]bool{}
	pk, _ := packages.Get("packages")
	var lerr error
	phpjson.Each(pk, func(_ string, versions phpjson.Value) {
		phpjson.Each(versions, func(_ string, vv phpjson.Value) {
			if lerr != nil {
				return
			}
			v := phpjson.ToMap(vv)
			if v == nil {
				return
			}
			normalizedName := strings.ToLower(v.Str("name"))
			if normalizedName != name {
				return
			}
			if hasPartial {
				if _, ok := r.partialByName[normalizedName]; ok {
					return
				}
			}
			uid := ""
			hasUID := false
			if raw, ok := v.Get("uid"); ok && raw != nil {
				uid, _ = phpjson.ScalarString(raw)
				hasUID = true
			}
			if hasUID && uids[uid] {
				return
			}
			v = v.Clone()
			if err := ensureNormalized(v); err != nil {
				lerr = err
				return
			}
			if alreadyLoaded[name][v.Str("version_normalized")] {
				return
			}
			ok, err := r.isVersionAcceptable(nil, normalizedName, v, acceptable, flags)
			if err != nil {
				lerr = err
				return
			}
			if ok {
				if hasUID {
					uids[uid] = true
				}
				toLoad = append(toLoad, v)
			}
		})
	})
	if lerr != nil {
		return nil, lerr
	}
	loaded, err := r.createPackages(toLoad, source)
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

// loadPackageList is ComposerRepository::loadPackageList.
func (r *Composer) loadPackageList(filter string) ([]string, error) {
	u := r.listURL
	if filter != "" {
		u += "?filter=" + url.QueryEscape(filter)
		return r.fetchPackageNames(u)
	}
	const cacheKey = "package-list.txt"
	if st, err := os.Stat(r.cache.Path(cacheKey)); err == nil && r.cache.IsEnabled() && time.Since(st.ModTime()) < 600*time.Second {
		if c, ok := r.cache.Read(cacheKey); ok && len(c) > 0 {
			return strings.Split(string(c), "\n"), nil
		}
	}
	names, err := r.fetchPackageNames(u)
	if err != nil {
		return nil, err
	}
	if !r.cache.ReadOnly {
		r.cache.Write(cacheKey, []byte(strings.Join(names, "\n")))
	}
	return names, nil
}

// fetchPackageNames GETs a list API URL and returns its packageNames.
func (r *Composer) fetchPackageNames(u string) ([]string, error) {
	resp, err := r.dl.Get(u, r.headers()...)
	if err != nil {
		return nil, err
	}
	v, err := decodeResponseJSON(resp, u)
	if err != nil {
		return nil, err
	}
	r.responseWarnings(v)
	var names []string
	if m := phpjson.AsMap(v); m != nil {
		pn, _ := m.Get("packageNames")
		phpjson.Each(pn, func(_ string, n phpjson.Value) {
			s, _ := phpjson.ScalarString(n)
			names = append(names, s)
		})
	}
	return names, nil
}

// decodeResponseJSON is Response::decodeJson.
func decodeResponseJSON(resp *httpdl.Response, u string) (phpjson.Value, error) {
	return jsonfile.Parse(resp.Body, u)
}

var reVendorQuery = regexp.MustCompile(`(?i)^\^(?P<query>(?P<vendor>[a-z0-9_.-]+)/[a-z0-9_.-]*)\*?$`)

// searchListAPI is the "^vendor/name" optimization of search() for
// repositories with providers and a list API; ok is false when it does not
// apply.
func (r *Composer) searchListAPI(query string) ([]SearchEntry, bool, error) {
	if r.listURL == "" {
		return nil, false, nil
	}
	m := reVendorQuery.FindStringSubmatch(query)
	if m == nil {
		return nil, false, nil
	}
	u := r.listURL + "?vendor=" + url.QueryEscape(m[2]) + "&filter=" + url.QueryEscape(m[1]+"*")
	names, err := r.fetchPackageNames(u)
	if err != nil {
		return nil, true, err
	}
	var out []SearchEntry
	for _, n := range names {
		out = append(out, SearchEntry{Name: n, Description: strp("")})
	}
	return out, true, nil
}
