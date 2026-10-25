package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Shell completion and command line details of Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (third stage, items 2 and 3):
//   - completion-*: the "completion" command (Symfony Console
//     DumpCompletionCommand, Resources/completion.bash);
//   - complete-*: the hidden "_complete" command called by the bash script
//     (CompleteCommand, CompletionInput) with Composer's suggestions
//     (CompletionTrait, ConfigCommand::suggestSettingKeys, ...);
//   - cli3-*: -vvv debug output of Factory/Config, the bin-compat
//     deprecation, command events of validate/reinstall, run-script on
//     Symfony Command classes and fund/search/browse without composer.json
//     against a loopback packagist.org stub.

func init() {
	extraScenarios = append(extraScenarios, finalCLIScenarios)
	extraOutputNormalizers = append(extraOutputNormalizers, normalizeCompletionScript)
}

// The completion script names its bash function and the completed command
// after basename($argv[0]): "composer" for the oracle (bin/composer),
// "gocomposer" for gocomposer. Only these two script lines are normalized.
var (
	reCompletionFunc     = regexp.MustCompile(`(?m)^_sf_gocomposer\(\) \{$`)
	reCompletionRegister = regexp.MustCompile(`(?m)^complete -F _sf_gocomposer gocomposer$`)
)

// composerUserDeprecations are the E_USER_DEPRECATED notices Composer
// itself triggers that gocomposer reproduces; unlike PHP runtime
// deprecations they are compared (dropOracleRuntimeDeprecations keeps them),
// with the PHP file and line normalized like exception locations.
var composerUserDeprecations = []string{`config.bin-compat "symlink" is deprecated since Composer 2.2`}

func isComposerUserDeprecation(line string) bool {
	for _, d := range composerUserDeprecations {
		if strings.Contains(line, d) {
			return true
		}
	}
	return false
}

var reDeprecationLocation = regexp.MustCompile(`(?m)^(Deprecation Notice: .*) in (?:/oracle-source/src/\S+:\d+|gocomposer)$`)

func normalizeCompletionScript(s string) string {
	s = reDeprecationLocation.ReplaceAllString(s, "${1} in <location>")
	s = reCompletionFunc.ReplaceAllString(s, "_sf_composer() {")
	return reCompletionRegister.ReplaceAllString(s, "complete -F _sf_composer composer")
}

func finalCLIScenarios() []Scenario {
	var out []Scenario
	out = append(out, completionScenarios()...)
	out = append(out, cli3Scenarios()...)
	return out
}

func completionScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	upd := []string{"update --no-progress"}
	extra := `"require-dev": {"acme/util": "^1.0"},
    "bin": ["bin/root-tool"],
    "scripts": {"test": "echo t", "build": ["@test"], "post-update-cmd": "echo done"},
    "scripts-aliases": {"build": ["b"]},
    "extra": {"x": {"y": 1}, "l": [1, 2]},
    "suggest": {"acme/http": "for http"}`
	base := withConfig(project(`{"acme/log": "^1.0", "acme/cli": "^1.0", "php": ">=7.4"}`, extra), `"platform": {"php": "8.1.0"}, "sort-packages": true`)
	withRepos := strings.Replace(base, `{"packagist.org": false}`, `{"packagist.org": false}, {"name": "zeta", "type": "path", "url": "pkgs/zeta"}, {"name": "alpha", "type": "path", "url": "pkgs/alpha"}`, 1)
	withPolicy := withConfig(base, `"policy": {"internal": {"block": true}, "advisories": true, "ignore-unreachable": true, "abc": {}}`)
	home := map[string]string{"COMPOSER_HOME": "/work/project/ghome"}
	globalHome := `{
    "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}],
    "config": {"secure-http": false},
    "require": {"acme/util": "^1.0"}
}
`
	cmp := func(name, args string) Scenario {
		return Scenario{Name: "complete-" + name, JSON: base, Files: files, Setup: upd, Command: "_complete -sbash " + args}
	}
	sc := []Scenario{
		{Name: "completion-bash", JSON: base, Command: "completion bash"},
		{Name: "completion-shell-env", JSON: base, Env: map[string]string{"SHELL": "/bin/bash"}, Command: "completion"},
		{Name: "completion-no-shell", JSON: base, Command: "completion"},
		{Name: "completion-zsh", JSON: base, Command: "completion zsh"},
		{Name: "completion-fish-env", JSON: base, Env: map[string]string{"SHELL": "/usr/bin/fish"}, Command: "completion"},
		{Name: "completion-quiet", JSON: base, Command: "completion bash -q"},

		// command names and option names
		cmp("commands", "-c1 -icomposer"),
		cmp("commands-prefix", "-c1 -icomposer -ius"),
		cmp("commands-abbrev", "-c1 -icomposer -iup"),
		cmp("commands-alias", "-c1 -icomposer -iwhy"),
		cmp("commands-exact", "-c1 -icomposer -iupdate"),
		cmp("app-options", "-c1 -icomposer -i--"),
		cmp("command-options", "-c2 -icomposer -iupdate -i--"),
		cmp("command-options-short", "-c2 -icomposer -irequire -i-"),
		cmp("working-dir-first", "-c3 -icomposer -i-d -i/work/project -isho"),
		cmp("verbose-first", "-c2 -icomposer -i-vvv -iinst"),
		cmp("after-separator", "-c3 -icomposer -iupdate -i-- -i--dev"),
		cmp("none", "-c3 -icomposer -ivalidate -icomposer.json"),
		// option values
		cmp("format", "-c3 -icomposer -ishow -i--format"),
		cmp("format-eq", "-c2 -icomposer -ishow -i--format=j"),
		cmp("format-short", "-c2 -icomposer -iarchive -i-f"),
		cmp("prefer-install", "-c3 -icomposer -iinstall -i--prefer-install"),
		cmp("audit-format", "-c3 -icomposer -iremove -i--audit-format"),
		cmp("bump-after-update", "-c2 -icomposer -iupdate -i--bump-after-update="),
		cmp("audit-options", "-c3 -icomposer -iaudit -i--ignore-severity"),
		cmp("audit-abandoned", "-c3 -icomposer -iaudit -i--abandoned"),
		cmp("licenses-format", "-c3 -icomposer -ilicenses -i--format"),
		cmp("fund-format", "-c3 -icomposer -ifund -i--format"),
		cmp("reinstall-type", "-c3 -icomposer -ireinstall -i--type"),
		cmp("show-ignore", "-c3 -icomposer -ishow -i--ignore"),
		cmp("init-require", "-c3 -icomposer -iinit -i--require -iacme/"),
		cmp("help-format", "-c3 -icomposer -ihelp -i--format"),
		cmp("list-format", "-c3 -icomposer -ilist -i--format"),
		cmp("global-option-value", "-c2 -icomposer -i--working-dir"),
		// argument values
		cmp("update-packages", "-c2 -icomposer -iupdate"),
		cmp("update-packages-prefix", "-c3 -icomposer -iupdate -iacme/log -iacme/c"),
		cmp("remove-packages", "-c2 -icomposer -iremove"),
		cmp("bump-packages", "-c2 -icomposer -ibump"),
		cmp("depends-package", "-c2 -icomposer -idepends"),
		cmp("depends-platform", "-c2 -icomposer -iwhy -iext-"),
		cmp("prohibits-package", "-c2 -icomposer -iprohibits"),
		cmp("require-vendors", "-c2 -icomposer -irequire"),
		cmp("require-vendor", "-c2 -icomposer -irequire -iacme"),
		cmp("require-names", "-c2 -icomposer -irequire -iacme/l"),
		cmp("require-platform", "-c2 -icomposer -irequire -iext-js"),
		cmp("require-php", "-c2 -icomposer -ir -iphp"),
		cmp("create-project-package", "-c2 -icomposer -icreate-project -iac"),
		cmp("archive-package", "-c2 -icomposer -iarchive -iacme/h"),
		cmp("show-installed", "-c2 -icomposer -ishow"),
		cmp("show-platform", "-c3 -icomposer -ishow -i--platform -ilib-"),
		cmp("show-available", "-c3 -icomposer -ishow -i--available -iacme/"),
		cmp("outdated-package", "-c2 -icomposer -ioutdated"),
		cmp("suggests-packages", "-c2 -icomposer -isuggests"),
		cmp("browse-packages", "-c2 -icomposer -ihome"),
		cmp("reinstall-packages", "-c2 -icomposer -ireinstall"),
		cmp("run-script", "-c2 -icomposer -irun-script"),
		cmp("run-script-alias", "-c2 -icomposer -irun"),
		cmp("exec-binary", "-c2 -icomposer -iexec"),
		cmp("help-command", "-c2 -icomposer -ihelp"),
		cmp("list-namespace", "-c2 -icomposer -ilist"),
		cmp("completion-shell", "-c2 -icomposer -icompletion"),
		cmp("global-command", "-c2 -icomposer -iglobal"),
		cmp("global-options", "-c3 -icomposer -iglobal -ishow -i--format"),
		cmp("repository-action", "-c2 -icomposer -irepository"),
		cmp("repository-enable", "-c3 -icomposer -irepo -ienable"),
		cmp("repository-add-type", "-c4 -icomposer -irepository -iadd -ifoo"),
		cmp("policy-action", "-c2 -icomposer -ipolicy"),
		cmp("policy-arg1", "-c4 -icomposer -ipolicy -iadd-source -ix"),
		cmp("script-command", "-c2 -icomposer -ibuild"),
		// config setting keys
		cmp("config-keys", "-c2 -icomposer -iconfig"),
		cmp("config-keys-prefix", "-c2 -icomposer -iconfig -ipla"),
		cmp("config-keys-extra", "-c2 -icomposer -iconfig -iex"),
		cmp("config-keys-unset", "-c3 -icomposer -iconfig -i--unset"),
		cmp("config-keys-global", "-c3 -icomposer -iconfig -i--global"),
		cmp("config-keys-list", "-c3 -icomposer -iconfig -i--list"),
		// failures (exit 2 without output, or the missing composer.json)
		cmp("no-shell", "-c1 -icomposer"),
		cmp("bad-shell", "-szsh -c1 -icomposer"),
		cmp("no-current", "-icomposer"),
		cmp("bad-current", "-c1x -icomposer"),
		cmp("index-out-of-range", "-c5 -icomposer -iupdate"),
		cmp("missing-input-value", "-c2 -icomposer -iupdate -i"),
		{Name: "complete-repository-names", JSON: withRepos, Files: files, Setup: upd, Command: "_complete -sbash -c3 -icomposer -irepository -iremove"},
		{Name: "complete-repository-before", JSON: withRepos, Files: files, Setup: upd, Command: "_complete -sbash -c5 -icomposer -irepository -iremove -ix -i--before"},
		{Name: "complete-policy-names", JSON: withPolicy, Files: files, Setup: upd, Command: "_complete -sbash -c3 -icomposer -ipolicy -iadd-source"},
		{Name: "complete-unlocked-installed", JSON: base, Files: files, Command: "_complete -sbash -c2 -icomposer -iupdate"},
		{Name: "complete-no-json-command", NoJSON: true, Command: "_complete -sbash -c1 -icomposer -ire"},
		{Name: "complete-no-json-update", NoJSON: true, Command: "_complete -sbash -c2 -icomposer -iupdate"},
		{Name: "complete-no-json-config", NoJSON: true, Command: "_complete -sbash -c2 -icomposer -iconfig -ire"},
		{Name: "complete-global-require", JSON: base, Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: home, Setup: []string{"global update --no-progress"}, Command: "_complete -sbash -c3 -icomposer -iglobal -iremove"},
		{Name: "complete-global-show", JSON: base, Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: home, Setup: []string{"global update --no-progress"}, Command: "_complete -sbash -c4 -icomposer -iglobal -ishow -iacme"},
		{Name: "complete-global-no-home-json", JSON: base, Env: home, Command: "_complete -sbash -c3 -icomposer -iglobal -iupdate"},
		{Name: "complete-plugin-commands", JSON: withConfig(project(`{"plugfix/commands-plugin": "^1.0"}`, ""), `"allow-plugins": {"plugfix/*": true}`), Files: files, Setup: upd, Command: "_complete -sbash -c1 -icomposer -iplug"},
		{Name: "complete-plugin-namespace", JSON: withConfig(project(`{"plugfix/commands-plugin": "^1.0"}`, ""), `"allow-plugins": {"plugfix/*": true}`), Files: files, Setup: upd, Command: "_complete -sbash -c2 -icomposer -ilist"},
		{Name: "complete-plugin-options", JSON: withConfig(project(`{"plugfix/commands-plugin": "^1.0"}`, ""), `"allow-plugins": {"plugfix/*": true}`), Files: files, Setup: upd, Command: "_complete -sbash -c2 -icomposer -iplugfix:echo -i--"},
		{Name: "complete-plugin-no-plugins", JSON: withConfig(project(`{"plugfix/commands-plugin": "^1.0"}`, ""), `"allow-plugins": {"plugfix/*": true}`), Files: files, Setup: upd, Command: "_complete --no-plugins -sbash -c1 -icomposer -iplug"},
	}
	return sc
}

