package main

import "strings"

// show --latest/--outdated/--all/--available and outdated scenarios
// (ShowCommand / OutdatedCommand of Composer f83588392df43d4fd027d84d26de459946ef4e80).
// They use the latestfix/* fixture packages so the acme/* scenarios are not affected.

func init() { extraScenarios = append(extraScenarios, showLatestScenarios) }

func showLatestScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	req := `{"php": ">=7.4", "latestfix/core": "^1.0", "latestfix/zero": "^0.1", "latestfix/stable": "^1.0", "latestfix/aband": "^1.0", "latestfix/gone": "^1.0", "latestfix/platneed": "^1.0"}`
	dev := `"require-dev": {"latestfix/devtool": "^1.0"}`
	base := project(req, dev)
	setup := []string{"update --no-progress --with latestfix/core:1.0.0 --with latestfix/dep:1.0.0 --with latestfix/zero:0.1.0 --with latestfix/aband:1.0.0 --with latestfix/devtool:1.0.0"}
	lockOnly := []string{setup[0] + " --no-install"}
	// every package at its newest installable version
	fresh := project(`{"latestfix/stable": "^1.0", "latestfix/zero": "^0.2"}`, "")
	// all releases date from 2024-01-02, so a long cooldown withholds them all
	// and a short one none
	platOnly := project(`{"latestfix/platonly": "^1.0", "latestfix/stable": "^1.0"}`, "")
	edge := project(`{"acme/edge": "dev-main", "latestfix/core": "^1.0"}`, `"minimum-stability": "dev", "prefer-stable": true`)
	v2 := withV2(base)
	// fall back to the package "time" (2024-01-02) so the cooldown applies
	cdTime := strings.Replace(base, `"url": "http://127.0.0.1:8899"}`, `"url": "http://127.0.0.1:8899", "require-published-time": false}`, 1)
	cdConfig := withConfig(cdTime, `"policy": {"cooldown": {"period": "5000 days"}}`)
	type c struct {
		name, json, cmd string
		setup           []string
		nojson          bool
		after           string
	}
	list := []c{
		// --- show --latest / --outdated listings ---
		{name: "show-latest", cmd: "show --latest"},
		{name: "show-latest-short", cmd: "show -l"},
		{name: "show-latest-outdated", cmd: "show -o"},
		{name: "show-latest-outdated-direct", cmd: "show -D -o"},
		{name: "show-latest-direct", cmd: "show -l -D"},
		{name: "show-latest-strict", cmd: "show -o --strict"},
		{name: "show-latest-strict-latest-only", cmd: "show -l --strict"},
		{name: "show-latest-strict-ignored", cmd: "show -o --strict --ignore 'latestfix/*'"},
		{name: "show-latest-strict-fresh", json: fresh, setup: []string{"update --no-progress"}, cmd: "show -o --strict"},
		{name: "show-latest-fresh", json: fresh, setup: []string{"update --no-progress"}, cmd: "show -l"},
		{name: "show-latest-strict-no-latest", cmd: "show --strict"},
		{name: "show-latest-major-only", cmd: "show -o -M"},
		{name: "show-latest-major-only-latest", cmd: "show -l --major-only"},
		{name: "show-latest-minor-only", cmd: "show -l -m"},
		{name: "show-latest-patch-only", cmd: "show -l --patch-only"},
		{name: "show-latest-patch-only-outdated", cmd: "show -o --patch-only"},
		{name: "show-latest-only-conflict", cmd: "show -l -M -m"},
		{name: "show-latest-json", cmd: "show -l --format=json"},
		{name: "show-latest-outdated-json", cmd: "show -o --format=json"},
		{name: "show-latest-outdated-json-strict", cmd: "show -o --format=json --strict"},
		{name: "show-latest-name-only", cmd: "show -l -N"},
		{name: "show-latest-name-only-json", cmd: "show -o -N --format=json"},
		{name: "show-latest-path", cmd: "show -o -P"},
		{name: "show-latest-wildcard", cmd: "show -l 'latestfix/c*'"},
		{name: "show-latest-ignore", cmd: "show -o --ignore latestfix/core --ignore 'latestfix/z*'"},
		{name: "show-latest-ignore-without-outdated", cmd: "show --ignore latestfix/core"},
		{name: "show-latest-no-dev", cmd: "show -o --no-dev"},
		{name: "show-latest-sort-by-age", cmd: "show -l -A"},
		{name: "show-latest-ansi", cmd: "show -l --ansi"},
		{name: "show-latest-outdated-ansi", cmd: "show -o --ansi"},
		{name: "show-latest-verbose", cmd: "show -l -v"},
		{name: "show-latest-ignore-platform-reqs", cmd: "show -l --ignore-platform-reqs"},
		{name: "show-latest-ignore-platform-req", cmd: "show -o --ignore-platform-req=php"},
		{name: "show-latest-platform-rejected", json: platOnly, setup: []string{"update --no-progress --ignore-platform-reqs"}, cmd: "show -l"},
		{name: "show-latest-platform-rejected-verbose", json: platOnly, setup: []string{"update --no-progress --ignore-platform-reqs"}, cmd: "show -o -v"},
		{name: "show-latest-locked", setup: lockOnly, cmd: "show -l --locked"},
		{name: "show-latest-locked-outdated", setup: lockOnly, cmd: "show -o --locked --no-dev"},
		{name: "show-latest-not-installed", setup: lockOnly, cmd: "show -l"},
		{name: "show-latest-tree", cmd: "show -l -t"},
		{name: "show-latest-platform", cmd: "show -l -p php"},
		{name: "show-latest-v2-metadata", json: v2, cmd: "show -o"},
		{name: "show-latest-dev-branch", json: edge, setup: []string{"update --no-progress --with latestfix/core:1.0.0"}, cmd: "show -l"},
		{name: "show-latest-dev-branch-major", json: edge, setup: []string{"update --no-progress --with latestfix/core:1.0.0"}, cmd: "show -l -M"},
		{name: "show-latest-dev-branch-package", json: edge, setup: []string{"update --no-progress --with latestfix/core:1.0.0"}, cmd: "show -l acme/edge"},
		{name: "show-latest-no-composer", nojson: true, cmd: "show -l -p --format=json"},
		// --- cooldown period ---
		{name: "show-latest-cooldown-long", cmd: "show -l --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-long-json", cmd: "show -o --cooldown-period='5000 days' --format=json"},
		{name: "show-latest-cooldown-short", cmd: "show -o --cooldown-period=60"},
		{name: "show-latest-cooldown-zero", cmd: "show -o --cooldown-period=0"},
		{name: "show-latest-cooldown-invalid", cmd: "show -l --cooldown-period=soon"},
		{name: "show-latest-cooldown-package", cmd: "show -l latestfix/core --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time", json: cdTime, cmd: "show -l --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-json", json: cdTime, cmd: "show -o --format=json --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-strict", json: cdTime, cmd: "show -o --strict --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-ansi", json: cdTime, cmd: "show -o --ansi --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-verbose", json: cdTime, cmd: "show -l -v --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-passed", json: cdTime, cmd: "show -o --cooldown-period='1000 days'"},
		{name: "show-latest-cooldown-time-package", json: cdTime, cmd: "show -l latestfix/core --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-time-package-json", json: cdTime, cmd: "show -l latestfix/zero --format=json --cooldown-period='5000 days'"},
		{name: "show-latest-cooldown-config", json: cdTime, after: cdConfig, cmd: "show -o"},
		{name: "show-latest-cooldown-config-json", json: cdTime, after: cdConfig, cmd: "show -l --format=json"},
		{name: "outdated-cooldown-config", json: cdTime, after: cdConfig, cmd: "outdated --strict"},
		{name: "outdated-cooldown-config-zero", json: cdTime, after: cdConfig, cmd: "outdated --cooldown-period=0"},
		// --- single package ---
		{name: "show-latest-package", cmd: "show -l latestfix/core"},
		{name: "show-latest-package-json", cmd: "show -l latestfix/core --format=json"},
		{name: "show-latest-package-uptodate", cmd: "show -l latestfix/stable"},
		{name: "show-latest-package-abandoned", cmd: "show -l latestfix/aband"},
		{name: "show-latest-package-abandoned-json", cmd: "show -l latestfix/aband --format=json"},
		{name: "show-latest-package-gone", cmd: "show -o --strict latestfix/gone"},
		{name: "show-latest-package-strict", cmd: "show -o --strict latestfix/core"},
		{name: "show-latest-package-strict-uptodate", cmd: "show -o --strict latestfix/stable"},
		{name: "show-latest-package-major", cmd: "show -l -M latestfix/zero"},
		{name: "show-latest-package-patch", cmd: "show -l --patch-only latestfix/core"},
		{name: "show-latest-package-path", cmd: "show -o --strict -P latestfix/core"},
		{name: "show-latest-package-direct-transitive", cmd: "show -l -D latestfix/dep"},
		{name: "show-latest-package-platform-rejected", cmd: "show -l latestfix/platneed"},
		// --- --all / --available ---
		{name: "show-latest-all", cmd: "show --all"},
		{name: "show-latest-all-json", cmd: "show --all --format=json"},
		{name: "show-latest-all-name-only", cmd: "show --all -N"},
		{name: "show-latest-all-wildcard", cmd: "show --all 'latestfix/*'"},
		{name: "show-latest-all-package", cmd: "show --all latestfix/core"},
		{name: "show-latest-all-package-version", cmd: "show --all latestfix/core 2.0.0"},
		{name: "show-latest-all-no-lock", setup: []string{}, cmd: "show --all 'latestfix/*'"},
		{name: "show-latest-all-latest", cmd: "show --all -l 'latestfix/*'"},
		{name: "show-latest-all-strict", cmd: "show --all --strict"},
		{name: "show-latest-available", cmd: "show -a"},
		{name: "show-latest-available-json", cmd: "show -a --format=json"},
		{name: "show-latest-available-wildcard", cmd: "show -a 'latestfix/z*'"},
		{name: "show-latest-available-package", cmd: "show -a latestfix/core"},
		{name: "show-latest-available-package-json", cmd: "show -a latestfix/zero --format=json"},
		{name: "show-latest-available-not-found", cmd: "show -a latestfix/nope"},
		{name: "show-latest-available-v2", json: v2, cmd: "show -a 'latestfix/*'"},
		{name: "show-latest-available-direct", cmd: "show -a -D"},
		{name: "show-latest-available-tree", cmd: "show -a -t"},
		{name: "show-latest-available-latest", cmd: "show -a -l"},
		// --- outdated command ---
		{name: "outdated-default", cmd: "outdated"},
		{name: "outdated-all", cmd: "outdated -a"},
		{name: "outdated-direct", cmd: "outdated -D"},
		{name: "outdated-strict", cmd: "outdated --strict"},
		{name: "outdated-strict-fresh", json: fresh, setup: []string{"update --no-progress"}, cmd: "outdated --strict"},
		{name: "outdated-major-only", cmd: "outdated -M"},
		{name: "outdated-minor-only", cmd: "outdated -m"},
		{name: "outdated-patch-only", cmd: "outdated -p"},
		{name: "outdated-only-conflict", cmd: "outdated -M -p"},
		{name: "outdated-json", cmd: "outdated --format=json"},
		{name: "outdated-all-json", cmd: "outdated -a -f json"},
		{name: "outdated-bad-format", cmd: "outdated --format=xml"},
		{name: "outdated-package", cmd: "outdated latestfix/core"},
		{name: "outdated-package-strict", cmd: "outdated --strict latestfix/zero"},
		{name: "outdated-package-not-found", cmd: "outdated latestfix/nope"},
		{name: "outdated-wildcard", cmd: "outdated 'latestfix/[cz]*'"},
		{name: "outdated-locked", setup: lockOnly, cmd: "outdated --locked"},
		{name: "outdated-no-dev", cmd: "outdated --no-dev"},
		{name: "outdated-ignore", cmd: "outdated --ignore latestfix/core --ignore latestfix/devtool"},
		{name: "outdated-ignore-platform-reqs", cmd: "outdated --ignore-platform-reqs"},
		{name: "outdated-ignore-platform-req", cmd: "outdated --ignore-platform-req=php+"},
		{name: "outdated-sort-by-age", cmd: "outdated -A"},
		{name: "outdated-cooldown", cmd: "outdated --cooldown-period='5000 days'"},
		{name: "outdated-ansi", cmd: "outdated --ansi"},
		{name: "outdated-not-installed", setup: lockOnly, cmd: "outdated"},
		{name: "outdated-no-composer", nojson: true, cmd: "outdated"},
		{name: "outdated-unknown-option", cmd: "outdated --latest"},
	}
	var out []Scenario
	for _, e := range list {
		js := e.json
		if js == "" {
			js = base
		}
		s := setup
		if e.setup != nil {
			s = e.setup
		}
		if e.nojson {
			s = nil
		}
		sc := Scenario{Name: e.name, JSON: js, Files: files, Setup: s, Command: e.cmd, NoJSON: e.nojson}
		if e.after != "" {
			sc.After = map[string]string{"composer.json": e.after}
		}
		out = append(out, sc)
	}
	return out
}
