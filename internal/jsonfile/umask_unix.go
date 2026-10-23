//go:build !windows

package jsonfile

import (
	"os"
	"syscall"
)

func currentUmask() int {
	m := syscall.Umask(0)
	syscall.Umask(m)
	return m
}

// sameOwner reports whether the temp file has the target's owner and group
// (or could be given them), as replaceFile requires.
func sameOwner(target os.FileInfo, tmp *os.File) bool {
	ts, ok := target.Sys().(*syscall.Stat_t)
	if !ok {
		return true
	}
	fi, err := tmp.Stat()
	if err != nil {
		return false
	}
	ns, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return true
	}
	if ts.Uid != ns.Uid || ts.Gid != ns.Gid {
		return tmp.Chown(int(ts.Uid), int(ts.Gid)) == nil
	}
	return true
}
