# Composer vs gocomposer comparison

- gocomposer_sha256: `f927922781922a813314e2e89f78606f6e1e7804eb27bcdd33dd8a6f810aa051`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| config-list-custom | MATCH | 0 | 0 |
| config-list-named-repos | MATCH | 0 | 0 |
| config-set-autoloader-suffix-null | MATCH | 1 | 1 |
| config-repo-remove | MATCH | 0 | 0 |
| config-repo-remove-assoc | MATCH | 0 | 0 |
| config-repo-get-named | MATCH | 0 | 0 |
| config-global-list-existing | MATCH | 0 | 0 |
