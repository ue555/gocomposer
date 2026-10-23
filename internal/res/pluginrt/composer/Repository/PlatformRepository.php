<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Repository\PlatformRepository with the
// static helpers of Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/PlatformRepository.php (MIT license, (c) Nils
// Adermann, Jordi Boggiano). The platform packages themselves are detected by
// gocomposer and are not available as a repository object to plugins.

namespace Composer\Repository;

use Composer\Pcre\Preg;
use GocomposerPluginHost\Unsupported;

class PlatformRepository extends ArrayRepository
{
    public const PLATFORM_PACKAGE_REGEX = '{^(?:php(?:-64bit|-ipv6|-zts|-debug)?|hhvm|(?:ext|lib)-[a-z0-9](?:[_.-]?[a-z0-9]+)*|composer(?:-(?:plugin|runtime)-api)?)$}iD';

    public function __construct(array $packages = [], array $overrides = [], $runtime = null, $hhvmDetector = null)
    {
        Unsupported::fail(__METHOD__);
    }

    public function getRepoName(): string
    {
        return 'platform repo';
    }

    public static function isPlatformPackage(string $name): bool
    {
        static $cache = [];

        if (isset($cache[$name])) {
            return $cache[$name];
        }

        return $cache[$name] = Preg::isMatch(PlatformRepository::PLATFORM_PACKAGE_REGEX, $name);
    }

    public static function getPlatformPhpVersion(): ?string
    {
        return null;
    }
}
