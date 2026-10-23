<?php declare(strict_types=1);
// gocomposer plugin runtime: ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 src/Composer/IO/ConsoleIO.php
// (MIT license, (c) Nils Adermann, Jordi Boggiano). The terminal program
// status (OSC) reporting was removed; output goes to the gocomposer
// process through GocomposerPluginHost\GoOutput.

namespace Composer\IO;

use Composer\Pcre\Preg;
use Composer\Question\StrictConfirmationQuestion;
use Symfony\Component\Console\Helper\HelperSet;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\ConsoleOutputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\ChoiceQuestion;
use Symfony\Component\Console\Question\Question;

class ConsoleIO extends BaseIO
{
    /** @var InputInterface */
    protected $input;
    /** @var OutputInterface */
    protected $output;
    /** @var HelperSet */
    protected $helperSet;
    /** @var string */
    protected $lastMessage = '';
    /** @var string */
    protected $lastMessageErr = '';

    /** @var float */
    private $startTime;
    /** @var array<IOInterface::*, OutputInterface::VERBOSITY_*> */
    private $verbosityMap;

    public function __construct(InputInterface $input, OutputInterface $output, HelperSet $helperSet)
    {
        $this->input = $input;
        $this->output = $output;
        $this->helperSet = $helperSet;
        $this->verbosityMap = [
            self::QUIET => OutputInterface::VERBOSITY_QUIET,
            self::NORMAL => OutputInterface::VERBOSITY_NORMAL,
            self::VERBOSE => OutputInterface::VERBOSITY_VERBOSE,
            self::VERY_VERBOSE => OutputInterface::VERBOSITY_VERY_VERBOSE,
            self::DEBUG => OutputInterface::VERBOSITY_DEBUG,
        ];
    }

    /**
     * @return void
     */
    public function enableDebugging(float $startTime)
    {
        $this->startTime = $startTime;
    }

    public function isInteractive()
    {
        return $this->input->isInteractive();
    }

    public function isDecorated()
    {
        return $this->output->isDecorated();
    }

    public function isVerbose()
    {
        return $this->output->isVerbose();
    }

    public function isVeryVerbose()
    {
        return $this->output->isVeryVerbose();
    }

    public function isDebug()
    {
        return $this->output->isDebug();
    }

    public function write($messages, bool $newline = true, int $verbosity = self::NORMAL)
    {
        $messages = self::sanitize($messages);

        $this->doWrite($messages, $newline, false, $verbosity);
    }

    public function writeError($messages, bool $newline = true, int $verbosity = self::NORMAL)
    {
        $messages = self::sanitize($messages);

        $this->doWrite($messages, $newline, true, $verbosity);
    }

    public function writeRaw($messages, bool $newline = true, int $verbosity = self::NORMAL)
    {
        $this->doWrite($messages, $newline, false, $verbosity, true);
    }

    public function writeErrorRaw($messages, bool $newline = true, int $verbosity = self::NORMAL)
    {
        $this->doWrite($messages, $newline, true, $verbosity, true);
    }

    /**
     * @param string[]|string $messages
     */
    private function doWrite($messages, bool $newline, bool $stderr, int $verbosity, bool $raw = false): void
    {
        $sfVerbosity = $this->verbosityMap[$verbosity];
        if ($sfVerbosity > $this->output->getVerbosity()) {
            return;
        }

        if ($raw) {
            $sfVerbosity |= OutputInterface::OUTPUT_RAW;
        }

        if (null !== $this->startTime) {
            $memoryUsage = memory_get_usage() / 1024 / 1024;
            $timeSpent = microtime(true) - $this->startTime;
            $messages = array_map(static function ($message) use ($memoryUsage, $timeSpent): string {
                return sprintf('[%.1fMiB/%.2fs] %s', $memoryUsage, $timeSpent, $message);
            }, (array) $messages);
        }

        if (true === $stderr && $this->output instanceof ConsoleOutputInterface) {
            $this->output->getErrorOutput()->write($messages, $newline, $sfVerbosity);
            $this->lastMessageErr = implode($newline ? "\n" : '', (array) $messages);

            return;
        }

        $this->output->write($messages, $newline, $sfVerbosity);
        $this->lastMessage = implode($newline ? "\n" : '', (array) $messages);
    }

    public function overwrite($messages, bool $newline = true, ?int $size = null, int $verbosity = self::NORMAL)
    {
        $this->doOverwrite($messages, $newline, $size, false, $verbosity);
    }

    public function overwriteError($messages, bool $newline = true, ?int $size = null, int $verbosity = self::NORMAL)
    {
        $this->doOverwrite($messages, $newline, $size, true, $verbosity);
    }

