package cli

// Port of Composer\Command\CreateProjectCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// In interactive mode --ask asks for the project directory and the existing
// VCS history removal is confirmed (internal/cio QuestionHelper port).

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/configsource"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

func init() { registerCommand(createProjectCommand) }

func createProjectCommand() *Command {
	return &Command{Name: "create-project", Run: runCreateProject}
}

// createProjectPreferredInstall is getPreferredInstallOptions($config,
// $input, $keepVcsRequiresPreferSource).
func createProjectPreferredInstall(preferred phpjson.Value, in *Input, keepVcs bool) (bool, bool, error) {
	preferSource, preferDist := false, false
	switch preferred {
	case "source":
		preferSource = true
	case "dist":
		preferDist = true
	}
	optSource, optDist := in.Bool("prefer-source"), in.Bool("prefer-dist")
	if pi, ok := in.String("prefer-install"); ok && in.present["prefer-install"] {
		if optSource {
			return false, false, &composer.UserError{Msg: "--prefer-source can not be used together with --prefer-install", Code: 1}
		}
		if optDist {
			return false, false, &composer.UserError{Msg: "--prefer-dist can not be used together with --prefer-install", Code: 1}
		}
		switch pi {
		case "dist":
			optDist = true
		case "source":
			optSource = true
		case "auto":
			preferSource, preferDist = false, false
		default:
			return false, false, &composer.UserError{Msg: `--prefer-install accepts one of "dist", "source" or "auto", got ` + pi, Code: 1}
		}
	}
	keep := keepVcs && in.Bool("keep-vcs")
	if optSource || optDist || keep {
		preferSource = optSource || keep
		preferDist = optDist
	}
	return preferSource, preferDist, nil
}

type createProject struct {
	ctx            *Context
	in             *Input
	config         *config.Config
	disablePlugins bool
	disableScripts bool
}

func runCreateProject(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	cwd, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	// Factory::createConfig
	cfg, err := config.CreateConfig(cwd, func(m string) { io.WriteError(m) })
	if err != nil {
		return 1, err
	}
	preferSource, preferDist, err := createProjectPreferredInstall(cfg.Get("preferred-install"), in, true)
	if err != nil {
		return 1, err
	}
	cp := &createProject{ctx: ctx, in: in, config: cfg, disablePlugins: in.Bool("no-plugins"), disableScripts: in.Bool("no-scripts")}
	if in.Bool("dev") {
		io.WriteError(`<warning>You are using the deprecated option "dev". Dev packages are installed by default now.</warning>`)
	}
	if in.Bool("no-custom-installers") {
		io.WriteError(`<warning>You are using the deprecated option "no-custom-installers". Use "no-plugins" instead.</warning>`)
		cp.disablePlugins = true
	}
	if io.Interactive && in.Bool("ask") {
		if _, ok := in.args["package"]; !ok || in.Arg("package") == "" {
			return 1, &composer.UserError{Msg: `Not enough arguments (missing: "package").`, Code: 1}
		}
		parts := strings.SplitN(strings.ToLower(in.Arg("package")), "/", 2)
		dir, ok, err := io.AskString("New project directory [<comment>"+parts[len(parts)-1]+"</comment>]: ", nil)
		if err != nil {
			return 1, err
		}
		if ok {
			in.args["directory"] = []string{dir}
		} else {
			delete(in.args, "directory")
		}
	}
	var repositories []string
	hasRepositories := false
	if r := in.Strings("repository"); len(r) > 0 {
		repositories, hasRepositories = r, true
	} else if u, ok := in.String("repository-url"); ok && in.present["repository-url"] {
		repositories, hasRepositories = []string{u}, true
	}
	return cp.installProject(in.Arg("package"), in.Arg("directory"), in.Arg("version"), preferSource, preferDist, !in.Bool("no-dev"), repositories, hasRepositories)
}

var reStabilitySuffix = regexp.MustCompile(`(?i)^[^,\s]*?@(stable|RC|beta|alpha|dev)$`)

