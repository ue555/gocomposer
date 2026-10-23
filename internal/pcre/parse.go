// Package pcre is a small backtracking regular expression engine implementing
// the PCRE2 (10.44, as bundled with PHP 8.4) subset used by Composer's
// JsonManipulator, JsonConfigSource and JsonFile, together with the PHP
// preg_* / Composer\Pcre\Preg calling conventions (delimited patterns with
// modifiers, preg_quote, preg_replace replacement back references,
// preg_replace_callback match arrays, PREG_UNMATCHED_AS_NULL).
//
// Ported behaviour follows src/Composer/Json/JsonManipulator.php and
// vendor/composer/pcre/src/Preg.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80 and PHP's ext/pcre/php_pcre.c.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// The engine is byte based (no /u). Supported syntax: literals and escapes,
// character classes, . ^ $ \A \z \Z \b \B, \d \D \s \S \w \W \h \H \v \V \N,
// groups (capturing, named (?<n>..)/(?P<n>..)/(?'n'..), non-capturing,
// atomic (?>..), look-ahead (?=..)/(?!..), inline flags), (?(DEFINE)..),
// subroutine calls (?&name) (?P>name) (?R) (?n) with PCRE2 >= 10.30
// semantics (not atomic, captures revert after the call), greedy/lazy/
// possessive quantifiers and the modifiers i m s x D U A. Anything else
// (back references, look-behind, \G, \K, \R, \p, /u, ...) is rejected at
// compile time instead of being silently mis-handled.
package pcre

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CompileError reports a pattern the engine cannot (or PHP would not) compile.
type CompileError struct {
	Pattern     string
	Msg         string
	Unsupported bool
}

func (e *CompileError) Error() string {
	return "pcre: " + e.Msg + " in pattern " + strconv.Quote(e.Pattern)
}

// Unwrap makes errors.Is(err, ErrUnsupported) work for unsupported syntax.
func (e *CompileError) Unwrap() error {
	if e.Unsupported {
		return ErrUnsupported
	}
	return nil
}

// ErrUnsupported is wrapped by compile errors for valid PCRE syntax that this
// engine does not implement.
var ErrUnsupported = errors.New("unsupported PCRE feature")

type flags struct {
	i, m, s, x, D, U bool
}

type parser struct {
	src     string
	pos     int
	ngroups int
	names   map[string]int
	groups  map[int]*capNode
	calls   []*callNode
	pattern string
}

func (p *parser) errorf(format string, a ...any) error {
	return &CompileError{Pattern: p.pattern, Msg: fmt.Sprintf(format, a...)}
}

func (p *parser) unsupported(what string) error {
	return &CompileError{Pattern: p.pattern, Msg: fmt.Sprintf("%s: %s", ErrUnsupported, what), Unsupported: true}
}

// splitDelimited separates a PHP pattern into body and modifiers following
// php_pcre.c pcre_get_compiled_regex_cache_ex.
func splitDelimited(pattern string) (body, mods string, err error) {
	p := strings.TrimLeft(pattern, " \t\n\r\v\f")
	if p == "" {
		return "", "", &CompileError{Pattern: pattern, Msg: "Empty regular expression"}
	}
	d := p[0]
	if isAlnum(d) || d == '\\' || d == 0 {
		return "", "", &CompileError{Pattern: pattern, Msg: "Delimiter must not be alphanumeric, backslash, or NUL"}
	}
	end := d
	switch d {
	case '(':
		end = ')'
	case '[':
		end = ']'
	case '{':
		end = '}'
	case '<':
		end = '>'
	}
	i := 1
	if end == d {
		for i < len(p) {
			if p[i] == '\\' && i+1 < len(p) {
				i += 2
				continue
			}
			if p[i] == d {
				break
			}
			i++
		}
	} else {
		brackets := 1
		for i < len(p) {
			if p[i] == '\\' && i+1 < len(p) {
				i += 2
				continue
			}
			if p[i] == end {
				brackets--
				if brackets <= 0 {
					break
				}
			} else if p[i] == d {
				brackets++
			}
			i++
		}
	}
	if i >= len(p) {
		return "", "", &CompileError{Pattern: pattern, Msg: "No ending delimiter found"}
	}
	return p[1:i], p[i+1:], nil
}

