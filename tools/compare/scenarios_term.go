package main

// Terminal display, transport error and super user scenarios (term-*,
// curl-*, root-*) of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80:
//
//   - term-href-*: decorated output with Symfony OutputFormatter <href=...>
//     tags (show, licenses, fund, search, config --list, audit) rendered as
//     OSC 8 hyperlinks, and the terminals OutputFormatterStyle::apply leaves
//     them out for (JetBrains JediTerm, old Konsole, IDEA_INITIAL_DIRECTORY).
//   - term-progress-*: the Symfony ProgressBar InstallationManager draws
//     while downloading and extracting, overwritten in place on a terminal.
//   - curl-*: CurlDownloader's "curl error N while downloading" messages for
//     loopback servers that fail in one specific way
//     (testdata/compare/vcs/term-fault-server.php).
//   - root-*: Application::doRunComposer's super user checks, run as uid 0.

import (
	"regexp"
	"strings"
)

func init() {
	extraScenarios = append(extraScenarios, termScenarios)
	extraOutputNormalizers = append(extraOutputNormalizers, normalizeTermTiming)
}

// faultStubScript starts the fault servers of the termfault fixture (ports
// 9201-9213, see termfault.sh and term-fault-server.php).
const faultStubScript = `if [ -f /fixtures/git/_fault/server.php ]; then p=9201; for m in hang close reset garbage partial tls-close tls-garbage tls-hang tls-self-signed tls-wrong-host tls-expired tls-unknown-ca tls-good; do php /fixtures/git/_fault/server.php $m $p /fixtures/git/_fault >/dev/null 2>&1 & p=$((p+1)); done
i=0; while ! php -r 'foreach (range(9201, 9213) as $p) { if (!@fsockopen("127.0.0.1", $p)) exit(1); }' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done; fi`

var (
	// libcurl reports how long a timed out transfer ran
	reCurlTimeout = regexp.MustCompile(`(timed out after )\d+( milliseconds)`)
	// intermediate (throttled, timing dependent) frames of a progress bar
	// overwritten in place on a terminal; the first and the final frame
	// are always drawn and stay compared
	reTermProgressMid = regexp.MustCompile(`\\e\[1G\\e\[2K +\d+/\d+ \[=*>?-*\] +([1-9]|[1-9]\d)%[^\\\n]*`)
	// the %elapsed:6s% field of the verbose progress formats
	reTermElapsed = regexp.MustCompile(`(\[[=>-]+\] +\d+%) +(?:< 1 sec|1 sec|\d+ secs)`)
	// ConsoleIO::writeProgress reports of Loop::runJobs: the intermediate
	// percentages are throttled (0.1s) and whether the last update already
	// saw every job done (an extra 100 before finish()'s) depends on timing
	reTermStatusMid  = regexp.MustCompile(`\\e\]7501;state=working:app=composer:progress=(?:\d|[1-9]\d)\\e\\`)
	reTermStatus100s = regexp.MustCompile(`(\\e\]7501;state=working:app=composer:progress=100\\e\\)+`)
)

// normalizeTermTiming removes the time dependent parts the term-* and curl-*
// scenarios show (the raw output stays in the result JSON).
func normalizeTermTiming(s string) string {
	s = reCurlTimeout.ReplaceAllString(s, "${1}N${2}")
	s = reTermProgressMid.ReplaceAllString(s, "")
	s = reTermStatusMid.ReplaceAllString(s, "")
	s = reTermStatus100s.ReplaceAllString(s, "${1}")
	return reTermElapsed.ReplaceAllString(s, "${1} <elapsed>")
}

func termScenarios() []Scenario {
	var out []Scenario
	out = append(out, termHrefScenarios()...)
	out = append(out, termProgressScenarios()...)
	out = append(out, curlScenarios()...)
	out = append(out, rootScenarios()...)
	return out
}

