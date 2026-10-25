package cli

// Port of Composer\Command\HomeCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// Without a composer.json the packages are looked up in the default
// repositories (RepositoryFactory::defaultReposWithDefaultManager).

import (
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
)

func init() { registerCommand(browseCommand) }

func browseCommand() *Command {
	return &Command{Name: "browse", Run: runBrowse}
}

func runBrowse(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	c := ctx.cache
	ret := 0
	packages := in.Args("packages")
	if len(packages) == 0 {
		io.WriteError("No package specified, opening homepage for the root package")
		rc, err := requireComposer(ctx)
		if err != nil {
			return exitFor(err)
		}
		c = rc
		packages = []string{c.Root().Name}
	}
	var finders []func(name string) ([]*pkg.Package, error)
	if c == nil {
		// HomeCommand::initializeRepos without a composer.json
		dc, err := defaultReposComposer(ctx)
		if err != nil {
			return 1, err
		}
		c = dc
	} else {
		finders = append(finders,
			func(name string) ([]*pkg.Package, error) {
				return repo.FindInList([]*pkg.Package{c.Package}, name, nil), nil
			},
			func(name string) ([]*pkg.Package, error) {
				ps, err := c.LocalRepo.Packages()
				if err != nil {
					return nil, err
				}
				return repo.FindInList(ps, name, nil), nil
			},
		)
	}
	for _, r := range c.Repositories {
		r := r
		finders = append(finders, func(name string) ([]*pkg.Package, error) { return r.FindPackages(name, nil) })
	}
	homepage, show := in.Bool("homepage"), in.Bool("show")
	for _, name := range packages {
		handled, exists := false, false
	outer:
		for _, find := range finders {
			found, err := find(name)
			if err != nil {
				return 1, err
			}
			for _, p := range found {
				exists = true
				if browsePackage(ctx, p, homepage, show) {
					handled = true
					break outer
				}
			}
		}
		if !exists {
			ret = 1
			io.WriteError("<warning>Package " + name + " not found</warning>")
		}
		if !handled {
			ret = 1
			what := "Invalid or missing repository URL"
			if homepage {
				what = "Invalid or missing homepage"
			}
			io.WriteError("<warning>" + what + " for " + name + "</warning>")
		}
	}
	return ret, nil
}

// browsePackage is HomeCommand::handlePackage.
func browsePackage(ctx *Context, p *pkg.Package, showHomepage, showOnly bool) bool {
	u := ""
	if s := p.Support(); s != nil {
		if v, ok := s.Get("source"); ok && v != nil {
			u, _ = phpjson.ScalarString(v)
		} else {
			u = p.SourceURL()
		}
	} else {
		u = p.SourceURL()
	}
	if u == "" || u == "0" || showHomepage {
		u = p.Homepage()
	}
	if u == "" || u == "0" || !filterValidateURL(u) {
		return false
	}
	if showOnly {
		ctx.IO.Write("<info>" + u + "</info>")
	} else {
		openBrowser(ctx, u)
	}
	return true
}

// openBrowser is HomeCommand::openBrowser.
func openBrowser(ctx *Context, u string) {
	if runtime.GOOS == "windows" {
		_ = exec.Command("cmd", "/c", "start", `"web"`, "explorer", u).Run()
		return
	}
	which := func(bin string) bool { return exec.Command("which", bin).Run() == nil }
	switch {
	case which("xdg-open"):
		_ = exec.Command("xdg-open", u).Run()
	case which("open"):
		_ = exec.Command("open", u).Run()
	default:
		ctx.IO.WriteError("No suitable browser opening command found, open yourself: " + u)
	}
}

var reURLScheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*$`)

// filterValidateURL approximates filter_var($url, FILTER_VALIDATE_URL):
// an ASCII URL with a scheme, and a host for http(s)/mailto-less schemes
// that use an authority.
func filterValidateURL(s string) bool {
	for _, r := range s {
		if r <= ' ' || r >= 0x7f {
			return false
		}
	}
	scheme, rest, ok := strings.Cut(s, ":")
	if !ok || !reURLScheme.MatchString(scheme) {
		return false
	}
	pu, err := url.Parse(s)
	if err != nil {
		return false
	}
	ls := strings.ToLower(scheme)
	if (ls == "http" || ls == "https") && pu.Host == "" {
		return false
	}
	if strings.HasPrefix(rest, "//") && pu.Host == "" && ls != "file" {
		return false
	}
	return true
}
