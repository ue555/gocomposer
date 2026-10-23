package cli

// Port of Composer\Command\BaseDependencyCommand, DependsCommand and
// ProhibitsCommand (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

func dependsCommand(prohibits bool) *Command {
	if prohibits {
		return &Command{Name: "prohibits", Run: func(ctx *Context) (int, error) { return runDepends(ctx, true) }}
	}
	return &Command{Name: "depends", Run: func(ctx *Context) (int, error) { return runDepends(ctx, false) }}
}

func runDepends(ctx *Context, inverted bool) (int, error) {
	io, in := ctx.IO, ctx.In
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	repos := []repo.Repository{repo.NewRootPackageRepository(c.Package.Clone())}
	if in.Bool("locked") {
		if locked, _ := c.Locker.IsLocked(); !locked {
			return 1, &composer.UserError{Msg: "A valid composer.lock file is required to run this command with --locked", Code: 1}
		}
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return 1, err
		}
		repos = append(repos, lr)
		pr, err := c.CreatePlatformRepo(c.Locker.PlatformOverrides())
		if err != nil {
			return 1, err
		}
		repos = append(repos, pr)
	} else {
		root := c.Package
		if len(c.LocalRepo.MustPackages()) == 0 && (len(root.Requires()) > 0 || len(root.DevRequires()) > 0) {
			io.Write("<warning>No dependencies installed. Try running composer install or update, or use --locked.</warning>")
			return 1, nil
		}
		repos = append(repos, c.LocalRepo)
		overrides := c.Config.Get("platform")
		if !phpjson.Truthy(overrides) {
			overrides = nil
		}
		pr, err := c.CreatePlatformRepo(overrides)
		if err != nil {
			return 1, err
		}
		repos = append(repos, pr)
	}
	installed := repo.NewInstalled(repos...)
	needle := in.Arg("package")
	textConstraint := "*"
	if ctx.Cmd.Name == "prohibits" {
		textConstraint = in.Arg("version")
	}
	packages := installed.FindPackagesWithReplacersAndProviders(needle, nil)
	if len(packages) == 0 {
		return 1, &composer.UserError{Msg: fmt.Sprintf(`Could not find package "%s" in your project`, needle), Code: 1}
	}
	tc, err := semver.ParseConstraints(textConstraint)
	if err != nil {
		return 1, err
	}
	matched := installed.FindPackage(needle, tc)
	if matched == nil {
		var match *pkg.Package
		for _, r := range c.Repositories {
			if m := repo.FindPackage(r, needle, tc); m != nil {
				match = m
				break
			}
		}
		if match != nil {
			arr := repo.NewInstalledArray()
			arr.AddPackage(match.Clone())
			installed.AddRepository(arr)
		} else if repo.IsPlatformPackage(needle) {
			// Bound::zero() returns a fresh instance, so the PHP `!==`
			// check always passes
			lb := tc.LowerBound().Version
			tmp := pkg.NewPackage(needle, lb, lb)
			arr := repo.NewInstalledArray()
			arr.AddPackage(tmp)
			installed.AddRepository(arr)
		} else {
			io.WriteError(`<error>Package "` + needle + `" could not be found with constraint "` + textConstraint + `", results below will most likely be incomplete.</error>`)
		}
	} else if repo.IsPlatformPackage(needle) {
		extra := ""
		if e := matched.Extra(); e != nil {
			if v, ok := e.Get("config.platform"); ok && v == true {
				extra = " (version provided by config.platform)"
			}
		}
		io.WriteError(`<info>Package "` + needle + " " + textConstraint + `" found in version "` + matched.PrettyVersion() + `"` + extra + ".</info>")
	} else if inverted {
		io.Write(`<comment>Package "` + needle + `" ` + matched.PrettyVersion() + " is already installed! To find out why, run `composer why " + needle + "`</comment>")
		return 0, nil
	}
	needles := []string{needle}
	if inverted {
		for _, p := range packages {
			for _, l := range p.Replaces() {
				needles = append(needles, l.Target)
			}
		}
	}
	var constraint semver.Constraint
	if textConstraint != "*" {
		constraint = tc
	}
	renderTree := in.Bool("tree")
	recursive := renderTree || in.Bool("recursive")
	ret := 0
	if inverted {
		ret = 1
	}
	results := installed.GetDependents(needles, constraint, inverted, recursive, nil)
	if len(results) == 0 {
		extra := ""
		if constraint != nil {
			not := ""
			if inverted {
				not = "not "
			}
			extra = fmt.Sprintf(" in versions %smatching %s", not, textConstraint)
		}
		io.WriteError(fmt.Sprintf(`<info>There is no installed package depending on "%s"%s</info>`, needle, extra))
		if inverted {
			ret = 0
		} else {
			ret = 1
		}
	} else if renderTree {
		for _, col := range treeColors {
			io.SetStyle(col, col)
		}
		root := packages[0]
		io.Write(fmt.Sprintf("<info>%s</info> %s %s", root.PrettyName, root.PrettyVersion(), root.Description()))
		printDependentsTree(io, results, "", 1)
	} else {
		printDependentsTable(io, results)
	}
	if inverted && !repo.IsPlatformPackage(needle) {
		cmd := "update"
		for _, r := range c.Package.Requires() {
			if r.Target == needle {
				cmd = "require"
				break
			}
		}
		for _, r := range c.Package.DevRequires() {
			if r.Target == needle {
				cmd = "require --dev"
				break
			}
		}
		io.WriteError("Not finding what you were looking for? Try calling `composer " + cmd + ` "` + needle + ":" + textConstraint + `" --dry-run` + "` to get another view on the problem.")
	}
	return ret, nil
}

