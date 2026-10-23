// Package repo ports Composer's repositories (Composer\Repository\*).
package repo

import (
	"regexp"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// Kind distinguishes repository classes where Composer uses instanceof checks.
type Kind int

const (
	KindArray Kind = iota
	KindPlatform
	KindLock
	KindRoot
	KindInstalled
	KindComposer
	KindPath
	KindVcs
	KindPackage
	KindArtifact
)

// NameConstraint is one entry of the ordered $packageNameMap.
type NameConstraint struct {
	Name       string
	Constraint semver.Constraint // nil means any
}

// LoadResult is the array returned by loadPackages().
type LoadResult struct {
	NamesFound []string
	Packages   []*pkg.Package
}

// ProviderInfo is one getProviders() entry.
type ProviderInfo struct {
	Name        string
	Description string
	Type        string
}

// Repository is RepositoryInterface.
type Repository interface {
	RepoName() string
	Kind() Kind
	LoadPackages(names []NameConstraint, acceptable map[string]int, flags map[string]int, alreadyLoaded map[string]map[string]bool) (LoadResult, error)
	Packages() ([]*pkg.Package, error)
	FindPackages(name string, c semver.Constraint) ([]*pkg.Package, error)
	Providers(name string) ([]ProviderInfo, error)
}

// IsPackageAcceptable ports StabilityFilter::isPackageAcceptable.
func IsPackageAcceptable(acceptable map[string]int, flags map[string]int, names []string, stability string) bool {
	for _, name := range names {
		if f, ok := flags[name]; ok {
			if pkg.Stabilities[stability] <= f {
				return true
			}
		} else if _, ok := acceptable[stability]; ok {
			return true
		}
	}
	return false
}

var platformPackageRe = regexp.MustCompile(`(?i)^(?:php(?:-64bit|-ipv6|-zts|-debug)?|hhvm|(?:ext|lib)-[a-z0-9](?:[_.-]?[a-z0-9]+)*|composer(?:-(?:plugin|runtime)-api)?)$`)

// IsPlatformPackage ports PlatformRepository::isPlatformPackage.
func IsPlatformPackage(name string) bool { return platformPackageRe.MatchString(name) }

// Array is ArrayRepository; other repositories embed it.
type Array struct {
	self     Repository
	packages []*pkg.Package
	kind     Kind
	name     string
	init     func() error
	inited   bool
}

func NewArray(packages ...*pkg.Package) *Array {
	a := &Array{kind: KindArray}
	for _, p := range packages {
		a.AddPackage(p)
	}
	return a
}

// NewLockArray is LockArrayRepository.
func NewLockArray() *Array { return &Array{kind: KindLock, inited: true} }

// NewPlatformArray is the storage of PlatformRepository.
func NewPlatformArray() *Array {
	return &Array{kind: KindPlatform, inited: true, name: "platform repo"}
}

// NewNamedArray creates an array repository with a fixed repo name and kind.
func NewNamedArray(kind Kind, name string) *Array {
	return &Array{kind: kind, inited: true, name: name}
}

// SetSelf makes packages added to the array point at the wrapping repository.
func (a *Array) SetSelf(r Repository) { a.self = r }

// SetName overrides getRepoName().
func (a *Array) SetName(n string) { a.name = n }

// NewInstalledArray is InstalledArrayRepository (used for mocks and dry runs).
func NewInstalledArray() *Array { return &Array{kind: KindInstalled, inited: true} }

// NewRootPackageRepository is RootPackageRepository.
func NewRootPackageRepository(root *pkg.Package) *Array {
	a := &Array{kind: KindRoot, inited: true}
	a.AddPackage(root)
	return a
}

func (a *Array) Kind() Kind { return a.kind }

func (a *Array) RepoName() string {
	if a.name != "" {
		return a.name
	}
	switch a.kind {
	case KindLock:
		return "lock repo"
	case KindRoot:
		return "root package repo"
	case KindInstalled:
		return "installed repo"
	case KindPackage:
		return PackageRepoName(a)
	}
	n := len(a.packages)
	s := ""
	if n > 1 {
		s = "s"
	}
	return "array repo (defining " + itoa(n) + " package" + s + ")"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

func (a *Array) ensure() error {
	if a.inited {
		return nil
	}
	a.inited = true
	if a.init != nil {
		return a.init()
	}
	return nil
}

// AddPackage is ArrayRepository::addPackage.
func (a *Array) AddPackage(p *pkg.Package) {
	if !a.inited && a.init == nil {
		a.inited = true
	}
	if a.self != nil {
		p.Repo = a.self
	} else {
		p.Repo = a
	}
	a.packages = append(a.packages, p)
	if p.AliasOf != nil && p.AliasOf.Repo == nil {
		a.AddPackage(p.AliasOf)
	}
}

func (a *Array) RemovePackage(p *pkg.Package) {
	id := p.UniqueName()
	for i, rp := range a.packages {
		if rp.UniqueName() == id {
			a.packages = append(a.packages[:i], a.packages[i+1:]...)
			return
		}
	}
}

func (a *Array) Packages() ([]*pkg.Package, error) {
	if err := a.ensure(); err != nil {
		return nil, err
	}
	return a.packages, nil
}

// MustPackages is for in-memory repositories that cannot fail.
func (a *Array) MustPackages() []*pkg.Package {
	p, _ := a.Packages()
	return p
}

func (a *Array) HasPackage(p *pkg.Package) bool {
	pkgs, _ := a.Packages()
	for _, rp := range pkgs {
		if rp.UniqueName() == p.UniqueName() {
			return true
		}
	}
	return false
}

// LoadPackages is ArrayRepository::loadPackages.
func (a *Array) LoadPackages(names []NameConstraint, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) (LoadResult, error) {
	pkgs, err := a.Packages()
	if err != nil {
		return LoadResult{}, err
	}
	return LoadFromList(pkgs, names, acceptable, flags, alreadyLoaded), nil
}

// LoadFromList implements the shared ArrayRepository::loadPackages logic.
func LoadFromList(pkgs []*pkg.Package, names []NameConstraint, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) LoadResult {
	nameMap := map[string]semver.Constraint{}
	present := map[string]bool{}
	for _, nc := range names {
		nameMap[nc.Name] = nc.Constraint
		present[nc.Name] = true
	}
	var result []*pkg.Package
	inResult := map[*pkg.Package]bool{}
	add := func(p *pkg.Package) {
		if !inResult[p] {
			inResult[p] = true
			result = append(result, p)
		}
	}
	var namesFound []string
	found := map[string]bool{}
	for _, p := range pkgs {
		if !present[p.Name] {
			continue
		}
		c := nameMap[p.Name]
		if (c == nil || c.Matches(semver.NewConstraint("==", p.Version()))) &&
			IsPackageAcceptable(acceptable, flags, p.Names(true), p.Stability()) &&
			!alreadyLoaded[p.Name][p.Version()] {
			add(p)
			if p.AliasOf != nil && !inResult[p.AliasOf] {
				add(p.AliasOf)
			}
		}
		if !found[p.Name] {
			found[p.Name] = true
			namesFound = append(namesFound, p.Name)
		}
	}
	for _, p := range pkgs {
		if p.AliasOf != nil && inResult[p.AliasOf] {
			add(p)
		}
	}
	return LoadResult{NamesFound: namesFound, Packages: result}
}

// FindPackages is ArrayRepository::findPackages.
func (a *Array) FindPackages(name string, c semver.Constraint) ([]*pkg.Package, error) {
	pkgs, err := a.Packages()
	if err != nil {
		return nil, err
	}
	return FindInList(pkgs, name, c), nil
}

func FindInList(pkgs []*pkg.Package, name string, c semver.Constraint) []*pkg.Package {
	name = strings.ToLower(name)
	var out []*pkg.Package
	for _, p := range pkgs {
		if p.Name == name && (c == nil || c.Matches(semver.NewConstraint("==", p.Version()))) {
			out = append(out, p)
		}
	}
	return out
}

// FindPackage is ArrayRepository::findPackage.
func FindPackage(r Repository, name string, c semver.Constraint) *pkg.Package {
	found, _ := r.FindPackages(name, c)
	if len(found) > 0 {
		return found[0]
	}
	return nil
}

// Providers is ArrayRepository::getProviders.
func (a *Array) Providers(name string) ([]ProviderInfo, error) {
	pkgs, err := a.Packages()
	if err != nil {
		return nil, err
	}
	return ProvidersInList(pkgs, name), nil
}

func ProvidersInList(pkgs []*pkg.Package, name string) []ProviderInfo {
	var out []ProviderInfo
	seen := map[string]bool{}
	for _, c := range pkgs {
		if seen[c.Name] {
			continue
		}
		for _, l := range c.Provides() {
			if l.Target == name {
				seen[c.Name] = true
				out = append(out, ProviderInfo{Name: c.Name, Description: c.Description(), Type: c.Type()})
				break
			}
		}
	}
	return out
}
