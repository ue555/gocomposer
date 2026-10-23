// Package validate ports Composer\Package\Loader\ValidatingArrayLoader and
// Composer\Util\ConfigValidator from composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/spdx"
)

// Flags of ValidatingArrayLoader.
const (
	CheckUnboundConstraints = 1
	CheckStrictConstraints  = 2
	CheckAll                = 3
)

// InvalidPackageError is InvalidPackageException.
type InvalidPackageError struct {
	Errors, Warnings []string
}

func (e *InvalidPackageError) Error() string {
	return "Invalid package information: \n" + strings.Join(append(append([]string(nil), e.Errors...), e.Warnings...), "\n")
}

// SecurityError is Composer\Exception\SecurityException (via
// ValidatingArrayLoader::validatePackage).
type SecurityError struct{ Msg string }

func (e *SecurityError) Error() string { return e.Msg }

// ArrayLoader validates a package config like ValidatingArrayLoader.
type ArrayLoader struct {
	Flags    int
	errors   []string
	warnings []string
	config   *phpjson.Map
	// Now is used for the license re-validation window (tests may set it).
	Now func() time.Time
}

// Warnings returns the warnings of the last Load.
func (l *ArrayLoader) Warnings() []string { return l.warnings }

// DebugType is PHP get_debug_type() for decoded JSON values.
func DebugType(v phpjson.Value) string {
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
	case []phpjson.Value, *phpjson.Map:
		return "array"
	}
	return fmt.Sprintf("%T", v)
}

// toStr is PHP's string conversion of a scalar.
func toStr(v phpjson.Value) string {
	s, _ := phpjson.ScalarString(v)
	return s
}

func isScalar(v phpjson.Value) bool {
	switch v.(type) {
	case bool, int64, float64, string:
		return true
	}
	return false
}

func isNumeric(v phpjson.Value) bool {
	switch t := v.(type) {
	case int64, float64:
		return true
	case string:
		s := strings.TrimLeft(t, " \t\n\r\v\f")
		if s == "" {
			return false
		}
		_, err := strconv.ParseFloat(strings.TrimRight(s, " \t\n\r\v\f"), 64)
		return err == nil && !strings.ContainsAny(s, "xXpP_") && !strings.EqualFold(strings.TrimLeft(s, "+-"), "inf") && !strings.EqualFold(strings.TrimLeft(s, "+-"), "infinity") && !strings.EqualFold(s, "nan")
	}
	return false
}

func phpTrim(s string) string { return strings.Trim(s, " \t\n\r\x00\x0B") }

// get returns config[key] with PHP isset semantics.
func (l *ArrayLoader) isset(key string) bool { return l.config.Isset(key) }

func (l *ArrayLoader) get(key string) phpjson.Value { v, _ := l.config.Get(key); return v }

func (l *ArrayLoader) errorf(format string, a ...any) {
	l.errors = append(l.errors, fmt.Sprintf(format, a...))
}

func (l *ArrayLoader) warnf(format string, a ...any) {
	l.warnings = append(l.warnings, fmt.Sprintf(format, a...))
}

func (l *ArrayLoader) validateString(property string, mandatory bool) bool {
	if l.isset(property) {
		if _, ok := l.get(property).(string); !ok {
			l.errors = append(l.errors, property+" : should be a string, "+DebugType(l.get(property))+" given")
			l.config.Delete(property)
			return false
		}
	}
	if !l.isset(property) || phpTrim(l.get(property).(string)) == "" {
		if mandatory {
			l.errors = append(l.errors, property+" : must be present")
		}
		l.config.Delete(property)
		return false
	}
	return true
}

func (l *ArrayLoader) validateArray(property string, mandatory bool) bool {
	if l.isset(property) && !phpjson.IsArray(l.get(property)) {
		l.errors = append(l.errors, property+" : should be an array, "+DebugType(l.get(property))+" given")
		l.config.Delete(property)
		return false
	}
	if !l.isset(property) || phpjson.ArrayLen(l.get(property)) == 0 {
		if mandatory {
			l.errors = append(l.errors, property+" : must be present and contain at least one element")
		}
		l.config.Delete(property)
		return false
	}
	return true
}

func (l *ArrayLoader) validateFlatArray(property string, re *regexp.Regexp, reSrc string, mandatory bool) bool {
	if !l.validateArray(property, mandatory) {
		return false
	}
	pass := true
	arr := l.get(property)
	var drop []string
	phpjson.Each(arr, func(k string, v phpjson.Value) {
		if _, isStr := v.(string); !isStr && !isNumeric(v) {
			l.errors = append(l.errors, property+"."+k+" : must be a string or int, "+DebugType(v)+" given")
			drop = append(drop, k)
			pass = false
			return
		}
		if re != nil && !re.MatchString(toStr(v)) {
			l.warnings = append(l.warnings, property+"."+k+" : invalid value ("+toStr(v)+"), must match "+reSrc)
			drop = append(drop, k)
			pass = false
		}
	})
	l.config.Set(property, unsetKeys(arr, drop))
	return pass
}

