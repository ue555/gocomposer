package solver

import (
	"fmt"
	"time"

	"gocomposer/internal/policy"
)

// decision is one entry of the decision queue.
type decision struct {
	literal int
	reason  *Rule
}

// Decisions is Composer\DependencyResolver\Decisions.
type Decisions struct {
	pool  *Pool
	dmap  map[int]int
	queue []decision
}

func newDecisions(pool *Pool) *Decisions { return &Decisions{pool: pool, dmap: map[int]int{}} }

// SolverBugError mirrors SolverBugException.
type SolverBugError struct{ Msg string }

func (e *SolverBugError) Error() string {
	return e.Msg + "\nThis exception was most likely caused by a bug in Composer.\nPlease report the command you ran, the exact error you received, and your composer.json on https://github.com/composer/composer/issues - thank you!\n"
}

func (d *Decisions) decide(lit, level int, why *Rule) {
	id := abs(lit)
	if prev := d.dmap[id]; prev != 0 {
		panic(&SolverBugError{fmt.Sprintf("Trying to decide %s on level %d, even though %s was previously decided as %d.", d.pool.LiteralToPrettyString(lit, nil), level, d.pool.LiteralToPackage(lit), prev)})
	}
	if lit > 0 {
		d.dmap[id] = level
	} else {
		d.dmap[id] = -level
	}
	d.queue = append(d.queue, decision{lit, why})
}

func (d *Decisions) satisfy(lit int) bool {
	v := d.dmap[abs(lit)]
	return (lit > 0 && v > 0) || (lit < 0 && v < 0)
}

func (d *Decisions) conflict(lit int) bool {
	v := d.dmap[abs(lit)]
	return (v > 0 && lit < 0) || (v < 0 && lit > 0)
}

func (d *Decisions) decided(lit int) bool        { return d.dmap[abs(lit)] != 0 }
func (d *Decisions) undecided(lit int) bool      { return d.dmap[abs(lit)] == 0 }
func (d *Decisions) decidedInstall(lit int) bool { return d.dmap[abs(lit)] > 0 }

func (d *Decisions) decisionLevel(lit int) int { return abs(d.dmap[abs(lit)]) }

func (d *Decisions) decisionRule(lit int) *Rule {
	id := abs(lit)
	for _, dec := range d.queue {
		if abs(dec.literal) == id {
			return dec.reason
		}
	}
	panic("Did not find a decision rule using " + fmt.Sprint(lit))
}

func (d *Decisions) count() int { return len(d.queue) }

func (d *Decisions) resetToOffset(offset int) {
	for len(d.queue) > offset+1 {
		last := d.queue[len(d.queue)-1]
		d.queue = d.queue[:len(d.queue)-1]
		d.dmap[abs(last.literal)] = 0
	}
}

func (d *Decisions) revertLast() {
	last := d.queue[len(d.queue)-1]
	d.dmap[abs(last.literal)] = 0
	d.queue = d.queue[:len(d.queue)-1]
}

// watchNode is RuleWatchNode.
type watchNode struct {
	watch1, watch2 int
	rule           *Rule
}

func newWatchNode(r *Rule) *watchNode {
	n := &watchNode{rule: r}
	if len(r.literals) > 0 {
		n.watch1 = r.literals[0]
	}
	if len(r.literals) > 1 {
		n.watch2 = r.literals[1]
	}
	return n
}

func (n *watchNode) watch2OnHighest(d *Decisions) {
	if len(n.rule.literals) < 3 || n.rule.IsMulti() {
		return
	}
	watchLevel := 0
	for _, lit := range n.rule.literals {
		if level := d.decisionLevel(lit); level > watchLevel {
			n.watch2 = lit
			watchLevel = level
		}
	}
}

func (n *watchNode) otherWatch(lit int) int {
	if n.watch1 == lit {
		return n.watch2
	}
	return n.watch1
}

func (n *watchNode) moveWatch(from, to int) {
	if n.watch1 == from {
		n.watch1 = to
	} else {
		n.watch2 = to
	}
}

// watchGraph is RuleWatchGraph. Each chain is stored reversed: the PHP
// SplDoublyLinkedList::unshift() becomes an append and iteration runs from
// the end of the slice.
type watchGraph struct {
	chains map[int][]*watchNode
}

func newWatchGraph() *watchGraph { return &watchGraph{chains: map[int][]*watchNode{}} }

