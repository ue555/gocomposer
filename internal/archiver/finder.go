// Package archiver ports Composer\Package\Archiver (ArchiveManager,
// ArchivableFilesFinder/Filter, BaseExcludeFilter, GitExcludeFilter,
// ComposerExcludeFilter and PharArchiver) of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT, (c) Nils Adermann,
// Jordi Boggiano), together with the parts of symfony/finder v5.4 (Glob,
// the VCS directory exclusion and sorting by name; MIT, (c) Fabien
// Potencier) that the finder relies on.
//
// Archives are written in the byte layout of PHP's PharData (ext/phar of
// PHP 8.4), the archiver Composer uses when ext-zip is not loaded (as in
// the official PHP image used as the comparison oracle).
package archiver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/pcre"
)

// vcsPatterns is Finder::$vcsPatterns (ignoreVCS(true)).
var vcsPatterns = map[string]bool{".svn": true, "_svn": true, "CVS": true, "_darcs": true, ".arch-params": true, ".monotone": true, ".bzr": true, ".git": true, ".hg": true}

// Item is one file or (empty) directory found by FindFiles.
type Item struct {
	// Pathname is the path below the sources directory as iterated.
	Pathname string
	// RealPath is the resolved path (symlinks followed).
	RealPath string
	IsDir    bool
	IsLink   bool
}

type excludePattern struct {
	pattern string
	negate  bool
}

// filter is BaseExcludeFilter.
type filter struct{ patterns []excludePattern }

func (f *filter) apply(relativePath string, exclude bool) (bool, error) {
	for _, p := range f.patterns {
		m, err := pcre.IsMatch(p.pattern, relativePath)
		var ce *pcre.CompileError
		if errors.As(err, &ce) && !ce.Unsupported {
			// the preg_match() compile warning becomes an ErrorException,
			// which is not the RuntimeException the filter suppresses
			return exclude, &PatternError{Msg: compileFailure(ce)}
		}
		// PcreException (backtrack limits) is suppressed
		if err == nil && m != nil {
			exclude = !p.negate
		}
	}
	return exclude, nil
}

// PatternError is the ErrorException of an invalid exclude pattern.
type PatternError struct{ Msg string }

func (e *PatternError) Error() string { return e.Msg }

// compileFailure renders PHP's "preg_match(): Compilation failed" warning.
// PCRE2 reports unterminated classes and groups at the end of the pattern.
func compileFailure(ce *pcre.CompileError) string {
	msg := "preg_match(): Compilation failed: " + ce.Msg
	if strings.HasPrefix(ce.Msg, "missing terminating ]") || strings.HasPrefix(ce.Msg, "missing closing parenthesis") {
		body := ce.Pattern
		if len(body) >= 2 {
			body = body[1:strings.LastIndexAny(body, "}")]
		}
		msg += fmt.Sprintf(" at offset %d", len(body))
	}
	return msg
}

// generatePattern is BaseExcludeFilter::generatePattern.
func generatePattern(rule string) excludePattern {
	negate := false
	pattern := ""
	if rule != "" && rule[0] == '!' {
		negate = true
		rule = strings.TrimLeft(rule, "!")
	}
	first := strings.IndexByte(rule, '/')
	if first == 0 {
		pattern = "^/"
	} else if first < 0 || len(rule)-1 == first {
		pattern = "/"
	}
	rule = strings.Trim(rule, "/")
	re := GlobToRegex(rule, true, true, "#")
	// remove delimiters as well as caret (^) and dollar sign ($) from the regex
	re = re[2 : len(re)-2]
	return excludePattern{pattern: "{" + pattern + re + "(?=$|/)}", negate: negate}
}

