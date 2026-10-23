package cli

// Port of Composer\Command\RepositoryCommand and the initialize() of
// Composer\Command\BaseConfigCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// Shell completion of repository names/types is not ported.

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/configsource"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
)

func init() { registerCommand(repositoryCommand) }

func repositoryCommand() *Command {
	return &Command{Name: "repository", Run: runRepository}
}

// baseConfig is the state BaseConfigCommand::initialize sets up.
type baseConfig struct {
	cfg        *config.Config
	configFile string
	source     *configsource.Source
}

// initBaseConfig is BaseConfigCommand::initialize.
func initBaseConfig(ctx *Context) (*baseConfig, error) {
	in := ctx.In
	global := in.Bool("global")
	fileOpt, hasFile := in.String("file")
	hasFile = hasFile && in.present["file"]
	if global && hasFile {
		return nil, configErr("--file and --global can not be combined")
	}
	cwd, _ := os.Getwd()
	cfg, err := config.CreateConfig(cwd, func(m string) { ctx.IO.WriteError(m) })
	if err != nil {
		return nil, err
	}
	home := cfg.GetString("home")
	if global {
		cfg.SetBaseDir(home)
	}
	b := &baseConfig{cfg: cfg}
	switch {
	case global:
		b.configFile = home + "/config.json"
	case hasFile:
		b.configFile = fileOpt
	default:
		f, err := composer.ComposerFileName()
		if err != nil {
			return nil, err
		}
		b.configFile = f
	}
	if b.configFile == "composer.json" || b.configFile == "./composer.json" {
		if _, err := os.Stat(b.configFile); err != nil {
			rc, e1 := filepath.EvalSymlinks(cwd)
			rh, e2 := filepath.EvalSymlinks(home)
			if e1 == nil && e2 == nil && rc == rh {
				_ = os.WriteFile(b.configFile, []byte("{\n}\n"), 0o666)
			}
		}
	}
	b.source = configsource.New(b.configFile, false)
	if global && !fileExists(b.configFile) {
		d := phpjson.NewMap()
		d.Set("config", emptyObject())
		if err := touchAndWrite(b.configFile, d); err != nil {
			return nil, err
		}
	}
	if !fileExists(b.configFile) {
		return nil, configErr(`File "` + b.configFile + `" cannot be found in the current directory`)
	}
	return b, nil
}

var reJSONObjectStart = regexp.MustCompile(`^\s*\{`)

func runRepository(ctx *Context) (int, error) {
	b, err := initBaseConfig(ctx)
	if err != nil {
		return 1, err
	}
	in := ctx.In
	action := "list"
	if _, ok := in.args["action"]; ok {
		action = strings.ToLower(in.Arg("action"))
	}
	_, hasName := in.args["name"]
	_, hasArg1 := in.args["arg1"]
	_, hasArg2 := in.args["arg2"]
	name, arg1, arg2 := in.Arg("name"), in.Arg("arg1"), in.Arg("arg2")

	data, err := readJSONMap(b.configFile)
	if err != nil {
		return 1, err
	}
	b.cfg.MergeSource(data, b.configFile)
	repos := b.cfg.Repositories()
	src := b.source

	switch action {
	case "list", "ls", "show":
		listRepositories(ctx, repos)
		return 0, nil

	case "add":
		if !hasName {
			return 1, configErr("You must pass a repository name. Example: composer repo add foo vcs https://example.org")
		}
		if !hasArg1 {
			return 1, configErr("You must pass the type and a url, or a JSON string.")
		}
		var repoConfig phpjson.Value
		if reJSONObjectStart.MatchString(arg1) {
			v, err := jsonfile.Parse([]byte(arg1), "")
			if err != nil {
				return 1, err
			}
			repoConfig = v
		} else {
			if !hasArg2 {
				return 1, configErr("You must pass the type and a url. Example: composer repo add foo vcs https://example.org")
			}
			m := phpjson.NewMap()
			m.Set("type", arg1)
			m.Set("url", arg2)
			repoConfig = m
		}
		before, hasBefore := in.String("before")
		hasBefore = hasBefore && in.present["before"]
		after, hasAfter := in.String("after")
		hasAfter = hasAfter && in.present["after"]
		if hasBefore && hasAfter {
			return 1, configErr("You can not combine --before and --after")
		}
		if hasBefore || hasAfter {
			if repoConfig == false {
				return 1, configErr("Cannot use --before/--after with boolean repository values")
			}
			ref, offset := before, 0
			if hasAfter {
				ref, offset = after, 1
			}
			return srcResult(src.InsertRepository(name, repoConfig, ref, offset))
		}
		return srcResult(src.AddRepository(name, repoConfig, in.Bool("append")))

	case "remove", "rm", "delete":
		if !hasName {
			return 1, configErr("You must pass the repository name to remove.")
		}
		if err := src.RemoveRepository(name); err != nil {
			return 1, withPrevious(err)
		}
		if name == "packagist" || name == "packagist.org" {
			return srcResult(src.AddRepository("packagist.org", false, false))
		}
		return 0, nil

	case "set-url", "seturl":
		if !hasName || !hasArg1 {
			return 1, configErr("Usage: composer repo set-url <name> <new-url>")
		}
		return srcResult(src.SetRepositoryUrl(name, arg1))

	case "get-url", "geturl":
		if !hasName {
			return 1, configErr("Usage: composer repo get-url <name>")
		}
		if repos != nil {
			if v, ok := repos.Get(name); ok && phpjson.IsArray(v) {
				u, _ := phpjson.ToMap(v).Get("url")
				s, isStr := u.(string)
				if !isStr {
					return 1, configErr("The " + name + " repository does not have a URL")
				}
				ctx.IO.Write(s)
				return 0, nil
			}
			for _, k := range repos.Keys() {
				v, _ := repos.Get(k)
				m := phpjson.ToMap(v)
				if m == nil || !phpjson.IsArray(v) {
					continue
				}
				if n, ok := m.Get("name"); ok && n == name {
					u, _ := m.Get("url")
					s, ok := u.(string)
					if !ok {
						return 1, configErr("The " + name + " repository does not have a URL")
					}
					ctx.IO.Write(s)
					return 0, nil
				}
			}
		}
		return 1, configErr("There is no " + name + " repository defined")

	case "disable":
		if !hasName {
			return 1, configErr("Usage: composer repo disable packagist.org")
		}
		if name == "packagist" || name == "packagist.org" {
			return srcResult(src.AddRepository("packagist.org", false, in.Bool("append")))
		}
		return 1, configErr("Only packagist.org can be enabled/disabled using this command. Use add/remove for other repositories.")

	case "enable":
		if !hasName {
			return 1, configErr("Usage: composer repo enable packagist.org")
		}
		if name == "packagist" || name == "packagist.org" {
			return srcResult(src.RemoveRepository("packagist.org"))
		}
		return 1, configErr("Only packagist.org can be enabled/disabled using this command.")
	}
	return 1, configErr(`Unknown action "` + action + `". Use list, add, remove, set-url, get-url, enable, disable`)
}

