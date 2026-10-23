package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"strings"
	"time"
)

// vcs2Scenarios compare the API based VCS drivers (GitHubDriver,
// GitLabDriver, GitBitbucketDriver) against the loopback forge stub
// (testdata/compare/vcs/forge-stub.php, HTTPS through a fixture CA given as
// "cafile") and the ArtifactRepository (Composer
// f83588392df43d4fd027d84d26de459946ef4e80).
func init() { extraScenarios = append(extraScenarios, vcs2Scenarios) }

// forgeToken is FORGE_TOKEN of forge-stub.php. The vcs-forge-*-auth
// scenarios check that it never shows up in any output (see the report in
// docs/validation/*-compare-vcs2*).
const forgeToken = "s3cr3t-forge-token"

const forgeCA = "/fixtures/git/_api/ca.pem"

// forgeProject has only the given repositories and trusts the fixture CA.
func forgeProject(requireJSON string, repos []string, config string) string {
	cfg := `"cafile": "` + forgeCA + `"`
	if config != "" {
		cfg += ", " + config
	}
	return "{\n    \"name\": \"test/project\",\n    \"repositories\": [\n        " + strings.Join(repos, ",\n        ") + ",\n        {\"packagist.org\": false}\n    ],\n    \"config\": {" + cfg + "},\n    \"require\": " + requireJSON + "\n}\n"
}

var fixedZipTime = time.Date(2024, 1, 2, 3, 4, 6, 0, time.UTC)

// makeZip builds a deterministic zip; a name ending in "/" is a directory.
func makeZip(entries ...[2]string) string {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e[0], Method: zip.Deflate, Modified: fixedZipTime}
		if strings.HasSuffix(e[0], "/") {
			h.Method = zip.Store
		}
		f, err := w.CreateHeader(h)
		if err != nil {
			panic(err)
		}
		_, _ = f.Write([]byte(e[1]))
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return buf.String()
}

// makeTarGz builds a deterministic tar.gz (directories end in "/").
func makeTarGz(entries ...[2]string) string {
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.ModTime = fixedZipTime
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		h := &tar.Header{Name: e[0], Mode: 0o644, Size: int64(len(e[1])), ModTime: fixedZipTime, Typeflag: tar.TypeReg, Format: tar.FormatUSTAR}
		if strings.HasSuffix(e[0], "/") {
			h.Typeflag, h.Mode, h.Size = tar.TypeDir, 0o755, 0
		}
		if err := tw.WriteHeader(h); err != nil {
			panic(err)
		}
		_, _ = tw.Write([]byte(e[1]))
	}
	_ = tw.Close()
	_ = gz.Close()
	return buf.String()
}

func artJSON(name, version string) string {
	return `{"name": "` + name + `", "version": "` + version + `", "license": "MIT", "autoload": {"psr-4": {"Art\\": "src/"}}}`
}

const artSrc = "<?php\nnamespace Art;\n\nclass Thing {}\n"

func artProject(requireJSON string, repos ...string) string {
	if len(repos) == 0 {
		repos = []string{`{"type": "artifact", "url": "artifacts/"}`}
	}
	return vcsOnlyProject(requireJSON, repos)
}

