// Package platform ports Composer\Repository\PlatformRepository. Facts about
// the PHP runtime (constants, loaded extensions, ReflectionExtension::info()
// output) are collected by running a small probe with the configured PHP
// binary; turning them into platform packages happens here in Go.
package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
)

// Versions reported for the virtual composer packages. They match the
// reference Composer checkout (Composer::getVersion() of a source checkout,
// PluginInterface::PLUGIN_API_VERSION and Composer::RUNTIME_API_VERSION).
const (
	ComposerVersion     = "2.11.999-dev+source"
	PluginAPIVersion    = "2.11.0"
	RuntimeAPIVersion   = "2.2.2"
	ComposerReleaseDate = "@release_date@"
)

// Facts is what the probe reports.
type Facts struct {
	PHPVersion  string            `json:"php_version"`
	PHPBinary   string            `json:"php_binary"`
	Debug       bool              `json:"debug"`
	ZTS         bool              `json:"zts"`
	IntSize     int               `json:"int_size"`
	IPv6        bool              `json:"ipv6"`
	Extensions  []string          `json:"extensions"`
	ExtVersions map[string]string `json:"ext_versions"`
	ExtInfo     map[string]string `json:"ext_info"`
	Constants   map[string]any    `json:"constants"`
	CurlVersion string            `json:"curl_version"`
	ICUCLDR     string            `json:"icu_cldr"`
	ICUUnicode  string            `json:"icu_unicode"`
	Imagick     string            `json:"imagick"`
	IniFiles    []string          `json:"ini_files"`
	OSFamily    string            `json:"os_family"`
	PHPOS       string            `json:"php_os"`
}

// probeScript prints runtime facts as JSON; it does not implement any Composer logic.
const probeScript = `
$c = array();
foreach (array('GD_VERSION','GMP_VERSION','ICONV_VERSION','INTL_ICU_VERSION','LIBXML_DOTTED_VERSION','MB_ONIGURUMA_VERSION','OPENSSL_VERSION_TEXT','PCRE_VERSION','PGSQL_LIBPQ_VERSION','RD_KAFKA_VERSION','SODIUM_LIBRARY_VERSION','LIBXSLT_DOTTED_VERSION','ZLIB_VERSION','AF_INET6','ZipArchive::LIBZIP_VERSION') as $k) {
  if (PHP_VERSION_ID >= 80600 && $k === 'MB_ONIGURUMA_VERSION') continue;
  if (defined($k)) { $c[$k] = constant($k); }
}
$exts = get_loaded_extensions();
$ver = array(); $info = array();
foreach ($exts as $e) {
  $v = phpversion($e); $ver[$e] = $v === false ? '0' : (string) $v;
  $r = new ReflectionExtension($e); ob_start(); $r->info(); $info[$e] = (string) ob_get_clean();
}
$curl = function_exists('curl_version') ? curl_version() : null;
$cldr = null;
if (class_exists('ResourceBundle', false)) { $rb = ResourceBundle::create('root', 'ICUDATA', false); if ($rb !== null) { $cldr = $rb->get('Version'); } }
$uni = class_exists('IntlChar', false) ? implode('.', array_slice(IntlChar::getUnicodeVersion(), 0, 3)) : null;
$im = null;
if (class_exists('Imagick', false)) { $iv = (new Imagick())->getVersion(); $im = $iv['versionString']; }
$inis = array(); $loaded = php_ini_loaded_file(); $inis[] = $loaded === false ? '' : $loaded;
$scanned = php_ini_scanned_files();
if ($scanned !== false && $scanned !== '') { foreach (explode(',', $scanned) as $s) { $s = trim($s); if ($s !== '') $inis[] = $s; } }
$ipv6 = defined('AF_INET6') || @inet_pton('::') !== false;
echo json_encode(array('php_version' => PHP_VERSION, 'php_binary' => PHP_BINARY, 'debug' => (bool) PHP_DEBUG, 'zts' => defined('PHP_ZTS') && PHP_ZTS,
  'int_size' => PHP_INT_SIZE, 'ipv6' => $ipv6, 'extensions' => $exts, 'ext_versions' => $ver, 'ext_info' => $info, 'constants' => $c,
  'curl_version' => $curl ? $curl['version'] : null, 'icu_cldr' => $cldr, 'icu_unicode' => $uni, 'imagick' => $im, 'ini_files' => $inis,
  'os_family' => PHP_OS_FAMILY, 'php_os' => PHP_OS), JSON_PARTIAL_OUTPUT_ON_ERROR | JSON_INVALID_UTF8_SUBSTITUTE);
`

