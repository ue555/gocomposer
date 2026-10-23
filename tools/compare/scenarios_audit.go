package main

import "strings"

func init() { extraScenarios = append(extraScenarios, auditScenarios) }

// auditScenarios covers the `audit` command (AuditCommand of Composer
// f83588392df43d4fd027d84d26de459946ef4e80) against the /v2adv (metadata
// advisories) and /v2api (security-advisories api-url) fixture variants.
func auditScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	vuln := `{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`
	clean := `{"acme/log-impl": "^1.0"}`
	adv := func(req, extra string) string { return withAdv(project(req, extra)) }
	api := func(req, extra string) string {
		return strings.Replace(project(req, extra), `"url": "http://127.0.0.1:8899"}`, `"url": "http://127.0.0.1:8899/v2api"}`, 1)
	}
	install := []string{"update --no-progress --no-blocking"}
	lockOnly := []string{"update --no-progress --no-blocking --no-install"}
	// an unreachable second repository, added after the install
	unreachable := func(s, url string) string {
		return strings.Replace(s, `{"packagist.org": false}`, `{"type": "composer", "url": "`+url+`"},
        {"packagist.org": false}`, 1)
	}

	type a struct {
		name, json, cmd string
		setup           []string
		env             map[string]string
		after           map[string]string
	}
	list := []a{
		// --- basics ---
		{name: "clean", json: adv(clean, ""), cmd: "audit", setup: install},
		{name: "clean-json", json: adv(clean, ""), cmd: "audit --format=json", setup: install},
		{name: "clean-summary", json: adv(clean, ""), cmd: "audit -f summary", setup: install},
		{name: "no-advisory-repo", json: project(vuln, ""), cmd: "audit", setup: install},
		{name: "table", json: adv(vuln, ""), cmd: "audit", setup: install},
		{name: "plain", json: adv(vuln, ""), cmd: "audit --format=plain", setup: install},
		{name: "json", json: adv(vuln, ""), cmd: "audit --format=json", setup: install},
		{name: "summary", json: adv(vuln, ""), cmd: "audit --format=summary", setup: install},
		{name: "format-short", json: adv(vuln, ""), cmd: "audit -f plain", setup: install},
		{name: "format-invalid", json: adv(vuln, ""), cmd: "audit --format=xml", setup: install},
		{name: "format-missing-value", json: adv(vuln, ""), cmd: "audit --format", setup: install},
		{name: "verbose", json: adv(vuln, ""), cmd: "audit -v", setup: install},
		// --- package selection ---
		{name: "no-vendor", json: adv(vuln, ""), cmd: "audit", setup: lockOnly},
		{name: "no-vendor-no-dev-only-dev", json: adv(`{}`, `"require-dev": {"acme/log-impl": "^1.0"}`), cmd: "audit --no-dev", setup: lockOnly},
		{name: "no-packages", json: adv(`{}`, ""), cmd: "audit"},
		{name: "no-lock-no-vendor", json: adv(vuln, ""), cmd: "audit"},
		{name: "locked", json: adv(vuln, ""), cmd: "audit --locked", setup: lockOnly},
		{name: "locked-json", json: adv(vuln, ""), cmd: "audit --locked --format=json", setup: lockOnly},
		{name: "locked-no-lock", json: adv(vuln, ""), cmd: "audit --locked"},
		{name: "no-dev", json: adv(clean, `"require-dev": {"acme/http": "^1.2", "acme/util": "1.0.0"}`), cmd: "audit --no-dev", setup: install},
		{name: "with-dev", json: adv(clean, `"require-dev": {"acme/http": "^1.2", "acme/util": "1.0.0"}`), cmd: "audit", setup: install},
		{name: "locked-no-dev", json: adv(clean, `"require-dev": {"acme/http": "^1.2", "acme/util": "1.0.0"}`), cmd: "audit --locked --no-dev", setup: lockOnly},
		{name: "direct", json: adv(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, ""), cmd: "audit --direct", setup: install},
		{name: "direct-short", json: adv(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, ""), cmd: "audit -D --format=plain", setup: install},
		// --- abandoned ---
		{name: "abandoned-default", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit", setup: install},
		{name: "abandoned-plain", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --format=plain", setup: install},
		{name: "abandoned-json", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --format=json", setup: install},
		{name: "abandoned-summary", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --format=summary", setup: install},
		{name: "abandoned-ignore", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --abandoned=ignore", setup: install},
		{name: "abandoned-report", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --abandoned=report", setup: install},
		{name: "abandoned-fail", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --abandoned=fail", setup: install},
		{name: "abandoned-invalid", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit --abandoned=nope", setup: install},
		{name: "abandoned-config-report", json: withConfig(adv(`{"acme/old": "^1.0"}`, ""), `"audit": {"abandoned": "report"}`), cmd: "audit", setup: install},
		{name: "abandoned-config-ignore", json: withConfig(adv(`{"acme/old": "^1.0"}`, ""), `"audit": {"abandoned": "ignore"}`), cmd: "audit --format=json", setup: install},
		{name: "abandoned-option-overrides-config", json: withConfig(adv(`{"acme/old": "^1.0"}`, ""), `"audit": {"abandoned": "ignore"}`), cmd: "audit --abandoned=fail", setup: install},
		{name: "abandoned-env-report", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit", setup: install, env: map[string]string{"COMPOSER_AUDIT_ABANDONED": "report"}},
		// (COMPOSER_AUDIT_ABANDONED=bogus would also break the setup install,
		// so the invalid value is set in composer.json after the install)
		{name: "abandoned-config-invalid", json: adv(`{"acme/old": "^1.0"}`, ""), cmd: "audit", setup: install,
			after: map[string]string{"composer.json": withConfig(adv(`{"acme/old": "^1.0"}`, ""), `"audit": {"abandoned": "bogus"}`)}},
		{name: "abandoned-env-overrides-config", json: withConfig(adv(`{"acme/old": "^1.0"}`, ""), `"audit": {"abandoned": "ignore"}`), cmd: "audit", setup: install, env: map[string]string{"COMPOSER_AUDIT_ABANDONED": "fail"}},
		{name: "abandoned-and-vulnerable", json: adv(`{"acme/old": "^1.0", "acme/http": "^1.2", "acme/log-impl": "^1.0"}`, ""), cmd: "audit --format=plain", setup: install},
		// --- ignore lists ---
		{name: "ignore-id", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["PKSA-test-0001"]}`), cmd: "audit", setup: install},
		{name: "ignore-id-reason", json: withConfig(adv(vuln, ""), `"audit": {"ignore": {"PKSA-test-0001": "accepted risk"}}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-id-json", json: withConfig(adv(vuln, ""), `"audit": {"ignore": {"PKSA-test-0001": "accepted risk"}}`), cmd: "audit --format=json", setup: install},
		{name: "ignore-cve", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["CVE-2024-0001"]}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-remote-id", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["GHSA-test-0001"]}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-package", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["acme/util"]}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-all", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["PKSA-test-0001", "PKSA-test-0002"]}`), cmd: "audit", setup: install},
		{name: "ignore-summary-cve", json: withConfig(adv(vuln, ""), `"audit": {"ignore": ["CVE-2024-0001"]}`), cmd: "audit --format=summary", setup: install},
		{name: "ignore-apply-block-only", json: withConfig(adv(vuln, ""), `"audit": {"ignore": {"PKSA-test-0001": {"apply": "block"}}}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-apply-audit", json: withConfig(adv(vuln, ""), `"audit": {"ignore": {"PKSA-test-0001": {"apply": "audit", "reason": "r"}}}`), cmd: "audit --format=plain", setup: install},
		{name: "ignore-severity-option", json: adv(vuln, ""), cmd: "audit --ignore-severity=high --format=plain", setup: install},
		{name: "ignore-severity-option-multi", json: adv(vuln, ""), cmd: "audit --ignore-severity=low --ignore-severity=high", setup: install},
		{name: "ignore-severity-option-other", json: adv(vuln, ""), cmd: "audit --ignore-severity=low --format=json", setup: install},
		{name: "ignore-severity-config", json: withConfig(adv(vuln, ""), `"audit": {"ignore-severity": ["high"]}`), cmd: "audit --format=json", setup: install},
		{name: "ignore-severity-config-reason", json: withConfig(adv(vuln, ""), `"audit": {"ignore-severity": {"high": "too noisy"}}`), cmd: "audit --format=plain", setup: install},
		{name: "policy-advisories-ignore", json: withConfig(adv(vuln, ""), `"policy": {"advisories": {"ignore": {"acme/http": "pkg reason"}}}`), cmd: "audit --format=plain", setup: install},
		{name: "policy-disabled", json: withConfig(adv(vuln, ""), `"policy": false`), cmd: "audit --format=plain", setup: install},
		// --- api-url advisories ---
		{name: "api-table", json: api(vuln, ""), cmd: "audit", setup: install},
		{name: "api-json", json: api(vuln, ""), cmd: "audit --format=json", setup: install},
		{name: "api-summary", json: api(vuln, ""), cmd: "audit --format=summary", setup: install},
		{name: "api-locked", json: api(vuln, ""), cmd: "audit --locked --format=plain", setup: lockOnly},
		{name: "api-unrequested-names", json: api(clean, ""), cmd: "audit", setup: install},
		{name: "api-ignore-severity", json: api(vuln, ""), cmd: "audit --ignore-severity=high --format=plain", setup: install},
		// --- unreachable repositories ---
		{name: "unreachable-404", json: adv(vuln, ""), cmd: "audit --format=plain", setup: install,
			after: map[string]string{"composer.json": unreachable(adv(vuln, ""), "http://127.0.0.1:8899/missing-repo")}},
		{name: "unreachable-404-ignored", json: adv(vuln, ""), cmd: "audit --format=plain --ignore-unreachable", setup: install,
			after: map[string]string{"composer.json": unreachable(adv(vuln, ""), "http://127.0.0.1:8899/missing-repo")}},
		{name: "unreachable-404-ignored-json", json: adv(vuln, ""), cmd: "audit --format=json --ignore-unreachable", setup: install,
			after: map[string]string{"composer.json": unreachable(adv(vuln, ""), "http://127.0.0.1:8899/missing-repo")}},
		{name: "unreachable-refused-ignored", json: adv(clean, ""), cmd: "audit --ignore-unreachable", setup: install,
			after: map[string]string{"composer.json": unreachable(adv(clean, ""), "http://127.0.0.1:1")}},
		{name: "unreachable-config-ignored", json: adv(clean, ""), cmd: "audit", setup: install,
			after: map[string]string{"composer.json": withConfig(unreachable(adv(clean, ""), "http://127.0.0.1:8899/missing-repo"), `"policy": {"ignore-unreachable": ["audit"]}`)}},
		{name: "unreachable-api-ignored", json: api(clean, ""), cmd: "audit --ignore-unreachable", setup: install,
			after: map[string]string{"composer.json": unreachable(api(clean, ""), "http://127.0.0.1:8899/missing-repo")}},
	}
	out := make([]Scenario, 0, len(list))
	for _, s := range list {
		out = append(out, Scenario{Name: "audit-" + s.name, JSON: s.json, Files: files, Setup: s.setup, Command: s.cmd, Env: s.env, After: s.after})
	}
	return out
}
