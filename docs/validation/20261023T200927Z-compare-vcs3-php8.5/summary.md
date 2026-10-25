# Composer vs gocomposer comparison

- gocomposer_sha256: `8a8f8d5fb28828fbc3b57432e29d52bc8926fe3ba0db114d8ba3b1345571d910`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.5`
- oracle_image_digest: `sha256:a9ce5d9c43bbc3ee3f07c81df52b3ffa6cf2de4beb6c803cf9ab8e3af016104f`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.5-zip`
- oracle_zip_image_digest: `sha256:9932fd56504e6222ca9dda0b0fc4d094e940ea6daaa0e18cfb89f98a2a24dcba`
- php_version: `8.5.12`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| vcs-svn-update | DIFF | 0 | 0 |
| vcs-svn-update-v | DIFF | 0 | 0 |
| vcs-svn-update-vv | DIFF | 0 | 0 |
| vcs-svn-detect | MATCH | 1 | 1 |
| vcs-svn-detect-serve | MATCH | 0 | 0 |
| vcs-svn-abs-path | DIFF | 0 | 0 |
| vcs-svn-trunk-url | DIFF | 0 | 0 |
| vcs-svn-dev-trunk | DIFF | 0 | 0 |
| vcs-svn-branch | DIFF | 0 | 0 |
| vcs-svn-show | MATCH | 0 | 0 |
| vcs-svn-package-path | DIFF | 0 | 0 |
| vcs-svn-no-tags-path | DIFF | 0 | 0 |
| vcs-svn-custom-trunk | DIFF | 0 | 0 |
| vcs-svn-missing | MATCH | 1 | 1 |
| vcs-svn-cache | DIFF | 0 | 0 |
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
| source-svn-install-v | DIFF | 0 | 0 |
| source-svn-update | DIFF | 0 | 0 |
| source-svn-update-v | KNOWN-DIFF | 1 | 1 |
| source-svn-to-trunk | DIFF | 0 | 0 |
| source-svn-local-changes | MATCH | 1 | 1 |
| source-svn-discard-true | DIFF | 0 | 0 |
| source-svn-discard-stash | MATCH | 1 | 1 |
| source-svn-uninstall-changes | MATCH | 1 | 1 |
| source-svn-uninstall-discard | MATCH | 0 | 0 |
| source-svn-remove | MATCH | 0 | 0 |
| source-svn-status-clean | MATCH | 0 | 0 |
| source-svn-status | MATCH | 1 | 1 |
| source-svn-status-v | MATCH | 1 | 1 |
| source-svn-tty-discard-yes | DIFF | 0 | 0 |
| source-svn-tty-discard-no | DIFF | 1 | 1 |
| source-svn-create-project | DIFF | 0 | 0 |
| source-svn-create-project-keep-vcs | DIFF | 0 | 0 |
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
| vcs-perforce-update | DIFF | 0 | 0 |
| vcs-perforce-update-vv | DIFF | 0 | 0 |
| vcs-perforce-detect | DIFF | 0 | 0 |
| vcs-perforce-detect-p4-url | MATCH | 1 | 1 |
| vcs-perforce-dev | DIFF | 0 | 0 |
| vcs-perforce-show | DIFF | 0 | 0 |
| vcs-perforce-env-user | DIFF | 0 | 0 |
| vcs-perforce-login | DIFF | 0 | 0 |
| vcs-perforce-login-wrong | MATCH | 1 | 1 |
| vcs-perforce-login-env | DIFF | 0 | 0 |
| vcs-perforce-stream | MATCH | 0 | 0 |
| vcs-perforce-no-depot | MATCH | 1 | 1 |
| vcs-perforce-bad-server | MATCH | 1 | 1 |
| vcs-perforce-invalid-port | MATCH | 1 | 1 |
| source-perforce-install | MATCH | 0 | 0 |
| source-perforce-update | DIFF | 0 | 0 |
| source-perforce-update-v | DIFF | 0 | 0 |
| source-perforce-status | MATCH | 0 | 0 |
| source-perforce-remove | MATCH | 0 | 0 |

