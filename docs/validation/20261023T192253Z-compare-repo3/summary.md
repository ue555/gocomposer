# Composer vs gocomposer comparison

- gocomposer_sha256: `7b27b923b1cf71afbd6ac15e193754510d251a287368b0c8fd02f7d248eab39d`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:cd7f1f040bba14a1974d175df2110347c2be1d880834ad6a106a3578bfd929fa`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| archive3-xz | MATCH | 0 | 0 |
| archive3-xz-broken | DIFF | 1 | 1 |
| archive3-rar-no-unrar | MATCH | 1 | 1 |
| archive3-rar | MATCH | 0 | 0 |
| archive3-rar-broken | MATCH | 1 | 1 |
| archive3-zip-7z | MATCH | 0 | 0 |
| archive3-zip-7z-zipext | MATCH | 0 | 0 |
| archive3-zip-7z-crc | DIFF | 1 | 1 |
| archive3-zip-7z-crc-zipext | DIFF | 1 | 1 |
| archive3-zip-no-command | MATCH | 1 | 1 |
| archive3-zip-no-command-zipext | MATCH | 0 | 0 |
| archive3-zip-no-command-crc-zipext | DIFF | 1 | 1 |
| warnings-update | MATCH | 0 | 0 |
| warnings-update-ansi | DIFF | 0 | 0 |
| warnings-show | MATCH | 0 | 0 |
| warnings-search-api | MATCH | 0 | 0 |
| warnings-search-list | MATCH | 0 | 0 |
| warnings-show-names-list | MATCH | 0 | 0 |
| warnings-audit-api | ERROR | 0 | 0 |
| warnings-update-audit-api | MATCH | 0 | 0 |
| warnings-filter-audit | MATCH | 0 | 0 |
| warnings-providers | MATCH | 0 | 0 |
| warnings-two-repos | MATCH | 0 | 0 |
| warnings-http-403 | MATCH | 100 | 100 |
| warnings-http-404-metadata | MATCH | 2 | 2 |

## archive3-xz-broken

```
stderr:
+ (go)   Failed to execute tar -xJf /work/project/vendor/composer/tmp-f78d3e2f39e87f  
+ (go)   a03f1b3f1d9ed9f391.xz -C /work/project/vendor/composer/445cd40e              
- (oracle)   Failed to execute tar -xJf /work/project/vendor/composer/tmp-4769c803730cc1  
- (oracle)   612668d52031d8b8cf.xz -C /work/project/vendor/composer/eb44ab2f              
```

## archive3-zip-7z-crc

```
stderr:
+ (go)   t/vendor/composer/f2c5e0d1                                                   
- (oracle)   t/vendor/composer/fb2259b0                                                   
```

## archive3-zip-7z-crc-zipext

```
stderr:
+ (go)     Failed to extract archfix/zipcrc: (2) /opt/arch/bin/7z x -bb0 -y /work/project/vendor/composer/tmp-HASH.zip -o/work/project/vendor/composer/68c94c62
- (oracle)     Failed to extract archfix/zipcrc: (2) /opt/arch/bin/7z x -bb0 -y /work/project/vendor/composer/tmp-HASH.zip -o/work/project/vendor/composer/479c6a6e
+ (go)   There was an error extracting the ZIP file for "archfix/zipcrc", it is eith  
+ (go)   er corrupted or using an invalid format.                                     
- (oracle)   The archive for "archfix/zipcrc" may contain identical file names with diff  
- (oracle)   erent capitalization (which fails on case insensitive filesystems): ZipArch  
- (oracle)   ive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json": "CRC  
- (oracle)    error"                                                                      
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   ZipArchive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json  
- (oracle)   ": "CRC error"                                                               
- (oracle)                                                                                
- (oracle) 
```

## archive3-zip-no-command-crc-zipext

```
stderr:
+ (go)   There was an error extracting the ZIP file for "archfix/zipcrc", it is eith  
+ (go)   er corrupted or using an invalid format.                                     
- (oracle)   The archive for "archfix/zipcrc" may contain identical file names with diff  
- (oracle)   erent capitalization (which fails on case insensitive filesystems): ZipArch  
- (oracle)   ive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json": "CRC  
- (oracle)    error"                                                                      
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   ZipArchive::extractTo(): Cannot extract "archfix-zipcrc-1.0.0/composer.json  
- (oracle)   ": "CRC error"                                                               
- (oracle)                                                                                
- (oracle) 
```

## warnings-update-ansi

```
stderr:
+ (go) [30;43mWarning from http://127.0.0.1:8899/git/repo3/warn: Root [31mwarning[0m text[39;49m
- (oracle) [30;43mWarning from http://127.0.0.1:8899/git/repo3/warn: Root warning text[39;49m
```

## warnings-audit-api

```
setup failed: update --no-progress --no-audit
Loading composer repositories with package information
Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
Warning from http://127.0.0.1:8899/git/repo3/warnadv: Advisory API warning
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires provfix/lib 1.0.0 (exact version match: 1.0.0 or 1.0.0.0), found provfix/lib[1.0.0] but these were not loaded, because they are affected by security advisories ("PKSA-prov-0001"). Go to https://packagist.org/security-advisories/ to find advisory details. To ignore the advisories, add their IDs to the "policy.advisories.ignore-id" config or add the package to "policy.advisories.ignore". To turn the feature off entirely, you can set "policy.advisories.block" to false.


```
