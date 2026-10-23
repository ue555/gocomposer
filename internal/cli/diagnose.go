package cli

// Port of Composer\Command\DiagnoseCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
//
// gocomposer differences, reported in Composer's line format:
//   - "Composer version" is the Composer version gocomposer is compatible with.
//   - PHP/OpenSSL/curl/zip/platform facts come from the PHP binary gocomposer
//     uses for scripts and platform detection (probed, no Composer logic).
//   - The self-audit of Composer's own dependencies cannot run (there is no
//     Composer installation), so it reports Composer's "non-standard
//     installation" warning when the network is enabled.
//   - The pubkey/version checks only run for the Composer phar and are
//     skipped like for a source checkout.

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/platform"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/validate"
)

func diagnoseCommand() *Command {
	return &Command{Name: "diagnose", Run: runDiagnose}
}

// diagnoseProbe prints facts about the PHP runtime as JSON; it does not
// implement any Composer logic.
const diagnoseProbe = `
$curl = null;
if (extension_loaded('curl')) {
  $v = curl_version();
  $curl = array(
    'version' => $v['version'],
    'libz_version' => isset($v['libz_version']) ? $v['libz_version'] : '',
    'brotli_version' => isset($v['brotli_version']) ? $v['brotli_version'] : '',
    'ssl_version' => isset($v['ssl_version']) ? $v['ssl_version'] : '',
    'zstd' => isset($v['features']) && defined('CURL_VERSION_ZSTD') && 0 !== ($v['features'] & CURL_VERSION_ZSTD),
    'http2' => isset($v['features']) && defined('CURL_VERSION_HTTP2') && defined('CURL_HTTP_VERSION_2_0') && (CURL_VERSION_HTTP2 & $v['features']) !== 0,
    'http3' => isset($v['features']) && defined('CURL_VERSION_HTTP3') && ($v['features'] & CURL_VERSION_HTTP3) !== 0,
    'multi' => function_exists('curl_multi_exec') && function_exists('curl_multi_init'),
  );
}
ob_start(); phpinfo(INFO_GENERAL); $phpinfo = ob_get_clean();
$configure = null;
if (is_string($phpinfo) && preg_match('{Configure Command(?: *</td><td class="v">| *=> *)(.*?)(?:</td>|$)}m', $phpinfo, $m)) { $configure = $m[1]; }
$ext = array();
foreach (array('zip','Phar','filter','hash','iconv','mbstring','openssl','apcu','zlib','xdebug','uopz','ionCube Loader') as $e) { $ext[$e] = extension_loaded($e); }
$ini = array();
foreach (array('allow_url_fopen','apc.enable_cli','xdebug.profiler_enabled','uopz.disable','uopz.exit') as $k) { $v = ini_get($k); $ini[$k] = $v === false ? null : $v; }
$scanned = php_ini_scanned_files();
echo json_encode(array(
  'binary' => PHP_BINARY,
  'version' => PHP_VERSION,
  'version_id' => PHP_VERSION_ID,
  'openssl_text' => defined('OPENSSL_VERSION_TEXT') ? OPENSSL_VERSION_TEXT : null,
  'openssl_number' => defined('OPENSSL_VERSION_NUMBER') ? OPENSSL_VERSION_NUMBER : null,
  'curl' => $curl,
  'ext' => $ext,
  'ini' => $ini,
  'json_decode' => function_exists('json_decode'),
  'proc_open' => function_exists('proc_open'),
  'configure' => $configure,
  'xdebug_active' => extension_loaded('xdebug') && (!function_exists('xdebug_info') || count(xdebug_info('mode')) > 0),
  'ioncube' => extension_loaded('ionCube Loader') ? array(ioncube_loader_iversion(), ioncube_loader_version()) : null,
  'windows_build' => defined('PHP_WINDOWS_VERSION_BUILD'),
  'ini_files' => array_merge(array((string) php_ini_loaded_file()), $scanned !== false ? array_map('trim', explode(',', $scanned)) : array()),
));
`

type diagnoseFacts struct {
	Binary        string             `json:"binary"`
	Version       string             `json:"version"`
	VersionID     int                `json:"version_id"`
	OpenSSLText   *string            `json:"openssl_text"`
	OpenSSLNumber *int64             `json:"openssl_number"`
	Curl          *diagnoseCurl      `json:"curl"`
	Ext           map[string]bool    `json:"ext"`
	Ini           map[string]*string `json:"ini"`
	JSONDecode    bool               `json:"json_decode"`
	ProcOpen      bool               `json:"proc_open"`
	Configure     *string            `json:"configure"`
	XdebugActive  bool               `json:"xdebug_active"`
	IonCube       []any              `json:"ioncube"`
	WindowsBuild  bool               `json:"windows_build"`
	IniFiles      []string           `json:"ini_files"`
}

