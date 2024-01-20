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
	"github.com/gorilla/websocket"
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
			for {
				select {
				case item := <-changeChan:
					client.SendMessage(v1.NewItemSavePayload(*item))
				}
			}
		}()
	}

	return nil
}

// processItemFetchMessage will start sending data to the client about the requested file
func processItemFetchMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var itemFetchPayload v1.ItemFetchPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemFetchPayload); err != nil {
		return err
	}

	item, err := client.StorageDriver.GetReader(itemFetchPayload.Item)
	if err != nil {
		return err
	}

	defer item.Close()

	if err := client.SendItem(item); err != nil {
		return err
	}

	return nil
}

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

		items := client.StorageDriver.GetAllItems()

		slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync)

		client.SendMessage(v1.NewInitialSyncMessage(items))
	}

	return nil
}

// processInitialSyncMessage adds items to the queue
// Check if sha256 matches locally
// Request File if it does not.
// If a file is not sent back in 30 seconds, close the connection
// This is done only once, after the initial sync, the client will watch the vault for changes
func processInitialSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var initialSyncPayload v1.InitialSyncPayload

	client.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer client.Conn.SetReadDeadline(time.Time{})

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncPayload); err != nil {
		return err
	}

	client.StorageDriver.Enqueue(initialSyncPayload.Items)
	if err := handleStorageItems(client, false); err != nil {
		return err
	}

	client.StorageDriver.EnqueueConflcits()
	if err := handleStorageItems(client, true); err != nil {
		return err
	}

	client.SendMessage(v1.NewInitialSyncDoneMessage(client.Client.LastSync))

	return nil
}

func handleStorageItems(client *socket.WebsocketClient, conflictMode bool) error {
processor:
	for client.StorageDriver.HasItemsToProcess() {
		item := client.StorageDriver.GetNext()

		if conflictMode {
			if item.ServerMTime < client.StorageDriver.GetMTime(*item) {
				slog.Debug("Skipping conflict", "item", item)
				continue processor
			}
		}

		slog.Debug("Fetching file from server", "item", item)

		client.SendMessage(v1.NewItemFetchMessage(*item))

		writer, err := client.StorageDriver.GetWriter(*item)
		if err != nil {
			return err
		}

		defer func() {
			writer.Close()
		}()

		if item.Size == 0 {
			slog.Debug("File Fetched Successfully", "item", item)
			continue
		}

		bytesRead := 0
		for {
			messageType, message, err := client.Conn.ReadMessage()
			if err != nil {
				return err
			}

			if messageType != websocket.BinaryMessage {
				return fmt.Errorf("invalid messageType received: %d, expected 2 (BinaryMessage)", messageType)
			}

			writer.Write(message)

			bytesRead += len(message)
			if bytesRead == item.Size {
				writer.Close()
				break
			}

			if bytesRead > item.Size {
				return fmt.Errorf("expected %d bytes, but got %d", item.Size, bytesRead)
			}
		}
		slog.Debug("File Fetched Successfully", "item", item)
	}

	return nil
}

func ProcessClientBinaryMessage(message []byte, client *socket.WebsocketClient) error {
	return nil
}
