<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Repository\InstalledFilesystemRepository
// for plugins (Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/InstalledFilesystemRepository.php and
// FilesystemRepository.php, MIT license, (c) Nils Adermann, Jordi Boggiano).
// The packages mirror gocomposer's local repository (vendor/composer/
// installed.json state of the running operation) and are refreshed from Go;
// adding/removing packages and writing the repository is done by gocomposer.

namespace Composer\Repository;

use Composer\Installer\InstallationManager;
use Composer\Package\PackageInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class InstalledFilesystemRepository extends WritableArrayRepository implements InstalledRepositoryInterface
{
    /** @var string|null */
    private $syncedIds = null;
    /** @var ?bool */
    private $devMode = null;
    /** @var bool */
    private $fresh = false;
    /** @var int gocomposer instance id */
    public $gcInst = 0;

    public function __construct($repositoryFile = null, bool $dumpVersions = false, $rootPackage = null, $filesystem = null)
    {
        parent::__construct([]);
    }

    public function getRepoName()
    {
        return 'installed '.parent::getRepoName();
    }

    public function isFresh()
    {
        $this->sync();

        return $this->fresh;
    }

    /**
     * @return ?bool true if dev requirements were installed, false if --no-dev was used, null if yet unknown
     */
    public function getDevMode()
    {
        $this->sync();

        return $this->devMode;
    }

    public function getDevPackageNames()
    {
        $this->sync();

        return parent::getDevPackageNames();
    }

    public function getPackages()
    {
        $this->sync();

        return parent::getPackages();
    }

    public function count(): int
    {
        $this->sync();

        return parent::count();
    }

    public function hasPackage(PackageInterface $package)
    {
        $this->sync();

        return parent::hasPackage($package);
    }

    protected function initialize()
    {
        parent::initialize();
        $this->syncedIds = null;
        $this->sync();
    }

    /**
     * @internal refreshes the packages from gocomposer when they changed
     */
    public function sync(): void
    {
        if (null === $this->packages) {
            $this->packages = [];
            $this->syncedIds = null;
        }
        $state = Host::cb('localRepo', ['known' => $this->syncedIds, 'inst' => $this->gcInst]);
        if (!is_array($state) || !empty($state['same'])) {
            return;
        }
        $this->syncedIds = $state['key'];
        $this->devMode = $state['dev'];
        $this->fresh = (bool) $state['fresh'];
        $this->packages = [];
        $this->packageMap = null;
        $this->devPackageNames = $state['devNames'];
        foreach ($state['pkgs'] as $ref) {
            $package = Host::pkg($ref);
            $this->packages[] = $package;
            if (!$package->getRepository()) {
                $package->setRepository($this);
            }
        }
    }

    public function reload()
    {
        $this->syncedIds = null;
        $this->sync();
    }

    /**
     * Writes vendor/composer/installed.json and installed.php (done by
     * gocomposer, with the dev package names set here).
     */
    public function write(bool $devMode, InstallationManager $installationManager)
    {
        Host::cb('localRepoOp', ['m' => 'write', 'dev' => $devMode, 'devNames' => $this->devPackageNames, 'inst' => $this->gcInst]);
        $this->syncedIds = null;
    }

    public function addPackage(PackageInterface $package)
    {
        Host::cb('localRepoOp', ['m' => 'add', 'p' => Host::ref($package), 'inst' => $this->gcInst]);
        $this->syncedIds = null;
        $this->sync();
    }

    public function removePackage(PackageInterface $package)
    {
        Host::cb('localRepoOp', ['m' => 'remove', 'p' => Host::ref($package), 'inst' => $this->gcInst]);
        $this->syncedIds = null;
        $this->sync();
    }

    public function setDevPackageNames(array $devPackageNames)
    {
        $this->sync();
        parent::setDevPackageNames($devPackageNames);
    }
}
