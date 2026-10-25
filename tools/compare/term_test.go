package main

import "testing"

func TestNormalizeTermTiming(t *testing.T) {
	oracle := ` 0/6 [>---]   0%\e[1G\e[2K 3/6 [==>-]  50%\e[1G\e[2K 6/6 [====] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\`
	gocomp := ` 0/6 [>---]   0%\e]7501;state=working:app=composer:progress=0\e\\e[1G\e[2K 6/6 [====] 100%\e]7501;state=working:app=composer:progress=100\e\\e]7501;state=working:app=composer\e\`
	if a, b := normalizeTermTiming(oracle), normalizeTermTiming(gocomp); a != b {
		t.Fatalf("\n%s\n%s", a, b)
	}
	if got := normalizeTermTiming("Operation timed out after 10002 milliseconds with 0 bytes received"); got != "Operation timed out after N milliseconds with 0 bytes received" {
		t.Fatal(got)
	}
	if got := normalizeTermTiming(` 6/6 [====] 100%  1 sec`); got != ` 6/6 [====] 100% <elapsed>` {
		t.Fatal(got)
	}
}
