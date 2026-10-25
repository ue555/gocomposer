# AGENTS.md — gocomposer 自律作業ルール

このファイルは gocomposer を操作するエージェント（Claude Code ほか）向けの行動規約です。基本方針は `CLAUDE.md` に従い、ここでは「無人で連続実行する」ための運用を定めます。

## 前提として読むもの
1. `CLAUDE.md`
2. `.claude/skills/gocomposer-development/SKILL.md`（開発ループと比較契約）
3. `.claude/skills/gocomposer-remaining/SKILL.md`（未対応項目の一括実装手順）
4. `docs/implementation-plan.md`、`docs/progress.md`、`docs/compatibility.md`、`docs/upstream.json`

## 確認フローなし
- ユーザーは不在とみなす。質問・承認待ち・「続けますか」・計画承認（EnterPlanMode/AskUserQuestion）を行わない。
- 仕様判断は基準SHAの Composer ソース（`.cache/oracle-source`、参照元 `/Users/kouji/oss/composer` は読取専用）と公式比較の結果で決め、理由を `docs/progress.md` に記録する。
- 判断に迷う機能は参照実装に合わせる。合わせられない部分は明示エラーにして次の項目へ進む。止まらない。

## してはいけないこと
- commit・push・タグ・リリース・GitHub Actions の起動・グローバルインストール。
- 参照元 `/Users/kouji/oss/composer`、ユーザーの Composer 設定・キャッシュ・`auth.json`、既存の gophp/gophpenv/Laravel アプリの変更。
- 認証や管理ポリシーの回避、同一原因での無限再試行。
- 受入範囲を縮めて完了扱いにすること。未実行の検証を成功と書くこと。
（`.claude/settings.json` の deny ルールでも同じ操作を禁止している）

## 作業単位
未対応項目ごとに「Go実装 → 単体テスト → `testdata/compare` シナリオ追加 → `scripts/compare.sh` → 差異修正 → `docs/compatibility.md`・`docs/progress.md` 更新」を1単位として完了させてから次へ進む。複数単位ごとに `scripts/verify.sh` を通し、既存の 372 MATCH・installer fixture 454件・実プロジェクト統合を退行させない。

## サブエージェント
調査（Composer ソースの読解、既存 Go コードの所在確認）は並列のサブエージェントに任せてよい。同じファイルを複数のエージェントが同時に編集しないこと。実装と検証の最終判断は親セッションが行う。

## 終了時
セッションを終える直前に、使っているスキルが指定する状態ファイルを書く（`gocomposer-remaining` は `.cache/remaining/status`、`gocomposer-next` は `.cache/next/status`、`gocomposer-final` は `.cache/final/status`）。1行目は `DONE` / `CONTINUE` / `BLOCKED`、2行目以降に理由と再開地点。`DONE` は対象項目すべての実装・比較・記録が済み、スキルが指定する最終検証が exit 0 のときだけ。
