// Command gencorpus writes internal/pcre/testdata/corpus.json, the regex
// corpus evaluated by oracle/scripts/pcre-dump.php on official PHP and by
// internal/pcre's oracle test. The JsonManipulator patterns are built exactly
// like src/Composer/Json/JsonManipulator.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80 builds them.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Usage: go run ./internal/pcre/testdata/gencorpus > internal/pcre/testdata/corpus.json
//
//	go run ./internal/pcre/testdata/gencorpus -limits > internal/pcre/testdata/limits-corpus.json
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"gocomposer/internal/pcre"
	"gocomposer/internal/phpjson"
)

type Case struct {
	Op    string `json:"op"`
	P     string `json:"p,omitempty"`
	S     string `json:"s"`
	R     string `json:"r,omitempty"`
	D     string `json:"d,omitempty"`
	C     string `json:"c,omitempty"`
	Limit *int   `json:"limit,omitempty"`
	// Unsupported marks valid PCRE syntax the engine rejects on purpose.
	Unsupported bool `json:"unsupported,omitempty"`
}

// DEFINES is a placeholder that pcre-dump.php and the Go tests replace with
// JsonManipulator::DEFINES as read from Composer by reflection.
const DEFINES = "@@DEFINES@@"

func enc(s string) string { return phpjson.MustEncode(s, phpjson.ComposerFlags) }

var cases []Case

func add(c Case) { cases = append(cases, c) }

func match(p string, subjects ...string) {
	for _, s := range subjects {
		add(Case{Op: "match", P: p, S: s})
	}
}

func both(p string, subjects ...string) {
	for _, s := range subjects {
		add(Case{Op: "match", P: p, S: s})
		add(Case{Op: "callback", P: p, S: s})
	}
}

func replace(p, r string, subjects ...string) {
	for _, s := range subjects {
		add(Case{Op: "replace", P: p, R: r, S: s})
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-limits" {
		limits()
		return
	}
	features()
	manipulatorPatterns()
	quoting()
	out, err := json.Marshal(map[string]any{"cases": cases})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))
}

