package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Util/SvnTest.php, Util/PerforceTest.php,
// Repository/Vcs/SvnDriverTest.php, HgDriverTest.php, FossilDriverTest.php,
// PerforceDriverTest.php). Copyright (c) Nils Adermann, Jordi Boggiano.
// MIT License.

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"gocomposer/internal/phpjson"
)

func svnArgs(t *testing.T, u string, cfg map[string]phpjson.Value, cache *bool) []string {
	t.Helper()
	env, _, _ := testEnv(t, false, cfg, nil, nil)
	s := NewSvn(u, env)
	if cache != nil {
		s.SetCacheCredentials(*cache)
	}
	cmd, _ := s.command([]string{"svn", "ls"}, u, nil)
	// drop the fixed parts: svn ls --non-interactive ... -- url
	return cmd[3 : len(cmd)-2]
}

func TestSvnCredentials(t *testing.T) {
	for _, c := range []struct {
		url  string
		want []string
	}{
		{"http://till:test@svn.example.org/", []string{"--username", "till", "--password", "test"}},
		{"http://svn.apache.org/", []string{}},
		{"svn://johndoe@example.org", []string{"--username", "johndoe", "--password", ""}},
	} {
		if got := svnArgs(t, c.url, nil, nil); !reflect.DeepEqual(append([]string{}, got...), c.want) && !(len(got) == 0 && len(c.want) == 0) {
			t.Errorf("%s: got %q, want %q", c.url, got, c.want)
		}
	}
}

func TestSvnInteractiveString(t *testing.T) {
	env, _, _ := testEnv(t, false, nil, nil, nil)
	cmd, _ := NewSvn("http://svn.example.org", env).command([]string{"svn", "ls"}, "http://svn.example.org", nil)
	want := []string{"svn", "ls", "--non-interactive", "--", "http://svn.example.org"}
	if !reflect.DeepEqual(cmd, want) {
		t.Errorf("got %q, want %q", cmd, want)
	}
}

func TestSvnCredentialsFromConfig(t *testing.T) {
	basic := map[string]phpjson.Value{"http-basic": func() phpjson.Value {
		m := phpjson.NewMap()
		c := phpjson.NewMap()
		c.Set("username", "foo")
		c.Set("password", "bar")
		m.Set("svn.apache.org", c)
		return m
	}()}
	yes, no := true, false
	if got := svnArgs(t, "http://svn.apache.org", basic, nil); !reflect.DeepEqual(got, []string{"--username", "foo", "--password", "bar"}) {
		t.Errorf("default: %q", got)
	}
	if got := svnArgs(t, "http://svn.apache.org", basic, &yes); !reflect.DeepEqual(got, []string{"--username", "foo", "--password", "bar"}) {
		t.Errorf("cache true: %q", got)
	}
	if got := svnArgs(t, "http://svn.apache.org", basic, &no); !reflect.DeepEqual(got, []string{"--no-auth-cache", "--username", "foo", "--password", "bar"}) {
		t.Errorf("cache false: %q", got)
	}
}

func TestSvnDriverWrongCredentialsInURL(t *testing.T) {
	svnVersionMu.Lock()
	svnVersion = ""
	svnVersionMu.Unlock()
	t.Cleanup(func() { svnVersionMu.Lock(); svnVersion = ""; svnVersionMu.Unlock() })
	out := "svn: OPTIONS of 'https://corp.svn.local/repo': authorization failed: Could not authenticate to server: rejected Basic challenge (https://corp.svn.local/)"
	authed := []string{"svn", "ls", "--verbose", "--non-interactive", "--username", "till", "--password", "secret", "--", "https://till:secret@corp.svn.local/repo/trunk"}
	var exp []procExpect
	for i := 0; i < 6; i++ {
		exp = append(exp, procExpect{cmd: authed, ret: 1, stderr: out})
	}
	exp = append(exp, procExpect{cmd: []string{"svn", "--version"}, stdout: "1.2.3"})
	proc := &procMock{strict: true, expects: exp}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewSvnDriver(repoConfig("https://till:secret@corp.svn.local/repo"), env)
	err := d.Initialize()
	want := "Repository https://till:***@corp.svn.local/repo could not be processed, wrong credentials provided (" + out + ")"
	if err == nil || err.Error() != want {
		t.Errorf("err = %v\nwant %s", err, want)
	}
	proc.assertComplete()
}

