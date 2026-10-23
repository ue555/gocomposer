package pkg

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// Port of Composer\Package\Loader\ArrayLoader.

// LoadError mirrors \UnexpectedValueException thrown while loading.
type LoadError struct {
	Msg string
	// Prev is the PHP exception's previous exception (rendered too).
	Prev error
}

func (e *LoadError) Error() string { return e.Msg }

// Previous returns the chained exception.
func (e *LoadError) Previous() error { return e.Prev }

// Loader is ArrayLoader; LoadOptions controls transport-options loading.
type Loader struct {
	LoadOptions bool
}

func jsonForMessage(v phpjson.Value) string {
	s, err := phpjson.Encode(v, 0)
	if err != nil {
		b, _ := json.Marshal(v)
		return string(b)
	}
	return s
}

// Load is ArrayLoader::load for CompletePackage.
func (l *Loader) Load(config *phpjson.Map) (*Package, error) {
	p, err := l.createObject(config)
	if err != nil {
		return nil, err
	}
	for _, kind := range LinkKinds {
		v, ok := config.Get(kind.Key)
		if !ok || !phpjson.IsArray(v) {
			continue
		}
		links, err := ParseLinks(p.Name, p.PrettyVersion(), kind.Type, v)
		if err != nil {
			return nil, err
		}
		p.SetLinksOf(kind.Type, links)
	}
	return l.configureObject(p, config)
}

// LoadPackages is ArrayLoader::loadPackages as used by repositories: links to
// the package itself are skipped (configureCachedLinks).
func (l *Loader) LoadPackages(versions []*phpjson.Map) ([]*Package, error) {
	var out []*Package
	for _, config := range versions {
		p, err := l.createObject(config)
		if err != nil {
			return nil, err
		}
		for _, kind := range LinkKinds {
			v, ok := config.Get(kind.Key)
			if !ok || v == nil {
				continue
			}
			var links Links
			var lerr error
			phpjson.Each(v, func(prettyTarget string, c phpjson.Value) {
				if lerr != nil {
					return
				}
				target := strings.ToLower(prettyTarget)
				if target == p.Name {
					return
				}
				cs, ok := c.(string)
				if !ok {
					lerr = &LoadError{Msg: "Link constraint in " + p.Name + " " + kind.Type + " > " + target + " should be a string"}
					return
				}
				link, err := CreateLink(p.Name, p.PrettyVersion(), kind.Type, target, cs)
				if err != nil {
					lerr = err
					return
				}
				links = links.set(link)
			})
			if lerr != nil {
				return nil, lerr
			}
			p.SetLinksOf(kind.Type, links)
		}
		pp, err := l.configureObject(p, config)
		if err != nil {
			return nil, err
		}
		out = append(out, pp)
	}
	return out, nil
}

func (l *Loader) createObject(config *phpjson.Map) (*Package, error) {
	nameV, ok := config.Get("name")
	if !ok || nameV == nil {
		return nil, &LoadError{Msg: "Unknown package has no name defined (" + jsonForMessage(config) + ")."}
	}
	name, _ := phpjson.ScalarString(nameV)
	verV, ok := config.Get("version")
	if !ok || verV == nil || phpjson.IsArray(verV) {
		return nil, &LoadError{Msg: "Package " + name + " has no version defined."}
	}
	prettyVersion, _ := phpjson.ScalarString(verV)
	var version string
	if vn, ok := config.Get("version_normalized"); ok {
		if s, isStr := vn.(string); isStr {
			version = s
			if version == semver.DefaultBranchAlias {
				n, err := semver.Normalize(prettyVersion)
				if err != nil {
					return nil, err
				}
				version = n
			}
		}
	}
	if version == "" {
		n, err := semver.Normalize(prettyVersion)
		if err != nil {
			return nil, &LoadError{Msg: `Failed to normalize version for package "` + name + `": ` + err.Error(), Prev: err}
		}
		version = n
	}
	return NewPackage(name, version, prettyVersion), nil
}

func toStringList(v phpjson.Value) []string {
	var out []string
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		s, _ := phpjson.ScalarString(e)
		out = append(out, s)
	})
	return out
}

