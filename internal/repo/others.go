package repo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
)

// NewPackageRepository is PackageRepository (type "package").
func NewPackageRepository(config *phpjson.Map) (*Array, error) {
	pv, _ := config.Get("package")
	var defs []phpjson.Value
	switch t := pv.(type) {
	case []phpjson.Value:
		defs = t
	case *phpjson.Map:
		if t.Len() > 0 {
			if _, isInt := phpjson.IntKey(t.Keys()[0]); isInt {
				phpjson.Each(t, func(_ string, v phpjson.Value) { defs = append(defs, v) })
			} else {
				defs = []phpjson.Value{t}
			}
		}
	}
	a := &Array{kind: KindPackage}
	a.init = func() error {
		loader := &pkg.Loader{LoadOptions: true}
		for _, d := range defs {
			m := phpjson.ToMap(d)
			if m == nil {
				return fmt.Errorf("A repository of type \"package\" contains an invalid package definition")
			}
			p, err := loader.Load(m)
			if err == nil {
				err = ValidatePackageDefinition(m)
			}
			if err != nil {
				enc, _ := phpjson.Encode(m, 0)
				return &InvalidRepositoryError{"A repository of type \"package\" contains an invalid package definition: " + err.Error() + "\n\nInvalid package definition:\n" + enc}
			}
			a.AddPackage(p)
		}
		return nil
	}
	pkgsCount := func() int {
		ps, _ := a.Packages()
		return len(ps)
	}
	_ = pkgsCount
	return a, nil
}

// InvalidRepositoryError mirrors InvalidRepositoryException.
type InvalidRepositoryError struct{ Msg string }

func (e *InvalidRepositoryError) Error() string { return e.Msg }

// PackageRepoName mirrors PackageRepository::getRepoName.
func PackageRepoName(a *Array) string {
	ps, _ := a.Packages()
	n := len(ps)
	s := ""
	if n > 1 {
		s = "s"
	}
	return "package repo (defining " + strconv.Itoa(n) + " package" + s + ")"
}

// ValidatePackageDefinition applies the parts of ValidatingArrayLoader that
// reject definitions (name format); warnings are not reported here.
func ValidatePackageDefinition(m *phpjson.Map) error {
	name := m.Str("name")
	if msg := PackageNamingError(name, false); msg != "" {
		return fmt.Errorf("name : %s", msg)
	}
	return nil
}

var (
	reNameValidRE = regexp.MustCompile(`(?i)^[a-z0-9](?:[_.-]?[a-z0-9]+)*/[a-z0-9](?:(?:[_.]|-{1,2})?[a-z0-9]+)*$`)
	reNameJSON    = regexp.MustCompile(`\.json$`)
	reNameUpper   = regexp.MustCompile(`[A-Z]`)
)

// suggestLowercase is Preg::replace('{(?:([a-z])([A-Z])|([A-Z])([A-Z][a-z]))}', '\\1\\3-\\2\\4', $name).
func suggestLowercase(name string) string {
	var sb strings.Builder
	i := 0
	isLower := func(c byte) bool { return c >= 'a' && c <= 'z' }
	isUpper := func(c byte) bool { return c >= 'A' && c <= 'Z' }
	for i < len(name) {
		if i+1 < len(name) && isLower(name[i]) && isUpper(name[i+1]) {
			sb.WriteByte(name[i])
			sb.WriteByte('-')
			sb.WriteByte(name[i+1])
			i += 2
			continue
		}
		if i+2 < len(name) && isUpper(name[i]) && isUpper(name[i+1]) && isLower(name[i+2]) {
			sb.WriteByte(name[i])
			sb.WriteByte('-')
			sb.WriteByte(name[i+1])
			sb.WriteByte(name[i+2])
			i += 3
			continue
		}
		sb.WriteByte(name[i])
		i++
	}
	return sb.String()
}

