<?php
namespace Plugfix3\Svc;

use Composer\Command\BaseCommand;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class SvcCommand extends BaseCommand
{
    protected function configure()
    {
        $this->setName('plugfix3:svc')
            ->setDescription('Exercises Composer services replaced or extended by plugins')
            ->addArgument('what', InputArgument::REQUIRED, 'dl-internals|autoload-helpers|autoload-target|autoload-filter|installer-fixed|installer-audit|installer-filter|lock-tx-none|factory|config|events|schema|archive|http-tls|set-services');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $composer = $this->requireComposer();
        $io = $this->getIO();
        switch ($input->getArgument('what')) {
            case 'dl-internals':
                $loop = $composer->getLoop();
                $dm = $composer->getDownloadManager();
                $process = $loop->getProcessExecutor();
                $zip = new InternalsZip($io, $composer->getConfig(), $loop->getHttpDownloader(), $composer->getEventDispatcher(), null, new \Composer\Util\Filesystem($process), $process);
                $local = $composer->getRepositoryManager()->getLocalRepository();
                $im = $composer->getInstallationManager();
                $log = $local->findPackage('acme/log', '*');
                $path = $im->getInstallPath($log);
                $name = $zip->fileName($log, $path);
                $io->write('fileName ' . (preg_match('{^' . preg_quote(realpath('vendor'), '{') . '/composer/tmp-[0-9a-f]{32}\.zip$}', $name) ? 'vendor/composer/tmp-HASH.zip' : $name));
                $io->write('processUrl ' . $zip->url($log, 'https://github.com/acme/log/archive/v1.zip'));
                $io->write('processUrl ' . $zip->url($log, 'http://127.0.0.1:8899/dists/acme-log-1.0.0.zip'));
                $io->write('changes before ' . var_export($zip->getLocalChanges($log, $path), true));
                file_put_contents($path . '/composer.json', "{}\n");
                unlink($path . '/src/Logger.php');
                file_put_contents($path . '/NEW.txt', "new\n");
                $io->write('changes after ' . var_export($zip->getLocalChanges($log, $path), true));
                $io->write('compare dir left ' . var_export(is_dir($path . '_compare'), true));
                $git = $dm->getDownloader('git');
                $io->write('git changes ' . var_export($git->getLocalChanges($log, $path), true));
                $io->write('git reference ' . var_export($git->getVcsReference($log, $path), true));
                break;
            case 'autoload-helpers':
                $gen = $composer->getAutoloadGenerator();
                $im = $composer->getInstallationManager();
                $root = $composer->getPackage();
                $local = $composer->getRepositoryManager()->getLocalRepository();
                $packageMap = $gen->buildPackageMap($im, $root, $local->getCanonicalPackages());
                foreach ($packageMap as $item) {
                    $io->write('map ' . $item[0]->getName() . ' => ' . var_export($item[1] === '' ? '' : str_replace(getcwd(), 'CWD', (string) $item[1]), true));
                }
                $autoloads = $gen->parseAutoloads($packageMap, $root);
                foreach (array('psr-0', 'psr-4', 'classmap', 'files', 'exclude-from-classmap') as $type) {
                    $io->write($type . ' ' . str_replace(getcwd(), 'CWD', json_encode($autoloads[$type], JSON_UNESCAPED_SLASHES)));
                }
                $filtered = $gen->parseAutoloads($packageMap, $root, true);
                $io->write('filtered psr-4 ' . implode(',', array_keys($filtered['psr-4'])));
                $loader = $gen->createLoader($autoloads, $composer->getConfig()->get('vendor-dir'));
                $prefixes = array_keys($loader->getPrefixesPsr4());
                sort($prefixes);
                $io->write('loader psr-4 ' . implode(',', $prefixes) . ' classmap=' . count($loader->getClassMap()));
                $io->write('findFile ' . str_replace(getcwd(), 'CWD', var_export($loader->findFile('Acme\\Log\\Logger'), true)));
                try {
                    $bad = new \Composer\Package\Package('plugfix3/bad', '1.0.0.0', '1.0.0');
                    $bad->setAutoload(array('psr-4' => array('Bad' => 'src/')));
                    $gen->buildPackageMap($im, $root, array($bad));
                } catch (\InvalidArgumentException $e) {
                    $io->write('invalid ' . $e->getMessage());
                }
                break;
            case 'autoload-target':
                $gen = $composer->getAutoloadGenerator();
                $map = $gen->dump($composer->getConfig(), $composer->getRepositoryManager()->getLocalRepository(), $composer->getPackage(), $composer->getInstallationManager(), 'plugfix3-al', true);
                $files = scandir('vendor/plugfix3-al');
                $io->write('target files ' . implode(',', array_diff($files, array('.', '..'))) . ' classes=' . count($map->getMap()));
                $io->write('autoload.php ' . (strpos(file_get_contents('vendor/autoload.php'), "/plugfix3-al/autoload_real.php") !== false ? 'uses plugfix3-al' : 'other'));
                break;
            case 'autoload-filter':
                $gen = $composer->getAutoloadGenerator();
                $gen->setPlatformRequirementFilter(new IgnoreExtFilter());
                $gen->dump($composer->getConfig(), $composer->getRepositoryManager()->getLocalRepository(), $composer->getPackage(), $composer->getInstallationManager(), 'composer');
                $check = file_get_contents('vendor/composer/platform_check.php');
                $io->write('platform_check json=' . var_export(strpos($check, "'json'") !== false, true) . ' ctype=' . var_export(strpos($check, "'ctype'") !== false, true));
                break;
            case 'installer-fixed':
                $installer = \Composer\Installer::create($io, $composer);
                $fixed = new \Composer\Repository\ArrayRepository();
                $util = new \Composer\Package\CompletePackage('acme/util', '1.0.0.0', '1.0.0');
                $util->setDistType('zip');
                $util->setDistUrl('http://127.0.0.1:8899/dists/acme-util-1.0.0.zip');
                $fixed->addPackage($util);
                $installer->setUpdate(true)->setDryRun(true);
                if (!getenv('PLUGFIX3_NOFIXED')) {
                    $installer->setAdditionalFixedRepository($fixed);
                }
                $code = $installer->run();
                $io->write('installer exit=' . $code);
                $tx = $installer->getLockTransaction();
                foreach ($tx->getOperations() as $op) {
                    $io->write('tx op ' . $op->show(false));
                }
                $names = array();
                foreach ($tx->getNewLockPackages(false) as $p) {
                    $names[] = $p->getPrettyName() . '@' . $p->getPrettyVersion();
                }
                $io->write('tx lock ' . implode(',', $names));

                return $code;
            case 'installer-audit':
                $installer = \Composer\Installer::create($io, $composer);
                $installer->setAuditConfig(new \Composer\Advisory\AuditConfig(false));
                $code = $installer->run();
                $io->write('installer (audit off) exit=' . $code);
                $installer = \Composer\Installer::create($io, $composer);
                $installer->setAuditConfig(new \Composer\Advisory\AuditConfig(false))->setAuditFormat('plain');
                $code = $installer->run();
                $io->write('installer (audit config invalidated) exit=' . $code);

                return $code;
            case 'installer-filter':
                $installer = \Composer\Installer::create($io, $composer);
                $installer->setUpdate(true)->setPlatformRequirementFilter(new IgnoreExtFilter());
                $code = $installer->run();
                $io->write('installer exit=' . $code);

                return $code;
            case 'factory':
                $other = \Composer\Factory::create($io, array(
                    'name' => 'test/array',
                    'require' => array('acme/util' => '^1.0'),
                    'repositories' => array(array('type' => 'composer', 'url' => 'http://127.0.0.1:8899'), array('packagist.org' => false)),
                    'config' => array('secure-http' => false, 'vendor-dir' => 'libs'),
                ));
                $io->write('array root=' . $other->getPackage()->getName() . ' requires=' . implode(',', array_keys($other->getPackage()->getRequires())) . ' locked=' . var_export($other->getLocker()->isLocked(), true) . ' vendor-dir=' . $other->getConfig()->get('vendor-dir', \Composer\Config::RELATIVE_PATHS) . ' global=' . var_export($other->isGlobal(), true));
                $p = $other->getRepositoryManager()->findPackage('acme/util', '^1.0');
                $io->write('array findPackage ' . ($p ? $p->getPrettyVersion() : 'null'));
                $global = \Composer\Factory::createGlobal($io);
                $io->write('global (none) ' . var_export($global, true));
                $home = $composer->getConfig()->get('home');
                file_put_contents($home . '/composer.json', json_encode(array('name' => 'test/global-home', 'require' => array('acme/log' => '^1.0'))));
                $global = \Composer\Factory::createGlobal($io, true);
                $io->write('global root=' . $global->getPackage()->getName() . ' isGlobal=' . var_export($global->isGlobal(), true) . ' requires=' . implode(',', array_keys($global->getPackage()->getRequires())));
                $factory = new \Composer\Factory();
                $sub = $factory->createComposer($io, 'sub/composer.json', true);
                $io->write('createComposer root=' . $sub->getPackage()->getName() . ' requires=' . implode(',', array_keys($sub->getPackage()->getRequires())));
                break;
            case 'config':
                $config = $composer->getConfig();
                $rel = function ($path) {
                    return str_replace(array(getcwd(), getenv('COMPOSER_HOME')), array('CWD', 'HOME'), $path);
                };
                $io->write('source ' . $rel($config->getConfigSource()->getName()) . ' auth ' . $rel($config->getAuthConfigSource()->getName()) . ' local auth ' . var_export($config->getLocalAuthConfigSource(), true));
                $config->getConfigSource()->addConfigSetting('sort-packages', true);
                $config->getConfigSource()->addProperty('extra.plugfix3', 'written');
                $io->write('composer.json sort-packages written=' . var_export(strpos(file_get_contents('composer.json'), '"sort-packages": true') !== false, true));
                $own = new \Composer\Config(false, getcwd());
                $io->write('own vendor-dir=' . $own->get('vendor-dir') . ' process-timeout=' . var_export($own->get('process-timeout'), true));
                $own->merge(array('config' => array('vendor-dir' => 'lib', 'process-timeout' => 5, 'bin-dir' => '{$vendor-dir}/tools')), 'plugfix3');
                $io->write('own merged vendor-dir=' . $rel($own->get('vendor-dir')) . ' relative=' . $own->get('vendor-dir', \Composer\Config::RELATIVE_PATHS) . ' bin-dir=' . $rel($own->get('bin-dir')) . ' process-timeout=' . var_export($own->get('process-timeout'), true) . ' source=' . $own->getSourceOfValue('process-timeout'));
                $io->write('main process-timeout=' . var_export($config->get('process-timeout'), true));
                try {
                    $own->getConfigSource();
                } catch (\TypeError $e) {
                    $io->write('own config source TypeError');
                }
                $created = \Composer\Factory::createConfig($io, getcwd());
                $io->write('createConfig home=' . $rel($created->get('home')) . ' source ' . $rel($created->getConfigSource()->getName()) . ' auth ' . $rel($created->getAuthConfigSource()->getName()));
                try {
                    $own->prohibitUrlByConfig('http://example.org/x.json', $io);
                } catch (\Composer\Downloader\TransportException $e) {
                    $io->write('own prohibit ' . $e->getMessage());
                }
                $config->prohibitUrlByConfig('http://127.0.0.1:8899/packages.json', $io);
                $config->prohibitUrlByConfig('http://localhost:8899/packages.json', $io);
                $config->prohibitUrlByConfig('http://localhost:8899/other.json', $io);
                $own->merge(array('config' => array('secure-http' => false)));
                $own->prohibitUrlByConfig('http://example.org/x.json', $io);
                $own->prohibitUrlByConfig('http://example.org/y.json', $io, array('ssl' => array('verify_peer' => false)));
                break;
            case 'events':
                $other = \Composer\Factory::create($io, 'sub/composer.json', true);
                $dispatcher = $other->getEventDispatcher();
                $io->write('listeners post-install-cmd=' . var_export($dispatcher->hasEventListeners(new \Composer\Script\Event('post-install-cmd', $other, $io)), true) . ' none=' . var_export($dispatcher->hasEventListeners(new \Composer\EventDispatcher\Event('plugfix3-none')), true));
                $io->write('dispatchScript=' . $dispatcher->dispatchScript('post-install-cmd', true));
                $io->write('dispatch=' . $dispatcher->dispatch('plugfix3-custom'));
                $dispatcher->setRunScripts(false);
                $io->write('dispatchScript without scripts=' . $dispatcher->dispatchScript('post-install-cmd'));
                $scriptless = \Composer\Factory::create($io, 'sub/composer.json', true, true);
                $io->write('disabled scripts=' . $scriptless->getEventDispatcher()->dispatchScript('post-install-cmd'));
                break;
            case 'schema':
                $json = new \Composer\Json\JsonFile('composer.json');
                foreach (array('plug3-schema.json', 'plug3-schema-ok.json') as $schema) {
                    try {
                        $json->validateSchema(\Composer\Json\JsonFile::LAX_SCHEMA, getcwd() . '/' . $schema);
                        $io->write($schema . ' valid');
                    } catch (\Composer\Json\JsonValidationException $e) {
                        $io->write($schema . ' ' . str_replace(getcwd(), 'CWD', $e->getMessage()) . ': ' . implode(' | ', $e->getErrors()));
                    }
                }
                try {
                    \Composer\Json\JsonFile::validateJsonSchema('data', (object) array('plug3' => 5), \Composer\Json\JsonFile::STRICT_SCHEMA, getcwd() . '/plug3-schema.json');
                } catch (\Composer\Json\JsonValidationException $e) {
                    $io->write('data ' . implode(' | ', $e->getErrors()));
                }
                break;
            case 'archive':
                $am = $composer->getArchiveManager();
                $io->write('archive manager ' . get_class($am));
                $log = $composer->getRepositoryManager()->getLocalRepository()->findPackage('acme/log', '*');
                $io->write('filename ' . $am->getPackageFilename($log));
                try {
                    $am->archive($log, 'txt', 'dist');
                } catch (\RuntimeException $e) {
                    $io->write('error ' . $e->getMessage());
                }
                $am->addArchiver(new TxtArchiver());
                $path = $am->archive($log, 'txt', 'dist');
                $io->write('txt ' . basename($path) . "\n" . file_get_contents($path));
                $path = $am->archive($composer->getPackage(), 'txt', 'dist', 'root-files');
                $io->write('root txt ' . basename($path) . "\n" . file_get_contents($path));
                $path = $am->archive($log, 'tar', 'dist');
                $io->write('tar ' . basename($path) . ' exists=' . var_export(is_file($path), true));
                $am->setOverwriteFiles(false);
                file_put_contents($path, 'kept');
                $path = $am->archive($log, 'tar', 'dist');
                $io->write('not overwritten=' . var_export(file_get_contents($path) === 'kept', true));
                unlink($path);
                break;
            case 'http-tls':
                $http = $composer->getLoop()->getHttpDownloader();
                $api = '/fixtures/git/_api/';
                $try = function ($label, $url, $options) use ($http, $io) {
                    try {
                        $response = $http->get($url, $options);
                        $io->write($label . ' ' . $response->getStatusCode() . ' ' . trim($response->getBody()));
                    } catch (\Composer\Downloader\TransportException $e) {
                        $io->write($label . ' ' . get_class($e) . ' status=' . var_export($e->getStatusCode(), true));
                    }
                };
                $try('cafile', 'https://127.0.0.1/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem')));
                $try('capath', 'https://127.0.0.1/hello.json', array('ssl' => array('capath' => $api)));
                $try('no-ca', 'https://127.0.0.1/hello.json', array());
                $try('verify_peer=false', 'https://127.0.0.1/hello.json', array('ssl' => array('verify_peer' => false)));
                $try('ignored-options', 'https://127.0.0.1/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem', 'allow_self_signed' => true, 'peer_fingerprint' => 'x'), 'http' => array('proxy' => 'tcp://127.0.0.1:9', 'request_fulluri' => true)));
                $fp = openssl_x509_fingerprint(file_get_contents($api . 'server.pem'), 'sha1');
                $try('self-signed-fingerprint', 'https://127.0.0.1/hello.json', array('ssl' => array('allow_self_signed' => true, 'cafile' => $api . 'ca.pem', 'peer_fingerprint' => $fp)));
                $try('self-signed-bad-fingerprint', 'https://127.0.0.1/hello.json', array('ssl' => array('allow_self_signed' => true, 'cafile' => $api . 'ca.pem', 'peer_fingerprint' => str_repeat('0', 40))));
                $try('curl-ignores-proxy', 'https://127.0.0.1/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem'), 'http' => array('proxy' => 'tcp://127.0.0.1:9')));
                $try('client-cert', 'https://127.0.0.1:8444/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem', 'local_cert' => $api . 'client.pem', 'local_pk' => $api . 'client.key')));
                $try('client-cert-passphrase', 'https://127.0.0.1:8444/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem', 'local_cert' => $api . 'client.pem', 'local_pk' => $api . 'client-enc.key', 'passphrase' => 'secret')));
                $try('no-client-cert', 'https://127.0.0.1:8444/hello.json', array('ssl' => array('cafile' => $api . 'ca.pem')));
                $target = sys_get_temp_dir() . '/plugfix3-tls.json';
                $response = $http->copy('https://127.0.0.1/hello.json', $target, array('ssl' => array('cafile' => $api . 'ca.pem')));
                $io->write('copy ' . $response->getStatusCode() . ' ' . trim(file_get_contents($target)));
                @unlink($target);
                break;
            case 'set-services':
                $loop = new \Composer\Util\Loop($composer->getLoop()->getHttpDownloader(), $composer->getLoop()->getProcessExecutor());
                $composer->setLoop($loop);
                $dm = new \Composer\Downloader\DownloadManager($io);
                $dm->setDownloader('plugfix3', new WrappingDownloader($io, $composer->getDownloadManager()->getDownloader('zip')));
                $composer->setDownloadManager($dm);
                $generator = new LoggingGenerator($composer->getEventDispatcher(), $io);
                $composer->setAutoloadGenerator($generator);
                $io->write('loop=' . var_export($composer->getLoop() === $loop, true) . ' dm=' . var_export($composer->getDownloadManager() === $dm, true) . ' generator=' . var_export($composer->getAutoloadGenerator() === $generator, true));
                $io->write('dm downloader ' . get_class($composer->getDownloadManager()->getDownloader('plugfix3')));
                try {
                    $composer->getDownloadManager()->getDownloader('zip');
                } catch (\InvalidArgumentException $e) {
                    $io->write('dm error ' . $e->getMessage());
                }
                $installer = \Composer\Installer::create($io, $composer);
                $code = $installer->run();
                $io->write('installer exit=' . $code);

                return $code;
            case 'lock-tx-none':
                $installer = \Composer\Installer::create($io, $composer);
                $code = $installer->run();
                $io->write('installer exit=' . $code);
                $io->write('lock transaction ' . var_export($installer->getLockTransaction(), true));

                return $code;
        }

        return 0;
    }
}
