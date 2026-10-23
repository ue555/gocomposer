//go:build !linux

// Command ptyrun only runs inside the Linux oracle container (see main_linux.go).
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "ptyrun: only supported on linux")
	os.Exit(125)
}
