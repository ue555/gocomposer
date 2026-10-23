// Ported for gocomposer (MIT) from justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016 justinrainbow/json-schema contributors):
// src/JsonSchema/Tool/DeepComparer.php, the type checks of
// src/JsonSchema/Constraints/TypeCheck/{Strict,Loose}TypeCheck.php, and the
// PHP 8.4 runtime semantics the library depends on (gettype, is_numeric,
// loose comparison, string/int conversion in vsprintf, json_encode).
package jsonschema

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"gocomposer/internal/phpjson"
)

// Value is a decoded JSON value as produced by json_decode($json) (objects as
// stdClass): nil, bool, int64, float64, string, []phpjson.Value (PHP array)
// or *phpjson.Map (stdClass object).
type Value = phpjson.Value

// undefinedT stands for the UndefinedConstraint instance the library passes
// as the value of a missing property. It is a PHP object (is_object() true).
type undefinedT struct{}

var undefined = &undefinedT{}

func isUndef(v Value) bool { _, ok := v.(*undefinedT); return ok }

// Exception is a PHP exception (or Error) that would escape
// Validator::validate(); Class is the fully qualified PHP class name.
type Exception struct {
	Class   string
	Message string
}

func (e *Exception) Error() string { return e.Class + ": " + e.Message }

func throw(class, msg string) { panic(&Exception{Class: class, Message: msg}) }

func normalize(v Value) Value {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float32:
		return float64(t)
	}
	return v
}

// gettype mirrors PHP gettype().
func gettype(v Value) string {
	switch normalize(v).(type) {
	case nil:
		return "NULL"
	case bool:
		return "boolean"
	case int64:
		return "integer"
	case float64:
		return "double"
	case string:
		return "string"
	case []Value:
		return "array"
	}
	return "object"
}

// typeName is used in TypeError messages (get_debug_type style).
func debugType(v Value) string {
	switch t := normalize(v).(type) {
	case nil:
		return "null"
	case bool:
		if t {
			return "true"
		}
		return "false"
	case int64:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case []Value:
		return "array"
	case *undefinedT:
		return "JsonSchema\\Constraints\\UndefinedConstraint"
	}
	return "stdClass"
}

func isObject(v Value) bool {
	switch v.(type) {
	case *phpjson.Map, *undefinedT:
		return true
	}
	return false
}

func isArray(v Value) bool { _, ok := v.([]Value); return ok }

// looseIsObject mirrors LooseTypeCheck::isObject: objects and empty or
// associative arrays (decoded JSON arrays are lists, so only empty ones).
func looseIsObject(v Value) bool {
	if isObject(v) {
		return true
	}
	if a, ok := v.([]Value); ok {
		return len(a) == 0
	}
	return false
}

func isScalar(v Value) bool {
	switch normalize(v).(type) {
	case bool, int64, float64, string:
		return true
	}
	return false
}

// truthy mirrors PHP's boolean conversion.
func truthy(v Value) bool {
	switch t := normalize(v).(type) {
	case nil:
		return false
	case bool:
		return t
	case int64:
		return t != 0
	case float64:
		return t != 0
	case string:
		return t != "" && t != "0"
	case []Value:
		return len(t) > 0
	}
	return true
}

// mapGet returns $obj->$k when the property exists.
func mapGet(v Value, k string) (Value, bool) {
	m, ok := v.(*phpjson.Map)
	if !ok {
		return nil, false
	}
	return m.Get(k)
}

// isset mirrors isset($obj->$k).
func isset(v Value, k string) bool {
	x, ok := mapGet(v, k)
	return ok && x != nil
}

func get(v Value, k string) Value { x, _ := mapGet(v, k); return x }

// iterValues mirrors foreach ($v as $x) for arrays and objects.
func iterValues(v Value) []Value {
	switch t := v.(type) {
	case []Value:
		return t
	case *phpjson.Map:
		out := make([]Value, 0, t.Len())
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			out = append(out, x)
		}
		return out
	}
	return nil
}

type kv struct {
	k Value // string or int64 key
	v Value
}

