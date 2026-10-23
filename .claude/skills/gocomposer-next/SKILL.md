---
name: gocomposer-next
description: gocomposerの第2段階（文書整理、symfony/flex対応とプラグイン用スタブAPIの拡充、対話プロンプトの実装）を確認なしで一括実装・公式比較・記録する。prompts/next.md から使う。
---

# 第2段階: 文書整理・flex/スタブAPI・対話プロンプト

`AGENTS.md`、`CLAUDE.md`、`gocomposer-development` スキル、`docs/progress.md`、`docs/compatibility.md` を先に読み、その原則（Go本体で処理、公式PHP/Composerとの差分比較、未対応は明示失敗、参照元は読取専用、commit/push/公開なし、確認フローなし）をすべて守る。
前段の `gocomposer-remaining` は完了済み（`.cache/remaining/status` は DONE）。その成果（MATCH 1180・KNOWN-DIFF 9・DIFF 0、installer fixture 454件、実プロジェクト8件）を退行させない。

## 1. 文書の整理
- `docs/compatibility.md` の「検証環境」節を現状に合わせる（CI はローカル再現済み・GitHub 上は未実行、実プロジェクト8件の一覧、PHP イメージの digest）。
- `docs/compatibility.md`・`docs/progress.md`・`README.md` を突き合わせ、件数・状態・コマンド表の矛盾をなくす。実際のログ（`docs/validation/`）にない数字は書かない。

## 2. symfony/flex 対応とプラグイン用スタブ API の拡充
基準SHAの `src/Composer/Plugin`、`Downloader`、`Autoload`、`Repository`、`Factory.php`、`Config.php`、`Installer/*` を読み、既存のプラグイン実行（常駐PHPプロセス＋スタブ経由でGoに戻す方式、`internal/composer/plugin.go`、`internal/res/pluginrt`）を拡張する。解決・lock・取得・展開・autoload は引き続き Go が行う。
- **イベント**: `pre-file-download`（PreFileDownloadEvent: processedUrl の書き換え、customCacheKey、transportOptions、type=package/metadata）、`post-file-download`（PostFileDownloadEvent）、`pre-pool-create`（PrePoolCreateEvent: packages の差し替え・除外が解決に反映されること）を発火し、Composer と同じタイミング・回数・引数にする。
- **スタブ API**: `getDownloadManager`（getDownloader、setOutputProgress 等）、`getAutoloadGenerator`、リモートリポジトリ（RepositoryManager/ComposerRepository の findPackage(s)・search 等）、`Factory::create`、`Config::merge`、`download`/`installCode` を上書きするインストーラ、LibraryInstaller を継承しないインストーラ。プラグインからの呼び出しで実際に起きる副作用（ダウンロード・ファイル配置）は Go 側で行う。
- **symfony/flex**: `extra.symfony.require` による symfony/* のバージョン制限、recipe の取得と適用（`config/`、`.env`、`symfony.lock` などの生成）、`symfony:recipes` 系コマンド、unpack を公式と比較する。recipe エンドポイントは loopback の HTTPS スタブを用意して `--network=none` の compare シナリオに入れ、加えて `scripts/integration.sh` に symfony/skeleton（固定バージョン、lock は公式 Composer で生成して `testdata/integration/locks` に固定）を実ネットワークのプロジェクトとして追加する。
- fixture プラグインで各 API・各イベントの正常系と失敗系のシナリオを追加する。それでも対応できない API は明示エラーのまま残し、理由を compatibility.md に書く。

## 3. 対話プロンプトの実装
Composer は入力が端末のときだけ対話モードになるため、比較ハーネスに疑似端末（pty）での実行と、決めた入力列を流す仕組みを追加する（`tools/compare` を拡張。外部 Go モジュールは使わず、必要なら oracle コンテナ内の Linux pty を直接扱う）。stdout/stderr の制御文字は意味を消さない範囲で正規化し、生ログは残す。
対象（Composer の QuestionHelper/IO の挙動、既定値、再質問、検証エラー、Ctrl-D/空入力を一致させる）:
- `init` の対話モード（パッケージ検索・選択を含む）
- `require` のパッケージ省略時の検索・選択、dev 推奨の確認
- `remove` の別名キー確認、`update --interactive`
- 認証情報の入力（ユーザー名・パスワード。パスワードは非表示入力で、ログやキャッシュに出ないこと）
- ローカル変更を破棄する前の確認（discard-changes）、`create-project` の VCS 履歴削除の確認
- その他、現在「対話が必須で明示エラー」になっている箇所すべて（`grep` で洗い出してから着手する）
`-n`/`--no-interaction` と非端末入力の動作は今のまま変えない。

## 作業単位と検証
各項目を「Go実装 → 単体テスト → `testdata/compare` シナリオ追加 → `scripts/compare.sh` → 差異修正 → docs 更新」の単位で完了させてから次へ進む。2 と 3 の完了時にはそれぞれ `scripts/verify.sh` 全体（PHP 8.3/8.5 の compare を含む）を通す。最後に `scripts/ci-local.sh` も通す。

## 状態ファイル
セッションを終える直前に必ず `.cache/next/status` を書く。1行目は次のいずれか、2行目以降に理由と再開地点。
- `DONE` — 1〜3 がすべて実装・比較・記録済みで、最後の `scripts/verify.sh` と `scripts/ci-local.sh` が exit 0。
- `CONTINUE` — まだ進められる作業が残っている。
- `BLOCKED` — 認証・利用上限・外部障害など、自力で解消できない理由だけが残っている。
受入範囲を縮めて `DONE` にしない。
