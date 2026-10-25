<?php
// gocomposer plugin host: the Composer services plugins reach beyond the
// basic API (see host.php for the protocol). Everything here is a thin layer
// over gocomposer's Go implementation: HTTP transfers, repositories, package
// downloads/extraction, the Installer and the autoload generator run in Go
// and are called back through Host::cb(). The event classes, DownloadManager
// and Request are Composer's own (copied unchanged from Composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT license, (c) Nils Adermann,
// Jordi Boggiano).

namespace GocomposerPluginHost {

use Composer\Composer;
use Composer\Config;
use Composer\DependencyResolver\Request;
use Composer\Downloader\DownloaderInterface;
use Composer\Downloader\DownloadManager;
use Composer\EventDispatcher\Event;
use Composer\IO\IOInterface;
use Composer\Package\PackageInterface;
use Composer\Package\Version\VersionParser;
use Composer\Plugin\PostFileDownloadEvent;
use Composer\Plugin\PreFileDownloadEvent;
use Composer\Plugin\PrePoolCreateEvent;
use Composer\Repository\ArrayRepository;
use Composer\Repository\ConfigurableRepositoryInterface;
use Composer\Repository\LockArrayRepository;
use Composer\Repository\PlatformRepository;
use Composer\Repository\RepositoryInterface;
use Composer\Repository\RootPackageRepository;
use Composer\Semver\Constraint\ConstraintInterface;
use Composer\Semver\Constraint\MatchAllConstraint;
use Composer\Util\Http\Response;
use Composer\Util\HttpDownloader;
use React\Promise\PromiseInterface;
use Symfony\Component\Console\Input\ArgvInput;

final class Api
{
    /** @var \Composer\Console\Application|null */
    public static $app = null;
    /** @var ArgvInput|null */
    public static $appInput = null;
    /** @var array<int, \Composer\Installer> */
    private static $installers = [];
    /** @var \Composer\Command\GlobalCommand|null */
    private static $global = null;
    /** @var array<int, RepositoryInterface> */
    private static $repos = [];

    /**
     * The command line gocomposer is running (Application::doRun's input):
     * bound to the definition of the running command when it is known.
     *
     * @param array<string, mixed>|null $ctx
     */
    public static function initApp(?array $ctx): void
    {
        if ($ctx === null) {
            return;
        }
        $app = new \Composer\Console\Application();
        $app->setAutoExit(false);
        $app->setCatchExceptions(false);
        $input = new ArgvInput(array_merge(['composer'], $ctx['argv'] ?? []));
        $input->setInteractive(Host::$input->isInteractive());
        if (isset($ctx['cmd']) && is_array($ctx['cmd'])) {
            $command = $app->gcProxy($ctx['cmd']);
            try {
                $command->mergeApplicationDefinition();
                $input->bind($command->getDefinition());
            } catch (\Throwable $e) {
                // gocomposer validated the input already
            }
        }
        self::$app = $app;
        self::$appInput = $input;
    }

    /**
     * @return array<string, mixed>
     */
    public static function appState(): array
    {
        if (self::$appInput === null) {
            return ['argv' => null];
        }
        $input = self::$appInput;
        $tokens = \Closure::bind(function () {
            return $this->tokens;
        }, $input, ArgvInput::class)();

        return ['argv' => array_values($tokens)];
    }

    /**
     * Commands plugins added to the application (Application::add()).
     *
     * @return list<\Symfony\Component\Console\Command\Command>
     */
    public static function appCommands(): array
    {
        return self::$app !== null ? self::$app->gcAdded() : [];
    }

    public static function globalCommand(): \Composer\Command\GlobalCommand
    {
        return self::$global ??= new \Composer\Command\GlobalCommand();
    }

    /**
     * The Composer\Installer object standing for a gocomposer Installer run.
     */
    public static function installerObj(int $id): \Composer\Installer
    {
        if (!isset(self::$installers[$id])) {
            self::$installers[$id] = \Composer\Installer::gcFor($id, []);
        }

        return self::$installers[$id];
    }

    /**
     * @param array<string, mixed> $opts
     */
    public static function installerStub(int $id, array $opts): void
    {
        self::$installers[$id] = \Composer\Installer::gcFor($id, $opts);
    }

    public static function composer(int $inst): Composer
    {
        if (!isset(Host::$instances[$inst])) {
            throw new \LogicException('Unknown gocomposer Composer instance '.$inst);
        }

        return Host::$instances[$inst];
    }

