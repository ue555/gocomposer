package composer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// LoadRootPackage ports RootPackageLoader::load.
func LoadRootPackage(data *phpjson.Map, cwd string, cfg *config.Config, io *cio.IO) (*pkg.Package, error) {
	conf := data.Clone()
	if !conf.Isset("name") {
		conf.Set("name", "__root__")
	} else if msg := repo.PackageNamingError(conf.Str("name"), false); msg != "" {
		return nil, &UserError{Msg: "Your package name " + msg, Code: 1}
	}
	autoVersioned := false
	if !conf.Isset("version") {
		if env := os.Getenv("COMPOSER_ROOT_VERSION"); env != "" {
			v, err := rootVersionFromEnv(env)
			if err != nil {
				return nil, err
			}
			conf.Set("version", v)
		} else if pretty, normalized, commit, ok := guessRootVersion(conf, cwd, cfg); ok {
			conf.Set("version", pretty)
			conf.Set("version_normalized", normalized)
			if commit != "" {
				src := phpjson.NewMap()
				src.Set("type", "")
				src.Set("url", "")
				src.Set("reference", commit)
				conf.Set("source", src)
				dist := phpjson.NewMap()
				dist.Set("type", "")
				dist.Set("url", "")
				dist.Set("reference", commit)
				conf.Set("dist", dist)
			}
		}
		if !conf.Isset("version") {
			if io != nil && conf.Str("name") != "__root__" && conf.Str("type") != "project" {
				io.WriteError("<warning>Composer could not detect the root package (" + conf.Str("name") + ") version, defaulting to '1.0.0'. See https://getcomposer.org/root-version</warning>")
			}
			conf.Set("version", "1.0.0")
			autoVersioned = true
		}
	}
	loader := &pkg.Loader{}
	p, err := loader.Load(conf)
	if err != nil {
		return nil, userErr(err)
	}
	real := p.Base()
	real.IsRoot = true
	p.IsRoot = true
	real.MinimumStability = "stable"
	if autoVersioned {
		real.ReplaceVersion(real.Version(), pkg.DefaultPrettyVersion)
	}
	if ms, ok := conf.Get("minimum-stability"); ok && ms != nil {
		s, _ := phpjson.ScalarString(ms)
		n, err := semver.NormalizeStability(s)
		if err != nil {
			return nil, userErr(err)
		}
		real.MinimumStability = n
	}
	var aliases []pkg.RootAlias
	flags := map[string]int{}
	refs := map[string]string{}
	for _, kind := range []struct{ key, typ string }{{"require", pkg.TypeRequire}, {"require-dev", pkg.TypeDevRequire}} {
		if !conf.Isset(kind.key) {
			continue
		}
		var names []string
		links := map[string]string{}
		for _, l := range real.LinksOf(kind.typ) {
			if _, ok := links[l.Target]; !ok {
				names = append(names, l.Target)
			}
			links[l.Target] = l.Constraint.PrettyString()
		}
		a, err := extractAliases(names, links, aliases)
		if err != nil {
			return nil, err
		}
		aliases = a
		flags = ExtractStabilityFlags(names, links, real.MinimumStability, flags)
		refs = extractReferences(names, links, refs)
		if _, ok := links[conf.Str("name")]; ok {
			return nil, &UserError{Msg: fmt.Sprintf("Root package '%s' cannot require itself in its composer.json\nDid you accidentally name your root package after an external package?", conf.Str("name")), Code: 1}
		}
	}
	for _, kind := range pkg.LinkKinds {
		v, ok := conf.Get(kind.Key)
		if !ok || v == nil {
			continue
		}
		var lerr error
		phpjson.Each(v, func(linkName string, _ phpjson.Value) {
			if lerr != nil {
				return
			}
			if msg := repo.PackageNamingError(linkName, true); msg != "" {
				lerr = &UserError{Msg: kind.Key + "." + msg, Code: 1}
			}
		})
		if lerr != nil {
			return nil, lerr
		}
	}
	real.Aliases = aliases
	real.StabilityFlags = flags
	real.References = refs
	if ps, ok := conf.Get("prefer-stable"); ok && ps != nil {
		real.PreferStable = phpjson.Truthy(ps)
	}
	if c := conf.Map("config"); c != nil {
		real.RootConfig = c
	}
	return p, nil
}

