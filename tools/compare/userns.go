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

// podmanRunRoot is podmanRun for a container process running as uid 0 (the
// root-* scenarios). Rootless podman on Linux maps the container's root to
// the host user by default, so files root writes into bind-mounted work
// directories stay owned (and removable) by the host user; keep-id is left
// out because it would map uid 0 to a subordinate id instead. CAP_SYS_ADMIN
// lets the scenario unmount podman's /run/.containerenv marker (inside the
// container's own mount namespace) so Platform::isDocker() is false.
func podmanRunRoot(args ...string) []string {
	return append([]string{"run", "--rm", "--user", "0", "--cap-add", "SYS_ADMIN"}, args...)
}
