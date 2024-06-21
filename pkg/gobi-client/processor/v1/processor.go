package processor_v1

import (
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/gobi-client/settings"
	"github.com/Michaelpalacce/gobi/pkg/socket"
)

// Processor is the processor for version 1 of the protocol
// It contains the business logic for the protocol
type Processor struct {
	WebsocketClient *socket.WebsocketClient
	LocalSettings   *settings.Store
}

// NewProcessor will create a new processor with the selected sync strategy in the client
func NewProcessor(client *socket.WebsocketClient, localSettings *settings.Store) *Processor {
	switch client.Client.SyncStrategy {
	default:
		slog.Info("Using LastModifiedTimeSyncStrategy")
	}

	return &Processor{
		WebsocketClient: client,
		LocalSettings:   localSettings,
	}
}
