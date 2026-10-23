// Ports of Composer f83588392df43d4fd027d84d26de459946ef4e80
// tests/Composer/Test/DependencyResolver/SecurityAdvisoryPoolFilterTest.php,
// FilterListPoolFilterTest.php, CooldownPoolFilterTest.php and
// PolicyRemovalReasonTest.php (expectations unchanged). Pool/Request are the
// policy.Pool / policy.Request abstractions; anonymous PHP repository
// subclasses are Go fakes.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// testRequest is Composer\DependencyResolver\Request.
type testRequest struct {
	locked     []*pkg.Package
	lockedRepo []*pkg.Package
}

func (r *testRequest) IsLockedPackage(p *pkg.Package) bool {
	for _, l := range r.locked {
		if l == p {
			return true
		}
	}
	return false
}
func (r *testRequest) LockedRepositoryPackages() []*pkg.Package { return r.lockedRepo }

func newPool(ps ...*pkg.Package) *Pool {
	return &Pool{Packages: ps, PolicyRemoved: NewPolicyRemovedVersions()}
}

func removalType(p *Pool, name, version string) string {
	if r := p.PolicyRemoved.Reason(name, version); r != nil {
		return r.Type()
	}
	return ""
}

func removals(p *Pool, typ string) []string {
	var out []string
	p.PolicyRemoved.All(func(name, version string, r *PolicyRemovalReason) {
		if r.Type() == typ {
			out = append(out, name+" "+version)
		}
	})
	return out
}

// ----------------------------------------------- SecurityAdvisoryPoolFilterTest

func advPolicyConfig(advisoriesBlock, abandonedBlock bool, advIgnore, abIgnore *IgnoreMap) *PolicyConfig {
	return NewPolicyConfig(true, NewAdvisoriesPolicyConfig(advisoriesBlock, AuditFail, advIgnore, nil, nil), MalwarePolicyConfigDisabled(),
		NewAbandonedPolicyConfig(abandonedBlock, AuditFail, abIgnore), CooldownPolicyConfigDisabled(), nil, IgnoreUnreachableDefault())
}

var uniq int

func genAdvisory(name, cve, affected string) (string, string) {
	uniq++
	id := fmt.Sprintf("PKSA-%013x", uniq)
	c := "null"
	if cve != "" {
		c = `"` + cve + `"`
	}
	return id, `{"advisoryId": "` + id + `", "packageName": "` + name + `", "remoteId": "test", "title": "Security Advisory", "link": null, "cve": ` + c + `, "affectedVersions": "` + affected + `", "source": "Tests", "reportedAt": "2024-04-31 12:37:47", "composerRepository": "Package Repository", "severity": "high", "sources": [{"name": "Security Advisory", "remoteId": "test"}]}`
}

func advRepo(t *testing.T, resolver *ProviderResolver, advisories string) *repo.Array {
	return packageRepository(t, resolver, `{"package": [], "security-advisories": `+advisories+`}`)
}

func TestFilterPackagesByAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	id1, a1 := genAdvisory("acme/package", "CVE-1999-1000", ">=1.0.0,<1.1.0")
	id2, a2 := genAdvisory("acme/package", "CVE-1999-1001", ">=1.0.0,<1.1.0")
	r := advRepo(t, resolver, `{"acme/package": [`+a1+`, `+a2+`]}`)
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	exp1, exp2 := newPkg("acme/package", "2.0.0.0", "2.0"), newPkg("acme/other", "1.0.0.0", "1.0")
	filtered, err := f.Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0"), exp1, exp2), []repo.Repository{r}, &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(exp1, exp2), filtered.Packages)
	if removalType(filtered, "acme/package", "1.0.0.0") != RemovalAdvisories || len(removals(filtered, RemovalAbandoned)) != 0 {
		t.Fatal("removals")
	}
	var ids []string
	for _, a := range filtered.PolicyRemoved.Reason("acme/package", "1.0.0.0").Advisories() {
		ids = append(ids, a.AdvisoryID)
	}
	if !reflect.DeepEqual(ids, []string{id1, id2}) {
		t.Fatal(ids)
	}
}

func TestDontFilterPackagesByIgnoredAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", ">=1.0.0,<1.1.0")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	ids := NewOMap[*IgnoreIdRule]()
	ids.Set("CVE-2024-1234", NewIgnoreIdRule("CVE-2024-1234", nil, true, true))
	pc := NewPolicyConfig(true, NewAdvisoriesPolicyConfig(true, AuditFail, nil, ids, nil), MalwarePolicyConfigDisabled(), NewAbandonedPolicyConfig(true, AuditFail, nil), CooldownPolicyConfigDisabled(), nil, IgnoreUnreachableDefault())
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), pc, nil, resolver)
	e1, e2 := newPkg("acme/package", "1.0.0.0", "1.0"), newPkg("acme/package", "1.1.0.0", "1.1")
	filtered, err := f.Filter(newPool(e1, e2), []repo.Repository{r}, &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(e1, e2), filtered.Packages)
	if len(removals(filtered, RemovalAbandoned))+len(removals(filtered, RemovalAdvisories)) != 0 {
		t.Fatal("removals")
	}
}

