<?php
namespace Plugfix\Events;

use Composer\Composer;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\Installer\InstallerEvent;
use Composer\Installer\PackageEvent;
use Composer\IO\IOInterface;
use Composer\Plugin\Capable;
use Composer\Plugin\CommandEvent;
use Composer\Plugin\PluginInterface;
use Composer\Plugin\PreCommandRunEvent;
use Composer\Script\Event;

class Plugin implements PluginInterface, EventSubscriberInterface, Capable
{
    const VERSION = '1.1.0';
    private $composer;
    private $io;
    private $count = 0;

    public function activate(Composer $composer, IOInterface $io)
    {
        $this->composer = $composer;
        $this->io = $io;
        $io->write('<info>[plugfix/events] activate ' . self::VERSION . '</info>');
        $io->write('[plugfix/events] verbose activate detail', true, IOInterface::VERBOSE);
    }

    public function deactivate(Composer $composer, IOInterface $io)
    {
        $io->write('[plugfix/events] deactivate ' . self::VERSION);
    }

    public function uninstall(Composer $composer, IOInterface $io)
    {
        $io->write('[plugfix/events] uninstall ' . self::VERSION);
    }

    public function getCapabilities()
    {
        return array('Composer\Plugin\Capability\CommandProvider' => 'Plugfix\Events\CommandProvider');
    }

    public static function getSubscribedEvents()
    {
        return array(
            'init' => 'onInit',
            'command' => 'onCommand',
            'pre-command-run' => 'onPreCommandRun',
            'pre-install-cmd' => 'onScript',
            'post-install-cmd' => array('onScript', 10),
            'pre-update-cmd' => 'onScript',
            'post-update-cmd' => array(array('onScript', 5), array('onLate', -5)),
            'pre-autoload-dump' => 'onScript',
            'post-autoload-dump' => 'onScript',
            'pre-operations-exec' => 'onOperations',
            'pre-package-install' => 'onPackage',
            'post-package-install' => 'onPackage',
            'pre-package-update' => 'onPackage',
            'post-package-update' => 'onPackage',
            'pre-package-uninstall' => 'onPackage',
            'post-package-uninstall' => 'onPackage',
            'plugfix-custom' => 'onCustom',
        );
    }

    public function onInit(\Composer\EventDispatcher\Event $event)
    {
        $this->io->write('[plugfix/events] init ' . get_class($event));
    }

    public function onCommand(CommandEvent $event)
    {
        $input = $event->getInput();
        $this->io->write('[plugfix/events] command ' . $event->getCommandName() . ' dry-run=' . var_export($input->hasOption('dry-run') ? $input->getOption('dry-run') : null, true));
    }

    public function onPreCommandRun(PreCommandRunEvent $event)
    {
        $this->io->write('[plugfix/events] pre-command-run ' . $event->getCommand());
    }

    public function onScript(Event $event)
    {
        $this->count++;
        $flags = $event->getFlags();
        ksort($flags);
        $this->io->write(sprintf('[plugfix/events] %s dev=%s count=%d flags=%s args=%s', $event->getName(), var_export($event->isDevMode(), true), $this->count, json_encode($flags), json_encode($event->getArguments())));
        $root = $event->getComposer()->getPackage();
        $this->io->write('[plugfix/events] root ' . $root->getName() . ' ' . $root->getPrettyVersion() . ' vendor-dir=' . basename($event->getComposer()->getConfig()->get('vendor-dir')));
    }

    public function onLate(Event $event)
    {
        $names = array();
        foreach ($event->getComposer()->getRepositoryManager()->getLocalRepository()->getCanonicalPackages() as $p) {
            $names[] = $p->getPrettyName() . '@' . $p->getPrettyVersion();
        }
        sort($names);
        $this->io->write('[plugfix/events] late ' . $event->getName() . ' installed=' . implode(',', $names));
    }

    public function onOperations(InstallerEvent $event)
    {
        $ops = array();
        foreach ($event->getTransaction()->getOperations() as $op) {
            $ops[] = $op->getOperationType() . ':' . $op->show(false);
        }
        $this->io->write('[plugfix/events] pre-operations-exec exec=' . var_export($event->isExecutingOperations(), true) . ' dev=' . var_export($event->isDevMode(), true) . ' ' . implode(' | ', $ops));
    }

    public function onPackage(PackageEvent $event)
    {
        $op = $event->getOperation();
        $package = $op->getOperationType() === 'update' ? $op->getTargetPackage() : $op->getPackage();
        $path = $event->getComposer()->getInstallationManager()->getInstallPath($package);
        $rel = $path === null ? 'null' : str_replace(getcwd() . '/', '', $path);
        $this->io->write(sprintf('[plugfix/events] %s: %s (%s) path=%s exists=%s ops=%d', $event->getName(), $op->getOperationType(), (string) $op, $rel, var_export($path !== null && is_dir($path), true), count($event->getOperations())));
    }

    public function onCustom(Event $event)
    {
        $this->io->write('[plugfix/events] custom ' . json_encode($event->getArguments()));
    }
}
