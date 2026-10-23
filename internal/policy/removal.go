// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/DependencyResolver/PolicyRemovalReason.php and the policy
// removal bookkeeping of src/Composer/DependencyResolver/Pool.php
// (getPolicyRemovedVersions, getPolicyRemovalReason, recordPolicyRemoval,
// getAllPolicyRemovedVersions).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"sort"
	"strings"
	"time"

	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// PolicyRemovalReason::* types.
const (
	RemovalAbandoned  = "abandoned"
	RemovalAdvisories = "advisories"
	RemovalFilterList = "filter-list"
	RemovalCooldown   = "cooldown"
)

// CooldownRemedy is PolicyRemovalReason::COOLDOWN_REMEDY.
const CooldownRemedy = `To exempt the package from the cooldown policy, add it to the "policy.cooldown.ignore" config, or run the command with --no-cooldown for a one-off bypass.`

// CooldownInfo is array{releaseDate, availableIn, source}.
type CooldownInfo struct {
	ReleaseDate string
	AvailableIn string
	Source      string
}

// PolicyRemovalReason is Composer\DependencyResolver\PolicyRemovalReason.
type PolicyRemovalReason struct {
	typ               string
	packageName       string
	prettyVersion     string
	advisories        []*Advisory
	filterListEntries []*FilterListEntry
	cooldown          *CooldownInfo
}

func RemovalReasonAbandoned(packageName, prettyVersion string) *PolicyRemovalReason {
	return &PolicyRemovalReason{typ: RemovalAbandoned, packageName: packageName, prettyVersion: prettyVersion}
}

func RemovalReasonAdvisories(packageName, prettyVersion string, advisories []*Advisory) *PolicyRemovalReason {
	return &PolicyRemovalReason{typ: RemovalAdvisories, packageName: packageName, prettyVersion: prettyVersion, advisories: append([]*Advisory(nil), advisories...)}
}

func RemovalReasonFilterList(packageName, prettyVersion string, entries []*FilterListEntry) *PolicyRemovalReason {
	return &PolicyRemovalReason{typ: RemovalFilterList, packageName: packageName, prettyVersion: prettyVersion, filterListEntries: append([]*FilterListEntry(nil), entries...)}
}

// RemovalReasonCooldown: source is "published-time" or "time".
func RemovalReasonCooldown(packageName, prettyVersion, releaseDate, availableIn, source string) *PolicyRemovalReason {
	return &PolicyRemovalReason{typ: RemovalCooldown, packageName: packageName, prettyVersion: prettyVersion, cooldown: &CooldownInfo{releaseDate, availableIn, source}}
}

// CombineRemovalReasons is PolicyRemovalReason::combine.
func CombineRemovalReasons(reasons []*PolicyRemovalReason) (*PolicyRemovalReason, error) {
	first := reasons[0]
	combined := &PolicyRemovalReason{typ: first.typ, packageName: first.packageName, prettyVersion: first.prettyVersion}
	seenAdv := map[string]bool{}
	seenEntry := map[*FilterListEntry]bool{}
	for _, r := range reasons {
		if r.typ != first.typ {
			return nil, newErr(LogicException, "Cannot combine reasons of different policies: "+first.typ+" and "+r.typ)
		}
		for _, a := range r.advisories {
			if !seenAdv[a.AdvisoryID] {
				seenAdv[a.AdvisoryID] = true
				combined.advisories = append(combined.advisories, a)
			}
		}
		for _, e := range r.filterListEntries {
			if !seenEntry[e] {
				seenEntry[e] = true
				combined.filterListEntries = append(combined.filterListEntries, e)
			}
		}
		if r.cooldown != nil && (combined.cooldown == nil || dateBefore(r.cooldown.ReleaseDate, combined.cooldown.ReleaseDate)) {
			combined.cooldown = r.cooldown
			combined.prettyVersion = r.prettyVersion
		}
	}
	return combined, nil
}

func dateBefore(a, b string) bool {
	ta, oka := parsePackageDate(a)
	tb, okb := parsePackageDate(b)
	if !oka || !okb {
		return a < b
	}
	return ta.Before(tb)
}

func (r *PolicyRemovalReason) Type() string                          { return r.typ }
func (r *PolicyRemovalReason) PackageName() string                   { return r.packageName }
func (r *PolicyRemovalReason) PrettyVersion() string                 { return r.prettyVersion }
func (r *PolicyRemovalReason) Advisories() []*Advisory               { return r.advisories }
func (r *PolicyRemovalReason) FilterListEntries() []*FilterListEntry { return r.filterListEntries }
func (r *PolicyRemovalReason) CooldownInfo() *CooldownInfo           { return r.cooldown }

