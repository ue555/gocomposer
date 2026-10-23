// PHP value semantics needed by the JsonManipulator port
// (src/Composer/Json/JsonManipulator.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80): isset/empty/array access on
// json_decode() results, json_decode in object and associative mode,
// JsonFile::encode and the exceptions PHP would throw.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package jsonmanip

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/pcre"
	"gocomposer/internal/phpjson"
)

// Exception is a PHP exception (or Error) the ported code throws. Manipulator
// methods panic with *Exception exactly where PHP throws; use Try to turn
// that into an error.
type Exception struct {
	Class   string // PHP class name, e.g. "InvalidArgumentException"
	Message string
	Code    int
	Err     error // underlying Go error, if any
}

func (e *Exception) Error() string { return e.Message }
func (e *Exception) Unwrap() error { return e.Err }

func throw(class, msg string) { panic(&Exception{Class: class, Message: msg}) }

// Try runs fn and converts a panic with *Exception into an error.
func Try(fn func() bool) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, isExc := r.(*Exception); isExc {
				ok, err = false, e
				return
			}
			panic(r)
		}
	}()
	return fn(), nil
}

// PREG error codes (preg_last_error()).
const (
	pregBacktrackLimitError = 2
	pregJitStackLimitError  = 6
)

func pcreException(fn, pattern string, err error) *Exception {
	code := 1 // PREG_INTERNAL_ERROR
	msg := "Internal error"
	switch {
	case errors.Is(err, pcre.ErrBacktrackLimit):
		code = pregBacktrackLimitError
		msg = "Backtrack limit exhausted"
	case errors.Is(err, pcre.ErrRecursionLimit):
		code = pregJitStackLimitError
		msg = "JIT stack limit exhausted"
	}
	return &Exception{Class: "Composer\\Pcre\\PcreException", Message: fn + "(): failed executing \"" + pattern + "\": " + msg, Code: code, Err: err}
}

// ---- preg wrappers (Composer\Pcre\Preg) ----

// isMatch is Preg::isMatch; PCRE failures throw.
func isMatch(pattern, subject string) *pcre.Match {
	m, err := pcre.IsMatch(pattern, subject)
	if err != nil {
		panic(pcreException("preg_match", pattern, err))
	}
	return m
}

// isMatchCatch is Preg::isMatch inside try/catch for the given PREG error
// codes, which make it report no match (ok=false).
func isMatchCatch(pattern, subject string, codes ...int) (m *pcre.Match, ok bool) {
	m, err := pcre.IsMatch(pattern, subject)
	if err != nil {
		e := pcreException("preg_match", pattern, err)
		for _, c := range codes {
			if e.Code == c {
				return nil, false
			}
		}
		panic(e)
	}
	return m, true
}

// isMatchStrictGroups is Preg::isMatchStrictGroups.
func isMatchStrictGroups(pattern, subject string) *pcre.Match {
	m := isMatch(pattern, subject)
	if m == nil {
		return nil
	}
	for _, e := range m.PHPArray(true) {
		if e.Value == nil {
			panic(&Exception{Class: "Composer\\Pcre\\UnexpectedNullMatchException", Message: "Pattern \"" + pattern + "\" had an unexpected unmatched group \"" + e.Key + "\", make sure the pattern always matches or use match() instead.", Err: errors.New("unexpected null match")})
		}
	}
	return m
}

func pregReplace(pattern, replacement, subject string) string {
	s, _, err := pcre.Replace(pattern, replacement, subject, -1)
	if err != nil {
		panic(pcreException("preg_replace", pattern, err))
	}
	return s
}

func pregReplaceCallback(pattern string, fn func(*pcre.Match) string, subject string) string {
	s, _, err := pcre.ReplaceCallback(pattern, fn, subject, -1)
	if err != nil {
		panic(pcreException("preg_replace_callback", pattern, err))
	}
	return s
}

func quote(s string) string { return pcre.Quote(s, "") }

func addcslashes(s string) string { return pcre.AddCSlashes(s, `\$`) }

// ---- JSON ----

// encode is JsonFile::encode($data) with the default flags.
func encode(v phpjson.Value) string {
	if f, ok := v.(float64); ok {
		return FormatFloat(f)
	}
	s, err := phpjson.Encode(v, phpjson.ComposerFlags)
	if err != nil {
		msg := "Unknown error"
		if errors.Is(err, phpjson.ErrInvalidUTF8) {
			msg = "Malformed UTF-8 characters, possibly incorrectly encoded"
		}
		panic(&Exception{Class: "RuntimeException", Message: "JSON encoding failed: " + msg, Err: err})
	}
	return s
}

