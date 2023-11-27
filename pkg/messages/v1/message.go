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

type ItemSyncPayload struct {
	Item models.Item `json:"item"`
}

// NewItemSyncMessage cotains data about an item that has had a change since the last reconcillation time
func NewItemSyncMessage(item models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemSyncType,
		Version: 1,
		Payload: ItemSyncPayload{
			Item: item,
		},
	}
}
