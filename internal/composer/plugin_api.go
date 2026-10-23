package composer

// The Composer services plugins use beyond the basic API, served to the PHP
// plugin runtime (internal/res/pluginrt/api.php): HttpDownloader
// (Util/HttpDownloader.php), RepositoryManager and repositories,
// Factory::create, Config::merge, DownloadManager and its downloaders,
// AutoloadGenerator::dump, Installer::run and the application/installer
// objects plugins find on the call stack. The semantics follow Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT license, (c) Nils Adermann,
// Jordi Boggiano); all work is done by gocomposer's Go implementation.

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// PluginApp describes the running command line for plugins: Composer runs
// plugin code inside Application::doRun($input), and plugins like
// symfony/flex look the application and its input up on the call stack.
type PluginApp struct {
	// Argv is the command line without the program name.
	Argv []string
	// Cmd is the command being run (nil while gocomposer looks for plugin
	// commands, when the input is not bound yet).
	Cmd *PluginAppCommand
	// Find resolves a command name like Application::find (nil, nil when
	// unknown).
	Find func(name string) (*PluginAppCommand, error)
	// Global is set while the global command runs.
	Global bool
}

// PluginAppCommand is a gocomposer command as described to plugins.
type PluginAppCommand struct {
	Name        string
	Aliases     []string
	Description string
	Args        []PluginCommandArg
	Opts        []PluginCommandOpt
}

var pluginApp *PluginApp

// SetPluginApp sets the command line plugins see (nil: none).
func SetPluginApp(a *PluginApp) { pluginApp = a }

// CurrentPluginApp returns the command line plugins see.
func CurrentPluginApp() *PluginApp { return pluginApp }

// PluginAppArgv returns the command line after the plugins of c were
// activated: plugins may rewrite the input (symfony/flex resolves package
// aliases of require/update/remove). ok is false when it did not change.
func PluginAppArgv(c *Composer) (argv []string, ok bool, err error) {
	pm := pluginManagerOf(c)
	if !pm.active() || pluginApp == nil {
		return nil, false, nil
	}
	v, err := pm.host.call("appState", nil)
	if err != nil {
		return nil, false, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		return nil, false, nil
	}
	list, ok2 := m.Get("argv")
	if !ok2 || list == nil {
		return nil, false, nil
	}
	phpjson.Each(list, func(_ string, t phpjson.Value) {
		s, _ := phpjson.ScalarString(t)
		argv = append(argv, s)
	})
	if len(argv) == len(pluginApp.Argv) {
		same := true
		for i := range argv {
			if argv[i] != pluginApp.Argv[i] {
				same = false
				break
			}
		}
		if same {
			return nil, false, nil
		}
	}
	pluginApp.Argv = argv
	return argv, true, nil
}

func appCommandSpec(cmd *PluginAppCommand) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("name", cmd.Name)
	aliases := []phpjson.Value{}
	for _, a := range cmd.Aliases {
		aliases = append(aliases, a)
	}
	m.Set("aliases", aliases)
	m.Set("description", cmd.Description)
	args := []phpjson.Value{}
	for _, a := range cmd.Args {
		am := phpjson.NewMap()
		am.Set("name", a.Name)
		am.Set("required", a.Required)
		am.Set("array", a.Array)
		am.Set("desc", a.Desc)
		am.Set("default", a.Default)
		args = append(args, am)
	}
	m.Set("args", args)
	opts := []phpjson.Value{}
	for _, o := range cmd.Opts {
		om := phpjson.NewMap()
		om.Set("name", o.Name)
		om.Set("shortcut", o.Shortcut)
		om.Set("none", o.None)
		om.Set("required", o.Required)
		om.Set("optional", o.Optional)
		om.Set("array", o.Array)
		om.Set("negatable", o.Negatable)
		om.Set("desc", o.Desc)
		om.Set("default", o.Default)
		opts = append(opts, om)
	}
	m.Set("opts", opts)
	return m
}

func appSpec() phpjson.Value {
	if pluginApp == nil {
		return nil
	}
	m := phpjson.NewMap()
	argv := []phpjson.Value{}
	for _, a := range pluginApp.Argv {
		argv = append(argv, a)
	}
	m.Set("argv", argv)
	if pluginApp.Cmd != nil {
		m.Set("cmd", appCommandSpec(pluginApp.Cmd))
	}
	return m
}

// running gocomposer Installers (Installer::run on the call stack)
var (
	installerMu    sync.Mutex
	installerSeq   int
	installerStack []*runningInstaller
	installerByID  = map[int]*runningInstaller{}
)

type runningInstaller struct {
	id int
	in *Installer
}

// pluginInstallerScope registers in as the Installer running until the
// returned function is called.
func pluginInstallerScope(in *Installer) func() {
	installerMu.Lock()
	installerSeq++
	ri := &runningInstaller{id: installerSeq, in: in}
	installerStack = append(installerStack, ri)
	installerByID[ri.id] = ri
	installerMu.Unlock()
	return func() {
		installerMu.Lock()
		for i, r := range installerStack {
			if r == ri {
				installerStack = append(installerStack[:i], installerStack[i+1:]...)
				break
			}
		}
		installerMu.Unlock()
	}
}

