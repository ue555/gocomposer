package solver

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// PlatformInfo answers the questions Problem asks about the running PHP.
type PlatformInfo interface {
	// PHPVersion is phpversion() of the PHP that Composer would run on.
	PHPVersion() string
	// ExtensionVersion is phpversion($ext); ok is false when not loaded.
	ExtensionLoaded(ext string) bool
	ExtensionVersion(ext string) (string, bool)
	// IniPaths is IniHelper::getAll().
	IniPaths() []string
}

// Problem is Composer\DependencyResolver\Problem.
type Problem struct {
	reasonSeen map[*Rule]bool
	reasons    [][]*Rule
	section    int
}

func (p *Problem) addRule(r *Rule) {
	if p.reasonSeen == nil {
		p.reasonSeen = map[*Rule]bool{}
	}
	if p.reasonSeen[r] {
		return
	}
	p.reasonSeen[r] = true
	for len(p.reasons) <= p.section {
		p.reasons = append(p.reasons, nil)
	}
	p.reasons[p.section] = append(p.reasons[p.section], r)
}

func (p *Problem) nextSection() { p.section++ }

// Reasons returns the problem rules grouped by section.
func (p *Problem) Reasons() [][]*Rule { return p.reasons }

// Context bundles what problem formatting needs.
type Context struct {
	RepoSet  *RepositorySet
	Request  *Request
	Pool     *Pool
	Verbose  bool
	Platform PlatformInfo
	// CompleteAdvisories is Problem::loadCompleteAdvisories (optional).
	CompleteAdvisories func([]policy.VersionReason) []policy.VersionReason
}

// PrettyString is SolverProblemsException::getPrettyString.
func (e *ProblemsError) PrettyString(ctx *Context, isDevExtraction bool) string {
	installedMap := ctx.Request.PresentMap()
	var missingExtensions []string
	isCausedByLock := false
	var problems []string
	for _, p := range e.Problems {
		problems = append(problems, p.PrettyString(ctx, installedMap, e.LearnedPool)+"\n")
		for _, section := range p.reasons {
			for _, r := range section {
				if req := r.RequiredPackage(); req != "" && strings.HasPrefix(req, "ext-") && !containsStr(missingExtensions, req) {
					missingExtensions = append(missingExtensions, req)
				}
			}
		}
		isCausedByLock = isCausedByLock || p.isCausedByLock(ctx)
	}
	text := "\n"
	i := 1
	seen := map[string]bool{}
	for _, pr := range problems {
		if seen[pr] {
			continue
		}
		seen[pr] = true
		text += fmt.Sprintf("  Problem %d%s", i, pr)
		i++
	}
	var hints []string
	if !isDevExtraction && (strings.Contains(text, "could not be found") || strings.Contains(text, "no matching package found")) {
		hints = append(hints, "Potential causes:\n - A typo in the package name\n - The package is not available in a stable-enough version according to your minimum-stability setting\n   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.\n - It's a private package and you forgot to add a custom repository to find it\n\nRead <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.")
	}
	if len(missingExtensions) > 0 {
		hints = append(hints, createExtensionHint(ctx, missingExtensions))
	}
	if isCausedByLock && !isDevExtraction && !ctx.Request.UpdateAllowTransitiveRootDependencies() {
		hints = append(hints, "Use the option --with-all-dependencies (-W) to allow upgrades, downgrades and removals for packages currently locked to specific versions.")
	}
	if strings.Contains(text, "found composer-plugin-api[2.0.0] but it does not match") && strings.Contains(text, "- ocramius/package-versions") {
		hints = append(hints, "<warning>ocramius/package-versions only provides support for Composer 2 in 1.8+, which requires PHP 7.4.</warning>\nIf you can not upgrade PHP you can require <info>composer/package-versions-deprecated</info> to resolve this with PHP 7.0+.")
	}
	if strings.Contains(text, "found composer-plugin-api[2.0.0] but it does not match") {
		hints = append(hints, "You are using Composer 2, which some of your plugins seem to be incompatible with. Make sure you update your plugins or report a plugin-issue to ask them to support Composer 2.")
	}
	if len(hints) > 0 {
		text += "\n" + strings.Join(hints, "\n\n")
	}
	return text
}

func createExtensionHint(ctx *Context, missing []string) string {
	var paths []string
	if ctx.Platform != nil {
		paths = ctx.Platform.IniPaths()
	}
	if len(paths) == 0 {
		return ""
	}
	if paths[0] == "" {
		if len(paths) == 1 {
			return ""
		}
		paths = paths[1:]
	}
	var args []string
	for _, ext := range missing {
		args = append(args, "--ignore-platform-req="+ext)
	}
	text := "To enable extensions, verify that they are enabled in your .ini files:\n    - "
	text += strings.Join(paths, "\n    - ")
	text += "\nYou can also run `php --ini` in a terminal to see which files are used by PHP in CLI mode."
	text += "\nAlternatively, you can run Composer with `" + strings.Join(args, " ") + "` to temporarily ignore these required extensions."
	return text
}

func (p *Problem) isCausedByLock(ctx *Context) bool {
	for _, section := range p.reasons {
		for _, r := range section {
			if r.isCausedByLock(ctx) {
				return true
			}
		}
	}
	return false
}

