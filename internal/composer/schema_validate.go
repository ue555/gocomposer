package composer

import (
	"os"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/validate"
)

// ValidateSchemaFile is JsonFile::validateSchema($schema) for a file path:
// syntax errors are returned as the jsonlint parse error, schema errors as
// *validate.JSONValidationError.
func ValidateSchemaFile(path string, schema int) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return &UserError{Msg: `The file "` + path + `" is not readable.`, Code: 1}
	}
	data, err := jsonfile.Parse(raw, path)
	if err != nil {
		return err
	}
	return validateSchemaData(path, data, schema)
}

func validateSchemaData(source string, data phpjson.Value, schema int) error {
	msgs, err := jsonschema.ValidateComposerJSON(data, schema)
	if err != nil {
		return err
	}
	if len(msgs) > 0 {
		return &validate.JSONValidationError{Msg: `"` + source + `" does not match the expected JSON schema`, Errors: msgs}
	}
	return nil
}

// validateComposerSchema is Factory::createComposer's LAX schema check.
func validateComposerSchema(data *phpjson.Map, file string, strict bool) error {
	mode := jsonschema.LaxSchema
	if strict {
		mode = jsonschema.StrictSchema
	}
	err := validateSchemaData(file, data, mode)
	if ve, ok := err.(*validate.JSONValidationError); ok {
		return &validate.JSONValidationError{Msg: ve.Msg + ":\n - " + strings.Join(ve.Errors, "\n - ")}
	}
	return err
}
