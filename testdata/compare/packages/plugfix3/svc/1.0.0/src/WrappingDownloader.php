<?php
namespace Plugfix3\Svc;

use Composer\Downloader\DownloaderInterface;
use Composer\Package\PackageInterface;
use React\Promise\PromiseInterface;

/**
 * A downloader not built on Composer's classes: logs and delegates to the
 * downloader it replaced.
 */
class WrappingDownloader implements DownloaderInterface
{
    private $inner;
    private $io;

    public function __construct($io, DownloaderInterface $inner)
    {
        $this->io = $io;
        $this->inner = $inner;
    }

    public function getInstallationSource(): string
    {
        return $this->inner->getInstallationSource();
    }

    public function download(PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        $this->io->write('[wrap] download ' . $package->getPrettyName());

        return $this->inner->download($package, $path, $prevPackage);
    }

    public function prepare(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        $this->io->write('[wrap] prepare ' . $type . ' ' . $package->getPrettyName());

        return $this->inner->prepare($type, $package, $path, $prevPackage);
    }

    public function install(PackageInterface $package, string $path): PromiseInterface
    {
        $this->io->write('[wrap] install ' . $package->getPrettyName());

        return $this->inner->install($package, $path);
    }

    public function update(PackageInterface $initial, PackageInterface $target, string $path): PromiseInterface
    {
        $this->io->write('[wrap] update ' . $target->getPrettyName());

        return $this->inner->update($initial, $target, $path);
    }

    public function remove(PackageInterface $package, string $path): PromiseInterface
    {
        $this->io->write('[wrap] remove ' . $package->getPrettyName());

        return $this->inner->remove($package, $path);
    }

    public function cleanup(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        $this->io->write('[wrap] cleanup ' . $type . ' ' . $package->getPrettyName());

        return $this->inner->cleanup($type, $package, $path, $prevPackage);
    }
}
