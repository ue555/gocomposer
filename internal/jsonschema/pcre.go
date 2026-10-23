// Ported for gocomposer (MIT). A small backtracking matcher implementing the
// subset of PCRE2 10.44 syntax and semantics that justinrainbow/json-schema
// 6.13.1 (MIT, Copyright (c) 2016) relies on through PHP's preg_match():
// patterns built by BaseConstraint::jsonPatternToPhpRegex ('~...~u'), the fixed
// patterns in FormatConstraint.php, Rfc3339.php, Tool/Validator/*.php and
// Uri/UriResolver.php, and the e-mail pattern of PHP's ext/filter
// (logical_filters.c, pattern by Michael Rushton). Go's regexp (RE2) cannot be
// used because Composer's schema relies on look-ahead ("^ignore(?!-unreachable$)")
// and because PCRE's "$" also matches before a final newline.
//
// PHP's /u modifier turns on both PCRE2_UTF and PCRE2_UCP, so \d, \w, \s, \b
// and POSIX classes are Unicode-aware in that mode; without /u the subject is
// matched byte by byte with ASCII-only classes and case folding.
package jsonschema

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type reFlags struct {
	caseless, multiline, dotall, extended, ungreedy, dollarEndOnly bool
}

type reNode interface{}

type (
	reLit struct {
		r    rune
		fold bool
	}
	reAny   struct{ dotall bool }
	reClass struct {
		items []reClassItem
		neg   bool
		fold  bool
	}
	reAssert struct {
		kind  byte // ^ $ A z Z b B G
		multi bool
		dEnd  bool
	}
	reGroup struct {
		idx int // -1 for non-capturing
		n   reNode
	}
	reConcat struct{ items []reNode }
	reAlt    struct{ alts []reNode }
	reRepeat struct {
		n        reNode
		min, max int // max -1 = unbounded
		greedy   bool
	}
	reLook struct {
		n          reNode
		ahead, neg bool
	}
	reAtomic  struct{ n reNode }
	reBackref struct {
		idx  int
		name string
		fold bool
	}
)

type reClassItem struct {
	lo, hi rune
	pred   func(rune) bool // when non-nil, lo/hi unused
	neg    bool            // negated predicate
	nofold bool            // predicates are not affected by caseless matching
}

// reProg is a compiled pattern.
type reProg struct {
	root   reNode
	ncap   int
	utf    bool
	names  map[string]int
	anchor bool
}

var errReLimit = errors.New("backtrack limit exhausted")

type reParser struct {
	src   []rune
	pos   int
	flags reFlags
	utf   bool
	ncap  int
	names map[string]int
	refs  []*reBackref
}

// compileRegex compiles a PCRE pattern body with PHP modifier letters
// (i, m, s, x, u, D, U). It returns an error where PCRE2 reports a compile
// error for the constructs this engine understands.
func compileRegex(pattern, mods string) (*reProg, error) {
	p := &reParser{names: map[string]int{}}
	for _, m := range mods {
		switch m {
		case 'i':
			p.flags.caseless = true
		case 'm':
			p.flags.multiline = true
		case 's':
			p.flags.dotall = true
		case 'x':
			p.flags.extended = true
		case 'u':
			p.utf = true
		case 'D':
			p.flags.dollarEndOnly = true
		case 'U':
			p.flags.ungreedy = true
		default:
			return nil, errors.New("unknown modifier")
		}
	}
	if p.utf {
		if !utf8.ValidString(pattern) {
			return nil, errors.New("UTF-8 error")
		}
		p.src = []rune(pattern)
	} else {
		p.src = make([]rune, len(pattern))
		for i := 0; i < len(pattern); i++ {
			p.src[i] = rune(pattern[i])
		}
	}
	root, err := p.parseAlt()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.src) {
		return nil, errors.New("unmatched closing parenthesis")
	}
	for _, r := range p.refs {
		if r.name != "" {
			idx, ok := p.names[r.name]
			if !ok {
				return nil, errors.New("reference to non-existent subpattern")
			}
			r.idx = idx
		}
		if r.idx > p.ncap {
			return nil, errors.New("reference to non-existent subpattern")
		}
	}
	prog := &reProg{root: root, ncap: p.ncap, utf: p.utf, names: p.names}
	prog.anchor = startsAnchored(root)
	return prog, nil
}

func startsAnchored(n reNode) bool {
	switch t := n.(type) {
	case *reAssert:
		return (t.kind == '^' && !t.multi) || t.kind == 'A'
	case *reConcat:
		return len(t.items) > 0 && startsAnchored(t.items[0])
	case *reGroup:
		return startsAnchored(t.n)
	case *reAtomic:
		return startsAnchored(t.n)
	case *reAlt:
		for _, a := range t.alts {
			if !startsAnchored(a) {
				return false
			}
		}
		return len(t.alts) > 0
	}
	return false
}

