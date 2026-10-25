<?php
// gocomposer plugin host: Composer services plugins replace or extend (see
// internal/composer/plugin_services.go). gocomposer keeps doing the work;
// where a plugin put its own object in place of a Composer service, the
// object is called at the point Composer calls it. The flows follow Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT license, (c) Nils Adermann,
// Jordi Boggiano).

namespace GocomposerPluginHost {

use Composer\Filter\PlatformRequirementFilter\IgnoreAllPlatformRequirementFilter;
use Composer\Filter\PlatformRequirementFilter\IgnoreNothingPlatformRequirementFilter;
use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterFactory;
use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterInterface;
use Composer\Package\AliasPackage;
use Composer\Package\Dumper\ArrayDumper;
use Composer\Package\PackageInterface;
use Composer\Repository\RepositoryInterface;

final class Services
{
    /**
     * An operation of a downloader replaced by a plugin, chained like
     * InstallationManager::executeBatch does for LibraryInstaller (the
     * prepare promise is settled already): the DownloadManager operation,
     * LibraryInstaller's binaries/repository step (gocomposer), the
     * cleanup of the download manager, the repository write (gocomposer);
     * settled with the batch (Host::settleOps).
     *
     * @param array<string, mixed> $a
     */
    public static function replacedOp(array $a): void
    {
        $dm = Host::$composer->gcExistingDownloadManager() ?? Host::$composer->getDownloadManager();
        $p = Host::pkg($a['p']);
        $i = isset($a['i']) ? Host::pkg($a['i']) : null;
        if (array_key_exists('psrc', $a) && $p->getInstallationSource() !== $a['psrc']) {
            $p->setInstallationSource($a['psrc']);
        }
        if ($i !== null && array_key_exists('isrc', $a) && $i->getInstallationSource() !== $a['isrc']) {
            $i->setInstallationSource($a['isrc']);
        }
        $m = $a['m'];
        $path = $a['path'];
        $key = $a['key'];
        $io = Host::$io;

        $promise = \React\Promise\resolve(null)->then(static function () use ($dm, $m, $p, $i, $path) {
            switch ($m) {
                case 'install':
                    return $dm->install($p, $path);
                case 'update':
                    return $dm->update($i, $p, $path);
                default:
                    return $dm->remove($p, $path);
            }
        })->then(static function () use ($key): void {
            Host::cb('replacedPost', ['key' => $key]);
        })->then(static function () use ($dm, $m, $p, $i, $path) {
            // avoid calling cleanup if the download was not even initialized for a package
            if (null === $p->getInstallationSource()) {
                return \React\Promise\resolve(null);
            }

            return $dm->cleanup($m, $p, $path, $i);
        })->then(static function () use ($key): void {
            Host::cb('opDone', ['key' => $key]);
        }, static function ($e) use ($m, $p, $io): void {
            $io->writeError('    <error>' . ucfirst($m) .' of '.$p->getPrettyName().' failed</error>');

            throw $e;
        });
        Host::$pendingOps[] = $promise;
    }

    /**
     * Factory::createArchiveManager for a Composer instance.
     */
    public static function archiveManager(\Composer\Composer $composer): \Composer\Package\Archiver\ArchiveManager
    {
        $am = new \Composer\Package\Archiver\ArchiveManager($composer->getDownloadManager(), $composer->getLoop());
        if (class_exists(\ZipArchive::class)) {
            $am->addArchiver(new \Composer\Package\Archiver\ZipArchiver());
        }
        if (class_exists(\Phar::class)) {
            $am->addArchiver(new \Composer\Package\Archiver\PharArchiver());
        }
        self::$builtinArchiveManagers[spl_object_id($am)] = count(self::archivers($am));

        return $am;
    }

    /** @var array<int, int> archive managers created for gocomposer and their number of archivers */
    private static $builtinArchiveManagers = [];

    /**
     * @return list<\Composer\Package\Archiver\ArchiverInterface>
     */
    private static function archivers(\Composer\Package\Archiver\ArchiveManager $am): array
    {
        return \Closure::bind(function (): array {
            return $this->archivers;
        }, $am, \Composer\Package\Archiver\ArchiveManager::class)();
    }

    /**
     * Whether the archive manager of the instance running the plugins does
     * something gocomposer's archive command does not: replaced by a plugin
     * (setArchiveManager), or archivers added (addArchiver).
     */
    public static function archiveManagerChanged(): bool
    {
        $am = Host::$composer->gcExistingArchiveManager();
        if ($am === null) {
            return false;
        }
        if (get_class($am) !== \Composer\Package\Archiver\ArchiveManager::class || !isset(self::$builtinArchiveManagers[spl_object_id($am)])) {
            return true;
        }

        return count(self::archivers($am)) !== self::$builtinArchiveManagers[spl_object_id($am)];
    }

