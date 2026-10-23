// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Advisory/PartialSecurityAdvisory.php, SecurityAdvisory.php,
// IgnoredSecurityAdvisory.php, AuditConfig.php,
// src/Composer/Repository/AdvisoryProviderInterface.php and
// RepositorySet::getMatchingSecurityAdvisories (src/Composer/Repository/RepositorySet.php).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"regexp"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// Advisory is a PartialSecurityAdvisory, SecurityAdvisory or
// IgnoredSecurityAdvisory. Full is false for partial advisories, which only
// carry PackageName, AdvisoryID and AffectedVersions.
type Advisory struct {
	AdvisoryID       string
	PackageName      string
	AffectedVersions semver.Constraint

	Full       bool
	Title      string
	CVE        *string
	Link       *string
	ReportedAt time.Time
	// Sources is the raw non-empty list of {name, remoteId} entries.
	Sources  []phpjson.Value
	Severity *string

	Ignored      bool
	IgnoreReason *string
}

// NewPartialSecurityAdvisory mirrors the PartialSecurityAdvisory constructor.
func NewPartialSecurityAdvisory(packageName, advisoryID string, affected semver.Constraint) *Advisory {
	return &Advisory{AdvisoryID: advisoryID, PackageName: packageName, AffectedVersions: affected}
}

// NewSecurityAdvisory mirrors the SecurityAdvisory constructor.
func NewSecurityAdvisory(packageName, advisoryID string, affected semver.Constraint, title string, sources []phpjson.Value, reportedAt time.Time, cve, link, severity *string) *Advisory {
	return &Advisory{AdvisoryID: advisoryID, PackageName: packageName, AffectedVersions: affected, Full: true, Title: title, Sources: sources, ReportedAt: reportedAt, CVE: cve, Link: link, Severity: severity}
}

// ToIgnoredAdvisory is SecurityAdvisory::toIgnoredAdvisory.
func (a *Advisory) ToIgnoredAdvisory(reason *string) *Advisory {
	c := *a
	c.Ignored = true
	c.IgnoreReason = reason
	return &c
}

var reAffectedEssential = regexp.MustCompile(`(^[>=<^~]*[\d.]+).*`)

// CreateAdvisory is PartialSecurityAdvisory::create.
func CreateAdvisory(packageName string, data phpjson.Value) (*Advisory, error) {
	affected := toStr(coalesce(data, "affectedVersions", nil))
	c, err := semver.ParseConstraints(affected)
	if err != nil {
		c, err = semver.ParseConstraints(reAffectedEssential.ReplaceAllString(affected, "$1"))
		if err != nil {
			c = semver.NewConstraint("==", "0.0.0-invalid-version")
		}
	}
	id := toStr(coalesce(data, "advisoryId", nil))
	if isset(data, "title") && isset(data, "sources") && isset(data, "reportedAt") {
		reported, ok := parsePackageDate(toStr(coalesce(data, "reportedAt", nil)))
		if !ok {
			return nil, newErr("Exception", "Failed to parse time string ("+toStr(coalesce(data, "reportedAt", nil))+")")
		}
		return NewSecurityAdvisory(packageName, id, c, toStr(coalesce(data, "title", nil)), arrayValues(coalesce(data, "sources", nil)), reported,
			nullableStr(coalesce(data, "cve", nil)), nullableStr(coalesce(data, "link", nil)), nullableStr(coalesce(data, "severity", nil))), nil
	}
	return NewPartialSecurityAdvisory(packageName, id, c), nil
}

// JSONValue is jsonSerialize() in (array) property order.
func (a *Advisory) JSONValue() *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("advisoryId", a.AdvisoryID)
	m.Set("packageName", a.PackageName)
	m.Set("affectedVersions", a.AffectedVersions.PrettyString())
	if !a.Full {
		return m
	}
	nullable := func(k string, v *string) {
		if v == nil {
			m.Set(k, nil)
		} else {
			m.Set(k, *v)
		}
	}
	m.Set("title", a.Title)
	nullable("cve", a.CVE)
	nullable("link", a.Link)
	m.Set("reportedAt", formatAtom(a.ReportedAt))
	m.Set("sources", a.Sources)
	nullable("severity", a.Severity)
	if a.Ignored && a.IgnoreReason != nil {
		m.Set("ignoreReason", *a.IgnoreReason)
	}
	return m
}

