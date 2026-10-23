# Composer vs gocomposer comparison

- gocomposer_sha256: `e6a0d9e271a6c37be9ff2fc01112d5a4100a720c732e083e4299efe54a079240`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| bump-basic | ERROR | 0 | 0 |
| bump-dry-run | ERROR | 0 | 0 |
| bump-dev-only | ERROR | 0 | 0 |
| bump-no-dev-only | ERROR | 0 | 0 |
| bump-library | ERROR | 0 | 0 |
| bump-library-dev-only | ERROR | 0 | 0 |
| bump-untyped | ERROR | 0 | 0 |
| bump-filter | ERROR | 0 | 0 |
| bump-filter-wildcard | ERROR | 0 | 0 |
| bump-nothing | ERROR | 0 | 0 |
| bump-nothing-dry-run | ERROR | 0 | 0 |
| bump-lock-outdated | ERROR | 0 | 0 |
| bump-lock-disabled | ERROR | 0 | 0 |
| bump-installed-only | ERROR | 0 | 0 |
| bump-not-installed | MATCH | 0 | 0 |
| bump-branch-alias | MATCH | 0 | 0 |
| bump-no-json | MATCH | 1 | 1 |
| bump-after-update-flag | MATCH | 2 | 2 |
| bump-after-update-dev | MATCH | 2 | 2 |
| bump-after-update-no-dev | MATCH | 2 | 2 |
| bump-after-update-dry-run | MATCH | 2 | 2 |
| bump-after-update-config | MATCH | 2 | 2 |
| bump-after-update-config-dev | MATCH | 2 | 2 |
| bump-after-update-untyped | MATCH | 2 | 2 |
| bump-after-update-partial | ERROR | 0 | 0 |
| bump-after-update-nothing | ERROR | 0 | 0 |
| bump-after-update-lock | ERROR | 0 | 0 |

## bump-basic

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-dry-run

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-dev-only

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-no-dev-only

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-library

```
setup failed: update --no-progress
Composer could not detect the root package (test/project) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-library-dev-only

```
setup failed: update --no-progress
Composer could not detect the root package (test/project) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-untyped

```
setup failed: update --no-progress
Composer could not detect the root package (test/project) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-filter

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-filter-wildcard

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-nothing

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-nothing-dry-run

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-lock-outdated

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-lock-disabled

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-installed-only

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-after-update-partial

```
setup failed: update --no-progress --with acme/log:1.0.0 --with acme/http:1.0.0
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0].
    - acme/http 1.0.0 requires psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-after-update-nothing

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```

## bump-after-update-lock

```
setup failed: update --no-progress
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires acme/http ~1.0 -> satisfiable by acme/http[1.0.0, 1.2.0].
    - acme/http[1.0.0, ..., 1.2.0] require psr/log-implementation ^1.0 -> could not be found in any version, but the following packages provide it:
      - acme/log-impl Provides psr/log-implementation
      Consider requiring one of these to satisfy the psr/log-implementation requirement.

Potential causes:
 - A typo in the package name
 - The package is not available in a stable-enough version according to your minimum-stability setting
   see <https://getcomposer.org/doc/04-schema.md#minimum-stability-root-only> for more details.
 - It's a private package and you forgot to add a custom repository to find it

Read <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.

```
