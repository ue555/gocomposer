// Package locker ports Composer\Package\Locker.
package locker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// ConflictedContentHash is JsonFile::CONFLICTED_CONTENT_HASH.
const ConflictedContentHash = "The composer.lock file contains unresolved merge conflicts."

// Locker manages composer.lock.
type Locker struct {
	File        string // "" or /dev/null disables writing
	hash        string
	contentHash string
	data        *phpjson.Map
	virtual     bool
	// PackageTime returns commit time for dev packages installed from source.
	PackageTime func(p *pkg.Package) string
	// ReadLog is called before the lock file is read (JsonFile::read's
	// -vvv "Reading" line); nil: no output.
	ReadLog func(path string)
}

// New creates a locker for lockFile from the raw composer.json contents.
func New(lockFile string, composerJSON []byte) (*Locker, error) {
	h := md5.Sum(composerJSON)
	ch, err := ContentHash(composerJSON)
	if err != nil {
		return nil, err
	}
	return &Locker{File: lockFile, hash: hex.EncodeToString(h[:]), contentHash: ch}, nil
}

// ContentHash is Locker::getContentHash.
func ContentHash(composerJSON []byte) (string, error) {
	v, err := phpjson.Decode(composerJSON)
	if err != nil {
		return "", err
	}
	content := phpjson.ToMap(v)
	if content == nil {
		content = phpjson.NewMap()
	}
	relevant := phpjson.NewMap()
	for _, key := range []string{"name", "version", "require", "require-dev", "conflict", "replace", "provide", "minimum-stability", "prefer-stable", "repositories", "extra"} {
		if val, ok := content.Get(key); ok {
			relevant.Set(key, val)
		}
	}
	if cfg := content.Map("config"); cfg != nil {
		c := phpjson.NewMap()
		if p, ok := cfg.Get("platform"); ok && p != nil {
			c.Set("platform", p)
		}
		if p, ok := cfg.Get("policy"); ok && p != nil {
			c.Set("policy", p)
		}
		if c.Len() > 0 {
			relevant.Set("config", c)
		}
	}
	relevant.KSort()
	enc, err := phpjson.Encode(relevant, 0)
	if err != nil {
		return "", err
	}
	sum := md5.Sum([]byte(enc))
	return hex.EncodeToString(sum[:]), nil
}

func (l *Locker) ContentHashValue() string { return l.contentHash }

func (l *Locker) enabled() bool { return l.File != "" && l.File != os.DevNull }

// Exists reports whether the lock file exists.
func (l *Locker) Exists() bool {
	if !l.enabled() {
		return false
	}
	_, err := os.Stat(l.File)
	return err == nil
}

// LockData reads and caches composer.lock.
func (l *Locker) LockData() (*phpjson.Map, error) {
	if l.data != nil {
		return l.data, nil
	}
	if !l.Exists() {
		return nil, fmt.Errorf("No lockfile found. Unable to read locked packages")
	}
	if l.ReadLog != nil {
		l.ReadLog(l.File)
	}
	raw, err := os.ReadFile(l.File)
	if err != nil {
		return nil, err
	}
	v, err := jsonfile.Parse(raw, l.File)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	l.data = m
	return m, nil
}

// IsLocked is Locker::isLocked.
func (l *Locker) IsLocked() (bool, error) {
	if !l.virtual && !l.Exists() {
		return false, nil
	}
	data, err := l.LockData()
	if err != nil {
		return false, err
	}
	return data.Has("packages"), nil
}

// IsFresh is Locker::isFresh.
func (l *Locker) IsFresh() bool {
	if l.data != nil && l.ReadLog != nil && l.Exists() {
		// Locker::isFresh reads the lock file again (not the cached data)
		l.ReadLog(l.File)
	}
	data, err := l.LockData()
	if err != nil {
		return false
	}
	if ch, ok := data.Get("content-hash"); ok && phpjson.Truthy(ch) {
		return l.contentHash == ch
	}
	if h, ok := data.Get("hash"); ok && phpjson.Truthy(h) {
		return l.hash == h
	}
	return false
}

