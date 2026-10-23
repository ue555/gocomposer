# Composer vs gocomposer integration (fixed real projects)

- gocomposer_linux_sha256: `dadc3e20e1ef9365f4d948e69929d625acb5cc18fb222aeaf2c1aff758c5a5b2`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| project | step | status | oracle exit | go exit |
|---|---|---|---|---|
| plugin-installers | install | DIFF | 0 | 1 |
| plugin-installers | update-no-install | MATCH | 0 | 0 |
| plugin-php-http-discovery | install | MATCH | 0 | 0 |
| plugin-php-http-discovery | update-no-install | MATCH | 0 | 0 |
| plugin-phpcs | install | DIFF | 0 | 1 |
| plugin-phpcs | update-no-install | MATCH | 0 | 0 |
| plugin-pest | install | MATCH | 0 | 0 |
| plugin-pest | update-no-install | MATCH | 0 | 0 |

## plugin-installers / install

```
exit code: oracle=0 go=1
stderr:
     + (go)     Install of roots/soil failed
     + (go) 
     + (go) In <location>:
     + (go)                                                                                
     + (go)   rename /work/project/vendor/composer/d27c98d0/roots-soil-e8e15c2 /work/proj  
     + (go)   ect/web/app/mu-plugins/soil/: file exists                                    
     + (go)                                                                                
     + (go) 
     + (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
     + (go) 
     - (oracle)   - Installing yoast/duplicate-post (4.6): Extracting archive
     - (oracle) Package roots/soil is abandoned, you should avoid using it. Use roots/acorn-prettify instead.
     - (oracle) Generating autoload files
     - (oracle) 2 packages you are using are looking for funding.
     - (oracle) Use the `composer fund` command to find out more!
autoload probe:
     - (oracle) resolved: 0 of 0
     - (oracle) root: gocomposer/integration-plugin-installers 1.0.0+no-version-set dev=true
     - (oracle) installed: composer/installers 'v2.3.0' '12fb2dfe5e16183de69e784a7b84046c43d97e8e'
     - (oracle) installed: gocomposer/integration-plugin-installers '1.0.0+no-version-set' NULL
     - (oracle) installed: roots/soil '4.1.1' 'e8e15c29f01414f86610445be8f7274b8b9a73b8'
     - (oracle) installed: yoast/duplicate-post '4.6' 'db85005366144846cac3b07632488a6da92876f5'
project tree:
     - (oracle) D web/app/mu-plugins/soil
     - (oracle) D web/app/mu-plugins/soil/resources
     - (oracle) D web/app/mu-plugins/soil/resources/views
     - (oracle) D web/app/mu-plugins/soil/src
     - (oracle) D web/app/mu-plugins/soil/src/Exceptions
     - (oracle) D web/app/mu-plugins/soil/src/Modules
     - (oracle) D web/app/mu-plugins/soil/tests
     - (oracle) D web/app/mu-plugins/soil/tests/Unit
     - (oracle) D web/app/mu-plugins/soil/tests/Unit/Helpers
     - (oracle) D web/app/mu-plugins/soil/tests/Unit/Modules
     - (oracle) D web/app/mu-plugins/soil/tests/__fixtures__
     - (oracle) D web/app/mu-plugins/soil/tests/__fixtures__/modules
     - (oracle) D web/app/plugins
     - (oracle) D web/app/plugins/duplicate-post
     - (oracle) D web/app/plugins/duplicate-post/compat
     - (oracle) D web/app/plugins/duplicate-post/css
     - (oracle) D web/app/plugins/duplicate-post/js
     - (oracle) D web/app/plugins/duplicate-post/js/dist
     - (oracle) D web/app/plugins/duplicate-post/src
     - (oracle) D web/app/plugins/duplicate-post/src/admin
     - (oracle) D web/app/plugins/duplicate-post/src/admin/views
     - (oracle) D web/app/plugins/duplicate-post/src/handlers
     - (oracle) D web/app/plugins/duplicate-post/src/ui
     - (oracle) D web/app/plugins/duplicate-post/src/watchers
     - (oracle) F vendor/autoload.php 644 c6f7348923cce470
     - (oracle) F vendor/composer/ClassLoader.php 644 d4b2313e562edf94
     + (go) F vendor/composer/installed.json 644 238424b973810024
     + (go) F vendor/composer/installed.php 644 1f7b43a4f66db5a2
     - (oracle) F vendor/composer/LICENSE 644 7855ac293067aebe
     - (oracle) F vendor/composer/autoload_classmap.php 644 e846382cea8953c1
     - (oracle) F vendor/composer/autoload_files.php 644 0796a206e0e9efa0
     - (oracle) F vendor/composer/autoload_namespaces.php 644 ede505407590638a
     - (oracle) F vendor/composer/autoload_psr4.php 644 8e2e8e96fe5cc19a
     - (oracle) F vendor/composer/autoload_real.php 644 3623d993763643dd
     - (oracle) F vendor/composer/autoload_static.php 644 a77140da4abde956
     - (oracle) F vendor/composer/installed.json 644 37dc69d78f6c8edb
     - (oracle) F vendor/composer/installed.php 644 b55fcbc3656ee527
     - (oracle) F vendor/composer/platform_check.php 644 efa64a484ec52731
     - (oracle) F web/app/mu-plugins/soil/CHANGELOG.md 644 4b90f6bbbfa1adaf
     - (oracle) F web/app/mu-plugins/soil/LICENSE.md 644 4e78fbeaeed2d876
     - (oracle) F web/app/mu-plugins/soil/README.md 644 e3617ed0d7748e2b
     - (oracle) F web/app/mu-plugins/soil/composer.json 644 02436cc125cd71a9
     - (oracle) F web/app/mu-plugins/soil/phpcs.xml.dist 644 f2a1a679bb7aaf70
     - (oracle) F web/app/mu-plugins/soil/phpunit.xml.dist 644 db1bc519e7471825
     - (oracle) F web/app/mu-plugins/soil/resources/views/google-analytics.php 644 c19b233276aec837
     - (oracle) F web/app/mu-plugins/soil/soil.php 644 7365aa42b1d21a51
     - (oracle) F web/app/mu-plugins/soil/src/DOM.php 644 d05f109c2e862c91
     - (oracle) F web/app/mu-plugins/soil/src/Exceptions/LifecycleException.php 644 7ac7723b3ceb35fd
     - (oracle) F web/app/mu-plugins/soil/src/Modules/AbstractModule.php 644 6dd93da53d3463d0
     - (oracle) F web/app/mu-plugins/soil/src/Modules/CleanUpModule.php 644 d963056539354978
     - (oracle) F web/app/mu-plugins/soil/src/Modules/DisableAssetVersioningModule.php 644 8713d22a6ab444de
     - (oracle) F web/app/mu-plugins/soil/src/Modules/DisableRestApiModule.php 644 21cc9651a615cfcf
     - (oracle) F web/app/mu-plugins/soil/src/Modules/DisableTrackbacksModule.php 644 2cd05a6e9d4fef34
     - (oracle) F web/app/mu-plugins/soil/src/Modules/GoogleAnalyticsModule.php 644 71ebff195e6a383a
     - (oracle) F web/app/mu-plugins/soil/src/Modules/JsToFooterModule.php 644 8d0f45b953c4f042
     - (oracle) F web/app/mu-plugins/soil/src/Modules/NavWalkerModule.php 644 168888c9f84a4465
     - (oracle) F web/app/mu-plugins/soil/src/Modules/NiceSearchModule.php 644 b1aa148b75b8c42c
     - (oracle) F web/app/mu-plugins/soil/src/Modules/RelativeUrlsModule.php 644 ef70d7f2a34cd6a9
     - (oracle) F web/app/mu-plugins/soil/src/NavWalker.php 644 ecf244cebf4639ac
     - (oracle) F web/app/mu-plugins/soil/src/Options.php 644 bf3b0eb3ee7294d1
     - (oracle) F web/app/mu-plugins/soil/src/Soil.php 644 a2e4bc2cdb69a8f1
     - (oracle) F web/app/mu-plugins/soil/src/autoload.php 644 9e348240676329d6
     - (oracle) F web/app/mu-plugins/soil/src/helpers.php 644 80c1685e3bb6b684
     - (oracle) F web/app/mu-plugins/soil/tests/TestCase.php 644 1955075d65c90183
     - (oracle) F web/app/mu-plugins/soil/tests/TestCaseLegacy.php 644 640fa94f71d8105b
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/DomTest.php 644 d00c14dd37e341e1
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/EntrypointTest.php 644 96ab48d18afe3ab2
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/Helpers/UrlCompareTest.php 644 598afa78cfaaca2f
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/ModuleTest.php 644 878c05b8b3e5906d
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/Modules/CleanUpModuleTest.php 644 7b5c77c094d917fa
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/Modules/RelativeUrlsModuleTest.php 644 347166af420e9820
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/OptionsTest.php 644 28d1f6e4dc0a466d
     - (oracle) F web/app/mu-plugins/soil/tests/Unit/SoilTest.php 644 c48696869fcc7f6c
     - (oracle) F web/app/mu-plugins/soil/tests/__fixtures__/modules/CustomHookModule.php 644 052a3f9363b44314
     - (oracle) F web/app/mu-plugins/soil/tests/__fixtures__/modules/CustomNameModule.php 644 46207c918424da73
     - (oracle) F web/app/mu-plugins/soil/tests/__fixtures__/modules/StubModule.php 644 e61a5a01591d9e8b
     - (oracle) F web/app/mu-plugins/soil/tests/api.php 644 911bea9625b49ebe
     - (oracle) F web/app/mu-plugins/soil/tests/bootstrap.php 644 0c27931c462db585
     - (oracle) F web/app/mu-plugins/soil/tests/helpers.php 644 546cc52a4676ef84
     - (oracle) F web/app/plugins/duplicate-post/admin-functions.php 644 7a0597210447d6c3
     - (oracle) F web/app/plugins/duplicate-post/common-functions.php 644 9fc94155fd410a2a
     - (oracle) F web/app/plugins/duplicate-post/compat/jetpack-functions.php 644 281ea5d7e268f0ff
     - (oracle) F web/app/plugins/duplicate-post/compat/wpml-functions.php 644 d4af56de00288fe4
     - (oracle) F web/app/plugins/duplicate-post/composer.json 644 589ca8b2f934ff98
     - (oracle) F web/app/plugins/duplicate-post/css/duplicate-post-options.css 644 7d4d0a421be4b993
     - (oracle) F web/app/plugins/duplicate-post/css/duplicate-post.css 644 d13f3f0174f75415
     - (oracle) F web/app/plugins/duplicate-post/duplicate-post.php 644 6950a5ad90772594
     - (oracle) F web/app/plugins/duplicate-post/duplicate_post_yoast_icon-125x125.png 644 537f6a71473cbb7b
     - (oracle) F web/app/plugins/duplicate-post/gpl-2.0.txt 644 8177f97513213526
     - (oracle) F web/app/plugins/duplicate-post/js/dist/duplicate-post-edit.js 644 4c33223a69e90963
     - (oracle) F web/app/plugins/duplicate-post/js/dist/duplicate-post-elementor.js 644 8dc1bcfee07158ea
     - (oracle) F web/app/plugins/duplicate-post/js/dist/duplicate-post-options.js 644 efbaf07826b6daf9
     - (oracle) F web/app/plugins/duplicate-post/js/dist/duplicate-post-quick-edit.js 644 6fbd1975c73c6546
     - (oracle) F web/app/plugins/duplicate-post/js/dist/duplicate-post-strings.js 644 b8e0acaef938079c
     - (oracle) F web/app/plugins/duplicate-post/options.php 644 f990001431f94b79
     - (oracle) F web/app/plugins/duplicate-post/readme.txt 644 57cb01cbaf87ec13
     - (oracle) F web/app/plugins/duplicate-post/src/admin/options-form-generator.php 644 a42076f340e745b7
     - (oracle) F web/app/plugins/duplicate-post/src/admin/options-inputs.php 644 628ff1c1a7ab282b
     - (oracle) F web/app/plugins/duplicate-post/src/admin/options-page.php 644 76adcb103a0259b6
     - (oracle) F web/app/plugins/duplicate-post/src/admin/options.php 644 25591974879d3017
     - (oracle) F web/app/plugins/duplicate-post/src/admin/views/options.php 644 feae98a44f74e9b0
     - (oracle) F web/app/plugins/duplicate-post/src/duplicate-post.php 644 c0f28df3642ae2a8
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/bulk-handler.php 644 dca1e1641c5c2d63
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/check-changes-handler.php 644 08d0d0675b85e21c
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/handler.php 644 31066e20c2eab2cc
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/link-handler.php 644 87ebc98074369bd3
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/rest-api-handler.php 644 8e5f216a5c27a714
     - (oracle) F web/app/plugins/duplicate-post/src/handlers/save-post-handler.php 644 520cd1acc3901fb8
     - (oracle) F web/app/plugins/duplicate-post/src/permissions-helper.php 644 dffab2c33869d523
     - (oracle) F web/app/plugins/duplicate-post/src/post-duplicator.php 644 dcf50c4f222d248b
     - (oracle) F web/app/plugins/duplicate-post/src/post-republisher.php 644 bd8c2bd1ffc7d10f
     - (oracle) F web/app/plugins/duplicate-post/src/revisions-migrator.php 644 41469c55455224ad
     - (oracle) F web/app/plugins/duplicate-post/src/ui/admin-bar.php 644 a5e6c9b147d7544d
     - (oracle) F web/app/plugins/duplicate-post/src/ui/asset-manager.php 644 e65c1027782c9c8a
     - (oracle) F web/app/plugins/duplicate-post/src/ui/block-editor.php 644 b200066824befe52
     - (oracle) F web/app/plugins/duplicate-post/src/ui/bulk-actions.php 644 495fedf574315c60
     - (oracle) F web/app/plugins/duplicate-post/src/ui/classic-editor.php 644 0f7185b556261558
     - (oracle) F web/app/plugins/duplicate-post/src/ui/column.php 644 acefe77f40b2d33b
     - (oracle) F web/app/plugins/duplicate-post/src/ui/link-builder.php 644 04cdcbd38d3916fc
     - (oracle) F web/app/plugins/duplicate-post/src/ui/metabox.php 644 30283347a6f37797
     - (oracle) F web/app/plugins/duplicate-post/src/ui/newsletter.php 644 fd6ad9ab8e58c853
     - (oracle) F web/app/plugins/duplicate-post/src/ui/post-states.php 644 d2fcde0fa583edf8
     - (oracle) F web/app/plugins/duplicate-post/src/ui/row-actions.php 644 d2d64ca16cd4e40d
     - (oracle) F web/app/plugins/duplicate-post/src/ui/user-interface.php 644 bc043e2fd634dbcd
     - (oracle) F web/app/plugins/duplicate-post/src/utils.php 644 d18e345f199b2f58
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/bulk-actions-watcher.php 644 0eafefa50e78e0c3
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/copied-post-watcher.php 644 e6a71903778c8433
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/link-actions-watcher.php 644 86733af1e8530492
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/original-post-watcher.php 644 1a52e36b324dea4e
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/republished-post-watcher.php 644 052e5c79b17f587d
     - (oracle) F web/app/plugins/duplicate-post/src/watchers/watchers.php 644 5bd4efc0dc978795
run `sh -c find web -maxdepth 4 -name '*.php' | sort | head -20`:
     - (oracle) web/app/mu-plugins/soil/soil.php
     - (oracle) web/app/plugins/duplicate-post/admin-functions.php
     - (oracle) web/app/plugins/duplicate-post/common-functions.php
     - (oracle) web/app/plugins/duplicate-post/duplicate-post.php
     - (oracle) web/app/plugins/duplicate-post/options.php
     - (oracle) exit=0
```

