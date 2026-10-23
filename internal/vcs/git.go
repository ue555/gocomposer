package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/Git.php, src/Composer/Util/GitHub.php,
// src/Composer/Util/GitLab.php, src/Composer/Util/Bitbucket.php
// authorizeOAuth). Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/fsutil"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// Git is Composer\Util\Git.
type Git struct {
	IO      *cio.IO
	Config  *config.Config
	Process *Process
	HTTP    *httpdl.Downloader
	Auth    *AuthStore
}

var (
	gitVersionOnce sync.Once
	gitVersion     string
	reGitVersion   = regexp.MustCompile(`(?m)^git version (\d+(?:\.\d+)+)`)
	reSSHNoPort    = regexp.MustCompile(`^ssh://[^@]+@[^:]+:[^0-9]+`)
	reRemoteCreds  = regexp.MustCompile(`(?im)^(?:composer|origin)\s+https?://(.+):(.+)@([^/]+)`)
	reHeadBranch   = regexp.MustCompile(`(?m)^\s*HEAD branch:\s(.+)\s*$`)
	reSha40        = regexp.MustCompile(`^[a-f0-9]{40}$`)
)

// GitVersion is Git::getVersion ("" when git is unavailable).
func GitVersion(p *Process) string {
	gitVersionOnce.Do(func() {
		code, out, _ := p.Execute([]string{"git", "--version"}, "")
		if code == 0 {
			if m := reGitVersion.FindStringSubmatch(out); m != nil {
				gitVersion = m[1]
			}
		}
	})
	return gitVersion
}

// CleanEnv is Git::cleanEnv.
func CleanEnv(p *Process) {
	v := GitVersion(p)
	if v != "" && semver.VersionCompareOp(v, "2.3.0", ">=") {
		if os.Getenv("GIT_TERMINAL_PROMPT") != "0" {
			os.Setenv("GIT_TERMINAL_PROMPT", "0")
		}
	} else if os.Getenv("GIT_ASKPASS") != "echo" {
		os.Setenv("GIT_ASKPASS", "echo")
	}
	if os.Getenv("GIT_DIR") != "" {
		os.Unsetenv("GIT_DIR")
	}
	if os.Getenv("GIT_WORK_TREE") != "" {
		os.Unsetenv("GIT_WORK_TREE")
	}
	if os.Getenv("LANGUAGE") != "C" {
		os.Setenv("LANGUAGE", "C")
	}
	os.Unsetenv("DYLD_LIBRARY_PATH")
}

// SupportsNoCommitHeader is Git::supportsNoCommitHeaderFlag.
func SupportsNoCommitHeader(p *Process) bool {
	v := GitVersion(p)
	return v != "" && semver.VersionCompareOp(v, "2.33.0-rc0", ">=")
}

// NoShowSignatureFlags is Git::getNoShowSignatureFlags.
func NoShowSignatureFlags(p *Process) []string {
	v := GitVersion(p)
	if v != "" && semver.VersionCompareOp(v, "2.10.0-rc0", ">=") {
		return []string{"--no-show-signature"}
	}
	return nil
}

// RevListCommand is Git::buildRevListCommand.
func RevListCommand(p *Process, args ...string) []string {
	cmd := []string{"git", "rev-list"}
	if SupportsNoCommitHeader(p) {
		cmd = append(cmd, "--no-commit-header")
	}
	return append(cmd, args...)
}

var reCommitHeader = regexp.MustCompile(`(?m)^commit [a-f0-9]{40}\n?`)

// ParseRevListOutput is Git::parseRevListOutput.
func ParseRevListOutput(p *Process, out string) string {
	if SupportsNoCommitHeader(p) {
		return out
	}
	return reCommitHeader.ReplaceAllString(out, "")
}

// CheckForRepoOwnershipError is Git::checkForRepoOwnershipError without IO
// (it throws).
func CheckForRepoOwnershipError(output, path string) error {
	if strings.Contains(output, "fatal: detected dubious ownership") {
		return &RuntimeError{Msg: `The repository at "` + path + `" does not have the correct ownership and git refuses to use it:` + "\n\n" + output}
	}
	return nil
}

func (g *Git) prohibit(u string) error {
	if g.HTTP != nil {
		return g.HTTP.ProhibitURL(u)
	}
	return nil
}

