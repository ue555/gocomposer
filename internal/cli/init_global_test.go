package cli

import (
	"reflect"
	"testing"
)

func TestSanitizePackageNameComponent(t *testing.T) {
	cases := map[string]string{
		"MyCool__Project.Name": "my-cool_project.name",
		"ACMECorp":             "acme-corp",
		"--foo--":              "foo",
		"Ünïcode Dir":          "ncodedir",
		"a..b--c":              "a.b-c",
	}
	for in, want := range cases {
		if got := sanitizePackageNameComponent(in); got != want {
			t.Errorf("sanitize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNamespaceFromPackageName(t *testing.T) {
	cases := map[string]string{
		"new_projects.acme-extra/package-name": `NewProjectsAcmeExtra\PackageName`,
		"acme/my-demo_pkg.x":                   `Acme\MyDemoPkgX`,
		"noslash":                              "",
	}
	for in, want := range cases {
		if got := namespaceFromPackageName(in); got != want {
			t.Errorf("namespace(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestStringInputRoundTrip(t *testing.T) {
	argv := []string{"-n", "global", "config", "extra.foo", "a b 'c'", "--json=[\"x y\"]", "0", "-vvv"}
	s := argvInputString(argv)
	got, err := stringInputTokenize(s)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, argv) {
		t.Errorf("round trip of %q: %q", s, got)
	}
	if got, _ := stringInputTokenize(`a\ b "c\td" e'f'g`); !reflect.DeepEqual(got, []string{"a b", "c\td", "efg"}) {
		t.Errorf("tokenize: %q", got)
	}
}
