package main

import "strings"

// Scenarios for run-script, exec, clear-cache, fund, search, browse,
// diagnose, self-update, repository and policy (Composer
// f83588392df43d4fd027d84d26de459946ef4e80 commands of the same names).

func init() { extraScenarios = append(extraScenarios, miscCommandScenarios) }

type miscSc struct {
	name, json, cmd string
	env             map[string]string
	files           map[string]string
	setup           []string
	nojson          bool
	known           string
	snap            []string
}

func miscCommandScenarios() []Scenario {
	var list []miscSc
	list = append(list, runScriptScenarios()...)
	list = append(list, execScenarios()...)
	list = append(list, clearCacheScenarios()...)
	list = append(list, fundScenarios()...)
	list = append(list, searchScenarios()...)
	list = append(list, browseScenarios()...)
	list = append(list, diagnoseScenarios()...)
	list = append(list, selfUpdateScenarios()...)
	list = append(list, repositoryScenarios()...)
	list = append(list, policyCmdScenarios()...)
	list = append(list, miscMoreScenarios...)
	var out []Scenario
	for _, s := range list {
		out = append(out, Scenario{
			Name:      s.name,
			JSON:      s.json,
			NoJSON:    s.nojson,
			Command:   s.cmd,
			Env:       s.env,
			Files:     s.files,
			Setup:     s.setup,
			KnownDiff: s.known,
			Snapshot:  s.snap,
		})
	}
	return out
}

// miscMoreScenarios is extended by later command groups in this file set.
var miscMoreScenarios []miscSc

func runScriptScenarios() []miscSc {
	scripts := func(s string) string { return project(`{}`, s) }
	base := scripts(`"scripts": {
        "test": "echo test-ran",
        "args": "echo args:",
        "devmode": "echo dev=$COMPOSER_DEV_MODE",
        "fail": "exit 3",
        "slow": "sleep 3; echo slow-done",
        "nested": ["@test", "echo nested-done"],
        "missing-ref": "@nope",
        "multi": ["echo one", "echo two @additional_args", "echo three"],
        "noargs": "echo noargs @no_additional_args",
        "post-install-cmd": "echo post-install",
        "pre-update-cmd": ["echo pre-update"],
        "post-package-install": "echo ppi",
        "env-set": ["@putenv FOO=bar", "echo foo=$FOO"],
        "composer-call": "@composer --version --no-ansi"
    },
    "scripts-descriptions": {"test": "Runs the tests", "nested": "Nested scripts"},
    "scripts-aliases": {"test": ["t", "tst"]}`)
	r := func(name, cmd string) miscSc { return miscSc{name: "run-script-" + name, json: base, cmd: cmd} }
	list := []miscSc{
		r("basic", "run-script test"),
		r("alias-run", "run test"),
		r("abbrev", "run-s test"),
		r("args", "run-script args a 'b c' --no-ansi"),
		r("args-dashdash", "run-script args -- --foo -x bar"),
		r("list", "run-script --list"),
		r("list-short", "run -l"),
		r("missing-argument", "run-script"),
		r("not-defined", "run-script nope"),
		r("event-constant-name", "run-script pre_install_cmd"),
		r("event", "run-script post-install-cmd"),
		r("event-array", "run-script pre-update-cmd -v"),
		r("package-event", "run-script post-package-install"),
		r("dev-mode", "run-script devmode"),
		r("no-dev", "run-script devmode --no-dev"),
		r("dev-and-no-dev", "run-script devmode --dev --no-dev"),
		{name: "run-script-no-dev-env", json: base, cmd: "run-script devmode", env: map[string]string{"COMPOSER_NO_DEV": "1"}},
		r("timeout-invalid", "run-script test --timeout=abc"),
		r("timeout-negative", "run-script test --timeout=-1"),
		r("timeout", "run-script slow --timeout=1"),
		r("timeout-zero", "run-script slow --timeout=0"),
		{name: "run-script-process-timeout-env", json: base, cmd: "run-script slow", env: map[string]string{"COMPOSER_PROCESS_TIMEOUT": "1"}},
		r("fail", "run-script fail"),
		r("nested", "run-script nested"),
		r("nested-verbose", "run-script nested -v"),
		r("missing-ref", "run-script missing-ref"),
		r("multi-args", "run-script multi x y"),
		r("no-additional-args", "run-script noargs x"),
		r("putenv", "run-script env-set"),
		r("no-scripts", "run-script test --no-scripts"),
		{name: "run-script-skip-env", json: base, cmd: "run-script test", env: map[string]string{"COMPOSER_SKIP_SCRIPTS": "test"}},
		{name: "run-script-skip-env-verbose", json: base, cmd: "run-script nested -v", env: map[string]string{"COMPOSER_SKIP_SCRIPTS": "test,args"}},
		r("quiet", "run-script test -q"),
		// scripts as commands (ScriptAliasCommand)
		r("direct", "test"),
		r("direct-args", "args a 'b c'"),
		r("direct-unknown-option", "args a --foo b"),
		r("direct-dashdash", "args -- --foo b"),
		r("direct-alias", "t"),
		r("direct-alias2", "tst"),
		r("direct-nested", "nested"),
		r("direct-fail", "fail"),
		r("direct-dev-mode", "devmode --no-dev"),
		r("direct-event-not-command", "post-install-cmd"),
		r("direct-abbrev", "nest"),
		r("unknown-command", "nope"),
		{name: "run-script-override-warning", json: scripts(`"scripts": {"install": "echo mine", "validate": "echo v"}`), cmd: "run-script install"},
		{name: "run-script-aliases-invalid", json: scripts(`"scripts": {"test": "echo t"}, "scripts-aliases": {"test": [1]}`), cmd: "run-script test"},
		{name: "run-script-list-empty", json: project(`{}`, ""), cmd: "run-script --list"},
		{name: "run-script-no-json", nojson: true, cmd: "run-script test"},
		{name: "run-script-disable-timeout", json: scripts(`"scripts": {"long": ["Composer\\Config::disableProcessTimeout", "sleep 2; echo long-done"]}`), cmd: "run-script long", env: map[string]string{"COMPOSER_PROCESS_TIMEOUT": "1"}},
		{name: "run-script-working-dir", json: project(`{}`, ""), files: map[string]string{"sub/composer.json": "{\"scripts\": {\"hello\": \"echo hello-from-sub\"}}\n"}, cmd: "hello -d sub"},
	}
	return list
}

