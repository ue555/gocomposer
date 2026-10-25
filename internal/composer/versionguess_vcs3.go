package composer

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Package/Version/VersionGuesser.php: guessHgVersion,
// guessFossilVersion, guessSvnVersion and guessFeatureVersion for hg).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"regexp"
	"sort"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
	"gocomposer/internal/vcs"
)

// guessOtherVcsVersion runs the Mercurial, Fossil and Subversion guesses of
// VersionGuesser::guessVersion (after git found nothing), in that order.
func guessOtherVcsVersion(dir string, conf *phpjson.Map) (version, pretty string, ok bool) {
	p := vcs.NewProcess(nil, 0)
	if v, pr, ok := guessHgVersion(p, dir, conf); ok && v != "" {
		return v, pr, true
	}
	if v, pr, ok := guessFossilVersion(p, dir); ok && v != "" {
		return v, pr, true
	}
	if v, pr, ok := guessSvnVersion(p, dir, conf); ok && v != "" {
		return v, pr, true
	}
	return "", "", false
}

// guessHgVersion is VersionGuesser::guessHgVersion.
func guessHgVersion(p *vcs.Process, dir string, conf *phpjson.Map) (string, string, bool) {
	code, out, err := p.Execute([]string{"hg", "branch"}, dir)
	if err != nil || code != 0 {
		return "", "", false
	}
	branch := strings.TrimSpace(out)
	version := semver.NormalizeBranch(branch)
	if version == semver.DefaultBranchAlias {
		return version, "dev-" + branch, true
	}
	if !strings.HasPrefix(version, "dev-") {
		return version, version, true
	}
	refs, err := vcs.HgBranches(p, "")
	if err != nil {
		return "", "", false
	}
	var branches []string
	for _, r := range refs {
		branches = append(branches, r.Name)
	}
	v, pr := guessFeatureVersionWith(conf, version, branches, func(candidate, branch string) (string, bool) {
		code, out, err := p.Execute([]string{"hg", "log", "-r", "not ancestors('" + candidate + "') and ancestors('" + branch + "')", "--template", `"{node}\n"`}, dir)
		return out, err == nil && code == 0
	})
	return v, pr, true
}

// guessFeatureVersionWith is VersionGuesser::guessFeatureVersion with the
// scm command given as a runner (candidate, branch) -> output.
func guessFeatureVersionWith(conf *phpjson.Map, version string, branches []string, run func(candidate, branch string) (string, bool)) (string, string) {
	pretty := version
	hasAlias := false
	if extra := conf.Map("extra"); extra != nil {
		if ba := extra.Map("branch-alias"); ba != nil && ba.Has(version) {
			hasAlias = true
		}
	}
	enc, _ := phpjson.Encode(conf, 0)
	if hasAlias && !strings.Contains(enc, `"self.version"`) {
		return version, pretty
	}
	branch := strings.TrimPrefix(version, "dev-")
	if !isFeatureBranch(conf, branch) {
		return version, pretty
	}
	sort.SliceStable(branches, func(i, j int) bool {
		a, b := branches[i], branches[j]
		ar, br := strings.HasPrefix(a, "remotes/"), strings.HasPrefix(b, "remotes/")
		if ar != br {
			return !ar
		}
		return semver.NatCompare(strings.ToLower(b), strings.ToLower(a)) < 0
	})
	length := int(^uint(0) >> 1)
	lastIndex := -1
	for i, cand := range branches {
		candVersion := reRemotesPrefix.ReplaceAllString(cand, "")
		if cand == branch || isFeatureBranch(conf, candVersion) {
			continue
		}
		out, ok := run(cand, branch)
		if !ok {
			continue
		}
		if len(out) < length || (len(out) == length && lastIndex < i) {
			lastIndex = i
			length = len(out)
			version = semver.NormalizeBranch(candVersion)
			pretty = "dev-" + candVersion
			if length == 0 {
				break
			}
		}
	}
	return version, pretty
}

var reRemotesPrefix = regexp.MustCompile(`^remotes/\S+/`)

// guessFossilVersion is VersionGuesser::guessFossilVersion.
func guessFossilVersion(p *vcs.Process, dir string) (string, string, bool) {
	var version, pretty string
	if code, out, err := p.Execute([]string{"fossil", "branch", "list"}, dir); err == nil && code == 0 {
		branch := strings.TrimSpace(out)
		version = semver.NormalizeBranch(branch)
		pretty = "dev-" + branch
	}
	if code, out, err := p.Execute([]string{"fossil", "tag", "list"}, dir); err == nil && code == 0 {
		if n, err := semver.Normalize(strings.TrimSpace(out)); err == nil {
			version = n
			pretty = strings.TrimSpace(out)
		}
	}
	return version, pretty, version != ""
}

// guessSvnVersion is VersionGuesser::guessSvnVersion.
func guessSvnVersion(p *vcs.Process, dir string, conf *phpjson.Map) (string, string, bool) {
	vcs.SvnCleanEnv()
	code, out, err := p.Execute([]string{"svn", "info", "--xml"}, dir)
	if err != nil || code != 0 {
		return "", "", false
	}
	quote := func(key, def string) string {
		if v, ok := conf.Get(key); ok && v != nil {
			s, _ := phpjson.ScalarString(v)
			return regexp.QuoteMeta(s)
		}
		return def
	}
	trunk, branches, tags := quote("trunk-path", "trunk"), quote("branches-path", "branches"), quote("tags-path", "tags")
	re, rerr := regexp.Compile(`<url>.*/(` + trunk + `|(` + branches + `|` + tags + `)/(.*))</url>`)
	if rerr != nil {
		return "", "", false
	}
	m := re.FindStringSubmatchIndex(out)
	if m == nil {
		return "", "", false
	}
	group := func(i int) (string, bool) {
		if m[2*i] < 0 {
			return "", false
		}
		return out[m[2*i]:m[2*i+1]], true
	}
	g2, ok2 := group(2)
	g3, ok3 := group(3)
	if ok2 && ok3 && (g2 == strings.ReplaceAll(branches, `\`, "") || g2 == strings.ReplaceAll(tags, `\`, "")) {
		return semver.NormalizeBranch(g3), "dev-" + g3, true
	}
	g1, _ := group(1)
	pretty := strings.TrimSpace(g1)
	if pretty == "trunk" {
		return "dev-trunk", pretty, true
	}
	n, err := semver.Normalize(pretty)
	if err != nil {
		return "", "", false
	}
	return n, pretty, true
}
