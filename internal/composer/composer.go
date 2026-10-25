// Package composer wires the Composer runtime together: Factory::createComposer,
// downloads, installation, the Installer, scripts and autoload dumping.
package composer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/installed"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/locker"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/platform"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// Composer is the loaded project (Composer\Composer).
type Composer struct {
	IO             *cio.IO
	Config         *config.Config
	Cwd            string
	ComposerFile   string
	RawJSON        []byte
	JSON           *phpjson.Map
	Package        *pkg.Package
	Repositories   []repo.Repository
	LocalRepo      *installed.Repository
	Locker         *locker.Locker
	HTTP           *httpdl.Downloader
	DisablePlugins bool
	DisableScripts bool
	// PolicyResolver maps repositories to their policy (filter list,
	// security advisory) providers.
	PolicyResolver *policy.ProviderResolver
	phpBinary      string
	phpFacts       *platform.Facts
	phpProbed      bool
	phpErr         error
}

// UserError is an error rendered like a Composer exception with a code.
type UserError struct {
	Msg  string
	Code int
	Prev error
}

func (e *UserError) Error() string { return e.Msg }

// Previous returns the chained (PHP previous) exception.
func (e *UserError) Previous() error { return e.Prev }

// userErr converts err to a *UserError keeping its previous-exception chain.
func userErr(err error) *UserError {
	ue := &UserError{Msg: err.Error(), Code: 1}
	if pe, ok := err.(interface{ Previous() error }); ok {
		ue.Prev = pe.Previous()
	}
	return ue
}

// ComposerFileName is Factory::getComposerFile.
func ComposerFileName() (string, error) {
	if env, ok := os.LookupEnv("COMPOSER"); ok {
		env = strings.TrimSpace(env)
		if env != "" {
			if st, err := os.Stat(env); err == nil && st.IsDir() {
				return "", &UserError{Msg: "The COMPOSER environment variable is set to " + env + " which is a directory, this variable should point to a composer.json or be left unset.", Code: 1}
			}
			return env, nil
		}
	}
	return "./composer.json", nil
}

// LockFileName is Factory::getLockFile.
func LockFileName(composerFile string) string {
	if strings.HasSuffix(composerFile, ".json") && filepath.Ext(composerFile) == ".json" {
		return composerFile[:len(composerFile)-4] + "lock"
	}
	return composerFile + ".lock"
}

// Options configure Create.
type Options struct {
	DisablePlugins bool
	DisableScripts bool
	FullLoad       bool
	// ComposerFile overrides Factory::getComposerFile() (Factory::create
	// with an explicit config path, e.g. validate <file>).
	ComposerFile string
	// PlatformFacts replaces probing the PHP runtime (tests).
	PlatformFacts *platform.Facts
	// SkipSchemaValidation mirrors Factory::createComposer being given the
	// configuration as an array (no file, so no LAX schema check).
	SkipSchemaValidation bool
	// SkipPurge leaves the local repository as loaded: the upstream
	// installer fixtures' InstallationManagerMock::isPackageInstalled only
	// checks the repository, so Factory::purgePackages never drops anything
	SkipPurge bool
	// Cwd is Factory::createComposer's $cwd (default: the working directory).
	Cwd string
}