// RunCommands is Git::runCommands: every command is run in order with
// %url% / %sanitizedUrl% replaced, retrying with credentials or other
// protocols like Composer. The concatenated output is returned.
func (g *Git) RunCommands(commands [][]string, u, cwd string, initialClone bool) (string, error) {
	if err := g.prohibit(u); err != nil {
		return "", err
	}
	origCwd := ""
	if initialClone {
		origCwd = cwd
	}
	var lastCommand []string
	var output string
	var procErr error
	run := func(u string) int {
		var outs []string
		status := 0
		for i, cmd := range commands {
			c := make([]string, len(cmd))
			for j, a := range cmd {
				switch a {
				case "%url%":
					c[j] = u
				case "%sanitizedUrl%":
					c[j] = StripCredentials(u)
				default:
					c[j] = a
				}
			}
			lastCommand = c
			dir := cwd
			if initialClone && i == 0 {
				dir = ""
			}
			code, out, err := g.Process.Execute(c, dir)
			outs = append(outs, out)
			if err != nil {
				procErr = err
				code = 1
			}
			status = code
			if code != 0 {
				break
			}
		}
		output = strings.Join(outs, "")
		return status
	}
	if reSSHNoPort.MatchString(u) {
		return "", &InvalidArgumentError{Msg: "The source URL " + Sanitize(u) + " is invalid, ssh URLs should have a port number after \":\".\nUse ssh://git@example.com:22/path or just git@example.com:path if you do not want to provide a password or custom port."}
	}
	if !initialClone {
		_, out, _ := g.Process.Execute([]string{"git", "remote", "-v"}, cwd)
		if m := reRemoteCreds.FindStringSubmatch(out); m != nil && !g.Auth.Has(m[3]) {
			user, _ := url.PathUnescape(m[1])
			pass, _ := url.PathUnescape(m[2])
			g.Auth.Set(m[3], user, pass)
		}
	}
	protocols := configList(g.Config, "github-protocols")
	ghRe := GitHubDomainsRegex(g.Config)
	if m := regexp.MustCompile(`^(?:https?|git)://` + ghRe + `/(.*)`).FindStringSubmatch(u); m != nil {
		var messages []string
		for _, proto := range protocols {
			var protoURL string
			if proto == "ssh" {
				protoURL = "git@" + m[1] + ":" + m[2]
			} else {
				protoURL = proto + "://" + m[1] + "/" + m[2]
			}
			if run(protoURL) == 0 {
				return output, nil
			}
			if procErr != nil {
				return output, procErr
			}
			messages = append(messages, "- "+protoURL+"\n"+regexp.MustCompile(`(?m)^`).ReplaceAllString(g.Process.ErrorOutput(), "  "))
			if initialClone && origCwd != "" {
				fsutil.RemoveAll(origCwd)
			}
		}
		if !g.Auth.Has(m[1]) && !g.IO.Interactive {
			return "", g.throwException("Failed to clone "+u+" via "+strings.Join(protocols, ", ")+" protocols, aborting.\n\n"+strings.Join(messages, "\n"), u)
		}
	}
	bypassSSH := regexp.MustCompile(`(?i)^git@`+ghRe+`:(.+?)\.git$`).MatchString(u) && !containsString(protocols, "ssh")
	var credentials []string
	if bypassSSH || run(u) != 0 {
		if procErr != nil {
			return output, procErr
		}
		errorMsg := g.Process.ErrorOutput()
		var m []string
		if m = regexp.MustCompile(`(?i)^git@` + ghRe + `:(.+?)\.git$`).FindStringSubmatch(u); m == nil {
			m = regexp.MustCompile(`(?i)^https?://` + ghRe + `/(.*?)(?:\.git)?$`).FindStringSubmatch(u)
		}
		if m != nil {
			if !g.Auth.Has(m[1]) {
				if !g.authorizeGitHub(m[1]) && g.interactive() {
					if _, err := g.HTTP.GitHubAuthorizeInteractively(m[1], "Cloning failed using an ssh key for authentication, enter your GitHub credentials to access private repos"); err != nil {
						return output, err
					}
				}
			}
			if g.Auth.Has(m[1]) {
				a := g.Auth.Get(m[1])
				authURL := "https://" + rawURLEncode(a.Username) + ":" + rawURLEncode(a.Password) + "@" + m[1] + "/" + m[2] + ".git"
				if run(authURL) == 0 {
					return output, nil
				}
				credentials = []string{rawURLEncode(a.Username), rawURLEncode(a.Password)}
				errorMsg = g.Process.ErrorOutput()
			}
		} else if bm := bitbucketMatch(u); bm != nil {
			domain := bm[2]
			repoWithGit := bm[3]
			if !strings.HasSuffix(repoWithGit, ".git") {
				repoWithGit += ".git"
			}
			var bitbucketUtil *httpdl.Bitbucket
			if g.HTTP != nil {
				bitbucketUtil = httpdl.NewBitbucket(g.HTTP)
			}
			if !g.Auth.Has(domain) {
				if !g.authorizeBitbucket(domain) && g.interactive() {
					if _, err := bitbucketUtil.AuthorizeOAuthInteractively(domain, "Enter your Bitbucket credentials to access private repos"); err != nil {
						return output, err
					}
					g.Auth.Set(domain, "x-token-auth", bitbucketUtil.Token())
				}
			}
			if g.Auth.Has(domain) {
				a := g.Auth.Get(domain)
				if strings.HasPrefix(a.Password, "ATAT") {
					a.Username = "x-bitbucket-api-token-auth"
				}
				authURL := "https://" + rawURLEncode(a.Username) + ":" + rawURLEncode(a.Password) + "@" + domain + "/" + repoWithGit
				if run(authURL) == 0 {
					return output, nil
				}
				if a.Username != "x-token-auth" && bitbucketUtil != nil {
					token, err := bitbucketUtil.RequestToken(domain, a.Username, a.Password)
					if err != nil {
						return output, err
					}
					if token != "" {
						g.Auth.Set(domain, "x-token-auth", token)
					}
				}
			}
			if g.Auth.Has(domain) {
				a := g.Auth.Get(domain)
				authURL := "https://" + rawURLEncode(a.Username) + ":" + rawURLEncode(a.Password) + "@" + domain + "/" + repoWithGit
				if run(authURL) == 0 {
					return output, nil
				}
				credentials = []string{rawURLEncode(a.Username), rawURLEncode(a.Password)}
			}
			sshURL := "git@ssh.bitbucket.org:" + repoWithGit
			if sshURL != u {
				g.IO.WriteError("    No bitbucket authentication configured. Falling back to ssh.")
				if run(sshURL) == 0 {
					return output, nil
				}
				errorMsg = g.Process.ErrorOutput()
			}
		} else if lm := gitlabMatch(g.Config, u); lm != nil {
			if lm[1] == "git" {
				lm[1] = "https"
			}
			if !g.Auth.Has(lm[2]) {
				if !g.authorizeGitLab(lm[2]) && g.interactive() {
					if _, err := g.HTTP.GitLabAuthorizeInteractively(lm[1], lm[2], "Cloning failed, enter your GitLab credentials to access private repos"); err != nil {
						return output, err
					}
				}
			}
			if g.Auth.Has(lm[2]) {
				a := g.Auth.Get(lm[2])
				var authURL string
				if a.Password == "private-token" || a.Password == "oauth2" || a.Password == "gitlab-ci-token" {
					authURL = lm[1] + "://" + rawURLEncode(a.Password) + ":" + rawURLEncode(a.Username) + "@" + lm[2] + "/" + lm[3]
				} else {
					authURL = lm[1] + "://" + rawURLEncode(a.Username) + ":" + rawURLEncode(a.Password) + "@" + lm[2] + "/" + lm[3]
				}
				if run(authURL) == 0 {
					return output, nil
				}
				credentials = []string{rawURLEncode(a.Username), rawURLEncode(a.Password)}
				errorMsg = g.Process.ErrorOutput()
			}
		} else if am := g.authenticationFailure(u); am != nil {
			authParts, hasAuthParts := "", false
			if i := strings.Index(am[2], "@"); i >= 0 {
				authParts, hasAuthParts = am[2][:i], true
				am[2] = am[2][i+1:]
			}
			var storeAuth any = false
			var user, pass any
			haveAuth := false
			if g.Auth.Has(am[2]) {
				a := g.Auth.Get(am[2])
				user, pass, haveAuth = a.Username, a.Password, true
			} else if g.interactive() {
				var defaultUser any
				if hasAuthParts && authParts != "" {
					defaultUser, _, _ = strings.Cut(authParts, ":")
				}
				g.IO.WriteError("    Authentication required (<info>" + am[2] + "</info>):")
				g.IO.WriteErrorV("<warning>"+strings.Trim(errorMsg, " \t\n\r\x00\x0B")+"</warning>", true, cio.Verbose)
				var err error
				if user, err = g.IO.Ask("      Username: ", defaultUser); err != nil {
					return output, err
				}
				if pass, err = g.IO.AskAndHideAnswer("      Password: "); err != nil {
					return output, err
				}
				haveAuth = true
				storeAuth = configStoreAuths(g.Config.Get("store-auths"))
			}
			if haveAuth {
				us, _ := cio.AnswerString(user)
				ps, _ := cio.AnswerString(pass)
				authURL := am[1] + rawURLEncode(us) + ":" + rawURLEncode(ps) + "@" + am[2] + am[3]
				if run(authURL) == 0 {
					g.Auth.SetPrompted(am[2], user, pass)
					if g.HTTP != nil {
						if err := g.HTTP.StoreAuth(am[2], storeAuth); err != nil {
							return output, err
						}
					}
					return output, nil
				}
				credentials = []string{rawURLEncode(us), rawURLEncode(ps)}
				errorMsg = g.Process.ErrorOutput()
			}
		}
		if initialClone && origCwd != "" {
			fsutil.RemoveAll(origCwd)
		}
		last := strings.Join(lastCommand, " ")
		if len(credentials) > 0 {
			last = maskCredentials(last, credentials)
			errorMsg = maskCredentials(errorMsg, credentials)
		}
		return output, g.throwException("Failed to execute "+last+"\n\n"+errorMsg, u)
	}
	return output, nil
}