func (p *reParser) eof() bool  { return p.pos >= len(p.src) }
func (p *reParser) peek() rune { return p.src[p.pos] }

func (p *reParser) skipExtended() {
	if !p.flags.extended {
		return
	}
	for !p.eof() {
		c := p.peek()
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			p.pos++
			continue
		}
		if c == '#' {
			for !p.eof() && p.peek() != '\n' {
				p.pos++
			}
			continue
		}
		break
	}
}

func (p *reParser) parseAlt() (reNode, error) {
	var alts []reNode
	for {
		seq, err := p.parseSeq()
		if err != nil {
			return nil, err
		}
		alts = append(alts, seq)
		if !p.eof() && p.peek() == '|' {
			p.pos++
			continue
		}
		break
	}
	if len(alts) == 1 {
		return alts[0], nil
	}
	return &reAlt{alts: alts}, nil
}

func (p *reParser) parseSeq() (reNode, error) {
	var items []reNode
	for {
		p.skipExtended()
		if p.eof() || p.peek() == '|' || p.peek() == ')' {
			break
		}
		atom, quantifiable, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		if atom == nil {
			continue
		}
		p.skipExtended()
		for !p.eof() {
			min, max, ok, err := p.parseQuantifier()
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
			if !quantifiable {
				return nil, errors.New("quantifier does not follow a repeatable item")
			}
			greedy := !p.flags.ungreedy
			possessive := false
			if !p.eof() && p.peek() == '?' {
				greedy = !greedy
				p.pos++
			} else if !p.eof() && p.peek() == '+' {
				possessive = true
				p.pos++
			}
			if max != -1 && min > max {
				return nil, errors.New("numbers out of order in {} quantifier")
			}
			var n reNode = &reRepeat{n: atom, min: min, max: max, greedy: greedy || possessive}
			if possessive {
				n = &reAtomic{n: n}
			}
			atom = n
			p.skipExtended()
		}
		items = append(items, atom)
	}
	if len(items) == 1 {
		return items[0], nil
	}
	return &reConcat{items: items}, nil
}

func (p *reParser) parseQuantifier() (int, int, bool, error) {
	switch p.peek() {
	case '*':
		p.pos++
		return 0, -1, true, nil
	case '+':
		p.pos++
		return 1, -1, true, nil
	case '?':
		p.pos++
		return 0, 1, true, nil
	case '{':
		// {n} {n,} {n,m} {,m}; anything else is a literal "{"
		j := p.pos + 1
		readNum := func() (int, bool) {
			start := j
			for j < len(p.src) && p.src[j] >= '0' && p.src[j] <= '9' {
				j++
			}
			if start == j {
				return 0, false
			}
			n, err := strconv.Atoi(string(p.src[start:j]))
			if err != nil || n > 65535 {
				return -2, true
			}
			return n, true
		}
		min, hasMin := readNum()
		max := min
		if j < len(p.src) && p.src[j] == ',' {
			j++
			var hasMax bool
			max, hasMax = readNum()
			if !hasMax {
				max = -1
			}
			if !hasMin && !hasMax {
				return 0, 0, false, nil
			}
			if !hasMin {
				min = 0
			}
		} else if !hasMin {
			return 0, 0, false, nil
		}
		if j >= len(p.src) || p.src[j] != '}' {
			return 0, 0, false, nil
		}
		if min == -2 || max == -2 {
			return 0, 0, false, errors.New("number too big in {} quantifier")
		}
		p.pos = j + 1
		return min, max, true, nil
	}
	return 0, 0, false, nil
}

// parseAtom returns nil atom for constructs that produce nothing (option
// settings, comments, \E).
func (p *reParser) parseAtom() (reNode, bool, error) {
	c := p.peek()
	switch c {
	case '(':
		return p.parseGroup()
	case '[':
		n, err := p.parseClass()
		return n, true, err
	case '.':
		p.pos++
		return &reAny{dotall: p.flags.dotall}, true, nil
	case '^':
		p.pos++
		return &reAssert{kind: '^', multi: p.flags.multiline}, false, nil
	case '$':
		p.pos++
		return &reAssert{kind: '$', multi: p.flags.multiline, dEnd: p.flags.dollarEndOnly}, false, nil
	case '\\':
		return p.parseEscape()
	case '*', '+', '?':
		return nil, false, errors.New("quantifier does not follow a repeatable item")
	case '{':
		save := p.pos
		if _, _, ok, _ := p.parseQuantifier(); ok {
			return nil, false, errors.New("quantifier does not follow a repeatable item")
		}
		p.pos = save + 1
		return &reLit{r: '{', fold: p.flags.caseless}, true, nil
	}
	p.pos++
	return &reLit{r: c, fold: p.flags.caseless}, true, nil
}

