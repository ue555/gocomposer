package cli

// Port of Composer\Command\ConfigCommand and Composer\Command\BaseConfigCommand
// (src/Composer/Command/ConfigCommand.php and BaseConfigCommand.php of
// composer/composer f83588392df43d4fd027d84d26de459946ef4e80).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Shell completion of setting keys (suggestSettingKeys) is not ported.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/configsource"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/policy"
)

func configCommand() *Command {
	return &Command{Name: "config", Run: runConfig}
}

// configurablePackageProperties is ConfigCommand::CONFIGURABLE_PACKAGE_PROPERTIES.
var configurablePackageProperties = []string{"name", "type", "description", "homepage", "version", "minimum-stability", "prefer-stable", "keywords", "license", "repositories", "suggest", "extra"}

type configCmd struct {
	ctx            *Context
	io             *cio.IO
	in             *Input
	cfg            *config.Config
	configFile     string
	configSource   *configsource.Source
	authConfigFile string
	authSource     *configsource.Source
}

func configErr(msg string) error { return &composer.UserError{Msg: msg, Code: 1} }

// fileExists is JsonFile::exists() (is_file).
func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.Mode().IsRegular()
}

// readJSON is JsonFile::read().
func readJSON(path string) (phpjson.Value, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, configErr("Could not read " + path + "\n\nThe file \"" + path + "\" is not readable.")
	}
	defer f.Close()
	var sb strings.Builder
	buf := make([]byte, 32*1024)
	for {
		n, rerr := f.Read(buf)
		sb.Write(buf[:n])
		if rerr != nil {
			break
		}
	}
	return jsonfile.Parse([]byte(sb.String()), path)
}

func readJSONMap(path string) (*phpjson.Map, error) {
	v, err := readJSON(path)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	return m, nil
}

// touchAndWrite is touch() + JsonFile::write() + chmod 0600 used to
// initialise the global config/auth files.
func touchAndWrite(path string, data *phpjson.Map) error {
	if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o666); err == nil {
		f.Close()
	}
	if err := jsonfile.Write(path, data, jsonfile.IndentDefault); err != nil {
		return err
	}
	_ = os.Chmod(path, 0o600)
	return nil
}

func emptyObject() *phpjson.Map { return &phpjson.Map{ForceObject: true} }

// initialize is BaseConfigCommand::initialize + ConfigCommand::initialize.
func (c *configCmd) initialize() error {
	global := c.in.Bool("global")
	fileOpt, hasFile := c.in.String("file")
	if global && hasFile {
		return configErr("--file and --global can not be combined")
	}
	cwd, _ := os.Getwd()
	cfg, err := config.CreateConfig(cwd, func(m string) { c.io.WriteError(m) })
	if err != nil {
		return err
	}
	c.cfg = cfg
	home := cfg.GetString("home")
	if global {
		cfg.SetBaseDir(home)
	}

	if global {
		c.configFile = home + "/config.json"
	} else if hasFile {
		c.configFile = fileOpt
	} else {
		f, err := composer.ComposerFileName()
		if err != nil {
			return err
		}
		c.configFile = f
	}

	// Create global composer.json if invoked using `composer global [config-cmd]`
	if c.configFile == "composer.json" || c.configFile == "./composer.json" {
		if _, err := os.Stat(c.configFile); err != nil {
			rc, e1 := filepath.EvalSymlinks(cwd)
			rh, e2 := filepath.EvalSymlinks(home)
			if e1 == nil && e2 == nil && rc == rh {
				_ = os.WriteFile(c.configFile, []byte("{\n}\n"), 0o666)
			}
		}
	}

	c.configSource = configsource.New(c.configFile, false)

	// Initialize the global file if it's not there, ignoring any warnings or notices
	if global && !fileExists(c.configFile) {
		d := phpjson.NewMap()
		d.Set("config", emptyObject())
		if err := touchAndWrite(c.configFile, d); err != nil {
			return err
		}
	}

	if !fileExists(c.configFile) {
		return configErr(fmt.Sprintf(`File "%s" cannot be found in the current directory`, c.configFile))
	}

	// ConfigCommand::initialize
	if global {
		c.authConfigFile = home + "/auth.json"
	} else {
		c.authConfigFile = phpDirname(c.configFile) + "/auth.json"
	}
	c.authSource = configsource.New(c.authConfigFile, true)
	if global && !fileExists(c.authConfigFile) {
		d := phpjson.NewMap()
		for _, k := range []string{"bitbucket-oauth", "github-oauth", "gitlab-oauth", "gitlab-token", "http-basic", "bearer", "forgejo-token"} {
			d.Set(k, emptyObject())
		}
		if err := touchAndWrite(c.authConfigFile, d); err != nil {
			return err
		}
	}
	return nil
}

// phpDirname is PHP's dirname() for unix paths.
func phpDirname(p string) string {
	if p == "" {
		return ""
	}
	trimmed := strings.TrimRight(p, "/")
	if trimmed == "" {
		return "/"
	}
	i := strings.LastIndex(trimmed, "/")
	if i < 0 {
		return "."
	}
	d := strings.TrimRight(trimmed[:i], "/")
	if d == "" {
		return "/"
	}
	return d
}

func runConfig(ctx *Context) (int, error) {
	c := &configCmd{ctx: ctx, io: ctx.IO, in: ctx.In}
	if err := c.initialize(); err != nil {
		return 1, err
	}
	return c.execute()
}

