# 互換性

基準: Composer `f83588392df43d4fd027d84d26de459946ef4e80`（2.11.999-dev+source）。gocomposer は Composer の解決・lock・取得・展開・autoload 生成を Go で再実装したもので、Composer 全仕様の完全互換ではない。下記の「対応」は公式Composerとの差分比較（`scripts/compare.sh`、PHP 8.3/8.4/8.5）または固定実プロジェクト統合（`scripts/integration.sh`）で一致を確認した範囲を示す。証拠は `docs/validation/` の各実行ディレクトリ（生ログ付き）。

## コマンド

| コマンド | 状態 | 確認範囲 |
|---|---|---|
| install | 対応 | lockからの導入、no-dev、二回目の冪等性、dry-run、progress表示、--audit（summary/table/plain/json）、ダウンロード失敗/checksum不一致/破損zip/展開先逸脱、失敗時のvendor保全、policy（malware filter list等）、Laravel/CLIアプリ/ライブラリの実プロジェクト |
| update | 対応 | 全体・部分（--with、-w/-W、--minimal-changes、--patch-only、--root-reqs、mirrors/lock/nothing）、--no-install、dry-run、platform/安定度/alias、解決不能診断、security advisory/abandoned/cooldownのblocking・audit、v1 inline と v2 metadata-url リポジトリ、--interactive のパッケージ選択（pty） |
| require | 対応 | バージョン推定（VersionSelector、cooldown）、--dev、--no-update、--dry-run、--fixed、--sort-packages、移動警告、失敗時のcomposer.json/lock復元、新規composer.json作成と削除、CRLF/タブのcomposer.json（JsonManipulator）、対話時のパッケージ省略検索・選択・dev 推奨・別キー移動（pty） |
| remove | 対応 | 通常・--dev・別セクション警告・未要求・ワイルドカード・--no-update・--dry-run・--unused・大文字小文字・allow-plugins設定の削除・依存で残る場合の終了2 |
| dump-autoload | 対応 | -o、--classmap-authoritative、--no-dev、PSR-0/4・classmap・files・exclude、platform_check、InstalledVersions |
| validate | 対応 | JSONスキーマ（lax/strict、justinrainbow/json-schema 6.13.1の文言と順序）、ConfigValidator、ValidatingArrayLoader、重複キー、lockの鮮度・不足要件、--with-dependencies、--strict |
| show / info | 対応 | 一覧・詳細・-t/--tree・json・-N・-P・-D・--no-dev・--locked・-s・-p・ワイルドカード、--latest/--outdated/--all/--available（--major-only/--minor-only/--patch-only、--direct、--strict、--ignore、--sort-by-age、cooldown、!/~/= 印と凡例、--ansi の色）。composer.json なしの --all/--available（packagist.org 前提）と、list API・v1 providers だけのリポジトリからの名前一覧は明示エラー |
| outdated | 対応 | ShowCommand への委譲（引数変換、--strict 終了コード、エラー時の show の synopsis） |
| audit | 対応 | table/plain/json/summary、--no-dev、--locked、-D、--abandoned、--ignore-severity、--ignore-unreachable、audit.* / policy.advisories 設定、advisory の metadata 形式と api-url 形式、到達不能リポジトリ（基準SHAの Auditor は終了コード 0/1 のみ） |
| run-script / run | 対応 | --list、--dev/--no-dev、--timeout、追加引数、予約名・未定義エラー、スクリプト名の直接コマンド化と scripts-aliases。Symfony Command クラスを指す script は明示エラー |
| exec | 対応 | --list、vendor/bin とルート bin、存在しない bin |
| clear-cache | 対応 | 全消去、--gc（ttl/maxsize/vcs）、cache-read-only。ttl 境界はファイル時刻の秒判定が実行ごとに揺れるため、決定的なシナリオと単体テストで確認 |
| fund / search / browse | 対応 | fund の text/json、search の -N/-O/--type/json、browse の -s/-H と起動コマンド不在時の警告。composer.json なしでの search/browse のパッケージ指定（packagist.org 前提）は明示エラー |
| diagnose | 対応（1項目差） | 全項目を移植。「Composer and its dependencies」の脆弱性確認は、gocomposer に監査対象の Composer インストールがないため Composer の非標準インストール警告を出す（KNOWN-DIFF） |
| self-update | 明示失敗 | gocomposer は phar ではなく自己更新機構がない。同じオプションを受けて「This instance of Composer does not have the self-update command.」と終了1。2行目の理由文だけ異なる（KNOWN-DIFF） |
| repository / repo, policy | 対応 | list/show/add/remove/set-url/get-url/enable/disable（-g、--file）、policy add-source |
| init | 対応 | 全オプション、git config・環境変数からの既定値、スキーマ検証失敗時の削除、--autoload。対話モード（質問一式、パッケージ検索・選択、検証エラーでの再質問、生成確認、.gitignore/.gitattributes、install の確認）を pty で比較 |
| bump | 対応 | --dev-only/--no-dev-only/--dry-run、終了コード、content-hash 更新、update の bump-after-update。lock 無効時は Composer が LogicException で落ちる（KNOWN-DIFF） |
| archive | 対応 | tar/zip/tar.gz（PharData 形式）、archive.exclude と export-ignore、--dir/--file/--ignore-filters、パッケージ指定。tar.bz2 は Go 標準にない圧縮のため、Composer（ext-bz2 なし）と同じエラー。ext-zip 環境の ZipArchiver 書式は未再現 |
| global | 対応 | 引数書き換え、ディレクトリ移動表示、COMPOSER の解除 |
| depends / why, prohibits / why-not | 対応 | -r、-t、--locked、platform要求、衝突・未導入の説明 |
| check-platform-reqs | 対応 | vendor/lock、--no-dev、--format=json |
| licenses | 対応 | text/json/summary、--no-dev、--locked |
| suggests | 対応 | --all、--list、--by-package、--by-suggestion、パッケージ指定 |
| about / list / help | 対応（gocomposer固有表示） | コマンド定義は公式から生成（commanddefs.json）。about/バージョン行はgocomposerと明記 |
| config | 対応 | 読取（既定値・ネスト・--source・--absolute・--list）、--unset/--file/--global/--auth/--json/--merge/--append、全スカラー/複数値キーの検証と文言、repositories・extra・suggest・scripts・platform・audit・policy・認証キー、スキーマ違反時の二重エラー表示（config-* 224シナリオ、20261023T060228Z-compare）。キー名のシェル補完、-vvv の読込ログ、bin-compat: symlink 非推奨警告は未移植 |
| reinstall | 対応 | 名前・ワイルドカード・--type、未インストール時の警告、依存順の uninstall/install とイベント、--prefer-source/--prefer-dist/--prefer-install、--no-autoloader/-o/-a、dist と git source の相互切替、ローカル変更時のエラーと discard-changes |
| status | 対応 | git source の未コミット変更・未 push・バージョン差（-v/-vv）、dist の再取得比較（_compare）、path（symlink/mirror）、終了コード 1/2/4 の加算、pre/post-status-cmd。svn/hg/fossil/perforce のパッケージは明示エラー |
| create-project | 対応 | バージョン指定・stability、--repository（composer/path/vcs）・--repository-url・--add-repository、--no-install/--no-scripts/--no-dev、--keep-vcs/--remove-vcs、self.version の置換、post-root-package-install/post-create-project-cmd、既存ディレクトリのエラー。リポジトリ指定なしは packagist.org への接続となり、ネットワークなしでは Composer と同じ curl error 6・終了100。対話端末での --ask と VCS 履歴削除の確認も pty で比較 |