type diagnoseCurl struct {
	Version       string `json:"version"`
	LibzVersion   string `json:"libz_version"`
	BrotliVersion string `json:"brotli_version"`
	SSLVersion    string `json:"ssl_version"`
	Zstd          bool   `json:"zstd"`
	HTTP2         bool   `json:"http2"`
	HTTP3         bool   `json:"http3"`
	Multi         bool   `json:"multi"`
}

func probeDiagnose() (*diagnoseFacts, error) {
	php, err := platform.FindPHP()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(php, "-d", "display_errors=stderr", "-d", "error_reporting=0", "-r", diagnoseProbe)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var f diagnoseFacts
	if err := json.Unmarshal(out.Bytes(), &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// diagResult is the $result of DiagnoseCommand::outputResult: true, a
// message, a list of messages or an exception.
type diagResult struct {
	ok   bool
	msgs []string
	err  error
}

var diagOK = diagResult{ok: true}

func diagMsg(m ...string) diagResult { return diagResult{msgs: m} }

type diagnoser struct {
	ctx      *Context
	io       *cio.IO
	cfg      *config.Config
	exitCode int
	facts    *diagnoseFacts
	dl       *httpdl.Downloader
}

func runDiagnose(ctx *Context) (int, error) {
	io := ctx.IO
	c := ctx.cache
	var cfg *config.Config
	if c != nil {
		cfg = c.Config
	} else {
		cwd, _ := os.Getwd()
		cc, err := config.CreateConfig(cwd, func(m string) { io.WriteError(m) })
		if err != nil {
			return 1, err
		}
		cfg = cc
	}
	merge := phpjson.NewMap()
	mc := phpjson.NewMap()
	mc.Set("secure-http", false)
	merge.Set("config", mc)
	cfg.MergeSource(merge, "command")
	if cfg.Get("disable-tls") == true {
		// Factory::createHttpDownloader
		io.WriteError("<warning>You are running Composer with SSL/TLS protection disabled.</warning>")
	}
	d := &diagnoser{ctx: ctx, io: io, cfg: cfg}
	d.facts, _ = probeDiagnose()

	io.Write(fmt.Sprintf("Composer version: <comment>%s</comment>", platform.ComposerVersion))

	io.WriteV("Checking Composer and its dependencies for vulnerabilities: ", false, cio.Normal)
	d.output(d.checkComposerAudit())

	phpVersion := "unknown"
	if c != nil {
		if pr, err := c.CreatePlatformRepo(cfg.Get("platform")); err == nil {
			if p := repo.FindPackage(pr, "php", &semver.MatchAll{}); p != nil {
				phpVersion = p.PrettyVersion()
				if strings.Contains(p.Description(), "overridden") {
					phpVersion += " - " + p.Description()
				}
			}
		}
	} else if d.facts != nil {
		phpVersion = d.facts.Version
	}
	if d.facts == nil {
		io.Write("PHP version: <error>missing (no PHP binary found)</error>")
	} else {
		io.Write(fmt.Sprintf("PHP version: <comment>%s</comment>", phpVersion))
		io.Write(fmt.Sprintf("PHP binary path: <comment>%s</comment>", d.facts.Binary))
	}
	if d.facts != nil && d.facts.OpenSSLText != nil {
		io.Write("OpenSSL version: <comment>" + *d.facts.OpenSSLText + "</comment>")
	} else {
		io.Write("OpenSSL version: <error>missing</error>")
	}
	io.Write("curl version: " + d.curlVersion())

	_, unzipErr := exec.LookPath("unzip")
	hasUnzip := unzipErr == nil
	bin7zip := ""
	for _, b := range []string{"7z", "7zz", "7za"} {
		if _, err := exec.LookPath(b); err == nil {
			bin7zip = b
			break
		}
	}
	zipLine := "zip: "
	if d.facts != nil && d.facts.Ext["zip"] {
		zipLine += "<comment>extension present</comment>"
	} else {
		zipLine += "<comment>extension not loaded</comment>"
	}
	if hasUnzip {
		zipLine += ", <comment>unzip present</comment>"
	} else {
		zipLine += ", <comment>unzip not available</comment>"
	}
	if bin7zip != "" {
		zipLine += ", <comment>7-Zip present (" + bin7zip + ")</comment>"
	} else {
		zipLine += ", <comment>7-Zip not available</comment>"
	}
	if (bin7zip != "" || hasUnzip) && d.facts != nil && !d.facts.ProcOpen {
		zipLine += ", <warning>proc_open is disabled or not present, unzip/7-z will not be usable</warning>"
	}
	io.Write(zipLine)

	if c != nil {
		io.Write("Active plugins: " + strings.Join(composer.RegisteredPlugins(c), ", "))
		io.WriteV("Checking composer.json: ", false, cio.Normal)
		d.output(d.checkComposerSchema())
		if locked, _ := c.Locker.IsLocked(); locked {
			io.WriteV("Checking composer.lock: ", false, cio.Normal)
			d.output(d.checkLockSchema(c.Locker.File))
		}
	}

	io.WriteV("Checking platform settings: ", false, cio.Normal)
	d.output(d.checkPlatform())

	io.WriteV("Checking git settings: ", false, cio.Normal)
	d.output(d.checkGit())

	io.WriteV("Checking http connectivity to packagist: ", false, cio.Normal)
	d.output(d.checkHTTP("http"))
	io.WriteV("Checking https connectivity to packagist: ", false, cio.Normal)
	d.output(d.checkHTTP("https"))

	if repos := cfg.Repositories(); repos != nil {
		for _, k := range repos.Keys() {
			rc := phpjson.ToMap(mustMapGet(repos, k))
			if rc == nil || rc.Str("type") != "composer" {
				continue
			}
			rawURL, ok := rc.Get("url")
			if !ok {
				continue
			}
			ru, _ := phpjson.ScalarString(rawURL)
			u := packagesJSONURL(ru)
			if !strings.HasPrefix(u, "http") || strings.HasPrefix(u, "https://repo.packagist.org") {
				continue
			}
			io.WriteV("Checking connectivity to "+ru+": ", false, cio.Normal)
			d.output(d.checkURL(u, strings.HasPrefix(u, "https://")))
		}
	}

	protos := []string{"http", "https"}
	if cfg.Get("disable-tls") == true {
		protos = []string{"http"}
	}
	if perr := proxyEnvError(); perr != "" {
		io.WriteV("Checking HTTP proxy: ", false, cio.Normal)
		if s, ok := d.networkSkip(); ok {
			d.output(diagMsg(s))
		} else {
			d.output(diagResult{err: &httpdl.TransportError{Msg: "Unable to use a proxy: " + perr}})
		}
	} else {
		for _, proto := range protos {
			status, secure, excluded := proxyStatusFor(proto)
			if status == "" {
				continue
			}
			typ := "HTTP"
			if secure {
				typ = "HTTPS"
			}
			io.WriteV("Checking "+typ+" proxy with "+proto+": ", false, cio.Normal)
			d.output(d.checkHTTPProxy(proto, status, excluded))
		}
	}

	oauth := phpjson.ToMap(cfg.Get("github-oauth"))
	if oauth != nil && oauth.Len() > 0 {
		for _, domain := range oauth.Keys() {
			io.WriteV("Checking "+domain+" oauth access: ", false, cio.Normal)
			d.output(d.checkGithubOauth(domain))
		}
	} else {
		io.WriteV("Checking github.com rate limit: ", false, cio.Normal)
		if s, ok := d.networkSkip(); ok {
			d.output(diagMsg(s))
		} else if r, printed := d.checkRateLimit(); !printed {
			d.output(r)
		}
	}

	io.WriteV("Checking disk free space: ", false, cio.Normal)
	d.output(d.checkDiskSpace())
	return d.exitCode, nil
}

func sha256Hex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func mustMapGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

// output is DiagnoseCommand::outputResult.
func (d *diagnoser) output(r diagResult) {
	if r.ok {
		d.io.Write("<info>OK</info>")
		return
	}
	msgs := r.msgs
	if r.err != nil {
		msgs = []string{"<error>[" + exceptionClass(r.err) + "] " + r.err.Error() + "</error>"}
	}
	hadError, hadWarning := len(msgs) == 0, false
	for _, m := range msgs {
		if strings.Contains(m, "<error>") {
			hadError = true
		} else if strings.Contains(m, "<warning>") {
			hadWarning = true
		}
	}
	if hadError {
		d.io.Write("<error>FAIL</error>")
		if d.exitCode < 2 {
			d.exitCode = 2
		}
	} else if hadWarning {
		d.io.Write("<warning>WARNING</warning>")
		if d.exitCode < 1 {
			d.exitCode = 1
		}
	}
	for _, m := range msgs {
		d.io.Write(strings.TrimSpace(m))
	}
}

func exceptionClass(err error) string {
	var te *httpdl.TransportError
	if errors.As(err, &te) {
		return `Composer\Downloader\TransportException`
	}
	return "RuntimeException"
}

// exceptionHints is HttpDownloader::getExceptionHints.
func exceptionHints(err error) []string {
	var te *httpdl.TransportError
	if !errors.As(err, &te) {
		return nil
	}
	if !strings.Contains(te.Msg, "Resolving timed out") && !strings.Contains(te.Msg, "Could not resolve host") {
		return nil
	}
	client := &http.Client{
		Timeout:       10 * time.Second,
		Transport:     &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	if resp, err := client.Get("https://8.8.8.8"); err == nil {
		resp.Body.Close()
		return []string{"<error>The following exception probably indicates you have misconfigured DNS resolver(s)</error>"}
	}
	return []string{"<error>The following exception probably indicates you are offline or have misconfigured DNS resolver(s)</error>"}
}

// networkSkip is checkConnectivityAndComposerNetworkHttpEnablement (the
// allow_url_fopen part does not apply: gocomposer does not use PHP streams).
func (d *diagnoser) networkSkip() (string, bool) {
	if envTruthy("COMPOSER_DISABLE_NETWORK") {
		return "<info>SKIP</> <comment>Network is disabled by COMPOSER_DISABLE_NETWORK.</>", true
	}
	return "", false
}

func (d *diagnoser) downloader() (*httpdl.Downloader, error) {
	if d.dl == nil {
		dl, err := httpdl.New(d.cfg)
		if err != nil {
			return nil, err
		}
		// prohibitUrlByConfig('http://repo.packagist.org', new NullIO)
		// marks the host as warned without output
		_ = dl.ProhibitURL("http://repo.packagist.org")
		dl.Warn = func(m string) { d.io.WriteError(m) }
		d.dl = dl
	}
	return d.dl, nil
}

func (d *diagnoser) checkComposerAudit() diagResult {
	if s, ok := d.networkSkip(); ok {
		return diagMsg(s)
	}
	return diagMsg("<warning>Could not find Composer's installed.json, this must be a non-standard Composer installation.</>")
}

func (d *diagnoser) curlVersion() string {
	f := d.facts
	if f == nil || f.Curl == nil {
		return "<error>missing, using php streams fallback, which reduces performance</error>"
	}
	if !f.Curl.Multi {
		return "<error>disabled via disable_functions, using php streams fallback, which reduces performance</error>"
	}
	or := func(s string) string {
		if s == "" {
			return "missing"
		}
		return s
	}
	zstd := "missing"
	if f.Curl.Zstd {
		zstd = "supported"
	}
	http := "1.0, 1.1"
	if f.Curl.HTTP2 {
		http += ", 2"
	}
	if f.Curl.HTTP3 {
		http += ", 3"
	}
	return "<comment>" + f.Curl.Version + "</comment> libz <comment>" + or(f.Curl.LibzVersion) + "</comment> brotli <comment>" + or(f.Curl.BrotliVersion) +
		"</comment> zstd <comment>" + zstd + "</comment> ssl <comment>" + or(f.Curl.SSLVersion) + "</comment> HTTP <comment>" + http + "</comment>"
}

func (d *diagnoser) checkComposerSchema() diagResult {
	file, err := composer.ComposerFileName()
	if err != nil {
		return diagResult{err: err}
	}
	errs, _, warnings, fatal := validate.ValidateConfig(file, validate.CheckAll, validate.CheckVersion, composer.ValidateSchemaFile)
	if fatal != nil {
		return diagResult{err: fatal}
	}
	if len(errs) == 0 && len(warnings) == 0 {
		return diagOK
	}
	var out string
	for _, m := range errs {
		out += "<error>" + m + "</error>\n"
	}
	for _, m := range warnings {
		out += "<warning>" + m + "</warning>\n"
	}
	return diagMsg(strings.TrimRight(out, " \t\n\r\x00\x0B"))
}

func (d *diagnoser) checkLockSchema(file string) diagResult {
	err := composer.ValidateSchemaFile(file, jsonschema.LockSchema)
	var ve *validate.JSONValidationError
	if errors.As(err, &ve) {
		var out string
		for _, e := range ve.Errors {
			out += "<error>" + e + "</error>\n"
		}
		return diagMsg(strings.TrimSpace(out))
	}
	if err != nil {
		return diagResult{err: err}
	}
	return diagOK
}

// iniBool is filter_var($v, FILTER_VALIDATE_BOOLEAN).
func iniBool(v *string) bool {
	if v == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*v)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

func (d *diagnoser) checkPlatform() diagResult {
	f := d.facts
	if f == nil {
		return diagMsg("<error>No PHP binary found (set GOCOMPOSER_PHP or add php to PATH)</error>")
	}
	output := ""
	out := func(msg, style string) { output += "<" + style + ">" + msg + "</" + style + ">\n" }
	type item struct {
		key     string
		current string
	}
	var errs, warnings []item
	displayIni := false
	if !f.JSONDecode {
		errs = append(errs, item{key: "json"})
	}
	if !f.Ext["Phar"] {
		errs = append(errs, item{key: "phar"})
	}
	if !f.Ext["filter"] {
		errs = append(errs, item{key: "filter"})
	}
	if !f.Ext["hash"] {
		errs = append(errs, item{key: "hash"})
	}
	if !f.Ext["iconv"] && !f.Ext["mbstring"] {
		errs = append(errs, item{key: "iconv_mbstring"})
	}
	if !iniBool(f.Ini["allow_url_fopen"]) {
		errs = append(errs, item{key: "allow_url_fopen"})
	}
	if len(f.IonCube) == 2 {
		if n, ok := f.IonCube[0].(float64); ok && n < 40009 {
			errs = append(errs, item{key: "ioncube", current: fmt.Sprint(f.IonCube[1])})
		}
	}
	if f.VersionID < 70205 {
		errs = append(errs, item{key: "php", current: f.Version})
	}
	if !f.Ext["openssl"] {
		errs = append(errs, item{key: "openssl"})
	}
	if f.Ext["openssl"] && f.OpenSSLNumber != nil && *f.OpenSSLNumber < 0x1000100f {
		warnings = append(warnings, item{key: "openssl_version"})
	}
	if !f.Ext["apcu"] && iniBool(f.Ini["apc.enable_cli"]) {
		warnings = append(warnings, item{key: "apc_cli"})
	}
	if !f.Ext["zlib"] {
		warnings = append(warnings, item{key: "zlib"})
	}
	if f.Configure != nil {
		if strings.Contains(*f.Configure, "--enable-sigchild") {
			warnings = append(warnings, item{key: "sigchild"})
		}
		if strings.Contains(*f.Configure, "--with-curlwrappers") {
			warnings = append(warnings, item{key: "curlwrappers"})
		}
	}
	if iniBool(f.Ini["xdebug.profiler_enabled"]) {
		warnings = append(warnings, item{key: "xdebug_profile"})
	} else if f.XdebugActive {
		warnings = append(warnings, item{key: "xdebug_loaded"})
	}
	if f.WindowsBuild && (phpVersionCompare(f.Version, "7.2.23") < 0 || (phpVersionCompare(f.Version, "7.3.0") >= 0 && phpVersionCompare(f.Version, "7.3.10") < 0)) {
		warnings = append(warnings, item{key: "onedrive", current: f.Version})
	}
	if f.Ext["uopz"] && !(iniBool(f.Ini["uopz.disable"]) || iniBool(f.Ini["uopz.exit"])) {
		warnings = append(warnings, item{key: "uopz"})
	}
	if len(errs) > 0 {
		for _, e := range errs {
			var text string
			switch e.key {
			case "json":
				text = "\nThe json extension is missing.\nInstall it or recompile php without --disable-json"
			case "phar":
				text = "\nThe phar extension is missing.\nInstall it or recompile php without --disable-phar"
			case "filter":
				text = "\nThe filter extension is missing.\nInstall it or recompile php without --disable-filter"
			case "hash":
				text = "\nThe hash extension is missing.\nInstall it or recompile php without --disable-hash"
			case "iconv_mbstring":
				text = "\nThe iconv OR mbstring extension is required and both are missing.\nInstall either of them or recompile php without --disable-iconv"
			case "php":
				text = "\nYour PHP (" + e.current + ") is too old, you must upgrade to PHP 7.2.5 or higher."
			case "allow_url_fopen":
				text = "\nThe allow_url_fopen setting is incorrect.\nAdd the following to the end of your `php.ini`:\n    allow_url_fopen = On"
				displayIni = true
			case "ioncube":
				text = "\nYour ionCube Loader extension (" + e.current + ") is incompatible with Phar files.\nUpgrade to ionCube 4.0.9 or higher or remove this line (path may be different) from your `php.ini` to disable it:\n    zend_extension = /usr/lib/php5/20090626+lfs/ioncube_loader_lin_5.3.so"
				displayIni = true
			case "openssl":
				text = "\nThe openssl extension is missing, which means that secure HTTPS transfers are impossible.\nIf possible you should enable it or recompile php with --with-openssl"
			}
			out(text, "error")
		}
		output += "\n"
	}
	for _, w := range warnings {
		var text string
		switch w.key {
		case "apc_cli":
			text = "The apc.enable_cli setting is incorrect.\nAdd the following to the end of your `php.ini`:\n  apc.enable_cli = Off"
			displayIni = true
		case "zlib":
			text = "The zlib extension is not loaded, this can slow down Composer a lot.\nIf possible, enable it or recompile php with --with-zlib\n"
			displayIni = true
		case "sigchild":
			text = "PHP was compiled with --enable-sigchild which can cause issues on some platforms.\nRecompile it without this flag if possible, see also:\n  https://bugs.php.net/bug.php?id=22999"
		case "curlwrappers":
			text = "PHP was compiled with --with-curlwrappers which will cause issues with HTTP authentication and GitHub.\n Recompile it without this flag if possible"
		case "openssl_version":
			v := ""
			if f.OpenSSLText != nil {
				t := *f.OpenSSLText
				if i := strings.Index(t, " "); i >= 0 {
					rest := strings.TrimSpace(t[i:])
					if j := strings.Index(rest, " "); j >= 0 {
						v = rest[:j]
					}
				}
				if v == "" {
					v = t
				}
			}
			text = "The OpenSSL library (" + v + ") used by PHP does not support TLSv1.2 or TLSv1.1.\nIf possible you should upgrade OpenSSL to version 1.0.1 or above."
		case "xdebug_loaded":
			text = "The xdebug extension is loaded, this can slow down Composer a little.\n Disabling it when using Composer is recommended."
		case "xdebug_profile":
			text = "The xdebug.profiler_enabled setting is enabled, this can slow down Composer a lot.\nAdd the following to the end of your `php.ini` to disable it:\n  xdebug.profiler_enabled = 0"
			displayIni = true
		case "onedrive":
			text = "The Windows OneDrive folder is not supported on PHP versions below 7.2.23 and 7.3.10.\nUpgrade your PHP (" + w.current + ") to use this location with Composer.\n"
		case "uopz":
			text = "The uopz extension ignores exit calls and may not work with all Composer commands.\nDisabling it when using Composer is recommended."
		}
		out(text, "comment")
	}
	if displayIni {
		out("\n\n"+iniMessage(f.IniFiles)+"\nIf you can not modify the ini file, you can also run `php -d option=value` to modify ini values on the fly. You can use -d multiple times.", "comment")
	}
	if v := os.Getenv("COMPOSER_IPRESOLVE"); v == "4" || v == "6" {
		warnings = append(warnings, item{key: "ipresolve"})
		out("The COMPOSER_IPRESOLVE env var is set to "+v+" which may result in network failures below.", "comment")
	}
	if len(warnings) == 0 && len(errs) == 0 {
		return diagOK
	}
	return diagMsg(output)
}

// iniMessage is IniHelper::getMessage.
func iniMessage(paths []string) string {
	if len(paths) > 0 && paths[0] == "" {
		paths = paths[1:]
	}
	if len(paths) == 0 || paths[0] == "" {
		return "A php.ini file does not exist. You will have to create one."
	}
	if len(paths) > 1 {
		return "Your command-line PHP is using multiple ini files. Run `php --ini` to show them."
	}
	return "The php.ini used by your command-line PHP is: " + paths[0]
}

var reGitVersion = regexp.MustCompile(`(?m)^git version (\d+(?:\.\d+)+)`)

func (d *diagnoser) checkGit() diagResult {
	if d.facts != nil && !d.facts.ProcOpen {
		return diagMsg("<comment>proc_open is not available, git cannot be used</comment>")
	}
	colorOut, _ := exec.Command("git", "config", "color.ui").Output()
	if strings.ToLower(strings.TrimSpace(string(colorOut))) == "always" {
		return diagMsg(`<comment>Your git color.ui setting is set to always, this is known to create issues. Use "git config --global color.ui true" to set it correctly.</comment>`)
	}
	verOut, err := exec.Command("git", "--version").Output()
	m := reGitVersion.FindStringSubmatch(string(verOut))
	if err != nil || m == nil {
		return diagMsg("<comment>No git process found</>")
	}
	if phpVersionCompare("2.24.0", m[1]) > 0 {
		return diagMsg("<warning>Your git version (" + m[1] + ") is too old and possibly will cause issues. Please upgrade to git 2.24 or above</>")
	}
	return diagMsg("<info>OK</> <comment>git version " + m[1] + "</>")
}

// phpVersionCompare is version_compare() for plain dotted numeric versions.
func phpVersionCompare(a, b string) int {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		if i >= len(pa) {
			return -1
		}
		if i >= len(pb) {
			return 1
		}
		na, _ := strconv.Atoi(pa[i])
		nb, _ := strconv.Atoi(pb[i])
		if na != nb {
			if na < nb {
				return -1
			}
			return 1
		}
	}
	return 0
}

const tlsDisabledWarning = "<warning>Composer is configured to disable SSL/TLS protection. This will leave remote HTTPS requests vulnerable to Man-In-The-Middle attacks.</warning>"

func (d *diagnoser) checkHTTP(proto string) diagResult {
	return d.checkURL(proto+"://repo.packagist.org/packages.json", proto == "https")
}

// checkURL is checkHttp / checkComposerRepo.
func (d *diagnoser) checkURL(u string, https bool) diagResult {
	if s, ok := d.networkSkip(); ok {
		return diagMsg(s)
	}
	var result []string
	dl, err := d.downloader()
	if err == nil {
		_, err = dl.Get(u)
	}
	if err != nil {
		result = append(result, exceptionHints(err)...)
		result = append(result, "<error>["+exceptionClass(err)+"] "+err.Error()+"</error>")
	}
	if https && d.cfg.Get("disable-tls") == true {
		result = append(result, tlsDisabledWarning)
	}
	if len(result) > 0 {
		return diagMsg(result...)
	}
	return diagOK
}

func (d *diagnoser) checkHTTPProxy(proto, status string, excluded bool) diagResult {
	if s, ok := d.networkSkip(); ok {
		return diagMsg(s)
	}
	if excluded {
		return diagMsg("<info>SKIP</> <comment>Because repo.packagist.org is " + status + "</>")
	}
	dl, err := d.downloader()
	if err != nil {
		return diagResult{err: err}
	}
	resp, err := dl.Get(proto + "://repo.packagist.org/packages.json")
	if err != nil {
		return diagResult{err: err}
	}
	v, err := phpjson.Decode(resp.Body)
	if err != nil {
		return diagResult{err: err}
	}
	if inc := phpjson.ToMap(mustMapGet(phpjson.ToMap(v), "provider-includes")); inc != nil && inc.Len() > 0 {
		path := inc.Keys()[0]
		hash := phpjson.ToMap(mustMapGet(inc, path)).Str("sha256")
		pr, err := dl.Get(proto + "://repo.packagist.org/" + strings.ReplaceAll(path, "%hash%", hash))
		if err != nil {
			return diagResult{err: err}
		}
		if sha256Hex(pr.Body) != hash {
			return diagMsg("<warning>It seems that your proxy (" + status + ") is modifying " + proto + " traffic on the fly</>")
		}
	}
	return diagMsg("<info>OK</> <comment>" + status + "</>")
}

func (d *diagnoser) checkGithubOauth(domain string) diagResult {
	if s, ok := d.networkSkip(); ok {
		return diagMsg(s)
	}
	dl, err := d.downloader()
	if err != nil {
		return diagResult{err: err}
	}
	u := "https://" + domain + "/api/v3/"
	if domain == "github.com" {
		u = "https://api.github.com/"
	}
	resp, err := dl.Get(u)
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) && te.StatusCode == 401 {
			return diagMsg(`<comment>The oauth token for ` + domain + ` seems invalid, run "composer config --global --unset github-oauth.` + domain + `" to remove it</comment>`)
		}
		return diagResult{err: err}
	}
	if exp := resp.Header.Get("github-authentication-token-expiration"); exp != "" {
		return diagMsg("<info>OK</> <comment>expires on " + exp + "</>")
	}
	return diagMsg("<info>OK</> <comment>does not expire</>")
}