func isAlnum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// Compile compiles a delimited PHP pattern such as "{^a+}sx".
func Compile(pattern string) (*Regexp, error) {
	body, mods, err := splitDelimited(pattern)
	if err != nil {
		return nil, err
	}
	var fl flags
	anchoredMod := false
	for i := 0; i < len(mods); i++ {
		switch mods[i] {
		case 'i':
			fl.i = true
		case 'm':
			fl.m = true
		case 's':
			fl.s = true
		case 'x':
			fl.x = true
		case 'D':
			fl.D = true
		case 'U':
			fl.U = true
		case 'A':
			anchoredMod = true
		case 'S', ' ', '\n', '\r':
		case 'u':
			return nil, &CompileError{Pattern: pattern, Msg: ErrUnsupported.Error() + ": /u modifier", Unsupported: true}
		default:
			return nil, &CompileError{Pattern: pattern, Msg: fmt.Sprintf("Unknown modifier '%c'", mods[i])}
		}
	}
	p := &parser{src: body, names: map[string]int{}, groups: map[int]*capNode{}, pattern: pattern}
	n, err := p.parseAlt(&fl, 0)
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.src) {
		return nil, p.errorf("unmatched closing parenthesis at offset %d", p.pos)
	}
	for _, c := range p.calls {
		if c.name != "" {
			idx, ok := p.names[c.name]
			if !ok {
				return nil, p.errorf("reference to non-existent subpattern %q", c.name)
			}
			c.num = idx
		}
		if c.num == 0 {
			c.target = n
			continue
		}
		g, ok := p.groups[c.num]
		if !ok {
			return nil, p.errorf("reference to non-existent subpattern %d", c.num)
		}
		c.target = g.body
	}
	names := make([]string, p.ngroups+1)
	for name, idx := range p.names {
		names[idx] = name
	}
	re := &Regexp{pattern: pattern, prog: n, ngroups: p.ngroups, names: names, anchored: anchoredMod || startsAnchored(n)}
	re.forceAnchored = anchoredMod
	return re, nil
}

// startsAnchored reports whether every match must start at the subject start
// (leading ^ without /m or \A, ignoring DEFINE groups).
func startsAnchored(n node) bool {
	switch t := n.(type) {
	case *seqNode:
		for _, it := range t.items {
			switch it.(type) {
			case defineNode:
				continue
			}
			return startsAnchored(it)
		}
		return false
	case assertNode:
		return t.kind == aStartSubject || t.kind == aCaret && !t.multiline
	case *capNode:
		return startsAnchored(t.body)
	case *altNode:
		for _, a := range t.alts {
			if !startsAnchored(a) {
				return false
			}
		}
		return len(t.alts) > 0
	}
	return false
}

func (p *parser) more() bool { return p.pos < len(p.src) }

// skipQuantGap skips what may separate an item from its quantifier:
// (?#...) comments always, white space and # comments with /x.
func (p *parser) skipQuantGap(fl *flags) {
	for {
		p.skipExtended(fl)
		if strings.HasPrefix(p.src[p.pos:], "(?#") {
			if end := strings.IndexByte(p.src[p.pos:], ')'); end >= 0 {
				p.pos += end + 1
				continue
			}
		}
		return
	}
}

func (p *parser) skipExtended(fl *flags) {
	if !fl.x {
		return
	}
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			p.pos++
			continue
		}
		if c == '#' {
			for p.pos < len(p.src) && p.src[p.pos] != '\n' {
				p.pos++
			}
			continue
		}
		return
	}
}

// parseAlt parses alternatives until ')' or the end. The flags pointer is
// shared with the enclosing group so inline (?i) changes stay scoped.
func (p *parser) parseAlt(fl *flags, depth int) (node, error) {
	if depth > 250 {
		return nil, p.errorf("parentheses are too deeply nested")
	}
	var alts []node
	for {
		seq, err := p.parseSeq(fl, depth)
		if err != nil {
			return nil, err
		}
		alts = append(alts, seq)
		if p.more() && p.src[p.pos] == '|' {
			p.pos++
			continue
		}
		break
	}
	if len(alts) == 1 {
		return alts[0], nil
	}
	return &altNode{alts: alts}, nil
}

