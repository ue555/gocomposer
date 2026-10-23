// Ports of Composer f83588392df43d4fd027d84d26de459946ef4e80
// tests/Composer/Test/FilterList/** (expectations unchanged).
// HttpDownloaderMock is replaced by a fake HTTPDownloader asserting the same
// url/options and returning the same body; PHP anonymous repository classes
// are Go fakes embedding repo.Array.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"reflect"
	"testing"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// ---------------------------------------- ComposerRepositoryFilterInformationTest

func TestComposerRepositoryFilterInformation(t *testing.T) {
	info := ComposerRepositoryFilterInformationFromData(j(t, `{"metadata": true, "lists": {"company-policy": {"enabled": true}, "aikido": {"enabled": true}}}`), nil)
	if !info.Metadata || !reflect.DeepEqual(info.Lists, []string{"company-policy", "aikido"}) {
		t.Fatalf("passes through: %+v", info)
	}
	for _, data := range []string{
		`{"lists": {"company-policy": {"enabled": true}, "aikido": {"enabled": false}}}`,
		`{"lists": {"company-policy": {"enabled": true}, "malware": []}}`,
		`{"lists": {"company-policy": {"enabled": true}, "malware": true}}`,
		`{"lists": {"advisories": {"enabled": true}, "company-policy": {"enabled": true}, "abandoned": {"enabled": true}, "cooldown": {"enabled": true}}}`,
		`{"lists": {"ignore-foo": {"enabled": true}, "ignoremalware": {"enabled": true}, "company-policy": {"enabled": true}}}`,
	} {
		if got := ComposerRepositoryFilterInformationFromData(j(t, data), nil).Lists; !reflect.DeepEqual(got, []string{"company-policy"}) {
			t.Fatalf("%s => %v", data, got)
		}
	}
	if ComposerRepositoryFilterInformationFromData(j(t, `{"lists": {"malware": {"enabled": true}}}`), nil).SummaryURL != nil {
		t.Fatal("summary default")
	}
	prefix := func(u string) string { return "https://example.org" + u }
	if s := ComposerRepositoryFilterInformationFromData(j(t, `{"lists": {"malware": {"enabled": true}}, "summary-url": "/p2/filter-summary.json"}`), prefix).SummaryURL; s == nil || *s != "https://example.org/p2/filter-summary.json" {
		t.Fatal("summary canonicalized")
	}
	if ComposerRepositoryFilterInformationFromData(j(t, `{"lists": {"malware": {"enabled": true}}, "summary-url": ["oops"]}`), func(u string) string { return u }).SummaryURL != nil {
		t.Fatal("non-string summary")
	}
	if ComposerRepositoryFilterInformationFromData(j(t, `{"lists": {"malware": {"enabled": true}}}`), nil).APIURL != nil {
		t.Fatal("api default")
	}
	if s := ComposerRepositoryFilterInformationFromData(j(t, `{"lists": {"malware": {"enabled": true}}, "api-url": "/api/filter"}`), prefix).APIURL; s == nil || *s != "https://example.org/api/filter" {
		t.Fatal("api canonicalized")
	}
}

// ----------------------------------------------------- FilterListApiClientTest

type httpExpectation struct {
	url     string
	options string // JSON of the expected options, "" = not checked
	body    string
}

// fakeHTTP is HttpDownloaderMock with strict expectations.
type fakeHTTP struct {
	t        *testing.T
	expected []httpExpectation
	calls    int
}

func (f *fakeHTTP) Request(u string, options *phpjson.Map) (*httpdl.Response, error) {
	f.t.Helper()
	if f.calls >= len(f.expected) {
		f.t.Fatalf("unexpected request to %s", u)
	}
	e := f.expected[f.calls]
	f.calls++
	if e.url != u {
		f.t.Fatalf("expected url %s, got %s", e.url, u)
	}
	if e.options != "" {
		got, _ := phpjson.Encode(options, 0)
		if got != e.options {
			f.t.Fatalf("options mismatch\nexpected %s\ngot      %s", e.options, got)
		}
	}
	return &httpdl.Response{URL: u, StatusCode: 200, Body: []byte(e.body)}, nil
}

