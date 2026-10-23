package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Util/GitTest.php: testSyncMirrorSanitizesUrl*).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
//
// The runCommand() tests are not ported: they drive the deprecated
// callable based Git::runCommand(), gocomposer only has the %url%
// template based runCommands().

import (
	"path/filepath"
	"testing"

	"gocomposer/internal/phpjson"
)

func syncMirrorEnv(t *testing.T, proc *procMock) *Git {
	t.Helper()
	t.Setenv("COMPOSER_DISABLE_NETWORK", "")
	env, _, _ := testEnv(t, false, map[string]phpjson.Value{
		"github-domains":   []phpjson.Value{"github.com"},
		"gitlab-domains":   []phpjson.Value{"gitlab.com"},
		"github-protocols": []phpjson.Value{"https"},
	}, nil, proc)
	return env.Git()
}

func TestGitSyncMirrorSanitizesURLAfterInitialClone(t *testing.T) {
	for _, u := range []string{
		"https://user:secret@example.com/repo.git",
		// a bare token in the user slot has no colon to key off
		"https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB@example.com/repo.git",
		// a password containing @ must not leave its tail behind
		"https://user:sec@ret@example.com/repo.git",
	} {
		t.Run(u, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "composer-test-nonexistent")
			proc := &procMock{strict: true, expects: []procExpect{
				{cmd: []string{"git", "clone", "--mirror", "--", u, dir}},
				{cmd: []string{"git", "remote", "-v"}},
				{cmd: []string{"git", "remote", "set-url", "origin", "--", "https://example.com/repo.git"}},
			}}
			g := syncMirrorEnv(t, proc)
			ok, err := g.SyncMirror(u, dir)
			if !ok || err != nil {
				t.Errorf("SyncMirror = %v, %v; want true", ok, err)
			}
			proc.assertComplete()
		})
	}
}

func TestGitSyncMirrorSanitizesURLEvenAfterFailedUpdate(t *testing.T) {
	dir := t.TempDir()
	proc := &procMock{strict: true, expects: []procExpect{
		{cmd: []string{"git", "rev-parse", "--git-dir"}, stdout: ".\n"},
		{cmd: []string{"git", "remote", "-v"}},
		{cmd: []string{"git", "remote", "set-url", "origin", "--", "https://user:secret@example.com/repo.git"}},
		{cmd: []string{"git", "remote", "update", "--prune", "origin"}, ret: 1},
		{cmd: []string{"git", "--version"}},
		{cmd: []string{"git", "remote", "-v"}},
		{cmd: []string{"git", "remote", "set-url", "origin", "--", "https://example.com/repo.git"}},
	}}
	g := syncMirrorEnv(t, proc)
	ok, _ := g.SyncMirror("https://user:secret@example.com/repo.git", dir)
	if ok {
		t.Errorf("SyncMirror = true, want false")
	}
	proc.assertComplete()
}
