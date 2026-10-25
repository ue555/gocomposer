package repo

// Port of RepositoryInterface::search with modes and type filtering for
// ArrayRepository, ComposerRepository, FilterRepository and
// PlatformRepository (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"net/url"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// Search modes (RepositoryInterface::SEARCH_*).
const (
	SearchFulltext = 0
	SearchName     = 1
	SearchVendor   = 2
)

// SearchEntry is one search() result row; Description is nil for PHP null,
// Abandoned is nil, true or the replacement package name.
type SearchEntry struct {
	Name        string
	Description *string
	Abandoned   phpjson.Value
	URL         string
	// Raw is the complete result object of a search API response.
	Raw *phpjson.Map
}

func strp(s string) *string { return &s }

func modeRegex(query string, mode int) *regexp.Regexp {
	return searchRegex(query, mode == SearchFulltext)
}

// arraySearch is ArrayRepository::search.
func arraySearch(pkgs []*pkg.Package, query string, mode int, typ *string) []SearchEntry {
	re := modeRegex(query, mode)
	seen := map[string]bool{}
	var out []SearchEntry
	for _, p := range pkgs {
		name := p.Name
		if mode == SearchVendor {
			name, _, _ = strings.Cut(name, "/")
		}
		if seen[name] {
			continue
		}
		if typ != nil {
			t := p.Type()
			if t == "" {
				t = "library"
			}
			if t != *typ {
				continue
			}
		}
		if re.MatchString(name) || (mode == SearchFulltext && re.MatchString(strings.Join(p.Keywords(), " ")+" "+p.Description())) {
			seen[name] = true
			if mode == SearchVendor {
				out = append(out, SearchEntry{Name: name})
				continue
			}
			e := SearchEntry{Name: p.PrettyName, Description: strp(p.Description())}
			if p.IsAbandoned() {
				if r := p.ReplacementPackage(); r != "" {
					e.Abandoned = r
				} else {
					e.Abandoned = true
				}
			}
			out = append(out, e)
		}
	}
	return out
}

// SearchWithMode runs $repository->search($query, $mode, $type).
func SearchWithMode(r Repository, query string, mode int, typ *string) ([]SearchEntry, error) {
	switch t := r.(type) {
	case *Composer:
		return t.searchMode(query, mode, typ)
	case *Filter:
		res, err := SearchWithMode(t.Repo, query, mode, typ)
		if err != nil {
			return nil, err
		}
		var out []SearchEntry
		for _, e := range res {
			if t.isAllowed(e.Name) {
				out = append(out, e)
			}
		}
		return out, nil
	}
	if r.Kind() == KindPlatform && mode == SearchVendor {
		return nil, nil
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	return arraySearch(ps, query, mode, typ), nil
}

func (r *Composer) searchMode(query string, mode int, typ *string) ([]SearchEntry, error) {
	data, err := r.loadRootServerFile()
	if err != nil {
		return nil, err
	}
	if s := data.Str("search"); s != "" && mode == SearchFulltext {
		t := ""
		if typ != nil {
			t = *typ
		}
		u := strings.NewReplacer("%query%", url.QueryEscape(query), "%type%", t).Replace(r.canonicalizeURL(s))
		resp, err := r.dl.Get(u, r.headers()...)
		if err != nil {
			return nil, err
		}
		v, err := phpjson.Decode(resp.Body)
		if err != nil {
			return nil, err
		}
		r.responseWarnings(v)
		var out []SearchEntry
		if res, ok := phpjson.AsMap(v).Get("results"); ok {
			phpjson.Each(res, func(_ string, e phpjson.Value) {
				m := phpjson.AsMap(e)
				if m == nil || phpjson.Truthy(mustGet(m, "virtual")) {
					return
				}
				se := SearchEntry{Name: m.Str("name"), URL: m.Str("url"), Raw: m}
				if d, ok := m.Get("description"); ok && d != nil {
					s, _ := phpjson.ScalarString(d)
					se.Description = &s
				}
				if a, ok := m.Get("abandoned"); ok {
					se.Abandoned = a
				}
				out = append(out, se)
			})
		}
		return out, nil
	}
	if mode == SearchVendor {
		names, err := r.allPackageNames()
		if err != nil {
			return nil, err
		}
		re := searchRegex(query, false)
		seen := map[string]bool{}
		var out []SearchEntry
		for _, n := range names {
			v, _, _ := strings.Cut(n, "/")
			if seen[v] {
				continue
			}
			seen[v] = true
			if re.MatchString(v) {
				out = append(out, SearchEntry{Name: v, Description: strp("")})
			}
		}
		return out, nil
	}
	if r.hasProviders || r.lazyProvidersURL != "" {
		if res, ok, err := r.searchListAPI(query); ok {
			return res, err
		}
		re := searchRegex(query, false)
		var out []SearchEntry
		names, err := r.PackageNames("")
		if err != nil {
			return nil, err
		}
		for _, n := range names {
			if re.MatchString(n) {
				out = append(out, SearchEntry{Name: n, Description: strp("")})
			}
		}
		return out, nil
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	// ComposerRepository calls parent::search($query, $mode) without the type
	return arraySearch(ps, query, mode, nil), nil
}

// allPackageNames is ComposerRepository::getPackageNames() without filter.
func (r *Composer) allPackageNames() ([]string, error) {
	if _, err := r.loadRootServerFile(); err != nil {
		return nil, err
	}
	if r.lazyProvidersURL != "" || r.hasProviders {
		return r.PackageNames("")
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range ps {
		out = append(out, p.PrettyName)
	}
	return out, nil
}
