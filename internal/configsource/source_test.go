// Port of tests/Composer/Test/Config/JsonConfigSourceTest.php
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80); fixtures are
// copied verbatim under testdata/Fixtures (see testdata/NOTICE.md).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package configsource

import (
	"os"
	"path/filepath"
	"testing"

	"gocomposer/internal/phpjson"
)

func fixturePath(name string) string { return filepath.Join("testdata", "Fixtures", name) }

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertFileEquals(t *testing.T, expected, actual string) {
	t.Helper()
	want, err := os.ReadFile(expected)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}
	if string(want) != string(got) {
		t.Fatalf("file %s differs from %s\nexpected:\n%s\nactual:\n%s", actual, expected, want, got)
	}
}

func assertSameFile(t *testing.T, expected, path string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(got) {
		t.Fatalf("expected:\n%s\nactual:\n%s", expected, got)
	}
}

func m(kv ...any) *phpjson.Map {
	out := phpjson.NewMap()
	for i := 0; i < len(kv); i += 2 {
		out.Set(kv[i].(string), kv[i+1])
	}
	return out
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddRepository(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("composer-repositories.json"), config)
	src := New(config, false)
	must(t, src.AddRepository("example_tld", m("type", "git", "url", "example.tld"), true))

	assertFileEquals(t, fixturePath("config/config-with-exampletld-repository.json"), config)
}

func TestAddRepositoryAsList(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("composer-repositories.json"), config)
	src := New(config, false)
	must(t, src.AddRepository("", m("type", "git", "url", "example.tld"), true))

	assertFileEquals(t, fixturePath("config/config-with-exampletld-repository-as-list.json"), config)
}

func TestAddRepositoryWithOptions(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("composer-repositories.json"), config)
	src := New(config, false)
	must(t, src.AddRepository("example_tld", m(
		"type", "composer",
		"url", "https://example.tld",
		"options", m(
			"ssl", m(
				"local_cert", "/home/composer/.ssl/composer.pem",
			),
		),
	), true))

	assertFileEquals(t, fixturePath("config/config-with-exampletld-repository-and-options.json"), config)
}

func TestRemoveRepository(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("config/config-with-exampletld-repository.json"), config)
	src := New(config, false)
	must(t, src.RemoveRepository("example_tld"))

	assertFileEquals(t, fixturePath("composer-empty.json"), config)
}

func TestAddPackagistRepositoryWithFalseValue(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("composer-repositories.json"), config)
	src := New(config, false)
	must(t, src.AddRepository("packagist", false, true))

	assertFileEquals(t, fixturePath("config/config-with-packagist-false.json"), config)
}

func TestRemovePackagist(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	copyFile(t, fixturePath("config/config-with-packagist-false.json"), config)
	src := New(config, false)
	must(t, src.RemoveRepository("packagist"))

	assertFileEquals(t, fixturePath("composer-empty.json"), config)
}

func TestAddPolicyListFieldPreservesFormatting(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	original := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor", "sort-packages": true,
        "policy": {
            "advisories": {
                "audit": "report"
            },
            "abandoned": {
                "block": true
            }
        }
    }
}
`
	must(t, os.WriteFile(config, []byte(original), 0o644))

	src := New(config, false)
	must(t, src.AddConfigSetting("policy.advisories.block", true))

	expected := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor", "sort-packages": true,
        "policy": {
            "advisories": {
                "audit": "report",
                "block": true
            },
            "abandoned": {
                "block": true
            }
        }
    }
}
`
	assertSameFile(t, expected, config)
}

func TestRemovePolicyListFieldPreservesFormatting(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	original := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor", "sort-packages": true,
        "policy": {
            "advisories": {
                "audit": "report",
                "block": true
            }
        }
    }
}
`
	must(t, os.WriteFile(config, []byte(original), 0o644))

	src := New(config, false)
	must(t, src.RemoveConfigSetting("policy.advisories.audit"))

	expected := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor", "sort-packages": true,
        "policy": {
            "advisories": {
                "block": true
            }
        }
    }
}
`
	assertSameFile(t, expected, config)
}

