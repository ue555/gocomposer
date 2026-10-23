//go:build !windows

package fsutil

import (
	"os"
	"syscall"
)

func applyUmask(perm os.FileMode) os.FileMode {
	mask := syscall.Umask(0)
	syscall.Umask(mask)
	return perm &^ os.FileMode(mask)
}
