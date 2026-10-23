# Composer vs gocomposer comparison

- gocomposer_sha256: `e209883dbed422ef8502abc210d92cebb934d43baf52a4383f6928c4a5bed300`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| show-latest-short | MATCH | 0 | 0 |
| show-latest-outdated | MATCH | 0 | 0 |
| show-latest-outdated-direct | MATCH | 0 | 0 |
| show-latest-direct | MATCH | 0 | 0 |
| show-latest-strict | MATCH | 1 | 1 |
| show-latest-strict-latest-only | MATCH | 1 | 1 |
| show-latest-strict-ignored | MATCH | 0 | 0 |
| show-latest-strict-fresh | MATCH | 0 | 0 |
| show-latest-fresh | MATCH | 0 | 0 |
| show-latest-strict-no-latest | MATCH | 1 | 1 |
| show-latest-major-only | MATCH | 0 | 0 |
| show-latest-major-only-latest | MATCH | 0 | 0 |
| show-latest-minor-only | MATCH | 0 | 0 |
| show-latest-patch-only | MATCH | 0 | 0 |
| show-latest-patch-only-outdated | MATCH | 0 | 0 |
| show-latest-only-conflict | MATCH | 1 | 1 |
| show-latest-json | MATCH | 0 | 0 |
| show-latest-outdated-json | MATCH | 0 | 0 |
| show-latest-outdated-json-strict | MATCH | 1 | 1 |
| show-latest-name-only | MATCH | 0 | 0 |
| show-latest-name-only-json | MATCH | 0 | 0 |
| show-latest-path | MATCH | 0 | 0 |
| show-latest-wildcard | MATCH | 0 | 0 |
| show-latest-ignore | MATCH | 0 | 0 |
| show-latest-ignore-without-outdated | MATCH | 0 | 0 |
| show-latest-no-dev | MATCH | 0 | 0 |
| show-latest-sort-by-age | MATCH | 0 | 0 |
| show-latest-ansi | DIFF | 0 | 0 |
| show-latest-outdated-ansi | DIFF | 0 | 0 |
| show-latest-verbose | MATCH | 0 | 0 |
| show-latest-ignore-platform-reqs | MATCH | 0 | 0 |
| show-latest-ignore-platform-req | MATCH | 0 | 0 |
| show-latest-platform-rejected | MATCH | 0 | 0 |
| show-latest-platform-rejected-verbose | MATCH | 0 | 0 |
| show-latest-locked | MATCH | 0 | 0 |
| show-latest-locked-outdated | MATCH | 0 | 0 |
| show-latest-not-installed | MATCH | 0 | 0 |
| show-latest-tree | MATCH | 1 | 1 |
| show-latest-platform | MATCH | 0 | 0 |
| show-latest-v2-metadata | MATCH | 0 | 0 |
| show-latest-dev-branch | MATCH | 0 | 0 |
| show-latest-dev-branch-major | MATCH | 0 | 0 |
| show-latest-dev-branch-package | MATCH | 0 | 0 |
| show-latest-no-composer | DIFF | 0 | 1 |
| show-latest-cooldown-long | MATCH | 1 | 1 |
| show-latest-cooldown-long-json | MATCH | 1 | 1 |
| show-latest-cooldown-short | MATCH | 1 | 1 |
| show-latest-cooldown-zero | MATCH | 0 | 0 |
| show-latest-cooldown-invalid | MATCH | 1 | 1 |
| show-latest-cooldown-package | MATCH | 1 | 1 |
| show-latest-package | MATCH | 0 | 0 |
| show-latest-package-json | MATCH | 0 | 0 |
| show-latest-package-uptodate | MATCH | 0 | 0 |
| show-latest-package-abandoned | MATCH | 0 | 0 |
| show-latest-package-abandoned-json | MATCH | 0 | 0 |
| show-latest-package-gone | MATCH | 0 | 0 |
| show-latest-package-strict | MATCH | 1 | 1 |
| show-latest-package-strict-uptodate | MATCH | 0 | 0 |
| show-latest-package-major | MATCH | 0 | 0 |
| show-latest-package-patch | MATCH | 0 | 0 |
| show-latest-package-path | MATCH | 1 | 1 |
| show-latest-package-direct-transitive | MATCH | 1 | 1 |
| show-latest-package-platform-rejected | MATCH | 0 | 0 |
| show-latest-all | MATCH | 0 | 0 |
| show-latest-all-json | MATCH | 0 | 0 |
| show-latest-all-name-only | MATCH | 0 | 0 |
| show-latest-all-wildcard | MATCH | 0 | 0 |
| show-latest-all-package | MATCH | 0 | 0 |
| show-latest-all-package-version | MATCH | 0 | 0 |
| show-latest-all-no-lock | MATCH | 0 | 0 |
| show-latest-all-latest | MATCH | 0 | 0 |
| show-latest-all-strict | MATCH | 1 | 1 |
| show-latest-available | MATCH | 0 | 0 |
| show-latest-available-json | MATCH | 0 | 0 |
| show-latest-available-wildcard | MATCH | 0 | 0 |
| show-latest-available-package | MATCH | 0 | 0 |
| show-latest-available-package-json | MATCH | 0 | 0 |
| show-latest-available-not-found | MATCH | 1 | 1 |
| show-latest-available-v2 | MATCH | 0 | 0 |
| show-latest-available-direct | MATCH | 1 | 1 |
| show-latest-available-tree | MATCH | 1 | 1 |
| show-latest-available-latest | MATCH | 0 | 0 |
| outdated-default | MATCH | 0 | 0 |
| outdated-all | MATCH | 0 | 0 |
| outdated-direct | MATCH | 0 | 0 |
| outdated-strict | MATCH | 1 | 1 |
| outdated-strict-fresh | MATCH | 0 | 0 |
| outdated-major-only | MATCH | 0 | 0 |
| outdated-minor-only | MATCH | 0 | 0 |
| outdated-patch-only | MATCH | 0 | 0 |
| outdated-only-conflict | MATCH | 1 | 1 |
| outdated-json | MATCH | 0 | 0 |
| outdated-all-json | MATCH | 0 | 0 |
| outdated-bad-format | MATCH | 1 | 1 |
| outdated-package | MATCH | 0 | 0 |
| outdated-package-strict | MATCH | 1 | 1 |
| outdated-package-not-found | MATCH | 1 | 1 |
| outdated-wildcard | MATCH | 0 | 0 |
| outdated-locked | MATCH | 0 | 0 |
| outdated-no-dev | MATCH | 0 | 0 |
| outdated-ignore | MATCH | 0 | 0 |
| outdated-ignore-platform-reqs | MATCH | 0 | 0 |
| outdated-ignore-platform-req | MATCH | 0 | 0 |
| outdated-sort-by-age | MATCH | 0 | 0 |
| outdated-cooldown | MATCH | 1 | 1 |
| outdated-ansi | DIFF | 0 | 0 |
| outdated-not-installed | MATCH | 0 | 0 |
| outdated-no-composer | DIFF | 1 | 1 |
| outdated-unknown-option | MATCH | 1 | 1 |

