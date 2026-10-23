//go:build darwin || freebsd || netbsd || openbsd

package cio

import "syscall"

const (
	ioctlGetTermios = syscall.TIOCGETA
	ioctlSetTermios = syscall.TIOCSETA
)