// ErrNoPHP is returned when no PHP binary could be found.
var ErrNoPHP = errors.New("no PHP binary found")

// FindPHP returns the PHP binary to use: GOCOMPOSER_PHP, then PHP_BINARY-like
// lookup of "php" on PATH.
func FindPHP() (string, error) {
	if p := os.Getenv("GOCOMPOSER_PHP"); p != "" {
		return p, nil
	}
	p, err := exec.LookPath("php")
	if err != nil {
		return "", ErrNoPHP
	}
	return p, nil
}

// ProbeScript returns the PHP code Probe runs (to record facts elsewhere).
func ProbeScript() string { return probeScript }

// Probe runs the probe script with the given PHP binary.
func Probe(php string) (*Facts, error) {
	cmd := exec.Command(php, "-d", "display_errors=stderr", "-d", "error_reporting=0", "-r", probeScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running %s failed: %v: %s", php, err, strings.TrimSpace(stderr.String()))
	}
	var f Facts
	if err := json.Unmarshal(stdout.Bytes(), &f); err != nil {
		return nil, fmt.Errorf("could not read PHP runtime information from %s: %v", php, err)
	}
	return &f, nil
}

// Override is one config.platform entry; Version "" with Disabled marks false.
type Override struct {
	Name     string
	Version  string
	Disabled bool
}

// Repository is the platform repository built from facts and overrides.
type Repository struct {
	*repo.Array
	facts     *Facts
	overrides map[string]Override
	order     []string
	disabled  map[string]*pkg.Package
	libraries map[string]bool
	lastPHP   string
}

// RepoName is "platform repo".
func (r *Repository) RepoName() string { return "platform repo" }

// Facts returns the probed runtime facts (nil when PHP was not available).
func (r *Repository) Facts() *Facts { return r.facts }

// Disabled reports packages disabled via config.platform.
func (r *Repository) Disabled() map[string]*pkg.Package { return r.disabled }

// IsPlatformPackageDisabled is PlatformRepository::isPlatformPackageDisabled.
func (r *Repository) IsPlatformPackageDisabled(name string) bool {
	_, ok := r.disabled[name]
	return ok
}

// LastSeenPlatformPHP is PlatformRepository::getPlatformPhpVersion.
func (r *Repository) LastSeenPlatformPHP() string { return r.lastPHP }

// OverrideError mirrors the UnexpectedValueException raised for bad overrides.
type OverrideError struct{ Msg string }

func (e *OverrideError) Error() string { return e.Msg }

// New builds the platform repository. facts may be nil when no PHP is
// available, in which case only overrides and the composer packages exist.
func New(facts *Facts, overrides []Override) (*Repository, error) {
	r := &Repository{Array: repo.NewPlatformArray(), facts: facts, overrides: map[string]Override{}, disabled: map[string]*pkg.Package{}, libraries: map[string]bool{}}
	r.Array.SetSelf(r)
	for _, o := range overrides {
		if o.Name == "php" && o.Disabled {
			return nil, &OverrideError{"config.platform." + o.Name + " cannot be set to false as you cannot disable php entirely."}
		}
		key := strings.ToLower(o.Name)
		if _, ok := r.overrides[key]; !ok {
			r.order = append(r.order, key)
		}
		r.overrides[key] = o
	}
	if err := r.initialize(); err != nil {
		return nil, err
	}
	return r, nil
}

func newPlatformPackage(name, version, pretty, description string) *pkg.Package {
	p := pkg.NewPackage(name, version, pretty)
	p.SetDescription(description)
	return p
}

