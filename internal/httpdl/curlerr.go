package httpdl

// The "curl error <errno> while downloading <url>: <curl_error()>" messages
// and the transport retry policy of Composer\Util\Http\CurlDownloader
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// src/Composer/Util/Http/CurlDownloader.php; MIT License, (c) Nils Adermann,
// Jordi Boggiano). Go reports network failures as error values, so the phase
// a transfer failed in (tracked with net/http/httptrace) and the error type
// are mapped to the libcurl error number and the error buffer text libcurl
// 8.14.1 with OpenSSL 3.5 (the oracle's PHP curl) writes for the same
// situation. Failures without a known libcurl counterpart keep the Go text
// and no number.

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Composer's CurlDownloader sets CURLOPT_CONNECTTIMEOUT 10 and
// CURLOPT_TIMEOUT max(default_socket_timeout, 300); a request's
// http.timeout option replaces the latter.
const (
	curlConnectTimeout = 10 * time.Second
	curlTimeout        = 300 * time.Second
	// CurlDownloader::$maxRetries
	curlMaxRetries = 3
)

// transfer follows one request through the phases libcurl distinguishes.
type transfer struct {
	mu        sync.Mutex
	start     time.Time
	connected bool // TCP (and TLS) connection established
	tlsStart  bool
	tlsDone   bool
	body      int64 // body bytes received (curl's bytecount)
	size      int64 // announced body size, -1 when unknown
}

func newTransfer() *transfer { return &transfer{start: time.Now(), size: -1} }

func (t *transfer) trace() *httptrace.ClientTrace {
	set := func(f func()) { t.mu.Lock(); f(); t.mu.Unlock() }
	return &httptrace.ClientTrace{
		GotConn:           func(httptrace.GotConnInfo) { set(func() { t.connected = true }) },
		TLSHandshakeStart: func() { set(func() { t.tlsStart = true }) },
		TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
			set(func() { t.tlsDone = err == nil })
		},
	}
}

func (t *transfer) addBody(n int64) { t.mu.Lock(); t.body += n; t.mu.Unlock() }

// countingWriter counts the body bytes written through it.
type countingWriter struct {
	w io.Writer
	t *transfer
}

func (c countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.t.addBody(int64(n))
	return n, err
}

// curlError maps a failed transfer to libcurl's error number and message
// (0 and the Go text when there is no libcurl counterpart). inBody is set
// for failures while the response body was read.
func (t *transfer) curlError(rawURL string, err error, inBody bool) (int, string) {
	t.mu.Lock()
	connected, tlsStart, tlsDone, body, size := t.connected, t.tlsStart, t.tlsDone, t.body, t.size
	elapsed := time.Since(t.start).Milliseconds()
	t.mu.Unlock()
	host, port := "", ""
	if u, perr := url.Parse(rawURL); perr == nil {
		host = u.Hostname()
		port = u.Port()
		if port == "" {
			port = map[string]string{"http": "80", "https": "443"}[u.Scheme]
		}
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) && !dnsErr.IsTimeout {
		return 6, "Could not resolve host: " + host
	}
	if isTimeout(err) {
		// multi.c multi_handle_timeout
		if !connected && !inBody {
			return 28, fmt.Sprintf("Connection timed out after %d milliseconds", elapsed)
		}
		if size >= 0 {
			return 28, fmt.Sprintf("Operation timed out after %d milliseconds with %d out of %d bytes received", elapsed, body, size)
		}
		return 28, fmt.Sprintf("Operation timed out after %d milliseconds with %d bytes received", elapsed, body)
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" && errors.Is(err, syscall.ECONNREFUSED) {
		return 7, "Failed to connect to " + host + " port " + port + " after 0 ms: Could not connect to server"
	}
	if tlsStart && !tlsDone && !inBody {
		if n, msg := tlsCurlError(host, err); n != 0 {
			return n, msg
		}
	}
	if errors.Is(err, syscall.ECONNRESET) {
		if tlsStart && !tlsDone && !inBody {
			return 35, "Recv failure: Connection reset by peer"
		}
		return 56, "Recv failure: Connection reset by peer"
	}
	if errors.Is(err, syscall.EPIPE) {
		return 55, "Send failure: Broken pipe"
	}
	if inBody && (errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF)) {
		// http.c / transfer.c: the connection closed before the body ended
		if size >= 0 {
			return 18, fmt.Sprintf("end of response with %d bytes missing", size-body)
		}
		return 18, "transfer closed with outstanding read data remaining"
	}
	if m := err.Error(); strings.Contains(m, "malformed HTTP ") && !strings.Contains(m, `"HTTP/`) {
		// a reply not starting with an HTTP/ status line is HTTP/0.9 to curl
		return 1, "Received HTTP/0.9 when not allowed"
	}
	if connected && (errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)) {
		return 52, "Empty reply from server"
	}
	return 0, Sanitize(err.Error())
}

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

