package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/messages/v1/rest"
	"github.com/Michaelpalacce/gobi/pkg/storage"
)

// ProcessClientTextMessage will decide how to process the text message.
func (p *Processor) ProcessClientTextMessage(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	// Called when the server wants to sync
	case v1.SyncType:
		if err := p.processSyncMessage(websocketMessage); err != nil {
			return err
		}
	case rest.SessionType:
		if err := p.processSessionMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// if p.WebsocketClient.InitialSync {
// 	p.WebsocketClient.InitialSync = false
//
// 	changeChan := make(chan *models.Item)
//
// 	go p.WebsocketClient.StorageDriver.WatchVault(p.WebsocketClient.Client.VaultName, changeChan)
// 	slog.Info("Starting to watch vault", "vaultName", p.WebsocketClient.Client.VaultName)
//
// 	go func() {
// 		// @TODO: send changes to the server
// 		// for item := range changeChan {
// 		// }
// 	}()
// }
//
// return nil

// processSyncMessage is a request from the server that it wants to sync.
// The client will send all items that have been modified since the last sync (provided by the server)
// @NOTE: Should this use the sync strategy?
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

	slog.Debug("Items found for sync since last reconcillation", "items", len(items), "lastSync", syncPayload.LastSync)
	// @TODO: Send new message telling the client the changed files
	// p.WebsocketClient.SendMessage()

	return nil
}

// processSessionMessage will process the session message from the server
func (p *Processor) processSessionMessage(websocketMessage messages.WebsocketMessage) error {
	var sessionPayload rest.SessionPayload

	if err := json.Unmarshal(websocketMessage.Payload, &sessionPayload); err != nil {
		return err
	}

	p.SessionID = sessionPayload.SessionId
	slog.Debug("Received session message", "sessionID", sessionPayload.SessionId)

	return nil
}
