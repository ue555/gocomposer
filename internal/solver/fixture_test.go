package solver

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

var reSection = regexp.MustCompile(`(?:^|\n*)--([A-Z-]+)--\n`)

// readFixture splits an upstream .test file into its sections.
func readFixture(t *testing.T, file string) map[string]string {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	locs := reSection.FindAllStringSubmatchIndex(s, -1)
	out := map[string]string{}
	for i, loc := range locs {
		name := s[loc[2]:loc[3]]
		end := len(s)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		out[name] = s[loc[1]:end]
	}
	return out
}

func mustDecode(t *testing.T, s string) phpjson.Value {
	v, err := phpjson.Decode([]byte(s))
	if err != nil {
		t.Fatalf("invalid json: %v\n%s", err, s)
	}
	return v
}

func TestUpstreamPoolBuilderFixtures(t *testing.T) {
	files, _ := filepath.Glob("testdata/poolbuilder/*.test")
	if len(files) == 0 {
		t.Fatal("no fixtures")
	}
	for _, file := range files {
		file := file
		t.Run(filepath.Base(file), func(t *testing.T) {
			runPoolBuilderFixture(t, file)
		})
	}
}

func runPoolBuilderFixture(t *testing.T, file string) {
	sec := readFixture(t, file)
	root := phpjson.NewMap()
	if strings.TrimSpace(sec["ROOT"]) != "" {
		root = phpjson.AsMap(mustDecode(t, sec["ROOT"]))
	}
	request := phpjson.AsMap(mustDecode(t, sec["REQUEST"]))
	packageRepos := mustDecode(t, sec["PACKAGE-REPOS"])
	var fixed phpjson.Value = []phpjson.Value{}
	if strings.TrimSpace(sec["FIXED"]) != "" {
		fixed = mustDecode(t, sec["FIXED"])
	}
	expect := mustDecode(t, sec["EXPECT"])
	expectOptimized := expect
	if strings.TrimSpace(sec["EXPECT-OPTIMIZED"]) != "" {
		expectOptimized = mustDecode(t, sec["EXPECT-OPTIMIZED"])
	}

	var rootAliases []pkg.RootAlias
	if a, ok := root.Get("aliases"); ok && phpjson.Truthy(a) {
		phpjson.Each(a, func(_ string, v phpjson.Value) {
			m := phpjson.AsMap(v)
			rootAliases = append(rootAliases, pkg.RootAlias{
				Package: m.Str("package"), Version: semver.MustNormalize(m.Str("version")),
				Alias: m.Str("alias"), AliasNormalized: semver.MustNormalize(m.Str("alias")),
			})
		})
	}
	minStab := "stable"
	if s := root.Str("minimum-stability"); s != "" {
		minStab = s
	}
	flags := map[string]int{}
	if f := root.Map("stability-flags"); f != nil {
		for _, k := range f.Keys() {
			flags[k] = pkg.Stabilities[f.Str(k)]
		}
	}
	refs := map[string]string{}
	if r := root.Map("references"); r != nil {
		for _, k := range r.Keys() {
			refs[k] = r.Str(k)
		}
	}

	loader := &pkg.Loader{LoadOptions: true}
	ids := map[*pkg.Package]string{}
	load := func(v phpjson.Value) *pkg.Package {
		m := phpjson.AsMap(v).Clone()
		id := ""
		if idv, ok := m.Get("id"); ok && phpjson.Truthy(idv) {
			id, _ = phpjson.ScalarString(idv)
			m.Delete("id")
		}
		p, err := loader.Load(m)
		if err != nil {
			t.Fatal(err)
		}
		if id != "" {
			ids[p] = id
		}
		return p
	}

	rs := NewRepositorySet(minStab, flags, rootAliases, refs, nil, nil)
	skip := ""
	phpjson.Each(packageRepos, func(_ string, packages phpjson.Value) {
		if m, ok := packages.(*phpjson.Map); ok {
			if m.Isset("type") || m.Isset("canonical") || m.Isset("only") || m.Isset("exclude") {
				skip = "fixture uses a typed or filtered repository"
				return
			}
		}
		a := repo.NewArray()
		rs.AddRepository(a)
		phpjson.Each(packages, func(_ string, pv phpjson.Value) { a.AddPackage(load(pv)) })
	})
	if skip != "" {
		t.Skip(skip)
	}
	locked := repo.NewLockArray()
	rs.AddRepository(locked)
	if l, ok := request.Get("locked"); ok {
		phpjson.Each(l, func(_ string, pv phpjson.Value) { locked.AddPackage(load(pv)) })
	}
	req := NewRequest(locked)
	if rq := request.Map("require"); rq != nil {
		for _, k := range rq.Keys() {
			req.RequireName(k, semver.MustParseConstraints(rq.Str(k)))
		}
	}
	if al, ok := request.Get("allowList"); ok {
		var list []string
		phpjson.Each(al, func(_ string, v phpjson.Value) { s, _ := phpjson.ScalarString(v); list = append(list, s) })
		mode := UpdateOnlyListed
		if v, ok := request.Get("allowTransitiveDepsNoRootRequire"); ok && v == true {
			mode = UpdateListedWithTransitiveDepsNoRoot
		}
		if v, ok := request.Get("allowTransitiveDeps"); ok && v == true {
			mode = UpdateListedWithTransitiveDeps
		}
		req.SetUpdateAllowList(list, mode)
	}
	phpjson.Each(fixed, func(_ string, pv phpjson.Value) { req.FixPackage(load(pv)) })

	pool, err := rs.CreatePool(req, NewPoolBuilder())
	if err != nil {
		t.Fatal(err)
	}
	check := func(label string, pool *Pool, want phpjson.Value) {
		var got []string
		for i := 1; i <= pool.Count(); i++ {
			got = append(got, describePoolPackage(pool.PackageByID(i), ids))
		}
		var exp []string
		phpjson.Each(want, func(_ string, v phpjson.Value) { s, _ := phpjson.ScalarString(v); exp = append(exp, s) })
		sort.Strings(got)
		sort.Strings(exp)
		if strings.Join(got, "\n") != strings.Join(exp, "\n") {
			t.Errorf("%s pool mismatch\n got: %q\nwant: %q", label, got, exp)
		}
	}
	check("unoptimized", pool, expect)
	optimized := NewPoolOptimizer(NewDefaultPolicy(false, false, nil)).Optimize(req, pool)
	check("optimized", optimized, expectOptimized)
}

