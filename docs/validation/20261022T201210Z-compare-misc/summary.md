# Composer vs gocomposer comparison

- gocomposer_sha256: `64b2ac1cf594354a7e992450c65d456d26196bc32183ff64da07362f14e46360`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| clear-cache-basic | MATCH | 0 | 0 |
| clear-cache-alias-cc | MATCH | 0 | 0 |
| clear-cache-alias-clearcache | MATCH | 0 | 0 |
| clear-cache-gc | MATCH | 0 | 0 |
| clear-cache-gc-maxsize | MATCH | 0 | 0 |
| clear-cache-gc-ttl-fresh | MATCH | 0 | 0 |
| clear-cache-gc-bad-maxsize | MATCH | 1 | 1 |
| clear-cache-read-only | MATCH | 0 | 0 |
| clear-cache-verbose | MATCH | 0 | 0 |
| clear-cache-empty | MATCH | 0 | 0 |
| clear-cache-empty-gc | MATCH | 0 | 0 |
| clear-cache-no-cache | MATCH | 0 | 0 |
| clear-cache-no-json | MATCH | 0 | 0 |
| clear-cache-custom-dirs | MATCH | 0 | 0 |