func features() {
	both(`/a+/`, "", "a", "baaab", "xyz", "aXa")
	both(`/a*/`, "", "baaab", "aaa")
	both(`/a*?/`, "baaab")
	both(`/a+?b/`, "aaab", "xab")
	both(`/(a)|(b)/`, "b", "a", "c", "ab")
	both(`/(a)?(b)?(c)/`, "c", "ac", "bc", "abc")
	both(`/(?<x>a)(?P<y>b)?/`, "a", "ab", "zzab")
	both(`/(?'q'z+)/`, "azzz")
	both(`/x(?:ab)+y/`, "xababy", "xy", "xaby")
	both(`/x(ab)*y/`, "xababy", "xy")
	both(`/(a|ab)(c|bcd)(d*)/`, "abcd")
	both(`/^(a+)+$/`, "aaaa", "aaab")
	both(`/(?>a+)b/`, "aaab", "aaa")
	both(`/a++b/`, "aaab")
	both(`/a++a/`, "aaaa")
	both(`/a?+a/`, "a", "aa")
	both(`/(?:ab)?+ab/`, "ab", "abab")
	both(`/(?:a|b)*+b/`, "aab")
	both(`/(?:a|b)*b/`, "aab")
	both(`/a{2}/`, "a", "aaa")
	both(`/a{2,}/`, "a", "aaaaa")
	both(`/a{2,3}/`, "aaaaa")
	both(`/a{2,3}?/`, "aaaaa")
	both(`/a{,3}/`, "aaaaa", "b")
	both(`/a{x}/`, "a{x}")
	both(`/a{/`, "a{")
	both(`/{a/`, "{a")
	both(`/a{1/`, "a{1")
	both(`/a{1,2/`, "a{1,2")
	both(`/x{0}y/`, "xy", "y")
	both(`/(ab){2}/`, "ababab")
	both(`/(?:a(b)?){2}/`, "aba", "aab", "abab")
	both(`/(a*)*/`, "b", "aab")
	both(`/(a*)+/`, "b", "aab")
	both(`/(a|)*b/`, "aab", "b")
	both(`/(?:(a)|b)*/`, "ab", "ba")
	both(`/^/`, "", "abc")
	both(`/$/`, "", "abc", "abc\n", "a\nb\n")
	both(`/^/m`, "a\nb\n", "\n\n")
	both(`/$/m`, "a\nb\n", "\n\n")
	both(`/a$/`, "a\n", "a\n\n", "a")
	both(`/a$/D`, "a\n", "a")
	both(`/a\z/`, "a\n", "a")
	both(`/a\Z/`, "a\n", "a", "a\n\n")
	both(`/\Aa/`, "aa", "ba")
	both(`/^a/m`, "b\na\na", "a")
	both(`/a$/m`, "a\nba\n", "ab")
	both(`/./`, "\n", "\na", "ab")
	both(`/./s`, "\n", "ab")
	both(`/.*/`, "ab\ncd")
	both(`/.*/s`, "ab\ncd")
	both(`/\bfoo\b/`, "foo", "a foo b", "afoo", "foo_")
	both(`/\Bo\B/`, "foo", "o")
	both(`/\d+/`, "abc123def45")
	both(`/\D+/`, "abc123def45")
	both(`/\s+/`, "a \t\n\v\f\rb", "a\u00a0b")
	both(`/\S+/`, " ab c ")
	both(`/\w+/`, "héllo wörld_1")
	both(`/\W+/`, "a, b!")
	both(`/\h+/`, "a \t\u00a0b")
	both(`/\v+/`, "a\n\v\f\rb")
	both(`/\N+/`, "ab\ncd")
	both(`/[abc]+/`, "xxabcabcyy")
	both(`/[^abc]+/`, "abxyzc")
	both(`/[a-c]+/`, "abcd")
	both(`/[a-]+/`, "a-a-b")
	both(`/[-a]+/`, "a-a-b")
	both(`/[]a]+/`, "a]a]b")
	both(`/[^]a]+/`, "a]bcd")
	both(`/[\]]+/`, "]]a")
	both(`/[\\]+/`, `a\\b`)
	both(`/[\d.]+/`, "v1.2.3x")
	both(`/[\s\S]+/`, "a\nb")
	both(`/[\w-]+/`, "foo-bar_baz qux")
	both(`/[^\s]+/`, "  ab  ")
	both(`/[\bx]+/`, "\x08x\x08")
	both(`/[\/]+/`, "a//b")
	both(`/["\\bfnrt\/]/`, `\`, "b", "/", "x")
	both(`/[0-9A-Fa-f]{4}/`, "u00F1", "zz12gh")
	both(`/[+-]?+\d/`, "+1", "-2", "3")
	both(`/[\x41-\x43]+/`, "ABCD")
	both(`/\x41\x{42}\101/`, "ABA")
	both(`/\0/`, "a\x00b")
	both(`/\000/`, "a\x00b")
	both(`/a\tb\nc\rd\fe/`, "a\tb\nc\rd\fe")
	both(`/\e\a/`, "\x1b\x07")
	both(`/\.\+\*\?\[\^\]\$\(\)\{\}\=\!\<\>\|\:\-\#\//`, `.+*?[^]$(){}=!<>|:-#/`)
	both(`/abc/i`, "ABC", "aBc", "abd")
	both(`/[a-c]+/i`, "ABCabcD")
	both(`/[^a]+/i`, "xAay")
	both(`/é/i`, "É", "é")
	both(`/a(?i)b/`, "aB", "AB")
	both(`/a(?i:b)c/`, "aBc", "aBC")
	both(`/(?i)a(?-i)b/`, "Ab", "AB")
	both(`/(a(?i)b)c/`, "aBc", "aBC")
	both(`/(?s).+/`, "a\nb")
	both(`/(?m)^b/`, "a\nb")
	both(`/(?x) a b /`, "ab")
	both(`/a b # comment
 c/x`, "abc", "a b c")
	both(`/[ ]a/x`, " a", "a")
	both(`/a\ b/x`, "a b")
	both(`/a\#b/x`, "a#b")
	both(`/a +/x`, "aaa")
	both(`/a+ +/x`, "aaa")
	both(`/a *+ b/x`, "aaab")
	both(`/(?: a | b ) {2}/x`, "ab")
	both(`/a+ ?a/x`, "aaa")
	both(`/a+ +a/x`, "aaa")
	both(`/a{2} +/x`, "aaa")
	both(`/a * ?b/x`, "aab")
	both(`/a (?#c) +/x`, "aaa")
	both(`/a # c
 +/x`, "aaa")
	both(`/\101\61\1234/`, "A1S4", "A1\x534")
	both(`/[\101-\103\8]+/`, "ABC8D")
	both(`/(a)\10/`, "a\x08")
	both(`/./`, "é")
	both(`/\xc3./`, "é")
	both(`/[\x80-\xff]+/`, "aéb")
	both(`/é+/`, "ééé")
	for _, p := range []string{`/(a)\1/`, `/a(?<=a)/`, `/a\Gb/`, `/a/u`, `/\Qa.\E/`, `/(a)?(?(1)b|c)/`, `/\p{L}/`, `/a\Kb/`, `/(*FAIL)/`, `/[[:alpha:]]/`} {
		add(Case{Op: "match", P: p, S: "aa.b", Unsupported: true})
	}
	both(`/a(?#comment)b/`, "ab")
	both(`/a(?=b)/`, "ab", "ac")
	both(`/a(?!b)/`, "ab", "ac")
	both(`/(?=(a))a/`, "a")
	both(`/(?!(a))b/`, "b")
	both(`/-? (?= [1-9]|0(?!\d) ) \d++/x`, "0", "01", "-12", "-0", "007")
	both(`/(?(DEFINE)(?<d>\d+))x(?&d)y/`, "x12y", "xy")
	both(`/(?(DEFINE)(?<d>\d+))x((?&d))y/`, "x12y")
	both(`/(?(DEFINE)(?<d>\d+))x(?P>d)/`, "x12")
	both(`/(?<p>\((?:[^()]|(?&p))*\))/`, "a(b(c)d)e", "((x)")
	both(`/\((?:[^()]|(?R))*\)/`, "(a(b)c)")
	both(`/^(a(?1)?b)$/`, "aabb", "aab")
	both(`/(?(DEFINE)(?<a>x+))^(?&a)x$/`, "xxx")
	both(`/(?(DEFINE)(?<w>a|ab))^(?&w)c/`, "abc")
	both(`/(?(DEFINE)(?<w>a+))^(?&w)a/`, "aaa")
	both(`/(?(DEFINE)(?<w>(?<inner>a)+))^(?&w)(?&w)/`, "aaa")
	both(`/(?:(?<n>a)|b)+/`, "ab", "ba")
	both(`#^\{(.*)\}$#s`, "{}", "{\n}", "{a}\n", "{a}b", " {a}")
	both(`/foo|foobar/`, "foobar")
	both(`/(foo|foobar)$/`, "foobar")
	both(`/x*/`, "axxb")
	both(`/\b/`, "ab cd")
	both(`/(?=a)/`, "bab")
	both(`/a|(?=b)/`, "xab")
	both(`/^\s*/`, "  a")
	both(`/\s*$/`, "a  \n")
	match(`/(/`, "a")
	match(`/a)/`, "a")
	match(`/[a/`, "a")
	match(`/a**/`, "a")
	match(`/a\/`, "a")
	match(`abc`, "abc")
	match(`/abc`, "abc")
	match(`/a/k`, "a")
	match(`/(?<n>a)(?<n>b)/`, "ab")
	match(`/(?&nope)/`, "a")
	match(`/[z-a]/`, "a")
	match(`/a{3,2}/`, "a")
	match(`/\i/`, "i")
	match(`{a{1}}`, "a")
	match(`(a(b)c)`, "abc")
	match(`[a]b]`, "ab")
	match(`<a>`, "a")
	match(` /a/`, "a")
	match("/a/ \n", "a")
	match(`#a\#b#`, "a#b")
	match(`{a\}b}`, "a}b")
	match(`/a.c/U`, "abcabc")
	match(`/a.*c/U`, "abcabc")
	match(`/a.*?c/U`, "abcabc")

	// replacement strings
	for _, r := range []string{
		"[$0]", "[\\0]", "[$1]", "[${1}]", "[${1]", "[$12]", "[${12}]", "[$9]", "[\\\\1]", "[\\$1]", "[$$1]",
		"[\\]", "[$]", "[\\x]", "[$a]", "[${a}]", "\\", "$", "[\\\\\\1]", "[\\\\$1]", "${1}0", "$10", "\\10", "\\1\\2\\3",
	} {
		replace(`/(a)(b)?/`, r, "xaby", "xay", "aa")
	}
	replace(`/x*/`, "-", "abc", "")
	replace(`/a|/`, "-", "bab")
	replace(`/\s*$/`, "!", "ab  \n")
	replace(`/^/m`, "> ", "a\nb\n")
	replace(`/$/m`, "<", "a\nb\n")
	for _, lim := range []int{0, 1, 2, -1} {
		l := lim
		add(Case{Op: "replace", P: `/a/`, R: "b", S: "aaaa", Limit: &l})
		add(Case{Op: "callback", P: `/a/`, S: "aaaa", Limit: &l})
	}
	// addcslashes'd replacements as JsonManipulator builds them
	for _, s := range []string{`foo\bar`, `$1 and \1`, `price $5.00`, `a\\$b`, `${1}`} {
		esc := pcre.AddCSlashes(s, `\$`)
		add(Case{Op: "addcslashes", S: s, C: `\$`})
		replace(`#\}$#`, esc+"}", "{x}")
	}
	add(Case{Op: "addcslashes", S: "a\nb\x01\xff$\\", C: "\n\x01\xff$\\"})
}

var docs = []string{
	`{}`,
	"{\n}",
	`{"require":{"foo/bar":"^1.0"}}`,
	"{\n    \"name\": \"vendor/pkg\",\n    \"require\": {\n        \"php\": \">=7.2\",\n        \"foo/bar\": \"^1.0\"\n    },\n    \"require-dev\": {\n        \"phpunit/phpunit\": \"^9\"\n    }\n}",
	"{\r\n    \"name\": \"vendor/pkg\",\r\n    \"require\": {\r\n        \"foo\\/bar\": \"^1.0\"\r\n    },\r\n    \"config\": {}\r\n}",
	"{\n\t\"config\": {\n\t\t\"platform\": {\"php\": \"8.1\"},\n\t\t\"sort-packages\": true\n\t},\n\t\"extra\": {\"a.b\": {\"c\": [1, 2, {\"d\": null}]}}\n}",
	"{\n  \"description\": \"Some U\\u00F1icode \\\"quoted\\\" \\\\ back\",\n  \"require\": {\n    \"foo/bar\": \"^1.0\"\n  }\n}",
	"{\n    \"name\": \"ünïcödé/pkg\",\n    \"keywords\": [\"日本語\", \"emoji 🎉\"],\n    \"repositories\": [\n        {\"type\": \"vcs\", \"url\": \"https://example.org/a\"},\n        {\"packagist.org\": false},\n        {\"name\": \"foo\", \"type\": \"composer\", \"url\": \"https://repo.example\"}\n    ]\n}",
	"{\n    \"repositories\": {\n        \"foo\": {\"type\": \"vcs\", \"url\": \"https://example.org/a\"},\n        \"packagist.org\": false,\n        \"0\": {\"type\": \"path\", \"url\": \"../x\"}\n    }\n}",
	"{\"numbers\": [0, -0, 1.5, -1.5e10, 2E-3, 1e+2, 123456789012345678901234567890], \"bools\": [true, false, null], \"nested\": [[[]], {}, [{}]]}",
	"{\n    \"a\": 1,\n    \"require\": {},\n    \"b\": 2\n}",
	"{\n    \"require\": {\n        \"foo/bar\": \"^1.0\",\n    }\n}",
	"{\n    // comment\n    \"require\": {\"foo/bar\": \"^1.0\"}\n}",
	"{\n    'require': {\"foo/bar\": \"^1.0\"}\n}",
	"{\n    \"require\": {\"foo/bar\": 01}\n}",
	"{\n    \"require\": {\"foo/bar\": \"^1.0\"\n}",
	"{\n    \"a b\": {\"x#y\": \"z\"},\n    \"x#y\": [1,2,3],\n    \"\": \"empty\",\n    \"0\": \"zero\",\n    \"\\\"q\": \"quote\"\n}",
	"{\"scripts\":{\"test\":[\"phpunit\",\"@lint\"],\"lint\":\"php -l\"},\"scripts-descriptions\":{\"test\":\"Run\"}}",
	"{   \"name\" :   \"x\"  ,   \"require\"   :   {   \"a/b\"   :   \"1\"   }   ,  \"config\":{\"require\":{\"a/b\":\"2\"}}   }",
	"{\n    \"extra\": {\n        \"require\": {\"nested\": \"value\"},\n        \"branch-alias\": {\"dev-main\": \"1.x-dev\"}\n    },\n    \"require\": {\"top/level\": \"*\"}\n}",
	"{\n    \"description\": \"a \\/ slash and \\u00e9 and \\t tab\",\n    \"autoload\": {\"psr-4\": {\"Foo\\\\\\\\\": \"src/\"}},\n    \"require\": [\"not\", \"an\", \"object\"]\n}",
	"{\"require\": \"string\", \"config\": 1, \"extra\": null, \"repositories\": []}",
	"{\n    \"require\": {\"foo/bar\": \"^1.0\"}, \"require\": {\"dup/key\": \"*\"}\n}",
}

func bigDoc(n int) string {
	var b strings.Builder
	b.WriteString("{\n    \"name\": \"big/doc\",\n    \"require\": {\n")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(",\n")
		}
		fmt.Fprintf(&b, "        \"vendor%d/package-%d\": \"^%d.0\"", i, i, i%7)
	}
	b.WriteString("\n    },\n    \"autoload\": {\n        \"psr-4\": {\n")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(",\n")
		}
		fmt.Fprintf(&b, "            \"Ns%d\\\\Sub\\\\\": \"packages/p%d/src\"", i, i)
	}
	b.WriteString("\n        }\n    },\n    \"extra\": {\"list\": [")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "{\"i\": %d, \"s\": \"v\\u00e9%d\", \"f\": %d.5e-1}", i, i, i)
	}
	b.WriteString("]}\n}")
	return b.String()
}

