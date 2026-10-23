// Non-unix helpers for the Filesystem::safeFilePutContents / is_writable port
// (src/Composer/Util/Filesystem.php, composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

//go:build !unix

package configsource

import "os"

// isWritable approximates PHP is_writable() by opening the file for writing.
func isWritable(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	if st.IsDir() {
		return st.Mode().Perm()&0o200 != 0
	}
	fh, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	fh.Close()
	return true
}

func keepOwner(*os.File, string, os.FileInfo) bool { return true }
