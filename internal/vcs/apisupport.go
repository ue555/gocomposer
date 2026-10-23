package vcs

// supports() of the API based drivers, ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/Repository/Vcs/
// GitHubDriver.php, GitLabDriver.php, GitBitbucketDriver.php,
// ForgejoDriver.php, src/Composer/Util/ForgejoUrl.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"regexp"
	"strings"
)

var (
	reGitHubURL    = regexp.MustCompile(`^((?:https?|git)://([^/]+)/|git@([^:]+):/?)([^/]+)/([^/]+?)(?:\.git|/)?$`)
	reGitLabURL    = regexp.MustCompile(`^(?:(?P<scheme>https?)://(?P<domain>.+?)(?::(?P<port>[0-9]+))?/|git@(?P<domain2>[^:]+):)(?P<parts>.+)/(?P<repo>[^/]+?)(?:\.git|/)?$`)
	reBitbucketURL = regexp.MustCompile(`(?i)^https?://bitbucket\.org/([^/]+)/([^/]+?)(\.git|/?)?$`)
	reForgejoURL   = regexp.MustCompile(`^(?:(?:https?|git)://([^/]+)/|git@([^:]+):/?)([^/]+)/([^/]+?)(?:\.git|/)?$`)
	reWWWPrefix    = regexp.MustCompile(`(?i)^www\.`)
	rePortSuffix   = regexp.MustCompile(`:\d+`)
)

// GitHubSupports is GitHubDriver::supports.
func GitHubSupports(env *Env, u string) bool {
	m := reGitHubURL.FindStringSubmatch(u)
	if m == nil {
		return false
	}
	origin := m[2]
	if origin == "" {
		origin = m[3]
	}
	return containsString(configList(env.Config, "github-domains"), strings.ToLower(reWWWPrefix.ReplaceAllString(origin, "")))
}

type gitlabURL struct {
	scheme, domain, port string
	parts                []string
	repo                 string
	hasPort              bool
}

func parseGitLabURL(u string) *gitlabURL {
	m := reGitLabURL.FindStringSubmatch(u)
	if m == nil {
		return nil
	}
	g := &gitlabURL{}
	for i, n := range reGitLabURL.SubexpNames() {
		switch n {
		case "scheme":
			g.scheme = m[i]
		case "domain":
			if m[i] != "" {
				g.domain = m[i]
			}
		case "domain2":
			if g.domain == "" {
				g.domain = m[i]
			}
		case "port":
			g.port = m[i]
			g.hasPort = m[i] != ""
		case "parts":
			g.parts = strings.Split(m[i], "/")
		case "repo":
			g.repo = m[i]
		}
	}
	return g
}

// gitlabDetermineOrigin is GitLabDriver::determineOrigin; it consumes the
// leading url parts that belong to the origin.
func gitlabDetermineOrigin(domains []string, guessed string, parts *[]string, port string, hasPort bool) (string, bool) {
	guessed = strings.ToLower(guessed)
	if containsString(domains, guessed) || (hasPort && containsString(domains, guessed+":"+port)) {
		if hasPort {
			return guessed + ":" + port, true
		}
		return guessed, true
	}
	if hasPort {
		guessed += ":" + port
	}
	for len(*parts) > 0 {
		part := (*parts)[0]
		*parts = (*parts)[1:]
		guessed += "/" + part
		if containsString(domains, guessed) || (hasPort && containsString(domains, rePortSuffix.ReplaceAllString(guessed, ""))) {
			return guessed, true
		}
	}
	return "", false
}

// GitLabSupports is GitLabDriver::supports.
func GitLabSupports(env *Env, u string) bool {
	g := parseGitLabURL(u)
	if g == nil {
		return false
	}
	parts := g.parts
	_, ok := gitlabDetermineOrigin(configList(env.Config, "gitlab-domains"), g.domain, &parts, g.port, g.hasPort)
	return ok
}

// BitbucketSupports is GitBitbucketDriver::supports.
func BitbucketSupports(u string) bool { return reBitbucketURL.MatchString(u) }

// forgejoSupports is ForgejoDriver::supports.
func forgejoSupports(env *Env, u string) bool {
	m := reForgejoURL.FindStringSubmatch(u)
	if m == nil {
		return false
	}
	origin := m[1]
	if origin == "" {
		origin = m[2]
	}
	return containsString(configList(env.Config, "forgejo-domains"), strings.ToLower(origin))
}