    /**
     * ArchiveCommand::archive's call of the archive manager of the instance
     * running the plugins.
     *
     * @param array<string, mixed> $a
     */
    public static function archiveVia(array $a): string
    {
        $package = ($a['p']['id'] ?? null) === 0 ? Host::$composer->getPackage() : Host::pkg($a['p']);

        return Host::$composer->getArchiveManager()->archive($package, $a['format'], $a['dest'], $a['fileName'] ?? null, (bool) $a['ignoreFilters']);
    }

    /** @var array<int, RepositoryInterface> repositories of plugin code gocomposer asks, by object id */
    private static $repos = [];
    /** @var int */
    private static $adoptSeq = 0;
    /** @var array<int, true> */
    private static $adoptSent = [];

    public static function registerRepo(RepositoryInterface $repo): int
    {
        $id = spl_object_id($repo);
        self::$repos[$id] = $repo;

        return $id;
    }

    public static function repoObject(int $id): RepositoryInterface
    {
        if (!isset(self::$repos[$id])) {
            throw new \LogicException('Unknown plugin repository '.$id);
        }

        return self::$repos[$id];
    }

    /**
     * Whether gocomposer implements the repository object itself (one of its
     * repositories, not subclassed nor changed by plugin code).
     */
    public static function isGoRepository(RepositoryInterface $repo): bool
    {
        if (Api::repoId($repo) === null) {
            return false;
        }
        if ($repo instanceof GoRepository) {
            return !$repo->gcIsMaterialized() && in_array(get_class($repo), [
                \Composer\Repository\ComposerRepository::class,
                \Composer\Repository\PathRepository::class,
                \Composer\Repository\VcsRepository::class,
                \Composer\Repository\ArtifactRepository::class,
                \Composer\Repository\PackageRepository::class,
            ], true);
        }

        return true;
    }

    /**
     * Describes a repository object for gocomposer.
     *
     * @return array<string, mixed>
     */
    public static function repoSpec(RepositoryInterface $repo): array
    {
        if (self::isGoRepository($repo)) {
            return ['rid' => Api::repoId($repo)];
        }

        return ['php' => self::registerRepo($repo), 'name' => $repo->getRepoName()];
    }

    /**
     * RepositoryManager::createRepository for a class registered with
     * setRepositoryClass().
     *
     * @param array<string, mixed> $config
     */
    public static function createRepository(int $inst, string $class, array $config, ?string $name): RepositoryInterface
    {
        $composer = Api::composer($inst);
        if (isset($config['packagist']) && false === $config['packagist']) {
            Host::$io->writeError('<warning>Repository "'.$name.'" ('.json_encode($config).') has a packagist key which should be in its own repository definition</warning>');
        }

        if (isset($config['only']) || isset($config['exclude']) || isset($config['canonical'])) {
            $filterConfig = $config;
            unset($config['only'], $config['exclude'], $config['canonical']);
        }

        $repository = new $class($config, Host::$io, $composer->getConfig(), Api::httpDownloader($inst), $composer->getEventDispatcher(), $composer->getLoop()->getProcessExecutor());

        if (isset($filterConfig)) {
            $repository = new \Composer\Repository\FilterRepository($repository, $filterConfig);
        }

        return $repository;
    }

    /**
     * The reference of a package of plugin code that gocomposer adopts
     * (negative ids: both sides keep the object/package for the id).
     *
     * @return array<string, mixed>
     */
    public static function adoptRef(PackageInterface $p): array
    {
        $oid = spl_object_id($p);
        if (isset(Host::$objIds[$oid]) && (Host::$pkgs[Host::$objIds[$oid]] ?? null) === $p) {
            $id = Host::$objIds[$oid];
            if ($id >= 0 || isset(self::$adoptSent[$id])) {
                return ['id' => $id];
            }
        } else {
            $id = --self::$adoptSeq;
            Host::$pkgs[$id] = $p;
            Host::$objIds[$oid] = $id;
        }
        self::$adoptSent[$id] = true;
        if ($p instanceof AliasPackage) {
            return ['id' => $id, 'd' => [
                '_alias_of' => self::adoptRef($p->getAliasOf()),
                'version' => $p->getPrettyVersion(),
                'version_normalized' => $p->getVersion(),
                '_root_alias' => $p->isRootPackageAlias(),
            ]];
        }

        return ['id' => $id, 'd' => (new ArrayDumper())->dump($p)];
    }

    /**
     * @param iterable<PackageInterface> $packages
     * @return list<array<string, mixed>>
     */
    private static function adoptRefs(iterable $packages): array
    {
        $out = [];
        foreach ($packages as $p) {
            $out[] = self::adoptRef($p);
        }

        return $out;
    }

