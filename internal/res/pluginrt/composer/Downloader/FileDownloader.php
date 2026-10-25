<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\FileDownloader for plugins,
// standing in for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/FileDownloader.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The operations (pre/post-file-download events, cache,
// checksums, extraction, removal) are performed by gocomposer's Go port of
// the downloader the class stands for; subclasses written by plugins can
// override them and call parent:: like with Composer. update() (removing,
// then installing through $this), processUrl() and the appendix of the
// operation line are Composer's code; getFileName() and getLocalChanges()
// answer with gocomposer's file names and change detection.

namespace Composer\Downloader;

use Composer\Cache;
use Composer\Config;
use Composer\EventDispatcher\EventDispatcher;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Util\Filesystem;
use Composer\Util\HttpDownloader;
use Composer\Util\ProcessExecutor;
use Composer\Util\Url as UrlUtil;
use Composer\DependencyResolver\Operation\UpdateOperation;
use GocomposerPluginHost\Api;
use GocomposerPluginHost\GoDownloader;
use GocomposerPluginHost\Host;
use React\Promise\PromiseInterface;

class FileDownloader implements DownloaderInterface, ChangeReportInterface
{
    use GoDownloader {
        update as gcGoUpdate;
    }

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

    /**
     * The temporary file the package is downloaded to (gocomposer's, named
     * like Composer's vendor/composer/tmp-<hash>.<extension>).
     */
    protected function getFileName(PackageInterface $package, string $path): string
    {
        return (string) Host::cb('dlFileName', ['p' => Host::ref($package), 'path' => $path]);
    }

    protected function processUrl(PackageInterface $package, string $url): string
    {
        if (!extension_loaded('openssl') && 0 === strpos($url, 'https:')) {
            throw new \RuntimeException('You must enable the openssl extension to download files via https');
        }

        if ($package->getDistReference() !== null) {
            $url = UrlUtil::updateDistReference($this->config, $url, $package->getDistReference());
        }

        return $url;
    }

    /**
     * Gets appendix message to add to the "- Upgrading x" string being output on update
     */
    protected function getInstallOperationAppendix(PackageInterface $package, string $path): string
    {
        return '';
    }

    public function update(PackageInterface $initial, PackageInterface $target, string $path): PromiseInterface
    {
        if (!in_array(Api::downloaderKind($this), ['zip', 'rar', 'tar', 'gzip', 'xz', 'phar', 'file'], true)) {
            return $this->gcGoUpdate($initial, $target, $path);
        }
        $this->io->writeError("  - " . UpdateOperation::format($initial, $target) . $this->getInstallOperationAppendix($target, $path));

        $promise = $this->remove($initial, $path, false);

        return $promise->then(function () use ($target, $path): PromiseInterface {
            return $this->install($target, $path, false);
        });
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
        return Host::cb('dlLocalChanges', ['kind' => Api::downloaderKind($this), 'p' => Host::ref($package), 'path' => $path]);
    }
}
