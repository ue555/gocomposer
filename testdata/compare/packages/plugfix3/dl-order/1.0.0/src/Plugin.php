<?php
namespace Plugfix3\DlOrder;

use Composer\Composer;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Plugin\PluginInterface;
use Composer\Plugin\PostFileDownloadEvent;
use Composer\Plugin\PreFileDownloadEvent;

/**
 * Logs the package download events: with parallel transfers Composer
 * dispatches post-file-download as the transfers complete.
 */
class Plugin implements PluginInterface, EventSubscriberInterface
{
    private $io;

    public function activate(Composer $composer, IOInterface $io)
    {
        $this->io = $io;
    }

    public function deactivate(Composer $composer, IOInterface $io)
    {
    }

    public function uninstall(Composer $composer, IOInterface $io)
    {
    }

    public static function getSubscribedEvents()
    {
        return array(
            'pre-file-download' => 'onPre',
            'post-file-download' => 'onPost',
        );
    }

    public function onPre(PreFileDownloadEvent $event)
    {
        $context = $event->getContext();
        if ($context instanceof PackageInterface) {
            $this->io->write('[order] pre ' . $context->getName());
        }
    }

    public function onPost(PostFileDownloadEvent $event)
    {
        $context = $event->getContext();
        if ($context instanceof PackageInterface) {
            $this->io->write('[order] post ' . $context->getName() . ' exists=' . var_export(is_file((string) $event->getFileName()), true));
            if (getenv('PLUGFIX3_ORDER_THROW') === $context->getName()) {
                throw new \RuntimeException('plugfix3 rejects ' . $context->getName());
            }
        }
    }
}
