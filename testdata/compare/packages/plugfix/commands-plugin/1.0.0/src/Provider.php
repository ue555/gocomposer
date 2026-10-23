<?php
namespace Plugfix\Commands;

use Composer\Plugin\Capability\CommandProvider;

class Provider implements CommandProvider
{
    private $args;

    public function __construct(array $args)
    {
        $this->args = $args;
    }

    public function getCommands()
    {
        return array(new ShowCommand(), new EchoCommand());
    }
}
