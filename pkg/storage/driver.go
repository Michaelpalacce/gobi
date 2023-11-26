package storage

import (
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// Event holds information about a file operation.
// Events could be deletes,updates,creations
type Event struct {
}

// Item represents metadata about an item
type Item struct {
	Item models.Item
}

// Driver interface holds the structure that all storage drivers must adhere to
type Driver interface {
	// PushFile will push the given file to the server and set it as the latest reconciled version
	PushFile(f Item) error

	// PullFile will fetch a file from the server and have it sent out
	PullFile(filePath string) (Item, error)
}
