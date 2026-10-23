// Ports of Composer f83588392df43d4fd027d84d26de459946ef4e80
// tests/Composer/Test/Advisory/AuditorTest.php and CooldownAuditorTest.php
// (expectations unchanged). PHPUnit mocks of ComposerRepository,
// RepositorySet and FilterListProviderSet are Go fakes implementing the
// same methods; IOMock expectations are checked with the same matching rules.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

const mockAdvisoriesJSON = `{
 "vendor1/package1": [
  {"advisoryId": "ID1", "packageName": "vendor1/package1", "title": "advisory1", "link": "https://advisory.example.com/advisory1", "cve": "CVE1", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source1", "remoteId": "RemoteID1"}], "reportedAt": "2022-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"},
  {"advisoryId": "ID4", "packageName": "vendor1/package1", "title": "advisory4", "link": "https://advisory.example.com/advisory4", "cve": "CVE3", "affectedVersions": ">=8,<8.2.2|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteID4"}], "reportedAt": "2022-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "high"},
  {"advisoryId": "ID5", "packageName": "vendor1/package1", "title": "advisory5", "link": "https://advisory.example.com/advisory5", "cve": "", "affectedVersions": ">=8,<8.2.2|>=1,<2.5.6", "sources": [{"name": "source1", "remoteId": "RemoteID3"}], "reportedAt": "2022-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ],
 "vendor1/package2": [
  {"advisoryId": "ID2", "packageName": "vendor1/package2", "title": "advisory2", "link": "https://advisory.example.com/advisory2", "cve": "", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source1", "remoteId": "RemoteID2"}], "reportedAt": "2022-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ],
 "vendorx/packagex": [
  {"advisoryId": "IDx", "packageName": "vendorx/packagex", "title": "advisory17", "link": "https://advisory.example.com/advisory17", "cve": "CVE5", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteIDx"}], "reportedAt": "2015-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ],
 "vendor2/package1": [
  {"advisoryId": "ID3", "packageName": "vendor2/package1", "title": "advisory3", "link": "https://advisory.example.com/advisory3", "cve": "CVE2", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteID1"}], "reportedAt": "2022-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"},
  {"advisoryId": "ID6", "packageName": "vendor2/package1", "title": "advisory6", "link": "https://advisory.example.com/advisory6", "cve": "CVE4", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteID3"}], "reportedAt": "2015-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ],
 "vendory/packagey": [
  {"advisoryId": "IDy", "packageName": "vendory/packagey", "title": "advisory7", "link": "https://advisory.example.com/advisory7", "cve": "CVE5", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteID4"}], "reportedAt": "2015-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ],
 "vendor3/package1": [
  {"advisoryId": "ID7", "packageName": "vendor3/package1", "title": "advisory7", "link": "https://advisory.example.com/advisory7", "cve": "CVE5", "affectedVersions": ">=3,<3.4.3|>=1,<2.5.6", "sources": [{"name": "source2", "remoteId": "RemoteID4"}], "reportedAt": "2015-05-25 13:21:00", "composerRepository": "https://packagist.org", "severity": "medium"}
 ]
}`

// mockComposerRepo is the ComposerRepository mock of AuditorTest::getRepoSet.
type mockComposerRepo struct {
	*repo.Array
	data phpjson.Value
}

func (m *mockComposerRepo) HasSecurityAdvisories() (bool, error) { return true, nil }

func (m *mockComposerRepo) GetSecurityAdvisories(pcm *ConstraintMap, allowPartial bool) (*AdvisoryMap, error) {
	out := NewOMap[[]*Advisory]()
	var err error
	eachArr(m.data, func(name string, list phpjson.Value) {
		c, ok := pcm.Get(name)
		if !ok || err != nil {
			return
		}
		var kept []*Advisory
		eachArr(list, func(_ string, data phpjson.Value) {
			a, e := CreateAdvisory(name, data)
			if e != nil {
				err = e
				return
			}
			if !allowPartial && !a.Full {
				err = newErr(RuntimeException, "Advisory for "+name+" could not be loaded as a full advisory from test repo")
				return
			}
			if a.AffectedVersions.Matches(c) {
				kept = append(kept, a)
			}
		})
		if len(kept) > 0 {
			out.Set(name, kept)
		}
	})
	return out, err
}

func getRepoSet(t *testing.T) AdvisorySource {
	return &RepositoryAdvisorySource{Repositories: []repo.Repository{&mockComposerRepo{Array: repo.NewArray(), data: j(t, mockAdvisoriesJSON)}}}
}

type createPolicyOpts struct {
	abandoned              string
	ignoreUnreachableAudit bool
	filteredAudit          string
	filteredListName       string
	ignoreAdvisories       [][2]*string // key, reason
	ignoreSeverities       []string
	ignoreAbandoned        [][2]*string
}

