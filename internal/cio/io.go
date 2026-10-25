// Package cio is the console IO (Composer\IO\ConsoleIO with Symfony console
// output formatting).
package cio

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Verbosity levels (OutputInterface::VERBOSITY_*).
const (
	Quiet       = 16
	Normal      = 32
	Verbose     = 64
	VeryVerbose = 128
	Debug       = 256
)

// IO writes formatted messages to stdout/stderr.
type IO struct {
	Out, Err    io.Writer
	Verbosity   int
	Decorated   bool
	Interactive bool
	mu          sync.Mutex
	lastErr     string
	custom      map[string]string
	// SymfonyStylesOnly drops Composer's own styles (warning, highlight),
	// like a BufferIO with a plain OutputFormatter: their tags stay literal.
	SymfonyStylesOnly bool
	// In is the input stream answers are read from (nil: os.Stdin).
	In       io.Reader
	inReader *bufio.Reader
	// program status reporting (status.go)
	statusEncoders []string
	programStatus  *ProgramStatus
	customUnset    map[string]string
}

// SetStyle registers a named foreground style (OutputFormatter::setStyle
// with new OutputFormatterStyle($color)).
func (o *IO) SetStyle(name, fg string) {
	if o.custom == nil {
		o.custom = map[string]string{}
	}
	o.custom[name] = colorCode(fg, 30)
}

// setStyleCodes registers a named style with explicit set/unset codes
// (OutputFormatter::setStyle with fg and bg colors).
func (o *IO) setStyleCodes(name, set, unset string) {
	if o.custom == nil {
		o.custom = map[string]string{}
	}
	if o.customUnset == nil {
		o.customUnset = map[string]string{}
	}
	o.custom[name] = set
	o.customUnset[name] = unset
}

// New creates an IO bound to the process streams.
func New() *IO {
	return &IO{Out: os.Stdout, Err: os.Stderr, Verbosity: Normal}
}

func (o *IO) IsVerbose() bool     { return o.Verbosity >= Verbose }
func (o *IO) IsVeryVerbose() bool { return o.Verbosity >= VeryVerbose }
func (o *IO) IsDebug() bool       { return o.Verbosity >= Debug }
func (o *IO) IsQuiet() bool       { return o.Verbosity <= Quiet }

var (
	// OutputFormatter::formatAndWrap: an opening tag may hold escaped
	// characters (\<, \>), a closing tag may be empty
	reTag     = regexp.MustCompile(`(?i)<(?:(/)([a-z][^<>]*)?|([a-z](?:[^\\<>]|\\.)*))>`)
	styles    = map[string]string{"info": "32", "comment": "33", "error": "37;41", "question": "30;46", "warning": "30;43", "highlight": "31"}
	reStyleKV = regexp.MustCompile(`([^=]+)=([^;]+)(;|$)`)
)

