// Ports of Composer f83588392df43d4fd027d84d26de459946ef4e80
// tests/Composer/Test/Policy/*Test.php (expectations unchanged).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"strings"
	"testing"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// ---------------------------------------------------------- PolicyConfigTest

func TestPolicyConfigRejectsReservedCustomListNames(t *testing.T) {
	for _, name := range []string{"ignore-foo", "ignoremalware", "package", "packages", "license", "licence", "support", "maintenance", "security", "minimum-release-age"} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			c := newConfig(t, `{"policy": {"`+name+`": {"block": true}}}`)
			_, err := FromConfig(c)
			expectError(t, err, UnexpectedValueException, name)
		})
	}
}

func TestPolicyConfigAllowsIgnoreUnreachableSiblingKey(t *testing.T) {
	clearEnv(t)
	pc := mustPolicy(t, newConfig(t, `{"policy": {"ignore-unreachable": true}}`))
	if !pc.IgnoreUnreachable.Update || pc.CustomLists.Len() != 0 {
		t.Fatal("unexpected")
	}
}

func TestPolicyConfigAllowsRegularCustomListName(t *testing.T) {
	clearEnv(t)
	pc := mustPolicy(t, newConfig(t, `{"policy": {"company-policy": {"block": true}}}`))
	if !pc.CustomLists.Has("company-policy") {
		t.Fatal("missing company-policy")
	}
}

func TestComposerPolicyAdvisoriesBlock(t *testing.T) {
	for name, tc := range map[string]struct {
		env      string
		block    bool
		expected bool
	}{"enable": {"1", false, true}, "disable": {"0", true, false}} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_POLICY_ADVISORIES_BLOCK", tc.env)
			pc := mustPolicy(t, newConfig(t, `{"policy": {"advisories": {"block": `+boolJSON(tc.block)+`, "audit": "report"}}}`))
			if pc.Advisories.Block != tc.expected || pc.Advisories.Audit != AuditReport {
				t.Fatalf("got %v %s", pc.Advisories.Block, pc.Advisories.Audit)
			}
		})
	}
}

func boolJSON(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

var abandonedBlockProvider = map[string]struct {
	env      string
	block    bool
	expected bool
}{"enable": {"1", false, true}, "disable": {"0", true, false}}

func TestComposerSecurityBlockingAbandoned(t *testing.T) {
	for name, tc := range abandonedBlockProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_SECURITY_BLOCKING_ABANDONED", tc.env)
			pc := mustPolicy(t, newConfig(t, `{"policy": {"abandoned": {"block": `+boolJSON(tc.block)+`, "audit": "report"}}}`))
			if pc.Abandoned.Block != tc.expected || pc.Abandoned.Audit != AuditReport {
				t.Fatal("mismatch")
			}
		})
	}
}

func TestComposerPolicyAbandonedBlock(t *testing.T) {
	for name, tc := range abandonedBlockProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_POLICY_ABANDONED_BLOCK", tc.env)
			pc := mustPolicy(t, newConfig(t, `{"policy": {"abandoned": {"block": `+boolJSON(tc.block)+`, "audit": "report"}}}`))
			if pc.Abandoned.Block != tc.expected || pc.Abandoned.Audit != AuditReport {
				t.Fatal("mismatch")
			}
		})
	}
}

func TestComposerPolicyAbandonedBlockTakesPrecedenceOverLegacyAlias(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_POLICY_ABANDONED_BLOCK", "1")
	putEnv(t, "COMPOSER_SECURITY_BLOCKING_ABANDONED", "0")
	if !mustPolicy(t, newConfig(t, "")).Abandoned.Block {
		t.Fatal("expected block")
	}
}

func TestLegacyAbandonedBlockEnvVarStillWorksWhenCanonicalUnset(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_SECURITY_BLOCKING_ABANDONED", "1")
	if !mustPolicy(t, newConfig(t, "")).Abandoned.Block {
		t.Fatal("expected block")
	}
}

func TestComposerSecurityBlockingAbandonedWithAuditConfig(t *testing.T) {
	for name, tc := range abandonedBlockProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_SECURITY_BLOCKING_ABANDONED", tc.env)
			pc := mustPolicy(t, newConfig(t, `{"audit": {"block-abandoned": `+boolJSON(tc.block)+`}}`))
			if pc.Abandoned.Block != tc.expected {
				t.Fatal("mismatch")
			}
		})
	}
}

var abandonedAuditProvider = map[string][3]string{"report": {"report", "fail", "report"}, "fail": {"fail", "report", "fail"}}

func TestComposerAuditAbandonedSetsAuditMode(t *testing.T) {
	for name, tc := range abandonedAuditProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_AUDIT_ABANDONED", tc[0])
			pc := mustPolicy(t, newConfig(t, `{"policy": {"abandoned": {"audit": "`+tc[1]+`"}}}`))
			if pc.Abandoned.Audit != tc[2] {
				t.Fatalf("got %s", pc.Abandoned.Audit)
			}
		})
	}
}

