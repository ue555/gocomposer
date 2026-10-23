package composer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/pkg"
)

// Port of Composer\Installer\BinaryInstaller (unix proxies; bin-compat
// "full" additionally writes .bat proxies).

var (
	reUnsafeProxy = regexp.MustCompile("[*$`\"&^|<>()%!;\x00-\x1f\x7f]")
	reShebangLine = regexp.MustCompile(`^#!(?:/[a-zA-Z0-9_.+-]+)*/([a-zA-Z0-9_.+][a-zA-Z0-9_.+-]*(?:[ \t].*)?)$`)
	reEnvPrefix   = regexp.MustCompile(`^env[ \t]+(?:-S[ \t]*)?`)
	rePHPInterp   = regexp.MustCompile(`^php[0-9.]*(?:[ \t]|$)`)
	rePHPBinStart = regexp.MustCompile(`^(#!.*\r?\n)?[\r\n\t ]*<\?php`)
	reBatchUnsafe = regexp.MustCompile("[\"\r\n\x1a]")
)

func isSafeBinPath(bin string) bool { return !reUnsafeProxy.MatchString(bin) }

func shebangCaller(shebang string) (string, bool) {
	m := reShebangLine.FindStringSubmatch(shebang)
	if m == nil || reUnsafeProxy.MatchString(m[1]) {
		return "", false
	}
	caller := reEnvPrefix.ReplaceAllString(m[1], "")
	if caller == "" || caller == "env" {
		return "", false
	}
	if caller[0] == '-' {
		return "env " + caller, true
	}
	return caller, true
}

func binInsidePackage(installPath, binPath string) bool {
	rb := fsutil.Realpath(binPath)
	ri := fsutil.Realpath(installPath)
	if rb == "" || ri == "" {
		return false
	}
	return strings.HasPrefix(rb, ri+string(filepath.Separator))
}

func escapeShellArg(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func (im *InstallationManager) realBinDir() string {
	_ = os.MkdirAll(im.binDir, 0o777)
	if r := fsutil.Realpath(im.binDir); r != "" {
		return r
	}
	return im.binDir
}

func (im *InstallationManager) installBinaries(p *pkg.Package, installPath string, warnOnOverwrite bool) {
	bins := p.Binaries()
	if len(bins) == 0 {
		return
	}
	io := im.c.IO
	for _, bin := range bins {
		if !isSafeBinPath(bin) {
			io.WriteError("    <warning>Skipped installation of bin " + bin + " for package " + p.Name + ": the bin path contains invalid characters</warning>")
			continue
		}
		binPath := installPath + "/" + bin
		st, err := os.Stat(binPath)
		if err != nil {
			io.WriteError("    <warning>Skipped installation of bin " + bin + " for package " + p.Name + ": file not found in package</warning>")
			continue
		}
		if st.IsDir() {
			io.WriteError("    <warning>Skipped installation of bin " + bin + " for package " + p.Name + ": found a directory at that path</warning>")
			continue
		}
		if !binInsidePackage(installPath, binPath) {
			io.WriteError("    <warning>Skipped installation of bin " + bin + " for package " + p.Name + ": the bin resolves to a path outside of the package directory</warning>")
			continue
		}
		if !fsutil.IsAbsolutePath(binPath) {
			binPath = fsutil.Realpath(binPath)
		}
		binDir := im.realBinDir()
		link := binDir + "/" + filepath.Base(bin)
		if lst, err := os.Lstat(link); err == nil {
			if lst.Mode()&os.ModeSymlink == 0 {
				if _, err := os.Stat(link); err == nil {
					if warnOnOverwrite {
						io.WriteError("    Skipped installation of bin " + bin + " for package " + p.Name + ": name conflicts with an existing file")
					}
					continue
				}
			}
			if fsutil.Realpath(link) == fsutil.Realpath(binPath) {
				os.Remove(link)
			}
		}
		compat := im.c.Config.GetString("bin-compat")
		if compat == "full" {
			if reBatchUnsafe.MatchString(binPath) || reBatchUnsafe.MatchString(binDir) {
				io.WriteError("    <warning>Skipped installation of bin " + bin + " for package " + p.Name + ": the path contains a double quote or line break, which cannot be used in a Windows bin proxy</warning>")
				continue
			}
			im.installUnixyProxy(binPath, link)
			if !strings.HasSuffix(binPath, ".bat") {
				if _, err := os.Stat(link + ".bat"); err == nil {
					io.WriteError("    Skipped installation of bin " + bin + ".bat proxy for package " + p.Name + ": a .bat proxy was already installed")
				} else {
					_ = os.WriteFile(link+".bat", []byte(im.windowsProxy(binPath, link+".bat")), 0o666)
					_ = os.Chmod(link+".bat", 0o777&^currentUmask())
				}
			}
		} else {
			im.installUnixyProxy(binPath, link)
		}
		_ = os.Chmod(binPath, 0o777&^currentUmask())
	}
}

func (im *InstallationManager) installUnixyProxy(binPath, link string) {
	code := im.unixyProxyCode(binPath, link)
	_ = os.Remove(link)
	_ = os.WriteFile(link, []byte(code), 0o666)
	_ = os.Chmod(link, 0o777&^currentUmask())
}

func (im *InstallationManager) removeBinaries(p *pkg.Package) {
	binDir := im.realBinDir()
	bins := p.Binaries()
	if len(bins) == 0 {
		return
	}
	for _, bin := range bins {
		link := binDir + "/" + filepath.Base(bin)
		if _, err := os.Lstat(link); err == nil {
			os.Remove(link)
		}
		if st, err := os.Stat(link + ".bat"); err == nil && st.Mode().IsRegular() {
			os.Remove(link + ".bat")
		}
	}
	if isDir(binDir) && fsutil.IsDirEmpty(binDir) {
		os.Remove(binDir)
	}
}

func determineBinaryCaller(bin string) string {
	if strings.HasSuffix(bin, ".bat") || strings.HasSuffix(bin, ".exe") {
		return "call"
	}
	data, err := os.ReadFile(bin)
	if err != nil {
		return "php"
	}
	line := string(data)
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = line[:i+1]
	}
	shebang := strings.TrimRight(line, " \t\r\n")
	if c, ok := shebangCaller(shebang); ok {
		return c
	}
	return "php"
}

