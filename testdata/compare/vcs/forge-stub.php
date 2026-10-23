<?php
// Loopback stub of the GitHub (v3), GitLab (v4) and Bitbucket (2.0) REST APIs
// used by the vcs-forge-* comparison scenarios (tools/compare/scenarios_vcs2.go).
//
// Every answer is computed from the bare fixture repositories in
// /fixtures/git/<repo>.git (built by testdata/compare/vcs/*.sh), so the
// GitHubDriver, GitLabDriver and GitBitbucketDriver of both Composer and
// gocomposer read the same refs, composer.json files, dates and zip archives.
//
// A GitHub repository name ending in "-ratelimited" answers 403 with an
// exhausted rate limit.
// A repository path ending in "-private" maps to the same bare repository but
// answers 404 (403 for Bitbucket) unless the request carries the fixture
// token FORGE_TOKEN (in Authorization: token/Bearer/Basic or PRIVATE-TOKEN).
//
// Interactive authentication (tools/compare/scenarios_tty2.go):
//   GET  https://api.github.com/          token check of GitHub::authorizeOAuthInteractively
//   POST <gitlab>/oauth/token             GitLab password grant (GITLAB_USER / GITLAB_PASSWORD)
//   POST https://bitbucket.org/site/oauth2/access_token
//                                         Bitbucket client credentials (BB_KEY / BB_SECRET)
//   GET  http://127.0.0.1:8900/gitauth/<repo>.git/...
//                                         dumb-http git behind basic auth (GIT_USER / GIT_PASSWORD)
// Each grant answers FORGE_TOKEN as the access token.
//
//   php -S 127.0.0.1:8900 forge-stub.php      plain HTTP (custom gitlab-domains)
//   php forge-stub.php tls 8443                HTTPS for github.com, api.github.com,
//                                              gitlab.com, bitbucket.org, ... (the
//                                              names resolve to 127.0.0.1 through
//                                              podman --add-host)

const FORGE_ROOT = '/fixtures/git';
const FORGE_TOKEN = 's3cr3t-forge-token';
const GITLAB_USER = 'fixture-user';
const GITLAB_PASSWORD = 'gl-s3cr3t-password';
const BB_KEY = 'fixture-consumer-key';
const BB_SECRET = 'bb-s3cr3t-consumer-secret';
const GIT_USER = 'fixture-user';
const GIT_PASSWORD = 'git-s3cr3t-password';

/** @return array{0: string, 1: string}|null basic auth user and password */
function forge_basic(array $h): ?array
{
    if (isset($h['authorization']) && preg_match('{^Basic\s+(.+)$}i', $h['authorization'], $m)) {
        $parts = explode(':', (string) base64_decode($m[1]), 2);

        return [$parts[0], $parts[1] ?? ''];
    }

    return null;
}

function forge_git(string $bare, string $args): string
{
    return (string) shell_exec('git --git-dir='.escapeshellarg($bare).' '.$args.' 2>/dev/null');
}

/** @return array{0: string, 1: bool}|null bare repository path, private flag */
function forge_repo(string $name): ?array
{
    $private = false;
    $name = preg_replace('{\.git$}', '', $name);
    if (str_ends_with($name, '-private')) {
        $private = true;
        $name = substr($name, 0, -8);
    }
    if (!preg_match('{^[A-Za-z0-9._-]+$}', $name) || !is_dir(FORGE_ROOT.'/'.$name.'.git')) {
        return null;
    }

    return [FORGE_ROOT.'/'.$name.'.git', $private];
}

function forge_authorized(array $h): bool
{
    foreach (['authorization', 'private-token'] as $k) {
        if (!isset($h[$k])) {
            continue;
        }
        $v = $h[$k];
        if (str_contains($v, FORGE_TOKEN)) {
            return true;
        }
        if (preg_match('{^Basic\s+(.+)$}i', $v, $m) && str_contains((string) base64_decode($m[1]), FORGE_TOKEN)) {
            return true;
        }
    }

    return false;
}

/** @return array<string, string> name => commit sha, sorted by name */
function forge_refs(string $bare, string $kind): array
{
    $out = [];
    foreach (explode("\n", trim(forge_git($bare, 'for-each-ref --sort=refname --format="%(refname)" refs/'.$kind))) as $ref) {
        if ($ref === '') {
            continue;
        }
        $out[substr($ref, strlen('refs/'.$kind.'/'))] = trim(forge_git($bare, 'rev-parse '.escapeshellarg($ref.'^{commit}')));
    }

    return $out;
}

