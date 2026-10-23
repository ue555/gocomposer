package repo

// Port of Composer\Repository\InstalledRepository (and the CompositeRepository
// behaviour it inherits) from composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

import (
	"strconv"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// Installed is InstalledRepository: a composite of the local/lock, root,
// and platform repositories.
type Installed struct {
	Repos []Repository
}

// NewInstalled creates the composite.
func NewInstalled(repos ...Repository) *Installed { return &Installed{Repos: repos} }

// AddRepository appends a repository.
func (r *Installed) AddRepository(rp Repository) { r.Repos = append(r.Repos, rp) }

// RepoName is InstalledRepository::getRepoName.
func (r *Installed) RepoName() string {
	var names []string
	for _, rp := range r.Repos {
		names = append(names, rp.RepoName())
	}
	return "installed repo (" + strings.Join(names, ", ") + ")"
}

// Packages is CompositeRepository::getPackages.
func (r *Installed) Packages() []*pkg.Package {
	var out []*pkg.Package
	for _, rp := range r.Repos {
		ps, _ := rp.Packages()
		out = append(out, ps...)
	}
	return out
}

// FindPackages is CompositeRepository::findPackages.
func (r *Installed) FindPackages(name string, c semver.Constraint) []*pkg.Package {
	var out []*pkg.Package
	for _, rp := range r.Repos {
		ps, _ := rp.FindPackages(name, c)
		out = append(out, ps...)
	}
	return out
}

// FindPackage is CompositeRepository::findPackage.
func (r *Installed) FindPackage(name string, c semver.Constraint) *pkg.Package {
	for _, rp := range r.Repos {
		if p := FindPackage(rp, name, c); p != nil {
			return p
		}
	}
	return nil
}

// FindPackagesWithReplacersAndProviders is the InstalledRepository method.
func (r *Installed) FindPackagesWithReplacersAndProviders(name string, c semver.Constraint) []*pkg.Package {
	name = strings.ToLower(name)
	var matches []*pkg.Package
	for _, rp := range r.Repos {
		ps, _ := rp.Packages()
		for _, cand := range ps {
			if cand.Name == name {
				if c == nil || c.Matches(semver.NewConstraint("==", cand.Version())) {
					matches = append(matches, cand)
				}
				continue
			}
			for _, l := range append(append(pkg.Links(nil), cand.Provides()...), cand.Replaces()...) {
				if l.Target == name && (c == nil || c.Matches(l.Constraint)) {
					matches = append(matches, cand)
					break
				}
			}
		}
	}
	return matches
}

// Dependent is one getDependents() result: [package, link, children|false].
type Dependent struct {
	Package *pkg.Package
	Link    *pkg.Link
	// Children are the recursive dependents; nil with Cut=true is PHP's false
	// (circular or not expanded).
	Children []Dependent
	Cut      bool
}

// linkKey is the PHP array key of a link in its links array.
func linkKeys(ls pkg.Links) []string {
	keys := make([]string, len(ls))
	idx := 0
	for i, l := range ls {
		if l.Indexed {
			keys[i] = "#" + strconv.Itoa(idx)
			idx++
		} else {
			keys[i] = l.Target
		}
	}
	return keys
}

// unionLinks is PHP's `$a += $b` on link arrays.
func unionLinks(a, b pkg.Links) pkg.Links {
	have := map[string]bool{}
	for _, k := range linkKeys(a) {
		have[k] = true
	}
	out := append(pkg.Links(nil), a...)
	bk := linkKeys(b)
	for i, l := range b {
		if !have[bk[i]] {
			have[bk[i]] = true
			out = append(out, l)
		}
	}
	return out
}

func inList(s string, list []string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// GetDependents ports InstalledRepository::getDependents.
func (r *Installed) GetDependents(needle []string, constraint semver.Constraint, invert, recurse bool, packagesFound []string) []Dependent {
	needles := make([]string, len(needle))
	for i, n := range needle {
		needles[i] = strings.ToLower(n)
	}
	var results []Dependent
	if packagesFound == nil {
		packagesFound = append([]string(nil), needles...)
	}
	all := r.Packages()
	var rootPackage *pkg.Package
	for _, p := range all {
		if p.IsRoot {
			rootPackage = p
			break
		}
	}
	for _, p := range all {
		links := p.Requires()
		packagesInTree := append([]string(nil), packagesFound...)
		if !invert {
			links = unionLinks(links, p.Replaces())
			for _, link := range p.Replaces() {
				for _, n := range append([]string(nil), needles...) {
					if link.Source == n {
						if constraint == nil || link.Constraint.Matches(constraint) {
							if inList(link.Target, packagesInTree) {
								results = append(results, Dependent{Package: p, Link: link, Cut: true})
								continue
							}
							packagesInTree = append(packagesInTree, link.Target)
							var deps []Dependent
							if recurse {
								deps = r.GetDependents([]string{link.Target}, nil, false, true, packagesInTree)
							}
							results = append(results, Dependent{Package: p, Link: link, Children: deps})
							needles = append(needles, link.Target)
						}
					}
				}
			}
		}
		if p.IsRoot {
			links = unionLinks(links, p.DevRequires())
		}
		for _, link := range links {
			for _, n := range needles {
				if link.Target == n {
					if constraint == nil || link.Constraint.Matches(constraint) == !invert {
						if inList(link.Source, packagesInTree) {
							results = append(results, Dependent{Package: p, Link: link, Cut: true})
							continue
						}
						packagesInTree = append(packagesInTree, link.Source)
						var deps []Dependent
						if recurse {
							deps = r.GetDependents([]string{link.Source}, nil, false, true, packagesInTree)
						}
						results = append(results, Dependent{Package: p, Link: link, Children: deps})
					}
				}
			}
		}
		if invert && inList(p.Name, needles) {
			for _, link := range p.Conflicts() {
				for _, other := range r.FindPackages(link.Target, nil) {
					if link.Constraint.Matches(semver.NewConstraint("=", other.Version())) == invert {
						results = append(results, Dependent{Package: p, Link: link, Cut: true})
					}
				}
			}
		}
		for _, link := range p.Conflicts() {
			if inList(link.Target, needles) {
				for _, other := range r.FindPackages(link.Target, nil) {
					if link.Constraint.Matches(semver.NewConstraint("=", other.Version())) == invert {
						results = append(results, Dependent{Package: p, Link: link, Cut: true})
					}
				}
			}
		}
		if invert && constraint != nil && inList(p.Name, needles) && constraint.Matches(semver.NewConstraint("=", p.Version())) {
		requires:
			for _, link := range p.Requires() {
				if IsPlatformPackage(link.Target) {
					if r.FindPackage(link.Target, link.Constraint) != nil {
						continue
					}
					platformPkg := r.FindPackage(link.Target, &semver.MatchAll{})
					description := "but it is missing"
					if platformPkg != nil {
						description = "but " + platformPkg.PrettyVersion() + " is installed"
					}
					results = append(results, Dependent{Package: p, Link: &pkg.Link{Source: p.Name, Target: link.Target, Constraint: &semver.MatchAll{}, Description: pkg.TypeRequire, PrettyConstraint: link.PrettyConstraint + " " + description}, Cut: true})
					continue
				}
				for _, other := range all {
					if !inList(link.Target, other.Names(true)) {
						continue
					}
					var version semver.Constraint = semver.NewConstraint("=", other.Version())
					if link.Target != other.Name {
						for _, prov := range append(append(pkg.Links(nil), other.Replaces()...), other.Provides()...) {
							if link.Target == prov.Target {
								version = prov.Constraint
								break
							}
						}
					}
					if !link.Constraint.Matches(version) {
						if rootPackage != nil {
							for _, rootReq := range append(append(pkg.Links(nil), rootPackage.Requires()...), rootPackage.DevRequires()...) {
								if inList(rootReq.Target, other.Names(true)) && !rootReq.Constraint.Matches(link.Constraint) {
									results = append(results, Dependent{Package: p, Link: link, Cut: true}, Dependent{Package: rootPackage, Link: rootReq, Cut: true})
									continue requires
								}
							}
							results = append(results, Dependent{Package: p, Link: link, Cut: true},
								Dependent{Package: rootPackage, Link: &pkg.Link{Source: rootPackage.Name, Target: link.Target, Constraint: &semver.MatchAll{}, Description: pkg.TypeDoesNotRequire, PrettyConstraint: "but " + other.PrettyVersion() + " is installed"}, Cut: true})
						} else {
							results = append(results, Dependent{Package: p, Link: link, Cut: true})
						}
					}
					continue requires
				}
			}
		}
	}
	return results
}

// Repositories returns the composed repositories.
func (r *Installed) Repositories() []Repository { return r.Repos }

// HasPackage is CompositeRepository::hasPackage (by unique name).
func (r *Installed) HasPackage(p *pkg.Package) bool {
	id := p.UniqueName()
	for _, c := range r.Packages() {
		if c.UniqueName() == id {
			return true
		}
	}
	return false
}
