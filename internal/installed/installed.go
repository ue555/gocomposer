// Package installed ports Composer\Repository\FilesystemRepository /
// InstalledFilesystemRepository: vendor/composer/installed.json, installed.php
// and the InstalledVersions class copy.
package installed

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/res"
	"gocomposer/internal/semver"
)

// Repository is the installed (local) repository.
type Repository struct {
	*repo.Array
	File            string
	DevPackageNames []string
	DevMode         *bool
	loaded          bool
	loadErr         error
}

// New creates the repository for vendor/composer/installed.json.
func New(file string) *Repository {
	r := &Repository{Array: repo.NewNamedArray(repo.KindInstalled, "installed repo"), File: file}
	r.Array.SetSelf(r)
	return r
}

// NewInMemory creates an InstalledArrayRepository-like repository that is
// never read from or written to disk (Installer::mockLocalRepositories).
func NewInMemory() *Repository {
	r := New(os.DevNull)
	r.loaded = true
	return r
}

// RepoName is getRepoName() for InstalledFilesystemRepository.
func (r *Repository) RepoName() string { return "installed repo" }

// Load reads installed.json (initialize()).
func (r *Repository) Load() error {
	if r.loaded {
		return r.loadErr
	}
	r.loaded = true
	data, err := os.ReadFile(r.File)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		r.loadErr = err
		return err
	}
	v, err := jsonfile.Parse(data, r.File)
	if err != nil {
		cls := `Seld\JsonLint\ParsingException`
		if _, ok := err.(*jsonfile.EncodingError); ok {
			cls = "UnexpectedValueException"
		}
		r.loadErr = fmt.Errorf("Invalid repository data in %s, packages could not be loaded: [%s] %s", r.File, cls, err.Error())
		return r.loadErr
	}
	var packages phpjson.Value = v
	if m := phpjson.AsMap(v); m != nil {
		if p, ok := m.Get("packages"); ok && p != nil {
			packages = p
		}
		if d, ok := m.Get("dev-package-names"); ok {
			r.DevPackageNames = nil
			phpjson.Each(d, func(_ string, e phpjson.Value) {
				s, _ := e.(string)
				r.DevPackageNames = append(r.DevPackageNames, s)
			})
		}
		if d, ok := m.Get("dev"); ok {
			if b, ok := d.(bool); ok {
				r.DevMode = &b
			}
		}
	}
	if !phpjson.IsArray(packages) {
		r.loadErr = fmt.Errorf("Invalid repository data in %s, packages could not be loaded: [UnexpectedValueException] Could not parse package list from the repository", r.File)
		return r.loadErr
	}
	loader := &pkg.Loader{LoadOptions: true}
	var lerr error
	phpjson.Each(packages, func(_ string, pv phpjson.Value) {
		if lerr != nil {
			return
		}
		m := phpjson.ToMap(pv)
		if m == nil {
			return
		}
		p, err := loader.Load(m)
		if err != nil {
			lerr = err
			return
		}
		r.Array.AddPackage(p)
	})
	if lerr != nil {
		r.loadErr = fmt.Errorf("Invalid repository data in %s, packages could not be loaded: %v", r.File, lerr)
	}
	return r.loadErr
}

// Reload re-reads the dev flag from disk (packages stay as in memory).
func (r *Repository) Reload() error {
	data, err := os.ReadFile(r.File)
	if err != nil {
		return nil
	}
	v, err := phpjson.Decode(data)
	if err != nil {
		return nil
	}
	if m := phpjson.AsMap(v); m != nil {
		if d, ok := m.Get("dev"); ok {
			if b, ok := d.(bool); ok {
				r.DevMode = &b
			}
		}
	}
	return nil
}

// Packages loads lazily.
func (r *Repository) Packages() ([]*pkg.Package, error) {
	if err := r.Load(); err != nil {
		return nil, err
	}
	return r.Array.Packages()
}

// MustPackages ignores load errors.
func (r *Repository) MustPackages() []*pkg.Package {
	p, _ := r.Packages()
	return p
}

// IsFresh is InstalledRepositoryInterface::isFresh: no installed.json yet.
func (r *Repository) IsFresh() bool {
	_, err := os.Stat(r.File)
	return os.IsNotExist(err)
}

// CanonicalPackages is CanonicalPackagesTrait::getCanonicalPackages.
func (r *Repository) CanonicalPackages() []*pkg.Package {
	var order []string
	byName := map[string]*pkg.Package{}
	for _, p := range r.MustPackages() {
		cur, ok := byName[p.Name]
		if !ok {
			order = append(order, p.Name)
		}
		if !ok || cur.AliasOf != nil {
			byName[p.Name] = p
		}
	}
	var out []*pkg.Package
	for _, n := range order {
		out = append(out, byName[n].Base())
	}
	return out
}

// HasPackage checks by unique name.
func (r *Repository) HasPackage(p *pkg.Package) bool {
	for _, e := range r.MustPackages() {
		if e.UniqueName() == p.UniqueName() {
			return true
		}
	}
	return false
}

