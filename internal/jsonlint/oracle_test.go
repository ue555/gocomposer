package jsonlint

import (
	"encoding/json"
	"os"
	"testing"
)

// testdata/oracle.json is produced by oracle/scripts/jsonlint-dump.php with
// the seld/jsonlint version locked by Composer f83588392df4; never edit it.
func TestOracle(t *testing.T) {
	var cases []string
	raw, _ := os.ReadFile("testdata/corpus.json")
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	type res struct {
		Msg  string `json:"msg"`
		Key  string `json:"key"`
		Line int    `json:"line"`
	}
	var want []map[string]*res
	raw, _ = os.ReadFile("testdata/oracle.json")
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	if len(want) != len(cases) {
		t.Fatalf("oracle has %d rows, corpus %d", len(want), len(cases))
	}
	for i, c := range cases {
		for mode, flags := range map[string]int{"default": 0, "dup": DetectKeyConflicts, "dupassoc": DetectKeyConflicts | ParseToAssoc} {
			err := Lint(c, flags)
			w := want[i][mode]
			if w == nil {
				if err != nil {
					t.Errorf("case %d %q (%s): unexpected error %q", i, c, mode, err)
				}
				continue
			}
			if err == nil {
				t.Errorf("case %d %q (%s): expected error %q", i, c, mode, w.Msg)
				continue
			}
			if err.Error() != w.Msg {
				t.Errorf("case %d %q (%s):\n got: %q\nwant: %q", i, c, mode, err.Error(), w.Msg)
			}
			if de, ok := err.(*DuplicateKeyError); ok != (w.Key != "") || (ok && (de.Key != w.Key || de.Line != w.Line)) {
				t.Errorf("case %d (%s): duplicate details mismatch: %+v vs %+v", i, mode, err, w)
			}
		}
	}
}
