package main

import "strings"

func init() { extraScenarios = append(extraScenarios, cmds2Scenarios) }

// cmds2Scenarios covers the init, bump (and the bump-after-update config of
// update), archive and global commands (InitCommand, BumpCommand,
// ArchiveCommand, GlobalCommand of Composer
// f83588392df43d4fd027d84d26de459946ef4e80).
func cmds2Scenarios() []Scenario {
	var out []Scenario
	out = append(out, bumpScenarios()...)
	out = append(out, initScenarios()...)
	out = append(out, globalScenarios()...)
	out = append(out, archiveScenarios()...)
	return out
}

func bumpScenarios() []Scenario {
	files := map[string]string{"src/Kernel.php": appSrc}
	req := `{"php": ">=7.4", "acme/log": "^1.0", "acme/http": "~1.0", "acme/log-impl": "^1.0.0", "acme/util": "*", "acme/cli": "1.*"}`
	dev := `"require-dev": {"acme/devtools": ">=1.0"}`
	proj := project(req, dev)
	library := strings.Replace(proj, `"type": "project",`, `"type": "library",`, 1)
	untyped := strings.Replace(proj, `"type": "project",`, ``, 1)
	upd := []string{"update --no-progress"}
	noLock := withConfig(proj, `"lock": false`)
	edge := project(`{"acme/edge": "^1.0@dev", "acme/log": "^1.0"}`, `"minimum-stability": "dev"`)
	cfgBump := func(v string) string { return withConfig(proj, `"bump-after-update": `+v) }
	type b struct {
		name, json, cmd string
		setup           []string
		after           map[string]string
		known           string
	}
	list := []b{
		{"bump-basic", proj, "bump", upd, nil, ""},
		{"bump-dry-run", proj, "bump --dry-run", upd, nil, ""},
		{"bump-dev-only", proj, "bump --dev-only", upd, nil, ""},
		{"bump-no-dev-only", proj, "bump -R", upd, nil, ""},
		{"bump-library", library, "bump", upd, nil, ""},
		{"bump-library-dev-only", library, "bump -D", upd, nil, ""},
		{"bump-untyped", untyped, "bump", upd, nil, ""},
		{"bump-filter", proj, "bump acme/log ACME/HTTP:^1", upd, nil, ""},
		{"bump-filter-wildcard", proj, "bump 'acme/h*'", upd, nil, ""},
		{"bump-nothing", proj, "bump", []string{"update --no-progress", "bump"}, nil, ""},
		{"bump-nothing-dry-run", proj, "bump --dry-run", []string{"update --no-progress", "bump"}, nil, ""},
		{"bump-lock-outdated", proj, "bump", upd, map[string]string{"composer.json": project(`{"acme/log": "^1.0", "acme/util": "^1.0"}`, "")}, ""},
		{"bump-lock-disabled", noLock, "bump", upd, nil, "Composer dies with an internal LogicException (BumpCommand reads the locked repository when the lock config is disabled) which Symfony renders with its class name and stack trace; gocomposer reports the same message with the same exit code 1."},
		{"bump-installed-only", noLock, "bump", upd, map[string]string{"composer.json": proj}, ""},
		{"bump-not-installed", proj, "bump", nil, nil, ""},
		{"bump-branch-alias", edge, "bump", upd, nil, ""},
		{"bump-no-json", "", "bump", nil, nil, ""},
		{"bump-after-update-flag", proj, "update --no-progress --bump-after-update", nil, nil, ""},
		{"bump-after-update-dev", untyped, "update --no-progress --bump-after-update=dev", nil, nil, ""},
		{"bump-after-update-no-dev", proj, "update --no-progress --bump-after-update=no-dev", nil, nil, ""},
		{"bump-after-update-dry-run", proj, "update --no-progress --dry-run --bump-after-update", nil, nil, ""},
		{"bump-after-update-config", cfgBump("true"), "update --no-progress", nil, nil, ""},
		{"bump-after-update-config-dev", cfgBump(`"dev"`), "update --no-progress", nil, nil, ""},
		{"bump-after-update-untyped", strings.Replace(cfgBump("true"), `"type": "project",`, ``, 1), "update --no-progress", nil, nil, ""},
		{"bump-after-update-partial", proj, "update acme/log --no-progress --bump-after-update", []string{"update --no-progress --with acme/log:1.0.0 --with acme/http:1.0.0"}, nil, ""},
		{"bump-after-update-nothing", proj, "update --no-progress --bump-after-update", upd, nil, ""},
		{"bump-after-update-lock", proj, "update --lock --no-progress --bump-after-update", upd, nil, ""},
	}
	var out []Scenario
	for _, s := range list {
		sc := Scenario{Name: s.name, JSON: s.json, Files: files, Command: s.cmd, Setup: s.setup, After: s.after, KnownDiff: s.known}
		if s.json == "" {
			sc.NoJSON = true
		}
		out = append(out, sc)
	}
	return out
}

