package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

type VaultNamePayload struct {
	VaultName string `json:"name"`
}

// NewVaultNameMessage is a message that the client sends to the server telling it which Vault to connect to
func NewVaultNameMessage(vaultName string) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: VaultNameType,
		Payload: VaultNamePayload{
			VaultName: vaultName,
		},
		Version: 1,
	}
}

type SyncPayload struct {
	// LastSync is timestamp in UTC
	LastSync int `json:"lastSync"`
}

// NewSyncMessage creates a message to send to the server telling it when was the last time the client synced.
// The LastSync is a timestamp
func NewSyncMessage(lastSync int) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: SyncType,
		Payload: SyncPayload{
			LastSync: lastSync,
		},
		Version: 1,
	}
}

type ItemsSyncPayload struct {
	Items []models.Item `json:"items"`
}

// NewItemsSyncMessage contains data about items that have had a change since the last reconcillation time
func NewItemsSyncMessage(items []models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemsSyncType,
		Version: 1,
		Payload: ItemsSyncPayload{
			Items: items,
		},
	}
}

type TransferStartPayload struct{}

// NewTransferStartPayload is the beginning handshake that the client or server wants to either push or pull data
func NewTransferStartPayload() messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    TransferStartType,
		Version: 1,
		Payload: TransferStartPayload{},
	}
}