func execScenarios() []miscSc {
	withCli := project(`{"acme/cli": "^1.0"}`, "")
	setup := []string{"update --no-progress"}
	local := project(`{}`, `"bin": ["bin/tool"]`)
	localFiles := map[string]string{"bin/tool": "#!/bin/sh\necho tool \"$@\"\n"}
	e := func(name, cmd string) miscSc {
		return miscSc{name: "exec-" + name, json: withCli, cmd: cmd, setup: setup}
	}
	return []miscSc{
		e("list", "exec --list"),
		e("no-binary", "exec"),
		e("run", "exec acme-sh"),
		e("run-args", "exec acme-sh -- a 'b c'"),
		e("run-php", "exec acme"),
		e("missing", "exec nope"),
		e("no-scripts", "exec acme-sh --no-scripts"),
		e("verbose", "exec acme-sh -v"),
		{name: "exec-no-bins", json: project(`{}`, ""), cmd: "exec --list"},
		{name: "exec-local-list", json: local, files: localFiles, cmd: "exec -l"},
		{name: "exec-no-json", nojson: true, cmd: "exec -l"},
	}
}

func clearCacheScenarios() []miscSc {
	cenv := map[string]string{"COMPOSER_CACHE_DIR": "/work/project/c"}
	big := strings.Repeat("0123456789abcdef", 128)
	files := map[string]string{
		"c/files/acme/log/abc.zip":      big,
		"c/repo/http---x/packages.json": "{}",
		"c/vcs/foo.git/HEAD":            "ref",
		"c/other.txt":                   "x",
		"c/.htaccess":                   "Deny from all",
	}
	base := project(`{}`, "")
	cc := func(name, cmd, json string) miscSc {
		return miscSc{name: "clear-cache-" + name, json: json, cmd: cmd, env: cenv, files: files}
	}
	return []miscSc{
		cc("basic", "clear-cache", base),
		cc("alias-cc", "cc", base),
		cc("alias-clearcache", "clearcache", base),
		cc("gc", "clear-cache --gc", base),
		cc("gc-maxsize", "clear-cache --gc", withConfig(base, `"cache-files-maxsize": "1KiB"`)),
		// a ttl of 0 is not compared: whether a file written just before the
		// run is older than "now" depends on timing (mtime < now, seconds);
		// expiry of old files is covered by TestCacheGCTTL instead
		cc("gc-ttl-fresh", "clear-cache --gc", withConfig(base, `"cache-ttl": 3600, "cache-files-ttl": 3600`)),
		cc("gc-bad-maxsize", "clear-cache --gc", withConfig(base, `"cache-files-maxsize": "lots"`)),
		cc("read-only", "clear-cache", withConfig(base, `"cache-read-only": true`)),
		cc("verbose", "clear-cache -v", base),
		{name: "clear-cache-empty", json: base, cmd: "clear-cache"},
		{name: "clear-cache-empty-gc", json: base, cmd: "clear-cache --gc"},
		{name: "clear-cache-no-cache", json: base, cmd: "clear-cache --no-cache"},
		{name: "clear-cache-no-json", nojson: true, cmd: "clear-cache", env: cenv, files: files},
		{name: "clear-cache-custom-dirs", json: withConfig(base, `"cache-files-dir": "c2/f", "cache-vcs-dir": "/dev/null"`), cmd: "clear-cache", files: map[string]string{"c2/f/a.zip": "x"}},
	}
}