func (r *Repository) initialize() error {
	for _, key := range r.order {
		o := r.overrides[key]
		if !repo.IsPlatformPackage(o.Name) {
			return &OverrideError{"Invalid platform package name in config.platform: " + o.Name}
		}
		if !o.Disabled {
			if _, err := r.addOverridden(o, ""); err != nil {
				return err
			}
		}
	}
	add := func(name, pretty, desc string) {
		v := semver.MustNormalize(pretty)
		r.addPackage(newPlatformPackage(name, v, pretty, desc))
	}
	add("composer", ComposerVersion, "Composer package")
	add("composer-plugin-api", PluginAPIVersion, "The Composer Plugin API")
	add("composer-runtime-api", RuntimeAPIVersion, "The Composer Runtime API")
	f := r.facts
	if f == nil {
		return nil
	}
	pretty := f.PHPVersion
	version, err := semver.Normalize(pretty)
	if err != nil {
		pretty = regexp.MustCompile(`^([^~+-]+).*$`).ReplaceAllString(f.PHPVersion, "$1")
		version, err = semver.Normalize(pretty)
		if err != nil {
			return err
		}
	}
	r.addPackage(newPlatformPackage("php", version, pretty, "The PHP interpreter"))
	if f.Debug {
		r.addPackage(newPlatformPackage("php-debug", version, pretty, "The PHP interpreter, with debugging symbols"))
	}
	if f.ZTS {
		r.addPackage(newPlatformPackage("php-zts", version, pretty, "The PHP interpreter, with Zend Thread Safety"))
	}
	if f.IntSize == 8 {
		r.addPackage(newPlatformPackage("php-64bit", version, pretty, "The PHP interpreter, 64bit"))
	}
	if f.IPv6 {
		r.addPackage(newPlatformPackage("php-ipv6", version, pretty, "The PHP interpreter, with IPv6 support"))
	}
	for _, name := range f.Extensions {
		if name == "standard" || name == "Core" {
			continue
		}
		r.addExtension(name, f.ExtVersions[name])
	}
	r.addLibraries()
	return nil
}

func (r *Repository) addPackage(p *pkg.Package) {
	if o, ok := r.overrides[p.Name]; ok {
		if o.Disabled {
			p.SetDescription(p.Description() + ". <warning>Package disabled via config.platform</warning>")
			r.markConfigPlatform(p)
			r.disabled[p.Name] = p
			return
		}
		overrider := repo.FindPackage(r.Array, p.Name, &semver.MatchAll{})
		actual := "actual: " + p.PrettyVersion()
		if p.Version() == overrider.Version() {
			actual = "same as actual"
		}
		overrider.SetDescription(overrider.Description() + ", " + actual)
		return
	}
	if o, ok := r.overrides["php"]; ok && strings.HasPrefix(p.Name, "php-") {
		overrider, err := r.addOverridden(o, p.PrettyName)
		if err != nil {
			return
		}
		actual := "actual: " + p.PrettyVersion()
		if p.Version() == overrider.Version() {
			actual = "same as actual"
		}
		overrider.SetDescription(overrider.Description() + ", " + actual)
		return
	}
	r.Array.AddPackage(p)
}

func (r *Repository) markConfigPlatform(p *pkg.Package) {
	extra := phpjson.NewMap()
	extra.Set("config.platform", true)
	p.SetExtra(extra)
}

func (r *Repository) addOverridden(o Override, name string) (*pkg.Package, error) {
	version, err := semver.Normalize(o.Version)
	if err != nil {
		return nil, err
	}
	n := o.Name
	if name != "" {
		n = name
	}
	p := newPlatformPackage(n, version, o.Version, "Package overridden via config.platform")
	r.markConfigPlatform(p)
	r.Array.AddPackage(p)
	if p.Name == "php" {
		parts := strings.Split(p.Version(), ".")
		if len(parts) > 3 {
			parts = parts[:3]
		}
		r.lastPHP = strings.Join(parts, ".")
	}
	return p, nil
}

var reExtVersionPrefix = regexp.MustCompile(`^(\d+\.\d+\.\d+(?:\.\d+)?)`)

