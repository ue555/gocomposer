# Composer vs gocomposer comparison

- gocomposer_sha256: `2f3d75915ee647f1248ac02751714cfced68e424ac5ee6799bbada66e1524e96`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| fund-text | MATCH | 0 | 0 |
| fund-json | MATCH | 0 | 0 |
| fund-json-short | MATCH | 0 | 0 |
| fund-bad-format | MATCH | 1 | 1 |
| fund-v2 | MATCH | 0 | 0 |
| fund-none | MATCH | 0 | 0 |
| fund-none-json | MATCH | 0 | 0 |
| fund-not-installed | DIFF | 0 | 0 |
| fund-no-json | MATCH | 1 | 1 |
| search-fulltext | MATCH | 0 | 0 |
| search-multi-token | MATCH | 0 | 0 |
| search-keyword | MATCH | 0 | 0 |
| search-only-name | MATCH | 0 | 0 |
| search-only-name-long | MATCH | 0 | 0 |
| search-only-vendor | MATCH | 0 | 0 |
| search-only-vendor-long | MATCH | 0 | 0 |
| search-name-and-vendor | MATCH | 1 | 1 |
| search-type | DIFF | 0 | 0 |
| search-type-other | DIFF | 0 | 0 |
| search-type-vendor | MATCH | 0 | 0 |
| search-json | MATCH | 0 | 0 |
| search-json-vendor | MATCH | 0 | 0 |
| search-json-empty | MATCH | 0 | 0 |
| search-bad-format | MATCH | 1 | 1 |
| search-no-results | MATCH | 0 | 0 |
| search-regex-chars | MATCH | 0 | 0 |
| search-regex-chars-name | MATCH | 0 | 0 |
| search-no-tokens | DIFF | 1 | 0 |
| search-very-verbose | DIFF | 0 | 0 |
| search-columns | MATCH | 0 | 0 |
| search-installed | MATCH | 0 | 0 |
| search-installed-vv | DIFF | 0 | 0 |
| search-v2 | MATCH | 0 | 0 |
| search-v2-vendor | MATCH | 0 | 0 |
| search-v2-json | MATCH | 0 | 0 |
| search-wide | MATCH | 0 | 0 |
| search-filtered-repo | MATCH | 0 | 0 |

## fund-not-installed

```
stderr:
- (oracle) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
```

## search-type

```
stdout:
- (oracle) fundco/beta     ! Abandoned ! Beta package of the fund c...
```

## search-type-other

```
stdout:
+ (go) fundco/beta ! Abandoned ! Beta package of the fund compa...
- (oracle) fundco/alpha    Alpha package with funding
- (oracle) fundco/beta     ! Abandoned ! Beta package of the fund c...
- (oracle) fundco/plain    No funding here
- (oracle) otherfund/gamma Gamma
```

## search-no-tokens