func (g *watchGraph) insert(n *watchNode) {
	if n.rule.IsAssertion() {
		return
	}
	if !n.rule.IsMulti() {
		g.chains[n.watch1] = append(g.chains[n.watch1], n)
		g.chains[n.watch2] = append(g.chains[n.watch2], n)
	} else {
		for _, lit := range n.rule.literals {
			g.chains[lit] = append(g.chains[lit], n)
		}
	}
}

func (g *watchGraph) propagateLiteral(decided, level int, d *Decisions) *Rule {
	lit := -decided
	if _, ok := g.chains[lit]; !ok {
		return nil
	}
	i := len(g.chains[lit]) - 1
	for i >= 0 {
		chain := g.chains[lit]
		node := chain[i]
		if !node.rule.IsMulti() {
			other := node.otherWatch(lit)
			if !node.rule.disabled && !d.satisfy(other) {
				alternative := 0
				found := false
				for _, rl := range node.rule.literals {
					if rl != lit && rl != other && !d.conflict(rl) {
						alternative = rl
						found = true
						break
					}
				}
				if found {
					// moveWatch: remove the node from this chain and unshift it on the other
					node.moveWatch(lit, alternative)
					g.chains[lit] = append(chain[:i:i], chain[i+1:]...)
					g.chains[alternative] = append(g.chains[alternative], node)
					i--
					continue
				}
				if d.conflict(other) {
					return node.rule
				}
				d.decide(other, level, node.rule)
			}
		} else {
			for _, other := range node.rule.literals {
				if lit != other && !d.satisfy(other) {
					if d.conflict(other) {
						return node.rule
					}
					d.decide(other, level, node.rule)
				}
			}
		}
		i--
	}
	return nil
}

type branch struct {
	literals []int
	level    int
}

// Solver is Composer\DependencyResolver\Solver.
type Solver struct {
	policy         *DefaultPolicy
	pool           *Pool
	rules          *RuleSet
	watchGraph     *watchGraph
	decisions      *Decisions
	fixedMap       map[int]bool
	propagateIndex int
	branches       []branch
	problems       []*Problem
	learnedPool    [][]*Rule
	learnedWhy     map[*Rule]int
	Debug          func(string)
	// Log writes to the error output at a verbosity (IOInterface::writeError).
	Log func(msg string, verbosity int)
}

func NewSolver(policy *DefaultPolicy, pool *Pool) *Solver {
	return &Solver{policy: policy, pool: pool, learnedWhy: map[*Rule]int{}}
}

func (s *Solver) RuleSetSize() int { return s.rules.Count() }
func (s *Solver) Pool() *Pool      { return s.pool }

// ProblemsError is SolverProblemsException.
type ProblemsError struct {
	Problems    []*Problem
	LearnedPool [][]*Rule
}

func (e *ProblemsError) Error() string {
	return fmt.Sprintf("Failed resolving dependencies with %d problems, call getPrettyString to get formatted details", len(e.Problems))
}

func (s *Solver) makeAssertionRuleDecisions() {
	decisionStart := s.decisions.count() - 1
	rulesCount := s.rules.Count()
	for ruleIndex := 0; ruleIndex < rulesCount; ruleIndex++ {
		rule := s.rules.RuleByID[ruleIndex]
		if !rule.IsAssertion() || rule.disabled {
			continue
		}
		lit := rule.literals[0]
		if !s.decisions.decided(lit) {
			s.decisions.decide(lit, 1, rule)
			continue
		}
		if s.decisions.satisfy(lit) {
			continue
		}
		if rule.Type == TypeLearned {
			rule.Disable()
			continue
		}
		conflict := s.decisions.decisionRule(lit)
		if conflict.Type == TypePackage {
			problem := &Problem{}
			problem.addRule(rule)
			problem.addRule(conflict)
			rule.Disable()
			s.problems = append(s.problems, problem)
			continue
		}
		problem := &Problem{}
		problem.addRule(rule)
		problem.addRule(conflict)
		for _, assertRule := range s.rules.OfType(TypeRequest) {
			if assertRule.disabled || !assertRule.IsAssertion() {
				continue
			}
			if abs(lit) != abs(assertRule.literals[0]) {
				continue
			}
			problem.addRule(assertRule)
			assertRule.Disable()
		}
		s.problems = append(s.problems, problem)
		s.decisions.resetToOffset(decisionStart)
		ruleIndex = -1
	}
}

