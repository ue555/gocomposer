# Composer vs gocomposer comparison

- gocomposer_sha256: `ae13d939edee6c42f53ca23e85a95f5fd539c53424fe5f3eb8251dc4e3dfe35c`
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
| completion-bash | MATCH | 0 | 0 |
| completion-shell-env | MATCH | 0 | 0 |
| completion-no-shell | MATCH | 2 | 2 |
| completion-zsh | MATCH | 2 | 2 |
| completion-fish-env | MATCH | 2 | 2 |
| completion-quiet | MATCH | 0 | 0 |
| complete-commands | MATCH | 0 | 0 |
| complete-commands-prefix | MATCH | 0 | 0 |
| complete-commands-abbrev | MATCH | 0 | 0 |
| complete-commands-alias | MATCH | 0 | 0 |
| complete-commands-exact | MATCH | 0 | 0 |
| complete-app-options | MATCH | 0 | 0 |
| complete-command-options | MATCH | 0 | 0 |
| complete-command-options-short | MATCH | 0 | 0 |
| complete-working-dir-first | MATCH | 0 | 0 |
| complete-verbose-first | MATCH | 0 | 0 |
| complete-after-separator | MATCH | 0 | 0 |
| complete-none | MATCH | 0 | 0 |
| complete-format | MATCH | 0 | 0 |
| complete-format-eq | MATCH | 0 | 0 |
| complete-format-short | MATCH | 0 | 0 |
| complete-prefer-install | MATCH | 0 | 0 |
| complete-audit-format | MATCH | 0 | 0 |
| complete-bump-after-update | MATCH | 0 | 0 |
| complete-audit-options | MATCH | 0 | 0 |
| complete-audit-abandoned | MATCH | 0 | 0 |
| complete-licenses-format | MATCH | 0 | 0 |
| complete-fund-format | MATCH | 0 | 0 |
| complete-reinstall-type | MATCH | 0 | 0 |
| complete-show-ignore | MATCH | 0 | 0 |
| complete-init-require | MATCH | 0 | 0 |
| complete-help-format | MATCH | 0 | 0 |
| complete-list-format | MATCH | 0 | 0 |
| complete-global-option-value | MATCH | 0 | 0 |
| complete-update-packages | MATCH | 0 | 0 |
| complete-update-packages-prefix | MATCH | 0 | 0 |
| complete-remove-packages | MATCH | 0 | 0 |
| complete-bump-packages | MATCH | 0 | 0 |
| complete-depends-package | MATCH | 0 | 0 |
| complete-depends-platform | MATCH | 0 | 0 |
| complete-prohibits-package | MATCH | 0 | 0 |
| complete-require-vendors | MATCH | 0 | 0 |
| complete-require-vendor | MATCH | 0 | 0 |
| complete-require-names | MATCH | 0 | 0 |
| complete-require-platform | MATCH | 0 | 0 |
| complete-require-php | MATCH | 0 | 0 |
| complete-create-project-package | MATCH | 0 | 0 |
| complete-archive-package | MATCH | 0 | 0 |
| complete-show-installed | MATCH | 0 | 0 |
| complete-show-platform | MATCH | 0 | 0 |
| complete-show-available | MATCH | 0 | 0 |
| complete-outdated-package | MATCH | 0 | 0 |
| complete-suggests-packages | MATCH | 0 | 0 |
| complete-browse-packages | MATCH | 0 | 0 |
| complete-reinstall-packages | MATCH | 0 | 0 |
| complete-run-script | MATCH | 0 | 0 |
| complete-run-script-alias | MATCH | 0 | 0 |
| complete-exec-binary | MATCH | 0 | 0 |
| complete-help-command | MATCH | 0 | 0 |
| complete-list-namespace | MATCH | 0 | 0 |
| complete-completion-shell | MATCH | 0 | 0 |
| complete-global-command | MATCH | 0 | 0 |
| complete-global-options | MATCH | 0 | 0 |
| complete-repository-action | MATCH | 0 | 0 |
| complete-repository-enable | MATCH | 0 | 0 |
| complete-repository-add-type | MATCH | 0 | 0 |
| complete-policy-action | MATCH | 0 | 0 |
| complete-policy-arg1 | MATCH | 0 | 0 |
| complete-script-command | MATCH | 0 | 0 |
| complete-config-keys | MATCH | 0 | 0 |
| complete-config-keys-prefix | MATCH | 0 | 0 |
| complete-config-keys-extra | MATCH | 0 | 0 |
| complete-config-keys-unset | MATCH | 0 | 0 |
| complete-config-keys-global | MATCH | 0 | 0 |
| complete-config-keys-list | MATCH | 0 | 0 |
| complete-no-shell | MATCH | 0 | 0 |
| complete-bad-shell | MATCH | 2 | 2 |
| complete-no-current | MATCH | 2 | 2 |
| complete-bad-current | MATCH | 2 | 2 |
| complete-index-out-of-range | MATCH | 2 | 2 |
| complete-missing-input-value | MATCH | 1 | 1 |
| complete-repository-names | MATCH | 0 | 0 |
| complete-repository-before | MATCH | 0 | 0 |
| complete-policy-names | MATCH | 0 | 0 |
| complete-unlocked-installed | MATCH | 0 | 0 |
| complete-no-json-command | MATCH | 0 | 0 |
| complete-no-json-update | MATCH | 1 | 1 |
| complete-no-json-config | MATCH | 0 | 0 |
| complete-global-require | MATCH | 0 | 0 |
| complete-global-show | MATCH | 0 | 0 |
| complete-global-no-home-json | MATCH | 1 | 1 |
| complete-plugin-commands | MATCH | 0 | 0 |
| complete-plugin-namespace | MATCH | 0 | 0 |
| complete-plugin-options | MATCH | 0 | 0 |
| complete-plugin-no-plugins | MATCH | 0 | 0 |
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
| plugin2-dl-update | MATCH | 0 | 0 |
| plugin2-dl-update-v2 | MATCH | 0 | 0 |
| plugin2-dl-rewrite | MATCH | 0 | 0 |
| plugin2-dl-meta | MATCH | 0 | 0 |
| plugin2-dl-pool | MATCH | 0 | 0 |
| plugin2-dl-pool-partial | MATCH | 2 | 2 |
| plugin2-dl-install-lock | MATCH | 0 | 0 |
| plugin2-dl-throw-pre | MATCH | 1 | 1 |
| plugin2-dl-throw-pool | MATCH | 1 | 1 |
| plugin2-dl-cache-hit | MATCH | 0 | 0 |
| plugin2-dl-require | MATCH | 0 | 0 |
| plugin2-dl-no-plugins | MATCH | 0 | 0 |
| plugin2-api-http | MATCH | 0 | 0 |
| plugin2-api-repos | MATCH | 0 | 0 |
| plugin2-api-create-repo | MATCH | 0 | 0 |
| plugin2-api-factory | MATCH | 0 | 0 |
| plugin2-api-merge | MATCH | 0 | 0 |
| plugin2-api-errors | MATCH | 0 | 0 |
| plugin2-api-autoload | DIFF | 0 | 1 |
| plugin2-api-dm | MATCH | 0 | 0 |
| plugin2-api-install | MATCH | 0 | 0 |
| plugin2-api-update | MATCH | 0 | 0 |
| plugin2-api-addrepo | MATCH | 0 | 0 |
| plugin2-api-merge-activate | MATCH | 0 | 0 |
| plugin2-api-downloader | MATCH | 0 | 0 |
| plugin2-script-listeners | MATCH | 0 | 0 |
| plugin2-dl-debug-events | MATCH | 0 | 0 |
| plugin2-installers-fresh | MATCH | 0 | 0 |
| plugin2-installers-installed | MATCH | 0 | 0 |
| plugin2-installers-upgrade | MATCH | 0 | 0 |
| plugin2-installers-remove | MATCH | 0 | 0 |
| plugin2-installers-install-lock | MATCH | 0 | 0 |
| plugin2-installers-fail-download | MATCH | 1 | 1 |
| plugin2-installers-fail-install | MATCH | 1 | 1 |
| flex-update-fresh | MATCH | 0 | 0 |
| flex-update-installed | MATCH | 0 | 0 |
| flex-require-alias | MATCH | 0 | 0 |
| flex-require-contrib | MATCH | 0 | 0 |
| flex-require-recipe-lib | MATCH | 0 | 0 |
| flex-require-auto-bundle | MATCH | 0 | 0 |
| flex-require-pack | MATCH | 0 | 0 |
| flex-require-split | MATCH | 0 | 0 |
| flex-update-split | MATCH | 0 | 0 |
| flex-remove | MATCH | 0 | 0 |
| flex-install-lock | MATCH | 0 | 0 |
| flex-recipes | MATCH | 0 | 0 |
| flex-recipes-one | MATCH | 0 | 0 |
| flex-recipes-install | MATCH | 0 | 0 |
| flex-sync-recipes-alias | MATCH | 0 | 0 |
| flex-dump-env | MATCH | 1 | 1 |
| flex-update-noop | MATCH | 0 | 0 |
| flex-no-plugins | MATCH | 0 | 0 |
| flex-create-project | MATCH | 0 | 0 |
| flex-recipes-update | MATCH | 100 | 100 |
| flex-recipes-update-no-changelog | MATCH | 0 | 0 |
| flex-not-in-root | MATCH | 0 | 0 |
| plugin3-dl-replace-install | MATCH | 0 | 0 |
| plugin3-dl-replace-update | MATCH | 0 | 0 |
| plugin3-dl-replace-remove | MATCH | 2 | 2 |
| plugin3-dl-replace-reinstall | MATCH | 0 | 0 |
| plugin3-dl-replace-fail | MATCH | 1 | 1 |
| plugin3-dl-extract | MATCH | 0 | 0 |
| plugin3-dl-extract-update | MATCH | 0 | 0 |
| plugin3-dl-extract-fail | MATCH | 1 | 1 |
| plugin3-dl-wrap | MATCH | 0 | 0 |
| plugin3-dl-wrap-update | MATCH | 0 | 0 |
| plugin3-dl-zip-parallel | MATCH | 0 | 0 |
| plugin3-dl-internals | MATCH | 0 | 0 |
| plugin3-dl-no-plugins | MATCH | 0 | 0 |
| plugin3-dl-order-parallel | MATCH | 0 | 0 |
| plugin3-dl-order-serial | MATCH | 0 | 0 |
| plugin3-dl-order-throw | MATCH | 1 | 1 |
| providers-update | MATCH | 0 | 0 |
| providers-update-verbose | MATCH | 0 | 0 |
| providers-install-lock | MATCH | 0 | 0 |
| providers-update-cached | MATCH | 0 | 0 |
| providers-root-listing | MATCH | 0 | 0 |
| providers-dev-alias | MATCH | 0 | 0 |
| providers-virtual | MATCH | 2 | 2 |
| providers-not-found | MATCH | 2 | 2 |
| providers-bad-hash | MATCH | 1 | 1 |
| providers-bad-hash-other | MATCH | 0 | 0 |
| providers-require | MATCH | 0 | 0 |
| providers-require-missing | MATCH | 1 | 1 |
| providers-show-available | MATCH | 0 | 0 |
| providers-show-all | MATCH | 0 | 0 |
| providers-show-all-names | MATCH | 0 | 0 |
| providers-outdated | MATCH | 0 | 0 |
| providers-update-one | MATCH | 0 | 0 |
| providers-depends | MATCH | 0 | 0 |
| providers-search | MATCH | 0 | 0 |
| providers-search-name | MATCH | 0 | 0 |
| providers-search-vendor | MATCH | 0 | 0 |
| providers-search-list | MATCH | 0 | 0 |
| providers-search-list-other | MATCH | 0 | 0 |
| providers-lazy-update | MATCH | 0 | 0 |
| providers-lazy-not-found | MATCH | 2 | 2 |
| providers-lazy-show | MATCH | 0 | 0 |
| providers-lazy-search | MATCH | 0 | 0 |
| providers-lazy-cached | MATCH | 0 | 0 |
| archive3-xz | MATCH | 0 | 0 |
| archive3-xz-status | MATCH | 1 | 1 |
| archive3-xz-reinstall | MATCH | 0 | 0 |
| archive3-xz-remove | MATCH | 0 | 0 |
| archive3-rar-status | MATCH | 1 | 1 |
| archive3-xz-broken | MATCH | 1 | 1 |
| archive3-rar-no-unrar | MATCH | 1 | 1 |
| archive3-rar | MATCH | 0 | 0 |
| archive3-rar-broken | MATCH | 1 | 1 |
| archive3-zip-7z | MATCH | 0 | 0 |
| archive3-zip-7z-zipext | MATCH | 0 | 0 |
| archive3-zip-7z-crc | MATCH | 1 | 1 |
| archive3-zip-7z-crc-zipext | MATCH | 1 | 1 |
| archive3-zip-no-command | MATCH | 1 | 1 |
| archive3-zip-no-command-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-crc-zipext | MATCH | 1 | 1 |
| warnings-update | MATCH | 0 | 0 |
| warnings-update-ansi | MATCH | 0 | 0 |
| warnings-show | MATCH | 0 | 0 |
| warnings-search-api | MATCH | 0 | 0 |
| warnings-search-list | MATCH | 0 | 0 |
| warnings-show-names-list | MATCH | 0 | 0 |
| warnings-audit-api | MATCH | 1 | 1 |
| warnings-update-audit-api | MATCH | 0 | 0 |
| warnings-filter-audit | MATCH | 0 | 0 |
| warnings-providers | MATCH | 0 | 0 |
| warnings-two-repos | MATCH | 0 | 0 |
| warnings-http-403 | MATCH | 100 | 100 |
| warnings-http-404-metadata | MATCH | 2 | 2 |
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
| term-href-show | MATCH | 0 | 0 |
| term-href-show-name | MATCH | 0 | 0 |
| term-href-show-latest | MATCH | 0 | 0 |
| term-href-licenses | MATCH | 0 | 0 |
| term-href-fund | MATCH | 0 | 0 |
| term-href-search | MATCH | 0 | 0 |
| term-href-depends | MATCH | 0 | 0 |
| term-href-config-list | MATCH | 0 | 0 |
| term-href-audit | DIFF | 1 | 1 |
| term-href-audit-plain | DIFF | 1 | 1 |
| term-href-jetbrains | MATCH | 0 | 0 |
| term-href-konsole-old | MATCH | 0 | 0 |
| term-href-konsole-new | MATCH | 0 | 0 |
| term-href-konsole-zero | MATCH | 0 | 0 |
| term-href-idea | MATCH | 0 | 0 |
| term-href-no-ansi | MATCH | 0 | 0 |
| term-href-ansi-pipe | MATCH | 0 | 0 |
| term-href-ansi-pipe-fund | MATCH | 0 | 0 |
| term-progress-install | MATCH | 0 | 0 |
| term-progress-update | DIFF | 0 | 0 |
| term-progress-require | DIFF | 0 | 0 |
| term-progress-verbose | MATCH | 0 | 0 |
| term-progress-very-verbose | DIFF | 0 | 0 |
| term-progress-narrow | DIFF | 0 | 0 |
| term-progress-tiny | DIFF | 0 | 2 |
| term-progress-no-progress | MATCH | 0 | 0 |
| term-progress-no-ansi | MATCH | 0 | 0 |
| term-progress-ci | MATCH | 0 | 0 |
| term-progress-ansi-pipe | MATCH | 0 | 0 |
| curl-refused | MATCH | 100 | 100 |
| curl-empty-reply | MATCH | 100 | 100 |
| curl-reset | MATCH | 100 | 100 |
| curl-reset-debug | DIFF | 100 | 100 |
| curl-http09 | MATCH | 100 | 100 |
| curl-partial | MATCH | 100 | 100 |
| curl-tls-eof | MATCH | 100 | 100 |
| curl-tls-wrong-version | MATCH | 100 | 100 |
| curl-tls-reset | MATCH | 100 | 100 |
| curl-tls-self-signed | MATCH | 100 | 100 |
| curl-tls-self-signed-system-ca | MATCH | 100 | 100 |
| curl-tls-wrong-host | MATCH | 100 | 100 |
| curl-tls-expired | MATCH | 100 | 100 |
| curl-tls-unknown-ca | MATCH | 100 | 100 |
| curl-tls-good | MATCH | 2 | 2 |
| curl-tls-show | MATCH | 100 | 100 |
| curl-timeout | MATCH | 100 | 100 |
| curl-tls-timeout | MATCH | 100 | 100 |
| curl-post-empty-reply | MATCH | 100 | 100 |
| curl-post-timeout-ignored | MATCH | 0 | 0 |
| root-validate | MATCH | 0 | 0 |
| root-validate-allowed | MATCH | 0 | 0 |
| root-validate-allowed-zero | MATCH | 0 | 0 |
| root-validate-no-plugins | MATCH | 0 | 0 |
| root-update | MATCH | 0 | 0 |
| root-list | DIFF | 0 | 0 |
| root-list-allowed | DIFF | 0 | 0 |
| root-script-alias | MATCH | 0 | 0 |
| root-tty-validate-yes | MATCH | 0 | 0 |
| root-tty-validate-default | MATCH | 0 | 0 |
| root-tty-validate-no | MATCH | 0 | 0 |
| root-tty-validate-n | MATCH | 0 | 0 |
| root-tty-validate-allowed | MATCH | 0 | 0 |
| root-tty-list-yes | DIFF | 0 | 0 |
| root-tty-list-no | DIFF | 0 | 0 |
| root-tty-script-alias-yes | MATCH | 0 | 0 |
| root-tty-update-yes | MATCH | 0 | 0 |
| root-tty-no-plugins | MATCH | 0 | 0 |
| tty-remove-alt-yes | MATCH | 0 | 0 |
| tty-remove-alt-default | MATCH | 0 | 0 |
| tty-remove-alt-no | MATCH | 0 | 0 |
| tty-remove-alt-invalid | MATCH | 0 | 0 |
| tty-remove-alt-eof | MATCH | 0 | 0 |
| tty-remove-alt-dev-wildcard | MATCH | 0 | 0 |
| tty-remove-alt-dry-run | MATCH | 0 | 0 |
| tty-remove-alt-no-interaction | MATCH | 0 | 0 |
| tty-remove-alt-no-ansi | MATCH | 0 | 0 |
| tty-init-full | MATCH | 0 | 0 |
| tty-init-defaults | MATCH | 0 | 0 |
| tty-init-options-as-defaults | MATCH | 0 | 0 |
| tty-init-require-options | MATCH | 0 | 0 |
| tty-init-abort | MATCH | 1 | 1 |
| tty-init-invalid-name | MATCH | 1 | 1 |
| tty-init-eof-first | MATCH | 0 | 0 |
| tty-init-eof-search | MATCH | 1 | 1 |
| tty-init-bad-license | MATCH | 1 | 1 |
| tty-init-invalid-author | MATCH | 0 | 0 |
| tty-init-pick-attempts | MATCH | 1 | 1 |
| tty-init-pick-empty | MATCH | 0 | 0 |
| tty-init-git | MATCH | 0 | 0 |
| tty-init-git-project-no | MATCH | 0 | 0 |
| tty-update-interactive-one | MATCH | 0 | 0 |
| tty-update-interactive-two | MATCH | 0 | 0 |
| tty-update-interactive-tab | MATCH | 0 | 0 |
| tty-update-interactive-arrows | MATCH | 0 | 0 |
| tty-update-interactive-backspace | MATCH | 1 | 1 |
| tty-update-interactive-invalid | MATCH | 1 | 1 |
| tty-update-interactive-empty | MATCH | 1 | 1 |
| tty-update-interactive-abort | MATCH | 1 | 1 |
| tty-update-interactive-filter | MATCH | 0 | 0 |
| tty-update-interactive-none | MATCH | 1 | 1 |
| tty-update-interactive-no-lock | MATCH | 3 | 3 |
| tty-update-interactive-n | MATCH | 1 | 1 |
| tty-require-dev-recommend-yes | MATCH | 0 | 0 |
| tty-require-dev-recommend-no | MATCH | 0 | 0 |
| tty-require-dev-recommend-two | MATCH | 0 | 0 |
| tty-require-dev-recommend-mixed | MATCH | 0 | 0 |
| tty-require-dev-recommend-update | MATCH | 0 | 0 |
| tty-require-move-default | MATCH | 0 | 0 |
| tty-require-move-yes | MATCH | 0 | 0 |
| tty-require-move-no-no | MATCH | 0 | 0 |
| tty-require-move-dev | MATCH | 0 | 0 |
| tty-require-search | MATCH | 137 | 137 |
| tty-require-search-none | MATCH | 0 | 0 |
| tty-require-search-eof | MATCH | 1 | 1 |
| tty-require-virtual | MATCH | 0 | 0 |
| tty-require-virtual-default | MATCH | 0 | 0 |
| tty-require-similar | MATCH | 0 | 0 |
| tty-require-similar-index | MATCH | 0 | 0 |
| tty-require-similar-empty | MATCH | 1 | 1 |
| tty-run-script-select | MATCH | 0 | 0 |
| tty-run-script-select-tab | MATCH | 0 | 0 |
| tty-run-script-select-invalid | MATCH | 1 | 1 |
| tty-run-script-select-empty | MATCH | 1 | 1 |
| tty-run-script-arg | MATCH | 0 | 0 |
| tty-run-script-no-scripts | MATCH | 1 | 1 |
| tty-exec-select | MATCH | 0 | 0 |
| tty-exec-select-name | MATCH | 0 | 0 |
| tty-exec-select-invalid | MATCH | 1 | 1 |
| tty-exec-no-bins | MATCH | 1 | 1 |
| tty-exec-list | MATCH | 0 | 0 |
| tty-create-project-ask | MATCH | 0 | 0 |
| tty-create-project-ask-default | MATCH | 0 | 0 |
| tty-create-project-ask-eof | MATCH | 1 | 1 |
| tty-create-project-ask-no-package | MATCH | 1 | 1 |
| tty-create-project-vcs-yes | MATCH | 0 | 0 |
| tty-create-project-vcs-default | MATCH | 0 | 0 |
| tty-create-project-vcs-no | MATCH | 0 | 0 |
| tty-create-project-vcs-remove-flag | MATCH | 0 | 0 |
| tty-plugin-allow-yes | MATCH | 0 | 0 |
| tty-plugin-allow-no | MATCH | 0 | 0 |
| tty-plugin-allow-discard | MATCH | 0 | 0 |
| tty-plugin-allow-help | MATCH | 0 | 0 |
| tty-plugin-allow-eof | MATCH | 1 | 1 |
| tty-plugin-allow-existing | MATCH | 0 | 0 |
| tty-parent-dir-yes | MATCH | 0 | 0 |
| tty-parent-dir-no | MATCH | 3 | 3 |
| tty-parent-dir-noninteractive | MATCH | 3 | 3 |
| tty-parent-dir-init-exempt | MATCH | 0 | 0 |
| tty-discard-update-yes | MATCH | 0 | 0 |
| tty-discard-update-no | MATCH | 1 | 1 |
| tty-discard-update-help | MATCH | 0 | 0 |
| tty-discard-update-view | MATCH | 0 | 0 |
| tty-discard-update-stash | MATCH | 1 | 1 |
| tty-discard-update-verbose | MATCH | 0 | 0 |
| tty-discard-update-no-interaction | MATCH | 1 | 1 |
| tty-discard-remove-yes | MATCH | 0 | 0 |
| tty-discard-remove-stash-is-help | MATCH | 1 | 1 |
| tty-discard-reinstall | MATCH | 0 | 0 |
| tty-discard-reinstall-no | MATCH | 1 | 1 |
| tty-discard-untracked-only | MATCH | 0 | 0 |
| tty-auth-basic-store-yes | MATCH | 2 | 2 |
| tty-auth-basic-store-no | MATCH | 2 | 2 |
| tty-auth-basic-store-default | MATCH | 2 | 2 |
| tty-auth-basic-store-invalid | MATCH | 2 | 2 |
| tty-auth-basic-store-true | MATCH | 2 | 2 |
| tty-auth-basic-store-false | MATCH | 2 | 2 |
| tty-auth-basic-403 | MATCH | 2 | 2 |
| tty-auth-basic-wrong | MATCH | 100 | 100 |
| tty-auth-basic-empty | MATCH | 100 | 100 |
| tty-auth-basic-eof | MATCH | 1 | 1 |
| tty-auth-basic-eof-password | MATCH | 1 | 1 |
| tty-auth-basic-vv | MATCH | 2 | 2 |
| tty-auth-basic-local-authjson | MATCH | 2 | 2 |
| tty-auth-basic-no-interaction | MATCH | 100 | 100 |
| tty-auth-github-token | MATCH | 0 | 0 |
| tty-auth-github-token-vv | MATCH | 0 | 0 |
| tty-auth-github-empty-token | MATCH | 100 | 100 |
| tty-auth-github-invalid-token | MATCH | 100 | 100 |
| tty-auth-github-eof | MATCH | 1 | 1 |
| tty-auth-github-local-yes | MATCH | 0 | 0 |
| tty-auth-github-local-no | MATCH | 0 | 0 |
| tty-auth-github-ratelimit | MATCH | 100 | 100 |
| tty-auth-github-no-interaction | MATCH | 1 | 1 |
| tty-auth-gitlab-password | MATCH | 0 | 0 |
| tty-auth-gitlab-bad-then-good | MATCH | 0 | 0 |
| tty-auth-gitlab-local-yes | MATCH | 0 | 0 |
| tty-auth-gitlab-eof | MATCH | 1 | 1 |
| tty-auth-gitlab-no-interaction | MATCH | 1 | 1 |
| tty-auth-gitlab-refresh | MATCH | 0 | 0 |
| tty-auth-gitlab-refresh-vv | MATCH | 0 | 0 |
| tty-auth-gitlab-refresh-missing | MATCH | 1 | 1 |
| tty-auth-gitlab-refresh-rejected | MATCH | 1 | 1 |
| tty-auth-bitbucket-consumer | MATCH | 0 | 0 |
| tty-auth-bitbucket-consumer-vv | MATCH | 0 | 0 |
| tty-auth-bitbucket-invalid | MATCH | 100 | 100 |
| tty-auth-bitbucket-no-key | MATCH | 100 | 100 |
| tty-auth-bitbucket-no-secret | MATCH | 100 | 100 |
| tty-auth-git-store-yes | MATCH | 0 | 0 |
| tty-auth-git-store-no | MATCH | 0 | 0 |
| tty-auth-git-store-true-v | MATCH | 0 | 0 |
| tty-auth-git-wrong | MATCH | 1 | 1 |
| tty-auth-git-default-user | MATCH | 1 | 1 |
| tty-auth-git-eof | MATCH | 1 | 1 |
| tty-auth-git-no-interaction | MATCH | 1 | 1 |
| tty-auth-git-github | MATCH | 0 | 0 |
| tty-auth-git-github-no-token | MATCH | 1 | 1 |
| tty-auth-git-gitlab | MATCH | 0 | 0 |
| tty-auth-git-bitbucket | MATCH | 0 | 0 |
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
| vcs-forge-github-ratelimit | MATCH | 1 | 1 |
| vcs-forge-github-ratelimit-auth | MATCH | 1 | 1 |
| vcs-forge-github-enterprise | MATCH | 0 | 0 |
| vcs-forge-github-not-domain | MATCH | 0 | 0 |
| vcs-forge-git-url-credentials | MATCH | 1 | 1 |
| vcs-forge-git-http-credentials | MATCH | 0 | 0 |
| vcs-forge-composer-401 | MATCH | 100 | 100 |
| vcs-forge-composer-403 | MATCH | 100 | 100 |
| vcs-forge-composer-401-basic | MATCH | 2 | 2 |
| vcs-forge-composer-403-bearer | MATCH | 2 | 2 |
| vcs-forge-composer-401-wrong | MATCH | 100 | 100 |
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
| vcs-detect-svn-scheme | MATCH | 1 | 1 |
| vcs-detect-svn-host | MATCH | 1 | 1 |
| vcs-detect-hg-ssh | MATCH | 1 | 1 |
| vcs-detect-fossil | MATCH | 1 | 1 |
| vcs-detect-perforce | MATCH | 1 | 1 |
| vcs-detect-git-ssh | MATCH | 1 | 1 |
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
| vcs-svn-update | MATCH | 0 | 0 |
| vcs-svn-update-v | MATCH | 0 | 0 |
| vcs-svn-update-vv | MATCH | 0 | 0 |
| vcs-svn-detect | MATCH | 1 | 1 |
| vcs-svn-detect-serve | MATCH | 0 | 0 |
| vcs-svn-abs-path | MATCH | 0 | 0 |
| vcs-svn-trunk-url | MATCH | 0 | 0 |
| vcs-svn-dev-trunk | MATCH | 0 | 0 |
| vcs-svn-branch | MATCH | 0 | 0 |
| vcs-svn-show | MATCH | 0 | 0 |
| vcs-svn-package-path | MATCH | 0 | 0 |
| vcs-svn-no-tags-path | MATCH | 0 | 0 |
| vcs-svn-custom-trunk | MATCH | 0 | 0 |
| vcs-svn-missing | MATCH | 1 | 1 |
| vcs-svn-cache | MATCH | 0 | 0 |
| vcs-svn-serve | MATCH | 0 | 0 |
| vcs-svn-serve-secure | MATCH | 1 | 1 |
| vcs-svn-serve-secure-domain | MATCH | 0 | 0 |
| vcs-svn-auth-noninteractive | MATCH | 1 | 1 |
| vcs-svn-auth-http-basic | MATCH | 0 | 0 |
| vcs-svn-auth-wrong | MATCH | 1 | 1 |
| vcs-svn-auth-url | MATCH | 0 | 0 |
| vcs-svn-auth-no-cache | MATCH | 0 | 0 |
| vcs-svn-auth-tty | MATCH | 0 | 0 |
| vcs-svn-auth-tty-wrong | MATCH | 1 | 1 |
| source-svn-install | MATCH | 0 | 0 |
| source-svn-install-v | MATCH | 0 | 0 |
| source-svn-update | MATCH | 0 | 0 |
| source-svn-update-v | KNOWN-DIFF | 1 | 1 |
| source-svn-to-trunk | MATCH | 0 | 0 |
| source-svn-local-changes | MATCH | 1 | 1 |
| source-svn-discard-true | MATCH | 0 | 0 |
| source-svn-discard-stash | MATCH | 1 | 1 |
| source-svn-uninstall-changes | MATCH | 1 | 1 |
| source-svn-uninstall-discard | MATCH | 0 | 0 |
| source-svn-remove | MATCH | 0 | 0 |
| source-svn-status-clean | MATCH | 0 | 0 |
| source-svn-status | MATCH | 1 | 1 |
| source-svn-status-v | MATCH | 1 | 1 |
| source-svn-tty-discard-yes | MATCH | 0 | 0 |
| source-svn-tty-discard-no | MATCH | 1 | 1 |
| source-svn-lock-reference | MATCH | 0 | 0 |
| vcs-hg-update | MATCH | 0 | 0 |
| vcs-hg-update-vv | MATCH | 0 | 0 |
| vcs-hg-detect | MATCH | 0 | 0 |
| vcs-hg-dev-default | MATCH | 0 | 0 |
| vcs-hg-branch | MATCH | 0 | 0 |
| vcs-hg-bookmark | MATCH | 0 | 0 |
| vcs-hg-show | MATCH | 0 | 0 |
| vcs-hg-serve | MATCH | 0 | 0 |
| vcs-hg-serve-vv | MATCH | 0 | 0 |
| vcs-hg-serve-detect | MATCH | 0 | 0 |
| vcs-hg-serve-cached | MATCH | 0 | 0 |
| vcs-hg-serve-secure | MATCH | 100 | 100 |
| vcs-hg-serve-missing | MATCH | 1 | 1 |
| vcs-hg-serve-auth | MATCH | 1 | 1 |
| vcs-hg-missing | MATCH | 1 | 1 |
| source-hg-install | MATCH | 0 | 0 |
| source-hg-install-serve | MATCH | 0 | 0 |
| source-hg-update | MATCH | 0 | 0 |
| source-hg-update-v | MATCH | 0 | 0 |
| source-hg-downgrade-v | MATCH | 0 | 0 |
| source-hg-to-dev | MATCH | 0 | 0 |
| source-hg-local-changes | MATCH | 1 | 1 |
| source-hg-discard-true | MATCH | 1 | 1 |
| source-hg-uninstall-changes | MATCH | 1 | 1 |
| source-hg-remove | MATCH | 0 | 0 |
| source-hg-status-clean | MATCH | 0 | 0 |
| source-hg-status | MATCH | 1 | 1 |
| source-hg-status-v | MATCH | 1 | 1 |
| vcs-fossil-update | MATCH | 1 | 1 |
| vcs-fossil-update-vv | DIFF | 1 | 1 |
| vcs-fossil-detect | DIFF | 1 | 1 |
| vcs-fossil-checkout-dir | MATCH | 0 | 0 |
| vcs-fossil-dev-trunk | MATCH | 1 | 1 |
| vcs-fossil-branch | MATCH | 1 | 1 |
| vcs-fossil-show | MATCH | 0 | 0 |
| vcs-fossil-serve | MATCH | 1 | 1 |
| vcs-fossil-serve-cached | ERROR | 0 | 0 |
| vcs-fossil-serve-secure | MATCH | 100 | 100 |
| vcs-fossil-missing | MATCH | 1 | 1 |
| source-fossil-install | MATCH | 1 | 1 |
| source-fossil-update | ERROR | 0 | 0 |
| source-fossil-update-v | ERROR | 0 | 0 |
| source-fossil-to-dev | ERROR | 0 | 0 |
| source-fossil-local-changes | ERROR | 0 | 0 |
| source-fossil-uninstall-changes | ERROR | 0 | 0 |
| source-fossil-remove | ERROR | 0 | 0 |
| source-fossil-status-clean | ERROR | 0 | 0 |
| source-fossil-status-v | ERROR | 0 | 0 |
| vcs-perforce-update | MATCH | 0 | 0 |
| vcs-perforce-update-vv | MATCH | 0 | 0 |
| vcs-perforce-detect | MATCH | 0 | 0 |
| vcs-perforce-detect-p4-url | MATCH | 1 | 1 |
| vcs-perforce-dev | MATCH | 0 | 0 |
| vcs-perforce-show | MATCH | 0 | 0 |
| vcs-perforce-env-user | MATCH | 0 | 0 |
| vcs-perforce-login | MATCH | 0 | 0 |
| vcs-perforce-login-wrong | MATCH | 1 | 1 |
| vcs-perforce-login-env | MATCH | 0 | 0 |
| vcs-perforce-stream | MATCH | 0 | 0 |
| vcs-perforce-no-depot | MATCH | 1 | 1 |
| vcs-perforce-bad-server | MATCH | 1 | 1 |
| vcs-perforce-invalid-port | MATCH | 1 | 1 |
| source-perforce-install | MATCH | 0 | 0 |
| source-perforce-update | MATCH | 0 | 0 |
| source-perforce-update-v | MATCH | 0 | 0 |
| source-perforce-status | MATCH | 0 | 0 |
| source-perforce-remove | MATCH | 0 | 0 |

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

