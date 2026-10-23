# Composer vs gocomposer comparison

- gocomposer_sha256: `bdd689d43a9bc2b632a3c0d7fc2ea757ee9c36c8625490df82d7a820f76cb1ab`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| bump-lock-disabled | KNOWN-DIFF | 1 | 1 |
| bump-after-update-dry-run | MATCH | 1 | 1 |

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
