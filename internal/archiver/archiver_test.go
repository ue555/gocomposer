package archiver

// Ported from tests/Composer/Test/Package/Archiver/ArchivableFilesFinderTest.php
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"archive/zip"
	"bytes"
	"hash/crc32"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var fileTree = []string{
	".foo", "A/prefixA.foo", "A/prefixB.foo", "A/prefixC.foo", "A/prefixD.foo", "A/prefixE.foo", "A/prefixF.foo",
	"B/sub/prefixA.foo", "B/sub/prefixB.foo", "B/sub/prefixC.foo", "B/sub/prefixD.foo", "B/sub/prefixE.foo", "B/sub/prefixF.foo",
	"C/prefixA.foo", "C/prefixB.foo", "C/prefixC.foo", "C/prefixD.foo", "C/prefixE.foo", "C/prefixF.foo",
	"D/prefixA", "D/prefixB", "D/prefixC", "D/prefixD", "D/prefixE", "D/prefixF",
	"E/subtestA.foo", "F/subtestA.foo", "G/subtestA.foo", "H/subtestA.foo", "I/J/subtestA.foo", "K/dirJ/subtestA.foo",
	"toplevelA.foo", "toplevelB.foo", "prefixA.foo", "prefixB.foo", "prefixC.foo", "prefixD.foo", "prefixE.foo", "prefixF.foo",
	"parameters.yml", "parameters.yml.dist", "!important!.txt", "!important_too!.txt", "#weirdfile",
}

func setupTree(t *testing.T) string {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range fileTree {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o666); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func archivable(t *testing.T, sources string, excludes []string, ignore bool) []string {
	items, base, err := FindFiles(sources, excludes, ignore)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, it := range items {
		if !it.IsDir {
			out = append(out, strings.TrimPrefix(it.RealPath, base))
		}
	}
	sort.Strings(out)
	return out
}

func TestManualExcludes(t *testing.T) {
	src := setupTree(t)
	got := archivable(t, src, []string{"prefixB.foo", "!/prefixB.foo", "/prefixA.foo", "prefixC.*", "!*/*/*/prefixC.foo", ".*"}, false)
	want := []string{"/!important!.txt", "/!important_too!.txt", "/#weirdfile", "/A/prefixA.foo", "/A/prefixD.foo", "/A/prefixE.foo", "/A/prefixF.foo",
		"/B/sub/prefixA.foo", "/B/sub/prefixC.foo", "/B/sub/prefixD.foo", "/B/sub/prefixE.foo", "/B/sub/prefixF.foo",
		"/C/prefixA.foo", "/C/prefixD.foo", "/C/prefixE.foo", "/C/prefixF.foo",
		"/D/prefixA", "/D/prefixB", "/D/prefixC", "/D/prefixD", "/D/prefixE", "/D/prefixF",
		"/E/subtestA.foo", "/F/subtestA.foo", "/G/subtestA.foo", "/H/subtestA.foo", "/I/J/subtestA.foo", "/K/dirJ/subtestA.foo",
		"/parameters.yml", "/parameters.yml.dist", "/prefixB.foo", "/prefixD.foo", "/prefixE.foo", "/prefixF.foo", "/toplevelA.foo", "/toplevelB.foo"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestSkipExcludes(t *testing.T) {
	src := setupTree(t)
	got := archivable(t, src, []string{"prefixB.foo"}, true)
	var want []string
	for _, f := range fileTree {
		want = append(want, "/"+f)
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestGitExcludes(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	src := setupTree(t)
	attrs := strings.Join([]string{"", "# gitattributes rules with comments and blank lines", "prefixB.foo export-ignore", "/prefixA.foo export-ignore", "prefixC.* export-ignore", "",
		"prefixE.foo export-ignore", "# and more", "# comments", "", "/prefixE.foo -export-ignore", "/prefixD.foo export-ignore", "prefixF.* export-ignore",
		"/*/*/prefixF.foo -export-ignore", "", "refixD.foo export-ignore", "/C export-ignore", "D/prefixA export-ignore", "E export-ignore", "F/ export-ignore",
		"G/* export-ignore", "J/ export-ignore", "parameters.yml export-ignore", `\!important!.txt export-ignore`, `\#* export-ignore`}, "\n")
	if err := os.WriteFile(filepath.Join(src, ".gitattributes"), []byte(attrs), 0o666); err != nil {
		t.Fatal(err)
	}
	got := archivable(t, src, nil, false)
	cmd := exec.Command("sh", "-c", `git init -q && git config user.email "you@example.com" && git config user.name "Your Name" && git config commit.gpgsign false && git add .git* && git commit -q -m "ignore rules" && git add . && git commit -q -m "init" && git archive --format=zip --prefix=archive/ -o archive.zip HEAD`)
	cmd.Dir = src
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git failed: %v %s", err, out)
	}
	zr, err := zip.OpenReader(filepath.Join(src, "archive.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var want []string
	for _, f := range zr.File {
		if !strings.HasSuffix(f.Name, "/") {
			want = append(want, strings.TrimPrefix(f.Name, "archive"))
		}
	}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestGlobToRegex(t *testing.T) {
	cases := map[string]string{
		"*.php":       `#^(?=[^\.])[^/]*\.php$#`,
		"/foo/**/bar": `#^(?=[^\.])/(?=[^\.])foo/(?:(?=[^\.])[^/]++/)*(?=[^\.])bar$#`,
		"{a,b}c?":     `#^(?=[^\.])(a|b)c[^/]$#`,
		".git":        `#^\.git$#`,
		`\*literal`:   `#^(?=[^\.])\*literal$#`,
		"docs/**":     `#^(?=[^\.])docs/(?:(?=[^\.])[^/]++/?)*$#`,
		"a(b)|c+^$#":  `#^(?=[^\.])a\(b\)\|c\+\^\$\#$#`,
	}
	for glob, want := range cases {
		if got := GlobToRegex(glob, true, true, "#"); got != want {
			t.Errorf("GlobToRegex(%q) = %q, want %q", glob, got, want)
		}
	}
}

func TestZipExtraField(t *testing.T) {
	// the ASi extra field PharData writes for a 0644 file
	asi := []byte{0xa4, 0x01}
	if got := crc32.ChecksumIEEE(asi); got != 0xfcb55b84 {
		t.Errorf("crc = %x", got)
	}
	data := buildZip([]*entry{{name: "a.txt", mode: 0o644, content: []byte("x\n")}, {name: "d", dir: true, mode: 0o777}})
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	if len(zr.File) != 2 || zr.File[0].Name != "a.txt" || zr.File[1].Name != "d/" {
		t.Fatalf("unexpected entries %v", zr.File)
	}
}
