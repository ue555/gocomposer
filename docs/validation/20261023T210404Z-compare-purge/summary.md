# Composer vs gocomposer comparison

- gocomposer_sha256: `1354dc29a272668426fe46931ea598838d0c6623d375d8a0c2d0872ca32bc98c`
- host_os: `darwin/arm64`
- oracle_arch_image: `localhost/gocomposer-oracle:8.4-arch`
- oracle_arch_image_digest: `sha256:93bc04d027eaeb74edeb8ea07a247c5f17cb2e91d4d17fe97681f949f96116e7`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:9cf420cd931841442e55100b424e04a60af91810d670826961c6583bd372308b`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:e3a99ff1ce724bd01dbd722f202bf4afa98eec9bde0789c99ad6f76b266e6f22`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| purge-install | MATCH | 0 | 0 |
| purge-install-all | DIFF | 0 | 0 |
| purge-show | MATCH | 0 | 0 |
| purge-show-name | MATCH | 1 | 1 |
| purge-update | MATCH | 0 | 0 |
| purge-dump-autoload | MATCH | 0 | 0 |
| purge-depends | MATCH | 0 | 0 |
| purge-licenses | MATCH | 0 | 0 |
| purge-status | MATCH | 0 | 0 |
| purge-remove | MATCH | 0 | 0 |

## purge-install-all

```
stderr:
+ (go)     Skipped installation of bin bin/acme for package acme/cli: name conflicts with an existing file
+ (go)     Skipped installation of bin bin/acme-sh for package acme/cli: name conflicts with an existing file
- (oracle)     Skipped installation of bin bin/acme for package acme/cli: name conflicts with an existing file
- (oracle)     Skipped installation of bin bin/acme-sh for package acme/cli: name conflicts with an existing file
```
