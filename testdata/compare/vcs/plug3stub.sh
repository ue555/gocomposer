# Loopback stub for the plugin3-* scenarios (plug3-stub.php): the throw-away
# CA and server certificate of forge.sh (valid for 127.0.0.1), with
# plug3-stub.php in place of the forge API stub. It serves the fixture dists
# with a fixed delay per URL (/delay/<ms>/dists/<file>, so that parallel
# downloads complete in a known order) on plain HTTP 127.0.0.1:8900 and a
# small JSON document over HTTPS on 127.0.0.1:443 (and on 127.0.0.1:8444 to
# clients presenting a certificate of the fixture CA: client.pem with
# client.key, or client-enc.key encrypted with the passphrase "secret").
. /vcsdefs/forge.sh
openssl req -newkey rsa:2048 -nodes -subj "/CN=plugfix3 client" -keyout client.key -out client.csr 2>/dev/null
printf '%s\n' "basicConstraints=CA:FALSE" "extendedKeyUsage=clientAuth" >client.ext
openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -extfile client.ext -out client.pem 2>/dev/null
openssl rsa -in client.key -aes128 -traditional -passout pass:secret -out client-enc.key 2>/dev/null
cp /vcsdefs/plug3-stub.php "$api/forge-stub.php"
chmod -R a+rX "$api"