func TestDontFilterPackagesWithBlockInsecureDisabled(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", ">=1.0.0,<1.1.0")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(false, true, nil, nil), nil, resolver)
	e1, e2 := newPkg("acme/package", "1.0.0.0", "1.0"), newPkg("acme/package", "1.1.0.0", "1.1")
	filtered, _ := f.Filter(newPool(e1, e2), []repo.Repository{r}, &testRequest{})
	assertSamePackages(t, pkgs(e1, e2), filtered.Packages)
	if len(removals(filtered, RemovalAbandoned))+len(removals(filtered, RemovalAdvisories)) != 0 {
		t.Fatal("removals")
	}
}

func TestDontFilterPackagesWithAbandonedPackage(t *testing.T) {
	ignoreName := "acme/ignore-abandoned"
	pc := advPolicyConfig(true, true, nil, ignoreMap(rule(ignoreName, nil, nil, true, true)))
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), pc, nil, nil)
	abandoned := loadPkg(t, "1.0.0.0", "1.0", `{"name": "acme/package", "version": "1.0", "abandoned": true}`)
	ignoreAbandoned := loadPkg(t, "1.0.0.0", "1.0", `{"name": "`+ignoreName+`", "version": "1.0", "abandoned": true}`)
	expected := newPkg("acme/other", "1.1.0.0", "1.1")
	filtered, err := f.Filter(newPool(expected, abandoned, ignoreAbandoned), nil, &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(expected, ignoreAbandoned), filtered.Packages)
	if len(removals(filtered, RemovalAbandoned)) != 1 || len(removals(filtered, RemovalAdvisories)) != 0 {
		t.Fatal("removals")
	}
}

// failingAdvisoryRepo is the anonymous PackageRepository whose getSecurityAdvisories throws.
type failingAdvisoryRepo struct {
	*repo.Array
	*PackageRepositoryProvider
	err error
}

func (f *failingAdvisoryRepo) GetSecurityAdvisories(*ConstraintMap, bool) (*AdvisoryMap, error) {
	return nil, f.err
}
func (f *failingAdvisoryRepo) RepoName() string { return "unreachable advisory repo" }

func TestWarnsWhenUnreachableAdvisoryRepositoriesAreIgnored(t *testing.T) {
	_, a := genAdvisory("acme/package", "CVE-2024-9999", ">=1.0.0,<2.0.0")
	r := &failingAdvisoryRepo{repo.NewArray(), &PackageRepositoryProvider{Config: j(t, `{"package": [], "security-advisories": {"acme/package": [`+a+`]}}`)},
		&httpdl.TransportError{Msg: `The "https://example.org/security.json" file could not be downloaded: HTTP/1.1 502 Bad Gateway`, StatusCode: 502}}
	io := &bufferIO{}
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), io, nil)
	if _, err := f.Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0")), []repo.Repository{r}, &testRequest{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(io.Output(), "Security advisory data could not be fetched from some repositories") || !strings.Contains(io.Output(), "HTTP/1.1 502 Bad Gateway") {
		t.Fatal(io.Output())
	}
}

func TestRethrowsAdvisoryTransportErrorWhenUnreachableIsNotIgnored(t *testing.T) {
	r := &failingAdvisoryRepo{repo.NewArray(), &PackageRepositoryProvider{Config: j(t, `{"package": [], "security-advisories": {"acme/package": []}}`)},
		&httpdl.TransportError{Msg: "boom", StatusCode: 500}}
	pc := NewPolicyConfig(true, NewAdvisoriesPolicyConfig(true, AuditFail, nil, nil, nil), MalwarePolicyConfigDisabled(), NewAbandonedPolicyConfig(true, AuditFail, nil), CooldownPolicyConfigDisabled(), nil, IgnoreUnreachableNone())
	// PackageRepository::hasSecurityAdvisories() is count() > 0, so the empty package list still counts
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), pc, &bufferIO{}, nil)
	_, err := f.Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0")), []repo.Repository{r}, &testRequest{})
	if _, ok := isTransportError(err); !ok {
		t.Fatalf("expected TransportException, got %v", err)
	}
}

// recordingAdvisoryRepo records the names requested from getSecurityAdvisories.
type recordingAdvisoryRepo struct {
	*repo.Array
	*PackageRepositoryProvider
	requested []string
}

func (r *recordingAdvisoryRepo) GetSecurityAdvisories(pcm *ConstraintMap, allowPartial bool) (*AdvisoryMap, error) {
	r.requested = append(r.requested, pcm.Keys()...)
	return r.PackageRepositoryProvider.GetSecurityAdvisories(pcm, allowPartial)
}

func TestDoesNotLoadAdvisoriesForDevPackages(t *testing.T) {
	_, a := genAdvisory("acme/package", "CVE-2024-1234", ">=2.0.0")
	r := &recordingAdvisoryRepo{Array: repo.NewArray(), PackageRepositoryProvider: &PackageRepositoryProvider{Config: j(t, `{"package": [], "security-advisories": {"acme/package": [`+a+`]}}`)}}
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, nil)
	_, err := f.Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0"), newPkg("acme/path-package", "dev-main", "dev-main"), newPkg("acme/branch-package", "3.3.9999999.9999999-dev", "3.3.x-dev")), []repo.Repository{r}, &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.requested, []string{"acme/package", "acme/branch-package"}) {
		t.Fatal(r.requested)
	}
}

func TestFilterNumericDevBranchByAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", "<3.5")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	filtered, _ := f.Filter(newPool(newPkg("acme/package", "3.3.9999999.9999999-dev", "3.3.x-dev"), newPkg("acme/package", "3.6.9999999.9999999-dev", "3.6.x-dev")), []repo.Repository{r}, &testRequest{})
	if len(filtered.Packages) != 1 || filtered.Packages[0].PrettyVersion() != "3.6.x-dev" || removalType(filtered, "acme/package", "3.3.9999999.9999999-dev") != RemovalAdvisories {
		t.Fatal(packageNames(filtered.Packages))
	}
}

func TestFilterBranchAliasOfDevPackageByAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", "<3.5")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	dev := newPkg("acme/package", "dev-main", "dev-main")
	branchAlias := pkg.NewAlias(dev, "3.3.9999999.9999999-dev", "3.3.x-dev")
	rootAlias := pkg.NewAlias(dev, "3.3.0.0", "3.3.0")
	rootAlias.RootPackageAlias = true
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	filtered, _ := f.Filter(newPool(dev, branchAlias, rootAlias), []repo.Repository{r}, &testRequest{})
	assertSamePackages(t, pkgs(dev, rootAlias), filtered.Packages)
	if removalType(filtered, "acme/package", "3.3.9999999.9999999-dev") != RemovalAdvisories {
		t.Fatal("removal")
	}
}

func TestDefaultBranchAliasIsNotFilteredByAdvisories(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", ">=2.0")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	dev := newPkg("acme/package", "dev-main", "dev-main")
	alias := pkg.NewAlias(dev, semver.DefaultBranchAlias, semver.DefaultBranchAlias)
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	filtered, _ := f.Filter(newPool(dev, alias), []repo.Repository{r}, &testRequest{})
	assertSamePackages(t, pkgs(dev, alias), filtered.Packages)
}

func TestRootAliasOfDevPackageIsNotFilteredByAliasVersion(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a := genAdvisory("acme/package", "CVE-2024-1234", "*")
	r := advRepo(t, resolver, `{"acme/package": [`+a+`]}`)
	dev := newPkg("acme/package", "dev-main", "dev-main")
	alias := pkg.NewAlias(dev, "1.0.0.0", "1.0.0")
	alias.RootPackageAlias = true
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	filtered, _ := f.Filter(newPool(dev, alias), []repo.Repository{r}, &testRequest{})
	assertSamePackages(t, pkgs(dev, alias), filtered.Packages)
}

func TestFilterRootAliasByAliasedPackageVersion(t *testing.T) {
	resolver := &ProviderResolver{}
	_, a1 := genAdvisory("acme/package", "CVE-2024-1234", "<3.5")
	_, a2 := genAdvisory("acme/other", "CVE-2024-1235", "<3.5")
	r := advRepo(t, resolver, `{"acme/package": [`+a1+`], "acme/other": [`+a2+`]}`)
	vulnerable := newPkg("acme/package", "3.3.9999999.9999999-dev", "3.3.x-dev")
	vulnerableAlias := pkg.NewAlias(vulnerable, "3.9.0.0", "3.9.0")
	vulnerableAlias.RootPackageAlias = true
	safe := newPkg("acme/other", "3.6.9999999.9999999-dev", "3.6.x-dev")
	safeAlias := pkg.NewAlias(safe, "3.0.0.0", "3.0.0")
	safeAlias.RootPackageAlias = true
	f := NewSecurityAdvisoryPoolFilter(NewAuditor(nil), advPolicyConfig(true, true, nil, nil), nil, resolver)
	filtered, _ := f.Filter(newPool(vulnerable, vulnerableAlias, safe, safeAlias), []repo.Repository{r}, &testRequest{})
	assertSamePackages(t, pkgs(safe, safeAlias), filtered.Packages)
}

// ---------------------------------------------------- FilterListPoolFilterTest

func generatePackageRepository(t *testing.T, resolver *ProviderResolver, constraint string) *repo.Array {
	return packageRepository(t, resolver, `{"package": [], "filter": {"test-list": [{"package": "acme/package", "constraint": "`+constraint+`", "reason": "malware", "url": "https://example.org/malware/acme/package"}]}}`)
}

func TestFilterListFilterPackages(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": true}}`))
	r := generatePackageRepository(t, resolver, "1.0")
	f := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver)
	e1, e2 := newPkg("acme/package", "2.0.0.0", "2.0"), newPkg("acme/other", "1.0.0.0", "1.0")
	filtered, err := f.Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0"), e1, e2), &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(e1, e2), filtered.Packages)
	if removalType(filtered, "acme/package", "1.0.0.0") != RemovalFilterList || len(removals(filtered, RemovalFilterList)) != 1 {
		t.Fatal("removals")
	}
}

