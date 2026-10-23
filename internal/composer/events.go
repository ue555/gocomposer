package composer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gocomposer/internal/autoload"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/res"
	"gocomposer/internal/solver"
)

// Script events (ScriptEvents::*).
const (
	PreInstallCmd    = "pre-install-cmd"
	PostInstallCmd   = "post-install-cmd"
	PreUpdateCmd     = "pre-update-cmd"
	PostUpdateCmd    = "post-update-cmd"
	PreAutoloadDump  = "pre-autoload-dump"
	PostAutoloadDump = "post-autoload-dump"
)

// ScriptError mirrors ScriptExecutionException.
type ScriptError struct {
	Msg  string
	Code int
}

func (e *ScriptError) Error() string { return e.Msg }

// CallbackError is an exception thrown by a PHP script callback, rethrown
// by the event dispatcher (rendered as an error box).
type CallbackError struct {
	Class, Msg string
}

func (e *CallbackError) Error() string { return e.Msg }

// EventDispatcher runs scripts defined by the root package.
type EventDispatcher struct {
	c          *Composer
	RunScripts bool
	stack      []string
	skip       map[string]bool
	// Binary is the path written to COMPOSER_BINARY for @composer calls.
	Binary string
	// Timeout is ProcessExecutor::$timeout for script processes (0 = none);
	// it defaults to the process-timeout config like Factory::createComposer.
	Timeout time.Duration
	// added holds listeners registered with AddListener (not subject to
	// RunScripts, like EventDispatcher::$listeners).
	added map[string][]string
}

// AddListener is EventDispatcher::addListener for command string listeners.
func (d *EventDispatcher) AddListener(name, listener string) {
	if d.added == nil {
		d.added = map[string][]string{}
	}
	d.added[name] = append(d.added[name], listener)
}

// ProcessTimeoutError mirrors Symfony's ProcessTimedOutException.
type ProcessTimeoutError struct {
	Command string
	Timeout time.Duration
}

// ExceptionHints are printed by Application::hintCommonErrors.
func (e *ProcessTimeoutError) ExceptionHints() []string {
	return []string{"<error>The following exception is caused by a process timeout</error>", "<error>Check https://getcomposer.org/doc/06-config.md#process-timeout for details</error>"}
}

func (e *ProcessTimeoutError) Error() string {
	return fmt.Sprintf("The process \"%s\" exceeded the timeout of %s seconds.", e.Command, strconv.FormatFloat(e.Timeout.Seconds(), 'f', -1, 64))
}

// NewEventDispatcher creates the dispatcher.
func NewEventDispatcher(c *Composer) *EventDispatcher {
	d := &EventDispatcher{c: c, RunScripts: true, skip: map[string]bool{}}
	if v := os.Getenv("COMPOSER_SKIP_SCRIPTS"); v != "" {
		for _, s := range strings.Split(v, ",") {
			d.skip[strings.TrimSpace(s)] = true
		}
	}
	if exe, err := os.Executable(); err == nil {
		d.Binary = exe
	}
	if c != nil && c.Config != nil {
		if n, ok := c.Config.Get("process-timeout").(int64); ok && n > 0 {
			d.Timeout = time.Duration(n) * time.Second
		}
	}
	return d
}

// ScriptEvent carries event data.
type ScriptEvent struct {
	Name      string
	DevMode   bool
	Args      []string
	Flags     map[string]any
	Operation *solver.Operation
	// Kind is the Composer event class for plugin listeners: "" (Script\Event),
	// "plain" (EventDispatcher\Event), "package" (PackageEvent), "installer"
	// (InstallerEvent), "command" (CommandEvent), "precommand"
	// (PreCommandRunEvent).
	Kind string
	// Ops are the operations of package/installer events.
	Ops    []*solver.Operation
	opsKey string
	// Exec is InstallerEvent::isExecutingOperations.
	Exec bool
	// Command and Input describe CommandEvent/PreCommandRunEvent.
	Command string
	Input   *CommandInput
	// eid identifies the event object in the plugin runtime.
	eid    string
	origin *ScriptEvent
}

func (d *EventDispatcher) listeners(name string) []string {
	out := append([]string(nil), d.added[name]...)
	return append(out, d.scriptListeners(name)...)
}

// listenerEntry is a script/command listener or a plugin listener.
type listenerEntry struct {
	script string
	plugin *pluginListener
}

