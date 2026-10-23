<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Filter/PlatformRequirementFilter/IgnoreAllPlatformRequirementFilter.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\Filter\PlatformRequirementFilter;

use Composer\Repository\PlatformRepository;

final class IgnoreAllPlatformRequirementFilter implements PlatformRequirementFilterInterface
{
    public function isIgnored(string $req): bool
    {
        return PlatformRepository::isPlatformPackage($req);
    }

    public function isUpperBoundIgnored(string $req): bool
    {
        return $this->isIgnored($req);
    }
}
