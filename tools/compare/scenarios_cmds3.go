package main

// cmds3Scenarios covers the reinstall, status and create-project commands
// (ReinstallCommand, StatusCommand, CreateProjectCommand of Composer
// f83588392df43d4fd027d84d26de459946ef4e80).
func init() { extraScenarios = append(extraScenarios, cmds3Scenarios) }

func cmds3Scenarios() []Scenario {
	var out []Scenario
	out = append(out, reinstallScenarios()...)
	out = append(out, statusScenarios()...)
	out = append(out, createProjectScenarios()...)
	return out
}

func reinstallScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	basic := project(`{"acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0"}`, `"require-dev": {"acme/devtools": "^1.0"}`)
	upd := []string{"update --no-progress"}
	scripts := project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, `"scripts": {"pre-install-cmd": "echo PRE-INSTALL $COMPOSER_DEV_MODE", "post-install-cmd": "echo POST-INSTALL", "pre-package-uninstall": "echo PRE-UNINSTALL", "post-package-install": "echo POST-PKG-INSTALL", "pre-autoload-dump": "echo PRE-DUMP", "post-autoload-dump": "echo POST-DUMP"}`)
	lib := gitRepo("vcs", "/fixtures/git/gitlib.git")
	libFx := []string{"gitlib"}
	both := pkgRepo(`{"name": "vcs/both", "version": "1.0.0", "source": {"type": "git", "url": "/fixtures/git/gitlib.git", "reference": "` + refLib100 + `"}, "dist": {"type": "zip", "url": "http://127.0.0.1:8899/dists/acme-log-1.0.0.zip"}, "autoload": {"psr-4": {"VcsLib\\": "src/"}}}`)
	bothJSON := vcsProject(`{"vcs/both": "1.0.0", "acme/util": "^1.0"}`, []string{both}, "")
	modified := map[string]string{"vendor/vcs/gitlib/src/Lib.php": "<?php\nnamespace VcsLib;\n\nclass Lib { const VERSION = \"local\"; }\n"}
	return []Scenario{
		{Name: "reinstall-basic", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress", Probe: true},
		{Name: "reinstall-progress", RacyRemovalBar: true, JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/log acme/http"},
		{Name: "reinstall-verbose", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress -v"},
		{Name: "reinstall-wildcard", JSON: basic, Files: files, Setup: upd, Command: "reinstall 'acme/l*' 'ACME/HTTP' --no-progress", Probe: true},
		{Name: "reinstall-all-wildcard", JSON: basic, Files: files, Setup: upd, Command: "reinstall '*' --no-progress", Probe: true},
		{Name: "reinstall-type", JSON: basic, Files: files, Setup: upd, Command: "reinstall --type=library --type=metapackage --no-progress", Probe: true},
		{Name: "reinstall-type-none", JSON: basic, Files: files, Setup: upd, Command: "reinstall --type=composer-plugin --no-progress"},
		{Name: "reinstall-type-and-names", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/log --type=library"},
		{Name: "reinstall-no-args", JSON: basic, Files: files, Setup: upd, Command: "reinstall"},
		{Name: "reinstall-no-match", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/nope 'other/*' --no-progress"},
		{Name: "reinstall-partial-match", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/nope acme/cli --no-progress"},
		{Name: "reinstall-not-installed", JSON: basic, Files: files, Command: "reinstall acme/log --no-progress"},
		{Name: "reinstall-lock-only", JSON: basic, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "reinstall acme/log --no-progress"},
		{Name: "reinstall-no-autoloader", JSON: basic, Files: files, Setup: upd, After: map[string]string{"vendor/composer/autoload_psr4.php": "<?php // stale\n"}, Command: "reinstall acme/log --no-progress --no-autoloader"},
		{Name: "reinstall-optimize", JSON: basic, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress -o --classmap-authoritative --apcu-autoloader-prefix=pfx", Probe: true},
		{Name: "reinstall-restores-files", JSON: basic, Files: files, Setup: upd, After: map[string]string{"vendor/acme/log/src/Logger.php": "<?php // broken\n", "vendor/acme/log/EXTRA.txt": "extra\n"}, Command: "reinstall acme/log --no-progress", Probe: true},
		{Name: "reinstall-no-dev-mode", JSON: basic, Files: files, Setup: []string{"update --no-progress --no-dev"}, Command: "reinstall acme/log --no-progress", Probe: true},
		{Name: "reinstall-scripts", JSON: scripts, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress"},
		{Name: "reinstall-no-scripts", JSON: scripts, Files: files, Setup: upd, Command: "reinstall acme/log --no-progress --no-scripts"},
		{Name: "reinstall-no-json", NoJSON: true, Command: "reinstall acme/log"},
		{Name: "reinstall-prefer-source", JSON: bothJSON, Files: files, GitFixtures: libFx, Setup: upd, Command: "reinstall vcs/both --no-progress --prefer-source", Probe: true},
		{Name: "reinstall-prefer-install-source", JSON: bothJSON, Files: files, GitFixtures: libFx, Setup: upd, Command: "reinstall vcs/both --no-progress --prefer-install=source"},
		{Name: "reinstall-source-to-dist", JSON: bothJSON, Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress --prefer-source"}, Command: "reinstall vcs/both --no-progress --prefer-dist", Probe: true},
		{Name: "reinstall-source-default", JSON: bothJSON, Files: files, GitFixtures: libFx, Setup: []string{"update --no-progress --prefer-source"}, Command: "reinstall vcs/both --no-progress"},
		{Name: "reinstall-prefer-conflict", JSON: bothJSON, Files: files, GitFixtures: libFx, Setup: upd, Command: "reinstall vcs/both --prefer-source --prefer-install=dist"},
		{Name: "reinstall-git", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: upd, Command: "reinstall vcs/gitlib --no-progress", Probe: true},
		{Name: "reinstall-git-changes", JSON: vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), Files: files, GitFixtures: libFx, Setup: upd, After: modified, Command: "reinstall vcs/gitlib --no-progress"},
		{Name: "reinstall-git-discard", JSON: withConfig(vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, ""), `"discard-changes": true`), Files: files, GitFixtures: libFx, Setup: upd, After: modified, Command: "reinstall vcs/gitlib --no-progress", Probe: true},
	}
}

func statusScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	upd := []string{"update --no-progress"}
	dist := project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, "")
	distChanged := map[string]string{"vendor/acme/log/src/Logger.php": "<?php // changed\n", "vendor/acme/log/src/New.php": "<?php // new\n"}
	lib := gitRepo("vcs", "/fixtures/git/gitlib.git")
	libFx := []string{"gitlib"}
	git := func(extra string) string {
		return vcsProject(`{"vcs/gitlib": "^1.0", "acme/util": "^1.0"}`, []string{lib}, extra)
	}
	gitScripts := git(`"scripts": {"mk-commit": "cd vendor/vcs/gitlib && echo '<?php // local' > src/Local.php && git add -A && git commit -qm local", "mk-checkout": "cd vendor/vcs/gitlib && git checkout -q v1.0.0", "mk-branch": "cd vendor/vcs/gitlib && git checkout -q -b topic && echo x > NOTE.txt && git add -A && git commit -qm topic"}`)
	modified := map[string]string{"vendor/vcs/gitlib/src/Lib.php": "<?php\nnamespace VcsLib;\n\nclass Lib { const VERSION = \"local\"; }\n"}
	pathFiles := map[string]string{"src/Kernel.php": appSrc, "pkgs/local/composer.json": `{"name": "path/local", "version": "1.0.0", "autoload": {"psr-4": {"PathLocal\\": "src/"}}}`, "pkgs/local/src/Local.php": "<?php\nnamespace PathLocal;\n\nclass Local {}\n"}
	pathRepo := func(opts string) string {
		return vcsProject(`{"path/local": "*", "acme/util": "^1.0"}`, []string{`{"type": "path", "url": "pkgs/*"` + opts + `}`}, "")
	}
	scripts := project(`{"acme/util": "^1.0"}`, `"scripts": {"pre-status-cmd": "echo PRE-STATUS $COMPOSER_DEV_MODE", "post-status-cmd": "echo POST-STATUS"}`)
	return []Scenario{
		{Name: "status-dist-clean", JSON: dist, Files: files, Setup: upd, Command: "status"},
		{Name: "status-dist-clean-v", JSON: dist, Files: files, Setup: upd, Command: "status -v"},
		{Name: "status-dist-changed", JSON: dist, Files: files, Setup: upd, After: distChanged, Command: "status"},
		{Name: "status-dist-changed-v", JSON: dist, Files: files, Setup: upd, After: distChanged, Command: "status -v"},
		{Name: "status-dist-changed-vv", JSON: dist, Files: files, Setup: upd, After: distChanged, Command: "status -vv"},
		{Name: "status-dist-removed-v", JSON: dist, Files: files, Setup: upd, After: map[string]string{"vendor/acme/util/composer.json": ""}, Command: "status -v"},
		{Name: "status-dist-two-packages", JSON: dist, Files: files, Setup: upd, After: map[string]string{"vendor/acme/log/src/Logger.php": "<?php // changed\n", "vendor/acme/util/composer.json": "{}\n"}, Command: "status"},
		{Name: "status-not-installed", JSON: dist, Files: files, Command: "status"},
		{Name: "status-lock-only", JSON: dist, Files: files, Setup: []string{"update --no-progress --no-install"}, Command: "status"},
		{Name: "status-no-json", NoJSON: true, Command: "status"},
		{Name: "status-scripts", JSON: scripts, Files: files, Setup: upd, Command: "status"},
		{Name: "status-no-scripts", JSON: scripts, Files: files, Setup: upd, Command: "status --no-scripts"},
		{Name: "status-git-clean", JSON: git(""), Files: files, GitFixtures: libFx, Setup: upd, Command: "status -v"},
		{Name: "status-git-modified", JSON: git(""), Files: files, GitFixtures: libFx, Setup: upd, After: modified, Command: "status"},
		{Name: "status-git-modified-v", JSON: git(""), Files: files, GitFixtures: libFx, Setup: upd, After: modified, Command: "status -v"},
		{Name: "status-git-untracked", JSON: git(""), Files: files, GitFixtures: libFx, Setup: upd, After: map[string]string{"vendor/vcs/gitlib/NOTES.txt": "notes\n"}, Command: "status -v"},
		{Name: "status-git-unpushed", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-commit"}, Command: "status"},
		{Name: "status-git-unpushed-v", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-commit"}, Command: "status -v"},
		{Name: "status-git-unpushed-vv", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-commit"}, Command: "status -vv"},
		{Name: "status-git-version-v", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-checkout"}, Command: "status -v"},
		{Name: "status-git-version-vv", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-checkout"}, Command: "status -vv"},
		{Name: "status-git-branch-v", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-branch"}, Command: "status -v"},
		{Name: "status-mixed", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-commit"}, After: map[string]string{"vendor/acme/util/composer.json": "{}\n", "vendor/vcs/gitlib/src/Lib.php": "<?php // changed\n"}, Command: "status"},
		{Name: "status-mixed-v", JSON: gitScripts, Files: files, GitFixtures: libFx, Env: gitIdent, Setup: []string{"update --no-progress", "run-script mk-commit"}, After: map[string]string{"vendor/acme/util/composer.json": "{}\n", "vendor/vcs/gitlib/src/Lib.php": "<?php // changed\n"}, Command: "status -v"},
		{Name: "status-path-symlink", JSON: pathRepo(""), Files: pathFiles, Setup: upd, Command: "status -v"},
		{Name: "status-path-mirror", JSON: pathRepo(`, "options": {"symlink": false}`), Files: pathFiles, Setup: upd, Command: "status -v"},
		{Name: "status-path-mirror-changed", JSON: pathRepo(`, "options": {"symlink": false}`), Files: pathFiles, Setup: upd, After: map[string]string{"vendor/path/local/src/Local.php": "<?php // changed\n"}, Command: "status -v"},
	}
}

func createProjectScenarios() []Scenario {
	const r = " --repository=http://127.0.0.1:8899 --no-secure-http --no-progress"
	sc := func(name, cmd string) Scenario {
		return Scenario{Name: "create-project-" + name, NoJSON: true, Command: "create-project " + cmd,
			Snapshot: []string{"app/composer.json", "app/composer.lock", "app/vendor/composer/installed.json", "app/vendor/composer/installed.php", "lib/composer.json"}}
	}
	withFiles := func(s Scenario, files map[string]string) Scenario { s.Files = files; return s }
	withGit := func(s Scenario) Scenario { s.GitFixtures = []string{"gitlib"}; return s }
	vcsRepo := ` --repository='{"type": "vcs", "url": "/fixtures/git/gitlib.git"}' --no-progress`
	pathFiles := map[string]string{
		"pkgs/app/composer.json": `{"name": "local/app", "version": "1.2.0", "type": "project", "repositories": [{"packagist.org": false}], "require": {"php": ">=7.4"}, "autoload": {"psr-4": {"LocalApp\\": "src/"}}}`,
		"pkgs/app/src/App.php":   "<?php\nnamespace LocalApp;\n\nclass App {}\n",
	}
	list := []Scenario{
		sc("basic", "cproj/app app"+r),
		sc("default-dir", "cproj/app"+r),
		sc("version-arg", "cproj/app app 1.0.0"+r),
		sc("version-colon", "cproj/app:^1.0 app"+r),
		sc("version-equals", "cproj/app=1.0.0 app"+r),
		sc("stability-flag", "cproj/app app -s beta"+r),
		sc("stability-at", "'cproj/app:^2.0@beta' app"+r),
		sc("version-beta", "cproj/app app 2.0.0-beta1"+r),
		sc("stability-invalid", "cproj/app app -s foo"+r),
		sc("not-found", "cproj/nope app"+r),
		sc("not-found-version", "cproj/app app 9.0.0"+r),
		sc("php-too-low", "cproj/needsphp app"+r),
		sc("php-ignore", "cproj/needsphp app --ignore-platform-reqs"+r),
		sc("secure-http", "cproj/app app --repository=http://127.0.0.1:8899 --no-progress"),
		sc("repository-url", "cproj/app app --repository-url=http://127.0.0.1:8899 --no-secure-http --no-progress"),
		sc("packagist-false-repo", `cproj/app app --repository='{"packagist.org": false}'`+r),
		sc("no-install", "cproj/app app --no-install"+r),
		sc("no-scripts", "cproj/app app --no-scripts"+r),
		sc("no-dev", "cproj/app app --no-dev"+r),
		sc("dev-option", "cproj/app app --dev"+r),
		sc("no-custom-installers", "cproj/app app --no-custom-installers"+r),
		sc("no-plugins", "cproj/app app --no-plugins"+r),
		sc("no-audit", "cproj/app app --no-audit"+r),
		sc("add-repository", "cproj/app app --add-repository"+r),
		sc("require", "cproj/app app --require acme/util:^1.0 --require acme/cli"+r),
		sc("ask-noninteractive", "cproj/app app --ask"+r),
		sc("prefer-conflict", "cproj/app app --prefer-source --prefer-install=dist"+r),
		sc("prefer-source-dist-only", "cproj/app app --prefer-source"+r),
		sc("verbose", "cproj/app app -v"+r),
		sc("progress", "cproj/app app --repository=http://127.0.0.1:8899 --no-secure-http"),
		sc("relative-dir", "cproj/app sub/dir/app/"+r),
		sc("absolute-dir", "cproj/app /work/project/abs"+r),
		sc("plain", "cproj/plain plain"+r),
		sc("bin", "cproj/nodeps nodeps"+r),
		sc("packagist", "cproj/app app --no-progress"),
		withFiles(sc("not-empty", "cproj/app app"+r), map[string]string{"app/keep.txt": "keep\n"}),
		withFiles(sc("dir-is-file", "cproj/app app"+r), map[string]string{"app": "file\n"}),
		withFiles(sc("path", `local/app app --repository='{"type": "path", "url": "pkgs/app"}' --no-progress`), pathFiles),
		withFiles(sc("path-symlink", `local/app app --repository='{"type": "path", "url": "pkgs/app", "options": {"symlink": true}}' --no-progress`), pathFiles),
		withGit(sc("vcs", "vcs/gitlib lib"+vcsRepo)),
		withGit(sc("vcs-keep", "vcs/gitlib lib --keep-vcs"+vcsRepo)),
		withGit(sc("vcs-remove", "vcs/gitlib lib --remove-vcs"+vcsRepo)),
		withGit(sc("vcs-dev", "vcs/gitlib lib -s dev"+vcsRepo)),
		withGit(sc("vcs-branch", "vcs/gitlib lib 1.x-dev --keep-vcs"+vcsRepo)),
		sc("dot", "cproj/app ."+r),
		sc("audit-format", "cproj/app app --audit-format=plain"+r),
		sc("prefer-install-auto", "cproj/app app --prefer-install=auto"+r),
		sc("prefer-install-bad", "cproj/app app --prefer-install=foo"+r),
		sc("stability-rc-upper", "cproj/app app -s RC"+r),
		{Name: "create-project-composer-env", NoJSON: true, Env: map[string]string{"COMPOSER": "other.json"}, Command: "create-project cproj/app app" + r},
		{Name: "create-project-no-args", JSON: project(`{"acme/log": "^1.0"}`, ""), Files: map[string]string{"src/Kernel.php": appSrc}, Command: "create-project --no-progress", Probe: true},
		{Name: "create-project-no-args-no-json", NoJSON: true, Command: "create-project"},
	}
	return list
}
