package storage

import (
	"os"
	"path/filepath"

	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

type LocalDriver struct {
	VaultPath string
	queue     []models.Item
}

// Enqueue adds the given items array to the queue for later processing.
func (d *LocalDriver) Enqueue(items []models.Item) error {
	if d.queue == nil {
		d.queue = make([]models.Item, len(items))
	}

	for _, item := range items {
		if ok := d.checkIfLocalMatch(item); !ok {
			d.queue = append(d.queue, item)
		}
	}

	return nil
}

// CheckIfLocalMatch will build up the correct filePath based on the item and check if what we have locally matches.
// Checks by filePath and SHA256
func (d *LocalDriver) checkIfLocalMatch(i models.Item) bool {
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

// TODO: finish me
func (d *LocalDriver) HasItemsToProcess() bool {
	return false
}