func containsString(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// rawURLEncode is PHP rawurlencode.
func rawURLEncode(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}

var (
	reBitbucketHTTP = regexp.MustCompile(`(?i)^(https?)://(bitbucket\.org)/(.*?)(?:\.git)?$`)
	reBitbucketGit  = regexp.MustCompile(`(?i)^(git)@(?:ssh\.)?(bitbucket\.org):(.+?\.git)$`)
	reBitbucketSSH  = regexp.MustCompile(`(?i)^(ssh)://git@(?:ssh\.)?(bitbucket\.org)/(.*?)(?:\.git)?$`)
	reAuthFailURL   = regexp.MustCompile(`(?i)^(https?://)([^/]+)(.*)$`)
)

func bitbucketMatch(u string) []string {
	for _, re := range []*regexp.Regexp{reBitbucketHTTP, reBitbucketGit, reBitbucketSSH} {
		if m := re.FindStringSubmatch(u); m != nil {
			return m
		}
	}
	return nil
}

func gitlabMatch(cfg *config.Config, u string) []string {
	re := GitLabDomainsRegex(cfg)
	if m := regexp.MustCompile(`(?i)^(git)@` + re + `:(.+?\.git)$`).FindStringSubmatch(u); m != nil {
		return m
	}
	return regexp.MustCompile(`(?i)^(https?)://` + re + `/(.*)`).FindStringSubmatch(u)
}

func (g *Git) authenticationFailure(u string) []string {
	m := reAuthFailURL.FindStringSubmatch(u)
	if m == nil {
		return nil
	}
	errOut := g.Process.ErrorOutput()
	for _, f := range []string{"fatal: Authentication failed", "remote error: Invalid username or password.", "error: 401 Unauthorized", "fatal: unable to access", "fatal: could not read Username"} {
		if strings.Contains(errOut, f) {
			return m
		}
	}
	return nil
}

// authorizeGitHub is GitHub::authorizeOAuth.
func (g *Git) authorizeGitHub(origin string) bool {
	if !containsString(configList(g.Config, "github-domains"), origin) {
		return false
	}
	if code, out, _ := g.Process.Execute([]string{"git", "config", "github.accesstoken"}, ""); code == 0 {
		g.Auth.Set(origin, strings.TrimSpace(out), "x-oauth-basic")
		return true
	}
	return false
}

// authorizeGitLab is GitLab::authorizeOAuth.
func (g *Git) authorizeGitLab(origin string) bool {
	bc := regexp.MustCompile(`:\d+`).ReplaceAllString(origin, "")
	domains := configList(g.Config, "gitlab-domains")
	if !containsString(domains, origin) && !containsString(domains, bc) {
		return false
	}
	if code, out, _ := g.Process.Execute([]string{"git", "config", "gitlab.accesstoken"}, ""); code == 0 {
		g.Auth.Set(origin, strings.TrimSpace(out), "oauth2")
		return true
	}
	if c1, user, _ := g.Process.Execute([]string{"git", "config", "gitlab.deploytoken.user"}, ""); c1 == 0 {
		if c2, tok, _ := g.Process.Execute([]string{"git", "config", "gitlab.deploytoken.token"}, ""); c2 == 0 {
			g.Auth.Set(origin, strings.TrimSpace(user), strings.TrimSpace(tok))
			return true
		}
	}
	// token from the composer config
	if tokens := phpjson.ToMap(g.Config.Get("gitlab-token")); tokens != nil {
		var tok phpjson.Value
		found := false
		for _, k := range []string{origin, bc} {
			if v, ok := tokens.Get(k); ok && v != nil {
				tok, found = v, true
			}
		}
		if found {
			user, pass := "", "private-token"
			if m := phpjson.ToMap(tok); m != nil && phpjson.IsArray(tok) {
				user, pass = m.Str("username"), m.Str("token")
			} else {
				user, _ = phpjson.ScalarString(tok)
			}
			if user == "private-token" || user == "gitlab-ci-token" || user == "oauth2" {
				g.Auth.Set(origin, pass, user)
			} else {
				g.Auth.Set(origin, user, pass)
			}
			return true
		}
	}
	return false
}

// authorizeBitbucket is Bitbucket::authorizeOAuth.
func (g *Git) authorizeBitbucket(origin string) bool {
	if origin != "bitbucket.org" {
		return false
	}
	if code, out, _ := g.Process.Execute([]string{"git", "config", "bitbucket.accesstoken"}, ""); code == 0 {
		g.Auth.Set(origin, "x-token-auth", strings.TrimSpace(out))
		return true
	}
	return false
}

var nonSecretCredentials = []string{"private-token", "x-token-auth", "oauth2", "gitlab-ci-token", "x-oauth-basic"}

func maskCredentials(s string, creds []string) string {
	var pairs []string
	for _, c := range creds {
		var masked string
		switch {
		case containsString(nonSecretCredentials, c):
			masked = c
		case len(c) > 6:
			masked = c[:3] + "..." + c[len(c)-3:]
		case len(c) > 3:
			masked = c[:3] + "..."
		default:
			masked = "XXX"
		}
		if c != "" {
			pairs = append(pairs, c, masked)
		}
	}
	return strings.NewReplacer(pairs...).Replace(s)
}

func (g *Git) throwException(message, u string) error {
	if code, _, _ := g.Process.Execute([]string{"git", "--version"}, ""); code != 0 {
		return &RuntimeError{Msg: Sanitize("Failed to clone " + u + ", git was not found, check that it is installed and in your PATH env.\n\n" + g.Process.ErrorOutput())}
	}
	return &RuntimeError{Msg: Sanitize(message)}
}

func isBareGitDir(p *Process, dir string) bool {
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return false
	}
	code, out, _ := p.Execute([]string{"git", "rev-parse", "--git-dir"}, dir)
	return code == 0 && strings.TrimSpace(out) == "."
}

