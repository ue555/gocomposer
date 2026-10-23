<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Installer\LibraryInstaller for
// plugins, based on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Installer/LibraryInstaller.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The methods are ported unchanged on top of the
// gocomposer-backed DownloadManager and BinaryInstaller. For subclasses that
// keep Composer's download/code hooks, install()/update()/uninstall() hand
// the whole operation to gocomposer's Go LibraryInstaller (same behaviour,
// done in one call).

namespace Composer\Installer;

use Composer\Composer;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\PartialComposer;
use Composer\Pcre\Preg;
use Composer\Util\Silencer;
use React\Promise\PromiseInterface;
use Composer\Repository\InstalledRepositoryInterface;
use Composer\Util\Filesystem;
use Composer\Util\Platform;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class LibraryInstaller implements InstallerInterface, BinaryPresenceInterface
{
    /** @var PartialComposer */
    protected $composer;
    /** @var string */
    protected $vendorDir;
    /** @var mixed */
    protected $downloadManager;
    /** @var IOInterface */
    protected $io;
    /** @var string */
    protected $type;
    /** @var Filesystem */
    protected $filesystem;
    /** @var BinaryInstaller */
    protected $binaryInstaller;

    public function __construct(IOInterface $io, PartialComposer $composer, ?string $type = 'library', ?Filesystem $filesystem = null, ?BinaryInstaller $binaryInstaller = null)
    {
        $this->composer = $composer;
        $this->downloadManager = $composer instanceof Composer ? $composer->getDownloadManager() : null;
        $this->io = $io;
        $this->type = $type;

        $this->filesystem = $filesystem ?: new Filesystem();
        $this->vendorDir = rtrim($composer->getConfig()->get('vendor-dir'), '/');
        $this->binaryInstaller = $binaryInstaller ?: new BinaryInstaller($this->io, rtrim($composer->getConfig()->get('bin-dir'), '/'), $composer->getConfig()->get('bin-compat'), $this->filesystem, $this->vendorDir);
    }

    public function supports(string $packageType)
    {
        return $packageType === $this->type || null === $this->type;
    }

    public function isInstalled(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        if (!$repo->hasPackage($package)) {
            return false;
        }

        $installPath = $this->getInstallPath($package);

        if (Filesystem::isReadable($installPath)) {
            return true;
        }

        if (Platform::isWindows() && $this->filesystem->isJunction($installPath)) {
            return true;
        }

        if (is_link($installPath)) {
            if (realpath($installPath) === false) {
                return false;
            }

            return true;
        }

        return false;
    }

    /** @var bool the subclass replaces one of the download/code hooks */
    private $gcPortedFlow = false;

    /**
     * @internal set by the gocomposer plugin host
     */
    public function gcSetPortedFlow(bool $ported): void
    {
        $this->gcPortedFlow = $ported;
    }

    public function download(PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->initializeVendorDir();
        $downloadPath = $this->getInstallPath($package);

        return $this->getDownloadManager()->download($package, $downloadPath, $prevPackage);
    }

    public function prepare($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->initializeVendorDir();
        $downloadPath = $this->getInstallPath($package);

        return $this->getDownloadManager()->prepare($type, $package, $downloadPath, $prevPackage);
    }

    public function cleanup($type, PackageInterface $package, ?PackageInterface $prevPackage = null)
    {
        $this->initializeVendorDir();
        $downloadPath = $this->getInstallPath($package);

        return $this->getDownloadManager()->cleanup($type, $package, $downloadPath, $prevPackage);
    }

    public function install(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        if (!$this->gcPortedFlow) {
            Host::cb('libOp', ['m' => 'install', 'p' => Host::ref($package)]);

            return \React\Promise\resolve(null);
        }

        $this->initializeVendorDir();
        $downloadPath = $this->getInstallPath($package);

        // remove the binaries if it appears the package files are missing
        if (!Filesystem::isReadable($downloadPath) && $repo->hasPackage($package)) {
            $this->binaryInstaller->removeBinaries($package);
        }

        $promise = $this->installCode($package);
        if (!$promise instanceof PromiseInterface) {
            $promise = \React\Promise\resolve(null);
        }

        $binaryInstaller = $this->binaryInstaller;
        $installPath = $this->getInstallPath($package);

        return $promise->then(static function () use ($binaryInstaller, $installPath, $package, $repo): void {
            $binaryInstaller->installBinaries($package, $installPath);
            if (!$repo->hasPackage($package)) {
                $repo->addPackage(clone $package);
            }
        });
    }

    public function update(InstalledRepositoryInterface $repo, PackageInterface $initial, PackageInterface $target)
    {
        if (!$repo->hasPackage($initial)) {
            throw new \InvalidArgumentException('Package is not installed: '.$initial);
        }
        if (!$this->gcPortedFlow) {
            Host::cb('libOp', ['m' => 'update', 'i' => Host::ref($initial), 'p' => Host::ref($target)]);

            return \React\Promise\resolve(null);
        }

        $this->initializeVendorDir();

        $this->binaryInstaller->removeBinaries($initial);
        $promise = $this->updateCode($initial, $target);
        if (!$promise instanceof PromiseInterface) {
            $promise = \React\Promise\resolve(null);
        }

        $binaryInstaller = $this->binaryInstaller;
        $installPath = $this->getInstallPath($target);

        return $promise->then(static function () use ($binaryInstaller, $installPath, $target, $initial, $repo): void {
            $binaryInstaller->installBinaries($target, $installPath);
            $repo->removePackage($initial);
            if (!$repo->hasPackage($target)) {
                $repo->addPackage(clone $target);
            }
        });
    }

    public function uninstall(InstalledRepositoryInterface $repo, PackageInterface $package)
    {
        if (!$repo->hasPackage($package)) {
            throw new \InvalidArgumentException('Package is not installed: '.$package);
        }
        if (!$this->gcPortedFlow) {
            Host::cb('libOp', ['m' => 'uninstall', 'p' => Host::ref($package)]);

            return \React\Promise\resolve(null);
        }

        $promise = $this->removeCode($package);
        if (!$promise instanceof PromiseInterface) {
            $promise = \React\Promise\resolve(null);
        }

        $binaryInstaller = $this->binaryInstaller;
        $downloadPath = $this->getPackageBasePath($package);
        $filesystem = $this->filesystem;

        return $promise->then(static function () use ($binaryInstaller, $filesystem, $downloadPath, $package, $repo): void {
            $binaryInstaller->removeBinaries($package);
            $repo->removePackage($package);

            if (strpos($package->getName(), '/')) {
                $packageVendorDir = dirname($downloadPath);
                if (is_dir($packageVendorDir) && $filesystem->isDirEmpty($packageVendorDir)) {
                    Silencer::call('rmdir', $packageVendorDir);
                }
            }
        });
    }

    public function getInstallPath(PackageInterface $package)
    {
        $this->initializeVendorDir();

        $basePath = ($this->vendorDir ? $this->vendorDir.'/' : '') . $package->getPrettyName();
        $targetDir = $package->getTargetDir();

        return $basePath . ($targetDir ? '/'.$targetDir : '');
    }

    public function ensureBinariesPresence(PackageInterface $package)
    {
        $this->binaryInstaller->installBinaries($package, $this->getInstallPath($package), false);
    }

    /**
     * @return string
     */
    protected function getPackageBasePath(PackageInterface $package)
    {
        $installPath = $this->getInstallPath($package);
        $targetDir = $package->getTargetDir();

        if ($targetDir) {
            return Preg::replace('{/*'.str_replace('/', '/+', preg_quote($targetDir)).'/?$}', '', $installPath);
        }

        return $installPath;
    }

    protected function installCode(PackageInterface $package)
    {
        $downloadPath = $this->getInstallPath($package);

        return $this->getDownloadManager()->install($package, $downloadPath);
    }

    protected function updateCode(PackageInterface $initial, PackageInterface $target)
    {
        $initialDownloadPath = $this->getInstallPath($initial);
        $targetDownloadPath = $this->getInstallPath($target);
        if ($targetDownloadPath !== $initialDownloadPath) {
            // if the target and initial dirs intersect, we force a remove + install
            // to avoid the rename wiping the target dir as part of the initial dir cleanup
            if (strpos($initialDownloadPath, $targetDownloadPath) === 0
                || strpos($targetDownloadPath, $initialDownloadPath) === 0
            ) {
                $promise = $this->removeCode($initial);
                if (!$promise instanceof PromiseInterface) {
                    $promise = \React\Promise\resolve(null);
                }

                return $promise->then(function () use ($target): PromiseInterface {
                    $promise = $this->installCode($target);
                    if ($promise instanceof PromiseInterface) {
                        return $promise;
                    }

                    return \React\Promise\resolve(null);
                });
            }

            $this->filesystem->rename($initialDownloadPath, $targetDownloadPath);
        }

        return $this->getDownloadManager()->update($initial, $target, $targetDownloadPath);
    }

    protected function removeCode(PackageInterface $package)
    {
        $downloadPath = $this->getPackageBasePath($package);

        return $this->getDownloadManager()->remove($package, $downloadPath);
    }

    protected function initializeVendorDir()
    {
        $this->filesystem->ensureDirectoryExists($this->vendorDir);
        $this->vendorDir = realpath($this->vendorDir);
    }

    protected function getDownloadManager()
    {
        assert($this->downloadManager instanceof \Composer\Downloader\DownloadManager, new \LogicException(self::class.' should be initialized with a fully loaded Composer instance to be able to install/... packages'));

        return $this->downloadManager;
    }
}
