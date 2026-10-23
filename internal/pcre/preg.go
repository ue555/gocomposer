// PHP preg_* calling conventions as used through Composer\Pcre\Preg
// (vendor/composer/pcre/src/Preg.php locked by composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80) and PHP's ext/pcre/php_pcre.c
// (preg_quote, preg_replace replacement parsing, global replacement loop).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package pcre

import (
	"strconv"
	"strings"
	"sync"
)

var cache sync.Map // pattern -> *Regexp or error

// Get compiles a pattern once and caches it (like PHP's regex cache).
func Get(pattern string) (*Regexp, error) {
	if v, ok := cache.Load(pattern); ok {
		if re, ok := v.(*Regexp); ok {
			return re, nil
		}
		return nil, v.(error)
	}
	re, err := Compile(pattern)
	if err != nil {
		cache.Store(pattern, err)
		return nil, err
	}
	cache.Store(pattern, re)
	return re, nil
}

// MustGet is Get that panics on invalid patterns (for fixed patterns).
func MustGet(pattern string) *Regexp {
	re, err := Get(pattern)
	if err != nil {
		panic(err)
	}
	return re
}

// Match is the result of a successful match.
type Match struct {
	re      *Regexp
	subject string
	caps    []int
}

// Group returns group i and whether it participated in the match.
func (m *Match) Group(i int) (string, bool) {
	if i < 0 || 2*i+1 >= len(m.caps) || m.caps[2*i] < 0 {
		return "", false
	}
	return m.subject[m.caps[2*i]:m.caps[2*i+1]], true
}

// Named returns the named group and whether it participated in the match.
func (m *Match) Named(name string) (string, bool) {
	for i, n := range m.re.names {
		if n == name && n != "" {
			return m.Group(i)
		}
	}
	return "", false
}

// Str returns group i or "" when unset.
func (m *Match) Str(i int) string { s, _ := m.Group(i); return s }

// N returns the named group or "" when unset.
func (m *Match) N(name string) string { s, _ := m.Named(name); return s }

// Offsets returns the byte offsets of group i (-1, -1 when unset).
func (m *Match) Offsets(i int) (int, int) { return m.caps[2*i], m.caps[2*i+1] }

// count is PCRE2's return code: 1 + the highest group number that is set.
func (m *Match) count() int {
	c := 1
	for i := 1; 2*i < len(m.caps); i++ {
		if m.caps[2*i] >= 0 {
			c = i + 1
		}
	}
	return c
}

// Entry is one element of a PHP match array.
type Entry struct {
	Key   string
	Value *string // nil = PHP null
}

// PHPArray returns the $matches array PHP builds. unmatchedAsNull mirrors
// PREG_UNMATCHED_AS_NULL: unset groups become null and trailing unset groups
// are included; without it unset groups are "" and trailing ones omitted.
func (m *Match) PHPArray(unmatchedAsNull bool) []Entry {
	var out []Entry
	count := m.count()
	n := m.re.ngroups + 1
	for i := 0; i < n; i++ {
		if i >= count && !unmatchedAsNull {
			break
		}
		var v *string
		if s, ok := m.Group(i); ok {
			v = &s
		} else if !unmatchedAsNull {
			e := ""
			v = &e
		}
		if name := m.re.names[i]; name != "" {
			out = append(out, Entry{Key: name, Value: v})
		}
		out = append(out, Entry{Key: strconv.Itoa(i), Value: v})
	}
	return out
}

// FindAt returns the first match at or after offset, or nil.
func (re *Regexp) FindAt(subject string, offset int) (*Match, error) {
	caps, err := re.exec(subject, offset, execOpts{})
	if err != nil || caps == nil {
		return nil, err
	}
	return &Match{re: re, subject: subject, caps: append([]int(nil), caps...)}, nil
}

// Find returns the first match or nil.
func (re *Regexp) Find(subject string) (*Match, error) { return re.FindAt(subject, 0) }

// IsMatch is Preg::isMatch: it returns the match (nil if none). Errors are
// compile errors or ErrBacktrackLimit/ErrRecursionLimit (PcreException).
func IsMatch(pattern, subject string) (*Match, error) {
	re, err := Get(pattern)
	if err != nil {
		return nil, err
	}
	return re.Find(subject)
}

