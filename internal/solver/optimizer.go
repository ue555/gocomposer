package solver

import (
	"sort"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// PoolOptimizer is Composer\DependencyResolver\PoolOptimizer.
type PoolOptimizer struct {
	policy                        *DefaultPolicy
	irremovable                   map[int]bool
	requireConstraintsPerPackage  map[string]*constraintSet
	conflictConstraintsPerPackage map[string]*constraintSet
	packagesToRemove              map[int]bool
	aliasesPerPackage             map[int][]*pkg.Package
	removedVersionsByPackage      map[*pkg.Package]*orderedVersions
}

// constraintSet is an ordered string(constraint) => constraint map.
type constraintSet struct {
	keys []string
	vals map[string]semver.Constraint
}

func (c *constraintSet) add(con semver.Constraint) {
	k := con.String()
	if _, ok := c.vals[k]; !ok {
		c.keys = append(c.keys, k)
	}
	c.vals[k] = con
}

func NewPoolOptimizer(policy *DefaultPolicy) *PoolOptimizer { return &PoolOptimizer{policy: policy} }

// Optimize is PoolOptimizer::optimize.
func (o *PoolOptimizer) Optimize(req *Request, pool *Pool) *Pool {
	o.irremovable = map[int]bool{}
	o.requireConstraintsPerPackage = map[string]*constraintSet{}
	o.conflictConstraintsPerPackage = map[string]*constraintSet{}
	o.packagesToRemove = map[int]bool{}
	o.aliasesPerPackage = map[int][]*pkg.Package{}
	o.removedVersionsByPackage = map[*pkg.Package]*orderedVersions{}
	o.prepare(req, pool)
	o.optimizeByIdenticalDependencies(req, pool)
	o.optimizeImpossiblePackagesAway(req, pool)
	return o.applyRemovals(pool)
}

func (o *PoolOptimizer) prepare(req *Request, pool *Pool) {
	groups := map[string][]semver.Constraint{}
	var groupOrder []string
	for _, p := range req.FixedOrLockedPackages() {
		if _, ok := groups[p.Name]; !ok {
			groupOrder = append(groupOrder, p.Name)
		}
		groups[p.Name] = append(groups[p.Name], semver.NewConstraint("==", p.Version()))
	}
	for _, r := range req.Requires() {
		o.extractRequire(r.Name, r.Constraint)
	}
	for _, p := range pool.Packages() {
		for _, l := range p.Requires() {
			o.extractRequire(l.Target, l.Constraint)
		}
		for _, l := range p.Conflicts() {
			o.extractConflict(l.Target, l.Constraint)
		}
		if p.AliasOf != nil {
			o.aliasesPerPackage[p.AliasOf.ID] = append(o.aliasesPerPackage[p.AliasOf.ID], p)
		}
	}
	irremovableConstraints := map[string]semver.Constraint{}
	for _, name := range groupOrder {
		cs := groups[name]
		if len(cs) == 1 {
			irremovableConstraints[name] = cs[0]
		} else {
			irremovableConstraints[name] = semver.NewMulti(cs, false)
		}
	}
	for _, p := range pool.Packages() {
		c, ok := irremovableConstraints[p.Name]
		if !ok {
			continue
		}
		if semver.MatchVersion(c, p.Version()) {
			o.markIrremovable(p)
		}
	}
}

func (o *PoolOptimizer) markIrremovable(p *pkg.Package) {
	o.irremovable[p.ID] = true
	if p.AliasOf != nil {
		o.markIrremovable(p.AliasOf)
	}
	for _, a := range o.aliasesPerPackage[p.ID] {
		o.irremovable[a.ID] = true
	}
}

func (o *PoolOptimizer) applyRemovals(pool *Pool) *Pool {
	var packages []*pkg.Package
	removedVersions := map[string]*orderedVersions{}
	for _, p := range pool.Packages() {
		if !o.packagesToRemove[p.ID] {
			packages = append(packages, p)
		} else {
			if removedVersions[p.Name] == nil {
				removedVersions[p.Name] = newOrderedVersions()
			}
			removedVersions[p.Name].set(p.Version(), p.PrettyVersion())
		}
	}
	np := NewPool(packages, pool.UnacceptableFixedOrLockedPackages())
	np.removedVersions = removedVersions
	np.removedVersionsByPackage = o.removedVersionsByPackage
	np.policyRemoved = pool.policyRemoved
	return np
}

type hashGroup struct {
	keys []string
	vals map[string][]int
}

func (o *PoolOptimizer) optimizeByIdenticalDependencies(req *Request, pool *Pool) {
	// packageName => groupHash => dependencyHash => ids, all ordered
	type groupMap struct {
		keys []string
		vals map[string]*hashGroup
	}
	identical := map[string]*groupMap{}
	var nameOrder []string
	for _, p := range pool.Packages() {
		if o.irremovable[p.ID] {
			continue
		}
		o.markForRemoval(p.ID)
		depHash := calculateDependencyHash(p)
		var replaceParts []string
		version := p.Version()
		for _, l := range p.Replaces() {
			if semver.MatchVersion(l.Constraint, version) {
				replaceParts = append(replaceParts, "require:"+l.Constraint.String())
			}
		}
		for _, name := range p.Names(false) {
			reqs, ok := o.requireConstraintsPerPackage[name]
			if !ok {
				continue
			}
			invariantParts := append([]string(nil), replaceParts...)
			if confs, ok := o.conflictConstraintsPerPackage[name]; ok {
				for _, k := range confs.keys {
					if semver.MatchVersion(confs.vals[k], version) {
						invariantParts = append(invariantParts, "conflict:"+k)
					}
				}
			}
			invariantHash := strings.Join(invariantParts, "")
			addedToInvariant := false
			for _, k := range reqs.keys {
				var groupHash string
				if semver.MatchVersion(reqs.vals[k], version) {
					groupHash = "require:" + k + invariantHash
				} else {
					if addedToInvariant || invariantHash == "" {
						continue
					}
					addedToInvariant = true
					groupHash = invariantHash
				}
				gm, ok := identical[name]
				if !ok {
					gm = &groupMap{vals: map[string]*hashGroup{}}
					identical[name] = gm
					nameOrder = append(nameOrder, name)
				}
				hg, ok := gm.vals[groupHash]
				if !ok {
					hg = &hashGroup{vals: map[string][]int{}}
					gm.vals[groupHash] = hg
					gm.keys = append(gm.keys, groupHash)
				}
				if _, ok := hg.vals[depHash]; !ok {
					hg.keys = append(hg.keys, depHash)
				}
				hg.vals[depHash] = append(hg.vals[depHash], p.ID)
			}
		}
	}
	for _, name := range nameOrder {
		gm := identical[name]
		for _, gk := range gm.keys {
			hg := gm.vals[gk]
			for _, dk := range hg.keys {
				ids := hg.vals[dk]
				versions := o.groupVersions(pool, ids)
				if len(ids) == 1 {
					o.keepPackageInGroup(pool.PackageByID(ids[0]), name, versions)
					continue
				}
				for _, lit := range o.policy.SelectPreferredPackages(pool, ids, "") {
					o.keepPackageInGroup(pool.LiteralToPackage(lit), name, versions)
				}
			}
		}
	}
}

func calculateDependencyHash(p *pkg.Package) string {
	var sb strings.Builder
	sections := []struct {
		key   string
		links pkg.Links
	}{{"requires", p.Requires()}, {"conflicts", p.Conflicts()}, {"replaces", p.Replaces()}, {"provides", p.Provides()}}
	for _, sec := range sections {
		if len(sec.links) == 0 {
			continue
		}
		sb.WriteString(sec.key + ":")
		sub := map[string]string{}
		var keys []string
		for _, l := range sec.links {
			if _, ok := sub[l.Target]; !ok {
				keys = append(keys, l.Target)
			}
			sub[l.Target] = l.Constraint.String()
		}
		sort.SliceStable(keys, func(i, j int) bool { return phpjson.PHPCompareKeys(keys[i], keys[j]) < 0 })
		for _, k := range keys {
			sb.WriteString(k + "@" + sub[k])
		}
	}
	return sb.String()
}

func (o *PoolOptimizer) markForRemoval(id int) {
	if o.irremovable[id] {
		panic("Attempted removing a package which was previously marked irremovable")
	}
	o.packagesToRemove[id] = true
}

func (o *PoolOptimizer) keepPackageInGroup(p *pkg.Package, name string, versions *orderedVersions) {
	o.recordRemovedVersions(p, name, versions)
	if !o.packagesToRemove[p.ID] {
		return
	}
	delete(o.packagesToRemove, p.ID)
	if p.AliasOf != nil {
		delete(o.packagesToRemove, p.AliasOf.ID)
		o.recordRemovedVersions(p.AliasOf, name, versions)
		for _, a := range o.aliasesPerPackage[p.AliasOf.ID] {
			delete(o.packagesToRemove, a.ID)
			o.recordRemovedVersions(a, name, versions)
		}
		return
	}
	for _, a := range o.aliasesPerPackage[p.ID] {
		delete(o.packagesToRemove, a.ID)
		o.recordRemovedVersions(a, name, versions)
	}
}

func (o *PoolOptimizer) groupVersions(pool *Pool, ids []int) *orderedVersions {
	v := newOrderedVersions()
	for _, id := range ids {
		gp := pool.PackageByID(id)
		if gp.AliasOf != nil && gp.PrettyVersion() == semver.DefaultBranchAlias {
			gp = gp.AliasOf
		}
		v.set(gp.Version(), gp.PrettyVersion())
	}
	return v
}

func (o *PoolOptimizer) recordRemovedVersions(p *pkg.Package, name string, versions *orderedVersions) {
	if !containsStr(p.Names(false), name) {
		return
	}
	existing, ok := o.removedVersionsByPackage[p]
	if !ok {
		o.removedVersionsByPackage[p] = versions.clone()
		return
	}
	for _, k := range versions.keys {
		existing.set(k, versions.vals[k])
	}
}

func (o *PoolOptimizer) optimizeImpossiblePackagesAway(req *Request, pool *Pool) {
	if len(req.LockedPackages()) == 0 {
		return
	}
	type idxEntry struct {
		ids  []int
		pkgs map[int]*pkg.Package
	}
	index := map[string]*idxEntry{}
	for _, p := range pool.Packages() {
		id := p.ID
		if o.irremovable[id] {
			continue
		}
		if _, ok := o.aliasesPerPackage[id]; ok || p.AliasOf != nil {
			continue
		}
		if req.IsFixedPackage(p) || req.IsLockedPackage(p) {
			continue
		}
		e, ok := index[p.Name]
		if !ok {
			e = &idxEntry{pkgs: map[int]*pkg.Package{}}
			index[p.Name] = e
		}
		e.ids = append(e.ids, id)
		e.pkgs[id] = p
	}
	for _, lp := range req.LockedPackages() {
		unused := true
		for _, n := range lp.Names(false) {
			if _, ok := o.requireConstraintsPerPackage[n]; ok {
				unused = false
				break
			}
		}
		if unused {
			continue
		}
		for _, l := range lp.Requires() {
			e, ok := index[l.Target]
			if !ok {
				continue
			}
			var keep []int
			for _, id := range e.ids {
				if !semver.MatchVersion(l.Constraint, e.pkgs[id].Version()) {
					o.markForRemoval(id)
					delete(e.pkgs, id)
				} else {
					keep = append(keep, id)
				}
			}
			e.ids = keep
		}
	}
}

func (o *PoolOptimizer) extractRequire(name string, c semver.Constraint) {
	for _, e := range expandDisjunctive(c) {
		cs, ok := o.requireConstraintsPerPackage[name]
		if !ok {
			cs = &constraintSet{vals: map[string]semver.Constraint{}}
			o.requireConstraintsPerPackage[name] = cs
		}
		cs.add(e)
	}
}

func (o *PoolOptimizer) extractConflict(name string, c semver.Constraint) {
	for _, e := range expandDisjunctive(c) {
		cs, ok := o.conflictConstraintsPerPackage[name]
		if !ok {
			cs = &constraintSet{vals: map[string]semver.Constraint{}}
			o.conflictConstraintsPerPackage[name] = cs
		}
		cs.add(e)
	}
}

func expandDisjunctive(c semver.Constraint) []semver.Constraint {
	c = semver.CompactConstraint(c)
	if m, ok := c.(*semver.Multi); ok && m.IsDisjunctive() {
		return m.Constraints()
	}
	return []semver.Constraint{c}
}
