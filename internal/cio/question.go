package cio

// Port of Symfony\Component\Console\Helper\QuestionHelper (symfony/console
// v5.4.54 as locked by composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Fabien Potencier), Symfony\Component\Console\Question\Question,
// ChoiceQuestion and Cursor, and of the question methods of
// Composer\IO\ConsoleIO and Composer\Question\StrictConfirmationQuestion
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).
//
// Answers are read from the IO's input (STDIN) through one shared buffer, like
// PHP's buffered STDIN stream. Hidden answers and the autocompletion used by
// choice questions change the terminal mode like `stty -echo` and
// `stty -icanon -echo` and restore it afterwards.

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// MissingInputError is Symfony's MissingInputException ("Aborted."): the
// input ended before an answer was given.
type MissingInputError struct{}

func (*MissingInputError) Error() string { return "Aborted." }

// InvalidArgumentError is an \InvalidArgumentException thrown by a question
// validator; it is shown as "<error>msg</error>" before asking again.
type InvalidArgumentError struct {
	Msg string
	// Console marks Symfony's Console\Exception\InvalidArgumentException
	// (ChoiceQuestion), rendered without the "In file line" title.
	Console bool
}

func (e *InvalidArgumentError) Error() string { return e.Msg }

// RuntimeError is a Symfony RuntimeException raised while asking (it aborts
// the validation attempts at once).
type RuntimeError struct{ Msg string }

func (e *RuntimeError) Error() string { return e.Msg }

// Question is Symfony\Component\Console\Question\Question. Answers are PHP
// values: nil, string, bool or []string.
type Question struct {
	Text    string
	Default any
	// Validator receives the answer (or the default) and returns the
	// validated value; an *InvalidArgumentError makes it ask again.
	Validator func(any) (any, error)
	// Normalizer is applied to the answer before validation.
	Normalizer func(any) any
	// MaxAttempts limits the validation attempts (0: unlimited).
	MaxAttempts int
	Hidden      bool
	// NoTrim keeps surrounding whitespace (Question::setTrimmable(false)).
	NoTrim bool
	// Autocomplete returns the autocompletion values for the typed text.
	Autocomplete func(string) []string

	choice *choiceQuestion
}

// Choice is one entry of a ChoiceQuestion (PHP array key => value).
type Choice struct{ Key, Value string }

type choiceQuestion struct {
	choices     []Choice
	assoc       bool
	multiselect bool
	prompt      string
	errorMsg    string
}

// input returns the shared buffered reader of the IO's input stream.
func (o *IO) input() *bufio.Reader {
	if o.inReader == nil {
		var r io.Reader = o.In
		if r == nil {
			r = os.Stdin
		}
		o.inReader = bufio.NewReaderSize(r, 8192)
	}
	return o.inReader
}

// inputFile is the terminal behind the input (nil when it is not a file).
func (o *IO) inputFile() *os.File {
	if o.In == nil {
		return os.Stdin
	}
	if f, ok := o.In.(*os.File); ok {
		return f
	}
	return nil
}

// AskQuestion is QuestionHelper::ask through ConsoleIO::askQuestion (the
// program status reports the blocked state while the user is asked).
func (o *IO) AskQuestion(q *Question, kind string) (any, error) {
	if !o.Interactive {
		return o.defaultAnswer(q)
	}
	previous := o.programStatus
	if previous == nil {
		previous = &ProgramStatus{State: "working"}
	}
	o.WriteProgramStatus(&ProgramStatus{State: "blocked", Message: q.Text, Kind: kind, hasMessage: true})
	defer o.WriteProgramStatus(previous)
	var ans any
	var err error
	if q.Validator == nil {
		ans, err = o.doAsk(q)
	} else {
		ans, err = o.validateAttempts(q)
	}
	var mi *MissingInputError
	if errors.As(err, &mi) {
		o.Interactive = false
		fallback, ferr := o.defaultAnswer(q)
		if ferr != nil {
			return nil, ferr
		}
		if fallback == nil {
			return nil, err
		}
		return fallback, nil
	}
	return ans, err
}

// defaultAnswer is QuestionHelper::getDefaultAnswer.
func (o *IO) defaultAnswer(q *Question) (any, error) {
	if q.Default == nil {
		return nil, nil
	}
	if q.Validator != nil {
		return q.Validator(q.Default)
	}
	return q.Default, nil
}

