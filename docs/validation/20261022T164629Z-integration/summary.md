# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `c455b969f4f4c61d665f71c968418554a6397d8ecb9065ce28b79f8e274e7e0b`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| project | step | status | oracle exit | go exit |
|---|---|---|---|---|
| cli-app | setup | ERROR | 0 | 0 |

## cli-app / setup

```
official composer update failed:
     Loading composer repositories with package information
     Updating dependencies
     Lock file operations: 27 installs, 0 updates, 0 removals
       - Locking carbonphp/carbon-doctrine-types (3.2.1)
       - Locking guzzlehttp/guzzle (7.16.0)
       - Locking guzzlehttp/promises (2.6.0)
       - Locking guzzlehttp/psr7 (2.13.2)
       - Locking monolog/monolog (3.12.1)
       - Locking nesbot/carbon (3.15.0)
       - Locking psr/clock (1.0.0)
       - Locking psr/container (2.0.2)
       - Locking psr/http-client (1.0.3)
       - Locking psr/http-factory (1.1.0)
       - Locking psr/http-message (2.0)
       - Locking psr/log (3.0.2)
       - Locking ralouphie/getallheaders (3.0.3)
       - Locking symfony/clock (v8.1.9)
       - Locking symfony/console (v7.4.21)
       - Locking symfony/deprecation-contracts (v3.7.1)
       - Locking symfony/polyfill-ctype (v1.37.0)
       - Locking symfony/polyfill-intl-grapheme (v1.43.0)
       - Locking symfony/polyfill-intl-normalizer (v1.43.0)
       - Locking symfony/polyfill-mbstring (v1.43.0)
       - Locking symfony/polyfill-php80 (v1.43.0)
       - Locking symfony/polyfill-time (v1.42.0)
       - Locking symfony/service-contracts (v3.7.3)
       - Locking symfony/string (v8.1.9)
       - Locking symfony/translation (v8.1.5)
       - Locking symfony/translation-contracts (v3.7.1)
       - Locking symfony/var-dumper (v7.4.21)
     Writing lock file
     15 package suggestions were added by new dependencies, use `composer suggest` to see details.
     > @php bin/app greet --name=update
     Script @php bin/app greet --name=update handling the post-update-cmd event returned with error code 255
     
```
