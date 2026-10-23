package main

import (
	"regexp"
	"strings"
)

// Second set of Composer plugin scenarios: the plugin events of downloads
// and pool creation (pre/post-file-download, pre-pool-create), the Composer
// services plugins use (HttpDownloader, RepositoryManager and repositories,
// Factory::create, Config::merge, DownloadManager, AutoloadGenerator,
// Installer), custom installers (plugin2-*), and symfony/flex 2.11.0 with a
// loopback HTTPS recipe endpoint (flex-*, testdata/compare/vcs/flexstub.sh).
// Behaviour reference: Composer f83588392df43d4fd027d84d26de459946ef4e80.

func init() {
	extraScenarios = append(extraScenarios, plugin2Scenarios)
	extraOutputNormalizers = append(extraOutputNormalizers, func(s string) string {
		return reFlexSession.ReplaceAllString(s, "${1}SESSION)")
	})
}

// extraOutputNormalizers are applied by normalizeOutput (main.go).
var extraOutputNormalizers []func(string) string

// symfony/flex prints a random session id (bin2hex(random_bytes(16))) with
// the recipes it applies.
var reFlexSession = regexp.MustCompile(`(Symfony operations: \d+ recipes? \()[0-9a-f]{32}\)`)

// flexJSON is a Symfony Flex project using the fixture recipe endpoint.
func flexJSON(require, extra string) string {
	return `{
    "name": "test/flex-app",
    "type": "project",
    "license": "proprietary",
    "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}],
    "config": {"secure-http": false, "cafile": "/fixtures/git/_api/ca.pem", "allow-plugins": {"symfony/flex": true}, "sort-packages": true},
    "require": ` + require + `,
    "scripts": {"auto-scripts": {"echo auto-ok": "script"}, "post-install-cmd": ["@auto-scripts"], "post-update-cmd": ["@auto-scripts"]},
    "extra": {"symfony": {"endpoint": ["https://127.0.0.1/flex/index.json", "https://127.0.0.1/flex/contrib.json"], "allow-contrib": false` + extra + `}}
}
`
}

var flexSnapshot = []string{"symfony.lock", ".env", ".gitignore", "config/bundles.php", "config/packages/hello.yaml", "config/services.yaml", "bin/recipe-lib"}

