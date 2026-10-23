<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\FileDownloader for plugins,
// standing in for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/FileDownloader.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The operations (pre/post-file-download events, cache,
// checksums, extraction, removal) are performed by gocomposer's Go port of
// the downloader the class stands for; subclasses written by plugins can
// override them and call parent:: like with Composer.

namespace Composer\Downloader;

use Composer\Cache;
use Composer\Config;
use Composer\EventDispatcher\EventDispatcher;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Util\Filesystem;
use Composer\Util\HttpDownloader;
use Composer\Util\ProcessExecutor;
use GocomposerPluginHost\GoDownloader;
use GocomposerPluginHost\Unsupported;

class FileDownloader implements DownloaderInterface, ChangeReportInterface
{
    use GoDownloader;

    /** @var array<string, int|string> */
    public static $downloadMetadata = [];
    /** @var array<string, string[]> */
    public static $responseHeaders = [];

    /** @var IOInterface */
    protected $io;
    /** @var Config */
    protected $config;
    /** @var HttpDownloader|null */
    protected $httpDownloader;
    /** @var Filesystem */
    protected $filesystem;
    /** @var ?Cache */
    protected $cache;
    /** @var ?EventDispatcher */
    protected $eventDispatcher;
    /** @var ProcessExecutor */
    protected $process;

    public function __construct(IOInterface $io, Config $config, ?HttpDownloader $httpDownloader = null, ?EventDispatcher $eventDispatcher = null, ?Cache $cache = null, ?Filesystem $filesystem = null, ?ProcessExecutor $process = null)
    {
        $this->io = $io;
        $this->config = $config;
        $this->eventDispatcher = $eventDispatcher;
        $this->httpDownloader = $httpDownloader;
        $this->cache = $cache;
        $this->process = $process ?? new ProcessExecutor($io);
        $this->filesystem = $filesystem ?? new Filesystem($this->process);
    }

    public function getInstallationSource(): string
    {
        return 'dist';
    }

    protected function getFileName(PackageInterface $package, string $path): string
    {
        Unsupported::fail(__METHOD__);
    }

    protected function processUrl(PackageInterface $package, string $url): string
    {
        Unsupported::fail(__METHOD__);
    }

    protected function clearLastCacheWrite(PackageInterface $package): void
    {
    }

    protected function addCleanupPath(PackageInterface $package, string $path): void
    {
    }

    protected function removeCleanupPath(PackageInterface $package, string $path): void
    {
    }

    public function getLocalChanges(PackageInterface $package, string $path): ?string
    {
        Unsupported::fail(__METHOD__);
    }
}
