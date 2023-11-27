package processor_v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
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

// sendBigFile will send an item to the client without storing bigger than 1024 chunks in memory
func sendBigFile(client *socket.WebsocketClient, item models.Item) error {
	file, err := os.Open(item.ServerPath)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading file: %s", err)
		}

		err = client.Conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			return fmt.Errorf("error reading file chunk: %s", err)
		}
	}

	return nil
}
