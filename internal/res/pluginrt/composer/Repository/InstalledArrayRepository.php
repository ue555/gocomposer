<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/InstalledArrayRepository.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\Repository;

/**
 * Installed array repository.
 *
 * This is used as an in-memory InstalledRepository mostly for testing purposes
 *
 * @author Jordi Boggiano <j.boggiano@seld.be>
 */
class InstalledArrayRepository extends WritableArrayRepository implements InstalledRepositoryInterface
{
    public function getRepoName(): string
    {
        return 'installed '.parent::getRepoName();
    }

    /**
     * @inheritDoc
     */
    public function isFresh(): bool
    {
        // this is not a completely correct implementation but there is no way to
        // distinguish an empty repo and a newly created one given this is all in-memory
        return $this->count() === 0;
    }
}
