package httpdl

// Unit tests of the interactive authentication helpers (ported behaviour of
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Util/
// AuthHelper.php, GitHub.php). Copyright (c) Nils Adermann, Jordi Boggiano.
// MIT License. The prompts themselves are compared on a pseudo terminal by
// the tty-auth-* scenarios (tools/compare/scenarios_tty2.go).

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
)

func TestPHPEncoders(t *testing.T) {
	if got := rawurlencode("Composer on h 2024-01-02 0304~"); got != "Composer%20on%20h%202024-01-02%200304~" {
		t.Errorf("rawurlencode = %q", got)
	}
	if got := urlencode("a b~c&d"); got != "a+b%7Ec%26d" {
		t.Errorf("urlencode = %q", got)
	}
}

func TestGitHubRateLimitHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("X-RateLimit-Limit", "60")
	h.Set("X-RateLimit-Remaining", "0")
	h.Set("X-RateLimit-Reset", "1704164645")
	if !GitHubIsRateLimited(h) {
		t.Fatal("rate limit not detected")
	}
	limit, reset := GitHubRateLimit(h)
	if limit != 60 || reset != "2024-01-02 03:04:05" {
		t.Errorf("GitHubRateLimit = %d %q", limit, reset)
	}
	if l, r := GitHubRateLimit(http.Header{}); l != 0 || r != "?" {
		t.Errorf("GitHubRateLimit(empty) = %d %q", l, r)
	}
	h2 := http.Header{}
	h2.Set("X-GitHub-SSO", "required; url=https://github.com/orgs/x/sso?authorization_request=abc")
	if !gitHubRequiresSso(h2) || gitHubSsoURL(h2) != "https://github.com/orgs/x/sso?authorization_request=abc" {
		t.Errorf("SSO header not parsed")
	}
}

// A 401 on a non-forge host asks for username and password, retries with
// them and stores them in auth.json when store-auths is true; the password
// never reaches the console.
func TestPromptBasicAuthAndStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("COMPOSER_HOME", home)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); ok && u == "user" && p == "pa ss" {
			w.Write([]byte(`{"packages": {}}`))
			return
		}
		w.WriteHeader(401)
	}))
	defer srv.Close()
	cfg := config.New(false, t.TempDir())
	cfg.Merge(mapOf("config", mapOf("home", home, "store-auths", true, "secure-http", false)))
	d, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	var errOut bytes.Buffer
	d.IO = &cio.IO{Out: &errOut, Err: &errOut, Verbosity: cio.Normal, Interactive: true, In: strings.NewReader("user\npa ss\n")}
	resp, err := d.Get(srv.URL + "/packages.json")
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Body) != `{"packages": {}}` {
		t.Fatalf("body %q", resp.Body)
	}
	if strings.Contains(errOut.String(), "pa ss") {
		t.Errorf("password written to the console: %q", errOut.String())
	}
	if !strings.Contains(errOut.String(), "Authentication required") {
		t.Errorf("missing prompt: %q", errOut.String())
	}
	b, err := os.ReadFile(filepath.Join(home, "auth.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"password": "pa ss"`) || !strings.Contains(string(b), `"username": "user"`) {
		t.Errorf("auth.json = %s", b)
	}
}

func mapOf(kv ...any) *phpjson.Map {
	m := phpjson.NewMap()
	for i := 0; i+1 < len(kv); i += 2 {
		m.Set(kv[i].(string), kv[i+1])
	}
	return m
}
