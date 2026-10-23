# Composer vs gocomposer comparison

- gocomposer_sha256: `0e9a127d41ed42c0d5034ea3600b2eae0bd3739b11c01200b3a9cf833f9bca9a`
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
| bump-basic | MATCH | 0 | 0 |
| bump-dry-run | MATCH | 1 | 1 |
| bump-dev-only | MATCH | 0 | 0 |
| bump-no-dev-only | MATCH | 0 | 0 |
| bump-library | MATCH | 0 | 0 |
| bump-library-dev-only | MATCH | 0 | 0 |
| bump-untyped | MATCH | 0 | 0 |
| bump-filter | MATCH | 0 | 0 |
| bump-filter-wildcard | MATCH | 0 | 0 |
| bump-nothing | MATCH | 0 | 0 |
| bump-nothing-dry-run | MATCH | 0 | 0 |
| bump-lock-outdated | MATCH | 2 | 2 |
| bump-lock-disabled | KNOWN-DIFF | 1 | 1 |
| bump-installed-only | MATCH | 0 | 0 |
| bump-not-installed | MATCH | 0 | 0 |
| bump-branch-alias | MATCH | 0 | 0 |
| bump-no-json | MATCH | 1 | 1 |
| bump-after-update-flag | MATCH | 0 | 0 |
| bump-after-update-dev | MATCH | 0 | 0 |
| bump-after-update-no-dev | MATCH | 0 | 0 |
| bump-after-update-dry-run | MATCH | 1 | 1 |
| bump-after-update-config | MATCH | 0 | 0 |
| bump-after-update-config-dev | MATCH | 0 | 0 |
| bump-after-update-untyped | MATCH | 0 | 0 |
| bump-after-update-partial | MATCH | 0 | 0 |
| bump-after-update-nothing | MATCH | 0 | 0 |
| bump-after-update-lock | MATCH | 0 | 0 |
| init-full | MATCH | 0 | 0 |
| init-full-reordered | MATCH | 0 | 0 |
| init-require | MATCH | 0 | 0 |
| init-require-pair-args | MATCH | 0 | 0 |
| init-require-no-version | MATCH | 1 | 1 |
| init-require-dev-no-version | MATCH | 1 | 1 |
| init-autoload | MATCH | 0 | 0 |
| init-autoload-nested | MATCH | 0 | 0 |
| init-autoload-with-deps | MATCH | 0 | 0 |
| init-autoload-dev-deps | MATCH | 0 | 0 |
| init-git-defaults | MATCH | 0 | 0 |
| init-git-defaults-description | MATCH | 0 | 0 |
| init-env-defaults | MATCH | 0 | 0 |
| init-env-overrides-git | MATCH | 0 | 0 |
| init-user-env | MATCH | 0 | 0 |
| init-no-defaults | MATCH | 0 | 0 |
| init-camel-dir | MATCH | 0 | 0 |
| init-invalid-name | MATCH | 1 | 1 |
| init-invalid-name-noslash | MATCH | 1 | 1 |
| init-invalid-author | MATCH | 1 | 1 |
| init-invalid-email | MATCH | 1 | 1 |
| init-empty-author | MATCH | 1 | 1 |
| init-author-name-only | MATCH | 0 | 0 |
| init-existing-json | MATCH | 0 | 0 |
| init-existing-git | MATCH | 0 | 0 |
| init-repository | MATCH | 0 | 0 |
| init-repository-invalid | MATCH | 1 | 1 |
| init-repository-bad-json | MATCH | 1 | 1 |
| init-repository-composer-file | MATCH | 0 | 0 |
| init-repository-filesystem-file | MATCH | 1 | 1 |
| init-repository-missing-file | MATCH | 1 | 1 |
| init-bad-stability | MATCH | 1 | 1 |
| init-bad-type-schema | MATCH | 1 | 1 |
| init-bad-license | MATCH | 0 | 0 |
| init-composer-env | MATCH | 0 | 0 |
| init-composer-env-subdir | MATCH | 0 | 0 |
| init-unicode | MATCH | 0 | 0 |
| global-require | MATCH | 0 | 0 |
| global-update | MATCH | 0 | 0 |
| global-show | MATCH | 0 | 0 |
| global-show-abbrev | MATCH | 0 | 0 |
| global-options-first | MATCH | 0 | 0 |
| global-remove | MATCH | 0 | 0 |
| global-no-args | MATCH | 1 | 1 |
| global-unknown-option-only | MATCH | 1 | 1 |
| global-unknown-command | MATCH | 1 | 1 |
| global-config-quoted | MATCH | 0 | 0 |
| global-config-equals | MATCH | 0 | 0 |
| global-new-home | MATCH | 0 | 0 |
| global-composer-env | MATCH | 0 | 0 |
| global-validate | MATCH | 2 | 2 |
| global-dump-autoload | MATCH | 0 | 0 |
| archive-root-tar | MATCH | 0 | 0 |
| archive-root-zip | MATCH | 0 | 0 |
| archive-root-tar-gz | MATCH | 0 | 0 |
| archive-root-tar-bz2 | MATCH | 1 | 1 |
| archive-root-bad-format | MATCH | 1 | 1 |
| archive-root-empty-format | MATCH | 1 | 1 |
| archive-root-dir | MATCH | 0 | 0 |
| archive-root-dir-relative | MATCH | 0 | 0 |
| archive-root-file | MATCH | 0 | 0 |
| archive-root-file-subdir | MATCH | 0 | 0 |
| archive-root-ignore-filters | MATCH | 0 | 0 |
| archive-root-named | MATCH | 0 | 0 |
| archive-root-installed | MATCH | 0 | 0 |
| archive-root-existing-archives | MATCH | 0 | 0 |
| archive-root-config-format | MATCH | 0 | 0 |
| archive-root-dir-is-file | MATCH | 1 | 1 |
| archive-root-scripts | MATCH | 0 | 0 |
| archive-package | MATCH | 0 | 0 |
| archive-package-installed | MATCH | 0 | 0 |
| archive-package-constraint | MATCH | 0 | 0 |
| archive-package-stability | MATCH | 0 | 0 |
| archive-package-stability-filtered | MATCH | 1 | 1 |
| archive-package-missing | MATCH | 1 | 1 |
| archive-package-file | MATCH | 1 | 1 |
| archive-package-case | MATCH | 0 | 0 |
| archive-no-json | MATCH | 1 | 1 |
| archive-empty-tar | MATCH | 0 | 0 |
| archive-empty-zip | MATCH | 0 | 0 |
| archive-empty-tar-gz | MATCH | 0 | 0 |
| archive-only-empty-dirs | MATCH | 0 | 0 |
| archive-globstar-attributes | MATCH | 1 | 1 |
| archive-globstar-valid | MATCH | 0 | 0 |
| archive-vendor-dir | MATCH | 0 | 0 |
| reinstall-basic | MATCH | 0 | 0 |
| reinstall-progress | MATCH | 0 | 0 |
| reinstall-verbose | MATCH | 0 | 0 |
| reinstall-wildcard | MATCH | 0 | 0 |
| reinstall-all-wildcard | MATCH | 0 | 0 |
| reinstall-type | MATCH | 0 | 0 |
| reinstall-type-none | MATCH | 1 | 1 |
| reinstall-type-and-names | MATCH | 1 | 1 |
| reinstall-no-args | MATCH | 1 | 1 |
| reinstall-no-match | MATCH | 1 | 1 |
| reinstall-partial-match | MATCH | 0 | 0 |
| reinstall-not-installed | MATCH | 1 | 1 |
| reinstall-lock-only | MATCH | 1 | 1 |
| reinstall-no-autoloader | MATCH | 0 | 0 |
| reinstall-optimize | MATCH | 0 | 0 |
| reinstall-restores-files | MATCH | 0 | 0 |
| reinstall-no-dev-mode | MATCH | 0 | 0 |
| reinstall-scripts | MATCH | 0 | 0 |
| reinstall-no-scripts | MATCH | 0 | 0 |
| reinstall-no-json | MATCH | 1 | 1 |
| reinstall-prefer-source | MATCH | 0 | 0 |
| reinstall-prefer-install-source | MATCH | 0 | 0 |
| reinstall-source-to-dist | MATCH | 0 | 0 |
| reinstall-source-default | MATCH | 0 | 0 |
| reinstall-prefer-conflict | MATCH | 1 | 1 |
| reinstall-git | MATCH | 0 | 0 |
| reinstall-git-changes | MATCH | 1 | 1 |
| reinstall-git-discard | MATCH | 0 | 0 |
| status-dist-clean | MATCH | 0 | 0 |
| status-dist-clean-v | MATCH | 0 | 0 |
| status-dist-changed | MATCH | 1 | 1 |
| status-dist-changed-v | MATCH | 1 | 1 |
| status-dist-changed-vv | MATCH | 1 | 1 |
| status-dist-removed-v | MATCH | 1 | 1 |
| status-dist-two-packages | MATCH | 1 | 1 |
| status-not-installed | MATCH | 0 | 0 |
| status-lock-only | MATCH | 0 | 0 |
| status-no-json | MATCH | 1 | 1 |
| status-scripts | MATCH | 0 | 0 |
| status-no-scripts | MATCH | 0 | 0 |
| status-git-clean | MATCH | 0 | 0 |
| status-git-modified | MATCH | 1 | 1 |
| status-git-modified-v | MATCH | 1 | 1 |
| status-git-untracked | MATCH | 0 | 0 |
| status-git-unpushed | MATCH | 0 | 0 |
| status-git-unpushed-v | MATCH | 0 | 0 |
| status-git-unpushed-vv | MATCH | 0 | 0 |
| status-git-version-v | MATCH | 4 | 4 |
| status-git-version-vv | MATCH | 4 | 4 |
| status-git-branch-v | MATCH | 6 | 6 |
| status-mixed | MATCH | 1 | 1 |
| status-mixed-v | MATCH | 1 | 1 |
| status-path-symlink | MATCH | 1 | 1 |
| status-path-mirror | MATCH | 0 | 0 |
| status-path-mirror-changed | MATCH | 1 | 1 |
| create-project-basic | MATCH | 0 | 0 |
| create-project-default-dir | MATCH | 0 | 0 |
| create-project-version-arg | MATCH | 0 | 0 |
| create-project-version-colon | MATCH | 0 | 0 |
| create-project-version-equals | MATCH | 0 | 0 |
| create-project-stability-flag | MATCH | 0 | 0 |
| create-project-stability-at | MATCH | 0 | 0 |
| create-project-version-beta | MATCH | 0 | 0 |
| create-project-stability-invalid | MATCH | 1 | 1 |
| create-project-not-found | MATCH | 1 | 1 |
| create-project-not-found-version | MATCH | 1 | 1 |
| create-project-php-too-low | MATCH | 1 | 1 |
| create-project-php-ignore | MATCH | 0 | 0 |
| create-project-secure-http | MATCH | 100 | 100 |
| create-project-repository-url | MATCH | 0 | 0 |
| create-project-packagist-false-repo | MATCH | 0 | 0 |
| create-project-no-install | MATCH | 0 | 0 |
| create-project-no-scripts | MATCH | 0 | 0 |
| create-project-no-dev | MATCH | 0 | 0 |
| create-project-dev-option | MATCH | 0 | 0 |
| create-project-no-custom-installers | MATCH | 0 | 0 |
| create-project-no-plugins | MATCH | 0 | 0 |
| create-project-no-audit | MATCH | 0 | 0 |
| create-project-add-repository | MATCH | 0 | 0 |
| create-project-require | MATCH | 0 | 0 |
| create-project-ask-noninteractive | MATCH | 0 | 0 |
| create-project-prefer-conflict | MATCH | 1 | 1 |
| create-project-prefer-source-dist-only | MATCH | 0 | 0 |
| create-project-verbose | MATCH | 0 | 0 |
| create-project-progress | MATCH | 0 | 0 |
| create-project-relative-dir | MATCH | 0 | 0 |
| create-project-absolute-dir | MATCH | 0 | 0 |
| create-project-plain | MATCH | 0 | 0 |
| create-project-bin | MATCH | 0 | 0 |
| create-project-packagist | MATCH | 100 | 100 |
| create-project-not-empty | MATCH | 1 | 1 |
| create-project-dir-is-file | MATCH | 1 | 1 |
| create-project-path | MATCH | 0 | 0 |
| create-project-path-symlink | MATCH | 0 | 0 |
| create-project-vcs | MATCH | 0 | 0 |
| create-project-vcs-keep | MATCH | 0 | 0 |
| create-project-vcs-remove | MATCH | 0 | 0 |
| create-project-vcs-dev | MATCH | 0 | 0 |
| create-project-vcs-branch | MATCH | 0 | 0 |
| create-project-dot | MATCH | 0 | 0 |
| create-project-audit-format | MATCH | 0 | 0 |
| create-project-prefer-install-auto | MATCH | 0 | 0 |
| create-project-prefer-install-bad | MATCH | 1 | 1 |
| create-project-stability-rc-upper | MATCH | 0 | 0 |
| create-project-composer-env | MATCH | 0 | 0 |
| create-project-no-args | MATCH | 0 | 0 |
| create-project-no-args-no-json | MATCH | 1 | 1 |
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
| exec-list | MATCH | 0 | 0 |
| exec-no-binary | MATCH | 0 | 0 |
| exec-run | MATCH | 0 | 0 |
| exec-run-args | MATCH | 0 | 0 |
| exec-run-php | MATCH | 0 | 0 |
| exec-missing | MATCH | 127 | 127 |
| exec-no-scripts | MATCH | 0 | 0 |
| exec-verbose | MATCH | 0 | 0 |
| exec-no-bins | MATCH | 1 | 1 |
| exec-local-list | MATCH | 0 | 0 |
| exec-no-json | MATCH | 1 | 1 |
| clear-cache-basic | MATCH | 0 | 0 |
| clear-cache-alias-cc | MATCH | 0 | 0 |
| clear-cache-alias-clearcache | MATCH | 0 | 0 |
| clear-cache-gc | MATCH | 0 | 0 |
| clear-cache-gc-maxsize | MATCH | 0 | 0 |
| clear-cache-gc-ttl-fresh | MATCH | 0 | 0 |
| clear-cache-gc-bad-maxsize | MATCH | 1 | 1 |
| clear-cache-read-only | MATCH | 0 | 0 |
| clear-cache-verbose | MATCH | 0 | 0 |
| clear-cache-empty | MATCH | 0 | 0 |
| clear-cache-empty-gc | MATCH | 0 | 0 |
| clear-cache-no-cache | MATCH | 0 | 0 |
| clear-cache-no-json | MATCH | 0 | 0 |
| clear-cache-custom-dirs | MATCH | 0 | 0 |
| fund-text | MATCH | 0 | 0 |
| fund-json | MATCH | 0 | 0 |
| fund-json-short | MATCH | 0 | 0 |
| fund-bad-format | MATCH | 1 | 1 |
| fund-v2 | MATCH | 0 | 0 |
| fund-none | MATCH | 0 | 0 |
| fund-none-json | MATCH | 0 | 0 |
| fund-not-installed | MATCH | 0 | 0 |
| fund-no-json | MATCH | 1 | 1 |
| search-fulltext | MATCH | 0 | 0 |
| search-multi-token | MATCH | 0 | 0 |
| search-keyword | MATCH | 0 | 0 |
| search-only-name | MATCH | 0 | 0 |
| search-only-name-long | MATCH | 0 | 0 |
| search-only-vendor | MATCH | 0 | 0 |
| search-only-vendor-long | MATCH | 0 | 0 |
| search-name-and-vendor | MATCH | 1 | 1 |
| search-type | MATCH | 0 | 0 |
| search-type-other | MATCH | 0 | 0 |
| search-type-vendor | MATCH | 0 | 0 |
| search-json | MATCH | 0 | 0 |
| search-json-vendor | MATCH | 0 | 0 |
| search-json-empty | MATCH | 0 | 0 |
| search-bad-format | MATCH | 1 | 1 |
| search-no-results | MATCH | 0 | 0 |
| search-regex-chars | MATCH | 0 | 0 |
| search-regex-chars-name | MATCH | 0 | 0 |
| search-no-tokens | MATCH | 1 | 1 |
| search-very-verbose | MATCH | 0 | 0 |
| search-columns | MATCH | 0 | 0 |
| search-installed | MATCH | 0 | 0 |
| search-installed-vv | MATCH | 0 | 0 |
| search-v2 | MATCH | 0 | 0 |
| search-v2-vendor | MATCH | 0 | 0 |
| search-v2-json | MATCH | 0 | 0 |
| search-wide | MATCH | 0 | 0 |
| search-filtered-repo | MATCH | 0 | 0 |
| browse-show | MATCH | 0 | 0 |
| browse-show-homepage | MATCH | 0 | 0 |
| browse-home-alias | MATCH | 0 | 0 |
| browse-homepage-missing | MATCH | 1 | 1 |
| browse-not-found | MATCH | 1 | 1 |
| browse-no-url | MATCH | 1 | 1 |
| browse-multi | MATCH | 1 | 1 |
| browse-open | MATCH | 0 | 0 |
| browse-root | MATCH | 1 | 1 |
| browse-root-homepage | MATCH | 0 | 0 |
| browse-root-source | MATCH | 0 | 0 |
| browse-root-open | MATCH | 0 | 0 |
| browse-root-invalid-url | MATCH | 1 | 1 |
| browse-installed | MATCH | 0 | 0 |
| browse-no-json | MATCH | 1 | 1 |
| diagnose-basic | MATCH | 1 | 1 |
| diagnose-clean | MATCH | 0 | 0 |
| diagnose-locked | MATCH | 0 | 0 |
| diagnose-lock-invalid | MATCH | 2 | 2 |
| diagnose-no-json | MATCH | 0 | 0 |
| diagnose-schema-error | MATCH | 0 | 0 |
| diagnose-platform-override | MATCH | 0 | 0 |
| diagnose-disable-tls | MATCH | 0 | 0 |
| diagnose-github-oauth | MATCH | 0 | 0 |
| diagnose-proxy | MATCH | 0 | 0 |
| diagnose-proxy-no-proxy | MATCH | 0 | 0 |
| diagnose-proxy-invalid | MATCH | 0 | 0 |
| diagnose-ipresolve | MATCH | 0 | 0 |
| diagnose-network | KNOWN-DIFF | 2 | 2 |
| self-update-basic | KNOWN-DIFF | 1 | 1 |
| self-update-alias | KNOWN-DIFF | 1 | 1 |
| self-update-rollback | KNOWN-DIFF | 1 | 1 |
| self-update-version | KNOWN-DIFF | 1 | 1 |
| self-update-no-json | KNOWN-DIFF | 1 | 1 |
| self-update-bad-option | MATCH | 1 | 1 |
| self-update-too-many | MATCH | 1 | 1 |
| repository-list | MATCH | 0 | 0 |
| repository-list-default | MATCH | 0 | 0 |
| repository-list-named | MATCH | 0 | 0 |
| repository-list-assoc | MATCH | 0 | 0 |
| repository-no-args | MATCH | 0 | 0 |
| repository-add | MATCH | 0 | 0 |
| repository-add-append | MATCH | 0 | 0 |
| repository-add-prepend-named | MATCH | 0 | 0 |
| repository-add-assoc | MATCH | 0 | 0 |
| repository-add-json | MATCH | 0 | 0 |
| repository-add-json-invalid | MATCH | 1 | 1 |
| repository-add-missing-name | MATCH | 1 | 1 |
| repository-add-missing-arg1 | MATCH | 1 | 1 |
| repository-add-missing-url | MATCH | 1 | 1 |
| repository-add-before | MATCH | 0 | 0 |
| repository-add-after | MATCH | 0 | 0 |
| repository-add-before-after | MATCH | 1 | 1 |
| repository-add-before-missing | MATCH | 1 | 1 |
| repository-add-existing | MATCH | 0 | 0 |
| repository-remove | MATCH | 0 | 0 |
| repository-remove-assoc | MATCH | 0 | 0 |
| repository-remove-missing | MATCH | 0 | 0 |
| repository-remove-packagist | MATCH | 0 | 0 |
| repository-remove-no-name | MATCH | 1 | 1 |
| repository-set-url | MATCH | 0 | 0 |
| repository-set-url-assoc | MATCH | 0 | 0 |
| repository-set-url-missing | MATCH | 1 | 1 |
| repository-set-url-unknown | MATCH | 0 | 0 |
| repository-get-url | MATCH | 0 | 0 |
| repository-get-url-assoc | MATCH | 0 | 0 |
| repository-get-url-no-url | MATCH | 1 | 1 |
| repository-get-url-missing | MATCH | 1 | 1 |
| repository-get-url-no-name | MATCH | 1 | 1 |
| repository-disable | MATCH | 0 | 0 |
| repository-disable-short | MATCH | 0 | 0 |
| repository-disable-other | MATCH | 1 | 1 |
| repository-disable-no-name | MATCH | 1 | 1 |
| repository-enable | MATCH | 0 | 0 |
| repository-enable-other | MATCH | 1 | 1 |
| repository-unknown-action | MATCH | 1 | 1 |
| repository-uppercase-action | MATCH | 0 | 0 |
| repository-global-add | MATCH | 0 | 0 |
| repository-global-list | MATCH | 0 | 0 |
| repository-file | MATCH | 1 | 1 |
| repository-file-and-global | MATCH | 1 | 1 |
| repository-file-missing | MATCH | 1 | 1 |
| repository-file-add | MATCH | 0 | 0 |
| repository-no-json | MATCH | 1 | 1 |
| policy-cmd-add-source | MATCH | 0 | 0 |
| policy-cmd-add-source-json | MATCH | 0 | 0 |
| policy-cmd-add-source-existing-list | MATCH | 0 | 0 |
| policy-cmd-add-source-duplicate | MATCH | 0 | 0 |
| policy-cmd-add-source-object-sources | MATCH | 1 | 1 |
| policy-cmd-add-source-uppercase-action | MATCH | 0 | 0 |
| policy-cmd-no-action | MATCH | 1 | 1 |
| policy-cmd-unknown-action | MATCH | 1 | 1 |
| policy-cmd-no-name | MATCH | 1 | 1 |
| policy-cmd-builtin | MATCH | 1 | 1 |
| policy-cmd-builtin-cooldown | MATCH | 1 | 1 |
| policy-cmd-reserved-prefix | MATCH | 1 | 1 |
| policy-cmd-reserved-name | MATCH | 1 | 1 |
| policy-cmd-dotted-name | MATCH | 1 | 1 |
| policy-cmd-empty-name | MATCH | 1 | 1 |
| policy-cmd-no-arg1 | MATCH | 1 | 1 |
| policy-cmd-no-url | MATCH | 1 | 1 |
| policy-cmd-bad-type | MATCH | 1 | 1 |
| policy-cmd-http-url | MATCH | 1 | 1 |
| policy-cmd-json-no-type | MATCH | 1 | 1 |
| policy-cmd-json-no-url | MATCH | 1 | 1 |
| policy-cmd-json-int-url | MATCH | 1 | 1 |
| policy-cmd-json-invalid | MATCH | 1 | 1 |
| policy-cmd-global | MATCH | 0 | 0 |
| policy-cmd-file-and-global | MATCH | 1 | 1 |
| policy-cmd-file | MATCH | 0 | 0 |
| policy-cmd-no-json | MATCH | 1 | 1 |
| plugin-blocked-update | DIFF | 1 | 0 |
| plugin-blocked-install-lock | DIFF | 1 | 0 |
| plugin-blocked-installed | DIFF | 1 | 0 |
| plugin-allow-false | MATCH | 0 | 0 |
| plugin-allow-all-false | MATCH | 0 | 0 |
| plugin-allow-true | DIFF | 0 | 0 |
| plugin-optional-blocked | MATCH | 0 | 0 |
| plugin-not-required-anymore | DIFF | 0 | 0 |
| plugin-old-lock-no-allow | DIFF | 1 | 0 |
| plugin-old-lock-new-api | DIFF | 1 | 0 |
| plugin-allowed-update | DIFF | 0 | 0 |
| plugin-allowed-update-verbose | DIFF | 0 | 0 |
| plugin-allowed-update-debug | DIFF | 0 | 0 |
| plugin-update-dry-run | MATCH | 0 | 0 |
| plugin-install-from-lock | DIFF | 0 | 1 |
| plugin-install-no-dev | DIFF | 0 | 1 |
| plugin-installed-noop | DIFF | 0 | 1 |
| plugin-installed-quiet | DIFF | 0 | 1 |
| plugin-no-scripts | DIFF | 0 | 1 |
| plugin-dump-autoload | DIFF | 0 | 0 |
| plugin-dump-autoload-optimize | DIFF | 0 | 0 |
| plugin-no-plugins-installed | MATCH | 0 | 0 |
| plugin-no-plugins-update | DIFF | 0 | 0 |
| plugin-update-version | DIFF | 0 | 1 |
| plugin-remove | DIFF | 0 | 1 |
| plugin-require | DIFF | 0 | 0 |
| plugin-require-installed | DIFF | 0 | 1 |
| plugin-run-script-custom | DIFF | 0 | 0 |
| plugin-show | DIFF | 0 | 0 |
| plugin-thrower-boom | ERROR | 0 | 0 |
| plugin-thrower-stop | DIFF | 0 | 0 |
| plugin-thrower-false | DIFF | 0 | 1 |
| plugin-old-api | MATCH | 2 | 2 |
| plugin-no-api | DIFF | 1 | 0 |
| plugin-no-class | DIFF | 1 | 0 |
| plugin-legacy-installer | DIFF | 0 | 0 |
| plugin-command | DIFF | 0 | 1 |
| plugin-command-default | DIFF | 0 | 1 |
| plugin-command-exit | DIFF | 3 | 1 |
| plugin-command-bad-option | DIFF | 1 | 1 |
| plugin-command-no-plugins | DIFF | 1 | 1 |
| plugin-installer-update | DIFF | 0 | 0 |
| plugin-installer-extra-path | DIFF | 0 | 0 |
| plugin-installer-install | DIFF | 0 | 1 |
| plugin-installer-upgrade | DIFF | 0 | 1 |
| plugin-installer-remove | DIFF | 0 | 1 |
| plugin-installer-no-plugins | DIFF | 0 | 0 |
| plugin-installer-dump | DIFF | 0 | 1 |
| show-latest | MATCH | 0 | 0 |
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
| vcs-verbose | MATCH | 0 | 0 |
| vcs-debug | MATCH | 0 | 0 |
| vcs-only-branches | MATCH | 0 | 0 |
| vcs-only-branches-invalid | MATCH | 1 | 1 |
| vcs-misc-vv | MATCH | 0 | 0 |
| vcs-misc-numeric-branch | MATCH | 0 | 0 |
| vcs-misc-slash-branch | MATCH | 0 | 0 |
| vcs-misc-old-tag | MATCH | 0 | 0 |
| vcs-ref-pin | MATCH | 0 | 0 |
| vcs-not-a-repo | MATCH | 1 | 1 |
| vcs-git-type-not-a-repo | MATCH | 1 | 1 |
| vcs-type-hg | MATCH | 1 | 1 |
| vcs-type-svn | MATCH | 1 | 1 |
| vcs-type-fossil | MATCH | 1 | 1 |
| vcs-type-perforce | MATCH | 1 | 1 |
| vcs-type-pear | MATCH | 1 | 1 |
| vcs-type-forgejo | MATCH | 1 | 1 |
| vcs-type-github-local | MATCH | 1 | 1 |
| vcs-insecure-http | MATCH | 100 | 100 |
| vcs-http-missing | MATCH | 1 | 1 |
| vcs-http-cached | MATCH | 0 | 0 |
| vcs-http-cached-debug | MATCH | 0 | 0 |
| vcs-http-dev-update | MATCH | 0 | 0 |
| source-update-checkout | MATCH | 0 | 0 |
| source-update-checkout-v | MATCH | 0 | 0 |
| source-downgrade-v | MATCH | 0 | 0 |
| source-dev-update | MATCH | 0 | 0 |
| source-local-changes | MATCH | 1 | 1 |
| source-discard-true | MATCH | 0 | 0 |
| source-discard-stash | MATCH | 1 | 1 |
| source-discard-stash-ident | MATCH | 1 | 1 |
| source-uninstall-changes | MATCH | 1 | 1 |
| source-uninstall-discard | MATCH | 0 | 0 |
| source-untracked-only | MATCH | 0 | 0 |
| source-prefer-source | MATCH | 0 | 0 |
| source-prefer-dist-default | MATCH | 0 | 0 |
| source-prefer-install-flag | MATCH | 0 | 0 |
| source-preferred-install-source | MATCH | 0 | 0 |
| source-preferred-install-pattern | MATCH | 0 | 0 |
| source-preferred-install-auto-dev | MATCH | 0 | 0 |
| source-dev-prefer-dist | MATCH | 0 | 0 |
| source-dist-to-source | MATCH | 0 | 0 |
| source-source-to-dist | MATCH | 0 | 0 |
| source-source-to-dist-changes | MATCH | 1 | 1 |
| source-fallback-dist | MATCH | 0 | 0 |
| source-install-lock-prefer-source | MATCH | 0 | 0 |
| source-reinstall-same | MATCH | 0 | 0 |
| vcs-git-install | MATCH | 0 | 0 |
| vcs-git-type | MATCH | 0 | 0 |
| vcs-file-url | MATCH | 0 | 0 |
| vcs-dev-main | MATCH | 0 | 0 |
| vcs-dev-branch | MATCH | 0 | 0 |
| vcs-dep | MATCH | 0 | 0 |
| vcs-install-from-lock | MATCH | 0 | 0 |
| vcs-dry-run | MATCH | 0 | 0 |
| vcs-no-composer-json | MATCH | 1 | 1 |
| vcs-missing-path | MATCH | 1 | 1 |
| vcs-odd-refs | MATCH | 0 | 0 |
| vcs-odd-refs-vv | MATCH | 0 | 0 |
| vcs-http-mirror | MATCH | 0 | 0 |
| vcs-forge-github | MATCH | 0 | 0 |
| vcs-forge-github-vv | MATCH | 0 | 0 |
| vcs-forge-github-type | MATCH | 0 | 0 |
| vcs-forge-github-1x-dev | MATCH | 0 | 0 |
| vcs-forge-github-install | MATCH | 0 | 0 |
| vcs-forge-github-missing | MATCH | 1 | 1 |
| vcs-forge-github-private-noauth | MATCH | 1 | 1 |
| vcs-forge-github-private-nofallback | MATCH | 1 | 1 |
| vcs-forge-github-private-auth | MATCH | 0 | 0 |
| vcs-forge-github-private-auth-lock | MATCH | 0 | 0 |
| vcs-forge-github-enterprise | DIFF | 0 | 0 |
| vcs-forge-github-not-domain | MATCH | 0 | 0 |
| vcs-forge-git-url-credentials | MATCH | 1 | 1 |
| vcs-forge-git-http-credentials | MATCH | 0 | 0 |
| vcs-forge-gitlab | MATCH | 0 | 0 |
| vcs-forge-gitlab-vv | MATCH | 0 | 0 |
| vcs-forge-gitlab-dev | MATCH | 0 | 0 |
| vcs-forge-gitlab-subgroup | MATCH | 0 | 0 |
| vcs-forge-gitlab-http-port | MATCH | 0 | 0 |
| vcs-forge-gitlab-install | MATCH | 0 | 0 |
| vcs-forge-gitlab-private-noauth | MATCH | 1 | 1 |
| vcs-forge-gitlab-private-token | MATCH | 0 | 0 |
| vcs-forge-gitlab-private-oauth | MATCH | 0 | 0 |
| vcs-forge-gitlab-protocol-invalid | MATCH | 1 | 1 |
| vcs-forge-gitlab-protocol-git | MATCH | 0 | 0 |
| vcs-forge-bitbucket | MATCH | 0 | 0 |
| vcs-forge-bitbucket-vv | MATCH | 0 | 0 |
| vcs-forge-bitbucket-dev | MATCH | 0 | 0 |
| vcs-forge-bitbucket-install | MATCH | 0 | 0 |
| vcs-forge-bitbucket-private-noauth | MATCH | 100 | 100 |
| vcs-forge-bitbucket-private-basic | MATCH | 0 | 0 |
| vcs-artifact-basic | MATCH | 0 | 0 |
| vcs-artifact-pin-old | MATCH | 0 | 0 |
| vcs-artifact-verbose | MATCH | 0 | 0 |
| vcs-artifact-tgz | MATCH | 0 | 0 |
| vcs-artifact-tar | MATCH | 0 | 0 |
| vcs-artifact-dep | MATCH | 0 | 0 |
| vcs-artifact-install-lock | MATCH | 0 | 0 |
| vcs-artifact-invalid-json | MATCH | 1 | 1 |
| vcs-artifact-no-composer-json | MATCH | 2 | 2 |
| vcs-artifact-deep-composer-json | MATCH | 2 | 2 |
| vcs-artifact-two-top-dirs | MATCH | 2 | 2 |
| vcs-artifact-macosx | MATCH | 0 | 0 |
| vcs-artifact-no-version | MATCH | 1 | 1 |
| vcs-artifact-bad-version | MATCH | 1 | 1 |
| vcs-artifact-empty-zip | MATCH | 2 | 2 |
| vcs-artifact-corrupt-zip | MATCH | 2 | 2 |
| vcs-artifact-corrupt-tgz | MATCH | 2 | 2 |
| vcs-artifact-other-files | MATCH | 0 | 0 |
| vcs-artifact-missing-dir | MATCH | 1 | 1 |
| vcs-artifact-absolute-url | MATCH | 0 | 0 |

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

