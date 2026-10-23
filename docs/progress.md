# 進捗

最終更新: 2026-10-23。基準: Composer `f83588392df43d4fd027d84d26de459946ef4e80`（参照元と一致、参照元はgit status clean・未変更）。

## 比較環境（oracle）

- 公式PHP: `docker.io/library/php:8.4-cli` digest `sha256:05c3a54f57053f6d9bb2a32e32d6ddbd7f5a443672385de37a2a8736a6ebbfd5`（PHP 8.4.27）を主比較。回帰用に php:8.3-cli（`sha256:f1b772b8…fb47`、PHP 8.3.36）と php:8.5-cli（`sha256:361f9916…09fd`、PHP 8.5.12）。`oracle/Containerfile` で git/unzip を追加した `localhost/gocomposer-oracle:{8.3,8.4,8.5}` と、ext-zip を追加した `-zip` 版（artifact リポジトリ用）。digest の一覧は compatibility.md の「検証環境」。再作成は `scripts/prepare-oracle.sh`。
- Composer: 参照元とは別の `.cache/oracle-source`（基準SHAのclone）。依存は同梱 `composer.lock`（sha256 `7e7ed6af…34c`）から composer.phar 2.10.3（sha256 `7a2d379d5b8ffdaa028580ef26494c36d2feef4b178d3dd1473a4dbc5e17c8d6`、公式sha256sumと一致）で `install --no-dev`。
- 参照元 `/Users/kouji/oss/composer` には書き込んでいない。

## 検証コマンドと最新結果

| 検証 | コマンド | 最新結果 |
|---|---|---|
| 公式差分比較 PHP 8.4 | `scripts/compare.sh` | 1386シナリオ中 MATCH 1377・KNOWN-DIFF 9・DIFF 0（20261023T093500Z-compare、最終 verify.sh 内） |
| 公式差分比較 PHP 8.3 | `scripts/compare.sh -image localhost/gocomposer-oracle:8.3` | MATCH 1377・KNOWN-DIFF 9・DIFF 0（20261023T101132Z-compare） |
| 公式差分比較 PHP 8.5 | `scripts/compare.sh -image localhost/gocomposer-oracle:8.5` | MATCH 1377・KNOWN-DIFF 9・DIFF 0（20261023T104618Z-compare） |
| 固定実プロジェクト統合 | `scripts/integration.sh` | symfony/console v7.4.9、CLIアプリ、laravel/laravel v13.9.0、live-psr、プラグイン 4 件（composer/installers、php-http/discovery、dealerdirect/phpcodesniffer-composer-installer、pestphp/pest-plugin）、symfony/skeleton v7.4.99（recipes:install を含む）の 9 プロジェクトの全ステップが MATCH（20261023T100947Z-integration） |
| upstream installer fixture | `go test ./internal/composer -run TestUpstreamInstallerFixtures` | Composer `tests/Composer/Test/Fixtures/installer` の227件×(pool optimizer 0/1)=454件すべて通過（期待値は無修正。PHPのplatformはoracle PHP 8.4.27のスナップショット＋Composer本体のテスト環境と同じくext-intlを追加）。過程で update mirrors の参照保持（LockTransaction::updateMirrorAndUrls）を修正 |
| 全体ゲート | `scripts/verify.sh` | exit 0（gofmt・go test（installer fixture 454件含む）・go test -race・go vet・build・compare・integration）。ログ docs/validation/20261023T060217Z-verify/verify.log |
| CI ローカル実行 | `scripts/ci-local.sh` | .github/workflows/ci.yml の全ステップ PASS（最新 20261023T112150Z-ci-local、第2段階の 1386 シナリオで実行）: unit macOS（go1.23.12）・unit Linux（podman golang:1.23 digest sha256:ac2b6b75…）・prepare・compare 8.4/8.3/8.5・integration。GitHub Actions 上は push が必要なため範囲外で未実行 |

## 実装済み（Go、主なもの）

- phpjson（PHP互換JSON、float 3書式をPHPと照合）、jsonlint、jsonfile、jsonschema（justinrainbow 6.13.1、コーパス340件一致）、pcre（JsonManipulator用バックトラッキングエンジン、4447件一致）、jsonmanip（JsonManipulatorTest 132データセット全一致）、configsource（JsonConfigSourceTest全一致）。
- semver（composer/semver、strnatcmp含む）、spdx、validate（ValidatingArrayLoader/ConfigValidator/email/URL）。
- solver（PoolBuilder/Optimizer/CDCL/Policy/Problem/Transaction、policy除去理由の問題文）、policy（security advisory/malware filter list/abandoned/cooldown、Auditor 4書式、208テスト）。
- repo（composer v1/v2、path、package、filter、search、InstalledRepository）、platform、config、cache、httpdl（認証・再試行・POST・curlエラー番号6/7）、archive（zip/tar、逸脱防止）、InstallationManager（並列DL、ProgressBar、bin proxy）、Installer、Locker、autoload生成、scripts（PHPコールバック含む）。
- CLI: Composer の全コマンド（self-update は明示失敗）。コマンド定義・ヘルプ・synopsisは公式から生成。Symfony互換の例外表示・Table・ProgressBar・色出力。
- VCS（git/github/gitlab/bitbucket）・artifact リポジトリ、git source インストールと変更検出、Composer プラグイン実行（常駐 PHP プロセス＋API スタブ）、files cache の確率的 GC。

