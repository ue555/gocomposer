// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Advisory/Auditor.php, CooldownAuditor.php and
// src/Composer/Util/PackageInfo.php (getViewSourceOrHomepageUrl).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"strings"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// Auditor::FORMAT_* and STATUS_*.
const (
	FormatTable   = "table"
	FormatPlain   = "plain"
	FormatJSON    = "json"
	FormatSummary = "summary"

	StatusOK     = 0
	StatusFailed = 1
)

// Formats is Auditor::FORMATS.
var Formats = []string{FormatTable, FormatPlain, FormatJSON, FormatSummary}

// CooldownEntry is one CooldownAuditor::collect() value.
type CooldownEntry struct {
	PrettyVersion string
	ReleaseDate   string
	AvailableIn   string
	Source        string
}

// CooldownAuditor is Composer\Advisory\CooldownAuditor.
type CooldownAuditor struct {
	Now time.Time
}

// NewCooldownAuditor: a zero now means the current time.
func NewCooldownAuditor(now time.Time) *CooldownAuditor {
	if now.IsZero() {
		now = timeNow()
	}
	return &CooldownAuditor{Now: now}
}

// Collect is CooldownAuditor::collect: package name => cooldown info for
// versions still within the cooldown.
func (c *CooldownAuditor) Collect(packages []*pkg.Package, cooldown *CooldownPolicyConfig) *OMap[CooldownEntry] {
	result := NewOMap[CooldownEntry]()
	if !cooldown.HasCooldown() {
		return result
	}
	for _, p := range packages {
		if p.AliasOf != nil {
			p = p.AliasOf
		}
		if isRootPackage(p) || isPlatformPackage(p.Name) || p.IsDev() || cooldown.IsIgnored(p, OperationAudit) {
			continue
		}
		releaseDate, ok := cooldown.GetEffectiveDate(p)
		if !ok || !cooldown.IsWithinCooldown(releaseDate, c.Now) {
			continue
		}
		source := "time"
		if _, ok := PublishedDate(p); ok {
			source = "published-time"
		}
		result.Set(p.Name, CooldownEntry{p.PrettyVersion(), formatAtom(releaseDate), cooldown.FormatTimeUntilAvailable(releaseDate, c.Now, false), source})
	}
	return result
}

// Auditor is Composer\Advisory\Auditor.
type Auditor struct {
	cooldownAuditor *CooldownAuditor
	// TableStyles are the formatter styles used for table width calculations
	// (ComposerStyles for Composer's ConsoleIO).
	TableStyles []string
}

// NewAuditor mirrors the constructor (nil = a CooldownAuditor at "now").
func NewAuditor(cooldownAuditor *CooldownAuditor) *Auditor {
	if cooldownAuditor == nil {
		cooldownAuditor = NewCooldownAuditor(time.Time{})
	}
	return &Auditor{cooldownAuditor: cooldownAuditor, TableStyles: ComposerStyles}
}

