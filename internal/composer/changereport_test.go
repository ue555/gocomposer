package composer

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Comparer::getChangedAsString(true): sections in first-seen order, items
// suffixed with "\r\n" and joined with "\r\n"; empty files are ignored.
func TestCompareTrees(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeTestFile(t, filepath.Join(src, "src/A.php"), "a")
	writeTestFile(t, filepath.Join(dst, "src/A.php"), "changed")
	writeTestFile(t, filepath.Join(dst, "src/New.php"), "new")
	writeTestFile(t, filepath.Join(src, "empty.txt"), "")
	if got, want := compareTrees(src, dst), "./src/A.php\r\n\r\n./src/New.php"; got != want {
		t.Fatalf("compareTrees = %q, want %q", got, want)
	}
	if got := compareTrees(src, src); got != "" {
		t.Fatalf("identical trees reported %q", got)
	}
}

func TestCompareTreesRemoved(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeTestFile(t, filepath.Join(src, "composer.json"), "{}")
	writeTestFile(t, filepath.Join(dst, "composer.json"), "")
	if got, want := compareTrees(src, dst), "./composer.json"; got != want {
		t.Fatalf("compareTrees = %q, want %q", got, want)
	}
}