func currentInstaller() *runningInstaller {
	installerMu.Lock()
	defer installerMu.Unlock()
	if n := len(installerStack); n > 0 {
		return installerStack[n-1]
	}
	return nil
}

func platformFilterSpec(f solver.PlatformFilter) phpjson.Value {
	switch t := f.(type) {
	case solver.IgnoreAll:
		return true
	case *solver.IgnoreList:
		// the runtime only hands the filter back to Installer::run (see
		// installerRun): the running installer's filter is reused
		_ = t
		return []phpjson.Value{baseFilterToken}
	}
	return false
}

// baseFilterToken stands for the platform requirement filter of the
// installer an Installer object was cloned from.
const baseFilterToken = "__gocomposer_base_filter__"

// frames describes the call stack plugin code runs under (Host::framed).
func (pm *PluginManager) frames() *phpjson.Map {
	m := phpjson.NewMap()
	if pluginApp != nil {
		m.Set("app", true)
		if pluginApp.Global {
			m.Set("global", true)
		}
	}
	if ri := currentInstaller(); ri != nil && pm.host != nil {
		if !pm.host.installersSent[ri.id] {
			pm.host.installersSent[ri.id] = true
			st := phpjson.NewMap()
			st.Set("dryRun", ri.in.DryRun)
			st.Set("platformFilter", platformFilterSpec(ri.in.PlatformFilter))
			args := phpjson.NewMap()
			args.Set("id", int64(ri.id))
			args.Set("opts", st)
			if _, err := pm.host.call("installerStub", args); err != nil {
				pm.setPendingErr(err)
			}
		}
		m.Set("installer", int64(ri.id))
	}
	return m
}

// instance returns the Composer instance a runtime object belongs to (0:
// the one running the plugins, others created with Factory::create).
func (pm *PluginManager) instance(a *phpjson.Map) (*Composer, error) {
	n, _ := mustGet(a, "inst").(int64)
	if n == 0 {
		return pm.c, nil
	}
	if c, ok := pm.instances[n]; ok {
		return c, nil
	}
	return nil, &UserError{Msg: fmt.Sprintf("Unknown gocomposer Composer instance %d", n), Code: 1}
}

// creatingFromPlugin is set while Factory::create runs for a plugin: the
// other Composer instances keep their plugins active.
var creatingFromPlugin bool

// apiCallback serves the callbacks of api.php (ok is false for other ops).
func (pm *PluginManager) apiCallback(op string, a *phpjson.Map) (phpjson.Value, bool, error) {
	switch op {
	case "http":
		v, err := pm.httpJobs(a)
		return v, true, err
	case "appFind":
		v, err := pm.appFind(a.Str("name"))
		return v, true, err
	case "repos":
		c, err := pm.instance(a)
		if err != nil {
			return nil, true, err
		}
		out := []phpjson.Value{}
		for i, r := range c.Repositories {
			out = append(out, pm.repoDesc(c, r, i))
		}
		return out, true, nil
	case "repoOp":
		v, err := pm.repoOp(a)
		return v, true, err
	case "repoCreate":
		v, err := pm.repoCreate(a)
		return v, true, err
	case "repoAdd":
		return nil, true, pm.repoAdd(a)
	case "configMerge":
		c, err := pm.instance(a)
		if err != nil {
			return nil, true, err
		}
		data := phpjson.ToMap(mustGet(a, "config"))
		if data == nil {
			return nil, true, nil
		}
		source := a.Str("source")
		if source == "" {
			source = "unknown"
		}
		// like Config::merge, the repositories already created are kept
		c.Config.MergeSource(data, source)
		return nil, true, nil
	case "factoryCreate":
		v, err := pm.factoryCreate(a)
		return v, true, err
	case "installerRun":
		v, err := pm.installerRun(a)
		return v, true, err
	case "installerSet":
		id, _ := mustGet(a, "id").(int64)
		installerMu.Lock()
		ri := installerByID[int(id)]
		installerMu.Unlock()
		if ri != nil && mustGet(a, "silentSuggestions") == true {
			ri.in.suggestions.io = &cio.IO{Out: discard{}, Err: discard{}, Verbosity: cio.Quiet}
		}
		return nil, true, nil
	case "dl":
		async, err := pm.downloaderOp(a)
		if err != nil {
			return nil, true, err
		}
		m := phpjson.NewMap()
		m.Set("async", async)
		return m, true, nil
	case "opDone":
		return nil, true, pm.opDone(a.Str("key"))
	case "imOutputProgress":
		if pm.curIM != nil {
			pm.curIM.OutputProgress = mustGet(a, "v") == true
		}
		pm.outputProgress = new(bool)
		*pm.outputProgress = mustGet(a, "v") == true
		return nil, true, nil
	case "autoloadDump":
		v, err := pm.autoloadDump(a)
		return v, true, err
	case "localRepoOp":
		return nil, true, pm.localRepoOp(a)
	case "validateSchema":
		data, err := phpjson.Decode([]byte(a.Str("json")))
		if err != nil {
			return nil, true, err
		}
		mode, _ := mustGet(a, "schema").(int64)
		msgs, err := jsonschema.ValidateComposerJSON(data, int(mode))
		if err != nil {
			return nil, true, err
		}
		out := []phpjson.Value{}
		for _, m := range msgs {
			out = append(out, m)
		}
		return out, true, nil
	}
	return nil, false, nil
}

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }

func (pm *PluginManager) appFind(name string) (phpjson.Value, error) {
	if pluginApp == nil || pluginApp.Find == nil {
		return nil, nil
	}
	cmd, err := pluginApp.Find(name)
	if err != nil || cmd == nil {
		return nil, nil
	}
	m := phpjson.NewMap()
	m.Set("cmd", appCommandSpec(cmd))
	return m, nil
}

// httpJobs runs the queued jobs of a runtime HttpDownloader concurrently.
func (pm *PluginManager) httpJobs(a *phpjson.Map) (phpjson.Value, error) {
	var jobs []*phpjson.Map
	phpjson.Each(mustGet(a, "jobs"), func(_ string, v phpjson.Value) {
		if m := phpjson.ToMap(v); m != nil {
			jobs = append(jobs, m)
		} else {
			jobs = append(jobs, phpjson.NewMap())
		}
	})
	max64, _ := mustGet(a, "max").(int64)
	maxJobs := int(max64)
	if maxJobs < 1 {
		maxJobs = 12
	}
	results := make([]phpjson.Value, len(jobs))
	sem := make(chan struct{}, maxJobs)
	var wg sync.WaitGroup
	for i, j := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, j *phpjson.Map) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = pm.httpJob(j)
		}(i, j)
	}
	wg.Wait()
	out := make([]phpjson.Value, len(results))
	copy(out, results)
	return out, nil
}

// unsupportedTransportOptions lists stream context options gocomposer's HTTP
// client can not apply per request.
var unsupportedTransportOptions = map[string][]string{
	"ssl":  {"cafile", "capath", "local_cert", "local_pk", "passphrase", "allow_self_signed", "peer_fingerprint"},
	"http": {"proxy", "request_fulluri"},
}

func (pm *PluginManager) httpJob(j *phpjson.Map) phpjson.Value {
	c := pm.c
	u := j.Str("url")
	options := phpjson.ToMap(mustGet(j, "options"))
	errResult := func(err error) phpjson.Value {
		e := phpjson.NewMap()
		e.Set("class", "Composer\\Downloader\\TransportException")
		e.Set("message", err.Error())
		// TransportException's default code is 400 (curl errors, blocked
		// URLs); HTTP errors carry the status code
		code := int64(400)
		if te, ok := err.(*httpdl.TransportError); ok && te.StatusCode != 0 {
			code = int64(te.StatusCode)
			e.Set("status", code)
		}
		e.Set("code", code)
		r := phpjson.NewMap()
		r.Set("error", e)
		return r
	}
	if options != nil {
		for group, keys := range unsupportedTransportOptions {
			gm := options.Map(group)
			if gm == nil {
				continue
			}
			for _, k := range keys {
				if v, ok := gm.Get(k); ok && v != nil && v != false {
					if group == "ssl" && (k == "cafile" || k == "capath") {
						if s, _ := v.(string); s != "" && s == c.Config.GetString(k) {
							continue
						}
					}
					return errResult(&UserError{Msg: "The " + group + "." + k + " transport option is not supported by gocomposer's HTTP client (requested by a Composer plugin for " + httpdl.Sanitize(u) + ")"})
				}
			}
		}
	}
	var resp *httpdl.Response
	var err error
	if to := j.Str("copyTo"); to != "" {
		resp, err = c.HTTP.Download(u, to, transportHeadersOf(options)...)
	} else {
		resp, err = c.HTTP.Request(u, options)
	}
	if err != nil {
		return errResult(err)
	}
	r := phpjson.NewMap()
	r.Set("status", int64(resp.StatusCode))
	headers := []phpjson.Value{}
	proto := resp.Proto
	if proto == "" {
		proto = "HTTP/1.1"
	}
	status := resp.Status
	if status == "" {
		status = fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	headers = append(headers, proto+" "+status)
	var keys []string
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range resp.Header[k] {
			headers = append(headers, strings.ToLower(k)+": "+v)
		}
	}
	r.Set("headers", headers)
	if j.Str("copyTo") == "" {
		r.Set("body", base64.StdEncoding.EncodeToString(resp.Body))
	}
	return r
}

func transportHeadersOf(options *phpjson.Map) []string {
	if options == nil {
		return nil
	}
	hm := options.Map("http")
	if hm == nil {
		return nil
	}
	var hs []string
	if h, ok := hm.Get("header"); ok {
		switch t := h.(type) {
		case string:
			hs = append(hs, strings.Split(strings.ReplaceAll(t, "\r\n", "\n"), "\n")...)
		default:
			phpjson.Each(t, func(_ string, v phpjson.Value) {
				if s, ok := v.(string); ok {
					hs = append(hs, s)
				}
			})
		}
	}
	return hs
}

