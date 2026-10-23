<?php declare(strict_types=1);
// Records Composer\Json\JsonManipulator behaviour on official PHP for
// internal/jsonmanip. Composer f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Usage:
//   php jsonmanip-dump.php tests
//       runs tests/Composer/Test/Json/JsonManipulatorTest.php unchanged against a
//       recording JsonManipulator wrapper and prints every constructor/method
//       call with arguments and results plus every assertion (expected value
//       and which observation it checked).
//   php jsonmanip-dump.php corpus file.json
//       replays {"cases":[{"contents":..., "ops":[{"m":..., "args":[...]}]}]}
//       and prints the observations of each case.
//   php jsonmanip-dump.php configsource file.json
//       runs Composer\Config\JsonConfigSource operations on a scratch file
//       ({"cases":[{"file":..., "contents":...|null, "auth":bool, "ops":[...]}]})
//       and records the file contents, permissions and exceptions after each
//       operation (the scratch directory is shown as @DIR@).
// Values are encoded with type tags so Go can rebuild PHP arrays/objects exactly.

namespace {
    require '/oracle-source/vendor/autoload.php';

    set_error_handler(static function (): bool { return true; });

    final class Rec
    {
        /** @var list<array<string, mixed>> */
        public static $events = [];
        /** @var list<mixed> */
        public static $values = [];

        public static function reset(): void
        {
            self::$events = [];
            self::$values = [];
        }

        /**
         * @param mixed[] $args
         * @return mixed
         */
        public static function obs(string $op, array $args, callable $fn)
        {
            $ev = ['op' => $op, 'args' => array_map([self::class, 'enc'], $args)];
            try {
                $r = $fn();
            } catch (\Throwable $e) {
                $ev['exc'] = ['class' => get_class($e), 'message' => $e->getMessage()];
                self::$events[] = $ev;
                self::$values[] = null;
                throw $e;
            }
            $ev['ret'] = self::enc($r);
            self::$events[] = $ev;
            self::$values[] = $r;

            return $r;
        }

        /** @param mixed $expected @param mixed $actual */
        public static function assert(string $kind, $expected, $actual): void
        {
            $idx = null;
            if (is_string($actual) || is_bool($actual)) {
                for ($i = count(self::$values) - 1; $i >= 0; $i--) {
                    if (self::$values[$i] === $actual) {
                        $idx = $i;
                        break;
                    }
                }
            }
            // index into the observation list (events of kind op)
            self::$events[] = ['assert' => $kind, 'expected' => self::enc($expected), 'obs' => $idx];
        }

        /** @param mixed $v */
        public static function enc($v)
        {
            if ($v === null || is_bool($v) || is_string($v)) {
                return is_string($v) && preg_match('//u', $v) !== 1 ? ['t' => 'bytes', 'b64' => base64_encode($v)] : $v;
            }
            if (is_int($v)) {
                return ['t' => 'int', 'v' => $v];
            }
            if (is_float($v)) {
                return ['t' => 'float', 'v' => var_export($v, true)];
            }
            if (is_array($v)) {
                $items = [];
                foreach ($v as $k => $x) {
                    $items[] = [(string) $k, self::enc($x)];
                }
                return ['t' => 'array', 'items' => $items];
            }
            if ($v instanceof \stdClass || $v instanceof \ArrayObject) {
                $items = [];
                foreach ((array) $v as $k => $x) {
                    $items[] = [(string) $k, self::enc($x)];
                }
                return ['t' => 'object', 'class' => get_class($v), 'items' => $items];
            }
            throw new \LogicException('cannot encode ' . get_debug_type($v));
        }

        /** @param mixed $v @return mixed */
        public static function dec($v)
        {
            if (!is_array($v)) {
                return $v;
            }
            switch ($v['t']) {
                case 'int':
                    return (int) $v['v'];
                case 'float':
                    return (float) $v['v'];
                case 'bytes':
                    return base64_decode($v['b64']);
                case 'array':
                    $a = [];
                    foreach ($v['items'] as [$k, $x]) {
                        $a[$k] = self::dec($x);
                    }
                    return $a;
                case 'object':
                    $o = $v['class'] === 'ArrayObject' ? new \ArrayObject() : new \stdClass();
                    foreach ($v['items'] as [$k, $x]) {
                        if ($o instanceof \ArrayObject) {
                            $o[$k] = self::dec($x);
                        } else {
                            $o->{$k} = self::dec($x);
                        }
                    }
                    return $o;
            }
            throw new \LogicException('bad tag');
        }
    }

    // load the real class under another name (not for configsource, which
    // uses Composer's classes unchanged)
    $src = ($argv[1] ?? '') === 'configsource' ? '' : file_get_contents('/oracle-source/src/Composer/Json/JsonManipulator.php');
    if ($src !== '') {
        $src = preg_replace('{^<\?php\s*}', '', $src);
        $src = str_replace("\nclass JsonManipulator\n", "\nclass RealJsonManipulator\n", $src, $n);
        if ($n !== 1) {
            throw new \LogicException('could not rename JsonManipulator');
        }
        eval($src);
    }
}

namespace Composer\Json {
    if (($argv[1] ?? '') !== 'configsource') {
    class JsonManipulator
    {
        /** @var RealJsonManipulator */
        private $real;

        public function __construct(string $contents)
        {
            \Rec::obs('new', [$contents], function () use ($contents) {
                $this->real = new RealJsonManipulator($contents);

                return null;
            });
        }

        /** @param mixed[] $args @return mixed */
        public function __call(string $name, array $args)
        {
            return \Rec::obs($name, $args, function () use ($name, $args) {
                return $this->real->{$name}(...$args);
            });
        }
    }
    }
}

