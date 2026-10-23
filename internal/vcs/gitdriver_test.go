package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Repository/Vcs/GitDriverTest.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"errors"
	"reflect"
	"testing"
)

func newGitDriver(t *testing.T, u string, proc *procMock) *GitDriver {
	t.Helper()
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewGitDriver(repoConfig(u), env)
	d.repoDir = env.Config.GetString("home") // setRepoDir($driver, $this->home)
	return d
}

func TestGitDriverGetRootIdentifierFromRemoteLocalRepository(t *testing.T) {
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"git", "branch", "--no-color"}, stdout: "* main\n  2.2\n  1.10"},
	}}
	home := t.TempDir()
	env, _, _ := testEnv(t, false, nil, nil, proc)
	d := NewGitDriver(repoConfig(home), env)
	d.repoDir = home
	if root, err := d.RootIdentifier(); root != "main" || err != nil {
		t.Errorf("root = %q (%v)", root, err)
	}
	proc.assertComplete()
}

func TestGitDriverGetRootIdentifierFromRemote(t *testing.T) {
	t.Setenv("COMPOSER_DISABLE_NETWORK", "")
	out := `* remote origin
  Fetch URL: https://example.org/acme.git
  Push  URL: https://example.org/acme.git
  HEAD branch: main
  Remote branches:
    1.10                       tracked
    2.2                        tracked
    main                       tracked`
	proc := &procMock{expects: []procExpect{
		{cmd: []string{"git", "remote", "-v"}},
		{cmd: []string{"git", "remote", "set-url", "origin", "--", "https://example.org/acme.git"}},
		{cmd: []string{"git", "remote", "show", "origin"}, stdout: out},
		{cmd: []string{"git", "remote", "set-url", "origin", "--", "https://example.org/acme.git"}},
	}}
	d := newGitDriver(t, "https://example.org/acme.git", proc)
	if root, err := d.RootIdentifier(); root != "main" || err != nil {
		t.Errorf("root = %q (%v)", root, err)
	}
	proc.assertComplete()
}

func TestGitDriverGetRootIdentifierFromLocalWithNetworkDisabled(t *testing.T) {
	t.Setenv("COMPOSER_DISABLE_NETWORK", "1")
	proc := &procMock{expects: []procExpect{
		{cmd: []string{"git", "branch", "--no-color"}, stdout: "* main\n  2.2\n  1.10"},
	}}
	d := newGitDriver(t, "https://example.org/acme.git", proc)
	if root, err := d.RootIdentifier(); root != "main" || err != nil {
		t.Errorf("root = %q (%v)", root, err)
	}
	proc.assertComplete()
}

func TestGitDriverGetBranchesFilterInvalidBranchNames(t *testing.T) {
	// Branches starting with a - character are not valid git branches names
	// Still assert that they get filtered to prevent issues later on
	proc := &procMock{expects: []procExpect{
		{cmd: []string{"git", "branch", "--no-color", "--no-abbrev", "-v"}, stdout: "* main 089681446ba44d6d9004350192486f2ceb4eaa06 commit\n  2.2  12681446ba44d6d9004350192486f2ceb4eaa06 commit\n  -h   089681446ba44d6d9004350192486f2ceb4eaa06 commit"},
	}}
	d := newGitDriver(t, "https://example.org/acme.git", proc)
	branches, err := d.Branches()
	want := []Ref{{"main", "089681446ba44d6d9004350192486f2ceb4eaa06"}, {"2.2", "12681446ba44d6d9004350192486f2ceb4eaa06"}}
	if err != nil || !reflect.DeepEqual(branches, want) {
		t.Errorf("branches = %v (%v), want %v", branches, err, want)
	}
}

func TestGitDriverFileGetContentInvalidIdentifier(t *testing.T) {
	d := newGitDriver(t, "https://example.org/acme.git", &procMock{})
	if content, ok, err := d.FileContent("file.txt", "h"); ok || content != "" || err != nil {
		t.Errorf("FileContent(h) = %q, %v, %v; want null", content, ok, err)
	}
	_, _, err := d.FileContent("file.txt", "-h")
	var re *RuntimeError
	if !errors.As(err, &re) {
		t.Errorf("FileContent(-h) error = %#v, want RuntimeException", err)
	}
}

func TestGitDriverGetChangeDateInvalidIdentifier(t *testing.T) {
	d := newGitDriver(t, "https://example.org/acme.git", &procMock{})
	_, _, err := d.ChangeDate("-n1 --format=%at HEAD")
	var re *RuntimeError
	if !errors.As(err, &re) {
		t.Errorf("ChangeDate error = %#v, want RuntimeException", err)
	}
}
