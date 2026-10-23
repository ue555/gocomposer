<?php
namespace Plugfix\Thrower;

use Composer\Composer;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\IO\IOInterface;
use Composer\Plugin\PluginInterface;
use Composer\Script\Event;

class Plugin implements PluginInterface, EventSubscriberInterface
{
    public function activate(Composer $composer, IOInterface $io) {}
    public function deactivate(Composer $composer, IOInterface $io) {}
    public function uninstall(Composer $composer, IOInterface $io) {}

    public static function getSubscribedEvents()
    {
        return array('post-autoload-dump' => array(array('stop', 20), array('boom', 10)), 'post-install-cmd' => 'fail');
    }

    public function stop(Event $event)
    {
        $event->getIO()->write('[plugfix/thrower] stop ' . $event->getName());
        if (getenv('PLUGFIX_STOP') || file_exists('plugfix-stop')) {
            $event->stopPropagation();
        }
    }

    public function boom(Event $event)
    {
        if (getenv('PLUGFIX_BOOM') || file_exists('plugfix-boom')) {
            throw new \RuntimeException('plugfix boom in ' . $event->getName());
        }
        $event->getIO()->write('[plugfix/thrower] no boom');
    }

    public function fail(Event $event)
    {
        return false;
    }
}
