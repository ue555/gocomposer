package composer

// Event dispatching to plugin subscribers and the callbacks the PHP plugin
// runtime makes into gocomposer (EventDispatcher::doDispatch of Composer
// f83588392df43d4fd027d84d26de459946ef4e80 for non-string listeners, MIT
// license, (c) Nils Adermann, Jordi Boggiano).

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/autoload"
	"gocomposer/internal/locker"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/solver"
)

// CommandInput is the command line of the running command, given to plugins
// as the Symfony input of CommandEvent/PreCommandRunEvent.
type CommandInput struct {
	Argv []string
	Args []InputArgDef
	Opts []InputOptDef
}

// InputArgDef is a Symfony InputArgument (Mode: REQUIRED=1, OPTIONAL=2,
// IS_ARRAY=4).
type InputArgDef struct {
	Name    string
	Mode    int
	Desc    string
	Default phpjson.Value
}

// InputOptDef is a Symfony InputOption (Mode: VALUE_NONE=1,
// VALUE_REQUIRED=2, VALUE_OPTIONAL=4, VALUE_IS_ARRAY=8, VALUE_NEGATABLE=16).
type InputOptDef struct {
	Name, Shortcut string
	Mode           int
	Desc           string
	Default        phpjson.Value
}

func (in *CommandInput) spec() *phpjson.Map {
	m := phpjson.NewMap()
	argv := []phpjson.Value{}
	for _, a := range in.Argv {
		argv = append(argv, a)
	}
	m.Set("argv", argv)
	args := []phpjson.Value{}
	for _, a := range in.Args {
		am := phpjson.NewMap()
		am.Set("name", a.Name)
		am.Set("mode", int64(a.Mode))
		am.Set("desc", a.Desc)
		am.Set("default", a.Default)
		args = append(args, am)
	}
	m.Set("args", args)
	opts := []phpjson.Value{}
	for _, o := range in.Opts {
		om := phpjson.NewMap()
		om.Set("name", o.Name)
		if o.Shortcut != "" {
			om.Set("shortcut", o.Shortcut)
		}
		om.Set("mode", int64(o.Mode))
		om.Set("desc", o.Desc)
		om.Set("default", o.Default)
		opts = append(opts, om)
	}
	m.Set("opts", opts)
	return m
}

// pluginListener is a listener registered in the plugin runtime.
type pluginListener struct {
	ID     int64
	Prio   int64
	Label  string
	Key    string
	String string // string listeners (commands) added by plugins
}

func (pm *PluginManager) listeners(name string) ([]pluginListener, error) {
	if !pm.active() {
		return nil, nil
	}
	args := phpjson.NewMap()
	args.Set("name", name)
	v, err := pm.host.call("listeners", args)
	if err != nil {
		return nil, err
	}
	var out []pluginListener
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		m := phpjson.ToMap(e)
		if m == nil {
			return
		}
		l := pluginListener{Label: m.Str("label"), Key: m.Str("key"), String: m.Str("string")}
		l.ID, _ = mustGet(m, "id").(int64)
		l.Prio, _ = mustGet(m, "prio").(int64)
		out = append(out, l)
	})
	return out, nil
}

func opSpec(h *pluginHost, op *solver.Operation) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", op.Type)
	if op.Type == solver.OpUpdate {
		m.Set("i", h.ref(op.Initial))
		m.Set("p", h.ref(op.Target))
	} else {
		m.Set("p", h.ref(op.Package))
	}
	return m
}

