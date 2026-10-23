package composer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
)

func timeNow() time.Time { return time.Now() }

type suggestion struct{ source, target, reason string }

// SuggestedPackagesReporter ports Composer\Installer\SuggestedPackagesReporter.
type SuggestedPackagesReporter struct {
	io          *cio.IO
	suggestions []suggestion
}

func NewSuggestedPackagesReporter(io *cio.IO) *SuggestedPackagesReporter {
	return &SuggestedPackagesReporter{io: io}
}

func (r *SuggestedPackagesReporter) AddSuggestionsFromPackage(p *pkg.Package) {
	s := p.Suggests()
	if s == nil {
		return
	}
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		reason := fmt.Sprint(v)
		if str, ok := v.(string); ok {
			reason = str
		}
		r.suggestions = append(r.suggestions, suggestion{p.PrettyName, k, reason})
	}
}

// providers maps a name to the packages providing it in the installed repo.
type providers map[string][]string

func (r *SuggestedPackagesReporter) filtered(installed providers) []suggestion {
	var out []suggestion
	for _, s := range r.suggestions {
		providedByOthers := false
		for _, prov := range installed[s.target] {
			if prov != strings.ToLower(s.source) {
				providedByOthers = true
				break
			}
		}
		if providedByOthers {
			continue
		}
		out = append(out, s)
	}
	return out
}

// OutputMinimalistic prints the summary line after an update.
func (r *SuggestedPackagesReporter) OutputMinimalistic(installedNames map[string]bool) {
	installed := providers{}
	for n := range installedNames {
		installed[n] = append(installed[n], "")
	}
	if f := r.filtered(installed); len(f) > 0 {
		r.io.WriteError(fmt.Sprintf("<info>%d package suggestions were added by new dependencies, use `composer suggest` to see details.</info>", len(f)))
	}
}

// Output modes of SuggestedPackagesReporter.
const (
	SuggestModeList         = 1
	SuggestModeByPackage    = 2
	SuggestModeBySuggestion = 4
)

func providersOf(installed *repo.Installed) providers {
	out := providers{}
	if installed == nil {
		return out
	}
	for _, p := range installed.Packages() {
		for _, n := range p.Names(true) {
			out[n] = append(out[n], p.Name)
		}
	}
	return out
}

// filteredSuggestions is getFilteredSuggestions().
func (r *SuggestedPackagesReporter) filteredSuggestions(installed *repo.Installed, onlyDependentsOf *pkg.Package) []suggestion {
	var prov providers
	if installed != nil && len(r.suggestions) > 0 {
		prov = providersOf(installed)
	}
	var sourceFilter []string
	if onlyDependentsOf != nil {
		for _, l := range append(append(pkg.Links(nil), onlyDependentsOf.Requires()...), onlyDependentsOf.DevRequires()...) {
			sourceFilter = append(sourceFilter, l.Target)
		}
		sourceFilter = append(sourceFilter, onlyDependentsOf.Name)
	}
	var out []suggestion
	for _, s := range r.filtered(prov) {
		if len(sourceFilter) > 0 {
			found := false
			for _, f := range sourceFilter {
				if f == s.source {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		out = append(out, s)
	}
	return out
}

var reCntrl = regexp.MustCompile(`[\x00-\x1f\x7f]`)

func escapeSuggestOutput(s string) string {
	s = reCntrl.ReplaceAllString(strings.ReplaceAll(s, "\n", " "), "")
	return strings.ReplaceAll(s, "<", `\<`)
}

type kv struct{ k, v string }

// orderedGroups emulates $a[$x][$y] = $reason followed by ksort($a).
func orderedGroups(list []suggestion, byTarget bool) ([]string, map[string][]kv) {
	groups := map[string][]kv{}
	var keys []string
	for _, s := range list {
		outer, inner := s.source, s.target
		if byTarget {
			outer, inner = s.target, s.source
		}
		if _, ok := groups[outer]; !ok {
			keys = append(keys, outer)
		}
		g := groups[outer]
		replaced := false
		for i := range g {
			if g[i].k == inner {
				g[i].v = s.reason
				replaced = true
			}
		}
		if !replaced {
			g = append(g, kv{inner, s.reason})
		}
		groups[outer] = g
	}
	sort.SliceStable(keys, func(i, j int) bool { return phpjson.PHPCompareKeys(keys[i], keys[j]) < 0 })
	return keys, groups
}

// Output is SuggestedPackagesReporter::output.
func (r *SuggestedPackagesReporter) Output(mode int, installed *repo.Installed, onlyDependentsOf *pkg.Package) {
	list := r.filteredSuggestions(installed, onlyDependentsOf)
	srcKeys, suggesters := orderedGroups(list, false)
	tgtKeys, suggested := orderedGroups(list, true)
	if mode&SuggestModeList != 0 {
		for _, name := range tgtKeys {
			r.io.Write("<info>" + name + "</info>")
		}
		return
	}
	line := func(name, reason string) string {
		if reason != "" && reason != "0" {
			return " - <info>" + name + "</info>: " + escapeSuggestOutput(reason)
		}
		return " - <info>" + name + "</info>"
	}
	if mode&SuggestModeByPackage != 0 {
		for _, s := range srcKeys {
			r.io.Write("<comment>" + s + "</comment> suggests:")
			for _, e := range suggesters[s] {
				r.io.Write(line(e.k, e.v))
			}
			r.io.Write("")
		}
	}
	if mode&SuggestModeBySuggestion != 0 {
		if mode&SuggestModeByPackage != 0 {
			r.io.Write(strings.Repeat("-", 78))
		}
		for _, t := range tgtKeys {
			r.io.Write("<comment>" + t + "</comment> is suggested by:")
			for _, e := range suggested[t] {
				r.io.Write(line(e.k, e.v))
			}
			r.io.Write("")
		}
	}
	if onlyDependentsOf != nil {
		if diff := len(r.filteredSuggestions(installed, nil)) - len(list); diff != 0 {
			r.io.Write(fmt.Sprintf("<info>%d additional suggestions</info> by transitive dependencies can be shown with <info>--all</info>", diff))
		}
	}
}
