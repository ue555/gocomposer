# Composer vs gocomposer comparison

- gocomposer_sha256: `b4d097afa1494f053879aa919f01faa68de52a74e496d24e5ebbed790115a086`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.3`
- oracle_image_digest: `sha256:91789a6658a5075533a2b26cfbdb2a264ab6051e638872129021bccc27665fce`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.3-zip`
- oracle_zip_image_digest: `sha256:08e40da7145cf527bf595df929ce3394a5ae7bcfd43a8e4bac8ec1243e406bf9`
- php_version: `8.3.36`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| providers-update | MATCH | 0 | 0 |
| providers-update-verbose | MATCH | 0 | 0 |
| providers-install-lock | MATCH | 0 | 0 |
| providers-update-cached | MATCH | 0 | 0 |
| providers-root-listing | MATCH | 0 | 0 |
| providers-dev-alias | MATCH | 0 | 0 |
| providers-virtual | MATCH | 2 | 2 |
| providers-not-found | MATCH | 2 | 2 |
| providers-bad-hash | MATCH | 1 | 1 |
| providers-bad-hash-other | MATCH | 0 | 0 |
| providers-require | MATCH | 0 | 0 |
| providers-require-missing | MATCH | 1 | 1 |
| providers-show-available | MATCH | 0 | 0 |
| providers-show-all | MATCH | 0 | 0 |
| providers-show-all-names | MATCH | 0 | 0 |
| providers-outdated | MATCH | 0 | 0 |
| providers-update-one | MATCH | 0 | 0 |
| providers-depends | MATCH | 0 | 0 |
| providers-search | MATCH | 0 | 0 |
| providers-search-name | MATCH | 0 | 0 |
| providers-search-vendor | MATCH | 0 | 0 |
| providers-search-list | MATCH | 0 | 0 |
| providers-search-list-other | MATCH | 0 | 0 |
| providers-lazy-update | MATCH | 0 | 0 |
| providers-lazy-not-found | MATCH | 2 | 2 |
| providers-lazy-show | MATCH | 0 | 0 |
| providers-lazy-search | MATCH | 0 | 0 |
| providers-lazy-cached | MATCH | 0 | 0 |
| archive3-xz | MATCH | 0 | 0 |
| archive3-xz-status | MATCH | 1 | 1 |
| archive3-xz-reinstall | MATCH | 0 | 0 |
| archive3-xz-remove | MATCH | 0 | 0 |
| archive3-rar-status | MATCH | 1 | 1 |
| archive3-xz-broken | MATCH | 1 | 1 |
| archive3-rar-no-unrar | MATCH | 1 | 1 |
| archive3-rar | MATCH | 0 | 0 |
| archive3-rar-broken | MATCH | 1 | 1 |
| archive3-zip-7z | MATCH | 0 | 0 |
| archive3-zip-7z-zipext | MATCH | 0 | 0 |
| archive3-zip-7z-crc | MATCH | 1 | 1 |
| archive3-zip-7z-crc-zipext | DIFF | 0 | 1 |
| archive3-zip-no-command | MATCH | 1 | 1 |
| archive3-zip-no-command-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-crc-zipext | DIFF | 0 | 1 |
| warnings-update | MATCH | 0 | 0 |
| warnings-update-ansi | MATCH | 0 | 0 |
| warnings-show | MATCH | 0 | 0 |
| warnings-search-api | MATCH | 0 | 0 |
| warnings-search-list | MATCH | 0 | 0 |
| warnings-show-names-list | MATCH | 0 | 0 |
| warnings-audit-api | MATCH | 1 | 1 |
| warnings-update-audit-api | MATCH | 0 | 0 |
| warnings-filter-audit | MATCH | 0 | 0 |
| warnings-providers | MATCH | 0 | 0 |
| warnings-two-repos | MATCH | 0 | 0 |
| warnings-http-403 | MATCH | 100 | 100 |
| warnings-http-404-metadata | MATCH | 2 | 2 |