func (p *reParser) parseGroup() (reNode, bool, error) {
	p.pos++ // (
	saved := p.flags
	idx := -1
	kind := byte(0) // 0 plain, '=', '!', '<'(lookbehind pos), '>'(lookbehind neg), 'a' atomic
	if !p.eof() && p.peek() == '?' {
		p.pos++
		if p.eof() {
			return nil, false, errors.New("unrecognized character after (? or (?-")
		}
		switch c := p.peek(); {
		case c == ':':
			p.pos++
		case c == '=':
			p.pos++
			kind = '='
		case c == '!':
			p.pos++
			kind = '!'
		case c == '>':
			p.pos++
			kind = 'a'
		case c == '#':
			for !p.eof() && p.peek() != ')' {
				p.pos++
			}
			if p.eof() {
				return nil, false, errors.New("missing ) after (?# comment")
			}
			p.pos++
			return nil, false, nil
		case c == '<' && p.pos+1 < len(p.src) && (p.src[p.pos+1] == '=' || p.src[p.pos+1] == '!'):
			if p.src[p.pos+1] == '=' {
				kind = '<'
			} else {
				kind = '>'
			}
			p.pos += 2
		case c == '<' || c == '\'' || c == 'P':
			if c == 'P' {
				p.pos++
				if p.eof() || p.peek() != '<' {
					return nil, false, errors.New("unrecognized character after (?P")
				}
			}
			term := '>'
			if p.peek() == '\'' {
				term = '\''
			}
			p.pos++
			start := p.pos
			for !p.eof() && p.peek() != term {
				r := p.peek()
				if !(r == '_' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
					return nil, false, errors.New("syntax error in subpattern name")
				}
				p.pos++
			}
			if p.eof() || start == p.pos {
				return nil, false, errors.New("syntax error in subpattern name")
			}
			name := string(p.src[start:p.pos])
			p.pos++
			p.ncap++
			idx = p.ncap
			if _, dup := p.names[name]; dup {
				return nil, false, errors.New("two named subpatterns have the same name")
			}
			p.names[name] = idx
		default:
			// option setting (?imsxU-imsxU) or (?imsx-imsx:...)
			on := true
			for {
				if p.eof() {
					return nil, false, errors.New("missing )")
				}
				r := p.peek()
				p.pos++
				switch r {
				case '-':
					on = false
					continue
				case 'i':
					p.flags.caseless = on
					continue
				case 'm':
					p.flags.multiline = on
					continue
				case 's':
					p.flags.dotall = on
					continue
				case 'x':
					p.flags.extended = on
					continue
				case 'U':
					p.flags.ungreedy = on
					continue
				case ')':
					// applies to the remainder of the enclosing group
					return nil, false, nil
				case ':':
				default:
					return nil, false, errors.New("unrecognized character after (? or (?-")
				}
				break
			}
		}
	} else {
		p.ncap++
		idx = p.ncap
	}
	inner, err := p.parseAlt()
	if err != nil {
		return nil, false, err
	}
	if p.eof() || p.peek() != ')' {
		return nil, false, errors.New("missing closing parenthesis")
	}
	p.pos++
	p.flags = saved
	switch kind {
	case '=':
		return &reLook{n: inner, ahead: true}, true, nil
	case '!':
		return &reLook{n: inner, ahead: true, neg: true}, true, nil
	case '<':
		return &reLook{n: inner}, true, nil
	case '>':
		return &reLook{n: inner, neg: true}, true, nil
	case 'a':
		return &reAtomic{n: inner}, true, nil
	}
	return &reGroup{idx: idx, n: inner}, true, nil
}

func isASCIIWord(r rune) bool {
	return r == '_' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func (p *reParser) isWord(r rune) bool {
	if !p.utf {
		return isASCIIWord(r)
	}
	return ucpWord(r)
}

func ucpWord(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Pc, r)
}

func hSpace(r rune) bool {
	switch r {
	case '\t', ' ', 0xa0, 0x1680, 0x180e, 0x202f, 0x205f, 0x3000:
		return true
	}
	return r >= 0x2000 && r <= 0x200a
}

func vSpace(r rune) bool {
	switch r {
	case '\n', '\v', '\f', '\r', 0x85, 0x2028, 0x2029:
		return true
	}
	return false
}

// classPred returns the predicate for \d \w \s \h \v (lower-case letter).
func (p *reParser) classPred(c rune) func(rune) bool {
	utf := p.utf
	switch c {
	case 'd':
		if utf {
			return func(r rune) bool { return unicode.Is(unicode.Nd, r) }
		}
		return func(r rune) bool { return r >= '0' && r <= '9' }
	case 'w':
		if utf {
			return ucpWord
		}
		return isASCIIWord
	case 's':
		if utf {
			return func(r rune) bool {
				return r == '\t' || r == '\n' || r == '\v' || r == '\f' || r == '\r' || r == ' ' || r == 0x85 || unicode.Is(unicode.Z, r)
			}
		}
		return func(r rune) bool { return r == '\t' || r == '\n' || r == '\v' || r == '\f' || r == '\r' || r == ' ' }
	case 'h':
		if utf {
			return hSpace
		}
		return func(r rune) bool { return r == '\t' || r == ' ' || r == 0xa0 }
	case 'v':
		if utf {
			return vSpace
		}
		return func(r rune) bool { return r == '\n' || r == '\v' || r == '\f' || r == '\r' || r == 0x85 }
	}
	return nil
}

