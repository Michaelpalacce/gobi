package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
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
		slog.Debug("Client sync attempt", "LastSync", syncPayload.LastSync)
	}

	return nil
}
