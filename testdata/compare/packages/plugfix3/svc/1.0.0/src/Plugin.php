<?php
namespace Plugfix3\Svc;

use Composer\Composer;
use Composer\IO\IOInterface;
use Composer\Plugin\Capable;
use Composer\Plugin\PluginInterface;

/**
 * Replaces Composer services when activated, depending on PLUGFIX3_SVC (a
 * comma separated list of modes):
 *   zip      the zip downloader logs its calls (LoggingZip)
 *   extract  the zip downloader overrides extract() (ExtractingZip)
 *   wrap     the zip downloader is wrapped by a DownloaderInterface
 *            implementation (WrappingDownloader)
 *   autoload the autoload generator logs its dumps (LoggingGenerator)
 *   repo     a repository type of the plugin provides plugfix3/mem-lib
 *            (MemRepository, prepended); repo-filter: with an exclude filter
 *   overlay  a package repository whose packages the plugin changed
 *   archiver an archiver for the "txt" format is added (TxtArchiver)
 *   archive-manager  the archive manager is replaced by one with the txt
 *            archiver only
 * The plugfix3:svc command exercises the services directly.
 */
class Plugin implements PluginInterface, Capable
{
    public function activate(Composer $composer, IOInterface $io)
    {
        $modes = array_filter(explode(',', (string) getenv('PLUGFIX3_SVC')));
        $dm = $composer->getDownloadManager();
        $loop = $composer->getLoop();
        $process = $loop->getProcessExecutor();
        $args = array($io, $composer->getConfig(), $loop->getHttpDownloader(), $composer->getEventDispatcher(), null, new \Composer\Util\Filesystem($process), $process);
        foreach ($modes as $mode) {
            switch ($mode) {
                case 'zip':
                    $dm->setDownloader('zip', new LoggingZip(...$args));
                    break;
                case 'extract':
                    $dm->setDownloader('ZIP', new ExtractingZip(...$args));
                    break;
                case 'wrap':
                    $dm->setDownloader('zip', new WrappingDownloader($io, $dm->getDownloader('zip')));
                    break;
                case 'autoload':
                    $composer->setAutoloadGenerator(new LoggingGenerator($composer->getEventDispatcher(), $io));
                    $io->write('[svc] autoload: generator is ' . get_class($composer->getAutoloadGenerator()));
                    continue 2;
                case 'repo':
                case 'repo-filter':
                    $rm = $composer->getRepositoryManager();
                    $rm->setRepositoryClass('plugfix3-mem', 'Plugfix3\Svc\MemRepository');
                    $config = array('name' => 'one', 'packages' => array(
                        array('name' => 'plugfix3/mem-lib', 'version' => '1.0.0', 'dist' => array('type' => 'zip', 'url' => 'http://127.0.0.1:8899/dists/acme-log-1.0.0.zip'), 'require' => array('acme/util' => '^1.0'), 'autoload' => array('psr-4' => array('Acme\\Log\\' => 'src/'))),
                        array('name' => 'plugfix3/mem-lib', 'version' => '1.1.0', 'dist' => array('type' => 'zip', 'url' => 'http://127.0.0.1:8899/dists/acme-log-1.1.0.zip')),
                        array('name' => 'plugfix3/mem-other', 'version' => '2.0.0', 'type' => 'metapackage'),
                    ));
                    if ($mode === 'repo-filter') {
                        $config['exclude'] = array('plugfix3/mem-other');
                    }
                    $repo = $rm->createRepository('plugfix3-mem', $config, 'mem-one');
                    $rm->prependRepository($repo);
                    $io->write('[svc] ' . $mode . ': added ' . get_class($repo) . ' ' . $repo->getRepoName() . ' repositories=' . count($rm->getRepositories()));
                    continue 2;
                case 'archive-manager':
                    $am = new \Composer\Package\Archiver\ArchiveManager($composer->getDownloadManager(), $composer->getLoop());
                    $am->addArchiver(new TxtArchiver());
                    $composer->setArchiveManager($am);
                    $io->write('[svc] archive-manager: replaced, same=' . var_export($composer->getArchiveManager() === $am, true));
                    continue 2;
                case 'archiver':
                    $composer->getArchiveManager()->addArchiver(new TxtArchiver());
                    $io->write('[svc] archiver: txt archiver added to ' . get_class($composer->getArchiveManager()));
                    continue 2;
                case 'overlay':
                    $rm = $composer->getRepositoryManager();
                    $repo = $rm->createRepository('package', array('package' => array('name' => 'plugfix3/listed', 'version' => '1.0.0', 'type' => 'metapackage')));
                    $added = new \Composer\Package\CompletePackage('plugfix3/added', '1.0.0.0', '1.0.0');
                    $added->setType('metapackage');
                    $added->setRequires(array('acme/util' => new \Composer\Package\Link('plugfix3/added', 'acme/util', new \Composer\Semver\Constraint\MatchAllConstraint(), \Composer\Package\Link::TYPE_REQUIRE, '*')));
                    $repo->addPackage($added);
                    $repo->removePackage($repo->findPackage('plugfix3/listed', '*'));
                    $rm->addRepository($repo);
                    $names = array();
                    foreach ($repo->getPackages() as $p) {
                        $names[] = $p->getPrettyName();
                    }
                    $io->write('[svc] overlay: ' . get_class($repo) . ' packages=' . implode(',', $names) . ' count=' . count($repo));
                    continue 2;
                default:
                    continue 2;
            }
            $io->write('[svc] ' . $mode . ': zip downloader is ' . get_class($dm->getDownloader('zip')) . ' type=' . $dm->getDownloaderType($dm->getDownloader('zip')));
        }
    }

    public function deactivate(Composer $composer, IOInterface $io)
    {
    }

    public function uninstall(Composer $composer, IOInterface $io)
    {
    }

    public function getCapabilities()
    {
        return array('Composer\Plugin\Capability\CommandProvider' => 'Plugfix3\Svc\Provider');
    }
}
