# Composer vs gocomposer comparison

- gocomposer_sha256: `7f74f6b2f14a112bfc05a18b26df4c21fab51066cbf126757c65089046ca2193`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
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
