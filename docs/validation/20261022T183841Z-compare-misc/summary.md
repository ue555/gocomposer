# Composer vs gocomposer comparison

- gocomposer_sha256: `3cf48393d65872e698823c59fd71e349e0e6e24bf4ce1cd14733b2cafa7bddef`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| diagnose-network | KNOWN-DIFF | 2 | 2 |

## diagnose-network

```
stdout:
+ (go) Checking Composer and its dependencies for vulnerabilities: WARNING
+ (go) Could not find Composer's installed.json, this must be a non-standard Composer installation.
- (oracle) Checking Composer and its dependencies for vulnerabilities: Failed performing audit: curl error 6 while downloading https://repo.packagist.org/packages.json: Could not resolve host: repo.packagist.org
```