func (cp *createProject) installProject(packageName, directory, packageVersion string, preferSource, preferDist, installDev bool, repositories []string, hasRepositories bool) (int, error) {
	io, in := cp.ctx.IO, cp.in
	oldCwd, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	defer os.Chdir(oldCwd)
	filter := platformFilter(in)
	var rootPkg *pkg.Package
	installedFromVcs := false
	if packageName != "" {
		rootPkg, err = cp.installRootPackage(packageName, directory, packageVersion, preferSource, preferDist, installDev, repositories, hasRepositories)
		if err != nil {
			return 1, err
		}
		installedFromVcs = rootPkg.InstallationSource() == "source"
	}
	opts := composer.Options{DisablePlugins: cp.disablePlugins, DisableScripts: cp.disableScripts, FullLoad: true}
	c, err := composer.Create(io, opts)
	if err != nil {
		return 1, err
	}
	flushLock := false
	if hasRepositories && in.Bool("add-repository") {
		for index, r := range repositories {
			rc, err := repoConfigFromString(r, true)
			if err != nil {
				return 1, err
			}
			name := generateRepositoryName(index, rc, c.Config.Repositories())
			src := configsource.New("composer.json", false)
			if isPackagistFalse(rc) {
				err = src.AddRepository("packagist.org", false, false)
			} else {
				err = src.AddRepository(name, rc, false)
			}
			if err != nil {
				return 1, err
			}
			flushLock = true
			if c, err = composer.Create(io, composer.Options{DisablePlugins: cp.disablePlugins, DisableScripts: cp.disableScripts, FullLoad: true}); err != nil {
				return 1, err
			}
		}
	}
	if reqs := in.Strings("require"); len(reqs) > 0 {
		src := configsource.New("composer.json", false)
		flushLock = true
		for _, r := range reqs {
			parts := strings.Split(r, ":")
			constraint := "*"
			if len(parts) > 1 {
				constraint = parts[1]
			}
			if err := src.AddLink("require", parts[0], constraint); err != nil {
				return 1, err
			}
		}
		if c, err = composer.Create(io, composer.Options{DisablePlugins: cp.disablePlugins, DisableScripts: cp.disableScripts, FullLoad: true}); err != nil {
			return 1, err
		}
	}
	if flushLock {
		if st, err := os.Stat("composer.lock"); err == nil && st.Mode().IsRegular() {
			os.Remove("composer.lock")
		}
	}
	events := composer.NewEventDispatcher(c)
	events.RunScripts = !c.DisableScripts
	if _, err := events.DispatchScript("post-root-package-install", installDev, nil, nil); err != nil {
		return scriptExit(err), err
	}
	ps, pd, err := createProjectPreferredInstall(c.Config.Get("preferred-install"), in, false)
	if err != nil {
		return 1, err
	}
	if !in.Bool("no-install") {
		if err := checkPlugins(cp.ctx, c); err != nil {
			return 1, err
		}
		inst := composer.NewInstaller(c)
		inst.IM.OutputProgress = !in.Bool("no-progress")
		inst.PreferSource = ps
		inst.PreferDist = pd
		inst.DevMode = installDev
		inst.PlatformFilter = filter
		inst.OptimizeAutoloader = c.Config.GetBool("optimize-autoloader")
		inst.ClassMapAuthoritative = c.Config.GetBool("classmap-authoritative")
		inst.ApcuAutoloader = c.Config.GetBool("apcu-autoloader")
		pc, err := composer.CreatePolicyConfig(c, in.Bool("no-blocking") || in.Bool("no-security-blocking"), in.Bool("no-cooldown"))
		if err != nil {
			return 1, err
		}
		inst.PolicyConfig = pc
		inst.Audit = !in.Bool("no-audit")
		if f, ok := in.String("audit-format"); ok && f != "" {
			inst.AuditFormat = f
		}
		if rootPkg != nil {
			inst.SuggestionsFromInstall = []*pkg.Package{rootPkg}
		}
		if locked, _ := c.Locker.IsLocked(); !locked {
			inst.Update = true
		}
		status, err := inst.Run()
		if err != nil {
			return status, err
		}
		if status != 0 {
			return status, nil
		}
	}
	hasVcs := installedFromVcs
	if !in.Bool("keep-vcs") && installedFromVcs {
		remove := true
		if !in.Bool("remove-vcs") && io.Interactive {
			var err error
			if remove, err = io.AskConfirmation("<info>Do you want to remove the existing VCS (.git, .svn..) history?</info> [<comment>y,n</comment>]? ", true); err != nil {
				return 1, err
			}
		}
		if remove {
			cwd, _ := os.Getwd()
			if err := removeVcsMetadata(cwd); err != nil {
				io.WriteError("<error>An error occurred while removing the VCS metadata: " + err.Error() + "</error>")
			}
			hasVcs = false
		}
	}
	if !hasVcs {
		root := c.Root()
		src := configsource.New("composer.json", false)
		for _, kind := range pkg.LinkKinds {
			for _, link := range root.LinksOf(kind.Type) {
				if link.PrettyConstraint == "self.version" {
					if err := src.AddLink(kind.Key, link.Target, root.PrettyVersion()); err != nil {
						return 1, err
					}
				}
			}
		}
	}
	if _, err := events.DispatchScript("post-create-project-cmd", installDev, nil, nil); err != nil {
		return scriptExit(err), err
	}
	return 0, nil
}

