package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
)

// Processor is the processor for version 1 of the protocol
// It contains the business logic for the protocol
type Processor struct {
	WebsocketClient *socket.WebsocketClient
}

// NewProcessor will create a new processor with the selected sync strategy in the client
func NewProcessor(client *socket.WebsocketClient) *Processor {
	switch client.Client.SyncStrategy {
	default:
		slog.Info("Using LastModifiedTimeSyncStrategy")
	}

	return &Processor{
		WebsocketClient: client,
	}
}

// ProcessClientTextMessage will decide how to process the text message.
func (p *Processor) ProcessClientTextMessage(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	// Server wants to perform initial sync
	case v1.InitialSyncType:
		if err := p.processInitialSyncMessage(websocketMessage); err != nil {
			return err
		}
		// Called when the server is done with the initial sync
	case v1.InitialSyncDoneType:
		if err := p.processInitialSyncDoneMessage(websocketMessage); err != nil {
			return err
		}
		// Called when the server wants to sync
	case v1.SyncType:
		if err := p.processSyncMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processInitialSyncMessage is received from the server telling the client witch items have changed since the last sync
// The client will do REST calls to fetch the items from the server according to the sync strategy
func (p *Processor) processInitialSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncPayload v1.InitialSyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncPayload); err != nil {
		return err
	}

	// @TODO: Fetch items from the server

	p.WebsocketClient.SendMessage(v1.NewInitialSyncDoneMessage(p.WebsocketClient.Client.LastSync))

	return nil
}

// processInitialSyncDoneMessage notifies the client that the server is done with the initial sync
func (p *Processor) processInitialSyncDoneMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncDonePayload v1.InitialSyncDonePayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncDonePayload); err != nil {
		return err
	}

	// TODO: This needs to be persisted, in the same place we retrieved it from
	p.WebsocketClient.Client.LastSync = initialSyncDonePayload.LastSync

	slog.Info("Initial Server Sync Done", "vaultName", p.WebsocketClient.Client.VaultName)
	slog.Info("Fully synced")

	// This needs to be persisted somehow
	p.WebsocketClient.Client.LastSync = int(time.Now().Unix())

	if p.WebsocketClient.InitialSync {
		p.WebsocketClient.InitialSync = false

		changeChan := make(chan *models.Item)

		go p.WebsocketClient.StorageDriver.WatchVault(p.WebsocketClient.Client.VaultName, changeChan)
		slog.Info("Starting to watch vault", "vaultName", p.WebsocketClient.Client.VaultName)

		go func() {
			// @TODO: send changes to the server
			// for item := range changeChan {
			// }
		}()
	}

	return nil
}

// processSyncMessage is a request from the server that it wants to sync.
// The client will send all items that have been modified since the last sync (provided by the server)
// @NOTE: Should this use the sync strategy?
func (p *Processor) processSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var syncPayload v1.SyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncPayload); err != nil {
		return err
	}

	p.WebsocketClient.StorageDriver.EnqueueItemsSince(
		syncPayload.LastSync,
		p.WebsocketClient.Client.VaultName,
	)

	items := p.WebsocketClient.StorageDriver.GetAllItems(storage.ConflictModeNo)

	slog.Debug("Items found for sync since last reconcillation", "items", len(items), "lastSync", syncPayload.LastSync)
	p.WebsocketClient.SendMessage(v1.NewInitialSyncMessage(items))

	return nil
}

func (p *Processor) ProcessClientBinaryMessage(message []byte) error {
	return fmt.Errorf("binary messages are not supported for version 1")
}
