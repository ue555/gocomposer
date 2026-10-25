package composer

// Scripts pointing to a Symfony Console command class: the isCommandClass
// branch of Composer\EventDispatcher\EventDispatcher::doDispatch (composer/
// composer f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
// Jordi Boggiano). The class is run with the configured PHP in a new
// Symfony Application (single command mode, exceptions not caught) on an
// output with Composer's verbosity, decoration and extra styles, using the
// symfony/console Composer bundles (internal/res/pluginrt/vendor) and the
// project's autoloader.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/res"
)

// scriptEventConstants are the Composer\Script\ScriptEvents values.
var scriptEventConstants = []string{
	"pre-install-cmd", "post-install-cmd", "pre-update-cmd", "post-update-cmd",
	"pre-status-cmd", "post-status-cmd", "pre-autoload-dump", "post-autoload-dump",
	"post-root-package-install", "post-create-project-cmd", "pre-archive-cmd", "post-archive-cmd",
}

// isScriptEventsConstant is defined('Composer\Script\ScriptEvents::'.str_replace('-', '_', strtoupper($name))).
func isScriptEventsConstant(name string) bool {
	c := strings.ReplaceAll(strings.ToUpper(name), "-", "_")
	for _, e := range scriptEventConstants {
		if strings.ReplaceAll(strings.ToUpper(e), "-", "_") == c {
			return true
		}
	}
	return false
}

// runCommandClass runs the Command class className for the event.
func (d *EventDispatcher) runCommandClass(ev *ScriptEvent, className string, additional []string) (int, error) {
	input := ""
	if v, ok := ev.Flags["script-alias-input"]; ok && v != nil {
		input, _ = v.(string)
	} else {
		quoted := make([]string, len(additional))
		for i, a := range additional {
			quoted[i] = escapeArg(a)
		}
		input = strings.Join(quoted, " ")
	}
	data := phpjson.NewMap()
	data.Set("mode", "event")
	data.Set("class", className)
	data.Set("name", ev.Name)
	data.Set("reserved", isScriptEventsConstant(ev.Name))
	data.Set("input", input)
	res, err := runCommandRunner(d.c, d.Binary, ev.DevMode, data)
	if err != nil {
		return 0, err
	}
	switch res.Str("status") {
	case "not-autoloadable":
		d.c.IO.WriteErrorV("<warning>Class "+className+" is not autoloadable, can not call "+ev.Name+" script</warning>", true, 16)
		return 0, nil
	case "not-command":
		d.c.IO.WriteErrorV("<warning>Class "+className+` does not extend Symfony\Component\Console\Command\Command, can not call `+ev.Name+" script</warning>", true, 16)
		return 0, nil
	case "reserved":
		d.c.IO.WriteErrorV("<warning>You cannot bind "+ev.Name+" to a Command class, use a non-reserved name</warning>", true, 16)
		return 0, nil
	case "exception", "error":
		if res.Str("status") == "exception" {
			d.c.IO.WriteErrorV(fmt.Sprintf("<error>Script %s handling the %s event terminated with an exception</error>", className, ev.Name), true, 16)
		}
		return 0, &CallbackError{Class: res.Str("class"), Msg: res.Str("message")}
	}
	return resultCode(res), nil
}

func resultCode(res *phpjson.Map) int {
	code, _ := res.Get("code")
	switch c := code.(type) {
	case int64:
		return int(c)
	case float64:
		return int(c)
	}
	return 0
}

// ScriptCommandClass is a composer.json script whose value names a Symfony
// Console command class, registered by Application::doRun as the command
// itself (with the warnings Composer writes while doing so).
type ScriptCommandClass struct {
	Script   string
	Command  PluginCommand
	Warnings []string
}

