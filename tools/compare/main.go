// Command compare runs official Composer (pinned source, official PHP image)
// and gocomposer on identical copies of fixture projects and diffs the
// observable results. Results are written below docs/validation/.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
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
	upstreamSHA   = "f83588392df43d4fd027d84d26de459946ef4e80"
	containerRoot = "/work"
)

type tool struct {
	name string
	cmd  string
}

var tools = []tool{
	{"oracle", "php /oracle-source/bin/composer"},
	{"go", "/gobin/gocomposer"},
}

type result struct {
	Scenario string            `json:"scenario"`
	Tool     string            `json:"tool"`
	Exit     int               `json:"exit"`
	Stdout   string            `json:"stdout"`
	Stderr   string            `json:"stderr"`
	Lock     string            `json:"lock"`
	Tree     []string          `json:"tree"`
	Probe    string            `json:"probe"`
	Files    map[string]string `json:"files"`
	// GitState is the state of every .git checkout below the project (HEAD,
	// branch, status, remotes, stashes); .git internals are excluded from the
	// compared tree (Tree keeps the raw, unnormalized manifest).
	GitState string `json:"git_state,omitempty"`
	// TTYRaw is the raw terminal transcript of a TTY scenario (tty.go).
	TTYRaw []byte `json:"tty_raw,omitempty"`
}

type comparison struct {
	Scenario    string   `json:"scenario"`
	Status      string   `json:"status"` // MATCH, DIFF, KNOWN-DIFF, ERROR
	KnownReason string   `json:"known_reason,omitempty"`
	Differences []string `json:"differences"`
	OracleExit  int      `json:"oracle_exit"`
	GoExit      int      `json:"go_exit"`
}

// oracleImage is the official PHP image (PHP 8.4 is the primary oracle;
// 8.3/8.5 images are regression targets for autoload and platform checks).
var oracleImage = "localhost/gocomposer-oracle:8.4"

func main() {
	root := flag.String("root", ".", "gocomposer repository root")
	only := flag.String("only", "", "comma separated scenario names")
	out := flag.String("out", "", "output dir (default docs/validation/<timestamp>-compare)")
	keep := flag.Bool("keep", false, "keep work directories")
	image := flag.String("image", oracleImage, "oracle image (localhost/gocomposer-oracle:8.3|8.4|8.5)")
	flag.Parse()
	oracleImage = *image
	absRoot, _ := filepath.Abs(*root)
	if err := run(absRoot, *only, *out, *keep); err != nil {
		fmt.Fprintln(os.Stderr, "compare:", err)
		os.Exit(2)
	}
}

