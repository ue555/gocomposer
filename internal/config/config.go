// Package config ports Composer\Config and the config parts of Composer\Factory.
package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
)

// Config is Composer\Config.
type Config struct {
	config       *phpjson.Map
	repositories *phpjson.Map // ordered name => repository definition
	BaseDir      string
	UseEnv       bool
	Getenv       func(string) (string, bool)
	// warn reports schema warnings (Factory::validateJsonSchema with an IO).
	warn func(string)
	// sources is Config::$sourceOfConfigValue.
	sources map[string]string
	// debug writes the IOInterface::DEBUG messages of Factory (config files
	// loaded, COMPOSER_AUTH); nil without an IO.
	debug func(string)
}

func defaultConfig() *phpjson.Map {
	m := phpjson.NewMap()
	set := func(k string, v phpjson.Value) { m.Set(k, v) }
	set("process-timeout", int64(300))
	set("use-include-path", false)
	set("allow-plugins", []phpjson.Value{})
	set("use-parent-dir", "prompt")
	set("preferred-install", "dist")
	audit := phpjson.NewMap()
	audit.Set("ignore", []phpjson.Value{})
	audit.Set("abandoned", "fail")
	set("audit", audit)
	set("policy", true)
	set("notify-on-install", true)
	set("github-protocols", []phpjson.Value{"https", "ssh", "git"})
	set("gitlab-protocol", nil)
	set("vendor-dir", "vendor")
	set("bin-dir", "{$vendor-dir}/bin")
	set("cache-dir", "{$home}/cache")
	set("data-dir", "{$home}")
	set("cache-files-dir", "{$cache-dir}/files")
	set("cache-repo-dir", "{$cache-dir}/repo")
	set("cache-vcs-dir", "{$cache-dir}/vcs")
	set("cache-ttl", int64(15552000))
	set("cache-files-ttl", nil)
	set("cache-files-maxsize", "300MiB")
	set("cache-read-only", false)
	set("bin-compat", "auto")
	set("discard-changes", false)
	set("autoloader-suffix", nil)
	set("sort-packages", false)
	set("optimize-autoloader", false)
	set("strict-psr-autoloader", false)
	set("classmap-authoritative", false)
	set("apcu-autoloader", false)
	set("prepend-autoloader", true)
	set("update-with-minimal-changes", false)
	set("github-domains", []phpjson.Value{"github.com"})
	set("bitbucket-expose-hostname", true)
	set("disable-tls", false)
	set("secure-http", true)
	set("secure-svn-domains", []phpjson.Value{})
	set("cafile", nil)
	set("capath", nil)
	set("github-expose-hostname", true)
	set("gitlab-domains", []phpjson.Value{"gitlab.com"})
	set("store-auths", "prompt")
	set("platform", []phpjson.Value{})
	set("archive-format", "tar")
	set("archive-dir", ".")
	set("htaccess-protect", true)
	set("use-github-api", true)
	set("lock", true)
	set("platform-check", "php-only")
	set("bitbucket-oauth", []phpjson.Value{})
	set("github-oauth", []phpjson.Value{})
	set("gitlab-oauth", []phpjson.Value{})
	set("gitlab-token", []phpjson.Value{})
	set("http-basic", []phpjson.Value{})
	set("bearer", []phpjson.Value{})
	set("custom-headers", []phpjson.Value{})
	set("bump-after-update", false)
	set("allow-missing-requirements", false)
	set("client-certificate", []phpjson.Value{})
	set("forgejo-domains", []phpjson.Value{"codeberg.org"})
	set("forgejo-token", []phpjson.Value{})
	set("source-fallback", false)
	return m
}

// New is new Config($useEnvironment, $baseDir).
func New(useEnv bool, baseDir string) *Config {
	repos := phpjson.NewMap()
	pk := phpjson.NewMap()
	pk.Set("type", "composer")
	pk.Set("url", "https://repo.packagist.org")
	repos.Set("packagist.org", pk)
	c := &Config{config: defaultConfig(), repositories: repos, BaseDir: baseDir, UseEnv: useEnv, Getenv: os.LookupEnv, sources: map[string]string{}}
	for _, k := range c.config.Keys() {
		v, _ := c.config.Get(k)
		c.setSource(v, k, SourceDefault)
	}
	for _, k := range c.repositories.Keys() {
		v, _ := c.repositories.Get(k)
		c.setSource(v, "repositories."+k, SourceDefault)
	}
	return c
}