// repository descriptors for the runtime

func (pm *PluginManager) repoID(r repo.Repository) int64 {
	if pm.repoIDs == nil {
		pm.repoIDs = map[repo.Repository]int64{}
		pm.repoByID = map[int64]repo.Repository{}
	}
	if id, ok := pm.repoIDs[r]; ok {
		return id
	}
	id := int64(len(pm.repoIDs) + 1)
	pm.repoIDs[r] = id
	pm.repoByID[id] = r
	return id
}

func repoKindName(k repo.Kind) string {
	switch k {
	case repo.KindComposer:
		return "composer"
	case repo.KindPath:
		return "path"
	case repo.KindVcs:
		return "vcs"
	case repo.KindPackage:
		return "package"
	case repo.KindArtifact:
		return "artifact"
	case repo.KindPlatform:
		return "platform"
	case repo.KindLock:
		return "lock"
	case repo.KindRoot:
		return "root"
	case repo.KindInstalled:
		return "installed"
	}
	return "array"
}

// repoDesc describes r; index is its position in the configured repository
// list (-1 for repositories that are not configured ones).
func (pm *PluginManager) repoDesc(c *Composer, r repo.Repository, index int) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("rid", pm.repoID(r))
	kind := repoKindName(r.Kind())
	m.Set("kind", kind)
	m.Set("name", r.RepoName())
	cfg := pm.repoConfigs[r]
	if cfg == nil && index >= 0 && c != nil {
		repos := c.Config.Repositories()
		if repos != nil && repos.Len() == len(c.Repositories) {
			cfg = phpjson.ToMap(mustGet(repos, repos.Keys()[index]))
		}
	}
	if cfg == nil {
		cfg = phpjson.NewMap()
	}
	if kind == "composer" && !cfg.Has("options") {
		// ComposerRepository::__construct defaults the options
		cfg = cfg.Clone()
		cfg.Set("options", []phpjson.Value{})
	}
	m.Set("config", cfg)
	switch kind {
	case "platform", "lock", "array":
		refs := []phpjson.Value{}
		if ps, err := r.Packages(); err == nil {
			for _, p := range ps {
				refs = append(refs, pm.host.ref(p))
			}
		}
		m.Set("pkgs", refs)
	case "installed", "root":
		for inst, ic := range pm.instances {
			if ic.LocalRepo == r {
				m.Set("inst", inst)
			}
		}
	}
	return m
}

func parseConstraintArg(v phpjson.Value) (semver.Constraint, error) {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil, nil
	}
	c, err := semver.ParseConstraints(s)
	if err != nil {
		return nil, &UserError{Msg: err.Error(), Code: 1}
	}
	return c, nil
}

func (pm *PluginManager) repoOp(a *phpjson.Map) (phpjson.Value, error) {
	id, _ := mustGet(a, "rid").(int64)
	r, ok := pm.repoByID[id]
	if !ok {
		return nil, &UserError{Msg: fmt.Sprintf("Unknown gocomposer repository %d", id), Code: 1}
	}
	h := pm.host
	refs := func(ps []*pkg.Package) []phpjson.Value {
		out := []phpjson.Value{}
		for _, p := range ps {
			out = append(out, h.ref(p))
		}
		return out
	}
	switch a.Str("m") {
	case "find":
		name := strings.ToLower(a.Str("name"))
		c, err := parseConstraintArg(mustGet(a, "c"))
		if err != nil {
			return nil, err
		}
		if mustGet(a, "one") == true {
			if c == nil {
				c = &semver.MatchAll{}
			}
			p := repo.FindPackage(r, name, c)
			if p == nil {
				return []phpjson.Value{}, nil
			}
			return refs([]*pkg.Package{p}), nil
		}
		ps, err := r.FindPackages(name, c)
		if err != nil {
			return nil, err
		}
		return refs(ps), nil
	case "packages":
		ps, err := r.Packages()
		if err != nil {
			return nil, err
		}
		return refs(ps), nil
	case "load":
		var names []repo.NameConstraint
		var perr error
		phpjson.Each(mustGet(a, "names"), func(_ string, v phpjson.Value) {
			m := phpjson.ToMap(v)
			if m == nil {
				return
			}
			c, err := parseConstraintArg(mustGet(m, "c"))
			if err != nil && perr == nil {
				perr = err
			}
			names = append(names, repo.NameConstraint{Name: strings.ToLower(m.Str("name")), Constraint: c})
		})
		if perr != nil {
			return nil, perr
		}
		intMap := func(v phpjson.Value) map[string]int {
			out := map[string]int{}
			phpjson.Each(v, func(k string, x phpjson.Value) {
				if n, ok := x.(int64); ok {
					out[k] = int(n)
				}
			})
			return out
		}
		loaded := map[string]map[string]bool{}
		phpjson.Each(mustGet(a, "loaded"), func(name string, x phpjson.Value) {
			phpjson.Each(x, func(ver string, _ phpjson.Value) {
				if loaded[name] == nil {
					loaded[name] = map[string]bool{}
				}
				loaded[name][ver] = true
			})
		})
		res, err := r.LoadPackages(names, intMap(mustGet(a, "acceptable")), intMap(mustGet(a, "flags")), loaded)
		if err != nil {
			return nil, err
		}
		out := phpjson.NewMap()
		found := []phpjson.Value{}
		for _, n := range res.NamesFound {
			found = append(found, n)
		}
		out.Set("namesFound", found)
		out.Set("packages", refs(res.Packages))
		return out, nil
	case "search":
		var typ *string
		if s, ok := mustGet(a, "type").(string); ok {
			typ = &s
		}
		mode, _ := mustGet(a, "mode").(int64)
		entries, err := repo.SearchWithMode(r, a.Str("query"), int(mode), typ)
		if err != nil {
			return nil, err
		}
		out := []phpjson.Value{}
		for _, e := range entries {
			if e.Raw != nil {
				out = append(out, e.Raw)
				continue
			}
			m := phpjson.NewMap()
			m.Set("name", e.Name)
			if e.Description != nil {
				m.Set("description", *e.Description)
			} else {
				m.Set("description", nil)
			}
			if e.Abandoned != nil {
				m.Set("abandoned", e.Abandoned)
			}
			if e.URL != "" {
				m.Set("url", e.URL)
			}
			out = append(out, m)
		}
		return out, nil
	case "providers":
		ps, err := r.Providers(strings.ToLower(a.Str("name")))
		if err != nil {
			return nil, err
		}
		out := phpjson.NewMap()
		for _, p := range ps {
			m := phpjson.NewMap()
			m.Set("name", p.Name)
			m.Set("description", p.Description)
			m.Set("type", p.Type)
			out.Set(p.Name, m)
		}
		return out, nil
	}
	return nil, &UserError{Msg: "Unknown repository operation " + a.Str("m"), Code: 1}
}

