# Composer vs gocomposer comparison

- gocomposer_sha256: `4a13e0d881aaaa4055e32bc3e4b848f562f664b40052995ff45c5d261b4895c4`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive-root-tar | DIFF | 0 | 0 |
| archive-root-zip | DIFF | 0 | 0 |
| archive-root-tar-gz | DIFF | 0 | 0 |
| archive-root-tar-bz2 | MATCH | 1 | 1 |
| archive-root-bad-format | MATCH | 1 | 1 |
| archive-root-empty-format | MATCH | 1 | 1 |
| archive-root-dir | DIFF | 0 | 0 |
| archive-root-dir-relative | DIFF | 0 | 0 |
| archive-root-file | DIFF | 0 | 0 |
| archive-root-file-subdir | DIFF | 0 | 0 |
| archive-root-ignore-filters | MATCH | 0 | 0 |
| archive-root-named | MATCH | 0 | 0 |
| archive-root-installed | DIFF | 0 | 0 |
| archive-root-existing-archives | DIFF | 0 | 0 |
| archive-root-config-format | DIFF | 0 | 0 |
| archive-root-dir-is-file | MATCH | 1 | 1 |
| archive-root-scripts | DIFF | 0 | 0 |
| archive-package | DIFF | 0 | 1 |
| archive-package-installed | DIFF | 0 | 1 |
| archive-package-constraint | DIFF | 0 | 1 |
| archive-package-stability | DIFF | 0 | 1 |
| archive-package-stability-filtered | MATCH | 1 | 1 |
| archive-package-missing | MATCH | 1 | 1 |
| archive-package-file | MATCH | 1 | 1 |
| archive-package-case | DIFF | 0 | 1 |
| archive-no-json | MATCH | 1 | 1 |

## archive-root-tar

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @9728
+ (go)   zero-block @10240
+ (go)   total=10752
- (oracle)   zero-block @7168
- (oracle)   zero-block @7680
- (oracle)   total=8192
```

## archive-root-zip

```
autoload probe:
+ (go)   local ver=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   local ver=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   local ver=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=942
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1007
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1079
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1179
+ (go)   eocd 504b0506000000000a000a0004030000070500000000
+ (go)   total=2081
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=974
- (oracle)   eocd 504b05060000000007000700130200003a0400000000
- (oracle)   total=1635
```

## archive-root-tar-gz

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @9728
+ (go)   zero-block @10240
+ (go)   total=10752
- (oracle)   zero-block @7168
- (oracle)   zero-block @7680
- (oracle)   total=8192
```

## archive-root-dir

```
autoload probe:
+ (go)   local ver=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   local ver=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   local ver=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=942
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1007
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1079
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1179
+ (go)   eocd 504b0506000000000a000a0004030000070500000000
+ (go)   total=2081
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=974
- (oracle)   eocd 504b05060000000007000700130200003a0400000000
- (oracle)   total=1635
```

## archive-root-dir-relative

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @10240
+ (go)   zero-block @10752
+ (go)   total=11264
- (oracle)   zero-block @7680
- (oracle)   zero-block @8192
- (oracle)   total=8704
```

## archive-root-file

```
autoload probe:
+ (go)   local ver=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   local ver=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   local ver=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=942
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1007
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1079
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1179
+ (go)   eocd 504b0506000000000a000a0004030000070500000000
+ (go)   total=2081
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=874
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=974
- (oracle)   eocd 504b05060000000007000700130200003a0400000000
- (oracle)   total=1635
```

## archive-root-file-subdir

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @10240
+ (go)   zero-block @10752
+ (go)   total=11264
- (oracle)   zero-block @7680
- (oracle)   zero-block @8192
- (oracle)   total=8704
```

## archive-root-installed