// listRepositories is RepositoryCommand::listRepositories.
func listRepositories(ctx *Context, repos *phpjson.Map) {
	io := ctx.IO
	type entry struct {
		key string
		val phpjson.Value
	}
	var list []entry
	nextInt := int64(0)
	if repos != nil {
		for _, k := range repos.Keys() {
			v, _ := repos.Get(k)
			list = append(list, entry{k, v})
			if n, ok := phpjson.IntKey(k); ok && n >= nextInt {
				nextInt = n + 1
			}
		}
	}
	packagistPresent := false
	for _, e := range list {
		m := phpjson.ToMap(e.val)
		if m == nil || !phpjson.IsArray(e.val) {
			continue
		}
		t, hasT := m.Get("type")
		u, hasU := m.Get("url")
		if hasT && hasU && t != nil && u != nil && t == "composer" {
			us, _ := phpjson.ScalarString(u)
			if strings.HasSuffix(parseURLHost(us), "packagist.org") {
				packagistPresent = true
				break
			}
		}
	}
	if !packagistPresent {
		m := phpjson.NewMap()
		m.Set("packagist.org", false)
		list = append(list, entry{strconv.FormatInt(nextInt, 10), m})
	}
	if len(list) == 0 {
		io.Write("No repositories configured")
		return
	}
	for _, e := range list {
		if e.val == false {
			io.Write("[" + e.key + "] <info>disabled</info>")
			continue
		}
		if !phpjson.IsArray(e.val) {
			continue
		}
		m := phpjson.ToMap(e.val)
		if m.Len() == 1 {
			first := m.Keys()[0]
			if v, _ := m.Get(first); v == false {
				io.Write("[" + first + "] <info>disabled</info>")
				continue
			}
		}
		name := e.key
		if v, ok := m.Get("name"); ok && v != nil {
			name, _ = phpjson.ScalarString(v)
		}
		typ := "unknown"
		if v, ok := m.Get("type"); ok && v != nil {
			typ, _ = phpjson.ScalarString(v)
		}
		var u string
		if v, ok := m.Get("url"); ok && v != nil {
			u, _ = phpjson.ScalarString(v)
		} else {
			u, _ = phpjson.Encode(m, phpjson.ComposerFlags)
		}
		io.Write("[" + name + "] <info>" + typ + "</info> " + u)
	}
}

// parseURLHost is (string) parse_url($url, PHP_URL_HOST).
func parseURLHost(s string) string {
	pu, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return pu.Hostname()
}

// srcResult maps a JsonConfigSource write result to the command result,
// keeping the chained JsonValidationException.
func srcResult(err error) (int, error) {
	if err != nil {
		return 1, withPrevious(err)
	}
	return 0, nil
}
