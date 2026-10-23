package solver

import (
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// RepositorySet flags.
const (
	AllowUnacceptableStabilities = 1
	AllowShadowedRepositories    = 2
)

// RootAliasEntry is one alias in RepositorySet::$rootAliases.
type RootAliasEntry struct {
	Alias           string
	AliasNormalized string
}

// RepositorySet is Composer\Repository\RepositorySet.
type RepositorySet struct {
	rootAliases          map[string]map[string]RootAliasEntry
	rootReferences       map[string]string
	repositories         []repo.Repository
	acceptable           map[string]int
	stabilityFlags       map[string]int
	rootRequires         []RequireEntry
	temporaryConstraints map[string]semver.Constraint
	locked               bool
	allowInstalled       bool
}

func NewRepositorySet(minimumStability string, stabilityFlags map[string]int, rootAliases []pkg.RootAlias, rootReferences map[string]string, rootRequires []RequireEntry, temporary map[string]semver.Constraint) *RepositorySet {
	rs := &RepositorySet{rootAliases: map[string]map[string]RootAliasEntry{}, rootReferences: rootReferences, acceptable: map[string]int{},
		stabilityFlags: stabilityFlags, temporaryConstraints: temporary}
	for _, a := range rootAliases {
		if rs.rootAliases[a.Package] == nil {
			rs.rootAliases[a.Package] = map[string]RootAliasEntry{}
		}
		rs.rootAliases[a.Package][a.Version] = RootAliasEntry{a.Alias, a.AliasNormalized}
	}
	for name, v := range pkg.Stabilities {
		if v <= pkg.Stabilities[minimumStability] {
			rs.acceptable[name] = v
		}
	}
	for _, r := range rootRequires {
		if !repo.IsPlatformPackage(r.Name) {
			rs.rootRequires = append(rs.rootRequires, r)
		}
	}
	if rs.stabilityFlags == nil {
		rs.stabilityFlags = map[string]int{}
	}
	return rs
}

func (rs *RepositorySet) AllowInstalledRepositories() { rs.allowInstalled = true }

func (rs *RepositorySet) RootRequire(name string) (semver.Constraint, bool) {
	for _, r := range rs.rootRequires {
		if r.Name == name {
			return r.Constraint, true
		}
	}
	return nil, false
}

func (rs *RepositorySet) AddRepository(r repo.Repository) {
	if rs.locked {
		panic("Pool has already been created from this repository set, it cannot be modified anymore.")
	}
	rs.repositories = append(rs.repositories, r)
}

func (rs *RepositorySet) Repositories() []repo.Repository { return rs.repositories }

var allStabilities = map[string]int{"stable": 0, "RC": 5, "beta": 10, "alpha": 15, "dev": 20}

// FindPackages is RepositorySet::findPackages.
func (rs *RepositorySet) FindPackages(name string, c semver.Constraint, flags int) []*pkg.Package {
	out, _ := rs.FindPackagesE(name, c, flags)
	return out
}

// FindPackagesE is RepositorySet::findPackages propagating repository
// errors (e.g. TransportException) like PHP does.
func (rs *RepositorySet) FindPackagesE(name string, c semver.Constraint, flags int) ([]*pkg.Package, error) {
	ignoreStability := flags&AllowUnacceptableStabilities != 0
	loadFromAll := flags&AllowShadowedRepositories != 0
	var candidates []*pkg.Package
	if loadFromAll {
		for _, r := range rs.repositories {
			found, err := r.FindPackages(name, c)
			if err != nil {
				return nil, err
			}
			candidates = append(candidates, found...)
		}
	} else {
	repos:
		for _, r := range rs.repositories {
			acc, fl := rs.acceptable, rs.stabilityFlags
			if ignoreStability {
				acc, fl = allStabilities, map[string]int{}
			}
			res, err := r.LoadPackages([]repo.NameConstraint{{Name: name, Constraint: c}}, acc, fl, nil)
			if err != nil {
				return nil, err
			}
			candidates = append(candidates, res.Packages...)
			for _, n := range res.NamesFound {
				if n == name {
					break repos
				}
			}
		}
	}
	if ignoreStability || !loadFromAll {
		return candidates, nil
	}
	var out []*pkg.Package
	for _, p := range candidates {
		if repo.IsPackageAcceptable(rs.acceptable, rs.stabilityFlags, p.Names(true), p.Stability()) {
			out = append(out, p)
		}
	}
	return out, nil
}

// Providers is RepositorySet::getProviders.
func (rs *RepositorySet) Providers(name string) []repo.ProviderInfo {
	var out []repo.ProviderInfo
	idx := map[string]int{}
	for _, r := range rs.repositories {
		ps, _ := r.Providers(name)
		for _, p := range ps {
			if i, ok := idx[p.Name]; ok {
				out[i] = p
				continue
			}
			idx[p.Name] = len(out)
			out = append(out, p)
		}
	}
	return out
}

// CreatePool is RepositorySet::createPool.
func (rs *RepositorySet) CreatePool(req *Request, b *PoolBuilder) (*Pool, error) {
	for _, r := range rs.repositories {
		if r.Kind() == repo.KindInstalled && !rs.allowInstalled {
			panic("The pool can not accept packages from an installed repository")
		}
	}
	rs.locked = true
	b.acceptable = rs.acceptable
	b.stabilityFlags = rs.stabilityFlags
	b.rootAliases = rs.rootAliases
	b.rootReferences = rs.rootReferences
	b.temporaryConstraints = rs.temporaryConstraints
	return b.BuildPool(rs.repositories, req)
}

// CreatePoolWithAllPackages is RepositorySet::createPoolWithAllPackages.
func (rs *RepositorySet) CreatePoolWithAllPackages() (*Pool, error) {
	rs.locked = true
	var packages []*pkg.Package
	for _, r := range rs.repositories {
		ps, err := r.Packages()
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			packages = append(packages, p)
			if a, ok := rs.rootAliases[p.Name][p.Version()]; ok {
				base := p.Base()
				alias := pkg.NewAlias(base, a.AliasNormalized, a.Alias)
				alias.RootPackageAlias = true
				packages = append(packages, alias)
			}
		}
	}
	return NewPool(packages, nil), nil
}

// IsPackageAcceptable is RepositorySet::isPackageAcceptable.
func (rs *RepositorySet) IsPackageAcceptable(names []string, stability string) bool {
	return repo.IsPackageAcceptable(rs.acceptable, rs.stabilityFlags, names, stability)
}