// Create is Factory::createComposer.
func Create(io *cio.IO, opts Options) (*Composer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if opts.Cwd != "" {
		cwd = opts.Cwd
	}
	composerFile := opts.ComposerFile
	if composerFile == "" {
		composerFile, err = ComposerFileName()
		if err != nil {
			return nil, err
		}
	}
	raw, err := os.ReadFile(composerFile)
	if err != nil {
		var msg string
		if composerFile == "./composer.json" || composerFile == "composer.json" {
			msg = "Composer could not find a composer.json file in " + cwd
		} else {
			msg = "Composer could not find the config file: " + composerFile
		}
		return nil, &UserError{Msg: msg + "\nTo initialize a project, please create a composer.json file. See https://getcomposer.org/basic-usage", Code: 1}
	}
	if !filepath.IsAbs(composerFile) && composerFile != "./composer.json" && composerFile != "composer.json" {
		cwd2 := filepath.Dir(composerFile)
		if abs, err := filepath.Abs(cwd2); err == nil {
			_ = abs
		}
	}
	data, err := ParseJSONFile(raw, composerFile)
	if err != nil {
		return nil, err
	}
	if !opts.SkipSchemaValidation {
		if err := ValidateSchemaLax(data, composerFile); err != nil {
			return nil, err
		}
	}
	// -vvv: JsonFile::read and the config files Factory loads (factorylog.go)
	if !opts.SkipSchemaValidation {
		ReadingLog(io, composerFile)
	}
	cfg, err := config.CreateConfigDebug(cwd, func(m string) { io.WriteError(m) }, debugLogger(io))
	if err != nil {
		return nil, err
	}
	cfg.Merge(data)
	if log := debugLogger(io); log != nil && !opts.SkipSchemaValidation {
		log("Loading config file " + composerFile + " (" + phpRealpath(composerFile) + ")")
	}
	// local auth.json next to composer.json
	authFile := filepath.Join(filepath.Dir(absPath(composerFile)), "auth.json")
	if _, err := os.Stat(authFile); err == nil {
		if log := debugLogger(io); log != nil {
			log("Loading config file " + authFile)
		}
		if err := config.ValidateSchemaWarn(func(m string) { io.WriteError(m) }, authFile, nil, jsonschema.AuthSchema); err != nil {
			return nil, err
		}
		ReadingLog(io, authFile)
		m, err := config.ReadJSONFile(authFile)
		if err != nil {
			return nil, err
		}
		wrap := phpjson.NewMap()
		wrap.Set("config", m)
		cfg.Merge(wrap)
	}
	if err := cfg.LoadAuthEnv(); err != nil {
		return nil, err
	}
	dl, err := httpdl.New(cfg)
	if err != nil {
		return nil, err
	}
	tlsDefaultsDebug(io, cfg)
	dl.Warn = func(s string) { io.WriteError(s) }
	// interactive credential prompts and the config sources they write to
	dl.IO = io
	dl.Sources = httpdl.ConfigSources{Config: absPath(composerFile)}
	if _, err := os.Stat(authFile); err == nil {
		dl.Sources.LocalAuth = authFile
	}
	c := &Composer{IO: io, Config: cfg, Cwd: cwd, ComposerFile: composerFile, RawJSON: raw, JSON: data, HTTP: dl,
		DisablePlugins: opts.DisablePlugins, DisableScripts: opts.DisableScripts}
	if opts.PlatformFacts != nil {
		c.phpProbed, c.phpFacts, c.phpBinary = true, opts.PlatformFacts, opts.PlatformFacts.PHPBinary
	}
	repos, err := c.createRepositories(cfg.Repositories())
	if err != nil {
		return nil, err
	}
	c.Repositories = repos
	root, err := LoadRootPackage(data, cwd, cfg, io)
	if err != nil {
		return nil, err
	}
	c.Package = root
	vendorDir := cfg.GetString("vendor-dir")
	c.LocalRepo = installed.New(vendorDir + "/composer/installed.json")
	if opts.FullLoad {
		if err := filesCacheGC(io, cfg); err != nil {
			return nil, err
		}
	}
	// Factory::createDefaultInstallers reads bin-compat for the
	// BinaryInstaller (its deprecation notice is triggered there)
	_, _ = cfg.GetFlags("bin-compat", 0)
	lockFile := LockFileName(composerFile)
	if !cfg.GetBool("lock") {
		if _, err := os.Stat(lockFile); err == nil {
			io.WriteError("<warning>" + lockFile + " is present but ignored as the \"lock\" config option is disabled.</warning>")
		}
		lockFile = os.DevNull
	}
	l, err := locker.New(lockFile, raw)
	if err != nil {
		return nil, err
	}
	c.Locker = l
	if io.IsDebug() {
		l.ReadLog = func(p string) { ReadingLog(io, p) }
	}
	if opts.FullLoad {
		if !c.isGlobal() {
			// Factory::createGlobalComposer (debug output only)
			globalComposerDebug(io, cfg)
		}
		// a custom composer.json path makes its directory Factory's cwd
		factoryCwd := ""
		if opts.ComposerFile != "" && opts.Cwd == "" && isFile(composerFile) {
			factoryCwd = filepath.Dir(composerFile)
		}
		pluginManagerDebug(io, c, factoryCwd)
	}
	if opts.FullLoad {
		if err := c.LocalRepo.Load(); err != nil {
			return nil, err
		}
	}
	// PluginManager::loadInstalledPlugins and the init event (plugin.go)
	if err := c.initPlugins(opts.FullLoad); err != nil {
		return nil, err
	}
	if opts.FullLoad && !opts.SkipPurge {
		c.purgePackages()
	}
	return c, nil
}

