<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Downloader\ArchiveDownloader, standing in
// for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/ArchiveDownloader.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). Extraction is performed by gocomposer. When a
// plugin's subclass overrides extract(), install() follows Composer's
// ArchiveDownloader::install (empty the target, extract() into a temporary
// directory, move the files into place) with gocomposer doing the file
// system steps and the parent extract().

namespace Composer\Downloader;

use Composer\DependencyResolver\Operation\InstallOperation;
use Composer\Package\PackageInterface;
use GocomposerPluginHost\Api;
use GocomposerPluginHost\Host;
use React\Promise\PromiseInterface;

abstract class ArchiveDownloader extends FileDownloader
{
    /** @var array<string, true> */
    protected $cleanupExecuted = [];

    public function install(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        $declaring = (new \ReflectionMethod($this, 'extract'))->getDeclaringClass()->getName();
        if (strpos($declaring, 'Composer\\Downloader\\') === 0) {
            return parent::install($package, $path, $output);
        }

        if ($output) {
            $this->io->writeError("  - " . InstallOperation::format($package) . $this->getInstallOperationAppendix($package, $path));
        }

        $ref = Host::ref($package);
        $st = Host::cb('dlArchive', ['m' => 'begin', 'p' => $ref, 'path' => $path]);
        $temporaryDir = $st['tmp'];
        $fileName = $this->getFileName($package, $path);
        $cleanup = static function () use ($ref, $path, $temporaryDir): void {
            Host::cb('dlArchive', ['m' => 'fail', 'p' => $ref, 'path' => $path, 'tmp' => $temporaryDir]);
        };

        try {
            $promise = $this->extract($package, $fileName, $temporaryDir);
        } catch (\Exception $e) {
            $cleanup();
            throw $e;
        }
        if (!$promise instanceof PromiseInterface) {
            $promise = \React\Promise\resolve(null);
        }

        return $promise->then(static function () use ($ref, $path, $temporaryDir): void {
            Host::cb('dlArchive', ['m' => 'finish', 'p' => $ref, 'path' => $path, 'tmp' => $temporaryDir]);
        }, static function ($e) use ($cleanup) {
            $cleanup();

            throw $e;
        });
    }

    /**
     * @inheritDoc
     */
    protected function getInstallOperationAppendix(PackageInterface $package, string $path): string
    {
        return ': Extracting archive';
    }

    /**
     * Extracts the archive with gocomposer's extractor of the archive type
     * the class stands for.
     *
     * @return PromiseInterface<void|null>
     */
    protected function extract(PackageInterface $package, string $file, string $path): PromiseInterface
    {
        $r = Host::cb('dlExtract', ['kind' => Api::downloaderKind($this), 'p' => Host::ref($package), 'file' => $file, 'path' => $path]);
        if (is_array($r) && !empty($r['async'])) {
            // ZipDownloader extracts with unzip in a background process
            return Api::defer();
        }

        return \React\Promise\resolve(null);
    }
}
