# Composer vs gocomposer comparison

- gocomposer_sha256: `f7e16b44e3f6e59eb00622d2e9437ac8e538c5bcdb81720b41b40fd1e1ecd868`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| update-basic | MATCH | 0 | 0 |
| install-from-lock | MATCH | 0 | 0 |
| install-no-dev | MATCH | 0 | 0 |
| install-twice | MATCH | 0 | 0 |
| dry-run-update | MATCH | 0 | 0 |
| partial-update | MATCH | 0 | 0 |
| failure-keeps-vendor | MATCH | 2 | 2 |
