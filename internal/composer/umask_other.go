//go:build windows

package composer

import "os"

func currentUmask() os.FileMode { return 0 }
