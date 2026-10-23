<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\ArchiveDownloader, standing in
// for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/ArchiveDownloader.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). Extraction is performed by gocomposer.

namespace Composer\Downloader;

use Composer\Package\PackageInterface;
use GocomposerPluginHost\Unsupported;

abstract class ArchiveDownloader extends FileDownloader
{
    /** @var array<string, true> */
    protected $cleanupExecuted = [];

    protected function extract(PackageInterface $package, string $file, string $path)
    {
        Unsupported::fail(__METHOD__);
    }
}