// createPolicyConfig is AuditorTest::createPolicyConfig.
func createPolicyConfig(o createPolicyOpts) *PolicyConfig {
	if o.abandoned == "" {
		o.abandoned = AuditIgnore
	}
	if o.filteredListName == "" {
		o.filteredListName = "test-list"
	}
	custom := NewOMap[*CustomListPolicyConfig]()
	if o.filteredAudit != "" {
		custom.Set(o.filteredListName, NewCustomListPolicyConfig(o.filteredListName, false, o.filteredAudit, BlockScopeUpdate, nil, nil))
	}
	pkgRules := NewOMap[[]*IgnorePackageRule]()
	ids := NewOMap[*IgnoreIdRule]()
	for _, kv := range o.ignoreAdvisories {
		if strings.Contains(*kv[0], "/") {
			appendTo(pkgRules, *kv[0], NewIgnorePackageRule(*kv[0], &semver.MatchAll{}, kv[1], true, true))
		} else {
			ids.Set(*kv[0], NewIgnoreIdRule(*kv[0], kv[1], true, true))
		}
	}
	sev := NewOMap[*IgnoreSeverityRule]()
	for _, s := range o.ignoreSeverities {
		sev.Set(s, NewIgnoreSeverityRule(s, nil, true, true))
	}
	abIgnore := NewOMap[[]*IgnorePackageRule]()
	for _, kv := range o.ignoreAbandoned {
		appendTo(abIgnore, *kv[0], NewIgnorePackageRule(*kv[0], &semver.MatchAll{}, kv[1], true, true))
	}
	iu := IgnoreUnreachableDefault()
	if o.ignoreUnreachableAudit {
		iu = IgnoreUnreachableAll()
	}
	return NewPolicyConfig(true, NewAdvisoriesPolicyConfig(false, AuditFail, pkgRules, ids, sev), MalwarePolicyConfigDisabled(),
		NewAbandonedPolicyConfig(false, o.abandoned, abIgnore), CooldownPolicyConfigDisabled(), custom, iu)
}

func abandonedPackages(t *testing.T) (*pkg.Package, *pkg.Package) {
	withReplacement := loadPkg(t, "1.0.0", "1.0.0", `{"name": "vendor/abandoned", "version": "1.0.0", "abandoned": "foo/bar"}`)
	noReplacement := loadPkg(t, "1.0.0", "1.0.0", `{"name": "vendor/abandoned2", "version": "1.0.0", "abandoned": true}`)
	return withReplacement, noReplacement
}

