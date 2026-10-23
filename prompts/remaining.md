/gocomposer-remaining

/Users/kouji/dev/gocomposer の gocomposer について、AGENTS.md の運用ルールと gocomposer-remaining スキルの手順に従い、docs/progress.md と docs/compatibility.md に残っている未対応項目をすべて一括で実装・検証してください。

対象（この順で進める）:
1. show の --latest、--outdated、--all、--available
2. outdated コマンド
3. audit コマンド
4. run-script、exec、reinstall、clear-cache、bump、status、create-project、global、init、search、fund、archive、repository、policy、diagnose、browse、self-update など未実装のコマンドすべて
5. VCS リポジトリ（git 優先）とソースからのインストール
6. プラグインの実行
7. CI 定義（.github/workflows/ci.yml）のローカル実行（macOS と podman の Linux コンテナ。GitHub 上の実行は push が必要なので範囲外として記録）

進め方:
- 確認・承認・質問は一切せず、最後まで連続して進めてください。ユーザーは不在です。
- 最初に docs/progress.md と .cache/remaining/status（あれば）を読み、完了済みの項目は飛ばして続きから始めてください。
- 前回までの成果（372 MATCH、installer fixture 454件、実プロジェクト統合）を退行させないでください。
- 各項目を「実装 → テスト → 公式比較 → 差異修正 → docs 更新」の単位で完了させ、ログは docs/validation/ に残してください。
- 未対応のまま残すものは成功を装わず明示エラーにし、理由を docs/compatibility.md に書いてください。
- commit・push・公開・グローバルインストール、参照元 /Users/kouji/oss/composer やユーザーの Composer 設定・キャッシュ・認証情報の変更はしないでください。
- コンテキストが長くなっても途中で止めず続けてください。セッションを終える前に、必ず .cache/remaining/status を AGENTS.md の規約どおり書いてください。
