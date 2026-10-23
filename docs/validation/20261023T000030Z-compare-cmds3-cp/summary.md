# Composer vs gocomposer comparison

- gocomposer_sha256: `4dc277e7cd04d5f7151109197d10baa0044bc86d344c172ea397cf236584b3c3`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| create-project-basic | DIFF | 0 | 0 |
| create-project-default-dir | DIFF | 0 | 0 |
| create-project-version-arg | DIFF | 0 | 0 |
| create-project-version-colon | DIFF | 0 | 0 |
| create-project-version-equals | DIFF | 0 | 0 |
| create-project-stability-flag | DIFF | 0 | 0 |
| create-project-stability-at | DIFF | 0 | 0 |
| create-project-version-beta | DIFF | 0 | 0 |
| create-project-stability-invalid | MATCH | 1 | 1 |
| create-project-not-found | MATCH | 1 | 1 |
| create-project-not-found-version | MATCH | 1 | 1 |
| create-project-php-too-low | MATCH | 1 | 1 |
| create-project-php-ignore | MATCH | 0 | 0 |
| create-project-secure-http | MATCH | 100 | 100 |
| create-project-repository-url | DIFF | 0 | 0 |
| create-project-packagist-false-repo | DIFF | 0 | 0 |
| create-project-no-install | DIFF | 0 | 0 |
| create-project-no-scripts | DIFF | 0 | 0 |
| create-project-no-dev | DIFF | 0 | 0 |
| create-project-dev-option | DIFF | 0 | 0 |
| create-project-no-custom-installers | DIFF | 0 | 0 |
| create-project-no-plugins | DIFF | 0 | 0 |
| create-project-no-audit | DIFF | 0 | 0 |
| create-project-add-repository | DIFF | 0 | 0 |
| create-project-require | DIFF | 0 | 0 |
| create-project-ask-noninteractive | DIFF | 0 | 0 |
| create-project-prefer-conflict | MATCH | 1 | 1 |
| create-project-prefer-source-dist-only | DIFF | 0 | 0 |
| create-project-verbose | DIFF | 0 | 0 |
| create-project-progress | DIFF | 0 | 0 |
| create-project-relative-dir | DIFF | 0 | 0 |
| create-project-absolute-dir | DIFF | 0 | 0 |
| create-project-plain | MATCH | 0 | 0 |
| create-project-bin | MATCH | 0 | 0 |
| create-project-packagist | DIFF | 100 | 100 |
| create-project-not-empty | MATCH | 1 | 1 |
| create-project-dir-is-file | MATCH | 1 | 1 |
| create-project-path | MATCH | 0 | 0 |
| create-project-path-symlink | MATCH | 0 | 0 |
| create-project-vcs | MATCH | 0 | 0 |
| create-project-vcs-keep | MATCH | 0 | 0 |
| create-project-vcs-remove | MATCH | 0 | 0 |
| create-project-vcs-dev | MATCH | 0 | 0 |
| create-project-vcs-branch | MATCH | 0 | 0 |
| create-project-no-args | MATCH | 0 | 0 |
| create-project-no-args-no-json | MATCH | 1 | 1 |

## create-project-basic

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-default-dir

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-version-arg

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 9cb3705b557cca23
```

## create-project-version-colon

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-version-equals

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 9cb3705b557cca23
```

## create-project-stability-flag

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 1e420dab392eeede
```

## create-project-stability-at

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 1e420dab392eeede
```

## create-project-version-beta

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 1e420dab392eeede
```

## create-project-repository-url

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-packagist-false-repo

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-install

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-scripts

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-dev

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-dev-option

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-custom-installers

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-plugins

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-no-audit

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-add-repository

```
project tree:
+ (go) F app/composer.json 644 4ec57af654bcb71b
- (oracle) F app/composer.json 644 aeb7df4647f32d98
```

## create-project-require

```
project tree:
+ (go) F app/composer.json 644 786305da7d989814
- (oracle) F app/composer.json 644 44328f36cf8f3d15
```

## create-project-ask-noninteractive

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-prefer-source-dist-only

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-verbose

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-progress

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
```

## create-project-relative-dir

```
project tree:
+ (go) F sub/dir/app/composer.json 644 d5dd35c84a9b491f
- (oracle) F sub/dir/app/composer.json 644 bd1767e5801599ed
```

## create-project-absolute-dir

```
project tree:
+ (go) F abs/composer.json 644 d5dd35c84a9b491f
- (oracle) F abs/composer.json 644 bd1767e5801599ed
```

## create-project-packagist

```
stderr:
- (oracle) The following exception probably indicates you are offline or have misconfigured DNS resolver(s)
- (oracle) If you intend to run Composer without connecting to the internet, run the command again prefixed with COMPOSER_DISABLE_NETWORK=1 to make Composer run in offline mode.
```
