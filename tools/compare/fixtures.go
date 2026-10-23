package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gocomposer/internal/phpjson"
)

const repoURL = "http://127.0.0.1:8899"

var fixedTime = time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

type fixturePkg struct {
	name, version, dir string
}

// buildRepo creates dists, packages.json (inline v1 format) and a v2
// metadata-url variant under out/.
func buildRepo(srcRoot, out string) error {
	var pkgs []fixturePkg
	vendors, err := os.ReadDir(srcRoot)
	if err != nil {
		return err
	}
	for _, v := range vendors {
		names, _ := os.ReadDir(filepath.Join(srcRoot, v.Name()))
		for _, n := range names {
			versions, _ := os.ReadDir(filepath.Join(srcRoot, v.Name(), n.Name()))
			for _, ver := range versions {
				pkgs = append(pkgs, fixturePkg{v.Name() + "/" + n.Name(), ver.Name(), filepath.Join(srcRoot, v.Name(), n.Name(), ver.Name())})
			}
		}
	}
	sort.Slice(pkgs, func(i, j int) bool {
		if pkgs[i].name != pkgs[j].name {
			return pkgs[i].name < pkgs[j].name
		}
		return pkgs[i].version < pkgs[j].version
	})
	if err := os.MkdirAll(filepath.Join(out, "dists"), 0o755); err != nil {
		return err
	}
	inline := phpjson.NewMap()
	byName := map[string][]phpjson.Value{}
	var names []string
	for _, p := range pkgs {
		raw, err := os.ReadFile(filepath.Join(p.dir, "composer.json"))
		if err != nil {
			return err
		}
		v, err := phpjson.Decode(raw)
		if err != nil {
			return fmt.Errorf("%s: %w", p.dir, err)
		}
		meta := phpjson.AsMap(v)
		meta.Set("version", p.version)
		sum := sha1.Sum([]byte(p.name + "@" + p.version))
		ref := hex.EncodeToString(sum[:])
		file := strings.ReplaceAll(p.name, "/", "-") + "-" + p.version + ".zip"
		zipPath := filepath.Join(out, "dists", file)
		top := strings.ReplaceAll(p.name, "/", "-") + "-" + ref[:7] + "/"
		if err := writeZip(zipPath, p.dir, top, p.name); err != nil {
			return err
		}
		shasum, err := fileSha1(zipPath)
		if err != nil {
			return err
		}
		distURL := repoURL + "/dists/" + file
		switch p.name {
		case "acme/badsum":
			shasum = strings.Repeat("0", 40)
		case "acme/missing":
			distURL = repoURL + "/dists/does-not-exist.zip"
		}
		dist := phpjson.NewMap()
		dist.Set("type", "zip")
		dist.Set("url", distURL)
		dist.Set("reference", ref)
		dist.Set("shasum", shasum)
		meta.Set("dist", dist)
		meta.Set("time", fixedTime.Format("2006-01-02T15:04:05-07:00"))
		if _, ok := byName[p.name]; !ok {
			names = append(names, p.name)
		}
		byName[p.name] = append(byName[p.name], meta)
	}
	pk := phpjson.NewMap()
	for _, n := range names {
		vm := phpjson.NewMap()
		for _, v := range byName[n] {
			vm.Set(phpjson.AsMap(v).Str("version"), v)
		}
		pk.Set(n, vm)
	}
	inline.Set("packages", pk)
	if err := writeJSON(filepath.Join(out, "packages.json"), inline); err != nil {
		return err
	}
	// v2 metadata-url variant (served from /v2)
	v2 := filepath.Join(out, "v2")
	if err := os.MkdirAll(filepath.Join(v2, "p2", "acme"), 0o755); err != nil {
		return err
	}
	root := phpjson.NewMap()
	root.Set("packages", []phpjson.Value{})
	root.Set("metadata-url", "/v2/p2/%package%.json")
	avail := []phpjson.Value{}
	for _, n := range names {
		avail = append(avail, n)
	}
	root.Set("available-packages", avail)
	if err := writeJSON(filepath.Join(v2, "packages.json"), root); err != nil {
		return err
	}
	for _, n := range names {
		var stable, dev []phpjson.Value
		for _, v := range byName[n] {
			m := phpjson.AsMap(v)
			ver := m.Str("version")
			if strings.HasPrefix(ver, "dev-") || strings.HasSuffix(ver, "-dev") {
				dev = append(dev, v)
			} else {
				stable = append(stable, v)
			}
		}
		for _, variant := range []struct {
			file string
			list []phpjson.Value
		}{{n + ".json", stable}, {n + "~dev.json", dev}} {
			if len(variant.list) == 0 {
				continue
			}
			// newest first like Packagist
			list := append([]phpjson.Value(nil), variant.list...)
			for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
				list[i], list[j] = list[j], list[i]
			}
			doc := phpjson.NewMap()
			p := phpjson.NewMap()
			p.Set(n, list)
			doc.Set("packages", p)
			if err := writeJSON(filepath.Join(v2, "p2", variant.file), doc); err != nil {
				return err
			}
		}
	}
	return buildAdvisoryVariant(out, names, byName)
}

// advisories served by the /v2adv variant through security-advisories
// metadata (complete for acme/http, partial for acme/util); acme/old is
// marked abandoned there.
var fixtureAdvisories = map[string][]map[string]any{
	"acme/http": {{
		"advisoryId": "PKSA-test-0001", "packageName": "acme/http", "remoteId": "GHSA-test-0001",
		"title": "Header injection in acme/http", "link": "https://example.org/advisories/1", "cve": "CVE-2024-0001",
		"affectedVersions": ">=1.2.0,<1.3.0", "source": "GitHub", "reportedAt": "2024-01-02 03:04:05",
		"composerRepository": "https://packagist.org", "severity": "high",
		"sources": []any{map[string]any{"name": "GitHub", "remoteId": "GHSA-test-0001"}},
	}},
	"acme/util": {{"advisoryId": "PKSA-test-0002", "packageName": "acme/util", "affectedVersions": "<1.1.0"}},
}

