<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Util\HttpDownloader for plugins, based
// on Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Util/HttpDownloader.php (MIT license, (c) Nils Adermann, Jordi
// Boggiano). The job queue, the promises and the sync/async API are ported;
// the transfers are made by gocomposer's Go HTTP client (the same one that
// downloads packages: auth, cafile, secure-http, retries and error messages
// like CurlDownloader). The jobs queued when wait()/countActiveJobs() runs are
// transferred concurrently by Go, then settled in queue order.

namespace Composer\Util;

use Composer\Config;
use Composer\Downloader\TransportException;
use Composer\IO\IOInterface;
use Composer\Pcre\Preg;
use Composer\Util\Http\Response;
use GocomposerPluginHost\Host;
use React\Promise\Promise;
use React\Promise\PromiseInterface;

class HttpDownloader
{
    private const STATUS_QUEUED = 1;
    private const STATUS_STARTED = 2;
    private const STATUS_COMPLETED = 3;
    private const STATUS_FAILED = 4;
    private const STATUS_ABORTED = 5;

    /** @var IOInterface */
    private $io;
    /** @var Config */
    private $config;
    /** @var array<int, array<string, mixed>> */
    private $jobs = [];
    /** @var mixed[] */
    private $options = [];
    /** @var int */
    private $runningJobs = 0;
    /** @var int */
    private $maxJobs = 12;
    /** @var int */
    private $idGen = 0;
    /** @var bool */
    private $disabled;
    /** @var bool */
    private $allowAsync = false;

    /**
     * @param mixed[] $options
     */
    public function __construct(IOInterface $io, Config $config, array $options = [], bool $disableTls = false)
    {
        $this->io = $io;

        $this->disabled = (bool) Platform::getEnv('COMPOSER_DISABLE_NETWORK');

        // the TLS defaults (cafile/capath) are applied by gocomposer's client
        $this->options = $options;
        $this->config = $config;

        if (is_numeric($maxJobs = Platform::getEnv('COMPOSER_MAX_PARALLEL_HTTP'))) {
            $this->maxJobs = max(1, min(50, (int) $maxJobs));
        }
    }

    /**
     * @param mixed[] $options
     * @return Response
     */
    public function get(string $url, array $options = [])
    {
        if ('' === $url) {
            throw new \InvalidArgumentException('$url must not be an empty string');
        }
        [$job, $promise] = $this->addJob(['url' => $url, 'options' => $options, 'copyTo' => null], true);
        $promise->then(null, static function (\Throwable $e) {
            // suppress error as it is rethrown to the caller by getResponse() a few lines below
        });
        $this->wait($job['id']);

        return $this->getResponse($job['id']);
    }

    /**
     * @param mixed[] $options
     * @return PromiseInterface
     */
    public function add(string $url, array $options = [])
    {
        if ('' === $url) {
            throw new \InvalidArgumentException('$url must not be an empty string');
        }
        [, $promise] = $this->addJob(['url' => $url, 'options' => $options, 'copyTo' => null]);

        return $promise;
    }

    /**
     * @param mixed[] $options
     * @return Response
     */
    public function copy(string $url, string $to, array $options = [])
    {
        if ('' === $url) {
            throw new \InvalidArgumentException('$url must not be an empty string');
        }
        [$job, $promise] = $this->addJob(['url' => $url, 'options' => $options, 'copyTo' => $to], true);
        $promise->then(null, static function (\Throwable $e) {
            // suppress error as it is rethrown to the caller by getResponse() a few lines below
        });
        $this->wait($job['id']);

        return $this->getResponse($job['id']);
    }

    /**
     * @param mixed[] $options
     * @return PromiseInterface
     */
    public function addCopy(string $url, string $to, array $options = [])
    {
        if ('' === $url) {
            throw new \InvalidArgumentException('$url must not be an empty string');
        }
        [, $promise] = $this->addJob(['url' => $url, 'options' => $options, 'copyTo' => $to]);

        return $promise;
    }

    /**
     * @return mixed[]
     */
    public function getOptions()
    {
        return $this->options;
    }

    /**
     * @param mixed[] $options
     * @return void
     */
    public function setOptions(array $options)
    {
        $this->options = array_replace_recursive($this->options, $options);
    }

    /**
     * @param array<string, mixed> $request
     * @return array{0: array<string, mixed>, 1: PromiseInterface}
     */
    private function addJob(array $request, bool $sync = false): array
    {
        $request['options'] = array_replace_recursive($this->options, $request['options']);

        $job = [
            'id' => $this->idGen++,
            'status' => self::STATUS_QUEUED,
            'request' => $request,
            'sync' => $sync,
            'origin' => Url::getOrigin($this->config, $request['url']),
        ];

        if (!$sync && !$this->allowAsync) {
            throw new \LogicException('You must use the HttpDownloader instance which is part of a Composer\Loop instance to be able to run async http requests');
        }

        // capture username/password from URL if there is one
        if (Preg::isMatchStrictGroups('{^https?://([^:/]+):([^@/]+)@([^/]+)}i', $request['url'], $match)) {
            $this->io->setAuthentication($job['origin'], rawurldecode($match[1]), rawurldecode($match[2]));
        }

        $resolver = static function ($resolve, $reject) use (&$job): void {
            $job['status'] = HttpDownloader::STATUS_QUEUED;
            $job['resolve'] = $resolve;
            $job['reject'] = $reject;
        };

        $canceler = static function () use (&$job): void {
            if ($job['status'] === HttpDownloader::STATUS_QUEUED) {
                $job['status'] = HttpDownloader::STATUS_ABORTED;
            }
            if ($job['status'] !== HttpDownloader::STATUS_STARTED) {
                return;
            }
            $job['status'] = HttpDownloader::STATUS_ABORTED;

            throw new \Composer\Downloader\TransportException('Download of ' . Url::sanitize($job['request']['url']) . ' canceled');
        };

        $promise = new Promise($resolver, $canceler);
        $promise = $promise->then(function ($response) use (&$job) {
            $job['status'] = HttpDownloader::STATUS_COMPLETED;
            $job['response'] = $response;

            $this->markJobDone();

            return $response;
        }, function ($e) use (&$job): void {
            $job['status'] = HttpDownloader::STATUS_FAILED;
            $job['exception'] = $e;

            $this->markJobDone();

            throw $e;
        });
        $this->jobs[$job['id']] = &$job;

        return [$job, $promise];
    }