// ReplaceFunc runs the preg_replace / preg_replace_callback loop. fn returns
// the replacement for a match. limit < 0 means no limit.
func (re *Regexp) ReplaceFunc(subject string, limit int, fn func(*Match) string) (string, int, error) {
	var b strings.Builder
	start := 0  // search offset
	copied := 0 // subject bytes already copied
	count := 0
	opts := execOpts{}
	for limit != 0 {
		caps, err := re.exec(subject, start, opts)
		if err != nil {
			return "", 0, err
		}
		if caps == nil {
			if opts.notEmptyAtStart && start < len(subject) {
				// after an empty match, advance one byte and search again
				start++
				opts = execOpts{}
				continue
			}
			break
		}
		m := &Match{re: re, subject: subject, caps: append([]int(nil), caps...)}
		b.WriteString(subject[copied:caps[0]])
		b.WriteString(fn(m))
		copied = caps[1]
		count++
		if limit > 0 {
			limit--
		}
		start = caps[1]
		if caps[0] == caps[1] {
			opts = execOpts{notEmptyAtStart: true, anchored: true}
		} else {
			opts = execOpts{}
		}
		if start > len(subject) {
			break
		}
	}
	b.WriteString(subject[copied:])
	return b.String(), count, nil
}

// Expand applies a preg_replace replacement string to a match ($n, ${n},
// \n back references; "\\" and "\$" escapes), like php_pcre_replace_impl.
func (m *Match) Expand(replacement string) string {
	var out []byte
	count := m.count()
	walkLast := byte(0)
	r := replacement
	for w := 0; w < len(r); {
		c := r[w]
		if c == '\\' || c == '$' {
			if walkLast == '\\' {
				out[len(out)-1] = c
				w++
				walkLast = 0
				continue
			}
			if ref, next, ok := getBackref(r, w); ok {
				if ref < count {
					out = append(out, m.Str(ref)...)
				}
				w = next
				continue
			}
		}
		out = append(out, c)
		w++
		walkLast = c
	}
	return string(out)
}

// getBackref is php_pcre.c preg_get_backref.
func getBackref(r string, w int) (int, int, bool) {
	if w+1 >= len(r) {
		return 0, 0, false
	}
	inBrace := false
	if r[w] == '$' && r[w+1] == '{' {
		inBrace = true
		w++
	}
	w++
	if w >= len(r) || !isDigit(r[w]) {
		return 0, 0, false
	}
	ref := int(r[w] - '0')
	w++
	if w < len(r) && isDigit(r[w]) {
		ref = ref*10 + int(r[w]-'0')
		w++
	}
	if inBrace {
		if w >= len(r) || r[w] != '}' {
			return 0, 0, false
		}
		w++
	}
	return ref, w, true
}

// Replace is preg_replace($pattern, $replacement, $subject, $limit) through
// Preg::replace.
func Replace(pattern, replacement, subject string, limit int) (string, int, error) {
	re, err := Get(pattern)
	if err != nil {
		return "", 0, err
	}
	return re.ReplaceFunc(subject, limit, func(m *Match) string { return m.Expand(replacement) })
}

// ReplaceAll is preg_replace with array patterns and replacements: each
// pattern is applied in turn to the result of the previous one.
func ReplaceAll(patterns, replacements []string, subject string) (string, error) {
	for i, p := range patterns {
		rep := ""
		if i < len(replacements) {
			rep = replacements[i]
		}
		var err error
		subject, _, err = Replace(p, rep, subject, -1)
		if err != nil {
			return "", err
		}
	}
	return subject, nil
}

// ReplaceCallback is Preg::replaceCallback without flags; the callback sees
// unset groups as "" (Match.PHPArray(false)).
func ReplaceCallback(pattern string, fn func(*Match) string, subject string, limit int) (string, int, error) {
	re, err := Get(pattern)
	if err != nil {
		return "", 0, err
	}
	return re.ReplaceFunc(subject, limit, fn)
}

// Quote is preg_quote($s, $delim); delim may be "" for none.
func Quote(s string, delim string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '.', '\\', '+', '*', '?', '[', '^', ']', '$', '(', ')', '{', '}', '=', '!', '>', '<', '|', ':', '-', '#':
			b.WriteByte('\\')
			b.WriteByte(c)
		case 0:
			b.WriteString(`\000`)
		default:
			if delim != "" && c == delim[0] {
				b.WriteByte('\\')
			}
			b.WriteByte(c)
		}
	}
	return b.String()
}

// AddCSlashes is addcslashes($s, $chars) for a plain character list without
// ranges.
func AddCSlashes(s, chars string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if strings.IndexByte(chars, c) >= 0 {
			b.WriteByte('\\')
			switch {
			case c == '\a':
				b.WriteByte('a')
			case c == '\b':
				b.WriteByte('b')
			case c == '\t':
				b.WriteByte('t')
			case c == '\n':
				b.WriteByte('n')
			case c == '\v':
				b.WriteByte('v')
			case c == '\f':
				b.WriteByte('f')
			case c == '\r':
				b.WriteByte('r')
			case c < 32 || c > 126:
				b.WriteString(strconv.FormatInt(int64(c)+01000, 8)[1:])
			default:
				b.WriteByte(c)
			}
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