// LockedRepository is Locker::getLockedRepository.
func (l *Locker) LockedRepository(withDev bool) (*repo.Array, error) {
	data, err := l.LockData()
	if err != nil {
		return nil, err
	}
	r := repo.NewLockArray()
	var locked []phpjson.Value
	phpjson.Each(mustGet(data, "packages"), func(_ string, v phpjson.Value) { locked = append(locked, v) })
	if withDev {
		dev, ok := data.Get("packages-dev")
		if !ok || dev == nil {
			return nil, fmt.Errorf("The lock file does not contain require-dev information, run install with the --no-dev option or delete it and run composer update to generate a new lock file.")
		}
		phpjson.Each(dev, func(_ string, v phpjson.Value) { locked = append(locked, v) })
	}
	if len(locked) == 0 {
		return r, nil
	}
	if first := phpjson.ToMap(locked[0]); first == nil || !first.Isset("name") {
		return nil, fmt.Errorf("Your composer.lock is invalid. Run \"composer update\" to generate a new one.")
	}
	loader := &pkg.Loader{LoadOptions: true}
	byName := map[string]*pkg.Package{}
	for _, info := range locked {
		p, err := loader.Load(phpjson.ToMap(info))
		if err != nil {
			return nil, err
		}
		r.AddPackage(p)
		byName[p.Name] = p
		if p.AliasOf != nil {
			byName[p.AliasOf.Name] = p.AliasOf
		}
	}
	if aliases, ok := data.Get("aliases"); ok {
		phpjson.Each(aliases, func(_ string, av phpjson.Value) {
			a := phpjson.ToMap(av)
			if a == nil {
				return
			}
			if base, ok := byName[a.Str("package")]; ok {
				alias := pkg.NewAlias(base, a.Str("alias_normalized"), a.Str("alias"))
				alias.RootPackageAlias = true
				r.AddPackage(alias)
			}
		})
	}
	return r, nil
}

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

// DevPackageNames is Locker::getDevPackageNames.
func (l *Locker) DevPackageNames() []string {
	data, err := l.LockData()
	if err != nil {
		return nil
	}
	var out []string
	phpjson.Each(mustGet(data, "packages-dev"), func(_ string, v phpjson.Value) {
		if m := phpjson.ToMap(v); m != nil {
			out = append(out, strings.ToLower(m.Str("name")))
		}
	})
	return out
}

// PlatformRequirements is Locker::getPlatformRequirements.
func (l *Locker) PlatformRequirements(withDev bool) (pkg.Links, error) {
	data, err := l.LockData()
	if err != nil {
		return nil, err
	}
	var out pkg.Links
	if p, ok := data.Get("platform"); ok && phpjson.Truthy(p) {
		links, err := pkg.ParseLinks("__root__", "1.0.0", pkg.TypeRequire, p)
		if err != nil {
			return nil, err
		}
		out = links
	}
	if withDev {
		if p, ok := data.Get("platform-dev"); ok && phpjson.Truthy(p) {
			links, err := pkg.ParseLinks("__root__", "1.0.0", pkg.TypeRequire, p)
			if err != nil {
				return nil, err
			}
			// array_merge with string keys: later entries override in place
			for _, nl := range links {
				replaced := false
				for i, e := range out {
					if e.Target == nl.Target {
						out[i] = nl
						replaced = true
						break
					}
				}
				if !replaced {
					out = append(out, nl)
				}
			}
		}
	}
	return out, nil
}

func (l *Locker) MinimumStability() string {
	data, _ := l.LockData()
	if data != nil {
		if s, ok := data.Get("minimum-stability"); ok && s != nil {
			str, _ := s.(string)
			return str
		}
	}
	return "stable"
}

func (l *Locker) StabilityFlags() map[string]int {
	out := map[string]int{}
	data, _ := l.LockData()
	if data == nil {
		return out
	}
	phpjson.Each(mustGet(data, "stability-flags"), func(k string, v phpjson.Value) {
		switch t := v.(type) {
		case int64:
			out[k] = int(t)
		case float64:
			out[k] = int(t)
		}
	})
	return out
}

// PreferStable returns nil when unset.
func (l *Locker) PreferStable() *bool { return l.boolKey("prefer-stable") }
func (l *Locker) PreferLowest() *bool { return l.boolKey("prefer-lowest") }

func (l *Locker) boolKey(k string) *bool {
	data, _ := l.LockData()
	if data == nil {
		return nil
	}
	v, ok := data.Get(k)
	if !ok || v == nil {
		return nil
	}
	b := phpjson.Truthy(v)
	return &b
}

