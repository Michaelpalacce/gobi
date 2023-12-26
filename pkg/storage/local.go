package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

type LocalDriver struct {
	VaultsPath string
	queue      []models.Item
}

func NewLocalDriver(vaultsPath string) *LocalDriver {
	return &LocalDriver{
		VaultsPath: vaultsPath,
		queue:      make([]models.Item, 0),
	}
}

// Enqueue adds the given items array to the queue for later processing.
// Will not add items that are already in the local storage, based on filePath and SHA256
func (d *LocalDriver) Enqueue(items []models.Item) {
	for _, item := range items {
		if ok := d.checkIfLocalMatch(item); !ok {
			d.queue = append(d.queue, item)
		}
	}
}

// GetNext will return the next File in the queue
func (d *LocalDriver) GetNext() *models.Item {
	if len(d.queue) > 0 {
		var current models.Item
		slog.Debug("Queue", "queue", d.queue)
		current, d.queue = d.queue[0], d.queue[1:]
		return &current
	}

	return nil
}

// CheckIfLocalMatch will build up the correct filePath based on the item and check if what we have locally matches.
// Checks by filePath and SHA256
func (d *LocalDriver) checkIfLocalMatch(i models.Item) bool {
	absFilePath := d.getFilePath(i)
	_, err := os.Stat(absFilePath)
	if err != nil {
		return false
	}

	sha256, err := digest.FileSHA256(absFilePath)
	if err != nil {
		return false
	}

	return sha256 == i.SHA256
}

// HasItemsToProcess will return true if there is more than one item in the queue
func (d *LocalDriver) HasItemsToProcess() bool {
	return len(d.queue) > 0
}

// GetReader should be used to get a reader for the given item, when you want to send it
func (d *LocalDriver) GetReader(i models.Item) (io.ReadCloser, error) {
	file, err := os.Open(d.getFilePath(i))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	return file, nil
}

// GetWriter should be used to get a writer for the given item, when you want to save it
func (d *LocalDriver) GetWriter(i models.Item) (io.WriteCloser, error) {
	file, err := os.OpenFile(d.getFilePath(i), os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	return file, nil
}

func (d *LocalDriver) Exists(i models.Item) bool {
	_, err := os.Stat(d.getFilePath(i))
	return err == nil
}

// getFilePath will return the absolute path to the file
func (d *LocalDriver) getFilePath(i models.Item) string {
	return filepath.Join(d.VaultsPath, i.VaultName, i.ServerPath)
}