// parseLines is BaseExcludeFilter::parseLines.
func parseLines(lines []string, parse func(string) *excludePattern) []excludePattern {
	var out []excludePattern
	for _, line := range lines {
		line = strings.Trim(line, " \t\n\r\x00\x0B")
		if line == "" || line == "0" || strings.HasPrefix(line, "#") {
			continue
		}
		if p := parse(line); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// newGitExcludeFilter is GitExcludeFilter: export-ignore rules of the
// .gitattributes file of the source directory.
func newGitExcludeFilter(sourcePath string) *filter {
	f := &filter{}
	raw, err := os.ReadFile(sourcePath + "/.gitattributes")
	if err != nil {
		return f
	}
	f.patterns = parseLines(phpFileLines(string(raw)), func(line string) *excludePattern {
		parts := splitWhitespace(line)
		if len(parts) == 2 && parts[1] == "export-ignore" {
			p := generatePattern(parts[0])
			return &p
		}
		if len(parts) == 2 && parts[1] == "-export-ignore" {
			p := generatePattern("!" + parts[0])
			return &p
		}
		return nil
	})
	return f
}

// newComposerExcludeFilter is ComposerExcludeFilter.
func newComposerExcludeFilter(rules []string) *filter {
	f := &filter{}
	for _, r := range rules {
		f.patterns = append(f.patterns, generatePattern(r))
	}
	return f
}

// phpFileLines is file(): lines keep their newline.
func phpFileLines(s string) []string {
	var out []string
	for s != "" {
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			out = append(out, s)
			break
		}
		out = append(out, s[:i+1])
		s = s[i+1:]
	}
	return out
}

// splitWhitespace is Preg::split('#\s+#', $line).
func splitWhitespace(s string) []string {
	var parts []string
	cur := strings.Builder{}
	inSpace := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			if !inSpace {
				parts = append(parts, cur.String())
				cur.Reset()
				inSpace = true
			}
			continue
		}
		inSpace = false
		cur.WriteByte(c)
	}
	return append(parts, cur.String())
}

// SourceError is the RuntimeException for an unresolvable source directory.
type SourceError struct{ Msg string }

func (e *SourceError) Error() string { return e.Msg }

// FindFiles is ArchivableFilesFinder: the files (and empty directories) of
// sources that are not excluded, sorted by (real) path name.
func FindFiles(sources string, excludes []string, ignoreFilters bool) ([]Item, string, error) {
	real := fsutil.Realpath(sources)
	if real == "" {
		return nil, "", &SourceError{Msg: `Could not realpath() the source directory "` + sources + `"`}
	}
	sources = fsutil.NormalizePath(real)
	var filters []*filter
	if !ignoreFilters {
		filters = []*filter{newGitExcludeFilter(sources), newComposerExcludeFilter(excludes)}
	}
	var filterErr error
	accept := func(it Item) bool {
		if it.RealPath == "" || filterErr != nil {
			return false
		}
		if it.IsLink && !strings.HasPrefix(it.RealPath, sources) {
			return false
		}
		rel := fsutil.NormalizePath(it.RealPath)
		if strings.HasPrefix(rel, sources) {
			rel = rel[len(sources):]
		}
		exclude := false
		for _, f := range filters {
			var err error
			if exclude, err = f.apply(rel, exclude); err != nil {
				filterErr = err
				return false
			}
		}
		return !exclude
	}
	var items []Item
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			p := dir + "/" + e.Name()
			lst, err := os.Lstat(p)
			if err != nil {
				continue
			}
			isLink := lst.Mode()&os.ModeSymlink != 0
			isDir := lst.IsDir()
			if isLink {
				if st, err := os.Stat(p); err == nil {
					isDir = st.IsDir()
				}
			}
			// ExcludeDirectoryFilterIterator: VCS directories are skipped
			// together with their contents
			if isDir && vcsPatterns[e.Name()] {
				continue
			}
			it := Item{Pathname: p, IsDir: isDir, IsLink: isLink}
			if r, err := filepath.EvalSymlinks(p); err == nil {
				if abs, err := filepath.Abs(r); err == nil {
					it.RealPath = abs
				}
			}
			if accept(it) {
				items = append(items, it)
			}
			// RecursiveDirectoryIterator does not follow directory links
			if isDir && !isLink {
				walk(p)
			}
		}
	}
	walk(sources)
	if filterErr != nil {
		return nil, "", filterErr
	}
	// SortableIterator::SORT_BY_NAME
	key := func(it Item) string {
		if it.RealPath != "" {
			return it.RealPath
		}
		return it.Pathname
	}
	sort.SliceStable(items, func(i, j int) bool { return key(items[i]) < key(items[j]) })
	// ArchivableFilesFinder::accept: directories only when empty
	var out []Item
	for _, it := range items {
		if it.IsDir {
			entries, err := os.ReadDir(it.Pathname)
			if err != nil || len(entries) > 0 {
				continue
			}
		}
		out = append(out, it)
	}
	return out, sources, nil
}
