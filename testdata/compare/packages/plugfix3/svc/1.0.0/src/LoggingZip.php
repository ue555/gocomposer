<?php
namespace Plugfix3\Svc;

use Composer\Downloader\ZipDownloader;
use Composer\Package\PackageInterface;
use React\Promise\PromiseInterface;

/**
 * Replaces the zip downloader: logs every call, then does what the zip
 * downloader does. PLUGFIX3_FAIL=<package> makes install() of that package
 * throw.
 */
class LoggingZip extends ZipDownloader
{
    private function log($what)
    {
        $this->io->write('[zip] ' . $what);
    }

    public function download(PackageInterface $package, string $path, ?PackageInterface $prevPackage = null, bool $output = true): PromiseInterface
    {
        $this->log('download ' . $package->getPrettyName() . ' ' . $package->getPrettyVersion() . ($prevPackage ? ' from ' . $prevPackage->getPrettyVersion() : '') . ' output=' . var_export($output, true));

        return parent::download($package, $path, $prevPackage, $output);
    }

    public function prepare(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        $this->log('prepare ' . $type . ' ' . $package->getPrettyName());

        return parent::prepare($type, $package, $path, $prevPackage);
    }

    public function install(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        $this->log('install ' . $package->getPrettyName() . ' ' . $package->getPrettyVersion() . ' output=' . var_export($output, true));
        if (getenv('PLUGFIX3_FAIL') === $package->getName()) {
            throw new \RuntimeException('plugfix3 refuses to install ' . $package->getName());
        }

        return parent::install($package, $path, $output);
    }

    public function update(PackageInterface $initial, PackageInterface $target, string $path): PromiseInterface
    {
        $this->log('update ' . $target->getPrettyName() . ' ' . $initial->getPrettyVersion() . ' -> ' . $target->getPrettyVersion());

        return parent::update($initial, $target, $path);
    }

    public function remove(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        $this->log('remove ' . $package->getPrettyName() . ' ' . $package->getPrettyVersion() . ' output=' . var_export($output, true));

        return parent::remove($package, $path, $output);
    }

    public function cleanup(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        $this->log('cleanup ' . $type . ' ' . $package->getPrettyName());

        return parent::cleanup($type, $package, $path, $prevPackage);
    }
}
