// Package file provides functions for file operations.
package file

import (
	"io/fs"
	"path/filepath"

	"github.com/nao1215/spare/utils/errfmt"
)

// WalkDir returns a list of files in the specified directory.
func WalkDir(rootDir string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errfmt.Wrap(err, "failed to walk directory")
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}