func lookupProperty(name string) (func(rune) bool, bool) {
	neg := false
	if strings.HasPrefix(name, "^") {
		neg = true
		name = name[1:]
	}
	var f func(rune) bool
	switch name {
	case "Any":
		f = func(rune) bool { return true }
	case "L&", "LC":
		f = func(r rune) bool { return unicode.In(r, unicode.Lu, unicode.Ll, unicode.Lt) }
	case "Xan":
		f = func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) }
	case "Xwd":
		f = ucpWord
	case "Xsp", "Xps":
		f = func(r rune) bool {
			return r == '\t' || r == '\n' || r == '\v' || r == '\f' || r == '\r' || unicode.Is(unicode.Z, r)
		}
	default:
		if t, ok := unicode.Categories[name]; ok {
			f = func(r rune) bool { return unicode.Is(t, r) }
		} else if t, ok := unicode.Scripts[name]; ok {
			f = func(r rune) bool { return unicode.Is(t, r) }
		} else {
			for k, t := range unicode.Scripts {
				if strings.EqualFold(strings.ReplaceAll(k, "_", ""), strings.ReplaceAll(strings.ReplaceAll(name, "_", ""), " ", "")) {
					t := t
					f = func(r rune) bool { return unicode.Is(t, r) }
					break
				}
			}
		}
	}
	if f == nil {
		return nil, false
	}
	if neg {
		g := f
		f = func(r rune) bool { return !g(r) }
	}
	return f, true
}

// parseProperty parses after \p or \P.
func (p *reParser) parseProperty(negated bool) (reClassItem, error) {
	if p.eof() {
		return reClassItem{}, errors.New("malformed \\P or \\p sequence")
	}
	var name string
	if p.peek() == '{' {
		end := p.pos + 1
		for end < len(p.src) && p.src[end] != '}' {
			end++
		}
		if end >= len(p.src) {
			return reClassItem{}, errors.New("malformed \\P or \\p sequence")
		}
		name = string(p.src[p.pos+1 : end])
		p.pos = end + 1
	} else {
		name = string(p.peek())
		p.pos++
	}
	f, ok := lookupProperty(name)
	if !ok {
		return reClassItem{}, errors.New("unknown property name after \\P or \\p")
	}
	return reClassItem{pred: f, neg: negated, nofold: true}, nil
}

// parseCharEscape handles escapes that denote a single character. ok=false if
// the escape is not a character escape.
func (p *reParser) parseCharEscape(c rune) (rune, bool, error) {
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
		for k := 0; k < 2 && !p.eof() && p.peek() >= '0' && p.peek() <= '7'; k++ {
			v = v*8 + int(p.peek()-'0')
			p.pos++
		}
		return rune(v), true, nil
	case 'o':
		if p.eof() || p.peek() != '{' {
			return 0, false, errors.New("missing opening brace after \\o")
		}
		end := p.pos + 1
		for end < len(p.src) && p.src[end] != '}' {
			end++
		}
		if end >= len(p.src) {
			return 0, false, errors.New("missing closing brace")
		}
		v, err := strconv.ParseUint(string(p.src[p.pos+1:end]), 8, 32)
		if err != nil {
			return 0, false, errors.New("non-octal character in \\o{}")
		}
		p.pos = end + 1
		return rune(v), true, nil
	case 'x':
		if !p.eof() && p.peek() == '{' {
			end := p.pos + 1
			for end < len(p.src) && p.src[end] != '}' {
				end++
			}
			if end >= len(p.src) {
				return 0, false, errors.New("missing closing brace")
			}
			v, err := strconv.ParseUint(string(p.src[p.pos+1:end]), 16, 32)
			if err != nil {
				return 0, false, errors.New("non-hex character in \\x{}")
			}
			if (!p.utf && v > 0xff) || v > 0x10ffff || (v >= 0xd800 && v <= 0xdfff) {
				return 0, false, errors.New("character code point value in \\x{} or \\o{} is too large")
			}
			p.pos = end + 1
			return rune(v), true, nil
		}
		v := 0
		for k := 0; k < 2 && !p.eof() && isHex(p.peek()); k++ {
			v = v*16 + hexVal(p.peek())
			p.pos++
		}
		return rune(v), true, nil
	case 'c':
		if p.eof() || p.peek() > 127 {
			return 0, false, errors.New("\\c must be followed by a printable ASCII character")
		}
		r := unicode.ToUpper(p.peek())
		p.pos++
		return r ^ 0x40, true, nil
	}
	if c < 128 && (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
		return 0, false, nil
	}
	return c, true, nil
}