// eventSpec describes the Composer event object for the plugin runtime.
func (pm *PluginManager) eventSpec(ev *ScriptEvent) *phpjson.Map {
	h := pm.host
	if ev.eid == "" {
		pm.eventSeq++
		ev.eid = "g" + itoa(pm.eventSeq)
	}
	m := phpjson.NewMap()
	m.Set("eid", ev.eid)
	m.Set("name", ev.Name)
	kind := ev.Kind
	if kind == "" {
		kind = "script"
	}
	m.Set("kind", kind)
	m.Set("dev", ev.DevMode)
	args := []phpjson.Value{}
	for _, a := range ev.Args {
		args = append(args, a)
	}
	m.Set("args", args)
	flags := phpjson.NewMap()
	for _, k := range sortedKeys(ev.Flags) {
		flags.Set(k, toPHPValue(ev.Flags[k]))
	}
	m.Set("flags", flags)
	if ev.origin != nil && ev.origin.eid != "" {
		m.Set("origin", ev.origin.eid)
	}
	switch kind {
	case "package", "installer":
		ops := ev.Ops
		key := ev.opsKey
		if key == "" {
			key = pm.curOpsKey
			if ops == nil {
				ops = pm.curOps
			}
			if key == "" {
				pm.opsSeq++
				key = "ops" + itoa(pm.opsSeq)
			}
		}
		m.Set("opsKey", key)
		list := []phpjson.Value{}
		for _, op := range ops {
			list = append(list, opSpec(h, op))
		}
		m.Set("ops", list)
		if ev.Operation != nil {
			m.Set("op", opSpec(h, ev.Operation))
		}
		m.Set("exec", ev.Exec)
		if kind == "installer" {
			present, result := []phpjson.Value{}, []phpjson.Value{}
			for _, p := range pm.txPresent {
				present = append(present, h.ref(p))
			}
			for _, p := range pm.txResult {
				result = append(result, h.ref(p))
			}
			m.Set("present", present)
			m.Set("result", result)
		}
	case "command", "precommand":
		m.Set("command", ev.Command)
		if ev.Input != nil {
			m.Set("input", ev.Input.spec())
		}
	}
	if x := pm.eventExtras[ev.eid]; x != nil {
		for _, k := range x.Keys() {
			m.Set(k, mustGet(x, k))
		}
	}
	return m
}

// pluginTransactionPackages records the packages of the LocalRepoTransaction
// Composer gives to pre-operations-exec listeners.
func pluginTransactionPackages(c *Composer, result, present []*pkg.Package) {
	if pm := pluginManagerOf(c); pm != nil {
		pm.txResult, pm.txPresent = result, present
	}
}

func sortedKeys(m map[string]any) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

func toPHPValue(v any) phpjson.Value {
	switch t := v.(type) {
	case nil, bool, string, int64, float64, *phpjson.Map, []phpjson.Value:
		return t
	case int:
		return int64(t)
	case []string:
		out := []phpjson.Value{}
		for _, s := range t {
			out = append(out, s)
		}
		return out
	}
	return nil
}

// callListener runs a plugin listener (with Composer's makeAutoloader of the
// project autoloader before its first call).
func (pm *PluginManager) callListener(d *EventDispatcher, l pluginListener, ev *ScriptEvent) (int, bool, error) {
	args := phpjson.NewMap()
	args.Set("id", l.ID)
	if loader := pm.projectLoader(l.Key, ev); loader != "" {
		args.Set("loader", loader)
	}
	args.Set("runScripts", d.RunScripts)
	args.Set("event", pm.eventSpec(ev))
	args.Set("frames", pm.frames())
	v, err := pm.host.call("callListener", args)
	if err != nil {
		return 0, false, err
	}
	m := phpjson.ToMap(v)
	ret := 0
	if n, ok := mustGet(m, "ret").(int64); ok {
		ret = int(n)
	}
	return ret, mustGet(m, "stopped") == true, nil
}

