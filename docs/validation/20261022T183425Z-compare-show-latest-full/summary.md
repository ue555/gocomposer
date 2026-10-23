# Composer vs gocomposer comparison

- gocomposer_sha256: `f17ab67478d944bf5e0d6cc8a6a918e8ad322b7ba2fcb32bc1d3c269990868be`
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
| clear-cache-gc-ttl | DIFF | 0 | 0 |
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
| vcs-debug | DIFF | 0 | 0 |
| vcs-only-branches | MATCH | 0 | 0 |
| vcs-only-branches-invalid | MATCH | 1 | 1 |
| vcs-misc-vv | MATCH | 0 | 0 |
| vcs-misc-numeric-branch | MATCH | 0 | 0 |
| vcs-misc-slash-branch | MATCH | 0 | 0 |
| vcs-misc-old-tag | MATCH | 0 | 0 |
| vcs-ref-pin | MATCH | 0 | 0 |
| vcs-not-a-repo | MATCH | 1 | 1 |
| vcs-git-type-not-a-repo | MATCH | 1 | 1 |
| vcs-type-hg | DIFF | 1 | 1 |
| vcs-type-svn | DIFF | 1 | 1 |
| vcs-type-fossil | DIFF | 1 | 1 |
| vcs-type-perforce | DIFF | 1 | 1 |
| vcs-type-pear | DIFF | 1 | 1 |
| vcs-type-forgejo | DIFF | 1 | 1 |
| vcs-type-github-local | DIFF | 1 | 1 |
| vcs-insecure-http | MATCH | 100 | 100 |
| vcs-http-missing | MATCH | 1 | 1 |
| vcs-http-cached | DIFF | 0 | 0 |
| vcs-http-cached-debug | DIFF | 0 | 0 |
| vcs-http-dev-update | MATCH | 0 | 0 |
| source-update-checkout | MATCH | 0 | 0 |
| source-update-checkout-v | MATCH | 0 | 0 |
| source-downgrade-v | MATCH | 0 | 0 |
| source-dev-update | MATCH | 0 | 0 |
| source-local-changes | DIFF | 1 | 1 |
| source-discard-true | MATCH | 0 | 0 |
| source-discard-stash | MATCH | 1 | 1 |
| source-discard-stash-ident | MATCH | 1 | 1 |
| source-uninstall-changes | DIFF | 1 | 1 |
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

## clear-cache-gc-ttl

```
project tree:
- (oracle) D c/vcs/foo.git
- (oracle) F c/files/acme/log/abc.zip 644 eb076a2ec6ced9ee
- (oracle) F c/repo/http---x/packages.json 644 44136fa355b3678a
- (oracle) F c/vcs/foo.git/HEAD 644 3ff6c05723bb069d
```

## vcs-debug

```
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
+ (go) Executing command (CWD): 'git' '--version'
- (oracle) Writing /work/cache/repo/-fixtures-git-gitlib.git/b2b1b720ad1e4e06cb3d641ea573be406bab1502 into cache
- (oracle) Writing /work/cache/repo/-fixtures-git-gitlib.git/702184333938c0ff49e87f0af7d9931d71486a52 into cache
- (oracle) Writing /work/cache/repo/-fixtures-git-gitlib.git/ba0333212d6fb4b08718b30b9dc442e929393f8b into cache
- (oracle) Writing /work/cache/repo/-fixtures-git-gitlib.git/24bec3fa0d593ddedbd07cb32c28449d25b4f5a0 into cache
- (oracle) Writing /work/cache/repo/-fixtures-git-gitdep.git/c16e0182384ca46e31fb08304ab23bc789f9e70a into cache
- (oracle) Writing /work/cache/repo/-fixtures-git-gitdep.git/95b74ce18d1c01a90035538fec87a2873119a4dd into cache
- (oracle) Built pool in 0.21747279167175s
- (oracle) Downloading http://127.0.0.1:8899/packages.json
- (oracle) [200] http://127.0.0.1:8899/packages.json
- (oracle) Writing /work/cache/repo/http---127.0.0.1-8899/packages.json into cache
- (oracle) Running pool optimizer.
- (oracle) Looking at all rules.
- (oracle) Reading ./composer.lock (/work/project/composer.lock)
```

## vcs-type-hg

```
stderr:
+ (go)                                                                                
+ (go)   The hg VCS driver is not supported by gocomposer (repository /fixtures/git/  
+ (go)   gitlib.git); only git, github, gitlab and bitbucket repositories can be rea  
+ (go)   d                                                                            
+ (go)                                                                                
- (oracle)                          
- (oracle)   Undefined array key 0  
- (oracle)                          
```

## vcs-type-svn