// DescribeScriptCommandClasses checks which of the scripts (name => class)
// are autoloadable Symfony Console commands (class_exists, is_subclass_of)
// and describes them after `new $class($script)`, the description fallback
// and the name check of Application::doRun.
func DescribeScriptCommandClasses(c *Composer, scripts []string, classes map[string]string, descriptions map[string]string) (map[string]ScriptCommandClass, error) {
	list := []phpjson.Value{}
	for _, s := range scripts {
		m := phpjson.NewMap()
		m.Set("script", s)
		m.Set("class", classes[s])
		m.Set("description", descriptions[s])
		list = append(list, m)
	}
	data := phpjson.NewMap()
	data.Set("mode", "describe")
	data.Set("scripts", list)
	res, err := runCommandRunner(c, "", true, data)
	if err != nil {
		return nil, err
	}
	out := map[string]ScriptCommandClass{}
	cmds, _ := res.Get("commands")
	phpjson.Each(cmds, func(_ string, e phpjson.Value) {
		m := phpjson.ToMap(e)
		if m == nil {
			return
		}
		sc := ScriptCommandClass{Script: m.Str("script")}
		phpjson.Each(mustGet(m, "warnings"), func(_ string, w phpjson.Value) {
			if s, ok := w.(string); ok {
				sc.Warnings = append(sc.Warnings, s)
			}
		})
		if dm := phpjson.ToMap(mustGet(m, "command")); dm != nil {
			sc.Command = decodeDescribedCommand(dm)
		}
		out[sc.Script] = sc
	})
	return out, nil
}

func decodeDescribedCommand(m *phpjson.Map) PluginCommand {
	pc := PluginCommand{Name: m.Str("name"), Class: m.Str("class"), Description: m.Str("description"), Help: m.Str("help"), Synopsis: m.Str("synopsis"), Hidden: mustGet(m, "hidden") == true}
	phpjson.Each(mustGet(m, "aliases"), func(_ string, a phpjson.Value) {
		if s, ok := a.(string); ok {
			pc.Aliases = append(pc.Aliases, s)
		}
	})
	phpjson.Each(mustGet(m, "usages"), func(_ string, a phpjson.Value) {
		if s, ok := a.(string); ok {
			pc.Usages = append(pc.Usages, s)
		}
	})
	phpjson.Each(mustGet(m, "args"), func(_ string, a phpjson.Value) {
		am := phpjson.ToMap(a)
		if am == nil {
			return
		}
		pc.Args = append(pc.Args, PluginCommandArg{Name: am.Str("name"), Required: mustGet(am, "required") == true, Array: mustGet(am, "array") == true, Desc: am.Str("desc"), Default: mustGet(am, "default")})
	})
	phpjson.Each(mustGet(m, "opts"), func(_ string, a phpjson.Value) {
		om := phpjson.ToMap(a)
		if om == nil {
			return
		}
		short, _ := mustGet(om, "shortcut").(string)
		pc.Opts = append(pc.Opts, PluginCommandOpt{Name: om.Str("name"), Shortcut: short, None: mustGet(om, "none") == true, Required: mustGet(om, "required") == true,
			Optional: mustGet(om, "optional") == true, Array: mustGet(om, "array") == true, Negatable: mustGet(om, "negatable") == true, Desc: om.Str("desc"), Default: mustGet(om, "default")})
	})
	return pc
}

// RunScriptCommandClass runs a script registered as its command class
// (`composer <script> ...`): the command runs in the Application with
// Composer's global options; exceptions are rendered by the caller.
func RunScriptCommandClass(c *Composer, script, className, input string) (int, error) {
	data := phpjson.NewMap()
	data.Set("mode", "direct")
	data.Set("class", className)
	data.Set("name", script)
	data.Set("input", input)
	res, err := runCommandRunner(c, "", true, data)
	if err != nil {
		return 1, err
	}
	switch res.Str("status") {
	case "exception", "error":
		code := resultCode(res)
		if code <= 0 {
			code = 1
		}
		return min(code, 255), &CallbackError{Class: res.Str("class"), Msg: res.Str("message")}
	}
	return resultCode(res), nil
}