func fundScenarios() []miscSc {
	deps := project(`{"fundco/alpha": "^1.0", "fundco/beta": "^1.0", "fundco/plain": "^1.0", "otherfund/gamma": "^1.0"}`, "")
	setup := []string{"update --no-progress"}
	f := func(name, cmd string) miscSc {
		return miscSc{name: "fund-" + name, json: deps, cmd: cmd, setup: setup}
	}
	return []miscSc{
		f("text", "fund"),
		f("json", "fund --format=json"),
		f("json-short", "fund -f json"),
		f("bad-format", "fund --format=xml"),
		{name: "fund-v2", json: withV2(deps), cmd: "fund", setup: setup},
		{name: "fund-none", json: project(`{"fundco/plain": "^1.0"}`, ""), cmd: "fund", setup: setup},
		{name: "fund-none-json", json: project(`{"fundco/plain": "^1.0"}`, ""), cmd: "fund --format json", setup: setup},
		{name: "fund-not-installed", json: deps, cmd: "fund"},
		{name: "fund-no-json", nojson: true, cmd: "fund"},
	}
}

func searchScenarios() []miscSc {
	base := project(`{}`, "")
	s := func(name, cmd string) miscSc { return miscSc{name: "search-" + name, json: base, cmd: cmd} }
	return []miscSc{
		s("fulltext", "search fundco"),
		s("multi-token", "search alpha gamma"),
		s("keyword", "search sponsor"),
		s("only-name", "search -N fundco"),
		s("only-name-long", "search --only-name fundco/al"),
		s("only-vendor", "search -O fund"),
		s("only-vendor-long", "search --only-vendor acme"),
		s("name-and-vendor", "search -N -O fund"),
		s("type", "search fundco --type library"),
		s("type-other", "search fund -t composer-plugin-like"),
		s("type-vendor", "search -O fund -t library"),
		s("json", "search fundco --format=json"),
		s("json-vendor", "search -O fund --format=json"),
		s("json-empty", "search zzzz --format=json"),
		s("bad-format", "search fundco --format=xml"),
		s("no-results", "search zzzz"),
		s("regex-chars", "search 'fund.o'"),
		s("regex-chars-name", "search -N 'fund.o'"),
		s("no-tokens", "search"),
		s("very-verbose", "search fundco -vv"),
		s("columns", "search fundco/beta"),
		{name: "search-installed", json: project(`{"fundco/alpha": "^1.0"}`, ""), cmd: "search alpha", setup: []string{"update --no-progress"}},
		{name: "search-installed-vv", json: project(`{"fundco/alpha": "^1.0"}`, ""), cmd: "search alpha -vv", setup: []string{"update --no-progress"}},
		{name: "search-v2", json: withV2(base), cmd: "search fundco"},
		{name: "search-v2-vendor", json: withV2(base), cmd: "search -O fund"},
		{name: "search-v2-json", json: withV2(base), cmd: "search alpha --format=json"},
		{name: "search-wide", json: base, cmd: "search fundco/beta", env: map[string]string{"COLUMNS": "200"}},
		{name: "search-filtered-repo", json: strings.Replace(base, `{"type": "composer", "url": "http://127.0.0.1:8899"}`, `{"type": "composer", "url": "http://127.0.0.1:8899", "only": ["fundco/*"]}`, 1), cmd: "search fund"},
	}
}