func (f *fakeHTTP) assertComplete() {
	f.t.Helper()
	if f.calls != len(f.expected) {
		f.t.Fatalf("expected %d requests, got %d", len(f.expected), f.calls)
	}
}

func pcmOf(pairs ...any) *ConstraintMap {
	m := NewOMap[semver.Constraint]()
	for i := 0; i < len(pairs); i += 2 {
		m.Set(pairs[i].(string), pairs[i+1].(semver.Constraint))
	}
	return m
}

func TestPostPurlsSendsPackagesAndListsAsBody(t *testing.T) {
	body := `{"packages":["pkg:\/\/composer\/vendor\/foo","pkg:\/\/composer\/vendor\/bar"],"lists":["malware","typosquatting"]}`
	opts, _ := phpjson.Encode(j(t, `{"http": {"method": "POST", "header": ["Content-type: application/json"], "timeout": 10, "content": `+jsonString(body)+`}}`), 0)
	h := &fakeHTTP{t: t, expected: []httpExpectation{{"https://example.org/api/filter", opts, `{"filter":[]}`}}}
	client := &FilterListApiClient{HTTP: h}
	resp, err := client.PostPurls("https://example.org/api/filter", pcmOf("vendor/foo", semver.NewConstraint("=", "1.0.0.0"), "vendor/bar", semver.NewConstraint("=", "2.0.0.0")), []string{"malware", "typosquatting"})
	if err != nil {
		t.Fatal(err)
	}
	decoded, _ := DecodeJSONResponse(resp, "https://example.org/api/filter")
	if enc, _ := phpjson.Encode(decoded, 0); enc != `{"filter":[]}` {
		t.Fatal(enc)
	}
	h.assertComplete()
}

func jsonString(s string) string {
	enc, _ := phpjson.Encode(s, phpjson.UnescapedSlashes)
	return enc
}

func TestPostPurlsSendsRepositoryTransportOptions(t *testing.T) {
	body := `{"packages":["pkg:\/\/composer\/vendor\/foo"],"lists":["malware"]}`
	opts, _ := phpjson.Encode(j(t, `{"ssl": {"verify_peer": false}, "http": {"header": ["X-Cops: S07E12", "Content-type: application/json"], "method": "POST", "timeout": 10, "content": `+jsonString(body)+`}}`), 0)
	// (object) ['X-Cops: S07E12'] has no Go equivalent; the other shapes are covered
	variations := []string{`"X-Cops: S07E12"`, `["X-Cops: S07E12"]`, `{"0": "X-Cops: S07E12"}`}
	h := &fakeHTTP{t: t}
	for range variations {
		h.expected = append(h.expected, httpExpectation{"https://example.org/api/filter", opts, `{"filter":[]}`})
	}
	for _, v := range variations {
		client := &FilterListApiClient{HTTP: h, Options: phpjson.AsMap(j(t, `{"ssl": {"verify_peer": false}, "http": {"header": `+v+`}}`))}
		resp, err := client.PostPurls("https://example.org/api/filter", pcmOf("vendor/foo", semver.NewConstraint("=", "1.0.0.0")), []string{"malware"})
		if err != nil {
			t.Fatal(err)
		}
		if string(resp.Body) != `{"filter":[]}` {
			t.Fatal("body")
		}
	}
	h.assertComplete()
}

// ------------------------------------------------------- FilterListAuditorTest

