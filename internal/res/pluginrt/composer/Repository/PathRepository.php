<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Repository\PathRepository, standing in for
// Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/PathRepository.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). Instances stand for gocomposer's "path" repositories
// (RepositoryManager::getRepositories()/createRepository()); queries are
// answered by the Go implementation. Constructing one creates the Go
// repository from the given repository configuration.

namespace Composer\Repository;

use GocomposerPluginHost\GoRepository;

class PathRepository extends GoRepository implements ConfigurableRepositoryInterface
{
    protected const GC_TYPE = 'path';
}