// viewSourceOrHomepageURL is PackageInfo::getViewSourceOrHomepageUrl.
func viewSourceOrHomepageURL(p *pkg.Package) string {
	u := viewSourceURL(p)
	if u == "" {
		u = p.Homepage()
	}
	return u
}

func nameWithLink(p *pkg.Package) string {
	if u := viewSourceOrHomepageURL(p); u != "" {
		return "<href=" + escapeTags(u) + ">" + p.PrettyName + "</>"
	}
	return p.PrettyName
}

func printDependentsTable(io *cio.IO, results []repo.Dependent) {
	var table [][]string
	doubles := map[string]bool{}
	for len(results) > 0 {
		var queue []repo.Dependent
		var rows [][]string
		for _, r := range results {
			unique := r.Link.String()
			if doubles[unique] {
				continue
			}
			doubles[unique] = true
			version := r.Package.PrettyVersion()
			if version == pkg.DefaultPrettyVersion {
				version = "-"
			}
			rows = append(rows, []string{nameWithLink(r.Package), version, r.Link.Description, fmt.Sprintf("%s (%s)", r.Link.Target, r.Link.PrettyConstraint)})
			if !r.Cut {
				queue = append(queue, r.Children...)
			}
		}
		results = queue
		table = append(rows, table...)
	}
	t := cio.NewTable("compact")
	t.Rows = table
	t.Render(io)
}

func printDependentsTree(io *cio.IO, results []repo.Dependent, prefix string, level int) {
	for idx, r := range results {
		color := treeColors[level%len(treeColors)]
		prevColor := treeColors[(level-1)%len(treeColors)]
		isLast := idx == len(results)-1
		versionText := r.Package.PrettyVersion()
		if versionText == pkg.DefaultPrettyVersion {
			versionText = ""
		}
		packageText := strings.TrimRight(fmt.Sprintf("<%s>%s</%s> %s", color, nameWithLink(r.Package), color, versionText), " \t\n\r\x00\x0B")
		linkText := fmt.Sprintf("%s <%s>%s</%s> %s", r.Link.Description, prevColor, r.Link.Target, prevColor, r.Link.PrettyConstraint)
		circular := ""
		if r.Cut {
			circular = "(circular dependency aborted here)"
		}
		bar := "├──"
		if isLast {
			bar = "└──"
		}
		line := strings.TrimRight(fmt.Sprintf("%s%s%s (%s) %s", prefix, bar, packageText, linkText, circular), " \t\n\r\x00\x0B")
		if !io.Decorated {
			line = strings.NewReplacer("└", "`-", "├", "|-", "──", "-", "│", "|").Replace(line)
		}
		io.Write(line)
		if !r.Cut {
			next := "│  "
			if isLast {
				next = "   "
			}
			printDependentsTree(io, r.Children, prefix+next, level+1)
		}
	}
}
