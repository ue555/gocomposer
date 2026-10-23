package httpdl

// Interactive authentication: AuthHelper::storeAuth, the credential setters
// of IOInterface, GitHub::authorizeOAuthInteractively (and the rate limit /
// SSO header helpers), GitLab::authorizeOAuthInteractively / createToken and
// Bitbucket (requestToken, authorizeOAuthInteractively, storeInAuthConfig).
// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/AuthHelper.php, GitHub.php, GitLab.php, Bitbucket.php,
// src/Composer/IO/BaseIO.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.
//
// Secrets typed by the user (passwords, tokens, consumer secrets) only go
// to request headers and, like Composer, to the auth config file chosen by
// the store-auths setting / the token prompts; they are never written to
// the console or to errors.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/configsource"
	"gocomposer/internal/phpjson"
)

// ConfigSources are the files of Config::getConfigSource(),
// getAuthConfigSource() and getLocalAuthConfigSource() (Factory): the
// project composer.json (or COMPOSER_HOME/config.json), COMPOSER_HOME/
// auth.json and the auth.json next to composer.json when it exists.
type ConfigSources struct {
	Config    string
	Auth      string
	LocalAuth string
}

func (d *Downloader) configSource() *configsource.Source {
	p := d.Sources.Config
	if p == "" {
		p = d.cfg.GetString("home") + "/config.json"
	}
	return configsource.New(p, false)
}

func (d *Downloader) authConfigSource() *configsource.Source {
	p := d.Sources.Auth
	if p == "" {
		p = d.cfg.GetString("home") + "/auth.json"
	}
	return configsource.New(p, true)
}

func (d *Downloader) localAuthConfigSource() *configsource.Source {
	if d.Sources.LocalAuth == "" {
		return nil
	}
	return configsource.New(d.Sources.LocalAuth, true)
}

// interactive is IOInterface::isInteractive.
func (d *Downloader) interactive() bool { return d.IO != nil && d.IO.Interactive }

func (d *Downloader) writeError(msg string) {
	if d.IO != nil {
		d.IO.WriteError(msg)
	} else if d.Warn != nil {
		d.Warn(msg)
	}
}

// HasAuth is IOInterface::hasAuthentication.
func (d *Downloader) HasAuth(origin string) bool { return d.hasAuth(origin) }

// GetAuth is IOInterface::getAuthentication (ok is false when unset).
func (d *Downloader) GetAuth(origin string) (Auth, bool) { return d.getAuth(origin) }

// SetAuth is IOInterface::setAuthentication.
func (d *Downloader) SetAuth(origin, user, pass string) {
	d.mu.Lock()
	delete(d.nullCreds, origin)
	d.mu.Unlock()
	d.setAuth(origin, user, pass)
}

// SetPromptedAuth is setAuthentication with answers of the IO (nil when the
// question was left empty: PHP null).
func (d *Downloader) SetPromptedAuth(origin string, user, pass any) {
	u, uok := cio.AnswerString(user)
	p, pok := cio.AnswerString(pass)
	d.SetAuth(origin, u, p)
	if !uok || !pok {
		d.mu.Lock()
		if d.nullCreds == nil {
			d.nullCreds = map[string][2]bool{}
		}
		d.nullCreds[origin] = [2]bool{!uok, !pok}
		d.mu.Unlock()
	}
}

// authenticationValue is getAuthentication() as the array stored by
// AuthHelper::storeAuth.
func (d *Downloader) authenticationValue(origin string) *phpjson.Map {
	a, _ := d.getAuth(origin)
	d.mu.Lock()
	nulls := d.nullCreds[origin]
	d.mu.Unlock()
	m := phpjson.NewMap()
	if nulls[0] {
		m.Set("username", nil)
	} else {
		m.Set("username", a.Username)
	}
	if nulls[1] {
		m.Set("password", nil)
	} else {
		m.Set("password", a.Password)
	}
	return m
}

