# Composer vs gocomposer comparison

- gocomposer_sha256: `f17ab67478d944bf5e0d6cc8a6a918e8ad322b7ba2fcb32bc1d3c269990868be`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| init-existing-git | ERROR | 0 | 0 |
| init-repository-missing-file | MATCH | 1 | 1 |
| global-unknown-command | MATCH | 1 | 1 |
| global-help | DIFF | 0 | 0 |

## init-existing-git

```
podman run failed: exit status 128

```

## global-help

```
stdout:
+ (go)   global [options] <command-name> [--] [<args>...]
- (oracle)   global <command-name> [<args>...]
+ (go)   command-name                 
+ (go)   args                         
- (oracle)   command-name                   
- (oracle)   args                           
+ (go)   -h, --help                   Display help for the given command. When no command is given display help for the list command
+ (go)   -q, --quiet                  Do not output any message
+ (go)   -v|-vv|-vvv, --verbose       Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   -V, --version                Display this application version
+ (go)       --ansi                   Force (or disable --no-ansi) ANSI output
+ (go)   -n, --no-interaction         Do not ask any interactive question
+ (go)       --profile                Display timing and memory usage information
+ (go)       --no-plugins             Whether to disable plugins.
+ (go)       --no-scripts             Skips the execution of all scripts defined in composer.json file.
+ (go)   -d, --working-dir=WORKING-DIR If specified, use the given directory as working directory.
+ (go)       --no-cache               Prevent use of the cache
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
```