func toValueList(v phpjson.Value) []phpjson.Value {
	var out []phpjson.Value
	phpjson.Each(v, func(_ string, e phpjson.Value) { out = append(out, e) })
	return out
}

func (l *Loader) configureObject(p *Package, config *phpjson.Map) (*Package, error) {
	if t, ok := config.Get("type"); ok && t != nil {
		s, _ := phpjson.ScalarString(t)
		p.typ = strings.ToLower(s)
	} else {
		p.typ = "library"
	}
	if config.Isset("target-dir") {
		p.targetDir = config.Str("target-dir")
	}
	if e, ok := config.Get("extra"); ok && phpjson.IsArray(e) {
		p.extra = phpjson.ToMap(e)
	}
	if b, ok := config.Get("bin"); ok && b != nil {
		var bins []string
		if phpjson.IsArray(b) {
			bins = toStringList(b)
		} else {
			s, _ := phpjson.ScalarString(b)
			bins = []string{s}
		}
		for i := range bins {
			bins[i] = strings.TrimLeft(bins[i], "/")
		}
		p.binaries = bins
	}
	if config.Isset("installation-source") {
		p.installationSource = config.Str("installation-source")
	}
	if v, ok := config.Get("default-branch"); ok && v == true {
		p.isDefaultBranch = true
	}
	if src := config.Map("source"); config.Isset("source") {
		if src == nil || !src.Isset("type") || !src.Isset("url") || !src.Isset("reference") {
			v, _ := config.Get("source")
			return nil, &LoadError{Msg: fmt.Sprintf("Package %s's source key should be specified as {\"type\": ..., \"url\": ..., \"reference\": ...},\n%s given.", configName(config), jsonForMessage(v))}
		}
		p.sourceType, _ = phpjson.ScalarString(mustGet(src, "type"))
		p.sourceURL, _ = phpjson.ScalarString(mustGet(src, "url"))
		p.sourceReference, _ = phpjson.ScalarString(mustGet(src, "reference"))
		if m, ok := src.Get("mirrors"); ok && m != nil {
			p.sourceMirrors = toValueList(m)
		}
	}
	if dist := config.Map("dist"); config.Isset("dist") {
		if dist == nil || !dist.Isset("type") || !dist.Isset("url") {
			v, _ := config.Get("dist")
			return nil, &LoadError{Msg: fmt.Sprintf("Package %s's dist key should be specified as {\"type\": ..., \"url\": ..., \"reference\": ..., \"shasum\": ...},\n%s given.", configName(config), jsonForMessage(v))}
		}
		p.distType, _ = phpjson.ScalarString(mustGet(dist, "type"))
		p.distURL, _ = phpjson.ScalarString(mustGet(dist, "url"))
		if dist.Isset("reference") {
			p.distReference, _ = phpjson.ScalarString(mustGet(dist, "reference"))
		}
		if dist.Isset("shasum") {
			p.distSha1, _ = phpjson.ScalarString(mustGet(dist, "shasum"))
			p.distSha1Present = true
		}
		if m, ok := dist.Get("mirrors"); ok && m != nil {
			p.distMirrors = toValueList(m)
		}
	}
	if s, ok := config.Get("suggest"); ok && phpjson.IsArray(s) {
		sm := phpjson.ToMap(s).Clone()
		for _, k := range sm.Keys() {
			v, _ := sm.Get(k)
			if str, ok := v.(string); ok && strings.TrimSpace(str) == "self.version" {
				sm.Set(k, p.prettyVersion)
			}
		}
		p.suggests = sm
	}
	if a, ok := config.Get("autoload"); ok && a != nil {
		p.autoload = phpjson.ToMap(a)
	}
	if a, ok := config.Get("autoload-dev"); ok && a != nil {
		p.devAutoload = phpjson.ToMap(a)
	}
	if ip, ok := config.Get("include-path"); ok && ip != nil {
		p.includePaths = toStringList(ip)
	}
	if pe, ok := config.Get("php-ext"); ok && pe != nil {
		p.phpExt = pe
	}
	if t, ok := config.Get("time"); ok && phpjson.Truthy(t) {
		switch tv := t.(type) {
		case string:
			p.releaseDate = tv
		case int64:
			p.releaseDate = strconv.FormatInt(tv, 10)
		}
	}
	if t, ok := config.Get("published-time"); ok {
		switch tv := t.(type) {
		case string:
			if tv != "" {
				p.publishedDate = tv
			}
		case int64:
			p.publishedDate = strconv.FormatInt(tv, 10)
		}
	}
	if n, ok := config.Get("notification-url"); ok && phpjson.Truthy(n) {
		p.notificationURL, _ = phpjson.ScalarString(n)
	}
	if ar := config.Map("archive"); ar != nil {
		if n, ok := ar.Get("name"); ok && phpjson.Truthy(n) {
			p.archiveName, _ = phpjson.ScalarString(n)
		}
		if ex, ok := ar.Get("exclude"); ok && phpjson.Truthy(ex) {
			p.archiveExcludes = toStringList(ex)
		}
	}
	if sc, ok := config.Get("scripts"); ok && phpjson.IsArray(sc) {
		scripts := phpjson.NewMap()
		phpjson.Each(sc, func(event string, listeners phpjson.Value) {
			if phpjson.IsArray(listeners) {
				scripts.Set(event, toValueList(listeners))
			} else if listeners == nil {
				scripts.Set(event, []phpjson.Value{})
			} else {
				scripts.Set(event, []phpjson.Value{listeners})
			}
		})
		p.scripts = scripts
	}
	if d, ok := config.Get("description"); ok && phpjson.Truthy(d) {
		if s, ok := d.(string); ok {
			p.description = s
		}
	}
	if h, ok := config.Get("homepage"); ok && phpjson.Truthy(h) {
		if s, ok := h.(string); ok {
			p.homepage = s
		}
	}
	if k, ok := config.Get("keywords"); ok && phpjson.Truthy(k) && phpjson.IsArray(k) {
		p.keywords = toStringList(k)
	}
	if lic, ok := config.Get("license"); ok && phpjson.Truthy(lic) {
		if phpjson.IsArray(lic) {
			p.license = toStringList(lic)
		} else {
			s, _ := phpjson.ScalarString(lic)
			p.license = []string{s}
		}
	}
	if a, ok := config.Get("authors"); ok && phpjson.Truthy(a) && phpjson.IsArray(a) {
		p.authors = toValueList(a)
	}
	if s, ok := config.Get("support"); ok && phpjson.IsArray(s) {
		p.support = phpjson.ToMap(s)
	}
	if f, ok := config.Get("funding"); ok && phpjson.Truthy(f) && phpjson.IsArray(f) {
		p.funding = toValueList(f)
	}
	if ab, ok := config.Get("abandoned"); ok && ab != nil {
		p.abandoned = ab
	}
	if r, ok := config.Get("repositories"); ok && phpjson.IsArray(r) {
		p.repositories = toValueList(r)
	}
	if l.LoadOptions {
		if to := config.Map("transport-options"); to != nil {
			p.transportOptions = to
		}
	}
	alias, err := GetBranchAlias(config)
	if err != nil {
		return nil, err
	}
	if alias != "" {
		pretty := reNines.ReplaceAllString(alias, ".x")
		return NewAlias(p, alias, pretty), nil
	}
	return p, nil
}

