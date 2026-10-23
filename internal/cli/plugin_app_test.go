package cli

import "testing"

// The command definition handed to plugins (Application::find proxies and
// the bound input of Application::doRun).
func TestDescribeForPlugins(t *testing.T) {
	cmd := &Command{Name: "require", Aliases: []string{"r"}, Description: "Adds packages",
		Args: []ArgDef{{Name: "packages", Array: true, Desc: "pkgs"}},
		Options: []OptDef{
			{Name: "dev", Kind: OptNone, Desc: "dev"},
			{Name: "prefer-install", Kind: OptRequired, Desc: "pref"},
			{Name: "with", Kind: OptArray, Desc: "with"},
			{Name: "bump-after-update", Kind: OptOptional, Default: "dev"},
		}}
	pc := describeForPlugins(cmd)
	if pc.Name != "require" || len(pc.Aliases) != 1 || pc.Aliases[0] != "r" {
		t.Fatalf("name/aliases: %+v", pc)
	}
	if len(pc.Args) != 1 || !pc.Args[0].Array || pc.Args[0].Required {
		t.Fatalf("args: %+v", pc.Args)
	}
	got := map[string]string{}
	for _, o := range pc.Opts {
		switch {
		case o.None:
			got[o.Name] = "none"
		case o.Required:
			got[o.Name] = "required"
		case o.Optional:
			got[o.Name] = "optional:" + o.Default.(string)
		case o.Array:
			got[o.Name] = "array"
		}
	}
	want := map[string]string{"dev": "none", "prefer-install": "required", "with": "array", "bump-after-update": "optional:dev"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: got %q, want %q", k, got[k], v)
		}
	}
}