    /**
     * A Composer instance created by gocomposer for Factory::create().
     *
     * @param array<string, mixed> $desc
     */
    public static function nestedComposer(array $desc): Composer
    {
        $inst = (int) $desc['inst'];
        if (isset(Host::$instances[$inst])) {
            return Host::$instances[$inst];
        }
        $composer = new Composer();
        $composer->gcInst = $inst;
        $config = new Config();
        $config->gcInst = $inst;
        $config->gcDetached = false;
        $composer->setConfig($config);
        $composer->setPackage(Host::loadRoot($desc['root']));
        $rm = new \Composer\Repository\RepositoryManager(null, $config);
        $local = new \Composer\Repository\InstalledFilesystemRepository();
        $local->gcInst = $inst;
        $rm->setLocalRepository($local);
        $composer->setRepositoryManager($rm);
        $im = new \Composer\Installer\InstallationManager();
        $im->gcInst = $inst;
        $composer->setInstallationManager($im);
        $dispatcher = new \GocomposerPluginHost\NestedEventDispatcher($composer, Host::$io);
        $dispatcher->setRunScripts(empty($desc['disableScripts']));
        $composer->setEventDispatcher($dispatcher);
        $locker = new \Composer\Package\Locker();
        $locker->gcInst = $inst;
        $composer->setLocker($locker);
        $composer->setPluginManager(new \Composer\Plugin\PluginManager());
        if (!empty($desc['global'])) {
            $composer->setGlobal();
        }
        $im->setDefaultInstallers(Host::$io, $composer);

        return Host::$instances[$inst] = $composer;
    }

    /**
     * The event objects of the download and pool events dispatched by
     * gocomposer.
     *
     * @param array<string, mixed> $spec
     */
    public static function event(array $spec): Event
    {
        switch ($spec['kind']) {
            case 'prepool':
                $repos = [];
                foreach ($spec['repos'] as $desc) {
                    $repos[] = self::repo($desc);
                }

                return new PrePoolCreateEvent(
                    $spec['name'],
                    $repos,
                    self::request($spec['request']),
                    $spec['acceptable'] ?? [],
                    $spec['flags'] ?? [],
                    $spec['rootAliases'] ?? [],
                    $spec['rootRefs'] ?? [],
                    self::pkgs($spec['packages'] ?? []),
                    self::pkgs($spec['unacceptable'] ?? [])
                );
            case 'prefiledl':
                $context = $spec['type'] === 'package' ? Host::pkg($spec['context']) : ['repository' => self::repo($spec['repo'])];
                $event = new PreFileDownloadEvent($spec['name'], self::httpDownloader(0), $spec['url'], $spec['type'], $context);
                if (isset($spec['options'])) {
                    $event->setTransportOptions($spec['options']);
                }

                return $event;
            case 'postfiledl':
                if ($spec['type'] === 'package') {
                    $context = Host::pkg($spec['context']);
                } else {
                    $response = new Response(['url' => $spec['url']], (int) $spec['status'], $spec['headers'] ?? [], (string) base64_decode((string) ($spec['body'] ?? '')));
                    $context = ['response' => $response, 'repository' => self::repo($spec['repo'])];
                }

                return new PostFileDownloadEvent($spec['name'], $spec['file'] ?? null, $spec['checksum'] ?? null, $spec['url'], $spec['type'], $context);
        }

        throw new \LogicException('Unknown event kind '.$spec['kind']);
    }

    /**
     * What gocomposer needs back from an event once its listeners ran.
     *
     * @return array<string, mixed>|null
     */
    public static function eventState(string $eid): ?array
    {
        $event = Host::$events[$eid] ?? null;
        if ($event instanceof PreFileDownloadEvent) {
            return ['url' => $event->getProcessedUrl(), 'cacheKey' => $event->getCustomCacheKey(), 'options' => $event->getTransportOptions()];
        }
        if ($event instanceof PrePoolCreateEvent) {
            $refs = static function (array $packages): array {
                $out = [];
                foreach ($packages as $p) {
                    $out[] = Host::ref($p);
                }

                return $out;
            };

            return ['packages' => $refs($event->getPackages()), 'unacceptable' => $refs($event->getUnacceptableFixedPackages())];
        }

        return null;
    }

    /**
     * @param list<array<string, mixed>> $refs
     * @return list<PackageInterface>
     */
    public static function pkgs(array $refs): array
    {
        $out = [];
        foreach ($refs as $ref) {
            $out[] = Host::pkg($ref);
        }

        return $out;
    }

