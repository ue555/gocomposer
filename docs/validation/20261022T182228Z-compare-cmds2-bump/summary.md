# Composer vs gocomposer comparison

- gocomposer_sha256: `e6a0d9e271a6c37be9ff2fc01112d5a4100a720c732e083e4299efe54a079240`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| bump-basic | MATCH | 0 | 0 |
| bump-dry-run | MATCH | 1 | 1 |
| bump-dev-only | MATCH | 0 | 0 |
| bump-no-dev-only | MATCH | 0 | 0 |
| bump-library | MATCH | 0 | 0 |
| bump-library-dev-only | MATCH | 0 | 0 |
| bump-untyped | MATCH | 0 | 0 |
| bump-filter | MATCH | 0 | 0 |
| bump-filter-wildcard | MATCH | 0 | 0 |
| bump-nothing | MATCH | 0 | 0 |
| bump-nothing-dry-run | MATCH | 0 | 0 |
| bump-lock-outdated | MATCH | 2 | 2 |
| bump-lock-disabled | DIFF | 1 | 1 |
| bump-installed-only | MATCH | 0 | 0 |
| bump-not-installed | MATCH | 0 | 0 |
| bump-branch-alias | MATCH | 0 | 0 |
| bump-no-json | MATCH | 1 | 1 |
| bump-after-update-flag | MATCH | 0 | 0 |
| bump-after-update-dev | MATCH | 0 | 0 |
| bump-after-update-no-dev | MATCH | 0 | 0 |
| bump-after-update-dry-run | DIFF | 1 | 1 |
| bump-after-update-config | MATCH | 0 | 0 |
| bump-after-update-config-dev | MATCH | 0 | 0 |
| bump-after-update-untyped | MATCH | 0 | 0 |
| bump-after-update-partial | MATCH | 0 | 0 |
| bump-after-update-nothing | MATCH | 0 | 0 |
| bump-after-update-lock | MATCH | 0 | 0 |

## bump-lock-disabled

```
stderr:
- (oracle)   [LogicException]                                   
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Package/Locker.php:342
- (oracle)  Composer\Package\Locker->getLockData() at /oracle-source/src/Composer/Package/Locker.php:168
- (oracle)  Composer\Package\Locker->getLockedRepository() at /oracle-source/src/Composer/Command/BumpCommand.php:126
- (oracle)  Composer\Command\BumpCommand->doBump() at /oracle-source/src/Composer/Command/BumpCommand.php:75
- (oracle)  Composer\Command\BumpCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## bump-after-update-dry-run

```
stdout:
+ (go) ./composer.json would be updated with:
+ (go)  - require.acme/log: ^1.1
+ (go)  - require.acme/http: ^1.2
+ (go)  - require.acme/util: >=1.1
stderr:
- (oracle) In <location>:
- (oracle)                                                
- (oracle)   Could not read ./composer.lock               
- (oracle)                                                
- (oracle)   The file "./composer.lock" is not readable.  
- (oracle)                                                
- (oracle) 
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
- (oracle) 
```
