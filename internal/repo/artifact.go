package repo

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/ArtifactRepository.php, src/Composer/Util/Zip.php,
// src/Composer/Util/Tar.php). Copyright (c) Nils Adermann, Jordi Boggiano.
// MIT License.

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/config"
	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// Artifact is ArtifactRepository.
type Artifact struct {
	*Array
	lookup     string
	repoConfig *phpjson.Map
	io         *cio.IO
}

// NewArtifact is ArtifactRepository::__construct.
func NewArtifact(repoConfig *phpjson.Map, io *cio.IO) (*Artifact, error) {
	a := &Artifact{Array: &Array{kind: KindArtifact}, lookup: config.ExpandPath(repoConfig.Str("url")), repoConfig: repoConfig, io: io}
	a.Array.self = a
	a.Array.init = func() error { return a.scanDirectory(a.lookup) }
	return a, nil
}

// RepoName is ArtifactRepository::getRepoName.
func (a *Artifact) RepoName() string { return "artifact repo (" + a.lookup + ")" }

// RepoConfig is getRepoConfig.
func (a *Artifact) RepoConfig() *phpjson.Map { return a.repoConfig }

var reArtifactFile = regexp.MustCompile(`(?i)^.+\.(zip|tar|gz|tgz)$`)

// readDirRaw lists a directory in the order the filesystem returns entries
// (like RecursiveDirectoryIterator, which does not sort).
func readDirRaw(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(-1)
}

// scanDirectory walks path recursively (following symlinks) like
// RecursiveIteratorIterator(RecursiveDirectoryIterator(FOLLOW_SYMLINKS)).
func (a *Artifact) scanDirectory(root string) error {
	given := root
	root = strings.TrimRight(root, "/")
	if root == "" {
		root = "/"
	}
	if _, err := readDirRaw(root); err != nil {
		reason := "No such file or directory"
		if os.IsPermission(err) {
			reason = "Permission denied"
		} else if st, serr := os.Stat(root); serr == nil && !st.IsDir() {
			reason = "Not a directory"
		}
		return &RuntimeError{Msg: "RecursiveDirectoryIterator::__construct(" + given + "): Failed to open directory: " + reason}
	}
	var walk func(dir string) error
	walk = func(dir string) error {
		names, err := readDirRaw(dir)
		if err != nil {
			return nil
		}
		for _, n := range names {
			p := dir + "/" + n
			if dir == "/" {
				p = "/" + n
			}
			st, err := os.Stat(p)
			if err != nil {
				continue
			}
			if st.IsDir() {
				if err := walk(p); err != nil {
					return err
				}
				continue
			}
			if !reArtifactFile.MatchString(p) || !st.Mode().IsRegular() {
				continue
			}
			p2, err := a.composerInformation(p)
			if err != nil {
				return err
			}
			if p2 == nil {
				a.io.WriteErrorV("File <comment>"+n+"</comment> doesn't seem to hold a package", true, cio.Verbose)
				continue
			}
			a.io.WriteErrorV(fmt.Sprintf("Found package <info>%s</info> (<comment>%s</comment>) in file <info>%s</info>", p2.Name, p2.PrettyVersion(), n), true, cio.Verbose)
			a.Array.AddPackage(p2)
		}
		return nil
	}
	return walk(root)
}

// RuntimeError is a PHP \RuntimeException / \UnexpectedValueException raised
// by a repository.
type RuntimeError struct {
	Msg  string
	Prev error
}

func (e *RuntimeError) Error() string { return e.Msg }

// Previous is the chained (PHP previous) exception, rendered as a further box.
func (e *RuntimeError) Previous() error { return e.Prev }

