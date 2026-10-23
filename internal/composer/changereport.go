package composer

// Ports of FileDownloader::getLocalChanges, Composer\Package\Comparer\Comparer
// and VcsDownloader/PathDownloader::getVcsReference (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
// Jordi Boggiano), used by the status command.

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// DownloaderType is the DownloadManager downloader name used for an
// installed package ("" for metapackages).
func (im *InstallationManager) DownloaderType(p *pkg.Package) string { return im.downloaderType(p) }

// FileLocalChanges is FileDownloader::getLocalChanges: the package is
// downloaded and installed again next to path (path + "_compare") with a
// silent IO and both trees are compared. It returns "" when nothing changed.
// Failures are reported in the result like Composer (or returned at debug
// verbosity, where Composer rethrows them).
func (im *InstallationManager) FileLocalChanges(p *pkg.Package, path string) (string, error) {
	target := strings.TrimRight(path, "/")
	cmp := target + "_compare"
	quiet := *im.c
	quiet.IO = &cio.IO{Out: io.Discard, Err: io.Discard, Verbosity: cio.Quiet}
	qim := NewInstallationManager(&quiet)
	qim.InstallPathOverride = func(*pkg.Package) string { return cmp }
	output := ""
	err := func() error {
		if isDir(cmp) {
			if err := fsutil.RemoveAll(cmp); err != nil {
				return err
			}
		}
		typ := im.downloaderType(p)
		switch typ {
		case "path":
			if err := qim.pathDownload(p); err != nil {
				return err
			}
		case "zip", "tar", "gzip", "file", "phar":
			transfer, err := qim.fileDownload(p)
			if err != nil {
				return err
			}
			if transfer != nil {
				if err := transfer(); err != nil {
					return err
				}
			}
		default:
			return &UserError{Msg: "gocomposer cannot detect local changes of " + p.PrettyName + ": the " + typ + " downloader is not supported", Code: 1}
		}
		var err error
		switch typ {
		case "path":
			err = qim.pathInstall(p, cmp, false)
		case "zip", "tar", "gzip":
			err = qim.archiveInstall(p, cmp, false)
		default:
			err = qim.fileInstall(p, cmp, false)
		}
		if err != nil {
			return err
		}
		output = compareTrees(cmp, target)
		return removeDirectory(cmp)
	}()
	if err != nil {
		if im.c.IO.IsDebug() {
			return "", err
		}
		return "Failed to detect changes: [" + exceptionClass(err) + "] " + err.Error(), nil
	}
	return strings.TrimSpace(output), nil
}

// removeDirectory is Filesystem::removeDirectory (a symlinked directory is
// only unlinked).
func removeDirectory(dir string) error {
	if st, err := os.Lstat(dir); err == nil && st.Mode()&os.ModeSymlink != 0 {
		return os.Remove(dir)
	}
	return fsutil.RemoveAll(dir)
}

func exceptionClass(err error) string {
	var te *httpdl.TransportError
	if errors.As(err, &te) {
		return `Composer\Downloader\TransportException`
	}
	return "RuntimeException"
}

// treeEntry is one file of Comparer::doTree.
type treeEntry struct{ file, hash string }

// doTree is Comparer::doTree: files (non-empty) and symlinks per directory,
// in readdir order.
type tree struct {
	dirs  []string
	files map[string][]treeEntry
}

func (t *tree) add(dir string, e treeEntry) {
	if _, ok := t.files[dir]; !ok {
		t.dirs = append(t.dirs, dir)
	}
	t.files[dir] = append(t.files[dir], e)
}

func (t *tree) get(dir, file string) (string, bool) {
	for _, e := range t.files[dir] {
		if e.file == file {
			return e.hash, true
		}
	}
	return "", false
}

func doTree(base, dir string, t *tree) bool {
	f, err := os.Open(filepath.Join(base, dir))
	if err != nil {
		return false
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return false
	}
	for _, name := range names {
		rel := dir + "/" + name
		full := filepath.Join(base, rel)
		st, err := os.Lstat(full)
		if err != nil {
			continue
		}
		if st.Mode()&os.ModeSymlink != 0 {
			l, _ := os.Readlink(full)
			t.add(dir, treeEntry{name, l})
			continue
		}
		if st.IsDir() {
			if !doTree(base, rel, t) {
				return false
			}
			continue
		}
		if st.Mode().IsRegular() && st.Size() > 0 {
			b, err := os.ReadFile(full)
			if err != nil {
				continue
			}
			sum := sha1.Sum(b)
			t.add(dir, treeEntry{name, hex.EncodeToString(sum[:])})
		}
	}
	return true
}

// compareTrees is Comparer::doCompare + getChangedAsString(true).
func compareTrees(source, update string) string {
	src := &tree{files: map[string][]treeEntry{}}
	if !doTree(source, ".", src) {
		return ""
	}
	dst := &tree{files: map[string][]treeEntry{}}
	if !doTree(update, ".", dst) {
		return ""
	}
	var sections []string
	changed := map[string][]string{}
	add := func(section, item string) {
		if _, ok := changed[section]; !ok {
			sections = append(sections, section)
		}
		changed[section] = append(changed[section], item)
	}
	for _, dir := range src.dirs {
		for _, e := range src.files[dir] {
			if h, ok := dst.get(dir, e.file); ok {
				if h != e.hash {
					add("changed", dir+"/"+e.file)
				}
			} else {
				add("removed", dir+"/"+e.file)
			}
		}
	}
	for _, dir := range dst.dirs {
		for _, e := range dst.files[dir] {
			if _, ok := src.get(dir, e.file); !ok {
				add("added", dir+"/"+e.file)
			}
		}
	}
	var items []string
	for _, s := range sections {
		for _, item := range changed[s] {
			items = append(items, item+"\r\n")
		}
	}
	return strings.TrimSpace(strings.Join(items, "\r\n"))
}

// GuessVersionCommit is VersionGuesser::guessVersion reduced to the
// pretty version and commit (VcsDownloader::getVcsReference and the status
// command). ok is false when no version could be guessed.
func GuessVersionCommit(dir string, p *pkg.Package, c *Composer) (pretty, commit string, ok bool) {
	conf := phpjson.NewMap()
	conf.Set("version", p.PrettyVersion())
	if extra := p.Extra(); extra != nil {
		conf.Set("extra", extra)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", "", false
	}
	pretty, _, ok = GuessVersion(abs, conf, c.Config)
	if !ok {
		return "", "", false
	}
	return pretty, gitCommit(abs), true
}
