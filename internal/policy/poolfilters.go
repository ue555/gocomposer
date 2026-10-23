// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/DependencyResolver/SecurityAdvisoryPoolFilter.php,
// FilterListPoolFilter.php, CooldownPoolFilter.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"time"

	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// ------------------------------------------------ SecurityAdvisoryPoolFilter

// SecurityAdvisoryPoolFilter is Composer\DependencyResolver\SecurityAdvisoryPoolFilter.
type SecurityAdvisoryPoolFilter struct {
	auditor  *Auditor
	config   *PolicyConfig
	io       IO
	resolver *ProviderResolver
}

// NewSecurityAdvisoryPoolFilter mirrors the constructor; resolver maps the
// repositories passed to Filter to advisory providers.
func NewSecurityAdvisoryPoolFilter(auditor *Auditor, config *PolicyConfig, io IO, resolver *ProviderResolver) *SecurityAdvisoryPoolFilter {
	if auditor == nil {
		auditor = NewAuditor(nil)
	}
	return &SecurityAdvisoryPoolFilter{auditor: auditor, config: config, io: orNullIO(io), resolver: resolver}
}

func advisoryTarget(p *pkg.Package) *pkg.Package {
	if isRootPackageAlias(p) {
		return p.AliasOf
	}
	return p
}

// Filter is SecurityAdvisoryPoolFilter::filter.
func (f *SecurityAdvisoryPoolFilter) Filter(pool *Pool, repositories []repo.Repository, request Request) (*Pool, error) {
	advisories := f.config.Advisories
	abandoned := f.config.Abandoned
	if !advisories.Block {
		return pool, nil
	}
	repoSet := &RepositoryAdvisorySource{Repositories: repositories, Resolver: f.resolver}

	var forAdvisories []*pkg.Package
	for _, p := range pool.Packages {
		p = advisoryTarget(p)
		if !isRootPackage(p) && !isPlatformPackage(p.Name) && !request.IsLockedPackage(p) {
			forAdvisories = append(forAdvisories, p)
		}
	}

	ignoreList := advisories.GetIgnoreListForOperation(OperationBlock)
	ignoreUnreachable := f.config.IgnoreUnreachable.Update

	all, err := repoSet.GetMatchingSecurityAdvisories(forAdvisories, true, ignoreUnreachable)
	if err != nil {
		return nil, err
	}
	if f.auditor.NeedsCompleteAdvisoryLoad(all.Advisories, ignoreList) {
		if all, err = repoSet.GetMatchingSecurityAdvisories(forAdvisories, false, ignoreUnreachable); err != nil {
			return nil, err
		}
	}
	if ignoreUnreachable && len(all.UnreachableRepos) > 0 {
		f.io.WriteError("<warning>Security advisory data could not be fetched from some repositories (ignored per policy.ignore-unreachable); matches may be incomplete:</warning>")
		for _, r := range all.UnreachableRepos {
			f.io.WriteError("  - " + r)
		}
	}

	advisoryMap, _ := f.auditor.ProcessAdvisories(all.Advisories, ignoreList, advisories.GetIgnoreSeverityForOperation(OperationBlock))
	ignoreAbandoned := abandoned.GetFlatIgnoreForOperation(OperationBlock)

	removed := pool.PolicyRemoved.Clone()
	var packages []*pkg.Package
	for _, p := range pool.Packages {
		if abandoned.Block && len(f.auditor.FilterAbandonedPackages([]*pkg.Package{p}, ignoreAbandoned)) != 0 {
			removed.Record(p, RemovalReasonAbandoned(p.Name, p.PrettyVersion()))
			continue
		}
		if matching := matchingAdvisories(p, advisoryMap); len(matching) > 0 {
			removed.Record(p, RemovalReasonAdvisories(p.Name, p.PrettyVersion(), matching))
			continue
		}
		packages = append(packages, p)
	}
	return &Pool{Packages: packages, PolicyRemoved: removed}, nil
}

