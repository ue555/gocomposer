<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Package/Archiver/ArchivableFilesFilter.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\Package\Archiver;

use FilterIterator;
use Iterator;
use PharData;
use SplFileInfo;

/**
 * @phpstan-extends FilterIterator<string, SplFileInfo, Iterator<string, SplFileInfo>>
 */
class ArchivableFilesFilter extends FilterIterator
{
    /** @var string[] */
    private $dirs = [];

    /**
     * @return bool true if the current element is acceptable, otherwise false.
     */
    public function accept(): bool
    {
        $file = $this->getInnerIterator()->current();
        if ($file->isDir()) {
            $this->dirs[] = (string) $file;

            return false;
        }

        return true;
    }

    public function addEmptyDir(PharData $phar, string $sources): void
    {
        foreach ($this->dirs as $filepath) {
            $localname = str_replace($sources . "/", '', $filepath);
            $phar->addEmptyDir($localname);
        }
    }
}
