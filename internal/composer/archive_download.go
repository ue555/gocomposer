package composer

// DownloadManager::download + install of a package into an arbitrary
// directory, as ArchiveManager::archive does for non-root packages
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

	"gocomposer/internal/fsutil"
	"gocomposer/internal/pkg"
)

// DownloadTo downloads p and installs its files into target. The files
// are first installed next to the vendor directory and then moved, as
// Filesystem::rename falls back to copying across devices (target is
// usually below the system temp directory).
func (im *InstallationManager) DownloadTo(p *pkg.Package, target string) error {
	transfer, err := im.download(p, nil)
	if err != nil {
		return err
	}
	if transfer != nil {
		if err := transfer(); err != nil {
			return err
		}
	}
	im.dropDownloadDir(p)
	var rnd [4]byte
	_, _ = rand.Read(rnd[:])
	stage := im.realVendorDir() + "/composer/archive-" + hex.EncodeToString(rnd[:])
	defer os.RemoveAll(stage)
	switch im.downloaderType(p) {
	case "":
		return nil
	case "path":
		err = im.pathInstall(p, stage, true)
	case "zip", "tar", "gzip":
		err = im.archiveInstall(p, stage, true)
	case "file", "phar":
		err = im.fileInstall(p, stage, true)
	default:
		return &UserError{Msg: "Unknown downloader type: " + im.downloaderType(p), Code: 1}
	}
	if err != nil {
		return err
	}
	return moveTree(stage, target)
}

// moveTree moves the contents of src into dst (rename, or copy+delete).
func moveTree(src, dst string) error {
	if err := os.MkdirAll(dst, 0o777); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		from, to := filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())
		if os.Rename(from, to) == nil {
			continue
		}
		if err := copyTreeAll(from, to); err != nil {
			return err
		}
	}
	return nil
}

func copyTreeAll(src, dst string) error {
	return filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		target := filepath.Join(dst, rel)
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			l, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(l, target)
		case info.IsDir():
			return os.MkdirAll(target, info.Mode().Perm()|0o700)
		default:
			return fsutil.CopyFile(p, target)
		}
	})
}
