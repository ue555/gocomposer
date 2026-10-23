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
| self-update-basic | KNOWN-DIFF | 1 | 1 |
| self-update-alias | KNOWN-DIFF | 1 | 1 |
| self-update-rollback | KNOWN-DIFF | 1 | 1 |
| self-update-version | KNOWN-DIFF | 1 | 1 |
| self-update-no-json | KNOWN-DIFF | 1 | 1 |
| self-update-bad-option | MATCH | 1 | 1 |
| self-update-too-many | MATCH | 1 | 1 |

## self-update-basic

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-alias

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-rollback

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-version

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```

## self-update-no-json

```
stdout:
+ (go) gocomposer does not implement self-update as it is not distributed as a composer.phar. Rebuild or reinstall gocomposer to update it.
- (oracle) This could be due to a number of reasons, such as Composer being installed as a system package on your OS, or Composer being installed as a package in the current project.
```
