package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Repository/Vcs/GitLabDriverTest.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
//
// testForwardsOptions is not ported: gocomposer has no per-repository
// stream context "options" for the GitLab API requests.

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"gocomposer/internal/phpjson"
)

// gitlabTestConfig is the setUp() config. secure-http is off because the
// upstream test replaces the whole HttpDownloader (which enforces it) by a
// mock; here only the transport is mocked.
func gitlabTestConfig() map[string]phpjson.Value {
	return map[string]phpjson.Value{"secure-http": false, "gitlab-domains": []phpjson.Value{
		"mycompany.com/gitlab",
		"gitlab.mycompany.com",
		"othercompany.com/nested/gitlab",
		"gitlab.com",
		"gitlab.mycompany.local",
	}}
}

func newGitLab(t *testing.T, u string, cfg map[string]phpjson.Value, expects []httpExpect) (*GitLabDriver, *httpMock) {
	t.Helper()
	if cfg == nil {
		cfg = gitlabTestConfig()
	}
	env, hm, _ := testEnv(t, false, cfg, expects, nil)
	d := NewGitLabDriver(repoConfig(u), env).(*GitLabDriver)
	if err := d.Initialize(); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	return d, hm
}

var gitlabInitializeURLs = [][2]string{
	{"https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject"},
	{"http://gitlab.com/mygroup/myproject", "http://gitlab.com/api/v4/projects/mygroup%2Fmyproject"},
	{"git@gitlab.com:mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject"},
}

const gitlabPrivateProject = `{
    "id": 17,
    "default_branch": "mymaster",
    "visibility": "private",
    "issues_enabled": true,
    "archived": false,
    "http_url_to_repo": "https://gitlab.com/mygroup/myproject.git",
    "ssh_url_to_repo": "git@gitlab.com:mygroup/myproject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/myproject",
    "web_url": "https://gitlab.com/mygroup/myproject"
}`

const gitlabPublicProject = `{
    "id": 17,
    "default_branch": "mymaster",
    "visibility": "public",
    "http_url_to_repo": "https://gitlab.com/mygroup/myproject.git",
    "ssh_url_to_repo": "git@gitlab.com:mygroup/myproject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/myproject",
    "web_url": "https://gitlab.com/mygroup/myproject"
}`

const gitlabAnonymousProject = `{
    "id": 17,
    "default_branch": "mymaster",
    "http_url_to_repo": "https://gitlab.com/mygroup/myproject.git",
    "ssh_url_to_repo": "git@gitlab.com:mygroup/myproject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/myproject",
    "web_url": "https://gitlab.com/mygroup/myproject"
}`

func assertGitLabInit(t *testing.T, d *GitLabDriver, apiURL, root, repoURL, webURL string) {
	t.Helper()
	if got := d.apiURL(); got != apiURL {
		t.Errorf("API URL = %q, want %q", got, apiURL)
	}
	if got, _ := d.RootIdentifier(); got != root {
		t.Errorf("root identifier = %q, want %q", got, root)
	}
	if got := d.repositoryURL(); got != repoURL {
		t.Errorf("repository URL = %q, want %q", got, repoURL)
	}
	if got := d.URL(); got != webURL {
		t.Errorf("URL = %q, want %q", got, webURL)
	}
}

func testGitLabInitialize(t *testing.T, u, apiURL string) (*GitLabDriver, *httpMock) {
	d, hm := newGitLab(t, u, nil, []httpExpect{{url: apiURL, body: gitlabPrivateProject}})
	assertGitLabInit(t, d, apiURL, "mymaster", "git@gitlab.com:mygroup/myproject.git", "https://gitlab.com/mygroup/myproject")
	hm.assertComplete()
	return d, hm
}

func TestGitLabInitialize(t *testing.T) {
	for _, c := range gitlabInitializeURLs {
		t.Run(c[0], func(t *testing.T) { testGitLabInitialize(t, c[0], c[1]) })
	}
}

func TestGitLabInitializePublicProject(t *testing.T) {
	for _, c := range gitlabInitializeURLs {
		t.Run(c[0], func(t *testing.T) {
			d, hm := newGitLab(t, c[0], nil, []httpExpect{{url: c[1], body: gitlabPublicProject}})
			assertGitLabInit(t, d, c[1], "mymaster", "https://gitlab.com/mygroup/myproject.git", "https://gitlab.com/mygroup/myproject")
			hm.assertComplete()
		})
	}
}

