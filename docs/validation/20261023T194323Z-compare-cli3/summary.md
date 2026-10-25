# Composer vs gocomposer comparison

- gocomposer_sha256: `32348ea841c6228e2045b8ee63c9fe2f32c0edf32c31a4c2d3faca2187eb52fe`
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
| cli3-vvv-validate | MATCH | 0 | 0 |
| cli3-vvv-show | MATCH | 0 | 0 |
| cli3-vvv-config-list | MATCH | 0 | 0 |
| cli3-vvv-config-get | MATCH | 0 | 0 |
| cli3-vvv-config-set | MATCH | 0 | 0 |
| cli3-vvv-config-extra | MATCH | 0 | 0 |
| cli3-vvv-config-global | MATCH | 0 | 0 |
| cli3-vvv-dump-autoload | MATCH | 0 | 0 |
| cli3-vvv-run-script | MATCH | 0 | 0 |
| cli3-vvv-script-command | MATCH | 0 | 0 |
| cli3-vvv-licenses | MATCH | 0 | 0 |
| cli3-vvv-depends | MATCH | 0 | 0 |
| cli3-vvv-suggests | MATCH | 0 | 0 |
| cli3-vvv-no-plugins | MATCH | 0 | 0 |
| cli3-vvv-shell-verbosity | MATCH | 0 | 0 |
| cli3-vvv-no-lock | MATCH | 0 | 0 |
| cli3-vvv-auth-files | DIFF | 0 | 0 |
| cli3-vvv-cafile | MATCH | 0 | 0 |
| cli3-vvv-ssl-cert-file | MATCH | 0 | 0 |
| cli3-vvv-ssl-cert-missing | MATCH | 0 | 0 |
| cli3-vvv-global-home | DIFF | 0 | 0 |
| cli3-vvv-allow-plugins | MATCH | 0 | 0 |
| cli3-vvv-validate-file | MATCH | 0 | 0 |
| cli3-bin-compat-get | MATCH | 0 | 0 |
| cli3-bin-compat-list | MATCH | 0 | 0 |
| cli3-bin-compat-set | MATCH | 0 | 0 |
| cli3-bin-compat-install | MATCH | 0 | 0 |
| cli3-bin-compat-env | MATCH | 0 | 0 |
| cli3-event-validate | DIFF | 0 | 0 |
| cli3-event-validate-no-plugins | MATCH | 0 | 0 |
| cli3-event-reinstall | MATCH | 0 | 0 |
| cli3-event-reinstall-none | MATCH | 1 | 1 |
| cli3-cmdclass-run | MATCH | 0 | 0 |
| cli3-cmdclass-run-args | MATCH | 0 | 0 |
| cli3-cmdclass-run-verbose | MATCH | 0 | 0 |
| cli3-cmdclass-run-fail | MATCH | 3 | 3 |
| cli3-cmdclass-run-boom | MATCH | 1 | 1 |
| cli3-cmdclass-run-missing | MATCH | 0 | 0 |
| cli3-cmdclass-run-notcmd | MATCH | 0 | 0 |
| cli3-cmdclass-run-named | MATCH | 0 | 0 |
| cli3-cmdclass-run-nested | MATCH | 0 | 0 |
| cli3-cmdclass-reserved | MATCH | 0 | 0 |
| cli3-cmdclass-list | MATCH | 0 | 0 |
| cli3-cmdclass-direct | MATCH | 0 | 0 |
| cli3-cmdclass-direct-options | MATCH | 0 | 0 |
| cli3-cmdclass-direct-bad-option | MATCH | 1 | 1 |
| cli3-cmdclass-direct-fail | MATCH | 3 | 3 |
| cli3-cmdclass-direct-boom | MATCH | 1 | 1 |
| cli3-cmdclass-direct-named | DIFF | 0 | 0 |
| cli3-cmdclass-not-installed | MATCH | 0 | 0 |
| cli3-packagist-search | MATCH | 0 | 0 |
| cli3-packagist-search-json | MATCH | 0 | 0 |
| cli3-packagist-search-all | DIFF | 1 | 1 |
| cli3-packagist-search-name | MATCH | 0 | 0 |
| cli3-packagist-search-vendor | MATCH | 0 | 0 |
| cli3-packagist-search-none | MATCH | 0 | 0 |
| cli3-packagist-browse-show | MATCH | 0 | 0 |
| cli3-packagist-browse-homepage | MATCH | 0 | 0 |
| cli3-packagist-browse-source | MATCH | 1 | 1 |
| cli3-packagist-browse-missing | MATCH | 1 | 1 |
| cli3-packagist-browse-root | MATCH | 1 | 1 |
| cli3-packagist-fund | DIFF | 0 | 0 |
| cli3-packagist-fund-json | DIFF | 0 | 0 |
| cli3-packagist-fund-no-json | MATCH | 1 | 1 |

## cli3-vvv-auth-files

```
stderr:
- (oracle) Reading /work/project/auth.json
```

## cli3-vvv-global-home

```
stderr:
- (oracle) Reading /work/project/ghome/vendor/composer/installed.json
```

## cli3-event-validate

```
stdout:
+ (go) [plugfix/events] deactivate 1.0.0
```

## cli3-cmdclass-direct-named

```
stdout:
+ (go) named command ran as other-name
- (oracle) named command ran as named
```

## cli3-packagist-search-all

```
stderr:
+ (go)                                              
+ (go)   Not enough arguments (missing: "tokens").  
+ (go)                                              
- (oracle) In <location>:
- (oracle)                                                           
- (oracle)   [Symfony\Component\Console\Exception\RuntimeException]  
- (oracle)   Not enough arguments (missing: "tokens").               
- (oracle)                                                           
- (oracle) Exception trace:
- (oracle)   at /oracle-source/vendor/symfony/console/Input/Input.php:76
- (oracle)  Symfony\Component\Console\Input\Input->validate() at /oracle-source/vendor/symfony/console/Command/Command.php:293
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## cli3-packagist-fund

```
stdout:
+ (go) The following packages were found in your dependencies which publish funding information:
+ (go) 
+ (go) stubvendor
+ (go)   alpha
+ (go)     https://github.com/sponsors/alpha
+ (go)     https://alpha.example.org/donate
+ (go) 
+ (go) Please consider following these links and sponsoring the work of package authors!
+ (go) Thank you!
- (oracle) No funding links were found in your package dependencies. This doesn't mean they don't need your support!
```

## cli3-packagist-fund-json

```
stdout:
+ (go) {
+ (go)     "stubvendor": {
+ (go)         "https://github.com/sponsors/alpha": [
+ (go)             "alpha"
+ (go)         ],
+ (go)         "https://alpha.example.org/donate": [
+ (go)             "alpha"
+ (go)         ]
+ (go)     }
+ (go) }
- (oracle) []
```