func (c *Config) env(name string) (string, bool) {
	if !c.UseEnv {
		return "", false
	}
	return c.Getenv(name)
}

// envValue mirrors Platform::getEnv returning false for unset variables.
func (c *Config) envValue(name string) (string, bool) {
	v, ok := c.env(name)
	if !ok {
		return "", false
	}
	return v, true
}

func arrayMerge(a, b phpjson.Value) phpjson.Value {
	// array_merge(): string keys override in place, int keys append
	out := phpjson.NewMap()
	n := 0
	add := func(v phpjson.Value) {
		phpjson.Each(v, func(k string, e phpjson.Value) {
			if _, isInt := phpjson.IntKey(k); isInt {
				out.Set(strconv.Itoa(n), e)
				n++
			} else {
				out.Set(k, e)
			}
		})
	}
	add(a)
	add(b)
	return out
}

var rePackagistURL = regexp.MustCompile(`^https?://(?:[a-z0-9-.]+\.)?packagist.org(/|$)`)

// Merge is Config::merge with the default (unknown) source.
func (c *Config) Merge(data *phpjson.Map) { c.MergeSource(data, SourceUnknown) }

// MergeSource is Config::merge($config, $source).
func (c *Config) MergeSource(data *phpjson.Map, source string) {
	if cfg, ok := data.Get("config"); ok && phpjson.Truthy(cfg) && phpjson.IsArray(cfg) {
		phpjson.Each(cfg, func(key string, val phpjson.Value) {
			cur, exists := c.config.Get(key)
			switch {
			case containsKey([]string{"bitbucket-oauth", "github-oauth", "gitlab-oauth", "gitlab-token", "http-basic", "bearer", "custom-headers", "client-certificate", "forgejo-token"}, key) && exists && cur != nil:
				c.config.Set(key, arrayMerge(cur, val))
				c.setSource(val, key, source)
			case key == "allow-plugins" && exists && phpjson.IsArray(cur) && phpjson.IsArray(val):
				c.config.Set(key, arrayMerge(arrayMerge(val, cur), val))
				c.setSource(val, key, source)
			case (key == "gitlab-domains" || key == "github-domains") && exists && cur != nil:
				// array_unique() keeps the keys of the first occurrences
				merged := arrayMerge(cur, val)
				seen := map[string]bool{}
				uniq := phpjson.NewMap()
				phpjson.Each(merged, func(k string, e phpjson.Value) {
					s, _ := phpjson.ScalarString(e)
					if !seen[s] {
						seen[s] = true
						uniq.Set(k, e)
					}
				})
				c.config.Set(key, uniq)
				c.setSource(val, key, source)
			case key == "preferred-install" && exists && cur != nil:
				if phpjson.IsArray(val) || phpjson.IsArray(cur) {
					if s, ok := val.(string); ok {
						m := phpjson.NewMap()
						m.Set("*", s)
						val = m
					}
					if s, ok := cur.(string); ok {
						m := phpjson.NewMap()
						m.Set("*", s)
						cur = m
						c.sources[key+"*"] = source
					}
					merged := phpjson.AsMap(arrayMerge(cur, val))
					if w, ok := merged.Get("*"); ok {
						merged.Delete("*")
						merged.Set("*", w)
					}
					c.config.Set(key, merged)
					c.setSource(val, key, source)
				} else {
					c.config.Set(key, val)
					c.setSource(val, key, source)
				}
			case key == "audit":
				curMap := phpjson.ToMap(cur)
				currentIgnores, _ := curMap.Get("ignore")
				merged := phpjson.AsMap(arrayMerge(curMap, val))
				var newIgnore phpjson.Value = []phpjson.Value{}
				if vm := phpjson.ToMap(val); vm != nil {
					if ig, ok := vm.Get("ignore"); ok {
						newIgnore = ig
					}
				}
				merged.Set("ignore", arrayMerge(currentIgnores, newIgnore))
				c.config.Set(key, merged)
				c.setSource(val, key, source)
			case key == "policy":
				if val == true {
					val = phpjson.NewMap()
				}
				if val == false {
					c.config.Set(key, false)
				} else if phpjson.IsArray(val) {
					current := phpjson.NewMap()
					if cm := phpjson.ToMap(cur); cm != nil && phpjson.IsArray(cur) {
						current = cm.Clone()
					}
					phpjson.Each(val, func(listName string, listConfig phpjson.Value) {
						if listName == "ignore-unreachable" {
							current.Set(listName, listConfig)
							return
						}
						if listConfig == true {
							listConfig = phpjson.NewMap()
						}
						existing, _ := current.Get(listName)
						if existing == true {
							existing = phpjson.NewMap()
						}
						switch {
						case listConfig == false:
							current.Set(listName, false)
						case existing == nil || existing == false:
							current.Set(listName, listConfig)
						case phpjson.IsArray(existing) && phpjson.IsArray(listConfig):
							merged := phpjson.AsMap(arrayMerge(existing, listConfig))
							for _, inner := range []string{"ignore", "ignore-id", "ignore-severity", "ignore-source"} {
								ei, _ := phpjson.ToMap(existing).Get(inner)
								ii, _ := phpjson.ToMap(listConfig).Get(inner)
								if phpjson.IsArray(ei) && phpjson.IsArray(ii) {
									merged.Set(inner, arrayMerge(ei, ii))
								}
							}
							current.Set(listName, merged)
						default:
							current.Set(listName, listConfig)
						}
					})
					c.config.Set(key, current)
				}
				c.setSource(val, key, source)
			default:
				c.config.Set(key, val)
				c.setSource(val, key, source)
			}
		})
	}
	if reposV, ok := data.Get("repositories"); ok && phpjson.Truthy(reposV) && phpjson.IsArray(reposV) {
		// array_reverse($this->repositories, true)
		keys := c.repositories.Keys()
		rev := phpjson.NewMap()
		for i := len(keys) - 1; i >= 0; i-- {
			v, _ := c.repositories.Get(keys[i])
			rev.Set(keys[i], v)
		}
		c.repositories = rev
		var newKeys []string
		newVals := map[string]phpjson.Value{}
		phpjson.Each(reposV, func(k string, v phpjson.Value) {
			newKeys = append(newKeys, k)
			newVals[k] = v
		})
		for i := len(newKeys) - 1; i >= 0; i-- {
			name := newKeys[i]
			repository := newVals[name]
			if repository == false {
				c.disableRepoByName(name)
				continue
			}
			if m := phpjson.ToMap(repository); m != nil && m.Len() == 1 {
				k := m.Keys()[0]
				if v, _ := m.Get(k); v == false {
					c.disableRepoByName(k)
					continue
				}
			}
			if m := phpjson.ToMap(repository); m != nil && m.Isset("type") && m.Isset("url") {
				if m.Str("type") == "composer" && rePackagistURL.MatchString(m.Str("url")) {
					c.disableRepoByName("packagist.org")
				}
			}
			if _, isInt := phpjson.IntKey(name); isInt {
				if !c.repositories.Has(name) {
					c.repositories.Set(name, repository)
				} else {
					// $this->repositories[] = $repository: next integer key
					next := 0
					for _, k := range c.repositories.Keys() {
						if n, ok := phpjson.IntKey(k); ok && int(n) >= next {
							next = int(n) + 1
						}
					}
					c.repositories.Set(strconv.Itoa(next), repository)
				}
				// 'repositories.' . array_search($repository, $this->repositories, true)
				found := ""
				for _, k := range c.repositories.Keys() {
					if v, _ := c.repositories.Get(k); strictEqual(v, repository) {
						found = k
						break
					}
				}
				c.setSource(repository, "repositories."+found, source)
			} else if name == "packagist" {
				c.repositories.Set("packagist.org", repository)
				c.setSource(repository, "repositories.packagist.org", source)
			} else {
				c.repositories.Set(name, repository)
				c.setSource(repository, "repositories."+name, source)
			}
		}
		keys = c.repositories.Keys()
		rev = phpjson.NewMap()
		for i := len(keys) - 1; i >= 0; i-- {
			v, _ := c.repositories.Get(keys[i])
			rev.Set(keys[i], v)
		}
		c.repositories = rev
	}
}

