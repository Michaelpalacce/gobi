package iops

import (
	"fmt"
	"path/filepath"
	"strings"
)

// JoinSafe joins the provided path segments and ensures the resulting path is within the base directory.
// It returns the absolute path if successful.
func JoinSafe(path ...string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("no path provided")
	}

	base := path[0]

	joinedPath := filepath.Join(path...)

	// Ensure the cleaned path is still within the base directory
	if !strings.HasPrefix(joinedPath, base) {
		return "", fmt.Errorf("path is outside of base directory")
	}

	absPath, err := filepath.Abs(joinedPath)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %w", err)
	}

	return absPath, nil
}
