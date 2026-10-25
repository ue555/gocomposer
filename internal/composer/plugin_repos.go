package composer

// Repositories implemented by plugin code (a class registered with
// RepositoryManager::setRepositoryClass, a RepositoryInterface of the
// plugin, an ArrayRepository it fills, or one of gocomposer's repositories
// whose packages it changed with addPackage/removePackage): added to the
// repository list (RepositoryManager::addRepository/prependRepository) or
// given to Installer::setAdditionalFixedRepository, they take part in the
// resolution like in Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/RepositoryManager.php, MIT license, (c) Nils
// Adermann, Jordi Boggiano): gocomposer asks the PHP object, which may call
// back into gocomposer for what it inherits.

import (
	"fmt"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// phpRepository is a repository object of the plugin runtime.
type phpRepository struct {
	pm   *PluginManager
	id   int64
	name string
}

var _ repo.Repository = (*phpRepository)(nil)

// phpRepo returns the repository for the runtime object id (one per id).
func (pm *PluginManager) phpRepo(id int64, name string) *phpRepository {
	if pm.phpRepos == nil {
		pm.phpRepos = map[int64]*phpRepository{}
	}
	if r, ok := pm.phpRepos[id]; ok {
		if name != "" {
			r.name = name
		}
		return r
	}
	r := &phpRepository{pm: pm, id: id, name: name}
	pm.phpRepos[id] = r
	return r
}

func (r *phpRepository) call(m string, a *phpjson.Map) (phpjson.Value, error) {
	if !r.pm.active() {
		return nil, &UserError{Msg: "The repository " + r.name + " is implemented by a Composer plugin, which is not running anymore", Code: 1}
	}
	if a == nil {
		a = phpjson.NewMap()
	}
	a.Set("id", r.id)
	a.Set("m", m)
	a.Set("frames", r.pm.frames())
	return r.pm.host.call("phpRepoOp", a)
}

func (r *phpRepository) packagesOf(v phpjson.Value) ([]*pkg.Package, error) {
	var out []*pkg.Package
	var perr error
	phpjson.Each(v, func(_ string, x phpjson.Value) {
		if perr != nil {
			return
		}
		p, err := r.pm.host.pkgFromRef(x)
		if err != nil {
			perr = err
			return
		}
		if p != nil {
			out = append(out, p)
		}
	})
	return out, perr
}

func (r *phpRepository) RepoName() string { return r.name }

func (r *phpRepository) Kind() repo.Kind { return repo.KindArray }

func (r *phpRepository) Packages() ([]*pkg.Package, error) {
	v, err := r.call("packages", nil)
	if err != nil {
		return nil, err
	}
	return r.packagesOf(v)
}

func constraintArg(c semver.Constraint) phpjson.Value {
	if c == nil {
		return nil
	}
	if _, ok := c.(*semver.MatchAll); ok {
		return "*"
	}
	if s := c.PrettyString(); s != "" {
		return s
	}
	return c.String()
}

func (r *phpRepository) FindPackages(name string, c semver.Constraint) ([]*pkg.Package, error) {
	a := phpjson.NewMap()
	a.Set("name", name)
	a.Set("c", constraintArg(c))
	v, err := r.call("find", a)
	if err != nil {
		return nil, err
	}
	return r.packagesOf(v)
}

func (r *phpRepository) LoadPackages(names []repo.NameConstraint, acceptable map[string]int, flags map[string]int, alreadyLoaded map[string]map[string]bool) (repo.LoadResult, error) {
	a := phpjson.NewMap()
	list := []phpjson.Value{}
	for _, n := range names {
		m := phpjson.NewMap()
		m.Set("name", n.Name)
		m.Set("c", constraintArg(n.Constraint))
		list = append(list, m)
	}
	a.Set("names", list)
	intMap := func(in map[string]int) *phpjson.Map {
		m := phpjson.NewMap()
		for k, v := range in {
			m.Set(k, int64(v))
		}
		return m
	}
	a.Set("acceptable", intMap(acceptable))
	a.Set("flags", intMap(flags))
	loaded := phpjson.NewMap()
	for name, vers := range alreadyLoaded {
		vm := phpjson.NewMap()
		for v := range vers {
			vm.Set(v, true)
		}
		loaded.Set(name, vm)
	}
	a.Set("loaded", loaded)
	v, err := r.call("load", a)
	if err != nil {
		return repo.LoadResult{}, err
	}
	res := phpjson.ToMap(v)
	if res == nil {
		return repo.LoadResult{}, fmt.Errorf("invalid loadPackages() result of %s", r.name)
	}
	var out repo.LoadResult
	phpjson.Each(mustGet(res, "namesFound"), func(_ string, x phpjson.Value) {
		if s, ok := x.(string); ok {
			out.NamesFound = append(out.NamesFound, strings.ToLower(s))
		}
	})
	out.Packages, err = r.packagesOf(mustGet(res, "packages"))
	return out, err
}

func (r *phpRepository) Providers(name string) ([]repo.ProviderInfo, error) {
	a := phpjson.NewMap()
	a.Set("name", name)
	v, err := r.call("providers", a)
	if err != nil {
		return nil, err
	}
	var out []repo.ProviderInfo
	phpjson.Each(v, func(k string, x phpjson.Value) {
		m := phpjson.ToMap(x)
		if m == nil {
			return
		}
		pi := repo.ProviderInfo{Name: k, Type: m.Str("type")}
		if n := m.Str("name"); n != "" {
			pi.Name = n
		}
		pi.Description, _ = mustGet(m, "description").(string)
		out = append(out, pi)
	})
	return out, nil
}

// reposCallback serves the repository callbacks of services.php (ok is
// false for other ops).
func (pm *PluginManager) reposCallback(op string, a *phpjson.Map) (phpjson.Value, bool, error) {
	switch op {
	case "lockTx":
		v, err := pm.lockTxValue(a)
		return v, true, err
	case "repoOverlay":
		// a plugin changed the packages of one of gocomposer's
		// repositories: the PHP object (keeping them) replaces it
		id, _ := mustGet(a, "rid").(int64)
		base := pm.repoByID[id]
		if base == nil {
			return nil, true, &UserError{Msg: fmt.Sprintf("Unknown gocomposer repository %d", id), Code: 1}
		}
		php, _ := mustGet(a, "php").(int64)
		r := pm.phpRepo(php, base.RepoName())
		for _, c := range append([]*Composer{pm.c}, pm.instanceList()...) {
			for i, e := range c.Repositories {
				if e == base {
					c.Repositories[i] = r
				}
			}
		}
		return nil, true, nil
	}
	return nil, false, nil
}

func (pm *PluginManager) instanceList() []*Composer {
	var out []*Composer
	for _, c := range pm.instances {
		out = append(out, c)
	}
	return out
}

// repoFromSpec is the repository a runtime repository description stands
// for ({rid}: one of gocomposer's, {php, name}: an object of plugin code).
func (pm *PluginManager) repoFromSpec(a *phpjson.Map) (repo.Repository, error) {
	if id, ok := mustGet(a, "rid").(int64); ok {
		r := pm.repoByID[id]
		if r == nil {
			return nil, &UserError{Msg: fmt.Sprintf("Unknown gocomposer repository %d", id), Code: 1}
		}
		return r, nil
	}
	if id, ok := mustGet(a, "php").(int64); ok {
		return pm.phpRepo(id, a.Str("name")), nil
	}
	return nil, nil
}

// lockTxValue is the data of a lock transaction for the runtime's
// LockTransaction (nil: the installer did not update).
func (pm *PluginManager) lockTxValue(a *phpjson.Map) (phpjson.Value, error) {
	var lt *solver.LockTransaction
	if base, _ := mustGet(a, "base").(int64); base != 0 {
		installerMu.Lock()
		ri := installerByID[int(base)]
		installerMu.Unlock()
		if ri != nil {
			lt = ri.in.lockTransaction
		}
	}
	if obj, ok := mustGet(a, "obj").(int64); ok && pm.lockTxs[obj] != nil {
		lt = pm.lockTxs[obj]
	}
	if lt == nil {
		return nil, nil
	}
	h := pm.host
	refs := func(ps []*pkg.Package) []phpjson.Value {
		out := []phpjson.Value{}
		for _, p := range ps {
			out = append(out, h.ref(p))
		}
		return out
	}
	m := phpjson.NewMap()
	ops := []phpjson.Value{}
	for _, op := range lt.Operations() {
		ops = append(ops, opSpec(h, op))
	}
	m.Set("ops", ops)
	m.Set("present", refs(lt.PresentPackages()))
	m.Set("result", refs(lt.ResultPackages()))
	m.Set("nonDev", refs(lt.NewLockPackagesMirrors(false, false)))
	m.Set("dev", refs(lt.NewLockPackagesMirrors(true, false)))
	m.Set("nonDevMirrors", refs(lt.NewLockPackagesMirrors(false, true)))
	m.Set("devMirrors", refs(lt.NewLockPackagesMirrors(true, true)))
	return m, nil
}