## bump-lock-disabled

```
stderr:
- (oracle)   [LogicException]                                   
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Package/Locker.php:342
- (oracle)  Composer\Package\Locker->getLockData() at /oracle-source/src/Composer/Package/Locker.php:168
- (oracle)  Composer\Package\Locker->getLockedRepository() at /oracle-source/src/Composer/Command/BumpCommand.php:126
- (oracle)  Composer\Command\BumpCommand->doBump() at /oracle-source/src/Composer/Command/BumpCommand.php:75
- (oracle)  Composer\Command\BumpCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## diagnose-network

```
stdout:
+ (go) Checking Composer and its dependencies for vulnerabilities: WARNING
+ (go) Could not find Composer's installed.json, this must be a non-standard Composer installation.
- (oracle) Checking Composer and its dependencies for vulnerabilities: Failed performing audit: curl error 6 while downloading https://repo.packagist.org/packages.json: Could not resolve host: repo.packagist.org
```

## self-update-basic

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-alias

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-rollback

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-version

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-no-json

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## plugin-blocked-update

```
exit code: oracle=1 go=0
stdout:
+ (go) No security vulnerability advisories found.
stderr:
+ (go)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
+ (go)   - Installing acme/log (1.1.0): Extracting archive
+ (go) 1 package suggestions were added by new dependencies, use `composer suggest` to see details.
+ (go) Generating autoload files
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   plugfix/events-plugin contains a Composer plugin which is blocked by your a  
- (oracle)   llow-plugins config. You may add it to the list if you consider it safe.     
- (oracle)   You can run "composer config --no-plugins allow-plugins.plugfix/events-plug  
- (oracle)   in [true|false]" to enable it (true) or disable it explicitly and suppress   
- (oracle)   this exception (false)                                                       
- (oracle)   See https://getcomposer.org/allow-plugins                                    
- (oracle)                                                                                
- (oracle) 
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor
+ (go) D vendor/acme
+ (go) D vendor/acme/log
+ (go) D vendor/acme/log/src
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/events-plugin
+ (go) D vendor/plugfix/events-plugin/src
+ (go) F vendor/acme/log/composer.json 644 b83d42805152f028
+ (go) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
+ (go) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
+ (go) F vendor/autoload.php 644 6c388903f320a749
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 0635e595c1470c97
+ (go) F vendor/composer/autoload_real.php 644 12937f6ad87f4318
+ (go) F vendor/composer/autoload_static.php 644 82493a6034d2c59c
+ (go) F vendor/composer/installed.json 644 c5a4556d1c372b7a
+ (go) F vendor/composer/installed.php 644 e3f9a94f37234285
+ (go) F vendor/composer/platform_check.php 644 efa64a484ec52731
+ (go) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
+ (go) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
+ (go) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         require __DIR__ . '/platform_check.php';
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInit6de87ac065017bb6df625f1693eb1809::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInit6de87ac065017bb6df625f1693eb1809
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\Events\\' => 15,
+ (go)         ),
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)             'Acme\\Log\\' => 9,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'Plugfix\\Events\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
+ (go)         ),
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)         'Acme\\Log\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/acme/log/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "acme/log",
+ (go)             "version": "1.1.0",
+ (go)             "version_normalized": "1.1.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip",
+ (go)                 "reference": "58b0f69fc5f96981e5f7fae32d592232a36012c2",
+ (go)                 "shasum": "c40f1901ba3e909555618861c594e339cd678f39"
+ (go)             },
+ (go)             "require": {
+ (go)                 "php": ">=7.4"
+ (go)             },
+ (go)             "suggest": {
+ (go)                 "acme/extra": "for extra stuff"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "library",
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Acme\\Log\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Logging",
+ (go)             "install-path": "../acme/log"
+ (go)         },
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "install-path": "../plugfix/events-plugin"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'acme/log' => array(
+ (go)             'pretty_version' => '1.1.0',
+ (go)             'version' => '1.1.0.0',
+ (go)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
+ (go)             'type' => 'library',
+ (go)             'install_path' => __DIR__ . '/../acme/log',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'plugfix/events-plugin' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
vendor/composer/platform_check.php:
+ (go) <?php
+ (go) // platform_check.php @generated by Composer
+ (go) 
+ (go) $issues = array();
+ (go) 
+ (go) if (!(PHP_VERSION_ID >= 70400)) {
+ (go)     $issues[] = 'Your Composer dependencies require a PHP version ">= 7.4.0". You are running ' . PHP_VERSION . '.';
+ (go) }
+ (go) 
+ (go) if ($issues) {
+ (go)     if (!headers_sent()) {
+ (go)         header('HTTP/1.1 500 Internal Server Error');
+ (go)     }
+ (go)     if (!ini_get('display_errors')) {
+ (go)         if (PHP_SAPI === 'cli' || PHP_SAPI === 'phpdbg') {
+ (go)             fwrite(STDERR, 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . implode(PHP_EOL, $issues) . PHP_EOL.PHP_EOL);
+ (go)         } elseif (!headers_sent()) {
+ (go)             echo 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . str_replace('You are running '.PHP_VERSION.'.', '', implode(PHP_EOL, $issues)) . PHP_EOL.PHP_EOL;
+ (go)         }
+ (go)     }
+ (go)     throw new \RuntimeException(
+ (go)         'Composer detected issues in your platform: ' . implode(' ', $issues)
+ (go)     );
+ (go) }
+ (go) 
```

## plugin-blocked-install-lock

```
exit code: oracle=1 go=0
stderr:
+ (go)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
+ (go)   - Installing acme/log (1.1.0): Extracting archive
+ (go) Generating autoload files
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   plugfix/events-plugin contains a Composer plugin which is blocked by your a  
- (oracle)   llow-plugins config. You may add it to the list if you consider it safe.     
- (oracle)   You can run "composer config --no-plugins allow-plugins.plugfix/events-plug  
- (oracle)   in [true|false]" to enable it (true) or disable it explicitly and suppress   
- (oracle)   this exception (false)                                                       
- (oracle)   See https://getcomposer.org/allow-plugins                                    
- (oracle)                                                                                
- (oracle) 
- (oracle) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor
+ (go) D vendor/acme
+ (go) D vendor/acme/log
+ (go) D vendor/acme/log/src
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/events-plugin
+ (go) D vendor/plugfix/events-plugin/src
+ (go) F vendor/acme/log/composer.json 644 b83d42805152f028
+ (go) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
+ (go) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
+ (go) F vendor/autoload.php 644 6c388903f320a749
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 0635e595c1470c97
+ (go) F vendor/composer/autoload_real.php 644 12937f6ad87f4318
+ (go) F vendor/composer/autoload_static.php 644 82493a6034d2c59c
+ (go) F vendor/composer/installed.json 644 c5a4556d1c372b7a
+ (go) F vendor/composer/installed.php 644 e3f9a94f37234285
+ (go) F vendor/composer/platform_check.php 644 efa64a484ec52731
+ (go) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
+ (go) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
+ (go) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         require __DIR__ . '/platform_check.php';
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInit6de87ac065017bb6df625f1693eb1809', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInit6de87ac065017bb6df625f1693eb1809::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInit6de87ac065017bb6df625f1693eb1809
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\Events\\' => 15,
+ (go)         ),
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)             'Acme\\Log\\' => 9,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'Plugfix\\Events\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
+ (go)         ),
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)         'Acme\\Log\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/acme/log/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInit6de87ac065017bb6df625f1693eb1809::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "acme/log",
+ (go)             "version": "1.1.0",
+ (go)             "version_normalized": "1.1.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip",
+ (go)                 "reference": "58b0f69fc5f96981e5f7fae32d592232a36012c2",
+ (go)                 "shasum": "c40f1901ba3e909555618861c594e339cd678f39"
+ (go)             },
+ (go)             "require": {
+ (go)                 "php": ">=7.4"
+ (go)             },
+ (go)             "suggest": {
+ (go)                 "acme/extra": "for extra stuff"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "library",
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Acme\\Log\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Logging",
+ (go)             "install-path": "../acme/log"
+ (go)         },
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "install-path": "../plugfix/events-plugin"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'acme/log' => array(
+ (go)             'pretty_version' => '1.1.0',
+ (go)             'version' => '1.1.0.0',
+ (go)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
+ (go)             'type' => 'library',
+ (go)             'install_path' => __DIR__ . '/../acme/log',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'plugfix/events-plugin' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
vendor/composer/platform_check.php:
+ (go) <?php
+ (go) // platform_check.php @generated by Composer
+ (go) 
+ (go) $issues = array();
+ (go) 
+ (go) if (!(PHP_VERSION_ID >= 70400)) {
+ (go)     $issues[] = 'Your Composer dependencies require a PHP version ">= 7.4.0". You are running ' . PHP_VERSION . '.';
+ (go) }
+ (go) 
+ (go) if ($issues) {
+ (go)     if (!headers_sent()) {
+ (go)         header('HTTP/1.1 500 Internal Server Error');
+ (go)     }
+ (go)     if (!ini_get('display_errors')) {
+ (go)         if (PHP_SAPI === 'cli' || PHP_SAPI === 'phpdbg') {
+ (go)             fwrite(STDERR, 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . implode(PHP_EOL, $issues) . PHP_EOL.PHP_EOL);
+ (go)         } elseif (!headers_sent()) {
+ (go)             echo 'Composer detected issues in your platform:' . PHP_EOL.PHP_EOL . str_replace('You are running '.PHP_VERSION.'.', '', implode(PHP_EOL, $issues)) . PHP_EOL.PHP_EOL;
+ (go)         }
+ (go)     }
+ (go)     throw new \RuntimeException(
+ (go)         'Composer detected issues in your platform: ' . implode(' ', $issues)
+ (go)     );
+ (go) }
+ (go) 
```

## plugin-blocked-installed

```
exit code: oracle=1 go=0
stdout:
+ (go) root-post-install
stderr:
+ (go) Installing dependencies from lock file (including require-dev)
+ (go) Verifying lock file contents can be installed on current platform.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Nothing to install, update or remove
+ (go) Generating autoload files
+ (go) > echo root-post-install
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   plugfix/events-plugin contains a Composer plugin which is blocked by your a  
- (oracle)   llow-plugins config. You may add it to the list if you consider it safe.     
- (oracle)   You can run "composer config --no-plugins allow-plugins.plugfix/events-plug  
- (oracle)   in [true|false]" to enable it (true) or disable it explicitly and suppress   
- (oracle)   this exception (false)                                                       
- (oracle)   See https://getcomposer.org/allow-plugins                                    
- (oracle)                                                                                
- (oracle) 
- (oracle) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
- (oracle) 
- (oracle) 
```

## plugin-allow-true

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=2
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,plugfix/events-plugin@1.0.0
```

## plugin-not-required-anymore

```
stderr:
- (oracle) The "plugfix/events-plugin" plugin was not loaded as it is not listed in allow-plugins and is not required by the root package anymore.
```

## plugin-old-lock-no-allow

```
exit code: oracle=1 go=0
stderr:
+ (go)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
+ (go) Generating autoload files
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Your composer.lock was generated before the allow-plugins security feature   
- (oracle)   was introduced and your composer.json does not define allow-plugins. Run "c  
- (oracle)   omposer update --lock" locally and commit the updated composer.lock, then a  
- (oracle)   dd an explicit allow-plugins section to composer.json. See https://getcompo  
- (oracle)   ser.org/allow-plugins                                                        
- (oracle)                                                                                
- (oracle) 
- (oracle) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/events-plugin
+ (go) D vendor/plugfix/events-plugin/src
+ (go) F vendor/autoload.php 644 037fb318047344d9
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 cf077d4728aedc14
+ (go) F vendor/composer/autoload_real.php 644 4719fb7042838000
+ (go) F vendor/composer/autoload_static.php 644 5024c222cb822903
+ (go) F vendor/composer/installed.json 644 a64dfcd2a0c28cdb
+ (go) F vendor/composer/installed.php 644 c8f18d553995a3e0
+ (go) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
+ (go) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
+ (go) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInit00000000000000000000000000000000
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInit00000000000000000000000000000000', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInit00000000000000000000000000000000', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInit00000000000000000000000000000000::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInit00000000000000000000000000000000
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\Events\\' => 15,
+ (go)         ),
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'Plugfix\\Events\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
+ (go)         ),
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInit00000000000000000000000000000000::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInit00000000000000000000000000000000::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInit00000000000000000000000000000000::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "install-path": "../plugfix/events-plugin"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'plugfix/events-plugin' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
```

## plugin-old-lock-new-api

```
exit code: oracle=1 go=0
stderr:
+ (go)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
+ (go) Generating autoload files
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   plugfix/events-plugin contains a Composer plugin which is blocked by your a  
- (oracle)   llow-plugins config. You may add it to the list if you consider it safe.     
- (oracle)   You can run "composer config --no-plugins allow-plugins.plugfix/events-plug  
- (oracle)   in [true|false]" to enable it (true) or disable it explicitly and suppress   
- (oracle)   this exception (false)                                                       
- (oracle)   See https://getcomposer.org/allow-plugins                                    
- (oracle)                                                                                
- (oracle) 
- (oracle) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/events-plugin
+ (go) D vendor/plugfix/events-plugin/src
+ (go) F vendor/autoload.php 644 037fb318047344d9
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 cf077d4728aedc14
+ (go) F vendor/composer/autoload_real.php 644 4719fb7042838000
+ (go) F vendor/composer/autoload_static.php 644 5024c222cb822903
+ (go) F vendor/composer/installed.json 644 a64dfcd2a0c28cdb
+ (go) F vendor/composer/installed.php 644 c8f18d553995a3e0
+ (go) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
+ (go) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
+ (go) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInit00000000000000000000000000000000
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInit00000000000000000000000000000000', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInit00000000000000000000000000000000', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInit00000000000000000000000000000000::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInit00000000000000000000000000000000
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\Events\\' => 15,
+ (go)         ),
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'Plugfix\\Events\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
+ (go)         ),
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInit00000000000000000000000000000000::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInit00000000000000000000000000000000::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInit00000000000000000000000000000000::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "install-path": "../plugfix/events-plugin"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'plugfix/events-plugin' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
```

## plugin-allowed-update

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,acme/util@1.1.0,plugfix/events-plugin@1.0.0
```

## plugin-allowed-update-verbose

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] verbose activate detail
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,acme/util@1.1.0,plugfix/events-plugin@1.0.0
stderr:
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > pre-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > pre-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > pre-autoload-dump: Plugfix\Events\Plugin->onScript
- (oracle) > post-autoload-dump: Plugfix\Events\Plugin->onScript
- (oracle) > post-update-cmd: Plugfix\Events\Plugin->onScript
- (oracle) > post-update-cmd: Plugfix\Events\Plugin->onLate
```

## plugin-allowed-update-debug

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] verbose activate detail
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,acme/util@1.1.0,plugfix/events-plugin@1.0.0
stderr:
- (oracle) Running 2.11.999-dev+source (@release_date@) with PHP 8.4.27 on Linux / 7.1.10-200.fc44.aarch64
- (oracle) Reading ./composer.json (/work/project/composer.json)
- (oracle) Loading config file ./composer.json (/work/project/composer.json)
- (oracle) Checked CA file /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a directory.
- (oracle) Checked CA file /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a directory.
- (oracle) Checked CA file /etc/ssl/certs/ca-certificates.crt: valid
- (oracle) Executing command (CWD): 'git' '--version'
- (oracle) Executing command (/work/project): 'git' 'branch' '-a' '--no-color' '--no-abbrev' '-v'
- (oracle) Executing command (/work/project): 'git' 'describe' '--exact-match' '--tags'
- (oracle) Executing command (/work/project): 'git' 'rev-list' '--no-commit-header' '--format=%H' '-n1' 'HEAD' '--no-show-signature'
- (oracle) Executing command (/work/project): 'hg' 'branch'
- (oracle) Executing command (/work/project): 'fossil' 'branch' 'list'
- (oracle) Executing command (/work/project): 'fossil' 'tag' 'list'
- (oracle) Executing command (/work/project): 'svn' 'info' '--xml'
- (oracle) Failed to initialize global composer: Composer could not find the config file: /work/home/composer.json
- (oracle) 
- (oracle) Dispatching init event
- (oracle) Dispatching pre-command-run (update) event
- (oracle) Dispatching command (update) event
- (oracle) Dispatching pre-file-download event
- (oracle) Downloading http://127.0.0.1:8899/packages.json
- (oracle) [200] http://127.0.0.1:8899/packages.json
- (oracle) Dispatching post-file-download event
- (oracle) Writing /work/cache/repo/http---127.0.0.1-8899/packages.json into cache
- (oracle) Dispatching pre-pool-create event
- (oracle) Built pool in 0.015841960906982s
- (oracle) Running pool optimizer.
- (oracle) Looking at all rules.
- (oracle) Looking at all rules.
- (oracle) Reading ./composer.lock (/work/project/composer.lock)
- (oracle) Dispatching pre-operations-exec event
- (oracle) Dispatching pre-file-download event
- (oracle) Downloading http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip
- (oracle) Dispatching pre-file-download event
- (oracle) Downloading http://127.0.0.1:8899/dists/acme-log-1.1.0.zip
- (oracle) Dispatching pre-file-download event
+ (go) Dispatching pre-package-install event
- (oracle) Downloading http://127.0.0.1:8899/dists/acme-util-1.1.0.zip
- (oracle) [200] http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip
- (oracle) Writing /work/cache/files/plugfix/events-plugin/922758884340d1b085f57f19fc1469d91c4df0ad.zip into cache from /work/project/vendor/composer/tmp-HASH.zip
- (oracle) Dispatching post-file-download event
- (oracle) [200] http://127.0.0.1:8899/dists/acme-log-1.1.0.zip
- (oracle) Writing /work/cache/files/acme/log/31e5146f9094134bf78a7c5baeb8e0a5b2e0e354.zip into cache from /work/project/vendor/composer/tmp-HASH.zip
- (oracle) Dispatching post-file-download event
- (oracle) [200] http://127.0.0.1:8899/dists/acme-util-1.1.0.zip
- (oracle) Writing /work/cache/files/acme/util/b63cd441757e627585620e4283db8491719f9a3a.zip into cache from /work/project/vendor/composer/tmp-HASH.zip
- (oracle) Dispatching post-file-download event
- (oracle) Dispatching pre-package-install (Installing plugfix/events-plugin (1.0.0)) event
+ (go) Dispatching post-package-install event
+ (go) Dispatching pre-package-install event
- (oracle) Executing async command (CWD): '/usr/bin/unzip' '-qq' '/work/project/vendor/composer/tmp-HASH.zip' '-d' '/work/project/vendor/composer/3e877efb'
- (oracle) Executing async command (CWD): 'rm' '-rf' '/work/project/vendor/composer/3e877efb'
- (oracle) Loading plugin Plugfix\Events\Plugin (from plugfix/events-plugin)
- (oracle) Dispatching post-package-install (Installing plugfix/events-plugin (1.0.0)) event
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) Dispatching pre-package-install (Installing acme/log (1.1.0)) event
- (oracle) > pre-package-install: Plugfix\Events\Plugin->onPackage
+ (go) Dispatching pre-package-install event
- (oracle) Executing async command (CWD): '/usr/bin/unzip' '-qq' '/work/project/vendor/composer/tmp-HASH.zip' '-d' '/work/project/vendor/composer/5ea54420'
- (oracle) Dispatching pre-package-install (Installing acme/util (1.1.0)) event
- (oracle) > pre-package-install: Plugfix\Events\Plugin->onPackage
+ (go) Dispatching post-package-install event
+ (go) Dispatching post-package-install event
- (oracle) Executing async command (CWD): '/usr/bin/unzip' '-qq' '/work/project/vendor/composer/tmp-HASH.zip' '-d' '/work/project/vendor/composer/798edb05'
- (oracle) Executing async command (CWD): 'rm' '-rf' '/work/project/vendor/composer/5ea54420'
- (oracle) Executing async command (CWD): 'rm' '-rf' '/work/project/vendor/composer/798edb05'
- (oracle) Dispatching post-package-install (Installing acme/log (1.1.0)) event
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) Dispatching post-package-install (Installing acme/util (1.1.0)) event
- (oracle) > post-package-install: Plugfix\Events\Plugin->onPackage
- (oracle) > pre-autoload-dump: Plugfix\Events\Plugin->onScript
- (oracle) > post-autoload-dump: Plugfix\Events\Plugin->onScript
- (oracle) > post-update-cmd: Plugfix\Events\Plugin->onScript
- (oracle) Executing command (CWD): echo root-post-update
- (oracle) > post-update-cmd: Plugfix\Events\Plugin->onLate
```

## plugin-install-from-lock

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=3
- (oracle) [plugfix/events] post-package-install: install (Installing acme/util (1.1.0)) path=vendor/acme/util exists=true ops=3
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-install-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) root-post-install
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Package operations: 3 installs, 0 updates, 0 removals
- (oracle)   - Downloading plugfix/events-plugin (1.0.0)
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle)   - Downloading acme/util (1.1.0)
- (oracle)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
- (oracle)   - Installing acme/util (1.1.0): Extracting archive
- (oracle) Generating autoload files
- (oracle) > echo root-post-install
project tree:
- (oracle) D vendor
- (oracle) D vendor/acme
- (oracle) D vendor/acme/log
- (oracle) D vendor/acme/log/src
- (oracle) D vendor/acme/util
- (oracle) D vendor/acme/util/lib
- (oracle) D vendor/acme/util/src
- (oracle) D vendor/composer
- (oracle) D vendor/plugfix
- (oracle) D vendor/plugfix/events-plugin
- (oracle) D vendor/plugfix/events-plugin/src
- (oracle) F vendor/acme/log/composer.json 644 b83d42805152f028
- (oracle) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
- (oracle) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
- (oracle) F vendor/acme/util/composer.json 644 ed976ef852108f1f
- (oracle) F vendor/acme/util/lib/legacy.php 644 aa94d8c20c2ac438
- (oracle) F vendor/acme/util/src/Str.php 644 bd3cbab92c3b093f
- (oracle) F vendor/acme/util/src/functions.php 644 f00c9f042b8e0674
- (oracle) F vendor/autoload.php 644 1262ac6dadab7d51
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 d795431bedb6f53e
- (oracle) F vendor/composer/autoload_files.php 644 0e951f0edb6be02b
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 787faed20fbe3acf
- (oracle) F vendor/composer/autoload_real.php 644 c5b3ecdb10011d67
- (oracle) F vendor/composer/autoload_static.php 644 728983b7dbecc048
- (oracle) F vendor/composer/installed.json 644 63e763232c78c45f
- (oracle) F vendor/composer/installed.php 644 58ca9ecaea3fdadf
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
- (oracle) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
- (oracle) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
- (oracle) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
- (oracle) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
autoload probe:
- (oracle) class AcmeAbstractThing: autoloaded
- (oracle) class AcmeSuit: autoloaded
- (oracle) class AcmeThingInterface: autoloaded
- (oracle) class AcmeThingTrait: autoloaded
- (oracle) class Acme\Log\Logger: autoloaded
- (oracle) class Acme\Log\LoggerInterface: autoloaded
- (oracle) class Acme\Util\Str: autoloaded
- (oracle) class Acme_Legacy_Helper: autoloaded
- (oracle) class App\Kernel: autoloaded
- (oracle) Fatal error: Uncaught Error: Interface "Composer\Plugin\Capability\CommandProvider" not found in /work/project/vendor/plugfix/events-plugin/src/CommandProvider.php:4
- (oracle) Stack trace:
- (oracle) #0 /work/project/vendor/composer/ClassLoader.php(576): include()
- (oracle) #1 /work/project/vendor/composer/ClassLoader.php(427): {closure:Composer\Autoload\ClassLoader::initializeIncludeClosure():575}('/work/project/v...')
- (oracle) #2 [internal function]: Composer\Autoload\ClassLoader->loadClass('Plugfix\\Events\\...')
- (oracle) #3 /gobin/probe.php(42): class_exists('Plugfix\\Events\\...')
- (oracle) #4 {main}
- (oracle)   thrown in /work/project/vendor/plugfix/events-plugin/src/CommandProvider.php on line 4
- (oracle) 
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'AcmeAbstractThing' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeSuit' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeThingInterface' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'AcmeThingTrait' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'Acme_Legacy_Helper' => $vendorDir . '/acme/util/lib/legacy.php',
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle) );
- (oracle) 
vendor/composer/autoload_files.php:
- (oracle) <?php
- (oracle) // autoload_files.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     '571dbb8238a79d30de0279182c1a3513' => $vendorDir . '/acme/util/src/functions.php',
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
- (oracle)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle)     'Acme\\Util\\' => array($vendorDir . '/acme/util/src'),
- (oracle)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         $filesToLoad = \Composer\Autoload\ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$files;
- (oracle)         $requireFile = \Closure::bind(static function ($fileIdentifier, $file) {
- (oracle)             if (empty($GLOBALS['__composer_autoload_files'][$fileIdentifier])) {
- (oracle)                 $GLOBALS['__composer_autoload_files'][$fileIdentifier] = true;
- (oracle) 
- (oracle)                 require $file;
- (oracle)             }
- (oracle)         }, null, null);
- (oracle)         foreach ($filesToLoad as $fileIdentifier => $file) {
- (oracle)             $requireFile($fileIdentifier, $file);
- (oracle)         }
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
- (oracle) class ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e
- (oracle) {
- (oracle)     public static $files = array (
- (oracle)         '571dbb8238a79d30de0279182c1a3513' => __DIR__ . '/..' . '/acme/util/src/functions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'P' =>
- (oracle)         array (
- (oracle)             'Plugfix\\Events\\' => 15,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)             'Acme\\Util\\' => 10,
- (oracle)             'Acme\\Log\\' => 9,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'Plugfix\\Events\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)         'Acme\\Util\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/util/src',
- (oracle)         ),
- (oracle)         'Acme\\Log\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'AcmeAbstractThing' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeSuit' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeThingInterface' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'AcmeThingTrait' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'Acme_Legacy_Helper' => __DIR__ . '/..' . '/acme/util/lib/legacy.php',
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "acme/log",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip",
- (oracle)                 "reference": "58b0f69fc5f96981e5f7fae32d592232a36012c2",
- (oracle)                 "shasum": "c40f1901ba3e909555618861c594e339cd678f39"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "php": ">=7.4"
- (oracle)             },
- (oracle)             "suggest": {
- (oracle)                 "acme/extra": "for extra stuff"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\Log\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Logging",
- (oracle)             "install-path": "../acme/log"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "acme/util",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-util-1.1.0.zip",
- (oracle)                 "reference": "010715fb4511de3a64a1da0e47aa42e3b0104c1c",
- (oracle)                 "shasum": "666855ca01fe54cab44e7fa49212e73bf4193aee"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "acme/log": "^1.1",
- (oracle)                 "ext-json": "*"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\Util\\": "src/"
- (oracle)                 },
- (oracle)                 "files": [
- (oracle)                     "src/functions.php"
- (oracle)                 ],
- (oracle)                 "classmap": [
- (oracle)                     "lib/"
- (oracle)                 ]
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Utilities",
- (oracle)             "install-path": "../acme/util"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/events-plugin",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
- (oracle)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
- (oracle)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "composer-plugin-api": "^2.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "composer-plugin",
- (oracle)             "extra": {
- (oracle)                 "class": "Plugfix\\Events\\Plugin"
- (oracle)             },
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Plugfix\\Events\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture plugin subscribing to events",
- (oracle)             "install-path": "../plugfix/events-plugin"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": [
- (oracle)         "acme/util"
- (oracle)     ]
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
- (oracle)         'acme/log' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/log',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'acme/util' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '010715fb4511de3a64a1da0e47aa42e3b0104c1c',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/util',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => true,
- (oracle)         ),
- (oracle)         'plugfix/events-plugin' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
- (oracle)             'type' => 'composer-plugin',
- (oracle)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
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

## plugin-install-no-dev

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=2
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] pre-autoload-dump dev=false count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=false count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-install-cmd dev=false count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) root-post-install
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Package operations: 2 installs, 0 updates, 0 removals
- (oracle)   - Downloading plugfix/events-plugin (1.0.0)
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle)   - Installing plugfix/events-plugin (1.0.0): Extracting archive
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
- (oracle) Generating autoload files
- (oracle) > echo root-post-install
project tree:
- (oracle) D vendor
- (oracle) D vendor/acme
- (oracle) D vendor/acme/log
- (oracle) D vendor/acme/log/src
- (oracle) D vendor/composer
- (oracle) D vendor/plugfix
- (oracle) D vendor/plugfix/events-plugin
- (oracle) D vendor/plugfix/events-plugin/src
- (oracle) F vendor/acme/log/composer.json 644 b83d42805152f028
- (oracle) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
- (oracle) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
- (oracle) F vendor/autoload.php 644 1262ac6dadab7d51
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 0635e595c1470c97
- (oracle) F vendor/composer/autoload_real.php 644 8e3def21c1292c9a
- (oracle) F vendor/composer/autoload_static.php 644 b08be4f4dd861d9e
- (oracle) F vendor/composer/installed.json 644 6ddb47fb66056d77
- (oracle) F vendor/composer/installed.php 644 e8c387e57c1328cd
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
- (oracle) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
- (oracle) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
- (oracle) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
- (oracle) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
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
- (oracle)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit8acf17641576c94ceca3fc9c9828ca9e', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::getInitializer($loader));
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
- (oracle) class ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'P' =>
- (oracle)         array (
- (oracle)             'Plugfix\\Events\\' => 15,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)             'Acme\\Log\\' => 9,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'Plugfix\\Events\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)         'Acme\\Log\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log/src',
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
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit8acf17641576c94ceca3fc9c9828ca9e::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "acme/log",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip",
- (oracle)                 "reference": "58b0f69fc5f96981e5f7fae32d592232a36012c2",
- (oracle)                 "shasum": "c40f1901ba3e909555618861c594e339cd678f39"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "php": ">=7.4"
- (oracle)             },
- (oracle)             "suggest": {
- (oracle)                 "acme/extra": "for extra stuff"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\Log\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Logging",
- (oracle)             "install-path": "../acme/log"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/events-plugin",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
- (oracle)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
- (oracle)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "composer-plugin-api": "^2.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "composer-plugin",
- (oracle)             "extra": {
- (oracle)                 "class": "Plugfix\\Events\\Plugin"
- (oracle)             },
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Plugfix\\Events\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture plugin subscribing to events",
- (oracle)             "install-path": "../plugfix/events-plugin"
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": false,
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
- (oracle)         'dev' => false,
- (oracle)     ),
- (oracle)     'versions' => array(
- (oracle)         'acme/log' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/log',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'plugfix/events-plugin' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
- (oracle)             'type' => 'composer-plugin',
- (oracle)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
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

## plugin-installed-noop

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run install
- (oracle) [plugfix/events] command install dry-run=false
- (oracle) [plugfix/events] pre-install-cmd dev=true count=1 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] pre-operations-exec exec=true dev=true 
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=3 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-install-cmd dev=true count=4 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) root-post-install
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Nothing to install, update or remove
- (oracle) Generating autoload files
- (oracle) > echo root-post-install
```