func (r *Rule) isCausedByLock(ctx *Context) bool {
	req := ctx.Request
	if r.Reason == RulePackageRequires {
		link := r.ReasonData.(*pkg.Link)
		if repo.IsPlatformPackage(link.Target) {
			return false
		}
		if req.LockedRepository != nil {
			for _, p := range req.LockedRepository.MustPackages() {
				if p.Name == link.Target {
					if ctx.Pool.IsUnacceptableFixedOrLockedPackage(p) {
						return true
					}
					if !link.Constraint.Matches(semver.NewConstraint("=", p.Version())) {
						return true
					}
					if !req.IsLockedPackage(p) {
						return true
					}
					break
				}
			}
		}
	}
	if r.Reason == RuleRootRequire {
		d := r.ReasonData.(RootRequireData)
		if repo.IsPlatformPackage(d.PackageName) {
			return false
		}
		if req.LockedRepository != nil {
			for _, p := range req.LockedRepository.MustPackages() {
				if p.Name == d.PackageName {
					if ctx.Pool.IsUnacceptableFixedOrLockedPackage(p) {
						return true
					}
					if !d.Constraint.Matches(semver.NewConstraint("=", p.Version())) {
						return true
					}
					break
				}
			}
		}
	}
	return false
}

// PrettyString is Problem::getPrettyString.
func (p *Problem) PrettyString(ctx *Context, installedMap map[int]*pkg.Package, learnedPool [][]*Rule) string {
	var reasons []*Rule
	for i := len(p.reasons) - 1; i >= 0; i-- {
		reasons = append(reasons, p.reasons[i]...)
	}
	if len(reasons) == 1 {
		rule := reasons[0]
		if rule.Reason == RuleRootRequire {
			d := rule.ReasonData.(RootRequireData)
			if len(ctx.Pool.WhatProvides(d.PackageName, d.Constraint)) == 0 {
				parts := getMissingPackageReason(ctx, d.PackageName, d.Constraint)
				return "\n    " + parts[0] + parts[1]
			}
		} else if rule.Reason == RuleLockedFilterListRemove {
			parts := getMissingLockedPackageReason(ctx.Pool, rule.ReasonData.(*pkg.Package))
			return "\n    " + parts[0] + parts[1]
		}
	}
	sort.SliceStable(reasons, func(i, j int) bool {
		p1, p2 := rulePriority(reasons[i]), rulePriority(reasons[j])
		if p1 != p2 {
			return p2 < p1
		}
		return phpjson.PHPCompareKeys(sortableString(ctx.Pool, reasons[i]), sortableString(ctx.Pool, reasons[j])) < 0
	})
	return formatDeduplicatedRules(reasons, "    ", ctx, installedMap, learnedPool) + formatPolicyExcludedConflictAlternatives(reasons, ctx)
}

func rulePriority(r *Rule) int {
	switch r.Reason {
	case RuleFixed, RuleLockedFilterListRemove:
		return 3
	case RuleRootRequire:
		return 2
	case RulePackageConflict, RulePackageRequires:
		return 1
	}
	return 0
}

func sortableString(pool *Pool, r *Rule) string {
	switch r.Reason {
	case RuleRootRequire:
		return r.ReasonData.(RootRequireData).PackageName
	case RuleFixed, RuleLockedFilterListRemove:
		return r.ReasonData.(*pkg.Package).String()
	case RulePackageConflict, RulePackageRequires:
		src := r.sourcePackage(pool)
		return src.String() + "//" + r.ReasonData.(*pkg.Link).PrettyString(src)
	case RulePackageSameName:
		return r.ReasonData.(string)
	case RulePackageAlias, RulePackageInverseAlias:
		return r.ReasonData.(*pkg.Package).String()
	case RuleLearned:
		parts := make([]string, len(r.literals))
		for i, l := range r.literals {
			parts[i] = fmt.Sprint(l)
		}
		return strings.Join(parts, "-")
	}
	return ""
}

func deduplicateDefaultBranchAlias(p *pkg.Package) *pkg.Package {
	if p.AliasOf != nil && p.PrettyVersion() == semver.DefaultBranchAlias {
		return p.AliasOf
	}
	return p
}

func (r *Rule) sourcePackage(pool *Pool) *pkg.Package {
	switch r.Reason {
	case RulePackageConflict:
		p1 := deduplicateDefaultBranchAlias(pool.LiteralToPackage(r.literals[0]))
		p2 := deduplicateDefaultBranchAlias(pool.LiteralToPackage(r.literals[1]))
		if r.ReasonData.(*pkg.Link).Source == p1.Name {
			p1, p2 = p2, p1
		}
		_ = p1
		return p2
	case RulePackageRequires:
		return deduplicateDefaultBranchAlias(pool.LiteralToPackage(r.literals[0]))
	}
	panic("Not implemented")
}

var (
	reDedupMessage = regexp.MustCompile(`^(\S+) (\S+) (requires|conflicts)`)
	reDedupPrefix  = regexp.MustCompile(`^\S+ \S+ `)
	reGrammar      = regexp.MustCompile(`^(%s%s (?:require|conflict))s`)
)