// versioned is a fixture project with a root version: Composer's
// VersionGuesser (git/hg/fossil/svn commands in -vvv output) is not run.
func versioned(require, extra string) string {
	return strings.Replace(project(require, extra), `"type": "project",`, `"type": "project",
    "version": "1.0.0",`, 1)
}

const helloCommandPHP = `<?php
namespace App;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

class HelloCommand extends Command
{
    protected function configure(): void
    {
        $this->setDescription('Says hello')
            ->addArgument('who', InputArgument::OPTIONAL, 'Who to greet', 'world')
            ->addOption('shout', 's', InputOption::VALUE_NONE, 'Shout the greeting');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $message = 'Hello '.$input->getArgument('who').' from '.$this->getName();
        if ($input->getOption('shout')) {
            $message = strtoupper($message);
        }
        $output->writeln('<info>'.$message.'</info>');
        $output->writeln('verbose detail', OutputInterface::VERBOSITY_VERBOSE);

        return 0;
    }
}
`

const failCommandPHP = `<?php
namespace App;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class FailCommand extends Command
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('failing on purpose');

        return 3;
    }
}
`

const boomCommandPHP = `<?php
namespace App;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class BoomCommand extends Command
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        throw new \RuntimeException('boom from the command');
    }
}
`

const namedCommandPHP = `<?php
namespace App;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class NamedCommand extends Command
{
    protected function configure(): void
    {
        $this->setName('other-name');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('named command ran as '.$this->getName());

        return 0;
    }
}
`

