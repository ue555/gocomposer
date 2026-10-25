# Composer vs gocomposer comparison

- gocomposer_sha256: `3d7493549c16f61be1d9db36c5429cec7b2a69149cb31ef2e94d002f8371f46c`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.5`
- oracle_image_digest: `sha256:8f997f444c893bd29a5143015dfa981662afe3e936c43d87adf8416d912f92e2`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.5-zip`
- oracle_zip_image_digest: `sha256:c0ec954e3ac6a8eeb8f9a5fb7ec90943a1414e08c86351127f1c15969e234528`
- php_version: `8.5.12`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| providers-update | MATCH | 0 | 0 |
| providers-update-verbose | MATCH | 0 | 0 |
| providers-install-lock | MATCH | 0 | 0 |
| providers-update-cached | MATCH | 0 | 0 |
| providers-root-listing | MATCH | 0 | 0 |
| providers-dev-alias | MATCH | 0 | 0 |
| providers-virtual | MATCH | 2 | 2 |
| providers-not-found | MATCH | 2 | 2 |
| providers-bad-hash | MATCH | 1 | 1 |
| providers-bad-hash-other | MATCH | 0 | 0 |
| providers-require | MATCH | 0 | 0 |
| providers-require-missing | MATCH | 1 | 1 |
| providers-show-available | MATCH | 0 | 0 |
| providers-show-all | MATCH | 0 | 0 |
| providers-show-all-names | MATCH | 0 | 0 |
| providers-outdated | MATCH | 0 | 0 |
| providers-update-one | MATCH | 0 | 0 |
| providers-depends | MATCH | 0 | 0 |
| providers-search | MATCH | 0 | 0 |
| providers-search-name | MATCH | 0 | 0 |
| providers-search-vendor | MATCH | 0 | 0 |
| providers-search-list | MATCH | 0 | 0 |
| providers-search-list-other | MATCH | 0 | 0 |
| providers-lazy-update | MATCH | 0 | 0 |
| providers-lazy-not-found | MATCH | 2 | 2 |
| providers-lazy-show | MATCH | 0 | 0 |
| providers-lazy-search | MATCH | 0 | 0 |
| providers-lazy-cached | MATCH | 0 | 0 |
| archive3-xz | MATCH | 0 | 0 |
| archive3-xz-status | MATCH | 1 | 1 |
| archive3-xz-reinstall | MATCH | 0 | 0 |
| archive3-xz-remove | MATCH | 0 | 0 |
| archive3-rar-status | MATCH | 1 | 1 |
| archive3-xz-broken | MATCH | 1 | 1 |
| archive3-rar-no-unrar | MATCH | 1 | 1 |
| archive3-rar | MATCH | 0 | 0 |
| archive3-rar-broken | MATCH | 1 | 1 |
| archive3-zip-7z | MATCH | 0 | 0 |
| archive3-zip-7z-zipext | MATCH | 0 | 0 |
| archive3-zip-7z-crc | MATCH | 1 | 1 |
| archive3-zip-7z-crc-zipext | MATCH | 1 | 1 |
| archive3-zip-no-command | MATCH | 1 | 1 |
| archive3-zip-no-command-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-crc-zipext | MATCH | 1 | 1 |
| warnings-update | MATCH | 0 | 0 |
| warnings-update-ansi | MATCH | 0 | 0 |
| warnings-show | MATCH | 0 | 0 |
| warnings-search-api | MATCH | 0 | 0 |
| warnings-search-list | MATCH | 0 | 0 |
| warnings-show-names-list | MATCH | 0 | 0 |
| warnings-audit-api | MATCH | 1 | 1 |
| warnings-update-audit-api | MATCH | 0 | 0 |
| warnings-filter-audit | MATCH | 0 | 0 |
| warnings-providers | MATCH | 0 | 0 |
| warnings-two-repos | MATCH | 0 | 0 |
| warnings-http-403 | MATCH | 100 | 100 |
| warnings-http-404-metadata | MATCH | 2 | 2 |
