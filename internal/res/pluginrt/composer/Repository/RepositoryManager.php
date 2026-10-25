<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Repository\RepositoryManager for
// plugins, based on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/RepositoryManager.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). The repositories are gocomposer's (composer,
// vcs, path, artifact, package): the objects returned here answer queries
// through the Go implementation, and repositories created or added by
// plugins are added to gocomposer's repository list. Repository classes
// registered by plugins (setRepositoryClass) are created like Composer
// does; gocomposer asks the repository objects of plugin code (and the
// ones it changed) when resolving.

namespace Composer\Repository;

use Composer\Package\PackageInterface;
use GocomposerPluginHost\Api;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Services;

class RepositoryManager
{
    /** @var InstalledRepositoryInterface */
    private $localRepository;
    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var array<string, class-string<RepositoryInterface>> classes registered by plugins */
    private $gcClasses = [];

    public function __construct($io = null, $config = null, $httpDownloader = null, $eventDispatcher = null, $process = null)
    {
        if ($config instanceof \Composer\Config) {
            $this->gcInst = $config->gcInst;
        }
    }

    /**
     * @param  string                                                 $name
     * @param  string|\Composer\Semver\Constraint\ConstraintInterface $constraint
     */
    public function findPackage(string $name, $constraint): ?PackageInterface
    {
        foreach ($this->getRepositories() as $repository) {
            /** @var RepositoryInterface $repository */
            if ($package = $repository->findPackage($name, $constraint)) {
                return $package;
            }
        }

        return null;
    }

    /**
     * @param  string|\Composer\Semver\Constraint\ConstraintInterface $constraint
     * @return PackageInterface[]
     */
    public function findPackages(string $name, $constraint): array
    {
        $packages = [];

        foreach ($this->getRepositories() as $repository) {
            $packages = array_merge($packages, $repository->findPackages($name, $constraint));
        }

        return $packages;
    }

    public function addRepository(RepositoryInterface $repository): void
    {
        $this->gcAdd($repository, false);
    }

    public function prependRepository(RepositoryInterface $repository): void
    {
        $this->gcAdd($repository, true);
    }

    private function gcAdd(RepositoryInterface $repository, bool $prepend): void
    {
        // one of gocomposer's repositories, or an object of plugin code
        // gocomposer asks (plugin_repos.go)
        Host::cb('repoAdd', ['inst' => $this->gcInst, 'prepend' => $prepend] + Services::repoSpec($repository));
    }

    /**
     * @param array<string, mixed> $config
     */
    public function createRepository(string $type, array $config, ?string $name = null): RepositoryInterface
    {
        if (isset($this->gcClasses[$type])) {
            return Services::createRepository($this->gcInst, $this->gcClasses[$type], $config, $name);
        }
        $desc = Host::cb('repoCreate', ['inst' => $this->gcInst, 'type' => $type, 'config' => $config, 'name' => $name]);

        return Api::repo($desc);
    }

    /**
     * @param class-string<RepositoryInterface> $class
     */
    public function setRepositoryClass(string $type, $class): void
    {
        $this->gcClasses[$type] = $class;
    }

    /**
     * @return RepositoryInterface[]
     */
    public function getRepositories(): array
    {
        $out = [];
        foreach (Host::cb('repos', ['inst' => $this->gcInst]) as $desc) {
            $out[] = Api::repo($desc);
        }

        return $out;
    }

    public function setLocalRepository(InstalledRepositoryInterface $repository): void
    {
        $this->localRepository = $repository;
    }

    public function getLocalRepository(): InstalledRepositoryInterface
    {
        return $this->localRepository;
    }

    public function getHttpDownloader()
    {
        return Api::httpDownloader($this->gcInst);
    }
}
