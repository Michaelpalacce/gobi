package v1_syncstrategies

import (
	"fmt"
	"log/slog"

	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
)

// DefaultSyncStrategy is the default sync strategy
// Resolution is done by accepting the latest version of the file based on mtime
type DefaultSyncStrategy struct {
	driver storage.Driver
}

func NewDefaultSyncStrategy(driver storage.Driver) *DefaultSyncStrategy {
	return &DefaultSyncStrategy{
		driver: driver,
	}
}

func (s *DefaultSyncStrategy) SendSingle(client *socket.WebsocketClient, item models.Item) error {
	slog.Debug("Sending item", "item", item)

	if err := client.SendItem(item); err != nil {
		return err
	}

	slog.Debug("Item Sent Successfully", "item", item)

	return nil
}

// Fetch will download all non-conflicting items
func (s DefaultSyncStrategy) Fetch(client *socket.WebsocketClient) error {
	return s.FetchMultiple(client, s.driver.GetAllItems(storage.ConflictModeNo), storage.ConflictModeNo)
}

// FetchConflicts will download all conflicting items
func (s *DefaultSyncStrategy) FetchConflicts(client *socket.WebsocketClient) error {
	return s.FetchMultiple(client, s.driver.GetAllItems(storage.ConflictModeYes), storage.ConflictModeYes)
}

// FetchMultiple will resolve multiple items, either by downloading all non-conflicting items, or all conflicting items
func (s *DefaultSyncStrategy) FetchMultiple(client *socket.WebsocketClient, items []models.Item, conflictMode bool) error {
	for _, item := range items {
		if err := s.FetchSingle(client, item, conflictMode); err != nil {
			return fmt.Errorf("error downloading file: %w", err)
		}
	}

	return nil
}

// FetchSingle will fetch a single item only, and write it using the storage driver
// Conflict resolution is done by accepting the latest version of the file based on mtime
func (s *DefaultSyncStrategy) FetchSingle(client *socket.WebsocketClient, item models.Item, conflictMode bool) error {
	slog.Debug("Fetching item from server", "item", item)

	if conflictMode && item.ServerMTime < client.StorageDriver.GetMTime(item) {
		slog.Debug("Skipping conflict", "item", item)
		return nil
	}

	client.SendMessage(v1.NewItemFetchMessage(item))

	return client.FetchItem(item)
}