func matchingAdvisories(p *pkg.Package, advisoryMap *AdvisoryMap) []*Advisory {
	p = advisoryTarget(p)
	if IsUnversionedBranch(p) {
		return nil
	}
	var out []*Advisory
	for _, name := range p.Names(false) {
		list, ok := advisoryMap.Get(name)
		if !ok {
			continue
		}
		pc := semver.NewConstraint("==", p.Version())
		for _, a := range list {
			if a.AffectedVersions.Matches(pc) {
				out = append(out, a)
			}
		}
	}
	return out
}

// ------------------------------------------------------ FilterListPoolFilter

// FilterListPoolFilter is Composer\DependencyResolver\FilterListPoolFilter.
type FilterListPoolFilter struct {
	config       *PolicyConfig
	auditor      FilterListAuditor
	http         HTTPDownloader
	blockScope   string
	repositories []repo.Repository
	io           IO
	resolver     *ProviderResolver
}

// NewFilterListPoolFilter mirrors the constructor (policyConfig,
// filterListAuditor, httpDownloader, blockScope, repositories, io).
func NewFilterListPoolFilter(config *PolicyConfig, http HTTPDownloader, blockScope string, repositories []repo.Repository, io IO, resolver *ProviderResolver) *FilterListPoolFilter {
	return &FilterListPoolFilter{config: config, http: http, blockScope: blockScope, repositories: repositories, io: orNullIO(io), resolver: resolver}
}

func isFilterable(p *pkg.Package) bool { return !isRootPackage(p) && !isPlatformPackage(p.Name) }

// Filter is FilterListPoolFilter::filter.
func (f *FilterListPoolFilter) Filter(pool *Pool, request Request) (*Pool, error) {
	checkLockedAgainstInstall := f.blockScope == BlockScopeUpdate

	configuredNames := f.config.GetActiveBlockFilterListNames(f.blockScope)
	var installNames []string
	if checkLockedAgainstInstall {
		installNames = f.config.GetActiveBlockFilterListNames(BlockScopeInstall)
	}
	var union []string
	for _, n := range append(append([]string(nil), configuredNames...), installNames...) {
		if !inStrings(n, union) {
			union = append(union, n)
		}
	}
	if len(union) == 0 {
		return pool, nil
	}

	ignoreUnreachable := f.config.IgnoreUnreachable.ForBlockScope(f.blockScope)
	if checkLockedAgainstInstall {
		ignoreUnreachable = ignoreUnreachable && f.config.IgnoreUnreachable.ForBlockScope(BlockScopeInstall)
	}

	providerSet, err := CreateFilterListProviderSet(f.config, f.repositories, f.http, f.resolver)
	if err != nil {
		return nil, err
	}
	var forFilter []*pkg.Package
	for _, p := range pool.Packages {
		if isFilterable(p) {
			forFilter = append(forFilter, p)
		}
	}
	fetched, err := f.auditor.CollectFilterLists(forFilter, providerSet, union, ignoreUnreachable)
	if err != nil {
		return nil, err
	}
	if len(fetched.UnreachableRepos) > 0 {
		f.io.WriteError("<warning>Filter list data could not be fetched from some sources (ignored per policy.ignore-unreachable); matches may be incomplete:</warning>")
		for _, r := range fetched.UnreachableRepos {
			f.io.WriteError("  - " + r)
		}
	}
	unionMap := fetched.Filter
	if unionMap.Len() == 0 {
		return pool, nil
	}

	configuredMap := filterMapByListNames(unionMap, configuredNames)
	installMap := NewOMap[*FilterByList]()
	var lockedMap map[string]map[string]bool
	if checkLockedAgainstInstall {
		installMap = filterMapByListNames(unionMap, installNames)
		lockedMap = map[string]map[string]bool{}
		for _, lp := range request.LockedRepositoryPackages() {
			if lockedMap[lp.Name] == nil {
				lockedMap[lp.Name] = map[string]bool{}
			}
			lockedMap[lp.Name][lp.Version()] = true
		}
	}

	removed := pool.PolicyRemoved.Clone()
	var packages []*pkg.Package
	for _, p := range pool.Packages {
		if !isFilterable(p) {
			packages = append(packages, p)
			continue
		}
		var matching []*FilterListEntry
		if checkLockedAgainstInstall && (request.IsLockedPackage(p) || lockedMap[p.Name][p.Version()]) {
			matching = f.auditor.GetMatchingBlockEntries(p, installMap, f.config, BlockScopeInstall)
		} else {
			matching = f.auditor.GetMatchingBlockEntries(p, configuredMap, f.config, f.blockScope)
		}
		if len(matching) > 0 {
			removed.Record(p, RemovalReasonFilterList(p.Name, p.PrettyVersion(), matching))
			continue
		}
		packages = append(packages, p)
	}
	return &Pool{Packages: packages, PolicyRemoved: removed}, nil
}

