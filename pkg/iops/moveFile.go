package iops

import (
	"io"
	"os"
	"path/filepath"
)

// MoveFile was implemented because the default os.Rename does not work cross device volumes
// It will create the destination file structure if it does not exist
func MoveFile(src, dst string) error {
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Optional: You can remove the source file after copying
	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}
