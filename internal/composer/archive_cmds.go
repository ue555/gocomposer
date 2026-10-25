// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Downloader/XzDownloader.php, RarDownloader.php and the
// unzip/7z/ZipArchive selection of ZipDownloader.php (download,
// extractWithSystemUnzip, extractWithZipArchive, getErrorMessage).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Like Composer, the xz and rar dists are extracted by external commands
// (tar -xJf, unrar) and zip dists by 7z when no unzip command exists.
// gocomposer has no PHP rar extension (RarArchive): when unrar fails it
// reports Composer's error for a PHP without ext-rar, and a PHP runtime that
// has ext-rar loaded is refused explicitly instead of being imitated.

package composer

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/pkg"
	"gocomposer/internal/vcs"
)

// unzipCommand is one ZipDownloader::$unzipCommands entry.
type unzipCommand struct {
	exe  string // unzip, 7z, 7zz or 7za
	path string // the resolved executable
}

// unzipCommands is ZipDownloader::$unzipCommands (static, computed once).
var unzipCommands = sync.OnceValue(func() []unzipCommand {
	var cmds []unzipCommand
	if p, err := exec.LookPath("unzip"); err == nil {
		cmds = append(cmds, unzipCommand{"unzip", p})
	}
	if runtime.GOOS != "windows" {
		// 7z linux/macOS support is only used if unzip is not present
		for _, c := range []string{"7z", "7zz", "7za"} {
			if p, err := exec.LookPath(c); err == nil {
				cmds = append(cmds, unzipCommand{c, p})
				break
			}
		}
	}
	return cmds
})

var (
	zipWarnOnce   sync.Once
	sevenZipWarn  sync.Once
	re7ZipVersion = regexp.MustCompile(`^\s*7-Zip(?: \[64\])? ([0-9.]+)`)
)

// hasZipArchive is class_exists('ZipArchive') of the PHP runtime.
func (im *InstallationManager) hasZipArchive() bool {
	_, facts, _ := im.c.PHP()
	return PlatformInfo{facts}.ExtensionLoaded("zip")
}

func (im *InstallationManager) iniMessage() string {
	_, facts, _ := im.c.PHP()
	paths := PlatformInfo{facts}.IniPaths()
	if len(paths) > 0 && paths[0] == "" {
		paths = paths[1:]
	}
	if len(paths) == 0 || paths[0] == "" {
		return "A php.ini file does not exist. You will have to create one."
	}
	if len(paths) > 1 {
		return "Your command-line PHP is using multiple ini files. Run `php --ini` to show them."
	}
	return "The php.ini used by your command-line PHP is: " + paths[0]
}

// zipDownloadCheck is the part of ZipDownloader::download run before the
// file download: it fails when neither ZipArchive nor an unzip/7z command is
// available, and warns once when only ZipArchive is.
func (im *InstallationManager) zipDownloadCheck() error {
	cmds := unzipCommands()
	if len(cmds) > 0 {
		return nil
	}
	if !im.hasZipArchive() {
		return &UserError{Msg: "The zip extension and unzip/7z commands are both missing, skipping.\n" + im.iniMessage(), Code: 1}
	}
	if runtime.GOOS != "windows" {
		zipWarnOnce.Do(func() {
			im.c.IO.WriteError("<warning>As there is no 'unzip' nor '7z' command installed zip files are being unpacked using the PHP zip extension.</warning>")
			im.c.IO.WriteError("<warning>This may cause invalid reports of corrupted archives. Besides, any UNIX permissions (e.g. executable) defined in the archives will be lost.</warning>")
			im.c.IO.WriteError("<warning>Installing 'unzip' or '7z' (21.01+) may remediate them.</warning>")
		})
	}
	return nil
}

func (im *InstallationManager) process() *vcs.Process {
	timeout := 300 * time.Second
	if t, ok := im.c.Config.Get("process-timeout").(int64); ok {
		timeout = time.Duration(t) * time.Second
	}
	return vcs.NewProcess(im.c.IO, timeout)
}

