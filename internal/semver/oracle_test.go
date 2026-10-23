package semver

import (
	"encoding/json"
	"os"
	"testing"
)

// testdata/oracle.json is produced by oracle/scripts/semver-dump.php running the
// composer/semver copy locked by Composer f83588392df43d4fd027d84d26de459946ef4e80
// on official PHP 8.4.27 (see docs/validation). Regenerate with
// scripts/oracle-semver.sh; never edit it by hand.

type corpus struct {
	Versions      []string    `json:"versions"`
	Compare       [][2]string `json:"compare"`
	Constraints   []string    `json:"constraints"`
	MatchVersions []string    `json:"matchVersions"`
	Subset        [][2]string `json:"subset"`
}

type oracle struct {
	Normalize []struct {
		OK bool   `json:"ok"`
		V  string `json:"v"`
		E  string `json:"e"`
	} `json:"normalize"`
	Stability   []string `json:"stability"`
	Branch      []string `json:"branch"`
	Compare     []int    `json:"compare"`
	Constraints []struct {
		OK bool   `json:"ok"`
		S  string `json:"s"`
		P  string `json:"p"`
		Lo string `json:"lo"`
		Hi string `json:"hi"`
		E  string `json:"e"`
	} `json:"constraints"`
	Matches [][]*bool `json:"matches"`
	Compact []*string `json:"compact"`
	Subset  []bool    `json:"subset"`
}

func load(t *testing.T) (corpus, oracle) {
	var c corpus
	var o oracle
	for name, dst := range map[string]any{"testdata/corpus.json": &c, "testdata/oracle.json": &o} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, dst); err != nil {
			t.Fatal(err)
		}
	}
	return c, o
}

func TestOracleNormalize(t *testing.T) {
	c, o := load(t)
	for i, v := range c.Versions {
		got, err := Normalize(v)
		want := o.Normalize[i]
		if want.OK != (err == nil) || (want.OK && got != want.V) || (!want.OK && err.Error() != want.E) {
			t.Errorf("normalize(%q) = %q, %v; want %+v", v, got, err, want)
		}
		if s := ParseStability(v); s != o.Stability[i] {
			t.Errorf("parseStability(%q) = %q; want %q", v, s, o.Stability[i])
		}
		if b := NormalizeBranch(v); b != o.Branch[i] {
			t.Errorf("normalizeBranch(%q) = %q; want %q", v, b, o.Branch[i])
		}
	}
}

func TestOracleVersionCompare(t *testing.T) {
	c, o := load(t)
	for i, p := range c.Compare {
		if got := VersionCompare(p[0], p[1]); got != o.Compare[i] {
			t.Errorf("version_compare(%q, %q) = %d; want %d", p[0], p[1], got, o.Compare[i])
		}
	}
}

func TestOracleConstraints(t *testing.T) {
	c, o := load(t)
	for i, s := range c.Constraints {
		want := o.Constraints[i]
		got, err := ParseConstraints(s)
		if want.OK != (err == nil) {
			t.Errorf("parseConstraints(%q) err=%v; want %+v", s, err, want)
			continue
		}
		if !want.OK {
			if err.Error() != want.E {
				t.Errorf("parseConstraints(%q) error %q; want %q", s, err.Error(), want.E)
			}
			continue
		}
		if got.String() != want.S || got.PrettyString() != want.P {
			t.Errorf("parseConstraints(%q) = %q (%q); want %q (%q)", s, got.String(), got.PrettyString(), want.S, want.P)
		}
		if lo := got.LowerBound().String(); lo != want.Lo {
			t.Errorf("lowerBound(%q) = %q; want %q", s, lo, want.Lo)
		}
		if hi := got.UpperBound().String(); hi != want.Hi {
			t.Errorf("upperBound(%q) = %q; want %q", s, hi, want.Hi)
		}
		for j, v := range c.MatchVersions {
			if w := o.Matches[i][j]; w != nil && got.Matches(NewConstraint("==", v)) != *w {
				t.Errorf("%q matches %q = %v; want %v", s, v, !*w, *w)
			}
		}
		if o.Compact[i] != nil {
			if cc := CompactConstraint(got).String(); cc != *o.Compact[i] {
				t.Errorf("compact(%q) = %q; want %q", s, cc, *o.Compact[i])
			}
		}
	}
}

func TestOracleSubset(t *testing.T) {
	c, o := load(t)
	for i, p := range c.Subset {
		a := MustParseConstraints(p[0])
		b := MustParseConstraints(p[1])
		if got := IsSubsetOf(a, b); got != o.Subset[i] {
			t.Errorf("isSubsetOf(%q, %q) = %v; want %v", p[0], p[1], got, o.Subset[i])
		}
	}
}

// Probed with php:8.4.27 version_compare(); see oracle/scripts notes.
func TestVersionCompareTrailing(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.", "1.0", 0}, {"1.", "1", 0}, {"1.0.", "1.0.0", -1}, {"1.0.", "1.0.a", 1},
		{"1.0-", "1.0.0", -1}, {"1-", "1-a", 1}, {"1-", "1-pl", -1}, {"1+", "1", 0},
		{"1.0_", "1.0", 0}, {"", "1", -1}, {"1.0.#", "1.0", 0}, {"1..0", "1.0", 0},
		{"1.0..", "1.0", 0}, {"a.", "a", 0}, {"1.0-", "1.0-dev", 1},
	}
	for _, c := range cases {
		if got := VersionCompare(c.a, c.b); got != c.want {
			t.Errorf("version_compare(%q, %q) = %d; want %d", c.a, c.b, got, c.want)
		}
	}
}
