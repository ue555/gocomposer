// Backtracking matcher for the PCRE2 subset parsed in parse.go.
// Ported semantics: PCRE2 10.44 as bundled with PHP 8.4 (used by
// src/Composer/Json/JsonManipulator.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package pcre

import "errors"

// ErrBacktrackLimit corresponds to PREG_BACKTRACK_LIMIT_ERROR: the matcher
// gave up because the pattern needed too many steps.
var ErrBacktrackLimit = errors.New("Backtrack limit exhausted")

// ErrRecursionLimit corresponds to PREG_RECURSION_LIMIT_ERROR /
// PREG_JIT_STACKLIMIT_ERROR: subroutine calls nested too deeply.
var ErrRecursionLimit = errors.New("Recursion limit exhausted")

// PHP runs PCRE2 with pcre.backtrack_limit=1000000 and pcre.jit=1. PCRE2's
// JIT counts "match calls" in its own way, so the engine approximates it with
// weighted step costs. The weights were fitted against official PHP 8.4.27 /
// PCRE2 10.44 JIT thresholds for JsonManipulator's patterns on growing
// documents (see internal/pcre/testdata/limits.json and TestLimitCalibration);
// thresholds agree within about 1.5%.
const (
	costAlt    = 3     // each alternative tried
	costCall   = 858   // subroutine call
	costRep    = 575   // group repetition step
	costSetPos = 44    // backtracking position of a single-byte repeat
	costCap    = 10000 // entering a capturing group (incl. inside calls)
	costAtomic = 1998  // atomic group / possessive group repeat
	costLook   = 1835  // look-ahead
	costPoss   = 1867  // possessive single-byte repeat
	costLit    = 90    // literal comparison
)

// MatchLimit is the cost budget equivalent to pcre.backtrack_limit=1000000.
var MatchLimit int64 = 1_433_000_000

// StackLimit approximates PHP's PCRE2 JIT stack (PREG_JIT_STACKLIMIT_ERROR):
// in stack units a nested array costs 8 and a nested object 22, matching the
// 3068 / 1115 nesting levels at which PHP 8.4 fails.
var StackLimit = 24544

const (
	stackCall = 1
	stackRep  = 3
	stackAlt  = 3
)

type cont func(int) bool

type node interface {
	match(m *machine, i int, k cont) bool
}

// simpleNode is a deterministic node: it either fails or advances to exactly
// one position, so sequences can step over it without closures.
type simpleNode interface {
	node
	step(m *machine, i int) (int, bool)
}

type machine struct {
	s     string
	caps  []int
	noCap int
	cost  int64
	err   error
	stack int
	maxSt int
	// unlimited disables MatchLimit/StackLimit (calibration tests)
	unlimited bool
}

func (m *machine) tick(c int64) bool {
	m.cost += c
	if m.cost > MatchLimit && !m.unlimited {
		if m.err == nil {
			m.err = ErrBacktrackLimit
		}
		return false
	}
	return m.err == nil
}

// push accounts for stack usage; it fails once StackLimit is exceeded.
func (m *machine) push(u int) bool {
	m.stack += u
	if m.stack > m.maxSt {
		m.maxSt = m.stack
	}
	if m.stack > StackLimit && !m.unlimited {
		if m.err == nil {
			m.err = ErrRecursionLimit
		}
		return false
	}
	return true
}

// ---- character sets ----

type charSet [4]uint64

func (s *charSet) add(c byte)      { s[c>>6] |= 1 << (c & 63) }
func (s *charSet) del(c byte)      { s[c>>6] &^= 1 << (c & 63) }
func (s *charSet) has(c byte) bool { return s[c>>6]&(1<<(c&63)) != 0 }
func (s *charSet) fill() {
	for i := range s {
		s[i] = ^uint64(0)
	}
}
func (s *charSet) invert() {
	for i := range s {
		s[i] = ^s[i]
	}
}

func (s *charSet) addFold(c byte, fold bool) {
	s.add(c)
	if fold {
		if c >= 'a' && c <= 'z' {
			s.add(c - 32)
		} else if c >= 'A' && c <= 'Z' {
			s.add(c + 32)
		}
	}
}

func (s *charSet) addClassEscape(e byte) {
	var t charSet
	switch e {
	case 'd', 'D':
		for c := '0'; c <= '9'; c++ {
			t.add(byte(c))
		}
	case 's', 'S':
		for _, c := range []byte{' ', '\t', '\n', '\v', '\f', '\r'} {
			t.add(c)
		}
	case 'w', 'W':
		for c := 0; c < 256; c++ {
			if isAlnum(byte(c)) || c == '_' {
				t.add(byte(c))
			}
		}
	case 'h', 'H':
		t.add(' ')
		t.add('\t')
		t.add(0xa0)
	case 'v', 'V':
		for _, c := range []byte{'\n', '\v', '\f', '\r', 0x85} {
			t.add(c)
		}
	}
	if e >= 'A' && e <= 'Z' {
		t.invert()
	}
	for i := range s {
		s[i] |= t[i]
	}
}

