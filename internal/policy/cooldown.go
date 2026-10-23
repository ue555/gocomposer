// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Policy/CooldownPolicyConfig.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// CooldownPolicyConfig is Composer\Policy\CooldownPolicyConfig.
type CooldownPolicyConfig struct {
	ListBase
	// Period is the cooldown in seconds; nil means no cooldown configured.
	Period *int64
}

func NewCooldownPolicyConfig(block bool, audit string, ignore *IgnoreMap, period *int64) *CooldownPolicyConfig {
	return &CooldownPolicyConfig{ListBase{CooldownName, block, audit, orEmpty(ignore)}, period}
}

// HasCooldown is CooldownPolicyConfig::hasCooldown.
func (c *CooldownPolicyConfig) HasCooldown() bool { return c.Period != nil && *c.Period > 0 }

// ShouldBlock is CooldownPolicyConfig::shouldBlock.
func (c *CooldownPolicyConfig) ShouldBlock(scope string) bool {
	if !c.HasCooldown() {
		return false
	}
	return c.baseShouldBlock(scope, false)
}

// PublishedTimeRequirer is implemented by Composer repositories
// (ComposerRepository::requiresPublishedTime).
type PublishedTimeRequirer interface {
	RepoName() string
	RequiresPublishedTime() bool
}

// GetEffectiveDate is CooldownPolicyConfig::getEffectiveDate.
func (c *CooldownPolicyConfig) GetEffectiveDate(p *pkg.Package) (time.Time, bool) {
	if t, ok := PublishedDate(p); ok {
		return t, true
	}
	return ReleaseDate(p)
}

// AssertPublishedTimeProvided is CooldownPolicyConfig::assertPublishedTimeProvided.
func (c *CooldownPolicyConfig) AssertPublishedTimeProvided(p *pkg.Package) error {
	if _, ok := PublishedDate(p); ok {
		return nil
	}
	r, ok := p.Base().Repo.(PublishedTimeRequirer)
	if !ok || !r.RequiresPublishedTime() {
		return nil
	}
	return newErr(RuntimeException, fmt.Sprintf(`%s %s from %s has no published-time, so the cooldown policy configured in "policy.cooldown" cannot tell whether it is new or old. Make sure the repository provides publication times (run "composer clear-cache" if it only recently started to), or set "require-published-time": false on that repository to fall back to the package-supplied "time" field, which offers weaker protection.`,
		p.PrettyName, p.PrettyVersion(), r.RepoName()))
}

// IsIgnored is CooldownPolicyConfig::isIgnored.
func (c *CooldownPolicyConfig) IsIgnored(p *pkg.Package, op string) bool {
	ignore := c.GetIgnoreForOperation(op)
	if ignore.Len() == 0 {
		return false
	}
	pc := semver.NewConstraint("==", p.Version())
	found := false
	ignore.Each(func(_ string, rules []*IgnorePackageRule) {
		for _, r := range rules {
			if !found && r.PackageNameRegex.MatchString(p.Name) && r.Constraint.Matches(pc) {
				found = true
			}
		}
	})
	return found
}

// IsWithinCooldown is CooldownPolicyConfig::isWithinCooldown.
func (c *CooldownPolicyConfig) IsWithinCooldown(effectiveDate, now time.Time) bool {
	if !c.HasCooldown() {
		return false
	}
	cutoff := now.Add(-time.Duration(*c.Period) * time.Second)
	return effectiveDate.After(cutoff)
}

// FormatTimeUntilAvailable is CooldownPolicyConfig::formatTimeUntilAvailable.
func (c *CooldownPolicyConfig) FormatTimeUntilAvailable(effectiveDate, now time.Time, short bool) string {
	var period int64
	if c.Period != nil {
		period = *c.Period
	}
	availableAt := effectiveDate.Add(time.Duration(period) * time.Second)
	d := availableAt.Sub(now)
	if d < 0 {
		d = -d
	}
	secs := int64(d / time.Second)
	days := secs / 86400
	h := (secs % 86400) / 3600
	i := (secs % 3600) / 60
	plural := func(n int64, word string) string {
		s := strconv.FormatInt(n, 10) + " " + word
		if n > 1 {
			s += "s"
		}
		return s
	}
	var parts []string
	if days > 0 {
		if short {
			parts = append(parts, strconv.FormatInt(days, 10)+"d")
		} else {
			parts = append(parts, plural(days, "day"))
		}
	}
	if h > 0 {
		if short {
			parts = append(parts, strconv.FormatInt(h, 10)+"h")
		} else {
			parts = append(parts, plural(h, "hour"))
		}
	}
	if len(parts) == 0 && i > 0 {
		if short {
			parts = append(parts, strconv.FormatInt(i, 10)+"m")
		} else {
			parts = append(parts, plural(i, "minute"))
		}
	}
	if r := strings.Join(parts, " "); r != "" {
		return r
	}
	if short {
		return "<1m"
	}
	return "less than a minute"
}

func (c *CooldownPolicyConfig) WithBlockingDisabled() *CooldownPolicyConfig {
	return NewCooldownPolicyConfig(false, c.Audit, c.Ignore, c.Period)
}

func (c *CooldownPolicyConfig) WithAudit(audit string) *CooldownPolicyConfig {
	return NewCooldownPolicyConfig(c.Block, audit, c.Ignore, c.Period)
}