// collectListeners is EventDispatcher::getListeners: plugin listeners by
// priority (highest first), with the root package scripts appended to the
// listeners of priority 0.
func (d *EventDispatcher) collectListeners(name string) ([]listenerEntry, error) {
	var pls []pluginListener
	if pm := pluginManagerOf(d.c); pm.active() {
		var err error
		if pls, err = pm.listeners(name); err != nil {
			return nil, err
		}
	}
	var out []listenerEntry
	add := func(l pluginListener) {
		l2 := l
		if l.String != "" {
			out = append(out, listenerEntry{script: l.String})
			return
		}
		out = append(out, listenerEntry{plugin: &l2})
	}
	i := 0
	for ; i < len(pls) && pls[i].Prio >= 0; i++ {
		add(pls[i])
	}
	for _, s := range d.listeners(name) {
		out = append(out, listenerEntry{script: s})
	}
	for ; i < len(pls); i++ {
		add(pls[i])
	}
	return out, nil
}

func (d *EventDispatcher) scriptListeners(name string) []string {
	if !d.RunScripts {
		return nil
	}
	scripts := d.c.Root().Scripts()
	if scripts == nil {
		return nil
	}
	v, ok := scripts.Get(name)
	if !ok || !phpjson.Truthy(v) {
		return nil
	}
	if d.skip[name] {
		d.c.IO.WriteErrorV("Skipped script listeners for <info>"+name+"</info> because of COMPOSER_SKIP_SCRIPTS", true, 64)
		return nil
	}
	var out []string
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		s, _ := phpjson.ScalarString(e)
		out = append(out, s)
	})
	return out
}

// HasListeners reports whether a script or plugin listener is defined for name.
func (d *EventDispatcher) HasListeners(name string) bool {
	l, err := d.collectListeners(name)
	return err == nil && len(l) > 0
}

// DispatchScript is EventDispatcher::dispatchScript.
func (d *EventDispatcher) DispatchScript(name string, devMode bool, args []string, flags map[string]any) (int, error) {
	return d.dispatch(&ScriptEvent{Name: name, DevMode: devMode, Args: args, Flags: flags})
}

// DispatchPackageEvent runs pre/post-package-* scripts and listeners.
func (d *EventDispatcher) DispatchPackageEvent(name string, devMode bool, op *solver.Operation) error {
	ev := &ScriptEvent{Name: name, DevMode: devMode, Operation: op, Kind: "package"}
	if pm := pluginManagerOf(d.c); pm != nil {
		ev.Ops, ev.opsKey = pm.curOps, pm.curOpsKey
	}
	_, err := d.dispatch(ev)
	return err
}

// DispatchInstallerEvent is EventDispatcher::dispatchInstallerEvent
// (pre-operations-exec).
func (d *EventDispatcher) DispatchInstallerEvent(name string, devMode, executeOperations bool, ops []*solver.Operation) error {
	if !d.HasListeners(name) && os.Getenv("COMPOSER_DEBUG_EVENTS") == "" {
		return nil
	}
	_, err := d.dispatch(&ScriptEvent{Name: name, DevMode: devMode, Kind: "installer", Exec: executeOperations, Ops: ops})
	return err
}

// DispatchCommandEvent dispatches PluginEvents::COMMAND for a command.
func (d *EventDispatcher) DispatchCommandEvent(command string, input *CommandInput) error {
	if pm := pluginManagerOf(d.c); !pm.active() && !d.HasListeners("command") {
		return nil
	}
	_, err := d.dispatch(&ScriptEvent{Name: "command", Kind: "command", Command: command, Input: input})
	return err
}

// DispatchPreCommandRun dispatches PluginEvents::PRE_COMMAND_RUN.
func (d *EventDispatcher) DispatchPreCommandRun(command string, input *CommandInput) error {
	if pm := pluginManagerOf(d.c); !pm.active() && !d.HasListeners("pre-command-run") {
		return nil
	}
	_, err := d.dispatch(&ScriptEvent{Name: "pre-command-run", Kind: "precommand", Command: command, Input: input})
	return err
}

