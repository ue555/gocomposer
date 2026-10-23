// Replays traces recorded on official PHP by oracle/scripts/jsonmanip-dump.php
// (Composer f83588392df43d4fd027d84d26de459946ef4e80). Never edit the
// recorded files by hand. Regenerate with:
//
//	go run ./internal/jsonmanip/testdata/gencorpus > internal/jsonmanip/testdata/corpus.json
//	R="podman run --rm --network=none -v $PWD/.cache/oracle-source:/oracle-source:ro -v $PWD/oracle/scripts:/scripts:ro -v $PWD/internal/jsonmanip/testdata:/data:ro localhost/gocomposer-oracle:8.4 php /scripts/jsonmanip-dump.php"
//	$R tests > internal/jsonmanip/testdata/upstream-tests.json
//	$R corpus /data/corpus.json > internal/jsonmanip/testdata/corpus-oracle.json
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package jsonmanip

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"gocomposer/internal/phpjson"
)

type event struct {
	Op       string            `json:"op"`
	Args     []json.RawMessage `json:"args"`
	Ret      json.RawMessage   `json:"ret"`
	Exc      *excRecord        `json:"exc"`
	Assert   string            `json:"assert"`
	Expected json.RawMessage   `json:"expected"`
	Obs      *int              `json:"obs"`
}

type excRecord struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

