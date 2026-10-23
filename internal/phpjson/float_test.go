package phpjson

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"testing"
)

// testdata/floats-oracle.json holds PHP 8.4 json_encode/var_export/serialize
// output for the IEEE bit patterns in testdata/floats.json; never edit it.
func TestFloatFormattingOracle(t *testing.T) {
	var bits []string
	raw, _ := os.ReadFile("testdata/floats.json")
	if err := json.Unmarshal(raw, &bits); err != nil {
		t.Fatal(err)
	}
	var want []struct {
		JSON      string `json:"json"`
		VarExport string `json:"var_export"`
		Serialize string `json:"serialize"`
	}
	raw, _ = os.ReadFile("testdata/floats-oracle.json")
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	for i, b := range bits {
		u, err := strconv.ParseUint(b, 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		f := math.Float64frombits(u)
		if got := FormatFloat(f); got != want[i].JSON {
			t.Errorf("json %v: %q want %q", f, got, want[i].JSON)
		}
		if got := FormatFloatVarExport(f); got != want[i].VarExport {
			t.Errorf("var_export %v: %q want %q", f, got, want[i].VarExport)
		}
		if got := "d:" + FormatFloatSerialize(f) + ";"; got != want[i].Serialize {
			t.Errorf("serialize %v: %q want %q", f, got, want[i].Serialize)
		}
	}
}