// parseJSON is JsonFile::parseJson($json) (associative arrays).
func parseJSON(s string) phpjson.Value {
	v, err := jsonfile.Parse([]byte(s), "")
	if err != nil {
		class := "Seld\\JsonLint\\ParsingException"
		var ee *jsonfile.EncodingError
		if errors.As(err, &ee) {
			class = "UnexpectedValueException"
		}
		panic(&Exception{Class: class, Message: err.Error(), Err: err})
	}
	return v
}

// decodeAssoc is @json_decode($s, true); ok=false means null due to an error.
func decodeAssoc(s string) (phpjson.Value, bool) {
	v, err := phpjson.Decode([]byte(s))
	if err != nil {
		return nil, false
	}
	return v, true
}

// decodeObject is @json_decode($s): JSON objects become stdClass, modelled
// as *phpjson.Map with ForceObject set.
func decodeObject(s string) (phpjson.Value, bool) {
	v, err := phpjson.Decode([]byte(s))
	if err != nil {
		return nil, false
	}
	markObjects(v)
	return v, true
}

func markObjects(v phpjson.Value) {
	switch t := v.(type) {
	case *phpjson.Map:
		t.ForceObject = true
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			markObjects(x)
		}
	case []phpjson.Value:
		for _, x := range t {
			markObjects(x)
		}
	}
}

// truthy is PHP's boolean conversion; objects (ForceObject maps) are
// always true.
func truthy(v phpjson.Value) bool {
	if m, ok := v.(*phpjson.Map); ok && m != nil && m.ForceObject {
		return true
	}
	return phpjson.Truthy(v)
}

// ---- arrays ----

func isArray(v phpjson.Value) bool {
	switch t := v.(type) {
	case *phpjson.Map:
		return t != nil && !t.ForceObject
	case []phpjson.Value:
		return true
	}
	return false
}

func isObject(v phpjson.Value) bool {
	m, ok := v.(*phpjson.Map)
	return ok && m != nil && m.ForceObject
}

// normKey canonicalises an array key the way PHP does (numeric strings
// become integers, which print back identically, so string keys suffice).
func normKey(k string) string { return k }

// index reads $arr[$key] for a PHP array (list or map).
func index(v phpjson.Value, key string) (phpjson.Value, bool) {
	switch t := v.(type) {
	case *phpjson.Map:
		if t == nil {
			return nil, false
		}
		return t.Get(normKey(key))
	case []phpjson.Value:
		if i, ok := phpjson.IntKey(key); ok && i >= 0 && i < int64(len(t)) {
			return t[i], true
		}
	}
	return nil, false
}

// isset mirrors isset($v[$k1][$k2]...), including string offsets.
func isset(v phpjson.Value, keys ...string) bool {
	cur := v
	for _, k := range keys {
		switch t := cur.(type) {
		case *phpjson.Map, []phpjson.Value:
			if isObject(t) {
				// stdClass does not support [] access; not used on objects
				return false
			}
			x, ok := index(t, k)
			if !ok || x == nil {
				return false
			}
			cur = x
		case string:
			i, ok := phpjson.IntKey(k)
			if !ok {
				return false
			}
			if i < 0 {
				i += int64(len(t))
			}
			if i < 0 || i >= int64(len(t)) {
				return false
			}
			cur = t[i : i+1]
		default:
			return false
		}
	}
	return cur != nil
}

// get is "$v[$k1][$k2]... ?? null".
func get(v phpjson.Value, keys ...string) phpjson.Value {
	if !isset(v, keys...) {
		return nil
	}
	cur := v
	for _, k := range keys {
		if s, ok := cur.(string); ok {
			i, _ := phpjson.IntKey(k)
			if i < 0 {
				i += int64(len(s))
			}
			cur = s[i : i+1]
			continue
		}
		cur, _ = index(cur, k)
	}
	return cur
}

func arrayKeyExists(v phpjson.Value, key string) bool {
	_, ok := index(v, key)
	return ok
}

func count(v phpjson.Value) int { return phpjson.ArrayLen(v) }

