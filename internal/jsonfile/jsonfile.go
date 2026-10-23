// Package jsonfile ports the decoding part of Composer\Json\JsonFile
// (parseJson/validateSyntax) from composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
package jsonfile

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"gocomposer/internal/jsonlint"
	"gocomposer/internal/phpjson"
)

// ConflictedContentHash is JsonFile::CONFLICTED_CONTENT_HASH.
const ConflictedContentHash = "VCS merge conflict detected. Please run `composer update --lock`."

// ParsingError is Seld\JsonLint\ParsingException as rethrown by
// JsonFile::validateSyntax.
type ParsingError struct {
	Msg  string
	Line int
}

func (e *ParsingError) Error() string { return e.Msg }

// EncodingError is the UnexpectedValueException for non UTF-8 input.
type EncodingError struct{ Msg string }

func (e *EncodingError) Error() string { return e.Msg }

var reLockConflict = regexp.MustCompile(`\r?\n<<<<<<< [^\r\n]+\r?\n\s+"content-hash": *"[0-9a-f]+", *\r?\n(?:\|{7} [^\r\n]+\r?\n\s+"content-hash": *"[0-9a-f]+", *\r?\n)?=======\r?\n\s+"content-hash": *"[0-9a-f]+", *\r?\n>>>>>>> [^\r\n]+(\r?\n)`)

// Parse is JsonFile::parseJson($json, $file). file may be "" (no file).
func Parse(raw []byte, file string) (phpjson.Value, error) {
	v, err := phpjson.Decode(raw)
	if err == nil {
		return v, nil
	}
	s := string(raw)
	if file != "" && strings.HasSuffix(file, ".lock") && strings.Contains(s, `"content-hash"`) {
		if loc := reLockConflict.FindStringSubmatchIndex(s); loc != nil {
			nl := s[loc[2]:loc[3]]
			replaced := s[:loc[0]] + "\n    \"content-hash\": \"" + ConflictedContentHash + "\"," + nl + s[loc[1]:]
			if v, err := phpjson.Decode([]byte(replaced)); err == nil && v != nil {
				return v, nil
			}
		}
	}
	return nil, ValidateSyntax(s, file)
}

// ValidateSyntax is JsonFile::validateSyntax; it returns the error Composer
// throws for input json_decode rejected.
func ValidateSyntax(s, file string) error {
	lerr := jsonlint.Lint(s, 0)
	if lerr == nil {
		if !utf8.ValidString(s) {
			if file == "" {
				return &EncodingError{"The input is not UTF-8, could not parse as JSON"}
			}
			return &EncodingError{`"` + file + `" is not UTF-8, could not parse as JSON`}
		}
		return nil
	}
	line := 0
	var pe *jsonlint.ParsingError
	if errors.As(lerr, &pe) {
		line = pe.Line
	}
	var de *jsonlint.DuplicateKeyError
	if errors.As(lerr, &de) {
		line = de.Line
	}
	if file == "" {
		return &ParsingError{Msg: "The input does not contain valid JSON\n" + lerr.Error(), Line: line}
	}
	return &ParsingError{Msg: `"` + file + "\" does not contain valid JSON\n" + lerr.Error(), Line: line}
}

// IndentDefault is JsonFile::INDENT_DEFAULT.
const IndentDefault = "    "

var reIndent = regexp.MustCompile(`(?m)^([ \t]+)"`)

// DetectIndenting is JsonFile::detectIndenting.
func DetectIndenting(content string) string {
	if m := reIndent.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return IndentDefault
}

var reLeading4 = regexp.MustCompile(`(?m)^ {4,}`)

// Encode is JsonFile::encode with the default flags and an indent.
func Encode(v phpjson.Value, indent string) (string, error) {
	s, err := phpjson.Encode(v, phpjson.ComposerFlags)
	if err != nil {
		return "", err
	}
	if indent != IndentDefault {
		s = reLeading4.ReplaceAllStringFunc(s, func(m string) string {
			return strings.Repeat(indent, len(m)/4)
		})
	}
	return s, nil
}

// Write is JsonFile::write with the default flags: the encoded value plus a
// newline is written (atomically) unless the file already has that content.
func Write(path string, v phpjson.Value, indent string) error {
	s, err := Encode(v, indent)
	if err != nil {
		return err
	}
	s += "\n"
	dir := filepath.Dir(path)
	if st, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0o777); err != nil {
			return &EncodingError{dir + " does not exist and could not be created."}
		}
	} else if !st.IsDir() {
		return &EncodingError{dir + " exists and is not a directory."}
	}
	if cur, err := os.ReadFile(path); err == nil && string(cur) == s {
		return nil
	}
	return SafeFilePutContents(path, []byte(s))
}

// SafeFilePutContents is Filesystem::safeFilePutContents/replaceFile: the
// content goes to "<path>~<hex>.tmp" (created 0600 since it may hold
// secrets) which is renamed over path with path's mode; symlinks are written
// through, and when replacing is not possible the file is written directly.
func SafeFilePutContents(path string, data []byte) error {
	if replaceFile(path, data) {
		return nil
	}
	return os.WriteFile(path, data, 0o666)
}

func replaceFile(path string, data []byte) bool {
	if strings.Contains(path, "://") {
		return false
	}
	if st, err := os.Lstat(path); err == nil && st.Mode()&os.ModeSymlink != 0 {
		real, err := filepath.EvalSymlinks(path)
		if err != nil {
			return false
		}
		path = real
	}
	st, statErr := os.Stat(path)
	if statErr == nil && !writable(path) {
		return false
	}
	var rnd [5]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return false
	}
	tmp := path + "~" + hex.EncodeToString(rnd[:]) + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return false
	}
	ok := os.Chmod(tmp, 0o600) == nil
	if ok && statErr == nil {
		ok = sameOwner(st, f)
	}
	if ok {
		n, err := f.Write(data)
		ok = err == nil && n == len(data)
	}
	if err := f.Close(); err != nil {
		ok = false
	}
	if ok {
		mode := os.FileMode(0o666) &^ os.FileMode(currentUmask())
		if statErr == nil {
			mode = st.Mode().Perm()
		}
		os.Chmod(tmp, mode)
		if os.Rename(tmp, path) == nil {
			return true
		}
	}
	os.Remove(tmp)
	return false
}

func writable(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