// PackageNamingError is ValidatingArrayLoader::hasPackageNamingError.
func PackageNamingError(name string, isLink bool) string {
	if IsPlatformPackage(name) {
		return ""
	}
	if !reNameValidRE.MatchString(name) {
		return name + ` is invalid, it should have a vendor name, a forward slash, and a package name. The vendor and package name can be words separated by -, . or _. The complete name should match "^[a-z0-9]([_.-]?[a-z0-9]+)*/[a-z0-9](([_.]?|-{0,2})[a-z0-9]+)*$".`
	}
	reserved := []string{"nul", "con", "prn", "aux", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	bits := strings.SplitN(strings.ToLower(name), "/", 2)
	for _, r := range reserved {
		if bits[0] == r || (len(bits) > 1 && bits[1] == r) {
			return name + " is reserved, package and vendor names can not match any of: " + strings.Join(reserved, ", ") + "."
		}
	}
	if reNameJSON.MatchString(name) {
		return name + " is invalid, package names can not end in .json, consider renaming it or perhaps using a -json suffix instead."
	}
	if reNameUpper.MatchString(name) {
		if isLink {
			return name + " is invalid, it should not contain uppercase characters. Please use " + strings.ToLower(name) + " instead."
		}
		return name + " is invalid, it should not contain uppercase characters. We suggest using " + strings.ToLower(suggestLowercase(name)) + " instead."
	}
	return ""
}

// Path is PathRepository.
type Path struct {
	*Array
	url        string
	rawURL     string
	options    *phpjson.Map
	repoConfig *phpjson.Map
	// GuessVersion returns (pretty, normalized, ok) for a directory.
	GuessVersion func(dir string, config *phpjson.Map) (string, string, bool)
}

// NewPath is PathRepository::__construct.
func NewPath(repoConfig *phpjson.Map, baseDir string) (*Path, error) {
	if !repoConfig.Isset("url") {
		return nil, fmt.Errorf("You must specify the `url` configuration for the path repository")
	}
	raw := repoConfig.Str("url")
	u := expandPath(raw)
	options := phpjson.NewMap()
	if o := repoConfig.Map("options"); o != nil {
		options = o.Clone()
	}
	if !options.Has("relative") {
		options.Set("relative", !filepath.IsAbs(u))
	}
	p := &Path{Array: &Array{kind: KindPath}, url: u, rawURL: raw, options: options, repoConfig: repoConfig}
	p.Array.self = p
	p.Array.name = "path repo (" + raw + ")"
	p.Array.init = func() error { return p.initialize(baseDir) }
	return p, nil
}

func (p *Path) RepoName() string { return "path repo (" + p.rawURL + ")" }

func expandPath(u string) string {
	if strings.HasPrefix(u, "~/") || u == "~" {
		if h := os.Getenv("HOME"); h != "" {
			return h + u[1:]
		}
	}
	return u
}

var reGlobChars = regexp.MustCompile(`[*{}]`)

func (p *Path) urlMatches(baseDir string) ([]string, error) {
	pattern := p.url
	abs := pattern
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(baseDir, abs)
	}
	var matches []string
	for _, alt := range expandBraces(abs) {
		ms, err := filepath.Glob(alt)
		if err != nil {
			return nil, err
		}
		for _, m := range ms {
			if st, err := os.Stat(m); err == nil && st.IsDir() {
				matches = append(matches, m)
			}
		}
	}
	sort.Strings(matches)
	// return paths in the same form as the configured url (relative stays relative)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if !filepath.IsAbs(pattern) {
			if rel, err := filepath.Rel(baseDir, m); err == nil {
				m = rel
			}
		}
		out = append(out, strings.TrimRight(filepath.ToSlash(m), "/"))
	}
	return out, nil
}

// expandBraces implements GLOB_BRACE expansion for {a,b} groups.
func expandBraces(s string) []string {
	start := strings.IndexByte(s, '{')
	if start < 0 {
		return []string{s}
	}
	depth := 0
	end := -1
	for i := start; i < len(s); i++ {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}
	if end < 0 {
		return []string{s}
	}
	var parts []string
	depth = 0
	last := start + 1
	for i := start + 1; i < end; i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[last:i])
				last = i + 1
			}
		}
	}
	parts = append(parts, s[last:end])
	var out []string
	for _, part := range parts {
		out = append(out, expandBraces(s[:start]+part+s[end+1:])...)
	}
	return out
}

