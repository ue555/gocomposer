# Composer vs gocomposer comparison

- gocomposer_sha256: `f7ee911ed591d230174992213d6bcfea0de22767693817e00d68025daa153ce7`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| show-list | MATCH | 0 | 0 |
| show-json | MATCH | 0 | 0 |
| show-name-only | MATCH | 0 | 0 |
| show-direct | MATCH | 0 | 0 |
| show-no-dev | MATCH | 0 | 0 |
| show-path | MATCH | 0 | 0 |
| show-package | MATCH | 0 | 0 |
| show-package-json | MATCH | 0 | 0 |
| show-tree | DIFF | 0 | 0 |
| show-tree-json | MATCH | 0 | 0 |
| show-package-tree | DIFF | 0 | 0 |
| show-locked | MATCH | 0 | 0 |
| show-self | DIFF | 0 | 0 |
| show-self-json | DIFF | 0 | 0 |
| show-platform-php | MATCH | 0 | 0 |
| show-wildcard | MATCH | 0 | 0 |
| show-not-found | MATCH | 1 | 1 |
| show-version-arg | MATCH | 0 | 0 |
| show-self-name-only | MATCH | 0 | 0 |
| show-locked-package | MATCH | 0 | 0 |
| show-platform-list-json | MATCH | 0 | 0 |
| show-bad-format | MATCH | 1 | 1 |
| show-tree-latest | MATCH | 1 | 1 |
| show-nothing-installed | MATCH | 0 | 0 |

## show-tree

```
stdout:
+ (go) `--<yellow>acme/util</yellow> ^1.0
+ (go)    |--<cyan>acme/log</cyan> ^1.1
+ (go)    |  `--<magenta>php</magenta> >=7.4
+ (go)    `--<cyan>ext-json</cyan> *
- (oracle) `--acme/util ^1.0
- (oracle)    |--acme/log ^1.1
- (oracle)    |  `--php >=7.4
- (oracle)    `--ext-json *
+ (go) `--<yellow>acme/log</yellow> ^1.0 || ^2.0
+ (go)    `--<cyan>php</cyan> >=7.4
- (oracle) `--acme/log ^1.0 || ^2.0
- (oracle)    `--php >=7.4
+ (go) |--<yellow>acme/util</yellow> ^1.0
+ (go) |  |--<cyan>acme/log</cyan> ^1.1
+ (go) |  |  `--<magenta>php</magenta> >=7.4
+ (go) |  `--<cyan>ext-json</cyan> *
+ (go) |--<yellow>php</yellow> ^7.4 || ^8.0
+ (go) `--<yellow>psr/log-implementation</yellow> ^1.0
- (oracle) |--acme/util ^1.0
- (oracle) |  |--acme/log ^1.1
- (oracle) |  |  `--php >=7.4
- (oracle) |  `--ext-json *
- (oracle) |--php ^7.4 || ^8.0
- (oracle) `--psr/log-implementation ^1.0
+ (go) `--<yellow>acme/log</yellow> ^1.0
+ (go)    `--<cyan>php</cyan> >=7.4
- (oracle) `--acme/log ^1.0
- (oracle)    `--php >=7.4
```

## show-package-tree

```
stdout:
+ (go) |--<yellow>acme/util</yellow> ^1.0
+ (go) |  |--<cyan>acme/log</cyan> ^1.1
+ (go) |  |  `--<magenta>php</magenta> >=7.4
+ (go) |  `--<cyan>ext-json</cyan> *
+ (go) |--<yellow>php</yellow> ^7.4 || ^8.0
+ (go) `--<yellow>psr/log-implementation</yellow> ^1.0
- (oracle) |--acme/util ^1.0
- (oracle) |  |--acme/log ^1.1
- (oracle) |  |  `--php >=7.4
- (oracle) |  `--ext-json *
- (oracle) |--php ^7.4 || ^8.0
- (oracle) `--psr/log-implementation ^1.0
```

## show-self

```
stdout:
+ (go) path     : /work/project/vendor/test/project
- (oracle) path     : 
```

## show-self-json

```
stdout:
+ (go)     "path": "/work/project/vendor/test/project",
```
