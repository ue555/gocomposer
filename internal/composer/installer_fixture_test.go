package composer

// Runner for Composer's installer integration fixtures
// (tests/Composer/Test/Fixtures/installer/*.test, copied unmodified to
// testdata/upstream-installer with a NOTICE). It ports
// tests/Composer/Test/InstallerTest.php::doTestIntegration and the mocks it
// uses (FactoryMock, InstallationManagerMock, InstalledFilesystemRepositoryMock):
// the installation manager only records operations, autoload/scripts are
// off, the local repository is in memory, and the PHP platform comes from a
// snapshot recorded with the oracle PHP 8.4 (testdata/php84-facts.json).
// composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT.

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	"gocomposer/internal/solver"
)

var reSection = regexp.MustCompile(`(?:^|\n*)--([A-Z-]+)--\n`)

func readFixture(t *testing.T, path string) map[string]string {
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	locs := reSection.FindAllStringSubmatchIndex(s, -1)
	data := map[string]string{}
	for i, l := range locs {
		end := len(s)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		data[s[l[2]:l[3]]] = s[l[1]:end]
	}
	return data
}

// formatToRegexp converts PHPUnit's assertStringMatchesFormat placeholders.
func formatToRegexp(format string) *regexp.Regexp {
	var sb strings.Builder
	sb.WriteString(`(?s)\A`)
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '%' && i+1 < len(format) {
			i++
			switch format[i] {
			case 'e':
				sb.WriteString(regexp.QuoteMeta(string(os.PathSeparator)))
			case 's':
				sb.WriteString(`[^\r\n]+`)
			case 'S':
				sb.WriteString(`[^\r\n]*`)
			case 'a':
				sb.WriteString(`.+`)
			case 'A':
				sb.WriteString(`.*`)
			case 'w':
				sb.WriteString(`\s*`)
			case 'i':
				sb.WriteString(`[+-]?\d+`)
			case 'd':
				sb.WriteString(`\d+`)
			case 'x':
				sb.WriteString(`[0-9a-fA-F]+`)
			case 'f':
				sb.WriteString(`[+-]?\.?\d+\.?\d*(?:[Ee][+-]?\d+)?`)
			case 'c':
				sb.WriteString(`.`)
			case '%':
				sb.WriteString(`%`)
			default:
				sb.WriteString(regexp.QuoteMeta("%" + string(format[i])))
			}
			continue
		}
		sb.WriteString(regexp.QuoteMeta(string(c)))
	}
	sb.WriteString(`\z`)
	return regexp.MustCompile(sb.String())
}

// looseEqual is PHPUnit assertEquals for decoded JSON (map order ignored).
func looseEqual(a, b phpjson.Value) bool {
	am, aIsMap := a.(*phpjson.Map)
	bm, bIsMap := b.(*phpjson.Map)
	al, aIsList := a.([]phpjson.Value)
	bl, bIsList := b.([]phpjson.Value)
	switch {
	case (aIsMap || aIsList) && (bIsMap || bIsList):
		if aIsList && bIsList {
			if len(al) != len(bl) {
				return false
			}
			for i := range al {
				if !looseEqual(al[i], bl[i]) {
					return false
				}
			}
			return true
		}
		if phpjson.ArrayLen(a) != phpjson.ArrayLen(b) {
			return false
		}
		ok := true
		phpjson.Each(a, func(k string, v phpjson.Value) {
			var other phpjson.Value
			var found bool
			if bIsMap {
				other, found = bm.Get(k)
			} else {
				phpjson.Each(b, func(k2 string, v2 phpjson.Value) {
					if k2 == k {
						other, found = v2, true
					}
				})
			}
			if !found || !looseEqual(v, other) {
				ok = false
			}
		})
		_ = am
		return ok
	}
	as, _ := phpjson.ScalarString(a)
	bs, _ := phpjson.ScalarString(b)
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return as == bs
}

