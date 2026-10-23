package composer

import "gocomposer/internal/phpjson"

// ValidateSchemaLax is JsonFile::validateSchema(LAX_SCHEMA). The JSON schema
// validator is implemented in internal/jsonschema; see schema_validate.go.
func ValidateSchemaLax(data *phpjson.Map, file string) error {
	return validateComposerSchema(data, file, false)
}
