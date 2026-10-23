//go:build linux

// Command ptyrun runs a command on a fresh pseudo terminal inside the oracle
// container and feeds it a scripted input sequence. It is used by the compare
// harness (tools/compare/tty.go) for both official Composer and gocomposer so
// that both see the same terminal and the same keystrokes.
//
// Usage: ptyrun <steps.json> <raw-transcript> -- <command> [args...]
//
// Each step waits for its prompt text (ANSI/OSC sequences and CR stripped) to
// appear after the previous match, for the output to settle, and optionally
// for a terminal mode (echo off / non-canonical), then writes its bytes to the
// master side. The raw byte transcript of the terminal is written to
// <raw-transcript>; diagnostics go to stderr; the exit code is the command's.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type step struct {
	Expect string `json:"expect,omitempty"`
	Send   string `json:"send,omitempty"`
	EOF    bool   `json:"eof,omitempty"`
	// Mode waits for a terminal state before sending: "noecho" (stty -echo,
	// hidden answers) or "raw" (stty -icanon, autocompleted answers).
	Mode string `json:"mode,omitempty"`
}

const (
	expectTimeout = 90 * time.Second
	settle        = 150 * time.Millisecond
	finishTimeout = 60 * time.Second
)

var reCtl = regexp.MustCompile("\x1b\\[[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]|\x1b\\][^\x07\x1b]*(?:\x07|\x1b\\\\)|\x1b[^\\[\\]]|\r")

type term struct {
	mu   sync.Mutex
	cond *sync.Cond
	raw  []byte
	last time.Time
	done bool
}

func ioctl(fd uintptr, req uintptr, arg unsafe.Pointer) error {
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(arg)); e != 0 {
		return e
	}
	return nil
}

func main() {
	if len(os.Args) < 5 || os.Args[3] != "--" {
		fmt.Fprintln(os.Stderr, "usage: ptyrun <steps.json> <raw-transcript> -- <command> [args...]")
		os.Exit(125)
	}
	var steps []step
	b, err := os.ReadFile(os.Args[1])
	if err != nil || json.Unmarshal(b, &steps) != nil {
		fmt.Fprintln(os.Stderr, "ptyrun: cannot read steps", err)
		os.Exit(125)
	}
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun:", err)
		os.Exit(125)
	}
	var unlock int32
	var ptn uint32
	if err := ioctl(master.Fd(), syscall.TIOCSPTLCK, unsafe.Pointer(&unlock)); err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun: unlockpt:", err)
		os.Exit(125)
	}
	if err := ioctl(master.Fd(), syscall.TIOCGPTN, unsafe.Pointer(&ptn)); err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun: ptsname:", err)
		os.Exit(125)
	}
	slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", ptn), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun:", err)
		os.Exit(125)
	}
	ws := struct{ row, col, x, y uint16 }{24, 80, 0, 0}
	_ = ioctl(slave.Fd(), syscall.TIOCSWINSZ, unsafe.Pointer(&ws))

	cmd := exec.Command(os.Args[4], os.Args[5:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = slave, slave, slave
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 0}
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun:", err)
		os.Exit(125)
	}
	t := &term{last: time.Now()}
	t.cond = sync.NewCond(&t.mu)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := master.Read(buf)
			t.mu.Lock()
			if n > 0 {
				t.raw = append(t.raw, buf[:n]...)
				t.last = time.Now()
			}
			if err != nil {
				t.done = true
			}
			t.cond.Broadcast()
			t.mu.Unlock()
			if err != nil {
				return
			}
		}
	}()
	// wake waiters periodically for timeouts and the settle period
	go func() {
		for {
			time.Sleep(20 * time.Millisecond)
			t.mu.Lock()
			t.cond.Broadcast()
			d := t.done
			t.mu.Unlock()
			if d {
				return
			}
		}
	}()
	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()

	failed := ""
	pos := 0
	for i, s := range steps {
		deadline := time.Now().Add(expectTimeout)
		t.mu.Lock()
		for {
			if t.done {
				failed = fmt.Sprintf("ptyrun: the command ended before input step %d (expect %q) was sent", i+1, s.Expect)
				break
			}
			if time.Now().After(deadline) {
				failed = fmt.Sprintf("ptyrun: timeout waiting for %q (step %d)", s.Expect, i+1)
				break
			}
			stripped := reCtl.ReplaceAllString(string(t.raw), "")
			if s.Expect != "" {
				if idx := strings.Index(stripped[min(pos, len(stripped)):], s.Expect); idx >= 0 {
					pos += idx + len(s.Expect)
					s.Expect = ""
				}
			}
			if s.Expect == "" && time.Since(t.last) >= settle && modeReady(master, s.Mode) {
				break
			}
			t.cond.Wait()
		}
		t.mu.Unlock()
		if failed != "" {
			break
		}
		data := s.Send
		if s.EOF {
			data += "\x04"
		}
		if _, err := master.Write([]byte(data)); err != nil {
			failed = "ptyrun: write: " + err.Error()
			break
		}
	}
	if failed != "" {
		fmt.Fprintln(os.Stderr, failed)
		killGroup(cmd)
	}
	// a command still waiting for input after the last step is stopped
	if failed == "" {
		select {
		case err = <-waitErr:
		case <-time.After(finishTimeout):
			fmt.Fprintln(os.Stderr, "ptyrun: the command did not finish after the last input step")
			killGroup(cmd)
			err = <-waitErr
		}
	} else {
		err = <-waitErr
	}
	// children left behind (e.g. php under sh -c) would keep the terminal open
	killGroup(cmd)
	slave.Close()
	// drain what is left until the terminal reports EIO
	deadline := time.Now().Add(5 * time.Second)
	t.mu.Lock()
	for !t.done && time.Now().Before(deadline) {
		t.cond.Wait()
	}
	raw := append([]byte(nil), t.raw...)
	t.mu.Unlock()
	if werr := os.WriteFile(os.Args[2], raw, 0o644); werr != nil {
		fmt.Fprintln(os.Stderr, "ptyrun:", werr)
	}
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
		if code < 0 {
			code = 128 + int(ee.Sys().(syscall.WaitStatus).Signal())
		}
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "ptyrun:", err)
		code = 125
	}
	os.Exit(code)
}

// killGroup kills the session started for the command (Setsid makes its pid
// the process group id), so grandchildren cannot outlive ptyrun.
func killGroup(cmd *exec.Cmd) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_ = cmd.Process.Kill()
}

// modeReady reports whether the terminal is in the requested mode.
func modeReady(master *os.File, mode string) bool {
	if mode == "" {
		return true
	}
	var tio syscall.Termios
	if ioctl(master.Fd(), syscall.TCGETS, unsafe.Pointer(&tio)) != nil {
		return true
	}
	switch mode {
	case "noecho":
		return tio.Lflag&syscall.ECHO == 0
	case "raw":
		return tio.Lflag&syscall.ICANON == 0 && tio.Lflag&syscall.ECHO == 0
	case "line":
		return tio.Lflag&syscall.ICANON != 0 && tio.Lflag&syscall.ECHO != 0
	}
	return true
}