func browseScenarios() []miscSc {
	base := project(`{}`, "")
	b := func(name, cmd string) miscSc { return miscSc{name: "browse-" + name, json: base, cmd: cmd} }
	return []miscSc{
		b("show", "browse fundco/alpha -s"),
		b("show-homepage", "browse fundco/alpha -H -s"),
		b("home-alias", "home fundco/plain --show"),
		b("homepage-missing", "browse fundco/plain --homepage --show"),
		b("not-found", "browse nope/nope -s"),
		b("no-url", "browse otherfund/gamma -s"),
		b("multi", "browse fundco/alpha fundco/plain nope/x -s"),
		b("open", "browse fundco/alpha"),
		b("root", "browse -s"),
		{name: "browse-root-homepage", json: project(`{}`, `"homepage": "https://example.org/root", "support": {"source": "https://example.org/root/src"}`), cmd: "browse -s -H"},
		{name: "browse-root-source", json: project(`{}`, `"homepage": "https://example.org/root", "support": {"source": "https://example.org/root/src"}`), cmd: "browse --show"},
		{name: "browse-root-open", json: project(`{}`, `"homepage": "https://example.org/root"`), cmd: "browse -H"},
		{name: "browse-root-invalid-url", json: project(`{}`, `"homepage": "example.com"`), cmd: "browse -s -H"},
		{name: "browse-installed", json: project(`{"fundco/alpha": "^1.0"}`, ""), cmd: "browse fundco/alpha -s", setup: []string{"update --no-progress"}},
		{name: "browse-no-json", nojson: true, cmd: "browse"},
	}
}

func diagnoseScenarios() []miscSc {
	nonet := map[string]string{"COMPOSER_DISABLE_NETWORK": "1"}
	env := func(kv ...string) map[string]string {
		m := map[string]string{"COMPOSER_DISABLE_NETWORK": "1"}
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return m
	}
	base := project(`{}`, "")
	clean := strings.Replace(base, `"type": "project",`, `"type": "project", "license": "MIT",`, 1)
	d := func(name, json string, e map[string]string) miscSc {
		return miscSc{name: "diagnose-" + name, json: json, cmd: "diagnose", env: e}
	}
	return []miscSc{
		d("basic", base, nonet),
		d("clean", clean, nonet),
		{name: "diagnose-locked", json: clean, cmd: "diagnose", env: nonet, files: map[string]string{"composer.lock": emptyLock}},
		{name: "diagnose-lock-invalid", json: clean, cmd: "diagnose", env: nonet, files: map[string]string{"composer.lock": strings.Replace(emptyLock, `"packages": []`, `"packages": "x"`, 1)}},
		{name: "diagnose-no-json", nojson: true, cmd: "diagnose", env: nonet},
		d("schema-error", strings.Replace(base, `"name": "test/project",`, `"name": "Test Project",`, 1), nonet),
		d("platform-override", withConfig(clean, `"platform": {"php": "8.1.0"}`), nonet),
		d("disable-tls", withConfig(clean, `"disable-tls": true`), nonet),
		d("github-oauth", withConfig(clean, `"github-oauth": {"github.com": "abcdef", "git.example.org": "xyz"}`), nonet),
		d("proxy", clean, env("http_proxy", "http://user:pass@proxy.example:3128", "HTTPS_PROXY", "https://proxy.example")),
		d("proxy-no-proxy", clean, env("http_proxy", "http://proxy.example:3128", "no_proxy", ".packagist.org")),
		d("proxy-invalid", clean, env("http_proxy", "http://proxy.example:0")),
		d("ipresolve", clean, env("COMPOSER_IPRESOLVE", "4")),
		{name: "diagnose-network", json: clean, cmd: "diagnose", known: "Checking Composer and its dependencies for vulnerabilities: Composer audits its own installation (vendor/composer/installed.json of Composer itself) against packagist.org; gocomposer has no Composer installation and reports Composer's non-standard-installation warning instead. All other network checks (packagist http/https, fixture repository, GitHub rate limit) match"},
	}
}

const emptyLock = `{
    "_readme": ["test"],
    "content-hash": "0123456789abcdef0123456789abcdef",
    "packages": [],
    "packages-dev": [],
    "aliases": [],
    "minimum-stability": "stable",
    "stability-flags": {},
    "prefer-stable": false,
    "prefer-lowest": false,
    "platform": {},
    "platform-dev": {},
    "plugin-api-version": "2.6.0"
}
`

