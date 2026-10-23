package cio

import (
	"strings"
)

// TableStyle is the subset of Symfony\Component\Console\Helper\TableStyle
// used by Composer commands.
type TableStyle struct {
	HorizontalOutside, HorizontalInside string
	VerticalOutside, VerticalInside     string
	// Crossing chars in Symfony order: cross, topLeft, topMid, topRight,
	// midRight, bottomRight, bottomMid, bottomLeft, midLeft,
	// topLeftBottom, topMidBottom, topRightBottom.
	Crossings         [12]string
	CellHeaderFormat  string
	CellRowFormat     string
	CellRowContentFmt string
	PadLeft           bool
}

func uniformCrossings(c string) [12]string {
	var a [12]string
	for i := range a {
		a[i] = c
	}
	return a
}

// TableStyles are Symfony's predefined styles.
var TableStyles = map[string]TableStyle{
	"default":             {HorizontalOutside: "-", HorizontalInside: "-", VerticalOutside: "|", VerticalInside: "|", Crossings: uniformCrossings("+"), CellHeaderFormat: "<info>%s</info>", CellRowFormat: "%s", CellRowContentFmt: " %s "},
	"borderless":          {HorizontalOutside: "=", HorizontalInside: "=", VerticalOutside: " ", VerticalInside: " ", Crossings: uniformCrossings(" "), CellHeaderFormat: "<info>%s</info>", CellRowFormat: "%s", CellRowContentFmt: " %s "},
	"compact":             {Crossings: uniformCrossings(""), CellHeaderFormat: "<info>%s</info>", CellRowFormat: "%s", CellRowContentFmt: "%s "},
	"symfony-style-guide": {HorizontalOutside: "-", HorizontalInside: "-", VerticalOutside: " ", VerticalInside: " ", Crossings: uniformCrossings(" "), CellHeaderFormat: "%s", CellRowFormat: "%s", CellRowContentFmt: " %s "},
	"box": {HorizontalOutside: "─", HorizontalInside: "─", VerticalOutside: "│", VerticalInside: "│",
		Crossings:        [12]string{"┼", "┌", "┬", "┐", "┤", "┘", "┴", "└", "├", "├", "┼", "┤"},
		CellHeaderFormat: "<info>%s</info>", CellRowFormat: "%s", CellRowContentFmt: " %s "},
}

// TableSeparator marks a separator row.
const TableSeparator = "\x00separator"

// Table renders rows like Symfony's Table helper (headers, multi-line cells,
// separators; no colspan/rowspan).
type Table struct {
	Style   TableStyle
	Headers []string
	Rows    [][]string
}

// NewTable creates a table with the named style.
func NewTable(style string) *Table {
	return &Table{Style: TableStyles[style]}
}

// Render writes the table to stdout (Write) of io.
func (t *Table) Render(o *IO) {
	t.RenderTo(o.Write)
}