// validateAttempts is QuestionHelper::validateAttempts.
func (o *IO) validateAttempts(q *Question) (any, error) {
	var lastErr error
	attempts := q.MaxAttempts
	for q.MaxAttempts == 0 || attempts > 0 {
		attempts--
		if lastErr != nil {
			o.writeQuestionError(lastErr)
		}
		ans, err := o.doAsk(q)
		if err != nil {
			return nil, err
		}
		v, err := q.Validator(ans)
		if err == nil {
			return v, nil
		}
		var re *RuntimeError
		var mi *MissingInputError
		if errors.As(err, &re) || errors.As(err, &mi) {
			return nil, err
		}
		lastErr = err
	}
	return nil, lastErr
}

// writeQuestionError is QuestionHelper::writeError (no formatter helper is
// registered in Composer's helper set).
func (o *IO) writeQuestionError(err error) {
	o.write(o.Err, "<error>"+err.Error()+"</error>", true, Normal)
}

// doAsk is QuestionHelper::doAsk.
func (o *IO) doAsk(q *Question) (any, error) {
	o.writePrompt(q)
	var ret string
	var autocomplete func(string) []string
	if q.Autocomplete != nil {
		autocomplete = q.Autocomplete
	} else if q.choice != nil {
		autocomplete = q.choice.autocompleteValues
	}
	f := o.inputFile()
	// Terminal::hasSttyAvailable(): stty works on the inherited STDIN
	sttyOK := f != nil && IsTTY(f)
	if autocomplete == nil || !sttyOK {
		got := false
		if q.Hidden {
			v, err := o.hiddenResponse(f, sttyOK, !q.NoTrim)
			if err != nil {
				return nil, err
			}
			ret, got = v, true
		}
		if !got {
			line, ok := o.readLine()
			if !ok {
				return nil, &MissingInputError{}
			}
			ret = line
			if !q.NoTrim {
				ret = phpTrim(ret)
			}
		}
	} else {
		v, err := o.autocomplete(q, f, autocomplete)
		if err != nil {
			return nil, err
		}
		ret = v
		if !q.NoTrim {
			ret = phpTrim(ret)
		}
	}
	var ans any = ret
	if len(ret) == 0 {
		ans = q.Default
	}
	if q.Normalizer != nil {
		return q.Normalizer(ans), nil
	}
	return ans, nil
}

// readLine is fgets($inputStream, 4096).
func (o *IO) readLine() (string, bool) {
	r := o.input()
	var sb strings.Builder
	for sb.Len() < 4095 {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		sb.WriteByte(b)
		if b == '\n' {
			break
		}
	}
	if sb.Len() == 0 {
		return "", false
	}
	return sb.String(), true
}

// hiddenResponse is QuestionHelper::getHiddenResponse.
func (o *IO) hiddenResponse(f *os.File, sttyOK, trimmable bool) (string, error) {
	restore := func() {}
	if sttyOK {
		restore, _ = sttyMode(f, false)
	} else if f != nil && IsTTY(f) {
		return "", &RuntimeError{Msg: "Unable to hide the response."}
	}
	value, ok := o.readLine()
	restore()
	if !ok {
		return "", &MissingInputError{}
	}
	if trimmable {
		value = phpTrim(value)
	}
	o.write(o.Err, "", true, Normal)
	return value, nil
}

// writePrompt is QuestionHelper::writePrompt.
func (o *IO) writePrompt(q *Question) {
	msg := q.Text
	if c := q.choice; c != nil {
		lines := []string{q.Text}
		maxWidth := 0
		for _, ch := range c.choices {
			if w := strWidth(ch.Key); w > maxWidth {
				maxWidth = w
			}
		}
		for _, ch := range c.choices {
			lines = append(lines, fmt.Sprintf("  [<info>%s%s</info>] %s", ch.Key, strings.Repeat(" ", maxWidth-strWidth(ch.Key)), ch.Value))
		}
		for _, l := range lines {
			o.write(o.Err, l, true, Normal)
		}
		msg = c.prompt
	}
	o.write(o.Err, msg, false, Normal)
}

func strWidth(s string) int { return utf8.RuneCountInString(s) }

// phpTrim is PHP trim() with the default character list.
func phpTrim(s string) string { return strings.Trim(s, " \t\n\r\x00\x0B") }

