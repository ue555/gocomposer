package semver

import (
	"sort"
	"sync"
)

// Port of composer/semver src/Intervals.php and Interval.php.

// Interval is a numeric range with constraint endpoints.
type Interval struct {
	Start *Single
	End   *Single
}

var (
	fromZero              = NewConstraint(">=", "0.0.0.0-dev")
	untilPositiveInfinity = NewConstraint("<", PHPIntMax+".0.0.0")
)

// branchSet mirrors the PHP ['names' => [...], 'exclude' => bool] arrays.
// keys hold the PHP array keys so that index-based comparisons in isSubsetOf
// behave like the original (array_intersect/array_diff preserve keys).
type branchSet struct {
	keys    []int
	names   []string
	exclude bool
}

func anyDev() branchSet { return branchSet{exclude: true} }
func noDev() branchSet  { return branchSet{exclude: false} }

func (b branchSet) contains(name string) bool {
	for _, n := range b.names {
		if n == name {
			return true
		}
	}
	return false
}

// intersect is array_intersect($a, $b) keeping $a's keys.
func intersect(a, b branchSet) ([]int, []string) {
	var keys []int
	var names []string
	for i, n := range a.names {
		if b.contains(n) {
			keys = append(keys, a.keys[i])
			names = append(names, n)
		}
	}
	return keys, names
}

// diff is array_diff($a, $b) keeping $a's keys.
func diff(a, b branchSet) ([]int, []string) {
	var keys []int
	var names []string
	for i, n := range a.names {
		if !b.contains(n) {
			keys = append(keys, a.keys[i])
			names = append(names, n)
		}
	}
	return keys, names
}

// merge is array_merge($a, $b) of two lists (renumbers keys).
func merge(a, b branchSet) ([]int, []string) {
	names := append(append([]string(nil), a.names...), b.names...)
	keys := make([]int, len(names))
	for i := range keys {
		keys[i] = i
	}
	return keys, names
}

func (b *branchSet) unique() {
	seen := map[string]bool{}
	var keys []int
	var names []string
	for i, n := range b.names {
		if seen[n] {
			continue
		}
		seen[n] = true
		keys = append(keys, b.keys[i])
		names = append(names, n)
	}
	b.keys, b.names = keys, names
}

type intervalsResult struct {
	numeric  []Interval
	branches branchSet
}

var (
	intervalsMu    sync.Mutex
	intervalsCache = map[string]intervalsResult{}
)

// ClearIntervals is Intervals::clear().
func ClearIntervals() {
	intervalsMu.Lock()
	intervalsCache = map[string]intervalsResult{}
	intervalsMu.Unlock()
}

func getIntervals(c Constraint) intervalsResult {
	key := c.String()
	intervalsMu.Lock()
	r, ok := intervalsCache[key]
	intervalsMu.Unlock()
	if ok {
		return r
	}
	r = generateIntervals(c, false)
	intervalsMu.Lock()
	intervalsCache[key] = r
	intervalsMu.Unlock()
	return r
}

