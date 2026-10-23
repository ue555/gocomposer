//go:build !(darwin || linux || freebsd || netbsd || openbsd)

package cio

func ttyWidth() int { return 0 }