func absPath(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	if r, err := filepath.EvalSymlinks(a); err == nil {
		return r
	}
	return a
}

// ParseJSONFile decodes composer.json like JsonFile::read.
func ParseJSONFile(raw []byte, file string) (*phpjson.Map, error) {
	v, err := jsonfile.Parse(raw, file)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		return nil, &UserError{Msg: "\"" + file + "\" does not contain valid JSON", Code: 1}
	}
	return m, nil
}

var reHTTPPrefix = regexp.MustCompile(`(?i)^https?://`)

// createRepositories is RepositoryFactory::defaultRepos + RepositoryManager::createRepository.
func (c *Composer) createRepositories(repos *phpjson.Map) ([]repo.Repository, error) {
	var out []repo.Repository
	names := map[string]bool{}
	for _, index := range repos.Keys() {
		rv, _ := repos.Get(index)
		if s, ok := rv.(string); ok {
			_ = s
			return nil, &UserError{Msg: `"repositories" should be an array of repository definitions, only a single repository was given`, Code: 1}
		}
		rc := phpjson.ToMap(rv)
		if rc == nil {
			enc, _ := phpjson.Encode(rv, 0)
			return nil, &UserError{Msg: fmt.Sprintf(`Repository "%s" (%s) should be an array, %s given`, index, enc, phpType(rv)), Code: 1}
		}
		if !rc.Isset("type") {
			enc, _ := phpjson.Encode(rc, 0)
			return nil, &UserError{Msg: fmt.Sprintf(`Repository "%s" (%s) must have a type defined`, index, enc), Code: 1}
		}
		name := index
		if _, isInt := phpjson.IntKey(index); isInt && rc.Isset("url") {
			name = reHTTPPrefix.ReplaceAllString(rc.Str("url"), "")
		}
		for names[name] {
			name += "2"
		}
		names[name] = true
		r, err := c.createRepository(rc.Str("type"), rc, index)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func phpType(v phpjson.Value) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "bool"
	case int64:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	}
	return "array"
}

func (c *Composer) createRepository(typ string, rc *phpjson.Map, name string) (repo.Repository, error) {
	if p, ok := rc.Get("packagist"); ok && p == false {
		enc, _ := phpjson.Encode(rc, 0)
		c.IO.WriteError(`<warning>Repository "` + name + `" (` + enc + `) has a packagist key which should be in its own repository definition</warning>`)
	}
	var filterConfig *phpjson.Map
	if rc.Isset("only") || rc.Isset("exclude") || rc.Isset("canonical") {
		filterConfig = rc
		rc = rc.Clone()
		rc.Delete("only")
		rc.Delete("exclude")
		rc.Delete("canonical")
	}
	var r repo.Repository
	switch typ {
	case "composer":
		cr, err := repo.NewComposer(rc, c.Config.GetString("cache-repo-dir"), c.Config.GetBool("cache-read-only"), c.HTTP)
		if err != nil {
			return nil, userErr(err)
		}
		cr.Warn = func(s string) { c.IO.WriteError(s) }
		r = cr
	case "package":
		pr, err := repo.NewPackageRepository(rc)
		if err != nil {
			return nil, err
		}
		c.policyResolver().RegisterPackageRepository(pr, rc)
		r = pr
	case "path":
		pr, err := repo.NewPath(rc, c.Cwd)
		if err != nil {
			return nil, userErr(err)
		}
		pr.GuessVersion = func(dir string, cfg *phpjson.Map) (string, string, bool) {
			return GuessVersion(dir, cfg, c.Config)
		}
		r = pr
	case "vcs", "git", "github", "gitlab", "bitbucket", "git-bitbucket", "svn", "hg", "fossil", "perforce":
		vr, err := repo.NewVcs(rc, c.VcsEnv())
		if err != nil {
			return nil, userErr(err)
		}
		r = vr
	case "artifact":
		ar, err := repo.NewArtifact(rc, c.IO)
		if err != nil {
			return nil, userErr(err)
		}
		r = ar
	case "pear":
		return nil, &FactoryInvalidArgumentError{Msg: "The PEAR repository has been removed from Composer 2.x"}
	default:
		return nil, &FactoryInvalidArgumentError{Msg: "Repository type is not registered: " + typ}
	}
	if filterConfig != nil {
		f, err := repo.NewFilter(r, filterConfig)
		if err != nil {
			return nil, userErr(err)
		}
		return f, nil
	}
	return r, nil
}

