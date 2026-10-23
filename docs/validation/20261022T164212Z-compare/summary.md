# Composer vs gocomposer comparison

- gocomposer_sha256: `fd048fa49815cb3c1c5d8b9207c064eb118eea9b2acbd7705ab159533ba8697b`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| policy-advisory-blocks | MATCH | 2 | 2 |
| policy-advisory-blocks-vv | DIFF | 2 | 2 |
| policy-advisory-only-vulnerable | MATCH | 2 | 2 |
| policy-no-blocking-audit | MATCH | 2 | 2 |
| policy-env-no-blocking | MATCH | 2 | 2 |
| policy-audit-install | ERROR | 0 | 0 |
| policy-audit-install-table | ERROR | 0 | 0 |
| policy-audit-install-json | ERROR | 0 | 0 |
| policy-audit-plain | ERROR | 0 | 0 |
| policy-abandoned | MATCH | 0 | 0 |
| policy-abandoned-ignore | MATCH | 0 | 0 |
| policy-ignore-advisory | MATCH | 2 | 2 |
| policy-disabled | MATCH | 2 | 2 |
| policy-cooldown-all | MATCH | 1 | 1 |
| policy-cooldown-none | MATCH | 1 | 1 |
| policy-cooldown-no-cooldown | MATCH | 0 | 0 |
| policy-locked-vulnerable-install | ERROR | 0 | 0 |
| policy-partial-update-vulnerable | ERROR | 0 | 0 |

## policy-advisory-blocks-vv

```
stdout:
- (oracle) Pool optimizer completed in N seconds
- (oracle) Found 61 package versions referenced in your dependency graph. 1 (2%) were optimized away.
stderr:
- (oracle) Dependency resolution completed in N seconds
```

## policy-audit-install

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.2 -> satisfiable by acme/http[1.2.0].
    - acme/http 1.2.0 requires psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## policy-audit-install-table

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.2 -> satisfiable by acme/http[1.2.0].
    - acme/http 1.2.0 requires psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## policy-audit-install-json

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.2 -> satisfiable by acme/http[1.2.0].
    - acme/http 1.2.0 requires psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## policy-audit-plain

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.2 -> satisfiable by acme/http[1.2.0].
    - acme/http 1.2.0 requires psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## policy-locked-vulnerable-install

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.2 -> satisfiable by acme/http[1.2.0].
    - acme/http 1.2.0 requires psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## policy-partial-update-vulnerable

```
setup failed: update --no-progress --no-blocking --no-install
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ^1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, there may be a typo in the package name.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```
