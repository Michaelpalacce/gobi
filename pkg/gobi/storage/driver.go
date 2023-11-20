package storage

import (
	"time"
)

// Event holds information about a file operation.
// Events could be deletes,updates,creations
type Event struct {
}

type File struct{}

// Driver interface holds the structure that all storage drivers must adhere to
type Driver interface {
	// Reconcile will give a list of changes from the given time until time.Now
	Reconcile(lastSync time.Time) ([]Event, error)

	// PushFile will push the given file to the server and set it as the latest reconciled version
	PushFile(f File) error

	// PullFile will fetch a file from the server and have it sent out
	PullFile(filePath string) (File, error)
}