// SyncMirror is Git::syncMirror.
func (g *Git) SyncMirror(u, dir string) (bool, error) {
	if v := os.Getenv("COMPOSER_DISABLE_NETWORK"); v != "" && v != "0" && v != "prime" {
		g.IO.WriteError("<warning>Aborting git mirror sync of " + Sanitize(u) + " as network is disabled</warning>")
		return false, nil
	}
	if isBareGitDir(g.Process, dir) {
		_, err := g.RunCommands([][]string{
			{"git", "remote", "set-url", "origin", "--", "%url%"},
			{"git", "remote", "update", "--prune", "origin"},
			{"git", "gc", "--auto"},
		}, u, dir, false)
		_, err2 := g.RunCommands([][]string{{"git", "remote", "set-url", "origin", "--", "%sanitizedUrl%"}}, u, dir, false)
		if err == nil {
			err = err2
		}
		if err != nil {
			g.IO.WriteErrorV("<error>Sync mirror failed: "+err.Error()+"</error>", true, cio.Debug)
			return false, nil
		}
		return true, nil
	}
	if err := CheckForRepoOwnershipError(g.Process.ErrorOutput(), dir); err != nil {
		return false, err
	}
	fsutil.RemoveAll(dir)
	if _, err := g.RunCommands([][]string{{"git", "clone", "--mirror", "--", "%url%", dir}}, u, dir, true); err != nil {
		return false, err
	}
	if _, err := g.RunCommands([][]string{{"git", "remote", "set-url", "origin", "--", "%sanitizedUrl%"}}, u, dir, false); err != nil {
		return false, err
	}
	return true, nil
}

