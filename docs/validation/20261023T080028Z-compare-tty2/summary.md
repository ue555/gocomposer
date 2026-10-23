# Composer vs gocomposer comparison

- gocomposer_sha256: `9df67b29bd97cd637def8a5e786ed33ad16e79d00d1cbfc5ed1fd317d87725fb`
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
| tty-discard-update-yes | MATCH | 0 | 0 |
| tty-discard-update-no | MATCH | 1 | 1 |
| tty-discard-update-help | MATCH | 0 | 0 |
| tty-discard-update-view | MATCH | 0 | 0 |
| tty-discard-update-stash | MATCH | 1 | 1 |
| tty-discard-update-verbose | MATCH | 0 | 0 |
| tty-discard-update-no-interaction | MATCH | 1 | 1 |
| tty-discard-remove-yes | MATCH | 0 | 0 |
| tty-discard-remove-stash-is-help | MATCH | 1 | 1 |
| tty-discard-reinstall | MATCH | 0 | 0 |
| tty-discard-reinstall-no | MATCH | 1 | 1 |
| tty-discard-untracked-only | MATCH | 0 | 0 |
