# Fossil fixtures of the vcs-fossil-* / source-fossil-* scenarios
# (tools/compare/scenarios_vcs3.go), built with the oracle image's fossil.
# Check-ins and tags use fixed dates and user; FossilDriver identifies
# versions by tag and branch names, which are what the lock file records.
#
#   /fixtures/git/fossil/fossillib.fossil  trunk: tags 1.0.0 and 1.1.0, a
#                                          later check-in; branch 2.x
#   /fixtures/git/fossil/checkout/         an open checkout of it (local
#                                          directory url)
# _serve/fossil.sh serves the repository on http://127.0.0.1:8902/ for the
# compared runs.
export HOME=/tmp/fixture-home USER=fixture TZ=UTC LANG=C LC_ALL=C
mkdir -p "$HOME"
FROOT=/fixtures/git/fossil
mkdir -p "$FROOT"
n=0
# next sets D to the next fixed date (not in a subshell: n must advance)
next() { n=$((n + 1)); D=$(date -u -d "@$((1704164645 + n * 3600))" '+%Y-%m-%dT%H:%M:%S'); }
# fossil detects edits by mtime and size: give every edit its own mtime
fput() { printf "$2" >"$1"; touch -d "@$((1704164645 + n * 3600 + 1800))" "$1"; }

next
fossil init --date-override "$D" "$FROOT/fossillib.fossil" >/dev/null
mkdir -p /tmp/fossilwc
cd /tmp/fossilwc
fossil open "$FROOT/fossillib.fossil" >/dev/null
printf '%s\n' '{
    "name": "vcs3/fossillib",
    "description": "fossil fixture library",
    "license": "MIT",
    "autoload": {"psr-4": {"Vcs3\\FossilLib\\": "src/"}}
}' >composer.json
mkdir -p src
printf '<?php\nnamespace Vcs3\\FossilLib;\n\nclass Lib { const VERSION = "1.0.0"; }\n' >src/Lib.php
fossil add composer.json src >/dev/null
next
fossil commit -m "Initial import" --date-override "$D" --no-warnings >/dev/null
next
fossil tag add --date-override "$D" 1.0.0 trunk >/dev/null
fput src/Lib.php '<?php\nnamespace Vcs3\\FossilLib;\n\nclass Lib { const VERSION = "1.1.0"; }\n'
next
fossil commit -m "Release 1.1.0" --date-override "$D" --no-warnings >/dev/null
next
fossil tag add --date-override "$D" 1.1.0 trunk >/dev/null
next
fossil branch new --date-override "$D" 2.x trunk >/dev/null
fput src/Lib.php '<?php\nnamespace Vcs3\\FossilLib;\n\nclass Lib { const VERSION = "dev-trunk"; }\n'
next
fossil commit -m "Work on trunk" --date-override "$D" --no-warnings >/dev/null
fossil close >/dev/null
cd /
rm -rf /tmp/fossilwc
mkdir -p "$FROOT/checkout"
cd "$FROOT/checkout"
fossil open "$FROOT/fossillib.fossil" >/dev/null
cd /

mkdir -p /fixtures/git/_serve
cat >/fixtures/git/_serve/fossil.sh <<'SERVE'
(HOME=/tmp/fossil-serve-home; mkdir -p $HOME; export HOME; fossil server --port 8902 --localhost /fixtures/git/fossil/fossillib.fossil >/dev/null 2>&1 &)
i=0; while ! php -r 'exit(@fsockopen("127.0.0.1", 8902) ? 0 : 1);' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done
SERVE
