<?php
namespace Plugfix\Commands;

use Composer\Command\BaseCommand;

class ShowCommand extends BaseCommand
{
    protected function configure()
    {
        $this->setName('show')->setDescription('Conflicts with a Composer command');
    }
}
