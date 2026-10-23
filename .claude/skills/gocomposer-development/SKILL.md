---
name: gocomposer-development
description: Composerの固定ソースを基準にGo製gocomposerの設計、実装、差分検証と統合テストを連続して進める。gocomposerプロジェクト内の開発で使用する。
---

# Composer互換CLIの継続開発

プロジェクトの `CLAUDE.md` と `docs/implementation-plan.md` に従う。承認済みの実装と検証は、段階ごとの質問や「続けますか」を挟まず継続する。

## 基準を固定する
`docs/upstream.json` のSHAと参照元の現在値を照合する。差があれば黙って基準を更新せず、そのSHAをプロジェクト内の `.cache/oracle-source` に独立checkoutする。参照元に依存導入しない。
Composerの `src/Composer` と `tests/Composer/Test` を調べ、依存パッケージ composer/semver、class-map-generator などもlockに固定された実装を確認する。
通常のPHPとGo製PHPを混同しない。公式PHPの実行パスまたはコンテナdigest、ComposerのSHA、依存lockハッシュを比較結果に保存する。

## 開発ループ
未完了の依存関係から次の実装を選び、正常・失敗・副作用のfixtureを用意する。Go実装→単体/統合テスト→公式比較→差異修正→再検証→進捗記録を連続して行う。
解決器、ファイル操作、autoloadは観測結果で検証する。未実装処理をmockしたテストだけで受入を満たさない。
ローカルfixture repositoryと固定archiveを優先し、Packagistの変動で回帰が揺れないようにする。実ネットワーク検証は別に記録する。
公式比較が実行不能でも独立したGo実装を続ける。期待値は暫定とし、環境エラーを一致・不一致に数えない。

## 比較契約
コマンド引数・環境・入力ファイルを共通化し、双方を別のコピーで実行する。stdout/stderr/終了コード、解決package/version/reference、lockの意味、vendor内容、autoloadで読み込めるクラス・関数を比較する。
絶対一時パス、所要時間、進捗表示だけを必要に応じ正規化する。package順序、エラーの種類、scripts順序など意味のある差を消さない。正規化前の生ログも残す。
PHP 8.4を主比較環境、8.3/8.5を生成autoloadとplatform条件の回帰対象にする。利用可能な正確なpatch版を固定して記録し、未実行の版を成功としない。
Composer更新とGo実装の更新を混ぜない。upstreamのテスト期待値はGoの都合で書き換えない。

## 終了条件
実装計画の受入条件、`scripts/verify.sh`、公式差分比較、固定実プロジェクトの再導入が成功していること。残るComposer全仕様との隔たりを `docs/compatibility.md` に列挙する。
計画・スキルの自己修正で受入範囲を縮めて完了扱いにしない。成果を次へ引き継げるよう毎回 `docs/progress.md` を更新する。