## plugin2-api-autoload

```
exit code: oracle=0 go=1
stdout:
- (oracle) classmap 8: Acme\Log\Logger,Acme\Log\LoggerInterface,App\Kernel,Composer\InstalledVersions,Plugfix\Api\ApiCommand,Plugfix\Api\LoggingZipDownloader,Plugfix\Api\Plugin,Plugfix\Api\Provider
stderr:
+ (go) In <location>:
+ (go)                                                                         
+ (go)   Call to undefined method GocomposerPluginHost\Services::filterSpec()  
+ (go)                                                                         
+ (go) 
+ (go) plugfix:api <what>
+ (go) 
+ (go) 
project tree:
+ (go) F vendor/composer/autoload_classmap.php 644 375085d0d6690301
- (oracle) F vendor/composer/autoload_classmap.php 644 4d70f24d28620c3f
+ (go) F vendor/composer/autoload_real.php 644 b2a11d3a6f497f86
+ (go) F vendor/composer/autoload_static.php 644 84f0d58a397c626a
- (oracle) F vendor/composer/autoload_real.php 644 d539165e56fab939
- (oracle) F vendor/composer/autoload_static.php 644 62f246111fb4945c
vendor/composer/autoload_classmap.php:
- (oracle)     'Acme\\Log\\Logger' => $vendorDir . '/acme/log/src/Logger.php',
- (oracle)     'Acme\\Log\\LoggerInterface' => $vendorDir . '/acme/log/src/LoggerInterface.php',
- (oracle)     'App\\Kernel' => $baseDir . '/src/Kernel.php',
- (oracle)     'Plugfix\\Api\\ApiCommand' => $vendorDir . '/plugfix/api-plugin/src/ApiCommand.php',
- (oracle)     'Plugfix\\Api\\LoggingZipDownloader' => $vendorDir . '/plugfix/api-plugin/src/LoggingZipDownloader.php',
- (oracle)     'Plugfix\\Api\\Plugin' => $vendorDir . '/plugfix/api-plugin/src/Plugin.php',
- (oracle)     'Plugfix\\Api\\Provider' => $vendorDir . '/plugfix/api-plugin/src/Provider.php',
vendor/composer/autoload_real.php:
- (oracle)         $loader->setClassMapAuthoritative(true);
vendor/composer/autoload_static.php:
- (oracle)         'Acme\\Log\\Logger' => __DIR__ . '/..' . '/acme/log/src/Logger.php',
- (oracle)         'Acme\\Log\\LoggerInterface' => __DIR__ . '/..' . '/acme/log/src/LoggerInterface.php',
- (oracle)         'App\\Kernel' => __DIR__ . '/../..' . '/src/Kernel.php',
- (oracle)         'Plugfix\\Api\\ApiCommand' => __DIR__ . '/..' . '/plugfix/api-plugin/src/ApiCommand.php',
- (oracle)         'Plugfix\\Api\\LoggingZipDownloader' => __DIR__ . '/..' . '/plugfix/api-plugin/src/LoggingZipDownloader.php',
- (oracle)         'Plugfix\\Api\\Plugin' => __DIR__ . '/..' . '/plugfix/api-plugin/src/Plugin.php',
- (oracle)         'Plugfix\\Api\\Provider' => __DIR__ . '/..' . '/plugfix/api-plugin/src/Provider.php',
```

