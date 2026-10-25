#!/usr/bin/env bash
# Runs every step of .github/workflows/ci.yml locally and records the result
# under docs/validation/<UTC>-ci-local/. GitHub Actions itself is not run
# (that needs a push, which is out of scope).
#
#   unit (macos)   : host macOS, Go toolchain pinned to the workflow's 1.23
#                    via GOTOOLCHAIN (downloaded into .cache/gomod)
#   unit (linux)   : podman, docker.io/library/golang:1.23 (Debian based;
#                    stands in for ubuntu-24.04 + setup-go 1.23)
#   oracle (linux) : the workflow's compare matrix (PHP 8.4/8.3/8.5) and
#                    integration jobs run in sequence on the host; podman
#                    executes the Linux oracle containers inside the podman
#                    machine VM
set -uo pipefail
cd "$(dirname "$0")/.."
root=$PWD
stamp=$(date -u +%Y%m%dT%H%M%SZ)
out="$root/docs/validation/$stamp-ci-local"
mkdir -p "$out" "$root/.cache/gomod" "$root/.cache/go-build-ci"
gotc=${CI_GOTOOLCHAIN:-go1.23.12}
linux_image=${CI_LINUX_IMAGE:-docker.io/library/golang:1.23}
summary="$out/summary.md"
status=0

unit_steps='set -e
echo "== go version"; go version
echo "== gofmt"; test -z "$(gofmt -l cmd internal tools)"
echo "== vet"; go vet ./...
echo "== test"; go test ./...
echo "== race"; go test -race ./...
echo "== build"; go build -o bin/gocomposer ./cmd/gocomposer'

record() { # name exit log
  local r=PASS
  [[ $2 -ne 0 ]] && r="FAIL (exit $2)" && status=1
  echo "| $1 | $r | $(basename "$3") |" >>"$summary"
  echo "$1: $r"
}

podman image inspect "$linux_image" >/dev/null 2>&1 || podman pull "$linux_image" >"$out/pull.log" 2>&1
digest=$(podman image inspect --format '{{.Digest}}' "$linux_image" 2>/dev/null)

{
  echo "# CI local run $stamp"
  echo
  echo "Workflow: .github/workflows/ci.yml (sha256 $(shasum -a 256 .github/workflows/ci.yml | cut -d' ' -f1))"
  echo "Host: $(uname -srm)"
  echo "Linux image: $linux_image ($digest)"
  echo "GitHub Actions: not run (requires push; out of scope)"
  echo
  echo "| job | result | log |"
  echo "|---|---|---|"
} >"$summary"

# unit (macos)
log="$out/unit-macos.log"
env GOTOOLCHAIN="$gotc" GOMODCACHE="$root/.cache/gomod" GOCACHE="$root/.cache/go-build-ci/macos" \
  bash -c "$unit_steps" >"$log" 2>&1
record "unit (macos, $gotc)" $? "$log"

# unit (linux): run on a copy so bin/ and caches stay separate
log="$out/unit-linux.log"
work=$(mktemp -d "$root/.cache/ci-linux.XXXXXX")
tar --exclude=./.cache --exclude=./bin --exclude=./docs/validation -cf - . | tar -xf - -C "$work"
podman run --rm -v "$work:/src" -w /src -e GOFLAGS=-buildvcs=false "$linux_image" bash -c "$unit_steps" >"$log" 2>&1
record "unit (linux, $linux_image)" $? "$log"
rm -rf "$work"

# oracle (linux): the workflow's per-PHP matrix jobs and the integration job,
# each preparing only the images it needs (ORACLE_VERSIONS), run in sequence
for step in "prepare-8.4|ORACLE_VERSIONS=8.4 scripts/prepare-oracle.sh" \
  "compare-8.4|scripts/compare.sh -image localhost/gocomposer-oracle:8.4" \
  "prepare-8.3|ORACLE_VERSIONS=8.3 scripts/prepare-oracle.sh" \
  "compare-8.3|scripts/compare.sh -image localhost/gocomposer-oracle:8.3" \
  "prepare-8.5|ORACLE_VERSIONS=8.5 scripts/prepare-oracle.sh" \
  "compare-8.5|scripts/compare.sh -image localhost/gocomposer-oracle:8.5" \
  "integration-prepare|ORACLE_VERSIONS=8.4 scripts/prepare-oracle.sh" \
  "integration|scripts/integration.sh"; do
  name=${step%%|*}
  cmd=${step#*|}
  log="$out/oracle-$name.log"
  bash -c "$cmd" >"$log" 2>&1
  record "oracle: $name" $? "$log"
done

echo "results: $out"
exit $status
