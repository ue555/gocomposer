package composer

// FileDownloader::__construct's cache garbage collection, run when
// Factory::createDownloadManager builds the files cache
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"gocomposer/internal/cache"
	"gocomposer/internal/cio"
	"gocomposer/internal/config"
)

// filesCacheGC creates the files cache only when cache-files-ttl > 0 and
// collects it with probability 1/51 per process (Cache::gcIsNecessary).
func filesCacheGC(io *cio.IO, cfg *config.Config) error {
	ttl, err := cfg.GetFlags("cache-files-ttl", 0)
	if err != nil {
		return err
	}
	if n, _ := ttl.(int64); n <= 0 || !cache.GcIsNecessary() {
		return nil
	}
	maxSize, err := cfg.GetFlags("cache-files-maxsize", 0)
	if err != nil {
		return err
	}
	if io.IsVeryVerbose() {
		io.WriteError("Running cache garbage collection")
	}
	c := cache.New(cfg.GetString("cache-files-dir"), "a-z0-9_./", cfg.GetBool("cache-read-only"))
	c.Warn = func(s string) { io.WriteError(s) }
	m, _ := maxSize.(int64)
	c.GC(ttl.(int64), m)
	return nil
}