// Audit is Auditor::audit: returns StatusOK or StatusFailed. filterListSource
// may be nil (no FilterListProviderSet).
func (a *Auditor) Audit(io IO, repoSet AdvisorySource, config *PolicyConfig, packages []*pkg.Package, format string, warningOnly bool, filterListSource FilterListSource) (int, error) {
	ignoreList := config.Advisories.GetIgnoreListForOperation(OperationAudit)
	ignoredSeverities := config.Advisories.GetIgnoreSeverityForOperation(OperationAudit)

	result, err := repoSet.GetMatchingSecurityAdvisories(packages, format == FormatSummary, config.IgnoreUnreachable.Audit)
	if err != nil {
		return 0, err
	}
	allAdvisories := result.Advisories
	unreachable := append([]string(nil), result.UnreachableRepos...)

	if format == FormatSummary && a.NeedsCompleteAdvisoryLoad(allAdvisories, ignoreList) {
		result, err = repoSet.GetMatchingSecurityAdvisories(packages, false, config.IgnoreUnreachable.Audit)
		if err != nil {
			return 0, err
		}
		allAdvisories = result.Advisories
		unreachable = append(unreachable, result.UnreachableRepos...)
	}
	advisories, ignoredAdvisories := a.ProcessAdvisories(allAdvisories, ignoreList, ignoredSeverities)

	abandonedCount := 0
	affectedPackagesCount := advisories.Len()
	var abandonedPackages []*pkg.Package
	if config.Abandoned.Audit != AuditIgnore {
		var canonical []*pkg.Package
		for _, p := range packages {
			if p.AliasOf == nil {
				canonical = append(canonical, p)
			}
		}
		abandonedPackages = a.FilterAbandonedPackages(canonical, config.Abandoned.GetFlatIgnoreForOperation(OperationAudit))
		if config.Abandoned.Audit == AuditFail {
			abandonedCount = len(abandonedPackages)
		}
	}

	filterAuditor := FilterListAuditor{}
	filtered := NewOMap[[]*FilterListEntry]()
	filteredCount := 0
	activeAuditLists := config.GetActiveAuditFilterLists()
	if filterListSource != nil && activeAuditLists.Len() > 0 {
		failing := map[string]bool{}
		activeAuditLists.Each(func(name string, l ListPolicy) {
			if l.AuditMode() == AuditFail {
				failing[name] = true
			}
		})
		fr, err := filterAuditor.CollectFilterLists(packages, filterListSource, activeAuditLists.Keys(), config.IgnoreUnreachable.Audit)
		if err != nil {
			return 0, err
		}
		unreachable = append(unreachable, fr.UnreachableRepos...)
		for _, p := range packages {
			for _, entry := range filterAuditor.GetMatchingAuditEntries(p, fr.Filter, config) {
				cur, _ := filtered.Get(p.Name)
				dup := false
				for _, e := range cur {
					if e == entry {
						dup = true
					}
				}
				if dup {
					continue
				}
				filtered.Set(p.Name, append(cur, entry))
				if failing[entry.ListName] {
					filteredCount++
				}
			}
		}
	}

	cooldownPackages := NewOMap[CooldownEntry]()
	cooldownCount := 0
	if config.Cooldown.Audit != AuditIgnore {
		cooldownPackages = a.cooldownAuditor.Collect(packages, config.Cooldown)
		if config.Cooldown.Audit == AuditFail {
			cooldownCount = cooldownPackages.Len()
		}
	}

	auditResult := StatusOK
	if affectedPackagesCount > 0 || abandonedCount > 0 || filteredCount > 0 || cooldownCount > 0 {
		auditResult = StatusFailed
	}

	if format == FormatJSON {
		json := phpjson.NewMap()
		json.Set("advisories", advisoryMapJSON(advisories))
		if ignoredAdvisories.Len() > 0 {
			json.Set("ignored-advisories", advisoryMapJSON(ignoredAdvisories))
		}
		if len(unreachable) > 0 {
			list := make([]phpjson.Value, len(unreachable))
			for i, u := range unreachable {
				list[i] = u
			}
			json.Set("unreachable-repositories", list)
		}
		ab := phpjson.NewMap()
		for _, p := range abandonedPackages {
			if r := p.ReplacementPackage(); r != "" {
				ab.Set(p.PrettyName, r)
			} else {
				ab.Set(p.PrettyName, nil)
			}
		}
		json.Set("abandoned", ab)
		fm := phpjson.NewMap()
		filtered.Each(func(name string, entries []*FilterListEntry) {
			list := make([]phpjson.Value, len(entries))
			for i, e := range entries {
				list[i] = e.jsonValue()
			}
			fm.Set(name, list)
		})
		json.Set("filter", fm)
		cm := phpjson.NewMap()
		cooldownPackages.Each(func(name string, e CooldownEntry) {
			m := phpjson.NewMap()
			m.Set("prettyVersion", e.PrettyVersion)
			m.Set("releaseDate", e.ReleaseDate)
			m.Set("availableIn", e.AvailableIn)
			m.Set("source", e.Source)
			cm.Set(name, m)
		})
		json.Set("cooldown", cm)
		enc, err := phpjson.Encode(json, phpjson.ComposerFlags)
		if err != nil {
			return 0, err
		}
		write(io, enc)
		return auditResult, nil
	}

	errorOrWarn := "error"
	if warningOnly {
		errorOrWarn = "warning"
	}
	if affectedPackagesCount > 0 || ignoredAdvisories.Len() > 0 {
		passes := []struct {
			advisories *AdvisoryMap
			message    string
		}{
			{ignoredAdvisories, "<info>Found %d ignored security vulnerability advisor%s affecting %d package%s%s</info>"},
			{advisories, "<" + errorOrWarn + ">Found %d security vulnerability advisor%s affecting %d package%s%s</" + errorOrWarn + ">"},
		}
		for _, pass := range passes {
			pkgCount, total := countAdvisories(pass.advisories)
			if pkgCount > 0 {
				plurality := "ies"
				if total == 1 {
					plurality = "y"
				}
				pkgPlurality := "s"
				if pkgCount == 1 {
					pkgPlurality = ""
				}
				punctuation := ":"
				if format == FormatSummary {
					punctuation = "."
				}
				write(io, fmt.Sprintf(pass.message, total, plurality, pkgCount, pkgPlurality, punctuation))
				if err := a.outputAdvisories(io, pass.advisories, format); err != nil {
					return 0, err
				}
			}
		}
		if format == FormatSummary {
			write(io, `Run "composer audit" for a full list of advisories.`)
		}
	} else {
		write(io, "<info>No security vulnerability advisories found.</info>")
	}

	if len(unreachable) > 0 {
		writeError(io, "<warning>The following repositories were unreachable:</warning>")
		for _, r := range unreachable {
			writeError(io, "  - "+r)
		}
	}

	if len(abandonedPackages) > 0 && format != FormatSummary {
		a.outputAbandonedPackages(io, abandonedPackages, format)
	}

	if filtered.Len() > 0 {
		plurality := "s"
		if filtered.Len() == 1 {
			plurality = ""
		}
		punctuation := ":"
		if format == FormatSummary {
			punctuation = "."
		}
		style := "warning"
		if filteredCount > 0 {
			style = "error"
		}
		write(io, fmt.Sprintf("<%s>Found %d package%s matching filters%s</%s>", style, filtered.Len(), plurality, punctuation, style))
		if format != FormatSummary {
			a.outputFilteredPackages(io, filtered, format)
		}
	}

	if cooldownPackages.Len() > 0 {
		plurality := "s"
		if cooldownPackages.Len() == 1 {
			plurality = ""
		}
		punctuation := ":"
		if format == FormatSummary {
			punctuation = "."
		}
		style := "warning"
		if cooldownCount > 0 {
			style = "error"
		}
		write(io, fmt.Sprintf(`<%s>Found %d package%s still in the cooldown period configured in "policy.cooldown"%s</%s>`, style, cooldownPackages.Len(), plurality, punctuation, style))
		if format != FormatSummary {
			cooldownPackages.Each(func(name string, info CooldownEntry) {
				write(io, fmt.Sprintf("  - %s (%s) published %s, available in %s", name, info.PrettyVersion, info.ReleaseDate, info.AvailableIn))
				if info.Source == "time" {
					write(io, "    The cooldown policy relied on the package-supplied time field because the repository provides no published-time.")
				}
			})
		}
	}

	return auditResult, nil
}

