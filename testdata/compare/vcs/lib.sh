# Helpers for the deterministic git fixture repositories used by the vcs-*
# and source-* comparison scenarios (tools/compare/scenarios_vcs.go).
#
# The scripts run inside the oracle container (POSIX sh, git 2.47) with
# /fixtures/git mounted read-write. Every commit gets a fixed author,
# committer and an increasing fixed timestamp so commit hashes are
# reproducible; both tools then read the very same bare repositories.
#
#   repo NAME          start (or continue) the work tree /fixtures/git/_work/NAME
#   put FILE CONTENT   write a file in the current work tree
#   commit MSG         commit everything with the next fixed timestamp
#   tag NAME [MSG]     lightweight tag (annotated when MSG is given)
#   branch NAME        create and switch to a branch (git checkout -b)
#   switch NAME        switch to an existing branch
#   publish NAME       mirror the work tree refs into /fixtures/git/NAME.git
#                      (bare, HEAD = main, dumb-http info refreshed)

set -e
export HOME=/tmp/fixture-home
mkdir -p "$HOME"
export GIT_AUTHOR_NAME="Fixture Author" GIT_AUTHOR_EMAIL="author@example.org"
export GIT_COMMITTER_NAME="Fixture Committer" GIT_COMMITTER_EMAIL="committer@example.org"
export GIT_CONFIG_NOSYSTEM=1 TZ=UTC LANG=C LC_ALL=C
FIXROOT=/fixtures/git
WORKROOT=$FIXROOT/_work
mkdir -p "$WORKROOT"

_next_date() {
	n=$(cat "$WORKROOT/.counter" 2>/dev/null || echo 0)
	n=$((n + 1))
	echo "$n" >"$WORKROOT/.counter"
	ts=$((1704164645 + n * 3600))
	export GIT_AUTHOR_DATE="@$ts +0000" GIT_COMMITTER_DATE="@$ts +0000"
}

repo() {
	cd "$WORKROOT"
	if [ ! -d "$1" ]; then
		git init -q -b main "$1"
		git -C "$1" config commit.gpgsign false
		git -C "$1" config tag.gpgsign false
	fi
	cd "$WORKROOT/$1"
}

put() {
	mkdir -p "$(dirname "$1")"
	printf '%s\n' "$2" >"$1"
}

commit() {
	_next_date
	git add -A
	git commit -q --allow-empty -m "$1"
}

tag() {
	_next_date
	if [ -n "$2" ]; then
		git tag -a -m "$2" "$1"
	else
		git tag "$1"
	fi
}

branch() { git checkout -q -b "$1"; }
switch() { git checkout -q "$1"; }

publish() {
	dest="$FIXROOT/$1.git"
	if [ ! -d "$dest" ]; then
		git init -q --bare -b main "$dest"
	fi
	git push -q --mirror --force "$dest"
	git -C "$dest" update-server-info
}