// StoreAuth is AuthHelper::storeAuth: storeAuth is true, false or "prompt".
func (d *Downloader) StoreAuth(origin string, storeAuth any) error {
	src := d.authConfigSource()
	var store *configsource.Source
	switch storeAuth {
	case true:
		store = src
	case "prompt":
		if d.IO == nil {
			break
		}
		answer, err := d.IO.AskAndValidate("Do you want to store credentials for "+origin+" in "+src.Name()+" ? [Yn] ", func(v any) (any, error) {
			s, _ := cio.AnswerString(v)
			s = strings.Trim(s, " \t\n\r\x00\x0B")
			if s != "" {
				s = strings.ToLower(s[:1])
			}
			if s == "y" || s == "n" {
				return s, nil
			}
			return nil, &cio.InvalidArgumentError{Msg: "Please answer (y)es or (n)o"}
		}, 0, "y")
		if err != nil {
			return err
		}
		if answer == "y" {
			store = src
		}
	}
	if store == nil {
		return nil
	}
	return store.AddConfigSetting("http-basic."+origin, d.authenticationValue(origin))
}

// --- GitHub ---

func headerLines(h http.Header) []string {
	var out []string
	for k, vs := range h {
		for _, v := range vs {
			out = append(out, k+": "+v)
		}
	}
	return out
}

var (
	reRateRemaining0 = regexp.MustCompile(`(?i)^x-ratelimit-remaining: *0$`)
	reSsoRequired    = regexp.MustCompile(`(?i)^x-github-sso: required`)
	reSsoURL         = regexp.MustCompile(`\burl=([^\s;]+)`)
)

// GitHubIsRateLimited is GitHub::isRateLimited.
func GitHubIsRateLimited(h http.Header) bool {
	for _, l := range headerLines(h) {
		if reRateRemaining0.MatchString(strings.TrimSpace(l)) {
			return true
		}
	}
	return false
}

// GitHubRateLimit is GitHub::getRateLimit as printed through sprintf('%d')
// ('?' gives 0) and the reset date ('?' when absent).
func GitHubRateLimit(h http.Header) (int, string) {
	limit, reset := 0, "?"
	if v := h.Get("X-RateLimit-Limit"); v != "" {
		limit = phpIntval(strings.TrimSpace(v))
	}
	if v := h.Get("X-RateLimit-Reset"); v != "" {
		reset = time.Unix(int64(phpIntval(strings.TrimSpace(v))), 0).UTC().Format("2006-01-02 15:04:05")
	}
	return limit, reset
}

func phpIntval(s string) int {
	end := 0
	for end < len(s) && (s[end] >= '0' && s[end] <= '9' || (end == 0 && (s[end] == '-' || s[end] == '+'))) {
		end++
	}
	n, _ := strconv.Atoi(s[:end])
	return n
}

func gitHubRequiresSso(h http.Header) bool {
	for _, l := range headerLines(h) {
		if reSsoRequired.MatchString(strings.TrimSpace(l)) {
			return true
		}
	}
	return false
}

func gitHubSsoURL(h http.Header) string {
	for _, l := range headerLines(h) {
		l = strings.TrimSpace(l)
		if !strings.Contains(strings.ToLower(l), "x-github-sso: required") {
			continue
		}
		if m := reSsoURL.FindStringSubmatch(l); m != nil {
			return m[1]
		}
	}
	return ""
}

// rawurlencode is PHP's rawurlencode.
func rawurlencode(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.' || c == '~' {
			sb.WriteByte(c)
		} else {
			fmt.Fprintf(&sb, "%%%02X", c)
		}
	}
	return sb.String()
}

// urlencode is PHP's urlencode (http_build_query's default encoding).
func urlencode(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.':
			sb.WriteByte(c)
		case c == ' ':
			sb.WriteByte('+')
		default:
			fmt.Fprintf(&sb, "%%%02X", c)
		}
	}
	return sb.String()
}