// checkRateLimit reports whether it already printed the low-limit warning.
func (d *diagnoser) checkRateLimit() (diagResult, bool) {
	dl, err := d.downloader()
	if err != nil {
		return diagResult{err: err}, false
	}
	resp, err := dl.Get("https://api.github.com/rate_limit")
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) && te.StatusCode == 401 {
			return diagMsg(`<comment>The oauth token for github.com seems invalid, run "composer config --global --unset github-oauth.github.com" to remove it</comment>`), false
		}
		return diagResult{err: err}, false
	}
	v, err := phpjson.Decode(resp.Body)
	if err != nil {
		return diagResult{err: err}, false
	}
	core := phpjson.ToMap(mustMapGet(phpjson.ToMap(mustMapGet(phpjson.ToMap(v), "resources")), "core"))
	if core == nil {
		return diagResult{}, false
	}
	remaining := phpInt(mustMapGet(core, "remaining"), nil)
	if remaining < 10 {
		d.io.Write("<warning>WARNING</warning>")
		d.io.Write(fmt.Sprintf("<comment>GitHub has a rate limit on their API. You currently have <options=bold>%d</options=bold> out of <options=bold>%d</options=bold> requests left.\nSee https://developer.github.com/v3/#rate-limiting and also\n    https://getcomposer.org/doc/articles/troubleshooting.md#api-rate-limit-and-oauth-tokens</comment>", remaining, phpInt(mustMapGet(core, "limit"), nil)))
		return diagResult{}, true
	}
	return diagOK, false
}