func selfUpdateScenarios() []miscSc {
	why := "gocomposer has no self-update mechanism (it is not a composer.phar): like Composer run from a source checkout it prints \"This instance of Composer does not have the self-update command.\" and exits 1, but the explanatory second line names gocomposer instead of Composer's system-package/project-package reasons"
	base := project(`{}`, "")
	return []miscSc{
		{name: "self-update-basic", json: base, cmd: "self-update", known: why},
		{name: "self-update-alias", json: base, cmd: "selfupdate --stable", known: why},
		{name: "self-update-rollback", json: base, cmd: "self-update -r", known: why},
		{name: "self-update-version", json: base, cmd: "self-update 2.8.0 --no-progress --2", known: why},
		{name: "self-update-no-json", nojson: true, cmd: "self-update --update-keys", known: why},
		{name: "self-update-bad-option", json: base, cmd: "self-update --nope"},
		{name: "self-update-too-many", json: base, cmd: "self-update 1 2"},
	}
}

func repositoryScenarios() []miscSc {
	snap := []string{"other.json", "../home/config.json"}
	base := project(`{}`, "")
	plain := "{\n    \"name\": \"test/plain\",\n    \"require\": {}\n}\n"
	named := strings.Replace(plain, `"require": {}`, `"repositories": [
        {"name": "foo", "type": "vcs", "url": "https://example.org/foo.git"},
        {"name": "bar", "type": "composer", "url": "https://repo.example.org"},
        {"type": "package", "package": {"name": "x/y", "version": "1.0.0"}}
    ],
    "require": {}`, 1)
	assoc := strings.Replace(plain, `"require": {}`, `"repositories": {
        "foo": {"type": "vcs", "url": "https://example.org/foo.git"},
        "pkg": {"type": "package", "package": {"name": "x/y", "version": "1.0.0"}},
        "packagist.org": false
    },
    "require": {}`, 1)
	r := func(name, json, cmd string) miscSc {
		return miscSc{name: "repository-" + name, json: json, cmd: cmd, snap: snap}
	}
	return []miscSc{
		r("list", base, "repository list"),
		r("list-default", plain, "repo list"),
		r("list-named", named, "repo ls"),
		r("list-assoc", assoc, "repo show"),
		r("no-args", base, "repository"),
		r("add", plain, "repo add foo vcs https://example.org/foo.git"),
		r("add-append", named, "repo add baz vcs https://example.org/baz.git --append"),
		r("add-prepend-named", named, "repo add baz composer https://baz.example.org"),
		r("add-assoc", assoc, "repo add baz vcs https://example.org/baz.git"),
		r("add-json", plain, `repo add zips '{"type": "artifact", "url": "/tmp/zips/"}'`),
		r("add-json-invalid", plain, `repo add bad '{nope'`),
		r("add-missing-name", plain, "repo add"),
		r("add-missing-arg1", plain, "repo add foo"),
		r("add-missing-url", plain, "repo add foo vcs"),
		r("add-before", named, "repo add baz vcs https://example.org/baz.git --before bar"),
		r("add-after", named, "repo add baz vcs https://example.org/baz.git --after foo"),
		r("add-before-after", named, "repo add baz vcs https://x --before foo --after bar"),
		r("add-before-missing", named, "repo add baz vcs https://x --before nope"),
		r("add-existing", named, "repo add foo vcs https://example.org/other.git"),
		r("remove", named, "repo remove foo"),
		r("remove-assoc", assoc, "repo rm foo"),
		r("remove-missing", named, "repo delete nope"),
		r("remove-packagist", plain, "repo remove packagist.org"),
		r("remove-no-name", plain, "repo remove"),
		r("set-url", named, "repo set-url foo https://git.example.org/foo.git"),
		r("set-url-assoc", assoc, "repo seturl foo https://git.example.org/foo.git"),
		r("set-url-missing", named, "repo set-url foo"),
		r("set-url-unknown", named, "repo set-url nope https://x"),
		r("get-url", named, "repo get-url bar"),
		r("get-url-assoc", assoc, "repo geturl foo"),
		r("get-url-no-url", assoc, "repo get-url pkg"),
		r("get-url-missing", named, "repo get-url nope"),
		r("get-url-no-name", named, "repo get-url"),
		r("disable", plain, "repo disable packagist.org"),
		r("disable-short", named, "repo disable packagist --append"),
		r("disable-other", named, "repo disable foo"),
		r("disable-no-name", named, "repo disable"),
		r("enable", assoc, "repo enable packagist.org"),
		r("enable-other", named, "repo enable foo"),
		r("unknown-action", plain, "repo frobnicate"),
		r("uppercase-action", plain, "repo LIST"),
		r("global-add", plain, "repo -g add foo vcs https://example.org/foo.git"),
		r("global-list", plain, "repository --global list"),
		r("file", plain, "repo --file other.json list"),
		r("file-and-global", plain, "repo -g -f other.json list"),
		r("file-missing", plain, "repo -f missing.json list"),
		{name: "repository-file-add", json: plain, cmd: "repo -f other.json add foo path ../foo", snap: snap, files: map[string]string{"other.json": "{\n    \"name\": \"other/pkg\"\n}\n"}},
		{name: "repository-no-json", nojson: true, cmd: "repo list", snap: snap},
	}
}

