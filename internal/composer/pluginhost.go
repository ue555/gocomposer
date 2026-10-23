package composer

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/res"
)

// pluginHost is the PHP process running the code of Composer plugins
// (internal/res/pluginrt/host.php). gocomposer keeps doing all of Composer's
// work in Go; the host only executes plugin classes (activate/deactivate/
// uninstall, event subscribers, command providers, custom installers) with
// a Composer API stub that calls back into Go. One host lives as long as the
// Composer instance whose plugins it runs, so plugin state survives between
// calls like inside Composer.
type pluginHost struct {
	pm      *PluginManager
	cmd     *exec.Cmd
	w       io.WriteCloser
	r       *bufio.Reader
	rf      *os.File
	dir     string
	dead    bool
	waitErr chan error

	ids    map[*pkg.Package]int
	byID   map[int]*pkg.Package
	sent   map[int]bool
	nextID int

	rootJSON string
	env      map[string]string
	// installersSent are the gocomposer Installers the runtime has objects for
	installersSent map[int]bool
}

// PHPException is an exception thrown by plugin code (rethrown like Composer
// does, rendered as an error box).
type PHPException struct {
	Class, Message, File string
	Code, Line           int
	Prev                 *PHPException
}

func (e *PHPException) Error() string { return e.Message }

// Unwrap exposes a TransportException thrown by plugin code as a transport
// error: Composer's Application exits with Installer::ERROR_TRANSPORT_EXCEPTION.
func (e *PHPException) Unwrap() error {
	if e.Class == "Composer\\Downloader\\TransportException" {
		return &httpdl.TransportError{Msg: e.Message, StatusCode: e.Code}
	}
	return nil
}

// Previous is the PHP previous exception (rendered as a further box).
func (e *PHPException) Previous() error {
	if e.Prev == nil {
		return nil
	}
	return e.Prev
}

func phpException(v phpjson.Value) *PHPException {
	m := phpjson.ToMap(v)
	if m == nil {
		return &PHPException{Class: "RuntimeException", Message: "unknown plugin error"}
	}
	e := &PHPException{Class: m.Str("class"), Message: m.Str("message"), File: m.Str("file")}
	if n, ok := mustGet(m, "code").(int64); ok {
		e.Code = int(n)
	}
	if n, ok := mustGet(m, "line").(int64); ok {
		e.Line = int(n)
	}
	if p := mustGet(m, "previous"); p != nil {
		e.Prev = phpException(p)
	}
	return e
}

func startPluginHost(pm *PluginManager) (*pluginHost, error) {
	c := pm.c
	// the PHP binary is enough (probing it may fail for unrelated reasons)
	php, _, _ := c.PHP()
	if php == "" {
		return nil, &UserError{Msg: "Composer plugins are executed with PHP, but no PHP binary was found (set GOCOMPOSER_PHP or add php to PATH). Run the command with --no-plugins to skip plugins.", Code: 1}
	}
	dir, err := os.MkdirTemp("", "gocomposer-plugins-")
	if err != nil {
		return nil, err
	}
	if err := writePluginRuntime(dir); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	reqR, reqW, err := os.Pipe()
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	respR, respW, err := os.Pipe()
	if err != nil {
		reqR.Close()
		reqW.Close()
		os.RemoveAll(dir)
		return nil, err
	}
	cmd := exec.Command(php, "-d", "memory_limit=1536M", filepath.Join(dir, "host.php"))
	cmd.Dir = c.Cwd
	env := os.Environ()
	if exe, err := os.Executable(); err == nil {
		env = append(env, "COMPOSER_BINARY="+exe)
	}
	cmd.Env = env
	cmd.ExtraFiles = []*os.File{reqR, respW}
	cmd.Stdout = c.IO.Out
	cmd.Stderr = c.IO.Err
	if c.IO.Interactive {
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Start(); err != nil {
		reqR.Close()
		reqW.Close()
		respR.Close()
		respW.Close()
		os.RemoveAll(dir)
		return nil, &UserError{Msg: "Could not start PHP to run Composer plugins: " + err.Error(), Code: 1}
	}
	reqR.Close()
	respW.Close()
	h := &pluginHost{pm: pm, cmd: cmd, w: reqW, r: bufio.NewReaderSize(respR, 1<<20), rf: respR, dir: dir,
		ids: map[*pkg.Package]int{}, byID: map[int]*pkg.Package{}, sent: map[int]bool{}, nextID: 1, waitErr: make(chan error, 1)}
	h.env = environMap()
	h.installersSent = map[int]bool{}
	go func() { h.waitErr <- cmd.Wait() }()
	return h, nil
}

// writePluginRuntime extracts the embedded PHP runtime (and Composer's
// ClassLoader/InstalledVersions) into dir.
func writePluginRuntime(dir string) error {
	err := fs.WalkDir(res.PluginRuntime, "pluginrt", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(p, "pluginrt"), "/")
		target := filepath.Join(dir, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		data, err := res.PluginRuntime.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o600)
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "composer", "Autoload"), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "composer", "Autoload", "ClassLoader.php"), res.ClassLoader, 0o600); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "composer", "InstalledVersions.php"), res.InstalledVersions, 0o600)
}