func TestAudit(t *testing.T) {
	withReplacement, noReplacement := abandonedPackages(t)
	cases := []struct {
		name            string
		packages        []*pkg.Package
		warningOnly     bool
		abandoned       string
		ignoreAbandoned [][2]*string
		format          string
		expected        int
		output          string
	}{
		{"Test no advisories returns 0", pkgs(newPkg("vendor1/package2", "9.0.0", "9.0.0"), newPkg("vendor1/package1", "9.0.0", "9.0.0"), newPkg("vendor3/package1", "9.0.0", "9.0.0")), true, "", nil, "", StatusOK,
			`No security vulnerability advisories found.`},
		{"Test with advisories returns 1", pkgs(newPkg("vendor1/package2", "9.0.0", "9.0.0"), newPkg("vendor1/package1", "8.2.1", "8.2.1"), newPkg("vendor3/package1", "9.0.0", "9.0.0")), true, "", nil, "", StatusFailed,
			`<warning>Found 2 security vulnerability advisories affecting 1 package:</warning>
Package: vendor1/package1
Severity: high
Advisory ID: ID4
CVE: CVE3
Title: advisory4
URL: https://advisory.example.com/advisory4
Affected versions: >=8,<8.2.2|>=1,<2.5.6
Reported at: 2022-05-25T13:21:00+00:00
--------
Package: vendor1/package1
Severity: medium
Advisory ID: ID5
CVE: ` + `
Title: advisory5
URL: https://advisory.example.com/advisory5
Affected versions: >=8,<8.2.2|>=1,<2.5.6
Reported at: 2022-05-25T13:21:00+00:00`},
		{"abandoned packages ignored", pkgs(withReplacement, noReplacement), false, AuditIgnore, nil, "", StatusOK, `No security vulnerability advisories found.`},
		{"abandoned packages individually ignored via full vendor", pkgs(withReplacement, noReplacement), false, AuditFail, [][2]*string{{sp("vendor/*"), nil}}, "", StatusOK, `No security vulnerability advisories found.`},
		{"abandoned packages individually ignored via package name", pkgs(withReplacement, noReplacement), false, AuditFail, [][2]*string{{sp(withReplacement.Name), nil}, {sp(noReplacement.Name), nil}}, "", StatusOK, `No security vulnerability advisories found.`},
		{"abandoned packages individually ignored not matching package name", pkgs(withReplacement, noReplacement), false, AuditFail, [][2]*string{{sp("acme/test"), sp("ignoring because yolo")}}, "", StatusFailed,
			`No security vulnerability advisories found.
Found 2 abandoned packages:
vendor/abandoned is abandoned. Use foo/bar instead.
vendor/abandoned2 is abandoned. No replacement was suggested.`},
		{"abandoned packages reported only", pkgs(withReplacement, noReplacement), true, AuditReport, nil, "", StatusOK,
			`No security vulnerability advisories found.
Found 2 abandoned packages:
vendor/abandoned is abandoned. Use foo/bar instead.
vendor/abandoned2 is abandoned. No replacement was suggested.`},
		{"abandoned packages fails", pkgs(withReplacement, noReplacement), false, AuditFail, nil, FormatTable, StatusFailed,
			`No security vulnerability advisories found.
Found 2 abandoned packages:
+-------------------+----------------------------------------------------------------------------------+
| Abandoned Package | Suggested Replacement                                                            |
+-------------------+----------------------------------------------------------------------------------+
| vendor/abandoned  | foo/bar                                                                          |
| vendor/abandoned2 | none                                                                             |
+-------------------+----------------------------------------------------------------------------------+`},
		{"vulnerable and abandoned packages fails", pkgs(newPkg("vendor1/package1", "8.2.1", "8.2.1"), withReplacement, noReplacement), false, AuditFail, nil, FormatTable, StatusFailed,
			`Found 2 security vulnerability advisories affecting 1 package:
+-------------------+----------------------------------------------------------------------------------+
| Package           | vendor1/package1                                                                 |
| Severity          | high                                                                             |
| Advisory ID       | ID4                                                                              |
| CVE               | CVE3                                                                             |
| Title             | advisory4                                                                        |
| URL               | https://advisory.example.com/advisory4                                           |
| Affected versions | >=8,<8.2.2|>=1,<2.5.6                                                            |
| Reported at       | 2022-05-25T13:21:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
+-------------------+----------------------------------------------------------------------------------+
| Package           | vendor1/package1                                                                 |
| Severity          | medium                                                                           |
| Advisory ID       | ID5                                                                              |
| CVE               |                                                                                  |
| Title             | advisory5                                                                        |
| URL               | https://advisory.example.com/advisory5                                           |
| Affected versions | >=8,<8.2.2|>=1,<2.5.6                                                            |
| Reported at       | 2022-05-25T13:21:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
Found 2 abandoned packages:
+-------------------+----------------------------------------------------------------------------------+
| Abandoned Package | Suggested Replacement                                                            |
+-------------------+----------------------------------------------------------------------------------+
| vendor/abandoned  | foo/bar                                                                          |
| vendor/abandoned2 | none                                                                             |
+-------------------+----------------------------------------------------------------------------------+`},
		{"abandoned packages fails with json format", pkgs(withReplacement, noReplacement), false, AuditFail, nil, FormatJSON, StatusFailed,
			`{
    "advisories": [],
    "abandoned": {
        "vendor/abandoned": "foo/bar",
        "vendor/abandoned2": null
    },
    "filter": [],
    "cooldown": []
}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.format == "" {
				tc.format = FormatPlain
			}
			pc := createPolicyConfig(createPolicyOpts{abandoned: tc.abandoned, ignoreAbandoned: tc.ignoreAbandoned})
			io := &bufferIO{}
			result, err := NewAuditor(nil).Audit(io, getRepoSet(t), pc, tc.packages, tc.format, tc.warningOnly, nil)
			if err != nil {
				t.Fatal(err)
			}
			if result != tc.expected {
				t.Fatalf("expected %d got %d", tc.expected, result)
			}
			if got := io.trimmed(); got != tc.output {
				t.Fatalf("output mismatch\nexpected:\n%s\ngot:\n%s", tc.output, got)
			}
		})
	}
}

// assertIOExpects ports IOMock::assertComplete for text expectations.
func assertIOExpects(t *testing.T, output string, expected []string, strict bool) {
	t.Helper()
	lines := regexp.MustCompile("\r?\n").Split(output, -1)
	for _, exp := range expected {
		matched := false
		for len(lines) > 0 {
			line := lines[0]
			lines = lines[1:]
			if line == exp {
				matched = true
				break
			}
			if strict {
				t.Fatalf("IO output mismatch. Expected:\n%s\nGot:\n%s\nComplete output:\n%s", exp, line, output)
			}
		}
		if !matched {
			t.Fatalf("Expected %q to be output still but there is no output left to consume. Complete output:\n%s", exp, output)
		}
	}
}

func advisoryLines(pkgName, sev, id, cve, title, url, affected, reported string, extra ...string) []string {
	return append([]string{"Package: " + pkgName, "Severity: " + sev, "Advisory ID: " + id, "CVE: " + cve, "Title: " + title, "URL: " + url, "Affected versions: " + affected, "Reported at: " + reported}, extra...)
}

func TestAuditWithIgnore(t *testing.T) {
	adv1 := func(extra ...string) []string {
		return advisoryLines("vendor1/package1", "medium", "ID1", "CVE1", "advisory1", "https://advisory.example.com/advisory1", ">=3,<3.4.3|>=1,<2.5.6", "2022-05-25T13:21:00+00:00", extra...)
	}
	one := "Found 1 ignored security vulnerability advisory affecting 1 package:"
	cases := []struct {
		name     string
		packages []*pkg.Package
		ignored  [][2]*string
		exitCode int
		expected []string
	}{
		{"ignore by CVE", pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0")), [][2]*string{{sp("CVE1"), nil}}, 0, append([]string{one}, adv1()...)},
		{"ignore by CVE with reasoning", pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0")), [][2]*string{{sp("CVE1"), sp("A good reason")}}, 0, append([]string{one}, adv1("Ignore reason: A good reason")...)},
		{"ignore by advisory id", pkgs(newPkg("vendor1/package2", "3.0.0.0", "3.0.0")), [][2]*string{{sp("ID2"), nil}}, 0,
			append([]string{one}, advisoryLines("vendor1/package2", "medium", "ID2", "", "advisory2", "https://advisory.example.com/advisory2", ">=3,<3.4.3|>=1,<2.5.6", "2022-05-25T13:21:00+00:00")...)},
		{"ignore by remote id", pkgs(newPkg("vendorx/packagex", "3.0.0.0", "3.0.0")), [][2]*string{{sp("RemoteIDx"), nil}}, 0,
			append([]string{one}, advisoryLines("vendorx/packagex", "medium", "IDx", "CVE5", "advisory17", "https://advisory.example.com/advisory17", ">=3,<3.4.3|>=1,<2.5.6", "2015-05-25T13:21:00+00:00")...)},
		{"ignore by package name", pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0")), [][2]*string{{sp("vendor1/package1"), nil}}, 0, append([]string{one}, adv1()...)},
		{"ignore by package name with reasoning", pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0")), [][2]*string{{sp("vendor1/package1"), sp("Package has known safe usage")}}, 0,
			append([]string{one}, adv1("Ignore reason: Package has known safe usage")...)},
		{"1 vulnerability, 0 ignored", pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0")), nil, 1,
			append([]string{"Found 1 security vulnerability advisory affecting 1 package:"}, adv1()...)},
		{"1 vulnerability, 3 ignored affecting 2 packages", pkgs(newPkg("vendor3/package1", "3.0.0.0", "3.0.0"), newPkg("vendorx/packagex", "3.0.0.0", "3.0.0"), newPkg("vendor2/package1", "3.0.0.0", "3.0.0")),
			[][2]*string{{sp("RemoteIDx"), nil}, {sp("ID3"), nil}, {sp("ID6"), nil}}, 1,
			concat([]string{"Found 3 ignored security vulnerability advisories affecting 2 packages:"},
				advisoryLines("vendor2/package1", "medium", "ID3", "CVE2", "advisory3", "https://advisory.example.com/advisory3", ">=3,<3.4.3|>=1,<2.5.6", "2022-05-25T13:21:00+00:00", "Ignore reason: None specified"),
				[]string{"--------"},
				advisoryLines("vendor2/package1", "medium", "ID6", "CVE4", "advisory6", "https://advisory.example.com/advisory6", ">=3,<3.4.3|>=1,<2.5.6", "2015-05-25T13:21:00+00:00", "Ignore reason: None specified"),
				[]string{"--------"},
				advisoryLines("vendorx/packagex", "medium", "IDx", "CVE5", "advisory17", "https://advisory.example.com/advisory17", ">=3,<3.4.3|>=1,<2.5.6", "2015-05-25T13:21:00+00:00", "Ignore reason: None specified"),
				[]string{"Found 1 security vulnerability advisory affecting 1 package:"},
				advisoryLines("vendor3/package1", "medium", "ID7", "CVE5", "advisory7", "https://advisory.example.com/advisory7", ">=3,<3.4.3|>=1,<2.5.6", "2015-05-25T13:21:00+00:00"))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pc := createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, ignoreAdvisories: tc.ignored})
			io := &bufferIO{}
			result, err := NewAuditor(nil).Audit(io, getRepoSet(t), pc, tc.packages, FormatPlain, false, nil)
			if err != nil {
				t.Fatal(err)
			}
			assertIOExpects(t, io.Output(), tc.expected, true)
			if result != tc.exitCode {
				t.Fatalf("exit %d", result)
			}
		})
	}
}

func concat(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func TestAuditWithIgnoreSeverity(t *testing.T) {
	cases := []struct {
		name     string
		sev      []string
		exitCode int
		first    string
	}{
		{"ignore medium", []string{"medium"}, 1, "Found 2 ignored security vulnerability advisories affecting 1 package:"},
		{"ignore high", []string{"high"}, 1, "Found 1 ignored security vulnerability advisory affecting 1 package:"},
		{"ignore high and medium", []string{"high", "medium"}, 0, "Found 3 ignored security vulnerability advisories affecting 1 package:"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pc := createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, ignoreSeverities: tc.sev})
			io := &bufferIO{}
			result, err := NewAuditor(nil).Audit(io, getRepoSet(t), pc, pkgs(newPkg("vendor1/package1", "2.0.0.0", "2.0.0")), FormatPlain, false, nil)
			if err != nil {
				t.Fatal(err)
			}
			assertIOExpects(t, io.Output(), []string{tc.first}, true)
			if result != tc.exitCode {
				t.Fatalf("exit %d", result)
			}
		})
	}
}

// unreachableRepoSet is the RepositorySet mock of testAuditWithIgnoreUnreachable.
type unreachableRepoSet struct{ msg string }

func (u *unreachableRepoSet) GetMatchingSecurityAdvisories(_ []*pkg.Package, _, ignoreUnreachable bool) (*AdvisoryResult, error) {
	if !ignoreUnreachable {
		return nil, &httpdl.TransportError{Msg: u.msg, StatusCode: 404}
	}
	reported := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	m := NewOMap[[]*Advisory]()
	m.Set("vendor1/package1", []*Advisory{
		NewSecurityAdvisory("vendor1/package1", "CVE-2023-12345", semver.NewConstraint("=", "3.0.0.0"), "First repo advisory", []phpjson.Value{map1("name", "test", "remoteId", "1")}, reported, sp("CVE-2023-12345"), sp("https://example.com/advisory/1"), sp("medium")),
		NewSecurityAdvisory("vendor1/package1", "CVE-2023-67890", semver.NewConstraint("=", "3.0.0.0"), "Third repo advisory", []phpjson.Value{map1("name", "test", "remoteId", "3")}, reported, sp("CVE-2023-67890"), sp("https://example.com/advisory/3"), sp("high")),
	})
	return &AdvisoryResult{Advisories: m, UnreachableRepos: []string{u.msg}}, nil
}

func map1(kv ...string) *phpjson.Map {
	m := phpjson.NewMap()
	for i := 0; i < len(kv); i += 2 {
		m.Set(kv[i], kv[i+1])
	}
	return m
}

func TestAuditWithIgnoreUnreachable(t *testing.T) {
	packages := pkgs(newPkg("vendor1/package1", "3.0.0.0", "3.0.0"))
	msg := `The "https://example.org/packages.json" file could not be downloaded: HTTP/1.1 404 Not Found`
	repoSet := &unreachableRepoSet{msg}
	auditor := NewAuditor(nil)

	_, err := auditor.Audit(&bufferIO{}, repoSet, createPolicyConfig(createPolicyOpts{}), packages, FormatPlain, false, nil)
	if _, ok := isTransportError(err); !ok || !strings.Contains(err.Error(), "HTTP/1.1 404 Not Found") {
		t.Fatalf("expected TransportException, got %v", err)
	}

	io := &bufferIO{}
	result, err := auditor.Audit(io, repoSet, createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, ignoreUnreachableAudit: true}), packages, FormatPlain, false, nil)
	if err != nil || result != StatusFailed {
		t.Fatalf("%v %d", err, result)
	}
	for _, s := range []string{"The following repositories were unreachable:", "HTTP/1.1 404 Not Found", "First repo advisory", "Third repo advisory", "CVE-2023-12345", "CVE-2023-67890"} {
		if !strings.Contains(io.Output(), s) {
			t.Fatalf("missing %q in %s", s, io.Output())
		}
	}

	io = &bufferIO{}
	result, err = auditor.Audit(io, repoSet, createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, ignoreUnreachableAudit: true}), packages, FormatJSON, false, nil)
	if err != nil || result != StatusFailed {
		t.Fatalf("%v %d", err, result)
	}
	v, err := phpjson.Decode([]byte(io.Output()))
	if err != nil {
		t.Fatal(err)
	}
	jm := phpjson.AsMap(v)
	ur, _ := jm.Get("unreachable-repositories")
	if phpjson.ArrayLen(ur) != 1 || !strings.Contains(toStr(coalesce(ur, "0", nil)), "HTTP/1.1 404 Not Found") {
		t.Fatal("unreachable-repositories")
	}
	list, _ := jm.Map("advisories").Get("vendor1/package1")
	if phpjson.ArrayLen(list) != 2 {
		t.Fatal("advisories count")
	}
	first, second := coalesce(list, "0", nil), coalesce(list, "1", nil)
	if toStr(coalesce(first, "cve", nil)) != "CVE-2023-12345" || toStr(coalesce(first, "title", nil)) != "First repo advisory" ||
		toStr(coalesce(second, "cve", nil)) != "CVE-2023-67890" || toStr(coalesce(second, "title", nil)) != "Third repo advisory" {
		t.Fatal("advisory json")
	}
}

// mockFilterListSource is the FilterListProviderSet mock.
type mockFilterListSource struct {
	filter *FilterByList
	calls  int
}

func (m *mockFilterListSource) GetMatchingFilterLists([]*pkg.Package, []string, bool) (*FilterListResult, error) {
	m.calls++
	return &FilterListResult{Filter: m.filter}, nil
}

func byList(listName string, entries ...*FilterListEntry) *FilterByList {
	m := NewOMap[[]*FilterListEntry]()
	m.Set(listName, entries)
	return m
}

func TestAuditWithFilter(t *testing.T) {
	matching := NewFilterListEntry("vendor/package", semver.NewConstraint(">=", "8.0.0.0"), "test-list", nil, sp("internal"), sp("ID-test-1"), nil)
	withDetails := NewFilterListEntry("vendor/package", semver.NewConstraint(">=", "8.0.0.0"), "test-list", sp("https://example.com/filtered"), sp("internal"), sp("ID-test-1"), nil)
	nonMatching := NewFilterListEntry("vendor/package", semver.NewConstraint(">=", "10.0.0.0"), "test-list", nil, sp("internal"), sp("ID-test-1"), nil)
	withSource := NewFilterListEntry("vendor/package", semver.NewConstraint(">=", "8.0.0.0"), "test-list", sp("https://example.com/filtered"), sp("internal"), sp("ID-test-1"), sp("aikido"))
	vp := func() []*pkg.Package { return pkgs(newPkg("vendor/package", "9.0.0", "9.0.0")) }
	cases := []struct {
		name     string
		packages []*pkg.Package
		filter   *FilterByList
		filtered string
		format   string
		expected int
		output   string
	}{
		{"AUDIT_IGNORE skips filter processing", vp(), byList("test-list", matching), AuditIgnore, FormatPlain, StatusOK, `No security vulnerability advisories found.`},
		{"AUDIT_FAIL with no matching entry returns STATUS_OK", vp(), byList("test-list", nonMatching), AuditFail, FormatPlain, StatusOK, `No security vulnerability advisories found.`},
		{"AUDIT_FAIL with matching entry returns STATUS_FAILED (plain)", vp(), byList("test-list", matching), AuditFail, FormatPlain, StatusFailed,
			"No security vulnerability advisories found.\nFound 1 package matching filters:\nvendor/package matched dependency policy \"test-list\". Reason: internal."},
		{"AUDIT_FAIL with matching entry shows url and reason (plain)", vp(), byList("test-list", withDetails), AuditFail, FormatPlain, StatusFailed,
			"No security vulnerability advisories found.\nFound 1 package matching filters:\nvendor/package matched dependency policy \"test-list\". Reason: internal. URL: https://example.com/filtered."},
		{"AUDIT_FAIL with matching entry shows source (plain)", vp(), byList("test-list", withSource), AuditFail, FormatPlain, StatusFailed,
			"No security vulnerability advisories found.\nFound 1 package matching filters:\nvendor/package matched dependency policy \"test-list\". Reason: internal. URL: https://example.com/filtered. Source: aikido."},
		{"AUDIT_REPORT with matching entry returns STATUS_OK", vp(), byList("test-list", matching), AuditReport, FormatPlain, StatusOK,
			"No security vulnerability advisories found.\n<warning>Found 1 package matching filters:</warning>\nvendor/package matched dependency policy \"test-list\". Reason: internal."},
		{"AUDIT_FAIL with matching entry shows summary line only (summary format)", vp(), byList("test-list", matching), AuditFail, FormatSummary, StatusFailed,
			"No security vulnerability advisories found.\nFound 1 package matching filters."},
		{"AUDIT_FAIL with multiple matching packages shows plural form", pkgs(newPkg("vendor/package", "9.0.0.0", "9.0.0"), newPkg("vendor/other", "1.0.0.0", "10.0.0")),
			byList("test-list", matching, NewFilterListEntry("vendor/other", semver.NewConstraint(">=", "1.0.0.0"), "test-list", nil, sp("internal"), sp("ID-TEST"), nil)), AuditFail, FormatPlain, StatusFailed,
			"No security vulnerability advisories found.\nFound 2 packages matching filters:\nvendor/package matched dependency policy \"test-list\". Reason: internal.\nvendor/other matched dependency policy \"test-list\". Reason: internal."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pc := createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, filteredAudit: tc.filtered})
			io := &bufferIO{}
			result, err := NewAuditor(nil).Audit(io, getRepoSet(t), pc, tc.packages, tc.format, true, &mockFilterListSource{filter: tc.filter})
			if err != nil {
				t.Fatal(err)
			}
			if result != tc.expected {
				t.Fatalf("result %d", result)
			}
			if got := io.trimmed(); got != tc.output {
				t.Fatalf("output mismatch\nexpected:\n%s\ngot:\n%s", tc.output, got)
			}
		})
	}
}

func decodeOutput(t *testing.T, io *bufferIO) *phpjson.Map {
	t.Helper()
	v, err := phpjson.Decode([]byte(io.Output()))
	if err != nil {
		t.Fatalf("invalid json %v: %s", err, io.Output())
	}
	return phpjson.AsMap(v)
}

func TestAuditWithFilterJson(t *testing.T) {
	entry := NewFilterListEntry("vendor/package", semver.NewConstraint(">=", "8.0.0.0"), "test-list", sp("https://example.com/filtered"), sp("Some reason"), sp("ID-test"), nil)
	src := &mockFilterListSource{filter: byList("test-list", entry)}
	io := &bufferIO{}
	result, err := NewAuditor(nil).Audit(io, getRepoSet(t), createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, filteredAudit: AuditFail}), pkgs(newPkg("vendor/package", "9.0.0", "9.0.0")), FormatJSON, true, src)
	if err != nil || result != StatusFailed || src.calls != 1 {
		t.Fatalf("%v %d %d", err, result, src.calls)
	}
	list, _ := decodeOutput(t, io).Map("filter").Get("vendor/package")
	if phpjson.ArrayLen(list) != 1 {
		t.Fatal("count")
	}
	e := phpjson.ToMap(coalesce(list, "0", nil))
	if e.Str("packageName") != "vendor/package" || e.Str("listName") != "test-list" || e.Str("url") != "https://example.com/filtered" || e.Str("reason") != "Some reason" {
		t.Fatal("entry")
	}
	if _, ok := coalesce(e, "constraint", nil).(string); !ok {
		t.Fatal("constraint must be a string")
	}
}

func TestAuditWithFilterReportsAliasedPackageOnce(t *testing.T) {
	entry := NewFilterListEntry("vendor/package", &semver.MatchAll{}, "test-list", sp("https://example.com/filtered"), sp("Some reason"), sp("ID-test"), nil)
	p := newPkg("vendor/package", "dev-main", "dev-main")
	alias := pkg.NewAlias(p, "9.1.9999999.9999999-dev", "9.1.x-dev")
	io := &bufferIO{}
	if _, err := NewAuditor(nil).Audit(io, getRepoSet(t), createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, filteredAudit: AuditFail}), pkgs(p, alias), FormatJSON, true, &mockFilterListSource{filter: byList("test-list", entry)}); err != nil {
		t.Fatal(err)
	}
	list, _ := decodeOutput(t, io).Map("filter").Get("vendor/package")
	if phpjson.ArrayLen(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", phpjson.ArrayLen(list))
	}
}

func TestAuditWithFilterAndVulnerabilities(t *testing.T) {
	entry := NewFilterListEntry("vendor1/package2", semver.NewConstraint(">=", "8.0.0.0"), "test-list", nil, sp("internal"), sp("ID-test"), nil)
	src := &mockFilterListSource{filter: byList("test-list", entry)}
	io := &bufferIO{}
	result, err := NewAuditor(nil).Audit(io, getRepoSet(t), createPolicyConfig(createPolicyOpts{abandoned: AuditIgnore, filteredAudit: AuditFail}),
		pkgs(newPkg("vendor1/package1", "8.2.1", "8.2.1"), newPkg("vendor1/package2", "9.0.0", "9.0.0")), FormatPlain, false, src)
	if err != nil || result != StatusFailed || src.calls != 1 {
		t.Fatal("result")
	}
	out := io.trimmed()
	for _, s := range []string{"Found 2 security vulnerability advisories affecting 1 package:", "Found 1 package matching filters:", `vendor1/package2 matched dependency policy "test-list". Reason: internal`} {
		if !strings.Contains(out, s) {
			t.Fatalf("missing %q", s)
		}
	}
}

func policyConfigWithCooldown(audit string) *PolicyConfig {
	return NewPolicyConfig(true, AdvisoriesPolicyConfigDisabled(), MalwarePolicyConfigDisabled(), AbandonedPolicyConfigDisabled(),
		NewCooldownPolicyConfig(true, audit, nil, i64(7*24*3600)), nil, IgnoreUnreachableDefault())
}

func TestCooldownAuditWiring(t *testing.T) {
	now := utc("2026-01-15 12:00:00")
	recent := newPkg("vendor/recent", "2.0.0.0", "2.0.0")
	recent.SetReleaseDate("2026-01-14 12:00:00")

	io := &bufferIO{}
	result, err := NewAuditor(NewCooldownAuditor(now)).Audit(io, getRepoSet(t), policyConfigWithCooldown(AuditFail), pkgs(recent), FormatJSON, false, nil)
	if err != nil || result != StatusFailed || !strings.Contains(io.Output(), `"cooldown"`) || !strings.Contains(io.Output(), "vendor/recent") {
		t.Fatalf("fail: %v %d %s", err, result, io.Output())
	}
	io = &bufferIO{}
	result, _ = NewAuditor(NewCooldownAuditor(now)).Audit(io, getRepoSet(t), policyConfigWithCooldown(AuditReport), pkgs(recent), FormatPlain, false, nil)
	if result != StatusOK || !strings.Contains(io.Output(), "still in the cooldown period") {
		t.Fatal("report")
	}
	io = &bufferIO{}
	result, _ = NewAuditor(NewCooldownAuditor(now)).Audit(io, getRepoSet(t), policyConfigWithCooldown(AuditIgnore), pkgs(recent), FormatPlain, false, nil)
	if result != StatusOK || strings.Contains(io.Output(), "still in the cooldown period") {
		t.Fatal("ignore")
	}
}

func advisoryData(name, id, affected string) string {
	return `{"advisoryId": "` + id + `", "packageName": "` + name + `", "remoteId": "` + id + `", "title": "advisory", "link": null, "cve": null, "affectedVersions": "` + affected + `", "source": "test", "reportedAt": "2024-04-31 12:37:47", "composerRepository": "Package Repository", "severity": "high", "sources": [{"name": "test", "remoteId": "` + id + `"}]}`
}

func TestAuditSkipsDevBranchesButMatchesTheirBranchAliases(t *testing.T) {
	resolver := &ProviderResolver{}
	r := packageRepository(t, resolver, `{"package": [], "security-advisories": {"acme/dev-only": [`+advisoryData("acme/dev-only", "ID-DEV", "*")+`], "acme/aliased": [`+advisoryData("acme/aliased", "ID-ALIAS", "<3.5")+`]}}`)
	repoSet := &RepositoryAdvisorySource{Repositories: []repo.Repository{r}, Resolver: resolver}

	devOnly := newPkg("acme/dev-only", "dev-main", "dev-main")
	aliased := loadPkg(t, "dev-main", "dev-main", `{"name": "acme/aliased", "version": "dev-main", "abandoned": true}`)
	branchAlias := pkg.NewAlias(aliased, "3.3.9999999.9999999-dev", "3.3.x-dev")

	io := &bufferIO{}
	result, err := NewAuditor(nil).Audit(io, repoSet, createPolicyConfig(createPolicyOpts{abandoned: AuditFail}), pkgs(devOnly, aliased, branchAlias), FormatJSON, false, nil)
	if err != nil || result != StatusFailed {
		t.Fatalf("%v %d", err, result)
	}
	out := decodeOutput(t, io)
	if adv := out.Map("advisories"); strings.Join(adv.Keys(), ",") != "acme/aliased" {
		t.Fatalf("advisories keys %v", adv.Keys())
	}
	list, _ := out.Map("advisories").Get("acme/aliased")
	if toStr(coalesce(coalesce(list, "0", nil), "advisoryId", nil)) != "ID-ALIAS" {
		t.Fatal("advisory id")
	}
	if ab, _ := phpjson.Encode(out.Map("abandoned"), phpjson.UnescapedSlashes); ab != `{"acme/aliased":null}` {
		t.Fatal(ab)
	}

	io = &bufferIO{}
	if _, err := NewAuditor(nil).Audit(io, repoSet, createPolicyConfig(createPolicyOpts{abandoned: AuditFail}), pkgs(devOnly, aliased, branchAlias), FormatPlain, false, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(io.Output(), "Found 1 abandoned package:") {
		t.Fatal(io.Output())
	}
}

func TestAuditOutputsPartialAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	r := packageRepository(t, resolver, `{"package": [], "security-advisories": {"acme/lib": [{"advisoryId": "ACME-2026-001", "affectedVersions": "<1.0.1"}]}}`)
	io := &bufferIO{}
	result, err := NewAuditor(nil).Audit(io, &RepositoryAdvisorySource{Repositories: []repo.Repository{r}, Resolver: resolver}, createPolicyConfig(createPolicyOpts{}), pkgs(newPkg("acme/lib", "1.0.0.0", "1.0.0")), FormatPlain, false, nil)
	if err != nil || result != StatusFailed {
		t.Fatalf("%v %d", err, result)
	}
	want := "Found 1 security vulnerability advisory affecting 1 package:\nPackage: acme/lib\nAdvisory ID: ACME-2026-001\nAffected versions: <1.0.1"
	if io.trimmed() != want {
		t.Fatalf("got:\n%s", io.trimmed())
	}
}

func TestNeedsCompleteAdvisoryLoad(t *testing.T) {
	full := func(id string) *Advisory {
		return NewSecurityAdvisory("foo/bar", id, semver.NewConstraint("=", "1.0.0.0"), "test", []phpjson.Value{map1("name", "foo", "remoteId", "remoteID")}, time.Now(), nil, nil, nil)
	}
	partial := NewPartialSecurityAdvisory("foo/bar", "1234", semver.NewConstraint("=", "1.0.0.0"))
	advs := func(pairs ...any) *AdvisoryMap {
		m := NewOMap[[]*Advisory]()
		for i := 0; i < len(pairs); i += 2 {
			m.Set(pairs[i].(string), pairs[i+1].([]*Advisory))
		}
		return m
	}
	ignore := func(keys ...string) *OMap[*string] {
		m := NewOMap[*string]()
		for _, k := range keys {
			m.Set(k, nil)
		}
		return m
	}
	cases := []struct {
		name     string
		adv      *AdvisoryMap
		ignore   *OMap[*string]
		expected bool
	}{
		{"no filter or advisories", advs(), ignore(), false},
		{"packagist filters are IDs so work fine with partial advisories", advs(), ignore("PKSA-foo-bar"), false},
		{"packagist filters are IDs so work fine with partial advisories/2", advs("vendor1/package1", []*Advisory{full("123"), partial}), ignore("PKSA-foo-bar"), false},
		{"no advisories no need to load any further", advs(), ignore("CVE-2025-1234"), false},
		{"no advisories no need to load any further/2", advs("vendor1/package1", []*Advisory{}), ignore("CVE-2025-1234"), false},
		{"CVE filter or other non-packagist ones might need to fully load for safety if partial advisories are present", advs("vendor1/package1", []*Advisory{full("123"), partial}), ignore("CVE-2025-1234"), true},
		{"filter does not trigger load if all advisories are fully loaded", advs("vendor1/package1", []*Advisory{full("123")}, "vendor1/package2", []*Advisory{full("1234")}), ignore("CVE-2025-1234"), false},
	}
	for _, tc := range cases {
		if got := NewAuditor(nil).NeedsCompleteAdvisoryLoad(tc.adv, tc.ignore); got != tc.expected {
			t.Fatalf("%s: got %v", tc.name, got)
		}
	}
}

// ------------------------------------------------------- CooldownAuditorTest

func cooldownPkg(t *testing.T, name, version string, releaseDate, publishedDate string) *pkg.Package {
	data := map1("name", name, "version", version)
	if releaseDate != "" {
		data.Set("time", releaseDate)
	}
	if publishedDate != "" {
		data.Set("published-time", publishedDate)
	}
	enc, _ := phpjson.Encode(data, 0)
	return loadPkg(t, version+".0", version, enc)
}

func cooldownConfig(audit string, period *int64, ignore *IgnoreMap) *CooldownPolicyConfig {
	return NewCooldownPolicyConfig(true, audit, ignore, period)
}

const cooldownNow = "2026-01-15 12:00:00"

func TestCooldownAuditorCollectsOnlyVersionsWithinCooldown(t *testing.T) {
	old := cooldownPkg(t, "vendor/old", "1.0.0", "2026-01-01 12:00:00", "")
	recent := cooldownPkg(t, "vendor/recent", "2.0.0", "2026-01-14 12:00:00", "")
	res := NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(old, recent), cooldownConfig(AuditReport, i64(604800), nil))
	e, ok := res.Get("vendor/recent")
	if res.Has("vendor/old") || !ok || e.PrettyVersion != "2.0.0" || e.Source != "time" || e.AvailableIn == "" {
		t.Fatalf("%+v", e)
	}
}

func TestCooldownAuditorPrefersPublishedDateOverReleaseDate(t *testing.T) {
	p := cooldownPkg(t, "vendor/pkg", "1.0.0", "2025-01-01 12:00:00", "2026-01-14 12:00:00")
	e, ok := NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(p), cooldownConfig(AuditReport, i64(604800), nil)).Get("vendor/pkg")
	if !ok || e.Source != "published-time" || e.ReleaseDate != formatAtom(utc("2026-01-14 12:00:00")) {
		t.Fatalf("%+v", e)
	}
}

func TestCooldownAuditorSkipsDevPlatformRootAndDatelessPackages(t *testing.T) {
	dev := cooldownPkg(t, "vendor/dev", "dev-main", cooldownNow, "")
	platform := cooldownPkg(t, "php", "8.3.0", cooldownNow, "")
	dateless := cooldownPkg(t, "vendor/nodate", "1.0.0", "", "")
	root := newPkg("vendor/root", "1.0.0.0", "1.0.0")
	root.IsRoot = true
	root.SetReleaseDate(cooldownNow)
	if res := NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(dev, platform, dateless, root), cooldownConfig(AuditReport, i64(604800), nil)); res.Len() != 0 {
		t.Fatal(res.Keys())
	}
}

func TestCooldownAuditorRespectsIgnoreRules(t *testing.T) {
	recent := cooldownPkg(t, "vendor/recent", "2.0.0", "2026-01-14 12:00:00", "")
	ignored := cooldownPkg(t, "vendor/ignored", "2.0.0", "2026-01-14 12:00:00", "")
	ignore, _ := ParseIgnoreMap(j(t, `{"vendor/ignored": "pinned by us"}`))
	res := NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(recent, ignored), cooldownConfig(AuditReport, i64(604800), ignore))
	if !res.Has("vendor/recent") || res.Has("vendor/ignored") {
		t.Fatal(res.Keys())
	}

	p := cooldownPkg(t, "vendor/pkg", "2.0.0", "2026-01-14 12:00:00", "")
	blockOnly := ignoreMap(rule("vendor/pkg", nil, nil, true, false))
	if !NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(p), cooldownConfig(AuditReport, i64(604800), blockOnly)).Has("vendor/pkg") {
		t.Fatal("block-only ignore must not apply to audit")
	}
}

func TestCooldownAuditorResolvesAliasesToAliasedPackage(t *testing.T) {
	dev := newPkg("vendor/dev", "dev-main", "dev-main")
	dev.SetReleaseDate("2026-01-14 12:00:00")
	stable := cooldownPkg(t, "vendor/stable", "1.2.0", "2026-01-14 12:00:00", "")
	res := NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(pkg.NewAlias(dev, "1.0.0.0", "1.0.0"), pkg.NewAlias(stable, "1.0.0.0", "1.0.0")), cooldownConfig(AuditReport, i64(604800), nil))
	e, ok := res.Get("vendor/stable")
	if res.Has("vendor/dev") || !ok || e.PrettyVersion != "1.2.0" {
		t.Fatal(res.Keys())
	}
}

func TestCooldownAuditorReturnsEmptyWhenNoCooldownConfigured(t *testing.T) {
	recent := cooldownPkg(t, "vendor/recent", "2.0.0", "2026-01-14 12:00:00", "")
	if NewCooldownAuditor(utc(cooldownNow)).Collect(pkgs(recent), cooldownConfig(AuditReport, nil, nil)).Len() != 0 {
		t.Fatal("expected empty")
	}
}
