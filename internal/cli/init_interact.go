package cli

// Port of the interactive parts of Composer\Command\InitCommand (interact(),
// the .gitignore/.gitattributes questions and updateDependencies) and of
// Symfony's FormatterHelper::formatBlock (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
// Jordi Boggiano; symfony/console v5.4.54, MIT, (c) Fabien Potencier).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/spdx"
)

var (
	reInitAutoloadPath = regexp.MustCompile(`^[^/][A-Za-z0-9\-_/]+/$`)
	reInitVendorIgnore = regexp.MustCompile(`^/?vendor(/\*?)?$`)
	initStabilities    = map[string]bool{"stable": true, "RC": true, "beta": true, "alpha": true, "dev": true}
)

// formatBlock is FormatterHelper::formatBlock for one message.
func formatBlock(message, style string, large bool) string {
	message = strings.ReplaceAll(message, "<", `\<`)
	line := " " + message + " "
	width := utf8.RuneCountInString(message) + 2
	if large {
		line = "  " + message + "  "
		width = utf8.RuneCountInString(message) + 4
	}
	var lines []string
	if large {
		lines = append(lines, strings.Repeat(" ", width))
	}
	lines = append(lines, line+strings.Repeat(" ", width-utf8.RuneCountInString(line)))
	if large {
		lines = append(lines, strings.Repeat(" ", width))
	}
	for i, l := range lines {
		lines[i] = "<" + style + ">" + l + "</" + style + ">"
	}
	return strings.Join(lines, "\n")
}

// configComposer is a Composer instance without a composer.json, standing
// for Factory::createConfig() and RepositoryFactory::createRepo/defaultRepos.
func (ic *initCmd) configComposer() (*composer.Composer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return composer.NewProjectComposer(ic.ctx.IO, phpjson.NewMap(), cwd, composer.Options{})
}

func strPtr(s string) *string { return &s }

