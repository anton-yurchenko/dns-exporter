package utils

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/spf13/afero"

	"github.com/pkg/errors"
)

// WriteToFile creates/writes content to file
func WriteToFile(domain string, content string, dir string, fs afero.Fs) (string, error) {
	// format filename
	file := fmt.Sprintf("%s/%s.txt", path.Clean(dir), strings.Replace(strings.TrimSuffix(domain, "."), ".", "-", 1))

	// create file
	out, err := fs.Create(file)
	if err != nil {
		return "", errors.Wrap(err, "error creating file")
	}
	defer out.Close()

	// write content
	_, err = io.WriteString(out, content)
	if err != nil {
		return "", errors.Wrap(err, "error writing to file")
	}

	return file, nil
}

// ValidateDir validates a directory exists, 'create=true' will force the creation.
func ValidateDir(dir string, create bool, fs afero.Fs) (bool, error) {
	b, err := afero.DirExists(fs, dir)
	if err != nil {
		return false, errors.Wrap(err, "error validating directory")
	}

	if !b && create {
		err = fs.MkdirAll(dir, 0777)
		if err != nil {
			return b, errors.Wrap(err, "error creating directory")
		}
	}

	return b, nil
}
