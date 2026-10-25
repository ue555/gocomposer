package cli

// Port of Composer\Command\ArchiveCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// The archive itself is built by internal/archiver (ArchiveManager and
// PharArchiver).

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/archiver"
	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

func init() { registerCommand(archiveCommand) }

func archiveCommand() *Command {
	return &Command{Name: "archive", Run: runArchive}
}

func runArchive(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	// tryComposer: a missing or invalid composer.json gives no instance
	c, err := ctx.composerFor(true)
	if err != nil {
		var pe *jsonfile.ParsingError
		if errors.As(err, &pe) {
			return 1, err
		}
		c = nil
	}
	var cfg *config.Config
	if c != nil {
		cfg = c.Config
		if code, err := composer.NewEventDispatcher(c).DispatchScript("pre-archive-cmd", true, nil, nil); err != nil || code != 0 {
			return code, err
		}
	} else {
		cwd, _ := os.Getwd()
		cfg, err = config.CreateConfig(cwd, func(m string) { io.WriteError(m) })
		if err != nil {
			return 1, err
		}
	}
	format, ok := in.String("format")
	if !ok {
		format = configString(cfg.Get("archive-format"))
	}
	dir, ok := in.String("dir")
	if !ok {
		dir = configString(cfg.Get("archive-dir"))
	}
	fileName, _ := in.String("file")
	code, err := archiveRun(ctx, c, in.Arg("package"), in.Args("version"), format, dir, fileName, in.Bool("ignore-filters"))
	if err != nil {
		return 1, err
	}
	if code == 0 && c != nil {
		if code, err := composer.NewEventDispatcher(c).DispatchScript("post-archive-cmd", true, nil, nil); err != nil || code != 0 {
			return code, err
		}
	}
	return code, nil
}

func configString(v phpjson.Value) string {
	s, _ := phpjson.ScalarString(v)
	return s
}

// archiveRun is ArchiveCommand::archive.
func archiveRun(ctx *Context, c *composer.Composer, packageName string, version []string, format, dest, fileName string, ignoreFilters bool) (int, error) {
	io := ctx.IO
	var p *pkg.Package
	if packageName != "" && packageName != "0" {
		var err error
		p, err = archiveSelectPackage(ctx, c, packageName, version)
		if err != nil {
			return 1, err
		}
		if p == nil {
			return 1, nil
		}
	} else {
		if c == nil {
			// requireComposer: the Factory error is printed as is
			_, err := ctx.composerFor(true)
			var ue *composer.UserError
			if errors.As(err, &ue) && strings.HasPrefix(ue.Msg, "Composer could not find ") {
				io.WriteError(ue.Msg)
				return 1, nil
			}
			if err == nil {
				err = &composer.UserError{Msg: "Composer could not be loaded", Code: 1}
			}
			return 1, err
		}
		p = c.Package
	}
	io.WriteError(`<info>Creating the archive into "` + dest + `".</info>`)
	// $composer->getArchiveManager() when plugins changed it (plugin_archive.go)
	packagePath, handled, err := composer.PluginArchive(c, p, format, dest, fileName, ignoreFilters)
	if !handled {
		packagePath, err = archivePackage(c, p, format, dest, fileName, ignoreFilters)
	}
	if err != nil {
		return 1, err
	}
	cwd, _ := os.Getwd()
	shortPath, err := fsutil.FindShortestPath(cwd, packagePath, true, false)
	if err != nil || len(shortPath) >= len(packagePath) {
		shortPath = packagePath
	}
	io.WriteErrorV("Created: ", false, 0)
	io.Write(shortPath)
	return 0, nil
}

