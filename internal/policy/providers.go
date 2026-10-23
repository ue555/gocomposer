// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80: the
// FilterListProviderInterface / AdvisoryProviderInterface implementations of
// src/Composer/Repository/ComposerRepository.php (hasFilter, getFilter,
// getFilterLists, loadFilterSummary, parseUserFilterConfig,
// hasSecurityAdvisories, getSecurityAdvisories), PackageRepository.php and
// FilterRepository.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// ProviderResolver maps repositories to their filter-list/advisory provider
// implementations, standing in for PHP's instanceof checks. A nil resolver
// only recognizes repositories that implement the interfaces themselves.
type ProviderResolver struct {
	// HTTP performs the POST requests of filter api-url and the
	// security-advisories API. Without it those code paths fail.
	HTTP HTTPDownloader
	// IO receives repository warnings (may be nil).
	IO IO

	mu           sync.Mutex
	packageRepos map[repo.Repository]phpjson.Value
	composer     map[*repo.Composer]*composerRepoProvider
}

// RegisterPackageRepository records the raw config of a "package" type
// repository (its "filter" and "security-advisories" keys), since
// repo.NewPackageRepository does not keep it.
func (r *ProviderResolver) RegisterPackageRepository(rp repo.Repository, config phpjson.Value) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.packageRepos == nil {
		r.packageRepos = map[repo.Repository]phpjson.Value{}
	}
	r.packageRepos[rp] = config
}

func (r *ProviderResolver) packageRepoConfig(rp repo.Repository) (phpjson.Value, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.packageRepos[rp]
	return c, ok
}

