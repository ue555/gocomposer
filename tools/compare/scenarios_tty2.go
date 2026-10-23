package main

// Interactive (pseudo terminal) scenarios, part 2: the discard-changes
// prompt of GitDownloader::cleanChanges and the credential prompts of
// AuthHelper, GitHub/GitLab/Bitbucket utils and Util\Git (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80) answered with scripted input.

func init() { extraScenarios = append(extraScenarios, tty2Scenarios) }

func tty2Scenarios() []Scenario {
	return append(ttyDiscardScenarios(), ttyAuthScenarios()...)
}

func ttyDiscardScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	lib := gitRepo("vcs", "/fixtures/git/gitlib.git")
	libFx := []string{"gitlib"}
	json := vcsProject(`{"vcs/gitlib": "^1.0"}`, []string{lib}, "")
	modified := map[string]string{"vendor/vcs/gitlib/src/Lib.php": "<?php\nnamespace VcsLib;\n\nclass Lib { const VERSION = \"local\"; }\n"}
	modified2 := map[string]string{
		"vendor/vcs/gitlib/src/Lib.php":   "<?php\nnamespace VcsLib;\n\nclass Lib { const VERSION = \"local\"; }\n",
		"vendor/vcs/gitlib/composer.json": "{\"name\": \"vcs/gitlib\", \"license\": \"MIT\", \"autoload\": {\"psr-4\": {\"VcsLib\\\\\": \"src/\"}}}\n",
	}
	lib100 := []string{"update --no-progress --with vcs/gitlib:1.0.0"}
	upd := []string{"update --no-progress"}
	qUpd := "Discard changes [y,n,v,d,s,?]? "
	qRm := "Discard changes [y,n,v,d,?]? "
	d := func(name string, setup []string, after map[string]string, cmd string, input ...TTYStep) Scenario {
		return Scenario{Name: name, JSON: json, Files: files, GitFixtures: libFx, Setup: setup, After: after, Command: cmd, TTY: true, Input: input}
	}
	yes := d("tty-discard-update-yes", lib100, modified, "update --no-progress", Answer(qUpd, "y"))
	yes.Probe = true
	stash := d("tty-discard-update-stash", lib100, modified, "update --no-progress", Answer(qUpd, "s"))
	stash.Env = gitIdent
	return []Scenario{
		yes,
		d("tty-discard-update-no", lib100, modified, "update --no-progress", Answer(qUpd, "n")),
		d("tty-discard-update-help", lib100, modified, "update --no-progress", Answer(qUpd, "?"), Answer(qUpd, "x"), Enter(qUpd), Answer(qUpd, "y")),
		d("tty-discard-update-view", lib100, modified2, "update --no-progress", Answer(qUpd, "v"), Answer(qUpd, "d"), Answer(qUpd, "y")),
		stash,
		d("tty-discard-update-verbose", lib100, modified, "update --no-progress -v", Answer(qUpd, "y")),
		d("tty-discard-update-no-interaction", lib100, modified, "update --no-progress -n"),
		d("tty-discard-remove-yes", upd, modified, "remove vcs/gitlib --no-progress", Answer(qRm, "y")),
		d("tty-discard-remove-stash-is-help", upd, modified, "remove vcs/gitlib --no-progress", Answer(qRm, "s"), Answer(qRm, "n")),
		d("tty-discard-reinstall", upd, modified, "reinstall vcs/gitlib --no-progress", Answer(qRm, "y")),
		d("tty-discard-reinstall-no", upd, modified, "reinstall vcs/gitlib --no-progress", Answer(qRm, "n")),
		d("tty-discard-untracked-only", lib100, map[string]string{"vendor/vcs/gitlib/NOTES.txt": "local notes\n"}, "update --no-progress"),
	}
}

// Secrets typed into the credential prompts (the forge stub accepts them,
// see testdata/compare/vcs/forge-stub.php). Passwords and consumer secrets
// must never show up anywhere; tokens only in the auth.json Composer writes.
const (
	ttyGitLabPassword = "gl-s3cr3t-password"
	ttyBBKey          = "fixture-consumer-key"
	ttyBBSecret       = "bb-s3cr3t-consumer-secret"
	ttyGitPassword    = "git-s3cr3t-password"
)