// CooldownPolicyConfigFromRawConfig is CooldownPolicyConfig::fromRawConfig.
func CooldownPolicyConfigFromRawConfig(policyConfig phpjson.Value) (*CooldownPolicyConfig, error) {
	cooldownConfig := coalesce(policyConfig, CooldownName, phpjson.NewMap())
	envPeriod, envSet, err := cooldownEnvPeriod()
	if err != nil {
		return nil, err
	}
	if cooldownConfig == false {
		if !envSet || envPeriod == nil {
			return CooldownPolicyConfigDisabled(), nil
		}
		cooldownConfig = phpjson.NewMap()
	}
	if !phpjson.IsArray(cooldownConfig) {
		cooldownConfig = phpjson.NewMap()
	}
	period, err := ParseDuration(coalesce(cooldownConfig, "period", nil))
	if err != nil {
		return nil, err
	}
	if envSet {
		period = envPeriod
	}
	ignore, err := ParseIgnoreMap(coalesce(cooldownConfig, "ignore", nil))
	if err != nil {
		return nil, err
	}
	return NewCooldownPolicyConfig(toBool(coalesce(cooldownConfig, "block", true)), auditValue(coalesce(cooldownConfig, "audit", nil), AuditIgnore), ignore, period), nil
}

// cooldownEnvPeriod is CooldownPolicyConfig::getEnvPeriod; set=false means
// the env var is unset or empty.
func cooldownEnvPeriod() (*int64, bool, error) {
	v, ok := getEnv("COMPOSER_POLICY_COOLDOWN_PERIOD")
	if !ok || v == "" {
		return nil, false, nil
	}
	p, err := ParseDuration(v)
	if err != nil {
		return nil, false, &Error{Class: RuntimeException, Msg: "Invalid value for COMPOSER_POLICY_COOLDOWN_PERIOD: " + v + ". Use formats like '7 days', '24 hours', or an integer number of seconds.", Prev: err}
	}
	return p, true, nil
}

func CooldownPolicyConfigDisabled() *CooldownPolicyConfig {
	return NewCooldownPolicyConfig(false, AuditIgnore, nil, nil)
}

var (
	reAllDigits   = regexp.MustCompile(`^\d+$`)
	reDurationFmt = regexp.MustCompile(`^(\d+)[ \t\n\v\f\r]*(second|minute|hour|day|week)s?$`)
	reIsNumeric   = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?[ \t\n\r\v\f]*$`)
)

var durationUnits = map[string]int64{"second": 1, "minute": 60, "hour": 3600, "day": 86400, "week": 604800}

func phpIntval(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return math.MaxInt64
	}
	return n
}

// ParseDuration is CooldownPolicyConfig::parseDuration (nil = no cooldown).
func ParseDuration(duration phpjson.Value) (*int64, error) {
	formatErr := func(d string) error {
		return newErr(RuntimeException, "Invalid policy.cooldown.period format: "+d+". Use an integer number of seconds or a duration like '7 days', '24 hours', '30 minutes' or '1 week'.")
	}
	switch t := duration.(type) {
	case nil:
		return nil, nil
	case int64:
		if t == 0 {
			return nil, nil
		}
		if t < 0 {
			return nil, newErr(RuntimeException, fmt.Sprintf("Invalid policy.cooldown.period: duration cannot be negative (%d).", t))
		}
		return &t, nil
	case int:
		return ParseDuration(int64(t))
	case string:
		if t == "" || t == "0" {
			return nil, nil
		}
		if reAllDigits.MatchString(t) {
			n := phpIntval(t)
			return &n, nil
		}
		if reIsNumeric.MatchString(t) {
			f, _ := strconv.ParseFloat(strings.Trim(t, " \t\n\r\v\f"), 64)
			if f < 0 {
				return nil, newErr(RuntimeException, "Invalid policy.cooldown.period: duration cannot be negative ("+t+").")
			}
			return nil, formatErr(t)
		}
		if m := reDurationFmt.FindStringSubmatch(t); m != nil {
			n := phpIntval(m[1]) * durationUnits[m[2]]
			return &n, nil
		}
		return nil, formatErr(t)
	}
	return nil, formatErr(toStr(duration))
}

// --- package date accessors ---

// PublishedDate returns the repository-set publication date
// (PackageInterface::getPublishedDate).
var PublishedDate = func(p *pkg.Package) (time.Time, bool) {
	raw := p.PublishedDate()
	if raw == "" {
		return time.Time{}, false
	}
	return parsePackageDate(raw)
}

// ReleaseDate is PackageInterface::getReleaseDate.
var ReleaseDate = func(p *pkg.Package) (time.Time, bool) {
	raw := p.ReleaseDate()
	if raw == "" {
		return time.Time{}, false
	}
	return parsePackageDate(raw)
}

var reLooseDate = regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})(?:[ T](\d{1,2}):(\d{2})(?::(\d{2}))?)?$`)

// parsePackageDate parses like new \DateTime($s, new \DateTimeZone('UTC')),
// including PHP's overflow of out-of-range days (2024-04-31 => 2024-05-01).
func parsePackageDate(raw string) (time.Time, bool) {
	if t, ok := pkg.ParsePHPDate(raw); ok {
		return t, true
	}
	if m := reLooseDate.FindStringSubmatch(strings.TrimSpace(raw)); m != nil {
		at := func(i int) int {
			n, _ := strconv.Atoi(m[i])
			return n
		}
		return time.Date(at(1), time.Month(at(2)), at(3), at(4), at(5), at(6), 0, time.UTC), true
	}
	return time.Time{}, false
}

// formatAtom is DateTimeInterface::ATOM / DATE_RFC3339.
func formatAtom(t time.Time) string { return t.Format("2006-01-02T15:04:05-07:00") }
