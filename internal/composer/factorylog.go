package composer

// The IOInterface::DEBUG output of Composer's Factory (-vvv): JsonFile::read
// "Reading ..." lines, the "Loading config file ..." lines of createComposer,
// the CA bundle lookup of StreamContextFactory::getTlsDefaults /
// composer/ca-bundle CaBundle::getSystemCaRootBundlePath and validateCaFile,
// and Factory::createGlobalComposer's "Failed to initialize global
// composer" (composer/composer f83588392df43d4fd027d84d26de459946ef4e80 and
// composer/ca-bundle 1.5.x as locked there, MIT, (c) Nils Adermann, Jordi
// Boggiano). These only write output; the actual work is done elsewhere.

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/phpjson"
)

// ReadingLog is the debug line of JsonFile::read for a JsonFile created
// with an IO: "Reading <path>" plus the realpath when it differs.
func ReadingLog(io *cio.IO, path string) {
	if io == nil || !io.IsDebug() {
		return
	}
	info := ""
	if rp := phpRealpath(path); rp != "" && rp != path {
		info = " (" + rp + ")"
	}
	io.WriteError("Reading " + path + info)
}

// phpRealpath is PHP's realpath() ("" when the path does not exist).
func phpRealpath(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	r, err := filepath.EvalSymlinks(a)
	if err != nil {
		return ""
	}
	return r
}

func debugLogger(io *cio.IO) func(string) {
	if io == nil || !io.IsDebug() {
		return nil
	}
	return func(m string) { io.WriteErrorV(m, true, cio.Debug) }
}

var (
	caMu       sync.Mutex
	caPath     *string
	caValidity = map[string]bool{}
)

var reTrustedCert = regexp.MustCompile(`(?m)^(-+(?:BEGIN|END))\s+TRUSTED\s+(CERTIFICATE-+)$`)

// validateCaFile is CaBundle::validateCaFile (cached per file name).
func validateCaFile(filename string, log func(string)) bool {
	if v, ok := caValidity[filename]; ok {
		return v
	}
	valid := false
	if b, err := os.ReadFile(filename); err == nil && len(b) > 0 {
		contents := reTrustedCert.ReplaceAll(b, []byte("$1 $2"))
		for len(contents) > 0 {
			var block *pem.Block
			block, contents = pem.Decode(contents)
			if block == nil {
				break
			}
			if block.Type == "CERTIFICATE" || block.Type == "X509 CERTIFICATE" {
				_, err := x509.ParseCertificate(block.Bytes)
				valid = err == nil
				break
			}
		}
	}
	if log != nil {
		state := "invalid"
		if valid {
			state = "valid"
		}
		log("Checked CA file " + phpRealpath(filename) + ": " + state)
	}
	caValidity[filename] = valid
	return valid
}