## 未対応項目の一括実装（2026-10-23、gocomposer-remaining）

並列エージェントで担当ファイルを分け、各自 `tools/compare/scenarios_<担当>.go`（`extraScenarios` に init() 登録）と `internal/cli/*.go`（`registerCommand` で登録）を所有した。仕様判断は基準SHAのソースと公式比較で行った。

| 項目 | 比較結果（PHP 8.4） | 判断・備考 |
|---|---|---|
| show --latest/--outdated/--all/--available、outdated | 担当時の報告は 127件 MATCH（20261022T183044Z-compare-show-latest 109件＋20261022T183325Z 18件）。最新の全体実行 20261023T060228Z では show-latest-* 92件・outdated-* 29件が MATCH（一部は show-* 名のシナリオ） | `--ansi` 一致のため cio の色出力を Symfony の区間別開始/解除コードに修正（全コマンドに影響、既存シナリオ退行なし）。cooldown 残り時間は時の境界をまたぐと揺れうる |
| audit | audit-* 68件 MATCH（20261022T182049Z-compare-audit） | 基準SHAの Auditor は 0/1 のみ返すので終了コードはビットフラグでなく 0/1。api-url 形式用 fixture /v2api を追加。curl error 7 の文言を oracle の libcurl 8.14.1 に合わせ「Could not connect to server」に変更 |
| run-script/exec/clear-cache/fund/search/browse/diagnose/self-update/repository/policy | 224件中 MATCH 218・KNOWN-DIFF 6（20261022T184437Z-compare-misc、20261022T201210Z-compare-misc） | self-update は自己更新機構がないため明示失敗（KNOWN-DIFF 5）。diagnose の Composer 自身の監査行のみ差（KNOWN-DIFF 1）。clear-cache-gc-ttl は秒境界で結果が揺れるため撤去し、決定的な gc-ttl-fresh と単体テスト TestCacheGCTTL で確認。scripts に process-timeout を適用 |
| init/bump/archive/global | 112件中 MATCH 111・KNOWN-DIFF 1（20261022T184324Z-compare-cmds2） | init は非対話のみ（対話は明示エラー）。基準SHAの init は既存 composer.json を上書きする（比較で確認）。bump-after-update を update に実装。archive は oracle に ext-zip がないため PharData 形式を再現。bump-lock-disabled は Composer が LogicException で落ちる |
| 全体 | 20261022T201228Z-compare-misc-full: 960件中 MATCH 938・KNOWN-DIFF 9・DIFF 13 | DIFF 13 はすべて作業中の vcs-*/source-*。既存 372 件の退行なし |
| VCS(git)・source install | vcs-*/source-* 124件 MATCH（20261023T004222Z-compare-vcs2-full） | 最初の担当は利用上限後に停止したまま再開しなかったため停止し、別担当が仕上げた。GitHub/GitLab/Bitbucket は loopback HTTPS スタブ（使い捨て CA、cafile）で比較。AuthHelper の非対話 401/403 処理を移植。トークンが出力・キャッシュ・vendor に出ないことを -vvv でも確認。config シナリオの repositories を vcs 型に戻して一致（root version 警告の差は解消） |
| artifact リポジトリ | vcs-artifact-* 20件 MATCH | Composer の ArtifactRepository は ext-zip 必須のため ext-zip 付き oracle（oracle/Containerfile.zip、`<image>-zip`）で比較。8.3/8.5 実行時も対応する -zip 版を使うようハーネスを修正し、digest を summary に記録 |
| reinstall/status/create-project | 107件 MATCH（20261023T003756Z-compare-cmds3） | status の dist 比較は FileDownloader::getLocalChanges と同じく _compare に再取得。create-project はリポジトリ指定なしだと packagist.org 接続となり、ネットワークなしで Composer と同じ curl error 6・終了100 |
| plugin 実行 | plugin-* 58件 MATCH（20261023T011003Z-compare-plugin）、統合 8 プロジェクト MATCH（20261023T011117Z-integration） | プラグインのコードだけを常駐 PHP プロセスで実行し、Composer API はスタブ経由で Go に戻す。Composer の API クラスと同梱ライブラリを基準 SHA / oracle lock の版で internal/res/pluginrt に埋め込み（THIRD_PARTY_NOTICES に記載、Composer の LICENSE を同梱）。実在プラグイン 4 件（composer/installers、php-http/discovery、dealerdirect/phpcodesniffer-composer-installer、pestphp/pest-plugin）を固定 lock で統合比較 |
| files cache の確率的 GC | 単体テスト | Composer は FileDownloader 生成時に 1/51 の確率で cache GC を行い -vv で「Running cache garbage collection」を出す。gocomposer に未実装だったため移植した（internal/cache/gc.go）。GC はキャッシュディレクトリを作成・削除し比較を揺らすため、比較では両ツールに COMPOSER_TEST_SUITE=1 を渡して無効化する（当初の行除外は撤去） |
| 全体（第2波後） | 20261023T011112Z-compare-plugin-full: MATCH 1180・KNOWN-DIFF 9・DIFF 0 | |
| CI ローカル実行 | 1回目 20261023T023955Z、2回目 20261023T035247Z、3回目 20261023T050038Z（全 PASS） | 1・2回目は compare-8.4 で各1件の非決定的な差。reinstall-progress は Composer の非同期 rm -rf の競合で oracle の進捗バーが揺れる（過去12回中11回は不定形）ため、そのシナリオだけ「Removing」直後のバーを不定形に揃えて比較。clear-cache-empty は 1/51 の確率的キャッシュ GC がキャッシュディレクトリを作るためで、比較では両ツールに COMPOSER_TEST_SUITE=1（基準SHAではこの GC 無効化にのみ使用）を渡す。PHP 8.5 では Composer 自身の Svn.php が出す PHP Deprecation Notice を oracle 側 stderr から除外。各正規化は compatibility.md に記載し、生出力は結果 JSON に残る |

