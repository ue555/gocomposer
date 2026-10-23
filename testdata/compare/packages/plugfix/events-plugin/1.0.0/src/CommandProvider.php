<?php
namespace Plugfix\Events;

class CommandProvider implements \Composer\Plugin\Capability\CommandProvider
{
    public function getCommands()
    {
        return array(new HelloCommand());
    }
}
