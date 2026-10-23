package config

// Port of the value-reporting parts of Composer\Config (src/Composer/Config.php
// of composer/composer f83588392df43d4fd027d84d26de459946ef4e80: get() with
// flags, all(), raw(), getSourceOfValue(), setSourceOfConfigValue()).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/phpjson"
)

// Config::SOURCE_* and Config::RELATIVE_PATHS.
const (
	SourceDefault = "default"
	SourceCommand = "command"
	SourceUnknown = "unknown"
	RelativePaths = 1
)

// RuntimeError is a \RuntimeException thrown by Config::get.
type RuntimeError struct{ Msg string }

func (e *RuntimeError) Error() string { return e.Msg }

// setSource is Config::setSourceOfConfigValue.
func (c *Config) setSource(v phpjson.Value, path, source string) {
	if c.sources == nil {
		c.sources = map[string]string{}
	}
	c.sources[path] = source
	if phpjson.IsArray(v) {
		phpjson.Each(v, func(k string, e phpjson.Value) {
			c.setSource(e, path+"."+k, source)
		})
	}
}

// strictEqual is PHP's === for decoded JSON values.
func strictEqual(a, b phpjson.Value) bool {
	if phpjson.IsArray(a) || phpjson.IsArray(b) {
		if !phpjson.IsArray(a) || !phpjson.IsArray(b) || phpjson.ArrayLen(a) != phpjson.ArrayLen(b) {
			return false
		}
		var ak, bk []string
		var av, bv []phpjson.Value
		phpjson.Each(a, func(k string, v phpjson.Value) { ak = append(ak, k); av = append(av, v) })
		phpjson.Each(b, func(k string, v phpjson.Value) { bk = append(bk, k); bv = append(bv, v) })
		for i := range ak {
			if ak[i] != bk[i] || !strictEqual(av[i], bv[i]) {
				return false
			}
		}
		return true
	}
	switch x := a.(type) {
	case int:
		a = int64(x)
	}
	switch x := b.(type) {
	case int:
		b = int64(x)
	}
	return a == b
}

// SetBaseDir is Config::setBaseDir.
func (c *Config) SetBaseDir(dir string) { c.BaseDir = dir }

// Has is Config::has.
func (c *Config) Has(key string) bool { return c.config.Has(key) }

// RawAll is Config::raw(): repositories and the unprocessed config.
func (c *Config) RawAll() *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("repositories", c.repositories)
	m.Set("config", c.config)
	return m
}

// All is Config::all($flags).
func (c *Config) All(flags int) (*phpjson.Map, error) {
	all := phpjson.NewMap()
	all.Set("repositories", c.repositories)
	cfg := phpjson.NewMap()
	for _, k := range c.config.Keys() {
		v, err := c.GetFlags(k, flags)
		if err != nil {
			return nil, err
		}
		cfg.Set(k, v)
	}
	all.Set("config", cfg)
	return all, nil
}

// SourceOfValue is Config::getSourceOfValue.
func (c *Config) SourceOfValue(key string) (string, error) {
	if _, err := c.GetFlags(key, 0); err != nil {
		return "", err
	}
	if s, ok := c.sources[key]; ok {
		return s, nil
	}
	return SourceUnknown, nil
}

// PHPString is PHP's string conversion of a config value (string
// interpolation / preg_replace_callback return values).
func PHPString(v phpjson.Value) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "1"
		}
		return ""
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	case float64:
		return phpjson.FormatFloatSerialize(t)
	case *phpjson.Map, []phpjson.Value:
		return "Array"
	}
	return fmt.Sprint(v)
}

// PHPIntval is PHP's (int) cast.
func PHPIntval(v phpjson.Value) int64 {
	switch t := v.(type) {
	case nil:
		return 0
	case bool:
		if t {
			return 1
		}
		return 0
	case int64:
		return t
	case int:
		return int64(t)
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) {
			return 0
		}
		if t >= math.MaxInt64 || t <= math.MinInt64 {
			return 0
		}
		return int64(t)
	case string:
		return stringIntval(t)
	case *phpjson.Map:
		if t.Len() > 0 {
			return 1
		}
		return 0
	case []phpjson.Value:
		if len(t) > 0 {
			return 1
		}
		return 0
	}
	return 0
}

var reLeadingNumber = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?`)
var reLeadingInt = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?[0-9]+`)