// extractZipWith is ZipDownloader::extract: the native extraction stands in
// for unzip (as before), 7z is executed like Composer does and ZipArchive is
// imitated when no command exists.
func (im *InstallationManager) extractZipWith(p *pkg.Package, file, path string) error {
	cmds := unzipCommands()
	if len(cmds) == 0 {
		return im.extractWithZipArchive(p, file, path)
	}
	if cmds[0].exe == "unzip" {
		if err := ExtractZip(file, path); err != nil {
			return &UserError{Msg: "Failed to extract " + p.Name + ": " + err.Error(), Code: 1}
		}
		return nil
	}
	isLastChance := !im.hasZipArchive()
	c := cmds[0]
	if runtime.GOOS != "windows" {
		sevenZipWarn.Do(func() {
			code, out, _ := im.process().Execute([]string{c.path}, "")
			if code == 0 {
				if m := re7ZipVersion.FindStringSubmatch(out); m != nil && versionLess(m[1], "21.01") {
					im.c.IO.WriteError("    <warning>Unzipping using " + c.exe + " " + m[1] + " may result in incorrect file permissions. Install " + c.exe + " 21.01+ or unzip to ensure you get correct permissions.</warning>")
				}
			}
		})
	}
	command := []string{c.path, "x", "-bb0", "-y", file, "-o" + path}
	proc := im.process()
	code, _, err := proc.Execute(command, "")
	var processError error
	if err != nil {
		processError = err
	} else if code != 0 {
		output := strings.ReplaceAll(proc.ErrorOutput(), ", "+file+".zip or "+file+".ZIP", "")
		processError = &UserError{Msg: fmt.Sprintf("Failed to extract %s: (%d) %s\n\n%s", p.Name, code, strings.Join(command, " "), output), Code: 1}
	} else {
		return nil
	}
	if isLastChance || strings.Contains(processError.Error(), "zip bomb") {
		return processError
	}
	if _, err := os.Stat(file); err != nil {
		im.c.IO.WriteError("    <warning>" + processError.Error() + "</warning>")
		im.c.IO.WriteError("    <warning>This most likely is due to a custom installer plugin not handling the returned Promise from the downloader</warning>")
		im.c.IO.WriteError("    <warning>See https://github.com/composer/installers/commit/5006d0c28730ade233a8f42ec31ac68fb1c5c9bb for an example fix</warning>")
	} else {
		im.c.IO.WriteError("    <warning>" + processError.Error() + "</warning>")
		im.c.IO.WriteError("    The archive may contain identical file names with different capitalization (which fails on case insensitive filesystems)")
		im.c.IO.WriteError("    Unzip with " + c.exe + " command failed, falling back to ZipArchive class")
	}
	return im.extractWithZipArchive(p, file, path)
}

// versionLess is version_compare($a, $b, '<') for dotted numeric versions.
func versionLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			fmt.Sscanf(pa[i], "%d", &x)
		}
		if i < len(pb) {
			fmt.Sscanf(pb[i], "%d", &y)
		}
		if x != y {
			return x < y
		}
	}
	return false
}