func advisoryMapJSON(m *AdvisoryMap) *phpjson.Map {
	out := phpjson.NewMap()
	m.Each(func(name string, list []*Advisory) {
		l := make([]phpjson.Value, len(list))
		for i, a := range list {
			l[i] = a.JSONValue()
		}
		out.Set(name, l)
	})
	return out
}

// NeedsCompleteAdvisoryLoad is Auditor::needsCompleteAdvisoryLoad.
func (a *Auditor) NeedsCompleteAdvisoryLoad(advisories *AdvisoryMap, ignoreList *OMap[*string]) bool {
	if advisories.Len() == 0 {
		return false
	}
	allFull := true
	advisories.Each(func(_ string, list []*Advisory) {
		for _, ad := range list {
			if !ad.Full {
				allFull = false
			}
		}
	})
	if allFull {
		return false
	}
	for _, id := range ignoreList.Keys() {
		if !strings.HasPrefix(id, "PKSA-") {
			return true
		}
	}
	return false
}

// FilterAbandonedPackages is Auditor::filterAbandonedPackages.
func (a *Auditor) FilterAbandonedPackages(packages []*pkg.Package, ignoreAbandoned *OMap[*string]) []*pkg.Package {
	var filter interface{ MatchString(string) bool }
	if ignoreAbandoned.Len() != 0 {
		filter = PackageNamesToRegexp(ignoreAbandoned.Keys())
	}
	var out []*pkg.Package
	for _, p := range packages {
		if p.IsAbandoned() && (filter == nil || !filter.MatchString(p.Name)) {
			out = append(out, p)
		}
	}
	return out
}

