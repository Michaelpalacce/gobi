package syncstrategies

import (
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// SyncStrategy is the interface that all sync strategies must implement
// Used to fetch and send data to the remote
// The sync strategies are responsible for ensuring thread safety,
// so only one sync strategy should be used per client
type SyncStrategy interface {
	// SendSingle will send a single item to the remote
	SendSingle(item models.Item) error

	// Fetch will fetch all non-conflicting items from the remote
	Fetch() error
	// FetchConflicts will fetch all conflicting items from the remote
	FetchConflicts() error
	// FetchSingle will fetch a single item from the remote
	FetchSingle(item models.Item, conflictMode bool) error
	// FetchMultiple will fetch multiple items from the remote
	FetchMultiple(items []models.Item, conflictMode bool) error
}
