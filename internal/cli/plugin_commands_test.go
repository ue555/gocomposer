package cli

import (
	"strings"
	"testing"
)

// Symfony Application::findAlternatives / findNamespace as used for plugin
// (namespaced) commands.
func TestFindAlternatives(t *testing.T) {
	got := strings.Join(findAlternatives("plugfix:nope", []string{"plugfix:echo", "install", "update", "pecho"}), ",")
	if got != "plugfix:echo" {
		t.Fatalf("got %q", got)
	}
	if got := findAlternatives("zzz", []string{"install", "update"}); len(got) != 0 {
		t.Fatalf("got %q", got)
	}
}

func TestNamespaceNotFound(t *testing.T) {
	a := &App{byName: map[string]*Command{"install": {Name: "install"}, "update": {Name: "update"}}}
	_, err := a.findCommand("plugfix:hello")
	if err == nil || err.Error() != `There are no commands defined in the "plugfix" namespace.` {
		t.Fatalf("got %v", err)
	}
	a.byName["plugfix:echo"] = &Command{Name: "plugfix:echo"}
	_, err = a.findCommand("plugfix:nope")
	if err == nil || err.Error() != "Command \"plugfix:nope\" is not defined.\n\nDid you mean this?\n    plugfix:echo" {
		t.Fatalf("got %v", err)
	}
}
