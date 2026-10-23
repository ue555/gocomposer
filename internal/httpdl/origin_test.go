package httpdl

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (tests/Composer/Test/Util/UrlTest.php: getOriginProvider) plus the
// special cases of Url::getOrigin.
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"testing"

	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
)

func TestOrigin(t *testing.T) {
	cases := []struct {
		want, url string
		domains   []phpjson.Value
	}{
		// a host that is only a shorter prefix of a configured domain must not resolve to it
		{"gitlab.example.co", "https://gitlab.example.co/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.com"}},
		{"gitlab.example.co", "https://gitlab.example.co/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.co.uk/gitlab"}},
		{"gitlab.example.com", "https://gitlab.example.com/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.com"}},
		{"gitlab.example.com/gitlab", "https://gitlab.example.com/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.com/gitlab"}},
		// a configured domain may spell out a port the URL omits
		{"gitlab.example.com:443", "https://gitlab.example.com/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.com:443"}},
		{"gitlab.example.com:443/gitlab", "https://gitlab.example.com/foo/bar/repository/archive.zip", []phpjson.Value{"gitlab.example.com:443/gitlab"}},
		// Url::getOrigin special cases
		{"github.com", "https://api.github.com/repos/a/b", nil},
		{"codeload.github.com", "https://codeload.github.com/a/b/zip/x", nil},
		{"packagist.org", "https://repo.packagist.org/p2/a/b.json", nil},
		{"127.0.0.1:8900", "http://127.0.0.1:8900/x", nil},
	}
	for _, c := range cases {
		cfg := config.New(false, t.TempDir())
		data := phpjson.NewMap()
		if c.domains != nil {
			data.Set("gitlab-domains", c.domains)
		}
		wrap := phpjson.NewMap()
		wrap.Set("config", data)
		cfg.Merge(wrap)
		d, err := New(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got := d.Origin(c.url); got != c.want {
			t.Errorf("Origin(%q) with %v = %q, want %q", c.url, c.domains, got, c.want)
		}
	}
}