func (r *Repository) addExtension(name, prettyVersion string) {
	extra := ""
	version, err := semver.Normalize(prettyVersion)
	if err != nil {
		extra = " (actual version: " + prettyVersion + ")"
		if m := reExtVersionPrefix.FindStringSubmatch(prettyVersion); m != nil {
			prettyVersion = m[1]
		} else {
			prettyVersion = "0"
		}
		version = semver.MustNormalize(prettyVersion)
	}
	pkgName := "ext-" + strings.ReplaceAll(strings.ToLower(name), " ", "-")
	ext := newPlatformPackage(pkgName, version, prettyVersion, "The "+name+" PHP extension"+extra)
	ext.SetType("php-ext")
	if name == "uuid" {
		c := semver.NewConstraint("=", version)
		ext.SetReplaces(pkg.Links{{Source: "ext-uuid", Target: "lib-uuid", Constraint: c, Description: pkg.TypeReplace, PrettyConstraint: ext.PrettyVersion()}})
	}
	r.addPackage(ext)
}

func (r *Repository) addLibrary(name, prettyVersion, description string, replaces, provides []string) {
	if prettyVersion == "" {
		return
	}
	version, err := semver.Normalize(prettyVersion)
	if err != nil {
		return
	}
	if !repo.IsPlatformPackage("lib-" + name) {
		return
	}
	if r.libraries["lib-"+name] {
		return
	}
	r.libraries["lib-"+name] = true
	if description == "" {
		description = "The " + name + " library"
	}
	lib := newPlatformPackage("lib-"+name, version, prettyVersion, description)
	var rl, pl pkg.Links
	for _, rp := range replaces {
		rp = strings.ToLower(rp)
		rl = append(rl, &pkg.Link{Source: "lib-" + name, Target: "lib-" + rp, Constraint: semver.NewConstraint("=", version), Description: pkg.TypeReplace, PrettyConstraint: lib.PrettyVersion()})
	}
	for _, pv := range provides {
		pv = strings.ToLower(pv)
		pl = append(pl, &pkg.Link{Source: "lib-" + name, Target: "lib-" + pv, Constraint: semver.NewConstraint("=", version), Description: pkg.TypeProvide, PrettyConstraint: lib.PrettyVersion()})
	}
	lib.SetReplaces(rl)
	lib.SetProvides(pl)
	r.addPackage(lib)
}

func (r *Repository) constant(name string) (string, bool) {
	v, ok := r.facts.Constants[name]
	if !ok || v == nil {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case float64:
		return strconv.FormatInt(int64(t), 10), true
	case bool:
		if t {
			return "1", true
		}
		return "", true
	}
	return fmt.Sprint(v), true
}

func match(pattern, s string) map[string]string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	out := map[string]string{}
	for i, n := range re.SubexpNames() {
		if n != "" {
			out[n] = m[i]
		}
	}
	return out
}

