package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Repository/Vcs/GitBitbucketDriverTest.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"gocomposer/internal/phpjson"
)

const bbRepoData = `{"mainbranch": {"name": "main"}, "scm":"git","website":"","has_wiki":false,"name":"repo","links":{"branches":{"href":"https:\/\/api.bitbucket.org\/2.0\/repositories\/user\/repo\/refs\/branches"},"tags":{"href":"https:\/\/api.bitbucket.org\/2.0\/repositories\/user\/repo\/refs\/tags"},"clone":[{"href":"https:\/\/user@bitbucket.org\/user\/repo.git","name":"https"},{"href":"ssh:\/\/git@bitbucket.org\/user\/repo.git","name":"ssh"}],"html":{"href":"https:\/\/bitbucket.org\/user\/repo"}},"language":"php","created_on":"2015-02-18T16:22:24.688+00:00","updated_on":"2016-05-17T13:20:21.993+00:00","is_private":true,"has_issues":false}`

const (
	bbRepoURL     = "https://api.bitbucket.org/2.0/repositories/user/repo?fields=-project%2C-owner"
	bbTagsURL     = "https://api.bitbucket.org/2.0/repositories/user/repo/refs/tags?pagelen=100&fields=values.name%2Cvalues.target.hash%2Cnext&sort=-target.date"
	bbBranchesURL = "https://api.bitbucket.org/2.0/repositories/user/repo/refs/branches?pagelen=100&fields=values.name%2Cvalues.target.hash%2Cvalues.heads%2Cnext&sort=-target.date"
	bbTagsBody    = `{"values":[{"name":"1.0.1","target":{"hash":"9b78a3932143497c519e49b8241083838c8ff8a1"}},{"name":"1.0.0","target":{"hash":"d3393d514318a9267d2f8ebbf463a9aaa389f8eb"}}]}`
	bbBranchBody  = `{"values":[{"name":"main","target":{"hash":"937992d19d72b5116c3e8c4a04f960e5fa270b22"}}]}`
)

func newBitbucket(t *testing.T, u string) (*BitbucketDriver, *httpMock, error) {
	t.Helper()
	env, hm, _ := testEnv(t, false, nil, nil, nil)
	d := NewBitbucketDriver(repoConfig(u), env).(*BitbucketDriver)
	return d, hm, d.Initialize()
}