func (s *Solver) log(msg string, verbosity int) {
	if s.Log != nil {
		s.Log(msg, verbosity)
	}
}

// checkForFilterListRemovedLockedPackages is
// Solver::checkForFilterListRemovedLockedPackages.
func (s *Solver) checkForFilterListRemovedLockedPackages(req *Request) {
	for _, p := range req.LockedPackages() {
		r := s.pool.PolicyRemovalReason(p.Name, p.Version())
		if r == nil || r.Type() != policy.RemovalFilterList {
			continue
		}
		prob := &Problem{}
		prob.addRule(newGenericRule(nil, RuleLockedFilterListRemove, p))
		s.problems = append(s.problems, prob)
	}
}

func (s *Solver) checkForRootRequireProblems(req *Request, filter PlatformFilter) {
	for _, r := range req.Requires() {
		if filter.IsIgnored(r.Name) {
			continue
		}
		c := filterConstraint(filter, r.Name, r.Constraint, true)
		if len(s.pool.WhatProvides(r.Name, c)) == 0 {
			problem := &Problem{}
			problem.addRule(newGenericRule(nil, RuleRootRequire, RootRequireData{r.Name, c}))
			s.problems = append(s.problems, problem)
		}
	}
}

// Solve is Solver::solve. Internal consistency failures surface as
// *SolverBugError values.
func (s *Solver) Solve(req *Request, filter PlatformFilter) (tx *LockTransaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			if be, ok := r.(*SolverBugError); ok {
				err = be
				return
			}
			panic(r)
		}
	}()
	if filter == nil {
		filter = IgnoreNothing{}
	}
	s.fixedMap = map[int]bool{}
	for _, p := range req.FixedPackages() {
		s.fixedMap[p.ID] = true
	}
	s.log("Generating rules", 256)
	s.rules = GenerateRules(s.pool, req, filter)
	s.checkForRootRequireProblems(req, filter)
	s.checkForFilterListRemovedLockedPackages(req)
	s.decisions = newDecisions(s.pool)
	s.watchGraph = newWatchGraph()
	for _, r := range s.rules.Iterate() {
		s.watchGraph.insert(newWatchNode(r))
	}
	s.makeAssertionRuleDecisions()
	s.log("Resolving dependencies through SAT", 256)
	before := time.Now()
	s.runSat()
	s.log("", 256)
	s.log(fmt.Sprintf("Dependency resolution completed in %.3f seconds", time.Since(before).Seconds()), 64)
	if len(s.problems) > 0 {
		return nil, &ProblemsError{Problems: s.problems, LearnedPool: s.learnedPool}
	}
	return newLockTransaction(s.pool, req.PresentList(), req.FixedPackagesMap(), s.decisions)
}

func (s *Solver) propagate(level int) *Rule {
	for s.propagateIndex >= 0 && s.propagateIndex < s.decisions.count() {
		dec := s.decisions.queue[s.propagateIndex]
		conflict := s.watchGraph.propagateLiteral(dec.literal, level, s.decisions)
		s.propagateIndex++
		if conflict != nil {
			return conflict
		}
	}
	return nil
}

func (s *Solver) revert(level int) {
	for s.decisions.count() > 0 {
		lit := s.decisions.queue[len(s.decisions.queue)-1].literal
		if s.decisions.undecided(lit) {
			break
		}
		if s.decisions.decisionLevel(lit) <= level {
			break
		}
		s.decisions.revertLast()
		s.propagateIndex = s.decisions.count()
	}
	for len(s.branches) > 0 && s.branches[len(s.branches)-1].level >= level {
		s.branches = s.branches[:len(s.branches)-1]
	}
}

func (s *Solver) setPropagateLearn(level, lit int, rule *Rule) int {
	level++
	s.decisions.decide(lit, level, rule)
	for {
		rule = s.propagate(level)
		if rule == nil {
			break
		}
		if level == 1 {
			s.analyzeUnsolvable(rule)
			return 0
		}
		learnLiteral, newLevel, newRule, why := s.analyze(level, rule)
		if newLevel <= 0 || newLevel >= level {
			panic(&SolverBugError{fmt.Sprintf("Trying to revert to invalid level %d from level %d.", newLevel, level)})
		}
		level = newLevel
		s.revert(level)
		s.rules.Add(newRule, TypeLearned)
		s.learnedWhy[newRule] = why
		node := newWatchNode(newRule)
		node.watch2OnHighest(s.decisions)
		s.watchGraph.insert(node)
		s.decisions.decide(learnLiteral, level, newRule)
	}
	return level
}

