package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// ------------------------------ Vault Name ------------------------------

type VaultNamePayload struct {
	VaultName string `json:"name"`
}

func NewVaultNameMessage(vaultName string) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: VaultNameType,
		Payload: VaultNamePayload{
			VaultName: vaultName,
		},
		Version: 1,
	}
}

// ------------------------------ Sync ------------------------------

type SyncPayload struct {
	// LastSync is timestamp in UTC
	LastSync int `json:"lastSync"`
}

func NewSyncMessage(lastSync int) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: SyncType,
		Payload: SyncPayload{
			LastSync: lastSync,
		},
		Version: 1,
	}
}

// ------------------------------ Initial Sync ------------------------------

type InitialSyncPayload struct {
	Items []models.Item `json:"items"`
}

func NewInitialSyncMessage(items []models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    InitialSyncType,
		Version: 1,
		Payload: InitialSyncPayload{
			Items: items,
		},
	}
}

// ------------------------------ Item Fetch ------------------------------

type ItemFetchPayload struct {
	Item models.Item `json:"item"`
}

func NewItemFetchMessage(item models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemFetchType,
		Version: 1,
		Payload: ItemFetchPayload{
			Item: item,
		},
	}
}

// ------------------------------ Item Save ------------------------------

type ItemSavePayload struct {
	Item models.Item `json:"item"`
}

func NewItemSavePayload(item models.Item) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    ItemSaveType,
		Version: 1,
		Payload: ItemSavePayload{
			Item: item,
		},
	}
}

// ------------------------------ Initial Sync Done ------------------------------

type InitialSyncDonePayload struct {
	// LastSync is timestamp in UTC
	LastSync int `json:"lastSync"`
}

func NewInitialSyncDoneMessage(lastSync int) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type:    InitialSyncDoneType,
		Version: 1,
		Payload: InitialSyncDonePayload{
			LastSync: lastSync,
		},
	}
}
