package fileutil

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Specify the number of bytes you want to be different when create partial
// matching files.
const numOfDiffBytes = 4

// Create will create a file in a specific dir, for the specific size. The
// data inside of file is completely random.
func Create(dir string, size int) (fileName string, err error) {
	tmpFile, err := ioutil.TempFile(dir, "")
	if err != nil {
		return "", errors.Wrapf(err, "failed to create tmp file inside of %s", dir)
	}
	defer func() {
		err = tmpFile.Close()
	}()

	var fileContent = make([]byte, size)
	_, err = rand.Read(fileContent)
	if err != nil {
		return "", errors.Wrap(err, "failed to create random data for file")
	}

	_, err = tmpFile.Write(fileContent)
	if err != nil {
		return "", errors.Wrapf(err, "failed to write random data to file %s", tmpFile.Name())
	}

	return tmpFile.Name(), err
}

// useful for fuzzing

// CreatePartialMatch will create two files that have the same size and the
// first few bytes but the final 4 bytes of the file are different. The 4 extra
// bytes are included into the specified size.
func CreatePartialMatch(dir string, size int) ([]string, error) {
	if !filepath.IsAbs(dir) {
		return nil, fmt.Errorf("cannot append to file, path %s is not absolute", dir)
	}

	// Create the identical files.
	originalFile, err := Create(dir, size-numOfDiffBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create file inside %s", dir)
	}

	var cpFile = fmt.Sprintf("%s_partial", originalFile)

	err = Copy(originalFile, cpFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to copy file to %s", dir)
	}

	err = appendToFile(originalFile, numOfDiffBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to append random data to %s", originalFile)
	}

	err = appendToFile(cpFile, numOfDiffBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to append random data to %s", originalFile)
	}

	return []string{originalFile, cpFile}, nil
}

func appendToFile(path string, size int) (err error) {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("cannot append to file, path is not absolute: %s", path)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", path)
	}
	defer func() {
		err = f.Close()
	}()

	var appendData = make([]byte, size)

	_, err = rand.Read(appendData)
	if err != nil {
		return errors.Wrapf(err, "failed to create random data for file")
	}

	if _, err := f.Write(appendData); err != nil {
		return errors.Wrapf(err, "failed to append data to file %s", path)
	}

	return err
}
