//go:build !unix

package cli

// diskFreeSpace reports no value where it is not implemented (like a
// disk_free_space() failure, which skips the check).
func diskFreeSpace(dir string) (uint64, bool) { return 0, false }