func TestUnfilteredPackagesConfig(t *testing.T) {
	for name, cfg := range map[string]string{
		"ignore-packages":         `{"ignore": ["acme/package"]}`,
		"ignore-packages-version": `{"ignore": {"acme/package": {"constraint": "*"}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			resolver := &ProviderResolver{}
			pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": `+cfg+`}}`))
			r := generatePackageRepository(t, resolver, "*")
			f := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver)
			e1, e2 := newPkg("acme/package", "1.0.0.0", "1.0"), newPkg("acme/package", "1.1.0.0", "1.1")
			filtered, err := f.Filter(newPool(e1, e2), &testRequest{})
			if err != nil {
				t.Fatal(err)
			}
			assertSamePackages(t, pkgs(e1, e2), filtered.Packages)
			if len(removals(filtered, RemovalFilterList)) != 0 {
				t.Fatal("removals")
			}
		})
	}
}

func TestUnfilteredPackagesConfigIntersection(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": {"ignore": {"acme/package": {"constraint": "<=1.0"}}}}}`))
	r := generatePackageRepository(t, resolver, ">=1.0")
	f := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver)
	expected := newPkg("acme/package", "1.0.0.0", "1.0")
	filtered, _ := f.Filter(newPool(expected, newPkg("acme/package", "1.1.0.0", "1.1")), &testRequest{})
	assertSamePackages(t, pkgs(expected), filtered.Packages)
	if len(removals(filtered, RemovalFilterList)) != 1 {
		t.Fatal("removals")
	}
}

func TestFilterWithAdditionalSources(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	h := &fakeHTTP{t: t, expected: []httpExpectation{{"https://example.org/malware/acme/package", "", `{"filter":[{"package":"acme/package","constraint":"3.0.0.0","reason":"malware"}]}`}}}
	pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": {"sources": {"source-list": {"type": "url", "url": "https://example.org/malware/acme/package"}}}}}`))
	r := generatePackageRepository(t, resolver, "0.0.1")
	f := NewFilterListPoolFilter(pc, h, BlockScopeUpdate, []repo.Repository{r}, nil, resolver)
	e1, e2 := newPkg("acme/package", "2.0.0.0", "2.0"), newPkg("acme/other", "1.0.0.0", "1.0")
	filtered, err := f.Filter(newPool(newPkg("acme/package", "3.0.0.0", "3.0"), e1, e2), &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(e1, e2), filtered.Packages)
	if removalType(filtered, "acme/package", "3.0.0.0") != RemovalFilterList || len(removals(filtered, RemovalFilterList)) != 1 {
		t.Fatal("removals")
	}
	h.assertComplete()
}

func malwareRepo(t *testing.T, resolver *ProviderResolver, name string) *repo.Array {
	return packageRepository(t, resolver, `{"package": [], "filter": {"malware": [{"package": "`+name+`", "constraint": "*", "reason": "malware", "url": "https://example.org/malware/`+name+`"}]}}`)
}

func TestInstallScopeFiltersLockedPackagesAgainstMalwareList(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	r := malwareRepo(t, resolver, "acme/locked")
	pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": true}}`))
	p := newPkg("acme/locked", "1.0.0.0", "1.0")
	req := &testRequest{locked: pkgs(p)} // fixLockedPackage

	installPool, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeInstall, []repo.Repository{r}, nil, resolver).Filter(newPool(p), req)
	if err != nil || len(installPool.Packages) != 0 || removalType(installPool, "acme/locked", "1.0.0.0") != RemovalFilterList {
		t.Fatal("install-scope filter must include locked packages")
	}
	updatePool, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver).Filter(newPool(p), req)
	if err != nil || len(updatePool.Packages) != 0 || removalType(updatePool, "acme/locked", "1.0.0.0") != RemovalFilterList {
		t.Fatal("update-scope filter must apply install-scope rules to locked packages")
	}
}

func TestUpdateScopeAppliesInstallScopeToPackagesInLockedRepository(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	r := malwareRepo(t, resolver, "acme/mirrored")
	pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": {"block": true, "block-scope": "install"}}}`))
	lockedPackage := newPkg("acme/mirrored", "1.0.0.0", "1.0")
	poolPackage := newPkg("acme/mirrored", "1.0.0.0", "1.0")
	req := &testRequest{lockedRepo: pkgs(lockedPackage)}
	updatePool, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver).Filter(newPool(poolPackage), req)
	if err != nil || len(updatePool.Packages) != 0 || removalType(updatePool, "acme/mirrored", "1.0.0.0") != RemovalFilterList {
		t.Fatal("packages from the locked repo must be checked against install-scope filter lists")
	}
}