// iterPairs mirrors foreach ($v as $k => $x).
func iterPairs(v Value) []kv {
	switch t := v.(type) {
	case []Value:
		out := make([]kv, len(t))
		for i, x := range t {
			out[i] = kv{int64(i), x}
		}
		return out
	case *phpjson.Map:
		out := make([]kv, 0, t.Len())
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			out = append(out, kv{k, x})
		}
		return out
	}
	return nil
}

// arrayKey converts a PHP array key candidate to its canonical form.
func arrayKey(k Value) (int64, bool) {
	switch t := normalize(k).(type) {
	case int64:
		return t, true
	case string:
		return phpjson.IntKey(t)
	case bool:
		if t {
			return 1, true
		}
		return 0, true
	case float64:
		return int64(t), true
	}
	return 0, false
}

// ---- numeric strings ----

func isWS(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

// numericPrefix parses the longest numeric prefix (PHP's
// _is_numeric_string_ex). It returns the parsed value, whether the whole
// string (allowing surrounding whitespace) is numeric, and whether any
// numeric prefix exists.
func numericPrefix(s string) (val Value, whole bool, prefix bool) {
	i := 0
	for i < len(s) && isWS(s[i]) {
		i++
	}
	start := i
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	digits := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
		digits++
	}
	isFloat := false
	frac := 0
	if i < len(s) && s[i] == '.' {
		j := i + 1
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
			frac++
		}
		if digits > 0 || frac > 0 {
			isFloat = true
			i = j
		}
	}
	if digits == 0 && frac == 0 {
		return nil, false, false
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		j := i + 1
		if j < len(s) && (s[j] == '+' || s[j] == '-') {
			j++
		}
		if j < len(s) && s[j] >= '0' && s[j] <= '9' {
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			isFloat = true
			i = j
		}
	}
	num := s[start:i]
	end := i
	for end < len(s) && isWS(s[end]) {
		end++
	}
	whole = end == len(s)
	if !isFloat {
		if n, err := strconv.ParseInt(num, 10, 64); err == nil {
			return n, whole, true
		}
	}
	f, _ := strconv.ParseFloat(num, 64)
	return f, whole, true
}

// isNumeric mirrors is_numeric().
func isNumeric(v Value) bool {
	switch t := normalize(v).(type) {
	case int64, float64:
		return true
	case string:
		_, whole, ok := numericPrefix(t)
		return ok && whole
	}
	return false
}

// toNumber returns the numeric value of an int, float or numeric string.
func toNumber(v Value) (Value, bool) {
	switch t := normalize(v).(type) {
	case int64, float64:
		return t, true
	case string:
		n, whole, ok := numericPrefix(t)
		if ok && whole {
			return n, true
		}
	}
	return nil, false
}

func toFloat(v Value) float64 {
	switch t := normalize(v).(type) {
	case int64:
		return float64(t)
	case float64:
		return t
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		n, _, ok := numericPrefix(t)
		if !ok {
			return 0
		}
		return toFloat(n)
	}
	return 0
}

// dvalToLval mirrors zend_dval_to_lval (modular for out-of-range values).
func dvalToLval(d float64) int64 {
	if math.IsInf(d, 0) || math.IsNaN(d) {
		return 0
	}
	if d >= -9.223372036854775808e18 && d < 9.223372036854775808e18 {
		return int64(d)
	}
	two64 := math.Pow(2, 64)
	dmod := math.Mod(d, two64)
	if dmod < 0 {
		if dmod == -math.Pow(2, 63) {
			return math.MinInt64
		}
		dmod += two64
	}
	if dmod > math.MaxInt64 {
		dmod -= two64
	}
	return int64(dmod)
}

// dvalToLvalCap mirrors zend_dval_to_lval_cap (used for numeric strings).
func dvalToLvalCap(d float64) int64 {
	if math.IsInf(d, 0) || math.IsNaN(d) {
		return 0
	}
	if d >= -9.223372036854775808e18 && d < 9.223372036854775808e18 {
		return int64(d)
	}
	if d > 0 {
		return math.MaxInt64
	}
	return math.MinInt64
}