func TestSvnDriverSupports(t *testing.T) {
	env, _, _ := testEnv(t, false, nil, nil, nil)
	for _, u := range []string{"http://svn.apache.org", "https://svn.sf.net", "svn://example.org", "svn+ssh://example.org"} {
		if !svnSupports(env, u, false) {
			t.Errorf("%s not supported", u)
		}
	}
}

func TestSvnDriverTagsBranchesAndIdentifiers(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"svn", "ls", "--verbose", "--non-interactive", "--", "https://svn.example.org/repo/trunk"}, stdout: "      7 fixture               Jan 02  2024 ./\n      7 fixture            19 Jan 02  2024 composer.json\n"},
		{cmd: []string{"svn", "ls", "--verbose", "--non-interactive", "--", "https://svn.example.org/repo/branches"}, stdout: "      6 fixture               Jan 02  2024 ./\n      2 fixture               Jan 02  2024 2.x/\n"},
		{cmd: []string{"svn", "ls", "--verbose", "--non-interactive", "--", "https://svn.example.org/repo/tags"}, stdout: "      5 fixture               Jan 02  2024 ./\n      3 fixture               Jan 02  2024 1.0.0/\n      9 fixture               Jan 02  2024 1.1.0/\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewSvnDriver(repoConfig("https://svn.example.org/repo/trunk/", "package-path", "/pkg/"), env)
	if err := d.Initialize(); err != nil {
		t.Fatal(err)
	}
	proc.assertComplete()
	if d.baseURL != "https://svn.example.org/repo" || d.URL() != "https://svn.example.org/repo/trunk" {
		t.Errorf("base %q url %q", d.baseURL, d.URL())
	}
	br, _ := d.Branches()
	if want := []Ref{{"trunk", "/trunk/pkg/@7"}, {"2.x", "/branches/2.x/pkg/@6"}}; !reflect.DeepEqual(br, want) {
		t.Errorf("branches %v", br)
	}
	tags, _ := d.Tags()
	if want := []Ref{{"1.0.0", "/tags/1.0.0/pkg/@5"}, {"1.1.0", "/tags/1.1.0/pkg/@9"}}; !reflect.DeepEqual(tags, want) {
		t.Errorf("tags %v", tags)
	}
	if root, _ := d.RootIdentifier(); root != "/trunk/pkg/@7" {
		t.Errorf("root %q", root)
	}
	if src := d.Source("/tags/1.0.0/pkg/@5"); src.Str("url") != "https://svn.example.org/repo" || src.Str("type") != "svn" {
		t.Errorf("source %v", src)
	}
	if p, r := svnSplitIdentifier("/tags/1.0.0/pkg/@5"); p != "/tags/1.0.0/pkg/" || r != "@5" {
		t.Errorf("split %q %q", p, r)
	}
}

func TestSvnDriverMissingComposerJSONIsRemembered(t *testing.T) {
	svnVersionMu.Lock()
	svnVersion = ""
	svnVersionMu.Unlock()
	t.Cleanup(func() { svnVersionMu.Lock(); svnVersion = ""; svnVersionMu.Unlock() })
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"svn", "cat", "--non-interactive", "--", "https://svn.example.org/repo/tags/1.0.0/composer.json@3"}, ret: 1, stderr: "svn: warning: W160013: Path '/tags/1.0.0/composer.json' not found\n"},
		{cmd: []string{"svn", "--version"}, stdout: "svn, version 1.14.5 (r1922182)\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewSvnDriver(repoConfig("https://svn.example.org/repo"), env)
	d.baseURL = "https://svn.example.org/repo"
	m, err := d.ComposerInformation("/tags/1.0.0/@3")
	if err != nil || m != nil {
		t.Fatalf("got %v %v", m, err)
	}
	// cached: no second svn cat
	if m, err := d.ComposerInformation("/tags/1.0.0/@3"); err != nil || m != nil {
		t.Fatalf("got %v %v", m, err)
	}
	proc.assertComplete()
}

func TestHgDriverSupports(t *testing.T) {
	env, _, _ := testEnv(t, false, nil, nil, nil)
	for _, u := range []string{"ssh://bitbucket.org/user/repo", "ssh://hg@bitbucket.org/user/repo", "ssh://user@bitbucket.org/user/repo", "https://bitbucket.org/user/repo", "https://user@bitbucket.org/user/repo"} {
		if !hgSupports(env, u, false) {
			t.Errorf("%s not supported", u)
		}
	}
}