## plugin-installed-quiet

```
exit code: oracle=0 go=1
stderr:
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
+ (go) 
```

## plugin-no-scripts

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run install
- (oracle) [plugfix/events] command install dry-run=false
- (oracle) [plugfix/events] pre-install-cmd dev=true count=1 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] pre-operations-exec exec=true dev=true 
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=3 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-install-cmd dev=true count=4 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Nothing to install, update or remove
- (oracle) Generating autoload files
```

## plugin-dump-autoload

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run dump-autoload
- (oracle) [plugfix/events] command dump-autoload dry-run=false
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
```

## plugin-dump-autoload-optimize

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run dump-autoload
- (oracle) [plugfix/events] command dump-autoload dry-run=false
- (oracle) [plugfix/events] pre-autoload-dump dev=false count=1 flags={"optimize":true} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=false count=2 flags={"optimize":true} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
```

## plugin-no-plugins-update

```
stderr:
- (oracle) The "plugfix/events-plugin" plugin was not loaded as plugins are disabled.
```

## plugin-update-version

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run update
- (oracle) [plugfix/events] command update dry-run=false
- (oracle) [plugfix/events] pre-update-cmd dev=true count=1 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] pre-operations-exec exec=true dev=true update:Upgrading plugfix/events-plugin (1.0.0 => 1.1.0)
- (oracle) [plugfix/events] pre-package-update: update (Upgrading plugfix/events-plugin (1.0.0 => 1.1.0)) path=vendor/plugfix/events-plugin exists=true ops=1
- (oracle) [plugfix/events] deactivate 1.0.0
- (oracle) [plugfix/events] activate 1.1.0
- (oracle) [plugfix/events] post-package-update: update (Upgrading plugfix/events-plugin (1.0.0 => 1.1.0)) path=vendor/plugfix/events-plugin exists=true ops=1
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) root-post-update
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,acme/util@1.1.0,plugfix/events-plugin@1.1.0
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 0 installs, 1 update, 0 removals
- (oracle)   - Upgrading plugfix/events-plugin (1.0.0 => 1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 0 installs, 1 update, 0 removals
- (oracle)   - Downloading plugfix/events-plugin (1.1.0)
- (oracle)   - Upgrading plugfix/events-plugin (1.0.0 => 1.1.0): Extracting archive
- (oracle) Generating autoload files
- (oracle) > @putenv PLUGFIX_X=1
- (oracle) > echo root-post-update
composer.lock:
+ (go)     "content-hash": "8acf17641576c94ceca3fc9c9828ca9e",
- (oracle)     "content-hash": "48d2cb2825a5b1e087d1127ec08cb9d7",
+ (go)             "version": "1.0.0",
- (oracle)             "version": "1.1.0",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.1.0.zip",
- (oracle)                 "reference": "a671cf94a114b1117a53605b378900c558353701",
- (oracle)                 "shasum": "ab85a0bcc7b858e33d23b505a77292891717a935"
project tree:
- (oracle) D vendor/bin
+ (go) F composer.lock 644 17fe9c22126f0126
- (oracle) F composer.lock 644 d837fbd1db014184
+ (go) F vendor/composer/installed.json 644 63e763232c78c45f
+ (go) F vendor/composer/installed.php 644 58ca9ecaea3fdadf
- (oracle) F vendor/composer/installed.json 644 17d96301b5c52adc
- (oracle) F vendor/composer/installed.php 644 f3fbc6568a86a375
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
- (oracle) F vendor/plugfix/events-plugin/src/Plugin.php 644 3a7be4f9daea7d35
vendor/composer/installed.json:
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.1.0.zip",
- (oracle)                 "reference": "a671cf94a114b1117a53605b378900c558353701",
- (oracle)                 "shasum": "ab85a0bcc7b858e33d23b505a77292891717a935"
vendor/composer/installed.php:
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => 'a671cf94a114b1117a53605b378900c558353701',
```

