package cli

// Port of the completion suggestions of Composer's commands
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano): Command\CompletionTrait, the suggested
// values of the InputOption/InputArgument definitions of each command
// (Console\Input\InputOption, InputArgument, BaseCommand::complete),
// ConfigCommand::suggestSettingKeys/flattenSettingKeys,
// RepositoryCommand::suggestRepoNames/suggestTypeForAdd,
// PolicyCommand::suggestListNames/suggestArg1, ShowCommand::
// suggestPackageBasedOnMode, RunScriptCommand::getScripts and
// ExecCommand::getBinaries.

import (
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
)

// createConfigNoIO is Factory::createConfig() without an IO.
func createConfigNoIO(cwd string) (*config.Config, error) { return config.CreateConfig(cwd, nil) }

type suggester func() ([]string, error)

func staticValues(vs ...string) suggester {
	return func() ([]string, error) { return vs, nil }
}

var (
	suggestPreferInstall = []string{"dist", "source", "auto"}
	suggestAuditFormat   = []string{"table", "plain", "json", "summary"}
)

// optionSuggester returns the suggested values of a command option
// (the $suggestedValues of its InputOption), or nil.
func (s *compSession) optionSuggester(cmd, opt string, ci *compInput) suggester {
	switch opt {
	case "prefer-install":
		switch cmd {
		case "install", "create-project", "update", "require", "reinstall":
			return staticValues(suggestPreferInstall...)
		}
	case "audit-format":
		switch cmd {
		case "install", "create-project", "update", "require", "remove":
			return staticValues(suggestAuditFormat...)
		}
	case "bump-after-update":
		if cmd == "update" {
			return staticValues("dev", "no-dev", "all")
		}
	case "format":
		switch cmd {
		case "search", "show", "outdated", "check-platform-reqs":
			return staticValues("json", "text")
		case "audit":
			return staticValues(suggestAuditFormat...)
		case "archive":
			return staticValues("tar", "tar.gz", "tar.bz2", "zip")
		case "licenses":
			return staticValues("text", "json", "summary")
		case "fund":
			return staticValues("text", "json")
		}
	case "abandoned":
		if cmd == "audit" {
			return staticValues("ignore", "report", "fail")
		}
	case "ignore-severity":
		if cmd == "audit" {
			return staticValues("low", "medium", "high", "critical")
		}
	case "ignore":
		if cmd == "show" || cmd == "outdated" {
			return s.installedPackage(ci, false, false)
		}
	case "require", "require-dev":
		if cmd == "init" {
			return s.availablePackageInclPlatform(ci)
		}
	case "type":
		if cmd == "reinstall" {
			return s.installedPackageTypes(false)
		}
	case "before", "after":
		if cmd == "repository" {
			return s.repoNames(ci)
		}
	}
	return nil
}

// argumentSuggester returns the suggested values of a command argument
// (the $suggestedValues of its InputArgument), or nil.
func (s *compSession) argumentSuggester(cmd, arg string, ci *compInput) suggester {
	switch cmd + " " + arg {
	case "config setting-key":
		return s.settingKeys(ci)
	case "depends package":
		return s.installedPackage(ci, true, true)
	case "prohibits package", "create-project package", "archive package":
		return s.availablePackage(ci, 99)
	case "update packages", "outdated package", "reinstall packages":
		return s.installedPackage(ci, false, false)
	case "show package":
		return func() ([]string, error) {
			if ci.optionTruthy("available") || ci.optionTruthy("all") {
				return s.availablePackageInclPlatform(ci)()
			}
			if ci.optionTruthy("platform") {
				return s.platformPackage(ci)()
			}
			return s.installedPackage(ci, false, false)()
		}
	case "suggests packages", "browse packages":
		return s.installedPackage(ci, true, false)
	case "require packages":
		return s.availablePackageInclPlatform(ci)
	case "run-script script":
		return s.scripts()
	case "remove packages", "bump packages":
		return s.rootRequirement()
	case "exec binary":
		return s.binaries()
	case "repository action":
		return staticValues("list", "add", "remove", "set-url", "get-url", "enable", "disable")
	case "repository name":
		return s.repoNames(ci)
	case "repository arg1":
		return func() ([]string, error) {
			if a := ci.getArgument("action"); a != nil && *a == "add" {
				return []string{"composer", "vcs", "artifact", "path"}, nil
			}
			return nil, nil
		}
	case "policy action":
		return staticValues("add-source")
	case "policy name":
		return s.policyListNames(ci)
	case "policy arg1":
		return func() ([]string, error) {
			if a := ci.getArgument("action"); a != nil && *a == "add-source" {
				return []string{"url"}, nil
			}
			return nil, nil
		}
	}
	return nil
}