// ProcessAdvisories is Auditor::processAdvisories: returns the active and
// the ignored advisories.
func (a *Auditor) ProcessAdvisories(all *AdvisoryMap, ignoreList, ignoredSeverities *OMap[*string]) (*AdvisoryMap, *AdvisoryMap) {
	if ignoreList.Len() == 0 && ignoredSeverities.Len() == 0 {
		return all, NewOMap[[]*Advisory]()
	}
	advisories := NewOMap[[]*Advisory]()
	ignored := NewOMap[[]*Advisory]()
	var ignoreReason *string
	all.Each(func(name string, list []*Advisory) {
		for _, adv := range list {
			active := true
			if r, ok := ignoreList.Get(name); ok {
				active = false
				ignoreReason = r
			}
			if r, ok := ignoreList.Get(adv.AdvisoryID); ok {
				active = false
				ignoreReason = r
			}
			if adv.Full {
				if adv.Severity != nil {
					if r, ok := ignoredSeverities.Get(*adv.Severity); ok {
						active = false
						if r != nil {
							ignoreReason = r
						} else {
							ignoreReason = strPtr(*adv.Severity + " severity is ignored")
						}
					}
				}
				if adv.CVE != nil {
					if r, ok := ignoreList.Get(*adv.CVE); ok {
						active = false
						ignoreReason = r
					}
				}
				for _, src := range adv.Sources {
					if r, ok := ignoreList.Get(toStr(coalesce(src, "remoteId", nil))); ok {
						active = false
						ignoreReason = r
						break
					}
				}
			}
			if active {
				appendTo(advisories, name, adv)
				continue
			}
			if adv.Full {
				adv = adv.ToIgnoredAdvisory(ignoreReason)
			}
			appendTo(ignored, name, adv)
		}
	})
	return advisories, ignored
}

func countAdvisories(m *AdvisoryMap) (int, int) {
	count := 0
	m.Each(func(_ string, l []*Advisory) { count += len(l) })
	return m.Len(), count
}

func (a *Auditor) outputAdvisories(io IO, advisories *AdvisoryMap, format string) error {
	switch format {
	case FormatTable:
		advisories.Each(func(_ string, list []*Advisory) {
			for _, adv := range list {
				labels, values := a.advisoryFields(adv)
				for i := range values {
					values[i] = Sanitize(values[i])
				}
				NewTable(a.TableStyles).SetHorizontal().SetHeaders(labels).AddRow(values).
					SetColumnWidth(1, 80).SetColumnMaxWidth(1, 80).Render(func(l string) { write(io, l) })
			}
		})
		return nil
	case FormatPlain:
		var lines []string
		first := true
		advisories.Each(func(_ string, list []*Advisory) {
			for _, adv := range list {
				if !first {
					lines = append(lines, "--------")
				}
				labels, values := a.advisoryFields(adv)
				for i, label := range labels {
					v := values[i]
					if label == "Title" || label == "Affected versions" {
						v = EscapeFormat(v)
					}
					lines = append(lines, label+": "+v)
				}
				first = false
			}
		})
		write(io, lines...)
		return nil
	case FormatSummary:
		return nil
	}
	return newErr(InvalidArgumentException, `Invalid format "`+format+`".`)
}