func termHrefScenarios() []Scenario {
	deps := project(`{"fundco/alpha": "^1.0", "fundco/beta": "^1.0", "fundco/plain": "^1.0", "otherfund/gamma": "^1.0"}`, "")
	setup := []string{"update --no-progress"}
	files := map[string]string{"src/Kernel.php": appSrc}
	vuln := withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, ""))
	auditSetup := []string{"update --no-progress --no-blocking"}
	t := func(name, json string, setup []string, cmd string, env map[string]string) Scenario {
		return Scenario{Name: "term-href-" + name, JSON: json, Files: files, Setup: setup, Command: cmd, TTY: true, Env: env}
	}
	return []Scenario{
		t("show", deps, setup, "show", nil),
		t("show-name", deps, setup, "show fundco/alpha", nil),
		t("show-latest", deps, setup, "show --latest", nil),
		t("licenses", deps, setup, "licenses", nil),
		t("fund", deps, setup, "fund", nil),
		t("search", deps, nil, "search fundco", nil),
		t("depends", deps, setup, "depends fundco/plain", nil),
		t("config-list", deps, nil, "config --list", nil),
		t("audit", vuln, auditSetup, "audit", nil),
		t("audit-plain", vuln, auditSetup, "audit --format=plain", nil),
		t("jetbrains", deps, setup, "show", map[string]string{"TERMINAL_EMULATOR": "JetBrains-JediTerm"}),
		t("konsole-old", deps, setup, "show", map[string]string{"KONSOLE_VERSION": "201100"}),
		t("konsole-new", deps, setup, "show", map[string]string{"KONSOLE_VERSION": "220400"}),
		t("konsole-zero", deps, setup, "show", map[string]string{"KONSOLE_VERSION": "0"}),
		t("idea", deps, setup, "show", map[string]string{"IDEA_INITIAL_DIRECTORY": "/tmp"}),
		t("no-ansi", deps, setup, "show --no-ansi", nil),
		// decorated output on a pipe
		{Name: "term-href-ansi-pipe", JSON: deps, Files: files, Setup: setup, Command: "show --ansi"},
		{Name: "term-href-ansi-pipe-fund", JSON: deps, Files: files, Setup: setup, Command: "fund --ansi"},
	}
}

func termProgressScenarios() []Scenario {
	basic := project(`{"php": ">=7.4", "acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0"}`, `"require-dev": {"acme/devtools": "^1.0"}`)
	files := map[string]string{"src/Kernel.php": appSrc}
	lockOnly := []string{"update --no-progress --no-install"}
	p := func(name string, setup []string, cmd string, env map[string]string) Scenario {
		return Scenario{Name: "term-progress-" + name, JSON: basic, Files: files, Setup: setup, Command: cmd, TTY: true, Env: env, Probe: true}
	}
	return []Scenario{
		p("install", lockOnly, "install", nil),
		p("update", nil, "update", nil),
		p("require", lockOnly, "require acme/util:^1.0", nil),
		p("verbose", lockOnly, "install -v", nil),
		p("very-verbose", lockOnly, "install -vv", nil),
		p("narrow", lockOnly, "install", map[string]string{"COLUMNS": "30"}),
		p("tiny", lockOnly, "install", map[string]string{"COLUMNS": "5"}),
		p("no-progress", lockOnly, "install --no-progress", nil),
		p("no-ansi", lockOnly, "install --no-ansi", nil),
		p("ci", lockOnly, "install", map[string]string{"CI": "1"}),
		// decorated output on a pipe (overwritten frames in a file)
		{Name: "term-progress-ansi-pipe", JSON: basic, Files: files, Setup: lockOnly, Command: "install --ansi", Probe: true},
	}
}

// faultRepo is a project whose only repository is a fault server.
func faultRepo(url, config string) string {
	s := strings.Replace(project(`{"acme/log": "^1.0"}`, ""), `"url": "http://127.0.0.1:8899"}`, `"url": "`+url+`"}`, 1)
	if config != "" {
		s = withConfig(s, config)
	}
	return s
}

