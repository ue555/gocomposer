<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\VcsDownloader, standing in
// for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/VcsDownloader.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The VCS operations are performed by gocomposer (git; the
// other VCS fail like gocomposer's own downloaders).

namespace Composer\Downloader;

use Composer\Config;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Util\Filesystem;
use Composer\Util\ProcessExecutor;
use GocomposerPluginHost\GoDownloader;
use GocomposerPluginHost\Unsupported;

abstract class VcsDownloader implements DownloaderInterface, ChangeReportInterface, VcsCapableDownloaderInterface
{
    use GoDownloader;

    /** @var IOInterface */
    protected $io;
    /** @var Config */
    protected $config;
    /** @var ProcessExecutor */
    protected $process;
    /** @var Filesystem */
    protected $filesystem;
    /** @var array<string, true> */
    protected $hasCleanedChanges = [];

    public function __construct(IOInterface $io, Config $config, ?ProcessExecutor $process = null, ?Filesystem $fs = null)
    {
        $this->io = $io;
        $this->config = $config;
        $this->process = $process ?? new ProcessExecutor($io);
        $this->filesystem = $fs ?? new Filesystem($this->process);
    }

    public function getInstallationSource(): string
    {
        return 'source';
    }

    public function getLocalChanges(PackageInterface $package, string $path): ?string
    {
        Unsupported::fail(__METHOD__);
    }

    public function getVcsReference(PackageInterface $package, string $path): ?string
    {
        Unsupported::fail(__METHOD__);
    }
}
