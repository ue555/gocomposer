<?php declare(strict_types=1);
// gocomposer plugin runtime: static helpers of Composer\Factory, ported from
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Factory.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). Factory::create() makes
// gocomposer create (and load the plugins of) another Composer instance
// from a composer.json file; the instance is exposed through the same
// gocomposer-backed stubs.

namespace Composer;

use Composer\Util\Platform;
use GocomposerPluginHost\Unsupported;
use Symfony\Component\Console\Formatter\OutputFormatter;
use Symfony\Component\Console\Formatter\OutputFormatterStyle;
use Symfony\Component\Console\Output\ConsoleOutput;

class Factory
{
    public static function getComposerFile(): string
    {
        $env = Platform::getEnv('COMPOSER');
        if (is_string($env)) {
            $env = trim($env);
            if ('' !== $env) {
                if (is_dir($env)) {
                    throw new \RuntimeException('The COMPOSER environment variable is set to '.$env.' which is a directory, this variable should point to a composer.json or be left unset.');
                }

                return $env;
            }
        }

        return './composer.json';
    }

    public static function getLockFile(string $composerFile): string
    {
        return "json" === pathinfo($composerFile, PATHINFO_EXTENSION)
                ? substr($composerFile, 0, -4).'lock'
                : $composerFile . '.lock';
    }

    /**
     * @return array{highlight: OutputFormatterStyle, warning: OutputFormatterStyle}
     */
    public static function createAdditionalStyles(): array
    {
        return [
            'highlight' => new OutputFormatterStyle('red'),
            'warning' => new OutputFormatterStyle('black', 'yellow'),
        ];
    }

    public static function createOutput(): ConsoleOutput
    {
        $styles = self::createAdditionalStyles();
        $formatter = new OutputFormatter(false, $styles);

        return new ConsoleOutput(ConsoleOutput::VERBOSITY_NORMAL, null, $formatter);
    }

    public static function createConfig($io = null, ?string $cwd = null): Config
    {
        return new Config();
    }

    /**
     * @param string|array<string, mixed>|null $config
     * @param bool|'local'|'global' $disablePlugins
     */
    public static function create($io, $config = null, $disablePlugins = false, bool $disableScripts = false): Composer
    {
        if (is_array($config)) {
            Unsupported::fail(__METHOD__.' with an array configuration');
        }
        // for BC reasons, if a config is passed in either as array or a path that is not the default composer.json path
        // we disable local plugins as they really should not be loaded from CWD
        if ($config !== null && $config !== self::getComposerFile() && $disablePlugins === false) {
            $disablePlugins = 'local';
        }
        $desc = \GocomposerPluginHost\Host::cb('factoryCreate', [
            'file' => $config,
            'disablePlugins' => $disablePlugins === true || $disablePlugins === 'local',
            'disableScripts' => $disableScripts,
        ]);

        return \GocomposerPluginHost\Api::nestedComposer($desc);
    }

    public static function createHttpDownloader($io, Config $config, array $options = []): \Composer\Util\HttpDownloader
    {
        return new \Composer\Util\HttpDownloader($io, $config, $options);
    }

    public static function createGlobal($io, bool $disablePlugins = false, bool $disableScripts = false)
    {
        Unsupported::fail(__METHOD__);
    }

    public function createComposer(...$args)
    {
        Unsupported::fail(__METHOD__);
    }

    public static function getHomeDir(): string
    {
        return (string) (new Config())->get('home');
    }
}