func TestUpdateScopeIgnoresInstallOnlyListsForNonLockedPackages(t *testing.T) {
	clearEnv(t)
	resolver := &ProviderResolver{}
	r := malwareRepo(t, resolver, "acme/free")
	pc := mustPolicy(t, newConfig(t, `{"policy": {"malware": {"block": true, "block-scope": "install"}}}`))
	p := newPkg("acme/free", "1.0.0.0", "1.0")
	updatePool, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, nil, resolver).Filter(newPool(p), &testRequest{})
	if err != nil {
		t.Fatal(err)
	}
	assertSamePackages(t, pkgs(p), updatePool.Packages)
	if len(removals(updatePool, RemovalFilterList)) != 0 {
		t.Fatal("install-only lists must not block non-locked packages during update scope")
	}
}

// failingFilterRepo is the anonymous PackageRepository whose getFilter throws.
type failingFilterRepo struct {
	*repo.Array
	*PackageRepositoryProvider
	err error
}

func (f *failingFilterRepo) GetFilter(*ConstraintMap, []string) (*FilterByList, error) {
	return nil, f.err
}
func (f *failingFilterRepo) RepoName() string { return "unreachable filter list repo" }

const testListFilterRepo = `{"package": [], "filter": {"test-list": [{"package": "acme/package", "constraint": "*", "reason": "malware"}]}}`

func TestWarnsWhenUnreachableSourcesAreIgnored(t *testing.T) {
	clearEnv(t)
	r := &failingFilterRepo{repo.NewArray(), &PackageRepositoryProvider{Config: j(t, testListFilterRepo)},
		&httpdl.TransportError{Msg: `The "https://example.org/filter.json" file could not be downloaded: HTTP/1.1 500 Internal Server Error`, StatusCode: 500}}
	pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": true}}`))
	io := &bufferIO{}
	if _, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, io, nil).Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0")), &testRequest{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(io.Output(), "Filter list data could not be fetched from some sources") || !strings.Contains(io.Output(), "HTTP/1.1 500 Internal Server Error") {
		t.Fatal(io.Output())
	}
}

func TestRethrowsFilterTransportErrorWhenUnreachableIsNotIgnored(t *testing.T) {
	clearEnv(t)
	r := &failingFilterRepo{repo.NewArray(), &PackageRepositoryProvider{Config: j(t, testListFilterRepo)}, &httpdl.TransportError{Msg: "boom", StatusCode: 500}}
	pc := mustPolicy(t, newConfig(t, `{"policy": {"test-list": true, "ignore-unreachable": false}}`))
	_, err := NewFilterListPoolFilter(pc, &fakeHTTP{t: t}, BlockScopeUpdate, []repo.Repository{r}, &bufferIO{}, nil).Filter(newPool(newPkg("acme/package", "1.0.0.0", "1.0")), &testRequest{})
	if _, ok := isTransportError(err); !ok {
		t.Fatalf("expected TransportException, got %v", err)
	}
}

// ------------------------------------------------------ CooldownPoolFilterTest

func sevenDays() *CooldownPolicyConfig {
	return NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(7*24*3600))
}

func released(p *pkg.Package, date string) *pkg.Package {
	p.SetReleaseDate(date)
	return p
}

func filterCooldown(t *testing.T, f *CooldownPoolFilter, pool *Pool, req Request) *Pool {
	t.Helper()
	out, err := f.Filter(pool, req)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestCooldownFilterNewPackages(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	old := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), "2026-01-01 12:00:00")
	recent := released(newPkg("vendor/pkg", "2.0.0.0", "2.0.0"), "2026-01-14 12:00:00")
	out := filterCooldown(t, f, newPool(old, recent), &testRequest{})
	assertSamePackages(t, pkgs(old), out.Packages)
	if removalType(out, "vendor/pkg", "2.0.0.0") != RemovalCooldown {
		t.Fatal("removal")
	}
}

func TestAvailableInReportsTotalDaysForLongWaitsSpanningMultipleMonths(t *testing.T) {
	f := NewCooldownPoolFilter(NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(60*24*3600)), utc(cooldownNow))
	out := filterCooldown(t, f, newPool(released(newPkg("vendor/pkg", "2.0.0.0", "2.0.0"), "2026-01-14 12:00:00")), &testRequest{})
	if info := out.PolicyRemoved.Reason("vendor/pkg", "2.0.0.0").CooldownInfo(); info == nil || info.AvailableIn != "59 days" {
		t.Fatal(info)
	}
}

func TestPublishedDateTakesPrecedenceOverReleaseDate(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	p := loadPkg(t, "2.0.0.0", "2.0.0", `{"name": "vendor/pkg", "version": "2.0.0", "time": "2025-01-01 12:00:00", "published-time": "2026-01-14 12:00:00"}`)
	out := filterCooldown(t, f, newPool(p), &testRequest{})
	if len(out.Packages) != 0 || out.PolicyRemoved.Reason("vendor/pkg", "2.0.0.0").CooldownInfo().Source != "published-time" {
		t.Fatal("mismatch")
	}
}

func TestFallsBackToReleaseDateWhenNoPublishedDate(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	out := filterCooldown(t, f, newPool(released(newPkg("vendor/pkg", "2.0.0.0", "2.0.0"), "2026-01-14 12:00:00")), &testRequest{})
	if len(out.Packages) != 0 || out.PolicyRemoved.Reason("vendor/pkg", "2.0.0.0").CooldownInfo().Source != "time" {
		t.Fatal("mismatch")
	}
}