## 機能別の制限

- **plugin**: 実行する。プラグインのコードだけをコマンド実行中に常駐する PHP プロセスで動かし、Composer API はスタブ経由で Go に戻す（解決・lock・取得・展開・autoload は Go）。allow-plugins 判定（非対話時の PluginBlockedException、対話時の y/n/d と書き戻し）、plugin-api-version、activate/deactivate/uninstall と更新時の再読込、イベント購読（init、command、pre/post-*-cmd、pre/post-package-*、pre-operations-exec、カスタムイベント）、CommandProvider、カスタムインストーラ（LibraryInstaller 継承・非継承、download/prepare/cleanup/installCode/updateCode/removeCode の上書き）に対応。pre-file-download（processedUrl・customCacheKey・transportOptions、type=package/metadata、再試行ごとの再発火）・post-file-download・pre-pool-create（パッケージの除外・差し替えを解決に反映）を発火する。スタブ API: getLoop/HttpDownloader、DownloadManager（組込みダウンローダのサブクラス、新しい dist type 用ダウンローダの追加）、getAutoloadGenerator()->dump()、RepositoryManager とリモートリポジトリ（findPackage(s)、search、loadPackages、getProviders、createRepository、add/prependRepository）、Factory::create（別 PHP ホストでそのプロジェクトのプラグインも読込）、Factory::createHttpDownloader、Config::merge、Installer::create()->run()、Locker・local repository の操作、JSON スキーマ検証。副作用（取得・配置）は Go が行う。symfony/flex v2.11.0（実ソースを fixture 化）の extra.symfony.require、recipe の取得・適用（loopback HTTPS の recipe エンドポイントスタブ）、pack の unpack、recipes/recipes:install/recipes:update、remove を比較。fixture プラグイン 58 + plugin2 34 + flex 22 シナリオと、composer/installers・php-http/discovery・dealerdirect/phpcodesniffer-composer-installer・pestphp/pest-plugin・symfony/skeleton（実 recipe 適用まで）の実プロジェクト統合で一致。
  - 明示エラー（Go 側の処理をプラグインの PHP 実装で置き換えられないため。`--no-plugins` で回避可）: getArchiveManager と Composer の set* によるサービス差し替え、AutoloadGenerator の buildPackageMap/parseAutoloads/createLoader・composer 以外への dump・独自 platform filter、setRepositoryClass・プラグイン定義のリポジトリクラス・リモートリポジトリ内パッケージの変更、配列設定での Factory::create と Factory::createGlobal、Config のソース操作と自作 Config への merge、setDownloader による組込みダウンローダの置換とダウンローダ内部（getFileName/extract/getLocalChanges）、HttpDownloader のリクエスト単位の ssl.cafile・クライアント証明書・proxy、Installer の setAdditionalFixedRepository/setAuditConfig/setPolicyConfig/getLockTransaction、Factory::create で作った Composer でのイベント発火、独自スキーマファイルでの検証、rar/xz と svn/hg/fossil/perforce のプラグインダウンローダ経由の取得、PHP が見つからない場合。
  - 既知の差: 並列転送時、Composer は post-file-download を完了順に発火するが、gocomposer は操作順（Composer の `COMPOSER_MAX_PARALLEL_HTTP=1` と同じ）。`disablePlugins='local'` の Factory::create ではグローバルプラグインも無効になる。flex が出すランダムな session id は比較時に `SESSION` へ正規化。
  - 未比較・差異: validate と reinstall では command イベントを発火しない。`-vvv` の Factory デバッグログは出さない。
