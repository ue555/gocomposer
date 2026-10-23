# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `c72c8df888a9f2f675084b95798a893517f882cd89132ed778b39e05b7c23eb4`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| project | step | status | oracle exit | go exit |
|---|---|---|---|---|
| library-symfony-console | install | MATCH | 0 | 0 |
| library-symfony-console | update-no-install | MATCH | 0 | 0 |
| cli-app | install | MATCH | 0 | 0 |
| cli-app | update-no-install | MATCH | 255 | 255 |
| laravel | install | MATCH | 0 | 0 |
| laravel | update-no-install | MATCH | 255 | 255 |
| live-psr | install | MATCH | 0 | 0 |
| live-psr | update-no-install | MATCH | 0 | 0 |