func TestHgDriverGetBranchesFilterInvalidBranchNames(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"hg", "branches"}, stdout: "default 1:dbf6c8acb640\n--help  1:dbf6c8acb640"},
		{cmd: []string{"hg", "bookmarks"}, stdout: "help    1:dbf6c8acb641\n--help  1:dbf6c8acb641\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewHgDriver(repoConfig("https://example.org/acme.git"), env)
	br, err := d.Branches()
	if err != nil {
		t.Fatal(err)
	}
	if want := []Ref{{"help", "dbf6c8acb641"}, {"default", "dbf6c8acb640"}}; !reflect.DeepEqual(br, want) {
		t.Errorf("branches %v", br)
	}
}

func TestHgDriverTagsDropTip(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"hg", "tags"}, stdout: "tip                                5:cbb3dc6fedb9\n1.1.0                              2:f11fae9f33d6\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewHgDriver(repoConfig("/tmp/hg"), env)
	tags, _ := d.Tags()
	if want := []Ref{{"1.1.0", "f11fae9f33d6"}}; !reflect.DeepEqual(tags, want) {
		t.Errorf("tags %v", tags)
	}
}

func TestHgDriverInvalidIdentifier(t *testing.T) {
	proc := &procMock{expects: []procExpect{{cmd: []string{"hg", "cat", "-r", "h", "--", "file.txt"}}}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewHgDriver(repoConfig("https://example.org/acme.git"), env)
	if _, ok, err := d.FileContent("file.txt", "h"); ok || err != nil {
		t.Errorf("h: %v %v", ok, err)
	}
	if _, _, err := d.FileContent("file.txt", "-h"); err == nil {
		t.Error("no error for -h")
	}
	if _, _, err := d.ChangeDate("-r foo"); err == nil {
		t.Error("no error for -r foo")
	}
}

func TestHgRunCommandRetriesWithCredentials(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"hg", "clone", "--noupdate", "--", "https://hg.example.org/repo", "/dir"}, ret: 255, stderr: "abort: authorization failed\n"},
		{cmd: []string{"hg", "clone", "--noupdate", "--", "https://us%40er:p%40ss@hg.example.org/repo", "/dir"}, ret: 255, stderr: "abort: authorization failed\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	env.Auth.Set("hg.example.org", "us@er", "p@ss")
	hgVersionOnce.Do(func() { hgVersion = "7.0.1" })
	err := NewHg(env).RunCommand(func(u string) []string { return []string{"hg", "clone", "--noupdate", "--", u, "/dir"} }, "https://hg.example.org/repo", "")
	if err == nil || err.Error() != "Failed to clone https://hg.example.org/repo, \n\nabort: authorization failed\n" {
		t.Errorf("err = %q", err)
	}
	proc.assertComplete()
}

func TestFossilDriverSupports(t *testing.T) {
	env, _, _ := testEnv(t, false, nil, nil, nil)
	for _, c := range []struct {
		url string
		ok  bool
	}{
		{"http://fossil.kd2.org/kd2fw/", true},
		{"https://chiselapp.com/user/rkeene/repository/flint/index", true},
		{"ssh://fossil.kd2.org/kd2fw.fossil", true},
	} {
		if got := fossilSupports(env, c.url); got != c.ok {
			t.Errorf("%s: %v", c.url, got)
		}
	}
}

func TestFossilDriverBranchesAndTags(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"fossil", "tag", "list"}, stdout: "1.0.0\n1.1.0\n2.x\ntrunk\n"},
		{cmd: []string{"fossil", "branch", "list"}, stdout: "   2.x\n * trunk\n"},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewFossilDriver(repoConfig("/tmp/fossil"), env)
	tags, _ := d.Tags()
	br, _ := d.Branches()
	if want := []Ref{{"1.0.0", "1.0.0"}, {"1.1.0", "1.1.0"}, {"2.x", "2.x"}, {"trunk", "trunk"}}; !reflect.DeepEqual(tags, want) {
		t.Errorf("tags %v", tags)
	}
	if want := []Ref{{"2.x", "2.x"}, {"trunk", "trunk"}}; !reflect.DeepEqual(br, want) {
		t.Errorf("branches %v", br)
	}
	if root, _ := d.RootIdentifier(); root != "trunk" {
		t.Errorf("root %q", root)
	}
}

