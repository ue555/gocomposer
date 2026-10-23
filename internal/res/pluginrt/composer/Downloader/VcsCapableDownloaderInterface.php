<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Downloader/VcsCapableDownloaderInterface.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

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
 * VCS Capable Downloader interface.
 *
 * @author Steve Buzonas <steve@fancyguy.com>
 */
interface VcsCapableDownloaderInterface
{
    /**
     * Gets the VCS Reference for the package at path
     *
     * @param  PackageInterface $package package instance
     * @param  string           $path    package directory
     * @return string|null      reference or null
     */
    public function getVcsReference(PackageInterface $package, string $path): ?string;
}
