// Console output helpers for the audit formats, ported from
// symfony/console v5.4.54 as vendored by Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (Formatter/OutputFormatter.php
// format/formatAndWrap/escape, Helper/Helper.php width/removeDecoration,
// Helper/Table.php and Helper/TableStyle.php "default" style,
// symfony/string AbstractString::wordwrap) and
// src/Composer/IO/ConsoleIO.php (sanitize, write).
// Symfony is MIT licensed, Copyright (c) Fabien Potencier.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"gocomposer/internal/cio"
)

// IO is the subset of Composer's IOInterface used by the policy layer.
// *cio.IO satisfies it.
type IO interface {
	Write(msg string)
	WriteError(msg string)
}

type nullIO struct{}

func (nullIO) Write(string)      {}
func (nullIO) WriteError(string) {}

// NullIO discards everything (Composer\IO\NullIO).
var NullIO IO = nullIO{}

func orNullIO(io IO) IO {
	if io == nil {
		return NullIO
	}
	return io
}

var reSanitize = regexp.MustCompile("\x1b\\[[\x30-\x3F]*[\x20-\x2F]*[\x40-\x7E]|\x1b\\].*?(?:\x1b\\\\|\x07)|\x1b.|[\x01-\x09\x0B\x0C\x0E-\x1A]|\r")

// Sanitize is ConsoleIO::sanitize with $allowNewlines = true: removes ANSI
// escape sequences and control characters (keeping \n and \r\n) and replaces
// invalid UTF-8 with '?'.
func Sanitize(s string) string {
	if !utf8.ValidString(s) {
		s = strings.ToValidUTF8(s, "?")
	}
	// Go regexp has no lookahead for "\r(?!\n)": sanitize around each \r\n
	parts := strings.Split(s, "\r\n")
	for i, p := range parts {
		parts[i] = reSanitize.ReplaceAllString(p, "")
	}
	return strings.Join(parts, "\r\n")
}

// write is ConsoleIO::write (sanitized; one call per message).
func write(io IO, msgs ...string) {
	for _, m := range msgs {
		io.Write(forIO(io, Sanitize(m)))
	}
}

// writeError is ConsoleIO::writeError.
func writeError(io IO, msgs ...string) {
	for _, m := range msgs {
		io.WriteError(forIO(io, Sanitize(m)))
	}
}