// newTestPerforce is PerforceTest::setUp on a non-Windows platform.
func newTestPerforce(t *testing.T, proc *procMock) (*Perforce, string) {
	t.Helper()
	env, _, _ := testEnv(t, false, nil, nil, proc)
	path := t.TempDir()
	rc := phpjson.NewMap()
	rc.Set("depot", "depot")
	rc.Set("branch", "branch")
	rc.Set("p4user", "user")
	rc.Set("unique_perforce_client_name", "TEST")
	p, err := NewPerforce(rc, "port", path, env)
	if err != nil {
		t.Fatal(err)
	}
	return p, path
}

func init() {
	// ExecutableFinder found no p4 in Composer's tests either
	p4ExecutableOnce.Do(func() { p4Executable = "p4" })
}

func TestPerforceClientAndStream(t *testing.T) {
	p, path := newTestPerforce(t, nil)
	if c := p.Client(); c != "composer_perforce_TEST_depot" {
		t.Errorf("client %q", c)
	}
	if s := p.Stream(); s != "//depot" {
		t.Errorf("stream %q", s)
	}
	if s := p.ClientSpec(); s != path+"/composer_perforce_TEST_depot.p4.spec" {
		t.Errorf("spec %q", s)
	}
	if got := p.generateP4Command([]string{"do", "something"}, true); !reflect.DeepEqual(got, []string{"p4", "-u", "user", "-c", "composer_perforce_TEST_depot", "-p", "port", "do", "something"}) {
		t.Errorf("command %q", got)
	}
	ps, _ := newTestPerforce(t, nil)
	ps.SetStream("//depot/branch")
	if c := ps.Client(); c != "composer_perforce_TEST_depot_branch" {
		t.Errorf("stream client %q", c)
	}
	if s := streamWithoutLabel("//depot/branching@label"); s != "//depot/branching" {
		t.Errorf("without label %q", s)
	}
}

func TestPerforceQueryP4UserEscapesInjectionOnUnix(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"/bin/sh", "-c", "echo $P4USER"}},
		{cmd: []string{"/bin/sh", "-c", "export P4USER='foo; id'"}},
	}}
	p, _ := newTestPerforce(t, proc)
	p.p4User = optStr{}
	p.env.IO.Interactive = true
	p.env.IO.In = strings.NewReader("foo; id\n")
	if err := p.queryP4User(); err != nil {
		t.Fatal(err)
	}
	proc.assertComplete()
}

func TestPerforceBranchesTagsAndComposerInformation(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"p4", "-u", "user", "-p", "port", "changes", "//depot/..."}, stdout: "Change 5678 on 2014/03/19 by Clark.Stuth@Clark.Stuth_test_client 'test changelist'"},
		{cmd: []string{"p4", "-u", "user", "-c", "composer_perforce_TEST_depot", "-p", "port", "labels"}, stdout: "Label 0.0.1 2013/07/31 'First Label!'\nLabel 0.0.2 2013/08/01 'Second Label!'\n"},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "files", "//depot/composer.json@0.0.1"}, stdout: "//depot/composer.json#1 - branch change 10001 (text)"},
		{cmd: []string{"p4", "-u", "user", "-c", "composer_perforce_TEST_depot", "-p", "port", "print", "//depot/composer.json@10001"}, stdout: `{"name": "test/perforce", "minimum-stability": "dev"}`},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "depots"}, stdout: "Depot depot 2013/06/25 stream /p4/1/depots/depot/... 'Created by Me'"},
	}}
	p, _ := newTestPerforce(t, proc)
	br, err := p.Branches()
	if err != nil || !reflect.DeepEqual(br, []Ref{{"master", "//depot@5678"}}) {
		t.Errorf("branches %v %v", br, err)
	}
	tags, _ := p.Tags()
	if want := []Ref{{"0.0.1", "//depot@0.0.1"}, {"0.0.2", "//depot@0.0.2"}}; !reflect.DeepEqual(tags, want) {
		t.Errorf("tags %v", tags)
	}
	info, err := p.ComposerInformation("//depot@0.0.1")
	if m := phpjson.ToMap(info); err != nil || m == nil || m.Str("name") != "test/perforce" {
		t.Errorf("info %v %v", info, err)
	}
	if ok, _ := p.CheckStream(); !ok || !p.IsStream() {
		t.Error("stream depot not detected")
	}
	proc.assertComplete()
}

