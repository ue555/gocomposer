package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGcIsNecessaryOncePerProcess(t *testing.T) {
	t.Setenv("COMPOSER_TEST_SUITE", "1")
	collected = false
	if GcIsNecessary() {
		t.Fatal("COMPOSER_TEST_SUITE must disable gc")
	}
	t.Setenv("COMPOSER_TEST_SUITE", "")
	if GcIsNecessary() {
		t.Fatal("gc must be considered at most once per process")
	}
}

func TestCacheGCRespectsReadOnly(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a", "old.zip")
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(f, old, old)
	if New(dir, "a-z0-9_./", true).GC(3600, 1<<30) {
		t.Fatal("read-only cache must not be collected")
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatal("file removed from read-only cache")
	}
	if !New(dir, "a-z0-9_./", false).GC(3600, 1<<30) {
		t.Fatal("writable cache must be collected")
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatal("expired file not removed")
	}
}
