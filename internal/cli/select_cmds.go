package cli

// Ports of RunScriptCommand::interact and ExecCommand::interact
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano): without a script/binary argument the
// user picks one from a ChoiceQuestion.

import (
	"strconv"

	"gocomposer/internal/cio"
)

// runScriptInteract is RunScriptCommand::interact.
func runScriptInteract(ctx *Context) error {
	in := ctx.In
	c, err := requireComposer(ctx)
	if err != nil {
		return err
	}
	scripts := c.Root().Scripts()
	if scripts == nil || scripts.Len() == 0 {
		return nil
	}
	if _, ok := in.args["script"]; ok || in.Bool("list") {
		return nil
	}
	var choices []cio.Choice
	for _, name := range scripts.Keys() {
		desc := ""
		if cmd, err := ctx.App.findCommand(name); err == nil {
			desc = cmd.Description
		}
		choices = setChoice(choices, name, desc)
	}
	res, err := ctx.IO.Select("Script to run: ", choices, true, "", 1, `Invalid script name "%s"`, false)
	if err != nil {
		return err
	}
	if s, ok := cio.AnswerString(res); ok {
		in.args["script"] = []string{s}
	}
	return nil
}

// execInteract is ExecCommand::interact.
func execInteract(ctx *Context) error {
	in := ctx.In
	c, err := requireComposer(ctx)
	if err != nil {
		return err
	}
	binaries := execBinaries(c, false)
	if len(binaries) == 0 {
		return nil
	}
	if _, ok := in.args["binary"]; ok || in.Bool("list") {
		return nil
	}
	choices := make([]cio.Choice, len(binaries))
	for i, b := range binaries {
		choices[i] = cio.Choice{Key: strconv.Itoa(i), Value: b}
	}
	res, err := ctx.IO.Select("Binary to run: ", choices, false, "", 1, `Invalid binary name "%s"`, false)
	if err != nil {
		return err
	}
	s, _ := cio.AnswerString(res)
	idx, _ := strconv.Atoi(s)
	in.args["binary"] = []string{binaries[idx]}
	return nil
}
