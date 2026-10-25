package cli

// Port of Composer\Command\StatusCommand (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).
// Local changes are detected for git sources (git status/unpushed commits/
// version variations), file based dists (re-download and compare) and path
// packages and the svn, hg, fossil and perforce source checkouts.

import (
	"fmt"
	"os"
	"strings"

	"gocomposer/internal/composer"
)

func init() { registerCommand(statusCommand) }

func statusCommand() *Command {
	return &Command{Name: "status", Run: runStatus}
}

const (
	statusExitErrors         = 1
	statusExitUnpushed       = 2
	statusExitVersionChanges = 4
)

type statusEntry struct {
	path, text string
}

type statusList struct{ entries []statusEntry }

// set is $array[$path] = $text (keeps the first insertion position).
func (l *statusList) set(path, text string) {
	for i := range l.entries {
		if l.entries[i].path == path {
			l.entries[i].text = text
			return
		}
	}
	l.entries = append(l.entries, statusEntry{path, text})
}

func runStatus(ctx *Context) (int, error) {
	c, err := requireComposer(ctx)
	if err != nil {
		return exitFor(err)
	}
	events := composer.NewEventDispatcher(c)
	events.RunScripts = !c.DisableScripts
	if _, err := events.DispatchScript("pre-status-cmd", true, nil, nil); err != nil {
		return scriptExit(err), err
	}
	code, err := statusExecute(ctx, c)
	if err != nil {
		return code, err
	}
	if _, err := events.DispatchScript("post-status-cmd", true, nil, nil); err != nil {
		return scriptExit(err), err
	}
	return code, nil
}

func indentChanges(changes string) string {
	lines := strings.Split(changes, "\n")
	for i, l := range lines {
		lines[i] = "    " + strings.TrimLeft(l, " \t\n\r\x00\x0B")
	}
	return strings.Join(lines, "\n")
}

func statusExecute(ctx *Context, c *composer.Composer) (int, error) {
	io := ctx.IO
	verbose := io.IsVerbose()
	im := composer.NewInstallationManager(c)
	var errs, unpushed statusList
	type versionChange struct {
		path                                     string
		prevVersion, prevRef, curVersion, curRef string
	}
	var versionChanges []versionChange

	for _, p := range c.LocalRepo.CanonicalPackages() {
		if p.Type() == "metapackage" {
			continue
		}
		src := p.InstallationSource()
		if src != "dist" && src != "source" {
			return 1, &composer.UserError{Msg: "Package " + p.String() + " does not have an installation source set", Code: 1}
		}
		typ := im.DownloaderType(p)
		targetDir := im.InstallPath(p)
		if targetDir == "" {
			continue
		}
		switch typ {
		case "git", "svn", "hg", "fossil", "perforce", "path", "zip", "tar", "gzip", "file", "phar", "xz", "rar":
		default:
			return 1, &composer.UserError{Msg: "Unknown downloader type: " + typ + ". Available types: git, svn, fossil, hg, perforce, zip, rar, tar, gzip, xz, phar, file, path.", Code: 1}
		}
		// ChangeReportInterface
		if st, err := os.Lstat(targetDir); err == nil && st.Mode()&os.ModeSymlink != 0 {
			errs.set(targetDir, targetDir+" is a symbolic link.")
		}
		var changes string
		var err error
		switch typ {
		case "git", "svn", "hg", "fossil", "perforce":
			changes, err = im.LocalChanges(p)
		default:
			changes, err = im.FileLocalChanges(p, targetDir)
		}
		if err != nil {
			return 1, err
		}
		if changes != "" {
			errs.set(targetDir, changes)
		}
		// VcsCapableDownloaderInterface
		if typ == "git" || typ == "path" {
			if _, commit, ok := composer.GuessVersionCommit(targetDir, p, c); ok && commit != "" {
				var prevRef string
				switch src {
				case "source":
					prevRef = p.SourceReference()
				case "dist":
					prevRef = p.DistReference()
				}
				curPretty, curCommit, ok := composer.GuessVersionCommit(targetDir, p, c)
				if prevRef != "" && ok && curCommit != prevRef && curPretty != prevRef {
					versionChanges = append(versionChanges, versionChange{targetDir, p.PrettyVersion(), prevRef, curPretty, curCommit})
				}
			}
		}
		// DvcsDownloaderInterface
		if typ == "git" {
			u, err := im.UnpushedChanges(p)
			if err != nil {
				return 1, err
			}
			if u != "" {
				unpushed.set(targetDir, u)
			}
		}
	}

	if len(errs.entries) == 0 && len(unpushed.entries) == 0 && len(versionChanges) == 0 {
		io.WriteError("<info>No local changes</info>")
		return 0, nil
	}
	writeList := func(l statusList) {
		for _, e := range l.entries {
			if verbose {
				io.Write("<info>" + e.path + "</info>:")
				io.Write(indentChanges(e.text))
			} else {
				io.Write(e.path)
			}
		}
	}
	if len(errs.entries) > 0 {
		io.WriteError("<error>You have changes in the following dependencies:</error>")
		writeList(errs)
	}
	if len(unpushed.entries) > 0 {
		io.WriteError("<warning>You have unpushed changes on the current branch in the following dependencies:</warning>")
		writeList(unpushed)
	}
	if len(versionChanges) > 0 {
		io.WriteError("<warning>You have version variations in the following dependencies:</warning>")
		for _, vc := range versionChanges {
			if verbose {
				cur, prev := vc.curVersion, vc.prevVersion
				if cur == "" {
					cur = vc.curRef
				}
				if prev == "" {
					prev = vc.prevRef
				}
				if io.IsVeryVerbose() {
					cur += fmt.Sprintf(" (%s)", vc.curRef)
					prev += fmt.Sprintf(" (%s)", vc.prevRef)
				}
				io.Write("<info>" + vc.path + "</info>:")
				io.Write(fmt.Sprintf("    From <comment>%s</comment> to <comment>%s</comment>", prev, cur))
			} else {
				io.Write(vc.path)
			}
		}
	}
	if !verbose {
		io.WriteError("Use --verbose (-v) to see a list of files")
	}
	code := 0
	if len(errs.entries) > 0 {
		code += statusExitErrors
	}
	if len(unpushed.entries) > 0 {
		code += statusExitUnpushed
	}
	if len(versionChanges) > 0 {
		code += statusExitVersionChanges
	}
	return code, nil
}
