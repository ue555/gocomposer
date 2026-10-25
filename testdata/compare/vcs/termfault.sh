# Loopback fault servers for the curl-* scenarios (scenarios_term.go):
# copies term-fault-server.php next to throw-away certificates. ca.pem is the
# trusted CA given to both tools through the scenario's "cafile" config;
# the server certificates are self-signed, for another host, expired, or
# issued by a CA that is not trusted.
fx=/fixtures/git/_fault
mkdir -p "$fx"
cp /vcsdefs/term-fault-server.php "$fx/server.php"
cd "$fx"
ca() { # name subject
	openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -subj "/CN=$2" \
		-addext "basicConstraints=critical,CA:TRUE" -addext "keyUsage=critical,keyCertSign,cRLSign" \
		-keyout "$1.key" -out "$1.pem" 2>/dev/null
}
leaf() { # name ca san days-option...
	n=$1 c=$2 san=$3
	shift 3
	openssl req -newkey rsa:2048 -nodes -subj "/CN=fault.test" -keyout "$n.key" -out "$n.csr" 2>/dev/null
	printf '%s\n' "subjectAltName=$san" "basicConstraints=CA:FALSE" "extendedKeyUsage=serverAuth" >"$n.ext"
	openssl x509 -req -in "$n.csr" -CA "$c.pem" -CAkey "$c.key" -CAcreateserial -extfile "$n.ext" -out "$n.pem" "$@" 2>/dev/null
}
ca ca "gocomposer fault CA"
ca other-ca "untrusted fault CA"
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -subj "/CN=127.0.0.1" -addext "subjectAltName=IP:127.0.0.1" \
	-keyout self-signed.key -out self-signed.pem 2>/dev/null
leaf good ca "IP:127.0.0.1" -days 3650
leaf wrong-host ca "DNS:other.test" -days 3650
leaf unknown-ca other-ca "IP:127.0.0.1" -days 3650
# expired: valid from 2020-01-01 to 2020-01-02
leaf expired ca "IP:127.0.0.1" -not_before 20200101000000Z -not_after 20200102000000Z
chmod a+r "$fx"/*
# composer repositories (served by the fixture web server below
# /git/_fault/) using the /v2adv package metadata, with an advisory api-url
# on a fault server
api() { # dir api-url
	mkdir -p "$fx/$1"
	printf '{"packages": [], "metadata-url": "/v2adv/p2/%%package%%.json", "security-advisories": {"metadata": false, "api-url": "%s"}}\n' "$2" >"$fx/$1/packages.json"
	chmod a+r "$fx/$1/packages.json"
}
api hang-api http://127.0.0.1:9201/advisories
api hang-tls-api https://127.0.0.1:9208/advisories
api close-api http://127.0.0.1:9202/advisories