## 第2段階（2026-10-23、gocomposer-next）

| 項目 | 状態 | 備考 |
|---|---|---|
| 1. 文書の整理 | 完了 | compatibility.md「検証環境」を現状化（PHP 8.3/8.4/8.5 のベース・oracle・ext-zip イメージ digest、実プロジェクト 8 件、KNOWN-DIFF 9 件の内訳、CI はローカル再現のみ）。README の比較件数（374→1189）・実プロジェクト件数・PHP なし時の記述の置き場所、progress.md のキャッシュ GC の扱い（行除外→COMPOSER_TEST_SUITE）、config-* 件数（221→224、最新ログ）を実ログに合わせて修正 |
| 2. symfony/flex・プラグイン用スタブ API | 完了 | pre/post-file-download・pre-pool-create の発火と反映、DownloadManager/HttpDownloader/AutoloadGenerator/RepositoryManager/Factory::create/Config::merge/Installer::create 等のスタブ、download・installCode 上書き/LibraryInstaller 非継承インストーラ、symfony/flex v2.11.0（loopback HTTPS recipe スタブ）。plugin2-* 34・flex-* 22 シナリオ追加、PHP 8.3/8.4/8.5 で MATCH（20261023T085804Z-compare-flex-full ほか）。integration に symfony/skeleton v7.4.99（install/update-no-install/recipes-install）を追加し 9 件全 MATCH（20261023T085132Z-integration）。残る明示エラーは compatibility.md の plugin 節 |
| 3. 対話プロンプト（pty 比較） | 完了 | tools/compare に疑似端末実行（ptyrun、入力列、生バイト保存）を追加。internal/cio に QuestionHelper を移植。init/require/remove/update -i/run-script/exec/create-project/allow-plugins/親ディレクトリ確認（tty-* 81）、認証入力・discard-changes（tty-auth-* 48、tty-discard-* 12）を追加し全 MATCH（20261023T084145Z-compare-tty-full、20261023T085613Z-compare-tty2-full）。-n/非端末入力の既存シナリオは不変 |
| 最終検証 | 完了 | `scripts/verify.sh` exit 0（20261023T093500Z-verify、compare 8.4 20261023T093500Z、integration 20261023T100947Z）、compare 8.3 20261023T101132Z・8.5 20261023T104618Z（いずれも 1386 中 MATCH 1377・KNOWN-DIFF 9・DIFF 0）、`scripts/ci-local.sh` 全 PASS（20261023T112150Z-ci-local） |

## 未実装・未比較

- config: キー名補完・-vvv 読込ログ・bin-compat symlink 警告は未移植。
- 対話: root 実行時の確認はシナリオなし（コンテナが uid 1000）。discard 確認の EOF は Composer と同じく無限再質問のため比較なし。
- self-update は明示失敗、diagnose の Composer 自身の監査行のみ差（KNOWN-DIFF）。
- plugin: サービス差し替え（set*、getArchiveManager）、setRepositoryClass、組込みダウンローダの置換などは明示エラー（compatibility.md）。post-file-download の並列時の発火順は操作順（既知の差）。
- VCS: svn/hg/fossil/perforce ドライバ本体、v1 providers-url は明示エラー。
- CI: GitHub Actions 上は未実行（push が範囲外）。ローカル再現のみ。

## 次の作業

gocomposer-remaining の対象 1〜7 はすべて実装・比較・記録済み。残る差は docs/compatibility.md に明示エラーまたは KNOWN-DIFF として列挙した仕様差のみ。
gocomposer-next の対象 1〜3 も実装・比較済み（最終検証は第2段階の表を参照）。
今後の候補: svn/hg ドライバ本体、v1 providers-url、残るプラグイン API（compatibility.md の明示エラー一覧）、GitHub Actions 上での CI 実行（push 後）。
