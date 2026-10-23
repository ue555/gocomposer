package cli

// Port of Composer\Command\ClearCacheCommand and the gc/clear parts of
// Composer\Cache (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/cache"
	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
)

func init() { registerCommand(clearCacheCommand) }

func clearCacheCommand() *Command {
	return &Command{Name: "clear-cache", Run: runClearCache}
}

func runClearCache(ctx *Context) (int, error) {
	io := ctx.IO
	var cfg *config.Config
	if ctx.cache != nil {
		cfg = ctx.cache.Config
	} else {
		cwd, _ := os.Getwd()
		c, err := config.CreateConfig(cwd, func(m string) { io.WriteError(m) })
		if err != nil {
			return 1, err
		}
		cfg = c
	}
	gc := ctx.In.Bool("gc")
	for _, key := range []string{"cache-vcs-dir", "cache-repo-dir", "cache-files-dir", "cache-dir"} {
		if key == "cache-dir" && gc {
			continue
		}
		cachePath := cfg.GetString(key)
		if !cacheUsable(cachePath) {
			io.WriteError(fmt.Sprintf("<info>Cache is not enabled (%s): %s</info>", key, cachePath))
			continue
		}
		real := phpRealpath(cachePath)
		if real == "" {
			io.WriteError(fmt.Sprintf("<info>Cache directory does not exist (%s): </info>", key))
			continue
		}
		readOnly := cfg.GetBool("cache-read-only")
		cc := cache.New(real, "a-z0-9._", readOnly)
		cc.Warn = func(m string) { io.WriteError(m) }
		if !cc.IsEnabled() {
			io.WriteError(fmt.Sprintf("<info>Cache is not enabled (%s): %s</info>", key, real))
			continue
		}
		root := cc.Root()
		if gc {
			io.WriteError(fmt.Sprintf("<info>Garbage-collecting cache (%s): %s</info>", key, real))
			switch key {
			case "cache-files-dir":
				maxSize, err := cacheFilesMaxSize(cfg)
				if err != nil {
					return 1, err
				}
				if !readOnly {
					cacheGC(root, phpInt(cfg.Get("cache-files-ttl"), cfg.Get("cache-ttl")), maxSize)
				}
			case "cache-repo-dir":
				if !readOnly {
					cacheGC(root, phpInt(cfg.Get("cache-ttl"), nil), 1024*1024*1024)
				}
			case "cache-vcs-dir":
				cacheGCVcs(root, phpInt(cfg.Get("cache-ttl"), nil))
			}
		} else {
			io.WriteError(fmt.Sprintf("<info>Clearing cache (%s): %s</info>", key, real))
			if !readOnly {
				emptyDirectory(root)
			}
		}
	}
	if gc {
		io.WriteError("<info>All caches garbage-collected.</info>")
	} else {
		io.WriteError("<info>All caches cleared.</info>")
	}
	return 0, nil
}

var reCacheUnusable = regexp.MustCompile(`(^|[\\/])(\$null|nul|NUL|/dev/null)([\\/]|$)`)

// cacheUsable is Cache::isUsable.
func cacheUsable(p string) bool { return !reCacheUnusable.MatchString(p) }

// phpRealpath is PHP's realpath(): "" when the path does not exist.
func phpRealpath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	r, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return ""
	}
	return r
}

// phpInt is max(0, (int) value) with a fallback for null values.
func phpInt(v, fallback phpjson.Value) int64 {
	if v == nil {
		v = fallback
	}
	var n int64
	switch t := v.(type) {
	case int64:
		n = t
	case float64:
		n = int64(t)
	case bool:
		if t {
			n = 1
		}
	case string:
		m := regexp.MustCompile(`^\s*[+-]?\d+`).FindString(t)
		n, _ = strconv.ParseInt(strings.TrimSpace(m), 10, 64)
	}
	if n < 0 {
		n = 0
	}
	return n
}

var reMaxSize = regexp.MustCompile(`(?i)^\s*([0-9.]+)\s*(?:([kmg])(?:i?b)?)?\s*$`)

// cacheFilesMaxSize is Config::get('cache-files-maxsize').
func cacheFilesMaxSize(cfg *config.Config) (int64, error) {
	raw := cfg.Raw()
	v, _ := raw.Get("cache-files-maxsize")
	s, _ := phpjson.ScalarString(v)
	m := reMaxSize.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("Could not parse the value of 'cache-files-maxsize': %s", s)
	}
	size, _ := strconv.ParseFloat(leadingFloat(m[1]), 64)
	switch strings.ToLower(m[2]) {
	case "g":
		size *= 1024 * 1024 * 1024
	case "m":
		size *= 1024 * 1024
	case "k":
		size *= 1024
	}
	if size < 0 {
		size = 0
	}
	return int64(size), nil
}

// leadingFloat is the numeric prefix PHP's (float) cast reads.
func leadingFloat(s string) string {
	if m := regexp.MustCompile(`^\d*\.?\d*`).FindString(s); m != "" && m != "." {
		return m
	}
	return "0"
}

// emptyDirectory is Filesystem::emptyDirectory.
func emptyDirectory(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		_ = os.MkdirAll(root, 0o777)
		return
	}
	for _, e := range entries {
		_ = os.RemoveAll(filepath.Join(root, e.Name()))
	}
}

func finderFiles(root string) []string        { return cache.FinderFiles(root) }
func dirSize(root string) int64               { return cache.DirSize(root) }
func cacheGC(root string, ttl, maxSize int64) { cache.GC(root, ttl, maxSize) }
func cacheGCVcs(root string, ttl int64)       { cache.GCVcs(root, ttl) }
