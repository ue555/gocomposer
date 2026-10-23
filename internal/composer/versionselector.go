package composer

// Port of Composer\Package\Version\VersionSelector (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// CooldownFilter filters candidates still in a cooldown period
// (VersionSelector::filterCooldownCandidates); nil disables the check.
type CooldownFilter func(candidates []*pkg.Package, io *cio.IO, showWarnings bool) []*pkg.Package

// VersionSelector selects the best version of a package.
type VersionSelector struct {
	repoSet             *solver.RepositorySet
	platformConstraints map[string][]semver.Constraint
	platformNames       []string
	Cooldown            CooldownFilter
	// PHPVersion is PHP_MAJOR.MINOR.RELEASE of the running PHP (for ext-*).
	PHPVersion       string
	PlatformRejected []*pkg.Package
	// FirstClearingWithheld / WithheldCandidates are filled by the cooldown filter.
	FirstClearingWithheld *Withheld
	WithheldCandidates    []Withheld
	cooldownPolicy        *policy.CooldownPolicyConfig
	err                   error
	// ShowWarningsFunc is the callable form of findBestCandidate's
	// $showWarnings: when set, a candidate warning is printed if the bool
	// argument is true or the callback returns true for the candidate.
	ShowWarningsFunc func(candidate *pkg.Package) bool
}

func (vs *VersionSelector) shouldWarn(p *pkg.Package, showWarnings bool) bool {
	return showWarnings || (vs.ShowWarningsFunc != nil && vs.ShowWarningsFunc(p))
}

// NewVersionSelector is VersionSelector::__construct.
func NewVersionSelector(rs *solver.RepositorySet, platformRepo repo.Repository) *VersionSelector {
	vs := &VersionSelector{repoSet: rs, platformConstraints: map[string][]semver.Constraint{}}
	if platformRepo != nil {
		ps, _ := platformRepo.Packages()
		for _, p := range ps {
			if _, ok := vs.platformConstraints[p.Name]; !ok {
				vs.platformNames = append(vs.platformNames, p.Name)
			}
			vs.platformConstraints[p.Name] = append(vs.platformConstraints[p.Name], semver.NewConstraint("==", p.Version()))
		}
	}
	return vs
}

// FindBestCandidate is VersionSelector::findBestCandidate; nil means false.
func (vs *VersionSelector) FindBestCandidate(name, targetVersion, preferredStability string, filter solver.PlatformFilter, flags int, io *cio.IO, showWarnings bool) (*pkg.Package, error) {
	minPriority, ok := pkg.Stabilities[preferredStability]
	if !ok {
		return nil, &UserError{Msg: "Expected a valid stability name as 3rd argument, got " + preferredStability, Code: 1}
	}
	if filter == nil {
		filter = solver.IgnoreNothing{}
	}
	var constraint semver.Constraint
	if targetVersion != "" {
		c, err := semver.ParseConstraints(targetVersion)
		if err != nil {
			return nil, err
		}
		constraint = c
	}
	found, err := vs.repoSet.FindPackagesE(strings.ToLower(name), constraint, flags)
	if err != nil {
		return nil, err
	}
	candidates := append([]*pkg.Package(nil), found...)
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		ap, bp := a.StabilityPriority(), b.StabilityPriority()
		cmp := 0
		switch {
		case minPriority < ap && bp < ap:
			cmp = 1
		case minPriority < ap && ap < bp:
			cmp = -1
		case minPriority >= ap && minPriority < bp:
			cmp = -1
		default:
			cmp = semver.VersionCompare(b.Version(), a.Version())
		}
		return cmp < 0
	})
	vs.PlatformRejected = nil
	vs.FirstClearingWithheld, vs.WithheldCandidates, vs.err = nil, nil, nil
	if vs.Cooldown != nil {
		candidates = vs.Cooldown(candidates, io, showWarnings)
		if vs.err != nil {
			return nil, vs.err
		}
	}
	var selected *pkg.Package
	_, ignoreAll := filter.(solver.IgnoreAll)
	if len(vs.platformConstraints) > 0 && !ignoreAll {
		warned := map[string]bool{}
		seen := map[string]bool{}
		for _, p := range candidates {
			skip := false
			for _, link := range p.Requires() {
				if !repo.IsPlatformPackage(link.Target) || filter.IsIgnored(link.Target) {
					continue
				}
				var reason string
				if provided, ok := vs.platformConstraints[link.Target]; ok {
					satisfied := false
					for _, pc := range provided {
						if link.Constraint.Matches(pc) {
							satisfied = true
							break
						}
						if il, ok := filter.(*solver.IgnoreList); ok && il.IsUpperBoundIgnored(link.Target) {
							if il.FilterConstraint(link.Target, link.Constraint, true).Matches(pc) {
								satisfied = true
								break
							}
						}
					}
					if satisfied {
						continue
					}
					reason = "is not satisfied by your platform"
				} else {
					reason = "is missing from your platform"
				}
				isLatest := !seen[p.Name]
				seen[p.Name] = true
				if io != nil && vs.shouldWarn(p, showWarnings) {
					key := p.Name + "/" + link.Target
					first := !warned[key]
					warned[key] = true
					latest := ""
					if isLatest {
						latest = "'s latest version"
					}
					verbosity := cio.Normal
					if !first {
						verbosity = cio.Verbose
					}
					io.WriteErrorV("<warning>Cannot use "+p.PrettyName+latest+" "+p.PrettyVersion()+" as it "+link.Description+" "+link.Target+" "+link.PrettyConstraint+" which "+reason+".</>", true, verbosity)
				}
				skip = true
			}
			if skip {
				vs.PlatformRejected = append(vs.PlatformRejected, p)
				continue
			}
			selected = p
			break
		}
	} else if len(candidates) > 0 {
		selected = candidates[0]
	}
	if selected == nil {
		return nil, nil
	}
	if selected.AliasOf != nil && selected.Version() == semver.DefaultBranchAlias {
		selected = selected.AliasOf
	}
	return selected, nil
}