// Format applies Symfony OutputFormatter rules: known style tags are
// replaced by ANSI codes when decorated and stripped otherwise; "\<" escapes.
// Like OutputFormatter::formatAndWrap, every text segment is wrapped in the
// set/unset codes of the innermost open style (Color::apply), and href
// styles become OSC 8 hyperlinks (OutputFormatterStyle::apply).
func (o *IO) Format(msg string) string {
	var sb, seg strings.Builder
	var stack []string
	flush := func() {
		if seg.Len() == 0 {
			return
		}
		t := seg.String()
		seg.Reset()
		if o.Decorated && len(stack) > 0 {
			sb.WriteString(o.applyStyle(stack[len(stack)-1], t))
		} else {
			sb.WriteString(t)
		}
	}
	i := 0
	for i < len(msg) {
		if msg[i] == '\\' && i+1 < len(msg) && msg[i+1] == '<' {
			seg.WriteByte('<')
			i += 2
			continue
		}
		if msg[i] == '<' {
			loc := reTag.FindStringSubmatchIndex(msg[i:])
			if loc != nil && loc[0] == 0 {
				closing := loc[2] >= 0
				// style names are looked up as written; only the keys and
				// fg/bg/options values of inline styles are case-insensitive
				// (OutputFormatter::createStyleFromString)
				name := ""
				if closing && loc[4] >= 0 {
					name = msg[i+loc[4] : i+loc[5]]
				} else if !closing {
					name = msg[i+loc[6] : i+loc[7]]
				}
				full := msg[i : i+loc[1]]
				if closing {
					if name == "" || (len(stack) > 0 && stack[len(stack)-1] == name) {
						flush()
						if len(stack) > 0 {
							stack = stack[:len(stack)-1]
						}
						i += loc[1]
						continue
					}
				} else if _, known := styles[name]; (known && !(o.SymfonyStylesOnly && (name == "warning" || name == "highlight"))) || inlineStyle(name) || o.custom[name] != "" {
					flush()
					stack = append(stack, name)
					i += loc[1]
					continue
				}
				// OutputFormatter::formatAndWrap: the text before the tag
				// and the unknown tag itself are styled separately
				flush()
				seg.WriteString(full)
				flush()
				i += loc[1]
				continue
			}
		}
		seg.WriteByte(msg[i])
		i++
	}
	flush()
	return sb.String()
}

// applyStyle is OutputFormatterStyle::apply for a style name.
func (o *IO) applyStyle(name, text string) string {
	set, unset, href := o.styleCodes(name)
	if href != "" && handlesHrefGracefully() {
		text = "\x1b]8;;" + href + "\x1b\\" + text + "\x1b]8;;\x1b\\"
	}
	if set == "" {
		return text
	}
	return "\x1b[" + set + "m" + text + "\x1b[" + unset + "m"
}

// handlesHrefGracefully is OutputFormatterStyle::apply's terminal check
// (symfony/console v5.4).
func handlesHrefGracefully() bool {
	if os.Getenv("TERMINAL_EMULATOR") == "JetBrains-JediTerm" {
		return false
	}
	if k := os.Getenv("KONSOLE_VERSION"); k != "" && k != "0" {
		if phpIntPrefix(k) <= 201100 {
			return false
		}
	}
	_, idea := os.LookupEnv("IDEA_INITIAL_DIRECTORY")
	return !idea
}

// builtinStyles are the set/unset codes of the default Symfony and Composer
// styles (fg/bg colors).
var builtinStyles = map[string][2]string{
	"info": {"32", "39"}, "comment": {"33", "39"}, "error": {"37;41", "39;49"},
	"question": {"30;46", "39;49"}, "warning": {"30;43", "39;49"}, "highlight": {"31", "39"},
}

var styleOptions = map[string][2]string{"bold": {"1", "22"}, "underscore": {"4", "24"}, "blink": {"5", "25"}, "reverse": {"7", "27"}, "conceal": {"8", "28"}}

// styleCodes returns the Color set/unset codes and href of a style.
func (o *IO) styleCodes(name string) (string, string, string) {
	if c, ok := o.custom[name]; ok {
		if u, ok := o.customUnset[name]; ok {
			return c, u, ""
		}
		return c, "39", ""
	}
	if c, ok := builtinStyles[name]; ok {
		return c[0], c[1], ""
	}
	var set, unset []string
	var opts [][2]string
	href, fg, bg := "", "", ""
	for _, part := range strings.Split(name, ";") {
		k, v, _ := strings.Cut(part, "=")
		k = strings.ToLower(k)
		if k != "href" {
			v = strings.ToLower(v)
		}
		switch k {
		case "fg":
			fg = v
		case "bg":
			bg = v
		case "options":
			for _, opt := range strings.Split(v, ",") {
				if c, ok := styleOptions[strings.TrimSpace(opt)]; ok {
					opts = append(opts, c)
				}
			}
		case "href":
			href = strings.NewReplacer(`\<`, "<", `\>`, ">").Replace(v)
		}
	}
	// Color::set(): foreground, background, then options
	if fg != "" {
		set = append(set, colorCode(fg, 30))
		unset = append(unset, "39")
	}
	if bg != "" {
		set = append(set, colorCode(bg, 40))
		unset = append(unset, "49")
	}
	for _, c := range opts {
		set = append(set, c[0])
		unset = append(unset, c[1])
	}
	return strings.Join(set, ";"), strings.Join(unset, ";"), href
}