func (p *Path) initialize(baseDir string) error {
	matches, err := p.urlMatches(baseDir)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		if reGlobChars.MatchString(p.url) {
			u := p.url
			for reGlobChars.MatchString(u) {
				u = filepath.Dir(u)
			}
			check := u
			if !filepath.IsAbs(check) {
				check = filepath.Join(baseDir, check)
			}
			if st, err := os.Stat(check); err == nil && st.IsDir() {
				return nil
			}
		}
		return fmt.Errorf("The `url` supplied for the path (%s) repository does not exist", p.url)
	}
	loader := &pkg.Loader{LoadOptions: true}
	for _, u := range matches {
		dir := u
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(baseDir, dir)
		}
		real, err := filepath.EvalSymlinks(dir)
		if err == nil {
			dir = real
		}
		composerFile := filepath.Join(dir, "composer.json")
		data, err := os.ReadFile(composerFile)
		if err != nil {
			continue
		}
		v, err := jsonfile.Parse(data, composerFile)
		if err != nil {
			return err
		}
		config := phpjson.ToMap(v)
		if config == nil {
			config = phpjson.NewMap()
		}
		config.Delete("published-time")
		dist := phpjson.NewMap()
		dist.Set("type", "path")
		dist.Set("url", u)
		reference := "auto"
		if r := p.options.Str("reference"); r != "" {
			reference = r
		}
		if reference == "none" {
			dist.Set("reference", nil)
		} else if reference == "config" || reference == "auto" {
			h := sha1.Sum(append(data, []byte(PHPSerialize(p.options))...))
			dist.Set("reference", hex.EncodeToString(h[:]))
		}
		config.Set("dist", dist)
		to := phpjson.NewMap()
		for _, k := range p.options.Keys() {
			if k == "symlink" || k == "relative" {
				val, _ := p.options.Get(k)
				to.Set(k, val)
			}
		}
		config.Set("transport-options", to)
		if versions := p.options.Map("versions"); versions != nil && config.Isset("name") {
			if ver, ok := versions.Get(config.Str("name")); ok {
				config.Set("version", ver)
			}
		}
		if reference == "auto" {
			if st, err := os.Stat(filepath.Join(dir, ".git")); err == nil && st.IsDir() {
				if ref, ok := gitHeadCommit(dir); ok {
					dist.Set("reference", ref)
				}
			}
		}
		if !config.Isset("version") {
			if p.GuessVersion != nil {
				if pretty, _, ok := p.GuessVersion(dir, config); ok && pretty != "" {
					config.Set("version", pretty)
				} else {
					config.Set("version", "dev-main")
				}
			} else {
				config.Set("version", "dev-main")
			}
		}
		pk, err := loader.Load(config)
		if err != nil {
			return fmt.Errorf("Failed loading the package in %s: %v", composerFile, err)
		}
		p.Array.AddPackage(pk)
	}
	return nil
}

// PHPSerialize implements PHP serialize() for JSON-like values.
func PHPSerialize(v phpjson.Value) string {
	switch t := v.(type) {
	case nil:
		return "N;"
	case bool:
		if t {
			return "b:1;"
		}
		return "b:0;"
	case int64:
		return "i:" + strconv.FormatInt(t, 10) + ";"
	case int:
		return "i:" + strconv.Itoa(t) + ";"
	case float64:
		return "d:" + phpjson.FormatFloatSerialize(t) + ";"
	case string:
		return "s:" + strconv.Itoa(len(t)) + ":\"" + t + "\";"
	case []phpjson.Value:
		var sb strings.Builder
		sb.WriteString("a:" + strconv.Itoa(len(t)) + ":{")
		for i, e := range t {
			sb.WriteString("i:" + strconv.Itoa(i) + ";")
			sb.WriteString(PHPSerialize(e))
		}
		sb.WriteString("}")
		return sb.String()
	case *phpjson.Map:
		var sb strings.Builder
		sb.WriteString("a:" + strconv.Itoa(t.Len()) + ":{")
		for _, k := range t.Keys() {
			if n, ok := phpjson.IntKey(k); ok {
				sb.WriteString("i:" + strconv.FormatInt(n, 10) + ";")
			} else {
				sb.WriteString("s:" + strconv.Itoa(len(k)) + ":\"" + k + "\";")
			}
			e, _ := t.Get(k)
			sb.WriteString(PHPSerialize(e))
		}
		sb.WriteString("}")
		return sb.String()
	}
	return "N;"
}