```
stderr:
+ (go)   The svn VCS driver is not supported by gocomposer (repository /fixtures/git  
+ (go)   /gitlib.git); only git, github, gitlab and bitbucket repositories can be re  
+ (go)   ad                                                                           
- (oracle)   Failed to load file:///fixtures/git/gitlib.git, svn was not found, check th  
- (oracle)   at it is installed and in your PATH env.                                     
- (oracle)   sh: 1: exec: svn: not found                                                  
- (oracle)                                                                                
```

## vcs-type-fossil

```
stderr:
+ (go)                                                                                
+ (go)   The fossil VCS driver is not supported by gocomposer (repository /fixtures/  
+ (go)   git/gitlib.git); only git, github, gitlab and bitbucket repositories can be  
+ (go)    read                                                                        
+ (go)                                                                                
- (oracle)                                                                           
- (oracle)   fossil was not found, check that it is installed and in your PATH env.  
- (oracle)                                                                           
- (oracle)   sh: 1: exec: fossil: not found                                          
- (oracle)                                                                           
```

## vcs-type-perforce

```
stderr:
+ (go)                                                                                
+ (go)   The perforce VCS driver is not supported by gocomposer (repository /fixture  
+ (go)   s/git/gitlib.git); only git, github, gitlab and bitbucket repositories can   
+ (go)   be read                                                                      
+ (go)                                                                                
- (oracle)                                
- (oracle)   Undefined array key "depot"  
- (oracle)                                
```

## vcs-type-pear

```
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                           
+ (go)   The PEAR repository has been removed from Composer 2.x  
+ (go)                                                           
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) The PEAR repository has been removed from Composer 2.x
```

## vcs-type-forgejo

```
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                               
+ (go)   Repository type is not registered: forgejo  
+ (go)                                               
+ (go) 
+ (go) update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]
+ (go) 
- (oracle) Repository type is not registered: forgejo
```

## vcs-type-github-local

```
stderr:
+ (go)                                                                                
+ (go)   The github VCS driver is not supported by gocomposer (repository /fixtures/  
+ (go)   git/gitlib.git); only git, github, gitlab and bitbucket repositories can be  
+ (go)    read                                                                        
+ (go)                                                                                
- (oracle)                                                                   
- (oracle)   The GitHub repository URL /fixtures/git/gitlib.git is invalid.  
- (oracle)                                                                   
```

## vcs-http-cached

```
stderr:
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## vcs-http-cached-debug

```
stderr:
+ (go) Loading composer repositories with package information
- (oracle) Running 2.11.999-dev+source (@release_date@) with PHP 8.4.27 on Linux / 7.1.10-200.fc44.aarch64
- (oracle) Reading ./composer.json (/work/project/composer.json)
- (oracle) Loading config file ./composer.json (/work/project/composer.json)
- (oracle) Checked CA file /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a directory.
- (oracle) Checked CA file /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a directory.
- (oracle) Checked CA file /etc/ssl/certs/ca-certificates.crt: valid
- (oracle) Executing command (/work/project): 'git' 'branch' '-a' '--no-color' '--no-abbrev' '-v'
- (oracle) Executing command (/work/project): 'git' 'describe' '--exact-match' '--tags'
- (oracle) Executing command (/work/project): 'git' 'rev-list' '--no-commit-header' '--format=%H' '-n1' 'HEAD' '--no-show-signature'
- (oracle) Executing command (/work/project): 'hg' 'branch'
- (oracle) Executing command (/work/project): 'fossil' 'branch' 'list'
- (oracle) Executing command (/work/project): 'fossil' 'tag' 'list'
- (oracle) Executing command (/work/project): 'svn' 'info' '--xml'
- (oracle) Failed to initialize global composer: Composer could not find the config file: /work/home/composer.json
- (oracle) 
- (oracle) Reading ./composer.lock (/work/project/composer.lock)
- (oracle) Loading composer repositories with package information
- (oracle) Reading /work/out/cache/repo/http---127.0.0.1-8899-git-gitlib.git/b2b1b720ad1e4e06cb3d641ea573be406bab1502 from cache
- (oracle) Reading /work/out/cache/repo/http---127.0.0.1-8899-git-gitlib.git/702184333938c0ff49e87f0af7d9931d71486a52 from cache
- (oracle) Reading /work/out/cache/repo/http---127.0.0.1-8899-git-gitlib.git/ba0333212d6fb4b08718b30b9dc442e929393f8b from cache
- (oracle) Reading /work/out/cache/repo/http---127.0.0.1-8899-git-gitlib.git/24bec3fa0d593ddedbd07cb32c28449d25b4f5a0 from cache
- (oracle) Built pool in 0.10237383842468s
- (oracle) Reading /work/out/cache/repo/http---127.0.0.1-8899/packages.json from cache
- (oracle) Running pool optimizer.
- (oracle) Looking at all rules.
- (oracle) Reading ./composer.lock (/work/project/composer.lock)
```

## source-local-changes

```
stderr:
+ (go)     Update of vcs/gitlib failed
```

## source-uninstall-changes

```
stderr:
+ (go)     Uninstall of vcs/gitlib failed
```
