<?php
namespace Plugfix\CustomInstallers;

use Composer\Composer;
use Composer\IO\IOInterface;
use Composer\Plugin\PluginInterface;

class Plugin implements PluginInterface
{
    public function activate(Composer $composer, IOInterface $io)
    {
        $im = $composer->getInstallationManager();
        $im->addInstaller(new HookInstaller($io, $composer, 'plugfix-lib'));
        $im->addInstaller(new MetaInstaller($io));
        $im->addInstaller(new AssetInstaller($io, $composer));
    }

    public function deactivate(Composer $composer, IOInterface $io)
    {
    }

    public function uninstall(Composer $composer, IOInterface $io)
    {
    }
}
