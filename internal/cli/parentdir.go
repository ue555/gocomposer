package cli

// Port of the "use the composer.json of a parent directory" prompt of
// Composer\Console\Application::doRunComposer (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann,
// Jordi Boggiano) and Application::getUseParentDirConfigValue.

import (
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/cio"
	"gocomposer/internal/composer"
	"gocomposer/internal/config"
	"gocomposer/internal/fsutil"
)

var parentDirExempt = map[string]bool{"": true, "list": true, "init": true, "about": true, "help": true, "diagnose": true, "self-update": true, "global": true, "create-project": true, "outdated": true}

// useParentDir switches to a parent directory holding a composer.json when
// the current one has none (config use-parent-dir: prompt/true).
func (a *App) useParentDir(io *cio.IO, argv []string, rawCmdName string) error {
	if workingDirArg(argv) != "" {
		return nil
	}
	cmdName := ""
	if rawCmdName != "" {
		c, err := a.findCommand(rawCmdName)
		if err == nil {
			cmdName = c.Name
		} else if strings.Contains(err.Error(), "is ambiguous") {
			// \InvalidArgumentException: the name stays ''
			cmdName = ""
		} else {
			cmdName = "\x00not-found"
		}
	}
	if parentDirExempt[cmdName] {
		return nil
	}
	file, err := composer.ComposerFileName()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(file); err == nil {
		return nil
	}
	for _, t := range argv {
		if t == "--" {
			break
		}
		if t == "--help" || t == "-h" || (cmdName == "config" && (t == "--file" || t == "-f")) {
			return nil
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	cfg, err := config.CreateConfig(cwd, nil)
	if err != nil {
		return nil
	}
	setting := cfg.Get("use-parent-dir")
	if b, ok := setting.(bool); ok && !b {
		return nil
	}
	always := false
	if b, ok := setting.(bool); ok && b {
		always = true
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		home = "/"
	}
	home = fsutil.Realpath(home)
	dir := filepath.Dir(cwd)
	for filepath.Dir(dir) != dir && dir != home {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			if !always && !io.Interactive {
				io.WriteError("<info>No composer.json in current directory, to use the one at " + dir + " run interactively or set config.use-parent-dir to true</info>")
				break
			}
			use := always
			if !use {
				if use, err = io.AskConfirmation("<info>No composer.json in current directory, do you want to use the one at "+dir+"?</info> [<comment>y,n</comment>]? ", true); err != nil {
					return err
				}
			}
			if use {
				if always {
					io.WriteError("<info>No composer.json in current directory, changing working directory to " + dir + "</info>")
				} else {
					io.WriteError(`<info>Always want to use the parent dir? Use "composer config --global use-parent-dir true" to change the default.</info>`)
				}
				_ = os.Chdir(dir)
			}
			break
		}
		dir = filepath.Dir(dir)
	}
	return nil
}
