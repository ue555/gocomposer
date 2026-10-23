package jsonlint

// Port of Seld\JsonLint\Lexer (seld/jsonlint 1.12.1, MIT, Copyright (c) 2011
// Jordi Boggiano). The PCRE rules are matched byte-wise by hand to keep the
// exact backtracking behaviour of the original patterns.

import (
	"strconv"
	"strings"
)

type lexer struct {
	flags    int
	input    string
	more     bool
	done     bool
	offset   int
	match    string
	yylineno int
	yyleng   int
	yytext   string
	yylloc   loc
}

func newLexer(flags int, input string) *lexer {
	return &lexer{flags: flags, input: input, yylloc: loc{firstLine: 1, firstColumn: 0, lastLine: 1, lastColumn: 0}}
}

func (l *lexer) parseError(msg string) {
	panic(&ParsingError{Msg: msg})
}

func (l *lexer) lex() int {
	for {
		symbol := l.next()
		switch symbol {
		case tSkipSpace, tBreakLine:
		case tComment, tOpenComment:
			if l.flags&AllowComments == 0 {
				l.parseError("Lexical error on line " + strconv.Itoa(l.yylineno+1) + ". Comments are not allowed.\n" + l.showPosition())
			}
			until := tCloseComment
			if symbol == tComment {
				until = tBreakLine
			}
			s := l.next()
			for s != until && !l.done {
				s = l.next()
			}
			if l.done {
				return symbolEOFToken
			}
		case tCloseComment:
			l.parseError("Lexical error on line " + strconv.Itoa(l.yylineno+1) + ". Unexpected token.\n" + l.showPosition())
		default:
			return symbol
		}
	}
}

func (l *lexer) showPosition() string {
	if l.yylineno == 0 && l.offset == 1 && l.match != "{" {
		return l.match + "...\n^"
	}
	pre := strings.ReplaceAll(l.pastInput(), "\n", "")
	n := len(pre) - 1
	if n < 0 {
		n = 0
	}
	return pre + strings.ReplaceAll(l.upcomingInput(), "\n", "") + "\n" + strings.Repeat("-", n) + "^"
}

