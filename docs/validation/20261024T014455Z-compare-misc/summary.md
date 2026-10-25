# Composer vs gocomposer comparison

- gocomposer_sha256: `9be6cabb6baf9b5ae2b45dfc6188021086fc008ccc8874d953fa5749ef69561d`
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
