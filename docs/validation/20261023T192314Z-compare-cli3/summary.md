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
| completion-bash | MATCH | 0 | 0 |
| completion-shell-env | MATCH | 0 | 0 |
| completion-no-shell | MATCH | 2 | 2 |
| completion-zsh | MATCH | 2 | 2 |
| completion-fish-env | MATCH | 2 | 2 |
| completion-quiet | MATCH | 0 | 0 |
| complete-commands | MATCH | 0 | 0 |
| complete-commands-prefix | MATCH | 0 | 0 |
| complete-commands-abbrev | MATCH | 0 | 0 |
| complete-commands-alias | MATCH | 0 | 0 |
| complete-commands-exact | MATCH | 0 | 0 |
| complete-app-options | MATCH | 0 | 0 |
| complete-command-options | MATCH | 0 | 0 |
| complete-command-options-short | MATCH | 0 | 0 |
| complete-working-dir-first | MATCH | 0 | 0 |
| complete-verbose-first | MATCH | 0 | 0 |
| complete-after-separator | MATCH | 0 | 0 |
| complete-none | MATCH | 0 | 0 |
| complete-format | MATCH | 0 | 0 |
| complete-format-eq | MATCH | 0 | 0 |
| complete-format-short | MATCH | 0 | 0 |
| complete-prefer-install | MATCH | 0 | 0 |
| complete-audit-format | MATCH | 0 | 0 |
| complete-bump-after-update | MATCH | 0 | 0 |
| complete-audit-options | MATCH | 0 | 0 |
| complete-audit-abandoned | MATCH | 0 | 0 |
| complete-licenses-format | MATCH | 0 | 0 |
| complete-fund-format | MATCH | 0 | 0 |
| complete-reinstall-type | MATCH | 0 | 0 |
| complete-show-ignore | MATCH | 0 | 0 |
| complete-init-require | MATCH | 0 | 0 |
| complete-help-format | MATCH | 0 | 0 |
| complete-list-format | MATCH | 0 | 0 |
| complete-global-option-value | MATCH | 0 | 0 |
| complete-update-packages | MATCH | 0 | 0 |
| complete-update-packages-prefix | MATCH | 0 | 0 |
| complete-remove-packages | MATCH | 0 | 0 |
| complete-bump-packages | MATCH | 0 | 0 |
| complete-depends-package | MATCH | 0 | 0 |
| complete-depends-platform | MATCH | 0 | 0 |
| complete-prohibits-package | MATCH | 0 | 0 |
| complete-require-vendors | MATCH | 0 | 0 |
| complete-require-vendor | MATCH | 0 | 0 |
| complete-require-names | MATCH | 0 | 0 |
| complete-require-platform | MATCH | 0 | 0 |
| complete-require-php | MATCH | 0 | 0 |
| complete-create-project-package | MATCH | 0 | 0 |
| complete-archive-package | MATCH | 0 | 0 |
| complete-show-installed | MATCH | 0 | 0 |
| complete-show-platform | MATCH | 0 | 0 |
| complete-show-available | MATCH | 0 | 0 |
| complete-outdated-package | MATCH | 0 | 0 |
| complete-suggests-packages | MATCH | 0 | 0 |
| complete-browse-packages | MATCH | 0 | 0 |
| complete-reinstall-packages | MATCH | 0 | 0 |
| complete-run-script | MATCH | 0 | 0 |
| complete-run-script-alias | MATCH | 0 | 0 |
| complete-exec-binary | MATCH | 0 | 0 |
| complete-help-command | MATCH | 0 | 0 |
| complete-list-namespace | MATCH | 0 | 0 |
| complete-completion-shell | MATCH | 0 | 0 |
| complete-global-command | MATCH | 0 | 0 |
| complete-global-options | MATCH | 0 | 0 |
| complete-repository-action | MATCH | 0 | 0 |
| complete-repository-enable | MATCH | 0 | 0 |
| complete-repository-add-type | MATCH | 0 | 0 |
| complete-policy-action | MATCH | 0 | 0 |
| complete-policy-arg1 | MATCH | 0 | 0 |
| complete-script-command | MATCH | 0 | 0 |
| complete-config-keys | MATCH | 0 | 0 |
| complete-config-keys-prefix | MATCH | 0 | 0 |
| complete-config-keys-extra | MATCH | 0 | 0 |
| complete-config-keys-unset | MATCH | 0 | 0 |
| complete-config-keys-global | MATCH | 0 | 0 |
| complete-config-keys-list | MATCH | 0 | 0 |
| complete-no-shell | MATCH | 0 | 0 |
| complete-bad-shell | MATCH | 2 | 2 |
| complete-no-current | MATCH | 2 | 2 |
| complete-bad-current | MATCH | 2 | 2 |
| complete-index-out-of-range | MATCH | 2 | 2 |
| complete-missing-input-value | MATCH | 1 | 1 |
| complete-debug-error | DIFF | 1 | 1 |
| complete-repository-names | MATCH | 0 | 0 |
| complete-repository-before | MATCH | 0 | 0 |
| complete-policy-names | MATCH | 0 | 0 |
| complete-unlocked-installed | MATCH | 0 | 0 |
| complete-no-json-command | MATCH | 0 | 0 |
| complete-no-json-update | MATCH | 1 | 1 |
| complete-no-json-config | MATCH | 0 | 0 |
| complete-global-require | MATCH | 0 | 0 |
| complete-global-show | MATCH | 0 | 0 |
| complete-global-no-home-json | MATCH | 1 | 1 |
| complete-plugin-commands | DIFF | 0 | 0 |
| complete-plugin-namespace | DIFF | 0 | 0 |
| complete-plugin-options | DIFF | 0 | 0 |
| complete-plugin-no-plugins | MATCH | 0 | 0 |

## complete-debug-error

```
stderr:
- (oracle) Running 2.11.999-dev+source (@release_date@) with PHP 8.4.27 on Linux / 7.1.10-200.fc44.aarch64
- (oracle)   [LogicException]                                                             
- (oracle) Exception trace:
- (oracle)   at /oracle-source/vendor/symfony/console/Completion/CompletionInput.php:224
- (oracle)  Symfony\Component\Console\Completion\CompletionInput->isCursorFree() at /oracle-source/vendor/symfony/console/Completion/CompletionInput.php:214
- (oracle)  Symfony\Component\Console\Completion\CompletionInput->getRelevantToken() at /oracle-source/vendor/symfony/console/Completion/CompletionInput.php:74
- (oracle)  Symfony\Component\Console\Completion\CompletionInput->bind() at /oracle-source/vendor/symfony/console/Command/CompleteCommand.php:174
- (oracle)  Symfony\Component\Console\Command\CompleteCommand->createCompletionInput() at /oracle-source/vendor/symfony/console/Command/CompleteCommand.php:87
- (oracle)  Symfony\Component\Console\Command\CompleteCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## complete-plugin-commands

```
stderr:
- (oracle) Plugin command show (Plugfix\Commands\ShowCommand) would override a Composer command and has been skipped
```

## complete-plugin-namespace

```
stderr:
- (oracle) Plugin command show (Plugfix\Commands\ShowCommand) would override a Composer command and has been skipped
```

## complete-plugin-options

```
stderr:
- (oracle) Plugin command show (Plugfix\Commands\ShowCommand) would override a Composer command and has been skipped
```
