package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Repository/Vcs/GitHubDriverTest.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
//
// Not ported: testPrivateRepository (interactive token prompt, which
// gocomposer reports as unsupported instead of asking).

import (
	"encoding/base64"
	"reflect"
	"testing"

	"gocomposer/internal/phpjson"
)

const (
	ghRepoURL    = "http://github.com/composer/packagist"
	ghRepoAPIURL = "https://api.github.com/repos/composer/packagist"
)

func b64(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }

func newGitHub(t *testing.T, interactive bool, expects []httpExpect, proc *procMock, cfg map[string]phpjson.Value) (*GitHubDriver, *httpMock) {
	t.Helper()
	env, hm, _ := testEnv(t, interactive, cfg, expects, proc)
	d := NewGitHubDriver(repoConfig(ghRepoURL), env).(*GitHubDriver)
	if err := d.Initialize(); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	return d, hm
}

func assertDistSource(t *testing.T, d Driver, sha, wantSourceURL string) {
	t.Helper()
	dist, err := d.Dist(sha)
	if err != nil || dist == nil {
		t.Fatalf("dist: %v %v", dist, err)
	}
	if dist.Str("type") != "zip" || dist.Str("url") != "https://api.github.com/repos/composer/packagist/zipball/SOMESHA" || dist.Str("reference") != sha {
		t.Errorf("dist = %v", dist)
	}
	src := d.Source(sha)
	if src.Str("type") != "git" || src.Str("url") != wantSourceURL || src.Str("reference") != sha {
		t.Errorf("source = %v", src)
	}
}

func TestGitHubPublicRepository(t *testing.T) {
	d, hm := newGitHub(t, true, []httpExpect{
		{url: ghRepoAPIURL, body: `{"master_branch": "test_master", "owner": {"login": "composer"}, "name": "packagist"}`},
	}, nil, nil)
	d.tags, d.tagsOK = []Ref{{"v0.0.0", "SOMESHA"}}, true
	if root, _ := d.RootIdentifier(); root != "test_master" {
		t.Errorf("root = %q", root)
	}
	assertDistSource(t, d, "SOMESHA", "https://github.com/composer/packagist.git")
	hm.assertComplete()
}

