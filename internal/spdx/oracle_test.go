package spdx

import (
	"encoding/json"
	"os"
	"testing"
)

// testdata/oracle.json comes from oracle/scripts/spdx-dump.php (PHP 8.4,
// composer/spdx-licenses locked by Composer f83588392df4); never edit it.
func TestOracle(t *testing.T) {
	var corpus []string
	raw, _ := os.ReadFile("testdata/corpus.json")
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
	}
	var want []struct {
		Valid bool          `json:"valid"`
		ByID  []interface{} `json:"byId"`
	}
	raw, _ = os.ReadFile("testdata/oracle.json")
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	for i, c := range corpus {
		if got := IsValidLicenseString(c); got != want[i].Valid {
			t.Errorf("%q: valid=%v want %v", c, got, want[i].Valid)
		}
		info := LicenseByIdentifier(c)
		if (info == nil) != (want[i].ByID == nil) {
			t.Errorf("%q: byId presence mismatch", c)
			continue
		}
		if info != nil && (info.Name != want[i].ByID[0] || info.OsiApproved != want[i].ByID[1] || info.URL != want[i].ByID[2] || info.Deprecated != want[i].ByID[3]) {
			t.Errorf("%q: byId %+v want %v", c, info, want[i].ByID)
		}
	}
}