func cli3Scenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	upd := []string{"update --no-progress"}
	vvvJSON := versioned(`{"acme/log": "^1.0"}`, `"scripts": {"greet": "echo hi"}, "extra": {"x": 1}`)
	ghome := map[string]string{"COMPOSER_HOME": "/work/project/ghome"}
	homeFiles := map[string]string{"src/Kernel.php": appSrc,
		"ghome/config.json": "{\n    \"config\": {\"sort-packages\": true}\n}\n",
		"ghome/auth.json":   "{\n    \"http-basic\": {}\n}\n",
		"auth.json":         "{\n    \"bearer\": {}\n}\n",
	}
	globalHome := `{
    "repositories": [{"type": "composer", "url": "http://127.0.0.1:8899"}, {"packagist.org": false}],
    "config": {"secure-http": false},
    "require": {"acme/util": "^1.0"}
}
`
	binCompat := withConfig(versioned(`{"acme/cli": "^1.0"}`, ""), `"bin-compat": "symlink"`)
	ev := `{"plugfix/events-plugin": "1.0.0", "acme/log": "^1.0"}`
	evProject := withConfig(project(ev, `"scripts": {"command": "echo command-script-listener"}`), `"allow-plugins": {"plugfix/*": true}`)
	classFiles := map[string]string{"src/Kernel.php": appSrc, "src/HelloCommand.php": helloCommandPHP, "src/FailCommand.php": failCommandPHP,
		"src/BoomCommand.php": boomCommandPHP, "src/NamedCommand.php": namedCommandPHP, "src/NotCommand.php": "<?php\nnamespace App;\n\nclass NotCommand {}\n"}
	classJSON := project(`{"acme/log": "^1.0"}`, `"scripts": {
        "hello": "App\\HelloCommand",
        "fail": "App\\FailCommand",
        "boom": "App\\BoomCommand",
        "named": "App\\NamedCommand",
        "missing": "App\\MissingCommand",
        "notcmd": "App\\NotCommand",
        "both": ["@hello", "echo after-hello"],
        "post-autoload-dump": "App\\HelloCommand"
    },
    "scripts-descriptions": {"fail": "Fails with code 3"}`)
	pkgst := []string{"packagiststub"}
	sslEnv := map[string]string{"SSL_CERT_FILE": "/fixtures/git/_api/ca.pem"}
	fundJSON := `{
    "name": "test/fund",
    "require": {"stubvendor/alpha": "^1.0", "stubvendor/beta": "^2.1", "othervendor/gamma": "^0.3"}
}
`
	fundInstalled := `{
    "packages": [
        {"name": "stubvendor/alpha", "version": "1.0.0", "version_normalized": "1.0.0.0", "type": "library", "funding": [{"type": "github", "url": "https://github.com/sponsors/alpha-installed"}]},
        {"name": "stubvendor/beta", "version": "2.1.0", "version_normalized": "2.1.0.0", "type": "library"},
        {"name": "othervendor/gamma", "version": "0.3.0", "version_normalized": "0.3.0.0", "type": "library", "funding": [{"type": "custom", "url": "https://gamma.example.org/fund"}]}
    ],
    "dev": true,
    "dev-package-names": []
}
`
	// the packages must be present in vendor/ (Factory::purgePackages drops
	// the others from the local repository)
	fundFiles := map[string]string{"vendor/composer/installed.json": fundInstalled, "vendor/stubvendor/alpha/README": "alpha\n",
		"vendor/stubvendor/beta/README": "beta\n", "vendor/othervendor/gamma/README": "gamma\n"}
	vvv := func(name, cmd string) Scenario {
		return Scenario{Name: "cli3-vvv-" + name, JSON: vvvJSON, Files: files, Setup: upd, Command: cmd}
	}
	return []Scenario{
		// -vvv output of Factory/Config (Reading, Loading config file, CA
		// bundle checks, global composer) and of the commands' own reads
		vvv("validate", "validate -vvv"),
		vvv("show", "show -vvv"),
		vvv("config-list", "config --list -vvv"),
		vvv("config-get", "config name -vvv"),
		vvv("config-set", "config name foo/bar -vvv"),
		vvv("config-extra", "config extra.x -vvv"),
		vvv("config-global", "config --global --list -vvv"),
		vvv("dump-autoload", "dump-autoload -vvv"),
		vvv("run-script", "run-script greet -vvv"),
		vvv("script-command", "greet -vvv"),
		vvv("licenses", "licenses -vvv"),
		vvv("depends", "depends acme/log -vvv"),
		vvv("suggests", "suggests -vvv"),
		vvv("no-plugins", "show --no-plugins -vvv"),
		{Name: "cli3-vvv-shell-verbosity", JSON: vvvJSON, Files: files, Setup: upd, Env: map[string]string{"SHELL_VERBOSITY": "3"}, Command: "show"},
		{Name: "cli3-vvv-no-lock", JSON: vvvJSON, Files: files, Command: "show -vvv"},
		{Name: "cli3-vvv-auth-files", JSON: vvvJSON, Files: homeFiles, Env: map[string]string{"COMPOSER_HOME": "/work/project/ghome", "COMPOSER_AUTH": `{"bearer": {"example.org": "t"}}`}, Setup: upd, Command: "show -vvv"},
		{Name: "cli3-vvv-cafile", JSON: withConfig(vvvJSON, `"cafile": "/etc/ssl/certs/ca-certificates.crt"`), Files: files, Setup: upd, Command: "show -vvv"},
		{Name: "cli3-vvv-ssl-cert-file", JSON: vvvJSON, Files: files, Env: map[string]string{"SSL_CERT_FILE": "/etc/ssl/certs/ca-certificates.crt"}, Setup: upd, Command: "show -vvv"},
		{Name: "cli3-vvv-ssl-cert-missing", JSON: vvvJSON, Files: files, Env: map[string]string{"SSL_CERT_FILE": "/nonexistent/ca.pem", "SSL_CERT_DIR": "/nonexistent-dir"}, Setup: upd, Command: "show -vvv"},
		{Name: "cli3-vvv-global-home", JSON: vvvJSON, Files: map[string]string{"src/Kernel.php": appSrc, "ghome/composer.json": globalHome}, Env: ghome, Setup: []string{"update --no-progress", "global update --no-progress"}, Command: "show -vvv"},
		{Name: "cli3-vvv-allow-plugins", JSON: withConfig(vvvJSON, `"allow-plugins": {"plugfix/*": true}`), Files: files, Setup: upd, Command: "show -vvv"},
		{Name: "cli3-vvv-validate-file", JSON: vvvJSON, Files: files, Setup: upd, Command: "validate composer.json -vvv"},
		// bin-compat "symlink" deprecation (Config::get, ErrorHandler)
		{Name: "cli3-bin-compat-get", JSON: binCompat, Files: files, Command: "config bin-compat"},
		{Name: "cli3-bin-compat-list", JSON: binCompat, Files: files, Command: "config --list"},
		{Name: "cli3-bin-compat-set", JSON: binCompat, Files: files, Command: "config bin-compat proxy"},
		{Name: "cli3-bin-compat-install", JSON: binCompat, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true},
		{Name: "cli3-bin-compat-env", JSON: versioned(`{"acme/log": "^1.0"}`, ""), Files: files, Env: map[string]string{"COMPOSER_BIN_COMPAT": "symlink"}, Command: "config bin-compat"},
		// PluginEvents::COMMAND of validate and reinstall
		{Name: "cli3-event-validate", JSON: evProject, Files: files, Setup: upd, Command: "validate"},
		{Name: "cli3-event-validate-no-plugins", JSON: evProject, Files: files, Setup: upd, Command: "validate --no-plugins"},
		{Name: "cli3-event-reinstall", JSON: evProject, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress"},
		{Name: "cli3-event-reinstall-none", JSON: evProject, Files: files, Setup: upd, Command: "reinstall nope/nope"},
		// scripts pointing to Symfony Console command classes
		{Name: "cli3-cmdclass-run", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script hello"},
		{Name: "cli3-cmdclass-run-args", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script hello -- bob --shout"},
		{Name: "cli3-cmdclass-run-verbose", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script hello -v"},
		{Name: "cli3-cmdclass-run-fail", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script fail"},
		{Name: "cli3-cmdclass-run-boom", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script boom"},
		{Name: "cli3-cmdclass-run-missing", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script missing"},
		{Name: "cli3-cmdclass-run-notcmd", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script notcmd"},
		{Name: "cli3-cmdclass-run-named", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script named"},
		{Name: "cli3-cmdclass-run-nested", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script both"},
		{Name: "cli3-cmdclass-reserved", JSON: classJSON, Files: classFiles, Setup: upd, Command: "dump-autoload"},
		{Name: "cli3-cmdclass-list", JSON: classJSON, Files: classFiles, Setup: upd, Command: "run-script --list"},
		{Name: "cli3-cmdclass-direct", JSON: classJSON, Files: classFiles, Setup: upd, Command: "hello alice"},
		{Name: "cli3-cmdclass-direct-options", JSON: classJSON, Files: classFiles, Setup: upd, Command: "hello --shout -v"},
		{Name: "cli3-cmdclass-direct-bad-option", JSON: classJSON, Files: classFiles, Setup: upd, Command: "hello --nope"},
		{Name: "cli3-cmdclass-direct-fail", JSON: classJSON, Files: classFiles, Setup: upd, Command: "fail"},
		{Name: "cli3-cmdclass-direct-boom", JSON: classJSON, Files: classFiles, Setup: upd, Command: "boom"},
		{Name: "cli3-cmdclass-direct-named", JSON: classJSON, Files: classFiles, Setup: upd, Command: "named"},
		{Name: "cli3-cmdclass-not-installed", JSON: classJSON, Files: classFiles, Command: "run-script hello"},
		// packagist.org (loopback stub) without composer.json
		{Name: "cli3-packagist-search", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search alpha"},
		{Name: "cli3-packagist-search-json", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search gamma --format=json"},
		{Name: "cli3-packagist-search-all", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search -vv alpha"},
		{Name: "cli3-packagist-search-name", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search --only-name alp"},
		{Name: "cli3-packagist-search-vendor", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search --only-vendor stub"},
		{Name: "cli3-packagist-search-none", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "search nothing-matches-this"},
		{Name: "cli3-packagist-browse-show", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "browse stubvendor/alpha --show"},
		{Name: "cli3-packagist-browse-homepage", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "browse -H stubvendor/beta -s"},
		{Name: "cli3-packagist-browse-source", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "browse stubvendor/beta othervendor/gamma -s"},
		{Name: "cli3-packagist-browse-missing", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "browse nope/nope -s"},
		{Name: "cli3-packagist-browse-root", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "browse"},
		{Name: "cli3-packagist-fund", JSON: fundJSON, Files: fundFiles, GitFixtures: pkgst, Env: sslEnv, Command: "fund"},
		{Name: "cli3-packagist-fund-json", JSON: fundJSON, Files: fundFiles, GitFixtures: pkgst, Env: sslEnv, Command: "fund --format=json"},
		{Name: "cli3-packagist-fund-no-json", NoJSON: true, GitFixtures: pkgst, Env: sslEnv, Command: "fund"},
	}
}

// apiHostArgs maps the further host names a stub fixture lists in
// _api/hosts (e.g. repo.packagist.org, packagist.org of packagiststub.sh)
// to the loopback stub.
func apiHostArgs(gitfx string) []string {
	b, err := os.ReadFile(filepath.Join(gitfx, "_api", "hosts"))
	if err != nil {
		return nil
	}
	var out []string
	for _, h := range strings.Fields(string(b)) {
		out = append(out, "--add-host="+h+":127.0.0.1")
	}
	return out
}
