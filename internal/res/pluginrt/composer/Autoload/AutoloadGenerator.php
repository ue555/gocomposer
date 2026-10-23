<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Autoload\AutoloadGenerator for plugins,
// standing in for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Autoload/AutoloadGenerator.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). dump() runs gocomposer's Go autoload generator
// (with the pre/post-autoload-dump events when runScripts is set, like
// Composer); the lower level helpers building package maps and loaders are
// not available.

namespace Composer\Autoload;

use Composer\ClassMapGenerator\ClassMap;
use Composer\Config;
use Composer\EventDispatcher\EventDispatcher;
use Composer\Filter\PlatformRequirementFilter\IgnoreAllPlatformRequirementFilter;
use Composer\Filter\PlatformRequirementFilter\IgnoreNothingPlatformRequirementFilter;
use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterFactory;
use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterInterface;
use Composer\Installer\InstallationManager;
use Composer\IO\IOInterface;
use Composer\Package\Locker;
use Composer\Package\PackageInterface;
use Composer\Package\RootPackageInterface;
use Composer\Repository\InstalledRepositoryInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class AutoloadGenerator
{
    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var EventDispatcher */
    private $eventDispatcher;
    /** @var IOInterface */
    private $io;
    /** @var ?bool */
    private $devMode = null;
    /** @var bool */
    private $classMapAuthoritative = false;
    /** @var bool */
    private $apcu = false;
    /** @var string|null */
    private $apcuPrefix;
    /** @var bool */
    private $dryRun = false;
    /** @var bool */
    private $runScripts = false;
    /** @var PlatformRequirementFilterInterface */
    private $platformRequirementFilter;

    public function __construct(EventDispatcher $eventDispatcher, ?IOInterface $io = null)
    {
        $this->eventDispatcher = $eventDispatcher;
        $this->io = $io ?? Host::$io;
        $this->platformRequirementFilter = PlatformRequirementFilterFactory::ignoreNothing();
    }

    public function setDevMode(bool $devMode = true)
    {
        $this->devMode = $devMode;
    }

    public function setClassMapAuthoritative(bool $classMapAuthoritative)
    {
        $this->classMapAuthoritative = $classMapAuthoritative;
    }

    public function setApcu(bool $apcu, ?string $apcuPrefix = null)
    {
        $this->apcu = $apcu;
        $this->apcuPrefix = $apcuPrefix;
    }

    public function setRunScripts(bool $runScripts = true)
    {
        $this->runScripts = $runScripts;
    }

    public function setDryRun(bool $dryRun = true)
    {
        $this->dryRun = $dryRun;
    }

    public function setIgnorePlatformRequirements($ignorePlatformReqs)
    {
        trigger_error('AutoloadGenerator::setIgnorePlatformRequirements is deprecated since Composer 2.2, use setPlatformRequirementFilter instead.', E_USER_DEPRECATED);

        $this->setPlatformRequirementFilter(PlatformRequirementFilterFactory::fromBoolOrList($ignorePlatformReqs));
    }

    public function setPlatformRequirementFilter(PlatformRequirementFilterInterface $platformRequirementFilter)
    {
        $this->platformRequirementFilter = $platformRequirementFilter;
    }

    /**
     * @return ClassMap
     */
    public function dump(Config $config, InstalledRepositoryInterface $localRepo, RootPackageInterface $rootPackage, InstallationManager $installationManager, string $targetDir, bool $scanPsrPackages = false, ?string $suffix = null, ?Locker $locker = null, bool $strictAmbiguous = false)
    {
        if ($targetDir !== 'composer') {
            Unsupported::fail(__METHOD__.' with a target directory other than "composer"');
        }
        if ($this->platformRequirementFilter instanceof IgnoreAllPlatformRequirementFilter) {
            $filter = true;
        } elseif ($this->platformRequirementFilter instanceof IgnoreNothingPlatformRequirementFilter) {
            $filter = false;
        } else {
            Unsupported::fail(__METHOD__.' with a '.get_class($this->platformRequirementFilter).' platform requirement filter');
        }
        $r = Host::cb('autoloadDump', [
            'inst' => $this->gcInst,
            'dev' => $this->devMode,
            'optimize' => $scanPsrPackages,
            'authoritative' => $this->classMapAuthoritative,
            'apcu' => $this->apcu,
            'apcuPrefix' => $this->apcuPrefix,
            'runScripts' => $this->runScripts,
            'dryRun' => $this->dryRun,
            'suffix' => $suffix,
            'strictAmbiguous' => $strictAmbiguous,
            'platformFilter' => $filter,
        ]);

        $classMap = new ClassMap();
        foreach ($r['map'] ?? [] as $class => $path) {
            $classMap->addClass($class, $path);
        }
        foreach ($r['ambiguous'] ?? [] as $class => $paths) {
            foreach ($paths as $path) {
                $classMap->addAmbiguousClass($class, $path);
            }
        }
        foreach ($r['psrViolations'] ?? [] as $path => $warning) {
            $classMap->addPsrViolation($warning, '', $path);
        }

        return $classMap;
    }

    public function buildPackageMap(InstallationManager $installationManager, PackageInterface $rootPackage, array $packages)
    {
        Unsupported::fail(__METHOD__);
    }

    public function parseAutoloads(array $packageMap, PackageInterface $rootPackage, $filteredDevPackages = false)
    {
        Unsupported::fail(__METHOD__);
    }

    public function createLoader(array $autoloads, ?string $vendorDir = null)
    {
        Unsupported::fail(__METHOD__);
    }
}
