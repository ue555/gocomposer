<?php
// Builds the static fixtures of the providers-*, archive3-* and warnings-*
// comparison scenarios (tools/compare/scenarios_repo3.go) below
// /fixtures/git/repo3, served by the fixture web server at
// http://127.0.0.1:8899/git/repo3/:
//
//   dists/      tar dists of the provfix/* packages
//   prov/       v1 provider repository: providers-url + provider-includes
//               (sha256 checked), provider files keyed by their hash
//   provroot/   providers listed directly in packages.json
//   provbad/    a provider hash that does not match the served file
//   provlazy/   v1 providers-lazy-url repository (no metadata-url)
//   provlist/   providers + list API
//   arch/       xz / rar / zip dists of the archfix/* packages
//   warn*/      repository responses with warning/info messages
//
// Usage: php repo3-gen.php <out dir>

$out = rtrim($argv[1], '/');
const BASE = 'http://127.0.0.1:8899/git/repo3';
const PATHBASE = '/git/repo3';

function put(string $file, string $data): void
{
    if (!is_dir(dirname($file))) {
        mkdir(dirname($file), 0777, true);
    }
    file_put_contents($file, $data);
}

function enc($data): string
{
    return json_encode($data, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT)."\n";
}

/** builds a tar (optionally xz compressed) of files below a top directory */
function make_tar(string $target, string $top, array $files, bool $xz = false): void
{
    $tmp = sys_get_temp_dir().'/repo3-'.bin2hex(random_bytes(4));
    foreach ($files as $name => $content) {
        put("$tmp/$top/$name", $content);
        if (str_starts_with($name, 'bin/')) {
            chmod("$tmp/$top/$name", 0755);
        }
    }
    if (!is_dir(dirname($target))) {
        mkdir(dirname($target), 0777, true);
    }
    $cmd = 'tar --sort=name --mtime=@1704164645 --owner=0 --group=0 --numeric-owner -C '.escapeshellarg($tmp).' -c'.($xz ? 'J' : '').'f '.escapeshellarg($target).' '.escapeshellarg($top);
    exec($cmd, $o, $rc);
    if ($rc !== 0) {
        fwrite(STDERR, "tar failed: $cmd\n");
        exit(1);
    }
    exec('rm -rf '.escapeshellarg($tmp));
}

/** a stored (uncompressed) RAR 4 archive, see the RAR 2.9 technote */
function make_rar(array $entries): string
{
    $hdr = static function (int $type, int $flags, string $rest): string {
        $body = chr($type).pack('v', $flags).pack('v', 7 + strlen($rest)).$rest;

        return pack('v', crc32($body) & 0xFFFF).$body;
    };
    $out = "Rar!\x1a\x07\x00".$hdr(0x73, 0, pack('v', 0).pack('V', 0));
    $dostime = ((2024 - 1980) << 25) | (1 << 21) | (2 << 16) | (3 << 11) | (4 << 5) | 2;
    foreach ($entries as $name => $data) {
        $isDir = str_ends_with($name, '/');
        $name = rtrim($name, '/');
        $mode = $isDir ? 040755 : (str_contains($name, '/bin/') ? 0100755 : 0100644);
        $data = $isDir ? '' : $data;
        $rest = pack('V', strlen($data)).pack('V', strlen($data)).chr(3).pack('V', $isDir ? 0 : crc32($data))
            .pack('V', $dostime).chr(20).chr(0x30).pack('v', strlen($name)).pack('V', $mode).$name;
        $out .= $hdr(0x74, 0x8000 | ($isDir ? 0x00E0 : 0), $rest).$data;
    }

    return $out.$hdr(0x7B, 0x4000, '');
}

