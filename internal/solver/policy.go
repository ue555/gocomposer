package solver

import (
	"os"
	"sort"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// DefaultPolicy is Composer\DependencyResolver\DefaultPolicy.
type DefaultPolicy struct {
	PreferStable            bool
	PreferLowest            bool
	preferDevOverPrerelease bool
	PreferredVersions       map[string]string
}

func NewDefaultPolicy(preferStable, preferLowest bool, preferred map[string]string) *DefaultPolicy {
	v := os.Getenv("COMPOSER_PREFER_DEV_OVER_PRERELEASE")
	return &DefaultPolicy{PreferStable: preferStable, PreferLowest: preferLowest, PreferredVersions: preferred,
		preferDevOverPrerelease: v != "" && v != "0"}
}

// VersionCompare is DefaultPolicy::versionCompare.
func (p *DefaultPolicy) VersionCompare(a, b *pkg.Package, op string) bool {
	if p.PreferStable {
		stabA, stabB := a.Stability(), b.Stability()
		if stabA != stabB {
			if p.PreferLowest && p.preferDevOverPrerelease && stabA != "stable" && stabB != "stable" {
				if stabA == "dev" {
					stabA = "stable"
				}
				if stabB == "dev" {
					stabB = "stable"
				}
			}
			return pkg.Stabilities[stabA] < pkg.Stabilities[stabB]
		}
	}
	if (a.IsDev() && strings.HasPrefix(a.Version(), "dev-")) || (b.IsDev() && strings.HasPrefix(b.Version(), "dev-")) {
		return semver.NewConstraint(op, b.Version()).MatchSpecific(semver.NewConstraint("==", a.Version()), true)
	}
	return semver.MatchVersion(semver.NewConstraint(op, b.Version()), a.Version())
}

// SelectPreferredPackages is DefaultPolicy::selectPreferredPackages.
func (p *DefaultPolicy) SelectPreferredPackages(pool *Pool, literals []int, requiredPackage string) []int {
	lits := append([]int(nil), literals...)
	sort.Ints(lits)
	var names []string
	groups := map[string][]int{}
	for _, l := range lits {
		n := pool.LiteralToPackage(l).Name
		if _, ok := groups[n]; !ok {
			names = append(names, n)
		}
		groups[n] = append(groups[n], l)
	}
	var selected []int
	for _, n := range names {
		g := groups[n]
		sort.SliceStable(g, func(i, j int) bool {
			return p.CompareByPriority(pool, pool.LiteralToPackage(g[i]), pool.LiteralToPackage(g[j]), requiredPackage, true) < 0
		})
		g = p.pruneToBestVersion(pool, g)
		g = p.pruneRemoteAliases(pool, g)
		selected = append(selected, g...)
	}
	sort.SliceStable(selected, func(i, j int) bool {
		return p.CompareByPriority(pool, pool.LiteralToPackage(selected[i]), pool.LiteralToPackage(selected[j]), requiredPackage, false) < 0
	})
	return selected
}

// CompareByPriority is DefaultPolicy::compareByPriority.
func (p *DefaultPolicy) CompareByPriority(pool *Pool, a, b *pkg.Package, requiredPackage string, ignoreReplace bool) int {
	if a.Name == b.Name {
		aAliased := a.AliasOf != nil
		bAliased := b.AliasOf != nil
		if aAliased && !bAliased {
			return -1
		}
		if !aAliased && bAliased {
			return 1
		}
	}
	if !ignoreReplace {
		if replaces(a, b) {
			return 1
		}
		if replaces(b, a) {
			return -1
		}
		if requiredPackage != "" {
			if pos := strings.IndexByte(requiredPackage, '/'); pos >= 0 {
				vendor := requiredPackage[:pos]
				aSame := strings.HasPrefix(a.Name, vendor)
				bSame := strings.HasPrefix(b.Name, vendor)
				if aSame != bSame {
					if aSame {
						return -1
					}
					return 1
				}
			}
		}
	}
	if a.ID == b.ID {
		return 0
	}
	if a.ID < b.ID {
		return -1
	}
	return 1
}

func replaces(source, target *pkg.Package) bool {
	for _, l := range source.Replaces() {
		if l.Target == target.Name {
			return true
		}
	}
	return false
}

func (p *DefaultPolicy) pruneToBestVersion(pool *Pool, literals []int) []int {
	if p.PreferredVersions != nil {
		name := pool.LiteralToPackage(literals[0]).Name
		if pv, ok := p.PreferredVersions[name]; ok {
			var best []int
			for _, l := range literals {
				if pool.LiteralToPackage(l).Version() == pv {
					best = append(best, l)
				}
			}
			if len(best) > 0 {
				return best
			}
		}
	}
	op := ">"
	if p.PreferLowest {
		op = "<"
	}
	best := []int{literals[0]}
	bestPkg := pool.LiteralToPackage(literals[0])
	for _, l := range literals[1:] {
		pk := pool.LiteralToPackage(l)
		if p.VersionCompare(pk, bestPkg, op) {
			bestPkg = pk
			best = []int{l}
		} else if p.VersionCompare(pk, bestPkg, "==") {
			best = append(best, l)
		}
	}
	return best
}

func (p *DefaultPolicy) pruneRemoteAliases(pool *Pool, literals []int) []int {
	hasLocal := false
	for _, l := range literals {
		pk := pool.LiteralToPackage(l)
		if pk.AliasOf != nil && pk.RootPackageAlias {
			hasLocal = true
			break
		}
	}
	if !hasLocal {
		return literals
	}
	var out []int
	for _, l := range literals {
		pk := pool.LiteralToPackage(l)
		if pk.AliasOf != nil && pk.RootPackageAlias {
			out = append(out, l)
		}
	}
	return out
}
