<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Composer for plugins, based on
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Composer.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The loop, download
// manager and autoload generator are gocomposer-backed stubs. Plugins may
// replace the services (set*()): like with Composer, the replacement is what
// the getters return from then on; gocomposer calls a replaced autoload
// generator where Composer's commands get it from the Composer object
// (Installer::create, dump-autoload, reinstall).

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
    /** @var \Composer\Package\Archiver\ArchiveManager|null */
    private $archiveManager;

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
            $this->downloadManager = $this->gcInstallerDownloadManager = \GocomposerPluginHost\Api::downloadManager($this);
        }

        return $this->downloadManager;
    }

    /** @var \Composer\Downloader\DownloadManager|null */
    private $gcInstallerDownloadManager;

    /**
     * @internal the download manager of the installers (created with the
     * Composer instance, the one a plugin replacing it with
     * setDownloadManager() leaves to the installers), when a plugin already
     * used it
     */
    public function gcExistingDownloadManager(): ?\Composer\Downloader\DownloadManager
    {
        return $this->gcInstallerDownloadManager;
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
        if ($this->archiveManager === null) {
            $this->archiveManager = \GocomposerPluginHost\Services::archiveManager($this);
        }

        return $this->archiveManager;
    }

    public function setLoop(\Composer\Util\Loop $loop): void
    {
        $this->loop = $loop;
    }

    public function setDownloadManager(\Composer\Downloader\DownloadManager $manager): void
    {
        $this->downloadManager = $manager;
        \GocomposerPluginHost\Host::cb('serviceSet', ['inst' => $this->gcInst, 'name' => 'downloadManager']);
    }

    public function setArchiveManager(\Composer\Package\Archiver\ArchiveManager $manager): void
    {
        $this->archiveManager = $manager;
        \GocomposerPluginHost\Host::cb('serviceSet', ['inst' => $this->gcInst, 'name' => 'archiveManager']);
    }

    public function setAutoloadGenerator(\Composer\Autoload\AutoloadGenerator $autoloadGenerator): void
    {
        $this->autoloadGenerator = $autoloadGenerator;
        \GocomposerPluginHost\Host::cb('serviceSet', ['inst' => $this->gcInst, 'name' => 'autoloadGenerator']);
    }

    /**
     * @internal the archive manager when a plugin already used it
     */
    public function gcExistingArchiveManager(): ?\Composer\Package\Archiver\ArchiveManager
    {
        return $this->archiveManager;
    }
}
