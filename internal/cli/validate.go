package cli

// Port of Composer\Command\ValidateCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/validate"
)

func validateCommand() *Command {
	return &Command{Name: "validate", Run: runValidate}
}

func runValidate(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	file := in.Arg("file")
	if file == "" {
		f, err := composer.ComposerFileName()
		if err != nil {
			return 1, err
		}
		file = f
	}
	st, err := os.Stat(file)
	if err != nil {
		io.WriteError("<error>" + file + " not found.</error>")
		return 3, nil
	}
	if f, err := os.Open(file); err != nil || st.IsDir() {
		io.WriteError("<error>" + file + " is not readable.</error>")
		return 3, nil
	} else {
		f.Close()
	}
	checkAll := validate.CheckAll
	if in.Bool("no-check-all") {
		checkAll = 0
	}
	checkPublish := !in.Bool("no-check-publish")
	checkLock := !in.Bool("no-check-lock")
	checkVersion := validate.CheckVersion
	if in.Bool("no-check-version") {
		checkVersion = 0
	}
	isStrict := in.Bool("strict")
	errs, publishErrors, warnings, fatal := validate.ValidateConfig(file, checkAll, checkVersion, composer.ValidateSchemaFile)
	if fatal != nil {
		return 1, fatal
	}
	c, err := composer.Create(io, composer.Options{DisablePlugins: in.Bool("no-plugins"), DisableScripts: in.Bool("no-scripts"), FullLoad: true, ComposerFile: file})
	if err != nil {
		return 1, err
	}
	checkLock = (checkLock && c.Config.GetBool("lock")) || in.Bool("check-lock")
	var lockErrors []string
	locked, _ := c.Locker.IsLocked()
	if locked && !c.Locker.IsFresh() {
		lockErrors = append(lockErrors, "- The lock file is not up to date with the latest changes in composer.json, it is recommended that you run `composer update` or `composer update <package name>`.")
	}
	if locked {
		info, err := composer.MissingRequirementInfo(c, true)
		if err != nil {
			return 1, err
		}
		lockErrors = append(lockErrors, info...)
	}
	outputValidateResult(io, file, &errs, &warnings, checkPublish, publishErrors, checkLock, lockErrors, true)
	exitCode := 0
	if len(errs) > 0 {
		exitCode = 2
	} else if isStrict && len(warnings) > 0 {
		exitCode = 1
	}
	if in.Bool("with-dependencies") {
		im := composer.NewInstallationManager(c)
		for _, p := range c.LocalRepo.MustPackages() {
			path := im.InstallPath(p)
			if path == "" {
				continue
			}
			depFile := filepath.Join(path, "composer.json")
			if st, err := os.Stat(path); err != nil || !st.IsDir() {
				continue
			}
			if _, err := os.Stat(depFile); err != nil {
				continue
			}
			derrs, dpub, dwarn, fatal := validate.ValidateConfig(path+"/composer.json", checkAll, checkVersion, composer.ValidateSchemaFile)
			if fatal != nil {
				return 1, fatal
			}
			outputValidateResult(io, p.PrettyName, &derrs, &dwarn, checkPublish, dpub, false, nil, false)
			depCode := 0
			if len(derrs) > 0 {
				depCode = 2
			} else if isStrict && len(dwarn) > 0 {
				depCode = 1
			}
			if depCode > exitCode {
				exitCode = depCode
			}
		}
	}
	return exitCode, nil
}

func outputValidateResult(io *cio.IO, name string, errs, warnings *[]string, checkPublish bool, publishErrors []string, checkLock bool, lockErrors []string, printSchemaURL bool) {
	doPrintSchemaURL := false
	switch {
	case len(*errs) > 0:
		io.WriteError("<error>" + name + " is invalid, the following errors/warnings were found:</error>")
	case len(publishErrors) > 0 && checkPublish:
		io.WriteError("<info>" + name + " is valid for simple usage with Composer but has</info>")
		io.WriteError("<info>strict errors that make it unable to be published as a package</info>")
		doPrintSchemaURL = printSchemaURL
	case len(*warnings) > 0:
		io.WriteError("<info>" + name + " is valid, but with a few warnings</info>")
		doPrintSchemaURL = printSchemaURL
	case len(lockErrors) > 0:
		kind := "warnings"
		if checkLock {
			kind = "errors"
		}
		io.Write("<info>" + name + " is valid but your composer.lock has some " + kind + "</info>")
	default:
		io.Write("<info>" + name + " is valid</info>")
	}
	if doPrintSchemaURL {
		io.WriteError("<warning>See https://getcomposer.org/doc/04-schema.md for details on the schema</warning>")
	}
	prefixed := func(list []string, header string) []string {
		out := []string{header}
		for _, e := range list {
			out = append(out, "- "+e)
		}
		return out
	}
	if len(*errs) > 0 {
		*errs = prefixed(*errs, "# General errors")
	}
	if len(*warnings) > 0 {
		*warnings = prefixed(*warnings, "# General warnings")
	}
	var extraWarnings []string
	if len(publishErrors) > 0 && checkPublish {
		*errs = append(*errs, prefixed(publishErrors, "# Publish errors")...)
	}
	if len(lockErrors) > 0 {
		if checkLock {
			*errs = append(*errs, append([]string{"# Lock file errors"}, lockErrors...)...)
		} else {
			extraWarnings = append(extraWarnings, append([]string{"# Lock file warnings"}, lockErrors...)...)
		}
	}
	for _, group := range []struct {
		style string
		msgs  []string
	}{{"error", *errs}, {"warning", append(append([]string(nil), *warnings...), extraWarnings...)}} {
		for _, msg := range group.msgs {
			if strings.HasPrefix(msg, "#") {
				io.WriteError("<" + group.style + ">" + msg + "</" + group.style + ">")
			} else {
				io.WriteError(msg)
			}
		}
	}
}
