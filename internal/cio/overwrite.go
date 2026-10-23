package cio

// ConsoleIO::overwriteError and the $lastMessageErr bookkeeping of
// ConsoleIO::doWrite, ported from composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/IO/ConsoleIO.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import "strings"

// noteErr records the last message written to stderr (ConsoleIO's
// $lastMessageErr) when it passed the verbosity filter.
func (o *IO) noteErr(msg string, verbosity int) {
	if verbosity > o.Verbosity || (o.Verbosity <= Quiet && verbosity != Quiet) {
		return
	}
	o.mu.Lock()
	o.lastErr = msg
	o.mu.Unlock()
}

// OverwriteError is ConsoleIO::overwriteError($message, $newline) at normal
// verbosity: on a decorated output the previous stderr message is erased
// with backspaces, otherwise a plain line is written.
func (o *IO) OverwriteError(msg string, newline bool) {
	if !o.Decorated {
		if msg != "" {
			o.WriteError(msg)
		}
		o.noteErr(msg, Normal)
		return
	}
	o.mu.Lock()
	size := len(StripTags(o.lastErr))
	o.mu.Unlock()
	o.WriteErrorV(strings.Repeat("\x08", size), false, Normal)
	o.WriteErrorV(msg, false, Normal)
	if fill := size - len(StripTags(msg)); fill > 0 {
		o.WriteErrorV(strings.Repeat(" ", fill), false, Normal)
		o.WriteErrorV(strings.Repeat("\x08", fill), false, Normal)
	}
	if newline {
		o.WriteErrorV("", true, Normal)
	}
	o.noteErr(msg, Normal)
}
