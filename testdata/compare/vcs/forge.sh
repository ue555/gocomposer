# Loopback forge API stub (forge-stub.php) for the vcs-forge-* scenarios:
# copies the stub next to the bare repositories and creates a throw-away CA
# plus a server certificate for the forge host names (they resolve to
# 127.0.0.1 inside the comparison container through podman --add-host).
# Both tools trust the CA only through the scenario's "cafile" config.
api=/fixtures/git/_api
mkdir -p "$api"
cp /vcsdefs/forge-stub.php "$api/forge-stub.php"
cd "$api"
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -subj "/CN=gocomposer fixture CA" \
	-addext "basicConstraints=critical,CA:TRUE" -addext "keyUsage=critical,keyCertSign,cRLSign" \
	-keyout ca.key -out ca.pem 2>/dev/null
openssl req -newkey rsa:2048 -nodes -subj "/CN=forge.test" -keyout server.key -out server.csr 2>/dev/null
printf '%s\n' "subjectAltName=DNS:github.com,DNS:api.github.com,DNS:gitlab.com,DNS:bitbucket.org,DNS:api.bitbucket.org,DNS:ghe.test,DNS:gitlab.test,IP:127.0.0.1" \
	"basicConstraints=CA:FALSE" "extendedKeyUsage=serverAuth" >server.ext
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -extfile server.ext -out server.pem 2>/dev/null
chmod a+r "$api"/*