func isHex(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func hexVal(r rune) int {
	switch {
	case r >= '0' && r <= '9':
		return int(r - '0')
	case r >= 'a' && r <= 'f':
		return int(r-'a') + 10
	}
	return int(r-'A') + 10
}

func (p *reParser) parseEscape() (reNode, bool, error) {
	p.pos++ // backslash
	if p.eof() {
		return nil, false, errors.New("\\ at end of pattern")
	}
	c := p.peek()
	p.pos++
	switch c {
	case 'd', 'w', 's', 'h', 'v':
		return &reClass{items: []reClassItem{{pred: p.classPred(c), nofold: true}}}, true, nil
	case 'D', 'W', 'S', 'H', 'V':
		return &reClass{items: []reClassItem{{pred: p.classPred(unicode.ToLower(c)), nofold: true}}, neg: true}, true, nil
	case 'N':
		return &reAny{dotall: false}, true, nil
	case 'p', 'P':
		it, err := p.parseProperty(c == 'P')
		if err != nil {
			return nil, false, err
		}
		return &reClass{items: []reClassItem{it}}, true, nil
	case 'b', 'B', 'A', 'z', 'Z', 'G':
		return &reAssert{kind: byte(c)}, false, nil
	case 'Q':
		var lits []reNode
		for !p.eof() {
			if p.peek() == '\\' && p.pos+1 < len(p.src) && p.src[p.pos+1] == 'E' {
				p.pos += 2
				break
			}
			lits = append(lits, &reLit{r: p.peek(), fold: p.flags.caseless})
			p.pos++
		}
		if len(lits) == 0 {
			return nil, false, nil
		}
		if len(lits) == 1 {
			return lits[0], true, nil
		}
		// only the last character is quantifiable; emulate by grouping all but last
		return &reConcat{items: lits}, true, nil
	case 'E':
		return nil, false, nil
	case 'g', 'k':
		return p.parseNamedRef(c)
	}
	if c >= '1' && c <= '9' {
		start := p.pos - 1
		for !p.eof() && p.peek() >= '0' && p.peek() <= '9' {
			p.pos++
		}
		n, _ := strconv.Atoi(string(p.src[start:p.pos]))
		ref := &reBackref{idx: n, fold: p.flags.caseless}
		p.refs = append(p.refs, ref)
		return ref, true, nil
	}
	r, ok, err := p.parseCharEscape(c)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, errors.New("unrecognized character follows \\")
	}
	return &reLit{r: r, fold: p.flags.caseless}, true, nil
}

func (p *reParser) parseNamedRef(c rune) (reNode, bool, error) {
	if p.eof() {
		return nil, false, errors.New("malformed reference")
	}
	open := p.peek()
	var close rune
	switch open {
	case '{':
		close = '}'
	case '<':
		close = '>'
	case '\'':
		close = '\''
	default:
		if c == 'g' && (open >= '0' && open <= '9' || open == '-') {
			start := p.pos
			p.pos++
			for !p.eof() && p.peek() >= '0' && p.peek() <= '9' {
				p.pos++
			}
			return p.numericRef(string(p.src[start:p.pos]))
		}
		return nil, false, errors.New("malformed reference")
	}
	end := p.pos + 1
	for end < len(p.src) && p.src[end] != close {
		end++
	}
	if end >= len(p.src) {
		return nil, false, errors.New("malformed reference")
	}
	body := string(p.src[p.pos+1 : end])
	p.pos = end + 1
	if c == 'g' && open == '{' && body != "" && (body[0] == '-' || body[0] >= '0' && body[0] <= '9') {
		return p.numericRef(body)
	}
	if c == 'g' && open != '{' {
		return nil, false, errors.New("subroutine calls are not supported")
	}
	ref := &reBackref{name: body, fold: p.flags.caseless}
	p.refs = append(p.refs, ref)
	return ref, true, nil
}

func (p *reParser) numericRef(s string) (reNode, bool, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n == 0 {
		return nil, false, errors.New("malformed reference")
	}
	if n < 0 {
		n = p.ncap + 1 + n
		if n <= 0 {
			return nil, false, errors.New("reference to non-existent subpattern")
		}
	}
	ref := &reBackref{idx: n, fold: p.flags.caseless}
	p.refs = append(p.refs, ref)
	return ref, true, nil
}

