package util

import (
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Removes(paths ...string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}