func TestComposerAuditAbandonedSetsAuditModeWithAuditConfig(t *testing.T) {
	for name, tc := range abandonedAuditProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_AUDIT_ABANDONED", tc[0])
			pc := mustPolicy(t, newConfig(t, `{"audit": {"abandoned": "`+tc[1]+`"}}`))
			if pc.Abandoned.Audit != tc[2] {
				t.Fatalf("got %s", pc.Abandoned.Audit)
			}
		})
	}
}

func TestComposerPolicyMalwareBlock(t *testing.T) {
	for name, tc := range abandonedBlockProvider {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, "COMPOSER_POLICY_MALWARE_BLOCK", tc.env)
			pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": {"block": `+boolJSON(tc.block)+`, "audit": "report"}}}`))
			if pc.Malware.Block != tc.expected || pc.Malware.Audit != AuditReport {
				t.Fatal("mismatch")
			}
		})
	}
}

func TestBothAbandonedEnvVarsApplyIndependently(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_SECURITY_BLOCKING_ABANDONED", "1")
	putEnv(t, "COMPOSER_AUDIT_ABANDONED", "report")
	pc := mustPolicy(t, newConfig(t, ""))
	if !pc.Abandoned.Block || pc.Abandoned.Audit != AuditReport {
		t.Fatal("mismatch")
	}
}

func TestEnvBlockOverridesWhenListExplicitlyDisabled(t *testing.T) {
	cases := []struct {
		env, cfg string
		check    func(*PolicyConfig) bool
	}{
		{"COMPOSER_POLICY_ADVISORIES_BLOCK", `{"policy": {"advisories": false}}`, func(p *PolicyConfig) bool { return p.Advisories.Block }},
		{"COMPOSER_POLICY_MALWARE_BLOCK", `{"policy": {"malware": false}}`, func(p *PolicyConfig) bool { return p.Malware.Block }},
		{"COMPOSER_POLICY_ABANDONED_BLOCK", `{"policy": {"abandoned": false}}`, func(p *PolicyConfig) bool { return p.Abandoned.Block }},
		{"COMPOSER_SECURITY_BLOCKING_ABANDONED", `{"policy": {"abandoned": false}}`, func(p *PolicyConfig) bool { return p.Abandoned.Block }},
	}
	for _, tc := range cases {
		t.Run(tc.env, func(t *testing.T) {
			clearEnv(t)
			putEnv(t, tc.env, "1")
			if !tc.check(mustPolicy(t, newConfig(t, tc.cfg))) {
				t.Fatal("expected block")
			}
		})
	}
}

func TestCooldownEnvAgeIgnoredWhenPolicyDisabled(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_POLICY_COOLDOWN_PERIOD", "2 days")
	if mustPolicy(t, newConfig(t, `{"policy": false}`)).Cooldown.HasCooldown() {
		t.Fatal("expected no cooldown")
	}
}

func TestComposerAuditAbandonedOverridesWhenAbandonedExplicitlyDisabled(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_AUDIT_ABANDONED", "fail")
	if mustPolicy(t, newConfig(t, `{"policy": {"abandoned": false}}`)).Abandoned.Audit != AuditFail {
		t.Fatal("expected fail")
	}
}

func TestWithIgnoreUnreachableOnlyAffectsRequestedScope(t *testing.T) {
	clearEnv(t)
	pc := mustPolicy(t, newConfig(t, `{"policy": {"ignore-unreachable": ["update"]}}`))
	if pc.IgnoreUnreachable != (IgnoreUnreachable{false, false, true}) {
		t.Fatalf("got %+v", pc.IgnoreUnreachable)
	}
	updated, err := pc.WithIgnoreUnreachable("audit")
	if err != nil || updated.IgnoreUnreachable != (IgnoreUnreachable{true, false, true}) {
		t.Fatalf("got %+v %v", updated.IgnoreUnreachable, err)
	}
}

// ---------------------------------------------------- AdvisoriesPolicyConfigTest

func TestAdvisoriesDefaultConfig(t *testing.T) {
	for _, raw := range []string{`[]`, `{"advisories": true}`, `{"advisories": []}`} {
		clearEnv(t)
		got, err := AdvisoriesPolicyConfigFromRawConfig(j(t, raw), j(t, `[]`))
		if err != nil {
			t.Fatal(err)
		}
		assertEqualPolicy(t, NewAdvisoriesPolicyConfig(true, AuditFail, nil, nil, nil), got)
	}
}

