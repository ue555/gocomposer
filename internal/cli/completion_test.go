package cli

import (
	"reflect"
	"strings"
	"testing"

	"gocomposer/internal/phpjson"
)

func bindCompletion(t *testing.T, cmd *Command, cur int, tokens ...string) *compInput {
	t.Helper()
	ci := newCompInput(tokens, cur)
	def := appDefinition()
	if cmd != nil {
		def = mergedDefinition(cmd)
	}
	if err := ci.bind(def); err != nil {
		t.Fatalf("bind(%v): %v", tokens, err)
	}
	return ci
}

func TestCompletionInputBind(t *testing.T) {
	app := NewApp()
	show := app.byName["show"]
	cases := []struct {
		cmd    *Command
		cur    int
		tokens []string
		typ    string
		name   string
		value  string
	}{
		{nil, 1, []string{"composer"}, compArgumentValue, "command", ""},
		{nil, 1, []string{"composer", "sh"}, compArgumentValue, "command", "sh"},
		{nil, 1, []string{"composer", "--"}, compOptionName, "", "--"},
		{show, 2, []string{"composer", "show"}, compArgumentValue, "package", ""},
		{show, 2, []string{"composer", "show", "acme/"}, compArgumentValue, "package", "acme/"},
		{show, 3, []string{"composer", "show", "--format"}, compOptionValue, "format", ""},
		{show, 2, []string{"composer", "show", "--format=j"}, compOptionValue, "format", "j"},
		{show, 2, []string{"composer", "show", "-fj"}, compOptionValue, "format", "j"},
		{show, 3, []string{"composer", "show", "acme/log", "1.0"}, compArgumentValue, "version", "1.0"},
		{show, 4, []string{"composer", "show", "acme/log", "1.0"}, compNone, "", ""},
		{show, 2, []string{"composer", "show", "--al"}, compOptionName, "", "--al"},
	}
	for _, c := range cases {
		ci := bindCompletion(t, c.cmd, c.cur, c.tokens...)
		if ci.ctype != c.typ || ci.completionName() != c.name || ci.cvalue != c.value {
			t.Errorf("%v@%d: got (%s, %q, %q), want (%s, %q, %q)", c.tokens, c.cur, ci.ctype, ci.completionName(), ci.cvalue, c.typ, c.name, c.value)
		}
	}
}

func TestCompletionInputErrors(t *testing.T) {
	for _, c := range []struct {
		cur    int
		tokens []string
	}{
		{3, []string{"composer"}},
		{2, []string{"composer", ""}},
	} {
		if err := newCompInput(c.tokens, c.cur).bind(appDefinition()); err == nil {
			t.Errorf("%v@%d: expected an error", c.tokens, c.cur)
		}
	}
}

func TestCompletionFirstArgument(t *testing.T) {
	for _, c := range []struct {
		tokens []string
		want   string
	}{
		{[]string{"composer", "-d", "/tmp", "show"}, "show"},
		{[]string{"composer", "--working-dir", "/tmp", "show"}, "show"},
		{[]string{"composer", "-vvv", "install"}, "install"},
		{[]string{"composer", "--working-dir=/tmp", "show"}, "show"},
	} {
		ci := bindCompletion(t, nil, len(c.tokens), c.tokens...)
		got, ok, err := ci.firstArgument()
		if err != nil || !ok || got != c.want {
			t.Errorf("%v: got %q %v %v, want %q", c.tokens, got, ok, err, c.want)
		}
	}
}

func TestCompletionInputString(t *testing.T) {
	if s := newCompInput([]string{"composer", "global", "show"}, 3).String(); s != "composer global show |" {
		t.Errorf("String() = %q", s)
	}
	if s := newCompInput([]string{"composer", "global", "show"}, 2).String(); s != "composer global show|" {
		t.Errorf("String() = %q", s)
	}
	ci := compInputFromString("composer show 'a|' \"b c\" d\\ e", 2)
	want := []string{"composer", "show", "'a|'", `"b c"`, `d\ e`}
	if !reflect.DeepEqual(ci.tokens, want) {
		t.Errorf("fromString tokens = %q, want %q", ci.tokens, want)
	}
}

func TestCompletionFind(t *testing.T) {
	app := NewApp()
	for name, want := range map[string]string{"up": "update", "upgrade": "update", "why": "depends", "dump": "dump-autoload", "_complete": "_complete"} {
		c := app.symfonyFind(name)
		if c == nil || c.Name != want {
			t.Errorf("find(%q) = %v, want %s", name, c, want)
		}
	}
	for _, name := range []string{"us", "s", "nope", "re"} {
		if c := app.symfonyFind(name); c != nil {
			t.Errorf("find(%q) = %s, want none", name, c.Name)
		}
	}
	entries := app.allEntries()
	var keys []string
	for _, e := range entries[:6] {
		keys = append(keys, e.key)
	}
	if strings.Join(keys, " ") != "help list _complete completion about config" {
		t.Errorf("all() order = %v", keys)
	}
}

func TestFlattenSettingKeys(t *testing.T) {
	v, err := phpjson.Decode([]byte(`{"a": 1, "b": {"c": {"d": true}, "e": [1, 2]}, "f": [], "g": {}}`))
	if err != nil {
		t.Fatal(err)
	}
	got := flattenSettingKeys(v, "")
	want := []string{"a", "b", "b.c", "b.c.d", "b.e", "f", "g"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("flatten = %v, want %v", got, want)
	}
	if got := flattenSettingKeys(v, "repositories."); !reflect.DeepEqual(got, []string{"repositories.a", "repositories.b", "repositories.f", "repositories.g"}) {
		t.Errorf("repositories flatten = %v", got)
	}
}
