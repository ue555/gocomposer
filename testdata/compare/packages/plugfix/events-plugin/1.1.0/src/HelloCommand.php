<?php
namespace Plugfix\Events;

use Composer\Command\BaseCommand;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

class HelloCommand extends BaseCommand
{
    protected function configure()
    {
        $this->setName('plugfix:hello')
            ->setDescription('Says hello from the fixture plugin')
            ->addArgument('who', InputArgument::OPTIONAL, 'Who to greet', 'world')
            ->addOption('shout', null, InputOption::VALUE_NONE, 'Greet loudly')
            ->addOption('times', 't', InputOption::VALUE_REQUIRED, 'Repeat count', '1');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $msg = 'Hello ' . $input->getArgument('who');
        if ($input->getOption('shout')) {
            $msg = strtoupper($msg) . '!';
        }
        for ($i = 0; $i < (int) $input->getOption('times'); $i++) {
            $output->writeln('<info>' . $msg . '</info>');
        }
        $composer = $this->requireComposer();
        $this->getIO()->writeError('root=' . $composer->getPackage()->getName() . ' plugins=' . count($composer->getPluginManager()->getPlugins()));
        $this->getIO()->write('verbose line', true, \Composer\IO\IOInterface::VERBOSE);
        return (int) $input->getOption('times') > 2 ? 3 : 0;
    }
}
