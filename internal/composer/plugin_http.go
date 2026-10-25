package composer

// Per request TLS options of plugin HTTP requests
// (HttpDownloader::get/add/copy with $options['ssl']): Composer's
// CurlDownloader applies cafile, capath, verify_peer, verify_peer_name,
// local_cert, local_pk and passphrase to the request (CurlDownloader::$options
// of Composer f83588392df43d4fd027d84d26de459946ef4e80, MIT license, (c) Nils
// Adermann, Jordi Boggiano); the other stream context options
// (peer_fingerprint, http.proxy, request_fulluri) have no curl equivalent
// and are ignored, the proxy coming from the environment (ProxyManager).
// A request with allow_self_signed is made with PHP streams instead
// (HttpDownloader::canUseCurl), which apply all of them: self-signed peer
// certificates are accepted, peer_fingerprint is checked and http.proxy is
// used. gocomposer makes such requests with a client of its own for those
// settings, otherwise like its other requests.

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
)

var reHashedCertName = regexp.MustCompile(`^[0-9a-f]{8}\.\d+$`)

// sslVerifyWarning is the warning Config::prohibitUrlByConfig prints (once
// per host and Config object, here the configuration of the running
// instance) when a request disables verify_peer or verify_peer_name.
func (pm *PluginManager) sslVerifyWarning(host string, ssl *phpjson.Map) {
	if host == "" || ssl == nil || pm.sslWarned[host] {
		return
	}
	warning := ""
	if v, ok := ssl.Get("verify_peer"); ok && !phpjson.Truthy(v) {
		warning = "verify_peer"
	}
	if v, ok := ssl.Get("verify_peer_name"); ok && !phpjson.Truthy(v) {
		if warning == "" {
			warning = "verify_peer_name"
		} else {
			warning += " and verify_peer_name"
		}
	}
	if warning == "" {
		return
	}
	pm.c.IO.WriteError("<warning>Warning: Accessing " + host + " with " + warning + " disabled.</warning>")
	if pm.sslWarned == nil {
		pm.sslWarned = map[string]bool{}
	}
	pm.sslWarned[host] = true
}

// tlsRequestOptions are the ssl options CurlDownloader maps to curl.
var tlsRequestOptions = []string{"cafile", "capath", "verify_peer", "verify_peer_name", "local_cert", "local_pk", "passphrase"}

// requestTLS returns the ssl options of a request that differ from what
// gocomposer's client applies (nil when none do).
func (pm *PluginManager) requestTLS(options *phpjson.Map) *phpjson.Map {
	if options == nil {
		return nil
	}
	ssl := options.Map("ssl")
	if ssl == nil {
		return nil
	}
	out := phpjson.NewMap()
	if v, ok := ssl.Get("allow_self_signed"); ok && phpjson.Truthy(v) {
		// RemoteFilesystem: the stream context options apply
		out.Set("allow_self_signed", true)
		if fp, ok := ssl.Get("peer_fingerprint"); ok && fp != nil {
			out.Set("peer_fingerprint", fp)
		}
		if hm := options.Map("http"); hm != nil {
			if p, ok := hm.Get("proxy"); ok && p != nil {
				out.Set("proxy", p)
			}
		}
	}
	for _, k := range tlsRequestOptions {
		v, ok := ssl.Get(k)
		if !ok || v == nil {
			continue
		}
		switch k {
		case "cafile", "capath":
			// the config's CA (HttpDownloader passes it with every request)
			if s, _ := v.(string); s == "" || s == pm.c.Config.GetString(k) {
				continue
			}
		case "verify_peer", "verify_peer_name":
			if phpjson.Truthy(v) {
				continue
			}
		}
		out.Set(k, v)
	}
	if out.Len() == 0 {
		return nil
	}
	return out
}

