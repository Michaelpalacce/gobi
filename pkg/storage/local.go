package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

func (d *LocalDriver) GetAllItems() []models.Item {
	queue := d.queue
	d.queue = make([]models.Item, 0)

	return queue
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
	path := d.getFilePath(i)
	dirPath := filepath.Dir(path)
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %s", err)
	}

	// Use os.MkdirAll to create the directory and its parents
	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %s", err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o666)
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

// CalculateSHA256 will return the SHA256 of the given item
func (d *LocalDriver) CalculateSHA256(i models.Item) string {
	digest, err := digest.FileSHA256(d.getFilePath(i))
	if err != nil {
		return ""
	}

	return digest
}

// EnqueueItemsSince will add all items that have been modified since the given lastSyncTime to the queue
func (d *LocalDriver) EnqueueItemsSince(lastSyncTime int, vaultName string) {
	vaultPath, err := filepath.Abs(filepath.Join(d.VaultsPath, vaultName))
	if err != nil {
		slog.Error("Error getting vault path", "error", err)
		return
	}

	filepath.WalkDir(vaultPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.ModTime().Unix() < int64(lastSyncTime) {
			return nil
		}

		item := models.Item{
			VaultName:  vaultName,
			ServerPath: strings.Replace(path, vaultPath+"/", "", 1),
			Size:       int(fileInfo.Size()),
		}

		item.SHA256 = d.CalculateSHA256(item)

		d.queue = append(d.queue, item)
		return nil
	})
	slog.Debug("EnqueueItemsSince", "queue", d.queue)
}