// archivePackage is ArchiveManager::archive.
func archivePackage(c *composer.Composer, p *pkg.Package, format, targetDir, fileName string, ignoreFilters bool) (string, error) {
	if format == "" {
		return "", &composer.UserError{Msg: "Format must be specified", Code: 1}
	}
	if !archiver.Supports(format) {
		return "", &composer.UserError{Msg: "No archiver found to support " + format + " format", Code: 1}
	}
	info := archiver.PackageInfo{Name: p.Name, PrettyVersion: p.PrettyVersion(), DistReference: p.DistReference(), DistType: p.DistType(), SourceReference: p.SourceReference()}
	if ar := pkg.Dump(p).Map("archive"); ar != nil {
		info.ArchiveName = ar.Str("name")
		info.ArchiveExcludes = stringList(get(ar, "exclude"))
	}
	var sourcePath string
	if p.IsRoot {
		cwd, _ := os.Getwd()
		sourcePath = fsutil.Realpath(cwd)
	} else {
		sourcePath = archiver.TempName()
		if err := archiver.EnsureDirectoryExists(sourcePath); err != nil {
			return "", err
		}
		im := composer.NewInstallationManager(c)
		if err := im.DownloadTo(p, sourcePath); err != nil {
			os.RemoveAll(sourcePath)
			return "", err
		}
		defer os.RemoveAll(sourcePath)
		// check exclude from downloaded composer.json
		if raw, err := os.ReadFile(sourcePath + "/composer.json"); err == nil {
			v, err := jsonfile.Parse(raw, sourcePath+"/composer.json")
			if err != nil {
				return "", err
			}
			if ar := phpjson.ToMap(v).Map("archive"); ar != nil {
				if n, ok := ar.Get("name"); ok && phpjson.Truthy(n) {
					info.ArchiveName, _ = phpjson.ScalarString(n)
				}
				if ex, ok := ar.Get("exclude"); ok && phpjson.Truthy(ex) {
					info.ArchiveExcludes = stringList(ex)
				}
			}
		}
	}
	target, err := archiver.ArchivePackage(info, sourcePath, format, targetDir, fileName, ignoreFilters)
	if err != nil {
		var ae *archiver.ArchiveError
		var de *archiver.DirError
		var se *archiver.SourceError
		var pe *archiver.PatternError
		if errors.As(err, &ae) || errors.As(err, &de) || errors.As(err, &se) || errors.As(err, &pe) {
			return "", &composer.UserError{Msg: err.Error(), Code: 1}
		}
		return "", err
	}
	return target, nil
}

// stringList reads a list of strings (dumped []string or decoded JSON).
func stringList(v phpjson.Value) []string {
	if l, ok := v.([]string); ok {
		return append([]string(nil), l...)
	}
	var out []string
	phpjson.Each(v, func(_ string, x phpjson.Value) {
		s, _ := phpjson.ScalarString(x)
		out = append(out, s)
	})
	return out
}

var reArchiveStabilitySuffix = regexp.MustCompile(`(?i)@(stable|RC|beta|alpha|dev)$`)

// archiveSelectPackage is ArchiveCommand::selectPackage.
func archiveSelectPackage(ctx *Context, c *composer.Composer, packageName string, versionArg []string) (*pkg.Package, error) {
	io := ctx.IO
	io.WriteError("<info>Searching for the specified package.</info>")
	if c == nil {
		return nil, &composer.UserError{Msg: "gocomposer cannot archive a package without a composer.json in the current directory (searching the default repositories is not supported)", Code: 1}
	}
	repos := append([]repo.Repository{c.LocalRepo}, c.Repositories...)
	minStability := c.Package.MinimumStability
	hasVersion := len(versionArg) > 0
	version := ""
	if hasVersion {
		version = versionArg[0]
		if m := reArchiveStabilitySuffix.FindStringSubmatch(version); m != nil {
			if s, err := semver.NormalizeStability(m[1]); err == nil {
				minStability = s
			}
			version = version[:len(version)-len(m[0])]
		}
	}
	rs := solver.NewRepositorySet(minStability, nil, nil, nil, nil, nil)
	// RepositorySet::addRepository flattens the CompositeRepository
	for _, r := range repos {
		rs.AddRepository(r)
	}
	var constraint semver.Constraint
	if hasVersion {
		cst, err := semver.ParseConstraints(version)
		if err != nil {
			return nil, err
		}
		constraint = cst
	}
	name := strings.ToLower(packageName)
	packages, err := rs.FindPackagesE(name, constraint, 0)
	if err != nil {
		return nil, err
	}
	var p *pkg.Package
	switch {
	case len(packages) > 1:
		vs := composer.NewVersionSelector(rs, nil)
		best, err := vs.FindBestCandidate(name, version, minStability, solver.IgnoreNothing{}, 0, nil, true)
		if err != nil {
			return nil, err
		}
		p = best
		if p == nil {
			p = packages[0]
		}
		io.WriteError("<info>Found multiple matches, selected " + p.PrettyString() + ".</info>")
		var alts []string
		for _, a := range packages {
			alts = append(alts, a.PrettyString())
		}
		io.WriteError("Alternatives were " + strings.Join(alts, ", ") + ".")
		io.WriteError("<comment>Please use a more specific constraint to pick a different package.</comment>")
	case len(packages) == 1:
		p = packages[0]
		io.WriteError("<info>Found an exact match " + p.PrettyString() + ".</info>")
	default:
		io.WriteError("<error>Could not find a package matching " + packageName + ".</error>")
		return nil, nil
	}
	return p, nil
}