/** a stored zip archive with unix modes; $badCrc corrupts the file CRCs */
function make_zip(array $entries, bool $badCrc = false): string
{
    $local = '';
    $central = '';
    $time = (3 << 11) | (4 << 5) | 2;
    $date = ((2024 - 1980) << 9) | (1 << 5) | 2;
    foreach ($entries as $name => $data) {
        $isDir = str_ends_with($name, '/');
        $mode = $isDir ? 040755 : (str_contains($name, '/bin/') ? 0100755 : 0100644);
        $crc = $isDir ? 0 : crc32($data);
        if ($badCrc && !$isDir) {
            $crc ^= 0x5A5A5A5A;
        }
        $offset = strlen($local);
        $common = pack('v', 10).pack('v', 0).pack('v', 0).pack('v', $time).pack('v', $date)
            .pack('V', $crc).pack('V', strlen($data)).pack('V', strlen($data)).pack('v', strlen($name));
        $local .= "PK\x03\x04".$common.pack('v', 0).$name.$data;
        $central .= "PK\x01\x02".pack('v', 0x031E).$common.pack('v', 0).pack('v', 0).pack('v', 0).pack('v', 0)
            .pack('V', ($mode << 16) | ($isDir ? 0x10 : 0)).pack('V', $offset).$name;
    }

    return $local.$central."PK\x05\x06".pack('v', 0).pack('v', 0).pack('v', count($entries)).pack('v', count($entries))
        .pack('V', strlen($central)).pack('V', strlen($local)).pack('v', 0);
}

// --------------------------------------------------------------- provfix/*
$prov = [
    'provfix/app' => ['1.0.0' => ['description' => 'Provider app', 'require' => ['provfix/lib' => '^1.0'], 'autoload' => ['psr-4' => ['Provfix\\App\\' => 'src/']]]],
    'provfix/lib' => [
        '1.0.0' => ['description' => 'Provider library', 'keywords' => ['provider'], 'autoload' => ['psr-4' => ['Provfix\\Lib\\' => 'src/']]],
        '1.1.0' => ['description' => 'Provider library', 'keywords' => ['provider'], 'autoload' => ['psr-4' => ['Provfix\\Lib\\' => 'src/']]],
        '2.0.0' => ['description' => 'Provider library', 'keywords' => ['provider'], 'autoload' => ['psr-4' => ['Provfix\\Lib\\' => 'src/']]],
    ],
    'provfix/tool' => [
        '1.0.0' => ['description' => 'Provider tool', 'bin' => ['bin/provtool'], 'provide' => ['provfix/tool-implementation' => '1.0']],
        'dev-main' => ['description' => 'Provider tool', 'bin' => ['bin/provtool'], 'extra' => ['branch-alias' => ['dev-main' => '1.1.x-dev']]],
    ],
];
$meta = [];
$uid = 1000;
foreach ($prov as $name => $versions) {
    foreach ($versions as $version => $data) {
        $ref = sha1($name.'@'.$version);
        $file = str_replace('/', '-', $name).'-'.$version.'.tar';
        $top = str_replace('/', '-', $name).'-'.substr($ref, 0, 7);
        $files = ['composer.json' => enc(['name' => $name] + $data)];
        $cls = str_replace(' ', '', ucwords(str_replace(['provfix/', '-'], ['', ' '], $name)));
        if (isset($data['autoload'])) {
            $ns = array_key_first($data['autoload']['psr-4']);
            $files['src/'.$cls.'.php'] = "<?php\nnamespace ".rtrim($ns, '\\').";\nclass $cls { const VERSION = '$version'; }\n";
        }
        if (isset($data['bin'])) {
            $files['bin/provtool'] = "#!/usr/bin/env php\n<?php echo \"provtool $version\\n\";\n";
        }
        make_tar("$out/dists/$file", $top, $files);
        $meta[$name][$version] = ['name' => $name, 'version' => $version] + $data + [
            'dist' => ['type' => 'tar', 'url' => BASE.'/dists/'.$file, 'reference' => $ref, 'shasum' => sha1_file("$out/dists/$file")],
            'type' => 'library',
            'time' => '2024-01-02T03:04:05+00:00',
            'uid' => $uid++,
        ];
    }
}

// provider files, keyed by their sha256 (Packagist v1 layout)
$hashes = [];
foreach ($meta as $name => $versions) {
    $body = enc(['packages' => [$name => $versions]]);
    $hashes[$name] = hash('sha256', $body);
    put("$out/prov/p/$name\$".$hashes[$name].'.json', $body);
}
$includes = [];
foreach (['provider-a' => ['provfix/app', 'provfix/lib'], 'provider-b' => ['provfix/tool']] as $inc => $names) {
    $list = [];
    foreach ($names as $n) {
        $list[$n] = ['sha256' => $hashes[$n]];
    }
    $body = enc(['providers' => $list]);
    $h = hash('sha256', $body);
    put("$out/prov/p/$inc\$$h.json", $body);
    $includes["p/$inc\$%hash%.json"] = ['sha256' => $h];
}
put("$out/prov/packages.json", enc(['packages' => [], 'providers-url' => PATHBASE.'/prov/p/%package%$%hash%.json', 'provider-includes' => $includes]));

