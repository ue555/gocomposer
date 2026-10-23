package semver

import (
	"encoding/json"
	"os"
	"testing"
)

// testdata/natoracle.json holds PHP 8.4 strnatcmp() results for the pairs
// of testdata/natcorpus.json (oracle image); never edit it by hand.
func TestNatCompareOracle(t *testing.T) {
	var pairs [][2]string
	var want []int
	raw, _ := os.ReadFile("testdata/natcorpus.json")
	if err := json.Unmarshal(raw, &pairs); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile("testdata/natoracle.json")
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	for i, p := range pairs {
		if got := NatCompare(p[0], p[1]); got != want[i] {
			t.Errorf("strnatcmp(%q, %q) = %d, want %d", p[0], p[1], got, want[i])
		}
	}
}