- **scripts**: コマンド文字列、`@php`、`@composer`、`@putenv`、スクリプト参照、PHPクラスコールバックに対応。クラスコールバックは別PHPプロセスでComposer APIのスタブ（Event/IO/Config/RootPackageの一部）を渡して実行する。スタブにないAPIを呼ぶコールバックは明示的な例外で失敗する。Laravelの `ComposerScripts::postAutoloadDump` と `artisan package:discover` は実プロジェクトで一致を確認。
- **リポジトリ**: composer（v1 inline/includes、v2 metadata-url、available-packages、providers-api、security-advisories metadata/api-url）、path、package、artifact（zip/tar/tgz、ext-zip 付き oracle で比較）、vcs/git（タグ・ブランチ・dev 版・reference、VCS キャッシュ、only-branches 等）、github/gitlab/bitbucket（API ドライバ。loopback の HTTPS スタブで比較、github-domains/gitlab-domains、非公開リポジトリと認証、レート制限）。v1 providers-url は明示的に未対応エラー。svn/hg/fossil/perforce/pear は Composer と同じ検出エラー・警告までを再現し、ドライバ本体は明示的に未対応エラー。
- **source インストール**: git のみ（--prefer-source/--prefer-dist、preferred-install の文字列・パターン、dist への fallback、dist⇔source 切替、update 時の checkout、未コミット変更の検出と discard-changes true/false/stash、対話時の y/n/v/d/s/? 確認（update・remove・reinstall）を pty で比較）。フォージの HTTPS 上の git clone は git が Composer の cafile を使えないため比較対象外。破棄確認で EOF を入力すると Composer と同じく再質問を繰り返す（シナリオなし）。
- **認証**: 401/403 の非対話処理（AuthHelper::promptAuthIfNeeded）、git config からのトークン取得、GitLab CI job token、Bitbucket の access token、GitLab OAuth トークンの更新。対話入力: 401/403 の Username/Password（パスワードは非表示入力）と store-auths（true/false/prompt）による保存、GitHub トークン入力と検証・保存先選択、GitLab の Username/Password からの OAuth 取得、Bitbucket OAuth consumer の token 交換、git 認証失敗時の入力を pty で比較。パスワード・トークンが auth.json（保存を選んだ場合のみ、公式と同じ）以外の出力・キャッシュ・vendor・lock に出ないことを確認。diagnose の HTTP 取得は対話認証を行わない。トークンが stdout/stderr・キャッシュ・vendor に出ないことを -vvv を含めて確認（URL に埋め込まれた認証情報は Composer と同じく lock/installed.json に残る）。
- **ダウンロード**: http(s)/file、Basic/Bearer/GitHub/GitLab/custom-headers 認証、再試行、checksum、キャッシュ（Composerと同じキー）。curl固有のエラー番号は DNS失敗(6)・接続拒否(7) のみ再現し、その他は Go のエラー文をそのまま含む。
- **アーカイブ**: zip/tar(.gz/.bz2)/単一gzip。zipの `../` は unzip と同様に除去してパッケージ内に展開し、パッケージ外を指すシンボリックリンクとシンボリックリンク経由の書き込みは拒否する（Composerより厳しい）。rar/xz/7z は未対応。
- **対話**: 入力が端末のときだけ対話モードになる（Composer と同じ isatty 判定。`-n`、`COMPOSER_NO_INTERACTION`、`-q`、非端末入力では従来どおり既定値で進む）。Symfony QuestionHelper（確認・質問・検証付き再質問・非表示入力・選択/複数選択・オートコンプリート、Ctrl-D で既定値または "Aborted."）と OSC 7501/9;4 の端末状態表示を移植。init、require（パッケージ省略時の検索・選択、dev 推奨、別キーからの移動、virtual パッケージのバージョン、類似名からの選択）、remove の別キー確認、update --interactive、run-script/exec の選択、create-project の --ask と VCS 履歴削除、allow-plugins の y/n/d/?、親ディレクトリの composer.json 利用確認、認証入力、discard-changes を、比較ハーネスの疑似端末（`tools/compare/ptyrun`、80x24、TERM=xterm）で決めた入力列を流して比較（tty-* 81 + tty-auth/tty-discard 60 シナリオ）。root 実行時の確認は実装済みだがコンテナが uid 1000 のため未比較。
- **端末表示**: 非端末出力はComposerと一致を確認。端末（decorated）出力のハイパーリンク（OSC 8）と OSC 9;4 進捗は未出力、ProgressBarの上書き描画は簡易実装。
- **時刻依存の文言**: `show` の相対日付等は同時刻実行で一致。
- **Composer自体がPHPエラーで落ちる入力**: composer.json の require 値が配列（ConfigValidatorのTypeError）や COMPOSER_AUTH の username が数値（BaseIOのTypeError）の場合、gocomposer はクラッシュせずスキーマエラー表示/無視とする（比較では KNOWN-DIFF として記録）。
- **ResponseWarnings**: リポジトリの warning/info 応答の表示は filter api-url・advisory API の応答について未実装。
- **キャッシュ GC**: Composer と同じく files cache の GC をプロセスごとに 1/51 の確率で実行し、`COMPOSER_TEST_SUITE` で無効になる。GC はキャッシュディレクトリを作成・削除し比較結果を揺らすため、比較ハーネスは両ツールに `COMPOSER_TEST_SUITE=1` を渡す（基準SHAの Composer ではこの変数はこの GC にしか使われない）。GC 自体は単体テストで確認。
- **比較での正規化（Composer 側の非決定性）**: (1) uninstall フェーズの進捗バーは Composer の非同期 `rm -rf` が待機開始時に残っているかで「0/N」か不定形かが変わる（過去12回の oracle 実行で 11 回が不定形）。`reinstall-progress` だけ「Removing」直後のバーを不定形に揃えて比較し、gocomposer は通常の不定形を出す。(2) PHP 8.5 で Composer 自身のソース（例 `src/Composer/Util/Svn.php:318`）が出す PHP の Deprecation Notice は oracle 側 stderr からその行だけ除外する。いずれも生出力は結果 JSON に残る。
- **時刻依存のシナリオ**: show の cooldown 残り時間（「Nd Nh left」）は、oracle と gocomposer の実行が時の境界をまたぐと表示が変わりうる。
- **PHP 拡張なしでの動作**: platform 検出には PHP 実行ファイルが必要（GOCOMPOSER_PHP または PATH の php）。見つからない場合は php/ext-* を提供しない platform として扱う。