    /**
     * @param array<string, mixed> $spec
     */
    private static function request(array $spec): Request
    {
        $request = new Request();
        foreach ($spec['fixed'] ?? [] as $ref) {
            $request->fixPackage(Host::pkg($ref));
        }
        foreach ($spec['locked'] ?? [] as $ref) {
            $request->lockPackage(Host::pkg($ref));
        }
        foreach ($spec['requires'] ?? [] as $req) {
            $request->requireName($req['name'], self::constraint($req));
        }
        if (isset($spec['allow'])) {
            $request->setUpdateAllowList($spec['allow'], (int) ($spec['transitive'] ?? 0));
        }

        return $request;
    }

    /**
     * @param array<string, mixed> $req
     */
    public static function constraint(array $req): ConstraintInterface
    {
        $parser = new VersionParser();
        foreach (['pretty', 'c'] as $k) {
            if (!isset($req[$k]) || $req[$k] === '') {
                continue;
            }
            try {
                $constraint = $parser->parseConstraints($req[$k]);
                if (isset($req['pretty'])) {
                    $constraint->setPrettyString($req['pretty']);
                }

                return $constraint;
            } catch (\Throwable $e) {
            }
        }

        return new MatchAllConstraint();
    }

    /**
     * A constraint given by plugin code, as the string gocomposer parses.
     *
     * @param mixed $constraint
     */
    public static function constraintString($constraint): ?string
    {
        if ($constraint === null) {
            return null;
        }
        if ($constraint instanceof ConstraintInterface) {
            $pretty = $constraint->getPrettyString();
            if ($constraint instanceof MatchAllConstraint) {
                return '*';
            }

            return $pretty !== '' ? $pretty : (string) $constraint;
        }

        return (string) $constraint;
    }

    /**
     * The repository object for a gocomposer repository descriptor.
     *
     * @param array<string, mixed> $desc
     */
    public static function repo(array $desc): RepositoryInterface
    {
        $rid = (int) $desc['rid'];
        if (isset(self::$repos[$rid])) {
            return self::$repos[$rid];
        }
        switch ($desc['kind']) {
            case 'composer':
                $repo = \Composer\Repository\ComposerRepository::gcCreate($desc);
                break;
            case 'path':
                $repo = \Composer\Repository\PathRepository::gcCreate($desc);
                break;
            case 'vcs':
                $repo = \Composer\Repository\VcsRepository::gcCreate($desc);
                break;
            case 'artifact':
                $repo = \Composer\Repository\ArtifactRepository::gcCreate($desc);
                break;
            case 'package':
                $repo = \Composer\Repository\PackageRepository::gcCreate($desc);
                break;
            case 'installed':
                $repo = self::composer((int) ($desc['inst'] ?? 0))->getRepositoryManager()->getLocalRepository();
                break;
            case 'php':
                $repo = Services::repoObject((int) $desc['id']);
                break;
            case 'root':
                $repo = new RootPackageRepository(self::composer((int) ($desc['inst'] ?? 0))->getPackage());
                break;
            case 'platform':
                $repo = (new \ReflectionClass(PlatformRepository::class))->newInstanceWithoutConstructor();
                \Closure::bind(function () {
                    $this->packages = [];
                }, $repo, ArrayRepository::class)();
                foreach (self::pkgs($desc['pkgs'] ?? []) as $p) {
                    $repo->addPackage($p);
                }
                break;
            case 'lock':
                $repo = new LockArrayRepository();
                foreach (self::pkgs($desc['pkgs'] ?? []) as $p) {
                    $repo->addPackage($p);
                }
                break;
            default:
                $repo = new ArrayRepository();
                foreach (self::pkgs($desc['pkgs'] ?? []) as $p) {
                    $repo->addPackage($p);
                }
        }

        return self::$repos[$rid] = $repo;
    }

    public static function adoptRepo(GoRepository $repo): void
    {
        self::$repos[$repo->gcRid] = $repo;
    }

    /**
     * The gocomposer id of a repository object given by plugin code.
     */
    public static function repoId(RepositoryInterface $repo): ?int
    {
        foreach (self::$repos as $rid => $r) {
            if ($r === $repo) {
                return $rid;
            }
        }

        return null;
    }

    public static function httpDownloader(int $inst): HttpDownloader
    {
        return self::composer($inst)->getLoop()->getHttpDownloader();
    }

