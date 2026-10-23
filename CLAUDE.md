# gocomposer

GoでComposer互換の依存管理CLIを実装する。作業先は /Users/kouji/dev/gocomposer。
参照元 /Users/kouji/oss/composer は読取専用として扱い、ソース・lock・vendor・Git状態を変更しない。
最初に `.claude/skills/gocomposer-development/SKILL.md`、`docs/implementation-plan.md`、`docs/upstream.json` を読む。
無人で連続実行するときの運用ルールは @AGENTS.md に従う。

## 作業の承認範囲
このプロジェクト内の設計、Go実装、テスト、fixture追加、文書更新、依存取得、一時的な公式PHP比較環境、loopbackサーバーによる検証は承認済み。
段階完了のたびにユーザー確認を求めず、次の未完了事項へ進む。仕様判断は参照実装と差分テストで決め、判断理由を記録する。
既存gophp/gophpenv/Laravelアプリ、手動DB、ユーザーのComposer設定・キャッシュ・認証情報は変更しない。commit・push・タグ・公開・グローバルインストールはこの作業に含まない。
認証や管理ポリシーの解除を独自に行わない。環境障害では独立した作業を続け、同一原因の無限再試行をしない。

## 実装原則
Composerを子プロセスで呼ぶだけのラッパーを成果物にしない。解決・lock・取得・展開・autoload生成はGo本体が行う。
公式PHP/Composerは比較基準として別の一時ディレクトリで実行する。生成するautoloadのPHPコードや、明示設定されたPHP script/plugin実行とは区別する。
PHPのsemverはGo module semverで代用しない。未対応機能は明示的に失敗させ、成功と偽らない。
既存ファイル、lock、vendorは失敗時も保全する。展開先逸脱・symlink逸脱・認証情報のログ流出を防ぐ。
参照/移植するComposerコード・fixtureには元ファイルと基準SHAを記録し、MIT表示を保持する。

## 継続と完了
`docs/progress.md` に実装済み・検証済み・未比較・阻害要因・次の作業を随時記録する。
単体テストが通っただけで全体完了としない。`docs/implementation-plan.md` の受入条件と実証ログを満たす。
実行上限・認証・外部環境で継続できない場合のみ、完了扱いにせず具体的な再開地点を残す。
