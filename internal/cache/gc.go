package cache

// Port of Composer\Cache::gc, ::gcVcsCache and ::gcIsNecessary
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	collectedMu sync.Mutex
	collected   bool
)

// GcIsNecessary is Cache::gcIsNecessary: at most once per process, never
// under COMPOSER_TEST_SUITE, otherwise with probability 1/51
// (!random_int(0, 50)).
func GcIsNecessary() bool {
	collectedMu.Lock()
	defer collectedMu.Unlock()
	if collected {
		return false
	}
	collected = true
	if os.Getenv("COMPOSER_TEST_SUITE") != "" {
		return false
	}
	n, err := rand.Int(rand.Reader, big.NewInt(51))
	return err == nil && n.Int64() == 0
}

// GC is Cache::gc for this cache (no-op when disabled or read-only).
func (c *Cache) GC(ttl, maxSize int64) bool {
	if !c.IsEnabled() || c.ReadOnly {
		return false
	}
	GC(c.root, ttl, maxSize)
	collectedMu.Lock()
	collected = true
	collectedMu.Unlock()
	return true
}

// GC removes files below root older than ttl seconds, then the least
// recently accessed files until the total size is at most maxSize.
// Finder's date('until X') keeps files whose mtime is not strictly before X
// (seconds precision).
func GC(root string, ttl, maxSize int64) {
	expire := time.Now().Add(-time.Duration(ttl) * time.Second).Truncate(time.Second)
	for _, f := range FinderFiles(root) {
		if st, err := os.Stat(f); err == nil && st.ModTime().Unix() < expire.Unix() {
			_ = os.Remove(f)
		}
	}
	total := DirSize(root)
	if total > maxSize {
		files := FinderFiles(root)
		atimes := map[string]int64{}
		for _, f := range files {
			atimes[f] = fileAtime(f)
		}
		sort.SliceStable(files, func(i, j int) bool { return atimes[files[i]] < atimes[files[j]] })
		for _, f := range files {
			if total <= maxSize {
				break
			}
			if st, err := os.Stat(f); err == nil {
				total -= st.Size()
			}
			_ = os.Remove(f)
		}
	}
}

// GCVcs is Cache::gcVcsCache.
func GCVcs(root string, ttl int64) {
	expire := time.Now().Add(-time.Duration(ttl) * time.Second).Truncate(time.Second)
	entries, _ := os.ReadDir(root)
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		p := filepath.Join(root, e.Name())
		if st, err := os.Stat(p); err == nil && st.ModTime().Unix() < expire.Unix() {
			_ = os.RemoveAll(p)
		}
	}
}

// FinderFiles lists files below root like Finder::create()->in()->files()
// (dot files and dot directories ignored).
func FinderFiles(root string) []string {
	var out []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if p == filepath.Clean(root) {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode().IsRegular() {
			out = append(out, p)
		}
		return nil
	})
	return out
}

// DirSize is Filesystem::size for a directory.
func DirSize(root string) int64 {
	var total int64
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && info.Mode().IsRegular() {
			total += info.Size()
		}
		return nil
	})
	return total
}

// fileAtime approximates SplFileInfo::getATime with the modification time
// (portable; caches are written once, so both normally agree).
func fileAtime(p string) int64 {
	st, err := os.Stat(p)
	if err != nil {
		return 0
	}
	return st.ModTime().UnixNano()
}
