// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/FilterList/FilterListEntry.php, FilterListEntryBuilder.php,
// ComposerRepositoryFilterInformation.php, FilterListApiClient.php,
// FilterListAuditor.php, FilterListProvider/FilterListProviderSet.php,
// FilterListProvider/UrlSourceFilterListProvider.php and
// src/Composer/Repository/FilterListProviderInterface.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"errors"
	"strings"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// FilterListEntry is Composer\FilterList\FilterListEntry.
type FilterListEntry struct {
	PackageName string
	ListName    string
	Constraint  semver.Constraint
	URL         *string
	Reason      *string
	ID          *string
	Source      *string
}

// NewFilterListEntry mirrors the constructor argument order
// (packageName, constraint, listName, url, reason, id, source).
func NewFilterListEntry(packageName string, constraint semver.Constraint, listName string, url, reason, id, source *string) *FilterListEntry {
	return &FilterListEntry{PackageName: packageName, ListName: listName, Constraint: constraint, URL: url, Reason: reason, ID: id, Source: source}
}

// CreateFilterListEntry is FilterListEntry::create.
func CreateFilterListEntry(listName string, data phpjson.Value) (*FilterListEntry, error) {
	c, err := semver.ParseConstraints(toStr(coalesce(data, "constraint", nil)))
	if err != nil {
		return nil, err
	}
	return NewFilterListEntry(toStr(coalesce(data, "package", nil)), c, listName,
		nullableStr(coalesce(data, "url", nil)), nullableStr(coalesce(data, "reason", nil)),
		nullableStr(coalesce(data, "id", nil)), nullableStr(coalesce(data, "source", nil))), nil
}

// jsonValue is the (array) cast used by the audit JSON output, with the
// constraint replaced by its pretty string.
func (e *FilterListEntry) jsonValue() *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("packageName", e.PackageName)
	m.Set("listName", e.ListName)
	m.Set("constraint", e.Constraint.PrettyString())
	set := func(k string, v *string) {
		if v == nil {
			m.Set(k, nil)
		} else {
			m.Set(k, *v)
		}
	}
	set("url", e.URL)
	set("reason", e.Reason)
	set("id", e.ID)
	set("source", e.Source)
	return m
}

// FilterByList is array<string(list name), list<FilterListEntry>>.
type FilterByList = OMap[[]*FilterListEntry]

// ConstraintMap is array<string(package name), ConstraintInterface>.
type ConstraintMap = OMap[semver.Constraint]

