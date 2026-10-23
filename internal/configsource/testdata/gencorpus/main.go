// Command gencorpus writes internal/configsource/testdata/corpus.json:
// deterministic Composer\Config\JsonConfigSource operation sequences on
// varied files, replayed on official PHP by
// `oracle/scripts/jsonmanip-dump.php configsource` (Composer
// f83588392df43d4fd027d84d26de459946ef4e80) and by the Go port.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Usage: go run ./internal/configsource/testdata/gencorpus > internal/configsource/testdata/corpus.json
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
	File     string `json:"file"`
	Contents any    `json:"contents"`
	Auth     bool   `json:"auth"`
	Ops      []op   `json:"ops"`
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

var authDocs = []string{
	"{\n    \"github-oauth\": {\n        \"github.com\": \"abc\"\n    },\n    \"http-basic\": {\"repo.example\": {\"username\": \"u\", \"password\": \"p\"}}\n}",
	"{\n}",
	"{\"bearer\": {}}",
	"{\n    \"gitlab-token\": [],\n    \"custom-headers\": {\"x.example\": [\"X-A: 1\"]}\n}",
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
	b.WriteString("\n    },\n    \"config\": [],\n    \"extra\": \"x\"")
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&b, ",\n    \"key%d\": \"v\"", i)
	}
	b.WriteString("\n}")
	return b.String()
}

var r = rand.New(rand.NewSource(20261024))

func pick[T any](xs ...T) T { return xs[r.Intn(len(xs))] }

func value() any {
	switch r.Intn(12) {
	case 0:
		return pick[any]("value", "ünï", "a/b", "$1 \\1", "", "0")
	case 1:
		return i64(pick(0, 1, 9000))
	case 2:
		return f64(pick("1.5", "1.0", "1.0E+100"))
	case 3:
		return pick[any](true, false)
	case 4:
		return nil
	case 5:
		return arr{}
	case 6:
		return obj{}
	case 7:
		return list("a", "b")
	case 8:
		return arr{{"type", "vcs"}, {"url", "https://example.org/x"}}
	case 9:
		return arr{{"nested", arr{{"deep", list(i64(1), obj{})}}}}
	default:
		return pick[any]("^1.0", "*", "dist", "source")
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
		return arr{{"type", "vcs"}, {"url", "https://example.org/" + pick("c", "d")}}
	}
}

func genOp(auth bool) op {
	repoNames := []string{"", "foo", "packagist.org", "packagist", "0", "bar", "explicit", "example_tld"}
	if auth {
		switch r.Intn(4) {
		case 0, 1:
			return op{"addConfigSetting", []any{pick("github-oauth.github.com", "http-basic.repo.example", "bearer.example.org", "gitlab-token.gitlab.com", "platform.php", "github-oauth", "custom-headers.x.example"), pick[any]("tok", arr{{"username", "u"}, {"password", "p"}}, list("X-B: 2"))}}
		default:
			return op{"removeConfigSetting", []any{pick("github-oauth.github.com", "http-basic.repo.example", "bearer.example.org", "gitlab-token.gitlab.com", "custom-headers.x.example", "github-oauth")}}
		}
	}
	switch r.Intn(14) {
	case 0, 1:
		return op{"addRepository", []any{pick(repoNames...), repoConfig(), r.Intn(3) != 0}}
	case 2:
		return op{"insertRepository", []any{pick(repoNames...), repoConfig(), pick("foo", "packagist.org", "bar", "explicit", "nope"), i64(pick(0, 1, 2))}}
	case 3:
		return op{"setRepositoryUrl", []any{pick(repoNames...), pick("https://new.example/", "$1\\back")}}
	case 4:
		return op{"removeRepository", []any{pick(repoNames...)}}
	case 5, 6:
		return op{"addConfigSetting", []any{pick("sort-packages", "platform.php", "policy.advisories.block", "policy.abandoned.audit", "policy.new.field", "policy.a.b.c", "github-oauth.github.com", "preferred-install.*", "allow-plugins.foo/bar", "process-timeout", "http-basic.x.y"), value()}}
	case 7:
		return op{"removeConfigSetting", []any{pick("sort-packages", "platform.php", "policy.advisories.block", "policy.advisories.audit", "policy.abandoned.block", "policy.x.y", "github-oauth.github.com", "platform", "preferred-install")}}
	case 8:
		return op{"addProperty", []any{pick("suggest.foo/bar", "extra.a.b", "extra.branch-alias", "scripts.test", "name", "description", "extra.0", "minimum-stability", "type"), value()}}
	case 9:
		return op{"removeProperty", []any{pick("suggest.foo/bar", "extra.a.b", "extra.require", "scripts.test", "autoload.psr-4", "autoload-dev.files", "name", "keywords", "Autoload.files", "extra")}}
	case 10, 11:
		return op{"addLink", []any{pick("require", "require-dev", "conflict", "suggest", "provide", "replace"), pick("foo/bar", "Foo/Bar", "php", "ext-json", "vendor7/package-7", "monolog/monolog", "my-vend/my-lib"), pick("^1.0", "*", "$1 ~2", ">=1 <2")}}
	default:
		return op{"removeLink", []any{pick("require", "require-dev", "conflict", "suggest", "provide", "replace", "config"), pick("foo/bar", "php", "ext-json", "monolog/monolog", "top/level", "a/b", "phpunit/phpunit")}}
	}
}

func main() {
	all := append([]string{}, docs...)
	all = append(all, bigDoc(30))
	var cases []tcase
	for round := 0; round < 12; round++ {
		for _, d := range all {
			c := tcase{File: "composer.json", Contents: enc(d)}
			for i := 1 + r.Intn(4); i > 0; i-- {
				c.Ops = append(c.Ops, genOp(false))
			}
			cases = append(cases, c)
		}
		for _, d := range authDocs {
			c := tcase{File: "auth.json", Contents: enc(d), Auth: true}
			for i := 1 + r.Intn(4); i > 0; i-- {
				c.Ops = append(c.Ops, genOp(true))
			}
			cases = append(cases, c)
		}
		// files that do not exist yet
		for _, auth := range []bool{false, true} {
			c := tcase{File: pick("config.json", "sub/dir/config.json"), Contents: nil, Auth: auth}
			for i := 1 + r.Intn(3); i > 0; i-- {
				c.Ops = append(c.Ops, genOp(auth))
			}
			cases = append(cases, c)
		}
	}
	for i := range cases {
		for j := range cases[i].Ops {
			for k := range cases[i].Ops[j].Args {
				cases[i].Ops[j].Args[k] = enc(cases[i].Ops[j].Args[k])
			}
		}
	}
	out, err := json.Marshal(map[string]any{"cases": cases})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n"))
}
