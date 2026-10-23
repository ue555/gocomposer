// Package jsonlint ports seld/jsonlint 1.12.1 (JsonParser and Lexer), the
// JSON linter Composer uses to produce its "Parse error on line N" messages
// and to detect duplicate keys. Ported from
// vendor/seld/jsonlint/src/Seld/JsonLint/{JsonParser,Lexer,ParsingException,
// DuplicateKeyException}.php as locked by composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// seld/jsonlint is MIT licensed, Copyright (c) 2011 Jordi Boggiano (see LICENSE).
package jsonlint

import (
	"regexp"
	"strconv"
	"strings"
)

// Flags of JsonParser.
const (
	DetectKeyConflicts = 1
	AllowDuplicateKeys = 2
	ParseToAssoc       = 4
	AllowComments      = 8
)

// ParsingError is Seld\JsonLint\ParsingException.
type ParsingError struct {
	Msg string
	// Line is details['line'] when known (1-based).
	Line int
}

func (e *ParsingError) Error() string { return e.Msg }

// DuplicateKeyError is Seld\JsonLint\DuplicateKeyException.
type DuplicateKeyError struct {
	ParsingError
	Key string
}

const (
	lexEOF         = 1
	tInvalid       = -1
	tSkipSpace     = 0
	tError         = 2
	tBreakLine     = 3
	tComment       = 30
	tOpenComment   = 31
	tCloseComment  = 32
	symbolEOFToken = 14
)

var terminals = map[int]string{2: "error", 4: "STRING", 6: "NUMBER", 8: "NULL", 10: "TRUE", 11: "FALSE", 14: "EOF", 17: "{", 18: "}", 21: ":", 22: ",", 23: "[", 24: "]"}

// productions_[n] = {lhs symbol, rhs length}
var productions = map[int][2]int{1: {3, 1}, 2: {5, 1}, 3: {7, 1}, 4: {9, 1}, 5: {9, 1}, 6: {12, 2}, 7: {13, 1}, 8: {13, 1}, 9: {13, 1}, 10: {13, 1}, 11: {13, 1}, 12: {13, 1}, 13: {15, 2}, 14: {15, 3}, 15: {20, 3}, 16: {19, 1}, 17: {19, 3}, 18: {16, 2}, 19: {16, 3}, 20: {25, 1}, 21: {25, 3}}

type act struct {
	kind, arg int // kind 0 = goto (arg = state), 1 shift, 2 reduce, 3 accept
}

// table is the parser table; the key order of each state matters for the
// "Expected" list, so it is kept as ordered slices.
type entry struct {
	sym int
	a   act
}

func g(s int) act            { return act{0, s} }
func sh(s int) act           { return act{1, s} }
func rd(p int) act           { return act{2, p} }
func e(sym int, a act) entry { return entry{sym, a} }

var valueStart = []entry{e(3, g(5)), e(4, sh(12)), e(5, g(6)), e(6, sh(13)), e(7, g(3)), e(8, sh(9)), e(9, g(4)), e(10, sh(10)), e(11, sh(11))}

