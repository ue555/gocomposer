package autoload

import (
	"strings"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/pkg"
)

// LoaderData renders a PHP file registering a ClassLoader built like
// AutoloadGenerator::createLoader (PSR-0, PSR-4 and classmap; no "files"),
// used to run script callbacks before or independently of vendor/autoload.php.
// The file expects ClassLoader.php next to it.
func LoaderData(root *pkg.Package, packages []*pkg.Package, installPath func(*pkg.Package) string, cwd, vendorDir string, devMode bool) (string, error) {
	g := &gen{basePath: fsutil.NormalizePath(fsutil.Realpath(cwd)), vendorPath: fsutil.NormalizePath(fsutil.Realpath(vendorDir)), devMode: devMode, root: root, cwd: cwd}
	packageMap := []PackageEntry{{root, pathPtr("")}}
	for _, p := range packages {
		if p.AliasOf != nil {
			continue
		}
		var ip *string
		if p.Type() != "metapackage" {
			ip = pathPtr(installPath(p))
		}
		packageMap = append(packageMap, PackageEntry{p, ip})
	}
	a := g.parseAutoloads(packageMap, false)
	var sb strings.Builder
	sb.WriteString("<?php\nif (!class_exists('Composer\\Autoload\\ClassLoader', false)) { require __DIR__ . '/ClassLoader.php'; }\n")
	sb.WriteString("$loader = new \\Composer\\Autoload\\ClassLoader(" + fsutil.VarExport(vendorDir) + ");\n")
	abs := func(p string) string {
		if fsutil.IsAbsolutePath(p) {
			return p
		}
		return g.basePath + "/" + p
	}
	for _, ns := range a.psr0.keys {
		var paths []string
		for _, p := range a.psr0.vals[ns] {
			paths = append(paths, fsutil.VarExport(abs(p)))
		}
		sb.WriteString("$loader->add(" + fsutil.VarExport(ns) + ", array(" + strings.Join(paths, ", ") + "));\n")
	}
	for _, ns := range a.psr4.keys {
		var paths []string
		for _, p := range a.psr4.vals[ns] {
			paths = append(paths, fsutil.VarExport(abs(p)))
		}
		sb.WriteString("$loader->addPsr4(" + fsutil.VarExport(ns) + ", array(" + strings.Join(paths, ", ") + "));\n")
	}
	if len(a.classmap) > 0 {
		cmg := NewGenerator(cwd)
		for _, dir := range a.classmap {
			ex, err := g.buildExclusionRegexWithVendorDir(dir, a.exclude)
			if err != nil {
				return "", err
			}
			_ = cmg.ScanPaths(dir, ex, "classmap", "")
		}
		sb.WriteString("$loader->addClassMap(array(\n")
		for _, c := range cmg.ClassMap.SortedClasses() {
			sb.WriteString("    " + fsutil.VarExport(c) + " => " + fsutil.VarExport(abs(cmg.ClassMap.Map[c])) + ",\n")
		}
		sb.WriteString("));\n")
	}
	sb.WriteString("$loader->register(false);\n")
	return sb.String(), nil
}
