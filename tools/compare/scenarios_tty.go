package main

// Interactive (pseudo terminal) scenarios: Composer's QuestionHelper/ConsoleIO
// prompts of composer/composer f83588392df43d4fd027d84d26de459946ef4e80 answered
// with scripted input (tty.go).

func init() { extraScenarios = append(extraScenarios, ttyScenarios) }

func ttyScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	twoKeys := project(`{"acme/log": "^1.0"}`, `"require-dev": {"acme/util": "^1.0"}`)
	removeQ := "Do you want to remove it from require-dev [yes]? "
	removeQd := "Do you want to remove it from require [yes]? "
	sc := []Scenario{
		// RemoveCommand: a package found in the other require key
		{Name: "tty-remove-alt-yes", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update", TTY: true, Input: []TTYStep{Answer(removeQ, "y")}},
		{Name: "tty-remove-alt-default", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update", TTY: true, Input: []TTYStep{Enter(removeQ)}},
		{Name: "tty-remove-alt-no", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update", TTY: true, Input: []TTYStep{Answer(removeQ, "no")}},
		{Name: "tty-remove-alt-invalid", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update", TTY: true, Input: []TTYStep{Answer(removeQ, "maybe"), Answer("Please answer yes, y, no, or n.", "N")}},
		{Name: "tty-remove-alt-eof", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update", TTY: true, Input: []TTYStep{CtrlD(removeQ)}},
		{Name: "tty-remove-alt-dev-wildcard", JSON: twoKeys, Files: files, Command: "remove --dev acme/l* --no-update", TTY: true, Input: []TTYStep{Answer(removeQd, "yes")}},
		{Name: "tty-remove-alt-dry-run", JSON: twoKeys, Files: files, Setup: []string{"update --no-progress"}, Command: "remove acme/util --dry-run --no-progress", TTY: true, Input: []TTYStep{Answer(removeQ, "y")}},
		{Name: "tty-remove-alt-no-interaction", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update -n", TTY: true},
		{Name: "tty-remove-alt-no-ansi", JSON: twoKeys, Files: files, Command: "remove acme/util --no-update --no-ansi", TTY: true, Input: []TTYStep{Answer(removeQ, "y")}},
	}
	sc = append(sc, ttyInitScenarios()...)
	sc = append(sc, ttyUpdateScenarios()...)
	sc = append(sc, ttyRequireScenarios()...)
	sc = append(sc, ttySelectScenarios()...)
	sc = append(sc, ttyCreateProjectScenarios()...)
	sc = append(sc, ttyPluginScenarios()...)
	parent := map[string]string{"../composer.json": project(`{"acme/log": "^1.0"}`, ""), "../src/Kernel.php": appSrc}
	parentQ := "No composer.json in current directory, do you want to use the one at /work? [y,n]? "
	sc = append(sc,
		Scenario{Name: "tty-parent-dir-yes", NoJSON: true, Files: parent, Command: "validate", TTY: true, Input: []TTYStep{Answer(parentQ, "y")}},
		Scenario{Name: "tty-parent-dir-no", NoJSON: true, Files: parent, Command: "validate", TTY: true, Input: []TTYStep{Answer(parentQ, "n")}},
		Scenario{Name: "tty-parent-dir-noninteractive", NoJSON: true, Files: parent, Command: "validate -n", TTY: true},
		Scenario{Name: "tty-parent-dir-init-exempt", NoJSON: true, Files: parent, Command: "init -n --name acme/x", TTY: true},
	)
	return sc
}

func ttyInitScenarios() []Scenario {
	// secure-http off for the loopback repository (Factory::createConfig)
	env := map[string]string{"GIT_CONFIG_GLOBAL": "/work/gitconfig", "GIT_CONFIG_NOSYSTEM": "1", "COMPOSER_HOME": "/work/ghome",
		"COMPOSER_DEFAULT_VENDOR": "acme", "COMPOSER_DEFAULT_AUTHOR": "Jane Doe", "COMPOSER_DEFAULT_EMAIL": "jane@example.com"}
	ghome := `{"config": {"secure-http": false, "autoloader-suffix": "InitTest"}}`
	files := map[string]string{"../gitconfig": "", "../ghome/config.json": ghome}
	gitFiles := map[string]string{"../gitconfig": "", "../ghome/config.json": ghome, ".git/HEAD": "ref: refs/heads/main\n", ".git/objects/.keep": "", ".git/refs/heads/.keep": "",
		".git/config": "[core]\n\trepositoryformatversion = 0\n\tbare = false\n", ".gitignore": "/node_modules"}
	snap := []string{".gitignore", ".gitattributes", "composer.lock"}
	repos := ` --repository http://127.0.0.1:8899 --repository '{"packagist.org":false}'`
	initCmd := "init" + repos
	name := "Package name (<vendor>/<name>) ["
	search := "Search for a package: "
	pick := "Enter package # to add, or the complete package name if it is not listed: "
	constraint := "Enter the version constraint to require (or leave blank to use the latest version): "
	reqQ := "Would you like to define your dependencies (require) interactively [yes]? "
	devQ := "Would you like to define your dev dependencies (require-dev) interactively [yes]? "
	autoload := "Add PSR-4 autoload mapping?"
	confirm := "Do you confirm generation [yes]? "
	install := "Would you like to install dependencies now [yes]? "
	head := func(n, desc, author, stab, typ, lic string) []TTYStep {
		return []TTYStep{Answer(name, n), Answer("Description [", desc), Answer("Author [", author), Answer("Minimum Stability [", stab), Answer("Package Type", typ), Answer("License [", lic)}
	}
	steps := func(parts ...[]TTYStep) []TTYStep {
		var out []TTYStep
		for _, p := range parts {
			out = append(out, p...)
		}
		return out
	}
	full := steps(
		[]TTYStep{Answer(name, "ttyfix/demo"), Answer("Description [", "Demo package"), Enter("Author ["), Answer("Minimum Stability [", "foo"), Answer("Minimum Stability [", "stable"), Answer("Package Type", "library"), Answer("License [", "MIT")},
		[]TTYStep{Answer(reqQ, "yes"), Answer(search, "tester"), Answer(pick, "0"), Enter(constraint),
			Answer(search, "ttyfix/lib"), Answer(constraint, "^1.0"), Answer(search, "nomatchzzz"), Enter(search)},
		[]TTYStep{Answer(devQ, "y"), Answer(search, "lint"), Answer(pick, "a b c"), Answer(pick, "ttyfix/linter ^2.0"), Enter(search)},
		[]TTYStep{Answer(autoload, "lib"), Answer(autoload, "lib/"), Answer(confirm, "yes"), Answer(install, "no")},
	)
	defaults := steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "n"), Answer(devQ, "n"), Enter(autoload), Enter(confirm)})
	sc := []Scenario{
		{Name: "tty-init-full", Command: initCmd, Input: full},
		{Name: "tty-init-defaults", Command: initCmd, Input: defaults},
		{Name: "tty-init-options-as-defaults", Command: initCmd + ` --name ttyfix/opt --description "From option" --author "Opt Author <opt@example.com>" --stability dev --type project --license GPL-2.0-only --autoload lib/`,
			Input: steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "no"), Answer(devQ, "no"), Enter(autoload), Answer(confirm, "y")})},
		{Name: "tty-init-require-options", Command: initCmd + ` --require ttyfix/lib --require-dev "ttyfix/tester:^1.0"`,
			Input: steps(head("ttyfix/req", "", "n", "", "", ""), []TTYStep{Answer(autoload, "n"), Enter(confirm), Answer(install, "n")})},
		{Name: "tty-init-abort", Command: initCmd, Input: steps(head("", "", "no", "", "", ""), []TTYStep{Answer(reqQ, "n"), Answer(devQ, "n"), Answer(autoload, "no"), Answer(confirm, "n")})},
		{Name: "tty-init-invalid-name", Command: initCmd, Input: []TTYStep{Answer(name, "Bad Name"), Answer(name, "acme/ok"), CtrlD("Description [")}},
		{Name: "tty-init-eof-first", Command: initCmd, Input: []TTYStep{CtrlD(name)}},
		{Name: "tty-init-eof-search", Command: initCmd, Input: steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "y"), CtrlD(search)})},
		{Name: "tty-init-bad-license", Command: initCmd, Input: head("", "", "", "", "", "Not-A-License")},
		{Name: "tty-init-invalid-author", Command: initCmd, Input: steps([]TTYStep{Enter(name), Enter("Description ["), Answer("Author [", "Jane <bad"), Answer("Author [", "Jane <not-an-email>"), Answer("Author [", "Solo Name"), Enter("Minimum Stability ["), Enter("Package Type"), Enter("License [")},
			[]TTYStep{Answer(reqQ, "n"), Answer(devQ, "n"), Enter(autoload), Enter(confirm)})},
		{Name: "tty-init-pick-attempts", Command: initCmd, Input: steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "y"), Answer(search, "ttyfix"), Answer(pick, "x y z"), Answer(pick, "a b c"), Answer(pick, "d e f")})},
		{Name: "tty-init-pick-empty", Command: initCmd, Input: steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "y"), Answer(search, "ttyfix"), Enter(pick), Answer(search, "oldlib"), Answer(pick, "1"), Answer(constraint, "  ^1.0  "), Enter(search), Answer(devQ, "n"), Enter(autoload), Enter(confirm), Answer(install, "no")})},
		{Name: "tty-init-git", Command: initCmd, Files: gitFiles, Input: steps(head("", "", "", "", "", ""), []TTYStep{Answer(reqQ, "y"), Answer(search, "ttyfix/lib"), Answer(constraint, "1.0.0"), Enter(search), Answer(devQ, "n"), Enter(autoload), Enter(confirm),
			Enter("Would you like the vendor directory added to your .gitignore [yes]? "), Answer("Would you like composer.lock marked as export-ignore in your .gitattributes [yes]? ", "y"), Enter(install)})},
		{Name: "tty-init-git-project-no", Command: initCmd, Files: gitFiles, Input: steps(head("", "", "", "", "project", ""), []TTYStep{Answer(reqQ, "n"), Answer(devQ, "n"), Answer(autoload, "n"), Enter(confirm),
			Answer("Would you like the vendor directory added to your .gitignore [yes]? ", "n")})},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].NoJSON = true
		sc[i].Env = env
		if sc[i].Files == nil {
			sc[i].Files = files
		}
		sc[i].Snapshot = snap
	}
	return sc
}

func ttyUpdateScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	proj := project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, "")
	old := []string{"update --no-progress --with acme/log:1.0.0 --with acme/util:1.0.0"}
	sel := "Select packages: (Select more than one value separated by comma) "
	cont := "Would you like to continue and update the above packages [yes]? "
	cont1 := "Would you like to continue and update the above package [yes]? "
	sc := []Scenario{
		{Name: "tty-update-interactive-one", Input: []TTYStep{Keys(sel, "acme/log\n"), Answer(cont1, "y")}},
		{Name: "tty-update-interactive-two", Input: []TTYStep{Keys(sel, "acme/util,acme/log\n"), Enter(cont)}},
		{Name: "tty-update-interactive-tab", Input: []TTYStep{Keys(sel, "acme/u\t\n"), Answer(cont1, "yes")}},
		{Name: "tty-update-interactive-arrows", Input: []TTYStep{Keys(sel, "a\x1b[B\x1b[B\x1b[A\n"), Answer(cont1, "y")}},
		{Name: "tty-update-interactive-backspace", Input: []TTYStep{Keys(sel, "acme/x\x7f\x7fl\n"), Answer(cont1, "y")}},
		{Name: "tty-update-interactive-invalid", Input: []TTYStep{Keys(sel, "nope\n")}},
		{Name: "tty-update-interactive-empty", Input: []TTYStep{Keys(sel, "\n")}},
		{Name: "tty-update-interactive-abort", Input: []TTYStep{Keys(sel, "acme/log\n"), Answer(cont1, "n")}},
		{Name: "tty-update-interactive-filter", Command: "update -i acme/u* --no-progress", Input: []TTYStep{Keys(sel, "acme/util\n"), Answer(cont1, "y")}},
		{Name: "tty-update-interactive-none", Setup: []string{"update --no-progress"}},
		{Name: "tty-update-interactive-no-lock", Setup: []string{}, Input: []TTYStep{Keys(sel, "acme/log\n"), Answer(cont1, "y")}},
		{Name: "tty-update-interactive-n", Command: "update -i -n --no-progress"},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].JSON = proj
		sc[i].Files = files
		if sc[i].Setup == nil {
			sc[i].Setup = old
		}
		if sc[i].Command == "" {
			sc[i].Command = "update -i --no-progress"
		}
	}
	return sc
}

func ttyRequireScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	base := project(`{"acme/log": "^1.0"}`, "")
	inDev := project(`{"acme/log": "^1.0"}`, `"require-dev": {"ttyfix/lib": "^1.0", "ttyfix/oldlib": "^1.0"}`)
	inReq := project(`{"acme/log": "^1.0", "ttyfix/tester": "^1.0"}`, "")
	devQ := "Do you want to re-run the command with --dev? [yes]? "
	moveQ := "Do you want to move this requirement? [no]? "
	moveQs := "Do you want to move these requirements? [no]? "
	rerunWith := "Do you want to re-run the command with --dev? [yes]? "
	rerunWithout := "Do you want to re-run the command without --dev? [yes]? "
	search := "Search for a package: "
	pick := "Enter package # to add, or the complete package name if it is not listed: "
	constraint := "Enter the version constraint to require (or leave blank to use the latest version): "
	sc := []Scenario{
		{Name: "tty-require-dev-recommend-yes", JSON: base, Command: "require ttyfix/tester --no-update", Input: []TTYStep{Enter(devQ)}},
		{Name: "tty-require-dev-recommend-no", JSON: base, Command: "require ttyfix/tester:^1.0 --no-update", Input: []TTYStep{Answer(devQ, "n")}},
		{Name: "tty-require-dev-recommend-two", JSON: base, Command: "require ttyfix/tester ttyfix/linter --no-update", Input: []TTYStep{Answer(devQ, "y")}},
		{Name: "tty-require-dev-recommend-mixed", JSON: base, Command: "require ttyfix/tester ttyfix/lib --no-update"},
		{Name: "tty-require-dev-recommend-update", JSON: base, Command: "require ttyfix/tester --no-progress", Input: []TTYStep{Answer(devQ, "yes")}},
		{Name: "tty-require-move-default", JSON: inDev, Command: "require ttyfix/lib:^1.3 --no-update", Input: []TTYStep{Enter(moveQ), Enter(rerunWith)}},
		{Name: "tty-require-move-yes", JSON: inDev, Command: "require ttyfix/lib:^1.3 --no-update", Input: []TTYStep{Answer(moveQ, "y")}},
		{Name: "tty-require-move-no-no", JSON: inDev, Command: "require ttyfix/lib:^1.3 ttyfix/oldlib:^1.0 --no-update", Input: []TTYStep{Answer(moveQs, "n"), Answer(rerunWith, "n")}},
		{Name: "tty-require-move-dev", JSON: inReq, Command: "require --dev ttyfix/tester:^1.2 --no-update", Input: []TTYStep{Answer(moveQ, "no"), Answer(rerunWithout, "yes")}},
		{Name: "tty-require-search", JSON: base, Command: "require --no-update", Input: []TTYStep{Answer(search, "lib"), Answer(pick, "1"), Enter(constraint), Answer(search, "acme/log"), Answer(search, "ttyfix/tester"), Answer(constraint, "1.0.0"), Enter(search), Answer(devQ, "n")}},
		{Name: "tty-require-search-none", JSON: base, Command: "require --no-update", Input: []TTYStep{Enter(search), Answer(devQ, "n")}},
		{Name: "tty-require-search-eof", JSON: base, Command: "require --no-update", Input: []TTYStep{CtrlD(search)}},
		{Name: "tty-require-virtual", JSON: base, Command: "require ttyfix/virtual --no-update", Input: []TTYStep{Answer("Which version constraint would you like to use? [*] ", "bad constraint!"), Answer("Which version constraint would you like to use? [*] ", "^1.0")}},
		{Name: "tty-require-virtual-default", JSON: base, Command: "require ttyfix/virtual --no-update", Input: []TTYStep{Enter("Which version constraint would you like to use? [*] ")}},
		{Name: "tty-require-similar", JSON: base, Command: "require ttyfix --no-update", Input: []TTYStep{Keys(" > ", "ttyfix/lib\n")}},
		{Name: "tty-require-similar-index", JSON: base, Command: "require ttyfix --no-update", Input: []TTYStep{Keys(" > ", "0\n")}},
		{Name: "tty-require-similar-empty", JSON: base, Command: "require ttyfix --no-update", Input: []TTYStep{Keys(" > ", "\n")}},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].Files = files
	}
	return sc
}

func ttySelectScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	scripts := project(`{"acme/log": "^1.0"}`, `"scripts": {"build": "echo building", "test": ["echo testing", "echo done"], "lint": "echo linting"}, "scripts-descriptions": {"build": "Build the app"}`)
	bins := project(`{"ttyfix/tester": "^1.0", "ttyfix/linter": "^2.0"}`, "")
	prompt := " > "
	sc := []Scenario{
		{Name: "tty-run-script-select", JSON: scripts, Command: "run-script", Input: []TTYStep{Keys(prompt, "build\n")}},
		{Name: "tty-run-script-select-tab", JSON: scripts, Command: "run-script", Input: []TTYStep{Keys(prompt, "te\t\n")}},
		{Name: "tty-run-script-select-invalid", JSON: scripts, Command: "run-script", Input: []TTYStep{Keys(prompt, "nope\n")}},
		{Name: "tty-run-script-select-empty", JSON: scripts, Command: "run-script", Input: []TTYStep{Keys(prompt, "\n")}},
		{Name: "tty-run-script-arg", JSON: scripts, Command: "run-script lint"},
		{Name: "tty-run-script-no-scripts", JSON: project(`{"acme/log": "^1.0"}`, ""), Command: "run-script"},
		{Name: "tty-exec-select", JSON: bins, Setup: []string{"update --no-progress"}, Command: "exec", Input: []TTYStep{Keys(prompt, "1\n")}},
		{Name: "tty-exec-select-name", JSON: bins, Setup: []string{"update --no-progress"}, Command: "exec", Input: []TTYStep{Keys(prompt, "li\x1b[B\n")}},
		{Name: "tty-exec-select-invalid", JSON: bins, Setup: []string{"update --no-progress"}, Command: "exec", Input: []TTYStep{Keys(prompt, "7\n")}},
		{Name: "tty-exec-no-bins", JSON: project(`{"acme/log": "^1.0"}`, ""), Command: "exec"},
		{Name: "tty-exec-list", JSON: bins, Setup: []string{"update --no-progress"}, Command: "exec --list"},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].Files = files
	}
	return sc
}