func TestGitHubPublicRepository2(t *testing.T) {
	repoURL := ghRepoURL
	d, hm := newGitHub(t, true, []httpExpect{
		{url: ghRepoAPIURL, body: `{"master_branch": "test_master", "owner": {"login": "composer"}, "name": "packagist"}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/composer.json?ref=feature%2F3.2-foo", body: `{"encoding":"base64","content":"` + b64(`{"support": {"source": "`+repoURL+`" }}`) + `"}`},
		{url: "https://api.github.com/repos/composer/packagist/commits/feature%2F3.2-foo", body: `{"commit": {"committer":{ "date": "2012-09-10"}}}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/.github/FUNDING.yml", body: `{"encoding": "base64", "content": "` + b64("custom: https://example.com") + `"}`},
	}, nil, nil)
	d.tags, d.tagsOK = []Ref{{"feature/3.2-foo", "SOMESHA"}}, true
	if root, _ := d.RootIdentifier(); root != "test_master" {
		t.Errorf("root = %q", root)
	}
	assertDistSource(t, d, "SOMESHA", "https://github.com/composer/packagist.git")
	data, err := d.ComposerInformation("feature/3.2-foo")
	if err != nil || data == nil {
		t.Fatalf("composer information: %v %v", data, err)
	}
	if data.Isset("abandoned") {
		t.Errorf("abandoned set: %v", data)
	}
	hm.assertComplete()
}

func TestGitHubInvalidSupportData(t *testing.T) {
	d, hm := newGitHub(t, true, []httpExpect{
		{url: ghRepoAPIURL, body: `{"master_branch": "test_master", "owner": {"login": "composer"}, "name": "packagist"}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/composer.json?ref=feature%2F3.2-foo", body: `{"encoding":"base64","content":"` + b64(`{"support": "`+ghRepoURL+`" }`) + `"}`},
		{url: "https://api.github.com/repos/composer/packagist/commits/feature%2F3.2-foo", body: `{"commit": {"committer":{ "date": "2012-09-10"}}}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/.github/FUNDING.yml", body: `{"encoding": "base64", "content": "` + b64("custom: https://example.com") + `"}`},
	}, nil, nil)
	d.tags, d.tagsOK = []Ref{{"feature/3.2-foo", "SOMESHA"}}, true
	d.branches, d.branchesOK = []Ref{{"test_master", "SOMESHA"}}, true
	data, err := d.ComposerInformation("feature/3.2-foo")
	if err != nil || data == nil {
		t.Fatalf("composer information: %v %v", data, err)
	}
	if v, _ := mapGet(data, "support", "source"); strOf(v) != "https://github.com/composer/packagist/tree/feature/3.2-foo" {
		t.Errorf("support.source = %v", v)
	}
	hm.assertComplete()
}

func TestGitHubFundingFormat(t *testing.T) {
	type f struct{ typ, url string }
	all := "community_bridge: project-name\ngithub: [userA, userB]\nissuehunt: userName\nko_fi: userName\nliberapay: userName\nopen_collective: userName\npatreon: userName\ntidelift: Platform/Package\npolar: userName\nbuy_me_a_coffee: userName\nthanks_dev: u/gh/userName\notechie: userName"
	cases := []struct {
		name, funding string
		want          []f
	}{
		{"All named platforms", all, []f{
			{"community_bridge", "https://funding.communitybridge.org/projects/project-name"},
			{"github", "https://github.com/userA"},
			{"github", "https://github.com/userB"},
			{"issuehunt", "https://issuehunt.io/r/userName"},
			{"ko_fi", "https://ko-fi.com/userName"},
			{"liberapay", "https://liberapay.com/userName"},
			{"open_collective", "https://opencollective.com/userName"},
			{"patreon", "https://www.patreon.com/userName"},
			{"tidelift", "https://tidelift.com/funding/github/Platform/Package"},
			{"polar", "https://polar.sh/userName"},
			{"buy_me_a_coffee", "https://www.buymeacoffee.com/userName"},
			{"thanks_dev", "https://thanks.dev/u/gh/userName"},
			{"otechie", "https://otechie.com/userName"},
		}},
		{"Custom: single schemaless URL", "custom: example.com", []f{{"custom", "https://example.com"}}},
		{"Custom: single schemaless URL in array format", "custom: [example.com]", []f{{"custom", "https://example.com"}}},
		{"Custom: double-quoted single URL", `custom: "https://example.com"`, []f{{"custom", "https://example.com"}}},
		{"Custom: double-quoted single URL in array format", `custom: ["https://example.com"]`, []f{{"custom", "https://example.com"}}},
		{"Custom: array with quoted URL and schemaless unquoted URL", `custom: ["https://example.com", example.org]`, []f{{"custom", "https://example.com"}, {"custom", "https://example.org"}}},
		{"Custom: array containing a non-simple scheme-less URL which will be discarded", `custom: [example.net/funding, "https://example.com", example.org]`, []f{{"custom", "https://example.com"}, {"custom", "https://example.org"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, hm := newGitHub(t, true, []httpExpect{
				{url: ghRepoAPIURL, body: `{"master_branch": "test_master", "owner": {"login": "composer"}, "name": "packagist"}`},
				{url: "https://api.github.com/repos/composer/packagist/contents/composer.json?ref=feature%2F3.2-foo", body: `{"encoding":"base64","content":"` + b64(`{"support": {"source": "`+ghRepoURL+`" }}`) + `"}`},
				{url: "https://api.github.com/repos/composer/packagist/commits/feature%2F3.2-foo", body: `{"commit": {"committer":{ "date": "2012-09-10"}}}`},
				{url: "https://api.github.com/repos/composer/packagist/contents/.github/FUNDING.yml", body: `{"encoding": "base64", "content": "` + b64(c.funding) + `"}`},
			}, nil, nil)
			d.tags, d.tagsOK = []Ref{{"feature/3.2-foo", "SOMESHA"}}, true
			d.branches, d.branchesOK = []Ref{{"test_master", "SOMESHA"}}, true
			data, err := d.ComposerInformation("feature/3.2-foo")
			if err != nil || data == nil {
				t.Fatalf("composer information: %v %v", data, err)
			}
			var got []f
			fv, _ := data.Get("funding")
			phpjson.Each(fv, func(_ string, e phpjson.Value) {
				m := phpjson.ToMap(e)
				got = append(got, f{m.Str("type"), m.Str("url")})
			})
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("funding = %v, want %v", got, c.want)
			}
			hm.assertComplete()
		})
	}
}

func TestGitHubPublicRepositoryArchived(t *testing.T) {
	sha := "SOMESHA"
	d, hm := newGitHub(t, true, []httpExpect{
		{url: ghRepoAPIURL, body: `{"master_branch": "test_master", "owner": {"login": "composer"}, "name": "packagist", "archived": true}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/composer.json?ref=" + sha, body: `{"encoding": "base64", "content": "` + b64(`{"name": "composer/packagist"}`) + `"}`},
		{url: "https://api.github.com/repos/composer/packagist/commits/" + sha, body: `{"commit": {"committer":{ "date": "2012-09-10"}}}`},
		{url: "https://api.github.com/repos/composer/packagist/contents/.github/FUNDING.yml", body: `{"encoding": "base64", "content": "` + b64("custom: https://example.com") + `"}`},
	}, nil, nil)
	d.tags, d.tagsOK = []Ref{{"v0.0.0", sha}}, true
	data, err := d.ComposerInformation(sha)
	if err != nil || data == nil {
		t.Fatalf("composer information: %v %v", data, err)
	}
	if v, _ := data.Get("abandoned"); v != true {
		t.Errorf("abandoned = %v", v)
	}
	hm.assertComplete()
}

func TestGitHubPrivateRepositoryNoInteraction(t *testing.T) {
	sha := "SOMESHA"
	repoSSHURL := "git@github.com:composer/packagist.git"
	proc := &procMock{strict: true}
	env, hm, _ := testEnv(t, false, nil, []httpExpect{{url: ghRepoAPIURL, status: 404}}, proc)
	cacheVcs := env.Config.GetString("cache-vcs-dir")
	proc.expects = []procExpect{
		{cmd: []string{"git", "config", "github.accesstoken"}, ret: 1},
		{cmd: []string{"git", "clone", "--mirror", "--", repoSSHURL, cacheVcs + "/git-github.com-composer-packagist.git/"}},
		{cmd: []string{"git", "remote", "-v"}},
		{cmd: []string{"git", "remote", "set-url", "origin", "--", repoSSHURL}},
		{cmd: []string{"git", "show-ref", "--tags", "--dereference"}, stdout: sha + " refs/tags/v0.0.0"},
		{cmd: []string{"git", "branch", "--no-color", "--no-abbrev", "-v"}, stdout: "  test_master     edf93f1fccaebd8764383dc12016d0a1a9672d89 Fix test & behavior"},
		{cmd: []string{"git", "branch", "--no-color"}, stdout: "* test_master"},
	}
	d := NewGitHubDriver(repoConfig(ghRepoURL), env).(*GitHubDriver)
	if err := d.Initialize(); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if root, err := d.RootIdentifier(); root != "test_master" {
		t.Errorf("root = %q (%v)", root, err)
	}
	dist, _ := d.Dist(sha)
	if dist.Str("type") != "zip" || dist.Str("url") != "https://api.github.com/repos/composer/packagist/zipball/SOMESHA" || dist.Str("reference") != sha {
		t.Errorf("dist = %v", dist)
	}
	src := d.Source(sha)
	if src.Str("type") != "git" || src.Str("url") != repoSSHURL || src.Str("reference") != sha {
		t.Errorf("source = %v", src)
	}
	hm.assertComplete()
}

func TestGitHubSupports(t *testing.T) {
	env, _, _ := testEnv(t, false, nil, nil, nil)
	cases := []struct {
		url  string
		want bool
	}{
		// GitHubDriver::supports accepts github-domains only (github.com by default)
		{"https://github.com/composer/packagist", true},
		{"http://github.com/composer/packagist.git", true},
		{"git@github.com:composer/packagist.git", true},
		{"git://github.com/composer/packagist", true},
		{"https://www.github.com/composer/packagist", true},
		{"https://gitlab.com/composer/packagist", false},
		{"https://github.com/composer", false},
		{"https://ghe.example.org/composer/packagist", false},
	}
	for _, c := range cases {
		if got := GitHubSupports(env, c.url); got != c.want {
			t.Errorf("GitHubSupports(%q) = %v, want %v", c.url, got, c.want)
		}
	}
	env2, _, _ := testEnv(t, false, map[string]phpjson.Value{"github-domains": []phpjson.Value{"github.com", "ghe.example.org"}}, nil, nil)
	if !GitHubSupports(env2, "https://ghe.example.org/composer/packagist") {
		t.Errorf("enterprise domain from github-domains not supported")
	}
}
