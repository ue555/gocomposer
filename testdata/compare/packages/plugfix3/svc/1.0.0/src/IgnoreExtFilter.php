<?php
namespace Plugfix3\Svc;

use Composer\Filter\PlatformRequirementFilter\PlatformRequirementFilterInterface;

/**
 * A platform requirement filter of its own: ignores the ext-plugfix3* and
 * ext-json requirements.
 */
class IgnoreExtFilter implements PlatformRequirementFilterInterface
{
    public function isIgnored(string $req): bool
    {
        return strpos($req, 'ext-plugfix3') === 0 || $req === 'ext-json';
    }

    public function isUpperBoundIgnored(string $req): bool
    {
        return $this->isIgnored($req);
    }
}