## show-latest-ansi

```
stdout:
+ (go) latestfix/aband    1.0.0 [31m1.1.0[0m Abandoned in its latest release
+ (go) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[0m
+ (go) latestfix/core[0m     1.0.0 [33m2.0.0[0m Core library for the latest checks
+ (go) latestfix/devtool  1.0.0 [31m1.0.1[0m A dev tool with a very long description that ...
+ (go) latestfix/gone     1.0.0 [32m1.0.0[0m Abandoned without replacement
+ (go) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[0m
+ (go) latestfix/platneed 1.0.0 [32m1.0.0[0m Newer release needs a future PHP
+ (go) latestfix/stable[0m   1.0.0 [32m1.0.0[0m Already up to date
+ (go) latestfix/zero     0.1.0 [33m0.2.0[0m Zero major package
+ (go) latestfix/dep      1.0.0 [31m1.2.0[0m Transitive dependency
- (oracle) latestfix/aband    1.0.0 [31m1.1.0[39m Abandoned in its latest release
- (oracle) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[39;49m
- (oracle) ]8;;https://example.org/core\latestfix/core]8;;\     1.0.0 [33m2.0.0[39m Core library for the latest checks
- (oracle) latestfix/devtool  1.0.0 [31m1.0.1[39m A dev tool with a very long description that ...
- (oracle) latestfix/gone     1.0.0 [32m1.0.0[39m Abandoned without replacement
- (oracle) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[39;49m
- (oracle) latestfix/platneed 1.0.0 [32m1.0.0[39m Newer release needs a future PHP
- (oracle) ]8;;https://example.org/stable/src\latestfix/stable]8;;\   1.0.0 [32m1.0.0[39m Already up to date
- (oracle) latestfix/zero     0.1.0 [33m0.2.0[39m Zero major package
- (oracle) latestfix/dep      1.0.0 [31m1.2.0[39m Transitive dependency
stderr:
+ (go) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[0m
+ (go) [32mColor legend:[0m
+ (go) - [31mpatch or minor[0m release available - update recommended
+ (go) - [33mmajor[0m release available - update possible
+ (go) - [32mup to date[0m version
- (oracle) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[39;49m
- (oracle) [32mColor legend:[39m
- (oracle) - [31mpatch or minor[39m release available - update recommended
- (oracle) - [33mmajor[39m release available - update possible
- (oracle) - [32mup to date[39m version
+ (go) [32mDirect dependencies required in composer.json:[0m
- (oracle) [32mDirect dependencies required in composer.json:[39m
+ (go) [32mTransitive dependencies not required in composer.json:[0m
- (oracle) [32mTransitive dependencies not required in composer.json:[39m
```