// repoCreate is RepositoryManager::createRepository.
func (pm *PluginManager) repoCreate(a *phpjson.Map) (phpjson.Value, error) {
	c, err := pm.instance(a)
	if err != nil {
		return nil, err
	}
	typ := a.Str("type")
	rc := phpjson.ToMap(mustGet(a, "config"))
	if rc == nil {
		rc = phpjson.NewMap()
	} else {
		rc = rc.Clone()
	}
	if !rc.Has("type") {
		rc.Set("type", typ)
	}
	name, _ := mustGet(a, "name").(string)
	r, err := c.createRepository(typ, rc, name)
	if err != nil {
		return nil, err
	}
	if pm.repoConfigs == nil {
		pm.repoConfigs = map[repo.Repository]*phpjson.Map{}
	}
	pm.repoConfigs[r] = rc
	pm.repoHook(r)
	return pm.repoDesc(c, r, -1), nil
}

// repoAdd is RepositoryManager::addRepository/prependRepository.
func (pm *PluginManager) repoAdd(a *phpjson.Map) error {
	c, err := pm.instance(a)
	if err != nil {
		return err
	}
	var r repo.Repository
	if id, ok := mustGet(a, "rid").(int64); ok {
		r = pm.repoByID[id]
		if r == nil {
			return &UserError{Msg: fmt.Sprintf("Unknown gocomposer repository %d", id), Code: 1}
		}
	} else {
		arr := repo.NewArray()
		var perr error
		phpjson.Each(mustGet(a, "pkgs"), func(_ string, v phpjson.Value) {
			p, err := pm.host.pkgFromRef(v)
			if err != nil {
				perr = err
				return
			}
			if p != nil {
				arr.AddPackage(p)
			}
		})
		if perr != nil {
			return perr
		}
		if n := a.Str("name"); n != "" {
			arr.SetName(n)
		}
		r = arr
	}
	for _, e := range c.Repositories {
		if e == r {
			return nil
		}
	}
	if mustGet(a, "prepend") == true {
		c.Repositories = append([]repo.Repository{r}, c.Repositories...)
	} else {
		c.Repositories = append(c.Repositories, r)
	}
	return nil
}

// factoryCreate is Factory::create from a plugin.
func (pm *PluginManager) factoryCreate(a *phpjson.Map) (phpjson.Value, error) {
	file, _ := mustGet(a, "file").(string)
	// Factory::createComposer: a custom composer.json path changes the
	// working directory of the instance to the file's directory
	cwd := ""
	if file != "" {
		if st, err := os.Stat(file); err == nil && !st.IsDir() {
			if abs, err := filepath.Abs(filepath.Dir(file)); err == nil {
				cwd = abs
			}
		}
	}
	prev := creatingFromPlugin
	creatingFromPlugin = true
	nc, err := Create(pm.c.IO, Options{ComposerFile: file, Cwd: cwd, DisablePlugins: mustGet(a, "disablePlugins") == true || pm.c.DisablePlugins, DisableScripts: mustGet(a, "disableScripts") == true, FullLoad: true})
	creatingFromPlugin = prev
	if err != nil {
		if ue, ok := err.(*UserError); ok && strings.HasPrefix(ue.Msg, "Composer could not find") {
			// Factory::createComposer throws InvalidArgumentException
			return nil, &PHPException{Class: "InvalidArgumentException", Message: ue.Msg}
		}
		return nil, err
	}
	if pm.instances == nil {
		pm.instances = map[int64]*Composer{}
	}
	pm.instSeq++
	id := pm.instSeq
	pm.instances[id] = nc
	out := phpjson.NewMap()
	out.Set("inst", id)
	out.Set("root", rootDumpOf(nc))
	out.Set("global", nc.isGlobal())
	return out, nil
}