func (r *Repository) addLibraries() {
	f := r.facts
	loaded := map[string]bool{}
	for _, e := range f.Extensions {
		loaded[e] = true
	}
	for _, name := range f.Extensions {
		info := f.ExtInfo[name]
		switch name {
		case "amqp":
			if m := match(`(?im)^librabbitmq version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-librabbitmq", m["version"], "AMQP librabbitmq version", nil, nil)
			}
			if m := match(`(?im)^AMQP protocol version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-protocol", strings.ReplaceAll(m["version"], "-", "."), "AMQP protocol version", nil, nil)
			}
		case "bz2":
			if m := match(`(?im)^BZip2 Version => (?P<version>.*),`, info); m != nil {
				r.addLibrary(name, m["version"], "", nil, nil)
			}
		case "curl":
			r.addLibrary(name, f.CurlVersion, "", nil, nil)
			if m := match(`(?im)^SSL Version => (?P<library>[^\r\n/]+)/(?P<version>[^\r\n]+?)\r?$`, info); m != nil {
				library := strings.ToLower(m["library"])
				if library == "openssl" {
					parsed, fips := ParseOpenssl(m["version"])
					if parsed != "" {
						libName := name + "-openssl"
						var provides []string
						if fips {
							libName += "-fips"
							provides = []string{"curl-openssl"}
						}
						r.addLibrary(libName, parsed, "curl OpenSSL version ("+parsed+")", nil, provides)
					}
				} else {
					shortlib := library
					sslLib := "curl-openssl"
					if strings.HasPrefix(library, "(securetransport)") {
						if sm := regexp.MustCompile(`^\(securetransport\) ([a-z0-9]+)`).FindStringSubmatch(library); sm != nil {
							shortlib = "securetransport"
							sslLib = "curl-" + sm[1]
						}
					}
					r.addLibrary(name+"-"+shortlib, m["version"], "curl "+library+" version ("+m["version"]+")", []string{sslLib}, nil)
				}
			}
			if m := match(`(?im)^libSSH Version => (?P<library>[^\r\n/]+)/(?P<version>.+?)(?:/.*)?$`, info); m != nil {
				r.addLibrary(name+"-"+strings.ToLower(m["library"]), m["version"], "curl "+m["library"]+" version", nil, nil)
			}
			if m := match(`(?im)^ZLib Version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-zlib", m["version"], "curl zlib version", nil, nil)
			}
		case "date":
			if m := match(`(?im)^timelib version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-timelib", m["version"], "date timelib version", nil, nil)
			}
			if m := match(`(?im)^Timezone Database => (?P<source>internal|external)$`, info); m != nil {
				external := m["source"] == "external"
				if zm := match(`(?im)^"Olson" Timezone Database Version => (?P<version>.+?)(?:\.system)?$`, info); zm != nil {
					if external && loaded["timezonedb"] {
						r.addLibrary("timezonedb-zoneinfo", zm["version"], `zoneinfo ("Olson") database for date (replaced by timezonedb)`, []string{name + "-zoneinfo"}, nil)
					} else {
						r.addLibrary(name+"-zoneinfo", zm["version"], `zoneinfo ("Olson") database for date`, nil, nil)
					}
				}
			}
		case "fileinfo":
			if m := match(`(?im)^libmagic => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libmagic", m["version"], "fileinfo libmagic version", nil, nil)
			}
		case "gd":
			if v, ok := r.constant("GD_VERSION"); ok {
				r.addLibrary(name, v, "", nil, nil)
			}
			if m := match(`(?im)^libJPEG Version => (?P<version>.+?)(?: compatible)?$`, info); m != nil {
				if v := ParseLibjpeg(m["version"]); v != "" {
					r.addLibrary(name+"-libjpeg", v, "libjpeg version for gd", nil, nil)
				}
			}
			if m := match(`(?im)^libPNG Version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libpng", m["version"], "libpng version for gd", nil, nil)
			}
			if m := match(`(?im)^FreeType Version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-freetype", m["version"], "freetype version for gd", nil, nil)
			}
			if m := match(`(?im)^libXpm Version => (?P<versionId>\d+)$`, info); m != nil {
				id, _ := strconv.Atoi(m["versionId"])
				r.addLibrary(name+"-libxpm", convertVersionID(id, 100), "libxpm version for gd", nil, nil)
			}
		case "gmp":
			if v, ok := r.constant("GMP_VERSION"); ok {
				r.addLibrary(name, v, "", nil, nil)
			}
		case "iconv":
			if v, ok := r.constant("ICONV_VERSION"); ok {
				r.addLibrary(name, v, "", nil, nil)
			}
		case "intl":
			desc := "The ICU unicode and globalization support library"
			if v, ok := r.constant("INTL_ICU_VERSION"); ok {
				r.addLibrary("icu", v, desc, nil, nil)
			} else if m := match(`(?im)^ICU version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary("icu", m["version"], desc, nil, nil)
			}
			if m := match(`(?im)^ICU TZData version => (?P<version>.*)$`, info); m != nil {
				if v := ParseZoneinfoVersion(m["version"]); v != "" {
					r.addLibrary("icu-zoneinfo", v, `zoneinfo ("Olson") database for icu`, nil, nil)
				}
			}
			if f.ICUCLDR != "" {
				r.addLibrary("icu-cldr", f.ICUCLDR, "ICU CLDR project version", nil, nil)
			}
			if f.ICUUnicode != "" {
				r.addLibrary("icu-unicode", f.ICUUnicode, "ICU unicode version", nil, nil)
			}
		case "imagick":
			if m := match(`^ImageMagick (?P<version>[\d.]+)(?:-(?P<patch>\d+))?`, f.Imagick); m != nil {
				v := m["version"]
				if m["patch"] != "" {
					v += "." + m["patch"]
				}
				r.addLibrary(name+"-imagemagick", v, "", []string{"imagick"}, nil)
			}
		case "ldap":
			m := match(`(?im)^Vendor Version => (?P<versionId>\d+)$`, info)
			vm := match(`(?im)^Vendor Name => (?P<vendor>.+)$`, info)
			if m != nil && vm != nil {
				id, _ := strconv.Atoi(m["versionId"])
				r.addLibrary(name+"-"+strings.ToLower(vm["vendor"]), convertVersionID(id, 100), vm["vendor"]+" version of ldap", nil, nil)
			}
		case "libxml":
			var provides []string
			for _, e := range f.Extensions {
				switch e {
				case "dom", "simplexml", "xml", "xmlreader", "xmlwriter":
					provides = append(provides, e+"-libxml")
				}
			}
			if v, ok := r.constant("LIBXML_DOTTED_VERSION"); ok {
				r.addLibrary(name, v, "libxml library version", nil, provides)
			}
		case "mbstring":
			if m := match(`(?im)^libmbfl version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libmbfl", m["version"], "mbstring libmbfl version", nil, nil)
			}
			if v, ok := r.constant("MB_ONIGURUMA_VERSION"); ok {
				r.addLibrary(name+"-oniguruma", v, "mbstring oniguruma version", nil, nil)
			} else if m := match(`(?im)^(?:oniguruma|Multibyte regex \(oniguruma\)) version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-oniguruma", m["version"], "mbstring oniguruma version", nil, nil)
			}
		case "memcached":
			if m := match(`(?im)^libmemcached version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libmemcached", m["version"], "libmemcached version", nil, nil)
			}
		case "openssl":
			if text, ok := r.constant("OPENSSL_VERSION_TEXT"); ok {
				if m := match(`(?i)^(?:OpenSSL|LibreSSL)?\s*(?P<version>\S+)`, text); m != nil {
					if parsed, fips := ParseOpenssl(m["version"]); parsed != "" {
						libName := name
						var provides []string
						if fips {
							libName += "-fips"
							provides = []string{name}
						}
						r.addLibrary(libName, parsed, text, nil, provides)
					}
				}
			}
		case "pcre":
			if v, ok := r.constant("PCRE_VERSION"); ok {
				r.addLibrary(name, regexp.MustCompile(`^(\S+).*`).ReplaceAllString(v, "$1"), "", nil, nil)
			}
			if m := match(`(?im)^PCRE Unicode Version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-unicode", m["version"], "PCRE Unicode version support", nil, nil)
			}
		case "mysqlnd", "pdo_mysql":
			if m := match(`(?im)^(?:Client API version|Version) => mysqlnd (?P<version>.+?) `, info); m != nil {
				r.addLibrary(name+"-mysqlnd", m["version"], "mysqlnd library version for "+name, nil, nil)
			}
		case "mongodb":
			if m := match(`(?im)^libmongoc bundled version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libmongoc", m["version"], "libmongoc version of mongodb", nil, nil)
			}
			if m := match(`(?im)^libbson bundled version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libbson", m["version"], "libbson version of mongodb", nil, nil)
			}
		case "pgsql", "pdo_pgsql":
			if name == "pgsql" {
				if v, ok := r.constant("PGSQL_LIBPQ_VERSION"); ok {
					r.addLibrary("pgsql-libpq", v, "libpq for pgsql", nil, nil)
					break
				}
			}
			if m := match(`(?im)^PostgreSQL\(libpq\) Version => (?P<version>.*)$`, info); m != nil {
				r.addLibrary(name+"-libpq", m["version"], "libpq for "+name, nil, nil)
			}
		case "pq":
			if m := match(`(?im)^libpq => (?P<compiled>.+) => (?P<linked>.+)$`, info); m != nil {
				r.addLibrary(name+"-libpq", m["linked"], "libpq for "+name, nil, nil)
			}
		case "rdkafka":
			if v, ok := r.constant("RD_KAFKA_VERSION"); ok {
				n, _ := strconv.ParseInt(v, 10, 64)
				r.addLibrary(name+"-librdkafka", fmt.Sprintf("%d.%d.%d", (n&0x7F000000)>>24, (n&0x00FF0000)>>16, (n&0x0000FF00)>>8), "librdkafka for "+name, nil, nil)
			}
		case "libsodium", "sodium":
			if v, ok := r.constant("SODIUM_LIBRARY_VERSION"); ok {
				r.addLibrary("libsodium", v, "", nil, nil)
			}
		case "sqlite3", "pdo_sqlite":
			if m := match(`(?im)^SQLite Library => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-sqlite", m["version"], "", nil, nil)
			}
		case "ssh2":
			if m := match(`(?im)^libssh2 version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name+"-libssh2", m["version"], "", nil, nil)
			}
		case "xsl":
			if v, ok := r.constant("LIBXSLT_DOTTED_VERSION"); ok {
				r.addLibrary("libxslt", v, "", []string{"xsl"}, nil)
			}
			if m := match(`(?im)^libxslt compiled against libxml Version => (?P<version>.+)$`, f.ExtInfo["xsl"]); m != nil {
				r.addLibrary("libxslt-libxml", m["version"], "libxml version libxslt is compiled against", nil, nil)
			}
		case "yaml":
			if m := match(`(?im)^LibYAML Version => (?P<version>.+)$`, f.ExtInfo["yaml"]); m != nil {
				r.addLibrary(name+"-libyaml", m["version"], "libyaml version of yaml", nil, nil)
			}
		case "zip":
			if v, ok := r.constant("ZipArchive::LIBZIP_VERSION"); ok {
				r.addLibrary(name+"-libzip", v, "", []string{"zip"}, nil)
			}
		case "zlib":
			if v, ok := r.constant("ZLIB_VERSION"); ok {
				r.addLibrary(name, v, "", nil, nil)
			} else if m := match(`(?im)^Linked Version => (?P<version>.+)$`, info); m != nil {
				r.addLibrary(name, m["version"], "", nil, nil)
			}
		}
	}
}

