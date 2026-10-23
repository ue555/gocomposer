package cli

// Port of Composer\Command\PackageDiscoveryTrait (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano),
// shared by the require and init commands, including the interactive package
// search ("Search for a package") and the similar-package selection.

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// discovery holds the PackageDiscoveryTrait state: the repositories
// searched ($this->repos, platform repository first) and how the minimum
// stability is determined.
type discovery struct {
	ctx          *Context
	repos        []repo.Repository
	platformRepo repo.Repository
	minStability func() string
}

var reTooStrict = regexp.MustCompile(`^\d+(\.\d+)?$`)

// determineRequirements is PackageDiscoveryTrait::determineRequirements.
func (d *discovery) determineRequirements(c *composer.Composer, requires []string, preferredStability string, useBest, fixed bool) ([]string, error) {
	io := d.ctx.IO
	if len(requires) == 0 {
		// the interactive package search (no question is answered in
		// non-interactive mode: the first ask returns null)
		return d.searchRequirements(c, preferredStability)
	}
	for _, r := range requires {
		if strings.ToLower(r) == "as" {
			return nil, &composer.UserError{Msg: fmt.Sprintf(`Cannot use "%s" as a separate argument. Quote the inline alias as one argument, e.g. "vendor/package:dev-main as 1.2.x-dev".`, r), Code: 1}
		}
	}
	var result []string
	for _, r := range parseNameVersionPairs(requires) {
		if r.HasVersion && reTooStrict.MatchString(r.Version) {
			io.WriteError(`<warning>The "` + r.Version + `" constraint for "` + r.Name + `" appears too strict and will likely not match what you want. See https://getcomposer.org/constraints</warning>`)
		}
		name, version := r.Name, r.Version
		if !r.HasVersion {
			n, v, err := d.findBestVersionAndName(c, name, preferredStability, fixed)
			if err != nil {
				return nil, err
			}
			name = n
			if useBest {
				version = v
				io.WriteError(fmt.Sprintf("Using version <info>%s</info> for <info>%s</info>", version, name))
			} else {
				version = "guess"
			}
		}
		result = append(result, name+" "+version)
	}
	return result, nil
}

func (d *discovery) repositorySet(minimumStability string) *solver.RepositorySet {
	rs := solver.NewRepositorySet(minimumStability, nil, nil, nil, nil, nil)
	for _, r := range d.repos {
		rs.AddRepository(r)
	}
	return rs
}

func (d *discovery) cooldownPolicy(c *composer.Composer) (*policy.CooldownPolicyConfig, error) {
	in := d.ctx.In
	cfg, err := composer.CreatePolicyConfig(c, in.Bool("no-blocking") || in.Bool("no-security-blocking"), in.Bool("no-cooldown"))
	if err != nil {
		return nil, err
	}
	if cfg.Enabled && cfg.Cooldown.ShouldBlock(policy.BlockScopeUpdate) {
		return cfg.Cooldown, nil
	}
	return nil, nil
}

func (d *discovery) newSelector(c *composer.Composer, rs *solver.RepositorySet) *composer.VersionSelector {
	vs := composer.NewVersionSelector(rs, d.platformRepo)
	if _, facts, _ := c.PHP(); facts != nil {
		parts := strings.SplitN(facts.PHPVersion, "-", 2)
		vp := strings.Split(parts[0], ".")
		if len(vp) >= 3 {
			vs.PHPVersion = strings.Join(vp[:3], ".")
		}
	}
	return vs
}