// installerRun is Installer::run for an installer configured by plugin code
// (Installer::create, or a clone of the running one like symfony/flex does).
func (pm *PluginManager) installerRun(a *phpjson.Map) (phpjson.Value, error) {
	c, err := pm.instance(a)
	if err != nil {
		return nil, err
	}
	in := NewInstaller(c)
	if base, _ := mustGet(a, "base").(int64); base != 0 {
		installerMu.Lock()
		ri := installerByID[int(base)]
		installerMu.Unlock()
		if ri != nil {
			src := ri.in
			in.PreferSource, in.PreferDist = src.PreferSource, src.PreferDist
			in.OptimizeAutoloader, in.ClassMapAuthoritative = src.OptimizeAutoloader, src.ClassMapAuthoritative
			in.ApcuAutoloader, in.ApcuPrefix, in.StrictPsr = src.ApcuAutoloader, src.ApcuPrefix, src.StrictPsr
			in.DevMode, in.DryRun, in.DownloadOnly = src.DevMode, src.DryRun, src.DownloadOnly
			in.Update, in.Install, in.DumpAutoloader, in.RunScripts = src.Update, src.Install, src.DumpAutoloader, src.RunScripts
			in.PreferStable, in.PreferLowest, in.MinimalUpdate = src.PreferStable, src.PreferLowest, src.MinimalUpdate
			in.WriteLock, in.ExecuteOperations = src.WriteLock, src.ExecuteOperations
			in.Audit, in.ErrorOnAudit, in.AuditFormat = src.Audit, src.ErrorOnAudit, src.AuditFormat
			in.UpdateMirrors, in.UpdateAllowList, in.UpdateAllowListSet, in.UpdateAllowTransitive = src.UpdateMirrors, src.UpdateAllowList, src.UpdateAllowListSet, src.UpdateAllowTransitive
			in.PlatformFilter, in.TemporaryConstraints = src.PlatformFilter, src.TemporaryConstraints
			in.IgnoredTypes, in.AllowedTypes, in.PolicyConfig = src.IgnoredTypes, src.AllowedTypes, src.PolicyConfig
			in.IM.SetPreferSource(src.PreferSource)
			in.IM.SetPreferDist(src.PreferDist)
			in.IM.OutputProgress = src.IM.OutputProgress
		}
	}
	opts := phpjson.ToMap(mustGet(a, "opts"))
	if opts == nil {
		opts = phpjson.NewMap()
	}
	b := func(k string, dst *bool) {
		if v, ok := opts.Get(k); ok {
			*dst = v == true
		}
	}
	b("dryRun", &in.DryRun)
	b("downloadOnly", &in.DownloadOnly)
	b("preferSource", &in.PreferSource)
	b("preferDist", &in.PreferDist)
	b("optimizeAutoloader", &in.OptimizeAutoloader)
	b("classMapAuthoritative", &in.ClassMapAuthoritative)
	b("apcuAutoloader", &in.ApcuAutoloader)
	b("strictPsr", &in.StrictPsr)
	b("update", &in.Update)
	b("install", &in.Install)
	b("devMode", &in.DevMode)
	b("dumpAutoloader", &in.DumpAutoloader)
	b("runScripts", &in.RunScripts)
	b("updateMirrors", &in.UpdateMirrors)
	b("preferStable", &in.PreferStable)
	b("preferLowest", &in.PreferLowest)
	b("minimalUpdate", &in.MinimalUpdate)
	b("writeLock", &in.WriteLock)
	b("executeOperations", &in.ExecuteOperations)
	b("audit", &in.Audit)
	b("errorOnAudit", &in.ErrorOnAudit)
	in.IM.SetPreferSource(in.PreferSource)
	in.IM.SetPreferDist(in.PreferDist)
	if s, ok := opts.Get("apcuPrefix"); ok {
		in.ApcuPrefix, _ = s.(string)
	}
	if s, ok := opts.Get("auditFormat"); ok {
		in.AuditFormat, _ = s.(string)
	}
	if v, ok := opts.Get("updateAllowList"); ok {
		in.UpdateAllowList = nil
		phpjson.Each(v, func(_ string, x phpjson.Value) {
			if s, ok := x.(string); ok {
				in.UpdateAllowList = append(in.UpdateAllowList, s)
			}
		})
		in.UpdateAllowListSet = true
	}
	if v, ok := opts.Get("updateAllowTransitive"); ok {
		n, _ := v.(int64)
		in.UpdateAllowTransitive = int(n)
	}
	if v, ok := opts.Get("platformFilter"); ok {
		switch t := v.(type) {
		case bool:
			if t {
				in.PlatformFilter = solver.IgnoreAll{}
			} else {
				in.PlatformFilter = solver.IgnoreNothing{}
			}
		default:
			var reqs []string
			phpjson.Each(t, func(_ string, x phpjson.Value) {
				if s, ok := x.(string); ok {
					reqs = append(reqs, s)
				}
			})
			if len(reqs) != 1 || reqs[0] != baseFilterToken {
				in.PlatformFilter = solver.NewIgnoreList(reqs)
			}
		}
	}
	strList := func(v phpjson.Value) []string {
		var out []string
		phpjson.Each(v, func(_ string, x phpjson.Value) {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		})
		return out
	}
	if v, ok := opts.Get("ignoredTypes"); ok {
		in.IgnoredTypes = strList(v)
	}
	if v, ok := opts.Get("allowedTypes"); ok {
		if v == nil {
			in.AllowedTypes = nil
		} else {
			in.AllowedTypes = strList(v)
			if in.AllowedTypes == nil {
				in.AllowedTypes = []string{}
			}
		}
	}
	if v, ok := opts.Get("temporaryConstraints"); ok {
		in.TemporaryConstraints = map[string]semver.Constraint{}
		var perr error
		phpjson.Each(v, func(name string, x phpjson.Value) {
			c, err := parseConstraintArg(x)
			if err != nil {
				perr = err
				return
			}
			if c == nil {
				c = &semver.MatchAll{}
			}
			in.TemporaryConstraints[strings.ToLower(name)] = c
		})
		if perr != nil {
			return nil, perr
		}
	}
	if opts.Has("silentSuggestions") && mustGet(opts, "silentSuggestions") == true {
		in.suggestions.io = &cio.IO{Out: discard{}, Err: discard{}, Verbosity: cio.Quiet}
	}
	code, err := in.Run()
	if err != nil {
		return nil, err
	}
	return int64(code), nil
}

