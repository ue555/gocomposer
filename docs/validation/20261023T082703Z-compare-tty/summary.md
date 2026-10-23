# Composer vs gocomposer comparison

- gocomposer_sha256: `6a0ed65979fc0c7c402072650995744d4e6933013828d30f864ff66b40bee12b`
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
| tty-init-full | MATCH | 0 | 0 |
| tty-init-defaults | MATCH | 0 | 0 |
| tty-init-options-as-defaults | MATCH | 0 | 0 |
| tty-init-require-options | MATCH | 0 | 0 |
| tty-init-abort | MATCH | 1 | 1 |
| tty-init-invalid-name | MATCH | 1 | 1 |
| tty-init-eof-first | MATCH | 0 | 0 |
| tty-init-eof-search | MATCH | 1 | 1 |
| tty-init-bad-license | MATCH | 1 | 1 |
| tty-init-invalid-author | MATCH | 0 | 0 |
| tty-init-pick-attempts | MATCH | 1 | 1 |
| tty-init-pick-empty | MATCH | 137 | 137 |
| tty-init-git | MATCH | 0 | 0 |
| tty-init-git-project-no | MATCH | 0 | 0 |
| tty-update-interactive-one | MATCH | 0 | 0 |
| tty-update-interactive-two | MATCH | 0 | 0 |
| tty-update-interactive-tab | MATCH | 0 | 0 |
| tty-update-interactive-arrows | MATCH | 0 | 0 |
| tty-update-interactive-backspace | DIFF | 1 | 1 |
| tty-update-interactive-invalid | DIFF | 1 | 1 |
| tty-update-interactive-empty | DIFF | 1 | 1 |
| tty-update-interactive-abort | MATCH | 1 | 1 |
| tty-update-interactive-filter | MATCH | 0 | 0 |
| tty-update-interactive-none | MATCH | 1 | 1 |
| tty-update-interactive-no-lock | MATCH | 3 | 3 |
| tty-update-interactive-n | MATCH | 1 | 1 |
| tty-require-dev-recommend-yes | MATCH | 0 | 0 |
| tty-require-dev-recommend-no | MATCH | 0 | 0 |
| tty-require-dev-recommend-two | MATCH | 0 | 0 |
| tty-require-dev-recommend-mixed | MATCH | 0 | 0 |
| tty-require-dev-recommend-update | MATCH | 0 | 0 |
| tty-require-move-default | MATCH | 0 | 0 |
| tty-require-move-yes | MATCH | 0 | 0 |
| tty-require-move-no-no | MATCH | 0 | 0 |
| tty-require-move-dev | MATCH | 0 | 0 |
| tty-require-search | MATCH | 137 | 137 |
| tty-require-search-none | MATCH | 0 | 0 |
| tty-require-search-eof | MATCH | 1 | 1 |
| tty-require-virtual | MATCH | 0 | 0 |
| tty-require-virtual-default | MATCH | 0 | 0 |
| tty-require-similar | MATCH | 0 | 0 |
| tty-require-similar-index | MATCH | 0 | 0 |
| tty-require-similar-empty | DIFF | 1 | 1 |
| tty-run-script-select | MATCH | 0 | 0 |
| tty-run-script-select-tab | MATCH | 0 | 0 |
| tty-run-script-select-invalid | DIFF | 1 | 1 |
| tty-run-script-select-empty | DIFF | 1 | 1 |
| tty-run-script-arg | MATCH | 0 | 0 |
| tty-run-script-no-scripts | MATCH | 1 | 1 |
| tty-exec-select | MATCH | 0 | 0 |
| tty-exec-select-name | MATCH | 0 | 0 |
| tty-exec-select-invalid | DIFF | 1 | 1 |
| tty-exec-no-bins | MATCH | 1 | 1 |
| tty-exec-list | MATCH | 0 | 0 |
| tty-create-project-ask | MATCH | 0 | 0 |
| tty-create-project-ask-default | MATCH | 0 | 0 |
| tty-create-project-ask-eof | MATCH | 1 | 1 |
| tty-create-project-ask-no-package | MATCH | 1 | 1 |
| tty-create-project-vcs-yes | MATCH | 0 | 0 |
| tty-create-project-vcs-default | MATCH | 0 | 0 |
| tty-create-project-vcs-no | MATCH | 0 | 0 |
| tty-create-project-vcs-remove-flag | MATCH | 0 | 0 |
| tty-plugin-allow-yes | MATCH | 0 | 0 |
| tty-plugin-allow-no | MATCH | 0 | 0 |
| tty-plugin-allow-discard | MATCH | 0 | 0 |
| tty-plugin-allow-help | MATCH | 0 | 0 |
| tty-plugin-allow-eof | MATCH | 1 | 1 |
| tty-plugin-allow-existing | MATCH | 0 | 0 |
| tty-parent-dir-yes | MATCH | 0 | 0 |
| tty-parent-dir-no | MATCH | 3 | 3 |
| tty-parent-dir-noninteractive | MATCH | 3 | 3 |
| tty-parent-dir-init-exempt | MATCH | 0 | 0 |

## tty-update-interactive-backspace

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-update-interactive-invalid

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-update-interactive-empty

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-require-similar-empty

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-run-script-select-invalid

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-run-script-select-empty

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```

## tty-exec-select-invalid

```
stdout:
+ (go) \e[33mIn <location>:\e[39m
```
