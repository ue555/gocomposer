package httpdl

import (
	"reflect"
	"testing"

	"gocomposer/internal/phpjson"
)

// Cases follow Composer f83588392df43d4fd027d84d26de459946ef4e80
// tests/Composer/Test/Util/Http/ResponseWarningsTest.php semantics.
func TestOutputResponseWarnings(t *testing.T) {
	decode := func(s string) phpjson.Value {
		v, err := phpjson.Decode([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	cases := []struct {
		name      string
		url       string
		body      string
		decorated bool
		want      []string
		wrote     bool
	}{
		{"legacy", "https://repo.example.org/packages.json", `{"warning": "w", "info": "i"}`, false,
			[]string{"<warning>Warning from https://repo.example.org/packages.json: w</warning>", "<info>Info from https://repo.example.org/packages.json: i</info>"}, true},
		{"versions", "https://a.example.org", `{"warning": "w", "warning-versions": "<2.0", "info": "i", "info-versions": ">=2.11"}`, false,
			[]string{"<info>Info from https://a.example.org: i</info>"}, true},
		{"modern", "https://b.example.org", `{"warnings": [{"versions": "*", "message": "m1"}, {"versions": "<1", "message": "m2"}], "infos": [{"versions": "^2.11", "message": "i1"}]}`, false,
			[]string{"<warning>Warning from https://b.example.org: m1</warning>", "<info>Info from https://b.example.org: i1</info>"}, true},
		{"ansi stripped", "https://c.example.org", "{\"warning\": \"a \\u001b[31mred\\u001b[0m b\"}", false,
			[]string{"<warning>Warning from https://c.example.org: a red b</warning>"}, true},
		{"ansi kept", "https://d.example.org", "{\"warning\": \"a \\u001b[31mred\\u001b[0m\"}", true,
			[]string{"<warning>Warning from https://d.example.org: a \x1b[31mred\x1b[0m</warning>"}, true},
		{"credentials sanitized", "https://user:secret@e.example.org/x", `{"info": "i"}`, false,
			[]string{"<info>Info from https://user:***@e.example.org/x: i</info>"}, true},
		{"nothing", "https://f.example.org", `{"packages": []}`, false, nil, false},
		{"not array", "https://g.example.org", `"x"`, false, nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ResetResponseWarnings()
			var got []string
			wrote := OutputResponseWarnings(func(s string) { got = append(got, s) }, c.decorated, c.url, decode(c.body))
			if wrote != c.wrote || !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %v %q, want %v %q", wrote, got, c.wrote, c.want)
			}
		})
	}
}

func TestOutputResponseWarningsShownOncePerHost(t *testing.T) {
	ResetResponseWarnings()
	var got []string
	w := func(s string) { got = append(got, s) }
	body, _ := phpjson.Decode([]byte(`{"warning": "same", "packages": {"a": 1}}`))
	other, _ := phpjson.Decode([]byte(`{"warning": "same", "packages": {"b": 2}}`))
	if !OutputResponseWarnings(w, false, "http://127.0.0.1:8899/repo/packages.json", body) {
		t.Fatal("expected output")
	}
	// same host and port, same warning keys: reported as shown, not repeated
	if !OutputResponseWarnings(w, false, "127.0.0.1:8899", other) {
		t.Fatal("expected shown before")
	}
	// another port is another host
	OutputResponseWarnings(w, false, "http://127.0.0.1:8900/repo", body)
	if len(got) != 2 {
		t.Fatalf("got %q", got)
	}
}
