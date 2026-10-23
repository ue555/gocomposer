# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `9057ffdf50c8fd556a4423c3f6157d524f323ad9cb7659cac1a45f2fd250e5a4`
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
| plugin-installers | install | MATCH | 0 | 0 |
| plugin-installers | update-no-install | MATCH | 0 | 0 |
| plugin-php-http-discovery | install | MATCH | 0 | 0 |
| plugin-php-http-discovery | update-no-install | MATCH | 0 | 0 |
| plugin-phpcs | install | MATCH | 0 | 0 |
| plugin-phpcs | update-no-install | MATCH | 0 | 0 |
| plugin-pest | install | MATCH | 0 | 0 |
| plugin-pest | update-no-install | MATCH | 0 | 0 |
