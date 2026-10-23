package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// gitStateScript records the observable state of every git checkout below
// the project after the compared command (vendor packages installed from
// source). The .git internals themselves (index timestamps, reflogs,
// FETCH_HEAD, packs) are not deterministic and are dropped from the compared
// tree by comparableTree; the raw tree stays in the result JSON.
const gitStateScript = `for g in $(find . -name .git -type d 2>/dev/null | sort); do d=$(dirname "$g"); echo "== $d"; (cd "$d" && echo "HEAD $(git rev-parse HEAD 2>&1)"; echo "ref $(git symbolic-ref -q HEAD)"; git status --porcelain --untracked-files=all; git remote -v; git config --get-all remote.origin.pushurl; git branch -a --no-color; git stash list) 2>&1; done >/work/out/gitstate
`

// comparableTree drops entries inside .git directories (the directory entry
// itself is kept so a missing/extra checkout is still reported).
func comparableTree(tree []string) []string {
	var out []string
	for _, l := range tree {
		if strings.Contains(l, "/.git/") || strings.HasPrefix(l, "F .git/") || strings.HasPrefix(l, "D .git/") {
			continue
		}
		out = append(out, l)
	}
	return out
}

func isDirPath(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

// runGitFixtureScript runs script (with lib.sh loaded) in the oracle
// container against <base>/gitfx mounted read-write at /fixtures/git.
func runGitFixtureScript(root, base, script string) error {
	gitfx := filepath.Join(base, "gitfx")
	if err := os.MkdirAll(gitfx, 0o777); err != nil {
		return err
	}
	full := ". /vcsdefs/lib.sh\n" + script
	args := podmanRun("--network=none",
		"-v", filepath.Join(root, "testdata", "compare", "vcs")+":/vcsdefs:ro",
		"-v", gitfx+":/fixtures/git",
		oracleImage, "sh", "-c", full)
	cmd := exec.Command("podman", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git fixture script failed: %v\n%s\n--- script:\n%s", err, out.String(), full)
	}
	return nil
}

// apiStubScript starts the loopback forge API stub when the scenario's git
// fixtures generated one (testdata/compare/vcs/forge.sh): plain HTTP on
// 127.0.0.1:8900 (custom gitlab-domains) and HTTPS on 127.0.0.1:443 for the
// forge host names mapped by forgePodmanArgs.
const apiStubScript = `if [ -f /fixtures/git/_api/forge-stub.php ]; then php -S 127.0.0.1:8900 /fixtures/git/_api/forge-stub.php >/dev/null 2>&1 &
php /fixtures/git/_api/forge-stub.php tls 443 >/dev/null 2>&1 &
i=0; while ! php -r 'exit(@fsockopen("127.0.0.1", 8900) && @fsockopen("127.0.0.1", 443) ? 0 : 1);' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done; fi`

// forgePodmanArgs resolve the public forge host names to the loopback stub
// (the container has no network) and let the unprivileged container user
// bind port 443. Only used when the forge fixture exists.
var forgePodmanArgs = []string{
	"--sysctl", "net.ipv4.ip_unprivileged_port_start=0",
	"--add-host=github.com:127.0.0.1", "--add-host=api.github.com:127.0.0.1",
	"--add-host=gitlab.com:127.0.0.1", "--add-host=bitbucket.org:127.0.0.1",
	"--add-host=api.bitbucket.org:127.0.0.1", "--add-host=ghe.test:127.0.0.1",
	"--add-host=gitlab.test:127.0.0.1",
}
