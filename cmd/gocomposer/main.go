// Command gocomposer is a Composer-compatible dependency manager for PHP
// written in Go.
package main

import (
	"os"

	"gocomposer/internal/cli"
)

func main() {
	os.Exit(cli.NewApp().Run(os.Args[1:]))
}
