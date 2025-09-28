package fileutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func EnsureDir(dir string) error {
	exists, err := CheckDir(dir)
	if err != nil {
		return fmt.Errorf("checking the dir: %v", err)
	}
	if !exists {
		innerErr := os.Mkdir(dir, 0750)
		if innerErr != nil {
			return fmt.Errorf("creating the dir: %v", innerErr)
		}
	}
	return nil
}

func CheckDir(dir string) (exists bool, err error) {
	_, err = os.Stat(dir)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}