func isReadable(p string) bool {
	f, err := os.Open(p)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

// systemCaRootBundlePath is CaBundle::getSystemCaRootBundlePath (the
// openssl.cafile/capath ini settings are not known to gocomposer).
func systemCaRootBundlePath(log func(string)) string {
	if caPath != nil {
		return *caPath
	}
	var candidates []string
	for _, env := range []string{"SSL_CERT_FILE", "SSL_CERT_DIR"} {
		if v, ok := os.LookupEnv(env); ok {
			candidates = append(candidates, v)
		} else {
			candidates = append(candidates, "")
		}
	}
	candidates = append(candidates,
		"/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
		"/etc/pki/tls/certs/ca-bundle.crt",
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/ssl/ca-bundle.pem",
		"/usr/ssl/certs/ca-bundle.crt",
		"/opt/local/share/curl/curl-ca-bundle.crt",
		"/usr/local/share/curl/curl-ca-bundle.crt",
		"/usr/share/ssl/certs/ca-bundle.crt",
		"/etc/ssl/cert.pem",
		"/usr/local/etc/openssl/cert.pem",
		"/usr/local/etc/openssl@1.1/cert.pem",
		"/opt/homebrew/etc/openssl@3/cert.pem",
		"/opt/homebrew/etc/openssl@1.1/cert.pem",
		"/etc/pki/tls/certs",
		"/etc/ssl/certs",
	)
	logf := func(m string) {
		if log != nil {
			log(m)
		}
	}
	for _, ca := range candidates {
		if ca == "" || ca == "0" {
			continue
		}
		// caFileUsable
		st, err := os.Stat(ca)
		if err != nil || !st.Mode().IsRegular() {
			logf("Checked CA file " + ca + " does not exist or it is not a file.")
		} else if !isReadable(ca) {
			logf("Checked file or directory " + ca + " is not readable.")
		} else if validateCaFile(ca, log) {
			caPath = &ca
			return ca
		}
		// caDirUsable
		if err != nil || !st.IsDir() {
			logf("Checked directory " + ca + " does not exist or it is not a directory.")
			continue
		}
		if !isReadable(ca) {
			logf("Checked file or directory " + ca + " is not readable.")
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(ca, "*"))
		if len(matches) == 0 {
			logf("No CA files found for pattern: " + ca + "/*")
			continue
		}
		caPath = &ca
		return ca
	}
	bundled := ""
	caPath = &bundled
	return bundled
}

// tlsDefaultsDebug writes the CA checks of HttpDownloader's
// StreamContextFactory::getTlsDefaults (Factory::createHttpDownloader).
func tlsDefaultsDebug(io *cio.IO, cfg *config.Config) {
	log := debugLogger(io)
	if log == nil || cfg.Get("disable-tls") == true {
		return
	}
	caMu.Lock()
	defer caMu.Unlock()
	cafile, _ := cfg.Get("cafile").(string)
	capath, _ := cfg.Get("capath").(string)
	if cafile == "" && capath == "" {
		result := systemCaRootBundlePath(log)
		if st, err := os.Stat(result); err == nil && st.IsDir() {
			return
		}
		cafile = result
	}
	if cafile != "" && isReadable(cafile) {
		validateCaFile(cafile, log)
	}
}

// globalComposerDebug is the debug output of Factory::createGlobalComposer
// (createComposer for COMPOSER_HOME/composer.json without full load).
func globalComposerDebug(io *cio.IO, cfg *config.Config) {
	log := debugLogger(io)
	if log == nil {
		return
	}
	home := cfg.GetString("home")
	file := home + "/composer.json"
	fail := func(msg string) { log("Failed to initialize global composer: " + msg) }
	if st, err := os.Stat(file); err != nil || !st.Mode().IsRegular() {
		fail("Composer could not find the config file: " + file + "\n")
		return
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		fail(err.Error())
		return
	}
	data, err := ParseJSONFile(raw, file)
	if err != nil {
		fail(err.Error())
		return
	}
	if err := ValidateSchemaLax(data, file); err != nil {
		fail(err.Error())
		return
	}
	ReadingLog(io, file)
	gcfg, err := config.CreateConfigDebug(home, nil, log)
	if err != nil {
		fail(err.Error())
		return
	}
	log("Loading config file " + file + " (" + phpRealpath(file) + ")")
	if auth := filepath.Dir(phpRealpath(file)) + "/auth.json"; isFile(auth) {
		log("Loading config file " + auth)
		if err := config.ValidateSchemaWarn(nil, auth, nil, jsonschema.AuthSchema); err != nil {
			fail(err.Error())
			return
		}
		ReadingLog(io, auth)
	}
	gcfg.SetDebug(log)
	if err := gcfg.LoadAuthEnv(); err != nil {
		fail(err.Error())
	}
}

// globalInstalledRepoFile is the installed.json of the global Composer
// instance PluginManager::loadInstalledPlugins reads for global plugins
// ("" when Factory::createGlobalComposer would fail).
func globalInstalledRepoFile(cfg *config.Config) string {
	home := cfg.GetString("home")
	file := home + "/composer.json"
	if !isFile(file) {
		return ""
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	data, err := ParseJSONFile(raw, file)
	if err != nil || ValidateSchemaLax(data, file) != nil {
		return ""
	}
	gcfg, err := config.CreateConfig(home, nil)
	if err != nil {
		return ""
	}
	gcfg.Merge(data)
	return gcfg.GetString("vendor-dir") + "/composer/installed.json"
}

func isFile(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}

// pluginManagerDebug is the -vvv output of creating the PluginManager and
// PluginManager::loadInstalledPlugins: parseAllowedPlugins reads the lock
// file when allow-plugins is empty, then the local repository
// (vendor/composer/installed.json) is read.
func pluginManagerDebug(io *cio.IO, c *Composer, factoryCwd string) {
	if !io.IsDebug() {
		return
	}
	if ap := c.Config.Get("allow-plugins"); phpjson.IsArray(ap) && phpjson.ArrayLen(ap) == 0 && c.Locker.Exists() {
		_, _ = c.Locker.LockData()
	}
	defer func() {
		// global plugins (createGlobalComposer's local repository)
		if !c.DisablePlugins && !c.isGlobal() {
			if g := globalInstalledRepoFile(c.Config); g != "" && isFile(g) {
				ReadingLog(io, g)
			}
		}
	}()
	if !isFile(c.LocalRepo.File) {
		return
	}
	path := c.LocalRepo.File
	if factoryCwd != "" {
		// Config::get('vendor-dir') relative to Factory's cwd
		if rel, ok := c.Config.GetRelative("vendor-dir").(string); ok && rel != "" && !filepath.IsAbs(rel) {
			path = factoryCwd + "/" + rel + "/composer/installed.json"
		}
	}
	ReadingLog(io, path)
}

var rePasswordOpt = regexp.MustCompile(`--password (.*[^\\]') `)

// scriptCommandDebug is ProcessExecutor::outputCommandRun for a command
// line run without a working directory ("CWD").
func scriptCommandDebug(io *cio.IO, line string) {
	if io == nil || !io.IsDebug() {
		return
	}
	safe := httpdl.Sanitize(line)
	safe = rePasswordOpt.ReplaceAllString(safe, "--password '***' ")
	io.WriteError("Executing command (CWD): " + safe)
}
