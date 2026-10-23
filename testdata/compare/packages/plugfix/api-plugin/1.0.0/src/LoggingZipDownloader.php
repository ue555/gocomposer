<?php
namespace Plugfix\Api;

use Composer\Downloader\ZipDownloader;
use Composer\Package\PackageInterface;
use React\Promise\PromiseInterface;

/**
 * Downloader registered for the custom dist type "plugfixzip": logs, then
 * downloads and extracts like the zip downloader.
 */
class LoggingZipDownloader extends ZipDownloader
{
    public function download(PackageInterface $package, string $path, ?PackageInterface $prevPackage = null, bool $output = true): PromiseInterface
    {
        $this->io->write('[zipdl] download ' . $package->getPrettyName() . ' ' . $package->getDistType() . ' ' . $package->getDistUrl());

        return parent::download($package, $path, $prevPackage, $output);
    }

    public function install(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        $this->io->write('[zipdl] install ' . $package->getPrettyName());

        return parent::install($package, $path, $output);
    }
}
