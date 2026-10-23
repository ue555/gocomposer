# Composer vs gocomposer comparison

- gocomposer_sha256: `9e970177a636b81c8bd3854b020cf951f27da0fb7db8f03a233eef6c8ced4ad8`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| init-full | MATCH | 0 | 0 |
| init-full-reordered | MATCH | 0 | 0 |
| init-require | MATCH | 0 | 0 |
| init-require-pair-args | MATCH | 0 | 0 |
| init-require-no-version | MATCH | 1 | 1 |
| init-require-dev-no-version | MATCH | 1 | 1 |
| init-autoload | MATCH | 0 | 0 |
| init-autoload-nested | MATCH | 0 | 0 |
| init-autoload-with-deps | MATCH | 0 | 0 |
| init-autoload-dev-deps | MATCH | 0 | 0 |
| init-git-defaults | MATCH | 0 | 0 |
| init-git-defaults-description | MATCH | 0 | 0 |
| init-env-defaults | MATCH | 0 | 0 |
| init-env-overrides-git | MATCH | 0 | 0 |
| init-user-env | MATCH | 0 | 0 |
| init-no-defaults | MATCH | 0 | 0 |
| init-camel-dir | MATCH | 0 | 0 |
| init-invalid-name | MATCH | 1 | 1 |
| init-invalid-name-noslash | MATCH | 1 | 1 |
| init-invalid-author | MATCH | 1 | 1 |
| init-invalid-email | MATCH | 1 | 1 |
| init-empty-author | MATCH | 1 | 1 |
| init-author-name-only | MATCH | 0 | 0 |
| init-existing-json | MATCH | 0 | 0 |
| init-existing-git | ERROR | 0 | 0 |
| init-repository | MATCH | 0 | 0 |
| init-repository-invalid | MATCH | 1 | 1 |
| init-repository-bad-json | MATCH | 1 | 1 |
| init-repository-composer-file | MATCH | 0 | 0 |
| init-repository-filesystem-file | MATCH | 1 | 1 |
| init-repository-missing-file | DIFF | 1 | 1 |
| init-bad-stability | MATCH | 1 | 1 |
| init-bad-type-schema | MATCH | 1 | 1 |
| init-bad-license | MATCH | 0 | 0 |
| init-composer-env | MATCH | 0 | 0 |
| init-composer-env-subdir | MATCH | 0 | 0 |
| init-unicode | MATCH | 0 | 0 |
| global-require | MATCH | 0 | 0 |
| global-update | MATCH | 0 | 0 |
| global-show | MATCH | 0 | 0 |
| global-show-abbrev | MATCH | 0 | 0 |
| global-options-first | MATCH | 0 | 0 |
| global-remove | MATCH | 0 | 0 |
| global-no-args | MATCH | 1 | 1 |
| global-unknown-option-only | MATCH | 1 | 1 |
| global-unknown-command | DIFF | 1 | 1 |
| global-config-quoted | MATCH | 0 | 0 |
| global-config-equals | MATCH | 0 | 0 |
| global-new-home | MATCH | 0 | 0 |
| global-composer-env | MATCH | 0 | 0 |
| global-validate | MATCH | 2 | 2 |
| global-dump-autoload | MATCH | 0 | 0 |
| global-help | DIFF | 0 | 1 |

## init-existing-git

```
podman run failed: exit status 128

```

## init-repository-missing-file

```
stderr:
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   The "missing.json" file could not be downloaded: Failed to open stream: No   
- (oracle)   such file or directory                                                       
- (oracle)                                                                                
- (oracle) 
```

## global-unknown-command

```
stderr:
+ (go) global <command-name> [<args>...]
+ (go) 
```

## global-help

```
exit code: oracle=0 go=1
stdout:
- (oracle) Description:
- (oracle)   Allows running commands in the global composer dir ($COMPOSER_HOME)
- (oracle) Usage:
- (oracle)   global <command-name> [<args>...]
- (oracle) 
- (oracle) Arguments:
- (oracle)   command-name                   
- (oracle)   args                           
- (oracle) 
- (oracle) Options:
- (oracle)   -h, --help                     Display help for the given command. When no command is given display help for the list command
- (oracle)   -q, --quiet                    Do not output any message
- (oracle)   -V, --version                  Display this application version
- (oracle)       --ansi|--no-ansi           Force (or disable --no-ansi) ANSI output
- (oracle)   -n, --no-interaction           Do not ask any interactive question
- (oracle)       --profile                  Display timing and memory usage information
- (oracle)       --no-plugins               Whether to disable plugins.
- (oracle)       --no-scripts               Skips the execution of all scripts defined in composer.json file.
- (oracle)   -d, --working-dir=WORKING-DIR  If specified, use the given directory as working directory.
- (oracle)       --no-cache                 Prevent use of the cache
- (oracle)   -v|vv|vvv, --verbose           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
- (oracle) 
- (oracle) Help:
- (oracle)   Use this command as a wrapper to run other Composer commands
- (oracle)   within the global context of COMPOSER_HOME.
- (oracle)   
- (oracle)   You can use this to install CLI utilities globally, all you need
- (oracle)   is to add the COMPOSER_HOME/vendor/bin dir to your PATH env var.
- (oracle)   
- (oracle)   COMPOSER_HOME is c:\Users<user>\AppData\Roaming\Composer on Windows
- (oracle)   and /home/<user>/.composer on unix systems.
- (oracle)   
- (oracle)   If your system uses freedesktop.org standards, then it will first check
- (oracle)   XDG_CONFIG_HOME or default to /home/<user>/.config/composer
- (oracle)   
- (oracle)   Note: This path may vary depending on customizations to bin-dir in
- (oracle)   composer.json or the environmental variable COMPOSER_BIN_DIR.
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#global
- (oracle) 
stderr:
+ (go)                                                    
+ (go)   Not enough arguments (missing: "command-name").  
+ (go)                                                    
+ (go) 
+ (go) global <command-name> [<args>...]
+ (go) 
+ (go) 
```
