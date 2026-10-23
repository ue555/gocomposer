package repo

// Port of RepositoryInterface::search (SEARCH_FULLTEXT) for the array and
// composer repositories (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"net/url"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// SearchResult is one search() entry.
type SearchResult struct {
	Name        string
	Description string
}

var reSpaces = regexp.MustCompile(`\s+`)

func searchRegex(query string, quote bool) *regexp.Regexp {
	if quote {
		query = regexp.QuoteMeta(query)
	}
	parts := reSpaces.Split(query, -1)
	re, err := regexp.Compile(`(?i)(?:` + strings.Join(parts, "|") + `)`)
	if err != nil {
		return regexp.MustCompile(`$^`)
	}
	return re
}

// searchPackages is ArrayRepository::search in fulltext mode.
func searchPackages(pkgs []*pkg.Package, query string) []SearchResult {
	re := searchRegex(query, true)
	seen := map[string]bool{}
	var out []SearchResult
	for _, p := range pkgs {
		if seen[p.Name] {
			continue
		}
		if re.MatchString(p.Name) || re.MatchString(strings.Join(p.Keywords(), " ")+" "+p.Description()) {
			seen[p.Name] = true
			out = append(out, SearchResult{Name: p.PrettyName, Description: p.Description()})
		}
	}
	return out
}

// Search runs a fulltext search on a repository.
func Search(r Repository, query string) ([]SearchResult, error) {
	switch t := r.(type) {
	case *Composer:
		return t.search(query)
	case *Filter:
		return Search(t.Repo, query)
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	return searchPackages(ps, query), nil
}

func (r *Composer) search(query string) ([]SearchResult, error) {
	data, err := r.loadRootServerFile()
	if err != nil {
		return nil, err
	}
	if s := data.Str("search"); s != "" {
		u := strings.NewReplacer("%query%", url.QueryEscape(query), "%type%", "").Replace(r.canonicalizeURL(s))
		resp, err := r.dl.Get(u, r.headers()...)
		if err != nil {
			return nil, err
		}
		v, err := phpjson.Decode(resp.Body)
		if err != nil {
			return nil, err
		}
		var out []SearchResult
		if res, ok := phpjson.AsMap(v).Get("results"); ok {
			if list, ok := res.([]phpjson.Value); ok {
				out = nil
				for _, e := range list {
					m := phpjson.AsMap(e)
					if m == nil || phpjson.Truthy(mustGet(m, "virtual")) {
						continue
					}
					out = append(out, SearchResult{Name: m.Str("name"), Description: m.Str("description")})
				}
			}
		}
		return out, nil
	}
	if r.hasProviders || r.lazyProvidersURL != "" {
		re := searchRegex(query, false)
		var out []SearchResult
		for _, n := range r.packageNames() {
			if re.MatchString(n) {
				out = append(out, SearchResult{Name: n})
			}
		}
		return out, nil
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	return searchPackages(ps, query), nil
}

// packageNames is getPackageNames() without a filter for the cases gocomposer
// supports (available-packages lists and partial packages).
func (r *Composer) packageNames() []string {
	if r.availablePackages != nil {
		return r.availableNames
	}
	if r.partialByName != nil {
		var out []string
		for _, k := range r.partialOrder {
			out = append(out, k)
		}
		return out
	}
	return nil
}

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }
