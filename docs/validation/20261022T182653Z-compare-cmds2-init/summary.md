# Composer vs gocomposer comparison

- gocomposer_sha256: `ec6e5725ad13205d6a8dbcddcb51c4b12a6ae904a2c3044b44736aa75984149a`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| init-full | MATCH | 0 | 0 |
| init-full-reordered | MATCH | 0 | 0 |
| init-require | MATCH | 0 | 0 |
| init-require-pair-args | MATCH | 0 | 0 |
| init-require-no-version | MATCH | 1 | 1 |
| init-require-dev-no-version | MATCH | 1 | 1 |
| init-autoload | DIFF | 0 | 0 |
| init-autoload-nested | DIFF | 0 | 0 |
| init-autoload-with-deps | MATCH | 0 | 0 |
| init-autoload-dev-deps | MATCH | 0 | 0 |
| init-git-defaults | MATCH | 0 | 0 |
| init-git-defaults-description | MATCH | 0 | 0 |
| init-env-defaults | MATCH | 0 | 0 |
| init-env-overrides-git | MATCH | 0 | 0 |
| init-user-env | MATCH | 0 | 0 |
| init-no-defaults | MATCH | 0 | 0 |
| init-camel-dir | MATCH | 0 | 0 |
| init-invalid-name | MATCH | 1 | 1 |
| init-invalid-name-noslash | MATCH | 1 | 1 |
| init-invalid-author | MATCH | 1 | 1 |
| init-invalid-email | MATCH | 1 | 1 |
| init-empty-author | MATCH | 1 | 1 |
| init-author-name-only | MATCH | 0 | 0 |
| init-existing-json | MATCH | 0 | 0 |
| init-existing-git | ERROR | 0 | 0 |
| init-repository | MATCH | 0 | 0 |
| init-repository-invalid | MATCH | 1 | 1 |
| init-repository-bad-json | MATCH | 1 | 1 |
| init-repository-composer-file | MATCH | 0 | 0 |
| init-repository-filesystem-file | MATCH | 1 | 1 |
| init-repository-missing-file | DIFF | 1 | 1 |
| init-bad-stability | MATCH | 1 | 1 |
| init-bad-type-schema | MATCH | 1 | 1 |
| init-bad-license | MATCH | 0 | 0 |
| init-composer-env | MATCH | 0 | 0 |
| init-composer-env-subdir | MATCH | 0 | 0 |
| init-unicode | MATCH | 0 | 0 |

## init-autoload

```
project tree:
+ (go) F vendor/autoload.php 644 c5389ec1f902acbb
- (oracle) F vendor/autoload.php 644 ecdc411359c71703
+ (go) F vendor/composer/autoload_real.php 644 dd861fc699f9becd
+ (go) F vendor/composer/autoload_static.php 644 80a418e2a9892e90
- (oracle) F vendor/composer/autoload_real.php 644 2d2c43bfa06aa006
- (oracle) F vendor/composer/autoload_static.php 644 003faeb59883964d
vendor/composer/autoload_real.php:
+ (go) class ComposerAutoloaderInitdca5a4afd19a5a8942d90c3bee7e5f80
- (oracle) class ComposerAutoloaderInit217d0600175e883bc03987b5d3529562
+ (go)         spl_autoload_register(array('ComposerAutoloaderInitdca5a4afd19a5a8942d90c3bee7e5f80', 'loadClassLoader'), true, true);
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit217d0600175e883bc03987b5d3529562', 'loadClassLoader'), true, true);
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInitdca5a4afd19a5a8942d90c3bee7e5f80', 'loadClassLoader'));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit217d0600175e883bc03987b5d3529562', 'loadClassLoader'));
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInitdca5a4afd19a5a8942d90c3bee7e5f80::getInitializer($loader));
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit217d0600175e883bc03987b5d3529562::getInitializer($loader));
vendor/composer/autoload_static.php:
+ (go) class ComposerStaticInitdca5a4afd19a5a8942d90c3bee7e5f80
- (oracle) class ComposerStaticInit217d0600175e883bc03987b5d3529562
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInitdca5a4afd19a5a8942d90c3bee7e5f80::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInitdca5a4afd19a5a8942d90c3bee7e5f80::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInitdca5a4afd19a5a8942d90c3bee7e5f80::$classMap;
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit217d0600175e883bc03987b5d3529562::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit217d0600175e883bc03987b5d3529562::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit217d0600175e883bc03987b5d3529562::$classMap;
```