// toMap returns a mutable copy of a PHP array as an ordered map.
func toMap(v phpjson.Value) *phpjson.Map {
	switch t := v.(type) {
	case *phpjson.Map:
		c := t.Clone()
		c.ForceObject = false
		return c
	case []phpjson.Value:
		return phpjson.ToMap(phpjson.CloneValue(t))
	}
	return phpjson.NewMap()
}

// mapOf builds [$k => $v].
func mapOf(k string, v phpjson.Value) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set(normKey(k), v)
	return m
}

// prependName is ['name' => $name] + $array.
func prependName(name phpjson.Value, arr phpjson.Value) *phpjson.Map {
	out := phpjson.NewMap()
	out.Set("name", name)
	phpjson.Each(arr, func(k string, v phpjson.Value) {
		if k != "name" {
			out.Set(k, v)
		}
	})
	return out
}

// isIntKey reports whether a PHP array key is stored as an integer.
func isIntKey(k string) bool {
	_, ok := phpjson.IntKey(k)
	return ok
}

// isNumeric is PHP 8 is_numeric() for strings.
func isNumeric(s string) bool {
	t := strings.TrimLeft(s, " \t\n\r\v\f")
	t = strings.TrimRight(t, " \t\n\r\v\f")
	if t == "" {
		return false
	}
	i := 0
	if t[i] == '+' || t[i] == '-' {
		i++
	}
	digits := 0
	for i < len(t) && t[i] >= '0' && t[i] <= '9' {
		i++
		digits++
	}
	if i < len(t) && t[i] == '.' {
		i++
		for i < len(t) && t[i] >= '0' && t[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return false
	}
	if i < len(t) && (t[i] == 'e' || t[i] == 'E') {
		j := i + 1
		if j < len(t) && (t[j] == '+' || t[j] == '-') {
			j++
		}
		ed := 0
		for j < len(t) && t[j] >= '0' && t[j] <= '9' {
			j++
			ed++
		}
		if ed > 0 {
			i = j
		}
	}
	return i == len(t)
}

// strictEqualFalseMap is ([$k => false] === $v).
func strictEqualFalseMap(k string, v phpjson.Value) bool {
	if !isArray(v) || count(v) != 1 {
		return false
	}
	x, ok := index(v, k)
	if !ok {
		return false
	}
	b, isBool := x.(bool)
	return isBool && !b
}

// keyString renders a PHP array key for JsonFile::encode: integer keys are
// encoded as numbers.
func encodeKey(k string) string {
	if i, ok := phpjson.IntKey(k); ok {
		return strconv.FormatInt(i, 10)
	}
	return encode(k)
}

// rtrimChars is rtrim($s, $chars) for a plain character list.
func rtrimChars(s, chars string) string {
	return strings.TrimRight(s, chars)
}

func typeName(v phpjson.Value) string {
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
	}
	if isObject(v) {
		return "stdClass"
	}
	return "array"
}

// FormatFloat is json_encode() of a float with serialize_precision=-1 and
// without JSON_PRESERVE_ZERO_FRACTION: zend_gcvt(value, 17, '.', 'e') on the
// shortest round-trip digits (smart_str_append_double).
func FormatFloat(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "0" // json_encode fails with JSON_ERROR_INF_OR_NAN; callers never pass these
	}
	if f == 0 {
		if math.Signbit(f) {
			return "-0"
		}
		return "0"
	}
	s := strconv.FormatFloat(f, 'e', -1, 64) // d.ddde±XX
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}
	mant, exp, _ := strings.Cut(s, "e")
	digits := strings.Replace(mant, ".", "", 1)
	e, _ := strconv.Atoi(exp)
	decpt := e + 1 // value = 0.digits * 10^decpt
	var out string
	if decpt < 0 && decpt < -3 || decpt >= 0 && decpt > 17 {
		// exponential format
		out = digits[:1] + "."
		if len(digits) > 1 {
			out += digits[1:]
		} else {
			out += "0"
		}
		x := decpt - 1
		if x < 0 {
			out += "e-" + strconv.Itoa(-x)
		} else {
			out += "e+" + strconv.Itoa(x)
		}
	} else if decpt <= 0 {
		out = "0." + strings.Repeat("0", -decpt) + digits
	} else {
		if len(digits) <= decpt {
			out = digits + strings.Repeat("0", decpt-len(digits))
		} else {
			out = digits[:decpt] + "." + digits[decpt:]
		}
	}
	if neg {
		out = "-" + out
	}
	return out
}