// IsSubsetOf is Intervals::isSubsetOf.
func IsSubsetOf(candidate, constraint Constraint) bool {
	if _, ok := constraint.(*MatchAll); ok {
		return true
	}
	if _, ok := candidate.(*MatchNone); ok {
		return false
	}
	if _, ok := constraint.(*MatchNone); ok {
		return false
	}
	inter := getIntervals(NewMulti([]Constraint{candidate, constraint}, true))
	cand := getIntervals(candidate)
	if len(inter.numeric) != len(cand.numeric) {
		return false
	}
	for i, iv := range inter.numeric {
		if i >= len(cand.numeric) {
			return false
		}
		if cand.numeric[i].Start.String() != iv.Start.String() || cand.numeric[i].End.String() != iv.End.String() {
			return false
		}
	}
	if inter.branches.exclude != cand.branches.exclude {
		return false
	}
	if len(inter.branches.names) != len(cand.branches.names) {
		return false
	}
	for i, name := range inter.branches.names {
		key := inter.branches.keys[i]
		found := false
		for j, k := range cand.branches.keys {
			if k == key {
				found = true
				if cand.branches.names[j] != name {
					return false
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// HaveIntersections is Intervals::haveIntersections.
func HaveIntersections(a, b Constraint) bool {
	if _, ok := a.(*MatchAll); ok {
		return true
	}
	if _, ok := b.(*MatchAll); ok {
		return true
	}
	if _, ok := a.(*MatchNone); ok {
		return false
	}
	if _, ok := b.(*MatchNone); ok {
		return false
	}
	r := generateIntervals(NewMulti([]Constraint{a, b}, true), true)
	return len(r.numeric) > 0 || r.branches.exclude || len(r.branches.names) > 0
}

// CompactConstraint is Intervals::compactConstraint.
func CompactConstraint(c Constraint) Constraint {
	if _, ok := c.(*Multi); !ok {
		return c
	}
	iv := generateIntervals(c, false)
	var constraints []Constraint
	hasNumericMatchAll := false
	if len(iv.numeric) == 1 && iv.numeric[0].Start.String() == fromZero.String() && iv.numeric[0].End.String() == untilPositiveInfinity.String() {
		constraints = append(constraints, iv.numeric[0].Start)
		hasNumericMatchAll = true
	} else {
		numeric := iv.numeric
		count := len(numeric)
		skipNotEqual := false
		if !iv.branches.exclude && count > 1 &&
			numeric[0].Start.String() == fromZero.String() &&
			numeric[count-1].End.String() == untilPositiveInfinity.String() {
			skipNotEqual = true
			for i := 0; i < count-1; i++ {
				if numeric[i].End.Operator() != "<" || numeric[i+1].Start.Operator() != ">" || numeric[i].End.Version() != numeric[i+1].Start.Version() {
					skipNotEqual = false
					break
				}
			}
		}
		var unEqual []Constraint
		for i := 0; i < count; i++ {
			interval := numeric[i]
			if !skipNotEqual && interval.End.Operator() == "<" && i+1 < count {
				next := numeric[i+1]
				if interval.End.Version() == next.Start.Version() && next.Start.Operator() == ">" {
					if len(unEqual) == 0 && interval.Start.String() != fromZero.String() {
						unEqual = append(unEqual, interval.Start)
					}
					unEqual = append(unEqual, NewConstraint("!=", interval.End.Version()))
					continue
				}
			}
			if len(unEqual) > 0 {
				if interval.End.String() != untilPositiveInfinity.String() {
					unEqual = append(unEqual, interval.End)
				}
				if len(unEqual) > 1 {
					constraints = append(constraints, NewMulti(unEqual, true))
				} else {
					constraints = append(constraints, unEqual[0])
				}
				unEqual = nil
				continue
			}
			if interval.Start.Version() == interval.End.Version() && interval.Start.Operator() == ">=" && interval.End.Operator() == "<=" {
				constraints = append(constraints, NewConstraint("==", interval.Start.Version()))
				continue
			}
			if interval.Start.String() == fromZero.String() {
				constraints = append(constraints, interval.End)
			} else if interval.End.String() == untilPositiveInfinity.String() {
				constraints = append(constraints, interval.Start)
			} else {
				constraints = append(constraints, NewMulti([]Constraint{interval.Start, interval.End}, true))
			}
		}
	}

	var devConstraints []Constraint
	if len(iv.branches.names) == 0 {
		if iv.branches.exclude && hasNumericMatchAll {
			return &MatchAll{}
		}
	} else {
		for _, name := range iv.branches.names {
			if iv.branches.exclude {
				devConstraints = append(devConstraints, NewConstraint("!=", name))
			} else {
				devConstraints = append(devConstraints, NewConstraint("==", name))
			}
		}
		if iv.branches.exclude {
			if len(constraints) > 1 {
				return NewMulti(append([]Constraint{NewMulti(constraints, false)}, devConstraints...), true)
			}
			if len(constraints) == 1 && constraints[0].String() == fromZero.String() {
				if len(devConstraints) > 1 {
					return NewMulti(devConstraints, true)
				}
				return devConstraints[0]
			}
			return NewMulti(append(append([]Constraint(nil), constraints...), devConstraints...), true)
		}
		constraints = append(constraints, devConstraints...)
	}
	if len(constraints) > 1 {
		return NewMulti(constraints, false)
	}
	if len(constraints) == 1 {
		return constraints[0]
	}
	return &MatchNone{}
}

type border struct {
	version  string
	operator string
	start    bool
}

var opSortOrder = map[string]int{">=": -3, "<": -2, ">": 2, "<=": 3}

func generateIntervals(c Constraint, stopOnFirstValid bool) intervalsResult {
	switch t := c.(type) {
	case *MatchAll:
		return intervalsResult{numeric: []Interval{{fromZero, untilPositiveInfinity}}, branches: anyDev()}
	case *MatchNone:
		return intervalsResult{branches: branchSet{exclude: false}}
	case *Single:
		return singleIntervals(t)
	case *Multi:
		return multiIntervals(t, stopOnFirstValid)
	}
	panic("unexpected constraint type")
}

func multiIntervals(c *Multi, stopOnFirstValid bool) intervalsResult {
	var numericGroups [][]Interval
	var constraintBranches []branchSet
	for _, sub := range c.constraints {
		r := getIntervals(sub)
		numericGroups = append(numericGroups, r.numeric)
		constraintBranches = append(constraintBranches, r.branches)
	}
	var branches branchSet
	if c.IsDisjunctive() {
		branches = noDev()
		for _, b := range constraintBranches {
			if b.exclude {
				if branches.exclude {
					branches.keys, branches.names = intersect(branches, b)
				} else {
					branches.exclude = true
					branches.keys, branches.names = diff(b, branches)
				}
			} else {
				if branches.exclude {
					branches.keys, branches.names = diff(branches, b)
				} else {
					branches.keys, branches.names = merge(branches, b)
				}
			}
		}
	} else {
		branches = anyDev()
		for _, b := range constraintBranches {
			if b.exclude {
				if branches.exclude {
					branches.keys, branches.names = merge(branches, b)
				} else {
					branches.keys, branches.names = diff(branches, b)
				}
			} else {
				if branches.exclude {
					branches.keys, branches.names = diff(b, branches)
					branches.exclude = false
				} else {
					branches.keys, branches.names = intersect(branches, b)
				}
			}
		}
	}
	branches.unique()

	if len(numericGroups) == 1 {
		return intervalsResult{numeric: numericGroups[0], branches: branches}
	}

	var borders []border
	for _, group := range numericGroups {
		for _, iv := range group {
			borders = append(borders, border{iv.Start.Version(), iv.Start.Operator(), true})
			borders = append(borders, border{iv.End.Version(), iv.End.Operator(), false})
		}
	}
	sort.SliceStable(borders, func(i, j int) bool {
		order := VersionCompare(borders[i].version, borders[j].version)
		if order == 0 {
			return opSortOrder[borders[i].operator] < opSortOrder[borders[j].operator]
		}
		return order < 0
	})

	active := 0
	var intervals []Interval
	threshold := 1
	if c.IsConjunctive() {
		threshold = len(numericGroups)
	}
	var start *Single
	for _, b := range borders {
		if b.start {
			active++
		} else {
			active--
		}
		if start == nil && active >= threshold {
			start = NewConstraint(b.operator, b.version)
		} else if start != nil && active < threshold {
			if VersionCompare(start.Version(), b.version) == 0 &&
				((start.Operator() == ">" && b.operator == "<=") || (start.Operator() == ">=" && b.operator == "<")) {
				// invalid interval, dropped
			} else {
				intervals = append(intervals, Interval{start, NewConstraint(b.operator, b.version)})
				if stopOnFirstValid {
					break
				}
			}
			start = nil
		}
	}
	return intervalsResult{numeric: intervals, branches: branches}
}

func singleIntervals(c *Single) intervalsResult {
	op := c.Operator()
	if isBranch(c.Version()) {
		var intervals []Interval
		branches := branchSet{exclude: false}
		if op == "!=" {
			intervals = append(intervals, Interval{fromZero, untilPositiveInfinity})
			branches = branchSet{keys: []int{0}, names: []string{c.Version()}, exclude: true}
		} else if op == "==" {
			branches.keys = []int{0}
			branches.names = []string{c.Version()}
		}
		return intervalsResult{numeric: intervals, branches: branches}
	}
	if op[0] == '>' {
		return intervalsResult{numeric: []Interval{{c, untilPositiveInfinity}}, branches: noDev()}
	}
	if op[0] == '<' {
		if op == "<" && c.Version() == "0.0.0.0-dev" {
			return intervalsResult{branches: noDev()}
		}
		return intervalsResult{numeric: []Interval{{fromZero, c}}, branches: noDev()}
	}
	if op == "!=" {
		return intervalsResult{numeric: []Interval{
			{fromZero, NewConstraint("<", c.Version())},
			{NewConstraint(">", c.Version()), untilPositiveInfinity},
		}, branches: anyDev()}
	}
	return intervalsResult{numeric: []Interval{{NewConstraint(">=", c.Version()), NewConstraint("<=", c.Version())}}, branches: noDev()}
}

// LastIntervalEnd returns the end of the last numeric interval of c
// (end(Intervals::get($c)['numeric'])->getEnd()).
func LastIntervalEnd(c Constraint) (*Single, bool) {
	r := getIntervals(c)
	if len(r.numeric) == 0 {
		return nil, false
	}
	return r.numeric[len(r.numeric)-1].End, true
}
