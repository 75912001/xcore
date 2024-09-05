package util

import (
	"os"
	"path/filepath"
)

func GetCurrentPath() (currentPath string, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
