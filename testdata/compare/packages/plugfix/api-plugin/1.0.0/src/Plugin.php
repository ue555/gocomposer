<?php
namespace Plugfix\Api;

use Composer\Composer;
use Composer\IO\IOInterface;
use Composer\Plugin\Capable;
use Composer\Plugin\PluginInterface;

class Plugin implements PluginInterface, Capable
{
    public function activate(Composer $composer, IOInterface $io)
    {
        if (getenv('PLUGFIX_API_ADDREPO')) {
            // a repository added by the plugin takes part in the resolution
            $rm = $composer->getRepositoryManager();
            $repo = $rm->createRepository('package', array('package' => array(
                array('name' => 'plugfix/virtual', 'version' => '1.0.0', 'type' => 'metapackage', 'require' => array('acme/log' => '^1.0')),
                array('name' => 'plugfix/virtual', 'version' => '1.2.0', 'type' => 'metapackage', 'require' => array('acme/util' => '^1.0')),
            )));
            $rm->prependRepository($repo);
            $io->write('[api] added ' . $repo->getRepoName());
        }
        if (getenv('PLUGFIX_API_DOWNLOADER')) {
            $loop = $composer->getLoop();
            $process = $loop->getProcessExecutor();
            $composer->getDownloadManager()->setDownloader('plugfixzip', new LoggingZipDownloader($io, $composer->getConfig(), $loop->getHttpDownloader(), $composer->getEventDispatcher(), null, new \Composer\Util\Filesystem($process), $process));
            $io->write('[api] registered the plugfixzip downloader');
        }
        if (getenv('PLUGFIX_API_MERGE')) {
            $composer->getConfig()->merge(array('config' => array('sort-packages' => true, 'process-timeout' => 77)));
            $io->write('[api] merged config process-timeout=' . var_export($composer->getConfig()->get('process-timeout'), true));
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
        return array('Composer\Plugin\Capability\CommandProvider' => 'Plugfix\Api\Provider');
    }
}
