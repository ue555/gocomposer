package composer

import (
	"strings"
	"testing"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

func testPkg(t *testing.T, name string, requires map[string]string) *pkg.Package {
	t.Helper()
	m := phpjson.NewMap()
	m.Set("name", name)
	m.Set("version", "1.0.0")
	if len(requires) > 0 {
		r := phpjson.NewMap()
		for k, v := range requires {
			r.Set(k, v)
		}
		m.Set("require", r)
	}
	p, err := (&pkg.Loader{}).Load(m)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func names(list []*pkg.Package) string {
	var out []string
	for _, p := range list {
		out = append(out, p.Name)
	}
	return strings.Join(out, ",")
}

// PackageSorter::sortPackages: dependencies first, weights (composer/installers
// -10000) before everything, natural name order otherwise.
func TestSortPackagesWeighted(t *testing.T) {
	a := testPkg(t, "vendor/a", map[string]string{"vendor/b": "*"})
	b := testPkg(t, "vendor/b", nil)
	inst := testPkg(t, "composer/installers", nil)
	c := testPkg(t, "vendor/c", nil)
	got := names(sortPackagesWeighted([]*pkg.Package{a, c, b, inst}, map[string]int{"composer/installers": -10000}))
	if want := "composer/installers,vendor/b,vendor/a,vendor/c"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

// RepositoryUtils::filterRequiredPackages follows the requirements of the
// root package recursively.
func TestFilterRequiredPackages(t *testing.T) {
	root := testPkg(t, "root/pkg", map[string]string{"vendor/a": "*"})
	a := testPkg(t, "vendor/a", map[string]string{"vendor/b": "*"})
	b := testPkg(t, "vendor/b", nil)
	c := testPkg(t, "vendor/c", nil)
	got := names(filterRequiredPackages([]*pkg.Package{c, b, a}, root, true, nil))
	if want := "vendor/a,vendor/b"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

// PluginManager::isPluginAllowed with allow-plugins patterns in
// non-interactive mode.
func TestIsPluginAllowedRules(t *testing.T) {
	cfg := config.New(false, t.TempDir())
	allow := phpjson.NewMap()
	allow.Set("acme/*", true)
	allow.Set("acme/blocked", false)
	allow.Set("other/off", false)
	wrap := phpjson.NewMap()
	c2 := phpjson.NewMap()
	c2.Set("allow-plugins", allow)
	wrap.Set("config", c2)
	cfg.Merge(wrap)
	io := cio.New()
	io.Interactive = false
	c := &Composer{IO: io, Config: cfg}
	pm := &PluginManager{c: c}
	pm.parseAllowedPlugins()
	cases := []struct {
		name     string
		optional bool
		allowed  bool
		blocked  bool
	}{
		{"acme/blocked", false, true, false}, // the first matching rule wins
		{"ACME/Thing", false, true, false},
		{"other/off", false, false, false},
		{"third/party", true, false, false},
		{"third/party", false, false, true},
		{"composer/package-versions-deprecated", false, false, false},
	}
	for _, tc := range cases {
		ok, err := pm.isPluginAllowed(tc.name, false, tc.optional, true)
		if tc.blocked {
			if _, isBlocked := err.(*PluginBlockedError); !isBlocked {
				t.Fatalf("%s: expected PluginBlockedError, got %v", tc.name, err)
			}
			if !strings.Contains(err.Error(), `You can run "composer config --no-plugins allow-plugins.third/party [true|false]"`) {
				t.Fatalf("unexpected message %q", err.Error())
			}
			continue
		}
		if err != nil || ok != tc.allowed {
			t.Fatalf("%s: got %v %v, want %v", tc.name, ok, err, tc.allowed)
		}
	}
}

func TestPluginDownloadCheck(t *testing.T) {
	p := testPkg(t, "vendor/plugin", nil)
	p.SetType("composer-plugin")
	err := pluginDownloadCheck(p)
	if err == nil || err.Error() != "Error while installing vendor/plugin, composer-plugin packages should have a class defined in their extra key to be usable." {
		t.Fatalf("unexpected %v", err)
	}
	e := phpjson.NewMap()
	e.Set("class", []phpjson.Value{"A\\Plugin", "B\\Plugin"})
	p.SetExtra(e)
	if err := pluginDownloadCheck(p); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(pluginClasses(p), ","); got != `A\Plugin,B\Plugin` {
		t.Fatalf("classes %s", got)
	}
}
