// Package pkg ports Composer's package model (Composer\Package\*) at
// f83588392df43d4fd027d84d26de459946ef4e80: BasePackage/Package/
// CompletePackage/RootPackage are one struct, AliasPackage is the same struct
// with AliasOf set (data getters delegate to the aliased package).
package pkg

import (
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// Link types (Link::TYPE_*).
const (
	TypeRequire    = "requires"
	TypeDevRequire = "devRequires"
	TypeProvide    = "provides"
	TypeConflict   = "conflicts"
	TypeReplace    = "replaces"
	// TypeDoesNotRequire is Link::TYPE_DOES_NOT_REQUIRE.
	TypeDoesNotRequire = "does not require"
)

// LinkKinds lists composer.json keys in BasePackage::$supportedLinkTypes order.
var LinkKinds = []struct {
	Key, Type, Description string
}{
	{"require", TypeRequire, "requires"},
	{"conflict", TypeConflict, "conflicts"},
	{"provide", TypeProvide, "provides"},
	{"replace", TypeReplace, "replaces"},
	{"require-dev", TypeDevRequire, "requires (for development)"},
}

// Stabilities maps stability names to BasePackage::STABILITIES priorities.
var Stabilities = map[string]int{"stable": 0, "RC": 5, "beta": 10, "alpha": 15, "dev": 20}

// StabilityNames lists stabilities from most to least stable.
var StabilityNames = []string{"stable", "RC", "beta", "alpha", "dev"}

const DefaultPrettyVersion = "1.0.0+no-version-set"

// Link is Composer\Package\Link.
type Link struct {
	Source           string
	Target           string
	Constraint       semver.Constraint
	Description      string
	PrettyConstraint string
	// Indexed marks links appended by AliasPackage (numeric keys in PHP).
	Indexed bool
}

// LinkDescription is the description Link::__construct stores for a link
// type (devRequires is shown as "requires (for development)").
func LinkDescription(linkType string) string {
	if linkType == TypeDevRequire {
		return "requires (for development)"
	}
	return linkType
}

func (l *Link) String() string {
	return l.Source + " " + l.Description + " " + l.Target + " (" + l.Constraint.String() + ")"
}

// PrettyString is Link::getPrettyString($sourcePackage).
func (l *Link) PrettyString(source *Package) string {
	return source.PrettyString() + " " + l.Description + " " + l.Target + " " + l.Constraint.PrettyString()
}

// Links is an ordered PHP array of links keyed by target.
type Links []*Link

// Get returns the link stored under key target (non-indexed entries).
func (ls Links) Get(target string) *Link {
	for _, l := range ls {
		if !l.Indexed && l.Target == target {
			return l
		}
	}
	return nil
}

func (ls Links) Has(target string) bool { return ls.Get(target) != nil }

// set replaces in place or appends, like $arr[$target] = $link.
func (ls Links) set(l *Link) Links {
	for i, e := range ls {
		if !e.Indexed && e.Target == l.Target {
			ls[i] = l
			return ls
		}
	}
	return append(ls, l)
}

// Repository is the subset of repository behaviour packages need.
type Repository interface {
	RepoName() string
}

// Package is Composer's package object.
type Package struct {
	Name          string
	PrettyName    string
	version       string
	prettyVersion string
	stability     string
	dev           bool

	typ                string
	targetDir          string
	extra              *phpjson.Map
	binaries           []string
	installationSource string
	sourceType         string
	sourceURL          string
	sourceReference    string
	sourceMirrors      []phpjson.Value
	distType           string
	distURL            string
	distReference      string
	distSha1           string
	distSha1Present    bool
	distMirrors        []phpjson.Value
	releaseDate        string
	publishedDate      string
	autoload           *phpjson.Map
	devAutoload        *phpjson.Map
	includePaths       []string
	phpExt             phpjson.Value
	notificationURL    string
	isDefaultBranch    bool
	transportOptions   *phpjson.Map
	suggests           *phpjson.Map

	requires    Links
	conflicts   Links
	provides    Links
	replaces    Links
	devRequires Links

	// CompletePackage
	scripts         *phpjson.Map
	license         []string
	authors         []phpjson.Value
	description     string
	homepage        string
	keywords        []string
	repositories    []phpjson.Value
	support         *phpjson.Map
	funding         []phpjson.Value
	abandoned       phpjson.Value
	archiveName     string
	archiveExcludes []string

	// RootPackage
	IsRoot           bool
	MinimumStability string
	PreferStable     bool
	StabilityFlags   map[string]int
	RootConfig       *phpjson.Map
	References       map[string]string
	Aliases          []RootAlias

	// AliasPackage
	AliasOf                *Package
	RootPackageAlias       bool
	hasSelfVersionRequires bool

	ID   int
	Repo Repository
}

// RootAlias is one entry of RootPackage::getAliases().
type RootAlias struct {
	Package         string
	Version         string
	Alias           string
	AliasNormalized string
}

// NewPackage creates a CompletePackage (BasePackage::__construct lowercases the name).
func NewPackage(name, version, prettyVersion string) *Package {
	p := &Package{Name: strings.ToLower(name), PrettyName: name, ID: -1, typ: "library"}
	p.setVersion(version, prettyVersion)
	return p
}

func (p *Package) setVersion(version, prettyVersion string) {
	p.version = version
	p.prettyVersion = prettyVersion
	p.stability = semver.ParseStability(version)
	p.dev = p.stability == "dev"
}

// ReplaceVersion is RootPackage::replaceVersion.
func (p *Package) ReplaceVersion(version, prettyVersion string) { p.setVersion(version, prettyVersion) }

// NewAlias creates an AliasPackage (or CompleteAliasPackage) of base.
func NewAlias(base *Package, version, prettyVersion string) *Package {
	a := &Package{Name: base.Name, PrettyName: base.PrettyName, ID: -1, AliasOf: base}
	a.setVersion(version, prettyVersion)
	a.requires = a.replaceSelfVersion(base.Requires(), TypeRequire)
	a.devRequires = a.replaceSelfVersion(base.DevRequires(), TypeDevRequire)
	a.conflicts = a.replaceSelfVersion(base.Conflicts(), TypeConflict)
	a.provides = a.replaceSelfVersion(base.Provides(), TypeProvide)
	a.replaces = a.replaceSelfVersion(base.Replaces(), TypeReplace)
	a.IsRoot = base.IsRoot
	return a
}

func (a *Package) replaceSelfVersion(links Links, linkType string) Links {
	prettyVersion := a.prettyVersion
	if prettyVersion == semver.DefaultBranchAlias {
		prettyVersion = a.AliasOf.PrettyVersion()
	}
	out := append(Links(nil), links...)
	if linkType == TypeConflict || linkType == TypeProvide || linkType == TypeReplace {
		for _, l := range links {
			if l.PrettyConstraint == "self.version" {
				c := semver.NewConstraint("=", a.version)
				c.SetPrettyString(prettyVersion)
				out = append(out, &Link{Source: l.Source, Target: l.Target, Constraint: c, Description: LinkDescription(linkType), PrettyConstraint: prettyVersion, Indexed: true})
			}
		}
		return out
	}
	for i, l := range out {
		if l.PrettyConstraint == "self.version" {
			if linkType == TypeRequire {
				a.hasSelfVersionRequires = true
			}
			c := semver.NewConstraint("=", a.version)
			c.SetPrettyString(prettyVersion)
			out[i] = &Link{Source: l.Source, Target: l.Target, Constraint: c, Description: LinkDescription(linkType), PrettyConstraint: prettyVersion}
		}
	}
	return out
}

// Base returns the non-alias package.
func (p *Package) Base() *Package {
	for p.AliasOf != nil {
		p = p.AliasOf
	}
	return p
}

func (p *Package) IsAlias() bool                { return p.AliasOf != nil }
func (p *Package) HasSelfVersionRequires() bool { return p.hasSelfVersionRequires }
func (p *Package) Version() string              { return p.version }
func (p *Package) PrettyVersion() string        { return p.prettyVersion }
func (p *Package) Stability() string            { return p.stability }
func (p *Package) IsDev() bool                  { return p.dev }
func (p *Package) Requires() Links              { return p.requires }
func (p *Package) DevRequires() Links           { return p.devRequires }
func (p *Package) Conflicts() Links             { return p.conflicts }
func (p *Package) Provides() Links              { return p.provides }
func (p *Package) Replaces() Links              { return p.replaces }

func (p *Package) SetRequires(l Links)    { p.requires = l }
func (p *Package) SetDevRequires(l Links) { p.devRequires = l }
func (p *Package) SetConflicts(l Links)   { p.conflicts = l }
func (p *Package) SetProvides(l Links)    { p.provides = l }
func (p *Package) SetReplaces(l Links)    { p.replaces = l }

// LinksOf returns the links of one type.
func (p *Package) LinksOf(linkType string) Links {
	switch linkType {
	case TypeRequire:
		return p.requires
	case TypeDevRequire:
		return p.devRequires
	case TypeConflict:
		return p.conflicts
	case TypeProvide:
		return p.provides
	case TypeReplace:
		return p.replaces
	}
	return nil
}

func (p *Package) SetLinksOf(linkType string, l Links) {
	switch linkType {
	case TypeRequire:
		p.requires = l
	case TypeDevRequire:
		p.devRequires = l
	case TypeConflict:
		p.conflicts = l
	case TypeProvide:
		p.provides = l
	case TypeReplace:
		p.replaces = l
	}
}

// Data getters delegate to the aliased package like AliasPackage does.

func (p *Package) Type() string               { return p.Base().typ }
func (p *Package) SetType(t string)           { p.Base().typ = t }
func (p *Package) TargetDir() string          { return p.Base().targetDir }
func (p *Package) Extra() *phpjson.Map        { return p.Base().extra }
func (p *Package) SetExtra(m *phpjson.Map)    { p.Base().extra = m }
func (p *Package) Binaries() []string         { return p.Base().binaries }
func (p *Package) InstallationSource() string { return p.Base().installationSource }
func (p *Package) SetInstallationSource(s string) {
	p.Base().installationSource = s
}
func (p *Package) SourceType() string             { return p.Base().sourceType }
func (p *Package) SourceURL() string              { return p.Base().sourceURL }
func (p *Package) SourceReference() string        { return p.Base().sourceReference }
func (p *Package) SourceMirrors() []phpjson.Value { return p.Base().sourceMirrors }
func (p *Package) DistType() string               { return p.Base().distType }
func (p *Package) DistURL() string                { return p.Base().distURL }
func (p *Package) DistReference() string          { return p.Base().distReference }
func (p *Package) DistSha1Checksum() string       { return p.Base().distSha1 }
func (p *Package) DistMirrors() []phpjson.Value   { return p.Base().distMirrors }

// PublishedDate is the repository-set publication time (raw string).
func (p *Package) PublishedDate() string     { return p.Base().publishedDate }
func (p *Package) ReleaseDate() string       { return p.Base().releaseDate }
func (p *Package) Autoload() *phpjson.Map    { return p.Base().autoload }
func (p *Package) DevAutoload() *phpjson.Map { return p.Base().devAutoload }
func (p *Package) IncludePaths() []string    { return p.Base().includePaths }
func (p *Package) NotificationURL() string   { return p.Base().notificationURL }
func (p *Package) IsDefaultBranch() bool     { return p.Base().isDefaultBranch }
func (p *Package) TransportOptions() *phpjson.Map {
	return p.Base().transportOptions
}
func (p *Package) Suggests() *phpjson.Map   { return p.Base().suggests }
func (p *Package) Scripts() *phpjson.Map    { return p.Base().scripts }
func (p *Package) License() []string        { return p.Base().license }
func (p *Package) Authors() []phpjson.Value { return p.Base().authors }
func (p *Package) Description() string      { return p.Base().description }
func (p *Package) Homepage() string         { return p.Base().homepage }
func (p *Package) Keywords() []string       { return p.Base().keywords }
func (p *Package) Support() *phpjson.Map    { return p.Base().support }
func (p *Package) Funding() []phpjson.Value { return p.Base().funding }

func (p *Package) SetAutoload(m *phpjson.Map)    { p.Base().autoload = m }
func (p *Package) SetDevAutoload(m *phpjson.Map) { p.Base().devAutoload = m }
func (p *Package) SetDist(typ, url, ref, sha string) {
	b := p.Base()
	b.distType, b.distURL, b.distReference, b.distSha1 = typ, url, ref, sha
}
func (p *Package) SetSource(typ, url, ref string) {
	b := p.Base()
	b.sourceType, b.sourceURL, b.sourceReference = typ, url, ref
}
func (p *Package) SetDistURL(u string)                { p.Base().distURL = u }
func (p *Package) SetSourceURL(u string)              { p.Base().sourceURL = u }
func (p *Package) SetDistReference(r string)          { p.Base().distReference = r }
func (p *Package) SetSourceReference(r string)        { p.Base().sourceReference = r }
func (p *Package) SetTransportOptions(m *phpjson.Map) { p.Base().transportOptions = m }
func (p *Package) SetDistMirrors(m []phpjson.Value)   { p.Base().distMirrors = m }
func (p *Package) SetSourceMirrors(m []phpjson.Value) { p.Base().sourceMirrors = m }
func (p *Package) SetNotificationURL(u string)        { p.Base().notificationURL = u }
func (p *Package) SetReleaseDate(s string)            { p.Base().releaseDate = s }
func (p *Package) SetDescription(s string)            { p.Base().description = s }
func (p *Package) SetBinaries(b []string)             { p.Base().binaries = b }
func (p *Package) SetScripts(m *phpjson.Map)          { p.Base().scripts = m }

// IsAbandoned and ReplacementPackage port CompletePackage.
func (p *Package) IsAbandoned() bool {
	switch v := p.Base().abandoned.(type) {
	case bool:
		return v
	case string:
		return true
	}
	return false
}

func (p *Package) ReplacementPackage() string {
	if s, ok := p.Base().abandoned.(string); ok {
		return s
	}
	return ""
}

// SetSourceDistReferences is Package::setSourceDistReferences.
func (p *Package) SetSourceDistReferences(reference string) {
	b := p.Base()
	b.sourceReference = reference
	// only bitbucket, github and gitlab have auto generated dist URLs that easily allow replacing the reference in the dist URL
	if b.distURL != "" && reGitHostDist.MatchString(b.distURL) {
		b.distReference = reference
		b.distURL = ReplaceSha40(b.distURL, reference)
	} else if b.distReference != "" {
		b.distReference = reference
	}
}

// Names is BasePackage::getNames.
func (p *Package) Names(withProvides bool) []string {
	seen := map[string]bool{p.Name: true}
	names := []string{p.Name}
	if withProvides {
		for _, l := range p.Provides() {
			if !seen[l.Target] {
				seen[l.Target] = true
				names = append(names, l.Target)
			}
		}
	}
	for _, l := range p.Replaces() {
		if !seen[l.Target] {
			seen[l.Target] = true
			names = append(names, l.Target)
		}
	}
	return names
}

func (p *Package) UniqueName() string   { return p.Name + "-" + p.version }
func (p *Package) PrettyString() string { return p.PrettyName + " " + p.prettyVersion }

func (p *Package) String() string {
	if p.AliasOf != nil {
		root := ""
		if p.RootPackageAlias {
			root = "root "
		}
		return p.UniqueName() + " (" + root + "alias of " + p.AliasOf.Version() + ")"
	}
	return p.UniqueName()
}

// FullPrettyVersion is getFullPrettyVersion($truncate, DISPLAY_SOURCE_REF_IF_DEV).
func (p *Package) FullPrettyVersion(truncate bool) string {
	st := p.SourceType()
	if !p.IsDev() || (st != "hg" && st != "git" && (st != "" || p.DistReference() == "")) {
		return p.prettyVersion
	}
	ref := p.SourceReference()
	if ref == "" {
		ref = p.DistReference()
	}
	if ref == "" {
		return p.prettyVersion
	}
	if truncate && len(ref) == 40 && st != "svn" {
		return p.prettyVersion + " " + ref[:7]
	}
	return p.prettyVersion + " " + ref
}

// StabilityPriority is getStabilityPriority.
func (p *Package) StabilityPriority() int { return Stabilities[p.stability] }

// Clone mirrors PHP clone: id reset and repository cleared.
func (p *Package) Clone() *Package {
	c := *p
	c.ID = -1
	c.Repo = nil
	return &c
}

var reGitHostDist = regexp.MustCompile(`(?i)^https?://(?:(?:www\.)?bitbucket\.org|(api\.)?github\.com|(?:www\.)?gitlab\.com)/`)

func isHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// ReplaceSha40 is Preg::replace('{(?<=/|sha=)[a-f0-9]{40}(?=/|$)}i', $ref, $url).
func ReplaceSha40(url, ref string) string {
	var sb strings.Builder
	i := 0
	for i < len(url) {
		if i+40 <= len(url) && (strings.HasSuffix(url[:i], "/") || strings.HasSuffix(strings.ToLower(url[:i]), "sha=")) {
			ok := true
			for j := i; j < i+40; j++ {
				if !isHex(url[j]) {
					ok = false
					break
				}
			}
			if ok && (i+40 == len(url) || url[i+40] == '/') {
				sb.WriteString(ref)
				i += 40
				continue
			}
		}
		sb.WriteByte(url[i])
		i++
	}
	return sb.String()
}