func run(root, only, outDir string, keep bool) error {
	stamp := time.Now().UTC().Format("20060102T150405Z")
	if outDir == "" {
		outDir = filepath.Join(root, "docs", "validation", stamp+"-compare")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	work := filepath.Join(root, ".cache", "compare", stamp)
	if err := os.MkdirAll(work, 0o755); err != nil {
		return err
	}
	if !keep {
		defer os.RemoveAll(work)
	}
	arch, err := podmanArch()
	if err != nil {
		return fmt.Errorf("oracle environment unavailable: %w", err)
	}
	gobin := filepath.Join(work, "gobin")
	if err := os.MkdirAll(gobin, 0o755); err != nil {
		return err
	}
	build := exec.Command("go", "build", "-trimpath", "-o", filepath.Join(gobin, "gocomposer"), "./cmd/gocomposer")
	build.Dir = root
	build.Env = append(os.Environ(), "GOOS=linux", "GOARCH="+arch, "CGO_ENABLED=0", "GOCACHE="+filepath.Join(root, ".cache", "go-build"))
	if outb, err := build.CombinedOutput(); err != nil {
		return fmt.Errorf("building linux binary: %v\n%s", err, outb)
	}
	if err := os.WriteFile(filepath.Join(gobin, "probe.php"), []byte(probePHP), 0o644); err != nil {
		return err
	}
	fixtures := filepath.Join(work, "fixtures")
	if err := buildRepo(filepath.Join(root, "testdata", "compare", "packages"), fixtures); err != nil {
		return err
	}
	// mount point for the per-scenario git fixtures (see gitfixtures.go)
	if err := os.MkdirAll(filepath.Join(fixtures, "git"), 0o755); err != nil {
		return err
	}
	scenarios := allScenarios()
	if only != "" {
		want := map[string]bool{}
		for _, s := range strings.Split(only, ",") {
			want[strings.TrimSpace(s)] = true
		}
		var sel []Scenario
		for _, s := range scenarios {
			match := want[s.Name]
			for w := range want {
				// "prefix*" selects every scenario with that prefix
				if strings.HasSuffix(w, "*") && strings.HasPrefix(s.Name, strings.TrimSuffix(w, "*")) {
					match = true
				}
			}
			if match {
				sel = append(sel, s)
			}
		}
		scenarios = sel
	}
	env := environmentInfo(root, gobin)
	var comps []comparison
	allMatch := true
	for _, sc := range scenarios {
		fmt.Fprintf(os.Stderr, "== %s\n", sc.Name)
		comp, results, err := runScenario(root, work, fixtures, gobin, sc)
		if err != nil {
			comp = comparison{Scenario: sc.Name, Status: "ERROR", Differences: []string{err.Error()}}
		}
		if comp.Status != "MATCH" && comp.Status != "KNOWN-DIFF" {
			allMatch = false
		}
		comps = append(comps, comp)
		raw := filepath.Join(outDir, sc.Name)
		_ = os.MkdirAll(raw, 0o755)
		for _, r := range results {
			b, _ := json.MarshalIndent(r, "", "  ")
			_ = os.WriteFile(filepath.Join(raw, r.Tool+".json"), b, 0o644)
		}
		fmt.Fprintf(os.Stderr, "   %s\n", comp.Status)
		for _, d := range comp.Differences {
			fmt.Fprintf(os.Stderr, "   - %s\n", firstLines(d, 12))
		}
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
		return fmt.Errorf("differences found")
	}
	return nil
}

func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append(lines[:n], fmt.Sprintf("... (%d more lines)", len(lines)-n))
	}
	return strings.Join(lines, "\n     ")
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
	if out, err := exec.Command("podman", "image", "inspect", oracleImage+zipImageSuffix, "--format", "{{.Digest}}").Output(); err == nil {
		info["oracle_zip_image"] = oracleImage + zipImageSuffix
		info["oracle_zip_image_digest"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("podman", "image", "inspect", oracleImage+"-arch", "--format", "{{.Digest}}").Output(); err == nil {
		info["oracle_arch_image"] = oracleImage + "-arch"
		info["oracle_arch_image_digest"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("podman", podmanRun("--network=none", oracleImage, "php", "-r", "echo PHP_VERSION;")...).Output(); err == nil {
		info["php_version"] = strings.TrimSpace(string(out))
	}
	if b, err := os.ReadFile(filepath.Join(gobin, "gocomposer")); err == nil {
		sum := sha256.Sum256(b)
		info["gocomposer_sha256"] = hex.EncodeToString(sum[:])
	}
	if b, err := os.ReadFile(filepath.Join(root, ".cache", "oracle-source", "composer.lock")); err == nil {
		sum := sha256.Sum256(b)
		info["oracle_composer_lock_sha256"] = hex.EncodeToString(sum[:])
	}
	if out, err := exec.Command("git", "-C", filepath.Join(root, ".cache", "oracle-source"), "rev-parse", "HEAD").Output(); err == nil {
		info["oracle_source_head"] = strings.TrimSpace(string(out))
	}
	return info
}

// Scenario describes one comparison.
type Scenario struct {
	// RacyRemovalBar marks scenarios whose uninstall phase draws a progress
	// bar over Composer's async `rm -rf` processes: how many are still
	// running when Loop::wait starts is timing dependent in Composer itself
	// (usually none, giving the indeterminate bar). The bar drawn right
	// after the "Removing" lines is compared in its indeterminate form.
	RacyRemovalBar bool
	Name           string
	JSON           string
	Files          map[string]string
	Setup          []string // oracle command args run before the snapshot (each a whole command line)
	Command        string
	Probe          bool
	Env            map[string]string
	Snapshot       []string          // extra files to capture
	After          map[string]string // files (re)written after Setup, before the compared command
	// KnownDiff documents a reviewed divergence (e.g. Composer itself
	// crashing with a PHP TypeError). Differences are still recorded but
	// reported as KNOWN-DIFF; the reason must also be in docs/compatibility.md.
	KnownDiff string
	// NoJSON leaves composer.json out of the project.
	NoJSON bool
	// GitFixtures names git fixture build scripts (testdata/compare/vcs/<name>.sh)
	// run inside the container before Setup; the resulting bare repositories
	// are mounted read-only at /fixtures/git for every run of the scenario.
	GitFixtures []string
	// GitAfter is a shell snippet (with testdata/compare/vcs/lib.sh loaded)
	// run against /fixtures/git after Setup, before the compared runs.
	GitAfter string
	// Image overrides the container image for this scenario's runs (both
	// tools), e.g. the ext-zip oracle variant built from oracle/Containerfile.zip.
	Image string
	// TTY runs the compared command on a pseudo terminal (interactive mode)
	// and feeds it Input (tty.go).
	TTY   bool
	Input []TTYStep
	// Root runs the compared command as uid 0 (podmanRunRoot) with the
	// container markers Platform::isDocker looks for removed, so Composer's
	// super user checks apply (scenarios_term.go).
	Root bool
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

func runScenario(root, work, fixtures, gobin string, sc Scenario) (comparison, []result, error) {
	base := filepath.Join(work, sc.Name)
	proj := filepath.Join(base, "template", "project")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		return comparison{}, nil, err
	}
	if !sc.NoJSON {
		if err := os.WriteFile(filepath.Join(proj, "composer.json"), []byte(sc.JSON), 0o644); err != nil {
			return comparison{}, nil, err
		}
	}
	for name, content := range sc.Files {
		p := filepath.Join(proj, name)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return comparison{}, nil, err
		}
	}
	if len(sc.GitFixtures) > 0 {
		var script strings.Builder
		for _, name := range sc.GitFixtures {
			script.WriteString("( . /vcsdefs/" + name + ".sh )\n")
		}
		if err := runGitFixtureScript(root, base, script.String()); err != nil {
			return comparison{}, nil, err
		}
	}
	for _, setup := range sc.Setup {
		r, err := runTool(root, filepath.Join(base, "template"), fixtures, gobin, tools[0], setup, sc, false)
		if err != nil {
			return comparison{}, nil, err
		}
		if r.Exit != 0 {
			return comparison{Scenario: sc.Name, Status: "ERROR", Differences: []string{"setup failed: " + setup + "\n" + r.Stderr}}, nil, nil
		}
	}
	for name, content := range sc.After {
		p := filepath.Join(proj, name)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return comparison{}, nil, err
		}
	}
	if sc.GitAfter != "" {
		if err := runGitFixtureScript(root, base, sc.GitAfter); err != nil {
			return comparison{}, nil, err
		}
	}
	var results []result
	for _, t := range tools {
		dir := filepath.Join(base, t.name)
		_ = os.RemoveAll(dir)
		if err := copyTree(filepath.Join(base, "template"), dir); err != nil {
			return comparison{}, nil, err
		}
		_ = os.RemoveAll(filepath.Join(dir, "home"))
		_ = os.RemoveAll(filepath.Join(dir, "cache"))
		r, err := runTool(root, dir, fixtures, gobin, t, sc.Command, sc, sc.Probe)
		if err != nil {
			return comparison{}, nil, err
		}
		results = append(results, r)
	}
	return compare(sc, results[0], results[1]), results, nil
}

func runTool(root, dir, fixtures, gobin string, t tool, command string, sc Scenario, probe bool) (result, error) {
	for _, d := range []string{"home", "cache", "out"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o777); err != nil {
			return result{}, err
		}
	}
	ttyRun := isTTYRun(sc, dir)
	rootRun := sc.Root && filepath.Base(dir) != "template"
	runLine := t.cmd + " " + command + ` >/work/out/stdout 2>/work/out/stderr; echo $? >/work/out/exit`
	if ttyRun {
		var err error
		if runLine, err = ttyCommandLine(root, gobin, dir, t, command, sc); err != nil {
			return result{}, err
		}
	}
	script := `php -S 127.0.0.1:8899 -t /fixtures >/dev/null 2>&1 &
i=0; while ! php -r 'exit(@fsockopen("127.0.0.1", 8899) ? 0 : 1);' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done
` + apiStubScript + `
` + faultStubScript + `
cd /work/project
` + runLine + `
`
	if rootRun {
		// podman bind-mounts /run/.containerenv; unmounting it needs
		// CAP_SYS_ADMIN (granted by podmanRunRoot)
		script = "umount /run/.containerenv 2>/dev/null; rm -f /run/.containerenv /.dockerenv\n" + script
	}
	if probe {
		script += `if [ -f vendor/autoload.php ]; then php /gobin/probe.php >/work/out/probe 2>&1; fi
`
	}
	script += gitStateScript
	run := podmanRun
	if rootRun {
		run = podmanRunRoot
	}
	args := run("--network=none",
		"-v", fixtures+":/fixtures:ro",
		"-v", filepath.Join(root, ".cache", "oracle-source")+":/oracle-source:ro",
		"-v", gobin+":/gobin:ro",
		"-v", dir+":/work",
	)
	if gitfx := filepath.Join(filepath.Dir(dir), "gitfx"); isDirPath(gitfx) {
		args = append(args, "-v", gitfx+":/fixtures/git:ro")
		if isDirPath(filepath.Join(gitfx, "_api")) {
			args = append(args, forgePodmanArgs...)
			// further host names of the stub (_api/hosts, scenarios_final_cli.go)
			args = append(args, apiHostArgs(gitfx)...)
		}
	}
	args = append(args,
		"-e", "COMPOSER_HOME=/work/home",
		"-e", "COMPOSER_CACHE_DIR=/work/cache",
	)
	if !ttyRun {
		args = append(args, "-e", "COMPOSER_NO_INTERACTION=1")
	}
	args = append(args,
		// Cache::gcIsNecessary runs the files cache GC with probability 1/51
		// per process; it creates and prunes the cache directory, which
		// scenarios observe. In this Composer SHA COMPOSER_TEST_SUITE only
		// disables that GC (gocomposer honours it the same way).
		"-e", "COMPOSER_TEST_SUITE=1",
		"-e", "COLUMNS=80",
		"-e", "TZ=UTC",
	)
	keys := make([]string, 0, len(sc.Env))
	for k := range sc.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		args = append(args, "-e", k+"="+sc.Env[k])
	}
	image := scenarioImage(sc)
	args = append(args, image, "sh", "-c", script)
	cmd := exec.Command("podman", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if out, err := cmd.Output(); err != nil {
		return result{}, fmt.Errorf("podman run failed: %v\n%s%s", err, out, stderr.String())
	}
	r := result{Scenario: sc.Name, Tool: t.name, Files: map[string]string{}}
	read := func(name string) string { b, _ := os.ReadFile(filepath.Join(dir, "out", name)); return string(b) }
	r.Stdout = read("stdout")
	r.Stderr = read("stderr")
	fmt.Sscanf(strings.TrimSpace(read("exit")), "%d", &r.Exit)
	r.Probe = read("probe")
	r.GitState = read("gitstate")
	if ttyRun {
		readTTYResult(dir, &r)
	}
	if b, err := os.ReadFile(filepath.Join(dir, "project", "composer.lock")); err == nil {
		r.Lock = string(b)
	}
	r.Tree = manifest(filepath.Join(dir, "project"))
	for _, f := range append([]string{"composer.json", "vendor/composer/installed.json", "vendor/composer/installed.php", "vendor/composer/autoload_static.php", "vendor/composer/autoload_psr4.php", "vendor/composer/autoload_classmap.php", "vendor/composer/autoload_files.php", "vendor/composer/autoload_real.php", "vendor/composer/platform_check.php"}, sc.Snapshot...) {
		if b, err := os.ReadFile(filepath.Join(dir, "project", f)); err == nil {
			r.Files[f] = string(b)
		}
	}
	return r, nil
}

