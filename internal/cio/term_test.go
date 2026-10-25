package cio

import (
	"os"
	"strings"
	"testing"
)

// Symfony OutputFormatter::createStyleFromString and
// OutputFormatterStyle::apply (symfony/console v5.4) details.
func TestFormatInlineStyleCase(t *testing.T) {
	o := &IO{Decorated: true}
	cases := map[string]string{
		// the href value keeps its case, keys and colors do not matter
		"<href=https://x.test/CVERecord?id=CVE-1>c</>": "\x1b]8;;https://x.test/CVERecord?id=CVE-1\x1b\\c\x1b]8;;\x1b\\",
		"<FG=RED>x</>":           "\x1b[31mx\x1b[39m",
		"<fg=red;bogus=1>x</>":   "<fg=red;bogus=1>x",
		"<Info>x</Info>":         "<Info>x</Info>",
		`<href=a\<b\>>x</>`:      "\x1b]8;;a<b>\x1b\\x\x1b]8;;\x1b\\",
		"<fg=yellow;href=u>k</>": "\x1b[33m\x1b]8;;u\x1b\\k\x1b]8;;\x1b\\\x1b[39m",
	}
	for in, want := range cases {
		if got := o.Format(in); got != want {
			t.Errorf("Format(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHrefTerminals(t *testing.T) {
	for _, c := range []struct {
		env  map[string]string
		want bool
	}{
		{map[string]string{}, true},
		{map[string]string{"TERMINAL_EMULATOR": "JetBrains-JediTerm"}, false},
		{map[string]string{"KONSOLE_VERSION": "201100"}, false},
		{map[string]string{"KONSOLE_VERSION": "201200"}, true},
		{map[string]string{"KONSOLE_VERSION": "0"}, true},
		{map[string]string{"KONSOLE_VERSION": "abc"}, false},
		{map[string]string{"IDEA_INITIAL_DIRECTORY": ""}, false},
	} {
		for _, k := range []string{"TERMINAL_EMULATOR", "KONSOLE_VERSION", "IDEA_INITIAL_DIRECTORY"} {
			if v, ok := c.env[k]; ok {
				t.Setenv(k, v)
			} else {
				t.Setenv(k, "")
				unsetenv(t, k)
			}
		}
		if got := handlesHrefGracefully(); got != c.want {
			t.Errorf("%v: got %v", c.env, got)
		}
	}
}

func TestTerminalWidthColumns(t *testing.T) {
	for in, want := range map[string]int{"100": 100, " 42 ": 42, "abc": 0, "": 0, "77x": 77} {
		t.Setenv("COLUMNS", in)
		if got := TerminalWidth(); got != want {
			t.Errorf("COLUMNS=%q: %d, want %d", in, got, want)
		}
	}
}

func TestProgressBarNarrowTerminal(t *testing.T) {
	t.Setenv("COLUMNS", "5")
	var sb stringsBuilder
	o := &IO{Err: &sb, Decorated: true, Verbosity: Normal}
	b := o.NewProgressBar(0)
	b.Start(6)
	b.SetProgress(6)
	want := " 0/6 [>]   0%\x1b[1G\x1b[2K 6/6 [=] 100%"
	if sb.String() != want {
		t.Fatalf("%q, want %q", sb.String(), want)
	}
}

type stringsBuilder = strings.Builder

func unsetenv(t *testing.T, k string) {
	t.Helper()
	os.Unsetenv(k)
}
