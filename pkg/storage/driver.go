package storage

import (
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// Event holds information about a file operation.
// Events could be deletes,updates,creations
type Event struct {
}

// Driver interface holds the structure that all storage drivers must adhere to
// Storage Drivers are responsible for storing what needs to be pushed/pulled and doing requests to sync what is needed
// Storage Drivers
type Driver interface {
	CheckIfLocalMatch(i models.Item) bool
}
