package main

import "regexp"

// Output of list/help (Symfony Console 5.4 ListCommand, HelpCommand and the
// txt/json/xml/md descriptors as used by Composer
// f83588392df43d4fd027d84d26de459946ef4e80): as root, on a terminal, for
// script and plugin commands, namespaces, --raw/--short/--format and the
// --version/--help flags.
func finalMiscScenarios() []Scenario {
	json := project(`{"acme/log": "^1.0"}`, `"scripts": {"greet": "echo hi", "both": ["@greet", "echo done"]}, "scripts-descriptions": {"greet": "Says hi"}`)
	nsJSON := project(`{"acme/log": "^1.0"}`, `"scripts": {"greet": "echo hi", "lint:php": "echo php", "lint:js": "echo js", "db:seed:all": "echo seed"}, "scripts-descriptions": {"lint:js": "Lints <info>JS</info> files"}, "scripts-aliases": {"lint:php": ["lp", "lint-php"]}`)
	files := map[string]string{"src/Kernel.php": appSrc}
	allow := map[string]string{"COMPOSER_ALLOW_SUPERUSER": "1"}
	plugin := withConfig(project(`{"plugfix/commands-plugin": "^1.0"}`, ""), `"allow-plugins": {"plugfix/*": true}`)
	pluginSetup := []string{"update --no-progress"}
	sc := func(name, cmd string) Scenario {
		return Scenario{Name: "misc-" + name, JSON: json, Files: files, Command: cmd}
	}
	ns := func(name, cmd string) Scenario {
		return Scenario{Name: "misc-" + name, JSON: nsJSON, Files: files, Command: cmd}
	}
	pl := func(name, cmd string) Scenario {
		return Scenario{Name: "misc-" + name, JSON: plugin, Files: files, Setup: pluginSetup, Command: cmd}
	}
	return []Scenario{
		{Name: "misc-list-root", JSON: json, Files: files, Command: "list", Root: true, Env: allow},
		{Name: "misc-list-root-noallow", JSON: json, Files: files, Command: "list -n", Root: true},
		{Name: "misc-help-root", JSON: json, Files: files, Command: "help", Root: true, Env: allow},
		{Name: "misc-list-tty", JSON: json, Files: files, Command: "list", TTY: true},
		{Name: "misc-help-tty", JSON: json, Files: files, Command: "help", TTY: true},
		{Name: "misc-help-install-tty", JSON: json, Files: files, Command: "help install", TTY: true},
		{Name: "misc-help-script", JSON: json, Files: files, Command: "help greet"},
		{Name: "misc-help-script-list", JSON: json, Files: files, Command: "help both"},
		{Name: "misc-help-script-tty", JSON: json, Files: files, Command: "help greet", TTY: true},
		{Name: "misc-list-scripts", JSON: json, Files: files, Command: "list"},
		// application flags
		sc("version", "--version"),
		sc("version-short", "-V"),
		sc("version-ansi", "--version --ansi"),
		sc("no-command", "--no-ansi"),
		sc("h-only", "-h"),
		sc("h-only-json", "-h --format=json"),
		// list
		sc("list-raw", "list --raw"),
		sc("list-ansi", "list --ansi"),
		sc("list-quiet", "list -q"),
		sc("list-json", "list --format=json"),
		sc("list-json-short", "list --format=json --short"),
		sc("list-xml", "list --format=xml"),
		sc("list-xml-short", "list --format=xml --short"),
		sc("list-md", "list --format=md"),
		sc("list-md-short", "list --format=md --short"),
		sc("list-txt-short", "list --format=txt --short"),
		sc("list-bad-format", "list --format=yaml"),
		sc("list-help", "list --help"),
		sc("list-missing-namespace", "list nope"),
		ns("list-ns-scripts", "list"),
		ns("list-ns", "list lint"),
		ns("list-ns-raw", "list lint --raw"),
		ns("list-ns-nested", "list db"),
		ns("list-ns-abbrev", "list li"),
		ns("list-ns-json", "list lint --format=json"),
		ns("list-ns-xml", "list lint --format=xml"),
		ns("list-ns-md", "list lint --format=md"),
		ns("list-ns-scripts-raw", "list --raw"),
		ns("list-ns-scripts-xml", "list --format=xml"),
		// help
		sc("help-list", "help list"),
		sc("help-raw", "help install --raw"),
		sc("help-json", "help install --format=json"),
		sc("help-xml", "help install --format=xml"),
		sc("help-md", "help install --format=md"),
		sc("help-md-ansi", "help validate --format=md --ansi"),
		sc("help-txt-ansi", "help require --ansi"),
		sc("help-bad-format", "help --format=yaml"),
		sc("help-missing", "help nope"),
		sc("help-abbrev", "help inst"),
		sc("help-ambiguous", "help re"),
		sc("help-alias", "help why"),
		sc("help-hidden", "help _complete"),
		sc("help-completion", "help completion"),
		sc("help-config", "help config"),
		sc("help-show-json", "help show --format=json"),
		sc("help-script-json", "help greet --format=json"),
		sc("help-script-xml", "help greet --format=xml"),
		sc("help-script-md", "help greet --format=md"),
		sc("help-flag", "install --help"),
		sc("help-flag-short", "validate -h"),
		sc("help-flag-json", "install --help --format=json"),
		sc("help-flag-raw", "remove -h --raw"),
		sc("help-flag-script", "greet --help"),
		ns("help-ns-script", "help lint:js"),
		ns("help-ns-alias", "help lp"),
		// plugin commands
		pl("list-plugin", "list"),
		pl("list-plugin-ns", "list plugfix"),
		pl("list-plugin-json", "list --format=json"),
		pl("list-plugin-raw", "list --raw"),
		pl("help-plugin", "help plugfix:echo"),
		pl("help-plugin-json", "help plugfix:echo --format=json"),
		pl("help-plugin-xml", "help plugfix:echo --format=xml"),
		pl("help-plugin-md", "help plugfix:echo --format=md"),
		pl("help-plugin-flag", "plugfix:echo --help"),
	}
}

