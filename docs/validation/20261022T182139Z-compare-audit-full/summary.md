# Composer vs gocomposer comparison

- gocomposer_sha256: `e6a0d9e271a6c37be9ff2fc01112d5a4100a720c732e083e4299efe54a079240`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| update-basic | MATCH | 0 | 0 |
| update-basic-progress | MATCH | 0 | 0 |
| update-v2-metadata | MATCH | 0 | 0 |
| install-from-lock | MATCH | 0 | 0 |
| install-no-dev | MATCH | 0 | 0 |
| install-twice | MATCH | 0 | 0 |
| dry-run-update | MATCH | 0 | 0 |
| dump-autoload-optimized | MATCH | 0 | 0 |
| dump-autoload-authoritative | MATCH | 0 | 0 |
| replace | MATCH | 0 | 0 |
| conflict-unresolvable | MATCH | 2 | 2 |
| conflict-resolves-lower | MATCH | 0 | 0 |
| platform-php-too-low | MATCH | 2 | 2 |
| missing-extension | MATCH | 2 | 2 |
| ignore-platform-reqs | MATCH | 0 | 0 |
| config-platform-php | MATCH | 0 | 0 |
| not-found-package | MATCH | 2 | 2 |
| stability-beta-flag | MATCH | 0 | 0 |
| stability-too-low | MATCH | 2 | 2 |
| dev-branch-alias | MATCH | 0 | 0 |
| root-alias | MATCH | 0 | 0 |
| partial-update | MATCH | 0 | 0 |
| partial-update-with-deps | MATCH | 0 | 0 |
| broken-archive | MATCH | 100 | 100 |
| checksum-mismatch | MATCH | 1 | 1 |
| download-404 | MATCH | 100 | 100 |
| path-traversal-archive | MATCH | 0 | 0 |
| failure-keeps-vendor | MATCH | 2 | 2 |
| no-lock-install | MATCH | 0 | 0 |
| check-platform-reqs-ok | MATCH | 0 | 0 |
| check-platform-reqs-fail | MATCH | 2 | 2 |
| check-platform-reqs-json | MATCH | 2 | 2 |
| check-platform-reqs-lock | MATCH | 2 | 2 |
| validate-basic | MATCH | 0 | 0 |
| validate-clean | MATCH | 0 | 0 |
| validate-warnings | MATCH | 1 | 1 |
| validate-strict-warnings | MATCH | 1 | 1 |
| validate-no-check | MATCH | 0 | 0 |
| validate-missing-file | MATCH | 3 | 3 |
| validate-invalid-json | MATCH | 1 | 1 |
| validate-duplicate-key | MATCH | 0 | 0 |
| validate-lock-stale | MATCH | 2 | 2 |
| validate-lock-stale-no-check | MATCH | 0 | 0 |
| validate-with-dependencies | MATCH | 0 | 0 |
| validate-bad-constraint | MATCH | 1 | 1 |
| show-list | MATCH | 0 | 0 |
| show-json | MATCH | 0 | 0 |
| show-name-only | MATCH | 0 | 0 |
| show-direct | MATCH | 0 | 0 |
| show-no-dev | MATCH | 0 | 0 |
| show-path | MATCH | 0 | 0 |
| show-package | MATCH | 0 | 0 |
| show-package-json | MATCH | 0 | 0 |
| show-tree | MATCH | 0 | 0 |
| show-tree-json | MATCH | 0 | 0 |
| show-package-tree | MATCH | 0 | 0 |
| show-locked | MATCH | 0 | 0 |
| show-self | MATCH | 0 | 0 |
| show-self-json | MATCH | 0 | 0 |
| show-platform-php | MATCH | 0 | 0 |
| show-wildcard | MATCH | 0 | 0 |
| show-not-found | MATCH | 1 | 1 |
| show-version-arg | MATCH | 0 | 0 |
| show-self-name-only | MATCH | 0 | 0 |
| show-locked-package | MATCH | 0 | 0 |
| show-platform-list-json | MATCH | 0 | 0 |
| show-bad-format | MATCH | 1 | 1 |
| show-tree-latest | MATCH | 1 | 1 |
| show-nothing-installed | MATCH | 0 | 0 |
| depends-basic | MATCH | 0 | 0 |
| depends-recursive | MATCH | 0 | 0 |
| depends-tree | MATCH | 0 | 0 |
| depends-php | MATCH | 0 | 0 |
| depends-notfound | MATCH | 1 | 1 |
| depends-locked | MATCH | 0 | 0 |
| prohibits-basic | MATCH | 1 | 1 |
| prohibits-installed | MATCH | 0 | 0 |
| prohibits-php | MATCH | 1 | 1 |
| prohibits-tree | MATCH | 1 | 1 |
| why-alias | MATCH | 0 | 0 |
| depends-nothing-installed | MATCH | 1 | 1 |
| licenses-text | MATCH | 0 | 0 |
| licenses-json | MATCH | 0 | 0 |
| licenses-summary | MATCH | 0 | 0 |
| licenses-no-dev | MATCH | 0 | 0 |
| licenses-locked | MATCH | 0 | 0 |
| licenses-bad-format | MATCH | 1 | 1 |
| suggests-default | MATCH | 0 | 0 |
| suggests-all | MATCH | 0 | 0 |
| suggests-list | MATCH | 0 | 0 |
| suggests-by-suggestion | MATCH | 0 | 0 |
| suggests-both | MATCH | 0 | 0 |
| suggests-missing | MATCH | 1 | 1 |
| suggests-filter | MATCH | 0 | 0 |
| validate-publish-errors | MATCH | 1 | 1 |
| validate-lax-fail | KNOWN-DIFF | 1 | 1 |
| validate-lax-fail-update | MATCH | 1 | 1 |
| composer-auth-invalid | KNOWN-DIFF | 255 | 0 |
| composer-auth-malformed | MATCH | 1 | 1 |
| policy-advisory-blocks | MATCH | 0 | 0 |
| policy-advisory-blocks-vv | MATCH | 0 | 0 |
| policy-advisory-only-vulnerable | MATCH | 2 | 2 |
| policy-no-blocking-audit | MATCH | 0 | 0 |
| policy-env-no-blocking | MATCH | 0 | 0 |
| policy-audit-install | MATCH | 5 | 5 |
| policy-audit-install-table | MATCH | 5 | 5 |
| policy-audit-install-json | MATCH | 5 | 5 |
| policy-audit-plain | MATCH | 5 | 5 |
| policy-abandoned | MATCH | 0 | 0 |
| policy-abandoned-ignore | MATCH | 0 | 0 |
| policy-ignore-advisory | MATCH | 0 | 0 |
| policy-disabled | MATCH | 0 | 0 |
| policy-cooldown-all | MATCH | 1 | 1 |
| policy-cooldown-none | MATCH | 1 | 1 |
| policy-cooldown-no-cooldown | MATCH | 0 | 0 |
| policy-locked-vulnerable-install | MATCH | 0 | 0 |
| policy-partial-update-vulnerable | MATCH | 0 | 0 |
| require-guess-version | MATCH | 0 | 0 |
| require-explicit | MATCH | 0 | 0 |
| require-dev | MATCH | 0 | 0 |
| require-no-update | MATCH | 0 | 0 |
| require-dry-run | MATCH | 0 | 0 |
| require-first-no-lock | MATCH | 0 | 0 |
| require-no-composer-json | MATCH | 1 | 1 |
| require-not-found | MATCH | 1 | 1 |
| require-min-stability | MATCH | 1 | 1 |
| require-platform-incompatible | MATCH | 1 | 1 |
| require-fails-reverts | MATCH | 2 | 2 |
| require-move-from-dev | MATCH | 0 | 0 |
| require-sort-packages | MATCH | 0 | 0 |
| require-too-strict | MATCH | 0 | 0 |
| require-self | MATCH | 1 | 1 |
| require-ext-missing | MATCH | 1 | 1 |
| require-alias-arg | MATCH | 1 | 1 |
| require-fixed-library | MATCH | 1 | 1 |
| require-fixed-project | MATCH | 0 | 0 |
| require-update-with-deps | MATCH | 0 | 0 |
| require-crlf-file | MATCH | 0 | 0 |
| require-tabs-file | MATCH | 0 | 0 |
| remove-basic | MATCH | 0 | 0 |
| remove-dev | MATCH | 0 | 0 |
| remove-wrong-section | MATCH | 2 | 2 |
| remove-not-required | MATCH | 0 | 0 |
| remove-wildcard | MATCH | 0 | 0 |
| remove-no-update | MATCH | 0 | 0 |
| remove-dry-run | MATCH | 0 | 0 |
| remove-still-required | MATCH | 2 | 2 |
| remove-unused | MATCH | 0 | 0 |
| remove-no-args | MATCH | 1 | 1 |
| remove-mixed-case | MATCH | 0 | 0 |
| remove-allow-plugins | MATCH | 0 | 0 |
| config-no-args | MATCH | 0 | 0 |
| config-get-vendor-dir | MATCH | 0 | 0 |
| config-get-bin-dir | MATCH | 0 | 0 |
| config-get-bin-dir-absolute | MATCH | 0 | 0 |
| config-get-bin-dir-source | MATCH | 0 | 0 |
| config-get-set-in-json | MATCH | 0 | 0 |
| config-get-bin-dir-custom-vendor | MATCH | 0 | 0 |
| config-get-process-timeout | MATCH | 0 | 0 |
| config-get-bool | MATCH | 0 | 0 |
| config-get-bool-default | MATCH | 0 | 0 |
| config-get-allow-plugins-empty | MATCH | 0 | 0 |
| config-get-platform-empty | MATCH | 0 | 0 |
| config-get-secure-svn-domains | MATCH | 0 | 0 |
| config-get-github-protocols | MATCH | 0 | 0 |
| config-get-github-protocols-secure | MATCH | 0 | 0 |
| config-get-github-protocols-git-first | MATCH | 0 | 0 |
| config-get-gitlab-domains-dedupe | MATCH | 0 | 0 |
| config-get-audit | MATCH | 0 | 0 |
| config-get-audit-nested | MATCH | 0 | 0 |
| config-get-nested-missing | MATCH | 1 | 1 |
| config-get-extra-nested | MATCH | 0 | 0 |
| config-get-extra-object | MATCH | 0 | 0 |
| config-get-extra | MATCH | 0 | 0 |
| config-get-extra-list-index | MATCH | 0 | 0 |
| config-get-suggest | MATCH | 0 | 0 |
| config-get-name | MATCH | 0 | 0 |
| config-get-license | MATCH | 0 | 0 |
| config-get-homepage-default | MATCH | 0 | 0 |
| config-get-min-stability-default | MATCH | 0 | 0 |
| config-get-prefer-stable-default | MATCH | 0 | 0 |
| config-get-keywords | MATCH | 0 | 0 |
| config-get-repositories | MATCH | 0 | 0 |
| config-get-repositories-default | MATCH | 0 | 0 |
| config-get-repo-index | MATCH | 0 | 0 |
| config-get-repo-missing | MATCH | 1 | 1 |
| config-get-repo-prefix | MATCH | 0 | 0 |
| config-get-undefined | MATCH | 1 | 1 |
| config-get-undefined-dotted | MATCH | 1 | 1 |
| config-get-leading-dot | MATCH | 1 | 1 |
| config-get-cache-files-maxsize | MATCH | 0 | 0 |
| config-get-cache-dir-env | MATCH | 0 | 0 |
| config-get-env-vendor-dir | MATCH | 0 | 0 |
| config-get-home | MATCH | 0 | 0 |
| config-get-data-dir | MATCH | 0 | 0 |
| config-get-cafile-empty | MATCH | 0 | 0 |
| config-get-gitlab-protocol-null | MATCH | 1 | 1 |
| config-get-preferred-install-map | MATCH | 0 | 0 |
| config-get-policy | MATCH | 0 | 0 |
| config-get-policy-nested | MATCH | 0 | 0 |
| config-get-auth-file | MATCH | 0 | 0 |
| config-get-unicode-slashes | MATCH | 0 | 0 |
| config-get-scalar-unicode | MATCH | 0 | 0 |
| config-get-scalar-float | MATCH | 0 | 0 |
| config-list | MATCH | 0 | 0 |
| config-list-source | MATCH | 0 | 0 |
| config-list-plain | MATCH | 0 | 0 |
| config-list-custom | MATCH | 0 | 0 |
| config-list-named-repos | MATCH | 0 | 0 |
| config-list-env | MATCH | 0 | 0 |
| config-list-bad-maxsize | MATCH | 1 | 1 |
| config-list-global | MATCH | 0 | 0 |
| config-list-auth | MATCH | 0 | 0 |
| config-set-sort-packages | MATCH | 0 | 0 |
| config-set-bool-invalid | MATCH | 1 | 1 |
| config-set-process-timeout | MATCH | 0 | 0 |
| config-set-process-timeout-invalid | MATCH | 1 | 1 |
| config-set-too-many | MATCH | 1 | 1 |
| config-set-vendor-dir | MATCH | 0 | 0 |
| config-set-preferred-install | MATCH | 0 | 0 |
| config-set-preferred-install-invalid | MATCH | 1 | 1 |
| config-set-preferred-install-pkg | MATCH | 0 | 0 |
| config-set-preferred-install-pkg-invalid | MATCH | 1 | 1 |
| config-set-store-auths | MATCH | 0 | 0 |
| config-set-store-auths-false | MATCH | 0 | 0 |
| config-set-discard-changes | MATCH | 0 | 0 |
| config-set-bump-after-update | MATCH | 0 | 0 |
| config-set-platform-check | MATCH | 0 | 0 |
| config-set-use-parent-dir | MATCH | 0 | 0 |
| config-set-cache-maxsize | MATCH | 0 | 0 |
| config-set-cache-maxsize-invalid | MATCH | 1 | 1 |
| config-set-bin-compat-invalid | MATCH | 1 | 1 |
| config-set-gitlab-protocol | MATCH | 0 | 0 |
| config-set-autoloader-suffix-null | MATCH | 1 | 1 |
| config-set-cafile-missing | MATCH | 1 | 1 |
| config-set-cafile | MATCH | 0 | 0 |
| config-set-capath-invalid | MATCH | 1 | 1 |
| config-set-github-protocols | MATCH | 0 | 0 |
| config-set-github-protocols-invalid | MATCH | 1 | 1 |
| config-set-github-domains | MATCH | 0 | 0 |
| config-set-disable-tls | MATCH | 0 | 0 |
| config-set-disable-tls-off | MATCH | 0 | 0 |
| config-unset-disable-tls | MATCH | 0 | 0 |
| config-unset-key | MATCH | 0 | 0 |
| config-unset-with-value | MATCH | 1 | 1 |
| config-set-unknown | MATCH | 1 | 1 |
| config-set-allow-plugins-bool | MATCH | 0 | 0 |
| config-set-lock-false | MATCH | 0 | 0 |
| config-set-audit-abandoned | MATCH | 0 | 0 |
| config-set-audit-abandoned-invalid | MATCH | 1 | 1 |
| config-set-audit-block-insecure | MATCH | 0 | 0 |
| config-set-audit-ignore | MATCH | 0 | 0 |
| config-set-audit-ignore-json | MATCH | 0 | 0 |
| config-set-audit-ignore-json-scalar | MATCH | 1 | 1 |
| config-set-audit-ignore-merge | MATCH | 0 | 0 |
| config-set-audit-ignore-merge-mixed | MATCH | 1 | 1 |
| config-set-audit-ignore-severity | MATCH | 0 | 0 |
| config-set-audit-ignore-severity-invalid | MATCH | 1 | 1 |
| config-unset-audit | MATCH | 0 | 0 |
| config-unset-audit-ignore | MATCH | 0 | 0 |
| config-set-policy-block | MATCH | 0 | 0 |
| config-set-policy-audit | MATCH | 0 | 0 |
| config-set-policy-audit-invalid | MATCH | 1 | 1 |
| config-set-policy-block-scope | MATCH | 0 | 0 |
| config-set-policy-block-scope-invalid | MATCH | 1 | 1 |
| config-set-policy-cooldown-period | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-int | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-null | MATCH | 0 | 0 |
| config-set-policy-cooldown-period-invalid | MATCH | 1 | 1 |
| config-set-policy-list-off | MATCH | 0 | 0 |
| config-set-policy-list-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-block | MATCH | 0 | 0 |
| config-set-policy-custom-scope | MATCH | 0 | 0 |
| config-set-policy-custom-scope-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-audit-invalid | MATCH | 1 | 1 |
| config-set-policy-custom-ignore | MATCH | 0 | 0 |
| config-set-policy-reserved | MATCH | 1 | 1 |
| config-set-policy-reserved-name | MATCH | 1 | 1 |
| config-set-policy-reserved-ignore | MATCH | 1 | 1 |
| config-set-policy-sources | MATCH | 1 | 1 |
| config-set-policy-unknown-field | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable | MATCH | 0 | 0 |
| config-set-policy-ignore-unreachable-bad | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable-json | MATCH | 0 | 0 |
| config-set-policy-ignore-unreachable-json-bool | MATCH | 1 | 1 |
| config-set-policy-ignore-unreachable-bool | MATCH | 0 | 0 |
| config-set-policy-ignore-id-json | MATCH | 0 | 0 |
| config-set-policy-ignore-merge | MATCH | 0 | 0 |
| config-set-policy-ignore-merge-mixed | MATCH | 1 | 1 |
| config-set-policy-ignore-severity | MATCH | 0 | 0 |
| config-set-policy-ignore-source | MATCH | 0 | 0 |
| config-unset-policy-nested | MATCH | 0 | 0 |
| config-unset-policy | MATCH | 0 | 0 |
| config-repo-add | MATCH | 0 | 0 |
| config-repo-add-append | MATCH | 0 | 0 |
| config-repo-add-plain | MATCH | 0 | 0 |
| config-repo-disable-packagist | MATCH | 0 | 0 |
| config-repo-add-json | MATCH | 0 | 0 |
| config-repo-add-invalid-json | MATCH | 1 | 1 |
| config-repo-add-true | MATCH | 1 | 1 |
| config-repo-add-three | MATCH | 1 | 1 |
| config-repo-remove | MATCH | 0 | 0 |
| config-repo-remove-assoc | MATCH | 0 | 0 |
| config-repo-get-named | MATCH | 0 | 0 |
| config-extra-set | MATCH | 0 | 0 |
| config-extra-set-new | MATCH | 0 | 0 |
| config-extra-json | MATCH | 0 | 0 |
| config-extra-json-merge-list | MATCH | 0 | 0 |
| config-extra-json-merge-object | MATCH | 0 | 0 |
| config-extra-json-merge-mixed | MATCH | 0 | 0 |
| config-extra-json-invalid | MATCH | 1 | 1 |
| config-extra-unset | MATCH | 0 | 0 |
| config-extra-unset-all | MATCH | 0 | 0 |
| config-suggest-set | MATCH | 0 | 0 |
| config-suggest-unset | MATCH | 0 | 0 |
| config-scripts-set | MATCH | 0 | 0 |
| config-scripts-set-one | MATCH | 0 | 0 |
| config-scripts-unset | MATCH | 0 | 0 |
| config-prop-name | MATCH | 0 | 0 |
| config-prop-description | MATCH | 0 | 0 |
| config-prop-min-stability | MATCH | 0 | 0 |
| config-prop-min-stability-invalid | MATCH | 1 | 1 |
| config-prop-prefer-stable | MATCH | 0 | 0 |
| config-prop-keywords | MATCH | 0 | 0 |
| config-prop-license | MATCH | 0 | 0 |
| config-prop-unset | MATCH | 0 | 0 |
| config-prop-unset-other | MATCH | 0 | 0 |
| config-prop-type-too-many | MATCH | 1 | 1 |
| config-platform-set | MATCH | 0 | 0 |
| config-platform-set-false | MATCH | 0 | 0 |
| config-platform-unset | MATCH | 0 | 0 |
| config-platform-unset-all | MATCH | 0 | 0 |
| config-allow-plugins-set | MATCH | 0 | 0 |
| config-allow-plugins-wildcard | MATCH | 0 | 0 |
| config-allow-plugins-invalid | MATCH | 1 | 1 |
| config-allow-plugins-unset | MATCH | 0 | 0 |
| config-auth-github-local | MATCH | 0 | 0 |
| config-auth-github-global | MATCH | 0 | 0 |
| config-auth-github-too-many | MATCH | 1 | 1 |
| config-auth-http-basic | MATCH | 0 | 0 |
| config-auth-http-basic-missing | MATCH | 1 | 1 |
| config-auth-bitbucket | MATCH | 0 | 0 |
| config-auth-bitbucket-missing | MATCH | 1 | 1 |
| config-auth-gitlab-token-pair | MATCH | 0 | 0 |
| config-auth-gitlab-token-three | MATCH | 1 | 1 |
| config-auth-gitlab-oauth | MATCH | 0 | 0 |
| config-auth-bearer | MATCH | 0 | 0 |
| config-auth-forgejo | MATCH | 0 | 0 |
| config-auth-forgejo-missing | MATCH | 1 | 1 |
| config-auth-custom-headers | MATCH | 0 | 0 |
| config-auth-custom-headers-invalid | MATCH | 1 | 1 |
| config-auth-unset | MATCH | 0 | 0 |
| config-auth-unset-missing-file | MATCH | 0 | 0 |
| config-auth-existing-in-json | MATCH | 0 | 0 |
| config-global-set | MATCH | 0 | 0 |
| config-global-set-modes | MATCH | 0 | 0 |
| config-global-get-existing | MATCH | 0 | 0 |
| config-global-list-existing | MATCH | 0 | 0 |
| config-global-merged-in-local | MATCH | 0 | 0 |
| config-global-dir-absolute | MATCH | 0 | 0 |
| config-global-property | MATCH | 1 | 1 |
| config-global-extra | MATCH | 1 | 1 |
| config-global-repo | MATCH | 0 | 0 |
| config-global-unset | MATCH | 0 | 0 |
| config-global-with-file | MATCH | 1 | 1 |
| config-global-no-json | MATCH | 0 | 0 |
| config-file-get | MATCH | 0 | 0 |
| config-file-set | MATCH | 0 | 0 |
| config-file-missing | MATCH | 1 | 1 |
| config-no-composer-json | MATCH | 1 | 1 |
| config-invalid-json | MATCH | 1 | 1 |
| config-composer-env | MATCH | 0 | 0 |
| config-editor | MATCH | 0 | 0 |
| config-editor-auth | MATCH | 0 | 0 |
| audit-clean | MATCH | 0 | 0 |
| audit-clean-json | MATCH | 0 | 0 |
| audit-clean-summary | MATCH | 0 | 0 |
| audit-no-advisory-repo | MATCH | 0 | 0 |
| audit-table | MATCH | 1 | 1 |
| audit-plain | MATCH | 1 | 1 |
| audit-json | MATCH | 1 | 1 |
| audit-summary | MATCH | 1 | 1 |
| audit-format-short | MATCH | 1 | 1 |
| audit-format-invalid | MATCH | 1 | 1 |
| audit-format-missing-value | MATCH | 1 | 1 |
| audit-verbose | MATCH | 1 | 1 |
| audit-no-vendor | MATCH | 1 | 1 |
| audit-no-vendor-no-dev-only-dev | MATCH | 0 | 0 |
| audit-no-packages | MATCH | 0 | 0 |
| audit-no-lock-no-vendor | MATCH | 1 | 1 |
| audit-locked | MATCH | 1 | 1 |
| audit-locked-json | MATCH | 1 | 1 |
| audit-locked-no-lock | MATCH | 1 | 1 |
| audit-no-dev | MATCH | 0 | 0 |
| audit-with-dev | MATCH | 1 | 1 |
| audit-locked-no-dev | MATCH | 0 | 0 |
| audit-direct | MATCH | 1 | 1 |
| audit-direct-short | MATCH | 1 | 1 |
| audit-abandoned-default | MATCH | 1 | 1 |
| audit-abandoned-plain | MATCH | 1 | 1 |
| audit-abandoned-json | MATCH | 1 | 1 |
| audit-abandoned-summary | MATCH | 1 | 1 |
| audit-abandoned-ignore | MATCH | 0 | 0 |
| audit-abandoned-report | MATCH | 0 | 0 |
| audit-abandoned-fail | MATCH | 1 | 1 |
| audit-abandoned-invalid | MATCH | 1 | 1 |
| audit-abandoned-config-report | MATCH | 0 | 0 |
| audit-abandoned-config-ignore | MATCH | 0 | 0 |
| audit-abandoned-option-overrides-config | MATCH | 1 | 1 |
| audit-abandoned-env-report | MATCH | 0 | 0 |
| audit-abandoned-config-invalid | MATCH | 1 | 1 |
| audit-abandoned-env-overrides-config | MATCH | 1 | 1 |
| audit-abandoned-and-vulnerable | MATCH | 1 | 1 |
| audit-ignore-id | MATCH | 1 | 1 |
| audit-ignore-id-reason | MATCH | 1 | 1 |
| audit-ignore-id-json | MATCH | 1 | 1 |
| audit-ignore-cve | MATCH | 1 | 1 |
| audit-ignore-remote-id | MATCH | 1 | 1 |
| audit-ignore-package | MATCH | 1 | 1 |
| audit-ignore-all | MATCH | 0 | 0 |
| audit-ignore-summary-cve | MATCH | 1 | 1 |
| audit-ignore-apply-block-only | MATCH | 1 | 1 |
| audit-ignore-apply-audit | MATCH | 1 | 1 |
| audit-ignore-severity-option | MATCH | 1 | 1 |
| audit-ignore-severity-option-multi | MATCH | 1 | 1 |
| audit-ignore-severity-option-other | MATCH | 1 | 1 |
| audit-ignore-severity-config | MATCH | 1 | 1 |
| audit-ignore-severity-config-reason | MATCH | 1 | 1 |
| audit-policy-advisories-ignore | MATCH | 1 | 1 |
| audit-policy-disabled | MATCH | 1 | 1 |
| audit-api-table | MATCH | 1 | 1 |
| audit-api-json | MATCH | 1 | 1 |
| audit-api-summary | MATCH | 1 | 1 |
| audit-api-locked | MATCH | 1 | 1 |
| audit-api-unrequested-names | MATCH | 0 | 0 |
| audit-api-ignore-severity | MATCH | 1 | 1 |
| audit-unreachable-404 | MATCH | 100 | 100 |
| audit-unreachable-404-ignored | MATCH | 1 | 1 |
| audit-unreachable-404-ignored-json | MATCH | 1 | 1 |
| audit-unreachable-refused-ignored | MATCH | 0 | 0 |
| audit-unreachable-config-ignored | MATCH | 0 | 0 |
| audit-unreachable-api-ignored | MATCH | 0 | 0 |
| vcs-git-install | DIFF | 0 | 1 |
| vcs-git-type | DIFF | 0 | 1 |
| vcs-file-url | DIFF | 0 | 1 |
| vcs-dev-main | DIFF | 0 | 1 |
| vcs-dev-branch | DIFF | 0 | 1 |
| vcs-dep | DIFF | 0 | 1 |
| vcs-install-from-lock | DIFF | 0 | 1 |
| vcs-dry-run | DIFF | 0 | 1 |
| vcs-no-composer-json | DIFF | 1 | 1 |
| vcs-missing-path | DIFF | 1 | 1 |
| vcs-odd-refs | DIFF | 0 | 1 |
| vcs-odd-refs-vvv | DIFF | 0 | 1 |
| vcs-http-mirror | DIFF | 0 | 1 |

