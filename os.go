package kvutils

import (
	"os"
)

func DoesFileExist(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsDir(path string) (bool, error) {
	doesFileExist, err := DoesFileExist(path)
	if err != nil {
		return false, err
	}
	if !doesFileExist {
		return false, nil
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
