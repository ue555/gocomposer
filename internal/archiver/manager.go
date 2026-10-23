package archiver

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/fsutil"
)

// PackageInfo is what ArchiveManager reads from a CompletePackageInterface.
type PackageInfo struct {
	Name            string
	PrettyVersion   string
	ArchiveName     string
	ArchiveExcludes []string
	DistReference   string
	DistType        string
	SourceReference string
}

var (
	reNonNameChar = regexp.MustCompile(`(?i)[^a-z0-9-_]`)
	reSha1Ref     = regexp.MustCompile(`^[a-f0-9]{40}$`)
)

type part struct{ key, value string }

// FilenameParts is ArchiveManager::getPackageFilenameParts.
func FilenameParts(p PackageInfo) []part {
	base := p.ArchiveName
	if base == "" {
		base = reNonNameChar.ReplaceAllString(p.Name, "-")
	}
	parts := []part{{"base", base}}
	if p.DistReference != "" && reSha1Ref.MatchString(p.DistReference) {
		parts = append(parts, part{"dist_reference", p.DistReference})
		if p.DistType != "" {
			parts = append(parts, part{"dist_type", p.DistType})
		}
	} else {
		parts = append(parts, part{"version", p.PrettyVersion})
		if p.DistReference != "" {
			parts = append(parts, part{"dist_reference", p.DistReference})
		}
	}
	if p.SourceReference != "" {
		sum := sha1.Sum([]byte(p.SourceReference))
		parts = append(parts, part{"source_reference", hex.EncodeToString(sum[:])[:6]})
	}
	for i := range parts {
		parts[i].value = strings.ReplaceAll(parts[i].value, "/", "-")
	}
	return parts
}

// PackageFilename is ArchiveManager::getPackageFilename.
func PackageFilename(p PackageInfo) string { return filenameFromParts(FilenameParts(p)) }

func filenameFromParts(parts []part) string {
	var s []string
	for _, p := range parts {
		s = append(s, p.value)
	}
	return strings.Join(s, "-")
}

// buildExcludePatterns is ArchiveManager::buildExcludePatterns.
func buildExcludePatterns(parts []part) []string {
	base := parts[0].value
	if len(parts) > 1 {
		base += "-*"
	}
	var out []string
	for _, f := range Formats {
		out = append(out, base+"."+f)
	}
	return out
}

// TempName is sys_get_temp_dir().'/composer_archive'.bin2hex(random_bytes(5)).
func TempName() string {
	var b [5]byte
	_, _ = rand.Read(b[:])
	return strings.TrimRight(os.TempDir(), "/") + "/composer_archive" + hex.EncodeToString(b[:])
}

// DirError is Filesystem::ensureDirectoryExists failing.
type DirError struct{ Msg string }

func (e *DirError) Error() string { return e.Msg }

// EnsureDirectoryExists is Filesystem::ensureDirectoryExists.
func EnsureDirectoryExists(dir string) error {
	st, err := os.Stat(dir)
	if err == nil {
		if !st.IsDir() {
			return &DirError{Msg: dir + " exists and is not a directory."}
		}
		return nil
	}
	if os.MkdirAll(dir, 0o777) != nil {
		return &DirError{Msg: dir + " does not exist and could not be created: "}
	}
	return nil
}

// Archive is ArchiveManager::archive once the sources are available in
// sourcePath (the project directory for the root package, a download
// directory otherwise). It returns the path of the created archive.
func ArchivePackage(p PackageInfo, sourcePath, format, targetDir, fileName string, ignoreFilters bool) (string, error) {
	parts := FilenameParts(p)
	if fileName != "" {
		parts = []part{{"base", fileName}}
	}
	packageName := filenameFromParts(parts)
	excludes := buildExcludePatterns(parts)
	if err := EnsureDirectoryExists(targetDir); err != nil {
		return "", err
	}
	target := fsutil.Realpath(targetDir) + "/" + packageName + "." + format
	if err := EnsureDirectoryExists(filepath.Dir(target)); err != nil {
		return "", err
	}
	tempTarget := TempName() + "." + format
	if err := EnsureDirectoryExists(filepath.Dir(tempTarget)); err != nil {
		return "", err
	}
	archivePath, err := Archive(sourcePath, tempTarget, format, append(excludes, p.ArchiveExcludes...), ignoreFilters)
	if err != nil {
		return "", err
	}
	if err := rename(archivePath, target); err != nil {
		return "", err
	}
	os.Remove(tempTarget)
	return target, nil
}

// rename is Filesystem::rename (falling back to copy+delete across devices).
func rename(from, to string) error {
	if err := os.Rename(from, to); err == nil {
		return nil
	}
	if err := fsutil.CopyFile(from, to); err != nil {
		return err
	}
	return os.Remove(from)
}