func TestGitLabInitializePublicProjectAsAnonymous(t *testing.T) {
	for _, c := range gitlabInitializeURLs {
		t.Run(c[0], func(t *testing.T) {
			d, hm := newGitLab(t, c[0], nil, []httpExpect{{url: c[1], body: gitlabAnonymousProject}})
			assertGitLabInit(t, d, c[1], "mymaster", "https://gitlab.com/mygroup/myproject.git", "https://gitlab.com/mygroup/myproject")
			hm.assertComplete()
		})
	}
}

func TestGitLabInitializeWithPortNumber(t *testing.T) {
	domain, port, namespace := "gitlab.mycompany.com", "5443", "mygroup/myproject"
	u := fmt.Sprintf("https://%s:%s/%s", domain, port, namespace)
	apiURL := fmt.Sprintf("https://%s:%s/api/v4/projects/%s", domain, port, url.QueryEscape(namespace))
	body := fmt.Sprintf(`{
    "default_branch": "1.0.x",
    "http_url_to_repo": "https://%[1]s:%[2]s/%[3]s.git",
    "path": "myproject",
    "path_with_namespace": "%[3]s",
    "web_url": "https://%[1]s:%[2]s/%[3]s"
}`, domain, port, namespace)
	d, hm := newGitLab(t, u, nil, []httpExpect{{url: apiURL, body: body}})
	assertGitLabInit(t, d, apiURL, "1.0.x", u+".git", u)
	hm.assertComplete()
}

