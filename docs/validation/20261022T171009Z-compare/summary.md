# Composer vs gocomposer comparison

- gocomposer_sha256: `8785a96420b40e3cc583f95d565a3781250374d2aa63094181f1d52ac6fbc4a6`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| remove-basic | MATCH | 0 | 0 |
| remove-dev | MATCH | 0 | 0 |
| remove-wrong-section | MATCH | 2 | 2 |
| remove-not-required | MATCH | 0 | 0 |
| remove-wildcard | MATCH | 0 | 0 |
| remove-no-update | MATCH | 0 | 0 |
| remove-dry-run | MATCH | 0 | 0 |
| remove-still-required | MATCH | 2 | 2 |
| remove-unused | MATCH | 0 | 0 |
| remove-no-args | MATCH | 1 | 1 |
| remove-mixed-case | MATCH | 0 | 0 |
| remove-allow-plugins | MATCH | 0 | 0 |