func (r *ProviderResolver) composerProvider(c *repo.Composer) *composerRepoProvider {
	if r == nil {
		return &composerRepoProvider{repo: c}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.composer == nil {
		r.composer = map[*repo.Composer]*composerRepoProvider{}
	}
	p, ok := r.composer[c]
	if !ok {
		p = &composerRepoProvider{repo: c, http: r.HTTP, io: r.IO}
		r.composer[c] = p
	}
	return p
}

// FilterListProviderFor resolves "$repository instanceof FilterListProviderInterface".
func (r *ProviderResolver) FilterListProviderFor(rp repo.Repository) FilterListProvider {
	if p, ok := rp.(FilterListProvider); ok {
		return p
	}
	switch t := rp.(type) {
	case *repo.Composer:
		return r.composerProvider(t)
	case *repo.Filter:
		return &filterRepoProvider{filter: t, resolver: r}
	}
	if c, ok := r.packageRepoConfig(rp); ok {
		return &PackageRepositoryProvider{Config: c}
	}
	return nil
}

// AdvisoryProviderFor resolves "$repository instanceof AdvisoryProviderInterface".
func (r *ProviderResolver) AdvisoryProviderFor(rp repo.Repository) AdvisoryProvider {
	if p, ok := rp.(AdvisoryProvider); ok {
		return p
	}
	switch t := rp.(type) {
	case *repo.Composer:
		return r.composerProvider(t)
	case *repo.Filter:
		return &filterRepoProvider{filter: t, resolver: r}
	}
	if c, ok := r.packageRepoConfig(rp); ok {
		return &PackageRepositoryProvider{Config: c}
	}
	return nil
}

// ------------------------------------------------------- PackageRepository

// PackageRepositoryProvider is the filter/advisory part of PackageRepository,
// built from the repository config.
type PackageRepositoryProvider struct {
	Config phpjson.Value
}

func (p *PackageRepositoryProvider) filter() phpjson.Value { return coalesce(p.Config, "filter", nil) }
func (p *PackageRepositoryProvider) advisories() phpjson.Value {
	return coalesce(p.Config, "security-advisories", nil)
}

func (p *PackageRepositoryProvider) HasFilter() (bool, error) {
	return phpjson.ArrayLen(p.filter()) > 0, nil
}

func (p *PackageRepositoryProvider) GetFilter(pcm *ConstraintMap, _ []string) (*FilterByList, error) {
	out := NewOMap[[]*FilterListEntry]()
	var err error
	eachArr(p.filter(), func(listName string, entries phpjson.Value) {
		eachArr(entries, func(_ string, data phpjson.Value) {
			if err != nil {
				return
			}
			e, ee := CreateFilterListEntry(listName, data)
			if ee != nil {
				err = ee
				return
			}
			if pcm.Has(e.PackageName) {
				appendTo(out, listName, e)
			}
		})
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (p *PackageRepositoryProvider) GetFilterLists() ([]string, error) {
	var names []string
	eachArr(p.filter(), func(k string, _ phpjson.Value) { names = append(names, k) })
	return names, nil
}

func (p *PackageRepositoryProvider) HasSecurityAdvisories() (bool, error) {
	return phpjson.ArrayLen(p.advisories()) > 0, nil
}

func (p *PackageRepositoryProvider) GetSecurityAdvisories(pcm *ConstraintMap, _ bool) (*AdvisoryMap, error) {
	return AdvisoriesFromRawData(p.advisories(), pcm)
}

// -------------------------------------------------------- FilterRepository

type filterRepoProvider struct {
	filter   *repo.Filter
	resolver *ProviderResolver
}

func (f *filterRepoProvider) allowed(pcm *ConstraintMap) *ConstraintMap {
	out := NewOMap[semver.Constraint]()
	pcm.Each(func(name string, c semver.Constraint) {
		if f.filter.AllowsPackage(name) {
			out.Set(name, c)
		}
	})
	return out
}

func (f *filterRepoProvider) HasFilter() (bool, error) {
	inner := f.resolver.FilterListProviderFor(f.filter.Repo)
	if inner == nil {
		return false, nil
	}
	return inner.HasFilter()
}

func (f *filterRepoProvider) GetFilter(pcm *ConstraintMap, lists []string) (*FilterByList, error) {
	inner := f.resolver.FilterListProviderFor(f.filter.Repo)
	if inner == nil {
		return NewOMap[[]*FilterListEntry](), nil
	}
	return inner.GetFilter(f.allowed(pcm), lists)
}

func (f *filterRepoProvider) GetFilterLists() ([]string, error) {
	inner := f.resolver.FilterListProviderFor(f.filter.Repo)
	if inner == nil {
		return nil, nil
	}
	return inner.GetFilterLists()
}

func (f *filterRepoProvider) HasSecurityAdvisories() (bool, error) {
	inner := f.resolver.AdvisoryProviderFor(f.filter.Repo)
	if inner == nil {
		return false, nil
	}
	return inner.HasSecurityAdvisories()
}

func (f *filterRepoProvider) GetSecurityAdvisories(pcm *ConstraintMap, allowPartial bool) (*AdvisoryMap, error) {
	inner := f.resolver.AdvisoryProviderFor(f.filter.Repo)
	if inner == nil {
		return NewOMap[[]*Advisory](), nil
	}
	return inner.GetSecurityAdvisories(f.allowed(pcm), allowPartial)
}

// ------------------------------------------------------ ComposerRepository

// AdvisoryAPIBatchSize is ComposerRepository::ADVISORY_API_BATCH_SIZE.
const AdvisoryAPIBatchSize = 500

type composerRepoProvider struct {
	repo *repo.Composer
	http HTTPDownloader
	io   IO

	userFilterParsed bool
	userFilterOff    bool
	userFilterSkip   map[string]bool
	userFilterErr    error
}

// parseUserFilterConfig is ComposerRepository::parseUserFilterConfig.
func (c *composerRepoProvider) parseUserFilterConfig() error {
	if c.userFilterParsed {
		return c.userFilterErr
	}
	c.userFilterParsed = true
	raw, _ := c.repo.PolicyRepoConfig().Get("filter")
	switch {
	case raw == false:
		c.userFilterOff = true
	case raw == nil || raw == true:
	case !phpjson.IsArray(raw):
		c.userFilterErr = newErr(UnexpectedValueException, `Repository "filter" must be a boolean or an object mapping advertised list names to false.`)
	default:
		c.userFilterSkip = map[string]bool{}
		eachArr(raw, func(listName string, value phpjson.Value) {
			if c.userFilterErr != nil {
				return
			}
			if isIntKey(listName) || listName == "" {
				c.userFilterErr = newErr(UnexpectedValueException, `Repository "filter" keys must be non-empty list-name strings.`)
				return
			}
			if value == true {
				return
			}
			if value != false {
				c.userFilterErr = newErr(UnexpectedValueException, fmt.Sprintf(`Repository "filter" entry for "%s" must be a boolean; got %s.`, listName, varExport(value)))
				return
			}
			c.userFilterSkip[listName] = true
		})
	}
	return c.userFilterErr
}

// varExport approximates var_export() for scalars.
func varExport(v phpjson.Value) string {
	switch t := v.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + strings.ReplaceAll(strings.ReplaceAll(t, `\`, `\\`), `'`, `\'`) + "'"
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return phpjson.FormatFloatVarExport(t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	}
	s, _ := phpjson.Encode(v, 0)
	return s
}

// filterInfo is $this->filterConfig, set from the root file's "filter" key
// for metadata-url (v2) repositories.
func (c *composerRepoProvider) filterInfo() (*ComposerRepositoryFilterInformation, error) {
	data, err := c.repo.PolicyRootData()
	if err != nil {
		return nil, err
	}
	if data.Str("metadata-url") == "" {
		return nil, nil
	}
	f, ok := data.Get("filter")
	if !ok || !phpjson.IsArray(f) {
		return nil, nil
	}
	return ComposerRepositoryFilterInformationFromData(f, c.repo.CanonicalizeURL), nil
}

func (c *composerRepoProvider) HasFilter() (bool, error) {
	if err := c.parseUserFilterConfig(); err != nil {
		return false, err
	}
	if c.userFilterOff {
		return false, nil
	}
	info, err := c.filterInfo()
	if err != nil {
		return false, err
	}
	return info != nil && info.Metadata, nil
}

func (c *composerRepoProvider) GetFilterLists() ([]string, error) {
	if err := c.parseUserFilterConfig(); err != nil {
		return nil, err
	}
	if c.userFilterOff {
		return nil, nil
	}
	info, err := c.filterInfo()
	if err != nil || info == nil {
		return nil, err
	}
	var out []string
	for _, l := range info.Lists {
		if !c.userFilterSkip[l] {
			out = append(out, l)
		}
	}
	return out, nil
}

func (c *composerRepoProvider) restrictToAvailable(pcm *ConstraintMap) *ConstraintMap {
	if !c.repo.HasAvailablePackageList() {
		return pcm
	}
	out := NewOMap[semver.Constraint]()
	pcm.Each(func(name string, cs semver.Constraint) {
		if c.repo.LazyRepoContains(strings.ToLower(name)) {
			out.Set(name, cs)
		}
	})
	return out
}

func (c *composerRepoProvider) GetFilter(pcm *ConstraintMap, configuredLists []string) (*FilterByList, error) {
	info, err := c.filterInfo()
	if err != nil {
		return nil, err
	}
	pcm = c.restrictToAvailable(pcm)
	fresh := c.repo.HasFreshMetadata()

	if info != nil && info.APIURL != nil && !fresh {
		client := &FilterListApiClient{HTTP: c.http, Options: c.repo.PolicyOptions()}
		resp, err := client.PostPurls(*info.APIURL, pcm, configuredLists)
		if err != nil {
			return nil, err
		}
		decoded, err := DecodeJSONResponse(resp, *info.APIURL)
		if err != nil {
			return nil, err
		}
		if !isset(decoded, "filter") || !phpjson.IsArray(coalesce(decoded, "filter", nil)) {
			return nil, &httpdl.TransportError{Msg: "Filter api-url " + *info.APIURL + " returned an unexpected response for " + c.repo.RepoName()}
		}
		return BuildFilterListEntries(coalesce(decoded, "filter", nil), pcm, "")
	}

	if info != nil && info.SummaryURL != nil && !fresh {
		summary, err := c.loadFilterSummary(*info.SummaryURL)
		if err != nil {
			return nil, err
		}
		candidates := map[string]bool{}
		for _, listName := range configuredLists {
			entries, ok := summary.Get(listName)
			if !ok {
				continue
			}
			var perr error
			entries.Each(func(packageName, summaryConstraint string) {
				if perr != nil {
					return
				}
				pc, ok := pcm.Get(packageName)
				if !ok {
					return
				}
				if _, all := pc.(*semver.MatchAll); !all {
					sc, e := semver.ParseConstraints(summaryConstraint)
					if e != nil {
						perr = e
						return
					}
					if !sc.Matches(pc) {
						return
					}
				}
				candidates[packageName] = true
			})
			if perr != nil {
				return nil, perr
			}
		}
		kept := NewOMap[semver.Constraint]()
		pcm.Each(func(name string, cs semver.Constraint) {
			if candidates[name] {
				kept.Set(name, cs)
			}
		})
		pcm = kept
	}

	filter := NewOMap[[]*FilterListEntry]()
	for _, origName := range pcm.Keys() {
		name := strings.ToLower(origName)
		if isPlatformPackage(name) || name == "__root__" {
			continue
		}
		resp, err := c.repo.PolicyPackageMetadata(name)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			continue
		}
		f, ok := resp.Get("filter")
		if !ok || f == nil || !phpjson.IsArray(f) {
			continue
		}
		built, err := BuildFilterListEntries(f, pcm, name)
		if err != nil {
			return nil, err
		}
		built.Each(func(listName string, entries []*FilterListEntry) {
			for _, e := range entries {
				appendTo(filter, listName, e)
			}
		})
	}
	return filter, nil
}

// loadFilterSummary is ComposerRepository::loadFilterSummary.
func (c *composerRepoProvider) loadFilterSummary(summaryURL string) (*OMap[*OMap[string]], error) {
	data, err := c.repo.PolicyFetchFilterSummary(summaryURL)
	if err != nil {
		return nil, err
	}
	if !phpjson.IsArray(data) || !isset(data, "filter") || !phpjson.IsArray(coalesce(data, "filter", nil)) {
		return nil, &httpdl.TransportError{Msg: "Filter summary URL " + summaryURL + " returned 404 for " + c.repo.RepoName(), StatusCode: 404}
	}
	summary := NewOMap[*OMap[string]]()
	eachArr(coalesce(data, "filter", nil), func(listName string, packages phpjson.Value) {
		if err != nil {
			return
		}
		if isIntKey(listName) || !phpjson.IsArray(packages) {
			shown := listName
			if isIntKey(listName) {
				shown = ""
			}
			err = newErr(UnexpectedValueException, "Invalid filter summary received from "+c.repo.RepoName()+`: list "`+shown+`" must map to an object of package => constraint`)
			return
		}
		eachArr(packages, func(packageName string, cs phpjson.Value) {
			if err != nil {
				return
			}
			s, ok := cs.(string)
			if isIntKey(packageName) || !ok {
				err = newErr(UnexpectedValueException, "Invalid filter summary received from "+c.repo.RepoName()+`: list "`+listName+`" entries must be strings`)
				return
			}
			m, ok := summary.Get(listName)
			if !ok {
				m = NewOMap[string]()
				summary.Set(listName, m)
			}
			m.Set(strings.ToLower(packageName), s)
		})
	})
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (c *composerRepoProvider) advisoryConfig() (configured, metadata bool, apiURL string, err error) {
	configured, metadata, apiURL, err = c.repo.PolicySecurityAdvisoryConfig()
	if err == nil && configured && apiURL == "" && !c.repo.HasAvailablePackageList() {
		err = newErr(UnexpectedValueException, "Invalid security advisory configuration on "+c.repo.RepoName()+": If the repository does not provide a security-advisories.api-url then available-packages or available-package-patterns are required to be provided for performance reason.")
	}
	return
}

func (c *composerRepoProvider) HasSecurityAdvisories() (bool, error) {
	configured, metadata, apiURL, err := c.advisoryConfig()
	if err != nil {
		return false, err
	}
	return configured && (metadata || apiURL != ""), nil
}

func (c *composerRepoProvider) GetSecurityAdvisories(pcm *ConstraintMap, allowPartial bool) (*AdvisoryMap, error) {
	configured, metadata, apiURL, err := c.advisoryConfig()
	if err != nil {
		return nil, err
	}
	advisories := NewOMap[[]*Advisory]()
	if !configured {
		return advisories, nil
	}
	pcm = c.restrictToAvailable(pcm).Clone()

	create := func(name string, list phpjson.Value) ([]*Advisory, error) {
		var out []*Advisory
		var cerr error
		eachArr(list, func(_ string, data phpjson.Value) {
			if cerr != nil {
				return
			}
			a, e := CreateAdvisory(name, data)
			if e != nil {
				cerr = e
				return
			}
			if pc, ok := pcm.Get(name); ok && a.AffectedVersions.Matches(pc) {
				out = append(out, a)
			}
		})
		return out, cerr
	}

	if metadata && (allowPartial || apiURL == "") {
		for _, origName := range pcm.Keys() {
			name := strings.ToLower(origName)
			if isPlatformPackage(name) || name == "__root__" {
				continue
			}
			resp, err := c.repo.PolicyPackageMetadata(name)
			if err != nil {
				return nil, err
			}
			if resp == nil {
				continue
			}
			sa, ok := resp.Get("security-advisories")
			if !ok || sa == nil || !phpjson.IsArray(sa) {
				continue
			}
			if phpjson.ArrayLen(sa) > 0 {
				list, err := create(name, sa)
				if err != nil {
					return nil, err
				}
				advisories.Set(name, list)
			}
			pcm.Delete(name)
		}
	}

	if apiURL != "" && pcm.Len() > 0 {
		if c.http == nil {
			return nil, newErr(LogicException, "No HTTP downloader configured to query the security-advisories API of "+c.repo.RepoName())
		}
		names := pcm.Keys()
		warned := false
		for start := 0; start < len(names); start += AdvisoryAPIBatchSize {
			end := start + AdvisoryAPIBatchSize
			if end > len(names) {
				end = len(names)
			}
			var form []string
			for i, n := range names[start:end] {
				form = append(form, url.QueryEscape("packages["+strconv.Itoa(i)+"]")+"="+url.QueryEscape(n))
			}
			resp, err := c.http.Request(apiURL, postOptions(c.repo.PolicyOptions(), "application/x-www-form-urlencoded", strings.Join(form, "&")))
			if err != nil {
				return nil, err
			}
			data, err := DecodeJSONResponse(resp, apiURL)
			if err != nil {
				return nil, err
			}
			var aerr error
			eachArr(coalesce(data, "advisories", nil), func(name string, list phpjson.Value) {
				if aerr != nil {
					return
				}
				if !pcm.Has(name) {
					if !warned && c.io != nil {
						requested := pcm.Keys()
						shown := strings.Join(requested, ", ")
						if len(requested) > 20 {
							shown = strings.Join(requested[:20], ", ") + " and " + strconv.Itoa(len(requested)-20) + " more"
						}
						c.io.WriteError("<warning>" + c.repo.RepoName() + " returned names which were not requested in response to the security-advisories API. " + name + " was not requested but is present in the response. Requested names were: " + shown + "</warning>")
					}
					warned = true
					return
				}
				if phpjson.ArrayLen(list) > 0 {
					l, e := create(name, list)
					if e != nil {
						aerr = e
						return
					}
					advisories.Set(name, l)
				}
			})
			if aerr != nil {
				return nil, aerr
			}
		}
	}

	out := NewOMap[[]*Advisory]()
	advisories.Each(func(name string, list []*Advisory) {
		if len(list) > 0 {
			out.Set(name, list)
		}
	})
	return out, nil
}