func manifest(dir string) []string {
	var out []string
	_ = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		if rel == "." {
			return nil
		}
		mode := info.Mode()
		switch {
		case mode&os.ModeSymlink != 0:
			l, _ := os.Readlink(p)
			out = append(out, fmt.Sprintf("L %s -> %s", rel, l))
		case info.IsDir():
			out = append(out, fmt.Sprintf("D %s", rel))
		default:
			b, _ := os.ReadFile(p)
			sum := sha256.Sum256(b)
			out = append(out, fmt.Sprintf("F %s %o %s", rel, mode.Perm(), hex.EncodeToString(sum[:])[:16]))
		}
		return nil
	})
	sort.Strings(out)
	return out
}

var (
	reTiming  = regexp.MustCompile(`\d+\.\d+ seconds`)
	reTmpHash = regexp.MustCompile(`tmp-[0-9a-f]{32}`)
	// (a terminal transcript has the title colored: \e[33m...\e[39m)
	reInFileLn = regexp.MustCompile(`(?m)^((?:\\e\[33m)?)In [A-Za-z]+\.php line \d+:((?:\\e\[39m)?)$`)
	reInGo     = regexp.MustCompile(`(?m)^((?:\\e\[33m)?)In gocomposer:((?:\\e\[39m)?)$`)
	// the token note of GitHub::authorizeOAuthInteractively has the minute
	// it was printed in ("Composer 2026-10-23 0815" url-encoded)
	reGitHubNote = regexp.MustCompile(`Composer\+\d{4}-\d{2}-\d{2}\+\d{4}`)
	// intermediate (throttled, timing dependent) progress bar redraws; the
	// first (0%) and final (100%) states are always printed and compared.
	reProgressMid = regexp.MustCompile(`^ +\d+/\d+ \[=*>-*\] +([1-9]|[1-9]\d)%`)
)