func formatDeduplicatedRules(rules []*Rule, indent string, ctx *Context, installedMap map[int]*pkg.Package, learnedPool [][]*Rule) string {
	var messages []string
	type tmplData struct {
		pkgOrder []string
		versions map[string]*orderedVersions
	}
	templates := map[string]*tmplData{}
	for _, r := range rules {
		message := r.PrettyString(ctx, installedMap, learnedPool)
		if (r.Reason == RulePackageRequires || r.Reason == RulePackageConflict) && reDedupMessage.MatchString(message) {
			m := reDedupMessage.FindStringSubmatch(message)
			message = strings.ReplaceAll(message, "%", "%%")
			template := reDedupPrefix.ReplaceAllLiteralString(message, "%s%s ")
			messages = append(messages, template)
			td, ok := templates[template]
			if !ok {
				td = &tmplData{versions: map[string]*orderedVersions{}}
				templates[template] = td
			}
			if _, ok := td.versions[m[1]]; !ok {
				td.pkgOrder = append(td.pkgOrder, m[1])
				td.versions[m[1]] = newOrderedVersions()
			}
			norm, err := semver.Normalize(m[2])
			if err != nil {
				norm = m[2]
			}
			td.versions[m[1]].set(norm, m[2])
			src := r.sourcePackage(ctx.Pool)
			rv := ctx.Pool.RemovedVersionsByPackage(src)
			for _, v := range rv.keys {
				td.versions[m[1]].set(v, rv.vals[v])
			}
		} else if message != "" {
			messages = append(messages, message)
		}
	}
	var result []string
	seen := map[string]bool{}
	for _, message := range messages {
		if seen[message] {
			continue
		}
		seen[message] = true
		if td, ok := templates[message]; ok {
			for _, pkgName := range td.pkgOrder {
				versions := td.versions[pkgName]
				sorted := sortVersionKeys(versions)
				list := make([]string, len(sorted))
				for i, k := range sorted {
					list[i] = versions.vals[k]
				}
				if !ctx.Verbose {
					list = condenseVersionList(sorted, versions, 1, 16)
				}
				if len(list) > 1 {
					msg := reGrammar.ReplaceAllString(message, "$1")
					result = append(result, phpSprintf(msg, pkgName, "["+strings.Join(list, ", ")+"]"))
				} else {
					result = append(result, phpSprintf(message, pkgName, " "+list[0]))
				}
			}
		} else {
			result = append(result, message)
		}
	}
	return "\n" + indent + "- " + strings.Join(result, "\n"+indent+"- ")
}

// phpSprintf handles the %s/%% subset used by templates.
func phpSprintf(format string, args ...string) string {
	var sb strings.Builder
	ai := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) {
			switch format[i+1] {
			case '%':
				sb.WriteByte('%')
				i++
				continue
			case 's':
				if ai < len(args) {
					sb.WriteString(args[ai])
				}
				ai++
				i++
				continue
			}
		}
		sb.WriteByte(format[i])
	}
	return sb.String()
}

// sortVersionKeys is uksort($versions, 'version_compare').
func sortVersionKeys(v *orderedVersions) []string {
	keys := append([]string(nil), v.keys...)
	sort.SliceStable(keys, func(i, j int) bool { return semver.VersionCompare(keys[i], keys[j]) < 0 })
	return keys
}

var reMajor = regexp.MustCompile(`^(\d+)\..*`)

// condenseVersionList ports Problem::condenseVersionList on sorted keys.
func condenseVersionList(keys []string, versions *orderedVersions, max, maxDev int) []string {
	if len(keys) <= max {
		out := make([]string, len(keys))
		for i, k := range keys {
			out[i] = versions.vals[k]
		}
		return out
	}
	var majorOrder []string
	byMajor := map[string][]string{}
	for _, k := range keys {
		var major string
		if len(k) >= 4 && strings.EqualFold(k[:4], "dev-") {
			major = "dev"
		} else {
			major = reMajor.ReplaceAllString(k, "$1")
		}
		if _, ok := byMajor[major]; !ok {
			majorOrder = append(majorOrder, major)
		}
		byMajor[major] = append(byMajor[major], versions.vals[k])
	}
	var filtered []string
	for _, major := range majorOrder {
		vs := byMajor[major]
		m := max
		if major == "dev" {
			m = maxDev
		}
		if len(vs) > m {
			filtered = append(filtered, vs[0], "...", vs[len(vs)-1])
		} else {
			filtered = append(filtered, vs...)
		}
	}
	return filtered
}

// GetPackageList is Problem::getPackageList.
func GetPackageList(packages []*pkg.Package, verbose bool, pool *Pool, c semver.Constraint, useRemovedVersionGroup bool) string {
	var nameOrder []string
	prettyNames := map[string]string{}
	versions := map[string]*orderedVersions{}
	hasDefaultBranch := map[string]bool{}
	for _, p := range packages {
		if _, ok := versions[p.Name]; !ok {
			nameOrder = append(nameOrder, p.Name)
			versions[p.Name] = newOrderedVersions()
		}
		prettyNames[p.Name] = p.PrettyName
		pv := p.PrettyVersion()
		if p.AliasOf != nil {
			pv += " (alias of " + p.AliasOf.PrettyVersion() + ")"
		}
		versions[p.Name].set(p.Version(), pv)
		if pool != nil && c != nil {
			rv := pool.RemovedVersions(p.Name, c)
			for _, k := range rv.keys {
				versions[p.Name].set(k, rv.vals[k])
			}
		}
		if pool != nil && useRemovedVersionGroup {
			rv := pool.RemovedVersionsByPackage(p)
			for _, k := range rv.keys {
				versions[p.Name].set(k, rv.vals[k])
			}
		}
		if p.IsDefaultBranch() {
			hasDefaultBranch[p.Name] = true
		}
	}
	var out []string
	for _, name := range nameOrder {
		v := versions[name]
		if _, ok := v.vals[semver.DefaultBranchAlias]; ok && hasDefaultBranch[name] {
			nv := newOrderedVersions()
			for _, k := range v.keys {
				if k != semver.DefaultBranchAlias {
					nv.set(k, v.vals[k])
				}
			}
			v = nv
		}
		keys := sortVersionKeys(v)
		var list []string
		if !verbose {
			list = condenseVersionList(keys, v, 4, 16)
		} else {
			for _, k := range keys {
				list = append(list, v.vals[k])
			}
		}
		out = append(out, prettyNames[name]+"["+strings.Join(list, ", ")+"]")
	}
	return strings.Join(out, ", ")
}

