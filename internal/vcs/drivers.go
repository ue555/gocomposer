package vcs

// Driver selection of VcsRepository::getDriver and the static supports()
// checks of the drivers, ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/Repository/Vcs/*).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
)

// IsDriverType reports whether typ names a driver directly
// (isset($this->drivers[$this->type])).
func IsDriverType(typ string) bool {
	switch typ {
	case "github", "gitlab", "bitbucket", "git-bitbucket", "forgejo", "git", "hg", "perforce", "fossil", "svn":
		return true
	}
	return false
}

// NewDriver instantiates the driver class registered for typ.
func NewDriver(typ string, repoConfig *phpjson.Map, env *Env) Driver {
	switch typ {
	case "git":
		return NewGitDriver(repoConfig, env)
	case "github":
		return NewGitHubDriver(repoConfig, env)
	case "gitlab":
		return NewGitLabDriver(repoConfig, env)
	case "bitbucket", "git-bitbucket":
		return NewBitbucketDriver(repoConfig, env)
	case "svn":
		return NewSvnDriver(repoConfig, env)
	case "hg":
		return NewHgDriver(repoConfig, env)
	case "fossil":
		return NewFossilDriver(repoConfig, env)
	case "perforce":
		return NewPerforceDriver(repoConfig, env)
	}
	return newUnsupportedDriver(typ, repoConfig, env)
}

// Supports dispatches to the driver's static supports().
func Supports(typ string, env *Env, u string, deep bool) (bool, error) {
	switch typ {
	case "github":
		return GitHubSupports(env, u), nil
	case "gitlab":
		return GitLabSupports(env, u), nil
	case "bitbucket", "git-bitbucket":
		return BitbucketSupports(u), nil
	case "forgejo":
		return forgejoSupports(env, u), nil
	case "git":
		return GitSupports(env, u, deep)
	case "hg":
		return hgSupports(env, u, deep), nil
	case "perforce":
		return perforceSupports(env, u, deep), nil
	case "fossil":
		return fossilSupports(env, u), nil
	case "svn":
		return svnSupports(env, u, deep), nil
	}
	return false, nil
}

var (
	reHgURL       = regexp.MustCompile(`(?i)(^(?:https?|ssh)://(?:[^@]+@)?bitbucket.org|https://(?:.*?)\.kilnhg.com)`)
	reFossilURL   = regexp.MustCompile(`(?i)(^(?:https?|ssh)://(?:[^@]@)?(?:chiselapp\.com|fossil\.))`)
	reFossilPath  = regexp.MustCompile(`/fossil/|\.fossil`)
	reSvnURL      = regexp.MustCompile(`(?i)(^svn://|^svn\+ssh://|svn\.)`)
	rePerforceURL = regexp.MustCompile(`(?i)\b(perforce|p4)\b`)
)

// dirWritable is is_writable() for a directory.
func dirWritable(p string) bool {
	f, err := os.CreateTemp(p, ".gocomposer-w-*")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}

func isDir(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

func hgSupports(env *Env, u string, deep bool) bool {
	if reHgURL.MatchString(u) {
		return true
	}
	if IsLocalPath(u) {
		u = PlatformPath(u)
		if !isDir(u) {
			return false
		}
		p := NewProcess(env.IO, env.Process.Timeout)
		if code, _, _ := p.Execute([]string{"hg", "summary"}, u); code == 0 {
			return true
		}
	}
	if !deep {
		return false
	}
	p := NewProcess(env.IO, env.Process.Timeout)
	code, _, _ := p.Execute([]string{"hg", "identify", "--", u}, "")
	return code == 0
}

func fossilSupports(env *Env, u string) bool {
	if reFossilURL.MatchString(u) || reFossilPath.MatchString(u) {
		return true
	}
	if IsLocalPath(u) {
		u = PlatformPath(u)
		if !isDir(u) {
			return false
		}
		p := NewProcess(env.IO, env.Process.Timeout)
		if code, _, _ := p.Execute([]string{"fossil", "info"}, u); code == 0 {
			return true
		}
	}
	return false
}

func svnSupports(env *Env, u string, deep bool) bool {
	u = svnNormalizeURL(u)
	if reSvnURL.MatchString(u) {
		return true
	}
	if !deep && !IsLocalPath(u) {
		return false
	}
	p := NewProcess(env.IO, env.Process.Timeout)
	code, _, _ := p.Execute([]string{"svn", "info", "--non-interactive", "--", u}, "")
	if code == 0 {
		return true
	}
	errOut := strings.ToLower(p.ErrorOutput())
	return strings.Contains(errOut, "authorization failed:") || strings.Contains(errOut, "authentication failed")
}

// svnNormalizeURL is SvnDriver::normalizeUrl: absolute paths become
// file:// urls.
func svnNormalizeURL(u string) string {
	if strings.HasPrefix(u, "/") || (len(u) > 1 && u[1] == ':') || strings.HasPrefix(u, `\\`) {
		return "file://" + strings.ReplaceAll(u, `\`, "/")
	}
	return u
}

func perforceSupports(env *Env, u string, deep bool) bool {
	if deep || rePerforceURL.MatchString(u) {
		return PerforceCheckServerExists(u, NewProcess(env.IO, env.Process.Timeout))
	}
	return false
}

// unsupportedDriver stands for the forgejo driver: its repositories are
// detected like Composer does, but reading them is not implemented and fails
// explicitly.
type unsupportedDriver struct {
	base
	typ string
}

func newUnsupportedDriver(typ string, repoConfig *phpjson.Map, env *Env) *unsupportedDriver {
	return &unsupportedDriver{base: newBase(repoConfig, env), typ: typ}
}

func (d *unsupportedDriver) err() error {
	return &RuntimeError{Msg: "The " + d.typ + " VCS driver is not supported by gocomposer (repository " + Sanitize(d.url) + "); only git, github, gitlab, bitbucket, svn, hg, fossil and perforce repositories can be read"}
}

func (d *unsupportedDriver) Initialize() error                    { return d.err() }
func (d *unsupportedDriver) RootIdentifier() (string, error)      { return "", d.err() }
func (d *unsupportedDriver) Branches() ([]Ref, error)             { return nil, d.err() }
func (d *unsupportedDriver) Tags() ([]Ref, error)                 { return nil, d.err() }
func (d *unsupportedDriver) Source(string) *phpjson.Map           { return nil }
func (d *unsupportedDriver) Dist(string) (*phpjson.Map, error)    { return nil, d.err() }
func (d *unsupportedDriver) HasComposerFile(string) (bool, error) { return false, d.err() }
func (d *unsupportedDriver) ComposerInformation(string) (*phpjson.Map, error) {
	return nil, d.err()
}