func TestPerforceIsLoggedInErrors(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"p4", "-u", "user", "-p", "port", "login", "-s"}, ret: 1, stderr: "User user doesn't exist."},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "login", "-s"}, ret: 127, stderr: "sh: 1: exec: p4: not found"},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "login", "-s"}, ret: 1, stderr: "Perforce password (P4PASSWD) invalid or unset."},
	}}
	p, _ := newTestPerforce(t, proc)
	if _, err := p.isLoggedIn(); err == nil || err.Error() != "Invalid user name: user" {
		t.Errorf("1: %v", err)
	}
	if _, err := p.isLoggedIn(); err == nil || !strings.HasPrefix(err.Error(), "p4 command not found in path: ") {
		t.Errorf("2: %v", err)
	}
	if ok, err := p.isLoggedIn(); ok || err != nil {
		t.Errorf("3: %v %v", ok, err)
	}
}

func TestPerforceSyncCodeBaseAndCleanup(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"p4", "-u", "user", "-c", "composer_perforce_TEST_depot", "-p", "port", "sync", "-f", "@label"}},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "client", "-d", "composer_perforce_TEST_depot"}},
	}}
	p, _ := newTestPerforce(t, proc)
	label := "label"
	if err := p.SyncCodeBase(&label); err != nil {
		t.Fatal(err)
	}
	if err := p.WriteP4ClientSpec(); err != nil {
		t.Fatal(err)
	}
	spec, _ := os.ReadFile(p.ClientSpec())
	for _, want := range []string{"Client: composer_perforce_TEST_depot\n\n", "Owner:  user\n\n", "View:  //depot/...  //composer_perforce_TEST_depot/... \n"} {
		if !strings.Contains(string(spec), want) {
			t.Errorf("spec lacks %q:\n%s", want, spec)
		}
	}
	if err := p.CleanupClientSpec(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p.ClientSpec()); err == nil {
		t.Error("client spec not removed")
	}
	proc.assertComplete()
}

func TestPerforceValidPorts(t *testing.T) {
	for _, port := range []string{"1666", "perforce:1666", "ssl:perforce.example.com:1666", "tcp6:[2001:db8::1]:1666", "perforce.example.com"} {
		if !PerforceIsValidPort(port) {
			t.Errorf("%s rejected", port)
		}
	}
	for _, port := range []string{"rsh:touch /tmp/pwned", " RSH : x", "jsh:x", "perforce:1666 extra", "-p"} {
		if PerforceIsValidPort(port) {
			t.Errorf("%s accepted", port)
		}
	}
	env, _, _ := testEnv(t, false, nil, nil, nil)
	if _, err := NewPerforce(nil, "rsh:touch /tmp/pwned", t.TempDir(), env); err == nil {
		t.Error("command executing port accepted")
	} else if _, ok := err.(*SecurityError); !ok {
		t.Errorf("err %T", err)
	}
	if PerforceCheckServerExists("rsh:touch /tmp/pwned", env.Process) {
		t.Error("rsh port probed")
	}
}

func TestPerforceDriverInitializeAndHasComposerFile(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"p4", "-u", "user", "-p", "port", "login", "-s"}},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "depots"}, stdout: "Depot depot 2013/06/25 local depot/... 'Created by Me'"},
		{cmd: []string{"p4", "-u", "user", "-c", "composer_perforce_TEST_depot", "-p", "port", "print", "//depot//composer.json"}, stdout: ""},
		{cmd: []string{"p4", "-u", "user", "-p", "port", "client", "-d", "composer_perforce_TEST_depot"}},
	}}
	env, _, _ := testEnv(t, false, nil, nil, proc)
	rc := repoConfig("port", "depot", "depot", "p4user", "user", "unique_perforce_client_name", "TEST")
	d := NewPerforceDriver(rc, env)
	if err := d.Initialize(); err != nil {
		t.Fatal(err)
	}
	if ok, err := d.HasComposerFile(""); ok || err != nil {
		t.Errorf("has composer file %v %v", ok, err)
	}
	src := d.Source("//depot@1.0.0")
	if src.Str("p4user") != "user" || src.Str("url") != "port" || src.Str("type") != "perforce" {
		t.Errorf("source %v", src)
	}
	d.Cleanup()
	proc.assertComplete()
}
