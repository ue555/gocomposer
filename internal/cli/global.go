package cli

// Port of Composer\Command\GlobalCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano),
// with the Symfony Console 5.4 ArgvInput::__toString / StringInput::tokenize
// round trip it relies on (symfony/console v5.4.54, MIT, (c) Fabien Potencier).

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/config"
)

func init() { registerCommand(globalCommand) }

func globalCommand() *Command {
	// a proxy command: its input is not validated against its definition
	return &Command{Name: "global", IgnoreValidationErrors: true, Run: runGlobal}
}

var (
	reGlobalWS       = regexp.MustCompile(`\s+`)
	reGlobalWord     = regexp.MustCompile(`\bg(?:l(?:o(?:b(?:a(?:l)?)?)?)?)?\b`)
	reArgvEqOption   = regexp.MustCompile(`^(-[^=]+=)(.+)`)
	reArgvPlainToken = regexp.MustCompile(`^[\w-]+\n?$`)
)

func runGlobal(ctx *Context) (int, error) {
	io := ctx.IO
	// Application::doRun: --help/-h shows the help of the proxy command
	for _, t := range ctx.In.Raw {
		if t == "--" {
			break
		}
		if t == "--help" || strings.HasPrefix(t, "--help=") || strings.HasPrefix(t, "-h") {
			return ctx.App.renderHelp(io, ctx.Cmd), nil
		}
	}
	str := argvInputString(ctx.In.Raw)
	// extract real command name
	var args []string
	for _, tok := range reGlobalWS.Split(str, -1) {
		if tok != "" && tok != "0" && tok[0] != '-' {
			args = append(args, tok)
			if len(args) >= 2 {
				break
			}
		}
	}
	// Command::run binding the input against the definition
	if len(args) < 2 {
		if _, err := parseInput(ctx.Cmd, ctx.In.Raw, "global"); err != nil {
			return 1, err
		}
		return 1, &CLIError{Msg: `Not enough arguments (missing: "command-name").`}
	}
	// the COMPOSER env var should not apply to the global execution scope
	if os.Getenv("COMPOSER") != "" {
		os.Unsetenv("COMPOSER")
	}
	cwd, _ := os.Getwd()
	cfg, err := config.CreateConfig(cwd, nil)
	if err != nil {
		return 1, err
	}
	home := cfg.GetString("home")
	if st, err := os.Stat(home); err != nil || !st.IsDir() {
		_ = os.MkdirAll(home, 0o777)
		if st, err := os.Stat(home); err != nil || !st.IsDir() {
			return 1, &composer.UserError{Msg: "Could not create home directory", Code: 1}
		}
	}
	if err := os.Chdir(home); err != nil {
		return 1, &composer.UserError{Msg: `Could not switch to home directory "` + home + `"`, Code: 1}
	}
	io.WriteError("<info>Changed current directory to " + home + "</info>")
	// create new input without "global" command prefix
	if loc := reGlobalWord.FindStringIndex(str); loc != nil {
		str = str[:loc[0]] + str[loc[1]:]
	}
	argv, err := stringInputTokenize(str)
	if err != nil {
		return 1, err
	}
	ctx.cache = nil
	// the nested Application::run has no running command until it found one
	ctx.App.running = nil
	return ctx.App.run(io, argv)
}

// argvInputString is Symfony's ArgvInput::__toString.
func argvInputString(tokens []string) string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		if m := reArgvEqOption.FindStringSubmatch(t); m != nil {
			out[i] = m[1] + escapeInputToken(m[2])
		} else if t != "" && t != "0" && t[0] != '-' {
			out[i] = escapeInputToken(t)
		} else {
			out[i] = t
		}
	}
	return strings.Join(out, " ")
}

// escapeInputToken is Input::escapeToken (escapeshellarg on unix).
func escapeInputToken(t string) string {
	if reArgvPlainToken.MatchString(t) {
		return t
	}
	return "'" + strings.ReplaceAll(t, "'", `'\''`) + "'"
}

const (
	reQuotedString = `(?:"([^"\\]*(?:\\.[^"\\]*)*)"|'([^'\\]*(?:\\.[^'\\]*)*)')`
)

var (
	reTokWS       = regexp.MustCompile(`^\s+`)
	reTokOptQuote = regexp.MustCompile(`^([^="'\s]+?)(=?)(` + reQuotedString + `+)`)
	reTokQuoted   = regexp.MustCompile(`^` + reQuotedString)
	reTokUnquoted = regexp.MustCompile(`^([^\s\\]+?)`)
)

// stringInputTokenize is StringInput::tokenize.
func stringInputTokenize(input string) ([]string, error) {
	var tokens []string
	var token *string
	appendTok := func(s string) {
		if token == nil {
			token = new(string)
		}
		*token += s
	}
	cursor := 0
	for cursor < len(input) {
		if input[cursor] == '\\' {
			cursor++
			if cursor < len(input) {
				appendTok(input[cursor : cursor+1])
			} else {
				appendTok("")
			}
			cursor++
			continue
		}
		rest := input[cursor:]
		var n int
		if m := reTokWS.FindString(rest); m != "" {
			if token != nil {
				tokens = append(tokens, *token)
				token = nil
			}
			n = len(m)
		} else if m := reTokOptQuote.FindStringSubmatch(rest); m != nil {
			q := m[3]
			q = q[1 : len(q)-1]
			q = strings.NewReplacer(`"'`, "", `'"`, "", `''`, "", `""`, "").Replace(q)
			appendTok(m[1] + m[2] + stripCSlashes(q))
			n = len(m[0])
		} else if m := reTokQuoted.FindString(rest); m != "" {
			appendTok(stripCSlashes(m[1 : len(m)-1]))
			n = len(m)
		} else if m := reTokUnquoted.FindStringSubmatch(rest); m != nil {
			appendTok(m[1])
			n = len(m[0])
		} else {
			end := cursor + 10
			if end > len(input) {
				end = len(input)
			}
			return nil, &CLIError{Msg: `Unable to parse input near "... ` + input[cursor:end] + ` ...".`}
		}
		cursor += n
	}
	if token != nil {
		tokens = append(tokens, *token)
	}
	return tokens, nil
}

// stripCSlashes is PHP stripcslashes().
func stripCSlashes(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case 'a':
			b.WriteByte('\a')
		case 'v':
			b.WriteByte('\v')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'x':
			j := i + 1
			for j < len(s) && j < i+3 && isHexDigit(s[j]) {
				j++
			}
			if j > i+1 {
				v, _ := strconv.ParseUint(s[i+1:j], 16, 8)
				b.WriteByte(byte(v))
				i = j - 1
			} else {
				b.WriteByte('x')
			}
		default:
			if s[i] >= '0' && s[i] <= '7' {
				j := i
				for j < len(s) && j < i+3 && s[j] >= '0' && s[j] <= '7' {
					j++
				}
				v, _ := strconv.ParseUint(s[i:j], 8, 16)
				b.WriteByte(byte(v))
				i = j - 1
			} else {
				b.WriteByte(s[i])
			}
		}
	}
	return b.String()
}

func isHexDigit(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
