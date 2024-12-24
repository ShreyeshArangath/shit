package utils

import (
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// IsFile checks if a given path is a file.
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// IsDirEmpty checks if a directory is empty.
func IsDirEmpty(name string) (bool, error) {
	entries, err := os.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
