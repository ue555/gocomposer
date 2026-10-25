<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Package\Archiver\PharArchiver, standing
// in for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Package/Archiver/PharArchiver.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). The archive is written by gocomposer's port of
// the archiver (internal/archiver).

namespace Composer\Package\Archiver;

use GocomposerPluginHost\Host;

class PharArchiver implements ArchiverInterface
{
    /** @var array<string, int> */
    protected static $formats = [
        'zip' => 1,
        'tar' => 2,
        'tar.gz' => 2,
        'tar.bz2' => 2,
    ];

    /**
     * @inheritDoc
     */
    public function archive(string $sources, string $target, string $format, array $excludes = [], bool $ignoreFilters = false): string
    {
        return (string) Host::cb('archiveRun', ['sources' => $sources, 'target' => $target, 'format' => $format, 'excludes' => array_values($excludes), 'ignoreFilters' => $ignoreFilters]);
    }

    /**
     * @inheritDoc
     */
    public function supports(string $format, ?string $sourceType): bool
    {
        return isset(static::$formats[$format]);
    }
}
