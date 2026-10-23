# Composer vs gocomposer comparison

- gocomposer_sha256: `92af647b627b5b935ba1f5a220f79e864f51331068a663d098c4f6e5a798f01e`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| validate-bad-constraint | DIFF | 1 | 2 |

## validate-bad-constraint

```
exit code: oracle=1 go=2
stderr:
+ (go) runtime: goroutine stack exceeds 1000000000-byte limit
+ (go) runtime: sp=0x7fa971560380 stack=[0x7fa971560000, 0x7fa991560000]
+ (go) fatal error: stack overflow
+ (go) runtime stack:
+ (go) runtime.throw({0x3e4cb6?, 0x6084c?})
+ (go) 	runtime/panic.go:1243 +0x38 fp=0xfffff691d0d0 sp=0xfffff691d0a0 pc=0x91108
+ (go) runtime.newstack()
+ (go) 	runtime/stack.go:1207 +0x480 fp=0xfffff691d200 sp=0xfffff691d0d0 pc=0x74940
+ (go) runtime.morestack()
+ (go) 	runtime/asm_arm64.s:487 +0x70 fp=0xfffff691d200 sp=0xfffff691d200 pc=0x96190
+ (go) 
+ (go) goroutine 1 gp=0x7fa9510ce1e0 m=0 mp=0x87da00 [running]:
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:60 +0xe0 fp=0x7fa971560380 sp=0x7fa971560380 pc=0x35b490
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715603c0 sp=0x7fa971560380 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560400 sp=0x7fa9715603c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560440 sp=0x7fa971560400 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560480 sp=0x7fa971560440 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715604c0 sp=0x7fa971560480 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560500 sp=0x7fa9715604c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560540 sp=0x7fa971560500 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560580 sp=0x7fa971560540 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715605c0 sp=0x7fa971560580 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560600 sp=0x7fa9715605c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560640 sp=0x7fa971560600 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560680 sp=0x7fa971560640 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715606c0 sp=0x7fa971560680 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560700 sp=0x7fa9715606c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560740 sp=0x7fa971560700 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560780 sp=0x7fa971560740 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715607c0 sp=0x7fa971560780 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560800 sp=0x7fa9715607c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560840 sp=0x7fa971560800 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560880 sp=0x7fa971560840 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715608c0 sp=0x7fa971560880 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560900 sp=0x7fa9715608c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560940 sp=0x7fa971560900 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560980 sp=0x7fa971560940 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa9715609c0 sp=0x7fa971560980 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560a00 sp=0x7fa9715609c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560a40 sp=0x7fa971560a00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560a80 sp=0x7fa971560a40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560ac0 sp=0x7fa971560a80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560b00 sp=0x7fa971560ac0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560b40 sp=0x7fa971560b00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560b80 sp=0x7fa971560b40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560bc0 sp=0x7fa971560b80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560c00 sp=0x7fa971560bc0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560c40 sp=0x7fa971560c00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560c80 sp=0x7fa971560c40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560cc0 sp=0x7fa971560c80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560d00 sp=0x7fa971560cc0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560d40 sp=0x7fa971560d00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560d80 sp=0x7fa971560d40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560dc0 sp=0x7fa971560d80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560e00 sp=0x7fa971560dc0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560e40 sp=0x7fa971560e00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560e80 sp=0x7fa971560e40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560ec0 sp=0x7fa971560e80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560f00 sp=0x7fa971560ec0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560f40 sp=0x7fa971560f00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560f80 sp=0x7fa971560f40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa971560fc0 sp=0x7fa971560f80 pc=0x35b3d4
+ (go) ...8388471 frames elided...
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155edc0 sp=0x7fa99155ed80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ee00 sp=0x7fa99155edc0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ee40 sp=0x7fa99155ee00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ee80 sp=0x7fa99155ee40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155eec0 sp=0x7fa99155ee80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ef00 sp=0x7fa99155eec0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ef40 sp=0x7fa99155ef00 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155ef80 sp=0x7fa99155ef40 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155efc0 sp=0x7fa99155ef80 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f000 sp=0x7fa99155efc0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f040 sp=0x7fa99155f000 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f080 sp=0x7fa99155f040 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f0c0 sp=0x7fa99155f080 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f100 sp=0x7fa99155f0c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f140 sp=0x7fa99155f100 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f180 sp=0x7fa99155f140 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f1c0 sp=0x7fa99155f180 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f200 sp=0x7fa99155f1c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f240 sp=0x7fa99155f200 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f280 sp=0x7fa99155f240 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f2c0 sp=0x7fa99155f280 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f300 sp=0x7fa99155f2c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f340 sp=0x7fa99155f300 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f380 sp=0x7fa99155f340 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f3c0 sp=0x7fa99155f380 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f400 sp=0x7fa99155f3c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f440 sp=0x7fa99155f400 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f480 sp=0x7fa99155f440 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f4c0 sp=0x7fa99155f480 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f500 sp=0x7fa99155f4c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f540 sp=0x7fa99155f500 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f580 sp=0x7fa99155f540 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f5c0 sp=0x7fa99155f580 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f600 sp=0x7fa99155f5c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f640 sp=0x7fa99155f600 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f680 sp=0x7fa99155f640 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f6c0 sp=0x7fa99155f680 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f700 sp=0x7fa99155f6c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f740 sp=0x7fa99155f700 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f780 sp=0x7fa99155f740 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f7c0 sp=0x7fa99155f780 pc=0x35b3d4
+ (go) gocomposer/internal/composer.userErr({0x7d4a88, 0x7fa9513999c0})
+ (go) 	gocomposer/internal/composer/composer.go:61 +0x24 fp=0x7fa99155f800 sp=0x7fa99155f7c0 pc=0x35b3d4
+ (go) gocomposer/internal/composer.LoadRootPackage(0x7fa9511c3200?, {0x7fa9510f207c, 0xd}, 0x7fa95143a630, 0x7fa9513feaf0)
+ (go) 	gocomposer/internal/composer/rootloader.go:60 +0x1560 fp=0x7fa99155fba0 sp=0x7fa99155f800 pc=0x376c80
+ (go) gocomposer/internal/composer.Create(0x7fa9513feaf0, {0x59?, 0x87?, 0x3e?, {0x0?, 0xffffa49544e8?}})
+ (go) 	gocomposer/internal/composer/composer.go:167 +0x5e4 fp=0x7fa99155fd10 sp=0x7fa99155fba0 pc=0x35bcd4
+ (go) gocomposer/internal/cli.(*Context).composerFor(0x7fa95143a1b0, 0x1)
+ (go) 	gocomposer/internal/cli/app.go:628 +0x98 fp=0x7fa99155fd50 sp=0x7fa99155fd10 pc=0x381fd8
+ (go) gocomposer/internal/cli.(*App).run(0x7fa95143a060, 0x7fa9513feaf0, {0x7fa95119c050?, 0x7fa951182648?, 0x1})
+ (go) 	gocomposer/internal/cli/app.go:354 +0x710 fp=0x7fa99155fe00 sp=0x7fa99155fd50 pc=0x37f740
+ (go) gocomposer/internal/cli.(*App).Run(0x7fa95143a060, {0x7fa95119c050, 0x1, 0x1})
+ (go) 	gocomposer/internal/cli/app.go:210 +0xfc fp=0x7fa99155fe70 sp=0x7fa99155fe00 pc=0x37ea5c
+ (go) main.main()
+ (go) 	gocomposer/cmd/gocomposer/main.go:12 +0x54 fp=0x7fa99155fea0 sp=0x7fa99155fe70 pc=0x38be14
+ (go) runtime.main()
+ (go) 	runtime/proc.go:302 +0x380 fp=0x7fa99155ffd0 sp=0x7fa99155fea0 pc=0x5bb20
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa99155ffd0 sp=0x7fa99155ffd0 pc=0x97f94
+ (go) 
+ (go) goroutine 2 gp=0x7fa9510ce780 m=nil [force gc (idle)]:
+ (go) runtime.gopark(0x0?, 0x0?, 0x0?, 0x0?, 0x0?)
+ (go) 	runtime/proc.go:474 +0xb8 fp=0x7fa951128f90 sp=0x7fa951128f70 pc=0x91218
+ (go) runtime.goparkunlock(...)
+ (go) 	runtime/proc.go:480
+ (go) runtime.forcegchelper()
+ (go) 	runtime/proc.go:387 +0xb4 fp=0x7fa951128fd0 sp=0x7fa951128f90 pc=0x5be04
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa951128fd0 sp=0x7fa951128fd0 pc=0x97f94
+ (go) created by runtime.init.7 in goroutine 1
+ (go) 	runtime/proc.go:375 +0x24
+ (go) 
+ (go) goroutine 3 gp=0x7fa9510ced20 m=nil [GC sweep wait]:
+ (go) runtime.gopark(0x0?, 0x0?, 0x0?, 0x0?, 0x0?)
+ (go) 	runtime/proc.go:474 +0xb8 fp=0x7fa951129770 sp=0x7fa951129750 pc=0x91218
+ (go) runtime.goparkunlock(...)
+ (go) 	runtime/proc.go:480
+ (go) runtime.bgsweep(0x7fa951148000)
+ (go) 	runtime/mgcsweep.go:279 +0x9c fp=0x7fa9511297b0 sp=0x7fa951129770 pc=0x4686c
+ (go) runtime.gcenable.gowrap1()
+ (go) 	runtime/mgc.go:214 +0x20 fp=0x7fa9511297d0 sp=0x7fa9511297b0 pc=0x88260
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa9511297d0 sp=0x7fa9511297d0 pc=0x97f94
+ (go) created by runtime.gcenable in goroutine 1
+ (go) 	runtime/mgc.go:214 +0x6c
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Link constraint in test/project requires > acme/log should be a valid versi  
- (oracle)   on constraint, got "~>1"                                                     
- (oracle)                                                                                
+ (go) goroutine 4 gp=0x7fa9510cef00 m=nil [GC scavenge wait]:
+ (go) runtime.gopark(0x7fa951148000?, 0x412d80?, 0x1?, 0x0?, 0x7fa9510cef00?)
+ (go) 	runtime/proc.go:474 +0xb8 fp=0x7fa951129f60 sp=0x7fa951129f40 pc=0x91218
+ (go) runtime.goparkunlock(...)
+ (go) 	runtime/proc.go:480
+ (go) runtime.(*scavengerState).park(0x87c9c0)
+ (go) 	runtime/mgcscavenge.go:425 +0x5c fp=0x7fa951129f90 sp=0x7fa951129f60 pc=0x444dc
+ (go) runtime.bgscavenge(0x7fa951148000)
+ (go) 	runtime/mgcscavenge.go:653 +0x44 fp=0x7fa951129fb0 sp=0x7fa951129f90 pc=0x449e4
+ (go) runtime.gcenable.gowrap2()
+ (go) 	runtime/mgc.go:215 +0x20 fp=0x7fa951129fd0 sp=0x7fa951129fb0 pc=0x88220
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa951129fd0 sp=0x7fa951129fd0 pc=0x97f94
+ (go) created by runtime.gcenable in goroutine 1
+ (go) 	runtime/mgc.go:215 +0xac
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Could not parse version constraint ~>1: Invalid operator "~>", you probably  
- (oracle)    meant to use the "~" operator                                               
- (oracle)                                                                                
+ (go) goroutine 17 gp=0x7fa9511843c0 m=nil [finalizer wait]:
+ (go) runtime.gopark(0x0?, 0x0?, 0x0?, 0x0?, 0x0?)
+ (go) 	runtime/proc.go:474 +0xb8 fp=0x7fa951124580 sp=0x7fa951124560 pc=0x91218
+ (go) runtime.runFinalizers()
+ (go) 	runtime/mfinal.go:210 +0x100 fp=0x7fa9511247d0 sp=0x7fa951124580 pc=0x37cf0
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa9511247d0 sp=0x7fa9511247d0 pc=0x97f94
+ (go) created by runtime.createfing in goroutine 1
+ (go) 	runtime/mfinal.go:172 +0x78
- (oracle) validate [--no-check-all] [--check-lock] [--no-check-lock] [--no-check-publish] [--no-check-version] [-A|--with-dependencies] [--strict] [--] [<file>]
+ (go) goroutine 18 gp=0x7fa951184960 m=nil [cleanup wait]:
+ (go) runtime.gopark(0x0?, 0x0?, 0x0?, 0x0?, 0x0?)
+ (go) 	runtime/proc.go:474 +0xb8 fp=0x7fa951124f50 sp=0x7fa951124f30 pc=0x91218
+ (go) runtime.goparkunlock(...)
+ (go) 	runtime/proc.go:480
+ (go) runtime.(*cleanupQueue).dequeue(0x87ce60)
+ (go) 	runtime/mcleanup.go:522 +0x160 fp=0x7fa951124f90 sp=0x7fa951124f50 pc=0x34850
+ (go) runtime.runCleanups()
+ (go) 	runtime/mcleanup.go:718 +0x40 fp=0x7fa951124fd0 sp=0x7fa951124f90 pc=0x34fa0
+ (go) runtime.goexit({})
+ (go) 	runtime/asm_arm64.s:1039 +0x4 fp=0x7fa951124fd0 sp=0x7fa951124fd0 pc=0x97f94
+ (go) created by runtime.(*cleanupQueue).createGs in goroutine 1
+ (go) 	runtime/mcleanup.go:672 +0x108
```
