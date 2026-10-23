<?php declare(strict_types=1);
// gocomposer plugin runtime: ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Command/BaseCommand.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano) for commands provided by
// plugins (CommandProvider capability). Helpers depending on Composer
// internals gocomposer does not expose (policy/audit configuration,
// completion) are not available.

namespace Composer\Command;

use Composer\Composer;
use Composer\Console\Application;
use Composer\Factory;
use Composer\IO\IOInterface;
use Composer\IO\NullIO;
use Composer\Package\Version\VersionParser;
use Composer\Plugin\PluginEvents;
use Composer\Plugin\PreCommandRunEvent;
use Composer\Util\Platform;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Terminal;

abstract class BaseCommand extends Command
{
    /** @var Composer|null */
    private $composer;

    /** @var IOInterface */
    private $io;

    public function getApplication(): Application
    {
        $application = parent::getApplication();
        if (!$application instanceof Application) {
            throw new \RuntimeException('Composer commands can only work with an '.Application::class.' instance set');
        }

        return $application;
    }

    /**
     * @deprecated since Composer 2.3.0 use requireComposer or tryComposer depending on whether you have $required set to true or false
     * @return Composer|null
     */
    public function getComposer(bool $required = true, ?bool $disablePlugins = null, ?bool $disableScripts = null)
    {
        if ($required) {
            return $this->requireComposer($disablePlugins, $disableScripts);
        }

        return $this->tryComposer($disablePlugins, $disableScripts);
    }

    public function requireComposer(?bool $disablePlugins = null, ?bool $disableScripts = null): Composer
    {
        if (null === $this->composer) {
            $application = parent::getApplication();
            if ($application instanceof Application) {
                $this->composer = $application->getComposer(true, $disablePlugins, $disableScripts);
                assert($this->composer instanceof Composer);
            } else {
                throw new \RuntimeException(
                    'Could not create a Composer\Composer instance, you must inject '.
                    'one if this command is not used with a Composer\Console\Application instance'
                );
            }
        }

        return $this->composer;
    }

    public function tryComposer(?bool $disablePlugins = null, ?bool $disableScripts = null): ?Composer
    {
        if (null === $this->composer) {
            $application = parent::getApplication();
            if ($application instanceof Application) {
                $this->composer = $application->getComposer(false, $disablePlugins, $disableScripts);
            }
        }

        return $this->composer;
    }

    /**
     * @return void
     */
    public function setComposer(Composer $composer)
    {
        $this->composer = $composer;
    }

    /**
     * @return void
     */
    public function resetComposer()
    {
        $this->composer = null;
        $this->getApplication()->resetComposer();
    }

    /**
     * @return bool
     */
    public function isProxyCommand()
    {
        return false;
    }

    /**
     * @return IOInterface
     */
    public function getIO()
    {
        if (null === $this->io) {
            $application = parent::getApplication();
            if ($application instanceof Application) {
                $this->io = $application->getIO();
            } else {
                $this->io = new NullIO();
            }
        }

        return $this->io;
    }

    /**
     * @return void
     */
    public function setIO(IOInterface $io)
    {
        $this->io = $io;
    }

    protected function initialize(InputInterface $input, OutputInterface $output): void
    {
        $composer = $this->tryComposer();
        $io = $this->getIO();

        if ($composer) {
            $preCommandRunEvent = new PreCommandRunEvent(PluginEvents::PRE_COMMAND_RUN, $input, (string) $this->getName());
            $composer->getEventDispatcher()->dispatch($preCommandRunEvent->getName(), $preCommandRunEvent);
        }

        if (true === $input->hasParameterOption(['--no-ansi']) && $input->hasOption('no-progress')) {
            $input->setOption('no-progress', true);
        }

        $envOptions = [
            'COMPOSER_NO_AUDIT' => ['no-audit'],
            'COMPOSER_NO_DEV' => ['no-dev', 'update-no-dev'],
            'COMPOSER_PREFER_STABLE' => ['prefer-stable'],
            'COMPOSER_PREFER_LOWEST' => ['prefer-lowest'],
            'COMPOSER_MINIMAL_CHANGES' => ['minimal-changes'],
            'COMPOSER_WITH_DEPENDENCIES' => ['with-dependencies'],
            'COMPOSER_WITH_ALL_DEPENDENCIES' => ['with-all-dependencies'],
            'COMPOSER_NO_SECURITY_BLOCKING' => ['no-security-blocking'],
            'COMPOSER_NO_BLOCKING' => ['no-blocking'],
        ];
        foreach ($envOptions as $envName => $optionNames) {
            foreach ($optionNames as $optionName) {
                if (true === $input->hasOption($optionName)) {
                    if (false === $input->getOption($optionName) && (bool) Platform::getEnv($envName)) {
                        $input->setOption($optionName, true);
                    }
                }
            }
        }

        if (true === $input->hasOption('ignore-platform-reqs')) {
            if (!$input->getOption('ignore-platform-reqs') && (bool) Platform::getEnv('COMPOSER_IGNORE_PLATFORM_REQS')) {
                $input->setOption('ignore-platform-reqs', true);

                $io->writeError('<warning>COMPOSER_IGNORE_PLATFORM_REQS is set. You may experience unexpected errors.</warning>');
            }
        }

        if (true === $input->hasOption('ignore-platform-req') && (!$input->hasOption('ignore-platform-reqs') || !$input->getOption('ignore-platform-reqs'))) {
            $ignorePlatformReqEnv = Platform::getEnv('COMPOSER_IGNORE_PLATFORM_REQ');
            if (0 === count($input->getOption('ignore-platform-req')) && is_string($ignorePlatformReqEnv) && '' !== $ignorePlatformReqEnv) {
                $input->setOption('ignore-platform-req', explode(',', $ignorePlatformReqEnv));

                $io->writeError('<warning>COMPOSER_IGNORE_PLATFORM_REQ is set to ignore '.$ignorePlatformReqEnv.'. You may experience unexpected errors.</warning>');
            }
        }

        parent::initialize($input, $output);
    }

    /**
     * @param array<string> $requirements
     * @return array<string, string>
     */
    protected function formatRequirements(array $requirements)
    {
        $requires = [];
        $requirements = $this->normalizeRequirements($requirements);
        foreach ($requirements as $requirement) {
            if (!isset($requirement['version'])) {
                throw new \UnexpectedValueException('Option '.$requirement['name'] .' is missing a version constraint, use e.g. '.$requirement['name'].':^1.0');
            }
            $requires[$requirement['name']] = $requirement['version'];
        }

        return $requires;
    }

    /**
     * @param array<string> $requirements
     * @return list<array{name: string, version?: string}>
     */
    protected function normalizeRequirements(array $requirements)
    {
        $parser = new VersionParser();

        return $parser->parseNameVersionPairs($requirements);
    }

    /**
     * @param array<int, array<int, string>> $table
     * @return void
     */
    protected function renderTable(array $table, OutputInterface $output)
    {
        $renderer = new Table($output);
        $renderer->setStyle('compact');
        $renderer->setRows($table)->render();
    }

    /**
     * @return int
     */
    protected function getTerminalWidth()
    {
        $terminal = new Terminal();
        $width = $terminal->getWidth();

        if (Platform::isWindows()) {
            $width--;
        } else {
            $width = max(80, $width);
        }

        return $width;
    }
}