```
autoload probe:
+ (go)   local ver=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   local ver=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   local ver=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2345
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2413
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2478
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2550
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2650
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=b53110b cs=190 us=190 name="vendor\/acme\/log\/composer.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2758
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f0f06175 cs=74 us=74 name="vendor\/acme\/log\/src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=3025
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=b86e719a cs=56 us=56 name="vendor\/acme\/log\/src\/LoggerInterface.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=3177
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c71de60e cs=748 us=748 name="vendor\/autoload.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=3320
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=751c4032 cs=16378 us=16378 name="vendor\/composer\/ClassLoader.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=4135
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=776e3c0f cs=17395 us=17395 name="vendor\/composer\/InstalledVersions.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=20592
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=82678656 cs=1068 us=1068 name="vendor\/composer\/LICENSE" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=38072
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=e6bdf14c cs=222 us=222 name="vendor\/composer\/autoload_classmap.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39211
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=74112f8a cs=139 us=139 name="vendor\/composer\/autoload_namespaces.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39518
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=35f3fb20 cs=232 us=232 name="vendor\/composer\/autoload_psr4.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39744
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=32d10b3a cs=1137 us=1137 name="vendor\/composer\/autoload_real.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=40057
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=11b8131 cs=1175 us=1175 name="vendor\/composer\/autoload_static.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=41275
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=76945bf8 cs=1038 us=1038 name="vendor\/composer\/installed.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=42533
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=48f48626 cs=1037 us=1037 name="vendor\/composer\/installed.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=43649
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a0a8a72d cs=917 us=917 name="vendor\/composer\/platform_check.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=44763
+ (go)   eocd 504b0506000000001a001a00ef080000c2b200000000
+ (go)   total=48071
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2345
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2445
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=b53110b cs=190 us=190 name="vendor\/acme\/log\/composer.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2553
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f0f06175 cs=74 us=74 name="vendor\/acme\/log\/src\/Logger.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2820
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=b86e719a cs=56 us=56 name="vendor\/acme\/log\/src\/LoggerInterface.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=2972
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=c71de60e cs=748 us=748 name="vendor\/autoload.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=3115
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=751c4032 cs=16378 us=16378 name="vendor\/composer\/ClassLoader.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=3930
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=776e3c0f cs=17395 us=17395 name="vendor\/composer\/InstalledVersions.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=20387
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=82678656 cs=1068 us=1068 name="vendor\/composer\/LICENSE" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=37867
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=e6bdf14c cs=222 us=222 name="vendor\/composer\/autoload_classmap.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39006
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=74112f8a cs=139 us=139 name="vendor\/composer\/autoload_namespaces.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39313
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=35f3fb20 cs=232 us=232 name="vendor\/composer\/autoload_psr4.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39539
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=32d10b3a cs=1137 us=1137 name="vendor\/composer\/autoload_real.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=39852
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=11b8131 cs=1175 us=1175 name="vendor\/composer\/autoload_static.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=41070
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=76945bf8 cs=1038 us=1038 name="vendor\/composer\/installed.json" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=42328
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=48f48626 cs=1037 us=1037 name="vendor\/composer\/installed.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=43444
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=a0a8a72d cs=917 us=917 name="vendor\/composer\/platform_check.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=44558
- (oracle)   eocd 504b05060000000017001700fe070000f5b100000000
- (oracle)   total=47625
```

## archive-root-existing-archives

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @11776
+ (go)   zero-block @12288
+ (go)   total=12800
- (oracle)   zero-block @9216
- (oracle)   zero-block @9728
- (oracle)   total=10240
```

## archive-root-config-format

```
autoload probe:
+ (go)   local ver=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   local ver=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   local ver=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=c53479c3 cs=7 us=7 name="docs\/index.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=922
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=a5ce2f53 cs=5 us=5 name="docs\/keep.md" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=990
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="empty-parent\/empty\/.keep" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1055
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1127
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1227
+ (go)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="dist\/" extra=6e750e001bdffaa5ff010000000000000000 comment= disk=0 ia=0 ea=0 lo=1335
+ (go)   eocd 504b0506000000000b000b00490300006c0500000000
+ (go)   total=2251
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=f748b50e cs=38 us=38 name="src\/Kernel.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=922
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=1c25ea06 cs=42 us=42 name="src\/Sub\/Helper.php" extra=6e750e00845bb5fca4010000000000000000 comment= disk=0 ia=0 ea=0 lo=1022
- (oracle)   central vm=0 vn=0 flag=2048 meth=0 crc=0 cs=0 us=0 name="dist\/" extra=6e750e001bdffaa5ff010000000000000000 comment= disk=0 ia=0 ea=0 lo=1130
- (oracle)   eocd 504b05060000000008000800580200009f0400000000
- (oracle)   total=1805
```

## archive-root-scripts

```
autoload probe:
+ (go)   entry name="docs\/index.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=4dfc8ac0c2fdb50826dd14a191d4863a0418a6f2
+ (go)   entry name="docs\/keep.md" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=0c251ef2fca485fd117c165c8a20692c0f1c7220
+ (go)   entry name="empty-parent\/empty\/.keep" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030303000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
+ (go)     sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709
+ (go)   zero-block @10240
+ (go)   zero-block @10752
+ (go)   total=11264
- (oracle)   zero-block @7680
- (oracle)   zero-block @8192
- (oracle)   total=8704
```

## archive-package

```
exit code: oracle=0 go=1
stdout:
- (oracle) ../arc/acme-log-f2d33608d115e75f71e7edf16a47742051393c05-zip.tar
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   rename /work/project/vendor/composer/2d82a6ae/acme-log-f2d3360 /tmp/compose  
+ (go)   r_archive047212b36c: invalid cross-device link                               
+ (go)                                                                                
+ (go) 
+ (go) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
+ (go) 
+ (go) 
- (oracle) Created: 
autoload probe:
- (oracle) archive /work/arc/acme-log-f2d33608d115e75f71e7edf16a47742051393c05-zip.tar perms=644
- (oracle)   entry name="composer.json" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303032323300 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=435ad5ac47ec2028dd166661120158fdea96731a
- (oracle)   entry name="src\/Logger.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303031313200 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=b784779c8d754389c8caf4fcb6c87635ebe95370
- (oracle)   entry name="src\/LoggerInterface.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030373000 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=6847ec6934007105b643c61576ecbd50adeb0f5d
- (oracle)   zero-block @3072
- (oracle)   zero-block @3584
- (oracle)   total=4096
```

## archive-package-installed

```
exit code: oracle=0 go=1
stdout:
- (oracle) ../arc/acme-log-58b0f69fc5f96981e5f7fae32d592232a36012c2-zip.zip
stderr:
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Found multiple matches, selected acme/log 2.0.0.
+ (go) Alternatives were acme/log 1.1.0, acme/log 1.0.0, acme/log 1.1.0, acme/log 2.0.0.
+ (go) Please use a more specific constraint to pick a different package.
- (oracle) Found an exact match acme/log 1.1.0.
+ (go)   - Downloading acme/log (2.0.0)
+ (go)   - Installing acme/log (2.0.0): Extracting archive
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   rename /work/project/vendor/composer/87195028/acme-log-4d1056e /tmp/compose  
+ (go)   r_archive5d54c4faa4: invalid cross-device link                               
+ (go)                                                                                
+ (go) 
+ (go) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
+ (go) 
+ (go) 
- (oracle)   - Downloading acme/log (1.1.0)
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle)   - Installing acme/log (1.1.0): Extracting archive
- (oracle) Created: 
autoload probe:
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

