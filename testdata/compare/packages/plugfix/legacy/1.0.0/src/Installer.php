<?php
namespace Plugfix\Legacy;

use Composer\Installer\LibraryInstaller;

class Installer extends LibraryInstaller
{
    public function supports($packageType)
    {
        return $packageType === 'plugfix-legacy-thing';
    }
}