$rootProviders = [];
foreach ($hashes as $n => $h) {
    $rootProviders[$n] = ['sha256' => $h];
}
put("$out/provroot/packages.json", enc(['packages' => [], 'providers-url' => PATHBASE.'/prov/p/%package%$%hash%.json', 'providers' => $rootProviders]));

$bad = $rootProviders;
$bad['provfix/lib'] = ['sha256' => hash('sha256', 'tampered')];
put("$out/provbad/p/provfix/lib\$".hash('sha256', 'tampered').'.json', enc(['packages' => ['provfix/lib' => $meta['provfix/lib']]]));
put("$out/provbad/packages.json", enc(['packages' => [], 'providers-url' => PATHBASE.'/provbad/p/%package%$%hash%.json', 'providers' => $bad]));
foreach ($hashes as $n => $h) {
    if ($n !== 'provfix/lib') {
        put("$out/provbad/p/$n\$$h.json", (string) file_get_contents("$out/prov/p/$n\$$h.json"));
    }
}

foreach ($meta as $name => $versions) {
    put("$out/provlazy/p/$name.json", enc(['packages' => [$name => $versions]]));
}
put("$out/provlazy/packages.json", enc(['packages' => [], 'providers-lazy-url' => PATHBASE.'/provlazy/p/%package%.json']));

put("$out/provlist/packages.json", enc(['packages' => [], 'providers-url' => PATHBASE.'/prov/p/%package%$%hash%.json', 'provider-includes' => $includes, 'list' => PATHBASE.'/provlist/list.json']));
put("$out/provlist/list.json", enc(['packageNames' => array_keys($meta)]));
foreach ($includes as $inc => $m) {
    put("$out/provlist/".str_replace('%hash%', $m['sha256'], $inc), (string) file_get_contents("$out/prov/".str_replace('%hash%', $m['sha256'], $inc)));
}

// -------------------------------------------------------------- archfix/*
$archFiles = static function (string $name, string $top): array {
    return [
        "$top/" => '',
        "$top/composer.json" => enc(['name' => $name, 'bin' => ['bin/archtool'], 'autoload' => ['psr-4' => ['Archfix\\' => 'src/']]]),
        "$top/src/" => '',
        "$top/src/".ucfirst(substr($name, 8)).'.php' => "<?php\nnamespace Archfix;\nclass ".ucfirst(substr($name, 8))." {}\n",
        "$top/bin/" => '',
        "$top/bin/archtool" => "#!/bin/sh\necho archtool\n",
    ];
};
$arch = [];
$addArch = static function (string $name, string $type, string $file) use (&$arch, $out): void {
    $arch[$name] = ['1.0.0' => [
        'name' => $name, 'version' => '1.0.0', 'type' => 'library', 'bin' => ['bin/archtool'],
        'autoload' => ['psr-4' => ['Archfix\\' => 'src/']],
        'dist' => ['type' => $type, 'url' => BASE.'/arch/dists/'.$file, 'reference' => sha1($name), 'shasum' => sha1_file("$out/arch/dists/$file")],
        'time' => '2024-01-02T03:04:05+00:00',
    ]];
};
foreach (['xzlib' => 'xz', 'rarlib' => 'rar', 'ziplib' => 'zip', 'zipcrc' => 'zip', 'rarbroken' => 'rar', 'xzbroken' => 'xz'] as $short => $type) {
    $name = "archfix/$short";
    $top = "archfix-$short-1.0.0";
    $entries = $archFiles($name, $top);
    if ($type === 'xz') {
        $file = "$top.tar.xz";
        $plain = [];
        foreach ($entries as $k => $v) {
            if (!str_ends_with($k, '/')) {
                $plain[substr($k, strlen($top) + 1)] = $v;
            }
        }
        make_tar("$out/arch/dists/$file", $top, $plain, true);
        if ($short === 'xzbroken') {
            $data = (string) file_get_contents("$out/arch/dists/$file");
            put("$out/arch/dists/$file", substr($data, 0, (int) (strlen($data) / 2)));
        }
    } elseif ($type === 'rar') {
        $file = "$top.rar";
        $data = make_rar($entries);
        if ($short === 'rarbroken') {
            $data = substr($data, 0, (int) (strlen($data) / 2));
        }
        put("$out/arch/dists/$file", $data);
    } else {
        $file = "$top.zip";
        put("$out/arch/dists/$file", make_zip($entries, $short === 'zipcrc'));
    }
    $addArch($name, $type, $file);
}
put("$out/arch/packages.json", enc(['packages' => $arch]));

