package composer

// Commands of the CommandProvider capability (Application::getPluginCommands
// of Composer f83588392df43d4fd027d84d26de459946ef4e80, MIT license,
// (c) Nils Adermann, Jordi Boggiano) and checks for plugin features
// gocomposer does not run.

import (
	"strings"

	"gocomposer/internal/phpjson"
)

// PluginCommand is a command provided by a plugin.
type PluginCommand struct {
	Index       int64
	Name        string
	Class       string
	Description string
	Help        string
	Aliases     []string
	Hidden      bool
	Synopsis    string
	Usages      []string
	Args        []PluginCommandArg
	Opts        []PluginCommandOpt
}

// PluginCommandArg is an InputArgument of a plugin command.
type PluginCommandArg struct {
	Name            string
	Required, Array bool
	Desc            string
	Default         phpjson.Value
}

// PluginCommandOpt is an InputOption of a plugin command.
type PluginCommandOpt struct {
	Name, Shortcut                             string
	None, Required, Optional, Array, Negatable bool
	Desc                                       string
	Default                                    phpjson.Value
}

// PluginCommands returns the commands of the active plugins' CommandProvider
// capabilities.
func PluginCommands(c *Composer) ([]PluginCommand, error) {
	pm := pluginManagerOf(c)
	if !pm.active() {
		return nil, nil
	}
	v, err := pm.host.call("commands", nil)
	if err != nil {
		return nil, err
	}
	var out []PluginCommand
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		m := phpjson.ToMap(e)
		if m == nil {
			return
		}
		pc := PluginCommand{Name: m.Str("name"), Class: m.Str("class"), Description: m.Str("description"), Help: m.Str("help"), Synopsis: m.Str("synopsis"), Hidden: mustGet(m, "hidden") == true}
		pc.Index, _ = mustGet(m, "index").(int64)
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
		out = append(out, pc)
	})
	return out, nil
}

// RunPluginCommand runs a plugin command with the raw command line (bound
// by the command's own Symfony definition in the plugin runtime).
func RunPluginCommand(c *Composer, index int64, argv []string) (int, error) {
	pm := pluginManagerOf(c)
	if !pm.active() {
		return 1, &UserError{Msg: "The plugins providing this command are not active", Code: 1}
	}
	args := phpjson.NewMap()
	args.Set("index", index)
	list := []phpjson.Value{}
	for _, a := range argv {
		list = append(list, a)
	}
	args.Set("argv", list)
	args.Set("frames", pm.frames())
	v, err := pm.host.call("runCommand", args)
	if err != nil {
		// Symfony's Application::run exits with the exception's code
		if pe, ok := err.(*PHPException); ok && pe.Code > 0 {
			return min(pe.Code, 255), err
		}
		return 1, err
	}
	n, _ := v.(int64)
	return int(n), nil
}

// unsupportedPluginEvents are events gocomposer does not dispatch to
// plugins (none: pre/post-file-download and pre-pool-create are dispatched,
// see plugin_dlevents.go).
var unsupportedPluginEvents []string

// CheckUnsupportedPluginEvents fails when an active plugin subscribed to an
// event gocomposer cannot dispatch, instead of silently skipping it.
func CheckUnsupportedPluginEvents(c *Composer) error {
	pm := pluginManagerOf(c)
	if !pm.active() {
		return nil
	}
	for _, name := range unsupportedPluginEvents {
		ls, err := pm.listeners(name)
		if err != nil {
			return err
		}
		if len(ls) == 0 {
			continue
		}
		var labels []string
		for _, l := range ls {
			if l.Label != "" {
				labels = append(labels, l.Label)
			} else if l.String != "" {
				labels = append(labels, l.String)
			} else {
				labels = append(labels, l.Key)
			}
		}
		return &UserError{Msg: "A Composer plugin subscribes to the \"" + name + "\" event (" + strings.Join(labels, ", ") + "), which gocomposer does not dispatch to plugins. Run the command with --no-plugins to skip plugins.", Code: 1}
	}
	return nil
}