    /**
     * Factory::createDownloadManager with gocomposer's downloaders.
     */
    public static function downloadManager(Composer $composer): DownloadManager
    {
        $io = Host::$io;
        $config = $composer->getConfig();
        $dm = new GoDownloadManager($io, false);
        $dm->gcInst = $composer->gcInst;
        $preferred = $config->get('preferred-install');
        switch ($preferred) {
            case 'dist':
                $dm->setPreferDist(true);
                break;
            case 'source':
                $dm->setPreferSource(true);
                break;
        }
        if (is_array($preferred)) {
            $dm->setPreferences($preferred);
        }
        $http = $composer->getLoop()->getHttpDownloader();
        $dm->setDownloader('git', new \Composer\Downloader\GitDownloader($io, $config));
        $dm->setDownloader('svn', new \Composer\Downloader\SvnDownloader($io, $config));
        $dm->setDownloader('fossil', new \Composer\Downloader\FossilDownloader($io, $config));
        $dm->setDownloader('hg', new \Composer\Downloader\HgDownloader($io, $config));
        $dm->setDownloader('perforce', new \Composer\Downloader\PerforceDownloader($io, $config));
        $dm->setDownloader('zip', new \Composer\Downloader\ZipDownloader($io, $config, $http));
        $dm->setDownloader('rar', new \Composer\Downloader\RarDownloader($io, $config, $http));
        $dm->setDownloader('tar', new \Composer\Downloader\TarDownloader($io, $config, $http));
        $dm->setDownloader('gzip', new \Composer\Downloader\GzipDownloader($io, $config, $http));
        $dm->setDownloader('xz', new \Composer\Downloader\XzDownloader($io, $config, $http));
        $dm->setDownloader('phar', new \Composer\Downloader\PharDownloader($io, $config, $http));
        $dm->setDownloader('file', new \Composer\Downloader\FileDownloader($io, $config, $http));
        $dm->setDownloader('path', new \Composer\Downloader\PathDownloader($io, $config, $http));
        $dm->gcReady = true;

        return $dm;
    }

    /**
     * The gocomposer implementation a downloader object stands for (its
     * closest Composer downloader class).
     */
    public static function downloaderKind(DownloaderInterface $downloader): string
    {
        $map = [
            \Composer\Downloader\ZipDownloader::class => 'zip',
            \Composer\Downloader\RarDownloader::class => 'rar',
            \Composer\Downloader\TarDownloader::class => 'tar',
            \Composer\Downloader\GzipDownloader::class => 'gzip',
            \Composer\Downloader\XzDownloader::class => 'xz',
            \Composer\Downloader\PharDownloader::class => 'phar',
            \Composer\Downloader\PathDownloader::class => 'path',
            \Composer\Downloader\FileDownloader::class => 'file',
            \Composer\Downloader\GitDownloader::class => 'git',
            \Composer\Downloader\SvnDownloader::class => 'svn',
            \Composer\Downloader\HgDownloader::class => 'hg',
            \Composer\Downloader\FossilDownloader::class => 'fossil',
            \Composer\Downloader\PerforceDownloader::class => 'perforce',
        ];
        foreach ($map as $class => $kind) {
            if ($downloader instanceof $class) {
                return $kind;
            }
        }

        throw new \RuntimeException('The downloader '.get_class($downloader).' does not extend one of Composer\'s downloaders, gocomposer can not run it. Run the command with --no-plugins to skip plugins.');
    }

    /**
     * A downloader operation performed by gocomposer (see GoDownloader).
     *
     * @param array<string, mixed> $args
     */
    public static function dl(DownloaderInterface $downloader, string $m, array $args): PromiseInterface
    {
        $args['kind'] = self::downloaderKind($downloader);
        $args['m'] = $m;
        $r = Host::cb('dl', $args);
        if (is_array($r) && !empty($r['async'])) {
            // Composer runs this step in a background process (unzip, rm -rf)
            // whose promise settles when the loop waits; gocomposer did the
            // work already, the promise settles at the same point
            return self::defer();
        }

        return \React\Promise\resolve(null);
    }

    /** @var array<int, array{d: \React\Promise\Deferred, started: bool}> background steps by job id */
    private static $jobs = [];
    /** @var int */
    private static $jobSeq = 0;
    /** @var int */
    private static $runningJobs = 0;

    /**
     * A promise standing for a step Composer runs in a background process
     * (ProcessExecutor::executeAsync): gocomposer did the work already, the
     * promise settles when the loop waits, in the order Composer's
     * ProcessExecutor completes its jobs (at most
     * COMPOSER_MAX_PARALLEL_PROCESSES started at a time, a job created while
     * a slot is free starts at once, queued ones start on the next tick).
     */
    public static function defer(): PromiseInterface
    {
        $deferred = new \React\Promise\Deferred();
        $id = self::$jobSeq++;
        self::$jobs[$id] = ['d' => $deferred, 'started' => false];
        if (self::$runningJobs < self::maxJobs()) {
            self::$jobs[$id]['started'] = true;
            self::$runningJobs++;
        }

        return $deferred->promise();
    }