func TestAdvisoriesFromRawConfig(t *testing.T) {
	raw := j(t, `{"advisories": {"block": true, "audit": "report",
		"ignore": {"acme/abandoned": "flagged by mistake", "acme/abandoned2": {"constraint": "1.0"}},
		"ignore-severity": {"low": {"reason": "reason", "on-block": false, "on-audit": false}, "high": "ignore"},
		"ignore-id": {"CVE-2024-1234": "flagged by mistake", "CVE-2024-1235": {"on-block": false}}}}`)
	got, err := AdvisoriesPolicyConfigFromRawConfig(raw, j(t, `[]`))
	if err != nil {
		t.Fatal(err)
	}
	ids := NewOMap[*IgnoreIdRule]()
	ids.Set("CVE-2024-1234", NewIgnoreIdRule("CVE-2024-1234", sp("flagged by mistake"), true, true))
	ids.Set("CVE-2024-1235", NewIgnoreIdRule("CVE-2024-1235", nil, false, true))
	sev := NewOMap[*IgnoreSeverityRule]()
	sev.Set("low", NewIgnoreSeverityRule("low", sp("reason"), false, false))
	sev.Set("high", NewIgnoreSeverityRule("high", sp("ignore"), true, true))
	assertEqualPolicy(t, NewAdvisoriesPolicyConfig(true, AuditReport, ignoreMap(
		rule("acme/abandoned", nil, sp("flagged by mistake"), true, true),
		rule("acme/abandoned2", semver.MustParseConstraints("1.0"), nil, true, true),
	), ids, sev), got)
}

func TestAdvisoriesShouldBlockNeverAppliesToInstallScope(t *testing.T) {
	a := NewAdvisoriesPolicyConfig(true, AuditFail, nil, nil, nil)
	if !a.ShouldBlock(BlockScopeUpdate) || a.ShouldBlock(BlockScopeInstall) {
		t.Fatal("mismatch")
	}
	a = NewAdvisoriesPolicyConfig(false, AuditFail, nil, nil, nil)
	if a.ShouldBlock(BlockScopeUpdate) || a.ShouldBlock(BlockScopeInstall) {
		t.Fatal("mismatch")
	}
}

func TestAdvisoriesFromAuditConfig(t *testing.T) {
	audit := j(t, `{"block": true, "ignore": {"acme/abandoned": "flagged by mistake", "acme/abandoned2": {"apply": "block"},
		"CVE-2024-1234": "flagged by mistake", "CVE-2024-1235": {"apply": "audit"}},
		"ignore-severity": {"low": {"apply": "block"}}}`)
	got, err := AdvisoriesPolicyConfigFromRawConfig(j(t, `[]`), audit)
	if err != nil {
		t.Fatal(err)
	}
	ids := NewOMap[*IgnoreIdRule]()
	ids.Set("CVE-2024-1234", NewIgnoreIdRule("CVE-2024-1234", sp("flagged by mistake"), true, true))
	ids.Set("CVE-2024-1235", NewIgnoreIdRule("CVE-2024-1235", nil, false, true))
	sev := NewOMap[*IgnoreSeverityRule]()
	sev.Set("low", NewIgnoreSeverityRule("low", nil, true, false))
	assertEqualPolicy(t, NewAdvisoriesPolicyConfig(true, AuditFail, ignoreMap(
		rule("acme/abandoned", nil, sp("flagged by mistake"), true, true),
		rule("acme/abandoned2", nil, nil, true, false),
	), ids, sev), got)
}

func TestLegacyAuditIgnore(t *testing.T) {
	cases := []struct {
		name, cfg, audit, block string
	}{
		{"simple array", `{"audit": {"ignore": ["CVE-2024-1234", "CVE-2024-5678"]}}`, `[CVE-2024-1234=>null, CVE-2024-5678=>null]`, `[CVE-2024-1234=>null, CVE-2024-5678=>null]`},
		{"apply audit only", `{"audit": {"ignore": {"CVE-2024-1234": {"apply": "audit", "reason": "Only ignore for auditing"}}}}`, `[CVE-2024-1234=>"Only ignore for auditing"]`, `[]`},
		{"apply block only", `{"audit": {"ignore": {"CVE-2024-1234": {"apply": "block", "reason": "Only ignore for blocking"}}}}`, `[]`, `[CVE-2024-1234=>"Only ignore for blocking"]`},
		{"mixed formats", `{"audit": {"ignore": {"0": "CVE-2024-1234", "CVE-2024-5678": "Simple reason", "CVE-2024-9999": {"apply": "audit", "reason": "Detailed reason"}, "CVE-2024-8888": {"apply": "block"}}}}`,
			`[CVE-2024-1234=>null, CVE-2024-5678=>"Simple reason", CVE-2024-9999=>"Detailed reason"]`,
			`[CVE-2024-1234=>null, CVE-2024-5678=>"Simple reason", CVE-2024-8888=>null]`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			a := mustPolicy(t, newConfig(t, tc.cfg)).Advisories
			if got := flat(a.GetIgnoreListForOperation(OperationAudit)); got != tc.audit {
				t.Fatalf("audit: %s", got)
			}
			if got := flat(a.GetIgnoreListForOperation(OperationBlock)); got != tc.block {
				t.Fatalf("block: %s", got)
			}
		})
	}
}

