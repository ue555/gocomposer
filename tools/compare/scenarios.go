package main

import "strings"

func project(requireJSON, extra string) string {
	s := `{
    "name": "test/project",
    "type": "project",
    "description": "fixture project",
    "repositories": [
        {"type": "composer", "url": "http://127.0.0.1:8899"},
        {"packagist.org": false}
    ],
    "config": {"secure-http": false},
    "autoload": {"psr-4": {"App\\": "src/"}},
    "require": ` + requireJSON
	if extra != "" {
		s += ",\n    " + extra
	}
	return s + "\n}\n"
}

const appSrc = "<?php\nnamespace App;\n\nclass Kernel {}\n"

func withAdv(s string) string {
	return strings.Replace(s, `"url": "http://127.0.0.1:8899"}`, `"url": "http://127.0.0.1:8899/v2adv"}`, 1)
}

func withConfig(s, extra string) string {
	return strings.Replace(s, `"config": {"secure-http": false}`, `"config": {"secure-http": false, `+extra+`}`, 1)
}

func withV2(s string) string {
	return strings.Replace(s, `"url": "http://127.0.0.1:8899"}`, `"url": "http://127.0.0.1:8899/v2"}`, 1)
}

func allScenarios() []Scenario {
	basicRequire := `{"php": ">=7.4", "acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0"}`
	basicDev := `"require-dev": {"acme/devtools": "^1.0"}`
	files := map[string]string{"src/Kernel.php": appSrc}
	sc := []Scenario{
		{Name: "update-basic", JSON: project(basicRequire, basicDev), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "update-basic-progress", JSON: project(basicRequire, basicDev), Files: files, Command: "update", Probe: false},
		{Name: "update-v2-metadata", JSON: withV2(project(basicRequire, basicDev)), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "install-from-lock", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true},
		{Name: "install-no-dev", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-dev --no-progress", Probe: true},
		{Name: "install-twice", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "install --no-progress", Probe: true},
		{Name: "dry-run-update", JSON: project(basicRequire, basicDev), Files: files, Command: "update --dry-run --no-progress"},
		{Name: "dump-autoload-optimized", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "dump-autoload -o", Probe: true},
		{Name: "dump-autoload-authoritative", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "dump-autoload --classmap-authoritative --no-dev", Probe: true},
		{Name: "replace", JSON: project(`{"acme/old": "^1.0", "acme/new": "^1.0"}`, ""), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "conflict-unresolvable", JSON: project(`{"acme/conflicter": "^1.0", "acme/log": "^2.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "conflict-resolves-lower", JSON: project(`{"acme/conflicter": "^1.0", "acme/log": "*"}`, ""), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "platform-php-too-low", JSON: project(`{"acme/needs-php9": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "missing-extension", JSON: project(`{"acme/needs-ext": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "ignore-platform-reqs", JSON: project(`{"acme/needs-ext": "^1.0", "acme/needs-php9": "^1.0"}`, ""), Files: files, Command: "update --no-progress --ignore-platform-reqs", Probe: true},
		{Name: "config-platform-php", JSON: project(`{"acme/log": "*"}`, `"extra": {"note": "platform pinned"}`), Files: files, Command: "update --no-progress", Probe: true, Env: nil},
		{Name: "not-found-package", JSON: project(`{"acme/does-not-exist": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "stability-beta-flag", JSON: project(`{"acme/beta": "^1.1@beta"}`, ""), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "stability-too-low", JSON: project(`{"acme/beta": "^1.1"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "dev-branch-alias", JSON: project(`{"acme/edge": "1.x-dev"}`, `"minimum-stability": "dev"`), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "root-alias", JSON: project(`{"acme/edge": "dev-main as 1.5.0"}`, ""), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "partial-update", JSON: project(`{"acme/util": "^1.0", "acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress --with acme/util:1.0.0 --with acme/log:1.0.0"}, Command: "update acme/util --no-progress", Probe: true},
		{Name: "partial-update-with-deps", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress --with acme/util:1.0.0 --with acme/log:1.0.0"}, Command: "update acme/util -w --no-progress", Probe: true},
		{Name: "broken-archive", JSON: project(`{"acme/broken": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "checksum-mismatch", JSON: project(`{"acme/badsum": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "download-404", JSON: project(`{"acme/missing": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "path-traversal-archive", JSON: project(`{"acme/evil": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "failure-keeps-vendor", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress", "require --no-update acme/needs-php9:^1.0"}, Command: "update --no-progress", Probe: true},
		{Name: "no-lock-install", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "install --no-progress", Probe: true},
		{Name: "check-platform-reqs-ok", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "check-platform-reqs"},
		{Name: "check-platform-reqs-fail", JSON: project(`{"php": ">=7.4", "acme/needs-ext": "^1.0", "acme/needs-php9": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress --ignore-platform-reqs"}, Command: "check-platform-reqs"},
		{Name: "check-platform-reqs-json", JSON: project(`{"php": ">=7.4", "acme/needs-ext": "^1.0", "acme/needs-php9": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress --ignore-platform-reqs"}, Command: "check-platform-reqs --format=json"},
		{Name: "check-platform-reqs-lock", JSON: project(`{"php": ">=7.4", "acme/needs-ext": "^1.0"}`, basicDev), Files: files, Setup: []string{"update --no-progress --no-install --ignore-platform-reqs"}, Command: "check-platform-reqs --no-dev"},
		{Name: "validate-basic", JSON: project(basicRequire, basicDev), Files: files, Command: "validate"},
		{Name: "validate-clean", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": "MIT",`, 1), Files: files, Command: "validate"},
		{Name: "validate-warnings", JSON: project(`{"php": ">=7.4", "acme/log": "*", "acme/util": "1.0.0", "acme/http": "dev-main#abc123", "acme/cli": "^1.0"}`,
			`"require-dev": {"acme/cli": "^1.0"}, "version": "1.0.0", "license": ["GPL-2.0", "Foo-License", " MIT "], "scripts": {"a": "echo"}, "scripts-descriptions": {"b": "nope"}, "keywords": ["ok", "bad!"], "homepage": "example.com", "authors": [{"name": "A", "email": "not-an-email"}], "provide": {"acme/cli": "1.0"}`), Files: files, Command: "validate"},
		{Name: "validate-strict-warnings", JSON: project(basicRequire, basicDev), Files: files, Command: "validate --strict"},
		{Name: "validate-no-check", JSON: project(`{"acme/log": "*"}`, `"version": "1.0.0"`), Files: files, Command: "validate --no-check-all --no-check-version --no-check-publish"},
		{Name: "validate-missing-file", JSON: project(basicRequire, basicDev), Files: files, Command: "validate nope.json"},
		{Name: "validate-invalid-json", JSON: "{\n    \"name\": \"test/project\",\n    \"require\": {\n        \"acme/log\": \"^1.0\",\n    }\n}\n", Command: "validate"},
		{Name: "validate-duplicate-key", JSON: "{\n    \"name\": \"test/project\",\n    \"license\": \"MIT\",\n    \"description\": \"x\",\n    \"repositories\": [{\"packagist.org\": false}],\n    \"require\": {},\n    \"name\": \"test/project\"\n}\n", Command: "validate"},
		{Name: "validate-lock-stale", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress --no-install"}, After: map[string]string{"composer.json": project(`{"php": ">=7.4", "acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0", "acme/util": "^2.0"}`, basicDev)}, Command: "validate"},
		{Name: "validate-lock-stale-no-check", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress --no-install"}, After: map[string]string{"composer.json": project(`{"php": ">=7.4", "acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0", "acme/util": "^2.0"}`, basicDev)}, Command: "validate --no-check-lock"},
		{Name: "validate-with-dependencies", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "validate -A"},
		{Name: "validate-bad-constraint", JSON: project(`{"acme/log": "~>1"}`, ""), Files: files, Command: "validate"},
		{Name: "show-list", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show"},
		{Name: "show-json", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --format=json"},
		{Name: "show-name-only", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -N"},
		{Name: "show-direct", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -D"},
		{Name: "show-no-dev", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --no-dev"},
		{Name: "show-path", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -P"},
		{Name: "show-package", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show acme/http"},
		{Name: "show-package-json", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show acme/http --format=json"},
		{Name: "show-tree", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -t"},
		{Name: "show-tree-json", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -t --format=json"},
		{Name: "show-package-tree", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show acme/http -t"},
		{Name: "show-locked", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --locked"},
		{Name: "show-self", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -s"},
		{Name: "show-self-json", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -s --format=json"},
		{Name: "show-platform-php", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -p php"},
		{Name: "show-wildcard", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show 'acme/l*'"},
		{Name: "show-not-found", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show acme/nope"},
		{Name: "show-version-arg", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show acme/log 1.1.0"},
		{Name: "show-self-name-only", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --self --name-only"},
		{Name: "show-locked-package", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --locked acme/util"},
		{Name: "show-platform-list-json", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -p --format=json"},
		{Name: "show-bad-format", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show --format=xml"},
		{Name: "show-tree-latest", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "show -t -l"},
		{Name: "show-nothing-installed", JSON: project(basicRequire, basicDev), Files: files, Command: "show"},
		{Name: "depends-basic", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends acme/log"},
		{Name: "depends-recursive", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends acme/log -r"},
		{Name: "depends-tree", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends acme/log -t"},
		{Name: "depends-php", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends php"},
		{Name: "depends-notfound", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends acme/nope"},
		{Name: "depends-locked", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "depends --locked acme/util"},
		{Name: "prohibits-basic", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "prohibits acme/log 3.0"},
		{Name: "prohibits-installed", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "prohibits acme/log 1.1.0"},
		{Name: "prohibits-php", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "prohibits php 9.0"},
		{Name: "prohibits-tree", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "prohibits acme/util 2.0 -t"},
		{Name: "why-alias", JSON: project(basicRequire, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "why acme/util"},
		{Name: "depends-nothing-installed", JSON: project(basicRequire, basicDev), Files: files, Command: "depends acme/log"},
		{Name: "licenses-text", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses"},
		{Name: "licenses-json", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses --format=json"},
		{Name: "licenses-summary", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses --format=summary"},
		{Name: "licenses-no-dev", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses --no-dev"},
		{Name: "licenses-locked", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses --locked"},
		{Name: "licenses-bad-format", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "licenses --format=xml"},
		{Name: "suggests-default", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests"},
		{Name: "suggests-all", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests --all"},
		{Name: "suggests-list", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests --list"},
		{Name: "suggests-by-suggestion", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests --by-suggestion"},
		{Name: "suggests-both", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests --by-package --by-suggestion"},
		{Name: "suggests-missing", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests acme/nope"},
		{Name: "suggests-filter", JSON: strings.Replace(project(basicRequire, basicDev), `"type": "project",`, `"type": "project", "license": ["MIT", "GPL-2.0-or-later"], "suggest": {"acme/devtools": "for <dev> work", "ext-foo": "speed"},`, 1), Files: files, Setup: []string{"update --no-progress"}, Command: "suggests acme/util"},
		{Name: "validate-publish-errors", JSON: strings.Replace(strings.Replace(project(`{"acme/log": "^1.0"}`, `"bogus-key": 1`), `"name": "test/project"`, `"name": "Test/MyProject"`, 1), `"description": "fixture project",`, "", 1), Files: files, Command: "validate"},
		{Name: "validate-lax-fail", JSON: project(`{"acme/log": ["^1.0"]}`, ""), Files: files, Command: "validate", KnownDiff: "Composer crashes with a PHP TypeError in ConfigValidator (Preg::isMatch on an array); gocomposer reports the LAX schema error instead. Both exit 1."},
		{Name: "validate-lax-fail-update", JSON: project(`{"acme/log": "^1.0"}`, `"autoload-dev": "nope"`), Files: files, Command: "update --no-progress"},
		{Name: "composer-auth-invalid", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Env: map[string]string{"COMPOSER_AUTH": `{"http-basic": {"example.org": {"username": 5}}}`}, Command: "validate", KnownDiff: "Composer prints the schema warnings and then dies with a PHP fatal TypeError (BaseIO::checkAndSetAuthentication with an int username, exit 255); gocomposer prints the same warnings and ignores the non-string credential."},
		{Name: "composer-auth-malformed", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Env: map[string]string{"COMPOSER_AUTH": `{nope`}, Command: "validate"},
		{Name: "policy-advisory-blocks", JSON: withAdv(project(`{"acme/http": "^1.0", "acme/log-impl": "^1.0"}`, "")), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "policy-advisory-blocks-vv", JSON: withAdv(project(`{"acme/http": "^1.0", "acme/log-impl": "^1.0"}`, "")), Files: files, Command: "update --no-progress -vv"},
		{Name: "policy-advisory-only-vulnerable", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), Files: files, Command: "update --no-progress"},
		{Name: "policy-no-blocking-audit", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, "")), Files: files, Command: "update --no-progress --no-blocking"},
		{Name: "policy-env-no-blocking", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), Files: files, Env: map[string]string{"COMPOSER_NO_BLOCKING": "1"}, Command: "update --no-progress"},
		{Name: "policy-audit-install", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "install --no-progress --audit"},
		{Name: "policy-audit-install-table", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "install --no-progress --audit --audit-format=table"},
		{Name: "policy-audit-install-json", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "install --no-progress --audit --audit-format=json"},
		{Name: "policy-audit-plain", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "install --no-progress --audit --audit-format=plain"},
		{Name: "policy-abandoned", JSON: withAdv(project(`{"acme/old": "^1.0"}`, "")), Files: files, Command: "update --no-progress"},
		{Name: "policy-abandoned-ignore", JSON: withConfig(withAdv(project(`{"acme/old": "^1.0"}`, "")), `"audit": {"abandoned": "ignore"}`), Files: files, Command: "update --no-progress"},
		{Name: "policy-ignore-advisory", JSON: withConfig(withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), `"audit": {"ignore": ["PKSA-test-0001"]}`), Files: files, Command: "update --no-progress"},
		{Name: "policy-disabled", JSON: withConfig(withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), `"policy": false`), Files: files, Command: "update --no-progress"},
		{Name: "policy-cooldown-all", JSON: withConfig(project(`{"acme/log": "^1.0"}`, ""), `"policy": {"cooldown": {"period": "5000 days"}}`), Files: files, Command: "update --no-progress"},
		{Name: "policy-cooldown-none", JSON: withConfig(project(`{"acme/log": "^1.0"}`, ""), `"policy": {"cooldown": {"period": "1 day"}}`), Files: files, Command: "update --no-progress"},
		{Name: "policy-cooldown-no-cooldown", JSON: withConfig(project(`{"acme/log": "^1.0"}`, ""), `"policy": {"cooldown": {"period": "5000 days"}}`), Files: files, Command: "update --no-progress --no-cooldown"},
		{Name: "policy-locked-vulnerable-install", JSON: withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "install --no-progress"},
		{Name: "policy-partial-update-vulnerable", JSON: withAdv(project(`{"acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/util": "^1.0"}`, "")), Files: files, Setup: []string{"update --no-progress --no-blocking --no-install"}, Command: "update acme/util --no-progress"},
		{Name: "require-guess-version", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "require acme/util --no-progress", Probe: true},
		{Name: "require-explicit", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "require acme/util:^1.0 --no-progress", Probe: true},
		{Name: "require-dev", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "require --dev acme/devtools --no-progress"},
		{Name: "require-no-update", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "require acme/http:^1.0 acme/log-impl --no-update"},
		{Name: "require-dry-run", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "require acme/util --dry-run --no-progress"},
		{Name: "require-first-no-lock", JSON: project(`{}`, ""), Files: files, Command: "require acme/log --no-progress", Probe: true},
		{Name: "require-no-composer-json", NoJSON: true, Files: files, Command: "require acme/log --no-progress"},
		{Name: "require-not-found", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "require acme/nope --no-progress"},
		{Name: "require-min-stability", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "require acme/edge --no-progress"},
		{Name: "require-platform-incompatible", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "require acme/needs-php9 --no-progress"},
		{Name: "require-fails-reverts", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "require acme/needs-php9:^1.0 --no-progress"},
		{Name: "require-move-from-dev", JSON: project(`{"acme/log": "^1.0"}`, `"require-dev": {"acme/util": "^1.0"}`), Files: files, Setup: []string{"update --no-progress"}, Command: "require acme/util:^1.1 --no-progress"},
		{Name: "require-sort-packages", JSON: project(`{"acme/util": "^1.0", "php": ">=7.4"}`, ""), Files: files, Command: "require acme/log:^1.0 ext-json:* --sort-packages --no-update"},
		{Name: "require-too-strict", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Command: "require acme/log:1 --no-update"},
		{Name: "require-self", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Command: "require test/project:^1.0 --no-update"},
		{Name: "require-ext-missing", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Command: "require ext-imaginary --no-update"},
		{Name: "require-alias-arg", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Command: "require acme/edge:dev-main as 1.0.0 --no-update"},
		{Name: "require-fixed-library", JSON: strings.Replace(project(`{"acme/util": "^1.0"}`, ""), `"type": "project",`, `"type": "library",`, 1), Files: files, Command: "require acme/log --fixed --no-update"},
		{Name: "require-fixed-project", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Command: "require acme/log --fixed --no-progress"},
		{Name: "require-update-with-deps", JSON: project(`{"acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress --with acme/util:1.0.0 --with acme/log:1.0.0"}, Command: "require acme/util:^1.1 -w --no-progress"},
		{Name: "require-crlf-file", JSON: strings.ReplaceAll(project(`{"acme/util": "^1.0"}`, ""), "\n", "\r\n"), Files: files, Command: "require acme/log:^1.0 --no-update"},
		{Name: "require-tabs-file", JSON: strings.ReplaceAll(project(`{"acme/util": "^1.0"}`, ""), "    ", "\t"), Files: files, Command: "require acme/log:^1.0 --no-update"},
		{Name: "remove-basic", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/util --no-progress", Probe: true},
		{Name: "remove-dev", JSON: project(`{"acme/log": "^1.0"}`, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "remove --dev acme/devtools --no-progress"},
		{Name: "remove-wrong-section", JSON: project(`{"acme/log": "^1.0"}`, basicDev), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/devtools --no-progress"},
		{Name: "remove-not-required", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/nope --no-progress"},
		{Name: "remove-wildcard", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "remove 'acme/u*' --no-progress"},
		{Name: "remove-no-update", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Command: "remove acme/util --no-update"},
		{Name: "remove-dry-run", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/util --dry-run --no-progress"},
		{Name: "remove-still-required", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/log --no-progress"},
		{Name: "remove-unused", JSON: project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Setup: []string{"update --no-progress", "remove acme/util --no-update"}, Command: "remove --unused --no-progress"},
		{Name: "remove-no-args", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "remove"},
		{Name: "remove-mixed-case", JSON: project(`{"Acme/Log": "^1.0", "acme/util": "^1.0"}`, ""), Files: files, Command: "remove ACME/LOG --no-update"},
		{Name: "remove-allow-plugins", JSON: withConfig(project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, ""), `"allow-plugins": {"acme/util": true, "acme/other": false}`), Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/util --no-progress"},
	}
	sc = append(sc, configScenarios()...)
	for _, f := range extraScenarios {
		sc = append(sc, f()...)
	}
	// config.platform scenario
	for i := range sc {
		if sc[i].Name == "config-platform-php" {
			sc[i].JSON = strings.Replace(sc[i].JSON, `"config": {"secure-http": false}`, `"config": {"secure-http": false, "platform": {"php": "7.4.0"}}`, 1)
		}
	}
	return sc
}

// extraScenarios is appended to from init() in per-feature scenario files.
var extraScenarios []func() []Scenario
