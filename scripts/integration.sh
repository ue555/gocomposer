#!/usr/bin/env bash
# Installs the fixed real projects (testdata/integration/projects.json:
# symfony/console as a library root, a CLI app, laravel/laravel, and a tiny
# project downloading live) with the official Composer and with gocomposer
# in independent copies, then compares exit codes, output, composer.lock,
# the vendor tree, the autoloader and project commands (artisan about/list)
# run with the official PHP. Requires podman, the oracle image and network
# access (Packagist, GitHub). Results go to
# docs/validation/<UTC timestamp>-integration/. Exits non-zero on any
# difference or when the comparison cannot run (recorded as 未比較).
set -euo pipefail
cd "$(dirname "$0")/.."
export GOCACHE="${GOCACHE:-$PWD/.cache/go-build}"

if ! command -v podman >/dev/null 2>&1; then
  echo "未比較: podman が見つかりません" >&2
  exit 2
fi
if ! podman image exists localhost/gocomposer-oracle:8.4; then
  echo "未比較: oracle image localhost/gocomposer-oracle:8.4 がありません" >&2
  exit 2
fi
if ! curl -fsS --max-time 20 -o /dev/null https://repo.packagist.org/packages.json; then
  echo "未比較: repo.packagist.org に接続できません (ネットワークが必要)" >&2
  exit 2
fi

exec go run ./tools/integration "$@"
