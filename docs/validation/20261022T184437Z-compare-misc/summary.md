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
| run-script-basic | MATCH | 0 | 0 |
| run-script-alias-run | MATCH | 0 | 0 |
| run-script-abbrev | MATCH | 0 | 0 |
| run-script-args | MATCH | 0 | 0 |
| run-script-args-dashdash | MATCH | 0 | 0 |
| run-script-list | MATCH | 0 | 0 |
| run-script-list-short | MATCH | 0 | 0 |
| run-script-missing-argument | MATCH | 1 | 1 |
| run-script-not-defined | MATCH | 1 | 1 |
| run-script-event-constant-name | MATCH | 1 | 1 |
| run-script-event | MATCH | 0 | 0 |
| run-script-event-array | MATCH | 0 | 0 |
| run-script-package-event | MATCH | 0 | 0 |
| run-script-dev-mode | MATCH | 0 | 0 |
| run-script-no-dev | MATCH | 0 | 0 |
| run-script-dev-and-no-dev | MATCH | 0 | 0 |
| run-script-no-dev-env | MATCH | 0 | 0 |
| run-script-timeout-invalid | MATCH | 1 | 1 |
| run-script-timeout-negative | MATCH | 1 | 1 |
| run-script-timeout | MATCH | 1 | 1 |
| run-script-timeout-zero | MATCH | 0 | 0 |
| run-script-process-timeout-env | MATCH | 1 | 1 |
| run-script-fail | MATCH | 3 | 3 |
| run-script-nested | MATCH | 0 | 0 |
| run-script-nested-verbose | MATCH | 0 | 0 |
| run-script-missing-ref | MATCH | 0 | 0 |
| run-script-multi-args | MATCH | 0 | 0 |
| run-script-no-additional-args | MATCH | 0 | 0 |
| run-script-putenv | MATCH | 0 | 0 |
| run-script-no-scripts | MATCH | 1 | 1 |
| run-script-skip-env | MATCH | 1 | 1 |
| run-script-skip-env-verbose | MATCH | 0 | 0 |
| run-script-quiet | MATCH | 0 | 0 |
| run-script-direct | MATCH | 0 | 0 |
| run-script-direct-args | MATCH | 0 | 0 |
| run-script-direct-unknown-option | MATCH | 0 | 0 |
| run-script-direct-dashdash | MATCH | 0 | 0 |
| run-script-direct-alias | MATCH | 0 | 0 |
| run-script-direct-alias2 | MATCH | 0 | 0 |
| run-script-direct-nested | MATCH | 0 | 0 |
| run-script-direct-fail | MATCH | 3 | 3 |
| run-script-direct-dev-mode | MATCH | 0 | 0 |
| run-script-direct-event-not-command | MATCH | 1 | 1 |
| run-script-direct-abbrev | MATCH | 0 | 0 |
| run-script-unknown-command | MATCH | 1 | 1 |
| run-script-override-warning | MATCH | 0 | 0 |
| run-script-aliases-invalid | MATCH | 1 | 1 |
| run-script-list-empty | MATCH | 0 | 0 |
| run-script-no-json | MATCH | 1 | 1 |
| run-script-disable-timeout | MATCH | 0 | 0 |
| run-script-working-dir | MATCH | 0 | 0 |
| exec-list | MATCH | 0 | 0 |
| exec-no-binary | MATCH | 0 | 0 |
| exec-run | MATCH | 0 | 0 |
| exec-run-args | MATCH | 0 | 0 |
| exec-run-php | MATCH | 0 | 0 |
| exec-missing | MATCH | 127 | 127 |
| exec-no-scripts | MATCH | 0 | 0 |
| exec-verbose | MATCH | 0 | 0 |
| exec-no-bins | MATCH | 1 | 1 |
| exec-local-list | MATCH | 0 | 0 |
| exec-no-json | MATCH | 1 | 1 |
| clear-cache-basic | MATCH | 0 | 0 |
| clear-cache-alias-cc | MATCH | 0 | 0 |
| clear-cache-alias-clearcache | MATCH | 0 | 0 |
| clear-cache-gc | MATCH | 0 | 0 |
| clear-cache-gc-maxsize | MATCH | 0 | 0 |
| clear-cache-gc-ttl | DIFF | 0 | 0 |
| clear-cache-gc-bad-maxsize | MATCH | 1 | 1 |
| clear-cache-read-only | MATCH | 0 | 0 |
| clear-cache-verbose | MATCH | 0 | 0 |
| clear-cache-empty | MATCH | 0 | 0 |
| clear-cache-empty-gc | MATCH | 0 | 0 |
| clear-cache-no-cache | MATCH | 0 | 0 |
| clear-cache-no-json | MATCH | 0 | 0 |
| clear-cache-custom-dirs | MATCH | 0 | 0 |
| fund-text | MATCH | 0 | 0 |
| fund-json | MATCH | 0 | 0 |
| fund-json-short | MATCH | 0 | 0 |
| fund-bad-format | MATCH | 1 | 1 |
| fund-v2 | MATCH | 0 | 0 |
| fund-none | MATCH | 0 | 0 |
| fund-none-json | MATCH | 0 | 0 |
| fund-not-installed | MATCH | 0 | 0 |
| fund-no-json | MATCH | 1 | 1 |
| search-fulltext | MATCH | 0 | 0 |
| search-multi-token | MATCH | 0 | 0 |
| search-keyword | MATCH | 0 | 0 |
| search-only-name | MATCH | 0 | 0 |
| search-only-name-long | MATCH | 0 | 0 |
| search-only-vendor | MATCH | 0 | 0 |
| search-only-vendor-long | MATCH | 0 | 0 |
| search-name-and-vendor | MATCH | 1 | 1 |
| search-type | MATCH | 0 | 0 |
| search-type-other | MATCH | 0 | 0 |
| search-type-vendor | MATCH | 0 | 0 |
| search-json | MATCH | 0 | 0 |
| search-json-vendor | MATCH | 0 | 0 |
| search-json-empty | MATCH | 0 | 0 |
| search-bad-format | MATCH | 1 | 1 |
| search-no-results | MATCH | 0 | 0 |
| search-regex-chars | MATCH | 0 | 0 |
| search-regex-chars-name | MATCH | 0 | 0 |
| search-no-tokens | MATCH | 1 | 1 |
| search-very-verbose | MATCH | 0 | 0 |
| search-columns | MATCH | 0 | 0 |
| search-installed | MATCH | 0 | 0 |
| search-installed-vv | MATCH | 0 | 0 |
| search-v2 | MATCH | 0 | 0 |
| search-v2-vendor | MATCH | 0 | 0 |
| search-v2-json | MATCH | 0 | 0 |
| search-wide | MATCH | 0 | 0 |
| search-filtered-repo | MATCH | 0 | 0 |
| browse-show | MATCH | 0 | 0 |
| browse-show-homepage | MATCH | 0 | 0 |
| browse-home-alias | MATCH | 0 | 0 |
| browse-homepage-missing | MATCH | 1 | 1 |
| browse-not-found | MATCH | 1 | 1 |
| browse-no-url | MATCH | 1 | 1 |
| browse-multi | MATCH | 1 | 1 |
| browse-open | MATCH | 0 | 0 |
| browse-root | MATCH | 1 | 1 |
| browse-root-homepage | MATCH | 0 | 0 |
| browse-root-source | MATCH | 0 | 0 |
| browse-root-open | MATCH | 0 | 0 |
| browse-root-invalid-url | MATCH | 1 | 1 |
| browse-installed | MATCH | 0 | 0 |
| browse-no-json | MATCH | 1 | 1 |
| diagnose-basic | MATCH | 1 | 1 |
| diagnose-clean | MATCH | 0 | 0 |
| diagnose-locked | MATCH | 0 | 0 |
| diagnose-lock-invalid | MATCH | 2 | 2 |
| diagnose-no-json | MATCH | 0 | 0 |
| diagnose-schema-error | MATCH | 0 | 0 |
| diagnose-platform-override | MATCH | 0 | 0 |
| diagnose-disable-tls | MATCH | 0 | 0 |
| diagnose-github-oauth | MATCH | 0 | 0 |
| diagnose-proxy | MATCH | 0 | 0 |
| diagnose-proxy-no-proxy | MATCH | 0 | 0 |
| diagnose-proxy-invalid | MATCH | 0 | 0 |
| diagnose-ipresolve | MATCH | 0 | 0 |
| diagnose-network | KNOWN-DIFF | 2 | 2 |
| self-update-basic | KNOWN-DIFF | 1 | 1 |
| self-update-alias | KNOWN-DIFF | 1 | 1 |
| self-update-rollback | KNOWN-DIFF | 1 | 1 |
| self-update-version | KNOWN-DIFF | 1 | 1 |
| self-update-no-json | KNOWN-DIFF | 1 | 1 |
| self-update-bad-option | MATCH | 1 | 1 |
| self-update-too-many | MATCH | 1 | 1 |
| repository-list | MATCH | 0 | 0 |
| repository-list-default | MATCH | 0 | 0 |
| repository-list-named | MATCH | 0 | 0 |
| repository-list-assoc | MATCH | 0 | 0 |
| repository-no-args | MATCH | 0 | 0 |
| repository-add | MATCH | 0 | 0 |
| repository-add-append | MATCH | 0 | 0 |
| repository-add-prepend-named | MATCH | 0 | 0 |
| repository-add-assoc | MATCH | 0 | 0 |
| repository-add-json | MATCH | 0 | 0 |
| repository-add-json-invalid | MATCH | 1 | 1 |
| repository-add-missing-name | MATCH | 1 | 1 |
| repository-add-missing-arg1 | MATCH | 1 | 1 |
| repository-add-missing-url | MATCH | 1 | 1 |
| repository-add-before | MATCH | 0 | 0 |
| repository-add-after | MATCH | 0 | 0 |
| repository-add-before-after | MATCH | 1 | 1 |
| repository-add-before-missing | MATCH | 1 | 1 |
| repository-add-existing | MATCH | 0 | 0 |
| repository-remove | MATCH | 0 | 0 |
| repository-remove-assoc | MATCH | 0 | 0 |
| repository-remove-missing | MATCH | 0 | 0 |
| repository-remove-packagist | MATCH | 0 | 0 |
| repository-remove-no-name | MATCH | 1 | 1 |
| repository-set-url | MATCH | 0 | 0 |
| repository-set-url-assoc | MATCH | 0 | 0 |
| repository-set-url-missing | MATCH | 1 | 1 |
| repository-set-url-unknown | MATCH | 0 | 0 |
| repository-get-url | MATCH | 0 | 0 |
| repository-get-url-assoc | MATCH | 0 | 0 |
| repository-get-url-no-url | MATCH | 1 | 1 |
| repository-get-url-missing | MATCH | 1 | 1 |
| repository-get-url-no-name | MATCH | 1 | 1 |
| repository-disable | MATCH | 0 | 0 |
| repository-disable-short | MATCH | 0 | 0 |
| repository-disable-other | MATCH | 1 | 1 |
| repository-disable-no-name | MATCH | 1 | 1 |
| repository-enable | MATCH | 0 | 0 |
| repository-enable-other | MATCH | 1 | 1 |
| repository-unknown-action | MATCH | 1 | 1 |
| repository-uppercase-action | MATCH | 0 | 0 |
| repository-global-add | MATCH | 0 | 0 |
| repository-global-list | MATCH | 0 | 0 |
| repository-file | MATCH | 1 | 1 |
| repository-file-and-global | MATCH | 1 | 1 |
| repository-file-missing | MATCH | 1 | 1 |
| repository-file-add | MATCH | 0 | 0 |
| repository-no-json | MATCH | 1 | 1 |
| policy-cmd-add-source | MATCH | 0 | 0 |
| policy-cmd-add-source-json | MATCH | 0 | 0 |
| policy-cmd-add-source-existing-list | MATCH | 0 | 0 |
| policy-cmd-add-source-duplicate | MATCH | 0 | 0 |
| policy-cmd-add-source-object-sources | MATCH | 1 | 1 |
| policy-cmd-add-source-uppercase-action | MATCH | 0 | 0 |
| policy-cmd-no-action | MATCH | 1 | 1 |
| policy-cmd-unknown-action | MATCH | 1 | 1 |
| policy-cmd-no-name | MATCH | 1 | 1 |
| policy-cmd-builtin | MATCH | 1 | 1 |
| policy-cmd-builtin-cooldown | MATCH | 1 | 1 |
| policy-cmd-reserved-prefix | MATCH | 1 | 1 |
| policy-cmd-reserved-name | MATCH | 1 | 1 |
| policy-cmd-dotted-name | MATCH | 1 | 1 |
| policy-cmd-empty-name | MATCH | 1 | 1 |
| policy-cmd-no-arg1 | MATCH | 1 | 1 |
| policy-cmd-no-url | MATCH | 1 | 1 |
| policy-cmd-bad-type | MATCH | 1 | 1 |
| policy-cmd-http-url | MATCH | 1 | 1 |
| policy-cmd-json-no-type | MATCH | 1 | 1 |
| policy-cmd-json-no-url | MATCH | 1 | 1 |
| policy-cmd-json-int-url | MATCH | 1 | 1 |
| policy-cmd-json-invalid | MATCH | 1 | 1 |
| policy-cmd-global | MATCH | 0 | 0 |
| policy-cmd-file-and-global | MATCH | 1 | 1 |
| policy-cmd-file | MATCH | 0 | 0 |
| policy-cmd-no-json | MATCH | 1 | 1 |

## clear-cache-gc-ttl

```
project tree:
- (oracle) D c/vcs/foo.git
- (oracle) F c/files/acme/log/abc.zip 644 eb076a2ec6ced9ee
- (oracle) F c/repo/http---x/packages.json 644 44136fa355b3678a
- (oracle) F c/vcs/foo.git/HEAD 644 3ff6c05723bb069d
```

## diagnose-network

```
stdout:
+ (go) Checking Composer and its dependencies for vulnerabilities: WARNING
+ (go) Could not find Composer's installed.json, this must be a non-standard Composer installation.
- (oracle) Checking Composer and its dependencies for vulnerabilities: Failed performing audit: curl error 6 while downloading https://repo.packagist.org/packages.json: Could not resolve host: repo.packagist.org
```

## self-update-basic

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-alias

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-rollback

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-version

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-no-json

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```
