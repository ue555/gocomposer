# Composer vs gocomposer comparison

- gocomposer_sha256: `2f3d75915ee647f1248ac02751714cfced68e424ac5ee6799bbada66e1524e96`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| show-latest-short | MATCH | 0 | 0 |
| show-latest-outdated | MATCH | 0 | 0 |
| show-latest-outdated-direct | MATCH | 0 | 0 |
| show-latest-direct | MATCH | 0 | 0 |
| show-latest-strict | MATCH | 1 | 1 |
| show-latest-strict-latest-only | MATCH | 1 | 1 |
| show-latest-strict-ignored | MATCH | 0 | 0 |
| show-latest-strict-fresh | MATCH | 0 | 0 |
| show-latest-fresh | MATCH | 0 | 0 |
| show-latest-strict-no-latest | MATCH | 1 | 1 |
| show-latest-major-only | MATCH | 0 | 0 |
| show-latest-major-only-latest | MATCH | 0 | 0 |
| show-latest-minor-only | MATCH | 0 | 0 |
| show-latest-patch-only | MATCH | 0 | 0 |
| show-latest-patch-only-outdated | MATCH | 0 | 0 |
| show-latest-only-conflict | MATCH | 1 | 1 |
| show-latest-json | MATCH | 0 | 0 |
| show-latest-outdated-json | MATCH | 0 | 0 |
| show-latest-outdated-json-strict | MATCH | 1 | 1 |
| show-latest-name-only | MATCH | 0 | 0 |
| show-latest-name-only-json | MATCH | 0 | 0 |
| show-latest-path | MATCH | 0 | 0 |
| show-latest-wildcard | MATCH | 0 | 0 |
| show-latest-ignore | MATCH | 0 | 0 |
| show-latest-ignore-without-outdated | MATCH | 0 | 0 |
| show-latest-no-dev | MATCH | 0 | 0 |
| show-latest-sort-by-age | MATCH | 0 | 0 |
| show-latest-ansi | MATCH | 0 | 0 |
| show-latest-outdated-ansi | MATCH | 0 | 0 |
| show-latest-verbose | MATCH | 0 | 0 |
| show-latest-ignore-platform-reqs | MATCH | 0 | 0 |
| show-latest-ignore-platform-req | MATCH | 0 | 0 |
| show-latest-platform-rejected | MATCH | 0 | 0 |
| show-latest-platform-rejected-verbose | MATCH | 0 | 0 |
| show-latest-locked | MATCH | 0 | 0 |
| show-latest-locked-outdated | MATCH | 0 | 0 |
| show-latest-not-installed | MATCH | 0 | 0 |
| show-latest-tree | MATCH | 1 | 1 |
| show-latest-platform | MATCH | 0 | 0 |
| show-latest-v2-metadata | MATCH | 0 | 0 |
| show-latest-dev-branch | MATCH | 0 | 0 |
| show-latest-dev-branch-major | MATCH | 0 | 0 |
| show-latest-dev-branch-package | MATCH | 0 | 0 |
| show-latest-no-composer | MATCH | 0 | 0 |
| show-latest-cooldown-long | MATCH | 1 | 1 |
| show-latest-cooldown-long-json | MATCH | 1 | 1 |
| show-latest-cooldown-short | MATCH | 1 | 1 |
| show-latest-cooldown-zero | MATCH | 0 | 0 |
| show-latest-cooldown-invalid | MATCH | 1 | 1 |
| show-latest-cooldown-package | MATCH | 1 | 1 |
| show-latest-package | MATCH | 0 | 0 |
| show-latest-package-json | MATCH | 0 | 0 |
| show-latest-package-uptodate | MATCH | 0 | 0 |
| show-latest-package-abandoned | MATCH | 0 | 0 |
| show-latest-package-abandoned-json | MATCH | 0 | 0 |
| show-latest-package-gone | MATCH | 0 | 0 |
| show-latest-package-strict | MATCH | 1 | 1 |
| show-latest-package-strict-uptodate | MATCH | 0 | 0 |
| show-latest-package-major | MATCH | 0 | 0 |
| show-latest-package-patch | MATCH | 0 | 0 |
| show-latest-package-path | MATCH | 1 | 1 |
| show-latest-package-direct-transitive | MATCH | 1 | 1 |
| show-latest-package-platform-rejected | MATCH | 0 | 0 |
| show-latest-all | MATCH | 0 | 0 |
| show-latest-all-json | MATCH | 0 | 0 |
| show-latest-all-name-only | MATCH | 0 | 0 |
| show-latest-all-wildcard | MATCH | 0 | 0 |
| show-latest-all-package | MATCH | 0 | 0 |
| show-latest-all-package-version | MATCH | 0 | 0 |
| show-latest-all-no-lock | MATCH | 0 | 0 |
| show-latest-all-latest | MATCH | 0 | 0 |
| show-latest-all-strict | MATCH | 1 | 1 |
| show-latest-available | MATCH | 0 | 0 |
| show-latest-available-json | MATCH | 0 | 0 |
| show-latest-available-wildcard | MATCH | 0 | 0 |
| show-latest-available-package | MATCH | 0 | 0 |
| show-latest-available-package-json | MATCH | 0 | 0 |
| show-latest-available-not-found | MATCH | 1 | 1 |
| show-latest-available-v2 | MATCH | 0 | 0 |
| show-latest-available-direct | MATCH | 1 | 1 |
| show-latest-available-tree | MATCH | 1 | 1 |
| show-latest-available-latest | MATCH | 0 | 0 |
| outdated-default | MATCH | 0 | 0 |
| outdated-all | MATCH | 0 | 0 |
| outdated-direct | MATCH | 0 | 0 |
| outdated-strict | MATCH | 1 | 1 |
| outdated-strict-fresh | MATCH | 0 | 0 |
| outdated-major-only | MATCH | 0 | 0 |
| outdated-minor-only | MATCH | 0 | 0 |
| outdated-patch-only | MATCH | 0 | 0 |
| outdated-only-conflict | MATCH | 1 | 1 |
| outdated-json | MATCH | 0 | 0 |
| outdated-all-json | MATCH | 0 | 0 |
| outdated-bad-format | MATCH | 1 | 1 |
| outdated-package | MATCH | 0 | 0 |
| outdated-package-strict | MATCH | 1 | 1 |
| outdated-package-not-found | MATCH | 1 | 1 |
| outdated-wildcard | MATCH | 0 | 0 |
| outdated-locked | MATCH | 0 | 0 |
| outdated-no-dev | MATCH | 0 | 0 |
| outdated-ignore | MATCH | 0 | 0 |
| outdated-ignore-platform-reqs | MATCH | 0 | 0 |
| outdated-ignore-platform-req | MATCH | 0 | 0 |
| outdated-sort-by-age | MATCH | 0 | 0 |
| outdated-cooldown | MATCH | 1 | 1 |
| outdated-ansi | MATCH | 0 | 0 |
| outdated-not-installed | MATCH | 0 | 0 |
| outdated-no-composer | MATCH | 1 | 1 |
| outdated-unknown-option | MATCH | 1 | 1 |