func rootVersionFromEnv(v string) (string, error) {
	n, err := semver.Normalize(v)
	if err != nil {
		return "", userErr(err)
	}
	_ = n
	return v, nil
}

var (
	reAliasDef    = regexp.MustCompile(`(?:^|\| *|, *)([^,\s#|]+)(?:#[^ ]+)? +as +([^,\s|]+)(?:$| *\|| *,)`)
	reStabFlag    = regexp.MustCompile(`(?i)^[^@]*?@(stable|RC|beta|alpha|dev)$`)
	reAsSuffix    = regexp.MustCompile(`^([^,\s@]+) as .+$`)
	reNoSpaceAt   = regexp.MustCompile(`^[^,\s@]+$`)
	reRefConstr   = regexp.MustCompile(`^[^,\s@]+?#([a-f0-9]+)$`)
	reOrSplitRoot = regexp.MustCompile(`\s*\|\|?\s*`)
)

func extractAliases(names []string, links map[string]string, aliases []pkg.RootAlias) ([]pkg.RootAlias, error) {
	for _, name := range names {
		reqVersion := links[name]
		if m := reAliasDef.FindStringSubmatch(reqVersion); m != nil {
			v, err := semver.NormalizeFull(m[1], reqVersion)
			if err != nil {
				return nil, userErr(err)
			}
			an, err := semver.NormalizeFull(m[2], reqVersion)
			if err != nil {
				return nil, userErr(err)
			}
			aliases = append(aliases, pkg.RootAlias{Package: strings.ToLower(name), Version: v, Alias: m[2], AliasNormalized: an})
		} else if strings.Contains(reqVersion, " as ") {
			return nil, &UserError{Msg: `Invalid alias definition in "` + name + `": "` + reqVersion + `". Aliases should be in the form "exact-version as other-exact-version".`, Code: 1}
		}
	}
	return aliases, nil
}

// ExtractStabilityFlags is RootPackageLoader::extractStabilityFlags.
func ExtractStabilityFlags(names []string, links map[string]string, minimumStability string, flags map[string]int) map[string]int {
	minStab := pkg.Stabilities[minimumStability]
	for _, reqName := range names {
		reqVersion := links[reqName]
		var constraints []string
		for _, or := range reOrSplitRoot.Split(strings.TrimSpace(reqVersion), -1) {
			constraints = append(constraints, semver.SplitAnd(or)...)
		}
		matched := false
		for _, c := range constraints {
			if m := reStabFlag.FindStringSubmatch(c); m != nil {
				name := strings.ToLower(reqName)
				ns, _ := semver.NormalizeStability(m[1])
				stab := pkg.Stabilities[ns]
				if cur, ok := flags[name]; ok && cur > stab {
					continue
				}
				flags[name] = stab
				matched = true
			}
		}
		if matched {
			continue
		}
		for _, c := range constraints {
			rv := reAsSuffix.ReplaceAllString(c, "$1")
			if reNoSpaceAt.MatchString(rv) {
				stabName := semver.ParseStability(rv)
				if stabName == "stable" {
					continue
				}
				name := strings.ToLower(reqName)
				stab := pkg.Stabilities[stabName]
				if cur, ok := flags[name]; (ok && cur > stab) || minStab > stab {
					continue
				}
				flags[name] = stab
			}
		}
	}
	return flags
}

func extractReferences(names []string, links map[string]string, refs map[string]string) map[string]string {
	for _, reqName := range names {
		rv := reAsSuffix.ReplaceAllString(links[reqName], "$1")
		if m := reRefConstr.FindStringSubmatch(rv); m != nil && semver.ParseStability(rv) == "dev" {
			refs[strings.ToLower(reqName)] = m[1]
		}
	}
	return refs
}

// guessRootVersion runs the VersionGuesser for the root package directory.
func guessRootVersion(conf *phpjson.Map, cwd string, cfg *config.Config) (string, string, string, bool) {
	pretty, normalized, ok := GuessVersion(cwd, conf, cfg)
	if !ok {
		return "", "", "", false
	}
	return pretty, normalized, gitCommit(cwd), true
}

// ExtractReferences is RootPackageLoader::extractReferences.
func ExtractReferences(names []string, links map[string]string, refs map[string]string) map[string]string {
	return extractReferences(names, links, refs)
}
