# Composer vs gocomposer comparison

- gocomposer_sha256: `c6c9ddacc6b49ef00d26c9cdd00f379b36593ce05ea28bbd85dafa310e9867d8`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
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
| create-project-no-args | MATCH | 0 | 0 |
| create-project-no-args-no-json | MATCH | 1 | 1 |
