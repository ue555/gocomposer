// Package phpjson decodes and encodes JSON the way Composer does through PHP's
// json_decode(..., true) and json_encode(): objects keep their key order,
// empty objects become empty PHP arrays (encoded as []), maps with keys
// "0".."n-1" encode as lists, and string escaping follows PHP's flags.
package phpjson

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Value is one of: nil, bool, int64, float64, string, []Value, *Map.
type Value = any

// Map is an ordered PHP associative array. ForceObject makes an empty map
// encode as {} (PHP stdClass) instead of [].
type Map struct {
	keys        []string
	vals        map[string]Value
	ForceObject bool
}

func NewMap() *Map { return &Map{vals: map[string]Value{}} }

func (m *Map) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

func (m *Map) Keys() []string {
	if m == nil {
		return nil
	}
	return append([]string(nil), m.keys...)
}

func (m *Map) Get(k string) (Value, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.vals[k]
	return v, ok
}

// Has reports whether k exists (PHP array_key_exists).
func (m *Map) Has(k string) bool { _, ok := m.Get(k); return ok }

// Isset mirrors PHP isset(): the key exists and the value is not null.
func (m *Map) Isset(k string) bool { v, ok := m.Get(k); return ok && v != nil }

// Set replaces the value in place or appends a new key at the end.
func (m *Map) Set(k string, v Value) {
	if m.vals == nil {
		m.vals = map[string]Value{}
	}
	if _, ok := m.vals[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = v
}

func (m *Map) Delete(k string) {
	if m == nil {
		return
	}
	if _, ok := m.vals[k]; !ok {
		return
	}
	delete(m.vals, k)
	for i, key := range m.keys {
		if key == k {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

// Clone performs a deep copy.
func (m *Map) Clone() *Map {
	if m == nil {
		return nil
	}
	c := &Map{keys: append([]string(nil), m.keys...), vals: make(map[string]Value, len(m.vals)), ForceObject: m.ForceObject}
	for k, v := range m.vals {
		c.vals[k] = CloneValue(v)
	}
	return c
}

// SortKeys reorders keys with less (PHP ksort semantics are provided by KSort).
func (m *Map) SortKeys(less func(a, b string) bool) {
	keys := m.keys
	// insertion sort keeps it stable without importing sort for small maps
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && less(keys[j], keys[j-1]); j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
}

// KSort sorts keys like PHP 8 ksort() with default flags: numeric keys compare
// numerically, numeric vs non-numeric and string vs string use PHP 8 rules.
func (m *Map) KSort() { m.SortKeys(func(a, b string) bool { return PHPCompareKeys(a, b) < 0 }) }

// PHPCompareKeys compares two array keys the way PHP 8 ksort() does.
func PHPCompareKeys(a, b string) int {
	ai, aInt := IntKey(a)
	bi, bInt := IntKey(b)
	switch {
	case aInt && bInt:
		return cmpInt(ai, bi)
	case aInt && !bInt:
		// int vs non-numeric string: PHP 8 compares as strings unless numeric
		if f, ok := numericString(b); ok {
			return cmpFloat(float64(ai), f)
		}
		return strings.Compare(a, b)
	case !aInt && bInt:
		if f, ok := numericString(a); ok {
			return cmpFloat(f, float64(bi))
		}
		return strings.Compare(a, b)
	default:
		fa, okA := numericString(a)
		fb, okB := numericString(b)
		if okA && okB {
			return cmpFloat(fa, fb)
		}
		return strings.Compare(a, b)
	}
}

func cmpInt(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func cmpFloat(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func numericString(s string) (float64, bool) {
	t := strings.TrimLeft(s, " \t\n\r\v\f")
	t = strings.TrimRight(t, " \t\n\r\v\f")
	if t == "" {
		return 0, false
	}
	for _, c := range t {
		if !(c >= '0' && c <= '9' || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-') {
			return 0, false
		}
	}
	f, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// IntKey reports whether a PHP array key string would be stored as an int.
func IntKey(k string) (int64, bool) {
	if k == "" || len(k) > 20 {
		return 0, false
	}
	if k == "0" {
		return 0, true
	}
	s := k
	if s[0] == '-' {
		s = s[1:]
	}
	if s == "" || s[0] < '1' || s[0] > '9' {
		return 0, false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
	}
	n, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

func CloneValue(v Value) Value {
	switch t := v.(type) {
	case *Map:
		return t.Clone()
	case []Value:
		c := make([]Value, len(t))
		for i := range t {
			c[i] = CloneValue(t[i])
		}
		return c
	default:
		return v
	}
}

// IsList reports whether a map has keys "0".."n-1" in order (PHP array_is_list).
func (m *Map) IsList() bool {
	for i, k := range m.keys {
		if k != strconv.Itoa(i) {
			return false
		}
	}
	return true
}

// Flags for Encode, matching PHP JSON_* constants.
const (
	UnescapedSlashes = 64
	PrettyPrint      = 128
	UnescapedUnicode = 256
	// ComposerFlags is JsonFile::encode's default.
	ComposerFlags = UnescapedSlashes | PrettyPrint | UnescapedUnicode
)

// Encode mirrors PHP json_encode($v, $flags).
func Encode(v Value, flags int) (string, error) {
	var buf bytes.Buffer
	if err := encodeValue(&buf, v, flags, 0); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MustEncode panics on invalid UTF-8; callers use it for data they built.
func MustEncode(v Value, flags int) string {
	s, err := Encode(v, flags)
	if err != nil {
		panic(err)
	}
	return s
}

func indent(buf *bytes.Buffer, flags, depth int) {
	if flags&PrettyPrint != 0 {
		buf.WriteByte('\n')
		for i := 0; i < depth; i++ {
			buf.WriteString("    ")
		}
	}
}

func encodeValue(buf *bytes.Buffer, v Value, flags, depth int) error {
	switch t := v.(type) {
	case nil:
		buf.WriteString("null")
	case bool:
		if t {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case int:
		buf.WriteString(strconv.Itoa(t))
	case int64:
		buf.WriteString(strconv.FormatInt(t, 10))
	case float64:
		buf.WriteString(FormatFloat(t))
	case string:
		return encodeString(buf, t, flags)
	case []string:
		vals := make([]Value, len(t))
		for i, s := range t {
			vals[i] = s
		}
		return encodeValue(buf, vals, flags, depth)
	case []Value:
		if len(t) == 0 {
			buf.WriteString("[]")
			return nil
		}
		buf.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			indent(buf, flags, depth+1)
			if err := encodeValue(buf, e, flags, depth+1); err != nil {
				return err
			}
		}
		indent(buf, flags, depth)
		buf.WriteByte(']')
	case *Map:
		if t == nil || t.Len() == 0 {
			if t != nil && t.ForceObject {
				buf.WriteString("{}")
			} else {
				buf.WriteString("[]")
			}
			return nil
		}
		if !t.ForceObject && t.IsList() {
			vals := make([]Value, t.Len())
			for i, k := range t.keys {
				vals[i] = t.vals[k]
			}
			return encodeValue(buf, vals, flags, depth)
		}
		buf.WriteByte('{')
		for i, k := range t.keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			indent(buf, flags, depth+1)
			if err := encodeString(buf, k, flags); err != nil {
				return err
			}
			buf.WriteByte(':')
			if flags&PrettyPrint != 0 {
				buf.WriteByte(' ')
			}
			if err := encodeValue(buf, t.vals[k], flags, depth+1); err != nil {
				return err
			}
		}
		indent(buf, flags, depth)
		buf.WriteByte('}')
	default:
		return fmt.Errorf("phpjson: unsupported type %T", v)
	}
	return nil
}

// ErrInvalidUTF8 matches PHP's JSON_ERROR_UTF8.
var ErrInvalidUTF8 = errors.New("Malformed UTF-8 characters, possibly incorrectly encoded")

func encodeString(buf *bytes.Buffer, s string, flags int) error {
	buf.WriteByte('"')
	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			switch c {
			case '"':
				buf.WriteString(`\"`)
			case '\\':
				buf.WriteString(`\\`)
			case '/':
				if flags&UnescapedSlashes != 0 {
					buf.WriteByte('/')
				} else {
					buf.WriteString(`\/`)
				}
			case '\b':
				buf.WriteString(`\b`)
			case '\f':
				buf.WriteString(`\f`)
			case '\n':
				buf.WriteString(`\n`)
			case '\r':
				buf.WriteString(`\r`)
			case '\t':
				buf.WriteString(`\t`)
			default:
				if c < 0x20 {
					fmt.Fprintf(buf, `\u%04x`, c)
				} else {
					buf.WriteByte(c)
				}
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			return ErrInvalidUTF8
		}
		if flags&UnescapedUnicode != 0 && r != 0x2028 && r != 0x2029 {
			buf.WriteString(s[i : i+size])
		} else if r >= 0x10000 {
			r -= 0x10000
			fmt.Fprintf(buf, `\u%04x\u%04x`, 0xD800+(r>>10), 0xDC00+(r&0x3FF))
		} else {
			fmt.Fprintf(buf, `\u%04x`, r)
		}
		i += size
	}
	buf.WriteByte('"')
	return nil
}

// FormatFloat mirrors PHP's json_encode float output with serialize_precision=-1.
func FormatFloat(f float64) string { return phpGcvt(f, 'e', false) }

// FormatFloatVarExport formats like var_export() (uppercase exponent, ".0"
// kept on integral values).
func FormatFloatVarExport(f float64) string { return phpGcvt(f, 'E', true) }

// FormatFloatSerialize formats like serialize()'s "d:" values.
func FormatFloatSerialize(f float64) string { return phpGcvt(f, 'E', false) }

// phpGcvt is php_gcvt() with serialize_precision -1 (shortest round-trip
// digits, mode 0): exponential notation when the decimal exponent is below
// -3 or above 17, e.g. 1.0e+25 and 1.5e-7, otherwise plain digits.
func phpGcvt(f float64, expChar byte, zeroFrac bool) string {
	switch {
	case math.IsNaN(f):
		return "NAN"
	case math.IsInf(f, 1):
		return "INF"
	case math.IsInf(f, -1):
		return "-INF"
	}
	neg := math.Signbit(f)
	if neg {
		f = -f
	}
	sign := ""
	if neg {
		sign = "-"
	}
	if f == 0 {
		if zeroFrac {
			return sign + "0.0"
		}
		return sign + "0"
	}
	e := strconv.FormatFloat(f, 'e', -1, 64) // d.ddde±XX
	mant, exp, _ := strings.Cut(e, "e")
	digits := strings.Replace(mant, ".", "", 1)
	x, _ := strconv.Atoi(exp)
	decpt := x + 1
	var out string
	if decpt < -3 || decpt > 17 {
		frac := digits[1:]
		if frac == "" {
			frac = "0"
		}
		esign := "+"
		ev := decpt - 1
		if ev < 0 {
			esign = "-"
			ev = -ev
		}
		out = digits[:1] + "." + frac + string(expChar) + esign + strconv.Itoa(ev)
		return sign + out
	}
	switch {
	case decpt <= 0:
		out = "0." + strings.Repeat("0", -decpt) + digits
	case decpt >= len(digits):
		out = digits + strings.Repeat("0", decpt-len(digits))
		if zeroFrac {
			out += ".0"
		}
	default:
		out = digits[:decpt] + "." + digits[decpt:]
	}
	return sign + out
}

// SyntaxError reports a JSON parse failure with its byte offset and line.
type SyntaxError struct {
	Msg    string
	Offset int
	Line   int
}

func (e *SyntaxError) Error() string { return e.Msg }

// Decode parses JSON like json_decode($s, true) but keeps key order.
func Decode(data []byte) (Value, error) {
	p := &parser{data: data}
	p.skipWS()
	v, err := p.value(0)
	if err != nil {
		return nil, err
	}
	p.skipWS()
	if p.pos != len(p.data) {
		return nil, p.errorf("Syntax error")
	}
	return v, nil
}

type parser struct {
	data []byte
	pos  int
}

func (p *parser) errorf(format string, a ...any) error {
	line := 1 + bytes.Count(p.data[:min(p.pos, len(p.data))], []byte("\n"))
	return &SyntaxError{Msg: fmt.Sprintf(format, a...), Offset: p.pos, Line: line}
}

func (p *parser) skipWS() {
	for p.pos < len(p.data) {
		switch p.data[p.pos] {
		case ' ', '\t', '\n', '\r':
			p.pos++
		default:
			return
		}
	}
}

func (p *parser) value(depth int) (Value, error) {
	if depth > 512 {
		return nil, p.errorf("Maximum stack depth exceeded")
	}
	if p.pos >= len(p.data) {
		return nil, p.errorf("Syntax error")
	}
	switch c := p.data[p.pos]; {
	case c == '{':
		p.pos++
		m := NewMap()
		p.skipWS()
		if p.pos < len(p.data) && p.data[p.pos] == '}' {
			p.pos++
			return m, nil
		}
		for {
			p.skipWS()
			if p.pos >= len(p.data) || p.data[p.pos] != '"' {
				return nil, p.errorf("Syntax error")
			}
			k, err := p.str()
			if err != nil {
				return nil, err
			}
			p.skipWS()
			if p.pos >= len(p.data) || p.data[p.pos] != ':' {
				return nil, p.errorf("Syntax error")
			}
			p.pos++
			p.skipWS()
			v, err := p.value(depth + 1)
			if err != nil {
				return nil, err
			}
			m.Set(k, v)
			p.skipWS()
			if p.pos >= len(p.data) {
				return nil, p.errorf("Syntax error")
			}
			if p.data[p.pos] == ',' {
				p.pos++
				continue
			}
			if p.data[p.pos] == '}' {
				p.pos++
				return m, nil
			}
			return nil, p.errorf("Syntax error")
		}
	case c == '[':
		p.pos++
		list := []Value{}
		p.skipWS()
		if p.pos < len(p.data) && p.data[p.pos] == ']' {
			p.pos++
			return list, nil
		}
		for {
			p.skipWS()
			v, err := p.value(depth + 1)
			if err != nil {
				return nil, err
			}
			list = append(list, v)
			p.skipWS()
			if p.pos >= len(p.data) {
				return nil, p.errorf("Syntax error")
			}
			if p.data[p.pos] == ',' {
				p.pos++
				continue
			}
			if p.data[p.pos] == ']' {
				p.pos++
				return list, nil
			}
			return nil, p.errorf("Syntax error")
		}
	case c == '"':
		return p.str()
	case c == 't':
		return p.lit("true", true)
	case c == 'f':
		return p.lit("false", false)
	case c == 'n':
		return p.lit("null", nil)
	case c == '-' || (c >= '0' && c <= '9'):
		return p.number()
	default:
		return nil, p.errorf("Syntax error")
	}
}

func (p *parser) lit(word string, v Value) (Value, error) {
	if bytes.HasPrefix(p.data[p.pos:], []byte(word)) {
		p.pos += len(word)
		return v, nil
	}
	return nil, p.errorf("Syntax error")
}

func (p *parser) number() (Value, error) {
	start := p.pos
	if p.data[p.pos] == '-' {
		p.pos++
	}
	if p.pos >= len(p.data) {
		return nil, p.errorf("Syntax error")
	}
	if p.data[p.pos] == '0' {
		p.pos++
	} else if p.data[p.pos] >= '1' && p.data[p.pos] <= '9' {
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	} else {
		return nil, p.errorf("Syntax error")
	}
	isFloat := false
	if p.pos < len(p.data) && p.data[p.pos] == '.' {
		isFloat = true
		p.pos++
		d := p.pos
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
		if d == p.pos {
			return nil, p.errorf("Syntax error")
		}
	}
	if p.pos < len(p.data) && (p.data[p.pos] == 'e' || p.data[p.pos] == 'E') {
		isFloat = true
		p.pos++
		if p.pos < len(p.data) && (p.data[p.pos] == '+' || p.data[p.pos] == '-') {
			p.pos++
		}
		d := p.pos
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
		if d == p.pos {
			return nil, p.errorf("Syntax error")
		}
	}
	text := string(p.data[start:p.pos])
	if !isFloat {
		if n, err := strconv.ParseInt(text, 10, 64); err == nil {
			return n, nil
		}
	}
	f, err := strconv.ParseFloat(text, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return nil, p.errorf("Syntax error")
	}
	return f, nil
}

func (p *parser) str() (string, error) {
	p.pos++ // opening quote
	var sb strings.Builder
	for {
		if p.pos >= len(p.data) {
			return "", p.errorf("Syntax error")
		}
		c := p.data[p.pos]
		switch {
		case c == '"':
			p.pos++
			s := sb.String()
			if !utf8.ValidString(s) {
				return "", p.errorf("Malformed UTF-8 characters, possibly incorrectly encoded")
			}
			return s, nil
		case c == '\\':
			p.pos++
			if p.pos >= len(p.data) {
				return "", p.errorf("Syntax error")
			}
			e := p.data[p.pos]
			p.pos++
			switch e {
			case '"', '\\', '/':
				sb.WriteByte(e)
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'u':
				r, err := p.hex4()
				if err != nil {
					return "", err
				}
				if r >= 0xD800 && r < 0xDC00 {
					if p.pos+1 < len(p.data) && p.data[p.pos] == '\\' && p.data[p.pos+1] == 'u' {
						p.pos += 2
						r2, err := p.hex4()
						if err != nil {
							return "", err
						}
						if r2 >= 0xDC00 && r2 < 0xE000 {
							r = 0x10000 + (r-0xD800)<<10 + (r2 - 0xDC00)
						} else {
							return "", p.errorf("Single unpaired UTF-16 surrogate in unicode escape")
						}
					} else {
						return "", p.errorf("Single unpaired UTF-16 surrogate in unicode escape")
					}
				} else if r >= 0xDC00 && r < 0xE000 {
					return "", p.errorf("Single unpaired UTF-16 surrogate in unicode escape")
				}
				sb.WriteRune(r)
			default:
				return "", p.errorf("Syntax error")
			}
		case c < 0x20:
			return "", p.errorf("Control character error, possibly incorrectly encoded")
		default:
			sb.WriteByte(c)
			p.pos++
		}
	}
}

func (p *parser) hex4() (rune, error) {
	if p.pos+4 > len(p.data) {
		return 0, p.errorf("Syntax error")
	}
	n, err := strconv.ParseUint(string(p.data[p.pos:p.pos+4]), 16, 32)
	if err != nil {
		return 0, p.errorf("Syntax error")
	}
	p.pos += 4
	return rune(n), nil
}

// Accessors used across the code base. They never panic on type mismatch.

func AsMap(v Value) *Map {
	m, _ := v.(*Map)
	return m
}

func (m *Map) Map(k string) *Map {
	v, _ := m.Get(k)
	return AsMap(v)
}

func (m *Map) String(k string) (string, bool) {
	v, _ := m.Get(k)
	s, ok := v.(string)
	return s, ok
}

func (m *Map) Str(k string) string { s, _ := m.String(k); return s }

// ToMap converts a list into a map keyed "0".."n-1", or returns the map.
func ToMap(v Value) *Map {
	switch t := v.(type) {
	case *Map:
		return t
	case []Value:
		m := NewMap()
		for i, e := range t {
			m.Set(strconv.Itoa(i), e)
		}
		return m
	}
	return nil
}

// IsArray reports whether v is a PHP array (list or map).
func IsArray(v Value) bool {
	switch t := v.(type) {
	case *Map:
		return t != nil
	case []Value:
		return true
	}
	return false
}

// ArrayLen returns count() of a PHP array value.
func ArrayLen(v Value) int {
	switch t := v.(type) {
	case *Map:
		return t.Len()
	case []Value:
		return len(t)
	}
	return 0
}

// Each iterates a PHP array value (list or map) in order.
func Each(v Value, fn func(k string, v Value)) {
	switch t := v.(type) {
	case *Map:
		if t == nil {
			return
		}
		for _, k := range t.keys {
			fn(k, t.vals[k])
		}
	case []Value:
		for i, e := range t {
			fn(strconv.Itoa(i), e)
		}
	}
}

// Truthy mirrors PHP's boolean conversion (used for !empty()).
func Truthy(v Value) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case int64:
		return t != 0
	case int:
		return t != 0
	case float64:
		return t != 0
	case string:
		return t != "" && t != "0"
	case []Value:
		return len(t) > 0
	case *Map:
		return t.Len() > 0
	}
	return true
}

// ScalarString converts a scalar like PHP (string) cast does.
func ScalarString(v Value) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case int64:
		return strconv.FormatInt(t, 10), true
	case int:
		return strconv.Itoa(t), true
	case float64:
		if t == math.Trunc(t) && math.Abs(t) < 1e15 {
			return strconv.FormatFloat(t, 'f', -1, 64), true
		}
		return strconv.FormatFloat(t, 'G', 14, 64), true
	case bool:
		if t {
			return "1", true
		}
		return "", true
	}
	return "", false
}
