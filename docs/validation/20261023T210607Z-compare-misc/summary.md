# Composer vs gocomposer comparison

- gocomposer_sha256: `5baa68ea675e11782e65091a8a3d30f03002b64b5c84ebc13143a96c42ec2352`
- host_os: `darwin/arm64`
- oracle_arch_image: `localhost/gocomposer-oracle:8.4-arch`
- oracle_arch_image_digest: `sha256:93bc04d027eaeb74edeb8ea07a247c5f17cb2e91d4d17fe97681f949f96116e7`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:9cf420cd931841442e55100b424e04a60af91810d670826961c6583bd372308b`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- oracle_zip_image: `localhost/gocomposer-oracle:8.4-zip`
- oracle_zip_image_digest: `sha256:e3a99ff1ce724bd01dbd722f202bf4afa98eec9bde0789c99ad6f76b266e6f22`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| misc-list-root | DIFF | 0 | 0 |
| misc-list-root-noallow | MATCH | 1 | 1 |
| misc-help-root | DIFF | 0 | 0 |
| misc-list-tty | DIFF | 0 | 0 |
| misc-help-tty | DIFF | 0 | 0 |
| misc-help-install-tty | DIFF | 0 | 0 |
| misc-help-script | DIFF | 0 | 0 |
| misc-help-script-list | DIFF | 0 | 0 |
| misc-help-script-tty | DIFF | 0 | 0 |
| misc-list-scripts | DIFF | 0 | 0 |

## misc-list-root

```
stdout:
+ (go) gocomposer version 0.1.0-dev (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle)    ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) Composer version 2.11.999-dev+source @release_date@
+ (go)   -h, --help                Display help for the given command. When no command is given display help for the list command
+ (go)   -q, --quiet               Do not output any message
+ (go)   -v|-vv|-vvv, --verbose    Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   -V, --version             Display this application version
+ (go)       --ansi                Force (or disable --no-ansi) ANSI output
+ (go)   -n, --no-interaction      Do not ask any interactive question
+ (go)       --profile             Display timing and memory usage information
+ (go)       --no-plugins          Whether to disable plugins.
+ (go)       --no-scripts          Skips the execution of all scripts defined in composer.json file.
+ (go)   -d, --working-dir         If specified, use the given directory as working directory.
+ (go)       --no-cache            Prevent use of the cache
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
+ (go)   browse               Opens the package's repository URL or homepage in your browser
- (oracle)   browse               [home] Opens the package's repository URL or homepage in your browser
+ (go)   clear-cache          Clears composer's internal package cache
- (oracle)   clear-cache          [clearcache|cc] Clears composer's internal package cache
+ (go)   depends              Shows which packages cause the given package to be installed
- (oracle)   depends              [why] Shows which packages cause the given package to be installed
+ (go)   dump-autoload        Dumps the autoloader
- (oracle)   dump-autoload        [dumpautoload] Dumps the autoloader
+ (go)   install              Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   install              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   prohibits            Shows which packages prevent the given package from being installed
- (oracle)   prohibits            [why-not] Shows which packages prevent the given package from being installed
+ (go)   remove               Removes a package from the require or require-dev
+ (go)   repository           Manages repositories
+ (go)   require              Adds required packages to your composer.json and installs them
+ (go)   run-script           Runs the scripts defined in composer.json
- (oracle)   remove               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   repository           [repo] Manages repositories
- (oracle)   require              [r] Adds required packages to your composer.json and installs them
- (oracle)   run-script           [run] Runs the scripts defined in composer.json
+ (go)   self-update          Updates composer.phar to the latest version
+ (go)   show                 Shows information about packages
- (oracle)   self-update          [selfupdate] Updates composer.phar to the latest version
- (oracle)   show                 [info] Shows information about packages
+ (go)   update               Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   update               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
```

## misc-help-root

```
stdout:
+ (go)   help [options] [<command_name>]
- (oracle)   help [options] [--] [<command_name>]
+ (go)   command_name                 The command name
- (oracle)   command_name                   The command name [default: "help"]
+ (go)       --format=FORMAT          The output format (txt, xml, json, or md)
+ (go)       --raw                    To output raw command help
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
- (oracle)       --format=FORMAT            The output format (txt, xml, json, or md) [default: "txt"]
- (oracle)       --raw                      To output raw command help
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
- (oracle)   The help command displays help for a given command:
- (oracle)   
- (oracle)     /oracle-source/bin/composer help list
- (oracle)   
- (oracle)   You can also output the help in other formats by using the --format option:
- (oracle)   
- (oracle)     /oracle-source/bin/composer help --format=xml list
- (oracle)   
- (oracle)   To display the list of available commands, please use the list command.
- (oracle) 
```

## misc-list-tty

```
stdout:
+ (go) \e]7501;state=working:app=composer\e\\e[32mgocomposer\e[39m version \e[33m0.1.0-dev\e[39m (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle) \e]7501;state=working:app=composer\e\   ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) \e[32mComposer\e[39m version \e[33m2.11.999-dev+source\e[39m @release_date@
+ (go)   \e[32m-h, --help               \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet              \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose   \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version            \e[39m Display this application version
+ (go)   \e[32m    --ansi               \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction     \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile            \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins         \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts         \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir        \e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache           \e[39m Prevent use of the cache
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32mabout               \e[39m Shows a short information about Composer
+ (go)   \e[32marchive             \e[39m Creates an archive of this composer package
+ (go)   \e[32maudit               \e[39m Checks for security vulnerability advisories for installed packages
+ (go)   \e[32mboth                \e[39m Runs the both script as defined in composer.json
+ (go)   \e[32mbrowse              \e[39m Opens the package's repository URL or homepage in your browser
+ (go)   \e[32mbump                \e[39m Increases the lower limit of your composer.json requirements to the currently installed versions
+ (go)   \e[32mcheck-platform-reqs \e[39m Check that platform requirements are satisfied
+ (go)   \e[32mclear-cache         \e[39m Clears composer's internal package cache
+ (go)   \e[32mcompletion          \e[39m Dump the shell completion script
+ (go)   \e[32mconfig              \e[39m Sets config options
+ (go)   \e[32mcreate-project      \e[39m Creates new project from a package into given directory
+ (go)   \e[32mdepends             \e[39m Shows which packages cause the given package to be installed
+ (go)   \e[32mdiagnose            \e[39m Diagnoses the system to identify common errors
+ (go)   \e[32mdump-autoload       \e[39m Dumps the autoloader
+ (go)   \e[32mexec                \e[39m Executes a vendored binary/script
+ (go)   \e[32mfund                \e[39m Discover how to help fund the maintenance of your dependencies
+ (go)   \e[32mglobal              \e[39m Allows running commands in the global composer dir ($COMPOSER_HOME)
+ (go)   \e[32mgreet               \e[39m Says hi
+ (go)   \e[32mhelp                \e[39m Display help for a command
+ (go)   \e[32minit                \e[39m Creates a basic composer.json file in current directory
+ (go)   \e[32minstall             \e[39m Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   \e[32mlicenses            \e[39m Shows information about licenses of dependencies
+ (go)   \e[32mlist                \e[39m List commands
+ (go)   \e[32moutdated            \e[39m Shows a list of installed packages that have updates available, including their latest version
+ (go)   \e[32mpolicy              \e[39m Manages custom dependency policies and their sources
+ (go)   \e[32mprohibits           \e[39m Shows which packages prevent the given package from being installed
+ (go)   \e[32mreinstall           \e[39m Uninstalls and reinstalls the given package names
+ (go)   \e[32mremove              \e[39m Removes a package from the require or require-dev
+ (go)   \e[32mrepository          \e[39m Manages repositories
+ (go)   \e[32mrequire             \e[39m Adds required packages to your composer.json and installs them
+ (go)   \e[32mrun-script          \e[39m Runs the scripts defined in composer.json
+ (go)   \e[32msearch              \e[39m Searches for packages
+ (go)   \e[32mself-update         \e[39m Updates composer.phar to the latest version
+ (go)   \e[32mshow                \e[39m Shows information about packages
+ (go)   \e[32mstatus              \e[39m Shows a list of locally modified packages
+ (go)   \e[32msuggests            \e[39m Shows package suggestions
+ (go)   \e[32mupdate              \e[39m Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
+ (go)   \e[32mvalidate            \e[39m Validates a composer.json and composer.lock
- (oracle)   \e[32mabout\e[39m                Shows a short information about Composer
- (oracle)   \e[32marchive\e[39m              Creates an archive of this composer package
- (oracle)   \e[32maudit\e[39m                Checks for security vulnerability advisories for installed packages
- (oracle)   \e[32mboth\e[39m                 Runs the both script as defined in composer.json
- (oracle)   \e[32mbrowse\e[39m               [home] Opens the package's repository URL or homepage in your browser
- (oracle)   \e[32mbump\e[39m                 Increases the lower limit of your composer.json requirements to the currently installed versions
- (oracle)   \e[32mcheck-platform-reqs\e[39m  Check that platform requirements are satisfied
- (oracle)   \e[32mclear-cache\e[39m          [clearcache|cc] Clears composer's internal package cache
- (oracle)   \e[32mcompletion\e[39m           Dump the shell completion script
- (oracle)   \e[32mconfig\e[39m               Sets config options
- (oracle)   \e[32mcreate-project\e[39m       Creates new project from a package into given directory
- (oracle)   \e[32mdepends\e[39m              [why] Shows which packages cause the given package to be installed
- (oracle)   \e[32mdiagnose\e[39m             Diagnoses the system to identify common errors
- (oracle)   \e[32mdump-autoload\e[39m        [dumpautoload] Dumps the autoloader
- (oracle)   \e[32mexec\e[39m                 Executes a vendored binary/script
- (oracle)   \e[32mfund\e[39m                 Discover how to help fund the maintenance of your dependencies
- (oracle)   \e[32mglobal\e[39m               Allows running commands in the global composer dir ($COMPOSER_HOME)
- (oracle)   \e[32mgreet\e[39m                Says hi
- (oracle)   \e[32mhelp\e[39m                 Display help for a command
- (oracle)   \e[32minit\e[39m                 Creates a basic composer.json file in current directory
- (oracle)   \e[32minstall\e[39m              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   \e[32mlicenses\e[39m             Shows information about licenses of dependencies
- (oracle)   \e[32mlist\e[39m                 List commands
- (oracle)   \e[32moutdated\e[39m             Shows a list of installed packages that have updates available, including their latest version
- (oracle)   \e[32mpolicy\e[39m               Manages custom dependency policies and their sources
- (oracle)   \e[32mprohibits\e[39m            [why-not] Shows which packages prevent the given package from being installed
- (oracle)   \e[32mreinstall\e[39m            Uninstalls and reinstalls the given package names
- (oracle)   \e[32mremove\e[39m               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   \e[32mrepository\e[39m           [repo] Manages repositories
- (oracle)   \e[32mrequire\e[39m              [r] Adds required packages to your composer.json and installs them
- (oracle)   \e[32mrun-script\e[39m           [run] Runs the scripts defined in composer.json
- (oracle)   \e[32msearch\e[39m               Searches for packages
- (oracle)   \e[32mself-update\e[39m          [selfupdate] Updates composer.phar to the latest version
- (oracle)   \e[32mshow\e[39m                 [info] Shows information about packages
- (oracle)   \e[32mstatus\e[39m               Shows a list of locally modified packages
- (oracle)   \e[32msuggests\e[39m             Shows package suggestions
- (oracle)   \e[32mupdate\e[39m               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   \e[32mvalidate\e[39m             Validates a composer.json and composer.lock
```

## misc-help-tty

```
stdout:
+ (go)   help [options] [<command_name>]
- (oracle)   help [options] [--] [<command_name>]
+ (go)   \e[32mcommand_name                \e[39m The command name
- (oracle)   \e[32mcommand_name\e[39m                   The command name\e[33m [default: "help"]\e[39m
+ (go)   \e[32m    --format=FORMAT         \e[39m The output format (txt, xml, json, or md)
+ (go)   \e[32m    --raw                   \e[39m To output raw command help
+ (go)   \e[32m-h, --help                  \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet                 \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose      \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version               \e[39m Display this application version
+ (go)   \e[32m    --ansi                  \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction        \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile               \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins            \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts            \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir=WORKING-DIR\e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache              \e[39m Prevent use of the cache
- (oracle)   \e[32m    --format=FORMAT\e[39m            The output format (txt, xml, json, or md)\e[33m [default: "txt"]\e[39m
- (oracle)   \e[32m    --raw\e[39m                      To output raw command help
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
- (oracle) 
- (oracle) \e[33mHelp:\e[39m
- (oracle)   The \e[32mhelp\e[39m command displays help for a given command:
- (oracle)   
- (oracle)     \e[32m/oracle-source/bin/composer help list\e[39m
- (oracle)   
- (oracle)   You can also output the help in other formats by using the \e[33m--format\e[39m option:
- (oracle)   
- (oracle)     \e[32m/oracle-source/bin/composer help --format=xml list\e[39m
- (oracle)   
- (oracle)   To display the list of available commands, please use the \e[32mlist\e[39m command.
```

## misc-help-install-tty

```
stdout:
+ (go)   \e[32mpackages                    \e[39m Should not be provided, use composer require instead to add a given package to composer.json.
- (oracle)   \e[32mpackages\e[39m                                             Should not be provided, use composer require instead to add a given package to composer.json.
+ (go)   \e[32m    --prefer-source         \e[39m Forces installation from package sources when possible, including VCS information.
+ (go)   \e[32m    --prefer-dist           \e[39m Forces installation from package dist (default behavior).
+ (go)   \e[32m    --prefer-install=PREFER-INSTALL\e[39m Forces installation from package dist|source|auto (auto chooses source for dev versions, dist for the rest).
+ (go)   \e[32m    --dry-run               \e[39m Outputs the operations but will not execute anything (implicitly enables --verbose).
+ (go)   \e[32m    --download-only         \e[39m Download only, do not install packages.
+ (go)   \e[32m    --dev                   \e[39m DEPRECATED: Enables installation of require-dev packages (enabled by default, only present for BC).
+ (go)   \e[32m    --no-suggest            \e[39m DEPRECATED: This flag does not exist anymore.
+ (go)   \e[32m    --no-dev                \e[39m Disables installation of require-dev packages.
+ (go)   \e[32m    --no-security-blocking  \e[39m DEPRECATED: use --no-blocking instead. Allows installing packages with security advisories, that are abandoned or that are still in the cooldown period (can also be set via the COMPOSER_NO_SECURITY_BLOCKING=1 env var).
+ (go)   \e[32m    --no-blocking           \e[39m Disables all policy blocking during this command, including the cooldown period (can also be set via the COMPOSER_NO_BLOCKING=1 env var).
+ (go)   \e[32m    --no-autoloader         \e[39m Skips autoloader generation
+ (go)   \e[32m    --no-progress           \e[39m Do not output download progress.
+ (go)   \e[32m    --no-install            \e[39m Do not use, only defined here to catch misuse of the install command.
+ (go)   \e[32m    --audit                 \e[39m Run an audit after installation is complete.
+ (go)   \e[32m    --audit-format=AUDIT-FORMAT\e[39m Audit output format. Must be "table", "plain", "json", or "summary".
+ (go)   \e[32m-v|-vv|-vvv, --verbose      \e[39m Shows more details including new commits pulled in when updating packages.
+ (go)   \e[32m-o, --optimize-autoloader   \e[39m Optimize autoloader during autoloader dump
+ (go)   \e[32m-a, --classmap-authoritative\e[39m Autoload classes from the classmap only. Implicitly enables `--optimize-autoloader`.
+ (go)   \e[32m    --strict-psr-autoloader \e[39m Return a failed status code (6) if PSR-4 or PSR-0 mapping errors are present. Requires --optimize-autoloader or --classmap-authoritative to work.
+ (go)   \e[32m    --apcu-autoloader       \e[39m Use APCu to cache found/not-found classes.
+ (go)   \e[32m    --apcu-autoloader-prefix=APCU-AUTOLOADER-PREFIX\e[39m Use a custom prefix for the APCu autoloader cache. Implicitly enables --apcu-autoloader
+ (go)   \e[32m    --ignore-platform-req=IGNORE-PLATFORM-REQ\e[39m Ignore a specific platform requirement (php & ext- packages).
+ (go)   \e[32m    --ignore-platform-reqs  \e[39m Ignore all platform requirements (php & ext- packages).
+ (go)   \e[32m-h, --help                  \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet                 \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose      \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version               \e[39m Display this application version
+ (go)   \e[32m    --ansi                  \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction        \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile               \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins            \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts            \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir=WORKING-DIR\e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache              \e[39m Prevent use of the cache
- (oracle)   \e[32m    --prefer-source\e[39m                                  Forces installation from package sources when possible, including VCS information.
- (oracle)   \e[32m    --prefer-dist\e[39m                                    Forces installation from package dist (default behavior).
- (oracle)   \e[32m    --prefer-install=PREFER-INSTALL\e[39m                  Forces installation from package dist|source|auto (auto chooses source for dev versions, dist for the rest).
- (oracle)   \e[32m    --dry-run\e[39m                                        Outputs the operations but will not execute anything (implicitly enables --verbose).
- (oracle)   \e[32m    --download-only\e[39m                                  Download only, do not install packages.
- (oracle)   \e[32m    --dev\e[39m                                            DEPRECATED: Enables installation of require-dev packages (enabled by default, only present for BC).
- (oracle)   \e[32m    --no-suggest\e[39m                                     DEPRECATED: This flag does not exist anymore.
- (oracle)   \e[32m    --no-dev\e[39m                                         Disables installation of require-dev packages.
- (oracle)   \e[32m    --no-security-blocking\e[39m                           DEPRECATED: use --no-blocking instead. Allows installing packages with security advisories, that are abandoned or that are still in the cooldown period (can also be set via the COMPOSER_NO_SECURITY_BLOCKING=1 env var).
- (oracle)   \e[32m    --no-blocking\e[39m                                    Disables all policy blocking during this command, including the cooldown period (can also be set via the COMPOSER_NO_BLOCKING=1 env var).
- (oracle)   \e[32m    --no-autoloader\e[39m                                  Skips autoloader generation
- (oracle)   \e[32m    --no-progress\e[39m                                    Do not output download progress.
- (oracle)   \e[32m    --no-install\e[39m                                     Do not use, only defined here to catch misuse of the install command.
- (oracle)   \e[32m    --audit\e[39m                                          Run an audit after installation is complete.
- (oracle)   \e[32m    --audit-format=AUDIT-FORMAT\e[39m                      Audit output format. Must be "table", "plain", "json", or "summary".\e[33m [default: "summary"]\e[39m
- (oracle)   \e[32m-o, --optimize-autoloader\e[39m                            Optimize autoloader during autoloader dump
- (oracle)   \e[32m-a, --classmap-authoritative\e[39m                         Autoload classes from the classmap only. Implicitly enables `--optimize-autoloader`.
- (oracle)   \e[32m    --strict-psr-autoloader\e[39m                          Return a failed status code (6) if PSR-4 or PSR-0 mapping errors are present. Requires --optimize-autoloader or --classmap-authoritative to work.
- (oracle)   \e[32m    --apcu-autoloader\e[39m                                Use APCu to cache found/not-found classes.
- (oracle)   \e[32m    --apcu-autoloader-prefix=APCU-AUTOLOADER-PREFIX\e[39m  Use a custom prefix for the APCu autoloader cache. Implicitly enables --apcu-autoloader
- (oracle)   \e[32m    --ignore-platform-req=IGNORE-PLATFORM-REQ\e[39m        Ignore a specific platform requirement (php & ext- packages).\e[33m (multiple values allowed)\e[39m
- (oracle)   \e[32m    --ignore-platform-reqs\e[39m                           Ignore all platform requirements (php & ext- packages).
- (oracle)   \e[32m-h, --help\e[39m                                           Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                                          Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                                        Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m                                 Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m                                 Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                                        Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m                                     Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m                                     Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m                        If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                                       Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m                                 Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
- (oracle) 
- (oracle) \e[33mHelp:\e[39m
- (oracle)   The \e[32minstall\e[39m command reads the composer.lock file from
- (oracle)   the current directory, processes it, and downloads and installs all the
- (oracle)   libraries and dependencies outlined in that file. If the file does not
- (oracle)   exist it will look for composer.json and do the same.
- (oracle)   
- (oracle)   \e[32mphp composer.phar install\e[39m
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#install--i
```

## misc-help-script

```
stdout:
+ (go)   args                         
- (oracle)   args                           
+ (go)       --dev                    Sets the dev mode.
+ (go)       --no-dev                 Disables the dev mode.
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
- (oracle)       --dev                      Sets the dev mode.
- (oracle)       --no-dev                   Disables the dev mode.
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
- (oracle)   The run-script command runs scripts defined in composer.json:
- (oracle)   
- (oracle)   php composer.phar run-script post-update-cmd
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#run-script--run
- (oracle) 
```

## misc-help-script-list

```
stdout:
+ (go)   args                         
- (oracle)   args                           
+ (go)       --dev                    Sets the dev mode.
+ (go)       --no-dev                 Disables the dev mode.
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
- (oracle)       --dev                      Sets the dev mode.
- (oracle)       --no-dev                   Disables the dev mode.
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
- (oracle)   The run-script command runs scripts defined in composer.json:
- (oracle)   
- (oracle)   php composer.phar run-script post-update-cmd
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#run-script--run
- (oracle) 
```

## misc-help-script-tty

```
stdout:
+ (go)   \e[32margs                        \e[39m 
- (oracle)   \e[32margs\e[39m                           
+ (go)   \e[32m    --dev                   \e[39m Sets the dev mode.
+ (go)   \e[32m    --no-dev                \e[39m Disables the dev mode.
+ (go)   \e[32m-h, --help                  \e[39m Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
+ (go)   \e[32m-q, --quiet                 \e[39m Do not output any message
+ (go)   \e[32m-v|-vv|-vvv, --verbose      \e[39m Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   \e[32m-V, --version               \e[39m Display this application version
+ (go)   \e[32m    --ansi                  \e[39m Force (or disable --no-ansi) ANSI output
+ (go)   \e[32m-n, --no-interaction        \e[39m Do not ask any interactive question
+ (go)   \e[32m    --profile               \e[39m Display timing and memory usage information
+ (go)   \e[32m    --no-plugins            \e[39m Whether to disable plugins.
+ (go)   \e[32m    --no-scripts            \e[39m Skips the execution of all scripts defined in composer.json file.
+ (go)   \e[32m-d, --working-dir=WORKING-DIR\e[39m If specified, use the given directory as working directory.
+ (go)   \e[32m    --no-cache              \e[39m Prevent use of the cache
- (oracle)   \e[32m    --dev\e[39m                      Sets the dev mode.
- (oracle)   \e[32m    --no-dev\e[39m                   Disables the dev mode.
- (oracle)   \e[32m-h, --help\e[39m                     Display help for the given command. When no command is given display help for the \e[32mlist\e[39m command
- (oracle)   \e[32m-q, --quiet\e[39m                    Do not output any message
- (oracle)   \e[32m-V, --version\e[39m                  Display this application version
- (oracle)   \e[32m    --ansi|--no-ansi\e[39m           Force (or disable --no-ansi) ANSI output
- (oracle)   \e[32m-n, --no-interaction\e[39m           Do not ask any interactive question
- (oracle)   \e[32m    --profile\e[39m                  Display timing and memory usage information
- (oracle)   \e[32m    --no-plugins\e[39m               Whether to disable plugins.
- (oracle)   \e[32m    --no-scripts\e[39m               Skips the execution of all scripts defined in composer.json file.
- (oracle)   \e[32m-d, --working-dir=WORKING-DIR\e[39m  If specified, use the given directory as working directory.
- (oracle)   \e[32m    --no-cache\e[39m                 Prevent use of the cache
- (oracle)   \e[32m-v|vv|vvv, --verbose\e[39m           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
- (oracle) 
- (oracle) \e[33mHelp:\e[39m
- (oracle)   The \e[32mrun-script\e[39m command runs scripts defined in composer.json:
- (oracle)   
- (oracle)   \e[32mphp composer.phar run-script post-update-cmd\e[39m
- (oracle)   
- (oracle)   Read more at https://getcomposer.org/doc/03-cli.md#run-script--run
```

## misc-list-scripts

```
stdout:
+ (go) gocomposer version 0.1.0-dev (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)
- (oracle)    ______
- (oracle)   / ____/___  ____ ___  ____  ____  ________  _____
- (oracle)  / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
- (oracle) / /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
- (oracle) \____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
- (oracle)                     /_/
- (oracle) Composer version 2.11.999-dev+source @release_date@
+ (go)   -h, --help                Display help for the given command. When no command is given display help for the list command
+ (go)   -q, --quiet               Do not output any message
+ (go)   -v|-vv|-vvv, --verbose    Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug
+ (go)   -V, --version             Display this application version
+ (go)       --ansi                Force (or disable --no-ansi) ANSI output
+ (go)   -n, --no-interaction      Do not ask any interactive question
+ (go)       --profile             Display timing and memory usage information
+ (go)       --no-plugins          Whether to disable plugins.
+ (go)       --no-scripts          Skips the execution of all scripts defined in composer.json file.
+ (go)   -d, --working-dir         If specified, use the given directory as working directory.
+ (go)       --no-cache            Prevent use of the cache
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
+ (go)   browse               Opens the package's repository URL or homepage in your browser
- (oracle)   browse               [home] Opens the package's repository URL or homepage in your browser
+ (go)   clear-cache          Clears composer's internal package cache
- (oracle)   clear-cache          [clearcache|cc] Clears composer's internal package cache
+ (go)   depends              Shows which packages cause the given package to be installed
- (oracle)   depends              [why] Shows which packages cause the given package to be installed
+ (go)   dump-autoload        Dumps the autoloader
- (oracle)   dump-autoload        [dumpautoload] Dumps the autoloader
+ (go)   install              Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
- (oracle)   install              [i] Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json
+ (go)   prohibits            Shows which packages prevent the given package from being installed
- (oracle)   prohibits            [why-not] Shows which packages prevent the given package from being installed
+ (go)   remove               Removes a package from the require or require-dev
+ (go)   repository           Manages repositories
+ (go)   require              Adds required packages to your composer.json and installs them
+ (go)   run-script           Runs the scripts defined in composer.json
- (oracle)   remove               [rm|uninstall] Removes a package from the require or require-dev
- (oracle)   repository           [repo] Manages repositories
- (oracle)   require              [r] Adds required packages to your composer.json and installs them
- (oracle)   run-script           [run] Runs the scripts defined in composer.json
+ (go)   self-update          Updates composer.phar to the latest version
+ (go)   show                 Shows information about packages
- (oracle)   self-update          [selfupdate] Updates composer.phar to the latest version
- (oracle)   show                 [info] Shows information about packages
+ (go)   update               Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
- (oracle)   update               [u|upgrade] Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file
```