## validate-lax-fail

```
stderr:
+ (go)                                                                      
+ (go)   "./composer.json" does not match the expected JSON schema:         
+ (go)    - require.acme/log : Array value found, but a string is required  
+ (go)                                                                      
- (oracle)                                                                                
- (oracle)   [TypeError]                                                                  
- (oracle)   Composer\Pcre\Preg::isMatch(): Argument #2 ($subject) must be of type strin  
- (oracle)   g, array given, called in /oracle-source/src/Composer/Util/ConfigValidator.  
- (oracle)   php on line 179                                                              
- (oracle)                                                                                
- (oracle) Exception trace:
- (oracle)   at /oracle-source/vendor/composer/pcre/src/Preg.php:289
- (oracle)  Composer\Pcre\Preg::isMatch() at /oracle-source/src/Composer/Util/ConfigValidator.php:179
- (oracle)  Composer\Util\ConfigValidator->validate() at /oracle-source/src/Composer/Command/ValidateCommand.php:89
- (oracle)  Composer\Command\ValidateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## composer-auth-invalid

```
exit code: oracle=255 go=0
stderr:
+ (go) ./composer.json is valid, but with a few warnings
+ (go) See https://getcomposer.org/doc/04-schema.md for details on the schema
+ (go) # General warnings
+ (go) - No license specified, it is recommended to do so. For closed-source software you may use "proprietary" as license.
- (oracle) Fatal error: Uncaught TypeError: Composer\IO\BaseIO::checkAndSetAuthentication(): Argument #2 ($username) must be of type string, int given, called in /oracle-source/src/Composer/IO/BaseIO.php on line 175 and defined in /oracle-source/src/Composer/IO/BaseIO.php:95
- (oracle) Stack trace:
- (oracle) #0 /oracle-source/src/Composer/IO/BaseIO.php(175): Composer\IO\BaseIO->checkAndSetAuthentication('example.org', 5, NULL)
- (oracle) #1 /oracle-source/src/Composer/Factory.php(357): Composer\IO\BaseIO->loadConfiguration(Object(Composer\Config))
- (oracle) #2 /oracle-source/src/Composer/Factory.php(627): Composer\Factory->createComposer(Object(Composer\IO\ConsoleIO), Array, true, '/work/project', true, false)
- (oracle) #3 /oracle-source/src/Composer/Console/Application.php(676): Composer\Factory::create(Object(Composer\IO\ConsoleIO), NULL, true, false)
- (oracle) #4 /oracle-source/src/Composer/Console/Application.php(585): Composer\Console\Application->getComposer(false, true)
- (oracle) #5 /oracle-source/src/Composer/Console/Application.php(532): Composer\Console\Application->hintCommonErrors(Object(ErrorException), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #6 /oracle-source/src/Composer/Console/Application.php(169): Composer\Console\Application->doRunComposer(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput), Object(Composer\IO\ConsoleIO))
- (oracle) #7 /oracle-source/vendor/symfony/console/Application.php(171): Composer\Console\Application->doRun(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #8 /oracle-source/src/Composer/Console/Application.php(142): Symfony\Component\Console\Application->run(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #9 /oracle-source/bin/composer(113): Composer\Console\Application->run()
- (oracle) #10 {main}
- (oracle)   thrown in /oracle-source/src/Composer/IO/BaseIO.php on line 95
```

## vcs-git-install

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (v1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.1.0): Cloning 7021843339
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "4e4b60712a57cfd50a089ebc4ffc3b4f",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 0923b38381f55bcd
- (oracle) F vendor/autoload.php 644 42f2f59dbd47fdd3
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 0565f15b17b8983f
- (oracle) F vendor/composer/autoload_static.php 644 2a4c2f75d89a8e81
- (oracle) F vendor/composer/installed.json 644 5c52719142765dff
- (oracle) F vendor/composer/installed.php 644 9fc3651dd661299a
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 0c67c498c8048495
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib 'v1.1.0' '702184333938c0ff49e87f0af7d9931d71486a52' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 702184333938c0ff49e87f0af7d9931d71486a52
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit4e4b60712a57cfd50a089ebc4ffc3b4f
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit4e4b60712a57cfd50a089ebc4ffc3b4f', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit4e4b60712a57cfd50a089ebc4ffc3b4f', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit4e4b60712a57cfd50a089ebc4ffc3b4f::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit4e4b60712a57cfd50a089ebc4ffc3b4f
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit4e4b60712a57cfd50a089ebc4ffc3b4f::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit4e4b60712a57cfd50a089ebc4ffc3b4f::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit4e4b60712a57cfd50a089ebc4ffc3b4f::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '702184333938c0ff49e87f0af7d9931d71486a52',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-git-type

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "git" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (v1.0.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.0.0): Cloning b2b1b720ad
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "6ff3631b94a6facc4da4bac7f3cb58d2",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "b2b1b720ad1e4e06cb3d641ea573be406bab1502"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T04:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 6fe7ca2fc48f92f9
- (oracle) F vendor/autoload.php 644 f5ed68943e5c1128
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 75e6825fa15d9aec
- (oracle) F vendor/composer/autoload_static.php 644 e3bffa73370a049d
- (oracle) F vendor/composer/installed.json 644 814b45f458ede728
- (oracle) F vendor/composer/installed.php 644 9e0a07571d264a72
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 c89d86f22ca56c90
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib 'v1.0.0' 'b2b1b720ad1e4e06cb3d641ea573be406bab1502' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD b2b1b720ad1e4e06cb3d641ea573be406bab1502
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.0.0)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit6ff3631b94a6facc4da4bac7f3cb58d2
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit6ff3631b94a6facc4da4bac7f3cb58d2', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit6ff3631b94a6facc4da4bac7f3cb58d2', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit6ff3631b94a6facc4da4bac7f3cb58d2::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit6ff3631b94a6facc4da4bac7f3cb58d2
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit6ff3631b94a6facc4da4bac7f3cb58d2::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit6ff3631b94a6facc4da4bac7f3cb58d2::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit6ff3631b94a6facc4da4bac7f3cb58d2::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "b2b1b720ad1e4e06cb3d641ea573be406bab1502"
- (oracle)             },
- (oracle)             "time": "2024-01-02T04:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'b2b1b720ad1e4e06cb3d641ea573be406bab1502',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-file-url

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (v1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.1.0): Cloning 7021843339
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "bda14524923cb05433e34845e099cc83",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 ae9da9a093b681aa
- (oracle) F vendor/autoload.php 644 e5ca23dd31b54815
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 abae6f2bd4b76c86
- (oracle) F vendor/composer/autoload_static.php 644 ed13c5ef2e1ad88d
- (oracle) F vendor/composer/installed.json 644 5c52719142765dff
- (oracle) F vendor/composer/installed.php 644 9fc3651dd661299a
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 0c67c498c8048495
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib 'v1.1.0' '702184333938c0ff49e87f0af7d9931d71486a52' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 702184333938c0ff49e87f0af7d9931d71486a52
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInitbda14524923cb05433e34845e099cc83
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInitbda14524923cb05433e34845e099cc83', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInitbda14524923cb05433e34845e099cc83', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInitbda14524923cb05433e34845e099cc83::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInitbda14524923cb05433e34845e099cc83
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInitbda14524923cb05433e34845e099cc83::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInitbda14524923cb05433e34845e099cc83::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInitbda14524923cb05433e34845e099cc83::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '702184333938c0ff49e87f0af7d9931d71486a52',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-dev-main

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (dev-main ba03332)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (dev-main ba03332): Cloning ba0333212d
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "6f0241d6931978f258899ba7478cd759",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "dev-main",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "ba0333212d6fb4b08718b30b9dc442e929393f8b"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "php": ">=7.4"
- (oracle)             },
- (oracle)             "default-branch": true,
- (oracle)             "type": "library",
- (oracle)             "extra": {
- (oracle)                 "branch-alias": {
- (oracle)                     "dev-main": "2.x-dev"
- (oracle)                 }
- (oracle)             },
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T10:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {
- (oracle)         "vcs/gitlib": 20
- (oracle)     },
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 c4c270fd5dd76f23
- (oracle) F vendor/autoload.php 644 31c0358b46ca6b63
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 56c494ebf2b0f57d
- (oracle) F vendor/composer/autoload_static.php 644 2692875d270326b0
- (oracle) F vendor/composer/installed.json 644 fda523144b83d4ee
- (oracle) F vendor/composer/installed.php 644 5b6b13b4c4031591
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
- (oracle) F vendor/vcs/gitlib/composer.json 644 24ab022a3246bf6c
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 c398df641a8d13d6
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib 'dev-main' 'ba0333212d6fb4b08718b30b9dc442e929393f8b' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD ba0333212d6fb4b08718b30b9dc442e929393f8b
- (oracle) ref refs/heads/main
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit6f0241d6931978f258899ba7478cd759
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         require __DIR__ . '/platform_check.php';
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit6f0241d6931978f258899ba7478cd759', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit6f0241d6931978f258899ba7478cd759', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit6f0241d6931978f258899ba7478cd759::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit6f0241d6931978f258899ba7478cd759
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit6f0241d6931978f258899ba7478cd759::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit6f0241d6931978f258899ba7478cd759::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit6f0241d6931978f258899ba7478cd759::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "dev-main",
- (oracle)             "version_normalized": "dev-main",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "ba0333212d6fb4b08718b30b9dc442e929393f8b"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "php": ">=7.4"
- (oracle)             },
- (oracle)             "time": "2024-01-02T10:04:05+00:00",
- (oracle)             "default-branch": true,
- (oracle)             "type": "library",
- (oracle)             "extra": {
- (oracle)                 "branch-alias": {
- (oracle)                     "dev-main": "2.x-dev"
- (oracle)                 }
- (oracle)             },
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'dev-main',
- (oracle)             'version' => 'dev-main',
- (oracle)             'reference' => 'ba0333212d6fb4b08718b30b9dc442e929393f8b',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(
- (oracle)                 0 => '2.x-dev',
- (oracle)             ),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
vendor/composer/platform_check.php:
- (oracle) <?php
- (oracle) // platform_check.php @generated by Composer
- (oracle) 
- (oracle) $issues = array();
- (oracle) 
- (oracle) if (!(PHP_VERSION_ID >= 70400)) {
- (oracle)     $issues[] = 'Your Composer dependencies require a PHP version ">= 7.4.0". You are running ' . PHP_VERSION . '.';
- (oracle) }
- (oracle) 
- (oracle) if ($issues) {
- (oracle)     if (!headers_sent()) {
- (oracle)         header('HTTP/1.1 500 Internal Server Error');
- (oracle)     }
- (oracle)     if (!ini_get('display_errors')) {
- (oracle)         if (PHP_SAPI === 'cli' || PHP_SAPI === 'phpdbg') {
- (oracle)             fwrite(STDERR, 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . implode(PHP_EOL, $issues) . PHP_EOL.PHP_EOL);
- (oracle)         } elseif (!headers_sent()) {
- (oracle)             echo 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . str_replace('You are running '.PHP_VERSION.'.', '', implode(PHP_EOL, $issues)) . PHP_EOL.PHP_EOL;
- (oracle)         }
- (oracle)     }
- (oracle)     throw new \RuntimeException(
- (oracle)         'Composer detected issues in your platform: ' . implode(' ', $issues)
- (oracle)     );
- (oracle) }
- (oracle) 
```

