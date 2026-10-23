# Composer vs gocomposer comparison

- gocomposer_sha256: `af3e748ed4a4d1038eb8556789ef799708a4e764e4dbe9864078e59c06efcc52`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| repository-list | MATCH | 0 | 0 |
| repository-list-default | MATCH | 0 | 0 |
| repository-list-named | MATCH | 0 | 0 |
| repository-list-assoc | MATCH | 0 | 0 |
| repository-no-args | MATCH | 0 | 0 |
| repository-add | MATCH | 0 | 0 |
| repository-add-append | MATCH | 0 | 0 |
| repository-add-prepend-named | MATCH | 0 | 0 |
| repository-add-assoc | MATCH | 0 | 0 |
| repository-add-json | MATCH | 0 | 0 |
| repository-add-json-invalid | MATCH | 1 | 1 |
| repository-add-missing-name | MATCH | 1 | 1 |
| repository-add-missing-arg1 | MATCH | 1 | 1 |
| repository-add-missing-url | MATCH | 1 | 1 |
| repository-add-before | MATCH | 0 | 0 |
| repository-add-after | MATCH | 0 | 0 |
| repository-add-before-after | MATCH | 1 | 1 |
| repository-add-before-missing | MATCH | 1 | 1 |
| repository-add-existing | MATCH | 0 | 0 |
| repository-remove | MATCH | 0 | 0 |
| repository-remove-assoc | MATCH | 0 | 0 |
| repository-remove-missing | MATCH | 0 | 0 |
| repository-remove-packagist | MATCH | 0 | 0 |
| repository-remove-no-name | MATCH | 1 | 1 |
| repository-set-url | MATCH | 0 | 0 |
| repository-set-url-assoc | MATCH | 0 | 0 |
| repository-set-url-missing | MATCH | 1 | 1 |
| repository-set-url-unknown | MATCH | 0 | 0 |
| repository-get-url | MATCH | 0 | 0 |
| repository-get-url-assoc | MATCH | 0 | 0 |
| repository-get-url-no-url | MATCH | 1 | 1 |
| repository-get-url-missing | MATCH | 1 | 1 |
| repository-get-url-no-name | MATCH | 1 | 1 |
| repository-disable | MATCH | 0 | 0 |
| repository-disable-short | MATCH | 0 | 0 |
| repository-disable-other | MATCH | 1 | 1 |
| repository-disable-no-name | MATCH | 1 | 1 |
| repository-enable | MATCH | 0 | 0 |
| repository-enable-other | MATCH | 1 | 1 |
| repository-unknown-action | MATCH | 1 | 1 |
| repository-uppercase-action | MATCH | 0 | 0 |
| repository-global-add | MATCH | 0 | 0 |
| repository-global-list | MATCH | 0 | 0 |
| repository-file | MATCH | 1 | 1 |
| repository-file-and-global | MATCH | 1 | 1 |
| repository-file-missing | MATCH | 1 | 1 |
| repository-file-add | MATCH | 0 | 0 |
| repository-no-json | MATCH | 1 | 1 |
