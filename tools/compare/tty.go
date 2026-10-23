package main

// Pseudo-terminal runs for the interactive scenarios (tty-*).
//
// Composer only asks questions when STDIN is a terminal and neither -n nor
// COMPOSER_NO_INTERACTION is given (Composer\Console\Application::doRunComposer),
// so these scenarios run the compared command through tools/compare/ptyrun,
// a small Go helper cross-compiled into /gobin that allocates a Linux pty
// inside the oracle container, makes it the controlling terminal of the
// command and feeds the scenario's input steps (each one waits for its prompt
// text to appear). Both tools get the same helper, terminal (80x24, TERM=xterm)
// and input sequence. stdout and stderr share the terminal, so the compared
// "stdout" is the whole terminal transcript normalized by normalizeTTY (CRLF
// to LF, remaining control bytes made visible); the raw bytes are kept in the
// result JSON (tty_raw, base64) and the helper's diagnostics are "stderr".

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// TTYStep is one scripted input of a terminal scenario.
type TTYStep struct {
	// Expect is text (ANSI/OSC sequences and CR removed) that must appear in
	// the terminal output after the previous step's match before sending.
	Expect string `json:"expect,omitempty"`
	// Send is written to the terminal as typed (include "\n" for Enter).
	Send string `json:"send,omitempty"`
	// EOF sends Ctrl-D after Send.
	EOF bool `json:"eof,omitempty"`
	// Mode waits for a terminal state: "noecho" (hidden answer, stty -echo),
	// "raw" (autocompletion, stty -icanon -echo) or "line" (canonical, echo).
	Mode string `json:"mode,omitempty"`
}

// Answer types the line after the prompt, then Enter.
func Answer(expect, text string) TTYStep { return TTYStep{Expect: expect, Send: text + "\n"} }

// Enter answers the prompt with an empty line.
func Enter(expect string) TTYStep { return TTYStep{Expect: expect, Send: "\n"} }

// CtrlD answers the prompt with end of input (Ctrl-D at the line start).
func CtrlD(expect string) TTYStep { return TTYStep{Expect: expect, EOF: true} }

// Hidden types an answer once echo is turned off (askAndHideAnswer).
func Hidden(expect, text string) TTYStep {
	return TTYStep{Expect: expect, Send: text + "\n", Mode: "noecho"}
}

// Keys types raw keystrokes once the terminal is non-canonical (the
// QuestionHelper autocompletion used by ChoiceQuestion/select).
func Keys(expect, keys string) TTYStep { return TTYStep{Expect: expect, Send: keys, Mode: "raw"} }

var (
	ptyOnce sync.Once
	ptyErr  error
)

// ensurePtyHelper cross-compiles tools/compare/ptyrun into gobin.
func ensurePtyHelper(root, gobin string) error {
	ptyOnce.Do(func() {
		arch, err := podmanArch()
		if err != nil {
			ptyErr = err
			return
		}
		build := exec.Command("go", "build", "-trimpath", "-o", filepath.Join(gobin, "ptyrun"), "./tools/compare/ptyrun")
		build.Dir = root
		build.Env = append(os.Environ(), "GOOS=linux", "GOARCH="+arch, "CGO_ENABLED=0", "GOCACHE="+filepath.Join(root, ".cache", "go-build"))
		if out, err := build.CombinedOutput(); err != nil {
			ptyErr = fmt.Errorf("building ptyrun: %v\n%s", err, out)
		}
	})
	return ptyErr
}

// isTTYRun reports whether this run of the scenario is the compared one on a
// terminal (Setup commands run in the template directory without a terminal).
func isTTYRun(sc Scenario, dir string) bool {
	return sc.TTY && filepath.Base(dir) != "template"
}

// ttyCommandLine writes the input steps and returns the shell line running
// the tool on a pseudo terminal.
func ttyCommandLine(root, gobin, dir string, t tool, command string, sc Scenario) (string, error) {
	if err := ensurePtyHelper(root, gobin); err != nil {
		return "", err
	}
	steps := sc.Input
	if steps == nil {
		steps = []TTYStep{}
	}
	b, _ := json.Marshal(steps)
	if err := os.WriteFile(filepath.Join(dir, "out", "ttyinput.json"), b, 0o644); err != nil {
		return "", err
	}
	inner := t.cmd + " " + command
	return "TERM=xterm /gobin/ptyrun /work/out/ttyinput.json /work/out/ttyraw -- sh -c " + shQuote(inner) +
		" >/work/out/stdout 2>/work/out/stderr; echo $? >/work/out/exit", nil
}

func shQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// readTTYResult fills the result from the terminal transcript.
func readTTYResult(dir string, r *result) {
	raw, _ := os.ReadFile(filepath.Join(dir, "out", "ttyraw"))
	r.TTYRaw = raw
	r.Stdout = normalizeTTY(string(raw))
}

// normalizeTTY turns a terminal transcript into comparable text: the pty's
// CRLF output translation is undone and every other control byte (escape
// sequences, lone CR, backspace, bell) is kept but written visibly, so
// colors, cursor movements and line rewrites still take part in the diff.
func normalizeTTY(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var sb strings.Builder
	for _, r := range s {
		switch {
		case r == '\n' || r == '\t':
			sb.WriteRune(r)
		case r == 0x1b:
			sb.WriteString(`\e`)
		case r == '\r':
			sb.WriteString(`\r`)
		case r < 0x20 || r == 0x7f:
			fmt.Fprintf(&sb, `\x%02x`, r)
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
