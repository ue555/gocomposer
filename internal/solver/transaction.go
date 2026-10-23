package solver

import (
	"regexp"
	"sort"
	"strings"

	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/validate"
)

// Operation types (OperationInterface::getOperationType()).
const (
	OpInstall              = "install"
	OpUpdate               = "update"
	OpUninstall            = "uninstall"
	OpMarkAliasInstalled   = "markAliasInstalled"
	OpMarkAliasUninstalled = "markAliasUninstalled"
)

// Operation is one of Composer's Operation\* classes.
type Operation struct {
	Type    string
	Package *pkg.Package // install/uninstall/alias ops
	Initial *pkg.Package // update
	Target  *pkg.Package // update
}

// MainPackage returns the package an operation is about (target for updates).
func (o *Operation) MainPackage() *pkg.Package {
	if o.Type == OpUpdate {
		return o.Target
	}
	return o.Package
}

// Show ports Operation::show($lock).
func (o *Operation) Show(lock bool) string {
	switch o.Type {
	case OpInstall:
		verb := "Installing"
		if lock {
			verb = "Locking"
		}
		return verb + " <info>" + o.Package.PrettyName + "</info> (<comment>" + o.Package.FullPrettyVersion(true) + "</comment>)"
	case OpUninstall:
		return "Removing <info>" + o.Package.PrettyName + "</info> (<comment>" + o.Package.FullPrettyVersion(true) + "</comment>)"
	case OpUpdate:
		return o.updateShow(lock)
	case OpMarkAliasInstalled:
		return "Marking <info>" + o.Package.PrettyName + "</info> (<comment>" + o.Package.FullPrettyVersion(true) + "</comment>) as installed, alias of <info>" + o.Package.AliasOf.PrettyName + "</info> (<comment>" + o.Package.AliasOf.FullPrettyVersion(true) + "</comment>)"
	case OpMarkAliasUninstalled:
		return "Marking <info>" + o.Package.PrettyName + "</info> (<comment>" + o.Package.FullPrettyVersion(true) + "</comment>) as uninstalled, alias of <info>" + o.Package.AliasOf.PrettyName + "</info> (<comment>" + o.Package.AliasOf.FullPrettyVersion(true) + "</comment>)"
	}
	return ""
}

func (o *Operation) updateShow(lock bool) string {
	fromVersion := o.Initial.FullPrettyVersion(true)
	toVersion := o.Target.FullPrettyVersion(true)
	if fromVersion == toVersion && o.Initial.SourceReference() != o.Target.SourceReference() {
		fromVersion = fullPrettyVersionMode(o.Initial, displaySourceRef)
		toVersion = fullPrettyVersionMode(o.Target, displaySourceRef)
	} else if fromVersion == toVersion && o.Initial.DistReference() != o.Target.DistReference() {
		fromVersion = fullPrettyVersionMode(o.Initial, displayDistRef)
		toVersion = fullPrettyVersionMode(o.Target, displayDistRef)
	}
	action := "Downgrading"
	if IsUpgrade(o.Initial.Version(), o.Target.Version()) {
		action = "Upgrading"
	}
	return action + " <info>" + o.Initial.PrettyName + "</info> (<comment>" + fromVersion + "</comment> => <comment>" + toVersion + "</comment>)"
}

const (
	displaySourceRef = 1
	displayDistRef   = 2
)

func fullPrettyVersionMode(p *pkg.Package, mode int) string {
	var ref string
	if mode == displaySourceRef {
		ref = p.SourceReference()
	} else {
		ref = p.DistReference()
	}
	if ref == "" {
		return p.PrettyVersion()
	}
	if len(ref) == 40 && p.SourceType() != "svn" {
		return p.PrettyVersion() + " " + ref[:7]
	}
	return p.PrettyVersion() + " " + ref
}

// IsUpgrade ports Composer\Package\Version\VersionParser::isUpgrade.
func IsUpgrade(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case "dev-master", "dev-trunk", "dev-default":
		from = semver.DefaultBranchAlias
	}
	switch to {
	case "dev-master", "dev-trunk", "dev-default":
		to = semver.DefaultBranchAlias
	}
	if strings.HasPrefix(from, "dev-") || strings.HasPrefix(to, "dev-") {
		return true
	}
	// Semver::sort([$to, $from])[0] === $from
	nf, err1 := semver.Normalize(from)
	nt, err2 := semver.Normalize(to)
	if err1 != nil || err2 != nil {
		return true
	}
	nf = semver.NormalizeDefaultBranch(nf)
	nt = semver.NormalizeDefaultBranch(nt)
	if nf == nt {
		return false
	}
	return !semver.Compare(nt, "<", nf)
}