// BuildFilterListEntries is FilterListEntryBuilder::build; defaultPackage ""
// means null.
func BuildFilterListEntries(rawByList phpjson.Value, packageConstraintMap *ConstraintMap, defaultPackage string) (*FilterByList, error) {
	result := NewOMap[[]*FilterListEntry]()
	var err error
	eachArr(rawByList, func(listName string, entries phpjson.Value) {
		if err != nil || isIntKey(listName) || !phpjson.IsArray(entries) {
			return
		}
		eachArr(entries, func(_ string, data phpjson.Value) {
			if err != nil || !phpjson.IsArray(data) || !isset(data, "constraint") {
				return
			}
			if !isset(data, "package") {
				if defaultPackage == "" {
					return
				}
				m := phpjson.ToMap(data).Clone()
				m.Set("package", defaultPackage)
				data = m
			}
			entry, e := CreateFilterListEntry(listName, data)
			if e != nil {
				err = e
				return
			}
			pc, ok := packageConstraintMap.Get(entry.PackageName)
			if !ok || !entry.Constraint.Matches(pc) {
				return
			}
			appendTo(result, listName, entry)
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ComposerRepositoryFilterInformation is the "filter" key of a Composer
// repository's packages.json.
type ComposerRepositoryFilterInformation struct {
	Metadata   bool
	Lists      []string
	SummaryURL *string
	APIURL     *string
}

// ComposerRepositoryFilterInformationFromData is ComposerRepositoryFilterInformation::fromData.
func ComposerRepositoryFilterInformationFromData(data phpjson.Value, canonicalizeURL func(string) string) *ComposerRepositoryFilterInformation {
	var lists []string
	if isset(data, "lists") && phpjson.IsArray(coalesce(data, "lists", nil)) {
		eachArr(coalesce(data, "lists", nil), func(name string, config phpjson.Value) {
			if isIntKey(name) || !phpjson.IsArray(config) || !toBool(coalesce(config, "enabled", false)) {
				return
			}
			lists = append(lists, name)
		})
	}
	var kept []string
	for _, name := range lists {
		if inStrings(name, ReservedNames) || inStrings(name, FutureReservedNames) {
			continue
		}
		reserved := false
		for _, prefix := range FutureReservedPrefixes {
			if strings.HasPrefix(name, prefix) {
				reserved = true
			}
		}
		if !reserved {
			kept = append(kept, name)
		}
	}
	urlOf := func(key string) *string {
		s, ok := coalesce(data, key, nil).(string)
		if !ok || s == "" {
			return nil
		}
		if canonicalizeURL != nil {
			s = canonicalizeURL(s)
		}
		return &s
	}
	return &ComposerRepositoryFilterInformation{Metadata: toBool(coalesce(data, "metadata", false)), Lists: kept, SummaryURL: urlOf("summary-url"), APIURL: urlOf("api-url")}
}

// HTTPDownloader is the part of Composer\Util\HttpDownloader used by the
// policy layer: Request performs one request with PHP stream-context style
// options (options.http.method/header/timeout/content, options.ssl, ...),
// like HttpDownloader::get($url, $options). Transport failures and non-2xx
// statuses must be returned as *httpdl.TransportError.
type HTTPDownloader interface {
	Request(url string, options *phpjson.Map) (*httpdl.Response, error)
}

// HTTPDownloaderFunc adapts a function to HTTPDownloader.
type HTTPDownloaderFunc func(url string, options *phpjson.Map) (*httpdl.Response, error)

func (f HTTPDownloaderFunc) Request(url string, options *phpjson.Map) (*httpdl.Response, error) {
	return f(url, options)
}

// postOptions builds the stream options of a POST like FilterListApiClient
// and ComposerRepository do: http.method=POST, http.header normalized to a
// list with the content type appended, http.timeout=10, http.content=body.
func postOptions(options *phpjson.Map, contentType, body string) *phpjson.Map {
	opts := phpjson.NewMap()
	if options != nil {
		opts = phpjson.AsMap(phpjson.CloneValue(options))
	}
	httpOpts := phpjson.NewMap()
	if h := opts.Map("http"); h != nil {
		httpOpts = h
	}
	httpOpts.Set("method", "POST")
	var headers []phpjson.Value
	if h, ok := httpOpts.Get("header"); ok {
		if phpjson.IsArray(h) {
			headers = arrayValues(h)
		} else if h != nil {
			headers = []phpjson.Value{h}
		}
	}
	headers = append(headers, "Content-type: "+contentType)
	httpOpts.Set("header", headers)
	httpOpts.Set("timeout", int64(10))
	httpOpts.Set("content", body)
	opts.Set("http", httpOpts)
	return opts
}

// DecodeJSONResponse is Response::decodeJson.
func DecodeJSONResponse(resp *httpdl.Response, url string) (phpjson.Value, error) {
	v, err := phpjson.Decode(resp.Body)
	if err != nil {
		return nil, newErr("Seld\\JsonLint\\ParsingException", `"`+httpdl.Sanitize(url)+`" does not contain valid JSON`)
	}
	return v, nil
}

// FilterListApiClient is Composer\FilterList\FilterListApiClient.
type FilterListApiClient struct {
	HTTP    HTTPDownloader
	Options *phpjson.Map
}

// PostPurls is FilterListApiClient::postPurls.
func (c *FilterListApiClient) PostPurls(url string, packageConstraintMap *ConstraintMap, configuredLists []string) (*httpdl.Response, error) {
	purls := []phpjson.Value{}
	for _, name := range packageConstraintMap.Keys() {
		purls = append(purls, "pkg://composer/"+name)
	}
	lists := []phpjson.Value{}
	for _, l := range configuredLists {
		lists = append(lists, l)
	}
	body := phpjson.NewMap()
	body.Set("packages", purls)
	body.Set("lists", lists)
	enc, err := phpjson.Encode(body, 0)
	if err != nil {
		return nil, err
	}
	if c.HTTP == nil {
		return nil, newErr(LogicException, "No HTTP downloader configured to POST to "+httpdl.Sanitize(url))
	}
	return c.HTTP.Request(url, postOptions(c.Options, "application/json", enc))
}

// FilterListProvider is Composer\Repository\FilterListProviderInterface.
type FilterListProvider interface {
	HasFilter() (bool, error)
	GetFilter(packageConstraintMap *ConstraintMap, configuredLists []string) (*FilterByList, error)
	GetFilterLists() ([]string, error)
}

// UrlSourceFilterListProvider is Composer\FilterList\FilterListProvider\UrlSourceFilterListProvider.
type UrlSourceFilterListProvider struct {
	client *FilterListApiClient
	source *UrlSource
}

func NewUrlSourceFilterListProvider(http HTTPDownloader, source *UrlSource) *UrlSourceFilterListProvider {
	return &UrlSourceFilterListProvider{client: &FilterListApiClient{HTTP: http}, source: source}
}

func (p *UrlSourceFilterListProvider) HasFilter() (bool, error) { return true, nil }

func (p *UrlSourceFilterListProvider) GetFilter(packageConstraintMap *ConstraintMap, _ []string) (*FilterByList, error) {
	resp, err := p.client.PostPurls(p.source.URL, packageConstraintMap, []string{p.source.ListName})
	if err != nil {
		return nil, err
	}
	decoded, err := DecodeJSONResponse(resp, p.source.URL)
	if err != nil {
		return nil, err
	}
	var entries phpjson.Value = []phpjson.Value{}
	if isset(decoded, "filter") && phpjson.IsArray(coalesce(decoded, "filter", nil)) {
		entries = coalesce(decoded, "filter", nil)
	}
	raw := phpjson.NewMap()
	raw.Set(p.source.ListName, entries)
	return BuildFilterListEntries(raw, packageConstraintMap, "")
}

func (p *UrlSourceFilterListProvider) GetFilterLists() ([]string, error) {
	return []string{p.source.ListName}, nil
}

// isTransportError reports whether err is a Composer TransportException.
func isTransportError(err error) (*httpdl.TransportError, bool) {
	var te *httpdl.TransportError
	if errors.As(err, &te) {
		return te, true
	}
	return nil, false
}

// FilterListProviderSet is Composer\FilterList\FilterListProvider\FilterListProviderSet.
type FilterListProviderSet struct {
	providers                 []FilterListProvider
	unreachableRepoExceptions []error
}

// NewFilterListProviderSet is the FilterListProviderSet constructor:
// repositories implementing FilterListProvider (directly or through
// ProviderResolver) with HasFilter() join the url-source providers.
// Transport errors raised by HasFilter are kept for getMatchingFilterLists.
func NewFilterListProviderSet(repositories []repo.Repository, sources []FilterListProvider, resolver *ProviderResolver) (*FilterListProviderSet, error) {
	set := &FilterListProviderSet{providers: append([]FilterListProvider(nil), sources...)}
	for _, r := range repositories {
		p := resolver.FilterListProviderFor(r)
		if p == nil {
			continue
		}
		has, err := p.HasFilter()
		if err != nil {
			if _, ok := isTransportError(err); ok {
				set.unreachableRepoExceptions = append(set.unreachableRepoExceptions, err)
				continue
			}
			return nil, err
		}
		if has {
			set.providers = append(set.providers, p)
		}
	}
	return set, nil
}

// CreateFilterListProviderSet is FilterListProviderSet::create.
func CreateFilterListProviderSet(config *PolicyConfig, repositories []repo.Repository, http HTTPDownloader, resolver *ProviderResolver) (*FilterListProviderSet, error) {
	var sources []FilterListProvider
	config.GetCustomListsWithSources().Each(func(_ string, l *CustomListPolicyConfig) {
		for _, s := range l.Sources {
			sources = append(sources, NewUrlSourceFilterListProvider(http, s))
		}
	})
	return NewFilterListProviderSet(repositories, sources, resolver)
}

// FilterListResult is array{filter, unreachableRepos}.
type FilterListResult struct {
	Filter           *FilterByList
	UnreachableRepos []string
}

// FilterListSource is the mockable part of FilterListProviderSet used by the Auditor.
type FilterListSource interface {
	GetMatchingFilterLists(packages []*pkg.Package, configuredLists []string, ignoreUnreachable bool) (*FilterListResult, error)
}

// GetMatchingFilterLists is FilterListProviderSet::getMatchingFilterLists.
func (s *FilterListProviderSet) GetMatchingFilterLists(packages []*pkg.Package, configuredLists []string, ignoreUnreachable bool) (*FilterListResult, error) {
	pcm := constraintMapFor(packages, false)
	var unreachable []string
	filters, err := s.getFilterListEntriesForConstraints(pcm, configuredLists, ignoreUnreachable, &unreachable)
	if err != nil {
		return nil, err
	}
	return &FilterListResult{Filter: filters, UnreachableRepos: unreachable}, nil
}

// constraintMapFor builds the flat name => OR(== version) map of
// RepositorySet::getMatchingSecurityAdvisories / getMatchingFilterLists.
func constraintMapFor(packages []*pkg.Package, skipUnversionedBranches bool) *ConstraintMap {
	byName := NewOMap[*OMap[semver.Constraint]]()
	for _, p := range packages {
		if isRootPackageAlias(p) {
			continue
		}
		if skipUnversionedBranches && IsUnversionedBranch(p) {
			continue
		}
		m, ok := byName.Get(p.Name)
		if !ok {
			m = NewOMap[semver.Constraint]()
			byName.Set(p.Name, m)
		}
		m.Set(p.Version(), semver.NewConstraint("=", p.Version()))
	}
	out := NewOMap[semver.Constraint]()
	byName.Each(func(name string, cs *OMap[semver.Constraint]) {
		var list []semver.Constraint
		cs.Each(func(_ string, c semver.Constraint) { list = append(list, c) })
		out.Set(name, semver.CreateMulti(list, false))
	})
	return out
}

func (s *FilterListProviderSet) getFilterListEntriesForConstraints(pcm *ConstraintMap, configuredLists []string, ignoreUnreachable bool, unreachable *[]string) (*FilterByList, error) {
	for _, e := range s.unreachableRepoExceptions {
		if !ignoreUnreachable {
			return nil, e
		}
		*unreachable = append(*unreachable, e.Error())
	}
	filters := NewOMap[[]*FilterListEntry]()
	for _, provider := range s.providers {
		providerLists, err := provider.GetFilterLists()
		if err != nil {
			return nil, err
		}
		var relevant []string
		for _, l := range configuredLists {
			if inStrings(l, providerLists) {
				relevant = append(relevant, l)
			}
		}
		if len(relevant) == 0 {
			continue
		}
		repoFilter, err := provider.GetFilter(pcm, relevant)
		if err != nil {
			if _, ok := isTransportError(err); ok {
				if !ignoreUnreachable {
					return nil, err
				}
				*unreachable = append(*unreachable, err.Error())
				continue
			}
			return nil, err
		}
		repoFilter.Each(func(listName string, entries []*FilterListEntry) {
			if !inStrings(listName, configuredLists) || !inStrings(listName, providerLists) {
				return
			}
			for _, entry := range entries {
				pc, ok := pcm.Get(entry.PackageName)
				if !ok || !entry.Constraint.Matches(pc) {
					continue
				}
				appendTo(filters, listName, entry)
			}
		})
	}
	filters.KSort()
	return filters, nil
}

// FilterListMap is array<packageName, array<listName, list<FilterListEntry>>>.
type FilterListMap = OMap[*FilterByList]

// FilterListAuditor is Composer\FilterList\FilterListAuditor.
type FilterListAuditor struct{}

// FilterListCollection is the collectFilterLists result.
type FilterListCollection struct {
	Filter           *FilterListMap
	UnreachableRepos []string
}

// CollectFilterLists is FilterListAuditor::collectFilterLists.
func (FilterListAuditor) CollectFilterLists(packages []*pkg.Package, providerSet FilterListSource, configuredLists []string, ignoreUnreachable bool) (*FilterListCollection, error) {
	result, err := providerSet.GetMatchingFilterLists(packages, configuredLists, ignoreUnreachable)
	if err != nil {
		return nil, err
	}
	m := NewOMap[*FilterByList]()
	result.Filter.Each(func(_ string, entries []*FilterListEntry) {
		for _, entry := range entries {
			byList, ok := m.Get(entry.PackageName)
			if !ok {
				byList = NewOMap[[]*FilterListEntry]()
				m.Set(entry.PackageName, byList)
			}
			appendTo(byList, entry.ListName, entry)
		}
	})
	m.KSort()
	return &FilterListCollection{Filter: m, UnreachableRepos: result.UnreachableRepos}, nil
}

// GetMatchingAuditEntries is FilterListAuditor::getMatchingAuditEntries.
func (a FilterListAuditor) GetMatchingAuditEntries(p *pkg.Package, m *FilterListMap, config *PolicyConfig) []*FilterListEntry {
	return a.matchingEntries(p, m, config.GetActiveAuditFilterLists(), OperationAudit)
}

// GetMatchingBlockEntries is FilterListAuditor::getMatchingBlockEntries.
func (a FilterListAuditor) GetMatchingBlockEntries(p *pkg.Package, m *FilterListMap, config *PolicyConfig, blockScope string) []*FilterListEntry {
	return a.matchingEntries(p, m, config.GetActiveBlockFilterLists(blockScope), OperationBlock)
}

func (a FilterListAuditor) matchingEntries(p *pkg.Package, m *FilterListMap, active *OMap[ListPolicy], op string) []*FilterListEntry {
	if isRootPackage(p) || m.Len() == 0 {
		return nil
	}
	if l, ok := active.Get(MalwareName); ok {
		if mw, ok := l.(*MalwarePolicyConfig); ok {
			m = applyMalwareIgnoreSource(m, mw)
		}
	}
	var names []string
	active.Each(func(_ string, l ListPolicy) { names = append(names, l.GetIgnoreForOperation(op).Keys()...) })
	ignoredRe := PackageNamesToRegexp(names)
	pc := semver.NewConstraint("==", p.Version())

	var matching []*FilterListEntry
	for _, name := range p.Names(false) {
		byList, ok := m.Get(name)
		if !ok {
			continue
		}
		entries := NewOMap[[]*FilterListEntry]()
		byList.Each(func(listName string, es []*FilterListEntry) {
			if active.Has(listName) {
				entries.Set(listName, es)
			}
		})
		if ignoredRe.MatchString(name) {
			unfiltered := NewOMap[[]*FilterListEntry]()
			entries.Each(func(listName string, es []*FilterListEntry) {
				l, _ := active.Get(listName)
				if isPackageIgnored(name, pc, l, op) {
					return
				}
				unfiltered.Set(listName, es)
			})
			entries = unfiltered
		}
		entries.Each(func(_ string, es []*FilterListEntry) {
			for _, e := range es {
				if e.Constraint.Matches(pc) {
					matching = append(matching, e)
				}
			}
		})
	}
	return matching
}

func isPackageIgnored(name string, pc semver.Constraint, l ListPolicy, op string) bool {
	found := false
	l.GetIgnoreForOperation(op).Each(func(_ string, rules []*IgnorePackageRule) {
		for _, r := range rules {
			if !found && r.PackageNameRegex.MatchString(name) && r.Constraint.Matches(pc) {
				found = true
			}
		}
	})
	return found
}

func applyMalwareIgnoreSource(m *FilterListMap, mw *MalwarePolicyConfig) *FilterListMap {
	if len(mw.IgnoreSource) == 0 {
		return m
	}
	out := NewOMap[*FilterByList]()
	m.Each(func(name string, byList *FilterByList) {
		entries, ok := byList.Get(MalwareName)
		if !ok {
			out.Set(name, byList)
			return
		}
		var kept []*FilterListEntry
		for _, e := range entries {
			if e.Source == nil || !inStrings(*e.Source, mw.IgnoreSource) {
				kept = append(kept, e)
			}
		}
		nb := byList.Clone()
		if len(kept) > 0 {
			nb.Set(MalwareName, kept)
		} else {
			nb.Delete(MalwareName)
			if nb.Len() == 0 {
				return
			}
		}
		out.Set(name, nb)
	})
	return out
}