func manipulatorPatterns() {
	allDocs := append([]string{}, docs...)
	allDocs = append(allDocs, bigDoc(60))
	keys := []string{"require", "require-dev", "config", "repositories", "extra", "name", "autoload", "a b", "x#y", "", "0", "\"q", "a.b", "missing", "scripts", "numbers"}
	for _, k := range keys {
		q := pcre.Quote(enc(k), "")
		pats := []string{
			`{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)` + `(?P<property>` + q + `\s*:\s*)(?P<value>(?&json))(?P<end>.*)}sx`,
			`{` + DEFINES + `^(?P<start> \s* \{ \s* (?: (?&string) \s* : (?&json) \s* , \s* )*?` + q + `\s*:\s*)(?P<content>(?&object))(?P<end>.*)}sx`,
			`{` + DEFINES + `^(?P<start> \s* \{ \s* (?: (?&string) \s* : (?&json) \s* , \s* )*?` + q + `\s*:\s*)(?P<content>(?&array))(?P<end>.*)}sx`,
			`{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)` + `(?P<key>` + q + `\s*:\s*(?&json))(?P<end>.*)}sx`,
			`{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)` + `(?P<removal>` + q + `\s*:\s*(?&json))\s*,?\s*(?P<end>.*)}sx`,
			`{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?` + q + `\s*:\s*)(?P<removal>\{(?P<removal_space>\s*+)\})(?P<end>\s*,?\s*.*)}sx`,
		}
		for _, p := range pats {
			match(p, allDocs...)
		}
		add(Case{Op: "callback", P: pats[1], S: allDocs[3]})
		add(Case{Op: "callback", P: pats[1], S: allDocs[4]})
		// matchChild
		key := strings.NewReplacer(" ", `\ `, "/", `\\?/`).Replace(q)
		child := `{` + DEFINES + `^(?P<start>(?P<prefix>\{ \s* (?: (?!` + key + `\s*:) (?&string) \s* : (?&json) \s* , \s* )*+) ` + key + `\s*:\s*)(?P<content>(?&json))(?P<end>,?)}x`
		match(child, objects()...)
	}
	// package regexes of addLink
	for _, pkg := range []string{"foo/bar", "FOO/Bar", "php", "a/b", "dup/key", "vendor7/package-7", "ext-json"} {
		pq := strings.ReplaceAll(pcre.Quote(pkg, ""), "/", `\\?/`)
		p1 := `{` + DEFINES + `"(?P<package>` + pq + `)"(\s*:\s*)(?&string)}ix`
		p2 := `{` + DEFINES + `"` + pq + `"(?P<separator>\s*:\s*)(?&string)}ix`
		for _, s := range objects() {
			add(Case{Op: "match", P: p1, S: s})
			add(Case{Op: "callback", P: p2, S: s})
		}
	}
	// setRepositoryUrl
	for idx := 0; idx < 4; idx++ {
		lr := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"repositories"\s*:\s*\[\s*((?&json)\s*+,\s*+){` + fmt.Sprint(idx) + `})(?P<repository>(?&object))(?P<end>.*)}sx`
		match(lr, docs[7], docs[8], docs[3])
	}
	for _, name := range []string{"foo", "packagist.org", "0"} {
		or := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"repositories"\s*:\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?` + pcre.Quote(enc(name), "") + `\s*:\s*)(?P<repository>(?&object))(?P<end>.*)}sx`
		match(or, docs[7], docs[8])
	}
	rr := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"url"\s*:\s*)(?P<url>(?&string))(?P<end>.*)}sx`
	both(rr, `{"type": "vcs", "url": "https://example.org/a"}`, `{"url":"x","type":"y"}`, `{"type": "vcs"}`, `{"name": {"url": "nested"}, "url": "top"}`)
	// insertListItem / removeListItem
	lists := []string{"[]", "[1]", "[1, 2, 3]", "[\n    {\"a\": 1},\n    {\"b\": [2, 3]},\n    \"c\"\n]", "[ 1 ,2,  3 , 4 ]", "[\"x\"]", "[{\"packagist.org\": false}, {\"type\": \"vcs\"}]"}
	for idx := 0; idx < 4; idx++ {
		p := `{` + DEFINES + `^(?P<start>\[\s*((?&json)\s*+,\s*?){` + fmt.Sprint(idx) + `})(?P<space_before_item>(\s*))(?P<end>.*)}sx`
		both(p, lists...)
		var startRegex, contentRegex, endRegex string
		contentRegex = `(?&json)`
		if idx > 1 {
			startRegex = `(?&json)\s*+(?:,(?&json)\s*+){` + fmt.Sprint(idx-1) + `}`
			contentRegex = `\s*+,?\s*+` + contentRegex
			endRegex = `(?:(\s*+,\s*+(?&json))*(?:\s*+(?&json))?)\s*+`
		} else if idx > 0 {
			startRegex = `(?&json)\s*+`
			contentRegex = `\s*+,?\s*+` + contentRegex
			endRegex = `(?:(\s*+,\s*+(?&json))*(?:\s*+(?&json))?)\s*+`
		} else {
			startRegex = `\s*+`
			contentRegex = contentRegex + `\s*+,?\s*+`
			endRegex = `(?:((?&json)\s*+,\s*+)*(?:\s*+(?&json))?)\s*+`
		}
		match(`{`+DEFINES+`(?P<start>\[`+startRegex+`)(?P<content>`+contentRegex+`)(?P<end>`+endRegex+`\])}sx`, lists...)
	}
	// small patterns
	small := []string{
		`#^\{(.*)\}$#s`,
		`#^\s*\{\s*\S+.*?(\s*\}\s*)$#s`,
		`#^\{(?P<leadingspace>\s*?)(?P<content>\S+.*?)?(?P<trailingspace>\s*)\}$#s`,
		`#^\[(?P<leadingspace>\s*?)(?P<content>\S+.*?)?(?P<trailingspace>\s*)\]$#s`,
		`#^\{\s*?(?P<content>\S+.*?)?(?P<trailingspace>\s*)\}$#s`,
		`#[^{\s](\s*)\}$#`,
		`{,\s*$}`,
		`{^\s*,?\s*}`,
		`#,\s*$#`,
		`#^\}$#`,
		`#,(\s*)$#`,
		`#^\{\s*\}\s*$#`,
		`#^([ \t]+)"#m`,
		`#^ {4,}#m`,
	}
	smallSubjects := append(append([]string{}, docs...), "{ }", "{  \n  }\n", "[]", "[ ]", "[\n    1\n  ]", "{\"a\": 1,\n    ", "}", "}\n", "{\n    \"a\": 1,\n", "    ,  \n  \"b\"", "{  \"a\": 1 }  ", "[ \"x\" ,\n  \"y\" ]")
	for _, p := range small {
		both(p, smallSubjects...)
	}
	for _, p := range []string{`/^php/`, `/^hhvm/`, `/^ext/`, `/^lib/`, `/^\D/`} {
		both(p, "php", "php-64bit", "hhvm", "ext-json", "lib-icu", "composer-plugin-api", "0-php", "2-ext-x")
	}
	replace(`#^ {4,}#m`, "\t", docs[3])
	both(`{^(bitbucket-oauth|github-oauth|gitlab-oauth|gitlab-token|bearer|http-basic|custom-headers|forgejo-token|platform)\.}`,
		"github-oauth.github.com", "platform.php", "platform", "http-basic.x.y", "bearer", "xbearer.y")
	// replacements as JsonManipulator performs them
	replace("#\n    }$#", pcre.AddCSlashes(",\n        \"new\": \"$1 \\\\ value\"\n    }", `\$`), "{\n        \"a\": \"b\"\n    }")
	replace("#^{\n        #", pcre.AddCSlashes("{\n        \"new\": \"x\",\n        ", `\$`), "{\n        \"a\": \"b\"\n    }")
	replace(`{,\s*$}`, "", "{\n    \"a\": 1,\n    ")
	replace(`{^\s*,?\s*}`, "", "  ,\n    \"b\": 2\n}")
	replace(`#,(\s*)$#`, "$1", "{\n    \"a\": 1,\n")
}

