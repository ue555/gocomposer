#!/usr/bin/env bash
# Differential comparison of gocomposer against the pinned official Composer
# (oracle image localhost/gocomposer-oracle:8.4 with Composer
# f83588392df43d4fd027d84d26de459946ef4e80 from .cache/oracle-source).
# Every scenario runs both tools on independent copies with --network=none
# and a loopback fixture repository. Results (raw logs included) are written
# to docs/validation/<UTC timestamp>-compare/. Exits non-zero when any
# scenario differs (documented KNOWN-DIFF scenarios excepted) or cannot run.
set -euo pipefail
cd "$(dirname "$0")/.."
export GOCACHE="${GOCACHE:-$PWD/.cache/go-build}"

if ! command -v podman >/dev/null 2>&1; then
  echo "未比較: podman が見つかりません (oracle image が必要)" >&2
  exit 2
fi
if ! podman image exists localhost/gocomposer-oracle:8.4; then
  echo "未比較: oracle image localhost/gocomposer-oracle:8.4 がありません (oracle/Containerfile から作成)" >&2
  exit 2
fi
if [[ ! -f .cache/oracle-source/bin/composer || ! -d .cache/oracle-source/vendor ]]; then
  echo "未比較: .cache/oracle-source (Composer f83588392df4 + vendor) がありません" >&2
  exit 2
fi
head=$(git -C .cache/oracle-source rev-parse HEAD)
want=$(python3 -c 'import json;print(json.load(open("docs/upstream.json"))["commit"])')
if [[ "$head" != "$want" ]]; then
  echo "未比較: .cache/oracle-source は $head (期待値 $want)" >&2
  exit 2
fi

exec go run ./tools/compare "$@"
