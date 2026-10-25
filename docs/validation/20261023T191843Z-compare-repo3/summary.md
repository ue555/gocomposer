# Composer vs gocomposer comparison

- gocomposer_sha256: `074dedd5a353d075f1905892096a828df8e4f7a82eb31fdfbcc0fb20680e4ea6`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:cd7f1f040bba14a1974d175df2110347c2be1d880834ad6a106a3578bfd929fa`
- php_version: `8.4.27`
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