var vcsMetadataNames = map[string]bool{".svn": true, "_svn": true, "CVS": true, "_darcs": true, ".arch-params": true, ".monotone": true, ".bzr": true, ".git": true, ".hg": true, ".fslckout": true, "_FOSSIL_": true}

// removeVcsMetadata removes the VCS directories at the top of dir (Finder
// depth 0 directories matching the VCS names).
func removeVcsMetadata(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if !vcsMetadataNames[e.Name()] {
			continue
		}
		p := filepath.Join(dir, e.Name())
		st, err := os.Stat(p)
		if err != nil || !st.IsDir() {
			continue
		}
		if err := fsutil.RemoveAll(p); err != nil {
			return &composer.UserError{Msg: "Could not remove " + p, Code: 1}
		}
	}
	return nil
}

func isPackagistFalse(rc phpjson.Value) bool {
	m := phpjson.ToMap(rc)
	if m == nil || m.Len() != 1 {
		return false
	}
	for _, k := range []string{"packagist", "packagist.org"} {
		if v, ok := m.Get(k); ok && v == false {
			return true
		}
	}
	return false
}

// generateRepositoryName is RepositoryFactory::generateRepositoryName.
func generateRepositoryName(index int, rc phpjson.Value, existing *phpjson.Map) string {
	name := ""
	if m := phpjson.ToMap(rc); m != nil && m.Isset("url") {
		u, _ := phpjson.ScalarString(get(m, "url"))
		name = regexp.MustCompile(`(?i)^https?://`).ReplaceAllString(u, "")
	} else {
		name = strconv.Itoa(index)
	}
	for existing != nil && existing.Has(name) {
		name += "2"
	}
	return name
}

