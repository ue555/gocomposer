package main

import (
	"strings"
)

// Third set of Composer plugin scenarios (plugin3-*): the order of the
// post-file-download events of parallel downloads (loopback stub with a
// fixed delay per file, testdata/compare/vcs/plug3stub.sh), and the Composer
// services plugins replace or extend (fixture plugins plugfix3/*).
// Behaviour reference: Composer f83588392df43d4fd027d84d26de459946ef4e80.

func init() {
	extraScenarios = append(extraScenarios, plugin3Scenarios)
}

// plug3Delayed lists packages whose dists are served by the delay stub:
// a-slow completes last, c-fast first.
const plug3Delayed = `{"type": "package", "package": [
            {"name": "plugfix3/a-slow", "version": "1.0.0", "dist": {"type": "zip", "url": "http://127.0.0.1:8901/delay/1800/dists/acme-log-1.0.0.zip"}},
            {"name": "plugfix3/b-mid", "version": "1.0.0", "dist": {"type": "zip", "url": "http://127.0.0.1:8901/delay/900/dists/acme-util-1.0.0.zip"}},
            {"name": "plugfix3/c-fast", "version": "1.0.0", "dist": {"type": "zip", "url": "http://127.0.0.1:8901/delay/50/dists/acme-log-1.0.0.zip"}}
        ]},
        `

func plug3Project(require string, delayed bool) string {
	s := project(require, "")
	if delayed {
		s = strings.Replace(s, `{"type": "composer", "url": "http://127.0.0.1:8899"},`, plug3Delayed+`{"type": "composer", "url": "http://127.0.0.1:8899"},`, 1)
	}
	return withConfig(s, `"allow-plugins": {"plugfix3/*": true}`)
}

