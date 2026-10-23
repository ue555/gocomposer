package cio

// ConsoleIO::sanitize (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	reSanitizeNL   = regexp.MustCompile("\x1b\\[[\x30-\x3F]*[\x20-\x2F]*[\x40-\x7E]|\x1b\\].*?(?:\x1b\\\\|\x07)|\x1b.|[\x01-\x09\x0B\x0C\x0E-\x1A]|\r")
	reSanitizeNoNL = regexp.MustCompile("\x1b\\[[\x30-\x3F]*[\x20-\x2F]*[\x40-\x7E]|\x1b\\].*?(?:\x1b\\\\|\x07)|\x1b.|[\x01-\x1A]")
)

// Sanitize is ConsoleIO::sanitize($message): ANSI escape sequences and
// control characters are removed (\n and \r\n are kept) and invalid UTF-8
// becomes '?'.
func Sanitize(s string) string { return sanitize(s, true) }

func sanitize(s string, allowNewlines bool) string {
	if !utf8.ValidString(s) {
		s = strings.ToValidUTF8(s, "?")
	}
	if !allowNewlines {
		return reSanitizeNoNL.ReplaceAllString(s, "")
	}
	// "\r(?!\n)": sanitize around each \r\n
	parts := strings.Split(s, "\r\n")
	for i, p := range parts {
		parts[i] = reSanitizeNL.ReplaceAllString(p, "")
	}
	return strings.Join(parts, "\r\n")
}
