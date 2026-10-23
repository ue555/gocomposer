# Composer vs gocomposer comparison

- gocomposer_sha256: `0e9a127d41ed42c0d5034ea3600b2eae0bd3739b11c01200b3a9cf833f9bca9a`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| create-project-dot | MATCH | 0 | 0 |
| create-project-audit-format | MATCH | 0 | 0 |
| create-project-prefer-install-auto | MATCH | 0 | 0 |
| create-project-prefer-install-bad | MATCH | 1 | 1 |
| create-project-stability-rc-upper | MATCH | 0 | 0 |
| create-project-composer-env | MATCH | 0 | 0 |
