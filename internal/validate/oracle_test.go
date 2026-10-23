package validate

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"gocomposer/internal/phpjson"
)

// testdata/oracle.json is produced by oracle/scripts/validate-dump.php
// (PHP 8.4, Composer f83588392df4); never edit it.
func TestOracle(t *testing.T) {
	raw, _ := os.ReadFile("testdata/corpus.json")
	corpus, err := phpjson.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	cm := phpjson.AsMap(corpus)
	var want struct {
		Emails   []bool    `json:"emails"`
		URLs     [][2]bool `json:"urls"`
		Packages []struct {
			Errors    []string `json:"errors"`
			Warnings  []string `json:"warnings"`
			Exception string   `json:"exception"`
		} `json:"packages"`
	}
	raw, _ = os.ReadFile("testdata/oracle.json")
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	i := 0
	phpjson.Each(mustGet(cm, "emails"), func(_ string, v phpjson.Value) {
		if got := ValidEmail(v.(string)); got != want.Emails[i] {
			t.Errorf("email %q: got %v want %v", v, got, want.Emails[i])
		}
		i++
	})
	i = 0
	phpjson.Each(mustGet(cm, "urls"), func(_ string, v phpjson.Value) {
		got := [2]bool{FilterURL(v.(string), []string{"http", "https"}), FilterURL(v.(string), []string{"irc", "ircs"})}
		if got != want.URLs[i] {
			t.Errorf("url %q: got %v want %v", v, got, want.URLs[i])
		}
		i++
	})
	i = 0
	phpjson.Each(mustGet(cm, "packages"), func(_ string, v phpjson.Value) {
		m := phpjson.AsMap(v).Clone()
		if !m.Isset("version") {
			m.Set("version", "1.0.0")
		}
		if !m.Isset("name") {
			m.Set("name", "dummy/dummy")
		}
		l := &ArrayLoader{Flags: CheckAll}
		_, err := l.Load(m)
		var gotErr []string
		if ie, ok := err.(*InvalidPackageError); ok {
			gotErr = ie.Errors
		} else if err != nil {
			t.Errorf("package %d: unexpected error %v", i, err)
		}
		w := want.Packages[i]
		if len(gotErr) == 0 && len(w.Errors) == 0 {
			gotErr, w.Errors = nil, nil
		}
		gotWarn := l.Warnings()
		if len(gotWarn) == 0 && len(w.Warnings) == 0 {
			gotWarn, w.Warnings = nil, nil
		}
		if !reflect.DeepEqual(gotErr, w.Errors) {
			t.Errorf("package %d errors:\n got %q\nwant %q", i, gotErr, w.Errors)
		}
		if !reflect.DeepEqual(gotWarn, w.Warnings) {
			t.Errorf("package %d warnings:\n got %q\nwant %q", i, gotWarn, w.Warnings)
		}
		i++
	})
}
