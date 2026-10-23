package archiver

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Formats is PharArchiver::$formats (the formats ArchiveManager supports).
var Formats = []string{"zip", "tar", "tar.gz", "tar.bz2"}

// Supports is PharArchiver::supports.
func Supports(format string) bool {
	for _, f := range Formats {
		if f == format {
			return true
		}
	}
	return false
}

// ArchiveError is a RuntimeException raised while archiving.
type ArchiveError struct{ Msg string }

func (e *ArchiveError) Error() string { return e.Msg }

type entry struct {
	name    string
	dir     bool
	mode    uint32
	content []byte
}

// now is the entry time (PharData stamps entries with the current time).
var now = time.Now

// Archive is PharArchiver::archive: it writes the archivable files of
// sources to target (whose name ends with "."+format) and returns the
// path of the written archive.
func Archive(sources, target, format string, excludes []string, ignoreFilters bool) (string, error) {
	if real := realOrEmpty(sources); real != "" {
		sources = real
	}
	if _, err := os.Stat(target); err == nil {
		os.Remove(target)
	}
	filename := target[:strings.LastIndex(target, format)-1]
	if format == "tar.gz" || format == "tar.bz2" {
		target = filename + ".tar"
	}
	items, base, err := FindFiles(sources, excludes, ignoreFilters)
	if err != nil {
		return "", err
	}
	// PharData::buildFromIterator then ArchivableFilesFilter::addEmptyDir
	var entries []*entry
	byName := map[string]*entry{}
	var dirs []Item
	for _, it := range items {
		if it.IsDir {
			dirs = append(dirs, it)
			continue
		}
		// phar stores files under their resolved path
		local := strings.TrimPrefix(strings.TrimPrefix(it.RealPath, base), "/")
		data, err := os.ReadFile(it.Pathname)
		if err != nil {
			return "", &ArchiveError{Msg: fmt.Sprintf("Could not create archive '%s' from '%s': %s", target, sources, err.Error())}
		}
		var mode uint32 = 0o644
		if st, err := os.Stat(it.Pathname); err == nil {
			mode = uint32(st.Mode().Perm())
		}
		if e, ok := byName[local]; ok {
			e.content, e.mode = data, mode
			continue
		}
		e := &entry{name: local, mode: mode, content: data}
		byName[local] = e
		entries = append(entries, e)
	}
	for _, d := range dirs {
		local := strings.ReplaceAll(d.Pathname, base+"/", "")
		if _, ok := byName[local+"/"]; ok {
			continue
		}
		e := &entry{name: local, dir: true, mode: 0o777}
		byName[local+"/"] = e
		entries = append(entries, e)
	}
	if len(entries) == 0 {
		// PharData did not write anything: Composer writes empty archives
		target = filename + "." + format
		var data []byte
		switch format {
		case "tar":
			data = make([]byte, 10240)
		case "zip":
			data = make([]byte, 22)
			binary.LittleEndian.PutUint32(data, 0x06054b50)
		case "tar.gz":
			// gzcompress(): zlib format
			var buf bytes.Buffer
			w := zlib.NewWriter(&buf)
			_, _ = w.Write(make([]byte, 10240))
			_ = w.Close()
			data = buf.Bytes()
		case "tar.bz2":
			return "", &ArchiveError{Msg: "Can not compress to tar.bz2 format"}
		}
		if err := os.WriteFile(target, data, 0o666); err != nil {
			return "", err
		}
		return target, nil
	}
	var data []byte
	if format == "zip" {
		data = buildZip(entries)
	} else {
		data, err = buildTar(entries, target)
		if err != nil {
			return "", err
		}
	}
	if format == "tar.bz2" {
		// PharData::canCompress(Phar::BZ2) is false without ext-bz2 (as in
		// the oracle PHP); gocomposer has no bzip2 compressor either
		return "", &ArchiveError{Msg: "Can not compress to tar.bz2 format"}
	}
	if format == "tar.gz" {
		var buf bytes.Buffer
		w, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
		w.OS = 3
		_, _ = w.Write(data)
		_ = w.Close()
		data = buf.Bytes()
		target = filename + "." + format
	}
	if err := os.WriteFile(target, data, 0o666); err != nil {
		return "", err
	}
	return target, nil
}

func realOrEmpty(p string) string {
	r, err := filepath.EvalSymlinks(p)
	if err != nil {
		return ""
	}
	abs, err := filepath.Abs(r)
	if err != nil {
		return ""
	}
	return abs
}

func octal(dst []byte, v int64) {
	s := fmt.Sprintf("%0*o", len(dst), v)
	copy(dst, s[len(s)-len(dst):])
}