// GitHubAuthorizeInteractively is GitHub::authorizeOAuthInteractively.
func (d *Downloader) GitHubAuthorizeInteractively(origin, message string) (bool, error) {
	io := d.IO
	if message != "" {
		io.WriteError(message)
	}
	note := "Composer"
	if d.cfg.Get("github-expose-hostname") == true {
		if out, err := exec.Command("hostname").Output(); err == nil {
			note += " on " + strings.TrimSpace(string(out))
		}
	}
	note += " " + time.Now().UTC().Format("2006-01-02 1504")
	noteEnc := strings.ReplaceAll(rawurlencode(note), "%20", "+")
	local := d.localAuthConfigSource()
	where := d.authConfigSource().Name()
	if local != nil {
		where = local.Name() + " OR " + where
	}
	io.WriteError(strings.Join([]string{
		"You need to provide a GitHub access token.",
		`Tokens will be stored in plain text in "` + where + `" for future use by Composer.`,
		"Due to the security risk of tokens being exfiltrated, use tokens with short expiration times and only the minimum permissions necessary.",
		"",
		"Carefully consider the following options in order:",
		"",
	}, "\n"))
	io.WriteError(strings.Join([]string{
		"1. When you don't use 'vcs'  type 'repositories'  in composer.json and do not need to clone source or download dist files",
		"from private GitHub repositories over HTTPS, use a fine-grained token with read-only access to public information.",
		"Use the following URL to create such a token:",
		"https://" + origin + "/settings/personal-access-tokens/new?name=" + noteEnc,
		"",
	}, "\n"))
	io.WriteError(strings.Join([]string{
		`2. When all relevant _private_ GitHub repositories belong to a single user or organisation, use a fine-grained token with`,
		`repository "content" read-only permissions. You can start with the following URL, but you may need to change the resource owner`,
		"to the right user or organisation. Additionally, you can scope permissions down to apply only to selected repositories.",
		"https://" + origin + "/settings/personal-access-tokens/new?contents=read&name=" + noteEnc,
		"",
	}, "\n"))
	io.WriteError(strings.Join([]string{
		`3. A "classic" token grants broad permissions on your behalf to all repositories accessible by you.`,
		"This may include write permissions, even though not needed by Composer. Use it only when you need to access",
		"private repositories across multiple organisations at the same time and using directory-specific authentication sources",
		"is not an option. You can generate a classic token here:",
		"https://" + origin + "/settings/tokens/new?scopes=repo&description=" + noteEnc,
		"",
	}, "\n"))
	io.WriteError("For additional information, check https://getcomposer.org/doc/articles/authentication-for-private-packages.md#github-oauth")
	storeLocal := false
	if local != nil {
		v, err := io.AskConfirmation("A local auth config source was found, do you want to store the token there?", true)
		if err != nil {
			return false, err
		}
		storeLocal = v
	}
	ans, err := io.AskAndHideAnswer("Token (hidden): ")
	if err != nil {
		return false, err
	}
	token, _ := cio.AnswerString(ans)
	token = strings.Trim(token, " \t\n\r\x00\x0B")
	if token == "" {
		io.WriteError("<warning>No token given, aborting.</warning>")
		io.WriteError(`You can also add it manually later by using "composer config --global --auth github-oauth.github.com <token>"`)
		return false, nil
	}
	d.SetAuth(origin, token, "x-oauth-basic")
	apiURL := origin + "/api/v3/"
	if origin == "github.com" {
		apiURL = "api.github.com/"
	}
	if _, err := d.GetNoAuthRetry("https://" + apiURL); err != nil {
		var te *TransportError
		if errors.As(err, &te) && (te.StatusCode == 403 || te.StatusCode == 401) {
			io.WriteError("<error>Invalid token provided.</error>")
			io.WriteError(`You can also add it manually later by using "composer config --global --auth github-oauth.github.com <token>"`)
			return false, nil
		}
		return false, err
	}
	target := d.authConfigSource()
	if storeLocal && local != nil {
		target = local
	}
	if err := d.configSource().RemoveConfigSetting("github-oauth." + origin); err != nil {
		return false, err
	}
	if err := target.AddConfigSetting("github-oauth."+origin, token); err != nil {
		return false, err
	}
	io.WriteError("<info>Token stored successfully.</info>")
	return true, nil
}

