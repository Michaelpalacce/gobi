package syncstrategies

import (
	"fmt"
	"log/slog"

	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
)

// LastModifiedTimeSyncStrategy is the default sync strategy
// Resolution is done by accepting the latest version of the file based on mtime
// TODO: Make this a struct with a mutex so that it can be used concurrently
type LastModifiedTimeSyncStrategy struct {
	Driver storage.Driver
	Client *socket.WebsocketClient
}

func NewLastModifiedTimeSyncStrategy(driver storage.Driver, client *socket.WebsocketClient) *LastModifiedTimeSyncStrategy {
	return &LastModifiedTimeSyncStrategy{
		Driver: driver,
		Client: client,
	}
}

// SendSingle will send a single item to the remote
func (s *LastModifiedTimeSyncStrategy) SendSingle(item models.Item) error {
	slog.Debug("Sending item", "item", item)

	if err := s.Client.SendItem(item); err != nil {
		return err
	}

	slog.Debug("Item Sent Successfully", "item", item)

	return nil
}

// Fetch will download all non-conflicting items
func (s LastModifiedTimeSyncStrategy) Fetch() error {
	return s.FetchMultiple(s.Driver.GetAllItems(storage.ConflictModeNo), storage.ConflictModeNo)
}

// FetchConflicts will download all conflicting items
func (s *LastModifiedTimeSyncStrategy) FetchConflicts() error {
	return s.FetchMultiple(s.Driver.GetAllItems(storage.ConflictModeYes), storage.ConflictModeYes)
}

// FetchMultiple will resolve multiple items, either by downloading all non-conflicting items, or all conflicting items
func (s *LastModifiedTimeSyncStrategy) FetchMultiple(items []models.Item, conflictMode bool) error {
	for _, item := range items {
		if err := s.FetchSingle(item, conflictMode); err != nil {
			return fmt.Errorf("error downloading file: %w", err)
		}
	}

	return nil
}

// FetchSingle will fetch a single item only, and write it using the storage driver
// Conflict resolution is done by accepting the latest version of the file based on mtime
// NOTE: This is where the conflict resolution happens
func (s *LastModifiedTimeSyncStrategy) FetchSingle(item models.Item, conflictMode bool) error {
	slog.Debug("Fetching item from server", "item", item)

	if conflictMode && item.ServerMTime < s.Client.StorageDriver.GetMTime(item) {
		slog.Debug("Skipping conflict", "item", item)
		return nil
	}

	s.Client.SendMessage(v1.NewItemFetchMessage(item))

	return s.Client.FetchItem(item)
}
