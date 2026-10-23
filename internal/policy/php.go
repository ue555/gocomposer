// Package policy ports Composer's dependency policy layer: Composer\Policy\*,
// Composer\FilterList\*, the policy pool filters of
// Composer\DependencyResolver (SecurityAdvisoryPoolFilter,
// FilterListPoolFilter, CooldownPoolFilter, PolicyRemovalReason) and the
// auditing parts of Composer\Advisory (Auditor, CooldownAuditor, advisories).
//
// Ported from Composer at f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// This file holds the PHP-semantics helpers shared by the port (ordered
// arrays, isset/??, get_debug_type, exception classes, env access).
package policy

import (
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// Exception class names carried by Error, mirroring the PHP exception types
// Composer throws so callers can render "[Class] message" like Composer.
const (
	UnexpectedValueException = "UnexpectedValueException"
	InvalidArgumentException = "InvalidArgumentException"
	RuntimeException         = "RuntimeException"
	LogicException           = "LogicException"
)

// Error is a PHP exception raised by the policy layer.
type Error struct {
	Class string
	Msg   string
	Prev  error
}

func (e *Error) Error() string { return e.Msg }
func (e *Error) Unwrap() error { return e.Prev }

func newErr(class, msg string) *Error { return &Error{Class: class, Msg: msg} }

// LookupEnv is Platform::getEnv (os.LookupEnv); replaceable for embedding.
var LookupEnv = os.LookupEnv

// getEnv mirrors Platform::getEnv: ok=false when the variable is unset.
func getEnv(name string) (string, bool) { return LookupEnv(name) }

// GetBoolEnv ports Platform::getBoolEnv: nil when unset or empty, an error
// for any value other than 0/1/false/true/off/on.
func GetBoolEnv(name string) (*bool, error) {
	v, ok := getEnv(name)
	if !ok || v == "" {
		return nil, nil
	}
	switch v {
	case "1", "true", "on":
		t := true
		return &t, nil
	case "0", "false", "off":
		f := false
		return &f, nil
	}
	return nil, newErr(RuntimeException, "Invalid value for "+name+": "+v+". Expected 0, 1, false, true, off, or on.")
}

// OMap is an ordered PHP array with string keys.
type OMap[V any] struct {
	keys []string
	vals map[string]V
}

// NewOMap creates an empty ordered map.
func NewOMap[V any]() *OMap[V] { return &OMap[V]{vals: map[string]V{}} }

func (m *OMap[V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

// Keys returns the keys in insertion order (array_keys).
func (m *OMap[V]) Keys() []string {
	if m == nil {
		return nil
	}
	return append([]string(nil), m.keys...)
}

func (m *OMap[V]) Get(k string) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	v, ok := m.vals[k]
	return v, ok
}

// Has is array_key_exists.
func (m *OMap[V]) Has(k string) bool { _, ok := m.Get(k); return ok }

// Set assigns in place or appends a new key at the end.
func (m *OMap[V]) Set(k string, v V) {
	if m.vals == nil {
		m.vals = map[string]V{}
	}
	if _, ok := m.vals[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = v
}

// Delete is unset().
func (m *OMap[V]) Delete(k string) {
	if m == nil {
		return
	}
	if _, ok := m.vals[k]; !ok {
		return
	}
	delete(m.vals, k)
	for i, kk := range m.keys {
		if kk == k {
			m.keys = append(m.keys[:i:i], m.keys[i+1:]...)
			break
		}
	}
}

// Each iterates in order.
func (m *OMap[V]) Each(fn func(k string, v V)) {
	if m == nil {
		return
	}
	for _, k := range m.keys {
		fn(k, m.vals[k])
	}
}

// Clone is a shallow copy (PHP array value semantics).
func (m *OMap[V]) Clone() *OMap[V] {
	c := NewOMap[V]()
	m.Each(func(k string, v V) { c.Set(k, v) })
	return c
}

// KSort sorts keys like PHP ksort() with SORT_REGULAR.
func (m *OMap[V]) KSort() {
	if m == nil {
		return
	}
	keys := m.keys
	// insertion sort keeps it stable and simple; maps here are small
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && phpjson.PHPCompareKeys(keys[j-1], keys[j]) > 0; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
}

func appendTo[V any](m *OMap[[]V], k string, v V) {
	cur, _ := m.Get(k)
	m.Set(k, append(cur, v))
}

// --- phpjson value helpers (PHP array access semantics) ---

// arrGet reads $arr[$key] from a PHP array value (list or map).
func arrGet(v phpjson.Value, key string) (phpjson.Value, bool) {
	switch t := v.(type) {
	case *phpjson.Map:
		return t.Get(key)
	case []phpjson.Value:
		if n, ok := phpjson.IntKey(key); ok && n >= 0 && int(n) < len(t) {
			return t[n], true
		}
	}
	return nil, false
}

// isset is isset($arr[$key]).
func isset(v phpjson.Value, key string) bool {
	x, ok := arrGet(v, key)
	return ok && x != nil
}

// coalesce is $arr[$key] ?? $def.
func coalesce(v phpjson.Value, key string, def phpjson.Value) phpjson.Value {
	if x, ok := arrGet(v, key); ok && x != nil {
		return x
	}
	return def
}

// isIntKey reports whether PHP would store the array key as an int.
func isIntKey(k string) bool { _, ok := phpjson.IntKey(k); return ok }

// eachArr iterates a PHP array value; non-arrays are skipped like foreach
// over a scalar (which only emits a warning).
func eachArr(v phpjson.Value, fn func(k string, v phpjson.Value)) { phpjson.Each(v, fn) }

// debugType is get_debug_type().
func debugType(v phpjson.Value) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "bool"
	case int64, int:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case []phpjson.Value, *phpjson.Map:
		return "array"
	}
	return "mixed"
}

// toBool is a (bool) cast.
func toBool(v phpjson.Value) bool { return phpjson.Truthy(v) }

// toStr is a (string) cast of a scalar (arrays become "Array").
func toStr(v phpjson.Value) string {
	if phpjson.IsArray(v) {
		return "Array"
	}
	s, _ := phpjson.ScalarString(v)
	return s
}

// nullableStr converts a ?string config value.
func nullableStr(v phpjson.Value) *string {
	if v == nil {
		return nil
	}
	s := toStr(v)
	return &s
}

func strPtr(s string) *string { return &s }

func inStrings(s string, list []string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// arrayValuesStrings returns the string values of a PHP array (in_array helpers).
func arrayValues(v phpjson.Value) []phpjson.Value {
	var out []phpjson.Value
	eachArr(v, func(_ string, x phpjson.Value) { out = append(out, x) })
	return out
}

// inArrayStrict is in_array($needle, $arr, true) for a string needle.
func inArrayStrict(needle string, arr phpjson.Value) bool {
	for _, x := range arrayValues(arr) {
		if s, ok := x.(string); ok && s == needle {
			return true
		}
	}
	return false
}

// PackageNameToRegexp is BasePackage::packageNameToRegexp.
func PackageNameToRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(pattern), `\*`, `.*`) + `$`)
}

// PackageNamesToRegexp is BasePackage::packageNamesToRegexp.
func PackageNamesToRegexp(names []string) *regexp.Regexp {
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = strings.ReplaceAll(regexp.QuoteMeta(n), `\*`, `.*`)
	}
	return regexp.MustCompile(`(?i)^(?:` + strings.Join(parts, "|") + `)$`)
}

// --- package helpers ---

// isRootPackage is "instanceof RootPackageInterface" (RootPackage and RootAliasPackage).
func isRootPackage(p *pkg.Package) bool { return p.IsRoot }

// isRootPackageAlias is "instanceof AliasPackage && isRootPackageAlias()".
func isRootPackageAlias(p *pkg.Package) bool { return p.AliasOf != nil && p.RootPackageAlias }

func isPlatformPackage(name string) bool { return repo.IsPlatformPackage(name) }

// IsUnversionedBranch is RepositorySet::isUnversionedBranch.
func IsUnversionedBranch(p *pkg.Package) bool {
	return strings.HasPrefix(p.Version(), "dev-") || p.Version() == semver.DefaultBranchAlias
}