func normalizeOutput(s string) string {
	s = reTiming.ReplaceAllString(s, "N seconds")
	s = reTmpHash.ReplaceAllString(s, "tmp-HASH")
	s = reGitHubNote.ReplaceAllString(s, "Composer+DATE")
	s = reInFileLn.ReplaceAllString(s, "${1}In <location>:${2}")
	s = reInGo.ReplaceAllString(s, "${1}In <location>:${2}")
	for _, f := range extraOutputNormalizers {
		s = f(s)
	}
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		if reProgressMid.MatchString(l) {
			continue
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func compare(sc Scenario, o, g result) comparison {
	c := comparison{Scenario: sc.Name, OracleExit: o.Exit, GoExit: g.Exit}
	if o.Exit != g.Exit {
		c.Differences = append(c.Differences, fmt.Sprintf("exit code: oracle=%d go=%d", o.Exit, g.Exit))
	}
	if d := diffText(normalizeOutput(o.Stdout), normalizeOutput(g.Stdout)); d != "" {
		c.Differences = append(c.Differences, "stdout:\n"+d)
	}
	oErr, gErr := dropOracleRuntimeDeprecations(o.Stderr), g.Stderr
	if sc.RacyRemovalBar {
		oErr, gErr = indeterminateRemovalBar(oErr), indeterminateRemovalBar(gErr)
	}
	if d := diffText(normalizeOutput(oErr), normalizeOutput(gErr)); d != "" {
		c.Differences = append(c.Differences, "stderr:\n"+d)
	}
	if d := diffText(o.Lock, g.Lock); d != "" {
		c.Differences = append(c.Differences, "composer.lock:\n"+d)
	}
	if d := diffText(strings.Join(comparableTree(o.Tree), "\n"), strings.Join(comparableTree(g.Tree), "\n")); d != "" {
		c.Differences = append(c.Differences, "project tree:\n"+d)
	}
	if d := diffText(o.Probe, g.Probe); d != "" {
		c.Differences = append(c.Differences, "autoload probe:\n"+d)
	}
	if d := diffText(o.GitState, g.GitState); d != "" {
		c.Differences = append(c.Differences, "git state:\n"+d)
	}
	keys := map[string]bool{}
	for k := range o.Files {
		keys[k] = true
	}
	for k := range g.Files {
		keys[k] = true
	}
	var names []string
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		if d := diffText(o.Files[k], g.Files[k]); d != "" {
			c.Differences = append(c.Differences, k+":\n"+d)
		}
	}
	if len(c.Differences) == 0 {
		c.Status = "MATCH"
	} else if sc.KnownDiff != "" {
		c.Status = "KNOWN-DIFF"
		c.KnownReason = sc.KnownDiff
	} else {
		c.Status = "DIFF"
	}
	return c
}

// diffText returns a compact line diff ("" when equal).
func diffText(a, b string) string {
	if a == b {
		return ""
	}
	al := strings.Split(a, "\n")
	bl := strings.Split(b, "\n")
	// LCS-based diff for moderate sizes
	n, m := len(al), len(bl)
	if n*m > 4_000_000 {
		return fmt.Sprintf("(texts differ, %d vs %d lines)", n, m)
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
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
		case j < m && (i >= n || dp[i][j+1] >= dp[i+1][j]):
			out = append(out, "+ (go) "+bl[j])
			j++
		default:
			out = append(out, "- (oracle) "+al[i])
			i++
		}
	}
	return strings.Join(out, "\n")
}

