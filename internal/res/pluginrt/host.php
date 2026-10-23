<?php
// gocomposer plugin host.
//
// gocomposer (Go) implements Composer's dependency management itself and
// starts this PHP process only to run the code of Composer plugins it has
// decided to load (allow-plugins, plugin-api-version and the package checks
// are done in Go, ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Plugin/PluginManager.php). The process hosts the plugin
// objects for the lifetime of the gocomposer command so that state kept by
// plugins between activate() and their event listeners is preserved.
//
// Protocol: one JSON document per line. gocomposer sends {"t":"call"} on fd 3
// and receives {"t":"ret"}/{"t":"exc"} on fd 4. While a call runs, this
// process may call back into gocomposer ({"t":"cb"} answered by
// {"t":"cbret"}/{"t":"cbexc"}); calls nest in both directions.
//
// The Composer API classes visible to plugins live in composer/ (copied or
// ported from the Composer sources of the same commit, see the headers of
// each file); libraries Composer bundles for plugins (symfony/console,
// symfony/process, symfony/finder, composer/semver, ...) are copied unchanged
// into vendor/ with their licenses.

(static function (): void {
    $rt = __DIR__;
    $prefixes = [
        'Composer\\Semver\\' => '/vendor/composer/semver/src/',
        'Composer\\ClassMapGenerator\\' => '/vendor/composer/class-map-generator/src/',
        'Composer\\Pcre\\' => '/vendor/composer/pcre/src/',
        'Composer\\' => '/composer/',
        'Psr\\Log\\' => '/vendor/psr/log/Psr/Log/',
        'React\\Promise\\' => '/vendor/react/promise/src/',
        'Seld\\JsonLint\\' => '/vendor/seld/jsonlint/src/Seld/JsonLint/',
        'Seld\\Signal\\' => '/vendor/seld/signal-handler/src/',
        'Symfony\\Component\\Console\\' => '/vendor/symfony/console/',
        'Symfony\\Component\\Filesystem\\' => '/vendor/symfony/filesystem/',
        'Symfony\\Component\\Finder\\' => '/vendor/symfony/finder/',
        'Symfony\\Component\\Process\\' => '/vendor/symfony/process/',
        'Symfony\\Component\\String\\' => '/vendor/symfony/string/',
        'Symfony\\Contracts\\Service\\' => '/vendor/symfony/service-contracts/',
        'Symfony\\Polyfill\\Ctype\\' => '/vendor/symfony/polyfill-ctype/',
        'Symfony\\Polyfill\\Intl\\Grapheme\\' => '/vendor/symfony/polyfill-intl-grapheme/',
        'Symfony\\Polyfill\\Intl\\Normalizer\\' => '/vendor/symfony/polyfill-intl-normalizer/',
        'Symfony\\Polyfill\\Mbstring\\' => '/vendor/symfony/polyfill-mbstring/',
        'Symfony\\Polyfill\\Php73\\' => '/vendor/symfony/polyfill-php73/',
        'Symfony\\Polyfill\\Php80\\' => '/vendor/symfony/polyfill-php80/',
        'Symfony\\Polyfill\\Php81\\' => '/vendor/symfony/polyfill-php81/',
        'Symfony\\Polyfill\\Php84\\' => '/vendor/symfony/polyfill-php84/',
    ];
    $stubs = [
        '/vendor/symfony/polyfill-intl-normalizer/Resources/stubs/',
        '/vendor/symfony/polyfill-php73/Resources/stubs/',
        '/vendor/symfony/polyfill-php80/Resources/stubs/',
        '/vendor/symfony/polyfill-php81/Resources/stubs/',
        '/vendor/symfony/polyfill-php84/Resources/stubs/',
    ];
    spl_autoload_register(static function (string $class) use ($rt, $prefixes, $stubs): void {
        foreach ($prefixes as $prefix => $dir) {
            if (strpos($class, $prefix) === 0) {
                $file = $rt.$dir.str_replace('\\', '/', substr($class, strlen($prefix))).'.php';
                if (is_file($file)) {
                    require $file;
                }

                return;
            }
        }
        if (strpos($class, '\\') === false) {
            foreach ($stubs as $dir) {
                $file = $rt.$dir.$class.'.php';
                if (is_file($file)) {
                    require $file;

                    return;
                }
            }
        }
    }, true, true);
    foreach ([
        '/vendor/symfony/deprecation-contracts/function.php',
        '/vendor/symfony/polyfill-ctype/bootstrap.php',
        '/vendor/symfony/polyfill-intl-grapheme/bootstrap.php',
        '/vendor/symfony/polyfill-intl-normalizer/bootstrap.php',
        '/vendor/symfony/polyfill-mbstring/bootstrap.php',
        '/vendor/symfony/polyfill-php73/bootstrap.php',
        '/vendor/symfony/polyfill-php80/bootstrap.php',
        '/vendor/symfony/polyfill-php81/bootstrap.php',
        '/vendor/symfony/polyfill-php84/bootstrap.php',
        '/vendor/react/promise/src/functions_include.php',
        '/vendor/symfony/string/Resources/functions.php',
    ] as $file) {
        require_once $rt.$file;
    }
})();

// Composer's ClassLoader / InstalledVersions are the runtime copies
require_once __DIR__.'/composer/Autoload/ClassLoader.php';
require_once __DIR__.'/lib.php';
require_once __DIR__.'/api.php';

set_error_handler([\GocomposerPluginHost\Host::class, 'errorHandler']);
error_reporting(E_ALL);

\GocomposerPluginHost\Host::$in = fopen('php://fd/3', 'r');
\GocomposerPluginHost\Host::$out = fopen('php://fd/4', 'w');
if (!\GocomposerPluginHost\Host::$in || !\GocomposerPluginHost\Host::$out) {
    fwrite(STDERR, "gocomposer plugin host: cannot open the control pipes\n");
    exit(1);
}

while (true) {
    $m = \GocomposerPluginHost\Host::recv();
    if (($m['t'] ?? null) === 'call') {
        if (($m['op'] ?? null) === 'shutdown') {
            \GocomposerPluginHost\Host::send(['t' => 'ret', 'v' => null]);
            exit(0);
        }
        \GocomposerPluginHost\Host::handleCall($m);
    }
}