// shellQuote is escapeshellarg().
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func (c *configCmd) execute() (int, error) {
	in := c.in
	// Open file in editor
	if in.Bool("editor") {
		editor, _ := os.LookupEnv("EDITOR")
		if editor == "" {
			for _, candidate := range []string{"editor", "vim", "vi", "nano", "pico", "ed"} {
				if out, err := exec.Command("sh", "-c", "which "+candidate).Output(); err == nil && strings.TrimSpace(string(out)) != "" {
					editor = candidate
					break
				}
			}
		} else {
			editor = shellQuote(editor)
		}
		file := c.configFile
		if in.Bool("auth") {
			file = c.authConfigFile
		}
		cmd := exec.Command("sh", "-c", editor+" "+shellQuote(file)+" > `tty`")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return 0, nil
	}

	if !in.Bool("global") {
		data, err := readJSONMap(c.configFile)
		if err != nil {
			return 1, err
		}
		c.cfg.MergeSource(data, c.configFile)
		authData := phpjson.NewMap()
		if fileExists(c.authConfigFile) {
			a, err := readJSONMap(c.authConfigFile)
			if err != nil {
				return 1, err
			}
			authData = a
		}
		wrap := phpjson.NewMap()
		wrap.Set("config", authData)
		c.cfg.MergeSource(wrap, c.authConfigFile)
	}

	// List the configuration of the file settings
	if in.Bool("list") {
		all, err := c.cfg.All(0)
		if err != nil {
			return 1, err
		}
		if err := c.listConfiguration(all, c.cfg.RawAll(), "", true, in.Bool("source")); err != nil {
			return 1, err
		}
		return 0, nil
	}

	keyArgs := in.Args("setting-key")
	if len(keyArgs) == 0 {
		return 0, nil
	}
	settingKey := keyArgs[0]
	values := in.Args("setting-value")
	unset := in.Bool("unset")

	if len(values) > 0 && unset {
		return 1, configErr("You can not combine a setting value with --unset")
	}

	// show the value if no value is provided
	if len(values) == 0 && !unset {
		return c.showValue(settingKey)
	}
	if err := c.setValue(settingKey, values); err != nil {
		return 1, withPrevious(err)
	}
	return 0, nil
}

// withPrevious keeps the JsonValidationException that JsonConfigSource
// chains as previous exception, so that it is rendered as a second box.
func withPrevious(err error) error {
	var je *jsonmanip.Exception
	var ve *configsource.ValidationError
	if errors.As(err, &je) && je.Err != nil && errors.As(je.Err, &ve) {
		return &composer.UserError{Msg: je.Message, Code: 1, Prev: ve}
	}
	return err
}

var reRepoKey = regexp.MustCompile(`^repos?(?:itories)?(?:\.(.+))?`)

func (c *configCmd) showValue(settingKey string) (int, error) {
	propertiesDefaults := phpjson.NewMap()
	propertiesDefaults.Set("type", "library")
	propertiesDefaults.Set("description", "")
	propertiesDefaults.Set("homepage", "")
	propertiesDefaults.Set("minimum-stability", "stable")
	propertiesDefaults.Set("prefer-stable", false)
	propertiesDefaults.Set("keywords", []phpjson.Value{})
	propertiesDefaults.Set("license", []phpjson.Value{})
	propertiesDefaults.Set("suggest", []phpjson.Value{})
	propertiesDefaults.Set("extra", []phpjson.Value{})

	rawData, err := readJSON(c.configFile)
	if err != nil {
		return 1, err
	}
	all, err := c.cfg.All(0)
	if err != nil {
		return 1, err
	}
	source, err := c.cfg.SourceOfValue(settingKey)
	if err != nil {
		return 1, err
	}
	allConfig, _ := all.Get("config")

	var value phpjson.Value
	if m := reRepoKey.FindStringSubmatch(settingKey); m != nil {
		repos, _ := all.Get("repositories")
		if m[1] == "" {
			value = repos
			if value == nil {
				value = []phpjson.Value{}
			}
		} else {
			v, ok := phpIsset(repos, m[1])
			if !ok {
				return 1, configErr("There is no " + m[1] + " repository defined")
			}
			value = v
		}
	} else if strings.Index(settingKey, ".") > 0 {
		bits := strings.Split(settingKey, ".")
		var data phpjson.Value
		if bits[0] == "extra" || bits[0] == "suggest" {
			data = rawData
		} else {
			data = allConfig
		}
		match := false
		key, keySet := "", false
		for _, bit := range bits {
			if keySet {
				key = key + "." + bit
			} else {
				key, keySet = bit, true
			}
			match = false
			if v, ok := phpIsset(data, key); ok {
				match = true
				data = v
				keySet = false
			}
		}
		if !match {
			return 1, configErr(settingKey + " is not defined.")
		}
		value = data
	} else if _, ok := phpIsset(allConfig, settingKey); ok {
		flags := config.RelativePaths
		if c.in.Bool("absolute") {
			flags = 0
		}
		v, err := c.cfg.GetFlags(settingKey, flags)
		if err != nil {
			return 1, err
		}
		value = v
		// ensure we get {} output for properties which are objects
		if phpjson.IsArray(value) && phpjson.ArrayLen(value) == 0 && schemaConfigTypeIsObject(settingKey) {
			value = emptyObject()
		}
	} else if v, ok := phpIsset(rawData, settingKey); ok && inStringList(settingKey, configurablePackageProperties) {
		value = v
		source = c.configFile
	} else if v, ok := phpIsset(propertiesDefaults, settingKey); ok {
		value = v
		source = "defaults"
	} else {
		return 1, configErr(settingKey + " is not defined")
	}

	var out string
	switch value.(type) {
	case *phpjson.Map, []phpjson.Value, bool:
		s, err := phpjson.Encode(value, phpjson.UnescapedSlashes|phpjson.UnescapedUnicode)
		if err != nil {
			return 1, err
		}
		out = s
	default:
		out = config.PHPString(value)
	}
	if c.in.Bool("source") {
		out += " (" + source + ")"
	}
	c.io.WriteV(out, true, cio.Quiet)
	return 0, nil
}