    /**
     * A query of gocomposer to a repository of plugin code.
     *
     * @param array<string, mixed> $a
     * @return mixed
     */
    public static function phpRepoOp(array $a)
    {
        $repo = self::repoObject((int) $a['id']);
        switch ($a['m']) {
            case 'packages':
                return self::adoptRefs($repo->getPackages());
            case 'find':
                return self::adoptRefs($repo->findPackages($a['name'], $a['c'] ?? null));
            case 'load':
                $map = [];
                foreach ($a['names'] as $n) {
                    $map[$n['name']] = $n['c'] === null ? null : Api::constraint(['pretty' => $n['c']]);
                }
                $r = $repo->loadPackages($map, $a['acceptable'] ?? [], $a['flags'] ?? [], $a['loaded'] ?? []);

                return ['namesFound' => array_values($r['namesFound'] ?? []), 'packages' => self::adoptRefs($r['packages'] ?? [])];
            case 'providers':
                return (object) $repo->getProviders($a['name']);
        }

        throw new \LogicException('Unknown repository operation '.$a['m']);
    }

    /** @var array<int, PlatformRequirementFilterInterface> filters of plugin code gocomposer asks */
    private static $filters = [];
    /** @var list<PlatformRequirementFilterInterface> filters created for gocomposer's own */
    private static $goFilters = [];

    /**
     * Describes a platform requirement filter for gocomposer.
     *
     * @return array<string, mixed>
     */
    public static function filterSpec(PlatformRequirementFilterInterface $filter): array
    {
        if (in_array($filter, self::$goFilters, true)) {
            return ['kind' => 'go'];
        }
        if ($filter instanceof IgnoreAllPlatformRequirementFilter) {
            return ['kind' => 'all'];
        }
        if ($filter instanceof IgnoreNothingPlatformRequirementFilter) {
            return ['kind' => 'nothing'];
        }
        $id = spl_object_id($filter);
        self::$filters[$id] = $filter;

        return ['kind' => 'php', 'id' => $id];
    }

    /**
     * The Composer filter object for a filter of gocomposer.
     *
     * @param array<string, mixed> $spec
     */
    public static function filterFromGo(array $spec): PlatformRequirementFilterInterface
    {
        switch ($spec['kind'] ?? 'nothing') {
            case 'all':
                $filter = PlatformRequirementFilterFactory::ignoreAll();
                break;
            case 'list':
                $filter = PlatformRequirementFilterFactory::fromBoolOrList($spec['reqs'] ?? []);
                break;
            case 'php':
                if (isset(self::$filters[$spec['id']])) {
                    return self::$filters[$spec['id']];
                }
                // no break
            default:
                $filter = PlatformRequirementFilterFactory::ignoreNothing();
        }
        self::$goFilters[] = $filter;

        return $filter;
    }

    /**
     * @param array<string, mixed> $a
     */
    public static function filterCall(array $a): bool
    {
        $filter = self::$filters[$a['id']] ?? null;
        if ($filter === null) {
            throw new \LogicException('Unknown platform requirement filter');
        }

        return $a['m'] === 'isUpperBoundIgnored' ? $filter->isUpperBoundIgnored($a['req']) : $filter->isIgnored($a['req']);
    }

    /**
     * A dump gocomposer runs through the autoload generator a plugin set
     * with Composer::setAutoloadGenerator(): the setters Composer's callers
     * use, then dump() (Installer::run / DumpAutoloadCommand /
     * ReinstallCommand).
     *
     * @param array<string, mixed> $a
     * @return array<string, mixed>
     */
    public static function autoloadVia(array $a): array
    {
        $composer = Host::$composer;
        $generator = $composer->getAutoloadGenerator();
        if (!empty($a['dryRun'])) {
            $generator->setDryRun(true);
        }
        if ($a['dev'] !== null) {
            $generator->setDevMode((bool) $a['dev']);
        }
        $generator->setClassMapAuthoritative((bool) $a['authoritative']);
        $generator->setApcu((bool) $a['apcu'], $a['apcuPrefix']);
        $generator->setRunScripts((bool) $a['runScripts']);
        $generator->setPlatformRequirementFilter(self::filterFromGo($a['filter']));
        $classMap = $generator->dump(
            $composer->getConfig(),
            $composer->getRepositoryManager()->getLocalRepository(),
            $composer->getPackage(),
            $composer->getInstallationManager(),
            'composer',
            (bool) $a['optimize'],
            $a['suffix'],
            $composer->getLocker(),
            (bool) $a['strictAmbiguous']
        );

        return ['classes' => is_object($classMap) && method_exists($classMap, 'getMap') ? count($classMap->getMap()) : null];
    }
}

}
