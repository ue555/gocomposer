// Package httpdl is the HTTP layer (Composer\Util\HttpDownloader, AuthHelper,
// Config::prohibitUrlByConfig and Url::sanitize). Credentials are only ever
// placed in request headers, never in messages.
package httpdl

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
)

// TransportError is Composer\Downloader\TransportException.
type TransportError struct {
	Msg        string
	StatusCode int
	// Header are the response headers (TransportException::getHeaders).
	Header http.Header
	// body is the response body of a failed POST (getResponse()).
	body []byte
	// curlErrno and curlMsg are the libcurl error of a failed transfer
	// (curlerr.go).
	curlErrno int
	curlMsg   string
}

func (e *TransportError) Error() string { return e.Msg }

// Response is a completed HTTP response.
type Response struct {
	URL        string
	StatusCode int
	Status     string
	Proto      string
	Header     http.Header
	Body       []byte
}

// Auth is one credential entry resolved for an origin.
type Auth struct {
	Username string
	Password string
}

// Downloader performs requests with Composer's configuration semantics.
type Downloader struct {
	cfg       *config.Config
	client    *http.Client
	auths     map[string]Auth
	mu        sync.Mutex
	warned    map[string]bool
	Warn      func(string)
	Debug     func(string)
	UserAgent string
	Timeout   time.Duration
	// MaxParallel bounds concurrent requests (Composer default 12).
	MaxParallel int
	// bitbucketRetry is AuthHelper::$bitbucketRetry.
	bitbucketRetry map[string]bool
	// IO asks for credentials when it is interactive (nil: never asks).
	IO *cio.IO
	// Sources are the config files credentials are stored in.
	Sources ConfigSources
	// promptMu serializes the authentication prompts of parallel requests.
	promptMu sync.Mutex
	// nullCreds marks prompted credentials answered with nothing (PHP null).
	nullCreds map[string][2]bool
}

// New builds a downloader from the config (auth, TLS, proxies from env).
func New(cfg *config.Config) (*Downloader, error) {
	d := &Downloader{cfg: cfg, auths: map[string]Auth{}, warned: map[string]bool{}, MaxParallel: 12}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment
	// CURLOPT_CONNECTTIMEOUT (curlerr.go)
	transport.DialContext = (&net.Dialer{Timeout: curlConnectTimeout, KeepAlive: 30 * time.Second}).DialContext
	transport.TLSHandshakeTimeout = curlConnectTimeout
	if cafile := cfg.GetString("cafile"); cafile != "" {
		pem, err := os.ReadFile(cafile)
		if err != nil {
			return nil, fmt.Errorf("could not read cafile %s: %w", cafile, err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		transport.TLSClientConfig = &tls.Config{RootCAs: pool}
	}
	if !cfg.GetBool("secure-http") && cfg.GetBool("disable-tls") {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // explicit disable-tls config
	}
	// CURLOPT_TIMEOUT; process-timeout only applies to processes
	timeout := curlTimeout
	d.Timeout = timeout
	d.client = &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 20 {
			return fmt.Errorf("too many redirects")
		}
		if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
			return &TransportError{Msg: `Could not follow the redirect to "` + Sanitize(req.URL.String()) + `" because only http and https redirects are supported.`}
		}
		// drop credentials when redirected to another host, like curl does
		if len(via) > 0 && req.URL.Host != via[0].URL.Host {
			req.Header.Del("Authorization")
			req.Header.Del("PRIVATE-TOKEN")
			req.Header.Del("JOB-TOKEN")
		}
		return nil
	}}
	d.loadAuth()
	return d, nil
}

// SetTransport replaces the HTTP transport; tests use it to serve canned
// responses the way Composer's HttpDownloaderMock does.
func (d *Downloader) SetTransport(rt http.RoundTripper) { d.client.Transport = rt }