// findBestVersionAndName is PackageDiscoveryTrait::findBestVersionAndNameForPackage.
func (d *discovery) findBestVersionAndName(c *composer.Composer, name, preferredStability string, fixed bool) (string, string, error) {
	io := d.ctx.IO
	filter := platformFilter(d.ctx.In)
	rs := d.repositorySet(d.minStability())
	cooldown, err := d.cooldownPolicy(c)
	if err != nil {
		return "", "", err
	}
	vs := d.newSelector(c, rs)
	vs.SetCooldownPolicy(cooldown)
	effectiveMin := d.minStability()
	p, err := vs.FindBestCandidate(name, "", preferredStability, filter, 0, io, true)
	if err != nil {
		return "", "", err
	}
	if p != nil {
		if fixed {
			return p.PrettyName, p.PrettyVersion(), nil
		}
		return p.PrettyName, vs.FindRecommendedRequireVersion(p), nil
	}
	if filter.IsIgnored(name) {
		return name, "*", nil
	}
	if providers := rs.Providers(name); len(providers) > 0 {
		constraint := "*"
		if io.Interactive {
			v, err := io.AskAndValidate(`Package "<info>`+name+`</info>" does not exist but is provided by `+fmt.Sprint(len(providers))+` packages. Which version constraint would you like to use? [<info>*</info>] `, func(v any) (any, error) {
				s, _ := cio.AnswerString(v)
				if _, err := semver.ParseConstraints(s); err != nil {
					return nil, &cio.InvalidArgumentError{Msg: err.Error()}
				}
				return v, nil
			}, 3, "*")
			if err != nil {
				return "", "", err
			}
			constraint, _ = cio.AnswerString(v)
		}
		return name, constraint, nil
	}
	unfiltered := d.newSelector(c, rs)
	if cooldown != nil {
		if cand, err := unfiltered.FindBestCandidate(name, "", preferredStability, filter, 0, nil, true); err == nil && cand != nil {
			wait := ""
			if w := vs.FirstClearingWithheld; w != nil {
				wait = " (version " + w.Package.PrettyVersion() + " becomes available in " + cooldown.FormatTimeUntilAvailable(w.ReleaseDate, time.Now(), false) + ")"
			}
			return "", "", &composer.UserError{Msg: fmt.Sprintf(`All versions of package %s are still in the cooldown period configured in "policy.cooldown"%s. `+policy.CooldownRemedy, name, wait), Code: 1}
		}
	}
	_, ignoreAll := filter.(solver.IgnoreAll)
	if !ignoreAll {
		if cand, err := unfiltered.FindBestCandidate(name, "", preferredStability, solver.IgnoreAll{}, 0, nil, true); err == nil && cand != nil {
			return "", "", &composer.UserError{Msg: fmt.Sprintf("Package %s has requirements incompatible with your PHP version, PHP extensions and Composer version"+d.platformExceptionDetails(cand), name), Code: 1}
		}
	}
	if cand, err := unfiltered.FindBestCandidate(name, "", preferredStability, filter, solver.AllowUnacceptableStabilities, nil, true); err == nil && cand != nil {
		if all, err := unfiltered.FindBestCandidate(name, "", preferredStability, filter, solver.AllowShadowedRepositories, nil, true); err == nil && all != nil {
			return "", "", &composer.UserError{Msg: "Package " + name + " exists in " + repoNameOf(all) + " and " + repoNameOf(cand) + " which has a higher repository priority. The packages from the higher priority repository do not match your minimum-stability and are therefore not installable. That repository is canonical so the lower priority repo's packages are not installable. See https://getcomposer.org/repoprio for details and assistance.", Code: 1}
		}
		return "", "", &composer.UserError{Msg: fmt.Sprintf("Could not find a version of package %s matching your minimum-stability (%s). Require it with an explicit version constraint allowing its desired stability.", name, effectiveMin), Code: 1}
	}
	if !ignoreAll {
		if cand, err := unfiltered.FindBestCandidate(name, "", preferredStability, solver.IgnoreAll{}, solver.AllowUnacceptableStabilities, nil, true); err == nil && cand != nil {
			additional := ""
			if c2, err := unfiltered.FindBestCandidate(name, "", preferredStability, solver.IgnoreAll{}, 0, nil, true); err == nil && c2 == nil {
				additional = "\n\nAdditionally, the package was only found with a stability of \"" + cand.Stability() + "\" while your minimum stability is \"" + effectiveMin + "\"."
			}
			return "", "", &composer.UserError{Msg: fmt.Sprintf("Could not find package %s in any version matching your PHP version, PHP extensions and Composer version"+d.platformExceptionDetails(cand)+"%s", name, additional), Code: 1}
		}
	}
	similar, err := d.findSimilar(c, name)
	if err != nil {
		return "", "", err
	}
	if len(similar) > 0 {
		for _, s := range similar {
			if s == name {
				return "", "", &composer.UserError{Msg: fmt.Sprintf("Could not find package %s. It was however found via repository search, which indicates a consistency issue with the repository.", name), Code: 1}
			}
		}
		if io.Interactive {
			choices := make([]cio.Choice, len(similar))
			for i, s := range similar {
				choices[i] = cio.Choice{Key: fmt.Sprint(i), Value: s}
			}
			res, err := io.Select("<error>Could not find package "+name+".</error>\nPick one of these or leave empty to abort:", choices, false, false, 1, `Value "%s" is invalid`, false)
			if err != nil {
				return "", "", err
			}
			if k, ok := res.(string); ok {
				idx := 0
				fmt.Sscan(k, &idx)
				return d.findBestVersionAndName(c, similar[idx], preferredStability, fixed)
			}
		}
		which := "this"
		if len(similar) > 1 {
			which = "one of these"
		}
		return "", "", &composer.UserError{Msg: fmt.Sprintf("Could not find package %s.\n\nDid you mean "+which+"?\n    %s", name, strings.Join(similar, "\n    ")), Code: 1}
	}
	return "", "", &composer.UserError{Msg: fmt.Sprintf("Could not find a matching version of package %s. Check the package spelling, your version constraint and that the package is available in a stability which matches your minimum-stability (%s).", name, effectiveMin), Code: 1}
}