// Transaction is Composer\DependencyResolver\Transaction.
type Transaction struct {
	operations           []*Operation
	present              []*pkg.Package
	resultPackages       []*pkg.Package
	resultPackagesByName map[string][]*pkg.Package
}

func packageSortLess(a, b *pkg.Package) int {
	if a.Name == b.Name {
		aAlias, bAlias := a.AliasOf != nil, b.AliasOf != nil
		if aAlias != bAlias {
			if aAlias {
				return -1
			}
			return 1
		}
		return strings.Compare(b.Version(), a.Version())
	}
	return strings.Compare(b.Name, a.Name)
}

func NewTransaction(present, result []*pkg.Package) *Transaction {
	t := &Transaction{present: present, resultPackagesByName: map[string][]*pkg.Package{}}
	t.resultPackages = append([]*pkg.Package(nil), result...)
	for _, p := range result {
		for _, n := range p.Names(true) {
			t.resultPackagesByName[n] = append(t.resultPackagesByName[n], p)
		}
	}
	sort.SliceStable(t.resultPackages, func(i, j int) bool { return packageSortLess(t.resultPackages[i], t.resultPackages[j]) < 0 })
	for n, list := range t.resultPackagesByName {
		sort.SliceStable(list, func(i, j int) bool { return packageSortLess(list[i], list[j]) < 0 })
		t.resultPackagesByName[n] = list
	}
	t.operations = t.calculateOperations()
	return t
}

func (t *Transaction) Operations() []*Operation { return t.operations }

func (t *Transaction) providersInResult(l *pkg.Link) []*pkg.Package {
	return t.resultPackagesByName[l.Target]
}

func (t *Transaction) rootPackages() []*pkg.Package {
	removed := map[*pkg.Package]bool{}
	for _, p := range t.resultPackages {
		if removed[p] {
			continue
		}
		for _, l := range p.Requires() {
			for _, r := range t.providersInResult(l) {
				if r != p {
					removed[r] = true
				}
			}
		}
	}
	var roots []*pkg.Package
	for _, p := range t.resultPackages {
		if !removed[p] {
			roots = append(roots, p)
		}
	}
	return roots
}

func (t *Transaction) calculateOperations() []*Operation {
	var ops []*Operation
	presentPackageMap := map[string]*pkg.Package{}
	var removeOrder []string
	removeMap := map[string]*pkg.Package{}
	presentAliasMap := map[string]*pkg.Package{}
	var removeAliasOrder []string
	removeAliasMap := map[string]*pkg.Package{}
	for _, p := range t.present {
		if p.AliasOf != nil {
			k := p.Name + "::" + p.Version()
			presentAliasMap[k] = p
			if _, ok := removeAliasMap[k]; !ok {
				removeAliasOrder = append(removeAliasOrder, k)
			}
			removeAliasMap[k] = p
		} else {
			presentPackageMap[p.Name] = p
			if _, ok := removeMap[p.Name]; !ok {
				removeOrder = append(removeOrder, p.Name)
			}
			removeMap[p.Name] = p
		}
	}
	stack := t.rootPackages()
	visited := map[*pkg.Package]bool{}
	processed := map[*pkg.Package]bool{}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if processed[p] {
			continue
		}
		if !visited[p] {
			visited[p] = true
			stack = append(stack, p)
			if p.AliasOf != nil {
				stack = append(stack, p.AliasOf)
			} else {
				for _, l := range p.Requires() {
					stack = append(stack, t.providersInResult(l)...)
				}
			}
		} else if !processed[p] {
			processed[p] = true
			if p.AliasOf != nil {
				k := p.Name + "::" + p.Version()
				if _, ok := presentAliasMap[k]; ok {
					delete(removeAliasMap, k)
				} else {
					ops = append(ops, &Operation{Type: OpMarkAliasInstalled, Package: p})
				}
			} else {
				if source, ok := presentPackageMap[p.Name]; ok {
					if p.Version() != source.Version() || p.DistReference() != source.DistReference() ||
						p.SourceReference() != source.SourceReference() ||
						p.IsAbandoned() != source.IsAbandoned() || p.ReplacementPackage() != source.ReplacementPackage() {
						ops = append(ops, &Operation{Type: OpUpdate, Initial: source, Target: p})
					}
					delete(removeMap, p.Name)
				} else {
					ops = append(ops, &Operation{Type: OpInstall, Package: p})
					delete(removeMap, p.Name)
				}
			}
		}
	}
	for _, name := range removeOrder {
		if p, ok := removeMap[name]; ok {
			ops = append([]*Operation{{Type: OpUninstall, Package: p}}, ops...)
		}
	}
	for _, k := range removeAliasOrder {
		if p, ok := removeAliasMap[k]; ok {
			ops = append(ops, &Operation{Type: OpMarkAliasUninstalled, Package: p})
		}
	}
	ops = movePluginsToFront(ops)
	return moveUninstallsToFront(ops)
}

