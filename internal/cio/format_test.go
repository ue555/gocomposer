package cio

import "testing"

// Expected sequences follow Symfony OutputFormatter / Color::apply
// (set codes before every text segment, unset codes after it).
func TestFormatDecorated(t *testing.T) {
	o := &IO{Decorated: true}
	cases := map[string]string{
		"<info>a</info>":                      "\x1b[32ma\x1b[39m",
		"<warning>w</warning>":                "\x1b[30;43mw\x1b[39;49m",
		"<info>a<comment>b</comment>c</info>": "\x1b[32ma\x1b[39m\x1b[33mb\x1b[39m\x1b[32mc\x1b[39m",
		"<href=https://x.test/a>name</>  ":    "\x1b]8;;https://x.test/a\x1b\\name\x1b]8;;\x1b\\  ",
		"<fg=red;options=bold>x</>":           "\x1b[31;1mx\x1b[39;22m",
		"<info></info>plain":                  "plain",
		"\\<info>literal":                     "<info>literal",
		"<unknown>t</unknown>":                "<unknown>t</unknown>",
	}
	for in, want := range cases {
		if got := o.Format(in); got != want {
			t.Errorf("Format(%q) = %q, want %q", in, got, want)
		}
	}
	plain := &IO{}
	if got := plain.Format("<info>a</info> <href=u>b</>"); got != "a b" {
		t.Errorf("undecorated: %q", got)
	}
}