func TestLegacyAuditIgnoreSeverity(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"audit": {"ignore-severity": ["low", "medium"]}}`)).Advisories
	if flat(a.GetIgnoreSeverityForOperation(OperationAudit)) != `[low=>null, medium=>null]` || flat(a.GetIgnoreSeverityForOperation(OperationBlock)) != `[low=>null, medium=>null]` {
		t.Fatal("simple array mismatch")
	}
	a = mustPolicy(t, newConfig(t, `{"audit": {"ignore-severity": {"low": {"apply": "audit", "reason": "We accept low severity issues"}, "medium": {"apply": "block"}}}}`)).Advisories
	if got := flat(a.GetIgnoreSeverityForOperation(OperationAudit)); got != `[low=>"We accept low severity issues"]` {
		t.Fatal(got)
	}
	if got := flat(a.GetIgnoreSeverityForOperation(OperationBlock)); got != `[medium=>null]` {
		t.Fatal(got)
	}
}

func TestLegacyAuditIgnoreInvalidApplyValue(t *testing.T) {
	clearEnv(t)
	_, err := FromConfig(newConfig(t, `{"audit": {"ignore": {"CVE-2024-1234": {"apply": "invalid"}}}}`))
	expectError(t, err, InvalidArgumentException, "Invalid 'apply' value for 'CVE-2024-1234': invalid. Expected 'audit', 'block', or 'all'.")
}

func TestGetIgnoreListForOperationMergesMultiRuleReasons(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"policy": {"advisories": {"ignore": {"vendor/multi": [{"constraint": "^1.0", "reason": "v1 patched"}, {"constraint": "^2.0", "reason": "v2 mitigated"}]}}}}`)).Advisories
	for _, op := range []string{OperationAudit, OperationBlock} {
		r, ok := a.GetIgnoreListForOperation(op).Get("vendor/multi")
		if !ok || r == nil || !strings.Contains(*r, "v1 patched") || !strings.Contains(*r, "v2 mitigated") {
			t.Fatalf("%s: %s", op, describeStr(r))
		}
	}
}

func TestGetIgnoreListForOperationPrefersConcreteReasonOverNull(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"policy": {"advisories": {"ignore": {"vendor/mixed": [{"constraint": "^1.0"}, {"constraint": "^2.0", "reason": "v2 mitigated"}]}}}}`)).Advisories
	if r, _ := a.GetIgnoreListForOperation(OperationAudit).Get("vendor/mixed"); r == nil || *r != "v2 mitigated" {
		t.Fatal(describeStr(r))
	}
}

func TestWithIgnoreSeverityAddsAuditScopedRulesForNewSeverities(t *testing.T) {
	u := AdvisoriesPolicyConfigDisabled().WithIgnoreSeverity([]string{"low", "medium"})
	if flat(u.GetIgnoreSeverityForOperation(OperationAudit)) != `[low=>null, medium=>null]` || flat(u.GetIgnoreSeverityForOperation(OperationBlock)) != `[]` {
		t.Fatal("mismatch")
	}
}

func TestPolicyAdvisoriesSetIgnoresLegacyAuditAdvisoriesKeys(t *testing.T) {
	audit := j(t, `{"block-insecure": true, "ignore": {"CVE-2024-1234": "should be ignored"}, "ignore-severity": {"low": "should be ignored"}}`)
	a, err := AdvisoriesPolicyConfigFromRawConfig(j(t, `{"advisories": {"block": false, "audit": "report"}}`), audit)
	if err != nil {
		t.Fatal(err)
	}
	if a.Block || a.Audit != AuditReport || a.Ignore.Len() != 0 || a.IgnoreID.Len() != 0 || a.IgnoreSeverity.Len() != 0 {
		t.Fatal("mismatch")
	}
	a, err = AdvisoriesPolicyConfigFromRawConfig(j(t, `{"advisories": false}`), audit)
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, AdvisoriesPolicyConfigDisabled(), a)
}

func TestWithIgnoreSeverityPreservesExistingRulesAndReasons(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"policy": {"advisories": {"ignore-severity": {"low": {"reason": "configured low", "on-block": false, "on-audit": true}}}}}`)).Advisories
	rules := a.WithIgnoreSeverity([]string{"low", "medium"}).GetIgnoreSeverityForOperation(OperationAudit)
	low, _ := rules.Get("low")
	medium, ok := rules.Get("medium")
	if low == nil || *low != "configured low" || !ok || medium != nil {
		t.Fatal("mismatch")
	}
}

// ----------------------------------------------------- AbandonedPolicyConfigTest

func TestAbandonedDefaultConfig(t *testing.T) {
	for _, raw := range []string{`[]`, `{"abandoned": true}`, `{"abandoned": []}`} {
		got, err := AbandonedPolicyConfigFromRawConfig(j(t, raw), j(t, `[]`))
		if err != nil {
			t.Fatal(err)
		}
		assertEqualPolicy(t, NewAbandonedPolicyConfig(false, AuditFail, nil), got)
	}
}