// interact is InitCommand::interact.
func (ic *initCmd) interact() error {
	io, in := ic.ctx.IO, ic.ctx.In
	cc, err := ic.configComposer()
	if err != nil {
		return err
	}
	var repos []repo.Repository
	platformRepo, err := cc.CreatePlatformRepo(nil)
	if err != nil {
		return err
	}
	if list := in.Strings("repository"); len(list) > 0 {
		repos = []repo.Repository{platformRepo}
		createPackagist := true
		for _, r := range list {
			rc, err := repoConfigFromString(r, true)
			if err != nil {
				return err
			}
			if isPackagistFalse(rc) {
				createPackagist = false
				continue
			}
			created, err := cc.CreateRepositoryFromConfig(rc)
			if err != nil {
				return err
			}
			repos = append(repos, created)
		}
		if createPackagist {
			m := phpjson.NewMap()
			m.Set("type", "composer")
			m.Set("url", "https://repo.packagist.org")
			created, err := cc.CreateRepositoryFromConfig(m)
			if err != nil {
				return err
			}
			repos = append(repos, created)
		}
	}

	io.WriteError("")
	io.WriteError(formatBlock("Welcome to the Composer config generator", "bg=blue;fg=white", true))
	io.WriteError("")
	io.WriteError("")
	io.WriteError("This command will guide you through creating your composer.json config.")
	io.WriteError("")

	// name
	name, ok := ic.str("name")
	if !ok {
		name = ic.defaultPackageName()
	}
	v, err := io.AskAndValidate("Package name (<vendor>/<name>) [<comment>"+name+"</comment>]: ", func(v any) (any, error) {
		s, ok := cio.AnswerString(v)
		if !ok {
			return name, nil
		}
		if !reInitPackageName.MatchString(s) {
			return nil, &cio.InvalidArgumentError{Msg: "The package name " + s + " is invalid, it should be lowercase and have a vendor name, a forward slash, and a package name, matching: [a-z0-9_.-]+/[a-z0-9_.-]+"}
		}
		return s, nil
	}, 0, name)
	if err != nil {
		return err
	}
	ic.setOpt("name", v)

	// description
	desc, ok := ic.str("description")
	var descDef any
	if ok && desc != "" && desc != "0" {
		descDef = desc
	} else {
		desc = ""
	}
	if v, err = io.Ask("Description [<comment>"+desc+"</comment>]: ", descDef); err != nil {
		return err
	}
	ic.setOpt("description", v)

	// author
	author, hasAuthor := ic.str("author")
	if !hasAuthor {
		author, hasAuthor = ic.defaultAuthor()
	}
	var authorDef any
	prompt := "Author ["
	if hasAuthor {
		authorDef = author
		prompt += "<comment>" + author + "</comment>, "
	}
	v, err = io.AskAndValidate(prompt+"n to skip]: ", func(v any) (any, error) {
		s, _ := cio.AnswerString(v)
		if s == "n" || s == "no" {
			return nil, nil
		}
		if s == "" || s == "0" {
			if !hasAuthor {
				s = ""
			} else {
				s = author
			}
		}
		parsed, err := formatAuthors(s)
		if err != nil {
			return nil, &cio.InvalidArgumentError{Msg: err.Error()}
		}
		a := phpjson.ToMap(parsed[0])
		if email, ok := a.Get("email"); ok {
			return fmt.Sprintf("%s <%s>", a.Str("name"), email), nil
		}
		return a.Str("name"), nil
	}, 0, authorDef)
	if err != nil {
		return err
	}
	ic.setOpt("author", v)

	// minimum stability
	stab, ok := ic.str("stability")
	var stabDef any
	if ok && stab != "" && stab != "0" {
		stabDef = stab
	} else {
		stab = ""
	}
	v, err = io.AskAndValidate("Minimum Stability [<comment>"+stab+"</comment>]: ", func(v any) (any, error) {
		s, ok := cio.AnswerString(v)
		if !ok {
			return stabDef, nil
		}
		if !initStabilities[s] {
			return nil, &cio.InvalidArgumentError{Msg: `Invalid minimum stability "` + s + `". Must be empty or one of: stable, RC, beta, alpha, dev`}
		}
		return s, nil
	}, 0, stabDef)
	if err != nil {
		return err
	}
	ic.setOpt("stability", v)
	minimumStability, _ := cio.AnswerString(v)

	// type
	typ, ok := ic.str("type")
	var typDef any
	if ok {
		typDef = typ
	}
	if v, err = io.Ask("Package Type (e.g. library, project, metapackage, composer-plugin) [<comment>"+typ+"</comment>]: ", typDef); err != nil {
		return err
	}
	if s, ok := v.(string); ok && s == "" {
		v = nil
	}
	ic.setOpt("type", v)

	// license
	license, ok := ic.str("license")
	var licDef any
	if ok {
		licDef = license
	} else if l := os.Getenv("COMPOSER_DEFAULT_LICENSE"); l != "" && l != "0" {
		license, licDef = l, l
	}
	if v, err = io.Ask("License [<comment>"+license+"</comment>]: ", licDef); err != nil {
		return err
	}
	if s, ok := cio.AnswerString(v); ok && !spdx.Validate([]string{s}) && s != "proprietary" {
		return &composer.UserError{Msg: "Invalid license provided: " + s + ". Only SPDX license identifiers (https://spdx.org/licenses/) or \"proprietary\" are accepted.", Code: 1}
	}
	ic.setOpt("license", v)

	io.WriteError("")
	io.WriteError("Define your dependencies.")
	io.WriteError("")

	// prepare to resolve dependencies ($this->getRepos())
	if repos == nil {
		defaults, err := cc.DefaultRepositories()
		if err != nil {
			return err
		}
		repos = append([]repo.Repository{platformRepo}, defaults...)
	}
	d := &discovery{ctx: ic.ctx, repos: repos, platformRepo: platformRepo, minStability: func() string {
		s, ok := ic.str("stability")
		if !ok {
			s = "stable"
		}
		if n, err := semver.NormalizeStability(s); err == nil {
			return n
		}
		return s
	}}
	preferred := minimumStability
	if preferred == "" || preferred == "0" {
		preferred = "stable"
	}
	for _, key := range []string{"require", "require-dev"} {
		what := "dependencies (require)"
		if key == "require-dev" {
			what = "dev dependencies (require-dev)"
		}
		given := ic.strs(key)
		var reqs []string
		ask := len(given) > 0
		if !ask {
			if ask, err = io.AskConfirmation("Would you like to define your "+what+" interactively [<comment>yes</comment>]? ", true); err != nil {
				return err
			}
		}
		if ask {
			if reqs, err = d.determineRequirements(cc, given, preferred, true, false); err != nil {
				return err
			}
		}
		if reqs == nil {
			reqs = []string{}
		}
		ic.overList[key] = reqs
	}

	// autoload
	autoload, ok := ic.str("autoload")
	if !ok || autoload == "" || autoload == "0" {
		autoload = "src/"
	}
	n, _ := ic.str("name")
	ns := namespaceFromPackageName(n)
	v, err = io.AskAndValidate(`Add PSR-4 autoload mapping? Maps namespace "`+ns+`" to the entered relative path. [<comment>`+autoload+`</comment>, n to skip]: `, func(v any) (any, error) {
		s, ok := cio.AnswerString(v)
		if !ok {
			return autoload, nil
		}
		if s == "n" || s == "no" {
			return nil, nil
		}
		if s == "" || s == "0" {
			s = autoload
		}
		if !reInitAutoloadPath.MatchString(s) {
			return nil, &cio.InvalidArgumentError{Msg: fmt.Sprintf(`The src folder name "%s" is invalid. Please add a relative path with tailing forward slash. [A-Za-z0-9_-/]+/`, s)}
		}
		return s, nil
	}, 0, autoload)
	if err != nil {
		return err
	}
	ic.setOpt("autoload", v)
	return nil
}