// advisoryFields is Auditor::getAdvisoryFields (labels and values in order).
func (a *Auditor) advisoryFields(adv *Advisory) ([]string, []string) {
	if !adv.Full {
		return []string{"Package", "Advisory ID", "Affected versions"},
			[]string{adv.PackageName, advisoryIDLink(adv), adv.AffectedVersions.PrettyString()}
	}
	severity := ""
	if adv.Severity != nil {
		severity = *adv.Severity
	}
	cve := "NO CVE"
	if adv.CVE != nil {
		cve = "<href=https://www.cve.org/CVERecord?id=" + *adv.CVE + ">" + *adv.CVE + "</>"
	}
	link := ""
	if adv.Link != nil {
		link = "<href=" + EscapeFormat(*adv.Link) + ">" + EscapeFormat(*adv.Link) + "</>"
	}
	labels := []string{"Package", "Severity", "Advisory ID", "CVE", "Title", "URL", "Affected versions", "Reported at"}
	values := []string{adv.PackageName, severity, advisoryIDLink(adv), cve, adv.Title, link, adv.AffectedVersions.PrettyString(), formatAtom(adv.ReportedAt)}
	if adv.Ignored {
		reason := "None specified"
		if adv.IgnoreReason != nil {
			reason = *adv.IgnoreReason
		}
		labels = append(labels, "Ignore reason")
		values = append(values, reason)
	}
	return labels, values
}

func advisoryIDLink(adv *Advisory) string {
	if strings.HasPrefix(adv.AdvisoryID, "PKSA-") {
		return "<href=https://packagist.org/security-advisories/" + adv.AdvisoryID + ">" + adv.AdvisoryID + "</>"
	}
	return adv.AdvisoryID
}

func (a *Auditor) outputAbandonedPackages(io IO, packages []*pkg.Package, format string) {
	plural := ""
	if len(packages) > 1 {
		plural = "s"
	}
	write(io, fmt.Sprintf("<error>Found %d abandoned package%s:</error>", len(packages), plural))
	if format == FormatPlain {
		for _, p := range packages {
			replacement := "No replacement was suggested"
			if r := p.ReplacementPackage(); r != "" {
				replacement = "Use " + r + " instead"
			}
			write(io, fmt.Sprintf("%s is abandoned. %s.", packageNameWithLink(p), replacement))
		}
		return
	}
	t := NewTable(a.TableStyles).SetHeaders([]string{"Abandoned Package", "Suggested Replacement"}).SetColumnWidth(1, 80).SetColumnMaxWidth(1, 80)
	for _, p := range packages {
		replacement := "none"
		if r := p.ReplacementPackage(); r != "" {
			replacement = r
		}
		t.AddRow([]string{Sanitize(packageNameWithLink(p)), Sanitize(replacement)})
	}
	t.Render(func(l string) { write(io, l) })
}

// viewSourceOrHomepageURL is PackageInfo::getViewSourceOrHomepageUrl ("" = null).
func viewSourceOrHomepageURL(p *pkg.Package) string {
	if src := p.Support().Str("source"); src != "" {
		return src
	}
	if u := p.SourceURL(); u != "" {
		return u
	}
	return p.Homepage()
}

func packageNameWithLink(p *pkg.Package) string {
	if u := viewSourceOrHomepageURL(p); u != "" {
		return "<href=" + EscapeFormat(u) + ">" + p.PrettyName + "</>"
	}
	return p.PrettyName
}

func (a *Auditor) outputFilteredPackages(io IO, filtered *OMap[[]*FilterListEntry], format string) {
	if format == FormatPlain {
		filtered.Each(func(_ string, entries []*FilterListEntry) {
			for _, e := range entries {
				parts := []string{e.PackageName + ` matched dependency policy "` + e.ListName + `"`}
				if e.Reason != nil {
					parts = append(parts, "Reason: "+*e.Reason)
				}
				if e.URL != nil {
					parts = append(parts, "URL: "+*e.URL)
				}
				if e.Source != nil {
					parts = append(parts, "Source: "+*e.Source)
				}
				write(io, strings.Join(parts, ". ")+".")
			}
		})
		return
	}
	t := NewTable(a.TableStyles).SetHeaders([]string{"Package", "Versions", "List", "URL", "Reason", "ID", "Source"}).SetColumnMaxWidth(4, 40)
	deref := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	filtered.Each(func(_ string, entries []*FilterListEntry) {
		for _, e := range entries {
			row := []string{e.PackageName, e.Constraint.PrettyString(), e.ListName, deref(e.URL), deref(e.Reason), deref(e.ID), deref(e.Source)}
			for i := range row {
				row[i] = Sanitize(row[i])
			}
			t.AddRow(row)
		}
	})
	t.Render(func(l string) { write(io, l) })
}
