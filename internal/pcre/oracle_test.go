// Oracle comparison for the regex engine. testdata/corpus.json is generated
// by testdata/gencorpus; testdata/oracle.json is the output of
// oracle/scripts/pcre-dump.php run on official PHP 8.4 (PCRE2 10.44) with
// Composer f83588392df43d4fd027d84d26de459946ef4e80's vendor/composer/pcre.
// Never edit oracle.json by hand. Regenerate with:
//
//	go run ./internal/pcre/testdata/gencorpus > internal/pcre/testdata/corpus.json
//	go run ./internal/pcre/testdata/gencorpus -limits > internal/pcre/testdata/limits-corpus.json
//	podman run --rm --network=none -v $PWD/.cache/oracle-source:/oracle-source:ro \
//	  -v $PWD/oracle/scripts:/scripts:ro -v $PWD/internal/pcre/testdata:/data:ro \
//	  localhost/gocomposer-oracle:8.4 php /scripts/pcre-dump.php /data/corpus.json > internal/pcre/testdata/oracle.json
//	(same with /data/limits-corpus.json > internal/pcre/testdata/limits.json)
//
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package pcre

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type oracleCase struct {
	Op    string `json:"op"`
	P     string `json:"p"`
	S     string `json:"s"`
	R     string `json:"r"`
	D     string `json:"d"`
	C     string `json:"c"`
	Limit *int   `json:"limit"`
	// Unsupported cases are valid PCRE the engine rejects explicitly.
	Unsupported bool `json:"unsupported"`
}

type oracleResult struct {
	R   json.RawMessage      `json:"r"`
	M   [][2]json.RawMessage `json:"m"`
	N   int                  `json:"n"`
	Err string               `json:"err"`
}

// oracleStr decodes a dump string: plain JSON string, null, or {"b64": ...}
// for bytes that are not valid UTF-8.
func oracleStr(raw json.RawMessage) (*string, error) {
	if string(raw) == "null" {
		return nil, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return &s, nil
	}
	var b struct {
		B64 string `json:"b64"`
	}
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	d, err := base64.StdEncoding.DecodeString(b.B64)
	if err != nil {
		return nil, err
	}
	s = string(d)
	return &s, nil
}

// serial mirrors serial() of pcre-dump.php.
func serial(m *Match) string {
	var parts []string
	for _, e := range m.PHPArray(false) {
		if e.Value == nil {
			parts = append(parts, e.Key+"=N")
		} else {
			parts = append(parts, e.Key+"=S"+strconv.Itoa(len(*e.Value))+":"+*e.Value)
		}
	}
	return "<" + strings.Join(parts, "|") + ">"
}

func loadOracle(t *testing.T) ([]oracleCase, []oracleResult) {
	t.Helper()
	var c struct {
		Cases []oracleCase `json:"cases"`
	}
	var o struct {
		PHP     string         `json:"php"`
		PCRE    string         `json:"pcre"`
		Defines string         `json:"defines"`
		Results []oracleResult `json:"results"`
	}
	for name, dst := range map[string]any{"testdata/corpus.json": &c, "testdata/oracle.json": &o} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, dst); err != nil {
			t.Fatal(err)
		}
	}
	if len(c.Cases) != len(o.Results) {
		t.Fatalf("corpus has %d cases, oracle %d results", len(c.Cases), len(o.Results))
	}
	if o.Defines == "" {
		t.Fatal("oracle lacks JsonManipulator::DEFINES")
	}
	for i := range c.Cases {
		c.Cases[i].P = strings.ReplaceAll(c.Cases[i].P, "@@DEFINES@@", o.Defines)
	}
	return c.Cases, o.Results
}

func pairsJSON(m *Match, unmatchedAsNull bool) string {
	out := []any{}
	for _, e := range m.PHPArray(unmatchedAsNull) {
		if e.Value == nil {
			out = append(out, []any{e.Key, nil})
		} else {
			out = append(out, []any{e.Key, *e.Value})
		}
	}
	return phpJSON(out)
}

