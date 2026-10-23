package main

import "strings"

// vcsScenarios compare VcsRepository/GitDriver, the GitDownloader and source
// installs (Composer f83588392df43d4fd027d84d26de459946ef4e80). The git
// fixture repositories are built by testdata/compare/vcs/*.sh.
func init() { extraScenarios = append(extraScenarios, vcsScenarios) }

// vcsProject is project() with extra repositories placed before the fixture
// composer repository.
func vcsProject(requireJSON string, repos []string, extra string) string {
	s := project(requireJSON, extra)
	if len(repos) > 0 {
		s = strings.Replace(s, `"repositories": [
`, `"repositories": [
        `+strings.Join(repos, ",\n        ")+",\n", 1)
	}
	return s
}

func gitRepo(typ, url string) string {
	return `{"type": "` + typ + `", "url": "` + url + `"}`
}

// deterministic commits of testdata/compare/vcs/gitlib.sh
const (
	refLib100  = "b2b1b720ad1e4e06cb3d641ea573be406bab1502"
	refLib110  = "702184333938c0ff49e87f0af7d9931d71486a52"
	refLibMain = "ba0333212d6fb4b08718b30b9dc442e929393f8b"
)

var gitIdent = map[string]string{"GIT_AUTHOR_NAME": "Tester", "GIT_AUTHOR_EMAIL": "tester@example.org", "GIT_COMMITTER_NAME": "Tester", "GIT_COMMITTER_EMAIL": "tester@example.org"}

// vcsOnlyProject has only the given repositories (no fixture composer
// repository, whose root packages.json would otherwise be revalidated).
func vcsOnlyProject(requireJSON string, repos []string) string {
	return "{\n    \"name\": \"test/project\",\n    \"repositories\": [\n        " + strings.Join(repos, ",\n        ") + ",\n        {\"packagist.org\": false}\n    ],\n    \"config\": {\"secure-http\": false},\n    \"require\": " + requireJSON + "\n}\n"
}

func pkgRepo(defs ...string) string {
	return `{"type": "package", "package": [` + strings.Join(defs, ", ") + `]}`
}

func vcsScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	lib := gitRepo("vcs", "/fixtures/git/gitlib.git")
	dep := gitRepo("vcs", "/fixtures/git/gitdep.git")
	libFx := []string{"gitlib"}
	both2 := []string{"gitlib", "gitdep"}
	misc := gitRepo("vcs", "/fixtures/git/gitmisc.git")
	miscFx := []string{"gitmisc"}
	work := gitRepo("vcs", "/fixtures/git/_work")
	httpLib := gitRepo("vcs", "http://127.0.0.1:8899/git/gitlib.git")
	cacheEnv := map[string]string{"COMPOSER_CACHE_DIR": "/work/out/cache"}
	modified := map[string]string{"vendor/vcs/gitlib/src/Lib.php": "<?php\nnamespace VcsLib;\n\nclass Lib { const VERSION = \"local\"; }\n"}
	devCommit := "repo gitlib\nput src/Extra.php '<?php\nnamespace VcsLib;\n\nclass Extra {}'\ncommit \"Add Extra\"\npublish gitlib\n"
	both := pkgRepo(`{"name": "vcs/both", "version": "1.0.0", "source": {"type": "git", "url": "/fixtures/git/gitlib.git", "reference": "`+refLib100+`"}, "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/acme-log-1.0.0.zip"}, "autoload": {"psr-4": {"VcsLib\\": "src/"}}}`,
		`{"name": "vcs/both", "version": "1.1.0", "source": {"type": "git", "url": "/fixtures/git/gitlib.git", "reference": "`+refLib110+`"}, "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip"}, "autoload": {"psr-4": {"VcsLib\\": "src/"}}}`,
		`{"name": "vcs/both", "version": "dev-main", "source": {"type": "git", "url": "/fixtures/git/gitlib.git", "reference": "`+refLibMain+`"}, "autoload": {"psr-4": {"VcsLib\\": "src/"}}}`)
	swap := pkgRepo(`{"name": "vcs/swap", "version": "1.0.0", "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/acme-log-1.0.0.zip"}}`,
		`{"name": "vcs/swap", "version": "dev-main", "source": {"type": "git", "url": "/fixtures/git/gitlib.git", "reference": "`+refLibMain+`"}}`)
	fallback := pkgRepo(`{"name": "vcs/fallback", "version": "1.0.0", "source": {"type": "git", "url": "http://127.0.0.1:8899/git/nope.git", "reference": "` + refLib100 + `"}, "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/acme-log-1.0.0.zip"}}`)
	lib100 := []string{"update --no-progress --with vcs/gitlib:1.0.0"}
	return []Scenario{
		{Name: "vcs-verbose", JSON: vcsProject(`{"vcs/gitdep": "^1.0"}`, []string{lib, dep}, ""), Files: files, GitFixtures: both2, Command: "update --no-progress -v"},
		{Name: "vcs-debug", JSON: vcsProject(`{"vcs/gitdep": "^1.0"}`, []string{lib, dep}, ""), Files: files, GitFixtures: both2, Command: "update --no-progress -vv"},
		{Name: "vcs-only-branches", JSON: vcsProject(`{"vcs/gitlib": "1.x-dev"}`, []string{`{"type": "vcs", "url": "/fixtures/git/gitlib.git", "only-branches": ["1.*", "nope"]}`}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress -vv"},
		{Name: "vcs-only-branches-invalid", JSON: vcsProject(`{"vcs/gitlib": "1.x-dev"}`, []string{`{"type": "vcs", "url": "/fixtures/git/gitlib.git", "only-branches": "1.x"}`}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-misc-vv", JSON: vcsProject(`{"vcs/gitmisc": "^1.0"}`, []string{misc}, ""), Files: files, GitFixtures: miscFx, Command: "update --no-progress -vv", Probe: true},
		{Name: "vcs-misc-numeric-branch", JSON: vcsProject(`{"vcs/gitmisc": "2.0.x-dev"}`, []string{misc}, ""), Files: files, GitFixtures: miscFx, Command: "update --no-progress"},
		{Name: "vcs-misc-slash-branch", JSON: vcsProject(`{"vcs/gitmisc": "dev-feature/x"}`, []string{misc}, ""), Files: files, GitFixtures: miscFx, Command: "update --no-progress"},
		{Name: "vcs-misc-old-tag", JSON: vcsProject(`{"vcs/gitmisc": "0.9.0"}`, []string{misc}, ""), Files: files, GitFixtures: miscFx, Command: "update --no-progress"},
		{Name: "vcs-ref-pin", JSON: vcsProject(`{"vcs/gitlib": "dev-main#`+refLib100+`"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-not-a-repo", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{work}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-git-type-not-a-repo", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("git", "/fixtures/git/_work")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-type-hg", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("hg", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-type-svn", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("svn", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-type-fossil", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("fossil", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-type-perforce", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("perforce", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-type-pear", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("pear", "/fixtures/git/gitlib.git")}, ""), Files: files, Command: "update --no-progress"},
		{Name: "vcs-type-forgejo", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("forgejo", "/fixtures/git/gitlib.git")}, ""), Files: files, Command: "update --no-progress"},
		{Name: "vcs-type-github-local", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("github", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-insecure-http", JSON: strings.Replace(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{httpLib}, ""), `"config": {"secure-http": false}`, `"config": {"secure-http": true}`, 1), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-http-missing", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("vcs", "http://127.0.0.1:8899/git/nope.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-http-cached", JSON: vcsOnlyProject(`{"vcs/gitlib": "^1.0"}`, []string{httpLib}), Files: files, GitFixtures: libFx, Env: cacheEnv, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true},
		{Name: "vcs-http-cached-debug", JSON: vcsOnlyProject(`{"vcs/gitlib": "^1.0"}`, []string{httpLib}), Files: files, GitFixtures: libFx, Env: cacheEnv, Setup: []string{"update --no-progress --no-install"}, Command: "update --no-progress -vv"},
		{Name: "vcs-http-dev-update", JSON: vcsProject(`{"vcs/gitlib": "dev-main"}`, []string{httpLib}, ""), Files: files, GitFixtures: libFx, Env: cacheEnv, Setup: []string{"update --no-progress"}, GitAfter: devCommit, Command: "update --no-progress"},
		// --- source installs / downloader ---
		{Name: "source-update-checkout", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: lib100, Command: "update --no-progress", Probe: true},
		{Name: "source-update-checkout-v", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: lib100, Command: "update --no-progress -v"},
		{Name: "source-downgrade-v", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, Command: "update --no-progress --with vcs/gitlib:1.0.0 -v"},
		{Name: "source-dev-update", JSON: vcsProject(`{"vcs/gitlib": "dev-main"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, GitAfter: devCommit, Command: "update --no-progress", Probe: true},
		{Name: "source-local-changes", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: lib100, After: modified, Command: "update --no-progress"},
		{Name: "source-discard-true", JSON: withConfig(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), `"discard-changes": true`), Files: files, GitFixtures: libFx, Setup: lib100, After: modified, Command: "update --no-progress", Probe: true},
		{Name: "source-discard-stash", JSON: withConfig(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), `"discard-changes": "stash"`), Files: files, GitFixtures: libFx, Setup: lib100, After: modified, Command: "update --no-progress"},
		{Name: "source-discard-stash-ident", JSON: withConfig(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), `"discard-changes": "stash"`), Files: files, GitFixtures: libFx, Env: gitIdent, Setup: lib100, After: modified, Command: "update --no-progress"},
		{Name: "source-uninstall-changes", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, After: modified, Command: "remove vcs/gitlib --no-progress"},
		{Name: "source-uninstall-discard", JSON: withConfig(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), `"discard-changes": true`), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, After: modified, Command: "remove vcs/gitlib --no-progress"},
		{Name: "source-untracked-only", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: lib100, After: map[string]string{"vendor/vcs/gitlib/NOTES.txt": "local notes\n"}, Command: "update --no-progress"},
		{Name: "source-prefer-source", JSON: vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress --prefer-source", Probe: true},
		{Name: "source-prefer-dist-default", JSON: vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "source-prefer-install-flag", JSON: vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress --prefer-install=source"},
		{Name: "source-preferred-install-source", JSON: withConfig(vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), `"preferred-install": "source"`), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "source-preferred-install-pattern", JSON: withConfig(vcsProject(`{"vcs/both": "1.0.0", "acme/log": "^1.0"}`, []string{both}, ""), `"preferred-install": {"vcs/*": "source", "*": "dist"}`), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "source-preferred-install-auto-dev", JSON: withConfig(vcsProject(`{"vcs/both": "dev-main"}`, []string{both}, ""), `"preferred-install": {"*": "auto"}`), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "source-dev-prefer-dist", JSON: vcsProject(`{"vcs/gitlib": "dev-main"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress --prefer-dist"},
		{Name: "source-dist-to-source", JSON: vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": vcsProject(`{"vcs/both": "dev-main"}`, []string{both}, "")}, Command: "update --no-progress", Probe: true},
		{Name: "source-source-to-dist", JSON: vcsProject(`{"vcs/swap": "dev-main"}`, []string{swap}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": vcsProject(`{"vcs/swap": "1.0.0"}`, []string{swap}, "")}, Command: "update --no-progress"},
		{Name: "source-source-to-dist-changes", JSON: vcsProject(`{"vcs/swap": "dev-main"}`, []string{swap}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, After: map[string]string{"composer.json": vcsProject(`{"vcs/swap": "1.0.0"}`, []string{swap}, ""), "vendor/vcs/swap/src/Lib.php": "<?php // local\n"}, Command: "update --no-progress"},
		{Name: "source-fallback-dist", JSON: vcsProject(`{"vcs/fallback": "1.0.0"}`, []string{fallback}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress --prefer-source"},
		{Name: "source-install-lock-prefer-source", JSON: vcsProject(`{"vcs/both": "1.0.0"}`, []string{both}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress --prefer-source"},
		{Name: "source-reinstall-same", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress"}, Command: "install --no-progress"},
		{Name: "vcs-git-install", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-git-type", JSON: vcsProject(`{"vcs/gitlib": "1.0.0"}`, []string{gitRepo("git", "/fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-file-url", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("vcs", "file:///fixtures/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-dev-main", JSON: vcsProject(`{"vcs/gitlib": "2.x-dev"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-dev-branch", JSON: vcsProject(`{"vcs/gitlib": "1.x-dev"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
		{Name: "vcs-dep", JSON: vcsProject(`{"vcs/gitdep": "^1.0"}`, []string{lib, dep}, ""), Files: files, GitFixtures: both2, Command: "update --no-progress", Probe: true},
		{Name: "vcs-install-from-lock", JSON: vcsProject(`{"vcs/gitdep": "^1.0"}`, []string{lib, dep}, ""), Files: files, GitFixtures: both2, Setup: []string{"update --no-progress --no-install"}, Command: "install --no-progress", Probe: true},
		{Name: "vcs-dry-run", JSON: vcsProject(`{"vcs/gitdep": "^1.0"}`, []string{lib, dep}, ""), Files: files, GitFixtures: both2, Command: "update --dry-run --no-progress"},
		{Name: "vcs-no-composer-json", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("vcs", "/fixtures/git/gitempty.git")}, ""), Files: files, GitFixtures: []string{"gitempty"}, Command: "update --no-progress"},
		{Name: "vcs-missing-path", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("vcs", "/fixtures/git/nope.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress"},
		{Name: "vcs-odd-refs", JSON: vcsProject(`{"vcs/gitodd": "*"}`, []string{gitRepo("vcs", "/fixtures/git/gitodd.git")}, ""), Files: files, GitFixtures: []string{"gitodd"}, Command: "update --no-progress"},
		{Name: "vcs-odd-refs-vv", JSON: vcsProject(`{"vcs/gitodd": "*"}`, []string{gitRepo("vcs", "/fixtures/git/gitodd.git")}, ""), Files: files, GitFixtures: []string{"gitodd"}, Command: "update --no-progress -vv"},
		{Name: "vcs-http-mirror", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{gitRepo("vcs", "http://127.0.0.1:8899/git/gitlib.git")}, ""), Files: files, GitFixtures: libFx, Command: "update --no-progress", Probe: true},
	}
}