// RemovePackage removes by unique name.
func (r *Repository) RemovePackage(p *pkg.Package) {
	_ = r.Load()
	r.Array.RemovePackage(p)
}

// AddPackage adds a package (and its aliased package).
func (r *Repository) AddPackage(p *pkg.Package) {
	_ = r.Load()
	r.Array.AddPackage(p)
}

// InstallPathFunc returns the absolute install path of a package ("" for none).
type InstallPathFunc func(p *pkg.Package) string

// Write is FilesystemRepository::write (dumpVersions=true behaviour when root != nil).
func (r *Repository) Write(devMode bool, installPath InstallPathFunc, root *pkg.Package, cwd string) error {
	if r.File == os.DevNull {
		return nil
	}
	repoDir := filepath.Dir(r.File)
	if err := os.MkdirAll(repoDir, 0o777); err != nil {
		return err
	}
	repoDir = fsutil.NormalizePath(fsutil.Realpath(repoDir))
	data := phpjson.NewMap()
	var packages []phpjson.Value
	var devNames []string
	installPaths := map[string]*string{}
	for _, p := range r.CanonicalPackages() {
		arr := pkg.Dump(p)
		ip := installPath(p)
		var installPathRel *string
		if ip != "" {
			abs := ip
			if !fsutil.IsAbsolutePath(abs) {
				abs = cwd + "/" + abs
			}
			sp, err := fsutil.FindShortestPath(repoDir, fsutil.NormalizePath(abs), true, false)
			if err != nil {
				return err
			}
			installPathRel = &sp
		}
		installPaths[p.Name] = installPathRel
		if installPathRel != nil {
			arr.Set("install-path", *installPathRel)
		} else {
			arr.Set("install-path", nil)
		}
		packages = append(packages, arr)
		for _, dn := range r.DevPackageNames {
			if dn == p.Name {
				devNames = append(devNames, p.Name)
				break
			}
		}
	}
	sort.Strings(devNames)
	sort.SliceStable(packages, func(i, j int) bool {
		return packages[i].(*phpjson.Map).Str("name") < packages[j].(*phpjson.Map).Str("name")
	})
	if packages == nil {
		packages = []phpjson.Value{}
	}
	data.Set("packages", packages)
	data.Set("dev", devMode)
	dn := make([]phpjson.Value, len(devNames))
	for i, n := range devNames {
		dn[i] = n
	}
	data.Set("dev-package-names", dn)
	enc, err := phpjson.Encode(data, phpjson.ComposerFlags)
	if err != nil {
		return err
	}
	if err := fsutil.FilePutContentsIfModified(r.File, []byte(enc+"\n")); err != nil {
		return err
	}
	b := devMode
	r.DevMode = &b
	if root == nil {
		return nil
	}
	versions, err := r.generateInstalledVersions(installPaths, devMode, repoDir, root, cwd)
	if err != nil {
		return err
	}
	php := "<?php return " + dumpToPHPCode(versions, 0) + ";\n"
	if err := fsutil.FilePutContentsIfModified(repoDir+"/installed.php", []byte(php)); err != nil {
		return err
	}
	return fsutil.FilePutContentsIfModified(repoDir+"/InstalledVersions.php", res.InstalledVersions)
}

// phpArray is an ordered PHP array used for installed.php generation.
type phpArray struct {
	keys []string
	vals map[string]any
}

func newPHPArray() *phpArray { return &phpArray{vals: map[string]any{}} }

func (a *phpArray) set(k string, v any) {
	if _, ok := a.vals[k]; !ok {
		a.keys = append(a.keys, k)
	}
	a.vals[k] = v
}

func (a *phpArray) get(k string) (any, bool) { v, ok := a.vals[k]; return v, ok }

func (a *phpArray) sub(k string) *phpArray {
	if v, ok := a.vals[k]; ok {
		if s, ok := v.(*phpArray); ok {
			return s
		}
	}
	s := newPHPArray()
	a.set(k, s)
	return s
}

func (a *phpArray) appendList(v any) {
	a.set(strconv.Itoa(len(a.keys)), v)
}

func (a *phpArray) ksort() {
	sort.SliceStable(a.keys, func(i, j int) bool { return phpjson.PHPCompareKeys(a.keys[i], a.keys[j]) < 0 })
}

func (r *Repository) dumpInstalledPackage(p *pkg.Package, installPaths map[string]*string, repoDir string, devPackages map[string]bool, cwd string) *phpArray {
	var reference *string
	if src := p.InstallationSource(); src != "" {
		var ref string
		if src == "source" {
			ref = p.SourceReference()
		} else {
			ref = p.DistReference()
		}
		if ref != "" {
			reference = &ref
		}
	}
	if reference == nil {
		ref := p.SourceReference()
		if ref == "" {
			ref = p.DistReference()
		}
		if ref != "" {
			reference = &ref
		}
	}
	var ip *string
	if p.IsRoot {
		to := fsutil.NormalizePath(fsutil.Realpath(cwd))
		s, _ := fsutil.FindShortestPath(repoDir, to, true, false)
		ip = &s
	} else {
		ip = installPaths[p.Name]
	}
	d := newPHPArray()
	d.set("pretty_version", p.PrettyVersion())
	d.set("version", p.Version())
	if reference != nil {
		d.set("reference", *reference)
	} else {
		d.set("reference", nil)
	}
	d.set("type", p.Type())
	if ip != nil {
		d.set("install_path", installPathValue(*ip))
	} else {
		d.set("install_path", nil)
	}
	d.set("aliases", newPHPArray())
	d.set("dev_requirement", devPackages[p.Name])
	return d
}