func (r *Rule) formatPackagesUnique(pool *Pool, items []*pkg.Package, verbose bool, c semver.Constraint, useRemovedVersionGroup bool) string {
	return GetPackageList(items, verbose, pool, c, useRemovedVersionGroup)
}

func (r *Rule) literalPackages(pool *Pool, lits []int) []*pkg.Package {
	out := make([]*pkg.Package, len(lits))
	for i, l := range lits {
		out[i] = pool.LiteralToPackage(l)
	}
	return out
}

// PrettyString is Rule::getPrettyString.
func (r *Rule) PrettyString(ctx *Context, installedMap map[int]*pkg.Package, learnedPool [][]*Rule) string {
	pool := ctx.Pool
	lits := r.literals
	switch r.Reason {
	case RuleRootRequire:
		d := r.ReasonData.(RootRequireData)
		packages := pool.WhatProvides(d.PackageName, d.Constraint)
		if len(packages) == 0 {
			return "No package found to satisfy root composer.json require " + d.PackageName + " " + d.Constraint.PrettyString()
		}
		var nonAlias []*pkg.Package
		for _, p := range packages {
			if p.AliasOf == nil {
				nonAlias = append(nonAlias, p)
			}
		}
		if len(nonAlias) == 1 && ctx.Request.IsLockedPackage(nonAlias[0]) {
			return nonAlias[0].PrettyName + " is locked to version " + nonAlias[0].PrettyVersion() + " and an update of this package was not requested."
		}
		return "Root composer.json requires " + d.PackageName + " " + d.Constraint.PrettyString() + " -> satisfiable by " + r.formatPackagesUnique(pool, packages, ctx.Verbose, d.Constraint, false) + "."
	case RuleFixed:
		p := deduplicateDefaultBranchAlias(r.ReasonData.(*pkg.Package))
		if ctx.Request.IsLockedPackage(p) {
			return p.PrettyName + " is locked to version " + p.PrettyVersion() + " and an update of this package was not requested."
		}
		return p.PrettyName + " is present at version " + p.PrettyVersion() + " and cannot be modified by Composer"
	case RulePackageConflict:
		p1 := deduplicateDefaultBranchAlias(pool.LiteralToPackage(lits[0]))
		p2 := deduplicateDefaultBranchAlias(pool.LiteralToPackage(lits[1]))
		conflictTarget := p1.PrettyString()
		link := r.ReasonData.(*pkg.Link)
		if link.Source == p1.Name {
			p1, p2 = p2, p1
			conflictTarget = p1.PrettyName + " " + link.PrettyConstraint
		}
		if link.Target != p1.Name {
			provideType := ""
			provided := ""
			for _, l := range p1.Provides() {
				if l.Target == link.Target {
					provideType = "provides"
					provided = l.PrettyConstraint
					break
				}
			}
			for _, l := range p1.Replaces() {
				if l.Target == link.Target {
					provideType = "replaces"
					provided = l.PrettyConstraint
					break
				}
			}
			if provideType != "" {
				conflictTarget = link.Target + " " + link.PrettyConstraint + " (" + p1.PrettyString() + " " + provideType + " " + link.Target + " " + provided + ")"
			}
		}
		return p2.PrettyString() + " conflicts with " + conflictTarget + "."
	case RulePackageRequires:
		source := deduplicateDefaultBranchAlias(pool.LiteralToPackage(lits[0]))
		link := r.ReasonData.(*pkg.Link)
		requires := r.literalPackages(pool, lits[1:])
		text := link.PrettyString(source)
		if len(requires) > 0 {
			text += " -> satisfiable by " + r.formatPackagesUnique(pool, requires, ctx.Verbose, link.Constraint, false) + "."
		} else {
			reason := getMissingPackageReason(ctx, link.Target, link.Constraint)
			return text + " -> " + reason[1]
		}
		return text
	case RulePackageSameName:
		names := map[string]bool{}
		var nameOrder []string
		for _, l := range lits {
			n := pool.LiteralToPackage(l).Name
			if !names[n] {
				names[n] = true
				nameOrder = append(nameOrder, n)
			}
		}
		replacedName := r.ReasonData.(string)
		if len(nameOrder) > 1 {
			var reason string
			if !names[replacedName] {
				word := "all"
				if len(lits) == 2 {
					word = "both"
				}
				reason = "They " + word + " replace " + replacedName + " and thus cannot coexist."
			} else {
				var replacers []string
				for _, n := range nameOrder {
					if n != replacedName {
						replacers = append(replacers, n)
					}
				}
				if len(replacers) == 1 {
					reason = replacers[0] + " replaces "
				} else {
					reason = "[" + strings.Join(replacers, ", ") + "] replace "
				}
				reason += replacedName + " and thus cannot coexist with it."
			}
			var installed, removable []*pkg.Package
			for _, l := range lits {
				if _, ok := installedMap[abs(l)]; ok {
					installed = append(installed, pool.LiteralToPackage(l))
				} else {
					removable = append(removable, pool.LiteralToPackage(l))
				}
			}
			if len(installed) > 0 && len(removable) > 0 {
				return r.formatPackagesUnique(pool, removable, ctx.Verbose, nil, true) + " cannot be installed as that would require removing " + r.formatPackagesUnique(pool, installed, ctx.Verbose, nil, true) + ". " + reason
			}
			return "Only one of these can be installed: " + r.formatPackagesUnique(pool, r.literalPackages(pool, lits), ctx.Verbose, nil, true) + ". " + reason
		}
		return "You can only install one version of a package, so only one of these can be installed: " + r.formatPackagesUnique(pool, r.literalPackages(pool, lits), ctx.Verbose, nil, true) + "."
	case RuleLearned:
		learned := " (conflict analysis result)"
		var ruleText string
		if len(lits) == 1 {
			ruleText = pool.LiteralToPrettyString(lits[0], installedMap)
		} else {
			var groupOrder []string
			groups := map[string][]*pkg.Package{}
			for _, l := range lits {
				p := pool.LiteralToPackage(l)
				var g string
				if _, ok := installedMap[p.ID]; ok {
					if l > 0 {
						g = "keep"
					} else {
						g = "remove"
					}
				} else if l > 0 {
					g = "install"
				} else {
					g = "don't install"
				}
				if _, ok := groups[g]; !ok {
					groupOrder = append(groupOrder, g)
				}
				groups[g] = append(groups[g], deduplicateDefaultBranchAlias(p))
			}
			var texts []string
			for _, g := range groupOrder {
				ps := groups[g]
				one := ""
				if len(ps) > 1 {
					one = " one of"
				}
				texts = append(texts, g+one+" "+r.formatPackagesUnique(pool, ps, ctx.Verbose, nil, false))
			}
			ruleText = strings.Join(texts, " | ")
		}
		return "Conclusion: " + ruleText + learned
	case RulePackageAlias:
		alias := pool.LiteralToPackage(lits[0])
		if alias.Version() == semver.DefaultBranchAlias {
			return ""
		}
		p := deduplicateDefaultBranchAlias(pool.LiteralToPackage(lits[1]))
		return alias.PrettyString() + " is an alias of " + p.PrettyString() + " and thus requires it to be installed too."
	case RulePackageInverseAlias:
		alias := pool.LiteralToPackage(lits[1])
		if alias.Version() == semver.DefaultBranchAlias {
			return ""
		}
		p := deduplicateDefaultBranchAlias(pool.LiteralToPackage(lits[0]))
		return alias.PrettyString() + " is an alias of " + p.PrettyString() + " and must be installed with it."
	}
	var parts []string
	for _, l := range lits {
		parts = append(parts, pool.LiteralToPrettyString(l, installedMap))
	}
	return "(" + strings.Join(parts, "|") + ")"
}