func containsKey(list []string, k string) bool {
	for _, e := range list {
		if e == k {
			return true
		}
	}
	return false
}

func (c *Config) disableRepoByName(name string) {
	if c.repositories.Has(name) {
		c.repositories.Delete(name)
	} else if name == "packagist" {
		c.repositories.Delete("packagist.org")
	}
}

// Repositories is getRepositories().
func (c *Config) Repositories() *phpjson.Map { return c.repositories }

// Raw returns the raw config map.
func (c *Config) Raw() *phpjson.Map { return c.config }

var reVar = regexp.MustCompile(`\{\$(.+)\}`)

// process is Config::process: {$key} references are resolved with the same
// flags (RELATIVE_PATHS) as the value being read.
func (c *Config) process(v phpjson.Value, relative bool) phpjson.Value {
	s, ok := v.(string)
	if !ok {
		return v
	}
	return reVar.ReplaceAllStringFunc(s, func(m string) string {
		key := reVar.FindStringSubmatch(m)[1]
		r := c.get(key, relative)
		str, _ := phpjson.ScalarString(r)
		return str
	})
}

var reAbsolute = regexp.MustCompile(`(?i)^(?:/|[a-z]:|[a-z0-9.]+://|\\\\)`)

func (c *Config) realpath(p string) string {
	if reAbsolute.MatchString(p) {
		return p
	}
	if c.BaseDir != "" {
		return c.BaseDir + "/" + p
	}
	return p
}

