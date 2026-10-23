<?php
// Runs a regex corpus through official PHP preg_* functions the way
// Composer\Pcre\Preg (vendor/composer/pcre of Composer
// f83588392df43d4fd027d84d26de459946ef4e80) calls them, and prints JSON results
// for internal/pcre's oracle test. Usage: php pcre-dump.php corpus.json
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
require '/oracle-source/vendor/autoload.php';

use Composer\Pcre\Preg;

set_error_handler(static function (): bool { return true; });

// strings that are not valid UTF-8 are emitted as {"b64": ...}
function str(?string $s)
{
    if ($s === null || preg_match('//u', $s) === 1) {
        return $s;
    }
    return ['b64' => base64_encode($s)];
}

function pairs(?array $m): ?array
{
    if ($m === null) {
        return null;
    }
    $out = [];
    foreach ($m as $k => $v) {
        $out[] = [(string) $k, str($v)];
    }
    return $out;
}

// byte-exact serialisation of a preg_replace_callback match array
function serial(array $m): string
{
    $parts = [];
    foreach ($m as $k => $v) {
        $parts[] = $k . '=' . ($v === null ? 'N' : 'S' . strlen($v) . ':' . $v);
    }
    return '<' . implode('|', $parts) . '>';
}

// Documents for the match/stack limit calibration (mirrored by limitDoc in
// internal/pcre/limits_test.go).
function limitPart(string $c, int $n): string
{
    $a = [];
    switch ($c) {
        case 'r': for ($i = 0; $i < $n; $i++) { $a[] = "        \"vendor$i/package-$i\": \"^".($i % 7).".0\""; } return "\"require\": {\n".implode(",\n", $a)."\n    }";
        case 'a': for ($i = 0; $i < $n; $i++) { $a[] = "            \"Ns$i\\\\\\\\Sub\\\\\\\\\": \"packages/p$i/src\""; } return "\"autoload\": {\"psr-4\": {\n".implode(",\n", $a)."\n    }}";
        case 'e': for ($i = 0; $i < $n; $i++) { $a[] = "{\"i\": $i, \"s\": \"v\\u00e9$i\", \"f\": $i.5e-1}"; } return "\"extra\": {\"list\": [".implode(", ", $a)."]}";
        case 'n': for ($i = 0; $i < $n; $i++) { $a[] = ($i % 3 == 0 ? "-$i.25e+3" : ($i % 3 == 1 ? "$i" : "0.$i")); } return "\"nums\": [".implode(",", $a)."]";
        case 's': for ($i = 0; $i < $n; $i++) { $a[] = "\"x\\n\\\"y\\\\z\\/\\u0041$i\""; } return "\"strs\": [".implode(", ", $a)."]";
        case 'd': for ($i = 0; $i < $n; $i++) { $a[] = "[[[true, null], {}], [false]]"; } return "\"deep\": [".implode(",", $a)."]";
        case 'w': for ($i = 0; $i < $n; $i++) { $a[] = "\n        {\n            \"k\" :  \"v\" ,\n            \"m\" :  [ 1 , 2 ]\n        }"; } return "\"ws\" : [".implode(" ,", $a)."\n    ]";
        case 'k': for ($i = 0; $i < $n; $i++) { $a[] = "\"key$i\": \"v\""; } return implode(",\n    ", $a);
        case 'A': return "\"x\": ".str_repeat('[', $n).str_repeat(']', $n).", \"require\": {}";
        case 'O': return "\"x\": ".str_repeat('{"a":', $n).'1'.str_repeat('}', $n).", \"require\": {}";
    }
    throw new \LogicException($c);
}

function limitDoc(string $parts, int $n): string
{
    $s = "{\n    \"name\": \"big/doc\"";
    foreach (str_split($parts) as $p) {
        $s .= ",\n    ".limitPart($p, $n);
    }
    return $s."\n}";
}

// patterns may contain @@DEFINES@@, replaced by Composer's JsonManipulator::DEFINES
$defines = (new \ReflectionClass(\Composer\Json\JsonManipulator::class))->getConstant('DEFINES');
$in = json_decode(file_get_contents($argv[1]), true, 512, JSON_THROW_ON_ERROR);
array_walk_recursive($in, static function (&$v, $k) use ($defines): void {
    if (($k === 'p' || $k === 'pattern') && is_string($v)) {
        $v = str_replace('@@DEFINES@@', $defines, $v);
    }
});
if (isset($in['limits'])) {
    // binary search the largest n for which preg_match does not fail
    $out = [];
    foreach ($in['limits'] as $c) {
        $re = $c['pattern'];
        $fails = static function (int $n) use ($re, $c): ?int {
            return preg_match($re, limitDoc($c['parts'], $n)) === false ? preg_last_error() : null;
        };
        $lo = 1;
        $hi = $c['max'];
        $err = null;
        while ($lo < $hi) {
            $mid = intdiv($lo + $hi + 1, 2);
            $e = $fails($mid);
            if ($e === null) {
                $lo = $mid;
            } else {
                $hi = $mid - 1;
            }
        }
        $out[] = ['n' => $lo, 'err' => $fails($lo + 1), 'sha1' => sha1(limitDoc($c['parts'], $lo))];
    }
    echo json_encode(['php' => PHP_VERSION, 'pcre' => PCRE_VERSION, 'jit' => (bool) ini_get('pcre.jit'), 'defines' => $defines, 'results' => $out], JSON_UNESCAPED_SLASHES), "\n";
    exit(0);
}
$out = [];
foreach ($in['cases'] as $case) {
    $res = [];
    try {
        switch ($case['op']) {
            case 'match':
                $matches = null;
                $r = Preg::match($case['p'], $case['s'], $matches);
                $res = ['r' => $r, 'm' => pairs($matches)];
                break;
            case 'replace':
                $count = 0;
                $r = Preg::replace($case['p'], $case['r'], $case['s'], $case['limit'] ?? -1, $count);
                $res = ['r' => str($r), 'n' => $count];
                break;
            case 'callback':
                $count = 0;
                $r = Preg::replaceCallback($case['p'], static function (array $m): string {
                    return serial($m);
                }, $case['s'], $case['limit'] ?? -1, $count);
                $res = ['r' => str($r), 'n' => $count];
                break;
            case 'quote':
                $res = ['r' => str(isset($case['d']) ? preg_quote($case['s'], $case['d']) : preg_quote($case['s']))];
                break;
            case 'addcslashes':
                $res = ['r' => str(addcslashes($case['s'], $case['c']))];
                break;
            default:
                throw new \LogicException('unknown op ' . $case['op']);
        }
    } catch (\Composer\Pcre\PcreException $e) {
        $res = ['err' => $e->getCode() === PREG_INTERNAL_ERROR ? 'compile' : 'limit:' . $e->getCode()];
    }
    $out[] = $res;
}
echo json_encode(['php' => PHP_VERSION, 'pcre' => PCRE_VERSION, 'jit' => (bool) ini_get('pcre.jit'), 'defines' => $defines, 'results' => $out], JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR), "\n";
