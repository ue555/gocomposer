<?php
// Loopback Symfony Flex recipe endpoint for the flex-* comparison scenarios
// (tools/compare/scenarios_plugin2.go). It serves the static endpoint files
// copied by flexstub.sh to /fixtures/git/_api/flex/ (index.json, the recipe
// manifests) under https://127.0.0.1/flex/... with the throw-away fixture CA
// of forge.sh. Both Composer and gocomposer trust it only through the
// scenario's "cafile" config. Same command line as forge-stub.php:
//
//   php -S 127.0.0.1:8900 flex-stub.php      plain HTTP
//   php flex-stub.php tls 443                 HTTPS

const FLEX_ROOT = '/fixtures/git/_api/flex';

/** @return array{0: int, 1: string, 2: string} */
function flex_handle(string $uri): array
{
    $path = (string) parse_url($uri, PHP_URL_PATH);
    if (preg_match('{^/flex/((?:archived/[A-Za-z0-9._-]+/)?[A-Za-z0-9._-]+\.json)$}', $path, $m) && is_file(FLEX_ROOT.'/'.$m[1])) {
        return [200, 'application/json', (string) file_get_contents(FLEX_ROOT.'/'.$m[1])];
    }

    return [404, 'application/json', '{"error":"not found"}'];
}

if (PHP_SAPI === 'cli-server') {
    [$status, $type, $body] = flex_handle($_SERVER['REQUEST_URI']);
    http_response_code($status);
    header('Content-Type: '.$type);
    echo $body;

    return true;
}

if (($argv[1] ?? '') !== 'tls') {
    fwrite(STDERR, "usage: php flex-stub.php tls PORT\n");
    exit(1);
}
$ctx = stream_context_create(['ssl' => [
    'local_cert' => '/fixtures/git/_api/server.pem',
    'local_pk' => '/fixtures/git/_api/server.key',
    'verify_peer' => false,
]]);
$server = stream_socket_server('tls://127.0.0.1:'.(int) $argv[2], $errno, $errstr, STREAM_SERVER_BIND | STREAM_SERVER_LISTEN, $ctx);
if (!$server) {
    fwrite(STDERR, "listen failed: $errstr\n");
    exit(1);
}
$reasons = [200 => 'OK', 404 => 'Not Found'];
while (true) {
    $conn = @stream_socket_accept($server, -1);
    if (!$conn) {
        continue;
    }
    stream_set_timeout($conn, 10);
    $line = fgets($conn);
    if ($line === false || !preg_match('{^(GET|HEAD) (\S+) HTTP/1\.[01]}', $line, $rm)) {
        fclose($conn);
        continue;
    }
    while (($l = fgets($conn)) !== false && trim($l) !== '') {
    }
    [$status, $type, $body] = flex_handle($rm[2]);
    $data = "HTTP/1.1 $status ".($reasons[$status] ?? 'Error')."\r\nContent-Type: $type\r\nContent-Length: ".strlen($body)."\r\nConnection: close\r\n\r\n".($rm[1] === 'HEAD' ? '' : $body);
    while ($data !== '') {
        $n = @fwrite($conn, $data);
        if ($n === false || $n === 0) {
            break;
        }
        $data = (string) substr($data, $n);
    }
    @fclose($conn);
}
