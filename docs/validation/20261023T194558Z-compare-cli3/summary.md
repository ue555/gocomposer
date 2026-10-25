# Composer vs gocomposer comparison

- gocomposer_sha256: `9cacc8b370abb33fa4afc72653083c0a581c887a9ca7f48522c925aea01fe159`
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
| cli3-vvv-auth-files | MATCH | 0 | 0 |
| cli3-vvv-global-home | MATCH | 0 | 0 |
| cli3-event-validate | MATCH | 0 | 0 |
| cli3-cmdclass-direct-named | MATCH | 0 | 0 |
| cli3-packagist-search-all | MATCH | 0 | 0 |
| cli3-packagist-fund | MATCH | 0 | 0 |
| cli3-packagist-fund-json | MATCH | 0 | 0 |