func (p *parser) parseSeq(fl *flags, depth int) (node, error) {
	seq := &seqNode{}
	for {
		p.skipExtended(fl)
		if !p.more() || p.src[p.pos] == '|' || p.src[p.pos] == ')' {
			break
		}
		start := p.pos
		atom, quantifiable, err := p.parseAtom(fl, depth)
		if err != nil {
			return nil, err
		}
		if atom == nil { // inline option setting or comment
			continue
		}
		p.skipQuantGap(fl)
		q, ok, err := p.parseQuantifier(fl)
		if err != nil {
			return nil, err
		}
		if ok {
			if !quantifiable {
				return nil, p.errorf("quantifier does not follow a repeatable item at offset %d", start)
			}
			atom = makeRepeat(atom, q)
			// a quantifier directly following a quantifier is an error
			p.skipQuantGap(fl)
			if p.more() && (p.src[p.pos] == '*' || p.src[p.pos] == '+' || p.src[p.pos] == '?') {
				return nil, p.errorf("quantifier does not follow a repeatable item at offset %d", p.pos)
			}
		}
		seq.items = append(seq.items, atom)
	}
	seq.items = mergeLiterals(seq.items)
	if len(seq.items) == 1 {
		return seq.items[0], nil
	}
	return seq, nil
}

type quant struct {
	min, max int // max < 0: unbounded
	lazy     bool
	poss     bool
}

func (p *parser) parseQuantifier(fl *flags) (quant, bool, error) {
	if !p.more() {
		return quant{}, false, nil
	}
	var q quant
	switch p.src[p.pos] {
	case '*':
		q = quant{0, -1, false, false}
		p.pos++
	case '+':
		q = quant{1, -1, false, false}
		p.pos++
	case '?':
		q = quant{0, 1, false, false}
		p.pos++
	case '{':
		mn, mx, n, ok := parseBraces(p.src[p.pos:])
		if !ok {
			return quant{}, false, nil
		}
		if mx >= 0 && mn > mx {
			return quant{}, false, p.errorf("numbers out of order in {} quantifier")
		}
		if mn > 65535 || mx > 65535 {
			return quant{}, false, p.errorf("number too big in {} quantifier")
		}
		q = quant{mn, mx, false, false}
		p.pos += n
	default:
		return quant{}, false, nil
	}
	p.skipQuantGap(fl)
	if p.more() {
		switch p.src[p.pos] {
		case '?':
			q.lazy = true
			p.pos++
		case '+':
			q.poss = true
			p.pos++
		}
	}
	if fl.U && !q.poss {
		q.lazy = !q.lazy
	}
	return q, true, nil
}

// parseBraces parses {n}, {n,}, {n,m} and (PCRE2 >= 10.43) {,m}.
func parseBraces(s string) (mn, mx, n int, ok bool) {
	i := 1
	readNum := func() (int, bool) {
		st := i
		for i < len(s) && isDigit(s[i]) {
			i++
		}
		if st == i {
			return 0, false
		}
		v, err := strconv.Atoi(s[st:i])
		if err != nil {
			return 1 << 30, true
		}
		return v, true
	}
	a, hasA := readNum()
	if i < len(s) && s[i] == '}' {
		if !hasA {
			return 0, 0, 0, false
		}
		return a, a, i + 1, true
	}
	if i >= len(s) || s[i] != ',' {
		return 0, 0, 0, false
	}
	i++
	b, hasB := readNum()
	if i >= len(s) || s[i] != '}' {
		return 0, 0, 0, false
	}
	if !hasA && !hasB {
		return 0, 0, 0, false
	}
	if !hasA {
		a = 0
	}
	if !hasB {
		b = -1
	}
	return a, b, i + 1, true
}

func makeRepeat(atom node, q quant) node {
	if q.min == 1 && q.max == 1 {
		return atom
	}
	var set *charSet
	if s, ok := atom.(*setNode); ok {
		set = &s.set
	} else if l, ok := atom.(*litNode); ok && len(l.s) == 1 {
		cs := l.charSet()
		set = &cs
	}
	if set != nil {
		r := repSetNode{set: *set, min: q.min, max: q.max, lazy: q.lazy && !q.poss}
		if q.poss {
			return &possSetNode{r}
		}
		return &r
	}
	r := &repNode{body: atom, min: q.min, max: q.max, lazy: q.lazy && !q.poss}
	if q.poss {
		hc := hasCaptures(atom)
		return &possRepNode{body: atom, min: q.min, max: q.max, hasCaps: hc, generic: &atomicNode{body: r, hasCaps: hc}}
	}
	return r
}