// decodeTagged rebuilds a PHP value encoded by Rec::enc.
func decodeTagged(raw json.RawMessage) (phpjson.Value, error) {
	var plain any
	if err := json.Unmarshal(raw, &plain); err != nil {
		return nil, err
	}
	switch t := plain.(type) {
	case nil, bool, string:
		return t, nil
	case map[string]any:
	default:
		return nil, fmt.Errorf("unexpected %T", plain)
	}
	var tagged struct {
		T     string              `json:"t"`
		V     json.RawMessage     `json:"v"`
		B64   string              `json:"b64"`
		Class string              `json:"class"`
		Items [][]json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(raw, &tagged); err != nil {
		return nil, err
	}
	switch tagged.T {
	case "int":
		var n int64
		err := json.Unmarshal(tagged.V, &n)
		return n, err
	case "float":
		var s string
		if err := json.Unmarshal(tagged.V, &s); err != nil {
			return nil, err
		}
		return strconv.ParseFloat(s, 64)
	case "bytes":
		b, err := base64.StdEncoding.DecodeString(tagged.B64)
		return string(b), err
	case "array", "object":
		m := phpjson.NewMap()
		for _, it := range tagged.Items {
			var k string
			if err := json.Unmarshal(it[0], &k); err != nil {
				return nil, err
			}
			v, err := decodeTagged(it[1])
			if err != nil {
				return nil, err
			}
			m.Set(k, v)
		}
		if tagged.T == "object" {
			m.ForceObject = true
			return m, nil
		}
		if m.IsList() {
			l := make([]phpjson.Value, 0, m.Len())
			phpjson.Each(m, func(_ string, v phpjson.Value) { l = append(l, v) })
			return l, nil
		}
		return m, nil
	}
	return nil, fmt.Errorf("unknown tag %q", tagged.T)
}

func argStr(a []phpjson.Value, i int) string {
	s, _ := a[i].(string)
	return s
}

func argBool(a []phpjson.Value, i int, def bool) bool {
	if i >= len(a) {
		return def
	}
	b, _ := a[i].(bool)
	return b
}

func argInt(a []phpjson.Value, i int, def int) int {
	if i >= len(a) {
		return def
	}
	n, _ := a[i].(int64)
	return int(n)
}

// call invokes the Go port of a JsonManipulator method by PHP name.
func call(m *Manipulator, op string, a []phpjson.Value) (phpjson.Value, error) {
	var ret phpjson.Value
	ok, err := Try(func() bool {
		switch op {
		case "addLink":
			ret = m.AddLink(argStr(a, 0), argStr(a, 1), argStr(a, 2), argBool(a, 3, false))
		case "removeSubNode":
			ret = m.RemoveSubNode(argStr(a, 0), argStr(a, 1))
		case "addSubNode":
			ret = m.AddSubNode(argStr(a, 0), argStr(a, 1), a[2], argBool(a, 3, true))
		case "addMainKey":
			ret = m.AddMainKey(argStr(a, 0), a[1])
		case "removeMainKey":
			ret = m.RemoveMainKey(argStr(a, 0))
		case "removeMainKeyIfEmpty":
			ret = m.RemoveMainKeyIfEmpty(argStr(a, 0))
		case "changeEmptyMainKeyFromAssocToList":
			ret = m.ChangeEmptyMainKeyFromAssocToList(argStr(a, 0))
		case "addConfigSetting":
			ret = m.AddConfigSetting(argStr(a, 0), a[1])
		case "removeConfigSetting":
			ret = m.RemoveConfigSetting(argStr(a, 0))
		case "addProperty":
			ret = m.AddProperty(argStr(a, 0), a[1])
		case "removeProperty":
			ret = m.RemoveProperty(argStr(a, 0))
		case "addRepository":
			ret = m.AddRepository(argStr(a, 0), a[1], argBool(a, 2, true))
		case "insertRepository":
			ret = m.InsertRepository(argStr(a, 0), a[1], argStr(a, 2), argInt(a, 3, 0))
		case "removeRepository":
			ret = m.RemoveRepository(argStr(a, 0))
		case "setRepositoryUrl":
			ret = m.SetRepositoryUrl(argStr(a, 0), argStr(a, 1))
		case "addListItem":
			ret = m.AddListItem(argStr(a, 0), a[1], argBool(a, 2, true))
		case "insertListItem":
			ret = m.InsertListItem(argStr(a, 0), a[1], argInt(a, 2, 0))
		case "removeListItem":
			ret = m.RemoveListItem(argStr(a, 0), argInt(a, 1, 0))
		case "getContents":
			ret = m.Contents()
		case "format":
			ret = m.Format(a[0], argInt(a, 1, 0), argBool(a, 2, false))
		default:
			panic("unknown method " + op)
		}
		return true
	})
	if !ok {
		return nil, err
	}
	return ret, nil
}

// replay runs the events and returns a description of the first mismatch.
func replay(events []event) (string, int) {
	var m *Manipulator
	var observed []phpjson.Value
	checked := 0
	for i, ev := range events {
		if ev.Assert != "" {
			if ev.Obs == nil {
				continue // derived value (e.g. json_decode of the contents); covered by exact contents checks
			}
			want, err := decodeTagged(ev.Expected)
			if err != nil {
				return fmt.Sprintf("event %d: %v", i, err), checked
			}
			got := observed[*ev.Obs]
			if ev.Assert == "json" {
				w, _ := phpjson.Decode([]byte(want.(string)))
				g, _ := phpjson.Decode([]byte(got.(string)))
				if phpjson.MustEncode(w, 0) != phpjson.MustEncode(g, 0) {
					return fmt.Sprintf("event %d: JSON assertion failed\nexpected %s\ngot      %s", i, want, got), checked
				}
			} else if want != got {
				return fmt.Sprintf("event %d: assertion failed\nexpected %#v\ngot      %#v", i, want, got), checked
			}
			checked++
			continue
		}
		args := make([]phpjson.Value, len(ev.Args))
		for j, raw := range ev.Args {
			v, err := decodeTagged(raw)
			if err != nil {
				return fmt.Sprintf("event %d: %v", i, err), checked
			}
			args[j] = v
		}
		var ret phpjson.Value
		var err error
		if ev.Op == "new" {
			m, err = New(args[0].(string))
		} else if m == nil {
			return fmt.Sprintf("event %d: %s without manipulator", i, ev.Op), checked
		} else {
			ret, err = call(m, ev.Op, args)
		}
		observed = append(observed, ret)
		if ev.Exc != nil {
			e, isExc := err.(*Exception)
			if !isExc || e.Class != ev.Exc.Class || e.Message != ev.Exc.Message {
				return fmt.Sprintf("event %d %s: PHP threw %s(%q), Go %v", i, ev.Op, ev.Exc.Class, ev.Exc.Message, err), checked
			}
			continue
		}
		if err != nil {
			return fmt.Sprintf("event %d %s: Go error %v", i, ev.Op, err), checked
		}
		if ev.Op == "new" {
			continue
		}
		want, derr := decodeTagged(ev.Ret)
		if derr != nil {
			return fmt.Sprintf("event %d: %v", i, derr), checked
		}
		if want != ret {
			return fmt.Sprintf("event %d %s(%s): result differs\nPHP %#v\nGo  %#v", i, ev.Op, argsString(ev.Args), want, ret), checked
		}
	}
	return "", checked
}

func argsString(args []json.RawMessage) string {
	s := ""
	for i, a := range args {
		if i > 0 {
			s += ", "
		}
		if len(a) > 200 {
			s += string(a[:200]) + "..."
		} else {
			s += string(a)
		}
	}
	return s
}

// knownDifferences lists upstream test cases that cannot pass, with the
// reason. They are skipped, never rewritten.
var knownDifferences = map[string]string{}

// TestUpstreamJsonManipulatorTest replays tests/Composer/Test/Json/JsonManipulatorTest.php
// (all 62 test methods and their data sets) as recorded by
// `php jsonmanip-dump.php tests` in testdata/upstream-tests.json.
func TestUpstreamJsonManipulatorTest(t *testing.T) {
	data, err := os.ReadFile("testdata/upstream-tests.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Results []struct {
			Test   string  `json:"test"`
			Set    string  `json:"set"`
			Status string  `json:"status"`
			Events []event `json:"events"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	methods := map[string]bool{}
	for _, r := range doc.Results {
		methods[r.Test] = true
		name := r.Test
		if r.Set != "" {
			name += "/" + r.Set
		}
		r := r
		t.Run(name, func(t *testing.T) {
			if reason, ok := knownDifferences[name]; ok {
				t.Skip(reason)
			}
			if r.Status != "pass" {
				t.Fatalf("upstream test did not pass on PHP: %s", r.Status)
			}
			if msg, _ := replay(r.Events); msg != "" {
				t.Fatal(msg)
			}
		})
	}
	if len(methods) != 62 {
		t.Errorf("expected 62 upstream test methods, trace has %d", len(methods))
	}
}

// TestOracleCorpus replays testdata/corpus.json (generated by
// testdata/gencorpus) against testdata/corpus-oracle.json, the output of
// `php jsonmanip-dump.php corpus` on official PHP 8.4.
func TestOracleCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/corpus-oracle.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Results []struct {
			Events []event `json:"events"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	failures, ops := 0, 0
	for i, r := range doc.Results {
		ops += len(r.Events)
		if msg, _ := replay(r.Events); msg != "" {
			failures++
			if failures <= 15 {
				t.Errorf("case %d: %s", i, msg)
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d cases differ", failures, len(doc.Results))
	}
	t.Logf("%d cases, %d recorded calls", len(doc.Results), ops)
}

// TestDefinesMatchesComposer checks the DEFINES constant against
// JsonManipulator::DEFINES as read by reflection on official PHP.
func TestDefinesMatchesComposer(t *testing.T) {
	data, err := os.ReadFile("testdata/upstream-tests.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Defines string `json:"defines"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Defines != DEFINES {
		t.Fatalf("DEFINES differs from Composer:\n%q\n%q", doc.Defines, DEFINES)
	}
}
