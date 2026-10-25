# Composer vs gocomposer comparison

- gocomposer_sha256: `3796dd581cbda81f9318b778b10829f14d7da44dbb0972afb8fe48fdfef2ad0d`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:cd7f1f040bba14a1974d175df2110347c2be1d880834ad6a106a3578bfd929fa`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| root-alias | MATCH | 0 | 0 |
| term-href-show | MATCH | 0 | 0 |
| term-href-show-name | MATCH | 0 | 0 |
| term-href-show-latest | MATCH | 0 | 0 |
| term-href-licenses | MATCH | 0 | 0 |
| term-href-fund | MATCH | 0 | 0 |
| term-href-search | MATCH | 0 | 0 |
| term-href-depends | MATCH | 0 | 0 |
| term-href-config-list | MATCH | 0 | 0 |
| term-href-audit | DIFF | 1 | 1 |
| term-href-audit-plain | DIFF | 1 | 1 |
| term-href-jetbrains | MATCH | 0 | 0 |
| term-href-konsole-old | MATCH | 0 | 0 |
| term-href-konsole-new | MATCH | 0 | 0 |
| term-href-konsole-zero | MATCH | 0 | 0 |
| term-href-idea | MATCH | 0 | 0 |
| term-href-no-ansi | MATCH | 0 | 0 |
| term-href-ansi-pipe | MATCH | 0 | 0 |
| term-href-ansi-pipe-fund | MATCH | 0 | 0 |
| term-progress-install | DIFF | 0 | 0 |
| term-progress-update | DIFF | 0 | 0 |
| term-progress-require | MATCH | 0 | 0 |
| term-progress-verbose | DIFF | 0 | 0 |
| term-progress-very-verbose | MATCH | 0 | 0 |
| term-progress-narrow | DIFF | 0 | 0 |
| term-progress-tiny | DIFF | 0 | 0 |
| term-progress-no-progress | MATCH | 0 | 0 |
| term-progress-no-ansi | MATCH | 0 | 0 |
| term-progress-ci | MATCH | 0 | 0 |
| term-progress-ansi-pipe | MATCH | 0 | 0 |
| curl-refused | MATCH | 100 | 100 |
| curl-empty-reply | MATCH | 100 | 100 |
| curl-reset | MATCH | 100 | 100 |
| curl-reset-debug | DIFF | 100 | 100 |
| curl-http09 | MATCH | 100 | 100 |
| curl-partial | MATCH | 100 | 100 |
| curl-tls-eof | MATCH | 100 | 100 |
| curl-tls-wrong-version | MATCH | 100 | 100 |
| curl-tls-reset | MATCH | 100 | 100 |
| curl-tls-self-signed | MATCH | 100 | 100 |
| curl-tls-self-signed-system-ca | MATCH | 100 | 100 |
| curl-tls-wrong-host | MATCH | 100 | 100 |
| curl-tls-expired | MATCH | 100 | 100 |
| curl-tls-unknown-ca | MATCH | 100 | 100 |
| curl-tls-good | MATCH | 2 | 2 |
| curl-tls-show | MATCH | 100 | 100 |
| curl-timeout | MATCH | 100 | 100 |
| curl-tls-timeout | MATCH | 100 | 100 |
| curl-post-empty-reply | MATCH | 100 | 100 |
| curl-post-timeout-ignored | MATCH | 0 | 0 |
| root-validate | MATCH | 0 | 0 |
| root-validate-allowed | MATCH | 0 | 0 |
| root-validate-allowed-zero | MATCH | 0 | 0 |
| root-validate-no-plugins | MATCH | 0 | 0 |
| root-update | MATCH | 0 | 0 |
| root-unknown-command | MATCH | 1 | 1 |
| root-unknown-command-allowed | MATCH | 1 | 1 |
| root-help | MATCH | 1 | 1 |
| root-script-alias | MATCH | 1 | 1 |
| root-tty-validate-yes | MATCH | 0 | 0 |
| root-tty-validate-default | MATCH | 0 | 0 |
| root-tty-validate-no | MATCH | 1 | 1 |
| root-tty-validate-n | MATCH | 0 | 0 |
| root-tty-validate-allowed | MATCH | 0 | 0 |
| root-tty-unknown-command-yes | MATCH | 1 | 1 |
| root-tty-unknown-command-no | MATCH | 1 | 1 |
| root-tty-help-yes | DIFF | 0 | 0 |
| root-tty-script-alias-yes | MATCH | 0 | 0 |
| root-tty-update-yes | MATCH | 0 | 0 |
| root-tty-no-plugins | MATCH | 0 | 0 |

## term-href-audit

```
stdout:
+ (go) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/pksa-test-0001\e\PKSA-test-0001\e]8;;\e\                                                                   |
+ (go) |\e[32m CVE               \e[39m| \e]8;;https://www.cve.org/cverecord?id=cve-2024-0001\e\CVE-2024-0001\e]8;;\e\                                                                    |
- (oracle) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/PKSA-test-0001\e\PKSA-test-0001\e]8;;\e\                                                                   |
- (oracle) |\e[32m CVE               \e[39m| \e]8;;https://www.cve.org/CVERecord?id=CVE-2024-0001\e\CVE-2024-0001\e]8;;\e\                                                                    |
+ (go) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/pksa-test-0002\e\PKSA-test-0002\e]8;;\e\                                                                   |
- (oracle) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/PKSA-test-0002\e\PKSA-test-0002\e]8;;\e\                                                                   |
```

## term-href-audit-plain

```
stdout:
+ (go) Advisory ID: \e]8;;https://packagist.org/security-advisories/pksa-test-0001\e\PKSA-test-0001\e]8;;\e\
+ (go) CVE: \e]8;;https://www.cve.org/cverecord?id=cve-2024-0001\e\CVE-2024-0001\e]8;;\e\
- (oracle) Advisory ID: \e]8;;https://packagist.org/security-advisories/PKSA-test-0001\e\PKSA-test-0001\e]8;;\e\
- (oracle) CVE: \e]8;;https://www.cve.org/CVERecord?id=CVE-2024-0001\e\CVE-2024-0001\e]8;;\e\
+ (go) Advisory ID: \e]8;;https://packagist.org/security-advisories/pksa-test-0002\e\PKSA-test-0002\e]8;;\e\
- (oracle) Advisory ID: \e]8;;https://packagist.org/security-advisories/PKSA-test-0002\e\PKSA-test-0002\e]8;;\e\
```

## term-progress-install

```
stdout:
+ (go)  0/6 [>---------------------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>---------------------------]   0%\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-update