function forge_commit(string $bare, string $id): ?string
{
    $sha = trim(forge_git($bare, 'rev-parse --verify -q '.escapeshellarg($id.'^{commit}')));

    return preg_match('{^[0-9a-f]{40}$}', $sha) ? $sha : null;
}

function forge_date(string $bare, string $sha, string $fmt): string
{
    return gmdate($fmt, (int) trim(forge_git($bare, 'log -1 --format=%ct '.$sha)));
}

function forge_file(string $bare, string $sha, string $file): ?string
{
    if (trim(forge_git($bare, 'cat-file -t '.escapeshellarg($sha.':'.$file))) !== 'blob') {
        return null;
    }

    return forge_git($bare, 'show '.escapeshellarg($sha.':'.$file));
}

function forge_zip(string $bare, string $sha, string $prefix): string
{
    return forge_git($bare, 'archive --format=zip --prefix='.escapeshellarg($prefix.'/').' '.$sha);
}

function forge_json($data): array
{
    return [200, 'application/json', json_encode($data, JSON_UNESCAPED_SLASHES), []];
}

function forge_404(string $msg = '{"message":"Not Found"}'): array
{
    return [404, 'application/json', $msg, []];
}

/** @return array{0: int, 1: string, 2: string, 3: array<string>} */
function forge_handle(string $host, string $uri, array $h, string $method = 'GET', string $body = ''): array
{
    $host = strtolower(preg_replace('{:\d+$}', '', $host));
    $path = (string) parse_url($uri, PHP_URL_PATH);
    parse_str((string) parse_url($uri, PHP_URL_QUERY), $q);

    // ---- OAuth token grants of the interactive authentications
    if ($method === 'POST') {
        parse_str($body, $form);
        if ($path === '/oauth/token') {
            if (($form['grant_type'] ?? '') === 'refresh_token' && ($form['refresh_token'] ?? null) === 'fixture-refresh-token') {
                return forge_json(['access_token' => FORGE_TOKEN, 'token_type' => 'Bearer', 'expires_in' => 7200, 'refresh_token' => 'fixture-refresh-token-2', 'scope' => 'api', 'created_at' => 1704168245]);
            }
            if (($form['grant_type'] ?? '') === 'password' && ($form['username'] ?? null) === GITLAB_USER && ($form['password'] ?? null) === GITLAB_PASSWORD) {
                return forge_json(['access_token' => FORGE_TOKEN, 'token_type' => 'Bearer', 'expires_in' => 7200, 'refresh_token' => 'fixture-refresh-token', 'scope' => 'api', 'created_at' => 1704164645]);
            }

            return [401, 'application/json', '{"error":"invalid_grant","error_description":"The provided authorization grant is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."}', []];
        }
        if ($host === 'bitbucket.org' && $path === '/site/oauth2/access_token') {
            if (($form['grant_type'] ?? '') === 'client_credentials' && forge_basic($h) === [BB_KEY, BB_SECRET]) {
                return forge_json(['access_token' => FORGE_TOKEN, 'scopes' => 'repository', 'token_type' => 'bearer', 'expires_in' => 7200, 'state' => 'client_credentials', 'refresh_token' => 'fixture-refresh-token']);
            }

            return [400, 'application/json', '{"error_description": "Invalid OAuth client credentials", "error": "invalid_client"}', []];
        }

        return forge_404();
    }

    // ---- token check of GitHub::authorizeOAuthInteractively
    if ($host === 'api.github.com' && $path === '/') {
        return forge_authorized($h)
            ? forge_json(['current_user_url' => 'https://api.github.com/user'])
            : [401, 'application/json', '{"message":"Bad credentials","documentation_url":"https://docs.github.com/rest"}', []];
    }

    // ---- dumb-http git clones from the forge hosts ("-private" ones need
    // FORGE_TOKEN in the basic auth credentials, like the APIs)
    if (in_array($host, ['github.com', 'gitlab.com', 'bitbucket.org'], true)
        && preg_match('{^/[A-Za-z0-9._-]+/([A-Za-z0-9._-]+\.git)(/[A-Za-z0-9._/-]*)$}', $path, $m) && !str_contains($m[2], '..')) {
        $repo = forge_repo($m[1]);
        if ($repo === null) {
            return [404, 'text/plain', 'Not Found', []];
        }
        if ($repo[1] && !forge_authorized($h)) {
            return [401, 'text/plain', 'Unauthorized', ['WWW-Authenticate: Basic realm="fixture git"']];
        }
        $file = $repo[0].$m[2];

        return is_file($file) ? [200, 'application/octet-stream', (string) file_get_contents($file), []] : [404, 'text/plain', 'Not Found', []];
    }

    // ---- dumb-http git repositories behind basic auth
    if (preg_match('{^/gitauth/([A-Za-z0-9._-]+\.git)(/[A-Za-z0-9._/-]*)$}', $path, $m) && !str_contains($m[2], '..')) {
        if (forge_basic($h) !== [GIT_USER, GIT_PASSWORD]) {
            return [401, 'text/plain', 'Unauthorized', ['WWW-Authenticate: Basic realm="fixture git"']];
        }
        $file = FORGE_ROOT.'/'.$m[1].$m[2];
        if (!is_file($file)) {
            return [404, 'text/plain', 'Not Found', []];
        }

        return [200, 'application/octet-stream', (string) file_get_contents($file), []];
    }

    // ---- a composer repository behind authentication (/auth401, /auth403)
    if (preg_match('{^/auth(401|403)/}', $path, $m)) {
        if (forge_authorized($h)) {
            return $path === "/auth$m[1]/packages.json" ? forge_json(['packages' => (object) []]) : forge_404();
        }

        return $m[1] === '401'
            ? [401, 'text/html', 'Unauthorized', ['WWW-Authenticate: Basic realm="fixture"']]
            : [403, 'text/html', 'Forbidden', []];
    }

    // ---- GitHub (api.github.com or <enterprise host>/api/v3)
    $gh = null;
    if ($host === 'api.github.com') {
        $gh = $path;
    } elseif (str_starts_with($path, '/api/v3/')) {
        $gh = substr($path, 7);
    }
    if ($gh !== null && preg_match('{^/repos/([^/]+)/([^/]+)(/.*)?$}', $gh, $m)) {
        [$owner, $name, $rest] = [$m[1], $m[2], $m[3] ?? ''];
        if (str_ends_with($name, '-ratelimited')) {
            // exhausted API rate limit (GitHub answers 403 with these headers)
            return [403, 'application/json', '{"message":"API rate limit exceeded for 127.0.0.1.","documentation_url":"https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting"}',
                ['X-RateLimit-Limit: 60', 'X-RateLimit-Remaining: 0', 'X-RateLimit-Reset: 1704164645']];
        }
        $repo = forge_repo($name);
        if ($repo === null || ($repo[1] && !forge_authorized($h))) {
            return forge_404('{"message":"Not Found","documentation_url":"https://docs.github.com/rest"}');
        }
        [$bare, $private] = $repo;
        if ($rest === '') {
            return forge_json(['name' => $name, 'full_name' => "$owner/$name", 'owner' => ['login' => $owner], 'private' => $private, 'default_branch' => 'main', 'has_issues' => true, 'archived' => false]);
        }
        if ($rest === '/tags') {
            $out = [];
            foreach (forge_refs($bare, 'tags') as $tag => $sha) {
                $out[] = ['name' => $tag, 'commit' => ['sha' => $sha, 'url' => "https://api.github.com/repos/$owner/$name/commits/$sha"]];
            }

            return forge_json($out);
        }
        if ($rest === '/git/refs/heads') {
            $out = [];
            foreach (forge_refs($bare, 'heads') as $branch => $sha) {
                $out[] = ['ref' => 'refs/heads/'.$branch, 'object' => ['sha' => $sha, 'type' => 'commit']];
            }

            return forge_json($out);
        }
        if (preg_match('{^/contents/(.+)$}', $rest, $mm)) {
            $sha = forge_commit($bare, $q['ref'] ?? 'main');
            $content = $sha === null ? null : forge_file($bare, $sha, rawurldecode($mm[1]));
            if ($content === null) {
                return forge_404();
            }

            return forge_json(['type' => 'file', 'encoding' => 'base64', 'content' => chunk_split(base64_encode($content), 60, "\n")]);
        }
        if (preg_match('{^/commits/(.+)$}', $rest, $mm)) {
            $sha = forge_commit($bare, rawurldecode($mm[1]));
            if ($sha === null) {
                return forge_404();
            }

            return forge_json(['sha' => $sha, 'commit' => ['committer' => ['date' => forge_date($bare, $sha, 'Y-m-d\TH:i:s\Z')]]]);
        }
        if (preg_match('{^/zipball/(.+)$}', $rest, $mm)) {
            $sha = forge_commit($bare, rawurldecode($mm[1]));
            if ($sha === null) {
                return forge_404();
            }

            return [200, 'application/zip', forge_zip($bare, $sha, "$owner-$name-".substr($sha, 0, 7)), []];
        }

        return forge_404();
    }

    // ---- GitLab (/api/v4/projects/<urlencoded namespace/project>)
    if (preg_match('{^/api/v4/projects/([^/]+)(/.*)?$}', $path, $m)) {
        $full = rawurldecode($m[1]);
        $rest = $m[2] ?? '';
        $name = basename($full);
        $repo = forge_repo($name);
        if ($repo === null || ($repo[1] && !forge_authorized($h))) {
            return forge_404('{"message":"404 Project Not Found"}');
        }
        [$bare, $private] = $repo;
        $scheme = isset($_SERVER['SERVER_SOFTWARE']) ? 'http' : 'https';
        $origin = $scheme === 'http' ? '127.0.0.1:8900' : $host;
        if ($rest === '') {
            return forge_json([
                'id' => 42, 'default_branch' => 'main', 'visibility' => $private ? 'private' : 'public',
                'web_url' => "$scheme://$origin/$full",
                'http_url_to_repo' => "$scheme://$origin/$full.git",
                'ssh_url_to_repo' => "git@".preg_replace('{:\d+$}', '', $origin).":$full.git",
                'issues_enabled' => true, 'archived' => false,
            ]);
        }
        if (preg_match('{^/repository/(tags|branches)$}', $rest, $mm)) {
            $out = [];
            foreach (forge_refs($bare, $mm[1] === 'tags' ? 'tags' : 'heads') as $ref => $sha) {
                $out[] = ['name' => $ref, 'commit' => ['id' => $sha, 'committed_date' => forge_date($bare, $sha, 'Y-m-d\TH:i:s.000+00:00')]];
            }

            return forge_json($out);
        }
        if (preg_match('{^/repository/files/(.+)/raw$}', $rest, $mm)) {
            $sha = forge_commit($bare, $q['ref'] ?? 'main');
            $content = $sha === null ? null : forge_file($bare, $sha, rawurldecode($mm[1]));
            if ($content === null) {
                return forge_404('{"message":"404 File Not Found"}');
            }

            return [200, 'text/plain', $content, []];
        }
        if ($rest === '/repository/archive.zip') {
            $sha = forge_commit($bare, $q['sha'] ?? 'main');
            if ($sha === null) {
                return forge_404();
            }

            return [200, 'application/zip', forge_zip($bare, $sha, "$name-$sha"), []];
        }

        return forge_404();
    }

    // ---- Bitbucket API
    if ($host === 'api.bitbucket.org' && preg_match('{^/2\.0/repositories/([^/]+)/([^/]+)(/.*)?$}', $path, $m)) {
        [$owner, $name, $rest] = [$m[1], $m[2], $m[3] ?? ''];
        $repo = forge_repo($name);
        if ($repo === null || ($repo[1] && !forge_authorized($h))) {
            return [403, 'application/json', '{"type":"error","error":{"message":"Access denied."}}', []];
        }
        $bare = $repo[0];
        $api = "https://api.bitbucket.org/2.0/repositories/$owner/$name";
        if ($rest === '') {
            return forge_json([
                'scm' => 'git', 'has_issues' => false, 'website' => '', 'mainbranch' => ['name' => 'main', 'type' => 'branch'],
                'links' => [
                    'clone' => [['name' => 'https', 'href' => "https://user@bitbucket.org/$owner/$name.git"], ['name' => 'ssh', 'href' => "git@bitbucket.org:$owner/$name.git"]],
                    'branches' => ['href' => "$api/refs/branches"], 'tags' => ['href' => "$api/refs/tags"],
                    'html' => ['href' => "https://bitbucket.org/$owner/$name"],
                ],
            ]);
        }
        if (preg_match('{^/refs/(tags|branches)$}', $rest, $mm)) {
            $values = [];
            foreach (forge_refs($bare, $mm[1] === 'tags' ? 'tags' : 'heads') as $ref => $sha) {
                $values[] = ['name' => $ref, 'target' => ['hash' => $sha, 'date' => forge_date($bare, $sha, 'c')]];
            }
            usort($values, static fn ($a, $b) => [$b['target']['date'], $a['name']] <=> [$a['target']['date'], $b['name']]);
            foreach ($values as &$v) {
                unset($v['target']['date']);
            }

            return forge_json(['values' => $values]);
        }
        if (preg_match('{^/src/([^/]+)/(.+)$}', $rest, $mm)) {
            $sha = forge_commit($bare, rawurldecode($mm[1]));
            $content = $sha === null ? null : forge_file($bare, $sha, rawurldecode($mm[2]));
            if ($content === null) {
                return forge_404('{"type":"error","error":{"message":"No such file or directory"}}');
            }

            return [200, 'text/plain', $content, []];
        }
        if (preg_match('{^/commit/([^/]+)$}', $rest, $mm)) {
            $sha = forge_commit($bare, rawurldecode($mm[1]));
            if ($sha === null) {
                return forge_404();
            }

            return forge_json(['date' => forge_date($bare, $sha, 'Y-m-d\TH:i:s+00:00')]);
        }

        return forge_404();
    }
    if ($host === 'bitbucket.org' && preg_match('{^/([^/]+)/([^/]+)/get/([^/]+)\.zip$}', $path, $m)) {
        $repo = forge_repo($m[2]);
        if ($repo === null || ($repo[1] && !forge_authorized($h))) {
            return [403, 'text/plain', 'Access denied', []];
        }
        $sha = forge_commit($repo[0], $m[3]);
        if ($sha === null) {
            return forge_404();
        }

        return [200, 'application/zip', forge_zip($repo[0], $sha, "$m[1]-$m[2]-".substr($sha, 0, 12)), []];
    }

    return forge_404();
}