func objects() []string {
	return []string{
		`{}`,
		`{"require": "x"}`,
		"{\n        \"foo/bar\": \"^1.0\",\n        \"foo\\/baz\": \"^2.0\",\n        \"FOO/BAR\": \"dup\"\n    }",
		`{"a b": 1, "x#y": [1, {"require": 2}], "require": {"php": ">=7", "vendor7/package-7": "1"}, "0": null, "": "e", "\"q": "q", "a.b": {"c": 1}}`,
		`{"extra": {"require": 1}, "require": 2}`,
		"{\n    \"name\": \"x\",\n    \"scripts\": {\"test\": \"phpunit\"},\n    \"missing\": [ ]\n}",
		`{"ext-json": "*", "php": "^8.1", "dup/key": "1", "dup\/key": "2"}`,
		`{"autoload": {"psr-4": {"A\\": "src/"}}, "config": {"platform": {"php": "8.1"}}, "repositories": [1, 2]}`,
		`{"bad": 'single', "require": 1}`,
		`{"numbers": [1, 2], "x": 01}`,
	}
}

func quoting() {
	for _, s := range []string{
		`.\+*?[^]$(){}=!<>|:-#/`, "a\x00b", "plain text", "foo/bar", `"require"`, "a b\tc\nd", "ünï", "@%&~`'\",;_",
	} {
		add(Case{Op: "quote", S: s})
		add(Case{Op: "quote", S: s, D: "/"})
		add(Case{Op: "quote", S: s, D: "#"})
		add(Case{Op: "quote", S: s, D: "x"})
	}
}