func environMap() map[string]string {
	out := map[string]string{}
	for _, kv := range os.Environ() {
		if k, v, ok := strings.Cut(kv, "="); ok {
			out[k] = v
		}
	}
	return out
}

// envDiff returns the environment changes since the last call (scripts run
// by gocomposer use @putenv and COMPOSER_DEV_MODE like Composer's process).
func (h *pluginHost) envDiff() *phpjson.Map {
	cur := environMap()
	set := phpjson.NewMap()
	var keys []string
	for k, v := range cur {
		if old, ok := h.env[k]; !ok || old != v {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		set.Set(k, cur[k])
	}
	var unset []phpjson.Value
	var gone []string
	for k := range h.env {
		if _, ok := cur[k]; !ok {
			gone = append(gone, k)
		}
	}
	sort.Strings(gone)
	for _, k := range gone {
		unset = append(unset, k)
	}
	h.env = cur
	if set.Len() == 0 && len(unset) == 0 {
		return nil
	}
	m := phpjson.NewMap()
	m.Set("set", set)
	if unset == nil {
		unset = []phpjson.Value{}
	}
	m.Set("unset", unset)
	return m
}

func (h *pluginHost) applyEnv(v phpjson.Value) {
	m := phpjson.ToMap(v)
	if m == nil {
		return
	}
	if set := m.Map("set"); set != nil {
		for _, k := range set.Keys() {
			s, _ := phpjson.ScalarString(mustGet(set, k))
			os.Setenv(k, s)
			h.env[k] = s
		}
	}
	phpjson.Each(mustGet(m, "unset"), func(_ string, k phpjson.Value) {
		if s, ok := k.(string); ok {
			os.Unsetenv(s)
			delete(h.env, s)
		}
	})
}

func (h *pluginHost) send(m *phpjson.Map) error {
	s, err := phpjson.Encode(m, 0)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(h.w, s+"\n"); err != nil {
		return h.died()
	}
	return nil
}

func (h *pluginHost) recv() (*phpjson.Map, error) {
	line, err := h.r.ReadBytes('\n')
	if err != nil {
		return nil, h.died()
	}
	v, err := phpjson.Decode(line)
	if err != nil {
		return nil, fmt.Errorf("gocomposer plugin host: invalid message: %v", err)
	}
	m := phpjson.ToMap(v)
	if m == nil {
		return nil, fmt.Errorf("gocomposer plugin host: invalid message")
	}
	return m, nil
}

// died handles the PHP process exiting under a call (fatal error, exit() in
// plugin code): the command ends with its exit code like Composer's process
// would, without further output.
func (h *pluginHost) died() error {
	h.dead = true
	code := 255
	if err := <-h.waitErr; err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		}
	} else {
		code = 0
	}
	h.waitErr <- nil
	return &ScriptError{Msg: "", Code: code}
}

