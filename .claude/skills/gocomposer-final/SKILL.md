---
name: gocomposer-final
description: gocomposerの第3段階（文書の更新、completion、svn/hg/fossil/perforce、v1 providers-url、rar/xz/7z、ResponseWarnings、端末表示、プラグインのサービス差し替え、root実行の確認など残る全項目）を確認なしで一括実装・公式比較・記録する。prompts/final.md から使う。
---

# 第3段階: 残る全項目

`AGENTS.md`、`CLAUDE.md`、`gocomposer-development` スキル、`docs/progress.md`、`docs/compatibility.md` を先に読み、その原則（Go本体で処理、公式PHP/Composerとの差分比較、未対応は明示失敗、参照元は読取専用、commit/push/公開なし、確認フローなし）をすべて守る。
前段の `gocomposer-remaining`・`gocomposer-next` は完了済み。その成果（PHP 8.3/8.4/8.5 とも 1386 中 MATCH 1377・KNOWN-DIFF 9・DIFF 0、installer fixture 454件、実プロジェクト9件、ローカル CI、GitHub Actions run 39983736909 の全ジョブ成功）を退行させない。

## 対象外（変えない）
設計上の差として既に記録済みのもの: self-update（phar ではない）、Composer 自体が PHP エラー/LogicException で落ちる入力（validate-lax-fail・composer-auth-invalid・bump-lock-disabled）、diagnose の Composer 自身の監査行。これらの KNOWN-DIFF は変えない。

## 対象（この順で進める）
1. **文書の更新**: GitHub Actions run 39983736909（commit 479e793、linux/amd64、1時間52分）の全ジョブ成功を `docs/progress.md`・`docs/compatibility.md`・`README.md` に反映し、「GitHub Actions は未実行」の記述をなくす。compatibility.md の「対話」節（OSC 9;4 移植済み）と「端末表示」節（OSC 9;4 未出力）など、文書内の矛盾を解消する。
2. **completion**: Symfony Console の `completion`（bash/zsh/fish のスクリプト出力）と隠しコマンド `_complete`、Composer の各コマンドの補完候補（パッケージ名、config のキー名、スクリプト名など）を移植する。`config` のキー名補完はここで対応する。
3. **細部の移植**: `config` の `-vvv` 読込ログと bin-compat symlink の非推奨警告、Factory の `-vvv` デバッグログ、validate と reinstall での command イベント、`run-script` が Symfony Command クラスを指す場合の実行、composer.json なしでの fund/search/browse のパッケージ指定（packagist.org を loopback スタブに向けて比較）。
4. **VCS: svn / hg / fossil / perforce**: ドライバ（SvnDriver、HgDriver、FossilDriver、PerforceDriver）、ダウンローダ、`status`、source インストールを移植する。oracle イメージに必要なクライアント（subversion、mercurial、fossil）を `oracle/Containerfile` に追加し、ローカルのリポジトリ（`file://` や loopback サーバー、`--network=none`）で比較する。Perforce のサーバー/クライアントが自由に入手・実行できない場合は、Go 実装を単体テストで検証し、公式比較できない理由を記録する（成功扱いにしない）。
5. **v1 providers-url リポジトリ**: Composer の ComposerRepository の providers-url / provider-includes 処理を移植し、loopback の fixture で比較する。
6. **アーカイブ: rar / xz / 7z**: Composer の RarDownloader（ext-rar）、XzDownloader（xz コマンド）、ZipDownloader の 7z フォールバックを調べ、Composer と同じ条件で対応する。外部 Go モジュールは使わない（必要なら標準ライブラリだけで展開処理を書く）。oracle 側に拡張やコマンドが必要なら Containerfile の別バリアントで用意する。実装できない形式は明示エラーと理由を残す。
7. **ResponseWarnings**: filter api-url と advisory API の応答に含まれる warning/info の表示を移植し、比較する。
8. **端末表示**: decorated 出力のハイパーリンク（OSC 8）、ProgressBar の上書き描画、curl のエラー番号（6/7 以外: 28 タイムアウト、35/60 TLS など、Go のエラーから対応づけられるもの）を Composer に合わせる。pty ハーネス（`tools/compare/ptyrun`）で比較する。
9. **プラグイン**: getArchiveManager、Composer の set* によるサービス差し替え、setRepositoryClass・プラグイン定義のリポジトリクラス、組込みダウンローダの置換、AutoloadGenerator の buildPackageMap/parseAutoloads/createLoader 等を、PHP 側の実装を呼び出す方式で可能な範囲まで対応する。並列ダウンロード時の post-file-download の発火順は、Composer の順序（完了順）が決定的になる条件で比較する。対応できないものは明示エラーのまま理由を記録する。
10. **root 実行時の確認プロンプト**: 比較ハーネスに root で実行するシナリオ（`podman run --user 0`、Linux では userns 設定との両立に注意）を追加し、公式と比較する。
11. **統合テストの改善**: cli-app と laravel の `update --no-install` は post-update-cmd が vendor なしで両者とも 255 になるだけで更新の検証になっていない。`--no-scripts` 付きのステップを追加する（既存ステップは残す）。
12. **CI の所要時間**: 項目追加で GitHub Actions の oracle ジョブが 6 時間の上限に近づかないよう、`.github/workflows/ci.yml` の oracle ジョブを PHP バージョンの matrix に分ける（各ジョブで prepare）。`scripts/ci-local.sh` も対応させる。GitHub 上の実行は push が必要で範囲外なので、ローカル再現のみと記録する。

## 作業単位と検証
各項目を「Go実装 → 単体テスト → `testdata/compare` シナリオ追加 → `scripts/compare.sh` → 差異修正 → docs 更新」の単位で完了させてから次へ進む。Containerfile を変えたら `scripts/prepare-oracle.sh` で 3 版のイメージを作り直し、digest を記録し直す。数項目ごとに `scripts/verify.sh` 全体（PHP 8.3/8.5 の compare を含む）を通す。最後に `scripts/verify.sh` と `scripts/ci-local.sh` を通す。
Mac（podman machine）と Linux（GitHub Actions の rootless podman）の両方で動くよう、podman の起動は `tools/*/userns.go` の `podmanRun` を通す。

## 状態ファイル
セッションを終える直前に必ず `.cache/final/status` を書く。1行目は次のいずれか、2行目以降に理由と再開地点。
- `DONE` — 1〜12 がすべて実装・比較・記録済み（公式比較できない理由が記録された項目を含む）で、最後の `scripts/verify.sh` と `scripts/ci-local.sh` が exit 0。
- `CONTINUE` — まだ進められる作業が残っている。
- `BLOCKED` — 認証・利用上限・外部障害など、自力で解消できない理由だけが残っている。
受入範囲を縮めて `DONE` にしない。
