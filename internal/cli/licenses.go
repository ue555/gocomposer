package cli

// Ports of Composer\Command\LicensesCommand and SuggestsCommand
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
)

func licensesCommand() *Command {
	return &Command{Name: "licenses", Run: runLicenses}
}

func suggestsCommand() *Command {
	return &Command{Name: "suggests", Run: runSuggests}
}

func runLicenses(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	root := c.Package
	var packages []*pkg.Package
	if in.Bool("locked") {
		if locked, _ := c.Locker.IsLocked(); !locked {
			return 1, &composer.UserError{Msg: "Valid composer.json and composer.lock files are required to run this command with --locked", Code: 1}
		}
		lr, err := c.Locker.LockedRepository(!in.Bool("no-dev"))
		if err != nil {
			return 1, err
		}
		packages = lr.MustPackages()
	} else if in.Bool("no-dev") {
		packages = filterRequiredPackages(c.LocalRepo.MustPackages(), root, false, nil)
	} else {
		packages = c.LocalRepo.MustPackages()
	}
	packages = append([]*pkg.Package(nil), packages...)
	sort.SliceStable(packages, func(i, j int) bool { return packages[i].Name < packages[j].Name })
	format, _ := in.String("format")
	switch format {
	case "text":
		lic := strings.Join(root.License(), ", ")
		if lic == "" {
			lic = "none"
		}
		io.Write("Name: <comment>" + root.PrettyName + "</comment>")
		io.Write("Version: <comment>" + root.FullPrettyVersion(true) + "</comment>")
		io.Write("Licenses: <comment>" + lic + "</comment>")
		io.Write("Dependencies:")
		io.Write("")
		t := cio.NewTable("compact")
		t.Headers = []string{"Name", "Version", "Licenses"}
		for _, p := range packages {
			l := strings.Join(p.License(), ", ")
			if l == "" {
				l = "none"
			}
			t.Rows = append(t.Rows, []string{nameWithLink(p), p.FullPrettyVersion(true), l})
		}
		t.Render(io)
	case "json":
		deps := phpjson.NewMap()
		for _, p := range packages {
			d := phpjson.NewMap()
			d.Set("version", p.FullPrettyVersion(true))
			d.Set("license", strList(p.License()))
			deps.Set(p.PrettyName, d)
		}
		out := phpjson.NewMap()
		out.Set("name", root.PrettyName)
		out.Set("version", root.FullPrettyVersion(true))
		out.Set("license", strList(root.License()))
		if deps.Len() == 0 {
			out.Set("dependencies", []phpjson.Value{})
		} else {
			out.Set("dependencies", deps)
		}
		io.Write(phpjson.MustEncode(out, phpjson.ComposerFlags))
	case "summary":
		counts := map[string]int{}
		var order []string
		for _, p := range packages {
			ls := p.License()
			if len(ls) == 0 {
				ls = []string{"none"}
			}
			for _, l := range ls {
				if _, ok := counts[l]; !ok {
					order = append(order, l)
				}
				counts[l]++
			}
		}
		// arsort($usedLicenses, SORT_NUMERIC) is stable since PHP 8
		sort.SliceStable(order, func(i, j int) bool { return counts[order[i]] > counts[order[j]] })
		st := cio.TableStyles["symfony-style-guide"]
		st.CellHeaderFormat = "<info>%s</info>"
		t := &cio.Table{Style: st, Headers: []string{"License", "Number of dependencies"}}
		for _, l := range order {
			t.Rows = append(t.Rows, []string{l, strconv.Itoa(counts[l])})
		}
		t.Render(io)
		io.Write("")
	default:
		return 1, &composer.UserError{Msg: fmt.Sprintf(`Unsupported format "%s".  See help for supported formats.`, format), Code: 1}
	}
	return 0, nil
}

func runSuggests(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	repos := []repo.Repository{repo.NewRootPackageRepository(c.Package.Clone())}
	locked, _ := c.Locker.IsLocked()
	if locked {
		pr, err := c.CreatePlatformRepo(c.Locker.PlatformOverrides())
		if err != nil {
			return 1, err
		}
		lr, err := c.Locker.LockedRepository(!in.Bool("no-dev"))
		if err != nil {
			return 1, err
		}
		repos = append(repos, pr, lr)
	} else {
		pr, err := c.CreatePlatformRepo(c.Config.Get("platform"))
		if err != nil {
			return 1, err
		}
		repos = append(repos, pr, c.LocalRepo)
	}
	installed := repo.NewInstalled(repos...)
	filter := in.Args("packages")
	var normalized []string
	for _, f := range filter {
		normalized = append(normalized, strings.ToLower(f))
	}
	validation := installed
	if locked && in.Bool("no-dev") {
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
		validation = repo.NewInstalled(lr)
	}
	var missing, excluded []string
	for i, p := range normalized {
		if len(installed.FindPackages(p, nil)) > 0 {
			continue
		}
		if len(validation.FindPackages(p, nil)) > 0 {
			excluded = append(excluded, filter[i])
		} else {
			missing = append(missing, filter[i])
		}
	}
	var errs []string
	if len(missing) == 1 {
		errs = append(errs, `Package "`+missing[0]+`" is not installed.`)
	} else if len(missing) > 1 {
		errs = append(errs, `Packages "`+strings.Join(missing, `", "`)+`" are not installed.`)
	}
	if len(excluded) == 1 {
		errs = append(errs, `Package "`+excluded[0]+`" is excluded by --no-dev.`)
	} else if len(excluded) > 1 {
		errs = append(errs, `Packages "`+strings.Join(excluded, `", "`)+`" are excluded by --no-dev.`)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			io.WriteError("<error>" + e + "</error>")
		}
		return 1, nil
	}
	reporter := composer.NewSuggestedPackagesReporter(io)
	packages := append(installed.Packages(), c.Package)
	for _, p := range packages {
		if len(normalized) > 0 && !inList(p.Name, normalized) {
			continue
		}
		reporter.AddSuggestionsFromPackage(p)
	}
	mode := composer.SuggestModeByPackage
	if in.Bool("by-suggestion") {
		mode = composer.SuggestModeBySuggestion
	}
	if in.Bool("by-package") {
		mode |= composer.SuggestModeByPackage
	}
	if in.Bool("list") {
		mode = composer.SuggestModeList
	}
	var only *pkg.Package
	if len(normalized) == 0 && !in.Bool("all") {
		only = c.Package
	}
	reporter.Output(mode, installed, only)
	return 0, nil
}