    /**
     * @param string[]|string $messages
     */
    private function doOverwrite($messages, bool $newline, ?int $size, bool $stderr, int $verbosity): void
    {
        // messages can be an array, let's convert it to string anyway
        $messages = implode($newline ? "\n" : '', (array) $messages);

        $decorated = $stderr ? $this->getErrorOutput()->isDecorated() : $this->output->isDecorated();

        // backspaces corrupt non-decorated output, so write a plain line instead
        if (!$decorated) {
            if ($messages !== '') {
                $this->doWrite($messages, true, $stderr, $verbosity);
            }
            if ($stderr) {
                $this->lastMessageErr = $messages;
            } else {
                $this->lastMessage = $messages;
            }

            return;
        }

        // since overwrite is supposed to overwrite last message...
        if (!isset($size)) {
            // removing possible formatting of lastMessage with strip_tags
            $size = strlen(strip_tags($stderr ? $this->lastMessageErr : $this->lastMessage));
        }
        // ...let's fill its length with backspaces
        $this->doWrite(str_repeat("\x08", $size), false, $stderr, $verbosity);

        // write the new message
        $this->doWrite($messages, false, $stderr, $verbosity);

        $fill = $size - strlen(strip_tags($messages));
        if ($fill > 0) {
            // whitespace whatever has left
            $this->doWrite(str_repeat(' ', $fill), false, $stderr, $verbosity);
            // move the cursor back
            $this->doWrite(str_repeat("\x08", $fill), false, $stderr, $verbosity);
        }

        if ($newline) {
            $this->doWrite('', true, $stderr, $verbosity);
        }

        if ($stderr) {
            $this->lastMessageErr = $messages;
        } else {
            $this->lastMessage = $messages;
        }
    }

    /**
     * @return ProgressBar
     */
    public function getProgressBar(int $max = 0)
    {
        return new ProgressBar($this->getErrorOutput(), $max);
    }

    public function ask($question, $default = null)
    {
        $question = new Question(self::sanitize($question), is_string($default) ? self::sanitize($default) : $default);

        return $this->askQuestion($question);
    }

    public function askConfirmation($question, $default = true)
    {
        $question = new StrictConfirmationQuestion(self::sanitize($question), is_string($default) ? self::sanitize($default) : $default);

        return $this->askQuestion($question);
    }

    public function askAndValidate($question, $validator, $attempts = null, $default = null)
    {
        $question = new Question(self::sanitize($question), is_string($default) ? self::sanitize($default) : $default);
        $question->setValidator($validator);
        $question->setMaxAttempts($attempts);

        return $this->askQuestion($question);
    }

    public function askAndHideAnswer($question)
    {
        $question = new Question(self::sanitize($question));
        $question->setHidden(true);

        return $this->askQuestion($question);
    }

    public function select($question, $choices, $default, $attempts = false, $errorMessage = 'Value "%s" is invalid', $multiselect = false)
    {
        $question = new ChoiceQuestion(self::sanitize($question), self::sanitize($choices), is_string($default) ? self::sanitize($default) : $default);
        $question->setMaxAttempts($attempts ?: null); // IOInterface requires false, and Question requires null or int
        $question->setErrorMessage($errorMessage);
        $question->setMultiselect($multiselect);

        $result = $this->askQuestion($question);

        $isAssoc = (bool) \count(array_filter(array_keys($choices), 'is_string'));
        if ($isAssoc) {
            return $result;
        }

        if (!is_array($result)) {
            return (string) array_search($result, $choices, true);
        }

        $results = [];
        foreach ($choices as $index => $choice) {
            if (in_array($choice, $result, true)) {
                $results[] = (string) $index;
            }
        }

        return $results;
    }

    /**
     * @return mixed
     */
    private function askQuestion(Question $question)
    {
        /** @var \Symfony\Component\Console\Helper\QuestionHelper $helper */
        $helper = $this->helperSet->get('question');

        return $helper->ask($this->input, $this->getErrorOutput(), $question);
    }

    public function getTable(): Table
    {
        return new Table($this->output);
    }

    private function getErrorOutput(): OutputInterface
    {
        if ($this->output instanceof ConsoleOutputInterface) {
            return $this->output->getErrorOutput();
        }

        return $this->output;
    }

    /**
     * Sanitize string to remove control characters
     *
     * @param string|iterable<string> $messages
     * @return string|array<string>
     */
    public static function sanitize($messages, bool $allowNewlines = true)
    {
        $escapePattern = '\x1B\[[\x30-\x3F]*[\x20-\x2F]*[\x40-\x7E]|\x1B\].*?(?:\x1B\\\\|\x07)|\x1B.';
        $pattern = $allowNewlines ? "{{$escapePattern}|[\x01-\x09\x0B\x0C\x0E-\x1A]|\r(?!\n)}u" : "{{$escapePattern}|[\x01-\x1A]}u";
        if (is_string($messages)) {
            $messages = self::ensureValidUtf8($messages);

            return Preg::replace($pattern, '', $messages);
        }

        $sanitized = [];
        foreach ($messages as $key => $message) {
            $message = self::ensureValidUtf8($message);
            $sanitized[$key] = Preg::replace($pattern, '', $message);
        }

        return $sanitized;
    }

    private static function ensureValidUtf8(string $string): string
    {
        if (function_exists('mb_check_encoding') && mb_check_encoding($string, 'UTF-8')) {
            return $string;
        }

        if (function_exists('mb_convert_encoding')) {
            return (string) mb_convert_encoding($string, 'UTF-8', 'UTF-8');
        }

        if (function_exists('iconv')) {
            $cleaned = @iconv('UTF-8', 'UTF-8//TRANSLIT', $string);
            if ($cleaned !== false) {
                return $cleaned;
            }
        }

        return $string;
    }
}
