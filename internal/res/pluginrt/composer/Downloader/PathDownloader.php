<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\PathDownloader, standing in
// for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/PathDownloader.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). Performed by gocomposer's path downloader.

namespace Composer\Downloader;

class PathDownloader extends FileDownloader implements VcsCapableDownloaderInterface
{
    public function getVcsReference(\Composer\Package\PackageInterface $package, string $path): ?string
    {
        return \GocomposerPluginHost\Host::cb('dlVcsReference', ['p' => \GocomposerPluginHost\Host::ref($package), 'path' => $path]);
    }
}