func (d *Downloader) loadAuth() {
	cfg := d.cfg
	// BaseIO::loadConfiguration: Bitbucket OAuth consumers first, the
	// other credential kinds override them
	if m := phpjson.ToMap(cfg.Get("bitbucket-oauth")); m != nil {
		for _, host := range m.Keys() {
			if cred := m.Map(host); cred != nil {
				d.auths[host] = Auth{cred.Str("consumer-key"), cred.Str("consumer-secret")}
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("http-basic")); m != nil {
		for _, host := range m.Keys() {
			cred := m.Map(host)
			if cred == nil {
				continue
			}
			d.auths[host] = Auth{cred.Str("username"), cred.Str("password")}
		}
	}
	if m := phpjson.ToMap(cfg.Get("bearer")); m != nil {
		for _, host := range m.Keys() {
			d.auths[host] = Auth{m.Str(host), "bearer"}
		}
	}
	if m := phpjson.ToMap(cfg.Get("github-oauth")); m != nil {
		for _, host := range m.Keys() {
			d.auths[host] = Auth{m.Str(host), "x-oauth-basic"}
		}
	}
	if m := phpjson.ToMap(cfg.Get("gitlab-oauth")); m != nil {
		for _, host := range m.Keys() {
			v, _ := m.Get(host)
			if s, ok := v.(string); ok {
				d.auths[host] = Auth{s, "oauth2"}
			} else if mm := phpjson.ToMap(v); mm != nil {
				d.auths[host] = Auth{mm.Str("token"), "oauth2"}
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("gitlab-token")); m != nil {
		for _, host := range m.Keys() {
			v, _ := m.Get(host)
			if s, ok := v.(string); ok {
				d.auths[host] = Auth{s, "private-token"}
			} else if mm := phpjson.ToMap(v); mm != nil {
				d.auths[host] = Auth{mm.Str("username"), mm.Str("token")}
			}
		}
	}
	if m := phpjson.ToMap(cfg.Get("custom-headers")); m != nil {
		for _, host := range m.Keys() {
			v, _ := m.Get(host)
			enc, _ := phpjson.Encode(v, 0)
			d.auths[host] = Auth{enc, "custom-headers"}
		}
	}
}

var (
	reAccessToken = regexp.MustCompile(`([&?]access_token=)[^&]+`)
	reUserInfo    = regexp.MustCompile(`(?i)(?:([a-z0-9][a-z0-9+.-]*://)|\A)([^:/\s?#]*)(?::([^\s/?#]*))?@`)
	reGHToken     = regexp.MustCompile(`^([a-f0-9]{12,}|gh[a-z]_[a-zA-Z0-9_]+|github_pat_[a-zA-Z0-9_]+)$`)
)

// Sanitize is Url::sanitize: hides passwords and tokens embedded in URLs.
func Sanitize(u string) string {
	u = reAccessToken.ReplaceAllString(u, "${1}***")
	return reUserInfo.ReplaceAllStringFunc(u, func(m string) string {
		sm := reUserInfo.FindStringSubmatch(m)
		user := sanitizeUsername(sm[2])
		if sm[3] != "" {
			return sm[1] + user + ":***@"
		}
		return sm[1] + user + "@"
	})
}

func sanitizeUsername(u string) string {
	switch u {
	case "x-token-auth", "x-oauth-basic", "oauth2", "gitlab-ci-token", "private-token", "token", "x-access-token":
		return u
	}
	if reGHToken.MatchString(u) || len(u) >= 12 {
		if len(u) > 3 {
			return u[:3] + "***"
		}
	}
	return u
}

var (
	reSchemeTransport = regexp.MustCompile(`(?i)^[a-z][a-z0-9+.-]*::(.+)$`)
	reScheme          = regexp.MustCompile(`(?i)^([a-z][a-z0-9+.-]*)://`)
)

// ProhibitURL is Config::prohibitUrlByConfig.
func (d *Downloader) ProhibitURL(u string) error {
	checked := u
	for {
		m := reSchemeTransport.FindStringSubmatch(checked)
		if m == nil {
			break
		}
		checked = m[1]
	}
	m := reScheme.FindStringSubmatch(checked)
	if m == nil {
		return nil
	}
	scheme := strings.ToLower(m[1])
	host := ""
	if pu, err := url.Parse(checked); err == nil {
		host = pu.Hostname()
	}
	switch scheme {
	case "http", "git", "ftp", "svn":
		if d.cfg.GetBool("secure-http") {
			if scheme == "svn" {
				return &TransportError{Msg: "Your configuration does not allow connections to " + Sanitize(u) + ". See https://getcomposer.org/doc/06-config.md#secure-svn-domains for details."}
			}
			return &TransportError{Msg: "Your configuration does not allow connections to " + Sanitize(u) + ". See https://getcomposer.org/doc/06-config.md#secure-http for details."}
		}
		if host != "" {
			d.mu.Lock()
			warned := d.warned[host]
			d.warned[host] = true
			d.mu.Unlock()
			if !warned && d.Warn != nil {
				d.Warn("<warning>Warning: Accessing " + host + " over " + scheme + " which is an insecure protocol.</warning>")
			}
		}
	}
	return nil
}

func (d *Downloader) origin(u *url.URL) string {
	host := u.Hostname()
	if p := u.Port(); p != "" && !((u.Scheme == "https" && p == "443") || (u.Scheme == "http" && p == "80")) {
		host += ":" + p
	}
	return host
}

func (d *Downloader) addAuth(req *http.Request) {
	origin := d.origin(req.URL)
	auth, ok := d.auths[origin]
	if !ok {
		// AuthHelper::findAuthOrigin: github.com / bitbucket.org credentials
		// apply to their api. hosts
		if origin == "api.github.com" || origin == "api.bitbucket.org" {
			canonical := strings.TrimPrefix(origin, "api.")
			if auth, ok = d.auths[canonical]; ok {
				origin = canonical
			}
		}
	}
	if req.URL.User != nil {
		pw, _ := req.URL.User.Password()
		auth, ok = Auth{req.URL.User.Username(), pw}, true
		req.URL.User = nil
	}
	if !ok {
		return
	}
	switch {
	case auth.Password == "bearer":
		req.Header.Add("Authorization", "Bearer "+auth.Username)
	case auth.Password == "custom-headers":
		v, err := phpjson.Decode([]byte(auth.Username))
		if err == nil {
			phpjson.Each(v, func(_ string, h phpjson.Value) {
				s, _ := h.(string)
				if k, val, ok := strings.Cut(s, ":"); ok {
					req.Header.Add(strings.TrimSpace(k), strings.TrimSpace(val))
				}
			})
		}
	case origin == "github.com" && auth.Password == "x-oauth-basic":
		if strings.HasPrefix(req.URL.String(), "https://api.github.com/") || strings.HasPrefix(req.URL.String(), "http://api.github.com/") {
			req.Header.Add("Authorization", "token "+auth.Username)
		}
	case (auth.Username == "gitlab-ci-token" || auth.Password == "gitlab-ci-token") && d.isGitlab(origin):
		tok := auth.Password
		if auth.Password == "gitlab-ci-token" {
			tok = auth.Username
		}
		req.Header.Add("JOB-TOKEN", tok)
	case auth.Password == "oauth2" && d.isGitlab(origin):
		req.Header.Add("Authorization", "Bearer "+auth.Username)
	case auth.Password == "private-token" && d.isGitlab(origin):
		req.Header.Add("PRIVATE-TOKEN", auth.Username)
	case origin == "bitbucket.org" && auth.Username == "x-token-auth" && req.URL.String() != "https://bitbucket.org/site/oauth2/access_token":
		if !isPublicBitbucketDownload(req.URL) {
			req.Header.Add("Authorization", "Bearer "+auth.Password)
		}
	default:
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth.Username+":"+auth.Password)))
	}
}

// isPublicBitbucketDownload is AuthHelper::isPublicBitBucketDownload.
func isPublicBitbucketDownload(u *url.URL) bool {
	if !strings.Contains(u.Hostname(), "bitbucket.org") {
		return true
	}
	parts := strings.Split(u.Path, "/")
	return len(parts) >= 4 && parts[3] == "downloads"
}

func (d *Downloader) isGitlab(origin string) bool {
	found := false
	phpjson.Each(d.cfg.Get("gitlab-domains"), func(_ string, v phpjson.Value) {
		if s, _ := v.(string); s == origin {
			found = true
		}
	})
	return found
}

// Get fetches a URL fully into memory. extraHeaders are appended verbatim.
func (d *Downloader) Get(rawURL string, extraHeaders ...string) (*Response, error) {
	return d.get(context.Background(), rawURL, "", extraHeaders)
}

// Download streams a URL into the file at dest.
func (d *Downloader) Download(rawURL, dest string, extraHeaders ...string) (*Response, error) {
	return d.get(context.Background(), rawURL, dest, extraHeaders)
}

// Request performs a request described by PHP stream-context style options
// (HttpDownloader::get($url, $options)): http.method, http.header (string or
// list) and http.content are honoured.
func (d *Downloader) Request(rawURL string, options *phpjson.Map) (*Response, error) {
	ctx := context.Background()
	method := http.MethodGet
	var headers []string
	var body []byte
	if options != nil {
		if hm := options.Map("http"); hm != nil {
			if m, ok := hm.Get("method"); ok {
				if s, ok := m.(string); ok && s != "" {
					method = strings.ToUpper(s)
				}
			}
			if h, ok := hm.Get("header"); ok {
				switch t := h.(type) {
				case string:
					headers = append(headers, strings.Split(strings.ReplaceAll(t, "\r\n", "\n"), "\n")...)
				default:
					phpjson.Each(t, func(_ string, v phpjson.Value) {
						if s, ok := v.(string); ok {
							headers = append(headers, s)
						}
					})
				}
			}
			if c, ok := hm.Get("content"); ok {
				if s, ok := c.(string); ok {
					body = []byte(s)
				}
			}
			// http.timeout replaces CURLOPT_TIMEOUT for this request
			if t, ok := hm.Get("timeout"); ok {
				if secs := timeoutSeconds(t); secs > 0 {
					c, cancel := context.WithTimeout(ctx, time.Duration(secs*float64(time.Second)))
					defer cancel()
					ctx = c
				}
			}
		}
	}
	if method == http.MethodGet && body == nil {
		return d.get(ctx, rawURL, "", headers)
	}
	if err := d.ProhibitURL(rawURL); err != nil {
		return nil, err
	}
	resp, err, retry := d.doMethod(ctx, method, rawURL, "", headers, body)
	// CURLE_SEND_ERROR is retried whatever the method
	for retries := 1; err != nil && retry && retries <= curlMaxRetries; retries++ {
		if te, ok := err.(*TransportError); !ok || te.curlErrno != 55 {
			break
		}
		d.debugRetry(rawURL, retries, err)
		time.Sleep(curlRetryDelay(retries))
		resp, err, retry = d.doMethod(ctx, method, rawURL, "", headers, body)
	}
	return resp, err
}

func (d *Downloader) get(ctx context.Context, rawURL, dest string, extraHeaders []string) (*Response, error) {
	if strings.HasPrefix(rawURL, "file://") || isLocalPathURL(rawURL) {
		return d.getFile(rawURL, dest)
	}
	if err := d.ProhibitURL(rawURL); err != nil {
		return nil, err
	}
	var lastErr error
	authRetries := 0
	var storeAuth any = false
	// CurlDownloader: a failed GET is restarted at most $maxRetries times,
	// after 0, 0.1 and 0.5 seconds
	for retries := 0; ; retries++ {
		resp, err, retry := d.doRequest(ctx, rawURL, dest, extraHeaders)
		if again, store, aerr := d.checkAuthFailure(ctx, rawURL, dest, resp, err, authRetries); aerr != nil {
			return nil, aerr
		} else if again {
			authRetries++
			retries--
			storeAuth = store
			continue
		}
		if err == nil {
			// CurlDownloader: credentials entered for this request are
			// stored once it succeeded
			if storeAuth != false {
				if serr := d.StoreAuth(d.Origin(rawURL), storeAuth); serr != nil {
					return nil, serr
				}
			}
			return resp, nil
		}
		lastErr = err
		if !retry || retries >= curlMaxRetries {
			break
		}
		d.debugRetry(rawURL, retries+1, err)
		time.Sleep(curlRetryDelay(retries + 1))
	}
	return nil, lastErr
}

// debugRetry writes CurlDownloader's "Retrying (N) <url> due to ..." line.
func (d *Downloader) debugRetry(rawURL string, retries int, err error) {
	if d.IO == nil {
		return
	}
	te, _ := err.(*TransportError)
	switch {
	case te != nil && te.curlErrno != 0:
		d.IO.WriteErrorV(fmt.Sprintf("Retrying (%d) %s due to curl error %d", retries, Sanitize(rawURL), te.curlErrno), true, cio.Debug)
	case te != nil && te.StatusCode != 0:
		d.IO.WriteErrorV(fmt.Sprintf("Retrying (%d) %s due to status code %d", retries, Sanitize(rawURL), te.StatusCode), true, cio.Debug)
	}
}

func (d *Downloader) doRequest(ctx context.Context, rawURL, dest string, extraHeaders []string) (*Response, error, bool) {
	return d.doMethod(ctx, http.MethodGet, rawURL, dest, extraHeaders, nil)
}

func (d *Downloader) doMethod(ctx context.Context, method, rawURL, dest string, extraHeaders []string, body []byte) (*Response, error, bool) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, reader)
	if err != nil {
		return nil, &TransportError{Msg: `The "` + Sanitize(rawURL) + `" file could not be downloaded (` + err.Error() + `)`}, false
	}
	ua := d.UserAgent
	if ua == "" {
		ua = "gocomposer (Composer-compatible; Go)"
	}
	req.Header.Set("User-Agent", ua)
	for _, h := range extraHeaders {
		if k, v, ok := strings.Cut(h, ":"); ok {
			req.Header.Add(strings.TrimSpace(k), strings.TrimSpace(v))
		}
	}
	d.addAuth(req)
	tr := newTransfer()
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), tr.trace()))
	res, err := d.client.Do(req)
	if err != nil {
		var te *TransportError
		if ue, ok := err.(*url.Error); ok {
			if t, ok := ue.Err.(*TransportError); ok {
				te = t
			}
		}
		if te != nil {
			return nil, te, false
		}
		te = tr.transportError(rawURL, err, false)
		return nil, te, curlRetryable(method, te.curlErrno, te.curlMsg)
	}
	defer res.Body.Close()
	tr.size = res.ContentLength
	out := &Response{URL: res.Request.URL.String(), StatusCode: res.StatusCode, Status: res.Status, Proto: res.Proto, Header: res.Header}
	if res.StatusCode >= 300 && res.StatusCode != 304 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
		out.Body = body
		reason := strings.TrimSpace(strings.TrimPrefix(res.Status, fmt.Sprint(res.StatusCode)))
		proto := strings.TrimPrefix(res.Proto, "HTTP/")
		if proto == "2.0" {
			proto = "2"
		}
		details := ""
		// CurlDownloader: the warnings/infos of a JSON error response are
		// shown (ResponseWarnings) instead of dumping the raw body
		warned := false
		if isJSONContentType(res.Header.Get("Content-Type")) {
			if v, derr := phpjson.Decode(body); derr == nil {
				warned = d.OutputResponseWarnings(d.Origin(rawURL), v)
			}
		}
		if !warned {
			details = jsonFailureDetails(res.Header.Get("Content-Type"), body)
		}
		msg := fmt.Sprintf(`The "%s" file could not be downloaded (HTTP/%s %d %s)`, Sanitize(rawURL), proto, res.StatusCode, reason) + details
		return out, &TransportError{Msg: msg, StatusCode: res.StatusCode, Header: res.Header}, res.StatusCode >= 500
	}
	if dest != "" && res.StatusCode == 200 {
		if err := os.MkdirAll(filepath.Dir(dest), 0o777); err != nil {
			return nil, &TransportError{Msg: `The "` + Sanitize(rawURL) + `" file could not be written to ` + dest + `: ` + err.Error()}, false
		}
		tmp := dest + ".part"
		f, err := os.Create(tmp)
		if err != nil {
			return nil, &TransportError{Msg: `The "` + Sanitize(rawURL) + `" file could not be written to ` + dest + `: ` + err.Error()}, false
		}
		n, err := io.Copy(countingWriter{f, tr}, res.Body)
		cerr := f.Close()
		if err != nil || cerr != nil {
			os.Remove(tmp)
			if err == nil {
				err = cerr
			}
			te := tr.transportError(rawURL, err, true)
			return nil, te, curlRetryable(method, te.curlErrno, te.curlMsg)
		}
		if res.ContentLength >= 0 && n != res.ContentLength {
			os.Remove(tmp)
			return nil, &TransportError{Msg: fmt.Sprintf("Content-Length mismatch, received %d bytes out of the expected %d for %s", n, res.ContentLength, Sanitize(rawURL))}, true
		}
		if err := os.Rename(tmp, dest); err != nil {
			os.Remove(tmp)
			return nil, &TransportError{Msg: err.Error()}, false
		}
		return out, nil, false
	}
	var buf bytes.Buffer
	_, err = io.Copy(countingWriter{&buf, tr}, res.Body)
	respBody := buf.Bytes()
	if err != nil {
		te := tr.transportError(rawURL, err, true)
		return nil, te, curlRetryable(method, te.curlErrno, te.curlMsg)
	}
	out.Body = respBody
	return out, nil, false
}

