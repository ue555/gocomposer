package semver

import (
	"strings"
	"sync"
)

// Operator values follow Composer\Semver\Constraint\Constraint::OP_*.
const (
	OpEQ = iota
	OpLT
	OpLE
	OpGT
	OpGE
	OpNE
)

var opFromString = map[string]int{
	"=": OpEQ, "==": OpEQ, "<": OpLT, "<=": OpLE, ">": OpGT, ">=": OpGE, "<>": OpNE, "!=": OpNE,
}

var opToString = [...]string{OpEQ: "==", OpLT: "<", OpLE: "<=", OpGT: ">", OpGE: ">=", OpNE: "!="}

// PHPIntMax is PHP_INT_MAX on 64-bit builds, used by Bound::positiveInfinity.
const PHPIntMax = "9223372036854775807"

// Constraint is Composer's ConstraintInterface.
type Constraint interface {
	Matches(provider Constraint) bool
	String() string
	PrettyString() string
	SetPrettyString(s string)
	LowerBound() Bound
	UpperBound() Bound
}

// Single is Composer\Semver\Constraint\Constraint.
type Single struct {
	op      int
	version string
	pretty  string
}

// NewConstraint builds a single constraint; it panics on an unknown operator
// because all callers pass literal operators.
func NewConstraint(op, version string) *Single {
	o, ok := opFromString[op]
	if !ok {
		panic("invalid operator " + op)
	}
	return &Single{op: o, version: version}
}

func (c *Single) Version() string  { return c.version }
func (c *Single) Operator() string { return opToString[c.op] }
func (c *Single) OpInt() int       { return c.op }

func (c *Single) String() string { return opToString[c.op] + " " + c.version }

func (c *Single) PrettyString() string {
	if c.pretty != "" {
		return c.pretty
	}
	return c.String()
}

func (c *Single) SetPrettyString(s string) { c.pretty = s }

func (c *Single) Matches(provider Constraint) bool {
	if p, ok := provider.(*Single); ok {
		return c.MatchSpecific(p, false)
	}
	return provider.Matches(c)
}

func isBranch(v string) bool { return strings.HasPrefix(v, "dev-") }

func versionCompareOp(a, b string, op string, compareBranches bool) bool {
	aIsBranch := isBranch(a)
	bIsBranch := isBranch(b)
	if op == "!=" && (aIsBranch || bIsBranch) {
		return a != b
	}
	if aIsBranch && bIsBranch {
		return op == "==" && a == b
	}
	if !compareBranches && (aIsBranch || bIsBranch) {
		return false
	}
	return VersionCompareOp(a, b, op)
}

func noEq(op int) string { return strings.ReplaceAll(opToString[op], "=", "") }

// MatchSpecific ports Constraint::matchSpecific.
func (c *Single) MatchSpecific(provider *Single, compareBranches bool) bool {
	noEqualOp := noEq(c.op)
	providerNoEqualOp := noEq(provider.op)
	isEqualOp := c.op == OpEQ
	isNonEqualOp := c.op == OpNE
	isProviderEqualOp := provider.op == OpEQ
	isProviderNonEqualOp := provider.op == OpNE

	if isNonEqualOp || isProviderNonEqualOp {
		if isNonEqualOp && !isProviderNonEqualOp && !isProviderEqualOp && isBranch(provider.version) {
			return false
		}
		if isProviderNonEqualOp && !isNonEqualOp && !isEqualOp && isBranch(c.version) {
			return false
		}
		if !isEqualOp && !isProviderEqualOp {
			return true
		}
		return versionCompareOp(provider.version, c.version, "!=", compareBranches)
	}

	if c.op != OpEQ && noEqualOp == providerNoEqualOp {
		return !(isBranch(c.version) || isBranch(provider.version))
	}

	version1, version2, operator := provider.version, c.version, c.op
	if isEqualOp {
		version1, version2, operator = c.version, provider.version, provider.op
	}
	if versionCompareOp(version1, version2, opToString[operator], compareBranches) {
		return !(opToString[provider.op] == providerNoEqualOp &&
			opToString[c.op] != noEqualOp &&
			VersionCompare(provider.version, c.version) == 0)
	}
	return false
}