// toInt mirrors zval_get_long() as used by sprintf('%d').
func toInt(v Value) int64 {
	switch t := normalize(v).(type) {
	case int64:
		return t
	case float64:
		return dvalToLval(t)
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		n, _, ok := numericPrefix(t)
		if !ok {
			return 0
		}
		if f, isF := n.(float64); isF {
			return dvalToLvalCap(f)
		}
		return n.(int64)
	case []Value:
		if len(t) > 0 {
			return 1
		}
		return 0
	}
	return 1
}

// phpFloatString mirrors PHP's float to string conversion (precision=14).
func phpFloatString(f float64) string { return gcvt(f, 14, 'E') }

// phpJSONFloat mirrors json_encode() of a float with serialize_precision=-1
// and no JSON_PRESERVE_ZERO_FRACTION (1.0 encodes as "1").
func phpJSONFloat(f float64) string { return gcvt(f, 0, 'e') }

// gcvt mirrors php_gcvt/zend_gcvt: ndigit significant digits (0 = shortest
// round-trip representation with the 17 digit exponent threshold).
func gcvt(f float64, ndigit int, expChar byte) string {
	if math.IsNaN(f) {
		return "NAN"
	}
	if math.IsInf(f, 1) {
		return "INF"
	}
	if math.IsInf(f, -1) {
		return "-INF"
	}
	if f == 0 {
		if math.Signbit(f) {
			return "-0"
		}
		return "0"
	}
	neg := f < 0
	if neg {
		f = -f
	}
	prec := ndigit - 1
	if ndigit == 0 {
		prec = -1
		ndigit = 17
	}
	e := strconv.FormatFloat(f, 'e', prec, 64)
	mant, exps, _ := strings.Cut(e, "e")
	exp, _ := strconv.Atoi(exps)
	digits := strings.TrimRight(strings.Replace(mant, ".", "", 1), "0")
	if digits == "" {
		digits = "0"
	}
	decpt := exp + 1
	var sb strings.Builder
	if neg {
		sb.WriteByte('-')
	}
	if (decpt < 0 && decpt < -3) || (decpt >= 0 && decpt > ndigit) {
		d := decpt - 1
		sign := '+'
		if d < 0 {
			sign = '-'
			d = -d
		}
		sb.WriteByte(digits[0])
		sb.WriteByte('.')
		if len(digits) == 1 {
			sb.WriteByte('0')
		} else {
			sb.WriteString(digits[1:])
		}
		sb.WriteByte(expChar)
		sb.WriteRune(sign)
		sb.WriteString(strconv.Itoa(d))
	} else if decpt < 0 {
		sb.WriteString("0.")
		for x := decpt; x < 0; x++ {
			sb.WriteByte('0')
		}
		sb.WriteString(digits)
	} else {
		for x := 0; x < decpt; x++ {
			if x < len(digits) {
				sb.WriteByte(digits[x])
			} else {
				sb.WriteByte('0')
			}
		}
		if len(digits) > decpt {
			if decpt == 0 {
				sb.WriteByte('0')
			}
			sb.WriteByte('.')
			sb.WriteString(digits[decpt:])
		}
	}
	return sb.String()
}

// toStr mirrors PHP's (string) conversion of scalars; arrays give "Array".
func toStr(v Value) string {
	switch t := normalize(v).(type) {
	case nil:
		return ""
	case bool:
		if t {
			return "1"
		}
		return ""
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return phpFloatString(t)
	case string:
		return t
	case []Value:
		return "Array"
	}
	throw("Error", "Object of class "+debugType(v)+" could not be converted to string")
	return ""
}

// ---- comparison ----

// looseEqual mirrors PHP 8 "==".
func looseEqual(a, b Value) bool { return compare(a, b) == 0 }

