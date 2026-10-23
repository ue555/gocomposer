package solver

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

const loadBatchSize = 50

// PoolBuilder is Composer\DependencyResolver\PoolBuilder (without the
// security-advisory, filter-list and cooldown filters, see docs/compatibility.md).
type PoolBuilder struct {
	acceptable           map[string]int
	stabilityFlags       map[string]int
	rootAliases          map[string]map[string]RootAliasEntry
	rootReferences       map[string]string
	temporaryConstraints map[string]semver.Constraint
	Optimizer            *PoolOptimizer
	IgnoredTypes         []string
	AllowedTypes         []string // nil means all
	Warn                 func(string)
	// PolicyFilter runs the dependency policy pool filters (security
	// advisories, filter lists, cooldown) before the optimizer.
	PolicyFilter func(pool *Pool, repos []repo.Repository, req *Request) (*Pool, error)
	// Write writes to the standard output at a verbosity (IOInterface::write).
	Write func(msg string, verbosity int)
	// PrePoolCreate dispatches PluginEvents::PRE_POOL_CREATE when the
	// builder has an event dispatcher: listeners may replace the packages
	// and the unacceptable fixed/locked packages the pool is created from.
	PrePoolCreate func(ev *PrePoolCreateEvent) error

	aliasMap          map[*pkg.Package][]aliasEntry
	packagesToLoad    *ncMap
	loadedPackages    map[string]semver.Constraint
	loadedPerRepo     map[int]map[string]map[string]bool
	packages          []*poolEntry
	unacceptable      []*pkg.Package
	skippedLoad       map[string][]*pkg.Package
	restricted        map[string]bool
	pathRepoUnlocked  map[string]bool
	maxExtendedReqs   map[string]bool
	updateAllowWarned map[string]bool
	indexCounter      int
	updateAllowList   []string
}

// PrePoolCreateEvent is the data of PluginEvents::PRE_POOL_CREATE.
type PrePoolCreateEvent struct {
	Repositories   []repo.Repository
	Request        *Request
	Acceptable     map[string]int
	StabilityFlags map[string]int
	RootAliases    map[string]map[string]RootAliasEntry
	RootReferences map[string]string
	Packages       []*pkg.Package
	Unacceptable   []*pkg.Package
}

type poolEntry struct {
	index   int
	pkg     *pkg.Package
	removed bool
}

type aliasEntry struct {
	index int
	pkg   *pkg.Package
}

// ncMap is an ordered name => constraint map (PHP array).
type ncMap struct {
	keys []string
	vals map[string]semver.Constraint
}

func newNCMap() *ncMap { return &ncMap{vals: map[string]semver.Constraint{}} }

