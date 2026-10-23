//go:build unix

package cli

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// currentScriptOwner mirrors PHP's get_current_user(): the owner of the
// running program file.
func currentScriptOwner() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	st, err := os.Stat(exe)
	if err != nil {
		return ""
	}
	sys, ok := st.Sys().(*syscall.Stat_t)
	if !ok {
		return ""
	}
	u, err := user.LookupId(strconv.FormatUint(uint64(sys.Uid), 10))
	if err != nil {
		return ""
	}
	return u.Username
}
