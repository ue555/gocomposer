package httpdl

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
)

// faultServer accepts connections and hands each one to handle.
func faultServer(t *testing.T, handle func(c net.Conn)) (string, *int32) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	var hits int32
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			atomic.AddInt32(&hits, 1)
			go handle(c)
		}
	}()
	return l.Addr().String(), &hits
}

func readRequestHead(c net.Conn) {
	r := bufio.NewReader(c)
	for {
		line, err := r.ReadString('\n')
		if err != nil || line == "\r\n" {
			return
		}
	}
}

func testDownloader(t *testing.T) *Downloader {
	cfg := config.New(false, t.TempDir())
	cfg.Set("secure-http", false)
	d, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestCurlErrors(t *testing.T) {
	cases := []struct {
		name   string
		handle func(c net.Conn)
		want   string
		hits   int32
	}{
		{"empty reply", func(c net.Conn) { readRequestHead(c); c.Close() }, "curl error 52 while downloading http://ADDR/p.json: Empty reply from server", 1},
		{"reset", func(c net.Conn) {
			readRequestHead(c)
			c.(*net.TCPConn).SetLinger(0)
			c.Close()
		}, "curl error 56 while downloading http://ADDR/p.json: Recv failure: Connection reset by peer", 4},
		{"http/0.9", func(c net.Conn) { readRequestHead(c); c.Write([]byte("this is not http\r\n")); c.Close() }, "curl error 1 while downloading http://ADDR/p.json: Received HTTP/0.9 when not allowed", 1},
		{"partial", func(c net.Conn) {
			readRequestHead(c)
			c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 100\r\nConnection: close\r\n\r\n{\"packages\""))
			c.Close()
		}, "curl error 18 while downloading http://ADDR/p.json: end of response with 89 bytes missing", 1},
		{"tls on plain", func(c net.Conn) {
			c.Read(make([]byte, 4096))
			c.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			c.Close()
		}, "curl error 35 while downloading https://ADDR/p.json: TLS connect error: error:0A00010B:SSL routines::wrong version number", 1},
		{"tls eof", func(c net.Conn) { c.Read(make([]byte, 4096)); c.Close() }, "curl error 35 while downloading https://ADDR/p.json: TLS connect error: error:0A000126:SSL routines::unexpected eof while reading", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addr, hits := faultServer(t, tc.handle)
			scheme := "http"
			if strings.Contains(tc.want, "https://") {
				scheme = "https"
			}
			_, err := testDownloader(t).Get(scheme + "://" + addr + "/p.json")
			want := strings.ReplaceAll(tc.want, "ADDR", addr)
			if err == nil || err.Error() != want {
				t.Fatalf("got %v\nwant %s", err, want)
			}
			if got := atomic.LoadInt32(hits); got != tc.hits {
				t.Errorf("attempts = %d, want %d", got, tc.hits)
			}
		})
	}
}

func TestCurlTLSCertificateErrors(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	_, err := testDownloader(t).Get(srv.URL + "/p.json")
	want := "curl error 60 while downloading " + srv.URL + "/p.json: SSL certificate problem: self-signed certificate"
	if err == nil || err.Error() != want {
		t.Fatalf("got %v\nwant %s", err, want)
	}
}

func TestCurlTimeout(t *testing.T) {
	addr, hits := faultServer(t, func(c net.Conn) { time.Sleep(3 * time.Second); c.Close() })
	opts := phpjson.NewMap()
	h := phpjson.NewMap()
	h.Set("method", "POST")
	h.Set("content", "a=b")
	h.Set("timeout", 0.3)
	opts.Set("http", h)
	_, err := testDownloader(t).Request("http://"+addr+"/api", opts)
	re := regexp.MustCompile(`^curl error 28 while downloading http://` + regexp.QuoteMeta(addr) + `/api: Operation timed out after \d+ milliseconds with 0 bytes received$`)
	if err == nil || !re.MatchString(err.Error()) {
		t.Fatalf("got %v", err)
	}
	if got := atomic.LoadInt32(hits); got != 1 {
		t.Errorf("a POST is not retried, got %d attempts", got)
	}
}

func TestCurlRetryable(t *testing.T) {
	for _, c := range []struct {
		method string
		errno  int
		msg    string
		want   bool
	}{
		{"GET", 28, "", true}, {"GET", 7, "", true}, {"GET", 52, "", false}, {"POST", 28, "", false},
		{"GET", 56, "Recv failure: Connection reset by peer", true}, {"GET", 56, "other", false},
		{"GET", 35, "Recv failure: Connection reset by peer", true}, {"GET", 60, "", false}, {"POST", 55, "", true},
	} {
		if got := curlRetryable(c.method, c.errno, c.msg); got != c.want {
			t.Errorf("curlRetryable(%s, %d, %q) = %v", c.method, c.errno, c.msg, got)
		}
	}
}
