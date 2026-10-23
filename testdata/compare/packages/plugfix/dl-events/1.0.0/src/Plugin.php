<?php
namespace Plugfix\DlEvents;

use Composer\Composer;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Plugin\PluginInterface;
use Composer\Plugin\PostFileDownloadEvent;
use Composer\Plugin\PreFileDownloadEvent;
use Composer\Plugin\PrePoolCreateEvent;

/**
 * Logs the download and pool events and, depending on PLUGFIX_DL (a comma
 * separated list of modes), changes them:
 *   rewrite     package URL of acme/log gets a query string, acme/util a
 *               custom cache key
 *   meta        the root packages.json request is redirected to /v2 and
 *               gets an extra header
 *   pool        acme/log versions >= 1.1 are removed from the pool
 *   unacceptable  the unacceptable fixed packages are listed
 *   throw-pre   pre-file-download of packages throws
 *   throw-pool  pre-pool-create throws
 */
class Plugin implements PluginInterface, EventSubscriberInterface
{
    private $io;
    private $modes = array();

    public function activate(Composer $composer, IOInterface $io)
    {
        $this->io = $io;
        $this->modes = array_filter(explode(',', (string) getenv('PLUGFIX_DL')));
        $io->write('[dl] activate modes=' . implode(',', $this->modes));
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
            'pre-file-download' => array(array('onPre', 5), array('onPreLate', -5)),
            'post-file-download' => 'onPost',
            'pre-pool-create' => 'onPool',
        );
    }

    private function has($mode)
    {
        return in_array($mode, $this->modes, true);
    }

    private function ctx($context)
    {
        if ($context instanceof PackageInterface) {
            return 'package ' . $context->getPrettyName() . ' ' . $context->getPrettyVersion();
        }
        if (is_array($context) && isset($context['repository'])) {
            return 'repository ' . get_class($context['repository']) . ' ' . $context['repository']->getRepoName();
        }

        return gettype($context);
    }

    public function onPre(PreFileDownloadEvent $event)
    {
        $this->io->write('[dl] pre ' . $event->getType() . ' ' . $event->getProcessedUrl() . ' (' . $this->ctx($event->getContext()) . ') downloader=' . get_class($event->getHttpDownloader()));
        if ($event->getType() === 'metadata') {
            $this->io->write('[dl]   options ' . json_encode($event->getTransportOptions()));
        }
        $context = $event->getContext();
        if ($this->has('throw-pre') && $event->getType() === 'package') {
            throw new \RuntimeException('plugfix refuses to download ' . $context->getName());
        }
        if ($this->has('rewrite') && $context instanceof PackageInterface) {
            if ($context->getName() === 'acme/log') {
                $event->setProcessedUrl($event->getProcessedUrl() . '?via=plugfix');
            }
            if ($context->getName() === 'acme/util') {
                $event->setCustomCacheKey('plugfix-key-' . $context->getName());
            }
        }
        if ($this->has('meta') && $event->getType() === 'metadata' && substr($event->getProcessedUrl(), -14) === '/packages.json' && strpos($event->getProcessedUrl(), '/v2/') === false) {
            $event->setProcessedUrl(str_replace('/packages.json', '/v2/packages.json', $event->getProcessedUrl()));
            $event->setTransportOptions(array('http' => array('header' => array('X-Plugfix: 1'))));
        }
    }

    public function onPreLate(PreFileDownloadEvent $event)
    {
        $this->io->write('[dl] pre-late ' . $event->getProcessedUrl() . ' key=' . var_export($event->getCustomCacheKey(), true));
    }

    public function onPost(PostFileDownloadEvent $event)
    {
        $line = '[dl] post ' . $event->getType() . ' ' . $event->getUrl() . ' (' . $this->ctx($event->getContext()) . ')';
        if ($event->getType() === 'package') {
            $line .= ' file=' . $event->getFileName() . ' exists=' . var_export(is_file((string) $event->getFileName()), true) . ' checksum=' . var_export($event->getChecksum(), true);
        } else {
            $context = $event->getContext();
            $response = $context['response'];
            $data = $response->decodeJson();
            $line .= ' status=' . $response->getStatusCode() . ' packages=' . (isset($data['packages']) ? count($data['packages']) : '-') . ' checksum=' . var_export($event->getChecksum(), true);
        }
        $this->io->write($line);
    }

    public function onPool(PrePoolCreateEvent $event)
    {
        if ($this->has('throw-pool')) {
            throw new \RuntimeException('plugfix refuses to create the pool');
        }
        $names = array();
        foreach ($event->getPackages() as $package) {
            $names[] = $package->getPrettyName() . ' ' . $package->getPrettyVersion() . ($package instanceof \Composer\Package\AliasPackage ? ' (alias)' : '');
        }
        sort($names);
        $this->io->write('[dl] pool ' . count($names) . ' packages: ' . implode(', ', $names));
        $repos = array();
        foreach ($event->getRepositories() as $repo) {
            $repos[] = get_class($repo) . '=' . $repo->getRepoName();
        }
        $this->io->write('[dl] pool repositories: ' . implode(' | ', $repos));
        $request = $event->getRequest();
        $fixed = array();
        foreach ($request->getFixedOrLockedPackages() as $package) {
            $fixed[] = $package->getPrettyName() . ' ' . $package->getPrettyVersion();
        }
        sort($fixed);
        $this->io->write('[dl] pool fixed/locked: ' . implode(', ', $fixed));
        $requires = array();
        foreach ($request->getRequires() as $name => $constraint) {
            $requires[] = $name . ' ' . $constraint->getPrettyString();
        }
        $this->io->write('[dl] pool requires: ' . implode(', ', $requires));
        $this->io->write('[dl] pool stabilities: ' . json_encode($event->getAcceptableStabilities()) . ' flags: ' . json_encode($event->getStabilityFlags()) . ' aliases: ' . json_encode($event->getRootAliases()) . ' refs: ' . json_encode($event->getRootReferences()));
        if ($this->has('unacceptable')) {
            $un = array();
            foreach ($event->getUnacceptableFixedPackages() as $package) {
                $un[] = $package->getPrettyName() . ' ' . $package->getPrettyVersion();
            }
            $this->io->write('[dl] pool unacceptable: ' . implode(', ', $un));
        }
        if ($this->has('pool')) {
            $kept = array();
            foreach ($event->getPackages() as $package) {
                if ($package->getName() === 'acme/log' && version_compare($package->getVersion(), '1.1.0.0', '>=')) {
                    continue;
                }
                $kept[] = $package;
            }
            $event->setPackages($kept);
            $this->io->write('[dl] pool kept ' . count($kept) . ' packages');
        }
    }
}
