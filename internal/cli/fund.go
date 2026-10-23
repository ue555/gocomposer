package cli

// Port of Composer\Command\FundCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
)

func init() { registerCommand(fundCommand) }

func fundCommand() *Command {
	return &Command{Name: "fund", Run: runFund}
}

func runFund(ctx *Context) (int, error) {
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	io := ctx.IO
	local, err := c.LocalRepo.Packages()
	if err != nil {
		return 1, err
	}
	toLoad := map[string]bool{}
	var names []repo.NameConstraint
	for _, p := range local {
		if p.AliasOf != nil {
			continue
		}
		if !toLoad[p.Name] {
			names = append(names, repo.NameConstraint{Name: p.Name})
		}
		toLoad[p.Name] = true
	}
	fundings := phpjson.NewMap()
	acceptable := map[string]int{"dev": pkg.Stabilities["dev"]}
	var loaded []*pkg.Package
	for _, r := range c.Repositories {
		res, err := r.LoadPackages(names, acceptable, map[string]int{}, nil)
		if err != nil {
			return 1, err
		}
		loaded = append(loaded, res.Packages...)
	}
	for _, p := range loaded {
		if p.AliasOf == nil && p.IsDefaultBranch() && len(p.Funding()) > 0 && toLoad[p.Name] {
			insertFunding(fundings, p)
			delete(toLoad, p.Name)
		}
	}
	for _, p := range local {
		if p.AliasOf != nil || !toLoad[p.Name] {
			continue
		}
		if len(p.Funding()) > 0 {
			insertFunding(fundings, p)
		}
	}
	fundings.KSort()

	format, _ := ctx.In.String("format")
	if format != "text" && format != "json" {
		io.WriteError(fmt.Sprintf(`Unsupported format "%s". See help for supported formats.`, format))
		return 1, nil
	}
	switch {
	case fundings.Len() > 0 && format == "text":
		prev := ""
		io.Write("The following packages were found in your dependencies which publish funding information:")
		for _, vendor := range fundings.Keys() {
			links := fundings.Map(vendor)
			io.Write("")
			io.Write("<comment>" + vendor + "</comment>")
			for _, u := range links.Keys() {
				var pkgs []string
				v, _ := links.Get(u)
				phpjson.Each(v, func(_ string, e phpjson.Value) { s, _ := phpjson.ScalarString(e); pkgs = append(pkgs, s) })
				line := "  <info>" + strings.Join(pkgs, ", ") + "</info>"
				if prev != line {
					io.Write(line)
					prev = line
				}
				io.Write("    <href=" + escapeFormatter(u) + ">" + u + "</>")
			}
		}
		io.Write("")
		io.Write("Please consider following these links and sponsoring the work of package authors!")
		io.Write("Thank you!")
	case format == "json":
		var v phpjson.Value = fundings
		if fundings.Len() == 0 {
			v = []phpjson.Value{}
		}
		s, err := phpjson.Encode(v, phpjson.ComposerFlags)
		if err != nil {
			return 1, err
		}
		io.Write(s)
	default:
		io.Write("No funding links were found in your package dependencies. This doesn't mean they don't need your support!")
	}
	return 0, nil
}

var reGithubSponsor = regexp.MustCompile(`^https://github.com/([^/]+)$`)

// insertFunding is FundCommand::insertFundingData.
func insertFunding(fundings *phpjson.Map, p *pkg.Package) {
	for _, opt := range p.Funding() {
		vendor, name, _ := strings.Cut(p.PrettyName, "/")
		m := phpjson.ToMap(opt)
		if m == nil {
			continue
		}
		uv, _ := m.Get("url")
		if !phpjson.Truthy(uv) {
			continue
		}
		u, _ := phpjson.ScalarString(uv)
		tv, _ := m.Get("type")
		if phpjson.Truthy(tv) && tv == "github" {
			if mm := reGithubSponsor.FindStringSubmatch(u); mm != nil {
				u = "https://github.com/sponsors/" + mm[1]
			}
		}
		links := fundings.Map(vendor)
		if links == nil {
			links = phpjson.NewMap()
			fundings.Set(vendor, links)
		}
		var list []phpjson.Value
		if v, ok := links.Get(u); ok {
			list, _ = v.([]phpjson.Value)
		}
		links.Set(u, append(list, name))
	}
}
