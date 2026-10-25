<?php declare(strict_types=1);
// gocomposer plugin runtime: static helpers of Composer\Factory, ported from
// Composer f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Factory.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). Factory::create(),
// createGlobal() and createComposer() make gocomposer create (and load the
// plugins of) another Composer instance from a composer.json file or a
// configuration array; the instance is exposed through the same
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
        return Config::gcCreate($cwd);
    }

    /**
     * @param string|array<string, mixed>|null $config
     * @param bool|'local'|'global' $disablePlugins
     */
    public static function create($io, $config = null, $disablePlugins = false, bool $disableScripts = false): Composer
    {
        // for BC reasons, if a config is passed in either as array or a path that is not the default composer.json path
        // we disable local plugins as they really should not be loaded from CWD
        if ($config !== null && $config !== self::getComposerFile() && $disablePlugins === false) {
            $disablePlugins = 'local';
        }

        return self::gcCreate($config, $disablePlugins, null, $disableScripts);
    }

    /**
     * @param string|array<string, mixed>|null $config
     * @param bool|'local'|'global' $disablePlugins
     */
    private static function gcCreate($config, $disablePlugins, ?string $cwd, bool $disableScripts): Composer
    {
        $args = [
            'disablePlugins' => $disablePlugins === true || $disablePlugins === 'local',
            'disableScripts' => $disableScripts,
            'cwd' => $cwd,
        ];
        if (is_array($config)) {
            $args['json'] = \Composer\Json\JsonFile::encode($config);
        } else {
            $args['file'] = $config;
        }
        $desc = \GocomposerPluginHost\Host::cb('factoryCreate', $args);

        return \GocomposerPluginHost\Api::nestedComposer($desc);
    }

    public static function createHttpDownloader($io, Config $config, array $options = []): \Composer\Util\HttpDownloader
    {
        return new \Composer\Util\HttpDownloader($io, $config, $options);
    }

    public static function createGlobal($io, bool $disablePlugins = false, bool $disableScripts = false): ?Composer
    {
        $desc = \GocomposerPluginHost\Host::cb('factoryCreate', ['global' => true, 'disablePlugins' => $disablePlugins, 'disableScripts' => $disableScripts]);

        return $desc === null ? null : \GocomposerPluginHost\Api::nestedComposer($desc);
    }

    /**
     * @param string|array<string, mixed>|null $localConfig
     * @param bool|'local'|'global' $disablePlugins
     */
    public function createComposer($io, $localConfig = null, $disablePlugins = false, ?string $cwd = null, bool $fullLoad = true, bool $disableScripts = false)
    {
        if (!$fullLoad) {
            Unsupported::fail(__METHOD__.' without a full load (PartialComposer)');
        }

        return self::gcCreate($localConfig, $disablePlugins, $cwd, $disableScripts);
    }

    public static function getHomeDir(): string
    {
        return (string) \GocomposerPluginHost\Host::cb('config', ['key' => 'home', 'flags' => 0, 'inst' => 0]);
    }
}