// FetchRefOrSyncMirror is Git::fetchRefOrSyncMirror.
func (g *Git) FetchRefOrSyncMirror(u, dir, ref, prettyVersion string) (bool, error) {
	in, err := g.checkRefIsInMirror(dir, ref)
	if err != nil {
		return false, err
	}
	if in {
		if reSha40.MatchString(ref) && prettyVersion != "" {
			branch := reDevAffix.ReplaceAllString(prettyVersion, "")
			var branches, tags *string
			if code, out, _ := g.Process.Execute([]string{"git", "branch"}, dir); code == 0 {
				branches = &out
			}
			if code, out, _ := g.Process.Execute([]string{"git", "tag"}, dir); code == 0 {
				tags = &out
			}
			if branches != nil && !regexp.MustCompile(`(?m)^[\s*]*v?`+regexp.QuoteMeta(branch)+`$`).MatchString(*branches) &&
				tags != nil && !regexp.MustCompile(`(?m)^[\s*]*`+regexp.QuoteMeta(branch)+`$`).MatchString(*tags) {
				if _, err := g.SyncMirror(u, dir); err != nil {
					return false, err
				}
			}
		}
		return true, nil
	}
	ok, err := g.SyncMirror(u, dir)
	if err != nil {
		return false, err
	}
	if ok {
		return g.checkRefIsInMirror(dir, ref)
	}
	return false, nil
}

