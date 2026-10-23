# Composer vs gocomposer comparison

- gocomposer_sha256: `554b103a6dfa04a5178732c182035c155e8eb33323634a215657b66541aa8a96`
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
| plugin-blocked-update | MATCH | 1 | 1 |
| plugin-blocked-install-lock | MATCH | 1 | 1 |
| plugin-blocked-installed | MATCH | 1 | 1 |
| plugin-allow-false | MATCH | 0 | 0 |
| plugin-allow-all-false | MATCH | 0 | 0 |
| plugin-allow-true | MATCH | 0 | 0 |
| plugin-optional-blocked | MATCH | 0 | 0 |
| plugin-not-required-anymore | MATCH | 0 | 0 |
| plugin-old-lock-no-allow | MATCH | 1 | 1 |
| plugin-old-lock-new-api | MATCH | 1 | 1 |
| plugin-allowed-update | MATCH | 0 | 0 |
| plugin-allowed-update-verbose | MATCH | 0 | 0 |
| plugin-update-dry-run | MATCH | 0 | 0 |
| plugin-install-from-lock | MATCH | 0 | 0 |
| plugin-install-no-dev | MATCH | 0 | 0 |
| plugin-installed-noop | MATCH | 0 | 0 |
| plugin-installed-quiet | MATCH | 0 | 0 |
| plugin-no-scripts | MATCH | 0 | 0 |
| plugin-dump-autoload | MATCH | 0 | 0 |
| plugin-dump-autoload-optimize | MATCH | 0 | 0 |
| plugin-no-plugins-installed | MATCH | 0 | 0 |
| plugin-no-plugins-update | MATCH | 0 | 0 |
| plugin-update-version | MATCH | 0 | 0 |
| plugin-remove | MATCH | 0 | 0 |
| plugin-require | MATCH | 0 | 0 |
| plugin-require-installed | MATCH | 0 | 0 |
| plugin-run-script-custom | MATCH | 0 | 0 |
| plugin-show | MATCH | 0 | 0 |
| plugin-thrower-boom | MATCH | 1 | 1 |
| plugin-thrower-stop | MATCH | 0 | 0 |
| plugin-thrower-false | MATCH | 0 | 0 |
| plugin-global | MATCH | 0 | 0 |
| plugin-global-blocked | MATCH | 1 | 1 |
| plugin-global-no-plugins | MATCH | 0 | 0 |
| plugin-autoload-update | MATCH | 0 | 0 |
| plugin-autoload-dump | MATCH | 0 | 0 |
| plugin-autoload-run-script | MATCH | 0 | 0 |
| plugin-old-api | MATCH | 2 | 2 |
| plugin-no-api | MATCH | 1 | 1 |
| plugin-no-class | MATCH | 1 | 1 |
| plugin-legacy-installer | MATCH | 0 | 0 |
| plugin-command | MATCH | 0 | 0 |
| plugin-command-default | MATCH | 0 | 0 |
| plugin-command-exit | MATCH | 3 | 3 |
| plugin-command-bad-option | MATCH | 1 | 1 |
| plugin-command-no-plugins | MATCH | 1 | 1 |
| plugin-command-conflict | MATCH | 0 | 0 |
| plugin-command-alias | MATCH | 0 | 0 |
| plugin-command-throws | MATCH | 1 | 1 |
| plugin-command-unknown-ns | MATCH | 1 | 1 |
| plugin-require-blocked | MATCH | 1 | 1 |
| plugin-installer-update | MATCH | 0 | 0 |
| plugin-installer-extra-path | MATCH | 0 | 0 |
| plugin-installer-install | MATCH | 0 | 0 |
| plugin-installer-upgrade | MATCH | 0 | 0 |
| plugin-installer-remove | MATCH | 0 | 0 |
| plugin-installer-no-plugins | MATCH | 0 | 0 |
| plugin-installer-dump | MATCH | 0 | 0 |
