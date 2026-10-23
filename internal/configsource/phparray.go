// PHP array semantics for the whole-file fallbacks of
// src/Composer/Config/JsonConfigSource.php (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80): nested writes with
// auto-vivification, unset, append, array_values/array_splice, isset.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package configsource

import (
	"strconv"

	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/phpjson"
)

func phpError(class, msg string) { panic(&jsonmanip.Exception{Class: class, Message: msg}) }

// normalize converts a json_decode($x, true) value so every PHP array is an
// ordered *phpjson.Map (lists keyed "0".."n-1").
func normalize(v phpjson.Value) phpjson.Value {
	switch t := v.(type) {
	case *phpjson.Map:
		out := phpjson.NewMap()
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			out.Set(k, normalize(x))
		}
		out.ForceObject = t.ForceObject
		return out
	case []phpjson.Value:
		out := phpjson.NewMap()
		for i, x := range t {
			out.Set(strconv.Itoa(i), normalize(x))
		}
		return out
	}
	return v
}

func isArr(v phpjson.Value) bool {
	m, ok := v.(*phpjson.Map)
	return ok && m != nil && !m.ForceObject
}

func arr(v phpjson.Value) *phpjson.Map {
	if isArr(v) {
		return v.(*phpjson.Map)
	}
	return nil
}

func isStringKey(k string) bool {
	_, isInt := phpjson.IntKey(k)
	return !isInt
}

// isset mirrors isset($v[$k1][$k2]...) (string offsets included).
func isset(v phpjson.Value, keys ...string) bool {
	cur := v
	for _, k := range keys {
		switch t := cur.(type) {
		case *phpjson.Map:
			if t == nil || t.ForceObject {
				return false
			}
			x, ok := t.Get(k)
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

// get is "$v[$k1]... ?? null" for arrays.
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
		cur, _ = cur.(*phpjson.Map).Get(k)
	}
	return cur
}

// slot returns $parent[$key] as an array for writing, creating or
// converting it the way PHP auto-vivifies (null/missing/false become []).
func slot(parent *phpjson.Map, key string) *phpjson.Map {
	x, ok := parent.Get(key)
	if !ok || x == nil {
		m := phpjson.NewMap()
		parent.Set(key, m)
		return m
	}
	switch t := x.(type) {
	case bool:
		if !t {
			m := phpjson.NewMap()
			parent.Set(key, m)
			return m
		}
	case *phpjson.Map:
		if !t.ForceObject {
			return t
		}
		phpError("Error", "Cannot use object of type stdClass as array")
	case string:
		phpError("TypeError", "Cannot access offset of type string on string")
	}
	phpError("Error", "Cannot use a scalar value as an array")
	return nil
}

// setPath is $root[$k1]...[$kn] = $val.
func setPath(root *phpjson.Map, val phpjson.Value, keys ...string) {
	cur := root
	for _, k := range keys[:len(keys)-1] {
		cur = slot(cur, k)
	}
	cur.Set(keys[len(keys)-1], val)
}

// unsetPath is unset($root[$k1]...[$kn]).
func unsetPath(root *phpjson.Map, keys ...string) {
	var cur phpjson.Value = root
	for i, k := range keys {
		switch t := cur.(type) {
		case nil:
			return
		case bool:
			if !t {
				return
			}
			phpError("Error", "Cannot unset offset in a non-array variable")
		case string:
			phpError("Error", "Cannot unset string offsets")
		case *phpjson.Map:
			if t.ForceObject {
				phpError("Error", "Cannot use object of type stdClass as array")
			}
			if i == len(keys)-1 {
				t.Delete(k)
				return
			}
			x, exists := t.Get(k)
			if !exists {
				return
			}
			cur = x
		default:
			phpError("Error", "Cannot unset offset in a non-array variable")
		}
	}
}

// push is $arr[] = $val.
func push(a *phpjson.Map, val phpjson.Value) {
	next := int64(0)
	for _, k := range a.Keys() {
		if i, ok := phpjson.IntKey(k); ok && i >= next {
			next = i + 1
		}
	}
	a.Set(strconv.FormatInt(next, 10), val)
}

// values is array_values().
func values(v phpjson.Value) []phpjson.Value {
	var out []phpjson.Value
	phpjson.Each(v, func(_ string, x phpjson.Value) { out = append(out, x) })
	return out
}

func listMap(vals []phpjson.Value) *phpjson.Map {
	m := phpjson.NewMap()
	for i, v := range vals {
		m.Set(strconv.Itoa(i), v)
	}
	return m
}

func mapOf(k string, v phpjson.Value) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set(k, v)
	return m
}

// prependName is ['name' => $name] + $array.
func prependName(name string, a *phpjson.Map) *phpjson.Map {
	out := phpjson.NewMap()
	out.Set("name", name)
	for _, k := range a.Keys() {
		if k != "name" {
			x, _ := a.Get(k)
			out.Set(k, x)
		}
	}
	return out
}

// isFalseMap is ($v === [$k => false]).
func isFalseMap(v phpjson.Value, k string) bool {
	m := arr(v)
	if m == nil || m.Len() != 1 {
		return false
	}
	x, ok := m.Get(k)
	b, isBool := x.(bool)
	return ok && isBool && !b
}

func isEmptyArray(v phpjson.Value) bool {
	m := arr(v)
	return m != nil && m.Len() == 0
}

func typeName(v phpjson.Value) string {
	switch t := v.(type) {
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
	case *phpjson.Map:
		if t.ForceObject {
			return "stdClass"
		}
	}
	return "array"
}

// requireArrayArg mirrors the TypeError of internal functions taking arrays.
func requireArrayArg(fn string, n int, name string, v phpjson.Value) *phpjson.Map {
	if !isArr(v) {
		phpError("TypeError", fn+"(): Argument #"+strconv.Itoa(n)+" ($"+name+") must be of type array, "+typeName(v)+" given")
	}
	return v.(*phpjson.Map)
}

// filterRepos is array_values(array_filter($repos ?? [], fn)) with the
// uniqueness filter of JsonConfigSource (which keeps every entry except
// one equal to [$name => false] that also has a 'name' key equal to $name).
func filterRepos(repos phpjson.Value, name string) *phpjson.Map {
	if repos == nil {
		repos = phpjson.NewMap()
	}
	a := requireArrayArg("array_filter", 1, "array", repos)
	var kept []phpjson.Value
	for _, k := range a.Keys() {
		val, _ := a.Get(k)
		n, hasName := get(val, "name").(string)
		keep := !isset(val, "name") || !hasName || n != name || !isFalseMap(val, name)
		if keep {
			kept = append(kept, val)
		}
	}
	return listMap(kept)
}
