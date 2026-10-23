// The parts of Composer\Json\JsonFile (read, write, encode, validateSchema)
// and Composer\Util\Filesystem::safeFilePutContents that JsonConfigSource
// relies on. Ported from src/Composer/Json/JsonFile.php and
// src/Composer/Util/Filesystem.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package configsource

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/jsonschema"
	"gocomposer/internal/pcre"
	"gocomposer/internal/phpjson"
)

const indentDefault = "    "

// jsonFile is a minimal Composer\Json\JsonFile for a local path.
type jsonFile struct {
	path   string
	indent string
}

// exists is JsonFile::exists() (is_file).
func (f *jsonFile) exists() bool {
	st, err := os.Stat(f.path)
	return err == nil && st.Mode().IsRegular()
}

// read is JsonFile::read(): it remembers the indentation and returns the
// decoded (associative) data.
func (f *jsonFile) read() phpjson.Value {
	if !isReadable(f.path) {
		phpError("RuntimeException", "Could not read "+f.path+"\n\nThe file \""+f.path+"\" is not readable.")
	}
	data, err := os.ReadFile(f.path)
	if err != nil {
		phpError("RuntimeException", "Could not read "+f.path)
	}
	f.indent = jsonmanip.DetectIndenting(string(data))
	v, err := jsonfile.Parse(data, f.path)
	if err != nil {
		class := "Seld\\JsonLint\\ParsingException"
		var ee *jsonfile.EncodingError
		if errors.As(err, &ee) {
			class = "UnexpectedValueException"
		}
		panic(&jsonmanip.Exception{Class: class, Message: err.Error(), Err: err})
	}
	return v
}

// write is JsonFile::write($hash) with the default options
// (JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE).
func (f *jsonFile) write(hash phpjson.Value) error {
	dir := filepath.Dir(f.path)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		if err == nil {
			real, rerr := filepath.EvalSymlinks(dir)
			if rerr != nil {
				real = dir
			} else if abs, aerr := filepath.Abs(real); aerr == nil {
				real = abs
			}
			phpError("UnexpectedValueException", real+" exists and is not a directory.")
		}
		if err := os.MkdirAll(dir, 0o777); err != nil {
			phpError("UnexpectedValueException", dir+" does not exist and could not be created.")
		}
	}
	indent := f.indent
	if indent == "" {
		indent = indentDefault
	}
	content := encode(hash, indent) + "\n"
	// filePutContentsIfModified
	if cur, err := os.ReadFile(f.path); err == nil && string(cur) == content {
		return nil
	}
	return safeFilePutContents(f.path, content)
}

// encode is JsonFile::encode($data, default flags, $indent).
func encode(data phpjson.Value, indent string) string {
	var b strings.Builder
	encodeValue(&b, data, 0)
	json := b.String()
	if indent != indentDefault {
		// Pretty printing and not using default indentation
		out, _, err := pcre.ReplaceCallback(`#^ {4,}#m`, func(m *pcre.Match) string {
			return strings.Repeat(indent, len(m.Str(0))/4)
		}, json, -1)
		if err != nil {
			phpError("Composer\\Pcre\\PcreException", "preg_replace_callback(): failed executing \"#^ {4,}#m\": "+err.Error())
		}
		return out
	}
	return json
}

func encodeValue(b *strings.Builder, v phpjson.Value, depth int) {
	nl := func(d int) {
		b.WriteByte('\n')
		for i := 0; i < d; i++ {
			b.WriteString(indentDefault)
		}
	}
	switch t := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if t {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case int64:
		b.WriteString(strconv.FormatInt(t, 10))
	case int:
		b.WriteString(strconv.Itoa(t))
	case float64:
		b.WriteString(jsonmanip.FormatFloat(t))
	case string:
		s, err := phpjson.Encode(t, phpjson.ComposerFlags)
		if err != nil {
			phpError("RuntimeException", "JSON encoding failed: Malformed UTF-8 characters, possibly incorrectly encoded")
		}
		b.WriteString(s)
	case []phpjson.Value:
		encodeValue(b, listMap(t), depth)
	case *phpjson.Map:
		if t.Len() == 0 {
			if t.ForceObject {
				b.WriteString("{}")
			} else {
				b.WriteString("[]")
			}
			return
		}
		if t.IsList() && !t.ForceObject {
			b.WriteByte('[')
			for i, k := range t.Keys() {
				if i > 0 {
					b.WriteByte(',')
				}
				nl(depth + 1)
				x, _ := t.Get(k)
				encodeValue(b, x, depth+1)
			}
			nl(depth)
			b.WriteByte(']')
			return
		}
		b.WriteByte('{')
		for i, k := range t.Keys() {
			if i > 0 {
				b.WriteByte(',')
			}
			nl(depth + 1)
			encodeValue(b, k, depth+1)
			b.WriteString(": ")
			x, _ := t.Get(k)
			encodeValue(b, x, depth+1)
		}
		nl(depth)
		b.WriteByte('}')
	default:
		phpError("RuntimeException", "JSON encoding failed: Unknown error")
	}
}