func (im *InstallationManager) windowsProxy(bin, link string) string {
	caller := determineBinaryCaller(bin)
	var target string
	if rePHPInterp.MatchString(caller) {
		target = strings.TrimSuffix(filepath.Base(link), ".bat")
	} else {
		target, _ = fsutil.FindShortestPath(link, bin, false, false)
	}
	return "@ECHO OFF\r\n" +
		"setlocal DISABLEDELAYEDEXPANSION\r\n" +
		"SET \"BIN_TARGET=%~dp0/" + strings.ReplaceAll(target, "%", "%%") + "\"\r\n" +
		"SET \"COMPOSER_RUNTIME_BIN_DIR=%~dp0\"\r\n" +
		caller + " \"%BIN_TARGET%\" %*\r\n"
}

func (im *InstallationManager) unixyProxyCode(bin, link string) string {
	binPath, _ := fsutil.FindShortestPath(link, bin, false, false)
	binDir := escapeShellArg(filepath.Dir(binPath))
	binFile := escapeShellArg(filepath.Base(binPath))
	contents, _ := readPrefix(bin, 500)
	if m := rePHPBinStart.FindStringSubmatch(contents); m != nil {
		shebang := ""
		if m[1] != "" {
			shebang = strings.TrimRight(m[1], " \t\r\n")
		}
		proxyCode := "#!/usr/bin/env php"
		if c, ok := shebangCaller(shebang); ok && rePHPInterp.MatchString(c) {
			proxyCode = shebang
		}
		binPathExported, _ := fsutil.FindShortestPathCode(link, bin, false, true, false)
		binPathComment := strings.ReplaceAll(strings.NewReplacer("\r", " ", "\n", " ").Replace(binPath), "*/", "* /")
		streamProxyCode, streamHint := "", ""
		globalsCode := "$GLOBALS['_composer_bin_dir'] = __DIR__;\n"
		phpunitHack1, phpunitHack2 := "", ""
		vendorReal := fsutil.Realpath(im.vendorDir)
		if vendorReal == "" {
			vendorReal = im.vendorDir
		}
		autoloadCode, _ := fsutil.FindShortestPathCode(link, vendorReal+"/autoload.php", false, true, false)
		globalsCode += "$GLOBALS['_composer_autoload_path'] = " + autoloadCode + ";\n"
		if fsutil.NormalizePath(bin) == fsutil.NormalizePath(im.vendorDir+"/phpunit/phpunit/phpunit") {
			globalsCode += "$GLOBALS['__PHPUNIT_ISOLATION_EXCLUDE_LIST'] = $GLOBALS['__PHPUNIT_ISOLATION_BLACKLIST'] = array(realpath(" + binPathExported + "));\n"
			phpunitHack1 = "'phpvfscomposer://'."
			phpunitHack2 = "\n                $data = str_replace('__DIR__', var_export(dirname($this->realpath), true), $data);\n                $data = str_replace('__FILE__', var_export($this->realpath, true), $data);"
		}
		if strings.TrimSpace(m[0]) != "<?php" {
			streamHint = " using a stream wrapper to prevent the shebang from being output on PHP<8\n *"
			streamProxyCode = strings.NewReplacer("{{HACK1}}", phpunitHack1, "{{HACK2}}", phpunitHack2, "{{BINPATH}}", binPathExported).Replace(streamProxyTemplate)
		}
		return proxyCode + "\n" + "<?php\n\n/**\n * Proxy PHP file generated by Composer\n *\n * This file includes the referenced bin path (" + binPathComment + ")\n *" + streamHint + "\n * @generated\n */\n\nnamespace Composer;\n\n" + globalsCode + "\n" + streamProxyCode + "\nreturn include " + binPathExported + ";\n"
	}
	return strings.NewReplacer("{{BINDIR}}", binDir, "{{BINFILE}}", binFile).Replace(shellProxyTemplate)
}

