# Composer vs gocomposer comparison

- gocomposer_sha256: `5899e972c99a3b25cca4025aa5e8c8f20fab6706f969a0217b717bbd64d8104a`
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
| term-href-audit | MATCH | 1 | 1 |
| term-href-audit-plain | MATCH | 1 | 1 |
| term-progress-install | MATCH | 0 | 0 |
| term-progress-update | MATCH | 0 | 0 |
| term-progress-require | MATCH | 0 | 0 |
| term-progress-verbose | MATCH | 0 | 0 |
| term-progress-very-verbose | MATCH | 0 | 0 |
| term-progress-narrow | MATCH | 0 | 0 |
| term-progress-tiny | MATCH | 0 | 0 |
| term-progress-no-progress | MATCH | 0 | 0 |
| term-progress-no-ansi | MATCH | 0 | 0 |
| term-progress-ci | MATCH | 0 | 0 |
| term-progress-ansi-pipe | MATCH | 0 | 0 |