// Verb is PolicyRemovalReason::getVerb.
func (r *PolicyRemovalReason) Verb(plural bool) string {
	if r.typ == RemovalFilterList {
		if plural {
			return "were"
		}
		return "was"
	}
	if plural {
		return "are"
	}
	return "is"
}

// Description is PolicyRemovalReason::getDescription.
func (r *PolicyRemovalReason) Description() string {
	switch r.typ {
	case RemovalAbandoned:
		return "abandoned"
	case RemovalAdvisories:
		return `affected by security advisories ("` + strings.Join(r.linkedAdvisoryIdentifiers(), `", "`) + `")`
	case RemovalFilterList:
		_, descs := r.filterListDescriptions()
		return strings.Join(descs, ", ")
	case RemovalCooldown:
		availableIn := ""
		if r.cooldown != nil {
			source := ""
			if r.cooldown.Source == "time" {
				source = ", based on the package-supplied time field as the repository provides no published-time"
			}
			availableIn = " (available in " + r.cooldown.AvailableIn + source + ")"
		}
		return `still in the cooldown period configured in "policy.cooldown"` + availableIn
	}
	return ""
}

// Remedy is PolicyRemovalReason::getRemedy.
func (r *PolicyRemovalReason) Remedy() string {
	switch r.typ {
	case RemovalAdvisories:
		return r.advisoryDetailsHint() + ` To ignore the advisories, add their IDs to the "policy.advisories.ignore-id" config or add the package to "policy.advisories.ignore".`
	case RemovalFilterList:
		names, _ := r.filterListDescriptions()
		quoted := make([]string, len(names))
		for i, n := range names {
			quoted[i] = `"policy.` + n + `.ignore"`
		}
		return "To ignore filters for this package, add the package to the " + strings.Join(quoted, " and ") + " config."
	case RemovalCooldown:
		return CooldownRemedy
	}
	return ""
}

// OffSwitch is PolicyRemovalReason::getOffSwitch.
func (r *PolicyRemovalReason) OffSwitch() string {
	if r.typ == RemovalFilterList {
		names, _ := r.filterListDescriptions()
		quoted := make([]string, len(names))
		for i, n := range names {
			quoted[i] = `"policy.` + n + `.block"`
		}
		return strings.Join(quoted, " and ")
	}
	return `"policy.` + r.typ + `.block"`
}

func (r *PolicyRemovalReason) linkedAdvisoryIdentifiers() []string {
	out := make([]string, len(r.advisories))
	for i, a := range r.advisories {
		switch {
		case a.Full && a.Link != nil && *a.Link != "":
			out[i] = "<href=" + EscapeFormat(*a.Link) + ">" + a.AdvisoryID + "</>"
		case strings.HasPrefix(a.AdvisoryID, "PKSA-"):
			out[i] = "<href=" + EscapeFormat("https://packagist.org/security-advisories/"+a.AdvisoryID) + ">" + a.AdvisoryID + "</>"
		default:
			out[i] = a.AdvisoryID
		}
	}
	return out
}

func (r *PolicyRemovalReason) advisoryDetailsHint() string {
	for _, a := range r.advisories {
		if !strings.HasPrefix(a.AdvisoryID, "PKSA-") {
			return "Review the advisory details above for more information."
		}
	}
	return "Go to https://packagist.org/security-advisories/ to find advisory details."
}

func (r *PolicyRemovalReason) filterListDescriptions() ([]string, []string) {
	lists := NewOMap[[]string]()
	for _, e := range r.filterListEntries {
		s := ""
		if e.Source != nil && *e.Source != "" && *e.Source != "0" {
			s += " reported by " + *e.Source
		}
		if e.URL != nil && *e.URL != "" && *e.URL != "0" {
			s += " (see " + *e.URL + ")"
		}
		if e.Reason != nil && *e.Reason != "" && *e.Reason != "0" {
			s += " reason: " + *e.Reason
		}
		appendTo(lists, e.ListName, s)
	}
	var names, descs []string
	lists.Each(func(name string, parts []string) {
		action := "filtered by "
		if name == MalwareName {
			action = "flagged as "
		}
		names = append(names, name)
		descs = append(descs, action+name+strings.Join(parts, ", "))
	})
	return names, descs
}

// PolicyRemovedVersions is Pool's $policyRemovedVersions:
// package name => normalized version => reason.
type PolicyRemovedVersions struct {
	m *OMap[*OMap[*PolicyRemovalReason]]
}

