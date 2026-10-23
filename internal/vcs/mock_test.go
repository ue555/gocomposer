package vcs

// Test doubles modelled on Composer's HttpDownloaderMock and
// ProcessExecutorMock (Composer f83588392df43d4fd027d84d26de459946ef4e80,
// tests/Composer/Test/Mock/*.php). Copyright (c) Nils Adermann, Jordi
// Boggiano. MIT License.

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
)

// Git::getVersion() is cached for the whole process like the PHP static;
// Composer's tests rely on an earlier real call having filled it, so the
// strict process mocks below never see `git --version`.
func init() {
	GitVersion(&Process{mock: func([]string, string) (int, string, string, error) { return 0, "git version 2.47.3\n", "", nil }})
}

// httpExpect is one HttpDownloaderMock expectation (status defaults to 200).
type httpExpect struct {
	url     string
	status  int
	body    string
	headers map[string]string
}

// httpMock serves the expectations in order (strict mode) and records
// every request it saw.
type httpMock struct {
	t        *testing.T
	mu       sync.Mutex
	expects  []httpExpect
	requests []*http.Request
}

func (m *httpMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests = append(m.requests, req)
	u := req.URL.String()
	if len(m.expects) == 0 {
		m.t.Errorf("unexpected request to %s (no expectation left)", u)
		return nil, fmt.Errorf("unexpected request %s", u)
	}
	e := m.expects[0]
	if e.url != u {
		m.t.Errorf("request to %s, expected %s", u, e.url)
		return nil, fmt.Errorf("unexpected request %s", u)
	}
	m.expects = m.expects[1:]
	status := e.status
	if status == 0 {
		status = 200
	}
	h := http.Header{}
	for k, v := range e.headers {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
		Header: h, Body: io.NopCloser(strings.NewReader(e.body)), ContentLength: int64(len(e.body)), Request: req}, nil
}

// assertComplete is the strict check that every expectation was consumed.
func (m *httpMock) assertComplete() {
	m.t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expects {
		m.t.Errorf("expected request to %s was not made", e.url)
	}
}

// procExpect is one ProcessExecutorMock expectation.
type procExpect struct {
	cmd    []string
	ret    int
	stdout string
	stderr string
}

type procMock struct {
	t       *testing.T
	expects []procExpect
	strict  bool
	log     []string
	// def is returned for commands without an expectation when not strict
	def procExpect
}

func (m *procMock) run(cmd []string, _ string) (int, string, string, error) {
	m.log = append(m.log, strings.Join(cmd, " "))
	// `git --version` only runs while Git::getVersion() is uncached; the Go
	// cache is primed (see init), so an expected call to it is optional
	if len(m.expects) > 1 && reflect.DeepEqual(m.expects[0].cmd, []string{"git", "--version"}) && !reflect.DeepEqual(cmd, m.expects[0].cmd) {
		m.expects = m.expects[1:]
	}
	if len(m.expects) > 0 && reflect.DeepEqual(m.expects[0].cmd, cmd) {
		e := m.expects[0]
		m.expects = m.expects[1:]
		return e.ret, e.stdout, e.stderr, nil
	}
	if m.strict {
		// ProcessExecutorMock throws an AssertionFailedError, which is an
		// \Exception: code that catches exceptions (getMirrorDefaultBranch)
		// swallows it, everything else fails with it
		want := "<none>"
		if len(m.expects) > 0 {
			want = strings.Join(m.expects[0].cmd, " ")
		}
		m.t.Logf("unexpected command %q, expected %s", strings.Join(cmd, " "), want)
		return 1, "", "", fmt.Errorf("Received unexpected command %q, expected %s", strings.Join(cmd, " "), want)
	}
	return m.def.ret, m.def.stdout, m.def.stderr, nil
}

func (m *procMock) assertComplete() {
	m.t.Helper()
	for _, e := range m.expects {
		m.t.Errorf("expected command %q was not run", strings.Join(e.cmd, " "))
	}
}

// testEnv builds a driver environment with a fresh home directory, the
// given config overrides (merged like Config::merge(['config' => ...])),
// mocked HTTP and process execution.
func testEnv(t *testing.T, interactive bool, cfg map[string]phpjson.Value, http []httpExpect, proc *procMock) (*Env, *httpMock, *bytes.Buffer) {
	t.Helper()
	home := t.TempDir()
	c := config.New(false, home)
	data := phpjson.NewMap()
	data.Set("home", home)
	data.Set("cache-dir", home+"/cache")
	data.Set("cache-repo-dir", home+"/cache/repo")
	data.Set("cache-vcs-dir", home+"/cache/vcs")
	for k, v := range cfg {
		data.Set(k, v)
	}
	wrap := phpjson.NewMap()
	wrap.Set("config", data)
	c.Merge(wrap)
	var out bytes.Buffer
	io := &cio.IO{Out: &out, Err: &out, Verbosity: cio.Normal, Interactive: interactive}
	dl, err := httpdl.New(c)
	if err != nil {
		t.Fatal(err)
	}
	hm := &httpMock{t: t, expects: http}
	dl.SetTransport(hm)
	p := NewProcess(io, 10*time.Second)
	if proc == nil {
		proc = &procMock{t: t, strict: true}
	}
	proc.t = t
	p.mock = proc.run
	return &Env{IO: io, Config: c, HTTP: dl, Process: p, Auth: NewAuthStore(c)}, hm, &out
}

// repoConfig is ['url' => $url] plus extra keys.
func repoConfig(url string, kv ...any) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("url", url)
	for i := 0; i+1 < len(kv); i += 2 {
		m.Set(kv[i].(string), kv[i+1])
	}
	return m
}

// mapGet walks nested keys of a decoded composer.json map.
func mapGet(m *phpjson.Map, keys ...string) (phpjson.Value, bool) {
	var v phpjson.Value = m
	for _, k := range keys {
		mm := phpjson.ToMap(v)
		if mm == nil {
			return nil, false
		}
		var ok bool
		if v, ok = mm.Get(k); !ok {
			return nil, false
		}
	}
	return v, true
}

func strOf(v phpjson.Value) string { s, _ := phpjson.ScalarString(v); return s }
