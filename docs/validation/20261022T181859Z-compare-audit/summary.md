# Composer vs gocomposer comparison

- gocomposer_sha256: `090ed25d2162010762a0f516e84d47a65dd1ff7d35452399995bbaf3cfaff838`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| audit-clean | MATCH | 0 | 0 |
| audit-clean-json | MATCH | 0 | 0 |
| audit-clean-summary | MATCH | 0 | 0 |
| audit-no-advisory-repo | MATCH | 0 | 0 |
| audit-table | MATCH | 1 | 1 |
| audit-plain | MATCH | 1 | 1 |
| audit-json | MATCH | 1 | 1 |
| audit-summary | MATCH | 1 | 1 |
| audit-format-short | MATCH | 1 | 1 |
| audit-format-invalid | MATCH | 1 | 1 |
| audit-format-missing-value | MATCH | 1 | 1 |
| audit-verbose | MATCH | 1 | 1 |
| audit-no-vendor | MATCH | 1 | 1 |
| audit-no-vendor-no-dev-only-dev | MATCH | 0 | 0 |
| audit-no-packages | MATCH | 0 | 0 |
| audit-no-lock-no-vendor | MATCH | 1 | 1 |
| audit-locked | MATCH | 1 | 1 |
| audit-locked-json | MATCH | 1 | 1 |
| audit-locked-no-lock | MATCH | 1 | 1 |
| audit-no-dev | MATCH | 0 | 0 |
| audit-with-dev | MATCH | 1 | 1 |
| audit-locked-no-dev | MATCH | 0 | 0 |
| audit-direct | MATCH | 1 | 1 |
| audit-direct-short | MATCH | 1 | 1 |
| audit-abandoned-default | MATCH | 1 | 1 |
| audit-abandoned-plain | MATCH | 1 | 1 |
| audit-abandoned-json | MATCH | 1 | 1 |
| audit-abandoned-summary | MATCH | 1 | 1 |
| audit-abandoned-ignore | MATCH | 0 | 0 |
| audit-abandoned-report | MATCH | 0 | 0 |
| audit-abandoned-fail | MATCH | 1 | 1 |
| audit-abandoned-invalid | MATCH | 1 | 1 |
| audit-abandoned-config-report | MATCH | 0 | 0 |
| audit-abandoned-config-ignore | MATCH | 0 | 0 |
| audit-abandoned-option-overrides-config | MATCH | 1 | 1 |
| audit-abandoned-env-report | MATCH | 0 | 0 |
| audit-abandoned-env-invalid | ERROR | 0 | 0 |
| audit-abandoned-and-vulnerable | MATCH | 1 | 1 |
| audit-ignore-id | MATCH | 1 | 1 |
| audit-ignore-id-reason | MATCH | 1 | 1 |
| audit-ignore-id-json | MATCH | 1 | 1 |
| audit-ignore-cve | MATCH | 1 | 1 |
| audit-ignore-remote-id | MATCH | 1 | 1 |
| audit-ignore-package | MATCH | 1 | 1 |
| audit-ignore-all | MATCH | 0 | 0 |
| audit-ignore-summary-cve | MATCH | 1 | 1 |
| audit-ignore-apply-block-only | MATCH | 1 | 1 |
| audit-ignore-apply-audit | MATCH | 1 | 1 |
| audit-ignore-severity-option | MATCH | 1 | 1 |
| audit-ignore-severity-option-multi | MATCH | 1 | 1 |
| audit-ignore-severity-option-other | MATCH | 1 | 1 |
| audit-ignore-severity-config | MATCH | 1 | 1 |
| audit-ignore-severity-config-reason | MATCH | 1 | 1 |
| audit-policy-advisories-ignore | MATCH | 1 | 1 |
| audit-policy-disabled | MATCH | 1 | 1 |
| audit-api-table | MATCH | 1 | 1 |
| audit-api-json | MATCH | 1 | 1 |
| audit-api-summary | MATCH | 1 | 1 |
| audit-api-locked | MATCH | 1 | 1 |
| audit-api-unrequested-names | MATCH | 0 | 0 |
| audit-api-ignore-severity | MATCH | 1 | 1 |
| audit-unreachable-404 | MATCH | 100 | 100 |
| audit-unreachable-404-ignored | MATCH | 1 | 1 |
| audit-unreachable-404-ignored-json | MATCH | 1 | 1 |
| audit-unreachable-refused-ignored | DIFF | 0 | 0 |
| audit-unreachable-config-ignored | MATCH | 0 | 0 |
| audit-unreachable-api-ignored | MATCH | 0 | 0 |

## audit-abandoned-env-invalid

```
setup failed: update --no-progress --no-blocking

In Config.php line 532:
                                                                               
  Invalid value for COMPOSER_AUDIT_ABANDONED: bogus. Expected one of ignore,   
  report, fail.                                                                
                                                                               

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## audit-unreachable-refused-ignored

```
stderr:
+ (go)   - curl error 7 while downloading http://127.0.0.1:1/packages.json: Failed to connect to 127.0.0.1 port 1 after 0 ms: Couldn't connect to server
+ (go)   - curl error 7 while downloading http://127.0.0.1:1/packages.json: Failed to connect to 127.0.0.1 port 1 after 0 ms: Couldn't connect to server
- (oracle)   - curl error 7 while downloading http://127.0.0.1:1/packages.json: Failed to connect to 127.0.0.1 port 1 after 0 ms: Could not connect to server
- (oracle)   - curl error 7 while downloading http://127.0.0.1:1/packages.json: Failed to connect to 127.0.0.1 port 1 after 0 ms: Could not connect to server
```