var posixClasses = map[string][2]func(rune) bool{
	// [ascii, ucp]
	"alpha":  {func(r rune) bool { return r < 128 && (r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') }, unicode.IsLetter},
	"digit":  {func(r rune) bool { return r >= '0' && r <= '9' }, func(r rune) bool { return unicode.Is(unicode.Nd, r) }},
	"alnum":  {func(r rune) bool { return r < 128 && (isASCIIWord(r) && r != '_') }, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) }},
	"space":  {func(r rune) bool { return r == ' ' || r >= '\t' && r <= '\r' }, func(r rune) bool { return r >= '\t' && r <= '\r' || r == ' ' || r == 0x85 || unicode.Is(unicode.Z, r) }},
	"word":   {isASCIIWord, ucpWord},
	"lower":  {func(r rune) bool { return r >= 'a' && r <= 'z' }, func(r rune) bool { return unicode.Is(unicode.Ll, r) }},
	"upper":  {func(r rune) bool { return r >= 'A' && r <= 'Z' }, func(r rune) bool { return unicode.Is(unicode.Lu, r) }},
	"xdigit": {isHex, isHex},
	"blank":  {func(r rune) bool { return r == ' ' || r == '\t' }, hSpace},
	"cntrl":  {func(r rune) bool { return r < 32 || r == 127 }, func(r rune) bool { return unicode.Is(unicode.Cc, r) }},
	"punct":  {func(r rune) bool { return r > 32 && r < 127 && !isASCIIWord(r) || r == '_' }, func(r rune) bool { return r > 32 && r < 127 && !isASCIIWord(r) || r == '_' }},
	"graph": {func(r rune) bool { return r > 32 && r < 127 }, func(r rune) bool {
		return r > 32 && r < 127 || r >= 128 && unicode.IsGraphic(r) && !unicode.Is(unicode.Z, r)
	}},
	"print": {func(r rune) bool { return r >= 32 && r < 127 }, func(r rune) bool { return r >= 32 && r < 127 || r >= 128 && unicode.IsPrint(r) }},
	"ascii": {func(r rune) bool { return r < 128 }, func(r rune) bool { return r < 128 }},
}

func (p *reParser) parseClass() (reNode, error) {
	p.pos++ // [
	cls := &reClass{fold: p.flags.caseless}
	if !p.eof() && p.peek() == '^' {
		cls.neg = true
		p.pos++
	}
	first := true
	for {
		if p.eof() {
			return nil, errors.New("missing terminating ] for character class")
		}
		c := p.peek()
		if c == ']' && !first {
			p.pos++
			break
		}
		first = false
		// POSIX class
		if c == '[' && p.pos+1 < len(p.src) && (p.src[p.pos+1] == ':' || p.src[p.pos+1] == '.' || p.src[p.pos+1] == '=') {
			delim := p.src[p.pos+1]
			end := p.pos + 2
			for end+1 < len(p.src) && !(p.src[end] == delim && p.src[end+1] == ']') {
				if p.src[end] == ']' {
					end = -1
					break
				}
				end++
			}
			if end > 0 && end+1 < len(p.src) {
				if delim != ':' {
					return nil, errors.New("POSIX collating elements are not supported")
				}
				name := string(p.src[p.pos+2 : end])
				neg := false
				if strings.HasPrefix(name, "^") {
					neg = true
					name = name[1:]
				}
				f, ok := posixClasses[name]
				if !ok {
					return nil, errors.New("unknown POSIX class name")
				}
				pred := f[0]
				if p.utf {
					pred = f[1]
				}
				nofold := name != "lower" && name != "upper"
				if !nofold {
					// [:lower:] / [:upper:] match both cases when caseless
					if p.flags.caseless {
						pred = posixClasses["alpha"][0]
						if p.utf {
							pred = func(r rune) bool { return unicode.IsUpper(r) || unicode.IsLower(r) || unicode.IsTitle(r) }
						}
					}
					nofold = true
				}
				cls.items = append(cls.items, reClassItem{pred: pred, neg: neg, nofold: nofold})
				p.pos = end + 2
				continue
			}
		}
		lo, item, isItem, err := p.classAtom()
		if err != nil {
			return nil, err
		}
		if isItem {
			cls.items = append(cls.items, item)
			continue
		}
		// range?
		if p.pos+1 < len(p.src) && p.peek() == '-' && p.src[p.pos+1] != ']' {
			save := p.pos
			p.pos++
			hi, _, hiItem, err := p.classAtom()
			if err != nil {
				return nil, err
			}
			if hiItem {
				// "a-\d": '-' is literal
				p.pos = save
			} else {
				if hi < lo {
					return nil, errors.New("range out of order in character class")
				}
				cls.items = append(cls.items, reClassItem{lo: lo, hi: hi})
				continue
			}
		}
		cls.items = append(cls.items, reClassItem{lo: lo, hi: lo})
	}
	return cls, nil
}

// classAtom reads one class member: either a single character (returned as
// lo) or a predicate item.
func (p *reParser) classAtom() (rune, reClassItem, bool, error) {
	c := p.peek()
	p.pos++
	if c != '\\' {
		return c, reClassItem{}, false, nil
	}
	if p.eof() {
		return 0, reClassItem{}, false, errors.New("\\ at end of pattern")
	}
	e := p.peek()
	p.pos++
	switch e {
	case 'd', 'w', 's', 'h', 'v':
		return 0, reClassItem{pred: p.classPred(e), nofold: true}, true, nil
	case 'D', 'W', 'S', 'H', 'V':
		return 0, reClassItem{pred: p.classPred(unicode.ToLower(e)), neg: true, nofold: true}, true, nil
	case 'p', 'P':
		it, err := p.parseProperty(e == 'P')
		return 0, it, true, err
	case 'b':
		return '\b', reClassItem{}, false, nil
	case 'Q', 'E':
		return 0, reClassItem{}, false, errors.New("\\Q/\\E inside a class is not supported")
	}
	r, ok, err := p.parseCharEscape(e)
	if err != nil {
		return 0, reClassItem{}, false, err
	}
	if !ok {
		return 0, reClassItem{}, false, errors.New("unrecognized character follows \\")
	}
	return r, reClassItem{}, false, nil
}

// ---- matching ----

const reStepLimit = 20_000_000

type reMatcher struct {
	prog  *reProg
	in    []rune
	caps  []int
	steps int
	start int
}

func (m *reMatcher) fold(r rune) []rune {
	if !m.prog.utf {
		switch {
		case r >= 'a' && r <= 'z':
			return []rune{r - 32}
		case r >= 'A' && r <= 'Z':
			return []rune{r + 32}
		}
		return nil
	}
	var out []rune
	for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
		out = append(out, f)
	}
	return out
}

