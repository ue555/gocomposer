<?php
// Validates corpus strings with the composer/spdx-licenses version locked by
// the oracle Composer checkout. Usage: php spdx-dump.php corpus.json
require '/oracle-source/vendor/autoload.php';
$s = new Composer\Spdx\SpdxLicenses();
$out = [];
foreach (json_decode(file_get_contents($argv[1]), true) as $c) {
    $out[] = ['valid' => $s->validate($c), 'byId' => $s->getLicenseByIdentifier($c)];
}
echo json_encode($out, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE), "\n";