var (
	reDevHash       = regexp.MustCompile(`^dev-.*#.*`)
	reAliasSuffix   = regexp.MustCompile(` +as +([^,\s|]+)$`)
	reValidPkgName  = regexp.MustCompile(`^[A-Za-z0-9_./-]+$`)
	reValidPkgChars = regexp.MustCompile(`[A-Za-z0-9_./-]+`)
	reNumericExact  = regexp.MustCompile(`^\d+(?:\.\d+)*$`)
)

func constraintToText(c semver.Constraint) string {
	if s, ok := c.(*semver.Single); ok && s.Operator() == "==" && !strings.HasPrefix(s.Version(), "dev-") {
		if !reNumericExact.MatchString(s.PrettyString()) {
			return " " + s.PrettyString() + " (exact version match)"
		}
		versions := []string{s.PrettyString()}
		for i := 3 - strings.Count(versions[0], "."); i > 0; i-- {
			versions = append(versions, versions[len(versions)-1]+".0")
		}
		if len(versions) > 1 {
			return " " + s.PrettyString() + " (exact version match: " + strings.Join(versions[:len(versions)-1], ", ") + " or " + versions[len(versions)-1] + ")"
		}
		return " " + s.PrettyString() + " (exact version match: " + versions[0] + ")"
	}
	if c != nil {
		return " " + c.PrettyString()
	}
	return ""
}

func getPlatformPackageVersion(pool *Pool, name, version string) (string, bool) {
	available := pool.WhatProvides(name, nil)
	if len(available) == 0 {
		return "", false
	}
	var selected *pkg.Package
	for _, p := range available {
		if isPlatformRepo(p) {
			selected = p
			break
		}
	}
	if selected == nil {
		selected = available[0]
	}
	if selected.Name != name {
		for _, l := range append(append(pkg.Links(nil), selected.Provides()...), selected.Replaces()...) {
			if l.Target == name {
				return l.PrettyConstraint + " " + l.Description[:len(l.Description)-1] + "d by " + selected.PrettyString(), true
			}
		}
	}
	version = selected.PrettyVersion()
	if extra := selected.Extra(); extra != nil {
		if v, ok := extra.Get("config.platform"); ok && v == true {
			version += "; " + strings.ReplaceAll(selected.Description(), "Package ", "")
		}
	}
	return version, true
}

func getProvidersList(rs *RepositorySet, name string, max int) (string, bool) {
	providers := rs.Providers(name)
	if len(providers) == 0 {
		return "", false
	}
	list := providers
	if len(providers) > max+1 {
		list = providers[:max]
	}
	var sb strings.Builder
	for _, p := range list {
		desc := ""
		if p.Description != "" {
			d := p.Description
			if len(d) > 100 {
				d = d[:100]
			}
			desc = " " + d
		}
		sb.WriteString("      - " + p.Name + desc + "\n")
	}
	if len(providers) > max+1 {
		sb.WriteString(fmt.Sprintf("      ... and %d more.\n", len(providers)-max))
	}
	return sb.String(), true
}

func hasMultipleNames(packages []*pkg.Package) bool {
	name := ""
	for _, p := range packages {
		if name == "" || name == p.Name {
			name = p.Name
		} else {
			return true
		}
	}
	return false
}

