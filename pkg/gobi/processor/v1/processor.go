package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage/metadata"
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

// processSyncMessage will start sending data to the client that needs to be synced up
func processSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var syncPayload v1.SyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncPayload); err != nil {
		return err
	} else {
		mongoDriver := metadata.MongoDriver{
			DB:     client.DB,
			Client: &client.Client,
		}

		items, err := mongoDriver.Reconcile(syncPayload.LastSync)
		if err != nil {
			return err
		}

		slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync)

		client.SendMessage(v1.NewItemsSyncMessage(items))
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

	slog.Debug("Fetching file from server", "item", item)

	client.SendMessage(v1.NewItemFetchMessage(item))

	writer, err := client.StorageDriver.GetWriter(item)
	if err != nil {
		return err
	}

	defer func() {
		writer.Close()
	}()

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

	itemsService := services.NewItemsService(client.DB)

	err = itemsService.Upsert(&item)
	if err != nil {
		return err
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
