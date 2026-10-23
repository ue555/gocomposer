<?php
namespace Plugfix\CustomInstallers;

use Composer\Installer\LibraryInstaller;
use Composer\Package\PackageInterface;

/**
 * A LibraryInstaller replacing the download and code hooks (calling the
 * parent implementations).
 */
class HookInstaller extends LibraryInstaller
{
    private function log($msg)
    {
        $this->io->write('[hook] ' . $msg);
    }

    public function download(PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->log('download ' . $package->getPrettyName() . ' ' . $package->getPrettyVersion() . ($prevPackage ? ' from ' . $prevPackage->getPrettyVersion() : ''));
        if (getenv('PLUGFIX_HOOK_FAIL') === 'download') {
            throw new \RuntimeException('plugfix hook installer download failure');
        }

        return parent::download($package, $prevPackage);
    }

    public function prepare($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->log('prepare ' . $type . ' ' . $package->getPrettyName());

        return parent::prepare($type, $package, $prevPackage);
    }

    public function cleanup($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->log('cleanup ' . $type . ' ' . $package->getPrettyName());

        return parent::cleanup($type, $package, $prevPackage);
    }

    protected function installCode(PackageInterface $package)
    {
        $this->log('installCode ' . $package->getPrettyName() . ' into ' . $this->getInstallPath($package));
        if (getenv('PLUGFIX_HOOK_FAIL') === 'install') {
            throw new \RuntimeException('plugfix hook installer install failure');
        }
        $promise = parent::installCode($package);
        $path = $this->getInstallPath($package);
        $io = $this->io;

        return $promise->then(function () use ($io, $path) {
            $io->write('[hook] installed files: ' . implode(',', array_diff(scandir($path), array('.', '..'))));
        });
    }

    protected function updateCode(PackageInterface $initial, PackageInterface $target)
    {
        $this->log('updateCode ' . $initial->getPrettyVersion() . ' => ' . $target->getPrettyVersion());

        return parent::updateCode($initial, $target);
    }

    protected function removeCode(PackageInterface $package)
    {
        $this->log('removeCode ' . $package->getPrettyName());

        return parent::removeCode($package);
    }
}
