<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Installer for plugins, standing in for
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Installer.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The installer that runs
// is gocomposer's Go port: an object of this class either stands for the
// gocomposer Installer currently running (it is on the call stack of plugin
// code like Composer's, see Host::framed()) or collects the options of an
// installer created by plugin code; run() runs gocomposer's Installer with
// those options for the Composer instance given to the constructor.

namespace Composer;

use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterFactory;
use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterInterface;
use Composer\Installer\SuggestedPackagesReporter;
use Composer\IO\IOInterface;
use Composer\IO\NullIO;
use Composer\Repository\RepositoryInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class Installer
{
    public const ERROR_NONE = 0;
    public const ERROR_GENERIC_FAILURE = 1;
    public const ERROR_NO_LOCK_FILE_FOR_PARTIAL_UPDATE = 3;
    public const ERROR_LOCK_FILE_INVALID = 4;
    public const ERROR_DEPENDENCY_RESOLUTION_FAILED = 2;
    public const ERROR_AUDIT_FAILED = 5;
    public const ERROR_TRANSPORT_EXCEPTION = 100;

    /** @var IOInterface|null */
    protected $io;
    /** @var Config|null */
    protected $config;
    /** @var PlatformRequirementFilterInterface */
    protected $platformRequirementFilter;
    /** @var SuggestedPackagesReporter|null */
    protected $suggestedPackagesReporter;

    /** @var int the gocomposer Installer this object stands for (0: none) */
    private $gcId = 0;
    /** @var int gocomposer instance id run() installs */
    private $gcInst = 0;
    /** @var array<string, mixed> options set on this object */
    private $gcOpts = [];
    /** @var bool */
    private $gcDryRun = false;
    /** @var bool */
    private $gcVerbose = false;
    /** @var array{0: PlatformRequirementFilterInterface, 1: bool|string[]}|null the filter gocomposer's Installer uses */
    private $gcFilter = null;

    public function __construct(IOInterface $io, Config $config, ...$services)
    {
        $this->io = $io;
        $this->config = $config;
        $this->gcInst = $config->gcInst;
        if ($this->platformRequirementFilter === null) {
            $this->platformRequirementFilter = PlatformRequirementFilterFactory::ignoreNothing();
        }
    }

    /**
     * @internal the object standing for gocomposer's Installer $id
     * @param array<string, mixed> $state
     */
    public static function gcFor(int $id, array $state): self
    {
        $installer = (new \ReflectionClass(self::class))->newInstanceWithoutConstructor();
        $installer->gcId = $id;
        $installer->io = Host::$io;
        $installer->config = Host::$composer->getConfig();
        $installer->gcDryRun = (bool) ($state['dryRun'] ?? false);
        $installer->gcVerbose = (bool) ($state['verbose'] ?? false);
        $filter = $state['platformFilter'] ?? false;
        $filter = is_array($filter) ? $filter : (bool) $filter;
        $installer->platformRequirementFilter = PlatformRequirementFilterFactory::fromBoolOrList($filter);
        $installer->gcFilter = [$installer->platformRequirementFilter, $filter];

        return $installer;
    }

    /**
     * @internal runs $fn with this object on the call stack (Installer::run)
     * @return mixed
     */
    public function gocomposerFrame(callable $fn)
    {
        return $fn();
    }

    public static function create(IOInterface $io, Composer $composer): self
    {
        return new static(
            $io,
            $composer->getConfig(),
            $composer->getPackage(),
            $composer->getDownloadManager(),
            $composer->getRepositoryManager(),
            $composer->getLocker(),
            $composer->getInstallationManager(),
            $composer->getEventDispatcher(),
            $composer->getAutoloadGenerator()
        );
    }

    public function run(): int
    {
        $opts = $this->gcOpts;
        $filter = $this->platformRequirementFilter;
        if ($this->gcFilter !== null && $filter === $this->gcFilter[0]) {
            $opts['platformFilter'] = $this->gcFilter[1];
        } elseif ($filter instanceof \Composer\Filter\PlatformRequirementFilter\IgnoreAllPlatformRequirementFilter) {
            $opts['platformFilter'] = true;
        } elseif ($filter instanceof \Composer\Filter\PlatformRequirementFilter\IgnoreNothingPlatformRequirementFilter) {
            $opts['platformFilter'] = false;
        } else {
            // a filter implemented by plugin code: gocomposer asks it
            $opts['platformFilter'] = \GocomposerPluginHost\Services::filterSpec($filter);
        }
        if ($this->gcAuditConfig !== null) {
            // getAuditConfig(): the AuditConfig set, else one made of audit/auditFormat
            $opts['audit'] = $this->gcAuditConfig->audit;
            $opts['auditFormat'] = $this->gcAuditConfig->auditFormat;
        }

        return (int) Host::cb('installerRun', ['base' => $this->gcId, 'inst' => $this->gcInst, 'obj' => spl_object_id($this), 'opts' => (object) $opts]);
    }

    /** @var \Composer\Advisory\AuditConfig|null */
    private $gcAuditConfig = null;

    /**
     * @param mixed $value
     * @return $this
     */
    private function gcSet(string $key, $value): self
    {
        $this->gcOpts[$key] = $value;

        return $this;
    }

    public function setIgnoredTypes(array $types): self
    {
        return $this->gcSet('ignoredTypes', array_values($types));
    }

    public function setAllowedTypes(?array $types): self
    {
        return $this->gcSet('allowedTypes', $types === null ? null : array_values($types));
    }

    public function setAdditionalFixedRepository(RepositoryInterface $additionalFixedRepository): self
    {
        // gocomposer resolves with the repository (plugin_repos.go)
        return $this->gcSet('additionalFixed', \GocomposerPluginHost\Services::repoSpec($additionalFixedRepository));
    }

    public function setTemporaryConstraints(array $constraints): self
    {
        $out = [];
        foreach ($constraints as $name => $constraint) {
            $out[$name] = \GocomposerPluginHost\Api::constraintString($constraint);
        }

        return $this->gcSet('temporaryConstraints', (object) $out);
    }

    public function setDryRun(bool $dryRun = true): self
    {
        $this->gcDryRun = $dryRun;

        return $this->gcSet('dryRun', $dryRun);
    }

    public function isDryRun(): bool
    {
        return $this->gcDryRun;
    }

    public function setDownloadOnly(bool $downloadOnly = true): self
    {
        return $this->gcSet('downloadOnly', $downloadOnly);
    }

    public function setPreferSource(bool $preferSource = true): self
    {
        return $this->gcSet('preferSource', $preferSource);
    }

    public function setPreferDist(bool $preferDist = true): self
    {
        return $this->gcSet('preferDist', $preferDist);
    }

    public function setOptimizeAutoloader(bool $optimizeAutoloader): self
    {
        return $this->gcSet('optimizeAutoloader', $optimizeAutoloader);
    }

    public function setClassMapAuthoritative(bool $classMapAuthoritative): self
    {
        if ($classMapAuthoritative) {
            $this->setOptimizeAutoloader(true);
        }

        return $this->gcSet('classMapAuthoritative', $classMapAuthoritative);
    }

    public function setApcuAutoloader(bool $apcuAutoloader, ?string $apcuAutoloaderPrefix = null): self
    {
        $this->gcSet('apcuPrefix', $apcuAutoloaderPrefix);

        return $this->gcSet('apcuAutoloader', $apcuAutoloader);
    }

    public function setStrictPsrAutoloader(bool $strictPsr): self
    {
        return $this->gcSet('strictPsr', $strictPsr);
    }

    public function setUpdate(bool $update): self
    {
        return $this->gcSet('update', $update);
    }

    public function setInstall(bool $install): self
    {
        return $this->gcSet('install', $install);
    }

    public function setDevMode(bool $devMode = true): self
    {
        return $this->gcSet('devMode', $devMode);
    }

    public function setDumpAutoloader(bool $dumpAutoloader = true): self
    {
        return $this->gcSet('dumpAutoloader', $dumpAutoloader);
    }

    public function setRunScripts(bool $runScripts = true): self
    {
        return $this->gcSet('runScripts', $runScripts);
    }

    public function setConfig(Config $config): self
    {
        $this->config = $config;

        return $this;
    }

    public function setVerbose(bool $verbose = true): self
    {
        $this->gcVerbose = $verbose;

        return $this->gcSet('verbose', $verbose);
    }

    public function isVerbose(): bool
    {
        return $this->gcVerbose;
    }

    /**
     * @param bool|string[] $ignorePlatformReqs
     */
    public function setIgnorePlatformRequirements($ignorePlatformReqs): self
    {
        trigger_error('Installer::setIgnorePlatformRequirements is deprecated since Composer 2.2, use setPlatformRequirementFilter instead.', E_USER_DEPRECATED);

        return $this->setPlatformRequirementFilter(PlatformRequirementFilterFactory::fromBoolOrList($ignorePlatformReqs));
    }

    public function setPlatformRequirementFilter(PlatformRequirementFilterInterface $platformRequirementFilter): self
    {
        $this->platformRequirementFilter = $platformRequirementFilter;

        return $this;
    }

    public function setUpdateMirrors(bool $updateMirrors): self
    {
        return $this->gcSet('updateMirrors', $updateMirrors);
    }

    public function setUpdateAllowList(array $packages): self
    {
        return $this->gcSet('updateAllowList', array_values(array_unique(array_map('strtolower', $packages))));
    }

    public function setUpdateAllowTransitiveDependencies(int $updateAllowTransitiveDependencies): self
    {
        return $this->gcSet('updateAllowTransitive', $updateAllowTransitiveDependencies);
    }

    public function setPreferStable(bool $preferStable = true): self
    {
        return $this->gcSet('preferStable', $preferStable);
    }

    public function setPreferLowest(bool $preferLowest = true): self
    {
        return $this->gcSet('preferLowest', $preferLowest);
    }

    public function setMinimalUpdate(bool $minimalUpdate = true): self
    {
        return $this->gcSet('minimalUpdate', $minimalUpdate);
    }

    public function setWriteLock(bool $writeLock = true): self
    {
        return $this->gcSet('writeLock', $writeLock);
    }

    public function setExecuteOperations(bool $executeOperations = true): self
    {
        return $this->gcSet('executeOperations', $executeOperations);
    }

    public function setAudit(bool $audit): self
    {
        $this->gcAuditConfig = null; // Invalidate cached config

        return $this->gcSet('audit', $audit);
    }

    public function setErrorOnAudit(bool $errorOnAudit): self
    {
        return $this->gcSet('errorOnAudit', $errorOnAudit);
    }

    public function setAuditFormat(string $auditFormat): self
    {
        $this->gcAuditConfig = null; // Invalidate cached config

        return $this->gcSet('auditFormat', $auditFormat);
    }

    public function setAuditConfig(\Composer\Advisory\AuditConfig $auditConfig): self
    {
        $this->gcAuditConfig = $auditConfig;

        return $this;
    }

    public function setPolicyConfig($policyConfig): self
    {
        Unsupported::fail(__METHOD__);
    }

    public function disablePlugins(): self
    {
        return $this;
    }

    public function setSuggestedPackagesReporter(SuggestedPackagesReporter $suggestedPackagesReporter): self
    {
        $this->suggestedPackagesReporter = $suggestedPackagesReporter;
        $io = \Closure::bind(function () {
            return $this->io;
        }, $suggestedPackagesReporter, SuggestedPackagesReporter::class)();
        $silent = $io instanceof NullIO;
        if ($this->gcId !== 0) {
            // the gocomposer Installer is running: its reporter is replaced
            Host::cb('installerSet', ['id' => $this->gcId, 'silentSuggestions' => $silent]);
        }

        return $this->gcSet('silentSuggestions', $silent);
    }

    /**
     * The lock transaction of the last update run() made (null when it did
     * not update).
     */
    public function getLockTransaction(): ?\Composer\DependencyResolver\LockTransaction
    {
        $data = Host::cb('lockTx', ['base' => $this->gcId, 'obj' => spl_object_id($this)]);

        return $data === null ? null : \Composer\DependencyResolver\LockTransaction::gcCreate($data);
    }
}
