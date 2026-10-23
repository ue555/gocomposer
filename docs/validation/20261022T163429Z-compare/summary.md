# Composer vs gocomposer comparison

- gocomposer_sha256: `cd65130634f94e65b31efa557c48ee6c044582a544ee79e94ad851b939fa05de`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| validate-publish-errors | MATCH | 1 | 1 |
| validate-lax-fail | DIFF | 1 | 1 |
| validate-lax-fail-update | MATCH | 1 | 1 |
| composer-auth-invalid | DIFF | 255 | 0 |
| composer-auth-malformed | MATCH | 1 | 1 |

## validate-lax-fail

```
stderr:
+ (go)                                                                      
+ (go)   "./composer.json" does not match the expected JSON schema:         
+ (go)    - require.acme/log : Array value found, but a string is required  
+ (go)                                                                      
- (oracle)                                                                                
- (oracle)   [TypeError]                                                                  
- (oracle)   Composer\Pcre\Preg::isMatch(): Argument #2 ($subject) must be of type strin  
- (oracle)   g, array given, called in /oracle-source/src/Composer/Util/ConfigValidator.  
- (oracle)   php on line 179                                                              
- (oracle)                                                                                
- (oracle) Exception trace:
- (oracle)   at /oracle-source/vendor/composer/pcre/src/Preg.php:289
- (oracle)  Composer\Pcre\Preg::isMatch() at /oracle-source/src/Composer/Util/ConfigValidator.php:179
- (oracle)  Composer\Util\ConfigValidator->validate() at /oracle-source/src/Composer/Command/ValidateCommand.php:89
- (oracle)  Composer\Command\ValidateCommand->execute() at /oracle-source/vendor/symfony/console/Command/Command.php:298
- (oracle)  Symfony\Component\Console\Command\Command->run() at /oracle-source/vendor/symfony/console/Application.php:1040
- (oracle)  Symfony\Component\Console\Application->doRunCommand() at /oracle-source/vendor/symfony/console/Application.php:301
- (oracle)  Symfony\Component\Console\Application->doRun() at /oracle-source/src/Composer/Console/Application.php:504
- (oracle)  Composer\Console\Application->doRunComposer() at /oracle-source/src/Composer/Console/Application.php:169
- (oracle)  Composer\Console\Application->doRun() at /oracle-source/vendor/symfony/console/Application.php:171
- (oracle)  Symfony\Component\Console\Application->run() at /oracle-source/src/Composer/Console/Application.php:142
- (oracle)  Composer\Console\Application->run() at /oracle-source/bin/composer:113
- (oracle) 
```

## composer-auth-invalid

```
exit code: oracle=255 go=0
stderr:
+ (go) ./composer.json is valid, but with a few warnings
+ (go) See https://getcomposer.org/doc/04-schema.md for details on the schema
+ (go) # General warnings
+ (go) - No license specified, it is recommended to do so. For closed-source software you may use "proprietary" as license.
- (oracle) Fatal error: Uncaught TypeError: Composer\IO\BaseIO::checkAndSetAuthentication(): Argument #2 ($username) must be of type string, int given, called in /oracle-source/src/Composer/IO/BaseIO.php on line 175 and defined in /oracle-source/src/Composer/IO/BaseIO.php:95
- (oracle) Stack trace:
- (oracle) #0 /oracle-source/src/Composer/IO/BaseIO.php(175): Composer\IO\BaseIO->checkAndSetAuthentication('example.org', 5, NULL)
- (oracle) #1 /oracle-source/src/Composer/Factory.php(357): Composer\IO\BaseIO->loadConfiguration(Object(Composer\Config))
- (oracle) #2 /oracle-source/src/Composer/Factory.php(627): Composer\Factory->createComposer(Object(Composer\IO\ConsoleIO), Array, true, '/work/project', true, false)
- (oracle) #3 /oracle-source/src/Composer/Console/Application.php(676): Composer\Factory::create(Object(Composer\IO\ConsoleIO), NULL, true, false)
- (oracle) #4 /oracle-source/src/Composer/Console/Application.php(585): Composer\Console\Application->getComposer(false, true)
- (oracle) #5 /oracle-source/src/Composer/Console/Application.php(532): Composer\Console\Application->hintCommonErrors(Object(ErrorException), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #6 /oracle-source/src/Composer/Console/Application.php(169): Composer\Console\Application->doRunComposer(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput), Object(Composer\IO\ConsoleIO))
- (oracle) #7 /oracle-source/vendor/symfony/console/Application.php(171): Composer\Console\Application->doRun(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #8 /oracle-source/src/Composer/Console/Application.php(142): Symfony\Component\Console\Application->run(Object(Symfony\Component\Console\Input\ArgvInput), Object(Symfony\Component\Console\Output\ConsoleOutput))
- (oracle) #9 /oracle-source/bin/composer(113): Composer\Console\Application->run()
- (oracle) #10 {main}
- (oracle)   thrown in /oracle-source/src/Composer/IO/BaseIO.php on line 95
```