## 検証環境

- 公式 PHP（ベースイメージの digest。`.cache/oracle-images.txt` と各 compare の `summary.json` に記録）:

  | PHP | ベースイメージ | oracle イメージ（`oracle/Containerfile`） | ext-zip 版（`oracle/Containerfile.zip`） |
  |---|---|---|---|
  | 8.4.27（主比較） | `php@sha256:a4ae1f57e22b89a0a69c5e085b3010ee727d094e1542d249e48a01bec3c44810`（multi-arch index。arm64: `05c3a54f…bfd5`、amd64: `a007a658…5a98`） | `localhost/gocomposer-oracle:8.4`（`sha256:7c3d8598…3607`） | `:8.4-zip`（`sha256:cd7f1f04…29fa`） |
  | 8.3.36 | `php@sha256:f1b772b8a2cc9ba2ba7415d21422f2c002c6f5789b53bcac2ababc5e5d89fb47` | `:8.3`（`sha256:91789a66…5fce`） | `:8.3-zip`（`sha256:08e40da7…6bf9`） |
  | 8.5.12 | `php@sha256:361f9916dee86b0c0b67dfd196ce5cd8ebcecb089724d867a57f564fa94909fd` | `:8.5`（`sha256:8f997f44…92e2`） | `:8.5-zip`（`sha256:c0ec954e…4528`） |

  oracle / ext-zip イメージの digest はローカルでビルドするたびに変わりうる（数値は 20261023T093500Z・101132Z・104618Z の compare 実行時）。compare はいずれも `--network=none` で、loopback の fixture リポジトリ、git bare リポジトリ、フォージ API の HTTPS スタブを使う。