func internalIgnore() *IgnoreMap {
	return ignoreMap(rule("internal/*", nil, sp("Internal packages"), true, true))
}

func TestCooldownIgnoredPackagesAreNotFiltered(t *testing.T) {
	f := NewCooldownPoolFilter(NewCooldownPolicyConfig(true, AuditIgnore, internalIgnore(), i64(7*24*3600)), utc(cooldownNow))
	p := released(newPkg("internal/pkg", "1.0.0.0", "1.0.0"), "2026-01-14 12:00:00")
	out := filterCooldown(t, f, newPool(p), &testRequest{})
	assertSamePackages(t, pkgs(p), out.Packages)
	if removalType(out, "internal/pkg", "1.0.0.0") != "" {
		t.Fatal("removal")
	}
}

func TestIgnoreRuleDoesNotApplyToPackagesReplacingTheIgnoredName(t *testing.T) {
	f := NewCooldownPoolFilter(NewCooldownPolicyConfig(true, AuditIgnore, internalIgnore(), i64(7*24*3600)), utc(cooldownNow))
	replacer := loadPkg(t, "2.0.0.0", "2.0.0", `{"name": "evil/pkg", "version": "2.0.0", "time": "2026-01-14 12:00:00", "replace": {"internal/pkg": "*"}}`)
	out := filterCooldown(t, f, newPool(replacer), &testRequest{})
	if len(out.Packages) != 0 || removalType(out, "evil/pkg", "2.0.0.0") != RemovalCooldown {
		t.Fatal("mismatch")
	}
}

func TestCooldownDevVersionsAndDatelessAreNotFiltered(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	dev := released(newPkg("vendor/pkg", "dev-main", "dev-main"), "2026-01-14 12:00:00")
	assertSamePackages(t, pkgs(dev), filterCooldown(t, f, newPool(dev), &testRequest{}).Packages)
	dateless := newPkg("vendor/pkg", "1.0.0.0", "1.0.0")
	assertSamePackages(t, pkgs(dateless), filterCooldown(t, f, newPool(dateless), &testRequest{}).Packages)
}

func TestCooldownDisabledOrZeroConfigDoesNotFilter(t *testing.T) {
	for _, period := range []*int64{nil, i64(0)} {
		f := NewCooldownPoolFilter(NewCooldownPolicyConfig(true, AuditIgnore, nil, period), time.Time{})
		p := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), time.Now().UTC().Format("2006-01-02 15:04:05"))
		assertSamePackages(t, pkgs(p), filterCooldown(t, f, newPool(p), &testRequest{}).Packages)
	}
}

func TestCooldownInfoIsStored(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	out := filterCooldown(t, f, newPool(released(newPkg("vendor/pkg", "2.0.0.0", "2.0.0"), "2026-01-14 12:00:00")), &testRequest{})
	info := out.PolicyRemoved.Reason("vendor/pkg", "2.0.0.0").CooldownInfo()
	vs := out.PolicyRemoved.Versions("vendor/pkg", nil)
	if len(out.Packages) != 0 || info == nil || len(vs) != 1 || vs[0].Version != "2.0.0.0" || vs[0].Reason.PrettyVersion() != "2.0.0" || info.ReleaseDate == "" || info.AvailableIn == "" {
		t.Fatal("mismatch")
	}
}

func TestCooldownOldEnoughAndExactCutoffAreKept(t *testing.T) {
	oneDay := NewCooldownPolicyConfig(true, AuditIgnore, nil, i64(24*3600))
	f := NewCooldownPoolFilter(oneDay, utc(cooldownNow))
	for _, date := range []string{"2026-01-13 12:00:00", "2026-01-14 12:00:00"} {
		p := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), date)
		assertSamePackages(t, pkgs(p), filterCooldown(t, f, newPool(p), &testRequest{}).Packages)
	}
}

func TestCooldownRootPlatformAndLockedPackagesNotFiltered(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	root := released(newPkg("my/project", "1.0.0.0", "1.0.0"), "2026-01-14 12:00:00")
	root.IsRoot = true
	assertSamePackages(t, pkgs(root), filterCooldown(t, f, newPool(root), &testRequest{}).Packages)

	php := released(newPkg("php", "8.3.0.0", "8.3.0"), "2026-01-14 12:00:00")
	ext := released(newPkg("ext-json", "8.3.0.0", "8.3.0"), "2026-01-14 12:00:00")
	lib := released(newPkg("lib-openssl", "3.0.0.0", "3.0.0"), "2026-01-14 12:00:00")
	if out := filterCooldown(t, f, newPool(php, ext, lib), &testRequest{}); len(out.Packages) != 3 {
		t.Fatal("platform")
	}

	locked := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), "2026-01-14 12:00:00")
	assertSamePackages(t, pkgs(locked), filterCooldown(t, f, newPool(locked), &testRequest{locked: pkgs(locked)}).Packages)
}