## plugin-phpcs / install

```
exit code: oracle=0 go=1
stdout:
     - (oracle) PHP CodeSniffer Config installed_paths set to ../../phpcompatibility/php-compatibility
stderr:
     + (go) In <location>:
     + (go)                                                
     + (go)   Class "Seld\Signal\SignalHandler" not found  
     + (go)                                                
     + (go) 
     + (go) install [--prefer-source] [--prefer-dist] [--prefer-install PREFER-INSTALL] [--dry-run] [--download-only] [--dev] [--no-suggest] [--no-dev] [--no-security-blocking] [--no-blocking] [--no-autoloader] [--no-progress] [--no-install] [--audit] [--audit-format AUDIT-FORMAT] [-v|vv|vvv|--verbose] [-o|--optimize-autoloader] [-a|--classmap-authoritative] [--strict-psr-autoloader] [--apcu-autoloader] [--apcu-autoloader-prefix APCU-AUTOLOADER-PREFIX] [--ignore-platform-req IGNORE-PLATFORM-REQ] [--ignore-platform-reqs] [--] [<packages>...]
     + (go) 
     + (go) 
project tree:
     - (oracle) F vendor/squizlabs/php_codesniffer/CodeSniffer.conf 644 69d9574f2b671500
run `php vendor/bin/phpcs --config-show`:
     + (go) Using config file: 
     - (oracle) Using config file: /work/project/vendor/squizlabs/php_codesniffer/CodeSniffer.conf
     - (oracle) installed_paths: ../../phpcompatibility/php-compatibility
run `php vendor/bin/phpcs -i`:
     + (go) The installed coding standards are MySource, PEAR, PSR1, PSR2, PSR12, Squiz and Zend
     - (oracle) The installed coding standards are MySource, PEAR, PSR1, PSR2, PSR12, Squiz, Zend and PHPCompatibility
```
