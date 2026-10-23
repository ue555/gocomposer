<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\DependencyResolver\LocalRepoTransaction
// (Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/DependencyResolver/LocalRepoTransaction.php, MIT license,
// (c) Nils Adermann, Jordi Boggiano). The transactions of gocomposer's
// installer events carry the operations computed by gocomposer together
// with the present (installed) and result (locked) packages, so that code
// reading Transaction's state (e.g. symfony/flex building its own
// Transaction from resultPackageMap) sees Composer's.

namespace Composer\DependencyResolver;

use Composer\DependencyResolver\Operation\OperationInterface;
use Composer\Package\PackageInterface;
use Composer\Repository\InstalledRepositoryInterface;
use Composer\Repository\RepositoryInterface;

class LocalRepoTransaction extends Transaction
{
    public function __construct(RepositoryInterface $lockedRepository, InstalledRepositoryInterface $localRepository)
    {
        parent::__construct(
            $localRepository->getPackages(),
            $lockedRepository->getPackages()
        );
    }

    /**
     * @internal
     * @param list<OperationInterface> $operations computed by gocomposer
     * @param list<PackageInterface>   $present
     * @param list<PackageInterface>   $result
     */
    public static function gcCreate(array $operations, array $present, array $result): self
    {
        $transaction = (new \ReflectionClass(self::class))->newInstanceWithoutConstructor();
        \Closure::bind(function () use ($operations, $present, $result): void {
            $this->presentPackages = $present;
            $this->setResultPackageMaps($result);
            $this->operations = $operations;
        }, $transaction, Transaction::class)();

        return $transaction;
    }
}