// asPlain converts a phpjson value to plain Go values (key order ignored,
// like assertEquals on associative arrays).
func asPlain(t *testing.T, v phpjson.Value) any {
	t.Helper()
	enc, err := phpjson.Encode(v, phpjson.UnescapedSlashes|phpjson.UnescapedUnicode)
	if err != nil {
		t.Fatal(err)
	}
	var out any
	if err := json.Unmarshal([]byte(enc), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestBitbucketGetRootIdentifierWrongScmType(t *testing.T) {
	d, hm, err := newBitbucket(t, "https://bitbucket.org/user/repo.git")
	if err != nil {
		t.Fatal(err)
	}
	hm.expects = []httpExpect{{url: bbRepoURL, body: `{"scm":"hg","website":"","has_wiki":false,"name":"repo","links":{"branches":{"href":"https:\/\/api.bitbucket.org\/2.0\/repositories\/user\/repo\/refs\/branches"},"tags":{"href":"https:\/\/api.bitbucket.org\/2.0\/repositories\/user\/repo\/refs\/tags"},"clone":[{"href":"https:\/\/user@bitbucket.org\/user\/repo","name":"https"},{"href":"ssh:\/\/hg@bitbucket.org\/user\/repo","name":"ssh"}],"html":{"href":"https:\/\/bitbucket.org\/user\/repo"}},"language":"php","created_on":"2015-02-18T16:22:24.688+00:00","updated_on":"2016-05-17T13:20:21.993+00:00","is_private":true,"has_issues":false}`}}
	_, err = d.RootIdentifier()
	var re *RuntimeError
	if !errors.As(err, &re) || re.Msg != "https://bitbucket.org/user/repo.git does not appear to be a git repository, use https://bitbucket.org/user/repo but remember that Bitbucket no longer supports the mercurial repositories. https://bitbucket.org/blog/sunsetting-mercurial-support-in-bitbucket" {
		t.Errorf("error = %v", err)
	}
}

func testBitbucketDriver(t *testing.T) *BitbucketDriver {
	d, hm, err := newBitbucket(t, "https://bitbucket.org/user/repo.git")
	if err != nil {
		t.Fatal(err)
	}
	hm.expects = []httpExpect{
		{url: bbRepoURL, body: bbRepoData},
		{url: bbTagsURL, body: bbTagsBody},
		{url: bbBranchesURL, body: bbBranchBody},
		{url: "https://api.bitbucket.org/2.0/repositories/user/repo/src/main/composer.json", body: `{"name": "user/repo","description": "test repo","license": "GPL","authors": [{"name": "Name","email": "local@domain.tld"}],"require": {"creator/package": "^1.0"},"require-dev": {"phpunit/phpunit": "~4.8"}}`},
		{url: "https://api.bitbucket.org/2.0/repositories/user/repo/commit/main?fields=date", body: `{"date": "2016-05-17T13:19:52+00:00"}`},
	}
	if root, err := d.RootIdentifier(); root != "main" {
		t.Errorf("root = %q (%v)", root, err)
	}
	if tags, err := d.Tags(); err != nil || !reflect.DeepEqual(tags, []Ref{{"1.0.1", "9b78a3932143497c519e49b8241083838c8ff8a1"}, {"1.0.0", "d3393d514318a9267d2f8ebbf463a9aaa389f8eb"}}) {
		t.Errorf("tags = %v (%v)", tags, err)
	}
	if br, err := d.Branches(); err != nil || !reflect.DeepEqual(br, []Ref{{"main", "937992d19d72b5116c3e8c4a04f960e5fa270b22"}}) {
		t.Errorf("branches = %v (%v)", br, err)
	}
	data, err := d.ComposerInformation("main")
	if err != nil {
		t.Fatal(err)
	}
	var want any
	_ = json.Unmarshal([]byte(`{
		"name": "user/repo",
		"description": "test repo",
		"license": "GPL",
		"authors": [{"name": "Name", "email": "local@domain.tld"}],
		"require": {"creator/package": "^1.0"},
		"require-dev": {"phpunit/phpunit": "~4.8"},
		"time": "2016-05-17T13:19:52+00:00",
		"support": {"source": "https://bitbucket.org/user/repo/src/937992d19d72b5116c3e8c4a04f960e5fa270b22/?at=main"},
		"homepage": "https://bitbucket.org/user/repo"
	}`), &want)
	if got := asPlain(t, data); !reflect.DeepEqual(got, want) {
		t.Errorf("composer information = %v, want %v", got, want)
	}
	hm.assertComplete()
	return d
}

func TestBitbucketDriver(t *testing.T) { testBitbucketDriver(t) }

func TestBitbucketGetParams(t *testing.T) {
	d := testBitbucketDriver(t)
	u := "https://bitbucket.org/user/repo.git"
	if d.URL() != u {
		t.Errorf("url = %q", d.URL())
	}
	dist, _ := d.Dist("reference")
	if got := flatMap(dist); !reflect.DeepEqual(got, map[string]string{"type": "zip", "url": "https://bitbucket.org/user/repo/get/reference.zip", "reference": "reference", "shasum": ""}) {
		t.Errorf("dist = %v", got)
	}
	if got := flatMap(d.Source("reference")); !reflect.DeepEqual(got, map[string]string{"type": "git", "url": u, "reference": "reference"}) {
		t.Errorf("source = %v", got)
	}
}

func TestBitbucketInitializeInvalidRepositoryURL(t *testing.T) {
	_, _, err := newBitbucket(t, "https://bitbucket.org/acme")
	var ia *InvalidArgumentError
	if !errors.As(err, &ia) {
		t.Errorf("error = %#v, want InvalidArgumentError", err)
	}
}

func TestBitbucketInvalidSupportData(t *testing.T) {
	repoURL := "https://bitbucket.org/user/repo.git"
	d, hm, err := newBitbucket(t, repoURL)
	if err != nil {
		t.Fatal(err)
	}
	hm.expects = []httpExpect{
		{url: bbRepoURL, body: bbRepoData},
		{url: "https://api.bitbucket.org/2.0/repositories/user/repo/src/main/composer.json", body: `{"support": "` + repoURL + `"}`},
		{url: "https://api.bitbucket.org/2.0/repositories/user/repo/commit/main?fields=date", body: `{"date": "2016-05-17T13:19:52+00:00"}`},
		{url: bbTagsURL, body: bbTagsBody},
		{url: bbBranchesURL, body: bbBranchBody},
	}
	if _, err := d.RootIdentifier(); err != nil {
		t.Fatal(err)
	}
	data, err := d.ComposerInformation("main")
	if err != nil || data == nil {
		t.Fatalf("composer information: %v %v", data, err)
	}
	if v, _ := mapGet(data, "support", "source"); strOf(v) != "https://bitbucket.org/user/repo/src/937992d19d72b5116c3e8c4a04f960e5fa270b22/?at=main" {
		t.Errorf("support.source = %v", v)
	}
	hm.assertComplete()
}

func TestBitbucketSupports(t *testing.T) {
	if !BitbucketSupports("https://bitbucket.org/user/repo.git") {
		t.Error("https url not supported")
	}
	// should not be changed, see https://github.com/composer/composer/issues/9400
	if BitbucketSupports("git@bitbucket.org:user/repo.git") {
		t.Error("ssh url supported")
	}
	if BitbucketSupports("https://github.com/user/repo.git") {
		t.Error("github url supported")
	}
}
