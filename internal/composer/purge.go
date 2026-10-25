// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Factory.php (purgePackages),
// src/Composer/Installer/InstallationManager.php (isPackageInstalled) and
// src/Composer/Installer/LibraryInstaller.php (isInstalled)
// (MIT, (c) Nils Adermann, Jordi Boggiano).

package composer

import (
	"os"
	"path/filepath"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// purgePackages drops packages whose install path was deleted on the
// filesystem from the local repository (Factory::purgePackages, run once
// everything is initialized on a full load).
func (c *Composer) purgePackages() {
	pkgs, err := c.LocalRepo.Packages()
	if err != nil || len(pkgs) == 0 {
		return
	}
	im := NewInstallationManager(c)
	for _, p := range append([]*pkg.Package(nil), pkgs...) {
		if !im.isPackageInstalled(c, p) {
			c.LocalRepo.RemovePackage(p)
		}
	}
}

// isPackageInstalled is InstallationManager::isPackageInstalled with the
// LibraryInstaller / MetapackageInstaller isInstalled checks.
func (im *InstallationManager) isPackageInstalled(c *Composer, p *pkg.Package) bool {
	if p.AliasOf != nil {
		return c.LocalRepo.HasPackage(p) && im.isPackageInstalled(c, p.AliasOf)
	}
	if !c.LocalRepo.HasPackage(p) {
		return false
	}
	// installers added by plugins answer themselves (plugin_install.go)
	if pm := pluginManagerOf(c); pm.active() {
		info, err := pm.installerFor(p.Type())
		if err == nil && info != nil {
			args := phpjson.NewMap()
			args.Set("id", info.ID)
			args.Set("p", pm.host.ref(p))
			v, err := pm.host.call("isInstalled", args)
			return err != nil || v == true
		}
	}
	if p.Type() == "metapackage" {
		return true
	}
	path := im.InstallPath(p)
	// Filesystem::isReadable
	if _, err := os.Stat(path); err == nil {
		return true
	}
	if fi, err := os.Lstat(path); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		_, err := filepath.EvalSymlinks(path)
		return err == nil
	}
	return false
}
