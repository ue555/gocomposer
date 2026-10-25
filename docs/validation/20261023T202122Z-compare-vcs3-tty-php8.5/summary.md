# Composer vs gocomposer comparison

- gocomposer_sha256: `8a8f8d5fb28828fbc3b57432e29d52bc8926fe3ba0db114d8ba3b1345571d910`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.5`
- oracle_image_digest: `sha256:a9ce5d9c43bbc3ee3f07c81df52b3ffa6cf2de4beb6c803cf9ab8e3af016104f`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.5-zip`
- oracle_zip_image_digest: `sha256:9932fd56504e6222ca9dda0b0fc4d094e940ea6daaa0e18cfb89f98a2a24dcba`
- php_version: `8.5.12`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| source-svn-tty-discard-yes | MATCH | 0 | 0 |
| source-svn-tty-discard-no | MATCH | 1 | 1 |
