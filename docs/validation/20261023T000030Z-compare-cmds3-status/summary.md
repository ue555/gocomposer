# Composer vs gocomposer comparison

- gocomposer_sha256: `7326390045d7352c2fdcd0a399303710897a2360eebb96224cb1a681e4e3dd3c`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
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
