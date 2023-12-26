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

// ProcessClientTextMessage will decide how to process the text message.
func ProcessClientTextMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	switch websocketMessage.Type {
	case v1.ItemsSyncType:
		if err := processItemSyncMessage(websocketMessage, client); err != nil {
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

// processItemSyncMessage adds items to the queue
// Check if sha256 matches locally
// Request File if it does not.
// If a file is not sent back in 30 seconds, close the connection
func processItemSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var itemsSyncPayload v1.ItemsSyncPayload

	client.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer client.Conn.SetReadDeadline(time.Time{})

	if err := json.Unmarshal(websocketMessage.Payload, &itemsSyncPayload); err != nil {
		return err
	}

	client.StorageDriver.Enqueue(itemsSyncPayload.Items)
	for client.StorageDriver.HasItemsToProcess() {
		item := client.StorageDriver.GetNext()
		slog.Debug("Fetching file from server", "item", item)

		client.SendMessage(v1.NewItemFetchMessage(*item))

		writer, err := client.StorageDriver.GetWriter(*item)
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
	}
	// After the queue is empty, check for any local changes and send them to the server
	client.StorageDriver.EnqueueItemsSince(client.Client.LastSync, client.Client.VaultName)

	client.Client.LastSync = int(time.Now().Unix())

	return nil
}

func ProcessClientBinaryMessage(message []byte, client *socket.WebsocketClient) error {
	return nil
}
