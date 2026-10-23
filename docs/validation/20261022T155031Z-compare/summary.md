# Composer vs gocomposer comparison

- gocomposer_sha256: `b158d4d09c910a5a0b211ea45975aed5d7871b99bafeb49598978c118096c6c7`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| update-basic | DIFF | 0 | 2 |

## update-basic

```
exit code: oracle=0 go=2
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) panic: runtime error: invalid memory address or nil pointer dereference
+ (go) [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x150150]
+ (go) 
+ (go) goroutine 1 [running]:
+ (go) gocomposer/internal/phpjson.Each({0x6ec4d8?, 0x0}, 0x82553d91f70)
+ (go) 	gocomposer/internal/phpjson/phpjson.go:779 +0x50
+ (go) gocomposer/internal/autoload.arrayMergeRecursive.func1(...)
+ (go) 	gocomposer/internal/autoload/generator.go:536
+ (go) gocomposer/internal/autoload.arrayMergeRecursive({0x6ec4d8, 0x82553ea9950}, {0x6ec4d8, 0x0})
+ (go) 	gocomposer/internal/autoload/generator.go:550 +0xc8
+ (go) gocomposer/internal/autoload.(*gen).parseType(0x82553d93120, {0x82553cfbc00?, 0x82553d920a8?, 0x75368?}, {0x3b09d2, 0x5}, 0x82553ed4cd0)
+ (go) 	gocomposer/internal/autoload/generator.go:598 +0x294
+ (go) gocomposer/internal/autoload.(*gen).parseAutoloads(0x82553d93120, {0x82553fc2480?, 0x3b13fc?, 0x7?}, {0x734098?, 0x3eedc0?})
+ (go) 	gocomposer/internal/autoload/generator.go:419 +0x354
+ (go) gocomposer/internal/autoload.Dump({{0x82553ca24e0, 0x14}, {0x82553ca207c, 0xd}, 0x82553fc91b1, 0x0, 0x0, {0x0, 0x0}, 0x0, ...}, ...)
+ (go) 	gocomposer/internal/autoload/generator.go:115 +0x6a8
+ (go) gocomposer/internal/composer.RunAutoloadDump(0x82553d4af00, 0x82553e96f80, 0x82553ed42d0, {0x0, 0x0, 0x0, 0x0, {0x0, 0x0}, 0x1, ...})
+ (go) 	gocomposer/internal/composer/installer.go:367 +0x580
+ (go) gocomposer/internal/composer.(*Installer).dumpAutoload(0x82553e0ec00, 0x82553ed42d0)
+ (go) 	gocomposer/internal/composer/installer.go:279 +0xd4
+ (go) gocomposer/internal/composer.(*Installer).Run(0x82553e0ec00)
+ (go) 	gocomposer/internal/composer/installer.go:179 +0x4f4
+ (go) gocomposer/internal/cli.updateCommand.func1(0x82553eac680)
+ (go) 	gocomposer/internal/cli/commands.go:391 +0x11a4
+ (go) gocomposer/internal/cli.(*App).run(0x82553eac540, 0x82553ed4190, {0x82553d20850?, 0x82553ea5680?, 0x2})
+ (go) 	gocomposer/internal/cli/app.go:271 +0x670
+ (go) gocomposer/internal/cli.(*App).Run(0x82553eac540, {0x82553d20850, 0x2, 0x2})
+ (go) 	gocomposer/internal/cli/app.go:143 +0xfc
+ (go) main.main()
+ (go) 	gocomposer/cmd/gocomposer/main.go:12 +0x54
- (oracle) 1 package you are using is looking for funding.
- (oracle) Use the `composer fund` command to find out more!
project tree:
- (oracle) F vendor/autoload.php 644 1353eb0f89b5ba0c
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 d795431bedb6f53e
- (oracle) F vendor/composer/autoload_files.php 644 0e951f0edb6be02b
- (oracle) F vendor/composer/autoload_namespaces.php 644 7f8af46ba6702f8f
- (oracle) F vendor/composer/autoload_psr4.php 644 23ac8d6d915546dc
- (oracle) F vendor/composer/autoload_real.php 644 aaaddfe74f1f126e
- (oracle) F vendor/composer/autoload_static.php 644 c870c5a9ec15a5b5
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
autoload probe:
- (oracle) class AcmeAbstractThing: autoloaded
- (oracle) class AcmeCli_Runner: autoloaded
- (oracle) class AcmeSuit: autoloaded
- (oracle) class AcmeThingInterface: autoloaded
- (oracle) class AcmeThingTrait: autoloaded
- (oracle) class Acme\DevTools\Dumper: autoloaded
- (oracle) class Acme\Http\Client: autoloaded
- (oracle) class Acme\LogImpl\Impl: autoloaded
- (oracle) class Acme\Log\Logger: autoloaded
- (oracle) class Acme\Log\LoggerInterface: autoloaded
- (oracle) class Acme\Util\Str: autoloaded
- (oracle) class Acme_Legacy_Helper: autoloaded
- (oracle) class App\Kernel: autoloaded
- (oracle) function acme_util_hello: hello from acme/util
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: acme/cli '1.0.0' '94cee099e231284142446146f8f9844a40206270' dev=true
- (oracle) installed: acme/devtools '1.0.0' '4107978d337771591283b614c4075e14c77e9a4c' dev=false
- (oracle) installed: acme/http '1.2.0' 'ce7601587b04bbcb0a274777c6c4a15fe82a9467' dev=true
- (oracle) installed: acme/log '1.1.0' '58b0f69fc5f96981e5f7fae32d592232a36012c2' dev=true
- (oracle) installed: acme/log-impl '1.0.0' 'be55dcec40f479aa52942c4aaa3f393405333530' dev=true
- (oracle) installed: acme/util '1.1.0' '010715fb4511de3a64a1da0e47aa42e3b0104c1c' dev=true
- (oracle) installed: psr/log-implementation NULL NULL dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) bin acme mode=0755
- (oracle)   exit=0 out=hello from acme/util
- (oracle) bin acme-sh mode=0755
- (oracle)   exit=0 out=acme-sh says hi
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'AcmeAbstractThing' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeSuit' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeThingInterface' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeThingTrait' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'Acme_Legacy_Helper' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_files.php:
- (oracle) <?php
- (oracle) // autoload_files.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     '571dbb8238a79d30de0279182c1a3513' => $vendorDir . '/acme/util/src/functions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle)     'Acme\\Util\\' => array($vendorDir . '/acme/util/src'),
- (oracle)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
- (oracle)     'Acme\\LogImpl\\' => array($vendorDir . '/acme/log-impl/src'),
- (oracle)     'Acme\\Http\\' => array($vendorDir . '/acme/http/src'),
- (oracle)     'Acme\\DevTools\\' => array($vendorDir . '/acme/devtools/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         require __DIR__ . '/platform_check.php';
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         $filesToLoad = \Composer\Autoload\ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$files;
- (oracle)         $requireFile = \Closure::bind(static function ($fileIdentifier, $file) {
- (oracle)             if (empty($GLOBALS['__composer_autoload_files'][$fileIdentifier])) {
- (oracle)                 $GLOBALS['__composer_autoload_files'][$fileIdentifier] = true;
- (oracle) 
- (oracle)                 require $file;
- (oracle)             }
- (oracle)         }, null, null);
- (oracle)         foreach ($filesToLoad as $fileIdentifier => $file) {
- (oracle)             $requireFile($fileIdentifier, $file);
- (oracle)         }
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca
- (oracle) {
- (oracle)     public static $files = array (
- (oracle)         '571dbb8238a79d30de0279182c1a3513' => __DIR__ . '/..' . '/acme/util/src/functions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)             'Acme\\Util\\' => 10,
- (oracle)             'Acme\\Log\\' => 9,
- (oracle)             'Acme\\LogImpl\\' => 13,
- (oracle)             'Acme\\Http\\' => 10,
- (oracle)             'Acme\\DevTools\\' => 14,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)         'Acme\\Util\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/util/src',
- (oracle)         ),
- (oracle)         'Acme\\Log\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log/src',
- (oracle)         ),
- (oracle)         'Acme\\LogImpl\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log-impl/src',
- (oracle)         ),
- (oracle)         'Acme\\Http\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/http/src',
- (oracle)         ),
- (oracle)         'Acme\\DevTools\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/devtools/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixesPsr0 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'AcmeCli_' =>
- (oracle)             array (
- (oracle)                 0 => __DIR__ . '/..' . '/acme/cli/lib',
- (oracle)             ),
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'AcmeAbstractThing' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeSuit' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeThingInterface' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeThingTrait' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'Acme_Legacy_Helper' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixDirsPsr4;
- (oracle)             $loader->prefixesPsr0 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixesPsr0;
- (oracle)             $loader->classMap = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/platform_check.php:
- (oracle) <?php
- (oracle) // platform_check.php @generated by Composer
- (oracle) 
- (oracle) $issues = array();
- (oracle) 
- (oracle) if (!(PHP_VERSION_ID >= 70400)) {
- (oracle)     $issues[] = 'Your Composer dependencies require a PHP version ">= 7.4.0". You are running ' . PHP_VERSION . '.';
- (oracle) }
- (oracle) 
- (oracle) if ($issues) {
- (oracle)     if (!headers_sent()) {
- (oracle)         header('HTTP/1.1 500 Internal Server Error');
- (oracle)     }
- (oracle)     if (!ini_get('display_errors')) {
- (oracle)         if (PHP_SAPI === 'cli' || PHP_SAPI === 'phpdbg') {
- (oracle)             fwrite(STDERR, 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . implode(PHP_EOL, $issues) . PHP_EOL.PHP_EOL);
- (oracle)         } elseif (!headers_sent()) {
- (oracle)             echo 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . str_replace('You are running '.PHP_VERSION.'.', '', implode(PHP_EOL, $issues)) . PHP_EOL.PHP_EOL;
- (oracle)         }
- (oracle)     }
- (oracle)     throw new \RuntimeException(
- (oracle)         'Composer detected issues in your platform: ' . implode(' ', $issues)
- (oracle)     );
- (oracle) }
- (oracle) 
```