// tlsCurlError maps a failed TLS handshake (OpenSSL backend of curl).
func tlsCurlError(host string, err error) (int, string) {
	var hostErr x509.HostnameError
	if errors.As(err, &hostErr) {
		// openssl.c ossl_verifyhost
		kind := "hostname"
		if ip := net.ParseIP(host); ip != nil {
			kind = "ipv4 address"
			if ip.To4() == nil {
				kind = "ipv6 address"
			}
		}
		c := hostErr.Certificate
		if c != nil && len(c.DNSNames) == 0 && len(c.IPAddresses) == 0 && c.Subject.CommonName != "" {
			return 60, "SSL: certificate subject name '" + c.Subject.CommonName + "' does not match target " + kind + " '" + host + "'"
		}
		return 60, "SSL: no alternative certificate subject name matches target " + kind + " '" + host + "'"
	}
	var invalid x509.CertificateInvalidError
	if errors.As(err, &invalid) && invalid.Reason == x509.Expired && invalid.Cert != nil {
		if time.Now().Before(invalid.Cert.NotBefore) {
			return 60, "SSL certificate problem: certificate is not yet valid"
		}
		return 60, "SSL certificate problem: certificate has expired"
	}
	var unknown x509.UnknownAuthorityError
	if errors.As(err, &unknown) {
		var chain []*x509.Certificate
		var verr *tls.CertificateVerificationError
		if errors.As(err, &verr) {
			chain = verr.UnverifiedCertificates
		}
		leaf := unknown.Cert
		if leaf == nil && len(chain) > 0 {
			leaf = chain[0]
		}
		switch {
		case leaf != nil && selfSigned(leaf):
			return 60, "SSL certificate problem: self-signed certificate"
		case len(chain) > 1 && selfSigned(chain[len(chain)-1]):
			return 60, "SSL certificate problem: self-signed certificate in certificate chain"
		}
		return 60, "SSL certificate problem: unable to get local issuer certificate"
	}
	var rec tls.RecordHeaderError
	if errors.As(err, &rec) || strings.Contains(err.Error(), "server gave HTTP response to HTTPS client") {
		return 35, "TLS connect error: error:0A00010B:SSL routines::wrong version number"
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return 35, "TLS connect error: error:0A000126:SSL routines::unexpected eof while reading"
	}
	return 0, ""
}

func selfSigned(c *x509.Certificate) bool {
	return bytes.Equal(c.RawIssuer, c.RawSubject) && c.CheckSignatureFrom(c) == nil
}

// curlRetryable is CurlDownloader's transport retry condition: GET requests
// failing with these errors are restarted up to $maxRetries times.
func curlRetryable(method string, errno int, msg string) bool {
	if errno == 55 {
		// CURLE_SEND_ERROR is retried whatever the method (curl#10591)
		return true
	}
	if method != "" && method != "GET" {
		return false
	}
	switch errno {
	case 7, 16, 92, 6, 28:
		return true
	case 56, 35:
		return strings.Contains(msg, "Connection reset by peer")
	}
	return false
}

// curlRetryDelay is CurlDownloader::restartJobWithDelay's default backoff
// for the given retry number.
func curlRetryDelay(retries int) time.Duration {
	switch {
	case retries >= 3:
		return 500 * time.Millisecond
	case retries >= 2:
		return 100 * time.Millisecond
	}
	return 0
}

// transportError builds the TransportException of a failed transfer.
func (t *transfer) transportError(rawURL string, err error, inBody bool) *TransportError {
	n, msg := t.curlError(rawURL, err, inBody)
	if n == 0 {
		return &TransportError{Msg: "curl error while downloading " + Sanitize(rawURL) + ": " + msg}
	}
	return &TransportError{Msg: fmt.Sprintf("curl error %d while downloading %s: %s", n, Sanitize(rawURL), msg), curlErrno: n, curlMsg: msg}
}

// curlMessage renders a Go network error like CurlDownloader does (for
// callers without a traced transfer).
func curlMessage(rawURL string, err error) string {
	return newTransfer().transportError(rawURL, err, false).Msg
}

// timeoutSeconds reads an http.timeout stream option (int or float).
func timeoutSeconds(v any) float64 {
	switch t := v.(type) {
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case float64:
		return t
	}
	return 0
}