func initScenarios() []Scenario {
	gitEnv := map[string]string{"GIT_CONFIG_GLOBAL": "/work/gitconfig", "GIT_CONFIG_NOSYSTEM": "1"}
	gitFiles := func(cfg string, extra map[string]string) map[string]string {
		m := map[string]string{"../gitconfig": cfg}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	noGit := gitFiles("", nil)
	// a fixed autoloader suffix (there is no lock file to derive it from)
	suffixHome := gitFiles("", map[string]string{"../ghome/config.json": `{"config": {"autoloader-suffix": "InitTest"}}`})
	suffixEnv := map[string]string{"COMPOSER_HOME": "/work/ghome"}
	withGit := gitFiles("[user]\n\tname = Git User\n\temail = git@example.com\n[github]\n\tuser = GitHubber\n", nil)
	full := `init --name acme/demo --description "Demo package" --author "Jane Doe <jane@example.com>" --type library --homepage https://example.com --license MIT --stability dev`
	env := func(extra map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range gitEnv {
			m[k] = v
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	type i struct {
		name, cmd string
		files     map[string]string
		env       map[string]string
		json      string
		probe     bool
	}
	existing := "{\n    \"name\": \"old/project\",\n    \"require\": {\"acme/log\": \"^1.0\"}\n}\n"
	list := []i{
		{"init-full", full, noGit, nil, "", false},
		{"init-full-reordered", `init --license MIT --stability dev --homepage https://example.com --type project --author "Jane Doe" --description "" --name acme/demo`, noGit, nil, "", false},
		{"init-require", `init --name acme/demo --require acme/log:^1.0 --require "acme/util 1.*" --require-dev "acme/devtools=*" --require acme/log:^1.1`, noGit, nil, "", false},
		{"init-require-pair-args", `init --name acme/demo --require acme/log --require ^1.0`, noGit, nil, "", false},
		{"init-require-no-version", `init --name acme/demo --require acme/log`, noGit, nil, "", false},
		{"init-require-dev-no-version", `init --name acme/demo --require-dev acme/log`, noGit, nil, "", false},
		{"init-autoload", `init --name acme/my-demo_pkg.x --autoload src/`, suffixHome, suffixEnv, "", true},
		{"init-autoload-nested", `init --name acme/demo --autoload lib/Sub/`, suffixHome, suffixEnv, "", true},
		{"init-autoload-with-deps", `init --name acme/demo --autoload lib/ --require acme/log:^1.0`, noGit, nil, "", false},
		{"init-autoload-dev-deps", `init --name acme/demo --autoload lib/ --require-dev acme/log:^1.0`, noGit, nil, "", false},
		{"init-git-defaults", `init`, withGit, nil, "", false},
		{"init-git-defaults-description", `init --description "From git"`, withGit, nil, "", false},
		{"init-env-defaults", `init`, noGit, map[string]string{"COMPOSER_DEFAULT_VENDOR": "My_Vendor", "COMPOSER_DEFAULT_AUTHOR": "Env Author", "COMPOSER_DEFAULT_EMAIL": "env@example.com"}, "", false},
		{"init-env-overrides-git", `init`, withGit, map[string]string{"COMPOSER_DEFAULT_VENDOR": "EnvVendor", "COMPOSER_DEFAULT_EMAIL": "env@example.com"}, "", false},
		{"init-user-env", `init`, noGit, map[string]string{"USER": "SomeUser"}, "", false},
		{"init-no-defaults", `init`, noGit, nil, "", false},
		{"init-camel-dir", `init -d MyCool__Project.Name`, gitFiles("[github]\n\tuser = ACMECorp\n", map[string]string{"MyCool__Project.Name/.keep": ""}), nil, "", false},
		{"init-invalid-name", `init --name Acme/Demo`, noGit, nil, "", false},
		{"init-invalid-name-noslash", `init --name demo`, noGit, nil, "", false},
		{"init-invalid-author", `init --name acme/demo --author "Jane <"`, noGit, nil, "", false},
		{"init-invalid-email", `init --name acme/demo --author "Jane Doe <not-an-email>"`, noGit, nil, "", false},
		{"init-empty-author", `init --name acme/demo --author ""`, noGit, nil, "", false},
		{"init-author-name-only", `init --name acme/demo --author "Jöhn O'Brien (dev)"`, noGit, nil, "", false},
		{"init-existing-json", `init --name acme/demo`, noGit, nil, existing, false},
		{"init-existing-git", `init --name acme/demo --autoload src/`, gitFiles("", map[string]string{".git/HEAD": "ref: refs/heads/main\n", ".git/objects/.keep": "", ".git/refs/heads/.keep": "", ".git/config": "[core]\n\trepositoryformatversion = 0\n\tbare = false\n", ".gitignore": "/node_modules\n", "../ghome/config.json": `{"config": {"autoloader-suffix": "InitTest"}}`}), suffixEnv, "", false},
		{"init-repository", `init --name acme/demo --repository http://127.0.0.1:8899 --repository '{"type":"path","url":"../pkgs","options":{"symlink":false}}' --repository '{"packagist.org":false}'`, noGit, nil, "", false},
		{"init-repository-invalid", `init --name acme/demo --repository foo://user:secret@example.com/repo`, noGit, nil, "", false},
		{"init-repository-bad-json", `init --name acme/demo --repository '{"type":'`, noGit, nil, "", false},
		{"init-repository-composer-file", `init --name acme/demo --repository repo.json`, gitFiles("", map[string]string{"repo.json": `{"packages": {"acme/x": {}}}`}), nil, "", false},
		{"init-repository-filesystem-file", `init --name acme/demo --repository other.json`, gitFiles("", map[string]string{"other.json": `{"foo": 1}`}), nil, "", false},
		{"init-repository-missing-file", `init --name acme/demo --repository missing.json`, noGit, nil, "", false},
		{"init-bad-stability", `init --name acme/demo --stability foo`, noGit, nil, "", false},
		{"init-bad-type-schema", `init --name acme/demo --homepage "not a url"`, noGit, nil, "", false},
		{"init-bad-license", `init --name acme/demo --license "Not A License"`, noGit, nil, "", false},
		{"init-composer-env", `init --name acme/demo --require acme/log:^1.0`, noGit, map[string]string{"COMPOSER": "custom.json"}, "", false},
		{"init-composer-env-subdir", `init --name acme/demo`, noGit, map[string]string{"COMPOSER": "sub/dir/composer.json"}, "", false},
		{"init-unicode", `init --name acme/demo --description "Ünïcødé / slashes \\ and \"quotes\""`, noGit, nil, "", false},
	}
	var out []Scenario
	for _, s := range list {
		sc := Scenario{Name: s.name, JSON: s.json, Files: s.files, Command: s.cmd, Env: env(s.env), Probe: s.probe, NoJSON: s.json == "",
			Snapshot: []string{"custom.json", "sub/dir/composer.json", ".gitignore", ".gitattributes", "MyCool__Project.Name/composer.json"}}
		out = append(out, sc)
	}
	return out
}

func globalScenarios() []Scenario {
	ghome := map[string]string{"COMPOSER_HOME": "/work/ghome"}
	globalJSON := "{\n    \"repositories\": [\n        {\"type\": \"composer\", \"url\": \"http://127.0.0.1:8899\"},\n        {\"packagist.org\": false}\n    ],\n    \"config\": {\"secure-http\": false}\n}\n"
	withHome := map[string]string{"../ghome/composer.json": globalJSON}
	installed := map[string]string{"../ghome/composer.json": strings.Replace(globalJSON, `"config"`, `"require": {"acme/log": "^1.0"},
    "config"`, 1)}
	snap := []string{"../ghome/composer.json", "../ghome/composer.lock", "../ghome/vendor/composer/installed.json", "../ghome/vendor/composer/installed.php", "../newhome/composer.json", "other.json"}
	env := func(extra map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range ghome {
			m[k] = v
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	type g struct {
		name, cmd string
		files     map[string]string
		setup     []string
		env       map[string]string
	}
	list := []g{
		{"global-require", "global require acme/log:^1.0 --no-progress", withHome, nil, nil},
		{"global-update", "global update --no-progress", installed, nil, nil},
		{"global-show", "global show", installed, []string{"global update --no-progress"}, nil},
		{"global-show-abbrev", "glo show --name-only", installed, []string{"global update --no-progress"}, nil},
		{"global-options-first", "-n --no-ansi global show -N", installed, []string{"global update --no-progress"}, nil},
		{"global-remove", "global remove acme/log --no-progress", installed, []string{"global update --no-progress"}, nil},
		{"global-no-args", "global", withHome, nil, nil},
		{"global-unknown-option-only", "global --foo", withHome, nil, nil},
		{"global-unknown-command", "global frobnicate", withHome, nil, nil},
		{"global-config-quoted", `global config extra.foo "a b 'c'"`, withHome, nil, nil},
		{"global-config-equals", `global config --json extra.list '["x y", 2]'`, withHome, nil, nil},
		{"global-new-home", "global config extra.created yes", nil, nil, map[string]string{"COMPOSER_HOME": "/work/newhome"}},
		{"global-composer-env", "global config extra.k v", map[string]string{"../ghome/composer.json": globalJSON, "other.json": "{}\n"}, nil, map[string]string{"COMPOSER": "other.json"}},
		{"global-validate", "global validate", installed, nil, nil},
		{"global-dump-autoload", "global dump-autoload", installed, []string{"global update --no-progress"}, nil},
	}
	var out []Scenario
	for _, s := range list {
		out = append(out, Scenario{Name: s.name, JSON: project(`{"acme/util": "^1.0"}`, ""), Files: s.files, Setup: s.setup, Command: s.cmd, Env: env(s.env), Snapshot: snap})
	}
	return out
}

// listArchivesPHP lists every archive below the given directories with all
// header fields except timestamps (and checksums derived from them) plus
// the sha1 of each entry, then deletes the archives so that the project
// trees compare equal.
const listArchivesPHP = `<?php
function hexs($s) { return bin2hex($s); }
function tarList($d) {
    $off = 0; $n = strlen($d);
    while ($off + 512 <= $n) {
        $h = substr($d, $off, 512);
        if ($h === str_repeat("\0", 512)) { echo "  zero-block @", $off, "\n"; $off += 512; continue; }
        $f = function ($a, $l) use ($h) { return substr($h, $a, $l); };
        $size = octdec(trim($f(124, 12), "\0 "));
        $type = $f(156, 1);
        echo "  entry name=", json_encode(rtrim($f(0, 100), "\0")), " prefix=", json_encode(rtrim($f(345, 155), "\0")),
            " mode=", hexs($f(100, 8)), " uid=", hexs($f(108, 8)), " gid=", hexs($f(116, 8)), " size=", hexs($f(124, 12)),
            " type=", hexs($type), " link=", json_encode(rtrim($f(157, 100), "\0")), " magic=", hexs($f(257, 6)), " ver=", hexs($f(263, 2)),
            " uname=", hexs(rtrim($f(265, 32), "\0")), " gname=", hexs(rtrim($f(297, 32), "\0")), " dev=", hexs($f(329, 16)), " pad=", hexs($f(500, 12)),
            " mtime-len=", strlen(rtrim($f(136, 12), "\0")), "\n";
        $data = substr($d, $off + 512, $size);
        echo "    sha1=", sha1($data), "\n";
        $off += 512 + (int) (ceil($size / 512) * 512);
    }
    echo "  total=", $n, "\n";
}
function zipList($d) {
    $off = 0; $n = strlen($d);
    while ($off < $n) {
        $sig = unpack('V', substr($d, $off, 4))[1];
        if ($sig === 0x04034b50) {
            $h = unpack('vver/vflag/vmeth/vmt/vmd/Vcrc/Vcs/Vus/vnl/vel', substr($d, $off + 4, 26));
            $name = substr($d, $off + 30, $h['nl']); $extra = substr($d, $off + 30 + $h['nl'], $h['el']);
            echo "  local ver=", $h['ver'], " flag=", $h['flag'], " meth=", $h['meth'], " crc=", dechex($h['crc']), " cs=", $h['cs'], " us=", $h['us'], " name=", json_encode($name), " extra=", hexs($extra), "\n";
            echo "    sha1=", sha1(substr($d, $off + 30 + $h['nl'] + $h['el'], $h['cs'])), "\n";
            $off += 30 + $h['nl'] + $h['el'] + $h['cs'];
        } elseif ($sig === 0x02014b50) {
            $h = unpack('vvm/vvn/vflag/vmeth/vmt/vmd/Vcrc/Vcs/Vus/vnl/vel/vcl/vdn/via/Vea/Vlo', substr($d, $off + 4, 42));
            $name = substr($d, $off + 46, $h['nl']); $extra = substr($d, $off + 46 + $h['nl'], $h['el']);
            echo "  central vm=", $h['vm'], " vn=", $h['vn'], " flag=", $h['flag'], " meth=", $h['meth'], " crc=", dechex($h['crc']), " cs=", $h['cs'], " us=", $h['us'],
                " name=", json_encode($name), " extra=", hexs($extra), " comment=", hexs(substr($d, $off + 46 + $h['nl'] + $h['el'], $h['cl'])), " disk=", $h['dn'], " ia=", $h['ia'], " ea=", dechex($h['ea']), " lo=", $h['lo'], "\n";
            $off += 46 + $h['nl'] + $h['el'] + $h['cl'];
        } elseif ($sig === 0x06054b50) {
            echo "  eocd ", hexs(substr($d, $off)), "\n";
            break;
        } else {
            echo "  unknown signature ", dechex($sig), " @", $off, "\n";
            break;
        }
    }
    echo "  total=", $n, "\n";
}
$found = [];
foreach (array_slice($argv, 1) as $dir) {
    if (!is_dir($dir)) continue;
    $it = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($dir, FilesystemIterator::SKIP_DOTS));
    foreach ($it as $f) {
        $p = $f->getPathname();
        if (strpos($p, '/vendor/') !== false) continue;
        if (preg_match('{\.(zip|tar|tar\.gz|tar\.bz2)$}', $p)) $found[] = $p;
    }
}
sort($found);
foreach ($found as $p) {
    $d = file_get_contents($p);
    echo "archive ", $p, " perms=", decoct(fileperms($p) & 0777), "\n";
    if (substr($p, -4) === '.zip') {
        zipList($d);
    } elseif (substr($p, -7) === '.tar.gz') {
        if (substr($d, 0, 2) === "\x1f\x8b") {
            echo "  gzip header ", hexs(substr($d, 0, 10)), "\n";
            tarList(gzdecode($d));
        } else {
            echo "  zlib header ", hexs(substr($d, 0, 2)), "\n";
            tarList(gzuncompress($d));
        }
    } elseif (substr($p, -8) === '.tar.bz2') {
        echo "  bzip2 ", hexs(substr($d, 0, 4)), "\n";
    } else {
        tarList($d);
    }
    unlink($p);
}
`

func archiveScenarios() []Scenario {
	// the archive command's own output goes to the compared files; the
	// archives are then listed into the probe output (and deleted), and the
	// rest of the harness line is commented out
	run := func(args string) string {
		return args + ` >/work/out/stdout 2>/work/out/stderr; echo $? >/work/out/exit; php /work/tools/list-archives.php /work/project /work/arc >/work/out/probe 2>&1 #`
	}
	src := map[string]string{
		"../tools/list-archives.php": listArchivesPHP,
		"src/Kernel.php":             appSrc,
		"src/Sub/Helper.php":         "<?php\nnamespace App\\Sub;\n\nclass Helper {}\n",
		"src/Sub/notes.txt":          "notes\n",
		"docs/index.md":              "# docs\n",
		"docs/keep.md":               "keep\n",
		"README.md":                  "readme\n",
		"tests/FooTest.php":          "<?php\n",
		".gitattributes":             "/tests export-ignore\n*.txt export-ignore\n/docs/keep.md -export-ignore\n# comment\n\n/missing export-ignore\n",
		".git/HEAD":                  "ref: refs/heads/main\n",
		".git/objects/.keep":         "",
		".git/refs/heads/.keep":      "",
		".git/config":                "[core]\n\trepositoryformatversion = 0\n\tbare = false\n",
		".hidden":                    "hidden\n",
		"empty-parent/empty/.keep":   "",
		"bin/tool":                   "#!/bin/sh\necho tool\n",
		"deep/" + strings.Repeat("d", 60) + "/" + strings.Repeat("e", 60) + "/file.txt": "long\n",
	}
	archiveJSON := func(extra string) string {
		return project(`{"acme/log": "^1.0"}`, `"archive": {`+extra+`}`)
	}
	base := archiveJSON(`"exclude": ["/docs", "!/docs/keep.md", "*.md", "!README.md", "/empty-parent/empty/.keep"]`)
	named := archiveJSON(`"name": "my-archive", "exclude": ["/tests"]`)
	upd := []string{"update --no-progress"}
	type a struct {
		name, json, cmd string
		setup           []string
		files           map[string]string
	}
	list := []a{
		{"archive-root-tar", base, "archive", nil, nil},
		{"archive-root-zip", base, "archive --format=zip", nil, nil},
		{"archive-root-tar-gz", base, "archive -f tar.gz", nil, nil},
		{"archive-root-tar-bz2", base, "archive --format tar.bz2", nil, nil},
		{"archive-root-bad-format", base, "archive --format=rar", nil, nil},
		{"archive-root-empty-format", base, "archive --format=", nil, nil},
		{"archive-root-dir", base, "archive --format=zip --dir=/work/arc", nil, nil},
		{"archive-root-dir-relative", base, "archive --dir=build/out", nil, nil},
		{"archive-root-file", base, "archive --format=zip --file=custom --dir=/work/arc", nil, nil},
		{"archive-root-file-subdir", base, "archive --file=sub/name", nil, nil},
		{"archive-root-ignore-filters", base, "archive --format=zip --ignore-filters --dir=/work/arc", nil, nil},
		{"archive-root-named", named, "archive --dir=/work/arc", nil, nil},
		{"archive-root-installed", base, "archive --format=zip", upd, nil},
		{"archive-root-existing-archives", base, "archive", nil, map[string]string{"test-project-old.tar": "old", "test-project.zip": "old", "other.zip": "keep"}},
		{"archive-root-config-format", withConfig(base, `"archive-format": "zip", "archive-dir": "dist"`), "archive", nil, nil},
		{"archive-root-dir-is-file", base, "archive --dir=README.md", nil, nil},
		{"archive-root-scripts", strings.Replace(base, `"require":`, `"scripts": {"pre-archive-cmd": "echo pre-archive", "post-archive-cmd": "echo post-archive"},
    "require":`, 1), "archive --dir=/work/arc", nil, nil},
		{"archive-package", base, "archive acme/log 1.0.0 --dir=/work/arc", nil, nil},
		{"archive-package-installed", base, "archive acme/log --format=zip --dir=/work/arc", upd, nil},
		{"archive-package-constraint", base, "archive acme/util '^1.0' --dir=/work/arc --format=tar.gz", nil, nil},
		{"archive-package-stability", base, "archive acme/beta '^1.1@beta' --dir=/work/arc", nil, nil},
		{"archive-package-stability-filtered", base, "archive acme/beta '^1.1' --dir=/work/arc", nil, nil},
		{"archive-package-missing", base, "archive acme/nope --dir=/work/arc", nil, nil},
		{"archive-package-file", base, "archive acme/edge --file=edge --dir=/work/arc", nil, nil},
		{"archive-package-case", base, "archive ACME/Util 1.1.0 --dir=/work/arc", nil, nil},
		{"archive-no-json", "", "archive", nil, nil},
		{"archive-empty-tar", archiveJSON(`"exclude": ["*", ".*"]`), "archive --dir=/work/arc", nil, nil},
		{"archive-empty-zip", archiveJSON(`"exclude": ["*", ".*"]`), "archive --format=zip --dir=/work/arc", nil, nil},
		{"archive-empty-tar-gz", archiveJSON(`"exclude": ["*", ".*"]`), "archive --format=tar.gz --dir=/work/arc", nil, nil},
		{"archive-only-empty-dirs", archiveJSON(`"exclude": ["*", ".*", "!/empty-parent", "!/empty-parent/empty", "/empty-parent/empty/.keep"]`), "archive --format=zip --dir=/work/arc", nil, nil},
		{"archive-globstar-attributes", base, "archive --format=zip --dir=/work/arc", nil, map[string]string{".gitattributes": "/deep/** export-ignore\n/deep/dd*/** -export-ignore\nsrc/{Sub,Other}/*.php export-ignore\n?EADME.md export-ignore\n\\#x export-ignore\nbad[ export-ignore\n.git* export-ignore\n"}},
		{"archive-globstar-valid", base, "archive --format=zip --dir=/work/arc", nil, map[string]string{".gitattributes": "/deep/** export-ignore\n/deep/dd*/** -export-ignore\nsrc/{Sub,Other}/*.php export-ignore\n?EADME.md export-ignore\n\\#x export-ignore\n.git* export-ignore\n.hid?en -export-ignore\n0\n  # indented comment\nsrc/Kernel.php   export-ignore extra\n"}},
		{"archive-vendor-dir", withConfig(base, `"vendor-dir": "lib/vendor"`), "archive --format=zip --dir=/work/arc", upd, nil},
	}
	var out []Scenario
	for _, s := range list {
		files := map[string]string{}
		for k, v := range src {
			files[k] = v
		}
		for k, v := range s.files {
			files[k] = v
		}
		out = append(out, Scenario{Name: s.name, JSON: s.json, NoJSON: s.json == "", Files: files, Setup: s.setup, Command: run(s.cmd)})
	}
	return out
}