func NewPolicyRemovedVersions() *PolicyRemovedVersions {
	return &PolicyRemovedVersions{m: NewOMap[*OMap[*PolicyRemovalReason]]()}
}

// Clone copies the outer and inner maps (PHP array value semantics).
func (v *PolicyRemovedVersions) Clone() *PolicyRemovedVersions {
	c := NewPolicyRemovedVersions()
	if v == nil {
		return c
	}
	v.m.Each(func(name string, versions *OMap[*PolicyRemovalReason]) { c.m.Set(name, versions.Clone()) })
	return c
}

// Set stores a reason directly (Pool constructor input).
func (v *PolicyRemovedVersions) Set(name, version string, reason *PolicyRemovalReason) {
	versions, ok := v.m.Get(name)
	if !ok {
		versions = NewOMap[*PolicyRemovalReason]()
		v.m.Set(name, versions)
	}
	versions.Set(version, reason)
}

// Record is Pool::recordPolicyRemoval: records under every name of the
// package, but a replacer never hides the removal of the package itself.
func (v *PolicyRemovedVersions) Record(p *pkg.Package, reason *PolicyRemovalReason) {
	for _, name := range p.Names(false) {
		existing := v.Reason(name, p.Version())
		if existing == nil || (existing.PackageName() != name && p.Name == name) {
			v.Set(name, p.Version(), reason)
		}
	}
}

// Reason is Pool::getPolicyRemovalReason.
func (v *PolicyRemovedVersions) Reason(name, version string) *PolicyRemovalReason {
	if v == nil {
		return nil
	}
	versions, ok := v.m.Get(name)
	if !ok {
		return nil
	}
	r, _ := versions.Get(version)
	return r
}

// VersionReason is one entry of getPolicyRemovedVersions.
type VersionReason struct {
	Version string
	Reason  *PolicyRemovalReason
}

// Versions is Pool::getPolicyRemovedVersions: the entries of a package
// matching the constraint (nil = all), in ascending version_compare order.
func (v *PolicyRemovedVersions) Versions(name string, constraint semver.Constraint) []VersionReason {
	if v == nil {
		return nil
	}
	versions, ok := v.m.Get(name)
	if !ok {
		return nil
	}
	var out []VersionReason
	versions.Each(func(ver string, r *PolicyRemovalReason) {
		if constraint == nil || constraint.Matches(semver.NewConstraint("==", ver)) {
			out = append(out, VersionReason{ver, r})
		}
	})
	sort.SliceStable(out, func(i, j int) bool { return semver.VersionCompare(out[i].Version, out[j].Version) < 0 })
	return out
}

// Names returns the package names with recorded removals, in insertion order.
func (v *PolicyRemovedVersions) Names() []string {
	if v == nil {
		return nil
	}
	return v.m.Keys()
}

// All iterates getAllPolicyRemovedVersions in insertion order.
func (v *PolicyRemovedVersions) All(fn func(name, version string, reason *PolicyRemovalReason)) {
	if v == nil {
		return
	}
	v.m.Each(func(name string, versions *OMap[*PolicyRemovalReason]) {
		versions.Each(func(ver string, r *PolicyRemovalReason) { fn(name, ver, r) })
	})
}

// Pool is the part of Composer\DependencyResolver\Pool the policy filters
// consume and produce. The caller rebuilds its solver pool from Packages
// (keeping unacceptable fixed/locked packages and removed versions of the
// input pool, as the PHP filters do) and keeps PolicyRemoved on it.
type Pool struct {
	Packages      []*pkg.Package
	PolicyRemoved *PolicyRemovedVersions
}

// Request is the part of Composer\DependencyResolver\Request used by the
// pool filters.
type Request interface {
	IsLockedPackage(p *pkg.Package) bool
	// LockedRepositoryPackages returns the packages of the request's locked
	// repository, nil when it has none.
	LockedRepositoryPackages() []*pkg.Package
}

// NewRequest adapts a solver request: isLocked is Request::isLockedPackage
// and locked the packages of Request::getLockedRepository() (or nil).
func NewRequest(isLocked func(*pkg.Package) bool, locked []*pkg.Package) Request {
	return &requestAdapter{isLocked: isLocked, locked: locked}
}

type requestAdapter struct {
	isLocked func(*pkg.Package) bool
	locked   []*pkg.Package
}

func (r *requestAdapter) IsLockedPackage(p *pkg.Package) bool {
	return r.isLocked != nil && r.isLocked(p)
}
func (r *requestAdapter) LockedRepositoryPackages() []*pkg.Package { return r.locked }

// now is overridable in tests.
var now = time.Now
