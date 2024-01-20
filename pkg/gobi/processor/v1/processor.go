package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/gorilla/websocket"
)

// ProcessServerBinaryMessage will decide how to process the binary message.
func ProcessServerBinaryMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	return nil
}

// ProcessServerTextMessage will decide how to process the text message.
func ProcessServerTextMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	if client.Client.Version == 0 {
		return fmt.Errorf("before communications can happen, client must send %s message to specify version to use for responses", messages.VersionType)
	}

	switch websocketMessage.Type {
	case v1.VaultNameType:
		if err := processVaultNameMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.SyncType:
		if err := processSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.InitialSyncType:
		if err := processInitialSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.InitialSyncDoneType:
		if err := processInitialSyncDoneMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.ItemFetchType:
		if err := processItemFetchMessage(websocketMessage, client); err != nil {
			return err
		}
	case v1.ItemSaveType:
		if err := processItemSaveMessage(websocketMessage, client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processVaultNameMessage will set the VaultName in the client if when it's sent
func processVaultNameMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var vaultNamePayload v1.VaultNamePayload

	if err := json.Unmarshal(websocketMessage.Payload, &vaultNamePayload); err != nil {
		return err
	} else {
		client.Client.VaultName = vaultNamePayload.VaultName
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
	// @TODO: Check me
	for client.StorageDriver.HasItemsToProcess(false) {
		item := client.StorageDriver.GetNext(false)
		slog.Debug("Fetching file from client", "item", item)

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

	client.SendMessage(v1.NewInitialSyncDoneMessage(client.Client.LastSync))

	slog.Info("Fully synced")

	return nil
}

func processInitialSyncDoneMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var initialSyncDonePayload v1.InitialSyncDonePayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncDonePayload); err != nil {
		return err
	}

	client.InitialSync = true
	// This is just for info
	client.Client.LastSync = initialSyncDonePayload.LastSync

	client.SendMessage(v1.NewSyncMessage(initialSyncDonePayload.LastSync))

	slog.Info("Initial Client Sync Done", "vaultName", client.Client.VaultName)

	return nil
}

// processSyncMessage will start sending data to the client that needs to be synced up
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

		items := client.StorageDriver.GetAllItems(false)

		slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync)

		client.SendMessage(v1.NewInitialSyncMessage(items))
	}

	return nil
}

func processItemSaveMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var (
		itemSavePayload v1.ItemSavePayload
		err             error
	)

	client.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer client.Conn.SetReadDeadline(time.Time{})

	if err = json.Unmarshal(websocketMessage.Payload, &itemSavePayload); err != nil {
		return err
	}
	item := itemSavePayload.Item

	if item.SHA256 == client.StorageDriver.CalculateSHA256(item) {
		// @TODO: Touch the file, update the mtime
		// Do the same for the client
		// Other clients should see the change and update their mtime
		slog.Debug("Item already exists locally", "item", item)
		return nil
	}

	slog.Debug("Fetching file from client", "item", item)

	client.SendMessage(v1.NewItemFetchMessage(item))

	writer, err := client.StorageDriver.GetWriter(item)
	if err != nil {
		return err
	}

	defer func() {
		writer.Close()
	}()

	if item.Size == 0 {
		slog.Debug("File is empty", "item", item)
		return nil
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

	return nil
}

// processItemFetchMessage will start sending data to the client about the requested file
func processItemFetchMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var itemFetchPayload v1.ItemFetchPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemFetchPayload); err != nil {
		return err
	}

	if err := client.SendItem(itemFetchPayload.Item); err != nil {
		return err
	}

	return nil
}
