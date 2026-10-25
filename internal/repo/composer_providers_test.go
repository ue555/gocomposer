package repo

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
)

func sha(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

// providerServer serves a v1 provider repository (providers-url +
// provider-includes) like Packagist did; tamper serves a provider file whose
// contents do not match the listed hash.
func providerServer(t *testing.T, tamper bool) (*httptest.Server, *int) {
	lib := `{"packages": {"v/lib": {"1.0.0": {"name": "v/lib", "version": "1.0.0", "uid": 1}, "1.1.0": {"name": "v/lib", "version": "1.1.0", "uid": 2}, "dup": {"name": "v/lib", "version": "1.1.0", "uid": 2}, "x": {"name": "v/other", "version": "9.0.0"}}}}`
	app := `{"packages": {"v/app": {"1.0.0": {"name": "v/app", "version": "1.0.0", "warning": "ignored"}}}, "warning": "provider file warning"}`
	libHash := sha(lib)
	if tamper {
		libHash = sha("tampered")
	}
	inc := `{"providers": {"v/lib": {"sha256": "` + libHash + `"}, "v/app": {"sha256": "` + sha(app) + `"}}}`
	root := `{"packages": [], "providers-url": "/p/%package%$%hash%.json", "provider-includes": {"p/inc$%hash%.json": {"sha256": "` + sha(inc) + `"}}}`
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch r.URL.Path {
		case "/packages.json":
			w.Write([]byte(root))
		case "/p/inc$" + sha(inc) + ".json":
			w.Write([]byte(inc))
		case "/p/v/lib$" + libHash + ".json":
			w.Write([]byte(lib))
		case "/p/v/app$" + sha(app) + ".json":
			w.Write([]byte(app))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

func newTestComposerRepo(t *testing.T, url string) (*Composer, *bytes.Buffer) {
	t.Helper()
	home := t.TempDir()
	c := config.New(false, home)
	data := phpjson.NewMap()
	data.Set("secure-http", false)
	wrap := phpjson.NewMap()
	wrap.Set("config", data)
	c.Merge(wrap)
	dl, err := httpdl.New(c)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	dl.IO = &cio.IO{Out: &out, Err: &out, Verbosity: cio.Normal}
	httpdl.ResetResponseWarnings()
	rc := phpjson.NewMap()
	rc.Set("url", url)
	r, err := NewComposer(rc, home+"/repo", false, dl)
	if err != nil {
		t.Fatal(err)
	}
	return r, &out
}

func TestProviderRepository(t *testing.T) {
	srv, hits := providerServer(t, false)
	r, out := newTestComposerRepo(t, srv.URL)
	names, err := r.PackageNames("")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(names, ",") != "v/lib,v/app" {
		t.Fatalf("names %v", names)
	}
	res, err := r.LoadPackages([]NameConstraint{{Name: "v/lib"}, {Name: "php"}}, map[string]int{"stable": 0}, map[string]int{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, p := range res.Packages {
		got = append(got, p.Name+"@"+p.PrettyVersion())
	}
	// uid 2 is loaded once, the v/other entry of the file is ignored
	if strings.Join(got, ",") != "v/lib@1.0.0,v/lib@1.1.0" || strings.Join(res.NamesFound, ",") != "v/lib" {
		t.Fatalf("loaded %v names %v", got, res.NamesFound)
	}
	ps, err := r.FindPackages("v/app", nil)
	if err != nil || len(ps) != 1 {
		t.Fatalf("find %v %v", ps, err)
	}
	if !strings.Contains(out.String(), "Warning from "+srv.URL+": provider file warning") {
		t.Fatalf("warning not shown: %q", out.String())
	}
	if _, err := r.Packages(); err == nil || !strings.Contains(err.Error(), "can not load the complete list of packages") {
		t.Fatalf("getPackages: %v", err)
	}
	// a second repository instance reads the hashed files from the cache
	before := *hits
	r2, _ := newTestComposerRepo(t, srv.URL)
	r2.cache = r.cache
	if _, err := r2.FindPackages("v/lib", nil); err != nil {
		t.Fatal(err)
	}
	if *hits-before != 1 { // only packages.json
		t.Fatalf("expected cached provider files, got %d requests", *hits-before)
	}
}

func TestProviderRepositoryHashMismatch(t *testing.T) {
	srv, hits := providerServer(t, true)
	r, _ := newTestComposerRepo(t, srv.URL)
	before := *hits
	_, err := r.FindPackages("v/lib", nil)
	if _, ok := err.(*RepositorySecurityError); !ok || !strings.Contains(err.Error(), "do not match its signature") || !strings.Contains(err.Error(), "/p/v/lib%24") {
		t.Fatalf("err %v", err)
	}
	// packages.json + include + three attempts of the provider file
	if *hits-before != 5 {
		t.Fatalf("requests %d", *hits-before)
	}
}