func concat(parts ...[]entry) []entry {
	var out []entry
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func endings(p int, syms ...int) []entry {
	var out []entry
	for _, s := range syms {
		out = append(out, e(s, rd(p)))
	}
	return out
}

var table = [][]entry{
	0:  concat(valueStart, []entry{e(12, g(1)), e(13, g(2)), e(15, g(7)), e(16, g(8)), e(17, sh(14)), e(23, sh(15))}),
	1:  {e(1, act{3, 0})},
	2:  {e(14, sh(16))},
	3:  endings(7, 14, 18, 22, 24),
	4:  endings(8, 14, 18, 22, 24),
	5:  endings(9, 14, 18, 22, 24),
	6:  endings(10, 14, 18, 22, 24),
	7:  endings(11, 14, 18, 22, 24),
	8:  endings(12, 14, 18, 22, 24),
	9:  endings(3, 14, 18, 22, 24),
	10: endings(4, 14, 18, 22, 24),
	11: endings(5, 14, 18, 22, 24),
	12: endings(1, 14, 18, 21, 22, 24),
	13: endings(2, 14, 18, 22, 24),
	14: {e(3, g(20)), e(4, sh(12)), e(18, sh(17)), e(19, g(18)), e(20, g(19))},
	15: concat(valueStart, []entry{e(13, g(23)), e(15, g(7)), e(16, g(8)), e(17, sh(14)), e(23, sh(15)), e(24, sh(21)), e(25, g(22))}),
	16: {e(1, rd(6))},
	17: endings(13, 14, 18, 22, 24),
	18: {e(18, sh(24)), e(22, sh(25))},
	19: endings(16, 18, 22),
	20: {e(21, sh(26))},
	21: endings(18, 14, 18, 22, 24),
	22: {e(22, sh(28)), e(24, sh(27))},
	23: endings(20, 22, 24),
	24: endings(14, 14, 18, 22, 24),
	25: {e(3, g(20)), e(4, sh(12)), e(20, g(29))},
	26: concat(valueStart, []entry{e(13, g(30)), e(15, g(7)), e(16, g(8)), e(17, sh(14)), e(23, sh(15))}),
	27: endings(19, 14, 18, 22, 24),
	28: concat(valueStart, []entry{e(13, g(31)), e(15, g(7)), e(16, g(8)), e(17, sh(14)), e(23, sh(15))}),
	29: endings(17, 18, 22),
	30: endings(15, 18, 22),
	31: endings(21, 22, 24),
}

func lookup(state, sym int) (act, bool) {
	for _, en := range table[state] {
		if en.sym == sym {
			return en.a, true
		}
	}
	return act{}, false
}

// Lint is JsonParser::lint: it returns nil for valid JSON, otherwise the
// *ParsingError (or *DuplicateKeyError with DetectKeyConflicts).
func Lint(input string, flags int) error {
	_, err := Parse(input, flags)
	return err
}

// value kinds tracked by the parser (only what key-conflict detection needs)
type object struct {
	keys map[string]bool // key -> value is non-null (PHP isset)
}

type member struct {
	key   string
	value interface{}
}

// Parse runs the parser. The returned value is only meaningful for
// distinguishing objects/nulls; callers use phpjson for real decoding.
func Parse(input string, flags int) (interface{}, error) {
	if strings.HasPrefix(input, "\xEF\xBB\xBF") {
		return nil, &ParsingError{Msg: "BOM detected, make sure your input does not include a Unicode Byte-Order-Mark"}
	}
	p := &parser{flags: flags}
	return p.parse(input)
}

type loc struct{ firstLine, lastLine, firstColumn, lastColumn int }

type parser struct {
	flags  int
	lexer  *lexer
	stack  []int
	vstack []interface{}
	lstack []loc
}

func (p *parser) popStack(n int) {
	p.stack = p.stack[:len(p.stack)-2*n]
	p.vstack = p.vstack[:len(p.vstack)-n]
	p.lstack = p.lstack[:len(p.lstack)-n]
}

func (p *parser) parse(input string) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if pe, ok := r.(error); ok {
				result, err = nil, pe
				return
			}
			panic(r)
		}
	}()
	p.stack = []int{0}
	p.vstack = []interface{}{nil}
	p.lexer = newLexer(p.flags, input)
	yytext := ""
	yylineno := 0
	recovering := 0
	yyloc := p.lexer.yylloc
	p.lstack = []loc{yyloc}
	symbol := 0
	haveSymbol := false
	preErrorSymbol := 0
	errStr := ""
	for {
		state := p.stack[len(p.stack)-1]
		var action act
		ok := false
		if state == 16 {
			action, ok = rd(6), true
		} else {
			if !haveSymbol {
				symbol = p.lexer.lex()
				haveSymbol = true
			}
			action, ok = lookup(state, symbol)
		}
		if !ok || action.kind == 0 {
			if recovering == 0 {
				var expected []string
				for _, en := range table[state] {
					if name, isTerm := terminals[en.sym]; isTerm && en.sym > 2 {
						expected = append(expected, "'"+name+"'")
					}
				}
				message := ""
				if contains(expected, "'STRING'") && len(p.lexer.match) > 0 && (p.lexer.match[0] == '"' || p.lexer.match[0] == '\'') {
					message = "Invalid string"
					full := p.lexer.fullUpcomingInput()
					if p.lexer.match[0] == '\'' {
						message += ", it appears you used single quotes instead of double quotes"
					} else if m := reUnescapedBackslash.FindStringSubmatch(toLatin1(full)); m != nil {
						message += ", it appears you have an unescaped backslash at: " + fromLatin1(m[1])
					} else if reUnterminated.MatchString(toLatin1(full)) {
						message += ", it appears you forgot to terminate a string, or attempted to write a multiline string which is invalid"
					}
				}
				errStr = "Parse error on line " + strconv.Itoa(yylineno+1) + ":\n"
				errStr += p.lexer.showPosition() + "\n"
				if message != "" {
					errStr += message
				} else {
					if len(expected) > 1 {
						errStr += "Expected one of: "
					} else {
						errStr += "Expected: "
					}
					errStr += strings.Join(expected, ", ")
				}
				if strings.HasSuffix(strings.Trim(p.lexer.pastInput(), " \t\n\r\x00\x0B"), ",") {
					errStr += " - It appears you have an extra trailing comma"
				}
				panic(&ParsingError{Msg: errStr, Line: p.lexer.yylineno})
			}
			// error recovery is unreachable: no state has an error action
			panic(&ParsingError{Msg: errStr})
		}
		switch action.kind {
		case 1: // shift
			p.stack = append(p.stack, symbol)
			p.vstack = append(p.vstack, p.lexer.yytext)
			p.lstack = append(p.lstack, p.lexer.yylloc)
			p.stack = append(p.stack, action.arg)
			haveSymbol = false
			if preErrorSymbol == 0 {
				yytext = p.lexer.yytext
				yylineno = p.lexer.yylineno
				yyloc = p.lexer.yylloc
				if recovering > 0 {
					recovering--
				}
			} else {
				symbol = preErrorSymbol
				haveSymbol = true
				preErrorSymbol = 0
			}
		case 2: // reduce
			n := productions[action.arg][1]
			currentToken := p.vstack[len(p.vstack)-n]
			ln := n
			if ln == 0 {
				ln = 1
			}
			position := loc{
				firstLine:   p.lstack[len(p.lstack)-ln].firstLine,
				lastLine:    p.lstack[len(p.lstack)-1].lastLine,
				firstColumn: p.lstack[len(p.lstack)-ln].firstColumn,
				lastColumn:  p.lstack[len(p.lstack)-1].lastColumn,
			}
			newToken, done, res := p.performAction(currentToken, yytext, yylineno, action.arg)
			if done {
				return res, nil
			}
			if n > 0 {
				p.popStack(n)
			}
			p.stack = append(p.stack, productions[action.arg][0])
			p.vstack = append(p.vstack, newToken)
			p.lstack = append(p.lstack, position)
			gt, _ := lookup(p.stack[len(p.stack)-2], p.stack[len(p.stack)-1])
			p.stack = append(p.stack, gt.arg)
		case 3:
			return true, nil
		}
	}
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