func (m *reMatcher) litEq(a, b rune, fold bool) bool {
	if a == b {
		return true
	}
	if !fold {
		return false
	}
	for _, f := range m.fold(b) {
		if f == a {
			return true
		}
	}
	return false
}

func (m *reMatcher) classHas(c *reClass, r rune) bool {
	in := m.classRaw(c, r, false)
	if !in && c.fold {
		for _, f := range m.fold(r) {
			if m.classRaw(c, f, true) {
				in = true
				break
			}
		}
	}
	return in != c.neg
}

func (m *reMatcher) classRaw(c *reClass, r rune, folded bool) bool {
	for _, it := range c.items {
		if it.pred != nil {
			if folded && it.nofold {
				continue
			}
			if it.pred(r) != it.neg {
				return true
			}
			continue
		}
		if r >= it.lo && r <= it.hi {
			return true
		}
	}
	return false
}

func (m *reMatcher) isWordAt(i int) bool {
	if i < 0 || i >= len(m.in) {
		return false
	}
	if m.prog.utf {
		return ucpWord(m.in[i])
	}
	return isASCIIWord(m.in[i])
}

func (m *reMatcher) run(n reNode, i int, k func(int) bool) bool {
	m.steps++
	if m.steps > reStepLimit {
		panic(errReLimit)
	}
	switch t := n.(type) {
	case *reLit:
		if i < len(m.in) && m.litEq(t.r, m.in[i], t.fold) {
			return k(i + 1)
		}
		return false
	case *reAny:
		if i < len(m.in) && (t.dotall || m.in[i] != '\n') {
			return k(i + 1)
		}
		return false
	case *reClass:
		if i < len(m.in) && m.classHas(t, m.in[i]) {
			return k(i + 1)
		}
		return false
	case *reAssert:
		if m.assert(t, i) {
			return k(i)
		}
		return false
	case *reConcat:
		return m.seq(t.items, 0, i, k)
	case *reAlt:
		for _, a := range t.alts {
			if m.run(a, i, k) {
				return true
			}
		}
		return false
	case *reGroup:
		if t.idx < 0 {
			return m.run(t.n, i, k)
		}
		return m.run(t.n, i, func(j int) bool {
			os, oe := m.caps[2*t.idx], m.caps[2*t.idx+1]
			m.caps[2*t.idx], m.caps[2*t.idx+1] = i, j
			if k(j) {
				return true
			}
			m.caps[2*t.idx], m.caps[2*t.idx+1] = os, oe
			return false
		})
	case *reRepeat:
		return m.rep(t, i, 0, k)
	case *reAtomic:
		saved := append([]int(nil), m.caps...)
		end := -1
		if !m.run(t.n, i, func(j int) bool { end = j; return true }) {
			copy(m.caps, saved)
			return false
		}
		if k(end) {
			return true
		}
		copy(m.caps, saved)
		return false
	case *reLook:
		saved := append([]int(nil), m.caps...)
		var ok bool
		if t.ahead {
			ok = m.run(t.n, i, func(int) bool { return true })
		} else {
			for st := i; st >= 0 && !ok; st-- {
				ok = m.run(t.n, st, func(j int) bool { return j == i })
			}
		}
		if t.neg {
			copy(m.caps, saved)
			if ok {
				return false
			}
			return k(i)
		}
		if !ok {
			copy(m.caps, saved)
			return false
		}
		if k(i) {
			return true
		}
		copy(m.caps, saved)
		return false
	case *reBackref:
		s, e := m.caps[2*t.idx], m.caps[2*t.idx+1]
		if s < 0 {
			return false // unset group never matches in PCRE2 default mode
		}
		j := i
		for x := s; x < e; x++ {
			if j >= len(m.in) || !m.litEq(m.in[x], m.in[j], t.fold) {
				return false
			}
			j++
		}
		return k(j)
	}
	panic("jsonschema: unknown regex node")
}