```
stdout:
+ (go)  0/6 [>---------------------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32m1 package suggestions were added by new dependencies, use `composer suggest` to see details.\e[39m
- (oracle)  0/6 [>---------------------------]   0%\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32m1 package suggestions were added by new dependencies, use `composer suggest` to see details.\e[39m
```

## term-progress-verbose

```
stdout:
+ (go)  0/6 [>---------------------------]   0% <elapsed>\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100% <elapsed>\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>---------------------------]   0% <elapsed>\e[1G\e[2K 6/6 [============================] 100% <elapsed>\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-narrow

```
stdout:
+ (go)  0/6 [>-----------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [==================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>-----------------]   0%\e[1G\e[2K 6/6 [==================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-tiny

```
stdout:
+ (go)  0/6 [>]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [=] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>]   0%\e[1G\e[2K 6/6 [=] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## curl-reset-debug

```
stderr:
- (oracle) Executing command (CWD): 'git' '--version'
- (oracle) Executing command (/work/project): 'git' 'branch' '-a' '--no-color' '--no-abbrev' '-v'
- (oracle) Executing command (/work/project): 'git' 'describe' '--exact-match' '--tags'
- (oracle) Executing command (/work/project): 'git' 'rev-list' '--no-commit-header' '--format=%H' '-n1' 'HEAD' '--no-show-signature'
- (oracle) Executing command (/work/project): 'hg' 'branch'
- (oracle) Executing command (/work/project): 'fossil' 'branch' 'list'
- (oracle) Executing command (/work/project): 'fossil' 'tag' 'list'
- (oracle) Executing command (/work/project): 'svn' 'info' '--xml'
- (oracle) Downloading http://127.0.0.1:9203/packages.json
- (oracle)   [Composer\Downloader\TransportException (100)]                               
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Util/Http/CurlDownloader.php:498
- (oracle)  Composer\Util\Http\CurlDownloader->tick() at /oracle-source/src/Composer/Util/HttpDownloader.php:406
- (oracle)  Composer\Util\HttpDownloader->countActiveJobs() at /oracle-source/src/Composer/Util/HttpDownloader.php:379
- (oracle)  Composer\Util\HttpDownloader->wait() at /oracle-source/src/Composer/Util/HttpDownloader.php:112
- (oracle)  Composer\Util\HttpDownloader->get() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1808
- (oracle)  Composer\Repository\ComposerRepository->fetchFile() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1514
- (oracle)  Composer\Repository\ComposerRepository->loadRootServerFile() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1136
- (oracle)  Composer\Repository\ComposerRepository->hasProviders() at /oracle-source/src/Composer/Repository/ComposerRepository.php:571
- (oracle)  Composer\Repository\ComposerRepository->loadPackages() at /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  Composer\DependencyResolver\PoolBuilder->loadPackagesMarkedForLoading() at /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  Composer\DependencyResolver\PoolBuilder->buildPool() at /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  Composer\Repository\RepositorySet->createPool() at /oracle-source/src/Composer/Installer.php:537
- (oracle)  Composer\Installer->doUpdate() at /oracle-source/src/Composer/Installer.php:327
- (oracle)  Composer\Installer->run() at /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  Composer\Command\UpdateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## root-tty-help-yes

```
stdout:
+ (go)   validate [options] [<file>]
- (oracle)   validate [options] [--] [<file>]
+ (go)   \e[32mfile                        \e[39m path to composer.json file
- (oracle)   \e[32mfile\e[39m                           path to composer.json file
+ (go)   \e[32m    --no-check-all          \e[39m Do not validate requires for overly strict/loose constraints
+ (go)   \e[32m    --check-lock            \e[39m Check if lock file is up to date (even when config.lock is false)
+ (go)   \e[32m    --no-check-lock         \e[39m Do not check if lock file is up to date
+ (go)   \e[32m    --no-check-publish      \e[39m Do not check for publish errors
+ (go)   \e[32m    --no-check-version      \e[39m Do not report a warning if the version field is present
+ (go)   \e[32m-A, --with-dependencies     \e[39m Also validate the composer.json of all installed dependencies
+ (go)   \e[32m    --strict                \e[39m Return a non-zero exit code for warnings as well as errors
+ (go)   \e[32m-h, --help                  \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet                 \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose      \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version               \e[39m Display this application version
+ (go)   \e[32m    --ansi                  \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction        \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile               \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins            \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts            \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir=WORKING-DIR\e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache              \e[39m Prevent use of the cache
- (oracle)   \e[32m    --no-check-all\e[39m             Do not validate requires for overly strict/loose constraints
- (oracle)   \e[32m    --check-lock\e[39m               Check if lock file is up to date (even when config.lock is false)
- (oracle)   \e[32m    --no-check-lock\e[39m            Do not check if lock file is up to date
- (oracle)   \e[32m    --no-check-publish\e[39m         Do not check for publish errors
- (oracle)   \e[32m    --no-check-version\e[39m         Do not report a warning if the version field is present
- (oracle)   \e[32m-A, --with-dependencies\e[39m        Also validate the composer.json of all installed dependencies
- (oracle)   \e[32m    --strict\e[39m                   Return a non-zero exit code for warnings as well as errors
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
- (oracle) 
- (oracle) \e[33mHelp:\e[39m
- (oracle)   The validate command validates a given composer.json and composer.lock
- (oracle)   
- (oracle)   Exit codes in case of errors are:
- (oracle)   1 validation warning(s), only when --strict is given
- (oracle)   2 validation error(s)
- (oracle)   3 file unreadable or missing
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#validate
```
