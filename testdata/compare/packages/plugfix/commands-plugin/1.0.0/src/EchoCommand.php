<?php
namespace Plugfix\Commands;

use Composer\Command\BaseCommand;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

class EchoCommand extends BaseCommand
{
    protected function configure()
    {
        $this->setName('plugfix:echo')
            ->setAliases(array('pecho'))
            ->setDescription('Echoes its input')
            ->addArgument('words', InputArgument::IS_ARRAY, 'Words')
            ->addOption('tag', null, InputOption::VALUE_REQUIRED | InputOption::VALUE_IS_ARRAY, 'Tags')
            ->addOption('fail', null, InputOption::VALUE_NONE, 'Throw an exception');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('words='.json_encode($input->getArgument('words')).' tags='.json_encode($input->getOption('tag')).' verbose='.var_export($output->isVerbose(), true).' name='.$this->getName());
        $output->writeln('<comment>composer '.$this->getApplication()->getName().'</comment>');
        if ($input->getOption('fail')) {
            throw new \RuntimeException('plugfix echo failed');
        }

        return 0;
    }
}