// phpJSON encodes like json_encode(..., JSON_UNESCAPED_SLASHES|JSON_UNESCAPED_UNICODE).
func phpJSON(v any) string {
	switch t := v.(type) {
	case nil:
		return "null"
	case string:
		b := []byte{'"'}
		for i := 0; i < len(t); i++ {
			c := t[i]
			switch c {
			case '"':
				b = append(b, '\\', '"')
			case '\\':
				b = append(b, '\\', '\\')
			case '\n':
				b = append(b, '\\', 'n')
			case '\r':
				b = append(b, '\\', 'r')
			case '\t':
				b = append(b, '\\', 't')
			case '\b':
				b = append(b, '\\', 'b')
			case '\f':
				b = append(b, '\\', 'f')
			default:
				if c < 0x20 {
					b = append(b, fmt.Sprintf(`\u%04x`, c)...)
				} else {
					b = append(b, c)
				}
			}
		}
		return string(append(b, '"'))
	case []any:
		s := "["
		for i, e := range t {
			if i > 0 {
				s += ","
			}
			s += phpJSON(e)
		}
		return s + "]"
	}
	panic("unsupported")
}

func errKind(err error) string {
	switch {
	case errors.Is(err, ErrBacktrackLimit), errors.Is(err, ErrRecursionLimit):
		return "limit"
	default:
		return "compile"
	}
}

func TestOracle(t *testing.T) {
	cases, results := loadOracle(t)
	failures, unsupported := 0, 0
	for i, c := range cases {
		want := results[i]
		got, err := runCase(c)
		if c.Unsupported {
			if want.Err != "" || !errors.Is(err, ErrUnsupported) {
				t.Errorf("case %d %q: expected PHP success and Go ErrUnsupported, got PHP %q / Go %v", i, c.P, want.Err, err)
				failures++
			}
			unsupported++
			continue
		}
		if want.Err != "" {
			if err == nil {
				t.Errorf("case %d %s %q on %q: PHP error %s, Go succeeded", i, c.Op, c.P, c.S, want.Err)
				failures++
			} else if want.Err == "compile" && errKind(err) != "compile" {
				t.Errorf("case %d: PHP %s, Go %v", i, want.Err, err)
				failures++
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d %s %q on %q: Go error %v", i, c.Op, c.P, c.S, err)
			failures++
			continue
		}
		if msg := compare(c, want, got); msg != "" {
			t.Errorf("case %d %s %q on %q: %s", i, c.Op, c.P, c.S, msg)
			failures++
		}
		if failures > 30 {
			t.Fatal("too many failures")
		}
	}
	t.Logf("%d oracle cases (%d explicitly unsupported)", len(cases), unsupported)
}

type goResult struct {
	matched bool
	pairs   string
	str     string
	n       int
}

func runCase(c oracleCase) (goResult, error) {
	limit := -1
	if c.Limit != nil {
		limit = *c.Limit
	}
	switch c.Op {
	case "match":
		m, err := IsMatch(c.P, c.S)
		if err != nil {
			return goResult{}, err
		}
		if m == nil {
			return goResult{pairs: "[]"}, nil
		}
		return goResult{matched: true, pairs: pairsJSON(m, true)}, nil
	case "replace":
		s, n, err := Replace(c.P, c.R, c.S, limit)
		return goResult{str: s, n: n}, err
	case "callback":
		s, n, err := ReplaceCallback(c.P, serial, c.S, limit)
		return goResult{str: s, n: n}, err
	case "quote":
		return goResult{str: Quote(c.S, c.D)}, nil
	case "addcslashes":
		return goResult{str: AddCSlashes(c.S, c.C)}, nil
	}
	return goResult{}, fmt.Errorf("unknown op %s", c.Op)
}

func compare(c oracleCase, want oracleResult, got goResult) string {
	switch c.Op {
	case "match":
		var r int
		if err := json.Unmarshal(want.R, &r); err != nil {
			return "bad oracle r"
		}
		if (r == 1) != got.matched {
			return fmt.Sprintf("matched: PHP %d, Go %v", r, got.matched)
		}
		wp := []any{}
		for _, p := range want.M {
			k, err1 := oracleStr(p[0])
			v, err2 := oracleStr(p[1])
			if err1 != nil || err2 != nil || k == nil {
				return "bad oracle m"
			}
			if v == nil {
				wp = append(wp, []any{*k, nil})
			} else {
				wp = append(wp, []any{*k, *v})
			}
		}
		if w := phpJSON(wp); w != got.pairs {
			return fmt.Sprintf("groups:\nPHP %s\nGo  %s", w, got.pairs)
		}
	default:
		sp, err := oracleStr(want.R)
		if err != nil || sp == nil {
			return "bad oracle r"
		}
		s := *sp
		if s != got.str {
			return fmt.Sprintf("result:\nPHP %q\nGo  %q", s, got.str)
		}
		if (c.Op == "replace" || c.Op == "callback") && want.N != got.n {
			return fmt.Sprintf("count: PHP %d, Go %d", want.N, got.n)
		}
	}
	return ""
}