    private static function maxJobs(): int
    {
        $maxJobs = \Composer\Util\Platform::getEnv('COMPOSER_MAX_PARALLEL_PROCESSES');
        if (is_numeric($maxJobs)) {
            return max(1, min(50, (int) $maxJobs));
        }

        return 10;
    }

    /**
     * Settles the background steps (also those queued meanwhile) like
     * ProcessExecutor::wait (countActiveJobs ticks; every started job is
     * complete by the next tick); returns whether there were any.
     */
    public static function drainDeferred(): bool
    {
        $any = false;
        while (self::$jobs !== []) {
            $any = true;
            foreach (array_keys(self::$jobs) as $id) {
                if (!isset(self::$jobs[$id])) {
                    continue;
                }
                if (self::$jobs[$id]['started']) {
                    $deferred = self::$jobs[$id]['d'];
                    unset(self::$jobs[$id]);
                    self::$runningJobs--;
                    $deferred->resolve(null);
                    continue;
                }
                if (self::$runningJobs < self::maxJobs()) {
                    self::$jobs[$id]['started'] = true;
                    self::$runningJobs++;
                }
            }
        }

        return $any;
    }

    /**
     * The downloader a plugin registered for a package type gocomposer does
     * not know, or registered in place of a built-in one
     * (DownloadManager::setDownloader()).
     *
     * @return array<string, mixed>|null
     */
    public static function downloaderFor(string $type): ?array
    {
        $dm = Host::$composer->gcExistingDownloadManager();
        if (!$dm instanceof GoDownloadManager || (in_array(strtolower($type), GoDownloadManager::BUILTIN, true) && !isset($dm->gcReplaced[strtolower($type)]))) {
            return null;
        }
        try {
            $downloader = $dm->getDownloader($type);
        } catch (\InvalidArgumentException $e) {
            return null;
        }

        return ['class' => get_class($downloader), 'source' => $downloader->getInstallationSource()];
    }

    /**
     * Runs a downloader registered by a plugin for a package type gocomposer
     * does not know.
     *
     * @param array<string, mixed> $a
     * @return mixed
     */
    public static function downloaderOp(array $a)
    {
        $dm = Host::$composer->gcExistingDownloadManager() ?? Host::$composer->getDownloadManager();
        $downloader = $dm->getDownloader($a['type']);
        $p = Host::pkg($a['p']);
        $i = isset($a['i']) ? Host::pkg($a['i']) : null;
        $path = $a['path'];
        switch ($a['m']) {
            case 'download':
                $promise = $downloader->download($p, $path, $i);
                break;
            case 'prepare':
            case 'cleanup':
                $method = $a['m'];
                $promise = $downloader->$method($a['op'], $p, $path, $i);
                break;
            case 'install':
                $promise = $downloader->install($p, $path);
                break;
            case 'update':
                $promise = $downloader->update($i, $p, $path);
                break;
            case 'remove':
                $promise = $downloader->remove($p, $path);
                break;
            default:
                throw new \LogicException('Unknown downloader operation '.$a['m']);
        }
        self::settle($promise);

        return null;
    }

    /**
     * Runs a promise chain to completion (the HTTP jobs it waits on are run
     * by the loop).
     *
     * @param mixed $promise
     */
    public static function settle($promise): void
    {
        if (!$promise instanceof PromiseInterface) {
            return;
        }
        $error = null;
        $done = false;
        $promise->then(static function () use (&$done): void {
            $done = true;
        }, static function ($e) use (&$error, &$done): void {
            $error = $e;
            $done = true;
        });
        if (!$done) {
            Host::$composer->getLoop()->wait([$promise]);
        }
        if (!$done) {
            throw new \LogicException('A promise of plugin code did not settle');
        }
        if ($error instanceof \Throwable) {
            throw $error;
        }
    }
}

/**
 * The DownloadManager of a Composer instance: gocomposer's downloaders handle
 * the package types it knows, plugins may add downloaders for other types or
 * replace the built-in ones; gocomposer then calls the plugin's downloader
 * for the operations of that type (see plugin_services.go).
 */
class GoDownloadManager extends DownloadManager
{
    public const BUILTIN = ['git', 'svn', 'fossil', 'hg', 'perforce', 'zip', 'rar', 'tar', 'gzip', 'xz', 'phar', 'file', 'path'];

