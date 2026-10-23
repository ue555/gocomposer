// Command gencorpus writes internal/jsonmanip/testdata/corpus.json: a
// deterministic set of JsonManipulator operation sequences on varied
// composer.json documents, replayed on official PHP by
// oracle/scripts/jsonmanip-dump.php (Composer
// f83588392df43d4fd027d84d26de459946ef4e80) and by the Go port.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Usage: go run ./internal/jsonmanip/testdata/gencorpus > internal/jsonmanip/testdata/corpus.json
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"
)

// tagged values understood by Rec::dec in jsonmanip-dump.php
type kv struct {
	K string
	V any
}

type arr []kv  // PHP array
type obj []kv  // stdClass
type aobj []kv // ArrayObject
type i64 int64
type f64 string

func enc(v any) any {
	switch t := v.(type) {
	case nil, bool:
		return t
	case string:
		if !utf8.ValidString(t) {
			return map[string]any{"t": "bytes", "b64": base64.StdEncoding.EncodeToString([]byte(t))}
		}
		return t
	case i64:
		return map[string]any{"t": "int", "v": int64(t)}
	case f64:
		return map[string]any{"t": "float", "v": string(t)}
	case arr:
		return map[string]any{"t": "array", "items": items(t)}
	case obj:
		return map[string]any{"t": "object", "class": "stdClass", "items": items(t)}
	case aobj:
		return map[string]any{"t": "object", "class": "ArrayObject", "items": items(t)}
	}
	panic(fmt.Sprintf("enc %T", v))
}

func items(l []kv) [][]any {
	out := [][]any{}
	for _, e := range l {
		out = append(out, []any{e.K, enc(e.V)})
	}
	return out
}

func list(vs ...any) arr {
	a := arr{}
	for i, v := range vs {
		a = append(a, kv{fmt.Sprint(i), v})
	}
	return a
}

type op struct {
	M    string `json:"m"`
	Args []any  `json:"args"`
}

type tcase struct {
	Contents any  `json:"contents"`
	Ops      []op `json:"ops"`
}

