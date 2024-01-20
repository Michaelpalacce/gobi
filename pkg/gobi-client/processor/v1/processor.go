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
	v1_syncstrategies "github.com/Michaelpalacce/gobi/pkg/syncStrategies/v1"
)

// ProcessClientTextMessage will decide how to process the text message.
func ProcessClientTextMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	switch websocketMessage.Type {
	// Called only once at the beginning
	case v1.InitialSyncType:
		if err := processInitialSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.InitialSyncDoneType:
		if err := processInitialSyncDoneMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.SyncType:
		if err := processSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.ItemFetchType:
		if err := processItemFetchMessage(websocketMessage, client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processInitialSyncDoneMessage notifies the client that the server is done with the initial sync
func processInitialSyncDoneMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var initialSyncDonePayload v1.InitialSyncDonePayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncDonePayload); err != nil {
		return err
	}

	client.InitialSync = true
	// TODO: This needs to be persisted
	client.Client.LastSync = initialSyncDonePayload.LastSync

	slog.Info("Initial Server Sync Done", "vaultName", client.Client.VaultName)
	slog.Info("Fully synced")

	// This needs to be persisted somehow
	client.Client.LastSync = int(time.Now().Unix())

	if client.InitialSync {
		client.InitialSync = false

		changeChan := make(chan *models.Item)

		go client.WatchVault(changeChan)
		slog.Info("Starting to watch vault", "vaultName", client.Client.VaultName)

		go func() {
			for item := range changeChan {
				client.SendMessage(v1.NewItemSavePayload(*item))
			}
		}()
	}

	return nil
}

// processItemFetchMessage will start sending the file to the server
func processItemFetchMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var itemFetchPayload v1.ItemFetchPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemFetchPayload); err != nil {
		return err
	}

	syncStrategy := getSyncStrategy(client)

	if err := syncStrategy.SendSingle(client, itemFetchPayload.Item); err != nil {
		return err
	}

	return nil
}

// processSyncMessage is a request from the server that it wants to sync.
// The client will send all items that have been modified since the last sync (provided by the server)
// @NOTE: Should this use the sync strategy?
func processSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var syncPayload v1.SyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncPayload); err != nil {
		return err
	} else {
		client.StorageDriver.EnqueueItemsSince(
			syncPayload.LastSync,
			client.Client.VaultName,
		)

		if err != nil {
			return err
		}

		items := client.StorageDriver.GetAllItems(storage.ConflictModeNo)

		slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync)

		client.SendMessage(v1.NewInitialSyncMessage(items))
	}

	return nil
}

// processInitialSyncMessage takes the list of items from the server and compares them to the local vault
// Check if sha256 matches locally
// Request File if it does not.
// If a file is not sent back in 30 seconds, close the connection
// This is done only once
func processInitialSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var initialSyncPayload v1.InitialSyncPayload

	client.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer client.Conn.SetReadDeadline(time.Time{})

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncPayload); err != nil {
		return err
	}

	client.StorageDriver.Enqueue(initialSyncPayload.Items)

	syncStrategy := getSyncStrategy(client)

	if err := syncStrategy.Fetch(client); err != nil {
		return err
	}
	if err := syncStrategy.FetchConflicts(client); err != nil {
		return err
	}

	client.SendMessage(v1.NewInitialSyncDoneMessage(client.Client.LastSync))

	return nil
}

func ProcessClientBinaryMessage(message []byte, client *socket.WebsocketClient) error {
	return nil
}

// @TODO: This should be configurable, but for now, this is the only sync strategy
func getSyncStrategy(client *socket.WebsocketClient) v1_syncstrategies.SyncStrategy {
	return v1_syncstrategies.NewDefaultSyncStrategy(client.StorageDriver)
}