// parseAtom returns the atom, whether a quantifier may follow, or a nil node
// for items that produce nothing (inline flags, comments).
func (p *parser) parseAtom(fl *flags, depth int) (node, bool, error) {
	c := p.src[p.pos]
	switch c {
	case '(':
		return p.parseGroup(fl, depth)
	case '[':
		s, err := p.parseClass(fl)
		if err != nil {
			return nil, false, err
		}
		return &setNode{set: s}, true, nil
	case '.':
		p.pos++
		var s charSet
		s.fill()
		if !fl.s {
			s.del('\n')
		}
		return &setNode{set: s}, true, nil
	case '^':
		p.pos++
		return assertNode{kind: aCaret, multiline: fl.m}, false, nil
	case '$':
		p.pos++
		return assertNode{kind: aDollar, multiline: fl.m, dollarEndOnly: fl.D}, false, nil
	case '\\':
		return p.parseEscape(fl)
	case '*', '+', '?':
		return nil, false, p.errorf("quantifier does not follow a repeatable item at offset %d", p.pos)
	case '{':
		if _, _, _, ok := parseBraces(p.src[p.pos:]); ok {
			return nil, false, p.errorf("quantifier does not follow a repeatable item at offset %d", p.pos)
		}
	}
	p.pos++
	return newLit(string([]byte{c}), fl.i), true, nil
}