// projectLoader is EventDispatcher::makeAutoloader: the project's class
// loader is (re)registered when a listener is called for the first time and
// the installed packages changed since the last registration.
func (pm *PluginManager) projectLoader(key string, ev *ScriptEvent) string {
	if pm.listenerKeys[key] {
		return ""
	}
	pm.listenerKeys[key] = true
	c := pm.c
	packages := c.LocalRepo.CanonicalPackages()
	var parts []string
	for _, p := range packages {
		parts = append(parts, p.Name+"/"+p.Version())
	}
	hash := strings.Join(parts, ",")
	switch ev.Kind {
	case "", "script", "package", "installer":
		pm.devMode = ev.DevMode
		if ev.DevMode {
			hash += "/dev"
		}
	}
	sum := sha256.Sum256([]byte(hash))
	h := hex.EncodeToString(sum[:])
	if h == pm.loaderHash {
		return ""
	}
	pm.loaderHash = h
	im := NewInstallationManager(c)
	code, err := autoload.LoaderData(c.Root(), packages, im.InstallPath, c.Cwd, c.VendorDir(), pm.devMode)
	if err != nil {
		pm.setPendingErr(err)
		return ""
	}
	pm.loaderSeq++
	file := filepath.Join(pm.host.dir, "project-loader-"+itoa(pm.loaderSeq)+".php")
	if err := os.WriteFile(file, []byte(code), 0o600); err != nil {
		pm.setPendingErr(err)
		return ""
	}
	return file
}

func (pm *PluginManager) pushDispatcher(d *EventDispatcher) {
	pm.dispatchers = append(pm.dispatchers, d)
}
func (pm *PluginManager) popDispatcher() {
	if len(pm.dispatchers) > 0 {
		pm.dispatchers = pm.dispatchers[:len(pm.dispatchers)-1]
	}
}

// dispatcher is the EventDispatcher plugins reach through
// $composer->getEventDispatcher().
func (pm *PluginManager) dispatcher() *EventDispatcher {
	if n := len(pm.dispatchers); n > 0 {
		return pm.dispatchers[n-1]
	}
	d := NewEventDispatcher(pm.c)
	d.RunScripts = !pm.c.DisableScripts
	if pm.runScripts != nil {
		d.RunScripts = *pm.runScripts
	}
	return d
}