func intersects(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}

func nonPlatformRequires(p *pkg.Package) []string {
	var out []string
	for _, l := range p.Requires() {
		if !repo.IsPlatformPackage(l.Target) {
			out = append(out, l.Target)
		}
	}
	return out
}

func movePluginsToFront(ops []*Operation) []*Operation {
	var dlNoDeps, dlWithDeps, noDeps, withDeps []*Operation
	var dlRequires, pluginRequires []string
	removed := map[int]bool{}
	for idx := len(ops) - 1; idx >= 0; idx-- {
		op := ops[idx]
		var p *pkg.Package
		switch op.Type {
		case OpInstall:
			p = op.Package
		case OpUpdate:
			p = op.Target
		default:
			continue
		}
		extra := p.Extra()
		isDl := false
		if p.Type() == "composer-plugin" && extra != nil {
			if v, ok := extra.Get("plugin-modifies-downloads"); ok && v == true {
				isDl = true
			}
		}
		if isDl || intersects(p.Names(true), dlRequires) {
			requires := nonPlatformRequires(p)
			if isDl && len(requires) == 0 {
				dlNoDeps = append([]*Operation{op}, dlNoDeps...)
			} else {
				dlRequires = append(dlRequires, requires...)
				dlWithDeps = append([]*Operation{op}, dlWithDeps...)
			}
			removed[idx] = true
			continue
		}
		isPlugin := p.Type() == "composer-plugin" || p.Type() == "composer-installer"
		if isPlugin || intersects(p.Names(true), pluginRequires) {
			requires := nonPlatformRequires(p)
			if isPlugin && len(requires) == 0 {
				noDeps = append([]*Operation{op}, noDeps...)
			} else {
				pluginRequires = append(pluginRequires, requires...)
				withDeps = append([]*Operation{op}, withDeps...)
			}
			removed[idx] = true
		}
	}
	out := append(append(append(append([]*Operation(nil), dlNoDeps...), dlWithDeps...), noDeps...), withDeps...)
	for i, op := range ops {
		if !removed[i] {
			out = append(out, op)
		}
	}
	return out
}

func moveUninstallsToFront(ops []*Operation) []*Operation {
	var uninst, rest []*Operation
	for _, op := range ops {
		if op.Type == OpUninstall || op.Type == OpMarkAliasUninstalled {
			uninst = append(uninst, op)
		} else {
			rest = append(rest, op)
		}
	}
	return append(uninst, rest...)
}

// LockTransaction is Composer\DependencyResolver\LockTransaction.
type LockTransaction struct {
	*Transaction
	presentList   []*pkg.Package
	unlockableMap map[int]*pkg.Package
	all           []*pkg.Package
	nonDev        []*pkg.Package
	dev           []*pkg.Package
}

func newLockTransaction(pool *Pool, present []*pkg.Package, unlockable map[int]*pkg.Package, d *Decisions) (*LockTransaction, error) {
	lt := &LockTransaction{presentList: present, unlockableMap: unlockable}
	for _, dec := range d.queue {
		if dec.literal > 0 {
			p := pool.LiteralToPackage(dec.literal)
			// setResultPackages validates every result package
			// (names, "-" prefixed urls/references, perforce ports, bin paths)
			if err := validate.ValidatePackage(p); err != nil {
				return nil, err
			}
			lt.all = append(lt.all, p)
			if _, ok := unlockable[p.ID]; !ok {
				lt.nonDev = append(lt.nonDev, p)
			}
		}
	}
	lt.Transaction = NewTransaction(present, lt.all)
	return lt, nil
}