// stringIntval is (int) "string" in PHP 8: the leading numeric part, with
// float notation (1e3, 1.5) truncated and integer overflow saturated.
func stringIntval(s string) int64 {
	num := reLeadingNumber.FindString(s)
	if num == "" {
		return 0
	}
	trimmed := strings.TrimLeft(num, " \t\n\r\v\f")
	if reLeadingInt.FindString(num) == num {
		n, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			if strings.HasPrefix(trimmed, "-") {
				return math.MinInt64
			}
			return math.MaxInt64
		}
		return n
	}
	f, _ := strconv.ParseFloat(trimmed, 64)
	return PHPIntval(f)
}

// PHPFloatval is (float) "string".
func PHPFloatval(s string) float64 {
	num := reLeadingNumber.FindString(s)
	if num == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(strings.TrimLeft(num, " \t\n\r\v\f"), 64)
	return f
}

func (c *Config) processFlags(v phpjson.Value, flags int) (phpjson.Value, error) {
	s, ok := v.(string)
	if !ok {
		return v, nil
	}
	var ferr error
	out := reVar.ReplaceAllStringFunc(s, func(m string) string {
		key := reVar.FindStringSubmatch(m)[1]
		r, err := c.GetFlags(key, flags)
		if err != nil && ferr == nil {
			ferr = err
		}
		return PHPString(r)
	})
	if ferr != nil {
		return nil, ferr
	}
	return out, nil
}

// getBoolEnv is Platform::getBoolEnv($name) (nil when unset or empty).
func (c *Config) getBoolEnv(name string) (*bool, error) {
	v, ok := c.Getenv(name)
	if !ok || v == "" {
		return nil, nil
	}
	switch v {
	case "1", "true", "on":
		b := true
		return &b, nil
	case "0", "false", "off":
		b := false
		return &b, nil
	}
	return nil, &RuntimeError{fmt.Sprintf("Invalid value for %s: %s. Expected 0, 1, false, true, off, or on.", name, v)}
}

var reMaxSize = regexp.MustCompile(`(?i)^\s*([0-9.]+)\s*(?:([kmg])(?:i?b)?)?\s*$`)

func phpBool(v phpjson.Value) bool { return v != "false" && phpjson.Truthy(v) }