// definitionSuggestions is BaseCommand::complete: the suggested values of
// the option or argument being completed.
func (s *compSession) definitionSuggestions(cmd *Command, ci *compInput) ([]string, error) {
	def := mergedDefinition(cmd)
	name := ci.completionName()
	var sg suggester
	switch ci.ctype {
	case compOptionValue:
		if _, ok := def.option(name); ok {
			sg = s.optionSuggester(cmd.Name, name, ci)
		}
	case compArgumentValue:
		if _, ok := def.argument(name); ok {
			sg = s.argumentSuggester(cmd.Name, name, ci)
		}
	}
	if sg == nil {
		return nil, nil
	}
	return sg()
}

// linkKeys are array_keys() of a links array: the targets, and numeric keys
// for the links an alias package appends.
func linkKeys(ls pkg.Links) []string {
	var out []string
	idx := 0
	for _, l := range ls {
		if l.Indexed {
			out = append(out, strconv.Itoa(idx))
			idx++
			continue
		}
		out = append(out, strings.ToLower(l.Target))
	}
	return out
}

// rootRequirement is CompletionTrait::suggestRootRequirement.
func (s *compSession) rootRequirement() suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		return append(linkKeys(c.Package.Requires()), linkKeys(c.Package.DevRequires())...), nil
	}
}

// rootPackageNames are the packages of a RootPackageRepository (an alias
// root package also adds the aliased package).
func rootPackageNames(p *pkg.Package) []string {
	if p.AliasOf != nil {
		return []string{p.Name, p.AliasOf.Name}
	}
	return []string{p.Name}
}

// installedPackages are the locked packages (with dev) when locked,
// else the installed ones.
func (s *compSession) installedPackages() ([]*pkg.Package, bool, error) {
	c, err := s.requireComposer()
	if err != nil {
		return nil, false, err
	}
	locked, err := c.Locker.IsLocked()
	if err != nil {
		return nil, false, err
	}
	if locked {
		r, err := c.Locker.LockedRepository(true)
		if err != nil {
			return nil, false, err
		}
		ps, err := r.Packages()
		return ps, true, err
	}
	ps, err := c.LocalRepo.Packages()
	return ps, false, err
}

// installedPackage is CompletionTrait::suggestInstalledPackage.
func (s *compSession) installedPackage(ci *compInput, includeRoot, includePlatform bool) suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		var names []string
		if includeRoot {
			names = append(names, rootPackageNames(c.Package)...)
		}
		ps, locked, err := s.installedPackages()
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			names = append(names, p.Name)
		}
		var hints []string
		if includePlatform {
			var overrides phpjson.Value
			if locked {
				overrides = c.Locker.PlatformOverrides()
			} else {
				overrides = c.Config.Get("platform")
			}
			pr, err := c.CreatePlatformRepo(overrides)
			if err != nil {
				return nil, err
			}
			pps, err := pr.Packages()
			if err != nil {
				return nil, err
			}
			if ci.cvalue == "" {
				// to reduce noise, when no text is yet entered we list only two entries for ext- and lib- prefixes
				type hint struct {
					prefix string
					count  int
					gone   bool
				}
				find := []*hint{{"ext-", 0, false}, {"lib-", 0, false}, {"php", 99, false}, {"composer", 99, false}}
				for _, p := range pps {
					for _, h := range find {
						if h.gone || !strings.HasPrefix(p.Name, h.prefix) {
							continue
						}
						if h.count == 0 || h.count >= 99 {
							hints = append(hints, p.Name)
							h.count++
						} else if h.count == 1 {
							h.gone = true
							n := len(p.Name) - 3
							if m := len(h.prefix) + 1; m > n {
								n = m
							}
							if n > len(p.Name) {
								n = len(p.Name)
							}
							hints = append(hints, p.Name[:n]+"...")
						}
						break
					}
				}
			} else {
				for _, p := range pps {
					names = append(names, p.Name)
				}
			}
		}
		return append(names, hints...), nil
	}
}