func writeMarkdown(path string, env map[string]string, comps []comparison) error {
	var sb strings.Builder
	sb.WriteString("# Composer vs gocomposer comparison\n\n")
	var keys []string
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("- %s: `%s`\n", k, env[k]))
	}
	sb.WriteString("\n| scenario | status | oracle exit | go exit |\n|---|---|---|---|\n")
	for _, c := range comps {
		sb.WriteString(fmt.Sprintf("| %s | %s | %d | %d |\n", c.Scenario, c.Status, c.OracleExit, c.GoExit))
	}
	for _, c := range comps {
		if len(c.Differences) == 0 {
			continue
		}
		sb.WriteString("\n## " + c.Scenario + "\n\n```\n")
		for _, d := range c.Differences {
			sb.WriteString(d + "\n")
		}
		sb.WriteString("```\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

var _ = io.Copy

// zipImageSuffix names the ext-zip variant of an oracle image
// (oracle/Containerfile.zip, built by scripts/prepare-oracle.sh).
const zipImageSuffix = "-zip"

// zipImage marks scenarios whose Composer code path needs ext-zip; they run
// on the ext-zip variant of the selected oracle image (8.3/8.4/8.5).
const zipImage = "+zip"

func scenarioImage(sc Scenario) string {
	switch sc.Image {
	case "":
		return oracleImage
	case zipImage:
		return oracleImage + zipImageSuffix
	}
	return sc.Image
}

// reOracleDeprecation matches PHP deprecation notices raised by Composer's
// own source under newer PHP versions (e.g. PHP 8.5's "Using null as an
// array offset" in src/Composer/Util/Svn.php). They report a PHP runtime
// issue in the oracle, not Composer behaviour, so only these lines (pointing
// into /oracle-source/src/) are dropped from the oracle's stderr; the raw
// output stays in the result JSON. Composer's bundled libraries
// (/oracle-source/vendor/, e.g. symfony/console's ArgvInput under PHP 8.5)
// are treated the same; gocomposer's PHP runtimes skip their copies' notices.
var reOracleDeprecation = regexp.MustCompile(`^Deprecation Notice: .* in /oracle-source/(?:src|vendor)/\S+:\d+$`)

func dropOracleRuntimeDeprecations(s string) string {
	var lines []string
	dropped, kept := 0, 0
	inTrace := false
	for _, l := range strings.Split(s, "\n") {
		if reOracleDeprecation.MatchString(l) && !isComposerUserDeprecation(l) {
			dropped++
			inTrace = true
			continue
		}
		// the stack trace printed after a dropped notice at -v and above
		if inTrace && (l == "Stack trace:" || strings.HasPrefix(l, " /")) {
			continue
		}
		inTrace = false
		if strings.HasPrefix(l, "Deprecation Notice: ") {
			kept++
		}
		lines = append(lines, l)
	}
	if dropped > 0 && kept == 0 {
		// the summary line printed when only dropped notices were hidden
		var out []string
		for _, l := range lines {
			if l != "More deprecation notices were hidden, run again with `-v` to show them." {
				out = append(out, l)
			}
		}
		lines = out
	}
	return strings.Join(lines, "\n")
}

// reRemovalBar is a determinate bar (start and finish frames) drawn
// directly after the "  - Removing ..." lines.
var reRemovalBar = regexp.MustCompile(`((?:^|\n)  - Removing [^\n]*\n)( *0/(\d+) \[>-*\] +0%\n *\d+/\d+ \[=+\] 100%\n)`)

func indeterminateRemovalBar(s string) string {
	return reRemovalBar.ReplaceAllString(s, "${1}    0 [>---------------------------]    0 [->--------------------------]\n")
}