// setOpt is $input->setOption() with an answer (nil: null).
func (ic *initCmd) setOpt(name string, v any) {
	s, ok := cio.AnswerString(v)
	if !ok {
		ic.over[name] = nil
		return
	}
	ic.over[name] = strPtr(s)
}

// gitQuestions asks for the vendor .gitignore entry and the composer.lock
// export-ignore attribute (interactive mode, in a git checkout).
func (ic *initCmd) gitQuestions() error {
	io := ic.ctx.IO
	if !io.Interactive || !isDirPath(".git") {
		return nil
	}
	ignoreFile := initRealpathOr(".gitignore")
	if !hasVendorIgnore(ignoreFile) {
		ok, err := io.AskConfirmation("Would you like the <info>vendor</info> directory added to your <info>.gitignore</info> [<comment>yes</comment>]? ", true)
		if err != nil {
			return err
		}
		if ok {
			if err := appendLine(ignoreFile, "/vendor/\n"); err != nil {
				return err
			}
		}
	}
	if !io.Interactive || !isDirPath(".git") {
		return nil
	}
	if typ, ok := ic.str("type"); ok && typ != "library" {
		return nil
	}
	attrFile := initRealpathOr(".gitattributes")
	if !hasExportIgnoreForLockfile(attrFile) {
		ok, err := io.AskConfirmation("Would you like <info>composer.lock</info> marked as export-ignore in your <info>.gitattributes</info> [<comment>yes</comment>]? ", true)
		if err != nil {
			return err
		}
		if ok {
			if err := appendLine(attrFile, "composer.lock export-ignore\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func isDirPath(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

// initRealpathOr is realpath($file) ?: realpath('.').'/'.$file.
func initRealpathOr(file string) string {
	if _, err := os.Lstat(file); err == nil {
		if r := fsutil.Realpath(file); r != "" {
			return r
		}
	}
	cwd, _ := os.Getwd()
	return filepath.Join(fsutil.Realpath(cwd), file)
}

func fileLines(path string) []string {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	s := strings.TrimSuffix(string(b), "\n")
	if s == "" && len(b) == 0 {
		return nil
	}
	return strings.Split(s, "\n")
}

// hasVendorIgnore is InitCommand::hasVendorIgnore.
func hasVendorIgnore(file string) bool {
	for _, l := range fileLines(file) {
		if reInitVendorIgnore.MatchString(l) {
			return true
		}
	}
	return false
}

func isPCRESpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// hasExportIgnoreForLockfile is InitCommand::hasExportIgnoreForLockfile
// ({^/?composer\.lock\s.*(?<!\S)export-ignore(?!\S)\s*$}).
func hasExportIgnoreForLockfile(file string) bool {
	for _, l := range fileLines(file) {
		rest := strings.TrimPrefix(l, "/")
		if !strings.HasPrefix(rest, "composer.lock") {
			continue
		}
		p := len(l) - len(rest) + len("composer.lock")
		if p >= len(l) || !isPCRESpace(l[p]) {
			continue
		}
		t := strings.TrimRightFunc(l, func(r rune) bool { return r < 128 && isPCRESpace(byte(r)) })
		if !strings.HasSuffix(t, "export-ignore") {
			continue
		}
		idx := len(t) - len("export-ignore")
		if idx-1 >= p && isPCRESpace(t[idx-1]) {
			return true
		}
	}
	return false
}

// appendLine is InitCommand::addVendorIgnore/addExportIgnoreForLockfile.
func appendLine(file, line string) error {
	contents := ""
	if b, err := os.ReadFile(file); err == nil {
		contents = string(b)
		if contents != "" && !strings.HasSuffix(contents, "\n") {
			contents += "\n"
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return &composer.UserError{Msg: "Could not read " + file, Code: 1}
	}
	return jsonfile.SafeFilePutContents(file, []byte(contents+line))
}

// updateDependencies is InitCommand::updateDependencies: update runs with an
// empty input on a fresh Composer instance.
func (ic *initCmd) updateDependencies() {
	ctx := ic.ctx
	cmd, err := ctx.App.findCommand("update")
	if err == nil {
		var in *Input
		in, err = parseInput(cmd, []string{"update"}, "update")
		if err == nil {
			applyEnvOptions(ctx.IO, cmd, in)
			_, err = cmd.Run(&Context{IO: ctx.IO, In: in, App: ctx.App, Cmd: cmd})
		}
	}
	if err != nil {
		ctx.IO.WriteError("Could not update dependencies. Run `composer update` to see more information.")
	}
}