## show-latest-outdated-ansi

```
stdout:
+ (go) latestfix/aband   1.0.0 [31m1.1.0[0m Abandoned in its latest release
+ (go) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[0m
+ (go) latestfix/core[0m    1.0.0 [33m2.0.0[0m Core library for the latest checks
+ (go) latestfix/devtool 1.0.0 [31m1.0.1[0m A dev tool with a very long description that w...
+ (go) latestfix/gone    1.0.0 [32m1.0.0[0m Abandoned without replacement
+ (go) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[0m
+ (go) latestfix/zero    0.1.0 [33m0.2.0[0m Zero major package
+ (go) latestfix/dep     1.0.0 [31m1.2.0[0m Transitive dependency
- (oracle) latestfix/aband   1.0.0 [31m1.1.0[39m Abandoned in its latest release
- (oracle) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[39;49m
- (oracle) ]8;;https://example.org/core\latestfix/core]8;;\    1.0.0 [33m2.0.0[39m Core library for the latest checks
- (oracle) latestfix/devtool 1.0.0 [31m1.0.1[39m A dev tool with a very long description that w...
- (oracle) latestfix/gone    1.0.0 [32m1.0.0[39m Abandoned without replacement
- (oracle) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[39;49m
- (oracle) latestfix/zero    0.1.0 [33m0.2.0[39m Zero major package
- (oracle) latestfix/dep     1.0.0 [31m1.2.0[39m Transitive dependency
stderr:
+ (go) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[0m
+ (go) [32mColor legend:[0m
+ (go) - [31mpatch or minor[0m release available - update recommended
+ (go) - [33mmajor[0m release available - update possible
- (oracle) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[39;49m
- (oracle) [32mColor legend:[39m
- (oracle) - [31mpatch or minor[39m release available - update recommended
- (oracle) - [33mmajor[39m release available - update possible
+ (go) [32mDirect dependencies required in composer.json:[0m
- (oracle) [32mDirect dependencies required in composer.json:[39m
+ (go) [32mTransitive dependencies not required in composer.json:[0m
- (oracle) [32mTransitive dependencies not required in composer.json:[39m
```

## show-latest-no-composer