func repoNameOf(p *pkg.Package) string {
	if p.Repo != nil {
		return p.Repo.RepoName()
	}
	return ""
}

// platformExceptionDetails is PackageDiscoveryTrait::getPlatformExceptionDetails.
func (d *discovery) platformExceptionDetails(cand *pkg.Package) string {
	if d.platformRepo == nil {
		return ""
	}
	var details []string
	for _, l := range cand.Requires() {
		if !repo.IsPlatformPackage(l.Target) {
			continue
		}
		pp := repo.FindPackage(d.platformRepo, l.Target, &semver.MatchAll{})
		if pp == nil {
			if disabled, ok := d.platformRepo.(interface{ IsPlatformPackageDisabled(string) bool }); ok && disabled.IsPlatformPackageDisabled(l.Target) {
				details = append(details, cand.PrettyName+" "+cand.PrettyVersion()+" requires "+l.Target+" "+l.PrettyConstraint+` but it is disabled by your platform config. Enable it again with "composer config platform.`+l.Target+` --unset".`)
			} else {
				details = append(details, cand.PrettyName+" "+cand.PrettyVersion()+" requires "+l.Target+" "+l.PrettyConstraint+" but it is not present.")
			}
			continue
		}
		if !l.Constraint.Matches(semver.NewConstraint("==", pp.Version())) {
			v := pp.PrettyVersion()
			if e := pp.Extra(); e != nil && e.Has("config.platform") {
				v += " (" + pp.Description() + ")"
			}
			details = append(details, cand.PrettyName+" "+cand.PrettyVersion()+" requires "+l.Target+" "+l.PrettyConstraint+" which does not match your installed version "+v+".")
		}
	}
	if len(details) == 0 {
		return ""
	}
	return ":\n  - " + strings.Join(details, "\n  - ")
}

// findSimilar is PackageDiscoveryTrait::findSimilar.
func (d *discovery) findSimilar(c *composer.Composer, name string) ([]string, error) {
	type res struct {
		name string
		dist int
	}
	var results []repo.SearchResult
	for _, r := range d.repos {
		found, err := repo.Search(r, name)
		if err != nil {
			// ignore search errors
			return nil, nil
		}
		results = append(results, found...)
	}
	var list []res
	seen := map[string]int{}
	for _, r := range results {
		if repo.FindPackage(c.LocalRepo, r.Name, &semver.MatchAll{}) != nil {
			continue
		}
		d := levenshtein(name, r.Name)
		if i, ok := seen[r.Name]; ok {
			list[i].dist = d
			continue
		}
		seen[r.Name] = len(list)
		list = append(list, res{r.Name, d})
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].dist < list[j].dist })
	var out []string
	for i, r := range list {
		if i == 5 {
			break
		}
		out = append(out, r.name)
	}
	return out, nil
}

var rePackageSelection = regexp.MustCompile(`^\s*(?P<name>[\S/]+)(?:\s+(?P<version>\S+))?\s*$`)

// existingPackages are the names of the installed packages
// ($composer->getRepositoryManager()->getLocalRepository()).
func (d *discovery) existingPackages(c *composer.Composer) []string {
	if c == nil || c.LocalRepo == nil {
		return nil
	}
	ps, err := c.LocalRepo.Packages()
	if err != nil {
		return nil
	}
	var out []string
	for _, p := range ps {
		out = append(out, p.Name)
	}
	return out
}

// search is CompositeRepository::search($query) over $this->repos.
func (d *discovery) search(query string) ([]repo.SearchEntry, error) {
	var out []repo.SearchEntry
	for _, r := range d.repos {
		res, err := repo.SearchWithMode(r, query, repo.SearchFulltext, nil)
		if err != nil {
			return nil, err
		}
		out = append(out, res...)
	}
	return out, nil
}

// providers is CompositeRepository::getProviders.
func (d *discovery) providers(name string) []repo.ProviderInfo {
	var out []repo.ProviderInfo
	idx := map[string]int{}
	for _, r := range d.repos {
		ps, _ := r.Providers(name)
		for _, p := range ps {
			if i, ok := idx[p.Name]; ok {
				out[i] = p
				continue
			}
			idx[p.Name] = len(out)
			out = append(out, p)
		}
	}
	return out
}

