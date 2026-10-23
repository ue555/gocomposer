# Composer vs gocomposer comparison

- gocomposer_sha256: `5ab34596da54fe4337559b4d9623e1b6627eaeb544111cdbeab8a81ae13aa943`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive-empty-tar | MATCH | 0 | 0 |
| archive-empty-zip | MATCH | 0 | 0 |
| archive-empty-tar-gz | MATCH | 0 | 0 |
| archive-only-empty-dirs | MATCH | 0 | 0 |
| archive-globstar-attributes | DIFF | 1 | 0 |
| archive-vendor-dir | MATCH | 0 | 0 |

## archive-globstar-attributes

```
exit code: oracle=1 go=0
stdout:
+ (go) ../arc/test-project-1.0.0+no-version-set.zip
stderr:
+ (go) Created: 
- (oracle) 
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   preg_match(): Compilation failed: missing terminating ] for character class  
- (oracle)    at offset 21                                                                
- (oracle)                                                                                
- (oracle) 
- (oracle) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
- (oracle) 
- (oracle) 
autoload probe:
+ (go) archive /work/arc/test-project-1.0.0+no-version-set.zip perms=644
+ (go)   local ver=0 flag=2048 meth=0 crc=878f944b cs=7 us=7 name=".hidden" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=333610b9a7955f54efdddea14634ab93d77a0d64
+ (go)   local ver=0 flag=2048 meth=0 crc=16e490b3 cs=7 us=7 name="README.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0753dda729217d9bd892d252bdd35f2ee6774a5b
+ (go)   local ver=0 flag=2048 meth=0 crc=7233fbdb cs=20 us=20 name="bin\/tool" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0d5f24983b6bf4de4df8206e45d9804d4959d66c
+ (go)   local ver=0 flag=2048 meth=0 crc=b8da4aa5 cs=445 us=445 name="composer.json" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=f2fc2eef688fb70fcfc3156c9a3ef7f0bee4b074
+ (go)   local ver=0 flag=2048 meth=0 crc=71857850 cs=5 us=5 name="deep\/dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd\/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee\/file.txt" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4aaf8f1bc676d6c000dff7c99abc088fb55d8758
+ (go)   local ver=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=059758be0d5079727a6a7b50280a53ae617f5c76
+ (go)   local ver=0 flag=2048 meth=0 crc=d6d8d23e cs=6 us=6 name="src\/Sub\/notes.txt" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=b9350f295d01cbab7589bc1c6850a621e86992ed
+ (go)   local ver=0 flag=2048 meth=0 crc=31ae3ce cs=6 us=6 name="tests\/FooTest.php" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=10ecb38ffc99bee1f8be94f945b8a87a71faef9c
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=878f944b cs=7 us=7 name=".hidden" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=0
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=16e490b3 cs=7 us=7 name="README.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=62
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=7233fbdb cs=20 us=20 name="bin\/tool" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=126
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=b8da4aa5 cs=445 us=445 name="composer.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=202
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=71857850 cs=5 us=5 name="deep\/dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd\/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee\/file.txt" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=708
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=896
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=d6d8d23e cs=6 us=6 name="src\/Sub\/notes.txt" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=996
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=31ae3ce cs=6 us=6 name="tests\/FooTest.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1067
+ (go)   eocd 504b05060000000008000800dc020000720400000000
+ (go)   total=1892
```