func TestGitLabInvalidSupportData(t *testing.T) {
	repoURL := "https://gitlab.com/mygroup/myproject"
	d, hm := testGitLabInitialize(t, repoURL, "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	d.branches, d.branchesOK = []Ref{{"main", "SOMESHA"}}, true
	d.tags, d.tagsOK = nil, true
	hm.expects = []httpExpect{{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/files/composer%2Ejson/raw?ref=SOMESHA", body: `{"support": "` + repoURL + `" }`}}
	data, err := d.ComposerInformation("main")
	if err != nil || data == nil {
		t.Fatalf("composer information: %v %v", data, err)
	}
	if v, _ := mapGet(data, "support", "source"); strOf(v) != "https://gitlab.com/mygroup/myproject/-/tree/main" {
		t.Errorf("support.source = %v", v)
	}
	hm.assertComplete()
}

func TestGitLabGetDist(t *testing.T) {
	d, _ := testGitLabInitialize(t, "https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	ref := "c3ebdbf9cceddb82cd2089aaef8c7b992e536363"
	dist, err := d.Dist(ref)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"type": "zip", "url": "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/archive.zip?sha=" + ref, "reference": ref, "shasum": ""}
	if got := flatMap(dist); !reflect.DeepEqual(got, want) {
		t.Errorf("dist = %v, want %v", got, want)
	}
}

func TestGitLabGetSource(t *testing.T) {
	d, _ := testGitLabInitialize(t, "https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	ref := "c3ebdbf9cceddb82cd2089aaef8c7b992e536363"
	want := map[string]string{"type": "git", "url": "git@gitlab.com:mygroup/myproject.git", "reference": ref}
	if got := flatMap(d.Source(ref)); !reflect.DeepEqual(got, want) {
		t.Errorf("source = %v, want %v", got, want)
	}
}

func TestGitLabGetSourceGivenPublicProject(t *testing.T) {
	d, _ := newGitLab(t, "https://gitlab.com/mygroup/myproject", nil, []httpExpect{{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject", body: gitlabPublicProject}})
	ref := "c3ebdbf9cceddb82cd2089aaef8c7b992e536363"
	want := map[string]string{"type": "git", "url": "https://gitlab.com/mygroup/myproject.git", "reference": ref}
	if got := flatMap(d.Source(ref)); !reflect.DeepEqual(got, want) {
		t.Errorf("source = %v, want %v", got, want)
	}
}

func TestGitLabGetTags(t *testing.T) {
	d, hm := testGitLabInitialize(t, "https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	hm.expects = []httpExpect{{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?per_page=100", body: `[
    {
       "name": "v1.0.0",
        "commit": {
            "id": "092ed2c762bbae331e3f51d4a17f67310bf99a81",
            "committed_date": "2012-05-28T04:42:42-07:00"
        }
    },
    {
        "name": "v2.0.0",
        "commit": {
            "id": "8e8f60b3ec86d63733db3bd6371117a758027ec6",
            "committed_date": "2014-07-06T12:59:11.000+02:00"
        }
    }
]`}}
	want := []Ref{{"v1.0.0", "092ed2c762bbae331e3f51d4a17f67310bf99a81"}, {"v2.0.0", "8e8f60b3ec86d63733db3bd6371117a758027ec6"}}
	for i := 0; i < 2; i++ { // the second call is served from the cache
		if got, err := d.Tags(); err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("tags = %v (%v), want %v", got, err, want)
		}
	}
	hm.assertComplete()
}

func TestGitLabGetPaginatedRefs(t *testing.T) {
	d, hm := testGitLabInitialize(t, "https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	type commit struct {
		ID            string `json:"id"`
		CommittedDate string `json:"committed_date"`
	}
	type branch struct {
		Name   string `json:"name"`
		Commit commit `json:"commit"`
	}
	data := []branch{
		{"mymaster", commit{"97eda36b5c1dd953a3792865c222d4e85e5f302e", "2013-01-03T21:04:07.000+01:00"}},
		{"staging", commit{"502cffe49f136443f2059803f2e7192d1ac066cd", "2013-03-09T16:35:23.000+01:00"}},
	}
	for i := 0; i < 98; i++ {
		data = append(data, branch{"stagingdupe", commit{"502cffe49f136443f2059803f2e7192d1ac066cd", "2013-03-09T16:35:23.000+01:00"}})
	}
	body, _ := json.Marshal(data)
	hm.expects = []httpExpect{
		{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/branches?per_page=100", body: string(body),
			headers: map[string]string{"Link": `<http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=2&per_page=20>; rel="next", <http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=1&per_page=20>; rel="first", <http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=3&per_page=20>; rel="last"`}},
		{url: "http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=2&per_page=20", body: string(body),
			headers: map[string]string{"Link": `<http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=2&per_page=20>; rel="prev", <http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=1&per_page=20>; rel="first", <http://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/tags?id=mygroup%2Fmyproject&page=3&per_page=20>; rel="last"`}},
	}
	want := []Ref{
		{"mymaster", "97eda36b5c1dd953a3792865c222d4e85e5f302e"},
		{"staging", "502cffe49f136443f2059803f2e7192d1ac066cd"},
		{"stagingdupe", "502cffe49f136443f2059803f2e7192d1ac066cd"},
	}
	for i := 0; i < 2; i++ {
		if got, err := d.Branches(); err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("branches = %v (%v), want %v", got, err, want)
		}
	}
	hm.assertComplete()
}

func TestGitLabGetBranches(t *testing.T) {
	d, hm := testGitLabInitialize(t, "https://gitlab.com/mygroup/myproject", "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject")
	hm.expects = []httpExpect{{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject/repository/branches?per_page=100", body: `[
    {
       "name": "mymaster",
        "commit": {
            "id": "97eda36b5c1dd953a3792865c222d4e85e5f302e",
            "committed_date": "2013-01-03T21:04:07.000+01:00"
        }
    },
    {
        "name": "staging",
        "commit": {
            "id": "502cffe49f136443f2059803f2e7192d1ac066cd",
            "committed_date": "2013-03-09T16:35:23.000+01:00"
        }
    }
]`}}
	want := []Ref{{"mymaster", "97eda36b5c1dd953a3792865c222d4e85e5f302e"}, {"staging", "502cffe49f136443f2059803f2e7192d1ac066cd"}}
	for i := 0; i < 2; i++ {
		if got, err := d.Branches(); err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("branches = %v (%v), want %v", got, err, want)
		}
	}
	hm.assertComplete()
}

func TestGitLabSupports(t *testing.T) {
	env, _, _ := testEnv(t, false, gitlabTestConfig(), nil, nil)
	cases := []struct {
		url  string
		want bool
	}{
		{"http://gitlab.com/foo/bar", true},
		{"http://gitlab.mycompany.com:5443/foo/bar", true},
		{"http://gitlab.com/foo/bar/", true},
		{"http://gitlab.com/foo/bar/", true},
		{"http://gitlab.com/foo/bar.git", true},
		{"http://gitlab.com/foo/bar.git", true},
		{"http://gitlab.com/foo/bar.baz.git", true},
		{"https://gitlab.com/foo/bar", true}, // extension_loaded('openssl')
		{"https://gitlab.mycompany.com:5443/foo/bar", true},
		{"git@gitlab.com:foo/bar.git", true},
		{"git@example.com:foo/bar.git", false},
		{"http://example.com/foo/bar", false},
		{"http://mycompany.com/gitlab/mygroup/myproject", true},
		{"https://mycompany.com/gitlab/mygroup/myproject", true},
		{"http://othercompany.com/nested/gitlab/mygroup/myproject", true},
		{"https://othercompany.com/nested/gitlab/mygroup/myproject", true},
		{"http://gitlab.com/mygroup/mysubgroup/mysubsubgroup/myproject", true},
		{"https://gitlab.com/mygroup/mysubgroup/mysubsubgroup/myproject", true},
	}
	for _, c := range cases {
		if got := GitLabSupports(env, c.url); got != c.want {
			t.Errorf("GitLabSupports(%q) = %v, want %v", c.url, got, c.want)
		}
	}
}

func TestGitLabSubDirectory(t *testing.T) {
	apiURL := "https://mycompany.com/gitlab/api/v4/projects/mygroup%2Fmy-pro%2Eject"
	d, hm := newGitLab(t, "https://mycompany.com/gitlab/mygroup/my-pro.ject", nil, []httpExpect{{url: apiURL, body: `{
    "id": 17,
    "default_branch": "mymaster",
    "visibility": "private",
    "http_url_to_repo": "https://gitlab.com/gitlab/mygroup/my-pro.ject",
    "ssh_url_to_repo": "git@gitlab.com:mygroup/my-pro.ject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/my-pro.ject",
    "web_url": "https://gitlab.com/gitlab/mygroup/my-pro.ject"
}`}})
	if got := d.apiURL(); got != apiURL {
		t.Errorf("API URL = %q, want %q", got, apiURL)
	}
	hm.assertComplete()
}

func TestGitLabSubGroup(t *testing.T) {
	apiURL := "https://gitlab.com/api/v4/projects/mygroup%2Fmysubgroup%2Fmyproject"
	d, hm := newGitLab(t, "https://gitlab.com/mygroup/mysubgroup/myproject", nil, []httpExpect{{url: apiURL, body: `{
    "id": 17,
    "default_branch": "mymaster",
    "visibility": "private",
    "http_url_to_repo": "https://gitlab.com/mygroup/mysubgroup/my-pro.ject",
    "ssh_url_to_repo": "git@gitlab.com:mygroup/mysubgroup/my-pro.ject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/mysubgroup/my-pro.ject",
    "web_url": "https://gitlab.com/mygroup/mysubgroup/my-pro.ject"
}`}})
	if got := d.apiURL(); got != apiURL {
		t.Errorf("API URL = %q, want %q", got, apiURL)
	}
	hm.assertComplete()
}

func TestGitLabSubDirectorySubGroup(t *testing.T) {
	apiURL := "https://mycompany.com/gitlab/api/v4/projects/mygroup%2Fmysubgroup%2Fmyproject"
	d, hm := newGitLab(t, "https://mycompany.com/gitlab/mygroup/mysubgroup/myproject", nil, []httpExpect{{url: apiURL, body: `{
    "id": 17,
    "default_branch": "mymaster",
    "visibility": "private",
    "http_url_to_repo": "https://mycompany.com/gitlab/mygroup/mysubgroup/my-pro.ject",
    "ssh_url_to_repo": "git@mycompany.com:mygroup/mysubgroup/my-pro.ject.git",
    "last_activity_at": "2014-12-01T09:17:51.000+01:00",
    "name": "My Project",
    "name_with_namespace": "My Group / My Project",
    "path": "myproject",
    "path_with_namespace": "mygroup/mysubgroup/my-pro.ject",
    "web_url": "https://mycompany.com/gitlab/mygroup/mysubgroup/my-pro.ject"
}`}})
	if got := d.apiURL(); got != apiURL {
		t.Errorf("API URL = %q, want %q", got, apiURL)
	}
	hm.assertComplete()
}

func TestGitLabProtocolOverrideRepositoryURLGeneration(t *testing.T) {
	cfg := gitlabTestConfig()
	cfg["gitlab-protocol"] = "http"
	d, hm := newGitLab(t, "git@gitlab.com:mygroup/myproject", cfg, []httpExpect{{url: "https://gitlab.com/api/v4/projects/mygroup%2Fmyproject", body: gitlabPrivateProject}})
	if got := d.repositoryURL(); got != "https://gitlab.com/mygroup/myproject.git" {
		t.Errorf("repository URL = %q (gitlab-protocol http)", got)
	}
	hm.assertComplete()
}

// flatMap renders a dist/source array as strings for comparison.
func flatMap(m *phpjson.Map) map[string]string {
	out := map[string]string{}
	if m == nil {
		return out
	}
	for _, k := range m.Keys() {
		v, _ := m.Get(k)
		out[k] = strings.TrimSpace(strOf(v))
	}
	return out
}