// autocomplete is QuestionHelper::autocomplete.
func (o *IO) autocomplete(q *Question, f *os.File, ac func(string) []string) (string, error) {
	fullChoice, ret := "", ""
	i, ofs := 0, -1
	matches := map[int]string{}
	setMatches := func(list []string) int {
		matches = map[int]string{}
		for k, v := range list {
			matches[k] = v
		}
		return len(list)
	}
	numMatches := setMatches(ac(ret))
	restore, _ := sttyMode(f, true)
	o.setStyleCodes("hl", "30;47", "39;49")
	r := o.input()
	eof := false
	readN := func(n int) string {
		var sb strings.Builder
		for sb.Len() < n {
			b, err := r.ReadByte()
			if err != nil {
				eof = true
				break
			}
			sb.WriteByte(b)
		}
		return sb.String()
	}
	isMulti := q.choice != nil && q.choice.multiselect
	for !eof {
		c := readN(1)
		if c == "" && ret == "" && q.Default == nil {
			restore()
			return "", &MissingInputError{}
		}
		switch {
		case c == "\x7f":
			if numMatches == 0 && i != 0 {
				i--
				_, last := lastRune(fullChoice)
				o.writeErrorRawFormatted(fmt.Sprintf("\x1b[%dD", strWidth(last)))
				fullChoice = mbSubstr(fullChoice, i)
			}
			if i == 0 {
				ofs = -1
				numMatches = setMatches(ac(ret))
			} else {
				numMatches = 0
			}
			ret = mbSubstr(ret, i)
		case c == "\x1b":
			c += readN(2)
			if len(c) > 2 && (c[2] == 'A' || c[2] == 'B') {
				if c[2] == 'A' && ofs == -1 {
					ofs = 0
				}
				if numMatches == 0 {
					continue
				}
				if c[2] == 'A' {
					ofs--
				} else {
					ofs++
				}
				ofs = phpMod(numMatches+ofs, numMatches)
			}
		case c == "" || c[0] < 32:
			if c == "\t" || c == "\n" {
				if numMatches > 0 && ofs != -1 {
					ret = matches[ofs]
					entered := len(phpTrim(mostRecentlyEnteredValue(fullChoice)))
					remaining := ""
					if entered < len(ret) {
						remaining = ret[entered:]
					}
					o.writeErrorFormatted(remaining)
					fullChoice += remaining
					i = utf8.RuneCountInString(fullChoice)
					kept := map[int]string{}
					for k, m := range ac(ret) {
						if ret == "" || strings.HasPrefix(m, ret) {
							kept[k] = m
						}
					}
					matches = kept
					numMatches = len(kept)
					ofs = -1
				}
				if c == "\n" {
					o.writeErrorFormatted(c)
					restore()
					return fullChoice, nil
				}
				numMatches = 0
			}
			continue
		default:
			if c[0] >= 0x80 {
				switch c[0] & 0xF0 {
				case 0xC0, 0xD0:
					c += readN(1)
				case 0xE0:
					c += readN(2)
				case 0xF0:
					c += readN(3)
				}
			}
			o.writeErrorFormatted(c)
			ret += c
			fullChoice += c
			i++
			tempRet := ret
			if isMulti {
				tempRet = mostRecentlyEnteredValue(fullChoice)
			}
			numMatches = 0
			ofs = 0
			for _, v := range ac(ret) {
				if strings.HasPrefix(v, tempRet) {
					matches[numMatches] = v
					numMatches++
				}
			}
		}
		o.writeErrorRawFormatted("\x1b[K")
		if numMatches > 0 && ofs != -1 {
			o.writeErrorRawFormatted("\x1b7")
			entered := len(phpTrim(mostRecentlyEnteredValue(fullChoice)))
			m := matches[ofs]
			rest := ""
			if entered < len(m) {
				rest = m[entered:]
			}
			o.writeErrorFormatted("<hl>" + rest + "</hl>")
			o.writeErrorRawFormatted("\x1b8")
		}
	}
	restore()
	return fullChoice, nil
}

func phpMod(a, b int) int { return a % b }

func lastRune(s string) (int, string) {
	r, size := utf8.DecodeLastRuneInString(s)
	if size == 0 {
		return 0, ""
	}
	return size, string(r)
}

// mbSubstr is mb_substr($s, 0, $n).
func mbSubstr(s string, n int) string {
	if n <= 0 {
		return ""
	}
	idx := 0
	for k := range s {
		if idx == n {
			return s[:k]
		}
		idx++
	}
	return s
}

