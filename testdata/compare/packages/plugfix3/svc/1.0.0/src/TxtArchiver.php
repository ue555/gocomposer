<?php
namespace Plugfix3\Svc;

use Composer\Package\Archiver\ArchivableFilesFinder;
use Composer\Package\Archiver\ArchiverInterface;

/**
 * An archiver for a format of its own ("txt"): the list of the archivable
 * files.
 */
class TxtArchiver implements ArchiverInterface
{
    public function archive(string $sources, string $target, string $format, array $excludes = array(), bool $ignoreFilters = false): string
    {
        $sources = realpath($sources);
        $files = array();
        foreach (new ArchivableFilesFinder($sources, $excludes, $ignoreFilters) as $file) {
            $files[] = substr($file->getPathname(), strlen($sources) + 1);
        }
        sort($files);
        file_put_contents($target, "plugfix3 txt archive\n" . implode("\n", $files) . "\n");

        return $target;
    }

    public function supports(string $format, ?string $sourceType): bool
    {
        return $format === 'txt';
    }
}
