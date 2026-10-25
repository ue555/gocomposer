<?php
namespace Plugfix3\Svc;

use Composer\Downloader\ZipDownloader;
use Composer\Package\PackageInterface;
use React\Promise\PromiseInterface;

/**
 * Replaces the zip downloader overriding extract(): logs, extracts with the
 * parent, then adds a marker file to the extracted package.
 */
class ExtractingZip extends ZipDownloader
{
    protected function extract(PackageInterface $package, string $file, string $path): PromiseInterface
    {
        $this->io->write('[extract] ' . $package->getPrettyName() . ' file=' . (preg_match('{/composer/tmp-[0-9a-f]{32}\.zip$}', $file) ? 'tmp.zip' : $file) . ' into vendor/composer/' . (preg_match('{/composer/[0-9a-f]{8}$}', $path) ? 'TMP' : $path) . ' exists=' . var_export(is_file($file), true));
        if (getenv('PLUGFIX3_FAIL') === $package->getName()) {
            throw new \RuntimeException('plugfix3 refuses to extract ' . $package->getName());
        }

        return parent::extract($package, $file, $path)->then(function () use ($package, $path) {
            $dirs = glob($path . '/*', GLOB_ONLYDIR);
            $target = count($dirs) === 1 ? $dirs[0] : $path;
            file_put_contents($target . '/PLUGFIX3-EXTRACTED', $package->getPrettyName() . "\n");
            $this->io->write('[extract] done ' . $package->getPrettyName());
        });
    }
}
