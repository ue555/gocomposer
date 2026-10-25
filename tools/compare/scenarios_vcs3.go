package main

// vcs3Scenarios compare the Subversion, Mercurial, Fossil and Perforce
// support (Composer f83588392df43d4fd027d84d26de459946ef4e80:
// Repository/Vcs/{Svn,Hg,Fossil,Perforce}Driver, Downloader/{Svn,Hg,Fossil,
// Perforce}Downloader, Util/{Svn,Hg,Perforce}). The repositories are built by
// testdata/compare/vcs/{svnlib,hglib,fossillib,p4stub}.sh inside the oracle
// image (which has the svn, hg and fossil clients) and are read through
// file:// urls, local paths and loopback servers (svnserve, hg serve,
// fossil server) with --network=none.
//
// Perforce: the Helix Core server and client (p4d/p4) are proprietary
// downloads under the Perforce EULA, not packaged by Debian and not freely
// redistributable, so the oracle image has no real p4. The vcs-perforce-*
// scenarios run both tools against the same scripted p4 stand-in
// (testdata/compare/vcs/p4stub.sh) that answers the commands Util/Perforce
// issues with fixed output; they compare Composer's and gocomposer's
// handling of that output, not the behaviour of a real Perforce server.
func init() { extraScenarios = append(extraScenarios, vcs3Scenarios) }

const (
	svnLibURL  = "file:///fixtures/git/svn/svnlib"
	svnMonoURL = "file:///fixtures/git/svn/svnmono"
	svnServe   = "svn://127.0.0.1:3690"
)

