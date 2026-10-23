// Package autoload ports Composer\Autoload\AutoloadGenerator and
// composer/class-map-generator 1.8.0 (the version locked by Composer
// f83588392df43d4fd027d84d26de459946ef4e80). php_strip_whitespace() is
// emulated by a small PHP lexer in stripWhitespace.
package autoload

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// latin1 maps each byte to the rune with the same value so Go regexps can
// emulate PCRE's byte-oriented character classes like [\x7f-\xff].
func latin1(b []byte) string {
	var sb strings.Builder
	sb.Grow(len(b))
	for _, c := range b {
		sb.WriteRune(rune(c))
	}
	return sb.String()
}

func unlatin1(s string) string {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		out = append(out, byte(r))
	}
	return string(out)
}

var (
	reQuickTypes = regexp.MustCompile(`(?i)\b(?:class|interface|trait|enum)\s`)
	reMain       = regexp.MustCompile(`(?i)(?:\b(?P<type>class|interface|trait|enum)\s+(?P<name>[a-zA-Z_\x{7f}-\x{ff}:][a-zA-Z0-9_\x{7f}-\x{ff}:\-]*)|\b(?P<ns>namespace)(?P<nsname>\s+[a-zA-Z_\x{7f}-\x{ff}][a-zA-Z0-9_\x{7f}-\x{ff}]*(?:\s*\\\s*[a-zA-Z_\x{7f}-\x{ff}][a-zA-Z0-9_\x{7f}-\x{ff}]*)*)?\s*[\{;])`)
	reHeredoc    = regexp.MustCompile(`^<<<[ \t]*(['"]?)([a-zA-Z_\x{80}-\x{ff}][a-zA-Z0-9_\x{80}-\x{ff}]*)(['"]?)(?:\r\n|\n|\r)`)
	typeConfig   = map[byte]typeCfg{}
	rejectChars  string
)

type typeCfg struct {
	name    string
	pattern *regexp.Regexp
}

func init() {
	for _, t := range []string{"class", "interface", "trait", "enum"} {
		typeConfig[t[0]] = typeCfg{t, regexp.MustCompile(`(?is)^.\b` + t + `\s+[a-zA-Z_\x{7f}-\x{ff}:][a-zA-Z0-9_\x{7f}-\x{ff}:\-]*`)}
	}
	rejectChars = `?"'</cite`
}

