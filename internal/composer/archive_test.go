package composer

import (
	"archive/tar"
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type entry struct {
	name, body string
	link       bool
}

func writeTestZip(t *testing.T, entries []entry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "a.zip")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e.name, Method: zip.Store}
		if e.link {
			h.SetMode(os.ModeSymlink | 0o777)
		} else {
			h.SetMode(0o644)
		}
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(e.body))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return p
}

// Composer extracts zips with unzip, which strips "../" components: the
// entry must land inside dest and nothing may be written outside of it.
func TestExtractZipStripsParentComponents(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "a", "b")
	os.MkdirAll(dest, 0o755)
	z := writeTestZip(t, []entry{{name: "pkg/ok.txt", body: "ok"}, {name: "../../evil.txt", body: "x"}, {name: "/abs.txt", body: "y"}})
	if err := ExtractZip(z, dest); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"pkg/ok.txt", "evil.txt", "abs.txt"} {
		if _, err := os.Stat(filepath.Join(dest, p)); err != nil {
			t.Errorf("%s not extracted inside dest: %v", p, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "evil.txt")); err == nil {
		t.Fatal("entry escaped the extraction directory")
	}
}

func TestExtractZipRejectsEscapingSymlink(t *testing.T) {
	dest := t.TempDir()
	z := writeTestZip(t, []entry{{name: "pkg/link", body: "../../../etc", link: true}})
	err := ExtractZip(z, dest)
	if err == nil || !strings.Contains(err.Error(), "points outside") {
		t.Fatalf("expected symlink escape error, got %v", err)
	}
}

func TestExtractZipRejectsWriteThroughSymlink(t *testing.T) {
	dest := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(dest, "dir")); err != nil {
		t.Fatal(err)
	}
	z := writeTestZip(t, []entry{{name: "dir/file.txt", body: "x"}})
	if err := ExtractZip(z, dest); err == nil {
		t.Fatal("expected error writing below a symlink")
	}
	if _, err := os.Stat(filepath.Join(outside, "file.txt")); err == nil {
		t.Fatal("file written through symlink")
	}
}

func TestExtractTarRejectsTraversal(t *testing.T) {
	p := filepath.Join(t.TempDir(), "a.tar")
	f, _ := os.Create(p)
	tw := tar.NewWriter(f)
	tw.WriteHeader(&tar.Header{Name: "../evil.txt", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	tw.Write([]byte("x"))
	tw.Close()
	f.Close()
	dest := filepath.Join(t.TempDir(), "d")
	os.MkdirAll(dest, 0o755)
	if err := ExtractTar(p, dest); err == nil {
		t.Fatal("expected traversal error")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(dest), "evil.txt")); err == nil {
		t.Fatal("tar entry escaped")
	}
}
