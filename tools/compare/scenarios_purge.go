package main

// Factory::purgePackages: packages whose install path was deleted from the
// filesystem are dropped from the local repository when Composer is created
// (full load), so install/show/fund/etc. see them as not installed.
func purgeScenarios() []Scenario {
	req := `{"acme/http": "^1.0", "acme/log-impl": "^1.0", "acme/cli": "^1.0"}`
	json := project(req, `"scripts": {"drop": "rm -rf vendor/acme/cli", "drop-all": "rm -rf vendor/acme"}`)
	files := map[string]string{"src/Kernel.php": appSrc}
	drop := []string{"update --no-progress", "run-script drop"}
	dropAll := []string{"update --no-progress", "run-script drop-all"}
	return []Scenario{
		{Name: "purge-install", JSON: json, Files: files, Setup: drop, Command: "install --no-progress", Probe: true},
		{Name: "purge-install-all", JSON: json, Files: files, Setup: dropAll, Command: "install --no-progress", Probe: true},
		{Name: "purge-show", JSON: json, Files: files, Setup: drop, Command: "show"},
		{Name: "purge-show-name", JSON: json, Files: files, Setup: drop, Command: "show acme/cli"},
		{Name: "purge-update", JSON: json, Files: files, Setup: drop, Command: "update --no-progress", Probe: true},
		{Name: "purge-dump-autoload", JSON: json, Files: files, Setup: drop, Command: "dump-autoload", Probe: true},
		{Name: "purge-depends", JSON: json, Files: files, Setup: drop, Command: "depends acme/log"},
		{Name: "purge-licenses", JSON: json, Files: files, Setup: drop, Command: "licenses"},
		{Name: "purge-status", JSON: json, Files: files, Setup: drop, Command: "status -v"},
		{Name: "purge-remove", JSON: json, Files: files, Setup: drop, Command: "remove acme/cli --no-progress", Probe: true},
	}
}

func init() { extraScenarios = append(extraScenarios, purgeScenarios) }
