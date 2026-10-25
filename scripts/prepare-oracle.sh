#!/usr/bin/env bash
# Prepares the official-Composer comparison environment from scratch:
#  - oracle images localhost/gocomposer-oracle:{8.3,8.4,8.5} (and their
#    -zip variants from oracle/Containerfile.zip and -arch variants from
#    oracle/Containerfile.arch) built from
#    oracle/Containerfile on the official php:<v>-cli images (pinned digest
#    for 8.4, recorded for all in .cache/oracle-images.txt)
#  - .cache/oracle-source: an independent clone of Composer at the commit in
#    docs/upstream.json (never the read-only reference checkout itself) with
#    its locked dependencies installed by a checksum-verified composer.phar.
# Usage: scripts/prepare-oracle.sh [composer git source]
# ORACLE_VERSIONS limits the images built (default "8.3 8.4 8.5"); 8.4 is
# always built because it installs .cache/oracle-source (CI matrix jobs).
set -euo pipefail
cd "$(dirname "$0")/.."

commit=$(python3 -c 'import json;print(json.load(open("docs/upstream.json"))["commit"])')
lock_sha=$(python3 -c 'import json;print(json.load(open("docs/upstream.json"))["composer_lock_sha256"])')
src="${1:-https://github.com/composer/composer.git}"
# multi-arch index of php:8.4.27-cli (arm64 entry sha256:05c3a54f5705…, amd64 entry sha256:a007a658…)
php84_digest="sha256:a4ae1f57e22b89a0a69c5e085b3010ee727d094e1542d249e48a01bec3c44810"
bootstrap_version="2.10.3"
bootstrap_sha256="7a2d379d5b8ffdaa028580ef26494c36d2feef4b178d3dd1473a4dbc5e17c8d6"

mkdir -p .cache
touch .cache/oracle-images.txt
versions=" ${ORACLE_VERSIONS:-8.3 8.4 8.5} "
[[ "$versions" == *" 8.4 "* ]] || versions="${versions}8.4 "
for v in 8.3 8.4 8.5; do
  [[ "$versions" == *" $v "* ]] || continue
  if [[ "$v" == 8.4 ]]; then
    base="docker.io/library/php@${php84_digest}"
    podman pull -q "$base" >/dev/null
  else
    podman pull -q "docker.io/library/php:${v}-cli" >/dev/null
    base="docker.io/library/php@$(podman image inspect "docker.io/library/php:${v}-cli" --format '{{.Digest}}')"
  fi
  podman build -q --build-arg "PHP_IMAGE=${base}" -t "localhost/gocomposer-oracle:${v}" -f oracle/Containerfile oracle >/dev/null
  # ext-zip variant for the scenarios whose Composer code path requires it
  # (ArtifactRepository); the default image lacks ext-zip on purpose.
  podman build -q --build-arg "BASE=localhost/gocomposer-oracle:${v}" -t "localhost/gocomposer-oracle:${v}-zip" -f oracle/Containerfile.zip oracle >/dev/null
  # archive-tools variant (7z, unrar, no unzip) for the archive3-* scenarios
  podman build -q --build-arg "BASE=localhost/gocomposer-oracle:${v}-zip" -t "localhost/gocomposer-oracle:${v}-arch" -f oracle/Containerfile.arch oracle >/dev/null
  { grep -v "^${v} " .cache/oracle-images.txt || true; echo "${v} ${base}"; } | sort > .cache/oracle-images.txt.new
  mv .cache/oracle-images.txt.new .cache/oracle-images.txt
done

if [[ ! -d .cache/oracle-source/.git ]]; then
  git clone --quiet "$src" .cache/oracle-source
fi
git -C .cache/oracle-source fetch --quiet origin "$commit" 2>/dev/null || true
git -C .cache/oracle-source checkout --quiet --detach "$commit"
actual_lock=$(shasum -a 256 .cache/oracle-source/composer.lock | cut -d' ' -f1)
if [[ "$actual_lock" != "$lock_sha" ]]; then
  echo "composer.lock sha256 $actual_lock does not match docs/upstream.json ($lock_sha)" >&2
  exit 1
fi

phar=.cache/composer-bootstrap-${bootstrap_version}.phar
if [[ ! -f "$phar" ]]; then
  curl -fsSL -o "$phar.tmp" "https://getcomposer.org/download/${bootstrap_version}/composer.phar"
  mv "$phar.tmp" "$phar"
fi
actual=$(shasum -a 256 "$phar" | cut -d' ' -f1)
if [[ "$actual" != "$bootstrap_sha256" ]]; then
  echo "bootstrap composer.phar sha256 $actual does not match $bootstrap_sha256" >&2
  exit 1
fi
cp "$phar" .cache/oracle-source/composer-bootstrap.phar
# Rootless podman on Linux (GitHub Actions): map the host user to the image's
# "oracle" user (uid/gid 1000) so the bind mount is writable and owned by it.
# podman machine on macOS needs no mapping (same as tools/*/userns.go).
userns=()
if [[ "$(uname -s)" == Linux && "$(id -u)" != 0 ]]; then
  userns=(--userns=keep-id:uid=1000,gid=1000)
fi
podman run --rm ${userns[@]+"${userns[@]}"} -v "$PWD/.cache/oracle-source:/src" -w /src -e COMPOSER_HOME=/tmp/home -e COMPOSER_CACHE_DIR=/tmp/cache \
  localhost/gocomposer-oracle:8.4 php composer-bootstrap.phar install --no-dev --no-interaction --no-progress
rm -f .cache/oracle-source/composer-bootstrap.phar
podman run --rm ${userns[@]+"${userns[@]}"} -v "$PWD/.cache/oracle-source:/oracle-source:ro" localhost/gocomposer-oracle:8.4 php /oracle-source/bin/composer --version
