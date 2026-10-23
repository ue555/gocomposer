<?php
namespace Plugfix\Installer;

use Composer\Installer\LibraryInstaller;
use Composer\Package\PackageInterface;
use Composer\Repository\InstalledRepositoryInterface;
use React\Promise\PromiseInterface;

class ModuleInstaller extends LibraryInstaller
{
    public function supports($packageType)
    {
        return $packageType === 'plugfix-module';
    }

    public function getInstallPath(PackageInterface $package)
    {
        $extra = $this->composer->getPackage()->getExtra();
        $base = isset($extra['plugfix-module-dir']) ? $extra['plugfix-module-dir'] : 'modules';
        $parts = explode('/', $package->getPrettyName());
        $name = isset($package->getExtra()['module-name']) ? $package->getExtra()['module-name'] : $parts[1];
        return $base . '/' . $name;
    }

    public function install(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        $io = $this->io;
        $path = $this->getInstallPath($package);
        $promise = parent::install($repo, $package);
        $done = function () use ($io, $path, $package) {
            $io->write(sprintf('[plugfix/installer] installed %s into %s (%s)', $package->getPrettyName(), $path, is_dir($path) ? 'present' : 'missing'));
        };
        if ($promise instanceof PromiseInterface) {
            return $promise->then($done);
        }
        $done();
        return null;
    }

    public function uninstall(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        $installPath = $this->getPackageBasePath($package);
        $io = $this->io;
        $outputStatus = function () use ($io, $installPath) {
            $io->write(sprintf('Deleting %s - %s', $installPath, !file_exists($installPath) ? '<comment>deleted</comment>' : '<error>not deleted</error>'));
        };
        $promise = parent::uninstall($repo, $package);
        if ($promise instanceof PromiseInterface) {
            return $promise->then($outputStatus);
        }
        $outputStatus();
        return null;
    }
}