func ttyAuthScenarios() []Scenario {
	forgeFx := []string{"gitlib", "forge"}
	req := `{"vcs/gitlib": "^1.0"}`
	noHost := `"github-expose-hostname": false`
	// COMPOSER_HOME inside the project so the stored auth.json is compared
	home := map[string]string{"COMPOSER_HOME": "/work/project/.home"}
	homeAuth := []string{".home/auth.json", "auth.json"}
	a := func(name, json, cmd string, input ...TTYStep) Scenario {
		return Scenario{Name: name, JSON: json, GitFixtures: forgeFx, Command: cmd, TTY: true, Input: input, Env: home, Snapshot: homeAuth}
	}
	withFiles := func(s Scenario, files map[string]string) Scenario { s.Files = files; return s }
	withProbe := func(s Scenario) Scenario { s.Probe = true; return s }
	plainHome := func(s Scenario) Scenario { s.Env, s.Snapshot = nil, nil; return s }
	withEnv := func(s Scenario, k, v string) Scenario {
		env := map[string]string{k: v}
		for ek, ev := range s.Env {
			env[ek] = ev
		}
		s.Env = env
		return s
	}

	// --- AuthHelper: a composer repository behind HTTP basic auth ---
	basic := func(path, storeAuths string) string {
		cfg := `"secure-http": false`
		if storeAuths != "" {
			cfg += `, "store-auths": ` + storeAuths
		}
		return forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/` + path + `"}`}, cfg)
	}
	user, pass := Answer("Username: ", "fixture-user"), Hidden("Password: ", forgeToken)
	storeQ := "Do you want to store credentials for 127.0.0.1:8900 in /work/project/.home/auth.json ? [Yn] "
	upd := "update --no-progress"
	localAuth := map[string]string{"auth.json": "{\n    \"bearer\": {\n        \"example.org\": \"unrelated\"\n    }\n}\n"}
	out := []Scenario{
		a("tty-auth-basic-store-yes", basic("auth401", ""), upd, user, pass, Answer(storeQ, "y")),
		a("tty-auth-basic-store-no", basic("auth401", ""), upd, user, pass, Answer(storeQ, "n")),
		a("tty-auth-basic-store-default", basic("auth401", ""), upd, user, pass, Enter(storeQ)),
		a("tty-auth-basic-store-invalid", basic("auth401", ""), upd, user, pass, Answer(storeQ, "maybe"), Answer("Please answer (y)es or (n)o", "Yes")),
		a("tty-auth-basic-store-true", basic("auth401", "true"), upd, user, pass),
		a("tty-auth-basic-store-false", basic("auth401", "false"), upd, user, pass),
		a("tty-auth-basic-403", basic("auth403", ""), upd, user, pass, Answer(storeQ, "y")),
		a("tty-auth-basic-wrong", basic("auth401", ""), upd, user, Hidden("Password: ", "wrong-password")),
		a("tty-auth-basic-empty", basic("auth401", ""), upd, Enter("Username: "), Hidden("Password: ", "")),
		a("tty-auth-basic-eof", basic("auth401", ""), upd, CtrlD("Username: ")),
		a("tty-auth-basic-eof-password", basic("auth401", ""), upd, user, CtrlD("Password: ")),
		a("tty-auth-basic-vv", basic("auth401", "true"), upd+" -vv", user, pass),
		withFiles(a("tty-auth-basic-local-authjson", basic("auth401", ""), upd, user, pass, Answer(storeQ, "y")), localAuth),
		a("tty-auth-basic-no-interaction", basic("auth401", ""), upd+" -n"),
	}

	// --- GitHub: token prompt (GitHub::authorizeOAuthInteractively) ---
	gh := forgeProject(req, []string{gitRepo("vcs", "https://github.com/vcs/gitlib-private")}, noHost)
	ghRL := forgeProject(req, []string{gitRepo("vcs", "https://github.com/vcs/gitlib-ratelimited")}, noHost)
	tokenQ := "Token (hidden): "
	localQ := "do you want to store the token there?"
	out = append(out,
		withProbe(a("tty-auth-github-token", gh, upd, Hidden(tokenQ, forgeToken))),
		a("tty-auth-github-token-vv", gh, upd+" -vv", Hidden(tokenQ, forgeToken)),
		a("tty-auth-github-empty-token", gh, upd, Hidden(tokenQ, "")),
		a("tty-auth-github-invalid-token", gh, upd, Hidden(tokenQ, "not-the-token")),
		a("tty-auth-github-eof", gh, upd, CtrlD(tokenQ)),
		withFiles(a("tty-auth-github-local-yes", gh, upd, Answer(localQ, "y"), Hidden(tokenQ, forgeToken)), localAuth),
		withFiles(a("tty-auth-github-local-no", gh, upd, Answer(localQ, "n"), Hidden(tokenQ, forgeToken)), localAuth),
		// the token does not lift the (stub's permanent) rate limit: asked
		// again by AuthHelper, then by GitHubDriver, then by AuthHelper
		a("tty-auth-github-ratelimit", ghRL, upd, Hidden(tokenQ, forgeToken), Hidden(tokenQ, ""), Hidden(tokenQ, ""), Hidden(tokenQ, "")),
		a("tty-auth-github-no-interaction", gh, upd+" -n"),
	)

	// --- GitLab: username/password exchanged for an OAuth token ---
	gl := forgeProject(req, []string{gitRepo("vcs", "https://gitlab.com/vcs/gitlib-private")}, "")
	glUser := Answer("Username: ", "fixture-user")
	out = append(out,
		withProbe(a("tty-auth-gitlab-password", gl, upd, glUser, Hidden("Password: ", ttyGitLabPassword))),
		a("tty-auth-gitlab-bad-then-good", gl, upd, glUser, Hidden("Password: ", "wrong"), Answer("Bad credentials", "fixture-user"), Hidden("Password: ", ttyGitLabPassword)),
		withFiles(a("tty-auth-gitlab-local-yes", gl, upd, Answer(localQ, "yes"), glUser, Hidden("Password: ", ttyGitLabPassword)), localAuth),
		a("tty-auth-gitlab-eof", gl, upd, CtrlD("Username: ")),
		a("tty-auth-gitlab-no-interaction", gl, upd+" -n"),
		withEnv(a("tty-auth-gitlab-refresh", gl, upd), "COMPOSER_AUTH", `{"gitlab-oauth": {"gitlab.com": {"token": "expired-token", "refresh-token": "fixture-refresh-token", "expires-at": 1}}}`),
		withEnv(a("tty-auth-gitlab-refresh-vv", gl, upd+" -vv"), "COMPOSER_AUTH", `{"gitlab-oauth": {"gitlab.com": {"token": "expired-token", "refresh-token": "fixture-refresh-token", "expires-at": 1}}}`),
		withEnv(a("tty-auth-gitlab-refresh-missing", gl, upd), "COMPOSER_AUTH", `{"gitlab-oauth": {"gitlab.com": {"token": "expired-token", "expires-at": 1}}}`),
		withEnv(a("tty-auth-gitlab-refresh-rejected", gl, upd, CtrlD("Username: ")), "COMPOSER_AUTH", `{"gitlab-oauth": {"gitlab.com": {"token": "expired-token", "refresh-token": "revoked", "expires-at": 1}}}`),
	)

	// --- Bitbucket: OAuth consumer exchanged for an access token (the
	// stored access-token-expiration is time() based, so COMPOSER_HOME stays
	// outside the compared tree) ---
	bb := forgeProject(req, []string{gitRepo("vcs", "https://bitbucket.org/vcs/gitlib-private")}, "")
	keyQ, secretQ := "Consumer Key (hidden): ", "Consumer Secret (hidden): "
	out = append(out,
		plainHome(a("tty-auth-bitbucket-consumer", bb, upd, Hidden(keyQ, ttyBBKey), Hidden(secretQ, ttyBBSecret))),
		plainHome(a("tty-auth-bitbucket-consumer-vv", bb, upd+" -vv", Hidden(keyQ, ttyBBKey), Hidden(secretQ, ttyBBSecret))),
		a("tty-auth-bitbucket-invalid", bb, upd, Hidden(keyQ, ttyBBKey), Hidden(secretQ, "wrong-secret")),
		a("tty-auth-bitbucket-no-key", bb, upd, Hidden(keyQ, "")),
		a("tty-auth-bitbucket-no-secret", bb, upd, Hidden(keyQ, ttyBBKey), Hidden(secretQ, "")),
	)

	// --- Util\Git: a private git repository over http ---
	git := func(u, storeAuths string) string {
		cfg := `"secure-http": false`
		if storeAuths != "" {
			cfg += `, "store-auths": ` + storeAuths
		}
		return forgeProject(req, []string{gitRepo("vcs", u)}, cfg)
	}
	gitURL := "http://127.0.0.1:8900/gitauth/gitlib.git"
	gitPass := Hidden("Password: ", ttyGitPassword)
	out = append(out,
		withProbe(a("tty-auth-git-store-yes", git(gitURL, ""), upd, user, gitPass, Answer(storeQ, "y"))),
		a("tty-auth-git-store-no", git(gitURL, ""), upd, user, gitPass, Answer(storeQ, "n")),
		a("tty-auth-git-store-true-v", git(gitURL, "true"), upd+" -v", user, gitPass),
		a("tty-auth-git-wrong", git(gitURL, ""), upd, user, Hidden("Password: ", "wrong-password")),
		a("tty-auth-git-default-user", git("http://fixture-user@127.0.0.1:8900/gitauth/gitlib.git", "false"), upd, Enter("Username: "), gitPass),
		a("tty-auth-git-eof", git(gitURL, ""), upd, CtrlD("Username: ")),
		a("tty-auth-git-no-interaction", git(gitURL, ""), upd+" -n"),
	)

	// --- Util\Git: clones from the forge hosts (GitHub token, GitLab
	// password grant, Bitbucket consumer) over https; git trusts the
	// fixture CA through GIT_SSL_CAINFO ---
	forgeGit := func(u string) string {
		return forgeProject(req, []string{`{"type": "git", "url": "` + u + `"}`}, noHost)
	}
	ca := func(s Scenario) Scenario { return withEnv(s, "GIT_SSL_CAINFO", forgeCA) }
	out = append(out,
		withProbe(ca(a("tty-auth-git-github", forgeGit("https://github.com/vcs/gitlib-private.git"), upd, Hidden(tokenQ, forgeToken)))),
		ca(a("tty-auth-git-github-no-token", forgeGit("https://github.com/vcs/gitlib-private.git"), upd, Hidden(tokenQ, ""))),
		withProbe(ca(a("tty-auth-git-gitlab", forgeGit("https://gitlab.com/vcs/gitlib-private.git"), upd, glUser, Hidden("Password: ", ttyGitLabPassword)))),
		ca(plainHome(a("tty-auth-git-bitbucket", forgeGit("https://bitbucket.org/vcs/gitlib-private.git"), upd, Hidden(keyQ, ttyBBKey), Hidden(secretQ, ttyBBSecret)))),
	)
	return out
}