var (
	schemaOnce sync.Once
	schemaData phpjson.Value
)

// schemaConfigTypeIsObject checks
// $schema['properties']['config']['properties'][$key]['type'] for 'object'.
func schemaConfigTypeIsObject(key string) bool {
	schemaOnce.Do(func() {
		schemaData, _ = phpjson.Decode(jsonschema.ComposerSchemaJSON)
	})
	t := dig(schemaData, "properties", "config", "properties", key, "type")
	if s, ok := t.(string); ok {
		return s == "object"
	}
	found := false
	phpjson.Each(t, func(_ string, v phpjson.Value) {
		if v == "object" {
			found = true
		}
	})
	return found
}

func dig(v phpjson.Value, keys ...string) phpjson.Value {
	for _, k := range keys {
		m := phpjson.ToMap(v)
		if m == nil {
			return nil
		}
		v, _ = m.Get(k)
	}
	return v
}

var reIntKey = regexp.MustCompile(`^(?:0|-?[1-9][0-9]*)$`)

// phpIsset is isset($data[$key]) returning $data[$key].
func phpIsset(data phpjson.Value, key string) (phpjson.Value, bool) {
	switch t := data.(type) {
	case *phpjson.Map:
		if t == nil {
			return nil, false
		}
		v, ok := t.Get(key)
		return v, ok && v != nil
	case []phpjson.Value:
		if !reIntKey.MatchString(key) {
			return nil, false
		}
		n, err := strconv.Atoi(key)
		if err != nil || n < 0 || n >= len(t) {
			return nil, false
		}
		return t[n], t[n] != nil
	case string:
		// string offsets: only integer(-like) offsets are set
		if !reIntKey.MatchString(key) {
			return nil, false
		}
		n, err := strconv.Atoi(key)
		if err != nil {
			return nil, false
		}
		if n < 0 {
			n += len(t)
		}
		if n < 0 || n >= len(t) {
			return nil, false
		}
		return t[n : n+1], true
	}
	return nil, false
}