func (m *reMatcher) seq(items []reNode, idx, i int, k func(int) bool) bool {
	if idx == len(items) {
		return k(i)
	}
	return m.run(items[idx], i, func(j int) bool { return m.seq(items, idx+1, j, k) })
}

func (m *reMatcher) rep(t *reRepeat, i, count int, k func(int) bool) bool {
	if count < t.min {
		return m.run(t.n, i, func(j int) bool {
			if j == i {
				// empty iteration: further mandatory iterations would also be empty
				return k(j)
			}
			return m.rep(t, j, count+1, k)
		})
	}
	if t.max != -1 && count >= t.max {
		return k(i)
	}
	next := func(j int) bool {
		if j == i {
			return false
		}
		return m.rep(t, j, count+1, k)
	}
	if t.greedy {
		if m.run(t.n, i, next) {
			return true
		}
		return k(i)
	}
	if k(i) {
		return true
	}
	return m.run(t.n, i, next)
}

func (m *reMatcher) assert(t *reAssert, i int) bool {
	n := len(m.in)
	switch t.kind {
	case '^':
		if t.multi {
			return i == 0 || (m.in[i-1] == '\n' && i < n)
		}
		return i == 0
	case '$':
		if t.multi {
			return i == n || m.in[i] == '\n'
		}
		if t.dEnd {
			return i == n
		}
		return i == n || (i == n-1 && m.in[i] == '\n')
	case 'A':
		return i == 0
	case 'G':
		return i == m.start
	case 'z':
		return i == n
	case 'Z':
		return i == n || (i == n-1 && m.in[i] == '\n')
	case 'b':
		return m.isWordAt(i-1) != m.isWordAt(i)
	case 'B':
		return m.isWordAt(i-1) == m.isWordAt(i)
	}
	return false
}

// reMatch is the result of a successful match: capture spans in subject
// units (runes in UTF mode, bytes otherwise), -1 when unset.
type reMatch struct {
	in   []rune
	caps []int
	utf  bool
}

// group returns the captured text of group i ("" when unset).
func (r *reMatch) group(i int) string {
	if 2*i+1 >= len(r.caps) || r.caps[2*i] < 0 {
		return ""
	}
	return runesToString(r.in[r.caps[2*i]:r.caps[2*i+1]], r.utf)
}

// count mirrors count($matches) of preg_match: trailing unset groups are
// dropped.
func (r *reMatch) count() int {
	n := len(r.caps) / 2
	for n > 1 && r.caps[2*(n-1)] < 0 {
		n--
	}
	return n
}

func runesToString(rs []rune, utf bool) string {
	// subject runes in byte mode are bytes 0..255; in UTF mode code points.
	if utf {
		return string(rs)
	}
	b := make([]byte, len(rs))
	for i, r := range rs {
		b[i] = byte(r)
	}
	return string(b)
}

// exec runs preg_match semantics: (match, ok, err). err is non-nil when PHP's
// preg_match would return false (invalid UTF-8 subject in /u mode, backtrack
// limit).
func (prog *reProg) exec(subject string) (res *reMatch, err error) {
	var in []rune
	if prog.utf {
		if !utf8.ValidString(subject) {
			return nil, errors.New("Malformed UTF-8 data")
		}
		in = []rune(subject)
	} else {
		in = make([]rune, len(subject))
		for i := 0; i < len(subject); i++ {
			in[i] = rune(subject[i])
		}
	}
	m := &reMatcher{prog: prog, in: in, caps: make([]int, 2*(prog.ncap+1))}
	defer func() {
		if r := recover(); r != nil {
			if r == errReLimit {
				res, err = nil, errReLimit
				return
			}
			panic(r)
		}
	}()
	last := len(in)
	if prog.anchor {
		last = 0
	}
	for st := 0; st <= last; st++ {
		for x := range m.caps {
			m.caps[x] = -1
		}
		m.start = 0
		if m.run(prog.root, st, func(j int) bool {
			m.caps[0], m.caps[1] = st, j
			return true
		}) {
			return &reMatch{in: in, caps: m.caps, utf: prog.utf}, nil
		}
	}
	return nil, nil
}

// pregMatch mirrors preg_match() returning 1, 0 or false (-1).
func pregMatch(prog *reProg, subject string) int {
	res, err := prog.exec(subject)
	if err != nil {
		return -1
	}
	if res == nil {
		return 0
	}
	return 1
}

// mustRegex compiles a fixed pattern used by the library itself.
func mustRegex(pattern, mods string) *reProg {
	p, err := compileRegex(pattern, mods)
	if err != nil {
		panic("jsonschema: bad internal regex " + pattern + ": " + err.Error())
	}
	return p
}