// RenderTo renders through writeln.
func (t *Table) RenderTo(writeln func(string)) {
	type row struct {
		cells     []string
		separator bool
		header    bool
	}
	var groups [][]row
	addGroup := func(cells []string, header bool) {
		if len(cells) == 1 && cells[0] == TableSeparator {
			groups = append(groups, []row{{separator: true}})
			return
		}
		// split multi-line cells into extra rows
		var extra [][]string
		first := make([]string, len(cells))
		for c, cell := range cells {
			lines := strings.Split(strings.ReplaceAll(cell, "\r\n", "\n"), "\n")
			first[c] = lines[0]
			for li, l := range lines[1:] {
				for len(extra) <= li {
					extra = append(extra, make([]string, len(cells)))
				}
				extra[li][c] = l
			}
		}
		g := []row{{cells: first, header: header}}
		for _, e := range extra {
			g = append(g, row{cells: e, header: header})
		}
		groups = append(groups, g)
	}
	if len(t.Headers) > 0 {
		addGroup(t.Headers, true)
	}
	headerGroups := len(groups)
	for _, r := range t.Rows {
		addGroup(r, false)
	}
	cols := 0
	for _, g := range groups {
		for _, r := range g {
			if len(r.cells) > cols {
				cols = len(r.cells)
			}
		}
	}
	if cols == 0 {
		return
	}
	widths := make([]int, cols)
	contentExtra := Width(strings.Replace(t.Style.CellRowContentFmt, "%s", "", 1))
	for c := 0; c < cols; c++ {
		max := 0
		for _, g := range groups {
			for _, r := range g {
				if r.separator || c >= len(r.cells) {
					continue
				}
				if w := Width(StripTags(r.cells[c])); w > max {
					max = w
				}
			}
		}
		widths[c] = max + contentExtra
	}
	st := t.Style
	separator := func(kind string) {
		var horizontal, left, mid, right string
		switch kind {
		case "mid":
			horizontal, left, mid, right = st.HorizontalInside, st.Crossings[8], st.Crossings[0], st.Crossings[4]
		case "top":
			horizontal, left, mid, right = st.HorizontalOutside, st.Crossings[1], st.Crossings[2], st.Crossings[3]
		case "topbottom":
			horizontal, left, mid, right = st.HorizontalOutside, st.Crossings[9], st.Crossings[10], st.Crossings[11]
		default:
			horizontal, left, mid, right = st.HorizontalOutside, st.Crossings[7], st.Crossings[6], st.Crossings[5]
		}
		if st.HorizontalOutside == "" && st.HorizontalInside == "" && st.Crossings[0] == "" {
			return
		}
		var sb strings.Builder
		sb.WriteString(left)
		for c := 0; c < cols; c++ {
			sb.WriteString(strings.Repeat(horizontal, widths[c]))
			if c == cols-1 {
				sb.WriteString(right)
			} else {
				sb.WriteString(mid)
			}
		}
		writeln(sb.String())
	}
	renderRow := func(r row) {
		var sb strings.Builder
		sb.WriteString(st.VerticalOutside)
		format := st.CellRowFormat
		if r.header {
			format = st.CellHeaderFormat
		}
		for c := 0; c < cols; c++ {
			cell := ""
			if c < len(r.cells) {
				cell = r.cells[c]
			}
			content := strings.Replace(st.CellRowContentFmt, "%s", cell, 1)
			pad := widths[c] - Width(StripTags(content))
			if pad < 0 {
				pad = 0
			}
			if st.PadLeft {
				content = strings.Repeat(" ", pad) + content
			} else {
				content += strings.Repeat(" ", pad)
			}
			sb.WriteString(strings.Replace(format, "%s", content, 1))
			if c == cols-1 {
				sb.WriteString(st.VerticalOutside)
			} else {
				sb.WriteString(st.VerticalInside)
			}
		}
		writeln(sb.String())
	}
	// Symfony renders headers, a divider, then rows: the header group opens
	// with a top border and the first row after the divider with a
	// top-bottom border.
	isHeader := true
	isFirstRow := false
	for gi, g := range groups {
		if gi == headerGroups {
			isHeader = false
			isFirstRow = true
		}
		headerSepRendered := false
		for _, r := range g {
			if r.separator {
				separator("mid")
				continue
			}
			if len(r.cells) == 0 {
				continue
			}
			if isHeader && !headerSepRendered {
				separator("top")
				headerSepRendered = true
			}
			if isFirstRow {
				separator("topbottom")
				isFirstRow = false
			}
			renderRow(r)
		}
	}
	separator("bottom")
}

var plainIO = &IO{}

// StripTags is Helper::removeDecoration: formatting tags are removed.
func StripTags(s string) string {
	return plainIO.Format(s)
}

// Width is Helper::width (mb_strwidth): East Asian wide runes count 2.
func Width(s string) int {
	w := 0
	for _, r := range s {
		if isWide(r) {
			w += 2
		} else {
			w++
		}
	}
	return w
}

func isWide(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) || (r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) ||
		(r >= 0xAC00 && r <= 0xD7A3) || (r >= 0xF900 && r <= 0xFAFF) || (r >= 0xFE30 && r <= 0xFE4F) ||
		(r >= 0xFF00 && r <= 0xFF60) || (r >= 0xFFE0 && r <= 0xFFE6) || (r >= 0x20000 && r <= 0x3FFFD)
}
