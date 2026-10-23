// Unix helpers for the Filesystem::safeFilePutContents / is_writable port
// (src/Composer/Util/Filesystem.php, composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

//go:build unix

package configsource

import (
	"os"
	"syscall"
)

// isWritable is PHP is_writable() (access(2) with W_OK).
func isWritable(path string) bool { return syscall.Access(path, 0x2) == nil }

// keepOwner gives the temporary file the owner and group of the target.
func keepOwner(fh *os.File, tempPath string, target os.FileInfo) bool {
	ts, ok := target.Sys().(*syscall.Stat_t)
	if !ok {
		return true
	}
	st, err := fh.Stat()
	if err != nil {
		return false
	}
	own, ok := st.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	if ts.Uid != own.Uid && os.Chown(tempPath, int(ts.Uid), -1) != nil {
		return false
	}
	if ts.Gid != own.Gid && os.Chown(tempPath, -1, int(ts.Gid)) != nil {
		return false
	}
	return true
}
