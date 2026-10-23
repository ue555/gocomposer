<?php
// Evaluates a corpus with the pinned composer/semver used by the oracle Composer
// checkout and prints JSON results. Usage: php semver-dump.php corpus.json
require '/oracle-source/vendor/autoload.php';
use Composer\Semver\VersionParser;
use Composer\Semver\Constraint\Constraint;
use Composer\Semver\Intervals;
use Composer\Semver\Comparator;

$in = json_decode(file_get_contents($argv[1]), true);
$vp = new VersionParser();
$out = ['normalize' => [], 'stability' => [], 'branch' => [], 'compare' => [], 'constraints' => [], 'matches' => [], 'compact' => [], 'subset' => []];
foreach ($in['versions'] as $v) {
    try { $out['normalize'][] = ['ok' => true, 'v' => $vp->normalize($v)]; }
    catch (\Throwable $e) { $out['normalize'][] = ['ok' => false, 'e' => $e->getMessage()]; }
    $out['stability'][] = VersionParser::parseStability($v);
    $out['branch'][] = $vp->normalizeBranch($v);
}
foreach ($in['compare'] as [$a, $b]) { $out['compare'][] = version_compare($a, $b); }
$parsed = [];
foreach ($in['constraints'] as $c) {
    try { $p = $vp->parseConstraints($c); $parsed[] = $p; $out['constraints'][] = ['ok' => true, 's' => (string) $p, 'p' => $p->getPrettyString(), 'lo' => (string) $p->getLowerBound(), 'hi' => (string) $p->getUpperBound()]; }
    catch (\Throwable $e) { $parsed[] = null; $out['constraints'][] = ['ok' => false, 'e' => $e->getMessage()]; }
}
foreach ($parsed as $i => $p) {
    $row = [];
    foreach ($in['matchVersions'] as $v) { $row[] = $p === null ? null : $p->matches(new Constraint('==', $v)); }
    $out['matches'][] = $row;
    $out['compact'][] = $p === null ? null : (string) Intervals::compactConstraint($p);
}
foreach ($in['subset'] as [$a, $b]) {
    $out['subset'][] = Intervals::isSubsetOf($vp->parseConstraints($a), $vp->parseConstraints($b));
}
echo json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES), "\n";
