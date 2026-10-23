package httpdl

// AuthHelper::promptAuthIfNeeded and the
// auth-failure checks of CurlDownloader::isAuthenticatedRetryNeeded, plus
// Url::getOrigin and the git-config token lookups of GitHub/GitLab/Bitbucket
// ::authorizeOAuth. Ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/Util/AuthHelper.php,
// src/Composer/Util/Http/CurlDownloader.php, src/Composer/Util/Url.php,
// src/Composer/Util/GitHub.php, GitLab.php, Bitbucket.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
//
// The interactive authorizations (GitHub/GitLab/Bitbucket::
// authorizeOAuthInteractively, Bitbucket::requestToken) and
// AuthHelper::storeAuth are in interactive.go.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/phpjson"
)

type ctxKey int

const noAuthRetryKey ctxKey = 1

// GetNoAuthRetry is HttpDownloader::get($url, ['retry-auth-failure' => false]).
func (d *Downloader) GetNoAuthRetry(rawURL string, extraHeaders ...string) (*Response, error) {
	return d.get(context.WithValue(context.Background(), noAuthRetryKey, true), rawURL, "", extraHeaders)
}

func authRetryDisabled(ctx context.Context) bool {
	v, _ := ctx.Value(noAuthRetryKey).(bool)
	return v
}

var rePortSuffix = regexp.MustCompile(`^([^/]+):\d+`)

// Origin is Url::getOrigin.
func (d *Downloader) Origin(rawURL string) string {
	if strings.HasPrefix(rawURL, "file://") {
		return rawURL
	}
	origin := ""
	if u, err := url.Parse(rawURL); err == nil {
		origin = u.Hostname()
		if p := u.Port(); p != "" {
			origin += ":" + p
		}
	}
	if strings.HasSuffix(origin, ".github.com") && origin != "codeload.github.com" {
		return "github.com"
	}
	if origin == "repo.packagist.org" {
		return "packagist.org"
	}
	if origin == "" {
		origin = rawURL
	}
	domains := configStrings(d.cfg.Get("gitlab-domains"))
	if !strings.Contains(origin, "/") && !containsStr(domains, origin) {
		for _, gd := range domains {
			bc := rePortSuffix.ReplaceAllString(gd, "$1")
			if gd != "" && (bc == origin || strings.HasPrefix(bc, origin+"/")) {
				return gd
			}
		}
	}
	return origin
}

func configStrings(v phpjson.Value) []string {
	var out []string
	phpjson.Each(v, func(_ string, e phpjson.Value) {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	})
	return out
}

func containsStr(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func (d *Downloader) hasAuth(origin string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.auths[origin]
	return ok
}

func (d *Downloader) getAuth(origin string) (Auth, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	a, ok := d.auths[origin]
	return a, ok
}

func (d *Downloader) setAuth(origin, user, pass string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.auths[origin] = Auth{Username: user, Password: pass}
}

// gitConfig runs `git config <key>` like the authorizeOAuth helpers
// (ProcessExecutor::execute, trimmed output, ok when the exit code is 0).
func gitConfig(key string) (string, bool) {
	out, err := exec.Command("git", "config", key).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), true
}

func (d *Downloader) authorizeGitHub(origin string) bool {
	if !containsStr(configStrings(d.cfg.Get("github-domains")), origin) {
		return false
	}
	if tok, ok := gitConfig("github.accesstoken"); ok {
		d.setAuth(origin, tok, "x-oauth-basic")
		return true
	}
	return false
}

