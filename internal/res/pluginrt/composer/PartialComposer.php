<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\PartialComposer for plugins, based on
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/PartialComposer.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The services are
// gocomposer stubs backed by the Go process.

namespace Composer;

use Composer\EventDispatcher\EventDispatcher;
use Composer\Installer\InstallationManager;
use Composer\Package\RootPackageInterface;
use Composer\Repository\RepositoryManager;

class PartialComposer
{
    /** @var bool */
    private $global = false;
    /** @var RootPackageInterface */
    private $package;
    /** @var Config */
    private $config;
    /** @var RepositoryManager */
    private $repositoryManager;
    /** @var InstallationManager */
    private $installationManager;
    /** @var EventDispatcher */
    private $eventDispatcher;

    public function setPackage(RootPackageInterface $package): void
    {
        $this->package = $package;
    }

    public function getPackage(): RootPackageInterface
    {
        return $this->package;
    }

    public function setConfig(Config $config): void
    {
        $this->config = $config;
    }

    public function getConfig(): Config
    {
        return $this->config;
    }

    public function setRepositoryManager(RepositoryManager $manager): void
    {
        $this->repositoryManager = $manager;
    }

    public function getRepositoryManager(): RepositoryManager
    {
        return $this->repositoryManager;
    }

    public function setInstallationManager(InstallationManager $manager): void
    {
        $this->installationManager = $manager;
    }

    public function getInstallationManager(): InstallationManager
    {
        return $this->installationManager;
    }

    public function setEventDispatcher(EventDispatcher $eventDispatcher): void
    {
        $this->eventDispatcher = $eventDispatcher;
    }

    public function getEventDispatcher(): EventDispatcher
    {
        return $this->eventDispatcher;
    }

    public function isGlobal(): bool
    {
        return $this->global;
    }

    public function setGlobal(): void
    {
        $this->global = true;
    }
}
