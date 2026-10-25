# Composer vs gocomposer comparison

- gocomposer_sha256: `a2388567361d201cefeda37fb0896192e93ee2dd3a20d2f174178dec88ff7430`
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
| plugin3-dl-order-parallel | DIFF | 0 | 0 |
| plugin3-dl-order-serial | MATCH | 0 | 0 |
| plugin3-dl-order-throw | DIFF | 1 | 1 |

## plugin3-dl-order-parallel

```
stdout:
+ (go) [order] post plugfix3/c-fast exists=true
- (oracle) [order] post plugfix3/c-fast exists=true
```

## plugin3-dl-order-throw

```
stdout:
+ (go) [order] post plugfix3/b-mid exists=true
- (oracle) [order] post plugfix3/b-mid exists=true
```