func (cp *createProject) installRootPackage(packageName, directory, packageVersion string, preferSource, preferDist, installDev bool, repositories []string, hasRepositories bool) (*pkg.Package, error) {
	io, in := cp.ctx.IO, cp.in
	pairs := parseNameVersionPairs([]string{packageName})
	name := strings.ToLower(pairs[0].Name)
	if packageVersion == "" && pairs[0].HasVersion {
		packageVersion = pairs[0].Version
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if directory == "" {
		parts := strings.SplitN(name, "/", 2)
		directory = cwd + "/" + parts[len(parts)-1]
	}
	directory = strings.TrimRight(directory, `/\`)
	if !fsutil.IsAbsolutePath(directory) {
		directory = cwd + "/" + directory
	}
	if directory == "" {
		return nil, &composer.UserError{Msg: "Got an empty target directory, something went wrong", Code: 1}
	}
	config := cp.config
	config.SetBaseDir(directory)
	if in.Bool("no-secure-http") {
		sec := phpjson.NewMap()
		sub := phpjson.NewMap()
		sub.Set("secure-http", false)
		sec.Set("config", sub)
		config.MergeSource(sec, "command-line")
	}
	short, err := fsutil.FindShortestPath(cwd, directory, true, false)
	if err != nil {
		short = directory
	}
	io.WriteError(`<info>Creating a "` + packageName + `" project at "` + short + `"</info>`)
	if st, err := os.Stat(directory); err == nil {
		if !st.IsDir() {
			return nil, &composer.UserError{Msg: `Cannot create project directory at "` + directory + `", it exists as a file.`, Code: 1}
		}
		if !fsutil.IsDirEmpty(directory) {
			return nil, &composer.UserError{Msg: `Project directory "` + directory + `" is not empty.`, Code: 1}
		}
	}
	stability, hasStability := in.String("stability")
	if !hasStability || !in.present["stability"] {
		if packageVersion == "" {
			stability = "stable"
		} else if m := reStabilitySuffix.FindStringSubmatch(packageVersion); m != nil {
			stability = m[1]
		} else {
			stability = semver.ParseStability(packageVersion)
		}
	}
	stability, err = semver.NormalizeStability(stability)
	if err != nil {
		return nil, &composer.UserError{Msg: err.Error(), Code: 1}
	}
	all, err := config.All(0)
	if err != nil {
		return nil, err
	}
	c, err := composer.NewProjectComposer(io, all, directory, composer.Options{DisablePlugins: cp.disablePlugins, DisableScripts: cp.disableScripts})
	if err != nil {
		return nil, err
	}
	rs := solver.NewRepositorySet(stability, nil, nil, nil, nil, nil)
	if !hasRepositories {
		repos, err := c.DefaultRepositories()
		if err != nil {
			return nil, err
		}
		// CompositeRepository of the default repositories
		for _, r := range repos {
			rs.AddRepository(r)
		}
	} else {
		for _, r := range repositories {
			rc, err := repoConfigFromString(r, true)
			if err != nil {
				return nil, err
			}
			if isPackagistFalse(rc) {
				continue
			}
			if m := phpjson.ToMap(rc); m != nil && m.Str("type") == "path" {
				opts := m.Map("options")
				if opts == nil || !opts.Has("symlink") {
					if opts == nil {
						opts = phpjson.NewMap()
					}
					opts.Set("symlink", false)
					m.Set("options", opts)
				}
			}
			r, err := c.CreateRepositoryFromConfig(rc)
			if err != nil {
				return nil, err
			}
			rs.AddRepository(r)
		}
	}
	platformRepo, err := c.CreatePlatformRepo(c.Config.Get("platform"))
	if err != nil {
		return nil, err
	}
	var cooldown *policy.CooldownPolicyConfig
	pc, err := composer.CreatePolicyConfig(c, in.Bool("no-blocking") || in.Bool("no-security-blocking"), in.Bool("no-cooldown"))
	if err != nil {
		return nil, err
	}
	if pc.Enabled && pc.Cooldown.ShouldBlock(policy.BlockScopeUpdate) {
		cooldown = pc.Cooldown
	}
	filter := platformFilter(in)
	vs := composer.NewVersionSelector(rs, platformRepo)
	vs.SetCooldownPolicy(cooldown)
	p, err := vs.FindBestCandidate(name, packageVersion, stability, filter, 0, io, true)
	if err != nil {
		return nil, err
	}
	if p == nil {
		desc := "stability " + stability
		if packageVersion != "" {
			desc = "version " + packageVersion
		}
		msg := "Could not find package " + name + " with " + desc
		if _, ignoreAll := filter.(solver.IgnoreAll); !ignoreAll {
			vs2 := composer.NewVersionSelector(rs, platformRepo)
			vs2.SetCooldownPolicy(cooldown)
			if alt, err := vs2.FindBestCandidate(name, packageVersion, stability, solver.IgnoreAll{}, 0, nil, false); err == nil && alt != nil {
				return nil, &composer.UserError{Msg: msg + " in a version installable using your PHP version, PHP extensions and Composer version.", Code: 1}
			}
		}
		if cooldown != nil {
			vs3 := composer.NewVersionSelector(rs, platformRepo)
			if alt, err := vs3.FindBestCandidate(name, packageVersion, stability, filter, 0, nil, false); err == nil && alt != nil {
				return nil, &composer.UserError{Msg: msg + ` that is past the cooldown period configured in "policy.cooldown". ` + policy.CooldownRemedy, Code: 1}
			}
		}
		return nil, &composer.UserError{Msg: msg + ".", Code: 1}
	}
	_ = os.MkdirAll(directory, 0o777)
	if p.AliasOf != nil && p.PrettyVersion() == semver.DefaultBranchAlias {
		p = p.AliasOf
	}
	io.WriteError("<info>Installing " + p.Name + " (" + p.FullPrettyVersion(false) + ")</info>")
	if cp.disablePlugins {
		io.WriteError("<info>Plugins have been disabled.</info>")
	}
	if p.AliasOf != nil {
		p = p.AliasOf
	}
	im := composer.NewInstallationManager(c)
	im.SetPreferSource(preferSource)
	im.SetPreferDist(preferDist)
	im.OutputProgress = !in.Bool("no-progress")
	installedPkg, err := im.InstallProject(p, directory, true)
	if err != nil {
		return nil, err
	}
	io.WriteError("<info>Created project in " + directory + "</info>")
	if err := os.Chdir(directory); err != nil {
		return nil, err
	}
	if _, err := os.Stat(directory + "/composer.json"); err == nil {
		if _, ok := os.LookupEnv("COMPOSER"); ok {
			os.Unsetenv("COMPOSER")
		}
	}
	os.Setenv("COMPOSER_ROOT_VERSION", installedPkg.PrettyVersion())
	return installedPkg, nil
}
