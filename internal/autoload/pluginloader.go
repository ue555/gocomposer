package autoload

import (
	"gocomposer/internal/fsutil"
	"gocomposer/internal/pkg"
)

// NamespacePaths is one psr-0/psr-4 prefix with its directories.
type NamespacePaths struct {
	Namespace string
	Paths     []string
}

// PluginAutoload is the class loader PluginManager::registerPackage builds
// for a plugin (AutoloadGenerator::parseAutoloads + createLoader of Composer
// f83588392df43d4fd027d84d26de459946ef4e80): paths are absolute, Files are
// (identifier, path) pairs required with composerRequire().
type PluginAutoload struct {
	PSR0, PSR4 []NamespacePaths
	ClassMap   [][2]string
	Files      [][2]string
}

// PluginLoader parses the autoload rules of entries like
// AutoloadGenerator::parseAutoloads($autoloads, $rootPackage) (the first
// entry takes the root position, as array_shift does there) and resolves them
// like createLoader, scanning classmap directories.
func PluginLoader(root *pkg.Package, entries []PackageEntry, cwd, vendorDir string, devMode bool, warn func(string)) *PluginAutoload {
	g := &gen{basePath: fsutil.NormalizePath(fsutil.Realpath(cwd)), vendorPath: fsutil.NormalizePath(fsutil.Realpath(vendorDir)), devMode: devMode, root: root, cwd: cwd}
	a := g.parseAutoloads(entries, false)
	abs := func(p string) string {
		if fsutil.IsAbsolutePath(p) {
			return p
		}
		return g.basePath + "/" + p
	}
	out := &PluginAutoload{}
	for _, ns := range a.psr0.keys {
		var paths []string
		for _, p := range a.psr0.vals[ns] {
			paths = append(paths, abs(p))
		}
		out.PSR0 = append(out.PSR0, NamespacePaths{ns, paths})
	}
	for _, ns := range a.psr4.keys {
		var paths []string
		for _, p := range a.psr4.vals[ns] {
			paths = append(paths, abs(p))
		}
		out.PSR4 = append(out.PSR4, NamespacePaths{ns, paths})
	}
	if len(a.classmap) > 0 {
		cmg := NewGenerator(cwd)
		for _, dir := range a.classmap {
			ex, err := g.buildExclusionRegexWithVendorDir(dir, a.exclude)
			if err != nil {
				if warn != nil {
					warn("<warning>" + err.Error() + "</warning>")
				}
				continue
			}
			if err := cmg.ScanPaths(dir, ex, "classmap", ""); err != nil && warn != nil {
				warn("<warning>" + err.Error() + "</warning>")
			}
		}
		for _, c := range cmg.ClassMap.SortedClasses() {
			out.ClassMap = append(out.ClassMap, [2]string{c, abs(cmg.ClassMap.Map[c])})
		}
	}
	for _, id := range a.files.keys {
		out.Files = append(out.Files, [2]string{id, abs(a.files.vals[id])})
	}
	return out
}
