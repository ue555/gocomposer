package composer

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"testing"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/solver"
)

// FileDownloader's cache key generator: <name>/<sha1(key)>.<dist type>.
func TestDistCacheKey(t *testing.T) {
	m := phpjson.NewMap()
	m.Set("name", "acme/log")
	m.Set("version", "1.0.0")
	dist := phpjson.NewMap()
	dist.Set("type", "zip")
	dist.Set("url", "http://example.org/acme-log.zip")
	m.Set("dist", dist)
	p, err := (&pkg.Loader{}).Load(m)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha1.Sum([]byte("custom-key"))
	if got, want := distCacheKey(p, "custom-key"), "acme/log/"+hex.EncodeToString(sum[:])+".zip"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

// http.header options are sent as header lines (string or list).
func TestTransportHeadersOf(t *testing.T) {
	opts := phpjson.NewMap()
	h := phpjson.NewMap()
	h.Set("header", []phpjson.Value{"X-A: 1", "X-B: 2"})
	opts.Set("http", h)
	if got := strings.Join(transportHeadersOf(opts), "|"); got != "X-A: 1|X-B: 2" {
		t.Fatalf("list: %q", got)
	}
	h.Set("header", "X-A: 1\r\nX-B: 2")
	if got := strings.Join(transportHeadersOf(opts), "|"); got != "X-A: 1|X-B: 2" {
		t.Fatalf("string: %q", got)
	}
	if transportHeadersOf(nil) != nil {
		t.Fatal("nil options")
	}
}

// The response header lines of the runtime's Response: status line first,
// then lower-cased names in a stable order.
func TestResponseHeaderLines(t *testing.T) {
	resp := &httpdl.Response{Proto: "HTTP/1.1", Status: "200 OK", StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}, "Etag": {`"x"`}}}
	var got []string
	for _, l := range responseHeaderLines(resp) {
		got = append(got, l.(string))
	}
	if want := `HTTP/1.1 200 OK|content-type: application/json|etag: "x"`; strings.Join(got, "|") != want {
		t.Fatalf("got %q", strings.Join(got, "|"))
	}
}

func TestRepoKindName(t *testing.T) {
	for k, want := range map[repo.Kind]string{repo.KindComposer: "composer", repo.KindPath: "path", repo.KindVcs: "vcs", repo.KindPackage: "package", repo.KindArtifact: "artifact", repo.KindPlatform: "platform", repo.KindLock: "lock", repo.KindRoot: "root", repo.KindInstalled: "installed", repo.KindArray: "array"} {
		if got := repoKindName(k); got != want {
			t.Errorf("%v: got %s, want %s", k, got, want)
		}
	}
}

// The platform requirement filter of a running installer as handed to the
// runtime: booleans, or the token standing for the installer's own list.
func TestPlatformFilterSpec(t *testing.T) {
	if platformFilterSpec(solver.IgnoreAll{}) != true {
		t.Fatal("ignore all")
	}
	if platformFilterSpec(solver.IgnoreNothing{}) != false {
		t.Fatal("ignore nothing")
	}
	list, ok := platformFilterSpec(solver.NewIgnoreList([]string{"ext-foo"})).([]phpjson.Value)
	if !ok || len(list) != 1 || list[0] != baseFilterToken {
		t.Fatalf("ignore list: %#v", list)
	}
}

// customFlow: installers not built on LibraryInstaller run every step
// themselves, LibraryInstaller subclasses only the hooks they replace.
func TestPluginInstallerCustomFlow(t *testing.T) {
	plain := &pluginInstallerInfo{Library: false}
	if !plain.customFlow("download") || !plain.customFlow("cleanup") {
		t.Fatal("non-library installer")
	}
	lib := &pluginInstallerInfo{Library: true, Overrides: []string{"installCode", "prepare"}}
	if lib.customFlow("download") || !lib.customFlow("prepare") || lib.customFlow("cleanup") {
		t.Fatal("library installer")
	}
}
