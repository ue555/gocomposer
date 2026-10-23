# Composer vs gocomposer comparison

- gocomposer_sha256: `cce199076e97cf1cb951f4c66df7971ce683ad8b9cc3751f26f4f2c3c7880eb9`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| fund-text | MATCH | 0 | 0 |
| fund-json | MATCH | 0 | 0 |
| fund-json-short | MATCH | 0 | 0 |
| fund-bad-format | MATCH | 1 | 1 |
| fund-v2 | MATCH | 0 | 0 |
| fund-none | MATCH | 0 | 0 |
| fund-none-json | MATCH | 0 | 0 |
| fund-not-installed | MATCH | 0 | 0 |
| fund-no-json | MATCH | 1 | 1 |
| search-fulltext | MATCH | 0 | 0 |
| search-multi-token | MATCH | 0 | 0 |
| search-keyword | MATCH | 0 | 0 |
| search-only-name | MATCH | 0 | 0 |
| search-only-name-long | MATCH | 0 | 0 |
| search-only-vendor | MATCH | 0 | 0 |
| search-only-vendor-long | MATCH | 0 | 0 |
| search-name-and-vendor | MATCH | 1 | 1 |
| search-type | MATCH | 0 | 0 |
| search-type-other | MATCH | 0 | 0 |
| search-type-vendor | MATCH | 0 | 0 |
| search-json | MATCH | 0 | 0 |
| search-json-vendor | MATCH | 0 | 0 |
| search-json-empty | MATCH | 0 | 0 |
| search-bad-format | MATCH | 1 | 1 |
| search-no-results | MATCH | 0 | 0 |
| search-regex-chars | MATCH | 0 | 0 |
| search-regex-chars-name | MATCH | 0 | 0 |
| search-no-tokens | MATCH | 1 | 1 |
| search-very-verbose | MATCH | 0 | 0 |
| search-columns | MATCH | 0 | 0 |
| search-installed | MATCH | 0 | 0 |
| search-installed-vv | MATCH | 0 | 0 |
| search-v2 | MATCH | 0 | 0 |
| search-v2-vendor | MATCH | 0 | 0 |
| search-v2-json | MATCH | 0 | 0 |
| search-wide | MATCH | 0 | 0 |
| search-filtered-repo | MATCH | 0 | 0 |
| browse-show | MATCH | 0 | 0 |
| browse-show-homepage | MATCH | 0 | 0 |
| browse-home-alias | MATCH | 0 | 0 |
| browse-homepage-missing | MATCH | 1 | 1 |
| browse-not-found | MATCH | 1 | 1 |
| browse-no-url | MATCH | 1 | 1 |
| browse-multi | MATCH | 1 | 1 |
| browse-open | MATCH | 0 | 0 |
| browse-root | MATCH | 1 | 1 |
| browse-root-homepage | MATCH | 0 | 0 |
| browse-root-source | MATCH | 0 | 0 |
| browse-root-open | MATCH | 0 | 0 |
| browse-root-invalid-url | MATCH | 1 | 1 |
| browse-installed | MATCH | 0 | 0 |
| browse-no-json | MATCH | 1 | 1 |