func (c *Single) LowerBound() Bound {
	if isBranch(c.version) {
		return ZeroBound()
	}
	switch c.op {
	case OpEQ, OpGE:
		return Bound{c.version, true}
	case OpGT:
		return Bound{c.version, false}
	}
	return ZeroBound()
}

func (c *Single) UpperBound() Bound {
	if isBranch(c.version) {
		return InfinityBound()
	}
	switch c.op {
	case OpEQ, OpLE:
		return Bound{c.version, true}
	case OpLT:
		return Bound{c.version, false}
	}
	return InfinityBound()
}

// Multi is Composer\Semver\Constraint\MultiConstraint.
type Multi struct {
	constraints []Constraint
	conjunctive bool
	pretty      string
	str         string
}

// NewMulti constructs a MultiConstraint; at least two constraints are required.
func NewMulti(cs []Constraint, conjunctive bool) *Multi {
	if len(cs) < 2 {
		panic("Must provide at least two constraints for a MultiConstraint")
	}
	return &Multi{constraints: cs, conjunctive: conjunctive}
}

func (m *Multi) Constraints() []Constraint { return m.constraints }
func (m *Multi) IsConjunctive() bool       { return m.conjunctive }
func (m *Multi) IsDisjunctive() bool       { return !m.conjunctive }

func (m *Multi) Matches(provider Constraint) bool {
	if !m.conjunctive {
		for _, c := range m.constraints {
			if provider.Matches(c) {
				return true
			}
		}
		return false
	}
	if pm, ok := provider.(*Multi); ok && pm.IsDisjunctive() {
		return provider.Matches(m)
	}
	for _, c := range m.constraints {
		if !provider.Matches(c) {
			return false
		}
	}
	return true
}

func (m *Multi) String() string {
	if m.str != "" {
		return m.str
	}
	parts := make([]string, len(m.constraints))
	for i, c := range m.constraints {
		parts[i] = c.String()
	}
	sep := " || "
	if m.conjunctive {
		sep = " "
	}
	m.str = "[" + strings.Join(parts, sep) + "]"
	return m.str
}

func (m *Multi) PrettyString() string {
	if m.pretty != "" {
		return m.pretty
	}
	return m.String()
}

func (m *Multi) SetPrettyString(s string) { m.pretty = s }

func (m *Multi) bounds() (Bound, Bound) {
	var lower, upper *Bound
	for _, c := range m.constraints {
		if lower == nil || upper == nil {
			l, u := c.LowerBound(), c.UpperBound()
			lower, upper = &l, &u
			continue
		}
		op1, op2 := "<", ">"
		if m.conjunctive {
			op1, op2 = ">", "<"
		}
		if l := c.LowerBound(); l.CompareTo(*lower, op1) {
			lower = &l
		}
		if u := c.UpperBound(); u.CompareTo(*upper, op2) {
			upper = &u
		}
	}
	return *lower, *upper
}

func (m *Multi) LowerBound() Bound { l, _ := m.bounds(); return l }
func (m *Multi) UpperBound() Bound { _, u := m.bounds(); return u }

// CreateMulti ports MultiConstraint::create including its OR-group merge.
func CreateMulti(cs []Constraint, conjunctive bool) Constraint {
	if len(cs) == 0 {
		return &MatchAll{}
	}
	if len(cs) == 1 {
		return cs[0]
	}
	if !conjunctive {
		left := cs[0]
		var merged []Constraint
		optimized := false
		for i := 1; i < len(cs); i++ {
			right := cs[i]
			lm, lok := left.(*Multi)
			rm, rok := right.(*Multi)
			if lok && lm.conjunctive && rok && rm.conjunctive && len(lm.constraints) == 2 && len(rm.constraints) == 2 {
				left0 := lm.constraints[0].String()
				left1 := lm.constraints[1].String()
				right0 := rm.constraints[0].String()
				right1 := rm.constraints[1].String()
				if strings.HasPrefix(left0, ">=") && strings.HasPrefix(left1, "<") &&
					strings.HasPrefix(right0, ">=") && strings.HasPrefix(right1, "<") &&
					safeSub(left1, 2) == safeSub(right0, 3) {
					optimized = true
					left = NewMulti([]Constraint{lm.constraints[0], rm.constraints[1]}, true)
					continue
				}
			}
			merged = append(merged, left)
			left = right
		}
		if optimized {
			merged = append(merged, left)
			if len(merged) == 1 {
				return merged[0]
			}
			return NewMulti(merged, false)
		}
	}
	return NewMulti(cs, conjunctive)
}