// SetNonDevPackages is LockTransaction::setNonDevPackages.
func (lt *LockTransaction) SetNonDevPackages(extraction *LockTransaction) {
	packages := extraction.NewLockPackages(false)
	lt.dev = lt.nonDev
	lt.nonDev = nil
	devRemoved := make([]bool, len(lt.dev))
	for _, p := range packages {
		for i, rp := range lt.dev {
			if devRemoved[i] {
				continue
			}
			if p.Name == rp.Name {
				lt.nonDev = append(lt.nonDev, rp)
				devRemoved[i] = true
			}
		}
	}
	var dev []*pkg.Package
	for i, p := range lt.dev {
		if !devRemoved[i] {
			dev = append(dev, p)
		}
	}
	lt.dev = dev
}

// NewLockPackages is getNewLockPackages($devMode) without mirror updates.
func (lt *LockTransaction) NewLockPackages(devMode bool) []*pkg.Package {
	return lt.NewLockPackagesMirrors(devMode, false)
}

// NewLockPackagesMirrors is LockTransaction::getNewLockPackages($devMode,
// $updateMirrors): when updating mirrors, packages not present before are
// replaced by the present package of the same version with refreshed
// URLs/mirrors so the locked references stay as they are.
func (lt *LockTransaction) NewLockPackagesMirrors(devMode, updateMirrors bool) []*pkg.Package {
	src := lt.nonDev
	if devMode {
		src = lt.dev
	}
	present := map[*pkg.Package]bool{}
	for _, p := range lt.presentList {
		present[p] = true
	}
	var out []*pkg.Package
	for _, p := range src {
		if p.AliasOf != nil {
			continue
		}
		if updateMirrors && !present[p] {
			p = lt.updateMirrorAndUrls(p)
		}
		out = append(out, p)
	}
	return out
}

var (
	reKnownDistHost = regexp.MustCompile(`(?i)^https?://(?:(?:www\.)?bitbucket\.org|(api\.)?github\.com|(?:www\.)?gitlab\.com)/`)
	reDistRef40     = regexp.MustCompile(`(?i)(/|sha=)[a-f0-9]{40}(/|$)`)
)

// updateMirrorAndUrls is LockTransaction::updateMirrorAndUrls.
func (lt *LockTransaction) updateMirrorAndUrls(p *pkg.Package) *pkg.Package {
	for _, present := range lt.presentList {
		if p.Name != present.Name || p.Version() != present.Version() {
			continue
		}
		if present.SourceReference() == "" {
			continue
		}
		if present.SourceType() != p.SourceType() {
			continue
		}
		if present.AliasOf == nil {
			present.SetSourceURL(p.SourceURL())
			present.SetSourceMirrors(p.SourceMirrors())
		}
		if present.DistType() != p.DistType() {
			return present
		}
		if p.DistURL() != "" && present.DistReference() != "" && reKnownDistHost.MatchString(p.DistURL()) {
			// preg_replace('{(?<=/|sha=)[a-f0-9]{40}(?=/|$)}i', ref, url)
			present.SetDistURL(reDistRef40.ReplaceAllString(p.DistURL(), "${1}"+present.DistReference()+"${2}"))
		}
		present.SetDistMirrors(p.DistMirrors())
		return present
	}
	return p
}

// Aliases is LockTransaction::getAliases.
func (lt *LockTransaction) Aliases(aliases []pkg.RootAlias) []pkg.RootAlias {
	var used []pkg.RootAlias
	remaining := append([]pkg.RootAlias(nil), aliases...)
	taken := make([]bool, len(remaining))
	for _, p := range lt.all {
		if p.AliasOf == nil {
			continue
		}
		for i, a := range remaining {
			if !taken[i] && a.Package == p.Name {
				used = append(used, a)
				taken[i] = true
			}
		}
	}
	sort.SliceStable(used, func(i, j int) bool { return used[i].Package < used[j].Package })
	return used
}

// AllResultPackages returns every decided package (including aliases).
func (lt *LockTransaction) AllResultPackages() []*pkg.Package { return lt.all }

// NewLocalRepoTransaction is LocalRepoTransaction: lock repo vs installed repo.
func NewLocalRepoTransaction(lockedPackages, localPackages []*pkg.Package) *Transaction {
	return NewTransaction(localPackages, lockedPackages)
}