// PlatformOverrides is getPlatformOverrides.
func (l *Locker) PlatformOverrides() *phpjson.Map {
	data, _ := l.LockData()
	if data == nil {
		return phpjson.NewMap()
	}
	if m := data.Map("platform-overrides"); m != nil {
		return m
	}
	return phpjson.NewMap()
}

// Aliases returns the lock file aliases.
func (l *Locker) Aliases() []pkg.RootAlias {
	data, _ := l.LockData()
	if data == nil {
		return nil
	}
	var out []pkg.RootAlias
	phpjson.Each(mustGet(data, "aliases"), func(_ string, v phpjson.Value) {
		a := phpjson.ToMap(v)
		if a == nil {
			return
		}
		out = append(out, pkg.RootAlias{Package: a.Str("package"), Version: a.Str("version"), Alias: a.Str("alias"), AliasNormalized: a.Str("alias_normalized")})
	})
	return out
}

// LockOptions bundles setLockData arguments.
type LockOptions struct {
	Packages          []*pkg.Package
	DevPackages       []*pkg.Package // nil means packages-dev: null
	DevPackagesSet    bool
	PlatformReqs      *phpjson.Map
	PlatformDevReqs   *phpjson.Map
	Aliases           []pkg.RootAlias
	MinimumStability  string
	StabilityFlags    map[string]int
	PreferStable      bool
	PreferLowest      bool
	PlatformOverrides *phpjson.Map
	Write             bool
}

// SetLockData is Locker::setLockData. It reports whether the lock changed.
func (l *Locker) SetLockData(o LockOptions) (bool, error) {
	lock := phpjson.NewMap()
	lock.Set("_readme", []phpjson.Value{
		"This file locks the dependencies of your project to a known state",
		"Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
		"This file is @generated automatically",
	})
	lock.Set("content-hash", l.contentHash)
	pk, err := l.lockPackages(o.Packages)
	if err != nil {
		return false, err
	}
	lock.Set("packages", pk)
	lock.Set("packages-dev", nil)
	var aliases []phpjson.Value
	for _, a := range o.Aliases {
		m := phpjson.NewMap()
		m.Set("package", a.Package)
		v := a.Version
		if v == "dev-master" || v == "dev-trunk" || v == "dev-default" {
			v = semver.DefaultBranchAlias
		}
		m.Set("version", v)
		m.Set("alias", a.Alias)
		m.Set("alias_normalized", a.AliasNormalized)
		aliases = append(aliases, m)
	}
	if aliases == nil {
		aliases = []phpjson.Value{}
	}
	lock.Set("aliases", aliases)
	lock.Set("minimum-stability", o.MinimumStability)
	flags := phpjson.NewMap()
	var names []string
	for n := range o.StabilityFlags {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		flags.Set(n, int64(o.StabilityFlags[n]))
	}
	flags.KSort()
	flags.ForceObject = true
	lock.Set("stability-flags", flags)
	lock.Set("prefer-stable", o.PreferStable)
	lock.Set("prefer-lowest", o.PreferLowest)
	if o.DevPackagesSet {
		dp, err := l.lockPackages(o.DevPackages)
		if err != nil {
			return false, err
		}
		lock.Set("packages-dev", dp)
	}
	plat := o.PlatformReqs
	if plat == nil {
		plat = phpjson.NewMap()
	}
	plat = plat.Clone()
	plat.ForceObject = true
	lock.Set("platform", plat)
	platDev := o.PlatformDevReqs
	if platDev == nil {
		platDev = phpjson.NewMap()
	}
	platDev = platDev.Clone()
	platDev.ForceObject = true
	lock.Set("platform-dev", platDev)
	if o.PlatformOverrides != nil && o.PlatformOverrides.Len() > 0 {
		lock.Set("platform-overrides", o.PlatformOverrides)
	}
	lock.Set("plugin-api-version", platform.PluginAPIVersion)
	isLocked, err := l.IsLocked()
	if err != nil {
		isLocked = false
	}
	newEnc, err := phpjson.Encode(lock, phpjson.ComposerFlags)
	if err != nil {
		return false, err
	}
	if isLocked {
		cur, _ := l.LockData()
		curEnc, _ := phpjson.Encode(cur, phpjson.ComposerFlags)
		if curEnc == newEnc && lockEqualStrict(cur, lock) {
			return false, nil
		}
	}
	if o.Write {
		if err := fsutil.FilePutContentsIfModified(l.File, []byte(newEnc+"\n")); err != nil {
			return false, err
		}
		l.data = nil
		l.virtual = false
	} else {
		l.virtual = true
		v, _ := phpjson.Decode([]byte(newEnc))
		l.data = phpjson.ToMap(v)
	}
	return true, nil
}