// extractWithZipArchive imitates ZipDownloader::extractWithZipArchive (PHP's
// ZipArchive::extractTo): UNIX permissions and symlinks of the archive are
// not restored, entries are written as plain files.
func (im *InstallationManager) extractWithZipArchive(p *pkg.Package, file, path string) error {
	st, err := os.Stat(file)
	if err != nil || st.Size() == 0 {
		return &UserError{Msg: fmt.Sprintf("'%s' is a corrupted zip archive (0 bytes), try again.", file), Code: 1}
	}
	zr, err := zip.OpenReader(file)
	if err != nil {
		// ZipArchive::ER_NOZIP
		return &UserError{Msg: fmt.Sprintf("'%s' is not a zip archive.", file), Code: 19}
	}
	defer zr.Close()
	var totalSize uint64
	inspectAll := false
	for _, f := range zr.File {
		if f.UncompressedSize64 > f.CompressedSize64*200 {
			inspectAll = true
		}
	}
	if inspectAll {
		for _, f := range zr.File {
			totalSize += f.UncompressedSize64
		}
		if totalSize > uint64(st.Size())*100 && totalSize > 50*1024*1024 {
			return &UserError{Msg: `Invalid zip file for "` + p.Name + `" with compression ratio >99% (possible zip bomb)`, Code: 1}
		}
	}
	// ext-zip before 1.22.4 (PHP 8.3) does not check the result of reading
	// an entry in ZipArchive::extractTo, newer versions (PHP 8.4+) warn
	// "Cannot extract" (observed with ext-zip 1.22.3 and 1.22.8)
	lenientCRC := false
	if _, facts, _ := im.c.PHP(); facts != nil {
		if v, ok := facts.ExtVersions["zip"]; ok && versionLess(v, "1.22.4") {
			lenientCRC = true
		}
	}
	extractErr := &UserError{Msg: `There was an error extracting the ZIP file for "` + p.Name + `", it is either corrupted or using an invalid format.`, Code: 1}
	for _, f := range zr.File {
		name := unzipEntryName(f.Name)
		if name == "" {
			continue
		}
		target, err := safeJoin(path, name)
		if err != nil {
			return extractErr
		}
		if strings.HasSuffix(name, "/") {
			if err := checkNoSymlinkParents(path, target); err != nil {
				return err
			}
			if err := os.MkdirAll(target, 0o777); err != nil {
				return extractErr
			}
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return extractErr
		}
		err = writeEntry(path, target, rc, 0o666&^currentUmask())
		rc.Close()
		if err != nil {
			// libzip's read error is raised as a PHP warning (ErrorException)
			// that ZipDownloader wraps
			reason := ""
			switch {
			case errors.Is(err, zip.ErrChecksum) && lenientCRC:
				// the entry was written completely, only its CRC is wrong
				continue
			case errors.Is(err, zip.ErrChecksum):
				reason = "CRC error"
			case errors.Is(err, zip.ErrFormat), errors.Is(err, io.ErrUnexpectedEOF):
				reason = "Read error"
			default:
				return err
			}
			inner := "ZipArchive::extractTo(): Cannot extract \"" + f.Name + "\": \"" + reason + "\""
			return &UserError{Msg: `The archive for "` + p.Name + `" may contain identical file names with different capitalization (which fails on case insensitive filesystems): ` + inner, Code: 1, Prev: &UserError{Msg: inner, Code: 1}}
		}
	}
	return nil
}

// extractXz is XzDownloader::extract.
func (im *InstallationManager) extractXz(file, path string) error {
	command := []string{"tar", "-xJf", file, "-C", path}
	proc := im.process()
	code, _, err := proc.Execute(command, "")
	if err == nil && code == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	return &UserError{Msg: "Failed to execute " + strings.Join(command, " ") + "\n\n" + proc.ErrorOutput(), Code: 1}
}

// extractRar is RarDownloader::extract for a PHP runtime without the rar
// extension (gocomposer never uses RarArchive).
func (im *InstallationManager) extractRar(file, path string) error {
	var processError string
	if runtime.GOOS != "windows" {
		command := []string{"sh", "-c", `unrar x -- "$0" "$1" >/dev/null && chmod -R u+w "$1"`, file, path}
		proc := im.process()
		code, _, err := proc.Execute(command, "")
		if err == nil && code == 0 {
			return nil
		}
		processError = "Failed to execute " + strings.Join(command, " ") + "\n\n" + proc.ErrorOutput()
	}
	_, facts, _ := im.c.PHP()
	if (PlatformInfo{facts}).ExtensionLoaded("rar") {
		return &UserError{Msg: "Could not decompress the archive: gocomposer does not use the PHP rar extension (RarArchive), install unrar to extract rar archives.", Code: 1}
	}
	if runtime.GOOS == "windows" {
		return &UserError{Msg: "Could not decompress the archive, enable the PHP rar extension or install unrar.\n" + im.iniMessage() + "\n" + processError, Code: 1}
	}
	return &UserError{Msg: "Could not decompress the archive, enable the PHP rar extension.\n" + im.iniMessage(), Code: 1}
}