// GetFlags is Config::get($key, $flags) with Composer's exact value
// semantics (int/size parsing, env overrides and their sources, errors).
func (c *Config) GetFlags(key string, flags int) (phpjson.Value, error) {
	cur, _ := c.config.Get(key)
	switch key {
	case "vendor-dir", "bin-dir", "process-timeout", "data-dir", "cache-dir", "cache-files-dir", "cache-repo-dir", "cache-vcs-dir", "cafile", "capath":
		envName := "COMPOSER_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		envVal, hasEnv := c.env(envName)
		if hasEnv {
			c.setSource(envVal, key, envName)
		}
		if key == "process-timeout" {
			if hasEnv {
				n := stringIntval(envVal)
				if n < 0 {
					n = 0
				}
				return n, nil
			}
			// max(0, $value) returns the value itself unless it is below 0
			switch t := cur.(type) {
			case int64:
				if t < 0 {
					return int64(0), nil
				}
			case float64:
				if t < 0 {
					return int64(0), nil
				}
			case string:
				if PHPFloatval(t) < 0 {
					return int64(0), nil
				}
			}
			return cur, nil
		}
		var raw phpjson.Value = cur
		if hasEnv {
			raw = envVal
		}
		processed, err := c.processFlags(raw, flags)
		if err != nil {
			return nil, err
		}
		s := strings.TrimRight(PHPString(processed), "/\\")
		s = ExpandPath(s)
		if !strings.HasSuffix(key, "-dir") {
			return s, nil
		}
		if flags&RelativePaths == RelativePaths {
			return s, nil
		}
		return c.realpath(s), nil
	case "cache-read-only", "htaccess-protect":
		envName := "COMPOSER_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		if v, ok := c.env(envName); ok {
			c.setSource(v, key, envName)
			return phpBool(v), nil
		}
		return phpBool(cur), nil
	case "disable-tls", "secure-http", "use-github-api", "lock", "source-fallback", "strict-psr-autoloader":
		if key == "secure-http" {
			if v, _ := c.GetFlags("disable-tls", 0); v == true {
				return false, nil
			}
		}
		return phpBool(cur), nil
	case "cache-ttl":
		n := PHPIntval(cur)
		if n < 0 {
			n = 0
		}
		return n, nil
	case "cache-files-maxsize":
		str := PHPString(cur)
		m := reMaxSize.FindStringSubmatch(str)
		if m == nil {
			return nil, &RuntimeError{"Could not parse the value of '" + key + "': " + str}
		}
		size := PHPFloatval(m[1])
		switch strings.ToLower(m[2]) {
		case "g":
			size *= 1024 * 1024 * 1024
		case "m":
			size *= 1024 * 1024
		case "k":
			size *= 1024
		}
		n := PHPIntval(size)
		if n < 0 {
			n = 0
		}
		return n, nil
	case "cache-files-ttl":
		if cur != nil {
			n := PHPIntval(cur)
			if n < 0 {
				n = 0
			}
			return n, nil
		}
		return c.GetFlags("cache-ttl", 0)
	case "home":
		processed, err := c.processFlags(ExpandPath(PHPString(cur)), flags)
		if err != nil {
			return nil, err
		}
		return strings.TrimRight(PHPString(processed), "/\\"), nil
	case "bin-compat":
		var value phpjson.Value = cur
		if v, ok := c.env("COMPOSER_BIN_COMPAT"); ok && phpjson.Truthy(v) {
			value = v
		}
		valid := value == true
		if s, ok := value.(string); ok {
			valid = s == "auto" || s == "full" || s == "proxy" || s == "symlink"
		}
		if !valid {
			return nil, &RuntimeError{"Invalid value for 'bin-compat': " + PHPString(value) + ". Expected auto, full or proxy"}
		}
		return value, nil
	case "discard-changes":
		if env, ok := c.env("COMPOSER_DISCARD_CHANGES"); ok {
			switch env {
			case "stash":
				return "stash", nil
			case "true", "false", "1", "0":
				return phpBool(env), nil
			}
			return nil, &RuntimeError{"Invalid value for COMPOSER_DISCARD_CHANGES: " + env + ". Expected 1, 0, true, false or stash"}
		}
		if cur != true && cur != false && cur != "stash" {
			return nil, &RuntimeError{"Invalid value for 'discard-changes': " + PHPString(cur) + ". Expected true, false or stash"}
		}
		return cur, nil
	case "github-protocols":
		protos := phpjson.NewMap()
		phpjson.Each(cur, func(k string, v phpjson.Value) { protos.Set(k, v) })
		secure, _ := c.config.Get("secure-http")
		if phpjson.Truthy(secure) {
			// array_search('git', $protos) (loose) then unset
			for _, k := range protos.Keys() {
				v, _ := protos.Get(k)
				if looseEqualsString(v, "git") {
					protos.Delete(k)
					break
				}
			}
		}
		if keys := protos.Keys(); len(keys) > 0 {
			if first, _ := protos.Get(keys[0]); first == "http" {
				return nil, &RuntimeError{`The http protocol for github is not available anymore, update your config's github-protocols to use "https", "git" or "ssh"`}
			}
		}
		return protos, nil
	case "autoloader-suffix":
		if cur == "" {
			return nil, nil
		}
		return c.processFlags(cur, flags)
	case "audit":
		result := phpjson.NewMap()
		if m := phpjson.ToMap(cur); m != nil {
			result = m.Clone()
		}
		if v, ok := c.env("COMPOSER_AUDIT_ABANDONED"); ok {
			if v != "ignore" && v != "report" && v != "fail" {
				return nil, &RuntimeError{"Invalid value for COMPOSER_AUDIT_ABANDONED: " + v + ". Expected one of ignore, report, fail."}
			}
			result.Set("abandoned", v)
		}
		if _, ok := c.env("COMPOSER_SECURITY_BLOCKING_ABANDONED"); ok {
			b, err := c.getBoolEnv("COMPOSER_SECURITY_BLOCKING_ABANDONED")
			if err != nil {
				return nil, err
			}
			if b == nil {
				result.Set("block-abandoned", nil)
			} else {
				result.Set("block-abandoned", *b)
			}
		}
		return result, nil
	case "policy":
		policyConfig := cur
		b, err := c.getBoolEnv("COMPOSER_POLICY")
		if err != nil {
			return nil, err
		}
		if b != nil && !*b {
			policyConfig = false
		} else if b != nil && *b && policyConfig == false {
			policyConfig = true
		}
		return policyConfig, nil
	}
	if cur == nil {
		return nil, nil
	}
	return c.processFlags(cur, flags)
}

// looseEqualsString is $v == $s for a non-numeric string $s.
func looseEqualsString(v phpjson.Value, s string) bool {
	switch t := v.(type) {
	case string:
		return t == s
	case bool:
		return t == (s != "" && s != "0")
	case nil:
		return s == ""
	}
	return false
}
