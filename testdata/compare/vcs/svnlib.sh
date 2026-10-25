# Subversion fixtures of the vcs-svn-* / source-svn-* scenarios
# (tools/compare/scenarios_vcs3.go), built with the oracle image's svn
# client. Every revision gets a fixed svn:date and the repositories a fixed
# uuid so both tools read the same, reproducible history.
#
#   /fixtures/git/svn/svnlib   standard layout: trunk, tags 1.0.0 and 1.1.0,
#                              branch 2.x, a later trunk commit
#   /fixtures/git/svn/wc-trunk, wc-2x  working copies of svnlib trunk and 2.x
#   /fixtures/git/svn/svnmono  "package-path" layout: trunk/pkg, tags/1.0.0/pkg
#   /fixtures/git/svn/svnpriv  copy of svnlib (vcs3/svnpriv) served by
#                              svnserve with authentication required
# _serve/svn.sh starts svnserve on svn://127.0.0.1:3690 for the compared runs.
export HOME=/tmp/fixture-home USER=fixture TZ=UTC LANG=C LC_ALL=C
mkdir -p "$HOME"
SVNROOT=/fixtures/git/svn
mkdir -p "$SVNROOT"

svn_dates() {
	repo=$1
	n=$(svnlook youngest "$repo")
	i=1
	while [ "$i" -le "$n" ]; do
		ts=$((1704164645 + i * 3600))
		printf '%s' "$(date -u -d "@$ts" '+%Y-%m-%dT%H:%M:%S.000000Z')" >/tmp/svndate
		svnadmin setrevprop "$repo" -r "$i" svn:date /tmp/svndate
		printf 'fixture' >/tmp/svnauthor
		svnadmin setrevprop "$repo" -r "$i" svn:author /tmp/svnauthor
		i=$((i + 1))
	done
}

svn_lib() {
	name=$1 pkg=$2 ns=$3
	repo=$SVNROOT/$name
	svnadmin create "$repo"
	svnadmin setuuid "$repo" "$4"
	url=file://$repo
	svn mkdir -q -m "Layout" "$url/trunk" "$url/branches" "$url/tags"
	wc=/tmp/svnwc-$name
	svn co -q "$url/trunk" "$wc"
	cd "$wc"
	printf '{\n    "name": "%s",\n    "description": "svn fixture library",\n    "license": "MIT",\n    "autoload": {"psr-4": {"%s\\\\": "src/"}}\n}\n' "$pkg" "$5" >composer.json
	mkdir -p src
	printf '<?php\nnamespace %s;\n\nclass Lib { const VERSION = "1.0.0"; }\n' "$ns" >src/Lib.php
	svn add -q composer.json src
	svn ci -q -m "Initial import"
	svn cp -q -m "Tag 1.0.0" "$url/trunk" "$url/tags/1.0.0"
	printf '<?php\nnamespace %s;\n\nclass Lib { const VERSION = "1.1.0"; }\n' "$ns" >src/Lib.php
	svn ci -q -m "Release 1.1.0"
	svn cp -q -m "Tag 1.1.0" "$url/trunk" "$url/tags/1.1.0"
	svn cp -q -m "Branch 2.x" "$url/trunk" "$url/branches/2.x"
	svn up -q
	printf '<?php\nnamespace %s;\n\nclass Lib { const VERSION = "dev-trunk"; }\n' "$ns" >src/Lib.php
	svn ci -q -m "Work on trunk"
	cd /
	rm -rf "$wc"
	svn_dates "$repo"
}

svn_lib svnlib vcs3/svnlib 'Vcs3\SvnLib' 0f6a3c8e-0000-4000-8000-000000000001 'Vcs3\\SvnLib'
# working copies for the VersionGuesser scenarios (path repositories)
svn co -q "file://$SVNROOT/svnlib/trunk" "$SVNROOT/wc-trunk"
svn co -q "file://$SVNROOT/svnlib/branches/2.x" "$SVNROOT/wc-2x"
svn_lib svnpriv vcs3/svnpriv 'Vcs3\SvnPriv' 0f6a3c8e-0000-4000-8000-000000000002 'Vcs3\\SvnPriv'
cat >"$SVNROOT/svnpriv/conf/svnserve.conf" <<'CONF'
[general]
anon-access = none
auth-access = read
password-db = passwd
realm = gocomposer-fixture
CONF
printf '[users]\nfixture-user = s3cr3t-svn\n' >"$SVNROOT/svnpriv/conf/passwd"
chmod a+r "$SVNROOT/svnpriv/conf/passwd"

# package-path layout
repo=$SVNROOT/svnmono
svnadmin create "$repo"
svnadmin setuuid "$repo" 0f6a3c8e-0000-4000-8000-000000000003
url=file://$repo
svn mkdir -q -m "Layout" "$url/trunk" "$url/branches" "$url/tags"
wc=/tmp/svnwc-mono
svn co -q "$url/trunk" "$wc"
cd "$wc"
mkdir -p pkg/src other
printf '%s\n' '{
    "name": "vcs3/svnmono",
    "license": "MIT",
    "autoload": {"psr-4": {"Vcs3\\SvnMono\\": "src/"}}
}' >pkg/composer.json
printf '<?php\nnamespace Vcs3\\SvnMono;\n\nclass Lib {}\n' >pkg/src/Lib.php
printf 'not a package\n' >other/README
svn add -q pkg other
svn ci -q -m "Initial import"
svn cp -q -m "Tag 1.0.0" "$url/trunk" "$url/tags/1.0.0"
cd /
rm -rf "$wc"
svn_dates "$repo"

mkdir -p /fixtures/git/_serve
cat >/fixtures/git/_serve/svn.sh <<'SERVE'
svnserve -d -r /fixtures/git/svn --listen-host 127.0.0.1 --listen-port 3690 >/dev/null 2>&1
i=0; while ! php -r 'exit(@fsockopen("127.0.0.1", 3690) ? 0 : 1);' 2>/dev/null; do i=$((i+1)); [ $i -gt 50 ] && break; sleep 0.1; done
SERVE
