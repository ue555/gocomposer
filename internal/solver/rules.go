package solver

import (
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/semver"
)

// Rule reasons (Rule::RULE_*).
const (
	RuleRootRequire            = 2
	RuleFixed                  = 3
	RulePackageConflict        = 6
	RulePackageRequires        = 7
	RulePackageSameName        = 10
	RuleLearned                = 12
	RulePackageAlias           = 13
	RulePackageInverseAlias    = 14
	RuleLockedFilterListRemove = 15
)

// Rule set types (RuleSet::TYPE_*).
const (
	TypePackage = 0
	TypeRequest = 1
	TypeLearned = 4
)

type ruleShape int

const (
	shapeGeneric ruleShape = iota
	shapeTwo
	shapeMulti
)

// RootRequireData is the reason data of RULE_ROOT_REQUIRE.
type RootRequireData struct {
	PackageName string
	Constraint  semver.Constraint
}

// Rule covers GenericRule, Rule2Literals and MultiConflictRule.
type Rule struct {
	literals []int
	shape    ruleShape
	Reason   int
	// ReasonData is *pkg.Link, RootRequireData, *pkg.Package, string or int
	ReasonData any
	Type       int
	disabled   bool
}

func newGenericRule(lits []int, reason int, data any) *Rule {
	l := append([]int(nil), lits...)
	sort.Ints(l)
	return &Rule{literals: l, shape: shapeGeneric, Reason: reason, ReasonData: data, Type: 255}
}

func newRule2(a, b int, reason int, data any) *Rule {
	if a > b {
		a, b = b, a
	}
	return &Rule{literals: []int{a, b}, shape: shapeTwo, Reason: reason, ReasonData: data, Type: 255}
}

func newMultiConflictRule(lits []int, reason int, data any) *Rule {
	l := append([]int(nil), lits...)
	sort.Ints(l)
	return &Rule{literals: l, shape: shapeMulti, Reason: reason, ReasonData: data, Type: 255}
}

func (r *Rule) Literals() []int  { return r.literals }
func (r *Rule) IsDisabled() bool { return r.disabled }
func (r *Rule) IsEnabled() bool  { return !r.disabled }
func (r *Rule) Disable() {
	if r.shape == shapeMulti {
		panic("Disabling multi conflict rules is not possible.")
	}
	r.disabled = true
}
func (r *Rule) Enable()           { r.disabled = false }
func (r *Rule) IsMulti() bool     { return r.shape == shapeMulti }
func (r *Rule) IsAssertion() bool { return r.shape == shapeGeneric && len(r.literals) == 1 }

func (r *Rule) key() string {
	var sb strings.Builder
	switch r.shape {
	case shapeGeneric:
		sb.WriteString("g:")
	case shapeTwo:
		sb.WriteString("2:")
	case shapeMulti:
		sb.WriteString("c:")
	}
	for i, l := range r.literals {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(l))
	}
	return sb.String()
}

// RequiredPackage is Rule::getRequiredPackage.
func (r *Rule) RequiredPackage() string {
	switch r.Reason {
	case RuleRootRequire:
		return r.ReasonData.(RootRequireData).PackageName
	case RuleFixed, RuleLockedFilterListRemove:
		return r.ReasonData.(*pkg.Package).Name
	case RulePackageRequires:
		return r.ReasonData.(*pkg.Link).Target
	}
	return ""
}

func (r *Rule) String() string {
	var sb strings.Builder
	if r.disabled {
		sb.WriteString("disabled(")
	} else {
		sb.WriteString("(")
	}
	if r.shape == shapeMulti {
		sb.WriteString("multi(")
	}
	for i, l := range r.literals {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(strconv.Itoa(l))
	}
	if r.shape == shapeMulti {
		sb.WriteString(")")
	}
	sb.WriteString(")")
	return sb.String()
}

// RuleSet is Composer\DependencyResolver\RuleSet.
type RuleSet struct {
	byType   map[int][]*Rule
	RuleByID []*Rule
	byKey    map[string]bool
}

func NewRuleSet() *RuleSet {
	return &RuleSet{byType: map[int][]*Rule{TypePackage: nil, TypeRequest: nil, TypeLearned: nil}, byKey: map[string]bool{}}
}

func (rs *RuleSet) Add(r *Rule, typ int) {
	k := r.key()
	if rs.byKey[k] {
		return
	}
	rs.byKey[k] = true
	rs.byType[typ] = append(rs.byType[typ], r)
	rs.RuleByID = append(rs.RuleByID, r)
	r.Type = typ
}

func (rs *RuleSet) Count() int { return len(rs.RuleByID) }

// Iterate returns rules ordered by type then insertion (RuleSetIterator).
func (rs *RuleSet) Iterate() []*Rule {
	var out []*Rule
	for _, t := range []int{TypePackage, TypeRequest, TypeLearned} {
		out = append(out, rs.byType[t]...)
	}
	return out
}

func (rs *RuleSet) OfType(t int) []*Rule { return rs.byType[t] }

// ruleSetGenerator is RuleSetGenerator.
type ruleSetGenerator struct {
	pool                 *Pool
	rules                *RuleSet
	addedMap             map[int]bool
	addedOrder           []*pkg.Package
	addedPackagesByNames map[string][]*pkg.Package
	addedNamesOrder      []string
}

func (g *ruleSetGenerator) addRule(typ int, r *Rule) {
	if r != nil {
		g.rules.Add(r, typ)
	}
}

