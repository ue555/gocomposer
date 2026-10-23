// Test helpers for the ports of Composer's PHPUnit tests (Composer
// f83588392df43d4fd027d84d26de459946ef4e80, tests/Composer/Test/...).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

var policyEnvVars = []string{
	"COMPOSER_POLICY", "COMPOSER_NO_BLOCKING", "COMPOSER_NO_SECURITY_BLOCKING",
	"COMPOSER_POLICY_ADVISORIES_BLOCK", "COMPOSER_POLICY_MALWARE_BLOCK", "COMPOSER_POLICY_ABANDONED_BLOCK",
	"COMPOSER_SECURITY_BLOCKING_ABANDONED", "COMPOSER_AUDIT_ABANDONED", "COMPOSER_POLICY_COOLDOWN_PERIOD",
	"COMPOSER_POLICY_COOLDOWN_BLOCK",
}

// clearEnv unsets the policy env vars for the duration of the test
// (Platform::clearEnv in tearDown).
func clearEnv(t *testing.T) {
	t.Helper()
	for _, n := range policyEnvVars {
		t.Setenv(n, "")
		os.Unsetenv(n)
	}
}

func putEnv(t *testing.T, name, value string) {
	t.Helper()
	t.Setenv(name, value)
}

// j decodes a JSON literal into a PHP-like value.
func j(t *testing.T, s string) phpjson.Value {
	t.Helper()
	v, err := phpjson.Decode([]byte(s))
	if err != nil {
		t.Fatalf("bad json %s: %v", s, err)
	}
	return v
}

// newConfig is "new Config(); $config->merge(['config' => <json>])".
func newConfig(t *testing.T, configJSON string) *config.Config {
	t.Helper()
	c := config.New(true, t.TempDir())
	if configJSON != "" {
		c.Merge(phpjson.AsMap(j(t, `{"config": `+configJSON+`}`)))
	}
	return c
}

func mustPolicy(t *testing.T, c *config.Config) *PolicyConfig {
	t.Helper()
	pc, err := FromConfig(c)
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}
	return pc
}

func expectError(t *testing.T, err error, class, contains string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s containing %q, got no error", class, contains)
	}
	var pe *Error
	if class != "" {
		if !errors.As(err, &pe) || pe.Class != class {
			t.Fatalf("expected %s, got %T %v", class, err, err)
		}
	}
	if !strings.Contains(err.Error(), contains) {
		t.Fatalf("expected message containing %q, got %q", contains, err.Error())
	}
}

func sp(s string) *string { return &s }

func i64(n int64) *int64 { return &n }

func describeStr(s *string) string {
	if s == nil {
		return "null"
	}
	return fmt.Sprintf("%q", *s)
}

func describeConstraint(c semver.Constraint) string {
	if _, ok := c.(*semver.MatchAll); ok {
		return "MatchAll"
	}
	return c.String() + "|" + c.PrettyString()
}

func describeIgnore(m *IgnoreMap) string {
	var sb strings.Builder
	m.Each(func(k string, rules []*IgnorePackageRule) {
		sb.WriteString(k + ":[")
		for _, r := range rules {
			fmt.Fprintf(&sb, "{%s %s %s %v %v %s}", r.PackageName, describeConstraint(r.Constraint), describeStr(r.Reason), r.OnBlock, r.OnAudit, r.PackageNameRegex)
		}
		sb.WriteString("]")
	})
	return sb.String()
}

func describeIDs(m *OMap[*IgnoreIdRule]) string {
	var sb strings.Builder
	m.Each(func(k string, r *IgnoreIdRule) {
		fmt.Fprintf(&sb, "%s:{%s %s %v %v}", k, r.ID, describeStr(r.Reason), r.OnBlock, r.OnAudit)
	})
	return sb.String()
}

func describeSeverities(m *OMap[*IgnoreSeverityRule]) string {
	var sb strings.Builder
	m.Each(func(k string, r *IgnoreSeverityRule) {
		fmt.Fprintf(&sb, "%s:{%s %s %v %v}", k, r.Severity, describeStr(r.Reason), r.OnBlock, r.OnAudit)
	})
	return sb.String()
}

