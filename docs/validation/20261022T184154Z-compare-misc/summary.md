# Composer vs gocomposer comparison

- gocomposer_sha256: `5ab34596da54fe4337559b4d9623e1b6627eaeb544111cdbeab8a81ae13aa943`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| policy-cmd-add-source | MATCH | 0 | 0 |
| policy-cmd-add-source-json | MATCH | 0 | 0 |
| policy-cmd-add-source-existing-list | MATCH | 0 | 0 |
| policy-cmd-add-source-duplicate | MATCH | 0 | 0 |
| policy-cmd-add-source-object-sources | DIFF | 1 | 1 |
| policy-cmd-add-source-uppercase-action | MATCH | 0 | 0 |
| policy-cmd-no-action | DIFF | 1 | 1 |
| policy-cmd-unknown-action | MATCH | 1 | 1 |
| policy-cmd-no-name | MATCH | 1 | 1 |
| policy-cmd-builtin | MATCH | 1 | 1 |
| policy-cmd-builtin-cooldown | MATCH | 1 | 1 |
| policy-cmd-reserved-prefix | MATCH | 1 | 1 |
| policy-cmd-reserved-name | MATCH | 1 | 1 |
| policy-cmd-dotted-name | MATCH | 1 | 1 |
| policy-cmd-empty-name | MATCH | 1 | 1 |
| policy-cmd-no-arg1 | MATCH | 1 | 1 |
| policy-cmd-no-url | MATCH | 1 | 1 |
| policy-cmd-bad-type | MATCH | 1 | 1 |
| policy-cmd-http-url | MATCH | 1 | 1 |
| policy-cmd-json-no-type | MATCH | 1 | 1 |
| policy-cmd-json-no-url | MATCH | 1 | 1 |
| policy-cmd-json-int-url | MATCH | 1 | 1 |
| policy-cmd-json-invalid | MATCH | 1 | 1 |
| policy-cmd-global | MATCH | 0 | 0 |
| policy-cmd-file-and-global | MATCH | 1 | 1 |
| policy-cmd-file | MATCH | 0 | 0 |
| policy-cmd-no-json | MATCH | 1 | 1 |

## policy-cmd-add-source-object-sources

```
stderr:
- (oracle) In <location>:
- (oracle)                                                              
- (oracle)   "./composer.json" does not match the expected JSON schema  
- (oracle)                                                              
- (oracle) 
```

## policy-cmd-no-action

```
stderr:
- (oracle) Composer could not detect the root package (test/plain) version, defaulting to '1.0.0'. See https://getcomposer.org/root-version
```
