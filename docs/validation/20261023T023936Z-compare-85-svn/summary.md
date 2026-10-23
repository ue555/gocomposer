# Composer vs gocomposer comparison

- gocomposer_sha256: `9057ffdf50c8fd556a4423c3f6157d524f323ad9cb7659cac1a45f2fd250e5a4`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.5`
- oracle_image_digest: `sha256:8f997f444c893bd29a5143015dfa981662afe3e936c43d87adf8416d912f92e2`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.5-zip`
- oracle_zip_image_digest: `sha256:c0ec954e3ac6a8eeb8f9a5fb7ec90943a1414e08c86351127f1c15969e234528`
- php_version: `8.5.12`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| vcs-type-hg | MATCH | 1 | 1 |
| vcs-type-svn | MATCH | 1 | 1 |
| vcs-type-fossil | MATCH | 1 | 1 |
| vcs-type-perforce | MATCH | 1 | 1 |
| vcs-type-pear | MATCH | 1 | 1 |
| vcs-type-forgejo | MATCH | 1 | 1 |
| vcs-type-github-local | MATCH | 1 | 1 |