// searchRequirements is the interactive part of determineRequirements: the
// user searches packages, picks one from the matches and gives a version
// constraint (or lets the best one be chosen).
func (d *discovery) searchRequirements(c *composer.Composer, preferredStability string) ([]string, error) {
	io := d.ctx.IO
	existing := d.existingPackages(c)
	isExisting := func(n string) bool {
		for _, e := range existing {
			if e == n {
				return true
			}
		}
		return false
	}
	var requires []string
	for {
		v, err := io.Ask("Search for a package: ", nil)
		if err != nil {
			return nil, err
		}
		query, ok := cio.AnswerString(v)
		if !ok {
			break
		}
		found, err := d.search(query)
		if err != nil {
			return nil, err
		}
		if len(found) == 0 {
			continue
		}
		var matches []repo.SearchEntry
		for _, m := range found {
			if !isExisting(m.Name) {
				matches = append(matches, m)
			}
		}
		exact := false
		for _, m := range matches {
			if m.Name == query {
				exact = true
				break
			}
		}
		pkgName, hasPkg := query, true
		if !exact {
			if len(d.providers(query)) > 0 {
				empty := ""
				matches = append([]repo.SearchEntry{{Name: query, Description: &empty}}, matches...)
			}
			var choices []string
			for pos, m := range matches {
				abandoned := ""
				if m.Abandoned != nil && phpjson.Truthy(m.Abandoned) {
					replacement := "No replacement was suggested"
					if s, ok := m.Abandoned.(string); ok {
						replacement = "Use " + s + " instead"
					}
					abandoned = "<warning>Abandoned. " + replacement + ".</warning>"
				}
				choices = append(choices, fmt.Sprintf(" <info>%5s</info> %s %s", fmt.Sprintf("[%d]", pos), m.Name, abandoned))
			}
			io.WriteError("")
			io.WriteError(fmt.Sprintf("Found <info>%d</info> packages matching <info>%s</info>", len(matches), query))
			io.WriteError("")
			for _, ch := range choices {
				io.WriteError(ch)
			}
			io.WriteError("")
			sel, err := io.AskAndValidate("Enter package # to add, or the complete package name if it is not listed: ", func(v any) (any, error) {
				selection, _ := cio.AnswerString(v)
				if selection == "" {
					return false, nil
				}
				if f, ok := phpNumeric(selection); ok {
					if i := int(f); i >= 0 && i < len(matches) {
						return matches[i].Name, nil
					}
				}
				if m := rePackageSelection.FindStringSubmatch(selection); m != nil {
					if ver := m[2]; ver != "" {
						if _, err := semver.ParseConstraints(ver); err != nil {
							return nil, &cio.InvalidArgumentError{Msg: err.Error()}
						}
						return m[1] + " " + ver, nil
					}
					return m[1], nil
				}
				return nil, &cio.InvalidArgumentError{Msg: "Not a valid selection"}
			}, 3, "")
			if err != nil {
				return nil, err
			}
			if s, ok := sel.(string); ok {
				pkgName = s
			} else {
				hasPkg = false
			}
		}
		if hasPkg && !strings.Contains(pkgName, " ") {
			cv, err := io.AskAndValidate("Enter the version constraint to require (or leave blank to use the latest version): ", func(v any) (any, error) {
				s, _ := cio.AnswerString(v)
				s = strings.Trim(s, " \t\n\r\x00\x0B")
				if s == "" {
					return false, nil
				}
				return s, nil
			}, 3, "")
			if err != nil {
				return nil, err
			}
			constraint, isStr := cv.(string)
			if !isStr {
				_, best, err := d.findBestVersionAndName(c, pkgName, preferredStability, false)
				if err != nil {
					return nil, err
				}
				constraint = best
				io.WriteError(fmt.Sprintf("Using version <info>%s</info> for <info>%s</info>", constraint, pkgName))
			}
			pkgName += " " + constraint
		}
		if hasPkg {
			requires = append(requires, pkgName)
			existing = append(existing, strings.SplitN(pkgName, " ", 2)[0])
		}
	}
	return requires, nil
}

var rePHPNumeric = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?[ \t\n\r\v\f]*$`)

// phpNumeric is is_numeric() with the (float) value.
func phpNumeric(s string) (float64, bool) {
	if !rePHPNumeric.MatchString(s) {
		return 0, false
	}
	var f float64
	_, err := fmt.Sscan(strings.TrimSpace(s), &f)
	return f, err == nil
}