func readPrefix(file string, n int) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, n)
	k, _ := f.Read(buf)
	return string(buf[:k]), nil
}

const streamProxyTemplate = `if (PHP_VERSION_ID < 80000) {
    if (!class_exists('Composer\BinProxyWrapper')) {
        /**
         * @internal
         */
        final class BinProxyWrapper
        {
            private $handle;
            private $position;
            private $realpath;

            public function stream_open($path, $mode, $options, &$opened_path)
            {
                // get rid of phpvfscomposer:// prefix for __FILE__ & __DIR__ resolution
                $opened_path = substr($path, 17);
                $this->realpath = realpath($opened_path) ?: $opened_path;
                $opened_path = {{HACK1}}$this->realpath;
                $this->handle = fopen($this->realpath, $mode);
                $this->position = 0;

                return (bool) $this->handle;
            }

            public function stream_read($count)
            {
                $data = fread($this->handle, $count);

                if ($this->position === 0) {
                    $data = preg_replace('{^#!.*\r?\n}', '', $data);
                }{{HACK2}}

                $this->position += strlen($data);

                return $data;
            }

            public function stream_cast($castAs)
            {
                return $this->handle;
            }

            public function stream_close()
            {
                fclose($this->handle);
            }

            public function stream_lock($operation)
            {
                return $operation ? flock($this->handle, $operation) : true;
            }

            public function stream_seek($offset, $whence)
            {
                if (0 === fseek($this->handle, $offset, $whence)) {
                    $this->position = ftell($this->handle);
                    return true;
                }

                return false;
            }

            public function stream_tell()
            {
                return $this->position;
            }

            public function stream_eof()
            {
                return feof($this->handle);
            }

            public function stream_stat()
            {
                return array();
            }

            public function stream_set_option($option, $arg1, $arg2)
            {
                return true;
            }

            public function url_stat($path, $flags)
            {
                $path = substr($path, 17);
                if (file_exists($path)) {
                    return stat($path);
                }

                return false;
            }
        }
    }

    if (
        (function_exists('stream_get_wrappers') && in_array('phpvfscomposer', stream_get_wrappers(), true))
        || (function_exists('stream_wrapper_register') && stream_wrapper_register('phpvfscomposer', 'Composer\BinProxyWrapper'))
    ) {
        return include("phpvfscomposer://" . {{BINPATH}});
    }
}
`

const shellProxyTemplate = `#!/usr/bin/env sh

# Support bash to support ` + "`source`" + ` with fallback on $0 if this does not run with bash
# https://stackoverflow.com/a/35006505/6512
selfArg="$BASH_SOURCE"
if [ -z "$selfArg" ]; then
    selfArg="$0"
fi

self=$(realpath "$selfArg" 2> /dev/null)
if [ -z "$self" ]; then
    self="$selfArg"
fi

dir=$(cd "${self%[/\\]*}" > /dev/null; cd {{BINDIR}} && pwd)

if [ -d /proc/cygdrive ]; then
    case $(which php) in
        $(readlink -n /proc/cygdrive)/*)
            # We are in Cygwin using Windows php, so the path must be translated
            dir=$(cygpath -m "$dir");
            ;;
    esac
fi

export COMPOSER_RUNTIME_BIN_DIR="$(cd "${self%[/\\]*}" > /dev/null; pwd)"

# If bash is sourcing this file, we have to source the target as well
bashSource="$BASH_SOURCE"
if [ -n "$bashSource" ]; then
    if [ "$bashSource" != "$0" ]; then
        source "${dir}/"{{BINFILE}} "$@"
        return
    fi
fi

exec "${dir}/"{{BINFILE}} "$@"
`