func (s *Solver) selectAndInstall(level int, decisionQueue []int, rule *Rule) int {
	literals := s.policy.SelectPreferredPackages(s.pool, decisionQueue, rule.RequiredPackage())
	selected := literals[0]
	literals = literals[1:]
	if len(literals) > 0 {
		s.branches = append(s.branches, branch{append([]int(nil), literals...), level})
	}
	return s.setPropagateLearn(level, selected, rule)
}

func (s *Solver) analyze(level int, rule *Rule) (int, int, *Rule, int) {
	analyzedRule := rule
	ruleLevel := 1
	num := 0
	l1num := 0
	seen := map[int]bool{}
	learnedLiteral := 0
	hasLearned := false
	var otherLearned []int
	decisionID := s.decisions.count()
	s.learnedPool = append(s.learnedPool, nil)
	last := len(s.learnedPool) - 1

outer:
	for {
		s.learnedPool[last] = append(s.learnedPool[last], rule)
		for _, lit := range rule.literals {
			if rule.IsMulti() && !s.decisions.decided(lit) {
				continue
			}
			if s.decisions.satisfy(lit) {
				continue
			}
			if seen[abs(lit)] {
				continue
			}
			seen[abs(lit)] = true
			l := s.decisions.decisionLevel(lit)
			if l == 1 {
				l1num++
			} else if level == l {
				num++
			} else {
				otherLearned = append(otherLearned, lit)
				if l > ruleLevel {
					ruleLevel = l
				}
			}
		}
		l1retry := true
		var lit int
		for l1retry {
			l1retry = false
			if num == 0 {
				l1num--
				if l1num == 0 {
					break outer
				}
			}
			for {
				if decisionID <= 0 {
					panic(&SolverBugError{fmt.Sprintf("Reached invalid decision id %d while looking through %s for a literal present in the analyzed rule %s.", decisionID, rule, analyzedRule)})
				}
				decisionID--
				lit = s.decisions.queue[decisionID].literal
				if seen[abs(lit)] {
					break
				}
			}
			delete(seen, abs(lit))
			if num != 0 {
				num--
				if num == 0 {
					learnedLiteral = -lit
					hasLearned = true
					if l1num == 0 {
						break outer
					}
					for _, o := range otherLearned {
						delete(seen, abs(o))
					}
					l1num++
					l1retry = true
					continue
				}
			}
			r := s.decisions.queue[decisionID].reason
			if r.IsMulti() {
				for _, rl := range r.literals {
					if !seen[abs(rl)] && s.decisions.satisfy(-rl) {
						s.learnedPool[last] = append(s.learnedPool[last], r)
						l := s.decisions.decisionLevel(rl)
						if l == 1 {
							l1num++
						} else if level == l {
							num++
						} else {
							otherLearned = append(otherLearned, rl)
							if l > ruleLevel {
								ruleLevel = l
							}
						}
						seen[abs(rl)] = true
						break
					}
				}
				l1retry = true
			}
		}
		rule = s.decisions.queue[decisionID].reason
	}
	why := len(s.learnedPool) - 1
	if !hasLearned {
		panic(&SolverBugError{fmt.Sprintf("Did not find a learnable literal in analyzed rule %s.", analyzedRule)})
	}
	lits := append([]int{learnedLiteral}, otherLearned...)
	return learnedLiteral, ruleLevel, newGenericRule(lits, RuleLearned, why), why
}

func (s *Solver) analyzeUnsolvableRule(problem *Problem, conflictRule *Rule, ruleSeen map[*Rule]bool) {
	ruleSeen[conflictRule] = true
	if conflictRule.Type == TypeLearned {
		learnedWhy := s.learnedWhy[conflictRule]
		for _, pr := range s.learnedPool[learnedWhy] {
			if !ruleSeen[pr] {
				s.analyzeUnsolvableRule(problem, pr, ruleSeen)
			}
		}
		return
	}
	if conflictRule.Type == TypePackage {
		return
	}
	problem.nextSection()
	problem.addRule(conflictRule)
}

