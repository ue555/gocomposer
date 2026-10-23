package cli

import (
	"testing"

	"gocomposer/internal/pkg"
)

// ShowCommand::getUpdateStatus (Composer f83588392df43d4fd027d84d26de459946ef4e80).
func TestUpdateStatus(t *testing.T) {
	cases := []struct{ installed, latest, want string }{
		{"1.0.0.0", "1.0.0.0", "up-to-date"},
		{"1.0.0.0", "1.2.0.0", "semver-safe-update"},
		{"1.0.0.0", "2.0.0.0", "update-possible"},
		{"0.1.0.0", "0.1.1.0", "semver-safe-update"},
		{"0.1.0.0", "0.2.0.0", "update-possible"},
	}
	for _, c := range cases {
		p := pkg.NewPackage("a/b", c.installed, c.installed)
		l := pkg.NewPackage("a/b", c.latest, c.latest)
		if got := updateStatus(l, p); got != c.want {
			t.Errorf("updateStatus(%s -> %s) = %s, want %s", c.installed, c.latest, got, c.want)
		}
	}
	if s := updateStatusToVersionStyle("semver-safe-update"); s != "highlight" {
		t.Errorf("style %q", s)
	}
	if s := updateStatusSigns.Replace("cooldown"); s != "c" {
		t.Errorf("sign %q", s)
	}
}
