<?php
// gocomposer plugin host: the host classes (see host.php). Composer's
// composerRequire() is from src/Composer/Autoload/AutoloadGenerator.php of
// Composer f83588392df43d4fd027d84d26de459946ef4e80 (MIT license, (c) Nils
// Adermann, Jordi Boggiano).

namespace Composer\Autoload {
    // AutoloadGenerator.php of Composer f83588392df43d4fd027d84d26de459946ef4e80
    // (MIT license, (c) Nils Adermann, Jordi Boggiano)
    function composerRequire(string $fileIdentifier, string $file): void
    {
        if (empty($GLOBALS['__composer_autoload_files'][$fileIdentifier])) {
            $GLOBALS['__composer_autoload_files'][$fileIdentifier] = true;

            require $file;
        }
    }
}

namespace GocomposerPluginHost {

use Composer\Composer;
use Composer\Config;
use Composer\DependencyResolver\LocalRepoTransaction;
use Composer\DependencyResolver\Operation\InstallOperation;
use Composer\DependencyResolver\Operation\MarkAliasInstalledOperation;
use Composer\DependencyResolver\Operation\MarkAliasUninstalledOperation;
use Composer\DependencyResolver\Operation\UninstallOperation;
use Composer\DependencyResolver\Operation\UpdateOperation;
use Composer\EventDispatcher\Event;
use Composer\EventDispatcher\EventDispatcher;
use Composer\EventDispatcher\EventSubscriberInterface;
use Composer\Installer\InstallationManager;
use Composer\Installer\InstallerEvent;
use Composer\Installer\InstallerInterface;
use Composer\Installer\LibraryInstaller;
use Composer\Installer\PackageEvent;
use Composer\IO\ConsoleIO;
use Composer\IO\IOInterface;
use Composer\Package\AliasPackage;
use Composer\Package\CompleteAliasPackage;
use Composer\Package\CompletePackage;
use Composer\Package\Dumper\ArrayDumper;
use Composer\Package\Loader\ArrayLoader;
use Composer\Package\Locker;
use Composer\Package\PackageInterface;
use Composer\Package\RootAliasPackage;
use Composer\Package\RootPackage;
use Composer\Package\RootPackageInterface;
use Composer\Pcre\Preg;
use Composer\Plugin\CommandEvent;
use Composer\Plugin\PluginInterface;
use Composer\Plugin\PluginManager;
use Composer\Plugin\PreCommandRunEvent;
use Composer\Repository\InstalledFilesystemRepository;
use Composer\Repository\RepositoryManager;
use Composer\Script\Event as ScriptEvent;
use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Helper\HelperSet;
use Symfony\Component\Console\Helper\QuestionHelper;
use Symfony\Component\Console\Input\ArgvInput;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputDefinition;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\ConsoleOutputInterface;
use Symfony\Component\Console\Output\Output;
use Symfony\Component\Console\Output\OutputInterface;

final class Unsupported
{
    /**
     * @return never
     */
    public static function fail(string $what): void
    {
        throw new \RuntimeException($what.'() is not supported by the gocomposer plugin runtime. Run the command with --no-plugins to skip Composer plugins.');
    }
}

/**
 * Symfony output writing to the gocomposer process (stream 1 = stdout,
 * 2 = stderr), formatted by Symfony like Composer's ConsoleOutput.
 */
class GoOutput extends Output
{
    /** @var int */
    private $stream;

    public function __construct(int $stream, int $verbosity, bool $decorated, ?\Symfony\Component\Console\Formatter\OutputFormatterInterface $formatter = null)
    {
        $this->stream = $stream;
        parent::__construct($verbosity, $decorated, $formatter);
    }

    protected function doWrite(string $message, bool $newline)
    {
        Host::cb('out', ['s' => $this->stream, 'd' => $message.($newline ? PHP_EOL : '')]);
    }
}

class GoConsoleOutput extends GoOutput implements ConsoleOutputInterface
{
    /** @var OutputInterface */
    private $stderr;
    /** @var array<int, mixed> */
    private $sections = [];

    public function __construct(int $verbosity, bool $decorated, \Symfony\Component\Console\Formatter\OutputFormatterInterface $formatter)
    {
        parent::__construct(1, $verbosity, $decorated, $formatter);
        $this->stderr = new GoOutput(2, $verbosity, $decorated, $formatter);
    }

    public function setDecorated(bool $decorated)
    {
        parent::setDecorated($decorated);
        $this->stderr->setDecorated($decorated);
    }

    public function setFormatter(\Symfony\Component\Console\Formatter\OutputFormatterInterface $formatter)
    {
        parent::setFormatter($formatter);
        $this->stderr->setFormatter($formatter);
    }

    public function setVerbosity(int $level)
    {
        parent::setVerbosity($level);
        $this->stderr->setVerbosity($level);
    }

    public function getErrorOutput()
    {
        return $this->stderr;
    }

    public function setErrorOutput(OutputInterface $error)
    {
        $this->stderr = $error;
    }

    public function section(): \Symfony\Component\Console\Output\ConsoleSectionOutput
    {
        Unsupported::fail(__METHOD__);
    }
}

final class Host
{
    /** @var resource */
    public static $in;
    /** @var resource */
    public static $out;
    /** @var Composer */
    public static $composer;
    /** @var IOInterface */
    public static $io;
    /** @var GoConsoleOutput */
    public static $output;
    /** @var ArgvInput */
    public static $input;
    /** @var array<int, PackageInterface> */
    public static $pkgs = [];
    /** @var array<int, int> spl_object_id => package id */
    public static $objIds = [];
    /** @var array<int, Event> */
    public static $events = [];
    /** @var array<string, array<int, list<array{id: int, cb: mixed}>>> */
    public static $listeners = [];
    /** @var int */
    public static $listenerSeq = 0;
    /** @var array<int, mixed> */
    public static $listenerById = [];
    /** @var array<string, list<object>> */
    public static $registered = [];
    /** @var list<PluginInterface> */
    public static $plugins = [];
    /** @var array<int, \Symfony\Component\Console\Command\Command> */
    public static $commands = [];
    /** @var int */
    public static $classCounter = 0;
    /** @var string|null */
    public static $rootSynced = null;
    /** @var \Composer\Autoload\ClassLoader|null */
    public static $projectLoader = null;
    /** @var array<string, list<\Composer\DependencyResolver\Operation\OperationInterface>> */
    public static $opsCache = [];
    /** @var array<int, InstallerInterface> */
    public static $installerIds = [];
    /** @var string */
    public static $tmpDir;
    /** @var int */
    public static $hasShownDeprecationNotice = 0;
    /** @var array<int, Composer> gocomposer Composer instances (0: the one running the plugins) */
    public static $instances = [];

