//go:build windows

package jsonfile

import "os"

func currentUmask() int { return 0 }

func sameOwner(target os.FileInfo, tmp *os.File) bool { return true }
