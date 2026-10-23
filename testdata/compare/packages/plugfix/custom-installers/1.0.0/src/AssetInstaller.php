<?php
namespace Plugfix\CustomInstallers;

use Composer\Composer;
use Composer\Installer\InstallerInterface;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Repository\InstalledRepositoryInterface;

/**
 * An installer written from scratch (not a LibraryInstaller): packages of
 * type plugfix-asset go to public/assets/<name> through the download manager.
 */
class AssetInstaller implements InstallerInterface
{
    private $io;
    private $composer;

    public function __construct(IOInterface $io, Composer $composer)
    {
        $this->io = $io;
        $this->composer = $composer;
    }

    public function supports(string $packageType)
    {
        return 'plugfix-asset' === $packageType;
    }

    public function isInstalled(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        return $repo->hasPackage($package) && is_dir($this->getInstallPath($package));
    }

    public function download(PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->io->write('[asset] download ' . $package->getPrettyName() . ' ' . $package->getPrettyVersion());

        return $this->composer->getDownloadManager()->download($package, $this->getInstallPath($package), $prevPackage);
    }

    public function prepare($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->io->write('[asset] prepare ' . $type);

        return $this->composer->getDownloadManager()->prepare($type, $package, $this->getInstallPath($package), $prevPackage);
    }

    public function install(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        $dm = $this->composer->getDownloadManager();
        $this->io->write('[asset] install with ' . get_class($dm->getDownloaderForPackage($package)) . ' type=' . $dm->getDownloaderType($dm->getDownloaderForPackage($package)));
        $path = $this->getInstallPath($package);

        return $dm->install($package, $path)->then(function () use ($repo, $package, $path) {
            $repo->addPackage(clone $package);
            $this->io->write('[asset] installed into ' . $path . ': ' . implode(',', array_diff(scandir($path), array('.', '..'))));
        });
    }

    public function update(InstalledRepositoryInterface $repo, PackageInterface $initial, PackageInterface $target)
    {
        $this->io->write('[asset] update ' . $initial->getPrettyVersion() . ' => ' . $target->getPrettyVersion());

        return $this->composer->getDownloadManager()->update($initial, $target, $this->getInstallPath($target))->then(function () use ($repo, $initial, $target) {
            $repo->removePackage($initial);
            $repo->addPackage(clone $target);
        });
    }

    public function uninstall(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        $this->io->write('[asset] uninstall ' . $package->getPrettyName());

        return $this->composer->getDownloadManager()->remove($package, $this->getInstallPath($package))->then(function () use ($repo, $package) {
            $repo->removePackage($package);
        });
    }

    public function cleanup($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->io->write('[asset] cleanup ' . $type);

        return $this->composer->getDownloadManager()->cleanup($type, $package, $this->getInstallPath($package), $prevPackage);
    }

    public function getInstallPath(PackageInterface $package)
    {
        $parts = explode('/', $package->getPrettyName());

        return 'public/assets/' . $parts[1];
    }
}
