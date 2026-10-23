<?php
namespace Plugfix\Api;

use Composer\Plugin\Capability\CommandProvider;

class Provider implements CommandProvider
{
    public function getCommands()
    {
        return array(new ApiCommand());
    }
}