## archive-package-constraint

```
exit code: oracle=0 go=1
stdout:
- (oracle) ../arc/acme-util-010715fb4511de3a64a1da0e47aa42e3b0104c1c-zip.tar.gz
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   rename /work/project/vendor/composer/e81bd479/acme-util-010715f /tmp/compos  
+ (go)   er_archive3023095c00: invalid cross-device link                              
+ (go)                                                                                
+ (go) 
+ (go) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
+ (go) 
+ (go) 
- (oracle) Created: 
autoload probe:
- (oracle) archive /work/arc/acme-util-010715fb4511de3a64a1da0e47aa42e3b0104c1c-zip.tar.gz perms=644
- (oracle)   gzip header 1f8b0800000000000003
- (oracle)   entry name="composer.json" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303033313500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=aa16b3227d8c8d4f67b454b1705a31fb16bbce12
- (oracle)   entry name="lib\/legacy.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303033313700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=d318682fb274e729fc37c253ecf93b268dcafdb3
- (oracle)   entry name="src\/Str.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303031373700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=e3bf97fe6b5e841e7cc69c288f135c3b0385c099
- (oracle)   entry name="src\/functions.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303031363600 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=1d3c7726027d411613282f544b0a892acd1839cf
- (oracle)   zero-block @4096
- (oracle)   zero-block @4608
- (oracle)   total=5120
```

## archive-package-stability

```
exit code: oracle=0 go=1
stdout:
- (oracle) ../arc/acme-beta-cd5d2f1c7e04187e8e66c25a11638c76a34c8eb7-zip.tar
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   rename /work/project/vendor/composer/482e8705/acme-beta-cd5d2f1 /tmp/compos  
+ (go)   er_archive91e4d793b0: invalid cross-device link                              
+ (go)                                                                                
+ (go) 
+ (go) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
+ (go) 
+ (go) 
- (oracle) Created: 
autoload probe:
- (oracle) archive /work/arc/acme-beta-cd5d2f1c7e04187e8e66c25a11638c76a34c8eb7-zip.tar perms=644
- (oracle)   entry name="composer.json" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303030353700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=41f7c8be25d6345892996b8f2ac540f060c076dd
- (oracle)   zero-block @1024
- (oracle)   zero-block @1536
- (oracle)   total=2048
```

## archive-package-case

```
exit code: oracle=0 go=1
stdout:
- (oracle) ../arc/acme-util-010715fb4511de3a64a1da0e47aa42e3b0104c1c-zip.tar
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   rename /work/project/vendor/composer/23ac4b77/acme-util-010715f /tmp/compos  
+ (go)   er_archive4f7e605962: invalid cross-device link                              
+ (go)                                                                                
+ (go) 
+ (go) archive [-f|--format FORMAT] [--dir DIR] [--file FILE] [--ignore-filters] [--] [<package> [<version>]]
+ (go) 
+ (go) 
- (oracle) Created: 
autoload probe:
- (oracle) archive /work/arc/acme-util-010715fb4511de3a64a1da0e47aa42e3b0104c1c-zip.tar perms=644
- (oracle)   entry name="composer.json" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303033313500 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=aa16b3227d8c8d4f67b454b1705a31fb16bbce12
- (oracle)   entry name="lib\/legacy.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303033313700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=d318682fb274e729fc37c253ecf93b268dcafdb3
- (oracle)   entry name="src\/Str.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303031373700 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=e3bf97fe6b5e841e7cc69c288f135c3b0385c099
- (oracle)   entry name="src\/functions.php" prefix="" mode=3030303036343400 uid=0000000000000000 gid=0000000000000000 size=303030303030303031363600 type=30 link="" magic=757374617200 ver=3030 uname= gname= dev=00000000000000000000000000000000 pad=000000000000000000000000 mtime-len=11
- (oracle)     sha1=1d3c7726027d411613282f544b0a892acd1839cf
- (oracle)   zero-block @4096
- (oracle)   zero-block @4608
- (oracle)   total=5120
```