func policyCmdScenarios() []miscSc {
	snap := []string{"other.json", "../home/config.json"}
	plain := "{\n    \"name\": \"test/plain\",\n    \"require\": {}\n}\n"
	withPolicy := strings.Replace(plain, `"require": {}`, `"config": {"policy": {"mylist": {"sources": [{"type": "url", "url": "https://example.org/a.json"}], "block": true}}},
    "require": {}`, 1)
	objSources := strings.Replace(plain, `"require": {}`, `"config": {"policy": {"mylist": {"sources": {"first": {"type": "url", "url": "https://example.org/a.json"}}}}},
    "require": {}`, 1)
	p := func(name, json, cmd string) miscSc {
		return miscSc{name: "policy-cmd-" + name, json: json, cmd: cmd, snap: snap}
	}
	return []miscSc{
		p("add-source", plain, "policy add-source mylist url https://example.org/list.json"),
		p("add-source-json", plain, `policy add-source mylist '{"type": "url", "url": "https://example.org/list.json"}'`),
		p("add-source-existing-list", withPolicy, "policy add-source mylist url https://example.org/b.json"),
		p("add-source-duplicate", withPolicy, "policy add-source mylist url https://example.org/a.json"),
		p("add-source-object-sources", objSources, "policy add-source mylist url https://example.org/b.json"),
		p("add-source-uppercase-action", plain, "policy ADD-SOURCE other url https://example.org/x.json"),
		p("no-action", plain, "policy"),
		p("unknown-action", plain, "policy list"),
		p("no-name", plain, "policy add-source"),
		p("builtin", plain, "policy add-source advisories url https://example.org/x.json"),
		p("builtin-cooldown", plain, "policy add-source cooldown url https://example.org/x.json"),
		p("reserved-prefix", plain, "policy add-source ignore-me url https://example.org/x.json"),
		p("reserved-name", plain, "policy add-source security url https://example.org/x.json"),
		p("dotted-name", plain, "policy add-source my.list url https://example.org/x.json"),
		p("empty-name", plain, "policy add-source '' url https://example.org/x.json"),
		p("no-arg1", plain, "policy add-source mylist"),
		p("no-url", plain, "policy add-source mylist url"),
		p("bad-type", plain, "policy add-source mylist file https://example.org/x.json"),
		p("http-url", plain, "policy add-source mylist url http://example.org/x.json"),
		p("json-no-type", plain, `policy add-source mylist '{"url": "https://example.org/x.json"}'`),
		p("json-no-url", plain, `policy add-source mylist '{"type": "url"}'`),
		p("json-int-url", plain, `policy add-source mylist '{"type": "url", "url": 5}'`),
		p("json-invalid", plain, `policy add-source mylist '{bad'`),
		p("global", plain, "policy -g add-source mylist url https://example.org/x.json"),
		p("file-and-global", plain, "policy -g -f other.json add-source mylist url https://example.org/x.json"),
		{name: "policy-cmd-file", json: plain, cmd: "policy --file other.json add-source mylist url https://example.org/x.json", snap: snap, files: map[string]string{"other.json": "{\n    \"name\": \"other/pkg\"\n}\n"}},
		{name: "policy-cmd-no-json", nojson: true, cmd: "policy add-source mylist url https://example.org/x.json", snap: snap},
	}
}
