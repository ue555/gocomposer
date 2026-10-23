package composer

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"

	"gocomposer/internal/config"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// GuessVersion ports VersionGuesser::guessVersion for git working copies.
// Mercurial, Fossil and Subversion detection are not implemented and yield
// no version (as when those tools are absent).
func GuessVersion(dir string, conf *phpjson.Map, cfg *config.Config) (string, string, bool) {
	version, pretty, ok := guessGitVersion(dir, conf)
	if !ok || version == "" {
		return "", "", false
	}
	if strings.HasSuffix(version, "-dev") && regexp.MustCompile(`\.9{7}`).MatchString(version) {
		pretty = regexp.MustCompile(`(\.9{7})+`).ReplaceAllString(version, ".x")
	}
	return pretty, version, true
}

func runGit(dir string, args ...string) (string, bool) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	env := os.Environ()
	var cleaned []string
	for _, e := range env {
		if strings.HasPrefix(e, "GIT_DIR=") || strings.HasPrefix(e, "GIT_WORK_TREE=") || strings.HasPrefix(e, "GIT_INDEX_FILE=") {
			continue
		}
		cleaned = append(cleaned, e)
	}
	cmd.Env = append(cleaned, "GIT_ASKPASS=echo", "LANGUAGE=C")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return out.String(), false
	}
	return out.String(), true
}

func gitCommit(dir string) string {
	out, ok := runGit(dir, "rev-list", "--format=%H", "-n1", "HEAD")
	if !ok {
		return ""
	}
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "commit ") {
			return l
		}
	}
	return ""
}

var (
	reCurrentBranch = regexp.MustCompile(`^(?:\* ) *(\(no branch\)|\(detached from \S+\)|\(HEAD detached at \S+\)|\S+) *([a-f0-9]+) .*$`)
	reRemoteHead    = regexp.MustCompile(`^ *.+/HEAD `)
	reBranchName    = regexp.MustCompile(`^(?:\* )? *((?:remotes/(?:origin|upstream)/)?[^\s/]+) *([a-f0-9]+) .*$`)
)

func isFeatureBranch(conf *phpjson.Map, branch string) bool {
	nonFeature := ""
	if v, ok := conf.Get("non-feature-branches"); ok && phpjson.Truthy(v) {
		var parts []string
		phpjson.Each(v, func(_ string, e phpjson.Value) { s, _ := e.(string); parts = append(parts, s) })
		nonFeature = strings.Join(parts, "|")
	}
	re, err := regexp.Compile(`^(` + nonFeature + `|master|main|latest|next|current|support|tip|trunk|default|develop|\d+\..+)$`)
	if err != nil {
		return true
	}
	return !re.MatchString(branch)
}

func guessGitVersion(dir string, conf *phpjson.Map) (string, string, bool) {
	var version, pretty string
	isDetached := false
	out, ok := runGit(dir, "branch", "-a", "--no-color", "--no-abbrev", "-v")
	if ok {
		var branches []string
		isFeature := false
		for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			line = strings.TrimRight(line, "\r")
			if line == "" {
				continue
			}
			if m := reCurrentBranch.FindStringSubmatch(line); m != nil {
				if m[1] == "(no branch)" || strings.HasPrefix(m[1], "(detached ") || strings.HasPrefix(m[1], "(HEAD detached at") {
					version = "dev-" + m[2]
					pretty = version
					isFeature = true
					isDetached = true
				} else {
					version = semver.NormalizeBranch(m[1])
					pretty = "dev-" + m[1]
					isFeature = isFeatureBranch(conf, m[1])
				}
			}
			if !reRemoteHead.MatchString(line) {
				if m := reBranchName.FindStringSubmatch(line); m != nil {
					branches = append(branches, m[1])
				}
			}
		}
		if isFeature {
			version, pretty = guessFeatureVersion(dir, conf, version, branches)
		}
	}
	if version == "" || isDetached {
		if tag, ok := runGit(dir, "describe", "--exact-match", "--tags"); ok {
			if n, err := semver.Normalize(strings.TrimSpace(tag)); err == nil {
				version = n
				pretty = strings.TrimSpace(tag)
			}
		}
	}
	return version, pretty, version != ""
}

func guessFeatureVersion(dir string, conf *phpjson.Map, version string, branches []string) (string, string) {
	pretty := version
	enc, _ := json.Marshal(conf)
	hasAlias := false
	if extra := conf.Map("extra"); extra != nil {
		if ba := extra.Map("branch-alias"); ba != nil && ba.Has(version) {
			hasAlias = true
		}
	}
	encPHP, _ := phpjson.Encode(conf, 0)
	_ = enc
	if hasAlias && !strings.Contains(encPHP, `"self.version"`) {
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
	type res struct {
		index int
		out   string
		ok    bool
		cand  string
	}
	var wg sync.WaitGroup
	results := make([]res, len(branches))
	for i, cand := range branches {
		candVersion := regexp.MustCompile(`^remotes/\S+/`).ReplaceAllString(cand, "")
		if cand == branch || isFeatureBranch(conf, candVersion) {
			continue
		}
		wg.Add(1)
		go func(i int, cand, candVersion string) {
			defer wg.Done()
			out, ok := runGit(dir, "rev-list", cand+".."+branch)
			results[i] = res{i, out, ok, candVersion}
		}(i, cand, candVersion)
	}
	wg.Wait()
	length := int(^uint(0) >> 1)
	lastIndex := -1
	for _, r := range results {
		if !r.ok {
			continue
		}
		if len(r.out) < length || (len(r.out) == length && lastIndex < r.index) {
			lastIndex = r.index
			length = len(r.out)
			version = semver.NormalizeBranch(r.cand)
			pretty = "dev-" + r.cand
		}
	}
	return version, pretty
}