```
exit code: oracle=0 go=1
stdout:
- (oracle) {
- (oracle)     "platform": [
- (oracle)         {
- (oracle)             "name": "composer",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2.11.999-dev+source",
- (oracle)             "description": "Composer package",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "composer-plugin-api",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2.11.0",
- (oracle)             "description": "The Composer Plugin API",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "composer-runtime-api",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2.2.2",
- (oracle)             "description": "The Composer Runtime API",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-ctype",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The ctype PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-curl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The curl PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-date",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The date PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-dom",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "20031129",
- (oracle)             "description": "The dom PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-fileinfo",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The fileinfo PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-filter",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The filter PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-hash",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The hash PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-iconv",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The iconv PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-json",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The json PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-libxml",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The libxml PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-mbstring",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The mbstring PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-mysqlnd",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "0",
- (oracle)             "description": "The mysqlnd PHP extension (actual version: mysqlnd 8.4.27)",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-openssl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The openssl PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-pcre",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The pcre PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-pdo",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The PDO PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-pdo_sqlite",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The pdo_sqlite PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-phar",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The Phar PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-posix",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The posix PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-random",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The random PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-readline",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The readline PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-reflection",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The Reflection PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-session",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The session PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-simplexml",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The SimpleXML PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-sodium",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The sodium PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-spl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The SPL PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-sqlite3",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The sqlite3 PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-tokenizer",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The tokenizer PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-xml",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The xml PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-xmlreader",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The xmlreader PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-xmlwriter",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The xmlwriter PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-zend-opcache",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The Zend OPcache PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "ext-zlib",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The zlib PHP extension",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-curl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.14.1",
- (oracle)             "description": "The curl library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-curl-libssh2",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "1.11.1",
- (oracle)             "description": "curl libssh2 version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-curl-openssl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "3.5.7",
- (oracle)             "description": "curl OpenSSL version (3.5.7)",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-curl-zlib",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "1.3.1",
- (oracle)             "description": "curl zlib version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-date-timelib",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2022.17",
- (oracle)             "description": "date timelib version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-date-zoneinfo",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2026.5",
- (oracle)             "description": "zoneinfo (\"Olson\") database for date",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-fileinfo-libmagic",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "545",
- (oracle)             "description": "fileinfo libmagic version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-iconv",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2.41",
- (oracle)             "description": "The iconv library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-libsodium",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "1.0.18",
- (oracle)             "description": "The libsodium library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-libxml",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "2.9.14",
- (oracle)             "description": "libxml library version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-mbstring-libmbfl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "1.3.2",
- (oracle)             "description": "mbstring libmbfl version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-mbstring-oniguruma",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "6.9.9",
- (oracle)             "description": "mbstring oniguruma version",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-openssl",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "3.5.7",
- (oracle)             "description": "OpenSSL 3.5.7 9 Jun 2026",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-pcre",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "10.44",
- (oracle)             "description": "The pcre library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-pcre-unicode",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "15.0.0",
- (oracle)             "description": "PCRE Unicode version support",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-pdo_sqlite-sqlite",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "3.46.1",
- (oracle)             "description": "The pdo_sqlite-sqlite library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-sqlite3-sqlite",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "3.46.1",
- (oracle)             "description": "The sqlite3-sqlite library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "lib-zlib",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "1.3.1",
- (oracle)             "description": "The zlib library",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "php",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The PHP interpreter",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "php-64bit",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The PHP interpreter, 64bit",
- (oracle)             "abandoned": false
- (oracle)         },
- (oracle)         {
- (oracle)             "name": "php-ipv6",
- (oracle)             "direct-dependency": false,
- (oracle)             "homepage": null,
- (oracle)             "source": null,
- (oracle)             "version": "8.4.27",
- (oracle)             "description": "The PHP interpreter, with IPv6 support",
- (oracle)             "abandoned": false
- (oracle)         }
- (oracle)     ]
- (oracle) }
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Composer could not find a composer.json file in /work/project                
+ (go)   To initialize a project, please create a composer.json file. See https://ge  
+ (go)   tcomposer.org/basic-usage                                                    
+ (go)                                                                                
+ (go) 
+ (go) show [--all] [--locked] [-i|--installed] [-p|--platform] [-a|--available] [-s|--self] [-N|--name-only] [-P|--path] [-t|--tree] [-l|--latest] [-o|--outdated] [--ignore IGNORE] [-M|--major-only] [-m|--minor-only] [--patch-only] [-A|--sort-by-age] [--cooldown-period COOLDOWN-PERIOD] [-D|--direct] [--strict] [-f|--format FORMAT] [--no-dev] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<package> [<version>]]
+ (go) 
- (oracle) No composer.json found in the current directory, disabling "latest" option
```

