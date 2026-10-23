<?php declare(strict_types=1);
// gocomposer plugin runtime: read-only Composer\Package\Locker for plugins,
// based on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Package/Locker.php (MIT license, (c) Nils Adermann, Jordi
// Boggiano). The Composer instance's locker reads the lock file through
// gocomposer (which writes it); a Locker constructed by plugin code with a
// JsonFile reads that file itself like Composer's. getLockedRepository() is
// ported unchanged.

namespace Composer\Package;

use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class Locker
{
    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var \Composer\Json\JsonFile|null */
    private $lockFile;
    /** @var array<string, mixed>|null */
    private $lockDataCache;

    public function __construct($io = null, $lockFile = null, $installationManager = null, string $composerFileContents = '', $process = null)
    {
        $this->lockFile = $lockFile;
    }

    public static function getContentHash(string $composerFileContents): string
    {
        return (string) Host::cb('contentHash', ['json' => $composerFileContents]);
    }

    public function isLocked(): bool
    {
        if ($this->lockFile !== null) {
            if (!$this->lockFile->exists()) {
                return false;
            }
            $data = $this->getLockData();

            return isset($data['packages']);
        }

        return (bool) Host::cb('lock', ['m' => 'isLocked', 'inst' => $this->gcInst]);
    }

    public function isFresh(): bool
    {
        if ($this->lockFile !== null) {
            Unsupported::fail(__METHOD__);
        }

        return (bool) Host::cb('lock', ['m' => 'isFresh', 'inst' => $this->gcInst]);
    }

    /**
     * @return array<string, mixed>
     */
    public function getLockData(): array
    {
        if ($this->lockFile !== null) {
            if (null !== $this->lockDataCache) {
                return $this->lockDataCache;
            }

            if (!$this->lockFile->exists()) {
                throw new \LogicException('No lockfile found. Unable to read locked packages');
            }

            return $this->lockDataCache = $this->lockFile->read();
        }
        $data = Host::cb('lock', ['m' => 'data', 'inst' => $this->gcInst]);
        if (!is_array($data)) {
            throw new \LogicException('No lockfile found. Unable to read locked packages');
        }

        return $data;
    }

    public function getPluginApi()
    {
        $data = $this->getLockData();

        return $data['plugin-api-version'] ?? '1.1.0';
    }

    /**
     * @return string[]
     */
    public function getDevPackageNames(): array
    {
        $names = [];
        $lockData = $this->getLockData();
        if (isset($lockData['packages-dev'])) {
            foreach ($lockData['packages-dev'] as $package) {
                $names[] = strtolower($package['name']);
            }
        }

        return $names;
    }

    public function getMinimumStability(): string
    {
        $lockData = $this->getLockData();

        return $lockData['minimum-stability'] ?? 'stable';
    }

    /**
     * @return array<string, string>
     */
    public function getStabilityFlags(): array
    {
        $lockData = $this->getLockData();

        return $lockData['stability-flags'] ?? [];
    }

    public function getPreferStable(): ?bool
    {
        $lockData = $this->getLockData();

        return $lockData['prefer-stable'] ?? null;
    }

    public function getPreferLowest(): ?bool
    {
        $lockData = $this->getLockData();

        return $lockData['prefer-lowest'] ?? null;
    }

    /**
     * @return array<string, string>
     */
    public function getPlatformOverrides(): array
    {
        $lockData = $this->getLockData();

        return $lockData['platform-overrides'] ?? [];
    }

    /**
     * @return string[][]
     */
    public function getAliases(): array
    {
        $lockData = $this->getLockData();

        return $lockData['aliases'] ?? [];
    }

    public function getJsonFile()
    {
        if ($this->lockFile !== null) {
            return $this->lockFile;
        }
        Unsupported::fail(__METHOD__);
    }

    public function getLockedRepository(bool $withDevReqs = false): \Composer\Repository\LockArrayRepository
    {
        $lockData = $this->getLockData();
        $packages = new \Composer\Repository\LockArrayRepository();

        $lockedPackages = $lockData['packages'];
        if ($withDevReqs) {
            if (isset($lockData['packages-dev'])) {
                $lockedPackages = array_merge($lockedPackages, $lockData['packages-dev']);
            } else {
                throw new \RuntimeException('The lock file does not contain require-dev information, run install with the --no-dev option or delete it and run composer update to generate a new lock file.');
            }
        }

        if (empty($lockedPackages)) {
            return $packages;
        }

        if (isset($lockedPackages[0]['name'])) {
            $loader = new \Composer\Package\Loader\ArrayLoader(null, true);
            $packageByName = [];
            foreach ($lockedPackages as $info) {
                $package = $loader->load($info);
                $packages->addPackage($package);
                $packageByName[$package->getName()] = $package;

                if ($package instanceof AliasPackage) {
                    $packageByName[$package->getAliasOf()->getName()] = $package->getAliasOf();
                }
            }

            if (isset($lockData['aliases'])) {
                foreach ($lockData['aliases'] as $alias) {
                    if (isset($packageByName[$alias['package']])) {
                        $aliasPkg = new CompleteAliasPackage($packageByName[$alias['package']], $alias['alias_normalized'], $alias['alias']);
                        $aliasPkg->setRootPackageAlias(true);
                        $packages->addPackage($aliasPkg);
                    }
                }
            }

            return $packages;
        }

        throw new \RuntimeException('Your composer.lock is invalid. Run "composer update" to generate a new one.');
    }

    public function getPlatformRequirements(bool $withDevReqs = false): array
    {
        Unsupported::fail(__METHOD__);
    }

    public function setLockData(...$args): bool
    {
        Unsupported::fail(__METHOD__);
    }

    public function updateHash($composerJson, ?callable $dataProcessor = null): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function getMissingRequirementInfo(RootPackageInterface $package, bool $includeDev): array
    {
        Unsupported::fail(__METHOD__);
    }
}
