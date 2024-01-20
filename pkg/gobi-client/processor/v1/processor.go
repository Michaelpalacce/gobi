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
	syncstrategies "github.com/Michaelpalacce/gobi/pkg/syncStrategies"
)

type Processor struct {
	SyncStrategy syncstrategies.SyncStrategy
	Client       *socket.WebsocketClient
}

// NewProcessor will create a new processor with the selected sync strategy in the client
func NewProcessor(client *socket.WebsocketClient) *Processor {
	var syncStrategy syncstrategies.SyncStrategy

	switch client.Client.SyncStrategy {
	case syncstrategies.LastModifiedTimeStrategy:
		fallthrough
	default:
		slog.Info("Using LastModifiedTimeSyncStrategy")
		syncStrategy = syncstrategies.NewLastModifiedTimeSyncStrategy(client.StorageDriver, client)
	}

	return &Processor{
		Client:       client,
		SyncStrategy: syncStrategy,
	}
}

// ProcessClientTextMessage will decide how to process the text message.
func (p *Processor) ProcessClientTextMessage(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	// Called only once at the beginning
	case v1.InitialSyncType:
		if err := p.processInitialSyncMessage(websocketMessage); err != nil {
			return err
		}
	case v1.InitialSyncDoneType:
		if err := p.processInitialSyncDoneMessage(websocketMessage); err != nil {
			return err
		}
	case v1.SyncType:
		if err := p.processSyncMessage(websocketMessage); err != nil {
			return err
		}
	case v1.ItemFetchType:
		if err := p.processItemFetchMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processInitialSyncDoneMessage notifies the client that the server is done with the initial sync
func (p *Processor) processInitialSyncDoneMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncDonePayload v1.InitialSyncDonePayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncDonePayload); err != nil {
		return err
	}

	p.Client.InitialSync = true
	// TODO: This needs to be persisted, in the same place we retrieved it from
	p.Client.Client.LastSync = initialSyncDonePayload.LastSync

	slog.Info("Initial Server Sync Done", "vaultName", p.Client.Client.VaultName)
	slog.Info("Fully synced")

	// This needs to be persisted somehow
	p.Client.Client.LastSync = int(time.Now().Unix())

	if p.Client.InitialSync {
		p.Client.InitialSync = false

		changeChan := make(chan *models.Item)

		go p.Client.WatchVault(changeChan)
		slog.Info("Starting to watch vault", "vaultName", p.Client.Client.VaultName)

		go func() {
			for item := range changeChan {
				p.Client.SendMessage(v1.NewItemSavePayload(*item))
			}
		}()
	}

	return nil
}

// processItemFetchMessage will start sending the file to the server
func (p *Processor) processItemFetchMessage(websocketMessage messages.WebsocketMessage) error {
	var itemFetchPayload v1.ItemFetchPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemFetchPayload); err != nil {
		return err
	}

	syncStrategy := p.SyncStrategy

	if err := syncStrategy.SendSingle(itemFetchPayload.Item); err != nil {
		return err
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
	} else {
		p.Client.StorageDriver.EnqueueItemsSince(
			syncPayload.LastSync,
			p.Client.Client.VaultName,
		)

		if err != nil {
			return err
		}

		items := p.Client.StorageDriver.GetAllItems(storage.ConflictModeNo)

		slog.Debug("Items found for sync since last reconcillation", "items", len(items), "lastSync", syncPayload.LastSync)

		p.Client.SendMessage(v1.NewInitialSyncMessage(items))
	}

	return nil
}

// processInitialSyncMessage takes the list of items from the server and compares them to the local vault
// Check if sha256 matches locally
// Request File if it does not.
// If a file is not sent back in 30 seconds, close the connection
// This is done only once
func (p *Processor) processInitialSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncPayload v1.InitialSyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncPayload); err != nil {
		return err
	}

	p.Client.StorageDriver.Enqueue(initialSyncPayload.Items)

	syncStrategy := p.SyncStrategy

	if err := syncStrategy.Fetch(); err != nil {
		return err
	}
	if err := syncStrategy.FetchConflicts(); err != nil {
		return err
	}

	p.Client.SendMessage(v1.NewInitialSyncDoneMessage(p.Client.Client.LastSync))

	return nil
}

func (p *Processor) ProcessClientBinaryMessage(message []byte) error {
	return nil
}
