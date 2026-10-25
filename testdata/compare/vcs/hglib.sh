# Mercurial fixtures of the vcs-hg-* / source-hg-* scenarios
# (tools/compare/scenarios_vcs3.go), built with the oracle image's hg
# client. Commits and tags have a fixed user and date so the changeset ids
# are reproducible.
#
#   /fixtures/git/hg/hglib  default branch: tags 1.0.0 and 1.1.0, a later
#                           commit; named branch 2.x; bookmark "feature"
#   /fixtures/git/hg/hgwc, hgwc-topic  working copies (default, branch topic)
# _serve/hg.sh starts "hg serve" for it on http://127.0.0.1:8901/ for the
# compared runs.
export HOME=/tmp/fixture-home TZ=UTC LANG=C LC_ALL=C HGPLAIN=1 HGUSER=fixture
mkdir -p "$HOME"
HGROOT=/fixtures/git/hg
mkdir -p "$HGROOT"
n=0
# next sets D to the next fixed date (not in a subshell: n must advance)
next() { n=$((n + 1)); D="$((1704164645 + n * 3600)) 0"; }

hg init "$HGROOT/hglib"
cd "$HGROOT/hglib"
printf '%s\n' '{
    "name": "vcs3/hglib",
    "description": "hg fixture library",
    "license": "MIT",
    "autoload": {"psr-4": {"Vcs3\\HgLib\\": "src/"}}
}' >composer.json
mkdir -p src
printf '<?php\nnamespace Vcs3\\HgLib;\n\nclass Lib { const VERSION = "1.0.0"; }\n' >src/Lib.php
hg add -q
next
hg commit -q -m "Initial import" -d "$D"
next
hg tag -q -d "$D" 1.0.0
printf '<?php\nnamespace Vcs3\\HgLib;\n\nclass Lib { const VERSION = "1.1.0"; }\n' >src/Lib.php
next
hg commit -q -m "Release 1.1.0" -d "$D"
next
hg tag -q -d "$D" 1.1.0
hg branch -q 2.x
printf '<?php\nnamespace Vcs3\\HgLib;\n\nclass Lib { const VERSION = "2.x-dev"; }\n' >src/Lib.php
next
hg commit -q -m "Work on 2.x" -d "$D"
hg update -q default
hg bookmark -q -i feature
printf '<?php\nnamespace Vcs3\\HgLib;\n\nclass Lib { const VERSION = "dev-default"; }\n' >src/Lib.php
next
hg commit -q -m "Work on default" -d "$D"
hg update -q null
# working copies for the VersionGuesser scenarios (path repositories):
# default branch, and a new named branch "topic" (a feature branch)
hg clone -q "$HGROOT/hglib" "$HGROOT/hgwc"
hg clone -q "$HGROOT/hglib" "$HGROOT/hgwc-topic"
cd "$HGROOT/hgwc-topic"
hg branch -q topic
cd /

mkdir -p /fixtures/git/_serve
cat >/fixtures/git/_serve/hg.sh <<'SERVE'
HOME=/tmp/hg-serve-home hg serve -R /fixtures/git/hg/hglib -a 127.0.0.1 -p 8901 -d >/dev/null 2>&1
i=0; while ! php -r 'exit(@fsockopen("127.0.0.1", 8901) ? 0 : 1);' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done
SERVE