## term-href-audit

```
stdout:
+ (go) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/pksa-test-0001\e\PKSA-test-0001\e]8;;\e\                                                                   |
+ (go) |\e[32m CVE               \e[39m| \e]8;;https://www.cve.org/cverecord?id=cve-2024-0001\e\CVE-2024-0001\e]8;;\e\                                                                    |
- (oracle) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/PKSA-test-0001\e\PKSA-test-0001\e]8;;\e\                                                                   |
- (oracle) |\e[32m CVE               \e[39m| \e]8;;https://www.cve.org/CVERecord?id=CVE-2024-0001\e\CVE-2024-0001\e]8;;\e\                                                                    |
+ (go) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/pksa-test-0002\e\PKSA-test-0002\e]8;;\e\                                                                   |
- (oracle) |\e[32m Advisory ID       \e[39m| \e]8;;https://packagist.org/security-advisories/PKSA-test-0002\e\PKSA-test-0002\e]8;;\e\                                                                   |
```

## term-href-audit-plain

```
stdout:
+ (go) Advisory ID: \e]8;;https://packagist.org/security-advisories/pksa-test-0001\e\PKSA-test-0001\e]8;;\e\
+ (go) CVE: \e]8;;https://www.cve.org/cverecord?id=cve-2024-0001\e\CVE-2024-0001\e]8;;\e\
- (oracle) Advisory ID: \e]8;;https://packagist.org/security-advisories/PKSA-test-0001\e\PKSA-test-0001\e]8;;\e\
- (oracle) CVE: \e]8;;https://www.cve.org/CVERecord?id=CVE-2024-0001\e\CVE-2024-0001\e]8;;\e\
+ (go) Advisory ID: \e]8;;https://packagist.org/security-advisories/pksa-test-0002\e\PKSA-test-0002\e]8;;\e\
- (oracle) Advisory ID: \e]8;;https://packagist.org/security-advisories/PKSA-test-0002\e\PKSA-test-0002\e]8;;\e\
```