func TestUpstreamInstallerFixtures(t *testing.T) {
	fixturesDir, err := filepath.Abs("../../testdata/upstream-installer")
	if err != nil {
		t.Fatal(err)
	}
	factsRaw, err := os.ReadFile("testdata/php84-facts.json")
	if err != nil {
		t.Fatal(err)
	}
	var facts platform.Facts
	if err := json.Unmarshal(factsRaw, &facts); err != nil {
		t.Fatal(err)
	}
	// Composer's test suite runs on PHP builds with ext-intl loaded (the
	// php:8.4-cli oracle snapshot has none); add it so lib-icu exists like
	// in that environment (solver-problems*.test expect "wrong version").
	facts.Extensions = append(facts.Extensions, "intl")
	if facts.ExtVersions == nil {
		facts.ExtVersions = map[string]string{}
	}
	facts.ExtVersions["intl"] = facts.PHPVersion
	if facts.Constants == nil {
		facts.Constants = map[string]any{}
	}
	facts.Constants["INTL_ICU_VERSION"] = "74.2"
	files, _ := filepath.Glob(filepath.Join(fixturesDir, "*.test"))
	sort.Strings(files)
	origWD, _ := os.Getwd()
	defer os.Chdir(origWD)
	for _, file := range files {
		name := filepath.Base(file)
		for _, optimizer := range []string{"0", "1"} {
			t.Run(name+"/optimizer="+optimizer, func(t *testing.T) {
				runInstallerFixture(t, file, fixturesDir, &facts, optimizer)
			})
		}
	}
}