var (
	reUnescapedBackslash = regexp.MustCompile(`".+?(\\[^"bfnrt/\\u](...)?)`)
	reUnterminated       = regexp.MustCompile(`(?m)"(?:[^"]+|\\")*$`)
)

// toLatin1/fromLatin1 map bytes to runes 0-255 so Go regexps work on bytes
// like PCRE without the /u modifier.
func toLatin1(s string) string {
	r := make([]rune, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = rune(s[i])
	}
	return string(r)
}

func fromLatin1(s string) string {
	b := make([]byte, 0, len(s))
	for _, r := range s {
		b = append(b, byte(r))
	}
	return string(b)
}

func (p *parser) performAction(currentToken interface{}, yytext string, yylineno int, yystate int) (interface{}, bool, interface{}) {
	token := currentToken
	ln := len(p.vstack) - 1
	switch yystate {
	case 1:
		token = yytext
	case 2:
		token = yytext
	case 3:
		token = nil
	case 4:
		token = true
	case 5:
		token = false
	case 6:
		return p.vstack[ln-1], true, p.vstack[ln-1]
	case 13:
		token = &object{keys: map[string]bool{}}
	case 14:
		token = p.vstack[ln-1]
	case 15:
		token = member{key: unescapeKey(p.vstack[ln-2]), value: p.vstack[ln]}
	case 16:
		m := p.vstack[ln].(member)
		o := &object{keys: map[string]bool{}}
		o.keys[m.key] = m.value != nil
		token = o
	case 17:
		o := p.vstack[ln-2].(*object)
		m := p.vstack[ln].(member)
		set, exists := o.keys[m.key]
		if p.flags&DetectKeyConflicts != 0 && exists && set {
			errStr := "Parse error on line " + strconv.Itoa(yylineno+1) + ":\n"
			errStr += p.lexer.showPosition() + "\n"
			errStr += "Duplicate key: " + m.key
			panic(&DuplicateKeyError{ParsingError: ParsingError{Msg: errStr, Line: yylineno + 1}, Key: m.key})
		}
		o.keys[m.key] = m.value != nil
		token = o
	case 18:
		token = []interface{}{}
	case 19:
		token = p.vstack[ln-1]
	case 20:
		token = []interface{}{p.vstack[ln]}
	case 21:
		token = append(p.vstack[ln-2].([]interface{}), p.vstack[ln])
	}
	return token, false, nil
}

var reEscape = regexp.MustCompile(`\\["bfnrt/\\]|\\u[a-fA-F0-9]{4}`)

// unescapeKey applies stringInterpolation to a member key (only the key
// text matters for duplicate detection messages).
func unescapeKey(v interface{}) string {
	s, _ := v.(string)
	return reEscape.ReplaceAllStringFunc(s, func(m string) string {
		switch m {
		case `\\`:
			return `\`
		case `\"`:
			return `"`
		case `\b`:
			return "\b"
		case `\f`:
			return "\f"
		case `\n`:
			return "\n"
		case `\r`:
			return "\r"
		case `\t`:
			return "\t"
		case `\/`:
			return "/"
		}
		n, _ := strconv.ParseUint(m[2:], 16, 32)
		return string(rune(n))
	})
}
