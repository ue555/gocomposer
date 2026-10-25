# gocomposer

Go で実装した Composer 互換の PHP 依存管理 CLI です。依存解決（CDCL SAT ソルバー）、`composer.lock` の生成・読込、ダウンロード・キャッシュ・展開、autoload 生成、scripts 実行を Go 本体で行います（Composer を子プロセスとして呼び出すラッパーではありません）。

基準は Composer `f83588392df43d4fd027d84d26de459946ef4e80`（2.11 開発版）で、公式 Composer との差分比較で出力・終了コード・lock・vendor・autoload の一致を確認した範囲だけを「対応」としています。Composer 全仕様との完全互換ではありません。対応範囲と制限は [docs/compatibility.md](docs/compatibility.md)、検証状況は [docs/progress.md](docs/progress.md) を参照してください。

## ビルド

```bash
go build -o bin/gocomposer ./cmd/gocomposer
```

Go 1.23 以上。外部 Go モジュールには依存しません（標準ライブラリのみ）。

## 実行に必要なもの

- PHP（CLI）: `php` が PATH にあるか `GOCOMPOSER_PHP` で指定。platform 要件（`php`、`ext-*`、`lib-*`）の検出、`@php` scripts、PHP クラスの script コールバック、プラグインの実行、`check-platform-reqs` に使います。
- git: VCS リポジトリと source インストールに使います（Composer と同じく git コマンドを実行）。
- PHP が無い場合は php/ext を提供しない platform として扱います（プラグインの実行と PHP を使う scripts はエラーになります）。
- unzip/7z は不要（zip/tar は Go で展開）。
- 設定・キャッシュ・認証情報は Composer と同じ場所（`COMPOSER_HOME`、`COMPOSER_CACHE_DIR`、`auth.json`、`COMPOSER_AUTH`）を使います。

## 使い方

```bash
gocomposer install                 # composer.lock から導入（lock が無ければ解決して作成）
gocomposer update [pkg...]         # 依存を解決し lock を更新（--with, -w/-W, --dry-run など）
gocomposer require vendor/pkg      # composer.json に追加して更新（バージョン推定あり）
gocomposer remove vendor/pkg
gocomposer create-project vendor/pkg dir
gocomposer show [-l|-o|-a|-t] [pkg]
gocomposer outdated | audit | why vendor/pkg | why-not vendor/pkg version
gocomposer dump-autoload [-o] | validate | check-platform-reqs | licenses | suggests | fund
gocomposer run-script name | exec bin | reinstall pkg | status [-v]
gocomposer config [--global] key [value...] | repository ... | policy ...
gocomposer init | bump | archive | global <command>
gocomposer search terms | browse pkg | clear-cache | diagnose
```

Composer の全コマンドを受け付け、オプション・ヘルプも Composer と同じです（`gocomposer help <command>`）。VCS（git、GitHub/GitLab/Bitbucket）リポジトリと source インストール、Composer プラグイン（symfony/flex を含む）の実行にも対応します。プラグインはコードだけを PHP プロセスで実行し、Composer API はスタブ経由で Go 本体が処理します。入力が端末のときは Composer と同じ対話プロンプト（init、require の検索、update --interactive、認証情報の入力など）を出し、`-n` や非端末入力では既定値で進みます。

対応できない部分は、成功を装わず明示的なエラーで終了します。主なものは次のとおりです。

- `self-update`（gocomposer は phar ではないため）
- forgejo ドライバ、ext-rar による rar 展開
- Composer のサービス自体を差し替えるなど、スタブにない API を使うプラグイン

一覧は [docs/compatibility.md](docs/compatibility.md) にあります。

## 検証

```bash
bash scripts/verify.sh
```

`gofmt`・`go vet`・`go test`（upstream の installer fixture 227件×2ほか、公式 PHP で生成した各種コーパスとの照合を含む）・`go test -race`・ビルドに続けて、次の 2 つを実行します。

- `scripts/compare.sh`: 公式 PHP 8.4（`-image` で 8.3/8.5 も可）と固定 SHA の Composer を podman コンテナで動かし、loopback の fixture リポジトリ（composer、git bare リポジトリ、フォージ API の HTTPS スタブ）に対して 1926 シナリオ（疑似端末での対話・端末表示を含む）で gocomposer と出力・終了コード・lock・プロジェクトツリー・autoload を比較します（`--network=none`）。最新結果は PHP 8.3/8.4/8.5 とも MATCH 1916・KNOWN-DIFF 10（理由は compatibility.md）・DIFF 0 です。
- `scripts/integration.sh`: symfony/console（ライブラリ）、CLI アプリ、laravel/laravel、小規模なライブダウンロード、主要プラグイン（composer/installers、php-http/discovery、dealerdirect/phpcodesniffer-composer-installer、pestphp/pest-plugin）、symfony/skeleton（symfony/flex の recipe 適用まで）を使うプロジェクトを公式 Composer と gocomposer で独立に導入し、vendor・autoload・`artisan about/list` を比較します（計 9 プロジェクト、ネットワーク必須）。

`scripts/ci-local.sh` は `.github/workflows/ci.yml` の全ステップを macOS（Go 1.23）と podman の Linux コンテナ（golang:1.23）でローカル実行します。GitHub Actions では、PHP 版ごとの matrix に分ける前のワークフローが run 39983736909（commit 479e793）で全ジョブ成功しています（分割後は未実行）。

前提の oracle 環境（公式 PHP イメージ、Composer の独立コピーと依存）は `scripts/prepare-oracle.sh` で作成します。結果は `docs/validation/<UTC時刻>-{compare,integration}/` に生ログ付きで保存されます。

## ライセンスと由来

Composer（MIT）と、その依存（composer/semver、composer/class-map-generator、composer/spdx-licenses、composer/pcre、seld/jsonlint、justinrainbow/json-schema、symfony/console）の挙動・データ・テスト fixture を移植しています。プラグイン実行用に、Composer の API クラスと Composer が同梱する PHP ライブラリを `internal/res/pluginrt` に埋め込んでいます（各 LICENSE 同梱）。詳細は [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) と `licenses/` を参照してください。

## 開発環境メモ

開発は Claude Code で行っています（`.claude/skills/gocomposer-development`、`prompts/`）。参照元 `/Users/kouji/oss/composer` は読み取り専用として扱い、変更しません。