func ansi(name string) string {
	var o IO
	set, _, _ := o.styleCodes(name)
	if set == "" {
		return ""
	}
	return "\x1b[" + set + "m"
}

func colorCode(name string, base int) string {
	colors := map[string]int{"black": 0, "red": 1, "green": 2, "yellow": 3, "blue": 4, "magenta": 5, "cyan": 6, "white": 7, "default": 9}
	if c, ok := colors[name]; ok {
		return fmt.Sprint(base + c)
	}
	return fmt.Sprint(base + 9)
}

// Write writes to stdout.
func (o *IO) Write(msg string) { o.write(o.Out, msg, true, Normal) }

// WriteV writes to stdout at a verbosity.
func (o *IO) WriteV(msg string, newline bool, verbosity int) { o.write(o.Out, msg, newline, verbosity) }

// WriteError writes to stderr at normal verbosity.
func (o *IO) WriteError(msg string) { o.write(o.Err, msg, true, Normal); o.noteErr(msg, Normal) }

// WriteErrorV writes to stderr at a verbosity.
func (o *IO) WriteErrorV(msg string, newline bool, verbosity int) {
	o.write(o.Err, msg, newline, verbosity)
	o.noteErr(msg, verbosity)
}

// WriteRaw writes without formatting (Composer's writeRaw).
func (o *IO) WriteRaw(msg string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.Verbosity <= Quiet {
		return
	}
	fmt.Fprintln(o.Out, msg)
}

func (o *IO) write(w io.Writer, msg string, newline bool, verbosity int) {
	if verbosity > o.Verbosity {
		return
	}
	if o.Verbosity <= Quiet && verbosity != Quiet {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	s := o.Format(msg)
	if newline {
		s += "\n"
	}
	io.WriteString(w, s)
}

// Debugf is a convenience for debug output.
func (o *IO) Debugf(format string, a ...any) {
	o.WriteErrorV(fmt.Sprintf(format, a...), true, Debug)
}

// WriteErrorRaw writes an unformatted fragment to stderr (no newline) at
// normal verbosity, as ProgressBar does through the error output.
func (o *IO) WriteErrorRaw(msg string) {
	if o.Verbosity <= Quiet {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	io.WriteString(o.Err, msg)
}

// TerminalWidth is Symfony\Component\Console\Terminal::getWidth(): COLUMNS,
// else the size of the controlling terminal, else 80.
func TerminalWidth() int {
	// getenv('COLUMNS') !== false: (int) trim($width), even when 0
	if c, ok := os.LookupEnv("COLUMNS"); ok {
		return phpIntPrefix(strings.TrimSpace(c))
	}
	if w := ttyWidth(); w > 0 {
		return w
	}
	return 80
}

// phpIntPrefix is PHP's (int) cast of a string: the leading optional sign
// and digits, 0 when there are none.
func phpIntPrefix(s string) int {
	s = strings.TrimLeft(s, " \t\n\r\v\f")
	end := 0
	if end < len(s) && (s[end] == '+' || s[end] == '-') {
		end++
	}
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	n, _ := strconv.Atoi(s[:end])
	return n
}

// inlineStyle reports whether a tag is an inline style
// (OutputFormatter::createStyleFromString): key=value pairs separated by
// ";" with the keys fg, bg, href and options in any case.
func inlineStyle(name string) bool {
	m := reStyleKV.FindAllStringSubmatch(name, -1)
	if m == nil {
		return false
	}
	for _, kv := range m {
		switch strings.ToLower(kv[1]) {
		case "fg", "bg", "href", "options":
		default:
			return false
		}
	}
	return true
}
