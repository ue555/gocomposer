// Package cache ports Composer\Cache: a directory of files with sanitized names.
package cache

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Cache is Composer\Cache.
type Cache struct {
	root      string
	allowlist *regexp.Regexp
	ReadOnly  bool
	enabled   *bool
	Warn      func(string)
}

var reUnusable = regexp.MustCompile(`(^|[\\/])(\$null|nul|NUL|/dev/null)([\\/]|$)`)

// New creates a cache rooted at dir; allowlist is a character class body like "a-z0-9._".
func New(dir, allowlist string, readOnly bool) *Cache {
	c := &Cache{root: strings.TrimRight(dir, "/\\") + "/", allowlist: regexp.MustCompile(`(?i)[^` + allowlist + `]`), ReadOnly: readOnly}
	if reUnusable.MatchString(dir) {
		f := false
		c.enabled = &f
	}
	return c
}

func (c *Cache) Root() string { return c.root }

func (c *Cache) IsEnabled() bool {
	if c.enabled == nil {
		ok := true
		if !c.ReadOnly {
			if err := os.MkdirAll(c.root, 0o777); err != nil || !writable(c.root) {
				if c.Warn != nil {
					c.Warn("<warning>Cannot create cache directory " + c.root + ", or directory is not writable. Proceeding without cache. See also cache-read-only config if your filesystem is read-only.</warning>")
				}
				ok = false
			}
		}
		c.enabled = &ok
	}
	return *c.enabled
}

func writable(dir string) bool {
	f, err := os.CreateTemp(dir, ".probe")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}

func (c *Cache) path(file string) string { return c.root + c.allowlist.ReplaceAllString(file, "-") }

// Read returns the cached contents or ok=false.
func (c *Cache) Read(file string) ([]byte, bool) {
	if !c.IsEnabled() {
		return nil, false
	}
	data, err := os.ReadFile(c.path(file))
	if err != nil {
		return nil, false
	}
	return data, true
}

// Write stores contents atomically (temp file + rename).
func (c *Cache) Write(file string, contents []byte) bool {
	if !c.IsEnabled() || c.ReadOnly {
		return false
	}
	p := c.path(file)
	if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
		return false
	}
	var rnd [5]byte
	_, _ = rand.Read(rnd[:])
	tmp := p + hex.EncodeToString(rnd[:]) + ".tmp"
	if err := os.WriteFile(tmp, contents, 0o666); err != nil {
		os.Remove(tmp)
		return false
	}
	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return false
	}
	return true
}

// CopyFrom copies source into the cache.
func (c *Cache) CopyFrom(file, source string) bool {
	if !c.IsEnabled() || c.ReadOnly {
		return false
	}
	p := c.path(file)
	if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
		return false
	}
	return copyFile(source, p) == nil
}

// CopyTo copies a cached file to target, refreshing its access time.
func (c *Cache) CopyTo(file, target string) bool {
	if !c.IsEnabled() {
		return false
	}
	p := c.path(file)
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	_ = os.Chtimes(p, time.Now(), st.ModTime())
	return copyFile(p, target) == nil
}

// Exists reports whether a file is present.
func (c *Cache) Exists(file string) bool {
	if !c.IsEnabled() {
		return false
	}
	_, err := os.Stat(c.path(file))
	return err == nil
}

// Path returns the absolute cache path of file.
func (c *Cache) Path(file string) string { return c.path(file) }

func (c *Cache) Remove(file string) {
	if c.IsEnabled() && !c.ReadOnly {
		os.Remove(c.path(file))
	}
}

// Sha1 / Sha256 of a cached file ("" if missing).
func (c *Cache) Sha1(file string) string   { return c.hash(file, sha1.New()) }
func (c *Cache) Sha256(file string) string { return c.hash(file, sha256.New()) }

func (c *Cache) hash(file string, h interface {
	io.Writer
	Sum([]byte) []byte
}) string {
	if !c.IsEnabled() {
		return ""
	}
	f, err := os.Open(c.path(file))
	if err != nil {
		return ""
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp := dst + ".tmp" + randHex()
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}

func randHex() string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
