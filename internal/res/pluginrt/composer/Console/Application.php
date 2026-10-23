<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Console\Application used to run
// commands provided by plugins, based on Composer
// f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Console/Application.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The real application is
// gocomposer itself; this one hosts plugin commands and stands for
// gocomposer's application when plugins look it up (commands they add are
// offered by gocomposer, find() resolves gocomposer's commands).

namespace Composer\Console;

use Composer\Composer;
use Composer\IO\IOInterface;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;
use Symfony\Component\Console\Application as BaseApplication;
use Symfony\Component\Console\Input\InputDefinition;
use Symfony\Component\Console\Input\InputOption;

class Application extends BaseApplication
{
    public function __construct(string $name = 'Composer', string $version = '')
    {
        parent::__construct($name, $version === '' ? Composer::getVersion() : $version);
    }

    public function getComposer(bool $required = true, ?bool $disablePlugins = null, ?bool $disableScripts = null): ?Composer
    {
        return Host::$composer;
    }

    public function resetComposer(): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function getIO(): IOInterface
    {
        return Host::$io;
    }

    public function getDisablePluginsByDefault(): bool
    {
        return false;
    }

    public function getDisableScriptsByDefault(): bool
    {
        return false;
    }

    /** @var list<\Symfony\Component\Console\Command\Command> commands added by plugins */
    private $gcAdded = [];
    /** @var array<string, \Composer\Command\BaseCommand> gocomposer commands looked up by plugins */
    private $gcProxies = [];

    /**
     * @internal runs $fn with this application running $input on the call
     * stack (Application::doRun), where plugins like symfony/flex look for it
     * @return mixed
     */
    public function gocomposerFrame(\Symfony\Component\Console\Input\InputInterface $input, callable $fn)
    {
        return $fn();
    }

    /**
     * @return \Symfony\Component\Console\Command\Command|null
     */
    public function add(\Symfony\Component\Console\Command\Command $command)
    {
        $result = parent::add($command);
        if ($result !== null && !$command instanceof \GocomposerPluginHost\ProxyCommand && !in_array(get_class($command), [
            \Symfony\Component\Console\Command\HelpCommand::class,
            \Symfony\Component\Console\Command\ListCommand::class,
            \Symfony\Component\Console\Command\CompleteCommand::class,
            \Symfony\Component\Console\Command\DumpCompletionCommand::class,
        ], true)) {
            $this->gcAdded[] = $command;
        }

        return $result;
    }

    /**
     * @return \Symfony\Component\Console\Command\Command
     */
    public function addCommand($command)
    {
        return $this->add($command);
    }

    /**
     * @internal
     * @return list<\Symfony\Component\Console\Command\Command>
     */
    public function gcAdded(): array
    {
        return $this->gcAdded;
    }

    /**
     * Commands added by plugins, then gocomposer's commands (looked up with
     * gocomposer's command resolution, abbreviations included).
     *
     * @return \Symfony\Component\Console\Command\Command
     */
    public function find(string $name)
    {
        foreach ($this->gcAdded as $command) {
            if ($command->getName() === $name || in_array($name, $command->getAliases(), true)) {
                return $command;
            }
        }
        $info = Host::cb('appFind', ['name' => $name]);
        if (is_array($info) && isset($info['cmd'])) {
            return $this->gcProxy($info['cmd']);
        }

        return parent::find($name);
    }

    /**
     * @internal
     * @param array<string, mixed> $info
     */
    public function gcProxy(array $info): \Composer\Command\BaseCommand
    {
        $name = $info['name'];
        if (!isset($this->gcProxies[$name])) {
            $command = new \GocomposerPluginHost\ProxyCommand($info);
            $command->setApplication($this);
            $this->gcProxies[$name] = $command;
        }

        return $this->gcProxies[$name];
    }

    protected function getDefaultInputDefinition(): InputDefinition
    {
        $definition = parent::getDefaultInputDefinition();
        $definition->addOption(new InputOption('--profile', null, InputOption::VALUE_NONE, 'Display timing and memory usage information'));
        $definition->addOption(new InputOption('--no-plugins', null, InputOption::VALUE_NONE, 'Whether to disable plugins.'));
        $definition->addOption(new InputOption('--no-scripts', null, InputOption::VALUE_NONE, 'Skips the execution of all scripts defined in composer.json file.'));
        $definition->addOption(new InputOption('--working-dir', '-d', InputOption::VALUE_REQUIRED, 'If specified, use the given directory as working directory.'));
        $definition->addOption(new InputOption('--no-cache', null, InputOption::VALUE_NONE, 'Prevent use of the cache'));

        return $definition;
    }
}
