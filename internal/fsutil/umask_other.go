//go:build windows

package fsutil

import "os"

func applyUmask(perm os.FileMode) os.FileMode { return perm }