// call runs op in the host and serves its callbacks until it returns.
func (h *pluginHost) call(op string, args *phpjson.Map) (phpjson.Value, error) {
	if h.dead {
		return nil, &UserError{Msg: "The gocomposer plugin host is not running", Code: 1}
	}
	if args == nil {
		args = phpjson.NewMap()
	}
	msg := phpjson.NewMap()
	msg.Set("t", "call")
	msg.Set("op", op)
	msg.Set("a", args)
	if root := h.rootData(); root != nil {
		msg.Set("root", root)
	}
	if env := h.envDiff(); env != nil {
		msg.Set("env", env)
	}
	if err := h.send(msg); err != nil {
		return nil, err
	}
	for {
		m, err := h.recv()
		if err != nil {
			return nil, err
		}
		if env, ok := m.Get("env"); ok {
			h.applyEnv(env)
		}
		switch m.Str("t") {
		case "ret", "exc":
			if root, ok := m.Get("root"); ok {
				h.pm.applyRootChanges(phpjson.ToMap(root))
			}
			if m.Str("t") == "exc" {
				return nil, phpException(mustGet(m, "e"))
			}
			v, _ := m.Get("v")
			return v, nil
		case "cb":
			v, cbErr := h.pm.callback(m.Str("op"), phpjson.ToMap(mustGet(m, "a")))
			reply := phpjson.NewMap()
			if cbErr != nil {
				if se, ok := cbErr.(*ScriptError); ok && h.dead {
					return nil, se
				}
				reply.Set("t", "cbexc")
				e := phpjson.NewMap()
				class := "RuntimeException"
				var pe *PHPException
				if asPHPException(cbErr, &pe) {
					class = pe.Class
				}
				if _, ok := cbErr.(*FactoryInvalidArgumentError); ok {
					class = "InvalidArgumentException"
				}
				e.Set("class", class)
				e.Set("message", cbErr.Error())
				code := int64(0)
				if ue, ok := cbErr.(*UserError); ok {
					code = int64(ue.Code)
				}
				if se, ok := cbErr.(*ScriptError); ok {
					class = "Composer\\Script\\ScriptExecutionException"
					e.Set("class", class)
					code = int64(se.Code)
				}
				e.Set("code", code)
				reply.Set("e", e)
				// remember the Go error so it surfaces unchanged when the
				// plugin does not catch it
				h.pm.lastCallbackErr = cbErr
			} else {
				reply.Set("t", "cbret")
				reply.Set("v", v)
			}
			if env := h.envDiff(); env != nil {
				reply.Set("env", env)
			}
			if err := h.send(reply); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("gocomposer plugin host: unexpected message %q", m.Str("t"))
		}
	}
}

func asPHPException(err error, target **PHPException) bool {
	pe, ok := err.(*PHPException)
	if ok {
		*target = pe
	}
	return ok
}

// rootData returns the root package data when it changed since it was last
// sent.
func (h *pluginHost) rootData() *phpjson.Map {
	d := h.pm.rootDump()
	s, _ := phpjson.Encode(d, 0)
	if s == h.rootJSON {
		return nil
	}
	h.rootJSON = s
	return d
}

// idOf returns the id of a Go package in the host (0 for the root package).
func (h *pluginHost) idOf(p *pkg.Package) int {
	if p == h.pm.c.Root() || p == h.pm.c.Package {
		return 0
	}
	id, ok := h.ids[p]
	if !ok {
		id = h.nextID
		h.nextID++
		h.ids[p] = id
		h.byID[id] = p
	}
	return id
}

// ref is the reference of a Go package for the host (its data is sent once).
func (h *pluginHost) ref(p *pkg.Package) phpjson.Value {
	if p == nil {
		return nil
	}
	m := phpjson.NewMap()
	id := h.idOf(p)
	m.Set("id", int64(id))
	if id == 0 {
		return m
	}
	if !h.sent[id] {
		h.sent[id] = true
		if p.AliasOf != nil {
			d := phpjson.NewMap()
			d.Set("_alias_of", h.ref(p.AliasOf))
			d.Set("version", p.PrettyVersion())
			d.Set("version_normalized", p.Version())
			d.Set("_root_alias", p.RootPackageAlias)
			m.Set("d", d)
		} else {
			m.Set("d", pkg.Dump(p))
		}
	}
	return m
}

// pkgFromRef resolves a package reference coming from the host.
func (h *pluginHost) pkgFromRef(v phpjson.Value) (*pkg.Package, error) {
	m := phpjson.ToMap(v)
	if m == nil {
		return nil, nil
	}
	if id, ok := mustGet(m, "id").(int64); ok {
		if id == 0 {
			return h.pm.c.Root(), nil
		}
		if p, ok := h.byID[int(id)]; ok {
			return p, nil
		}
		return nil, fmt.Errorf("unknown package reference %d", id)
	}
	d := m.Map("d")
	if d == nil {
		return nil, fmt.Errorf("invalid package reference")
	}
	p, err := (&pkg.Loader{}).Load(d)
	if err != nil {
		return nil, err
	}
	if a := m.Map("alias"); a != nil {
		p = pkg.NewAlias(p, a.Str("version_normalized"), a.Str("version"))
	}
	return p, nil
}

// close shuts the host down and removes its runtime directory.
func (h *pluginHost) close() {
	if h == nil {
		return
	}
	if !h.dead {
		msg := phpjson.NewMap()
		msg.Set("t", "call")
		msg.Set("op", "shutdown")
		msg.Set("a", phpjson.NewMap())
		_ = h.send(msg)
		h.w.Close()
		<-h.waitErr
		h.dead = true
	} else {
		h.w.Close()
	}
	h.rf.Close()
	os.RemoveAll(h.dir)
}