func mustEntry(t *testing.T, listName, data string) *FilterListEntry {
	t.Helper()
	e, err := CreateFilterListEntry(listName, j(t, data))
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func filterMap(entries ...*FilterListEntry) *FilterListMap {
	m := NewOMap[*FilterByList]()
	for _, e := range entries {
		byList, ok := m.Get(e.PackageName)
		if !ok {
			byList = NewOMap[[]*FilterListEntry]()
			m.Set(e.PackageName, byList)
		}
		appendTo(byList, e.ListName, e)
	}
	return m
}

func TestGetMatchingEntriesUnfilteredPackages(t *testing.T) {
	cases := map[string]struct {
		ignore string
		count  int
	}{
		"acme/other fully ignore":   {`["acme/other"]`, 1},
		"acme/package fully ignore": {`["acme/package"]`, 0},
		"acme/package fully ignore but not on-block with block operation": {`{"acme/package": {"on-block": false}}`, 1},
		"acme/package fully ignore but not on-audit with block operation": {`{"acme/package": {"on-audit": false}}`, 0},
		"acme/* fully ignore":                                {`["acme/*"]`, 0},
		"acme/package 1.0 ignore":                            {`{"acme/package": {"constraint": "1.0"}}`, 0},
		"acme/* 1.0 ignore":                                  {`{"acme/*": {"constraint": "1.0"}}`, 0},
		"acme/package 1.1 ignore":                            {`{"acme/package": {"constraint": "1.1"}}`, 1},
		"acme/* 1.1 ignore":                                  {`{"acme/*": {"constraint": "1.1"}}`, 1},
		"multiple acme/package entries, first rule matches":  {`{"acme/package": [{"constraint": "1.0"}, {"constraint": "1.1"}]}`, 0},
		"multiple acme/package entries, second rule matches": {`{"acme/package": [{"constraint": "1.1"}, {"constraint": "1.0"}]}`, 0},
		"multiple acme/package entries, no rule matches":     {`{"acme/package": [{"constraint": "1.1"}, {"constraint": "1.2"}]}`, 1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			p := newPkg("acme/package", "1.0.0.0", "1.0")
			m := filterMap(mustEntry(t, "list", `{"package": "acme/package", "constraint": "*"}`), mustEntry(t, "list", `{"package": "acme/other", "constraint": "*"}`))
			pc := mustPolicy(t, newConfig(t, `{"policy": {"list": {"ignore": `+tc.ignore+`}}}`))
			if got := len((FilterListAuditor{}).GetMatchingBlockEntries(p, m, pc, BlockScopeUpdate)); got != tc.count {
				t.Fatalf("got %d", got)
			}
		})
	}
}

func TestGetMatchingEntriesIgnoreSource(t *testing.T) {
	for _, op := range []string{OperationBlock, OperationAudit} {
		for _, tc := range []struct {
			source, ignoreSource string
			count                int
		}{
			{`"untrusted"`, `["untrusted"]`, 0},
			{`"trusted"`, `["untrusted"]`, 1},
			{`null`, `["untrusted"]`, 1},
			{`"untrusted"`, `[]`, 1},
		} {
			clearEnv(t)
			p := newPkg("acme/package", "1.0.0.0", "1.0")
			m := filterMap(mustEntry(t, "malware", `{"package": "acme/package", "constraint": "*", "source": `+tc.source+`}`))
			pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": {"ignore-source": `+tc.ignoreSource+`}}}`))
			var got []*FilterListEntry
			if op == OperationAudit {
				got = FilterListAuditor{}.GetMatchingAuditEntries(p, m, pc)
			} else {
				got = FilterListAuditor{}.GetMatchingBlockEntries(p, m, pc, BlockScopeUpdate)
			}
			if len(got) != tc.count {
				t.Fatalf("%s %+v: got %d", op, tc, len(got))
			}
		}
	}
}