func isWord(c byte) bool { return isAlnum(c) || c == '_' }

// ---- nodes ----

type litNode struct {
	s    string
	fold bool
}

func newLit(s string, fold bool) *litNode {
	if fold {
		hasLetter := false
		for i := 0; i < len(s); i++ {
			c := s[i] | 0x20
			if c >= 'a' && c <= 'z' {
				hasLetter = true
			}
		}
		fold = hasLetter
	}
	return &litNode{s: s, fold: fold}
}

func (n *litNode) charSet() charSet {
	var s charSet
	s.addFold(n.s[0], n.fold)
	return s
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

func (n *litNode) step(m *machine, i int) (int, bool) {
	m.cost += costLit
	if len(m.s)-i < len(n.s) {
		return 0, false
	}
	if !n.fold {
		if m.s[i:i+len(n.s)] != n.s {
			return 0, false
		}
		return i + len(n.s), true
	}
	for k := 0; k < len(n.s); k++ {
		if lower(m.s[i+k]) != lower(n.s[k]) {
			return 0, false
		}
	}
	return i + len(n.s), true
}

func (n *litNode) match(m *machine, i int, k cont) bool {
	j, ok := n.step(m, i)
	return ok && k(j)
}

type setNode struct{ set charSet }

func (n *setNode) step(m *machine, i int) (int, bool) {
	if i < len(m.s) && n.set.has(m.s[i]) {
		return i + 1, true
	}
	return 0, false
}

func (n *setNode) match(m *machine, i int, k cont) bool {
	j, ok := n.step(m, i)
	return ok && k(j)
}

const (
	aCaret = iota
	aDollar
	aStartSubject
	aEndSubject
	aEndSubjectNL
	aWordB
	aNotWordB
)

type assertNode struct {
	kind          int
	multiline     bool
	dollarEndOnly bool
}

func (n assertNode) step(m *machine, i int) (int, bool) {
	s := m.s
	switch n.kind {
	case aCaret:
		if i == 0 {
			return i, true
		}
		if n.multiline && s[i-1] == '\n' && i < len(s) {
			// PCRE2: in multiline mode ^ does not match after a final newline
			return i, true
		}
		return 0, false
	case aDollar:
		if i == len(s) {
			return i, true
		}
		if n.multiline {
			return i, s[i] == '\n'
		}
		if !n.dollarEndOnly && i == len(s)-1 && s[i] == '\n' {
			return i, true
		}
		return 0, false
	case aStartSubject:
		return i, i == 0
	case aEndSubject:
		return i, i == len(s)
	case aEndSubjectNL:
		return i, i == len(s) || i == len(s)-1 && s[i] == '\n'
	case aWordB, aNotWordB:
		a := i > 0 && isWord(s[i-1])
		b := i < len(s) && isWord(s[i])
		return i, (a != b) == (n.kind == aWordB)
	}
	return 0, false
}

func (n assertNode) match(m *machine, i int, k cont) bool {
	j, ok := n.step(m, i)
	return ok && k(j)
}

// defineNode is (?(DEFINE)...): it never matches its body.
type defineNode struct{}

func (defineNode) step(m *machine, i int) (int, bool) { return i, true }
func (defineNode) match(m *machine, i int, k cont) bool {
	return k(i)
}

type seqNode struct{ items []node }

func (n *seqNode) match(m *machine, i int, k cont) bool { return matchSeq(m, n.items, i, k) }

func matchSeq(m *machine, items []node, i int, k cont) bool {
	for len(items) > 0 {
		if s, ok := items[0].(simpleNode); ok {
			j, ok := s.step(m, i)
			if !ok {
				return false
			}
			i = j
			items = items[1:]
			continue
		}
		if len(items) == 1 {
			return items[0].match(m, i, k)
		}
		rest := items[1:]
		return items[0].match(m, i, func(j int) bool { return matchSeq(m, rest, j, k) })
	}
	return k(i)
}

type altNode struct{ alts []node }

func (n *altNode) match(m *machine, i int, k cont) bool {
	if !m.push(stackAlt) {
		m.stack -= stackAlt
		return false
	}
	defer func() { m.stack -= stackAlt }()
	for _, a := range n.alts {
		if !m.tick(costAlt) {
			return false
		}
		if a.match(m, i, k) {
			return true
		}
		if m.err != nil {
			return false
		}
	}
	return false
}

type capNode struct {
	idx  int
	body node
}

func (n *capNode) match(m *machine, i int, k cont) bool {
	if !m.tick(costCap) {
		return false
	}
	if m.noCap > 0 {
		return n.body.match(m, i, k)
	}
	return n.body.match(m, i, func(j int) bool {
		os, oe := m.caps[2*n.idx], m.caps[2*n.idx+1]
		m.caps[2*n.idx], m.caps[2*n.idx+1] = i, j
		if k(j) {
			return true
		}
		m.caps[2*n.idx], m.caps[2*n.idx+1] = os, oe
		return false
	})
}

// callNode is a subroutine call. Since PCRE2 10.30 calls are not atomic;
// captures set during the call revert when it returns, which is modelled by
// suppressing capture recording inside the call.
type callNode struct {
	name   string
	num    int
	target node
}

func (n *callNode) match(m *machine, i int, k cont) bool {
	if !m.tick(costCall) {
		return false
	}
	if !m.push(stackCall) {
		m.stack -= stackCall
		return false
	}
	m.noCap++
	r := n.target.match(m, i, func(j int) bool {
		m.noCap--
		ok := k(j)
		m.noCap++
		return ok
	})
	m.noCap--
	m.stack -= stackCall
	return r
}

type atomicNode struct {
	body    node
	hasCaps bool
}

func (n *atomicNode) match(m *machine, i int, k cont) bool {
	if !m.tick(costAtomic) {
		return false
	}
	var saved []int
	if n.hasCaps && m.noCap == 0 {
		saved = append([]int(nil), m.caps...)
	}
	end := -1
	if !n.body.match(m, i, func(j int) bool { end = j; return true }) {
		return false
	}
	if k(end) {
		return true
	}
	if saved != nil {
		copy(m.caps, saved)
	}
	return false
}

type lookNode struct {
	body    node
	neg     bool
	hasCaps bool
}

func (n *lookNode) match(m *machine, i int, k cont) bool {
	if !m.tick(costLook) {
		return false
	}
	var saved []int
	if n.hasCaps && m.noCap == 0 {
		saved = append([]int(nil), m.caps...)
	}
	matched := n.body.match(m, i, func(int) bool { return true })
	if m.err != nil {
		return false
	}
	if n.neg {
		if saved != nil {
			copy(m.caps, saved)
		}
		if matched {
			return false
		}
		return k(i)
	}
	if !matched {
		return false
	}
	if k(i) {
		return true
	}
	if saved != nil {
		copy(m.caps, saved)
	}
	return false
}

// repSetNode repeats a single-byte matcher.
type repSetNode struct {
	set      charSet
	min, max int
	lazy     bool
}

func (n *repSetNode) run(m *machine, i int) int {
	lim := len(m.s)
	if n.max >= 0 && i+n.max < lim {
		lim = i + n.max
	}
	j := i
	for j < lim && n.set.has(m.s[j]) {
		j++
	}
	return j - i
}

func (n *repSetNode) match(m *machine, i int, k cont) bool {
	c := n.run(m, i)
	if c < n.min {
		return false
	}
	if n.lazy {
		for x := n.min; x <= c; x++ {
			if !m.tick(costSetPos) {
				return false
			}
			if k(i + x) {
				return true
			}
			if m.err != nil {
				return false
			}
		}
		return false
	}
	for x := c; x >= n.min; x-- {
		if !m.tick(costSetPos) {
			return false
		}
		if k(i + x) {
			return true
		}
		if m.err != nil {
			return false
		}
	}
	return false
}

// possSetNode is a possessive single-byte repeat; it is deterministic.
type possSetNode struct{ repSetNode }

func (n *possSetNode) step(m *machine, i int) (int, bool) {
	m.cost += costPoss
	c := n.run(m, i)
	if c < n.min {
		return 0, false
	}
	return i + c, true
}

func (n *possSetNode) match(m *machine, i int, k cont) bool {
	j, ok := n.step(m, i)
	return ok && k(j)
}

type repNode struct {
	body     node
	min, max int
	lazy     bool
}

func (n *repNode) match(m *machine, i int, k cont) bool { return n.iter(m, i, 0, k) }

func (n *repNode) iter(m *machine, i, count int, k cont) bool {
	if !m.tick(costRep) {
		return false
	}
	if !m.push(stackRep) {
		m.stack -= stackRep
		return false
	}
	defer func() { m.stack -= stackRep }()
	if n.lazy {
		if count >= n.min {
			if k(i) {
				return true
			}
			if m.err != nil {
				return false
			}
		}
		if n.max >= 0 && count >= n.max {
			return false
		}
		return n.body.match(m, i, func(j int) bool {
			if j == i && count >= n.min {
				return false // empty iteration adds nothing after k(i) failed
			}
			return n.iter(m, j, count+1, k)
		})
	}
	if n.max < 0 || count < n.max {
		if n.body.match(m, i, func(j int) bool {
			if j == i && count+1 >= n.min {
				// an empty iteration ends an unbounded loop
				return k(j)
			}
			return n.iter(m, j, count+1, k)
		}) {
			return true
		}
		if m.err != nil {
			return false
		}
	}
	if count >= n.min {
		return k(i)
	}
	return false
}

// possRepNode is a possessive group repeat (X)*+. Each iteration is
// committed as soon as it succeeds, like PCRE2's OP_BRAPOS, so earlier
// iterations leave no backtracking state behind. If the minimum cannot be
// reached that way it falls back to the general atomic(repeat) semantics.
type possRepNode struct {
	body     node
	min, max int
	hasCaps  bool
	generic  *atomicNode
}

func (n *possRepNode) match(m *machine, i int, k cont) bool {
	if !m.tick(costAtomic) {
		return false
	}
	if !m.push(stackRep) {
		m.stack -= stackRep
		return false
	}
	var saved []int
	if n.hasCaps && m.noCap == 0 {
		saved = append([]int(nil), m.caps...)
	}
	pos, count := i, 0
	for {
		if !m.tick(costRep) {
			m.stack -= stackRep
			return false
		}
		if n.max >= 0 && count >= n.max {
			break
		}
		end := -1
		if !n.body.match(m, pos, func(j int) bool { end = j; return true }) {
			break
		}
		count++
		if end == pos {
			break
		}
		pos = end
	}
	if m.err != nil {
		m.stack -= stackRep
		return false
	}
	m.stack -= stackRep
	if count < n.min {
		if saved != nil {
			copy(m.caps, saved)
		}
		return n.generic.body.match(m, i, func(j int) bool { pos = j; return true }) && m.err == nil && n.finish(m, pos, k, saved)
	}
	return n.finish(m, pos, k, saved)
}

func (n *possRepNode) finish(m *machine, pos int, k cont, saved []int) bool {
	if k(pos) {
		return true
	}
	if saved != nil {
		copy(m.caps, saved)
	}
	return false
}

// ---- top level ----

// Regexp is a compiled delimited pattern.
type Regexp struct {
	pattern       string
	prog          node
	ngroups       int
	names         []string // by group number, "" when unnamed
	anchored      bool
	forceAnchored bool
}

// Pattern returns the source pattern.
func (re *Regexp) Pattern() string { return re.pattern }

// NumGroups returns the number of capturing groups.
func (re *Regexp) NumGroups() int { return re.ngroups }

// GroupNames returns the group names indexed by number ("" if unnamed).
func (re *Regexp) GroupNames() []string { return append([]string(nil), re.names...) }

type execOpts struct {
	notEmptyAtStart bool
	anchored        bool
}

// measure runs a full match attempt without limits and reports its cost and
// peak stack usage (used to calibrate against PHP).
func (re *Regexp) measure(subject string) (cost int64, maxStack int, matched bool) {
	m := &machine{s: subject, caps: make([]int, 2*(re.ngroups+1)), unlimited: true}
	last := len(subject)
	if re.anchored {
		last = 0
	}
	for st := 0; st <= last && !matched; st++ {
		for i := range m.caps {
			m.caps[i] = -1
		}
		matched = re.prog.match(m, st, func(int) bool { return true })
	}
	return m.cost, m.maxSt, matched
}

// exec finds the first match starting at or after offset.
func (re *Regexp) exec(subject string, offset int, o execOpts) ([]int, error) {
	m := &machine{s: subject, caps: make([]int, 2*(re.ngroups+1))}
	last := len(subject)
	if o.anchored || re.forceAnchored {
		last = offset
	} else if re.anchored {
		if offset > 0 {
			return nil, nil
		}
		last = 0
	}
	for st := offset; st <= last; st++ {
		for i := range m.caps {
			m.caps[i] = -1
		}
		end := -1
		ok := re.prog.match(m, st, func(j int) bool {
			if o.notEmptyAtStart && j == st {
				return false
			}
			end = j
			return true
		})
		if m.err != nil {
			return nil, m.err
		}
		if ok {
			m.caps[0], m.caps[1] = st, end
			return m.caps, nil
		}
	}
	return nil, nil
}