// installedPackageTypes is CompletionTrait::suggestInstalledPackageTypes.
func (s *compSession) installedPackageTypes(includeRoot bool) suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		var types []string
		if includeRoot {
			types = append(types, c.Package.Type())
			if c.Package.AliasOf != nil {
				types = append(types, c.Package.AliasOf.Type())
			}
		}
		ps, _, err := s.installedPackages()
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			types = append(types, p.Type())
		}
		return phpArrayUnique(types), nil
	}
}

func phpArrayUnique(vs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range vs {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// searchAll is CompositeRepository::search over the configured repositories.
func (s *compSession) searchAll(query string, mode int) ([]string, error) {
	c, err := s.requireComposer()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, r := range c.Repositories {
		res, err := repo.SearchWithMode(r, query, mode, nil)
		if err != nil {
			return nil, err
		}
		for _, e := range res {
			names = append(names, e.Name)
		}
	}
	return names, nil
}

// availablePackage is CompletionTrait::suggestAvailablePackage.
func (s *compSession) availablePackage(ci *compInput, max int) suggester {
	return func() ([]string, error) {
		if max < 1 {
			return nil, nil
		}
		value := ci.cvalue
		var results []string
		showVendors := false
		if !strings.Contains(value, "/") {
			r, err := s.searchAll("^"+pregQuote(value), repo.SearchVendor)
			if err != nil {
				return nil, err
			}
			results = r
			showVendors = true
		}
		// if we get a single vendor, we expand it into its contents already
		if len(results) <= 1 {
			r, err := s.searchAll("^"+pregQuote(value), repo.SearchName)
			if err != nil {
				return nil, err
			}
			results = r
			showVendors = false
		}
		if showVendors {
			for i := range results {
				results[i] += "/"
			}
			// sort shorter results first to avoid auto-expanding the completion to a longer string than needed
			sort.SliceStable(results, func(i, j int) bool {
				a, b := results[i], results[j]
				if len(a) == len(b) {
					return phpjson.PHPCompareKeys(a, b) < 0
				}
				return len(a) < len(b)
			})
			var pinned []string
			// ensure if the input is an exact match that it is always in the result set
			exact := value + "/"
			for i, r := range results {
				if r == exact {
					pinned = append(pinned, exact)
					results = append(results[:i:i], results[i+1:]...)
					break
				}
			}
			n := max - len(pinned)
			if n < len(results) {
				if n < 0 {
					n = 0
				}
				results = results[:n]
			}
			return append(pinned, results...), nil
		}
		if len(results) > max {
			results = results[:max]
		}
		return results, nil
	}
}

var rePlatformPrefix = regexp.MustCompile(`^(ext|lib|php)(-|$)|^com`)

// availablePackageInclPlatform is CompletionTrait::suggestAvailablePackageInclPlatform.
func (s *compSession) availablePackageInclPlatform(ci *compInput) suggester {
	return func() ([]string, error) {
		var matches []string
		if rePlatformPrefix.MatchString(ci.cvalue) {
			m, err := s.platformPackage(ci)()
			if err != nil {
				return nil, err
			}
			matches = m
		}
		rest, err := s.availablePackage(ci, 99-len(matches))()
		if err != nil {
			return nil, err
		}
		return append(matches, rest...), nil
	}
}

// platformPackage is CompletionTrait::suggestPlatformPackage.
func (s *compSession) platformPackage(ci *compInput) suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		pr, err := c.CreatePlatformRepo(c.Config.Get("platform"))
		if err != nil {
			return nil, err
		}
		pps, err := pr.Packages()
		if err != nil {
			return nil, err
		}
		// BasePackage::packageNameToRegexp($value.'*')
		re, err := regexp.Compile(`(?i)^` + strings.ReplaceAll(pregQuote(ci.cvalue+"*"), `\*`, ".*") + `$`)
		if err != nil {
			return nil, &errCompletion{err.Error()}
		}
		var out []string
		for _, p := range pps {
			if re.MatchString(p.Name) {
				out = append(out, p.Name)
			}
		}
		return out, nil
	}
}

