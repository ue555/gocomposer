package composer

// Support for Composer\Command\CreateProjectCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
// Jordi Boggiano): Factory::createComposer given the configuration as an
// array, RepositoryFactory::createRepo/defaultRepos and the ProjectInstaller.

import (
	"os"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/installed"
	"gocomposer/internal/locker"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/solver"
)

// NewProjectComposer is Factory::createComposer($io, $config->all()) as
// used by CreateProjectCommand::installRootPackage: the configuration comes
// from data (no composer.json is read) and paths resolve against baseDir.
func NewProjectComposer(io *cio.IO, data *phpjson.Map, baseDir string, opts Options) (*Composer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg, err := config.CreateConfig(cwd, func(m string) { io.WriteError(m) })
	if err != nil {
		return nil, err
	}
	cfg.Merge(data)
	cfg.SetBaseDir(baseDir)
	if err := cfg.LoadAuthEnv(); err != nil {
		return nil, err
	}
	dl, err := httpdl.New(cfg)
	if err != nil {
		return nil, err
	}
	dl.Warn = func(s string) { io.WriteError(s) }
	dl.IO = io
	c := &Composer{IO: io, Config: cfg, Cwd: cwd, JSON: data, HTTP: dl,
		DisablePlugins: opts.DisablePlugins, DisableScripts: opts.DisableScripts}
	root, err := LoadRootPackage(data, cwd, cfg, io)
	if err != nil {
		return nil, err
	}
	c.Package = root
	c.LocalRepo = installed.NewInMemory()
	l, err := locker.New(os.DevNull, []byte("{}"))
	if err != nil {
		return nil, err
	}
	c.Locker = l
	return c, nil
}

// CreateRepositoryFromConfig is RepositoryFactory::createRepo for one
// repository configuration (validated like RepositoryFactory::createRepos).
func (c *Composer) CreateRepositoryFromConfig(rc phpjson.Value) (repo.Repository, error) {
	m := phpjson.NewMap()
	m.Set("0", rc)
	repos, err := c.createRepositories(m)
	if err != nil {
		return nil, err
	}
	return repos[0], nil
}

// DefaultRepositories is RepositoryFactory::defaultRepos (the repositories
// of the configuration, packagist.org unless disabled).
func (c *Composer) DefaultRepositories() ([]repo.Repository, error) {
	return c.createRepositories(c.Config.Repositories())
}

// InstallProject is InstallationManager::execute of the install operation
// through the ProjectInstaller: the package is downloaded and installed into
// dir (no binaries, no installed.json). It returns the installed package
// (with its installation source set).
func (im *InstallationManager) InstallProject(p *pkg.Package, dir string, devMode bool) (*pkg.Package, error) {
	p = p.Clone()
	p.SetBinaries(nil)
	im.InstallPathOverride = func(*pkg.Package) string { return dir }
	op := &solver.Operation{Type: solver.OpInstall, Package: p}
	if err := im.Execute(installed.NewInMemory(), []*solver.Operation{op}, devMode, true, false); err != nil {
		return nil, err
	}
	return p, nil
}