func (d *diagnoser) checkDiskSpace() diagResult {
	const minFree = 1024 * 1024
	for _, key := range []string{"home", "vendor-dir"} {
		dir := d.cfg.GetString(key)
		if free, ok := diskFreeSpace(dir); ok && free < minFree {
			return diagMsg("<error>The disk hosting " + dir + " is full</error>")
		}
	}
	return diagOK
}

// packagesJSONURL is new ComposerRepository($repo)->getPackagesJsonUrl().
func packagesJSONURL(raw string) string {
	u := raw
	if !regexp.MustCompile(`^[\w.]+\??://`).MatchString(u) {
		if rp := phpRealpath(u); rp != "" {
			u = "file://" + rp
		} else {
			u = "http://" + u
		}
	}
	u = strings.TrimRight(u, "/")
	if strings.HasPrefix(u, "https?") {
		u = "https" + u[6:]
	}
	if m := regexp.MustCompile(`(?i)^(https?)://packagist\.org/?$`).FindStringSubmatch(u); m != nil {
		u = m[1] + "://repo.packagist.org"
	}
	if pu, err := url.Parse(strings.ReplaceAll(u, `\`, "/")); err == nil && strings.Contains(pu.Path, ".json") {
		return u
	}
	return u + "/packages.json"
}

// proxyEnv is ProxyManager::getProxyEnv.
func proxyEnv(name string) (string, string) {
	for _, n := range []string{strings.ToLower(name), strings.ToUpper(name)} {
		if v := os.Getenv(n); v != "" {
			return v, n
		}
	}
	return "", ""
}

type proxyItem struct {
	url, safeURL string
}

// newProxyItem is ProxyItem::__construct.
func newProxyItem(proxyURL, envName string) (*proxyItem, string) {
	syntaxErr := "unsupported `" + envName + "` syntax"
	if strings.ContainsAny(proxyURL, "\r\n\t") {
		return nil, syntaxErr
	}
	raw := proxyURL
	if !strings.Contains(raw, "://") {
		raw = "//" + raw
	}
	pu, err := url.Parse(raw)
	if err != nil {
		return nil, syntaxErr
	}
	if pu.Hostname() == "" {
		return nil, "unable to find proxy host in " + envName
	}
	scheme := "http://"
	if pu.Scheme != "" {
		scheme = strings.ToLower(pu.Scheme) + "://"
	}
	safe := ""
	if pu.User != nil {
		safe = "***"
		if _, has := pu.User.Password(); has {
			safe += ":***"
		}
		safe += "@"
	}
	port := pu.Port()
	if port == "" {
		switch scheme {
		case "http://":
			port = "80"
		case "https://":
			port = "443"
		default:
			return nil, "unable to find proxy port in " + envName
		}
	}
	if n, _ := strconv.Atoi(port); n == 0 {
		return nil, "port 0 is reserved in " + envName
	}
	return &proxyItem{url: scheme + pu.Hostname() + ":" + port, safeURL: scheme + safe + pu.Hostname() + ":" + port}, ""
}

func proxyItems() (httpP, httpsP *proxyItem, errMsg string) {
	if v, n := proxyEnv("http_proxy"); v != "" {
		p, e := newProxyItem(v, n)
		if e != "" {
			return nil, nil, e
		}
		httpP = p
	}
	if httpP == nil {
		if v, n := proxyEnv("cgi_http_proxy"); v != "" {
			p, e := newProxyItem(v, n)
			if e != "" {
				return nil, nil, e
			}
			httpP = p
		}
	}
	if v, n := proxyEnv("https_proxy"); v != "" {
		p, e := newProxyItem(v, n)
		if e != "" {
			return nil, nil, e
		}
		httpsP = p
	}
	return httpP, httpsP, ""
}

func proxyEnvError() string {
	_, _, e := proxyItems()
	return e
}

// proxyStatusFor is getProxyForRequest($proto.'://repo.packagist.org').
func proxyStatusFor(proto string) (status string, secure, excluded bool) {
	httpP, httpsP, _ := proxyItems()
	p := httpP
	if proto == "https" {
		p = httpsP
	}
	if p == nil {
		return "", false, false
	}
	if noProxyMatches("repo.packagist.org") {
		return "excluded by no_proxy", false, true
	}
	return p.safeURL, strings.HasPrefix(p.url, "https://"), false
}

// noProxyMatches is a host-name subset of NoProxyPattern::test.
func noProxyMatches(host string) bool {
	v, _ := proxyEnv("no_proxy")
	if v == "" {
		return false
	}
	for _, e := range regexp.MustCompile(`[\s,]+`).Split(v, -1) {
		e = strings.ToLower(strings.TrimSpace(e))
		if e == "" {
			continue
		}
		if e == "*" {
			return true
		}
		if h, _, ok := strings.Cut(e, ":"); ok {
			e = h
		}
		e = strings.TrimPrefix(e, ".")
		if host == e || strings.HasSuffix(host, "."+e) {
			return true
		}
	}
	return false
}