// buildTar writes the entries like ext/phar's tar writer (ustar headers,
// zero uid/gid fields, terminated by two zero blocks).
func buildTar(entries []*entry, target string) ([]byte, error) {
	var buf bytes.Buffer
	mtime := now().Unix()
	for _, e := range entries {
		var h [512]byte
		name := e.name
		if len(name) > 100 {
			if len(name) > 256 {
				return nil, &ArchiveError{Msg: fmt.Sprintf(`tar-based phar "%s" cannot be created, filename "%s" is too long for tar file format`, target, name)}
			}
			b := len(name) - 101
			for b < len(name) && name[b] != '/' {
				b++
			}
			if b >= len(name) || b > 155 {
				return nil, &ArchiveError{Msg: fmt.Sprintf(`tar-based phar "%s" cannot be created, filename "%s" is too long for tar file format`, target, name)}
			}
			copy(h[345:500], name[:b])
			copy(h[0:100], name[b+1:])
		} else {
			copy(h[0:100], name)
		}
		octal(h[100:107], int64(e.mode))
		size := int64(len(e.content))
		if e.dir {
			size = 0
		}
		octal(h[124:135], size)
		octal(h[136:147], mtime)
		if e.dir {
			h[156] = '5'
		} else {
			h[156] = '0'
		}
		copy(h[257:263], "ustar\x00")
		copy(h[263:265], "00")
		for i := 148; i < 156; i++ {
			h[i] = ' '
		}
		sum := 0
		for _, c := range h {
			sum += int(c)
		}
		octal(h[148:155], int64(sum))
		buf.Write(h[:])
		if !e.dir {
			buf.Write(e.content)
			if pad := len(e.content) % 512; pad != 0 {
				buf.Write(make([]byte, 512-pad))
			}
		}
	}
	buf.Write(make([]byte, 1024))
	return buf.Bytes(), nil
}

func dosDateTime(t time.Time) (uint16, uint16) {
	if t.Year() < 1980 {
		return 0, 0x21
	}
	tm := uint16(t.Hour()<<11 | t.Minute()<<5 | t.Second()/2)
	dt := uint16((t.Year()-1980)<<9 | int(t.Month())<<5 | t.Day())
	return tm, dt
}

// buildZip writes the entries like ext/phar's zip writer (stored entries,
// UTF-8 flag, an "nu" (ASi Unix) extra field with the permissions).
func buildZip(entries []*entry) []byte {
	var buf, central bytes.Buffer
	tm, dt := dosDateTime(now())
	le := binary.LittleEndian
	for _, e := range entries {
		name := e.name
		var content []byte
		if e.dir {
			name += "/"
		} else {
			content = e.content
		}
		crc := crc32.ChecksumIEEE(content)
		var asi [10]byte
		le.PutUint16(asi[0:], uint16(e.mode))
		extra := make([]byte, 18)
		le.PutUint16(extra[0:], 0x756e)
		le.PutUint16(extra[2:], 14)
		// ext/phar's CRC covers the permission bytes only
		le.PutUint32(extra[4:], crc32.ChecksumIEEE(asi[0:2]))
		copy(extra[8:], asi[:])
		offset := buf.Len()
		var lh [30]byte
		le.PutUint32(lh[0:], 0x04034b50)
		le.PutUint16(lh[4:], 0)
		le.PutUint16(lh[6:], 0x0800)
		le.PutUint16(lh[8:], 0)
		le.PutUint16(lh[10:], tm)
		le.PutUint16(lh[12:], dt)
		le.PutUint32(lh[14:], crc)
		le.PutUint32(lh[18:], uint32(len(content)))
		le.PutUint32(lh[22:], uint32(len(content)))
		le.PutUint16(lh[26:], uint16(len(name)))
		le.PutUint16(lh[28:], uint16(len(extra)))
		buf.Write(lh[:])
		buf.WriteString(name)
		buf.Write(extra)
		buf.Write(content)
		var ch [46]byte
		le.PutUint32(ch[0:], 0x02014b50)
		le.PutUint16(ch[4:], 0)
		le.PutUint16(ch[6:], 0)
		le.PutUint16(ch[8:], 0x0800)
		le.PutUint16(ch[10:], 0)
		le.PutUint16(ch[12:], tm)
		le.PutUint16(ch[14:], dt)
		le.PutUint32(ch[16:], crc)
		le.PutUint32(ch[20:], uint32(len(content)))
		le.PutUint32(ch[24:], uint32(len(content)))
		le.PutUint16(ch[28:], uint16(len(name)))
		le.PutUint16(ch[30:], uint16(len(extra)))
		le.PutUint32(ch[42:], uint32(offset))
		central.Write(ch[:])
		central.WriteString(name)
		central.Write(extra)
	}
	cdOffset := buf.Len()
	buf.Write(central.Bytes())
	var eocd [22]byte
	le.PutUint32(eocd[0:], 0x06054b50)
	le.PutUint16(eocd[8:], uint16(len(entries)))
	le.PutUint16(eocd[10:], uint16(len(entries)))
	le.PutUint32(eocd[12:], uint32(central.Len()))
	le.PutUint32(eocd[16:], uint32(cdOffset))
	buf.Write(eocd[:])
	return buf.Bytes()
}
