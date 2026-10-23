package main

import "testing"

func TestIndeterminateRemovalBar(t *testing.T) {
	oracle := "  - Removing acme/http (1.2.0)\n  - Removing acme/log (1.1.0)\n 0/2 [>---------------------------]   0%\n 2/2 [============================] 100%\n  - Downloading acme/log (1.1.0)\n 0/2 [>---------------------------]   0%\n 2/2 [============================] 100%\n"
	want := "  - Removing acme/http (1.2.0)\n  - Removing acme/log (1.1.0)\n    0 [>---------------------------]    0 [->--------------------------]\n  - Downloading acme/log (1.1.0)\n 0/2 [>---------------------------]   0%\n 2/2 [============================] 100%\n"
	if got := indeterminateRemovalBar(oracle); got != want {
		t.Fatalf("got %q", got)
	}
	if got := indeterminateRemovalBar(want); got != want {
		t.Fatalf("indeterminate form changed: %q", got)
	}
}