func ttyCreateProjectScenarios() []Scenario {
	const r = " --repository=http://127.0.0.1:8899 --no-secure-http --no-progress"
	vcsRepo := ` --repository='{"type": "vcs", "url": "/fixtures/git/gitlib.git"}' --no-progress`
	snap := []string{"app/composer.json", "app/composer.lock", "lib/composer.json", "mydir/composer.json", "cproj-app/composer.json"}
	ask := "New project directory ["
	vcsQ := "Do you want to remove the existing VCS (.git, .svn..) history? [y,n]? "
	sc := []Scenario{
		{Name: "tty-create-project-ask", Command: "create-project cproj/app --ask" + r, Input: []TTYStep{Answer(ask, "mydir")}},
		{Name: "tty-create-project-ask-default", Command: "create-project cproj/app --ask" + r, Input: []TTYStep{Enter(ask)}},
		{Name: "tty-create-project-ask-eof", Command: "create-project cproj/app --ask" + r, Input: []TTYStep{CtrlD(ask)}},
		{Name: "tty-create-project-ask-no-package", Command: "create-project --ask" + r},
		{Name: "tty-create-project-vcs-yes", Command: "create-project vcs/gitlib lib" + vcsRepo, GitFixtures: []string{"gitlib"}, Input: []TTYStep{Answer(vcsQ, "y")}},
		{Name: "tty-create-project-vcs-default", Command: "create-project vcs/gitlib lib" + vcsRepo, GitFixtures: []string{"gitlib"}, Input: []TTYStep{Enter(vcsQ)}},
		{Name: "tty-create-project-vcs-no", Command: "create-project vcs/gitlib lib" + vcsRepo, GitFixtures: []string{"gitlib"}, Input: []TTYStep{Answer(vcsQ, "n")}},
		{Name: "tty-create-project-vcs-remove-flag", Command: "create-project vcs/gitlib lib --remove-vcs" + vcsRepo, GitFixtures: []string{"gitlib"}},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].NoJSON = true
		sc[i].Snapshot = snap
	}
	return sc
}

func ttyPluginScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	ev := project(`{"plugfix/events-plugin": "1.0.0", "acme/log": "^1.0"}`, "")
	trust := `Do you trust "plugfix/events-plugin" to execute code and wish to enable it now? (writes "allow-plugins" to composer.json) [y,n,d,?] `
	sc := []Scenario{
		{Name: "tty-plugin-allow-yes", Input: []TTYStep{Answer(trust, "y")}},
		{Name: "tty-plugin-allow-no", Input: []TTYStep{Answer(trust, "n")}},
		{Name: "tty-plugin-allow-discard", Input: []TTYStep{Answer(trust, "d")}},
		{Name: "tty-plugin-allow-help", Input: []TTYStep{Answer(trust, "?"), Enter(trust), Answer(trust, "y")}},
		{Name: "tty-plugin-allow-eof", Input: []TTYStep{CtrlD(trust)}},
		{Name: "tty-plugin-allow-existing", JSON: withConfig(ev, `"allow-plugins": {"acme/*": true}, "sort-packages": true`), Input: []TTYStep{Answer(trust, "y")}},
	}
	for i := range sc {
		sc[i].TTY = true
		sc[i].Files = files
		if sc[i].JSON == "" {
			sc[i].JSON = ev
		}
		sc[i].Command = "update --no-progress"
	}
	return sc
}