## plugin-remove

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run remove
- (oracle) [plugfix/events] deactivate 1.0.0
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] command remove dry-run=false
- (oracle) [plugfix/events] pre-update-cmd dev=true count=1 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] pre-operations-exec exec=true dev=true uninstall:Removing plugfix/events-plugin (1.0.0)
- (oracle) [plugfix/events] pre-package-uninstall: uninstall (Removing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=1
- (oracle) [plugfix/events] deactivate 1.0.0
- (oracle) [plugfix/events] uninstall 1.0.0
- (oracle) root-post-update
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) remove [--dev] [--dry-run] [--no-progress] [--no-update] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--update-no-dev] [-w|--update-with-dependencies] [-W|--update-with-all-dependencies] [--with-all-dependencies] [--no-update-with-dependencies] [-m|--minimal-changes] [--unused] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--] [<packages>...]
+ (go) 
- (oracle) Running composer update plugfix/events-plugin
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 0 installs, 0 updates, 1 removal
- (oracle)   - Removing plugfix/events-plugin (1.0.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 0 installs, 0 updates, 1 removal
- (oracle)   - Removing plugfix/events-plugin (1.0.0)
- (oracle) Generating autoload files
- (oracle) > @putenv PLUGFIX_X=1
- (oracle) > echo root-post-update
composer.lock:
+ (go)     "content-hash": "8acf17641576c94ceca3fc9c9828ca9e",
- (oracle)     "content-hash": "a26f612fc650ab601e9efc55d57ba2b5",
+ (go)         },
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "time": "2024-01-02T03:04:05+00:00"
project tree:
- (oracle) D vendor/bin
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/events-plugin
+ (go) D vendor/plugfix/events-plugin/src
+ (go) F composer.lock 644 17fe9c22126f0126
- (oracle) F composer.lock 644 1299b1766444b5f7
+ (go) F vendor/composer/autoload_psr4.php 644 787faed20fbe3acf
- (oracle) F vendor/composer/autoload_psr4.php 644 494862b89f0088b2
+ (go) F vendor/composer/autoload_static.php 644 728983b7dbecc048
+ (go) F vendor/composer/installed.json 644 63e763232c78c45f
+ (go) F vendor/composer/installed.php 644 58ca9ecaea3fdadf
- (oracle) F vendor/composer/autoload_static.php 644 d66a65bfdff508c1
- (oracle) F vendor/composer/installed.json 644 e789f9f1a5d5d793
- (oracle) F vendor/composer/installed.php 644 cd911c910d26559f
+ (go) F vendor/plugfix/events-plugin/composer.json 644 8bea6ab006d8c822
+ (go) F vendor/plugfix/events-plugin/src/CommandProvider.php 644 491b7467d2361913
+ (go) F vendor/plugfix/events-plugin/src/HelloCommand.php 644 bd6acffd1ea24c2f
+ (go) F vendor/plugfix/events-plugin/src/Plugin.php 644 119b2173b85d09c9
vendor/composer/autoload_psr4.php:
+ (go)     'Plugfix\\Events\\' => array($vendorDir . '/plugfix/events-plugin/src'),
vendor/composer/autoload_static.php:
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\Events\\' => 15,
+ (go)         ),
+ (go)         'Plugfix\\Events\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/events-plugin/src',
+ (go)         ),
vendor/composer/installed.json:
+ (go)         },
+ (go)         {
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "install-path": "../plugfix/events-plugin"
vendor/composer/installed.php:
+ (go)         'plugfix/events-plugin' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '6b533a89a79c1d22e1063facddfd920f480e19b1',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/events-plugin',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
```

## plugin-require

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] post-package-install: install (Installing plugfix/events-plugin (1.0.0)) path=vendor/plugfix/events-plugin exists=true ops=2
- (oracle) [plugfix/events] pre-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] post-package-install: install (Installing acme/log (1.1.0)) path=vendor/acme/log exists=true ops=2
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=1 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=3 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,plugfix/events-plugin@1.0.0
```

## plugin-require-installed

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run require
- (oracle) [plugfix/events] deactivate 1.0.0
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] command require dry-run=false
- (oracle) [plugfix/events] pre-update-cmd dev=true count=1 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] pre-operations-exec exec=true dev=true 
- (oracle) [plugfix/events] pre-autoload-dump dev=true count=2 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-autoload-dump dev=true count=3 flags={"optimize":false} args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) [plugfix/events] post-update-cmd dev=true count=4 flags=[] args=[]
- (oracle) [plugfix/events] root test/project 1.0.0+no-version-set vendor-dir=vendor
- (oracle) root-post-update
- (oracle) [plugfix/events] late post-update-cmd installed=acme/log@1.1.0,acme/util@1.1.0,plugfix/events-plugin@1.0.0
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) Installation failed, reverting ./composer.json and ./composer.lock to their original content.
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/events-plugin contains a Composer plugin which is allowed by your a  
+ (go)   llow-plugins config, but gocomposer cannot execute Composer plugins. Run wi  
+ (go)   th --no-plugins to install without activating plugins.                       
+ (go)                                                                                
+ (go) 
+ (go) require [--dev] [--dry-run] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--fixed] [--no-suggest] [--no-progress] [--no-update] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--update-no-dev] [-w|--update-with-dependencies] [-W|--update-with-all-dependencies] [--with-dependencies] [--with-all-dependencies] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--sort-packages] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Nothing to modify in lock file
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Nothing to install, update or remove
- (oracle) Generating autoload files
- (oracle) > @putenv PLUGFIX_X=1
- (oracle) > echo root-post-update
composer.lock:
+ (go)     "content-hash": "8acf17641576c94ceca3fc9c9828ca9e",
- (oracle)     "content-hash": "161135a5cdaa00449c6fd30db2946261",
+ (go)             "name": "plugfix/events-plugin",
+ (go)             "version": "1.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
+ (go)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
+ (go)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\Events\\Plugin"
+ (go)             },
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\Events\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture plugin subscribing to events",
+ (go)             "time": "2024-01-02T03:04:05+00:00"
+ (go)         }
+ (go)     ],
+ (go)     "packages-dev": [
+ (go)         {
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/events-plugin",
- (oracle)             "version": "1.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-events-plugin-1.0.0.zip",
- (oracle)                 "reference": "6b533a89a79c1d22e1063facddfd920f480e19b1",
- (oracle)                 "shasum": "7a1aee28d76ec742d7a0c06c96acb2c38720ccf5"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "composer-plugin-api": "^2.0"
- (oracle)             },
- (oracle)             "type": "composer-plugin",
- (oracle)             "extra": {
- (oracle)                 "class": "Plugfix\\Events\\Plugin"
- (oracle)             },
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Plugfix\\Events\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture plugin subscribing to events",
- (oracle)             "time": "2024-01-02T03:04:05+00:00"
- (oracle)     "packages-dev": [],
project tree:
+ (go) F composer.json 644 22bf21e32a689a84
+ (go) F composer.lock 644 17fe9c22126f0126
- (oracle) F composer.json 644 a368e08bcadefbb8
- (oracle) F composer.lock 644 ee906ccc04f00a4f
+ (go) F vendor/composer/installed.json 644 63e763232c78c45f
+ (go) F vendor/composer/installed.php 644 58ca9ecaea3fdadf
- (oracle) F vendor/composer/installed.json 644 f93203af6d84986d
- (oracle) F vendor/composer/installed.php 644 98685ae95f18ab9e
composer.json:
+ (go)     "require": {"plugfix/events-plugin": "1.0.0", "acme/log": "^1.0"},
+ (go)     "require-dev": {"acme/util": "^1.0"},
- (oracle)     "require": {"plugfix/events-plugin": "1.0.0", "acme/log": "^1.0",
- (oracle)         "acme/util": "^1.0"},
vendor/composer/installed.json:
+ (go)     "dev-package-names": [
+ (go)         "acme/util"
+ (go)     ]
- (oracle)     "dev-package-names": []
vendor/composer/installed.php:
+ (go)             'dev_requirement' => true,
- (oracle)             'dev_requirement' => false,
```

## plugin-run-script-custom

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run run-script
- (oracle) [plugfix/events] custom ["a","b"]
```

## plugin-show

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run show
- (oracle) [plugfix/events] command show dry-run=NULL
```

## plugin-thrower-boom

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking plugfix/thrower (1.0.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Downloading plugfix/thrower (1.0.0)
  - Installing plugfix/thrower (1.0.0): Extracting archive
Generating autoload files

In Plugin.php line 32:
                                      
  plugfix boom in post-autoload-dump  
                                      

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## plugin-thrower-stop

```
stdout:
- (oracle) [plugfix/thrower] stop post-autoload-dump
```

## plugin-thrower-false

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/thrower] stop post-autoload-dump
- (oracle) [plugfix/thrower] no boom
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/thrower contains a Composer plugin which is allowed by your allow-p  
+ (go)   lugins config, but gocomposer cannot execute Composer plugins. Run with --n  
+ (go)   o-plugins to install without activating plugins.                             
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Nothing to install, update or remove
- (oracle) Generating autoload files
```

## plugin-no-api

```
exit code: oracle=1 go=0
stdout:
+ (go) No security vulnerability advisories found.
stderr:
+ (go) Generating autoload files
- (oracle) Plugin initialization failed (Plugin plugfix/no-api is missing a require statement for a version of the composer-plugin-api package.), uninstalling plugin
- (oracle)   - Removing plugfix/no-api (1.0.0)
- (oracle)     Install of plugfix/no-api failed
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Plugin plugfix/no-api is missing a require statement for a version of the c  
- (oracle)   omposer-plugin-api package.                                                  
- (oracle)                                                                                
- (oracle) 
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/no-api
+ (go) D vendor/plugfix/no-api/src
- (oracle) D vendor/bin
+ (go) F vendor/autoload.php 644 a490ec5f52e52017
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 d9321e0c90ae332e
+ (go) F vendor/composer/autoload_real.php 644 2395ad2215df18a1
+ (go) F vendor/composer/autoload_static.php 644 4956b247241df6cb
+ (go) F vendor/composer/installed.json 644 c9ee2ff849cbe60a
+ (go) F vendor/composer/installed.php 644 3e50b64c7a6a8aeb
+ (go) F vendor/plugfix/no-api/composer.json 644 072424a7bbd350a5
+ (go) F vendor/plugfix/no-api/src/Plugin.php 644 ca26e9b8bcb5e803
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Plugfix\\NoApi\\' => array($vendorDir . '/plugfix/no-api/src'),
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInitb8d80a5b99322d360d2fc8e6ea4bb56e
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInitb8d80a5b99322d360d2fc8e6ea4bb56e', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInitb8d80a5b99322d360d2fc8e6ea4bb56e', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInitb8d80a5b99322d360d2fc8e6ea4bb56e::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInitb8d80a5b99322d360d2fc8e6ea4bb56e
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'P' =>
+ (go)         array (
+ (go)             'Plugfix\\NoApi\\' => 14,
+ (go)         ),
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'Plugfix\\NoApi\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/..' . '/plugfix/no-api/src',
+ (go)         ),
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInitb8d80a5b99322d360d2fc8e6ea4bb56e::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInitb8d80a5b99322d360d2fc8e6ea4bb56e::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInitb8d80a5b99322d360d2fc8e6ea4bb56e::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "plugfix/no-api",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-no-api-1.0.0.zip",
+ (go)                 "reference": "e8e2752b86d97a79ddfae658f61bbadb5af084d5",
+ (go)                 "shasum": "b1e71157f92c3611f293928a82051f550eb9d847"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "composer-plugin",
+ (go)             "extra": {
+ (go)                 "class": "Plugfix\\NoApi\\Plugin"
+ (go)             },
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\NoApi\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Plugin without a composer-plugin-api requirement",
+ (go)             "install-path": "../plugfix/no-api"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'plugfix/no-api' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => 'e8e2752b86d97a79ddfae658f61bbadb5af084d5',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/no-api',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
```

## plugin-no-class

```
exit code: oracle=1 go=0
stdout:
+ (go) No security vulnerability advisories found.
stderr:
+ (go)   - Downloading plugfix/no-class (1.0.0)
+ (go)   - Installing plugfix/no-class (1.0.0): Extracting archive
+ (go) Generating autoload files
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Error while installing plugfix/no-class, composer-plugin packages should ha  
- (oracle)   ve a class defined in their extra key to be usable.                          
- (oracle)                                                                                
- (oracle) 
- (oracle) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
- (oracle) 
- (oracle) 
project tree:
+ (go) D vendor
+ (go) D vendor/composer
+ (go) D vendor/plugfix
+ (go) D vendor/plugfix/no-class
+ (go) F vendor/autoload.php 644 0c9e14a4364432b0
+ (go) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
+ (go) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
+ (go) F vendor/composer/LICENSE 644 7855ac293067aebe
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
+ (go) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
+ (go) F vendor/composer/autoload_psr4.php 644 c35c405f584d1f28
+ (go) F vendor/composer/autoload_real.php 644 2b0b4d09ed3a236c
+ (go) F vendor/composer/autoload_static.php 644 2c4237d84430dc68
+ (go) F vendor/composer/installed.json 644 90989c0209eb93d3
+ (go) F vendor/composer/installed.php 644 af110557dfc83730
+ (go) F vendor/plugfix/no-class/composer.json 644 eaf85e23c40677cf
vendor/composer/autoload_classmap.php:
+ (go) <?php
+ (go) // autoload_classmap.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
+ (go) );
+ (go) 
vendor/composer/autoload_psr4.php:
+ (go) <?php
+ (go) // autoload_psr4.php @generated by Composer
+ (go) 
+ (go) $vendorDir = dirname(__DIR__);
+ (go) $baseDir = dirname($vendorDir);
+ (go) 
+ (go) return array(
+ (go)     'App\\' => array($baseDir . '/src'),
+ (go) );
+ (go) 
vendor/composer/autoload_real.php:
+ (go) <?php
+ (go) // autoload_real.php @generated by Composer
+ (go) 
+ (go) class ComposerAutoloaderInitb1f2561fd2c0862a42432e4f878d9da9
+ (go) {
+ (go)     private static $loader;
+ (go) 
+ (go)     public static function loadClassLoader($class)
+ (go)     {
+ (go)         if ('Composer\Autoload\ClassLoader' === $class) {
+ (go)             require __DIR__ . '/ClassLoader.php';
+ (go)         }
+ (go)     }
+ (go) 
+ (go)     /**
+ (go)      * @return \Composer\Autoload\ClassLoader
+ (go)      */
+ (go)     public static function getLoader()
+ (go)     {
+ (go)         if (null !== self::$loader) {
+ (go)             return self::$loader;
+ (go)         }
+ (go) 
+ (go)         spl_autoload_register(array('ComposerAutoloaderInitb1f2561fd2c0862a42432e4f878d9da9', 'loadClassLoader'), true, true);
+ (go)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInitb1f2561fd2c0862a42432e4f878d9da9', 'loadClassLoader'));
+ (go) 
+ (go)         require __DIR__ . '/autoload_static.php';
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInitb1f2561fd2c0862a42432e4f878d9da9::getInitializer($loader));
+ (go) 
+ (go)         $loader->register(true);
+ (go) 
+ (go)         return $loader;
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/autoload_static.php:
+ (go) <?php
+ (go) // autoload_static.php @generated by Composer
+ (go) 
+ (go) namespace Composer\Autoload;
+ (go) 
+ (go) class ComposerStaticInitb1f2561fd2c0862a42432e4f878d9da9
+ (go) {
+ (go)     public static $prefixLengthsPsr4 = array (
+ (go)         'A' =>
+ (go)         array (
+ (go)             'App\\' => 4,
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $prefixDirsPsr4 = array (
+ (go)         'App\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/src',
+ (go)         ),
+ (go)     );
+ (go) 
+ (go)     public static $classMap = array (
+ (go)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
+ (go)     );
+ (go) 
+ (go)     public static function getInitializer(ClassLoader $loader)
+ (go)     {
+ (go)         return \Closure::bind(function () use ($loader) {
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInitb1f2561fd2c0862a42432e4f878d9da9::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInitb1f2561fd2c0862a42432e4f878d9da9::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInitb1f2561fd2c0862a42432e4f878d9da9::$classMap;
+ (go) 
+ (go)         }, null, ClassLoader::class);
+ (go)     }
+ (go) }
+ (go) 
vendor/composer/installed.json:
+ (go) {
+ (go)     "packages": [
+ (go)         {
+ (go)             "name": "plugfix/no-class",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-no-class-1.0.0.zip",
+ (go)                 "reference": "d415761d27aa67736e8807c25fb180c5c712a21c",
+ (go)                 "shasum": "3402f4235ab2cd20473b0b6a1d27de7bbc9d1e41"
+ (go)             },
+ (go)             "require": {
+ (go)                 "composer-plugin-api": "^2.0"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "type": "composer-plugin",
+ (go)             "installation-source": "dist",
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Plugin without extra.class",
+ (go)             "install-path": "../plugfix/no-class"
+ (go)         }
+ (go)     ],
+ (go)     "dev": true,
+ (go)     "dev-package-names": []
+ (go) }
vendor/composer/installed.php:
+ (go) <?php return array(
+ (go)     'root' => array(
+ (go)         'name' => 'test/project',
+ (go)         'pretty_version' => '1.0.0+no-version-set',
+ (go)         'version' => '1.0.0.0',
+ (go)         'reference' => null,
+ (go)         'type' => 'project',
+ (go)         'install_path' => __DIR__ . '/../../',
+ (go)         'aliases' => array(),
+ (go)         'dev' => true,
+ (go)     ),
+ (go)     'versions' => array(
+ (go)         'plugfix/no-class' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => 'd415761d27aa67736e8807c25fb180c5c712a21c',
+ (go)             'type' => 'composer-plugin',
+ (go)             'install_path' => __DIR__ . '/../plugfix/no-class',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)         'test/project' => array(
+ (go)             'pretty_version' => '1.0.0+no-version-set',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => null,
+ (go)             'type' => 'project',
+ (go)             'install_path' => __DIR__ . '/../../',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
+ (go)     ),
+ (go) );
```

## plugin-legacy-installer

```
stderr:
- (oracle) Loading "plugfix/legacy" which is a legacy composer-installer built for Composer 1.x, it is likely to cause issues as you are running Composer 2.x.
```

## plugin-command

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run plugfix:hello
- (oracle) HELLO BOB!
- (oracle) HELLO BOB!
stderr:
+ (go) 
+ (go)                                            
+ (go)   Command "plugfix:hello" is not defined.  
+ (go)                                            
+ (go) 
- (oracle) root=test/project plugins=1
```

## plugin-command-default

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] verbose activate detail
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run plugfix:hello
- (oracle) Hello world
- (oracle) verbose line
stderr:
+ (go) 
+ (go)                                            
+ (go)   Command "plugfix:hello" is not defined.  
+ (go)                                            
+ (go) 
- (oracle) > init: Plugfix\Events\Plugin->onInit
- (oracle) > pre-command-run: Plugfix\Events\Plugin->onPreCommandRun
- (oracle) root=test/project plugins=1
```

## plugin-command-exit

```
exit code: oracle=3 go=1
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
- (oracle) [plugfix/events] pre-command-run plugfix:hello
- (oracle) Hello world
- (oracle) Hello world
- (oracle) Hello world
stderr:
+ (go) 
+ (go)                                            
+ (go)   Command "plugfix:hello" is not defined.  
+ (go)                                            
+ (go) 
- (oracle) root=test/project plugins=1
```

## plugin-command-bad-option

```
stdout:
- (oracle) [plugfix/events] activate 1.0.0
- (oracle) [plugfix/events] init Composer\EventDispatcher\Event
stderr:
+ (go)                                            
+ (go)   Command "plugfix:hello" is not defined.  
+ (go)                                            
- (oracle)                                        
- (oracle)   The "--nope" option does not exist.  
- (oracle)                                        
- (oracle) plugfix:hello [--shout] [-t|--times TIMES] [--] [<who>]
- (oracle) 
```

## plugin-command-no-plugins

```
stderr:
- (oracle) Plugins have been disabled, which may be why some commands are missing, unless you made a typo
+ (go)                                            
+ (go)   Command "plugfix:hello" is not defined.  
+ (go)                                            
- (oracle)                                                              
- (oracle)   There are no commands defined in the "plugfix" namespace.  
- (oracle)                                                              
```

## plugin-installer-update

```
stdout:
- (oracle) [plugfix/installer] installed plugfix/module-b into modules/bee (present)
- (oracle) [plugfix/installer] installed plugfix/module-a into modules/module-a (present)
project tree:
- (oracle) D modules
- (oracle) D modules/bee
- (oracle) D modules/bee/src
- (oracle) D modules/module-a
- (oracle) D modules/module-a/bin
- (oracle) D modules/module-a/src
+ (go) D vendor/plugfix/module-a
+ (go) D vendor/plugfix/module-a/bin
+ (go) D vendor/plugfix/module-a/src
+ (go) D vendor/plugfix/module-b
+ (go) D vendor/plugfix/module-b/src
- (oracle) F modules/bee/composer.json 644 4d2fc4f52472a068
- (oracle) F modules/bee/src/Bee.php 644 95a279ce42abc9be
- (oracle) F modules/module-a/bin/module-a 755 497f6811dd69c6c6
- (oracle) F modules/module-a/composer.json 644 42213e55ec719364
- (oracle) F modules/module-a/src/Module.php 644 5fb3895bcb248f97
+ (go) F vendor/bin/module-a 755 1ff37d44f82d39b0
- (oracle) F vendor/bin/module-a 755 3adc06b2b2fd7aad
+ (go) F vendor/composer/autoload_classmap.php 644 7e4440b9eb4ba974
- (oracle) F vendor/composer/autoload_classmap.php 644 c11bdea06020b395
+ (go) F vendor/composer/autoload_psr4.php 644 f71f006bbf58f79b
- (oracle) F vendor/composer/autoload_psr4.php 644 fd5a0a19865aceb2
+ (go) F vendor/composer/autoload_static.php 644 81e21c8102a302d3
+ (go) F vendor/composer/installed.json 644 774804a8e83bb9ed
+ (go) F vendor/composer/installed.php 644 291c79ec46b39aba
- (oracle) F vendor/composer/autoload_static.php 644 67977963a6de732c
- (oracle) F vendor/composer/installed.json 644 e3809b62ec53d12c
- (oracle) F vendor/composer/installed.php 644 d1384fc76d15b4a8
+ (go) F vendor/plugfix/module-a/bin/module-a 755 497f6811dd69c6c6
+ (go) F vendor/plugfix/module-a/composer.json 644 42213e55ec719364
+ (go) F vendor/plugfix/module-a/src/Module.php 644 5fb3895bcb248f97
+ (go) F vendor/plugfix/module-b/composer.json 644 4d2fc4f52472a068
+ (go) F vendor/plugfix/module-b/src/Bee.php 644 95a279ce42abc9be
vendor/composer/autoload_classmap.php:
+ (go)     'Plugfix\\ModuleB\\Bee' => $vendorDir . '/plugfix/module-b/src/Bee.php',
- (oracle)     'Plugfix\\ModuleB\\Bee' => $baseDir . '/modules/bee/src/Bee.php',
vendor/composer/autoload_psr4.php:
+ (go)     'Plugfix\\ModuleA\\' => array($vendorDir . '/plugfix/module-a/src'),
- (oracle)     'Plugfix\\ModuleA\\' => array($baseDir . '/modules/module-a/src'),
vendor/composer/autoload_static.php:
+ (go)             0 => __DIR__ . '/..' . '/plugfix/module-a/src',
- (oracle)             0 => __DIR__ . '/../..' . '/modules/module-a/src',
+ (go)         'Plugfix\\ModuleB\\Bee' => __DIR__ . '/..' . '/plugfix/module-b/src/Bee.php',
- (oracle)         'Plugfix\\ModuleB\\Bee' => __DIR__ . '/../..' . '/modules/bee/src/Bee.php',
vendor/composer/installed.json:
+ (go)             "install-path": "../plugfix/module-a"
- (oracle)             "install-path": "../../modules/module-a"
+ (go)             "install-path": "../plugfix/module-b"
- (oracle)             "install-path": "../../modules/bee"
vendor/composer/installed.php:
+ (go)             'install_path' => __DIR__ . '/../plugfix/module-a',
- (oracle)             'install_path' => __DIR__ . '/../../modules/module-a',
+ (go)             'install_path' => __DIR__ . '/../plugfix/module-b',
- (oracle)             'install_path' => __DIR__ . '/../../modules/bee',
```

## plugin-installer-extra-path

```
stdout:
- (oracle) [plugfix/installer] installed plugfix/module-b into app/mods/bee (present)
- (oracle) [plugfix/installer] installed plugfix/module-a into app/mods/module-a (present)
project tree:
- (oracle) D app
- (oracle) D app/mods
- (oracle) D app/mods/bee
- (oracle) D app/mods/bee/src
- (oracle) D app/mods/module-a
- (oracle) D app/mods/module-a/bin
- (oracle) D app/mods/module-a/src
+ (go) D vendor/plugfix/module-a
+ (go) D vendor/plugfix/module-a/bin
+ (go) D vendor/plugfix/module-a/src
+ (go) D vendor/plugfix/module-b
+ (go) D vendor/plugfix/module-b/src
- (oracle) F app/mods/bee/composer.json 644 4d2fc4f52472a068
- (oracle) F app/mods/bee/src/Bee.php 644 95a279ce42abc9be
- (oracle) F app/mods/module-a/bin/module-a 755 497f6811dd69c6c6
- (oracle) F app/mods/module-a/composer.json 644 42213e55ec719364
- (oracle) F app/mods/module-a/src/Module.php 644 5fb3895bcb248f97
+ (go) F vendor/bin/module-a 755 1ff37d44f82d39b0
- (oracle) F vendor/bin/module-a 755 8e7ad711bd91ba99
+ (go) F vendor/composer/autoload_classmap.php 644 7e4440b9eb4ba974
- (oracle) F vendor/composer/autoload_classmap.php 644 4cdb763ed6027135
+ (go) F vendor/composer/autoload_psr4.php 644 f71f006bbf58f79b
- (oracle) F vendor/composer/autoload_psr4.php 644 dc58efa9fdf5277c
+ (go) F vendor/composer/autoload_static.php 644 84e2752267b2f433
+ (go) F vendor/composer/installed.json 644 774804a8e83bb9ed
+ (go) F vendor/composer/installed.php 644 291c79ec46b39aba
- (oracle) F vendor/composer/autoload_static.php 644 072201add50227e5
- (oracle) F vendor/composer/installed.json 644 a67ab9c6acb9db65
- (oracle) F vendor/composer/installed.php 644 257bb6ff9787dba2
+ (go) F vendor/plugfix/module-a/bin/module-a 755 497f6811dd69c6c6
+ (go) F vendor/plugfix/module-a/composer.json 644 42213e55ec719364
+ (go) F vendor/plugfix/module-a/src/Module.php 644 5fb3895bcb248f97
+ (go) F vendor/plugfix/module-b/composer.json 644 4d2fc4f52472a068
+ (go) F vendor/plugfix/module-b/src/Bee.php 644 95a279ce42abc9be
vendor/composer/autoload_classmap.php:
+ (go)     'Plugfix\\ModuleB\\Bee' => $vendorDir . '/plugfix/module-b/src/Bee.php',
- (oracle)     'Plugfix\\ModuleB\\Bee' => $baseDir . '/app/mods/bee/src/Bee.php',
vendor/composer/autoload_psr4.php:
+ (go)     'Plugfix\\ModuleA\\' => array($vendorDir . '/plugfix/module-a/src'),
- (oracle)     'Plugfix\\ModuleA\\' => array($baseDir . '/app/mods/module-a/src'),
vendor/composer/autoload_static.php:
+ (go)             0 => __DIR__ . '/..' . '/plugfix/module-a/src',
- (oracle)             0 => __DIR__ . '/../..' . '/app/mods/module-a/src',
+ (go)         'Plugfix\\ModuleB\\Bee' => __DIR__ . '/..' . '/plugfix/module-b/src/Bee.php',
- (oracle)         'Plugfix\\ModuleB\\Bee' => __DIR__ . '/../..' . '/app/mods/bee/src/Bee.php',
vendor/composer/installed.json:
+ (go)             "install-path": "../plugfix/module-a"
- (oracle)             "install-path": "../../app/mods/module-a"
+ (go)             "install-path": "../plugfix/module-b"
- (oracle)             "install-path": "../../app/mods/bee"
vendor/composer/installed.php:
+ (go)             'install_path' => __DIR__ . '/../plugfix/module-a',
- (oracle)             'install_path' => __DIR__ . '/../../app/mods/module-a',
+ (go)             'install_path' => __DIR__ . '/../plugfix/module-b',
- (oracle)             'install_path' => __DIR__ . '/../../app/mods/bee',
```

## plugin-installer-install

```
exit code: oracle=0 go=1
stdout:
- (oracle) [plugfix/installer] installed plugfix/module-a into modules/module-a (present)
- (oracle) [plugfix/installer] installed plugfix/module-b into modules/bee (present)
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/installer-plugin contains a Composer plugin which is allowed by you  
+ (go)   r allow-plugins config, but gocomposer cannot execute Composer plugins. Run  
+ (go)    with --no-plugins to install without activating plugins.                    
+ (go)                                                                                
+ (go) 
+ (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
+ (go) 
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Verifying lock file contents can be installed on current platform.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Package operations: 4 installs, 0 updates, 0 removals
- (oracle)   - Downloading plugfix/installer-plugin (1.0.0)
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle)   - Downloading plugfix/module-a (1.0.0)
- (oracle)   - Downloading plugfix/module-b (1.0.0)
- (oracle)   - Installing plugfix/installer-plugin (1.0.0): Extracting archive
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
- (oracle)   - Installing plugfix/module-a (1.0.0): Extracting archive
- (oracle)   - Installing plugfix/module-b (1.0.0): Extracting archive
- (oracle) Generating autoload files
project tree:
- (oracle) D modules
- (oracle) D modules/bee
- (oracle) D modules/bee/src
- (oracle) D modules/module-a
- (oracle) D modules/module-a/bin
- (oracle) D modules/module-a/src
- (oracle) D vendor
- (oracle) D vendor/acme
- (oracle) D vendor/acme/log
- (oracle) D vendor/acme/log/src
- (oracle) D vendor/bin
- (oracle) D vendor/composer
- (oracle) D vendor/plugfix
- (oracle) D vendor/plugfix/installer-plugin
- (oracle) D vendor/plugfix/installer-plugin/src
- (oracle) F modules/bee/composer.json 644 4d2fc4f52472a068
- (oracle) F modules/bee/src/Bee.php 644 95a279ce42abc9be
- (oracle) F modules/module-a/bin/module-a 755 497f6811dd69c6c6
- (oracle) F modules/module-a/composer.json 644 42213e55ec719364
- (oracle) F modules/module-a/src/Module.php 644 5fb3895bcb248f97
- (oracle) F vendor/acme/log/composer.json 644 b83d42805152f028
- (oracle) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
- (oracle) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
- (oracle) F vendor/autoload.php 644 1f62180e98d7d489
- (oracle) F vendor/bin/module-a 755 3adc06b2b2fd7aad
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 c11bdea06020b395
- (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
- (oracle) F vendor/composer/autoload_psr4.php 644 fd5a0a19865aceb2
- (oracle) F vendor/composer/autoload_real.php 644 06156c11d4e6c368
- (oracle) F vendor/composer/autoload_static.php 644 67977963a6de732c
- (oracle) F vendor/composer/installed.json 644 e3809b62ec53d12c
- (oracle) F vendor/composer/installed.php 644 d1384fc76d15b4a8
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
- (oracle) F vendor/plugfix/installer-plugin/composer.json 644 8ffc5973dfda0eeb
- (oracle) F vendor/plugfix/installer-plugin/src/ModuleInstaller.php 644 9b1097da74043115
- (oracle) F vendor/plugfix/installer-plugin/src/Plugin.php 644 6eac7b557ee82279
autoload probe:
- (oracle) class Acme\Log\Logger: autoloaded
- (oracle) class Acme\Log\LoggerInterface: autoloaded
- (oracle) class App\Kernel: autoloaded
- (oracle) Fatal error: Uncaught Error: Class "Composer\Installer\LibraryInstaller" not found in /work/project/vendor/plugfix/installer-plugin/src/ModuleInstaller.php:9
- (oracle) Stack trace:
- (oracle) #0 /work/project/vendor/composer/ClassLoader.php(576): include()
- (oracle) #1 /work/project/vendor/composer/ClassLoader.php(427): {closure:Composer\Autoload\ClassLoader::initializeIncludeClosure():575}('/work/project/v...')
- (oracle) #2 [internal function]: Composer\Autoload\ClassLoader->loadClass('Plugfix\\Install...')
- (oracle) #3 /gobin/probe.php(42): class_exists('Plugfix\\Install...')
- (oracle) #4 {main}
- (oracle)   thrown in /work/project/vendor/plugfix/installer-plugin/src/ModuleInstaller.php on line 9
- (oracle) 
vendor/composer/autoload_classmap.php:
- (oracle) <?php
- (oracle) // autoload_classmap.php @generated by Composer
- (oracle) 
- (oracle) $vendorDir = dirname(__DIR__);
- (oracle) $baseDir = dirname($vendorDir);
- (oracle) 
- (oracle) return array(
- (oracle)     'Composer\\InstalledVersions' => $vendorDir . '/composer/InstalledVersions.php',
- (oracle)     'Plugfix\\ModuleB\\Bee' => $baseDir . '/modules/bee/src/Bee.php',
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
- (oracle)     'Plugfix\\ModuleA\\' => array($baseDir . '/modules/module-a/src'),
- (oracle)     'Plugfix\\Installer\\' => array($vendorDir . '/plugfix/installer-plugin/src'),
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInitf5b73ea3441589fa0e80e3fff1fca18d
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInitf5b73ea3441589fa0e80e3fff1fca18d', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInitf5b73ea3441589fa0e80e3fff1fca18d', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInitf5b73ea3441589fa0e80e3fff1fca18d::getInitializer($loader));
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
- (oracle) class ComposerStaticInitf5b73ea3441589fa0e80e3fff1fca18d
- (oracle) {
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'P' =>
- (oracle)         array (
- (oracle)             'Plugfix\\ModuleA\\' => 16,
- (oracle)             'Plugfix\\Installer\\' => 18,
- (oracle)         ),
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)             'Acme\\Log\\' => 9,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
- (oracle)         'Plugfix\\ModuleA\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/modules/module-a/src',
- (oracle)         ),
- (oracle)         'Plugfix\\Installer\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/plugfix/installer-plugin/src',
- (oracle)         ),
- (oracle)         'App\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/../..' . '/src',
- (oracle)         ),
- (oracle)         'Acme\\Log\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $classMap = array (
- (oracle)         'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
- (oracle)         'Plugfix\\ModuleB\\Bee' => __DIR__ . '/../..' . '/modules/bee/src/Bee.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static function getInitializer(ClassLoader $loader)
- (oracle)     {
- (oracle)         return \Closure::bind(function () use ($loader) {
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInitf5b73ea3441589fa0e80e3fff1fca18d::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInitf5b73ea3441589fa0e80e3fff1fca18d::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInitf5b73ea3441589fa0e80e3fff1fca18d::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "acme/log",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-log-1.1.0.zip",
- (oracle)                 "reference": "58b0f69fc5f96981e5f7fae32d592232a36012c2",
- (oracle)                 "shasum": "c40f1901ba3e909555618861c594e339cd678f39"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "php": ">=7.4"
- (oracle)             },
- (oracle)             "suggest": {
- (oracle)                 "acme/extra": "for extra stuff"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\Log\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Logging",
- (oracle)             "install-path": "../acme/log"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/installer-plugin",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-installer-plugin-1.0.0.zip",
- (oracle)                 "reference": "bddb4bec7433b4d507cb1fa318c64a9089775057",
- (oracle)                 "shasum": "a006b525b6afa5d20997b156f1b1b84122fe22e7"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "composer-plugin-api": "^2.1"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "composer-plugin",
- (oracle)             "extra": {
- (oracle)                 "class": "Plugfix\\Installer\\Plugin",
- (oracle)                 "plugin-modifies-install-path": true
- (oracle)             },
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Plugfix\\Installer\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture plugin providing a custom installer",
- (oracle)             "install-path": "../plugfix/installer-plugin"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/module-a",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.0.0.zip",
- (oracle)                 "reference": "25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9",
- (oracle)                 "shasum": "fb1a008a58736c94e8d9f38dddedf3b40d4175ca"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "plugfix/installer-plugin": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "bin": [
- (oracle)                 "bin/module-a"
- (oracle)             ],
- (oracle)             "type": "plugfix-module",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Plugfix\\ModuleA\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture module installed by the custom installer",
- (oracle)             "install-path": "../../modules/module-a"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "plugfix/module-b",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-b-1.0.0.zip",
- (oracle)                 "reference": "6ff5d5ac34200d3b06845803fd86da504556b97b",
- (oracle)                 "shasum": "e51c205982294cf6bc7998824bcc3fb4ba492072"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "plugfix/installer-plugin": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "plugfix-module",
- (oracle)             "extra": {
- (oracle)                 "module-name": "bee"
- (oracle)             },
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "classmap": [
- (oracle)                     "src/"
- (oracle)                 ]
- (oracle)             },
- (oracle)             "license": [
- (oracle)                 "MIT"
- (oracle)             ],
- (oracle)             "description": "Fixture module with a custom module name",
- (oracle)             "install-path": "../../modules/bee"
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
- (oracle)         'acme/log' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/log',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'plugfix/installer-plugin' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'bddb4bec7433b4d507cb1fa318c64a9089775057',
- (oracle)             'type' => 'composer-plugin',
- (oracle)             'install_path' => __DIR__ . '/../plugfix/installer-plugin',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'plugfix/module-a' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9',
- (oracle)             'type' => 'plugfix-module',
- (oracle)             'install_path' => __DIR__ . '/../../modules/module-a',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'plugfix/module-b' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '6ff5d5ac34200d3b06845803fd86da504556b97b',
- (oracle)             'type' => 'plugfix-module',
- (oracle)             'install_path' => __DIR__ . '/../../modules/bee',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'test/project' => array(
- (oracle)             'pretty_version' => '1.0.0+no-version-set',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => null,
- (oracle)             'type' => 'project',
- (oracle)             'install_path' => __DIR__ . '/../../',
- (oracle)             'aliases' => array(),
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

## plugin-installer-upgrade

```
exit code: oracle=0 go=1
stdout:
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/installer-plugin contains a Composer plugin which is allowed by you  
+ (go)   r allow-plugins config, but gocomposer cannot execute Composer plugins. Run  
+ (go)    with --no-plugins to install without activating plugins.                    
+ (go)                                                                                
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 0 installs, 1 update, 0 removals
- (oracle)   - Upgrading plugfix/module-a (1.0.0 => 1.1.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 0 installs, 1 update, 0 removals
- (oracle)   - Downloading plugfix/module-a (1.1.0)
- (oracle)   - Upgrading plugfix/module-a (1.0.0 => 1.1.0): Extracting archive
- (oracle) Generating autoload files
composer.lock:
+ (go)     "content-hash": "f5b73ea3441589fa0e80e3fff1fca18d",
- (oracle)     "content-hash": "fcac42c7dcfe7ba099f944250feed70b",
+ (go)             "version": "1.0.0",
- (oracle)             "version": "1.1.0",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.0.0.zip",
+ (go)                 "reference": "25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9",
+ (go)                 "shasum": "fb1a008a58736c94e8d9f38dddedf3b40d4175ca"
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.1.0.zip",
- (oracle)                 "reference": "dae582f92819352a86a522d1299b525aa61d1779",
- (oracle)                 "shasum": "bd779d0baa928bdd60d061bc03339ed2d3df5fbd"
project tree:
+ (go) F composer.lock 644 1bc00b4e27a0c3e7
- (oracle) F composer.lock 644 1864c34a63d3f042
+ (go) F modules/module-a/bin/module-a 755 497f6811dd69c6c6
- (oracle) F modules/module-a/bin/module-a 755 970172fd55d554fc
+ (go) F modules/module-a/src/Module.php 644 5fb3895bcb248f97
- (oracle) F modules/module-a/src/Module.php 644 fc168f7d81ae87c4
+ (go) F vendor/composer/installed.json 644 e3809b62ec53d12c
+ (go) F vendor/composer/installed.php 644 d1384fc76d15b4a8
- (oracle) F vendor/composer/installed.json 644 45346ab9d6f8f2a4
- (oracle) F vendor/composer/installed.php 644 436da8ce610782b4
vendor/composer/installed.json:
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
- (oracle)             "version": "1.1.0",
- (oracle)             "version_normalized": "1.1.0.0",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.0.0.zip",
+ (go)                 "reference": "25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9",
+ (go)                 "shasum": "fb1a008a58736c94e8d9f38dddedf3b40d4175ca"
- (oracle)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.1.0.zip",
- (oracle)                 "reference": "dae582f92819352a86a522d1299b525aa61d1779",
- (oracle)                 "shasum": "bd779d0baa928bdd60d061bc03339ed2d3df5fbd"
vendor/composer/installed.php:
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9',
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => 'dae582f92819352a86a522d1299b525aa61d1779',
```

## plugin-installer-remove

```
exit code: oracle=0 go=1
stdout:
- (oracle) Deleting modules/module-a - deleted
- (oracle) No security vulnerability advisories found.
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   plugfix/installer-plugin contains a Composer plugin which is allowed by you  
+ (go)   r allow-plugins config, but gocomposer cannot execute Composer plugins. Run  
+ (go)    with --no-plugins to install without activating plugins.                    
+ (go)                                                                                
+ (go) 
+ (go) remove [--dev] [--dry-run] [--no-progress] [--no-update] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--update-no-dev] [-w|--update-with-dependencies] [-W|--update-with-all-dependencies] [--with-all-dependencies] [--no-update-with-dependencies] [-m|--minimal-changes] [--unused] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--] [<packages>...]
+ (go) 
- (oracle) Running composer update plugfix/module-a
- (oracle) Loading composer repositories with package information
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Updating dependencies
- (oracle) Lock file operations: 0 installs, 0 updates, 1 removal
- (oracle)   - Removing plugfix/module-a (1.0.0)
- (oracle) Writing lock file
- (oracle) Installing dependencies from lock file (including require-dev)
- (oracle) Package operations: 0 installs, 0 updates, 1 removal
- (oracle)   - Removing plugfix/module-a (1.0.0)
- (oracle) Generating autoload files
composer.lock:
+ (go)     "content-hash": "f5b73ea3441589fa0e80e3fff1fca18d",
- (oracle)     "content-hash": "bf4774f249da4ac0db76ac25ebbc687b",
+ (go)             "name": "plugfix/module-a",
+ (go)             "version": "1.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.0.0.zip",
+ (go)                 "reference": "25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9",
+ (go)                 "shasum": "fb1a008a58736c94e8d9f38dddedf3b40d4175ca"
+ (go)             },
+ (go)             "require": {
+ (go)                 "plugfix/installer-plugin": "^1.0"
+ (go)             },
+ (go)             "bin": [
+ (go)                 "bin/module-a"
+ (go)             ],
+ (go)             "type": "plugfix-module",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\ModuleA\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture module installed by the custom installer",
+ (go)             "time": "2024-01-02T03:04:05+00:00"
+ (go)         },
+ (go)         {
project tree:
+ (go) D modules/module-a
+ (go) D modules/module-a/bin
+ (go) D modules/module-a/src
+ (go) D vendor/bin
+ (go) F composer.lock 644 1bc00b4e27a0c3e7
- (oracle) F composer.lock 644 319130f2471e7cb9
+ (go) F modules/module-a/bin/module-a 755 497f6811dd69c6c6
+ (go) F modules/module-a/composer.json 644 42213e55ec719364
+ (go) F modules/module-a/src/Module.php 644 5fb3895bcb248f97
+ (go) F vendor/bin/module-a 755 3adc06b2b2fd7aad
+ (go) F vendor/composer/autoload_psr4.php 644 fd5a0a19865aceb2
- (oracle) F vendor/composer/autoload_psr4.php 644 d85a6e4f8252a552
+ (go) F vendor/composer/autoload_static.php 644 67977963a6de732c
+ (go) F vendor/composer/installed.json 644 e3809b62ec53d12c
+ (go) F vendor/composer/installed.php 644 d1384fc76d15b4a8
- (oracle) F vendor/composer/autoload_static.php 644 9326a2cd8b66f335
- (oracle) F vendor/composer/installed.json 644 3c3ecad23a4ca074
- (oracle) F vendor/composer/installed.php 644 48eb0a57bf64b3aa
vendor/composer/autoload_psr4.php:
+ (go)     'Plugfix\\ModuleA\\' => array($baseDir . '/modules/module-a/src'),
vendor/composer/autoload_static.php:
+ (go)             'Plugfix\\ModuleA\\' => 16,
+ (go)         'Plugfix\\ModuleA\\' =>
+ (go)         array (
+ (go)             0 => __DIR__ . '/../..' . '/modules/module-a/src',
+ (go)         ),
vendor/composer/installed.json:
+ (go)             "name": "plugfix/module-a",
+ (go)             "version": "1.0.0",
+ (go)             "version_normalized": "1.0.0.0",
+ (go)             "dist": {
+ (go)                 "type": "zip",
+ (go)                 "url": "http://127.0.0.1:8899/dists/plugfix-module-a-1.0.0.zip",
+ (go)                 "reference": "25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9",
+ (go)                 "shasum": "fb1a008a58736c94e8d9f38dddedf3b40d4175ca"
+ (go)             },
+ (go)             "require": {
+ (go)                 "plugfix/installer-plugin": "^1.0"
+ (go)             },
+ (go)             "time": "2024-01-02T03:04:05+00:00",
+ (go)             "bin": [
+ (go)                 "bin/module-a"
+ (go)             ],
+ (go)             "type": "plugfix-module",
+ (go)             "installation-source": "dist",
+ (go)             "autoload": {
+ (go)                 "psr-4": {
+ (go)                     "Plugfix\\ModuleA\\": "src/"
+ (go)                 }
+ (go)             },
+ (go)             "license": [
+ (go)                 "MIT"
+ (go)             ],
+ (go)             "description": "Fixture module installed by the custom installer",
+ (go)             "install-path": "../../modules/module-a"
+ (go)         },
+ (go)         {
vendor/composer/installed.php:
+ (go)         'plugfix/module-a' => array(
+ (go)             'pretty_version' => '1.0.0',
+ (go)             'version' => '1.0.0.0',
+ (go)             'reference' => '25a36bbe4a1b2d576e4cde6893a7afe5e9f584c9',
+ (go)             'type' => 'plugfix-module',
+ (go)             'install_path' => __DIR__ . '/../../modules/module-a',
+ (go)             'aliases' => array(),
+ (go)             'dev_requirement' => false,
+ (go)         ),
```

## plugin-installer-no-plugins

```
stderr:
- (oracle) The "plugfix/installer-plugin" plugin was not loaded as plugins are disabled.
```

## plugin-installer-dump

```
exit code: oracle=0 go=1
stdout:
+ (go) Not all dependencies are installed. Make sure to run a "composer install" to install missing dependencies
- (oracle) Generated optimized autoload files containing 8 classes
stderr:
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Could not scan for classes inside "/work/project/vendor/plugfix/module-b/sr  
+ (go)   c/" which does not appear to be a file nor a folder                          
+ (go)                                                                                
+ (go) 
+ (go) dump-autoload [-o|--optimize] [-a|--classmap-authoritative] [--apcu] [--apcu-prefix APCU-PREFIX] [--dry-run] [--dev] [--no-dev] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--strict-psr] [--strict-ambiguous]
+ (go) 
+ (go) 
project tree:
+ (go) F vendor/composer/autoload_classmap.php 644 c11bdea06020b395
- (oracle) F vendor/composer/autoload_classmap.php 644 b363023dd50ae7de
+ (go) F vendor/composer/autoload_static.php 644 67977963a6de732c
- (oracle) F vendor/composer/autoload_static.php 644 d564d4a4bd434f3b
vendor/composer/autoload_classmap.php:
- (oracle)     'Acme\\Log\\Logger' => $vendorDir . '/acme/log/src/Logger.php',
- (oracle)     'Acme\\Log\\LoggerInterface' => $vendorDir . '/acme/log/src/LoggerInterface.php',
- (oracle)     'App\\Kernel' => $baseDir . '/src/Kernel.php',
- (oracle)     'Plugfix\\Installer\\ModuleInstaller' => $vendorDir . '/plugfix/installer-plugin/src/ModuleInstaller.php',
- (oracle)     'Plugfix\\Installer\\Plugin' => $vendorDir . '/plugfix/installer-plugin/src/Plugin.php',
- (oracle)     'Plugfix\\ModuleA\\Module' => $baseDir . '/modules/module-a/src/Module.php',
vendor/composer/autoload_static.php:
- (oracle)         'Acme\\Log\\Logger' => __DIR__ . '/..' . '/acme/log/src/Logger.php',
- (oracle)         'Acme\\Log\\LoggerInterface' => __DIR__ . '/..' . '/acme/log/src/LoggerInterface.php',
- (oracle)         'App\\Kernel' => __DIR__ . '/../..' . '/src/Kernel.php',
- (oracle)         'Plugfix\\Installer\\ModuleInstaller' => __DIR__ . '/..' . '/plugfix/installer-plugin/src/ModuleInstaller.php',
- (oracle)         'Plugfix\\Installer\\Plugin' => __DIR__ . '/..' . '/plugfix/installer-plugin/src/Plugin.php',
- (oracle)         'Plugfix\\ModuleA\\Module' => __DIR__ . '/../..' . '/modules/module-a/src/Module.php',
```

## vcs-forge-github-enterprise

```
stderr:
- (oracle) Running cache garbage collection
```
