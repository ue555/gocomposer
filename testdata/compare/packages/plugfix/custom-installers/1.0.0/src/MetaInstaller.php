<?php
namespace Plugfix\CustomInstallers;

use Composer\Installer\MetapackageInstaller;

/**
 * Installs packages of type plugfix-meta like metapackages (like
 * symfony/flex's SymfonyPackInstaller).
 */
class MetaInstaller extends MetapackageInstaller
{
    public function supports(string $packageType)
    {
        return 'plugfix-meta' === $packageType;
    }
}
