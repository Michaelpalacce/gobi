package storage

import "github.com/Michaelpalacce/gobi/pkg/models"

type LocalDriver struct {
	VaultPath string
}

// CheckIfLocalMatch will build up the correct filePath based on the item and check if what we have locally matches.
// Checks by filePath and SHA256
func (d *LocalDriver) CheckIfLocalMatch(i models.Item) bool {
	return false
}
