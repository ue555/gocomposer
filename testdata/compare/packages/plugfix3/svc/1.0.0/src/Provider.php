<?php
namespace Plugfix3\Svc;

use Composer\Plugin\Capability\CommandProvider;

class Provider implements CommandProvider
{
    public function getCommands()
    {
        return array(new SvcCommand());
    }
}
