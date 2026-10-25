# Loopback packagist.org (packagist-stub.php) for the cli3-packagist-*
# scenarios: a throw-away CA and a server certificate for repo.packagist.org
# and packagist.org, which resolve to 127.0.0.1 inside the comparison
# container (podman --add-host from _api/hosts). Both tools trust the CA only
# through the scenario's SSL_CERT_FILE.
api=/fixtures/git/_api
mkdir -p "$api"
cp /vcsdefs/packagist-stub.php "$api/forge-stub.php"
cd "$api"
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -subj "/CN=gocomposer packagist fixture CA" \
	-addext "basicConstraints=critical,CA:TRUE" -addext "keyUsage=critical,keyCertSign,cRLSign" \
	-keyout ca.key -out ca.pem 2>/dev/null
openssl req -newkey rsa:2048 -nodes -subj "/CN=packagist.test" -keyout server.key -out server.csr 2>/dev/null
printf '%s\n' "subjectAltName=DNS:repo.packagist.org,DNS:packagist.org" \
	"basicConstraints=CA:FALSE" "extendedKeyUsage=serverAuth" >server.ext
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -extfile server.ext -out server.pem 2>/dev/null
printf '%s\n' repo.packagist.org packagist.org >hosts
chmod a+r "$api"/*
