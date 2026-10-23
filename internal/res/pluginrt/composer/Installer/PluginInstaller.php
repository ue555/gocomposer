<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Installer\PluginInstaller, based on
// Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Installer/PluginInstaller.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). Plugin registration on install is driven by gocomposer.

namespace Composer\Installer;

use Composer\IO\IOInterface;
use Composer\PartialComposer;
use Composer\Util\Filesystem;

class PluginInstaller extends LibraryInstaller
{
    public function __construct(IOInterface $io, PartialComposer $composer, ?Filesystem $fs = null, ?BinaryInstaller $binaryInstaller = null)
    {
        parent::__construct($io, $composer, 'composer-plugin', $fs, $binaryInstaller);
    }

    public function supports(string $packageType)
    {
        return $packageType === 'composer-plugin' || $packageType === 'composer-installer';
    }

    public function disablePlugins(): void
    {
    }
}
