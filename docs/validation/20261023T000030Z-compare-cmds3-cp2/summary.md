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
| create-project-packagist | DIFF | 100 | 100 |

## create-project-basic

```
project tree:
+ (go) F app/composer.json 644 d5dd35c84a9b491f
- (oracle) F app/composer.json 644 bd1767e5801599ed
app/composer.json:
+ (go)     "require": {"acme/log": "^1.0", "cproj/lib": "self.version"},
- (oracle)     "require": {"acme/log": "^1.0", "cproj/lib": "1.1.0"},
```

## create-project-packagist

```
stderr:
- (oracle) The following exception probably indicates you are offline or have misconfigured DNS resolver(s)
- (oracle) If you intend to run Composer without connecting to the internet, run the command again prefixed with COMPOSER_DISABLE_NETWORK=1 to make Composer run in offline mode.
```