// --- GitLab ---

// GitLabAuthorizeInteractively is GitLab::authorizeOAuthInteractively.
func (d *Downloader) GitLabAuthorizeInteractively(scheme, origin, message string) (bool, error) {
	io := d.IO
	if message != "" {
		io.WriteError(message)
	}
	local := d.localAuthConfigSource()
	patLink := scheme + "://" + origin + "/-/user_settings/personal_access_tokens"
	revokeLink := scheme + "://" + origin + "/-/user_settings/applications"
	where := d.authConfigSource().Name()
	if local != nil {
		where = local.Name() + " OR " + where
	}
	io.WriteError(`A token will be created and stored in "` + where + `", your password will never be stored`)
	io.WriteError("To revoke access to this token you can visit:")
	io.WriteError(revokeLink)
	io.WriteError("Alternatively you can setup an personal access token on:")
	io.WriteError(patLink)
	io.WriteError(`and store it under "gitlab-token" see https://getcomposer.org/doc/articles/authentication-for-private-packages.md#gitlab-token for more details.`)
	io.WriteError("https://getcomposer.org/doc/articles/authentication-for-private-packages.md#gitlab-token")
	io.WriteError("for more details.")
	storeLocal := false
	if local != nil {
		v, err := io.AskConfirmation("A local auth config source was found, do you want to store the token there?", true)
		if err != nil {
			return false, err
		}
		storeLocal = v
	}
	for attempt := 0; attempt < 5; attempt++ {
		token, err := d.gitLabCreateToken(scheme, origin)
		if err != nil {
			var te *TransportError
			if errors.As(err, &te) && (te.StatusCode == 403 || te.StatusCode == 401) {
				if te.StatusCode == 401 {
					invalidGrant := false
					if v, jerr := phpjson.Decode(te.body); jerr == nil {
						if m := phpjson.ToMap(v); m != nil && m.Str("error") == "invalid_grant" {
							invalidGrant = true
						}
					}
					if invalidGrant {
						io.WriteError("Bad credentials. If you have two factor authentication enabled you will have to manually create a personal access token")
					} else {
						io.WriteError("Bad credentials.")
					}
				} else {
					io.WriteError("Maximum number of login attempts exceeded. Please try again later.")
				}
				io.WriteError(`You can also manually create a personal access token enabling the "read_api" scope at:`)
				io.WriteError(patLink)
				io.WriteError(`Add it using "composer config --global --auth gitlab-token.` + origin + ` <token>"`)
				continue
			}
			return false, err
		}
		access := scalarString(token, "access_token")
		d.SetAuth(origin, access, "oauth2")
		target := d.authConfigSource()
		if storeLocal && local != nil {
			target = local
		}
		var value phpjson.Value = access
		if v, ok := token.Get("expires_in"); ok && v != nil {
			m := phpjson.NewMap()
			m.Set("expires-at", int64(phpIntval(scalarString(token, "created_at"))+phpIntval(scalarString(token, "expires_in"))))
			if rt, ok := token.Get("refresh_token"); ok {
				m.Set("refresh-token", rt)
			} else {
				m.Set("refresh-token", nil)
			}
			m.Set("token", access)
			value = m
		}
		if err := target.AddConfigSetting("gitlab-oauth."+origin, value); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, &AuthRuntimeError{Msg: "Invalid GitLab credentials 5 times in a row, aborting."}
}

// AuthRuntimeError is a \RuntimeException / \LogicException raised by the
// authentication helpers.
type AuthRuntimeError struct{ Msg string }

func (e *AuthRuntimeError) Error() string { return e.Msg }

func scalarString(m *phpjson.Map, key string) string {
	v, _ := m.Get(key)
	s, _ := phpjson.ScalarString(v)
	return s
}

// gitLabCreateToken is GitLab::createToken.
func (d *Downloader) gitLabCreateToken(scheme, origin string) (*phpjson.Map, error) {
	io := d.IO
	user, err := io.Ask("Username: ", nil)
	if err != nil {
		return nil, err
	}
	pass, err := io.AskAndHideAnswer("Password: ")
	if err != nil {
		return nil, err
	}
	var parts []string
	if s, ok := cio.AnswerString(user); ok {
		parts = append(parts, "username="+urlencode(s))
	}
	if s, ok := cio.AnswerString(pass); ok {
		parts = append(parts, "password="+urlencode(s))
	}
	parts = append(parts, "grant_type=password")
	resp, err := d.post(scheme+"://"+origin+"/oauth/token", strings.Join(parts, "&"), "Content-Type: application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	token, err := decodeJSONMap(resp)
	if err != nil {
		return nil, err
	}
	io.WriteError("Token successfully created")
	return token, nil
}

// post is HttpDownloader::get with an http POST context and
// 'retry-auth-failure' => false. A failed response keeps its body on the
// error (TransportException::getResponse).
func (d *Downloader) post(rawURL, body string, headers ...string) (*Response, error) {
	if err := d.ProhibitURL(rawURL); err != nil {
		return nil, err
	}
	resp, err, _ := d.doMethod(context.Background(), http.MethodPost, rawURL, "", headers, []byte(body))
	if err != nil {
		var te *TransportError
		if errors.As(err, &te) && resp != nil {
			te.body = resp.Body
		}
		return nil, err
	}
	return resp, nil
}

func decodeJSONMap(resp *Response) (*phpjson.Map, error) {
	v, err := phpjson.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	return m, nil
}

// --- Bitbucket ---

// BitbucketAccessTokenURL is Bitbucket::OAUTH2_ACCESS_TOKEN_URL.
const BitbucketAccessTokenURL = "https://bitbucket.org/site/oauth2/access_token"

// Bitbucket is Composer\Util\Bitbucket (one instance keeps the token it
// requested).
type Bitbucket struct {
	d     *Downloader
	token *phpjson.Map
}

// NewBitbucket returns a Bitbucket util bound to the downloader.
func NewBitbucket(d *Downloader) *Bitbucket { return &Bitbucket{d: d} }

// Token is Bitbucket::getToken.
func (b *Bitbucket) Token() string {
	if b.token == nil {
		return ""
	}
	return scalarString(b.token, "access_token")
}

func (b *Bitbucket) requestAccessToken() (bool, error) {
	resp, err := b.d.post(BitbucketAccessTokenURL, "grant_type=client_credentials", "Content-Type: application/x-www-form-urlencoded")
	if err != nil {
		var te *TransportError
		if errors.As(err, &te) {
			if te.StatusCode == 400 {
				b.d.writeError("<error>Invalid OAuth consumer provided.</error>")
				b.d.writeError("This can have three reasons:")
				b.d.writeError("1. You are authenticating with a bitbucket username/password combination")
				b.d.writeError("2. You are using an OAuth consumer, but didn't configure a (dummy) callback url")
				b.d.writeError("3. You are using an OAuth consumer, but didn't configure it as private consumer")
				return false, nil
			}
			if te.StatusCode == 403 || te.StatusCode == 401 {
				b.d.writeError("<error>Invalid OAuth consumer provided.</error>")
				b.d.writeError(`You can also add it manually later by using "composer config --global --auth bitbucket-oauth.bitbucket.org <consumer-key> <consumer-secret>"`)
				return false, nil
			}
		}
		return false, err
	}
	token, err := decodeJSONMap(resp)
	if err != nil {
		return false, err
	}
	if !token.Isset("expires_in") || !token.Isset("access_token") {
		enc, _ := phpjson.Encode(token, 0)
		return false, &AuthRuntimeError{Msg: "Expected a token configured with expires_in and access_token present, got " + enc}
	}
	b.token = token
	return true, nil
}

// AuthorizeOAuthInteractively is Bitbucket::authorizeOAuthInteractively.
func (b *Bitbucket) AuthorizeOAuthInteractively(origin, message string) (bool, error) {
	d := b.d
	io := d.IO
	if message != "" {
		io.WriteError(message)
	}
	local := d.localAuthConfigSource()
	where := d.authConfigSource().Name()
	if local != nil {
		where = local.Name() + " OR " + where
	}
	io.WriteError("Follow the instructions here:")
	io.WriteError("https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud/")
	io.WriteError(`to create a consumer. It will be stored in "` + where + `" for future use by Composer.`)
	io.WriteError(`Ensure you enter a "Callback URL" (http://example.com is fine) or it will not be possible to create an Access Token (this callback url will not be used by composer)`)
	storeLocal := false
	if local != nil {
		v, err := io.AskConfirmation("A local auth config source was found, do you want to store the token there?", true)
		if err != nil {
			return false, err
		}
		storeLocal = v
	}
	hint := `You can also add it manually later by using "composer config --global --auth bitbucket-oauth.bitbucket.org <consumer-key> <consumer-secret>"`
	ans, err := io.AskAndHideAnswer("Consumer Key (hidden): ")
	if err != nil {
		return false, err
	}
	key, _ := cio.AnswerString(ans)
	key = strings.Trim(key, " \t\n\r\x00\x0B")
	if key == "" || key == "0" {
		io.WriteError("<warning>No consumer key given, aborting.</warning>")
		io.WriteError(hint)
		return false, nil
	}
	ans, err = io.AskAndHideAnswer("Consumer Secret (hidden): ")
	if err != nil {
		return false, err
	}
	secret, _ := cio.AnswerString(ans)
	secret = strings.Trim(secret, " \t\n\r\x00\x0B")
	if secret == "" || secret == "0" {
		io.WriteError("<warning>No consumer secret given, aborting.</warning>")
		io.WriteError(hint)
		return false, nil
	}
	d.SetAuth(origin, key, secret)
	if ok, err := b.requestAccessToken(); err != nil || !ok {
		return false, err
	}
	target := d.authConfigSource()
	if storeLocal && local != nil {
		target = local
	}
	if err := b.storeInAuthConfig(target, origin, key, secret); err != nil {
		return false, err
	}
	if err := d.authConfigSource().RemoveConfigSetting("http-basic." + origin); err != nil {
		return false, err
	}
	io.WriteError("<info>Consumer stored successfully.</info>")
	return true, nil
}

// RequestToken is Bitbucket::requestToken.
func (b *Bitbucket) RequestToken(origin, key, secret string) (string, error) {
	if b.token != nil || b.tokenFromConfig(origin) {
		return scalarString(b.token, "access_token"), nil
	}
	b.d.SetAuth(origin, key, secret)
	ok, err := b.requestAccessToken()
	if err != nil || !ok {
		return "", err
	}
	src := b.d.localAuthConfigSource()
	if src == nil {
		src = b.d.authConfigSource()
	}
	if err := b.storeInAuthConfig(src, origin, key, secret); err != nil {
		return "", err
	}
	return scalarString(b.token, "access_token"), nil
}

// storeInAuthConfig is Bitbucket::storeInAuthConfig; like Composer the
// consumer always ends up in getAuthConfigSource() whatever source is given.
func (b *Bitbucket) storeInAuthConfig(_ *configsource.Source, origin, key, secret string) error {
	d := b.d
	if err := d.configSource().RemoveConfigSetting("bitbucket-oauth." + origin); err != nil {
		return err
	}
	if b.token == nil || !b.token.Isset("expires_in") {
		enc := "null"
		if b.token != nil {
			enc, _ = phpjson.Encode(b.token, 0)
		}
		return &AuthRuntimeError{Msg: "Expected a token configured with expires_in present, got " + enc}
	}
	consumer := phpjson.NewMap()
	consumer.Set("consumer-key", key)
	consumer.Set("consumer-secret", secret)
	consumer.Set("access-token", mustGet(b.token, "access_token"))
	consumer.Set("access-token-expiration", time.Now().Unix()+int64(phpIntval(scalarString(b.token, "expires_in"))))
	return d.authConfigSource().AddConfigSetting("bitbucket-oauth."+origin, consumer)
}

func mustGet(m *phpjson.Map, k string) phpjson.Value { v, _ := m.Get(k); return v }

// tokenFromConfig is Bitbucket::getTokenFromConfig.
func (b *Bitbucket) tokenFromConfig(origin string) bool {
	cfg := phpjson.ToMap(b.d.cfg.Get("bitbucket-oauth"))
	if cfg == nil {
		return false
	}
	c := cfg.Map(origin)
	if c == nil || !c.Isset("access-token") || !c.Isset("access-token-expiration") {
		return false
	}
	exp, _ := phpjson.ScalarString(mustGet(c, "access-token-expiration"))
	if n, err := strconv.ParseFloat(exp, 64); err != nil || float64(time.Now().Unix()) > n {
		return false
	}
	t := phpjson.NewMap()
	t.Set("access_token", mustGet(c, "access-token"))
	b.token = t
	return true
}

// GitLabIsOAuthExpired is GitLab::isOAuthExpired.
func (d *Downloader) GitLabIsOAuthExpired(origin string) bool {
	tokens := phpjson.ToMap(d.cfg.Get("gitlab-oauth"))
	if tokens == nil {
		return false
	}
	t := tokens.Map(origin)
	if t == nil || !t.Isset("expires-at") {
		return false
	}
	exp, _ := phpjson.ScalarString(mustGet(t, "expires-at"))
	n, err := strconv.ParseFloat(exp, 64)
	return err == nil && n < float64(time.Now().Unix())
}

// GitLabAuthorizeOAuthRefresh is GitLab::authorizeOAuthRefresh (with
// GitLab::refreshToken).
func (d *Downloader) GitLabAuthorizeOAuthRefresh(scheme, origin string) (bool, error) {
	var refresh phpjson.Value
	if tokens := phpjson.ToMap(d.cfg.Get("gitlab-oauth")); tokens != nil {
		if t := tokens.Map(origin); t != nil && t.Isset("refresh-token") {
			refresh = mustGet(t, "refresh-token")
		}
	}
	if refresh == nil {
		return false, &AuthRuntimeError{Msg: "No GitLab refresh token present for " + origin + "."}
	}
	rs, _ := phpjson.ScalarString(refresh)
	resp, err := d.post(scheme+"://"+origin+"/oauth/token", "refresh_token="+urlencode(rs)+"&grant_type=refresh_token", "Content-Type: application/x-www-form-urlencoded")
	var token *phpjson.Map
	if err == nil {
		token, err = decodeJSONMap(resp)
	}
	if err != nil {
		var te *TransportError
		if errors.As(err, &te) {
			d.writeError("Couldn't refresh access token: " + te.Msg)
			return false, nil
		}
		return false, err
	}
	if d.IO != nil {
		d.IO.WriteErrorV("GitLab token successfully refreshed", true, cio.VeryVerbose)
		d.IO.WriteErrorV("To revoke access to this token you can visit "+scheme+"://"+origin+"/-/user_settings/applications", true, cio.VeryVerbose)
	}
	access := scalarString(token, "access_token")
	d.SetAuth(origin, access, "oauth2")
	m := phpjson.NewMap()
	m.Set("expires-at", int64(phpIntval(scalarString(token, "created_at"))+phpIntval(scalarString(token, "expires_in"))))
	m.Set("refresh-token", mustGet(token, "refresh_token"))
	m.Set("token", access)
	if err := d.authConfigSource().AddConfigSetting("gitlab-oauth."+origin, m); err != nil {
		return false, err
	}
	return true, nil
}