## init-autoload-nested

```
project tree:
+ (go) F vendor/autoload.php 644 ee4b236a4a7eae8c
- (oracle) F vendor/autoload.php 644 3d6abc0ab2f79199
+ (go) F vendor/composer/autoload_real.php 644 d74d54f1d4eebabb
+ (go) F vendor/composer/autoload_static.php 644 3da47750f913c8cb
- (oracle) F vendor/composer/autoload_real.php 644 840e3d1dcca6c00d
- (oracle) F vendor/composer/autoload_static.php 644 b5dcbef1489ce26b
vendor/composer/autoload_real.php:
+ (go) class ComposerAutoloaderInitaaf099c69a1a2232674d3b61f01ebfaa
- (oracle) class ComposerAutoloaderInit89660340584e6d27f1d8d9b620e79a96
+ (go)         spl_autoload_register(array('ComposerAutoloaderInitaaf099c69a1a2232674d3b61f01ebfaa', 'loadClassLoader'), true, true);
- (oracle)         spl_autoload_register(array('ComposerAutoloaderInit89660340584e6d27f1d8d9b620e79a96', 'loadClassLoader'), true, true);
+ (go)         spl_autoload_unregister(array('ComposerAutoloaderInitaaf099c69a1a2232674d3b61f01ebfaa', 'loadClassLoader'));
- (oracle)         spl_autoload_unregister(array('ComposerAutoloaderInit89660340584e6d27f1d8d9b620e79a96', 'loadClassLoader'));
+ (go)         call_user_func(\Composer\Autoload\ComposerStaticInitaaf099c69a1a2232674d3b61f01ebfaa::getInitializer($loader));
- (oracle)         call_user_func(\Composer\Autoload\ComposerStaticInit89660340584e6d27f1d8d9b620e79a96::getInitializer($loader));
vendor/composer/autoload_static.php:
+ (go) class ComposerStaticInitaaf099c69a1a2232674d3b61f01ebfaa
- (oracle) class ComposerStaticInit89660340584e6d27f1d8d9b620e79a96
+ (go)             $loader->prefixLengthsPsr4 = ComposerStaticInitaaf099c69a1a2232674d3b61f01ebfaa::$prefixLengthsPsr4;
+ (go)             $loader->prefixDirsPsr4 = ComposerStaticInitaaf099c69a1a2232674d3b61f01ebfaa::$prefixDirsPsr4;
+ (go)             $loader->classMap = ComposerStaticInitaaf099c69a1a2232674d3b61f01ebfaa::$classMap;
- (oracle)             $loader->prefixLengthsPsr4 = ComposerStaticInit89660340584e6d27f1d8d9b620e79a96::$prefixLengthsPsr4;
- (oracle)             $loader->prefixDirsPsr4 = ComposerStaticInit89660340584e6d27f1d8d9b620e79a96::$prefixDirsPsr4;
- (oracle)             $loader->classMap = ComposerStaticInit89660340584e6d27f1d8d9b620e79a96::$classMap;
```

## init-existing-git

```
podman run failed: exit status 128

```

## init-repository-missing-file

```
stderr:
+ (go)                                             
+ (go)   Could not read missing.json               
+ (go)                                             
+ (go)   The file "missing.json" is not readable.  
+ (go)                                             
- (oracle)                                                                                
- (oracle)   The "missing.json" file could not be downloaded: Failed to open stream: No   
- (oracle)   such file or directory                                                       
- (oracle)                                                                                
- (oracle) In <location>:
- (oracle)                                                                                
- (oracle)   The "missing.json" file could not be downloaded: Failed to open stream: No   
- (oracle)   such file or directory                                                       
- (oracle)                                                                                
- (oracle) 
```