type installPathValue string

func (r *Repository) generateInstalledVersions(installPaths map[string]*string, devMode bool, repoDir string, rootPkg *pkg.Package, cwd string) (*phpArray, error) {
	devPackages := map[string]bool{}
	for _, n := range r.DevPackageNames {
		devPackages[n] = true
	}
	packages := append([]*pkg.Package(nil), r.MustPackages()...)
	packages = append(packages, rootPkg)
	root := rootPkg
	for root.AliasOf != nil {
		root = root.AliasOf
		packages = append(packages, root)
	}
	versions := newPHPArray()
	rd := r.dumpInstalledPackage(root, installPaths, repoDir, devPackages, cwd)
	rootData := newPHPArray()
	rootData.set("name", root.Name)
	for _, k := range []string{"pretty_version", "version", "reference", "type", "install_path", "aliases"} {
		v, _ := rd.get(k)
		rootData.set(k, v)
	}
	rootData.set("dev", devMode)
	versions.set("root", rootData)
	vv := newPHPArray()
	versions.set("versions", vv)
	for _, p := range packages {
		if p.AliasOf != nil {
			continue
		}
		vv.set(p.Name, r.dumpInstalledPackage(p, installPaths, repoDir, devPackages, cwd))
	}
	for _, p := range packages {
		isDev := devPackages[p.Name]
		for _, kind := range []struct {
			links pkg.Links
			key   string
		}{{p.Replaces(), "replaced"}, {p.Provides(), "provided"}} {
			for _, l := range kind.links {
				if repo.IsPlatformPackage(l.Target) {
					continue
				}
				entry := vv.sub(l.Target)
				if dr, ok := entry.get("dev_requirement"); !ok {
					entry.set("dev_requirement", isDev)
				} else if !isDev {
					_ = dr
					entry.set("dev_requirement", false)
				}
				val := l.PrettyConstraint
				if val == "self.version" {
					val = p.PrettyVersion()
				}
				list := entry.sub(kind.key)
				found := false
				for _, k := range list.keys {
					if list.vals[k] == val {
						found = true
					}
				}
				if !found {
					list.appendList(val)
				}
			}
		}
	}
	for _, p := range packages {
		if p.AliasOf == nil {
			continue
		}
		vv.sub(p.Name).sub("aliases").appendList(p.PrettyVersion())
		if p.IsRoot {
			rootData.sub("aliases").appendList(p.PrettyVersion())
		}
	}
	vv.ksort()
	versions.ksort()
	for _, name := range vv.keys {
		entry, ok := vv.vals[name].(*phpArray)
		if !ok {
			continue
		}
		for _, key := range []string{"aliases", "replaced", "provided"} {
			if lv, ok := entry.vals[key].(*phpArray); ok {
				var items []string
				for _, k := range lv.keys {
					items = append(items, lv.vals[k].(string))
				}
				sort.SliceStable(items, func(i, j int) bool { return semver.NatCompare(items[i], items[j]) < 0 })
				nl := newPHPArray()
				for _, it := range items {
					nl.appendList(it)
				}
				entry.vals[key] = nl
			}
		}
	}
	return versions, nil
}

func dumpToPHPCode(a *phpArray, level int) string {
	var sb strings.Builder
	sb.WriteString("array(\n")
	level++
	for _, k := range a.keys {
		v := a.vals[k]
		sb.WriteString(strings.Repeat("    ", level))
		if _, isInt := phpjson.IntKey(k); isInt {
			sb.WriteString(k + " => ")
		} else {
			sb.WriteString(fsutil.VarExport(k) + " => ")
		}
		switch t := v.(type) {
		case *phpArray:
			if len(t.keys) > 0 {
				sb.WriteString(dumpToPHPCode(t, level))
			} else {
				sb.WriteString("array(),\n")
			}
		case installPathValue:
			s := string(t)
			if fsutil.IsAbsolutePath(s) {
				sb.WriteString(fsutil.VarExport(s) + ",\n")
			} else {
				sb.WriteString("__DIR__ . " + fsutil.VarExport("/"+s) + ",\n")
			}
		case string:
			sb.WriteString(fsutil.VarExport(t) + ",\n")
		case bool:
			if t {
				sb.WriteString("true,\n")
			} else {
				sb.WriteString("false,\n")
			}
		case nil:
			sb.WriteString("null,\n")
		}
	}
	sb.WriteString(strings.Repeat("    ", level-1) + ")")
	if level-1 != 0 {
		sb.WriteString(",\n")
	}
	return sb.String()
}
