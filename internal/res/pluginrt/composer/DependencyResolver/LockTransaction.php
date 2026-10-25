<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\DependencyResolver\LockTransaction
// (Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/DependencyResolver/LockTransaction.php, MIT license, (c) Nils
// Adermann, Jordi Boggiano); the lock transactions are gocomposer's: an
// object carries the operations, present and result packages of one, and
// the packages it locks (getNewLockPackages).

namespace Composer\DependencyResolver;

use Composer\Package\PackageInterface;
use GocomposerPluginHost\Api;

class LockTransaction extends Transaction
{
    /** @var array<string, list<PackageInterface>> */
    private $gcLockPackages = [];

    /**
     * @internal
     * @param array<string, mixed> $data gocomposer's lock transaction
     */
    public static function gcCreate(array $data): self
    {
        $ops = [];
        foreach ($data['ops'] ?? [] as $spec) {
            $ops[] = \GocomposerPluginHost\Host::opFromSpec($spec);
        }
        $present = Api::pkgs($data['present'] ?? []);
        $result = Api::pkgs($data['result'] ?? []);
        $transaction = (new \ReflectionClass(self::class))->newInstanceWithoutConstructor();
        \Closure::bind(function () use ($ops, $present, $result): void {
            $this->presentPackages = $present;
            $this->setResultPackageMaps($result);
            $this->operations = $ops;
        }, $transaction, Transaction::class)();
        foreach (['nonDev', 'dev', 'nonDevMirrors', 'devMirrors'] as $k) {
            $transaction->gcLockPackages[$k] = Api::pkgs($data[$k] ?? []);
        }

        return $transaction;
    }

    /**
     * @return BasePackage[]
     */
    public function getNewLockPackages(bool $devMode, bool $updateMirrors = false): array
    {
        return $this->gcLockPackages[($devMode ? 'dev' : 'nonDev').($updateMirrors ? 'Mirrors' : '')];
    }
}
