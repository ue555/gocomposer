// Ported for gocomposer from Composer (MIT, Copyright (c) Nils Adermann,
// Jordi Boggiano) at f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Json/JsonFile.php (validateJsonSchema). The embedded schemas
// res/composer-schema.json and res/composer-lock-schema.json are verbatim
// copies of Composer's res/ files at the same commit (see res/NOTICE).
// Validation itself uses the port of justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016) in this package.
package jsonschema

import (
	_ "embed"
	"fmt"

	"gocomposer/internal/phpjson"
)

// JsonFile::*_SCHEMA constants.
const (
	LaxSchema    = 1
	StrictSchema = 2
	AuthSchema   = 3
	LockSchema   = 4
)

// URIs under which the embedded schemas are registered. Composer uses
// file:// URIs of its res/ directory; the path never shows up in errors.
const (
	ComposerSchemaURI = "file:///composer/res/composer-schema.json"
	LockSchemaURI     = "file:///composer/res/composer-lock-schema.json"
)

//go:embed res/composer-schema.json
var ComposerSchemaJSON []byte

//go:embed res/composer-lock-schema.json
var LockSchemaJSON []byte

func composerSources() map[string][]byte {
	return map[string][]byte{ComposerSchemaURI: ComposerSchemaJSON, LockSchemaURI: LockSchemaJSON}
}

// ComposerSchemaData returns the schema value JsonFile::validateJsonSchema
// passes to Validator::validate() for the given *_SCHEMA mode (with the
// default schema files).
func ComposerSchemaData(mode int) (Value, error) {
	ref := func(uri string) Value {
		m := phpjson.NewMap()
		m.Set("$ref", uri)
		m.Set("$schema", "https://json-schema.org/draft-04/schema#")
		return m
	}
	switch mode {
	case LaxSchema:
		return ref(ComposerSchemaURI), nil
	case StrictSchema:
		v, err := phpjson.Decode(ComposerSchemaJSON)
		if err != nil {
			return nil, err
		}
		m := phpjson.AsMap(v)
		m.Set("additionalProperties", false)
		m.Set("required", []phpjson.Value{"name", "description"})
		return m, nil
	case AuthSchema:
		return ref(ComposerSchemaURI + "#/properties/config"), nil
	case LockSchema:
		return ref(LockSchemaURI), nil
	}
	return nil, fmt.Errorf("jsonschema: unknown schema mode %d", mode)
}

// ValidateComposerSchema validates data (decoded with objects as
// *phpjson.Map, as JsonFile::validateSchema does with json_decode($content))
// the way JsonFile::validateJsonSchema($source, $data, $mode) does and returns
// the raw error list. Like JsonFile it first applies
// json_decode(json_encode($data)) (so 1.0 becomes 1). Data with PHP
// associative-array semantics (empty map = []) must go through FromAssoc
// first.
func ValidateComposerSchema(data Value, mode int) ([]Error, error) {
	schema, err := ComposerSchemaData(mode)
	if err != nil {
		return nil, err
	}
	vd := &Validator{Sources: composerSources()}
	// convert assoc arrays to objects: json_decode((string) json_encode($data))
	return vd.Validate(jsonRoundTrip(data), schema)
}

// FormatError formats an error like JsonFile::validateJsonSchema:
// "property : message", or just the message when the property is empty.
func FormatError(e Error) string {
	if truthy(e.Property) {
		return e.Property + " : " + e.Message
	}
	return e.Message
}

// FormatErrors formats every error with FormatError.
func FormatErrors(errs []Error) []string {
	out := make([]string, 0, len(errs))
	for _, e := range errs {
		out = append(out, FormatError(e))
	}
	return out
}

// ValidateComposerJSON returns the error strings JsonValidationException
// would carry (nil when the document is valid).
func ValidateComposerJSON(data Value, mode int) ([]string, error) {
	errs, err := ValidateComposerSchema(data, mode)
	if err != nil {
		return nil, err
	}
	if len(errs) == 0 {
		return nil, nil
	}
	return FormatErrors(errs), nil
}
