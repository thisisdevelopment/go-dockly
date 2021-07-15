package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Copy will copy the file from src to dst, the paths have to be absolute to
// ensure consistent behavior.
func Copy(src string, dst string) (err error) {
	if !filepath.IsAbs(src) || !filepath.IsAbs(dst) {
		return fmt.Errorf("can't copy src to dst paths not abosulte between src: %s and dst: %s", src, dst)
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("failed to copy file %s not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "failed to open file to copy")
	}
	defer func() {
		err = srcFile.Close()
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create file to copy to for %s", src)
	}
	defer func() {
		err = dstFile.Close()
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file src: %s dst: %s", src, dstFile.Name())
	}

	return err
}
