# Composer vs gocomposer comparison

- gocomposer_sha256: `98791a1c013aa911f495808acb5fc8d19f2459501cdada771322535bbf531f81`
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
| root-alias | MATCH | 0 | 0 |
| misc-list-root | MATCH | 0 | 0 |
| misc-list-root-noallow | MATCH | 1 | 1 |
| misc-help-root | MATCH | 0 | 0 |
| misc-list-tty | MATCH | 0 | 0 |
| misc-help-tty | MATCH | 0 | 0 |
| misc-help-install-tty | MATCH | 0 | 0 |
| misc-help-script | MATCH | 0 | 0 |
| misc-help-script-list | MATCH | 0 | 0 |
| misc-help-script-tty | MATCH | 0 | 0 |
| misc-list-scripts | MATCH | 0 | 0 |
| misc-version | MATCH | 0 | 0 |
| misc-version-short | MATCH | 0 | 0 |
| misc-version-ansi | MATCH | 0 | 0 |
| misc-no-command | MATCH | 0 | 0 |
| misc-h-only | MATCH | 0 | 0 |
| misc-h-only-json | MATCH | 0 | 0 |
| misc-list-raw | MATCH | 0 | 0 |
| misc-list-ansi | MATCH | 0 | 0 |
| misc-list-quiet | MATCH | 0 | 0 |
| misc-list-json | MATCH | 0 | 0 |
| misc-list-json-short | MATCH | 0 | 0 |
| misc-list-xml | MATCH | 0 | 0 |
| misc-list-xml-short | MATCH | 0 | 0 |
| misc-list-md | MATCH | 0 | 0 |
| misc-list-md-short | MATCH | 0 | 0 |
| misc-list-txt-short | MATCH | 0 | 0 |
| misc-list-bad-format | MATCH | 1 | 1 |
| misc-list-help | MATCH | 0 | 0 |
| misc-list-missing-namespace | MATCH | 1 | 1 |
| misc-list-ns-scripts | MATCH | 0 | 0 |
| misc-list-ns | MATCH | 0 | 0 |
| misc-list-ns-raw | MATCH | 0 | 0 |
| misc-list-ns-nested | MATCH | 0 | 0 |
| misc-list-ns-abbrev | MATCH | 0 | 0 |
| misc-list-ns-json | MATCH | 0 | 0 |
| misc-list-ns-xml | MATCH | 0 | 0 |
| misc-list-ns-md | MATCH | 0 | 0 |
| misc-list-ns-scripts-raw | MATCH | 0 | 0 |
| misc-list-ns-scripts-xml | MATCH | 0 | 0 |
| misc-help-list | MATCH | 0 | 0 |
| misc-help-raw | MATCH | 0 | 0 |
| misc-help-json | MATCH | 0 | 0 |
| misc-help-xml | MATCH | 0 | 0 |
| misc-help-md | MATCH | 0 | 0 |
| misc-help-md-ansi | MATCH | 0 | 0 |
| misc-help-txt-ansi | MATCH | 0 | 0 |
| misc-help-bad-format | MATCH | 1 | 1 |
| misc-help-missing | MATCH | 1 | 1 |
| misc-help-abbrev | MATCH | 0 | 0 |
| misc-help-ambiguous | MATCH | 1 | 1 |
| misc-help-alias | MATCH | 0 | 0 |
| misc-help-hidden | MATCH | 0 | 0 |
| misc-help-completion | MATCH | 0 | 0 |
| misc-help-config | MATCH | 0 | 0 |
| misc-help-show-json | MATCH | 0 | 0 |
| misc-help-script-json | MATCH | 0 | 0 |
| misc-help-script-xml | MATCH | 0 | 0 |
| misc-help-script-md | MATCH | 0 | 0 |
| misc-help-flag | MATCH | 0 | 0 |
| misc-help-flag-short | MATCH | 0 | 0 |
| misc-help-flag-json | MATCH | 0 | 0 |
| misc-help-flag-raw | MATCH | 0 | 0 |
| misc-help-flag-script | MATCH | 0 | 0 |
| misc-help-ns-script | MATCH | 0 | 0 |
| misc-help-ns-alias | MATCH | 0 | 0 |
| misc-list-plugin | MATCH | 0 | 0 |
| misc-list-plugin-ns | MATCH | 0 | 0 |
| misc-list-plugin-json | MATCH | 0 | 0 |
| misc-list-plugin-raw | MATCH | 0 | 0 |
| misc-help-plugin | MATCH | 0 | 0 |
| misc-help-plugin-json | MATCH | 0 | 0 |
| misc-help-plugin-xml | MATCH | 0 | 0 |
| misc-help-plugin-md | MATCH | 0 | 0 |
| misc-help-plugin-flag | MATCH | 0 | 0 |
| run-script-basic | MATCH | 0 | 0 |
| run-script-alias-run | MATCH | 0 | 0 |
| run-script-abbrev | MATCH | 0 | 0 |
| run-script-args | MATCH | 0 | 0 |
| run-script-args-dashdash | MATCH | 0 | 0 |
| run-script-list | MATCH | 0 | 0 |
| run-script-list-short | MATCH | 0 | 0 |
| run-script-missing-argument | MATCH | 1 | 1 |
| run-script-not-defined | MATCH | 1 | 1 |
| run-script-event-constant-name | MATCH | 1 | 1 |
| run-script-event | MATCH | 0 | 0 |
| run-script-event-array | MATCH | 0 | 0 |
| run-script-package-event | MATCH | 0 | 0 |
| run-script-dev-mode | MATCH | 0 | 0 |
| run-script-no-dev | MATCH | 0 | 0 |
| run-script-dev-and-no-dev | MATCH | 0 | 0 |
| run-script-no-dev-env | MATCH | 0 | 0 |
| run-script-timeout-invalid | MATCH | 1 | 1 |
| run-script-timeout-negative | MATCH | 1 | 1 |
| run-script-timeout | MATCH | 1 | 1 |
| run-script-timeout-zero | MATCH | 0 | 0 |
| run-script-process-timeout-env | MATCH | 1 | 1 |
| run-script-fail | MATCH | 3 | 3 |
| run-script-nested | MATCH | 0 | 0 |
| run-script-nested-verbose | MATCH | 0 | 0 |
| run-script-missing-ref | MATCH | 0 | 0 |
| run-script-multi-args | MATCH | 0 | 0 |
| run-script-no-additional-args | MATCH | 0 | 0 |
| run-script-putenv | MATCH | 0 | 0 |
| run-script-no-scripts | MATCH | 1 | 1 |
| run-script-skip-env | MATCH | 1 | 1 |
| run-script-skip-env-verbose | MATCH | 0 | 0 |
| run-script-quiet | MATCH | 0 | 0 |
| run-script-direct | MATCH | 0 | 0 |
| run-script-direct-args | MATCH | 0 | 0 |
| run-script-direct-unknown-option | MATCH | 0 | 0 |
| run-script-direct-dashdash | MATCH | 0 | 0 |
| run-script-direct-alias | MATCH | 0 | 0 |
| run-script-direct-alias2 | MATCH | 0 | 0 |
| run-script-direct-nested | MATCH | 0 | 0 |
| run-script-direct-fail | MATCH | 3 | 3 |
| run-script-direct-dev-mode | MATCH | 0 | 0 |
| run-script-direct-event-not-command | MATCH | 1 | 1 |
| run-script-direct-abbrev | MATCH | 0 | 0 |
| run-script-unknown-command | MATCH | 1 | 1 |
| run-script-override-warning | MATCH | 0 | 0 |
| run-script-aliases-invalid | MATCH | 1 | 1 |
| run-script-list-empty | MATCH | 0 | 0 |
| run-script-no-json | MATCH | 1 | 1 |
| run-script-disable-timeout | MATCH | 0 | 0 |
| run-script-working-dir | MATCH | 0 | 0 |
| plugin-command-unknown-ns | MATCH | 1 | 1 |
| root-validate | MATCH | 0 | 0 |
| root-validate-allowed | MATCH | 0 | 0 |
| root-validate-allowed-zero | MATCH | 0 | 0 |
| root-validate-no-plugins | MATCH | 0 | 0 |
| root-update | MATCH | 0 | 0 |
| root-unknown-command | MATCH | 1 | 1 |
| root-unknown-command-allowed | MATCH | 1 | 1 |
| root-help | MATCH | 1 | 1 |
| root-script-alias | MATCH | 1 | 1 |
| root-tty-validate-yes | MATCH | 0 | 0 |
| root-tty-validate-default | MATCH | 0 | 0 |
| root-tty-validate-no | MATCH | 1 | 1 |
| root-tty-validate-n | MATCH | 0 | 0 |
| root-tty-validate-allowed | MATCH | 0 | 0 |
| root-tty-unknown-command-yes | MATCH | 1 | 1 |
| root-tty-unknown-command-no | MATCH | 1 | 1 |
| root-tty-script-alias-yes | MATCH | 0 | 0 |
| root-tty-update-yes | MATCH | 0 | 0 |
| root-tty-no-plugins | MATCH | 0 | 0 |
