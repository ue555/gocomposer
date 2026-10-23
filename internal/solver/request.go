package solver

import (
	"regexp"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// Update allow-list modes (Request::UPDATE_*).
const (
	UpdateOnlyListed                     = 0
	UpdateListedWithTransitiveDepsNoRoot = 1
	UpdateListedWithTransitiveDeps       = 2
)

// RequireEntry is one root requirement in insertion order.
type RequireEntry struct {
	Name       string
	Constraint semver.Constraint
}

// Request is Composer\DependencyResolver\Request.
type Request struct {
	LockedRepository *repo.Array
	requires         []RequireEntry
	requireIndex     map[string]int
	fixed            []*pkg.Package
	locked           []*pkg.Package
	fixedLocked      []*pkg.Package
	updateAllowList  []string
	updateAllowRe    *regexp.Regexp
	updateTransitive int
	restricted       map[string]bool
}

func NewRequest(locked *repo.Array) *Request {
	return &Request{LockedRepository: locked, requireIndex: map[string]int{}}
}

// RequireName is Request::requireName.
func (r *Request) RequireName(name string, c semver.Constraint) {
	name = strings.ToLower(name)
	if c == nil {
		c = &semver.MatchAll{}
	}
	if _, ok := r.requireIndex[name]; ok {
		panic("Overwriting requires seems like a bug (" + name + ")")
	}
	r.requireIndex[name] = len(r.requires)
	r.requires = append(r.requires, RequireEntry{name, c})
}

func (r *Request) Requires() []RequireEntry { return r.requires }

func (r *Request) Require(name string) (semver.Constraint, bool) {
	i, ok := r.requireIndex[name]
	if !ok {
		return nil, false
	}
	return r.requires[i].Constraint, true
}

func addUnique(list []*pkg.Package, p *pkg.Package) []*pkg.Package {
	for _, e := range list {
		if e == p {
			return list
		}
	}
	return append(list, p)
}

func removeFrom(list []*pkg.Package, p *pkg.Package) []*pkg.Package {
	for i, e := range list {
		if e == p {
			return append(list[:i:i], list[i+1:]...)
		}
	}
	return list
}

func contains(list []*pkg.Package, p *pkg.Package) bool {
	for _, e := range list {
		if e == p {
			return true
		}
	}
	return false
}

func (r *Request) FixPackage(p *pkg.Package)  { r.fixed = addUnique(r.fixed, p) }
func (r *Request) LockPackage(p *pkg.Package) { r.locked = addUnique(r.locked, p) }
func (r *Request) FixLockedPackage(p *pkg.Package) {
	r.fixed = addUnique(r.fixed, p)
	r.fixedLocked = addUnique(r.fixedLocked, p)
}
func (r *Request) UnlockPackage(p *pkg.Package) { r.locked = removeFrom(r.locked, p) }

// SetUpdateAllowList is Request::setUpdateAllowList.
func (r *Request) SetUpdateAllowList(list []string, transitive int) {
	r.updateAllowList = list
	if len(list) > 0 {
		r.updateAllowRe = PackageNamesToRegexp(list)
	} else {
		r.updateAllowRe = nil
	}
	r.updateTransitive = transitive
}

func (r *Request) IsInUpdateAllowList(name string) bool {
	return r.updateAllowRe != nil && r.updateAllowRe.MatchString(name)
}

func (r *Request) UpdateAllowList() []string { return r.updateAllowList }
func (r *Request) UpdateAllowTransitiveDependencies() bool {
	return r.updateTransitive != UpdateOnlyListed
}
func (r *Request) UpdateAllowTransitiveRootDependencies() bool {
	return r.updateTransitive == UpdateListedWithTransitiveDeps
}

func (r *Request) FixedPackages() []*pkg.Package      { return r.fixed }
func (r *Request) IsFixedPackage(p *pkg.Package) bool { return contains(r.fixed, p) }
func (r *Request) LockedPackages() []*pkg.Package {
	// array_merge() of int-keyed arrays appends, so a package present in both lists appears twice
	return append(append([]*pkg.Package(nil), r.locked...), r.fixedLocked...)
}
func (r *Request) IsLockedPackage(p *pkg.Package) bool {
	return contains(r.locked, p) || contains(r.fixedLocked, p)
}
func (r *Request) FixedOrLockedPackages() []*pkg.Package {
	return append(append([]*pkg.Package(nil), r.fixed...), r.locked...)
}

// PresentMap is Request::getPresentMap(true) keyed by package id.
func (r *Request) PresentMap() map[int]*pkg.Package {
	m := map[int]*pkg.Package{}
	if r.LockedRepository != nil {
		for _, p := range r.LockedRepository.MustPackages() {
			m[p.ID] = p
		}
	}
	for _, p := range r.fixed {
		m[p.ID] = p
	}
	return m
}

// PresentList is getPresentMap() (by object) as an ordered list.
func (r *Request) PresentList() []*pkg.Package {
	var out []*pkg.Package
	if r.LockedRepository != nil {
		for _, p := range r.LockedRepository.MustPackages() {
			out = addUnique(out, p)
		}
	}
	for _, p := range r.fixed {
		out = addUnique(out, p)
	}
	return out
}

func (r *Request) FixedPackagesMap() map[int]*pkg.Package {
	m := map[int]*pkg.Package{}
	for _, p := range r.fixed {
		m[p.ID] = p
	}
	return m
}

func (r *Request) RestrictPackages(names []string) {
	r.restricted = map[string]bool{}
	for _, n := range names {
		r.restricted[n] = true
	}
}

// PackageNameToRegexp is BasePackage::packageNameToRegexp with the default wrap.
func PackageNameToRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(pattern), `\*`, `.*`) + `$`)
}