// unsetKeys removes keys from a decoded array (lists keep their indexes
// as an object to mirror PHP's sparse arrays).
func unsetKeys(v phpjson.Value, keys []string) phpjson.Value {
	if len(keys) == 0 {
		return v
	}
	drop := map[string]bool{}
	for _, k := range keys {
		drop[k] = true
	}
	out := phpjson.NewMap()
	phpjson.Each(v, func(k string, val phpjson.Value) {
		if !drop[k] {
			out.Set(k, val)
		}
	})
	return out
}

func (l *ArrayLoader) validateRegex(property, re string, mandatory bool) bool {
	if !l.validateString(property, mandatory) {
		return false
	}
	s := l.get(property).(string)
	if !regexp.MustCompile(`^(?:` + re + `)\n?$`).MatchString(s) {
		msg := property + " : invalid value (" + s + "), must match " + re
		if mandatory {
			l.errors = append(l.errors, msg)
		} else {
			l.warnings = append(l.warnings, msg)
		}
		l.config.Delete(property)
		return false
	}
	return true
}

func (l *ArrayLoader) validateURL(property string, mandatory bool) bool {
	if !l.validateString(property, mandatory) {
		return false
	}
	s := l.get(property).(string)
	if !FilterURL(s, []string{"http", "https"}) {
		l.warnings = append(l.warnings, property+" : invalid value ("+s+"), must be an http/https URL")
		l.config.Delete(property)
		return false
	}
	return true
}

var reSchemeHost = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9+.-]*):`)

// FilterURL is ValidatingArrayLoader::filterUrl; scheme and host are
// extracted like PHP's parse_url (which, unlike net/url, accepts any
// characters in the host).
func FilterURL(value string, schemes []string) bool {
	if value == "" {
		return true
	}
	scheme, host, ok := parseURLSchemeHost(value)
	if !ok || scheme == "" || host == "" {
		return false
	}
	for _, s := range schemes {
		if scheme == s {
			return true
		}
	}
	return false
}

func parseURLSchemeHost(v string) (string, string, bool) {
	m := reSchemeHost.FindStringSubmatch(v)
	if m == nil {
		return "", "", true
	}
	rest := v[len(m[0]):]
	if !strings.HasPrefix(rest, "//") {
		return m[1], "", true
	}
	auth := rest[2:]
	if i := strings.IndexAny(auth, "/?#"); i >= 0 {
		auth = auth[:i]
	}
	if i := strings.LastIndexByte(auth, '@'); i >= 0 {
		auth = auth[i+1:]
	}
	host := auth
	if strings.HasPrefix(host, "[") {
		if i := strings.IndexByte(host, ']'); i >= 0 {
			port := host[i+1:]
			host = host[:i+1]
			if port != "" && !validPort(strings.TrimPrefix(port, ":")) {
				return "", "", false
			}
		}
	} else if i := strings.LastIndexByte(host, ':'); i >= 0 {
		if !validPort(host[i+1:]) {
			return "", "", false
		}
		host = host[:i]
	}
	return m[1], host, true
}

func validPort(p string) bool {
	if p == "" {
		return true
	}
	n, err := strconv.Atoi(p)
	return err == nil && n >= 0 && n <= 65535
}

var (
	reType         = `[A-Za-z0-9-]+`
	reKeywords     = `[\p{N}\p{L} ._-]+`
	reKeywordsRE   = regexp.MustCompile(`^(?:` + reKeywords + `)\n?$`)
	rePkgKey       = regexp.MustCompile(`^[A-Za-z0-9_./-]+$`)
	reDashStart    = regexp.MustCompile(`^\s*-`)
	reBinDotDot    = regexp.MustCompile(`(?:^|[\\/])\.\.(?:[\\/]|$)`)
	reUnsafeBin    = regexp.MustCompile("[*$`\"&^|<>()%!;\x00-\x1f\x7f]")
	rePerforceRsh  = regexp.MustCompile(`(?i)^\s*(?:rsh|jsh)\s*:`)
	rePerforcePort = regexp.MustCompile(`(?i)^(?:(?:tcp|ssl)(?:4|6|46|64)?:)?(?:\[[0-9a-f:.]+\]|[a-z0-9._][a-z0-9._-]*)(?::[a-z0-9._][a-z0-9._-]*)?$`)
	reGPLPlus      = regexp.MustCompile(`(?i)^[AL]?GPL-[123](\.[01])?\+$`)
	reGPLBare      = regexp.MustCompile(`(?i)^[AL]?GPL-[123](\.[01])?$`)
)

// IsValidPerforcePort is Perforce::isValidPort.
func IsValidPerforcePort(u string) bool {
	if rePerforceRsh.MatchString(u) {
		return false
	}
	return rePerforcePort.MatchString(u)
}

