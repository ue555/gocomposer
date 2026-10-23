// Command integration installs fixed real projects (a library, a CLI app and
// a Laravel application, see testdata/integration/projects.json) with the
// official Composer (oracle image, pinned Composer f83588392df4) and with
// gocomposer, each in an independent copy, and compares the observable
// results: exit codes, output, composer.lock, vendor tree (paths, modes,
// content hashes), the autoloader (class map resolution via the generated
// ClassLoader, InstalledVersions) and project commands run with official PHP
// (e.g. `php artisan about`). Network access is required (Packagist,
// GitHub codeload); unauthenticated GitHub API zipball downloads are rate
// limited, so dists of the larger projects are seeded into a fresh per-run
// cache from codeload (Composer's own cache layout and keys) and one small
// project downloads everything live. Results and raw logs go to
// docs/validation/<UTC timestamp>-integration/.
package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	oracleImage = "localhost/gocomposer-oracle:8.4"
	upstreamSHA = "f83588392df43d4fd027d84d26de459946ef4e80"
)

type project struct {
	Name    string     `json:"name"`
	Git     string     `json:"git"`
	Ref     string     `json:"ref"`
	Commit  string     `json:"commit"`
	Dir     string     `json:"dir"`
	Live    bool       `json:"live"`
	Laravel bool       `json:"laravel"`
	Run     [][]string `json:"run"`
	// LockCommand is the official Composer command generating the lock
	// (default: update --no-install --no-scripts).
	LockCommand string `json:"lockCommand"`
	// Steps are further comparisons: command sequences run on a copy of the
	// project with the fixed lock (e.g. installing Symfony Flex recipes).
	Steps []customStep `json:"steps"`
	// TreeSkip are top-level directories left out of the tree comparison
	// (caches written by the project's own commands).
	TreeSkip []string `json:"treeSkip"`
	// Normalize replaces random values in files before hashing them.
	Normalize []normRule `json:"normalize"`
}

type customStep struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
	Runs     bool     `json:"runs"`
}

type normRule struct {
	// Glob matches the file name (path.Match on the base name).
	Glob    string `json:"glob"`
	Regex   string `json:"regex"`
	Replace string `json:"replace"`
}

type stepResult struct {
	Tool    string            `json:"tool"`
	Step    string            `json:"step"`
	Command string            `json:"command"`
	Exit    int               `json:"exit"`
	Stdout  string            `json:"stdout"`
	Stderr  string            `json:"stderr"`
	Lock    string            `json:"lock,omitempty"`
	Tree    []string          `json:"tree,omitempty"`
	Probe   string            `json:"probe,omitempty"`
	Runs    map[string]string `json:"runs,omitempty"`
}

type comparison struct {
	Project     string   `json:"project"`
	Step        string   `json:"step"`
	Status      string   `json:"status"`
	Differences []string `json:"differences,omitempty"`
	OracleExit  int      `json:"oracle_exit"`
	GoExit      int      `json:"go_exit"`
}

type tool struct{ name, cmd string }

var tools = []tool{{"oracle", "php /oracle-source/bin/composer"}, {"go", "/gobin/gocomposer"}}

func main() {
	root := flag.String("root", ".", "repository root")
	only := flag.String("only", "", "comma separated project names")
	keep := flag.Bool("keep", false, "keep work directories")
	genLocks := flag.Bool("generate-locks", false, "(re)generate missing composer.lock files with the official Composer")
	flag.Parse()
	abs, _ := filepath.Abs(*root)
	if err := run(abs, *only, *keep, *genLocks); err != nil {
		fmt.Fprintln(os.Stderr, "integration:", err)
		os.Exit(1)
	}
}