func (d *Downloader) authorizeGitLab(origin string) bool {
	bc := regexp.MustCompile(`:\d+`).ReplaceAllString(origin, "")
	domains := configStrings(d.cfg.Get("gitlab-domains"))
	if !containsStr(domains, origin) && !containsStr(domains, bc) {
		return false
	}
	if tok, ok := gitConfig("gitlab.accesstoken"); ok {
		d.setAuth(origin, tok, "oauth2")
		return true
	}
	if user, ok := gitConfig("gitlab.deploytoken.user"); ok {
		if tok, ok := gitConfig("gitlab.deploytoken.token"); ok {
			d.setAuth(origin, user, tok)
			return true
		}
	}
	tokens := phpjson.ToMap(d.cfg.Get("gitlab-token"))
	if tokens == nil {
		return false
	}
	var tok phpjson.Value
	found := false
	for _, k := range []string{origin, bc} {
		if v, ok := tokens.Get(k); ok && v != nil {
			tok, found = v, true
			break
		}
	}
	if !found {
		return false
	}
	user, pass := "", "private-token"
	if m := phpjson.ToMap(tok); m != nil && phpjson.IsArray(tok) {
		user, pass = m.Str("username"), m.Str("token")
	} else {
		user, _ = phpjson.ScalarString(tok)
	}
	if user == "private-token" || user == "gitlab-ci-token" || user == "oauth2" {
		user, pass = pass, user
	}
	d.setAuth(origin, user, pass)
	return true
}

func (d *Downloader) authorizeBitbucket(origin string) bool {
	if origin != "bitbucket.org" {
		return false
	}
	if tok, ok := gitConfig("bitbucket.accesstoken"); ok {
		d.setAuth(origin, "x-token-auth", tok)
		return true
	}
	return false
}

