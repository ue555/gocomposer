<?php
// Loopback fault server for the curl-* comparison scenarios
// (tools/compare/scenarios_term.go): every mode misbehaves in one way a
// transfer can fail, so the "curl error N while downloading ..." messages of
// Composer's CurlDownloader can be compared with gocomposer's.
//
//   php term-fault-server.php MODE PORT [DIR]
//
// DIR holds the certificates made by termfault.sh (tls-* modes).
//   hang          accepts and never answers
//   close         reads the request, closes without a reply (curl 52)
//   reset         closes with the request unread, the kernel resets (curl 56)
//   garbage       answers with a line that is not HTTP (curl 1)
//   partial       announces 100 bytes, sends 10 and closes (curl 18)
//   tls-close     reads the ClientHello and closes (curl 35)
//   tls-garbage   answers the ClientHello with plain HTTP (curl 35)
//   tls-hang      accepts and never starts the TLS handshake
//   tls-<cert>    a TLS server presenting DIR/<cert>.pem (self-signed,
//                 wrong-host, expired, unknown-ca: curl 60)

$mode = $argv[1] ?? '';
$port = (int) ($argv[2] ?? 0);
$dir = $argv[3] ?? __DIR__;

function readHead($conn): void
{
    $buf = '';
    while (!feof($conn) && strpos($buf, "\r\n\r\n") === false) {
        $chunk = fread($conn, 8192);
        if ($chunk === false || $chunk === '') {
            break;
        }
        $buf .= $chunk;
    }
}

if (strpos($mode, 'tls-') === 0 && !in_array($mode, ['tls-close', 'tls-garbage', 'tls-hang'], true)) {
    $cert = substr($mode, 4);
    $ctx = stream_context_create(['ssl' => [
        'local_cert' => $dir.'/'.$cert.'.pem',
        'local_pk' => $dir.'/'.$cert.'.key',
        'verify_peer' => false,
    ]]);
    $server = stream_socket_server('tls://127.0.0.1:'.$port, $errno, $errstr, STREAM_SERVER_BIND | STREAM_SERVER_LISTEN, $ctx);
    if (!$server) {
        fwrite(STDERR, "listen: $errstr\n");
        exit(1);
    }
    while (true) {
        $conn = @stream_socket_accept($server, -1);
        if (!$conn) {
            continue;
        }
        readHead($conn);
        fwrite($conn, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 2\r\nConnection: close\r\n\r\n{}");
        fclose($conn);
    }
}

$server = stream_socket_server('tcp://127.0.0.1:'.$port, $errno, $errstr);
if (!$server) {
    fwrite(STDERR, "listen: $errstr\n");
    exit(1);
}
$held = [];
while (true) {
    $conn = @stream_socket_accept($server, -1);
    if (!$conn) {
        continue;
    }
    switch ($mode) {
        case 'hang':
        case 'tls-hang':
            $held[] = $conn;
            break;
        case 'close':
            readHead($conn);
            fclose($conn);
            break;
        case 'reset':
            // unread data in the receive buffer makes close() send a RST
            $r = [$conn];
            $w = $e = null;
            stream_select($r, $w, $e, 5);
            usleep(50000);
            fclose($conn);
            break;
        case 'garbage':
            readHead($conn);
            fwrite($conn, "this is not http\r\n");
            fclose($conn);
            break;
        case 'partial':
            readHead($conn);
            fwrite($conn, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 100\r\nConnection: close\r\n\r\n{\"packages\"");
            fflush($conn);
            usleep(100000);
            fclose($conn);
            break;
        case 'tls-close':
            fread($conn, 8192);
            fclose($conn);
            break;
        case 'tls-garbage':
            fread($conn, 8192);
            fwrite($conn, "HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\nConnection: close\r\n\r\n");
            fclose($conn);
            break;
        default:
            fwrite(STDERR, "unknown mode $mode\n");
            exit(1);
    }
}