// The application banner differs by design: Composer prints its ASCII logo
// and "Composer version <version> <date>" (Application::getHelp and
// getLongVersion), gocomposer prints its own name and version line and does
// not claim to be Composer. Both banners, and the application name/version
// fields of the json/xml/md descriptors, are replaced by the same
// placeholders; the raw outputs are kept in the result files. The
// %command.full_name% / completion help paths ($_SERVER['PHP_SELF'] vs
// os.Args[0]) are normalized the same way.
var (
	reAnsiOpt = `(?:(?:\x1b|\\e)\[[0-9;]*m)*`
	// a terminal transcript may start with the OSC 9;4/7501 program status
	reOSCPrefix       = `((?:(?:\x1b|\\e)\][0-9;]+[^\x1b\\\n]*(?:\x1b|\\e)\\)*)`
	reComposerBanner  = regexp.MustCompile(`(?m)^` + reOSCPrefix + ` {3}______\n {2}/ ____/___  ____ ___  ____  ____  ________  _____\n / /   / __ \\/ __ ` + "`" + `__ \\/ __ \\/ __ \\/ ___/ _ \\/ ___/\n/ /___/ /_/ / / / / / / /_/ / /_/ \(__  \)  __/ /\n\\____/\\____/_/ /_/ /_/ \.___/\\____/____/\\___/_/\n {20}/_/\n` + reAnsiOpt + `Composer` + reAnsiOpt + ` version ` + reAnsiOpt + `2\.11\.999-dev\+source` + reAnsiOpt + ` @release_date@$`)
	reComposerVersion = regexp.MustCompile(`(?m)^` + reOSCPrefix + reAnsiOpt + `Composer` + reAnsiOpt + ` version ` + reAnsiOpt + `2\.11\.999-dev\+source` + reAnsiOpt + ` @release_date@$`)
	reGoBanner        = regexp.MustCompile(`(?m)^` + reOSCPrefix + reAnsiOpt + `gocomposer` + reAnsiOpt + ` version ` + reAnsiOpt + `[^ \x1b\\]+` + reAnsiOpt + ` \(Composer-compatible, based on Composer 2\.11\.999-dev\+source [0-9a-f]{40}\)$`)
	reAppJSON         = regexp.MustCompile(`"application":\{"name":"(?:Composer|gocomposer)","version":"[^"]*"\}`)
	reAppXML          = regexp.MustCompile(`<symfony name="(?:Composer|gocomposer)" version="[^"]*">`)
	reAppMD           = regexp.MustCompile(`(?m)^(?:Composer 2\.11\.999-dev\+source|gocomposer \S+)\n=+$`)
	rePHPSelf         = regexp.MustCompile(`(?:/oracle-source/bin/composer|/gobin/gocomposer|\\/oracle-source\\/bin\\/composer|\\/gobin\\/gocomposer)\b`)
	reCompletionDir   = regexp.MustCompile(`(/etc/bash_completion\.d/|\\/etc\\/bash_completion\.d\\/)(?:gocomposer|composer)\b`)
)

func normalizeBanner(s string) string {
	s = reComposerBanner.ReplaceAllString(s, "${1}<APPLICATION BANNER>")
	s = reComposerVersion.ReplaceAllString(s, "${1}<APPLICATION BANNER>")
	s = reGoBanner.ReplaceAllString(s, "${1}<APPLICATION BANNER>")
	s = reAppJSON.ReplaceAllString(s, `"application":{"name":"<APP>","version":"<VERSION>"}`)
	s = reAppXML.ReplaceAllString(s, `<symfony name="<APP>" version="<VERSION>">`)
	s = reAppMD.ReplaceAllString(s, "<APP> <VERSION>\n=")
	s = rePHPSelf.ReplaceAllString(s, "<PHP_SELF>")
	return reCompletionDir.ReplaceAllString(s, "${1}<NAME>")
}

func init() {
	extraScenarios = append(extraScenarios, finalMiscScenarios)
	extraOutputNormalizers = append(extraOutputNormalizers, normalizeBanner)
}
