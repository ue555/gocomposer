package cli

// Port of UpdateCommand::getPackagesInteractively and createVersionSelector
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano): `update --interactive` lists the
// packages with a newer version and lets the user pick them (a multiselect
// ChoiceQuestion with autocompletion).

import (
	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/solver"
)

func updatePackagesInteractively(ctx *Context, c *composer.Composer, packages []string) ([]string, error) {
	io, in := ctx.IO, ctx.In
	if !io.Interactive {
		return nil, &composer.UserError{Msg: "--interactive cannot be used in non-interactive terminals.", Code: 1}
	}
	filter := platformFilter(in)
	root := c.Root()
	requires := map[string]string{}
	var reqOrder []string
	for _, l := range append(append(pkg.Links(nil), root.Requires()...), root.DevRequires()...) {
		if _, ok := requires[l.Target]; !ok {
			reqOrder = append(reqOrder, l.Target)
		}
		requires[l.Target] = l.PrettyConstraint
	}
	var nameFilter interface{ MatchString(string) bool }
	if len(packages) > 0 {
		nameFilter = solver.PackageNamesToRegexp(packages)
	}
	io.WriteError("<info>Loading packages that can be updated...</info>")
	var installed []*pkg.Package
	if locked, _ := c.Locker.IsLocked(); locked {
		lr, err := c.Locker.LockedRepository(true)
		if err != nil {
			return nil, err
		}
		installed = lr.MustPackages()
	} else {
		ps, err := c.LocalRepo.Packages()
		if err != nil {
			return nil, err
		}
		installed = ps
	}
	// createVersionSelector: every repository but the platform one
	rs := solver.NewRepositorySet("stable", nil, nil, nil, nil, nil)
	for _, r := range c.Repositories {
		if r.Kind() != repo.KindPlatform {
			rs.AddRepository(r)
		}
	}
	vs := composer.NewVersionSelector(rs, nil)
	var choices []cio.Choice
	for _, p := range installed {
		if nameFilter != nil && !nameFilter.MatchString(p.Name) {
			continue
		}
		constraint := requires[p.Name]
		stability := root.MinimumStability
		if f, ok := root.StabilityFlags[p.Name]; ok {
			stability = ""
			for name, prio := range pkg.Stabilities {
				if prio == f {
					stability = name
				}
			}
		}
		latest, err := vs.FindBestCandidate(p.Name, constraint, stability, filter, 0, nil, true)
		if err != nil {
			return nil, err
		}
		if latest != nil && (p.Version() != latest.Version() || latest.IsDev()) {
			choices = setChoice(choices, p.Name, "<comment>"+p.PrettyVersion()+"</comment> => <comment>"+latest.PrettyVersion()+"</comment>")
		}
	}
	if len(installed) == 0 {
		for _, req := range reqOrder {
			if repo.IsPlatformPackage(req) {
				continue
			}
			choices = setChoice(choices, req, "")
		}
	}
	if len(choices) == 0 {
		return nil, &composer.UserError{Msg: "Could not find any package with new versions available", Code: 1}
	}
	res, err := io.Select("Select packages: (Select more than one value separated by comma) ", choices, true, false, 1, `No package named "%s" is installed.`, true)
	if err != nil {
		return nil, err
	}
	selected, _ := res.([]string)
	t := cio.NewTable("default")
	t.Headers = []string{"Selected packages"}
	for _, s := range selected {
		t.Rows = append(t.Rows, []string{s})
	}
	t.Render(io)
	plural := "s"
	if len(selected) == 1 {
		plural = ""
	}
	ok, err := io.AskConfirmation("Would you like to continue and update the above package"+plural+" [<comment>yes</comment>]? ", true)
	if err != nil {
		return nil, err
	}
	if ok {
		return selected, nil
	}
	return nil, &composer.UserError{Msg: "Installation aborted.", Code: 1}
}

// setChoice sets $choices[$key] = $value keeping the PHP array order.
func setChoice(choices []cio.Choice, key, value string) []cio.Choice {
	for i, c := range choices {
		if c.Key == key {
			choices[i].Value = value
			return choices
		}
	}
	return append(choices, cio.Choice{Key: key, Value: value})
}