func plugin2Scenarios() []Scenario {
	flexFx := []string{"flexstub"}
	flexOnly := flexJSON(`{"symfony/flex": "^2.11"}`, "")
	withHello := flexJSON(`{"flexfix/hello-bundle": "^1.0", "symfony/flex": "^2.11"}`, "")
	flex := func(name, json, cmd string, setup []string, after map[string]string) Scenario {
		return Scenario{Name: name, JSON: json, GitFixtures: flexFx, Setup: setup, After: after, Command: cmd, Snapshot: flexSnapshot}
	}
	// the container configurator of the recipe-lib recipe edits config/services.yaml
	withServices := func(s Scenario) Scenario {
		s.Files = map[string]string{"config/services.yaml": "parameters:\n\nservices:\n    _defaults:\n        autowire: true\n"}
		return s
	}
	files := map[string]string{"src/Kernel.php": appSrc}
	allow := func(s, v string) string { return withConfig(s, `"allow-plugins": `+v) }
	// download and pool events (plugfix/dl-events)
	dlBase := allow(project(`{"plugfix/dl-events": "^1.0"}`, ""), `{"plugfix/*": true}`)
	dlFull := allow(project(`{"plugfix/dl-events": "^1.0", "acme/log": "^1.0", "acme/util": "^1.0"}`, ""), `{"plugfix/*": true}`)
	dlLocalCache := func(s string) string { return withConfig(s, `"cache-files-dir": "fcache"`) }
	dlEnv := func(modes string) map[string]string {
		return map[string]string{"PLUGFIX_DL": modes, "COMPOSER_MAX_PARALLEL_HTTP": "1"}
	}
	dl := func(name, modes, json, cmd string, setup []string, after map[string]string) Scenario {
		env := dlEnv(modes)
		if strings.HasSuffix(name, "-debug-events") {
			env["COMPOSER_DEBUG_EVENTS"] = "1"
		}
		return Scenario{Name: name, JSON: json, Files: files, Setup: setup, After: after, Command: cmd, Env: env}
	}
	dlInstalled := []string{"update --no-progress"}
	// Composer services used by plugins (plugfix/api-plugin)
	apiJSON := allow(project(`{"plugfix/api-plugin": "^1.0", "acme/log": "^1.0"}`, ""), `{"plugfix/*": true}`)
	apiFiles := map[string]string{"src/Kernel.php": appSrc, "sub/composer.json": `{"name": "test/sub", "require": {"acme/util": "^1.0"}, "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}], "config": {"secure-http": false}}` + "\n"}
	api := func(name, cmd string, env map[string]string) Scenario {
		return Scenario{Name: name, JSON: apiJSON, Files: apiFiles, Setup: []string{"update --no-progress"}, Command: cmd, Env: env}
	}
	// installers replacing hooks or not built on LibraryInstaller (plugfix/custom-installers)
	ciJSON := allow(project(`{"plugfix/custom-installers": "^1.0", "plugfix/lib-a": "1.0.0", "plugfix/asset-a": "1.0.0", "plugfix/meta-a": "^1.0"}`, ""), `{"plugfix/*": true}`)
	ciOnly := allow(project(`{"plugfix/custom-installers": "^1.0"}`, ""), `{"plugfix/*": true}`)
	// lib-a is the first download (no other transfer pending when its download hook throws)
	ciLib := allow(project(`{"plugfix/custom-installers": "^1.0", "plugfix/lib-a": "1.0.0"}`, ""), `{"plugfix/*": true}`)
	seq := map[string]string{"COMPOSER_MAX_PARALLEL_PROCESSES": "1", "COMPOSER_MAX_PARALLEL_HTTP": "1"}
	ci := func(name, json, cmd string, setup []string, after map[string]string, env map[string]string) Scenario {
		if env == nil {
			env = seq
		}
		return Scenario{Name: name, JSON: json, Files: files, Setup: setup, After: after, Command: cmd, Env: env, Probe: true}
	}
	sc := []Scenario{
		dl("plugin2-dl-update", "", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		dl("plugin2-dl-update-v2", "", withV2(dlBase), "update --no-progress", dlInstalled, map[string]string{"composer.json": withV2(dlFull)}),
		dl("plugin2-dl-rewrite", "rewrite", dlLocalCache(dlBase), "update --no-progress", dlInstalled, map[string]string{"composer.json": dlLocalCache(dlFull)}),
		dl("plugin2-dl-meta", "meta", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		dl("plugin2-dl-pool", "pool,unacceptable", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		dl("plugin2-dl-pool-partial", "pool", dlFull, "update acme/log --no-progress", []string{"update --no-progress"}, nil),
		dl("plugin2-dl-install-lock", "", dlFull, "install --no-progress", []string{"update --no-progress --no-install"}, nil),
		dl("plugin2-dl-throw-pre", "throw-pre", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		dl("plugin2-dl-throw-pool", "throw-pool", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		dl("plugin2-dl-cache-hit", "rewrite", dlLocalCache(dlFull), "update --no-progress -vv", []string{"update --no-progress", "remove acme/log acme/util --no-progress"}, map[string]string{"composer.json": dlLocalCache(dlFull)}),
		dl("plugin2-dl-require", "", dlBase, "require acme/log:^1.0 --no-progress", dlInstalled, nil),
		dl("plugin2-dl-no-plugins", "", dlBase, "update --no-progress --no-plugins", dlInstalled, map[string]string{"composer.json": dlFull}),
		// services
		api("plugin2-api-http", "plugfix:api http", map[string]string{"COMPOSER_MAX_PARALLEL_HTTP": "1"}),
		api("plugin2-api-repos", "plugfix:api repos", nil),
		api("plugin2-api-create-repo", "plugfix:api create-repo", nil),
		api("plugin2-api-factory", "plugfix:api factory", nil),
		api("plugin2-api-merge", "plugfix:api merge", nil),
		api("plugin2-api-errors", "plugfix:api errors", nil),
		api("plugin2-api-autoload", "plugfix:api autoload", nil),
		api("plugin2-api-dm", "plugfix:api dm", map[string]string{"COMPOSER_MAX_PARALLEL_HTTP": "1"}),
		api("plugin2-api-install", "plugfix:api install", nil),
		api("plugin2-api-update", "plugfix:api update", nil),
		{Name: "plugin2-api-addrepo", JSON: apiJSON, Files: apiFiles, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": strings.Replace(apiJSON, `"acme/log": "^1.0"`, `"plugfix/virtual": "^1.0"`, 1)}, Command: "update --no-progress", Env: map[string]string{"PLUGFIX_API_ADDREPO": "1"}, Probe: true},
		{Name: "plugin2-api-merge-activate", JSON: apiJSON, Files: apiFiles, Setup: []string{"update --no-progress"}, Command: "require acme/util:^1.0 --no-progress", Env: map[string]string{"PLUGFIX_API_MERGE": "1"}},
		// a downloader registered by a plugin for a custom dist type
		{Name: "plugin2-api-downloader", JSON: strings.Replace(apiJSON, `{"type": "composer", "url": "http://127.0.0.1:8899"},`, `{"type": "package", "package": {"name": "plugfix/custom-dist", "version": "1.0.0", "dist": {"type": "plugfixzip", "url": "http://127.0.0.1:8899/dists/acme-util-1.0.0.zip"}}}, {"type": "composer", "url": "http://127.0.0.1:8899"},`, 1), Files: apiFiles, Setup: []string{"update --no-progress"}, Command: "require plugfix/custom-dist:1.0.0 --no-progress", Env: map[string]string{"PLUGFIX_API_DOWNLOADER": "1"}, Probe: true},
		// root scripts listening to the download and pool events, without plugins
		{Name: "plugin2-script-listeners", JSON: project(`{"acme/log": "^1.0"}`, `"scripts": {"pre-file-download": "echo script-pre-file-download", "post-file-download": "echo script-post-file-download", "pre-pool-create": "echo script-pre-pool-create"}`), Files: files, Command: "update --no-progress", Env: map[string]string{"COMPOSER_MAX_PARALLEL_HTTP": "1"}},
		dl("plugin2-dl-debug-events", "", dlBase, "update --no-progress", dlInstalled, map[string]string{"composer.json": dlFull}),
		// custom installers
		ci("plugin2-installers-fresh", ciJSON, "update --no-progress", nil, nil, nil),
		ci("plugin2-installers-installed", ciOnly, "update --no-progress", []string{"update --no-progress"}, map[string]string{"composer.json": ciJSON}, nil),
		ci("plugin2-installers-upgrade", ciJSON, "update --no-progress", []string{"update --no-progress"}, map[string]string{"composer.json": strings.Replace(strings.Replace(ciJSON, `"plugfix/lib-a": "1.0.0"`, `"plugfix/lib-a": "1.1.0"`, 1), `"plugfix/asset-a": "1.0.0"`, `"plugfix/asset-a": "1.1.0"`, 1)}, nil),
		ci("plugin2-installers-remove", ciJSON, "remove plugfix/lib-a plugfix/asset-a plugfix/meta-a --no-progress", []string{"update --no-progress"}, nil, nil),
		ci("plugin2-installers-install-lock", ciJSON, "install --no-progress", []string{"update --no-progress --no-install"}, nil, nil),
		ci("plugin2-installers-fail-download", ciOnly, "update --no-progress", []string{"update --no-progress"}, map[string]string{"composer.json": ciLib}, map[string]string{"PLUGFIX_HOOK_FAIL": "download", "COMPOSER_MAX_PARALLEL_PROCESSES": "1"}),
		ci("plugin2-installers-fail-install", ciOnly, "update --no-progress", []string{"update --no-progress"}, map[string]string{"composer.json": ciJSON}, map[string]string{"PLUGFIX_HOOK_FAIL": "install", "COMPOSER_MAX_PARALLEL_PROCESSES": "1"}),
		// flex installed by the run itself: reinstall with Factory::create and the cloned Installer
		flex("flex-update-fresh", withHello, "update --no-progress", nil, nil),
		// flex active from the start: recipes of the new packages are applied
		flex("flex-update-installed", flexOnly, "update --no-progress", []string{"update --no-progress"}, map[string]string{"composer.json": withHello}),
		flex("flex-require-alias", flexOnly, "require hello --no-progress", []string{"update --no-progress"}, nil),
		flex("flex-require-contrib", flexOnly, "require flexfix/contrib-lib --no-progress", []string{"update --no-progress"}, nil),
		withServices(flex("flex-require-recipe-lib", flexOnly, "require flexfix/recipe-lib --no-progress", []string{"update --no-progress"}, nil)),
		flex("flex-require-auto-bundle", flexOnly, "require flexfix/auto-bundle --no-progress", []string{"update --no-progress"}, nil),
		withServices(flex("flex-require-pack", flexOnly, "require fpack --no-progress", []string{"update --no-progress"}, nil)),
		flex("flex-require-split", flexJSON(`{"symfony/flex": "^2.11"}`, `, "require": "1.*"`), "require flexfix/split-a --no-progress", []string{"update --no-progress"}, nil),
		flex("flex-update-split", flexJSON(`{"flexfix/split-a": "*", "symfony/flex": "^2.11"}`, `, "require": "1.*"`), "update --no-progress", []string{"update --no-progress"}, nil),
		flex("flex-remove", withHello, "remove flexfix/hello-bundle --no-progress", []string{"update --no-progress", "update --no-progress"}, nil),
		flex("flex-install-lock", withHello, "install --no-progress", []string{"update --no-progress --no-install"}, nil),
		flex("flex-recipes", withHello, "recipes", []string{"update --no-progress", "update --no-progress"}, nil),
		flex("flex-recipes-one", withHello, "recipes flexfix/hello-bundle", []string{"update --no-progress", "update --no-progress"}, nil),
		flex("flex-recipes-install", withHello, "symfony:recipes:install --no-interaction", []string{"update --no-progress --no-scripts"}, nil),
		flex("flex-sync-recipes-alias", withHello, "sync-recipes -n", []string{"update --no-progress --no-scripts"}, nil),
		flex("flex-dump-env", withHello, "dump-env prod", []string{"update --no-progress"}, nil),
		flex("flex-update-noop", withHello, "update --no-progress", []string{"update --no-progress", "update --no-progress"}, nil),
		// (without flex, the auto-scripts values would run as commands: `script`)
		flex("flex-no-plugins", withHello, "update --no-progress --no-plugins --no-scripts", []string{"update --no-progress"}, nil),
		// create-project of a flex skeleton: flex-require merged (with
		// extra.symfony.require), reinstall through Factory::create and the
		// cloned Installer, configureProject (license, name, LICENSE removed)
		{Name: "flex-create-project", NoJSON: true, GitFixtures: flexFx, Command: "create-project flexfix/skeleton app --repository=http://127.0.0.1:8899 --no-secure-http --no-progress",
			Snapshot: []string{"app/composer.json", "app/composer.lock", "app/symfony.lock", "app/.env", "app/config/bundles.php", "app/vendor/composer/installed.json"}},
		// recipes:update of a recipe installed from an older (archived) ref,
		// in a git repository (RecipePatcher builds and applies a patch)
		{Name: "flex-recipes-update", JSON: strings.Replace(withHello, `"post-install-cmd"`, `"mk-old": ["sed -i 's/2a2b2c2d2e2f303132333435363738393a3b3c3d/0a0b0c0d0e0f101112131415161718191a1b1c1d/' symfony.lock", "printf 'hello:\\n    greeting: hi\\n' > config/packages/hello.yaml", "git init -q", "git add -A", "git commit -qm init"], "post-install-cmd"`, 1),
			GitFixtures: flexFx, Setup: []string{"update --no-progress", "run-script mk-old"}, Command: "recipes:update flexfix/hello-bundle", Env: gitIdent, Snapshot: flexSnapshot},
		{Name: "flex-recipes-update-no-changelog", JSON: strings.Replace(withHello, `"post-install-cmd"`, `"mk-old": ["sed -i 's/2a2b2c2d2e2f303132333435363738393a3b3c3d/0a0b0c0d0e0f101112131415161718191a1b1c1d/' symfony.lock", "printf 'hello:\\n    greeting: hi\\n' > config/packages/hello.yaml", "git init -q", "git add -A", "git commit -qm init"], "post-install-cmd"`, 1),
			GitFixtures: flexFx, Setup: []string{"update --no-progress", "run-script mk-old"}, Command: "recipes:update flexfix/hello-bundle --no-changelog", Env: gitIdent, Snapshot: flexSnapshot},
		// flex installed as a dependency only: recipes and package.json sync are disabled
		flex("flex-not-in-root", flexJSON(`{"flexfix/needs-flex": "^1.0"}`, ""), "require flexfix/hello-bundle --no-progress", []string{"update --no-progress"}, nil),
	}
	return sc
}