// downloaderOp performs an operation of the runtime's downloaders (the
// Composer downloader classes the runtime objects stand for). async reports
// whether Composer performs the step in a background process (unzip of a
// zip archive, rm -rf of a package directory), whose promise settles later.
func (pm *PluginManager) downloaderOp(a *phpjson.Map) (async bool, err error) {
	c := pm.c
	im := pm.curIM
	if im == nil || im.c != c {
		// one manager for the calls outside of an install (download, then
		// install use the same temporary file)
		if pm.pluginIM == nil {
			pm.pluginIM = NewInstallationManager(c)
		}
		im = pm.pluginIM
		if pm.outputProgress != nil {
			im.OutputProgress = *pm.outputProgress
		}
	}
	h := pm.host
	p, err := h.pkgFromRef(mustGet(a, "p"))
	if err != nil {
		return false, err
	}
	var prev *pkg.Package
	if v := mustGet(a, "i"); v != nil {
		if prev, err = h.pkgFromRef(v); err != nil {
			return false, err
		}
	}
	if src := a.Str("src"); src != "" {
		p.SetInstallationSource(src)
	}
	if prev != nil {
		if src := a.Str("isrc"); src != "" {
			prev.SetInstallationSource(src)
		}
	}
	kind := a.Str("kind")
	path := a.Str("path")
	output := mustGet(a, "output") != false
	switch kind {
	case "svn", "hg", "fossil", "perforce":
		return false, unsupportedVcsDownloader(p, kind)
	case "rar", "xz":
		return false, &UserError{Msg: "The " + kind + " archive format is not supported by gocomposer", Code: 1}
	}
	if path != "" && a.Str("m") != "remove" {
		if want := strings.TrimRight(path, "/"); want != "" && fsutilNormalize(im.InstallPath(p)) != fsutilNormalize(want) {
			return false, &UserError{Msg: "gocomposer can only download " + p.PrettyName + " into its install path (" + im.InstallPath(p) + "), not " + path, Code: 1}
		}
	}
	existingDir := func(q *pkg.Package) bool {
		st, err := os.Lstat(im.packageBasePath(q))
		return err == nil && st.IsDir() && st.Mode()&os.ModeSymlink == 0
	}
	switch a.Str("m") {
	case "download":
		var transfer func() error
		switch kind {
		case "zip", "tar", "gzip", "file", "phar":
			transfer, err = im.fileDownload(p)
		case "path":
			err = im.pathDownload(p)
		case "git":
			err = im.gitDL().download(p, prev)
		}
		if err != nil {
			return false, err
		}
		if transfer != nil {
			if err := transfer(); err != nil {
				return false, err
			}
			return false, im.flushPluginPostDownloads()
		}
		return false, nil
	case "prepare":
		if kind == "git" {
			return false, im.gitDL().prepare(a.Str("op"), p, prev, im.InstallPath(p))
		}
		return false, nil
	case "cleanup":
		if kind == "git" {
			return false, im.gitDL().cleanup(a.Str("op"), p, prev, im.InstallPath(p))
		}
		im.cleanup(p)
		return false, nil
	case "install":
		async = kind == "zip" && unzipAvailable()
		return async, im.installCodeAs(kind, p, im.InstallPath(p), output)
	case "update":
		async = kind != "path" && kind != "git" && existingDir(prev) || kind == "zip" && unzipAvailable()
		return async, im.updateCodeAs(kind, prev, p, im.InstallPath(p))
	case "remove":
		async = kind != "path" && existingDir(p)
		return async, im.removeCode(p, output)
	}
	return false, &UserError{Msg: "Unknown downloader operation " + a.Str("m"), Code: 1}
}