// mostRecentlyEnteredValue is QuestionHelper::mostRecentlyEnteredValue.
func mostRecentlyEnteredValue(entered string) string {
	if !strings.Contains(entered, ",") {
		return entered
	}
	choices := strings.Split(entered, ",")
	if last := phpTrim(choices[len(choices)-1]); last != "" {
		return last
	}
	return entered
}

// writeErrorFormatted is $output->write($msg) on the error output.
func (o *IO) writeErrorFormatted(msg string) { o.write(o.Err, msg, false, Normal) }

// writeErrorRawFormatted writes a cursor control sequence through the error
// output (Cursor methods use $output->write()).
func (o *IO) writeErrorRawFormatted(seq string) { o.write(o.Err, seq, false, Normal) }

// --- ConsoleIO question methods ---------------------------------------------

// Program status kinds (Osc7501ProgramStatusEncoder::KIND_*).
const (
	KindQuestion   = "question"
	KindPermission = "permission"
	KindAuth       = "auth"
)

// Ask is ConsoleIO::ask: the answer, or def ("" answers give the default);
// nil means PHP null.
func (o *IO) Ask(question string, def any) (any, error) {
	return o.AskQuestion(&Question{Text: Sanitize(question), Default: sanitizeDefault(def)}, KindQuestion)
}

// AskString is Ask returning the answer as a string (null as "") and
// whether it was non-null.
func (o *IO) AskString(question string, def any) (string, bool, error) {
	v, err := o.Ask(question, def)
	if err != nil {
		return "", false, err
	}
	s, ok := AnswerString(v)
	return s, ok, nil
}

// AnswerString converts an answer to a PHP string cast; ok is false for null.
func AnswerString(v any) (string, bool) {
	switch x := v.(type) {
	case nil:
		return "", false
	case string:
		return x, true
	case bool:
		if x {
			return "1", true
		}
		return "", true
	}
	return fmt.Sprint(v), true
}

var (
	reConfirmTrue  = regexp.MustCompile(`(?i)^y(?:es)?$`)
	reConfirmFalse = regexp.MustCompile(`(?i)^no?$`)
)

// AskConfirmation is ConsoleIO::askConfirmation (StrictConfirmationQuestion).
func (o *IO) AskConfirmation(question string, def bool) (bool, error) {
	q := &Question{Text: Sanitize(question), Default: def}
	q.Normalizer = func(answer any) any {
		if b, ok := answer.(bool); ok {
			return b
		}
		s, _ := AnswerString(answer)
		if phpEmpty(s) && def {
			return def
		}
		if reConfirmTrue.MatchString(s) {
			return true
		}
		if reConfirmFalse.MatchString(s) {
			return false
		}
		return nil
	}
	q.Validator = func(answer any) (any, error) {
		if b, ok := answer.(bool); ok {
			return b, nil
		}
		return nil, &InvalidArgumentError{Msg: "Please answer yes, y, no, or n."}
	}
	v, err := o.AskQuestion(q, KindPermission)
	if err != nil {
		return false, err
	}
	b, _ := v.(bool)
	return b, nil
}

func phpEmpty(s string) bool { return s == "" || s == "0" }

// AskAndValidate is ConsoleIO::askAndValidate (attempts 0: unlimited).
func (o *IO) AskAndValidate(question string, validator func(any) (any, error), attempts int, def any) (any, error) {
	q := &Question{Text: Sanitize(question), Default: sanitizeDefault(def), Validator: validator, MaxAttempts: attempts}
	return o.AskQuestion(q, KindQuestion)
}

// AskAndHideAnswer is ConsoleIO::askAndHideAnswer.
func (o *IO) AskAndHideAnswer(question string) (any, error) {
	return o.AskQuestion(&Question{Text: Sanitize(question), Hidden: true}, KindAuth)
}