// PackageNamesToRegexp is BasePackage::packageNamesToRegexp.
func PackageNamesToRegexp(names []string) *regexp.Regexp {
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = strings.ReplaceAll(regexp.QuoteMeta(n), `\*`, `.*`)
	}
	return regexp.MustCompile(`(?i)^(?:` + strings.Join(parts, "|") + `)$`)
}

// PlatformFilter is PlatformRequirementFilterInterface.
type PlatformFilter interface {
	IsIgnored(req string) bool
	IsUpperBoundIgnored(req string) bool
}

type IgnoreNothing struct{}

func (IgnoreNothing) IsIgnored(string) bool           { return false }
func (IgnoreNothing) IsUpperBoundIgnored(string) bool { return false }

type IgnoreAll struct{}

func (IgnoreAll) IsIgnored(req string) bool           { return repo.IsPlatformPackage(req) }
func (IgnoreAll) IsUpperBoundIgnored(req string) bool { return repo.IsPlatformPackage(req) }

// IgnoreList is IgnoreListPlatformRequirementFilter.
type IgnoreList struct {
	ignore, upper *regexp.Regexp
}

func NewIgnoreList(reqs []string) *IgnoreList {
	var all, upper []string
	for _, r := range reqs {
		if strings.HasSuffix(r, "+") {
			upper = append(upper, r[:len(r)-1])
		} else {
			all = append(all, r)
		}
	}
	return &IgnoreList{ignore: PackageNamesToRegexp(all), upper: PackageNamesToRegexp(upper)}
}

func (f *IgnoreList) IsIgnored(req string) bool {
	return repo.IsPlatformPackage(req) && f.ignore.MatchString(req)
}

func (f *IgnoreList) IsUpperBoundIgnored(req string) bool {
	if !repo.IsPlatformPackage(req) {
		return false
	}
	return f.IsIgnored(req) || f.upper.MatchString(req)
}

// FilterConstraint is IgnoreListPlatformRequirementFilter::filterConstraint.
func (f *IgnoreList) FilterConstraint(req string, c semver.Constraint, allowUpperBoundOverride bool) semver.Constraint {
	if !repo.IsPlatformPackage(req) {
		return c
	}
	if !allowUpperBoundOverride || !f.upper.MatchString(req) {
		return c
	}
	if f.ignore.MatchString(req) {
		return &semver.MatchAll{}
	}
	end, ok := semver.LastIntervalEnd(c)
	if ok && end.String() != "< "+semver.PHPIntMax+".0.0.0" {
		return semver.NewMulti([]semver.Constraint{c, semver.NewConstraint(">=", end.Version())}, false)
	}
	return c
}

// filterConstraint applies an IgnoreList filter if one is active.
func filterConstraint(f PlatformFilter, name string, c semver.Constraint, allowUpper bool) semver.Constraint {
	if il, ok := f.(*IgnoreList); ok {
		return il.FilterConstraint(name, c, allowUpper)
	}
	return c
}