func isPhpScript(s string) bool { return !strings.Contains(s, " ") && strings.Contains(s, "::") }
func isCommandClass(s string) bool {
	return strings.Contains(s, `\`) && !strings.Contains(s, " ") && strings.HasSuffix(s, "Command")
}
func isComposerScript(s string) bool {
	return strings.HasPrefix(s, "@") && !strings.HasPrefix(s, "@php ") && !strings.HasPrefix(s, "@putenv ")
}

var reNoAdditionalArgs = regexp.MustCompile(` ?@no_additional_args`)

func (d *EventDispatcher) dispatch(ev *ScriptEvent) (int, error) {
	if os.Getenv("COMPOSER_DEBUG_EVENTS") != "" {
		details := ""
		switch ev.Kind {
		case "package":
			if ev.Operation != nil {
				details = ev.Operation.Show(false)
			}
		case "command", "precommand":
			details = ev.Command
		}
		if details != "" {
			d.c.IO.WriteError("Dispatching <info>" + ev.Name + "</info> (" + details + ") event")
		} else {
			d.c.IO.WriteError("Dispatching <info>" + ev.Name + "</info> event")
		}
	}
	entries, err := d.collectListeners(ev.Name)
	if err != nil {
		return 0, err
	}
	for _, s := range d.stack {
		if s == ev.Name {
			return 0, &UserError{Msg: fmt.Sprintf("Circular call to script handler '%s' detected", ev.Name), Code: 1}
		}
	}
	d.stack = append(d.stack, ev.Name)
	defer func() { d.stack = d.stack[:len(d.stack)-1] }()
	pm := pluginManagerOf(d.c)
	if pm != nil {
		pm.pushDispatcher(d)
		defer pm.popDispatcher()
		ownEvent := ev.eid == ""
		defer func() {
			if ownEvent {
				pm.freeEvent(ev)
			}
		}()
	}
	returnMax := 0
	for _, entry := range entries {
		if entry.plugin != nil {
			d.ensureBinDirIsInPath()
			ret, stopped, err := pm.callListener(d, *entry.plugin, ev)
			if err != nil {
				return 0, err
			}
			if ret > returnMax {
				returnMax = ret
			}
			if stopped {
				break
			}
			continue
		}
		callable := entry.script
		ret := 0
		d.ensureBinDirIsInPath()
		additional := ev.Args
		if strings.Contains(callable, "@no_additional_args") {
			callable = reNoAdditionalArgs.ReplaceAllString(callable, "")
			additional = nil
		}
		formatted := ev.Name
		if len(additional) > 0 {
			formatted += " (" + strings.Join(additional, ", ") + ")"
		}
		switch {
		case isComposerScript(callable):
			d.c.IO.WriteErrorV("> "+formatted+": "+callable, true, 64)
			parts := strings.Split(callable[1:], " ")
			scriptName := parts[0]
			script := parts[1:]
			idx := -1
			for i, s := range script {
				if s == "@additional_args" {
					idx = i
					break
				}
			}
			var args []string
			if idx >= 0 {
				args = append(append(append([]string(nil), script[:idx]...), additional...), script[idx:]...)
			} else {
				args = append(append([]string(nil), script...), additional...)
			}
			if strings.HasPrefix(callable, "@composer ") {
				execLine := d.phpExecCommand() + " " + escapeArg(d.Binary) + " " + strings.Join(args, " ")
				if strings.HasPrefix(callable, "@composer ") && !d.binaryIsPHP() {
					execLine = escapeArg(d.Binary) + " " + strings.Join(args, " ")
				}
				code, err := d.executeShell(execLine)
				if err != nil {
					return 0, err
				}
				if code != 0 {
					d.c.IO.WriteErrorV(fmt.Sprintf("<error>Script %s handling the %s event returned with error code %d</error>", callable, ev.Name, code), true, 16)
					return 0, &ScriptError{Msg: "Error Output: ", Code: code}
				}
			} else {
				if !d.HasListeners(scriptName) {
					d.c.IO.WriteErrorV("<warning>You made a reference to a non-existent script "+callable+"</warning>", true, 16)
				}
				r, err := d.dispatch(&ScriptEvent{Name: scriptName, DevMode: ev.DevMode, Args: args, Flags: ev.Flags, origin: ev})
				if err != nil {
					if _, ok := err.(*ScriptError); ok {
						d.c.IO.WriteErrorV(fmt.Sprintf("<error>Script %s was called via %s</error>", callable, ev.Name), true, 16)
					}
					return 0, err
				}
				ret = r
			}
		case callable == `Composer\Config::disableProcessTimeout`:
			// Config::disableProcessTimeout() only resets ProcessExecutor's
			// static timeout, which lives in this process for gocomposer
			if d.c.IO.IsVerbose() {
				d.c.IO.WriteError("> " + ev.Name + ": " + callable)
			} else if d.needsOutput(ev) {
				d.c.IO.WriteError("> " + callable)
			}
			d.Timeout = 0
		case isPhpScript(callable):
			className := callable[:strings.Index(callable, "::")]
			methodName := callable[strings.Index(callable, "::")+2:]
			if d.c.IO.IsVerbose() {
				// written by the runner
			}
			r, err := d.runPhpCallback(ev, className, methodName)
			if err != nil {
				return 0, err
			}
			ret = r
		case isCommandClass(callable):
			return 0, &UserError{Msg: "Script " + callable + " is a Symfony Console command class; running command classes from scripts is not supported by gocomposer", Code: 1}
		default:
			if callable == "" {
				continue
			}
			argStr := ""
			var quoted []string
			for _, a := range additional {
				quoted = append(quoted, escapeArg(a))
			}
			argStr = strings.Join(quoted, " ")
			execLine := callable
			if !strings.HasPrefix(callable, "@putenv ") {
				if strings.Contains(callable, "@additional_args") {
					execLine = strings.ReplaceAll(callable, "@additional_args", argStr)
				} else if argStr != "" {
					execLine = callable + " " + argStr
				}
			}
			if d.c.IO.IsVerbose() {
				d.c.IO.WriteError("> " + ev.Name + ": " + execLine)
			} else if d.needsOutput(ev) {
				d.c.IO.WriteError("> " + execLine)
			}
			for _, localExec := range d.c.Root().Binaries() {
				if regexp.MustCompile(`\b` + regexp.QuoteMeta(callable) + `$`).MatchString(localExec) {
					caller := determineBinaryCaller(localExec)
					execLine = regexp.MustCompile(`^`+regexp.QuoteMeta(callable)).ReplaceAllLiteralString(execLine, caller+" "+localExec)
					break
				}
			}
			if strings.HasPrefix(execLine, "@putenv ") {
				kv := execLine[8:]
				if k, v, ok := strings.Cut(kv, "="); ok {
					os.Setenv(k, v)
				} else {
					os.Unsetenv(kv)
				}
				continue
			}
			if strings.HasPrefix(execLine, "@php ") {
				pathAndArgs := execLine[5:]
				if m := regexp.MustCompile(`^[^'"\s/\\]+`).FindString(pathAndArgs); m != "" {
					if _, err := os.Stat(m); err != nil {
						if p, err := exec.LookPath(m); err == nil {
							pathAndArgs = p + pathAndArgs[len(m):]
						}
					}
				}
				execLine = d.phpExecCommand() + " " + pathAndArgs
			} else {
				if php, _, _ := d.c.PHP(); php != "" {
					os.Setenv("PHP_BINARY", php)
				}
			}
			if strings.HasPrefix(execLine, "composer ") {
				if d.binaryIsPHP() {
					execLine = d.phpExecCommand() + " " + escapeArg(d.Binary) + execLine[8:]
				} else {
					execLine = escapeArg(d.Binary) + execLine[8:]
				}
			}
			code, err := d.executeShell(execLine)
			if err != nil {
				return 0, err
			}
			if code != 0 {
				d.c.IO.WriteErrorV(fmt.Sprintf("<error>Script %s handling the %s event returned with error code %d</error>", callable, ev.Name, code), true, 16)
				return 0, &ScriptError{Msg: "Error Output: ", Code: code}
			}
		}
		if ret > returnMax {
			returnMax = ret
		}
	}
	return returnMax, nil
}