// Select is ConsoleIO::select. choices keep the PHP array order; assoc
// tells whether any key is a string (otherwise keys are 0..n-1). The result
// is the selected key (or keys when multiselect) as strings; nil/false-ish
// defaults are returned unchanged when nothing was selected.
func (o *IO) Select(question string, choices []Choice, assoc bool, def any, attempts int, errorMessage string, multiselect bool) (any, error) {
	sc := make([]Choice, len(choices))
	for i, c := range choices {
		sc[i] = Choice{Key: c.Key, Value: Sanitize(c.Value)}
	}
	cq := &choiceQuestion{choices: sc, assoc: assoc, multiselect: multiselect, prompt: " > ", errorMsg: errorMessage}
	q := &Question{Text: Sanitize(question), Default: sanitizeDefault(def), MaxAttempts: attempts, choice: cq}
	q.Validator = cq.validator(q)
	result, err := o.AskQuestion(q, KindQuestion)
	if err != nil {
		return nil, err
	}
	if assoc {
		return result, nil
	}
	if list, ok := result.([]string); ok {
		var out []string
		for _, c := range sc {
			for _, r := range list {
				if r == c.Value {
					out = append(out, c.Key)
					break
				}
			}
		}
		return out, nil
	}
	// (string) array_search($result, $choices, true)
	for _, c := range sc {
		if s, ok := result.(string); ok && s == c.Value {
			return c.Key, nil
		}
	}
	return "", nil
}

func sanitizeDefault(def any) any {
	if s, ok := def.(string); ok {
		return Sanitize(s)
	}
	return def
}

// autocompleteValues is Question::setAutocompleterValues($choices).
func (c *choiceQuestion) autocompleteValues(string) []string {
	var out []string
	if c.assoc {
		for _, ch := range c.choices {
			out = append(out, ch.Key)
		}
	}
	for _, ch := range c.choices {
		out = append(out, ch.Value)
	}
	return out
}

var (
	reMultiselect = regexp.MustCompile(`^[^,]+(?:,[^,]+)*$`)
	reIntKey      = regexp.MustCompile(`^(?:0|-?[1-9][0-9]*)$`)
)

// validator is ChoiceQuestion::getDefaultValidator.
func (c *choiceQuestion) validator(q *Question) func(any) (any, error) {
	return func(selected any) (any, error) {
		sel, _ := AnswerString(selected)
		var selectedChoices []string
		if c.multiselect {
			if !reMultiselect.MatchString(sel) {
				return nil, &InvalidArgumentError{Msg: phpSprintf(c.errorMsg, sel), Console: true}
			}
			selectedChoices = strings.Split(sel, ",")
		} else {
			selectedChoices = []string{sel}
		}
		if !q.NoTrim {
			for k, v := range selectedChoices {
				selectedChoices[k] = phpTrim(v)
			}
		}
		var results []string
		for _, value := range selectedChoices {
			var keys []string
			for _, ch := range c.choices {
				if ch.Value == value {
					keys = append(keys, ch.Key)
				}
			}
			if len(keys) > 1 {
				return nil, &InvalidArgumentError{Msg: fmt.Sprintf(`The provided answer is ambiguous. Value should be one of "%s".`, strings.Join(keys, `" or "`)), Console: true}
			}
			found := false
			result := ""
			// array_search($value, $choices) (loose comparison)
			for _, ch := range c.choices {
				if phpLooseEqual(ch.Value, value) {
					found, result = true, ch.Key
					break
				}
			}
			if !c.assoc {
				if found {
					idx, _ := strconv.Atoi(result)
					result = c.choices[idx].Value
				} else if reIntKey.MatchString(value) {
					if idx, err := strconv.Atoi(value); err == nil && idx >= 0 && idx < len(c.choices) {
						found, result = true, c.choices[idx].Value
					}
				}
			} else if !found {
				for _, ch := range c.choices {
					if ch.Key == value {
						found, result = true, value
						break
					}
				}
			}
			if !found {
				return nil, &InvalidArgumentError{Msg: phpSprintf(c.errorMsg, value), Console: true}
			}
			results = append(results, result)
		}
		if c.multiselect {
			return results, nil
		}
		return results[0], nil
	}
}

// phpSprintf formats a message with one %s placeholder.
func phpSprintf(format, arg string) string {
	s := strings.Replace(strings.ReplaceAll(format, "%%", "\x00"), "%s", strings.ReplaceAll(arg, "%", "\x01"), 1)
	return strings.NewReplacer("\x00", "%", "\x01", "%").Replace(s)
}

var reNumeric = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?[ \t\n\r\v\f]*$`)

// phpLooseEqual is PHP 8 string == string (numeric strings compare as numbers).
func phpLooseEqual(a, b string) bool {
	if a == b {
		return true
	}
	if reNumeric.MatchString(a) && reNumeric.MatchString(b) {
		fa, err1 := strconv.ParseFloat(strings.TrimSpace(a), 64)
		fb, err2 := strconv.ParseFloat(strings.TrimSpace(b), 64)
		return err1 == nil && err2 == nil && fa == fb
	}
	return false
}