func run(root, only string, keep, genLocks bool) error {
	stamp := time.Now().UTC().Format("20060102T150405Z")
	outDir := filepath.Join(root, "docs", "validation", stamp+"-integration")
	work := filepath.Join(root, ".cache", "integration", stamp)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	if !keep {
		defer os.RemoveAll(work)
	}
	var cfg struct {
		Projects []project `json:"projects"`
	}
	raw, err := os.ReadFile(filepath.Join(root, "testdata", "integration", "projects.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	arch, err := podmanArch()
	if err != nil {
		return fmt.Errorf("未比較: oracle environment unavailable: %w", err)
	}
	gobin := filepath.Join(work, "gobin")
	if err := os.MkdirAll(gobin, 0o755); err != nil {
		return err
	}
	build := exec.Command("go", "build", "-trimpath", "-o", filepath.Join(gobin, "gocomposer"), "./cmd/gocomposer")
	build.Dir = root
	build.Env = append(os.Environ(), "GOOS=linux", "GOARCH="+arch, "CGO_ENABLED=0", "GOCACHE="+filepath.Join(root, ".cache", "go-build"))
	if out, err := build.CombinedOutput(); err != nil {
		return fmt.Errorf("building linux binary: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(gobin, "probe.php"), []byte(probePHP), 0o644); err != nil {
		return err
	}
	want := map[string]bool{}
	for _, s := range strings.Split(only, ",") {
		if s = strings.TrimSpace(s); s != "" {
			want[s] = true
		}
	}
	env := environmentInfo(root, gobin)
	var comps []comparison
	allMatch := true
	for _, p := range cfg.Projects {
		if len(want) > 0 && !want[p.Name] {
			continue
		}
		fmt.Fprintf(os.Stderr, "== %s\n", p.Name)
		cs, err := runProject(root, work, outDir, gobin, p, genLocks)
		if err != nil {
			cs = append(cs, comparison{Project: p.Name, Step: "setup", Status: "ERROR", Differences: []string{err.Error()}})
		}
		for _, c := range cs {
			fmt.Fprintf(os.Stderr, "   %-28s %s\n", c.Step, c.Status)
			for _, d := range c.Differences {
				fmt.Fprintf(os.Stderr, "     - %s\n", firstLines(d, 10))
			}
			if c.Status != "MATCH" {
				allMatch = false
			}
		}
		comps = append(comps, cs...)
	}
	summary := map[string]any{"environment": env, "results": comps, "all_match": allMatch, "generated": stamp}
	b, _ := json.MarshalIndent(summary, "", "  ")
	if err := os.WriteFile(filepath.Join(outDir, "summary.json"), b, 0o644); err != nil {
		return err
	}
	if err := writeMarkdown(filepath.Join(outDir, "summary.md"), env, comps); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "results:", outDir)
	if !allMatch {
		return errors.New("differences found")
	}
	return nil
}

func podmanArch() (string, error) {
	out, err := exec.Command("podman", "image", "inspect", oracleImage, "--format", "{{.Architecture}}").Output()
	if err != nil {
		return "", fmt.Errorf("podman image %s not available (%v)", oracleImage, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func environmentInfo(root, gobin string) map[string]string {
	info := map[string]string{"upstream_commit": upstreamSHA, "oracle_image": oracleImage, "host_os": runtime.GOOS + "/" + runtime.GOARCH}
	if out, err := exec.Command("podman", "image", "inspect", oracleImage, "--format", "{{.Digest}}").Output(); err == nil {
		info["oracle_image_digest"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("git", "-C", filepath.Join(root, ".cache", "oracle-source"), "rev-parse", "HEAD").Output(); err == nil {
		info["oracle_source_head"] = strings.TrimSpace(string(out))
	}
	if b, err := os.ReadFile(filepath.Join(gobin, "gocomposer")); err == nil {
		sum := sha256.Sum256(b)
		info["gocomposer_linux_sha256"] = hex.EncodeToString(sum[:])
	}
	if b, err := os.ReadFile(filepath.Join(root, ".cache", "oracle-source", "composer.lock")); err == nil {
		sum := sha256.Sum256(b)
		info["oracle_composer_lock_sha256"] = hex.EncodeToString(sum[:])
	}
	return info
}

// prepareSource creates the pristine project directory.
func prepareSource(root, dst string, p project) error {
	if p.Git != "" {
		cmd := exec.Command("git", "clone", "--quiet", "--depth", "1", "--branch", p.Ref, p.Git, dst)
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git clone %s %s: %v\n%s", p.Git, p.Ref, err, out)
		}
		head, err := exec.Command("git", "-C", dst, "rev-parse", "HEAD").Output()
		if err != nil {
			return err
		}
		if got := strings.TrimSpace(string(head)); got != p.Commit {
			return fmt.Errorf("%s %s resolved to %s, expected pinned commit %s", p.Git, p.Ref, got, p.Commit)
		}
		return nil
	}
	return copyTree(filepath.Join(root, "testdata", "integration", p.Dir), dst)
}

func lockPath(root string, p project) string {
	return filepath.Join(root, "testdata", "integration", "locks", p.Name+".lock")
}

func runProject(root, work, outDir, gobin string, p project, genLocks bool) ([]comparison, error) {
	base := filepath.Join(work, p.Name)
	src := filepath.Join(base, "src")
	if err := prepareSource(root, src, p); err != nil {
		return nil, err
	}
	if p.Laravel {
		if err := prepareLaravel(src); err != nil {
			return nil, err
		}
	}
	lock := lockPath(root, p)
	if _, err := os.Stat(lock); err != nil {
		if !genLocks {
			return nil, fmt.Errorf("missing %s (run with -generate-locks once and commit it)", lock)
		}
		if err := generateLock(root, base, src, gobin, lock, p.LockCommand); err != nil {
			return nil, err
		}
	}
	lockData, err := os.ReadFile(lock)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(src, "composer.lock"), lockData, 0o644); err != nil {
		return nil, err
	}
	seed := filepath.Join(base, "seed-cache")
	if !p.Live {
		if err := seedCache(seed, lockData); err != nil {
			return nil, err
		}
	}
	var comps []comparison
	steps := []struct {
		name, command string
		useLock       bool
		withRuns      bool
	}{
		{"install", "install --no-progress --no-interaction", true, true},
		{"update-no-install", "update --no-progress --no-install --no-interaction", true, false},
	}
	for _, cs := range p.Steps {
		steps = append(steps, struct {
			name, command string
			useLock       bool
			withRuns      bool
		}{cs.Name, strings.Join(cs.Commands, "\n"), true, cs.Runs})
	}
	for _, st := range steps {
		var results []stepResult
		for _, t := range tools {
			dir := filepath.Join(base, st.name, t.name)
			if err := copyTree(src, filepath.Join(dir, "project")); err != nil {
				return nil, err
			}
			if !st.useLock {
				os.Remove(filepath.Join(dir, "project", "composer.lock"))
			}
			if !p.Live {
				if err := copyTree(seed, filepath.Join(dir, "cache")); err != nil {
					return nil, err
				}
			}
			r, err := runTool(root, dir, gobin, t, st.command, p, st.withRuns)
			if err != nil {
				return nil, err
			}
			r.Step = st.name
			results = append(results, r)
			b, _ := json.MarshalIndent(r, "", "  ")
			_ = os.MkdirAll(filepath.Join(outDir, p.Name), 0o755)
			_ = os.WriteFile(filepath.Join(outDir, p.Name, st.name+"-"+t.name+".json"), b, 0o644)
		}
		comps = append(comps, compare(p.Name, st.name, results[0], results[1]))
	}
	return comps, nil
}

func prepareLaravel(src string) error {
	ex, err := os.ReadFile(filepath.Join(src, ".env.example"))
	if err != nil {
		return err
	}
	env := regexp.MustCompile(`(?m)^APP_KEY=.*$`).ReplaceAllString(string(ex), "APP_KEY=base64:Z29jb21wb3Nlci1pbnRlZ3JhdGlvbi10ZXN0LWtleSE=")
	return os.WriteFile(filepath.Join(src, ".env"), []byte(env), 0o644)
}

func generateLock(root, base, src, gobin, lock, command string) error {
	dir := filepath.Join(base, "lockgen")
	if err := copyTree(src, filepath.Join(dir, "project")); err != nil {
		return err
	}
	if command == "" {
		command = "update --no-progress --no-install --no-scripts --no-interaction"
	}
	r, err := runTool(root, dir, gobin, tools[0], command, project{}, false)
	if err != nil {
		return err
	}
	if r.Exit != 0 {
		return fmt.Errorf("official composer update failed:\n%s", r.Stderr)
	}
	if err := os.MkdirAll(filepath.Dir(lock), 0o755); err != nil {
		return err
	}
	return os.WriteFile(lock, []byte(r.Lock), 0o644)
}

var reGitHubZipball = regexp.MustCompile(`^https://api\.github\.com/repos/([^/]+)/([^/]+)/zipball/([^/]+)$`)

// seedCache downloads the GitHub zipballs of the lock from codeload (not
// subject to the API rate limit) into Composer's files cache layout:
// files/<name>/<sha1(dist url)>.zip.
func seedCache(dir string, lockData []byte) error {
	var lock struct {
		Packages    []lockPkg `json:"packages"`
		PackagesDev []lockPkg `json:"packages-dev"`
	}
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return err
	}
	client := &http.Client{Timeout: 120 * time.Second}
	for _, p := range append(lock.Packages, lock.PackagesDev...) {
		if p.Dist.Type != "zip" {
			continue
		}
		m := reGitHubZipball.FindStringSubmatch(p.Dist.URL)
		if m == nil {
			continue
		}
		sum := sha1.Sum([]byte(p.Dist.URL))
		target := filepath.Join(dir, "files", p.Name, hex.EncodeToString(sum[:])+".zip")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		u := "https://codeload.github.com/" + m[1] + "/" + m[2] + "/legacy.zip/" + m[3]
		if err := fetch(client, u, target); err != nil {
			return fmt.Errorf("seeding %s: %w", p.Name, err)
		}
	}
	return nil
}

type lockPkg struct {
	Name string `json:"name"`
	Dist struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"dist"`
}

func fetch(client *http.Client, u, target string) error {
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := client.Get(u)
		if err != nil {
			last = err
			time.Sleep(time.Second)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			last = fmt.Errorf("%s: HTTP %d", u, resp.StatusCode)
			time.Sleep(time.Second)
			continue
		}
		f, err := os.Create(target)
		if err != nil {
			resp.Body.Close()
			return err
		}
		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err == nil {
			return nil
		}
		last = err
	}
	return last
}

func runTool(root, dir, gobin string, t tool, command string, p project, withRuns bool) (stepResult, error) {
	for _, d := range []string{"home", "cache", "out"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o777); err != nil {
			return stepResult{}, err
		}
	}
	// a step may run several commands (one per line): their output is
	// concatenated, the exit code is the first non-zero one
	script := "cd /work/project\n: >/work/out/stdout; : >/work/out/stderr; echo 0 >/work/out/exit\n"
	cmds := strings.Split(command, "\n")
	for _, c := range cmds {
		if len(cmds) > 1 {
			script += "echo '## " + c + "' >>/work/out/stdout; echo '## " + c + "' >>/work/out/stderr\n"
		}
		script += t.cmd + " " + c + " >>/work/out/stdout 2>>/work/out/stderr; rc=$?; if [ $rc -ne 0 ] && [ \"$(cat /work/out/exit)\" = 0 ]; then echo $rc >/work/out/exit; fi\n"
	}
	script += "if [ -f vendor/autoload.php ]; then php /gobin/probe.php >/work/out/probe 2>&1; fi\n"
	if withRuns {
		for i, r := range p.Run {
			script += "if [ -f vendor/autoload.php ]; then " + shellJoin(r) + fmt.Sprintf(" >/work/out/run%d 2>&1; echo \"exit=$?\" >>/work/out/run%d; fi\n", i, i)
		}
	}
	args := podmanRun(
		"-v", filepath.Join(root, ".cache", "oracle-source")+":/oracle-source:ro",
		"-v", gobin+":/gobin:ro",
		"-v", dir+":/work",
		"-e", "COMPOSER_HOME=/work/home",
		"-e", "COMPOSER_CACHE_DIR=/work/cache",
		"-e", "COMPOSER_NO_INTERACTION=1",
		"-e", "COLUMNS=80",
		"-e", "TZ=UTC",
		"-e", "GIT_CONFIG_COUNT=1",
		"-e", "GIT_CONFIG_KEY_0=safe.directory",
		"-e", "GIT_CONFIG_VALUE_0=*",
		oracleImage, "sh", "-c", script)
	cmd := exec.Command("podman", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if out, err := cmd.Output(); err != nil {
		return stepResult{}, fmt.Errorf("podman run failed: %v\n%s%s", err, out, stderr.String())
	}
	read := func(name string) string { b, _ := os.ReadFile(filepath.Join(dir, "out", name)); return string(b) }
	r := stepResult{Tool: t.name, Command: command, Stdout: read("stdout"), Stderr: read("stderr"), Probe: read("probe")}
	fmt.Sscanf(strings.TrimSpace(read("exit")), "%d", &r.Exit)
	if b, err := os.ReadFile(filepath.Join(dir, "project", "composer.lock")); err == nil {
		r.Lock = string(b)
	}
	r.Tree = treeOf(filepath.Join(dir, "project"), p)
	if withRuns {
		r.Runs = map[string]string{}
		for i, rr := range p.Run {
			r.Runs[strings.Join(rr, " ")] = read(fmt.Sprintf("run%d", i))
		}
	}
	return r, nil
}

func shellJoin(args []string) string {
	var out []string
	for _, a := range args {
		out = append(out, "'"+strings.ReplaceAll(a, "'", `'\''`)+"'")
	}
	return strings.Join(out, " ")
}

// treeOf lists vendor/ and the files commands may generate (paths, modes and
// content hashes); .git and node_modules are skipped.
func treeOf(dir string, proj project) []string {
	var out []string
	var rules []struct {
		glob string
		re   *regexp.Regexp
		repl string
	}
	for _, n := range proj.Normalize {
		rules = append(rules, struct {
			glob string
			re   *regexp.Regexp
			repl string
		}{n.Glob, regexp.MustCompile(n.Regex), n.Replace})
	}
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		if rel == "." {
			return nil
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		for _, skip := range proj.TreeSkip {
			if d.IsDir() && rel == skip {
				return filepath.SkipDir
			}
		}
		info, err := os.Lstat(p)
		if err != nil {
			return nil
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			l, _ := os.Readlink(p)
			out = append(out, fmt.Sprintf("L %s -> %s", rel, l))
		case info.IsDir():
			out = append(out, "D "+rel)
		default:
			b, _ := os.ReadFile(p)
			for _, r := range rules {
				if ok, _ := filepath.Match(r.glob, d.Name()); ok {
					b = r.re.ReplaceAll(b, []byte(r.repl))
				}
			}
			sum := sha256.Sum256(b)
			out = append(out, fmt.Sprintf("F %s %o %s", rel, info.Mode().Perm(), hex.EncodeToString(sum[:])[:16]))
		}
		return nil
	})
	sort.Strings(out)
	return out
}

var (
	reTiming   = regexp.MustCompile(`\d+\.\d+ seconds`)
	reInFileLn = regexp.MustCompile(`(?m)^In [A-Za-z]+\.php line \d+:$`)
	reInGo     = regexp.MustCompile(`(?m)^In gocomposer:$`)
	// symfony/flex prints a random session id (bin2hex(random_bytes(16)))
	reFlexSession = regexp.MustCompile(`(Symfony operations: \d+ recipes? \()[0-9a-f]{32}\)`)
)

func normalize(s string) string {
	s = reTiming.ReplaceAllString(s, "N seconds")
	s = reInFileLn.ReplaceAllString(s, "In <location>:")
	s = reFlexSession.ReplaceAllString(s, "${1}SESSION)")
	return reInGo.ReplaceAllString(s, "In <location>:")
}

func compare(project, step string, o, g stepResult) comparison {
	c := comparison{Project: project, Step: step, OracleExit: o.Exit, GoExit: g.Exit}
	if o.Exit != g.Exit {
		c.Differences = append(c.Differences, fmt.Sprintf("exit code: oracle=%d go=%d", o.Exit, g.Exit))
	}
	for _, f := range []struct {
		name string
		a, b string
	}{{"stdout", o.Stdout, g.Stdout}, {"stderr", o.Stderr, g.Stderr}, {"composer.lock", o.Lock, g.Lock}, {"autoload probe", o.Probe, g.Probe}, {"project tree", strings.Join(o.Tree, "\n"), strings.Join(g.Tree, "\n")}} {
		if d := diffText(normalize(f.a), normalize(f.b)); d != "" {
			c.Differences = append(c.Differences, f.name+":\n"+d)
		}
	}
	keys := make([]string, 0, len(o.Runs))
	for k := range o.Runs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if d := diffText(normalize(o.Runs[k]), normalize(g.Runs[k])); d != "" {
			c.Differences = append(c.Differences, "run `"+k+"`:\n"+d)
		}
	}
	if len(c.Differences) == 0 {
		c.Status = "MATCH"
	} else {
		c.Status = "DIFF"
	}
	return c
}

// diffText is a small line diff (LCS) listing removed/added lines.
func diffText(a, b string) string {
	if a == b {
		return ""
	}
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	if len(al)*len(bl) > 25_000_000 {
		return multisetDiff(al, bl)
	}
	n, m := len(al), len(bl)
	dp := make([][]int32, n+1)
	for i := range dp {
		dp[i] = make([]int32, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if al[i] == bl[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	var out []string
	i, j := 0, 0
	for i < n || j < m {
		switch {
		case i < n && j < m && al[i] == bl[j]:
			i++
			j++
		case j < m && (i == n || dp[i][j+1] >= dp[i+1][j]):
			out = append(out, "+ (go) "+bl[j])
			j++
		default:
			out = append(out, "- (oracle) "+al[i])
			i++
		}
	}
	return strings.Join(out, "\n")
}

// multisetDiff lists lines present only on one side (used for large,
// order-insensitive inputs such as sorted trees).
func multisetDiff(al, bl []string) string {
	count := map[string]int{}
	for _, l := range al {
		count[l]++
	}
	for _, l := range bl {
		count[l]--
	}
	var out []string
	for _, l := range al {
		if count[l] > 0 {
			out = append(out, "- (oracle) "+l)
			count[l]--
		}
	}
	for _, l := range bl {
		if count[l] < 0 {
			out = append(out, "+ (go) "+l)
			count[l]++
		}
	}
	return strings.Join(out, "\n")
}

func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append(lines[:n], fmt.Sprintf("... (%d more lines)", len(lines)-n))
	}
	return strings.Join(lines, "\n     ")
}

func writeMarkdown(path string, env map[string]string, comps []comparison) error {
	var sb strings.Builder
	sb.WriteString("# Composer vs gocomposer integration (fixed real projects)\n\n")
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString("- " + k + ": `" + env[k] + "`\n")
	}
	sb.WriteString("\n| project | step | status | oracle exit | go exit |\n|---|---|---|---|---|\n")
	for _, c := range comps {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d |\n", c.Project, c.Step, c.Status, c.OracleExit, c.GoExit))
	}
	for _, c := range comps {
		if len(c.Differences) == 0 {
			continue
		}
		sb.WriteString("\n## " + c.Project + " / " + c.Step + "\n\n```\n")
		for _, d := range c.Differences {
			sb.WriteString(firstLines(d, 200) + "\n")
		}
		sb.WriteString("```\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		target := filepath.Join(dst, rel)
		if info.Mode()&os.ModeSymlink != 0 {
			l, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(l, target)
		}
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm()|0o700)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}