func TestGetMatchingEntriesKeepsNonIgnoredSourcesAndDropsIgnored(t *testing.T) {
	clearEnv(t)
	p := newPkg("acme/package", "1.0.0.0", "1.0")
	m := filterMap(
		mustEntry(t, "malware", `{"package": "acme/package", "constraint": "*", "source": "untrusted"}`),
		mustEntry(t, "malware", `{"package": "acme/package", "constraint": "*", "source": "trusted"}`),
	)
	pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": {"ignore-source": ["untrusted"]}}}`))
	got := FilterListAuditor{}.GetMatchingBlockEntries(p, m, pc, BlockScopeUpdate)
	if len(got) != 1 || *got[0].Source != "trusted" {
		t.Fatal("mismatch")
	}
}

func TestGetMatchingEntriesIgnoresUnconfiguredLists(t *testing.T) {
	clearEnv(t)
	p := newPkg("acme/package", "1.0.0.0", "1.0")
	m := filterMap(mustEntry(t, "unconfigured", `{"package": "acme/package", "constraint": "*"}`))
	pc := mustPolicy(t, newConfig(t, `{"policy": {"list": {"ignore": ["acme/package"]}}}`))
	if got := (FilterListAuditor{}).GetMatchingBlockEntries(p, m, pc, BlockScopeUpdate); len(got) != 0 {
		t.Fatal("expected none")
	}
}

func TestGetMatchingEntriesDropsUnconfiguredListEntriesForNonIgnoredPackage(t *testing.T) {
	clearEnv(t)
	p := newPkg("acme/package", "1.0.0.0", "1.0")
	m := filterMap(mustEntry(t, "configured", `{"package": "acme/package", "constraint": "*"}`), mustEntry(t, "unconfigured", `{"package": "acme/package", "constraint": "*"}`))
	pc := mustPolicy(t, newConfig(t, `{"policy": {"configured": true}}`))
	got := FilterListAuditor{}.GetMatchingBlockEntries(p, m, pc, BlockScopeUpdate)
	if len(got) != 1 || got[0].ListName != "configured" {
		t.Fatal("mismatch")
	}
}

// -------------------------------------------------- FilterListEntryBuilderTest

func TestBuildFiltersByConstraintMatch(t *testing.T) {
	res, err := BuildFilterListEntries(j(t, `{"malware": [
		{"package": "vendor/foo", "constraint": "^1.0", "url": "https://example.org/foo", "reason": "bad", "id": "A"},
		{"package": "vendor/foo", "constraint": "^9.0", "url": null, "reason": "ignored", "id": "B"},
		{"package": "vendor/unknown", "constraint": "*", "url": null, "reason": "ignored", "id": "C"}]}`), pcmOf("vendor/foo", semver.NewConstraint("=", "1.2.3.0")), "")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := res.Get("malware")
	if !ok || len(l) != 1 || l[0].PackageName != "vendor/foo" || *l[0].ID != "A" {
		t.Fatal("mismatch")
	}
}

func TestBuildFillsInDefaultPackageWhenMissing(t *testing.T) {
	res, _ := BuildFilterListEntries(j(t, `{"malware": [{"constraint": "^1.0", "url": "https://example.org/foo", "reason": "bad", "id": "PKFE-001"}]}`), pcmOf("vendor/foo", semver.NewConstraint("=", "1.2.3.0")), "vendor/foo")
	if l, _ := res.Get("malware"); len(l) != 1 || l[0].PackageName != "vendor/foo" {
		t.Fatal("mismatch")
	}
}

func TestBuildPrefersExplicitPackageOverDefault(t *testing.T) {
	res, _ := BuildFilterListEntries(j(t, `{"malware": [{"package": "vendor/foo", "constraint": "*", "id": "X"}]}`), pcmOf("vendor/foo", semver.NewConstraint("=", "1.0.0.0")), "vendor/bar")
	if l, _ := res.Get("malware"); len(l) != 1 || l[0].PackageName != "vendor/foo" {
		t.Fatal("mismatch")
	}
}

