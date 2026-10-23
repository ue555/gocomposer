// Package solver ports Composer\DependencyResolver at
// f83588392df43d4fd027d84d26de459946ef4e80: pool building and optimisation,
// rule generation, the CDCL SAT solver, policy, problem reporting and
// transactions.
package solver

import (
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// Pool is Composer\DependencyResolver\Pool.
type Pool struct {
	packages                 []*pkg.Package
	byName                   map[string][]*pkg.Package
	providerCache            map[string]map[string][]*pkg.Package
	unacceptable             []*pkg.Package
	removedVersions          map[string]*orderedVersions
	removedVersionsByPackage map[*pkg.Package]*orderedVersions
	policyRemoved            *policy.PolicyRemovedVersions
}

// WithPackages is `new Pool($packages, $pool->getUnacceptableFixedOrLockedPackages(),
// $pool->getAllRemovedVersions(), $pool->getAllRemovedVersionsByPackage(),
// $policyRemovedVersions)` as used by the policy pool filters.
func (p *Pool) WithPackages(packages []*pkg.Package, policyRemoved *policy.PolicyRemovedVersions) *Pool {
	np := NewPool(packages, p.unacceptable)
	np.removedVersions = p.removedVersions
	np.removedVersionsByPackage = p.removedVersionsByPackage
	np.policyRemoved = policyRemoved
	return np
}

// AllPolicyRemovedVersions is Pool::getAllPolicyRemovedVersions.
func (p *Pool) AllPolicyRemovedVersions() *policy.PolicyRemovedVersions { return p.policyRemoved }

// PolicyRemovalReason is Pool::getPolicyRemovalReason.
func (p *Pool) PolicyRemovalReason(name, version string) *policy.PolicyRemovalReason {
	if p.policyRemoved == nil {
		return nil
	}
	return p.policyRemoved.Reason(name, version)
}

// orderedVersions is an ordered PHP array version => prettyVersion.
type orderedVersions struct {
	keys []string
	vals map[string]string
}

func newOrderedVersions() *orderedVersions { return &orderedVersions{vals: map[string]string{}} }

func (o *orderedVersions) set(k, v string) {
	if _, ok := o.vals[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.vals[k] = v
}

func (o *orderedVersions) clone() *orderedVersions {
	c := newOrderedVersions()
	for _, k := range o.keys {
		c.set(k, o.vals[k])
	}
	return c
}

// NewPool assigns ids 1..n to the packages like Pool::setPackages.
func NewPool(packages []*pkg.Package, unacceptable []*pkg.Package) *Pool {
	p := &Pool{byName: map[string][]*pkg.Package{}, providerCache: map[string]map[string][]*pkg.Package{}, unacceptable: unacceptable,
		removedVersions: map[string]*orderedVersions{}, removedVersionsByPackage: map[*pkg.Package]*orderedVersions{}}
	for i, pk := range packages {
		p.packages = append(p.packages, pk)
		pk.ID = i + 1
		for _, n := range pk.Names(true) {
			p.byName[n] = append(p.byName[n], pk)
		}
	}
	return p
}

func (p *Pool) Packages() []*pkg.Package        { return p.packages }
func (p *Pool) Count() int                      { return len(p.packages) }
func (p *Pool) PackageByID(id int) *pkg.Package { return p.packages[id-1] }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p *Pool) LiteralToPackage(lit int) *pkg.Package { return p.PackageByID(abs(lit)) }

// LiteralToPrettyString is Pool::literalToPrettyString.
func (p *Pool) LiteralToPrettyString(lit int, installedMap map[int]*pkg.Package) string {
	pk := p.LiteralToPackage(lit)
	var prefix string
	if _, ok := installedMap[pk.ID]; ok {
		if lit > 0 {
			prefix = "keep"
		} else {
			prefix = "remove"
		}
	} else if lit > 0 {
		prefix = "install"
	} else {
		prefix = "don't install"
	}
	return prefix + " " + pk.PrettyString()
}

// WhatProvides is Pool::whatProvides with its per-constraint cache.
func (p *Pool) WhatProvides(name string, c semver.Constraint) []*pkg.Package {
	key := ""
	if c != nil {
		key = c.String()
	}
	if m, ok := p.providerCache[name]; ok {
		if r, ok := m[key]; ok {
			return r
		}
	} else {
		p.providerCache[name] = map[string][]*pkg.Package{}
	}
	var matches []*pkg.Package
	for _, cand := range p.byName[name] {
		if p.Match(cand, name, c) {
			matches = append(matches, cand)
		}
	}
	p.providerCache[name][key] = matches
	return matches
}

// Match is Pool::match.
func (p *Pool) Match(cand *pkg.Package, name string, c semver.Constraint) bool {
	if cand.Name == name {
		return c == nil || semver.MatchVersion(c, cand.Version())
	}
	for _, l := range cand.Provides() {
		if l.Target == name && (c == nil || c.Matches(l.Constraint)) {
			return true
		}
	}
	for _, l := range cand.Replaces() {
		if l.Target == name && (c == nil || c.Matches(l.Constraint)) {
			return true
		}
	}
	return false
}

func (p *Pool) IsUnacceptableFixedOrLockedPackage(pk *pkg.Package) bool {
	for _, u := range p.unacceptable {
		if u == pk {
			return true
		}
	}
	return false
}

func (p *Pool) UnacceptableFixedOrLockedPackages() []*pkg.Package { return p.unacceptable }

// RemovedVersions is Pool::getRemovedVersions (version => pretty).
func (p *Pool) RemovedVersions(name string, c semver.Constraint) *orderedVersions {
	out := newOrderedVersions()
	rv, ok := p.removedVersions[name]
	if !ok {
		return out
	}
	for _, v := range rv.keys {
		if c.Matches(semver.NewConstraint("==", v)) {
			out.set(v, rv.vals[v])
		}
	}
	return out
}

func (p *Pool) RemovedVersionsByPackage(pk *pkg.Package) *orderedVersions {
	if v, ok := p.removedVersionsByPackage[pk]; ok {
		return v
	}
	return newOrderedVersions()
}

// isPlatformRepo reports PlatformRepository instances.
func isPlatformRepo(pk *pkg.Package) bool {
	r, ok := pk.Repo.(repo.Repository)
	return ok && r.Kind() == repo.KindPlatform
}