func substr(s string, start, length int) string {
	if start >= len(s) {
		return ""
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}

func (l *lexer) pastInput() string {
	pastLength := l.offset - len(l.match)
	prefix := ""
	if pastLength > 20 {
		prefix = "..."
	}
	start := pastLength - 20
	if start < 0 {
		start = 0
	}
	n := pastLength
	if n > 20 {
		n = 20
	}
	if n < 0 {
		n = 0
	}
	return prefix + substr(l.input, start, n)
}

func (l *lexer) upcomingInput() string {
	next := l.match
	if len(next) < 20 {
		next += substr(l.input, l.offset, 20-len(next))
	}
	suffix := ""
	if len(next) > 20 {
		suffix = "..."
	}
	return substr(next, 0, 20) + suffix
}

func (l *lexer) fullUpcomingInput() string {
	next := l.match
	if strings.HasPrefix(next, `"`) && strings.Count(next, `"`) == 1 {
		n := len(l.input)
		strEnd := n
		if n != l.offset {
			q := strPos(l.input, `"`, l.offset+1)
			if q <= 0 {
				q = n
			}
			nl := strPos(l.input, "\n", l.offset+1)
			if nl <= 0 {
				nl = n
			}
			strEnd = q
			if nl < strEnd {
				strEnd = nl
			}
		}
		next += substr(l.input, l.offset, strEnd-l.offset)
	} else if len(next) < 20 {
		next += substr(l.input, l.offset, 20-len(next))
	}
	return next
}

// strPos is PHP strpos (returns -1 when not found or offset out of range).
func strPos(s, needle string, offset int) int {
	if offset > len(s) {
		return -1
	}
	i := strings.Index(s[offset:], needle)
	if i < 0 {
		return -1
	}
	return offset + i
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

func isWord(c byte) bool {
	return c == '_' || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// wordBoundaryAt is PCRE \b at position i (preceded by a word char here).
func (l *lexer) wordBoundaryAt(i int) bool {
	before := i > 0 && isWord(l.input[i-1])
	after := i < len(l.input) && isWord(l.input[i])
	return before != after
}

// matchNumber implements -?([0-9]|[1-9][0-9]+)(\.[0-9]+)?([eE][+-]?[0-9]+)?\b
// with PCRE backtracking order.
func (l *lexer) matchNumber(pos int) int {
	s := l.input
	i := pos
	if i < len(s) && s[i] == '-' {
		i++
	}
	var intEnds []int
	if i < len(s) && isDigit(s[i]) {
		intEnds = append(intEnds, i+1)
	}
	if i < len(s) && s[i] >= '1' && s[i] <= '9' {
		j := i + 1
		for j < len(s) && isDigit(s[j]) {
			j++
		}
		for k := j; k >= i+2; k-- {
			intEnds = append(intEnds, k)
		}
	}
	for _, ie := range intEnds {
		var fracEnds []int
		if ie < len(s) && s[ie] == '.' {
			j := ie + 1
			for j < len(s) && isDigit(s[j]) {
				j++
			}
			for k := j; k >= ie+2; k-- {
				fracEnds = append(fracEnds, k)
			}
		}
		fracEnds = append(fracEnds, ie)
		for _, fe := range fracEnds {
			var expEnds []int
			if fe < len(s) && (s[fe] == 'e' || s[fe] == 'E') {
				j := fe + 1
				if j < len(s) && (s[j] == '+' || s[j] == '-') {
					j++
				}
				d := j
				for d < len(s) && isDigit(s[d]) {
					d++
				}
				for k := d; k >= j+1; k-- {
					expEnds = append(expEnds, k)
				}
			}
			expEnds = append(expEnds, fe)
			for _, ee := range expEnds {
				if l.wordBoundaryAt(ee) {
					return ee - pos
				}
			}
		}
	}
	return -1
}

// matchString implements the possessive string rule.
func (l *lexer) matchString(pos int) int {
	s := l.input
	if pos >= len(s) || s[pos] != '"' {
		return -1
	}
	i := pos + 1
	for i < len(s) {
		c := s[i]
		if c == '\\' && i+1 < len(s) && strings.IndexByte(`"bfnrt/\`, s[i+1]) >= 0 {
			i += 2
			continue
		}
		if c == '\\' && i+5 < len(s) && s[i+1] == 'u' && isHex(s[i+2]) && isHex(s[i+3]) && isHex(s[i+4]) && isHex(s[i+5]) {
			i += 6
			continue
		}
		if c > 0x1f && c != '\\' && c != '"' {
			i++
			continue
		}
		break
	}
	if i < len(s) && s[i] == '"' {
		return i + 1 - pos
	}
	return -1
}

func isHex(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (l *lexer) matchLiteral(pos int, lit string) int {
	if strings.HasPrefix(l.input[pos:], lit) && l.wordBoundaryAt(pos+len(lit)) {
		return len(lit)
	}
	return -1
}

// rule returns the match length of rule i at the current offset or -1.
func (l *lexer) rule(i int) int {
	s, pos := l.input, l.offset
	switch i {
	case 0: // \s*\n\r?
		j := pos
		lastNL := -1
		for j < len(s) && isSpace(s[j]) {
			if s[j] == '\n' {
				lastNL = j
			}
			j++
		}
		if lastNL < 0 {
			return -1
		}
		end := lastNL + 1
		if end < len(s) && s[end] == '\r' {
			end++
		}
		return end - pos
	case 1:
		j := pos
		for j < len(s) && isSpace(s[j]) {
			j++
		}
		if j == pos {
			return -1
		}
		return j - pos
	case 2:
		return l.matchNumber(pos)
	case 3:
		return l.matchString(pos)
	case 4, 5, 6, 7, 8, 9:
		if pos < len(s) && s[pos] == "{}[],:"[i-4] {
			return 1
		}
		return -1
	case 10:
		return l.matchLiteral(pos, "true")
	case 11:
		return l.matchLiteral(pos, "false")
	case 12:
		return l.matchLiteral(pos, "null")
	case 13: // \G$
		if pos == len(s) || (pos == len(s)-1 && s[pos] == '\n') {
			return 0
		}
		return -1
	case 14:
		if strings.HasPrefix(s[pos:], "//") {
			return 2
		}
		return -1
	case 15:
		if strings.HasPrefix(s[pos:], "/*") {
			return 2
		}
		return -1
	case 16:
		if strings.HasPrefix(s[pos:], "*/") {
			return 2
		}
		return -1
	case 17:
		if pos < len(s) && s[pos] != '\n' {
			return 1
		}
		return -1
	}
	return -1
}

var ruleTokens = []int{tBreakLine, tSkipSpace, 6, 4, 17, 18, 23, 24, 22, 21, 10, 11, 8, 14, tComment, tOpenComment, tCloseComment, tInvalid}

func (l *lexer) next() int {
	if l.done {
		return lexEOF
	}
	if l.offset == len(l.input) {
		l.done = true
	}
	if !l.more {
		l.yytext = ""
		l.match = ""
	}
	for i := 0; i < 18; i++ {
		n := l.rule(i)
		if n < 0 {
			continue
		}
		m := l.input[l.offset : l.offset+n]
		lines := strings.Split(m, "\n")[1:]
		l.yylineno += len(lines)
		lastCol := l.yylloc.lastColumn + len(m)
		if len(lines) > 0 {
			lastCol = len(lines[len(lines)-1])
		}
		l.yylloc = loc{firstLine: l.yylloc.lastLine, lastLine: l.yylineno + 1, firstColumn: l.yylloc.lastColumn, lastColumn: lastCol}
		l.yytext += m
		l.match += m
		l.yyleng = len(l.yytext)
		l.more = false
		l.offset += n
		tok := ruleTokens[i]
		if i == 3 {
			l.yytext = substr(l.yytext, 1, l.yyleng-2)
		}
		return tok
	}
	if l.offset == len(l.input) {
		return lexEOF
	}
	l.parseError("Lexical error on line " + strconv.Itoa(l.yylineno+1) + ". Unrecognized text.\n" + l.showPosition())
	return tInvalid
}