func (a *Artifact) composerInformation(file string) (*pkg.Package, error) {
	ext := strings.TrimPrefix(path.Ext(file), ".")
	var json []byte
	var err error
	switch ext {
	case "gz", "tar", "tgz":
		json, err = tarComposerJSON(file)
	case "zip":
		json, err = zipComposerJSON(file)
	default:
		return nil, &RuntimeError{Msg: `Files with "` + ext + `" extensions aren't supported. Only ZIP and TAR/TAR.GZ/TGZ archives are supported.`}
	}
	if err != nil {
		a.io.WriteV("Failed loading package "+file+": "+err.Error(), false, cio.Verbose)
		json = nil
	}
	if json == nil {
		return nil, nil
	}
	v, err := jsonfile.Parse(json, file+"#composer.json")
	if err != nil {
		return nil, err
	}
	m := phpjson.ToMap(v)
	if m == nil {
		m = phpjson.NewMap()
	}
	m.Delete("published-time")
	sum, err := fileSha1(file)
	if err != nil {
		return nil, err
	}
	dist := phpjson.NewMap()
	typ := "zip"
	if ext != "zip" {
		typ = "tar"
	}
	dist.Set("type", typ)
	dist.Set("url", strings.ReplaceAll(file, `\`, "/"))
	dist.Set("shasum", sum)
	m.Set("dist", dist)
	p, err := (&pkg.Loader{}).Load(m)
	if err != nil {
		return nil, &RuntimeError{Msg: "Failed loading package in " + file + ": " + err.Error(), Prev: err}
	}
	return p, nil
}

func fileSha1(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

const multipleTopLevel = "Archive has more than one top level directories, and no composer.json was found on the top level, so it's an invalid archive. Top level paths found were: "

// zipComposerJSON is Zip::getComposerJson (nil when the file is not a
// readable archive).
func zipComposerJSON(file string) ([]byte, error) {
	zr, err := zip.OpenReader(file)
	if err != nil {
		return nil, nil
	}
	defer zr.Close()
	if len(zr.File) == 0 {
		return nil, nil
	}
	read := func(name string) ([]byte, bool) {
		for _, f := range zr.File {
			if f.Name == name {
				rc, err := f.Open()
				if err != nil {
					return nil, false
				}
				defer rc.Close()
				b, err := io.ReadAll(rc)
				if err != nil {
					return nil, false
				}
				return b, true
			}
		}
		return nil, false
	}
	if b, ok := read("composer.json"); ok {
		return b, nil
	}
	var top []string
	seen := map[string]bool{}
	addTop := func(p string) error {
		if !seen[p] {
			seen[p] = true
			top = append(top, p)
		}
		if len(top) > 1 {
			return &RuntimeError{Msg: multipleTopLevel + strings.Join(top, ",")}
		}
		return nil
	}
	for _, f := range zr.File {
		name := f.Name
		dir := phpDirname(name)
		if strings.Contains(name, "__MACOSX") {
			continue
		}
		if dir == "." {
			if err := addTop(name); err != nil {
				return nil, err
			}
			continue
		}
		if !strings.Contains(dir, `\`) && !strings.Contains(dir, "/") {
			if err := addTop(dir + "/"); err != nil {
				return nil, err
			}
		}
	}
	if len(top) > 0 {
		if b, ok := read(top[0] + "composer.json"); ok {
			return b, nil
		}
	}
	return nil, &RuntimeError{Msg: "No composer.json found either at the top level or within the topmost directory"}
}

// phpDirname is PHP dirname() for archive entry names.
func phpDirname(p string) string {
	t := strings.TrimRight(p, "/")
	if t == "" {
		return "/"
	}
	i := strings.LastIndexByte(t, '/')
	if i < 0 {
		return "."
	}
	d := strings.TrimRight(t[:i], "/")
	if d == "" {
		return "/"
	}
	return d
}

// tarComposerJSON is Tar::getComposerJson (PharData over tar, tar.gz and
// tar.bz2 archives).
func tarComposerJSON(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// PharData reports the archive by its absolute (cwd expanded) path
	phar := file
	if abs, err := filepath.Abs(file); err == nil {
		phar = abs
	}
	br := bufio.NewReader(f)
	var r io.Reader = br
	if magic, _ := br.Peek(3); len(magic) >= 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		gz, err := gzip.NewReader(br)
		if err != nil {
			return nil, &RuntimeError{Msg: "internal corruption of phar \"" + phar + "\" (__HALT_COMPILER(); not found)"}
		}
		r = gz
	} else if len(magic) == 3 && string(magic) == "BZh" {
		r = bzip2.NewReader(br)
	}
	type entry struct {
		dir     bool
		content []byte
	}
	entries := map[string]*entry{}
	var order []string
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			if len(order) == 0 {
				return nil, &RuntimeError{Msg: "internal corruption of phar \"" + phar + "\" (__HALT_COMPILER(); not found)"}
			}
			break
		}
		name := strings.TrimPrefix(filepath.ToSlash(h.Name), "./")
		isDir := h.Typeflag == tar.TypeDir || strings.HasSuffix(name, "/")
		name = strings.TrimRight(name, "/")
		if name == "" {
			continue
		}
		e := &entry{dir: isDir}
		if !isDir && (name == "composer.json" || strings.Count(name, "/") == 1 && strings.HasSuffix(name, "/composer.json")) {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tr); err == nil {
				e.content = buf.Bytes()
			}
		}
		if _, ok := entries[name]; !ok {
			order = append(order, name)
		}
		entries[name] = e
		// implicit parent directories
		if i := strings.IndexByte(name, '/'); i > 0 {
			parent := name[:i]
			if _, ok := entries[parent]; !ok {
				entries[parent] = &entry{dir: true}
				order = append(order, parent)
			}
		}
	}
	if e, ok := entries["composer.json"]; ok && !e.dir {
		return e.content, nil
	}
	var top []string
	for _, n := range order {
		if strings.Contains(n, "/") {
			continue
		}
		if entries[n].dir {
			top = append(top, n)
			if len(top) > 1 {
				return nil, &RuntimeError{Msg: multipleTopLevel + strings.Join(top, ",")}
			}
		}
	}
	if len(top) > 0 {
		if e, ok := entries[top[0]+"/composer.json"]; ok && !e.dir {
			return e.content, nil
		}
	}
	return nil, &RuntimeError{Msg: "No composer.json found either at the top level or within the topmost directory"}
}
