# Scripted stand-in for the Perforce command line client used by the
# vcs-perforce-* scenarios (tools/compare/scenarios_vcs3.go). The real
# Helix Core client/server are proprietary (Perforce EULA) and not in the
# oracle image, so both tools run this same script as "p4" (PATH is set by
# the scenarios). It answers the commands Composer's Util/Perforce issues
# with fixed output, like the expectations of Composer's own PerforceTest:
#
#   server 127.0.0.1:1666, user "fixture" (logged in) or "locked" (needs the
#   password "p4-s3cr3t"); depot "depot" (local) with labels 1.0.0 and 1.1.0
#   and head change 5; depot "streams" (stream) with stream //streams/main.
#   "p4 print" prints the file content only (like p4 print -q).
mkdir -p /fixtures/git/_p4/bin
cat >/fixtures/git/_p4/bin/p4 <<'P4'
#!/bin/sh
user= client= port=
while [ $# -gt 0 ]; do
	case "$1" in
	-u) user=$2; shift 2 ;;
	-c) client=$2; shift 2 ;;
	-p) port=$2; shift 2 ;;
	*) break ;;
	esac
done
cmd=$1
[ $# -gt 0 ] && shift
if [ "$port" != "127.0.0.1:1666" ]; then
	printf 'Perforce client error:\n\tConnect to server failed; check $P4PORT.\n\tTCP connect to %s failed.\n' "$port" >&2
	exit 1
fi
composer_json() {
	printf '{\n    "name": "vcs3/p4lib",\n    "description": "perforce fixture library",\n    "license": "MIT",\n    "autoload": {"psr-4": {"Vcs3\\\\P4Lib\\\\": "src/"}}\n}\n'
}
version_of() {
	case "$1" in
	*@1.0.0 | *@2) echo 1.0.0 ;;
	*@1.1.0 | *@4) echo 1.1.0 ;;
	*) echo dev-master ;;
	esac
}
case "$cmd" in
info)
	printf 'User name: %s\nServer address: 127.0.0.1:1666\nServer version: P4D/STUB\n' "$user"
	;;
login)
	if [ "$1" = "-s" ]; then
		if [ "$user" = "locked" ]; then
			echo "Perforce password (P4PASSWD) invalid or unset." >&2
			exit 1
		fi
		echo "User $user ticket expires in 12 hours 0 minutes."
	else
		read -r pw
		if [ "$pw" != "p4-s3cr3t" ]; then
			echo "Password invalid." >&2
			exit 1
		fi
		echo "User $user logged in."
	fi
	;;
depots)
	printf "Depot depot 2024/01/02 local depot/... 'Default depot '\nDepot streams 2024/01/02 stream streams/... 'Stream depot '\n"
	;;
streams)
	printf "Stream //streams/main mainline none 'main' \n"
	;;
client)
	if [ "$1" = "-d" ]; then
		echo "Client $2 deleted."
	else
		cat >/dev/null
		echo "Client $client saved."
	fi
	;;
labels)
	printf "Label 1.0.0 2024/01/02 'Release 1.0.0 '\nLabel 1.1.0 2024/01/02 'Release 1.1.0 '\n"
	;;
changes)
	if [ "$1" = "-m1" ]; then
		case "$2" in
		@1.0.0) echo "Change 2 on 2024/01/02 by fixture@ws 'Release 1.0.0 '" ;;
		@1.1.0) echo "Change 4 on 2024/01/02 by fixture@ws 'Release 1.1.0 '" ;;
		*) echo "$2 - no such file(s)." >&2; exit 1 ;;
		esac
	else
		printf "Change 5 on 2024/01/02 by fixture@ws 'Work on main '\nChange 4 on 2024/01/02 by fixture@ws 'Release 1.1.0 '\nChange 2 on 2024/01/02 by fixture@ws 'Release 1.0.0 '\n"
	fi
	;;
files)
	case "$1" in
	*/composer.json@*) echo "${1%@*}#1 - add change 2 (text)" ;;
	*) echo "$1 - no such file(s)." ;;
	esac
	;;
print)
	case "$1" in
	*/composer.json*) composer_json ;;
	*) echo "$1 - no such file(s)." >&2 ;;
	esac
	;;
sync)
	v=$(version_of "$2")
	mkdir -p src
	composer_json >composer.json
	printf '<?php\nnamespace Vcs3\\P4Lib;\n\nclass Lib { const VERSION = "%s"; }\n' "$v" >src/Lib.php
	echo "//depot/composer.json#1 - added as $PWD/composer.json"
	echo "//depot/src/Lib.php#1 - added as $PWD/src/Lib.php"
	;;
filelog)
	printf '//depot/src/Lib.php\n... #2 change 4 edit on 2024/01/02 by fixture@ws (text) '"'"'Release 1.1.0 '"'"'\n'
	;;
*)
	echo "Unknown command.  Try 'p4 help' for info." >&2
	exit 1
	;;
esac
P4
chmod 755 /fixtures/git/_p4/bin/p4
