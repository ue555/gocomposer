package composer

// PluginEvents::COMMAND dispatched by commands that use the listeners'
// return value (ValidateCommand: max($eventCode, $exitCode)), Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT, (c) Nils Adermann, Jordi
// Boggiano).

import "gocomposer/internal/cio"

// DispatchCommandEventCode dispatches the command event and returns the
// EventDispatcher::dispatch result (1 when a listener returned false).
func (d *EventDispatcher) DispatchCommandEventCode(command string, input *CommandInput) (int, error) {
	if pm := pluginManagerOf(d.c); !pm.active() && !d.HasListeners("command") {
		return 0, nil
	}
	return d.dispatch(&ScriptEvent{Name: "command", Kind: "command", Command: command, Input: input})
}

// CreateInstance is BaseCommand::createComposerInstance: a further Composer
// instance next to the command's one (ValidateCommand). Unlike the instances
// Require/RemoveCommand create for the update, the plugins of the existing
// instance stay active (no PluginManager::deactivateInstalledPlugins).
func CreateInstance(io *cio.IO, opts Options) (*Composer, error) {
	prev := creatingFromPlugin
	creatingFromPlugin = true
	defer func() { creatingFromPlugin = prev }()
	return Create(io, opts)
}
