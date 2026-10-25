<?php
namespace Plugfix3\Svc;

use Composer\Downloader\ZipDownloader;
use Composer\Package\PackageInterface;

/**
 * Exposes the protected helpers of the zip downloader.
 */
class InternalsZip extends ZipDownloader
{
    public function fileName(PackageInterface $package, string $path)
    {
        return $this->getFileName($package, $path);
    }

    public function url(PackageInterface $package, string $url)
    {
        return $this->processUrl($package, $url);
    }
}