## archive3-zip-7z-crc-zipext

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go)     Install of archfix/zipcrc failed
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   The archive for "archfix/zipcrc" may contain identical file names with diff  
+ (go)   erent capitalization (which fails on case insensitive filesystems): ZipArch  
+ (go)   ive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json": "CRC  
+ (go)    error"                                                                      
+ (go)                                                                                
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   ZipArchive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json  
+ (go)   ": "CRC error"                                                               
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Generating autoload files
project tree:
- (oracle) D vendor
- (oracle) D vendor/archfix
- (oracle) D vendor/archfix/zipcrc
- (oracle) D vendor/archfix/zipcrc/bin
- (oracle) D vendor/archfix/zipcrc/src
- (oracle) D vendor/bin
- (oracle) D vendor/composer
- (oracle) F vendor/archfix/zipcrc/bin/archtool 755 63e5cd20333a2449
- (oracle) F vendor/archfix/zipcrc/composer.json 644 eabbb5b5ab82f5c6
- (oracle) F vendor/archfix/zipcrc/src/Zipcrc.php 644 717fa4118cb95c22
- (oracle) F vendor/autoload.php 644 e4df171b3ebd67ad
- (oracle) F vendor/bin/archtool 755 4973c40d2f80670e
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 d5274594f6204537
- (oracle) F vendor/composer/autoload_real.php 644 a18af7b42f2a3a7a
- (oracle) F vendor/composer/autoload_static.php 644 584da5c998a0b31c
- (oracle) F vendor/composer/installed.json 644 f9604e5fba77b64b
- (oracle) F vendor/composer/installed.php 644 a980ee71a4a0a74b
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
- (oracle)     'Archfix\\' => array($vendorDir . '/archfix/zipcrc/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::getInitializer($loader));
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
- (oracle) class ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'Archfix\\' => 8,
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'Archfix\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/archfix/zipcrc/src',
- (oracle)         ),
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
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "archfix/zipcrc",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/git/repo3/arch/dists/archfix-zipcrc-1.0.0.zip",
- (oracle)                 "reference": "718127c0f5bd7429bbd699b5d2aece2f60ed5687",
- (oracle)                 "shasum": "d60bf86d31591f2bd15e27698649bd400b8cfa1b"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "bin": [
- (oracle)                 "bin/archtool"
- (oracle)             ],
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Archfix\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "install-path": "../archfix/zipcrc"
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
- (oracle)         'archfix/zipcrc' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '718127c0f5bd7429bbd699b5d2aece2f60ed5687',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../archfix/zipcrc',
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

## archive3-zip-no-command-crc-zipext

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go)     Install of archfix/zipcrc failed
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   The archive for "archfix/zipcrc" may contain identical file names with diff  
+ (go)   erent capitalization (which fails on case insensitive filesystems): ZipArch  
+ (go)   ive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json": "CRC  
+ (go)    error"                                                                      
+ (go)                                                                                
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   ZipArchive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json  
+ (go)   ": "CRC error"                                                               
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Generating autoload files
project tree:
- (oracle) D vendor
- (oracle) D vendor/archfix
- (oracle) D vendor/archfix/zipcrc
- (oracle) D vendor/archfix/zipcrc/bin
- (oracle) D vendor/archfix/zipcrc/src
- (oracle) D vendor/bin
- (oracle) D vendor/composer
- (oracle) F vendor/archfix/zipcrc/bin/archtool 755 63e5cd20333a2449
- (oracle) F vendor/archfix/zipcrc/composer.json 644 eabbb5b5ab82f5c6
- (oracle) F vendor/archfix/zipcrc/src/Zipcrc.php 644 717fa4118cb95c22
- (oracle) F vendor/autoload.php 644 e4df171b3ebd67ad
- (oracle) F vendor/bin/archtool 755 4973c40d2f80670e
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 d5274594f6204537
- (oracle) F vendor/composer/autoload_real.php 644 a18af7b42f2a3a7a
- (oracle) F vendor/composer/autoload_static.php 644 584da5c998a0b31c
- (oracle) F vendor/composer/installed.json 644 f9604e5fba77b64b
- (oracle) F vendor/composer/installed.php 644 a980ee71a4a0a74b
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
- (oracle)     'Archfix\\' => array($vendorDir . '/archfix/zipcrc/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit8fbfb2363b8aeecb805302f10d456dd9', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::getInitializer($loader));
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
- (oracle) class ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'Archfix\\' => 8,
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'Archfix\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/archfix/zipcrc/src',
- (oracle)         ),
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
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit8fbfb2363b8aeecb805302f10d456dd9::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "archfix/zipcrc",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/git/repo3/arch/dists/archfix-zipcrc-1.0.0.zip",
- (oracle)                 "reference": "718127c0f5bd7429bbd699b5d2aece2f60ed5687",
- (oracle)                 "shasum": "d60bf86d31591f2bd15e27698649bd400b8cfa1b"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "bin": [
- (oracle)                 "bin/archtool"
- (oracle)             ],
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Archfix\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "install-path": "../archfix/zipcrc"
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
- (oracle)         'archfix/zipcrc' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '718127c0f5bd7429bbd699b5d2aece2f60ed5687',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../archfix/zipcrc',
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
