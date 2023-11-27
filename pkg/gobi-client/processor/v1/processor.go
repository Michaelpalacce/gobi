package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Michaelpalacce/gobi/pkg/iops"
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
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processItemSyncMessage adds items to the queue
// Check if sha256 matches locally
// Request File if it does not.
func processItemSyncMessage(websocketMessage messages.WebsocketMessage, client *socket.WebsocketClient) error {
	var itemsSyncPayload v1.ItemsSyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemsSyncPayload); err != nil {
		return err
	}

	client.StorageDriver.Enqueue(itemsSyncPayload.Items)
	for client.StorageDriver.HasItemsToProcess() {
		item := client.StorageDriver.GetNext()
		slog.Debug("Fetching file from server", "item", item)

		client.SendMessage(v1.NewItemFetchMessage(*item))

		// TODO: Move this to the driver
		var tempFile *os.File
		var err error
		if tempFile, err = os.CreateTemp("", "websocket_upload_"); err != nil {
			return fmt.Errorf("could not create a temp file: %s", err)
		}
		defer func() {
			tempFile.Close()
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

			tempFile.Write(message)

			bytesRead += len(message)
			fmt.Println(bytesRead)
			if bytesRead == item.Size {
				tempFile.Close()
				break
			}

			if bytesRead > item.Size {
				return fmt.Errorf("expected %d bytes, but got %d", item.Size, bytesRead)
			}
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %s", err)
		}
		filePath := filepath.Join(cwd, "./.dev/clientFolder", item.VaultName, item.ServerPath)
		slog.Debug("File Fetched Successfully", "item", item, "filePath", filePath)
		if err := iops.MoveFile(tempFile.Name(), filePath); err != nil {
			return err
		}
	}

	return nil
}

func ProcessClientBinaryMessage(message []byte, client *socket.WebsocketClient) error {
	return nil
}