// compare mirrors PHP 8 "<=>" for the value kinds that occur in decoded JSON.
// Uncomparable pairs return 1.
func compare(a, b Value) int {
	a, b = normalize(a), normalize(b)
	switch x := a.(type) {
	case nil:
		switch y := b.(type) {
		case nil:
			return 0
		case bool:
			return cmpBool(false, y)
		case string:
			if y == "" {
				return 0
			}
			return -1
		case *phpjson.Map, *undefinedT:
			return -1
		}
		return cmpBool(false, truthy(b))
	case bool:
		return cmpBool(x, truthy(b))
	}
	switch y := b.(type) {
	case nil:
		return -compare(b, a)
	case bool:
		return cmpBool(truthy(a), y)
	}
	switch x := a.(type) {
	case int64, float64:
		switch y := b.(type) {
		case int64, float64:
			return cmpNum(x, y)
		case string:
			if n, ok := toNumber(y); ok {
				return cmpNum(x, n)
			}
			return strings.Compare(toStr(x), y)
		case []Value:
			return -1
		}
		return -1
	case string:
		switch y := b.(type) {
		case int64, float64:
			return -compare(b, a)
		case string:
			nx, okx := toNumber(x)
			ny, oky := toNumber(y)
			if okx && oky {
				return cmpNum(nx, ny)
			}
			return sign(strings.Compare(x, y))
		case []Value:
			return -1
		}
		return -1
	case []Value:
		y, ok := b.([]Value)
		if !ok {
			return 1
		}
		if len(x) != len(y) {
			return sign(len(x) - len(y))
		}
		for i := range x {
			if c := compare(x[i], y[i]); c != 0 {
				return c
			}
		}
		return 0
	case *phpjson.Map:
		y, ok := b.(*phpjson.Map)
		if !ok {
			return 1
		}
		if x == y {
			return 0
		}
		if x.Len() != y.Len() {
			return sign(x.Len() - y.Len())
		}
		for _, k := range x.Keys() {
			yv, ok := y.Get(k)
			if !ok {
				return 1
			}
			xv, _ := x.Get(k)
			if c := compare(xv, yv); c != 0 {
				return c
			}
		}
		return 0
	}
	if a == b {
		return 0
	}
	return 1
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}

func cmpBool(a, b bool) int {
	if a == b {
		return 0
	}
	if !a {
		return -1
	}
	return 1
}

func cmpNum(a, b Value) int {
	if x, ok := a.(int64); ok {
		if y, ok := b.(int64); ok {
			switch {
			case x < y:
				return -1
			case x > y:
				return 1
			}
			return 0
		}
	}
	fa, fb := toFloat(a), toFloat(b)
	switch {
	case fa < fb:
		return -1
	case fa > fb:
		return 1
	case fa == fb:
		return 0
	}
	return 1
}

// deepEqual mirrors JsonSchema\Tool\DeepComparer::isEqual.
func deepEqual(left, right Value) bool {
	left, right = normalize(left), normalize(right)
	if left == nil && right == nil {
		return true
	}
	ls, rs := isScalar(left), isScalar(right)
	if ls && rs {
		_, li := left.(int64)
		_, lf := left.(float64)
		_, ri := right.(int64)
		_, rf := right.(float64)
		if (li || lf) && (ri || rf) && toFloat(left) == toFloat(right) {
			return true
		}
		return left == right
	}
	if ls != rs {
		return false
	}
	if la, ok := left.([]Value); ok {
		if ra, ok := right.([]Value); ok {
			if len(la) != len(ra) {
				return false
			}
			for i := range la {
				if !deepEqual(la[i], ra[i]) {
					return false
				}
			}
			return true
		}
		return false
	}
	lm, ok1 := left.(*phpjson.Map)
	rm, ok2 := right.(*phpjson.Map)
	if ok1 && ok2 {
		// (array) casts turn numeric property names into integer keys on both
		// sides, so plain key comparison is equivalent.
		if lm.Len() != rm.Len() {
			return false
		}
		for _, k := range lm.Keys() {
			rv, ok := rm.Get(k)
			if !ok {
				return false
			}
			lv, _ := lm.Get(k)
			if !deepEqual(lv, rv) {
				return false
			}
		}
		return true
	}
	return false
}

// ---- json_encode ----

var errJSONEncode = errors.New("json_encode failed")

// jsonEncode mirrors json_encode($v) with default flags for values decoded
// as objects (every *phpjson.Map is a stdClass and encodes as {}). On failure
// (invalid UTF-8) it returns "" like (string) false.
func jsonEncode(v Value) (out string) {
	defer func() {
		if r := recover(); r != nil {
			if r != errJSONEncode {
				panic(r)
			}
			out = ""
		}
	}()
	var sb strings.Builder
	encodeObjMode(&sb, v)
	return sb.String()
}

