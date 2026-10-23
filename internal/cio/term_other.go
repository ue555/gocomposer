//go:build !(darwin || linux || freebsd || netbsd || openbsd)

package cio

import "os"

// IsTTY reports false where terminal modes are not supported.
func IsTTY(f *os.File) bool { return false }

func sttyMode(f *os.File, raw bool) (func(), bool) { return func() {}, false }
