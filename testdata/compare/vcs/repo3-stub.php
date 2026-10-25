<?php
// Loopback router for the warnings-* comparison scenarios that need error
// responses carrying repository warnings (tools/compare/scenarios_repo3.go),
// installed as /fixtures/git/_api/forge-stub.php by repo3stub.sh so the
// harness starts it on 127.0.0.1:8900:
//
//   /w403/packages.json      403 application/json with warning + info
//   /w404/packages.json      v2 repository (provfix/lib available)
//   /w404/p2/<name>.json     404 application/json with an info message
//
//   php -S 127.0.0.1:8900 forge-stub.php   the router
//   php forge-stub.php tls <port>          only holds the port open (the
//                                          harness waits for it; no TLS
//                                          endpoint is used by these scenarios)

if (PHP_SAPI === 'cli') {
    if (($argv[1] ?? '') === 'tls') {
        $server = @stream_socket_server('tcp://127.0.0.1:'.(int) ($argv[2] ?? 443));
        while ($server && ($c = @stream_socket_accept($server, -1))) {
            fclose($c);
        }
    }
    exit(0);
}

$path = rawurldecode((string) parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH));
$json = static function (int $status, array $data): bool {
    http_response_code($status);
    header('Content-Type: application/json');
    echo json_encode($data, JSON_UNESCAPED_SLASHES);

    return true;
};

if ($path === '/w403/packages.json') {
    return $json(403, ['warning' => 'Access to this repository is blocked by policy', 'info' => 'Ask your administrator for access']);
}
if ($path === '/w404/packages.json') {
    return $json(200, ['packages' => [], 'metadata-url' => '/w404/p2/%package%.json', 'available-packages' => ['provfix/lib']]);
}
if (str_starts_with($path, '/w404/p2/')) {
    return $json(404, ['info' => 'This package is not mirrored here']);
}
http_response_code(404);
header('Content-Type: text/plain');
echo "not found\n";
return true;
