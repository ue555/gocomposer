package cli

// Port of the latest-version lookup and package listing of
// Composer\Command\ShowCommand (execute() list part, printPackages,
// findLatestPackage, getUpdateStatus, getCooldownPolicyConfig,
// getRepositorySet) (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// showWithheld is findLatestPackage()'s 'withheld' entry: the newer release
// still in the cooldown period that becomes installable first.
type showWithheld struct {
	p           *pkg.Package
	releaseDate time.Time
	availableIn string
}

const atomFormat = "2006-01-02T15:04:05-07:00"

func releaseTime(p *pkg.Package) (time.Time, bool) {
	if p.ReleaseDate() == "" {
		return time.Time{}, false
	}
	return pkg.ParsePHPDate(p.ReleaseDate())
}

// cooldownPolicyConfig is ShowCommand::getCooldownPolicyConfig.
func (sc *showCtx) cooldownPolicyConfig() (*policy.CooldownPolicyConfig, error) {
	in := sc.ctx.In
	if override, ok := in.String("cooldown-period"); ok {
		period, err := policy.ParseDuration(override)
		if err != nil {
			return nil, &composer.UserError{Msg: `Invalid --cooldown-period value "` + override + `", use an integer number of seconds or a duration like "7 days".`, Code: 1, Prev: err}
		}
		cfg, err := composer.CreatePolicyConfig(sc.c, false, false)
		if err != nil {
			return nil, err
		}
		cd := policy.NewCooldownPolicyConfig(true, cfg.Cooldown.Audit, cfg.Cooldown.Ignore, period)
		if cd.HasCooldown() {
			return cd, nil
		}
		return nil, nil
	}
	// BaseCommand::getBlockingCooldownPolicy
	cfg, err := composer.CreatePolicyConfig(sc.c, false, false)
	if err != nil {
		return nil, err
	}
	if cfg.Enabled && cfg.Cooldown.ShouldBlock(policy.BlockScopeUpdate) {
		return cfg.Cooldown, nil
	}
	return nil, nil
}

// repositorySet is ShowCommand::getRepositorySet.
func (sc *showCtx) repositorySet() *solver.RepositorySet {
	if sc.latestRS == nil {
		root := sc.c.Root()
		min := root.MinimumStability
		if min == "" {
			min = "stable"
		}
		rs := solver.NewRepositorySet(min, root.StabilityFlags, nil, nil, nil, nil)
		for _, r := range sc.c.Repositories {
			rs.AddRepository(r)
		}
		sc.latestRS = rs
	}
	return sc.latestRS
}

var reMajorOnly = regexp.MustCompile(`^((?:0\.)+)?(\d+)\.`)
var reTrailingZeros = regexp.MustCompile(`(\.0)+$`)