func describePoolPackage(p *pkg.Package, ids map[*pkg.Package]string) string {
	if id, ok := ids[p]; ok {
		return id
	}
	suffix := ""
	if p.SourceReference() != "" {
		suffix = "#" + p.SourceReference()
	}
	if r, ok := p.Repo.(repo.Repository); ok && r.Kind() == repo.KindLock {
		suffix += " (locked)"
	}
	if p.AliasOf != nil {
		if id, ok := ids[p.AliasOf]; ok {
			return p.Name + "-" + p.Version() + suffix + " (alias of " + id + ")"
		}
		return p.Name + "-" + p.Version() + suffix + " (alias of " + p.AliasOf.Version() + ")"
	}
	return p.Name + "-" + p.Version() + suffix
}

func TestUpstreamPoolOptimizerFixtures(t *testing.T) {
	files, _ := filepath.Glob("testdata/pooloptimizer/*.test")
	if len(files) == 0 {
		t.Fatal("no fixtures")
	}
	for _, file := range files {
		file := file
		t.Run(filepath.Base(file), func(t *testing.T) {
			sec := readFixture(t, file)
			request := phpjson.AsMap(mustDecode(t, sec["REQUEST"]))
			loader := &pkg.Loader{}
			load := func(v phpjson.Value) *pkg.Package {
				p, err := loader.Load(phpjson.AsMap(v))
				if err != nil {
					t.Fatal(err)
				}
				return p
			}
			loadAll := func(v phpjson.Value) []*pkg.Package {
				var out []*pkg.Package
				phpjson.Each(v, func(_ string, pv phpjson.Value) {
					p := load(pv)
					out = append(out, p)
					if p.AliasOf != nil {
						out = append(out, p.AliasOf)
					}
				})
				return out
			}
			before := loadAll(mustDecode(t, sec["POOL-BEFORE"]))
			after := loadAll(mustDecode(t, sec["POOL-AFTER"]))
			req := NewRequest(repo.NewLockArray())
			if l, ok := request.Get("locked"); ok {
				phpjson.Each(l, func(_ string, pv phpjson.Value) { req.LockPackage(load(pv)) })
			}
			if l, ok := request.Get("fixed"); ok {
				phpjson.Each(l, func(_ string, pv phpjson.Value) { req.FixPackage(load(pv)) })
			}
			if rq := request.Map("require"); rq != nil {
				for _, k := range rq.Keys() {
					req.RequireName(k, semver.MustParseConstraints(rq.Str(k)))
				}
			}
			preferStable, _ := request.Get("preferStable")
			preferLowest, _ := request.Get("preferLowest")
			pool := NewPool(before, nil)
			out := NewPoolOptimizer(NewDefaultPolicy(preferStable == true, preferLowest == true, nil)).Optimize(req, pool)
			reduce := func(ps []*pkg.Package) []string {
				var r []string
				for _, p := range ps {
					s := p.Name + "@" + p.Version()
					if p.AliasOf != nil {
						s += " (alias of " + p.AliasOf.Version() + ")"
					}
					r = append(r, s)
				}
				sort.Strings(r)
				return r
			}
			got, want := reduce(out.Packages()), reduce(after)
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Errorf("optimized pool mismatch\n got: %q\nwant: %q", got, want)
			}
		})
	}
}
