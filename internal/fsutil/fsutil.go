// Package fsutil ports the path helpers of Composer\Util\Filesystem and PHP's
// var_export() for strings, used when generating PHP files.
package fsutil

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// IsAbsolutePath is Filesystem::isAbsolutePath.
func IsAbsolutePath(p string) bool {
	return strings.HasPrefix(p, "/") || (len(p) > 1 && p[1] == ':') || strings.HasPrefix(p, `\\`)
}

var (
	rePrefix  = regexp.MustCompile(`(?i)^([0-9a-z]{2,}:(?://(?:[a-z]:)?)?|[a-z]:)`)
	reDriveUp = regexp.MustCompile(`(?i)(^|://)[a-z]:$`)
)

// NormalizePath is Filesystem::normalizePath.
func NormalizePath(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	prefix := ""
	absolute := ""
	if strings.HasPrefix(p, "//") && len(p) > 2 {
		absolute = "//"
		p = p[2:]
	}
	if m := rePrefix.FindStringSubmatch(p); m != nil {
		prefix = m[1]
		p = p[len(prefix):]
	}
	if strings.HasPrefix(p, "/") {
		absolute = "/"
		p = p[1:]
	}
	var parts []string
	up := false
	for _, chunk := range strings.Split(p, "/") {
		if chunk == ".." && (len(absolute) > 0 || up) {
			if len(parts) > 0 {
				parts = parts[:len(parts)-1]
			}
			up = !(len(parts) == 0 || parts[len(parts)-1] == "..")
		} else if chunk != "." && chunk != "" {
			parts = append(parts, chunk)
			up = chunk != ".."
		}
	}
	prefix = reDriveUp.ReplaceAllStringFunc(prefix, strings.ToUpper)
	return prefix + absolute + strings.Join(parts, "/")
}

func phpDirname(p string) string {
	d := path.Dir(p)
	return d
}

var reDriveRoot = regexp.MustCompile(`(?i)^[A-Z]:/?$`)

// FindShortestPath is Filesystem::findShortestPath.
func FindShortestPath(from, to string, directories, preferRelative bool) (string, error) {
	if !IsAbsolutePath(from) || !IsAbsolutePath(to) {
		return "", fmt.Errorf("$from (%s) and $to (%s) must be absolute paths.", from, to)
	}
	from = NormalizePath(from)
	to = NormalizePath(to)
	if directories {
		from = strings.TrimRight(from, "/") + "/dummy_file"
	}
	if phpDirname(from) == phpDirname(to) {
		return "./" + path.Base(to), nil
	}
	common := to
	for !strings.HasPrefix(from+"/", common+"/") && common != "/" && !reDriveRoot.MatchString(common) {
		common = phpDirname(common)
	}
	if !strings.HasPrefix(from, common) {
		return to, nil
	}
	common = strings.TrimRight(common, "/") + "/"
	depth := strings.Count(safeSub(from, len(common)), "/")
	code := strings.Repeat("../", depth)
	if !preferRelative && common == "/" && depth > 1 {
		return to, nil
	}
	result := code + safeSub(to, len(common))
	if result == "" {
		return "./", nil
	}
	return result, nil
}

func safeSub(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[n:]
}

// FindShortestPathCode is Filesystem::findShortestPathCode.
func FindShortestPathCode(from, to string, directories, staticCode, preferRelative bool) (string, error) {
	if !IsAbsolutePath(from) || !IsAbsolutePath(to) {
		return "", fmt.Errorf("$from (%s) and $to (%s) must be absolute paths.", from, to)
	}
	from = NormalizePath(from)
	to = NormalizePath(to)
	if from == to {
		if directories {
			return "__DIR__", nil
		}
		return "__FILE__", nil
	}
	common := to
	for !strings.HasPrefix(from+"/", common+"/") && common != "/" && !reDriveRoot.MatchString(common) && common != "." {
		common = phpDirname(common)
	}
	if !strings.HasPrefix(from, common) || common == "." {
		return VarExport(to), nil
	}
	common = strings.TrimRight(common, "/") + "/"
	if strings.HasPrefix(to, from+"/") {
		return "__DIR__ . " + VarExport(to[len(from):]), nil
	}
	depth := strings.Count(safeSub(from, len(common)), "/")
	if directories {
		depth++
	}
	if !preferRelative && common == "/" && depth > 1 {
		return VarExport(to), nil
	}
	var code string
	if staticCode {
		code = "__DIR__ . '" + strings.Repeat("/..", depth) + "'"
	} else {
		code = strings.Repeat("dirname(", depth) + "__DIR__" + strings.Repeat(")", depth)
	}
	rel := safeSub(to, len(common))
	if rel != "" {
		return code + "." + VarExport("/"+rel), nil
	}
	return code, nil
}

// VarExport is var_export() of a PHP string.
func VarExport(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	// var_export emits NUL bytes as '' . "\0" . ''
	s = strings.ReplaceAll(s, "\x00", `' . "\0" . '`)
	return "'" + s + "'"
}

// IsLocalPath is Filesystem::isLocalPath (non-Windows).
func IsLocalPath(p string) bool {
	return regexp.MustCompile(`(?i)^(file://|/|/?[a-z]:[\\/]|\.\.[\\/]|[a-z0-9_.-]+[\\/])`).MatchString(p)
}

// GetPlatformPath is Filesystem::getPlatformPath (non-Windows).
func GetPlatformPath(p string) string {
	return regexp.MustCompile(`(?i)^file://`).ReplaceAllString(p, "")
}

// Realpath resolves symlinks and makes p absolute; it returns "" when p does not exist.
func Realpath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	r, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return ""
	}
	return r
}

// FilePutContentsIfModified writes data when the file content differs.
func FilePutContentsIfModified(p string, data []byte) error {
	if cur, err := os.ReadFile(p); err == nil && string(cur) == string(data) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
		return err
	}
	return WriteFileAtomic(p, data, 0o666)
}

// WriteFileAtomic writes through a temp file in the same directory.
func WriteFileAtomic(p string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(p)
	f, err := os.CreateTemp(dir, "."+filepath.Base(p)+".tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if st, err := os.Stat(p); err == nil {
		perm = st.Mode().Perm()
	} else {
		perm = applyUmask(perm)
	}
	_ = os.Chmod(tmp, perm)
	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// CopyFile copies a regular file preserving the permission bits.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	st, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o777); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, st.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// RemoveAll removes a path; symlinks are removed, not followed.
func RemoveAll(p string) error {
	if _, err := os.Lstat(p); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(p)
}

// IsDirEmpty reports whether a directory has no entries.
func IsDirEmpty(p string) bool {
	entries, err := os.ReadDir(p)
	return err == nil && len(entries) == 0
}

// WithinDir reports whether target is inside base after cleaning (no symlink resolution).
func WithinDir(base, target string) bool {
	base = filepath.Clean(base)
	target = filepath.Clean(target)
	if base == target {
		return true
	}
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, "../")
}
