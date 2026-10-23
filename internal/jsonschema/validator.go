// Package jsonschema is a port of justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016 justinrainbow/json-schema contributors), the JSON
// schema validator Composer uses, restricted to the code path Composer
// exercises: Validator::validate() with the default check mode (draft-03/04
// keywords handled by the legacy Constraints classes). It reproduces the
// library's error list exactly: order, property paths, JSON pointers,
// messages, constraint names and parameters, including its quirks (duplicate
// removal with array_unique, the oneOf error reordering, numeric strings
// being range-checked, ...).
//
// This file ports src/JsonSchema/Validator.php.
//
// Value model: documents and schemas are the values json_decode($json)
// returns, expressed with gocomposer/internal/phpjson: *phpjson.Map is a
// stdClass object (also when empty, so "{}" stays an object), []phpjson.Value
// is a PHP array (JSON array, "[]"), plus int64, float64, string, bool and nil.
// Use FromAssoc for data that has PHP associative-array semantics (empty maps
// become [] and maps with keys "0".."n-1" become lists, exactly like
// json_decode(json_encode($assoc))).
package jsonschema

import (
	"errors"
	"fmt"

	"gocomposer/internal/phpjson"
)

// Check modes (Constraint::CHECK_MODE_*).
const (
	CheckModeNormal               = 0x00000001
	CheckModeTypeCast             = 0x00000002
	CheckModeCoerceTypes          = 0x00000004
	CheckModeApplyDefaults        = 0x00000008
	CheckModeExceptions           = 0x00000010
	CheckModeDisableFormat        = 0x00000020
	CheckModeEarlyCoerce          = 0x00000040
	CheckModeOnlyRequiredDefaults = 0x00000080
	CheckModeValidateSchema       = 0x00000100
	CheckModeStrict               = 0x00000200
)

// ErrUnsupportedCheckMode is returned for check modes that are not ported
// (type casting, coercion, defaults, exceptions, schema validation, strict
// dialect mode). Composer never uses them.
var ErrUnsupportedCheckMode = errors.New("jsonschema: unsupported check mode")

// Validator mirrors JsonSchema\Validator. Each Validate call uses a fresh
// SchemaStorage, like Composer which creates a new Validator per validation.
type Validator struct {
	// CheckMode is a CheckMode* bitmask; 0 means CheckModeNormal.
	CheckMode int
	// Sources maps retrievable URIs (without fragment, e.g.
	// "file:///res/schema.json") to raw JSON documents. Other file:// URIs
	// are read from disk; other schemes fail.
	Sources map[string][]byte
}

// Validate validates data against schema and returns Validator::getErrors()
// (after the final array_unique). A non-nil error is an *Exception that the
// PHP library would have thrown, or ErrUnsupportedCheckMode. The schema value
// is not modified.
func (vd *Validator) Validate(data, schema Value) (errs []Error, err error) {
	mode := vd.CheckMode
	if mode == 0 {
		mode = CheckModeNormal
	}
	if mode&^(CheckModeNormal|CheckModeDisableFormat) != 0 {
		return nil, fmt.Errorf("%w: 0x%x", ErrUnsupportedCheckMode, mode)
	}
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*Exception); ok {
				errs, err = nil, e
				return
			}
			panic(r)
		}
	}()
	v := &validator{st: newSchemaStorage(vd.Sources), disableFormat: mode&CheckModeDisableFormat != 0}
	schema = phpjson.CloneValue(normalizeDeep(schema))
	data = normalizeDeep(data)

	schemaURI := internalProvidedSchemaURI
	for _, key := range []string{"id", "$id"} {
		if idv, ok := mapGet(schema, key); ok {
			s, isStr := idv.(string)
			if !isStr {
				throw("TypeError", "JsonSchema\\SchemaStorage::addSchema(): Argument #1 ($id) must be of type string, "+debugType(idv)+" given")
			}
			schemaURI = s
		}
	}
	v.st.addSchema(schemaURI, schema)
	sch := v.st.getSchema(schemaURI)

	top := &inst{}
	if b, ok := sch.(bool); ok {
		if !b {
			v.addError(top, "false", jsonPointer{})
		}
		return top.errs, nil
	}
	// SchemaConstraint::check
	validationSchema := sch
	if validationSchema == nil {
		switch data.(type) {
		case *phpjson.Map, *undefinedT:
			if s, ok := mapGet(data, "$schema"); ok {
				validationSchema = s
			}
		case string:
		default:
			throw("TypeError", "property_exists(): Argument #1 ($object_or_class) must be of type object|string, "+debugType(data)+" given")
		}
		if validationSchema == nil {
			throw("JsonSchema\\Exception\\InvalidArgumentException", "no schema found to verify against")
		}
	}
	if a, ok := validationSchema.([]Value); ok {
		validationSchema = arrayToObject(a)
	}
	v.checkUndefined(top, data, validationSchema, jsonPointer{}, nil)
	return arrayUnique(top.errs), nil
}

// arrayUnique mirrors array_unique($errors, SORT_REGULAR): the first of each
// group of equal errors is kept, in original order.
func arrayUnique(errs []Error) []Error {
	seen := map[string]bool{}
	var out []Error
	for _, e := range errs {
		key := e.Property + "\x00" + e.Pointer + "\x00" + e.Message + "\x00" + e.Constraint + "\x00" + e.ParamsJSON()
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, e)
	}
	return out
}

func normalizeDeep(v Value) Value {
	switch t := v.(type) {
	case int, int32, float32:
		return normalize(t)
	case []Value:
		for i := range t {
			t[i] = normalizeDeep(t[i])
		}
	case *phpjson.Map:
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			switch x.(type) {
			case int, int32, float32:
				t.Set(k, normalize(x))
			default:
				normalizeDeep(x)
			}
		}
	}
	return v
}

// Schema is a loaded schema document.
type Schema struct {
	root    Value
	baseURI string
	sources map[string][]byte
}

// LoadSchema decodes a schema document. With a non-empty baseURI the document
// is registered under that URI and validation uses {"$ref": baseURI}, the way
// Composer validates against schema files (relative "$ref"s then resolve
// against baseURI). Without it the schema is used as the provided schema
// (internal://provided-schema/).
func LoadSchema(raw []byte, baseURI string) (*Schema, error) {
	v, err := phpjson.Decode(raw)
	if err != nil {
		return nil, err
	}
	s := &Schema{root: v, baseURI: baseURI, sources: map[string][]byte{}}
	if baseURI != "" {
		s.sources[baseURI] = raw
		ref := phpjson.NewMap()
		ref.Set("$ref", baseURI)
		s.root = ref
	}
	return s, nil
}

// AddSource registers another document retrievable by URI (for "$ref"s to
// other files).
func (s *Schema) AddSource(uri string, raw []byte) { s.sources[uri] = raw }

// Validate validates data in the default check mode.
func (s *Schema) Validate(data Value) ([]Error, error) {
	vd := &Validator{Sources: s.sources}
	return vd.Validate(data, s.root)
}

// FromAssoc mirrors json_decode(json_encode($data)) for data with PHP
// associative-array semantics: empty maps become [] and list-like maps become
// arrays. Encoding failures (invalid UTF-8) yield nil, like PHP's
// json_decode(”).
func FromAssoc(data Value) Value {
	return jsonRoundTrip(assocToObjects(data))
}

// JSONRoundTrip mirrors json_decode((string) json_encode($data)) for data
// already in object form (*phpjson.Map = stdClass): integral floats become
// integers (json_encode(1.0) is "1"), invalid UTF-8 yields nil.
func JSONRoundTrip(data Value) Value { return jsonRoundTrip(data) }
