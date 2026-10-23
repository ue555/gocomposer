//go:build darwin || linux || freebsd || netbsd || openbsd

package cio

import (
	"os"
	"syscall"
	"unsafe"
)

// IsTTY is stream_isatty()/posix_isatty(): the file is a terminal (unlike a
// character-device test, /dev/null is not one).
func IsTTY(f *os.File) bool {
	if f == nil {
		return false
	}
	var tio syscall.Termios
	return termiosIoctl(f.Fd(), ioctlGetTermios, &tio) == nil
}

func termiosIoctl(fd uintptr, req uintptr, tio *syscall.Termios) error {
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(unsafe.Pointer(tio))); e != 0 {
		return e
	}
	return nil
}

// sttyMode changes the terminal mode of f like `stty -echo` (hide) or
// `stty -icanon -echo` (raw) and returns the function restoring the saved
// mode (`stty $(stty -g)`); ok is false when the terminal mode cannot be
// read (Terminal::hasSttyAvailable() false).
func sttyMode(f *os.File, raw bool) (restore func(), ok bool) {
	var saved syscall.Termios
	if f == nil || termiosIoctl(f.Fd(), ioctlGetTermios, &saved) != nil {
		return func() {}, false
	}
	tio := saved
	tio.Lflag &^= syscall.ECHO
	if raw {
		tio.Lflag &^= syscall.ICANON
	}
	_ = termiosIoctl(f.Fd(), ioctlSetTermios, &tio)
	return func() { _ = termiosIoctl(f.Fd(), ioctlSetTermios, &saved) }, true
}
