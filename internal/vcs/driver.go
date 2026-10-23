package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/VcsDriver.php, VcsDriverInterface.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"errors"
	"regexp"
	"time"

	"gocomposer/internal/cache"
	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
)

// Ref is one tag or branch: name => identifier.
type Ref struct {
	Name string
	ID   string
}

// Driver is VcsDriverInterface.
type Driver interface {
	Initialize() error
	RootIdentifier() (string, error)
	Branches() ([]Ref, error)
	Tags() ([]Ref, error)
	// ComposerInformation returns nil when the identifier has no (or an
	// empty) composer.json.
	ComposerInformation(identifier string) (*phpjson.Map, error)
	// HasComposerFile swallows transport errors like Composer.
	HasComposerFile(identifier string) (bool, error)
	Source(identifier string) *phpjson.Map
	Dist(identifier string) (*phpjson.Map, error)
	URL() string
	Cleanup()
}

// Env groups the collaborators every driver receives
// (VcsDriver::__construct arguments).
type Env struct {
	IO      *cio.IO
	Config  *config.Config
	HTTP    *httpdl.Downloader
	Process *Process
	Auth    *AuthStore
}

// Git returns a Git util bound to the environment.
func (e *Env) Git() *Git {
	return &Git{IO: e.IO, Config: e.Config, Process: e.Process, HTTP: e.HTTP, Auth: e.Auth}
}

// contentReader is implemented by concrete drivers (getFileContent /
// getChangeDate).
type contentReader interface {
	FileContent(file, identifier string) (string, bool, error)
	ChangeDate(identifier string) (time.Time, bool, error)
}

type infoEntry struct {
	data *phpjson.Map
}

// base is the abstract VcsDriver.
type base struct {
	*Env
	url        string
	originURL  string
	repoConfig *phpjson.Map
	infoCache  map[string]infoEntry
	cache      *cache.Cache
	reader     contentReader
}

func newBase(repoConfig *phpjson.Map, env *Env) base {
	u := repoConfig.Str("url")
	if IsLocalPath(u) {
		u = PlatformPath(u)
	}
	return base{Env: env, url: u, originURL: u, repoConfig: repoConfig, infoCache: map[string]infoEntry{}}
}

var reSha40i = regexp.MustCompile(`(?i)^[a-f0-9]{40}$`)

func (b *base) shouldCache(identifier string) bool {
	return b.cache != nil && reSha40i.MatchString(identifier)
}

// ComposerInformation is VcsDriver::getComposerInformation.
func (b *base) ComposerInformation(identifier string) (*phpjson.Map, error) {
	if e, ok := b.infoCache[identifier]; ok {
		return e.data, nil
	}
	if b.shouldCache(identifier) {
		if res, ok := b.cache.Read(identifier); ok && len(res) > 0 {
			v, err := jsonfile.Parse(res, "")
			if err != nil {
				return nil, err
			}
			m := phpjson.ToMap(v)
			b.infoCache[identifier] = infoEntry{m}
			return m, nil
		}
	}
	composer, err := b.baseComposerInformation(identifier)
	if err != nil {
		return nil, err
	}
	if b.shouldCache(identifier) {
		var v phpjson.Value
		if composer != nil {
			v = composer
		}
		if enc, err := phpjson.Encode(v, phpjson.UnescapedUnicode|phpjson.UnescapedSlashes); err == nil {
			b.cache.Write(identifier, []byte(enc))
		}
	}
	b.infoCache[identifier] = infoEntry{composer}
	return composer, nil
}

func (b *base) baseComposerInformation(identifier string) (*phpjson.Map, error) {
	content, ok, err := b.reader.FileContent("composer.json", identifier)
	if err != nil {
		return nil, err
	}
	if !ok || content == "" || content == "0" {
		return nil, nil
	}
	v, err := jsonfile.Parse([]byte(content), identifier+":composer.json")
	if err != nil {
		return nil, err
	}
	if !phpjson.IsArray(v) || phpjson.ArrayLen(v) == 0 {
		return nil, nil
	}
	composer := phpjson.ToMap(v)
	if t, ok := composer.Get("time"); !ok || !phpjson.Truthy(t) {
		d, ok, err := b.reader.ChangeDate(identifier)
		if err != nil {
			return nil, err
		}
		if ok {
			composer.Set("time", d.UTC().Format("2006-01-02T15:04:05-07:00"))
		}
	}
	return composer, nil
}

// HasComposerFile is VcsDriver::hasComposerFile: transport errors count as
// "no composer.json", everything else propagates.
func (b *base) HasComposerFile(identifier string) (bool, error) {
	m, err := b.ComposerInformation(identifier)
	if err != nil {
		var te *httpdl.TransportError
		if errors.As(err, &te) {
			return false, nil
		}
		return false, err
	}
	return m != nil, nil
}

func (b *base) Cleanup() {}

// URL is getUrl for drivers that return their (possibly rewritten) url.
func (b *base) URL() string { return b.url }
