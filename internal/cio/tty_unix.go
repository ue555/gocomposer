//go:build darwin || linux || freebsd || netbsd || openbsd

package cio

import (
	"os"
	"syscall"
	"unsafe"
)

func ttyWidth() int {
	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		var ws struct{ row, col, x, y uint16 }
		if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws))); e == 0 && ws.col > 0 {
			return int(ws.col)
		}
	}
	return 0
}