// ExpandPath is Platform::expandPath: a leading "~/" becomes the user
// directory and a leading $VAR or %VAR% is replaced by its value.
func ExpandPath(p string) string {
	if len(p) >= 2 && p[0] == '~' && (p[1] == '/' || p[1] == '\\') {
		home, ok := os.LookupEnv("HOME")
		if !ok {
			home, _ = os.UserHomeDir()
		}
		return home + p[1:]
	}
	m := reEnvVar.FindStringSubmatch(p)
	if m == nil {
		return p
	}
	name := m[1]
	if name == "" {
		name = m[2]
	}
	v, _ := os.LookupEnv(name)
	return v + m[3]
}

// #^(\$|(?P<percent>%))(?P<var>\w++)(?(percent)%)(?P<path>.*)#
var reEnvVar = regexp.MustCompile(`(?s)^(?:\$([A-Za-z0-9_]+)|%([A-Za-z0-9_]+)%)(.*)$`)

// Get is Config::get for every key; relativePaths mirrors RELATIVE_PATHS.
func (c *Config) Get(key string) phpjson.Value { return c.get(key, false) }

// GetRelative is get($key, Config::RELATIVE_PATHS).
func (c *Config) GetRelative(key string) phpjson.Value { return c.get(key, true) }

func (c *Config) get(key string, relative bool) phpjson.Value {
	cur, _ := c.config.Get(key)
	switch key {
	case "vendor-dir", "bin-dir", "process-timeout", "data-dir", "cache-dir", "cache-files-dir", "cache-repo-dir", "cache-vcs-dir", "cafile", "capath":
		envName := "COMPOSER_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		envVal, hasEnv := c.envValue(envName)
		if key == "process-timeout" {
			n := int64(0)
			if hasEnv {
				n, _ = strconv.ParseInt(strings.TrimSpace(envVal), 10, 64)
			} else {
				switch t := cur.(type) {
				case int64:
					n = t
				case string:
					n, _ = strconv.ParseInt(t, 10, 64)
				case float64:
					n = int64(t)
				}
			}
			if n < 0 {
				n = 0
			}
			return n
		}
		var raw phpjson.Value = cur
		if hasEnv {
			raw = envVal
		}
		if raw == nil {
			return nil
		}
		processed := c.process(raw, relative)
		s, _ := phpjson.ScalarString(processed)
		s = strings.TrimRight(s, "/\\")
		s = ExpandPath(s)
		if !strings.HasSuffix(key, "-dir") {
			if s == "" && cur == nil && !hasEnv {
				return nil
			}
			return s
		}
		if relative {
			return s
		}
		return c.realpath(s)
	case "cache-read-only", "htaccess-protect":
		envName := "COMPOSER_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		if v, ok := c.envValue(envName); ok {
			return v != "false" && phpjson.Truthy(v)
		}
		return cur != "false" && phpjson.Truthy(cur)
	case "disable-tls", "secure-http", "use-github-api", "lock", "source-fallback", "strict-psr-autoloader":
		if key == "secure-http" && c.Get("disable-tls") == true {
			return false
		}
		return cur != "false" && phpjson.Truthy(cur)
	case "home":
		s, _ := phpjson.ScalarString(c.process(ExpandPath(asString(cur)), relative))
		return strings.TrimRight(s, "/\\")
	case "autoloader-suffix":
		if cur == "" {
			return nil
		}
		return c.process(cur, relative)
	case "policy":
		if v, ok := c.envValue("COMPOSER_POLICY"); ok {
			if b, isBool := boolEnv(v); isBool {
				if !b {
					return false
				}
				if cur == false {
					return true
				}
			}
		}
		return cur
	}
	if cur == nil {
		return nil
	}
	return c.process(cur, relative)
}