// getMissingPackageReason is Problem::getMissingPackageReason.
func getMissingPackageReason(ctx *Context, name string, c semver.Constraint) [2]string {
	pool := ctx.Pool
	rs := ctx.RepoSet
	verbose := ctx.Verbose
	lname := strings.ToLower(name)
	if repo.IsPlatformPackage(name) {
		if strings.HasPrefix(lname, "php") || name == "hhvm" {
			phpv := ""
			if ctx.Platform != nil {
				phpv = ctx.Platform.PHPVersion()
			}
			version, ok := getPlatformPackageVersion(pool, name, phpv)
			msg := "- Root composer.json requires " + name + constraintToText(c) + " but "
			if name == "hhvm" && len(pool.WhatProvides(name, nil)) > 0 {
				return [2]string{msg, "your HHVM version does not satisfy that requirement."}
			}
			if name == "hhvm" {
				return [2]string{msg, "HHVM was not detected on this machine, make sure it is in your PATH."}
			}
			if !ok {
				return [2]string{msg, `the ` + name + ` package is disabled by your platform config. Enable it again with "composer config platform.` + name + ` --unset".`}
			}
			return [2]string{msg, "your " + name + " version (" + version + ") does not satisfy that requirement."}
		}
		if strings.HasPrefix(lname, "ext-") {
			if strings.Contains(name, " ") {
				return [2]string{"- ", "PHP extension " + name + " should be required as " + strings.ReplaceAll(name, " ", "-") + "."}
			}
			ext := name[4:]
			msg := "- Root composer.json requires PHP extension " + name + constraintToText(c) + " but "
			version, ok := getPlatformPackageVersion(pool, name, "0")
			if !ok {
				providersStr, has := getProvidersList(rs, name, 5)
				if has {
					providersStr = "\n\n      Alternatively you can require one of these packages that provide the extension (or parts of it):\n" +
						"      <warning>Keep in mind that the suggestions are automated and may not be valid or safe to use</warning>\n" + providersStr
				}
				if ctx.Platform != nil && ctx.Platform.ExtensionLoaded(ext) {
					return [2]string{msg, `the ` + name + ` package is disabled by your platform config. Enable it again with "composer config platform.` + name + ` --unset".` + providersStr}
				}
				return [2]string{msg, "it is missing from your system. Install or enable PHP's " + ext + " extension." + providersStr}
			}
			return [2]string{msg, "it has the wrong version installed (" + version + ")."}
		}
		if strings.HasPrefix(lname, "lib-") {
			if lname == "lib-icu" {
				errMsg := "it is missing from your system, make sure the intl extension is loaded."
				if ctx.Platform != nil && ctx.Platform.ExtensionLoaded("intl") {
					errMsg = "it has the wrong version installed, try upgrading the intl extension."
				}
				return [2]string{"- Root composer.json requires linked library " + name + constraintToText(c) + " but ", errMsg}
			}
			providersStr, has := getProvidersList(rs, name, 5)
			if has {
				providersStr = "\n\n      Alternatively you can require one of these packages that provide the library (or parts of it):\n" +
					"      <warning>Keep in mind that the suggestions are automated and may not be valid or safe to use</warning>\n" + providersStr
			}
			return [2]string{"- Root composer.json requires linked library " + name + constraintToText(c) + " but ", "it has the wrong version installed or is missing from your system, make sure to load the extension providing it." + providersStr}
		}
	}
	var lockedPackage *pkg.Package
	for _, p := range ctx.Request.LockedPackages() {
		if p.Name == name {
			lockedPackage = p
			if pool.IsUnacceptableFixedOrLockedPackage(p) {
				return [2]string{"- ", p.PrettyName + " is fixed to " + p.PrettyVersion() + " (lock file version) by a partial update but that version is rejected by your minimum-stability. Make sure you list it as an argument for the update command."}
			}
			break
		}
	}
	if s, ok := c.(*semver.Single); ok && s.Operator() == "==" && reDevHash.MatchString(s.PrettyString()) {
		newC := reAliasSuffix.ReplaceAllString(s.PrettyString(), "")
		packages := rs.FindPackages(name, semver.NewMulti([]semver.Constraint{semver.NewConstraint("==", newC), semver.NewConstraint("==", strings.ReplaceAll(newC, "#", "+"))}, false), 0)
		if len(packages) > 0 {
			return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + `. The # character in branch names is replaced by a + character. Make sure to require it as "` + strings.ReplaceAll(s.PrettyString(), "#", "+") + `".`}
		}
	}
	packages := rs.FindPackages(name, c, 0)
	if len(packages) > 0 {
		if rr, ok := rs.RootRequire(name); ok {
			matched := false
			for _, p := range packages {
				if rr.Matches(semver.NewConstraint("==", p.Version())) {
					matched = true
					break
				}
			}
			if !matched {
				conf := "it conflicts"
				if hasMultipleNames(packages) {
					conf = "these conflict"
				}
				return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but " + conf + " with your root composer.json require (" + rr.PrettyString() + ")."}
			}
		}
		for _, n := range packages[0].Names(true) {
			if tc, ok := rs.temporaryConstraints[n]; ok {
				matched := false
				for _, p := range packages {
					if tc.Matches(semver.NewConstraint("==", p.Version())) {
						matched = true
						break
					}
				}
				if !matched {
					conf := "it conflicts"
					if hasMultipleNames(packages) {
						conf = "these conflict"
					}
					return [2]string{"- Root composer.json requires " + n + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but " + conf + " with your temporary update constraint (" + n + ":" + tc.PrettyString() + ")."}
				}
			}
		}
		if lockedPackage != nil {
			fixed := semver.NewConstraint("==", lockedPackage.Version())
			matched := false
			for _, p := range packages {
				if fixed.Matches(semver.NewConstraint("==", p.Version())) {
					matched = true
					break
				}
			}
			if !matched {
				return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but the package is fixed to " + lockedPackage.PrettyVersion() + " (lock file version) by a partial update and that version does not match. Make sure you list it as an argument for the update command."}
			}
		}
		if removals := ctx.Pool.PolicyRemovedVersions(name, c); len(removals) > 0 {
			return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but these were not loaded, " + formatPolicyRemovalReasons(ctx, removals, true)}
		}
		allLock := true
		for _, p := range packages {
			if repoKind(p) != repo.KindLock {
				allLock = false
				break
			}
		}
		if allLock {
			return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " in the lock file but not in remote repositories, make sure you avoid updating this package to keep the one from the lock file."}
		}
		conf := "it conflicts"
		if hasMultipleNames(packages) {
			conf = "they conflict"
		}
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but these were not loaded, likely because " + conf + " with another require."}
	}
	packages = rs.FindPackages(name, c, AllowUnacceptableStabilities)
	if len(packages) > 0 {
		all := rs.FindPackages(name, c, AllowShadowedRepositories)
		if len(all) > 0 {
			return computeCheckForLowerPrioRepo(pool, verbose, name, packages, all, "minimum-stability", c)
		}
		does := "it does"
		if hasMultipleNames(packages) {
			does = "these do"
		}
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but " + does + " not match your minimum-stability."}
	}
	packages = rs.FindPackages(name, nil, AllowUnacceptableStabilities)
	if len(packages) > 0 {
		all := rs.FindPackages(name, c, AllowShadowedRepositories)
		if len(all) > 0 {
			return computeCheckForLowerPrioRepo(pool, verbose, name, packages, all, "constraint", c)
		}
		suffix := ""
		if s, ok := c.(*semver.Single); ok && s.Version() == "dev-master" {
			for _, cand := range packages {
				if cand.Version() == "dev-default" || cand.Version() == "dev-main" {
					suffix = " Perhaps dev-master was renamed to " + cand.PrettyVersion() + "?"
					break
				}
			}
		}
		if packages[0].IsRoot {
			suffix = " See https://getcomposer.org/dep-on-root for details and assistance."
		}
		does := "it does"
		if hasMultipleNames(packages) {
			does = "these do"
		}
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ", "found " + GetPackageList(packages, verbose, pool, c, false) + " but " + does + " not match the constraint." + suffix}
	}
	if !reValidPkgName.MatchString(name) {
		illegal := reValidPkgChars.ReplaceAllString(name, "")
		return [2]string{"- Root composer.json requires " + name + ", it ", `could not be found, it looks like its name is invalid, "` + illegal + `" is not allowed in package names.`}
	}
	if providersStr, ok := getProvidersList(rs, name, 15); ok {
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", it ", "could not be found in any version, but the following packages provide it:\n" + providersStr + "      Consider requiring one of these to satisfy the " + name + " requirement."}
	}
	return [2]string{"- Root composer.json requires " + name + ", it ", "could not be found in any version, there may be a typo in the package name."}
}

