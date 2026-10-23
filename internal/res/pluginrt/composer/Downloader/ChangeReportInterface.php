<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/ChangeReportInterface.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\Downloader;

use Composer\Package\PackageInterface;

/**
 * ChangeReport interface.
 *
 * @author Sascha Egerer <sascha.egerer@dkd.de>
 */
interface ChangeReportInterface
{
    /**
     * Checks for changes to the local copy
     *
     * @param  PackageInterface $package package instance
     * @param  string           $path    package directory
     * @return string|null      changes or null
     */
    public function getLocalChanges(PackageInterface $package, string $path): ?string;
}