var docs = []string{
	``,
	`{}`,
	"{\n}",
	"  {  }  ",
	`[]`,
	`"string"`,
	`{"require":{"foo/bar":"^1.0"}}`,
	"{\n    \"name\": \"vendor/pkg\",\n    \"description\": \"A package\",\n    \"require\": {\n        \"php\": \">=7.2\",\n        \"foo/bar\": \"^1.0\"\n    },\n    \"require-dev\": {\n        \"phpunit/phpunit\": \"^9\"\n    }\n}",
	"{\r\n    \"name\": \"vendor/pkg\",\r\n    \"require\": {\r\n        \"foo\\/bar\": \"^1.0\"\r\n    },\r\n    \"config\": {}\r\n}",
	"{\n\t\"config\": {\n\t\t\"platform\": {\"php\": \"8.1\"},\n\t\t\"sort-packages\": true\n\t},\n\t\"extra\": {\"a\": {\"b\": [1, 2, {\"d\": null}]}}\n}",
	"{\n  \"description\": \"Some U\\u00F1icode \\\"quoted\\\" \\\\ back\",\n  \"require\": {\n    \"foo/bar\": \"^1.0\"\n  }\n}",
	"{\n    \"name\": \"ünïcödé/pkg\",\n    \"keywords\": [\"日本語\", \"emoji 🎉\"],\n    \"repositories\": [\n        {\"type\": \"vcs\", \"url\": \"https://example.org/a\"},\n        {\"packagist.org\": false},\n        {\"name\": \"foo\", \"type\": \"composer\", \"url\": \"https://repo.example\"}\n    ]\n}",
	"{\n    \"repositories\": {\n        \"foo\": {\"type\": \"vcs\", \"url\": \"https://example.org/a\"},\n        \"packagist.org\": false,\n        \"0\": {\"type\": \"path\", \"url\": \"../x\", \"options\": {}}\n    }\n}",
	"{\n    \"repositories\": [{\"type\": \"vcs\", \"url\": \"https://a\"}, {\"type\": \"vcs\", \"url\": \"https://b\"}],\n    \"keywords\": [\"a\", \"b\", \"c\"]\n}",
	"{\n    \"a\": 1,\n    \"require\": {},\n    \"config\": {\n    },\n    \"b\": 2\n}",
	"{\n    \"require\": {\n        \"foo/bar\": \"^1.0\",\n    }\n}",
	"{\n    // comment\n    \"require\": {\"foo/bar\": \"^1.0\"}\n}",
	"{\n    \"require\": {\"foo/bar\": \"^1.0\"\n}",
	"{\n    \"name\": \"bad\xffutf8\"\n}",
	"{\n    \"a b\": {\"x#y\": \"z\"},\n    \"x#y\": [1,2,3],\n    \"\": \"empty\",\n    \"0\": \"zero\",\n    \"\\\"q\": \"quote\"\n}",
	"{\"scripts\":{\"test\":[\"phpunit\",\"@lint\"],\"lint\":\"php -l\"},\"scripts-descriptions\":{\"test\":\"Run\"}}",
	"{   \"name\" :   \"x\"  ,   \"require\"   :   {   \"a/b\"   :   \"1\"   }   ,  \"config\":{\"require\":{\"a/b\":\"2\"}}   }",
	"{\n    \"extra\": {\n        \"require\": {\"nested\": \"value\"},\n        \"branch-alias\": {\"dev-main\": \"1.x-dev\"}\n    },\n    \"require\": {\"top/level\": \"*\"}\n}",
	"{\n    \"description\": \"a \\/ slash and \\u00e9 and \\t tab\",\n    \"autoload\": {\"psr-4\": {\"Foo\\\\\\\\\": \"src/\"}, \"files\": [\"a.php\"]},\n    \"suggest\": [\"not\", \"an\", \"object\"]\n}",
	"{\"require\": \"string\", \"config\": 1, \"extra\": null, \"repositories\": [], \"suggest\": {}}",
	"{\n    \"require\": {\"foo/bar\": \"^1.0\"}, \"require\": {\"dup/key\": \"*\"}\n}",
	"{\n    \"config\": {\n        \"policy\": {\n            \"advisories\": {\"audit\": \"report\"},\n            \"abandoned\": {\"block\": true}\n        },\n        \"github-oauth\": {\"github.com\": \"token\"}\n    }\n}",
	"{\n    \"name\": \"x/y\",\n    \"repositories\": [\n        {\n            \"type\": \"composer\",\n            \"url\": \"https://repo.packagist.com/acme/\"\n        },\n        {\n            \"packagist.org\": false\n        }\n    ],\n    \"require\": {\n        \"php\": \"^8.1\",\n        \"ext-json\": \"*\",\n        \"monolog/monolog\": \"^3.0\"\n    }\n}",
	"{\n    \"numbers\": [0, -0, 1.5, -1.5e10, 2E-3, 1e+2, 123456789012345678901234567890],\n    \"extra\": {\"x\": 1.0, \"y\": 1e100}\n}",
	"{\"repositories\":{\"a\":{\"type\":\"vcs\",\"url\":\"u\"},\"b\":false}}",
	"{\n  \"name\": \"two/space\",\n  \"require\": {\n    \"a/a\": \"1\"\n  },\n  \"scripts\": {\n    \"post-install-cmd\": \"echo $1\"\n  }\n}",
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
	b.WriteString("\n    },\n    \"extra\": [")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "{\"i\": %d}", i)
	}
	b.WriteString("]")
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&b, ",\n    \"key%d\": \"v\"", i)
	}
	b.WriteString("\n}")
	return b.String()
}

func deepDoc(n int) string {
	return "{\n    \"extra\": {\"deep\": " + strings.Repeat(`{"a":`, n) + "1" + strings.Repeat("}", n) + "},\n    \"require\": {}\n}"
}

var r = rand.New(rand.NewSource(20261023))

func pick[T any](xs ...T) T { return xs[r.Intn(len(xs))] }

func value() any {
	switch r.Intn(16) {
	case 0:
		return pick[any]("value", "ünï", "a/b", "$1 \\1 ${1}", `back\slash`, "", "0", "quote\"d", "line\nbreak")
	case 1:
		return i64(pick(0, 1, -5, 9000))
	case 2:
		return f64(pick("1.5", "1.0", "1.0E+100", "-0.25"))
	case 3:
		return pick[any](true, false)
	case 4:
		return nil
	case 5:
		return arr{}
	case 6:
		return obj{}
	case 7:
		return aobj{}
	case 8:
		return list("a", "b")
	case 9:
		return arr{{"type", "vcs"}, {"url", "https://example.org/x"}}
	case 10:
		return arr{{"nested", arr{{"deep", list(i64(1), obj{})}}}}
	case 11:
		return obj{{"k", "v"}}
	case 12:
		return arr{{"5", "five"}, {"x", nil}}
	case 13:
		return list(arr{{"a", i64(1)}}, list())
	default:
		return pick[any]("^1.0", "*", "dev-main")
	}
}

