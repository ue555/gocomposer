# Loopback Symfony Flex recipe endpoint (flex-stub.php) for the flex-*
# scenarios: the throw-away CA and server certificate of forge.sh (valid for
# 127.0.0.1), with the flex stub in place of the forge API stub and the
# endpoint files of testdata/compare/vcs/flex/ (index.json and the recipes).
. /vcsdefs/forge.sh
cp /vcsdefs/flex-stub.php "$api/forge-stub.php"
mkdir -p "$api/flex"
cp -R /vcsdefs/flex/. "$api/flex/"
chmod -R a+rX "$api"