func (d *Downloader) getFile(rawURL, dest string) (*Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, &TransportError{Msg: `The "` + rawURL + `" file could not be downloaded: ` + err.Error()}
	}
	p := u.Path
	if !strings.HasPrefix(rawURL, "file://") {
		// a plain filesystem path (e.g. an artifact repository dist url),
		// read by RemoteFilesystem's file_get_contents in Composer
		p = rawURL
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if p == rawURL {
			// RemoteFilesystem::get: the PHP 8 stream warning, error code 0
			return nil, &TransportError{Msg: `The "` + rawURL + `" file could not be downloaded: Failed to open stream: No such file or directory`}
		}
		code := 0
		if os.IsNotExist(err) {
			code = 404
		}
		return nil, &TransportError{Msg: `The "` + rawURL + `" file could not be downloaded: failed to open stream: No such file or directory`, StatusCode: code}
	}
	if dest != "" {
		if err := os.MkdirAll(filepath.Dir(dest), 0o777); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dest, data, 0o666); err != nil {
			return nil, err
		}
		return &Response{URL: rawURL, StatusCode: 200}, nil
	}
	return &Response{URL: rawURL, StatusCode: 200, Body: data, Header: http.Header{}}, nil
}

// jsonFailureDetails is the body excerpt CurlDownloader::failResponse appends
// to the message of a failed JSON response (Composer
// f83588392df43d4fd027d84d26de459946ef4e80). It is not added when
// ResponseWarnings presented the body's warnings/infos (see doMethod).
func jsonFailureDetails(contentType string, body []byte) string {
	if !isJSONContentType(contentType) {
		return ""
	}
	if len(body) > 200 {
		return ":\n" + string(body[:200]) + "..."
	}
	return ":\n" + string(body)
}

// isLocalPathURL reports a url without a scheme (a plain local path), which
// Composer's HttpDownloader hands to RemoteFilesystem instead of curl.
func isLocalPathURL(u string) bool {
	return u != "" && !strings.Contains(u, "://")
}

// Job is one entry for GetAll.
type Job struct {
	URL     string
	Dest    string
	Headers []string
}

// Result pairs a Job with its outcome.
type Result struct {
	Resp *Response
	Err  error
}

// GetAll runs jobs concurrently and returns results in job order.
func (d *Downloader) GetAll(jobs []Job) []Result {
	results := make([]Result, len(jobs))
	sem := make(chan struct{}, max(1, d.MaxParallel))
	var wg sync.WaitGroup
	for i, j := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, j Job) {
			defer wg.Done()
			defer func() { <-sem }()
			resp, err := d.get(context.Background(), j.URL, j.Dest, j.Headers)
			results[i] = Result{resp, err}
		}(i, j)
	}
	wg.Wait()
	return results
}
