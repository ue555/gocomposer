/gocomposer-final

/Users/kouji/dev/gocomposer の gocomposer について、AGENTS.md の運用ルールと gocomposer-final スキルの手順に従い、残っている作業をすべて一括で実装・検証してください。

対象（この順で進める。詳細はスキルを参照）:
1. 文書の更新（GitHub Actions run 39983736909 の成功を反映し、文書内の矛盾を解消）
2. completion コマンド（bash/zsh/fish、_complete、config のキー名補完など）
3. 細部の移植（config の -vvv 読込ログと bin-compat 警告、Factory の -vvv ログ、validate/reinstall の command イベント、Symfony Command クラスを指す run-script、composer.json なしの fund/search/browse）
4. VCS の svn・hg・fossil・perforce
5. v1 providers-url リポジトリ
6. アーカイブの rar・xz・7z
7. ResponseWarnings（filter api-url と advisory API の警告表示）
8. 端末表示（OSC 8 ハイパーリンク、ProgressBar の上書き描画、curl のエラー番号）
9. プラグインのサービス差し替えなど、残っている明示エラーの API と post-file-download の発火順
10. root 実行時の確認プロンプトの公式比較
11. 統合テストに --no-scripts 付きの update ステップを追加
12. CI の oracle ジョブを PHP バージョンごとの matrix に分割

進め方:
- 確認・承認・質問は一切せず、最後まで連続して進めてください。ユーザーは不在です。
- 最初に docs/progress.md と .cache/final/status（あれば）を読み、完了済みの項目は飛ばして続きから始めてください。
- 前回までの成果（3版とも MATCH 1377・DIFF 0、installer fixture 454件、実プロジェクト9件、ローカル CI と GitHub Actions の成功）を退行させないでください。
- self-update、Composer 自体が落ちる入力、diagnose の監査行の KNOWN-DIFF は設計上の差なので変えないでください。
- 各項目を「実装 → テスト → 公式比較 → 差異修正 → docs 更新」の単位で完了させ、ログは docs/validation/ に残してください。
- 対応できないもの・公式と比較できないものは成功を装わず、明示エラーにするか理由を docs/compatibility.md に書いてください。
- commit・push・公開・グローバルインストール、参照元 /Users/kouji/oss/composer やユーザーの Composer 設定・キャッシュ・認証情報の変更はしないでください。
- コンテキストが長くなっても途中で止めず続けてください。セッションを終える前に、必ず .cache/final/status をスキルの規約どおり書いてください。
