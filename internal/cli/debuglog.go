package cli

// Port of the -vvv "Running ..." line of Composer\Console\Application::doRun
// and of the deprecation notice handling of Composer\Util\ErrorHandler
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/platform"
)

// registerErrorHandler is ErrorHandler::register($io) for the deprecation
// notices gocomposer reproduces (config.bin-compat "symlink"): the first is
// shown, then a single "More deprecation notices were hidden" line unless
// verbose. Composer prints the PHP file and line of trigger_error() and,
// when verbose, a PHP stack trace; gocomposer prints "in gocomposer" and no
// trace.
func registerErrorHandler(io *cio.IO) {
	shown := 0
	config.Deprecation = func(msg string) {
		if shown > 0 && !io.IsVerbose() {
			if shown == 1 {
				io.WriteError("<warning>More deprecation notices were hidden, run again with `-v` to show them.</warning>")
				shown = 2
			}
			return
		}
		shown = 1
		io.WriteError("<warning>Deprecation Notice: " + msg + " in gocomposer</warning>")
	}
}

// debugVerbosity is the output verbosity Symfony's configureIO sets from
// the raw command line before Application::doRun.
func debugVerbosity(argv []string) int {
	if argvHas(argv, "-q", "--quiet") {
		return cio.Quiet
	}
	if v := verbosityFrom(argv); v > 0 {
		return v
	}
	switch os.Getenv("SHELL_VERBOSITY") {
	case "-1":
		return cio.Quiet
	case "1":
		return cio.Verbose
	case "2":
		return cio.VeryVerbose
	case "3":
		return cio.Debug
	}
	return cio.Normal
}

// debugRunningLine writes Application::doRun's
// "Running <version> (<date>) with PHP <version> on <os>" debug line for
// commands that are not proxy commands (global, outdated).
func (a *App) debugRunningLine(io *cio.IO, argv []string, cmdName string) {
	if debugVerbosity(argv) < cio.Debug {
		return
	}
	if c, err := a.findCommand(cmdName); err == nil && (c.Name == "global" || c.Name == "outdated") {
		return
	}
	saved := io.Verbosity
	io.Verbosity = cio.Debug
	defer func() { io.Verbosity = saved }()
	io.WriteErrorV("Running "+platform.ComposerVersion+" ("+platform.ComposerReleaseDate+") with "+runtimeDescription()+" on "+osDescription(), true, cio.Debug)
}

// runtimeDescription is 'PHP '.PHP_VERSION of the PHP binary gocomposer
// uses for platform packages and scripts.
func runtimeDescription() string {
	if php, err := platform.FindPHP(); err == nil {
		if out, err := exec.Command(php, "-r", "echo PHP_VERSION;").Output(); err == nil && len(out) > 0 {
			return "PHP " + strings.TrimSpace(string(out))
		}
	}
	return "gocomposer " + Version + " (" + runtime.Version() + ", no PHP found)"
}

// osDescription is php_uname('s').' / '.php_uname('r').
func osDescription() string {
	s, r := unameSR()
	return s + " / " + r
}

func unameSR() (string, string) {
	sys := ""
	rel := ""
	if out, err := exec.Command("uname", "-s").Output(); err == nil {
		sys = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		rel = strings.TrimSpace(string(out))
	}
	if sys == "" {
		sys = strings.ToUpper(runtime.GOOS[:1]) + runtime.GOOS[1:]
	}
	return sys, rel
}
