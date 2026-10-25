# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `7f4937d7dbb52ad9625817fdf2940c03b572907c6a00d7234e43e7b34f92449d`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:9cf420cd931841442e55100b424e04a60af91810d670826961c6583bd372308b`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| project | step | status | oracle exit | go exit |
|---|---|---|---|---|
| library-symfony-console | install | MATCH | 0 | 0 |
| library-symfony-console | update-no-install | MATCH | 0 | 0 |
| library-symfony-console | update-no-install-no-scripts | MATCH | 0 | 0 |
| cli-app | install | MATCH | 0 | 0 |
| cli-app | update-no-install | MATCH | 255 | 255 |
| cli-app | update-no-install-no-scripts | MATCH | 0 | 0 |
| laravel | install | MATCH | 0 | 0 |
| laravel | update-no-install | MATCH | 255 | 255 |
| laravel | update-no-install-no-scripts | MATCH | 0 | 0 |
| live-psr | install | MATCH | 0 | 0 |
| live-psr | update-no-install | MATCH | 0 | 0 |
| live-psr | update-no-install-no-scripts | MATCH | 0 | 0 |
| plugin-installers | install | MATCH | 0 | 0 |
| plugin-installers | update-no-install | MATCH | 0 | 0 |
| plugin-installers | update-no-install-no-scripts | MATCH | 0 | 0 |
| plugin-php-http-discovery | install | MATCH | 0 | 0 |
| plugin-php-http-discovery | update-no-install | MATCH | 0 | 0 |
| plugin-php-http-discovery | update-no-install-no-scripts | MATCH | 0 | 0 |
| plugin-phpcs | install | MATCH | 0 | 0 |
| plugin-phpcs | update-no-install | MATCH | 0 | 0 |
| plugin-phpcs | update-no-install-no-scripts | MATCH | 0 | 0 |
| plugin-pest | install | MATCH | 0 | 0 |
| plugin-pest | update-no-install | MATCH | 0 | 0 |
| plugin-pest | update-no-install-no-scripts | MATCH | 0 | 0 |
| symfony-skeleton | install | MATCH | 0 | 0 |
| symfony-skeleton | update-no-install | MATCH | 0 | 0 |
| symfony-skeleton | update-no-install-no-scripts | MATCH | 0 | 0 |
| symfony-skeleton | recipes-install | MATCH | 0 | 0 |