func repoConfig() any {
	switch r.Intn(5) {
	case 0:
		return false
	case 1:
		return arr{{"type", "composer"}, {"url", "https://repo.example/" + pick("a", "b")}}
	case 2:
		return arr{{"type", "path"}, {"url", "../pkg"}, {"options", arr{{"symlink", false}}}}
	case 3:
		return arr{{"name", "explicit"}, {"type", "vcs"}, {"url", "git@x:y.git"}}
	default:
		return list("odd")
	}
}

func genOp() op {
	mainNodes := []string{"config", "extra", "scripts", "autoload", "suggest", "require", "repositories", "keywords", "missing"}
	subNames := []string{"sort-packages", "platform.php", "a.b", "foo", "0", "policy.x", "with space", "x/y", "require", "github-oauth.github.com", "deep.a.b"}
	keys := []string{"name", "description", "require", "require-dev", "config", "extra", "repositories", "keywords", "a b", "", "0", "new-key", "minimum-stability", "scripts"}
	repoNames := []string{"", "foo", "packagist.org", "0", "bar", "a", "b", "explicit", "5"}
	switch r.Intn(22) {
	case 0, 1:
		return op{"addLink", []any{pick("require", "require-dev", "conflict", "suggest", "new-type", "config"), pick("foo/bar", "Foo/Bar", "php", "ext-json", "vendor7/package-7", "a b/c", "lib-icu", "composer-plugin-api", "zzz/last", "monolog/monolog"), pick("^1.0", "*", "$1 \\1 ~2", ">=1 <2"), r.Intn(2) == 0}}
	case 2, 3:
		return op{"addSubNode", []any{pick(mainNodes...), pick(subNames...), value(), r.Intn(3) != 0}}
	case 4, 5:
		return op{"removeSubNode", []any{pick(mainNodes...), pick(subNames...)}}
	case 6:
		return op{"addMainKey", []any{pick(keys...), value()}}
	case 7:
		return op{"removeMainKey", []any{pick(keys...)}}
	case 8:
		return op{"removeMainKeyIfEmpty", []any{pick(keys...)}}
	case 9:
		return op{"addConfigSetting", []any{pick("sort-packages", "platform.php", "policy.advisories.block", "policy.abandoned.audit", "policy.new.field", "policy.a.b.c", "github-oauth.github.com", "preferred-install.*", "allow-plugins.foo/bar"), value()}}
	case 10:
		return op{"removeConfigSetting", []any{pick("sort-packages", "platform.php", "policy.advisories.block", "policy.advisories.audit", "policy.abandoned.block", "policy.x.y", "github-oauth.github.com", "platform")}}
	case 11:
		return op{"addProperty", []any{pick("suggest.foo/bar", "extra.a.b", "extra.branch-alias", "scripts.test", "name", "description", "extra.0"), value()}}
	case 12:
		return op{"removeProperty", []any{pick("suggest.foo/bar", "extra.a.b", "extra.require", "scripts.test", "autoload.psr-4", "autoload-dev.files", "name", "keywords", "autoload.files")}}
	case 13, 14:
		return op{"addRepository", []any{pick(repoNames...), repoConfig(), r.Intn(3) != 0}}
	case 15:
		return op{"insertRepository", []any{pick(repoNames...), repoConfig(), pick("foo", "packagist.org", "bar", "explicit", "nope", "0"), i64(pick(0, 1, 2))}}
	case 16:
		return op{"removeRepository", []any{pick(repoNames...)}}
	case 17:
		return op{"setRepositoryUrl", []any{pick(repoNames...), pick("https://new.example/", "$1\\back")}}
	case 18:
		return op{"addListItem", []any{pick("repositories", "keywords", "missing", "require"), value(), r.Intn(2) == 0}}
	case 19:
		return op{"insertListItem", []any{pick("repositories", "keywords", "missing"), value(), i64(pick(0, 1, 2, 3))}}
	case 20:
		return op{"removeListItem", []any{pick("repositories", "keywords", "missing", "extra"), i64(pick(-1, 0, 1, 2, 5))}}
	default:
		return op{"format", []any{value(), i64(pick(0, 1, 2)), r.Intn(2) == 0}}
	}
}

func main() {
	all := append([]string{}, docs...)
	all = append(all, bigDoc(40), deepDoc(1200), deepDoc(30))
	var cases []tcase
	for round := 0; round < 40; round++ {
		for _, d := range all {
			n := 1 + r.Intn(5)
			c := tcase{Contents: enc(d)}
			for i := 0; i < n; i++ {
				c.Ops = append(c.Ops, genOp())
			}
			for i := range c.Ops {
				for j := range c.Ops[i].Args {
					c.Ops[i].Args[j] = enc(c.Ops[i].Args[j])
				}
			}
			cases = append(cases, c)
		}
	}
	out, err := json.Marshal(map[string]any{"cases": cases})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))
}
