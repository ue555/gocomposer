# Composer vs gocomposer comparison

- gocomposer_sha256: `8a8f8d5fb28828fbc3b57432e29d52bc8926fe3ba0db114d8ba3b1345571d910`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.3`
- oracle_image_digest: `sha256:2b1b802373a38be255fdd8574e7461226aeaccbd1e06295cd90c4a0616ab4280`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.3-zip`
- oracle_zip_image_digest: `sha256:f7513be70340918a49f0c6b8ae972d068b0422f51fbb77268d7303a9de4aebd6`
- php_version: `8.3.36`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| vcs-svn-update | MATCH | 0 | 0 |
| vcs-svn-update-v | MATCH | 0 | 0 |
| vcs-svn-update-vv | MATCH | 0 | 0 |
| vcs-svn-detect | MATCH | 1 | 1 |
| vcs-svn-detect-serve | MATCH | 0 | 0 |
| vcs-svn-abs-path | MATCH | 0 | 0 |
| vcs-svn-trunk-url | MATCH | 0 | 0 |
| vcs-svn-dev-trunk | MATCH | 0 | 0 |
| vcs-svn-branch | MATCH | 0 | 0 |
| vcs-svn-show | MATCH | 0 | 0 |
| vcs-svn-package-path | MATCH | 0 | 0 |
| vcs-svn-no-tags-path | MATCH | 0 | 0 |
| vcs-svn-custom-trunk | MATCH | 0 | 0 |
| vcs-svn-missing | MATCH | 1 | 1 |
| vcs-svn-cache | MATCH | 0 | 0 |
| vcs-svn-serve | MATCH | 0 | 0 |
| vcs-svn-serve-secure | MATCH | 1 | 1 |
| vcs-svn-serve-secure-domain | MATCH | 0 | 0 |
| vcs-svn-auth-noninteractive | MATCH | 1 | 1 |
| vcs-svn-auth-http-basic | MATCH | 0 | 0 |
| vcs-svn-auth-wrong | MATCH | 1 | 1 |
| vcs-svn-auth-url | MATCH | 0 | 0 |
| vcs-svn-auth-no-cache | MATCH | 0 | 0 |
| vcs-svn-auth-tty | MATCH | 0 | 0 |
| vcs-svn-auth-tty-wrong | MATCH | 1 | 1 |
| vcs-svn-path-guess | MATCH | 0 | 0 |
| vcs-svn-path-guess-branch | MATCH | 0 | 0 |
| source-svn-install | MATCH | 0 | 0 |
| source-svn-install-v | MATCH | 0 | 0 |
| source-svn-update | MATCH | 0 | 0 |
| source-svn-update-v | KNOWN-DIFF | 1 | 1 |
| source-svn-to-trunk | MATCH | 0 | 0 |
| source-svn-local-changes | MATCH | 1 | 1 |
| source-svn-discard-true | MATCH | 0 | 0 |
| source-svn-discard-stash | MATCH | 1 | 1 |
| source-svn-uninstall-changes | MATCH | 1 | 1 |
| source-svn-uninstall-discard | MATCH | 0 | 0 |
| source-svn-remove | MATCH | 0 | 0 |
| source-svn-status-clean | MATCH | 0 | 0 |
| source-svn-status | MATCH | 1 | 1 |
| source-svn-status-v | MATCH | 1 | 1 |
| source-svn-tty-discard-yes | MATCH | 0 | 0 |
| source-svn-tty-discard-no | MATCH | 1 | 1 |
| source-svn-create-project | MATCH | 0 | 0 |
| source-svn-create-project-keep-vcs | MATCH | 0 | 0 |
| source-svn-lock-reference | MATCH | 0 | 0 |
| vcs-hg-update | MATCH | 0 | 0 |
| vcs-hg-update-vv | MATCH | 0 | 0 |
| vcs-hg-detect | MATCH | 0 | 0 |
| vcs-hg-dev-default | MATCH | 0 | 0 |
| vcs-hg-branch | MATCH | 0 | 0 |
| vcs-hg-bookmark | MATCH | 0 | 0 |
| vcs-hg-show | MATCH | 0 | 0 |
| vcs-hg-serve | MATCH | 0 | 0 |
| vcs-hg-serve-vv | MATCH | 0 | 0 |
| vcs-hg-serve-detect | MATCH | 0 | 0 |
| vcs-hg-serve-cached | MATCH | 0 | 0 |
| vcs-hg-serve-secure | MATCH | 100 | 100 |
| vcs-hg-serve-missing | MATCH | 1 | 1 |
| vcs-hg-serve-auth | MATCH | 1 | 1 |
| vcs-hg-missing | MATCH | 1 | 1 |
| vcs-hg-path-guess | MATCH | 0 | 0 |
| vcs-hg-path-guess-feature | MATCH | 0 | 0 |
| source-hg-create-project | MATCH | 0 | 0 |
| source-hg-install | MATCH | 0 | 0 |
| source-hg-install-serve | MATCH | 0 | 0 |
| source-hg-update | MATCH | 0 | 0 |
| source-hg-update-v | MATCH | 0 | 0 |
| source-hg-downgrade-v | MATCH | 0 | 0 |
| source-hg-to-dev | MATCH | 0 | 0 |
| source-hg-local-changes | MATCH | 1 | 1 |
| source-hg-discard-true | MATCH | 1 | 1 |
| source-hg-uninstall-changes | MATCH | 1 | 1 |
| source-hg-remove | MATCH | 0 | 0 |
| source-hg-status-clean | MATCH | 0 | 0 |
| source-hg-status | MATCH | 1 | 1 |
| source-hg-status-v | MATCH | 1 | 1 |
| vcs-fossil-update | MATCH | 0 | 0 |
| vcs-fossil-update-vv | MATCH | 0 | 0 |
| vcs-fossil-detect | MATCH | 0 | 0 |
| vcs-fossil-checkout-dir | MATCH | 0 | 0 |
| vcs-fossil-dev-trunk | MATCH | 0 | 0 |
| vcs-fossil-branch | MATCH | 0 | 0 |
| vcs-fossil-show | MATCH | 0 | 0 |
| vcs-fossil-serve | MATCH | 1 | 1 |
| vcs-fossil-no-repo-cache-dir | MATCH | 1 | 1 |
| vcs-fossil-serve-secure | MATCH | 100 | 100 |
| vcs-fossil-missing | MATCH | 1 | 1 |
| vcs-fossil-path-guess | MATCH | 0 | 0 |
| source-fossil-no-user | MATCH | 1 | 1 |
| source-fossil-create-project | MATCH | 0 | 0 |
| source-fossil-install | MATCH | 0 | 0 |
| source-fossil-update | MATCH | 0 | 0 |
| source-fossil-update-v | MATCH | 0 | 0 |
| source-fossil-to-dev | MATCH | 0 | 0 |
| source-fossil-local-changes | MATCH | 1 | 1 |
| source-fossil-uninstall-changes | MATCH | 1 | 1 |
| source-fossil-remove | MATCH | 0 | 0 |
| source-fossil-status-clean | MATCH | 0 | 0 |
| source-fossil-status-v | MATCH | 1 | 1 |
| vcs-perforce-update | MATCH | 0 | 0 |
| vcs-perforce-update-vv | MATCH | 0 | 0 |
| vcs-perforce-detect | MATCH | 0 | 0 |
| vcs-perforce-detect-p4-url | MATCH | 1 | 1 |
| vcs-perforce-dev | MATCH | 0 | 0 |
| vcs-perforce-show | MATCH | 0 | 0 |
| vcs-perforce-env-user | MATCH | 0 | 0 |
| vcs-perforce-login | MATCH | 0 | 0 |
| vcs-perforce-login-wrong | MATCH | 1 | 1 |
| vcs-perforce-login-env | MATCH | 0 | 0 |
| vcs-perforce-stream | MATCH | 0 | 0 |
| vcs-perforce-no-depot | MATCH | 1 | 1 |
| vcs-perforce-bad-server | MATCH | 1 | 1 |
| vcs-perforce-invalid-port | MATCH | 1 | 1 |
| source-perforce-install | MATCH | 0 | 0 |
| source-perforce-update | MATCH | 0 | 0 |
| source-perforce-update-v | MATCH | 0 | 0 |
| source-perforce-status | MATCH | 0 | 0 |
| source-perforce-remove | MATCH | 0 | 0 |

## source-svn-update-v

```
stderr:
- (oracle)   [RuntimeException]                                                           
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Downloader/SvnDownloader.php:223
- (oracle)  Composer\Downloader\SvnDownloader->getCommitLogs() at /oracle-source/src/Composer/Downloader/VcsDownloader.php:195
- (oracle)  Composer\Downloader\VcsDownloader->update() at /oracle-source/src/Composer/Downloader/DownloadManager.php:319
- (oracle)  Composer\Downloader\DownloadManager->update() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:313
- (oracle)  Composer\Installer\LibraryInstaller->updateCode() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:177
- (oracle)  Composer\Installer\LibraryInstaller->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:502
- (oracle)  Composer\Installer\InstallationManager->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  Composer\Installer\InstallationManager->Composer\Installer\{closure}() at /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  React\Promise\Internal\FulfilledPromise->then() at /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  Composer\Installer\InstallationManager->executeBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  Composer\Installer\InstallationManager->downloadAndExecuteBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  Composer\Installer\InstallationManager->execute() at /oracle-source/src/Composer/Installer.php:885
- (oracle)  Composer\Installer->doInstall() at /oracle-source/src/Composer/Installer.php:689
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
