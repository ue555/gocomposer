<?php
// Loopback stub of the plugin3-* comparison scenarios
// (tools/compare/scenarios_plugin3.go), started like forge-stub.php:
//
//   php -S 127.0.0.1:8900 plug3-stub.php     plain HTTP
//   php plug3-stub.php tls 443               HTTPS (throw-away CA of forge.sh);
//                                            first starts the delay server:
//   php plug3-stub.php delay 8901            plain HTTP, one process serving
//                                            all connections concurrently
//
// /delay/<ms>/<path> answers after <ms> milliseconds with /fixtures/<path>
// (the fixture dists): parallel downloads complete in the order of their
// delays (php -S serves the connections a worker accepted one after the
// other, so it can not be used for that). /hello.json answers
// {"hello": "world"}.

/** @return array{0: int, 1: string, 2: string} */
function plug3_handle(string $uri): array
{
    $path = (string) parse_url($uri, PHP_URL_PATH);
    if ($path === '/hello.json') {
        return [200, 'application/json', '{"hello": "world"}'];
    }

    return [404, 'application/json', '{"error":"not found"}'];
}

if (PHP_SAPI === 'cli-server') {
    [$status, $type, $body] = plug3_handle($_SERVER['REQUEST_URI']);
    http_response_code($status);
    header('Content-Type: '.$type);
    echo $body;

    return true;
}

/** @return array{0: int, 1: string, 2: string, 3: int} status, type, body, delay in ms */
function plug3_route(string $uri): array
{
    $path = (string) parse_url($uri, PHP_URL_PATH);
    $delay = 0;
    if (preg_match('{^/delay/(\d+)/(dists/[A-Za-z0-9._-]+)$}', $path, $m)) {
        $delay = (int) $m[1];
        $path = '/'.$m[2];
    }
    if (preg_match('{^/(dists/[A-Za-z0-9._-]+)$}', $path, $m) && is_file('/fixtures/'.$m[1])) {
        return [200, 'application/zip', (string) file_get_contents('/fixtures/'.$m[1]), $delay];
    }

    return [404, 'application/json', '{"error":"not found"}', $delay];
}

/**
 * Event loop HTTP/1.0 server: every connection gets its response once its
 * delay elapsed, independently of the others.
 */
function plug3_delay_server(int $port): void
{
    $server = stream_socket_server('tcp://127.0.0.1:'.$port, $errno, $errstr);
    if (!$server) {
        fwrite(STDERR, "listen failed: $errstr\n");
        exit(1);
    }
    stream_set_blocking($server, false);
    $conns = [];
    $seq = 0;
    while (true) {
        $now = microtime(true);
        $read = [$server];
        $write = [];
        $timeout = 1.0;
        foreach ($conns as $id => $c) {
            if ($c['out'] === null) {
                if ($c['due'] === null) {
                    $read[] = $c['s'];
                } elseif ($c['due'] <= $now) {
                    $conns[$id]['out'] = $c['resp'];
                    $write[] = $c['s'];
                } else {
                    $timeout = min($timeout, $c['due'] - $now);
                }
            } else {
                $write[] = $c['s'];
            }
        }
        $except = null;
        $sec = (int) $timeout;
        if (false === @stream_select($read, $write, $except, $sec, (int) (($timeout - $sec) * 1000000))) {
            continue;
        }
        foreach ($read as $s) {
            if ($s === $server) {
                $conn = @stream_socket_accept($server, 0);
                if ($conn) {
                    stream_set_blocking($conn, false);
                    $conns[++$seq] = ['s' => $conn, 'buf' => '', 'due' => null, 'resp' => '', 'out' => null];
                }
                continue;
            }
            foreach ($conns as $id => $c) {
                if ($c['s'] !== $s) {
                    continue;
                }
                $data = fread($s, 65536);
                if ($data === '' || $data === false) {
                    if (feof($s)) {
                        fclose($s);
                        unset($conns[$id]);
                    }
                    break;
                }
                $conns[$id]['buf'] .= $data;
                if (strpos($conns[$id]['buf'], "\r\n\r\n") !== false) {
                    if (!preg_match('{^(GET|HEAD) (\S+) HTTP/1\.[01]}', $conns[$id]['buf'], $rm)) {
                        fclose($s);
                        unset($conns[$id]);
                        break;
                    }
                    [$status, $type, $body, $delay] = plug3_route($rm[2]);
                    $conns[$id]['due'] = microtime(true) + $delay / 1000;
                    $conns[$id]['resp'] = 'HTTP/1.1 '.$status.' '.($status === 200 ? 'OK' : 'Not Found')."\r\nContent-Type: $type\r\nContent-Length: ".strlen($body)."\r\nConnection: close\r\n\r\n".($rm[1] === 'HEAD' ? '' : $body);
                }
                break;
            }
        }
        foreach ($write as $s) {
            foreach ($conns as $id => $c) {
                if ($c['s'] !== $s || $conns[$id]['out'] === null) {
                    continue;
                }
                $n = @fwrite($s, $conns[$id]['out']);
                if ($n === false) {
                    fclose($s);
                    unset($conns[$id]);
                    break;
                }
                $conns[$id]['out'] = (string) substr($conns[$id]['out'], $n);
                if ($conns[$id]['out'] === '') {
                    fclose($s);
                    unset($conns[$id]);
                }
                break;
            }
        }
    }
}

if (($argv[1] ?? '') === 'delay') {
    plug3_delay_server((int) $argv[2]);
    exit(0);
}

if (($argv[1] ?? '') !== 'tls' && ($argv[1] ?? '') !== 'tlsclient') {
    fwrite(STDERR, "usage: php plug3-stub.php tls|tlsclient|delay PORT\n");
    exit(1);
}
$clientAuth = $argv[1] === 'tlsclient';
if (!$clientAuth) {
    // the delay and client certificate servers are ready before the HTTPS port
    // the comparison script waits for
    exec('php '.escapeshellarg(__FILE__).' delay 8901 >/dev/null 2>&1 &');
    exec('php '.escapeshellarg(__FILE__).' tlsclient 8444 >/dev/null 2>&1 &');
    for ($i = 0; $i < 50 && !(@fsockopen('127.0.0.1', 8901) && @fsockopen('127.0.0.1', 8444)); $i++) {
        usleep(100000);
    }
}
$ctx = stream_context_create(['ssl' => $clientAuth ? [
    'local_cert' => '/fixtures/git/_api/server.pem',
    'local_pk' => '/fixtures/git/_api/server.key',
    'verify_peer' => true,
    'verify_peer_name' => false,
    'cafile' => '/fixtures/git/_api/ca.pem',
] : [
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
    [$status, $type, $body] = plug3_handle($rm[2]);
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
