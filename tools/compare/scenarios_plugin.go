package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Composer plugin scenarios (PluginManager, PluginInstaller, plugin event
// subscribers, CommandProvider capabilities and custom installers of Composer
// f83588392df43d4fd027d84d26de459946ef4e80) using the loopback fixture
// plugins testdata/compare/packages/plugfix/*.

func init() { extraScenarios = append(extraScenarios, pluginScenarios) }

func pluginScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	allow := func(s, v string) string { return withConfig(s, `"allow-plugins": `+v) }
	ev := `{"plugfix/events-plugin": "1.0.0", "acme/log": "^1.0"}`
	evDev := `"require-dev": {"acme/util": "^1.0"}`
	scripts := `"scripts": {"post-update-cmd": ["@putenv PLUGFIX_X=1", "echo root-post-update"], "post-install-cmd": "echo root-post-install", "plugfix-custom": "echo custom-script"}`
	evProject := allow(project(ev, evDev+",\n    "+scripts), `{"plugfix/*": true}`)
	evInstalled := []string{"update --no-progress"}
	mods := `{"plugfix/module-a": "1.0.0", "plugfix/module-b": "^1.0", "acme/log": "^1.0"}`
	// the zip extractions of Composer run in parallel and the custom
	// installer reports their completion: keep them sequential
	seq := map[string]string{"COMPOSER_MAX_PARALLEL_PROCESSES": "1"}
	modProject := allow(project(mods, ""), `{"plugfix/installer-plugin": true}`)
	single := func(name string) string {
		return allow(project(`{"plugfix/`+name+`": "^1.0"}`, ""), `{"plugfix/*": true}`)
	}
	notRequired := allow(project(`{"acme/log": "^1.0"}`, evDev), `{}`)
	lockOldJSON := project(`{"plugfix/events-plugin": "1.0.0"}`, "")
	alProject := allow(project(`{"plugfix/autoload-plugin": "^1.0"}`, `"scripts": {"plugfix-custom": ["echo custom-from-plugin $PLUGFIX_FROM_PLUGIN", "@putenv PLUGFIX_X=from-script"]}`), `{"plugfix/autoload-plugin": true}`)
	ghome := map[string]string{"COMPOSER_HOME": "/work/project/ghome"}
	globalHome := `{
    "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}],
    "config": {"secure-http": false, "allow-plugins": {"plugfix/*": true}},
    "require": {"plugfix/events-plugin": "1.0.0"}
}
`
	sc := []Scenario{
		// allow-plugins handling
		{Name: "plugin-blocked-update", JSON: project(ev, ""), Files: files, Command: "update --no-progress"},
		{Name: "plugin-blocked-install-lock", JSON: project(ev, ""), Files: files, Setup: []string{"update --no-progress --no-install --no-plugins"}, Command: "install --no-progress"},
		{Name: "plugin-blocked-installed", JSON: evProject, Files: files, Setup: evInstalled, After: map[string]string{"composer.json": project(ev, evDev+",\n    "+scripts)}, Command: "install --no-progress"},
		{Name: "plugin-allow-false", JSON: allow(project(ev, ""), `{"plugfix/events-plugin": false}`), Files: files, Command: "update --no-progress", Probe: true},
		{Name: "plugin-allow-all-false", JSON: allow(project(ev, ""), `false`), Files: files, Command: "update --no-progress"},
		{Name: "plugin-allow-true", JSON: allow(project(ev, ""), `true`), Files: files, Command: "update --no-progress"},
		{Name: "plugin-optional-blocked", JSON: project(`{"plugfix/optional": "^1.0"}`, ""), Files: files, Command: "update --no-progress"},
		{Name: "plugin-not-required-anymore", JSON: evProject, Files: files, Setup: evInstalled, After: map[string]string{"composer.json": notRequired}, Command: "install --no-progress"},
		{Name: "plugin-old-lock-no-allow", JSON: lockOldJSON, Files: map[string]string{"src/Kernel.php": appSrc, "composer.lock": oldPluginLock("2.1.0")}, Command: "install --no-progress"},
		{Name: "plugin-old-lock-new-api", JSON: lockOldJSON, Files: map[string]string{"src/Kernel.php": appSrc, "composer.lock": oldPluginLock("2.2.0")}, Command: "install --no-progress"},
		// event subscribers
		{Name: "plugin-allowed-update", JSON: evProject, Files: files, Command: "update --no-progress", Probe: true},
		{Name: "plugin-allowed-update-verbose", JSON: evProject, Files: files, Command: "update --no-progress -v"},
		{Name: "plugin-update-dry-run", JSON: evProject, Files: files, Command: "update --no-progress --dry-run"},
		{Name: "plugin-install-from-lock", JSON: evProject, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true},
		{Name: "plugin-install-no-dev", JSON: evProject, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress --no-dev"},
		{Name: "plugin-installed-noop", JSON: evProject, Files: files, Setup: evInstalled, Command: "install --no-progress"},
		{Name: "plugin-installed-quiet", JSON: evProject, Files: files, Setup: evInstalled, Command: "install --no-progress -q"},
		{Name: "plugin-no-scripts", JSON: evProject, Files: files, Setup: evInstalled, Command: "install --no-progress --no-scripts"},
		{Name: "plugin-dump-autoload", JSON: evProject, Files: files, Setup: evInstalled, Command: "dump-autoload", Probe: true},
		{Name: "plugin-dump-autoload-optimize", JSON: evProject, Files: files, Setup: evInstalled, Command: "dump-autoload -o --no-dev", Probe: true},
		{Name: "plugin-no-plugins-installed", JSON: evProject, Files: files, Setup: evInstalled, Command: "install --no-progress --no-plugins"},
		{Name: "plugin-no-plugins-update", JSON: evProject, Files: files, Command: "update --no-progress --no-plugins", Probe: true},
		{Name: "plugin-update-version", JSON: evProject, Files: files, Setup: evInstalled, After: map[string]string{"composer.json": strings.Replace(evProject, `"plugfix/events-plugin": "1.0.0"`, `"plugfix/events-plugin": "1.1.0"`, 1)}, Command: "update plugfix/events-plugin --no-progress"},
		{Name: "plugin-remove", JSON: evProject, Files: files, Setup: evInstalled, Command: "remove plugfix/events-plugin --no-progress"},
		{Name: "plugin-require", JSON: allow(project(`{"acme/log": "^1.0"}`, ""), `{"plugfix/*": true}`), Files: files, Command: "require plugfix/events-plugin:1.0.0 --no-progress"},
		{Name: "plugin-require-installed", JSON: evProject, Files: files, Setup: evInstalled, Command: "require acme/util:^1.0 --no-progress"},
		{Name: "plugin-run-script-custom", JSON: evProject, Files: files, Setup: evInstalled, Command: "run-script plugfix-custom -- a b"},
		{Name: "plugin-show", JSON: evProject, Files: files, Setup: evInstalled, Command: "show"},
		{Name: "plugin-thrower-boom", JSON: single("thrower"), Files: files, Setup: []string{"update --no-progress"}, After: map[string]string{"plugfix-boom": "1"}, Command: "dump-autoload"},
		{Name: "plugin-thrower-stop", JSON: single("thrower"), Files: files, Setup: []string{"update --no-progress"}, After: map[string]string{"plugfix-boom": "1", "plugfix-stop": "1"}, Command: "dump-autoload"},
		{Name: "plugin-thrower-false", JSON: single("thrower"), Files: files, Setup: []string{"update --no-progress"}, Command: "install --no-progress"},
		// plugins installed in COMPOSER_HOME (global plugins)
		{Name: "plugin-global", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: ghome, Setup: []string{"global update --no-progress"}, Command: "update --no-progress"},
		{Name: "plugin-global-blocked", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: ghome, Setup: []string{"global update --no-progress"}, After: map[string]string{"ghome/composer.json": strings.Replace(globalHome, `"plugfix/*": true`, `"other/*": true`, 1)}, Command: "update --no-progress"},
		{Name: "plugin-global-no-plugins", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: ghome, Setup: []string{"global update --no-progress"}, Command: "update --no-progress --no-plugins"},
		// root package changes and events dispatched by plugins
		{Name: "plugin-autoload-update", JSON: alProject, Files: files, Command: "update --no-progress", Probe: true, Snapshot: []string{"lib/generated/Info.php"}},
		{Name: "plugin-autoload-dump", JSON: alProject, Files: files, Setup: []string{"update --no-progress"}, Command: "dump-autoload --no-dev", Probe: true, Snapshot: []string{"lib/generated/Info.php"}},
		{Name: "plugin-autoload-run-script", JSON: alProject, Files: files, Setup: []string{"update --no-progress"}, Command: "run-script plugfix-custom -- a"},
		// plugin-api-version and package checks
		{Name: "plugin-old-api", JSON: single("old-api"), Files: files, Command: "update --no-progress"},
		{Name: "plugin-no-api", JSON: single("no-api"), Files: files, Command: "update --no-progress"},
		{Name: "plugin-no-class", JSON: single("no-class"), Files: files, Command: "update --no-progress"},
		{Name: "plugin-legacy-installer", JSON: single("legacy"), Files: files, Command: "update --no-progress"},
		// CommandProvider capability
		{Name: "plugin-command", JSON: evProject, Files: files, Setup: evInstalled, Command: "plugfix:hello Bob --shout -t 2"},
		{Name: "plugin-command-default", JSON: evProject, Files: files, Setup: evInstalled, Command: "plugfix:hello -v"},
		{Name: "plugin-command-exit", JSON: evProject, Files: files, Setup: evInstalled, Command: "plugfix:hello --times=3"},
		{Name: "plugin-command-bad-option", JSON: evProject, Files: files, Setup: evInstalled, Command: "plugfix:hello --nope"},
		{Name: "plugin-command-no-plugins", JSON: evProject, Files: files, Setup: evInstalled, Command: "plugfix:hello --no-plugins"},
		{Name: "plugin-command-conflict", JSON: single("commands-plugin"), Files: files, Setup: []string{"update --no-progress"}, Command: "plugfix:echo a b --tag=x --tag y -v"},
		{Name: "plugin-command-alias", JSON: single("commands-plugin"), Files: files, Setup: []string{"update --no-progress"}, Command: "pecho"},
		{Name: "plugin-command-throws", JSON: single("commands-plugin"), Files: files, Setup: []string{"update --no-progress"}, Command: "plugfix:echo --fail"},
		{Name: "plugin-command-unknown-ns", JSON: single("commands-plugin"), Files: files, Setup: []string{"update --no-progress"}, Command: "plugfix:nope"},
		{Name: "plugin-require-blocked", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: files, Command: "require plugfix/events-plugin:1.0.0 --no-progress"},
		// custom installers
		{Name: "plugin-installer-update", JSON: modProject, Files: files, Command: "update --no-progress", Probe: true, Env: seq},
		{Name: "plugin-installer-extra-path", JSON: allow(project(mods, `"extra": {"plugfix-module-dir": "app/mods"}`), `{"plugfix/installer-plugin": true}`), Files: files, Command: "update --no-progress", Probe: true, Env: seq},
		{Name: "plugin-installer-install", JSON: modProject, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true, Env: seq},
		{Name: "plugin-installer-upgrade", JSON: modProject, Files: files, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": strings.Replace(modProject, `"plugfix/module-a": "1.0.0"`, `"plugfix/module-a": "1.1.0"`, 1)}, Command: "update plugfix/module-a --no-progress", Probe: true},
		{Name: "plugin-installer-remove", JSON: modProject, Files: files, Setup: []string{"update --no-progress"}, Command: "remove plugfix/module-a --no-progress", Probe: true},
		{Name: "plugin-installer-no-plugins", JSON: modProject, Files: files, Command: "update --no-progress --no-plugins", Probe: true},
		{Name: "plugin-installer-dump", JSON: modProject, Files: files, Setup: []string{"update --no-progress"}, Command: "dump-autoload -o", Probe: true},
	}
	return sc
}

// oldPluginLock is a lock file for plugfix/events-plugin 1.0.0 written by a
// Composer with the given plugin API version (pre-2.2 locks enable the
// allow-plugins BC mode). References follow fixtures.go (sha1 of
// name@version); the content-hash is deliberately stale for both tools.
func oldPluginLock(api string) string {
	sum := sha1.Sum([]byte("plugfix/events-plugin@1.0.0"))
	ref := hex.EncodeToString(sum[:])
	return `{
    "_readme": ["fixture"],
    "content-hash": "00000000000000000000000000000000",
    "packages": [
        {
            "name": "plugfix/events-plugin",
            "version": "1.0.0",
            "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip", "reference": "` + ref + `"},
            "require": {"composer-plugin-api": "^2.0"},
            "type": "composer-plugin",
            "extra": {"class": "Plugfix\\Events\\Plugin"},
            "autoload": {"psr-4": {"Plugfix\\Events\\": "src/"}},
            "license": ["MIT"],
            "description": "Fixture plugin subscribing to events"
        }
    ],
    "packages-dev": [],
    "aliases": [],
    "minimum-stability": "stable",
    "stability-flags": {},
    "prefer-stable": false,
    "prefer-lowest": false,
    "platform": {},
    "platform-dev": {},
    "plugin-api-version": "` + api + `"
}
`
}
