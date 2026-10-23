# Composer vs gocomposer comparison

- gocomposer_sha256: `75952f361d1af71f56a4760d9d264f67fe8c9e2cb8e221153cf6c95e48327830`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive-root-tar | MATCH | 0 | 0 |
| archive-root-zip | MATCH | 0 | 0 |
| archive-root-tar-gz | MATCH | 0 | 0 |
| archive-root-tar-bz2 | MATCH | 1 | 1 |
| archive-root-bad-format | MATCH | 1 | 1 |
| archive-root-empty-format | MATCH | 1 | 1 |
| archive-root-dir | MATCH | 0 | 0 |
| archive-root-dir-relative | MATCH | 0 | 0 |
| archive-root-file | MATCH | 0 | 0 |
| archive-root-file-subdir | MATCH | 0 | 0 |
| archive-root-ignore-filters | MATCH | 0 | 0 |
| archive-root-named | MATCH | 0 | 0 |
| archive-root-installed | MATCH | 0 | 0 |
| archive-root-existing-archives | MATCH | 0 | 0 |
| archive-root-config-format | MATCH | 0 | 0 |
| archive-root-dir-is-file | MATCH | 1 | 1 |
| archive-root-scripts | MATCH | 0 | 0 |
| archive-package | MATCH | 0 | 0 |
| archive-package-installed | MATCH | 0 | 0 |
| archive-package-constraint | MATCH | 0 | 0 |
| archive-package-stability | MATCH | 0 | 0 |
| archive-package-stability-filtered | MATCH | 1 | 1 |
| archive-package-missing | MATCH | 1 | 1 |
| archive-package-file | MATCH | 1 | 1 |
| archive-package-case | MATCH | 0 | 0 |
| archive-no-json | MATCH | 1 | 1 |
| archive-empty-tar | MATCH | 0 | 0 |
| archive-empty-zip | MATCH | 0 | 0 |
| archive-empty-tar-gz | MATCH | 0 | 0 |
| archive-only-empty-dirs | MATCH | 0 | 0 |
| archive-globstar-attributes | MATCH | 1 | 1 |
| archive-globstar-valid | MATCH | 0 | 0 |
| archive-vendor-dir | MATCH | 0 | 0 |