```
exit code: oracle=1 go=0
stdout:
+ (go) composer               Composer package
+ (go) composer-plugin-api    The Composer Plugin API
+ (go) composer-runtime-api   The Composer Runtime API
+ (go) php                    The PHP interpreter
+ (go) php-64bit              The PHP interpreter, 64bit
+ (go) php-ipv6               The PHP interpreter, with IPv6 support
+ (go) ext-date               The date PHP extension
+ (go) ext-libxml             The libxml PHP extension
+ (go) ext-openssl            The openssl PHP extension
+ (go) ext-pcre               The pcre PHP extension
+ (go) ext-sqlite3            The sqlite3 PHP extension
+ (go) ext-zlib               The zlib PHP extension
+ (go) ext-ctype              The ctype PHP extension
+ (go) ext-curl               The curl PHP extension
+ (go) ext-dom                The dom PHP extension
+ (go) ext-fileinfo           The fileinfo PHP extension
+ (go) ext-filter             The filter PHP extension
+ (go) ext-hash               The hash PHP extension
+ (go) ext-iconv              The iconv PHP extension
+ (go) ext-json               The json PHP extension
+ (go) ext-mbstring           The mbstring PHP extension
+ (go) ext-spl                The SPL PHP extension
+ (go) ext-session            The session PHP extension
+ (go) ext-pdo                The PDO PHP extension
+ (go) ext-pdo_sqlite         The pdo_sqlite PHP extension
+ (go) ext-posix              The posix PHP extension
+ (go) ext-random             The random PHP extension
+ (go) ext-readline           The readline PHP extension
+ (go) ext-reflection         The Reflection PHP extension
+ (go) ext-phar               The Phar PHP extension
+ (go) ext-simplexml          The SimpleXML PHP extension
+ (go) ext-tokenizer          The tokenizer PHP extension
+ (go) ext-xml                The xml PHP extension
+ (go) ext-xmlreader          The xmlreader PHP extension
+ (go) ext-xmlwriter          The xmlwriter PHP extension
+ (go) ext-mysqlnd            The mysqlnd PHP extension (actual version: mysqlnd 8...
+ (go) ext-sodium             The sodium PHP extension
+ (go) ext-zend-opcache       The Zend OPcache PHP extension
+ (go) lib-date-timelib       date timelib version
+ (go) lib-date-zoneinfo      zoneinfo ("Olson") database for date
+ (go) lib-libxml             libxml library version
+ (go) lib-openssl            OpenSSL 3.5.7 9 Jun 2026
+ (go) lib-pcre               The pcre library
+ (go) lib-pcre-unicode       PCRE Unicode version support
+ (go) lib-sqlite3-sqlite     The sqlite3-sqlite library
+ (go) lib-zlib               The zlib library
+ (go) lib-curl               The curl library
+ (go) lib-curl-openssl       curl OpenSSL version (3.5.7)
+ (go) lib-curl-libssh2       curl libssh2 version
+ (go) lib-curl-zlib          curl zlib version
+ (go) lib-fileinfo-libmagic  fileinfo libmagic version
+ (go) lib-iconv              The iconv library
+ (go) lib-mbstring-libmbfl   mbstring libmbfl version
+ (go) lib-mbstring-oniguruma mbstring oniguruma version
+ (go) lib-pdo_sqlite-sqlite  The pdo_sqlite-sqlite library
+ (go) lib-libsodium          The libsodium library
+ (go) acme/badsum            Dist checksum is wrong
+ (go) acme/beta              Has betas
+ (go) acme/broken            Dist archive is corrupt
+ (go) acme/cli               CLI tool
+ (go) acme/conflicter        Conflicts with acme/log 2
+ (go) acme/devtools          Dev tools
+ (go) acme/edge              Branch only
+ (go) acme/evil              Archive with path traversal
+ (go) acme/http              HTTP client
+ (go) acme/log               Logging
+ (go) acme/log-impl          Provides psr/log-implementation
+ (go) acme/missing           Dist returns 404
+ (go) acme/needs-ext         Needs an extension that is missing
+ (go) acme/needs-php9        Needs a future PHP
+ (go) acme/new               Replaces acme/old
+ (go) acme/old               Old package
+ (go) acme/util              Utilities
+ (go) fundco/alpha           Alpha package with funding
+ (go) fundco/beta            ! Abandoned ! Beta package of the...
+ (go) fundco/plain           No funding here
+ (go) latestfix/aband        Abandoned in its latest release
+ (go) latestfix/core         Core library for the latest checks
+ (go) latestfix/dep          Transitive dependency
+ (go) latestfix/devtool      A dev tool with a very long description that will ne...
+ (go) latestfix/gone         ! Abandoned ! Abandoned without r...
+ (go) latestfix/platneed     Newer release needs a future PHP
+ (go) latestfix/platonly     Only release needs a future PHP
+ (go) latestfix/stable       Already up to date
+ (go) latestfix/zero         Zero major package
+ (go) otherfund/gamma        Gamma
stderr:
+ (go) Warning: Accessing 127.0.0.1 over http which is an insecure protocol.
- (oracle)                                              
- (oracle)   Not enough arguments (missing: "tokens").  
- (oracle)                                              
- (oracle) 
- (oracle) search [-N|--only-name] [-O|--only-vendor] [-t|--type TYPE] [-f|--format FORMAT] [--] <tokens>...
- (oracle) 
- (oracle) 
```

## search-very-verbose

```
stderr:
+ (go) Searched composite repo (installed array repo (defining 0 package), platform repo), found 0 result(s)
- (oracle) Searched installed array repo (defining 0 package), found 0 result(s)
- (oracle) Searched platform repo, found 0 result(s)
```

## search-installed-vv

```
stderr:
+ (go) Searched composite repo (installed array repo (defining 1 package), platform repo), found 1 result(s)
- (oracle) Searched installed array repo (defining 1 package), found 1 result(s)
- (oracle) Searched platform repo, found 0 result(s)
```