var reNines = regexp.MustCompile(`(\.9{7})+`)

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

func configName(config *phpjson.Map) string {
	v, _ := config.Get("name")
	s, _ := phpjson.ScalarString(v)
	return s
}

// ParseLinks is ArrayLoader::parseLinks.
func ParseLinks(source, sourceVersion, linkType string, links phpjson.Value) (Links, error) {
	var res Links
	var err error
	phpjson.Each(links, func(target string, c phpjson.Value) {
		if err != nil {
			return
		}
		constraint, ok := c.(string)
		if !ok {
			return
		}
		var link *Link
		link, err = CreateLink(source, sourceVersion, linkType, strings.ToLower(target), constraint)
		if err == nil {
			res = res.set(link)
		}
	})
	return res, err
}

// CreateLink is ArrayLoader::createLink.
func CreateLink(source, sourceVersion, linkType, target, prettyConstraint string) (*Link, error) {
	constraint := prettyConstraint
	if prettyConstraint == "self.version" {
		constraint = sourceVersion
	}
	parsed, err := semver.ParseConstraints(constraint)
	if err != nil {
		return nil, &LoadError{Msg: "Link constraint in " + source + " " + linkType + " > " + target + ` should be a valid version constraint, got "` + constraint + `"`, Prev: err}
	}
	return &Link{Source: source, Target: target, Constraint: parsed, Description: LinkDescription(linkType), PrettyConstraint: prettyConstraint}, nil
}