// findLatestPackage is ShowCommand::findLatestPackage.
func (sc *showCtx) findLatestPackage(p *pkg.Package, platformRepo repo.Repository, majorOnly, minorOnly, patchOnly bool, filter solver.PlatformFilter, cooldown *policy.CooldownPolicyConfig) (*pkg.Package, *showWithheld, error) {
	io := sc.ctx.IO
	name := p.Name
	vs := composer.NewVersionSelector(sc.repositorySet(), platformRepo)
	vs.SetCooldownPolicy(cooldown)
	root := sc.c.Root()
	stability := root.MinimumStability
	if stability == "" {
		stability = "stable"
	}
	if flag, ok := root.StabilityFlags[name]; ok {
		for sname, v := range pkg.Stabilities {
			if v == flag {
				stability = sname
				break
			}
		}
	}
	bestStability := stability
	if root.PreferStable {
		bestStability = p.Stability()
	}
	targetVersion := ""
	if strings.HasPrefix(p.Version(), "dev-") {
		targetVersion = p.Version()
		// dev-x branches are considered to be on the latest major version always
		if majorOnly {
			return nil, nil, nil
		}
	}
	if targetVersion == "" {
		if majorOnly {
			if m := reMajorOnly.FindStringSubmatch(p.Version()); m != nil {
				first, _ := strconv.Atoi(m[2])
				targetVersion = ">=" + m[1] + strconv.Itoa(first+1) + ",<9999999-dev"
			}
		}
		if minorOnly {
			targetVersion = "^" + p.Version()
		}
		if patchOnly {
			trimmed := reTrailingZeros.ReplaceAllString(p.Version(), "")
			partsNeeded := 3
			if strings.HasPrefix(trimmed, "0") {
				partsNeeded = 4
			}
			for strings.Count(trimmed, ".")+1 < partsNeeded {
				trimmed += ".0"
			}
			targetVersion = "~" + trimmed
		}
	}
	showWarnings := io.IsVerbose()
	if !showWarnings {
		vs.ShowWarningsFunc = func(candidate *pkg.Package) bool {
			if strings.HasPrefix(candidate.Version(), "dev-") || strings.HasPrefix(p.Version(), "dev-") {
				return false
			}
			return semver.VersionCompare(candidate.Version(), p.Version()) <= 0
		}
	}
	candidate, err := vs.FindBestCandidate(name, targetVersion, bestStability, filter, 0, io, showWarnings)
	if err != nil {
		return nil, nil, err
	}
	for candidate != nil && candidate.AliasOf != nil {
		candidate = candidate.AliasOf
	}
	latest := candidate
	// versions exist but none is installable on this platform, so there is nothing to update to
	if latest == nil && len(vs.PlatformRejected) > 0 {
		latest = p
	}
	var withheld *showWithheld
	if cooldown != nil && !strings.HasPrefix(p.Version(), "dev-") {
		baseline := p
		if latest != nil && semver.VersionCompare(latest.Version(), p.Version()) > 0 {
			baseline = latest
		}
		installedIsWithheld := false
		for _, wc := range vs.WithheldCandidates {
			wp := wc.Package
			for wp.AliasOf != nil {
				wp = wp.AliasOf
			}
			if wp.Version() == p.Version() {
				installedIsWithheld = true
			}
			if !(semver.VersionCompare(wp.Version(), baseline.Version()) > 0) {
				continue
			}
			if withheld == nil || wc.ReleaseDate.Before(withheld.releaseDate) {
				withheld = &showWithheld{p: wp, releaseDate: wc.ReleaseDate}
			}
		}
		if installedIsWithheld && (latest == nil || semver.VersionCompare(latest.Version(), p.Version()) < 0) {
			latest = p
		}
		if withheld != nil {
			withheld.availableIn = cooldown.FormatTimeUntilAvailable(withheld.releaseDate, time.Now(), false)
		}
	}
	return latest, withheld, nil
}

// semverSatisfies is Composer\Semver\Semver::satisfies.
func semverSatisfies(version, constraints string) bool {
	n, err := semver.Normalize(version)
	if err != nil {
		return false
	}
	c, err := semver.ParseConstraints(constraints)
	if err != nil {
		return false
	}
	return c.Matches(semver.NewConstraint("==", n))
}

// updateStatus is ShowCommand::getUpdateStatus.
func updateStatus(latest, p *pkg.Package) string {
	if latest.FullPrettyVersion(true) == p.FullPrettyVersion(true) {
		return "up-to-date"
	}
	constraint := p.Version()
	if !strings.HasPrefix(constraint, "dev-") {
		constraint = "^" + constraint
	}
	if latest.Version() != "" && latest.Version() != "0" && semverSatisfies(latest.Version(), constraint) {
		return "semver-safe-update"
	}
	return "update-possible"
}

var updateStatusReplacer = strings.NewReplacer("up-to-date", "info", "semver-safe-update", "highlight", "update-possible", "comment", "cooldown", "")
var updateStatusSigns = strings.NewReplacer("up-to-date", "=", "semver-safe-update", "!", "update-possible", "~", "cooldown", "c")

