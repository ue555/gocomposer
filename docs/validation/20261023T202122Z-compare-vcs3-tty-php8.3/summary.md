# Composer vs gocomposer comparison

- gocomposer_sha256: `8a8f8d5fb28828fbc3b57432e29d52bc8926fe3ba0db114d8ba3b1345571d910`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.3`
- oracle_image_digest: `sha256:2b1b802373a38be255fdd8574e7461226aeaccbd1e06295cd90c4a0616ab4280`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.3-zip`
- oracle_zip_image_digest: `sha256:f7513be70340918a49f0c6b8ae972d068b0422f51fbb77268d7303a9de4aebd6`
- php_version: `8.3.36`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| source-svn-tty-discard-yes | MATCH | 0 | 0 |
| source-svn-tty-discard-no | MATCH | 1 | 1 |
