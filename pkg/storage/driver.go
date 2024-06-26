package storage

import (
	"io"

	"github.com/Michaelpalacce/gobi/pkg/models"
)

// Event holds information about a file operation.
// Events could be deletes,updates,creations
type Event struct{}

// Driver interface holds the structure that all storage drivers must adhere to
// Storage Drivers are responsible for storing what needs to be pushed/pulled and doing requests to sync what is needed
// Storage Drivers are also responsible for handling the actual file operations
// Conflicts are files changed on both the server and the client
type Driver interface {
	Enqueue(items []models.Item)

	HasItemsToProcess(conflictMode bool) bool

	GetAllItems(conflictMode bool) []models.Item

	GetMTime(i models.Item) int64

	GetNext(conflictMode bool) *models.Item

	EnqueueItemsSince(lastSyncTime int, vaultName string)

	GetReader(i models.Item) (io.ReadCloser, error)

	GetWriter(i models.Item) (io.WriteCloser, error)

	Exists(i models.Item) bool

	Touch(i models.Item) error

	CalculateSHA256(i models.Item) string

	WatchVault(vaultName string, changeChan chan<- *models.Item) error
}

const (
	ConflictModeNo  bool = false
	ConflictModeYes bool = true
)
