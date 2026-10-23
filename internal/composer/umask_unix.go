//go:build !windows

package composer

import (
	"os"
	"sync"
	"syscall"
)

var umaskOnce sync.Once
var umaskVal os.FileMode

func currentUmask() os.FileMode {
	umaskOnce.Do(func() {
		m := syscall.Umask(0)
		syscall.Umask(m)
		umaskVal = os.FileMode(m)
	})
	return umaskVal
}
