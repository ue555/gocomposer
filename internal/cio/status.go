package cio

// Program status reporting of Composer\IO\ConsoleIO with
// Composer\Console\TerminalStatusPolicy, OscHelper, ProgramStatus,
// Osc7501ProgramStatusEncoder and Osc9ProgressEncoder
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano): OSC sequences written to a terminal
// error output telling the terminal whether Composer works, waits for an
// answer or finished.

import (
	"encoding/base64"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ProgramStatus is Composer\Console\ProgramStatus.
type ProgramStatus struct {
	State      string // idle, working, done, blocked, error, clear
	Message    string
	hasMessage bool
	Kind       string // osc7501 "kind" protocol option
	Progress   *int
}

// StatusFromExitCode is ProgramStatus::fromExitCode.
func StatusFromExitCode(code int) *ProgramStatus {
	switch {
	case code == 129 || code == 130 || code == 143:
		return &ProgramStatus{State: "idle"}
	case code == 0:
		return &ProgramStatus{State: "done"}
	}
	return &ProgramStatus{State: "error"}
}

// StatusError is ProgramStatus::error().
func StatusError() *ProgramStatus { return &ProgramStatus{State: "error"} }

// StatusIdle is ProgramStatus::idle().
func StatusIdle() *ProgramStatus { return &ProgramStatus{State: "idle"} }

var reStatusSemver = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
var reStatusMux = regexp.MustCompile(`^(?:screen|tmux)(?:[-.]|$)`)

// InitProgramStatus selects the encoders like TerminalStatusPolicy::fromEnvironment.
func (o *IO) InitProgramStatus() {
	env := func(k string) string { return os.Getenv(k) }
	setting := strings.TrimSpace(os.Getenv("COMPOSER_TERMINAL_STATUS"))
	_, set := os.LookupEnv("COMPOSER_TERMINAL_STATUS")
	if !set {
		setting = "auto"
	}
	var protocols []string
	switch setting {
	case "", "auto":
		protocols = []string{"7501"}
		if supportsProgress(env) {
			protocols = []string{"7501", "9;4"}
		}
	case "off":
	default:
		seen := map[string]bool{}
		for _, p := range strings.Split(setting, ",") {
			p = strings.TrimSpace(p)
			if (p == "7501" || p == "9;4") && !seen[p] {
				seen[p] = true
				protocols = append(protocols, p)
			}
		}
	}
	o.statusEncoders = nil
	has := func(p string) bool {
		for _, x := range protocols {
			if x == p {
				return true
			}
		}
		return false
	}
	if has("9;4") {
		o.statusEncoders = append(o.statusEncoders, "9;4")
	}
	if has("7501") {
		o.statusEncoders = append(o.statusEncoders, "7501")
	}
}

func supportsProgress(env func(string) string) bool {
	for _, k := range []string{"TMUX", "STY", "SSH_TTY", "SSH_CONNECTION"} {
		if env(k) != "" {
			return false
		}
	}
	if reStatusMux.MatchString(env("TERM")) {
		return false
	}
	if env("WT_SESSION") != "" || env("ConEmuANSI") == "ON" {
		return true
	}
	min := map[string]string{"ghostty": "1.2.0", "iTerm.app": "3.6.6"}
	m, ok := min[env("TERM_PROGRAM")]
	v := env("TERM_PROGRAM_VERSION")
	return ok && reStatusSemver.MatchString(v) && versionGE(v, m)
}

func versionGE(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		x, _ := strconv.Atoi(pa[i])
		y, _ := strconv.Atoi(pb[i])
		if x != y {
			return x > y
		}
	}
	return true
}

// supportsOutput is OscHelper::supportsOutput for the error output.
func (o *IO) supportsOutput() bool {
	if !o.Decorated || o.Verbosity <= Quiet {
		return false
	}
	f, ok := o.Err.(*os.File)
	return !ok || IsTTY(f)
}

// IsProgramStatusEnabled is ConsoleIO::isProgramStatusEnabled.
func (o *IO) IsProgramStatusEnabled() bool {
	return len(o.statusEncoders) > 0 && o.supportsOutput()
}

// StartProgramStatus is ConsoleIO::startProgramStatus.
func (o *IO) StartProgramStatus() {
	o.programStatus = &ProgramStatus{State: "working"}
	for _, enc := range o.statusEncoders {
		if enc == "9;4" {
			continue // Osc9ProgressEncoder::encodeStart
		}
		o.writeOsc(enc, encodeStatus(enc, o.programStatus))
	}
}

// WriteProgramStatus is ConsoleIO::writeProgramStatus.
func (o *IO) WriteProgramStatus(st *ProgramStatus) {
	o.programStatus = st
	for _, enc := range o.statusEncoders {
		o.writeOsc(enc, encodeStatus(enc, st))
	}
}

// WriteProgress is ConsoleIO::writeProgress (nil clears the progress).
func (o *IO) WriteProgress(progress *int) {
	st := o.programStatus
	if st == nil {
		st = &ProgramStatus{State: "working"}
	}
	c := *st
	c.Progress = nil
	if progress != nil && *progress >= 0 && *progress <= 100 && (c.State == "working" || c.State == "blocked") {
		p := *progress
		c.Progress = &p
	}
	o.WriteProgramStatus(&c)
}

func (o *IO) writeOsc(command, payload string) {
	if !o.supportsOutput() {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	_, _ = o.Err.Write([]byte("\x1b]" + command + ";" + payload + "\x1b\\"))
}

func encodeStatus(enc string, st *ProgramStatus) string {
	if enc == "9;4" {
		switch st.State {
		case "working":
			if st.Progress == nil {
				return "3"
			}
			return "1;" + strconv.Itoa(*st.Progress)
		case "blocked":
			if st.Progress == nil {
				return "4"
			}
			return "4;" + strconv.Itoa(*st.Progress)
		case "error":
			return "2"
		}
		return "0"
	}
	payload := "state=" + st.State + ":app=composer"
	if st.Progress != nil {
		payload += ":progress=" + strconv.Itoa(*st.Progress)
	}
	if st.State == "blocked" && (st.Kind == KindPermission || st.Kind == KindQuestion || st.Kind == KindAuth) {
		payload += ":kind=" + st.Kind
	}
	if st.hasMessage {
		payload += ":msg=" + base64.StdEncoding.EncodeToString([]byte(statusMessage(st.Message, 2048)))
	}
	return payload
}

var (
	reStatusCtl     = regexp.MustCompile("[\x00-\x1f\x7f]|\u0080|[\u0081-\u009f]")
	reStatusTrailMB = regexp.MustCompile(`(?:[\xc2-\xdf]|[\xe0-\xef][\x80-\xbf]{0,1}|[\xf0-\xf4][\x80-\xbf]{0,2})$`)
)

// statusMessage is AbstractOscProgramStatusEncoder::sanitizeMessage.
func statusMessage(msg string, maxBytes int) string {
	msg = sanitize(msg, false)
	msg = reStatusCtl.ReplaceAllString(msg, "")
	msg = (&IO{}).Format(msg)
	if len(msg) > maxBytes {
		msg = msg[:maxBytes]
		msg = reStatusTrailMB.ReplaceAllString(msg, "")
	}
	return phpTrim(msg)
}
