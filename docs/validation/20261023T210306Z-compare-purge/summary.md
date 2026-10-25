# Composer vs gocomposer comparison

- gocomposer_sha256: `8a8f8d5fb28828fbc3b57432e29d52bc8926fe3ba0db114d8ba3b1345571d910`
- host_os: `darwin/arm64`
- oracle_arch_image: `localhost/gocomposer-oracle:8.4-arch`
- oracle_arch_image_digest: `sha256:93bc04d027eaeb74edeb8ea07a247c5f17cb2e91d4d17fe97681f949f96116e7`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:9cf420cd931841442e55100b424e04a60af91810d670826961c6583bd372308b`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:e3a99ff1ce724bd01dbd722f202bf4afa98eec9bde0789c99ad6f76b266e6f22`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| purge-install | DIFF | 0 | 0 |
| purge-install-all | DIFF | 0 | 1 |
| purge-show | DIFF | 0 | 0 |
| purge-show-name | DIFF | 1 | 0 |
| purge-update | DIFF | 0 | 0 |
| purge-dump-autoload | DIFF | 0 | 1 |
| purge-depends | MATCH | 0 | 0 |
| purge-licenses | DIFF | 0 | 0 |
| purge-status | MATCH | 0 | 0 |
| purge-remove | DIFF | 0 | 0 |

## purge-install

```
stderr:
+ (go) Nothing to install, update or remove
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Downloading acme/cli (1.0.0)
- (oracle)   - Installing acme/cli (1.0.0): Extracting archive
- (oracle)     Skipped installation of bin bin/acme for package acme/cli: name conflicts with an existing file
- (oracle)     Skipped installation of bin bin/acme-sh for package acme/cli: name conflicts with an existing file
+ (go)     Skipped installation of bin bin/acme for package acme/cli: file not found in package
+ (go)     Skipped installation of bin bin/acme-sh for package acme/cli: file not found in package
project tree:
- (oracle) D vendor/acme/cli
- (oracle) D vendor/acme/cli/bin
- (oracle) D vendor/acme/cli/lib
- (oracle) D vendor/acme/cli/lib/AcmeCli
- (oracle) F vendor/acme/cli/bin/acme 755 802bc66258c986f2
- (oracle) F vendor/acme/cli/bin/acme-sh 755 9f2574ae8fe83002
- (oracle) F vendor/acme/cli/composer.json 644 382e59b9dededdc3
- (oracle) F vendor/acme/cli/lib/AcmeCli/Runner.php 644 da07eeda61fd87ba
autoload probe:
- (oracle) class AcmeCli_Runner: autoloaded
+ (go)   exit=0 out= | Warning: include(/work/project/vendor/bin/../acme/cli/bin/acme): Failed to open stream: No such file or directory in /work/project/vendor/bin/acme on line 119 |  | Warning: include(): Failed opening '/work/project/vendor/bin/../acme/cli/bin/acme' for inclusion (include_path='.:/usr/local/lib/php') in /work/project/vendor/bin/acme on line 119
- (oracle)   exit=0 out=hello from acme/util
+ (go)   exit=127 out=vendor/bin/acme-sh: 15: cd: can't cd to ../acme/cli/bin | vendor/bin/acme-sh: 37: exec: /acme-sh: not found
- (oracle)   exit=0 out=acme-sh says hi
```

## purge-install-all