func plugin3Scenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	stub := []string{"plug3stub"}
	// download order (plugfix3/dl-order)
	orderBase := plug3Project(`{"plugfix3/dl-order": "^1.0"}`, true)
	orderFull := plug3Project(`{"plugfix3/dl-order": "^1.0", "plugfix3/a-slow": "1.0.0", "plugfix3/b-mid": "1.0.0", "plugfix3/c-fast": "1.0.0"}`, true)
	order := func(name string, env map[string]string) Scenario {
		e := map[string]string{}
		for k, v := range env {
			e[k] = v
		}
		return Scenario{Name: name, JSON: orderBase, Files: files, GitFixtures: stub, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": orderFull}, Command: "update --no-progress", Env: e}
	}
	// replaced downloaders and their internals (plugfix3/svc)
	svcBase := plug3Project(`{"plugfix3/svc": "^1.0"}`, false)
	svcLog := plug3Project(`{"plugfix3/svc": "^1.0", "acme/log": "1.0.0", "acme/util": "^1.0"}`, false)
	svcLog11 := plug3Project(`{"plugfix3/svc": "^1.0", "acme/log": "1.1.0", "acme/util": "^1.0"}`, false)
	seq := func(modes string, extra map[string]string) map[string]string {
		e := map[string]string{"PLUGFIX3_SVC": modes, "COMPOSER_MAX_PARALLEL_HTTP": "1", "COMPOSER_MAX_PARALLEL_PROCESSES": "1"}
		for k, v := range extra {
			e[k] = v
		}
		return e
	}
	svc := func(name, json, cmd string, setup []string, after map[string]string, env map[string]string) Scenario {
		return Scenario{Name: name, JSON: json, Files: files, Setup: setup, After: after, Command: cmd, Env: env, Probe: true}
	}
	installed := []string{"update --no-progress"}
	// sub/composer.json (Factory::create, events) and schema files
	svcExtra := map[string]string{
		"src/Kernel.php":       appSrc,
		"sub/composer.json":    `{"name": "test/sub", "require": {"acme/util": "^1.0"}, "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}], "config": {"secure-http": false}, "scripts": {"post-install-cmd": "echo sub-post-install", "plugfix3-custom": "echo sub-custom-event"}}` + "\n",
		"plug3-schema.json":    `{"type": "object", "required": ["plug3"], "properties": {"plug3": {"type": "string"}, "name": {"type": "integer"}}}` + "\n",
		"plug3-schema-ok.json": `{"type": "object", "properties": {"name": {"type": "string"}}}` + "\n",
	}
	svcFiles := func(name, json, cmd string, setup []string, env map[string]string) Scenario {
		return Scenario{Name: name, JSON: json, Files: svcExtra, Setup: setup, Command: cmd, Env: env}
	}
	sc := []Scenario{
		svc("plugin3-dl-replace-install", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("zip", nil)),
		svc("plugin3-dl-replace-update", svcLog, "update --no-progress", installed, map[string]string{"composer.json": svcLog11}, seq("zip", nil)),
		svc("plugin3-dl-replace-remove", svcLog, "remove acme/log --no-progress", installed, nil, seq("zip", nil)),
		svc("plugin3-dl-replace-reinstall", svcLog, "reinstall acme/log", installed, nil, seq("zip", nil)),
		svc("plugin3-dl-replace-fail", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("zip", map[string]string{"PLUGFIX3_FAIL": "acme/util"})),
		svc("plugin3-dl-extract", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("extract", nil)),
		svc("plugin3-dl-extract-update", svcLog, "update --no-progress", installed, map[string]string{"composer.json": svcLog11}, seq("extract", nil)),
		svc("plugin3-dl-extract-fail", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("extract", map[string]string{"PLUGFIX3_FAIL": "acme/util"})),
		svc("plugin3-dl-wrap", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("wrap", nil)),
		svc("plugin3-dl-wrap-update", svcLog, "update --no-progress", installed, map[string]string{"composer.json": svcLog11}, seq("wrap", nil)),
		svc("plugin3-dl-zip-parallel", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, map[string]string{"PLUGFIX3_SVC": "zip"}),
		svc("plugin3-dl-internals", svcLog, "plugfix3:svc dl-internals", installed, nil, nil),
		svc("plugin3-dl-no-plugins", svcBase, "update --no-progress --no-plugins", installed, map[string]string{"composer.json": svcLog}, seq("zip", nil)),
		// autoload generator replaced, its helpers, other target directory, custom platform filter
		svc("plugin3-al-replace-update", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("autoload", nil)),
		svc("plugin3-al-replace-dump", svcLog, "dump-autoload -o", installed, nil, seq("autoload", nil)),
		svc("plugin3-al-replace-dump-dev", svcLog, "dump-autoload --no-dev --apcu --dry-run", installed, nil, seq("autoload", nil)),
		svc("plugin3-al-replace-reinstall", svcLog, "reinstall acme/log", installed, nil, seq("autoload", nil)),
		svc("plugin3-al-replace-skip", svcBase, "update --no-progress", installed, map[string]string{"composer.json": svcLog}, seq("autoload", map[string]string{"PLUGFIX3_AL_SKIP": "1"})),
		svc("plugin3-al-replace-no-plugins", svcLog, "dump-autoload --no-plugins", installed, nil, seq("autoload", nil)),
		svc("plugin3-al-helpers", svcLog, "plugfix3:svc autoload-helpers", installed, nil, nil),
		svc("plugin3-al-target", svcLog, "plugfix3:svc autoload-target", installed, nil, nil),
		svc("plugin3-al-filter", plug3Project(`{"plugfix3/svc": "^1.0", "ext-json": "*", "ext-ctype": "*"}`, false), "plugfix3:svc autoload-filter", installed, nil, nil),
		// repository classes of plugins and changed repositories
		svc("plugin3-repo-class", svcBase, "update --no-progress", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "plugfix3/mem-lib": "^1.0"}`, false)}, seq("repo", nil)),
		svc("plugin3-repo-class-show", plug3Project(`{"plugfix3/svc": "^1.0"}`, false), "show -a plugfix3/mem-lib", installed, nil, seq("repo", nil)),
		svc("plugin3-repo-class-filter", svcBase, "update --no-progress", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "plugfix3/mem-other": "*"}`, false)}, seq("repo-filter", nil)),
		svc("plugin3-repo-overlay", svcBase, "update --no-progress", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "plugfix3/added": "^1.0"}`, false)}, seq("overlay", nil)),
		svc("plugin3-repo-overlay-removed", svcBase, "update --no-progress", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "plugfix3/listed": "^1.0"}`, false)}, seq("overlay", nil)),
		// Installer: additional fixed repository, lock transaction, audit config, platform filter of a plugin
		svc("plugin3-installer-fixed", plug3Project(`{"plugfix3/svc": "^1.0", "acme/util": "1.0.0"}`, false), "plugfix3:svc installer-fixed", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "acme/util": "^1.0", "acme/log": "^1.0"}`, false)}, nil),
		svc("plugin3-installer-unfixed", plug3Project(`{"plugfix3/svc": "^1.0", "acme/util": "1.0.0"}`, false), "plugfix3:svc installer-fixed", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "acme/util": "^1.0", "acme/log": "^1.0"}`, false)}, map[string]string{"PLUGFIX3_NOFIXED": "1"}),
		svc("plugin3-installer-audit", svcLog, "plugfix3:svc installer-audit", installed, nil, nil),
		svc("plugin3-installer-filter", svcBase, "plugfix3:svc installer-filter", installed, map[string]string{"composer.json": plug3Project(`{"plugfix3/svc": "^1.0", "acme/log": "^1.0", "ext-plugfix3": "*"}`, false)}, nil),
		svc("plugin3-lock-tx-none", svcLog, "plugfix3:svc lock-tx-none", installed, nil, nil),
		// Factory (configuration array, global, createComposer), Config sources and configs of plugins, nested events, schema files
		svcFiles("plugin3-factory", svcLog, "plugfix3:svc factory", installed, nil),
		svcFiles("plugin3-config", svcLog, "plugfix3:svc config", installed, nil),
		svcFiles("plugin3-events", svcLog, "plugfix3:svc events", installed, nil),
		svcFiles("plugin3-schema", svcLog, "plugfix3:svc schema", installed, nil),
		// ArchiveManager: the API, archivers of plugins in the archive command
		svcFiles("plugin3-archive-api", svcLog, "plugfix3:svc archive", installed, nil),
		svcFiles("plugin3-archive-cmd", svcLog, "archive --format=txt --dir=dist", installed, seq("archiver", nil)),
		svcFiles("plugin3-archive-cmd-package", svcLog, "archive acme/log 1.0.0 --format=txt --dir=dist", installed, seq("archiver", nil)),
		svcFiles("plugin3-archive-cmd-replaced", svcLog, "archive --format=txt --dir=dist", installed, seq("archive-manager", nil)),
		svcFiles("plugin3-archive-cmd-replaced-tar", svcLog, "archive --format=tar --dir=dist", installed, seq("archive-manager", nil)),
		svcFiles("plugin3-set-services", svcLog, "plugfix3:svc set-services", installed, nil),
		svcFiles("plugin3-archive-cmd-unknown", svcLog, "archive --format=nope --dir=dist", installed, seq("archiver", nil)),
		// TLS options of plugin HTTP requests (loopback HTTPS stub)
		{Name: "plugin3-http-tls", JSON: svcLog, Files: files, GitFixtures: stub, Setup: installed, Command: "plugfix3:svc http-tls"},
		order("plugin3-dl-order-parallel", nil),
		order("plugin3-dl-order-serial", map[string]string{"COMPOSER_MAX_PARALLEL_HTTP": "1"}),
		order("plugin3-dl-order-throw", map[string]string{"PLUGFIX3_ORDER_THROW": "plugfix3/b-mid"}),
	}
	return sc
}