// FindClasses is PhpFileParser::findClasses.
func FindClasses(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(`File at "%s" does not exist, check your classmap definitions`, path)
		}
		return nil, fmt.Errorf(`File at "%s" is not readable, check its permissions`, path)
	}
	stripped := stripWhitespace(raw)
	if len(stripped) == 0 {
		if strings.TrimSpace(string(raw)) == "" {
			return nil, nil
		}
		return nil, fmt.Errorf(`File at "%s" could not be parsed as PHP, it may be binary or corrupted`, path)
	}
	contents := latin1(stripped)
	matches := reQuickTypes.FindAllStringIndex(contents, -1)
	if len(matches) == 0 {
		return nil, nil
	}
	contents = cleanPHP(contents, len(matches))
	var classes []string
	namespace := ""
	off := 0
	for off <= len(contents) {
		loc := reMain.FindStringSubmatchIndex(contents[off:])
		if loc == nil {
			break
		}
		start := off + loc[0]
		if start > 0 {
			prev := contents[start-1]
			if prev == '\\' || prev == '$' || prev == ':' || prev == '>' {
				off = start + 1
				continue
			}
		}
		get := func(i int) string {
			if loc[2*i] < 0 {
				return ""
			}
			return contents[off+loc[2*i] : off+loc[2*i+1]]
		}
		if ns := get(3); ns != "" {
			namespace = strings.NewReplacer(" ", "", "\t", "", "\r", "", "\n", "").Replace(get(4)) + `\`
		} else {
			name := get(2)
			typ := strings.ToLower(get(1))
			if name != "extends" && name != "implements" {
				if name[0] == ':' {
					name = "xhp" + strings.NewReplacer("-", "_", ":", "__").Replace(name)[1:]
				} else if typ == "enum" {
					if i := strings.LastIndexByte(name, ':'); i >= 0 {
						name = name[:i]
					}
				}
				classes = append(classes, unlatin1(strings.TrimLeft(namespace+name, `\`)))
			}
		}
		off = off + loc[1]
		if loc[1] == loc[0] {
			off++
		}
	}
	return classes, nil
}

// cleanPHP is PhpFileCleaner::clean on latin1-mapped contents.
func cleanPHP(contents string, maxMatches int) string {
	var clean strings.Builder
	n := len(contents)
	i := 0
	peek := func(c byte) bool { return i+1 < n && contents[i+1] == c }
	for i < n {
		// skipToPhp
		for i < n {
			if contents[i] == '<' && peek('?') {
				i += 2
				break
			}
			i++
		}
		clean.WriteString("<?")
	inner:
		for i < n {
			c := contents[i]
			if c == '?' && peek('>') {
				clean.WriteString("?>")
				i += 2
				break inner
			}
			if c == '"' || c == '\'' {
				i = skipString(contents, i, c)
				clean.WriteString("null")
				continue
			}
			if c == '<' && peek('<') {
				if m := reHeredoc.FindStringSubmatch(contents[i:]); m != nil && m[1] == m[3] {
					i += len(m[0])
					i = skipHeredoc(contents, i, m[2])
					clean.WriteString("null")
					continue
				}
			}
			if c == '/' {
				if peek('/') {
					i += strcspn(contents[i:], "\r\n")
					continue
				}
				if peek('*') {
					i += 2
					for i < n {
						i += strcspn(contents[i:], "*")
						if i+1 < n && contents[i+1] == '/' {
							i += 2
							break
						}
						i++
					}
					continue
				}
			}
			if maxMatches == 1 {
				if tc, ok := typeConfig[c]; ok {
					if strings.HasPrefix(contents[i:], tc.name) && i > 0 {
						if m := tc.pattern.FindString(contents[i-1:]); m != "" {
							prev := contents[i-1]
							if prev != '$' && prev != ':' && prev != '>' {
								return clean.String() + m
							}
						}
					}
				}
			}
			i++
			skip := strcspn(contents[i:], rejectChars)
			clean.WriteByte(c)
			if skip > 0 {
				clean.WriteString(contents[i : i+skip])
				i += skip
			}
		}
	}
	return clean.String()
}

func strcspn(s, reject string) int {
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(reject, s[i]) >= 0 {
			return i
		}
	}
	return len(s)
}

func skipString(s string, i int, delim byte) int {
	n := len(s)
	i++
	for i < n {
		i += strcspn(s[i:], `\`+string(delim))
		if i >= n {
			break
		}
		if s[i] == '\\' && i+1 < n && (s[i+1] == '\\' || s[i+1] == delim) {
			i += 2
			continue
		}
		if s[i] == delim {
			return i + 1
		}
		i++
	}
	return i
}

func isIdentByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c >= 0x80
}

func skipHeredoc(s string, i int, delim string) int {
	n := len(s)
	for i < n {
		switch s[i] {
		case '\t', ' ':
			i++
			continue
		case delim[0]:
			if strings.HasPrefix(s[i:], delim) {
				end := i + len(delim)
				if end >= n || !isIdentByte(s[end]) {
					return end
				}
			}
		}
		i += strcspn(s[i:], "\r\n")
		for i < n && (s[i] == '\r' || s[i] == '\n') {
			i++
		}
	}
	return i
}

// ClassMap is Composer\ClassMapGenerator\ClassMap.
type ClassMap struct {
	Map            map[string]string
	order          []string
	ambiguous      map[string][]string
	ambiguousOrder []string
	psrViolations  map[string][]string
	psrOrder       []string
}

func NewClassMap() *ClassMap {
	return &ClassMap{Map: map[string]string{}, ambiguous: map[string][]string{}, psrViolations: map[string][]string{}}
}

func (m *ClassMap) AddClass(class, path string) {
	key := strings.ReplaceAll(path, `\`, "/")
	if _, ok := m.psrViolations[key]; ok {
		delete(m.psrViolations, key)
		for i, k := range m.psrOrder {
			if k == key {
				m.psrOrder = append(m.psrOrder[:i:i], m.psrOrder[i+1:]...)
				break
			}
		}
	}
	if _, ok := m.Map[class]; !ok {
		m.order = append(m.order, class)
	}
	m.Map[class] = path
}

func (m *ClassMap) addAmbiguous(class, path string) {
	if _, ok := m.ambiguous[class]; !ok {
		m.ambiguousOrder = append(m.ambiguousOrder, class)
	}
	m.ambiguous[class] = append(m.ambiguous[class], path)
}

func (m *ClassMap) addPsrViolation(warning, path string) {
	key := strings.TrimRight(strings.ReplaceAll(path, `\`, "/"), "/")
	if _, ok := m.psrViolations[key]; !ok {
		m.psrOrder = append(m.psrOrder, key)
	}
	m.psrViolations[key] = append(m.psrViolations[key], warning)
}

// ClearPsrViolationsByPath removes violations under a path prefix.
func (m *ClassMap) ClearPsrViolationsByPath(prefix string) {
	prefix = strings.TrimRight(strings.ReplaceAll(prefix, `\`, "/"), "/")
	var keep []string
	for _, k := range m.psrOrder {
		if k == prefix || strings.HasPrefix(k, prefix+"/") {
			delete(m.psrViolations, k)
			continue
		}
		keep = append(keep, k)
	}
	m.psrOrder = keep
}

// PsrViolations returns the warnings in insertion order.
func (m *ClassMap) PsrViolations() []string {
	var out []string
	for _, k := range m.psrOrder {
		out = append(out, m.psrViolations[k]...)
	}
	return out
}

var reDuplicatesFilter = regexp.MustCompile(`(?i)(?:^|/)(test|fixture|example|stub)s?/`)

// AmbiguousClasses is getAmbiguousClasses() with the default filter (or none).
func (m *ClassMap) AmbiguousClasses(filter bool) ([]string, map[string][]string) {
	out := map[string][]string{}
	var order []string
	for _, c := range m.ambiguousOrder {
		var paths []string
		for _, p := range m.ambiguous[c] {
			if filter && reDuplicatesFilter.MatchString(strings.ReplaceAll(p, `\`, "/")) {
				continue
			}
			paths = append(paths, p)
		}
		if len(paths) > 0 {
			out[c] = paths
			order = append(order, c)
		}
	}
	return order, out
}

func foldCase(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

var reRootPrefix = regexp.MustCompile(`(?i)^(?:[0-9a-z]{2,}:(?://(?:[a-z]:)?)?|[a-z]:)`)

// AmbiguousFolders is getAmbiguousFolders().
func (m *ClassMap) AmbiguousFolders(filter bool) [][]string {
	var paths []string
	for _, c := range m.order {
		paths = append(paths, strings.ReplaceAll(m.Map[c], `\`, "/"))
	}
	for _, c := range m.ambiguousOrder {
		for _, p := range m.ambiguous[c] {
			paths = append(paths, strings.ReplaceAll(p, `\`, "/"))
		}
	}
	sort.Strings(paths)
	first := map[string]string{}
	variants := map[string]map[string]string{}
	variantOrder := map[string][]string{}
	var variantKeys []string
	previous := ""
	for _, p := range paths {
		common := 0
		for common < len(p) && common < len(previous) && p[common] == previous[common] {
			common++
		}
		skipUntil := 0
		if common > 0 {
			if idx := strings.LastIndexByte(p[:common], '/'); idx >= 0 {
				skipUntil = idx
			}
		}
		previous = p
		folded := foldCase(p)
		rootLen := 0
		if m := reRootPrefix.FindString(p); m != "" {
			rootLen = len(m)
		}
		offset := 0
		for _, name := range strings.Split(p, "/") {
			end := offset + len(name)
			if name != "" && end > rootLen && end > skipUntil {
				fp := folded[:end]
				if f, ok := first[fp]; !ok {
					first[fp] = p[:offset] + name
				} else if f[offset:] != name {
					if _, ok := variants[fp]; !ok {
						variants[fp] = map[string]string{}
						variantKeys = append(variantKeys, fp)
					}
					if _, ok := variants[fp][name]; !ok {
						variants[fp][name] = p[:offset] + name
						variantOrder[fp] = append(variantOrder[fp], name)
					}
				}
			}
			offset = end + 1
		}
	}
	if len(variants) == 0 {
		return nil
	}
	var relevant map[string]bool
	if filter {
		relevant = map[string]bool{}
		for _, p := range paths {
			if reDuplicatesFilter.MatchString(p) {
				continue
			}
			folded := foldCase(p)
			off := 0
			for {
				sep := strings.IndexByte(folded[off:], '/')
				if sep < 0 {
					break
				}
				relevant[folded[:off+sep]] = true
				off += sep + 1
			}
			relevant[folded] = true
		}
	}
	var keys []string
	for _, k := range variantKeys {
		if relevant == nil || relevant[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var out [][]string
	for _, k := range keys {
		group := []string{first[k]}
		for _, n := range variantOrder[k] {
			group = append(group, variants[k][n])
		}
		out = append(out, group)
	}
	return out
}

// SortedClasses returns class names in ksort order.
func (m *ClassMap) SortedClasses() []string {
	keys := make([]string, 0, len(m.Map))
	for k := range m.Map {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Generator is ClassMapGenerator with avoidDuplicateScans().
type Generator struct {
	extensions []string
	scanned    map[string]bool
	ClassMap   *ClassMap
	Cwd        string
}

func NewGenerator(cwd string) *Generator {
	return &Generator{extensions: []string{"php", "inc", "hh"}, scanned: map[string]bool{}, ClassMap: NewClassMap(), Cwd: cwd}
}

var vcsDirs = map[string]bool{".svn": true, "_svn": true, "CVS": true, "_darcs": true, ".arch-params": true, ".monotone": true, ".bzr": true, ".git": true, ".hg": true}

// findFiles emulates Finder::create()->files()->followLinks()->name(...)->in($path)
// (VCS and dot files ignored); entries are visited in name order.
func (g *Generator) findFiles(root string) ([]string, error) {
	var out []string
	var walk func(dir, rel string, depth int) error
	visited := map[string]bool{}
	walk = func(dir, rel string, depth int) error {
		real, err := filepath.EvalSymlinks(dir)
		if err == nil {
			if visited[real] {
				return nil
			}
			visited[real] = true
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") || vcsDirs[name] {
				continue
			}
			p := dir + "/" + name
			st, err := os.Stat(p)
			if err != nil {
				continue
			}
			if st.IsDir() {
				if err := walk(p, rel+name+"/", depth+1); err != nil {
					return err
				}
				continue
			}
			ext := ""
			if i := strings.LastIndexByte(name, '.'); i >= 0 {
				ext = name[i+1:]
			}
			for _, x := range g.extensions {
				if ext == x {
					out = append(out, p)
					break
				}
			}
		}
		return nil
	}
	return out, walk(strings.TrimRight(root, "/"), "", 0)
}

// ScanPaths is ClassMapGenerator::scanPaths for a single path.
func (g *Generator) ScanPaths(path string, excluded *regexp.Regexp, autoloadType, namespace string) error {
	var files []string
	basePath := path
	if st, err := os.Stat(path); err == nil && !st.IsDir() {
		files = []string{path}
	} else if err == nil && st.IsDir() {
		f, err := g.findFiles(path)
		if err != nil {
			return err
		}
		files = f
	} else if strings.Contains(path, "*") {
		matches, _ := filepath.Glob(path)
		sort.Strings(matches)
		for _, m := range matches {
			f, err := g.findFiles(m)
			if err == nil {
				files = append(files, f...)
			}
		}
	} else {
		return fmt.Errorf(`Could not scan for classes inside "%s" which does not appear to be a file nor a folder`, path)
	}
	cwd := g.Cwd
	for _, filePath := range files {
		ext := ""
		if i := strings.LastIndexByte(filePath, '.'); i >= 0 && !strings.Contains(filePath[i:], "/") {
			ext = filePath[i+1:]
		}
		ok := false
		for _, x := range g.extensions {
			if ext == x {
				ok = true
			}
		}
		if !ok {
			continue
		}
		if !strings.HasPrefix(filePath, "/") {
			filePath = normalizePath(cwd + "/" + filePath)
		} else {
			filePath = reMultiSlash.ReplaceAllString(filePath, "${1}/")
		}
		realPath, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			return fmt.Errorf("realpath of %s failed to resolve, got false", filePath)
		}
		if abs, err := filepath.Abs(realPath); err == nil {
			realPath = abs
		}
		if g.scanned[realPath] {
			continue
		}
		if excluded != nil && (excluded.MatchString(strings.ReplaceAll(realPath, `\`, "/")) || excluded.MatchString(strings.ReplaceAll(filePath, `\`, "/"))) {
			continue
		}
		classes, err := FindClasses(filePath)
		if err != nil {
			return err
		}
		if autoloadType != "classmap" {
			classes = g.filterByNamespace(classes, filePath, namespace, autoloadType, basePath)
			if len(classes) > 0 {
				g.scanned[realPath] = true
			}
		} else {
			g.scanned[realPath] = true
		}
		for _, c := range classes {
			if cur, ok := g.ClassMap.Map[c]; !ok {
				g.ClassMap.AddClass(c, filePath)
			} else if filePath != cur {
				g.ClassMap.addAmbiguous(c, filePath)
			}
		}
	}
	return nil
}

var reMultiSlash = regexp.MustCompile(`([^:/])[\\/]{2,}`)

func (g *Generator) filterByNamespace(classes []string, filePath, baseNamespace, typ, basePath string) []string {
	var valid, rejected []string
	realSubPath := ""
	if len(filePath) > len(basePath)+1 {
		realSubPath = filePath[len(basePath)+1:]
	}
	if i := strings.LastIndexByte(realSubPath, '.'); i >= 0 {
		realSubPath = realSubPath[:i]
	}
	for _, class := range classes {
		var subPath string
		if typ == "psr-0" {
			if baseNamespace != "" && !strings.HasPrefix(class, baseNamespace) {
				rejected = append(rejected, class)
				continue
			}
			if i := strings.LastIndexByte(class, '\\'); i >= 0 {
				ns := class[:i+1]
				cn := class[i+1:]
				subPath = strings.ReplaceAll(ns, `\`, "/") + strings.ReplaceAll(cn, "_", "/")
			} else {
				subPath = strings.ReplaceAll(class, "_", "/")
			}
		} else {
			sub := class
			if baseNamespace != "" {
				if len(class) >= len(baseNamespace) {
					sub = class[len(baseNamespace):]
				} else {
					sub = ""
				}
			}
			subPath = strings.ReplaceAll(sub, `\`, "/")
		}
		if subPath == realSubPath {
			valid = append(valid, class)
		} else {
			rejected = append(rejected, class)
		}
	}
	if len(valid) == 0 {
		cwd := normalizePath(g.Cwd)
		shortPath := replacePrefixOnce(normalizePath(filePath), cwd, ".")
		shortBase := replacePrefixOnce(normalizePath(basePath), cwd, ".")
		for _, c := range rejected {
			g.ClassMap.addPsrViolation(fmt.Sprintf("Class %s located in %s does not comply with %s autoloading standard (rule: %s => %s). Skipping.", c, shortPath, typ, baseNamespace, shortBase), filePath)
		}
		return nil
	}
	return valid
}

func replacePrefixOnce(s, prefix, repl string) string {
	if strings.HasPrefix(s, prefix) {
		return repl + s[len(prefix):]
	}
	return s
}

// normalizePath is ClassMapGenerator::normalizePath (same as Filesystem's).
func normalizePath(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	absolute := ""
	if strings.HasPrefix(p, "/") {
		absolute = "/"
		p = p[1:]
	}
	var parts []string
	up := false
	for _, chunk := range strings.Split(p, "/") {
		if chunk == ".." && (absolute != "" || up) {
			if len(parts) > 0 {
				parts = parts[:len(parts)-1]
			}
			up = !(len(parts) == 0 || parts[len(parts)-1] == "..")
		} else if chunk != "." && chunk != "" {
			parts = append(parts, chunk)
			up = chunk != ".."
		}
	}
	return absolute + strings.Join(parts, "/")
}