```
exit code: oracle=0 go=1
stderr:
+ (go) Nothing to install, update or remove
- (oracle) Package operations: 5 installs, 0 updates, 0 removals
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle)   - Downloading acme/util (1.1.0)
- (oracle)   - Downloading acme/cli (1.0.0)
- (oracle)   - Downloading acme/log-impl (1.0.0)
- (oracle)   - Downloading acme/http (1.2.0)
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
- (oracle)   - Installing acme/util (1.1.0): Extracting archive
- (oracle)   - Installing acme/cli (1.0.0): Extracting archive
- (oracle)   - Installing acme/log-impl (1.0.0): Extracting archive
- (oracle)   - Installing acme/http (1.2.0): Extracting archive
- (oracle)     Skipped installation of bin bin/acme for package acme/cli: name conflicts with an existing file
- (oracle)     Skipped installation of bin bin/acme-sh for package acme/cli: name conflicts with an existing file
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Could not scan for classes inside "/work/project/vendor/acme/util/lib/" whi  
+ (go)   ch does not appear to be a file nor a folder                                 
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) 1 package you are using is looking for funding.
- (oracle) Use the `composer fund` command to find out more!
project tree:
- (oracle) D vendor/acme
- (oracle) D vendor/acme/cli
- (oracle) D vendor/acme/cli/bin
- (oracle) D vendor/acme/cli/lib
- (oracle) D vendor/acme/cli/lib/AcmeCli
- (oracle) D vendor/acme/http
- (oracle) D vendor/acme/http/src
- (oracle) D vendor/acme/log
- (oracle) D vendor/acme/log-impl
- (oracle) D vendor/acme/log-impl/src
- (oracle) D vendor/acme/log/src
- (oracle) D vendor/acme/util
- (oracle) D vendor/acme/util/lib
- (oracle) D vendor/acme/util/src
- (oracle) F vendor/acme/cli/bin/acme 755 802bc66258c986f2
- (oracle) F vendor/acme/cli/bin/acme-sh 755 9f2574ae8fe83002
- (oracle) F vendor/acme/cli/composer.json 644 382e59b9dededdc3
- (oracle) F vendor/acme/cli/lib/AcmeCli/Runner.php 644 da07eeda61fd87ba
- (oracle) F vendor/acme/http/composer.json 644 dff053c25f488c21
- (oracle) F vendor/acme/http/src/Client.php 644 138fee57b895c83e
- (oracle) F vendor/acme/log-impl/composer.json 644 bd9c251ce1f7b2fc
- (oracle) F vendor/acme/log-impl/src/Impl.php 644 fd59a8b43f5299f1
- (oracle) F vendor/acme/log/composer.json 644 b83d42805152f028
- (oracle) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
- (oracle) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
- (oracle) F vendor/acme/util/composer.json 644 ed976ef852108f1f
- (oracle) F vendor/acme/util/lib/legacy.php 644 aa94d8c20c2ac438
- (oracle) F vendor/acme/util/src/Str.php 644 bd3cbab92c3b093f
- (oracle) F vendor/acme/util/src/functions.php 644 f00c9f042b8e0674
autoload probe:
+ (go) 
+ (go) Warning: require(/work/project/vendor/composer/../acme/util/src/functions.php): Failed to open stream: No such file or directory in /work/project/vendor/composer/autoload_real.php on line 41
+ (go) 
+ (go) Fatal error: Uncaught Error: Failed opening required '/work/project/vendor/composer/../acme/util/src/functions.php' (include_path='.:/usr/local/lib/php') in /work/project/vendor/composer/autoload_real.php:41
+ (go) Stack trace:
+ (go) #0 /work/project/vendor/composer/autoload_real.php(45): {closure:ComposerAutoloaderInit43054bc41b61d0247ad104e2c7c26012::getLoader():37}('571dbb8238a79d3...', '/work/project/v...')
+ (go) #1 /work/project/vendor/autoload.php(22): ComposerAutoloaderInit43054bc41b61d0247ad104e2c7c26012::getLoader()
+ (go) #2 /gobin/probe.php(3): require('/work/project/v...')
+ (go) #3 {main}
+ (go)   thrown in /work/project/vendor/composer/autoload_real.php on line 41
- (oracle) class AcmeAbstractThing: autoloaded
- (oracle) class AcmeCli_Runner: autoloaded
- (oracle) class AcmeSuit: autoloaded
- (oracle) class AcmeThingInterface: autoloaded
- (oracle) class AcmeThingTrait: autoloaded
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
```

## purge-show

```
stdout:
+ (go) acme/cli      1.0.0 CLI tool
```

## purge-show-name

```
exit code: oracle=1 go=0
stdout:
+ (go) name     : acme/cli
+ (go) descrip. : CLI tool
+ (go) keywords : 
+ (go) versions : * 1.0.0
+ (go) released : 2024-01-02, 2 years ago
+ (go) type     : library
+ (go) homepage : 
+ (go) source   : []  
+ (go) dist     : [zip] http://127.0.0.1:8899/dists/acme-cli-1.0.0.zip 94cee099e231284142446146f8f9844a40206270
+ (go) path     : 
+ (go) names    : acme/cli
+ (go) autoload
+ (go) psr-0
+ (go) AcmeCli_ => lib/
+ (go) 
+ (go) requires
+ (go) acme/util ^1.0
+ (go) 
stderr:
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Package "acme/cli" not found, try using --available (-a) to show all availa  
- (oracle)   ble packages.                                                                
- (oracle)                                                                                
- (oracle) 
- (oracle) show [--all] [--locked] [-i|--installed] [-p|--platform] [-a|--available] [-s|--self] [-N|--name-only] [-P|--path] [-t|--tree] [-l|--latest] [-o|--outdated] [--ignore IGNORE] [-M|--major-only] [-m|--minor-only] [--patch-only] [-A|--sort-by-age] [--cooldown-period COOLDOWN-PERIOD] [-D|--direct] [--strict] [-f|--format FORMAT] [--no-dev] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<package> [<version>]]
- (oracle) 
- (oracle) 
```

