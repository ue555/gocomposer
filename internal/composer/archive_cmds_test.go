package composer

import "testing"

func TestVersionLess(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{{"16.02", "21.01", true}, {"25.01", "21.01", false}, {"21.01", "21.01", false}, {"1.22.3", "1.22.4", true}, {"1.22.8", "1.22.4", false}, {"1.22", "1.22.4", true}} {
		if got := versionLess(c.a, c.b); got != c.want {
			t.Errorf("versionLess(%s, %s) = %v", c.a, c.b, got)
		}
	}
}

func TestSevenZipVersionRegex(t *testing.T) {
	m := re7ZipVersion.FindStringSubmatch("\n7-Zip 25.01 (arm64) : Copyright (c) 1999-2025 Igor Pavlov")
	if m == nil || m[1] != "25.01" {
		t.Fatalf("got %v", m)
	}
	m = re7ZipVersion.FindStringSubmatch("\n7-Zip [64] 16.02 : Copyright (c) 1999-2016 Igor Pavlov")
	if m == nil || m[1] != "16.02" {
		t.Fatalf("got %v", m)
	}
}