// jsonRoundTrip mirrors json_decode((string) json_encode($v)) for a value
// whose objects are stdClass: integral floats such as 1.0 or 1e15 become
// integers, failures yield nil.
func jsonRoundTrip(v Value) Value {
	s := jsonEncode(v)
	if s == "" {
		return nil
	}
	out, err := phpjson.Decode([]byte(s))
	if err != nil {
		return nil
	}
	return out
}

// assocToObjects converts a value with PHP associative-array semantics into
// the json_decode(json_encode($v)) object form structure: empty maps and maps
// with keys "0".."n-1" become arrays (ForceObject maps stay objects).
func assocToObjects(v Value) Value {
	switch t := v.(type) {
	case []Value:
		out := make([]Value, len(t))
		for i, x := range t {
			out[i] = assocToObjects(x)
		}
		return out
	case *phpjson.Map:
		if t == nil {
			return nil
		}
		if !t.ForceObject && (t.Len() == 0 || t.IsList()) {
			out := make([]Value, 0, t.Len())
			for _, k := range t.Keys() {
				x, _ := t.Get(k)
				out = append(out, assocToObjects(x))
			}
			return out
		}
		m := phpjson.NewMap()
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			m.Set(k, assocToObjects(x))
		}
		return m
	}
	return normalize(v)
}

func encodeObjMode(sb *strings.Builder, v Value) {
	switch t := normalize(v).(type) {
	case nil:
		sb.WriteString("null")
	case bool:
		if t {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	case int64:
		sb.WriteString(strconv.FormatInt(t, 10))
	case float64:
		sb.WriteString(phpJSONFloat(t))
	case string:
		s, err := phpjson.Encode(t, 0)
		if err != nil {
			panic(errJSONEncode)
		}
		sb.WriteString(s)
	case []Value:
		sb.WriteByte('[')
		for i, x := range t {
			if i > 0 {
				sb.WriteByte(',')
			}
			encodeObjMode(sb, x)
		}
		sb.WriteByte(']')
	case *phpjson.Map:
		sb.WriteByte('{')
		for i, k := range t.Keys() {
			if i > 0 {
				sb.WriteByte(',')
			}
			ks, _ := phpjson.Encode(k, 0)
			sb.WriteString(ks)
			sb.WriteByte(':')
			x, _ := t.Get(k)
			encodeObjMode(sb, x)
		}
		sb.WriteByte('}')
	case *undefinedT:
		sb.WriteString("{}")
	}
}

// encodeParams encodes the constraint params like json_encode of the PHP
// associative array ($more): [] when empty, an object otherwise.
func encodeParams(params []param) string {
	if len(params) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, p := range params {
		if i > 0 {
			sb.WriteByte(',')
		}
		ks, _ := phpjson.Encode(p.Name, 0)
		sb.WriteString(ks)
		sb.WriteByte(':')
		encodeObjMode(&sb, p.Value)
	}
	sb.WriteByte('}')
	return sb.String()
}

// ucfirst mirrors PHP ucfirst() (ASCII only).
func ucfirst(s string) string {
	if s != "" && s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}

// vsprintf supports the %s, %d and %% conversions used by ConstraintError.
func vsprintf(format string, args []Value) string {
	var sb strings.Builder
	n := 0
	need := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) && format[i+1] != '%' {
			need++
		}
		if format[i] == '%' && i+1 < len(format) && format[i+1] == '%' {
			i++
		}
	}
	if need > len(args) {
		throw("ValueError", "The arguments array must contain "+strconv.Itoa(need)+" items, "+strconv.Itoa(len(args))+" given")
	}
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' || i+1 >= len(format) {
			sb.WriteByte(c)
			continue
		}
		i++
		switch format[i] {
		case '%':
			sb.WriteByte('%')
		case 's':
			sb.WriteString(toStr(args[n]))
			n++
		case 'd':
			sb.WriteString(strconv.FormatInt(toInt(args[n]), 10))
			n++
		default:
			sb.WriteByte('%')
			sb.WriteByte(format[i])
		}
	}
	return sb.String()
}
