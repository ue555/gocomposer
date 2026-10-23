<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Repository\PackageRepository, standing in for
// Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/PackageRepository.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). Instances stand for gocomposer's "package" repositories
// (RepositoryManager::getRepositories()/createRepository()); queries are
// answered by the Go implementation. Constructing one creates the Go
// repository from the given repository configuration.

namespace Composer\Repository;

use GocomposerPluginHost\GoRepository;

class PackageRepository extends GoRepository
{
    protected const GC_TYPE = 'package';
}
