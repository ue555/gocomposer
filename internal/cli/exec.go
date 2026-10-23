package cli

// Port of Composer\Command\ExecCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// The interactive binary selection is in select_cmds.go.

import (
	"path/filepath"
	"sort"
	"strings"

	"gocomposer/internal/composer"
)

func init() { registerCommand(execCommand) }

func execCommand() *Command {
	return &Command{Name: "exec", Run: runExec}
}

func runExec(ctx *Context) (int, error) {
	if ctx.IO.Interactive {
		if err := execInteract(ctx); err != nil {
			return exitFor(err)
		}
	}
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	in := ctx.In
	if _, hasBinary := in.args["binary"]; in.Bool("list") || !hasBinary {
		bins := execBinaries(c, true)
		if len(bins) == 0 {
			return 1, &composer.UserError{Msg: "No binaries found in composer.json or in bin-dir (" + c.Config.GetString("bin-dir") + ")", Code: 1}
		}
		ctx.IO.Write("<comment>Available binaries:</comment>")
		for _, b := range bins {
			ctx.IO.Write("<info>- " + b + "</info>")
		}
		return 0, nil
	}
	binary := in.Arg("binary")
	d := composer.NewEventDispatcher(c)
	d.RunScripts = !c.DisableScripts
	d.AddListener("__exec_command", binary)
	return d.DispatchScript("__exec_command", true, in.Args("args"), nil)
}

// execBinaries is ExecCommand::getBinaries.
func execBinaries(c *composer.Composer, forDisplay bool) []string {
	binDir := c.Config.GetString("bin-dir")
	matches, _ := filepath.Glob(globEscape(binDir) + "/*")
	sort.Strings(matches)
	var bins []string
	for _, m := range matches {
		// glob(3) skips dot files
		if !strings.HasPrefix(filepath.Base(m), ".") {
			bins = append(bins, m)
		}
	}
	local := c.Root().Binaries()
	all := append([]string(nil), bins...)
	for _, b := range local {
		if forDisplay {
			b += " (local)"
		}
		all = append(all, b)
	}
	var out []string
	prev, hasPrev := "", false
	for _, b := range all {
		if hasPrev && b == prev+".bat" {
			continue
		}
		prev, hasPrev = b, true
		out = append(out, filepath.Base(b))
	}
	return out
}

// globEscape escapes glob metacharacters of a literal path prefix.
func globEscape(p string) string {
	var out []rune
	for _, r := range p {
		switch r {
		case '*', '?', '[', '\\':
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	return string(out)
}
