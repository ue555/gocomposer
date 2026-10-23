# Composer vs gocomposer comparison

- gocomposer_sha256: `f0d2d56d3d63e991a2a741acc3767ba76d9c6532e6c4abcf33baf4e2c7d55e4b`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive-root-tar | MATCH | 0 | 0 |
| archive-root-zip | MATCH | 0 | 0 |
| archive-root-tar-gz | MATCH | 0 | 0 |
| archive-root-tar-bz2 | MATCH | 1 | 1 |
| archive-root-bad-format | MATCH | 1 | 1 |
| archive-root-empty-format | MATCH | 1 | 1 |
| archive-root-dir | MATCH | 0 | 0 |
| archive-root-dir-relative | MATCH | 0 | 0 |
| archive-root-file | MATCH | 0 | 0 |
| archive-root-file-subdir | MATCH | 0 | 0 |
| archive-root-ignore-filters | MATCH | 0 | 0 |
| archive-root-named | MATCH | 0 | 0 |
| archive-root-installed | MATCH | 0 | 0 |
| archive-root-existing-archives | MATCH | 0 | 0 |
| archive-root-config-format | MATCH | 0 | 0 |
| archive-root-dir-is-file | MATCH | 1 | 1 |
| archive-root-scripts | MATCH | 0 | 0 |
| archive-package | MATCH | 0 | 0 |
| archive-package-installed | DIFF | 0 | 0 |
| archive-package-constraint | MATCH | 0 | 0 |
| archive-package-stability | MATCH | 0 | 0 |
| archive-package-stability-filtered | MATCH | 1 | 1 |
| archive-package-missing | MATCH | 1 | 1 |
| archive-package-file | MATCH | 1 | 1 |
| archive-package-case | MATCH | 0 | 0 |
| archive-no-json | MATCH | 1 | 1 |

## archive-package-installed

```
stdout:
+ (go) ../arc/acme-log-4d1056e36a149970eff22bce6c72fd8f5cb360f1-zip.zip
- (oracle) ../arc/acme-log-58b0f69fc5f96981e5f7fae32d592232a36012c2-zip.zip
stderr:
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Found multiple matches, selected acme/log 2.0.0.
+ (go) Alternatives were acme/log 1.1.0, acme/log 1.0.0, acme/log 1.1.0, acme/log 2.0.0.
+ (go) Please use a more specific constraint to pick a different package.
- (oracle) Found an exact match acme/log 1.1.0.
+ (go)   - Downloading acme/log (2.0.0)
+ (go)   - Installing acme/log (2.0.0): Extracting archive
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
autoload probe:
+ (go) archive /work/arc/acme-log-4d1056e36a149970eff22bce6c72fd8f5cb360f1-zip.zip perms=644
+ (go)   local ver=0 flag=2048 meth=0 crc=8ecd0e50 cs=147 us=147 name="composer.json" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=e5915504a4439e9ed2d75d7f8535dd9c3352c10b
+ (go)   local ver=0 flag=2048 meth=0 crc=d76ddae8 cs=74 us=74 name="src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=fa0e2296c1595a9e54b2ee77061c063e56772f78
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=8ecd0e50 cs=147 us=147 name="composer.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=0
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=d76ddae8 cs=74 us=74 name="src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=208
+ (go)   eocd 504b050600000000020002009b000000580100000000
+ (go)   total=521
- (oracle) archive /work/arc/acme-log-58b0f69fc5f96981e5f7fae32d592232a36012c2-zip.zip perms=644
- (oracle)   local ver=0 flag=2048 meth=0 crc=b53110b cs=190 us=190 name="composer.json" extra=6e750e00845bb5fca4010000000000000000
- (oracle)     sha1=06ba33838c61ff04892bcdea9c695f13ed21df81
- (oracle)   local ver=0 flag=2048 meth=0 crc=f0f06175 cs=74 us=74 name="src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000
- (oracle)     sha1=591ca1ecc6d2a078274743cb21468f8911e136dd
- (oracle)   local ver=0 flag=2048 meth=0 crc=b86e719a cs=56 us=56 name="src\/LoggerInterface.php" extra=6e750e00845bb5fca4010000000000000000
- (oracle)     sha1=6847ec6934007105b643c61576ecbd50adeb0f5d
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=b53110b cs=190 us=190 name="composer.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=0
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f0f06175 cs=74 us=74 name="src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=251
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=b86e719a cs=56 us=56 name="src\/LoggerInterface.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=387
- (oracle)   eocd 504b05060000000003000300f2000000020200000000
- (oracle)   total=778
```