    private function markJobDone(): void
    {
        $this->runningJobs--;
    }

    /**
     * Transfers the queued jobs (gocomposer runs them concurrently, at most
     * maxJobs at a time) and settles their promises in queue order.
     */
    private function runQueued(): void
    {
        $batch = [];
        foreach ($this->jobs as $id => $job) {
            if ($job['status'] !== self::STATUS_QUEUED) {
                continue;
            }
            $this->jobs[$id]['status'] = self::STATUS_STARTED;
            $this->runningJobs++;
            if ($this->disabled) {
                $headers = $job['request']['options']['http']['header'] ?? [];
                if (false !== stripos(implode('', (array) $headers), 'if-modified-since')) {
                    $this->jobs[$id]['resolve'](new Response(['url' => $job['request']['url']], 304, [], ''));
                } else {
                    $e = new TransportException('Network disabled, request canceled: '.Url::sanitize($job['request']['url']), 499);
                    $e->setStatusCode(499);
                    $this->jobs[$id]['reject']($e);
                }
                continue;
            }
            $batch[$id] = ['url' => $job['request']['url'], 'options' => $job['request']['options'], 'copyTo' => $job['request']['copyTo']];
        }
        if ($batch === []) {
            return;
        }
        $results = Host::cb('http', ['jobs' => array_values($batch), 'max' => $this->maxJobs]);
        $i = 0;
        foreach ($batch as $id => $request) {
            $r = $results[$i++] ?? null;
            $job = $this->jobs[$id];
            if (!is_array($r)) {
                $job['reject'](new TransportException('gocomposer: no response for '.Url::sanitize($request['url'])));
                continue;
            }
            if (isset($r['error'])) {
                $job['reject'](self::exception($r['error']));
                continue;
            }
            $body = null;
            if (isset($r['body'])) {
                $body = (string) base64_decode($r['body']);
            }
            $job['resolve'](new Response(['url' => $request['url']], (int) $r['status'], $r['headers'] ?? [], $body));
        }
    }

    /**
     * @param array<string, mixed> $e
     */
    private static function exception(array $e): \Exception
    {
        $class = $e['class'] ?? TransportException::class;
        if ($class !== TransportException::class) {
            return new \RuntimeException((string) $e['message'], (int) ($e['code'] ?? 0));
        }
        $ex = new TransportException((string) $e['message'], (int) ($e['code'] ?? 0));
        if (isset($e['status'])) {
            $ex->setStatusCode((int) $e['status']);
        }
        if (isset($e['headers'])) {
            $ex->setHeaders($e['headers']);
        }
        if (isset($e['body'])) {
            $ex->setResponse((string) base64_decode($e['body']));
        }

        return $ex;
    }

    /**
     * @param int|null $index For internal use only, the job id
     * @return void
     */
    public function wait(?int $index = null)
    {
        do {
            $jobCount = $this->countActiveJobs($index);
        } while ($jobCount);
    }

    /**
     * @internal
     */
    public function enableAsync(): void
    {
        $this->allowAsync = true;
    }

    /**
     * @internal
     *
     * @param  int|null $index For internal use only, the job id
     * @return int      number of active (queued or started) jobs
     */
    public function countActiveJobs(?int $index = null): int
    {
        // the background steps gocomposer performed already (see Api::defer)
        \GocomposerPluginHost\Api::drainDeferred();
        $this->runQueued();

        if (null !== $index) {
            return $this->jobs[$index]['status'] < self::STATUS_COMPLETED ? 1 : 0;
        }

        $active = 0;
        foreach ($this->jobs as $job) {
            if ($job['status'] < self::STATUS_COMPLETED) {
                $active++;
            } elseif (!$job['sync']) {
                unset($this->jobs[$job['id']]);
            }
        }

        return $active;
    }

    private function getResponse(int $index): Response
    {
        if (!isset($this->jobs[$index])) {
            throw new \LogicException('Invalid request id');
        }

        if ($this->jobs[$index]['status'] === self::STATUS_FAILED) {
            assert(isset($this->jobs[$index]['exception']));
            throw $this->jobs[$index]['exception'];
        }

        if (!isset($this->jobs[$index]['response'])) {
            throw new \LogicException('Response not available yet, call wait() first');
        }

        $resp = $this->jobs[$index]['response'];

        unset($this->jobs[$index]);

        return $resp;
    }

    /**
     * @internal
     * @deprecated use ResponseWarnings::output() instead
     * @param  mixed $data
     */
    public static function outputWarnings(IOInterface $io, string $url, $data): bool
    {
        return \Composer\Util\Http\ResponseWarnings::output($io, $url, $data);
    }

    /**
     * @internal
     * @return ?string[]
     */
    public static function getExceptionHints(\Throwable $e): ?array
    {
        return null;
    }

    /**
     * @internal
     */
    public static function isCurlEnabled(): bool
    {
        return \extension_loaded('curl') && \function_exists('curl_multi_exec') && \function_exists('curl_multi_init');
    }
}
