---
name: gocomposer-remaining
description: gocomposerの未対応項目（show --latest/--outdated/--all/--available、outdated/audit/run-script/create-project等の未実装コマンド、plugin実行、VCSリポジトリ、source install、CI実行）を確認なしで一括実装・公式比較・記録する。prompts/remaining.md から使う。
---

# 未対応項目の一括実装

`gocomposer-development` スキル、`CLAUDE.md`、`docs/implementation-plan.md`、`docs/progress.md`、`docs/compatibility.md` を先に読み、その原則（Go本体で処理、公式PHP/Composerとの差分比較、未対応は明示失敗、参照元は読取専用、commit/push/公開なし）をすべて守る。

## 確認フローなし
ユーザーへの質問・承認待ち・「続けますか」は一切行わない（ユーザーは不在）。仕様判断は基準SHAのComposerソースと公式比較の結果で決め、判断理由を `docs/progress.md` に書く。
AskUserQuestion・EnterPlanMode は使わない。迷ったら参照実装に合わせ、合わせられない部分は明示失敗にして次へ進む。

## 対象と順序
依存関係の少ない順に進める。各項目は「Go実装 → 単体テスト → `testdata/compare` に正常・失敗・副作用シナリオ追加 → `scripts/compare.sh` で公式比較 → 差異修正 → `docs/compatibility.md`/`docs/progress.md` 更新」を1単位として完了させてから次へ行く。

1. **show --latest / --outdated / --all / --available**: Composer `ShowCommand` の最新版探索（VersionSelector、`--major-only/--minor-only/--patch-only`、`--direct`、`--strict`、`--ignore`、`--format=json`、色・semver互換の印）を移植。
2. **outdated**: ShowCommand への委譲（`show --latest --outdated`）として実装し、終了コード（`--strict`）を一致させる。
3. **audit**: 既存 `internal/policy` の Auditor を使い、`--format`、`--no-dev`、`--locked`、`--abandoned`、`--ignore-severity`、advisory API/metadata の両方を比較。
4. **run-script / run / exec / reinstall / clear-cache / status / init / bump / fund / search / browse / diagnose / archive / repository / policy / global / create-project / self-update**: 影響の大きい順（run-script, exec, reinstall, clear-cache, bump, status, create-project, global, init, search, fund, archive, repository, policy, diagnose, browse, self-update）に移植。
   - `self-update` は gocomposer 自身の更新機構がないため、Composer と同じ引数を受けて「gocomposer は self-update 非対応」と明示失敗するのが正しい挙動なら、その理由を compatibility.md に書く。成功を装わない。
   - `diagnose` / `browse` のようにネットワークや外部プログラムに依存するものは、loopback fixture とスタブ環境で比較できる範囲を比較し、残りを記録する。
5. **VCSリポジトリ（git優先、次にsvn/hg/perforceは可能な範囲）と source インストール**: Composer の VcsRepository/GitDriver/GitHubDriver/GitLabDriver/BitbucketDriver、`GitDownloader`、`--prefer-source`、`preferred-install`、dist/source 切替、reference の lock 記録を移植。比較は `git` を持つ oracle コンテナ内のローカル bare リポジトリ（`file://`、`--network=none`）で行う。GitHub API 依存のドライバはスタブHTTPサーバーで検証する。source を含む vendor の再現・更新・未コミット変更検出（`status`、`discard-changes`）も比較する。
6. **plugin 実行**: Composer plugin は PHP コードで Composer API を直接使うため、既存の script コールバック方式（別PHPプロセス＋Composer APIスタブ）を拡張する方針を基本とする。最低限、実在の主要プラグイン（composer/installers の installer 型、php-http/discovery、dealerdirect/phpcodesniffer-composer-installer、pestphp/pest-plugin など）の lock/vendor/autoload/出力が公式と一致するところまで比較する。スタブ外のAPIを呼ぶプラグインは明示的に失敗させ、`--no-plugins` での続行を維持する。`allow-plugins` の判定・プロンプト非対話時の既定値も公式に合わせる。
7. **CI**: `.github/workflows/ci.yml` の全ステップをローカルで実行する。unit ジョブは macOS 上でそのまま、Linux は podman の `docker.io/library/golang:1.23` 系コンテナ（ubuntu-24.04 相当）で実行し、oracle ジョブは podman で実行する（`act` が使えるなら併用可）。結果は `docs/validation/<UTC>-ci-local/` に保存する。GitHub Actions 上の実行は push が必要で本作業の範囲外なので、「ローカル再現のみ・GitHub 上は未実行」と明記する。成功扱いを偽らない。

## 検証と記録
- 項目ごとに `go test ./...` と該当シナリオの `scripts/compare.sh` を通す。複数項目がまとまったら `scripts/verify.sh` 全体（PHP 8.3/8.5 の compare も含む）を通す。
- 既存の 372 MATCH を退行させない。退行したら先に直す。
- 各項目の完了時に `docs/progress.md` の「実装済み」「未実装・未比較」「次の作業」と `docs/compatibility.md` のコマンド表を更新する。
- 移植したComposerコード・fixtureには元ファイルと基準SHAを記録する。

## セッション終了時の状態ファイル
セッションを終える直前に必ず `.cache/remaining/status` を書く。1行目は次のいずれか、2行目以降に理由と再開地点。
- `DONE` — 上記1〜7がすべて実装・比較・記録済みで、最後の `scripts/verify.sh` が exit 0。
- `CONTINUE` — まだ進められる作業が残っている（コンテキストや時間の都合で区切った）。
- `BLOCKED` — 認証・利用上限・外部障害など、自力で解消できない理由だけが残っている。
受入範囲を縮めて `DONE` にしない。
