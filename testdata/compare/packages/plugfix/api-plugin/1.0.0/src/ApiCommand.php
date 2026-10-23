<?php
namespace Plugfix\Api;

use Composer\Command\BaseCommand;
use Composer\Downloader\TransportException;
use Composer\Factory;
use Composer\Installer;
use Composer\Semver\Constraint\MatchAllConstraint;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class ApiCommand extends BaseCommand
{
    protected function configure()
    {
        $this->setName('plugfix:api')
            ->setDescription('Exercises Composer services')
            ->addArgument('what', InputArgument::REQUIRED, 'http|repos|create-repo|factory|merge|autoload|dm|install|update');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $composer = $this->requireComposer();
        $io = $this->getIO();
        switch ($input->getArgument('what')) {
            case 'http':
                $http = $composer->getLoop()->getHttpDownloader();
                $response = $http->get('http://127.0.0.1:8899/packages.json');
                $data = $response->decodeJson();
                $io->write('get status=' . $response->getStatusCode() . ' type=' . $response->getHeader('content-type') . ' packages=' . count($data['packages']));
                $promises = array();
                $promises[] = $http->add('http://127.0.0.1:8899/v2/p2/acme/log.json')->then(function ($response) use ($io) {
                    $data = $response->decodeJson();
                    $io->write('async log versions=' . count($data['packages']['acme/log']));
                });
                $promises[] = $http->add('http://127.0.0.1:8899/v2/p2/acme/nope.json')->then(null, function ($e) use ($io) {
                    $io->write('async 404 ' . get_class($e) . ' code=' . $e->getCode() . ' status=' . var_export($e instanceof TransportException ? $e->getStatusCode() : null, true) . ' message=' . $e->getMessage());
                });
                $composer->getLoop()->wait($promises);
                $target = sys_get_temp_dir() . '/plugfix-copy.zip';
                $response = $http->copy('http://127.0.0.1:8899/dists/acme-log-1.0.0.zip', $target);
                $io->write('copy status=' . $response->getStatusCode() . ' size=' . filesize($target) . ' sha1=' . sha1_file($target));
                @unlink($target);
                $own = Factory::createHttpDownloader($io, $composer->getConfig());
                try {
                    $own->get('http://127.0.0.1:8899/missing.json');
                } catch (TransportException $e) {
                    $io->write('sync 404 code=' . $e->getCode() . ' status=' . $e->getStatusCode() . ' message=' . $e->getMessage());
                }
                break;
            case 'repos':
                $rm = $composer->getRepositoryManager();
                foreach ($rm->getRepositories() as $repo) {
                    $config = $repo instanceof \Composer\Repository\ConfigurableRepositoryInterface ? $repo->getRepoConfig() : array();
                    $io->write('repo ' . get_class($repo) . ' ' . $repo->getRepoName() . ' config=' . json_encode($config));
                }
                $p = $rm->findPackage('acme/log', '^1.0');
                $io->write('findPackage acme/log ^1.0 = ' . ($p ? $p->getPrettyVersion() . ' ' . get_class($p) . ' repo=' . $p->getRepository()->getRepoName() : 'null'));
                $versions = array();
                foreach ($rm->findPackages('acme/log', '*') as $p) {
                    $versions[] = $p->getPrettyVersion();
                }
                $io->write('findPackages acme/log = ' . implode(',', $versions));
                $io->write('findPackage missing = ' . var_export($rm->findPackage('acme/nope', new MatchAllConstraint()), true));
                $repos = $rm->getRepositories();
                $repo = $repos[0];
                $found = array();
                foreach ($repo->search('log') as $r) {
                    $found[] = $r['name'];
                }
                $io->write('search log = ' . implode(',', $found));
                $loaded = $repo->loadPackages(array('acme/util' => new MatchAllConstraint()), array('stable' => 0), array());
                $lv = array();
                foreach ($loaded['packages'] as $p) {
                    $lv[] = $p->getPrettyVersion();
                }
                $io->write('loadPackages acme/util names=' . implode(',', $loaded['namesFound']) . ' versions=' . implode(',', $lv));
                $io->write('providers log-impl = ' . json_encode($repo->getProviders('acme/log-impl')));
                $io->write('packages = ' . count($repo->getPackages()));
                break;
            case 'create-repo':
                $rm = $composer->getRepositoryManager();
                $repo = $rm->createRepository('package', array('package' => array('name' => 'plugfix/virtual', 'version' => '2.0.0', 'type' => 'metapackage')));
                $io->write('created ' . get_class($repo) . ' ' . $repo->getRepoName());
                $rm->addRepository($repo);
                $p = $rm->findPackage('plugfix/virtual', '*');
                $io->write('found ' . ($p ? $p->getPrettyName() . ' ' . $p->getPrettyVersion() : 'null') . ' repositories=' . count($rm->getRepositories()));
                break;
            case 'factory':
                $other = Factory::create($io, 'sub/composer.json', true);
                $io->write('factory root=' . $other->getPackage()->getName() . ' requires=' . implode(',', array_keys($other->getPackage()->getRequires())));
                $io->write('factory vendor-dir=' . $other->getConfig()->get('vendor-dir', \Composer\Config::RELATIVE_PATHS) . ' locked=' . var_export($other->getLocker()->isLocked(), true));
                $io->write('factory installed=' . count($other->getRepositoryManager()->getLocalRepository()->getPackages()) . ' repositories=' . count($other->getRepositoryManager()->getRepositories()));
                $io->write('main root=' . $composer->getPackage()->getName());
                break;
            case 'errors':
                try {
                    Factory::create($io, 'nope/composer.json');
                } catch (\Exception $e) {
                    $io->write('factory error ' . get_class($e) . ': ' . $e->getMessage());
                }
                try {
                    $composer->getRepositoryManager()->createRepository('nope', array());
                } catch (\Exception $e) {
                    $io->write('createRepository error ' . get_class($e) . ': ' . $e->getMessage());
                }
                try {
                    $composer->getDownloadManager()->getDownloader('nope');
                } catch (\Exception $e) {
                    $io->write('getDownloader error ' . get_class($e) . ': ' . $e->getMessage());
                }
                try {
                    $composer->getLoop()->getHttpDownloader()->get('http://127.0.0.1:1/unreachable.json');
                } catch (\Composer\Downloader\TransportException $e) {
                    $io->write('unreachable error ' . get_class($e) . ' code=' . $e->getCode());
                }
                break;
            case 'merge':
                $config = $composer->getConfig();
                $io->write('before process-timeout=' . var_export($config->get('process-timeout'), true) . ' bin-dir=' . $config->get('bin-dir', \Composer\Config::RELATIVE_PATHS));
                $config->merge(array('config' => array('process-timeout' => 42, 'bin-dir' => 'tools/bin')));
                $io->write('after process-timeout=' . var_export($config->get('process-timeout'), true) . ' bin-dir=' . $config->get('bin-dir', \Composer\Config::RELATIVE_PATHS) . ' source=' . $config->getSourceOfValue('process-timeout'));
                $config->merge(array('repositories' => array(array('type' => 'package', 'package' => array('name' => 'plugfix/merged', 'version' => '1.0.0')))));
                $p = $composer->getRepositoryManager()->findPackage('plugfix/merged', '*');
                $io->write('merged repository package=' . ($p ? $p->getPrettyVersion() : 'null'));
                break;
            case 'autoload':
                $generator = $composer->getAutoloadGenerator();
                $generator->setDevMode(false);
                $generator->setClassMapAuthoritative(true);
                $generator->setRunScripts(true);
                $map = $generator->dump($composer->getConfig(), $composer->getRepositoryManager()->getLocalRepository(), $composer->getPackage(), $composer->getInstallationManager(), 'composer', true);
                $classes = array_keys($map->getMap());
                sort($classes);
                $io->write('classmap ' . count($classes) . ': ' . implode(',', $classes));
                break;
            case 'dm':
                $dm = $composer->getDownloadManager();
                $io->write('zip=' . get_class($dm->getDownloader('zip')) . ' git=' . get_class($dm->getDownloader('git')) . ' path=' . get_class($dm->getDownloader('path')));
                $p = $composer->getRepositoryManager()->findPackage('acme/util', '1.0.0');
                $im = $composer->getInstallationManager();
                $path = $im->getInstallPath($p);
                $p->setInstallationSource('dist');
                $downloader = $dm->getDownloaderForPackage($p);
                $io->write('downloader for ' . $p->getPrettyName() . '=' . get_class($downloader) . ' type=' . $dm->getDownloaderType($downloader));
                $loop = $composer->getLoop();
                $loop->wait(array($dm->download($p, $path)));
                $loop->wait(array($dm->prepare('install', $p, $path)));
                $loop->wait(array($dm->install($p, $path)));
                $loop->wait(array($dm->cleanup('install', $p, $path)));
                $io->write('installed ' . implode(',', array_diff(scandir($path), array('.', '..'))));
                $loop->wait(array($dm->remove($p, $path)));
                $io->write('removed ' . var_export(is_dir($path), true));
                try {
                    $dm->getDownloader('nope');
                } catch (\InvalidArgumentException $e) {
                    $io->write('error ' . $e->getMessage());
                }
                break;
            case 'install':
                $installer = Installer::create($io, $composer);
                $installer->setDevMode(false)->setDumpAutoloader(true)->setOptimizeAutoloader(true);
                $code = $installer->run();
                $io->write('installer exit=' . $code);

                return $code;
            case 'update':
                $installer = Installer::create($io, $composer);
                $installer->setUpdate(true)->setUpdateAllowList(array('acme/log'))->setInstall(false);
                $code = $installer->run();
                $io->write('installer exit=' . $code);

                return $code;
        }

        return 0;
    }
}