## vcs-svn-update

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-svn-update-v

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-svn-update-vv

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-svn-abs-path

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-svn-trunk-url

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-svn-dev-trunk

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-svn-branch

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-svn-package-path

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-svn-no-tags-path

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-svn-custom-trunk

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-svn-cache

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## source-svn-install-v

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:61
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:329
- (oracle)  /oracle-source/src/Composer/Command/InstallCommand.php:153
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## source-svn-update

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-svn-update-v

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:387
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:274
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/SvnDriver.php:89
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:149
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:203
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:103
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:85
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:174
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:319
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:313
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:177
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:502
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:289
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:242
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:222
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:129
- (oracle)  /oracle-source/src/Composer/Util/Svn.php:120
- (oracle)  /oracle-source/src/Composer/Downloader/SvnDownloader.php:221
- (oracle)  /oracle-source/src/Composer/Downloader/VcsDownloader.php:195
- (oracle)  /oracle-source/src/Composer/Downloader/DownloadManager.php:319
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:313
- (oracle)  /oracle-source/src/Composer/Installer/LibraryInstaller.php:177
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:502
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  /oracle-source/src/Composer/Installer.php:885
- (oracle)  /oracle-source/src/Composer/Installer.php:689
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle)   [RuntimeException]                                                           
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Downloader/SvnDownloader.php:223
- (oracle)  Composer\Downloader\SvnDownloader->getCommitLogs() at /oracle-source/src/Composer/Downloader/VcsDownloader.php:195
- (oracle)  Composer\Downloader\VcsDownloader->update() at /oracle-source/src/Composer/Downloader/DownloadManager.php:319
- (oracle)  Composer\Downloader\DownloadManager->update() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:313
- (oracle)  Composer\Installer\LibraryInstaller->updateCode() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:177
- (oracle)  Composer\Installer\LibraryInstaller->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:502
- (oracle)  Composer\Installer\InstallationManager->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  Composer\Installer\InstallationManager->{closure:Composer\Installer\InstallationManager::executeBatch():385}() at /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
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

## source-svn-to-trunk

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-svn-discard-true

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-svn-tty-discard-yes

```
stdout:
+ (go) \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08                                                      \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\e[30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.\e[39;49m
- (oracle) \e[30;43mDeprecation Notice: Using null as an array offset is deprecated, use an empty string instead in /oracle-source/src/Composer/Util/Svn.php:318\e[39;49m
- (oracle) \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08                                                                                                                                            \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\e[30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.\e[39;49m
- (oracle) \e[30;43mMore deprecation notices were hidden, run again with `-v` to show them.\e[39;49m
```

## source-svn-tty-discard-no

```
stdout:
+ (go) \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08                                                      \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\e[30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.\e[39;49m
- (oracle) \e[30;43mDeprecation Notice: Using null as an array offset is deprecated, use an empty string instead in /oracle-source/src/Composer/Util/Svn.php:318\e[39;49m
- (oracle) \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08                                                                                                                                            \x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\e[30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.\e[39;49m
```

## source-svn-create-project

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-svn-create-project-keep-vcs

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-update

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-update-vv

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-perforce-detect

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## vcs-perforce-dev

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-show

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-env-user

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-login

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## vcs-perforce-login-env

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-perforce-update

```
stderr:
- (oracle) More deprecation notices were hidden, run again with `-v` to show them.
```

## source-perforce-update-v

```
stderr:
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Repository/Vcs/PerforceDriver.php:100
- (oracle)  /oracle-source/src/Composer/Repository/VcsRepository.php:337
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:308
- (oracle)  /oracle-source/src/Composer/Repository/ArrayRepository.php:62
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  /oracle-source/src/Composer/Installer.php:537
- (oracle)  /oracle-source/src/Composer/Installer.php:327
- (oracle)  /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```
