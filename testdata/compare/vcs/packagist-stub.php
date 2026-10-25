<?php
// Loopback stub of packagist.org for the cli3-packagist-* comparison
// scenarios (tools/compare/scenarios_final_cli.go): commands run without a
// composer.json (search, browse) or with the default repository enabled
// (fund) use Composer's default repository https://repo.packagist.org.
// packagiststub.sh maps repo.packagist.org and packagist.org to 127.0.0.1
// (podman --add-host through _api/hosts) with a throw-away CA that both
// tools trust only through SSL_CERT_FILE.
//
//   https://repo.packagist.org/packages.json        root file (metadata-url, search, list)
//   https://repo.packagist.org/p2/<name>.json        tagged versions
//   https://repo.packagist.org/p2/<name>~dev.json    dev versions
//   https://packagist.org/search.json?q=&type=       search API
//   https://packagist.org/packages/list.json         package names (?vendor=&filter=)
//
//   php -S 127.0.0.1:8900 packagist-stub.php     (unused, 404)
//   php packagist-stub.php tls 443               HTTPS

const PACKAGES = [
    'stubvendor/alpha' => [
        'description' => 'Alpha library for the packagist stub',
        'versions' => [
            ['version' => '1.0.0', 'homepage' => 'https://alpha.example.org', 'support' => ['source' => 'https://github.com/stubvendor/alpha/tree/1.0.0', 'issues' => 'https://github.com/stubvendor/alpha/issues'], 'funding' => [['type' => 'github', 'url' => 'https://github.com/sponsors/alpha-old']]],
            ['version' => 'dev-main', 'default-branch' => true, 'homepage' => 'https://alpha.example.org', 'support' => ['source' => 'https://github.com/stubvendor/alpha/tree/main'], 'funding' => [['type' => 'github', 'url' => 'https://github.com/sponsors/alpha'], ['type' => 'custom', 'url' => 'https://alpha.example.org/donate']]],
        ],
    ],
    'stubvendor/beta' => [
        'description' => 'Beta tools',
        'versions' => [
            ['version' => '2.1.0', 'homepage' => 'https://beta.example.org', 'source' => ['type' => 'git', 'url' => 'https://github.com/stubvendor/beta.git', 'reference' => 'abcdef0123456789abcdef0123456789abcdef01']],
        ],
    ],
    'othervendor/gamma' => [
        'description' => 'Gamma package (abandoned)',
        'abandoned' => 'stubvendor/alpha',
        'versions' => [
            ['version' => '0.3.0', 'funding' => [['type' => 'patreon', 'url' => 'https://www.patreon.com/gamma']]],
        ],
    ],
];

/** @return array{0: int, 1: string, 2: string} */
function stub_handle(string $host, string $uri): array
{
    $json = static function ($data): array {
        return [200, 'application/json', json_encode($data, JSON_UNESCAPED_SLASHES)];
    };
    $path = (string) parse_url($uri, PHP_URL_PATH);
    parse_str((string) parse_url($uri, PHP_URL_QUERY), $q);
    if ($host === 'repo.packagist.org' && $path === '/packages.json') {
        return $json([
            'packages' => [],
            'metadata-url' => '/p2/%package%.json',
            'search' => 'https://packagist.org/search.json?q=%query%&type=%type%',
            'list' => 'https://packagist.org/packages/list.json',
        ]);
    }
    if ($host === 'repo.packagist.org' && preg_match('{^/p2/([a-z0-9_.-]+/[a-z0-9_.-]+?)(~dev)?\.json$}', $path, $m)) {
        if (!isset(PACKAGES[$m[1]])) {
            return [404, 'application/json', '{"error":"not found"}'];
        }
        $dev = isset($m[2]) && $m[2] === '~dev';
        $list = [];
        foreach (PACKAGES[$m[1]]['versions'] as $v) {
            if (str_starts_with($v['version'], 'dev-') !== $dev) {
                continue;
            }
            $entry = ['name' => $m[1], 'description' => PACKAGES[$m[1]]['description']] + $v;
            if (isset(PACKAGES[$m[1]]['abandoned'])) {
                $entry['abandoned'] = PACKAGES[$m[1]]['abandoned'];
            }
            $list[] = $entry;
        }

        return $json(['packages' => [$m[1] => $list]]);
    }
    if ($host === 'packagist.org' && $path === '/search.json') {
        $query = strtolower((string) ($q['q'] ?? ''));
        $results = [];
        foreach (PACKAGES as $name => $p) {
            if ($query !== '' && !str_contains($name, $query) && !str_contains(strtolower($p['description']), $query)) {
                continue;
            }
            $r = ['name' => $name, 'description' => $p['description'], 'url' => 'https://packagist.org/packages/'.$name, 'repository' => 'https://github.com/'.$name, 'downloads' => 10, 'favers' => 1];
            if (isset($p['abandoned'])) {
                $r['abandoned'] = $p['abandoned'];
            }
            $results[] = $r;
        }

        return $json(['results' => $results, 'total' => count($results)]);
    }
    if ($host === 'packagist.org' && $path === '/packages/list.json') {
        $names = array_keys(PACKAGES);
        if (isset($q['vendor'])) {
            $names = array_values(array_filter($names, static function ($n) use ($q) {
                return str_starts_with($n, $q['vendor'].'/');
            }));
        }
        if (isset($q['filter'])) {
            $re = '{^'.str_replace('\\*', '.*', preg_quote($q['filter'])).'$}i';
            $names = array_values(array_filter($names, static function ($n) use ($re) {
                return (bool) preg_match($re, $n);
            }));
        }

        return $json(['packageNames' => $names]);
    }

    return [404, 'application/json', '{"error":"not found"}'];
}

if (PHP_SAPI === 'cli-server') {
    http_response_code(404);

    return true;
}

if (($argv[1] ?? '') !== 'tls') {
    fwrite(STDERR, "usage: php packagist-stub.php tls PORT\n");
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
    $host = '';
    while (($l = fgets($conn)) !== false && trim($l) !== '') {
        [$k, $v] = array_pad(explode(':', $l, 2), 2, '');
        if (strtolower(trim($k)) === 'host') {
            $host = strtolower(preg_replace('{:\d+$}', '', trim($v)));
        }
    }
    [$status, $type, $body] = stub_handle($host, $rm[2]);
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
