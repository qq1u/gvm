package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func CurrentDir() (string, error) {
	current, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("excutable failed: %w", err)
	}

	return filepath.Dir(current), nil
}