func createRequireRule(p *pkg.Package, providers []*pkg.Package, reason int, data any) *Rule {
	lits := []int{-p.ID}
	for _, prov := range providers {
		if prov == p {
			return nil
		}
		lits = append(lits, prov.ID)
	}
	return newGenericRule(lits, reason, data)
}

func createInstallOneOfRule(packages []*pkg.Package, reason int, data any) *Rule {
	lits := make([]int, len(packages))
	for i, p := range packages {
		lits[i] = p.ID
	}
	return newGenericRule(lits, reason, data)
}

func createRule2Literals(issuer, provider *pkg.Package, reason int, data any) *Rule {
	if issuer == provider {
		return nil
	}
	return newRule2(-issuer.ID, -provider.ID, reason, data)
}

func createMultiConflictRule(packages []*pkg.Package, reason int, data any) *Rule {
	lits := make([]int, len(packages))
	for i, p := range packages {
		lits[i] = -p.ID
	}
	if len(lits) == 2 {
		return newRule2(lits[0], lits[1], reason, data)
	}
	return newMultiConflictRule(lits, reason, data)
}

func (g *ruleSetGenerator) addRulesForPackage(start *pkg.Package, filter PlatformFilter) {
	queue := []*pkg.Package{start}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if g.addedMap[p.ID] {
			continue
		}
		g.addedMap[p.ID] = true
		g.addedOrder = append(g.addedOrder, p)
		if p.AliasOf == nil {
			for _, name := range p.Names(false) {
				if _, ok := g.addedPackagesByNames[name]; !ok {
					g.addedNamesOrder = append(g.addedNamesOrder, name)
				}
				g.addedPackagesByNames[name] = append(g.addedPackagesByNames[name], p)
			}
		} else {
			queue = append(queue, p.AliasOf)
			g.addRule(TypePackage, createRequireRule(p, []*pkg.Package{p.AliasOf}, RulePackageAlias, p))
			g.addRule(TypePackage, createRequireRule(p.AliasOf, []*pkg.Package{p}, RulePackageInverseAlias, p.AliasOf))
			if !p.HasSelfVersionRequires() {
				continue
			}
		}
		for _, link := range p.Requires() {
			c := link.Constraint
			if filter.IsIgnored(link.Target) {
				continue
			}
			c = filterConstraint(filter, link.Target, c, true)
			possible := g.pool.WhatProvides(link.Target, c)
			g.addRule(TypePackage, createRequireRule(p, possible, RulePackageRequires, link))
			queue = append(queue, possible...)
		}
	}
}

func (g *ruleSetGenerator) addConflictRules(filter PlatformFilter) {
	for _, p := range g.addedOrder {
		for _, link := range p.Conflicts() {
			if _, ok := g.addedPackagesByNames[link.Target]; !ok {
				continue
			}
			if filter.IsIgnored(link.Target) {
				continue
			}
			c := filterConstraint(filter, link.Target, link.Constraint, false)
			for _, conflict := range g.pool.WhatProvides(link.Target, c) {
				if conflict.AliasOf == nil || conflict.Name == link.Target {
					g.addRule(TypePackage, createRule2Literals(p, conflict, RulePackageConflict, link))
				}
			}
		}
	}
	for _, name := range g.addedNamesOrder {
		packages := g.addedPackagesByNames[name]
		if len(packages) > 1 {
			g.addRule(TypePackage, createMultiConflictRule(packages, RulePackageSameName, name))
		}
	}
}

func (g *ruleSetGenerator) addRulesForRequest(req *Request, filter PlatformFilter) {
	for _, p := range req.FixedPackages() {
		// a locked package dropped by a policy filter list is reported by
		// Solver.checkForFilterListRemovedLockedPackages instead
		if req.IsLockedPackage(p) {
			if r := g.pool.PolicyRemovalReason(p.Name, p.Version()); r != nil && r.Type() == policy.RemovalFilterList {
				continue
			}
		}
		if p.ID == -1 {
			if g.pool.IsUnacceptableFixedOrLockedPackage(p) {
				continue
			}
			panic("Fixed package " + p.PrettyString() + " was not added to solver pool.")
		}
		g.addRulesForPackage(p, filter)
		g.addRule(TypeRequest, createInstallOneOfRule([]*pkg.Package{p}, RuleFixed, p))
	}
	for _, r := range req.Requires() {
		if filter.IsIgnored(r.Name) {
			continue
		}
		c := filterConstraint(filter, r.Name, r.Constraint, true)
		packages := g.pool.WhatProvides(r.Name, c)
		if len(packages) > 0 {
			for _, p := range packages {
				g.addRulesForPackage(p, filter)
			}
			g.addRule(TypeRequest, createInstallOneOfRule(packages, RuleRootRequire, RootRequireData{r.Name, c}))
		}
	}
}

func (g *ruleSetGenerator) addRulesForRootAliases(filter PlatformFilter) {
	for _, p := range g.pool.Packages() {
		if !g.addedMap[p.ID] && p.AliasOf != nil && (p.RootPackageAlias || g.addedMap[p.AliasOf.ID]) {
			g.addRulesForPackage(p, filter)
		}
	}
}

// GenerateRules is RuleSetGenerator::getRulesFor.
func GenerateRules(pool *Pool, req *Request, filter PlatformFilter) *RuleSet {
	if filter == nil {
		filter = IgnoreNothing{}
	}
	g := &ruleSetGenerator{pool: pool, rules: NewRuleSet(), addedMap: map[int]bool{}, addedPackagesByNames: map[string][]*pkg.Package{}}
	g.addRulesForRequest(req, filter)
	g.addRulesForRootAliases(filter)
	g.addConflictRules(filter)
	return g.rules
}