func inStringList(s string, list []string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

var reIsNumeric = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?[ \t\n\r\v\f]*$`)

// isNumeric is is_numeric() for strings.
func isNumeric(s string) bool { return reIsNumeric.MatchString(s) }

var reConfigPrefix = regexp.MustCompile(`^config\.`)
var reNonAlnum = regexp.MustCompile(`(?i)[^a-z0-9]`)
var reDashes = regexp.MustCompile(`-+`)

// listConfiguration is ConfigCommand::listConfiguration; kNull is $k === null.
func (c *configCmd) listConfiguration(contents, rawContents phpjson.Value, k string, kNull bool, showSource bool) error {
	origK, origNull := k, kNull
	var ferr error
	phpjson.Each(contents, func(key string, value phpjson.Value) {
		if ferr != nil {
			return
		}
		if kNull && key != "config" && key != "repositories" {
			return
		}
		rawVal, _ := phpIsset(rawContents, key)

		if phpjson.IsArray(value) {
			firstKey, hasFirst := "", false
			phpjson.Each(value, func(fk string, _ phpjson.Value) {
				if !hasFirst {
					firstKey, hasFirst = fk, true
				}
			})
			if !(hasFirst && isNumeric(firstKey)) || (key == "repositories" && kNull) {
				k += reConfigPrefix.ReplaceAllString(key+".", "")
				kNull = false
				if !phpjson.IsArray(rawVal) {
					rawVal = nil
				}
				ferr = c.listConfiguration(value, rawVal, k, false, showSource)
				k, kNull = origK, origNull
				return
			}
		}

		valStr, valIsString := "", false
		switch t := value.(type) {
		case *phpjson.Map, []phpjson.Value:
			var parts []string
			phpjson.Each(t, func(_ string, e phpjson.Value) {
				if phpjson.IsArray(e) {
					s, _ := phpjson.Encode(e, 0)
					parts = append(parts, s)
				} else {
					parts = append(parts, config.PHPString(e))
				}
			})
			valStr, valIsString = "["+strings.Join(parts, ", ")+"]", true
		case bool:
			valStr, valIsString = "false", true
			if t {
				valStr = "true"
			}
		case string:
			valStr, valIsString = t, true
		default:
			valStr = config.PHPString(t)
		}

		source := ""
		if showSource {
			s, err := c.cfg.SourceOfValue(k + key)
			if err != nil {
				ferr = err
				return
			}
			source = " (" + s + ")"
		}

		var link string
		if !kNull && strings.HasPrefix(k, "repositories") {
			link = "https://getcomposer.org/doc/05-repositories.md"
		} else {
			id := k
			if k == "" || kNull {
				id = key
			}
			if i := strings.Index(id, "."); i >= 0 {
				id = id[:i]
			}
			id = reNonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(id)), "-")
			id = reDashes.ReplaceAllString(id, "-")
			link = "https://getcomposer.org/doc/06-config.md#" + id
		}
		if rs, ok := rawVal.(string); ok && (!valIsString || rs != valStr) {
			c.io.WriteV("[<fg=yellow;href="+link+">"+k+key+"</>] <info>"+rs+" ("+valStr+")</info>"+source, true, cio.Quiet)
		} else {
			c.io.WriteV("[<fg=yellow;href="+link+">"+k+key+"</>] <info>"+valStr+"</info>"+source, true, cio.Quiet)
		}
	})
	return ferr
}

// validator returns true, false or an error message (PHP's mixed result).
type validator func(v phpjson.Value) (bool, string)
type normalizer func(v phpjson.Value) phpjson.Value
type valueHandler struct {
	validate  validator
	normalize normalizer
}

func boolValidator(v phpjson.Value) (bool, string) {
	return inValues(v, "true", "false", "1", "0"), ""
}

func boolNormalizer(v phpjson.Value) phpjson.Value {
	s, _ := v.(string)
	return s != "false" && phpjson.Truthy(s)
}

func keepAsIs(v phpjson.Value) phpjson.Value { return v }

func inValues(v phpjson.Value, list ...string) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return inStringList(s, list)
}

func enumValidator(list ...string) validator {
	return func(v phpjson.Value) (bool, string) { return inValues(v, list...), "" }
}

func isStringValidator(v phpjson.Value) (bool, string) { _, ok := v.(string); return ok, "" }

func isNumericValidator(v phpjson.Value) (bool, string) {
	s, ok := v.(string)
	return ok && isNumeric(s), ""
}

func intvalNormalizer(v phpjson.Value) phpjson.Value { return config.PHPIntval(v) }

// boolOrKeyword normalizes like the store-auths/discard-changes/... closures.
func boolOrKeyword(keywords ...string) normalizer {
	return func(v phpjson.Value) phpjson.Value {
		if inValues(v, keywords...) {
			return v
		}
		return boolNormalizer(v)
	}
}

func nullStringNormalizer(v phpjson.Value) phpjson.Value {
	if v == "null" {
		return nil
	}
	return v
}

func isReadable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

var reCacheSize = regexp.MustCompile(`(?i)^\s*([0-9.]+)\s*(?:([kmg])(?:i?b)?)?\s*$`)

func uniqueConfigValues() map[string]valueHandler {
	audit := enumValidator(policy.AuditIgnore, policy.AuditReport, policy.AuditFail)
	blockScope := enumValidator(policy.AvailableBlockScopes...)
	b := valueHandler{boolValidator, boolNormalizer}
	str := valueHandler{isStringValidator, keepAsIs}
	num := valueHandler{isNumericValidator, intvalNormalizer}
	return map[string]valueHandler{
		"process-timeout":   num,
		"use-include-path":  b,
		"use-github-api":    b,
		"preferred-install": {enumValidator("auto", "source", "dist"), keepAsIs},
		"gitlab-protocol":   {enumValidator("git", "http", "https"), keepAsIs},
		"store-auths":       {enumValidator("true", "false", "prompt"), boolOrKeyword("prompt")},
		"notify-on-install": b,
		"vendor-dir":        str,
		"bin-dir":           str,
		"archive-dir":       str,
		"archive-format":    str,
		"data-dir":          str,
		"cache-dir":         str,
		"cache-files-dir":   str,
		"cache-repo-dir":    str,
		"cache-vcs-dir":     str,
		"cache-ttl":         num,
		"cache-files-ttl":   num,
		"cache-files-maxsize": {func(v phpjson.Value) (bool, string) {
			s, _ := v.(string)
			return reCacheSize.MatchString(s), ""
		}, keepAsIs},
		"bin-compat":                  {enumValidator("auto", "full", "proxy", "symlink"), keepAsIs},
		"discard-changes":             {enumValidator("stash", "true", "false", "1", "0"), boolOrKeyword("stash")},
		"autoloader-suffix":           {isStringValidator, nullStringNormalizer},
		"sort-packages":               b,
		"optimize-autoloader":         b,
		"strict-psr-autoloader":       b,
		"classmap-authoritative":      b,
		"apcu-autoloader":             b,
		"prepend-autoloader":          b,
		"update-with-minimal-changes": b,
		"disable-tls":                 b,
		"secure-http":                 b,
		"source-fallback":             b,
		"bump-after-update":           {enumValidator("dev", "no-dev", "true", "false", "1", "0"), boolOrKeyword("dev", "no-dev")},
		"cafile": {func(v phpjson.Value) (bool, string) {
			s, _ := v.(string)
			_, err := os.Stat(s)
			return s != "" && err == nil && isReadable(s), ""
		}, nullStringNormalizer},
		"capath": {func(v phpjson.Value) (bool, string) {
			s, _ := v.(string)
			st, err := os.Stat(s)
			return s != "" && err == nil && st.IsDir() && isReadable(s), ""
		}, nullStringNormalizer},
		"github-expose-hostname":     b,
		"htaccess-protect":           b,
		"lock":                       b,
		"allow-plugins":              b,
		"platform-check":             {enumValidator("php-only", "true", "false", "1", "0"), boolOrKeyword("php-only")},
		"use-parent-dir":             {enumValidator("true", "false", "prompt"), boolOrKeyword("prompt")},
		"audit.abandoned":            {audit, keepAsIs},
		"audit.ignore-unreachable":   b,
		"audit.block-insecure":       b,
		"audit.block-abandoned":      b,
		"policy.advisories.block":    b,
		"policy.advisories.audit":    {audit, keepAsIs},
		"policy.malware.block":       b,
		"policy.malware.block-scope": {blockScope, keepAsIs},
		"policy.malware.audit":       {audit, keepAsIs},
		"policy.abandoned.block":     b,
		"policy.abandoned.audit":     {audit, keepAsIs},
		"policy.cooldown.block":      b,
		"policy.cooldown.audit":      {audit, keepAsIs},
		"policy.cooldown.period": {func(v phpjson.Value) (bool, string) {
			if v == "null" || v == "" {
				return true, ""
			}
			_, err := policy.ParseDuration(v)
			return err == nil, ""
		}, func(v phpjson.Value) phpjson.Value {
			if v == "null" || v == "" {
				return nil
			}
			if s, ok := v.(string); ok && isNumeric(s) {
				return config.PHPIntval(s)
			}
			return v
		}},
		"policy.ignore-unreachable": b,
	}
}

func listValues(vals []string) []phpjson.Value {
	out := make([]phpjson.Value, len(vals))
	for i, v := range vals {
		out[i] = v
	}
	return out
}

func arrayExpected(v phpjson.Value) (bool, string) { return true, "" }

func multiConfigValues() map[string]valueHandler {
	severity := valueHandler{func(v phpjson.Value) (bool, string) {
		var bad bool
		phpjson.Each(v, func(_ string, e phpjson.Value) {
			if !inValues(e, "low", "medium", "high", "critical") {
				bad = true
			}
		})
		if bad {
			return false, "valid severities include: low, medium, high, critical"
		}
		return true, ""
	}, keepAsIs}
	return map[string]valueHandler{
		"github-protocols": {func(v phpjson.Value) (bool, string) {
			var bad bool
			phpjson.Each(v, func(_ string, e phpjson.Value) {
				if !inValues(e, "git", "https", "ssh") {
					bad = true
				}
			})
			if bad {
				return false, "valid protocols include: git, https, ssh"
			}
			return true, ""
		}, keepAsIs},
		"github-domains":                    {arrayExpected, keepAsIs},
		"gitlab-domains":                    {arrayExpected, keepAsIs},
		"audit.ignore-severity":             severity,
		"policy.advisories.ignore-severity": severity,
		"policy.malware.ignore-source":      {arrayExpected, keepAsIs},
	}
}

func uniqueProps() map[string]valueHandler {
	str := valueHandler{isStringValidator, keepAsIs}
	return map[string]valueHandler{
		"name":        str,
		"type":        str,
		"description": str,
		"homepage":    str,
		"version":     str,
		"minimum-stability": {func(v phpjson.Value) (bool, string) {
			// VersionParser::normalizeStability throws for unknown stabilities
			return true, ""
		}, func(v phpjson.Value) phpjson.Value {
			s := strings.ToLower(v.(string))
			if s == "rc" {
				return "RC"
			}
			return s
		}},
		"prefer-stable": {boolValidator, boolNormalizer},
	}
}

func multiProps() map[string]valueHandler {
	return map[string]valueHandler{
		"keywords": {arrayExpected, keepAsIs},
		"license":  {arrayExpected, keepAsIs},
	}
}

// handleSingleValue is ConfigCommand::handleSingleValue.
func (c *configCmd) handleSingleValue(key string, h valueHandler, values []string, add func(string, phpjson.Value) error) error {
	if len(values) != 1 {
		return configErr("You can only pass one value. Example: php composer.phar config process-timeout 300")
	}
	if key == "minimum-stability" {
		s := strings.ToLower(values[0])
		if !inStringList(s, []string{"stable", "rc", "beta", "alpha", "dev"}) {
			return configErr(`Invalid stability string "` + s + `", expected one of stable, RC, beta, alpha or dev`)
		}
	}
	if ok, msg := h.validate(values[0]); !ok {
		suffix := ""
		if msg != "" {
			suffix = " (" + msg + ")"
		}
		return configErr(fmt.Sprintf(`"%s" is an invalid value`+suffix, values[0]))
	}
	normalized := h.normalize(values[0])
	if key == "disable-tls" {
		cur, err := c.cfg.GetFlags("disable-tls", 0)
		if err != nil {
			return err
		}
		if !phpjson.Truthy(normalized) && phpjson.Truthy(cur) {
			c.io.WriteError("<info>You are now running Composer with SSL/TLS protection enabled.</info>")
		} else if phpjson.Truthy(normalized) && !phpjson.Truthy(cur) {
			c.io.WriteError("<warning>You are now running Composer with SSL/TLS protection disabled.</warning>")
		}
	}
	return add(key, normalized)
}

// handleMultiValue is ConfigCommand::handleMultiValue.
func (c *configCmd) handleMultiValue(key string, h valueHandler, values []string, add func(string, phpjson.Value) error) error {
	list := listValues(values)
	if ok, msg := h.validate(list); !ok {
		suffix := ""
		if msg != "" {
			suffix = " (" + msg + ")"
		}
		enc, _ := phpjson.Encode(list, 0)
		return configErr(fmt.Sprintf(`%s is an invalid value`+suffix, enc))
	}
	return add(key, h.normalize(list))
}

// parseJSONArg is JsonFile::parseJson($value) without a file name.
func parseJSONArg(s string) (phpjson.Value, error) {
	return jsonfile.Parse([]byte(s), "")
}

// isList is array_is_list for decoded values.
func isList(v phpjson.Value) bool {
	switch t := v.(type) {
	case []phpjson.Value:
		return true
	case *phpjson.Map:
		return t.IsList()
	}
	return false
}

// mergeArrays implements the --merge combination: array_merge for two
// lists, $value + $current for two maps; mixed returns ok=false.
func mergeArrays(current, value phpjson.Value) (phpjson.Value, bool) {
	cl, vl := isList(current), isList(value)
	if cl && vl {
		var out []phpjson.Value
		phpjson.Each(current, func(_ string, e phpjson.Value) { out = append(out, e) })
		phpjson.Each(value, func(_ string, e phpjson.Value) { out = append(out, e) })
		if out == nil {
			out = []phpjson.Value{}
		}
		return out, true
	}
	if !cl && !vl {
		return arrayUnion(value, current), true
	}
	return nil, false
}

// arrayUnion is PHP's $a + $b.
func arrayUnion(a, b phpjson.Value) phpjson.Value {
	out := phpjson.NewMap()
	phpjson.Each(a, func(k string, e phpjson.Value) { out.Set(k, e) })
	phpjson.Each(b, func(k string, e phpjson.Value) {
		if !out.Has(k) {
			out.Set(k, e)
		}
	})
	return out
}

var (
	rePolicyList        = regexp.MustCompile(`^policy\.([^.]+)$`)
	rePolicyIgnore      = regexp.MustCompile(`^policy\.([^.]+)\.ignore$`)
	rePolicyCustomField = regexp.MustCompile(`^policy\.([^.]+)\.(block|block-scope|audit)$`)
	rePolicySources     = regexp.MustCompile(`^policy\.([^.]+)\.sources$`)
	rePreferredInstall  = regexp.MustCompile(`^preferred-install\.(.+)`)
	reAllowPlugins      = regexp.MustCompile(`^allow-plugins\.([a-zA-Z0-9/*-]+)`)
	reRepoSet           = regexp.MustCompile(`^repos?(?:itories)?\.(.+)`)
	reExtra             = regexp.MustCompile(`^extra\.(.+)`)
	reSuggest           = regexp.MustCompile(`^suggest\.(.+)`)
	rePlatform          = regexp.MustCompile(`^platform\.(.+)`)
	reAuthKey           = regexp.MustCompile(`^(bitbucket-oauth|github-oauth|gitlab-oauth|gitlab-token|http-basic|custom-headers|bearer|forgejo-token)\.(.+)`)
	reScripts           = regexp.MustCompile(`^scripts\.(.+)`)
	reHeader            = regexp.MustCompile(`^[^:]+:\s*.+$`)
)

func (c *configCmd) setValue(settingKey string, values []string) error {
	in := c.in
	unset := in.Bool("unset")
	src := c.configSource
	uniqueCfg := uniqueConfigValues()
	multiCfg := multiConfigValues()

	// unsetting any audit/policy.* key resolves to the same removeConfigSetting call
	if unset && (settingKey == "audit" || settingKey == "policy" || strings.HasPrefix(settingKey, "policy.")) {
		return src.RemoveConfigSetting(settingKey)
	}

	nonCustom := append(append([]string(nil), policy.BuiltinListNames...), policy.NonListKeys...)
	policyJSONMergeKeys := []string{"policy.advisories.ignore-id"}
	for _, n := range policy.BuiltinListNames {
		policyJSONMergeKeys = append(policyJSONMergeKeys, "policy."+n+".ignore")
	}
	customIgnore := ""
	isCustomPolicyIgnore := false
	if m := rePolicyIgnore.FindStringSubmatch(settingKey); m != nil && !inStringList(m[1], nonCustom) {
		isCustomPolicyIgnore, customIgnore = true, m[1]
	}
	if inStringList(settingKey, policyJSONMergeKeys) || isCustomPolicyIgnore {
		if isCustomPolicyIgnore {
			if e := policy.GetFutureReservedListNameError(customIgnore); e != "" {
				return configErr("Invalid dependency policy name: " + e)
			}
		}
		var value phpjson.Value = listValues(values)
		if in.Bool("json") {
			v, err := parseJSONArg(values[0])
			if err != nil {
				return err
			}
			if !phpjson.IsArray(v) {
				return configErr("Expected an array or object for " + settingKey)
			}
			value = v
		}
		if in.Bool("merge") {
			current, err := readJSON(c.configFile)
			if err != nil {
				return err
			}
			cur, _ := phpIsset(current, "config")
			for _, bit := range strings.Split(settingKey, ".") {
				if !phpjson.IsArray(cur) {
					cur = nil
					break
				}
				v, ok := phpIsset(cur, bit)
				if !ok {
					cur = nil
					break
				}
				cur = v
			}
			if cur != nil && phpjson.IsArray(cur) && phpjson.IsArray(value) {
				merged, ok := mergeArrays(cur, value)
				if !ok {
					return configErr("Cannot merge array and object for " + settingKey)
				}
				value = merged
			}
		}
		return src.AddConfigSetting(settingKey, value)
	}

	if settingKey == "policy.ignore-unreachable" {
		scopesErr := configErr("valid values for " + settingKey + " include: " + strings.Join(policy.IgnoreUnreachableScopes, ", "))
		if in.Bool("json") {
			v, err := parseJSONArg(values[0])
			if err != nil {
				return err
			}
			if !phpjson.IsArray(v) {
				return configErr("Expected a boolean or array for " + settingKey)
			}
			bad := false
			phpjson.Each(v, func(_ string, e phpjson.Value) {
				if !inValues(e, policy.IgnoreUnreachableScopes...) {
					bad = true
				}
			})
			if bad {
				return scopesErr
			}
			return src.AddConfigSetting(settingKey, v)
		}
		if len(values) > 0 && inStringList(values[0], policy.IgnoreUnreachableScopes) {
			for _, v := range values {
				if !inStringList(v, policy.IgnoreUnreachableScopes) {
					return scopesErr
				}
			}
			return src.AddConfigSetting(settingKey, listValues(values))
		}
	}

	if m := rePolicyList.FindStringSubmatch(settingKey); m != nil && !inStringList(m[1], policy.NonListKeys) {
		if e := policy.GetFutureReservedListNameError(m[1]); e != "" {
			return configErr("Invalid dependency policy name: " + e)
		}
		if ok, _ := boolValidator(values[0]); !ok {
			return configErr(fmt.Sprintf(`"%s" is an invalid value for %s, expected a boolean`, values[0], settingKey))
		}
		return src.AddConfigSetting(settingKey, boolNormalizer(values[0]))
	}

	if m := rePolicyCustomField.FindStringSubmatch(settingKey); m != nil && !inStringList(m[1], nonCustom) {
		if e := policy.GetFutureReservedListNameError(m[1]); e != "" {
			return configErr("Invalid dependency policy name: " + e)
		}
		switch m[2] {
		case "block":
			if ok, _ := boolValidator(values[0]); !ok {
				return configErr(fmt.Sprintf(`"%s" is an invalid value for %s, expected a boolean`, values[0], settingKey))
			}
			return src.AddConfigSetting(settingKey, boolNormalizer(values[0]))
		case "block-scope":
			if !inStringList(values[0], policy.AvailableBlockScopes) {
				return configErr(fmt.Sprintf(`"%s" is an invalid value for %s, must be one of: all, update, install`, values[0], settingKey))
			}
			return src.AddConfigSetting(settingKey, values[0])
		default:
			if !inStringList(values[0], []string{policy.AuditIgnore, policy.AuditReport, policy.AuditFail}) {
				return configErr(fmt.Sprintf(`"%s" is an invalid value for %s, must be one of: ignore, report, fail`, values[0], settingKey))
			}
			return src.AddConfigSetting(settingKey, values[0])
		}
	}

	if m := rePolicySources.FindStringSubmatch(settingKey); m != nil {
		return configErr(fmt.Sprintf("Setting dependency policy sources is not supported by `composer config`. Use `composer policy add-source %s url <https-url>` instead.", m[1]))
	}

	_, isUnique := uniqueCfg[settingKey]
	_, isMulti := multiCfg[settingKey]
	if unset && (isUnique || isMulti) {
		if settingKey == "disable-tls" {
			cur, err := c.cfg.GetFlags("disable-tls", 0)
			if err != nil {
				return err
			}
			if phpjson.Truthy(cur) {
				c.io.WriteError("<info>You are now running Composer with SSL/TLS protection enabled.</info>")
			}
		}
		return src.RemoveConfigSetting(settingKey)
	}
	if isUnique {
		return c.handleSingleValue(settingKey, uniqueCfg[settingKey], values, src.AddConfigSetting)
	}
	if isMulti {
		return c.handleMultiValue(settingKey, multiCfg[settingKey], values, src.AddConfigSetting)
	}

	// handle preferred-install per-package config
	if rePreferredInstall.MatchString(settingKey) {
		if unset {
			return src.RemoveConfigSetting(settingKey)
		}
		if !inStringList(values[0], []string{"auto", "source", "dist"}) {
			return configErr("Invalid value for " + settingKey + ". Should be one of: auto, source, or dist")
		}
		return src.AddConfigSetting(settingKey, values[0])
	}

	// handle allow-plugins config setting elements true or false to add/remove
	if reAllowPlugins.MatchString(settingKey) {
		if unset {
			return src.RemoveConfigSetting(settingKey)
		}
		if ok, _ := boolValidator(values[0]); !ok {
			return configErr(fmt.Sprintf(`"%s" is an invalid value`, values[0]))
		}
		return src.AddConfigSetting(settingKey, boolNormalizer(values[0]))
	}

	// handle properties
	uProps := uniqueProps()
	mProps := multiProps()
	_, isUProp := uProps[settingKey]
	_, isMProp := mProps[settingKey]
	if in.Bool("global") && (isUProp || isMProp || strings.HasPrefix(settingKey, "extra.")) {
		return configErr("The " + settingKey + " property can not be set in the global config.json file. Use `composer global config` to apply changes to the global composer.json")
	}
	if unset && (isUProp || isMProp) {
		return src.RemoveProperty(settingKey)
	}
	if isUProp {
		return c.handleSingleValue(settingKey, uProps[settingKey], values, src.AddProperty)
	}
	if isMProp {
		return c.handleMultiValue(settingKey, mProps[settingKey], values, src.AddProperty)
	}

	// handle repositories
	if m := reRepoSet.FindStringSubmatch(settingKey); m != nil {
		name := m[1]
		if unset {
			return src.RemoveRepository(name)
		}
		appendRepo := in.Bool("append")
		if len(values) == 2 {
			repo := phpjson.NewMap()
			repo.Set("type", values[0])
			repo.Set("url", values[1])
			return src.AddRepository(name, repo, appendRepo)
		}
		if len(values) == 1 {
			value := strings.ToLower(values[0])
			if ok, _ := boolValidator(value); ok {
				if boolNormalizer(value) == false {
					return src.AddRepository(name, false, appendRepo)
				}
			} else {
				v, err := parseJSONArg(values[0])
				if err != nil {
					return err
				}
				return src.AddRepository(name, v, appendRepo)
			}
		}
		return configErr("You must pass the type and a url. Example: php composer.phar config repositories.foo vcs https://bar.com")
	}

	// handle extra
	if reExtra.MatchString(settingKey) {
		if unset {
			return src.RemoveProperty(settingKey)
		}
		var value phpjson.Value = values[0]
		if in.Bool("json") {
			v, err := parseJSONArg(values[0])
			if err != nil {
				return err
			}
			value = v
			if in.Bool("merge") {
				current, err := readJSON(c.configFile)
				if err != nil {
					return err
				}
				cur := current
				for _, bit := range strings.Split(settingKey, ".") {
					v, _ := phpIsset(cur, bit)
					cur = v
				}
				if phpjson.IsArray(cur) && phpjson.IsArray(value) {
					if isList(cur) && isList(value) {
						value, _ = mergeArrays(cur, value)
					} else {
						value = arrayUnion(value, cur)
					}
				}
			}
		}
		return src.AddProperty(settingKey, value)
	}

	// handle suggest
	if reSuggest.MatchString(settingKey) {
		if unset {
			return src.RemoveProperty(settingKey)
		}
		return src.AddProperty(settingKey, strings.Join(values, " "))
	}

	// handle unsetting extra/suggest
	if (settingKey == "suggest" || settingKey == "extra") && unset {
		return src.RemoveProperty(settingKey)
	}

	// handle platform
	if rePlatform.MatchString(settingKey) {
		if unset {
			return src.RemoveConfigSetting(settingKey)
		}
		var v phpjson.Value = values[0]
		if values[0] == "false" {
			v = false
		}
		return src.AddConfigSetting(settingKey, v)
	}

	// handle unsetting platform
	if settingKey == "platform" && unset {
		return src.RemoveConfigSetting(settingKey)
	}

	// handle audit.ignore and audit.ignore-abandoned with --merge support
	if settingKey == "audit.ignore" || settingKey == "audit.ignore-abandoned" {
		if unset {
			return src.RemoveConfigSetting(settingKey)
		}
		var value phpjson.Value = listValues(values)
		if in.Bool("json") {
			v, err := parseJSONArg(values[0])
			if err != nil {
				return err
			}
			if !phpjson.IsArray(v) {
				return configErr("Expected an array or object for " + settingKey)
			}
			value = v
		}
		if in.Bool("merge") {
			current, err := readJSON(c.configFile)
			if err != nil {
				return err
			}
			cur := dig2(current, "config", "audit", strings.TrimPrefix(settingKey, "audit."))
			if cur != nil && phpjson.IsArray(cur) && phpjson.IsArray(value) {
				merged, ok := mergeArrays(cur, value)
				if !ok {
					return configErr("Cannot merge array and object for " + settingKey)
				}
				value = merged
			}
		}
		return src.AddConfigSetting(settingKey, value)
	}

	// handle auth
	if m := reAuthKey.FindStringSubmatch(settingKey); m != nil {
		key := m[1] + "." + m[2]
		auth := c.authSource
		if unset {
			if err := auth.RemoveConfigSetting(key); err != nil {
				return err
			}
			return src.RemoveConfigSetting(key)
		}
		pair := func(a, b string) *phpjson.Map {
			p := phpjson.NewMap()
			p.Set(a, values[0])
			p.Set(b, values[1])
			return p
		}
		switch {
		case m[1] == "bitbucket-oauth":
			if len(values) != 2 {
				return configErr("Expected two arguments (consumer-key, consumer-secret), got " + strconv.Itoa(len(values)))
			}
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, pair("consumer-key", "consumer-secret"))
		case m[1] == "gitlab-token" && len(values) == 2:
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, pair("username", "token"))
		case inStringList(m[1], []string{"github-oauth", "gitlab-oauth", "gitlab-token", "bearer"}):
			if len(values) != 1 {
				return configErr("Too many arguments, expected only one token")
			}
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, values[0])
		case m[1] == "http-basic":
			if len(values) != 2 {
				return configErr("Expected two arguments (username, password), got " + strconv.Itoa(len(values)))
			}
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, pair("username", "password"))
		case m[1] == "custom-headers":
			if len(values) == 0 {
				return configErr("Expected at least one argument (header), got none")
			}
			for _, h := range values {
				if !reHeader.MatchString(h) {
					return configErr(`Header "` + h + `" is not in "Header-Name: Header-Value" format`)
				}
			}
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, listValues(values))
		case m[1] == "forgejo-token":
			if len(values) != 2 {
				return configErr("Expected two arguments (username, access token), got " + strconv.Itoa(len(values)))
			}
			if err := src.RemoveConfigSetting(key); err != nil {
				return err
			}
			return auth.AddConfigSetting(key, pair("username", "token"))
		}
		return nil
	}

	// handle script
	if reScripts.MatchString(settingKey) {
		if unset {
			return src.RemoveProperty(settingKey)
		}
		if len(values) > 1 {
			return src.AddProperty(settingKey, listValues(values))
		}
		return src.AddProperty(settingKey, values[0])
	}

	// handle unsetting other top level properties
	if unset {
		return src.RemoveProperty(settingKey)
	}

	return configErr("Setting " + settingKey + " does not exist or is not supported by this command")
}

// dig2 is $data['a']['b']['c'] ?? null.
func dig2(v phpjson.Value, keys ...string) phpjson.Value {
	for _, k := range keys {
		next, ok := phpIsset(v, k)
		if !ok {
			return nil
		}
		v = next
	}
	return v
}