func (d *EventDispatcher) needsOutput(ev *ScriptEvent) bool {
	if ev.Name == "__exec_command" {
		return false
	}
	if ev.Flags != nil {
		if v, ok := ev.Flags["script-alias-input"]; ok && v != nil {
			return false
		}
	}
	return true
}

func (d *EventDispatcher) binaryIsPHP() bool { return strings.HasSuffix(d.Binary, ".phar") }

func escapeArg(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func (d *EventDispatcher) phpExecCommand() string {
	php, _, err := d.c.PHP()
	if err != nil || php == "" {
		php = "php"
	}
	return escapeArg(php) + " -d allow_url_fopen=" + escapeArg("1") + " -d disable_functions=" + escapeArg("") + " -d memory_limit=" + escapeArg("1536M")
}

func (d *EventDispatcher) ensureBinDirIsInPath() {
	binDir := d.c.Config.GetString("bin-dir")
	if st, err := os.Stat(binDir); err == nil && st.IsDir() {
		real := fsutil.Realpath(binDir)
		path := os.Getenv("PATH")
		re := regexp.MustCompile(`(^|:)` + regexp.QuoteMeta(real) + `($|:)`)
		if !re.MatchString(path) {
			os.Setenv("PATH", real+":"+path)
		}
	}
}

func (d *EventDispatcher) executeShell(line string) (int, error) {
	cmd := exec.Command("/bin/sh", "-c", line)
	// Symfony Process runs the command line through "sh -c"
	cmd.Args[0] = "sh"
	cmd.Dir = d.c.Cwd
	cmd.Env = os.Environ()
	if d.Binary != "" {
		cmd.Env = append(cmd.Env, "COMPOSER_BINARY="+d.Binary)
	}
	cmd.Stdout = d.c.IO.Out
	cmd.Stderr = d.c.IO.Err
	if d.c.IO.IsQuiet() {
		// ProcessExecutor::outputHandler writes through the IO, which drops
		// output when quiet
		cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	}
	if d.c.IO.Interactive {
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Start(); err != nil {
		return 1, nil
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var timeout <-chan time.Time
	if d.Timeout > 0 {
		t := time.NewTimer(d.Timeout)
		defer t.Stop()
		timeout = t.C
	}
	var err error
	select {
	case err = <-done:
	case <-timeout:
		// Process::stop(): SIGTERM, then SIGKILL after 10 seconds
		_ = cmd.Process.Signal(syscall.SIGTERM)
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
		return 1, &ProcessTimeoutError{Command: line, Timeout: d.Timeout}
	}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode(), nil
		}
		return 1, nil
	}
	return 0, nil
}

// runPhpCallback executes Class::method($event) with the configured PHP. The
// event object is a gocomposer stub implementing the Composer\Script\Event
// API subset scripts commonly use; unsupported calls fail explicitly.
func (d *EventDispatcher) runPhpCallback(ev *ScriptEvent, className, methodName string) (int, error) {
	php, _, err := d.c.PHP()
	if err != nil || php == "" {
		return 0, &UserError{Msg: "Script " + className + "::" + methodName + " requires a PHP binary to run (set GOCOMPOSER_PHP or add php to PATH)", Code: 1}
	}
	tmpDir, err := os.MkdirTemp("", "gocomposer-script")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tmpDir)
	loaderFile := filepath.Join(tmpDir, "loader.php")
	if err := d.writeScriptLoader(loaderFile, ev.DevMode); err != nil {
		return 0, err
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "ClassLoader.php"), res.ClassLoader, 0o600); err != nil {
		return 0, err
	}
	cfg := phpjson.NewMap()
	for _, k := range []string{"vendor-dir", "bin-dir", "cache-dir", "data-dir", "home", "process-timeout"} {
		cfg.Set(k, d.c.Config.Get(k))
	}
	data := phpjson.NewMap()
	data.Set("name", ev.Name)
	data.Set("dev", ev.DevMode)
	args := make([]phpjson.Value, len(ev.Args))
	for i, a := range ev.Args {
		args[i] = a
	}
	data.Set("args", args)
	data.Set("config", cfg)
	data.Set("class", className)
	data.Set("method", methodName)
	data.Set("loader", loaderFile)
	data.Set("verbosity", int64(d.c.IO.Verbosity))
	data.Set("decorated", d.c.IO.Decorated)
	data.Set("verbose", d.c.IO.IsVerbose())
	data.Set("output", d.needsOutput(ev))
	root := d.c.Root()
	rp := phpjson.NewMap()
	rp.Set("name", root.Name)
	rp.Set("pretty_name", root.PrettyName)
	rp.Set("version", root.Version())
	rp.Set("pretty_version", root.PrettyVersion())
	rp.Set("extra", nonNilMap(root.Extra()))
	rp.Set("type", root.Type())
	data.Set("root", rp)
	resultFile := filepath.Join(tmpDir, "exception.json")
	data.Set("result", resultFile)
	enc, _ := phpjson.Encode(data, 0)
	runner := filepath.Join(tmpDir, "runner.php")
	if err := os.WriteFile(runner, []byte(scriptRunnerPHP), 0o600); err != nil {
		return 0, err
	}
	cmd := exec.Command(php, "-d", "memory_limit=1536M", runner)
	cmd.Dir = d.c.Cwd
	cmd.Env = append(os.Environ(), "GOCOMPOSER_SCRIPT_EVENT="+enc)
	if d.Binary != "" {
		cmd.Env = append(cmd.Env, "COMPOSER_BINARY="+d.Binary)
	}
	var stderr bytes.Buffer
	cmd.Stdout = d.c.IO.Out
	cmd.Stderr = io.MultiWriter(d.c.IO.Err, &stderr)
	if d.c.IO.Interactive {
		cmd.Stdin = os.Stdin
	}
	err = cmd.Run()
	if err != nil {
		code := 1
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		}
		switch code {
		case 3:
			d.c.IO.WriteErrorV("<warning>Class "+className+" is not autoloadable, can not call "+ev.Name+" script</warning>", true, 16)
			return 0, nil
		case 4:
			d.c.IO.WriteErrorV("<warning>Method "+className+"::"+methodName+" is not callable, can not call "+ev.Name+" script</warning>", true, 16)
			return 0, nil
		case 5:
			return 1, nil
		}
		d.c.IO.WriteErrorV(fmt.Sprintf("<error>Script %s::%s handling the %s event terminated with an exception</error>", className, methodName, ev.Name), true, 16)
		// the callback's exception is rethrown like EventDispatcher does
		if raw, rerr := os.ReadFile(resultFile); rerr == nil {
			if v, derr := phpjson.Decode(raw); derr == nil {
				if m := phpjson.AsMap(v); m != nil {
					return 0, &CallbackError{Class: m.Str("class"), Msg: m.Str("message")}
				}
			}
		}
		return 0, &CallbackError{Msg: "Script " + className + "::" + methodName + " failed (exit code " + fmt.Sprint(code) + ")"}
	}
	return 0, nil
}