func describePolicy(v any) string {
	switch c := v.(type) {
	case *AdvisoriesPolicyConfig:
		return fmt.Sprintf("advisories{%v %s ignore=%s ids=%s sev=%s}", c.Block, c.Audit, describeIgnore(c.Ignore), describeIDs(c.IgnoreID), describeSeverities(c.IgnoreSeverity))
	case *MalwarePolicyConfig:
		return fmt.Sprintf("malware{%v %s %s ignore=%s src=%v}", c.Block, c.Audit, c.BlockScope, describeIgnore(c.Ignore), c.IgnoreSource)
	case *AbandonedPolicyConfig:
		return fmt.Sprintf("abandoned{%v %s ignore=%s}", c.Block, c.Audit, describeIgnore(c.Ignore))
	case *CooldownPolicyConfig:
		p := "null"
		if c.Period != nil {
			p = fmt.Sprint(*c.Period)
		}
		return fmt.Sprintf("cooldown{%v %s ignore=%s period=%s}", c.Block, c.Audit, describeIgnore(c.Ignore), p)
	case *CustomListPolicyConfig:
		var srcs []string
		for _, s := range c.Sources {
			srcs = append(srcs, s.ListName+"="+s.URL)
		}
		return fmt.Sprintf("custom{%s %v %s %s ignore=%s src=%v}", c.Name, c.Block, c.Audit, c.BlockScope, describeIgnore(c.Ignore), srcs)
	}
	return fmt.Sprintf("%#v", v)
}

func assertEqualPolicy(t *testing.T, expected, actual any) {
	t.Helper()
	if e, a := describePolicy(expected), describePolicy(actual); e != a {
		t.Fatalf("policy mismatch\nexpected: %s\nactual:   %s", e, a)
	}
}

func ignoreMap(rules ...*IgnorePackageRule) *IgnoreMap {
	m := NewOMap[[]*IgnorePackageRule]()
	for _, r := range rules {
		appendTo(m, r.PackageName, r)
	}
	return m
}

func rule(name string, c semver.Constraint, reason *string, onBlock, onAudit bool) *IgnorePackageRule {
	if c == nil {
		c = &semver.MatchAll{}
	}
	return NewIgnorePackageRule(name, c, reason, onBlock, onAudit)
}

func flat(m *OMap[*string]) string {
	var parts []string
	m.Each(func(k string, v *string) { parts = append(parts, k+"=>"+describeStr(v)) })
	return "[" + strings.Join(parts, ", ") + "]"
}

// ---- packages ----

// newPkg is "new Package($name, $version, $prettyVersion)".
func newPkg(name, version, pretty string) *pkg.Package { return pkg.NewPackage(name, version, pretty) }

// loadPkg builds a package through the array loader (for fields without
// setters such as published-time, abandoned or replace), then forces the
// exact version strings like "new Package(...)" would.
func loadPkg(t *testing.T, version, pretty, data string) *pkg.Package {
	t.Helper()
	p, err := (&pkg.Loader{}).Load(phpjson.AsMap(j(t, data)))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	p.ReplaceVersion(version, pretty)
	return p
}

func utc(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.UTC)
	if err != nil {
		panic(err)
	}
	return t
}

// ---- IO ----

// bufferIO is Composer\IO\BufferIO: stdout and stderr share one
// non-decorated stream whose formatter only knows Symfony's default styles.
type bufferIO struct {
	sb strings.Builder
}

func (b *bufferIO) Write(msg string) {
	b.sb.WriteString(Formatter{Styles: SymfonyDefaultStyles}.Format(msg) + "\n")
}
func (b *bufferIO) WriteError(msg string) { b.Write(msg) }
func (b *bufferIO) Output() string        { return b.sb.String() }

// trimmed is trim(str_replace("\r", ”, $io->getOutput())).
func (b *bufferIO) trimmed() string {
	return strings.TrimSpace(strings.ReplaceAll(b.sb.String(), "\r", ""))
}

func packageNames(ps []*pkg.Package) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.PrettyString())
	}
	return out
}

func assertSamePackages(t *testing.T, expected, actual []*pkg.Package) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("expected packages %v, got %v", packageNames(expected), packageNames(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("expected packages %v, got %v", packageNames(expected), packageNames(actual))
		}
	}
}

func pkgs(ps ...*pkg.Package) []*pkg.Package { return ps }