    /** @var bool */
    public $gcReady = false;
    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var array<string, true> built-in types replaced by plugins */
    public $gcReplaced = [];

    public function setDownloader(string $type, DownloaderInterface $downloader): DownloadManager
    {
        parent::setDownloader($type, $downloader);
        $type = strtolower($type);
        if ($this->gcReady && in_array($type, self::BUILTIN, true) && !isset($this->gcReplaced[$type])) {
            $this->gcReplaced[$type] = true;
            if ($this === Api::composer($this->gcInst)->gcExistingDownloadManager() && $this->gcInst === 0) {
                Host::cb('dlReplaced', ['inst' => $this->gcInst, 'type' => $type]);
            }
        }

        return $this;
    }
}

/**
 * A gocomposer command as seen by plugins through Application::find(): name,
 * aliases and input definition only (it is run by gocomposer).
 */
class ProxyCommand extends \Composer\Command\BaseCommand
{
    /** @var array<string, mixed> */
    private $gcInfo;

    /**
     * @param array<string, mixed> $info
     */
    public function __construct(array $info)
    {
        $this->gcInfo = $info;
        parent::__construct($info['name']);
    }

    protected function configure(): void
    {
        $info = $this->gcInfo;
        $this->setAliases($info['aliases'] ?? []);
        $this->setDescription($info['description'] ?? '');
        foreach ($info['args'] ?? [] as $arg) {
            $mode = $arg['required'] ? \Symfony\Component\Console\Input\InputArgument::REQUIRED : \Symfony\Component\Console\Input\InputArgument::OPTIONAL;
            if ($arg['array']) {
                $mode |= \Symfony\Component\Console\Input\InputArgument::IS_ARRAY;
            }
            $this->addArgument($arg['name'], $mode, $arg['desc'] ?? '', $arg['required'] ? null : ($arg['default'] ?? ($arg['array'] ? [] : null)));
        }
        foreach ($info['opts'] ?? [] as $opt) {
            if (in_array($opt['name'], ['help', 'quiet', 'verbose', 'version', 'ansi', 'no-ansi', 'no-interaction', 'profile', 'no-plugins', 'no-scripts', 'working-dir', 'no-cache'], true)) {
                continue;
            }
            $mode = \Symfony\Component\Console\Input\InputOption::VALUE_NONE;
            if ($opt['array']) {
                $mode = \Symfony\Component\Console\Input\InputOption::VALUE_REQUIRED | \Symfony\Component\Console\Input\InputOption::VALUE_IS_ARRAY;
            } elseif ($opt['required']) {
                $mode = \Symfony\Component\Console\Input\InputOption::VALUE_REQUIRED;
            } elseif ($opt['optional']) {
                $mode = \Symfony\Component\Console\Input\InputOption::VALUE_OPTIONAL;
            } elseif (!empty($opt['negatable'])) {
                $mode |= \Symfony\Component\Console\Input\InputOption::VALUE_NEGATABLE;
            }
            $default = null;
            if ($opt['array']) {
                $default = [];
            } elseif ($mode !== \Symfony\Component\Console\Input\InputOption::VALUE_NONE && isset($opt['default']) && $opt['default'] !== '') {
                $default = $opt['default'];
            }
            $this->addOption($opt['name'], ($opt['shortcut'] ?? '') !== '' ? $opt['shortcut'] : null, $mode, $opt['desc'] ?? '', $default);
        }
    }

    protected function execute(\Symfony\Component\Console\Input\InputInterface $input, \Symfony\Component\Console\Output\OutputInterface $output): int
    {
        Unsupported::fail('Running gocomposer\'s '.$this->getName().' command from plugin code');
    }
}

/**
 * The event dispatcher of a Composer instance created with Factory::create():
 * gocomposer dispatches the events with that instance (its root scripts and
 * the listeners of its plugins, which run in their own gocomposer plugin
 * host). Listeners can not be added from here (they would have to join the
 * listeners of the other host).
 */
class NestedEventDispatcher extends \Composer\EventDispatcher\EventDispatcher
{
    public function dispatch(?string $eventName, ?Event $event = null): int
    {
        if (null === $event) {
            if (null === $eventName) {
                throw new \InvalidArgumentException('If no $event is passed in to '.__METHOD__.' you have to pass in an $eventName, got null.');
            }
            $event = new Event($eventName);
        }

        return Host::dispatchFromPlugin($event, $this->composer->gcInst, $this->runScripts);
    }

