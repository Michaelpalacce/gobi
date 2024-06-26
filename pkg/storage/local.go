package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/Michaelpalacce/gobi/pkg/iops"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/fsnotify/fsnotify"
)

// @TODO: IMPORTANT! Make sure we don't allow any paths that are outside of the vaults path and that we don't allow any paths that are not in the vault path

var localVaultsLocation = os.Getenv("LOCAL_VAULTS_LOCATION")

// LocalDriver is a storage driver that stores files locally on the disk.
type LocalDriver struct {
	VaultPath string
	queue     []models.Item
	conflicts []models.Item
}

// NewLocalDriver creates a new LocalDriver for the given Vault
func NewLocalDriver(vaultName string) (*LocalDriver, error) {
	path, err := iops.JoinSafe(localVaultsLocation, vaultName)
	slog.Debug("LocalDriver", "path", path)

	if err != nil {
		return nil, fmt.Errorf("error getting vault path: %w", err)
	}

	storageDriver := &LocalDriver{
		VaultPath: path,
		queue:     make([]models.Item, 0),
		conflicts: make([]models.Item, 0),
	}

	err = storageDriver.EnsureVault()
	if err != nil {
		return nil, fmt.Errorf("error ensuring that the vault exists: %w", err)
	}

	return storageDriver, nil
}

// EnsureVault will make sure that the vault exists on the disk
func (d *LocalDriver) EnsureVault() error {
	_, err := os.Stat(d.VaultPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(d.VaultPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating vault: %w", err)
		}

		slog.Info("Created vault", "vault", d.VaultPath)
	}

	return nil
}

// Enqueue adds the given items array to the queue for later processing.
// Will not add items that are already in the local storage, based on filePath and SHA256
func (d *LocalDriver) Enqueue(items []models.Item) {
	for _, item := range items {
		if ok := d.checkIfLocalMatch(item); !ok {
			fileInfo, err := os.Stat(d.getFilePath(item))
			if err == nil && fileInfo.ModTime().Unix() > item.ServerMTime {
				d.conflicts = append(d.conflicts, item)
				continue
			}

			d.queue = append(d.queue, item)
		}
	}
}

func (d *LocalDriver) GetMTime(i models.Item) int64 {
	fileInfo, err := os.Stat(d.getFilePath(i))
	if err != nil {
		return 0
	}

	return fileInfo.ModTime().Unix()
}

// GetNext will return the next File in the queue
func (d *LocalDriver) GetNext(conflictMode bool) *models.Item {
	if conflictMode && len(d.conflicts) > 0 {
		var current models.Item
		current, d.conflicts = d.conflicts[0], d.conflicts[1:]
		return &current
	}

	if !conflictMode && len(d.queue) > 0 {
		var current models.Item
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
func (d *LocalDriver) HasItemsToProcess(conflictMode bool) bool {
	if conflictMode {
		return len(d.conflicts) > 0
	}

	return len(d.queue) > 0
}

func (d *LocalDriver) GetAllItems(conflictMode bool) []models.Item {
	if conflictMode {
		queue := d.conflicts
		d.conflicts = make([]models.Item, 0)
		return queue
	}

	queue := d.queue
	d.queue = make([]models.Item, 0)

	return queue
}

// GetReader opens a file in the local storage and returns a reader for it.
// It returns an error if the file cannot be opened.
// The caller is responsible for closing the reader.
func (d *LocalDriver) GetReader(i models.Item) (io.ReadCloser, error) {
	file, err := os.Open(d.getFilePath(i))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return file, nil
}

// Touch will update the mtime of the given item to the server mtime
func (d *LocalDriver) Touch(i models.Item) error {
	// int64 to time.Time
	t := time.Unix(i.ServerMTime, 0)

	return os.Chtimes(d.getFilePath(i), t, t)
}

// GetWriter should be used to get a writer for the given item, when you want to save it
func (d *LocalDriver) GetWriter(i models.Item) (io.WriteCloser, error) {
	path := d.getFilePath(i)
	dirPath := filepath.Dir(path)
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	// Use os.MkdirAll to create the directory and its parents
	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return file, nil
}

func (d *LocalDriver) Exists(i models.Item) bool {
	_, err := os.Stat(d.getFilePath(i))
	return err == nil
}

// getFilePath will return the absolute path to the file
func (d *LocalDriver) getFilePath(i models.Item) string {
	return filepath.Join(d.VaultPath, i.ServerPath)
}

// CalculateSHA256 will return the SHA256 of the given item
func (d *LocalDriver) CalculateSHA256(i models.Item) string {
	digest, err := digest.FileSHA256(d.getFilePath(i))
	if err != nil {
		slog.Error("Error calculating SHA256", "error", err)
		return ""
	}

	return digest
}

// EnqueueItemsSince will add all items that have been modified since the given lastSyncTime to the queue
func (d *LocalDriver) EnqueueItemsSince(lastSyncTime int, vaultName string) {
	vaultPath, err := filepath.Abs(d.VaultPath)
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
			ServerPath:  strings.Replace(path, vaultPath+"/", "", 1),
			Size:        int(fileInfo.Size()),
			ServerMTime: fileInfo.ModTime().Unix(),
		}

		item.SHA256 = d.CalculateSHA256(item)

		d.queue = append(d.queue, item)
		return nil
	})
}

// WatchVault will watch the given vault for changes and add them to the queue
// @TODO: Deletions. Send a message to the server that the file was deleted
func (d *LocalDriver) WatchVault(vaultName string, changeChan chan<- *models.Item) error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {

					if err != nil {
						return
					}

					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						return
					}

					item := models.Item{
						ServerPath: strings.Replace(event.Name, d.VaultPath, "", 1),
						Size:       int(fileInfo.Size()),
					}

					item.SHA256 = d.CalculateSHA256(item)
					slog.Debug("File changed", "item", item)
					changeChan <- &item
					// d.Enqueue([]models.Item{item})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				slog.Error("error:", err)
			}
		}
	}()

	slog.Info("Watching path", "path", d.VaultPath)

	err = d.addRecursiveWatchers(watcher, d.VaultPath, false)
	if err != nil {
		return err
	}
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}

// From https://github.com/farmergreg/rfsnotify/blob/master/rfsnotify.go
func (d *LocalDriver) addRecursiveWatchers(watcher *fsnotify.Watcher, path string, unWatch bool) error {
	err := filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if unWatch {
				if err = watcher.Remove(walkPath); err != nil {
					return err
				}
			} else {
				if err = watcher.Add(walkPath); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