// callback serves a call from the plugin runtime.
func (pm *PluginManager) callback(op string, a *phpjson.Map) (phpjson.Value, error) {
	if a == nil {
		a = phpjson.NewMap()
	}
	if v, ok, err := pm.apiCallback(op, a); ok {
		return v, err
	}
	c, err := pm.instance(a)
	if err != nil {
		return nil, err
	}
	switch op {
	case "out":
		s, _ := phpjson.ScalarString(mustGet(a, "d"))
		w := c.IO.Out
		if n, _ := mustGet(a, "s").(int64); n == 2 {
			w = c.IO.Err
		}
		_, _ = w.Write([]byte(s))
		return nil, nil
	case "config":
		key := a.Str("key")
		if n, _ := mustGet(a, "flags").(int64); n&1 != 0 {
			return c.Config.GetRelative(key), nil
		}
		return c.Config.Get(key), nil
	case "configHas":
		return c.Config.Raw().Has(a.Str("key")), nil
	case "configAll", "configRaw":
		out := phpjson.NewMap()
		out.Set("repositories", c.Config.Repositories())
		if op == "configRaw" {
			out.Set("config", c.Config.Raw())
			return out, nil
		}
		relative := false
		if n, _ := mustGet(a, "flags").(int64); n&1 != 0 {
			relative = true
		}
		cfg := phpjson.NewMap()
		for _, k := range c.Config.Raw().Keys() {
			if relative {
				cfg.Set(k, c.Config.GetRelative(k))
			} else {
				cfg.Set(k, c.Config.Get(k))
			}
		}
		out.Set("config", cfg)
		return out, nil
	case "configSource":
		src, err := c.Config.SourceOfValue(a.Str("key"))
		if err != nil {
			return nil, err
		}
		return src, nil
	case "localRepo":
		return pm.localRepoState(c, a.Str("known"))
	case "lock":
		switch a.Str("m") {
		case "isLocked":
			locked, _ := c.Locker.IsLocked()
			return locked, nil
		case "isFresh":
			return c.Locker.IsFresh(), nil
		default:
			if locked, _ := c.Locker.IsLocked(); !locked {
				return nil, nil
			}
			data, err := c.Locker.LockData()
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	case "contentHash":
		s, err := locker.ContentHash([]byte(a.Str("json")))
		if err != nil {
			return nil, err
		}
		return s, nil
	case "libOp":
		return pm.libOp(a)
	case "binaryCaller":
		return determineBinaryCaller(a.Str("bin")), nil
	case "dispatch":
		return pm.dispatchFromPlugin(a)
	case "hasListeners":
		d := pm.dispatcher()
		return d.HasListeners(a.Str("name")), nil
	case "setRunScripts":
		b := mustGet(a, "run") == true
		pm.runScripts = &b
		if n := len(pm.dispatchers); n > 0 {
			pm.dispatchers[n-1].RunScripts = b
		}
		return nil, nil
	case "installersChanged":
		pm.installers = nil
		pm.installPaths = nil
		return nil, nil
	case "isPluginAllowed":
		return pm.isPluginAllowed(a.Str("name"), mustGet(a, "global") == true, mustGet(a, "optional") == true, true)
	}
	return nil, &UserError{Msg: "Unknown gocomposer plugin callback " + op, Code: 1}
}

// localRepoState describes the local repository for the runtime's
// InstalledFilesystemRepository (only when it changed since known).
func (pm *PluginManager) localRepoState(c *Composer, known string) (phpjson.Value, error) {
	repo := pm.curRepo
	if repo == nil || c != pm.c {
		repo = c.LocalRepo
	}
	packages, err := repo.Packages()
	if err != nil {
		return nil, err
	}
	h := pm.host
	var ids []string
	for _, p := range packages {
		ids = append(ids, itoa(h.idOf(p)))
	}
	key := strings.Join(ids, ",")
	devNames := []phpjson.Value{}
	for _, n := range repo.DevPackageNames {
		devNames = append(devNames, n)
		key += "|" + n
	}
	var dev phpjson.Value
	if repo.DevMode != nil {
		dev = *repo.DevMode
		if *repo.DevMode {
			key += "|dev"
		} else {
			key += "|nodev"
		}
	}
	fresh := repo.IsFresh()
	if fresh {
		key += "|fresh"
	}
	out := phpjson.NewMap()
	if key == known {
		// the data of new packages was not needed
		out.Set("same", true)
		return out, nil
	}
	refs := []phpjson.Value{}
	for _, p := range packages {
		refs = append(refs, h.ref(p))
	}
	out.Set("key", key)
	out.Set("pkgs", refs)
	out.Set("dev", dev)
	out.Set("devNames", devNames)
	out.Set("fresh", fresh)
	return out, nil
}

// dispatchFromPlugin dispatches an event dispatched by plugin code
// ($composer->getEventDispatcher()->dispatch()).
func (pm *PluginManager) dispatchFromPlugin(a *phpjson.Map) (phpjson.Value, error) {
	kind := a.Str("kind")
	if kind == "unsupported" {
		return nil, &UserError{Msg: "Dispatching package or installer events from plugins is not supported by gocomposer. Run the command with --no-plugins to skip plugins.", Code: 1}
	}
	ev := &ScriptEvent{Name: a.Str("name"), Kind: kind, DevMode: mustGet(a, "dev") == true, eid: a.Str("eid")}
	ev.Command = a.Str("command")
	phpjson.Each(mustGet(a, "args"), func(_ string, v phpjson.Value) {
		s, _ := phpjson.ScalarString(v)
		ev.Args = append(ev.Args, s)
	})
	if f := phpjson.ToMap(mustGet(a, "flags")); f != nil && f.Len() > 0 {
		ev.Flags = map[string]any{}
		for _, k := range f.Keys() {
			ev.Flags[k] = mustGet(f, k)
		}
	}
	code, err := pm.dispatcher().dispatch(ev)
	if err != nil {
		return nil, err
	}
	return int64(code), nil
}

// freeEvent releases the runtime event object of a finished dispatch.
func (pm *PluginManager) freeEvent(ev *ScriptEvent) {
	if !pm.active() || ev.eid == "" || !strings.HasPrefix(ev.eid, "g") {
		return
	}
	args := phpjson.NewMap()
	args.Set("eid", ev.eid)
	_, _ = pm.host.call("freeEvent", args)
}