// Filter is FilterRepository.
type Filter struct {
	Repo      Repository
	only      *regexp.Regexp
	exclude   *regexp.Regexp
	canonical bool
}

// NewFilter wraps a repository with only/exclude/canonical options.
func NewFilter(r Repository, options *phpjson.Map) (*Filter, error) {
	f := &Filter{Repo: r, canonical: true}
	toNames := func(v phpjson.Value) []string {
		var out []string
		phpjson.Each(v, func(_ string, e phpjson.Value) { s, _ := e.(string); out = append(out, s) })
		return out
	}
	if v, ok := options.Get("only"); ok && v != nil {
		if !phpjson.IsArray(v) {
			return nil, fmt.Errorf(`"only" key for repository %s should be an array`, r.RepoName())
		}
		f.only = namesRegexp(toNames(v))
	}
	if v, ok := options.Get("exclude"); ok && v != nil {
		if !phpjson.IsArray(v) {
			return nil, fmt.Errorf(`"exclude" key for repository %s should be an array`, r.RepoName())
		}
		f.exclude = namesRegexp(toNames(v))
	}
	if f.only != nil && f.exclude != nil {
		return nil, fmt.Errorf(`Only one of "only" and "exclude" can be specified for repository %s`, r.RepoName())
	}
	if v, ok := options.Get("canonical"); ok && v != nil {
		b, isBool := v.(bool)
		if !isBool {
			return nil, fmt.Errorf(`"canonical" key for repository %s should be a boolean`, r.RepoName())
		}
		f.canonical = b
	}
	return f, nil
}

func namesRegexp(names []string) *regexp.Regexp {
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = strings.ReplaceAll(regexp.QuoteMeta(n), `\*`, `.*`)
	}
	return regexp.MustCompile(`(?i)^(?:` + strings.Join(parts, "|") + `)$`)
}

func (f *Filter) isAllowed(name string) bool {
	if f.only == nil && f.exclude == nil {
		return true
	}
	if f.only != nil {
		return f.only.MatchString(name)
	}
	return !f.exclude.MatchString(name)
}

func (f *Filter) RepoName() string { return f.Repo.RepoName() }
func (f *Filter) Kind() Kind       { return f.Repo.Kind() }

func (f *Filter) LoadPackages(names []NameConstraint, acceptable, flags map[string]int, alreadyLoaded map[string]map[string]bool) (LoadResult, error) {
	var allowed []NameConstraint
	for _, nc := range names {
		if f.isAllowed(nc.Name) {
			allowed = append(allowed, nc)
		}
	}
	if len(allowed) == 0 {
		return LoadResult{}, nil
	}
	res, err := f.Repo.LoadPackages(allowed, acceptable, flags, alreadyLoaded)
	if err != nil {
		return res, err
	}
	if !f.canonical {
		res.NamesFound = nil
	}
	return res, nil
}

func (f *Filter) Packages() ([]*pkg.Package, error) {
	ps, err := f.Repo.Packages()
	if err != nil {
		return nil, err
	}
	var out []*pkg.Package
	for _, p := range ps {
		if f.isAllowed(p.Name) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *Filter) FindPackages(name string, c semver.Constraint) ([]*pkg.Package, error) {
	if !f.isAllowed(name) {
		return nil, nil
	}
	return f.Repo.FindPackages(name, c)
}

func (f *Filter) Providers(name string) ([]ProviderInfo, error) {
	ps, err := f.Repo.Providers(name)
	if err != nil {
		return nil, err
	}
	var out []ProviderInfo
	for _, p := range ps {
		if f.isAllowed(p.Name) {
			out = append(out, p)
		}
	}
	return out, nil
}

func gitHeadCommit(dir string) (string, bool) {
	out, err := exec.Command("git", "-C", dir, "rev-list", "-n1", "--format=%H", "HEAD").Output()
	if err != nil {
		return "", false
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "commit ") {
			return l, true
		}
	}
	return "", false
}