// ---------------------------------------------------------------- warn*/
$v2 = static function (string $dir, array $extra, array $metaExtra = []) use ($out, $meta): void {
    foreach ($meta as $name => $versions) {
        $stable = [];
        $dev = [];
        foreach (array_reverse($versions, true) as $v => $data) {
            if (str_starts_with($v, 'dev-')) {
                $dev[] = $data;
            } else {
                $stable[] = $data;
            }
        }
        put("$out/$dir/p2/$name.json", enc(['packages' => [$name => $stable]] + $metaExtra));
        if ($dev) {
            put("$out/$dir/p2/$name~dev.json", enc(['packages' => [$name => $dev]] + $metaExtra));
        }
    }
    put("$out/$dir/packages.json", enc(['packages' => [], 'metadata-url' => PATHBASE."/$dir/p2/%package%.json"] + $extra));
};
$avail = ['available-packages' => array_keys($meta)];
$v2('warn', $avail + [
    'warning' => "Root \u{1b}[31mwarning\u{1b}[0m text",
    'warning-versions' => '>=2.0',
    'info' => 'Root info for old Composer versions',
    'info-versions' => '<2.0',
    'warnings' => [
        ['versions' => '>=2.5', 'message' => 'Modern warning'],
        ['versions' => '<2.0', 'message' => 'Hidden modern warning'],
    ],
    'infos' => [['versions' => '*', 'message' => 'Modern info']],
], ['info' => 'Metadata info']);
$v2('warnsearch', [
    'search' => PATHBASE.'/warnsearch/search.json?q=%query%&type=%type%',
    'list' => PATHBASE.'/warnsearch/list.json',
]);
put("$out/warnsearch/search.json", enc([
    'results' => [['name' => 'provfix/lib', 'description' => 'Provider library', 'url' => BASE.'/provfix/lib', 'repository' => '', 'downloads' => 0, 'favers' => 0]],
    'total' => 1,
    'warning' => 'Search API warning',
]));
put("$out/warnsearch/list.json", enc(['packageNames' => array_keys($meta), 'info' => 'List API info']));
$v2('warnadv', $avail + ['security-advisories' => ['metadata' => false, 'api-url' => PATHBASE.'/warnadv/advisories.json']]);
put("$out/warnadv/advisories.json", enc([
    'advisories' => ['provfix/lib' => [[
        'advisoryId' => 'PKSA-prov-0001', 'packageName' => 'provfix/lib', 'remoteId' => 'GHSA-prov-0001',
        'title' => 'Provider library issue', 'link' => 'https://example.org/advisories/prov-1', 'cve' => 'CVE-2024-0002',
        'affectedVersions' => '>=1.0.0,<1.1.0', 'source' => 'GitHub', 'reportedAt' => '2024-01-02 03:04:05',
        'composerRepository' => 'https://packagist.org', 'severity' => 'medium',
        'sources' => [['name' => 'GitHub', 'remoteId' => 'GHSA-prov-0001']],
    ]]],
    'warnings' => [['versions' => '>=2.0', 'message' => 'Advisory API warning']],
]));
$v2('warnfilter', $avail + ['filter' => ['metadata' => true, 'lists' => ['malware' => ['enabled' => true]], 'api-url' => PATHBASE.'/warnfilter/filter.json']]);
put("$out/warnfilter/filter.json", enc(['filter' => ['malware' => []], 'info' => 'Filter API info', 'warning' => 'Filter API warning']));
// v1 provider repository whose provider includes carry a warning
put("$out/warnprov/packages.json", enc(['packages' => [], 'providers-url' => PATHBASE.'/prov/p/%package%$%hash%.json', 'provider-includes' => $includes, 'warning' => 'Provider root warning']));
foreach ($includes as $inc => $m) {
    put("$out/warnprov/".str_replace('%hash%', $m['sha256'], $inc), (string) file_get_contents("$out/prov/".str_replace('%hash%', $m['sha256'], $inc)));
}