    public function dispatchScript(string $eventName, bool $devMode = false, array $additionalArgs = [], array $flags = []): int
    {
        assert($this->composer instanceof Composer, new \LogicException('This should only be reached with a fully loaded Composer'));

        return Host::dispatchFromPlugin(new \Composer\Script\Event($eventName, $this->composer, $this->io, $devMode, $additionalArgs, $flags), $this->composer->gcInst, $this->runScripts);
    }

    public function hasEventListeners(Event $event): bool
    {
        return (bool) Host::cb('nestedHasListeners', ['inst' => $this->composer->gcInst, 'name' => $event->getName(), 'runScripts' => $this->runScripts]);
    }

    public function addListener(string $eventName, $listener, int $priority = 0): void
    {
        Unsupported::fail(__METHOD__.' on a Composer instance created by a plugin');
    }

    public function setRunScripts(bool $runScripts = true): self
    {
        $this->runScripts = $runScripts;

        return $this;
    }
}

/**
 * A repository implemented by gocomposer (composer, vcs, path, artifact and
 * package repositories): queries are answered by the Go repository.
 */
class GoRepository extends ArrayRepository
{
    /** @var int */
    public $gcRid = 0;
    /** @var string */
    protected $gcName = '';
    /** @var array<string, mixed> */
    protected $gcConfig = [];

    protected const GC_TYPE = '';

    /**
     * new ComposerRepository($repoConfig, ...) etc.: the repository is
     * created by gocomposer from the configuration.
     */
    public function __construct(...$args)
    {
        $config = $args[0] ?? null;
        if (static::GC_TYPE === '' || !is_array($config)) {
            Unsupported::fail(static::class.'::__construct');
        }
        $desc = Host::cb('repoCreate', ['inst' => 0, 'type' => static::GC_TYPE, 'config' => $config, 'name' => null]);
        $this->gcRid = (int) $desc['rid'];
        $this->gcName = (string) $desc['name'];
        $this->gcConfig = $desc['config'] ?? [];
        Api::adoptRepo($this);
    }

    /**
     * @param array<string, mixed> $desc
     * @return static
     */
    public static function gcCreate(array $desc)
    {
        $repo = (new \ReflectionClass(static::class))->newInstanceWithoutConstructor();
        $repo->gcRid = (int) $desc['rid'];
        $repo->gcName = (string) $desc['name'];
        $repo->gcConfig = $desc['config'] ?? [];

        return $repo;
    }

    /**
     * @param array<string, mixed> $a
     * @return mixed
     */
    protected function gcOp(string $m, array $a = [])
    {
        return Host::cb('repoOp', ['rid' => $this->gcRid, 'm' => $m] + $a);
    }

    /**
     * @param list<array<string, mixed>> $refs
     * @return list<PackageInterface>
     */
    protected function gcPackages(array $refs): array
    {
        $out = [];
        foreach ($refs as $ref) {
            $p = Host::pkg($ref);
            if ($p->getRepository() === null) {
                $p->setRepository($this);
            }
            $out[] = $p;
        }

        return $out;
    }

    public function getRepoName()
    {
        return $this->gcName;
    }

    public function getRepoConfig()
    {
        return $this->gcConfig;
    }

    public function findPackage(string $name, $constraint)
    {
        if ($this->gcMaterialized) {
            return parent::findPackage($name, $constraint);
        }
        $found = $this->gcPackages($this->gcOp('find', ['name' => $name, 'c' => Api::constraintString($constraint), 'one' => true]));

        return $found[0] ?? null;
    }

    public function findPackages(string $name, $constraint = null)
    {
        if ($this->gcMaterialized) {
            return parent::findPackages($name, $constraint);
        }
        return $this->gcPackages($this->gcOp('find', ['name' => $name, 'c' => Api::constraintString($constraint), 'one' => false]));
    }

    public function getPackages()
    {
        if ($this->gcMaterialized) {
            return parent::getPackages();
        }
        return $this->gcPackages($this->gcOp('packages'));
    }

    public function loadPackages(array $packageNameMap, array $acceptableStabilities, array $stabilityFlags, array $alreadyLoaded = [])
    {
        if ($this->gcMaterialized) {
            return parent::loadPackages($packageNameMap, $acceptableStabilities, $stabilityFlags, $alreadyLoaded);
        }
        $names = [];
        foreach ($packageNameMap as $name => $constraint) {
            $names[] = ['name' => (string) $name, 'c' => Api::constraintString($constraint)];
        }
        $r = $this->gcOp('load', ['names' => $names, 'acceptable' => (object) $acceptableStabilities, 'flags' => (object) $stabilityFlags, 'loaded' => (object) $alreadyLoaded]);

        return ['namesFound' => $r['namesFound'] ?? [], 'packages' => $this->gcPackages($r['packages'] ?? [])];
    }