func TestPackageWithMultipleNamesTrackedCorrectly(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	p := loadPkg(t, "2.0.0.0", "2.0.0", `{"name": "vendor/pkg", "version": "2.0.0", "time": "2026-01-14 12:00:00", "replace": {"vendor/replaced": "2.0.0"}}`)
	out := filterCooldown(t, f, newPool(p), &testRequest{})
	if len(out.Packages) != 0 || removalType(out, "vendor/pkg", "2.0.0.0") != RemovalCooldown || removalType(out, "vendor/replaced", "2.0.0.0") != RemovalCooldown {
		t.Fatal("removals")
	}
	vs := out.PolicyRemoved.Versions("vendor/replaced", nil)
	if out.PolicyRemoved.Reason("vendor/replaced", "2.0.0.0").CooldownInfo() == nil || len(vs) != 1 || vs[0].Reason.PackageName() != "vendor/pkg" {
		t.Fatal("primary name")
	}
}

func TestReplacerDoesNotHideTheReasonOfThePackageItReplaces(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	replacer := loadPkg(t, "2.0.0.0", "2.0.0", `{"name": "vendor/pkg", "version": "2.0.0", "time": "2026-01-14 12:00:00", "replace": {"vendor/replaced": "2.0.0"}}`)
	pool := newPool(replacer)
	pool.PolicyRemoved.Set("vendor/replaced", "2.0.0.0", RemovalReasonFilterList("vendor/replaced", "2.0.0", []*FilterListEntry{NewFilterListEntry("vendor/replaced", &semver.MatchAll{}, "malware", nil, nil, nil, nil)}))
	out := filterCooldown(t, f, pool, &testRequest{})
	if removalType(out, "vendor/pkg", "2.0.0.0") != RemovalCooldown || removalType(out, "vendor/replaced", "2.0.0.0") != RemovalFilterList {
		t.Fatal("mismatch")
	}
}

func TestDevAliasIsKeptTogetherWithItsDevTarget(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	branch := released(newPkg("vendor/pkg", "dev-main", "dev-main"), "2026-01-14 12:00:00")
	alias := pkg.NewAlias(branch, "1.0.0.0", "1.0.0")
	assertSamePackages(t, pkgs(branch, alias), filterCooldown(t, f, newPool(branch, alias), &testRequest{}).Packages)
}

func TestAliasOfWithheldVersionIsWithheldToo(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	base := released(newPkg("vendor/pkg", "2.0.0.0", "2.0.0"), "2026-01-14 12:00:00")
	alias := pkg.NewAlias(base, "2.0.9999999.9999999-dev", "2.0.x-dev")
	out := filterCooldown(t, f, newPool(base, alias), &testRequest{})
	if len(out.Packages) != 0 || removalType(out, "vendor/pkg", "2.0.0.0") != RemovalCooldown || removalType(out, "vendor/pkg", "2.0.9999999.9999999-dev") != RemovalCooldown {
		t.Fatal("mismatch")
	}
}

func createComposerRepository(t *testing.T, requirePublishedTime bool) *repo.Composer {
	t.Helper()
	cfg := phpjson.NewMap()
	cfg.Set("url", "https://repo.example.org")
	cfg.Set("require-published-time", requirePublishedTime)
	r, err := repo.NewComposer(cfg, t.TempDir(), true, nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestComposerRepositoryVersionWithoutPublishedTimeIsRejected(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	p := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), "2020-01-01 00:00:00")
	p.Repo = createComposerRepository(t, true)
	_, err := f.Filter(newPool(p), &testRequest{})
	expectError(t, err, RuntimeException, "vendor/pkg 1.0.0 from composer repo (https://repo.example.org) has no published-time")
}

func TestComposerRepositoryCanWaivePublishedTime(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	p := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), "2026-01-14 12:00:00")
	p.Repo = createComposerRepository(t, false)
	out := filterCooldown(t, f, newPool(p), &testRequest{})
	if len(out.Packages) != 0 || out.PolicyRemoved.Reason("vendor/pkg", "1.0.0.0").CooldownInfo().Source != "time" {
		t.Fatal("mismatch")
	}
}

func TestOtherRepositoriesNeedNoPublishedTime(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	p := released(newPkg("vendor/pkg", "1.0.0.0", "1.0.0"), "2020-01-01 00:00:00")
	repo.NewArray().AddPackage(p)
	assertSamePackages(t, pkgs(p), filterCooldown(t, f, newPool(p), &testRequest{}).Packages)
}

func TestDevAndIgnoredVersionsFromComposerRepositoriesNeedNoPublishedTime(t *testing.T) {
	f := NewCooldownPoolFilter(NewCooldownPolicyConfig(true, AuditIgnore, internalIgnore(), i64(7*24*3600)), utc(cooldownNow))
	r := createComposerRepository(t, true)
	branch := newPkg("vendor/pkg", "dev-main", "dev-main")
	branch.Repo = r
	ignored := newPkg("internal/pkg", "1.0.0.0", "1.0.0")
	ignored.Repo = r
	assertSamePackages(t, pkgs(branch, ignored), filterCooldown(t, f, newPool(branch, ignored), &testRequest{}).Packages)
}