- Composer: `.cache/oracle-source`（基準SHAの独立 clone、composer.lock sha256 `7e7ed6af…34c`）。
- 実プロジェクト（`scripts/integration.sh`、実ネットワーク。lock は公式 Composer で生成して `testdata/integration/locks` に固定）: symfony/console v7.4.9（ライブラリとして root）、CLI アプリ（`testdata/integration/cli-app`）、laravel/laravel v13.9.0、live-psr（小規模なライブダウンロード）、plugin-installers（composer/installers）、plugin-php-http-discovery（php-http/discovery）、plugin-phpcs（dealerdirect/phpcodesniffer-composer-installer）、plugin-pest（pestphp/pest-plugin）、symfony/skeleton v7.4.99（symfony/flex、install・update --no-install に加えて実 recipe エンドポイントからの recipes:install）の 9 件。
- 最新結果: compare は PHP 8.4/8.3/8.5 とも 1386 シナリオ（うち疑似端末 141）中 MATCH 1377・KNOWN-DIFF 9・DIFF 0（20261023T093500Z / 101132Z / 104618Z）、integration は 9 件の全ステップが MATCH（20261023T100947Z）、`scripts/verify.sh` exit 0（20261023T093500Z-verify）。
- KNOWN-DIFF 9 件: validate-lax-fail・composer-auth-invalid・bump-lock-disabled（Composer 自体が PHP エラー/LogicException で落ちる）、diagnose-network（Composer 自身の監査行）、self-update-* 5 件（自己更新機構なし）。
- CI: `.github/workflows/ci.yml` の全ステップを `scripts/ci-local.sh` でローカル再現し全 PASS（20261023T112150Z-ci-local。compare は 3 版とも 1386 中 MATCH 1377・KNOWN-DIFF 9・DIFF 0、integration 9 件全 MATCH。unit は macOS の go1.23.12 と podman の `golang:1.23`（`sha256:ac2b6b75…`）、oracle ジョブはホスト上で podman を使用）。GitHub Actions 上では push が範囲外のため未実行。
