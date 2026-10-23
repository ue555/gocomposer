package cli

// Port of Composer\Command\CheckPlatformReqsCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"sort"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

func checkPlatformCommand() *Command {
	return &Command{Name: "check-platform-reqs", Run: runCheckPlatformReqs}
}

type platformResult struct {
	name, version string
	link          *pkg.Link
	status        string
	provider      string
}

func runCheckPlatformReqs(ctx *Context) (int, error) {
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	io, in := ctx.IO, ctx.In
	noDev := in.Bool("no-dev")
	nonDev := ""
	if noDev {
		nonDev = "non-dev "
	}
	requires := map[string]pkg.Links{}
	var removePackages []string
	var installedRepo repo.Repository
	if in.Bool("lock") {
		io.WriteError("<info>Checking " + nonDev + "platform requirements using the lock file</info>")
		lr, err := c.Locker.LockedRepository(!noDev)
		if err != nil {
			return 1, err
		}
		installedRepo = lr
	} else {
		if len(c.LocalRepo.MustPackages()) == 0 {
			io.WriteError("<warning>No vendor dir present, checking " + nonDev + "platform requirements from the lock file</warning>")
			lr, err := c.Locker.LockedRepository(!noDev)
			if err != nil {
				return 1, err
			}
			installedRepo = lr
		} else {
			if noDev {
				removePackages = c.LocalRepo.DevPackageNames
			}
			io.WriteError("<info>Checking " + nonDev + "platform requirements for packages in the vendor dir</info>")
			installedRepo = c.LocalRepo
		}
	}
	if !noDev {
		for _, l := range c.Root().DevRequires() {
			requires[l.Target] = pkg.Links{l}
		}
	}
	inst := repo.NewInstalled(installedRepo, repo.NewRootPackageRepository(c.Root()))
	for _, p := range inst.Packages() {
		if inList(p.Name, removePackages) {
			continue
		}
		for _, l := range p.Requires() {
			requires[l.Target] = append(requires[l.Target], l)
		}
	}
	names := make([]string, 0, len(requires))
	for n := range requires {
		names = append(names, n)
	}
	sort.Strings(names)
	platformRepo, err := c.CreatePlatformRepo(nil)
	if err != nil {
		return 1, err
	}
	inst.AddRepository(platformRepo)

	var results []platformResult
	exitCode := 0
	for _, require := range names {
		links := requires[require]
		if !repo.IsPlatformPackage(require) {
			continue
		}
		candidates := inst.FindPackagesWithReplacersAndProviders(require, nil)
		if len(candidates) == 0 {
			results = append(results, platformResult{require, "n/a", links[0], "<error>missing</error>", ""})
			if exitCode < 2 {
				exitCode = 2
			}
			continue
		}
		var reqResults []platformResult
		matched := false
	candidates:
		for _, cand := range candidates {
			var candConstraint semver.Constraint
			if cand.Name == require {
				s := semver.NewConstraint("=", cand.Version())
				s.SetPrettyString(cand.PrettyVersion())
				candConstraint = s
			} else {
				for _, l := range append(append(pkg.Links(nil), cand.Provides()...), cand.Replaces()...) {
					if l.Target == require {
						candConstraint = l.Constraint
						break
					}
				}
			}
			if candConstraint == nil {
				continue
			}
			name, provider := require, "<comment>provided by "+cand.PrettyName+"</comment>"
			if cand.Name == require {
				name, provider = cand.PrettyName, ""
			}
			for _, l := range links {
				if !l.Constraint.Matches(candConstraint) {
					reqResults = append(reqResults, platformResult{name, candConstraint.PrettyString(), l, "<error>failed</error>", provider})
					continue candidates
				}
			}
			results = append(results, platformResult{name, candConstraint.PrettyString(), nil, "<info>success</info>", provider})
			matched = true
			break
		}
		if matched {
			continue
		}
		results = append(results, reqResults...)
		if exitCode < 1 {
			exitCode = 1
		}
	}
	format, _ := in.String("format")
	if format == "json" {
		rows := []phpjson.Value{}
		for _, r := range results {
			m := phpjson.NewMap()
			m.Set("name", r.name)
			m.Set("version", r.version)
			m.Set("status", cio.StripTags(r.status))
			if r.link != nil {
				fr := phpjson.NewMap()
				fr.Set("source", r.link.Source)
				fr.Set("type", r.link.Description)
				fr.Set("target", r.link.Target)
				fr.Set("constraint", r.link.PrettyConstraint)
				m.Set("failed_requirement", fr)
			} else {
				m.Set("failed_requirement", nil)
			}
			if r.provider == "" {
				m.Set("provider", nil)
			} else {
				m.Set("provider", cio.StripTags(r.provider))
			}
			rows = append(rows, m)
		}
		s, err := phpjson.Encode(rows, phpjson.ComposerFlags)
		if err != nil {
			return 1, err
		}
		io.Write(s)
	} else {
		t := cio.NewTable("compact")
		for _, r := range results {
			linkStr, pretty := "", ""
			if r.link != nil {
				linkStr = r.link.String()
				pretty = r.link.Source + " " + r.link.Description + " " + r.link.Target + " (" + r.link.PrettyConstraint + ")"
			}
			t.Rows = append(t.Rows, []string{r.name, r.version, linkStr, pretty, strings.TrimRight(r.status+" "+r.provider, " \t\n\r\x00\x0B")})
		}
		t.Render(io)
	}
	return exitCode, nil
}

func inList(s string, list []string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}
