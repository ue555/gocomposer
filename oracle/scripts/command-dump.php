<?php
// Dumps the console definitions (names, aliases, descriptions, arguments,
// options, help) of the pinned Composer checkout as JSON. The output is
// embedded by gocomposer as internal/cli/commanddefs.json.
// Usage: php command-dump.php > commanddefs.json
require '/oracle-source/vendor/autoload.php';
putenv('COMPOSER_HOME=/tmp/cd-home');
$app = new Composer\Console\Application();
$app->setAutoExit(false);
$ref = new ReflectionMethod($app, 'getDefaultCommands');
$ref->setAccessible(true);
$cmds = [];
$all = $app->all();
ksort($all);
$seen = [];
foreach ($all as $name => $cmd) {
    if (isset($seen[spl_object_id($cmd)])) continue;
    $seen[spl_object_id($cmd)] = true;
    $def = $cmd->getDefinition();
    $args = [];
    foreach ($def->getArguments() as $a) {
        $args[] = ['name' => $a->getName(), 'required' => $a->isRequired(), 'array' => $a->isArray(), 'description' => $a->getDescription(), 'default' => $a->getDefault()];
    }
    $opts = [];
    foreach ($def->getOptions() as $o) {
        $opts[] = ['name' => $o->getName(), 'shortcut' => $o->getShortcut(), 'accept_value' => $o->acceptValue(), 'value_required' => $o->isValueRequired(), 'value_optional' => $o->isValueOptional(), 'array' => $o->isArray(), 'negatable' => $o->isNegatable(), 'description' => $o->getDescription(), 'default' => $o->getDefault()];
    }
    $cmds[] = ['name' => $cmd->getName(), 'aliases' => $cmd->getAliases(), 'hidden' => $cmd->isHidden(), 'description' => $cmd->getDescription(), 'help' => $cmd->getHelp(), 'usages' => $cmd->getUsages(), 'arguments' => $args, 'options' => $opts, 'synopsis' => $cmd->getSynopsis()];
}
$g = [];
foreach ($app->getDefinition()->getOptions() as $o) {
    $g[] = ['name' => $o->getName(), 'shortcut' => $o->getShortcut(), 'accept_value' => $o->acceptValue(), 'value_required' => $o->isValueRequired(), 'value_optional' => $o->isValueOptional(), 'array' => $o->isArray(), 'negatable' => $o->isNegatable(), 'description' => $o->getDescription(), 'default' => $o->getDefault()];
}
echo json_encode(['source' => 'composer/composer', 'commands' => $cmds, 'global_options' => $g], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE), "\n";