func TestBuildIgnoresMalformedListShapes(t *testing.T) {
	res, _ := BuildFilterListEntries(j(t, `{"malware": "not-an-array", "5": [{"package": "vendor/foo", "constraint": "*"}], "typosquatting": ["not-an-entry", {"package": "vendor/foo", "constraint": "*"}]}`), pcmOf("vendor/foo", semver.NewConstraint("=", "1.0.0.0")), "")
	if !reflect.DeepEqual(res.Keys(), []string{"typosquatting"}) {
		t.Fatal(res.Keys())
	}
	if l, _ := res.Get("typosquatting"); len(l) != 1 {
		t.Fatal("count")
	}
}

func TestBuildReturnsEmptyForEmptyInput(t *testing.T) {
	res, _ := BuildFilterListEntries(j(t, `[]`), pcmOf("vendor/foo", semver.NewConstraint("=", "1.0.0.0")), "")
	if res.Len() != 0 {
		t.Fatal("expected empty")
	}
}

// --------------------------------------------------- FilterListProviderSetTest

// packageRepository is "new PackageRepository($config)" registered with the resolver.
func packageRepository(t *testing.T, resolver *ProviderResolver, config string) *repo.Array {
	t.Helper()
	cfg := phpjson.AsMap(j(t, config))
	a, err := repo.NewPackageRepository(cfg)
	if err != nil {
		t.Fatal(err)
	}
	resolver.RegisterPackageRepository(a, cfg)
	return a
}

func TestGetMatchingFilterListsOnlyReturnsConfiguredLists(t *testing.T) {
	resolver := &ProviderResolver{}
	r := packageRepository(t, resolver, `{"package": [], "filter": {
		"malware": [{"package": "acme/package", "constraint": "1.0", "reason": "malware"}],
		"typosquatting": [{"package": "acme/package", "constraint": "1.0", "reason": "typosquatting"}]}}`)
	set, err := NewFilterListProviderSet([]repo.Repository{r}, nil, resolver)
	if err != nil {
		t.Fatal(err)
	}
	res, err := set.GetMatchingFilterLists(pkgs(newPkg("acme/package", "1.0.0.0", "1.0")), []string{"malware"}, false)
	if err != nil {
		t.Fatal(err)
	}
	l, ok := res.Filter.Get("malware")
	if !ok || res.Filter.Has("typosquatting") || len(l) != 1 || *l[0].Reason != "malware" {
		t.Fatal("mismatch")
	}
}

// unreachableFilterRepo is the anonymous ArrayRepository whose hasFilter() throws.
type unreachableFilterRepo struct {
	*repo.Array
	message string
}

func (u *unreachableFilterRepo) HasFilter() (bool, error) {
	return false, &httpdl.TransportError{Msg: u.message}
}
func (u *unreachableFilterRepo) GetFilter(*ConstraintMap, []string) (*FilterByList, error) {
	return NewOMap[[]*FilterListEntry](), nil
}
func (u *unreachableFilterRepo) GetFilterLists() ([]string, error) { return nil, nil }

func TestUnreachableRepositoryDuringConstruction(t *testing.T) {
	r := &unreachableFilterRepo{Array: repo.NewArray(), message: "repo.example.com could not be reached"}
	set, err := NewFilterListProviderSet([]repo.Repository{r}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := set.GetMatchingFilterLists(pkgs(newPkg("acme/package", "1.0.0.0", "1.0")), []string{"malware"}, true)
	if err != nil || res.Filter.Len() != 0 || !reflect.DeepEqual(res.UnreachableRepos, []string{"repo.example.com could not be reached"}) {
		t.Fatalf("ignored: %v %+v", err, res)
	}
	_, err = set.GetMatchingFilterLists(pkgs(newPkg("acme/package", "1.0.0.0", "1.0")), []string{"malware"}, false)
	if _, ok := isTransportError(err); !ok || err.Error() != "repo.example.com could not be reached" {
		t.Fatalf("expected TransportException, got %v", err)
	}
}
