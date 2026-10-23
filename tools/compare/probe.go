package main

// probePHP loads vendor/autoload.php with official PHP and reports which
// declared classes/functions resolve, InstalledVersions data and bin output.
const probePHP = `<?php
error_reporting(E_ALL);
$loader = require 'vendor/autoload.php';
$names = array();
$dirs = array('vendor', 'src', 'lib', 'app');
foreach ($dirs as $dir) {
    if (!is_dir($dir)) continue;
    $files = array();
    $it = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($dir, FilesystemIterator::SKIP_DOTS));
    foreach ($it as $f) {
        $p = $f->getPathname();
        if (substr($p, -4) !== '.php') continue;
        if (strpos($p, 'vendor/composer/') === 0 || strpos($p, 'vendor/bin/') === 0) continue;
        $files[] = $p;
    }
    sort($files);
    foreach ($files as $p) {
        $tokens = token_get_all(file_get_contents($p));
        $ns = '';
        for ($i = 0, $n = count($tokens); $i < $n; $i++) {
            if (!is_array($tokens[$i])) continue;
            if ($tokens[$i][0] === T_NAMESPACE) {
                $ns = '';
                for ($j = $i + 1; $j < $n; $j++) {
                    if (is_array($tokens[$j]) && in_array($tokens[$j][0], array(T_STRING, T_NAME_QUALIFIED), true)) { $ns .= $tokens[$j][1]; }
                    elseif ($tokens[$j] === ';' || $tokens[$j] === '{') break;
                }
            }
            if (in_array($tokens[$i][0], array(T_CLASS, T_INTERFACE, T_TRAIT, T_ENUM), true)) {
                if ($i > 0 && is_array($tokens[$i-1]) && $tokens[$i-1][0] === T_DOUBLE_COLON) continue;
                for ($j = $i + 1; $j < $n; $j++) {
                    if (is_array($tokens[$j]) && $tokens[$j][0] === T_STRING) { $names[] = ltrim($ns . '\\' . $tokens[$j][1], '\\'); break; }
                    if (!is_array($tokens[$j]) || $tokens[$j][0] !== T_WHITESPACE) break;
                }
            }
        }
    }
}
$names = array_unique($names);
sort($names);
foreach ($names as $c) {
    $ok = class_exists($c) || interface_exists($c) || trait_exists($c) || (function_exists('enum_exists') && enum_exists($c));
    echo 'class ', $c, ': ', $ok ? 'autoloaded' : 'MISSING', "\n";
}
echo 'function acme_util_hello: ', function_exists('acme_util_hello') ? acme_util_hello() : 'missing', "\n";
if (class_exists('Composer\InstalledVersions')) {
    $root = Composer\InstalledVersions::getRootPackage();
    echo 'root: ', $root['name'], ' ', $root['pretty_version'], ' dev=', var_export($root['dev'], true), "\n";
    $pkgs = Composer\InstalledVersions::getInstalledPackages();
    sort($pkgs);
    foreach ($pkgs as $p) {
        echo 'installed: ', $p, ' ', var_export(Composer\InstalledVersions::getPrettyVersion($p), true), ' ', var_export(Composer\InstalledVersions::getReference($p), true), ' dev=', var_export(Composer\InstalledVersions::isInstalled($p, false), true), "\n";
    }
}
if (is_dir('vendor/bin')) {
    $bins = scandir('vendor/bin');
    foreach ($bins as $b) {
        if ($b[0] === '.') continue;
        echo 'bin ', $b, ' mode=', substr(sprintf('%o', fileperms('vendor/bin/' . $b)), -4), "\n";
        $out = array(); $code = 0;
        exec('vendor/bin/' . escapeshellarg($b) . ' 2>&1', $out, $code);
        echo '  exit=', $code, ' out=', implode(' | ', $out), "\n";
    }
}
`