## vcs-dev-branch

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (1.x-dev 24bec3f)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (1.x-dev 24bec3f): Cloning 24bec3fa0d
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "7ebd767fe6cdf442cb1715a889cfcf40",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "1.x-dev",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "24bec3fa0d593ddedbd07cb32c28449d25b4f5a0"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T09:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {
- (oracle)         "vcs/gitlib": 20
- (oracle)     },
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 1626bdeeb2a1cd36
- (oracle) F vendor/autoload.php 644 01a72c5c4ea6b021
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 39676327ac126e4b
- (oracle) F vendor/composer/autoload_static.php 644 db7fddaeaceeedbc
- (oracle) F vendor/composer/installed.json 644 e0981106f4939f63
- (oracle) F vendor/composer/installed.php 644 223b5db582dc0cbf
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 351c96f263c90b60
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib '1.x-dev' '24bec3fa0d593ddedbd07cb32c28449d25b4f5a0' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 24bec3fa0d593ddedbd07cb32c28449d25b4f5a0
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at 24bec3f)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit7ebd767fe6cdf442cb1715a889cfcf40
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit7ebd767fe6cdf442cb1715a889cfcf40', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit7ebd767fe6cdf442cb1715a889cfcf40', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit7ebd767fe6cdf442cb1715a889cfcf40::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit7ebd767fe6cdf442cb1715a889cfcf40
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit7ebd767fe6cdf442cb1715a889cfcf40::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit7ebd767fe6cdf442cb1715a889cfcf40::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit7ebd767fe6cdf442cb1715a889cfcf40::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "1.x-dev",
- (oracle)             "version_normalized": "1.9999999.9999999.9999999-dev",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "24bec3fa0d593ddedbd07cb32c28449d25b4f5a0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T09:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => '1.x-dev',
- (oracle)             'version' => '1.9999999.9999999.9999999-dev',
- (oracle)             'reference' => '24bec3fa0d593ddedbd07cb32c28449d25b4f5a0',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-dep

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitdep (1.0.0)
- (oracle)   - Locking vcs/gitlib (v1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.1.0): Cloning 7021843339
- (oracle)   - Installing vcs/gitdep (1.0.0): Cloning c16e018238
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "bfc6217487d11fba7def9afd37c65fcd",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitdep",
- (oracle)             "version": "1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitdep.git",
- (oracle)                 "reference": "c16e0182384ca46e31fb08304ab23bc789f9e70a"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "vcs/gitlib": "^1.0"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "classmap": [
- (oracle)                     "lib/"
- (oracle)                 ]
- (oracle)             },
- (oracle)             "description": "depends on gitlib",
- (oracle)             "time": "2024-01-02T11:04:05+00:00"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitdep
- (oracle) D vendor/vcs/gitdep/.git
- (oracle) D vendor/vcs/gitdep/lib
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 dca0d00e0f1a9a8f
- (oracle) F vendor/autoload.php 644 5860f1fddbcae241
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 1a67220a20f47f86
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 f9e1f3494eb10f40
- (oracle) F vendor/composer/autoload_static.php 644 59563898451382f9
- (oracle) F vendor/composer/installed.json 644 86c2031c73c816ed
- (oracle) F vendor/composer/installed.php 644 3c42d922cf93d54d
- (oracle) F vendor/vcs/gitdep/composer.json 644 790efae931d3b928
- (oracle) F vendor/vcs/gitdep/lib/Dep.php 644 3dd3ec586a79dc23
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 0c67c498c8048495
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsDep: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitdep '1.0.0' 'c16e0182384ca46e31fb08304ab23bc789f9e70a' dev=true
- (oracle) installed: vcs/gitlib 'v1.1.0' '702184333938c0ff49e87f0af7d9931d71486a52' dev=true
git state:
- (oracle) == ./vendor/vcs/gitdep
- (oracle) HEAD c16e0182384ca46e31fb08304ab23bc789f9e70a
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitdep.git (fetch)
- (oracle) composer	/fixtures/git/gitdep.git (push)
- (oracle) origin	/fixtures/git/gitdep.git (fetch)
- (oracle) origin	/fixtures/git/gitdep.git (push)
- (oracle) * (HEAD detached at 1.0.0)
- (oracle)   main
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 702184333938c0ff49e87f0af7d9931d71486a52
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle)     'VcsDep' => $vendorDir . '/vcs/gitdep/lib/Dep.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)         'VcsDep' => __DIR__ . '/..' . '/vcs/gitdep/lib/Dep.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitdep",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitdep.git",
- (oracle)                 "reference": "c16e0182384ca46e31fb08304ab23bc789f9e70a"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "vcs/gitlib": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T11:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "classmap": [
- (oracle)                     "lib/"
- (oracle)                 ]
- (oracle)             },
- (oracle)             "description": "depends on gitlib",
- (oracle)             "install-path": "../vcs/gitdep"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitdep' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'c16e0182384ca46e31fb08304ab23bc789f9e70a',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitdep',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '702184333938c0ff49e87f0af7d9931d71486a52',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-install-from-lock

