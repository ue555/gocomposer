#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ ! -f go.mod || ! -d cmd/gocomposer ]]; then
  echo '未実装: Go moduleとcmd/gocomposerが必要です' >&2
  exit 2
fi
mkdir -p .cache/go-build bin
export GOCACHE="$PWD/.cache/go-build"
unformatted=$(gofmt -l cmd internal 2>/dev/null) || exit 2
if [[ -n "$unformatted" ]]; then echo "$unformatted"; exit 1; fi
go test ./...
go test -race ./...
go vet ./...
go build -o bin/gocomposer ./cmd/gocomposer
for gate in scripts/compare.sh scripts/integration.sh; do
  if [[ ! -x "$gate" ]]; then echo "未実装の検証: $gate" >&2; exit 2; fi
  "$gate"
done
if git rev-parse --show-toplevel >/dev/null 2>&1; then git diff --check; fi
