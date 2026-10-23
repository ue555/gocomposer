# Composer vs gocomposer comparison

- gocomposer_sha256: `bb7a455bddda59b1d0c8f5ae9d5d826b67beebabb06b5b61ac75b60144401bb1`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| config-platform-php | MATCH | 0 | 0 |
| config-no-args | MATCH | 0 | 0 |
| config-get-vendor-dir | MATCH | 0 | 0 |
| config-get-bin-dir | MATCH | 0 | 0 |
| config-get-bin-dir-absolute | MATCH | 0 | 0 |
| config-get-bin-dir-source | MATCH | 0 | 0 |
| config-get-set-in-json | MATCH | 0 | 0 |
| config-get-bin-dir-custom-vendor | MATCH | 0 | 0 |
| config-get-process-timeout | MATCH | 0 | 0 |
| config-get-bool | MATCH | 0 | 0 |
| config-get-bool-default | MATCH | 0 | 0 |
| config-get-allow-plugins-empty | MATCH | 0 | 0 |
| config-get-platform-empty | MATCH | 0 | 0 |
| config-get-secure-svn-domains | MATCH | 0 | 0 |
| config-get-github-protocols | MATCH | 0 | 0 |
| config-get-github-protocols-secure | MATCH | 0 | 0 |
| config-get-github-protocols-git-first | MATCH | 0 | 0 |
| config-get-gitlab-domains-dedupe | MATCH | 0 | 0 |
| config-get-audit | MATCH | 0 | 0 |
| config-get-audit-nested | MATCH | 0 | 0 |
| config-get-nested-missing | MATCH | 1 | 1 |
| config-get-extra-nested | MATCH | 0 | 0 |
| config-get-extra-object | MATCH | 0 | 0 |
| config-get-extra | MATCH | 0 | 0 |
| config-get-extra-list-index | MATCH | 0 | 0 |
| config-get-suggest | MATCH | 0 | 0 |
| config-get-name | MATCH | 0 | 0 |
| config-get-license | MATCH | 0 | 0 |
| config-get-homepage-default | MATCH | 0 | 0 |
| config-get-min-stability-default | MATCH | 0 | 0 |
| config-get-prefer-stable-default | MATCH | 0 | 0 |
| config-get-keywords | MATCH | 0 | 0 |
| config-get-repositories | MATCH | 0 | 0 |
| config-get-repositories-default | MATCH | 0 | 0 |
| config-get-repo-index | MATCH | 0 | 0 |
| config-get-repo-missing | MATCH | 1 | 1 |
| config-get-repo-prefix | MATCH | 0 | 0 |
| config-get-undefined | MATCH | 1 | 1 |
| config-get-undefined-dotted | MATCH | 1 | 1 |
| config-get-leading-dot | MATCH | 1 | 1 |
| config-get-cache-files-maxsize | MATCH | 0 | 0 |
| config-get-cache-dir-env | MATCH | 0 | 0 |
| config-get-env-vendor-dir | MATCH | 0 | 0 |
| config-get-home | MATCH | 0 | 0 |
| config-get-data-dir | MATCH | 0 | 0 |
| config-get-cafile-empty | MATCH | 0 | 0 |
| config-get-gitlab-protocol-null | MATCH | 1 | 1 |
| config-get-preferred-install-map | MATCH | 0 | 0 |
| config-get-policy | MATCH | 0 | 0 |
| config-get-policy-nested | MATCH | 0 | 0 |
| config-get-auth-file | MATCH | 0 | 0 |
| config-get-unicode-slashes | MATCH | 0 | 0 |
| config-get-scalar-unicode | MATCH | 0 | 0 |
| config-get-scalar-float | MATCH | 0 | 0 |
| config-list | MATCH | 0 | 0 |
| config-list-source | MATCH | 0 | 0 |
| config-list-plain | MATCH | 0 | 0 |
| config-list-custom | DIFF | 1 | 0 |
| config-list-named-repos | DIFF | 0 | 0 |
| config-list-env | MATCH | 0 | 0 |
| config-list-bad-maxsize | MATCH | 1 | 1 |
| config-list-global | MATCH | 0 | 0 |
| config-list-auth | MATCH | 0 | 0 |
| config-set-sort-packages | MATCH | 0 | 0 |
| config-set-bool-invalid | MATCH | 1 | 1 |
| config-set-process-timeout | MATCH | 0 | 0 |
| config-set-process-timeout-invalid | MATCH | 1 | 1 |
| config-set-too-many | MATCH | 1 | 1 |
| config-set-vendor-dir | MATCH | 0 | 0 |
| config-set-preferred-install | MATCH | 0 | 0 |
| config-set-preferred-install-invalid | MATCH | 1 | 1 |
| config-set-preferred-install-pkg | MATCH | 0 | 0 |
| config-set-preferred-install-pkg-invalid | MATCH | 1 | 1 |
| config-set-store-auths | MATCH | 0 | 0 |
| config-set-store-auths-false | MATCH | 0 | 0 |
| config-set-discard-changes | MATCH | 0 | 0 |
| config-set-bump-after-update | MATCH | 0 | 0 |
| config-set-platform-check | MATCH | 0 | 0 |
| config-set-use-parent-dir | MATCH | 0 | 0 |
| config-set-cache-maxsize | MATCH | 0 | 0 |
| config-set-cache-maxsize-invalid | MATCH | 1 | 1 |
| config-set-bin-compat-invalid | MATCH | 1 | 1 |
| config-set-gitlab-protocol | MATCH | 0 | 0 |
| config-set-autoloader-suffix-null | DIFF | 1 | 1 |
| config-set-cafile-missing | MATCH | 1 | 1 |
| config-set-cafile | MATCH | 0 | 0 |
| config-set-capath-invalid | MATCH | 1 | 1 |
| config-set-github-protocols | MATCH | 0 | 0 |
| config-set-github-protocols-invalid | MATCH | 1 | 1 |
| config-set-github-domains | MATCH | 0 | 0 |
| config-set-disable-tls | MATCH | 0 | 0 |
| config-set-disable-tls-off | MATCH | 0 | 0 |
| config-unset-disable-tls | MATCH | 0 | 0 |
| config-unset-key | MATCH | 0 | 0 |
| config-unset-with-value | MATCH | 1 | 1 |
| config-set-unknown | MATCH | 1 | 1 |
| config-set-allow-plugins-bool | MATCH | 0 | 0 |
| config-set-lock-false | MATCH | 0 | 0 |
| config-set-audit-abandoned | MATCH | 0 | 0 |
| config-set-audit-abandoned-invalid | MATCH | 1 | 1 |
| config-set-audit-block-insecure | MATCH | 0 | 0 |
| config-set-audit-ignore | MATCH | 0 | 0 |
| config-set-audit-ignore-json | MATCH | 0 | 0 |
| config-set-audit-ignore-json-scalar | MATCH | 1 | 1 |
| config-set-audit-ignore-merge | MATCH | 0 | 0 |
| config-set-audit-ignore-merge-mixed | MATCH | 1 | 1 |
| config-set-audit-ignore-severity | MATCH | 0 | 0 |
| config-set-audit-ignore-severity-invalid | MATCH | 1 | 1 |
| config-unset-audit | MATCH | 0 | 0 |
| config-unset-audit-ignore | MATCH | 0 | 0 |
| config-set-policy-block | MATCH | 0 | 0 |
| config-set-policy-audit | MATCH | 0 | 0 |
| config-set-policy-audit-invalid | MATCH | 1 | 1 |
| config-set-policy-block-scope | MATCH | 0 | 0 |
| config-set-policy-block-scope-invalid | MATCH | 1 | 1 |
| config-set-policy-cooldown-period | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-int | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-null | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-invalid | MATCH | 1 | 1 |
| config-set-policy-list-off | MATCH | 0 | 0 |
| config-set-policy-list-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-block | MATCH | 0 | 0 |
| config-set-policy-custom-scope | MATCH | 0 | 0 |
| config-set-policy-custom-scope-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-audit-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-ignore | MATCH | 0 | 0 |
| config-set-policy-reserved | MATCH | 1 | 1 |
| config-set-policy-reserved-name | MATCH | 1 | 1 |
| config-set-policy-reserved-ignore | MATCH | 1 | 1 |
| config-set-policy-sources | MATCH | 1 | 1 |
| config-set-policy-unknown-field | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable | MATCH | 0 | 0 |
| config-set-policy-ignore-unreachable-bad | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable-json | MATCH | 0 | 0 |
| config-set-policy-ignore-unreachable-json-bool | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable-bool | MATCH | 0 | 0 |
| config-set-policy-ignore-id-json | MATCH | 0 | 0 |
| config-set-policy-ignore-merge | MATCH | 0 | 0 |
| config-set-policy-ignore-merge-mixed | MATCH | 1 | 1 |
| config-set-policy-ignore-severity | MATCH | 0 | 0 |
| config-set-policy-ignore-source | MATCH | 0 | 0 |
| config-unset-policy-nested | MATCH | 0 | 0 |
| config-unset-policy | MATCH | 0 | 0 |
| config-repo-add | MATCH | 0 | 0 |
| config-repo-add-append | MATCH | 0 | 0 |
| config-repo-add-plain | MATCH | 0 | 0 |
| config-repo-disable-packagist | MATCH | 0 | 0 |
| config-repo-add-json | MATCH | 0 | 0 |
| config-repo-add-invalid-json | MATCH | 1 | 1 |
| config-repo-add-true | MATCH | 1 | 1 |
| config-repo-add-three | MATCH | 1 | 1 |
| config-repo-remove | MATCH | 0 | 0 |
| config-repo-remove-assoc | DIFF | 0 | 0 |
| config-repo-get-named | MATCH | 0 | 0 |
| config-extra-set | MATCH | 0 | 0 |
| config-extra-set-new | MATCH | 0 | 0 |
| config-extra-json | MATCH | 0 | 0 |
| config-extra-json-merge-list | MATCH | 0 | 0 |
| config-extra-json-merge-object | MATCH | 0 | 0 |
| config-extra-json-merge-mixed | MATCH | 0 | 0 |
| config-extra-json-invalid | MATCH | 1 | 1 |
| config-extra-unset | MATCH | 0 | 0 |
| config-extra-unset-all | MATCH | 0 | 0 |
| config-suggest-set | MATCH | 0 | 0 |
| config-suggest-unset | MATCH | 0 | 0 |
| config-scripts-set | MATCH | 0 | 0 |
| config-scripts-set-one | MATCH | 0 | 0 |
| config-scripts-unset | MATCH | 0 | 0 |
| config-prop-name | MATCH | 0 | 0 |
| config-prop-description | MATCH | 0 | 0 |
| config-prop-min-stability | MATCH | 0 | 0 |
| config-prop-min-stability-invalid | MATCH | 1 | 1 |
| config-prop-prefer-stable | MATCH | 0 | 0 |
| config-prop-keywords | MATCH | 0 | 0 |
| config-prop-license | MATCH | 0 | 0 |
| config-prop-unset | MATCH | 0 | 0 |
| config-prop-unset-other | MATCH | 0 | 0 |
| config-prop-type-too-many | MATCH | 1 | 1 |
| config-platform-set | MATCH | 0 | 0 |
| config-platform-set-false | MATCH | 0 | 0 |
| config-platform-unset | MATCH | 0 | 0 |
| config-platform-unset-all | MATCH | 0 | 0 |
| config-allow-plugins-set | MATCH | 0 | 0 |
| config-allow-plugins-wildcard | MATCH | 0 | 0 |
| config-allow-plugins-invalid | MATCH | 1 | 1 |
| config-allow-plugins-unset | MATCH | 0 | 0 |
| config-auth-github-local | MATCH | 0 | 0 |
| config-auth-github-global | MATCH | 0 | 0 |
| config-auth-github-too-many | MATCH | 1 | 1 |
| config-auth-http-basic | MATCH | 0 | 0 |
| config-auth-http-basic-missing | MATCH | 1 | 1 |
| config-auth-bitbucket | MATCH | 0 | 0 |
| config-auth-bitbucket-missing | MATCH | 1 | 1 |
| config-auth-gitlab-token-pair | MATCH | 0 | 0 |
| config-auth-gitlab-token-three | MATCH | 1 | 1 |
| config-auth-gitlab-oauth | MATCH | 0 | 0 |
| config-auth-bearer | MATCH | 0 | 0 |
| config-auth-forgejo | MATCH | 0 | 0 |
| config-auth-forgejo-missing | MATCH | 1 | 1 |
| config-auth-custom-headers | MATCH | 0 | 0 |
| config-auth-custom-headers-invalid | MATCH | 1 | 1 |
| config-auth-unset | MATCH | 0 | 0 |
| config-auth-unset-missing-file | MATCH | 0 | 0 |
| config-auth-existing-in-json | MATCH | 0 | 0 |
| config-global-set | MATCH | 0 | 0 |
| config-global-set-modes | MATCH | 0 | 0 |
| config-global-get-existing | MATCH | 0 | 0 |
| config-global-list-existing | MATCH | 0 | 0 |
| config-global-merged-in-local | MATCH | 0 | 0 |
| config-global-dir-absolute | MATCH | 0 | 0 |
| config-global-property | MATCH | 1 | 1 |
| config-global-extra | MATCH | 1 | 1 |
| config-global-repo | MATCH | 0 | 0 |
| config-global-unset | MATCH | 0 | 0 |
| config-global-with-file | MATCH | 1 | 1 |
| config-global-no-json | MATCH | 0 | 0 |
| config-file-get | MATCH | 0 | 0 |
| config-file-set | MATCH | 0 | 0 |
| config-file-missing | MATCH | 1 | 1 |
| config-no-composer-json | MATCH | 1 | 1 |
| config-invalid-json | MATCH | 1 | 1 |
| config-composer-env | MATCH | 0 | 0 |
| config-editor | MATCH | 0 | 0 |
| config-editor-auth | MATCH | 0 | 0 |

