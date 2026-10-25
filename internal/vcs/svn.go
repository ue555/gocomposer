package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/Svn.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
)

const svnMaxAuthTries = 5

var (
	svnVersionMu sync.Mutex
	svnVersion   string
	reSvnVersion = regexp.MustCompile(`(\d+(?:\.\d+)+)`)
)

// Svn is Composer\Util\Svn.
type Svn struct {
	env              *Env
	url              string
	hasAuth          *bool
	username         string
	password         string
	cacheCredentials bool
	qtyAuthTries     int
}

// NewSvn is new Svn($url, $io, $config, $process).
func NewSvn(u string, env *Env) *Svn {
	return &Svn{env: env, url: u, cacheCredentials: true}
}

// SvnCleanEnv is Svn::cleanEnv.
func SvnCleanEnv() { os.Unsetenv("DYLD_LIBRARY_PATH") }

// SetCacheCredentials is Svn::setCacheCredentials.
func (s *Svn) SetCacheCredentials(b bool) { s.cacheCredentials = b }

// Execute is Svn::execute (a remote command: the url is checked against
// the secure-http config first).
func (s *Svn) Execute(cmd []string, u, cwd string, path *string, verbose bool) (string, error) {
	if err := ProhibitURL(s.env, u); err != nil {
		return "", err
	}
	return s.executeWithAuthRetry(cmd, cwd, u, path, verbose)
}

// ExecuteLocal is Svn::executeLocal.
func (s *Svn) ExecuteLocal(cmd []string, path, cwd string, verbose bool) (string, error) {
	return s.executeWithAuthRetry(cmd, cwd, "", &path, verbose)
}

func (s *Svn) executeWithAuthRetry(svnCmd []string, cwd, u string, path *string, verbose bool) (string, error) {
	command, err := s.command(svnCmd, u, path)
	if err != nil {
		return "", err
	}
	code, out, err := s.env.Process.Execute(command, cwd)
	if err != nil {
		return "", err
	}
	// the output handler drops the redirect notices and echoes the rest in
	// verbose mode
	var b strings.Builder
	for _, chunk := range strings.SplitAfter(out, "\n") {
		if chunk == "" || strings.HasPrefix(chunk, "Redirecting to URL ") {
			continue
		}
		b.WriteString(chunk)
	}
	output := b.String()
	if verbose && output != "" {
		s.env.IO.WriteErrorV(output, false, cio.Normal)
	}
	if code == 0 {
		return output, nil
	}
	full := strings.TrimSpace(output + "\n" + s.env.Process.ErrorOutput())
	lower := strings.ToLower(full)
	if !strings.Contains(lower, "could not authenticate to server:") &&
		!strings.Contains(lower, "authorization failed") &&
		!strings.Contains(lower, "svn: e170001:") &&
		!strings.Contains(lower, "svn: e215004:") {
		return "", &RuntimeError{Msg: full}
	}
	if !s.HasAuth() {
		if err := s.doAuthDance(); err != nil {
			return "", err
		}
	}
	s.qtyAuthTries++
	if s.qtyAuthTries-1 < svnMaxAuthTries {
		return s.executeWithAuthRetry(svnCmd, cwd, u, path, verbose)
	}
	return "", &RuntimeError{Msg: "wrong credentials provided (" + full + ")"}
}

// doAuthDance is Svn::doAuthDance.
func (s *Svn) doAuthDance() error {
	io := s.env.IO
	if io == nil || !io.Interactive {
		return &RuntimeError{Msg: "can not ask for authentication in non interactive mode"}
	}
	io.WriteError("The Subversion server (" + Sanitize(s.url) + ") requested credentials:")
	t := true
	s.hasAuth = &t
	user, err := io.Ask("Username: ", "")
	if err != nil {
		return err
	}
	pass, err := io.AskAndHideAnswer("Password: ")
	if err != nil {
		return err
	}
	s.username = phpStringCast(user)
	s.password = phpStringCast(pass)
	cache, err := io.AskConfirmation("Should Subversion cache these credentials? (yes/no) ", true)
	if err != nil {
		return err
	}
	s.cacheCredentials = cache
	return nil
}

// phpStringCast is (string) $value for an IO answer.
func phpStringCast(v any) string {
	s, _ := cio.AnswerString(v)
	return s
}

// command is Svn::getCommand.
func (s *Svn) command(cmd []string, u string, path *string) ([]string, error) {
	out := append(append([]string{}, cmd...), "--non-interactive")
	if s.HasAuth() {
		if !s.cacheCredentials {
			out = append(out, "--no-auth-cache")
		}
		out = append(out, "--username", s.username, "--password", s.password)
	}
	out = append(out, "--", u)
	if path != nil {
		out = append(out, *path)
	}
	return out, nil
}

// HasAuth is Svn::hasAuth: credentials come from the http-basic config of
// the repository host or from the url itself.
func (s *Svn) HasAuth() bool {
	if s.hasAuth != nil {
		return *s.hasAuth
	}
	ok := s.authFromConfig()
	if !ok {
		ok = s.authFromURL()
	}
	s.hasAuth = &ok
	return ok
}

func (s *Svn) authFromConfig() bool {
	if s.env.Config == nil {
		return false
	}
	auth := phpjson.ToMap(s.env.Config.Get("http-basic"))
	if auth == nil {
		return false
	}
	host := phpURLHost(s.url)
	c := auth.Map(host)
	if c == nil {
		return false
	}
	s.username, s.password = c.Str("username"), c.Str("password")
	return true
}

func (s *Svn) authFromURL() bool {
	pu, err := url.Parse(s.url)
	if err != nil || pu.User == nil || pu.User.Username() == "" || pu.User.Username() == "0" {
		return false
	}
	s.username = pu.User.Username()
	if p, ok := pu.User.Password(); ok && p != "" && p != "0" {
		s.password = p
	} else {
		s.password = ""
	}
	return true
}

// phpURLHost is parse_url($url, PHP_URL_HOST) ("" when absent).
func phpURLHost(u string) string {
	pu, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return pu.Hostname()
}

// BinaryVersion is Svn::binaryVersion ("" when svn is not available; the
// found version is kept for the whole process like the static property).
func (s *Svn) BinaryVersion() string {
	svnVersionMu.Lock()
	defer svnVersionMu.Unlock()
	if svnVersion == "" {
		if code, out, _ := s.env.Process.Execute([]string{"svn", "--version"}, ""); code == 0 {
			if m := reSvnVersion.FindStringSubmatch(out); m != nil {
				svnVersion = m[1]
			}
		}
	}
	return svnVersion
}

// ProhibitURL is Config::prohibitUrlByConfig for the VCS utilities,
// including the secure-svn-domains exception for svn:// urls.
func ProhibitURL(env *Env, u string) error {
	if env == nil || env.HTTP == nil {
		return nil
	}
	if env.Config != nil && strings.HasPrefix(strings.ToLower(u), "svn://") && env.Config.GetBool("secure-http") {
		host := phpURLHost(u)
		for _, d := range configList(env.Config, "secure-svn-domains") {
			if d == host {
				return nil
			}
		}
	}
	return env.HTTP.ProhibitURL(u)
}