// GetBranchAlias is ArrayLoader::getBranchAlias.
func GetBranchAlias(config *phpjson.Map) (string, error) {
	v, ok := config.Get("version")
	if !ok || v == nil || phpjson.IsArray(v) {
		return "", &LoadError{Msg: "no/invalid version defined"}
	}
	version, _ := phpjson.ScalarString(v)
	if !strings.HasPrefix(version, "dev-") && !strings.HasSuffix(version, "-dev") {
		return "", nil
	}
	if extra := config.Map("extra"); extra != nil {
		if ba, ok := extra.Get("branch-alias"); ok && phpjson.IsArray(ba) {
			result := ""
			phpjson.Each(ba, func(sourceBranch string, t phpjson.Value) {
				if result != "" {
					return
				}
				targetBranch, _ := phpjson.ScalarString(t)
				if !strings.HasSuffix(targetBranch, "-dev") {
					return
				}
				var validated string
				if targetBranch == semver.DefaultBranchAlias {
					validated = semver.DefaultBranchAlias
				} else {
					validated = semver.NormalizeBranch(targetBranch[:len(targetBranch)-4])
				}
				if !strings.HasSuffix(validated, "-dev") {
					return
				}
				if strings.ToLower(version) != strings.ToLower(sourceBranch) {
					return
				}
				if sp, ok := semver.ParseNumericAliasPrefix(sourceBranch); ok {
					if tp, ok := semver.ParseNumericAliasPrefix(targetBranch); ok && !strings.HasPrefix(strings.ToLower(tp), strings.ToLower(sp)) {
						return
					}
				}
				result = validated
			})
			if result != "" {
				return result, nil
			}
		}
	}
	if db, ok := config.Get("default-branch"); ok && db == true {
		if _, ok := semver.ParseNumericAliasPrefix(strings.TrimPrefix(version, "v")); !ok {
			return semver.DefaultBranchAlias, nil
		}
	}
	return "", nil
}

// ReleaseTime parses the `time` value like new \DateTime($s, UTC) and returns
// it formatted as DATE_RFC3339, or "" when PHP would fail to parse it.
func ReleaseTime(raw string) string {
	if raw == "" {
		return ""
	}
	t, ok := ParsePHPDate(raw)
	if !ok {
		return ""
	}
	return t.Format("2006-01-02T15:04:05-07:00")
}

// ParsePHPDate covers the date formats found in package metadata.
func ParsePHPDate(raw string) (time.Time, bool) {
	s := strings.TrimSpace(raw)
	allDigits := s != ""
	for _, c := range s {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return time.Time{}, false
		}
		return time.Unix(n, 0).UTC(), true
	}
	layoutsTZ := []string{
		time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05-0700", "2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05 -0700", "2006-01-02 15:04:05 -07:00", "2006-01-02T15:04:05.000000-07:00",
		"2006-01-02 15:04:05Z07:00", "Mon, 02 Jan 2006 15:04:05 -0700",
	}
	for _, layout := range layoutsTZ {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	layoutsUTC := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02 15:04", "2006-01-02", "2006-01-02 15:04:05.000000", "2006/01/02 15:04:05", "2006/01/02"}
	for _, layout := range layoutsUTC {
		if t, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			return t, true
		}
	}
	if strings.HasSuffix(s, " UTC") {
		return ParsePHPDate(strings.TrimSuffix(s, " UTC") + "+00:00")
	}
	return time.Time{}, false
}
