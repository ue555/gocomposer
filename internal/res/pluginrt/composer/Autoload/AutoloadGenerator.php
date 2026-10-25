<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Autoload\AutoloadGenerator for plugins,
// standing in for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Autoload/AutoloadGenerator.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). dump() runs gocomposer's Go autoload generator
// (with the pre/post-autoload-dump events when runScripts is set, like
// Composer) into the given target directory, with any platform requirement
// filter (filters implemented by plugins are asked by gocomposer). The
// helpers building package maps and loaders (buildPackageMap,
// parseAutoloads, createLoader and what they use) are Composer's code,
// copied unchanged.

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
use Composer\Package\AliasPackage;
use Composer\Pcre\Preg;
use Composer\Util\Filesystem;
use Composer\Util\PackageSorter;
use Composer\Util\Platform;
use Composer\ClassMapGenerator\ClassMapGenerator;
use Composer\Package\Locker;
use Composer\Package\PackageInterface;
use Composer\Package\RootPackageInterface;
use Composer\Repository\InstalledRepositoryInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Services;

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
            'targetDir' => $targetDir,
            'platformFilter' => Services::filterSpec($this->platformRequirementFilter),
        ]);

        // dump() sets the dev mode from installed.json when it was not set
        if ($this->devMode === null && isset($r['devMode'])) {
            $this->devMode = (bool) $r['devMode'];
        }

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

    // Composer's code below: AutoloadGenerator.php of Composer
    // f83588392df43d4fd027d84d26de459946ef4e80 (MIT license, (c) Nils
    // Adermann, Jordi Boggiano), copied unchanged

    /**
     * Same as buildExclusionRegex, but also excludes the vendor dir if $dir contains it
     *
     * @param array<string> $excluded
     * @return non-empty-string|null
     */
    private function buildExclusionRegexWithVendorDir(string $dir, array $excluded, ?string $vendorPath): ?string
    {
        if (null !== $vendorPath) {
            $fs = new Filesystem();
            $resolvedDir = $fs->normalizePath($fs->isAbsolutePath($dir) ? $dir : realpath(Platform::getCwd()).'/'.$dir);
            if (str_starts_with($vendorPath, $resolvedDir.'/')) {
                $excluded[] = preg_quote($vendorPath.'/');
            }
        }

        return $this->buildExclusionRegex($dir, $excluded);
    }

    /**
     * @param array<string> $excluded
     * @return non-empty-string|null
     */
    private function buildExclusionRegex(string $dir, array $excluded): ?string
    {
        if ([] === $excluded) {
            return null;
        }

        // filter excluded patterns here to only use those matching $dir
        // exclude-from-classmap patterns are all realpath'd so we can only filter them if $dir exists so that realpath($dir) will work
        // if $dir does not exist, it should anyway not find anything there so no trouble
        if (file_exists($dir)) {
            // transform $dir in the same way that exclude-from-classmap patterns are transformed so we can match them against each other
            $dirMatch = preg_quote(strtr(realpath($dir), '\\', '/'));
            // also match against the non-realpath version for symlinks
            $fs = new Filesystem();
            $absDir = $fs->isAbsolutePath($dir) ? $dir : realpath(Platform::getCwd()).'/'.$dir;
            $dirMatchNormalized = preg_quote(strtr($fs->normalizePath($absDir), '\\', '/'));
            $isSymlink = $dirMatch !== $dirMatchNormalized;

            foreach ($excluded as $index => $pattern) {
                // extract the constant string prefix of the pattern here, until we reach a non-escaped regex special character
                $pattern = Preg::replace('{^(([^.+*?\[^\]$(){}=!<>|:\\\\#-]+|\\\\[.+*?\[^\]$(){}=!<>|:#-])*).*}', '$1', $pattern);
                // if the pattern is not a subset or superset of $dir, it is unrelated and we skip it
                if (
                    (!str_starts_with($pattern, $dirMatch) && !str_starts_with($dirMatch, $pattern))
                    && (
                        // only check dirMatchNormalized if it is actually different from dirMatch otherwise there is no point
                        !$isSymlink
                        || (!str_starts_with($pattern, $dirMatchNormalized) && !str_starts_with($dirMatchNormalized, $pattern))
                    )
                ) {
                    unset($excluded[$index]);
                }
            }
        }

        return \count($excluded) > 0 ? '{(' . implode('|', $excluded) . ')}' : null;
    }

    /**
     * @param PackageInterface[] $packages
     * @return non-empty-array<int, array{0: PackageInterface, 1: string|null}>
     */
    public function buildPackageMap(InstallationManager $installationManager, PackageInterface $rootPackage, array $packages)
    {
        // build package => install path map
        $packageMap = [[$rootPackage, '']];

        foreach ($packages as $package) {
            if ($package instanceof AliasPackage) {
                continue;
            }
            $this->validatePackage($package);
            $packageMap[] = [
                $package,
                $installationManager->getInstallPath($package),
            ];
        }

        return $packageMap;
    }

    /**
     * @return void
     * @throws \InvalidArgumentException Throws an exception, if the package has illegal settings.
     */
    protected function validatePackage(PackageInterface $package)
    {
        $autoload = $package->getAutoload();
        if (!empty($autoload['psr-4']) && null !== $package->getTargetDir()) {
            $name = $package->getName();
            $package->getTargetDir();
            throw new \InvalidArgumentException("PSR-4 autoloading is incompatible with the target-dir property, remove the target-dir in package '$name'.");
        }
        if (!empty($autoload['psr-4'])) {
            foreach ($autoload['psr-4'] as $namespace => $dirs) {
                if ($namespace !== '' && '\\' !== substr($namespace, -1)) {
                    throw new \InvalidArgumentException("psr-4 namespaces must end with a namespace separator, '$namespace' does not, use '$namespace\\'.");
                }
            }
        }
    }

    /**
     * Compiles an ordered list of namespace => path mappings
     *
     * @param non-empty-array<int, array{0: PackageInterface, 1: string|null}> $packageMap array of array(package, installDir-relative-to-composer.json or null for metapackages)
     * @param RootPackageInterface $rootPackage root package instance
     * @param bool|string[] $filteredDevPackages If an array, the list of packages that must be removed. If bool, whether to filter out require-dev packages
     * @return array
     * @phpstan-return array{
     *     'psr-0': array<string, array<string>>,
     *     'psr-4': array<string, array<string>>,
     *     'classmap': array<int, string>,
     *     'files': array<string, string>,
     *     'exclude-from-classmap': array<int, string>,
     * }
     */
    public function parseAutoloads(array $packageMap, PackageInterface $rootPackage, $filteredDevPackages = false)
    {
        if (isset($packageMap[0][0]) && $packageMap[0][0] instanceof RootPackageInterface) {
            $rootPackageMap = array_shift($packageMap);
        } else {
            $rootPackageMap = null;
        }
        if (is_array($filteredDevPackages)) {
            $packageMap = array_filter($packageMap, static function ($item) use ($filteredDevPackages): bool {
                return !in_array($item[0]->getName(), $filteredDevPackages, true);
            });
        } elseif ($filteredDevPackages) {
            $packageMap = $this->filterPackageMap($packageMap, $rootPackage);
        }
        $sortedPackageMap = $this->sortPackageMap($packageMap);
        if ($rootPackageMap !== null) {
            $sortedPackageMap[] = $rootPackageMap;
        }
        $reverseSortedMap = array_reverse($sortedPackageMap);

        // reverse-sorted means root first, then dependents, then their dependents, etc.
        // which makes sense to allow root to override classmap or psr-0/4 entries with higher precedence rules
        $psr0 = $this->parseAutoloadsType($reverseSortedMap, 'psr-0', $rootPackage);
        $psr4 = $this->parseAutoloadsType($reverseSortedMap, 'psr-4', $rootPackage);
        $classmap = $this->parseAutoloadsType($reverseSortedMap, 'classmap', $rootPackage);

        // sorted (i.e. dependents first) for files to ensure that dependencies are loaded/available once a file is included
        $files = $this->parseAutoloadsType($sortedPackageMap, 'files', $rootPackage);
        // using sorted here but it does not really matter as all are excluded equally
        $exclude = $this->parseAutoloadsType($sortedPackageMap, 'exclude-from-classmap', $rootPackage);

        krsort($psr0);
        krsort($psr4);

        return [
            'psr-0' => $psr0,
            'psr-4' => $psr4,
            'classmap' => $classmap,
            'files' => $files,
            'exclude-from-classmap' => $exclude,
        ];
    }

    /**
     * Registers an autoloader based on an autoload-map returned by parseAutoloads
     *
     * @param array<string, mixed[]> $autoloads see parseAutoloads return value
     * @return ClassLoader
     */
    public function createLoader(array $autoloads, ?string $vendorDir = null)
    {
        $loader = new ClassLoader($vendorDir);

        if (isset($autoloads['psr-0'])) {
            foreach ($autoloads['psr-0'] as $namespace => $path) {
                $loader->add($namespace, $path);
            }
        }

        if (isset($autoloads['psr-4'])) {
            foreach ($autoloads['psr-4'] as $namespace => $path) {
                $loader->addPsr4($namespace, $path);
            }
        }

        if (isset($autoloads['classmap'])) {
            $excluded = [];
            if (!empty($autoloads['exclude-from-classmap'])) {
                $excluded = $autoloads['exclude-from-classmap'];
            }

            $vendorPath = null;
            if (null !== $vendorDir && false !== ($realVendorDir = realpath($vendorDir))) {
                $vendorPath = (new Filesystem())->normalizePath($realVendorDir);
            }

            $classMapGenerator = new ClassMapGenerator(['php', 'inc', 'hh']);
            $classMapGenerator->avoidDuplicateScans();

            foreach ($autoloads['classmap'] as $dir) {
                try {
                    $classMapGenerator->scanPaths($dir, $this->buildExclusionRegexWithVendorDir($dir, $excluded, $vendorPath));
                } catch (\RuntimeException $e) {
                    $this->io->writeError('<warning>'.$e->getMessage().'</warning>');
                }
            }

            $loader->addClassMap($classMapGenerator->getClassMap()->getMap());
        }

        return $loader;
    }

    /**
     * @param array<int, array{0: PackageInterface, 1: string|null}> $packageMap
     * @param string $type one of: 'psr-0'|'psr-4'|'classmap'|'files'
     * @return array<int, string>|array<string, array<string>>|array<string, string>
     */
    protected function parseAutoloadsType(array $packageMap, string $type, RootPackageInterface $rootPackage)
    {
        $autoloads = [];

        foreach ($packageMap as $item) {
            [$package, $installPath] = $item;

            // packages that are not installed cannot autoload anything
            if (null === $installPath) {
                continue;
            }

            $autoload = $package->getAutoload();
            if ($this->devMode && $package === $rootPackage) {
                $autoload = array_merge_recursive($autoload, $package->getDevAutoload());
            }

            // skip misconfigured packages
            if (!isset($autoload[$type]) || !is_array($autoload[$type])) {
                continue;
            }
            if (null !== $package->getTargetDir() && $package !== $rootPackage) {
                $installPath = substr($installPath, 0, -strlen('/'.$package->getTargetDir()));
            }

            foreach ($autoload[$type] as $namespace => $paths) {
                if (in_array($type, ['psr-4', 'psr-0'], true)) {
                    // normalize namespaces to ensure "\" becomes "" and others do not have leading separators as they are not needed
                    $namespace = ltrim($namespace, '\\');
                }
                foreach ((array) $paths as $path) {
                    if (($type === 'files' || $type === 'classmap' || $type === 'exclude-from-classmap') && $package->getTargetDir() && !Filesystem::isReadable($installPath.'/'.$path)) {
                        // remove target-dir from file paths of the root package
                        if ($package === $rootPackage) {
                            $targetDir = str_replace('\\<dirsep\\>', '[\\\\/]', preg_quote(str_replace(['/', '\\'], '<dirsep>', $package->getTargetDir())));
                            $path = ltrim(Preg::replace('{^'.$targetDir.'}', '', ltrim($path, '\\/')), '\\/');
                        } else {
                            // add target-dir from file paths that don't have it
                            $path = $package->getTargetDir() . '/' . $path;
                        }
                    }

                    if ($type === 'exclude-from-classmap') {
                        // first escape user input
                        $path = Preg::replace('{/+}', '/', preg_quote(trim(strtr($path, '\\', '/'), '/')));

                        // add support for wildcards * and **
                        $path = strtr($path, ['\\*\\*' => '.+?', '\\*' => '[^/]+?']);

                        // add support for up-level relative paths
                        $updir = null;
                        $path = Preg::replaceCallback(
                            '{^((?:(?:\\\\\\.){1,2}+/)+)}',
                            static function ($matches) use (&$updir): string {
                                // undo preg_quote for the matched string
                                $updir = str_replace('\\.', '.', $matches[1]);

                                return '';
                            },
                            $path
                        );
                        if (empty($installPath)) {
                            $installPath = strtr(Platform::getCwd(), '\\', '/');
                        }

                        $resolvedPath = realpath($installPath . '/' . $updir);
                        if (false === $resolvedPath) {
                            continue;
                        }
                        $autoloads[] = preg_quote(strtr($resolvedPath, '\\', '/')) . '/' . $path . '($|/)';
                        continue;
                    }

                    $relativePath = empty($installPath) ? (empty($path) ? '.' : $path) : $installPath.'/'.$path;

                    if ($type === 'files') {
                        $autoloads[$this->getFileIdentifier($package, $path)] = $relativePath;
                        continue;
                    }
                    if ($type === 'classmap') {
                        $autoloads[] = $relativePath;
                        continue;
                    }

                    $autoloads[$namespace][] = $relativePath;
                }
            }
        }

        return $autoloads;
    }

    /**
     * @return string
     */
    protected function getFileIdentifier(PackageInterface $package, string $path)
    {
        // TODO composer v3 change this to sha1 or xxh3? Possibly not worth the potential breakage though
        return hash('md5', $package->getName() . ':' . $path);
    }

    /**
     * Filters out dev-dependencies
     *
     * @param array<int, array{0: PackageInterface, 1: string|null}> $packageMap
     * @return array<int, array{0: PackageInterface, 1: string|null}>
     */
    protected function filterPackageMap(array $packageMap, RootPackageInterface $rootPackage)
    {
        $packages = [];
        $include = [];
        $replacedBy = [];

        foreach ($packageMap as $item) {
            $package = $item[0];
            $name = $package->getName();
            $packages[$name] = $package;
            foreach ($package->getReplaces() as $replace) {
                $replacedBy[$replace->getTarget()] = $name;
            }
        }

        $add = static function (PackageInterface $package) use (&$add, $packages, &$include, $replacedBy): void {
            foreach ($package->getRequires() as $link) {
                $target = $link->getTarget();
                if (isset($replacedBy[$target])) {
                    $target = $replacedBy[$target];
                }
                if (!isset($include[$target])) {
                    $include[$target] = true;
                    if (isset($packages[$target])) {
                        $add($packages[$target]);
                    }
                }
            }
        };
        $add($rootPackage);

        return array_filter(
            $packageMap,
            static function ($item) use ($include): bool {
                $package = $item[0];
                foreach ($package->getNames() as $name) {
                    if (isset($include[$name])) {
                        return true;
                    }
                }

                return false;
            }
        );
    }

    /**
     * Sorts packages by dependency weight
     *
     * Packages of equal weight are sorted alphabetically
     *
     * @param array<int, array{0: PackageInterface, 1: string|null}> $packageMap
     * @return array<int, array{0: PackageInterface, 1: string|null}>
     */
    protected function sortPackageMap(array $packageMap)
    {
        $packages = [];
        $paths = [];

        foreach ($packageMap as $item) {
            [$package, $path] = $item;
            $name = $package->getName();
            $packages[$name] = $package;
            $paths[$name] = $path;
        }

        $sortedPackages = PackageSorter::sortPackages($packages);

        $sortedPackageMap = [];

        foreach ($sortedPackages as $package) {
            $name = $package->getName();
            $sortedPackageMap[] = [$packages[$name], $paths[$name]];
        }

        return $sortedPackageMap;
    }
}