## term-progress-update

```
stdout:
+ (go)  0/6 [>---------------------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32m1 package suggestions were added by new dependencies, use `composer suggest` to see details.\e[39m
- (oracle)  0/6 [>---------------------------]   0%\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32m1 package suggestions were added by new dependencies, use `composer suggest` to see details.\e[39m
```

## term-progress-require

```
stdout:
+ (go)  0/6 [>---------------------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>---------------------------]   0%\e[1G\e[2K 6/6 [============================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-very-verbose

```
stdout:
+ (go)  0/6 [>---------------------------]   0% <elapsed>/< 1 sec\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [============================] 100% <elapsed>/< 1 sec\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>---------------------------]   0% <elapsed>/< 1 sec\e[1G\e[2K 6/6 [============================] 100% <elapsed>/< 1 sec\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-narrow

```
stdout:
+ (go)  0/6 [>-----------------]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [==================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle)  0/6 [>-----------------]   0%\e]7501;state=working:app=composer:progress=83\e\\e[1G\e[2K 6/6 [==================] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
```

## term-progress-tiny

```
exit code: oracle=0 go=2
stdout:
+ (go)  0/6 []   0%\e]7501;state=working:app=composer:progress=0\e\panic: strings: negative Repeat count
+ (go) 
+ (go) goroutine 1 [running]:
+ (go) strings.Repeat({0x843028?, 0x5903a5?}, 0x3?)
+ (go) 	strings/strings.go:609 +0x4a0
+ (go) gocomposer/internal/cio.(*ProgressBar).bar(0x4ff42e30c160)
+ (go) 	gocomposer/internal/cio/progress.go:223 +0xb0
+ (go) gocomposer/internal/cio.(*ProgressBar).render.func1({0x5adafb, 0x5})
+ (go) 	gocomposer/internal/cio/progress.go:171 +0x9c
+ (go) regexp.(*Regexp).ReplaceAllStringFunc.func1({0x4ff42e350368, 0x6, 0x8}, {0x4ff42e350380?, 0x0?, 0x0?})
+ (go) 	regexp/regexp.go:577 +0x70
+ (go) regexp.(*Regexp).replaceAll(0x4ff42e116c80, {0x0, 0x0, 0x0}, {0x5adae9, 0x26}, 0x2, 0x4ff42e23b198)
+ (go) 	regexp/regexp.go:615 +0x2d4
+ (go) regexp.(*Regexp).ReplaceAllStringFunc(0xeb7a8a54d061?, {0x5adae9?, 0x4ff42e23b238?}, 0x12f534?)
+ (go) 	regexp/regexp.go:576 +0x54
+ (go) gocomposer/internal/cio.(*ProgressBar).render(...)
+ (go) 	gocomposer/internal/cio/progress.go:166
+ (go) gocomposer/internal/cio.(*ProgressBar).buildLine(0x4ff42e30c160)
+ (go) 	gocomposer/internal/cio/progress.go:157 +0x44
+ (go) gocomposer/internal/cio.(*ProgressBar).display(0x4ff42e30c160)
+ (go) 	gocomposer/internal/cio/progress.go:151 +0x44
+ (go) gocomposer/internal/cio.(*ProgressBar).SetProgress(0x4ff42e30c160, 0x4ff42e098380?)
+ (go) 	gocomposer/internal/cio/progress.go:85 +0x110
+ (go) gocomposer/internal/cio.(*ProgressBar).Finish(0x4ff42e23b340?)
+ (go) 	gocomposer/internal/cio/progress.go:110 +0x48
+ (go) gocomposer/internal/composer.(*InstallationManager).waitOnJobs(0x4ff42e2457c0, 0xd18628?, 0x6, 0x6, 0x4ff42e0cc6c0)
+ (go) 	gocomposer/internal/composer/install.go:1513 +0x1e8
+ (go) gocomposer/internal/composer.(*InstallationManager).downloadAndExecuteBatch(0x4ff42e2457c0, 0x4ff42e0cf220, {0x4ff42e0cc600, 0x6, 0xca5d20?}, 0x1, 0x1, 0x0, {0x5905dc?, 0x7659c?, ...})
+ (go) 	gocomposer/internal/composer/install.go:1071 +0x45c
+ (go) gocomposer/internal/composer.(*InstallationManager).Execute(0x4ff42e2457c0, 0x4ff42e0cf220, {0x4ff42e4141b0, 0x6, 0x6}, 0x1, 0x1, 0x0)
+ (go) 	gocomposer/internal/composer/install.go:1011 +0x314
+ (go) gocomposer/internal/composer.(*Installer).doInstall(0x4ff42e2f50e0, 0x4ff42e0cf220, 0x0)
+ (go) 	gocomposer/internal/composer/installer.go:992 +0x11d8
+ (go) gocomposer/internal/composer.(*Installer).Run(0x4ff42e2f50e0)
+ (go) 	gocomposer/internal/composer/installer.go:163 +0x3a4
+ (go) gocomposer/internal/cli.installCommand.func1(0x4ff42e58c240)
+ (go) 	gocomposer/internal/cli/commands.go:160 +0x154
+ (go) gocomposer/internal/cli.(*App).run(0x4ff42e58c0c0, 0x4ff42e5900a0, {0x4ff42e11e050?, 0x1?, 0x1})
+ (go) 	gocomposer/internal/cli/app.go:519 +0x96c
+ (go) gocomposer/internal/cli.(*App).Run(0x4ff42e58c0c0, {0x4ff42e11e050, 0x1, 0x1})
+ (go) 	gocomposer/internal/cli/app.go:244 +0x22c
+ (go) main.main()
+ (go) 	gocomposer/cmd/gocomposer/main.go:12 +0x54
+ (go) 
- (oracle)  0/6 [>]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [=] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K  - Installing \e[32macme/log\e[39m (\e[33m1.1.0\e[39m): Extracting archive
- (oracle)   - Installing \e[32macme/util\e[39m (\e[33m1.1.0\e[39m): Extracting archive
- (oracle)   - Installing \e[32macme/cli\e[39m (\e[33m1.0.0\e[39m): Extracting archive
- (oracle)   - Installing \e[32macme/devtools\e[39m (\e[33m1.0.0\e[39m): Extracting archive
- (oracle)   - Installing \e[32macme/log-impl\e[39m (\e[33m1.0.0\e[39m): Extracting archive
- (oracle)   - Installing \e[32macme/http\e[39m (\e[33m1.2.0\e[39m): Extracting archive
- (oracle)  0/6 [>]   0%\e]7501;state=working:app=composer:progress=83\e\\e[1G\e[2K 6/6 [=] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\\e[1G\e[2K\e[32mGenerating autoload files\e[39m
- (oracle) \e[32m1 package you are using is looking for funding.\e[39m
- (oracle) \e[32mUse the `composer fund` command to find out more!\e[39m
- (oracle) \e]7501;state=done:app=composer\e\
project tree:
- (oracle) D vendor/acme/cli/bin
- (oracle) D vendor/acme/cli/lib
- (oracle) D vendor/acme/cli/lib/AcmeCli
- (oracle) D vendor/acme/devtools/src
- (oracle) D vendor/acme/http/src
- (oracle) D vendor/acme/log-impl/src
- (oracle) D vendor/acme/log/src
- (oracle) D vendor/acme/util/lib
- (oracle) D vendor/acme/util/src
- (oracle) D vendor/bin
+ (go) F vendor/composer/tmp-738a6cd8d0c7a4e5f0499adbcf340e06.zip 644 69c81f62851beea9
+ (go) F vendor/composer/tmp-8a72f9edf6428127608b91f3d9c762ad.zip 644 3a6f208a987e171a
+ (go) F vendor/composer/tmp-95a583271b74f229331637fdd9b899c5.zip 644 16528c7ac0e58d80
+ (go) F vendor/composer/tmp-acf0a99ece71eb607ba6a1a979aa98d7.zip 644 1e7379706bf83ed0
+ (go) F vendor/composer/tmp-bebfe48fd4f809fde3eef82cb658a2a8.zip 644 a39ee1bf5095ad72
+ (go) F vendor/composer/tmp-e54550513a119e4a4980947557c52d60.zip 644 e3ac8085b239bbc6
- (oracle) F vendor/acme/cli/bin/acme 755 802bc66258c986f2
- (oracle) F vendor/acme/cli/bin/acme-sh 755 9f2574ae8fe83002
- (oracle) F vendor/acme/cli/composer.json 644 382e59b9dededdc3
- (oracle) F vendor/acme/cli/lib/AcmeCli/Runner.php 644 da07eeda61fd87ba
- (oracle) F vendor/acme/devtools/composer.json 644 8deb9729ab693910
- (oracle) F vendor/acme/devtools/src/Dumper.php 644 89f8cbaefda7a093
- (oracle) F vendor/acme/http/composer.json 644 dff053c25f488c21
- (oracle) F vendor/acme/http/src/Client.php 644 138fee57b895c83e
- (oracle) F vendor/acme/log-impl/composer.json 644 bd9c251ce1f7b2fc
- (oracle) F vendor/acme/log-impl/src/Impl.php 644 fd59a8b43f5299f1
- (oracle) F vendor/acme/log/composer.json 644 b83d42805152f028
- (oracle) F vendor/acme/log/src/Logger.php 644 13ef68bda8291625
- (oracle) F vendor/acme/log/src/LoggerInterface.php 644 beb8d938c4444c81
- (oracle) F vendor/acme/util/composer.json 644 ed976ef852108f1f
- (oracle) F vendor/acme/util/lib/legacy.php 644 aa94d8c20c2ac438
- (oracle) F vendor/acme/util/src/Str.php 644 bd3cbab92c3b093f
- (oracle) F vendor/acme/util/src/functions.php 644 f00c9f042b8e0674
- (oracle) F vendor/autoload.php 644 1353eb0f89b5ba0c
- (oracle) F vendor/bin/acme 755 0ec11a39f30b01c4
- (oracle) F vendor/bin/acme-sh 755 9a2820eb637c9499
- (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
- (oracle) F vendor/composer/InstalledVersions.php 644 1e4525198301d1ed
- (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
- (oracle) F vendor/composer/autoload_classmap.php 644 d795431bedb6f53e
- (oracle) F vendor/composer/autoload_files.php 644 0e951f0edb6be02b
- (oracle) F vendor/composer/autoload_namespaces.php 644 7f8af46ba6702f8f
- (oracle) F vendor/composer/autoload_psr4.php 644 23ac8d6d915546dc
- (oracle) F vendor/composer/autoload_real.php 644 aaaddfe74f1f126e
- (oracle) F vendor/composer/autoload_static.php 644 c870c5a9ec15a5b5
- (oracle) F vendor/composer/installed.json 644 1950f68107e3c5fe
- (oracle) F vendor/composer/installed.php 644 01774ecdf9dbdcc3
- (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
autoload probe:
- (oracle) class AcmeAbstractThing: autoloaded
- (oracle) class AcmeCli_Runner: autoloaded
- (oracle) class AcmeSuit: autoloaded
- (oracle) class AcmeThingInterface: autoloaded
- (oracle) class AcmeThingTrait: autoloaded
- (oracle) class Acme\DevTools\Dumper: autoloaded
- (oracle) class Acme\Http\Client: autoloaded
- (oracle) class Acme\LogImpl\Impl: autoloaded
- (oracle) class Acme\Log\Logger: autoloaded
- (oracle) class Acme\Log\LoggerInterface: autoloaded
- (oracle) class Acme\Util\Str: autoloaded
- (oracle) class Acme_Legacy_Helper: autoloaded
- (oracle) class App\Kernel: autoloaded
- (oracle) function acme_util_hello: hello from acme/util
- (oracle) root: test/project 1.0.0+no-version-set dev=true
- (oracle) installed: acme/cli '1.0.0' '94cee099e231284142446146f8f9844a40206270' dev=true
- (oracle) installed: acme/devtools '1.0.0' '4107978d337771591283b614c4075e14c77e9a4c' dev=false
- (oracle) installed: acme/http '1.2.0' 'ce7601587b04bbcb0a274777c6c4a15fe82a9467' dev=true
- (oracle) installed: acme/log '1.1.0' '58b0f69fc5f96981e5f7fae32d592232a36012c2' dev=true
- (oracle) installed: acme/log-impl '1.0.0' 'be55dcec40f479aa52942c4aaa3f393405333530' dev=true
- (oracle) installed: acme/util '1.1.0' '010715fb4511de3a64a1da0e47aa42e3b0104c1c' dev=true
- (oracle) installed: psr/log-implementation NULL NULL dev=true
- (oracle) installed: test/project '1.0.0+no-version-set' NULL dev=true
- (oracle) bin acme mode=0755
- (oracle)   exit=0 out=hello from acme/util
- (oracle) bin acme-sh mode=0755
- (oracle)   exit=0 out=acme-sh says hi
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
- (oracle)     'App\\' => array($baseDir . '/src'),
- (oracle)     'Acme\\Util\\' => array($vendorDir . '/acme/util/src'),
- (oracle)     'Acme\\Log\\' => array($vendorDir . '/acme/log/src'),
- (oracle)     'Acme\\LogImpl\\' => array($vendorDir . '/acme/log-impl/src'),
- (oracle)     'Acme\\Http\\' => array($vendorDir . '/acme/http/src'),
- (oracle)     'Acme\\DevTools\\' => array($vendorDir . '/acme/devtools/src'),
- (oracle) );
- (oracle) 
vendor/composer/autoload_real.php:
- (oracle) <?php
- (oracle) // autoload_real.php @generated by Composer
- (oracle) 
- (oracle) class ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca
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
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca', 'loadClassLoader'), true, true);
- (oracle)         self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit17b0dbd82a93c6a1761bf6ffafc822ca', 'loadClassLoader'));
- (oracle) 
- (oracle)         require __DIR__ . '/autoload_static.php';
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::getInitializer($loader));
- (oracle) 
- (oracle)         $loader->register(true);
- (oracle) 
- (oracle)         $filesToLoad = \Composer\Autoload\ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$files;
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
- (oracle) class ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca
- (oracle) {
- (oracle)     public static $files = array (
- (oracle)         '571dbb8238a79d30de0279182c1a3513' => __DIR__ . '/..' . '/acme/util/src/functions.php',
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixLengthsPsr4 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'App\\' => 4,
- (oracle)             'Acme\\Util\\' => 10,
- (oracle)             'Acme\\Log\\' => 9,
- (oracle)             'Acme\\LogImpl\\' => 13,
- (oracle)             'Acme\\Http\\' => 10,
- (oracle)             'Acme\\DevTools\\' => 14,
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixDirsPsr4 = array (
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
- (oracle)         'Acme\\LogImpl\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/log-impl/src',
- (oracle)         ),
- (oracle)         'Acme\\Http\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/http/src',
- (oracle)         ),
- (oracle)         'Acme\\DevTools\\' =>
- (oracle)         array (
- (oracle)             0 => __DIR__ . '/..' . '/acme/devtools/src',
- (oracle)         ),
- (oracle)     );
- (oracle) 
- (oracle)     public static $prefixesPsr0 = array (
- (oracle)         'A' =>
- (oracle)         array (
- (oracle)             'AcmeCli_' =>
- (oracle)             array (
- (oracle)                 0 => __DIR__ . '/..' . '/acme/cli/lib',
- (oracle)             ),
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
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixDirsPsr4;
- (oracle)             $loader->prefixesPsr0 = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$prefixesPsr0;
- (oracle)             $loader->classMap = ComposerStaticInit17b0dbd82a93c6a1761bf6ffafc822ca::$classMap;
- (oracle) 
- (oracle)         }, null, ClassLoader::class);
- (oracle)     }
- (oracle) }
- (oracle) 
vendor/composer/installed.json:
- (oracle) {
- (oracle)     "packages": [
- (oracle)         {
- (oracle)             "name": "acme/cli",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-cli-1.0.0.zip",
- (oracle)                 "reference": "94cee099e231284142446146f8f9844a40206270",
- (oracle)                 "shasum": "19122d362bcf3265801223914bcc367b825b3a00"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "acme/util": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "bin": [
- (oracle)                 "bin/acme",
- (oracle)                 "bin/acme-sh"
- (oracle)             ],
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-0": {
- (oracle)                     "AcmeCli_": "lib/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "description": "CLI tool",
- (oracle)             "install-path": "../acme/cli"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "acme/devtools",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-devtools-1.0.0.zip",
- (oracle)                 "reference": "4107978d337771591283b614c4075e14c77e9a4c",
- (oracle)                 "shasum": "6165796917c9d36a7319eb43e0e5319a0f231cd9"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "acme/log": "^1.0 || ^2.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\DevTools\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "description": "Dev tools",
- (oracle)             "install-path": "../acme/devtools"
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "acme/http",
- (oracle)             "version": "1.2.0",
- (oracle)             "version_normalized": "1.2.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-http-1.2.0.zip",
- (oracle)                 "reference": "ce7601587b04bbcb0a274777c6c4a15fe82a9467",
- (oracle)                 "shasum": "9f24d825f6e1e2e9a6757bafc819af8a3011eeac"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "acme/util": "^1.0",
- (oracle)                 "php": "^7.4 || ^8.0",
- (oracle)                 "psr/log-implementation": "^1.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\Http\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "description": "HTTP client",
- (oracle)             "funding": [
- (oracle)                 {
- (oracle)                     "type": "github",
- (oracle)                     "url": "https://github.com/sponsors/acme"
- (oracle)                 }
- (oracle)             ],
- (oracle)             "install-path": "../acme/http"
- (oracle)         },
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
- (oracle)             "name": "acme/log-impl",
- (oracle)             "version": "1.0.0",
- (oracle)             "version_normalized": "1.0.0.0",
- (oracle)             "dist": {
- (oracle)                 "type": "zip",
- (oracle)                 "url": "http://127.0.0.1:8899/dists/acme-log-impl-1.0.0.zip",
- (oracle)                 "reference": "be55dcec40f479aa52942c4aaa3f393405333530",
- (oracle)                 "shasum": "d97ff4a53fa4a053860a2f51825fb249e5eded98"
- (oracle)             },
- (oracle)             "require": {
- (oracle)                 "acme/log": "^1.0"
- (oracle)             },
- (oracle)             "provide": {
- (oracle)                 "psr/log-implementation": "1.0.0"
- (oracle)             },
- (oracle)             "time": "2024-01-02T03:04:05+00:00",
- (oracle)             "type": "library",
- (oracle)             "installation-source": "dist",
- (oracle)             "autoload": {
- (oracle)                 "psr-4": {
- (oracle)                     "Acme\\LogImpl\\": "src/"
- (oracle)                 }
- (oracle)             },
- (oracle)             "description": "Provides psr/log-implementation",
- (oracle)             "install-path": "../acme/log-impl"
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
- (oracle)         }
- (oracle)     ],
- (oracle)     "dev": true,
- (oracle)     "dev-package-names": [
- (oracle)         "acme/devtools"
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
- (oracle)         'acme/cli' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '94cee099e231284142446146f8f9844a40206270',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/cli',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'acme/devtools' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => '4107978d337771591283b614c4075e14c77e9a4c',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/devtools',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => true,
- (oracle)         ),
- (oracle)         'acme/http' => array(
- (oracle)             'pretty_version' => '1.2.0',
- (oracle)             'version' => '1.2.0.0',
- (oracle)             'reference' => 'ce7601587b04bbcb0a274777c6c4a15fe82a9467',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/http',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'acme/log' => array(
- (oracle)             'pretty_version' => '1.1.0',
- (oracle)             'version' => '1.1.0.0',
- (oracle)             'reference' => '58b0f69fc5f96981e5f7fae32d592232a36012c2',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/log',
- (oracle)             'aliases' => array(),
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'acme/log-impl' => array(
- (oracle)             'pretty_version' => '1.0.0',
- (oracle)             'version' => '1.0.0.0',
- (oracle)             'reference' => 'be55dcec40f479aa52942c4aaa3f393405333530',
- (oracle)             'type' => 'library',
- (oracle)             'install_path' => __DIR__ . '/../acme/log-impl',
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
- (oracle)             'dev_requirement' => false,
- (oracle)         ),
- (oracle)         'psr/log-implementation' => array(
- (oracle)             'dev_requirement' => false,
- (oracle)             'provided' => array(
- (oracle)                 0 => '1.0.0',
- (oracle)             ),
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

## curl-reset-debug

```
stderr:
- (oracle) Executing command (CWD): 'git' '--version'
- (oracle) Executing command (/work/project): 'git' 'branch' '-a' '--no-color' '--no-abbrev' '-v'
- (oracle) Executing command (/work/project): 'git' 'describe' '--exact-match' '--tags'
- (oracle) Executing command (/work/project): 'git' 'rev-list' '--no-commit-header' '--format=%H' '-n1' 'HEAD' '--no-show-signature'
- (oracle) Executing command (/work/project): 'hg' 'branch'
- (oracle) Executing command (/work/project): 'fossil' 'branch' 'list'
- (oracle) Executing command (/work/project): 'fossil' 'tag' 'list'
- (oracle) Executing command (/work/project): 'svn' 'info' '--xml'
- (oracle) Downloading http://127.0.0.1:9203/packages.json
- (oracle)   [Composer\Downloader\TransportException (100)]                               
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Util/Http/CurlDownloader.php:498
- (oracle)  Composer\Util\Http\CurlDownloader->tick() at /oracle-source/src/Composer/Util/HttpDownloader.php:406
- (oracle)  Composer\Util\HttpDownloader->countActiveJobs() at /oracle-source/src/Composer/Util/HttpDownloader.php:379
- (oracle)  Composer\Util\HttpDownloader->wait() at /oracle-source/src/Composer/Util/HttpDownloader.php:112
- (oracle)  Composer\Util\HttpDownloader->get() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1808
- (oracle)  Composer\Repository\ComposerRepository->fetchFile() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1514
- (oracle)  Composer\Repository\ComposerRepository->loadRootServerFile() at /oracle-source/src/Composer/Repository/ComposerRepository.php:1136
- (oracle)  Composer\Repository\ComposerRepository->hasProviders() at /oracle-source/src/Composer/Repository/ComposerRepository.php:571
- (oracle)  Composer\Repository\ComposerRepository->loadPackages() at /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:458
- (oracle)  Composer\DependencyResolver\PoolBuilder->loadPackagesMarkedForLoading() at /oracle-source/src/Composer/DependencyResolver/PoolBuilder.php:294
- (oracle)  Composer\DependencyResolver\PoolBuilder->buildPool() at /oracle-source/src/Composer/Repository/RepositorySet.php:379
- (oracle)  Composer\Repository\RepositorySet->createPool() at /oracle-source/src/Composer/Installer.php:537
- (oracle)  Composer\Installer->doUpdate() at /oracle-source/src/Composer/Installer.php:327
- (oracle)  Composer\Installer->run() at /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  Composer\Command\UpdateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## root-list

```
stdout:
+ (go) gocomposer version 0.1.0-dev (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle)    ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) Composer version 2.11.999-dev+source @release_date@
+ (go)   -h, --help                Display help for the given command. When no command is given display help for the list command
+ (go)   -q, --quiet               Do not output any message
+ (go)   -v|-vv|-vvv, --verbose    Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   -V, --version             Display this application version
+ (go)       --ansi                Force (or disable --no-ansi) ANSI output
+ (go)   -n, --no-interaction      Do not ask any interactive question
+ (go)       --profile             Display timing and memory usage information
+ (go)       --no-plugins          Whether to disable plugins.
+ (go)       --no-scripts          Skips the execution of all scripts defined in composer.json file.
+ (go)   -d, --working-dir         If specified, use the given directory as working directory.
+ (go)       --no-cache            Prevent use of the cache
- (oracle)   -h, --help                     Display help for the given command. When no command is given display help for the list command
- (oracle)   -q, --quiet                    Do not output any message
- (oracle)   -V, --version                  Display this application version
- (oracle)       --ansi|--no-ansi           Force (or disable --no-ansi) ANSI output
- (oracle)   -n, --no-interaction           Do not ask any interactive question
- (oracle)       --profile                  Display timing and memory usage information
- (oracle)       --no-plugins               Whether to disable plugins.
- (oracle)       --no-scripts               Skips the execution of all scripts defined in composer.json file.
- (oracle)   -d, --working-dir=WORKING-DIR  If specified, use the given directory as working directory.
- (oracle)       --no-cache                 Prevent use of the cache
- (oracle)   -v|vv|vvv, --verbose           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   browse               Opens the package's repository URL or homepage in your browser
- (oracle)   browse               [home] Opens the package's repository URL or homepage in your browser
+ (go)   clear-cache          Clears composer's internal package cache
- (oracle)   clear-cache          [clearcache|cc] Clears composer's internal package cache
+ (go)   depends              Shows which packages cause the given package to be installed
- (oracle)   depends              [why] Shows which packages cause the given package to be installed
+ (go)   dump-autoload        Dumps the autoloader
- (oracle)   dump-autoload        [dumpautoload] Dumps the autoloader
+ (go)   install              Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   install              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   prohibits            Shows which packages prevent the given package from being installed
- (oracle)   prohibits            [why-not] Shows which packages prevent the given package from being installed
+ (go)   remove               Removes a package from the require or require-dev
+ (go)   repository           Manages repositories
+ (go)   require              Adds required packages to your composer.json and installs them
+ (go)   run-script           Runs the scripts defined in composer.json
- (oracle)   remove               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   repository           [repo] Manages repositories
- (oracle)   require              [r] Adds required packages to your composer.json and installs them
- (oracle)   run-script           [run] Runs the scripts defined in composer.json
+ (go)   self-update          Updates composer.phar to the latest version
+ (go)   show                 Shows information about packages
- (oracle)   self-update          [selfupdate] Updates composer.phar to the latest version
- (oracle)   show                 [info] Shows information about packages
+ (go)   update               Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   update               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
```

## root-list-allowed

```
stdout:
+ (go) gocomposer version 0.1.0-dev (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle)    ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) Composer version 2.11.999-dev+source @release_date@
+ (go)   -h, --help                Display help for the given command. When no command is given display help for the list command
+ (go)   -q, --quiet               Do not output any message
+ (go)   -v|-vv|-vvv, --verbose    Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   -V, --version             Display this application version
+ (go)       --ansi                Force (or disable --no-ansi) ANSI output
+ (go)   -n, --no-interaction      Do not ask any interactive question
+ (go)       --profile             Display timing and memory usage information
+ (go)       --no-plugins          Whether to disable plugins.
+ (go)       --no-scripts          Skips the execution of all scripts defined in composer.json file.
+ (go)   -d, --working-dir         If specified, use the given directory as working directory.
+ (go)       --no-cache            Prevent use of the cache
- (oracle)   -h, --help                     Display help for the given command. When no command is given display help for the list command
- (oracle)   -q, --quiet                    Do not output any message
- (oracle)   -V, --version                  Display this application version
- (oracle)       --ansi|--no-ansi           Force (or disable --no-ansi) ANSI output
- (oracle)   -n, --no-interaction           Do not ask any interactive question
- (oracle)       --profile                  Display timing and memory usage information
- (oracle)       --no-plugins               Whether to disable plugins.
- (oracle)       --no-scripts               Skips the execution of all scripts defined in composer.json file.
- (oracle)   -d, --working-dir=WORKING-DIR  If specified, use the given directory as working directory.
- (oracle)       --no-cache                 Prevent use of the cache
- (oracle)   -v|vv|vvv, --verbose           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   browse               Opens the package's repository URL or homepage in your browser
- (oracle)   browse               [home] Opens the package's repository URL or homepage in your browser
+ (go)   clear-cache          Clears composer's internal package cache
- (oracle)   clear-cache          [clearcache|cc] Clears composer's internal package cache
+ (go)   depends              Shows which packages cause the given package to be installed
- (oracle)   depends              [why] Shows which packages cause the given package to be installed
+ (go)   dump-autoload        Dumps the autoloader
- (oracle)   dump-autoload        [dumpautoload] Dumps the autoloader
+ (go)   install              Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   install              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   prohibits            Shows which packages prevent the given package from being installed
- (oracle)   prohibits            [why-not] Shows which packages prevent the given package from being installed
+ (go)   remove               Removes a package from the require or require-dev
+ (go)   repository           Manages repositories
+ (go)   require              Adds required packages to your composer.json and installs them
+ (go)   run-script           Runs the scripts defined in composer.json
- (oracle)   remove               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   repository           [repo] Manages repositories
- (oracle)   require              [r] Adds required packages to your composer.json and installs them
- (oracle)   run-script           [run] Runs the scripts defined in composer.json
+ (go)   self-update          Updates composer.phar to the latest version
+ (go)   show                 Shows information about packages
- (oracle)   self-update          [selfupdate] Updates composer.phar to the latest version
- (oracle)   show                 [info] Shows information about packages
+ (go)   update               Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   update               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
```

## root-tty-list-yes

```
stdout:
+ (go) \e]7501;state=working:app=composer\e\\e[32mgocomposer\e[39m version \e[33m0.1.0-dev\e[39m (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle) \e]7501;state=working:app=composer\e\   ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) \e[32mComposer\e[39m version \e[33m2.11.999-dev+source\e[39m @release_date@
+ (go)   \e[32m-h, --help               \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet              \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose   \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version            \e[39m Display this application version
+ (go)   \e[32m    --ansi               \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction     \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile            \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins         \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts         \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir        \e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache           \e[39m Prevent use of the cache
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32mabout               \e[39m Shows a short information about Composer
+ (go)   \e[32marchive             \e[39m Creates an archive of this composer package
+ (go)   \e[32maudit               \e[39m Checks for security vulnerability advisories for installed packages
+ (go)   \e[32mbrowse              \e[39m Opens the package's repository URL or homepage in your browser
+ (go)   \e[32mbump                \e[39m Increases the lower limit of your composer.json requirements to the currently installed versions
+ (go)   \e[32mcheck-platform-reqs \e[39m Check that platform requirements are satisfied
+ (go)   \e[32mclear-cache         \e[39m Clears composer's internal package cache
+ (go)   \e[32mcompletion          \e[39m Dump the shell completion script
+ (go)   \e[32mconfig              \e[39m Sets config options
+ (go)   \e[32mcreate-project      \e[39m Creates new project from a package into given directory
+ (go)   \e[32mdepends             \e[39m Shows which packages cause the given package to be installed
+ (go)   \e[32mdiagnose            \e[39m Diagnoses the system to identify common errors
+ (go)   \e[32mdump-autoload       \e[39m Dumps the autoloader
+ (go)   \e[32mexec                \e[39m Executes a vendored binary/script
+ (go)   \e[32mfund                \e[39m Discover how to help fund the maintenance of your dependencies
+ (go)   \e[32mglobal              \e[39m Allows running commands in the global composer dir ($COMPOSER_HOME)
+ (go)   \e[32mhello               \e[39m Runs the hello script as defined in composer.json
+ (go)   \e[32mhelp                \e[39m Display help for a command
+ (go)   \e[32minit                \e[39m Creates a basic composer.json file in current directory
+ (go)   \e[32minstall             \e[39m Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   \e[32mlicenses            \e[39m Shows information about licenses of dependencies
+ (go)   \e[32mlist                \e[39m List commands
+ (go)   \e[32moutdated            \e[39m Shows a list of installed packages that have updates available, including their latest version
+ (go)   \e[32mpolicy              \e[39m Manages custom dependency policies and their sources
+ (go)   \e[32mprohibits           \e[39m Shows which packages prevent the given package from being installed
+ (go)   \e[32mreinstall           \e[39m Uninstalls and reinstalls the given package names
+ (go)   \e[32mremove              \e[39m Removes a package from the require or require-dev
+ (go)   \e[32mrepository          \e[39m Manages repositories
+ (go)   \e[32mrequire             \e[39m Adds required packages to your composer.json and installs them
+ (go)   \e[32mrun-script          \e[39m Runs the scripts defined in composer.json
+ (go)   \e[32msearch              \e[39m Searches for packages
+ (go)   \e[32mself-update         \e[39m Updates composer.phar to the latest version
+ (go)   \e[32mshow                \e[39m Shows information about packages
+ (go)   \e[32mstatus              \e[39m Shows a list of locally modified packages
+ (go)   \e[32msuggests            \e[39m Shows package suggestions
+ (go)   \e[32mupdate              \e[39m Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
+ (go)   \e[32mvalidate            \e[39m Validates a composer.json and composer.lock
- (oracle)   \e[32mabout\e[39m                Shows a short information about Composer
- (oracle)   \e[32marchive\e[39m              Creates an archive of this composer package
- (oracle)   \e[32maudit\e[39m                Checks for security vulnerability advisories for installed packages
- (oracle)   \e[32mbrowse\e[39m               [home] Opens the package's repository URL or homepage in your browser
- (oracle)   \e[32mbump\e[39m                 Increases the lower limit of your composer.json requirements to the currently installed versions
- (oracle)   \e[32mcheck-platform-reqs\e[39m  Check that platform requirements are satisfied
- (oracle)   \e[32mclear-cache\e[39m          [clearcache|cc] Clears composer's internal package cache
- (oracle)   \e[32mcompletion\e[39m           Dump the shell completion script
- (oracle)   \e[32mconfig\e[39m               Sets config options
- (oracle)   \e[32mcreate-project\e[39m       Creates new project from a package into given directory
- (oracle)   \e[32mdepends\e[39m              [why] Shows which packages cause the given package to be installed
- (oracle)   \e[32mdiagnose\e[39m             Diagnoses the system to identify common errors
- (oracle)   \e[32mdump-autoload\e[39m        [dumpautoload] Dumps the autoloader
- (oracle)   \e[32mexec\e[39m                 Executes a vendored binary/script
- (oracle)   \e[32mfund\e[39m                 Discover how to help fund the maintenance of your dependencies
- (oracle)   \e[32mglobal\e[39m               Allows running commands in the global composer dir ($COMPOSER_HOME)
- (oracle)   \e[32mhello\e[39m                Runs the hello script as defined in composer.json
- (oracle)   \e[32mhelp\e[39m                 Display help for a command
- (oracle)   \e[32minit\e[39m                 Creates a basic composer.json file in current directory
- (oracle)   \e[32minstall\e[39m              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   \e[32mlicenses\e[39m             Shows information about licenses of dependencies
- (oracle)   \e[32mlist\e[39m                 List commands
- (oracle)   \e[32moutdated\e[39m             Shows a list of installed packages that have updates available, including their latest version
- (oracle)   \e[32mpolicy\e[39m               Manages custom dependency policies and their sources
- (oracle)   \e[32mprohibits\e[39m            [why-not] Shows which packages prevent the given package from being installed
- (oracle)   \e[32mreinstall\e[39m            Uninstalls and reinstalls the given package names
- (oracle)   \e[32mremove\e[39m               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   \e[32mrepository\e[39m           [repo] Manages repositories
- (oracle)   \e[32mrequire\e[39m              [r] Adds required packages to your composer.json and installs them
- (oracle)   \e[32mrun-script\e[39m           [run] Runs the scripts defined in composer.json
- (oracle)   \e[32msearch\e[39m               Searches for packages
- (oracle)   \e[32mself-update\e[39m          [selfupdate] Updates composer.phar to the latest version
- (oracle)   \e[32mshow\e[39m                 [info] Shows information about packages
- (oracle)   \e[32mstatus\e[39m               Shows a list of locally modified packages
- (oracle)   \e[32msuggests\e[39m             Shows package suggestions
- (oracle)   \e[32mupdate\e[39m               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   \e[32mvalidate\e[39m             Validates a composer.json and composer.lock
```

## root-tty-list-no

```
stdout:
+ (go) \e]7501;state=working:app=composer\e\\e[32mgocomposer\e[39m version \e[33m0.1.0-dev\e[39m (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle) \e]7501;state=working:app=composer\e\   ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) \e[32mComposer\e[39m version \e[33m2.11.999-dev+source\e[39m @release_date@
+ (go)   \e[32m-h, --help               \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet              \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose   \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version            \e[39m Display this application version
+ (go)   \e[32m    --ansi               \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction     \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile            \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins         \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts         \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir        \e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache           \e[39m Prevent use of the cache
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32mabout               \e[39m Shows a short information about Composer
+ (go)   \e[32marchive             \e[39m Creates an archive of this composer package
+ (go)   \e[32maudit               \e[39m Checks for security vulnerability advisories for installed packages
+ (go)   \e[32mbrowse              \e[39m Opens the package's repository URL or homepage in your browser
+ (go)   \e[32mbump                \e[39m Increases the lower limit of your composer.json requirements to the currently installed versions
+ (go)   \e[32mcheck-platform-reqs \e[39m Check that platform requirements are satisfied
+ (go)   \e[32mclear-cache         \e[39m Clears composer's internal package cache
+ (go)   \e[32mcompletion          \e[39m Dump the shell completion script
+ (go)   \e[32mconfig              \e[39m Sets config options
+ (go)   \e[32mcreate-project      \e[39m Creates new project from a package into given directory
+ (go)   \e[32mdepends             \e[39m Shows which packages cause the given package to be installed
+ (go)   \e[32mdiagnose            \e[39m Diagnoses the system to identify common errors
+ (go)   \e[32mdump-autoload       \e[39m Dumps the autoloader
+ (go)   \e[32mexec                \e[39m Executes a vendored binary/script
+ (go)   \e[32mfund                \e[39m Discover how to help fund the maintenance of your dependencies
+ (go)   \e[32mglobal              \e[39m Allows running commands in the global composer dir ($COMPOSER_HOME)
+ (go)   \e[32mhello               \e[39m Runs the hello script as defined in composer.json
+ (go)   \e[32mhelp                \e[39m Display help for a command
+ (go)   \e[32minit                \e[39m Creates a basic composer.json file in current directory
+ (go)   \e[32minstall             \e[39m Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   \e[32mlicenses            \e[39m Shows information about licenses of dependencies
+ (go)   \e[32mlist                \e[39m List commands
+ (go)   \e[32moutdated            \e[39m Shows a list of installed packages that have updates available, including their latest version
+ (go)   \e[32mpolicy              \e[39m Manages custom dependency policies and their sources
+ (go)   \e[32mprohibits           \e[39m Shows which packages prevent the given package from being installed
+ (go)   \e[32mreinstall           \e[39m Uninstalls and reinstalls the given package names
+ (go)   \e[32mremove              \e[39m Removes a package from the require or require-dev
+ (go)   \e[32mrepository          \e[39m Manages repositories
+ (go)   \e[32mrequire             \e[39m Adds required packages to your composer.json and installs them
+ (go)   \e[32mrun-script          \e[39m Runs the scripts defined in composer.json
+ (go)   \e[32msearch              \e[39m Searches for packages
+ (go)   \e[32mself-update         \e[39m Updates composer.phar to the latest version
+ (go)   \e[32mshow                \e[39m Shows information about packages
+ (go)   \e[32mstatus              \e[39m Shows a list of locally modified packages
+ (go)   \e[32msuggests            \e[39m Shows package suggestions
+ (go)   \e[32mupdate              \e[39m Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
+ (go)   \e[32mvalidate            \e[39m Validates a composer.json and composer.lock
- (oracle)   \e[32mabout\e[39m                Shows a short information about Composer
- (oracle)   \e[32marchive\e[39m              Creates an archive of this composer package
- (oracle)   \e[32maudit\e[39m                Checks for security vulnerability advisories for installed packages
- (oracle)   \e[32mbrowse\e[39m               [home] Opens the package's repository URL or homepage in your browser
- (oracle)   \e[32mbump\e[39m                 Increases the lower limit of your composer.json requirements to the currently installed versions
- (oracle)   \e[32mcheck-platform-reqs\e[39m  Check that platform requirements are satisfied
- (oracle)   \e[32mclear-cache\e[39m          [clearcache|cc] Clears composer's internal package cache
- (oracle)   \e[32mcompletion\e[39m           Dump the shell completion script
- (oracle)   \e[32mconfig\e[39m               Sets config options
- (oracle)   \e[32mcreate-project\e[39m       Creates new project from a package into given directory
- (oracle)   \e[32mdepends\e[39m              [why] Shows which packages cause the given package to be installed
- (oracle)   \e[32mdiagnose\e[39m             Diagnoses the system to identify common errors
- (oracle)   \e[32mdump-autoload\e[39m        [dumpautoload] Dumps the autoloader
- (oracle)   \e[32mexec\e[39m                 Executes a vendored binary/script
- (oracle)   \e[32mfund\e[39m                 Discover how to help fund the maintenance of your dependencies
- (oracle)   \e[32mglobal\e[39m               Allows running commands in the global composer dir ($COMPOSER_HOME)
- (oracle)   \e[32mhello\e[39m                Runs the hello script as defined in composer.json
- (oracle)   \e[32mhelp\e[39m                 Display help for a command
- (oracle)   \e[32minit\e[39m                 Creates a basic composer.json file in current directory
- (oracle)   \e[32minstall\e[39m              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   \e[32mlicenses\e[39m             Shows information about licenses of dependencies
- (oracle)   \e[32mlist\e[39m                 List commands
- (oracle)   \e[32moutdated\e[39m             Shows a list of installed packages that have updates available, including their latest version
- (oracle)   \e[32mpolicy\e[39m               Manages custom dependency policies and their sources
- (oracle)   \e[32mprohibits\e[39m            [why-not] Shows which packages prevent the given package from being installed
- (oracle)   \e[32mreinstall\e[39m            Uninstalls and reinstalls the given package names
- (oracle)   \e[32mremove\e[39m               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   \e[32mrepository\e[39m           [repo] Manages repositories
- (oracle)   \e[32mrequire\e[39m              [r] Adds required packages to your composer.json and installs them
- (oracle)   \e[32mrun-script\e[39m           [run] Runs the scripts defined in composer.json
- (oracle)   \e[32msearch\e[39m               Searches for packages
- (oracle)   \e[32mself-update\e[39m          [selfupdate] Updates composer.phar to the latest version
- (oracle)   \e[32mshow\e[39m                 [info] Shows information about packages
- (oracle)   \e[32mstatus\e[39m               Shows a list of locally modified packages
- (oracle)   \e[32msuggests\e[39m             Shows package suggestions
- (oracle)   \e[32mupdate\e[39m               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   \e[32mvalidate\e[39m             Validates a composer.json and composer.lock
```

## source-svn-update-v

```
stderr:
- (oracle)   [RuntimeException]                                                           
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Downloader/SvnDownloader.php:223
- (oracle)  Composer\Downloader\SvnDownloader->getCommitLogs() at /oracle-source/src/Composer/Downloader/VcsDownloader.php:195
- (oracle)  Composer\Downloader\VcsDownloader->update() at /oracle-source/src/Composer/Downloader/DownloadManager.php:319
- (oracle)  Composer\Downloader\DownloadManager->update() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:313
- (oracle)  Composer\Installer\LibraryInstaller->updateCode() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:177
- (oracle)  Composer\Installer\LibraryInstaller->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:502
- (oracle)  Composer\Installer\InstallationManager->update() at /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  Composer\Installer\InstallationManager->{closure:Composer\Installer\InstallationManager::executeBatch():385}() at /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  React\Promise\Internal\FulfilledPromise->then() at /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  Composer\Installer\InstallationManager->executeBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  Composer\Installer\InstallationManager->downloadAndExecuteBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  Composer\Installer\InstallationManager->execute() at /oracle-source/src/Composer/Installer.php:885
- (oracle)  Composer\Installer->doInstall() at /oracle-source/src/Composer/Installer.php:689
- (oracle)  Composer\Installer->doUpdate() at /oracle-source/src/Composer/Installer.php:327
- (oracle)  Composer\Installer->run() at /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  Composer\Command\UpdateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## vcs-fossil-update-vv

```
stderr:
- (oracle)   [RuntimeException]                        
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Downloader/FossilDownloader.php:117
- (oracle)  Composer\Downloader\FossilDownloader->execute() at /oracle-source/src/Composer/Downloader/FossilDownloader.php:48
- (oracle)  Composer\Downloader\FossilDownloader->doInstall() at /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  Composer\Downloader\VcsDownloader->install() at /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  Composer\Downloader\DownloadManager->install() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  Composer\Installer\LibraryInstaller->installCode() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  Composer\Installer\LibraryInstaller->install() at /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  Composer\Installer\InstallationManager->install() at /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  Composer\Installer\InstallationManager->{closure:Composer\Installer\InstallationManager::executeBatch():385}() at /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  React\Promise\Internal\FulfilledPromise->then() at /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  Composer\Installer\InstallationManager->executeBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  Composer\Installer\InstallationManager->downloadAndExecuteBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  Composer\Installer\InstallationManager->execute() at /oracle-source/src/Composer/Installer.php:885
- (oracle)  Composer\Installer->doInstall() at /oracle-source/src/Composer/Installer.php:689
- (oracle)  Composer\Installer->doUpdate() at /oracle-source/src/Composer/Installer.php:327
- (oracle)  Composer\Installer->run() at /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  Composer\Command\UpdateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## vcs-fossil-detect

```
stderr:
- (oracle)   [RuntimeException]                        
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Downloader/FossilDownloader.php:117
- (oracle)  Composer\Downloader\FossilDownloader->execute() at /oracle-source/src/Composer/Downloader/FossilDownloader.php:48
- (oracle)  Composer\Downloader\FossilDownloader->doInstall() at /oracle-source/src/Composer/Downloader/VcsDownloader.php:137
- (oracle)  Composer\Downloader\VcsDownloader->install() at /oracle-source/src/Composer/Downloader/DownloadManager.php:278
- (oracle)  Composer\Downloader\DownloadManager->install() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:278
- (oracle)  Composer\Installer\LibraryInstaller->installCode() at /oracle-source/src/Composer/Installer/LibraryInstaller.php:149
- (oracle)  Composer\Installer\LibraryInstaller->install() at /oracle-source/src/Composer/Installer/InstallationManager.php:479
- (oracle)  Composer\Installer\InstallationManager->install() at /oracle-source/src/Composer/Installer/InstallationManager.php:386
- (oracle)  Composer\Installer\InstallationManager->{closure:Composer\Installer\InstallationManager::executeBatch():385}() at /oracle-source/vendor/react/promise/src/Internal/FulfilledPromise.php:47
- (oracle)  React\Promise\Internal\FulfilledPromise->then() at /oracle-source/src/Composer/Installer/InstallationManager.php:385
- (oracle)  Composer\Installer\InstallationManager->executeBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:327
- (oracle)  Composer\Installer\InstallationManager->downloadAndExecuteBatch() at /oracle-source/src/Composer/Installer/InstallationManager.php:225
- (oracle)  Composer\Installer\InstallationManager->execute() at /oracle-source/src/Composer/Installer.php:885
- (oracle)  Composer\Installer->doInstall() at /oracle-source/src/Composer/Installer.php:689
- (oracle)  Composer\Installer->doUpdate() at /oracle-source/src/Composer/Installer.php:327
- (oracle)  Composer\Installer->run() at /oracle-source/src/Composer/Command/UpdateCommand.php:308
- (oracle)  Composer\Command\UpdateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## vcs-fossil-serve-cached

```
setup failed: update --no-progress --no-install
Composer could not detect the root package (test/project) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.

In FossilDriver.php line 109:
                                                                               
  Failed to clone http://127.0.0.1:8902/ to repository /work/project/.cache-k  
  eep/repo/http---127-0-0-1-8902-.fossil                                       
                                                                               
  SQLITE_CANTOPEN(14): cannot open file at line 45249 of [20abf1ec10]          
  SQLITE_CANTOPEN(14): os_unix.c:45249: (2) open(/work/project/.cache-keep/re  
  po/http---127-0-0-1-8902-.fossil) - No such file or directory                
  Database error: [/work/project/.cache-keep/repo/http---127-0-0-1-8902-.foss  
  il]: unable to open database file                                            
                                                                               

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-update

```
setup failed: update --no-progress --with vcs3/fossillib:1.0.0
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.0.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.0.0): Cloning 1.0.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.0.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-update-v

```
setup failed: update --no-progress --with vcs3/fossillib:1.0.0
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.0.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.0.0): Cloning 1.0.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.0.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-to-dev

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.1.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.1.0): Cloning 1.1.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.1.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-local-changes

```
setup failed: update --no-progress --with vcs3/fossillib:1.0.0
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.0.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.0.0): Cloning 1.0.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.0.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-uninstall-changes

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.1.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.1.0): Cloning 1.1.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.1.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-remove

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.1.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.1.0): Cloning 1.1.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.1.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-status-clean

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.1.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.1.0): Cloning 1.1.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.1.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```

## source-fossil-status-v

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking vcs3/fossillib (1.1.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing vcs3/fossillib (1.1.0): Cloning 1.1.0
    Install of vcs3/fossillib failed

In FossilDownloader.php line 117:
                                            
  Failed to execute fossil update -- 1.1.0  
                                            
  cannot determine user                     
                                            

update [--with WITH] [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--dev] [--no-dev] [--lock] [--no-install] [--no-audit] [--audit-format AUDIT-FORMAT] [--no-security-blocking] [--no-blocking] [--no-cooldown] [--no-autoloader] [--no-suggest] [--no-progress] [-w|--with-dependencies] [-W|--with-all-dependencies] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--prefer-stable] [--prefer-lowest] [-m|--minimal-changes] [--patch-only] [-i|--interactive] [--root-reqs] [--bump-after-update [BUMP-AFTER-UPDATE]] [--] [<packages>...]


```