// tlsDownloader returns gocomposer's HTTP client set up with the ssl
// options of a request.
func (pm *PluginManager) tlsDownloader(ssl *phpjson.Map) (*httpdl.Downloader, error) {
	key, _ := phpjson.Encode(ssl, 0)
	if d, ok := pm.tlsClients[key]; ok {
		return d, nil
	}
	cfg := &tls.Config{}
	str := func(k string) string { s, _ := mustGet(ssl, k).(string); return s }
	cafile, capath := str("cafile"), str("capath")
	if cafile == "" && capath == "" {
		cafile, capath = pm.c.Config.GetString("cafile"), pm.c.Config.GetString("capath")
	}
	if cafile != "" || capath != "" {
		pool := x509.NewCertPool()
		if cafile == "" {
			// curl keeps its default CA bundle next to CURLOPT_CAPATH
			if sys, err := x509.SystemCertPool(); err == nil {
				pool = sys
			}
		}
		if cafile != "" {
			raw, err := os.ReadFile(cafile)
			if err != nil {
				return nil, &UserError{Msg: fmt.Sprintf("Could not read the cafile %s: %v", cafile, err), Code: 1}
			}
			pool.AppendCertsFromPEM(raw)
		}
		if capath != "" {
			// OpenSSL looks certificates up by their subject hash file
			// names (<hash>.<n>, as created by c_rehash)
			entries, _ := os.ReadDir(capath)
			for _, e := range entries {
				if !reHashedCertName.MatchString(e.Name()) {
					continue
				}
				if raw, err := os.ReadFile(filepath.Join(capath, e.Name())); err == nil {
					pool.AppendCertsFromPEM(raw)
				}
			}
		}
		cfg.RootCAs = pool
	}
	verifyPeer := true
	if v, ok := ssl.Get("verify_peer"); ok && !phpjson.Truthy(v) {
		verifyPeer = false
	}
	verifyName := true
	if v, ok := ssl.Get("verify_peer_name"); ok && !phpjson.Truthy(v) {
		verifyName = false
	}
	selfSigned := mustGet(ssl, "allow_self_signed") == true
	fingerprint := mustGet(ssl, "peer_fingerprint")
	if !verifyPeer {
		cfg.InsecureSkipVerify = true //nolint:gosec // the plugin asked for verify_peer=false
	} else if selfSigned || fingerprint != nil {
		// verify_peer with allow_self_signed (a self-signed peer
		// certificate is accepted) and peer_fingerprint
		roots := cfg.RootCAs
		cfg.InsecureSkipVerify = true //nolint:gosec // verified below
		cfg.VerifyConnection = func(cs tls.ConnectionState) error {
			certs := cs.PeerCertificates
			if len(certs) == 0 {
				return fmt.Errorf("no peer certificate")
			}
			if fingerprint != nil && !fingerprintMatches(certs[0], fingerprint) {
				return fmt.Errorf("peer_fingerprint match failure")
			}
			leaf := certs[0]
			if selfSigned && leaf.CheckSignatureFrom(leaf) == nil {
				return nil
			}
			inter := x509.NewCertPool()
			for _, c := range certs[1:] {
				inter.AddCert(c)
			}
			opts := x509.VerifyOptions{Roots: roots, Intermediates: inter}
			if verifyName {
				opts.DNSName = cs.ServerName
			}
			_, err := leaf.Verify(opts)
			return err
		}
	} else if !verifyName {
		// CURLOPT_SSL_VERIFYHOST off: the chain is verified, not the name
		roots := cfg.RootCAs
		cfg.InsecureSkipVerify = true //nolint:gosec // the chain is verified below
		cfg.VerifyPeerCertificate = func(raw [][]byte, _ [][]*x509.Certificate) error {
			if len(raw) == 0 {
				return fmt.Errorf("no peer certificate")
			}
			certs := make([]*x509.Certificate, 0, len(raw))
			for _, r := range raw {
				c, err := x509.ParseCertificate(r)
				if err != nil {
					return err
				}
				certs = append(certs, c)
			}
			inter := x509.NewCertPool()
			for _, c := range certs[1:] {
				inter.AddCert(c)
			}
			_, err := certs[0].Verify(x509.VerifyOptions{Roots: roots, Intermediates: inter})
			return err
		}
	}
	if cert := str("local_cert"); cert != "" {
		pair, err := loadClientCert(cert, str("local_pk"), str("passphrase"))
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []tls.Certificate{pair}
	}
	d, err := httpdl.New(pm.c.Config)
	if err != nil {
		return nil, err
	}
	d.Warn, d.Debug, d.IO, d.Sources = pm.c.HTTP.Warn, pm.c.HTTP.Debug, pm.c.HTTP.IO, pm.c.HTTP.Sources
	d.UserAgent, d.Timeout = pm.c.HTTP.UserAgent, pm.c.HTTP.Timeout
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment
	transport.DialContext = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext
	transport.TLSClientConfig = cfg
	if p, _ := mustGet(ssl, "proxy").(string); p != "" {
		// the http.proxy stream context option (tcp://host:port)
		pu, err := url.Parse(strings.Replace(p, "tcp://", "http://", 1))
		if err != nil {
			return nil, &UserError{Msg: "Invalid http.proxy option " + p, Code: 1}
		}
		transport.Proxy = http.ProxyURL(pu)
	}
	d.SetTransport(transport)
	if pm.tlsClients == nil {
		pm.tlsClients = map[string]*httpdl.Downloader{}
	}
	pm.tlsClients[key] = d
	return d, nil
}