// runCommandRunner runs command-runner.php with the project autoloader and
// returns its result document.
func runCommandRunner(c *Composer, binary string, devMode bool, data *phpjson.Map) (*phpjson.Map, error) {
	php, _, err := c.PHP()
	if err != nil || php == "" {
		return nil, &UserError{Msg: "Running Symfony Console command classes from scripts requires a PHP binary (set GOCOMPOSER_PHP or add php to PATH)", Code: 1}
	}
	tmpDir, err := os.MkdirTemp("", "gocomposer-command-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)
	if err := writePluginRuntime(tmpDir); err != nil {
		return nil, err
	}
	loaderFile := filepath.Join(tmpDir, "project-loader.php")
	d := &EventDispatcher{c: c}
	if err := d.writeScriptLoader(loaderFile, devMode); err != nil {
		return nil, err
	}
	// the generated loader requires ClassLoader.php next to itself
	if err := os.WriteFile(filepath.Join(tmpDir, "ClassLoader.php"), res.ClassLoader, 0o600); err != nil {
		return nil, err
	}
	resultFile := filepath.Join(tmpDir, "result.json")
	data.Set("loader", loaderFile)
	data.Set("verbosity", int64(c.IO.Verbosity))
	data.Set("decorated", c.IO.Decorated)
	data.Set("result", resultFile)
	enc, _ := phpjson.Encode(data, 0)
	runner := filepath.Join(tmpDir, "command-runner.php")
	if err := os.WriteFile(runner, []byte(commandRunnerPHP), 0o600); err != nil {
		return nil, err
	}
	cmd := exec.Command(php, "-d", "memory_limit=1536M", runner)
	cmd.Dir = c.Cwd
	cmd.Env = append(os.Environ(), "GOCOMPOSER_SCRIPT_COMMAND="+enc)
	if binary != "" {
		cmd.Env = append(cmd.Env, "COMPOSER_BINARY="+binary)
	}
	cmd.Stdout = c.IO.Out
	cmd.Stderr = c.IO.Err
	if c.IO.Interactive {
		cmd.Stdin = os.Stdin
	}
	runErr := cmd.Run()
	raw, rerr := os.ReadFile(resultFile)
	if rerr != nil {
		code := 1
		if ee, ok := runErr.(*exec.ExitError); ok {
			code = ee.ExitCode()
		}
		return nil, &CallbackError{Msg: fmt.Sprintf("The PHP command runner failed (exit code %d)", code)}
	}
	v, err := phpjson.Decode(raw)
	if err != nil {
		return nil, err
	}
	m := phpjson.AsMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	return m, nil
}

const commandRunnerPHP = `<?php
// gocomposer: runs a script that points to a Symfony Console command class
// like Composer's EventDispatcher (isCommandClass branch).
$data = json_decode((string) getenv('GOCOMPOSER_SCRIPT_COMMAND'), true);
$result = static function (array $r) use ($data): void {
    file_put_contents($data['result'], json_encode($r));
};
$rt = __DIR__;
// Composer\Util\ErrorHandler::handle of Composer f83588392df43d4fd027d84d26de459946ef4e80
// (MIT license, (c) Nils Adermann, Jordi Boggiano). Deprecations raised in
// the bundled copies of Composer's libraries (this runtime directory) are
// PHP version artifacts of Composer itself and are not repeated here.
$deprecationShown = 0;
set_error_handler(static function (int $level, string $message, string $file, int $line) use ($rt, $data, &$deprecationShown): bool {
    $isDeprecationNotice = $level === E_DEPRECATED || $level === E_USER_DEPRECATED;
    if (!$isDeprecationNotice && 0 === (error_reporting() & $level)) {
        return true;
    }
    if (!$isDeprecationNotice) {
        throw new \ErrorException($message, 0, $level, $file, $line);
    }
    if (strpos($file, $rt.'/') === 0) {
        return true;
    }
    $verbose = (int) $data['verbosity'] >= 64;
    if ($deprecationShown > 0 && !$verbose) {
        if ($deprecationShown === 1) {
            fwrite(STDERR, (!empty($data['decorated']) ? "\033[30;43m" : '').'More deprecation notices were hidden, run again with '.chr(96).'-v'.chr(96).' to show them.'.(!empty($data['decorated']) ? "\033[39;49m" : '').PHP_EOL);
            $deprecationShown = 2;
        }

        return true;
    }
    $deprecationShown = 1;
    fwrite(STDERR, (!empty($data['decorated']) ? "\033[30;43m" : '').'Deprecation Notice: '.$message.' in '.$file.':'.$line.(!empty($data['decorated']) ? "\033[39;49m" : '').PHP_EOL);

    return true;
});
error_reporting(E_ALL);
$prefixes = [
    'Symfony\\Component\\Console\\' => '/vendor/symfony/console/',
    'Symfony\\Component\\String\\' => '/vendor/symfony/string/',
    'Symfony\\Contracts\\Service\\' => '/vendor/symfony/service-contracts/',
    'Symfony\\Polyfill\\Ctype\\' => '/vendor/symfony/polyfill-ctype/',
    'Symfony\\Polyfill\\Intl\\Grapheme\\' => '/vendor/symfony/polyfill-intl-grapheme/',
    'Symfony\\Polyfill\\Intl\\Normalizer\\' => '/vendor/symfony/polyfill-intl-normalizer/',
    'Symfony\\Polyfill\\Mbstring\\' => '/vendor/symfony/polyfill-mbstring/',
    'Symfony\\Polyfill\\Php73\\' => '/vendor/symfony/polyfill-php73/',
    'Symfony\\Polyfill\\Php80\\' => '/vendor/symfony/polyfill-php80/',
    'Symfony\\Polyfill\\Php81\\' => '/vendor/symfony/polyfill-php81/',
    'Symfony\\Polyfill\\Php84\\' => '/vendor/symfony/polyfill-php84/',
];
spl_autoload_register(static function (string $class) use ($rt, $prefixes): void {
    foreach ($prefixes as $prefix => $dir) {
        if (strpos($class, $prefix) === 0) {
            $file = $rt.$dir.str_replace('\\', '/', substr($class, strlen($prefix))).'.php';
            if (is_file($file)) {
                require $file;
            }

            return;
        }
    }
}, true, true);
foreach ([
    '/vendor/symfony/deprecation-contracts/function.php',
    '/vendor/symfony/polyfill-ctype/bootstrap.php',
    '/vendor/symfony/polyfill-intl-grapheme/bootstrap.php',
    '/vendor/symfony/polyfill-intl-normalizer/bootstrap.php',
    '/vendor/symfony/polyfill-mbstring/bootstrap.php',
    '/vendor/symfony/polyfill-php73/bootstrap.php',
    '/vendor/symfony/polyfill-php80/bootstrap.php',
    '/vendor/symfony/polyfill-php81/bootstrap.php',
    '/vendor/symfony/polyfill-php84/bootstrap.php',
    '/vendor/symfony/string/Resources/functions.php',
] as $file) {
    if (is_file($rt.$file)) {
        require_once $rt.$file;
    }
}
// the project autoloader is appended after Composer's own classes
require $data['loader'];

$Command = 'Symfony\\Component\\Console\\Command\\Command';
$describe = static function ($command): array {
    $def = $command->getDefinition();
    $args = [];
    foreach ($def->getArguments() as $arg) {
        $args[] = ['name' => $arg->getName(), 'required' => $arg->isRequired(), 'array' => $arg->isArray(), 'desc' => $arg->getDescription(), 'default' => $arg->getDefault()];
    }
    $opts = [];
    foreach ($def->getOptions() as $opt) {
        $opts[] = ['name' => $opt->getName(), 'shortcut' => $opt->getShortcut(), 'none' => !$opt->acceptValue(), 'required' => $opt->isValueRequired(), 'optional' => $opt->isValueOptional(), 'array' => $opt->isArray(), 'negatable' => $opt->isNegatable(), 'desc' => $opt->getDescription(), 'default' => $opt->getDefault()];
    }

    return ['name' => $command->getName(), 'class' => get_class($command), 'description' => $command->getDescription(), 'help' => $command->getHelp(),
        'aliases' => $command->getAliases(), 'hidden' => $command->isHidden(), 'synopsis' => $command->getSynopsis(false), 'usages' => $command->getUsages(), 'args' => $args, 'opts' => $opts];
};

if ($data['mode'] === 'describe') {
    // Application::doRun: scripts pointing to a valid Command subclass are
    // imported as commands
    $out = [];
    foreach ($data['scripts'] as $s) {
        $dummy = $s['class'];
        if (!(class_exists($dummy) && is_subclass_of($dummy, $Command))) {
            continue;
        }
        $warnings = [];
        if (is_subclass_of($dummy, 'Symfony\\Component\\Console\\SingleCommandApplication')) {
            $warnings[] = '<warning>The script named '.$s['script'].' extends SingleCommandApplication which is not compatible with Composer 2.9+, make sure you extend Symfony\\Component\\Console\\Command instead.</warning>';
        }
        $cmd = new $dummy($s['script']);
        if ($cmd->getName() !== '' && $cmd->getName() !== null && $cmd->getName() !== $s['script']) {
            $warnings[] = '<warning>The script named '.$s['script'].' in composer.json has a mismatched name in its class definition. For consistency, either use the same name, or do not define one inside the class.</warning>';
            $cmd->setName($s['script']);
        }
        if ($cmd->getDescription() === '' && is_string($s['description'])) {
            $cmd->setDescription($s['description']);
        }
        $out[] = ['script' => $s['script'], 'warnings' => $warnings, 'command' => $describe($cmd)];
    }
    $result(['status' => 'ok', 'commands' => $out]);
    exit(0);
}

$className = $data['class'];
if (!class_exists($className)) {
    $result(['status' => 'not-autoloadable']);
    exit(0);
}
if (!is_a($className, $Command, true)) {
    $result(['status' => 'not-command']);
    exit(0);
}
if (!empty($data['reserved'])) {
    $result(['status' => 'reserved']);
    exit(0);
}

// Composer's output: Factory::createOutput (additional styles) with the
// verbosity and decoration of the running Composer; Symfony's configureIO
// reads SHELL_VERBOSITY, which Composer's own configureIO has set.
$verbosity = (int) $data['verbosity'];
$shell = [16 => -1, 32 => 0, 64 => 1, 128 => 2, 256 => 3][$verbosity] ?? 0;
putenv('SHELL_VERBOSITY='.$shell);
$_ENV['SHELL_VERBOSITY'] = $_SERVER['SHELL_VERBOSITY'] = $shell;
$styles = [
    'highlight' => new Symfony\Component\Console\Formatter\OutputFormatterStyle('red'),
    'warning' => new Symfony\Component\Console\Formatter\OutputFormatterStyle('black', 'yellow'),
];
$formatter = new Symfony\Component\Console\Formatter\OutputFormatter(false, $styles);
$output = new Symfony\Component\Console\Output\ConsoleOutput($verbosity, (bool) $data['decorated'], $formatter);

$app = new Symfony\Component\Console\Application();
$app->setCatchExceptions(false);
if (method_exists($app, 'setCatchErrors')) {
    $app->setCatchErrors(false);
}
$app->setAutoExit(false);
$cmd = new $className($data['name']);
if ($data['mode'] === 'direct' && $cmd->getName() !== '' && $cmd->getName() !== null && $cmd->getName() !== $data['name']) {
    // Application::doRun registered the command under the script name
    $cmd->setName($data['name']);
}
if (method_exists($app, 'addCommand')) {
    $app->addCommand($cmd);
} else {
    $app->add($cmd);
}
$app->setDefaultCommand((string) $cmd->getName(), true);
if ($data['mode'] === 'direct') {
    // the command runs within Composer's Application: its global options
    $app->getDefinition()->addOptions([
        new Symfony\Component\Console\Input\InputOption('--profile', null, Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'Display timing and memory usage information'),
        new Symfony\Component\Console\Input\InputOption('--no-plugins', null, Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'Whether to disable plugins.'),
        new Symfony\Component\Console\Input\InputOption('--no-scripts', null, Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'Skips the execution of all scripts defined in composer.json file.'),
        new Symfony\Component\Console\Input\InputOption('--working-dir', '-d', Symfony\Component\Console\Input\InputOption::VALUE_REQUIRED, 'If specified, use the given directory as working directory.'),
        new Symfony\Component\Console\Input\InputOption('--no-cache', null, Symfony\Component\Console\Input\InputOption::VALUE_NONE, 'Prevent use of the cache'),
    ]);
}
try {
    $code = $app->run(new Symfony\Component\Console\Input\StringInput($data['input']), $output);
    $result(['status' => 'ok', 'code' => $code]);
} catch (\Exception $e) {
    $result(['status' => 'exception', 'class' => get_class($e), 'message' => $e->getMessage(), 'code' => $e->getCode()]);
} catch (\Throwable $e) {
    $result(['status' => 'error', 'class' => get_class($e), 'message' => $e->getMessage(), 'code' => $e->getCode()]);
}
`
