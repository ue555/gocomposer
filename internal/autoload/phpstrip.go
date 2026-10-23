package autoload

import "bytes"

// stripWhitespace emulates php_strip_whitespace(): comments are dropped,
// runs of whitespace in PHP code become a single space, strings, heredocs
// and inline HTML are kept verbatim. It follows zend_strip() token handling
// closely enough for class discovery.
func stripWhitespace(src []byte) []byte {
	var out bytes.Buffer
	n := len(src)
	i := 0
	inPHP := false
	prevSpace := false
	isWS := func(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' }
	hasPrefixFold := func(b []byte, at int, s string) bool {
		if at+len(s) > len(b) {
			return false
		}
		return bytes.EqualFold(b[at:at+len(s)], []byte(s))
	}
	for i < n {
		if !inPHP {
			// inline HTML until an open tag
			start := i
			for i < n {
				if src[i] == '<' && i+1 < n && src[i+1] == '?' {
					break
				}
				i++
			}
			out.Write(src[start:i])
			if i >= n {
				break
			}
			// open tag: "<?php" + one whitespace char, "<?=" or "<?"
			if hasPrefixFold(src, i, "<?php") && (i+5 == n || isWS(src[i+5])) {
				end := i + 5
				if end < n {
					if src[end] == '\r' && end+1 < n && src[end+1] == '\n' {
						end += 2
					} else {
						end++
					}
				}
				out.Write(src[i:end])
				i = end
			} else if i+2 < n && src[i+2] == '=' {
				out.Write(src[i : i+3])
				i += 3
			} else {
				out.Write(src[i : i+2])
				i += 2
			}
			inPHP = true
			prevSpace = false
			continue
		}
		c := src[i]
		switch {
		case isWS(c):
			for i < n && isWS(src[i]) {
				i++
			}
			if !prevSpace {
				out.WriteByte(' ')
				prevSpace = true
			}
			continue
		case c == '?' && i+1 < n && src[i+1] == '>':
			end := i + 2
			if end < n && src[end] == '\n' {
				end++
			} else if end+1 < n && src[end] == '\r' && src[end+1] == '\n' {
				end += 2
			}
			out.Write(src[i:end])
			i = end
			inPHP = false
			prevSpace = false
			continue
		case c == '#' && !(i+1 < n && src[i+1] == '['), c == '/' && i+1 < n && src[i+1] == '/':
			for i < n && src[i] != '\n' && src[i] != '\r' {
				if src[i] == '?' && i+1 < n && src[i+1] == '>' {
					break
				}
				i++
			}
			continue
		case c == '/' && i+1 < n && src[i+1] == '*':
			end := bytes.Index(src[i+2:], []byte("*/"))
			if end < 0 {
				i = n
			} else {
				i = i + 2 + end + 2
			}
			continue
		case c == '\'' || c == '"' || c == '`':
			start := i
			i++
			for i < n {
				if src[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if src[i] == c {
					i++
					break
				}
				i++
			}
			out.Write(src[start:i])
			prevSpace = false
			continue
		case c == '<' && i+2 < n && src[i+1] == '<' && src[i+2] == '<':
			if end, label, ok := heredocBounds(src, i); ok {
				out.Write(src[i:end])
				i = end
				_ = label
				// zend_strip: write the following token unless it is whitespace, then "\n"
				if i < n && !isWS(src[i]) {
					// write the next single-char token (typically ";" or ",")
					out.WriteByte(src[i])
					i++
				} else {
					for i < n && isWS(src[i]) {
						i++
					}
				}
				out.WriteByte('\n')
				prevSpace = true
				continue
			}
		}
		out.WriteByte(c)
		i++
		prevSpace = false
	}
	return out.Bytes()
}

// heredocBounds returns the end offset of a heredoc/nowdoc starting at i.
func heredocBounds(src []byte, i int) (int, string, bool) {
	n := len(src)
	j := i + 3
	for j < n && (src[j] == ' ' || src[j] == '\t') {
		j++
	}
	quote := byte(0)
	if j < n && (src[j] == '\'' || src[j] == '"') {
		quote = src[j]
		j++
	}
	start := j
	for j < n && isIdentByte(src[j]) && !(j == start && src[j] >= '0' && src[j] <= '9') {
		j++
	}
	if j == start {
		return 0, "", false
	}
	label := string(src[start:j])
	if quote != 0 {
		if j >= n || src[j] != quote {
			return 0, "", false
		}
		j++
	}
	if j >= n || (src[j] != '\n' && src[j] != '\r') {
		return 0, "", false
	}
	// body: find a line whose first non-blank token is the label followed by a non-identifier char
	pos := j
	for pos < n {
		// move to start of next line
		for pos < n && src[pos] != '\n' && src[pos] != '\r' {
			pos++
		}
		for pos < n && (src[pos] == '\n' || src[pos] == '\r') {
			pos++
		}
		k := pos
		for k < n && (src[k] == ' ' || src[k] == '\t') {
			k++
		}
		if bytes.HasPrefix(src[k:], []byte(label)) {
			e := k + len(label)
			if e >= n || !isIdentByte(src[e]) {
				return e, label, true
			}
		}
	}
	return n, label, true
}
