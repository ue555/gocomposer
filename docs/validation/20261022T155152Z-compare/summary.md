# Composer vs gocomposer comparison

- gocomposer_sha256: `59417ced76b9e99d09167ce79986aa56dfe336b508338eb0b7f0ce038fe159dc`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| update-basic | MATCH | 0 | 0 |
| update-basic-progress | DIFF | 0 | 0 |
| update-v2-metadata | MATCH | 0 | 0 |
| install-from-lock | DIFF | 0 | 0 |
| install-no-dev | DIFF | 0 | 0 |
| install-twice | DIFF | 0 | 0 |
| dry-run-update | DIFF | 0 | 0 |
| dump-autoload-optimized | MATCH | 0 | 0 |
| dump-autoload-authoritative | MATCH | 0 | 0 |
| replace | MATCH | 0 | 0 |
| conflict-unresolvable | MATCH | 2 | 2 |
| conflict-resolves-lower | MATCH | 0 | 0 |
| platform-php-too-low | MATCH | 2 | 2 |
| missing-extension | MATCH | 2 | 2 |
| ignore-platform-reqs | MATCH | 0 | 0 |
| config-platform-php | MATCH | 0 | 0 |
| not-found-package | MATCH | 2 | 2 |
| stability-beta-flag | MATCH | 0 | 0 |
| stability-too-low | MATCH | 2 | 2 |
| dev-branch-alias | MATCH | 0 | 0 |
| root-alias | MATCH | 0 | 0 |
| partial-update | MATCH | 0 | 0 |
| partial-update-with-deps | MATCH | 0 | 0 |
| broken-archive | DIFF | 100 | 1 |
| checksum-mismatch | DIFF | 1 | 1 |
| download-404 | DIFF | 100 | 1 |
| path-traversal-archive | DIFF | 0 | 1 |
| failure-keeps-vendor | MATCH | 2 | 2 |
| no-lock-install | DIFF | 0 | 0 |

## update-basic-progress

```
stderr:
- (oracle)  0/6 [>---------------------------]   0%
- (oracle)  6/6 [============================] 100%
- (oracle)  0/6 [>---------------------------]   0%
- (oracle)  6/6 [============================] 100%
```

## install-from-lock

```
stdout:
+ (go) No security vulnerability advisories found.
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## install-no-dev

```
stdout:
+ (go) No security vulnerability advisories found.
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## install-twice

```
stdout:
+ (go) No security vulnerability advisories found.
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## dry-run-update

```
stderr:
+ (go) Analyzed 65 packages to resolve dependencies
+ (go) Analyzed 75 rules to resolve dependencies
+ (go) Installs: acme/log:1.1.0, acme/util:1.1.0, acme/cli:1.0.0, acme/devtools:1.0.0, acme/log-impl:1.0.0, acme/http:1.2.0
+ (go) Installs: acme/log:1.1.0, acme/util:1.1.0, acme/cli:1.0.0, acme/devtools:1.0.0, acme/log-impl:1.0.0, acme/http:1.2.0
```

## broken-archive

```
exit code: oracle=100 go=1
stderr:
- (oracle)   - Downloading acme/broken (1.0.0)
- (oracle)   - Downloading acme/broken (1.0.0)
- (oracle)   - Downloading acme/broken (1.0.0)
+ (go)   The downloaded archive for acme/broken is truncated or corrupt, zip end of c
+ (go)   entral directory not found (file size: 201 bytes)
- (oracle)   The downloaded archive for acme/broken is truncated or corrupt, zip end of
- (oracle)   central directory not found (file size: 201 bytes)
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
```

## checksum-mismatch

```
stderr:
+ (go)   The checksum verification of the file failed (downloaded from http://127.0.0
+ (go)   .1:8899/dists/acme-badsum-1.0.0.zip)
- (oracle)   The checksum verification of the file failed (downloaded from http://127.0.
- (oracle)   0.1:8899/dists/acme-badsum-1.0.0.zip)
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
```

## download-404

```
exit code: oracle=100 go=1
stderr:
+ (go)   The "http://127.0.0.1:8899/dists/does-not-exist.zip" file could not be downl
+ (go)   oaded (HTTP/1.1 404 Not Found)
- (oracle)   The "http://127.0.0.1:8899/dists/does-not-exist.zip" file could not be down
- (oracle)   loaded (HTTP/1.1 404 Not Found)
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
```

## path-traversal-archive

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go)     Install of acme/evil failed
+ (go) 
+ (go) In <location>:
+ (go) 
+ (go)   Failed to extract acme/evil: archive entry "../../evil-escaped.txt" escapes
+ (go)   the extraction directory
+ (go) 
+ (go) 
- (oracle) Generating autoload files
project tree:
- (oracle) D vendor
- (oracle) D vendor/acme
- (oracle) D vendor/acme/evil
- (oracle) D vendor/acme/evil/acme-evil-efa6f12
- (oracle) D vendor/composer
- (oracle) F vendor/acme/evil/acme-evil-efa6f12/composer.json 644 ee71c15e5d5d2a26
- (oracle) F vendor/acme/evil/evil-escaped.txt 644 6f20d6e2bf761eb1
- (oracle) F vendor/autoload.php 644 2990f90bc2ff17f2
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 c35c405f584d1f28
- (oracle) F vendor/composer/autoload_real.php 644 fe0d4d531038d171
- (oracle) F vendor/composer/autoload_static.php 644 06a52ea3a393820d
- (oracle) F vendor/composer/installed.json 644 2ae0185f79798182
- (oracle) F vendor/composer/installed.php 644 a8fd330f4f65cf07
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
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
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInita347aad43421fccd839f53909d7cd6e5
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInita347aad43421fccd839f53909d7cd6e5', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInita347aad43421fccd839f53909d7cd6e5', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInita347aad43421fccd839f53909d7cd6e5::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
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
- (oracle) class ComposerStaticInita347aad43421fccd839f53909d7cd6e5
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInita347aad43421fccd839f53909d7cd6e5::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInita347aad43421fccd839f53909d7cd6e5::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInita347aad43421fccd839f53909d7cd6e5::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "acme/evil",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-evil-1.0.0.zip",
- (oracle)                 "reference": "efa6f128215acc47aa16ce8c95bba29d063b6d6e",
- (oracle)                 "shasum": "0c1f3358693eca9995de74ea35706c0641c0457c"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "description": "Archive with path traversal",
- (oracle)             "install-path": "../acme/evil"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'acme/evil' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'efa6f128215acc47aa16ce8c95bba29d063b6d6e',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/evil',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## no-lock-install

```
stdout:
+ (go) No security vulnerability advisories found.
```