func runInstallerFixture(t *testing.T, file, fixturesDir string, facts *platform.Facts, optimizer string) {
	data := readFixture(t, file)
	if cond := strings.TrimSpace(data["CONDITION"]); cond != "" {
		switch cond {
		case "!defined('HHVM_VERSION')":
		case "putenv('COMPOSER_FUND=0')":
			t.Setenv("COMPOSER_FUND", "0")
		case "putenv('COMPOSER_FUND=1')":
			t.Setenv("COMPOSER_FUND", "1")
		default:
			t.Skipf("unsupported CONDITION %q", cond)
		}
	}
	t.Setenv("COMPOSER_POOL_OPTIMIZER", optimizer)
	if data["EXPECT-OUTPUT-OPTIMIZED"] != "" && optimizer == "1" {
		data["EXPECT-OUTPUT"] = data["EXPECT-OUTPUT-OPTIMIZED"]
	}
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	os.MkdirAll(home, 0o755)
	t.Setenv("COMPOSER_HOME", home)
	t.Setenv("COMPOSER_CACHE_DIR", filepath.Join(tmp, "cache"))
	t.Setenv("COMPOSER", "")
	os.Unsetenv("COMPOSER")
	// FactoryMock::createConfig merges 'repositories' => ['packagist' => false]
	if err := os.WriteFile(filepath.Join(home, "config.json"), []byte(`{"repositories": {"packagist": false}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	composerData, err := phpjson.Decode([]byte(data["COMPOSER"]))
	if err != nil {
		t.Fatalf("COMPOSER: %v", err)
	}
	cm := phpjson.AsMap(composerData)
	phpjson.Each(get(cm, "repositories"), func(_ string, r phpjson.Value) {
		rm := phpjson.AsMap(r)
		if rm == nil || rm.Str("type") != "composer" {
			return
		}
		if u := rm.Str("url"); regexp.MustCompile(`^file://[^/]`).MatchString(u) {
			rm.Set("url", "file://"+filepath.ToSlash(fixturesDir)+"/"+u[7:])
		}
	})
	// the Locker of InstallerTest hashes json_encode($composerConfig)
	enc, _ := phpjson.Encode(cm, 0)
	project := filepath.Join(tmp, "project")
	os.MkdirAll(filepath.Join(project, "vendor", "composer"), 0o755)
	if err := os.WriteFile(filepath.Join(project, "composer.json"), []byte(enc), 0o644); err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(project, "composer.lock")
	if strings.TrimSpace(data["LOCK"]) != "" {
		lv, err := phpjson.Decode([]byte(data["LOCK"]))
		if err != nil {
			t.Fatalf("LOCK: %v", err)
		}
		lm := phpjson.AsMap(lv)
		if !lm.Has("hash") {
			flat, _ := phpjson.Encode(cm, 0)
			sum := md5.Sum([]byte(flat))
			lm.Set("hash", hex.EncodeToString(sum[:]))
		}
		le, _ := phpjson.Encode(lm, phpjson.PrettyPrint)
		os.WriteFile(lockPath, []byte(le), 0o644)
	}
	if strings.TrimSpace(data["INSTALLED"]) != "" {
		os.WriteFile(filepath.Join(project, "vendor", "composer", "installed.json"), []byte(data["INSTALLED"]), 0o644)
	}
	// InstallerTest runs from tests/Composer/Test, where relative path
	// repositories point to Fixtures/functional/...
	for _, d := range []string{"Fixtures", "DependencyResolver"} {
		if err := os.Symlink(filepath.Join(fixturesDir, d), filepath.Join(project, d)); err != nil {
			t.Fatal(err)
		}
	}
	lockBefore, _ := os.ReadFile(lockPath)
	if err := os.Chdir(project); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	io := &cio.IO{Out: &out, Err: &out, Verbosity: cio.Normal, SymfonyStylesOnly: true}
	c, err := Create(io, Options{FullLoad: true, PlatformFacts: facts, DisablePlugins: true, SkipSchemaValidation: true, SkipPurge: true})
	expectException := strings.TrimSpace(data["EXPECT-EXCEPTION"]) != ""
	if err != nil {
		if expectException && strings.Contains(err.Error(), strings.TrimRight(data["EXPECT"], " \n")) {
			return
		}
		t.Fatalf("create: %v\n%s", err, out.String())
	}
	var trace []string
	inst := NewInstaller(c)
	inst.IM.Trace = &trace
	inst.IM.InstallPathOverride = func(p *pkg.Package) string { return "vendor/" + p.Name }
	inst.DevMode = true
	inst.RunScripts = false
	inst.MockAutoloadGenerator = true
	inst.Audit = false
	args := strings.Fields(data["RUN"])
	cmd := args[0]
	var packages []string
	var ignoreList []string
	ignoreAll := false
	for _, a := range args[1:] {
		switch {
		case a == "--no-dev":
			inst.DevMode = false
		case a == "--dry-run":
			inst.DryRun = true
		case a == "--ignore-platform-reqs":
			ignoreAll = true
		case strings.HasPrefix(a, "--ignore-platform-req="):
			ignoreList = append(ignoreList, strings.TrimPrefix(a, "--ignore-platform-req="))
		case a == "--no-install":
			inst.Install = false
		case a == "--lock":
			inst.UpdateMirrors = true
		case a == "--with-all-dependencies":
			inst.UpdateAllowTransitive = solver.UpdateListedWithTransitiveDeps
		case a == "--with-dependencies":
			if inst.UpdateAllowTransitive != solver.UpdateListedWithTransitiveDeps {
				inst.UpdateAllowTransitive = solver.UpdateListedWithTransitiveDepsNoRoot
			}
		case a == "--minimal-changes":
			inst.MinimalUpdate = true
		case a == "--prefer-stable":
			inst.PreferStable = true
		case a == "--prefer-lowest":
			inst.PreferLowest = true
		case a == "-v" || a == "-vv" || a == "-vvv":
			// the Application output verbosity; the installer IO stays NORMAL
		case strings.HasPrefix(a, "-"):
			t.Fatalf("unsupported RUN option %s", a)
		default:
			packages = append(packages, a)
		}
	}
	switch {
	case ignoreAll:
		inst.PlatformFilter = solver.IgnoreAll{}
	case len(ignoreList) > 0:
		inst.PlatformFilter = solver.NewIgnoreList(ignoreList)
	}
	if cmd == "update" {
		inst.Update = true
		var filtered []string
		for _, p := range packages {
			if p != "lock" && p != "nothing" && p != "mirrors" {
				filtered = append(filtered, p)
			}
		}
		if len(filtered) != len(packages) {
			inst.UpdateMirrors = true
		}
		if len(filtered) > 0 {
			inst.UpdateAllowList = filtered
			inst.UpdateAllowListSet = true
		}
	} else if cmd != "install" {
		t.Fatalf("The run command only supports install and update")
	}
	code, runErr := inst.Run()
	output := strings.ReplaceAll(out.String(), "\r", "")
	if expectException {
		want := strings.TrimRight(data["EXPECT"], " \n")
		if runErr == nil || !strings.Contains(runErr.Error(), want) {
			t.Fatalf("expected exception %s %q, got %v (exit %d)\n%s", data["EXPECT-EXCEPTION"], want, runErr, code, output)
		}
		return
	}
	if runErr != nil && code == 0 {
		// the exception is rendered by the console application, not the IO
		code = 1
	}
	wantCode := 0
	if s := strings.TrimSpace(data["EXPECT-EXIT-CODE"]); s != "" {
		if v, err := phpjson.Decode([]byte(s)); err == nil {
			if n, ok := v.(int64); ok {
				wantCode = int(n)
			}
		}
	}
	if code != wantCode {
		t.Errorf("exit code %d, want %d (error: %v)\n%s", code, wantCode, runErr, output)
	}
	lockAfter, lockErr := os.ReadFile(lockPath)
	written := lockErr == nil && !bytes.Equal(lockBefore, lockAfter)
	if el := strings.TrimSpace(data["EXPECT-LOCK"]); el != "" {
		if el == "false" {
			if written {
				t.Errorf("lock written although EXPECT-LOCK is false")
			}
		} else {
			if !written && lockBefore == nil {
				t.Errorf("lock not written")
			}
			want, err := phpjson.Decode([]byte(el))
			if err != nil {
				t.Fatalf("EXPECT-LOCK: %v", err)
			}
			got, err := phpjson.Decode(lockAfter)
			if err != nil {
				t.Fatalf("lock: %v", err)
			}
			gm := phpjson.AsMap(got)
			for _, k := range []string{"hash", "content-hash", "_readme", "plugin-api-version"} {
				gm.Delete(k)
			}
			if !looseEqual(want, gm) {
				ge, _ := phpjson.Encode(gm, phpjson.ComposerFlags)
				we, _ := phpjson.Encode(want, phpjson.ComposerFlags)
				t.Errorf("lock mismatch\n got: %s\nwant: %s", ge, we)
			}
		}
	}
	if ei := strings.TrimSpace(data["EXPECT-INSTALLED"]); ei != "" {
		var actual []phpjson.Value
		for _, p := range c.LocalRepo.CanonicalPackages() {
			d := pkg.Dump(p)
			d.Delete("version_normalized")
			actual = append(actual, d)
		}
		sort.SliceStable(actual, func(i, j int) bool {
			return phpjson.AsMap(actual[i]).Str("name") < phpjson.AsMap(actual[j]).Str("name")
		})
		if actual == nil {
			actual = []phpjson.Value{}
		}
		want, err := phpjson.Decode([]byte(ei))
		if err != nil {
			t.Fatalf("EXPECT-INSTALLED: %v", err)
		}
		ge, _ := phpjson.Encode(actual, phpjson.ComposerFlags)
		we, _ := phpjson.Encode(want, phpjson.ComposerFlags)
		if ge != we {
			t.Errorf("installed mismatch\n got: %s\nwant: %s", ge, we)
		}
	}
	if got, want := strings.Join(trace, "\n"), strings.TrimRight(data["EXPECT"], " \t\n\r\x00\x0B"); got != want {
		t.Errorf("trace mismatch\n got:\n%s\nwant:\n%s\noutput:\n%s", got, want, output)
	}
	if eo := data["EXPECT-OUTPUT"]; eo != "" {
		output = regexp.MustCompile(`(?m)^    - .*?\.ini$`).ReplaceAllString(output, "__inilist__")
		output = regexp.MustCompile(`(__inilist__\r?\n)+`).ReplaceAllString(output, "__inilist__\n")
		re := formatToRegexp(strings.TrimRight(eo, " \t\n\r\x00\x0B"))
		if !re.MatchString(strings.TrimRight(output, " \t\n\r\x00\x0B")) {
			t.Errorf("output mismatch\n got:\n%s\nwant (format):\n%s", output, eo)
		}
	}
}

func get(m *phpjson.Map, k string) phpjson.Value {
	if m == nil {
		return nil
	}
	v, _ := m.Get(k)
	return v
}
