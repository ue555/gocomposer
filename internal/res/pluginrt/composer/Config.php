<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Config for plugins. The values are
// resolved by gocomposer's Go implementation of Composer\Config
// (f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/Config.php, MIT
// license, (c) Nils Adermann, Jordi Boggiano): the configuration of a
// Composer instance, or one created by plugin code (new Config(),
// Factory::createConfig()), which gocomposer keeps for the object.
// merge() merges into it. The config sources are Composer's
// JsonConfigSource objects for the files gocomposer loaded (composer.json,
// the home config.json and auth.json, a local auth.json).
// prohibitUrlByConfig() is Composer's code (the insecure protocol warning
// of a Composer instance's configuration is shared with gocomposer's HTTP
// client, like Composer's Config object is shared with its HttpDownloader).

namespace Composer;

use Composer\Config\ConfigSourceInterface;
use Composer\Config\JsonConfigSource;
use Composer\Downloader\TransportException;
use Composer\IO\IOInterface;
use Composer\Json\JsonFile;
use Composer\Pcre\Preg;
use Composer\Util\Url;
use GocomposerPluginHost\Host;

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
    /** @var int|null gocomposer's configuration of a Config created by plugin code */
    private $gcCfg = null;
    /** @var bool */
    private $gcUseEnvironment;
    /** @var string|null */
    private $gcBaseDir;
    /** @var bool whether the config sources of gocomposer's files were set */
    private $gcSourcesLoaded = false;

    /** @var ?ConfigSourceInterface */
    private $configSource;
    /** @var ?ConfigSourceInterface */
    private $authConfigSource;
    /** @var ?ConfigSourceInterface */
    private $localAuthConfigSource = null;
    /** @var array<string, true> */
    private $warnedHosts = [];
    /** @var array<string, true> */
    private $sslVerifyWarnedHosts = [];

    public function __construct(bool $useEnvironment = true, ?string $baseDir = null)
    {
        $this->gcUseEnvironment = $useEnvironment;
        $this->gcBaseDir = $baseDir;
    }

    /**
     * @internal Factory::createConfig(): the configuration of the home
     * directory (config.json, auth.json, environment) loaded by gocomposer
     */
    public static function gcCreate(?string $cwd): self
    {
        $config = new self();
        $r = Host::cb('configCreate', ['cwd' => $cwd]);
        $config->gcCfg = (int) $r['cfg'];
        $config->gcSources($r);

        return $config;
    }

    /**
     * @return array<string, mixed>
     */
    private function gcArgs(): array
    {
        if (!$this->gcDetached) {
            return ['inst' => $this->gcInst];
        }
        if ($this->gcCfg === null) {
            $this->gcCfg = (int) Host::cb('configNew', ['useEnv' => $this->gcUseEnvironment, 'baseDir' => $this->gcBaseDir]);
            $this->gcSourcesLoaded = true;
        }

        return ['cfg' => $this->gcCfg];
    }

    /**
     * @param array<string, mixed> $r
     */
    private function gcSources(array $r): void
    {
        $this->gcSourcesLoaded = true;
        if (isset($r['config']) && $this->configSource === null) {
            $this->configSource = new JsonConfigSource(new JsonFile($r['config'], null, Host::$io));
        }
        if (isset($r['auth']) && $this->authConfigSource === null) {
            $this->authConfigSource = new JsonConfigSource(new JsonFile($r['auth'], null, Host::$io), true);
        }
        if (isset($r['localAuth']) && $this->localAuthConfigSource === null) {
            $this->localAuthConfigSource = new JsonConfigSource(new JsonFile($r['localAuth'], null, Host::$io), true);
        }
    }

    private function gcLoadSources(): void
    {
        if ($this->gcSourcesLoaded || $this->gcDetached) {
            return;
        }
        $this->gcSources(Host::cb('configSources', ['inst' => $this->gcInst]));
    }

    public function setBaseDir(?string $baseDir): void
    {
        $this->gcBaseDir = $baseDir;
        Host::cb('configSetBaseDir', ['baseDir' => $baseDir] + $this->gcArgs());
    }

    public function setConfigSource(ConfigSourceInterface $source): void
    {
        $this->gcLoadSources();
        $this->configSource = $source;
    }

    public function getConfigSource(): ConfigSourceInterface
    {
        $this->gcLoadSources();

        return $this->configSource;
    }

    public function setAuthConfigSource(ConfigSourceInterface $source): void
    {
        $this->gcLoadSources();
        $this->authConfigSource = $source;
    }

    public function getAuthConfigSource(): ConfigSourceInterface
    {
        $this->gcLoadSources();

        return $this->authConfigSource;
    }

    public function setLocalAuthConfigSource(ConfigSourceInterface $source): void
    {
        $this->gcLoadSources();
        $this->localAuthConfigSource = $source;
    }

    public function getLocalAuthConfigSource(): ?ConfigSourceInterface
    {
        $this->gcLoadSources();

        return $this->localAuthConfigSource;
    }

    /**
     * @return mixed
     */
    public function get(string $key, int $flags = 0)
    {
        return Host::cb('config', ['key' => $key, 'flags' => $flags] + $this->gcArgs());
    }

    public function has(string $key): bool
    {
        return (bool) Host::cb('configHas', ['key' => $key] + $this->gcArgs());
    }

    /**
     * @return array<string, mixed[]>
     */
    public function all(int $flags = 0): array
    {
        return Host::cb('configAll', ['flags' => $flags] + $this->gcArgs());
    }

    /**
     * @return array<string, mixed[]>
     */
    public function raw(): array
    {
        return Host::cb('configRaw', $this->gcArgs());
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
        return (string) Host::cb('configSource', ['key' => $key] + $this->gcArgs());
    }

    /**
     * @param array<string, mixed> $config
     */
    public function merge(array $config, string $source = self::SOURCE_UNKNOWN): void
    {
        // merged into gocomposer's configuration (Config::merge semantics,
        // including repositories)
        Host::cb('configMerge', ['config' => $config, 'source' => $source] + $this->gcArgs());
    }

    /**
     * Validates that the passed URL is allowed to be used by current config, or throws an exception.
     *
     * @param IOInterface $io
     * @param mixed[]     $repoOptions
     */
    public function prohibitUrlByConfig(string $url, ?IOInterface $io = null, array $repoOptions = []): void
    {
        // Config::prohibitUrlByConfig of Composer f83588392df43d4fd027d84d26de459946ef4e80
        // (MIT license, (c) Nils Adermann, Jordi Boggiano)

        // git's <transport>::<address> syntax connects to the address, so check that
        $checkedUrl = $url;
        while (Preg::isMatchStrictGroups('{^[a-z][a-z0-9+.-]*::(.+)$}i', $checkedUrl, $match)) {
            $checkedUrl = $match[1];
        }

        // Return right away for custom URLs without a scheme like local paths (see issue #5173), but do not rely
        // on filter_var for this as it rejects many malformed URLs which tools would still happily connect to
        if (!Preg::isMatchStrictGroups('{^([a-z][a-z0-9+.-]*)://}i', $checkedUrl, $match)) {
            return;
        }

        // Throw exception on known insecure protocols
        $scheme = strtolower($match[1]);
        $hostname = parse_url($checkedUrl, PHP_URL_HOST);
        if (in_array($scheme, ['http', 'git', 'ftp', 'svn'])) {
            if ($this->get('secure-http')) {
                if ($scheme === 'svn') {
                    if (in_array($hostname, $this->get('secure-svn-domains'), true)) {
                        return;
                    }

                    throw new TransportException("Your configuration does not allow connections to " . Url::sanitize($url) . ". See https://getcomposer.org/doc/06-config.md#secure-svn-domains for details.");
                }

                throw new TransportException("Your configuration does not allow connections to " . Url::sanitize($url) . ". See https://getcomposer.org/doc/06-config.md#secure-http for details.");
            }
            if ($io !== null) {
                if (is_string($hostname)) {
                    if (!$this->gcDetached) {
                        // the warned hosts of gocomposer's HTTP client
                        Host::cb('configWarnInsecure', ['url' => $checkedUrl, 'inst' => $this->gcInst]);
                    } else {
                        if (!isset($this->warnedHosts[$hostname])) {
                            $io->writeError("<warning>Warning: Accessing $hostname over $scheme which is an insecure protocol.</warning>");
                        }
                        $this->warnedHosts[$hostname] = true;
                    }
                }
            }
        }

        if ($io !== null && is_string($hostname) && !isset($this->sslVerifyWarnedHosts[$hostname])) {
            $warning = null;
            if (isset($repoOptions['ssl']['verify_peer']) && !(bool) $repoOptions['ssl']['verify_peer']) {
                $warning = 'verify_peer';
            }

            if (isset($repoOptions['ssl']['verify_peer_name']) && !(bool) $repoOptions['ssl']['verify_peer_name']) {
                $warning = $warning === null ? 'verify_peer_name' : $warning . ' and verify_peer_name';
            }

            if ($warning !== null) {
                if (!$this->gcDetached) {
                    // the hosts warned about by gocomposer's HTTP requests
                    Host::cb('configWarnSsl', ['host' => $hostname, 'option' => $warning === 'verify_peer and verify_peer_name' ? 'both' : $warning, 'inst' => $this->gcInst]);
                } else {
                    $io->writeError("<warning>Warning: Accessing $hostname with $warning disabled.</warning>");
                    $this->sslVerifyWarnedHosts[$hostname] = true;
                }
            }
        }
    }

    public static function disableProcessTimeout(): void
    {
        \Composer\Util\ProcessExecutor::setTimeout(0);
    }
}
