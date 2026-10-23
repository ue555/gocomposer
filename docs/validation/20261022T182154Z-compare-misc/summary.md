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
| run-script-basic | DIFF | 0 | 0 |
| run-script-alias-run | DIFF | 0 | 0 |
| run-script-abbrev | DIFF | 0 | 0 |
| run-script-args | DIFF | 0 | 0 |
| run-script-args-dashdash | DIFF | 0 | 0 |
| run-script-list | DIFF | 0 | 0 |
| run-script-list-short | DIFF | 0 | 0 |
| run-script-missing-argument | DIFF | 1 | 1 |
| run-script-not-defined | DIFF | 1 | 1 |
| run-script-event-constant-name | DIFF | 1 | 1 |
| run-script-event | DIFF | 0 | 0 |
| run-script-event-array | DIFF | 0 | 0 |
| run-script-package-event | DIFF | 0 | 0 |
| run-script-dev-mode | DIFF | 0 | 0 |
| run-script-no-dev | DIFF | 0 | 0 |
| run-script-dev-and-no-dev | DIFF | 0 | 0 |
| run-script-no-dev-env | DIFF | 0 | 0 |
| run-script-timeout-invalid | DIFF | 1 | 1 |
| run-script-timeout-negative | DIFF | 1 | 1 |
| run-script-timeout | DIFF | 1 | 1 |
| run-script-timeout-zero | DIFF | 0 | 0 |
| run-script-process-timeout-env | DIFF | 1 | 1 |
| run-script-fail | DIFF | 3 | 3 |
| run-script-nested | DIFF | 0 | 0 |
| run-script-nested-verbose | DIFF | 0 | 0 |
| run-script-missing-ref | DIFF | 0 | 0 |
| run-script-multi-args | DIFF | 0 | 0 |
| run-script-no-additional-args | DIFF | 0 | 0 |
| run-script-putenv | DIFF | 0 | 0 |
| run-script-no-scripts | DIFF | 1 | 1 |
| run-script-skip-env | DIFF | 1 | 1 |
| run-script-skip-env-verbose | DIFF | 1 | 1 |
| run-script-quiet | DIFF | 0 | 0 |
| run-script-direct | DIFF | 0 | 0 |
| run-script-direct-args | DIFF | 0 | 0 |
| run-script-direct-unknown-option | DIFF | 0 | 0 |
| run-script-direct-dashdash | DIFF | 0 | 0 |
| run-script-direct-alias | DIFF | 0 | 0 |
| run-script-direct-alias2 | DIFF | 0 | 0 |
| run-script-direct-nested | DIFF | 0 | 0 |
| run-script-direct-fail | DIFF | 3 | 3 |
| run-script-direct-dev-mode | DIFF | 0 | 0 |
| run-script-direct-event-not-command | DIFF | 1 | 1 |
| run-script-direct-abbrev | DIFF | 0 | 0 |
| run-script-unknown-command | DIFF | 1 | 1 |
| run-script-override-warning | MATCH | 0 | 0 |
| run-script-aliases-invalid | MATCH | 1 | 1 |
| run-script-list-empty | MATCH | 0 | 0 |
| run-script-no-json | DIFF | 1 | 1 |
| run-script-disable-timeout | MATCH | 0 | 0 |
| run-script-working-dir | MATCH | 0 | 0 |
| exec-list | MATCH | 0 | 0 |
| exec-no-binary | MATCH | 0 | 0 |
| exec-run | MATCH | 0 | 0 |
| exec-run-args | MATCH | 0 | 0 |
| exec-run-php | MATCH | 0 | 0 |
| exec-missing | DIFF | 127 | 127 |
| exec-no-scripts | MATCH | 0 | 0 |
| exec-verbose | MATCH | 0 | 0 |
| exec-no-bins | MATCH | 1 | 1 |
| exec-local-list | MATCH | 0 | 0 |
| exec-no-json | DIFF | 1 | 1 |

## run-script-basic

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-alias-run

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-abbrev

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-args

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-args-dashdash

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-list

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-list-short

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-missing-argument

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-not-defined

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-event-constant-name

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-event

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-event-array

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:80
- (oracle)  /oracle-source/src/Composer/Package/Loader/RootPackageLoader.php:134
- (oracle)  /oracle-source/src/Composer/Factory.php:394
- (oracle)  /oracle-source/src/Composer/Factory.php:627
- (oracle)  /oracle-source/src/Composer/Console/Application.php:676
- (oracle)  /oracle-source/src/Composer/Console/Application.php:443
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## run-script-package-event

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-dev-mode

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-no-dev

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-dev-and-no-dev

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-no-dev-env

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-timeout-invalid

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-timeout-negative

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-timeout

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle) The following exception is caused by a process timeout
- (oracle) Check https://getcomposer.org/doc/06-config.md#process-timeout for details
```

## run-script-timeout-zero

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-process-timeout-env

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle) The following exception is caused by a process timeout
- (oracle) Check https://getcomposer.org/doc/06-config.md#process-timeout for details
```

## run-script-fail

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-nested

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-nested-verbose

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:80
- (oracle)  /oracle-source/src/Composer/Package/Loader/RootPackageLoader.php:134
- (oracle)  /oracle-source/src/Composer/Factory.php:394
- (oracle)  /oracle-source/src/Composer/Factory.php:627
- (oracle)  /oracle-source/src/Composer/Console/Application.php:676
- (oracle)  /oracle-source/src/Composer/Console/Application.php:443
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
```

## run-script-missing-ref

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-multi-args

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-no-additional-args

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-putenv

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-no-scripts

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-skip-env

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-skip-env-verbose

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle) Stack trace:
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
- (oracle)  /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:80
- (oracle)  /oracle-source/src/Composer/Package/Loader/RootPackageLoader.php:134
- (oracle)  /oracle-source/src/Composer/Factory.php:394
- (oracle)  /oracle-source/src/Composer/Factory.php:627
- (oracle)  /oracle-source/src/Composer/Console/Application.php:676
- (oracle)  /oracle-source/src/Composer/Console/Application.php:443
- (oracle)  /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  /oracle-source/bin/composer:113
- (oracle)   [InvalidArgumentException]                    
- (oracle) Exception trace:
- (oracle)   at /oracle-source/src/Composer/Command/RunScriptCommand.php:125
- (oracle)  Composer\Command\RunScriptCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## run-script-quiet

```
stdout:
+ (go) test-ran
```

## run-script-direct

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-args

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-unknown-option

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-dashdash

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-alias

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-alias2

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-nested

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-fail

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-dev-mode

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-event-not-command

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-direct-abbrev

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-unknown-command

```
stderr:
- (oracle) Deprecation Notice: The `putenv` script name is reserved for internal use, please avoid defining it in /oracle-source/src/Composer/Package/Loader/ArrayLoader.php:281
```

## run-script-no-json

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
+ (go) run-script [--timeout TIMEOUT] [--dev] [--no-dev] [-l|--list] [--] [<script> [<args>...]]
+ (go) 
- (oracle) Composer could not find a composer.json file in /work/project
- (oracle) To initialize a project, please create a composer.json file. See https://getcomposer.org/basic-usage
```

## exec-missing

```
stderr:
+ (go) /bin/sh: 1: nope: not found
- (oracle) sh: 1: nope: not found
```

## exec-no-json

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
+ (go) exec [-l|--list] [--] [<binary> [<args>...]]
+ (go) 
- (oracle) Composer could not find a composer.json file in /work/project
- (oracle) To initialize a project, please create a composer.json file. See https://getcomposer.org/basic-usage
```