// LimitCase asks the oracle for the largest n at which preg_match on
// limitDoc(Parts, n) does not fail with a PCRE limit error.
type LimitCase struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Parts   string `json:"parts"`
	Max     int    `json:"max"`
}

func limits() {
	pats := map[string]string{
		"property": `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)(?P<property>"require\-dev"\s*:\s*)(?P<value>(?&json))(?P<end>.*)}sx`,
		"node":     `{` + DEFINES + `^(?P<start> \s* \{ \s* (?: (?&string) \s* : (?&json) \s* , \s* )*?"config"\s*:\s*)(?P<content>(?&object))(?P<end>.*)}sx`,
		"removal":  `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)(?P<removal>"zzz"\s*:\s*(?&json))\s*,?\s*(?P<end>.*)}sx`,
		"nodearr":  `{` + DEFINES + `^(?P<start> \s* \{ \s* (?: (?&string) \s* : (?&json) \s* , \s* )*?"name"\s*:\s*)(?P<content>(?&array))(?P<end>.*)}sx`,
		"present":  `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)(?P<property>"require"\s*:\s*)(?P<value>(?&json))(?P<end>.*)}sx`,
	}
	sets := map[string][]string{
		"property": {"r", "a", "e", "n", "s", "d", "w", "ra", "re", "ae", "rae", "ns", "sd", "wn", "ew", "rk", "kk", "rn"},
		"node":     {"r", "e", "s", "n", "d", "w", "ra", "rk", "kk", "sd", "ew"},
		"removal":  {"r", "e", "s", "n", "d", "rk", "kk"},
		"nodearr":  {"r", "e", "s", "rk"},
		"present":  {"A", "O", "rk", "ek"},
	}
	names := make([]string, 0, len(sets))
	for n := range sets {
		names = append(names, n)
	}
	sort.Strings(names)
	var cs []LimitCase
	for _, n := range names {
		for _, parts := range sets[n] {
			cs = append(cs, LimitCase{Name: n + ":" + parts, Pattern: pats[n], Parts: parts, Max: 40000})
		}
	}
	out, err := json.Marshal(map[string]any{"limits": cs})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))
}