func filterMapByListNames(m *FilterListMap, names []string) *FilterListMap {
	out := NewOMap[*FilterByList]()
	if len(names) == 0 {
		return out
	}
	m.Each(func(pkgName string, byList *FilterByList) {
		byList.Each(func(listName string, entries []*FilterListEntry) {
			if !inStrings(listName, names) {
				return
			}
			nb, ok := out.Get(pkgName)
			if !ok {
				nb = NewOMap[[]*FilterListEntry]()
				out.Set(pkgName, nb)
			}
			nb.Set(listName, entries)
		})
	})
	return out
}

// -------------------------------------------------------- CooldownPoolFilter

// CooldownPoolFilter is Composer\DependencyResolver\CooldownPoolFilter.
type CooldownPoolFilter struct {
	config     *CooldownPolicyConfig
	now        time.Time
	unverified *OMap[[]string]
}

// NewCooldownPoolFilter mirrors the constructor; a zero now means "now".
func NewCooldownPoolFilter(config *CooldownPolicyConfig, now time.Time) *CooldownPoolFilter {
	if now.IsZero() {
		now = timeNow()
	}
	return &CooldownPoolFilter{config: config, now: now, unverified: NewOMap[[]string]()}
}

// Filter is CooldownPoolFilter::filter.
func (f *CooldownPoolFilter) Filter(pool *Pool, request Request) (*Pool, error) {
	if !f.config.HasCooldown() || !f.config.Block {
		return pool, nil
	}
	f.unverified = NewOMap[[]string]()
	removed := pool.PolicyRemoved.Clone()
	var packages []*pkg.Package
	for _, p := range pool.Packages {
		target := p
		if p.AliasOf != nil {
			target = p.AliasOf
		}
		if isRootPackage(target) || isPlatformPackage(target.Name) || request.IsLockedPackage(p) || request.IsLockedPackage(target) ||
			target.IsDev() || f.config.IsIgnored(target, OperationBlock) {
			packages = append(packages, p)
			continue
		}
		if err := f.config.AssertPublishedTimeProvided(target); err != nil {
			return nil, err
		}
		_, hasPublished := PublishedDate(target)
		if !hasPublished && p.AliasOf == nil && target.Base().Repo != nil {
			repoName := target.Base().Repo.RepoName()
			appendTo(f.unverified, repoName, p.PrettyName+" ("+p.PrettyVersion()+")")
		}
		releaseDate, ok := f.config.GetEffectiveDate(target)
		if !ok || !f.config.IsWithinCooldown(releaseDate, f.now) {
			packages = append(packages, p)
			continue
		}
		source := "time"
		if hasPublished {
			source = "published-time"
		}
		removed.Record(p, RemovalReasonCooldown(p.Name, p.PrettyVersion(), formatAtom(releaseDate), f.config.FormatTimeUntilAvailable(releaseDate, f.now, false), source))
	}
	return &Pool{Packages: packages, PolicyRemoved: removed}, nil
}

// UnverifiedVersions is CooldownPoolFilter::getUnverifiedVersions:
// repository name => "name (version)" judged without a repository date.
func (f *CooldownPoolFilter) UnverifiedVersions() *OMap[[]string] { return f.unverified }

func timeNow() time.Time { return now() }