func TestAbandonedFromRawConfig(t *testing.T) {
	got, err := AbandonedPolicyConfigFromRawConfig(j(t, `{"abandoned": {"block": true, "audit": "report", "ignore": {"acme/abandoned": "flagged by mistake", "acme/abandoned2": {"constraint": "1.0"}}}}`), j(t, `[]`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, NewAbandonedPolicyConfig(true, AuditReport, ignoreMap(
		rule("acme/abandoned", nil, sp("flagged by mistake"), true, true),
		rule("acme/abandoned2", semver.MustParseConstraints("1.0"), nil, true, true),
	)), got)
}

func TestAbandonedFromAuditConfig(t *testing.T) {
	got, err := AbandonedPolicyConfigFromRawConfig(j(t, `[]`), j(t, `{"block-abandoned": true, "abandoned": "report", "ignore-abandoned": {"acme/abandoned": "flagged by mistake", "acme/abandoned2": {"apply": "block"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, NewAbandonedPolicyConfig(true, AuditReport, ignoreMap(
		rule("acme/abandoned", nil, sp("flagged by mistake"), true, true),
		rule("acme/abandoned2", nil, nil, true, false),
	)), got)
}

func TestLegacyIgnoreAbandoned(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"audit": {"ignore-abandoned": ["vendor/package1", "vendor/package2"]}}`)).Abandoned
	if flat(a.GetFlatIgnoreForOperation(OperationAudit)) != `[vendor/package1=>null, vendor/package2=>null]` || flat(a.GetFlatIgnoreForOperation(OperationBlock)) != `[vendor/package1=>null, vendor/package2=>null]` {
		t.Fatal("simple mismatch")
	}
	a = mustPolicy(t, newConfig(t, `{"audit": {"ignore-abandoned": {"vendor/package1": {"apply": "audit", "reason": "Report but do not block"}, "vendor/package2": {"apply": "block", "reason": "Block but do not report"}}}}`)).Abandoned
	if flat(a.GetFlatIgnoreForOperation(OperationAudit)) != `[vendor/package1=>"Report but do not block"]` || flat(a.GetFlatIgnoreForOperation(OperationBlock)) != `[vendor/package2=>"Block but do not report"]` {
		t.Fatal("detailed mismatch")
	}
}

func TestPolicyAbandonedIgnoresLegacyAuditAbandonedKeys(t *testing.T) {
	a, err := AbandonedPolicyConfigFromRawConfig(j(t, `{"abandoned": {"block": true, "audit": "report"}}`), j(t, `{"block-abandoned": false, "abandoned": "ignore", "ignore-abandoned": {"acme/abandoned": "should be ignored"}}`))
	if err != nil || !a.Block || a.Audit != AuditReport || a.Ignore.Len() != 0 {
		t.Fatal("mismatch")
	}
	a, err = AbandonedPolicyConfigFromRawConfig(j(t, `{"abandoned": false}`), j(t, `{"block-abandoned": true, "abandoned": "fail", "ignore-abandoned": {"acme/abandoned": "should be ignored"}}`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, AbandonedPolicyConfigDisabled(), a)
}

func TestGetFlatIgnoreForOperationMergesMultiRuleReasons(t *testing.T) {
	clearEnv(t)
	a := mustPolicy(t, newConfig(t, `{"policy": {"abandoned": {"ignore": {"vendor/multi-abandoned": [{"constraint": "^1.0", "reason": "fork ready"}, {"constraint": "^2.0", "reason": "maintained downstream"}]}}}}`)).Abandoned
	r, _ := a.GetFlatIgnoreForOperation(OperationAudit).Get("vendor/multi-abandoned")
	if r == nil || !strings.Contains(*r, "fork ready") || !strings.Contains(*r, "maintained downstream") {
		t.Fatal(describeStr(r))
	}
}

// ------------------------------------------------------ CooldownPolicyConfigTest

func TestParseDuration(t *testing.T) {
	for _, v := range []phpjson.Value{nil, "", int64(0), "0"} {
		if p, err := ParseDuration(v); err != nil || p != nil {
			t.Fatalf("%#v => %v %v", v, p, err)
		}
	}
	ok := map[any]int64{int64(3600): 3600, "86400": 86400, "90 seconds": 90, "1 minute": 60, "30 minutes": 1800, "1 hour": 3600,
		"1 day": 86400, "7 days": 604800, "1 week": 604800, "2 weeks": 1209600, "1week": 604800}
	for in, want := range ok {
		p, err := ParseDuration(in)
		if err != nil || p == nil || *p != want {
			t.Fatalf("%#v => %v %v", in, p, err)
		}
	}
	_, err := ParseDuration("not a duration")
	expectError(t, err, RuntimeException, "Invalid policy.cooldown.period format")
	for _, in := range []string{"tomorrow", "next week", "next month", "today", "midnight", "1 month", "1 year", "2 days ago", "7 days extra", "1.5", "7 Days", "1 WEEK", "  2 hours  ", "   ", "7 days\n"} {
		_, err := ParseDuration(in)
		expectError(t, err, RuntimeException, "Invalid policy.cooldown.period format")
	}
	for _, in := range []phpjson.Value{int64(-3600), "-3600"} {
		_, err := ParseDuration(in)
		expectError(t, err, RuntimeException, "cannot be negative")
	}
}

func TestCooldownDefaultIsInactiveWithoutAge(t *testing.T) {
	for _, raw := range []string{`[]`, `{"cooldown": false}`, `{"cooldown": []}`, `{"cooldown": true}`} {
		clearEnv(t)
		c, err := CooldownPolicyConfigFromRawConfig(j(t, raw))
		if err != nil || c.Period != nil || c.HasCooldown() || c.ShouldBlock(BlockScopeUpdate) {
			t.Fatalf("%s: %v", raw, err)
		}
	}
}

func TestCooldownFromRawConfig(t *testing.T) {
	clearEnv(t)
	got, err := CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": {"period": "7 days", "block": true, "audit": "report", "ignore": {"acme/pkg": "pinned by us", "acme/pkg2": {"constraint": "1.0"}}}}`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, NewCooldownPolicyConfig(true, AuditReport, ignoreMap(
		rule("acme/pkg", nil, sp("pinned by us"), true, true),
		rule("acme/pkg2", semver.MustParseConstraints("1.0"), nil, true, true),
	), i64(604800)), got)

	c, err := CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": {"period": 3600}}`))
	if err != nil || *c.Period != 3600 || !c.HasCooldown() {
		t.Fatal("integer age")
	}
}

