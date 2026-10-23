# 実装計画と受入範囲

これはComposer全仕様の完成を意味しない初期互換リリースの条件。順番に進めつつ、後段の実プロジェクトで必要になった前段機能を追加する。段階間の承認待ちは不要。

1. **基準と骨格**: Go module、cmd/gocomposer、CLI出力・終了コード、composer.json読取/validate、固定SHAの公式比較環境、ローカルfixture repositoryを作る。
2. **モデル・バージョン・解決器**: constraints/安定度/dev branch/alias、require/require-dev/conflict/provide/replace、platform/config.platform、優先度と依存解決、解決不能診断を実装。solverの代表的な成功/競合/循環ケースを移植する。単純な貪欲法だけで完了にしない。
3. **install/update/lock**: lockありの再現install、lockなしの解決、no-dev/dry-run/prefer-dist、限定update、content-hash、dist/source取得、cache、archive検証・展開、path repository、削除と更新、失敗時の整合性を実装。二度目のinstallの冪等性を検証する。
4. **autoloadと操作CLI**: PSR-4/PSR-0/classmap/files、autoload-dev、最適化、vendor/bin、installed.json/InstalledVersions互換を実装。require/remove/show/dump-autoload/check-platform-reqsと必要なconfig操作を接続する。autoloadは公式PHPで実際に読んで確認する。
5. **実プロジェクト**: 固定した小規模ライブラリ、CLIアプリ、Laravelアプリの3ケースをプロジェクト内sandboxに配置し、公式ComposerとGoで独立install/updateする。Laravelのvendor/autoload.phpとartisan about/listを公式PHPで起動し、gophpが利用可能なら追加検証する。既存アプリ本体でinstallしない。必要な標準scripts/eventsは明示設定PHPで実行し、no-scriptsを守る。pluginは無視せず、初期範囲ではno-pluginsでの検証と未対応時の明示失敗を区別する。
6. **回帰・再構築**: go test ./...、go test -race ./...、go vet ./...、gofmt、build、差分検証、空cacheからの再構築を自動化する。MacとLinuxのCI定義を作り、未実行CIを成功扱いにしない。READMEにCLI、必要な外部実行環境、制限と実証範囲を書く。

## 完了判定
- 上記1〜6の初期受入範囲について実装と実比較を完了する。条件を満たせない項目は明示し、未完了を維持する。
- 公式比較はinstall/update/require/remove/dump-autoload/validate/show/check-platform-reqsを含み、lock、依存バージョン、失敗時の状態とautoload動作を観測する。
- 解決不能、platform不一致、破損archive、展開先逸脱、ダウンロード失敗、既存vendorを持つ失敗ケースを収録する。
- 公式環境の版/SHA、Goバイナリのhash、入力hash、コマンドとexit codeを一意な検証ディレクトリに残す。秘密値は記録しない。
- `scripts/compare.sh` と `scripts/integration.sh` を実装し、実処理が失敗したら非0を返す。準備段階では存在しないのでverifyは未完了を返す。

## 後続互換範囲
任意PHP plugin、global/self-update/create-project/audit等の未実装コマンド、全repository形式・認証方式、全Composer設定との完全互換は別途追跡する。固定アプリで必要になれば今回の範囲で補完する。部分対応を「Composer完全互換」と表記しない。
