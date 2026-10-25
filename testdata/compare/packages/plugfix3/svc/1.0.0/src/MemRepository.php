<?php
namespace Plugfix3\Svc;

use Composer\Package\Loader\ArrayLoader;
use Composer\Repository\ArrayRepository;

/**
 * A repository type registered by the plugin: the packages of its
 * "packages" configuration, logging the loadPackages() calls.
 */
class MemRepository extends ArrayRepository
{
    private $io;
    private $name;

    public function __construct(array $repoConfig, $io, $config, $httpDownloader = null, $eventDispatcher = null, $process = null)
    {
        $this->io = $io;
        $this->name = $repoConfig['name'] ?? 'mem';
        parent::__construct();
        $loader = new ArrayLoader();
        foreach ($repoConfig['packages'] ?? array() as $data) {
            $this->addPackage($loader->load($data));
        }
        $io->write('[mem] created ' . $this->name . ' with ' . count($this->getPackages()) . ' packages config=' . get_class($config));
    }

    public function getRepoName()
    {
        return 'plugfix3 mem repo ' . $this->name;
    }

    public function loadPackages(array $packageNameMap, array $acceptableStabilities, array $stabilityFlags, array $alreadyLoaded = array())
    {
        $names = array_keys($packageNameMap);
        sort($names);
        $result = parent::loadPackages($packageNameMap, $acceptableStabilities, $stabilityFlags, $alreadyLoaded);
        $found = array();
        foreach ($result['packages'] as $p) {
            $found[] = $p->getPrettyName() . '@' . $p->getPrettyVersion();
        }
        if ($found) {
            $this->io->write('[mem] loadPackages found ' . implode(',', $found));
        }

        return $result;
    }
}