func TestRemovePolicyListFieldCascadesEmptyAncestors(t *testing.T) {
	config := filepath.Join(t.TempDir(), "composer.json")
	original := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor", "sort-packages": true,
        "policy": {
            "advisories": {
                "block": true
            }
        }
    }
}
`
	must(t, os.WriteFile(config, []byte(original), 0o644))

	src := New(config, false)
	must(t, src.RemoveConfigSetting("policy.advisories.block"))

	expected := `{
    "name": "vendor/pkg", "config": {
        "vendor-dir": "vendor",
        "sort-packages": true,
        "policy": {
        }
    }
}
`
	assertSameFile(t, expected, config)
}

type addLinkCase struct {
	before, typ, name, value, compareAgainst string
}

func addLinkDataArguments(typ, name, value, fixtureBasename, before string) addLinkCase {
	return addLinkCase{before, typ, name, value, fixturePath("addLink/" + fixtureBasename + ".json")}
}

func provideAddLinkData() []addLinkCase {
	empty := fixturePath("composer-empty.json")
	oneOfEverything := fixturePath("composer-one-of-everything.json")
	twoOfEverything := fixturePath("composer-two-of-everything.json")

	return []addLinkCase{
		addLinkDataArguments("require", "my-vend/my-lib", "1.*", "require-from-empty", empty),
		addLinkDataArguments("require", "my-vend/my-lib", "1.*", "require-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("require", "my-vend/my-lib", "1.*", "require-from-twoOfEverything", twoOfEverything),

		addLinkDataArguments("require-dev", "my-vend/my-lib-tests", "1.*", "require-dev-from-empty", empty),
		addLinkDataArguments("require-dev", "my-vend/my-lib-tests", "1.*", "require-dev-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("require-dev", "my-vend/my-lib-tests", "1.*", "require-dev-from-twoOfEverything", twoOfEverything),

		addLinkDataArguments("provide", "my-vend/my-lib-interface", "1.*", "provide-from-empty", empty),
		addLinkDataArguments("provide", "my-vend/my-lib-interface", "1.*", "provide-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("provide", "my-vend/my-lib-interface", "1.*", "provide-from-twoOfEverything", twoOfEverything),

		addLinkDataArguments("suggest", "my-vend/my-optional-extension", "1.*", "suggest-from-empty", empty),
		addLinkDataArguments("suggest", "my-vend/my-optional-extension", "1.*", "suggest-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("suggest", "my-vend/my-optional-extension", "1.*", "suggest-from-twoOfEverything", twoOfEverything),

		addLinkDataArguments("replace", "my-vend/other-app", "1.*", "replace-from-empty", empty),
		addLinkDataArguments("replace", "my-vend/other-app", "1.*", "replace-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("replace", "my-vend/other-app", "1.*", "replace-from-twoOfEverything", twoOfEverything),

		addLinkDataArguments("conflict", "my-vend/my-old-app", "1.*", "conflict-from-empty", empty),
		addLinkDataArguments("conflict", "my-vend/my-old-app", "1.*", "conflict-from-oneOfEverything", oneOfEverything),
		addLinkDataArguments("conflict", "my-vend/my-old-app", "1.*", "conflict-from-twoOfEverything", twoOfEverything),
	}
}

func TestAddLink(t *testing.T) {
	for _, c := range provideAddLinkData() {
		t.Run(filepath.Base(c.compareAgainst), func(t *testing.T) {
			composerJSON := filepath.Join(t.TempDir(), "composer.json")
			copyFile(t, c.before, composerJSON)
			src := New(composerJSON, false)

			must(t, src.AddLink(c.typ, c.name, c.value))

			assertFileEquals(t, c.compareAgainst, composerJSON)
		})
	}
}

type removeLinkCase struct {
	sourceFile, typ, name, compareAgainst string
}

func removeLinkDataArguments(typ, name, fixtureBasename, after string) removeLinkCase {
	if after == "" {
		after = fixturePath("removeLink/" + fixtureBasename + "-after.json")
	}
	return removeLinkCase{fixturePath("removeLink/" + fixtureBasename + ".json"), typ, name, after}
}

func provideRemoveLinkData() []removeLinkCase {
	oneOfEverything := fixturePath("composer-one-of-everything.json")
	twoOfEverything := fixturePath("composer-two-of-everything.json")

	return []removeLinkCase{
		removeLinkDataArguments("require", "my-vend/my-lib", "require-to-empty", ""),
		removeLinkDataArguments("require", "my-vend/my-lib", "require-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("require", "my-vend/my-lib", "require-to-twoOfEverything", twoOfEverything),

		removeLinkDataArguments("require-dev", "my-vend/my-lib-tests", "require-dev-to-empty", ""),
		removeLinkDataArguments("require-dev", "my-vend/my-lib-tests", "require-dev-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("require-dev", "my-vend/my-lib-tests", "require-dev-to-twoOfEverything", twoOfEverything),

		removeLinkDataArguments("provide", "my-vend/my-lib-interface", "provide-to-empty", ""),
		removeLinkDataArguments("provide", "my-vend/my-lib-interface", "provide-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("provide", "my-vend/my-lib-interface", "provide-to-twoOfEverything", twoOfEverything),

		removeLinkDataArguments("suggest", "my-vend/my-optional-extension", "suggest-to-empty", ""),
		removeLinkDataArguments("suggest", "my-vend/my-optional-extension", "suggest-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("suggest", "my-vend/my-optional-extension", "suggest-to-twoOfEverything", twoOfEverything),

		removeLinkDataArguments("replace", "my-vend/other-app", "replace-to-empty", ""),
		removeLinkDataArguments("replace", "my-vend/other-app", "replace-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("replace", "my-vend/other-app", "replace-to-twoOfEverything", twoOfEverything),

		removeLinkDataArguments("conflict", "my-vend/my-old-app", "conflict-to-empty", ""),
		removeLinkDataArguments("conflict", "my-vend/my-old-app", "conflict-to-oneOfEverything", oneOfEverything),
		removeLinkDataArguments("conflict", "my-vend/my-old-app", "conflict-to-twoOfEverything", twoOfEverything),
	}
}

func TestRemoveLink(t *testing.T) {
	for _, c := range provideRemoveLinkData() {
		t.Run(filepath.Base(c.sourceFile), func(t *testing.T) {
			composerJSON := filepath.Join(t.TempDir(), "composer.json")
			copyFile(t, c.sourceFile, composerJSON)
			src := New(composerJSON, false)

			must(t, src.RemoveLink(c.typ, c.name))

			assertFileEquals(t, c.compareAgainst, composerJSON)
		})
	}
}