## purge-update

```
stderr:
+ (go) Nothing to install, update or remove
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Downloading acme/cli (1.0.0)
- (oracle)   - Installing acme/cli (1.0.0): Extracting archive
- (oracle)     Skipped installation of bin bin/acme for package acme/cli: name conflicts with an existing file
- (oracle)     Skipped installation of bin bin/acme-sh for package acme/cli: name conflicts with an existing file
+ (go)     Skipped installation of bin bin/acme for package acme/cli: file not found in package
+ (go)     Skipped installation of bin bin/acme-sh for package acme/cli: file not found in package
project tree:
- (oracle) D vendor/acme/cli
- (oracle) D vendor/acme/cli/bin
- (oracle) D vendor/acme/cli/lib
- (oracle) D vendor/acme/cli/lib/AcmeCli
- (oracle) F vendor/acme/cli/bin/acme 755 802bc66258c986f2
- (oracle) F vendor/acme/cli/bin/acme-sh 755 9f2574ae8fe83002
- (oracle) F vendor/acme/cli/composer.json 644 382e59b9dededdc3
- (oracle) F vendor/acme/cli/lib/AcmeCli/Runner.php 644 da07eeda61fd87ba
autoload probe:
- (oracle) class AcmeCli_Runner: autoloaded
+ (go)   exit=0 out= | Warning: include(/work/project/vendor/bin/../acme/cli/bin/acme): Failed to open stream: No such file or directory in /work/project/vendor/bin/acme on line 119 |  | Warning: include(): Failed opening '/work/project/vendor/bin/../acme/cli/bin/acme' for inclusion (include_path='.:/usr/local/lib/php') in /work/project/vendor/bin/acme on line 119
- (oracle)   exit=0 out=hello from acme/util
+ (go)   exit=127 out=vendor/bin/acme-sh: 15: cd: can't cd to ../acme/cli/bin | vendor/bin/acme-sh: 37: exec: /acme-sh: not found
- (oracle)   exit=0 out=acme-sh says hi
```

## purge-dump-autoload

```
exit code: oracle=0 go=1
stdout:
+ (go) Not all dependencies are installed. Make sure to run a "composer install" to install missing dependencies
project tree:
+ (go) F vendor/composer/autoload_namespaces.php 644 7f8af46ba6702f8f
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_static.php 644 6d9175620c62440e
- (oracle) F vendor/composer/autoload_static.php 644 0bcf0f386273271e
vendor/composer/autoload_static.php:
+ (go)     public static $prefixesPsr0 = array (
+ (go)         'A' =>
+ (go)         array (
+ (go)             'AcmeCli_' =>
+ (go)             array (
+ (go)                 0 => __DIR__ . '/..' . '/acme/cli/lib',
+ (go)             ),
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)             $loader->prefixesPsr0 = ComposerStaticInit43054bc41b61d0247ad104e2c7c26012::$prefixesPsr0;
```

## purge-licenses

```
stdout:
+ (go) acme/cli      1.0.0   none     
```

## purge-remove

```
stderr:
+ (go) Package operations: 0 installs, 0 updates, 1 removal
+ (go)   - Removing acme/cli (1.0.0)
- (oracle) Nothing to install, update or remove
project tree:
- (oracle) D vendor/bin
- (oracle) F vendor/bin/acme 755 0ec11a39f30b01c4
- (oracle) F vendor/bin/acme-sh 755 9a2820eb637c9499
autoload probe:
- (oracle) bin acme mode=0755
- (oracle)   exit=0 out= | Warning: include(/work/project/vendor/bin/../acme/cli/bin/acme): Failed to open stream: No such file or directory in /work/project/vendor/bin/acme on line 119 |  | Warning: include(): Failed opening '/work/project/vendor/bin/../acme/cli/bin/acme' for inclusion (include_path='.:/usr/local/lib/php') in /work/project/vendor/bin/acme on line 119
- (oracle) bin acme-sh mode=0755
- (oracle)   exit=127 out=vendor/bin/acme-sh: 15: cd: can't cd to ../acme/cli/bin | vendor/bin/acme-sh: 37: exec: /acme-sh: not found
```
