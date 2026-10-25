# Composer vs gocomposer comparison

- gocomposer_sha256: `7f4937d7dbb52ad9625817fdf2940c03b572907c6a00d7234e43e7b34f92449d`
- host_os: `darwin/arm64`
- oracle_arch_image: `localhost/gocomposer-oracle:8.4-arch`
- oracle_arch_image_digest: `sha256:93bc04d027eaeb74edeb8ea07a247c5f17cb2e91d4d17fe97681f949f96116e7`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:9cf420cd931841442e55100b424e04a60af91810d670826961c6583bd372308b`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:e3a99ff1ce724bd01dbd722f202bf4afa98eec9bde0789c99ad6f76b266e6f22`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| cli3-vvv-validate | MATCH | 0 | 0 |
| cli3-vvv-show | MATCH | 0 | 0 |
| cli3-vvv-config-list | MATCH | 0 | 0 |
| cli3-vvv-config-get | MATCH | 0 | 0 |
| cli3-vvv-config-set | MATCH | 0 | 0 |
| cli3-vvv-config-extra | MATCH | 0 | 0 |
| cli3-vvv-config-global | MATCH | 0 | 0 |
| cli3-vvv-dump-autoload | MATCH | 0 | 0 |
| cli3-vvv-run-script | MATCH | 0 | 0 |
| cli3-vvv-script-command | MATCH | 0 | 0 |
| cli3-vvv-licenses | MATCH | 0 | 0 |
| cli3-vvv-depends | MATCH | 0 | 0 |
| cli3-vvv-suggests | MATCH | 0 | 0 |
| cli3-vvv-no-plugins | MATCH | 0 | 0 |
| cli3-vvv-shell-verbosity | MATCH | 0 | 0 |
| cli3-vvv-no-lock | MATCH | 0 | 0 |
| cli3-vvv-auth-files | MATCH | 0 | 0 |
| cli3-vvv-cafile | MATCH | 0 | 0 |
| cli3-vvv-ssl-cert-file | MATCH | 0 | 0 |
| cli3-vvv-ssl-cert-missing | MATCH | 0 | 0 |
| cli3-vvv-global-home | MATCH | 0 | 0 |
| cli3-vvv-allow-plugins | MATCH | 0 | 0 |
| cli3-vvv-validate-file | MATCH | 0 | 0 |
| cli3-bin-compat-get | MATCH | 0 | 0 |
| cli3-bin-compat-list | MATCH | 0 | 0 |
| cli3-bin-compat-set | MATCH | 0 | 0 |
| cli3-bin-compat-install | MATCH | 0 | 0 |
| cli3-bin-compat-env | MATCH | 0 | 0 |
| cli3-event-validate | MATCH | 0 | 0 |
| cli3-event-validate-no-plugins | MATCH | 0 | 0 |
| cli3-event-reinstall | MATCH | 0 | 0 |
| cli3-event-reinstall-none | MATCH | 1 | 1 |
| cli3-cmdclass-run | MATCH | 0 | 0 |
| cli3-cmdclass-run-args | MATCH | 0 | 0 |
| cli3-cmdclass-run-verbose | MATCH | 0 | 0 |
| cli3-cmdclass-run-fail | MATCH | 3 | 3 |
| cli3-cmdclass-run-boom | MATCH | 1 | 1 |
| cli3-cmdclass-run-missing | MATCH | 0 | 0 |
| cli3-cmdclass-run-notcmd | MATCH | 0 | 0 |
| cli3-cmdclass-run-named | MATCH | 0 | 0 |
| cli3-cmdclass-run-nested | MATCH | 0 | 0 |
| cli3-cmdclass-reserved | MATCH | 0 | 0 |
| cli3-cmdclass-list | MATCH | 0 | 0 |
| cli3-cmdclass-direct | MATCH | 0 | 0 |
| cli3-cmdclass-direct-options | MATCH | 0 | 0 |
| cli3-cmdclass-direct-bad-option | MATCH | 1 | 1 |
| cli3-cmdclass-direct-fail | MATCH | 3 | 3 |
| cli3-cmdclass-direct-boom | MATCH | 1 | 1 |
| cli3-cmdclass-direct-named | MATCH | 0 | 0 |
| cli3-cmdclass-not-installed | MATCH | 0 | 0 |
| cli3-packagist-search | MATCH | 0 | 0 |
| cli3-packagist-search-json | MATCH | 0 | 0 |
| cli3-packagist-search-all | MATCH | 0 | 0 |
| cli3-packagist-search-name | MATCH | 0 | 0 |
| cli3-packagist-search-vendor | MATCH | 0 | 0 |
| cli3-packagist-search-none | MATCH | 0 | 0 |
| cli3-packagist-browse-show | MATCH | 0 | 0 |
| cli3-packagist-browse-homepage | MATCH | 0 | 0 |
| cli3-packagist-browse-source | MATCH | 1 | 1 |
| cli3-packagist-browse-missing | MATCH | 1 | 1 |
| cli3-packagist-browse-root | MATCH | 1 | 1 |
| cli3-packagist-fund | MATCH | 0 | 0 |
| cli3-packagist-fund-json | MATCH | 0 | 0 |
| cli3-packagist-fund-no-json | MATCH | 1 | 1 |