    /** @var array<string, string>|null */
    public static $envSnapshot = null;

    /**
     * Environment changes made by plugin code since the last message
     * (Platform::putEnv() calls are seen by gocomposer's scripts too).
     */
    private static function envDiff(): ?array
    {
        $cur = getenv();
        if (self::$envSnapshot === null) {
            self::$envSnapshot = $cur;

            return null;
        }
        $set = [];
        foreach ($cur as $k => $v) {
            if (!array_key_exists($k, self::$envSnapshot) || self::$envSnapshot[$k] !== $v) {
                $set[$k] = $v;
            }
        }
        $unset = array_values(array_diff(array_keys(self::$envSnapshot), array_keys($cur)));
        // Composer reads its own settings with Platform::getEnv(), which
        // prefers $_SERVER: COMPOSER_* values set there by plugin code (e.g.
        // symfony/flex setting COMPOSER_PREFER_DEV_OVER_PRERELEASE) apply to
        // gocomposer too
        foreach ($_SERVER as $k => $v) {
            if (is_string($k) && is_string($v) && strncmp($k, 'COMPOSER_', 9) === 0 && ($cur[$k] ?? null) !== $v) {
                putenv($k.'='.$v);
                $cur[$k] = $v;
                $set[$k] = $v;
            }
        }
        self::$envSnapshot = $cur;
        if ($set === [] && $unset === []) {
            return null;
        }

        return ['set' => (object) $set, 'unset' => $unset];
    }

    private static function applyEnv(array $env): void
    {
        foreach ($env['set'] ?? [] as $k => $v) {
            putenv($k.'='.$v);
            $_SERVER[$k] = $_ENV[$k] = $v;
        }
        foreach ($env['unset'] ?? [] as $k) {
            putenv($k);
            unset($_SERVER[$k], $_ENV[$k]);
        }
        self::$envSnapshot = getenv();
    }

    public static function send(array $msg): void
    {
        $env = self::envDiff();
        if ($env !== null) {
            $msg['env'] = $env;
        }
        $json = json_encode($msg, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_INVALID_UTF8_SUBSTITUTE | JSON_PRESERVE_ZERO_FRACTION | JSON_PARTIAL_OUTPUT_ON_ERROR);
        fwrite(self::$out, $json."\n");
        fflush(self::$out);
    }

    public static function recv(): array
    {
        $line = fgets(self::$in);
        if ($line === false) {
            exit(0);
        }

        return json_decode($line, true);
    }

    /**
     * Calls back into gocomposer.
     *
     * @return mixed
     */
    public static function cb(string $op, array $args = [])
    {
        self::send(['t' => 'cb', 'op' => $op, 'a' => $args]);
        while (true) {
            $m = self::recv();
            if (isset($m['env'])) {
                self::applyEnv($m['env']);
            }
            switch ($m['t']) {
                case 'call':
                    self::handleCall($m);
                    break;
                case 'cbret':
                    return $m['v'] ?? null;
                case 'cbexc':
                    $e = $m['e'];
                    $class = $e['class'] ?? \RuntimeException::class;
                    if (!class_exists($class) || !is_a($class, \Throwable::class, true)) {
                        $class = \RuntimeException::class;
                    }
                    throw new $class($e['message'] ?? '', (int) ($e['code'] ?? 0));
            }
        }
    }

    public static function handleCall(array $m): void
    {
        if (isset($m['env'])) {
            self::applyEnv($m['env']);
        }
        try {
            if (isset($m['root'])) {
                self::syncRoot($m['root']);
            }
            $v = self::call($m['op'], $m['a'] ?? []);
            $ret = ['t' => 'ret', 'v' => $v];
        } catch (\Throwable $e) {
            $ret = ['t' => 'exc', 'e' => self::exception($e)];
        }
        if (self::$composer !== null) {
            $rc = self::rootChanges();
            if ($rc !== null) {
                $ret['root'] = $rc;
            }
        }
        self::send($ret);
    }

    /**
     * @return array<string, mixed>
     */
    public static function exception(\Throwable $e): array
    {
        return [
            'class' => get_class($e),
            'message' => $e->getMessage(),
            'code' => $e->getCode(),
            'file' => $e->getFile(),
            'line' => $e->getLine(),
            'trace' => $e->getTraceAsString(),
            'previous' => $e->getPrevious() !== null ? self::exception($e->getPrevious()) : null,
        ];
    }