var reBranchAliasDev = regexp.MustCompile(`^(\d+\.\d+\.\d+)(\.9999999)-dev$`)
var reSemverPart = regexp.MustCompile(`^\d+\D?`)

// FindRecommendedRequireVersion is VersionSelector::findRecommendedRequireVersion.
func (vs *VersionSelector) FindRecommendedRequireVersion(p *pkg.Package) string {
	if strings.HasPrefix(p.Name, "ext-") && vs.PHPVersion != "" {
		parts := strings.Split(p.Version(), ".")
		if len(parts) > 3 {
			parts = parts[:3]
		}
		if strings.Join(parts, ".") == vs.PHPVersion {
			return "*"
		}
	}
	if !p.IsDev() {
		return transformVersion(p.Version(), p.PrettyVersion(), p.Stability())
	}
	extra, _ := pkg.GetBranchAlias(pkg.Dump(p))
	if extra != "" && extra != semver.DefaultBranchAlias {
		if reBranchAliasDev.MatchString(extra) {
			extra = reBranchAliasDev.ReplaceAllString(extra, "${1}.0")
			extra = strings.ReplaceAll(extra, ".9999999", ".0")
			return transformVersion(extra, extra, "dev")
		}
	}
	return p.PrettyVersion()
}

func transformVersion(version, prettyVersion, stability string) string {
	parts := strings.Split(version, ".")
	if len(parts) == 4 && reSemverPart.MatchString(parts[3]) {
		if parts[0] == "0" {
			parts = parts[:3]
		} else {
			parts = parts[:2]
		}
		version = strings.Join(parts, ".")
	} else {
		return prettyVersion
	}
	if stability != "stable" {
		version += "@" + stability
	}
	return "^" + version
}

// Withheld is a candidate held back by the cooldown policy.
type Withheld struct {
	Package     *pkg.Package
	ReleaseDate time.Time
}

// SetCooldownPolicy enables VersionSelector::filterCooldownCandidates for a
// blocking cooldown policy (nil disables it).
func (vs *VersionSelector) SetCooldownPolicy(c *policy.CooldownPolicyConfig) {
	vs.cooldownPolicy = c
	if c == nil || !c.HasCooldown() || !c.Block {
		vs.Cooldown = nil
		return
	}
	vs.Cooldown = func(candidates []*pkg.Package, io *cio.IO, showWarnings bool) []*pkg.Package {
		now := time.Now()
		seen := map[string]bool{}
		accepted := map[string]bool{}
		var result []*pkg.Package
		for _, p := range candidates {
			target := p
			if p.AliasOf != nil {
				target = p.AliasOf
			}
			if target.IsDev() || c.IsIgnored(target, "block") {
				result = append(result, p)
				accepted[p.Name] = true
				continue
			}
			if err := c.AssertPublishedTimeProvided(target); err != nil {
				vs.err = err
				return nil
			}
			date, ok := c.GetEffectiveDate(target)
			if !ok || !c.IsWithinCooldown(date, now) {
				result = append(result, p)
				accepted[p.Name] = true
				continue
			}
			if vs.FirstClearingWithheld == nil || date.Before(vs.FirstClearingWithheld.ReleaseDate) {
				vs.FirstClearingWithheld = &Withheld{p, date}
			}
			vs.WithheldCandidates = append(vs.WithheldCandidates, Withheld{p, date})
			if accepted[p.Name] {
				continue
			}
			isLatest := !seen[p.Name]
			seen[p.Name] = true
			if io != nil && vs.shouldWarn(p, showWarnings) {
				source := "time"
				if target.PublishedDate() != "" {
					source = "published-time"
				}
				reason := policy.RemovalReasonCooldown(p.Name, p.PrettyVersion(), date.Format("2006-01-02T15:04:05-07:00"), c.FormatTimeUntilAvailable(date, now, false), source)
				latest := ""
				verbosity := cio.Verbose
				if isLatest {
					latest = "'s latest version"
					verbosity = cio.Normal
				}
				io.WriteErrorV("<warning>Cannot use "+p.PrettyName+latest+" "+p.PrettyVersion()+" as it "+reason.Verb(false)+" "+reason.Description()+".</>", true, verbosity)
			}
		}
		return result
	}
}