func TestCooldownShouldBlock(t *testing.T) {
	active := NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(604800))
	if !active.ShouldBlock(BlockScopeUpdate) || active.ShouldBlock(BlockScopeInstall) {
		t.Fatal("active")
	}
	if NewCooldownPolicyConfig(false, AuditIgnore, nil, i64(604800)).ShouldBlock(BlockScopeUpdate) {
		t.Fatal("block off")
	}
	if NewCooldownPolicyConfig(true, AuditIgnore, nil, nil).ShouldBlock(BlockScopeUpdate) {
		t.Fatal("no age")
	}
}

func TestFormatTimeUntilAvailable(t *testing.T) {
	c := NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(604800))
	now := utc("2026-01-10 00:00:00")
	cases := []struct {
		date  string
		short bool
		want  string
	}{
		{"2026-01-06 04:00:00", false, "3 days 4 hours"},
		{"2026-01-06 04:00:00", true, "3d 4h"},
		{"2026-01-04 00:00:00", false, "1 day"},
		{"2026-01-04 00:00:00", true, "1d"},
		{"2026-01-03 00:05:00", true, "5m"},
		{"2026-01-03 00:00:30", true, "<1m"},
		{"2026-01-03 00:00:30", false, "less than a minute"},
	}
	for _, tc := range cases {
		if got := c.FormatTimeUntilAvailable(utc(tc.date), now, tc.short); got != tc.want {
			t.Fatalf("%s %v: %q", tc.date, tc.short, got)
		}
	}
}

func TestCooldownEnv(t *testing.T) {
	clearEnv(t)
	putEnv(t, "COMPOSER_POLICY_COOLDOWN_PERIOD", "2 days")
	c, err := CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": {"period": "7 days", "ignore": {"acme/pkg": "reason"}}}`))
	if err != nil || *c.Period != 172800 || !c.Ignore.Has("acme/pkg") {
		t.Fatal("env overrides age")
	}
	c, err = CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": false}`))
	if err != nil || *c.Period != 172800 || !c.Block {
		t.Fatal("env enables explicitly disabled cooldown")
	}

	putEnv(t, "COMPOSER_POLICY_COOLDOWN_PERIOD", "0")
	c, err = CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": {"period": "7 days"}}`))
	if err != nil || c.Period != nil || c.HasCooldown() {
		t.Fatal("env disables cooldown")
	}
	c, err = CooldownPolicyConfigFromRawConfig(j(t, `{"cooldown": false}`))
	if err != nil || c.HasCooldown() {
		t.Fatal("zero env age with cooldown disabled")
	}

	putEnv(t, "COMPOSER_POLICY_COOLDOWN_PERIOD", "not a duration")
	_, err = CooldownPolicyConfigFromRawConfig(j(t, `[]`))
	expectError(t, err, RuntimeException, "Invalid value for COMPOSER_POLICY_COOLDOWN_PERIOD")
}

func TestCooldownWithBlockingDisabled(t *testing.T) {
	d := NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(604800)).WithBlockingDisabled()
	if d.Block || *d.Period != 604800 {
		t.Fatal("mismatch")
	}
}

// ---------------------------------------------------- CustomListPolicyConfigTest

func TestCustomListDefaultConfig(t *testing.T) {
	for _, raw := range []string{`[]`, `true`} {
		got, err := CustomListPolicyConfigFromRawConfig("test", j(t, raw))
		if err != nil {
			t.Fatal(err)
		}
		assertEqualPolicy(t, NewCustomListPolicyConfig("test", true, AuditFail, BlockScopeUpdate, nil, nil), got)
	}
}

var shouldBlockMatrix = []struct {
	configured, query string
	expected          bool
}{
	{BlockScopeAll, BlockScopeUpdate, true}, {BlockScopeAll, BlockScopeInstall, true},
	{BlockScopeUpdate, BlockScopeUpdate, true}, {BlockScopeUpdate, BlockScopeInstall, false},
	{BlockScopeInstall, BlockScopeUpdate, false}, {BlockScopeInstall, BlockScopeInstall, true},
}