func safeSub(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[n:]
}

// MatchAll is MatchAllConstraint.
type MatchAll struct{ pretty string }

func (m *MatchAll) Matches(Constraint) bool { return true }
func (m *MatchAll) String() string          { return "*" }
func (m *MatchAll) PrettyString() string {
	if m.pretty != "" {
		return m.pretty
	}
	return "*"
}
func (m *MatchAll) SetPrettyString(s string) { m.pretty = s }
func (m *MatchAll) LowerBound() Bound        { return ZeroBound() }
func (m *MatchAll) UpperBound() Bound        { return InfinityBound() }

// MatchNone is MatchNoneConstraint.
type MatchNone struct{ pretty string }

func (m *MatchNone) Matches(Constraint) bool { return false }
func (m *MatchNone) String() string          { return "[]" }
func (m *MatchNone) PrettyString() string {
	if m.pretty != "" {
		return m.pretty
	}
	return "[]"
}
func (m *MatchNone) SetPrettyString(s string) { m.pretty = s }
func (m *MatchNone) LowerBound() Bound        { return Bound{"0.0.0.0-dev", false} }
func (m *MatchNone) UpperBound() Bound        { return Bound{"0.0.0.0-dev", false} }

// Bound is Composer\Semver\Constraint\Bound.
type Bound struct {
	Version   string
	Inclusive bool
}

func ZeroBound() Bound     { return Bound{"0.0.0.0-dev", true} }
func InfinityBound() Bound { return Bound{PHPIntMax + ".0.0.0", false} }

func (b Bound) IsZero() bool             { return b.Version == "0.0.0.0-dev" && b.Inclusive }
func (b Bound) IsPositiveInfinity() bool { return b.Version == PHPIntMax+".0.0.0" && !b.Inclusive }

func (b Bound) CompareTo(other Bound, op string) bool {
	if b == other {
		return false
	}
	c := VersionCompare(b.Version, other.Version)
	if c != 0 {
		if op == ">" {
			return c == 1
		}
		return c == -1
	}
	if op == ">" {
		return other.Inclusive
	}
	return !other.Inclusive
}

func (b Bound) String() string {
	if b.Inclusive {
		return b.Version + " [inclusive]"
	}
	return b.Version + " [exclusive]"
}

// matchCache memoizes CompilingMatcher::match results per constraint string.
// The solver evaluates the same pairs many times.
type matchKey struct {
	constraint string
	version    string
}

var (
	matchMu    sync.Mutex
	matchCache = map[matchKey]bool{}
)

// MatchVersion is CompilingMatcher::match($constraint, OP_EQ, $version).
// The resolver is single-threaded per process; the cache is not shared across
// goroutines.
func MatchVersion(c Constraint, version string) bool {
	if c == nil {
		return true
	}
	if _, ok := c.(*MatchAll); ok {
		return true
	}
	key := matchKey{c.String(), version}
	matchMu.Lock()
	v, ok := matchCache[key]
	matchMu.Unlock()
	if ok {
		return v
	}
	r := c.Matches(&Single{op: OpEQ, version: version})
	matchMu.Lock()
	if len(matchCache) > 1<<20 {
		matchCache = map[matchKey]bool{}
	}
	matchCache[key] = r
	matchMu.Unlock()
	return r
}

// ClearMatchCache drops memoized results (CompilingMatcher::clear).
func ClearMatchCache() {
	matchMu.Lock()
	matchCache = map[matchKey]bool{}
	matchMu.Unlock()
}

// Compare is Comparator::compare.
func Compare(v1, op, v2 string) bool {
	return NewConstraint(op, v2).MatchSpecific(NewConstraint("==", v1), true)
}