// lockEqualStrict approximates PHP's $lock !== $this->getLockData(): the
// decoded file compares as arrays, so empty objects equal empty arrays.
func lockEqualStrict(a, b *phpjson.Map) bool {
	ea, _ := phpjson.Encode(normalizeObjects(a), 0)
	eb, _ := phpjson.Encode(normalizeObjects(b), 0)
	return ea == eb
}

func normalizeObjects(v phpjson.Value) phpjson.Value {
	switch t := v.(type) {
	case *phpjson.Map:
		c := phpjson.NewMap()
		for _, k := range t.Keys() {
			e, _ := t.Get(k)
			c.Set(k, normalizeObjects(e))
		}
		return c
	case []phpjson.Value:
		out := make([]phpjson.Value, len(t))
		for i, e := range t {
			out[i] = normalizeObjects(e)
		}
		return out
	}
	return v
}

func (l *Locker) lockPackages(packages []*pkg.Package) ([]phpjson.Value, error) {
	var locked []*phpjson.Map
	for _, p := range packages {
		if p.AliasOf != nil {
			continue
		}
		if p.PrettyName == "" || p.PrettyVersion() == "" {
			return nil, fmt.Errorf("Package \"%s\" has no version or name and can not be locked", p)
		}
		spec := pkg.Dump(p)
		spec.Delete("version_normalized")
		t, hasTime := spec.Get("time")
		pt, hasPT := spec.Get("published-time")
		spec.Delete("time")
		spec.Delete("published-time")
		if p.IsDev() && p.InstallationSource() == "source" && l.PackageTime != nil {
			if pt := l.PackageTime(p); pt != "" {
				t, hasTime = pt, true
			}
		}
		if hasTime && t != nil {
			spec.Set("time", t)
		}
		if hasPT && pt != nil {
			spec.Set("published-time", pt)
		}
		spec.Delete("installation-source")
		locked = append(locked, spec)
	}
	sort.SliceStable(locked, func(i, j int) bool {
		a, b := locked[i], locked[j]
		if c := strings.Compare(a.Str("name"), b.Str("name")); c != 0 {
			return c < 0
		}
		return strings.Compare(a.Str("version"), b.Str("version")) < 0
	})
	out := make([]phpjson.Value, len(locked))
	for i, m := range locked {
		out[i] = m
	}
	return out, nil
}

// UpdateHash is Locker::updateHash: rewrite only the content-hash.
func (l *Locker) UpdateHash(composerJSON []byte, process func(*phpjson.Map) *phpjson.Map) error {
	ch, err := ContentHash(composerJSON)
	if err != nil {
		return err
	}
	st, statErr := os.Stat(l.File)
	l.data = nil
	data, err := l.LockData()
	if err != nil {
		return err
	}
	data.Set("content-hash", ch)
	if process != nil {
		data = process(data)
	}
	for _, k := range []string{"stability-flags", "platform", "platform-dev"} {
		if m := data.Map(k); m != nil && m.Len() == 0 {
			m.ForceObject = true
		} else if v, ok := data.Get(k); ok && phpjson.IsArray(v) && phpjson.ArrayLen(v) == 0 {
			nm := phpjson.NewMap()
			nm.ForceObject = true
			data.Set(k, nm)
		}
	}
	if m := data.Map("stability-flags"); m != nil {
		m.KSort()
	}
	enc, err := phpjson.Encode(data, phpjson.ComposerFlags)
	if err != nil {
		return err
	}
	if err := fsutil.FilePutContentsIfModified(l.File, []byte(enc+"\n")); err != nil {
		return err
	}
	l.data = nil
	l.contentHash = ch
	if statErr == nil {
		_ = os.Chtimes(l.File, st.ModTime(), st.ModTime())
	}
	return nil
}
