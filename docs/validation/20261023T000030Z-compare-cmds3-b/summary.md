# Composer vs gocomposer comparison

- gocomposer_sha256: `10ca324ff80a10f21c46f78cac848655995f63c93e7b2d69a3c4bf3ec6f22d83`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| reinstall-progress | DIFF | 0 | 0 |
| reinstall-wildcard | MATCH | 0 | 0 |
| reinstall-all-wildcard | MATCH | 0 | 0 |
| reinstall-type | MATCH | 0 | 0 |

## reinstall-progress

```
stderr:
+ (go)     0 [>---------------------------]
- (oracle)     0 [>---------------------------]    0 [->--------------------------]
```
