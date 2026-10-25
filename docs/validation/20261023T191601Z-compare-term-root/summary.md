# Composer vs gocomposer comparison

- gocomposer_sha256: `d5225c3feebc6b2c3f6b9a4e3cfe120e7fc20ebdea8e178661047bde16d972ad`
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
| root-alias | MATCH | 0 | 0 |
| root-validate | MATCH | 0 | 0 |
| root-validate-allowed | MATCH | 0 | 0 |
| root-validate-allowed-zero | MATCH | 0 | 0 |
| root-validate-no-plugins | MATCH | 0 | 0 |
| root-update | MATCH | 0 | 0 |
| root-list | DIFF | 0 | 0 |
| root-list-allowed | DIFF | 0 | 0 |
| root-script-alias | MATCH | 0 | 0 |
| root-tty-validate-yes | MATCH | 0 | 0 |
| root-tty-validate-default | MATCH | 0 | 0 |
| root-tty-validate-no | MATCH | 0 | 0 |
| root-tty-validate-n | MATCH | 0 | 0 |
| root-tty-validate-allowed | MATCH | 0 | 0 |
| root-tty-list-yes | DIFF | 0 | 0 |
| root-tty-list-no | DIFF | 0 | 0 |
| root-tty-script-alias-yes | MATCH | 0 | 0 |
| root-tty-update-yes | MATCH | 0 | 0 |
| root-tty-no-plugins | MATCH | 0 | 0 |

## root-list

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

## root-list-allowed

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

## root-tty-list-yes

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
+ (go)   \e[32mhello               \e[39m Runs the hello script as defined in composer.json
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
- (oracle)   \e[32mhello\e[39m                Runs the hello script as defined in composer.json
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

## root-tty-list-no

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
+ (go)   \e[32mhello               \e[39m Runs the hello script as defined in composer.json
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
- (oracle)   \e[32mhello\e[39m                Runs the hello script as defined in composer.json
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
