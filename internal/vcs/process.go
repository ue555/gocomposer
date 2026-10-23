// Package vcs ports Composer's version control helpers: the subset of
// Util\ProcessExecutor used to run VCS commands, Util\Git, Util\Url and the
// Repository\Vcs drivers.
//
// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/ProcessExecutor.php, src/Composer/Util/Git.php,
// src/Composer/Util/Url.php, src/Composer/Repository/Vcs/*.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
package vcs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/httpdl"
)

// Process is the ProcessExecutor subset used by the VCS code: commands are
// argument lists run without a shell, their output is captured.
type Process struct {
	IO *cio.IO
	// Timeout is ProcessExecutor::$timeout (config process-timeout); 0
	// disables it.
	Timeout time.Duration
	mu      sync.Mutex
	errOut  string
	// mock replaces command execution in tests (ProcessExecutorMock).
	mock func(cmd []string, cwd string) (int, string, string, error)
}

// NewProcess creates an executor writing debug output to io.
func NewProcess(io *cio.IO, timeout time.Duration) *Process {
	return &Process{IO: io, Timeout: timeout}
}

// TimeoutError is Symfony's ProcessTimedOutException.
type TimeoutError struct{ Msg string }

func (e *TimeoutError) Error() string { return e.Msg }

var gitCmdsNeedGitDir = [][]string{{"show"}, {"log"}, {"branch"}, {"remote", "set-url"}}

func requiresGitDirEnv(cmd []string) bool {
	if len(cmd) == 0 || cmd[0] != "git" {
		return false
	}
	for _, need := range gitCmdsNeedGitDir {
		all := true
		for _, n := range need {
			found := false
			for _, c := range cmd {
				if c == n {
					found = true
					break
				}
			}
			if !found {
				all = false
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}

// Escape is ProcessExecutor::escape on non-Windows platforms.
func Escape(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", `'\''`) + "'"
}

var rePasswordArg = regexp.MustCompile(`--password (.*[^\\]') `)

func (p *Process) outputCommandRun(cmd []string, cwd string) {
	if p.IO == nil || !p.IO.IsDebug() {
		return
	}
	parts := make([]string, len(cmd))
	for i, a := range cmd {
		parts[i] = Escape(a)
	}
	safe := httpdl.Sanitize(strings.Join(parts, " "))
	safe = rePasswordArg.ReplaceAllString(safe, "--password '***' ")
	if cwd == "" {
		cwd = "CWD"
	}
	p.IO.WriteErrorV("Executing command ("+cwd+"): "+safe, true, cio.Debug)
}

// Execute is ProcessExecutor::execute with captured output. It returns the
// exit code and stdout; stderr is available from ErrorOutput. A timeout
// returns a *TimeoutError.
func (p *Process) Execute(cmd []string, cwd string) (int, string, error) {
	if p.mock != nil {
		code, out, errOut, err := p.mock(cmd, cwd)
		p.mu.Lock()
		p.errOut = errOut
		p.mu.Unlock()
		return code, out, err
	}
	p.outputCommandRun(cmd, cwd)
	p.mu.Lock()
	p.errOut = ""
	p.mu.Unlock()
	var env []string
	if cwd != "" && requiresGitDirEnv(cmd) {
		if st, err := os.Stat(strings.TrimRight(cwd, "/") + "/.git"); err != nil || !st.IsDir() {
			_, out, errOut, _ := p.run([]string{"git", "config", "safe.bareRepository"}, cwd, []string{"GIT_DIR=" + cwd})
			_ = errOut
			if strings.TrimSpace(out) == "explicit" {
				env = []string{"GIT_DIR=" + cwd}
			}
		}
	}
	code, out, errOut, err := p.run(cmd, cwd, env)
	p.mu.Lock()
	p.errOut = errOut
	p.mu.Unlock()
	return code, out, err
}

// ErrorOutput is ProcessExecutor::getErrorOutput (of the last Execute).
func (p *Process) ErrorOutput() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.errOut
}

func (p *Process) run(cmd []string, cwd string, extraEnv []string) (int, string, string, error) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if p.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, p.Timeout)
		defer cancel()
	}
	c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	if cwd != "" {
		c.Dir = cwd
	}
	if len(extraEnv) > 0 {
		c.Env = append(os.Environ(), extraEnv...)
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if ctx.Err() == context.DeadlineExceeded {
		parts := make([]string, len(cmd))
		for i, a := range cmd {
			parts[i] = Escape(a)
		}
		return -1, stdout.String(), stderr.String(), &TimeoutError{Msg: fmt.Sprintf(`The process "%s" exceeded the timeout of %d seconds.`, strings.Join(parts, " "), int(p.Timeout/time.Second))}
	}
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return ee.ExitCode(), stdout.String(), stderr.String(), nil
		}
		// the command could not be started: Symfony Process runs the
		// command line through /bin/sh, which reports a missing binary
		// with exit code 127
		if cwd != "" {
			if st, serr := os.Stat(cwd); serr != nil || !st.IsDir() {
				return 1, "", "", &RuntimeError{Msg: `The provided cwd "` + cwd + `" does not exist.`}
			}
		}
		return 127, stdout.String(), "sh: 1: exec: " + cmd[0] + ": not found\n", nil
	}
	return 0, stdout.String(), stderr.String(), nil
}

// RuntimeError is a PHP \RuntimeException raised by the VCS layer.
type RuntimeError struct {
	Msg  string
	Prev error
}

func (e *RuntimeError) Error() string   { return e.Msg }
func (e *RuntimeError) Previous() error { return e.Prev }

// InvalidArgumentError is a PHP \InvalidArgumentException.
type InvalidArgumentError struct{ Msg string }

func (e *InvalidArgumentError) Error() string { return e.Msg }

var reLineSplit = regexp.MustCompile(`\r?\n`)

// SplitLines is ProcessExecutor::splitLines.
func SplitLines(output string) []string {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	return reLineSplit.Split(output, -1)
}
