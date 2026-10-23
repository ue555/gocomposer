package composer

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ExtractError reports an unsafe or corrupt archive.
type ExtractError struct{ Msg string }

func (e *ExtractError) Error() string { return e.Msg }

// safeJoin returns dest/name after verifying the entry stays inside dest.
func safeJoin(dest, name string) (string, error) {
	n := strings.ReplaceAll(name, `\`, "/")
	if strings.HasPrefix(n, "/") || (len(n) > 1 && n[1] == ':') {
		return "", &ExtractError{fmt.Sprintf("archive entry %q has an absolute path", name)}
	}
	clean := path.Clean("/" + n)
	for _, part := range strings.Split(n, "/") {
		if part == ".." {
			return "", &ExtractError{fmt.Sprintf("archive entry %q escapes the extraction directory", name)}
		}
	}
	return filepath.Join(dest, filepath.FromSlash(strings.TrimPrefix(clean, "/"))), nil
}

// checkNoSymlinkParents verifies that no existing parent of target (below
// dest) is a symlink, so writes cannot be redirected outside dest.
func checkNoSymlinkParents(dest, target string) error {
	rel, err := filepath.Rel(dest, filepath.Dir(target))
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	cur := dest
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		cur = filepath.Join(cur, part)
		st, err := os.Lstat(cur)
		if err != nil {
			return nil
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return &ExtractError{fmt.Sprintf("archive entry %q is placed below a symlink", strings.TrimPrefix(target, dest))}
		}
	}
	return nil
}

// symlinkInside reports whether a link target resolves inside root.
func symlinkInside(root, linkPath, target string) bool {
	if filepath.IsAbs(target) {
		return false
	}
	resolved := filepath.Clean(filepath.Join(filepath.Dir(linkPath), target))
	rel, err := filepath.Rel(root, resolved)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

const maxEntrySize = 4 << 30

func writeEntry(dest, target string, r io.Reader, mode os.FileMode) error {
	if err := checkNoSymlinkParents(dest, target); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
		return err
	}
	if st, err := os.Lstat(target); err == nil && st.Mode()&os.ModeSymlink != 0 {
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, io.LimitReader(r, maxEntrySize)); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Chmod(target, mode)
}

// unzipEntryName mirrors Info-ZIP unzip (which Composer prefers for zip
// dists): leading "/" or drive letters and "." / ".." path components are
// dropped, so an entry like "../../x" is extracted as "x" inside dest instead
// of failing the install.
func unzipEntryName(name string) string {
	n := strings.ReplaceAll(name, `\`, "/")
	if len(n) > 1 && n[1] == ':' {
		n = n[2:]
	}
	var parts []string
	for _, part := range strings.Split(n, "/") {
		if part == "" || part == "." || part == ".." {
			continue
		}
		parts = append(parts, part)
	}
	out := strings.Join(parts, "/")
	if strings.HasSuffix(n, "/") && out != "" {
		out += "/"
	}
	return out
}

// ExtractZip extracts a zip archive into dest (which must exist).
func ExtractZip(file, dest string) error {
	zr, err := zip.OpenReader(file)
	if err != nil {
		return &ExtractError{fmt.Sprintf("the file is not a valid zip archive (%v)", err)}
	}
	defer zr.Close()
	var links []*zip.File
	for _, f := range zr.File {
		name := unzipEntryName(f.Name)
		if name == "" {
			continue
		}
		target, err := safeJoin(dest, name)
		if err != nil {
			return err
		}
		mode := f.Mode()
		switch {
		case mode&os.ModeSymlink != 0:
			links = append(links, f)
			continue
		case f.FileInfo().IsDir() || strings.HasSuffix(f.Name, "/"):
			if err := checkNoSymlinkParents(dest, target); err != nil {
				return err
			}
			if err := os.MkdirAll(target, 0o777); err != nil {
				return err
			}
			continue
		}
		perm := zipPerm(f)
		rc, err := f.Open()
		if err != nil {
			return &ExtractError{fmt.Sprintf("could not read %s from archive: %v", f.Name, err)}
		}
		err = writeEntry(dest, target, rc, perm)
		rc.Close()
		if err != nil {
			return err
		}
	}
	for _, f := range links {
		target, _ := safeJoin(dest, unzipEntryName(f.Name))
		rc, err := f.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, io.LimitReader(rc, 4096))
		rc.Close()
		if err != nil {
			return err
		}
		if err := makeSymlink(dest, target, buf.String()); err != nil {
			return err
		}
	}
	return nil
}

func zipPerm(f *zip.File) os.FileMode {
	// Unix permissions are stored in the high 16 bits when created on Unix
	if f.CreatorVersion>>8 == 3 { // Unix
		p := os.FileMode(f.ExternalAttrs>>16) & 0o777
		if p != 0 {
			return p
		}
	}
	return 0o666 &^ currentUmask()
}

func makeSymlink(dest, target, linkTarget string) error {
	if !symlinkInside(dest, target, linkTarget) {
		return &ExtractError{fmt.Sprintf("archive symlink %q points outside of the package (%s)", strings.TrimPrefix(target, dest+string(filepath.Separator)), linkTarget)}
	}
	if err := checkNoSymlinkParents(dest, target); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
		return err
	}
	_ = os.Remove(target)
	return os.Symlink(linkTarget, target)
}

// ExtractTar extracts tar, tar.gz or tar.bz2 archives.
func ExtractTar(file, dest string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	var r io.Reader = f
	head := make([]byte, 3)
	n, _ := io.ReadFull(f, head)
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	switch {
	case n >= 2 && head[0] == 0x1f && head[1] == 0x8b:
		gz, err := gzip.NewReader(f)
		if err != nil {
			return &ExtractError{err.Error()}
		}
		defer gz.Close()
		r = gz
	case n >= 3 && string(head) == "BZh":
		r = bzip2.NewReader(f)
	}
	tr := tar.NewReader(r)
	type link struct{ target, name string }
	var links []link
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &ExtractError{fmt.Sprintf("invalid tar archive: %v", err)}
		}
		if h.Name == "pax_global_header" || h.Typeflag == tar.TypeXGlobalHeader {
			continue
		}
		target, err := safeJoin(dest, h.Name)
		if err != nil {
			return err
		}
		switch h.Typeflag {
		case tar.TypeDir:
			if err := checkNoSymlinkParents(dest, target); err != nil {
				return err
			}
			if err := os.MkdirAll(target, 0o777); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			mode := os.FileMode(h.Mode) & 0o777
			if mode == 0 {
				mode = 0o644
			}
			if err := writeEntry(dest, target, tr, mode); err != nil {
				return err
			}
		case tar.TypeSymlink:
			links = append(links, link{target, h.Linkname})
		case tar.TypeLink:
			src, err := safeJoin(dest, h.Linkname)
			if err != nil {
				return err
			}
			in, err := os.Open(src)
			if err != nil {
				return &ExtractError{fmt.Sprintf("hard link %s points to a missing entry", h.Name)}
			}
			st, _ := in.Stat()
			err = writeEntry(dest, target, in, st.Mode().Perm())
			in.Close()
			if err != nil {
				return err
			}
		default:
			// devices, fifos etc. are skipped like PharData does not create them
		}
	}
	for _, l := range links {
		if err := makeSymlink(dest, l.target, l.name); err != nil {
			return err
		}
	}
	return nil
}

// ExtractGzip decompresses a single gzip file to target.
func ExtractGzip(file, target string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return &ExtractError{err.Error()}
	}
	defer gz.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, io.LimitReader(gz, maxEntrySize)); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// ZipHasEndOfCentralDirectory is ZipDownloader::validateDownloadedFile.
func ZipHasEndOfCentralDirectory(file string) (bool, int64) {
	st, err := os.Stat(file)
	if err != nil {
		return true, 0
	}
	size := st.Size()
	f, err := os.Open(file)
	if err != nil {
		return true, size
	}
	defer f.Close()
	off := size - 22 - 65535
	if off < 0 {
		off = 0
	}
	if _, err := f.Seek(off, io.SeekStart); err != nil {
		return true, size
	}
	tail, _ := io.ReadAll(f)
	return bytes.Contains(tail, []byte("PK\x05\x06")), size
}
