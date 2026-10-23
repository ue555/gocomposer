package composer

// Global plugins: the plugins installed in COMPOSER_HOME (Factory::
// createGlobalComposer and PluginManager::loadInstalledPlugins of Composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT license, (c) Nils Adermann,
// Jordi Boggiano).

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/config"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/installed"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// globalComposer is the PartialComposer of the global (COMPOSER_HOME)
// project: its configuration and local repository.
type globalComposer struct {
	cfg       *config.Config
	repo      *installed.Repository
	vendorDir string
	rules     []allowRule
}

// isGlobal is Composer::isGlobal(): the project is COMPOSER_HOME itself.
func (c *Composer) isGlobal() bool {
	home := c.Config.GetString("home")
	if home == "" || c.ComposerFile == "" {
		return false
	}
	a := fsutil.Realpath(home)
	b := fsutil.Realpath(filepath.Dir(absPath(c.ComposerFile)))
	return a != "" && a == b
}

// loadGlobalComposer is Factory::createGlobalComposer (failures are only
// reported in debug output, like Composer).
func (pm *PluginManager) loadGlobalComposer() *globalComposer {
	c := pm.c
	if c.isGlobal() {
		return nil
	}
	home := c.Config.GetString("home")
	file := home + "/composer.json"
	// Composer writes "Failed to initialize global composer: ..." in debug
	// output; gocomposer's -vvv output does not reproduce Factory's debug log
	debug := func(string) {}
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	data, err := ParseJSONFile(raw, file)
	if err != nil {
		debug(err.Error())
		return nil
	}
	cfg, err := config.CreateConfig(home, nil)
	if err != nil {
		debug(err.Error())
		return nil
	}
	cfg.Merge(data)
	vendorDir := cfg.GetString("vendor-dir")
	repo := installed.New(vendorDir + "/composer/installed.json")
	if err := repo.Load(); err != nil {
		debug(err.Error())
		return nil
	}
	g := &globalComposer{cfg: cfg, repo: repo, vendorDir: vendorDir}
	v := cfg.Get("allow-plugins")
	switch t := v.(type) {
	case bool:
		g.rules = []allowRule{{regexp.MustCompile(`^`), t}}
	default:
		phpjson.Each(v, func(pattern string, allow phpjson.Value) {
			b, _ := allow.(bool)
			g.rules = append(g.rules, allowRule{packageNameToRegexp(pattern), b})
		})
	}
	return g
}

// installPath is the global LibraryInstaller path of a global package.
func (g *globalComposer) installPath(p *pkg.Package) string {
	if p.Type() == "metapackage" {
		return ""
	}
	vendor := strings.TrimRight(g.vendorDir, "/")
	_ = os.MkdirAll(vendor, 0o777)
	if r := fsutil.Realpath(vendor); r != "" {
		vendor = r
	}
	base := vendor + "/" + p.PrettyName
	if td := p.TargetDir(); td != "" {
		return base + "/" + td
	}
	return base
}

func (g *globalComposer) has(p *pkg.Package) bool {
	return g != nil && g.repo.HasPackage(p)
}
