package cli

// Port of Composer\Command\SearchCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// Without a composer.json the default repositories (packagist.org) are
// searched (createComposerInstance with an empty configuration).

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/repo"
)

func init() { registerCommand(searchCommand) }

func searchCommand() *Command {
	return &Command{Name: "search", Run: runSearch}
}

func runSearch(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	format, _ := in.String("format")
	if format != "text" && format != "json" {
		io.WriteError(fmt.Sprintf(`Unsupported format "%s". See help for supported formats.`, format))
		return 1, nil
	}
	c := ctx.cache
	if c == nil {
		// createComposerInstance($input, $io, []): the default repositories
		nc, err := defaultReposComposer(ctx)
		if err != nil {
			return 1, err
		}
		c = nc
	}
	localPkgs, err := c.LocalRepo.Packages()
	if err != nil {
		return 1, err
	}
	// CompositeRepository flattens the installed composite: local repository,
	// platform repository, then the configured repositories
	plural := ""
	if len(localPkgs) > 1 {
		plural = "s"
	}
	type searchRepo struct {
		name string
		r    repo.Repository
	}
	repos := []searchRepo{{fmt.Sprintf("installed array repo (defining %d package%s)", len(localPkgs), plural), repo.NewArray(localPkgs...)}}
	if pr, err := c.CreatePlatformRepo(nil); err == nil {
		repos = append(repos, searchRepo{pr.RepoName(), pr})
	}
	for _, r := range c.Repositories {
		repos = append(repos, searchRepo{r.RepoName(), r})
	}

	mode := repo.SearchFulltext
	if in.Bool("only-name") {
		if in.Bool("only-vendor") {
			return 1, &composer.UserError{Msg: "--only-name and --only-vendor cannot be used together", Code: 1}
		}
		mode = repo.SearchName
	} else if in.Bool("only-vendor") {
		mode = repo.SearchVendor
	}
	var typ *string
	if t, ok := in.String("type"); ok && in.present["type"] {
		typ = &t
	}
	query := strings.Join(in.Args("tokens"), " ")
	if mode != repo.SearchFulltext {
		query = pregQuote(query)
	}
	var results []repo.SearchEntry
	for _, sr := range repos {
		res, err := repo.SearchWithMode(sr.r, query, mode, typ)
		if err != nil {
			return 1, err
		}
		io.WriteErrorV(fmt.Sprintf("Searched %s, found <info>%d</info> result(s)", sr.name, len(res)), true, cio.VeryVerbose)
		results = append(results, res...)
	}
	if len(results) > 0 && format == "text" {
		width := cio.TerminalWidth()
		if width < 80 {
			width = 80
		}
		nameLength := 0
		for _, r := range results {
			if len(r.Name) > nameLength {
				nameLength = len(r.Name)
			}
		}
		nameLength++
		for _, r := range results {
			desc := ""
			if r.Description != nil {
				desc = *r.Description
			}
			warning := ""
			if phpjson.Truthy(r.Abandoned) {
				warning = "<warning>! Abandoned !</warning> "
			}
			remaining := width - nameLength - len(warning) - 2
			if len(desc) > remaining {
				desc = phpSubstr(desc, remaining-3) + "..."
			}
			if r.URL != "" {
				io.Write("<href=" + escapeFormatter(r.URL) + ">" + r.Name + "</>" + strings.Repeat(" ", nameLength-len(r.Name)) + warning + desc)
			} else {
				io.Write(r.Name + strings.Repeat(" ", nameLength-len(r.Name)) + warning + desc)
			}
		}
	} else if format == "json" {
		list := []phpjson.Value{}
		for _, r := range results {
			if r.Raw != nil {
				list = append(list, r.Raw)
				continue
			}
			m := phpjson.NewMap()
			m.Set("name", r.Name)
			if r.Description != nil {
				m.Set("description", *r.Description)
			} else {
				m.Set("description", nil)
			}
			if r.Abandoned != nil {
				m.Set("abandoned", r.Abandoned)
			}
			list = append(list, m)
		}
		s, err := phpjson.Encode(list, phpjson.ComposerFlags)
		if err != nil {
			return 1, err
		}
		io.Write(s)
	}
	return 0, nil
}

var rePregQuote = regexp.MustCompile(`[.\\+*?\[^\]$(){}=!<>|:#-]`)

// pregQuote is preg_quote($s) (no delimiter).
func pregQuote(s string) string {
	s = rePregQuote.ReplaceAllString(s, `\$0`)
	return strings.ReplaceAll(s, "\x00", `\000`)
}

// phpSubstr is substr($s, 0, $n) on bytes (negative n counts from the end).
func phpSubstr(s string, n int) string {
	if n < 0 {
		n = len(s) + n
		if n < 0 {
			n = 0
		}
	}
	if n > len(s) {
		n = len(s)
	}
	return s[:n]
}

// escapeFormatter is OutputFormatter::escape.
func escapeFormatter(s string) string {
	return regexp.MustCompile(`([^\\]|^)([<>])`).ReplaceAllString(s, `$1\$2`)
}

// defaultReposComposer is a Composer instance created from an empty
// configuration (Factory::create($io, [])), whose repositories are the
// default ones of the (global) configuration: packagist.org unless disabled.
func defaultReposComposer(ctx *Context) (*composer.Composer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	c, err := composer.NewProjectComposer(ctx.IO, phpjson.NewMap(), cwd, composer.Options{DisablePlugins: true, DisableScripts: ctx.In.Bool("no-scripts")})
	if err != nil {
		return nil, err
	}
	repos, err := c.DefaultRepositories()
	if err != nil {
		return nil, err
	}
	c.Repositories = repos
	return c, nil
}
