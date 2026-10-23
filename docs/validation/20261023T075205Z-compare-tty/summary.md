# Composer vs gocomposer comparison

- gocomposer_sha256: `e88d0cd8f2bdd4c11d21477985e067a0dc17c0e28da985334c416ba14e0cb5f5`
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
| tty-remove-alt-yes | MATCH | 0 | 0 |
| tty-remove-alt-default | MATCH | 0 | 0 |
| tty-remove-alt-no | MATCH | 0 | 0 |
| tty-remove-alt-invalid | MATCH | 0 | 0 |
| tty-remove-alt-eof | MATCH | 0 | 0 |
| tty-remove-alt-dev-wildcard | MATCH | 0 | 0 |
| tty-remove-alt-dry-run | MATCH | 0 | 0 |
| tty-remove-alt-no-interaction | MATCH | 0 | 0 |
| tty-remove-alt-no-ansi | MATCH | 0 | 0 |
