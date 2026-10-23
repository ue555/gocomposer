package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScriptAliasInput(t *testing.T) {
	cases := map[string][]string{
		"a 'b c'":             {"test", "a", "b c"},
		"--foo=bar 'x y'":     {"test", "--foo=bar", "x y"},
		"--foo='a b'":         {"test", "--foo=a b"},
		"test":                {"-v", "test"},
		"":                    {"test"},
		"-- --foo 'it'\\''s'": {"test", "--", "--foo", "it's"},
	}
	for want, argv := range cases {
		if got := scriptAliasInput(argv); got != want {
			t.Errorf("scriptAliasInput(%q) = %q, want %q", argv, got, want)
		}
	}
}

func TestPregQuote(t *testing.T) {
	if got := pregQuote("fund.o/x-y#z"); got != `fund\.o/x\-y\#z` {
		t.Errorf("pregQuote = %q", got)
	}
}

func TestFilterValidateURL(t *testing.T) {
	for u, want := range map[string]bool{
		"https://example.org/x": true,
		"example.com":           false,
		"http://":               false,
		"mailto:a@example.org":  true,
		"https://exa mple.org":  false,
	} {
		if got := filterValidateURL(u); got != want {
			t.Errorf("filterValidateURL(%q) = %v", u, got)
		}
	}
}

func TestPackagesJSONURL(t *testing.T) {
	for in, want := range map[string]string{
		"http://127.0.0.1:8899":         "http://127.0.0.1:8899/packages.json",
		"https://packagist.org":         "https://repo.packagist.org/packages.json",
		"https://example.org/repo.json": "https://example.org/repo.json",
		"https://example.org/sub/":      "https://example.org/sub/packages.json",
		"repo.example.org/composer":     "http://repo.example.org/composer/packages.json",
	} {
		if got := packagesJSONURL(in); got != want {
			t.Errorf("packagesJSONURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNewProxyItem(t *testing.T) {
	p, err := newProxyItem("http://user:pass@proxy.example:3128", "http_proxy")
	if err != "" || p.safeURL != "http://***:***@proxy.example:3128" || p.url != "http://proxy.example:3128" {
		t.Fatalf("got %+v %q", p, err)
	}
	if _, err := newProxyItem("http://proxy.example:0", "http_proxy"); err != "port 0 is reserved in http_proxy" {
		t.Fatalf("port 0: %q", err)
	}
	if p, _ := newProxyItem("https://proxy.example", "HTTPS_PROXY"); p.url != "https://proxy.example:443" {
		t.Fatalf("https default port: %+v", p)
	}
}

func TestPhpVersionCompare(t *testing.T) {
	if phpVersionCompare("2.24.0", "2.47.3") >= 0 || phpVersionCompare("2.24.0", "2.24") <= 0 || phpVersionCompare("1.2", "1.2") != 0 {
		t.Fatal("phpVersionCompare")
	}
}

func TestIsScriptEventConstant(t *testing.T) {
	if !isScriptEventConstant("pre_install_cmd") || !isScriptEventConstant("post-install-cmd") || isScriptEventConstant("post-package-install") {
		t.Fatal("isScriptEventConstant")
	}
}

func TestCacheGCTTL(t *testing.T) {
	root := t.TempDir()
	old := filepath.Join(root, "a", "old.zip")
	fresh := filepath.Join(root, "fresh.zip")
	dot := filepath.Join(root, ".keep")
	for _, f := range []string{old, fresh, dot} {
		_ = os.MkdirAll(filepath.Dir(f), 0o755)
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	past := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(old, past, past)
	_ = os.Chtimes(dot, past, past)
	cacheGC(root+"/", 3600, 1<<30)
	if _, err := os.Stat(old); err == nil {
		t.Error("expired file kept")
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Error("fresh file removed")
	}
	if _, err := os.Stat(dot); err != nil {
		t.Error("dot file removed (Finder ignores dot files)")
	}
	vcs := t.TempDir()
	oldRepo := filepath.Join(vcs, "old.git")
	newRepo := filepath.Join(vcs, "new.git")
	_ = os.Mkdir(oldRepo, 0o755)
	_ = os.Mkdir(newRepo, 0o755)
	_ = os.Chtimes(oldRepo, past, past)
	cacheGCVcs(vcs+"/", 3600)
	if _, err := os.Stat(oldRepo); err == nil {
		t.Error("expired vcs dir kept")
	}
	if _, err := os.Stat(newRepo); err != nil {
		t.Error("fresh vcs dir removed")
	}
}
