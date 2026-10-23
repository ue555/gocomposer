# Composer vs gocomposer comparison

- gocomposer_sha256: `9512f8d65cd3640152466542b9854ec146a8c53265384f3bbfcdf99805adeb95`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| reinstall-basic | MATCH | 0 | 0 |
| reinstall-progress | DIFF | 0 | 0 |
| reinstall-verbose | DIFF | 0 | 0 |
| reinstall-wildcard | DIFF | 0 | 0 |
| reinstall-all-wildcard | DIFF | 0 | 0 |
| reinstall-type | DIFF | 0 | 0 |
| reinstall-type-none | MATCH | 1 | 1 |
| reinstall-type-and-names | MATCH | 1 | 1 |
| reinstall-no-args | MATCH | 1 | 1 |
| reinstall-no-match | MATCH | 1 | 1 |
| reinstall-partial-match | MATCH | 0 | 0 |
| reinstall-not-installed | MATCH | 1 | 1 |
| reinstall-lock-only | MATCH | 1 | 1 |
| reinstall-no-autoloader | MATCH | 0 | 0 |
| reinstall-optimize | MATCH | 0 | 0 |
| reinstall-restores-files | MATCH | 0 | 0 |
| reinstall-no-dev-mode | MATCH | 0 | 0 |
| reinstall-scripts | MATCH | 0 | 0 |
| reinstall-no-scripts | MATCH | 0 | 0 |
| reinstall-no-json | MATCH | 1 | 1 |
| reinstall-prefer-source | MATCH | 0 | 0 |
| reinstall-prefer-install-source | MATCH | 0 | 0 |
| reinstall-source-to-dist | MATCH | 0 | 0 |
| reinstall-source-default | MATCH | 0 | 0 |
| reinstall-prefer-conflict | MATCH | 1 | 1 |
| reinstall-git | MATCH | 0 | 0 |
| reinstall-git-changes | MATCH | 1 | 1 |
| reinstall-git-discard | MATCH | 0 | 0 |

## reinstall-progress

```
stderr:
+ (go)  0/2 [>---------------------------]   0%
+ (go)  2/2 [============================] 100%
- (oracle)     0 [>---------------------------]    0 [->--------------------------]
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go)  0/2 [>---------------------------]   0%Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) 
- (oracle)  0/2 [>---------------------------]   0%
```

## reinstall-verbose

```
stderr:
- (oracle) Running 2.11.999-dev+source (@release_date@) with PHP 8.4.27 on Linux / 7.1.10-200.fc44.aarch64
- (oracle) Reading ./composer.json (/work/project/composer.json)
- (oracle) Loading config file ./composer.json (/work/project/composer.json)
- (oracle) Checked CA file /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem does not exist or it is not a directory.
- (oracle) Checked CA file /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a file.
- (oracle) Checked directory /etc/pki/tls/certs/ca-bundle.crt does not exist or it is not a directory.
- (oracle) Checked CA file /etc/ssl/certs/ca-certificates.crt: valid
- (oracle) Executing command (CWD): 'git' '--version'
- (oracle) Executing command (/work/project): 'git' 'branch' '-a' '--no-color' '--no-abbrev' '-v'
- (oracle) Executing command (/work/project): 'git' 'describe' '--exact-match' '--tags'
- (oracle) Executing command (/work/project): 'git' 'rev-list' '--no-commit-header' '--format=%H' '-n1' 'HEAD' '--no-show-signature'
- (oracle) Executing command (/work/project): 'hg' 'branch'
- (oracle) Executing command (/work/project): 'fossil' 'branch' 'list'
- (oracle) Executing command (/work/project): 'fossil' 'tag' 'list'
- (oracle) Executing command (/work/project): 'svn' 'info' '--xml'
- (oracle) Failed to initialize global composer: Composer could not find the config file: /work/home/composer.json
- (oracle) 
- (oracle) Reading ./composer.lock (/work/project/composer.lock)
- (oracle) Reading /work/project/vendor/composer/installed.json
- (oracle) Executing async command (CWD): 'rm' '-rf' '/work/project/vendor/acme/log'
- (oracle) Downloading http://127.0.0.1:8899/dists/acme-log-1.1.0.zip
- (oracle) [200] http://127.0.0.1:8899/dists/acme-log-1.1.0.zip
- (oracle) Writing /work/cache/files/acme/log/31e5146f9094134bf78a7c5baeb8e0a5b2e0e354.zip into cache from /work/project/vendor/composer/tmp-HASH.zip
- (oracle) Executing async command (CWD): '/usr/bin/unzip' '-qq' '/work/project/vendor/composer/tmp-HASH.zip' '-d' '/work/project/vendor/composer/364670b0'
- (oracle) Executing async command (CWD): 'rm' '-rf' '/work/project/vendor/composer/364670b0'
```

## reinstall-wildcard

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## reinstall-all-wildcard

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## reinstall-type

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```