func curlScenarios() []Scenario {
	fx := []string{"termfault"}
	ca := `"cafile": "/fixtures/git/_fault/ca.pem"`
	c := func(name, url, config, cmd string) Scenario {
		return Scenario{Name: "curl-" + name, JSON: faultRepo(url, config), GitFixtures: fx, Command: cmd}
	}
	files := map[string]string{"src/Kernel.php": appSrc}
	vuln := withAdv(project(`{"acme/http": "^1.2", "acme/log-impl": "^1.0", "acme/util": "1.0.0"}`, ""))
	// audit against a repository whose advisory api-url misbehaves (a POST
	// with Composer's 10 second http.timeout)
	api := func(name, repo, cmd string) Scenario {
		after := strings.Replace(vuln, `"url": "http://127.0.0.1:8899/v2adv"}`, `"url": "http://127.0.0.1:8899/git/_fault/`+repo+`"}`, 1)
		return Scenario{Name: "curl-" + name, JSON: vuln, Files: files, GitFixtures: fx, Setup: []string{"update --no-progress --no-blocking"},
			After: map[string]string{"composer.json": after}, Command: cmd}
	}
	return []Scenario{
		c("refused", "http://127.0.0.1:9", "", "update --no-progress"),
		c("empty-reply", "http://127.0.0.1:9202", "", "update --no-progress"),
		c("reset", "http://127.0.0.1:9203", "", "update --no-progress"),
		c("http09", "http://127.0.0.1:9204", "", "update --no-progress"),
		c("partial", "http://127.0.0.1:9205", "", "update --no-progress"),
		c("tls-eof", "https://127.0.0.1:9206", ca, "update --no-progress"),
		c("tls-wrong-version", "https://127.0.0.1:9207", ca, "update --no-progress"),
		c("tls-reset", "https://127.0.0.1:9203", ca, "update --no-progress"),
		c("tls-self-signed", "https://127.0.0.1:9209", ca, "update --no-progress"),
		c("tls-self-signed-system-ca", "https://127.0.0.1:9209", "", "update --no-progress"),
		c("tls-wrong-host", "https://127.0.0.1:9210", ca, "update --no-progress"),
		c("tls-expired", "https://127.0.0.1:9211", ca, "update --no-progress"),
		c("tls-unknown-ca", "https://127.0.0.1:9212", ca, "update --no-progress"),
		c("tls-good", "https://127.0.0.1:9213", ca, "update --no-progress"),
		c("tls-show", "https://127.0.0.1:9209", ca, "show -a acme/log"),
		api("timeout", "hang-api", "audit"),
		api("tls-timeout", "hang-tls-api", "audit --format=plain"),
		api("post-empty-reply", "close-api", "audit"),
		api("post-timeout-ignored", "hang-api", "audit --ignore-unreachable"),
	}
}

func rootScenarios() []Scenario {
	base := project(`{"acme/log": "^1.0"}`, `"scripts": {"hello": "echo hello"}`)
	files := map[string]string{"src/Kernel.php": appSrc}
	q := "Continue as root/super user"
	r := func(name, cmd string, tty bool, env map[string]string, input ...TTYStep) Scenario {
		return Scenario{Name: "root-" + name, JSON: base, Files: files, Command: cmd, Root: true, TTY: tty, Env: env, Input: input}
	}
	allow := map[string]string{"COMPOSER_ALLOW_SUPERUSER": "1"}
	return []Scenario{
		r("validate", "validate", false, nil),
		r("validate-allowed", "validate", false, allow),
		r("validate-allowed-zero", "validate", false, map[string]string{"COMPOSER_ALLOW_SUPERUSER": "0"}),
		r("validate-no-plugins", "validate --no-plugins", false, nil),
		r("update", "update --no-progress", false, nil),
		// an unknown command name may be a plugin command: plugins are
		// needed before the name can be resolved
		r("unknown-command", "frobnicate", false, nil),
		r("unknown-command-allowed", "frobnicate", false, allow),
		r("help", "help validate", false, nil),
		r("script-alias", "hello", false, nil),
		r("tty-validate-yes", "validate", true, nil, Answer(q, "y")),
		r("tty-validate-default", "validate", true, nil, Enter(q)),
		r("tty-validate-no", "validate", true, nil, Answer(q, "n")),
		r("tty-validate-n", "validate -n", true, nil),
		r("tty-validate-allowed", "validate", true, allow),
		r("tty-unknown-command-yes", "frobnicate", true, nil, Answer(q, "yes")),
		r("tty-unknown-command-no", "frobnicate", true, nil, Answer(q, "no")),
		r("tty-script-alias-yes", "hello", true, nil, Answer(q, "y")),
		r("tty-update-yes", "update --no-progress", true, nil, Answer(q, "y")),
		r("tty-no-plugins", "validate --no-plugins", true, nil, Answer(q, "y")),
	}
}