func fsutilNormalize(p string) string { return filepath.Clean(p) }

// installCodeAs is DownloaderInterface::install for the downloader kind.
func (im *InstallationManager) installCodeAs(kind string, p *pkg.Package, target string, output bool) error {
	switch kind {
	case "path":
		return im.pathInstall(p, target, output)
	case "zip", "tar", "gzip":
		return im.archiveInstallAs(kind, p, target, output)
	case "file", "phar":
		return im.fileInstall(p, target, output)
	case "git":
		return im.gitDL().install(p, target)
	}
	return &UserError{Msg: "Unknown downloader type: " + kind, Code: 1}
}

// updateCodeAs is DownloaderInterface::update for the downloader kind (the
// same-type branch of DownloadManager::update).
func (im *InstallationManager) updateCodeAs(kind string, initial, target *pkg.Package, targetPath string) error {
	switch kind {
	case "git":
		return im.gitDL().update(initial, target, targetPath)
	case "path":
		im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUpdate, Initial: initial, Target: target}).Show(false) + im.pathAppendix(target))
		if err := im.removeCode(initial, false); err != nil {
			return err
		}
		return im.pathInstall(target, targetPath, false)
	}
	appendix := ""
	switch kind {
	case "zip", "tar", "gzip":
		appendix = ": Extracting archive"
	}
	im.c.IO.WriteError("  - " + (&solver.Operation{Type: solver.OpUpdate, Initial: initial, Target: target}).Show(false) + appendix)
	if err := im.removeCode(initial, false); err != nil {
		return err
	}
	return im.installCodeAs(kind, target, targetPath, false)
}

// autoloadDump is AutoloadGenerator::dump called by plugin code.
func (pm *PluginManager) autoloadDump(a *phpjson.Map) (phpjson.Value, error) {
	c, err := pm.instance(a)
	if err != nil {
		return nil, err
	}
	var dev *bool
	if v, ok := mustGet(a, "dev").(bool); ok {
		dev = &v
	}
	o := AutoloadOptions{
		DevMode:               dev,
		Optimize:              mustGet(a, "optimize") == true,
		ClassMapAuthoritative: mustGet(a, "authoritative") == true,
		Apcu:                  mustGet(a, "apcu") == true,
		RunScripts:            mustGet(a, "runScripts") == true,
		DryRun:                mustGet(a, "dryRun") == true,
		StrictAmbiguous:       mustGet(a, "strictAmbiguous") == true,
	}
	o.ApcuPrefix, _ = mustGet(a, "apcuPrefix").(string)
	o.Suffix, _ = mustGet(a, "suffix").(string)
	if mustGet(a, "platformFilter") == true {
		o.PlatformFilter = solver.IgnoreAll{}
	}
	repo := c.LocalRepo
	if c == pm.c && pm.curRepo != nil {
		repo = pm.curRepo
	}
	d := NewEventDispatcher(c)
	d.RunScripts = !c.DisableScripts
	res, err := RunAutoloadDump(c, d, repo, o)
	if err != nil {
		return nil, err
	}
	out := phpjson.NewMap()
	cm := phpjson.NewMap()
	if res != nil && res.ClassMap != nil {
		for _, class := range res.ClassMap.SortedClasses() {
			cm.Set(class, res.ClassMap.Map[class])
		}
	}
	out.Set("map", cm)
	return out, nil
}

// localRepoOp applies InstalledFilesystemRepository changes made by plugin
// installers (addPackage/removePackage/write).
func (pm *PluginManager) localRepoOp(a *phpjson.Map) error {
	c, err := pm.instance(a)
	if err != nil {
		return err
	}
	r := c.LocalRepo
	if c == pm.c && pm.curRepo != nil {
		r = pm.curRepo
	}
	switch a.Str("m") {
	case "add", "remove":
		p, err := pm.host.pkgFromRef(mustGet(a, "p"))
		if err != nil {
			return err
		}
		if a.Str("m") == "add" {
			if !r.HasPackage(p) {
				r.AddPackage(p.Clone())
			}
		} else {
			r.RemovePackage(p)
		}
		return nil
	case "write":
		if v, ok := a.Get("devNames"); ok {
			var names []string
			phpjson.Each(v, func(_ string, x phpjson.Value) {
				if s, ok := x.(string); ok {
					names = append(names, s)
				}
			})
			r.DevPackageNames = names
		}
		im := NewInstallationManager(c)
		return r.Write(mustGet(a, "dev") == true, im.InstallPath, c.Package, c.Cwd)
	}
	return &UserError{Msg: "Unknown local repository operation " + a.Str("m"), Code: 1}
}