// scripts is the script argument suggestion of RunScriptCommand.
func (s *compSession) scripts() suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		if sc := c.Package.Scripts(); sc != nil {
			return sc.Keys(), nil
		}
		return nil, nil
	}
}

// binaries is ExecCommand::getBinaries(false).
func (s *compSession) binaries() suggester {
	return func() ([]string, error) {
		c, err := s.requireComposer()
		if err != nil {
			return nil, err
		}
		binDir := c.Config.GetString("bin-dir")
		var all []string
		if entries, err := os.ReadDir(binDir); err == nil {
			for _, e := range entries {
				if strings.HasPrefix(e.Name(), ".") {
					continue
				}
				all = append(all, binDir+"/"+e.Name())
			}
		}
		all = append(all, c.Package.Binaries()...)
		var out []string
		prev := ""
		hasPrev := false
		for _, b := range all {
			// skip .bat copies
			if hasPrev && b == prev+".bat" {
				continue
			}
			prev, hasPrev = b, true
			out = append(out, phpBasename(b))
		}
		return out, nil
	}
}

// composerConfigFile is BaseConfigCommand::getComposerConfigFile.
func composerConfigFile(ci *compInput, cfg *config.Config) (string, error) {
	if ci.optionTruthy("global") {
		return cfg.GetString("home") + "/config.json", nil
	}
	if f := ci.optionString("file"); f != nil {
		return *f, nil
	}
	return composerFileName()
}

// settingKeys is ConfigCommand::suggestSettingKeys.
func (s *compSession) settingKeys(ci *compInput) suggester {
	return func() ([]string, error) {
		if ci.optionTruthy("list") || ci.optionTruthy("editor") || ci.optionTruthy("auth") {
			return nil, nil
		}
		cwd, _ := os.Getwd()
		cfg, err := createConfigNoIO(cwd)
		if err != nil {
			return nil, err
		}
		configFile, err := composerConfigFile(ci, cfg)
		if err != nil {
			return nil, err
		}
		var authFile string
		if ci.optionTruthy("global") {
			authFile = cfg.GetString("home") + "/auth.json"
		} else {
			authFile = phpDirname(configFile) + "/auth.json"
		}
		var fileData *phpjson.Map
		if fileExists(configFile) {
			m, err := readJSONMap(configFile)
			if err != nil {
				return nil, err
			}
			fileData = m
			cfg.MergeSource(m, configFile)
		}
		if fileExists(authFile) {
			m, err := readJSONMap(authFile)
			if err != nil {
				return nil, err
			}
			wrap := phpjson.NewMap()
			wrap.Set("config", m)
			cfg.MergeSource(wrap, authFile)
		}
		// collect all configuration setting-keys
		keys := append(flattenSettingKeys(cfg.Raw(), ""), flattenSettingKeys(cfg.Repositories(), "repositories.")...)
		if ci.optionTruthy("unset") {
			// keep only the currently customized setting-keys
			var kept []string
			for _, k := range keys {
				src, err := cfg.SourceOfValue(k)
				if err != nil {
					return nil, err
				}
				if src == configFile || src == authFile {
					kept = append(kept, k)
				}
			}
			keys = kept
		} else {
			keys = append(keys, configurablePackageProperties...)
		}
		// add all existing configurable package-properties
		if fileData != nil {
			props := phpjson.NewMap()
			for _, k := range fileData.Keys() {
				if containsStr(configurablePackageProperties, k) {
					v, _ := fileData.Get(k)
					props.Set(k, v)
				}
			}
			keys = append(keys, flattenSettingKeys(props, "")...)
		}
		// filter settings-keys by completion value
		if v := ci.cvalue; v != "" {
			var kept []string
			for _, k := range keys {
				if strings.HasPrefix(k, v) {
					kept = append(kept, k)
				}
			}
			keys = kept
		}
		sort.SliceStable(keys, func(i, j int) bool { return phpjson.PHPCompareKeys(keys[i], keys[j]) < 0 })
		return phpArrayUnique(keys), nil
	}
}