func vcs2Scenarios() []Scenario {
	// --- forge APIs ---
	forgeFx := []string{"gitlib", "forge"}
	gh := gitRepo("vcs", "https://github.com/vcs/gitlib")
	ghPriv := gitRepo("vcs", "https://github.com/vcs/gitlib-private")
	gl := gitRepo("vcs", "https://gitlab.com/vcs/gitlib")
	glPriv := gitRepo("vcs", "https://gitlab.com/vcs/gitlib-private")
	bb := gitRepo("vcs", "https://bitbucket.org/vcs/gitlib")
	bbPriv := gitRepo("vcs", "https://bitbucket.org/vcs/gitlib-private")
	req := `{"vcs/gitlib": "^1.0"}`
	auth := func(j string) map[string]string { return map[string]string{"COMPOSER_AUTH": j} }
	ghAuth := auth(`{"github-oauth": {"github.com": "` + forgeToken + `"}}`)
	glAuth := auth(`{"gitlab-token": {"gitlab.com": "` + forgeToken + `"}}`)
	glOAuth := auth(`{"gitlab-oauth": {"gitlab.com": "` + forgeToken + `"}}`)
	bbAuth := auth(`{"http-basic": {"bitbucket.org": {"username": "fixture-user", "password": "` + forgeToken + `"}}}`)
	fs := func(name, json, cmd string, extra func(*Scenario)) Scenario {
		s := Scenario{Name: name, JSON: json, GitFixtures: forgeFx, Command: cmd}
		if extra != nil {
			extra(&s)
		}
		return s
	}
	probe := func(s *Scenario) { s.Probe = true }
	withEnv := func(env map[string]string) func(*Scenario) { return func(s *Scenario) { s.Env = env } }
	setupNoInstall := func(s *Scenario) { s.Setup = []string{"update --no-progress --no-install"}; s.Probe = true }

	out := []Scenario{
		// GitHub (github.com -> api.github.com)
		fs("vcs-forge-github", forgeProject(req, []string{gh}, ""), "update --no-progress", probe),
		fs("vcs-forge-github-vv", forgeProject(req, []string{gh}, ""), "update --no-progress -vv", nil),
		fs("vcs-forge-github-type", forgeProject(`{"vcs/gitlib": "dev-main"}`, []string{gitRepo("github", "https://github.com/vcs/gitlib.git")}, ""), "update --no-progress", probe),
		fs("vcs-forge-github-1x-dev", forgeProject(`{"vcs/gitlib": "1.x-dev"}`, []string{gh}, ""), "update --no-progress", probe),
		fs("vcs-forge-github-install", forgeProject(req, []string{gh}, ""), "install --no-progress", setupNoInstall),
		fs("vcs-forge-github-missing", forgeProject(req, []string{gitRepo("vcs", "https://github.com/vcs/nope")}, ""), "update --no-progress", nil),
		fs("vcs-forge-github-private-noauth", forgeProject(req, []string{ghPriv}, ""), "update --no-progress", nil),
		fs("vcs-forge-github-private-nofallback", forgeProject(req, []string{`{"type": "github", "url": "https://github.com/vcs/gitlib-private", "allow-git-fallback": false}`}, ""), "update --no-progress", nil),
		fs("vcs-forge-github-private-auth", forgeProject(req, []string{ghPriv}, ""), "update --no-progress -vv", withEnv(ghAuth)),
		fs("vcs-forge-github-private-auth-lock", forgeProject(req, []string{ghPriv}, ""), "update --no-progress", func(s *Scenario) { s.Env = ghAuth; s.Probe = true }),
		fs("vcs-forge-github-ratelimit", forgeProject(req, []string{gitRepo("vcs", "https://github.com/vcs/gitlib-ratelimited")}, ""), "update --no-progress", nil),
		fs("vcs-forge-github-ratelimit-auth", forgeProject(req, []string{gitRepo("vcs", "https://github.com/vcs/gitlib-ratelimited")}, ""), "update --no-progress", withEnv(ghAuth)),
		fs("vcs-forge-github-enterprise", forgeProject(req, []string{gitRepo("vcs", "https://ghe.test/vcs/gitlib")}, `"github-domains": ["github.com", "ghe.test"]`), "update --no-progress -vv", probe),
		fs("vcs-forge-github-not-domain", forgeProject(req, []string{gitRepo("github", "https://ghe.test/vcs/gitlib")}, ""), "update --no-progress", nil),
		// credentials embedded in a plain git url (GitDriver): https through
		// the untrusted-by-git stub fails, dumb http through the fixture
		// server succeeds; both outputs must mask the password
		fs("vcs-forge-git-url-credentials", forgeProject(req, []string{gitRepo("vcs", "https://fixture-user:"+forgeToken+"@gitlab.test/vcs/gitlib.git")}, ""), "update --no-progress", nil),
		fs("vcs-forge-git-http-credentials", forgeProject(req, []string{gitRepo("vcs", "http://fixture-user:"+forgeToken+"@127.0.0.1:8899/git/gitlib.git")}, `"secure-http": false`), "update --no-progress -vv", probe),
		// composer repositories answering 401/403 (AuthHelper::promptAuthIfNeeded,
		// non-interactive) and the same with credentials
		fs("vcs-forge-composer-401", forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/auth401"}`}, `"secure-http": false`), "update --no-progress", nil),
		fs("vcs-forge-composer-403", forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/auth403"}`}, `"secure-http": false`), "update --no-progress", nil),
		fs("vcs-forge-composer-401-basic", forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/auth401"}`}, `"secure-http": false`), "update --no-progress", withEnv(auth(`{"http-basic": {"127.0.0.1:8900": {"username": "fixture-user", "password": "`+forgeToken+`"}}}`))),
		fs("vcs-forge-composer-403-bearer", forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/auth403"}`}, `"secure-http": false`), "update --no-progress", withEnv(auth(`{"bearer": {"127.0.0.1:8900": "`+forgeToken+`"}}`))),
		fs("vcs-forge-composer-401-wrong", forgeProject(`{"vcs/none": "*"}`, []string{`{"type": "composer", "url": "http://127.0.0.1:8900/auth401"}`}, `"secure-http": false`), "update --no-progress", withEnv(auth(`{"http-basic": {"127.0.0.1:8900": {"username": "fixture-user", "password": "wrong-password"}}}`))),
		// GitLab
		fs("vcs-forge-gitlab", forgeProject(req, []string{gl}, ""), "update --no-progress", probe),
		fs("vcs-forge-gitlab-vv", forgeProject(req, []string{gl}, ""), "update --no-progress -vv", nil),
		fs("vcs-forge-gitlab-dev", forgeProject(`{"vcs/gitlib": "dev-main"}`, []string{gitRepo("gitlab", "https://gitlab.com/vcs/gitlib")}, ""), "update --no-progress", probe),
		fs("vcs-forge-gitlab-subgroup", forgeProject(req, []string{gitRepo("vcs", "https://gitlab.com/group/sub/gitlib")}, ""), "update --no-progress -vv", probe),
		fs("vcs-forge-gitlab-http-port", forgeProject(req, []string{gitRepo("vcs", "http://127.0.0.1:8900/vcs/gitlib")}, `"secure-http": false, "gitlab-domains": ["gitlab.com", "127.0.0.1:8900"]`), "update --no-progress -vv", probe),
		fs("vcs-forge-gitlab-install", forgeProject(req, []string{gl}, ""), "install --no-progress", setupNoInstall),
		fs("vcs-forge-gitlab-private-noauth", forgeProject(req, []string{glPriv}, ""), "update --no-progress", nil),
		fs("vcs-forge-gitlab-private-token", forgeProject(req, []string{glPriv}, ""), "update --no-progress -vv", withEnv(glAuth)),
		fs("vcs-forge-gitlab-private-oauth", forgeProject(req, []string{glPriv}, ""), "update --no-progress", func(s *Scenario) { s.Env = glOAuth; s.Probe = true }),
		fs("vcs-forge-gitlab-protocol-invalid", forgeProject(req, []string{gl}, `"gitlab-protocol": "ftp"`), "update --no-progress", nil),
		fs("vcs-forge-gitlab-protocol-git", forgeProject(req, []string{gl}, `"gitlab-protocol": "git"`), "update --no-progress --no-install", nil),
		// Bitbucket
		fs("vcs-forge-bitbucket", forgeProject(req, []string{bb}, ""), "update --no-progress", probe),
		fs("vcs-forge-bitbucket-vv", forgeProject(req, []string{bb}, ""), "update --no-progress -vv", nil),
		fs("vcs-forge-bitbucket-dev", forgeProject(`{"vcs/gitlib": "dev-main"}`, []string{gitRepo("bitbucket", "https://bitbucket.org/vcs/gitlib.git")}, ""), "update --no-progress", probe),
		fs("vcs-forge-bitbucket-install", forgeProject(req, []string{bb}, ""), "install --no-progress", setupNoInstall),
		fs("vcs-forge-bitbucket-private-noauth", forgeProject(req, []string{bbPriv}, ""), "update --no-progress", nil),
		fs("vcs-forge-bitbucket-private-basic", forgeProject(req, []string{bbPriv}, ""), "update --no-progress -vv", withEnv(bbAuth)),
	}

	// --- vcs type detection of the drivers gocomposer does not implement
	// (svn, hg, fossil, perforce): neither binary exists in the container
	for _, c := range [][2]string{
		{"svn-scheme", "svn://127.0.0.1/repo/trunk"},
		{"svn-host", "http://svn.example.org/repo"},
		{"hg-ssh", "ssh://hg@bitbucket.org/acme/repo"},
		{"fossil", "https://chiselapp.com/user/acme/repository/repo"},
		{"perforce", "p4://127.0.0.1:1666/depot"},
		{"git-ssh", "git@127.0.0.1:acme/repo.git"},
	} {
		out = append(out, fs("vcs-detect-"+c[0], forgeProject(req, []string{gitRepo("vcs", c[1])}, `"secure-http": false`), "update --no-progress", nil))
	}

	// --- artifact repository ---
	files := func(m map[string]string) map[string]string { return m }
	zip100 := makeZip([2]string{"composer.json", artJSON("art/thing", "1.0.0")}, [2]string{"src/Thing.php", artSrc})
	zip110 := makeZip([2]string{"thing-1.1.0/", ""}, [2]string{"thing-1.1.0/composer.json", artJSON("art/thing", "1.1.0")}, [2]string{"thing-1.1.0/src/Thing.php", artSrc})
	tgz120 := makeTarGz([2]string{"package/", ""}, [2]string{"package/composer.json", artJSON("art/thing", "1.2.0")}, [2]string{"package/src/Thing.php", artSrc})
	tar130 := func() string {
		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		for _, e := range [][2]string{{"composer.json", artJSON("art/thing", "1.3.0")}, {"src/Thing.php", artSrc}} {
			_ = tw.WriteHeader(&tar.Header{Name: e[0], Mode: 0o644, Size: int64(len(e[1])), ModTime: fixedZipTime, Typeflag: tar.TypeReg, Format: tar.FormatUSTAR})
			_, _ = tw.Write([]byte(e[1]))
		}
		_ = tw.Close()
		return buf.String()
	}()
	other := makeZip([2]string{"composer.json", `{"name": "art/other", "version": "2.0.0", "require": {"art/thing": "^1.0"}}`})
	badJSON := makeZip([2]string{"composer.json", `{"name": "art/bad", "version": "1.0.0",}`})
	noJSON := makeZip([2]string{"README.md", "nothing\n"})
	deep := makeZip([2]string{"a/", ""}, [2]string{"a/b/", ""}, [2]string{"a/b/composer.json", artJSON("art/deep", "1.0.0")})
	twoTop := makeZip([2]string{"one/composer.json", artJSON("art/two", "1.0.0")}, [2]string{"two/README.md", "x\n"})
	noVersion := makeZip([2]string{"composer.json", `{"name": "art/nover"}`})
	badVersion := makeZip([2]string{"composer.json", `{"name": "art/badver", "version": "not a version"}`})
	empty := makeZip()
	macosx := makeZip([2]string{"__MACOSX/", ""}, [2]string{"__MACOSX/._composer.json", "junk"}, [2]string{"pkg/", ""}, [2]string{"pkg/composer.json", artJSON("art/mac", "1.0.0")})
	art := func(name, require, cmd string, f map[string]string, extra func(*Scenario)) Scenario {
		s := Scenario{Name: name, JSON: artProject(require), Files: f, Command: cmd, Image: zipImage}
		if extra != nil {
			extra(&s)
		}
		return s
	}
	two := files(map[string]string{"artifacts/thing-1.0.0.zip": zip100, "artifacts/sub/thing-1.1.0.zip": zip110})
	out = append(out,
		art("vcs-artifact-basic", `{"art/thing": "^1.0"}`, "update --no-progress", two, probe),
		art("vcs-artifact-pin-old", `{"art/thing": "1.0.0"}`, "update --no-progress", two, probe),
		art("vcs-artifact-verbose", `{"art/thing": "^1.0"}`, "update --no-progress -v", files(map[string]string{"artifacts/thing-1.1.0.zip": zip110}), nil),
		art("vcs-artifact-tgz", `{"art/thing": "^1.2"}`, "update --no-progress -v", files(map[string]string{"artifacts/thing-1.2.0.tgz": tgz120}), probe),
		art("vcs-artifact-tar", `{"art/thing": "^1.3"}`, "update --no-progress -v", files(map[string]string{"artifacts/thing-1.3.0.tar": tar130}), probe),
		art("vcs-artifact-dep", `{"art/other": "^2.0"}`, "update --no-progress", files(map[string]string{"artifacts/other.zip": other, "artifacts/thing-1.0.0.zip": zip100}), probe),
		art("vcs-artifact-install-lock", `{"art/thing": "^1.0"}`, "install --no-progress", two, setupNoInstall),
		art("vcs-artifact-invalid-json", `{"art/bad": "*"}`, "update --no-progress", files(map[string]string{"artifacts/bad.zip": badJSON}), nil),
		art("vcs-artifact-no-composer-json", `{"art/thing": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/none.zip": noJSON}), nil),
		art("vcs-artifact-deep-composer-json", `{"art/deep": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/deep.zip": deep}), nil),
		art("vcs-artifact-two-top-dirs", `{"art/two": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/two.zip": twoTop}), nil),
		art("vcs-artifact-macosx", `{"art/mac": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/mac.zip": macosx}), probe),
		art("vcs-artifact-no-version", `{"art/nover": "*"}`, "update --no-progress", files(map[string]string{"artifacts/nover.zip": noVersion}), nil),
		art("vcs-artifact-bad-version", `{"art/badver": "*"}`, "update --no-progress", files(map[string]string{"artifacts/badver.zip": badVersion}), nil),
		art("vcs-artifact-empty-zip", `{"art/thing": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/empty.zip": empty}), nil),
		art("vcs-artifact-corrupt-zip", `{"art/thing": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/corrupt.zip": "this is not a zip archive\n"}), nil),
		art("vcs-artifact-corrupt-tgz", `{"art/thing": "*"}`, "update --no-progress -v", files(map[string]string{"artifacts/corrupt.tgz": "this is not a tarball\n"}), nil),
		art("vcs-artifact-other-files", `{"art/thing": "^1.0"}`, "update --no-progress -v", files(map[string]string{"artifacts/thing-1.0.0.zip": zip100, "artifacts/notes.txt": "ignored\n", "artifacts/readme.md": "ignored\n"}), nil),
		art("vcs-artifact-missing-dir", `{"art/thing": "*"}`, "update --no-progress", files(map[string]string{"src/.keep": ""}), nil),
		Scenario{Name: "vcs-artifact-absolute-url", JSON: artProject(`{"art/thing": "^1.0"}`, `{"type": "artifact", "url": "/work/project/artifacts"}`), Files: two, Command: "update --no-progress", Image: zipImage, Probe: true},
	)
	return out
}