func (p *parser) parseEscape(fl *flags) (node, bool, error) {
	if p.pos+1 >= len(p.src) {
		return nil, false, p.errorf("\\ at end of pattern")
	}
	c := p.src[p.pos+1]
	p.pos += 2
	var s charSet
	switch c {
	case 'd', 'D', 's', 'S', 'w', 'W', 'h', 'H', 'v', 'V':
		s.addClassEscape(c)
		return &setNode{set: s}, true, nil
	case 'N':
		s.fill()
		s.del('\n')
		return &setNode{set: s}, true, nil
	case 'A':
		return assertNode{kind: aStartSubject}, false, nil
	case 'z':
		return assertNode{kind: aEndSubject}, false, nil
	case 'Z':
		return assertNode{kind: aEndSubjectNL}, false, nil
	case 'b':
		return assertNode{kind: aWordB}, false, nil
	case 'B':
		return assertNode{kind: aNotWordB}, false, nil
	case 'Q':
		return nil, false, p.unsupported(`\Q..\E quoting`)
	case 'E':
		return nil, false, nil // \E without \Q is ignored
	}
	if b, ok, err := p.charEscape(c); err != nil {
		return nil, false, err
	} else if ok {
		return newLit(string([]byte{b}), fl.i), true, nil
	}
	if c >= '1' && c <= '9' {
		// PCRE2: a decimal number below 10, starting with 8 or 9, or not
		// above the number of groups so far is a back reference; otherwise
		// up to three octal digits are read.
		st := p.pos - 1
		e := st
		for e < len(p.src) && isDigit(p.src[e]) {
			e++
		}
		n, _ := strconv.Atoi(p.src[st:e])
		if n < 10 || c == '8' || c == '9' || n <= p.ngroups {
			return nil, false, p.unsupported("back reference")
		}
		p.pos = st
		return newLit(string([]byte{p.readOctal(3)}), fl.i), true, nil
	}
	if isAlnum(c) {
		switch c {
		case 'G', 'K', 'R', 'X', 'p', 'P', 'g', 'k', 'C':
			return nil, false, p.unsupported(`\` + string([]byte{c}))
		}
		return nil, false, p.errorf("unrecognized character follows \\")
	}
	return newLit(string([]byte{c}), fl.i), true, nil
}

// charEscape decodes single-character escapes shared by classes and atoms.
// p.pos is just after the escape letter c.
func (p *parser) charEscape(c byte) (byte, bool, error) {
	switch c {
	case 'n':
		return '\n', true, nil
	case 'r':
		return '\r', true, nil
	case 't':
		return '\t', true, nil
	case 'f':
		return '\f', true, nil
	case 'e':
		return 0x1b, true, nil
	case 'a':
		return 0x07, true, nil
	case '0':
		v := 0
		for k := 0; k < 2 && p.more() && p.src[p.pos] >= '0' && p.src[p.pos] <= '7'; k++ {
			v = v*8 + int(p.src[p.pos]-'0')
			p.pos++
		}
		return byte(v), true, nil
	case 'o':
		if !p.more() || p.src[p.pos] != '{' {
			return 0, false, p.errorf("missing opening brace after \\o")
		}
		end := strings.IndexByte(p.src[p.pos:], '}')
		if end < 0 {
			return 0, false, p.errorf("missing closing brace")
		}
		v, err := strconv.ParseUint(p.src[p.pos+1:p.pos+end], 8, 32)
		if err != nil || v > 255 {
			return 0, false, p.errorf("invalid octal escape")
		}
		p.pos += end + 1
		return byte(v), true, nil
	case 'x':
		if p.more() && p.src[p.pos] == '{' {
			end := strings.IndexByte(p.src[p.pos:], '}')
			if end < 0 {
				return 0, false, p.errorf("missing closing brace")
			}
			v, err := strconv.ParseUint(p.src[p.pos+1:p.pos+end], 16, 32)
			if err != nil || v > 255 {
				return 0, false, p.errorf("character code point value in \\x{} is too large")
			}
			p.pos += end + 1
			return byte(v), true, nil
		}
		v := 0
		for k := 0; k < 2 && p.more() && isHex(p.src[p.pos]); k++ {
			v = v*16 + hexVal(p.src[p.pos])
			p.pos++
		}
		return byte(v), true, nil
	}
	return 0, false, nil
}

// readOctal reads up to max octal digits at p.pos.
func (p *parser) readOctal(max int) byte {
	v := 0
	for k := 0; k < max && p.more() && p.src[p.pos] >= '0' && p.src[p.pos] <= '7'; k++ {
		v = v*8 + int(p.src[p.pos]-'0')
		p.pos++
	}
	return byte(v)
}

func isHex(c byte) bool {
	return isDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func hexVal(c byte) int {
	switch {
	case isDigit(c):
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	default:
		return int(c-'A') + 10
	}
}

func (p *parser) parseClass(fl *flags) (charSet, error) {
	var s charSet
	p.pos++ // [
	neg := false
	if p.more() && p.src[p.pos] == '^' {
		neg = true
		p.pos++
	}
	first := true
	for {
		if !p.more() {
			return s, p.errorf("missing terminating ] for character class")
		}
		c := p.src[p.pos]
		if c == ']' && !first {
			p.pos++
			break
		}
		first = false
		if c == '[' && p.pos+1 < len(p.src) && (p.src[p.pos+1] == ':' || p.src[p.pos+1] == '.' || p.src[p.pos+1] == '=') {
			if strings.Contains(p.src[p.pos+2:], string([]byte{p.src[p.pos+1]})+"]") {
				return s, p.unsupported("POSIX character class")
			}
		}
		var lo int
		if c == '\\' {
			if p.pos+1 >= len(p.src) {
				return s, p.errorf("\\ at end of pattern")
			}
			e := p.src[p.pos+1]
			p.pos += 2
			switch e {
			case 'd', 'D', 's', 'S', 'w', 'W', 'h', 'H', 'v', 'V':
				s.addClassEscape(e)
				continue
			case 'b':
				lo = 0x08
			case 'Q', 'E':
				return s, p.unsupported(`\Q..\E in a character class`)
			default:
				b, err := p.classChar(e)
				if err != nil {
					return s, err
				}
				lo = int(b)
			}
		} else {
			p.pos++
			lo = int(c)
		}
		// range?
		if p.pos+1 < len(p.src) && p.src[p.pos] == '-' && p.src[p.pos+1] != ']' {
			save := p.pos
			p.pos++
			hc := p.src[p.pos]
			hi := -1
			if hc == '\\' {
				if p.pos+1 >= len(p.src) {
					return s, p.errorf("\\ at end of pattern")
				}
				e := p.src[p.pos+1]
				switch e {
				case 'd', 'D', 's', 'S', 'w', 'W', 'h', 'H', 'v', 'V':
					return s, p.errorf("invalid range in character class")
				case 'b':
					hi = 0x08
					p.pos += 2
				default:
					p.pos += 2
					b, err := p.classChar(e)
					if err != nil {
						return s, err
					}
					hi = int(b)
				}
			} else if hc == '[' && p.pos+1 < len(p.src) && p.src[p.pos+1] == ':' {
				p.pos = save
			} else {
				hi = int(hc)
				p.pos++
			}
			if hi >= 0 {
				if hi < lo {
					return s, p.errorf("range out of order in character class")
				}
				for ch := lo; ch <= hi; ch++ {
					s.addFold(byte(ch), fl.i)
				}
				continue
			}
		}
		s.addFold(byte(lo), fl.i)
	}
	if neg {
		s.invert()
	}
	return s, nil
}

// classChar decodes a single-character escape inside a class; p.pos is just
// after the escape letter e.
func (p *parser) classChar(e byte) (byte, error) {
	b, ok, err := p.charEscape(e)
	switch {
	case err != nil:
		return 0, err
	case ok:
		return b, nil
	case e == '8' || e == '9':
		return e, nil
	case isDigit(e):
		p.pos--
		return p.readOctal(3), nil
	case isAlnum(e):
		return 0, p.errorf("unrecognized character follows \\")
	}
	return e, nil
}

func (p *parser) readName(term byte) (string, error) {
	st := p.pos
	for p.pos < len(p.src) && (isAlnum(p.src[p.pos]) || p.src[p.pos] == '_') {
		p.pos++
	}
	if p.pos >= len(p.src) || p.src[p.pos] != term {
		return "", p.errorf("syntax error in subpattern name (missing terminator?)")
	}
	name := p.src[st:p.pos]
	p.pos++
	if name == "" {
		return "", p.errorf("group name expected")
	}
	if isDigit(name[0]) {
		return "", p.errorf("subpattern name must start with a non-digit")
	}
	if len(name) > 32 {
		return "", p.errorf("subpattern name is too long")
	}
	return name, nil
}

func (p *parser) parseGroup(fl *flags, depth int) (node, bool, error) {
	p.pos++ // (
	if !p.more() {
		return nil, false, p.errorf("missing closing parenthesis")
	}
	inner := *fl
	if p.src[p.pos] != '?' {
		if p.src[p.pos] == '*' {
			return nil, false, p.unsupported("backtracking control verb")
		}
		return p.capture(&inner, depth, "")
	}
	p.pos++
	if !p.more() {
		return nil, false, p.errorf("missing closing parenthesis")
	}
	c := p.src[p.pos]
	switch c {
	case ':':
		p.pos++
		return p.finishGroup(&inner, depth, func(n node) node { return n })
	case '>':
		p.pos++
		return p.finishGroup(&inner, depth, func(n node) node { return &atomicNode{body: n, hasCaps: hasCaptures(n)} })
	case '=', '!':
		p.pos++
		neg := c == '!'
		return p.finishGroup(&inner, depth, func(n node) node { return &lookNode{body: n, neg: neg, hasCaps: hasCaptures(n)} })
	case '#':
		end := strings.IndexByte(p.src[p.pos:], ')')
		if end < 0 {
			return nil, false, p.errorf("missing ) after (?# comment")
		}
		p.pos += end + 1
		return nil, false, nil
	case '<':
		if p.pos+1 < len(p.src) && (p.src[p.pos+1] == '=' || p.src[p.pos+1] == '!') {
			return nil, false, p.unsupported("look-behind assertion")
		}
		p.pos++
		name, err := p.readName('>')
		if err != nil {
			return nil, false, err
		}
		return p.capture(&inner, depth, name)
	case '\'':
		p.pos++
		name, err := p.readName('\'')
		if err != nil {
			return nil, false, err
		}
		return p.capture(&inner, depth, name)
	case 'P':
		p.pos++
		if !p.more() {
			return nil, false, p.errorf("unrecognized character after (?P")
		}
		switch p.src[p.pos] {
		case '<':
			p.pos++
			name, err := p.readName('>')
			if err != nil {
				return nil, false, err
			}
			return p.capture(&inner, depth, name)
		case '>':
			p.pos++
			name, err := p.readName(')')
			if err != nil {
				return nil, false, err
			}
			return p.call(name, 0), true, nil
		case '=':
			return nil, false, p.unsupported("back reference")
		}
		return nil, false, p.errorf("unrecognized character after (?P")
	case '&':
		p.pos++
		name, err := p.readName(')')
		if err != nil {
			return nil, false, err
		}
		return p.call(name, 0), true, nil
	case 'R':
		if p.pos+1 < len(p.src) && p.src[p.pos+1] == ')' {
			p.pos += 2
			return p.call("", 0), true, nil
		}
	case '(':
		// conditional group; only (?(DEFINE)...) is supported
		if strings.HasPrefix(p.src[p.pos:], "(DEFINE)") {
			p.pos += len("(DEFINE)")
			body, err := p.parseAlt(&inner, depth+1)
			if err != nil {
				return nil, false, err
			}
			if _, ok := body.(*altNode); ok {
				return nil, false, p.errorf("DEFINE subpattern contains more than one branch")
			}
			if !p.more() || p.src[p.pos] != ')' {
				return nil, false, p.errorf("missing closing parenthesis")
			}
			p.pos++
			_ = body
			return defineNode{}, false, nil
		}
		return nil, false, p.unsupported("conditional group")
	}
	if isDigit(c) || c == '+' || c == '-' && p.pos+1 < len(p.src) && isDigit(p.src[p.pos+1]) {
		if c == '+' || c == '-' {
			return nil, false, p.unsupported("relative subroutine call")
		}
		st := p.pos
		for p.more() && isDigit(p.src[p.pos]) {
			p.pos++
		}
		if !p.more() || p.src[p.pos] != ')' {
			return nil, false, p.errorf("(?R (recursive pattern call) must be followed by )")
		}
		n, _ := strconv.Atoi(p.src[st:p.pos])
		p.pos++
		return p.call("", n), true, nil
	}
	// inline option setting: (?imsxU-imsxU) or (?imsx:...)
	on := true
	nf := inner
	for p.more() {
		ch := p.src[p.pos]
		switch ch {
		case 'i':
			nf.i = on
		case 'm':
			nf.m = on
		case 's':
			nf.s = on
		case 'x':
			nf.x = on
		case 'U':
			nf.U = on
		case '-':
			if !on {
				return nil, false, p.errorf("invalid option")
			}
			on = false
		case ')':
			p.pos++
			*fl = nf // applies to the rest of the enclosing group
			return nil, false, nil
		case ':':
			p.pos++
			return p.finishGroup(&nf, depth, func(n node) node { return n })
		default:
			return nil, false, p.errorf("unrecognized character after (? or (?-")
		}
		p.pos++
	}
	return nil, false, p.errorf("missing closing parenthesis")
}

func (p *parser) call(name string, num int) node {
	c := &callNode{name: name, num: num}
	p.calls = append(p.calls, c)
	return c
}

func (p *parser) capture(fl *flags, depth int, name string) (node, bool, error) {
	p.ngroups++
	idx := p.ngroups
	if name != "" {
		if _, dup := p.names[name]; dup {
			return nil, false, p.errorf("two named subpatterns have the same name (PCRE2_DUPNAMES not set)")
		}
		p.names[name] = idx
	}
	g := &capNode{idx: idx}
	p.groups[idx] = g
	n, q, err := p.finishGroup(fl, depth, func(n node) node { g.body = n; return g })
	return n, q, err
}

func (p *parser) finishGroup(fl *flags, depth int, wrap func(node) node) (node, bool, error) {
	body, err := p.parseAlt(fl, depth+1)
	if err != nil {
		return nil, false, err
	}
	if !p.more() || p.src[p.pos] != ')' {
		return nil, false, p.errorf("missing closing parenthesis")
	}
	p.pos++
	return wrap(body), true, nil
}

// mergeLiterals joins adjacent literal nodes with equal case folding.
func mergeLiterals(items []node) []node {
	out := items[:0]
	for _, it := range items {
		if l, ok := it.(*litNode); ok && len(out) > 0 {
			if prev, ok := out[len(out)-1].(*litNode); ok && prev.fold == l.fold {
				out[len(out)-1] = &litNode{s: prev.s + l.s, fold: l.fold}
				continue
			}
		}
		out = append(out, it)
	}
	return out
}

// hasCaptures reports whether matching n can set capture groups (captures set
// inside subroutine calls are reverted, so calls do not count).
func hasCaptures(n node) bool {
	switch t := n.(type) {
	case *capNode:
		return true
	case *seqNode:
		for _, it := range t.items {
			if hasCaptures(it) {
				return true
			}
		}
	case *altNode:
		for _, it := range t.alts {
			if hasCaptures(it) {
				return true
			}
		}
	case *repNode:
		return hasCaptures(t.body)
	case *atomicNode:
		return t.hasCaps
	case *possRepNode:
		return t.hasCaps
	case *lookNode:
		return t.hasCaps
	}
	return false
}