func computeCheckForLowerPrioRepo(pool *Pool, verbose bool, name string, higher, all []*pkg.Package, reason string, c semver.Constraint) [2]string {
	var nextRepoPackages []*pkg.Package
	var nextRepo any
	for _, p := range all {
		if nextRepo == nil || nextRepo == p.Repo {
			nextRepoPackages = append(nextRepoPackages, p)
			nextRepo = p.Repo
		} else {
			break
		}
	}
	nr, _ := nextRepo.(repo.Repository)
	if len(higher) > 0 && higher[0].IsRoot {
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", it is ",
			"satisfiable by " + GetPackageList(nextRepoPackages, verbose, pool, c, false) + " from " + nr.RepoName() + " but " + higher[0].PrettyName + " " + higher[0].PrettyVersion() + " is the root package and cannot be modified. See https://getcomposer.org/dep-on-root for details and assistance."}
	}
	if nr != nil && nr.Kind() == repo.KindLock {
		singular := len(higher) == 1
		suggestion := "Make sure you either fix the " + reason + " or avoid updating this package to keep the one present in the lock file (" + GetPackageList(nextRepoPackages, verbose, pool, c, false) + ")."
		if nextRepoPackages[0].DistType() == "path" {
			to := nextRepoPackages[0].TransportOptions()
			sym, has := to.Get("symlink")
			if !has || sym != false {
				suggestion = "Make sure you fix the " + reason + " as packages installed from symlinked path repos are updated even in partial updates and the one from the lock file can thus not be used."
			}
		}
		itDoes, isAre := "these do", "are"
		if singular {
			itDoes, isAre = "it does", "is"
		}
		return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", ",
			"found " + GetPackageList(higher, verbose, pool, c, false) + " but " + itDoes + " not match your " + reason + " and " + isAre + " therefore not installable. " + suggestion}
	}
	higherRepo := ""
	if hr, ok := higher[0].Repo.(repo.Repository); ok {
		higherRepo = hr.RepoName()
	}
	nextName := ""
	if nr != nil {
		nextName = nr.RepoName()
	}
	return [2]string{"- Root composer.json requires " + name + constraintToText(c) + ", it is ", "satisfiable by " + GetPackageList(nextRepoPackages, verbose, pool, c, false) + " from " + nextName + " but " + GetPackageList(higher, verbose, pool, c, false) + " from " + higherRepo + " has higher repository priority. The packages from the higher priority repository do not match your " + reason + " and are therefore not installable. That repository is canonical so the lower priority repo's packages are not installable. See https://getcomposer.org/repoprio for details and assistance."}
}

// PolicyRemovedVersions is Pool::getPolicyRemovedVersions.
func (p *Pool) PolicyRemovedVersions(name string, c semver.Constraint) []policy.VersionReason {
	return p.policyRemoved.Versions(name, c)
}

// getMissingLockedPackageReason is Problem::getMissingLockedPackageReason.
func getMissingLockedPackageReason(pool *Pool, p *pkg.Package) [2]string {
	prefix := "- Package " + p.Name + " " + p.PrettyVersion() + " (in the lock file) "
	if r := pool.PolicyRemovalReason(p.Name, p.Version()); r != nil {
		return [2]string{prefix, "was not loaded, because it " + formatPolicyRemovalSentence(r, false)}
	}
	panic(&SolverBugError{"Filter list removed locked package must have version removed from pool."})
}

// formatPolicyRemovalSentence is Problem::formatPolicyRemovalSentence.
func formatPolicyRemovalSentence(r *policy.PolicyRemovalReason, plural bool) string {
	remedy := r.Remedy()
	s := r.Verb(plural) + " " + r.Description() + "."
	if remedy != "" {
		s += " " + remedy
	}
	return s + " To turn the feature off entirely, you can set " + r.OffSwitch() + " to false."
}

// formatPolicyRemovalReasons is Problem::formatPolicyRemovalReasons.
func formatPolicyRemovalReasons(ctx *Context, removals []policy.VersionReason, plural bool) string {
	if ctx.CompleteAdvisories != nil {
		removals = ctx.CompleteAdvisories(removals)
	}
	types := map[string]bool{}
	for _, r := range removals {
		types[r.Reason.Type()] = true
	}
	if len(types) == 1 {
		reasons := make([]*policy.PolicyRemovalReason, len(removals))
		for i, r := range removals {
			reasons[i] = r.Reason
		}
		combined, err := policy.CombineRemovalReasons(reasons)
		if err != nil {
			panic(&SolverBugError{err.Error()})
		}
		they := "it "
		if plural {
			they = "they "
		}
		return "because " + they + formatPolicyRemovalSentence(combined, plural)
	}
	var descOrder []string
	versionsByReason := map[string][]string{}
	var remedies, offSwitches []string
	addUnique := func(list []string, v string) []string {
		for _, e := range list {
			if e == v {
				return list
			}
		}
		return append(list, v)
	}
	for _, r := range removals {
		d := r.Reason.Description()
		if _, ok := versionsByReason[d]; !ok {
			descOrder = append(descOrder, d)
		}
		versionsByReason[d] = addUnique(versionsByReason[d], r.Reason.PrettyVersion())
		if rem := r.Reason.Remedy(); rem != "" {
			remedies = addUnique(remedies, rem)
		}
		offSwitches = addUnique(offSwitches, r.Reason.OffSwitch())
	}
	var lines []string
	for _, d := range descOrder {
		lines = append(lines, strings.Join(versionsByReason[d], ", ")+": "+d)
	}
	remedies = addUnique(remedies, "To turn a policy off entirely, you can set "+strings.Join(offSwitches, " or ")+" to false.")
	return "because:\n      - " + strings.Join(lines, "\n      - ") + "\n      " + strings.Join(remedies, " ")
}

// formatPolicyExcludedConflictAlternatives is
// Problem::formatPolicyExcludedConflictAlternatives.
func formatPolicyExcludedConflictAlternatives(reasons []*Rule, ctx *Context) string {
	if ctx.Pool.policyRemoved == nil {
		return ""
	}
	rootConstraints := map[string]semver.Constraint{}
	requireConstraints := map[string][]semver.Constraint{}
	conflictConstraints := map[string][]semver.Constraint{}
	var conflictOrder []string
	for _, r := range reasons {
		switch r.Reason {
		case RuleRootRequire:
			d := r.ReasonData.(RootRequireData)
			rootConstraints[d.PackageName] = d.Constraint
		case RulePackageRequires:
			l := r.ReasonData.(*pkg.Link)
			requireConstraints[l.Target] = append(requireConstraints[l.Target], l.Constraint)
		case RulePackageConflict:
			l := r.ReasonData.(*pkg.Link)
			if _, ok := conflictConstraints[l.Target]; !ok {
				conflictOrder = append(conflictOrder, l.Target)
			}
			conflictConstraints[l.Target] = append(conflictConstraints[l.Target], l.Constraint)
		}
	}
	var messages []string
	for _, name := range conflictOrder {
		constraints := conflictConstraints[name]
		_, isRoot := rootConstraints[name]
		reqs, isReq := requireConstraints[name]
		if !isRoot && !isReq {
			continue
		}
		var removals []policy.VersionReason
		for _, vr := range ctx.Pool.PolicyRemovedVersions(name, rootConstraints[name]) {
			if vr.Reason.PackageName() != name {
				continue
			}
			pv := semver.NewConstraint("==", vr.Version)
			anyMatch := func(list []semver.Constraint) bool {
				for _, c := range list {
					if c.Matches(pv) {
						return true
					}
				}
				return false
			}
			if anyMatch(constraints) || (isReq && !anyMatch(reqs)) {
				continue
			}
			removals = append(removals, vr)
		}
		if len(removals) == 0 {
			continue
		}
		plural := len(removals) > 1
		versions := newOrderedVersions()
		for _, r := range removals {
			versions.set(r.Version, r.Reason.PrettyVersion())
		}
		var pretty []string
		if ctx.Verbose {
			for _, r := range removals {
				pretty = append(pretty, r.Reason.PrettyVersion())
			}
		} else {
			pretty = condenseVersionList(versions.keys, versions, 1, 1)
		}
		w := map[bool][3]string{true: {"match", "are", "were"}, false: {"matches", "is", "was"}}[plural]
		messages = append(messages, "\n    - "+name+"["+strings.Join(pretty, ", ")+"] "+w[0]+" the requirements and "+w[1]+" not excluded by the conflicts above, but "+w[2]+" not loaded, "+formatPolicyRemovalReasons(ctx, removals, plural))
	}
	return strings.Join(messages, "")
}
