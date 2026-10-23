<?php
// Dumps FILTER_VALIDATE_EMAIL, ValidatingArrayLoader::filterUrl and
// ValidatingArrayLoader (CHECK_ALL) results for a corpus using the oracle
// Composer checkout. Usage: php validate-dump.php corpus.json
require '/oracle-source/vendor/autoload.php';
use Composer\Package\Loader\ValidatingArrayLoader;
use Composer\Package\Loader\ArrayLoader;
use Composer\Package\Loader\InvalidPackageException;
$in = json_decode(file_get_contents($argv[1]), true);
$out = ['emails' => [], 'urls' => [], 'packages' => []];
foreach ($in['emails'] as $e) { $out['emails'][] = filter_var($e, FILTER_VALIDATE_EMAIL) !== false; }
$l = new ValidatingArrayLoader(new ArrayLoader(), true, null, ValidatingArrayLoader::CHECK_ALL);
$m = new ReflectionMethod($l, 'filterUrl'); $m->setAccessible(true);
foreach ($in['urls'] as $u) { $out['urls'][] = [$m->invoke($l, $u), $m->invoke($l, $u, ['irc', 'ircs'])]; }
foreach ($in['packages'] as $p) {
    if (!isset($p['version'])) { $p['version'] = '1.0.0'; }
    if (!isset($p['name'])) { $p['name'] = 'dummy/dummy'; }
    $l = new ValidatingArrayLoader(new ArrayLoader(), true, null, ValidatingArrayLoader::CHECK_ALL);
    try { $l->load($p); $out['packages'][] = ['errors' => [], 'warnings' => $l->getWarnings()]; }
    catch (InvalidPackageException $e) { $out['packages'][] = ['errors' => $e->getErrors(), 'warnings' => $e->getWarnings()]; }
    catch (\Throwable $e) { $out['packages'][] = ['exception' => get_class($e).': '.$e->getMessage()]; }
}
echo json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE), "\n";
