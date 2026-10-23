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
		// Perforce::checkServerExists
		p := NewProcess(env.IO, env.Process.Timeout)
		code, _, _ := p.Execute([]string{"p4", "-p", u, "info", "-s"}, "")
		return code == 0
	}
	return false
}

// unsupportedDriver stands for the hg, svn, fossil, perforce and forgejo
// drivers: their repositories are detected like Composer does and a missing
// client binary fails with Composer's own error, but reading them is not
// implemented and fails explicitly.
type unsupportedDriver struct {
	base
	typ string
	// hgMissing: HgDriver on a local path without hg installed reads no
	// tags/branches and then fails in getRootIdentifier.
	hgMissing bool
}

func newUnsupportedDriver(typ string, repoConfig *phpjson.Map, env *Env) *unsupportedDriver {
	return &unsupportedDriver{base: newBase(repoConfig, env), typ: typ}
}

func (d *unsupportedDriver) err() error {
	return &RuntimeError{Msg: "The " + d.typ + " VCS driver is not supported by gocomposer (repository " + Sanitize(d.url) + "); only git, github, gitlab and bitbucket repositories can be read"}
}

func (d *unsupportedDriver) binaryMissing(bin string) (bool, string) {
	code, _, _ := d.Process.Execute([]string{bin, "--version"}, "")
	return code == 127, d.Process.ErrorOutput()
}

func (d *unsupportedDriver) Initialize() error {
	switch d.typ {
	case "fossil":
		// FossilDriver::checkFossil
		code, _, _ := d.Process.Execute([]string{"fossil", "version"}, "")
		if code != 0 {
			return &RuntimeError{Msg: "fossil was not found, check that it is installed and in your PATH env.\n\n" + d.Process.ErrorOutput()}
		}
	case "perforce":
		if _, ok := d.repoConfig.Get("depot"); !ok {
			return &RuntimeError{Msg: `Undefined array key "depot"`}
		}
	case "svn":
		u := strings.TrimRight(svnNormalizeURL(d.url), "/")
		// Svn::execute() checks the url first (prohibitUrlByConfig); its
		// TransportException is a RuntimeException that SvnDriver::execute
		// replaces by the missing binary message below
		if d.HTTP != nil {
			_ = d.HTTP.ProhibitURL(u)
		}
		if missing, errOut := d.binaryMissing("svn"); missing {
			return &RuntimeError{Msg: "Failed to load " + Sanitize(u) + ", svn was not found, check that it is installed and in your PATH env.\n\n" + errOut}
		}
	case "hg":
		if IsLocalPath(d.url) {
			if missing, _ := d.binaryMissing("hg"); missing {
				d.hgMissing = true
				return nil
			}
			break
		}
		// HgDriver::initialize for a remote url: mirror into cache-vcs-dir
		cacheDir := d.Config.GetString("cache-vcs-dir")
		if !CacheUsable(cacheDir) {
			return &RuntimeError{Msg: "HgDriver requires a usable cache directory, and it looks like you set it to be disabled"}
		}
		_ = os.MkdirAll(cacheDir, 0o777)
		if !dirWritable(cacheDir) {
			return &RuntimeError{Msg: "Can not clone " + Sanitize(d.url) + ` to access package information. The "` + cacheDir + `" directory is not writable by the current user.`}
		}
		if d.HTTP != nil {
			if err := d.HTTP.ProhibitURL(d.url); err != nil {
				return err
			}
		}
		// Hg::runCommand fails the clone; Hg::throwException reports a
		// missing binary through Hg::getVersion()
		if missing, errOut := d.binaryMissing("hg"); missing {
			return &RuntimeError{Msg: Sanitize("Failed to clone " + d.url + ", hg was not found, check that it is installed and in your PATH env.\n\n" + errOut)}
		}
	}
	return d.err()
}

func (d *unsupportedDriver) RootIdentifier() (string, error) {
	if d.hgMissing {
		return "", &RuntimeError{Msg: "Undefined array key 0"}
	}
	return "", d.err()
}
func (d *unsupportedDriver) Branches() ([]Ref, error) {
	if d.hgMissing {
		return nil, nil
	}
	return nil, d.err()
}
func (d *unsupportedDriver) Tags() ([]Ref, error) {
	if d.hgMissing {
		return nil, nil
	}
	return nil, d.err()
}
func (d *unsupportedDriver) Source(string) *phpjson.Map           { return nil }
func (d *unsupportedDriver) Dist(string) (*phpjson.Map, error)    { return nil, d.err() }
func (d *unsupportedDriver) HasComposerFile(string) (bool, error) { return false, d.err() }
func (d *unsupportedDriver) ComposerInformation(string) (*phpjson.Map, error) {
	return nil, d.err()
}
