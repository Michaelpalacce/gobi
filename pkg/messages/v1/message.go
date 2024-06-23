package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/messages"
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
		Version: Version,
	}
}

// ------------------------------ Sync Strategy ------------------------------

type SyncStrategyPayload struct {
	SyncStrategy int `json:"syncStrategy"`
}

func NewSyncStrategyMessage(syncStrategy int) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: SyncStrategyType,
		Payload: SyncStrategyPayload{
			SyncStrategy: syncStrategy,
		},
		Version: Version,
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
		Version: Version,
	}
}