// updateStatusToVersionStyle is ShowCommand::updateStatusToVersionStyle.
func updateStatusToVersionStyle(status string) string { return updateStatusReplacer.Replace(status) }

type showMeta struct {
	nameLength, versionLength, latestLength, releaseDateLength int
	writeLatest, writeReleaseDate                              bool
}

// listPackages is the "list packages" part of ShowCommand::execute.
func (sc *showCtx) listPackages(packages map[string]*showTyped, format string, platformRepo repo.Repository, filter solver.PlatformFilter, cooldown *policy.CooldownPolicyConfig) (int, error) {
	io, in := sc.ctx.IO, sc.ctx.In
	showAllTypes := in.Bool("all")
	showLatest := in.Bool("latest")
	majorOnly, minorOnly, patchOnly := in.Bool("major-only"), in.Bool("minor-only"), in.Bool("patch-only")
	var ignores []string
	for _, s := range in.Strings("ignore") {
		ignores = append(ignores, strings.ToLower(s))
	}
	ignoredRe := solver.PackageNamesToRegexp(ignores)
	indent := ""
	if showAllTypes {
		indent = "  "
	}
	latestPackages := map[string]*pkg.Package{}
	withheldPackages := map[string]*showWithheld{}
	exitCode := 0
	viewData := phpjson.NewMap()
	metas := map[string]showMeta{}
	var typeOrder []string
	writeVersion, writeDescription := false, false
	nameOnly, pathOpt := in.Bool("name-only"), in.Bool("path")
	rootRequires := sc.rootRequires()

	latestLookups := map[string][]*pkg.Package{}
	if showLatest {
		for _, typ := range []string{"platform", "locked", "installed"} {
			t := packages[typ]
			if t == nil {
				continue
			}
			for _, name := range t.order {
				if e := t.m[name]; e.p != nil && !ignoredRe.MatchString(e.p.PrettyName) {
					latestLookups[typ] = append(latestLookups[typ], e.p)
				}
			}
		}
	}

	for _, tv := range []struct {
		typ         string
		showVersion bool
	}{{"platform", true}, {"locked", true}, {"available", false}, {"installed", true}} {
		t := packages[tv.typ]
		if t == nil {
			continue
		}
		names := append([]string(nil), t.order...)
		sort.SliceStable(names, func(i, j int) bool { return phpjson.PHPCompareKeys(names[i], names[j]) < 0 })
		var m showMeta
		for _, p := range latestLookups[tv.typ] {
			delete(withheldPackages, p.PrettyName)
			l, w, err := sc.findLatestPackage(p, platformRepo, majorOnly, minorOnly, patchOnly, filter, cooldown)
			if err != nil {
				return 1, err
			}
			if l != nil {
				latestPackages[p.PrettyName] = l
			}
			if w != nil {
				withheldPackages[p.PrettyName] = w
			}
		}
		writePath := !nameOnly && pathOpt
		writeVersion = !nameOnly && !pathOpt && tv.showVersion
		writeLatest := writeVersion && showLatest
		writeDescription = !nameOnly && !pathOpt
		writeReleaseDate := writeLatest && (in.Bool("sort-by-age") || format == "json")
		hasOutdated := false
		if in.Bool("sort-by-age") {
			sort.SliceStable(names, func(i, j int) bool {
				a, b := t.m[names[i]].p, t.m[names[j]].p
				if a == nil || b == nil {
					return false
				}
				ta, oka := releaseTime(a)
				tb, okb := releaseTime(b)
				switch {
				case !oka && !okb:
					return false
				case !oka:
					return true
				case !okb:
					return false
				}
				return ta.Before(tb)
			})
		}
		list := []phpjson.Value{}
		for _, name := range names {
			e := t.m[name]
			d := phpjson.NewMap()
			if e.p == nil {
				d.Set("name", e.name)
				if len(e.name) > m.nameLength {
					m.nameLength = len(e.name)
				}
				list = append(list, d)
				continue
			}
			p := e.p
			var latest *pkg.Package
			if showLatest {
				latest = latestPackages[p.PrettyName]
			}
			upToDate := latest != nil && latest.FullPrettyVersion(true) == p.FullPrettyVersion(true) && !latest.IsAbandoned()
			upToDate = upToDate || (latest == nil && majorOnly)
			ignored := ignoredRe.MatchString(p.PrettyName)
			withheld := withheldPackages[p.PrettyName]
			if in.Bool("outdated") && (ignored || (upToDate && withheld == nil)) {
				continue
			}
			if (in.Bool("outdated") || in.Bool("strict")) && !upToDate {
				hasOutdated = true
			}
			d.Set("name", p.PrettyName)
			d.Set("direct-dependency", inList(p.Name, rootRequires))
			if format != "json" || !nameOnly {
				d.Set("homepage", nullable(p.Homepage()))
				d.Set("source", nullable(viewSourceURL(p)))
			}
			if len(p.PrettyName) > m.nameLength {
				m.nameLength = len(p.PrettyName)
			}
			if writeVersion {
				v := p.FullPrettyVersion(true)
				if format == "text" {
					v = strings.TrimLeft(v, "v")
				}
				d.Set("version", v)
				if len(v) > m.versionLength {
					m.versionLength = len(v)
				}
			}
			if writeReleaseDate {
				if rt, ok := releaseTime(p); ok {
					age := strings.ReplaceAll(relativeTime(rt), " ago", " old")
					if !strings.Contains(age, " old") {
						age = "from " + age
					}
					d.Set("release-age", age)
					if len(age) > m.releaseDateLength {
						m.releaseDateLength = len(age)
					}
					d.Set("release-date", rt.Format(atomFormat))
				} else {
					d.Set("release-age", "")
					d.Set("release-date", "")
				}
			}
			hasInstallable := latest != nil && latest.FullPrettyVersion(true) != p.FullPrettyVersion(true)
			shown := latest
			if withheld != nil && !hasInstallable {
				shown = withheld.p
			}
			if writeLatest && shown != nil {
				short := ""
				if withheld != nil && cooldown != nil {
					short = cooldown.FormatTimeUntilAvailable(withheld.releaseDate, time.Now(), true)
				}
				lv := shown.FullPrettyVersion(true)
				if format == "text" {
					lv = strings.TrimLeft(lv, "v")
				}
				d.Set("latest", lv)
				if rt, ok := releaseTime(shown); ok {
					d.Set("latest-release-date", rt.Format(atomFormat))
				} else {
					d.Set("latest-release-date", "")
				}
				if withheld != nil && !hasInstallable {
					d.Set("latest-status", "cooldown")
					d.Set("cooldown-available-in", withheld.availableIn)
					if format == "text" {
						lv += " (" + short + " left)"
					}
				} else {
					d.Set("latest-status", updateStatus(shown, p))
					if withheld != nil {
						cl := withheld.p.FullPrettyVersion(true)
						d.Set("cooldown-latest", cl)
						d.Set("cooldown-available-in", withheld.availableIn)
						if format == "text" {
							lv += " (" + strings.TrimLeft(cl, "v") + " in " + short + ")"
						}
					}
				}
				d.Set("latest", lv)
				if len(lv) > m.latestLength {
					m.latestLength = len(lv)
				}
			} else if writeLatest {
				d.Set("latest", "[none matched]")
				d.Set("latest-status", "up-to-date")
				if len("[none matched]") > m.latestLength {
					m.latestLength = len("[none matched]")
				}
			}
			if writeDescription {
				d.Set("description", nullable(p.Description()))
			}
			if writePath {
				path := sc.installPath(p)
				if path != "" {
					d.Set("path", firstLine(realpath(path)))
				} else {
					d.Set("path", nil)
				}
			}
			var abandoned phpjson.Value = false
			if latest != nil && latest.IsAbandoned() {
				repl := latest.ReplacementPackage()
				replacement := "No replacement was suggested"
				abandoned = true
				if repl != "" {
					replacement = "Use " + repl + " instead"
					abandoned = repl
				}
				d.Set("warning", "Package "+p.PrettyName+" is abandoned, you should avoid using it. "+replacement+".")
			}
			d.Set("abandoned", abandoned)
			list = append(list, d)
		}
		m.writeLatest, m.writeReleaseDate = writeLatest, writeReleaseDate
		viewData.Set(tv.typ, list)
		metas[tv.typ] = m
		typeOrder = append(typeOrder, tv.typ)
		if in.Bool("strict") && hasOutdated {
			exitCode = 1
			break
		}
	}

	if format == "json" {
		if viewData.Len() == 0 {
			io.Write(phpjson.MustEncode([]phpjson.Value{}, phpjson.ComposerFlags))
		} else {
			io.Write(phpjson.MustEncode(viewData, phpjson.ComposerFlags))
		}
		return exitCode, nil
	}
	anyData := false
	for _, typ := range typeOrder {
		v, _ := viewData.Get(typ)
		if len(v.([]phpjson.Value)) > 0 {
			anyData = true
		}
	}
	if in.Bool("latest") && anyData {
		if !io.Decorated {
			io.WriteError("Legend:")
			io.WriteError("! patch or minor release available - update recommended")
			io.WriteError("~ major release available - update possible")
			if cooldown != nil {
				io.WriteError("c newer release still in the cooldown period - not installable yet")
			}
			if !in.Bool("outdated") {
				io.WriteError("= up to date version")
			}
		} else {
			io.WriteError("<info>Color legend:</info>")
			io.WriteError("- <highlight>patch or minor</highlight> release available - update recommended")
			io.WriteError("- <comment>major</comment> release available - update possible")
			if cooldown != nil {
				io.WriteError("- release still in the cooldown period - not installable yet (shown without color)")
			}
			if !in.Bool("outdated") {
				io.WriteError("- <info>up to date</info> version")
			}
		}
	}
	width := baseTerminalWidth()
	for _, typ := range typeOrder {
		v, _ := viewData.Get(typ)
		list := v.([]phpjson.Value)
		m := metas[typ]
		versionFits := m.nameLength+m.versionLength+3 <= width
		latestFits := m.nameLength+m.versionLength+m.latestLength+3 <= width
		releaseDateFits := m.nameLength+m.versionLength+m.latestLength+m.releaseDateLength+3 <= width
		descriptionFits := m.nameLength+m.versionLength+m.latestLength+m.releaseDateLength+24 <= width
		latestLength := m.latestLength
		if latestFits && !io.Decorated {
			latestLength += 2
		}
		if showAllTypes {
			if typ == "available" {
				io.Write("<comment>" + typ + "</comment>:")
			} else {
				io.Write("<info>" + typ + "</info>:")
			}
		}
		opts := showPrintOpts{indent: indent, writeVersion: writeVersion && versionFits, writeDescription: writeDescription && descriptionFits, width: width,
			versionLength: m.versionLength, nameLength: m.nameLength, latestLength: latestLength, writeReleaseDate: m.writeReleaseDate && releaseDateFits, releaseDateLength: m.releaseDateLength}
		if m.writeLatest && !in.Bool("direct") {
			var direct, transitive []phpjson.Value
			for _, e := range list {
				if dd, ok := phpjson.AsMap(e).Get("direct-dependency"); ok && dd == true {
					direct = append(direct, e)
				} else {
					transitive = append(transitive, e)
				}
			}
			opts.writeLatest = latestFits
			io.WriteError("")
			io.WriteError("<info>Direct dependencies required in composer.json:</>")
			if len(direct) > 0 {
				sc.printPackages(direct, opts)
			} else {
				io.WriteError("Everything up to date")
			}
			io.WriteError("")
			io.WriteError("<info>Transitive dependencies not required in composer.json:</>")
			if len(transitive) > 0 {
				sc.printPackages(transitive, opts)
			} else {
				io.WriteError("Everything up to date")
			}
		} else {
			if m.writeLatest && len(list) == 0 {
				io.WriteError("All your direct dependencies are up to date")
			} else {
				opts.writeLatest = m.writeLatest && latestFits
				sc.printPackages(list, opts)
			}
		}
		if showAllTypes {
			io.Write("")
		}
	}
	return exitCode, nil
}