// promptAuth is AuthHelper::promptAuthIfNeeded. It reports whether the
// request is to be retried and the storeAuth value of the retry (false, true
// or "prompt"); a non-nil error replaces the failed response. resp is nil
// for the auth failures detected from the response content.
func (d *Downloader) promptAuth(rawURL, origin string, status int, reason string, resp *Response, retryCount int) (bool, any, error) {
	d.promptMu.Lock()
	defer d.promptMu.Unlock()
	switch {
	case containsStr(configStrings(d.cfg.Get("github-domains")), origin):
		if d.hasAuth(origin) && retryCount == 0 {
			return true, false, nil
		}
		var header http.Header
		if resp != nil {
			header = resp.Header
		}
		message := "\n"
		if gitHubRequiresSso(header) {
			message = "GitHub API token requires SSO authorization. Authorize this token at " + Sanitize(gitHubSsoURL(header)) + "\n"
			d.writeError(message)
			if !d.interactive() {
				return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 403}
			}
			if _, err := d.IO.Ask("After authorizing your token, confirm that you would like to retry the request", nil); err != nil {
				return false, false, err
			}
			return true, false, nil
		}
		if GitHubIsRateLimited(header) {
			limit, reset := GitHubRateLimit(header)
			if d.hasAuth(origin) {
				message = "Review your configured GitHub OAuth token or enter a new one to go over the API rate limit."
			} else {
				message = "Create a GitHub OAuth token to go over the API rate limit."
			}
			message = fmt.Sprintf("GitHub API limit (%d calls/hr) is exhausted, could not fetch "+Sanitize(rawURL)+". "+message+" You can also wait until %s for the rate limit to reset.", limit, reset) + "\n"
		} else {
			apiMessage := ""
			if resp != nil && resp.Body != nil {
				if v, err := phpjson.Decode(resp.Body); err == nil {
					if m := phpjson.ToMap(v); m != nil {
						if mv, ok := m.Get("message"); ok {
							apiMessage, _ = mv.(string)
						}
					}
				}
			}
			if apiMessage != "" {
				message += "Could not fetch " + Sanitize(rawURL) + ": " + apiMessage
			} else {
				message += "Could not fetch " + Sanitize(rawURL) + ", please "
				if d.hasAuth(origin) {
					message += "review your configured GitHub OAuth token or enter a new one to access private repos"
				} else {
					message += "create a GitHub OAuth token to access private repos"
				}
			}
		}
		if !d.authorizeGitHub(origin) {
			if !d.interactive() {
				return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
			}
			ok, err := d.GitHubAuthorizeInteractively(origin, message)
			if err != nil {
				return false, false, err
			}
			if !ok {
				return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
			}
		}
		return true, false, nil
	case containsStr(configStrings(d.cfg.Get("gitlab-domains")), origin):
		message := "\nCould not fetch " + Sanitize(rawURL) + ", enter your " + origin + " credentials "
		if status == 401 {
			message += "to access private repos"
		} else {
			message += "to go over the API rate limit"
		}
		auth, had := d.getAuth(origin)
		if had && (auth.Password == "gitlab-ci-token" || auth.Password == "private-token" || auth.Password == "oauth2" || auth.Username == "gitlab-ci-token") {
			return false, false, &TransportError{Msg: "Invalid credentials for '" + Sanitize(rawURL) + "', aborting.", StatusCode: status}
		}
		if !d.authorizeGitLab(origin) {
			if !d.interactive() {
				return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
			}
			scheme := ""
			if u, err := url.Parse(rawURL); err == nil {
				scheme = u.Scheme
			}
			ok, err := d.GitLabAuthorizeInteractively(scheme, origin, message)
			if err != nil {
				return false, false, err
			}
			if !ok {
				return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
			}
		}
		if now, ok := d.getAuth(origin); had && ok && now == auth {
			return false, false, &TransportError{Msg: "Invalid credentials for '" + Sanitize(rawURL) + "', aborting.", StatusCode: status}
		}
		return true, false, nil
	case origin == "bitbucket.org" || origin == "api.bitbucket.org":
		askForOAuthToken := true
		origin = "bitbucket.org"
		if auth, ok := d.getAuth(origin); ok {
			if auth.Username != "x-token-auth" {
				token, err := NewBitbucket(d).RequestToken(origin, auth.Username, auth.Password)
				if err != nil {
					return false, false, err
				}
				if token != "" {
					d.setAuth(origin, "x-token-auth", token)
					askForOAuthToken = false
				}
			} else {
				d.mu.Lock()
				if d.bitbucketRetry == nil {
					d.bitbucketRetry = map[string]bool{}
				}
				retried := d.bitbucketRetry[rawURL]
				d.bitbucketRetry[rawURL] = true
				d.mu.Unlock()
				if retried {
					return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
				}
				askForOAuthToken = false
			}
		}
		if askForOAuthToken {
			message := "\nCould not fetch " + Sanitize(rawURL) + ", please create a bitbucket OAuth token to "
			if status == 401 || status == 403 {
				message += "access private repos"
			} else {
				message += "go over the API rate limit"
			}
			if !d.authorizeBitbucket(origin) {
				if !d.interactive() {
					return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
				}
				ok, err := NewBitbucket(d).AuthorizeOAuthInteractively(origin, message)
				if err != nil {
					return false, false, err
				}
				if !ok {
					return false, false, &TransportError{Msg: "Could not authenticate against " + origin, StatusCode: 401}
				}
			}
		}
		return true, false, nil
	}
	if status == 404 {
		return false, false, nil
	}
	if !d.interactive() {
		switch status {
		case 401:
			return false, false, &TransportError{Msg: "The '" + Sanitize(rawURL) + "' URL required authentication (HTTP 401).\nYou must be using the interactive console to authenticate", StatusCode: status}
		case 403:
			return false, false, &TransportError{Msg: "The '" + Sanitize(rawURL) + "' URL could not be accessed (HTTP 403): " + reason, StatusCode: status}
		}
		return false, false, &TransportError{Msg: "Unknown error code '" + strconv.Itoa(status) + "', reason: " + reason, StatusCode: status}
	}
	if d.hasAuth(origin) {
		if retryCount == 0 {
			return true, false, nil
		}
		return false, false, &TransportError{Msg: "Invalid credentials (HTTP " + strconv.Itoa(status) + ") for '" + Sanitize(rawURL) + "', aborting.", StatusCode: status}
	}
	d.IO.WriteError("    Authentication required (<info>" + origin + "</info>):")
	user, err := d.IO.Ask("      Username: ", nil)
	if err != nil {
		return false, false, err
	}
	pass, err := d.IO.AskAndHideAnswer("      Password: ")
	if err != nil {
		return false, false, err
	}
	d.SetPromptedAuth(origin, user, pass)
	return true, configStoreAuths(d.cfg.Get("store-auths")), nil
}

