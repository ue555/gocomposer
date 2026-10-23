package cio

import (
	"os"
	"regexp"
	"strings"
)

var reColorTerm = regexp.MustCompile(`^((screen|xterm|vt100|vt220|putty|rxvt|ansi|cygwin|linux).*)|(.*-256(color)?(-bce)?)$`)

// HasColorSupport is Symfony StreamOutput::hasColorSupport (symfony/console
// v5.4.54, MIT, (c) Fabien Potencier) for a non-Windows stream.
func HasColorSupport(f *os.File) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if !IsTTY(f) {
		m := strings.ToUpper(os.Getenv("MSYSTEM"))
		if m != "MINGW32" && m != "MINGW64" {
			return false
		}
	}
	if os.Getenv("TERM_PROGRAM") == "Hyper" {
		return true
	}
	for _, k := range []string{"COLORTERM", "ANSICON"} {
		if _, ok := os.LookupEnv(k); ok {
			return true
		}
	}
	if os.Getenv("ConEmuANSI") == "ON" {
		return true
	}
	term := os.Getenv("TERM")
	if term == "dumb" {
		return false
	}
	return reColorTerm.MatchString(term)
}
