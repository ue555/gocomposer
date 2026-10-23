<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Installer\InstallationManager for
// plugins, based on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Installer/InstallationManager.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). Installers added by plugins live here and are
// consulted by gocomposer for install paths and operations; the default
// installers delegate to gocomposer's Go implementation.

namespace Composer\Installer;

use Composer\IO\IOInterface;
use Composer\Package\AliasPackage;
use Composer\Package\PackageInterface;
use Composer\PartialComposer;
use Composer\Repository\InstalledRepositoryInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class InstallationManager
{
    /** @var list<InstallerInterface> */
    private $installers = [];
    /** @var array<string, InstallerInterface> */
    private $cache = [];
    /** @var list<InstallerInterface> */
    private $defaults = [];
    /** @var int gocomposer instance id */
    public $gcInst = 0;

    public function __construct($loop = null, ?IOInterface $io = null, $eventDispatcher = null)
    {
    }

    /**
     * @internal
     */
    public function setDefaultInstallers(IOInterface $io, PartialComposer $composer): void
    {
        $this->defaults = [
            new LibraryInstaller($io, $composer, null),
            new PluginInstaller($io, $composer),
            new MetapackageInstaller($io),
        ];
    }

    public function reset(): void
    {
    }

    public function addInstaller(InstallerInterface $installer): void
    {
        array_unshift($this->installers, $installer);
        $this->cache = [];
        if ($this->gcInst === 0) {
            Host::installersChanged();
        }
    }

    public function removeInstaller(InstallerInterface $installer): void
    {
        if (false !== ($key = array_search($installer, $this->installers, true))) {
            array_splice($this->installers, $key, 1);
            $this->cache = [];
            if ($this->gcInst === 0) {
                Host::installersChanged();
            }
        }
    }

    public function disablePlugins(): void
    {
    }

    /**
     * @internal the installer added by a plugin which handles the type, if any
     */
    public function getPluginInstaller(string $type): ?InstallerInterface
    {
        $type = strtolower($type);
        foreach ($this->installers as $installer) {
            if ($installer->supports($type)) {
                return $installer;
            }
        }

        return null;
    }

    public function getInstaller(string $type): InstallerInterface
    {
        $type = strtolower($type);

        if (isset($this->cache[$type])) {
            return $this->cache[$type];
        }

        foreach (array_merge($this->installers, array_reverse($this->defaults)) as $installer) {
            if ($installer->supports($type)) {
                return $this->cache[$type] = $installer;
            }
        }

        throw new \InvalidArgumentException('Unknown installer type: '.$type);
    }

    public function isPackageInstalled(InstalledRepositoryInterface $repo, PackageInterface $package): bool
    {
        if ($package instanceof AliasPackage) {
            return $repo->hasPackage($package) && $this->isPackageInstalled($repo, $package->getAliasOf());
        }

        return $this->getInstaller($package->getType())->isInstalled($repo, $package);
    }

    public function ensureBinariesPresence(PackageInterface $package): void
    {
        try {
            $installer = $this->getInstaller($package->getType());
        } catch (\InvalidArgumentException $e) {
            // no installer found for the current package type (@see `getInstaller()`)
            return;
        }

        // if the given installer support installing binaries
        if ($installer instanceof BinaryPresenceInterface) {
            $installer->ensureBinariesPresence($package);
        }
    }

    public function getInstallPath(PackageInterface $package): ?string
    {
        $installer = $this->getInstaller($package->getType());

        return $installer->getInstallPath($package);
    }

    public function setOutputProgress(bool $outputProgress): void
    {
        Host::cb('imOutputProgress', ['v' => $outputProgress, 'inst' => $this->gcInst]);
    }

    public function execute(InstalledRepositoryInterface $repo, array $operations, bool $devMode = true, bool $runScripts = true, bool $downloadOnly = false): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function download(PackageInterface $package)
    {
        Unsupported::fail(__METHOD__);
    }

    public function install(InstalledRepositoryInterface $repo, $operation)
    {
        Unsupported::fail(__METHOD__);
    }

    public function update(InstalledRepositoryInterface $repo, $operation)
    {
        Unsupported::fail(__METHOD__);
    }

    public function uninstall(InstalledRepositoryInterface $repo, $operation)
    {
        Unsupported::fail(__METHOD__);
    }

    public function notifyInstalls(IOInterface $io): void
    {
    }
}
