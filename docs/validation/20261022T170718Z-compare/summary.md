# Composer vs gocomposer comparison

- gocomposer_sha256: `55f1b78d85ee4e9847519cc8f47f6797b942f6dcd453b720fcfc4918e641a311`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| dry-run-update | MATCH | 0 | 0 |
| require-guess-version | MATCH | 0 | 0 |
| require-explicit | MATCH | 0 | 0 |
| require-dev | MATCH | 0 | 0 |
| require-no-update | MATCH | 0 | 0 |
| require-dry-run | MATCH | 0 | 0 |
| require-first-no-lock | MATCH | 0 | 0 |
| require-no-composer-json | MATCH | 1 | 1 |
| require-not-found | MATCH | 1 | 1 |
| require-min-stability | MATCH | 1 | 1 |
| require-platform-incompatible | MATCH | 1 | 1 |
| require-fails-reverts | MATCH | 2 | 2 |
| require-move-from-dev | MATCH | 0 | 0 |
| require-sort-packages | MATCH | 0 | 0 |
| require-too-strict | MATCH | 0 | 0 |
| require-self | MATCH | 1 | 1 |
| require-ext-missing | MATCH | 1 | 1 |
| require-alias-arg | MATCH | 1 | 1 |
| require-fixed-library | MATCH | 1 | 1 |
| require-fixed-project | MATCH | 0 | 0 |
| require-update-with-deps | MATCH | 0 | 0 |
| require-crlf-file | MATCH | 0 | 0 |
| require-tabs-file | MATCH | 0 | 0 |
