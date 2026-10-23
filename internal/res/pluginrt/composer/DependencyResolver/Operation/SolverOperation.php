<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/DependencyResolver/Operation/SolverOperation.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\DependencyResolver\Operation;

/**
 * Abstract operation class.
 *
 * @author Aleksandr Bezpiatov <aleksandr.bezpiatov@spryker.com>
 */
abstract class SolverOperation implements OperationInterface
{
    /**
     * @abstract must be redefined by extending classes
     */
    protected const TYPE = '';

    /**
     * Returns operation type.
     */
    public function getOperationType(): string
    {
        return static::TYPE;
    }

    /**
     * @inheritDoc
     */
    public function __toString()
    {
        return $this->show(false);
    }
}
