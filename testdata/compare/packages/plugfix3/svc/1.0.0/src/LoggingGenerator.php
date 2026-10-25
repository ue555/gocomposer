<?php
namespace Plugfix3\Svc;

use Composer\Autoload\AutoloadGenerator;
use Composer\Config;
use Composer\Installer\InstallationManager;
use Composer\Package\Locker;
use Composer\Package\RootPackageInterface;
use Composer\Repository\InstalledRepositoryInterface;

/**
 * Replaces the autoload generator: logs the dumps, then generates like
 * Composer's (PLUGFIX3_AL_SKIP: does not dump at all).
 */
class LoggingGenerator extends AutoloadGenerator
{
    private $logIo;

    public function __construct($eventDispatcher, $io)
    {
        parent::__construct($eventDispatcher, $io);
        $this->logIo = $io;
    }

    public function dump(Config $config, InstalledRepositoryInterface $localRepo, RootPackageInterface $rootPackage, InstallationManager $installationManager, string $targetDir, bool $scanPsrPackages = false, ?string $suffix = null, ?Locker $locker = null, bool $strictAmbiguous = false)
    {
        $this->logIo->write('[al] dump ' . $targetDir . ' root=' . $rootPackage->getName() . ' installed=' . count($localRepo->getPackages()) . ' scanPsr=' . var_export($scanPsrPackages, true) . ' suffix=' . var_export($suffix, true) . ' locker=' . var_export($locker !== null, true) . ' strictAmbiguous=' . var_export($strictAmbiguous, true));
        if (getenv('PLUGFIX3_AL_SKIP')) {
            return new \Composer\ClassMapGenerator\ClassMap();
        }
        $map = parent::dump($config, $localRepo, $rootPackage, $installationManager, $targetDir, $scanPsrPackages, $suffix, $locker, $strictAmbiguous);
        $this->logIo->write('[al] dumped classes=' . count($map->getMap()));

        return $map;
    }
}