func vcs3Scenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	svnFx := []string{"svnlib"}
	svn := gitRepo("svn", svnLibURL)
	sv := func(name, json, cmd string, extra func(*Scenario)) Scenario {
		s := Scenario{Name: name, JSON: json, Files: files, GitFixtures: svnFx, Command: cmd}
		if extra != nil {
			extra(&s)
		}
		return s
	}
	probe := func(s *Scenario) { s.Probe = true }
	setup := func(cmds ...string) func(*Scenario) {
		return func(s *Scenario) { s.Setup = cmds; s.Probe = true }
	}
	after := func(cmds []string, files map[string]string) func(*Scenario) {
		return func(s *Scenario) { s.Setup = cmds; s.After = files }
	}
	env := func(m map[string]string) func(*Scenario) { return func(s *Scenario) { s.Env = m } }
	req := `{"vcs3/svnlib": "^1.0"}`
	svnModified := map[string]string{"vendor/vcs3/svnlib/src/Lib.php": "<?php // local change\n"}
	svn100 := []string{"update --no-progress --with vcs3/svnlib:1.0.0"}
	// the pty scenarios use svnserve: with a file:// url (no host) Util/Svn
	// raises PHP 8.5's "null as an array offset" deprecation, which the pty
	// transcript (stdout and stderr merged, decorated) cannot filter out
	svnSrv := gitRepo("svn", svnServe+"/svnlib")

	out := []Scenario{
		// --- SvnDriver ---
		sv("vcs-svn-update", vcsProject(req, []string{svn}, ""), "update --no-progress", probe),
		sv("vcs-svn-update-v", vcsProject(req, []string{svn}, ""), "update --no-progress -v", nil),
		sv("vcs-svn-update-vv", vcsProject(req, []string{svn}, ""), "update --no-progress -vv", nil),
		sv("vcs-svn-detect", vcsProject(req, []string{gitRepo("vcs", svnLibURL)}, ""), "update --no-progress", nil),
		sv("vcs-svn-detect-serve", vcsProject(req, []string{gitRepo("vcs", svnServe+"/svnlib")}, ""), "update --no-progress -vv", nil),
		sv("vcs-svn-abs-path", vcsProject(req, []string{gitRepo("svn", "/fixtures/git/svn/svnlib")}, ""), "update --no-progress", probe),
		sv("vcs-svn-trunk-url", vcsProject(req, []string{gitRepo("svn", svnLibURL+"/trunk")}, ""), "update --no-progress", nil),
		sv("vcs-svn-dev-trunk", vcsProject(`{"vcs3/svnlib": "dev-trunk"}`, []string{svn}, ""), "update --no-progress", probe),
		sv("vcs-svn-branch", vcsProject(`{"vcs3/svnlib": "2.x-dev"}`, []string{svn}, ""), "update --no-progress", probe),
		sv("vcs-svn-show", vcsProject(req, []string{svn}, ""), "show vcs3/svnlib --all", nil),
		sv("vcs-svn-package-path", vcsProject(`{"vcs3/svnmono": "^1.0"}`, []string{`{"type": "svn", "url": "` + svnMonoURL + `", "package-path": "pkg"}`}, ""), "update --no-progress -vv", probe),
		sv("vcs-svn-no-tags-path", vcsProject(`{"vcs3/svnlib": "dev-trunk"}`, []string{`{"type": "svn", "url": "` + svnLibURL + `", "tags-path": false, "branches-path": false}`}, ""), "update --no-progress -vv", nil),
		sv("vcs-svn-custom-trunk", vcsProject(`{"vcs3/svnlib": "*"}`, []string{`{"type": "svn", "url": "` + svnLibURL + `", "trunk-path": "branches/2.x"}`}, ""), "update --no-progress -vv", nil),
		sv("vcs-svn-missing", vcsProject(req, []string{gitRepo("svn", "file:///fixtures/git/svn/nope")}, ""), "update --no-progress", nil),
		sv("vcs-svn-cache", vcsOnlyProject(req, []string{svn}), "update --no-progress -vv", func(s *Scenario) {
			s.Setup = []string{"update --no-progress --no-install"}
			s.Env = map[string]string{"COMPOSER_CACHE_DIR": "/work/project/.cache-keep"}
		}),
		// svnserve on loopback: secure-http, secure-svn-domains, credentials
		sv("vcs-svn-serve", vcsProject(req, []string{gitRepo("svn", svnServe+"/svnlib")}, ""), "update --no-progress", probe),
		sv("vcs-svn-serve-secure", withConfig(vcsProject(req, []string{gitRepo("svn", svnServe+"/svnlib")}, ""), `"secure-http": true`), "update --no-progress", nil),
		sv("vcs-svn-serve-secure-domain", withConfig(vcsProject(req, []string{gitRepo("svn", svnServe+"/svnlib")}, ""), `"secure-http": true, "secure-svn-domains": ["127.0.0.1"]`), "update --no-progress", nil),
		sv("vcs-svn-auth-noninteractive", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", svnServe+"/svnpriv")}, ""), "update --no-progress", nil),
		sv("vcs-svn-auth-http-basic", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", svnServe+"/svnpriv")}, ""), "update --no-progress -vv", env(map[string]string{"COMPOSER_AUTH": `{"http-basic": {"127.0.0.1": {"username": "fixture-user", "password": "s3cr3t-svn"}}}`})),
		sv("vcs-svn-auth-wrong", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", svnServe+"/svnpriv")}, ""), "update --no-progress", env(map[string]string{"COMPOSER_AUTH": `{"http-basic": {"127.0.0.1": {"username": "fixture-user", "password": "wrong"}}}`})),
		sv("vcs-svn-auth-url", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", "svn://fixture-user:s3cr3t-svn@127.0.0.1:3690/svnpriv")}, ""), "update --no-progress", probe),
		sv("vcs-svn-auth-no-cache", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{`{"type": "svn", "url": "svn://fixture-user:s3cr3t-svn@127.0.0.1:3690/svnpriv", "svn-cache-credentials": false}`}, ""), "update --no-progress -v", nil),
		sv("vcs-svn-auth-tty", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", svnServe+"/svnpriv")}, ""), "update --no-progress --no-install", func(s *Scenario) {
			s.TTY = true
			s.Input = []TTYStep{Answer("Username: ", "fixture-user"), Hidden("Password: ", "s3cr3t-svn"), Answer("Should Subversion cache these credentials? (yes/no) ", "no")}
		}),
		sv("vcs-svn-auth-tty-wrong", vcsProject(`{"vcs3/svnpriv": "^1.0"}`, []string{gitRepo("svn", svnServe+"/svnpriv")}, ""), "update --no-progress --no-install", func(s *Scenario) {
			s.TTY = true
			s.Input = []TTYStep{Answer("Username: ", "fixture-user"), Hidden("Password: ", "nope"), Answer("Should Subversion cache these credentials? (yes/no) ", "no")}
		}),

		// VersionGuesser::guessSvnVersion through a path repository
		sv("vcs-svn-path-guess", vcsProject(`{"vcs3/svnlib": "@dev"}`, []string{`{"type": "path", "url": "/fixtures/git/svn/wc-trunk"}`}, ""), "update --no-progress", probe),
		sv("vcs-svn-path-guess-branch", vcsProject(`{"vcs3/svnlib": "@dev"}`, []string{`{"type": "path", "url": "/fixtures/git/svn/wc-2x"}`}, ""), "update --no-progress", probe),

		// --- SvnDownloader ---
		sv("source-svn-install", vcsOnlyProject(req, []string{svn}), "install --no-progress", setup("update --no-progress --no-install")),
		sv("source-svn-install-v", vcsOnlyProject(req, []string{svn}), "install --no-progress -v", setup("update --no-progress --no-install")),
		sv("source-svn-update", vcsProject(req, []string{svn}, ""), "update --no-progress", setup(svn100...)),
		// VcsDownloader::update prints the commit logs in verbose mode:
		// SvnDownloader::getCommitLogs runs Svn::executeLocal, which appends
		// "-- <empty url> <path>" and svn rejects it (E200007), so Composer
		// itself fails the update; gocomposer fails the same way, but prints
		// no PHP exception trace (the only difference)
		sv("source-svn-update-v", vcsProject(req, []string{svn}, ""), "update --no-progress -v", func(s *Scenario) {
			s.Setup = svn100
			s.KnownDiff = "Composer's SvnDownloader::getCommitLogs fails (svn E200007: Svn::executeLocal passes an empty url and the path as two targets); both tools fail identically, only Composer's -v exception trace is absent from gocomposer"
		}),
		sv("source-svn-to-trunk", vcsProject(req, []string{svn}, ""), "update --no-progress", func(s *Scenario) {
			s.Setup = []string{"update --no-progress"}
			s.After = map[string]string{"composer.json": vcsProject(`{"vcs3/svnlib": "dev-trunk"}`, []string{svn}, "")}
			s.Probe = true
		}),
		sv("source-svn-local-changes", vcsProject(req, []string{svn}, ""), "update --no-progress", after(svn100, svnModified)),
		sv("source-svn-discard-true", withConfig(vcsProject(req, []string{svn}, ""), `"discard-changes": true`), "update --no-progress", func(s *Scenario) { s.Setup = svn100; s.After = svnModified; s.Probe = true }),
		sv("source-svn-discard-stash", withConfig(vcsProject(req, []string{svn}, ""), `"discard-changes": "stash"`), "update --no-progress", after(svn100, svnModified)),
		sv("source-svn-uninstall-changes", vcsProject(req, []string{svn}, ""), "remove vcs3/svnlib --no-progress", after([]string{"update --no-progress"}, svnModified)),
		sv("source-svn-uninstall-discard", withConfig(vcsProject(req, []string{svn}, ""), `"discard-changes": true`), "remove vcs3/svnlib --no-progress", after([]string{"update --no-progress"}, svnModified)),
		sv("source-svn-remove", vcsProject(req, []string{svn}, ""), "remove vcs3/svnlib --no-progress", setup("update --no-progress")),
		sv("source-svn-status-clean", vcsProject(req, []string{svn}, ""), "status", setup("update --no-progress")),
		sv("source-svn-status", vcsProject(req, []string{svn}, ""), "status", after([]string{"update --no-progress"}, svnModified)),
		sv("source-svn-status-v", vcsProject(req, []string{svn}, ""), "status -v", after([]string{"update --no-progress"}, svnModified)),
		sv("source-svn-tty-discard-yes", vcsProject(req, []string{svnSrv}, ""), "update --no-progress", func(s *Scenario) {
			s.Setup, s.After, s.TTY = svn100, svnModified, true
			s.Input = []TTYStep{Answer("Discard changes [y,n,v,?]? ", "?"), Answer("Discard changes [y,n,v,?]? ", "v"), Answer("Discard changes [y,n,v,?]? ", "y")}
		}),
		sv("source-svn-tty-discard-no", vcsProject(req, []string{svnSrv}, ""), "update --no-progress", func(s *Scenario) {
			s.Setup, s.After, s.TTY = svn100, svnModified, true
			s.Input = []TTYStep{Answer("Discard changes [y,n,v,?]? ", "n")}
		}),
		sv("source-svn-create-project", "", `create-project vcs3/svnlib app --repository='{"type": "svn", "url": "`+svnLibURL+`"}' --no-progress`, func(s *Scenario) {
			s.NoJSON = true
			s.Snapshot = []string{"app/composer.json", "app/src/Lib.php"}
		}),
		sv("source-svn-create-project-keep-vcs", "", `create-project vcs3/svnlib app 1.0.0 --keep-vcs --repository='{"type": "svn", "url": "`+svnLibURL+`"}' --no-progress`, func(s *Scenario) {
			s.NoJSON = true
			s.Snapshot = []string{"app/composer.json", "app/src/Lib.php"}
		}),
		sv("source-svn-lock-reference", vcsOnlyProject(`{"vcs3/svnlib": "dev-trunk"}`, []string{svn}), "update --no-progress --lock", setup("update --no-progress")),
	}
	out = append(out, hgScenarios(files)...)
	out = append(out, fossilScenarios(files)...)
	out = append(out, perforceScenarios(files)...)
	return out
}

// vcs3Family builds scenarios over one fixture script.
func vcs3Family(files map[string]string, fx string) func(name, json, cmd string, extra func(*Scenario)) Scenario {
	return func(name, json, cmd string, extra func(*Scenario)) Scenario {
		s := Scenario{Name: name, JSON: json, Files: files, GitFixtures: []string{fx}, Command: cmd}
		if extra != nil {
			extra(&s)
		}
		return s
	}
}

func withSetup(cmds ...string) func(*Scenario) {
	return func(s *Scenario) { s.Setup = cmds; s.Probe = true }
}

func withSetupAfter(cmds []string, files map[string]string) func(*Scenario) {
	return func(s *Scenario) { s.Setup = cmds; s.After = files }
}

const (
	hgLibPath = "/fixtures/git/hg/hglib"
	hgServe   = "http://127.0.0.1:8901/"
)

func hgScenarios(files map[string]string) []Scenario {
	hs := vcs3Family(files, "hglib")
	probe := func(s *Scenario) { s.Probe = true }
	hg := gitRepo("hg", hgLibPath)
	srv := gitRepo("hg", hgServe)
	req := `{"vcs3/hglib": "^1.0"}`
	modified := map[string]string{"vendor/vcs3/hglib/src/Lib.php": "<?php // local change\n"}
	hg100 := []string{"update --no-progress --with vcs3/hglib:1.0.0"}
	return []Scenario{
		// HgDriver: local repository, hg serve (cache-vcs-dir mirror)
		hs("vcs-hg-update", vcsProject(req, []string{hg}, ""), "update --no-progress", probe),
		hs("vcs-hg-update-vv", vcsProject(req, []string{hg}, ""), "update --no-progress -vv", nil),
		hs("vcs-hg-detect", vcsProject(req, []string{gitRepo("vcs", hgLibPath)}, ""), "update --no-progress -v", nil),
		hs("vcs-hg-dev-default", vcsProject(`{"vcs3/hglib": "dev-default"}`, []string{hg}, ""), "update --no-progress", probe),
		hs("vcs-hg-branch", vcsProject(`{"vcs3/hglib": "2.x-dev"}`, []string{hg}, ""), "update --no-progress", probe),
		hs("vcs-hg-bookmark", vcsProject(`{"vcs3/hglib": "dev-feature"}`, []string{hg}, ""), "update --no-progress", probe),
		hs("vcs-hg-show", vcsProject(req, []string{hg}, ""), "show vcs3/hglib --all", nil),
		hs("vcs-hg-serve", vcsProject(req, []string{srv}, ""), "update --no-progress", probe),
		hs("vcs-hg-serve-vv", vcsProject(req, []string{srv}, ""), "update --no-progress -vv", nil),
		hs("vcs-hg-serve-detect", vcsProject(req, []string{gitRepo("vcs", hgServe)}, ""), "update --no-progress -v", nil),
		hs("vcs-hg-serve-cached", vcsOnlyProject(req, []string{srv}), "update --no-progress -v", func(s *Scenario) {
			s.Setup = []string{"update --no-progress --no-install"}
			s.Env = map[string]string{"COMPOSER_CACHE_DIR": "/work/project/.cache-keep"}
		}),
		hs("vcs-hg-serve-secure", withConfig(vcsProject(req, []string{srv}, ""), `"secure-http": true`), "update --no-progress", nil),
		hs("vcs-hg-serve-missing", vcsProject(req, []string{gitRepo("hg", "http://127.0.0.1:8901/nope")}, ""), "update --no-progress", nil),
		hs("vcs-hg-serve-auth", vcsProject(req, []string{gitRepo("hg", "http://127.0.0.1:8901/nope")}, ""), "update --no-progress", func(s *Scenario) {
			s.Env = map[string]string{"COMPOSER_AUTH": `{"http-basic": {"127.0.0.1:8901": {"username": "fixture-user", "password": "s3cr3t-hg"}}}`}
		}),
		hs("vcs-hg-missing", vcsProject(req, []string{gitRepo("hg", "/fixtures/git/hg/nope")}, ""), "update --no-progress", nil),
		// VersionGuesser::guessHgVersion through a path repository
		hs("vcs-hg-path-guess", vcsProject(`{"vcs3/hglib": "@dev"}`, []string{`{"type": "path", "url": "/fixtures/git/hg/hgwc"}`}, ""), "update --no-progress", probe),
		hs("vcs-hg-path-guess-feature", vcsProject(`{"vcs3/hglib": "@dev"}`, []string{`{"type": "path", "url": "/fixtures/git/hg/hgwc-topic"}`}, ""), "update --no-progress", probe),
		// HgDownloader
		hs("source-hg-create-project", "", `create-project vcs3/hglib app --repository='{"type": "hg", "url": "`+hgLibPath+`"}' --no-progress`, func(s *Scenario) {
			s.NoJSON = true
			s.Snapshot = []string{"app/composer.json", "app/src/Lib.php"}
		}),
		hs("source-hg-install", vcsOnlyProject(req, []string{hg}), "install --no-progress", withSetup("update --no-progress --no-install")),
		hs("source-hg-install-serve", vcsOnlyProject(req, []string{srv}), "install --no-progress -v", withSetup("update --no-progress --no-install")),
		hs("source-hg-update", vcsProject(req, []string{hg}, ""), "update --no-progress", withSetup(hg100...)),
		hs("source-hg-update-v", vcsProject(req, []string{hg}, ""), "update --no-progress -v", withSetup(hg100...)),
		hs("source-hg-downgrade-v", vcsProject(req, []string{hg}, ""), "update --no-progress --with vcs3/hglib:1.0.0 -v", withSetup("update --no-progress")),
		hs("source-hg-to-dev", vcsProject(req, []string{hg}, ""), "update --no-progress", func(s *Scenario) {
			s.Setup = []string{"update --no-progress"}
			s.After = map[string]string{"composer.json": vcsProject(`{"vcs3/hglib": "dev-default"}`, []string{hg}, "")}
			s.Probe = true
		}),
		hs("source-hg-local-changes", vcsProject(req, []string{hg}, ""), "update --no-progress", withSetupAfter(hg100, modified)),
		hs("source-hg-discard-true", withConfig(vcsProject(req, []string{hg}, ""), `"discard-changes": true`), "update --no-progress", withSetupAfter(hg100, modified)),
		hs("source-hg-uninstall-changes", vcsProject(req, []string{hg}, ""), "remove vcs3/hglib --no-progress", withSetupAfter([]string{"update --no-progress"}, modified)),
		hs("source-hg-remove", vcsProject(req, []string{hg}, ""), "remove vcs3/hglib --no-progress", withSetup("update --no-progress")),
		hs("source-hg-status-clean", vcsProject(req, []string{hg}, ""), "status", withSetup("update --no-progress")),
		hs("source-hg-status", vcsProject(req, []string{hg}, ""), "status", withSetupAfter([]string{"update --no-progress"}, modified)),
		hs("source-hg-status-v", vcsProject(req, []string{hg}, ""), "status -v", withSetupAfter([]string{"update --no-progress"}, modified)),
	}
}

const (
	fossilRepoPath = "/fixtures/git/fossil/fossillib.fossil"
	fossilServe    = "http://127.0.0.1:8902/"
)

func fossilScenarios(files map[string]string) []Scenario {
	family := vcs3Family(files, "fossillib")
	// fossil refuses to open/update checkouts without a user name (USER is
	// not set in the container): every scenario but source-fossil-no-user
	// sets one
	fs := func(name, json, cmd string, extra func(*Scenario)) Scenario {
		return family(name, json, cmd, func(s *Scenario) {
			s.Env = map[string]string{"USER": "fixture"}
			if extra != nil {
				extra(s)
			}
		})
	}
	probe := func(s *Scenario) { s.Probe = true }
	fl := gitRepo("fossil", fossilRepoPath)
	srv := gitRepo("fossil", fossilServe)
	req := `{"vcs3/fossillib": "^1.0"}`
	modified := map[string]string{"vendor/vcs3/fossillib/src/Lib.php": "<?php // local change\n"}
	f100 := []string{"update --no-progress --with vcs3/fossillib:1.0.0"}
	return []Scenario{
		// FossilDriver: repository file (cloned into the cache), checkout
		// directory, fossil server
		fs("vcs-fossil-update", vcsProject(req, []string{fl}, ""), "update --no-progress", probe),
		fs("vcs-fossil-update-vv", vcsProject(req, []string{fl}, ""), "update --no-progress -vv", nil),
		fs("vcs-fossil-detect", vcsProject(req, []string{gitRepo("vcs", fossilRepoPath)}, ""), "update --no-progress -v", nil),
		fs("vcs-fossil-checkout-dir", vcsProject(req, []string{gitRepo("fossil", "/fixtures/git/fossil/checkout")}, ""), "update --no-progress --no-install -v", nil),
		fs("vcs-fossil-dev-trunk", vcsProject(`{"vcs3/fossillib": "dev-trunk"}`, []string{fl}, ""), "update --no-progress", probe),
		fs("vcs-fossil-branch", vcsProject(`{"vcs3/fossillib": "2.x-dev"}`, []string{fl}, ""), "update --no-progress", probe),
		fs("vcs-fossil-show", vcsProject(req, []string{fl}, ""), "show vcs3/fossillib --all", nil),
		fs("vcs-fossil-serve", vcsProject(req, []string{srv}, ""), "update --no-progress", probe),
		// without a composer repository nothing creates cache-repo-dir
		// before FossilDriver clones into it: Composer fails (SQLITE_CANTOPEN)
		fs("vcs-fossil-no-repo-cache-dir", vcsOnlyProject(req, []string{srv}), "update --no-progress", nil),
		fs("vcs-fossil-serve-secure", withConfig(vcsProject(req, []string{srv}, ""), `"secure-http": true`), "update --no-progress", nil),
		fs("vcs-fossil-missing", vcsProject(req, []string{gitRepo("fossil", "/fixtures/git/fossil/nope.fossil")}, ""), "update --no-progress", nil),
		// VersionGuesser::guessFossilVersion through a path repository
		fs("vcs-fossil-path-guess", vcsProject(`{"vcs3/fossillib": "@dev"}`, []string{`{"type": "path", "url": "/fixtures/git/fossil/checkout"}`}, ""), "update --no-progress", probe),
		// FossilDownloader
		family("source-fossil-no-user", vcsProject(req, []string{fl}, ""), "update --no-progress", nil),
		fs("source-fossil-create-project", "", `create-project vcs3/fossillib app --repository='{"type": "fossil", "url": "`+fossilRepoPath+`"}' --no-progress`, func(s *Scenario) {
			s.NoJSON = true
			s.Snapshot = []string{"app/composer.json", "app/src/Lib.php"}
		}),
		fs("source-fossil-install", vcsOnlyProject(req, []string{fl}), "install --no-progress", withSetup("update --no-progress --no-install")),
		fs("source-fossil-update", vcsProject(req, []string{fl}, ""), "update --no-progress", withSetup(f100...)),
		fs("source-fossil-update-v", vcsProject(req, []string{fl}, ""), "update --no-progress -v", withSetup(f100...)),
		fs("source-fossil-to-dev", vcsProject(req, []string{fl}, ""), "update --no-progress", func(s *Scenario) {
			s.Setup = []string{"update --no-progress"}
			s.After = map[string]string{"composer.json": vcsProject(`{"vcs3/fossillib": "dev-trunk"}`, []string{fl}, "")}
			s.Probe = true
		}),
		fs("source-fossil-local-changes", vcsProject(req, []string{fl}, ""), "update --no-progress", withSetupAfter(f100, modified)),
		fs("source-fossil-uninstall-changes", vcsProject(req, []string{fl}, ""), "remove vcs3/fossillib --no-progress", withSetupAfter([]string{"update --no-progress"}, modified)),
		fs("source-fossil-remove", vcsProject(req, []string{fl}, ""), "remove vcs3/fossillib --no-progress", withSetup("update --no-progress")),
		fs("source-fossil-status-clean", vcsProject(req, []string{fl}, ""), "status", withSetup("update --no-progress")),
		fs("source-fossil-status-v", vcsProject(req, []string{fl}, ""), "status -v", withSetupAfter([]string{"update --no-progress"}, modified)),
	}
}

// p4Path puts the scripted p4 of p4stub.sh first in PATH.
var p4Env = map[string]string{"PATH": "/fixtures/git/_p4/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}

func perforceScenarios(files map[string]string) []Scenario {
	ps := vcs3Family(files, "p4stub")
	withEnv := func(extra map[string]string, then func(*Scenario)) func(*Scenario) {
		return func(s *Scenario) {
			s.Env = map[string]string{}
			for k, v := range p4Env {
				s.Env[k] = v
			}
			for k, v := range extra {
				s.Env[k] = v
			}
			if then != nil {
				then(s)
			}
		}
	}
	probe := func(s *Scenario) { s.Probe = true }
	repoCfg := func(extra string) string {
		return `{"type": "perforce", "url": "127.0.0.1:1666", "depot": "depot", "p4user": "fixture"` + extra + `}`
	}
	p4 := repoCfg("")
	req := `{"vcs3/p4lib": "^1.0"}`
	userEnv := map[string]string{"P4USER": "fixture"}
	return []Scenario{
		ps("vcs-perforce-update", vcsProject(req, []string{p4}, ""), "update --no-progress", withEnv(nil, probe)),
		ps("vcs-perforce-update-vv", vcsProject(req, []string{p4}, ""), "update --no-progress -vv", withEnv(nil, nil)),
		ps("vcs-perforce-detect", vcsProject(req, []string{`{"type": "vcs", "url": "127.0.0.1:1666", "depot": "depot", "p4user": "fixture"}`}, ""), "update --no-progress --no-install -v", withEnv(nil, nil)),
		ps("vcs-perforce-detect-p4-url", vcsProject(req, []string{`{"type": "vcs", "url": "p4.test:1666", "depot": "depot"}`}, ""), "update --no-progress --no-install", withEnv(nil, nil)),
		ps("vcs-perforce-dev", vcsProject(`{"vcs3/p4lib": "dev-master"}`, []string{p4}, ""), "update --no-progress", withEnv(nil, probe)),
		ps("vcs-perforce-show", vcsProject(req, []string{p4}, ""), "show vcs3/p4lib --all", withEnv(nil, nil)),
		ps("vcs-perforce-env-user", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1666", "depot": "depot"}`}, ""), "update --no-progress --no-install", withEnv(userEnv, nil)),
		ps("vcs-perforce-login", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1666", "depot": "depot", "p4user": "locked", "p4password": "p4-s3cr3t"}`}, ""), "update --no-progress --no-install", withEnv(nil, nil)),
		ps("vcs-perforce-login-wrong", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1666", "depot": "depot", "p4user": "locked", "p4password": "nope"}`}, ""), "update --no-progress --no-install", withEnv(nil, nil)),
		ps("vcs-perforce-login-env", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1666", "depot": "depot", "p4user": "locked"}`}, ""), "update --no-progress --no-install", withEnv(map[string]string{"P4PASSWD": "p4-s3cr3t"}, nil)),
		ps("vcs-perforce-stream", vcsProject(`{"vcs3/p4lib": "*"}`, []string{`{"type": "perforce", "url": "127.0.0.1:1666", "depot": "streams", "branch": "main", "p4user": "fixture"}`}, ""), "update --no-progress --no-install -v", withEnv(nil, nil)),
		ps("vcs-perforce-no-depot", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1666"}`}, ""), "update --no-progress", withEnv(nil, nil)),
		ps("vcs-perforce-bad-server", vcsProject(req, []string{`{"type": "perforce", "url": "127.0.0.1:1667", "depot": "depot", "p4user": "fixture"}`}, ""), "update --no-progress", withEnv(nil, nil)),
		ps("vcs-perforce-invalid-port", vcsProject(req, []string{`{"type": "perforce", "url": "rsh:/bin/true", "depot": "depot", "p4user": "fixture"}`}, ""), "update --no-progress", withEnv(nil, nil)),
		ps("source-perforce-install", vcsOnlyProject(req, []string{p4}), "install --no-progress", withEnv(userEnv, func(s *Scenario) {
			s.Setup = []string{"update --no-progress --no-install"}
			s.Probe = true
		})),
		ps("source-perforce-update", vcsProject(req, []string{p4}, ""), "update --no-progress", withEnv(nil, func(s *Scenario) {
			s.Setup = []string{"update --no-progress --with vcs3/p4lib:1.0.0"}
			s.Probe = true
		})),
		ps("source-perforce-update-v", vcsProject(req, []string{p4}, ""), "update --no-progress -v", withEnv(nil, func(s *Scenario) {
			s.Setup = []string{"update --no-progress --with vcs3/p4lib:1.0.0"}
		})),
		ps("source-perforce-status", vcsProject(req, []string{p4}, ""), "status -v", withEnv(nil, func(s *Scenario) {
			s.Setup = []string{"update --no-progress"}
			s.After = map[string]string{"vendor/vcs3/p4lib/src/Lib.php": "<?php // local change\n"}
		})),
		ps("source-perforce-remove", vcsProject(req, []string{p4}, ""), "remove vcs3/p4lib --no-progress", withEnv(nil, func(s *Scenario) {
			s.Setup = []string{"update --no-progress"}
		})),
	}
}