// AdvisoryMap is array<string(package name), list<advisory>>.
type AdvisoryMap = OMap[[]*Advisory]

// AdvisoryProvider is Composer\Repository\AdvisoryProviderInterface.
type AdvisoryProvider interface {
	HasSecurityAdvisories() (bool, error)
	// GetSecurityAdvisories returns the 'advisories' part (namesFound is unused by callers).
	GetSecurityAdvisories(packageConstraintMap *ConstraintMap, allowPartial bool) (*AdvisoryMap, error)
}

// AdvisoryResult is array{advisories, unreachableRepos}.
type AdvisoryResult struct {
	Advisories       *AdvisoryMap
	UnreachableRepos []string
}

// AdvisorySource is the mockable RepositorySet::getMatchingSecurityAdvisories.
type AdvisorySource interface {
	GetMatchingSecurityAdvisories(packages []*pkg.Package, allowPartial, ignoreUnreachable bool) (*AdvisoryResult, error)
}

// RepositoryAdvisorySource is a RepositorySet over the given repositories
// for advisory lookups.
type RepositoryAdvisorySource struct {
	Repositories []repo.Repository
	Resolver     *ProviderResolver
}

// GetMatchingSecurityAdvisories is RepositorySet::getMatchingSecurityAdvisories.
func (s *RepositoryAdvisorySource) GetMatchingSecurityAdvisories(packages []*pkg.Package, allowPartial, ignoreUnreachable bool) (*AdvisoryResult, error) {
	pcm := constraintMapFor(packages, true)
	var unreachable []string
	var repoAdvisories []*AdvisoryMap
	for _, r := range s.Repositories {
		p := s.Resolver.AdvisoryProviderFor(r)
		if p == nil {
			continue
		}
		has, err := p.HasSecurityAdvisories()
		if err == nil && has {
			var adv *AdvisoryMap
			adv, err = p.GetSecurityAdvisories(pcm, allowPartial)
			if err == nil {
				repoAdvisories = append(repoAdvisories, adv)
			}
		}
		if err != nil {
			if _, ok := isTransportError(err); ok && ignoreUnreachable {
				unreachable = append(unreachable, err.Error())
				continue
			}
			return nil, err
		}
	}
	// array_merge_recursive then ksort
	merged := NewOMap[[]*Advisory]()
	for _, ra := range repoAdvisories {
		ra.Each(func(name string, list []*Advisory) {
			cur, _ := merged.Get(name)
			merged.Set(name, append(cur, list...))
		})
	}
	merged.KSort()
	return &AdvisoryResult{Advisories: merged, UnreachableRepos: unreachable}, nil
}

// AdvisoriesFromRawData builds the advisories of an AdvisoryProvider from raw
// "security-advisories" data (name => list of advisory arrays), keeping
// those whose affected versions match the requested constraint, like
// PackageRepository::getSecurityAdvisories and ComposerRepository's $create.
func AdvisoriesFromRawData(raw phpjson.Value, pcm *ConstraintMap) (*AdvisoryMap, error) {
	out := NewOMap[[]*Advisory]()
	var err error
	eachArr(raw, func(name string, list phpjson.Value) {
		if err != nil {
			return
		}
		c, ok := pcm.Get(name)
		if !ok {
			return
		}
		var kept []*Advisory
		eachArr(list, func(_ string, data phpjson.Value) {
			if err != nil {
				return
			}
			a, e := CreateAdvisory(name, data)
			if e != nil {
				err = e
				return
			}
			if a.AffectedVersions.Matches(c) {
				kept = append(kept, a)
			}
		})
		if len(kept) > 0 {
			out.Set(name, kept)
		}
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuditConfig is Composer\Advisory\AuditConfig.
type AuditConfig struct {
	Audit       bool
	AuditFormat string
}

// NewAuditConfig mirrors the defaults (audit=true, format=summary).
func NewAuditConfig() AuditConfig { return AuditConfig{Audit: true, AuditFormat: FormatSummary} }