// loadClientCert loads the client certificate of local_cert/local_pk (the
// key may be in the certificate file, encrypted with passphrase in the
// traditional PEM format).
func loadClientCert(certFile, keyFile, passphrase string) (tls.Certificate, error) {
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return tls.Certificate{}, &UserError{Msg: fmt.Sprintf("Could not read the local_cert %s: %v", certFile, err), Code: 1}
	}
	keyPEM := certPEM
	if keyFile != "" {
		if keyPEM, err = os.ReadFile(keyFile); err != nil {
			return tls.Certificate{}, &UserError{Msg: fmt.Sprintf("Could not read the local_pk %s: %v", keyFile, err), Code: 1}
		}
	}
	var keyOut []byte
	for rest := keyPEM; ; {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			continue
		}
		if block.Type == "ENCRYPTED PRIVATE KEY" {
			return tls.Certificate{}, &UserError{Msg: "The encrypted PKCS#8 private key " + keyFile + " is not supported by gocomposer's HTTP client (use a traditional PEM key)", Code: 1}
		}
		//nolint:staticcheck // legacy PEM encryption is what curl's SSLKEYPASSWD decrypts too
		if x509.IsEncryptedPEMBlock(block) {
			der, err := x509.DecryptPEMBlock(block, []byte(passphrase)) //nolint:staticcheck
			if err != nil {
				return tls.Certificate{}, &UserError{Msg: "Could not decrypt the private key " + keyFile + ": " + err.Error(), Code: 1}
			}
			block = &pem.Block{Type: block.Type, Bytes: der}
		}
		keyOut = pem.EncodeToMemory(block)
		break
	}
	if keyOut == nil {
		return tls.Certificate{}, &UserError{Msg: "No private key found for the local_cert " + certFile, Code: 1}
	}
	pair, err := tls.X509KeyPair(certPEM, keyOut)
	if err != nil {
		return tls.Certificate{}, &UserError{Msg: "Could not load the client certificate " + certFile + ": " + err.Error(), Code: 1}
	}
	return pair, nil
}

// fingerprintMatches checks the peer_fingerprint stream context option: a
// hash string (md5 when 32 characters long, sha1 when 40) or an array of
// algorithm => hash.
func fingerprintMatches(cert *x509.Certificate, fp phpjson.Value) bool {
	check := func(algo, want string) bool {
		var sum []byte
		switch strings.ToLower(algo) {
		case "md5":
			s := md5.Sum(cert.Raw)
			sum = s[:]
		case "sha1":
			s := sha1.Sum(cert.Raw)
			sum = s[:]
		case "sha256":
			s := sha256.Sum256(cert.Raw)
			sum = s[:]
		default:
			return false
		}
		return strings.EqualFold(hex.EncodeToString(sum), want)
	}
	if s, ok := fp.(string); ok {
		switch len(s) {
		case 32:
			return check("md5", s)
		case 40:
			return check("sha1", s)
		}
		return false
	}
	m := phpjson.ToMap(fp)
	if m == nil || m.Len() == 0 {
		return false
	}
	for _, k := range m.Keys() {
		s, _ := mustGet(m, k).(string)
		if !check(k, s) {
			return false
		}
	}
	return true
}
