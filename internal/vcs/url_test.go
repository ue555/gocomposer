package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Util/UrlTest.php: sanitizeProvider,
// stripCredentialsProvider). Copyright (c) Nils Adermann, Jordi Boggiano.
// MIT License.

import "testing"

func TestSanitize(t *testing.T) {
	cases := [][2]string{
		// empty input safe (callers may pass `$x ?? ''` for nullable URLs)
		{"", ""},
		// with scheme
		{"https://foo:***@example.org/", "https://foo:bar@example.org/"},
		{"https://foo@example.org/", "https://foo@example.org/"},
		{"https://example.org/", "https://example.org/"},
		{"http://10a***:***@example.org", "http://10a8f08e8d7b7b9:foo@example.org"},
		{"https://foo:***@example.org:123/", "https://foo:bar@example.org:123/"},
		{"https://example.org/foo/bar?access_token=***", "https://example.org/foo/bar?access_token=abcdef"},
		{"https://example.org/foo/bar?foo=bar&access_token=***", "https://example.org/foo/bar?foo=bar&access_token=abcdef"},
		{"https://ghp***:***@github.com/acme/repo", "https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB:x-oauth-basic@github.com/acme/repo"},
		{"https://git***:***@github.com/acme/repo", "https://github_pat_1234567890abcdefghijkl_1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVW:x-oauth-basic@github.com/acme/repo"},
		{"http://abc***:***@example.org:123/", "http://abcdefghijkl:bar@example.org:123/"},
		{"https://abc***:***@example.org:123/", "https://abcdefghijklmnop:bar@example.org:123/"},
		// token/long username in the user slot without a password (e.g. https://TOKEN@host)
		{"https://ghp***@github.com/acme/repo", "https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB@github.com/acme/repo"},
		{"https://git***@github.com/acme/repo", "https://github_pat_1234567890abcdefghijkl_1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVW@github.com/acme/repo"},
		{"http://10a***@example.org", "http://10a8f08e8d7b7b9@example.org"},
		{"https://abc***@example.org:123/", "https://abcdefghijklmnop@example.org:123/"},
		// well-known non-secret credential markers are shown verbatim even though they are 12char+
		{"https://x-token-auth:***@bitbucket.org/acme/repo", "https://x-token-auth:secret@bitbucket.org/acme/repo"},
		{"https://gitlab-ci-token:***@gitlab.example.org/", "https://gitlab-ci-token:realtoken@gitlab.example.org/"},
		// without scheme
		{"foo:***@example.org/", "foo:bar@example.org/"},
		{"foo@example.org/", "foo@example.org/"},
		{"example.org/", "example.org/"},
		{"10a***:***@example.org", "10a8f08e8d7b7b9:foo@example.org"},
		{"foo:***@example.org:123/", "foo:bar@example.org:123/"},
		{"example.org/foo/bar?access_token=***", "example.org/foo/bar?access_token=abcdef"},
		{"example.org/foo/bar?foo=bar&access_token=***", "example.org/foo/bar?foo=bar&access_token=abcdef"},
		{"abc***:***@example.org:123/", "abcdefghijkl:bar@example.org:123/"},
		{"abc***:***@example.org:123/", "abcdefghijklmnop:bar@example.org:123/"},
		{"ghp***@github.com/acme/repo", "ghp_1234567890abcdefghijklmnopqrstuvwxyzAB@github.com/acme/repo"},
		{"10a***@example.org", "10a8f08e8d7b7b9@example.org"},
		{"abc***@example.org:123/", "abcdefghijklmnop@example.org:123/"},
		// anywhere in the string, as URLs usually reach sanitize() embedded in a longer message
		{"Failed to execute git clone --mirror -- https://foo:***@example.org/private/repo.git /cache/repo", "Failed to execute git clone --mirror -- https://foo:bar@example.org/private/repo.git /cache/repo"},
		{"tried https://foo:***@example.org/a and https://baz:***@example.com/b", "tried https://foo:bar@example.org/a and https://baz:qux@example.com/b"},
		{"fatal: unable to access 'https://gitlab-ci-token:***@example.org/g/r.git/'", "fatal: unable to access 'https://gitlab-ci-token:realtoken@example.org/g/r.git/'"},
		// passwords/usernames containing @ are masked up to the last @ of the authority
		{"https://foo:***@example.org/repo.git", "https://foo:bar@baz@example.org/repo.git"},
		{"https://use***:***@example.org/repo.git", "https://user@corp.example:realtoken@example.org/repo.git"},
		// empty user slot, e.g. https://:TOKEN@host
		{"https://:***@example.org/", "https://:realtoken@example.org/"},
		// empty password slot, e.g. https://TOKEN:@host
		{"https://ghp***@github.com/acme/repo", "https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB:@github.com/acme/repo"},
		// schemes containing + . -
		{"git+ssh://foo:***@example.org/repo.git", "git+ssh://foo:bar@example.org/repo.git"},
	}
	for _, c := range cases {
		if got := Sanitize(c[1]); got != c[0] {
			t.Errorf("Sanitize(%q) = %q, want %q", c[1], got, c[0])
		}
		// testSanitizeIsIdempotent
		if got := Sanitize(Sanitize(c[1])); got != c[0] {
			t.Errorf("Sanitize(Sanitize(%q)) = %q, want %q", c[1], got, c[0])
		}
	}
}

func TestStripCredentials(t *testing.T) {
	cases := [][2]string{
		{"https://example.org/repo.git", "https://user:pass@example.org/repo.git"},
		// a bare token in the user slot is a credential too, and has no colon to key off
		{"https://github.com/acme/repo.git", "https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB@github.com/acme/repo.git"},
		{"https://example.org/repo.git", "https://user:pa@ss@example.org/repo.git"},
		{"https://example.org:8080/repo.git", "https://user:pass@example.org:8080/repo.git"},
		{"https://example.org/repo.git?ref=a@b", "https://user:pass@example.org/repo.git?ref=a@b"},
		// token in the user slot with an empty password, as Git::runCommands() builds for token-only http-basic auth
		{"https://github.com/acme/repo.git", "https://ghp_1234567890abcdefghijklmnopqrstuvwxyzAB:@github.com/acme/repo.git"},
		{"https://example.org/repo.git", "https://:@example.org/repo.git"},
		// nothing to strip
		{"https://example.org/repo.git", "https://example.org/repo.git"},
		{"https://example.org/@scope/repo.git", "https://example.org/@scope/repo.git"},
		{"git@example.org:acme/repo.git", "git@example.org:acme/repo.git"},
		// ssh login name is not a credential and must be preserved
		{"ssh://git@example.org/repo.git", "ssh://git@example.org/repo.git"},
		{"ssh://gogs@git.int.example/group/pkg.git", "ssh://gogs@git.int.example/group/pkg.git"},
		{"ssh://git.int.example/group/pkg.git", "ssh://gogs:secret@git.int.example/group/pkg.git"},
		{"ssh://gogs@git.int.example/group/pkg.git", "ssh://gogs:@git.int.example/group/pkg.git"},
		{"git+ssh://gogs@git.int.example/group/pkg.git", "git+ssh://gogs@git.int.example/group/pkg.git"},
		{"git+ssh://git.int.example/group/pkg.git", "git+ssh://gogs:secret@git.int.example/group/pkg.git"},
	}
	for _, c := range cases {
		if got := StripCredentials(c[1]); got != c[0] {
			t.Errorf("StripCredentials(%q) = %q, want %q", c[1], got, c[0])
		}
	}
}