func TestCustomListShouldBlockHonoursConfiguredBlockScope(t *testing.T) {
	for _, tc := range shouldBlockMatrix {
		if got := NewCustomListPolicyConfig("company-policy", true, AuditFail, tc.configured, nil, nil).ShouldBlock(tc.query); got != tc.expected {
			t.Fatalf("%+v", tc)
		}
	}
}

func TestCustomListFromRawConfig(t *testing.T) {
	got, err := CustomListPolicyConfigFromRawConfig("test", j(t, `{"block": false, "block-scope": "install", "audit": "report",
		"ignore": {"acme/test": "flagged by mistake", "acme/test2": {"constraint": "1.0"}}, "sources": [{"type": "url", "url": "https://example.com"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, NewCustomListPolicyConfig("test", false, AuditReport, BlockScopeInstall, ignoreMap(
		rule("acme/test", nil, sp("flagged by mistake"), true, true),
		rule("acme/test2", semver.MustParseConstraints("1.0"), nil, true, true),
	), []*UrlSource{{ListName: "test", URL: "https://example.com"}}), got)
}

func TestCustomListRejectsNonHttpsSourceUrl(t *testing.T) {
	for _, u := range []string{"http://insecure.example.org/list.json", "ftp://example.org/list.json", "file:///etc/list.json", "//example.org/list.json"} {
		_, err := CustomListPolicyConfigFromRawConfig("company-policy", j(t, `{"sources": [{"type": "url", "url": "`+u+`"}]}`))
		expectError(t, err, RuntimeException, `must start with "https://"`)
	}
}

func TestCustomListRejectsInvalidBlockScope(t *testing.T) {
	for _, bs := range []string{`"installs"`, `true`} {
		_, err := CustomListPolicyConfigFromRawConfig("company-policy", j(t, `{"block-scope": `+bs+`}`))
		expectError(t, err, UnexpectedValueException, `Invalid block-scope for policy list "company-policy"`)
	}
}

// ------------------------------------------------------- MalwarePolicyConfigTest

func TestMalwareDefaultConfig(t *testing.T) {
	for _, raw := range []string{`[]`, `{"malware": true}`, `{"malware": []}`} {
		got, err := MalwarePolicyConfigFromRawConfig(j(t, raw))
		if err != nil {
			t.Fatal(err)
		}
		assertEqualPolicy(t, NewMalwarePolicyConfig(true, AuditFail, BlockScopeAll, nil, nil), got)
	}
}

func TestMalwareFromRawConfig(t *testing.T) {
	got, err := MalwarePolicyConfigFromRawConfig(j(t, `{"malware": {"block": false, "block-scope": "update", "audit": "report",
		"ignore": {"acme/malware": "flagged by mistake", "acme/malware2": {"constraint": "1.0"}}, "ignore-source": ["untrusted"]}}`))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualPolicy(t, NewMalwarePolicyConfig(false, AuditReport, BlockScopeUpdate, ignoreMap(
		rule("acme/malware", nil, sp("flagged by mistake"), true, true),
		rule("acme/malware2", semver.MustParseConstraints("1.0"), nil, true, true),
	), []string{"untrusted"}), got)
}

func TestMalwareShouldBlock(t *testing.T) {
	off := NewMalwarePolicyConfig(false, AuditFail, BlockScopeAll, nil, nil)
	if off.ShouldBlock(BlockScopeUpdate) || off.ShouldBlock(BlockScopeInstall) {
		t.Fatal("off")
	}
	for _, tc := range shouldBlockMatrix {
		if got := NewMalwarePolicyConfig(true, AuditFail, tc.configured, nil, nil).ShouldBlock(tc.query); got != tc.expected {
			t.Fatalf("%+v", tc)
		}
	}
}

// ------------------------------------------------------------- Ignore*RuleTest

func TestParseIgnoreIdMap(t *testing.T) {
	rules, err := ParseIgnoreIdMap(j(t, `[]`))
	if err != nil || rules.Len() != 0 {
		t.Fatal("empty")
	}
	rules, _ = ParseIgnoreIdMap(j(t, `["CVE-123", "GHSA-456"]`))
	if describeIDs(rules) != `CVE-123:{CVE-123 null true true}GHSA-456:{GHSA-456 null true true}` {
		t.Fatal(describeIDs(rules))
	}
	rules, _ = ParseIgnoreIdMap(j(t, `{"CVE-123": "reason", "CVE-456": {"on-block": false, "reason": "other reason"}, "CVE-789": null, "CVE-012": {"on-audit": false}}`))
	if describeIDs(rules) != `CVE-123:{CVE-123 "reason" true true}CVE-456:{CVE-456 "other reason" false true}CVE-789:{CVE-789 null true true}CVE-012:{CVE-012 null true false}` {
		t.Fatal(describeIDs(rules))
	}
	for _, bad := range []string{`[null]`, `[{"id": "CVE-1"}]`, `[true]`, `[42]`, `{"CVE-1": true}`, `{"CVE-1": 42}`} {
		_, err := ParseIgnoreIdMap(j(t, bad))
		expectError(t, err, UnexpectedValueException, "")
	}
}

func TestParseIgnoreSeverityMap(t *testing.T) {
	rules, err := ParseIgnoreSeverityMap(j(t, `[]`))
	if err != nil || rules.Len() != 0 {
		t.Fatal("empty")
	}
	rules, _ = ParseIgnoreSeverityMap(j(t, `["low", "medium"]`))
	if describeSeverities(rules) != `low:{low null true true}medium:{medium null true true}` {
		t.Fatal(describeSeverities(rules))
	}
	rules, _ = ParseIgnoreSeverityMap(j(t, `{"low": "reason", "medium": {"on-block": false, "reason": "other reason"}, "high": null, "critical": {"on-audit": false}}`))
	if describeSeverities(rules) != `low:{low "reason" true true}medium:{medium "other reason" false true}high:{high null true true}critical:{critical null true false}` {
		t.Fatal(describeSeverities(rules))
	}
	for _, bad := range []string{`[null]`, `[{"severity": "low"}]`, `[true]`, `[42]`, `{"low": true}`, `{"low": 42}`} {
		_, err := ParseIgnoreSeverityMap(j(t, bad))
		expectError(t, err, UnexpectedValueException, "")
	}
}

func TestParseIgnoreMap(t *testing.T) {
	rules, err := ParseIgnoreMap(j(t, `[]`))
	if err != nil || rules.Len() != 0 {
		t.Fatal("empty")
	}
	rules, _ = ParseIgnoreMap(j(t, `["vendor/foo", "vendor/bar"]`))
	assertEqualIgnore(t, ignoreMap(rule("vendor/foo", nil, nil, true, true), rule("vendor/bar", nil, nil, true, true)), rules)

	rules, _ = ParseIgnoreMap(j(t, `{"vendor/foo": "reason", "vendor/bar": {"constraint": "^2.0", "on-block": false, "reason": "other reason"}, "vendor/baz": null, "vendor/qux": {"on-audit": false}}`))
	assertEqualIgnore(t, ignoreMap(
		rule("vendor/foo", nil, sp("reason"), true, true),
		rule("vendor/bar", semver.MustParseConstraints("^2.0"), sp("other reason"), false, true),
		rule("vendor/baz", nil, nil, true, true),
		rule("vendor/qux", nil, nil, true, false),
	), rules)
	if r, _ := rules.Get("vendor/bar"); r[0].Constraint.PrettyString() != "^2.0" {
		t.Fatal("pretty")
	}

	rules, _ = ParseIgnoreMap(j(t, `{"vendor/foo": [{"constraint": "^1.0", "reason": "old version"}, {"constraint": "^3.0", "on-block": false}]}`))
	assertEqualIgnore(t, ignoreMap(
		rule("vendor/foo", semver.MustParseConstraints("^1.0"), sp("old version"), true, true),
		rule("vendor/foo", semver.MustParseConstraints("^3.0"), nil, false, true),
	), rules)

	for _, bad := range []string{`[{"package": "vendor/foo", "constraint": "^1.0"}]`, `[true]`, `{"vendor/foo": true}`, `[42]`} {
		_, err := ParseIgnoreMap(j(t, bad))
		expectError(t, err, UnexpectedValueException, "")
	}
}

func assertEqualIgnore(t *testing.T, expected, actual *IgnoreMap) {
	t.Helper()
	if describeIgnore(expected) != describeIgnore(actual) {
		t.Fatalf("expected %s\ngot %s", describeIgnore(expected), describeIgnore(actual))
	}
}

func TestIgnoreUnreachable(t *testing.T) {
	if got := IgnoreUnreachableFromRawAuditConfig(j(t, `{"ignore-unreachable": true}`)); got != (IgnoreUnreachable{true, false, false}) {
		t.Fatal("fromRawAuditConfig")
	}
	iu := IgnoreUnreachable{false, true, false}
	if !iu.ForBlockScope(BlockScopeInstall) || iu.ForBlockScope(BlockScopeUpdate) {
		t.Fatal("forBlockScope 1")
	}
	iu = IgnoreUnreachable{false, false, true}
	if iu.ForBlockScope(BlockScopeInstall) || !iu.ForBlockScope(BlockScopeUpdate) {
		t.Fatal("forBlockScope 2")
	}
	if u, err := iu.With("audit"); err != nil || u != (IgnoreUnreachable{true, false, true}) {
		t.Fatal("with")
	}
	if u, err := IgnoreUnreachableNone().With("audit", "install"); err != nil || u != (IgnoreUnreachable{true, true, false}) {
		t.Fatal("with multiple")
	}
	_, err := IgnoreUnreachableNone().With("not-a-scope")
	expectError(t, err, InvalidArgumentException, "")
	_, err = IgnoreUnreachableNone().With()
	expectError(t, err, InvalidArgumentException, "")
}

var _ = time.Now