func (m *ncMap) set(k string, c semver.Constraint) {
	if _, ok := m.vals[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = c
}

func (m *ncMap) get(k string) (semver.Constraint, bool) { c, ok := m.vals[k]; return c, ok }

func (m *ncMap) del(k string) {
	if _, ok := m.vals[k]; !ok {
		return
	}
	delete(m.vals, k)
	for i, key := range m.keys {
		if key == k {
			m.keys = append(m.keys[:i:i], m.keys[i+1:]...)
			return
		}
	}
}

func (m *ncMap) len() int { return len(m.keys) }

func NewPoolBuilder() *PoolBuilder { return &PoolBuilder{} }

func (b *PoolBuilder) warn(msg string) {
	if b.Warn != nil {
		b.Warn(msg)
	}
}

func (b *PoolBuilder) reset() {
	b.aliasMap = map[*pkg.Package][]aliasEntry{}
	b.packagesToLoad = newNCMap()
	b.loadedPackages = map[string]semver.Constraint{}
	b.loadedPerRepo = map[int]map[string]map[string]bool{}
	b.packages = nil
	b.unacceptable = nil
	b.skippedLoad = map[string][]*pkg.Package{}
	b.pathRepoUnlocked = map[string]bool{}
	b.maxExtendedReqs = map[string]bool{}
	b.updateAllowWarned = map[string]bool{}
	b.indexCounter = 0
}

func (b *PoolBuilder) livePackages() []*pkg.Package {
	var out []*pkg.Package
	for _, e := range b.packages {
		if !e.removed {
			out = append(out, e.pkg)
		}
	}
	return out
}

// BuildPool is PoolBuilder::buildPool.
func (b *PoolBuilder) BuildPool(repositories []repo.Repository, req *Request) (*Pool, error) {
	b.reset()
	b.restricted = req.restricted
	if len(req.UpdateAllowList()) > 0 {
		b.updateAllowList = req.UpdateAllowList()
		b.warnAboutNonMatchingUpdateAllowList(req)
		if req.LockedRepository == nil {
			panic("No lock repo present and yet a partial update was requested.")
		}
		for _, lp := range req.LockedRepository.MustPackages() {
			if !req.IsInUpdateAllowList(lp.Name) {
				b.skippedLoad[lp.Name] = append(b.skippedLoad[lp.Name], lp)
				for _, l := range lp.Replaces() {
					b.skippedLoad[l.Target] = append(b.skippedLoad[l.Target], lp)
				}
				if lp.DistType() == "path" {
					to := lp.TransportOptions()
					sym, has := to.Get("symlink")
					if !has || sym != false {
						b.pathRepoUnlocked[lp.Name] = true
						continue
					}
				}
				req.LockPackage(lp)
			}
		}
	}
	for _, p := range req.FixedOrLockedPackages() {
		b.loadedPackages[p.Name] = &semver.MatchAll{}
		for _, l := range p.Replaces() {
			b.loadedPackages[l.Target] = &semver.MatchAll{}
		}
		kind := repoKind(p)
		if kind == repo.KindRoot || kind == repo.KindPlatform || repo.IsPackageAcceptable(b.acceptable, b.stabilityFlags, p.Names(true), p.Stability()) {
			if err := b.loadPackage(req, repositories, p, false); err != nil {
				return nil, err
			}
		} else {
			b.unacceptable = append(b.unacceptable, p)
		}
	}
	for _, r := range req.Requires() {
		if _, ok := b.loadedPackages[r.Name]; ok {
			continue
		}
		b.packagesToLoad.set(r.Name, r.Constraint)
		b.maxExtendedReqs[r.Name] = true
	}
	for _, name := range append([]string(nil), b.packagesToLoad.keys...) {
		if _, ok := b.loadedPackages[name]; ok {
			b.packagesToLoad.del(name)
		}
	}
	for b.packagesToLoad.len() > 0 {
		if err := b.loadPackagesMarkedForLoading(req, repositories); err != nil {
			return nil, err
		}
	}
	if len(b.temporaryConstraints) > 0 {
		for _, e := range b.packages {
			if e.removed || e.pkg.AliasOf != nil {
				continue
			}
			for _, name := range e.pkg.Names(true) {
				c, ok := b.temporaryConstraints[name]
				if !ok {
					continue
				}
				group := []*poolEntry{e}
				for _, a := range b.aliasMap[e.pkg] {
					for _, e2 := range b.packages {
						if e2.index == a.index {
							group = append(group, e2)
						}
					}
				}
				found := false
				for _, g := range group {
					if semver.MatchVersion(c, g.pkg.Version()) {
						found = true
					}
				}
				if !found {
					for _, g := range group {
						g.removed = true
					}
				}
			}
		}
	}
	packages, unacceptable := b.livePackages(), b.unacceptable
	if b.PrePoolCreate != nil {
		ev := &PrePoolCreateEvent{Repositories: repositories, Request: req, Acceptable: b.acceptable, StabilityFlags: b.stabilityFlags,
			RootAliases: b.rootAliases, RootReferences: b.rootReferences, Packages: packages, Unacceptable: unacceptable}
		if err := b.PrePoolCreate(ev); err != nil {
			return nil, err
		}
		packages, unacceptable = ev.Packages, ev.Unacceptable
	}
	pool := NewPool(packages, unacceptable)
	if b.PolicyFilter != nil {
		filtered, err := b.PolicyFilter(pool, repositories, req)
		if err != nil {
			return nil, err
		}
		pool = filtered
	}
	if b.Optimizer != nil {
		before := time.Now()
		total := len(pool.Packages())
		pool = b.Optimizer.Optimize(req, pool)
		if filtered := total - len(pool.Packages()); filtered > 0 && b.Write != nil {
			b.Write(fmt.Sprintf("Pool optimizer completed in %.3f seconds", time.Since(before).Seconds()), 128)
			b.Write(fmt.Sprintf("<info>Found %s package versions referenced in your dependency graph. %s (%d%%) were optimized away.</info>", numberFormat(total), numberFormat(filtered), int(100/float64(total)*float64(filtered)+0.5)), 128)
		}
	}
	semver.ClearIntervals()
	return pool, nil
}

func repoKind(p *pkg.Package) repo.Kind {
	if r, ok := p.Repo.(repo.Repository); ok {
		return r.Kind()
	}
	return -1
}

func (b *PoolBuilder) markPackageNameForLoading(req *Request, name string, c semver.Constraint) {
	if repo.IsPlatformPackage(name) {
		return
	}
	if b.maxExtendedReqs[name] {
		return
	}
	if rr, ok := req.Require(name); ok && !semver.IsSubsetOf(c, rr) {
		c = rr
	}
	if _, loaded := b.loadedPackages[name]; !loaded {
		if existing, ok := b.packagesToLoad.get(name); ok {
			if semver.IsSubsetOf(c, existing) {
				return
			}
			c = semver.CompactConstraint(semver.CreateMulti([]semver.Constraint{existing, c}, false))
		}
		b.packagesToLoad.set(name, c)
		return
	}
	if semver.IsSubsetOf(c, b.loadedPackages[name]) {
		return
	}
	b.packagesToLoad.set(name, semver.CompactConstraint(semver.CreateMulti([]semver.Constraint{b.loadedPackages[name], c}, false)))
	delete(b.loadedPackages, name)
}

func (b *PoolBuilder) loadPackagesMarkedForLoading(req *Request, repositories []repo.Repository) error {
	for _, name := range append([]string(nil), b.packagesToLoad.keys...) {
		c, _ := b.packagesToLoad.get(name)
		if b.restricted != nil && !b.restricted[name] {
			b.packagesToLoad.del(name)
			continue
		}
		b.loadedPackages[name] = c
	}
	// array_chunk(..., 50, true)
	var all []repo.NameConstraint
	for _, k := range b.packagesToLoad.keys {
		all = append(all, repo.NameConstraint{Name: k, Constraint: b.packagesToLoad.vals[k]})
	}
	batches := chunk(all)
	b.packagesToLoad = newNCMap()
	for repoIndex, r := range repositories {
		if r.Kind() == repo.KindPlatform || (req.LockedRepository != nil && r == repo.Repository(req.LockedRepository)) {
			continue
		}
		if len(batches) == 0 {
			break
		}
		for bi, batch := range batches {
			res, err := r.LoadPackages(batch, b.acceptable, b.stabilityFlags, b.loadedPerRepo[repoIndex])
			if err != nil {
				return err
			}
			for _, nameFound := range res.NamesFound {
				var kept []repo.NameConstraint
				for _, nc := range batches[bi] {
					if nc.Name != nameFound {
						kept = append(kept, nc)
					}
				}
				batches[bi] = kept
			}
			for _, p := range res.Packages {
				if b.loadedPerRepo[repoIndex] == nil {
					b.loadedPerRepo[repoIndex] = map[string]map[string]bool{}
				}
				if b.loadedPerRepo[repoIndex][p.Name] == nil {
					b.loadedPerRepo[repoIndex][p.Name] = map[string]bool{}
				}
				b.loadedPerRepo[repoIndex][p.Name][p.Version()] = true
				if containsStr(b.IgnoredTypes, p.Type()) || (b.AllowedTypes != nil && !containsStr(b.AllowedTypes, p.Type())) {
					continue
				}
				if err := b.loadPackage(req, repositories, p, !b.pathRepoUnlocked[p.Name]); err != nil {
					return err
				}
			}
		}
		var merged []repo.NameConstraint
		for _, bt := range batches {
			merged = append(merged, bt...)
		}
		batches = chunk(merged)
	}
	return nil
}

func chunk(all []repo.NameConstraint) [][]repo.NameConstraint {
	var out [][]repo.NameConstraint
	for i := 0; i < len(all); i += loadBatchSize {
		end := i + loadBatchSize
		if end > len(all) {
			end = len(all)
		}
		out = append(out, all[i:end:end])
	}
	return out
}

func containsStr(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

func (b *PoolBuilder) loadPackage(req *Request, repositories []repo.Repository, p *pkg.Package, propagateUpdate bool) error {
	index := b.indexCounter
	b.indexCounter++
	b.packages = append(b.packages, &poolEntry{index: index, pkg: p})
	if p.AliasOf != nil {
		b.aliasMap[p.AliasOf] = append(b.aliasMap[p.AliasOf], aliasEntry{index, p})
	}
	name := p.Name
	if ref, ok := b.rootReferences[name]; ok {
		if !req.IsLockedPackage(p) && !req.IsFixedPackage(p) {
			p.SetSourceDistReferences(ref)
		}
	}
	if propagateUpdate || b.pathRepoUnlocked[p.Name] {
		if alias, ok := b.rootAliases[name][p.Version()]; ok {
			base := p
			if p.AliasOf != nil {
				base = p.AliasOf
			}
			aliasPkg := pkg.NewAlias(base, alias.AliasNormalized, alias.Alias)
			aliasPkg.RootPackageAlias = true
			newIndex := b.indexCounter
			b.indexCounter++
			b.packages = append(b.packages, &poolEntry{index: newIndex, pkg: aliasPkg})
			b.aliasMap[aliasPkg.AliasOf] = append(b.aliasMap[aliasPkg.AliasOf], aliasEntry{newIndex, aliasPkg})
		}
	}
	for _, link := range p.Requires() {
		require := link.Target
		if _, skipped := b.skippedLoad[require]; skipped {
			if propagateUpdate && req.UpdateAllowTransitiveDependencies() {
				skippedRoot := b.getSkippedRootRequires(req, require)
				if req.UpdateAllowTransitiveRootDependencies() || len(skippedRoot) == 0 {
					if err := b.unlockPackage(req, repositories, require); err != nil {
						return err
					}
					b.markPackageNameForLoading(req, require, link.Constraint)
				} else {
					for _, rr := range skippedRoot {
						if !b.updateAllowWarned[rr] {
							b.updateAllowWarned[rr] = true
							b.warn("<warning>Dependency " + rr + " is also a root requirement. Package has not been listed as an update argument, so keeping locked at old version. Use --with-all-dependencies (-W) to include root dependencies.</warning>")
						}
					}
				}
			} else if b.pathRepoUnlocked[require] {
				if _, loaded := b.loadedPackages[require]; !loaded {
					b.markPackageNameForLoading(req, require, link.Constraint)
				}
			}
		} else {
			b.markPackageNameForLoading(req, require, link.Constraint)
		}
	}
	if propagateUpdate && req.UpdateAllowTransitiveDependencies() {
		for _, link := range p.Replaces() {
			replace := link.Target
			_, loaded := b.loadedPackages[replace]
			_, skipped := b.skippedLoad[replace]
			if loaded && skipped {
				skippedRoot := b.getSkippedRootRequires(req, replace)
				if req.UpdateAllowTransitiveRootDependencies() || len(skippedRoot) == 0 {
					if err := b.unlockPackage(req, repositories, replace); err != nil {
						return err
					}
					b.markPackageNameForLoadingIfRequired(req, replace)
				} else {
					for _, rr := range skippedRoot {
						if !b.updateAllowWarned[rr] {
							b.updateAllowWarned[rr] = true
							b.warn("<warning>Dependency " + rr + " is also a root requirement. Package has not been listed as an update argument, so keeping locked at old version. Use --with-all-dependencies (-W) to include root dependencies.</warning>")
						}
					}
				}
			}
		}
	}
	return nil
}

func (b *PoolBuilder) isRootRequire(req *Request, name string) bool {
	_, ok := req.Require(name)
	return ok
}

func (b *PoolBuilder) getSkippedRootRequires(req *Request, name string) []string {
	skipped, ok := b.skippedLoad[name]
	if !ok {
		return nil
	}
	var matches []string
	if b.isRootRequire(req, name) {
		for _, p := range skipped {
			if name != p.Name {
				matches = append(matches, p.Name+" (via replace of "+name+")")
			} else {
				matches = append(matches, p.Name)
			}
		}
		return matches
	}
	for _, p := range skipped {
		if b.isRootRequire(req, p.Name) {
			matches = append(matches, p.Name)
		}
		for _, l := range p.Replaces() {
			if b.isRootRequire(req, l.Target) {
				if name != p.Name {
					matches = append(matches, p.Name+" (via replace of "+name+")")
				} else {
					matches = append(matches, p.Name)
				}
				break
			}
		}
	}
	return matches
}

func (b *PoolBuilder) warnAboutNonMatchingUpdateAllowList(req *Request) {
	if req.LockedRepository == nil {
		panic("No lock repo present and yet a partial update was requested.")
	}
patterns:
	for _, pattern := range b.updateAllowList {
		matchedPlatform := false
		re := PackageNameToRegexp(pattern)
		for _, p := range req.LockedRepository.MustPackages() {
			if re.MatchString(p.Name) {
				continue patterns
			}
		}
		for _, r := range req.Requires() {
			if re.MatchString(r.Name) {
				if repo.IsPlatformPackage(r.Name) {
					matchedPlatform = true
					continue
				}
				continue patterns
			}
		}
		if matchedPlatform {
			b.warn(`<warning>Pattern "` + pattern + `" listed for update matches platform packages, but these cannot be updated by Composer.</warning>`)
		} else if strings.Contains(pattern, "*") {
			b.warn(`<warning>Pattern "` + pattern + `" listed for update does not match any locked packages.</warning>`)
		} else {
			b.warn(`<warning>Package "` + pattern + `" listed for update is not locked.</warning>`)
		}
	}
}

func (b *PoolBuilder) unlockPackage(req *Request, repositories []repo.Repository, name string) error {
	for _, por := range append([]*pkg.Package(nil), b.skippedLoad[name]...) {
		if por.Name != name {
			if _, ok := b.skippedLoad[por.Name]; ok {
				replacerName := por.Name
				if req.UpdateAllowTransitiveRootDependencies() || (!b.isRootRequire(req, name) && !b.isRootRequire(req, replacerName)) {
					if err := b.unlockPackage(req, repositories, replacerName); err != nil {
						return err
					}
					if b.isRootRequire(req, replacerName) {
						b.markPackageNameForLoading(req, replacerName, &semver.MatchAll{})
					} else {
						for _, lp := range b.livePackages() {
							if l := lp.Requires().Get(replacerName); l != nil {
								b.markPackageNameForLoading(req, replacerName, l.Constraint)
							}
						}
					}
				}
			}
		}
	}
	if b.pathRepoUnlocked[name] {
		for _, e := range b.packages {
			if !e.removed && e.pkg.Name == name {
				b.removeLoadedPackage(repositories, e)
			}
		}
	}
	delete(b.skippedLoad, name)
	delete(b.loadedPackages, name)
	delete(b.maxExtendedReqs, name)
	delete(b.pathRepoUnlocked, name)
	for _, lp := range req.LockedPackages() {
		if lp.AliasOf != nil || lp.Name != name {
			continue
		}
		var entry *poolEntry
		for _, e := range b.packages {
			if !e.removed && e.pkg == lp {
				entry = e
				break
			}
		}
		if entry == nil {
			continue
		}
		req.UnlockPackage(lp)
		b.removeLoadedPackage(repositories, entry)
		for _, fol := range req.FixedOrLockedPackages() {
			if fol == lp {
				continue
			}
			if _, ok := b.skippedLoad[fol.Name]; ok {
				requires := fol.Requires()
				if l := requires.Get(lp.Name); l != nil {
					b.markPackageNameForLoading(req, lp.Name, l.Constraint)
				}
				for _, replace := range lp.Replaces() {
					if requires.Get(replace.Target) != nil {
						if _, ok := b.skippedLoad[replace.Target]; ok {
							if err := b.unlockPackage(req, repositories, replace.Target); err != nil {
								return err
							}
							b.markPackageNameForLoading(req, replace.Target, replace.Constraint)
						}
					}
				}
			}
		}
	}
	return nil
}

func (b *PoolBuilder) markPackageNameForLoadingIfRequired(req *Request, name string) {
	if rr, ok := req.Require(name); ok {
		b.markPackageNameForLoading(req, name, rr)
	}
	for _, p := range b.livePackages() {
		for _, l := range p.Requires() {
			if name == l.Target {
				b.markPackageNameForLoading(req, l.Target, l.Constraint)
			}
		}
	}
}

func (b *PoolBuilder) removeLoadedPackage(repositories []repo.Repository, e *poolEntry) {
	p := e.pkg
	repoIndex := -1
	for i, r := range repositories {
		if pr, ok := p.Repo.(repo.Repository); ok && pr == r {
			repoIndex = i
			break
		}
	}
	if m := b.loadedPerRepo[repoIndex]; m != nil && m[p.Name] != nil {
		delete(m[p.Name], p.Version())
	}
	e.removed = true
	if aliases, ok := b.aliasMap[p]; ok {
		for _, a := range aliases {
			if m := b.loadedPerRepo[repoIndex]; m != nil && m[a.pkg.Name] != nil {
				delete(m[a.pkg.Name], a.pkg.Version())
			}
			for _, e2 := range b.packages {
				if e2.index == a.index {
					e2.removed = true
				}
			}
		}
		delete(b.aliasMap, p)
	}
}

// numberFormat is PHP number_format($n).
func numberFormat(n int) string {
	str := strconv.Itoa(n)
	var out []byte
	for i := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, str[i])
	}
	return string(out)
}
