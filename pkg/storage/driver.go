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
	CheckIfLocalMatch(i models.Item) bool
}