func nonNilMap(m *phpjson.Map) *phpjson.Map {
	if m == nil {
		return phpjson.NewMap()
	}
	return m
}

// writeScriptLoader writes a PSR-0/PSR-4/classmap loader like
// AutoloadGenerator::createLoader (no "files" autoloads).
func (d *EventDispatcher) writeScriptLoader(file string, devMode bool) error {
	im := NewInstallationManager(d.c)
	var packages []*pkg.Package
	for _, p := range d.c.LocalRepo.CanonicalPackages() {
		packages = append(packages, p)
	}
	data, err := autoload.LoaderData(d.c.Root(), packages, im.InstallPath, d.c.Cwd, d.c.VendorDir(), devMode)
	if err != nil {
		return err
	}
	return os.WriteFile(file, []byte(data), 0o600)
}

const scriptRunnerPHP = `<?php
// gocomposer script runner: executes a Composer script callback with stub
// implementations of the Composer API objects passed to scripts.
namespace Composer\IO {
    interface IOInterface {}
    class GocomposerIO implements IOInterface {
        private $v; private $d;
        public function __construct($v, $d) { $this->v = $v; $this->d = $d; }
        public function isInteractive() { return false; }
        public function isVerbose() { return $this->v >= 64; }
        public function isVeryVerbose() { return $this->v >= 128; }
        public function isDebug() { return $this->v >= 256; }
        public function isDecorated() { return $this->d; }
        private function fmt($m) { $m = is_array($m) ? implode(PHP_EOL, $m) : $m; return preg_replace('{</?(?:info|comment|error|question|warning|highlight|(?:fg|bg|options|href)=[^>]*)?>}', '', $m); }
        public function write($messages, $newline = true, $verbosity = 32) { if ($verbosity <= $this->v) { echo $this->fmt($messages), $newline ? PHP_EOL : ''; } }
        public function writeError($messages, $newline = true, $verbosity = 32) { if ($verbosity <= $this->v) { fwrite(STDERR, $this->fmt($messages) . ($newline ? PHP_EOL : '')); } }
        public function writeRaw($messages, $newline = true, $verbosity = 32) { $this->write($messages, $newline, $verbosity); }
        public function writeErrorRaw($messages, $newline = true, $verbosity = 32) { $this->writeError($messages, $newline, $verbosity); }
        public function overwrite($messages, $newline = true, $size = null, $verbosity = 32) { $this->write($messages, $newline, $verbosity); }
        public function overwriteError($messages, $newline = true, $size = null, $verbosity = 32) { $this->writeError($messages, $newline, $verbosity); }
        public function ask($question, $default = null) { return $default; }
        public function askConfirmation($question, $default = true) { return $default; }
        public function askAndValidate($question, $validator, $attempts = null, $default = null) { return $default; }
        public function askAndHideAnswer($question) { return null; }
        public function select($question, $choices, $default, $attempts = false, $errorMessage = 'Value "%s" is invalid', $multiselect = false) { return $default; }
        public function __call($name, $args) { throw new \RuntimeException('Composer\IO\IOInterface::' . $name . '() is not supported by gocomposer script stubs'); }
    }
}
namespace Composer {
    class Config {
        private $c;
        public function __construct(array $c) { $this->c = $c; }
        public function get($key, $flags = 0) { if (!array_key_exists($key, $this->c)) { throw new \RuntimeException('Composer\Config::get(' . $key . ') is not supported by gocomposer script stubs'); } return $this->c[$key]; }
        public function has($key) { return array_key_exists($key, $this->c); }
        public function all($flags = 0) { return array('config' => $this->c); }
        public function __call($name, $args) { throw new \RuntimeException('Composer\Config::' . $name . '() is not supported by gocomposer script stubs'); }
    }
    class RootPackageStub {
        private $p;
        public function __construct(array $p) { $this->p = $p; }
        public function getName() { return $this->p['name']; }
        public function getPrettyName() { return $this->p['pretty_name']; }
        public function getVersion() { return $this->p['version']; }
        public function getPrettyVersion() { return $this->p['pretty_version']; }
        public function getExtra() { return $this->p['extra']; }
        public function getType() { return $this->p['type']; }
        public function __call($name, $args) { throw new \RuntimeException('Composer\Package\RootPackage::' . $name . '() is not supported by gocomposer script stubs'); }
    }
    class Composer {
        private $config; private $package;
        public function __construct($config, $package) { $this->config = $config; $this->package = $package; }
        public function getConfig() { return $this->config; }
        public function getPackage() { return $this->package; }
        public function __call($name, $args) { throw new \RuntimeException('Composer\Composer::' . $name . '() is not supported by gocomposer script stubs'); }
    }
}
namespace Composer\EventDispatcher {
    class Event {
        protected $name; protected $args; protected $flags; private $stopped = false;
        public function __construct($name, array $args = array(), array $flags = array()) { $this->name = $name; $this->args = $args; $this->flags = $flags; }
        public function getName() { return $this->name; }
        public function getArguments() { return $this->args; }
        public function getFlags() { return $this->flags; }
        public function isPropagationStopped() { return $this->stopped; }
        public function stopPropagation() { $this->stopped = true; }
    }
}
namespace Composer\Script {
    class Event extends \Composer\EventDispatcher\Event {
        private $composer; private $io; private $devMode;
        public function __construct($name, $composer, $io, $devMode = false, array $args = array(), array $flags = array()) { parent::__construct($name, $args, $flags); $this->composer = $composer; $this->io = $io; $this->devMode = $devMode; }
        public function getComposer() { return $this->composer; }
        public function getIO() { return $this->io; }
        public function isDevMode() { return $this->devMode; }
        public function getOriginatingEvent() { return null; }
    }
}
namespace {
    $d = json_decode(getenv('GOCOMPOSER_SCRIPT_EVENT'), true);
    putenv('GOCOMPOSER_SCRIPT_EVENT');
    require $d['loader'];
    $io = new \Composer\IO\GocomposerIO($d['verbosity'], $d['decorated']);
    $composer = new \Composer\Composer(new \Composer\Config($d['config']), new \Composer\RootPackageStub($d['root']));
    $event = new \Composer\Script\Event($d['name'], $composer, $io, $d['dev'], $d['args']);
    $class = $d['class']; $method = $d['method'];
    if (!class_exists($class)) { exit(3); }
    if (!is_callable(array($class, $method))) { exit(4); }
    if ($d['verbose']) { fwrite(STDERR, '> ' . $d['name'] . ': ' . $class . '::' . $method . PHP_EOL); }
    elseif ($d['output']) { fwrite(STDERR, '> ' . $class . '::' . $method . PHP_EOL); }
    try {
        $r = $class::$method($event);
    } catch (\Throwable $e) {
        file_put_contents($d['result'], json_encode(array('class' => get_class($e), 'message' => $e->getMessage())));
        exit(6);
    }
    exit($r === false ? 5 : 0);
}
`