// configStoreAuths is $config->get('store-auths'): true, false or "prompt".
func configStoreAuths(v phpjson.Value) any {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		if t == "prompt" {
			return "prompt"
		}
		return t != "" && t != "0"
	}
	return phpjson.Truthy(v)
}

// statusLine is Response::getStatusMessage ("HTTP/1.1 403 Forbidden").
func statusLine(r *Response) string {
	proto := r.Proto
	if proto == "HTTP/2.0" {
		proto = "HTTP/2"
	}
	return strings.TrimSpace(proto + " " + r.Status)
}

// needsAuthRetry is the auth failure detection of
// CurlDownloader::isAuthenticatedRetryNeeded beyond plain 401/403: the
// Bitbucket login page served for a private zip and GitLab's 404 for an
// archive without credentials.
func (d *Downloader) needsAuthRetry(rawURL, origin string, resp *Response) string {
	if resp == nil {
		return ""
	}
	if origin == "bitbucket.org" && strings.HasSuffix(rawURL, ".zip") {
		if u, err := url.Parse(rawURL); err == nil && !isPublicBitbucketDownload(u) {
			loc := resp.Header.Get("Location")
			if (loc == "" || !strings.HasSuffix(loc, ".zip")) && regexp.MustCompile(`(?i)^text/html\b`).MatchString(resp.Header.Get("Content-Type")) {
				return "Bitbucket requires authentication and it was not provided"
			}
		}
	}
	if resp.StatusCode == 404 && containsStr(configStrings(d.cfg.Get("gitlab-domains")), origin) && strings.Contains(rawURL, "archive.zip") {
		return "GitLab requires authentication and it was not provided"
	}
	return ""
}

// maxAuthRetries guards the retry loop (Composer relies on the prompts).
const maxAuthRetries = 5

// checkAuthFailure is the auth handling of CurlDownloader for one response:
// 401/403 go through promptAuthIfNeeded, the Bitbucket login page and
// GitLab archive 404 fail with their own reason. It reports whether the
// request is to be repeated with the (possibly new) credentials, and the
// storeAuth value of that retry.
func (d *Downloader) checkAuthFailure(ctx context.Context, rawURL, dest string, resp *Response, err error, authRetries int) (bool, any, error) {
	if resp == nil || authRetries >= maxAuthRetries {
		return false, false, nil
	}
	retryAuth := !authRetryDisabled(ctx)
	origin := d.Origin(rawURL)
	var te *TransportError
	if retryAuth && err != nil && errors.As(err, &te) && (resp.StatusCode == 401 || resp.StatusCode == 403) {
		again, store, perr := d.promptAuth(rawURL, origin, resp.StatusCode, statusLine(resp), resp, authRetries)
		// CurlDownloader attaches the response headers to the exception
		var pte *TransportError
		if errors.As(perr, &pte) && pte.Header == nil {
			pte.Header = resp.Header
		}
		if perr != nil || again {
			return again, store, perr
		}
	}
	msg := d.needsAuthRetry(rawURL, origin, resp)
	if msg == "" {
		return false, false, nil
	}
	if retryAuth {
		again, store, perr := d.promptAuth(rawURL, origin, 401, "", nil, authRetries)
		if perr != nil || again {
			return again, store, perr
		}
	}
	if dest != "" {
		_ = os.Remove(dest)
	}
	return false, false, &TransportError{Msg: `The "` + Sanitize(rawURL) + `" file could not be downloaded (` + msg + `)` + jsonFailureDetails(resp.Header.Get("Content-Type"), resp.Body), StatusCode: resp.StatusCode}
}