    /**
     * @return mixed
     */
    private static function call(string $op, array $a)
    {
        switch ($op) {
            case 'init':
                return self::init($a);
            case 'register':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return self::register($a);
                });
            case 'deactivatePackage':
                return self::deactivatePackage($a['name'], false);
            case 'uninstallPackage':
                return self::deactivatePackage($a['name'], true);
            case 'listeners':
                return self::listListeners($a['name']);
            case 'callListener':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return self::callListener($a);
                });
            case 'commands':
                return self::commands();
            case 'runCommand':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return self::runCommand($a);
                });
            case 'installerFor':
                return self::installerFor($a['type']);
            case 'installPath':
                return self::$composer->getInstallationManager()->getInstallPath(self::pkg($a['p']));
            case 'isInstalled':
                // Factory::purgePackages asking a plugin installer
                $installer = self::$installerIds[$a['id']] ?? null;
                if ($installer === null) {
                    throw new \LogicException('Unknown installer');
                }

                return $installer->isInstalled(self::$composer->getRepositoryManager()->getLocalRepository(), self::pkg($a['p']));
            case 'installerOp':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return self::installerOp($a);
                });
            case 'settleOps':
                self::settleOps();

                return null;
            case 'appState':
                return Api::appState();
            case 'eventState':
                return Api::eventState($a['eid']);
            case 'downloaderFor':
                return Api::downloaderFor($a['type']);
            case 'downloaderOp':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return Api::downloaderOp($a);
                });
            case 'replacedOp':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return Services::replacedOp($a);
                });
            case 'autoloadVia':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return Services::autoloadVia($a);
                });
            case 'filterCall':
                return Services::filterCall($a);
            case 'archiveChanged':
                return Services::archiveManagerChanged();
            case 'archiveVia':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return Services::archiveVia($a);
                });
            case 'phpRepoOp':
                return self::framed($a['frames'] ?? null, static function () use ($a) {
                    return Services::phpRepoOp($a);
                });
            case 'installerStub':
                Api::installerStub((int) $a['id'], $a['opts'] ?? []);

                return null;
            case 'freeEvent':
                unset(self::$events[$a['eid']]);

                return null;
            case 'freeOps':
                unset(self::$opsCache[$a['key']]);

                return null;
            case 'ping':
                return 'pong';
        }

        throw new \LogicException('Unknown gocomposer plugin host call '.$op);
    }

    private static function init(array $a)
    {
        self::$tmpDir = $a['tmpDir'];
        $formatter = new OutputFormatter((bool) $a['decorated'], \Composer\Factory::createAdditionalStyles());
        self::$output = new GoConsoleOutput((int) $a['verbosity'], (bool) $a['decorated'], $formatter);
        self::$input = new ArgvInput(['composer']);
        self::$input->setInteractive((bool) $a['interactive']);
        $helperSet = new HelperSet([new QuestionHelper()]);
        self::$io = new ConsoleIO(self::$input, self::$output, $helperSet);

        $composer = new Composer();
        $config = new Config();
        $config->gcDetached = false;
        $composer->setConfig($config);
        self::$rootSynced = null;
        self::syncRoot($a['root']);
        $rm = new RepositoryManager();
        $rm->setLocalRepository(new InstalledFilesystemRepository());
        $composer->setRepositoryManager($rm);
        $im = new InstallationManager();
        $composer->setInstallationManager($im);
        $composer->setEventDispatcher(new EventDispatcher($composer, self::$io));
        $composer->setLocker(new Locker());
        $composer->setPluginManager(new PluginManager());
        if (!empty($a['global'])) {
            $composer->setGlobal();
        }
        self::$composer = $composer;
        self::$instances[0] = $composer;
        self::setRoot(self::$pendingRoot);
        $im->setDefaultInstallers(self::$io, $composer);
        Api::initApp($a['app'] ?? null);

        return null;
    }

    /** @var RootPackageInterface|null */
    private static $pendingRoot = null;

    private static function setRoot(RootPackageInterface $root): void
    {
        if (self::$composer !== null) {
            self::$composer->setPackage($root);
        }
        self::$pendingRoot = $root;
        self::$pkgs[0] = $root;
        self::$objIds[spl_object_id($root)] = 0;
    }

    /**
     * Loads the root package data sent by gocomposer (whenever gocomposer
     * changed it, e.g. require --dry-run).
     */
    private static function syncRoot(array $data): void
    {
        self::setRoot(self::loadRoot($data));
        self::$rootSynced = self::rootDump();
    }

    /**
     * Loads root package data sent by gocomposer.
     *
     * @param array<string, mixed> $data
     */
    public static function loadRoot(array $data): RootPackageInterface
    {
        $extra = $data['_root'] ?? [];
        unset($data['_root']);
        $loader = new ArrayLoader(null, true);
        $new = $loader->load($data, RootPackage::class);
        $base = $new instanceof RootAliasPackage ? $new->getAliasOf() : $new;
        if ($base instanceof RootPackage) {
            $base->setMinimumStability($extra['minimum-stability'] ?? 'stable');
            $base->setPreferStable((bool) ($extra['prefer-stable'] ?? false));
            $base->setStabilityFlags($extra['stability-flags'] ?? []);
            $base->setReferences($extra['references'] ?? []);
            $base->setAliases($extra['aliases'] ?? []);
            $base->setConfig($extra['config'] ?? []);
        }

        return $new;
    }

    private static function rootDump(): string
    {
        $dumper = new ArrayDumper();
        $root = self::$pendingRoot;
        if ($root instanceof RootAliasPackage) {
            $root = $root->getAliasOf();
        }

        return (string) json_encode($dumper->dump($root), JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_INVALID_UTF8_SUBSTITUTE | JSON_PARTIAL_OUTPUT_ON_ERROR);
    }

    /**
     * @return array<string, mixed>|null the root package data when a plugin modified it
     */
    private static function rootChanges(): ?array
    {
        if (self::$pendingRoot === null) {
            return null;
        }
        $cur = self::rootDump();
        if ($cur === self::$rootSynced) {
            return null;
        }
        self::$rootSynced = $cur;

        return json_decode($cur, true);
    }

    /**
     * Resolves a package reference sent by gocomposer.
     */
    public static function pkg(?array $ref): ?PackageInterface
    {
        if ($ref === null) {
            return null;
        }
        $id = $ref['id'] ?? null;
        if ($id !== null && isset(self::$pkgs[$id])) {
            return self::$pkgs[$id];
        }
        if (!isset($ref['d'])) {
            throw new \LogicException('Unknown package reference '.json_encode($ref));
        }
        $d = $ref['d'];
        if (isset($d['_alias_of'])) {
            $base = self::pkg($d['_alias_of']);
            if ($base instanceof RootPackage) {
                $p = new RootAliasPackage($base, $d['version_normalized'], $d['version']);
            } elseif ($base instanceof CompletePackage) {
                $p = new CompleteAliasPackage($base, $d['version_normalized'], $d['version']);
            } else {
                $p = new AliasPackage($base, $d['version_normalized'], $d['version']);
            }
            if (!empty($d['_root_alias'])) {
                $p->setRootPackageAlias(true);
            }
        } else {
            unset($d['_alias_of'], $d['_root_alias']);
            $loader = new ArrayLoader(null, true);
            $p = $loader->load($d);
            if ($p instanceof AliasPackage) {
                $p = $p->getAliasOf();
            }
        }
        if ($id !== null) {
            self::$pkgs[$id] = $p;
            self::$objIds[spl_object_id($p)] = $id;
        }

        return $p;
    }

    /**
     * Reference of a package for gocomposer: the id of packages it sent, or
     * the dumped package data for packages created by plugins.
     *
     * @return array<string, mixed>
     */
    public static function ref(PackageInterface $p): array
    {
        $oid = spl_object_id($p);
        if (isset(self::$objIds[$oid]) && isset(self::$pkgs[self::$objIds[$oid]]) && self::$pkgs[self::$objIds[$oid]] === $p) {
            return ['id' => self::$objIds[$oid]];
        }
        $dumper = new ArrayDumper();
        $d = $dumper->dump($p instanceof AliasPackage ? $p->getAliasOf() : $p);
        if ($p instanceof AliasPackage) {
            return ['d' => $d, 'alias' => ['version' => $p->getPrettyVersion(), 'version_normalized' => $p->getVersion()]];
        }

        return ['d' => $d];
    }

    /**
     * Runs $fn under the stack frames Composer would have at that point
     * (Application::doRun with the command input, GlobalCommand, Installer),
     * which plugins such as symfony/flex look up with debug_backtrace().
     *
     * @param array<string, mixed>|null $frames
     * @return mixed
     */
    public static function framed(?array $frames, callable $fn)
    {
        $call = $fn;
        if ($frames !== null) {
            if (isset($frames['installer'])) {
                $installer = Api::installerObj((int) $frames['installer']);
                $inner = $call;
                $call = static function () use ($installer, $inner) {
                    return $installer->gocomposerFrame($inner);
                };
            }
            if (!empty($frames['global'])) {
                $global = Api::globalCommand();
                $inner = $call;
                $call = static function () use ($global, $inner) {
                    return $global->gocomposerFrame($inner);
                };
            }
            if (!empty($frames['app']) && Api::$app !== null) {
                $app = Api::$app;
                $input = Api::$appInput;
                $inner = $call;
                $call = static function () use ($app, $input, $inner) {
                    return $app->gocomposerFrame($input, $inner);
                };
            }
        }

        return $call();
    }

    private static function register(array $a): int
    {
        $package = self::pkg($a['pkg']);
        $isGlobal = (bool) $a['global'];
        $classLoader = new \Composer\Autoload\ClassLoader($a['vendorDir']);
        foreach ($a['psr0'] as $entry) {
            $classLoader->add($entry[0], $entry[1]);
        }
        foreach ($a['psr4'] as $entry) {
            $classLoader->addPsr4($entry[0], $entry[1]);
        }
        if (!empty($a['classmap'])) {
            $classLoader->addClassMap($a['classmap']);
        }
        $classLoader->register(false);

        foreach ($a['files'] as $entry) {
            // exclude laminas/laminas-zendframework-bridge:src/autoload.php as it breaks Composer in some conditions
            if ($entry[0] === '7e9bd612cc444b3eed788ebbe46263a0') {
                continue;
            }
            \Composer\Autoload\composerRequire($entry[0], $entry[1]);
        }

        $count = 0;
        foreach ($a['classes'] as $class) {
            if (class_exists($class, false)) {
                $class = trim($class, '\\');
                $path = $classLoader->findFile($class);
                $code = file_get_contents($path);
                $separatorPos = strrpos($class, '\\');
                $className = $class;
                if ($separatorPos) {
                    $className = substr($class, $separatorPos + 1);
                }
                $code = Preg::replace('{^((?:(?:final|readonly)\s+)*(?:\s*))class\s+('.preg_quote($className).')}mi', '$1class $2_composer_tmp'.self::$classCounter, $code, 1);
                $code = strtr($code, [
                    '__FILE__' => var_export($path, true),
                    '__DIR__' => var_export(dirname($path), true),
                    '__CLASS__' => var_export($class, true),
                ]);
                $code = Preg::replace('/^\s*<\?(php)?/i', '', $code, 1);
                eval($code);
                $class .= '_composer_tmp'.self::$classCounter;
                self::$classCounter++;
            }

            if (!empty($a['legacy'])) {
                if (!is_a($class, 'Composer\Installer\InstallerInterface', true)) {
                    throw new \RuntimeException('Could not activate plugin "'.$package->getName().'" as "'.$class.'" does not implement Composer\Installer\InstallerInterface');
                }
                self::$io->writeError('<warning>Loading "'.$package->getName() . '" '.($isGlobal ? '(installed globally) ' : '').'which is a legacy composer-installer built for Composer 1.x, it is likely to cause issues as you are running Composer 2.x.</warning>');
                $installer = new $class(self::$io, self::$composer);
                self::$composer->getInstallationManager()->addInstaller($installer);
                self::$registered[$package->getName()][] = $installer;
                $count++;
            } elseif (class_exists($class)) {
                if (!is_a($class, 'Composer\Plugin\PluginInterface', true)) {
                    throw new \RuntimeException('Could not activate plugin "'.$package->getName().'" as "'.$class.'" does not implement Composer\Plugin\PluginInterface');
                }
                $plugin = new $class();
                self::addPlugin($plugin, $isGlobal, $package);
                self::$registered[$package->getName()][] = $plugin;
                $count++;
            } elseif (!empty($a['failOnMissing'])) {
                throw new \UnexpectedValueException('Plugin '.$package->getName().' could not be initialized, class not found: '.$class);
            }
        }

        return $count;
    }

    private static function addPlugin(PluginInterface $plugin, bool $isGlobalPlugin, PackageInterface $sourcePackage): void
    {
        $details = ['from '.$sourcePackage->getName()];
        if ($isGlobalPlugin) {
            $details[] = 'installed globally';
        }
        self::$io->writeError('Loading plugin '.get_class($plugin).' ('.implode(', ', $details).')', true, IOInterface::DEBUG);
        self::$plugins[] = $plugin;
        $plugin->activate(self::$composer, self::$io);

        if ($plugin instanceof EventSubscriberInterface) {
            self::$composer->getEventDispatcher()->addSubscriber($plugin);
        }
    }

    private static function deactivatePackage(string $name, bool $uninstall): bool
    {
        if (!isset(self::$registered[$name])) {
            return false;
        }

        $plugins = self::$registered[$name];
        foreach ($plugins as $plugin) {
            if ($plugin instanceof InstallerInterface) {
                self::$composer->getInstallationManager()->removeInstaller($plugin);
            } else {
                self::removePlugin($plugin);
                if ($uninstall) {
                    self::$io->writeError('Uninstalling plugin '.get_class($plugin), true, IOInterface::DEBUG);
                    $plugin->uninstall(self::$composer, self::$io);
                }
            }
        }
        unset(self::$registered[$name]);

        return true;
    }

    private static function removePlugin(PluginInterface $plugin): void
    {
        $index = array_search($plugin, self::$plugins, true);
        if ($index === false) {
            return;
        }

        self::$io->writeError('Unloading plugin '.get_class($plugin), true, IOInterface::DEBUG);
        unset(self::$plugins[$index]);
        self::$plugins = array_values(self::$plugins);
        $plugin->deactivate(self::$composer, self::$io);

        self::removeListener($plugin);
    }

    /**
     * @param mixed $listener
     */
    public static function addListener(string $eventName, $listener, int $priority = 0): void
    {
        $id = ++self::$listenerSeq;
        self::$listeners[$eventName][$priority][] = ['id' => $id, 'cb' => $listener];
        self::$listenerById[$id] = $listener;
    }

    /**
     * @param mixed $listener
     */
    public static function removeListener($listener): void
    {
        foreach (self::$listeners as $eventName => $priorities) {
            foreach ($priorities as $priority => $listeners) {
                foreach ($listeners as $index => $candidate) {
                    $cb = $candidate['cb'];
                    if ($listener === $cb || (is_array($cb) && is_object($listener) && $cb[0] === $listener)) {
                        unset(self::$listeners[$eventName][$priority][$index]);
                        unset(self::$listenerById[$candidate['id']]);
                    }
                }
            }
        }
    }

    /**
     * Listeners of an event ordered by priority (highest first), with their
     * priority so gocomposer can merge the root package scripts at 0.
     *
     * @return list<array{id: int, prio: int, label: string|null}>
     */
    private static function listListeners(string $name): array
    {
        if (!isset(self::$listeners[$name])) {
            return [];
        }
        $byPrio = self::$listeners[$name];
        krsort($byPrio);
        $out = [];
        foreach ($byPrio as $prio => $listeners) {
            foreach ($listeners as $l) {
                $cb = $l['cb'];
                $label = null;
                $key = 'unknown';
                if (is_array($cb) && isset($cb[0], $cb[1]) && (is_string($cb[0]) || is_object($cb[0])) && is_string($cb[1])) {
                    $label = (is_object($cb[0]) ? get_class($cb[0]) : $cb[0]).'->'.$cb[1];
                    $key = (is_object($cb[0]) ? get_class($cb[0]) : $cb[0]).'::'.$cb[1];
                } elseif (is_string($cb)) {
                    $label = null;
                    $key = $cb;
                } elseif ($cb instanceof \Closure) {
                    $key = 'closure';
                }
                $out[] = ['id' => $l['id'], 'prio' => (int) $prio, 'label' => $label, 'key' => $key, 'string' => is_string($cb) ? $cb : null];
            }
        }

        return $out;
    }

    /**
     * Creates (or reuses) the event object of a dispatch run by gocomposer.
     */
    public static function event(array $spec): Event
    {
        $eid = $spec['eid'] ?? null;
        if ($eid !== null && isset(self::$events[$eid])) {
            return self::$events[$eid];
        }
        $name = $spec['name'];
        $args = $spec['args'] ?? [];
        $flags = $spec['flags'] ?? [];
        $dev = (bool) ($spec['dev'] ?? false);
        switch ($spec['kind'] ?? 'script') {
            case 'plain':
                $event = new Event($name, $args, $flags);
                break;
            case 'package':
                $event = new PackageEvent($name, self::$composer, self::$io, $dev, self::$composer->getRepositoryManager()->getLocalRepository(), self::ops($spec['opsKey'], $spec['ops'] ?? null), self::op($spec['op']));
                break;
            case 'installer':
                $event = new InstallerEvent($name, self::$composer, self::$io, $dev, (bool) $spec['exec'], LocalRepoTransaction::gcCreate(self::ops($spec['opsKey'], $spec['ops'] ?? null), Api::pkgs($spec['present'] ?? []), Api::pkgs($spec['result'] ?? [])));
                break;
            case 'command':
                $event = new CommandEvent($name, $spec['command'], self::input($spec['input'] ?? null), self::$output, $args, $flags);
                break;
            case 'precommand':
                $event = new PreCommandRunEvent($name, self::input($spec['input'] ?? null), $spec['command']);
                break;
            case 'prepool':
            case 'prefiledl':
            case 'postfiledl':
                $event = Api::event($spec);
                break;
            default:
                $event = new ScriptEvent($name, self::$composer, self::$io, $dev, $args, $flags);
                if (isset($spec['origin']) && isset(self::$events[$spec['origin']])) {
                    $event->setOriginatingEvent(self::$events[$spec['origin']]);
                }
        }
        if ($eid !== null) {
            self::$events[$eid] = $event;
        }

        return $event;
    }

    /**
     * @return list<\Composer\DependencyResolver\Operation\OperationInterface>
     */
    private static function ops(string $key, ?array $specs): array
    {
        if (!isset(self::$opsCache[$key])) {
            $ops = [];
            foreach ($specs ?? [] as $spec) {
                $ops[] = self::op($spec);
            }
            self::$opsCache[$key] = $ops;
        }

        return self::$opsCache[$key];
    }

    /**
     * @param array<string, mixed> $spec
     */
    public static function opFromSpec(array $spec): \Composer\DependencyResolver\Operation\OperationInterface
    {
        return self::op($spec);
    }

    private static function op(array $spec): \Composer\DependencyResolver\Operation\OperationInterface
    {
        switch ($spec['type']) {
            case 'install':
                return new InstallOperation(self::pkg($spec['p']));
            case 'update':
                return new UpdateOperation(self::pkg($spec['i']), self::pkg($spec['p']));
            case 'uninstall':
                return new UninstallOperation(self::pkg($spec['p']));
            case 'markAliasInstalled':
                return new MarkAliasInstalledOperation(self::pkg($spec['p']));
            case 'markAliasUninstalled':
                return new MarkAliasUninstalledOperation(self::pkg($spec['p']));
        }

        throw new \LogicException('Unknown operation type '.$spec['type']);
    }

    /**
     * The command line input of the running gocomposer command, bound to its
     * definition like Composer's input.
     */
    private static function input(?array $spec): \Symfony\Component\Console\Input\InputInterface
    {
        if ($spec === null) {
            return self::$input;
        }
        $argv = array_merge(['composer'], $spec['argv']);
        try {
            $definition = new InputDefinition();
            $definition->addArgument(new InputArgument('command', InputArgument::REQUIRED, 'The command to execute'));
            foreach ($spec['args'] as $arg) {
                $definition->addArgument(new InputArgument($arg['name'], $arg['mode'], $arg['desc'] ?? '', $arg['default'] ?? null));
            }
            foreach ($spec['opts'] as $opt) {
                $definition->addOption(new InputOption($opt['name'], $opt['shortcut'] ?? null, $opt['mode'], $opt['desc'] ?? '', $opt['default'] ?? null));
            }
            $input = new ArgvInput($argv, $definition);
        } catch (\Throwable $e) {
            $input = new ArgvInput($argv);
        }
        $input->setInteractive(self::$input->isInteractive());

        return $input;
    }

    private static function callListener(array $a): array
    {
        if (isset($a['loader'])) {
            self::projectLoader($a['loader']);
        }
        if (isset($a['runScripts'])) {
            $dispatcher = self::$composer->getEventDispatcher();
            \Closure::bind(static function () use ($dispatcher, $a): void {
                $dispatcher->runScripts = (bool) $a['runScripts'];
            }, null, EventDispatcher::class)();
        }
        $event = self::event($a['event']);
        if (!isset(self::$listenerById[$a['id']])) {
            return ['ret' => 0, 'stopped' => $event->isPropagationStopped()];
        }
        $callable = self::$listenerById[$a['id']];
        $additionalArgs = $event->getArguments();
        $formattedEventNameWithArgs = $event->getName() . ($additionalArgs !== [] ? ' (' . implode(', ', $additionalArgs) . ')' : '');
        if (!is_callable($callable)) {
            $className = is_object($callable[0]) ? get_class($callable[0]) : $callable[0];

            throw new \RuntimeException('Subscriber '.$className.'::'.$callable[1].' for event '.$event->getName().' is not callable, make sure the function is defined and public');
        }
        if (is_array($callable) && (is_string($callable[0]) || is_object($callable[0])) && is_string($callable[1])) {
            self::$io->writeError(sprintf('> %s: %s', $formattedEventNameWithArgs, (is_object($callable[0]) ? get_class($callable[0]) : $callable[0]).'->'.$callable[1]), true, IOInterface::VERBOSE);
        }
        $return = false === $callable($event) ? 1 : 0;

        return ['ret' => $return, 'stopped' => $event->isPropagationStopped()];
    }

    /**
     * EventDispatcher::makeAutoloader(): the project autoloader (without
     * files) generated by gocomposer for the listener calls.
     */
    private static function projectLoader(string $file): void
    {
        if (self::$projectLoader !== null) {
            self::$projectLoader->unregister();
        }
        $loader = (static function (string $file) {
            $loader = null;
            require $file;

            return $loader;
        })($file);
        self::$projectLoader = $loader;
    }

    /**
     * Dispatches an event a plugin dispatched itself through gocomposer
     * (script listeners and plugin listeners alike).
     */
    public static function dispatchFromPlugin(Event $event, int $inst = 0, ?bool $runScripts = null): int
    {
        $eid = 'p'.spl_object_id($event).'-'.(++self::$listenerSeq);
        self::$events[$eid] = $event;
        $spec = ['eid' => $eid, 'name' => $event->getName(), 'args' => $event->getArguments(), 'flags' => $event->getFlags()];
        if ($event instanceof ScriptEvent) {
            $spec['kind'] = 'script';
            $spec['dev'] = $event->isDevMode();
        } elseif ($event instanceof PackageEvent || $event instanceof InstallerEvent) {
            $spec['kind'] = 'unsupported';
        } elseif ($event instanceof PreCommandRunEvent) {
            $spec['kind'] = 'precommand';
            $spec['command'] = $event->getCommand();
        } elseif ($event instanceof CommandEvent) {
            $spec['kind'] = 'command';
            $spec['command'] = $event->getCommandName();
        } else {
            $spec['kind'] = 'plain';
        }
        try {
            if ($inst !== 0) {
                // the event dispatcher of a Composer instance created with
                // Factory::create(): gocomposer dispatches with that instance
                return (int) self::cb('nestedDispatch', $spec + ['inst' => $inst, 'runScripts' => $runScripts]);
            }

            return (int) self::cb('dispatch', $spec);
        } finally {
            unset(self::$events[$eid]);
        }
    }

    public static function installersChanged(): void
    {
        if (self::$composer !== null) {
            self::cb('installersChanged', []);
        }
    }

    /**
     * @return array<string, mixed>|null
     */
    private static function installerFor(string $type): ?array
    {
        $installer = self::$composer->getInstallationManager()->getPluginInstaller($type);
        if ($installer === null) {
            return null;
        }
        $id = spl_object_id($installer);
        self::$installerIds[$id] = $installer;
        $library = $installer instanceof LibraryInstaller;
        $overrides = [];
        $ref = new \ReflectionClass($installer);
        $methods = $library
            ? ['download', 'prepare', 'cleanup', 'installCode', 'updateCode', 'removeCode', 'install', 'update', 'uninstall']
            : ['download', 'prepare', 'cleanup', 'install', 'update', 'uninstall'];
        foreach ($methods as $m) {
            if (!$ref->hasMethod($m)) {
                continue;
            }
            $declaring = $ref->getMethod($m)->getDeclaringClass()->getName();
            if (!$library || !in_array($declaring, [LibraryInstaller::class, \Composer\Installer\PluginInstaller::class], true)) {
                $overrides[] = $m;
            }
        }
        if ($library && $installer instanceof LibraryInstaller) {
            $installer->gcSetPortedFlow(array_values(array_intersect($overrides, ['download', 'prepare', 'cleanup', 'installCode', 'updateCode', 'removeCode'])) !== []);
        }

        return ['id' => $id, 'class' => get_class($installer), 'library' => $library, 'overrides' => $overrides];
    }

    /**
     * @return mixed
     */
    private static function installerOp(array $a)
    {
        $installer = self::$installerIds[$a['id']] ?? null;
        if ($installer === null) {
            throw new \LogicException('Unknown installer');
        }
        $repo = self::$composer->getRepositoryManager()->getLocalRepository();
        $p = self::pkg($a['p']);
        $i = self::pkg($a['i'] ?? null);
        // the installation source gocomposer chose for the operation
        if (array_key_exists('psrc', $a) && $p !== null && $p->getInstallationSource() !== $a['psrc']) {
            $p->setInstallationSource($a['psrc']);
        }
        if (array_key_exists('isrc', $a) && $i !== null && $i->getInstallationSource() !== $a['isrc']) {
            $i->setInstallationSource($a['isrc']);
        }
        if (!empty($a['defer'])) {
            // executeBatch runs the operation in the prepare promise's then():
            // exceptions become rejections and the next operations still start
            $m = $a['m'];
            $promise = \React\Promise\resolve(null)->then(static function () use ($installer, $m, $repo, $p, $i) {
                return $m === 'update' ? $installer->update($repo, $i, $p) : $installer->$m($repo, $p);
            });
            $a['m'] = 'deferred';
        }
        switch ($a['m']) {
            case 'deferred':
                $a['m'] = $m;
                break;
            case 'install':
                $promise = $installer->install($repo, $p);
                break;
            case 'update':
                $promise = $installer->update($repo, $i, $p);
                break;
            case 'uninstall':
                $promise = $installer->uninstall($repo, $p);
                break;
            case 'download':
                $promise = $installer->download($p, $i);
                break;
            case 'prepare':
            case 'cleanup':
                $method = $a['m'];
                $promise = $installer->$method($a['type'], $p, $i);
                break;
            default:
                throw new \LogicException('Unknown installer operation '.$a['m']);
        }
        if (empty($a['defer'])) {
            Api::settle($promise);

            return null;
        }
        // InstallationManager::executeBatch: the operation, the cleanup of the
        // installer chosen when downloading, then the repository write
        // (done by gocomposer), settled with the batch's other operations
        if (!$promise instanceof \React\Promise\PromiseInterface) {
            $promise = \React\Promise\resolve(null);
        }
        $opType = $a['m'];
        $cleanup = isset($a['cleanup']) ? (self::$installerIds[$a['cleanup']] ?? null) : null;
        $io = self::$io;
        $key = $a['key'];
        $promise = $promise->then(static function () use ($cleanup, $opType, $p, $i) {
            if ($cleanup === null || null === $p->getInstallationSource()) {
                return \React\Promise\resolve(null);
            }

            return $cleanup->cleanup($opType, $p, $i);
        })->then(static function () use ($key): void {
            Host::cb('opDone', ['key' => $key]);
        }, static function ($e) use ($opType, $p, $io): void {
            $io->writeError('    <error>' . ucfirst($opType) .' of '.$p->getPrettyName().' failed</error>');

            throw $e;
        });
        self::$pendingOps[] = $promise;

        return null;
    }

    /** @var list<\React\Promise\PromiseInterface> operations started by installerOp with defer */
    public static $pendingOps = [];

    /**
     * InstallationManager::waitOnPromises for the deferred operations.
     */
    private static function settleOps(): void
    {
        $pending = self::$pendingOps;
        self::$pendingOps = [];
        if ($pending === []) {
            return;
        }
        self::$composer->getLoop()->wait($pending);
    }

    /**
     * Runs a promise chain to completion (the promises of the runtime
     * installers are already resolved; plugin callbacks run synchronously).
     *
     * @param mixed $promise
     */
    private static function settle($promise): void
    {
        if (!$promise instanceof \React\Promise\PromiseInterface) {
            return;
        }
        $error = null;
        $promise->then(null, static function ($e) use (&$error): void {
            $error = $e;
        });
        if ($error instanceof \Throwable) {
            throw $error;
        }
    }

    /**
     * The commands of the CommandProvider capabilities of the active plugins
     * (Application::getPluginCommands()).
     *
     * @return list<array<string, mixed>>
     */
    private static function commands(): array
    {
        $out = [];
        self::$commands = [];
        $app = new \Composer\Console\Application();
        $app->setAutoExit(false);
        $app->setCatchExceptions(false);
        $pm = self::$composer->getPluginManager();
        foreach ($pm->getPluginCapabilities('Composer\Plugin\Capability\CommandProvider', ['composer' => self::$composer, 'io' => self::$io]) as $capability) {
            $newCommands = $capability->getCommands();
            if (!is_array($newCommands)) {
                throw new \UnexpectedValueException('Plugin capability '.get_class($capability).' failed to return an array from getCommands');
            }
            foreach ($newCommands as $command) {
                if (!$command instanceof \Composer\Command\BaseCommand) {
                    throw new \UnexpectedValueException('Plugin capability '.get_class($capability).' returned an invalid value, we expected an array of Composer\Command\BaseCommand objects');
                }
            }
            foreach ($newCommands as $command) {
                $idx = count(self::$commands);
                self::$commands[$idx] = $command;
                $def = $command->getDefinition();
                $args = [];
                foreach ($def->getArguments() as $arg) {
                    $args[] = ['name' => $arg->getName(), 'required' => $arg->isRequired(), 'array' => $arg->isArray(), 'desc' => $arg->getDescription(), 'default' => $arg->getDefault()];
                }
                $opts = [];
                foreach ($def->getOptions() as $opt) {
                    $opts[] = ['name' => $opt->getName(), 'shortcut' => $opt->getShortcut(), 'none' => !$opt->acceptValue(), 'required' => $opt->isValueRequired(), 'optional' => $opt->isValueOptional(), 'array' => $opt->isArray(), 'negatable' => $opt->isNegatable(), 'desc' => $opt->getDescription(), 'default' => $opt->getDefault()];
                }
                $command->setApplication($app);
                $out[] = self::describeCommand($idx, $command, $args, $opts);
            }
        }
        foreach (Api::appCommands() as $command) {
            $idx = count(self::$commands);
            self::$commands[$idx] = $command;
            $def = $command->getDefinition();
            $args = [];
            foreach ($def->getArguments() as $arg) {
                $args[] = ['name' => $arg->getName(), 'required' => $arg->isRequired(), 'array' => $arg->isArray(), 'desc' => $arg->getDescription(), 'default' => $arg->getDefault()];
            }
            $opts = [];
            foreach ($def->getOptions() as $opt) {
                $opts[] = ['name' => $opt->getName(), 'shortcut' => $opt->getShortcut(), 'none' => !$opt->acceptValue(), 'required' => $opt->isValueRequired(), 'optional' => $opt->isValueOptional(), 'array' => $opt->isArray(), 'negatable' => $opt->isNegatable(), 'desc' => $opt->getDescription(), 'default' => $opt->getDefault()];
            }
            $out[] = self::describeCommand($idx, $command, $args, $opts);
        }

        return $out;
    }

    /**
     * @param list<array<string, mixed>> $args
     * @param list<array<string, mixed>> $opts
     * @return array<string, mixed>
     */
    private static function describeCommand(int $idx, \Symfony\Component\Console\Command\Command $command, array $args, array $opts): array
    {
        return [
            'index' => $idx,
            'name' => $command->getName(),
            'class' => get_class($command),
            'description' => $command->getDescription(),
            'help' => $command->getHelp(),
            'aliases' => $command->getAliases(),
            'hidden' => $command->isHidden(),
            'synopsis' => $command->getSynopsis(false),
            'usages' => $command->getUsages(),
            'args' => $args,
            'opts' => $opts,
        ];
    }

    private static function runCommand(array $a): int
    {
        $command = self::$commands[$a['index']] ?? null;
        if ($command === null) {
            throw new \LogicException('Unknown plugin command');
        }
        $app = $command->getApplication();
        $app->setAutoExit(false);
        $app->setCatchExceptions(false);
        if (!$app->has((string) $command->getName())) {
            $app->add($command);
        }
        $input = new ArgvInput(array_merge(['composer'], $a['argv']));
        $input->setInteractive(self::$input->isInteractive());

        return $app->run($input, self::$output);
    }

    public static function errorHandler(int $level, string $message, string $file, int $line): bool
    {
        // Composer\Util\ErrorHandler::handle of Composer f83588392df43d4fd027d84d26de459946ef4e80
        // (MIT license, (c) Nils Adermann, Jordi Boggiano)
        $isDeprecationNotice = $level === E_DEPRECATED || $level === E_USER_DEPRECATED;

        if (!$isDeprecationNotice && 0 === (error_reporting() & $level)) {
            return true;
        }

        if (filter_var(ini_get('xdebug.scream'), FILTER_VALIDATE_BOOLEAN)) {
            $message .= "\n\nWarning: You have xdebug.scream enabled, the warning above may be".
            "\na legitimately suppressed error that you were not supposed to see.";
        }

        if (!$isDeprecationNotice) {
            if (in_array($level, [E_WARNING, E_USER_WARNING], true) && str_contains($message, 'should either be used or intentionally ignored by casting it as (void)')) {
                self::outputWarning('Ignored new PHP warning but it should be reported and fixed: '.$message.' in '.$file.':'.$line);

                return true;
            }

            throw new \ErrorException($message, 0, $level, $file, $line);
        }

        if (self::$io !== null) {
            if (self::$hasShownDeprecationNotice > 0 && !self::$io->isVerbose()) {
                if (self::$hasShownDeprecationNotice === 1) {
                    self::$io->writeError('<warning>More deprecation notices were hidden, run again with `-v` to show them.</warning>');
                    self::$hasShownDeprecationNotice = 2;
                }

                return true;
            }
            self::$hasShownDeprecationNotice = 1;
            self::outputWarning('Deprecation Notice: '.$message.' in '.$file.':'.$line);
        }

        return true;
    }

    private static function outputWarning(string $message): void
    {
        if (self::$io !== null) {
            self::$io->writeError('<warning>'.$message.'</warning>');
            if (self::$io->isVerbose()) {
                self::$io->writeError('<warning>Stack trace:</warning>');
                self::$io->writeError(array_filter(array_map(static function ($a): ?string {
                    if (isset($a['line'], $a['file'])) {
                        return '<warning> '.$a['file'].':'.$a['line'].'</warning>';
                    }

                    return null;
                }, array_slice(debug_backtrace(), 2)), static function (?string $line) {
                    return $line !== null;
                }));
            }
        }
    }
}

}
