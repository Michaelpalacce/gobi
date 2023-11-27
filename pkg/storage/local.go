package storage

import (
	"os"
	"path/filepath"

	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

type LocalDriver struct {
	VaultPath string
}

// CheckIfLocalMatch will build up the correct filePath based on the item and check if what we have locally matches.
// Checks by filePath and SHA256
func (d *LocalDriver) CheckIfLocalMatch(i models.Item) bool {
	absFilePath := filepath.Join(d.VaultPath, i.ServerPath)
	_, err := os.Stat(absFilePath)

	// This can be returned directly, yes, leave it for now
	if err != nil {
		return false
	}

	sha256, err := digest.FileSHA256(absFilePath)

	if err != nil {
		return false
	}

	return sha256 == i.SHA256
}
