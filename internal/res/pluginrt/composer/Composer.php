<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Composer for plugins, based on
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Composer.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The loop, download
// manager and autoload generator are gocomposer-backed stubs; the archive
// manager and replacing services fail explicitly.

namespace Composer;

use Composer\Package\Locker;
use Composer\Pcre\Preg;
use Composer\Plugin\PluginManager;
use GocomposerPluginHost\Unsupported;

class Composer extends PartialComposer
{
    public const VERSION = '@package_version@';
    public const BRANCH_ALIAS_VERSION = '@package_branch_alias_version@';
    public const RELEASE_DATE = '@release_date@';
    public const SOURCE_VERSION = '2.11.999-dev+source';

    public const RUNTIME_API_VERSION = '2.2.2';

    public static function getVersion(): string
    {
        // no replacement done, this must be a source checkout
        if (self::VERSION === '@package_version'.'@') {
            return self::SOURCE_VERSION;
        }

        // we have a branch alias and version is a commit id, this must be a snapshot build
        if (self::BRANCH_ALIAS_VERSION !== '' && Preg::isMatch('{^[a-f0-9]{40}$}', self::VERSION)) {
            return self::BRANCH_ALIAS_VERSION.'+'.self::VERSION;
        }

        return self::VERSION;
    }

    /** @var Locker */
    private $locker;
    /** @var PluginManager */
    private $pluginManager;

    public function setLocker(Locker $locker): void
    {
        $this->locker = $locker;
    }

    public function getLocker(): Locker
    {
        return $this->locker;
    }

    public function setPluginManager(PluginManager $manager): void
    {
        $this->pluginManager = $manager;
    }

    public function getPluginManager(): PluginManager
    {
        return $this->pluginManager;
    }

    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var \Composer\Util\Loop|null */
    private $loop;
    /** @var \Composer\Downloader\DownloadManager|null */
    private $downloadManager;
    /** @var \Composer\Autoload\AutoloadGenerator|null */
    private $autoloadGenerator;

    public function getLoop()
    {
        if ($this->loop === null) {
            $http = new \Composer\Util\HttpDownloader(\GocomposerPluginHost\Host::$io, $this->getConfig());
            $this->loop = new \Composer\Util\Loop($http, new \Composer\Util\ProcessExecutor(\GocomposerPluginHost\Host::$io));
        }

        return $this->loop;
    }

    public function getDownloadManager()
    {
        if ($this->downloadManager === null) {
            $this->downloadManager = \GocomposerPluginHost\Api::downloadManager($this);
        }

        return $this->downloadManager;
    }

    /**
     * @internal the download manager when a plugin already used it
     */
    public function gcExistingDownloadManager(): ?\Composer\Downloader\DownloadManager
    {
        return $this->downloadManager;
    }

    public function getAutoloadGenerator()
    {
        if ($this->autoloadGenerator === null) {
            $this->autoloadGenerator = new \Composer\Autoload\AutoloadGenerator($this->getEventDispatcher(), \GocomposerPluginHost\Host::$io);
            $this->autoloadGenerator->gcInst = $this->gcInst;
        }

        return $this->autoloadGenerator;
    }

    public function getArchiveManager()
    {
        Unsupported::fail(__METHOD__);
    }

    public function setLoop($loop): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function setDownloadManager($manager): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function setArchiveManager($manager): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function setAutoloadGenerator($autoloadGenerator): void
    {
        Unsupported::fail(__METHOD__);
    }
}
