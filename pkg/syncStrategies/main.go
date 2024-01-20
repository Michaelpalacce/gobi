package syncstrategies

import (
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// SyncStrategy is the interface that all sync strategies must implement
// Used to fetch and send data to the server
// The sync strategies are responsible for ensuring thread safety,
// so only one sync strategy should be used per client
type SyncStrategy interface {
	SendSingle(item models.Item) error
	FetchSingle(item models.Item, conflictMode bool) error

	Fetch() error
	FetchConflicts() error
	FetchMultiple(items []models.Item, conflictMode bool) error
}
