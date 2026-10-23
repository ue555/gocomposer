package cli

// Port of Composer\Command\InitCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano)
// and RepositoryFactory::configFromString.
// The interactive questionnaire (InitCommand::interact) is in init_interact.go.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/validate"
)

func init() { registerCommand(initCommand) }

func initCommand() *Command {
	return &Command{Name: "init", Run: runInit}
}

var (
	reInitPackageName = regexp.MustCompile(`^[a-z0-9]([_.-]?[a-z0-9]+)*/[a-z0-9](([_.]|-{1,2})?[a-z0-9]+)*$`)
	reInitAuthor      = regexp.MustCompile(`^(?P<name>[- .,\p{L}\p{N}\p{Mn}'’"()]+)(?:\s+<(?P<email>.+?)>)?\n?$`)
	reInitCamel       = regexp.MustCompile(`(?:([a-z])([A-Z])|([A-Z])([A-Z][a-z]))`)
	reInitTrim        = regexp.MustCompile(`^[_.-]+|[_.-]+$|[^a-z0-9_.-]`)
	reInitRepeated    = regexp.MustCompile(`([_.-]){2,}`)
	reInitNonAlnum    = regexp.MustCompile(`(?i)[^a-z0-9]`)
	reInitGitConfig   = regexp.MustCompile(`(?m)^([^=\n]+)=(.*)$`)
)

type initCmd struct {
	ctx       *Context
	gitConfig map[string]string
	// option values set by interact() ($input->setOption); a nil value is
	// an option set to null
	over     map[string]*string
	overList map[string][]string
}

// str is $input->getOption() of a string option (ok false for null).
func (ic *initCmd) str(name string) (string, bool) {
	if v, ok := ic.over[name]; ok {
		if v == nil {
			return "", false
		}
		return *v, true
	}
	return ic.ctx.In.String(name)
}

// strs is $input->getOption() of an array option.
func (ic *initCmd) strs(name string) []string {
	if v, ok := ic.overList[name]; ok {
		return v
	}
	return ic.ctx.In.Strings(name)
}

