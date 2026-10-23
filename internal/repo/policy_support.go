// Accessors used by internal/policy to port the FilterListProviderInterface
// and AdvisoryProviderInterface parts of ComposerRepository and
// FilterRepository (Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Repository/ComposerRepository.php, FilterRepository.php).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package repo

import (
	"gocomposer/internal/phpjson"
)

// RequiresPublishedTime is ComposerRepository::requiresPublishedTime: false
// only when the repository config sets "require-published-time": false.
func (r *Composer) RequiresPublishedTime() bool {
	v, ok := r.repoConfig.Get("require-published-time")
	return !(ok && v == false)
}

// PolicyRepoConfig returns the repository config (getRepoConfig).
func (r *Composer) PolicyRepoConfig() *phpjson.Map { return r.repoConfig }

// PolicyOptions returns the repository transport options ($this->options).
func (r *Composer) PolicyOptions() *phpjson.Map { return r.options }

// PolicyRootData loads packages.json (loadRootServerFile) and returns it.
func (r *Composer) PolicyRootData() (*phpjson.Map, error) { return r.loadRootServerFile() }

// CanonicalizeURL is ComposerRepository::canonicalizeUrl.
func (r *Composer) CanonicalizeURL(u string) string { return r.canonicalizeURL(u) }

// HasAvailablePackageList reports available-packages/-patterns presence.
func (r *Composer) HasAvailablePackageList() bool { return r.hasAvailableList }

// LazyRepoContains is lazyProvidersRepoContains for a lowercase name.
func (r *Composer) LazyRepoContains(name string) bool { return r.lazyRepoContains(name) }

// HasFreshMetadata reports whether package metadata was fetched in this run
// ($this->freshMetadataUrls !== []).
func (r *Composer) HasFreshMetadata() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.fresh) > 0
}

// PolicyPackageMetadata is startCachedAsyncDownload($name, $name): the
// decoded per-package metadata file, or nil when it has none for the name.
func (r *Composer) PolicyPackageMetadata(name string) (*phpjson.Map, error) {
	resp, _, err := r.cachedAsyncDownload(name, name)
	return resp, err
}

// PolicyFetchFilterSummary fetches a filter summary URL through the
// repository cache ("filter-summary.json") with If-Modified-Since, returning
// the cached contents when the server reports it fresh (loadFilterSummary).
func (r *Composer) PolicyFetchFilterSummary(u string) (phpjson.Value, error) {
	const cacheKey = "filter-summary.json"
	var contents phpjson.Value
	lastModified := ""
	if c, ok := r.cache.Read(cacheKey); ok {
		if v, err := phpjson.Decode(c); err == nil && phpjson.IsArray(v) {
			contents = v
			lastModified = phpjson.ToMap(v).Str("last-modified")
		}
	}
	resp, fresh, err := r.asyncFetchFile(u, cacheKey, lastModified)
	if err != nil {
		return nil, err
	}
	if fresh {
		return contents, nil
	}
	if resp == nil {
		return nil, nil
	}
	return resp, nil
}

// PolicySecurityAdvisoryConfig returns the parsed "security-advisories"
// root config: configured, metadata flag and canonical api-url ("" = null).
func (r *Composer) PolicySecurityAdvisoryConfig() (configured, metadata bool, apiURL string, err error) {
	if _, err = r.loadRootServerFile(); err != nil {
		return false, false, "", err
	}
	return r.hasAdvCfg, r.advMetadata, r.advAPIURL, nil
}

// AllowsPackage is FilterRepository::isAllowed.
func (f *Filter) AllowsPackage(name string) bool { return f.isAllowed(name) }
