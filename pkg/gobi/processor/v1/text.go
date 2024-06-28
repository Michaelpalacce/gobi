package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/storage"
)

// ProcessServerTextMessage will decide how to process the text message.
func (p *Processor) ProcessServerTextMessage(websocketMessage messages.WebsocketMessage) error {
	if p.WebsocketClient.Client.Version == 0 {
		return fmt.Errorf("before communications can happen, client must send %s message to specify version to use for responses", messages.VersionType)
	}

	switch websocketMessage.Type {
	// The client tells us what the vault name is
	case v1.VaultNameType:
		if err := p.processVaultNameMessage(websocketMessage); err != nil {
			return err
		}
		// The client tells us what the sync strategy is
	case v1.SyncStrategyType:
		if err := p.processSyncStrategyMessage(websocketMessage); err != nil {
			return err
		}
	case v1.SyncType:
		if err := p.processSyncMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

func (p *Processor) processSyncStrategyMessage(websocketMessage messages.WebsocketMessage) error {
	var syncStrategyPayload v1.SyncStrategyPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncStrategyPayload); err != nil {
		return err
	}

	// switch syncStrategyPayload.SyncStrategy {
	// default:
	// 	return fmt.Errorf("unknown sync strategy: %d", syncStrategyPayload.SyncStrategy)
	// }
	p.UpdateSession()

	return nil
}

// processVaultNameMessage will set the VaultName in the client if when it's sent
// This is also when the Storage Driver is created
func (p *Processor) processVaultNameMessage(websocketMessage messages.WebsocketMessage) error {
	var vaultNamePayload v1.VaultNamePayload

	if err := json.Unmarshal(websocketMessage.Payload, &vaultNamePayload); err != nil {
		return err
	}

	p.WebsocketClient.Client.VaultName = vaultNamePayload.VaultName
	storageDriver, err := storage.NewLocalDriver(vaultNamePayload.VaultName)
	if err != nil {
		return err
	}

	p.WebsocketClient.StorageDriver = storageDriver
	p.UpdateSession()

	return nil
}

// processSyncMessage will enqueue items since the last sync and send the metadata to the client
func (p *Processor) processSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var syncPayload v1.SyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncPayload); err != nil {
		return err
	}

	p.WebsocketClient.StorageDriver.EnqueueItemsSince(
		syncPayload.LastSync,
		p.WebsocketClient.Client.VaultName,
	)

	items := p.WebsocketClient.StorageDriver.GetAllItems(storage.ConflictModeNo)

	slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync, "vaultName", p.WebsocketClient.Client.VaultName)

	// @TODO: send items to the client in a new message format
	// p.WebsocketClient.SendMessage()

	return nil
}