var reDevAffix = regexp.MustCompile(`(?i)(?:^dev-|(?:\.x)?-dev$)`)

func (g *Git) checkRefIsInMirror(dir, ref string) (bool, error) {
	if isBareGitDir(g.Process, dir) {
		if code, _, _ := g.Process.Execute([]string{"git", "rev-parse", "--quiet", "--verify", ref + "^{commit}"}, dir); code == 0 {
			return true, nil
		}
	}
	if err := CheckForRepoOwnershipError(g.Process.ErrorOutput(), dir); err != nil {
		return false, err
	}
	return false, nil
}

// MirrorDefaultBranch is Git::getMirrorDefaultBranch.
func (g *Git) MirrorDefaultBranch(u, dir string, isLocal bool) string {
	if v := os.Getenv("COMPOSER_DISABLE_NETWORK"); v != "" && v != "0" {
		return ""
	}
	var out string
	if isLocal {
		_, out, _ = g.Process.Execute([]string{"git", "remote", "show", "origin"}, dir)
	} else {
		o, err := g.RunCommands([][]string{
			{"git", "remote", "set-url", "origin", "--", "%url%"},
			{"git", "remote", "show", "origin"},
		}, u, dir, false)
		_, err2 := g.RunCommands([][]string{{"git", "remote", "set-url", "origin", "--", "%sanitizedUrl%"}}, u, dir, false)
		if err == nil {
			err = err2
		}
		if err != nil {
			g.IO.WriteErrorV("<error>Failed to fetch root identifier from remote: "+err.Error()+"</error>", true, cio.Debug)
			return ""
		}
		out = o
	}
	for _, line := range SplitLines(out) {
		if m := reHeadBranch.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}

// interactive is IOInterface::isInteractive for the auth prompts (they need
// the HTTP layer to validate and store the credentials).
func (g *Git) interactive() bool { return g.IO != nil && g.IO.Interactive && g.HTTP != nil }

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