func runInit(ctx *Context) (int, error) {
	io := ctx.IO
	ic := &initCmd{ctx: ctx, over: map[string]*string{}, overList: map[string][]string{}}
	if !io.Interactive {
		// InitCommand::initialize (non-interactive defaults)
		if _, ok := ic.str("name"); !ok {
			n := ic.defaultPackageName()
			ic.over["name"] = &n
		}
		if _, ok := ic.str("author"); !ok {
			if a, ok := ic.defaultAuthor(); ok {
				ic.over["author"] = &a
			}
		}
	} else if err := ic.interact(); err != nil {
		return 1, err
	}
	name, _ := ic.str("name")
	author, hasAuthor := ic.str("author")

	// $options in definition order, without null/[] values
	options := phpjson.NewMap()
	options.Set("name", name)
	if v, ok := ic.str("description"); ok {
		options.Set("description", v)
	}
	if hasAuthor {
		options.Set("author", author)
	}
	for _, k := range []string{"type", "homepage"} {
		if v, ok := ic.str(k); ok {
			options.Set(k, v)
		}
	}
	requires := ic.strs("require")
	if len(requires) > 0 {
		options.Set("require", nil)
	}
	requireDevs := ic.strs("require-dev")
	if len(requireDevs) > 0 {
		options.Set("require-dev", nil)
	}
	stability, hasStability := ic.str("stability")
	if hasStability {
		options.Set("stability", stability)
	}
	if v, ok := ic.str("license"); ok {
		options.Set("license", v)
	}
	autoloadPath, hasAutoload := ic.str("autoload")
	if hasAutoload {
		options.Set("autoload", autoloadPath)
	}

	if !reInitPackageName.MatchString(name) {
		return 1, &composer.UserError{Msg: "The package name " + name + " is invalid, it should be lowercase and have a vendor name, a forward slash, and a package name, matching: [a-z0-9_.-]+/[a-z0-9_.-]+", Code: 1}
	}
	if hasAuthor {
		authors, err := formatAuthors(author)
		if err != nil {
			return 1, err
		}
		options.Delete("author")
		options.Set("authors", authors)
	}
	if repos := ic.strs("repository"); len(repos) > 0 {
		var list []phpjson.Value
		for _, r := range repos {
			cfg, err := repoConfigFromString(r, true)
			if err != nil {
				return 1, err
			}
			list = append(list, cfg)
		}
		options.Set("repositories", list)
	}
	if hasStability {
		options.Delete("stability")
		options.Set("minimum-stability", stability)
	}
	requireMap, err := initRequirements(requires)
	if err != nil {
		return 1, err
	}
	options.Set("require", requireMap)
	if len(requireDevs) > 0 {
		devMap, err := initRequirements(requireDevs)
		if err != nil {
			return 1, err
		}
		options.Set("require-dev", devMap)
	}
	if hasAutoload {
		psr4 := phpjson.NewMap()
		psr4.Set(namespaceFromPackageName(name)+`\`, autoloadPath)
		al := phpjson.NewMap()
		al.Set("psr-4", psr4)
		al.ForceObject = true
		options.Set("autoload", al)
	}

	file, err := composer.ComposerFileName()
	if err != nil {
		return 1, err
	}
	if io.Interactive {
		js, err := jsonfile.Encode(options, jsonfile.IndentDefault)
		if err != nil {
			return 1, err
		}
		io.WriteError("")
		io.WriteError(js)
		io.WriteError("")
		ok, err := io.AskConfirmation("Do you confirm generation [<comment>yes</comment>]? ", true)
		if err != nil {
			return 1, err
		}
		if !ok {
			io.WriteError("<error>Command aborted</error>")
			return 1, nil
		}
	} else {
		io.WriteError("Writing " + file)
	}
	if err := jsonfile.Write(file, options, jsonfile.IndentDefault); err != nil {
		return 1, err
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return 1, err
	}
	data, err := composer.ParseJSONFile(raw, file)
	if err != nil {
		return 1, err
	}
	if err := composer.ValidateSchemaLax(data, file); err != nil {
		io.WriteError("<error>Schema validation error, aborting</error>")
		io.WriteError(err.Error())
		os.Remove(file)
		return 1, nil
	}

	// PHP truthiness of $autoloadPath
	autoloadSet := hasAutoload && autoloadPath != "" && autoloadPath != "0"
	hasDeps := requireMap.Len() > 0 || len(requireDevs) > 0
	if autoloadSet {
		if err := os.MkdirAll(autoloadPath, 0o777); err != nil {
			return 1, err
		}
		if !hasDeps {
			ic.runDumpAutoload()
		}
	}
	if err := ic.gitQuestions(); err != nil {
		return 1, err
	}
	if io.Interactive && hasDeps {
		ok, err := io.AskConfirmation("Would you like to install dependencies now [<comment>yes</comment>]? ", true)
		if err != nil {
			return 1, err
		}
		if ok {
			ic.updateDependencies()
		}
	}
	if autoloadSet {
		ns := namespaceFromPackageName(name)
		io.WriteError(`PSR-4 autoloading configured. Use "<comment>namespace ` + ns + `;</comment>" in ` + autoloadPath)
		io.WriteError(`Include the Composer autoloader with: <comment>require 'vendor/autoload.php';</comment>`)
	}
	return 0, nil
}

// initRequirements is BaseCommand::formatRequirements, with an empty result
// encoded as an object.
func initRequirements(reqs []string) (*phpjson.Map, error) {
	m := phpjson.NewMap()
	m.ForceObject = true
	if len(reqs) == 0 {
		return m, nil
	}
	names, versions, err := formatRequirements(reqs)
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		m.Set(n, versions[n])
	}
	return m, nil
}

// runDumpAutoload is InitCommand::runDumpAutoloadCommand: dump-autoload is
// run with an empty input on a fresh Composer instance.
func (ic *initCmd) runDumpAutoload() {
	ctx := ic.ctx
	cmd, err := ctx.App.findCommand("dump-autoload")
	if err == nil {
		var in *Input
		in, err = parseInput(cmd, []string{"dump-autoload"}, "dump-autoload")
		if err == nil {
			applyEnvOptions(ctx.IO, cmd, in)
			_, err = cmd.Run(&Context{IO: ctx.IO, In: in, App: ctx.App, Cmd: cmd})
		}
	}
	if err != nil {
		ctx.IO.WriteError("Could not run dump-autoload.")
	}
}

// formatAuthors is InitCommand::formatAuthors / parseAuthorString.
func formatAuthors(author string) ([]phpjson.Value, error) {
	m := reInitAuthor.FindStringSubmatch(author)
	if m == nil {
		return nil, &composer.UserError{Msg: "Invalid author string.  Must be in the formats: Jane Doe or John Smith <john@example.com>", Code: 1}
	}
	a := phpjson.NewMap()
	a.Set("name", strings.TrimSpace(m[1]))
	if email := m[reInitAuthor.SubexpIndex("email")]; email != "" {
		if !validate.ValidEmail(email) {
			return nil, &composer.UserError{Msg: `Invalid email "` + email + `"`, Code: 1}
		}
		a.Set("email", email)
	}
	return []phpjson.Value{a}, nil
}

// namespaceFromPackageName is InitCommand::namespaceFromPackageName.
func namespaceFromPackageName(name string) string {
	if name == "" || !strings.Contains(name, "/") {
		return ""
	}
	var parts []string
	for _, part := range strings.Split(name, "/") {
		part = reInitNonAlnum.ReplaceAllString(part, " ")
		part = ucwords(part)
		parts = append(parts, strings.ReplaceAll(part, " ", ""))
	}
	return strings.Join(parts, `\`)
}

// ucwords is PHP ucwords() with the default delimiters.
func ucwords(s string) string {
	b := []byte(s)
	upNext := true
	for i, c := range b {
		if upNext && c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
		upNext = c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\f' || c == '\v'
	}
	return string(b)
}

func asciiLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

// sanitizePackageNameComponent is InitCommand::sanitizePackageNameComponent.
func sanitizePackageNameComponent(name string) string {
	name = reInitCamel.ReplaceAllStringFunc(name, func(m string) string {
		sm := reInitCamel.FindStringSubmatch(m)
		return sm[1] + sm[3] + "-" + sm[2] + sm[4]
	})
	name = asciiLower(name)
	name = reInitTrim.ReplaceAllString(name, "")
	return reInitRepeated.ReplaceAllString(name, "$1")
}

// getGitConfig is InitCommand::getGitConfig (`git config -l`).
func (ic *initCmd) getGitConfig() map[string]string {
	if ic.gitConfig != nil {
		return ic.gitConfig
	}
	ic.gitConfig = map[string]string{}
	out, err := exec.Command("git", "config", "-l").Output()
	if err != nil {
		return ic.gitConfig
	}
	for _, m := range reInitGitConfig.FindAllStringSubmatch(string(out), -1) {
		ic.gitConfig[m[1]] = m[2]
	}
	return ic.gitConfig
}

func envNonEmpty(name string) (string, bool) {
	v := os.Getenv(name)
	return v, v != "" && v != "0"
}

// defaultPackageName is InitCommand::getDefaultPackageName.
func (ic *initCmd) defaultPackageName() string {
	git := ic.getGitConfig()
	cwd, _ := os.Getwd()
	name := sanitizePackageNameComponent(filepath.Base(fsutil.Realpath(cwd)))
	vendor := name
	if v, ok := envNonEmpty("COMPOSER_DEFAULT_VENDOR"); ok {
		vendor = v
	} else if v, ok := git["github.user"]; ok {
		vendor = v
	} else if v, ok := envNonEmpty("USERNAME"); ok {
		vendor = v
	} else if v, ok := envNonEmpty("USER"); ok {
		vendor = v
	} else if u := currentScriptOwner(); u != "" {
		vendor = u
	}
	return sanitizePackageNameComponent(vendor) + "/" + name
}

// defaultAuthor is InitCommand::getDefaultAuthor.
func (ic *initCmd) defaultAuthor() (string, bool) {
	git := ic.getGitConfig()
	name, hasName := envNonEmpty("COMPOSER_DEFAULT_AUTHOR")
	if !hasName {
		name, hasName = git["user.name"]
	}
	email, hasEmail := envNonEmpty("COMPOSER_DEFAULT_EMAIL")
	if !hasEmail {
		email, hasEmail = git["user.email"]
	}
	if hasName && hasEmail {
		return fmt.Sprintf("%s <%s>", name, email), true
	}
	return "", false
}

// repoConfigFromString is RepositoryFactory::configFromString.
func repoConfigFromString(repository string, allowFilesystem bool) (phpjson.Value, error) {
	if len(repository) >= 4 && strings.EqualFold(repository[:4], "http") {
		m := phpjson.NewMap()
		m.Set("type", "composer")
		m.Set("url", repository)
		return m, nil
	}
	if filepath.Ext(repository) == ".json" {
		// JsonFile::read through the HttpDownloader (RemoteFilesystem
		// reading a local path with file_get_contents)
		raw, err := os.ReadFile(repository)
		if err != nil {
			reason := "Failed to open stream: No such file or directory"
			if os.IsPermission(err) {
				reason = "Failed to open stream: Permission denied"
			} else if !os.IsNotExist(err) {
				reason = err.Error()
			}
			msg := `The "` + repository + `" file could not be downloaded: ` + reason
			// RuntimeException wrapping the TransportException
			return nil, &composer.UserError{Msg: msg, Code: 1, Prev: &composer.UserError{Msg: msg, Code: 1}}
		}
		v, err := jsonfile.Parse(raw, repository)
		if err != nil {
			return nil, err
		}
		data := phpjson.ToMap(v)
		if data != nil && (phpjson.Truthy(get(data, "packages")) || phpjson.Truthy(get(data, "includes")) || phpjson.Truthy(get(data, "provider-includes"))) {
			m := phpjson.NewMap()
			m.Set("type", "composer")
			m.Set("url", "file://"+strings.ReplaceAll(fsutil.Realpath(repository), `\`, "/"))
			return m, nil
		}
		if allowFilesystem {
			// ['type' => 'filesystem', 'json' => JsonFile]: the JsonFile
			// object has no public properties and encodes as {}
			obj := phpjson.NewMap()
			obj.ForceObject = true
			m := phpjson.NewMap()
			m.Set("type", "filesystem")
			m.Set("json", obj)
			return m, nil
		}
		return nil, &composer.UserError{Msg: "Invalid repository URL (" + repository + ") given. This file does not contain a valid composer repository.", Code: 1}
	}
	if strings.HasPrefix(repository, "{") {
		return jsonfile.Parse([]byte(repository), "")
	}
	return nil, &composer.UserError{Msg: "Invalid repository url (" + httpdl.Sanitize(repository) + ") given. Has to be a .json file, an http url or a JSON object.", Code: 1}
}
