/gocomposer-next

/Users/kouji/dev/gocomposer の gocomposer について、AGENTS.md の運用ルールと gocomposer-next スキルの手順に従い、次の3項目を一括で実装・検証してください。

対象（この順で進める）:
1. 文書の整理（docs/compatibility.md の「検証環境」節の更新と、progress.md・README.md との矛盾の解消）
2. symfony/flex への対応とプラグイン用スタブ API の拡充（pre-file-download・post-file-download・pre-pool-create イベント、getDownloadManager・getAutoloadGenerator・リモートリポジトリ・Factory::create・Config::merge、download/installCode を上書きするインストーラ、symfony/skeleton の実プロジェクト統合）
3. 対話プロンプトの実装（pty を使った比較ハーネスを追加し、init・require・remove・update --interactive・認証情報の入力・ローカル変更の破棄確認など、現在「対話が必須で明示エラー」になっている箇所すべて）

進め方:
- 確認・承認・質問は一切せず、最後まで連続して進めてください。ユーザーは不在です。
- 最初に docs/progress.md と .cache/next/status（あれば）を読み、完了済みの項目は飛ばして続きから始めてください。
- 前回までの成果（MATCH 1180・DIFF 0、installer fixture 454件、実プロジェクト8件、ローカル CI）を退行させないでください。
- 各項目を「実装 → テスト → 公式比較 → 差異修正 → docs 更新」の単位で完了させ、ログは docs/validation/ に残してください。
- 対応できないものは成功を装わず明示エラーにし、理由を docs/compatibility.md に書いてください。
- commit・push・公開・グローバルインストール、参照元 /Users/kouji/oss/composer やユーザーの Composer 設定・キャッシュ・認証情報の変更はしないでください。
- コンテキストが長くなっても途中で止めず続けてください。セッションを終える前に、必ず .cache/next/status をスキルの規約どおり書いてください。