```
exit code: oracle=0 go=1
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Package operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.1.0): Cloning 7021843339
- (oracle)   - Installing vcs/gitdep (1.0.0): Cloning c16e018238
- (oracle) Generating autoload files
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitdep
- (oracle) D vendor/vcs/gitdep/.git
- (oracle) D vendor/vcs/gitdep/lib
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F vendor/autoload.php 644 5860f1fddbcae241
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 1a67220a20f47f86
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 f9e1f3494eb10f40
- (oracle) F vendor/composer/autoload_static.php 644 59563898451382f9
- (oracle) F vendor/composer/installed.json 644 86c2031c73c816ed
- (oracle) F vendor/composer/installed.php 644 3c42d922cf93d54d
- (oracle) F vendor/vcs/gitdep/composer.json 644 790efae931d3b928
- (oracle) F vendor/vcs/gitdep/lib/Dep.php 644 3dd3ec586a79dc23
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 0c67c498c8048495
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsDep: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitdep '1.0.0' 'c16e0182384ca46e31fb08304ab23bc789f9e70a' dev=true
- (oracle) installed: vcs/gitlib 'v1.1.0' '702184333938c0ff49e87f0af7d9931d71486a52' dev=true
git state:
- (oracle) == ./vendor/vcs/gitdep
- (oracle) HEAD c16e0182384ca46e31fb08304ab23bc789f9e70a
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitdep.git (fetch)
- (oracle) composer	/fixtures/git/gitdep.git (push)
- (oracle) origin	/fixtures/git/gitdep.git (fetch)
- (oracle) origin	/fixtures/git/gitdep.git (push)
- (oracle) * (HEAD detached at 1.0.0)
- (oracle)   main
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 702184333938c0ff49e87f0af7d9931d71486a52
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitlib.git (fetch)
- (oracle) composer	/fixtures/git/gitlib.git (push)
- (oracle) origin	/fixtures/git/gitlib.git (fetch)
- (oracle) origin	/fixtures/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/1.x
- (oracle)   remotes/composer/main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle)     'VcsDep' => $vendorDir . '/vcs/gitdep/lib/Dep.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInitbfc6217487d11fba7def9afd37c65fcd', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)         'VcsDep' => __DIR__ . '/..' . '/vcs/gitdep/lib/Dep.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInitbfc6217487d11fba7def9afd37c65fcd::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitdep",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitdep.git",
- (oracle)                 "reference": "c16e0182384ca46e31fb08304ab23bc789f9e70a"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "vcs/gitlib": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T11:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "classmap": [
- (oracle)                     "lib/"
- (oracle)                 ]
- (oracle)             },
- (oracle)             "description": "depends on gitlib",
- (oracle)             "install-path": "../vcs/gitdep"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitdep' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'c16e0182384ca46e31fb08304ab23bc789f9e70a',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitdep',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '702184333938c0ff49e87f0af7d9931d71486a52',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-dry-run

```
exit code: oracle=0 go=1
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitdep (1.0.0)
- (oracle)   - Locking vcs/gitlib (v1.1.0)
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitlib (v1.1.0)
- (oracle)   - Installing vcs/gitdep (1.0.0)
- (oracle) No installed packages - skipping audit.
```

## vcs-no-composer-json

```
stderr:
- (oracle) Loading composer repositories with package information
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
- (oracle)   No valid composer.json was found in any branch or tag of /fixtures/git/gite  
- (oracle)   mpty.git, could not load a package from it.                                  
```

## vcs-missing-path

```
stderr:
- (oracle) Loading composer repositories with package information
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
- (oracle)   Failed to read package information from /fixtures/git/nope.git as the path   
- (oracle)   does not exist                                                               
```

## vcs-odd-refs

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
- (oracle) Loading composer repositories with package information
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) Skipped branch broken, "7cfde7d4abc9c4b59358ea94bbbc654ec84968ff:composer.json" does not contain valid JSON
- (oracle) Parse error on line 1:
- (oracle) ...ame": "vcs/gitodd",
- (oracle) ---------------------^
- (oracle) Expected: 'STRING' - It appears you have an extra trailing comma
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitodd (1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Installing vcs/gitodd (1.1.0): Cloning f2374cc590
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "81d19246bdae8fe3bfdf249e44212bdd",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitodd",
- (oracle)             "version": "1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitodd.git",
- (oracle)                 "reference": "f2374cc590ddd630879007769852ad69dbe909ce"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitodd
- (oracle) D vendor/vcs/gitodd/.git
- (oracle) F composer.lock 644 70aa700e6204a13d
- (oracle) F vendor/autoload.php 644 84633b7a2e33cc18
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 c35c405f584d1f28
- (oracle) F vendor/composer/autoload_real.php 644 2858cb0a1dad48fb
- (oracle) F vendor/composer/autoload_static.php 644 79636c6aa0c5b10c
- (oracle) F vendor/composer/installed.json 644 1f23045c31ea3815
- (oracle) F vendor/composer/installed.php 644 0d7641fa4e1b8e91
- (oracle) F vendor/vcs/gitodd/composer.json 644 fe0acac8b19576ef
git state:
- (oracle) == ./vendor/vcs/gitodd
- (oracle) HEAD f2374cc590ddd630879007769852ad69dbe909ce
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitodd.git (fetch)
- (oracle) composer	/fixtures/git/gitodd.git (push)
- (oracle) origin	/fixtures/git/gitodd.git (fetch)
- (oracle) origin	/fixtures/git/gitodd.git (push)
- (oracle) * (HEAD detached at 1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/broken
- (oracle)   remotes/composer/main
- (oracle)   remotes/composer/nojson
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/broken
- (oracle)   remotes/origin/main
- (oracle)   remotes/origin/nojson
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitodd",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitodd.git",
- (oracle)                 "reference": "f2374cc590ddd630879007769852ad69dbe909ce"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "install-path": "../vcs/gitodd"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitodd' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => 'f2374cc590ddd630879007769852ad69dbe909ce',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitodd',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-odd-refs-vvv

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) Loading composer repositories with package information
- (oracle) Reading composer.json of vcs/gitodd (1.0.0)
- (oracle) Skipped tag 1.0.0, tag (1.0.0.0) does not match version (3.0.0.0) in composer.json
- (oracle) Reading composer.json of vcs/gitodd (1.1.0)
- (oracle) Importing tag 1.1.0 (1.1.0.0)
- (oracle) Reading composer.json of vcs/gitodd (1.2.0-dev)
- (oracle) Skipped tag 1.2.0-dev, invalid tag name, tags can not use dev prefixes or suffixes
- (oracle) Reading composer.json of vcs/gitodd (main)
- (oracle) Importing branch main (dev-main)
- (oracle) Reading composer.json of vcs/gitodd (broken)
- (oracle) Skipped branch broken, "7cfde7d4abc9c4b59358ea94bbbc654ec84968ff:composer.json" does not contain valid JSON
- (oracle) Parse error on line 1:
- (oracle) ...ame": "vcs/gitodd",
- (oracle) ---------------------^
- (oracle) Expected: 'STRING' - It appears you have an extra trailing comma
- (oracle) Reading composer.json of vcs/gitodd (nojson)
- (oracle) Skipped branch nojson, no composer file
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Dependency resolution completed in N seconds
- (oracle) Analyzed 58 packages to resolve dependencies
- (oracle) Analyzed 58 rules to resolve dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle) Installs: vcs/gitodd:1.1.0
- (oracle)   - Locking vcs/gitodd (1.1.0) from vcs repo (git /fixtures/git/gitodd.git)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle) Installs: vcs/gitodd:1.1.0
- (oracle)   - Installing vcs/gitodd (1.1.0): Cloning f2374cc590ddd630879007769852ad69dbe909ce
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "81d19246bdae8fe3bfdf249e44212bdd",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitodd",
- (oracle)             "version": "1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitodd.git",
- (oracle)                 "reference": "f2374cc590ddd630879007769852ad69dbe909ce"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitodd
- (oracle) D vendor/vcs/gitodd/.git
- (oracle) F composer.lock 644 70aa700e6204a13d
- (oracle) F vendor/autoload.php 644 84633b7a2e33cc18
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 c35c405f584d1f28
- (oracle) F vendor/composer/autoload_real.php 644 2858cb0a1dad48fb
- (oracle) F vendor/composer/autoload_static.php 644 79636c6aa0c5b10c
- (oracle) F vendor/composer/installed.json 644 1f23045c31ea3815
- (oracle) F vendor/composer/installed.php 644 0d7641fa4e1b8e91
- (oracle) F vendor/vcs/gitodd/composer.json 644 fe0acac8b19576ef
git state:
- (oracle) == ./vendor/vcs/gitodd
- (oracle) HEAD f2374cc590ddd630879007769852ad69dbe909ce
- (oracle) ref 
- (oracle) composer	/fixtures/git/gitodd.git (fetch)
- (oracle) composer	/fixtures/git/gitodd.git (push)
- (oracle) origin	/fixtures/git/gitodd.git (fetch)
- (oracle) origin	/fixtures/git/gitodd.git (push)
- (oracle) * (HEAD detached at 1.1.0)
- (oracle)   main
- (oracle)   remotes/composer/broken
- (oracle)   remotes/composer/main
- (oracle)   remotes/composer/nojson
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/broken
- (oracle)   remotes/origin/main
- (oracle)   remotes/origin/nojson
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit81d19246bdae8fe3bfdf249e44212bdd', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit81d19246bdae8fe3bfdf249e44212bdd::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitodd",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "/fixtures/git/gitodd.git",
- (oracle)                 "reference": "f2374cc590ddd630879007769852ad69dbe909ce"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "install-path": "../vcs/gitodd"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitodd' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => 'f2374cc590ddd630879007769852ad69dbe909ce',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitodd',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```

