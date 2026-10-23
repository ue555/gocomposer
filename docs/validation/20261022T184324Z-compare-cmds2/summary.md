# Composer vs gocomposer comparison

- gocomposer_sha256: `72ff76c24eae6d632d9a3849548bf3666f6f3610c97efa24ba131a4fdbac6798`
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
| bump-lock-disabled | KNOWN-DIFF | 1 | 1 |
| bump-installed-only | MATCH | 0 | 0 |
| bump-not-installed | MATCH | 0 | 0 |
| bump-branch-alias | MATCH | 0 | 0 |
| bump-no-json | MATCH | 1 | 1 |
| bump-after-update-flag | MATCH | 0 | 0 |
| bump-after-update-dev | MATCH | 0 | 0 |
| bump-after-update-no-dev | MATCH | 0 | 0 |
| bump-after-update-dry-run | MATCH | 1 | 1 |
| bump-after-update-config | MATCH | 0 | 0 |
| bump-after-update-config-dev | MATCH | 0 | 0 |
| bump-after-update-untyped | MATCH | 0 | 0 |
| bump-after-update-partial | MATCH | 0 | 0 |
| bump-after-update-nothing | MATCH | 0 | 0 |
| bump-after-update-lock | MATCH | 0 | 0 |
| init-full | MATCH | 0 | 0 |
| init-full-reordered | MATCH | 0 | 0 |
| init-require | MATCH | 0 | 0 |
| init-require-pair-args | MATCH | 0 | 0 |
| init-require-no-version | MATCH | 1 | 1 |
| init-require-dev-no-version | MATCH | 1 | 1 |
| init-autoload | MATCH | 0 | 0 |
| init-autoload-nested | MATCH | 0 | 0 |
| init-autoload-with-deps | MATCH | 0 | 0 |
| init-autoload-dev-deps | MATCH | 0 | 0 |
| init-git-defaults | MATCH | 0 | 0 |
| init-git-defaults-description | MATCH | 0 | 0 |
| init-env-defaults | MATCH | 0 | 0 |
| init-env-overrides-git | MATCH | 0 | 0 |
| init-user-env | MATCH | 0 | 0 |
| init-no-defaults | MATCH | 0 | 0 |
| init-camel-dir | MATCH | 0 | 0 |
| init-invalid-name | MATCH | 1 | 1 |
| init-invalid-name-noslash | MATCH | 1 | 1 |
| init-invalid-author | MATCH | 1 | 1 |
| init-invalid-email | MATCH | 1 | 1 |
| init-empty-author | MATCH | 1 | 1 |
| init-author-name-only | MATCH | 0 | 0 |
| init-existing-json | MATCH | 0 | 0 |
| init-existing-git | MATCH | 0 | 0 |
| init-repository | MATCH | 0 | 0 |
| init-repository-invalid | MATCH | 1 | 1 |
| init-repository-bad-json | MATCH | 1 | 1 |
| init-repository-composer-file | MATCH | 0 | 0 |
| init-repository-filesystem-file | MATCH | 1 | 1 |
| init-repository-missing-file | MATCH | 1 | 1 |
| init-bad-stability | MATCH | 1 | 1 |
| init-bad-type-schema | MATCH | 1 | 1 |
| init-bad-license | MATCH | 0 | 0 |
| init-composer-env | MATCH | 0 | 0 |
| init-composer-env-subdir | MATCH | 0 | 0 |
| init-unicode | MATCH | 0 | 0 |
| global-require | MATCH | 0 | 0 |
| global-update | MATCH | 0 | 0 |
| global-show | MATCH | 0 | 0 |
| global-show-abbrev | MATCH | 0 | 0 |
| global-options-first | MATCH | 0 | 0 |
| global-remove | MATCH | 0 | 0 |
| global-no-args | MATCH | 1 | 1 |
| global-unknown-option-only | DIFF | 1 | 1 |
| global-unknown-command | MATCH | 1 | 1 |
| global-config-quoted | MATCH | 0 | 0 |
| global-config-equals | MATCH | 0 | 0 |
| global-new-home | MATCH | 0 | 0 |
| global-composer-env | MATCH | 0 | 0 |
| global-validate | MATCH | 2 | 2 |
| global-dump-autoload | MATCH | 0 | 0 |
| archive-root-tar | MATCH | 0 | 0 |
| archive-root-zip | MATCH | 0 | 0 |
| archive-root-tar-gz | MATCH | 0 | 0 |
| archive-root-tar-bz2 | MATCH | 1 | 1 |
| archive-root-bad-format | MATCH | 1 | 1 |
| archive-root-empty-format | MATCH | 1 | 1 |
| archive-root-dir | MATCH | 0 | 0 |
| archive-root-dir-relative | MATCH | 0 | 0 |
| archive-root-file | MATCH | 0 | 0 |
| archive-root-file-subdir | MATCH | 0 | 0 |
| archive-root-ignore-filters | MATCH | 0 | 0 |
| archive-root-named | MATCH | 0 | 0 |
| archive-root-installed | MATCH | 0 | 0 |
| archive-root-existing-archives | MATCH | 0 | 0 |
| archive-root-config-format | MATCH | 0 | 0 |
| archive-root-dir-is-file | MATCH | 1 | 1 |
| archive-root-scripts | MATCH | 0 | 0 |
| archive-package | MATCH | 0 | 0 |
| archive-package-installed | MATCH | 0 | 0 |
| archive-package-constraint | MATCH | 0 | 0 |
| archive-package-stability | MATCH | 0 | 0 |
| archive-package-stability-filtered | MATCH | 1 | 1 |
| archive-package-missing | MATCH | 1 | 1 |
| archive-package-file | MATCH | 1 | 1 |
| archive-package-case | MATCH | 0 | 0 |
| archive-no-json | MATCH | 1 | 1 |
| archive-empty-tar | MATCH | 0 | 0 |
| archive-empty-zip | MATCH | 0 | 0 |
| archive-empty-tar-gz | MATCH | 0 | 0 |
| archive-only-empty-dirs | MATCH | 0 | 0 |
| archive-globstar-attributes | MATCH | 1 | 1 |
| archive-globstar-valid | MATCH | 0 | 0 |
| archive-vendor-dir | MATCH | 0 | 0 |

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

## global-unknown-option-only

```
stderr:
+ (go)                                                    
+ (go)   Not enough arguments (missing: "command-name").  
+ (go)                                                    
- (oracle)                                       
- (oracle)   The "--foo" option does not exist.  
- (oracle)                                       
```
