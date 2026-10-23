package cli

// Cases from tests/Composer/Test/Package/Version/VersionBumperTest.php
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"testing"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

func TestBumpRequirement(t *testing.T) {
	cases := []struct{ name, req, version, expected, alias string }{
		{"upgrade caret", "^1.0", "1.2.1", "^1.2.1", ""},
		{"upgrade caret with v", "^v1.0", "1.2.1", "^1.2.1", ""},
		{"skip trailing .0s", "^1.0", "1.0.0", "^1.0", ""},
		{"skip trailing .0s/2", "^1.2", "1.2.0", "^1.2", ""},
		{"preserve major.minor.patch format when installed minor is 0", "^1.0.0", "1.2.0", "^1.2.0", ""},
		{"preserve major.minor.patch format when installed minor is 1", "^1.0.0", "1.2.1", "^1.2.1", ""},
		{"preserve multi constraints", "^1.2 || ^2.3", "1.3.2", "^1.3.2 || ^2.3", ""},
		{"preserve multi constraints/2", "^1.2 || ^2.3", "2.4.0", "^1.2 || ^2.4", ""},
		{"preserve multi constraints/3", "^1.2 || ^2.3 || ^2", "2.4.0", "^1.2 || ^2.4 || ^2.4", ""},
		{"preserve multi constraints/4", "^1.2 || ^2.3.3 || ^2", "2.4.0", "^1.2 || ^2.4.0 || ^2.4", ""},
		{"@dev is preserved", "^3@dev", "3.2.x-dev", "^3.2@dev", ""},
		{"non-stable versions abort upgrades", "~2", "2.1-beta.1", "~2", ""},
		{"dev reqs are skipped", "dev-main", "dev-foo", "dev-main", ""},
		{"dev version does not upgrade", "^3.2", "dev-main", "^3.2", ""},
		{"upgrade dev version if aliased", "^3.2", "dev-main", "^3.3", "3.3.x-dev"},
		{"upgrade major wildcard to caret", "2.*", "2.4.0", "^2.4", ""},
		{"upgrade major wildcard to caret with v", "v2.*", "2.4.0", "^2.4", ""},
		{"upgrade major wildcard as x to caret", "2.x", "2.4.0", "^2.4", ""},
		{"upgrade major wildcard as x to caret/2", "2.x.x", "2.4.0", "^2.4.0", ""},
		{"leave minor wildcard alone", "2.4.*", "2.4.3", "2.4.*", ""},
		{"leave patch wildcard alone", "2.4.3.*", "2.4.3.2", "2.4.3.*", ""},
		{"leave single tilde alone", "~2", "2.4.3", "~2", ""},
		{"upgrade tilde to caret when compatible", "~2.2", "2.4.3", "^2.4.3", ""},
		{"upgrade patch-only-tilde, longer version", "~2.2.3", "2.2.6.2", "~2.2.6", ""},
		{"upgrade patch-only-tilde", "~2.2.3", "2.2.6", "~2.2.6", ""},
		{"upgrade patch-only-tilde, also .0s", "~2.0.0", "2.0.0", "~2.0.0", ""},
		{"upgrade patch-only-tilde - with year based versions", "~2025.1.561", "2025.1.583", "~2025.1.583", ""},
		{"upgrade 4 bits tilde", "~2.2.3.1", "2.2.4", "~2.2.4.0", ""},
		{"upgrade 4 bits tilde/2", "~2.2.3.1", "2.2.4.0", "~2.2.4.0", ""},
		{"upgrade 4 bits tilde/3", "~2.2.3.1", "2.2.4.5", "~2.2.4.5", ""},
		{"upgrade bigger-or-eq to latest", ">=3.0", "3.4.5", ">=3.4.5", ""},
		{"upgrade bigger-or-eq to latest with v", ">=v3.0", "3.4.5", ">=3.4.5", ""},
		{"upgrade bigger-or-eq with a multi-digit major", ">=10.0", "10.5.2", ">=10.5.2", ""},
		{"upgrade bigger-or-eq with a multi-digit major in an or-constraint", "^9.0 || >=10", "12.1.3", "^9.0 || >=12.1.3", ""},
		{"upgrade bigger-or-eq with a calendar version", ">=2025.1", "2025.3.0", ">=2025.3", ""},
		{"leave a bigger-or-eq above the installed version untouched", "^8.5 || >=10", "8.5.30", "^8.5.30 || >=10", ""},
		{"leave a single-digit bigger-or-eq above the installed version untouched", "^1.0 || >=3", "1.5.0", "^1.5 || >=3", ""},
		{"upgrade bigger-or-eq on a dev version if aliased", ">=3.2", "dev-main", ">=3.3", "3.3.x-dev"},
		{"leave a bigger-or-eq below the installed version in another range untouched", ">=10.0 <11.0 || >=12.0", "12.1.0", ">=10.0 <11.0 || >=12.1", ""},
		{"leave a single-digit bigger-or-eq below the installed version in another range untouched", ">=2.0 <3.0 || >=5.0", "5.1.0", ">=2.0 <3.0 || >=5.1", ""},
		{"upgrade bigger-or-eq with a multi-digit major and an upper bound", ">=10.0 <11", "10.5.2", ">=10.5.2 <11", ""},
		{"upgrade bigger-or-eq with a multi-digit major and a comma upper bound", ">=10.0,<11", "10.5.2", ">=10.5.2,<11", ""},
		{"upgrade bigger-or-eq with a multi-digit major and a stability flag", ">=10.0@dev", "10.5.2", ">=10.5.2@dev", ""},
		{"upgrade caret and leave a bigger-or-eq above the installed version untouched", "^10.0 || >=11", "10.5.0", "^10.5 || >=11", ""},
		{"upgrade bigger-or-eq in a single-pipe or-constraint", ">=10.0 <11|>=12", "12.1.0", ">=10.0 <11|>=12.1", ""},
		{"leave bigger-than untouched", ">2.2.3", "2.2.6", ">2.2.3", ""},
		{"skip pre-stable releases", "^0.3 || ^0.4", "0.4.3", "^0.3 || ^0.4.3", ""},
		{"upgrade full wildcard to bigger-or-eq", "*", "1.2.3", ">=1.2.3", ""},
	}
	for _, c := range cases {
		p := pkg.NewPackage("foo/bar", semver.MustNormalize(c.version), c.version)
		if c.alias != "" {
			ba := phpjson.NewMap()
			ba.Set(c.version, c.alias)
			extra := phpjson.NewMap()
			extra.Set("branch-alias", ba)
			p.SetExtra(extra)
		}
		cst := semver.MustParseConstraints(c.req)
		if got := bumpRequirement(cst, c.req, p); got != c.expected {
			t.Errorf("%s: bump(%q, %q) = %q, want %q", c.name, c.req, c.version, got, c.expected)
		}
	}
}