type showPrintOpts struct {
	indent                                                        string
	writeVersion, writeLatest, writeDescription, writeReleaseDate bool
	width, versionLength, nameLength, latestLength                int
	releaseDateLength                                             int
}

// printPackages is ShowCommand::printPackages.
func (sc *showCtx) printPackages(list []phpjson.Value, o showPrintOpts) {
	io := sc.ctx.IO
	padName := o.writeVersion || o.writeLatest || o.writeReleaseDate || o.writeDescription
	padVersion := o.writeLatest || o.writeReleaseDate || o.writeDescription
	padLatest := o.writeDescription || o.writeReleaseDate
	padReleaseDate := o.writeDescription
	for _, v := range list {
		p := phpjson.AsMap(v)
		name := p.Str("name")
		link := ""
		if s, ok := p.Get("source"); ok && s != nil {
			link, _ = s.(string)
		} else if h, ok := p.Get("homepage"); ok && h != nil {
			link, _ = h.(string)
		}
		if link != "" {
			pad := 0
			if padName {
				pad = o.nameLength - len(name)
			}
			io.WriteV(o.indent+"<href="+escapeTags(link)+">"+name+"</>"+strings.Repeat(" ", max0(pad)), false, 32)
		} else {
			w := 0
			if padName {
				w = o.nameLength
			}
			io.WriteV(o.indent+strPad(name, w), false, 32)
		}
		if ver, ok := p.Get("version"); ok && ver != nil && o.writeVersion {
			s, _ := ver.(string)
			w := 0
			if padVersion {
				w = o.versionLength
			}
			io.WriteV(" "+strPad(s, w), false, 32)
		}
		lat, hasLatest := p.Get("latest")
		st, hasStatus := p.Get("latest-status")
		if hasLatest && lat != nil && hasStatus && st != nil && o.writeLatest {
			latestVersion, _ := lat.(string)
			status, _ := st.(string)
			style := updateStatusToVersionStyle(status)
			if !io.Decorated {
				latestVersion = updateStatusSigns.Replace(status) + " " + latestVersion
			}
			w := 0
			if padLatest {
				w = o.latestLength
			}
			latestVersion = strPad(latestVersion, w)
			if style != "" {
				io.WriteV(" <"+style+">"+latestVersion+"</"+style+">", false, 32)
			} else {
				io.WriteV(" "+latestVersion, false, 32)
			}
			if age, ok := p.Get("release-age"); ok && age != nil && o.writeReleaseDate {
				s, _ := age.(string)
				w := 0
				if padReleaseDate {
					w = o.releaseDateLength
				}
				io.WriteV(" "+strPad(s, w), false, 32)
			}
		}
		if d, ok := p.Get("description"); ok && d != nil && o.writeDescription {
			desc := strtok(d.(string))
			remaining := o.width - o.nameLength - o.versionLength - o.releaseDateLength - 4
			if o.writeLatest {
				remaining -= o.latestLength
			}
			if remaining <= 0 {
				desc = ""
			} else if strWidth(desc) > remaining {
				desc = strimwidth(desc, remaining, "...")
			}
			io.WriteV(" "+desc, false, 32)
		}
		if path, ok := p.Get("path"); ok {
			if s, isStr := path.(string); isStr {
				io.WriteV(" "+s, false, 32)
			} else {
				io.WriteV(" null", false, 32)
			}
		}
		io.Write("")
		if w, ok := p.Get("warning"); ok && w != nil {
			io.Write("<warning>" + w.(string) + "</warning>")
		}
	}
}
