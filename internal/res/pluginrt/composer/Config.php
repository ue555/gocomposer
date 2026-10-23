<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Config for plugins. The values are
// resolved by gocomposer's Go implementation of Composer\Config
// (f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Config.php, MIT
// license, (c) Nils Adermann, Jordi Boggiano), merge() merges into it; the
// config sources (writing composer.json/config.json) are not available.

namespace Composer;

use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class Config
{
    public const SOURCE_DEFAULT = 'default';
    public const SOURCE_COMMAND = 'command';
    public const SOURCE_UNKNOWN = 'unknown';

    public const RELATIVE_PATHS = 1;

    /** @var array<string, mixed> */
    public static $defaultConfig = [];

    /** @var array<string, mixed> */
    public static $defaultRepositories = [];

    /** @var int gocomposer instance id */
    public $gcInst = 0;
    /** @var bool created by plugin code (not one of gocomposer's instances) */
    public $gcDetached = true;

    public function __construct(bool $useEnvironment = true, ?string $baseDir = null)
    {
    }

    /**
     * @return mixed
     */
    public function get(string $key, int $flags = 0)
    {
        return Host::cb('config', ['key' => $key, 'flags' => $flags, 'inst' => $this->gcInst]);
    }

    public function has(string $key): bool
    {
        return (bool) Host::cb('configHas', ['key' => $key, 'inst' => $this->gcInst]);
    }

    /**
     * @return array<string, mixed[]>
     */
    public function all(int $flags = 0): array
    {
        return Host::cb('configAll', ['flags' => $flags, 'inst' => $this->gcInst]);
    }

    /**
     * @return array<string, mixed[]>
     */
    public function raw(): array
    {
        return Host::cb('configRaw', ['inst' => $this->gcInst]);
    }

    /**
     * @return array<int|string, mixed>
     */
    public function getRepositories(): array
    {
        $raw = $this->raw();

        return $raw['repositories'] ?? [];
    }

    public function getSourceOfValue(string $key): string
    {
        return (string) Host::cb('configSource', ['key' => $key, 'inst' => $this->gcInst]);
    }

    public function setConfigSource($source): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function getConfigSource()
    {
        Unsupported::fail(__METHOD__);
    }

    public function setAuthConfigSource($source): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function getAuthConfigSource()
    {
        Unsupported::fail(__METHOD__);
    }

    public function setLocalAuthConfigSource($source): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function getLocalAuthConfigSource()
    {
        Unsupported::fail(__METHOD__);
    }

    /**
     * @param array<string, mixed> $config
     */
    public function merge(array $config, string $source = self::SOURCE_UNKNOWN): void
    {
        // merged into gocomposer's configuration (Config::merge semantics,
        // including repositories)
        if ($this->gcDetached) {
            Unsupported::fail(__METHOD__.' on a Config created by plugin code');
        }
        Host::cb('configMerge', ['config' => $config, 'source' => $source, 'inst' => $this->gcInst]);
    }

    public function prohibitUrlByConfig(string $url, $io = null, array $repoOptions = []): void
    {
        Unsupported::fail(__METHOD__);
    }

    public static function disableProcessTimeout(): void
    {
        \Composer\Util\ProcessExecutor::setTimeout(0);
    }
}
