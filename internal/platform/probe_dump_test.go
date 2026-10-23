package platform

import (
	"os"
	"testing"
)

// TestWriteProbeScript writes the probe script for recording a facts
// snapshot with an oracle PHP (GOCOMPOSER_WRITE_PROBE=<path>).
func TestWriteProbeScript(t *testing.T) {
	p := os.Getenv("GOCOMPOSER_WRITE_PROBE")
	if p == "" {
		t.Skip("GOCOMPOSER_WRITE_PROBE not set")
	}
	if err := os.WriteFile(p, []byte("<?php\n"+ProbeScript()), 0o644); err != nil {
		t.Fatal(err)
	}
}