## outdated-ansi

```
stdout:
+ (go) latestfix/aband   1.0.0 [31m1.1.0[0m Abandoned in its latest release
+ (go) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[0m
+ (go) latestfix/core[0m    1.0.0 [33m2.0.0[0m Core library for the latest checks
+ (go) latestfix/devtool 1.0.0 [31m1.0.1[0m A dev tool with a very long description that w...
+ (go) latestfix/gone    1.0.0 [32m1.0.0[0m Abandoned without replacement
+ (go) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[0m
+ (go) latestfix/zero    0.1.0 [33m0.2.0[0m Zero major package
+ (go) latestfix/dep     1.0.0 [31m1.2.0[0m Transitive dependency
- (oracle) latestfix/aband   1.0.0 [31m1.1.0[39m Abandoned in its latest release
- (oracle) [30;43mPackage latestfix/aband is abandoned, you should avoid using it. Use latestfix/stable instead.[39;49m
- (oracle) ]8;;https://example.org/core\latestfix/core]8;;\    1.0.0 [33m2.0.0[39m Core library for the latest checks
- (oracle) latestfix/devtool 1.0.0 [31m1.0.1[39m A dev tool with a very long description that w...
- (oracle) latestfix/gone    1.0.0 [32m1.0.0[39m Abandoned without replacement
- (oracle) [30;43mPackage latestfix/gone is abandoned, you should avoid using it. No replacement was suggested.[39;49m
- (oracle) latestfix/zero    0.1.0 [33m0.2.0[39m Zero major package
- (oracle) latestfix/dep     1.0.0 [31m1.2.0[39m Transitive dependency
stderr:
+ (go) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[0m
+ (go) [32mColor legend:[0m
+ (go) - [31mpatch or minor[0m release available - update recommended
+ (go) - [33mmajor[0m release available - update possible
- (oracle) [30;43mWarning: Accessing 127.0.0.1 over http which is an insecure protocol.[39;49m
- (oracle) [32mColor legend:[39m
- (oracle) - [31mpatch or minor[39m release available - update recommended
- (oracle) - [33mmajor[39m release available - update possible
+ (go) [32mDirect dependencies required in composer.json:[0m
- (oracle) [32mDirect dependencies required in composer.json:[39m
+ (go) [32mTransitive dependencies not required in composer.json:[0m
- (oracle) [32mTransitive dependencies not required in composer.json:[39m
```

## outdated-no-composer

```
stderr:
+ (go) 
+ (go) In <location>:
+ (go)                                                                                
+ (go)   Composer could not find a composer.json file in /work/project                
+ (go)   To initialize a project, please create a composer.json file. See https://ge  
+ (go)   tcomposer.org/basic-usage                                                    
+ (go)                                                                                
+ (go) 
+ (go) show [--all] [--locked] [-i|--installed] [-p|--platform] [-a|--available] [-s|--self] [-N|--name-only] [-P|--path] [-t|--tree] [-l|--latest] [-o|--outdated] [--ignore IGNORE] [-M|--major-only] [-m|--minor-only] [--patch-only] [-A|--sort-by-age] [--cooldown-period COOLDOWN-PERIOD] [-D|--direct] [--strict] [-f|--format FORMAT] [--no-dev] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<package> [<version>]]
+ (go) 
- (oracle) Composer could not find a composer.json file in /work/project
- (oracle) To initialize a project, please create a composer.json file. See https://getcomposer.org/basic-usage
```