namespace Composer\Test {
    final class AssertionFailed extends \Exception
    {
    }

    abstract class TestCase
    {
        public function __construct()
        {
        }

        /** @param mixed $e @param mixed $a */
        public static function assertEquals($e, $a, string $m = ''): void
        {
            \Rec::assert('equals', $e, $a);
            if ($e != $a) {
                throw new AssertionFailed('assertEquals failed');
            }
        }

        /** @param mixed $e @param mixed $a */
        public static function assertSame($e, $a, string $m = ''): void
        {
            \Rec::assert('same', $e, $a);
            if ($e !== $a) {
                throw new AssertionFailed('assertSame failed');
            }
        }

        /** @param mixed $a */
        public static function assertTrue($a, string $m = ''): void
        {
            \Rec::assert('same', true, $a);
            if ($a !== true) {
                throw new AssertionFailed('assertTrue failed');
            }
        }

        public static function assertJsonStringEqualsJsonString(string $e, string $a, string $m = ''): void
        {
            \Rec::assert('json', $e, $a);
            if (json_decode($e) != json_decode($a)) {
                throw new AssertionFailed('assertJsonStringEqualsJsonString failed');
            }
        }
    }
}

namespace {
    $mode = $argv[1] ?? '';
    $out = [];
    if ($mode === 'tests') {
        require '/oracle-source/tests/Composer/Test/Json/JsonManipulatorTest.php';
        $class = \Composer\Test\Json\JsonManipulatorTest::class;
        foreach ((new \ReflectionClass($class))->getMethods() as $rm) {
            $method = $rm->getName();
            if (strpos($method, 'test') !== 0) {
                continue;
            }
            $sets = ['' => []];
            if (preg_match('{@dataProvider\s+(\S+)}', (string) $rm->getDocComment(), $pm)) {
                $sets = [];
                foreach ($class::{$pm[1]}() as $setName => $args) {
                    $sets[(string) $setName] = $args;
                }
            }
            foreach ($sets as $setName => $args) {
                \Rec::reset();
                $status = 'pass';
                $error = null;
                try {
                    (new $class())->{$method}(...array_values($args));
                } catch (\Throwable $e) {
                    $status = 'fail';
                    $error = get_class($e) . ': ' . $e->getMessage();
                }
                $out[] = ['test' => $method, 'set' => (string) $setName, 'status' => $status, 'error' => $error, 'events' => \Rec::$events];
            }
        }
    } elseif ($mode === 'corpus') {
        $in = json_decode(file_get_contents($argv[2]), true, 512, JSON_THROW_ON_ERROR);
        foreach ($in['cases'] as $case) {
            \Rec::reset();
            try {
                $m = new \Composer\Json\JsonManipulator(\Rec::dec($case['contents']));
                foreach ($case['ops'] as $op) {
                    $args = array_map([\Rec::class, 'dec'], $op['args']);
                    try {
                        $m->{$op['m']}(...$args);
                    } catch (\Throwable $e) {
                        // recorded; the remaining operations still run
                    }
                    $m->getContents();
                }
            } catch (\Throwable $e) {
            }
            $out[] = ['events' => \Rec::$events];
        }
    } elseif ($mode === 'configsource') {
        $in = json_decode(file_get_contents($argv[2]), true, 512, JSON_THROW_ON_ERROR);
        $fs = new \Composer\Util\Filesystem();
        umask(022);
        foreach ($in['cases'] as $case) {
            $dir = sys_get_temp_dir() . '/cs-' . bin2hex(random_bytes(6));
            mkdir($dir);
            $path = $dir . '/' . $case['file'];
            if ($case['contents'] !== null) {
                file_put_contents($path, \Rec::dec($case['contents']));
                chmod($path, 0644);
            }
            $source = new \Composer\Config\JsonConfigSource(new \Composer\Json\JsonFile($path), $case['auth']);
            $events = [];
            foreach ($case['ops'] as $op) {
                $args = array_map([\Rec::class, 'dec'], $op['args']);
                $ev = ['op' => $op['m']];
                try {
                    $source->{$op['m']}(...$args);
                } catch (\Throwable $e) {
                    $ev['exc'] = ['class' => get_class($e), 'message' => str_replace($dir, '@DIR@', $e->getMessage())];
                }
                clearstatcache();
                $ev['contents'] = is_file($path) ? \Rec::enc(file_get_contents($path)) : null;
                $ev['mode'] = is_file($path) ? fileperms($path) & 0777 : null;
                $events[] = $ev;
            }
            $fs->removeDirectory($dir);
            $out[] = ['events' => $events];
        }
    } else {
        fwrite(STDERR, "usage: jsonmanip-dump.php tests|corpus|configsource file\n");
        exit(1);
    }
    $defines = (new \ReflectionClass($mode === 'configsource' ? \Composer\Json\JsonManipulator::class : \Composer\Json\RealJsonManipulator::class))->getConstant('DEFINES');
    echo json_encode(['php' => PHP_VERSION, 'pcre' => PCRE_VERSION, 'jit' => (bool) ini_get('pcre.jit'), 'defines' => $defines, 'results' => $out], JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE | JSON_THROW_ON_ERROR), "\n";
}
