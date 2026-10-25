# Composer vs gocomposer comparison

- gocomposer_sha256: `4c25dab6af44adad23259a499209495b3af1c890fb5af8491e5737b8b54453b9`
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
| archive3-xz-broken | MATCH | 1 | 1 |
| archive3-zip-7z-crc | MATCH | 1 | 1 |
| archive3-zip-7z-crc-zipext | MATCH | 1 | 1 |
| archive3-zip-no-command-crc-zipext | MATCH | 1 | 1 |
| warnings-update-ansi | MATCH | 0 | 0 |
| warnings-audit-api | MATCH | 1 | 1 |
