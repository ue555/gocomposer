package composer

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gocomposer/internal/cio"
)

func debugIO() (*cio.IO, *bytes.Buffer) {
	var buf bytes.Buffer
	return &cio.IO{Out: &bytes.Buffer{}, Err: &buf, Verbosity: cio.Debug}, &buf
}

func TestReadingLog(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "composer.json")
	if err := os.WriteFile(f, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	io, buf := debugIO()
	ReadingLog(io, f)
	if got := buf.String(); got != "Reading "+f+"\n" && !strings.HasPrefix(got, "Reading "+f+" (") {
		t.Errorf("absolute path: %q", got)
	}
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	ReadingLog(io, "./composer.json")
	real, _ := filepath.EvalSymlinks(f)
	if got, want := buf.String(), "Reading ./composer.json ("+real+")\n"; got != want {
		t.Errorf("relative path: %q, want %q", got, want)
	}
	io.Verbosity = cio.VeryVerbose
	buf.Reset()
	ReadingLog(io, "./composer.json")
	if buf.Len() != 0 {
		t.Errorf("non-debug output: %q", buf.String())
	}
}

func TestValidateCaFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.pem")
	if err := os.WriteFile(bad, []byte("not a certificate"), 0o644); err != nil {
		t.Fatal(err)
	}
	var logs []string
	if validateCaFile(bad, func(m string) { logs = append(logs, m) }) {
		t.Error("garbage accepted as CA file")
	}
	real, _ := filepath.EvalSymlinks(bad)
	if len(logs) != 1 || logs[0] != "Checked CA file "+real+": invalid" {
		t.Errorf("logs = %q", logs)
	}
	// cached: no second log line
	validateCaFile(bad, func(m string) { logs = append(logs, m) })
	if len(logs) != 1 {
		t.Errorf("validity not cached: %q", logs)
	}
}

func TestIsScriptEventsConstant(t *testing.T) {
	for name, want := range map[string]bool{"post-autoload-dump": true, "POST_INSTALL_CMD": true, "hello": false, "command": false} {
		if got := isScriptEventsConstant(name); got != want {
			t.Errorf("%s: %v", name, got)
		}
	}
}
