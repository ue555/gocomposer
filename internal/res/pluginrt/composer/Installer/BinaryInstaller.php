<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Installer\BinaryInstaller for plugins,
// based on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Installer/BinaryInstaller.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The proxies are written by gocomposer's Go port.

namespace Composer\Installer;

use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Util\Filesystem;
use GocomposerPluginHost\Host;

class BinaryInstaller
{
    /** @var string */
    protected $binDir;
    /** @var string */
    protected $binCompat;
    /** @var IOInterface */
    protected $io;
    /** @var Filesystem */
    protected $filesystem;
    /** @var string|null */
    private $vendorDir;

    public function __construct(IOInterface $io, string $binDir, string $binCompat, ?Filesystem $filesystem = null, ?string $vendorDir = null)
    {
        $this->binDir = $binDir;
        $this->binCompat = $binCompat;
        $this->io = $io;
        $this->filesystem = $filesystem ?: new Filesystem();
        $this->vendorDir = $vendorDir;
    }

    public function installBinaries(PackageInterface $package, string $installPath, bool $warnOnOverwrite = true): void
    {
        Host::cb('libOp', ['m' => 'installBinaries', 'p' => Host::ref($package), 'path' => $installPath, 'warn' => $warnOnOverwrite]);
    }

    public function removeBinaries(PackageInterface $package): void
    {
        Host::cb('libOp', ['m' => 'removeBinaries', 'p' => Host::ref($package)]);
    }

    public static function determineBinaryCaller(string $bin): string
    {
        return (string) Host::cb('binaryCaller', ['bin' => $bin]);
    }

    /**
     * @return string[]
     */
    protected function getBinaries(PackageInterface $package): array
    {
        return $package->getBinaries();
    }
}
