# Composer vs gocomposer comparison

- gocomposer_sha256: `9de9e3fe9f1d667ab69be8aa6a7afd7e7e7f527ca12f183d94992edfb9a6a002`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
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
