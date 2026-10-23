<?php
// Lints each corpus string with the seld/jsonlint version locked by the
// oracle Composer checkout (default flags and DETECT_KEY_CONFLICTS|PARSE_TO_ASSOC).
// Usage: php jsonlint-dump.php corpus.json
require '/oracle-source/vendor/autoload.php';
use Seld\JsonLint\JsonParser;
use Seld\JsonLint\DuplicateKeyException;
$cases = json_decode(file_get_contents($argv[1]), true);
$out = [];
foreach ($cases as $c) {
    $row = [];
    foreach (['default' => 0, 'dup' => JsonParser::DETECT_KEY_CONFLICTS, 'dupassoc' => JsonParser::DETECT_KEY_CONFLICTS | JsonParser::PARSE_TO_ASSOC] as $k => $flags) {
        $p = new JsonParser();
        try { $p->parse($c, $flags); $row[$k] = null; }
        catch (DuplicateKeyException $e) { $row[$k] = ['msg' => $e->getMessage(), 'key' => $e->getKey(), 'line' => $e->getDetails()['line']]; }
        catch (\Seld\JsonLint\ParsingException $e) { $row[$k] = ['msg' => $e->getMessage()]; }
    }
    $out[] = $row;
}
echo json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_INVALID_UTF8_SUBSTITUTE), "\n";
