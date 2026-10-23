# Composer vs gocomposer comparison

- gocomposer_sha256: `cce199076e97cf1cb951f4c66df7971ce683ad8b9cc3751f26f4f2c3c7880eb9`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| show-latest-cooldown-long | MATCH | 1 | 1 |
| show-latest-cooldown-long-json | MATCH | 1 | 1 |
| show-latest-cooldown-short | MATCH | 1 | 1 |
| show-latest-cooldown-zero | MATCH | 0 | 0 |
| show-latest-cooldown-invalid | MATCH | 1 | 1 |
| show-latest-cooldown-package | MATCH | 1 | 1 |
| show-latest-cooldown-time | MATCH | 0 | 0 |
| show-latest-cooldown-time-json | MATCH | 0 | 0 |
| show-latest-cooldown-time-strict | MATCH | 1 | 1 |
| show-latest-cooldown-time-ansi | MATCH | 0 | 0 |
| show-latest-cooldown-time-verbose | MATCH | 0 | 0 |
| show-latest-cooldown-time-passed | MATCH | 0 | 0 |
| show-latest-cooldown-time-package | MATCH | 0 | 0 |
| show-latest-cooldown-time-package-json | MATCH | 0 | 0 |
| show-latest-cooldown-config | MATCH | 0 | 0 |
| show-latest-cooldown-config-json | MATCH | 0 | 0 |
| outdated-cooldown-config | MATCH | 1 | 1 |
| outdated-cooldown-config-zero | MATCH | 0 | 0 |