// flattenSettingKeys is ConfigCommand::flattenSettingKeys.
func flattenSettingKeys(v phpjson.Value, prefix string) []string {
	var keys []string
	phpjson.Each(v, func(k string, val phpjson.Value) {
		keys = append(keys, prefix+k)
		// array-lists must not be added to completion
		// sub-keys of repository-keys must not be added to completion
		if m, ok := val.(*phpjson.Map); ok && m != nil && m.Len() > 0 && !m.IsList() && prefix != "repositories." {
			keys = append(keys, flattenSettingKeys(m, prefix+k+".")...)
		}
	})
	return keys
}

// composerFileName is Factory::getComposerFile.
func composerFileName() (string, error) {
	if env, ok := os.LookupEnv("COMPOSER"); ok {
		env = strings.TrimSpace(env)
		if env != "" {
			if st, err := os.Stat(env); err == nil && st.IsDir() {
				return "", &errCompletion{"The COMPOSER environment variable is set to " + env + " which is a directory, this variable should point to a composer.json or be left unset."}
			}
			return env, nil
		}
	}
	return "./composer.json", nil
}

// repoNames is RepositoryCommand::suggestRepoNames.
func (s *compSession) repoNames(ci *compInput) suggester {
	return func() ([]string, error) {
		action := ci.getArgument("action")
		if action != nil && (*action == "enable" || *action == "disable") {
			return []string{"packagist.org"}, nil
		}
		if action == nil || (*action != "remove" && *action != "set-url" && *action != "get-url") {
			return nil, nil
		}
		cwd, _ := os.Getwd()
		cfg, err := createConfigNoIO(cwd)
		if err != nil {
			return nil, err
		}
		file, err := composerConfigFile(ci, cfg)
		if err != nil {
			return nil, err
		}
		data, err := readJSONMap(file)
		if err != nil {
			return nil, err
		}
		var repos []string
		if rv, ok := data.Get("repositories"); ok {
			phpjson.Each(rv, func(_ string, r phpjson.Value) {
				if m, ok := r.(*phpjson.Map); ok && m != nil {
					if n, ok := m.Get("name"); ok && n != nil {
						if s, ok := phpjson.ScalarString(n); ok {
							repos = append(repos, s)
						}
					}
				}
			})
		}
		sort.SliceStable(repos, func(i, j int) bool { return phpjson.PHPCompareKeys(repos[i], repos[j]) < 0 })
		return repos, nil
	}
}

// policyListNames is PolicyCommand::suggestListNames.
func (s *compSession) policyListNames(ci *compInput) suggester {
	return func() ([]string, error) {
		action := ci.getArgument("action")
		if action == nil || *action != "add-source" {
			return nil, nil
		}
		cwd, _ := os.Getwd()
		cfg, err := createConfigNoIO(cwd)
		if err != nil {
			return nil, err
		}
		file, err := composerConfigFile(ci, cfg)
		if err != nil {
			return nil, err
		}
		if !fileExists(file) {
			return nil, nil
		}
		data, err := readJSONMap(file)
		if err != nil {
			return nil, err
		}
		pv, _ := data.Get("config")
		var pol phpjson.Value
		if m := phpjson.AsMap(pv); m != nil {
			pol, _ = m.Get("policy")
		}
		if pol != nil && !phpjson.IsArray(pol) {
			return nil, nil
		}
		var names []string
		phpjson.Each(pol, func(k string, _ phpjson.Value) {
			if containsStr(policy.BuiltinListNames, k) || containsStr(policy.NonListKeys, k) {
				return
			}
			names = append(names, k)
		})
		sort.SliceStable(names, func(i, j int) bool { return phpjson.PHPCompareKeys(names[i], names[j]) < 0 })
		return names, nil
	}
}
