package main

import (
	"os"
	"runtime"
)

// podmanRun returns the leading "podman run --rm" arguments. On Linux with
// rootless podman (GitHub Actions), the host user is mapped to the oracle
// image's "oracle" user (uid/gid 1000) so bind-mounted work directories stay
// writable and git does not reject them as dubiously owned. podman machine
// on macOS shares files without ownership checks and needs no mapping.
func podmanRun(args ...string) []string {
	run := []string{"run", "--rm"}
	if runtime.GOOS == "linux" && os.Getuid() != 0 {
		run = append(run, "--userns=keep-id:uid=1000,gid=1000")
	}
	return append(run, args...)
}