func asString(v phpjson.Value) string { s, _ := phpjson.ScalarString(v); return s }

// boolEnv mirrors Platform::getBoolEnv.
func boolEnv(v string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off", "":
		return false, true
	}
	return false, false
}

// GetString returns a string config value.
func (c *Config) GetString(key string) string { return asString(c.Get(key)) }

// GetBool returns a boolean-ish config value.
func (c *Config) GetBool(key string) bool {
	v := c.Get(key)
	if v == "false" {
		return false
	}
	return phpjson.Truthy(v)
}

// Set overrides a raw config value (used by --option style overrides).
func (c *Config) Set(key string, v phpjson.Value) { c.config.Set(key, v) }

// HomeDir is Factory::getHomeDir.
func HomeDir(getenv func(string) (string, bool)) string {
	if h, ok := getenv("COMPOSER_HOME"); ok && h != "" {
		return h
	}
	userDir := userDir(getenv)
	var dirs []string
	if useXdg(getenv) {
		xdg, _ := getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			xdg = userDir + "/.config"
		}
		dirs = append(dirs, xdg+"/composer")
	}
	dirs = append(dirs, userDir+"/.composer")
	for _, d := range dirs {
		if st, err := os.Stat(d); err == nil && st.IsDir() {
			return d
		}
	}
	return dirs[0]
}

func userDir(getenv func(string) (string, bool)) string {
	home, _ := getenv("HOME")
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	return strings.TrimRight(strings.ReplaceAll(home, "\\", "/"), "/")
}

func useXdg(getenv func(string) (string, bool)) bool {
	for _, k := range []string{"XDG_CONFIG_HOME", "XDG_CACHE_HOME", "XDG_DATA_HOME", "XDG_RUNTIME_DIR", "XDG_SESSION_ID"} {
		if v, ok := getenv(k); ok && v != "" {
			return true
		}
	}
	if st, err := os.Stat("/etc/xdg"); err == nil && st.IsDir() {
		return true
	}
	return false
}

// CacheDir is Factory::getCacheDir (without migrating old cache dirs).
func CacheDir(home string, getenv func(string) (string, bool)) string {
	if d, ok := getenv("COMPOSER_CACHE_DIR"); ok && d != "" {
		return d
	}
	if h, ok := getenv("COMPOSER_HOME"); ok && h != "" {
		return h + "/cache"
	}
	ud := userDir(getenv)
	if runtime.GOOS == "darwin" {
		return ud + "/Library/Caches/composer"
	}
	if home == ud+"/.composer" {
		if st, err := os.Stat(home + "/cache"); err == nil && st.IsDir() {
			return home + "/cache"
		}
	}
	if useXdg(getenv) {
		x, _ := getenv("XDG_CACHE_HOME")
		if x == "" {
			x = ud + "/.cache"
		}
		return x + "/composer"
	}
	return home + "/cache"
}

// DataDir is Factory::getDataDir.
func DataDir(home string, getenv func(string) (string, bool)) string {
	if h, ok := getenv("COMPOSER_HOME"); ok && h != "" {
		return h
	}
	ud := userDir(getenv)
	if home != ud+"/.composer" && useXdg(getenv) {
		x, _ := getenv("XDG_DATA_HOME")
		if x == "" {
			x = ud + "/.local/share"
		}
		return x + "/composer"
	}
	return home
}

// ReadJSONFile reads and decodes a JSON object file.
func ReadJSONFile(path string) (*phpjson.Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v, err := jsonfile.Parse(data, path)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	return m, nil
}

// CreateConfig is Factory::createConfig: defaults, home dirs, global
// config.json and auth.json, COMPOSER_AUTH. It writes the .htaccess files
// Composer writes when htaccess-protect is enabled.
func CreateConfig(cwd string, warn func(string)) (*Config, error) {
	return CreateConfigDebug(cwd, warn, nil)
}