    public function search(string $query, int $mode = 0, ?string $type = null)
    {
        if ($this->gcMaterialized) {
            return parent::search($query, $mode, $type);
        }
        return $this->gcOp('search', ['query' => $query, 'mode' => $mode, 'type' => $type]);
    }

    public function getProviders(string $packageName)
    {
        if ($this->gcMaterialized) {
            return parent::getProviders($packageName);
        }
        return $this->gcOp('providers', ['name' => $packageName]);
    }

    public function hasPackage(PackageInterface $package)
    {
        if ($this->gcMaterialized) {
            return parent::hasPackage($package);
        }
        foreach ($this->findPackages($package->getName(), new \Composer\Semver\Constraint\Constraint('==', $package->getVersion())) as $p) {
            if ($p->getUniqueName() === $package->getUniqueName()) {
                return true;
            }
        }

        return false;
    }

    public function count(): int
    {
        if ($this->gcMaterialized) {
            return parent::count();
        }
        return count($this->getPackages());
    }

    /** @var bool the object keeps the packages: a plugin added or removed some */
    protected $gcMaterialized = false;

    public function gcIsMaterialized(): bool
    {
        return $this->gcMaterialized;
    }

    /**
     * The repository becomes an ArrayRepository holding its packages, which
     * gocomposer asks from then on (the repositories built on Composer's
     * ArrayRepository: package, path, vcs, artifact).
     */
    protected function gcMaterialize(string $method): void
    {
        if ($this->gcMaterialized) {
            return;
        }
        if (static::GC_TYPE === 'composer' || static::GC_TYPE === '') {
            Unsupported::fail(static::class.'::'.$method);
        }
        $packages = $this->getPackages();
        $this->packages = [];
        foreach ($packages as $package) {
            parent::addPackage($package);
        }
        $this->gcMaterialized = true;
        Host::cb('repoOverlay', ['rid' => $this->gcRid, 'php' => Services::registerRepo($this)]);
    }

    public function addPackage(PackageInterface $package)
    {
        $this->gcMaterialize('addPackage');
        parent::addPackage($package);
    }

    public function removePackage(PackageInterface $package): void
    {
        $this->gcMaterialize('removePackage');
        parent::removePackage($package);
    }
}

/**
 * The downloaders of the DownloadManager: each operation is performed by
 * gocomposer's implementation of the corresponding Composer downloader.
 */
trait GoDownloader
{
    public function download(PackageInterface $package, string $path, ?PackageInterface $prevPackage = null, bool $output = true): PromiseInterface
    {
        return Api::dl($this, 'download', ['p' => Host::ref($package), 'i' => $prevPackage !== null ? Host::ref($prevPackage) : null, 'path' => $path, 'src' => $package->getInstallationSource(), 'output' => $output]);
    }

    public function prepare(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        return Api::dl($this, 'prepare', ['op' => $type, 'p' => Host::ref($package), 'i' => $prevPackage !== null ? Host::ref($prevPackage) : null, 'path' => $path, 'src' => $package->getInstallationSource()]);
    }

    public function cleanup(string $type, PackageInterface $package, string $path, ?PackageInterface $prevPackage = null): PromiseInterface
    {
        return Api::dl($this, 'cleanup', ['op' => $type, 'p' => Host::ref($package), 'i' => $prevPackage !== null ? Host::ref($prevPackage) : null, 'path' => $path, 'src' => $package->getInstallationSource()]);
    }

    public function install(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        return Api::dl($this, 'install', ['p' => Host::ref($package), 'path' => $path, 'src' => $package->getInstallationSource(), 'output' => $output]);
    }

    public function update(PackageInterface $initial, PackageInterface $target, string $path): PromiseInterface
    {
        return Api::dl($this, 'update', ['p' => Host::ref($target), 'i' => Host::ref($initial), 'path' => $path, 'src' => $target->getInstallationSource(), 'isrc' => $initial->getInstallationSource()]);
    }

    public function remove(PackageInterface $package, string $path, bool $output = true): PromiseInterface
    {
        return Api::dl($this, 'remove', ['p' => Host::ref($package), 'path' => $path, 'src' => $package->getInstallationSource(), 'output' => $output]);
    }
}

}