func findBinError(bin string) string {
	if reBinDotDot.MatchString(bin) {
		return `must not contain a ".." path component`
	}
	if reUnsafeBin.MatchString(bin) {
		return "must not contain any of the characters *$`\"&^|<>()%!; nor control characters"
	}
	return ""
}

// Load is ValidatingArrayLoader::load. It returns the loaded package, or an
// *InvalidPackageError listing the errors (warnings are in Warnings()).
func (l *ArrayLoader) Load(config *phpjson.Map) (*pkg.Package, error) {
	l.errors, l.warnings = nil, nil
	orig := config
	l.config = phpjson.AsMap(phpjson.CloneValue(config))
	cfg := l.config

	l.validateString("name", true)
	if n, ok := orig.Get("name"); ok && n != nil {
		if err := repo.PackageNamingError(toStr(n), false); err != "" {
			l.errors = append(l.errors, "name : "+err)
		}
	}
	if l.isset("version") {
		v := l.get("version")
		if !isScalar(v) {
			l.validateString("version", false)
		} else {
			s := toStr(v)
			cfg.Set("version", s)
			if _, err := semver.Normalize(s); err != nil {
				l.errors = append(l.errors, "version : invalid value ("+s+"): "+err.Error())
				cfg.Delete("version")
			}
		}
	}
	if c := cfg.Map("config"); c != nil && c.Isset("platform") {
		pv, _ := c.Get("platform")
		eachCast(pv, func(key string, platform phpjson.Value) {
			if platform == false {
				return
			}
			s, ok := platform.(string)
			if !ok {
				l.errors = append(l.errors, "config.platform."+key+" : invalid value ("+DebugType(platform)+" "+varExport(platform)+"): expected string or false")
				return
			}
			if _, err := semver.Normalize(s); err != nil {
				l.errors = append(l.errors, "config.platform."+key+" : invalid value ("+s+"): "+err.Error())
			}
		})
	}

	l.validateRegex("type", reType, false)
	l.validateString("target-dir", false)
	l.validateArray("extra", false)

	if l.isset("bin") {
		if _, ok := l.get("bin").(string); ok {
			l.validateString("bin", false)
		} else {
			l.validateFlatArray("bin", nil, "", false)
		}
		if s, ok := l.get("bin").(string); ok && l.isset("bin") {
			if e := findBinError(s); e != "" {
				l.errors = append(l.errors, "bin : invalid value ("+s+"), "+e)
				cfg.Delete("bin")
			}
		} else if l.isset("bin") && phpjson.IsArray(l.get("bin")) {
			var drop []string
			phpjson.Each(l.get("bin"), func(k string, b phpjson.Value) {
				if s, ok := b.(string); ok {
					if e := findBinError(s); e != "" {
						l.errors = append(l.errors, "bin."+k+" : invalid value ("+s+"), "+e)
						drop = append(drop, k)
					}
				}
			})
			cfg.Set("bin", unsetKeys(l.get("bin"), drop))
		}
	}

	l.validateArray("scripts", false)
	l.validateString("description", false)
	l.validateURL("homepage", false)
	l.validateFlatArray("keywords", reKeywordsRE, reKeywords, false)

	var releaseDate *time.Time
	l.validateString("time", false)
	if l.isset("time") {
		s := l.get("time").(string)
		if t, ok := pkg.ParsePHPDate(s); ok {
			releaseDate = &t
		} else {
			l.errors = append(l.errors, "time : invalid value ("+s+"): "+dateError(s))
			cfg.Delete("time")
		}
	}

	if l.isset("license") {
		lic := l.get("license")
		_, isStr := lic.(string)
		if phpjson.IsArray(lic) || isStr {
			var licenses []string
			var kept []phpjson.Value
			if isStr {
				kept = []phpjson.Value{lic}
			} else {
				phpjson.Each(lic, func(_ string, v phpjson.Value) { kept = append(kept, v) })
			}
			var final []phpjson.Value
			for _, v := range kept {
				s, ok := v.(string)
				if !ok {
					l.warnings = append(l.warnings, "License "+jsonEncode(v)+" should be a string.")
					continue
				}
				licenses = append(licenses, s)
				final = append(final, s)
			}
			now := time.Now
			if l.Now != nil {
				now = l.Now
			}
			if releaseDate == nil || !releaseDate.Before(now().Add(-8*24*time.Hour)) {
				for _, license := range licenses {
					if license == "proprietary" {
						continue
					}
					toValidate := strings.ReplaceAll(license, "proprietary", "MIT")
					if !spdx.IsValidLicenseString(toValidate) {
						if spdx.IsValidLicenseString(phpTrim(toValidate)) {
							l.warnings = append(l.warnings, "License "+jsonEncode(license)+" must not contain extra spaces, make sure to trim it.")
						} else {
							l.warnings = append(l.warnings, "License "+jsonEncode(license)+" is not a valid SPDX license identifier, see https://spdx.org/licenses/ if you use an open license.\nIf the software is closed-source, you may use \"proprietary\" as license.")
						}
					}
				}
			}
			if final == nil {
				final = []phpjson.Value{}
			}
			cfg.Set("license", final)
		} else {
			l.warnings = append(l.warnings, "License must be a string or array of strings, got "+jsonEncode(lic)+".")
			cfg.Delete("license")
		}
	}

	if l.validateArray("authors", false) {
		authors := phpjson.NewMap()
		phpjson.Each(l.get("authors"), func(key string, author phpjson.Value) {
			if !phpjson.IsArray(author) {
				l.errors = append(l.errors, "authors."+key+" : should be an array, "+DebugType(author)+" given")
				return
			}
			am := phpjson.ToMap(author)
			if am == nil {
				am = phpjson.NewMap()
			}
			am = am.Clone()
			for _, f := range []string{"homepage", "email", "name", "role"} {
				if am.Isset(f) {
					if _, ok := mustGet(am, f).(string); !ok {
						l.errors = append(l.errors, "authors."+key+"."+f+" : invalid value, must be a string")
						am.Delete(f)
					}
				}
			}
			orig := phpjson.ToMap(author)
			if orig != nil && orig.Isset("homepage") {
				hp := toStr(mustGet(orig, "homepage"))
				if !FilterURL(hp, []string{"http", "https"}) {
					l.warnings = append(l.warnings, "authors."+key+".homepage : invalid value ("+hp+"), must be an http/https URL")
					am.Delete("homepage")
				}
			}
			if orig != nil && orig.Isset("email") {
				em := toStr(mustGet(orig, "email"))
				if !ValidEmail(em) {
					l.warnings = append(l.warnings, "authors."+key+".email : invalid value ("+em+"), must be a valid email address")
					am.Delete("email")
				}
			}
			if am.Len() > 0 {
				authors.Set(key, am)
			}
		})
		if authors.Len() == 0 {
			cfg.Delete("authors")
		} else {
			cfg.Set("authors", authors)
		}
	}

	if l.validateArray("support", false) && phpjson.ArrayLen(l.get("support")) > 0 {
		sup := phpjson.ToMap(l.get("support"))
		if sup == nil {
			sup = phpjson.NewMap()
		}
		sup = sup.Clone()
		for _, k := range []string{"issues", "forum", "wiki", "source", "email", "irc", "docs", "rss", "chat", "security"} {
			if sup.Isset(k) {
				if _, ok := mustGet(sup, k).(string); !ok {
					l.errors = append(l.errors, "support."+k+" : invalid value, must be a string")
					sup.Delete(k)
				}
			}
		}
		if sup.Isset("email") && !ValidEmail(toStr(mustGet(sup, "email"))) {
			l.warnings = append(l.warnings, "support.email : invalid value ("+toStr(mustGet(sup, "email"))+"), must be a valid email address")
			sup.Delete("email")
		}
		if sup.Isset("irc") && !FilterURL(toStr(mustGet(sup, "irc")), []string{"irc", "ircs"}) {
			l.warnings = append(l.warnings, "support.irc : invalid value ("+toStr(mustGet(sup, "irc"))+"), must be a irc://<server>/<channel> or ircs:// URL")
			sup.Delete("irc")
		}
		for _, k := range []string{"issues", "forum", "wiki", "source", "docs", "chat", "security"} {
			if sup.Isset(k) && !FilterURL(toStr(mustGet(sup, k)), []string{"http", "https"}) {
				l.warnings = append(l.warnings, "support."+k+" : invalid value ("+toStr(mustGet(sup, k))+"), must be an http/https URL")
				sup.Delete(k)
			}
		}
		if sup.Len() == 0 {
			cfg.Delete("support")
		} else {
			cfg.Set("support", sup)
		}
	}

	if l.validateArray("funding", false) && phpjson.ArrayLen(l.get("funding")) > 0 {
		funding := phpjson.NewMap()
		phpjson.Each(l.get("funding"), func(key string, opt phpjson.Value) {
			if !phpjson.IsArray(opt) {
				l.errors = append(l.errors, "funding."+key+" : should be an array, "+DebugType(opt)+" given")
				return
			}
			om := phpjson.ToMap(opt)
			if om == nil {
				om = phpjson.NewMap()
			}
			orig := om
			om = om.Clone()
			for _, f := range []string{"type", "url"} {
				if om.Isset(f) {
					if _, ok := mustGet(om, f).(string); !ok {
						l.errors = append(l.errors, "funding."+key+"."+f+" : invalid value, must be a string")
						om.Delete(f)
					}
				}
			}
			if orig.Isset("url") && !FilterURL(toStr(mustGet(orig, "url")), []string{"http", "https"}) {
				l.warnings = append(l.warnings, "funding."+key+".url : invalid value ("+toStr(mustGet(orig, "url"))+"), must be an http/https URL")
				om.Delete("url")
			}
			if om.Len() > 0 {
				funding.Set(key, om)
			}
		})
		if funding.Len() == 0 {
			cfg.Delete("funding")
		} else {
			cfg.Set("funding", funding)
		}
	}

	if l.isset("php-ext") && l.validateArray("php-ext", false) {
		l.validatePhpExt()
	}

	unbound := semver.NewConstraint("=", "10000000-dev")
	for _, lk := range pkg.LinkKinds {
		linkType := lk.Key
		if !(l.validateArray(linkType, false) && l.isset(linkType)) {
			continue
		}
		links := phpjson.ToMap(l.get(linkType))
		if links == nil {
			links = phpjson.NewMap()
		}
		links = links.Clone()
		phpjson.Each(l.get(linkType), func(pkgName string, constraint phpjson.Value) {
			if l.isset("name") && strings.EqualFold(pkgName, toStr(l.get("name"))) {
				l.errors = append(l.errors, linkType+"."+pkgName+" : a package cannot set a "+linkType+" on itself")
				links.Delete(pkgName)
				return
			}
			if err := repo.PackageNamingError(pkgName, true); err != "" {
				l.warnings = append(l.warnings, linkType+"."+err)
			} else if !rePkgKey.MatchString(pkgName) {
				l.errors = append(l.errors, linkType+"."+pkgName+" : invalid key, package names must be strings containing only [A-Za-z0-9_./-]")
			}
			cs, isStr := constraint.(string)
			if !isStr {
				l.errors = append(l.errors, linkType+"."+pkgName+" : invalid value, must be a string containing a version constraint")
				links.Delete(pkgName)
			} else if cs != "self.version" {
				lc, err := semver.ParseConstraints(cs)
				if err != nil {
					l.errors = append(l.errors, linkType+"."+pkgName+" : invalid version constraint ("+err.Error()+")")
					links.Delete(pkgName)
					return
				}
				if l.Flags&CheckUnboundConstraints != 0 && linkType == "require" && lc.Matches(unbound) && !repo.IsPlatformPackage(pkgName) {
					l.warnings = append(l.warnings, linkType+"."+pkgName+" : unbound version constraints ("+cs+") should be avoided")
				} else if single, ok := lc.(*semver.Single); ok && l.Flags&CheckStrictConstraints != 0 && linkType == "require" && (single.Operator() == "==" || single.Operator() == "=") && semver.NewConstraint(">=", "1.0.0.0-dev").Matches(single) {
					l.warnings = append(l.warnings, linkType+"."+pkgName+" : exact version constraints ("+cs+") should be avoided if the package follows semantic versioning")
				}
				if _, none := semver.CompactConstraint(lc).(*semver.MatchNone); none {
					l.warnings = append(l.warnings, linkType+"."+pkgName+" : this version constraint cannot possibly match anything ("+cs+")")
				}
			}
			if linkType == "conflict" && l.isset("replace") {
				replace := phpjson.ToMap(l.get("replace"))
				if replace != nil {
					overlap := false
					phpjson.Each(l.get("conflict"), func(ck string, _ phpjson.Value) {
						if replace.Has(ck) {
							overlap = true
						}
					})
					if overlap {
						l.errors = append(l.errors, linkType+"."+pkgName+" : you cannot conflict with a package that is also replaced, as replace already creates an implicit conflict rule")
						links.Delete(pkgName)
					}
				}
			}
		})
		cfg.Set(linkType, links)
	}

	if l.validateArray("suggest", false) && l.isset("suggest") {
		sug := phpjson.ToMap(l.get("suggest"))
		phpjson.Each(l.get("suggest"), func(p string, d phpjson.Value) {
			if _, ok := d.(string); !ok {
				l.errors = append(l.errors, "suggest."+p+" : invalid value, must be a string describing why the package is suggested")
				if sug != nil {
					sug.Delete(p)
				}
			}
		})
	}

	if l.validateString("minimum-stability", false) && l.isset("minimum-stability") {
		ms := l.get("minimum-stability").(string)
		// isset(STABILITIES[strtolower($ms)]): "rc" is rejected, "RC" allowed
		if _, ok := pkg.Stabilities[strings.ToLower(ms)]; !ok && ms != "RC" {
			l.errors = append(l.errors, "minimum-stability : invalid value ("+ms+"), must be one of "+strings.Join(pkg.StabilityNames, ", "))
			cfg.Delete("minimum-stability")
		}
	}

	if l.validateArray("autoload", false) && l.isset("autoload") {
		types := []string{"psr-0", "psr-4", "classmap", "files", "exclude-from-classmap"}
		al := phpjson.ToMap(l.get("autoload"))
		phpjson.Each(l.get("autoload"), func(typ string, tc phpjson.Value) {
			if !inList(typ, types) {
				l.errors = append(l.errors, "autoload : invalid value ("+typ+"), must be one of "+strings.Join(types, ", "))
				if al != nil {
					al.Delete(typ)
				}
			}
			if typ == "psr-4" {
				phpjson.Each(tc, func(ns string, _ phpjson.Value) {
					if ns != "" && !strings.HasSuffix(ns, `\`) {
						l.errors = append(l.errors, "autoload.psr-4 : invalid value ("+ns+"), namespaces must end with a namespace separator, should be "+ns+`\\`)
					}
				})
			}
		})
	}

	if al := cfg.Map("autoload"); al != nil && al.Isset("psr-4") && l.isset("target-dir") {
		l.errors = append(l.errors, "target-dir : this can not be used together with the autoload.psr-4 setting, remove target-dir to upgrade to psr-4")
		al.Delete("psr-4")
	}

	for _, srcType := range []string{"source", "dist"} {
		if !(l.validateArray(srcType, false) && phpjson.ArrayLen(l.get(srcType)) > 0) {
			continue
		}
		sm := phpjson.ToMap(l.get(srcType))
		if sm == nil {
			sm = phpjson.NewMap()
		}
		if !sm.Isset("type") {
			l.errors = append(l.errors, srcType+".type : must be present")
		}
		if !sm.Isset("url") {
			l.errors = append(l.errors, srcType+".url : must be present")
		}
		if srcType == "source" && !sm.Isset("reference") {
			l.errors = append(l.errors, srcType+".reference : must be present")
		}
		if sm.Isset("type") {
			if _, ok := mustGet(sm, "type").(string); !ok {
				l.errors = append(l.errors, srcType+".type : should be a string, "+DebugType(mustGet(sm, "type"))+" given")
			}
		}
		if sm.Isset("url") {
			if _, ok := mustGet(sm, "url").(string); !ok {
				l.errors = append(l.errors, srcType+".url : should be a string, "+DebugType(mustGet(sm, "url"))+" given")
			}
		}
		if sm.Isset("reference") {
			ref := mustGet(sm, "reference")
			_, isS := ref.(string)
			_, isI := ref.(int64)
			if !isS && !isI {
				l.errors = append(l.errors, srcType+".reference : should be a string or int, "+DebugType(ref)+" given")
			}
			if reDashStart.MatchString(toStr(ref)) {
				l.errors = append(l.errors, srcType+`.reference : must not start with a "-", "`+toStr(ref)+`" given`)
			}
		}
		if sm.Isset("url") && reDashStart.MatchString(toStr(mustGet(sm, "url"))) {
			l.errors = append(l.errors, srcType+`.url : must not start with a "-", "`+toStr(mustGet(sm, "url"))+`" given`)
		}
		if srcType == "source" && toStr(mustGet(sm, "type")) == "perforce" {
			if u, ok := mustGet(sm, "url").(string); ok && !IsValidPerforcePort(u) {
				l.errors = append(l.errors, srcType+`.url : invalid Perforce port ("`+u+`"), it must be of the form [tcp|ssl:][host:]port`)
			}
		}
	}

	l.validateFlatArray("include-path", nil, "", false)
	l.validateArray("transport-options", false)

	if extra := cfg.Map("extra"); extra != nil && extra.Isset("branch-alias") {
		ba, _ := extra.Get("branch-alias")
		if !phpjson.IsArray(ba) {
			l.errors = append(l.errors, "extra.branch-alias : must be an array of versions => aliases")
		} else {
			phpjson.Each(ba, func(source string, target phpjson.Value) {
				t, ok := target.(string)
				if !ok {
					l.warnings = append(l.warnings, "extra.branch-alias."+source+" : the target branch ("+jsonEncode(target)+") must be a string, \""+DebugType(target)+"\" received.")
					return
				}
				if !strings.HasSuffix(t, "-dev") {
					l.warnings = append(l.warnings, "extra.branch-alias."+source+" : the target branch ("+t+") must end in -dev")
					return
				}
				validated := semver.NormalizeBranch(t[:len(t)-4])
				if !strings.HasSuffix(validated, "-dev") {
					l.warnings = append(l.warnings, "extra.branch-alias."+source+" : the target branch ("+t+") must be a parseable number like 2.0-dev")
					return
				}
				sp, sok := semver.ParseNumericAliasPrefix(source)
				tp, tok := semver.ParseNumericAliasPrefix(t)
				if sok && sp != "" && tok && tp != "" && !strings.HasPrefix(strings.ToLower(tp), strings.ToLower(sp)) {
					l.warnings = append(l.warnings, "extra.branch-alias."+source+" : the target branch ("+t+") is not a valid numeric alias for this version")
				}
			})
		}
	}

	if len(l.errors) > 0 {
		return nil, &InvalidPackageError{Errors: l.errors, Warnings: l.warnings}
	}
	loader := &pkg.Loader{}
	p, err := loader.Load(cfg)
	l.config = nil
	return p, err
}

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

func inList(s string, list []string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// eachCast iterates (array) $value.
func eachCast(v phpjson.Value, fn func(k string, v phpjson.Value)) {
	if phpjson.IsArray(v) {
		phpjson.Each(v, fn)
		return
	}
	fn("0", v)
}

func jsonEncode(v phpjson.Value) string {
	s, _ := phpjson.Encode(v, 0)
	return s
}

func varExport(v phpjson.Value) string {
	switch t := v.(type) {
	case nil:
		return "NULL"
	case bool:
		if t {
			return "true"
		}
		return "false"
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return phpjson.FormatFloatVarExport(t)
	case string:
		return "'" + strings.ReplaceAll(strings.ReplaceAll(t, `\`, `\\`), `'`, `\'`) + "'"
	}
	var sb strings.Builder
	sb.WriteString("array (\n")
	phpjson.Each(v, func(k string, val phpjson.Value) {
		key := k
		if _, isInt := phpjson.IntKey(k); !isInt {
			key = "'" + strings.ReplaceAll(strings.ReplaceAll(k, `\`, `\\`), `'`, `\'`) + "'"
		}
		inner := strings.ReplaceAll(varExport(val), "\n", "\n  ")
		sb.WriteString("  " + key + " => " + inner + ",\n")
	})
	sb.WriteString(")")
	return sb.String()
}

// dateError approximates the DateMalformedStringException message of
// `new DateTime($time)` (PHP reports the first parse error).
func dateError(s string) string {
	c := ""
	if len(s) > 0 {
		c = s[:1]
	}
	return "Failed to parse time string (" + s + ") at position 0 (" + c + "): The timezone could not be found in the database"
}

func (l *ArrayLoader) validatePhpExt() {
	cfg := l.config
	if t := toStr(l.get("type")); t != "php-ext" && t != "php-ext-zend" {
		l.errors = append(l.errors, `php-ext can only be set by packages of type "php-ext" or "php-ext-zend" which must be C extensions`)
		cfg.Delete("php-ext")
	}
	phpExt := phpjson.ToMap(l.get("php-ext"))
	if phpExt == nil {
		return
	}
	if phpExt.Isset("extension-name") {
		if _, ok := mustGet(phpExt, "extension-name").(string); !ok {
			l.errors = append(l.errors, "php-ext.extension-name : should be a string, "+DebugType(mustGet(phpExt, "extension-name"))+" given")
			phpExt.Delete("extension-name")
		}
	}
	if phpExt.Isset("priority") {
		if _, ok := mustGet(phpExt, "priority").(int64); !ok {
			l.errors = append(l.errors, "php-ext.priority : should be an integer, "+DebugType(mustGet(phpExt, "priority"))+" given")
			phpExt.Delete("priority")
		}
	}
	for _, f := range []string{"support-zts", "support-nts"} {
		if phpExt.Isset(f) {
			if _, ok := mustGet(phpExt, f).(bool); !ok {
				l.errors = append(l.errors, "php-ext."+f+" : should be a boolean, "+DebugType(mustGet(phpExt, f))+" given")
				phpExt.Delete(f)
			}
		}
	}
	if phpExt.Isset("build-path") {
		if _, ok := mustGet(phpExt, "build-path").(string); !ok {
			l.errors = append(l.errors, "php-ext.build-path : should be a string or null, "+DebugType(mustGet(phpExt, "build-path"))+" given")
			phpExt.Delete("build-path")
		}
	}
	if phpExt.Isset("download-url-method") {
		m := mustGet(phpExt, "download-url-method")
		_, isStr := m.(string)
		if !phpjson.IsArray(m) && !isStr {
			l.errors = append(l.errors, "php-ext.download-url-method : should be an array or a string, "+DebugType(m)+" given")
			phpExt.Delete("download-url-method")
		} else {
			valid := []string{"composer-default", "pre-packaged-source", "pre-packaged-binary"}
			defined := m
			if isStr {
				defined = []phpjson.Value{m}
			}
			if phpjson.ArrayLen(defined) == 0 {
				l.errors = append(l.errors, "php-ext.download-url-method : must contain at least one element")
				phpExt.Delete("download-url-method")
			} else {
				phpjson.Each(defined, func(k string, v phpjson.Value) {
					s, ok := v.(string)
					if !ok {
						l.errors = append(l.errors, "php-ext.download-url-method."+k+" : should be a string, "+DebugType(v)+" given")
						phpExt.Delete("download-url-method")
					} else if !inList(s, valid) {
						l.errors = append(l.errors, "php-ext.download-url-method."+k+" : invalid value ("+s+"), must be one of "+strings.Join(valid, ", "))
						phpExt.Delete("download-url-method")
					}
				})
			}
		}
	}
	if phpExt.Isset("os-families") && phpExt.Isset("os-families-exclude") {
		l.errors = append(l.errors, "php-ext : os-families and os-families-exclude cannot both be specified")
		phpExt.Delete("os-families")
		phpExt.Delete("os-families-exclude")
	} else {
		valid := []string{"windows", "bsd", "darwin", "solaris", "linux", "unknown"}
		for _, field := range []string{"os-families", "os-families-exclude"} {
			if !phpExt.Isset(field) {
				continue
			}
			v := mustGet(phpExt, field)
			if !phpjson.IsArray(v) {
				l.errors = append(l.errors, "php-ext."+field+" : should be an array, "+DebugType(v)+" given")
				phpExt.Delete(field)
			} else if phpjson.ArrayLen(v) == 0 {
				l.errors = append(l.errors, "php-ext."+field+" : must contain at least one element")
				phpExt.Delete(field)
			} else {
				var drop []string
				phpjson.Each(v, func(k string, os phpjson.Value) {
					s, ok := os.(string)
					if !ok {
						l.errors = append(l.errors, "php-ext."+field+"."+k+" : should be a string, "+DebugType(os)+" given")
						drop = append(drop, k)
					} else if !inList(s, valid) {
						l.errors = append(l.errors, "php-ext."+field+"."+k+" : invalid value ("+s+"), must be one of "+strings.Join(valid, ", "))
						drop = append(drop, k)
					}
				})
				nv := unsetKeys(v, drop)
				if phpjson.ArrayLen(nv) == 0 {
					phpExt.Delete(field)
				} else {
					phpExt.Set(field, nv)
				}
			}
		}
	}
	if phpExt.Isset("configure-options") {
		co := mustGet(phpExt, "configure-options")
		if !phpjson.IsArray(co) {
			l.errors = append(l.errors, "php-ext.configure-options : should be an array, "+DebugType(co)+" given")
			phpExt.Delete("configure-options")
		} else {
			var drop []string
			phpjson.Each(co, func(k string, opt phpjson.Value) {
				if !phpjson.IsArray(opt) {
					l.errors = append(l.errors, "php-ext.configure-options."+k+" : should be an array, "+DebugType(opt)+" given")
					drop = append(drop, k)
					return
				}
				om := phpjson.ToMap(opt)
				if om == nil || !om.Isset("name") {
					l.errors = append(l.errors, "php-ext.configure-options."+k+".name : must be present")
					drop = append(drop, k)
					return
				}
				if _, ok := mustGet(om, "name").(string); !ok {
					l.errors = append(l.errors, "php-ext.configure-options."+k+".name : should be a string, "+DebugType(mustGet(om, "name"))+" given")
					drop = append(drop, k)
					return
				}
				if om.Isset("needs-value") {
					if _, ok := mustGet(om, "needs-value").(bool); !ok {
						l.errors = append(l.errors, "php-ext.configure-options."+k+".needs-value : should be a boolean, "+DebugType(mustGet(om, "needs-value"))+" given")
					}
				}
				if om.Isset("description") {
					if _, ok := mustGet(om, "description").(string); !ok {
						l.errors = append(l.errors, "php-ext.configure-options."+k+".description : should be a string, "+DebugType(mustGet(om, "description"))+" given")
					}
				}
			})
			nv := unsetKeys(co, drop)
			if phpjson.ArrayLen(nv) == 0 {
				phpExt.Delete("configure-options")
			}
		}
	}
	if phpExt.Len() == 0 {
		cfg.Delete("php-ext")
	}
}

// ValidatePackage is ValidatingArrayLoader::validatePackage, applied by
// LockTransaction to every package of a solver result.
func ValidatePackage(p *pkg.Package) error {
	if p.IsRoot {
		return nil
	}
	if err := repo.PackageNamingError(p.Name, false); err != "" {
		return &SecurityError{"Invalid package found during dependency resolution, aborting: " + err}
	}
	for _, f := range []struct{ field, value string }{
		{"source.url", p.SourceURL()}, {"source.reference", p.SourceReference()},
		{"dist.url", p.DistURL()}, {"dist.reference", p.DistReference()},
	} {
		if f.value != "" && reDashStart.MatchString(f.value) {
			return &SecurityError{p.Name + " has an invalid " + f.field + `, it must not start with a "-": ` + f.value}
		}
	}
	if p.SourceType() == "perforce" && p.SourceURL() != "" && !IsValidPerforcePort(p.SourceURL()) {
		return &SecurityError{p.Name + " has an invalid source.url, it must be a Perforce port of the form [tcp|ssl:][host:]port: " + p.SourceURL()}
	}
	for _, bin := range p.Binaries() {
		if reBinDotDot.MatchString(bin) {
			return &SecurityError{p.Name + " has an invalid bin " + bin + `, it must not contain ".." path segments`}
		}
	}
	return nil
}
