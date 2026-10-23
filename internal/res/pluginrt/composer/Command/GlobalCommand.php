<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Command\GlobalCommand, standing in for
// Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Command/GlobalCommand.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The global command is run by gocomposer; an instance is
// on the call stack of plugin code while it runs (plugins like symfony/flex
// look for it with debug_backtrace()).

namespace Composer\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class GlobalCommand extends BaseCommand
{
    protected function configure(): void
    {
        $this->setName('global');
    }

    /**
     * @internal
     * @return mixed
     */
    public function gocomposerFrame(callable $fn)
    {
        return $fn();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        \GocomposerPluginHost\Unsupported::fail(__METHOD__);
    }
}