func TestUnverifiedVersionsAreGroupedByRepository(t *testing.T) {
	f := NewCooldownPoolFilter(sevenDays(), utc(cooldownNow))
	repoA, repoB := repo.NewArray(), repo.NewArray()
	timeOnly := released(newPkg("vendor/a", "1.0.0.0", "1.0.0"), "2020-01-01 00:00:00")
	repoA.AddPackage(timeOnly)
	noDate := newPkg("vendor/b", "1.0.0.0", "1.0.0")
	repoB.AddPackage(noDate)
	verified := loadPkg(t, "1.0.0.0", "1.0.0", `{"name": "vendor/c", "version": "1.0.0", "published-time": "2020-01-01 00:00:00"}`)
	repoB.AddPackage(verified)
	branch := newPkg("vendor/d", "dev-main", "dev-main")
	repoB.AddPackage(branch)
	filterCooldown(t, f, newPool(timeOnly, noDate, verified, branch), &testRequest{})

	var got []string
	f.UnverifiedVersions().Each(func(k string, v []string) { got = append(got, k+"=>"+strings.Join(v, ",")) })
	want := []string{repoA.RepoName() + "=>vendor/a (1.0.0)", repoB.RepoName() + "=>vendor/b (1.0.0)"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

// ---------------------------------------------------- PolicyRemovalReasonTest

func TestCombinePicksChronologicallyEarliestCooldownAcrossOffsets(t *testing.T) {
	c, err := CombineRemovalReasons([]*PolicyRemovalReason{
		RemovalReasonCooldown("vendor/pkg", "1.0.0", "2026-01-12T01:00:00+02:00", "5 days", "time"),
		RemovalReasonCooldown("vendor/pkg", "2.0.0", "2026-01-12T00:30:00+00:00", "6 days", "time"),
	})
	if err != nil || c.PrettyVersion() != "1.0.0" || c.CooldownInfo().AvailableIn != "5 days" {
		t.Fatal("mismatch")
	}
}

func TestCombineMergesAdvisoriesOfAllVersions(t *testing.T) {
	shared := NewPartialSecurityAdvisory("vendor/pkg", "PKSA-aaaa-aaaa-aaaa", semver.NewConstraint("<", "2.0.0.0"))
	c, err := CombineRemovalReasons([]*PolicyRemovalReason{
		RemovalReasonAdvisories("vendor/pkg", "1.0.0", []*Advisory{shared, NewPartialSecurityAdvisory("vendor/pkg", "CVE-2026-1", semver.NewConstraint("<", "1.1.0.0"))}),
		RemovalReasonAdvisories("vendor/pkg", "1.5.0", []*Advisory{shared}),
	})
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, a := range c.Advisories() {
		ids = append(ids, a.AdvisoryID)
	}
	if !reflect.DeepEqual(ids, []string{"PKSA-aaaa-aaaa-aaaa", "CVE-2026-1"}) || !strings.Contains(c.Remedy(), "Review the advisory details above") {
		t.Fatal(ids)
	}
}

func TestCombineRejectsMixedPolicies(t *testing.T) {
	_, err := CombineRemovalReasons([]*PolicyRemovalReason{
		RemovalReasonAbandoned("vendor/pkg", "1.0.0"),
		RemovalReasonCooldown("vendor/pkg", "2.0.0", "2026-01-12T00:30:00+00:00", "6 days", "time"),
	})
	expectError(t, err, LogicException, "")
}

func TestFilterListRendersEachList(t *testing.T) {
	r := RemovalReasonFilterList("vendor/pkg", "1.0.0", []*FilterListEntry{
		NewFilterListEntry("vendor/pkg", &semver.MatchAll{}, "malware", sp("https://example.org/m"), sp("looks suspicious"), sp("PKG-1"), sp("aikido")),
		NewFilterListEntry("vendor/pkg", &semver.MatchAll{}, "internal", nil, nil, nil, nil),
	})
	if r.Verb(true) != "were" {
		t.Fatal("verb")
	}
	if got := r.Description(); got != "flagged as malware reported by aikido (see https://example.org/m) reason: looks suspicious, filtered by internal" {
		t.Fatal(got)
	}
	if got := r.Remedy(); got != `To ignore filters for this package, add the package to the "policy.malware.ignore" and "policy.internal.ignore" config.` {
		t.Fatal(got)
	}
	if got := r.OffSwitch(); got != `"policy.malware.block" and "policy.internal.block"` {
		t.Fatal(got)
	}
}

func TestCooldownMentionsWeakerTimeSource(t *testing.T) {
	verified := RemovalReasonCooldown("vendor/pkg", "1.0.0", "2026-01-12T00:30:00+00:00", "5 days", "published-time")
	unverified := RemovalReasonCooldown("vendor/pkg", "1.0.0", "2026-01-12T00:30:00+00:00", "5 days", "time")
	if verified.Description() != `still in the cooldown period configured in "policy.cooldown" (available in 5 days)` ||
		!strings.Contains(unverified.Description(), "based on the package-supplied time field") || verified.OffSwitch() != `"policy.cooldown.block"` {
		t.Fatal("mismatch")
	}
}