if (PHP_SAPI === 'cli-server') {
    $h = array_change_key_case(function_exists('getallheaders') ? getallheaders() : [], CASE_LOWER);
    [$status, $type, $body, $extra] = forge_handle($_SERVER['HTTP_HOST'] ?? '127.0.0.1', $_SERVER['REQUEST_URI'], $h, $_SERVER['REQUEST_METHOD'] ?? 'GET', (string) file_get_contents('php://input'));
    http_response_code($status);
    header('Content-Type: '.$type);
    foreach ($extra as $line) {
        header($line);
    }
    echo $body;

    return true;
}

if (($argv[1] ?? '') !== 'tls') {
    fwrite(STDERR, "usage: php forge-stub.php tls PORT\n");
    exit(1);
}
$ctx = stream_context_create(['ssl' => [
    'local_cert' => FORGE_ROOT.'/_api/server.pem',
    'local_pk' => FORGE_ROOT.'/_api/server.key',
    'verify_peer' => false,
]]);
$server = stream_socket_server('tls://127.0.0.1:'.(int) $argv[2], $errno, $errstr, STREAM_SERVER_BIND | STREAM_SERVER_LISTEN, $ctx);
if (!$server) {
    fwrite(STDERR, "listen failed: $errstr\n");
    exit(1);
}
$reasons = [200 => 'OK', 400 => 'Bad Request', 401 => 'Unauthorized', 403 => 'Forbidden', 404 => 'Not Found'];
while (true) {
    $conn = @stream_socket_accept($server, -1);
    if (!$conn) {
        continue;
    }
    stream_set_timeout($conn, 10);
    $line = fgets($conn);
    if ($line === false || !preg_match('{^(GET|HEAD|POST) (\S+) HTTP/1\.[01]}', $line, $rm)) {
        fclose($conn);
        continue;
    }
    $h = [];
    while (($l = fgets($conn)) !== false && trim($l) !== '') {
        [$k, $v] = array_pad(explode(':', $l, 2), 2, '');
        $h[strtolower(trim($k))] = trim($v);
    }
    $reqBody = '';
    $len = (int) ($h['content-length'] ?? 0);
    while ($len > 0 && !feof($conn)) {
        $chunk = fread($conn, $len);
        if ($chunk === false || $chunk === '') {
            break;
        }
        $reqBody .= $chunk;
        $len -= strlen($chunk);
    }
    [$status, $type, $body, $extra] = forge_handle($h['host'] ?? '', $rm[2], $h, $rm[1], $reqBody);
    $head = "HTTP/1.1 $status ".($reasons[$status] ?? 'Error')."\r\nContent-Type: $type\r\nContent-Length: ".strlen($body)."\r\nConnection: close\r\n";
    foreach ($extra as $x) {
        $head .= $x."\r\n";
    }
    $data = $head."\r\n".($rm[1] === 'HEAD' ? '' : $body);
    while ($data !== '') {
        $n = @fwrite($conn, $data);
        if ($n === false || $n === 0) {
            break;
        }
        $data = (string) substr($data, $n);
    }
    @fclose($conn);
}
