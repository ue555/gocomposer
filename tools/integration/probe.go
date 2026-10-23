package main

// probePHP loads vendor/autoload.php with the official PHP and reports, for
// every class-like declaration found in the project and vendor sources,
// which file the generated ClassLoader resolves it to (findFile, so nothing
// is executed), plus the InstalledVersions data.
const probePHP = `<?php
error_reporting(E_ALL);
$loader = require 'vendor/autoload.php';
$cwd = getcwd() . '/';
$names = array();
foreach (array('vendor', 'src', 'app', 'lib') as $dir) {
    if (!is_dir($dir)) continue;
    $files = array();
    $it = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($dir, FilesystemIterator::SKIP_DOTS));
    foreach ($it as $f) {
        $p = $f->getPathname();
        if (substr($p, -4) !== '.php' || strpos($p, 'vendor/composer/') === 0) continue;
        $files[] = $p;
    }
    sort($files);
    foreach ($files as $p) {
        $tokens = @token_get_all(file_get_contents($p));
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
                if ($i > 0 && is_array($tokens[$i-1]) && in_array($tokens[$i-1][0], array(T_DOUBLE_COLON, T_NEW), true)) continue;
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
$resolved = 0;
foreach ($names as $c) {
    $f = $loader->findFile($c);
    if ($f === false) { echo 'unresolved: ', $c, "\n"; continue; }
    $resolved++;
    $real = realpath($f);
    echo 'class ', $c, ' => ', $real !== false && strpos($real, $cwd) === 0 ? substr($real, strlen($cwd)) : $f, "\n";
}
echo 'resolved: ', $resolved, ' of ', count($names), "\n";
if (class_exists('Composer\InstalledVersions')) {
    $root = Composer\InstalledVersions::getRootPackage();
    echo 'root: ', $root['name'], ' ', $root['pretty_version'], ' dev=', var_export($root['dev'], true), "\n";
    $pkgs = Composer\InstalledVersions::getInstalledPackages();
    sort($pkgs);
    foreach ($pkgs as $p) {
        echo 'installed: ', $p, ' ', var_export(Composer\InstalledVersions::getPrettyVersion($p), true), ' ', var_export(Composer\InstalledVersions::getReference($p), true), "\n";
    }
}
`