func (s *Solver) analyzeUnsolvable(conflictRule *Rule) {
	problem := &Problem{}
	problem.addRule(conflictRule)
	ruleSeen := map[*Rule]bool{}
	s.analyzeUnsolvableRule(problem, conflictRule, ruleSeen)
	s.problems = append(s.problems, problem)
	seen := map[int]bool{}
	for _, lit := range conflictRule.literals {
		if s.decisions.satisfy(lit) {
			continue
		}
		seen[abs(lit)] = true
	}
	for i := len(s.decisions.queue) - 1; i >= 0; i-- {
		dec := s.decisions.queue[i]
		if !seen[abs(dec.literal)] {
			continue
		}
		why := dec.reason
		problem.addRule(why)
		s.analyzeUnsolvableRule(problem, why, ruleSeen)
		for _, lit := range why.literals {
			if s.decisions.satisfy(lit) {
				continue
			}
			seen[abs(lit)] = true
		}
	}
}

func (s *Solver) runSat() {
	s.propagateIndex = 0
	level := 1
	systemLevel := level + 1
	for {
		if level == 1 {
			if conflictRule := s.propagate(level); conflictRule != nil {
				s.analyzeUnsolvable(conflictRule)
				return
			}
		}
		if level < systemLevel {
			requestRules := s.rules.OfType(TypeRequest)
			brokeAt := -1
			for idx, rule := range requestRules {
				if !rule.IsEnabled() {
					continue
				}
				var decisionQueue []int
				noneSatisfied := true
				for _, lit := range rule.literals {
					if s.decisions.satisfy(lit) {
						noneSatisfied = false
						break
					}
					if lit > 0 && s.decisions.undecided(lit) {
						decisionQueue = append(decisionQueue, lit)
					}
				}
				if noneSatisfied && len(decisionQueue) > 0 {
					var pruned []int
					for _, lit := range decisionQueue {
						if s.fixedMap[abs(lit)] {
							pruned = append(pruned, lit)
						}
					}
					if len(pruned) > 0 {
						decisionQueue = pruned
					}
				}
				if noneSatisfied && len(decisionQueue) > 0 {
					oLevel := level
					level = s.selectAndInstall(level, decisionQueue, rule)
					if level == 0 {
						return
					}
					if level <= oLevel {
						brokeAt = idx
						break
					}
				}
			}
			systemLevel = level + 1
			// $iterator->next(); if ($iterator->valid()) continue;
			if brokeAt >= 0 && brokeAt+1 < len(requestRules) {
				continue
			}
		}
		if level < systemLevel {
			systemLevel = level
		}
		rulesCount := s.rules.Count()
		for i, n := 0, 0; n < rulesCount; i, n = i+1, n+1 {
			if i == rulesCount {
				i = 0
			}
			rule := s.rules.RuleByID[i]
			if rule.disabled {
				continue
			}
			var decisionQueue []int
			skip := false
			for _, lit := range rule.literals {
				if lit <= 0 {
					if !s.decisions.decidedInstall(lit) {
						skip = true
						break
					}
				} else {
					if s.decisions.decidedInstall(lit) {
						skip = true
						break
					}
					if s.decisions.undecided(lit) {
						decisionQueue = append(decisionQueue, lit)
					}
				}
			}
			if skip || len(decisionQueue) < 2 {
				continue
			}
			level = s.selectAndInstall(level, decisionQueue, rule)
			if level == 0 {
				return
			}
			rulesCount = s.rules.Count()
			n = -1
		}
		if level < systemLevel {
			continue
		}
		if len(s.branches) > 0 {
			lastLiteral := 0
			lastLevel := 0
			lastBranchIndex := 0
			lastBranchOffset := 0
			found := false
			for i := len(s.branches) - 1; i >= 0; i-- {
				b := s.branches[i]
				for offset, lit := range b.literals {
					if lit > 0 && s.decisions.decisionLevel(lit) > b.level+1 {
						lastLiteral = lit
						lastBranchIndex = i
						lastBranchOffset = offset
						lastLevel = b.level
						found = true
					}
				}
			}
			if found {
				b := &s.branches[lastBranchIndex]
				b.literals = append(b.literals[:lastBranchOffset:lastBranchOffset], b.literals[lastBranchOffset+1:]...)
				level = lastLevel
				s.revert(level)
				why := s.decisions.queue[len(s.decisions.queue)-1].reason
				level = s.setPropagateLearn(level, lastLiteral, why)
				if level == 0 {
					return
				}
				continue
			}
		}
		break
	}
}