func toValue(v any) phpjson.Value {
	switch t := v.(type) {
	case map[string]any:
		m := phpjson.NewMap()
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			m.Set(k, toValue(t[k]))
		}
		return m
	case []any:
		out := []phpjson.Value{}
		for _, e := range t {
			out = append(out, toValue(e))
		}
		return out
	}
	return v
}

func buildAdvisoryVariant(out string, names []string, byName map[string][]phpjson.Value) error {
	dir := filepath.Join(out, "v2adv")
	if err := os.MkdirAll(filepath.Join(dir, "p2", "acme"), 0o755); err != nil {
		return err
	}
	root := phpjson.NewMap()
	root.Set("packages", []phpjson.Value{})
	root.Set("metadata-url", "/v2adv/p2/%package%.json")
	avail := []phpjson.Value{}
	for _, n := range names {
		avail = append(avail, n)
	}
	root.Set("available-packages", avail)
	sa := phpjson.NewMap()
	sa.Set("metadata", true)
	root.Set("security-advisories", sa)
	if err := writeJSON(filepath.Join(dir, "packages.json"), root); err != nil {
		return err
	}
	for _, n := range names {
		var list []phpjson.Value
		for i := len(byName[n]) - 1; i >= 0; i-- {
			v := phpjson.AsMap(phpjson.CloneValue(byName[n][i]))
			ver := v.Str("version")
			if strings.HasPrefix(ver, "dev-") || strings.HasSuffix(ver, "-dev") {
				continue
			}
			if n == "acme/old" {
				v.Set("abandoned", "acme/new")
			}
			list = append(list, v)
		}
		doc := phpjson.NewMap()
		p := phpjson.NewMap()
		if list == nil {
			list = []phpjson.Value{}
		}
		p.Set(n, list)
		doc.Set("packages", p)
		advs := []phpjson.Value{}
		for _, a := range fixtureAdvisories[n] {
			advs = append(advs, toValue(a))
		}
		doc.Set("security-advisories", advs)
		if err := writeJSON(filepath.Join(dir, "p2", n+".json"), doc); err != nil {
			return err
		}
	}
	return buildAdvisoryAPIVariant(out)
}

// buildAdvisoryAPIVariant writes /v2api: the /v2adv package metadata, but
// advisories only through the security-advisories "api-url" (a static
// response holding every fixtureAdvisories entry, served for the POST).
func buildAdvisoryAPIVariant(out string) error {
	dir := filepath.Join(out, "v2api")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	adv, err := phpjson.Decode(mustReadFile(filepath.Join(out, "v2adv", "packages.json")))
	if err != nil {
		return err
	}
	root := phpjson.AsMap(adv)
	sa := phpjson.NewMap()
	sa.Set("metadata", false)
	sa.Set("api-url", "/v2api/advisories.json")
	root.Set("security-advisories", sa)
	if err := writeJSON(filepath.Join(dir, "packages.json"), root); err != nil {
		return err
	}
	names := make([]string, 0, len(fixtureAdvisories))
	for n := range fixtureAdvisories {
		names = append(names, n)
	}
	sort.Strings(names)
	all := phpjson.NewMap()
	for _, n := range names {
		list := []phpjson.Value{}
		for _, a := range fixtureAdvisories[n] {
			list = append(list, toValue(a))
		}
		all.Set(n, list)
	}
	resp := phpjson.NewMap()
	resp.Set("advisories", all)
	return writeJSON(filepath.Join(dir, "advisories.json"), resp)
}

func mustReadFile(p string) []byte {
	b, _ := os.ReadFile(p)
	return b
}

func writeJSON(path string, v phpjson.Value) error {
	s, err := phpjson.Encode(v, phpjson.ComposerFlags)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(s+"\n"), 0o644)
}

func fileSha1(p string) (string, error) {
	f, err := os.Open(p)
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

func writeZip(zipPath, srcDir, top, name string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(f)
	var files []string
	_ = filepath.Walk(srcDir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	sort.Strings(files)
	addDir := func(dir string) error {
		h := &zip.FileHeader{Name: dir, Method: zip.Store, Modified: fixedTime}
		h.SetMode(os.ModeDir | 0o755)
		_, err := zw.CreateHeader(h)
		return err
	}
	if err := addDir(top); err != nil {
		return err
	}
	dirs := map[string]bool{}
	for _, p := range files {
		rel, _ := filepath.Rel(srcDir, p)
		rel = filepath.ToSlash(rel)
		parts := strings.Split(rel, "/")
		for i := 1; i < len(parts); i++ {
			d := top + strings.Join(parts[:i], "/") + "/"
			if !dirs[d] {
				dirs[d] = true
				if err := addDir(d); err != nil {
					return err
				}
			}
		}
		st, _ := os.Stat(p)
		h := &zip.FileHeader{Name: top + rel, Method: zip.Deflate, Modified: fixedTime}
		h.SetMode(st.Mode().Perm())
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		data, _ := os.ReadFile(p)
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	if name == "acme/evil" {
		h := &zip.FileHeader{Name: "../../evil-escaped.txt", Method: zip.Store, Modified: fixedTime}
		h.SetMode(0o644)
		w, _ := zw.CreateHeader(h)
		w.Write([]byte("should never be written\n"))
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if name == "acme/broken" {
		data, _ := os.ReadFile(zipPath)
		return os.WriteFile(zipPath, data[:len(data)/2], 0o644)
	}
	return nil
}
