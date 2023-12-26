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
type Driver interface {
	Enqueue(items []models.Item)

	HasItemsToProcess() bool

	GetNext() *models.Item

	EnqueueItemsSince(lastSyncTime int, vaultName string)

	GetReader(i models.Item) (io.ReadCloser, error)

	GetWriter(i models.Item) (io.WriteCloser, error)

	Exists(i models.Item) bool

	CalculateSHA256(i models.Item) string
}
