<?php
namespace Plugfix\AutoloadPlugin;

use Composer\Composer;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\IO\IOInterface;
use Composer\Plugin\PluginInterface;
use Composer\Script\Event;
use Composer\Util\Filesystem;
use Composer\Util\Platform;

class Plugin implements PluginInterface, EventSubscriberInterface
{
    public function activate(Composer $composer, IOInterface $io) {}
    public function deactivate(Composer $composer, IOInterface $io) {}
    public function uninstall(Composer $composer, IOInterface $io) {}

    public static function getSubscribedEvents()
    {
        return array(
            'pre-autoload-dump' => 'generate',
            'post-autoload-dump' => 'dispatchCustom',
            'plugfix-custom' => array('onCustom', 5),
        );
    }

    public function generate(Event $event)
    {
        $fs = new Filesystem();
        $dir = getcwd().'/lib/generated';
        $fs->ensureDirectoryExists($dir);
        file_put_contents($dir.'/Info.php', "<?php\nnamespace Plugfix\\Generated;\n\nclass Info\n{\n    const DEV = ".var_export($event->isDevMode(), true).";\n}\n");
        $root = $event->getComposer()->getPackage();
        $autoload = $root->getAutoload();
        $autoload['classmap'][] = 'lib/generated/Info.php';
        $root->setAutoload($autoload);
        $event->getIO()->write('[plugfix/autoload] generated, dev env='.var_export(getenv('COMPOSER_DEV_MODE'), true).' rel='.$fs->findShortestPath(getcwd().'/vendor/composer', $dir, true));
    }

    public function dispatchCustom(Event $event)
    {
        Platform::putEnv('PLUGFIX_FROM_PLUGIN', 'yes');
        $code = $event->getComposer()->getEventDispatcher()->dispatchScript('plugfix-custom', $event->isDevMode(), array('x', 'y'));
        $event->getIO()->write('[plugfix/autoload] custom returned '.$code.' x-env='.var_export(getenv('PLUGFIX_X'), true));
    }

    public function onCustom(Event $event)
    {
        $event->getIO()->write('[plugfix/autoload] onCustom '.json_encode($event->getArguments()).' origin='.($event->getOriginatingEvent() ? $event->getOriginatingEvent()->getName() : 'none'));
    }
}