// CreateConfigDebug is Factory::createConfig($io): debug receives the
// messages Composer writes at IOInterface::DEBUG verbosity ("Loading config
// file ..."); it is kept for Factory::loadComposerAuthEnv.
func CreateConfigDebug(cwd string, warn func(string), debug func(string)) (*Config, error) {
	c := New(true, cwd)
	c.debug = debug
	home := HomeDir(c.Getenv)
	base := phpjson.NewMap()
	cfg := phpjson.NewMap()
	cfg.Set("home", home)
	cfg.Set("cache-dir", CacheDir(home, c.Getenv))
	cfg.Set("data-dir", DataDir(home, c.Getenv))
	base.Set("config", cfg)
	c.MergeSource(base, SourceDefault)
	globalFile := c.GetString("home") + "/config.json"
	if _, err := os.Stat(globalFile); err == nil {
		c.debugf("Loading config file " + globalFile)
		if err := ValidateSchemaWarn(warn, globalFile, nil, jsonschema.LaxSchema); err != nil {
			return nil, err
		}
		m, err := ReadJSONFile(globalFile)
		if err != nil {
			return nil, err
		}
		c.MergeSource(m, globalFile)
	}
	if c.GetBool("htaccess-protect") {
		for _, d := range []string{c.GetString("home"), c.GetString("cache-dir"), c.GetString("data-dir")} {
			if _, err := os.Stat(filepath.Join(d, ".htaccess")); err != nil {
				_ = os.MkdirAll(d, 0o777)
				_ = os.WriteFile(filepath.Join(d, ".htaccess"), []byte("Deny from all"), 0o666)
			}
		}
	}
	authFile := c.GetString("home") + "/auth.json"
	if _, err := os.Stat(authFile); err == nil {
		c.debugf("Loading config file " + authFile)
		if err := ValidateSchemaWarn(warn, authFile, nil, jsonschema.AuthSchema); err != nil {
			return nil, err
		}
		m, err := ReadJSONFile(authFile)
		if err != nil {
			return nil, err
		}
		wrap := phpjson.NewMap()
		wrap.Set("config", m)
		c.MergeSource(wrap, authFile)
	}
	c.warn = warn
	if err := c.LoadAuthEnv(); err != nil {
		return nil, err
	}
	return c, nil
}

// ValidateSchemaWarn is Factory::validateJsonSchema: schema errors of a
// config/auth file (or of data when non-nil) are reported as a warning, or
// returned as an error when there is no IO (warn == nil).
func ValidateSchemaWarn(warn func(string), source string, data phpjson.Value, mode int) error {
	if data == nil {
		raw, err := os.ReadFile(source)
		if err != nil {
			return &EnvError{`The file "` + source + `" is not readable.`}
		}
		v, err := jsonfile.Parse(raw, source)
		if err != nil {
			return err
		}
		data = v
	}
	msgs, err := jsonschema.ValidateComposerJSON(data, mode)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}
	msg := `"` + source + `" does not match the expected JSON schema, this may result in errors and should be resolved:` + "\n - " + strings.Join(msgs, "\n - ")
	if warn != nil {
		warn("<warning>" + msg + "</>")
		return nil
	}
	return &EnvError{msg}
}

// LoadAuthEnv merges COMPOSER_AUTH (Factory::loadComposerAuthEnv).
func (c *Config) LoadAuthEnv() error {
	v, ok := c.env("COMPOSER_AUTH")
	if !ok || v == "" {
		return nil
	}
	decoded, err := phpjson.Decode([]byte(v))
	if err != nil {
		return &EnvError{"COMPOSER_AUTH environment variable is malformed, should be a valid JSON object"}
	}
	c.debugf("Loading auth config from COMPOSER_AUTH")
	if err := ValidateSchemaWarn(c.warn, "COMPOSER_AUTH", decoded, jsonschema.AuthSchema); err != nil {
		return err
	}
	m := phpjson.ToMap(decoded)
	if m == nil {
		return &EnvError{"COMPOSER_AUTH environment variable is malformed, should be a valid JSON object"}
	}
	wrap := phpjson.NewMap()
	wrap.Set("config", m)
	c.MergeSource(wrap, "COMPOSER_AUTH")
	return nil
}

// EnvError reports malformed environment configuration.
type EnvError struct{ Msg string }

func (e *EnvError) Error() string { return e.Msg }
