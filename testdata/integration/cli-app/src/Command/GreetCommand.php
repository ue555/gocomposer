<?php

namespace App\Command;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(name: 'greet', description: 'Greets someone')]
final class GreetCommand extends Command
{
    protected function configure(): void
    {
        $this->addOption('name', null, InputOption::VALUE_REQUIRED, 'Who to greet', 'world');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln(app_greeting((string) $input->getOption('name')) . ' (' . \LegacyGreeter::suffix() . ')');

        return Command::SUCCESS;
    }
}
