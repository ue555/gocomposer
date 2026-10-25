# Composer vs gocomposer comparison

- gocomposer_sha256: `4d0afc98a995df94e74c2167d3bc989df18996f366e48ed2da699b2691d1fbf4`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.3`
- oracle_image_digest: `sha256:91789a6658a5075533a2b26cfbdb2a264ab6051e638872129021bccc27665fce`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.3-zip`
- oracle_zip_image_digest: `sha256:08e40da7145cf527bf595df929ce3394a5ae7bcfd43a8e4bac8ec1243e406bf9`
- php_version: `8.3.36`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive3-zip-7z-crc-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-crc-zipext | MATCH | 0 | 0 |
