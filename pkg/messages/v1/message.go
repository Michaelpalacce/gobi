package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// VaultNamePayload stores the name of the vault that the client wants to connect to
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

// SyncPayload stores the last time the client synced with the server
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

// ItemsSyncPayload stores the items that have changed since the last sync
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

// ItemFetchPayload stores the item that the client wants to get
type ItemFetchPayload struct {
	Item models.Item `json:"item"`
}

// NewItemFetchMessage is a message that notifies the server that the client wants to get a specific file
func NewItemFetchMessage(item models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemFetchType,
		Version: 1,
		Payload: ItemFetchPayload{
			Item: item,
		},
	}
}

// ItemSavePayload stores the item that the client wants to save to the server
type ItemSavePayload struct {
	Item models.Item `json:"item"`
}

// NewItemSavePayload is a message that notifies the server that the client wants to save a specific file
func NewItemSavePayload(item models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemSaveType,
		Version: 1,
		Payload: ItemSavePayload{
			Item: item,
		},
	}
}
