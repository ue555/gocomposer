# Composer vs gocomposer comparison

- gocomposer_sha256: `28088f3378cc3fd74c352ef2cb08451945227b678e7cf03fb55824a1a7387ac7`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| update-basic | MATCH | 0 | 0 |
| update-basic-progress | MATCH | 0 | 0 |
| update-v2-metadata | MATCH | 0 | 0 |
| install-from-lock | DIFF | 0 | 0 |
| install-no-dev | DIFF | 0 | 0 |
| install-twice | DIFF | 0 | 0 |
| dry-run-update | MATCH | 0 | 0 |
| dump-autoload-optimized | MATCH | 0 | 0 |
| dump-autoload-authoritative | MATCH | 0 | 0 |
| replace | MATCH | 0 | 0 |
| conflict-unresolvable | MATCH | 2 | 2 |
| conflict-resolves-lower | MATCH | 0 | 0 |
| platform-php-too-low | MATCH | 2 | 2 |
| missing-extension | MATCH | 2 | 2 |
| ignore-platform-reqs | MATCH | 0 | 0 |
| config-platform-php | MATCH | 0 | 0 |
| not-found-package | MATCH | 2 | 2 |
| stability-beta-flag | MATCH | 0 | 0 |
| stability-too-low | MATCH | 2 | 2 |
| dev-branch-alias | MATCH | 0 | 0 |
| root-alias | MATCH | 0 | 0 |
| partial-update | MATCH | 0 | 0 |
| partial-update-with-deps | MATCH | 0 | 0 |
| broken-archive | MATCH | 100 | 100 |
| checksum-mismatch | MATCH | 1 | 1 |
| download-404 | MATCH | 100 | 100 |
| path-traversal-archive | MATCH | 0 | 0 |
| failure-keeps-vendor | MATCH | 2 | 2 |
| no-lock-install | MATCH | 0 | 0 |
| check-platform-reqs-ok | MATCH | 0 | 0 |
| check-platform-reqs-fail | MATCH | 2 | 2 |
| check-platform-reqs-json | MATCH | 2 | 2 |
| check-platform-reqs-lock | MATCH | 2 | 2 |
| validate-basic | MATCH | 0 | 0 |
| validate-clean | MATCH | 0 | 0 |
| validate-warnings | DIFF | 1 | 0 |
| validate-strict-warnings | MATCH | 1 | 1 |
| validate-no-check | MATCH | 0 | 0 |
| validate-missing-file | MATCH | 3 | 3 |
| validate-invalid-json | MATCH | 1 | 1 |
| validate-duplicate-key | MATCH | 0 | 0 |
| validate-lock-stale | MATCH | 2 | 2 |
| validate-lock-stale-no-check | MATCH | 0 | 0 |
| validate-with-dependencies | MATCH | 0 | 0 |
| validate-bad-constraint | DIFF | 1 | 1 |

## install-from-lock

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## install-no-dev

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## install-twice

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## validate-warnings

```
exit code: oracle=1 go=0
stderr:
+ (go) ./composer.json is valid, but with a few warnings
+ (go) See https://getcomposer.org/doc/04-schema.md for details on the schema
+ (go) # General warnings
+ (go) - License "GPL-2.0" is a deprecated SPDX license identifier, use "GPL-2.0-only" or "GPL-2.0-or-later" instead
+ (go) - The version field is present, it is recommended to leave it out if the package is published on Packagist.
+ (go) - acme/cli is required both in require and require-dev, this can lead to unexpected behavior
+ (go) - The package acme/cli in require is also listed in provide which satisfies the requirement. Remove it from provide if you wish to install it.
+ (go) - The package acme/cli in require-dev is also listed in provide which satisfies the requirement. Remove it from provide if you wish to install it.
+ (go) - The package "acme/http" is pointing to a commit-ref, this is bad practice and can cause unforeseen issues.
+ (go) - Description for non-existent script "b" found in "scripts-descriptions"
+ (go) - homepage : invalid value (example.com), must be an http/https URL
+ (go) - keywords.1 : invalid value (bad!), must match [\p{N}\p{L} ._-]+
+ (go) - License "Foo-License" is not a valid SPDX license identifier, see https://spdx.org/licenses/ if you use an open license.
+ (go) If the software is closed-source, you may use "proprietary" as license.
+ (go) - License " MIT " must not contain extra spaces, make sure to trim it.
+ (go) - authors.0.email : invalid value (not-an-email), must be a valid email address
+ (go) - require.acme/log : unbound version constraints (*) should be avoided
+ (go) - require.acme/util : exact version constraints (1.0.0) should be avoided if the package follows semantic versioning
- (oracle) In <location>:
- (oracle)                                                               
- (oracle)   "./composer.json" does not match the expected JSON schema:  
- (oracle)    - authors[0].email : Invalid email                         
- (oracle)    - homepage : Invalid URL format                            
- (oracle)                                                               
- (oracle) 
- (oracle) validate [--no-check-all] [--check-lock] [--no-check-lock] [--no-check-publish] [--no-check-version] [-A|--with-dependencies] [--strict] [--] [<file>]
- (oracle) 
- (oracle) 
```

## validate-bad-constraint

```
stderr:
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   Could not parse version constraint ~>1: Invalid operator "~>", you probably  
- (oracle)    meant to use the "~" operator                                               
- (oracle)                                                                                
- (oracle) 
```
