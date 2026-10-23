# Composer vs gocomposer comparison

- gocomposer_sha256: `c0002892352c26301aacdce924ba7a47d2f220270ad588f412c73edab23450f8`
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
| tty-auth-basic-store-yes | MATCH | 2 | 2 |
| tty-auth-basic-store-no | MATCH | 2 | 2 |
| tty-auth-basic-store-default | MATCH | 2 | 2 |
| tty-auth-basic-store-invalid | MATCH | 2 | 2 |
| tty-auth-basic-store-true | MATCH | 2 | 2 |
| tty-auth-basic-store-false | MATCH | 2 | 2 |
| tty-auth-basic-403 | MATCH | 2 | 2 |
| tty-auth-basic-wrong | MATCH | 100 | 100 |
| tty-auth-basic-empty | MATCH | 100 | 100 |
| tty-auth-basic-eof | MATCH | 1 | 1 |
| tty-auth-basic-eof-password | MATCH | 1 | 1 |
| tty-auth-basic-vv | MATCH | 2 | 2 |
| tty-auth-basic-local-authjson | MATCH | 2 | 2 |
| tty-auth-basic-no-interaction | MATCH | 100 | 100 |
| tty-auth-github-token | MATCH | 0 | 0 |
| tty-auth-github-token-vv | MATCH | 0 | 0 |
| tty-auth-github-empty-token | MATCH | 100 | 100 |
| tty-auth-github-invalid-token | MATCH | 100 | 100 |
| tty-auth-github-eof | MATCH | 1 | 1 |
| tty-auth-github-local-yes | MATCH | 0 | 0 |
| tty-auth-github-local-no | MATCH | 0 | 0 |
| tty-auth-github-ratelimit | MATCH | 100 | 100 |
| tty-auth-github-no-interaction | MATCH | 1 | 1 |
| tty-auth-gitlab-password | MATCH | 0 | 0 |
| tty-auth-gitlab-bad-then-good | MATCH | 0 | 0 |
| tty-auth-gitlab-local-yes | MATCH | 0 | 0 |
| tty-auth-gitlab-eof | MATCH | 1 | 1 |
| tty-auth-gitlab-no-interaction | MATCH | 1 | 1 |
| tty-auth-gitlab-refresh | MATCH | 0 | 0 |
| tty-auth-gitlab-refresh-vv | MATCH | 0 | 0 |
| tty-auth-gitlab-refresh-missing | MATCH | 1 | 1 |
| tty-auth-gitlab-refresh-rejected | MATCH | 1 | 1 |
| tty-auth-bitbucket-consumer | MATCH | 0 | 0 |
| tty-auth-bitbucket-consumer-vv | MATCH | 0 | 0 |
| tty-auth-bitbucket-invalid | MATCH | 100 | 100 |
| tty-auth-bitbucket-no-key | MATCH | 100 | 100 |
| tty-auth-bitbucket-no-secret | MATCH | 100 | 100 |
| tty-auth-git-store-yes | MATCH | 0 | 0 |
| tty-auth-git-store-no | MATCH | 0 | 0 |
| tty-auth-git-store-true-v | MATCH | 0 | 0 |
| tty-auth-git-wrong | MATCH | 1 | 1 |
| tty-auth-git-default-user | MATCH | 1 | 1 |
| tty-auth-git-eof | MATCH | 1 | 1 |
| tty-auth-git-no-interaction | MATCH | 1 | 1 |
| tty-auth-git-github | MATCH | 0 | 0 |
| tty-auth-git-github-no-token | MATCH | 1 | 1 |
| tty-auth-git-gitlab | MATCH | 0 | 0 |
| tty-auth-git-bitbucket | MATCH | 0 | 0 |
