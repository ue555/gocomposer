# Composer vs gocomposer comparison

- gocomposer_sha256: `fa37a7f5ad570c358063405cd5d827ccba41de81fd094b84229ad4baf2292618`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| diagnose-basic | MATCH | 1 | 1 |
| diagnose-clean | MATCH | 0 | 0 |
| diagnose-locked | ERROR | 0 | 0 |
| diagnose-no-json | MATCH | 0 | 0 |
| diagnose-schema-error | MATCH | 0 | 0 |
| diagnose-platform-override | MATCH | 0 | 0 |
| diagnose-disable-tls | DIFF | 0 | 0 |
| diagnose-github-oauth | MATCH | 0 | 0 |
| diagnose-proxy | MATCH | 0 | 0 |
| diagnose-proxy-no-proxy | MATCH | 0 | 0 |
| diagnose-proxy-invalid | MATCH | 0 | 0 |
| diagnose-ipresolve | MATCH | 0 | 0 |
| diagnose-network | KNOWN-DIFF | 2 | 2 |

## diagnose-locked

```
setup failed: update --no-progress
Loading composer repositories with package information

In HttpDownloader.php line 330:
                                                                           
  Network disabled, request canceled: http://127.0.0.1:8899/packages.json  
                                                                           

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## diagnose-disable-tls

```
stderr:
- (oracle) You are running Composer with SSL/TLS protection disabled.
```

## diagnose-network

```
stdout:
+ (go) Checking Composer and its dependencies for vulnerabilities: WARNING
+ (go) Could not find Composer's installed.json, this must be a non-standard Composer installation.
- (oracle) Checking Composer and its dependencies for vulnerabilities: Failed performing audit: curl error 6 while downloading https://repo.packagist.org/packages.json: Could not resolve host: repo.packagist.org
- (oracle) The following exception probably indicates you are offline or have misconfigured DNS resolver(s)
- (oracle) The following exception probably indicates you are offline or have misconfigured DNS resolver(s)
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```