## vcs-http-mirror

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Repository type "vcs" is not supported by gocomposer yet (supported: compos  
+ (go)   er, package, path)                                                           
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 1 install, 0 updates, 0 removals
- (oracle)   - Locking vcs/gitlib (v1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 1 install, 0 updates, 0 removals
- (oracle)   - Syncing vcs/gitlib (v1.1.0) into cache
- (oracle)   - Installing vcs/gitlib (v1.1.0): Cloning 7021843339 from cache
- (oracle) Generating autoload files
composer.lock:
- (oracle) {
- (oracle)     "_readme": [
- (oracle)         "This file locks the dependencies of your project to a known state",
- (oracle)         "Read more about it at https://getcomposer.org/doc/01-basic-usage.md#installing-dependencies",
- (oracle)         "This file is @generated automatically"
- (oracle)     ],
- (oracle)     "content-hash": "e3db336dc4b3a3443a893eecbdc2b32e",
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "http://127.0.0.1:8899/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "type": "library",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "time": "2024-01-02T06:04:05+00:00"
- (oracle)         }
- (oracle)     ],
- (oracle)     "packages-dev": [],
- (oracle)     "aliases": [],
- (oracle)     "minimum-stability": "stable",
- (oracle)     "stability-flags": {},
- (oracle)     "prefer-stable": false,
- (oracle)     "prefer-lowest": false,
- (oracle)     "platform": {},
- (oracle)     "platform-dev": {},
- (oracle)     "plugin-api-version": "2.11.0"
- (oracle) }
project tree:
- (oracle) D vendor
- (oracle) D vendor/composer
- (oracle) D vendor/vcs
- (oracle) D vendor/vcs/gitlib
- (oracle) D vendor/vcs/gitlib/.git
- (oracle) D vendor/vcs/gitlib/src
- (oracle) F composer.lock 644 e49a323e2c917ccc
- (oracle) F vendor/autoload.php 644 aa06d43855cb1f5a
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 ed8130bc1b5d18c0
- (oracle) F vendor/composer/autoload_real.php 644 0e5f28360cddf0f5
- (oracle) F vendor/composer/autoload_static.php 644 d1b3285321dad3ef
- (oracle) F vendor/composer/installed.json 644 09f645b5c2b548fc
- (oracle) F vendor/composer/installed.php 644 9fc3651dd661299a
- (oracle) F vendor/vcs/gitlib/composer.json 644 07b24d8c59029ef6
- (oracle) F vendor/vcs/gitlib/src/Lib.php 644 0c67c498c8048495
autoload probe:
- (oracle) class App\Kernel: autoloaded
- (oracle) class VcsLib\Lib: autoloaded
- (oracle) function acme_util_hello: missing
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) installed: vcs/gitlib 'v1.1.0' '702184333938c0ff49e87f0af7d9931d71486a52' dev=true
git state:
- (oracle) == ./vendor/vcs/gitlib
- (oracle) HEAD 702184333938c0ff49e87f0af7d9931d71486a52
- (oracle) ref 
- (oracle) composer	http://127.0.0.1:8899/git/gitlib.git (fetch)
- (oracle) composer	http://127.0.0.1:8899/git/gitlib.git (push)
- (oracle) origin	http://127.0.0.1:8899/git/gitlib.git (fetch)
- (oracle) origin	http://127.0.0.1:8899/git/gitlib.git (push)
- (oracle) * (HEAD detached at v1.1.0)
- (oracle)   main
- (oracle)   remotes/origin/1.x
- (oracle)   remotes/origin/HEAD -> origin/main
- (oracle)   remotes/origin/main
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_psr4.php:
- (oracle) <?php
- (oracle) // autoload_psr4.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'VcsLib\\' => array($vendorDir . '/vcs/gitlib/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInite3db336dc4b3a3443a893eecbdc2b32e
- (oracle) {
- (oracle)     private static $loader;
- (oracle) 
- (oracle)     public static function loadClassLoader($class)
- (oracle)     {
- (oracle)         if ('Composer\Autoload\ClassLoader' === $class) {
- (oracle)             require __DIR__ . '/ClassLoader.php';
- (oracle)         }
- (oracle)     }
- (oracle) 
- (oracle)     /**
- (oracle)      * @return \Composer\Autoload\ClassLoader
- (oracle)      */
- (oracle)     public static function getLoader()
- (oracle)     {
- (oracle)         if (null !== self::$loader) {
- (oracle)             return self::$loader;
- (oracle)         }
- (oracle) 
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInite3db336dc4b3a3443a893eecbdc2b32e', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInite3db336dc4b3a3443a893eecbdc2b32e', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInite3db336dc4b3a3443a893eecbdc2b32e::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         return $loader;
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/autoload_static.php:
- (oracle) <?php
- (oracle) // autoload_static.php @generated by Composer
- (oracle) 
- (oracle) namespace Composer\Autoload;
- (oracle) 
- (oracle) class ComposerStaticInite3db336dc4b3a3443a893eecbdc2b32e
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'V' =>
- (oracle)         array (
- (oracle)             'VcsLib\\' => 7,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'VcsLib\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/vcs/gitlib/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInite3db336dc4b3a3443a893eecbdc2b32e::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInite3db336dc4b3a3443a893eecbdc2b32e::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInite3db336dc4b3a3443a893eecbdc2b32e::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "vcs/gitlib",
- (oracle)             "version": "v1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "source": {
- (oracle)                 "type": "git",
- (oracle)                 "url": "http://127.0.0.1:8899/git/gitlib.git",
- (oracle)                 "reference": "702184333938c0ff49e87f0af7d9931d71486a52"
- (oracle)             },
- (oracle)             "time": "2024-01-02T06:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "source",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "VcsLib\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "git fixture library",
- (oracle)             "install-path": "../vcs/gitlib"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": []
- (oracle) }
vendor/composer/installed.php:
- (oracle) <?php return array(
- (oracle)     'root' => array(
- (oracle)         'name' => 'test/project',
- (oracle)         'pretty_version' => '1.0.0+no-version-set',
- (oracle)         'version' => '1.0.0.0',
- (oracle)         'reference' => null,
- (oracle)         'type' => 'project',
- (oracle)         'install_path' => __DIR__ . '/../../',
- (oracle)         'aliases' => array(),
- (oracle)         'dev' => true,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'vcs/gitlib' => array(
- (oracle)             'pretty_version' => 'v1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '702184333938c0ff49e87f0af7d9931d71486a52',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../vcs/gitlib',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)     ),
- (oracle) );
```
