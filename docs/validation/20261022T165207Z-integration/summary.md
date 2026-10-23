# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `8df50073ab060c6fab1c12450c3d32831f685bc24819e3daf469d8c94eb03f6d`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| project | step | status | oracle exit | go exit |
|---|---|---|---|---|
| laravel | install | MATCH | 0 | 0 |
| laravel | update-no-install | MATCH | 255 | 255 |
