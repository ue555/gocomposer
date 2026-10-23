<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\EventDispatcher\EventDispatcher for
// plugins. Listeners registered here are kept in this process; dispatching is
// done by gocomposer's Go port of
// f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/EventDispatcher/EventDispatcher.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano), which merges them with
// the root package scripts and calls them back here.

namespace Composer\EventDispatcher;

use Composer\Composer;
use Composer\DependencyResolver\Operation\OperationInterface;
use Composer\DependencyResolver\Transaction;
use Composer\Installer\InstallerEvent;
use Composer\Installer\PackageEvent;
use Composer\IO\IOInterface;
use Composer\PartialComposer;
use Composer\Repository\RepositoryInterface;
use Composer\Script\Event as ScriptEvent;
use GocomposerPluginHost\Host;

class EventDispatcher
{
    /** @var PartialComposer */
    protected $composer;
    /** @var IOInterface */
    protected $io;
    /** @var bool */
    protected $runScripts = true;

    public function __construct(PartialComposer $composer, IOInterface $io, $process = null)
    {
        $this->composer = $composer;
        $this->io = $io;
    }

    /**
     * @return $this
     */
    public function setRunScripts(bool $runScripts = true): self
    {
        $this->runScripts = $runScripts;
        Host::cb('setRunScripts', ['run' => $runScripts]);

        return $this;
    }

    public function dispatch(?string $eventName, ?Event $event = null): int
    {
        if (null === $event) {
            if (null === $eventName) {
                throw new \InvalidArgumentException('If no $event is passed in to '.__METHOD__.' you have to pass in an $eventName, got null.');
            }
            $event = new Event($eventName);
        }

        return Host::dispatchFromPlugin($event);
    }

    /**
     * @param array<int, mixed>    $additionalArgs
     * @param array<string, mixed> $flags
     */
    public function dispatchScript(string $eventName, bool $devMode = false, array $additionalArgs = [], array $flags = []): int
    {
        assert($this->composer instanceof Composer, new \LogicException('This should only be reached with a fully loaded Composer'));

        return Host::dispatchFromPlugin(new ScriptEvent($eventName, $this->composer, $this->io, $devMode, $additionalArgs, $flags));
    }

    /**
     * @param OperationInterface[] $operations
     */
    public function dispatchPackageEvent(string $eventName, bool $devMode, RepositoryInterface $localRepo, array $operations, OperationInterface $operation): int
    {
        assert($this->composer instanceof Composer, new \LogicException('This should only be reached with a fully loaded Composer'));

        return Host::dispatchFromPlugin(new PackageEvent($eventName, $this->composer, $this->io, $devMode, $localRepo, $operations, $operation));
    }

    public function dispatchInstallerEvent(string $eventName, bool $devMode, bool $executeOperations, Transaction $transaction): int
    {
        assert($this->composer instanceof Composer, new \LogicException('This should only be reached with a fully loaded Composer'));

        return Host::dispatchFromPlugin(new InstallerEvent($eventName, $this->composer, $this->io, $devMode, $executeOperations, $transaction));
    }

    /**
     * @param callable|string $listener
     */
    public function addListener(string $eventName, $listener, int $priority = 0): void
    {
        Host::addListener($eventName, $listener, $priority);
    }

    /**
     * @param callable|object $listener
     */
    public function removeListener($listener): void
    {
        Host::removeListener($listener);
    }

    public function addSubscriber(EventSubscriberInterface $subscriber): void
    {
        foreach ($subscriber->getSubscribedEvents() as $eventName => $params) {
            if (is_string($params)) {
                $this->addListener($eventName, [$subscriber, $params]);
            } elseif (is_string($params[0])) {
                $this->addListener($eventName, [$subscriber, $params[0]], $params[1] ?? 0);
            } else {
                foreach ($params as $listener) {
                    $this->addListener($eventName, [$subscriber, $listener[0]], $listener[1] ?? 0);
                }
            }
        }
    }

    public function hasEventListeners(Event $event): bool
    {
        return (bool) Host::cb('hasListeners', ['name' => $event->getName()]);
    }
}