// forIO adapts Symfony escapes for internal/cio, whose formatter only
// understands "\<": "\>" and the "\0" trailing-backslash markers are resolved
// here like OutputFormatter::formatAndWrap's final strtr.
func forIO(io IO, msg string) string {
	if _, ok := io.(*cio.IO); !ok {
		return msg
	}
	msg = strings.ReplaceAll(msg, `\>`, ">")
	return strings.ReplaceAll(msg, "\x00", `\`)
}

// EscapeFormat is OutputFormatter::escape.
func EscapeFormat(text string) string {
	text = strings.NewReplacer("<", `\<`, ">", `\>`).Replace(text)
	return escapeTrailingBackslash(text)
}

func escapeTrailingBackslash(text string) string {
	if strings.HasSuffix(text, `\`) {
		l := len(text)
		text = strings.TrimRight(text, `\`)
		text = strings.ReplaceAll(text, "\x00", "")
		text += strings.Repeat("\x00", l-len(text))
	}
	return text
}

// ComposerStyles are the formatter styles of Composer's console output
// (Symfony defaults plus Composer's "warning" and "highlight").
var ComposerStyles = []string{"error", "info", "comment", "question", "warning", "highlight"}

// SymfonyDefaultStyles are the styles of a bare OutputFormatter (BufferIO).
var SymfonyDefaultStyles = []string{"error", "info", "comment", "question"}

// Formatter is a non-decorated Symfony OutputFormatter.
type Formatter struct {
	Styles []string
}

var (
	reFormatTag = regexp.MustCompile(`(?i)<(([a-z](?:[^\\<>]*|\\.)*)|/([a-z][^<>]*)?)>`)
	reStyleAttr = regexp.MustCompile(`([^=]+)=([^;]+)(;|$)`)
)

func (f Formatter) isStyle(tag string) bool {
	if inStrings(tag, f.Styles) {
		return true
	}
	ms := reStyleAttr.FindAllStringSubmatch(tag, -1)
	if ms == nil {
		return false
	}
	for _, m := range ms {
		switch strings.ToLower(m[1]) {
		case "fg", "bg", "href", "options":
		default:
			return false
		}
	}
	return true
}

// Format is OutputFormatter::format in non-decorated mode.
func (f Formatter) Format(msg string) string { return f.FormatAndWrap(msg, 0) }

// FormatAndWrap is OutputFormatter::formatAndWrap in non-decorated mode.
func (f Formatter) FormatAndWrap(message string, width int) string {
	offset := 0
	var output strings.Builder
	currentLineLength := 0
	for _, m := range reFormatTag.FindAllStringSubmatchIndex(message, -1) {
		pos := m[0]
		text := message[m[0]:m[1]]
		if pos != 0 && message[pos-1] == '\\' {
			continue
		}
		output.WriteString(applyCurrentStyle(message[offset:pos], output.String(), width, &currentLineLength))
		offset = pos + len(text)
		open := text[1] != '/'
		tag := ""
		if open {
			tag = message[m[4]:m[5]]
		} else if m[6] >= 0 {
			tag = message[m[6]:m[7]]
		}
		if !open && tag == "" {
			continue
		}
		if !f.isStyle(tag) {
			output.WriteString(applyCurrentStyle(text, output.String(), width, &currentLineLength))
		}
	}
	output.WriteString(applyCurrentStyle(message[offset:], output.String(), width, &currentLineLength))
	return strings.NewReplacer("\x00", `\`, `\<`, "<", `\>`, ">").Replace(output.String())
}

func applyCurrentStyle(text, current string, width int, currentLineLength *int) string {
	if text == "" {
		return ""
	}
	if width == 0 {
		return text
	}
	if *currentLineLength == 0 && current != "" {
		text = strings.TrimLeft(text, " \t\n\r\x00\x0B")
	}
	prefix := ""
	if *currentLineLength != 0 {
		i := width - *currentLineLength
		if i > len(text) {
			i = len(text)
		}
		if i < 0 {
			i = 0
		}
		prefix = text[:i] + "\n"
		text = text[i:]
	}
	trailingNL := strings.HasSuffix(text, "\n")
	text = prefix + addLineBreaks(text, width)
	text = strings.TrimRight(text, "\n")
	if trailingNL {
		text += "\n"
	}
	if *currentLineLength == 0 && current != "" && !strings.HasSuffix(current, "\n") {
		text = "\n" + text
	}
	for _, line := range strings.Split(text, "\n") {
		*currentLineLength += len(line)
		if width <= *currentLineLength {
			*currentLineLength = 0
		}
	}
	return text
}

// addLineBreaks is AbstractString::wordwrap($width, "\n", true) on code points.
func addLineBreaks(text string, width int) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 1 && lines[0] == "" {
		return ""
	}
	var chars []string
	var mask []byte
	for i, line := range lines {
		if i > 0 {
			chars = append(chars, "\n")
			mask = append(mask, '#')
		}
		for _, r := range line {
			chars = append(chars, string(r))
			if r == ' ' {
				mask = append(mask, ' ')
			} else {
				mask = append(mask, '?')
			}
		}
	}
	wrapped := phpWordwrap(string(mask), width, "#", true)
	var out strings.Builder
	j := 0
	i := -1
	for b := 0; b < len(wrapped); b++ {
		if wrapped[b] != '#' {
			continue
		}
		for i++; i < b; i++ {
			if j < len(chars) {
				out.WriteString(chars[j])
			}
			j++
		}
		if j < len(chars) && (chars[j] == "\n" || chars[j] == " ") {
			j++
		}
		out.WriteString("\n")
	}
	for ; j < len(chars); j++ {
		out.WriteString(chars[j])
	}
	return out.String()
}

// phpWordwrap is PHP's wordwrap() for a single-byte break.
func phpWordwrap(text string, width int, brk string, cut bool) string {
	if text == "" {
		return ""
	}
	if width == 0 && cut {
		return text
	}
	var out strings.Builder
	laststart, lastspace := 0, 0
	current := 0
	for ; current < len(text); current++ {
		c := text[current]
		switch {
		case c == brk[0] && current+len(brk) < len(text) && strings.HasPrefix(text[current:], brk):
			out.WriteString(text[laststart : current+len(brk)])
			current += len(brk) - 1
			laststart = current + 1
			lastspace = laststart
		case c == ' ':
			if current-laststart >= width {
				out.WriteString(text[laststart:current])
				out.WriteString(brk)
				laststart = current + 1
			}
			lastspace = current
		case current-laststart >= width && cut && laststart >= lastspace:
			out.WriteString(text[laststart:current])
			out.WriteString(brk)
			laststart = current
			lastspace = current
		case current-laststart >= width && laststart < lastspace:
			out.WriteString(text[laststart:lastspace])
			out.WriteString(brk)
			laststart = lastspace + 1
			lastspace = laststart
		}
	}
	if laststart != current {
		out.WriteString(text[laststart:current])
	}
	return out.String()
}

var (
	reAnsiSGR = regexp.MustCompile("\x1b\\[[^m]*m")
	reOSC8    = regexp.MustCompile("\x1b]8;[^;]*;[^\x1b]*\x1b\\\\")
)

// RemoveDecoration is Helper::removeDecoration.
func (f Formatter) RemoveDecoration(s string) string {
	s = f.Format(s)
	s = reAnsiSGR.ReplaceAllString(s, "")
	return reOSC8.ReplaceAllString(s, "")
}

// displayWidth approximates Helper::width (terminal columns).
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		switch {
		case r < 0x20 || (r >= 0x7f && r < 0xa0):
		case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || r == 0x200B:
		case isWide(r):
			w += 2
		default:
			w++
		}
	}
	return w
}

func isWide(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) || (r >= 0x2E80 && r <= 0x303E) || (r >= 0x3041 && r <= 0x33FF) ||
		(r >= 0x3400 && r <= 0x4DBF) || (r >= 0x4E00 && r <= 0x9FFF) || (r >= 0xA000 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) || (r >= 0xF900 && r <= 0xFAFF) || (r >= 0xFE30 && r <= 0xFE4F) ||
		(r >= 0xFF00 && r <= 0xFF60) || (r >= 0xFFE0 && r <= 0xFFE6) || (r >= 0x1F300 && r <= 0x1F64F) ||
		(r >= 0x1F900 && r <= 0x1F9FF) || (r >= 0x20000 && r <= 0x3FFFD)
}

// Table is a Symfony console Table with the "default" style, supporting the
// features used by the Auditor: headers, rows, horizontal mode,
// setColumnWidth and setColumnMaxWidth.
type Table struct {
	formatter       Formatter
	headers         []string
	rows            [][]string
	horizontal      bool
	columnWidths    map[int]int
	columnMaxWidths map[int]int
}

// NewTable creates a table whose width calculations use the given styles.
func NewTable(styles []string) *Table {
	return &Table{formatter: Formatter{Styles: styles}, columnWidths: map[int]int{}, columnMaxWidths: map[int]int{}}
}

func (t *Table) SetHeaders(h []string) *Table        { t.headers = h; return t }
func (t *Table) AddRow(r []string) *Table            { t.rows = append(t.rows, r); return t }
func (t *Table) SetHorizontal() *Table               { t.horizontal = true; return t }
func (t *Table) SetColumnWidth(col, w int) *Table    { t.columnWidths[col] = w; return t }
func (t *Table) SetColumnMaxWidth(col, w int) *Table { t.columnMaxWidths[col] = w; return t }

type tableRow struct {
	cells   []string
	divider bool
}

// Render is Table::render; each output line is passed to writeln.
func (t *Table) Render(writeln func(string)) {
	var rows []tableRow
	if t.horizontal {
		for i, h := range t.headers {
			r := []string{h}
			for _, row := range t.rows {
				if i < len(row) {
					r = append(r, row[i])
				} else {
					r = append(r, "")
				}
			}
			rows = append(rows, tableRow{cells: r})
		}
	} else {
		if len(t.headers) > 0 {
			rows = append(rows, tableRow{cells: t.headers})
		}
		rows = append(rows, tableRow{divider: true})
		for _, r := range t.rows {
			rows = append(rows, tableRow{cells: r})
		}
	}

	numberOfColumns := 0
	for _, r := range rows {
		if !r.divider && len(r.cells) > numberOfColumns {
			numberOfColumns = len(r.cells)
		}
	}

	// buildTableRows: wrap cells over their max width and split multi-line cells
	var groups [][]tableRow
	for _, r := range rows {
		if r.divider {
			groups = append(groups, []tableRow{r})
			continue
		}
		cells := append([]string(nil), r.cells...)
		extra := map[int][]string{}
		maxLines := 1
		for col, cell := range cells {
			if mw, ok := t.columnMaxWidths[col]; ok && displayWidth(t.formatter.RemoveDecoration(cell)) > mw {
				cell = t.formatter.FormatAndWrap(cell, mw)
			}
			if !strings.Contains(cell, "\n") {
				cells[col] = cell
				continue
			}
			eol := "\n"
			if strings.Contains(cell, "\r\n") {
				eol = "\r\n"
			}
			parts := strings.Split(cell, eol)
			for i, p := range parts {
				parts[i] = escapeTrailingBackslash(p)
			}
			cell = strings.Join(parts, eol)
			lines := strings.Split(strings.ReplaceAll(cell, eol, "<fg=default;bg=default></>"+eol), eol)
			cells[col] = lines[0]
			extra[col] = lines[1:]
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}
		group := []tableRow{{cells: cells}}
		for l := 1; l < maxLines; l++ {
			row := make([]string, len(cells))
			for col := range cells {
				if l-1 < len(extra[col]) {
					row[col] = extra[col][l-1]
				}
			}
			group = append(group, tableRow{cells: row})
		}
		groups = append(groups, group)
	}

	// calculateColumnsWidth
	widths := make([]int, numberOfColumns)
	for col := 0; col < numberOfColumns; col++ {
		maxLen := 0
		for _, g := range groups {
			for _, r := range g {
				if r.divider {
					continue
				}
				w := 0
				if col < len(r.cells) {
					w = displayWidth(t.formatter.RemoveDecoration(r.cells[col]))
				}
				if cw := t.columnWidths[col]; cw > w {
					w = cw
				}
				if mw, ok := t.columnMaxWidths[col]; ok && mw < w {
					w = mw
				}
				if w > maxLen {
					maxLen = w
				}
			}
		}
		widths[col] = maxLen + 2
	}

	separator := func() {
		if numberOfColumns == 0 {
			return
		}
		var sb strings.Builder
		sb.WriteString("+")
		for col := 0; col < numberOfColumns; col++ {
			sb.WriteString(strings.Repeat("-", widths[col]))
			sb.WriteString("+")
		}
		writeln(sb.String())
	}
	renderRow := func(cells []string, cellFormat, firstCellFormat string) {
		var sb strings.Builder
		sb.WriteString("|")
		for col := 0; col < numberOfColumns; col++ {
			cell := ""
			if col < len(cells) {
				cell = cells[col]
			}
			pad := widths[col] - 2 - displayWidth(t.formatter.RemoveDecoration(cell))
			if pad < 0 {
				pad = 0
			}
			content := " " + cell + " " + strings.Repeat(" ", pad)
			format := cellFormat
			if firstCellFormat != "" && col == 0 {
				format = firstCellFormat
			}
			sb.WriteString(strings.Replace(format, "%s", content, 1))
			sb.WriteString("|")
		}
		writeln(sb.String())
	}

	isHeader := !t.horizontal
	isFirstRow := t.horizontal
	for _, g := range groups {
		headerSeparatorRendered := false
		for _, r := range g {
			if r.divider {
				isHeader = false
				isFirstRow = true
				continue
			}
			if len(r.cells) == 0 {
				continue
			}
			if isHeader && !headerSeparatorRendered {
				separator()
				headerSeparatorRendered = true
			}
			if isFirstRow {
				separator()
				isFirstRow = false
			}
			if t.horizontal {
				renderRow(r.cells, "%s", "<info>%s</info>")
			} else if isHeader {
				renderRow(r.cells, "<info>%s</info>", "")
			} else {
				renderRow(r.cells, "%s", "")
			}
		}
	}
	separator()
}
