package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	"github.com/Michaelpalacce/gobi/pkg/storage/metadata"
	"github.com/gorilla/websocket"
)

// ProcessServerTextMessage will decide how to process the text message.
func ProcessServerTextMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	if client.Client.Version == 0 {
		return fmt.Errorf("before communications can happen, client must send %s message to specify version to use for responses", messages.VersionType)
	}

	switch websocketMessage.Type {
	case VaultNameType:
		if err := processVaultNameMessage(websocketMessage, client); err != nil {
			return err
		}
	case SyncType:
		if err := processSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processVaultNameMessage will set the VaultName in the client if when it's sent
func processVaultNameMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	var vaultNamePayload VaultNamePayload

	if err := json.Unmarshal(websocketMessage.Payload, &vaultNamePayload); err != nil {
		return err
	} else {
		client.Client.VaultName = vaultNamePayload.VaultName
	}

	return nil
}

// processSyncMessage will start sending data to the client that needs to be synced up
func processSyncMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	var syncPayload SyncPayload

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

		slog.Debug("Items Found For Sync", "items", items)
		for _, item := range items {
			client.SendMessage(NewItemSyncMessage(item.Item))
			sendBigFile(client, item)
		}
	}

	return nil
}

// sendBigFile will send an item to the client
func sendBigFile(client *client.WebsocketClient, item storage.Item) error {
	file, err := os.Open(item.Item.ServerPath)
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