// PHP returns the PHP binary and runtime facts used for platform checks.
func (c *Composer) PHP() (string, *platform.Facts, error) {
	if !c.phpProbed {
		c.phpProbed = true
		bin, err := platform.FindPHP()
		if err != nil {
			c.phpErr = err
		} else {
			c.phpBinary = bin
			facts, err := platform.Probe(bin)
			if err != nil {
				c.phpErr = err
			} else {
				c.phpFacts = facts
			}
		}
	}
	return c.phpBinary, c.phpFacts, c.phpErr
}

// PlatformOverrides converts a config.platform value to overrides.
func PlatformOverrides(v phpjson.Value) ([]platform.Override, error) {
	var out []platform.Override
	var err error
	phpjson.Each(v, func(name string, val phpjson.Value) {
		if err != nil {
			return
		}
		switch t := val.(type) {
		case string:
			out = append(out, platform.Override{Name: name, Version: t})
		case bool:
			if t {
				err = &UserError{Msg: "config.platform." + name + " should be a string or false, but got bool true", Code: 1}
				return
			}
			out = append(out, platform.Override{Name: name, Disabled: true})
		default:
			err = &UserError{Msg: fmt.Sprintf("config.platform.%s should be a string or false, but got %s", name, phpType(val)), Code: 1}
		}
	})
	return out, err
}

// CreatePlatformRepo is Installer::createPlatformRepo.
func (c *Composer) CreatePlatformRepo(overrides phpjson.Value) (*platform.Repository, error) {
	ov, err := PlatformOverrides(overrides)
	if err != nil {
		return nil, err
	}
	_, facts, _ := c.PHP()
	r, err := platform.New(facts, ov)
	if err != nil {
		return nil, userErr(err)
	}
	return r, nil
}

// PlatformInfo adapts probed facts for problem reporting.
type PlatformInfo struct{ Facts *platform.Facts }

func (p PlatformInfo) PHPVersion() string {
	if p.Facts == nil {
		return ""
	}
	return p.Facts.PHPVersion
}

func (p PlatformInfo) ExtensionLoaded(ext string) bool {
	if p.Facts == nil {
		return false
	}
	for _, e := range p.Facts.Extensions {
		if strings.EqualFold(e, ext) {
			return true
		}
	}
	return false
}

func (p PlatformInfo) ExtensionVersion(ext string) (string, bool) {
	if p.Facts == nil {
		return "", false
	}
	v, ok := p.Facts.ExtVersions[ext]
	return v, ok
}

func (p PlatformInfo) IniPaths() []string {
	if p.Facts == nil {
		return nil
	}
	return p.Facts.IniFiles
}

var _ solver.PlatformInfo = PlatformInfo{}

// VendorDir returns the absolute vendor dir.
func (c *Composer) VendorDir() string { return c.Config.GetString("vendor-dir") }

// RootAliases returns the root package aliases.
func (c *Composer) RootAliases() []pkg.RootAlias { return c.Package.Base().Aliases }

func (c *Composer) policyResolver() *policy.ProviderResolver {
	if c.PolicyResolver == nil {
		c.PolicyResolver = &policy.ProviderResolver{IO: c.IO}
		if c.HTTP != nil {
			c.PolicyResolver.HTTP = c.HTTP
		}
	}
	return c.PolicyResolver
}

// Root returns the non-alias root package.
func (c *Composer) Root() *pkg.Package { return c.Package.Base() }

var _ = semver.Normalize