// validateSchemaLax is JsonFile::validateSchema(JsonFile::LAX_SCHEMA). It
// returns *jsonmanip.Exception for unreadable/invalid files and
// *ValidationError for schema violations.
func validateSchemaLax(path string) error {
	if !isReadable(path) {
		return &jsonmanip.Exception{Class: "RuntimeException", Message: "The file \"" + path + "\" is not readable."}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return &jsonmanip.Exception{Class: "RuntimeException", Message: "The file \"" + path + "\" is not readable.", Err: err}
	}
	data, derr := phpjson.Decode(content)
	if derr != nil {
		data = nil
		if string(content) != "null" {
			if err := jsonfile.ValidateSyntax(string(content), path); err != nil {
				class := "Seld\\JsonLint\\ParsingException"
				var ee *jsonfile.EncodingError
				if errors.As(err, &ee) {
					class = "UnexpectedValueException"
				}
				return &jsonmanip.Exception{Class: class, Message: err.Error(), Err: err}
			}
		}
	}
	errs, err := jsonschema.ValidateComposerJSON(data, jsonschema.LaxSchema)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return &ValidationError{Message: "\"" + path + "\" does not match the expected JSON schema", Errors: errs}
	}
	return nil
}

// ValidationError is Composer\Json\JsonValidationException.
type ValidationError struct {
	Message string
	Errors  []string
}

func (e *ValidationError) Error() string { return e.Message }

// safeFilePutContents is Filesystem::safeFilePutContents: write to a
// temporary sibling and rename it over the target (keeping mode, owner and
// group, writing through symlinks); fall back to a direct write when that is
// not possible. Unlike PHP, a failed direct write is reported as an error
// instead of a warning.
func safeFilePutContents(path, content string) error {
	if replaceFile(path, content) {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0o666)
}

func replaceFile(path, content string) bool {
	// stream wrappers (php://memory, ...) do not support renaming
	if strings.Contains(path, "://") {
		return false
	}
	// write through symlinks instead of replacing them
	if lst, err := os.Lstat(path); err == nil && lst.Mode()&os.ModeSymlink != 0 {
		real, err := filepath.EvalSymlinks(path)
		if err != nil {
			return false
		}
		path = real
	}
	// rename() only needs the dir to be writable, so read-only files are left to the direct write
	if _, err := os.Stat(path); err == nil && !isWritable(path) {
		return false
	}

	var rnd [5]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return false
	}
	tempPath := path + "~" + hex.EncodeToString(rnd[:]) + ".tmp"
	fh, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		return false
	}
	newMode := os.FileMode(0o666)
	if st, err := fh.Stat(); err == nil {
		newMode = st.Mode().Perm() // 0666 & ~umask
	}
	// restrict permissions before writing as the content may be sensitive, e.g. auth.json
	written := os.Chmod(tempPath, 0o600) == nil

	st, statErr := os.Stat(path)
	if written && statErr == nil {
		// the renamed file would be ours, so only replace it if its owner and group can be kept
		written = keepOwner(fh, tempPath, st)
	}

	if written {
		n, err := fh.WriteString(content)
		written = err == nil && n == len(content)
	}
	written = fh.Close() == nil && written

	if written {
		mode := newMode
		if statErr == nil {
			mode = st.Mode().Perm()
		}
		_ = os.Chmod(tempPath, mode)
		if os.Rename(tempPath, path) == nil {
			return true
		}
	}
	_ = os.Remove(tempPath)
	return false
}

// isReadable is Filesystem::isReadable.
func isReadable(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	if st.IsDir() {
		d, err := os.Open(path)
		if err != nil {
			return false
		}
		d.Close()
		return true
	}
	fh, err := os.Open(path)
	if err != nil {
		return false
	}
	fh.Close()
	return true
}