// ParseOpenssl is Composer\Platform\Version::parseOpenssl.
func ParseOpenssl(v string) (string, bool) {
	m := match(`^(?P<version>[0-9.]+)(?P<patch>[a-z]{0,2})(?P<suffix>(?:-?(?:dev|pre|alpha|beta|rc|fips)[\d]*)*)(?:-\w+)?(?: \(.+?\))?$`, v)
	if m == nil {
		return "", false
	}
	patch := ""
	if semver.VersionCompare(m["version"], "3.0.0") < 0 {
		patch = "." + strconv.Itoa(convertAlpha(m["patch"]))
	}
	fips := strings.Contains(m["suffix"], "fips")
	suffix := "-" + strings.TrimLeft(m["suffix"], "-")
	suffix = strings.NewReplacer("-fips", "", "-pre", "-alpha").Replace(suffix)
	return strings.TrimRight(m["version"]+patch+suffix, "-"), fips
}

// ParseLibjpeg is Version::parseLibjpeg.
func ParseLibjpeg(v string) string {
	m := match(`^(?P<major>\d+)(?P<minor>[a-z]*)$`, v)
	if m == nil {
		return ""
	}
	return m["major"] + "." + strconv.Itoa(convertAlpha(m["minor"]))
}

// ParseZoneinfoVersion is Version::parseZoneinfoVersion.
func ParseZoneinfoVersion(v string) string {
	m := match(`^(?P<year>\d{4})(?P<revision>[a-z]*)$`, v)
	if m == nil {
		return ""
	}
	return m["year"] + "." + strconv.Itoa(convertAlpha(m["revision"]))
}

func convertAlpha(alpha string) int {
	sum := 0
	for _, c := range alpha {
		sum += int(c)
	}
	return len(alpha)*(-int('a')+1) + sum
}

func convertVersionID(id, base int) string {
	return fmt.Sprintf("%d.%d.%d", id/(base*base), (id/base)%base, id%base)
}

// Names returns sorted names of the platform packages (for diagnostics).
func (r *Repository) Names() []string {
	var out []string
	for _, p := range r.MustPackages() {
		out = append(out, p.Name)
	}
	sort.Strings(out)
	return out
}