## config-list-custom

```
exit code: oracle=1 go=0
stdout:
+ (go) [repositories.0.type] composer (./composer.json)
+ (go) [repositories.0.url] http://127.0.0.1:8899 (./composer.json)
+ (go) [process-timeout] 300 (./composer.json)
+ (go) [use-include-path] false (default)
+ (go) [allow-plugins.acme/plugin] true (./composer.json)
+ (go) [use-parent-dir] prompt (default)
+ (go) [preferred-install.acme/*] source (./composer.json)
+ (go) [preferred-install.*] dist (unknown)
+ (go) [audit.ignore] [PKSA-1] (./composer.json)
+ (go) [audit.abandoned] fail (default)
+ (go) [policy.cooldown.period] 7 days (./composer.json)
+ (go) [notify-on-install] true (default)
+ (go) [github-protocols] [git, https] (./composer.json)
+ (go) [gitlab-protocol]  (default)
+ (go) [vendor-dir] lib (/work/project/lib) (./composer.json)
+ (go) [bin-dir] {$vendor-dir}/bin (/work/project/lib/bin) (default)
+ (go) [cache-dir] /work/cache (COMPOSER_CACHE_DIR)
+ (go) [data-dir] /work/home (default)
+ (go) [cache-files-dir] {$cache-dir}/files (/work/cache/files) (default)
+ (go) [cache-repo-dir] {$cache-dir}/repo (/work/cache/repo) (default)
+ (go) [cache-vcs-dir] {$cache-dir}/vcs (/work/cache/vcs) (default)
+ (go) [cache-ttl] 15552000 (default)
+ (go) [cache-files-ttl] 15552000 (default)
+ (go) [cache-files-maxsize] 300MiB (314572800) (default)
+ (go) [cache-read-only] false (default)
+ (go) [bin-compat] auto (default)
+ (go) [discard-changes] false (default)
+ (go) [autoloader-suffix]  (default)
+ (go) [sort-packages] false (default)
+ (go) [optimize-autoloader] false (default)
+ (go) [strict-psr-autoloader] false (default)
+ (go) [classmap-authoritative] false (default)
+ (go) [apcu-autoloader] false (default)
+ (go) [prepend-autoloader] true (default)
+ (go) [update-with-minimal-changes] false (default)
+ (go) [github-domains] [github.com] (default)
+ (go) [bitbucket-expose-hostname] true (default)
+ (go) [disable-tls] false (default)
+ (go) [secure-http] false (./composer.json)
+ (go) [cafile]  (default)
+ (go) [capath]  (default)
+ (go) [github-expose-hostname] true (default)
+ (go) [gitlab-domains] [gitlab.com] (default)
+ (go) [store-auths] prompt (default)
+ (go) [platform.php] 8.1.0 (./composer.json)
+ (go) [archive-format] tar (default)
+ (go) [archive-dir] . (default)
+ (go) [htaccess-protect] true (default)
+ (go) [use-github-api] true (default)
+ (go) [lock] true (default)
+ (go) [platform-check] php-only (default)
+ (go) [bump-after-update] false (default)
+ (go) [allow-missing-requirements] false (default)
+ (go) [forgejo-domains] [codeberg.org] (default)
+ (go) [source-fallback] false (default)
+ (go) [home] /work/home (default)
stderr:
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   [TypeError]                                                                  
- (oracle)   Composer\Util\ProcessExecutor::setTimeout(): Argument #1 ($timeout) must be  
- (oracle)    of type int, string given, called in /oracle-source/src/Composer/IO/BaseIO  
- (oracle)   .php on line 212                                                             
- (oracle)                                                                                
- (oracle) 
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Util/ProcessExecutor.php:463
- (oracle)  Composer\Util\ProcessExecutor::setTimeout() at /oracle-source/src/Composer/IO/BaseIO.php:212
- (oracle)  Composer\IO\BaseIO->loadConfiguration() at /oracle-source/src/Composer/Command/ConfigCommand.php:214
- (oracle)  Composer\Command\ConfigCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
- (oracle) config [-g|--global] [-e|--editor] [-a|--auth] [--unset] [-l|--list] [-f|--file FILE] [--absolute] [-j|--json] [-m|--merge] [--append] [--source] [--] [<setting-key> [<setting-value>...]]
- (oracle) 
- (oracle) 
```

## config-list-named-repos

```
stderr:
- (oracle) Composer could not detect the root package (test/plain) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
```

## config-set-autoloader-suffix-null

```
stderr:
- (oracle) In <location>:
- (oracle)                                                              
- (oracle)   "./composer.json" does not match the expected JSON schema  
- (oracle)                                                              
- (oracle) 
```

## config-repo-remove-assoc

```
stderr:
- (oracle) Composer could not detect the root package (test/plain) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
```
