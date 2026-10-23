<?php
// Runs a corpus of JSON schema validation cases through the pinned
// justinrainbow/json-schema of the oracle Composer checkout, exactly like
// Composer\Json\JsonFile::validateJsonSchema does, and prints the raw error
// lists as JSON. Usage: php jsonschema-dump.php corpus.json
//
// corpus.json: [{"name": ..., "schema": <inline schema> | "composer" |
//   "composer-strict" | "composer-auth" | "lock", "data": <json>,
//   "assoc": true|false (decode data as assoc array first, as callers that
//   pass arrays do), "mode": <check mode, optional>}]
require '/oracle-source/vendor/autoload.php';

use JsonSchema\Validator;

// Warnings (e.g. preg_match() compilation failures) are swallowed: the
// library itself does not handle them; Composer's own ErrorHandler is not
// part of the library behaviour being compared.
set_error_handler(static function () { return true; });

$composerSchema = 'file:///oracle-source/res/composer-schema.json';
$lockSchema = 'file:///oracle-source/res/composer-lock-schema.json';

$cases = json_decode(file_get_contents($argv[1]));
$out = [];
foreach ($cases as $case) {
    $data = $case->data;
    if (!empty($case->assoc)) {
        $data = json_decode(json_encode($data), true);
    }
    $schema = $case->schema;
    if (is_string($schema)) {
        switch ($schema) {
            case 'composer':
                $schema = (object) ['$ref' => $composerSchema, '$schema' => 'https://json-schema.org/draft-04/schema#'];
                break;
            case 'composer-strict':
                $schema = json_decode((string) file_get_contents($composerSchema));
                $schema->additionalProperties = false;
                $schema->required = ['name', 'description'];
                break;
            case 'composer-auth':
                $schema = (object) ['$ref' => $composerSchema.'#/properties/config', '$schema' => 'https://json-schema.org/draft-04/schema#'];
                break;
            case 'lock':
                $schema = (object) ['$ref' => $lockSchema, '$schema' => 'https://json-schema.org/draft-04/schema#'];
                break;
            default:
                fwrite(STDERR, "unknown schema $schema\n");
                exit(1);
        }
    }
    $row = ['name' => $case->name];
    try {
        $validator = new Validator();
        // convert assoc arrays to objects (JsonFile::validateJsonSchema)
        $data = json_decode((string) json_encode($data));
        if (isset($case->mode)) {
            $validator->validate($data, $schema, $case->mode);
        } else {
            $validator->validate($data, $schema);
        }
        $errors = [];
        $formatted = [];
        foreach ($validator->getErrors() as $error) {
            $errors[] = [
                'property' => $error['property'],
                'pointer' => $error['pointer'],
                'message' => $error['message'],
                'constraint' => $error['constraint']['name'],
                'params' => json_encode($error['constraint']['params']),
            ];
            $formatted[] = ($error['property'] ? $error['property'].' : ' : '').$error['message'];
        }
        $row['errors'] = $errors;
        $row['formatted'] = $formatted;
    } catch (\Throwable $e) {
        $row['exception'] = ['class' => get_class($e), 'message' => $e->getMessage()];
    }
    $out[] = $row;
}
echo json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_PRESERVE_ZERO_FRACTION), "\n";
